package manager

import (
	cpmanager "coffee/app/campaignprofiles/manager"
	discoverymanager "coffee/app/discovery/manager"
	partnerusagemanager "coffee/app/partnerusage/manager"
	"encoding/json"
	"strings"

	log "github.com/sirupsen/logrus"

	"coffee/app/profilecollection/dao"
	"coffee/app/profilecollection/domain"
	"coffee/client/jobtracker"
	"coffee/constants"
	"coffee/core/appcontext"
	coredomain "coffee/core/domain"
	"coffee/core/rest"
	"coffee/helpers"
	"context"
	"errors"
	"strconv"

	"go.step.sm/crypto/randutil"
)

// -------------------------------
// profile collection service
// -------------------------------

type Manager struct {
	*rest.Manager[domain.ProfileCollectionEntry, dao.ProfileCollectionEntity, int64]
	dao                    *dao.Dao
	itemDAO                *dao.ItemDao
	winklCollectionInfoDAO *dao.WinklCollectionInfoDao
	socialSearchManager    *discoverymanager.SearchManager
	partnerUsageManager    *partnerusagemanager.PartnerUsageManager
}

func (m *Manager) init(ctx context.Context) {
	pcdao := dao.CreateDao(ctx)
	itemDao := dao.CreateItemDao(ctx)
	winklCollectionInfoDAO := dao.CreateWinklCollectionInfoDao(ctx)
	m.Manager = rest.NewManager(ctx, pcdao.Dao, ToEntry, ToEntity)
	m.dao = pcdao
	m.itemDAO = itemDao
	m.winklCollectionInfoDAO = winklCollectionInfoDAO
	socialSearchManager := discoverymanager.CreateSearchManager(ctx)
	m.socialSearchManager = socialSearchManager
	m.partnerUsageManager = partnerusagemanager.CreatePartnerUsageManager(ctx)

}

func CreateManager(ctx context.Context) *Manager {
	manager := &Manager{}
	manager.init(ctx)
	return manager
}

func (m *Manager) Create(ctx context.Context, input *domain.ProfileCollectionEntry) (*domain.ProfileCollectionEntry, error) {
	createdEntry, err := m.Manager.Create(ctx, input)
	if err != nil {
		// record exists
		appCtx := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)
		appCtx.Session.Close(ctx) // TODO: Figure out a better way to reset session. Need to do this because previous session has been marked as rollback-only
		appCtx.Session = nil
		errMsg := err.Error()
		var entry *domain.ProfileCollectionEntry
		if strings.Contains(errMsg, "profile_collection_pkey") {
			entry, _ = m.FindById(ctx, *input.Id, true)
		} else if strings.Contains(errMsg, "profile_collection_share_id_key") {
			entry, _ = m.FindByShareId(ctx, *input.ShareId)
		} else if strings.Contains(errMsg, "profile_collection_source_key") {
			entry, _ = m.FindBySourceAndSourceId(ctx, *input.Source, *input.SourceId)
		}
		if entry != nil {
			return entry, nil
		}
		return nil, err
	}
	if createdEntry == nil || createdEntry.Id == nil {
		return nil, errors.New("failed to create profile collection")
	}

	if *input.Source == string(constants.SaasAtCollection) {
		if input.ItemsAssetId != nil {
			entry, err2 := m.addItemsFromAssetId(ctx, input, createdEntry)
			if err2 != nil {
				return entry, err2
			}
		} else if input.Items != nil && len(input.Items) > 0 {
			_, err, _ := m.AddItemsToCollection(ctx, *createdEntry, input.Items)
			if err != nil {
				return nil, err
			}
		} else if input.ImportCollectionId != nil {
			entry, err2 := m.addItemsFromOtherCollection(ctx, input, createdEntry)
			if err2 != nil {
				return entry, err2
			}
		} else {
			return nil, errors.New("no Items for tracking")
		}
	}
	return createdEntry, err
}

func (m *Manager) Update(ctx context.Context, id int64, input *domain.ProfileCollectionEntry) (*domain.ProfileCollectionEntry, error) {
	updatedEntry, err := m.Manager.Update(ctx, id, input)
	if err != nil {
		return nil, err
	}
	if updatedEntry == nil || updatedEntry.Id == nil {
		return nil, errors.New("failed to update profile collection")
	}
	entry, err := m.Manager.FindById(ctx, id)
	if err == nil && entry != nil && (*entry.Source == string(constants.SaasAtCollection) || *entry.Source == string(constants.SaasCollection)) {
		if input.ItemsAssetId != nil {
			input.Name = entry.Name
			entry, err2 := m.addItemsFromAssetId(ctx, input, entry)
			if err2 != nil {
				return entry, err2
			}
		} else if input.Items != nil && len(input.Items) > 0 {
			_, err, _ := m.AddItemsToCollection(ctx, *entry, input.Items)
			if err != nil {
				return nil, err
			}
		} else if input.ImportCollectionId != nil {
			entry, err2 := m.addItemsFromOtherCollection(ctx, input, entry)
			if err2 != nil {
				return entry, err2
			}
		}
	}
	return updatedEntry, err
}

func (m *Manager) addItemsFromOtherCollection(ctx context.Context, input *domain.ProfileCollectionEntry, createdEntry *domain.ProfileCollectionEntry) (*domain.ProfileCollectionEntry, error) {
	jobEntry := jobtracker.JobEntry{}
	jobName := "Import Items From Collection - " + strconv.FormatInt(*input.ImportCollectionId, 10)
	jobType := "IMPORT_ITEMS_FROM_PROFILE_COLLECTION"
	jobEntry.JobType = &jobType
	jobEntry.JobName = &jobName
	jobEntry.Input = map[string]string{}
	jobEntry.Input["toCollectionId"] = strconv.FormatInt(*createdEntry.Id, 10)
	jobEntry.Input["fromCollectionId"] = strconv.FormatInt(*input.ImportCollectionId, 10)
	if input.Platform != nil {
		jobEntry.Input["platform"] = *input.Platform
	}
	client := jobtracker.New(ctx)
	jobResponse, err := client.CreateJob(jobEntry)
	if err != nil {
		return nil, err
	}
	if jobResponse == nil || jobResponse.Status.Type != "SUCCESS" || len(jobResponse.Data) == 0 || (jobResponse.Data)[0].Id == nil {
		return nil, errors.New("failed to create job for bulk item import")
	}
	jobId := (jobResponse.Data)[0].Id

	updateCollectionEntry := domain.ProfileCollectionEntry{JobId: jobId}
	m.Manager.Update(ctx, *createdEntry.Id, &updateCollectionEntry)
	createdEntry.JobId = jobId
	return nil, nil
}

func (m *Manager) addItemsFromAssetId(ctx context.Context, input *domain.ProfileCollectionEntry, createdEntry *domain.ProfileCollectionEntry) (*domain.ProfileCollectionEntry, error) {
	jobEntry := jobtracker.JobEntry{}
	jobName := "Add Items Collection - " + *input.Name
	jobType := "ADD_ITEM_TO_PROFILE_COLLECTION"
	jobEntry.JobType = &jobType
	jobEntry.JobName = &jobName
	jobEntry.Input = map[string]string{}
	jobEntry.Input["collectionId"] = strconv.FormatInt(*createdEntry.Id, 10)
	jobEntry.InputAssetInformation = &coredomain.AssetInfo{Id: *input.ItemsAssetId}
	if input.Platform != nil {
		jobEntry.Input["collectionPlatform"] = *input.Platform
	}

	client := jobtracker.New(ctx)
	jobResponse, err := client.CreateJob(jobEntry)
	if err != nil {
		return nil, err
	}
	if jobResponse == nil || jobResponse.Status.Type != "SUCCESS" || len(jobResponse.Data) == 0 || (jobResponse.Data)[0].Id == nil {
		return nil, errors.New("failed to create job for bulk item import")
	}
	jobId := (jobResponse.Data)[0].Id

	updateCollectionEntry := domain.ProfileCollectionEntry{JobId: jobId}
	m.Manager.Update(ctx, *createdEntry.Id, &updateCollectionEntry)
	createdEntry.JobId = jobId
	return nil, nil
}

func (m *Manager) FindById(ctx context.Context, id int64, enrichProfilesSummary bool) (*domain.ProfileCollectionEntry, error) {
	entity, err := m.dao.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	appCtx := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)

	// non featured collections need to be shown to relevant users/partners only
	if entity.Featured == nil || !*entity.Featured {
		if appCtx.PartnerId == nil || appCtx.AccountId == nil {
			// logged out user trying to access non-featured collections
			return nil, errors.New("you are not authorized to access this collection")
		}

		if entity.PartnerId != *appCtx.PartnerId && *appCtx.PartnerId != -1 {
			return nil, errors.New("you are not authorized to delete this collection")
		}
	}

	entry, err := ToEntry(entity)

	if err != nil {
		return nil, err
	}

	if enrichProfilesSummary {
		entry = m.EnrichCollectionWithItemsSummary(ctx, *entry, false)
		// enrich each custom column with computed value
		entry = m.EnrichCollectionWithCustomColumnComputedValues(ctx, *entry, false)
	}
	return entry, nil
}

func (m *Manager) FindByIdInternal(ctx context.Context, id int64) (*domain.ProfileCollectionEntry, error) {
	entity, err := m.dao.FindById(ctx, id)
	if err != nil {
		return nil, err
	}
	entry, err := ToEntry(entity)
	if err != nil {
		return nil, err
	}
	return entry, nil
}

func (m *Manager) FindByShareId(ctx context.Context, shareId string) (*domain.ProfileCollectionEntry, error) {
	entity, err := m.dao.FindByShareId(ctx, shareId)
	if err != nil {
		return nil, err
	}
	entry, _ := ToEntry(entity)
	if err != nil {
		return nil, err
	}
	entry = m.EnrichCollectionWithItemsSummary(ctx, *entry, false)
	// enrich each custom column with computed value
	entry = m.EnrichCollectionWithCustomColumnComputedValues(ctx, *entry, true)
	return entry, nil
}

func (m *Manager) FindBySourceAndSourceId(ctx context.Context, source string, sourceId string) (*domain.ProfileCollectionEntry, error) {
	entity, err := m.dao.FindBySourceAndSourceId(ctx, source, sourceId)
	if err != nil {
		return nil, err
	}
	entry, _ := ToEntry(entity)
	if err != nil {
		return nil, err
	}
	entry = m.EnrichCollectionWithItemsSummary(ctx, *entry, false)
	// enrich each custom column with computed value
	entry = m.EnrichCollectionWithCustomColumnComputedValues(ctx, *entry, false)
	return entry, nil
}

func (m *Manager) FindCollectionByCampaignId(ctx context.Context, campaignId string) (*domain.ProfileCollectionEntry, error) {
	entity, err := m.dao.FindCollectionByCampaignId(ctx, campaignId)
	if err != nil {
		return nil, err
	}
	entry, _ := ToEntry(entity)
	if err != nil {
		return nil, err
	}
	return entry, nil
}

func (m *Manager) FindCollectionByShortlistId(ctx context.Context, shortlistId string) (*domain.ProfileCollectionEntry, error) {
	entity, err := m.dao.FindCollectionByShortlistId(ctx, shortlistId)
	if err != nil {
		return nil, err
	}
	if entity != nil {
		entry, _ := ToEntry(entity)
		return entry, nil
	}
	return nil, nil
}

func (m *Manager) RenewShareId(ctx context.Context, id int64) (*domain.ProfileCollectionEntry, error) {
	shareId, _ := randutil.Alphanumeric(7)
	entry := domain.ProfileCollectionEntry{
		ShareId: &shareId,
	}

	updatedEntry, err := m.Update(ctx, id, &entry)
	if err != nil {
		return nil, err
	}

	return updatedEntry, nil
}

func (m *Manager) Delete(ctx context.Context, id int64) (*domain.ProfileCollectionEntry, error) {
	entry := domain.ProfileCollectionEntry{
		Enabled: helpers.ToBool(false),
	}
	updatedEntry, err := m.Manager.Update(ctx, id, &entry)
	if err != nil {
		return nil, err
	}
	return updatedEntry, nil
}

func (m *Manager) FetchRecentCollectionsForPartner(ctx context.Context) (*[]domain.ProfileCollectionEntry, error) {
	appCtx := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)
	// 2 most recently modified (item added/removed) collections for this partner
	updatedItemEntites := m.itemDAO.FindMostRecentCollectionByItemUpdate(ctx, *appCtx.PartnerId)

	// fetch most recent created collection
	recentlyAddedCollectionEntity, _ := m.dao.FindMostRecentAddedCollection(ctx, *appCtx.PartnerId)
	recentlyAddedCollectionEntry, _ := ToEntry(recentlyAddedCollectionEntity)

	applicableCollections := []domain.ProfileCollectionEntry{}
	if len(updatedItemEntites) > 0 {
		entry, _ := m.Manager.FindById(ctx, updatedItemEntites[0].ProfileCollectionId)
		applicableCollections = append(applicableCollections, *entry)

		if recentlyAddedCollectionEntry != nil && recentlyAddedCollectionEntry.Id != entry.Id {
			applicableCollections = append(applicableCollections, *recentlyAddedCollectionEntry)
		}
		if len(applicableCollections) == 1 && len(updatedItemEntites) > 1 {
			entry2, _ := m.Manager.FindById(ctx, updatedItemEntites[1].ProfileCollectionId)
			applicableCollections = append(applicableCollections, *entry2)
		}
	} else {
		if recentlyAddedCollectionEntry != nil {
			applicableCollections = append(applicableCollections, *recentlyAddedCollectionEntry)
		}
	}

	return &applicableCollections, nil
}

func (m *Manager) Search(ctx context.Context, query coredomain.SearchQuery, sortBy string, sortDir string, page int, size int) ([]domain.ProfileCollectionEntry, int64, error) {

	appCtx := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)
	partnerIdForFilter := "-1"

	val := rest.GetFilterValueForKey(query.Filters, "featured")
	if val != nil && *val == "true" {
		// searching for featured collections, allow irrespective of login or not
		partnerIdForFilter = "-1"
	} else {
		// non featured collections query
		if appCtx.PartnerId == nil || appCtx.AccountId == nil {
			return nil, 0, errors.New("you are not authorized to access collections data")
		}
		partnerIdForFilter = strconv.FormatInt(*appCtx.PartnerId, 10)
	}
	// remove any filter with partnerId coming from anywhere
	query = rest.RemoveKeysFromQuery(query, "partner_id")
	query.Filters = append(query.Filters, coredomain.SearchFilter{FilterType: "EQ", Field: "partner_id", Value: partnerIdForFilter})

	platformAccountCode := rest.GetFilterValueForKey(query.Filters, "platformAccountCode")
	query = rest.RemoveKeysFromQuery(query, "platformAccountCode")
	if platformAccountCode != nil {
		pAccountCode, _ := strconv.ParseInt(*platformAccountCode, 10, 64)
		pId, _ := strconv.ParseInt(partnerIdForFilter, 10, 64)
		collectionIds, err := m.itemDAO.FetchCollectionIdsContainingPlatformAccountCode(ctx, pAccountCode, pId)
		if err != nil {
			return nil, 0, err
		}
		if len(collectionIds) == 0 {
			return nil, 0, nil
		}
		query.Filters = append(query.Filters, coredomain.SearchFilter{FilterType: "IN", Field: "id", Value: strings.Join(collectionIds, ",")})
	}

	entities, filteredCount, err := m.dao.Search(ctx, query, sortBy, sortDir, page, size)

	if err != nil {
		return nil, 0, err
	}

	entries := []domain.ProfileCollectionEntry{}
	for i := 0; i < len(entities); i++ {
		collectionEntry, _ := ToEntry(&entities[i])
		if collectionEntry != nil {
			collectionEntry = m.EnrichCollectionWithItemsSummary(ctx, *collectionEntry, true)
			entries = append(entries, *collectionEntry)
		}
	}
	return entries, filteredCount, err
}

func SearchSocialProfile(ctx context.Context, platform string, profileAccountId *int64, campaignProfileId *int64, profile *coredomain.Profile, source string) (*coredomain.CampaignProfileEntry, error) {
	cpManager := cpmanager.CreateManager(ctx)

	var cp *coredomain.CampaignProfileEntry
	var err error
	if campaignProfileId != nil {
		cp, err = cpManager.FindById(ctx, *campaignProfileId, true)
		if err != nil {
			return nil, err
		}
		if cp == nil {
			return nil, errors.New("No Campaign profile found for id - " + strconv.FormatInt(*campaignProfileId, 10))
		}
	} else if profileAccountId != nil {
		cp, _ = cpManager.FindCampaignProfileByPlatformAccountCode(ctx, platform, *profileAccountId, string(constants.SaasCollection))
		if cp == nil {
			return nil, errors.New("No Campaign profile found for platform account id - " + strconv.FormatInt(*profileAccountId, 10))
		}
	} else if profile != nil && profile.Handle != nil {
		// fetch details by platform and handle
		cp, err = cpManager.FindCampaignProfileByPlatformHandle(ctx, platform, *profile.Handle, source, false) // last parameter is fullRefresh.
		// if fulRefresh is 1, then it will return entire data from beat.
		if err != nil {
			return nil, err
		}
		if cp == nil {
			return nil, errors.New("Invalid Handle - " + *profile.Handle)
		}
	}
	return cp, nil
}

// -------------------------
// Winkl collection migration info
// ------------------------
func (m *Manager) FetchWinklCollectionInfoById(ctx context.Context, collectionId int64) (*domain.WinklCollectionInfoEntry, error) {
	entity, err := m.winklCollectionInfoDAO.FetchWinklCollectionInfoById(ctx, collectionId)

	if entity != nil {
		entry := domain.WinklCollectionInfoEntry{
			WinklCollectionId:        entity.WinklCollectionId,
			WinklCollectionShareId:   entity.WinklCollectionShareId,
			WinklCampaignId:          entity.WinklCampaignId,
			WinklShortlistId:         entity.WinklShortlistId,
			ProfileCollectionId:      entity.ProfileCollectionId,
			ProfileCollectionShareId: entity.ProfileCollectionShareId,
		}
		return &entry, nil
	}

	return nil, err
}

func (m *Manager) FetchWinklCollectionInfoByShareId(ctx context.Context, collectionShareId string) (*domain.WinklCollectionInfoEntry, error) {
	entity, err := m.winklCollectionInfoDAO.FetchWinklCollectionInfoByShareId(ctx, collectionShareId)

	if entity != nil {
		entry := domain.WinklCollectionInfoEntry{
			WinklCollectionId:        entity.WinklCollectionId,
			WinklCollectionShareId:   entity.WinklCollectionShareId,
			WinklCampaignId:          entity.WinklCampaignId,
			WinklShortlistId:         entity.WinklShortlistId,
			ProfileCollectionId:      entity.ProfileCollectionId,
			ProfileCollectionShareId: entity.ProfileCollectionShareId,
		}
		return &entry, nil
	}

	return nil, err
}

// -------------------------------
// item functions
// -------------------------------

func (m *Manager) AddItemsToCollection(ctx context.Context, collection domain.ProfileCollectionEntry, itemsToAdd []domain.ProfileCollectionItemEntry) ([]domain.ProfileCollectionItemEntry, error, []domain.ProfileCollectionItemEntry) {
	appCtx := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)
	var err error

	if len(itemsToAdd) == 0 {
		err := errors.New("no profiles to add")
		return nil, err, nil
	}

	if appCtx.PlanType != nil && (*appCtx.PlanType == constants.PaidPlan || *appCtx.PlanType == constants.SaasPlan) && collection.Source != nil && *collection.Source == string(constants.SaasAtCollection) {
		activityId := ""
		platform := ""
		if len(itemsToAdd) > 0 && itemsToAdd[0].Platform != nil {
			platform = *itemsToAdd[0].Platform
		}
		activityMeta := m.MakeProfileCollectionActivityMeta(platform, *collection.Name)
		_, err := m.partnerUsageManager.LogFreeActivity(ctx, constants.AccountTracking, constants.AddProfileToCollectionActivity, int64(len(itemsToAdd)), *appCtx.PartnerId, *appCtx.AccountId, activityId, activityMeta, "", string(*appCtx.PlanType))
		if err != nil {
			return nil, err, nil
		}
	} else if appCtx.PlanType != nil && *appCtx.PlanType == constants.FreePlan && appCtx.PartnerId != nil && collection.Source != nil && *collection.Source == string(constants.SaasAtCollection) {
		activityId := ""
		platform := ""
		if len(itemsToAdd) > 0 && itemsToAdd[0].Platform != nil {
			platform = *itemsToAdd[0].Platform
		}
		activityMeta := m.MakeProfileCollectionActivityMeta(platform, *collection.Name)
		_, err := m.partnerUsageManager.LogFreeActivity(ctx, constants.AccountTracking, constants.AddProfileToCollectionActivity, int64(len(itemsToAdd)), *appCtx.PartnerId, *appCtx.AccountId, activityId, activityMeta, "", string(*appCtx.PlanType))
		if err != nil {
			return nil, err, nil
		}
	}

	rankInstagram, rankYoutube := m.itemDAO.FindMaxItemRankForCollectionByPlatform(ctx, *collection.Id)
	var itemEntries []domain.ProfileCollectionItemEntry
	var failedItems []domain.ProfileCollectionItemEntry
	for _, itemEntry := range itemsToAdd {
		// check if its a valid social profile in system
		if itemEntry.Platform == nil {
			return nil, errors.New("missing attribute platform"), nil
		}
		cp, err := SearchSocialProfile(ctx, *itemEntry.Platform, itemEntry.PlatformAccountCode, itemEntry.CampaignProfileId, itemEntry.Profile, string(constants.SaasCollection))
		if err != nil {
			failedItems = append(failedItems, itemEntry)
			continue
		}
		if cp == nil {
			failedItems = append(failedItems, itemEntry)
			continue
		}
		itemEntry.CampaignProfileId = &cp.Id
		itemEntry.PlatformAccountCode = &cp.PlatformAccountId
		itemEntry.ProfileSocialId = cp.Handle

		if itemEntry.ShortlistId != nil && *itemEntry.ShortlistId == "-1" {
			itemEntry.ShortlistId = helpers.ToString("")
		}

		// if itemEntry already exist, just make sure to enable it
		existingItemEntity := m.itemDAO.FindCollectionItem(ctx, *collection.Id, *itemEntry.Platform, *itemEntry.PlatformAccountCode)
		if existingItemEntity != nil {
			updateEntity := dao.ProfileCollectionItemEntity{Enabled: helpers.ToBool(true)}
			// Increase ordering value by 1
			if existingItemEntity.Platform == string(constants.InstagramPlatform) {
				rankInstagram++
				updateEntity.Rank = rankInstagram
			} else if existingItemEntity.Platform == string(constants.YoutubePlatform) {
				rankYoutube++
				updateEntity.Rank = rankYoutube
			} else if existingItemEntity.Platform == string(constants.GCCPlatform) {
				updateEntity.Rank = 1
			}
			_, err := m.itemDAO.Update(ctx, existingItemEntity.Id, &updateEntity)
			if err != nil {
				err := errors.New("error while adding profiles")
				return nil, err, nil
			}
			existingItemEntity.Enabled = updateEntity.Enabled
			existingItemEntity.Rank = updateEntity.Rank
		} else {
			if itemEntry.CreatedBy == nil {
				userAccountId := strconv.FormatInt(*appCtx.AccountId, 10)
				itemEntry.CreatedBy = &userAccountId
			}
			itemEntry.PartnerId = collection.PartnerId
			itemEntry.ProfileCollectionId = collection.Id
			itemEntry.Enabled = helpers.ToBool(true)
			if itemEntry.Hidden == nil {
				itemEntry.Hidden = helpers.ToBool(false)
			}
			if itemEntry.CustomColumnsData == nil {
				itemEntry.CustomColumnsData = []domain.CustomColumn{}
			}
			itemEntity, _ := ToItemEntity(&itemEntry)

			// Increase ordering value by 1
			if *itemEntry.Platform == string(constants.InstagramPlatform) {
				rankInstagram++
				itemEntity.Rank = rankInstagram
			} else if *itemEntry.Platform == string(constants.YoutubePlatform) {
				rankYoutube++
				itemEntity.Rank = rankYoutube
			} else if *itemEntry.Platform == string(constants.GCCPlatform) {
				itemEntity.Rank = 1
			}

			existingItemEntity, err = m.itemDAO.Create(ctx, itemEntity)
			if err != nil {
				return nil, err, nil
			} else if existingItemEntity == nil {
				return nil, errors.New("failed to add item"), nil
			}
		}
		if itemEntry.CustomColumnsData != nil {
			err = m.itemDAO.UpdateItemCustomColumnsData(ctx, existingItemEntity.Id, itemEntry.CustomColumnsData)
			if err != nil {
				return nil, err, nil
			}
		}
		entry, _ := ToItemEntry(existingItemEntity)
		itemEntries = append(itemEntries, *entry)
	}
	return itemEntries, err, failedItems
}

func (m *Manager) DeleteItemsFromCollection(ctx context.Context, collectionId int64, items []domain.ProfileCollectionItemEntry) ([]domain.ProfileCollectionItemEntry, error) {
	if len(items) == 0 {
		err := errors.New("no items to delete")
		return nil, err
	}

	updatedEntities := []domain.ProfileCollectionItemEntry{}
	for _, item := range items {
		itemEntity := dao.ProfileCollectionItemEntity{
			Enabled: helpers.ToBool(false),
		}
		var err error
		var updatedEntity *dao.ProfileCollectionItemEntity
		if item.Id != nil {
			updatedEntity, err = m.itemDAO.Update(ctx, *item.Id, &itemEntity)
		} else if item.CampaignProfileId != nil {
			updatedEntity, err = m.itemDAO.UpdateByCollectionItemCampaignProfileId(ctx, collectionId, *item.CampaignProfileId, itemEntity)
		} else if item.PlatformAccountCode != nil {
			updatedEntity, err = m.itemDAO.UpdateCollectionItemByPlatformAccountCode(ctx, collectionId, *item.PlatformAccountCode, itemEntity)
		} else if item.Platform != nil && item.Profile != nil && item.Profile.Handle != nil {
			profile, searchErr := m.socialSearchManager.FindProfileByPlatformHandle(ctx, *item.Platform, *item.Profile.Handle, string(constants.SaasCollection), false)
			if searchErr == nil && profile != nil {
				platformAccountCode, _ := strconv.ParseInt(profile.PlatformCode, 10, 64)
				updatedEntity, err = m.itemDAO.UpdateCollectionItemByPlatformAccountCode(ctx, collectionId, platformAccountCode, itemEntity)
				if err != nil {
					log.Error(err)
					err := errors.New("error while deleting profiles")
					return nil, err
				}
			} else {
				log.Error(searchErr)
				err := errors.New("error while finding profile")
				return nil, err
			}
		}
		if err != nil {
			log.Error(err)
			err := errors.New("error while deleting profiles")
			return nil, err
		}
		if updatedEntity == nil {
			log.Error(err)
			err := errors.New("error while deleting profiles")
			return nil, err
		}
		updatedEntry, _ := ToItemEntry(updatedEntity)
		updatedEntities = append(updatedEntities, *updatedEntry)
	}
	return updatedEntities, nil
}

func (m *Manager) UpdateCollectionItemsInBulk(ctx context.Context, collectionId int64, items []domain.ProfileCollectionItemEntry) ([]domain.ProfileCollectionItemEntry, error) {
	if len(items) == 0 {
		err := errors.New("no items to update")
		return nil, err
	}

	updatedEntities := []domain.ProfileCollectionItemEntry{}
	for _, item := range items {
		if item.ShortlistId != nil && *item.ShortlistId == "-1" {
			item.ShortlistId = helpers.ToString("")
		}
		inputEntity, _ := ToItemEntity(&item)
		var err error
		var updatedEntity *dao.ProfileCollectionItemEntity
		var errorMsg string
		if item.Id != nil {
			updatedEntity, err = m.itemDAO.Update(ctx, *item.Id, inputEntity)
			errorMsg = "invalid id"
		} else if item.CampaignProfileId != nil {
			updatedEntity, err = m.itemDAO.UpdateByCollectionItemCampaignProfileId(ctx, collectionId, *item.CampaignProfileId, *inputEntity)
			errorMsg = "invalid campaignProfileId"
		} else if item.PlatformAccountCode != nil {
			updatedEntity, err = m.itemDAO.UpdateCollectionItemByPlatformAccountCode(ctx, collectionId, *item.PlatformAccountCode, *inputEntity)
			errorMsg = "invalid platformAccountCode"
		}
		if err != nil {
			err := errors.New(errorMsg)
			return nil, err
		}
		if updatedEntity == nil {
			err := errors.New("error while Updating profiles")
			return nil, err
		}
		if item.CustomColumnsData != nil {
			err = m.itemDAO.UpdateItemCustomColumnsData(ctx, updatedEntity.Id, item.CustomColumnsData)
			if err != nil {
				return nil, err
			}
		}
		updatedEntry, _ := ToItemEntry(updatedEntity)
		updatedEntities = append(updatedEntities, *updatedEntry)
	}
	return updatedEntities, nil
}

func (m *Manager) SearchItemsInCollection(ctx context.Context, query coredomain.SearchQuery, sortBy string, sortDir string, page int, size int) ([]domain.ProfileCollectionItemEntry, int64, error) {

	appCtx := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)

	if len(query.Filters) == 0 {
		return nil, 0, errors.New("empty filters")
	}

	collectionId := rest.GetFilterValueForKey(query.Filters, "collectionId")
	collectionShareId := rest.GetFilterValueForKey(query.Filters, "collectionShareId")
	campaignId := rest.GetFilterValueForKey(query.Filters, "campaignId")
	shortlistId := rest.GetFilterValueForKey(query.Filters, "shortlistId")
	if collectionId == nil && collectionShareId == nil && campaignId == nil && shortlistId == nil {
		return nil, 0, errors.New("empty collection/campaign/shortlist id")
	}

	platform := rest.GetFilterValueForKey(query.Filters, "platform")
	query = rest.RemoveKeysFromQuery(query, "platform")
	if platform == nil && (campaignId != nil || shortlistId != nil) {
		var c *domain.ProfileCollectionEntry
		if campaignId != nil {
			c, _ = m.FindCollectionByCampaignId(ctx, *campaignId)
		} else if shortlistId != nil {
			c, _ = m.FindCollectionByShortlistId(ctx, *shortlistId)
		}
		if c == nil {
			return []domain.ProfileCollectionItemEntry{}, 0, nil
		} else {
			platform = c.CampaignPlatform
			cIdStr := strconv.FormatInt(*c.Id, 10)
			collectionId = &cIdStr
		}
	}
	if platform == nil {
		return nil, 0, errors.New("empty platform")
	}

	query.Filters = append(query.Filters, coredomain.SearchFilter{FilterType: "EQ", Field: "itemEnabled", Value: "true"})
	if collectionId == nil && collectionShareId != nil {
		c, _ := m.FindByShareId(ctx, *collectionShareId)
		cIdStr := strconv.FormatInt(*c.Id, 10)
		collectionId = &cIdStr
		query.Filters = append(query.Filters, coredomain.SearchFilter{FilterType: "EQ", Field: "itemHidden", Value: "false"})
	}

	var profiles *[]coredomain.Profile
	filteredCount := int64(0)
	if *platform == string(constants.InstagramPlatform) || *platform == string(constants.YoutubePlatform) {
		socialDiscoveryManager := discoverymanager.CreateSearchManager(ctx)
		profiles, filteredCount, _ = socialDiscoveryManager.SearchByPlatform(ctx, *platform, query, sortBy, sortDir, page, size, string(constants.SaasCollection), true, true)
	} else if *platform == string(constants.GCCPlatform) {
		query.Filters = append(query.Filters, coredomain.SearchFilter{FilterType: "EQ", Field: "platform", Value: string(constants.GCCPlatform)})
		cpManager := cpmanager.CreateManager(ctx)
		var campaignProfiles []coredomain.CampaignProfileEntry
		var err error
		campaignProfiles, filteredCount, err = cpManager.Search(ctx, query, sortBy, sortDir, page, size, false)
		if err == nil && campaignProfiles != nil {
			var userProfiles []coredomain.Profile
			for i := range campaignProfiles {
				cp := campaignProfiles[i]
				userProfile := coredomain.Profile{
					PlatformCode:    strconv.FormatInt(cp.PlatformAccountId, 10),
					Handle:          cp.Handle,
					Platform:        cp.Platform,
					Name:            cp.UserDetails.Name,
					Phone:           cp.UserDetails.Phone,
					Email:           cp.UserDetails.Email,
					CampaignProfile: &cp,
				}
				userProfiles = append(userProfiles, userProfile)
			}
			profiles = &userProfiles
		}
	}

	var platformAccountCodes []int64
	if profiles != nil && len(*profiles) > 0 {
		for i := range *profiles {
			socialProfile := (*profiles)[i]
			platformAccountId, _ := strconv.ParseInt(socialProfile.PlatformCode, 10, 64)
			platformAccountCodes = append(platformAccountCodes, platformAccountId)
		}
	}

	var items []domain.ProfileCollectionItemEntry
	if len(platformAccountCodes) > 0 {
		cId, _ := strconv.ParseInt(*collectionId, 10, 64)
		profileCollectionItemEntities, err := m.itemDAO.FetchCollectionItemsByPlatformAccountCode(ctx, cId, platformAccountCodes, *platform)
		if err != nil {
			return nil, 0, err
		}
		if profileCollectionItemEntities == nil || len(profileCollectionItemEntities) != len(platformAccountCodes) {
			return nil, 0, errors.New("some profiles are not part of collection items. Data Issue")
		}
		var platformAccountCodeToItemMap = map[int64]domain.ProfileCollectionItemEntry{}
		for i := range profileCollectionItemEntities {
			entity := (profileCollectionItemEntities)[i]
			itemEntry, _ := ToItemEntry(&entity)
			platformAccountCodeToItemMap[*itemEntry.PlatformAccountCode] = *itemEntry
		}

		for i := range *profiles {
			socialProfile := (*profiles)[i]
			platformAccountId, _ := strconv.ParseInt(socialProfile.PlatformCode, 10, 64)
			if itemEntry, ok := platformAccountCodeToItemMap[platformAccountId]; ok {
				itemEntry.Profile = &socialProfile
				itemEntry.CustomColumnsData = m.FetchItemCustomColumnsData(ctx, *itemEntry.Id)
				if appCtx.PartnerId == nil || *appCtx.PartnerId != -1 {
					itemEntry.Commercials = nil
				}
				items = append(items, itemEntry)
			}
		}
	}

	return items, filteredCount, nil
}

func (m *Manager) CountCollectionsForPartner(ctx context.Context, partnerId int64) (int64, error) {
	count, err := m.dao.CountCollectionsForPartner(ctx, partnerId)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (m *Manager) FetchItemCustomColumnsData(ctx context.Context, profileCollectionItemId int64) []domain.CustomColumn {
	rows, err := m.itemDAO.FetchExistingCustomColumnRowsForItem(ctx, profileCollectionItemId)
	var customColumnsArray []domain.CustomColumn
	if err == nil && rows != nil {
		for i := range rows {
			row := rows[i]
			customColumnsArray = append(customColumnsArray, domain.CustomColumn{
				Id:    &row.Id,
				Key:   row.Key,
				Value: row.Value,
			})
		}
	}
	return customColumnsArray
}

func (m *Manager) MakeProfileCollectionActivityMeta(platform string, name string) string {
	var requestMetaStr string
	var requestMeta = make(map[string]interface{})
	requestMeta["platform"] = platform
	requestMeta["name"] = name

	requestMetaJson, err := json.Marshal(requestMeta)
	if err != nil {
		log.Error(err)
		return requestMetaStr
	}
	requestMetaStr = string(requestMetaJson)
	return requestMetaStr
}

package manager

import (
	partnerusagemanager "coffee/app/partnerusage/manager"
	"coffee/app/postcollection/dao"
	"coffee/app/postcollection/domain"
	"coffee/constants"
	"coffee/core/appcontext"
	"coffee/core/rest"
	"coffee/helpers"
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type PostCollectionItemManager struct {
	*rest.Manager[domain.PostCollectionItemEntry, dao.PostCollectionItemEntity, int64]
	dao                 *dao.PostCollectionItemDAO
	partnerUsageManager *partnerusagemanager.PartnerUsageManager
}

func (m *PostCollectionItemManager) init(ctx context.Context) {
	newDAO := dao.CreatePostCollectionItemDAO(ctx)
	m.Manager = rest.NewManager(ctx, newDAO.Dao, m.ToEntry, m.ToEntity)
	m.dao = newDAO
	m.partnerUsageManager = partnerusagemanager.CreatePartnerUsageManager(ctx)
}

func CreatePostCollectionItemManager(ctx context.Context) *PostCollectionItemManager {
	manager := &PostCollectionItemManager{}
	manager.init(ctx)
	return manager
}

func (m *PostCollectionItemManager) ToEntity(entry *domain.PostCollectionItemEntry) (*dao.PostCollectionItemEntity, error) {
	entity := &dao.PostCollectionItemEntity{}
	if entry.PostCollectionId != nil {
		entity.PostCollectionId = *entry.PostCollectionId
	}
	if entry.Platform != nil {
		entity.Platform = *entry.Platform
	}
	if entry.ShortCode != nil {
		entity.ShortCode = *entry.ShortCode
	}
	if entry.PostType != nil {
		entity.PostType = *entry.PostType
	}
	if entry.Enabled != nil {
		entity.Enabled = entry.Enabled
	}
	if entry.Bookmarked != nil {
		entity.Bookmarked = entry.Bookmarked
	}
	if entry.Cost != nil {
		entity.Cost = entry.Cost
	}
	if entry.PostedByHandle != nil {
		entity.PostedByHandle = entry.PostedByHandle
	}
	if entry.CreatedBy != nil {
		entity.CreatedBy = *entry.CreatedBy
	}
	if entry.SponsorLinks != nil {
		entity.SponsorLinks.Set(entry.SponsorLinks)
	}
	if entry.PostTitle != nil {
		entity.PostTitle = entry.PostTitle
	}
	if entry.PostThumbnail != nil {
		entity.PostThumbnail = entry.PostThumbnail
	}
	if entry.PostLink != nil {
		entity.PostLink = entry.PostLink
	}
	if entry.MetricsIngestionFreq != nil {
		entity.MetricsIngestionFreq = *entry.MetricsIngestionFreq
	}
	if entry.RetrieveData != nil {
		entity.RetrieveData = entry.RetrieveData
	}
	if entry.ShowInReport != nil {
		entity.ShowInReport = entry.ShowInReport
	}
	if entry.PostedByCpId != nil {
		entity.PostedByCpId = entry.PostedByCpId
	}
	return entity, nil
}

// Entry Methods

func (m *PostCollectionItemManager) ToEntry(entity *dao.PostCollectionItemEntity) (*domain.PostCollectionItemEntry, error) {
	if entity != nil {
		entry := &domain.PostCollectionItemEntry{
			Id:                   &entity.Id,
			Platform:             &entity.Platform,
			ShortCode:            &entity.ShortCode,
			PostCollectionId:     &entity.PostCollectionId,
			Cost:                 entity.Cost,
			Bookmarked:           entity.Bookmarked,
			Enabled:              entity.Enabled,
			PostType:             &entity.PostType,
			PostedByHandle:       entity.PostedByHandle,
			PostThumbnail:        entity.PostThumbnail,
			PostTitle:            entity.PostTitle,
			PostLink:             entity.PostLink,
			CreatedBy:            &entity.CreatedBy,
			MetricsIngestionFreq: &entity.MetricsIngestionFreq,
			RetrieveData:         entity.RetrieveData,
			ShowInReport:         entity.ShowInReport,
			PostedByCpId:         entity.PostedByCpId,

			CreatedAt: helpers.ToInt64(entity.CreatedAt.Unix()),
			UpdatedAt: helpers.ToInt64(entity.UpdatedAt.Unix()),
		}
		entity.SponsorLinks.AssignTo(&entry.SponsorLinks)
		return entry, nil
	}
	return nil, nil
}

func (m *PostCollectionItemManager) AddItemsToCollection(ctx context.Context, collection domain.PostCollectionEntry, items []domain.PostCollectionItemEntry) (*[]domain.PostCollectionItemEntry, error) {
	// Check For Consumption For The CSV FLOW

	if len(items) == 0 {
		err := errors.New("no items to add")
		return nil, err
	}

	appCtx := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)
	if appCtx.PlanType != nil && (*appCtx.PlanType == constants.PaidPlan || *appCtx.PlanType == constants.SaasPlan) && collection.Source != nil && *collection.Source == string(constants.SaasCollection) {
		activityId := ""
		platform := ""
		if len(items) > 0 && items[0].Platform != nil {
			platform = *items[0].Platform
		}
		campaignName := collection.Name
		activityMeta := m.MakePostCollectionActivityMeta(*collection.Id, platform, *campaignName)
		_, err := m.partnerUsageManager.LogFreeActivity(ctx, constants.CampaignReport, constants.AddPostToCollectionActivity, int64(len(items)), *appCtx.PartnerId, *appCtx.AccountId, activityId, activityMeta, "", string(*appCtx.PlanType))
		if err != nil {
			return nil, err
		}
	} else if appCtx.PlanType != nil && *appCtx.PlanType == constants.FreePlan && appCtx.PartnerId != nil && collection.Source != nil && *collection.Source == string(constants.SaasCollection) {
		activityId := ""
		platform := ""
		if len(items) > 0 && items[0].Platform != nil {
			platform = *items[0].Platform
		}
		campaignName := collection.Name
		activityMeta := m.MakePostCollectionActivityMeta(*collection.Id, platform, *campaignName)
		_, err := m.partnerUsageManager.LogFreeActivity(ctx, constants.CampaignReport, constants.AddPostToCollectionActivity, int64(len(items)), *appCtx.PartnerId, *appCtx.AccountId, activityId, activityMeta, "", string(*appCtx.PlanType))
		if err != nil {
			return nil, err
		}
	}

	createdEntries := []domain.PostCollectionItemEntry{}
	for _, itemEntry := range items {
		// if itemEntry already exist, just make sure to enable it
		if itemEntry.PostType != nil && *itemEntry.PostType == "story" && itemEntry.ShortCode == nil && itemEntry.PostedByHandle != nil { // story items need to be filled with dummy shortcode if not present
			dummyShortCode := helpers.GenerateStoryDummyShortCode(*itemEntry.PostedByHandle)
			itemEntry.ShortCode = &dummyShortCode
		}
		if itemEntry.Platform == nil || itemEntry.ShortCode == nil || itemEntry.PostType == nil {
			return nil, errors.New("Missing attributes platform/shortCode/platformAccountHandle/postType")
		}
		platformUpper := strings.ToUpper(*itemEntry.Platform)
		itemEntry.Platform = &platformUpper
		existingItemEntity, err := m.dao.FindCollectionItemByPlatformAndShortCode(ctx, *collection.Id, *itemEntry.Platform, *itemEntry.ShortCode)
		if err != nil {
			return nil, err
		}
		if existingItemEntity != nil && existingItemEntity.Id > 0 {
			updateEntity := dao.PostCollectionItemEntity{Enabled: helpers.ToBool(true)}
			if itemEntry.ShowInReport != nil {
				updateEntity.ShowInReport = itemEntry.ShowInReport
			}
			if itemEntry.RetrieveData != nil {
				updateEntity.RetrieveData = itemEntry.RetrieveData
			}
			updatedEntity, err := m.dao.Update(ctx, existingItemEntity.Id, &updateEntity)
			if err != nil {
				return nil, err
			} else if updatedEntity == nil {
				return nil, errors.New("error while adding items")
			}
			createdEntry, _ := m.ToEntry(updatedEntity)
			createdEntries = append(createdEntries, *createdEntry)
		} else {
			appCtx := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)
			userAccountId := strconv.FormatInt(*appCtx.AccountId, 10)
			itemEntry.CreatedBy = &userAccountId

			itemEntry.PostCollectionId = collection.Id
			itemEntry.MetricsIngestionFreq = collection.MetricsIngestionFreq
			itemEntry.Enabled = helpers.ToBool(true)
			itemEntry.Bookmarked = helpers.ToBool(false)
			if itemEntry.SponsorLinks == nil {
				itemEntry.SponsorLinks = []string{}
			}

			itemEntity, err := m.ToEntity(&itemEntry)

			createdEntity, err := m.dao.Create(ctx, itemEntity)
			if err != nil {
				return nil, err
			} else if createdEntity == nil {
				return nil, errors.New("Failed to Add Item")
			}
			createdEntry, _ := m.ToEntry(createdEntity)
			createdEntries = append(createdEntries, *createdEntry)
		}

	}
	return &createdEntries, nil
}

func (m *PostCollectionItemManager) RemoveItemsFromCollection(ctx context.Context, items []domain.PostCollectionItemEntry) (*[]domain.PostCollectionItemEntry, error) {
	if len(items) == 0 {
		err := errors.New("no items to remove")
		return nil, err
	}

	itemIds := []int64{}
	for _, itemEntry := range items {
		itemIds = append(itemIds, *itemEntry.Id)
	}
	existingItemEntities, err := m.dao.FindByIds(ctx, itemIds)
	if err != nil {
		return nil, err
	}
	if len(existingItemEntities) == 0 || len(existingItemEntities) != len(itemIds) {
		return nil, errors.New("Invalid Item Ids")
	}

	itemId2EntityMap := map[int64]dao.PostCollectionItemEntity{}
	for _, entity := range existingItemEntities {
		itemId2EntityMap[entity.Id] = entity
	}

	removedItemEntries := []domain.PostCollectionItemEntry{}
	for _, itemEntry := range items {
		existingItemEntity, ok := itemId2EntityMap[*itemEntry.Id]
		if !ok {
			return nil, errors.New("Invalid item id " + strconv.FormatInt(*itemEntry.Id, 10))
		}
		updateEntity := dao.PostCollectionItemEntity{Enabled: helpers.ToBool(false)}
		updatedEntity, err := m.dao.Update(ctx, existingItemEntity.Id, &updateEntity)
		if err != nil {
			return nil, err
		} else if updatedEntity == nil {
			err := errors.New("error while removing item")
			return nil, err
		}
		updatedEntry, _ := m.ToEntry(updatedEntity)
		removedItemEntries = append(removedItemEntries, *updatedEntry)
	}
	return &removedItemEntries, nil
}

func (m *PostCollectionItemManager) UpdateItemByPlatformAndShortCode(ctx context.Context, collectionId string, platform string, shortCode string, input domain.PostCollectionItemEntry) (*domain.PostCollectionItemEntry, error) {
	item, err := m.dao.FindCollectionItemByPlatformAndShortCode(ctx, collectionId, platform, shortCode)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, errors.New("Invalid Item Id")
	}
	updatedEntry, err := m.Update(ctx, item.Id, &input)
	if err != nil {
		return nil, err
	}
	return updatedEntry, nil
}

func (m *PostCollectionItemManager) RemoveItemsFromCollectionByShortCode(ctx context.Context, collectionId string, items []domain.PostCollectionItemEntry) (*[]domain.PostCollectionItemEntry, error) {
	if len(items) == 0 {
		return nil, errors.New("no items to remove")
	}

	itemId2EntityMap := map[string]dao.PostCollectionItemEntity{}
	for _, itemEntry := range items {
		if itemEntry.Platform == nil || itemEntry.ShortCode == nil {
			return nil, errors.New("Bad Request: Invalid item details")
		}

		entity, err := m.dao.FindCollectionItemByPlatformAndShortCode(ctx, collectionId, *itemEntry.Platform, *itemEntry.ShortCode)
		if err != nil {
			return nil, err
		}
		if entity == nil {
			return nil, errors.New("No item found in collection for code " + *itemEntry.ShortCode + " on platform " + *itemEntry.Platform)
		}
		itemId2EntityMap[*itemEntry.ShortCode] = *entity
	}

	removedItemEntries := []domain.PostCollectionItemEntry{}
	for _, itemEntry := range items {
		existingItemEntity, _ := itemId2EntityMap[*itemEntry.ShortCode]
		updateEntity := dao.PostCollectionItemEntity{Enabled: helpers.ToBool(false)}
		updatedEntity, err := m.dao.Update(ctx, existingItemEntity.Id, &updateEntity)
		if err != nil {
			return nil, err
		} else if updatedEntity == nil {
			err := errors.New("error while removing item")
			return nil, err
		}
		updatedEntry, _ := m.ToEntry(updatedEntity)
		removedItemEntries = append(removedItemEntries, *updatedEntry)
	}
	return &removedItemEntries, nil
}

func (m *PostCollectionItemManager) CountCollectionItemsInACollection(ctx context.Context, collectionId string) (int64, error) {
	count, err := m.dao.CountCollectionItemsInACollection(ctx, collectionId)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (m *PostCollectionItemManager) GetLatestItemAddedOn(ctx context.Context, collectionId string) (*time.Time, error) {
	timestamp, err := m.dao.GetLatestItemAddedOn(ctx, collectionId)
	if err != nil {
		return timestamp, err
	}
	return timestamp, nil
}

func (m *PostCollectionItemManager) MakePostCollectionActivityMeta(id string, platform string, name string) string {
	var requestMetaStr string
	var requestMeta domain.CampaignReportActivityMeta
	requestMeta.Id = id
	requestMeta.Platform = platform
	requestMeta.Result.Name = name

	requestMetaJson, err := json.Marshal(requestMeta)
	if err != nil {
		log.Error(err)
		return requestMetaStr
	}
	requestMetaStr = string(requestMetaJson)
	return requestMetaStr
}

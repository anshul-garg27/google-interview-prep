package api

import (
	camanager "coffee/app/collectionanalytics/manager"
	PartnerusageManager "coffee/app/partnerusage/manager"
	"coffee/app/profilecollection/dao"
	domain "coffee/app/profilecollection/domain"
	"coffee/app/profilecollection/manager"
	"coffee/constants"
	"coffee/core/appcontext"
	coredomain "coffee/core/domain"
	"coffee/core/rest"
	"coffee/helpers"
	"context"
	"errors"
	"strconv"
	"strings"
)

// common helper methods

func getPrimaryDataStore() coredomain.DataStore {
	return coredomain.Postgres
}

// -------------------------------
// profile collection service
// -------------------------------

type Service struct {
	*rest.Service[domain.ProfileCollectionResponse, domain.ProfileCollectionEntry, dao.ProfileCollectionEntity, int64]
	Manager                    *manager.Manager
	CollectionAnalyticsManager *camanager.CollectionAnalyticsManager
	PartnerusageManager        *PartnerusageManager.PartnerUsageManager
}

func (s *Service) Init(ctx context.Context) {
	manager := manager.CreateManager(ctx)
	service := rest.NewService(ctx, manager.Manager, getPrimaryDataStore, domain.CreateResponse, domain.CreateErrorResponse)
	s.Service = service
	s.Manager = manager
	s.CollectionAnalyticsManager = camanager.CreateCollectionAnalyticsManager(ctx)
	s.PartnerusageManager = PartnerusageManager.CreatePartnerUsageManager(ctx)
}

func CreateService(ctx context.Context) *Service {
	service := &Service{}
	service.Init(ctx)
	return service
}

func NewProfileCollectionService() *Service {
	service := CreateService(context.Background())
	return service
}

func (s *Service) Create(ctx context.Context, entry *domain.ProfileCollectionEntry) domain.ProfileCollectionResponse {
	appCtx := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)
	if (appCtx.PlanType == nil || *appCtx.PlanType == constants.FreePlan) && entry.Source != nil && *entry.Source == string(constants.SaasAtCollection) {
		// check no of collections partner has created till now
		count, err := s.Manager.CountCollectionsForPartner(ctx, *appCtx.PartnerId)
		if err != nil {
			return domain.CreateErrorResponse(err, 0)
		}

		if count >= 5 {
			return domain.CreateErrorResponse(errors.New("you are not allowed to create more than 5 profile trackers. Please upgrade your account"), 0)
		}

		// check limit for no of items that can be added
		if entry.JobId == nil && entry.Items != nil && len(entry.Items) > 5 {
			return domain.CreateErrorResponse(errors.New("you are not allowed to add more than 5 profiles in one tracker. Please upgrade your account"), 0)
		}
	}

	if entry.CustomColumns != nil {
		for i := range entry.CustomColumns {
			for j := range entry.CustomColumns[i].Columns {
				if entry.CustomColumns[i].Columns[j].Key == nil || *entry.CustomColumns[i].Columns[j].Key == "" {
					key := strings.ReplaceAll(strings.ToLower(entry.CustomColumns[i].Columns[j].Title), " ", "_")
					entry.CustomColumns[i].Columns[j].Key = &key
				}
			}
		}
	}
	if appCtx.PlanType != nil && (*appCtx.PlanType == constants.PaidPlan || *appCtx.PlanType == constants.SaasPlan) && appCtx.PartnerId != nil && entry.Source != nil && *entry.Source == string(constants.SaasAtCollection) {
		activityId := ""
		platform := ""
		if len(entry.Items) > 0 && entry.Items[0].Platform != nil {
			platform = *entry.Items[0].Platform
		}
		activityMeta := s.Manager.MakeProfileCollectionActivityMeta(platform, *entry.Name)
		_, err := s.PartnerusageManager.LogPaidActivity(ctx, constants.AccountTracking, constants.AccountTrackingActivity, 0, *appCtx.PartnerId, *appCtx.AccountId, activityId, activityMeta, "", string(*appCtx.PlanType))
		if err != nil {
			return domain.CreateErrorResponse(err, 0)
		}
	} else if appCtx.PlanType != nil && *appCtx.PlanType == constants.FreePlan && appCtx.PartnerId != nil && entry.Source != nil && *entry.Source == string(constants.SaasAtCollection) {
		activityId := ""
		platform := ""
		if len(entry.Items) > 0 && entry.Items[0].Platform != nil {
			platform = *entry.Items[0].Platform
		}
		activityMeta := s.Manager.MakeProfileCollectionActivityMeta(platform, *entry.Name)
		_, err := s.PartnerusageManager.LogFreeActivity(ctx, constants.AccountTracking, constants.AccountTrackingActivity, 0, *appCtx.PartnerId, *appCtx.AccountId, activityId, activityMeta, "", string(*appCtx.PlanType))
		if err != nil {
			return domain.CreateErrorResponse(err, 0)
		}
	}
	e, err := s.Manager.Create(ctx, entry)
	if err != nil {
		return domain.CreateErrorResponse(err, 0)
	}
	return domain.CreateResponse([]domain.ProfileCollectionEntry{*e}, 1, "", "Record Created Successfully")
}

func (s *Service) Update(ctx context.Context, id int64, entry *domain.ProfileCollectionEntry) domain.ProfileCollectionResponse {
	_, err := s.Manager.FindById(ctx, id, false)
	if err != nil {
		return domain.CreateErrorResponse(err, 0)
	}

	if entry.CustomColumns != nil {
		for i := range entry.CustomColumns {
			for j := range entry.CustomColumns[i].Columns {
				if entry.CustomColumns[i].Columns[j].Key == nil || *entry.CustomColumns[i].Columns[j].Key == "" {
					key := strings.ReplaceAll(strings.ToLower(entry.CustomColumns[i].Columns[j].Title), " ", "_")
					entry.CustomColumns[i].Columns[j].Key = &key
				}
			}
		}
	}

	_, err = s.Manager.Update(ctx, id, entry)
	if err != nil {
		return domain.CreateErrorResponse(err, 0)
	}
	return domain.CreateResponse(nil, 1, "", "Record Updated Successfully")
}

func (s *Service) FindById(ctx context.Context, id int64) domain.ProfileCollectionResponse {
	entry, err := s.Manager.FindById(ctx, id, true)
	if err != nil {
		return domain.CreateErrorResponse(err, 0)
	}
	return domain.CreateResponse([]domain.ProfileCollectionEntry{*entry}, 1, "", "Record Retrieved Successfully")
}

func (s *Service) FindByShareId(ctx context.Context, shareId string) domain.ProfileCollectionResponse {
	collection, err := s.Manager.FindByShareId(ctx, shareId)
	if collection != nil && !*collection.ShareEnabled {
		return domain.CreateErrorResponse(errors.New("this collection can not be shared with anyone"), 0)
	}
	count := int64(1)
	collections := []domain.ProfileCollectionEntry{}
	if collection != nil {
		collections = append(collections, *collection)
	}
	if err != nil {
		return domain.CreateErrorResponse(err, 0)
	}
	return domain.CreateResponse(collections, count, "", "Record Retrieved Successfully")
}

func (s *Service) RenewShareId(ctx context.Context, id int64) domain.ProfileCollectionResponse {

	_, err := s.Manager.FindById(ctx, id, false)
	if err != nil {
		return domain.CreateErrorResponse(err, 0)
	}

	updatedEntry, err := s.Manager.RenewShareId(ctx, id)
	if err != nil {
		return domain.CreateErrorResponse(err, 0)
	}
	return domain.CreateResponse([]domain.ProfileCollectionEntry{*updatedEntry}, 1, "", "Record Updated Successfully")
}

func (s *Service) Delete(ctx context.Context, id int64) domain.ProfileCollectionResponse {

	_, err := s.Manager.FindById(ctx, id, false)
	if err != nil {
		return domain.CreateErrorResponse(err, 0)
	}

	_, err = s.Manager.Delete(ctx, id)
	if err != nil {
		return domain.CreateErrorResponse(err, 0)
	}
	return domain.CreateResponse(nil, 1, "", "Record Deleted Successfully")
}

func (s *Service) FetchRecentCollectionsForPartner(ctx context.Context) domain.ProfileCollectionResponse {
	collections, err := s.Manager.FetchRecentCollectionsForPartner(ctx)
	if err != nil {
		return domain.CreateErrorResponse(err, 0)
	}
	return domain.CreateResponse(*collections, 1, "", "Record Retrieved Successfully")
}

func (s *Service) Search(ctx context.Context, searchQuery coredomain.SearchQuery, sortBy string, sortDir string, page int, size int) domain.ProfileCollectionResponse {
	entries, filteredCount, err := s.Manager.Search(ctx, searchQuery, sortBy, sortDir, page, size)
	if err != nil {
		return domain.CreateErrorResponse(err, 0)
	}
	var nextCursor string
	if len(entries) >= size {
		nextCursor = strconv.Itoa(page + 1)
	} else {
		nextCursor = ""
	}
	entries, _ = s.EnrichMetricsAndItems(ctx, entries)
	return domain.CreateResponse(entries, filteredCount, nextCursor, "Record(s) Retrieved Successfully")
}

func (s *Service) AddItemsToCollection(ctx context.Context, id int64, input []domain.ProfileCollectionItemEntry) domain.ProfileCollectionItemResponse {
	appCtx := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)
	collection, err := s.Manager.FindById(ctx, id, false)
	if err != nil {
		return domain.CreateItemAddResponse(nil, err, nil)
	}
	if collection == nil || !*collection.Enabled || *collection.PartnerId != *appCtx.PartnerId {
		return domain.CreateItemAddResponse(nil, errors.New("invalid Collection"), nil)
	}

	items, err, failedItems := s.Manager.AddItemsToCollection(ctx, *collection, input)
	return domain.CreateItemAddResponse(&items, err, failedItems)
}

func (s *Service) DeleteItemsFromCollection(ctx context.Context, id int64, input []domain.ProfileCollectionItemEntry) domain.ProfileCollectionItemResponse {
	_, err := s.Manager.FindById(ctx, id, false)
	if err != nil {
		return domain.CreateItemDeleteResponse(err)
	}
	_, err = s.Manager.DeleteItemsFromCollection(ctx, id, input)
	return domain.CreateItemDeleteResponse(err)
}

func (s *Service) UpdateCollectionItemsInBulk(ctx context.Context, id int64, input []domain.ProfileCollectionItemEntry) domain.ProfileCollectionItemResponse {
	_, err := s.Manager.FindById(ctx, id, false)
	if err != nil {
		return domain.CreateItemUpdateResponse(err)
	}
	_, err = s.Manager.UpdateCollectionItemsInBulk(ctx, id, input)
	return domain.CreateItemUpdateResponse(err)
}

func (s *Service) UpdateCollectionItemsInBulkByShareId(ctx context.Context, shareId string, items []domain.ProfileCollectionItemEntry) domain.ProfileCollectionItemResponse {
	c, err := s.Manager.FindByShareId(ctx, shareId)
	if err != nil {
		return domain.CreateItemUpdateResponse(err)
	}
	if c == nil {
		return domain.CreateItemUpdateResponse(errors.New("invalid Collection Share Id"))
	}

	// in update by shareId, only certain fields are allowed to be updated. Filter out rest of the fields
	var input []domain.ProfileCollectionItemEntry
	for _, item := range items {
		input = append(input, domain.ProfileCollectionItemEntry{
			Id:                    item.Id,
			CampaignProfileId:     item.CampaignProfileId,
			PlatformAccountCode:   item.PlatformAccountCode,
			CustomColumnsData:     item.CustomColumnsData,
			BrandSelectionStatus:  item.BrandSelectionStatus,
			BrandSelectionRemarks: item.BrandSelectionRemarks,
		})
	}

	_, err = s.Manager.UpdateCollectionItemsInBulk(ctx, *c.Id, input)
	return domain.CreateItemUpdateResponse(err)
}

func (s *Service) SearchItemsInCollection(ctx context.Context, query coredomain.SearchQuery, sortBy string, sortDir string, page int, size int) domain.ProfileCollectionItemResponse {
	items, filteredCount, err := s.Manager.SearchItemsInCollection(ctx, query, sortBy, sortDir, page, size)
	return domain.CreateItemResponse(items, filteredCount, page, size, err)
}

func (s *Service) EnrichMetricsAndItems(ctx context.Context, entries []domain.ProfileCollectionEntry) ([]domain.ProfileCollectionEntry, error) {
	var enrichedEntries []domain.ProfileCollectionEntry
	for _, entry := range entries {
		if strings.ToUpper(*entry.Source) == string(constants.SaasAtCollection) {
			query := coredomain.SearchQuery{Filters: []coredomain.SearchFilter{}}
			query.Filters = append(query.Filters, coredomain.SearchFilter{FilterType: "EQ", Field: "collectionId", Value: strconv.FormatInt(*entry.Id, 10)})
			query.Filters = append(query.Filters, coredomain.SearchFilter{FilterType: "EQ", Field: "collectionType", Value: "PROFILE"})
			summary, err1 := s.CollectionAnalyticsManager.FetchCollectionMetricsSummary(ctx, query)
			// profiles, err2 := s.CollectionAnalyticsManager.FetchCollectionProfilesWithMetricsSummary(ctx, query, "reach", "DESC", 1, 5)
			if err1 == nil && summary != nil {
				if summary.TotalProfiles > 0 {
					entry.ReportReady = helpers.ToBool(true)
				}
				entry.Reach = summary.Reach
				entry.CostPerReach = summary.CostPerReach
				entry.TotalPosts = summary.TotalPosts
				entry.EngagementRate = summary.EngagementRate
			}
		}
		enrichedEntries = append(enrichedEntries, entry)
	}
	return enrichedEntries, nil
}

func (s *Service) FetchWinklCollectionInfoById(ctx context.Context, collectionId int64) domain.WinklCollectionInfoResponse {
	collectionInfo, err := s.Manager.FetchWinklCollectionInfoById(ctx, collectionId)
	return domain.CreateCollectionMigrationInfoResponse(collectionInfo, err)
}

func (s *Service) FetchWinklCollectionInfoByShareId(ctx context.Context, collectionShareId string) domain.WinklCollectionInfoResponse {
	collectionInfo, err := s.Manager.FetchWinklCollectionInfoByShareId(ctx, collectionShareId)
	return domain.CreateCollectionMigrationInfoResponse(collectionInfo, err)
}

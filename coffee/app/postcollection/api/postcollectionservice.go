package api

import (
	manager2 "coffee/app/collectionanalytics/manager"
	partnermanager "coffee/app/partnerusage/manager"
	"coffee/app/postcollection/dao"
	"coffee/app/postcollection/domain"
	"coffee/app/postcollection/manager"
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

func getPrimaryDataStore() coredomain.DataStore {
	return coredomain.Postgres
}

type PostCollectionService struct {
	*rest.Service[domain.PostCollectionResponse, domain.PostCollectionEntry, dao.PostCollectionEntity, string]
	Manager                    *manager.PostCollectionManager
	ItemManager                *manager.PostCollectionItemManager
	CollectionAnalyticsManager *manager2.CollectionAnalyticsManager
	PartnerUsageManager        *partnermanager.PartnerUsageManager
}

func (s *PostCollectionService) Init(ctx context.Context) {
	cmanager := manager.CreatePostCollectionManager(ctx)
	service := rest.NewService(ctx, cmanager.Manager, getPrimaryDataStore, s.createPostCollectionResponse, s.createPostCollectionErrorResponse)
	s.Service = service
	s.Manager = cmanager
	s.ItemManager = manager.CreatePostCollectionItemManager(ctx)
	s.CollectionAnalyticsManager = manager2.CreateCollectionAnalyticsManager(ctx)
	s.PartnerUsageManager = partnermanager.CreatePartnerUsageManager(ctx)
}

func CreatePostCollectionService(ctx context.Context) *PostCollectionService {
	service := &PostCollectionService{}
	service.Init(ctx)
	return service
}

func (s *PostCollectionService) createPostCollectionResponse(collections []domain.PostCollectionEntry, filterCount int64, nextCursor string, message string) domain.PostCollectionResponse {
	return domain.PostCollectionResponse{
		Status: coredomain.Status{
			Status:        "SUCCESS",
			Message:       message,
			FilteredCount: filterCount,
			HasNextPage:   !(nextCursor == ""),
			NextCursor:    nextCursor,
		},
		Collections: &collections,
	}
}

func (s *PostCollectionService) createPostCollectionErrorResponse(err error, code int) domain.PostCollectionResponse {
	return domain.PostCollectionResponse{
		Status: coredomain.Status{
			Status:        "ERROR",
			Message:       err.Error(),
			FilteredCount: 0,
			HasNextPage:   false,
			NextCursor:    "",
		},
		Collections: nil,
	}
}

func (s *PostCollectionService) createPostCollectionItemResponse(collectionItems *[]domain.PostCollectionItemEntry, filterCount int64, nextCursor string, message string) domain.PostCollectionItemResponse {
	return domain.PostCollectionItemResponse{
		Status: coredomain.Status{
			Status:        "SUCCESS",
			Message:       message,
			FilteredCount: filterCount,
			HasNextPage:   !(nextCursor == ""),
			NextCursor:    nextCursor,
		},
		CollectionItems: collectionItems,
	}
}

func (s *PostCollectionService) createPostCollectionItemErrorResponse(err error, code int) domain.PostCollectionItemResponse {
	return domain.PostCollectionItemResponse{
		Status: coredomain.Status{
			Status:        "ERROR",
			Message:       err.Error(),
			FilteredCount: 0,
			HasNextPage:   false,
			NextCursor:    "",
		},
		CollectionItems: nil,
	}
}

func (s *PostCollectionService) Create(ctx context.Context, entry *domain.PostCollectionEntry) domain.PostCollectionResponse {
	appCtx := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)
	if (appCtx.PlanType == nil || *appCtx.PlanType == constants.FreePlan) && entry.Source != nil && *entry.Source != string(constants.GccCampaignCollection) {
		// check no of collections partner has created till now
		count, err := s.Manager.CountCollectionsForPartner(ctx, *appCtx.PartnerId)
		if err != nil {
			return s.createPostCollectionErrorResponse(err, 0)
		}

		if count >= 5 {
			return s.createPostCollectionErrorResponse(errors.New("you are not allowed to create more than 5 collections. Please upgrade your account"), 0)
		}

		// check limit for no of collections and no of items that can be added
		if entry.JobId == nil && entry.Items != nil && len(entry.Items) > 10 {
			return s.createPostCollectionErrorResponse(errors.New("you are not allowed to add more than 10 posts in one collection. Please upgrade your account"), 0)
		}
	}
	e, err := s.Manager.Create(ctx, entry)
	if err != nil {
		return s.createPostCollectionErrorResponse(err, 0)
	}
	if appCtx.PlanType != nil && (*appCtx.PlanType == constants.PaidPlan || *appCtx.PlanType == constants.SaasPlan) && appCtx.PartnerId != nil && entry.Source != nil && *entry.Source == string(constants.SaasCollection) {
		activityId := ""
		platform := ""
		if len(entry.Items) > 0 && entry.Items[0].Platform != nil {
			platform = *entry.Items[0].Platform
		}
		campaignName := entry.Name
		activityMeta := s.ItemManager.MakePostCollectionActivityMeta(*e.Id, platform, *campaignName)
		_, err := s.PartnerUsageManager.LogPaidActivity(ctx, constants.CampaignReport, constants.CampaignReportActivity, 0, *appCtx.PartnerId, *appCtx.AccountId, activityId, activityMeta, "", string(*appCtx.PlanType))
		if err != nil {
			return s.createPostCollectionErrorResponse(err, 0)
		}
	} else if appCtx.PlanType != nil && *appCtx.PlanType == constants.FreePlan && appCtx.PartnerId != nil && entry.Source != nil && *entry.Source == string(constants.SaasCollection) {
		activityId := ""
		platform := ""
		if len(entry.Items) > 0 && entry.Items[0].Platform != nil {
			platform = *entry.Items[0].Platform
		}
		campaignName := entry.Name
		activityMeta := s.ItemManager.MakePostCollectionActivityMeta(*e.Id, platform, *campaignName)
		_, err := s.PartnerUsageManager.LogFreeActivity(ctx, constants.CampaignReport, constants.CampaignReportActivity, 0, *appCtx.PartnerId, *appCtx.AccountId, activityId, activityMeta, "", string(*appCtx.PlanType))
		if err != nil {
			return s.createPostCollectionErrorResponse(err, 0)
		}
	}
	return s.createPostCollectionResponse([]domain.PostCollectionEntry{*e}, 1, "", "Record Created Successfully")
}

func (s *PostCollectionService) FindById(ctx context.Context, collectionId string) domain.PostCollectionResponse {
	e, err := s.Manager.FindById(ctx, collectionId, false)
	if err != nil {
		return s.createPostCollectionErrorResponse(err, 0)
	}
	return s.createPostCollectionResponse([]domain.PostCollectionEntry{*e}, 1, "", "Record Retrieved Successfully")
}

func (s *PostCollectionService) FindByShareId(ctx context.Context, shareId string) domain.PostCollectionResponse {
	e, err := s.Manager.FindByShareId(ctx, shareId)
	if err != nil {
		return s.createPostCollectionErrorResponse(err, 0)
	}
	return s.createPostCollectionResponse([]domain.PostCollectionEntry{*e}, 1, "", "Record Retrieved Successfully")
}

func (s *PostCollectionService) SearchCollections(ctx context.Context, searchQuery coredomain.SearchQuery, sortBy string, sortDir string, page int, size int) domain.PostCollectionResponse {
	entries, filteredCount, err := s.Manager.Search(ctx, searchQuery, sortBy, sortDir, page, size)
	if err != nil {
		return s.createPostCollectionErrorResponse(err, 0)
	}
	entries, _ = s.EnrichMetricsAndItems(ctx, entries)
	var nextCursor string
	if len(entries) >= size {
		nextCursor = strconv.Itoa(page + 1)
	} else {
		nextCursor = ""
	}
	return s.createPostCollectionResponse(entries, filteredCount, nextCursor, "Record(s) Retrieved Successfully")
}

func (s *PostCollectionService) AddItemsToCollection(ctx context.Context, collectionId string, items []domain.PostCollectionItemEntry) domain.PostCollectionItemResponse {
	collection, err := s.Manager.FindById(ctx, collectionId, false)
	if err != nil {
		return s.createPostCollectionItemErrorResponse(err, 0)
	}
	if collection == nil {
		return s.createPostCollectionItemErrorResponse(errors.New("Invalid Collection Id"), 0)
	}

	createdItems, err := s.ItemManager.AddItemsToCollection(ctx, *collection, items)
	if err != nil {
		return s.createPostCollectionItemErrorResponse(err, 0)
	}
	itemsCreated := int64(0)
	if createdItems != nil {
		itemsCreated = int64(len(*createdItems))
	}

	return s.createPostCollectionItemResponse(createdItems, itemsCreated, "", "Items(s) Added Successfully")
}

func (s *PostCollectionService) DeleteItemsFromCollection(ctx context.Context, collectionId string, items []domain.PostCollectionItemEntry) domain.PostCollectionItemResponse {
	collection, err := s.Manager.FindById(ctx, collectionId, false)
	if err != nil {
		return s.createPostCollectionItemErrorResponse(err, 0)
	}
	if collection == nil {
		return s.createPostCollectionItemErrorResponse(errors.New("Invalid Collection Id"), 0)
	}

	deletedItems, err := s.ItemManager.RemoveItemsFromCollection(ctx, items)
	if err != nil {
		return s.createPostCollectionItemErrorResponse(err, 0)
	}
	itemsRemoved := int64(0)
	if deletedItems != nil {
		itemsRemoved = int64(len(*deletedItems))
	}

	return s.createPostCollectionItemResponse(deletedItems, itemsRemoved, "", "Items(s) Removed Successfully")
}

func (s *PostCollectionService) UpdatePostCollectionItemById(ctx context.Context, collectionId string, itemId int64, input domain.PostCollectionItemEntry) domain.PostCollectionItemResponse {
	collection, err := s.Manager.FindById(ctx, collectionId, false)
	if err != nil {
		return s.createPostCollectionItemErrorResponse(err, 0)
	}
	if collection == nil {
		return s.createPostCollectionItemErrorResponse(errors.New("Invalid Collection Id"), 0)
	}

	updatedEntry, err := s.ItemManager.Update(ctx, itemId, &input)
	if err != nil {
		return s.createPostCollectionItemErrorResponse(err, 0)
	}
	return s.createPostCollectionItemResponse(&[]domain.PostCollectionItemEntry{*updatedEntry}, 1, "", "Item Updated Successfully")
}

func (s *PostCollectionService) UpdatePostCollectionItemByShortCode(ctx context.Context, collectionId string, platform string, shortCode string, input domain.PostCollectionItemEntry) domain.PostCollectionItemResponse {
	collection, err := s.Manager.FindById(ctx, collectionId, false)
	if err != nil {
		return s.createPostCollectionItemErrorResponse(err, 0)
	}
	if collection == nil {
		return s.createPostCollectionItemErrorResponse(errors.New("Invalid Collection Id"), 0)
	}

	updatedEntry, err := s.ItemManager.UpdateItemByPlatformAndShortCode(ctx, collectionId, platform, shortCode, input)
	if err != nil {
		return s.createPostCollectionItemErrorResponse(err, 0)
	}
	return s.createPostCollectionItemResponse(&[]domain.PostCollectionItemEntry{*updatedEntry}, 1, "", "Item Updated Successfully")
}

func (s *PostCollectionService) DeleteItemsFromCollectionByShortCode(ctx context.Context, collectionId string, items []domain.PostCollectionItemEntry) domain.PostCollectionItemResponse {
	collection, err := s.Manager.FindById(ctx, collectionId, false)
	if err != nil {
		return s.createPostCollectionItemErrorResponse(err, 0)
	}
	if collection == nil {
		return s.createPostCollectionItemErrorResponse(errors.New("Invalid Collection Id"), 0)
	}

	deletedItems, err := s.ItemManager.RemoveItemsFromCollectionByShortCode(ctx, collectionId, items)
	if err != nil {
		return s.createPostCollectionItemErrorResponse(err, 0)
	}
	itemsRemoved := int64(0)
	if deletedItems != nil {
		itemsRemoved = int64(len(*deletedItems))
	}

	return s.createPostCollectionItemResponse(deletedItems, itemsRemoved, "", "Items(s) Removed Successfully")
}

func (r *PostCollectionService) Update(ctx context.Context, id string, entry *domain.PostCollectionEntry) domain.PostCollectionResponse {
	_, err := r.Manager.Update(ctx, id, entry)
	if err != nil {
		return r.createPostCollectionErrorResponse(err, 0)
	}
	return r.createPostCollectionResponse(nil, 1, "", "Record Updated Successfully")
}

func (s *PostCollectionService) SearchPostCollectionItems(ctx context.Context, collectionId string, searchQuery coredomain.SearchQuery, sortBy string, sortDir string, page int, size int) domain.PostCollectionItemResponse {

	collection, err := s.Manager.FindById(ctx, collectionId, false)
	if err != nil {
		return s.createPostCollectionItemErrorResponse(err, 0)
	}
	if collection == nil {
		return s.createPostCollectionItemErrorResponse(errors.New("Invalid Collection Id"), 0)
	}

	searchQuery.Filters = append(searchQuery.Filters, coredomain.SearchFilter{FilterType: "EQ", Field: "post_collection_id", Value: collectionId})
	entries, filteredCount, err := s.ItemManager.Manager.Search(ctx, searchQuery, sortBy, sortDir, page, size)
	if err != nil {
		return s.createPostCollectionItemErrorResponse(err, 0)
	}
	var nextCursor string
	if len(entries) >= size {
		nextCursor = strconv.Itoa(page + 1)
	} else {
		nextCursor = ""
	}
	return s.createPostCollectionItemResponse(&entries, filteredCount, nextCursor, "Record(s) Retrieved Successfully")
}

func (s *PostCollectionService) EnrichMetricsAndItems(ctx context.Context, entries []domain.PostCollectionEntry) ([]domain.PostCollectionEntry, error) {
	var enrichedEntries []domain.PostCollectionEntry
	for _, entry := range entries {
		if strings.ToUpper(*entry.Source) == "SAAS" {
			query := coredomain.SearchQuery{Filters: []coredomain.SearchFilter{}}
			query.Filters = append(query.Filters, coredomain.SearchFilter{FilterType: "EQ", Field: "collectionId", Value: *entry.Id})
			query.Filters = append(query.Filters, coredomain.SearchFilter{FilterType: "EQ", Field: "collectionType", Value: "POST"})
			summary, err1 := s.CollectionAnalyticsManager.FetchCollectionMetricsSummary(ctx, query)
			posts, err2 := s.CollectionAnalyticsManager.FetchCollectionPostsWithMetricsSummary(ctx, query, "reach", "DESC", 1, 5)
			if err1 == nil && summary != nil {
				entry.TotalPosts = summary.TotalPosts
				if entry.TotalPosts > 0 {
					entry.ReportReady = helpers.ToBool(true)
				}
				entry.Reach = summary.Reach
				entry.CostPerReach = summary.CostPerReach
			}

			if err2 == nil && posts != nil {
				var items []domain.PostCollectionItemEntry
				for _, post := range posts {
					items = append(items, domain.PostCollectionItemEntry{
						Platform:         &post.Platform,
						ShortCode:        &post.PostShortCode,
						PostType:         &post.PostType,
						PostedByHandle:   &post.ProfileHandle,
						HandleProfilePic: post.ProfilePic,
						PostTitle:        &post.PostTitle,
						PostLink:         &post.PostLink,
						PostThumbnail:    post.PostThumbnail,
					})
				}
				entry.Items = items
				if len(items) == 0 {
					c, _ := s.ItemManager.CountCollectionItemsInACollection(ctx, *entry.Id)
					entry.TotalPosts = c
				}
			}
		} else {
			c, _ := s.ItemManager.CountCollectionItemsInACollection(ctx, *entry.Id)
			entry.TotalPosts = c
			entry.ReportReady = helpers.ToBool(true)
		}
		enrichedEntries = append(enrichedEntries, entry)
	}
	return enrichedEntries, nil
}

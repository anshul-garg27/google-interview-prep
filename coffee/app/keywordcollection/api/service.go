package api

import (
	"coffee/app/keywordcollection/dao"
	"coffee/app/keywordcollection/domain"
	"coffee/app/keywordcollection/manager"
	partnermanager "coffee/app/partnerusage/manager"
	"coffee/constants"
	"coffee/core/appcontext"
	coredomain "coffee/core/domain"
	"coffee/core/rest"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	log "github.com/sirupsen/logrus"
)

func getPrimaryDataStore() coredomain.DataStore {
	return coredomain.Postgres
}

type KeywordCollectionService struct {
	*rest.Service[domain.KeywordCollectionResponse, domain.KeywordCollectionEntry, dao.KeywordCollectionEntity, string]
	Manager             *manager.KeywordCollectionManager
	PartnerUsageManager *partnermanager.PartnerUsageManager
	Cache               map[string]interface{}
}

func (s *KeywordCollectionService) Init(ctx context.Context) {
	cmanager := manager.CreateKeywordCollectionManager(ctx)
	service := rest.NewService(ctx, cmanager.Manager, getPrimaryDataStore, s.createKeywordCollectionResponse, s.createKeywordCollectionErrorResponse)
	s.Service = service
	s.Manager = cmanager
	s.PartnerUsageManager = partnermanager.CreatePartnerUsageManager(ctx)
	s.Cache = make(map[string]interface{})
}

func CreateKeywordCollectionService(ctx context.Context) *KeywordCollectionService {
	service := &KeywordCollectionService{}
	service.Init(ctx)
	return service
}

func (s *KeywordCollectionService) createKeywordCollectionResponse(collections []domain.KeywordCollectionEntry, filterCount int64, nextCursor string, message string) domain.KeywordCollectionResponse {
	return domain.KeywordCollectionResponse{
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

func (s *KeywordCollectionService) createKeywordCollectionPostResponse(posts []domain.KeywordCollectionPost, filterCount int64, nextCursor string, message string) domain.KeywordCollectionPostResponse {
	return domain.KeywordCollectionPostResponse{
		Status: coredomain.Status{
			Status:        "SUCCESS",
			Message:       message,
			FilteredCount: filterCount,
			HasNextPage:   !(nextCursor == ""),
			NextCursor:    nextCursor,
		},
		Posts: &posts,
	}
}

func (s *KeywordCollectionService) createKeywordCollectionProfileResponse(profiles []domain.KeywordCollectionProfile, filterCount int64, nextCursor string, message string) domain.KeywordCollectionProfileResponse {
	return domain.KeywordCollectionProfileResponse{
		Status: coredomain.Status{
			Status:        "SUCCESS",
			Message:       message,
			FilteredCount: filterCount,
			HasNextPage:   !(nextCursor == ""),
			NextCursor:    nextCursor,
		},
		Profiles: &profiles,
	}
}

func (s *KeywordCollectionService) createKeywordCollectionErrorResponse(err error, code int) domain.KeywordCollectionResponse {
	return domain.KeywordCollectionResponse{
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

func (s *KeywordCollectionService) createKeywordCollectionPostErrorResponse(err error, code int) domain.KeywordCollectionPostResponse {
	return domain.KeywordCollectionPostResponse{
		Status: coredomain.Status{
			Status:        "ERROR",
			Message:       err.Error(),
			FilteredCount: 0,
			HasNextPage:   false,
			NextCursor:    "",
		},
		Posts: nil,
	}
}

func (s *KeywordCollectionService) createKeywordCollectionProfileErrorResponse(err error, code int) domain.KeywordCollectionProfileResponse {
	return domain.KeywordCollectionProfileResponse{
		Status: coredomain.Status{
			Status:        "ERROR",
			Message:       err.Error(),
			FilteredCount: 0,
			HasNextPage:   false,
			NextCursor:    "",
		},
		Profiles: nil,
	}
}

func (s *KeywordCollectionService) Create(ctx context.Context, entry *domain.KeywordCollectionEntry) domain.KeywordCollectionResponse {
	appCtx := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)
	if appCtx.PlanType == nil || *appCtx.PlanType == constants.FreePlan {
		// check no of collections partner has created till now
		count, err := s.Manager.CountCollectionsForPartner(ctx, *appCtx.PartnerId)
		if err != nil {
			return s.createKeywordCollectionErrorResponse(err, 0)
		}

		if count >= 5 {
			return s.createKeywordCollectionErrorResponse(errors.New("you are not allowed to create more than 5 collections. Please upgrade your account"), 0)
		}

		// check limit for no of collections and no of items that can be added
		if entry.JobId == nil && entry.Keywords != nil && len(*entry.Keywords) > 5 {
			return s.createKeywordCollectionErrorResponse(errors.New("you are not allowed to add more than 5 keywords in one collection. Please upgrade your account"), 0)
		}
	}

	e, err := s.Manager.Create(ctx, entry)
	if err != nil {
		return s.createKeywordCollectionErrorResponse(err, 0)
	}

	if appCtx.PlanType != nil && (*appCtx.PlanType == constants.PaidPlan || *appCtx.PlanType == constants.SaasPlan) && appCtx.PartnerId != nil {
		activityId := ""
		activityMeta := makeKeywordCollectionActivityMeta(*e.Id, *entry.Platform, *entry.Name, entry.Keywords)
		_, err := s.PartnerUsageManager.LogPaidActivity(ctx, constants.TopicResearch, constants.TopicResearchActivity, int64(len(*entry.Keywords)), *appCtx.PartnerId, *appCtx.AccountId, activityId, activityMeta, "", string(*appCtx.PlanType))
		if err != nil {
			return s.createKeywordCollectionErrorResponse(err, 0)
		}
	} else if appCtx.PlanType != nil && *appCtx.PlanType == constants.FreePlan && appCtx.PartnerId != nil {
		activityId := ""
		activityMeta := makeKeywordCollectionActivityMeta(*e.Id, *entry.Platform, *entry.Name, entry.Keywords)

		_, err := s.PartnerUsageManager.LogFreeActivity(ctx, constants.TopicResearch, constants.TopicResearchActivity, int64(len(*entry.Keywords)), *appCtx.PartnerId, *appCtx.AccountId, activityId, activityMeta, "", string(*appCtx.PlanType))
		if err != nil {
			return s.createKeywordCollectionErrorResponse(err, 0)
		}
	}
	return s.createKeywordCollectionResponse([]domain.KeywordCollectionEntry{*e}, 1, "", "Record Created Successfully")
}

func (s *KeywordCollectionService) CreateReport(ctx context.Context, collectionId string) domain.KeywordCollectionResponse {
	e, err := s.Manager.CreateReport(ctx, collectionId)
	if err != nil {
		return s.createKeywordCollectionErrorResponse(err, 0)
	}
	return s.createKeywordCollectionResponse([]domain.KeywordCollectionEntry{*e}, 1, "", "Report Job Created Successfully")
}

func (s *KeywordCollectionService) FindById(ctx context.Context, collectionId string) domain.KeywordCollectionResponse {
	e, err := s.Manager.FindById(ctx, collectionId, false)
	if err != nil {
		return s.createKeywordCollectionErrorResponse(err, 0)
	}
	return s.createKeywordCollectionResponse([]domain.KeywordCollectionEntry{*e}, 1, "", "Record Retrieved Successfully")
}

func (s *KeywordCollectionService) FindByShareId(ctx context.Context, shareId string) domain.KeywordCollectionResponse {
	e, err := s.Manager.FindByShareId(ctx, shareId)
	if err != nil {
		return s.createKeywordCollectionErrorResponse(err, 0)
	}
	return s.createKeywordCollectionResponse([]domain.KeywordCollectionEntry{*e}, 1, "", "Record Retrieved Successfully")
}

func (s *KeywordCollectionService) SearchCollections(ctx context.Context, searchQuery coredomain.SearchQuery, sortBy string, sortDir string, page int, size int) domain.KeywordCollectionResponse {
	entries, filteredCount, err := s.Manager.Search(ctx, searchQuery, sortBy, sortDir, page, size)
	if err != nil {
		return s.createKeywordCollectionErrorResponse(err, 0)
	}
	entries, _ = s.EnrichMetricsAndItems(ctx, entries)
	var nextCursor string
	if len(entries) >= size {
		nextCursor = strconv.Itoa(page + 1)
	} else {
		nextCursor = ""
	}
	return s.createKeywordCollectionResponse(entries, filteredCount, nextCursor, "Record(s) Retrieved Successfully")
}

func (s *KeywordCollectionService) GetCollectionReport(ctx context.Context, searchQuery coredomain.SearchQuery) domain.KeywordCollectionResponse {
	var collectionId string
	collectionIdFilter := rest.GetFilterForKey(searchQuery.Filters, "collectionId")
	collectionId = collectionIdFilter.Value
	entry, err := s.Manager.GetCollectionReport(ctx, collectionId)
	if err != nil || entry == nil {
		return s.createKeywordCollectionErrorResponse(err, 0)
	}
	return s.createKeywordCollectionResponse([]domain.KeywordCollectionEntry{*entry}, 1, "", "Record(s) Retrieved Successfully")
}

func (s *KeywordCollectionService) GetCollectionPosts(ctx context.Context, searchQuery coredomain.SearchQuery, sortBy string, sortDir string, page int, size int) domain.KeywordCollectionPostResponse {
	entries, filteredCount, err := s.Manager.GetCollectionPosts(ctx, searchQuery, sortBy, sortDir, page, size)
	if err != nil {
		return s.createKeywordCollectionPostErrorResponse(err, 0)
	}
	return s.createKeywordCollectionPostResponse(entries, filteredCount, "", "Record(s) Retrieved Successfully")
}

func (s *KeywordCollectionService) GetCollectionProfiles(ctx context.Context, searchQuery coredomain.SearchQuery, sortBy string, sortDir string, page int, size int) domain.KeywordCollectionProfileResponse {
	entries, filteredCount, err := s.Manager.GetCollectionProfiles(ctx, searchQuery, sortBy, sortDir, page, size)
	if err != nil {
		return s.createKeywordCollectionProfileErrorResponse(err, 0)
	}
	return s.createKeywordCollectionProfileResponse(entries, filteredCount, "", "Record(s) Retrieved Successfully")
}

func (r *KeywordCollectionService) Update(ctx context.Context, id string, entry *domain.KeywordCollectionEntry) domain.KeywordCollectionResponse {
	_, err := r.Manager.Update(ctx, id, entry)
	if err != nil {
		return r.createKeywordCollectionErrorResponse(err, 0)
	}
	return r.createKeywordCollectionResponse(nil, 1, "", "Record Updated Successfully")
}

func (s *KeywordCollectionService) EnrichMetricsAndItems(ctx context.Context, entries []domain.KeywordCollectionEntry) ([]domain.KeywordCollectionEntry, error) {
	for i := range entries {
		entries[i].ReportStatus = "PENDING"
		cacheKey := fmt.Sprintf("collection-report-%s", *entries[i].Id)
		cacheValue, present := s.Cache[cacheKey]
		if present {
			entries[i].Report = cacheValue.(*domain.KeywordCollectionReport)
		} else {
			enrichedCollection, _ := s.Manager.GetCollectionReportPartial(ctx, &entries[i])
			if enrichedCollection != nil && enrichedCollection.Report != nil {
				entries[i].Report = enrichedCollection.Report
				s.Cache[cacheKey] = enrichedCollection.Report
				entries[i].ReportStatus = enrichedCollection.ReportStatus
			}
		}
	}
	return entries, nil
}

func makeKeywordCollectionActivityMeta(id string, platform string, name string, keywords *[]string) string {
	var requestMetaStr string
	var requestMeta domain.TopicResearchActivityMeta
	requestMeta.Id = id
	requestMeta.Platform = platform
	requestMeta.Result.Name = name
	requestMeta.Result.Keywords = keywords

	requestMetaJson, err := json.Marshal(requestMeta)
	if err != nil {
		log.Error(err)
		return requestMetaStr
	}
	requestMetaStr = string(requestMetaJson)
	return requestMetaStr
}

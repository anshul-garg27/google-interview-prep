package api

import (
	"coffee/app/discovery/domain"
	"coffee/app/discovery/manager"
	"coffee/app/partnerusage/dao"
	partnermanager "coffee/app/partnerusage/manager"
	"coffee/constants"
	"coffee/core/appcontext"
	coredomain "coffee/core/domain"
	"coffee/core/rest"
	"context"
	"encoding/json"
	"strconv"

	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
)

type Service struct {
	Manager             *manager.SearchManager
	ManagerTimeSeries   *manager.TimeSeriesManager
	ManagerHashtags     *manager.HashtagsManager
	ManagerAudience     *manager.AudienceManager
	ManagerLocation     *manager.LocationManager
	PartnerUsageManager *partnermanager.PartnerUsageManager
}

func NewSocialDiscoveryService() *Service {
	ctx := context.TODO()
	seachmanager := manager.CreateSearchManager(ctx)
	managerTimeSeries := manager.CreateManagerTimeSeries(ctx)
	managerHashtags := manager.CreateHashtagsManager(ctx)
	managerAudience := manager.CreateAudienceManager(ctx)
	locationsDao := manager.CreateLocationManager(ctx)
	partnermanager := partnermanager.CreatePartnerUsageManager(ctx)
	return &Service{
		Manager:             seachmanager,
		ManagerTimeSeries:   managerTimeSeries,
		ManagerHashtags:     managerHashtags,
		ManagerAudience:     managerAudience,
		ManagerLocation:     locationsDao,
		PartnerUsageManager: partnermanager,
	}
}

func (s *Service) FindByProfileId(ctx context.Context, profileId string, linkedSource string) domain.SocialProfileResponse {
	var profiles []coredomain.Profile
	profile, err := s.Manager.FindByProfileId(ctx, profileId, linkedSource)
	filteredCount := int64(0)
	if profile != nil {
		filteredCount = int64(1)
		profiles = append(profiles, *profile)
	}
	profiles = blockProfileDataForFreeUsers(ctx, profiles)
	return domain.CreateSoialResponse(&profiles, filteredCount, 1, 1, err, int64(0))
}

func (s *Service) FindByPlatformProfileId(ctx context.Context, platform string, platformProfileId int64, linkedSource string, campaignProfileJoin bool) domain.SocialProfileResponse {
	var profiles []coredomain.Profile
	profile, err := s.Manager.FindByPlatformProfileId(ctx, platform, platformProfileId, linkedSource, campaignProfileJoin)
	filteredCount := int64(0)
	if err != nil {
		return domain.CreateSoialResponse(&profiles, filteredCount, 1, 1, err, int64(0))

	}
	appCtx := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)
	if appCtx.PlanType != nil && (*appCtx.PlanType == constants.PaidPlan || *appCtx.PlanType == constants.SaasPlan) && appCtx.PartnerId != nil && *appCtx.PartnerId != -1 {
		key := "profile-" + platform + "-" + strconv.FormatInt(platformProfileId, 10)
		profileTrack, _ := s.PartnerUsageManager.FindProfilePageByKey(ctx, *appCtx.PartnerId, key)
		if profileTrack == nil {
			profile.Locked = true
		}
	}
	filteredCount = int64(1)
	profiles = append(profiles, *profile)

	profiles = blockOverviewDataForFreeUsers(ctx, profiles)
	return domain.CreateSoialResponse(&profiles, filteredCount, 1, 1, err, int64(0))
}

func (s *Service) FindAudienceByPlatformProfileId(ctx context.Context, platform string, platformProfileId int64) domain.ProfileAudienceRespnse {
	var audienceData []domain.SocialProfileAudienceInfoEntry
	audience, err := s.ManagerAudience.FindAudienceByPlatformProfileId(ctx, platform, platformProfileId)
	filteredCount := int64(0)
	if audience != nil {
		filteredCount = int64(1)
		audienceData = append(audienceData, *audience)
	}
	profilePageTrackKey := "profile-" + platform + "-" + strconv.FormatInt(platformProfileId, 10)
	audienceData = blockAudienceDataForFreeUsers(ctx, audienceData, profilePageTrackKey)
	return domain.CreateProfileAudienceResponse(err, &audienceData, filteredCount, "")
}

func (s *Service) FindTimeSeriesByPlatformProfileId(ctx context.Context, platform string, platformProfileId int64, searchQuery coredomain.SearchQuery, sortBy string, sortDir string, page int, size int) domain.ProfileGrowthRespnse {
	var growthData []domain.Growth
	var filteredCount int64
	var err error
	profilePageTrackKey := "profile-" + platform + "-" + strconv.FormatInt(platformProfileId, 10)
	err = checkTimeSeriesValidityForFreeUsers(ctx, profilePageTrackKey)
	if err != nil {
		return domain.CreateGrowthResponse(err, &growthData, filteredCount, "")
	}

	growth, filteredCount, err := s.ManagerTimeSeries.FindTimeSeriesByPlatformProfileId(ctx, platform, platformProfileId, searchQuery, sortBy, sortDir, page, size)
	if growth != nil {
		growthData = append(growthData, *growth)
		filteredCount = int64(1)
	}
	return domain.CreateGrowthResponse(err, &growthData, filteredCount, "")
}

func (s *Service) FindContentByPlatformProfileId(ctx context.Context, platform string, platformProfileId int64, searchQuery coredomain.SearchQuery, sortBy string, sortDir string, page int, size int, linkedSource string) domain.ProfileContentResponse {
	var contentData []domain.Content
	content, filteredCount, err := s.Manager.FindContentByPlatformProfileId(ctx, platform, platformProfileId, searchQuery, sortBy, sortDir, page, size, linkedSource)
	if content != nil {
		contentData = append(contentData, *content)
	}
	return domain.CreateContentResponse(err, &contentData, filteredCount, "")
}

func (s *Service) FindHashTagsByPlatformProfileId(ctx context.Context, platform string, platformProfileId int64) domain.ProfileHashtagResponse {
	var hashtagsData []domain.SocialProfileHashatagsEntry

	hashtags, err := s.ManagerHashtags.FindHashTagsByPlatformProfileId(ctx, platform, platformProfileId)
	filteredCount := int64(0)
	if hashtags != nil {
		hashtagsData = append(hashtagsData, *hashtags)
		filteredCount = int64(1)
	}

	return domain.CreateHashtagResponse(err, &hashtagsData, filteredCount, "")
}

func (s *Service) FindSimilarAccountsByPlatformProfileId(ctx context.Context, platform string, platformProfileId int64, searchQuery coredomain.SearchQuery, sortBy string, sortDir string, page int, size int, linkedSource string) domain.SocialProfileResponse {
	profiles, filteredCount, err := s.Manager.FindSimilarAccountsByPlatformProfileId(ctx, platform, platformProfileId, searchQuery, sortBy, sortDir, page, size, linkedSource)
	return domain.CreateSoialResponse(profiles, filteredCount, 1, 1, err, int64(0))
}

func (s *Service) SearchByPlatform(ctx context.Context, platform string, searchQuery coredomain.SearchQuery, sortBy string, sortDir string, page int, size int, linkedSource string, campaignProfileJoin bool, activityId string, computeCounts bool) domain.SocialProfileResponse {
	var profiles *[]coredomain.Profile
	var filteredCount int64
	var entity *dao.AcitivityTrackerEntity
	var err error
	var blockFreeUsers, blockSaasUsers bool

	searchQuery = rest.ParseGenderAgeFilters(searchQuery)
	searchQueryForActivityMeta := searchQuery

	if linkedSource != string(constants.GCCPlatform) {
		blockFreeUsers, blockSaasUsers = checkForPaidRequestParams(searchQuery, page, size)
		err = validateDiscoverySearch(ctx, page, size, searchQuery, blockFreeUsers, blockSaasUsers)
		if err != nil {
			return domain.CreateSoialResponse(profiles, filteredCount, page, size, err, int64(0))
		}
		searchQuery.Filters = append([]coredomain.SearchFilter{{
			FilterType: "EQ",
			Field:      "country",
			Value:      "IN",
		}}, searchQuery.Filters...)
	}
	if computeCounts {
		computeCounts = shouldComputeCounts(searchQuery, page)
	}
	profiles, filteredCount, err = s.Manager.SearchByPlatform(ctx, platform, searchQuery, sortBy, sortDir, page, size, linkedSource, campaignProfileJoin, computeCounts)
	if err != nil {
		return domain.CreateSoialResponse(profiles, filteredCount, page, size, err, int64(0))
	}

	appCtx := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)

	if page == 1 && appCtx.PlanType != nil && (*appCtx.PlanType == constants.PaidPlan || *appCtx.PlanType == constants.SaasPlan) && appCtx.PartnerId != nil && linkedSource != string(constants.GccLinkedProfile) {
		activityMeta := makeDiscoveryActivityMeta(size, page, sortBy, sortDir, linkedSource, strconv.FormatBool(campaignProfileJoin), searchQueryForActivityMeta, platform, int(filteredCount))
		entity, err = s.PartnerUsageManager.LogPaidActivity(ctx, constants.Discovery, constants.DiscoveryActivity, int64(1), *appCtx.PartnerId, *appCtx.AccountId, activityId, activityMeta, "", string(*appCtx.PlanType))
		if err != nil && err.Error() == "limit consumed" && (!blockFreeUsers) {
			entity, err = s.PartnerUsageManager.LogFreeActivity(ctx, constants.Discovery, constants.DiscoveryActivity, int64(1), *appCtx.PartnerId, *appCtx.AccountId, activityId, activityMeta, "", string(*appCtx.PlanType))
			if err != nil {
				return domain.CreateSoialResponse(profiles, int64(0), page, size, err, int64(0))
			}
		} else if err != nil && err.Error() == "limit consumed" && blockFreeUsers {
			return domain.CreateSoialResponse(profiles, int64(0), page, size, err, int64(0))
		}
	} else if page == 1 && appCtx.PlanType != nil && *appCtx.PlanType == constants.FreePlan && appCtx.PartnerId != nil && linkedSource != string(constants.GccLinkedProfile) {
		activityMeta := makeDiscoveryActivityMeta(size, page, sortBy, sortDir, linkedSource, strconv.FormatBool(campaignProfileJoin), searchQueryForActivityMeta, platform, int(filteredCount))
		entity, err = s.PartnerUsageManager.LogFreeActivity(ctx, constants.Discovery, constants.DiscoveryActivity, int64(1), *appCtx.PartnerId, *appCtx.AccountId, activityId, activityMeta, "", string(*appCtx.PlanType))
		if err != nil {
			return domain.CreateSoialResponse(profiles, int64(0), page, size, err, int64(0))
		}
	}
	if entity != nil {
		return domain.CreateSoialResponse(profiles, filteredCount, page, size, err, entity.Id)
	}
	return domain.CreateSoialResponse(profiles, filteredCount, page, size, err, int64(0))
}

func (s *Service) FindByPlatformHandle(ctx context.Context, platform string, handle string, linkedSource string, campaignProfileJoin bool) domain.SocialProfileResponse {
	profile, err := s.Manager.FindProfileByPlatformHandle(ctx, platform, handle, linkedSource, campaignProfileJoin)
	var filteredCount int64
	if profile != nil {
		filteredCount = int64(1)
		return domain.CreateSoialResponse(&[]coredomain.Profile{*profile}, filteredCount, 1, 1, err, int64(0))
	}
	return domain.CreateSoialResponse(&[]coredomain.Profile{*profile}, filteredCount, 1, 1, err, int64(0))
}

func (s *Service) SearchGccLocations(ctx context.Context, name string) domain.GccLocationsResponse {
	locations, err := s.ManagerLocation.SearchGccLocations(ctx, name)
	return domain.CreateGccLocationsResponse(err, locations)
}

func (s *Service) SearchLocations(ctx context.Context, searchQuery coredomain.SearchQuery, sortBy string, sortDir string, page int, size int) domain.LocationsResponse {
	locations, filteredCount, err := s.ManagerLocation.SearchLocations(ctx, searchQuery, sortBy, sortDir, page, size)
	return domain.CreateLocationsResponse(locations, filteredCount, page, size, err)
}

func (s *Service) SearchLocationsNew(ctx context.Context, searchQuery coredomain.SearchQuery, size int) domain.LocationsResponse {
	locations, err := s.ManagerLocation.SearchLocationsNew(ctx, searchQuery, size)
	filteredCount := size
	if len(locations) <= size {
		filteredCount = len(locations)
	}
	return domain.CreateLocationsResponse(&locations, int64(filteredCount), 1, size, err)
}

func makeDiscoveryActivityMeta(size int, cursor int, sortBy string, sortDir string, source string, campaignProfileJoinStr string, query coredomain.SearchQuery, platform string, filteredCount int) string {
	// Make Separate function To Make Meta In Rest
	var requestMetaStr string
	var requestMeta domain.DiscoveryActivityMeta
	requestMeta.QueryParams = make(map[string]string)
	requestMeta.QueryParams["size"] = strconv.Itoa(size)
	requestMeta.QueryParams["cursor"] = strconv.Itoa(cursor)
	requestMeta.QueryParams["sortBy"] = sortBy
	requestMeta.QueryParams["sortDir"] = sortDir
	requestMeta.QueryParams["source"] = source
	requestMeta.QueryParams["campaignProfileJoin"] = campaignProfileJoinStr

	requestMeta.BodyParam = query
	requestMeta.Platform = platform

	resultMap := make(map[string]interface{})
	resultMap["resultCount"] = filteredCount
	requestMeta.Result = resultMap

	requestMetaJson, err := json.Marshal(requestMeta)
	if err != nil {
		log.Error(err)
		return requestMetaStr
	}
	requestMetaStr = string(requestMetaJson)
	return requestMetaStr
}

// func makeDiscoveryAudienceActivityMeta(platform string, platformProfileId int64) string {
// 	var requestMetaStr string
// 	var requestMeta = make(map[string]interface{})
// 	requestMeta["platform"] = platform
// 	requestMeta["platformProfileId"] = platformProfileId
// 	requestMetaJson, err := json.Marshal(requestMeta)
// 	if err != nil {
// 		log.Error(err)
// 		return requestMetaStr
// 	}
// 	requestMetaStr = string(requestMetaJson)
// 	return requestMetaStr
// }

func shouldComputeCounts(query coredomain.SearchQuery, page int) bool {
	computeCounts := false
	campaignIdFilter := rest.GetFilterForKey(query.Filters, "campaignId")
	collectionIdFilter := rest.GetFilterForKey(query.Filters, "collectionId")
	if campaignIdFilter != nil || collectionIdFilter != nil { // Campaign Related Pages
		computeCounts = true
	}
	if collectionIdFilter != nil && collectionIdFilter.FilterType == "NE" { // Disable for eligible influencers
		computeCounts = false
	}
	if campaignIdFilter != nil && campaignIdFilter.FilterType == "NE" { // Disable for eligible influencers
		computeCounts = false
	}
	coreFilters := []string{"campaignId", "collectionId", "onGcc", "followers", "isBlacklisted"}
	nonCoreFilterCount := 0
	for i := range query.Filters {
		if !slices.Contains(coreFilters, query.Filters[i].Field) {
			nonCoreFilterCount += 1
		}
	}
	if nonCoreFilterCount > 0 {
		computeCounts = true
	}
	if page >= 2 {
		computeCounts = false
	}
	return computeCounts
}

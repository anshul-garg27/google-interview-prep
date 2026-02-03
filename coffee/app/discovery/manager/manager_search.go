package manager

import (
	"coffee/app/discovery/dao"
	"coffee/app/discovery/domain"
	beatservice "coffee/client/beat"
	"coffee/constants"
	"coffee/core/appcontext"
	coredomain "coffee/core/domain"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"

	"github.com/pkg/errors"
)

type SearchManager struct {
	instaDao        *dao.InstagramAccountDao
	ytDao           *dao.YoutubeAccountDao
	groupMetricsDao *dao.GroupMetricsDao
}

func CreateSearchManager(ctx context.Context) *SearchManager {
	instaDao := dao.CreateInstagramAccountDao(ctx)
	ytDao := dao.CreateYoutubeAccountDao(ctx)
	groupMetricsDao := dao.CreateGroupMetricsDao(ctx)
	manager := &SearchManager{
		instaDao:        instaDao,
		ytDao:           ytDao,
		groupMetricsDao: groupMetricsDao,
	}
	return manager
}

/*
*
Return Profile Information for a given Profile ID
Render Influencer Page on SaaS and CM
  - Cross Platform Linking Details of Influencer
  - Basic metrics of each
*/
func (m *SearchManager) FindByProfileId(ctx context.Context, id string, linkedSource string) (*coredomain.Profile, error) {
	var profiles []coredomain.Profile
	var isAdmin bool
	appCtx := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)
	if appCtx.PartnerId != nil && *appCtx.PartnerId == -1 {
		isAdmin = true
	}
	var instaAccount *dao.InstagramAccountEntity
	if linkedSource == string(constants.GCCPlatform) {
		instaAccount, _ = m.instaDao.FindByGccProfileId(ctx, id, isAdmin)
	} else {
		instaAccount, _ = m.instaDao.FindByProfileId(ctx, id, isAdmin)
	}

	if instaAccount != nil {
		var linkedChannels []*string
		profile := GetProfileFromInstagramEntity(ctx, *instaAccount, linkedSource)
		profile.Metrics = coredomain.ProfileMetrics{}
		profiles = append(profiles, *profile)

		if linkedSource == string(constants.GCCPlatform) && instaAccount.GccLinkedChannelId != nil && *instaAccount.GccLinkedChannelId != "" {
			linkedChannels = append(linkedChannels, instaAccount.GccLinkedChannelId)
		} else if linkedSource != string(constants.GCCPlatform) && instaAccount.LinkedChannelId != nil && *instaAccount.LinkedChannelId != "" {
			linkedChannels = append(linkedChannels, instaAccount.LinkedChannelId)
		}
		if len(linkedChannels) > 0 {
			linkedDetails, err := m.GetYoutubeDetalsByChannelIds(ctx, linkedChannels)
			if err == nil {

				if linkedSource == string(constants.GCCPlatform) && instaAccount.GccLinkedChannelId != nil && *instaAccount.GccLinkedChannelId != "" {
					profiles = EnrichProfileWithGccLinkedSocials(profiles, linkedDetails)
				} else {
					profiles = EnrichProfileWithSaasLinkedSocials(profiles, linkedDetails)
				}
			}
		}
		return &profiles[0], nil
	}
	var ytAccount *dao.YoutubeAccountEntity
	if linkedSource == string(constants.GCCPlatform) {
		ytAccount, _ = m.ytDao.FindByGccProfileId(ctx, id, isAdmin)
	} else {
		ytAccount, _ = m.ytDao.FindByProfileId(ctx, id, isAdmin)
	}
	if ytAccount != nil {
		var linkedHandle []*string
		profile := GetProfileFromYoutubeEntity(ctx, *ytAccount, linkedSource)

		profile.Metrics = coredomain.ProfileMetrics{}
		profiles = append(profiles, *profile)

		if linkedSource == string(constants.GCCPlatform) && ytAccount.GccLinkedInstagramHandle != nil && *ytAccount.GccLinkedInstagramHandle != "" {
			linkedHandle = append(linkedHandle, ytAccount.GccLinkedInstagramHandle)
		} else if ytAccount.LinkedHandle != nil && *ytAccount.LinkedHandle != "" {
			linkedHandle = append(linkedHandle, ytAccount.LinkedHandle)
		}
		if len(linkedHandle) > 0 {
			linkedDetails, err := m.GetInstagramDetailsByHandle(ctx, linkedHandle)
			if err == nil {

				if linkedSource == string(constants.GCCPlatform) && ytAccount.GccLinkedInstagramHandle != nil && *ytAccount.GccLinkedInstagramHandle != "" {
					profiles = EnrichProfileWithGccLinkedSocials(profiles, linkedDetails)
				} else {
					profiles = EnrichProfileWithSaasLinkedSocials(profiles, linkedDetails)
				}
			}
		}

		return &profiles[0], nil
	}
	if instaAccount == nil && ytAccount == nil {
		return nil, errors.New("profile not found")
	}
	return &profiles[0], nil
}

// FindByPlatformProfileId
/**
Will be called by Identity when
- Identity needs metrics for a given social account (IA or YA)
Logic
	- Return Data present in DB
*/
func (m *SearchManager) FindByPlatformProfileId(ctx context.Context, platform string, platformProfileId int64, linkedSource string, campaignProfileJoin bool) (*coredomain.Profile, error) {
	var profiles []coredomain.Profile
	var profile *coredomain.Profile
	var groupKey string
	var linkedChannels, linkedHandle []*string
	if platform == string(constants.InstagramPlatform) {
		instaAccount, err := m.instaDao.FindById(ctx, platformProfileId, campaignProfileJoin)
		if err != nil {
			return nil, err
		}
		if instaAccount.GroupKey != nil {
			groupKey = *instaAccount.GroupKey
		}

		profile = GetProfileFromInstagramEntity(ctx, *instaAccount, linkedSource)
		profiles = append(profiles, *profile)

		if linkedSource == string(constants.GCCPlatform) && instaAccount.GccLinkedChannelId != nil && *instaAccount.GccLinkedChannelId != "" {
			linkedChannels = append(linkedChannels, instaAccount.GccLinkedChannelId)
		} else if linkedSource != string(constants.GCCPlatform) && instaAccount.LinkedChannelId != nil && *instaAccount.LinkedChannelId != "" {
			linkedChannels = append(linkedChannels, instaAccount.LinkedChannelId)
		}

		if len(linkedChannels) > 0 {
			linkedDetails, err := m.GetYoutubeDetalsByChannelIds(ctx, linkedChannels)
			if err == nil {
				if linkedSource == string(constants.GCCPlatform) && instaAccount.GccLinkedChannelId != nil && *instaAccount.GccLinkedChannelId != "" {
					profiles = EnrichProfileWithGccLinkedSocials(profiles, linkedDetails)
				} else {
					profiles = EnrichProfileWithSaasLinkedSocials(profiles, linkedDetails)
				}
			}
		}
	} else if platform == string(constants.YoutubePlatform) {
		ytAccount, err := m.ytDao.FindById(ctx, platformProfileId, campaignProfileJoin)
		if err != nil {
			return nil, err
		}
		if ytAccount.GroupKey != nil {
			groupKey = *ytAccount.GroupKey

		}
		profile = GetProfileFromYoutubeEntity(ctx, *ytAccount, linkedSource)
		profiles = append(profiles, *profile)
		if ytAccount.LinkedHandle != nil && *ytAccount.LinkedHandle != "" {
			linkedHandle = append(linkedHandle, ytAccount.LinkedHandle)
			linkedDetails, err := m.GetInstagramDetailsByHandle(ctx, linkedHandle)
			if err == nil {
				profiles = EnrichProfileWithSaasLinkedSocials(profiles, linkedDetails)
			}
		}
	}
	if groupKey != "" {
		similarGroupData, err := m.groupMetricsDao.FindByGroupKey(ctx, groupKey)
		if err == nil {
			entry, _ := ToGroupMetricsEntry(*similarGroupData)
			profiles[0].SimilarProfileGroupData = entry
		}
	}
	return &profiles[0], nil
}

func (m *SearchManager) InsertDataFromBeat(ctx context.Context, platform string, handle string) (*beatservice.ProfileResponse, error) {
	beatClient := beatservice.New(ctx)
	handle = strings.TrimSpace(handle)
	if handle == "" {
		message := "handle is empty. Trying to get empty handle from beat"
		log.Error(message)
		return nil, errors.New(message)
	}
	apiData, err := beatClient.FindProfileByHandle(platform, handle)
	if err != nil {
		log.Error("error fetching data from beat: ", err)
		return nil, err
	}
	if apiData == nil {
		message := "error due to beat returning no data"
		log.Error(message)
		return nil, errors.New(message)
	}
	if apiData.Status.Type == "ERROR" || apiData.Profile == nil {
		return nil, errors.New("profile not found")
	}
	return apiData, nil
}

func (m *SearchManager) FindByInstagramHandle(ctx context.Context, platform string, handle string, linkedSource string, campaignProfileJoin bool) ([]coredomain.Profile, error) {
	var profiles []coredomain.Profile
	var instagramEntity dao.InstagramAccountEntity
	instaAccount, err := m.instaDao.FindByHandle(ctx, handle, campaignProfileJoin)

	if err != nil && err.Error() != "record not found" {
		err = fmt.Errorf("FindByInstagramHandle error = %s handle = %s ", err, handle)
		log.Error(err)
		return nil, err
	} else if err != nil && err.Error() == "record not found" {
		beatData, err := m.InsertDataFromBeat(ctx, platform, handle)
		if err != nil {
			return nil, err
		}
		instagramEntity = m.MakeInstagramEntityFromBeat(beatData)
		instaAccount, err := m.instaDao.FindByIgId(ctx, beatData.Profile.ProfileId)
		if err != nil && instaAccount == nil {
			_, err = m.instaDao.Create(ctx, &instagramEntity)
			if err != nil {
				log.Error("error creating instagram account: ", err)
				return nil, err
			}
		} else {
			instagramEntity.ID = instaAccount.ID
		}
		profileId := fmt.Sprintf("IA_%d", instagramEntity.ID)
		instagramEntity.GccProfileId = &profileId
		instagramEntity.ProfileId = instagramEntity.GccProfileId
		instagramEntity.Enabled = true
		flag := false
		instagramEntity.Deleted = &flag
		_, err = m.instaDao.Update(ctx, instagramEntity.ID, &instagramEntity)
		if err != nil {
			log.Error("error setting profile_id and gcc_profile_id of IA: ", err)
			return nil, err
		}
		profile := GetProfileFromInstagramEntity(ctx, instagramEntity, linkedSource)
		profiles = append(profiles, *profile)

		return profiles, nil
	}

	var linkedChannels []*string
	if instaAccount.GccProfileId == nil || *instaAccount.GccProfileId == "" {
		profileId := fmt.Sprintf("IA_%d", instaAccount.ID)
		instaAccount.GccProfileId = &profileId
		instaAccountJson, _ := json.Marshal(instaAccount)
		instaAccount.Enabled = true
		flag := false
		instaAccount.Deleted = &flag
		_, err := m.instaDao.Update(ctx, instaAccount.ID, instaAccount)
		if err != nil {
			err = fmt.Errorf("error setting gcc_profile_id of IA: %s. instagram Account is %s", err, string(instaAccountJson))
			log.Error(err)
			return nil, err
		}
	}
	if instaAccount.ProfileId == nil || *instaAccount.ProfileId == "" {
		profileId := fmt.Sprintf("IA_%d", instaAccount.ID)
		instaAccount.ProfileId = &profileId
		instaAccount.Enabled = true
		flag := false
		instaAccount.Deleted = &flag
		_, err := m.instaDao.Update(ctx, instaAccount.ID, instaAccount)
		if err != nil {
			log.Error("error setting profile_id of IA: ", err)
			return nil, err
		}
	}
	profile := GetProfileFromInstagramEntity(ctx, *instaAccount, linkedSource)
	profiles = append(profiles, *profile)
	if instaAccount.LinkedChannelId != nil && *instaAccount.LinkedChannelId != "" {
		linkedChannels = append(linkedChannels, instaAccount.LinkedChannelId)
		linkedDetails, err := m.GetYoutubeDetalsByChannelIds(ctx, linkedChannels)
		if err == nil {
			profiles = EnrichProfileWithSaasLinkedSocials(profiles, linkedDetails)
		}
	}

	return profiles, nil
}

func (m *SearchManager) FindByYoutubeHandle(ctx context.Context, platform string, handle string, linkedSource string, campaignProfileJoin bool) ([]coredomain.Profile, error) {
	var profiles []coredomain.Profile
	ytAccount, err := m.ytDao.FindByHandle(ctx, handle, campaignProfileJoin)
	if err != nil && err.Error() != "record not found" {
		err = fmt.Errorf("FindByYoutubeHandle error = %s handle = %s ", err, handle)
		log.Error(err)
		return nil, err
	} else if err != nil && err.Error() == "record not found" {
		beatData, err := m.InsertDataFromBeat(ctx, platform, handle)
		if err != nil {
			return nil, err
		}
		youtubeEntity := m.MakeYoutubeEntityFromBeat(beatData)
		_, err = m.ytDao.Create(ctx, &youtubeEntity)
		if err != nil {
			log.Error("error creating youtube account: ", err)
			return nil, err
		}
		profileId := fmt.Sprintf("YA_%d", youtubeEntity.ID)
		youtubeEntity.GccProfileId = &profileId
		youtubeEntity.ProfileId = &profileId
		youtubeEntity.Enabled = true
		flag := false
		youtubeEntity.Deleted = &flag
		_, err = m.ytDao.Update(ctx, youtubeEntity.ID, &youtubeEntity)
		if err != nil {
			log.Error("error setting profile_id and gcc_profile_id of YA: ", err)
			return nil, err
		}
		profile := GetProfileFromYoutubeEntity(ctx, youtubeEntity, linkedSource)
		profiles = append(profiles, *profile)

		return profiles, nil
	}

	var linkedHandle []*string
	if ytAccount.GccProfileId == nil || *ytAccount.GccProfileId == "" {
		profileId := fmt.Sprintf("YA_%d", ytAccount.ID)
		ytAccount.GccProfileId = &profileId
		ytAccountJson, _ := json.Marshal(ytAccount)
		ytAccount.Enabled = true
		flag := false
		ytAccount.Deleted = &flag
		_, err := m.ytDao.Update(ctx, ytAccount.ID, ytAccount)
		if err != nil {
			err = fmt.Errorf("error setting gcc_profile_id of YA: %s. Youtube Account is %s", err, string(ytAccountJson))
			log.Error(err)
			return nil, err
		}
	}
	if ytAccount.ProfileId == nil || *ytAccount.ProfileId == "" {
		profileId := fmt.Sprintf("YA_%d", ytAccount.ID)
		ytAccount.ProfileId = &profileId
		ytAccount.Enabled = true
		flag := false
		ytAccount.Deleted = &flag
		_, err := m.ytDao.Update(ctx, ytAccount.ID, ytAccount)
		if err != nil {
			log.Error("error setting profile_id of YA: ", err)
			return nil, err
		}
	}
	profile := GetProfileFromYoutubeEntity(ctx, *ytAccount, linkedSource)
	profiles = append(profiles, *profile)
	if ytAccount.LinkedHandle != nil && *ytAccount.LinkedHandle != "" {
		linkedHandle = append(linkedHandle, ytAccount.LinkedHandle)
		linkedDetails, err := m.GetInstagramDetailsByHandle(ctx, linkedHandle)
		if err == nil {
			profiles = EnrichProfileWithSaasLinkedSocials(profiles, linkedDetails)
		}
	}
	return profiles, nil
}

// will be called from Campaign-Profiles Manager
func (m *SearchManager) FindProfileByPlatformHandle(ctx context.Context, platform string, handle string, linkedSource string, campaignProfileJoin bool) (*coredomain.Profile, error) {
	var profiles []coredomain.Profile
	var err error
	if platform == string(constants.InstagramPlatform) {
		profiles, err = m.FindByInstagramHandle(ctx, platform, handle, linkedSource, campaignProfileJoin)
		if err != nil {
			return nil, err
		}
	} else if platform == string(constants.YoutubePlatform) {
		profiles, err = m.FindByYoutubeHandle(ctx, platform, handle, linkedSource, campaignProfileJoin)
		if err != nil {
			return nil, err
		}
	}
	return &profiles[0], nil
}

func (m *SearchManager) MakeYoutubeEntityFromBeat(apiData *beatservice.ProfileResponse) dao.YoutubeAccountEntity {
	var youtubeAccount dao.YoutubeAccountEntity
	youtubeAccount.ChannelId = &apiData.Profile.YoutubeChannelId
	youtubeAccount.Title = &apiData.Profile.YoutubeTitle
	youtubeAccount.Followers = &apiData.Profile.YoutubeSubscribers
	youtubeAccount.UploadsCount = &apiData.Profile.YoutubeUploads
	youtubeAccount.ViewsCount = &apiData.Profile.YoutubeViews
	youtubeAccount.Thumbnail = &apiData.Profile.YoutubeThumbnail
	youtubeAccount.EnabledForSaas = EnabledForSaasYoutube(&youtubeAccount)
	youtubeAccount.Enabled = true
	flag := false
	youtubeAccount.Deleted = &flag
	name := ""
	if youtubeAccount.Title != nil {
		name = strings.ToLower(*youtubeAccount.Title)
	}
	totalPhrase := name
	keywords := strings.Split(totalPhrase, " ")

	youtubeAccount.Keywords = pq.StringArray(keywords)
	youtubeAccount.SearchPhrase = &totalPhrase
	youtubeAccount.KeywordsAdmin = pq.StringArray(keywords)
	youtubeAccount.SearchPhraseAdmin = &totalPhrase

	if youtubeAccount.Followers != nil {
		followers := *youtubeAccount.Followers
		var label string
		if followers < 1000 {
			label = "Nano"
		} else if followers < 75000 {
			label = "Micro"
		} else if followers < 1000000 {
			label = "Macro"
		} else {
			label = "Mega"
		}
		youtubeAccount.Label = &label
	}
	return youtubeAccount
}

func (m *SearchManager) MakeInstagramEntityFromBeat(apiData *beatservice.ProfileResponse) dao.InstagramAccountEntity {
	var instagramAccount dao.InstagramAccountEntity
	instagramAccount.Handle = &apiData.Profile.Handle
	instagramAccount.BusinessID = &apiData.Profile.FbId
	instagramAccount.Thumbnail = &apiData.Profile.ProfilePicUrl
	instagramAccount.IgID = &apiData.Profile.ProfileId
	instagramAccount.Bio = &apiData.Profile.Biography
	instagramAccount.Following = &apiData.Profile.Following
	instagramAccount.Followers = &apiData.Profile.Followers
	instagramAccount.Name = &apiData.Profile.FullName
	instagramAccount.IsPrivate = apiData.Profile.IsPrivate
	instagramAccount.EnabledForSaas = EnabledForSaasInstagram(&instagramAccount)
	instagramAccount.Enabled = true
	flag := false
	instagramAccount.Deleted = &flag
	name := ""
	handle := ""
	if instagramAccount.Name != nil {
		name = strings.ToLower(*instagramAccount.Name)
	}
	if instagramAccount.Handle != nil {
		handle = strings.ToLower(*instagramAccount.Handle)
	}
	totalPhrase := name + " " + handle
	keywords := strings.Split(totalPhrase, " ")

	instagramAccount.Keywords = pq.StringArray(keywords)
	instagramAccount.SearchPhrase = &totalPhrase
	instagramAccount.KeywordsAdmin = pq.StringArray(keywords)
	instagramAccount.SearchPhraseAdmin = &totalPhrase

	if instagramAccount.Followers != nil {
		followers := *instagramAccount.Followers
		var label string
		if followers < 1000 {
			label = "Nano"
		} else if followers < 75000 {
			label = "Micro"
		} else if followers < 1000000 {
			label = "Macro"
		} else {
			label = "Mega"
		}
		instagramAccount.Label = &label
	}

	instagramAccount.PostCount = &apiData.Profile.Uploads
	return instagramAccount
}

func (m *SearchManager) FindYoutubeAccountsByProfileId(ctx context.Context, gccProfileId string, linkedSource string) (*[]coredomain.Profile, error) {
	var profiles []coredomain.Profile
	ytEtities, err := m.ytDao.FindAllYtAccountsByGccProfileId(ctx, gccProfileId)
	if err != nil {
		return nil, err
	}
	for _, entity := range *ytEtities {
		profile := GetProfileFromYoutubeEntity(ctx, entity, linkedSource)
		profiles = append(profiles, *profile)
	}
	return &profiles, nil
}

func (m *SearchManager) UpdateMultipleFieldsYA(ctx context.Context, newGccProfileId string, iGHandle *string, youtubeId string, keywordsAdmin *[]string, searchPhraseAdmin *string) (*dao.YoutubeAccountEntity, error) {
	id, _ := strconv.ParseInt(youtubeId, 10, 64)
	flag := false
	entity := dao.YoutubeAccountEntity{
		GccProfileId:             &newGccProfileId,
		GccLinkedInstagramHandle: iGHandle,
		Deleted:                  &flag,
		Enabled:                  true,
	}
	if keywordsAdmin != nil {
		entity.KeywordsAdmin = *keywordsAdmin
	}
	if searchPhraseAdmin != nil {
		entity.SearchPhraseAdmin = searchPhraseAdmin
	}
	updatedEntity, err := m.ytDao.Update(ctx, id, &entity)
	if err != nil {
		return nil, err
	}
	return updatedEntity, nil
}

func (m *SearchManager) UpdateMultipleFieldsIA(ctx context.Context, ytChannelId *string, instaId string, keywordsAdmin *[]string, searchPhraseAdmin *string) (*dao.InstagramAccountEntity, error) {
	id, _ := strconv.ParseInt(instaId, 10, 64)
	flag := false
	entity := dao.InstagramAccountEntity{
		GccLinkedChannelId: ytChannelId,
		Deleted:            &flag,
		Enabled:            true,
	}
	if keywordsAdmin != nil {
		entity.KeywordsAdmin = *keywordsAdmin
	}
	if searchPhraseAdmin != nil {
		entity.SearchPhraseAdmin = searchPhraseAdmin
	}
	updatedEntity, err := m.instaDao.Update(ctx, id, &entity)
	if err != nil {
		return nil, err
	}
	return updatedEntity, nil
}

func (m *SearchManager) FindContentByPlatformProfileId(ctx context.Context, platform string, platformProfileId int64, query coredomain.SearchQuery, sortBy string, sortDir string, page int, size int, linkedSource string) (*domain.Content, int64, error) {
	// var contentData []*Content
	var filteredCount int64
	var err error
	var content *domain.Content
	if platform == string(constants.InstagramPlatform) {
		profiles, err := m.FindByPlatformProfileId(ctx, platform, platformProfileId, linkedSource, false)
		if err != nil {
			return content, filteredCount, err
		}
		igId := profiles.IgId
		platformCode := profiles.PlatformCode
		handle := profiles.Handle
		beatClient := beatservice.New(ctx)
		apiData, err := beatClient.RecentPostsByPlatformAndPlatformId(platform, *igId)
		if err != nil {
			return content, filteredCount, err
		}
		if apiData.Status.Type == "SUCCESS" {
			posts := TransformDataToSocialProfilePosts(platformProfileId, apiData.Data, platform, *handle)
			content = &domain.Content{
				PlatformCode: platformCode,
				Platform:     platform,
				Posts:        posts,
			}
			filteredCount = int64(len(posts))
		}
	}
	return content, filteredCount, err
}

func (m *SearchManager) FindSimilarAccountsByPlatformProfileId(ctx context.Context, platform string, platformProfileId int64, query coredomain.SearchQuery, sortBy string, sortDir string, page int, size int, linkedSource string) (*[]coredomain.Profile, int64, error) {
	var profiles []coredomain.Profile
	var followers *int64
	var categories []string
	if platform == string(constants.InstagramPlatform) {
		instaAccount, err := m.instaDao.FindById(ctx, platformProfileId, false)
		if err != nil {
			return &profiles, 0, err
		}
		followers = instaAccount.Followers
		categories = instaAccount.Categories
	} else if platform == string(constants.YoutubePlatform) {
		ytAccount, err := m.ytDao.FindById(ctx, platformProfileId, false)
		if err != nil {
			return &profiles, 0, err
		}
		followers = ytAccount.Followers
		categories = ytAccount.Categories
	}
	categoriesStr := strings.Join(categories, ", ")
	query = coredomain.SearchQuery{}
	if followers != nil {
		minFollowers := *followers - int64(0.60*float64(*followers))
		maxFollowers := *followers + int64(0.60*float64(*followers))

		query.Filters = append(query.Filters, coredomain.SearchFilter{
			FilterType: "GTE",
			Field:      "followers",
			Value:      strconv.Itoa(int(minFollowers)),
		})
		query.Filters = append(query.Filters, coredomain.SearchFilter{
			FilterType: "LTE",
			Field:      "followers",
			Value:      strconv.Itoa(int(maxFollowers)),
		})
		query.Filters = append(query.Filters, coredomain.SearchFilter{
			FilterType: "ARRAY",
			Field:      "categories",
			Value:      categoriesStr,
		})
		query.Filters = append(query.Filters, coredomain.SearchFilter{
			FilterType: "NE",
			Field:      "id",
			Value:      strconv.FormatInt(platformProfileId, 10),
		})
	}
	computeCounts := false
	return m.SearchByPlatform(ctx, platform, query, sortBy, sortDir, page, size, linkedSource, false, computeCounts)
}

// SearchByPlatform
/**
Search API to power Creator Discovery
	- Input
		- Filters (list of AND based criteria)
		- Sort
		- Size
		- Page
	- Returns
		- Profiles matching criteria
		- Filtered Count
		- Total Count
*/
func getInstaYtSortingOrder(sortBy string) []string {
	sortByArray := strings.Split(sortBy, ",")
	for i, component := range sortByArray {
		if component == "engagement_rate_grade_order" {
			replacement := "array_position(ARRAY['Poor', 'Average', 'Good', 'Very Good', 'Excellent'], engagement_rate_grade)"
			sortByArray[i] = replacement
		} else if component == "avg_likes_grade_order" {
			replacement := "array_position(ARRAY['Poor', 'Average', 'Good', 'Very Good', 'Excellent'], avg_likes_grade)"
			sortByArray[i] = replacement
		} else if component == "views_30d_grade_order" {
			replacement := "array_position(ARRAY['Poor', 'Average', 'Good', 'Very Good', 'Excellent'], views_30d_grade)"
			sortByArray[i] = replacement
		}
	}
	return sortByArray
}

func (m *SearchManager) SearchByPlatform(ctx context.Context, platform string, query coredomain.SearchQuery, sortBy string, sortDir string, page int, size int, linkedSource string, campaignProfileJoin bool, computeCounts bool) (*[]coredomain.Profile, int64, error) {

	var profiles *[]coredomain.Profile
	var filteredCount int64
	var err error
	if platform == string(constants.InstagramPlatform) {
		profiles, filteredCount, err = m.searchInstagram(ctx, query, sortBy, sortDir, page, size, linkedSource, campaignProfileJoin, computeCounts)
	} else if platform == string(constants.YoutubePlatform) {
		profiles, filteredCount, err = m.searchYoutube(ctx, query, sortBy, sortDir, page, size, linkedSource, campaignProfileJoin, computeCounts)
	} else {
		profiles, filteredCount, err = m.searchCrossPlatform(ctx, query, sortBy, sortDir, page, size, linkedSource)
	}
	return profiles, filteredCount, err
}

func makeInstagramFilters(ctx context.Context, query coredomain.SearchQuery, campaignProfileJoin bool) (coredomain.SearchQuery, []coredomain.JoinClauses) {

	var newSearchQuery coredomain.SearchQuery
	var filterMap map[string]string
	var profileCollectionItemJoin, profileCollectionJoin, CampaignIdNotInClause bool
	appCtx := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)
	for _, filter := range query.Filters {
		if filter.Field == "campaignId" && filter.FilterType == "NE" {
			CampaignIdNotInClause = true
		}
	}
	for _, filter := range query.Filters {
		var newFilter coredomain.SearchFilter

		field := filter.Field
		// This is for Audience Location
		if strings.Contains(field, "city") || strings.Contains(field, "state") || strings.Contains(field, "country") {
			newFilter.Field = "instagram_account." + field
			newFilter.FilterType = filter.FilterType
			newFilter.Value = filter.Value
			newSearchQuery.Filters = append(newSearchQuery.Filters, newFilter)
		}
		if filter.Field == "engagement_rate" || filter.Field == "averageEngagement" {
			filterValfloat, _ := strconv.ParseFloat(filter.Value, 64)
			val := filterValfloat / 100
			filter.Value = strconv.FormatFloat(val, 'f', -1, 64)
		}
		if filter.Field == "gender" {
			filter.Value = strings.ToUpper(filter.Value)
		}

		if appCtx.PartnerId != nil && *appCtx.PartnerId == -1 {
			filterMap = domain.InstagramSearchFilterAdminMap
		} else {
			filterMap = domain.InstagramSearchFilterMap
		}

		value, ok := filterMap[field]
		if ok {
			newFilter.Field = value
			newFilter.FilterType = filter.FilterType
			newFilter.Value = filter.Value
			newSearchQuery.Filters = append(newSearchQuery.Filters, newFilter)

			if strings.Contains(value, "campaign_profiles.") {
				campaignProfileJoin = true
			}

			if strings.Contains(value, "profile_collection_item.") {
				profileCollectionItemJoin = true
			}

			if strings.Contains(value, "profile_collection.") {
				profileCollectionItemJoin = true
				profileCollectionJoin = true
			}
		}
	}
	var joinClauseStruct []coredomain.JoinClauses
	if campaignProfileJoin {
		joinClauseStruct = append(joinClauseStruct, coredomain.JoinClauses{
			PreloadVariable:  "CampaignProfile",
			PreloadCondition: "platform = ?",
			PreloadValue:     string(constants.InstagramPlatform),
			JoinClause:       "LEFT JOIN campaign_profiles ON instagram_account.id = campaign_profiles.platform_account_id AND campaign_profiles.platform = 'INSTAGRAM'",
		})
	}
	if profileCollectionItemJoin && !CampaignIdNotInClause {
		joinClauseStruct = append(joinClauseStruct, coredomain.JoinClauses{
			JoinClause: "LEFT JOIN profile_collection_item ON instagram_account.id = profile_collection_item.platform_account_code AND profile_collection_item.platform = 'INSTAGRAM'",
		})
	}
	if profileCollectionJoin && !CampaignIdNotInClause {
		joinClauseStruct = append(joinClauseStruct, coredomain.JoinClauses{
			JoinClause: "LEFT JOIN profile_collection ON profile_collection_item.profile_collection_id = profile_collection.id AND profile_collection_item.platform = 'INSTAGRAM'",
		})
	}

	return newSearchQuery, joinClauseStruct
}

func (m *SearchManager) searchInstagram(ctx context.Context, query coredomain.SearchQuery, sortBy string, sortDir string, page int, size int, linkedSource string, campaignProfileJoin bool, computeCounts bool) (*[]coredomain.Profile, int64, error) {

	newSearchQuery, joinClauseStruct := makeInstagramFilters(ctx, query, campaignProfileJoin)
	sortMap := domain.InstagramSortingMap
	_, ok := sortMap[sortBy]
	if !ok {
		sortBy = "followers"
	}
	sortByArray := getInstaYtSortingOrder(sortBy)

	entities, filteredCount, err := m.instaDao.SearchInstagramProfiles(ctx, newSearchQuery, sortByArray, sortDir, page, size, joinClauseStruct, computeCounts)
	var profiles []coredomain.Profile
	if err != nil {
		return &profiles, filteredCount, err
	}
	var linkedChannels []*string
	for _, entity := range entities {
		profile := GetProfileFromInstagramEntity(ctx, entity, linkedSource)
		profiles = append(profiles, *profile)

		if linkedSource == string(constants.GCCPlatform) && entity.GccLinkedChannelId != nil && *entity.GccLinkedChannelId != "" {
			linkedChannels = append(linkedChannels, entity.GccLinkedChannelId)
		} else if entity.LinkedChannelId != nil && *entity.LinkedChannelId != "" {
			linkedChannels = append(linkedChannels, entity.LinkedChannelId)
		}
	}
	if len(linkedChannels) > 0 {
		linkedDetails, err := m.GetYoutubeDetalsByChannelIds(ctx, linkedChannels)
		if err == nil {

			if linkedSource == string(constants.GCCPlatform) {
				profiles = EnrichProfileWithGccLinkedSocials(profiles, linkedDetails)
			} else {
				profiles = EnrichProfileWithSaasLinkedSocials(profiles, linkedDetails)
			}
		}
	}
	return &profiles, filteredCount, err
}
func makeYoutubeFilters(ctx context.Context, query coredomain.SearchQuery, linkedSource string, campaignProfileJoin bool) (coredomain.SearchQuery, []coredomain.JoinClauses) {
	var newSearchQuery coredomain.SearchQuery
	var filterMap map[string]string
	var profileCollectionItemJoin, profileCollectionJoin, CampaignIdNotInClause bool
	appCtx := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)
	for _, filter := range query.Filters {
		if filter.Field == "campaignId" && filter.FilterType == "NE" {
			CampaignIdNotInClause = true
		}
	}
	for _, filter := range query.Filters {
		var newFilter coredomain.SearchFilter

		field := filter.Field
		// This is for Audience Location
		if strings.Contains(field, "city") || strings.Contains(field, "state") || strings.Contains(field, "country") {
			newFilter.Field = "youtube_account." + field
			newFilter.FilterType = filter.FilterType
			newFilter.Value = filter.Value
			newSearchQuery.Filters = append(newSearchQuery.Filters, newFilter)
		}
		if filter.Field == "engagement_rate" {
			filterValfloat, _ := strconv.ParseFloat(filter.Value, 64)
			val := filterValfloat / 100
			filter.Value = strconv.FormatFloat(val, 'f', -1, 64)
		}
		if filter.Field == "gender" {
			filter.Value = strings.ToUpper(filter.Value)
		}

		if appCtx.PartnerId != nil && *appCtx.PartnerId == -1 {
			filterMap = domain.YoutubeSearchFilterAdminMap
		} else {
			filterMap = domain.YoutubeSearchFilterMap
		}

		value, ok := filterMap[field]
		if ok {
			newFilter.Field = value
			newFilter.FilterType = filter.FilterType
			newFilter.Value = filter.Value
			newSearchQuery.Filters = append(newSearchQuery.Filters, newFilter)

			if strings.Contains(value, "campaign_profiles.") {
				campaignProfileJoin = true
			}

			if strings.Contains(value, "profile_collection_item.") {
				profileCollectionItemJoin = true
			}
			if strings.Contains(value, "profile_collection.") {
				profileCollectionItemJoin = true
				profileCollectionJoin = true
			}
		}
	}
	var joinClauseStruct []coredomain.JoinClauses
	if campaignProfileJoin {
		joinClauseStruct = append(joinClauseStruct, coredomain.JoinClauses{
			PreloadVariable:  "CampaignProfile",
			PreloadCondition: "platform = ?",
			PreloadValue:     string(constants.YoutubePlatform),
			JoinClause:       "LEFT JOIN campaign_profiles ON youtube_account.id = campaign_profiles.platform_account_id AND campaign_profiles.platform = 'YOUTUBE'",
		})
	}
	if profileCollectionItemJoin && !CampaignIdNotInClause {
		joinClauseStruct = append(joinClauseStruct, coredomain.JoinClauses{
			JoinClause: "LEFT JOIN profile_collection_item ON youtube_account.id = profile_collection_item.platform_account_code::int8 AND profile_collection_item.platform = 'YOUTUBE'",
		})
	}
	if profileCollectionJoin && !CampaignIdNotInClause {
		joinClauseStruct = append(joinClauseStruct, coredomain.JoinClauses{
			JoinClause: "LEFT JOIN profile_collection ON profile_collection_item.profile_collection_id = profile_collection.id AND profile_collection_item.platform = 'YOUTUBE'",
		})
	}
	return newSearchQuery, joinClauseStruct
}

func (m *SearchManager) searchYoutube(ctx context.Context, query coredomain.SearchQuery, sortBy string, sortDir string, page int, size int, linkedSource string, campaignProfileJoin bool, computeCounts bool) (*[]coredomain.Profile, int64, error) {

	newSearchQuery, joinClauseStruct := makeYoutubeFilters(ctx, query, linkedSource, campaignProfileJoin)
	sortMap := domain.YoutubeSortingMap
	_, ok := sortMap[sortBy]
	if !ok {
		sortBy = "followers"
	}
	sortByArray := getInstaYtSortingOrder(sortBy)

	entities, filteredCount, err := m.ytDao.SearchYoutubeProfiles(ctx, newSearchQuery, sortByArray, sortDir, page, size, joinClauseStruct, computeCounts)
	var profiles []coredomain.Profile
	if err != nil {
		return &profiles, filteredCount, err
	}
	var linkedHandle []*string
	for _, entity := range entities {
		profile := GetProfileFromYoutubeEntity(ctx, entity, linkedSource)
		profiles = append(profiles, *profile)

		if linkedSource == string(constants.GCCPlatform) && entity.GccLinkedInstagramHandle != nil && *entity.GccLinkedInstagramHandle != "" {
			linkedHandle = append(linkedHandle, entity.GccLinkedInstagramHandle)
		} else if entity.LinkedHandle != nil && *entity.LinkedHandle != "" {
			linkedHandle = append(linkedHandle, entity.LinkedHandle)
		}
	}
	if len(linkedHandle) > 0 {
		linkedDetails, err := m.GetInstagramDetailsByHandle(ctx, linkedHandle)
		if err == nil {

			if linkedSource == string(constants.GCCPlatform) {
				profiles = EnrichProfileWithGccLinkedSocials(profiles, linkedDetails)
			} else {
				profiles = EnrichProfileWithSaasLinkedSocials(profiles, linkedDetails)
			}
		}
	}
	return &profiles, filteredCount, err
}

func (m *SearchManager) searchCrossPlatform(ctx context.Context, query coredomain.SearchQuery, sortBy string, sortDir string, page int, size int, linkedSource string) (*[]coredomain.Profile, int64, error) {
	ytsize := size / 2
	instaSize := size - ytsize

	var ytErr, instaErr error
	var profiles []coredomain.Profile
	var ytprofiles, instaprofiles *[]coredomain.Profile
	// var wg sync.WaitGroup

	if ytsize > 0 {
		// wg.Add(1)
		// go func() {
		// 	defer wg.Done()
		ytprofiles, _, ytErr = m.searchYoutube(ctx, query, sortBy, sortDir, page, ytsize, linkedSource, false, false)
		// }()

	}
	if instaSize > 0 {
		// wg.Add(1)
		// go func() {
		// 	defer wg.Done()
		instaprofiles, _, ytErr = m.searchInstagram(ctx, query, sortBy, sortDir, page, instaSize, linkedSource, false, false)
		// }()
	}
	// wg.Wait()
	if instaErr != nil && ytErr != nil {
		instaErr = errors.Wrap(instaErr, ytErr.Error())
		return &profiles, 0, instaErr
	}
	if ytprofiles != nil {
		profiles = append(profiles, *ytprofiles...)
	}
	if instaprofiles != nil {
		profiles = append(profiles, *instaprofiles...)
	}
	return &profiles, int64(len(profiles)), instaErr
}

func (m *SearchManager) FindByPlatformProfileIds(ctx context.Context, platform string, platformProfileId []int64, linkedSource string) (*[]coredomain.Profile, error) {
	var profiles []coredomain.Profile
	if platform == string(constants.InstagramPlatform) {
		instaEntities, err := m.instaDao.FindByPlatformProfileIds(ctx, platformProfileId)
		if err != nil {
			return nil, err
		}
		var linkedChannels []*string
		for _, entity := range *instaEntities {
			profile := GetProfileFromInstagramEntity(ctx, entity, linkedSource)
			profiles = append(profiles, *profile)
			if entity.LinkedChannelId != nil && *entity.LinkedChannelId != "" {
				linkedChannels = append(linkedChannels, entity.LinkedChannelId)
			}
		}
		if len(linkedChannels) > 0 {
			linkedDetails, err := m.GetYoutubeDetalsByChannelIds(ctx, linkedChannels)
			if err == nil {
				profiles = EnrichProfileWithSaasLinkedSocials(profiles, linkedDetails)
			}
		}
	} else if platform == string(constants.YoutubePlatform) {
		ytEtities, err := m.ytDao.FindByPlatformProfileIds(ctx, platformProfileId)
		if err != nil {
			return nil, err
		}
		var linkedHandle []*string
		for _, entity := range *ytEtities {
			profile := GetProfileFromYoutubeEntity(ctx, entity, linkedSource)
			profiles = append(profiles, *profile)
			if entity.LinkedHandle != nil && *entity.LinkedHandle != "" {
				linkedHandle = append(linkedHandle, entity.LinkedHandle)
			}
		}
		if len(linkedHandle) > 0 {
			linkedDetails, err := m.GetInstagramDetailsByHandle(ctx, linkedHandle)
			if err == nil {
				profiles = EnrichProfileWithSaasLinkedSocials(profiles, linkedDetails)
			}
		}
	}
	return &profiles, nil
}

func (m *SearchManager) GetYoutubeDetalsByChannelIds(ctx context.Context, channelIdArray []*string) (map[string]coredomain.SocialAccount, error) {
	entities, err := m.ytDao.GetYoutubeDetalsByChannelIds(ctx, channelIdArray)
	if err != nil {
		return nil, err
	}
	var linkedMap = make(map[string]coredomain.SocialAccount)
	for _, entity := range *entities {
		linked := MakeYoutubeLinkedSocials(entity)
		linkedMap[*entity.ChannelId] = linked
	}
	return linkedMap, nil
}

func (m *SearchManager) GetInstagramDetailsByHandle(ctx context.Context, handleArray []*string) (map[string]coredomain.SocialAccount, error) {
	entities, err := m.instaDao.GetInstagramDetalsByHandle(ctx, handleArray)
	if err != nil {
		return nil, err
	}
	var linkedMap = make(map[string]coredomain.SocialAccount)
	if entities != nil {
		for _, entity := range *entities {
			linked := MakeInstagramLinkedSocials(entity)
			linkedMap[*entity.Handle] = linked
		}
	}

	return linkedMap, nil
}

func (m *SearchManager) GetYoutubeDetailsByChannelIds(ctx context.Context, ChannelIds []string) (map[string]coredomain.SocialAccount, error) {
	entities, err := m.ytDao.FindByChanneldIds(ctx, ChannelIds)
	if err != nil {
		return nil, err
	}
	var linkedMap = make(map[string]coredomain.SocialAccount)
	if entities != nil {
		for _, entity := range *entities {
			linked := MakeYoutubeLinkedSocials(entity)
			linkedMap[*entity.ChannelId] = linked
		}
	}

	return linkedMap, nil
}
func (m *SearchManager) GetInstagramDetailsByIgIds(ctx context.Context, igIds []string) (map[string]coredomain.SocialAccount, error) {
	entities, err := m.instaDao.FindByIgIds(ctx, igIds)
	if err != nil {
		return nil, err
	}
	var linkedMap = make(map[string]coredomain.SocialAccount)
	if entities != nil {
		for _, entity := range *entities {
			linked := MakeInstagramLinkedSocials(entity)
			linkedMap[*entity.Handle] = linked
		}
	}

	return linkedMap, nil
}

func (m *SearchManager) FindNotableFollowersDetails(ctx context.Context, idArray []string, platform string) *[]coredomain.SocialAccount {
	var notableFollowers []coredomain.SocialAccount
	if platform == string(constants.InstagramPlatform) && len(idArray) > 0 {
		instaEntities, err := m.instaDao.FindByIgIds(ctx, idArray)
		if err == nil {
			for _, entity := range *instaEntities {
				followerProfile := MakeInstagramLinkedSocials(entity)
				notableFollowers = append(notableFollowers, followerProfile)
			}
		}
	} else if platform == string(constants.YoutubePlatform) && len(idArray) > 0 {
		ytEntities, err := m.ytDao.FindByChanneldIds(ctx, idArray)
		if err == nil {
			for _, entity := range *ytEntities {
				followerProfile := MakeYoutubeLinkedSocials(entity)
				notableFollowers = append(notableFollowers, followerProfile)
			}
		}
	}
	return &notableFollowers
}

func (m *SearchManager) EnrichProfilesWithProfileCollectionItemId(profiles []coredomain.Profile, itemsMap map[string]int64) *[]coredomain.Profile {
	for i := 0; i < len(profiles); i++ {
		value, ok := itemsMap[profiles[i].PlatformCode]
		if ok {
			profiles[i].ProfileCollectionItemId = value
		}
	}
	return &profiles
}

func EnabledForSaasInstagram(instagramAccount *dao.InstagramAccountEntity) bool {

	flag := true
	if (instagramAccount.Country == nil || *instagramAccount.Country != "IN") || instagramAccount.Followers == nil ||
		instagramAccount.EngagementRate == 0 || instagramAccount.AvgReach == nil || instagramAccount.AvgLikes == nil || instagramAccount.AvgComments == nil {
		flag = false
	}
	return flag
}

func EnabledForSaasYoutube(youtubeAccount *dao.YoutubeAccountEntity) bool {

	flag := true
	if (youtubeAccount.Country == nil || *youtubeAccount.Country != "IN") || youtubeAccount.Followers == nil ||
		youtubeAccount.ViewsCount == nil || youtubeAccount.AvgViews == nil {
		flag = false
	}
	return flag
}

func (m *SearchManager) UpdateInstagramForFullRefresh(ctx context.Context, instaId string, entity dao.InstagramAccountEntity) (*dao.InstagramAccountEntity, error) {
	id, _ := strconv.ParseInt(instaId, 10, 64)
	flag := false
	entity.Deleted = &flag
	entity.Enabled = true
	updatedEntity, err := m.instaDao.Update(ctx, id, &entity)
	if err != nil {
		return nil, err
	}
	return updatedEntity, nil
}

func (m *SearchManager) UpdateYoutubeForFullRefresh(ctx context.Context, instaId string, entity dao.YoutubeAccountEntity) (*dao.YoutubeAccountEntity, error) {
	id, _ := strconv.ParseInt(instaId, 10, 64)
	flag := false
	entity.Deleted = &flag
	entity.Enabled = true
	updatedEntity, err := m.ytDao.Update(ctx, id, &entity)
	if err != nil {
		return nil, err
	}
	return updatedEntity, nil
}

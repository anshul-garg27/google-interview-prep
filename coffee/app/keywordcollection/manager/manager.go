package manager

import (
	discoverydomain "coffee/app/discovery/domain"
	discoverymanager "coffee/app/discovery/manager"
	"coffee/app/keywordcollection/dao"
	"coffee/app/keywordcollection/domain"
	"coffee/constants"
	"coffee/core/appcontext"
	coredomain "coffee/core/domain"
	"coffee/core/rest"
	"coffee/helpers"
	"coffee/publishers"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/viper"
)

type KeywordCollectionManager struct {
	*rest.Manager[domain.KeywordCollectionEntry, dao.KeywordCollectionEntity, string]
	dao              *dao.KeywordCollectionDAO
	metricsDao       *dao.MetricsDao
	discoveryManager *discoverymanager.SearchManager
}

func (m *KeywordCollectionManager) init(ctx context.Context) {
	newDAO := dao.CreateKeywordCollectionDAO(ctx)
	m.Manager = rest.NewManager(ctx, newDAO.Dao, m.ToEntry, m.ToEntity)
	m.dao = newDAO
	m.discoveryManager = discoverymanager.CreateSearchManager(ctx)
	m.metricsDao = dao.CreateKeywordCollectionMetricsDao(ctx)
}

func CreateKeywordCollectionManager(ctx context.Context) *KeywordCollectionManager {
	manager := &KeywordCollectionManager{}
	manager.init(ctx)
	return manager
}

func (m *KeywordCollectionManager) ToEntry(entity *dao.KeywordCollectionEntity) (*domain.KeywordCollectionEntry, error) {
	entry := &domain.KeywordCollectionEntry{
		Id:                   &entity.Id,
		ShareId:              &entity.ShareId,
		Name:                 &entity.Name,
		Platform:             &entity.Platform,
		PartnerId:            &entity.PartnerId,
		Enabled:              entity.Enabled,
		ReportBucket:         entity.ReportBucket,
		ReportPath:           entity.ReportPath,
		DemographyReportPath: entity.DemographyReportPath,
		CreatedBy:            &entity.CreatedBy,
		Keywords:             (*[]string)(&entity.Keywords),
		CreatedAt:            helpers.ToInt64(entity.CreatedAt.Unix()),
		UpdatedAt:            helpers.ToInt64(entity.UpdatedAt.Unix()),
		JobId:                entity.JobId,
	}
	if entity.StartTime != nil {
		ts := (*entity.StartTime).UnixMilli()
		entry.StartTime = &ts
	}
	if entity.EndTime != nil {
		ts := (*entity.EndTime).UnixMilli()
		entry.EndTime = &ts
	}
	return entry, nil
}

func (m *KeywordCollectionManager) ToEntity(entry *domain.KeywordCollectionEntry) (*dao.KeywordCollectionEntity, error) {
	entity := &dao.KeywordCollectionEntity{}
	if entry.Id != nil {
		entity.Id = *entry.Id
	}
	if entry.ShareId != nil {
		entity.ShareId = *entry.ShareId
	}
	if entry.PartnerId != nil {
		entity.PartnerId = *entry.PartnerId
	}
	if entry.Name != nil {
		entity.Name = *entry.Name
	}
	if entry.Platform != nil {
		entity.Platform = *entry.Platform
	}
	if entry.Keywords != nil {
		entity.Keywords = *entry.Keywords
	}
	if entry.ReportBucket != nil {
		entity.ReportBucket = entry.ReportBucket
	}
	if entry.ReportPath != nil {
		entity.ReportPath = entry.ReportPath
	}
	if entry.DemographyReportPath != nil {
		entity.DemographyReportPath = entry.DemographyReportPath
	}
	if entry.Enabled != nil {
		entity.Enabled = entry.Enabled
	}
	if entry.StartTime != nil {
		st := time.UnixMilli(*entry.StartTime)
		entity.StartTime = &st
	}
	if entry.EndTime != nil {
		et := time.UnixMilli(*entry.EndTime)
		entity.EndTime = &et
	}
	if entry.CreatedBy != nil {
		entity.CreatedBy = *entry.CreatedBy
	}
	if entry.JobId != nil {
		entity.JobId = entry.JobId
	}
	return entity, nil
}

func (m *KeywordCollectionManager) getKeywordCollectionName(list []string) string {
	if len(list) == 0 {
		return ""
	}

	sort.Sort(sort.Reverse(sort.StringSlice(list)))

	var result string
	if len(list) > 3 {
		result = strings.Join(list[:3], ", ") + "..."
	} else {
		result = strings.Join(list, ", ")
	}

	return result
}

func (m *KeywordCollectionManager) Create(ctx context.Context, input *domain.KeywordCollectionEntry) (*domain.KeywordCollectionEntry, error) {
	if input.Name == nil || *input.Name == "" {
		name := m.getKeywordCollectionName(*input.Keywords)
		input.Name = &name
	}
	if input.Name == nil {
		return nil, errors.New("bad Request - Missing mandatory fields")
	}

	uuid := uuid.New().String()
	input.Id = &uuid

	createdEntry, err := m.Manager.Create(ctx, input)
	if err != nil {
		return nil, err
	}
	if createdEntry == nil || createdEntry.Id == nil {
		return nil, errors.New("failed to create post collection")
	}
	if err != nil {
		panic(err)
	}
	err = m.createReportJob(createdEntry)
	if err != nil {
		panic(err)
	}
	return createdEntry, err
}

func (m *KeywordCollectionManager) createReportJob(collectionEntry *domain.KeywordCollectionEntry) error {
	messagePayload := m.createPayloadForBeat(collectionEntry)
	jsonBytes, _ := json.Marshal(messagePayload)
	topic := "beat.dx___keyword_collection_rk"
	err := publishers.PublishMessage(jsonBytes, topic)
	return err
}

func (m *KeywordCollectionManager) createPayloadForBeat(createdEntry *domain.KeywordCollectionEntry) map[string]interface{} {
	messagePayload := make(map[string]interface{})
	messagePayload["job_id"] = fmt.Sprintf("%s_KC_%s", viper.GetString("ENV"), *createdEntry.Id)
	messagePayload["keywords"] = createdEntry.Keywords
	messagePayload["platform"] = createdEntry.Platform
	startTimestamp := time.UnixMilli(*createdEntry.StartTime)
	endTimestamp := time.UnixMilli(*createdEntry.EndTime)
	messagePayload["start_date"] = fmt.Sprint(startTimestamp.Format("2006-01-02"))
	messagePayload["end_date"] = fmt.Sprint(endTimestamp.Format("2006-01-02"))
	return messagePayload
}

func (m *KeywordCollectionManager) Update(ctx context.Context, id string, input *domain.KeywordCollectionEntry) (*domain.KeywordCollectionEntry, error) {
	updatedEntry, err := m.Manager.Update(ctx, id, input)
	if err != nil {
		return nil, err
	}
	if updatedEntry == nil || updatedEntry.Id == nil {
		return nil, errors.New("failed to create profile collection")
	}
	entry, err := m.Manager.FindById(ctx, id)
	return entry, err
}

func (m *KeywordCollectionManager) FindById(ctx context.Context, collectionId string, fetchTotalItems bool) (*domain.KeywordCollectionEntry, error) {
	e, err := m.Manager.FindById(ctx, collectionId)
	if err != nil {
		return nil, err
	}
	appCtx := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)
	if e != nil && *appCtx.PartnerId != -1 && *e.PartnerId != *appCtx.PartnerId {
		return nil, errors.New("You are not authorized to view this collection")
	}
	return e, nil
}

func (m *KeywordCollectionManager) CreateReport(ctx context.Context, collectionId string) (*domain.KeywordCollectionEntry, error) {
	e, err := m.Manager.FindById(ctx, collectionId)
	if err != nil {
		return nil, err
	}
	appCtx := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)
	if e != nil && *appCtx.PartnerId != -1 && *e.PartnerId != *appCtx.PartnerId {
		return nil, errors.New("You are not authorized to view this collection")
	}
	err = m.createReportJob(e)
	if err != nil {
		panic(err)
	}
	return e, nil
}

func (m *KeywordCollectionManager) GetCollectionReport(ctx context.Context, collectionId string) (*domain.KeywordCollectionEntry, error) {
	collection, err := m.Manager.FindById(ctx, collectionId)
	if err != nil {
		return nil, err
	}
	if collection.ReportBucket == nil || collection.ReportPath == nil || collection.DemographyReportPath == nil {
		return nil, errors.New("report not generated yet")
	}
	collection.ReportStatus = "COMPLETE"
	summary, err := m.metricsDao.GetSummary(ctx, *collection.ReportBucket, *collection.ReportPath)
	if err != nil {
		return nil, err
	}
	if summary == nil {
		return nil, errors.New("unable to fetch summary")
	}
	posts, err := m.metricsDao.GetPosts(ctx, *collection.ReportBucket, *collection.ReportPath, "reach", "DESC")
	if err != nil {
		return nil, err
	}
	if posts == nil {
		return nil, errors.New("unable to fetch posts")
	}

	profiles, err := m.metricsDao.GetProfiles(ctx, *collection.ReportBucket, *collection.ReportPath, "reach", "DESC")
	if err != nil {
		return nil, err
	}
	if profiles == nil {
		return nil, errors.New("unable to fetch profiles")
	}

	rollup := "monthly"
	duration := time.Duration(*collection.EndTime-*collection.StartTime) * time.Millisecond
	if duration <= time.Hour*24*60 {
		rollup = "weekly"
	}
	monthlyStats, err := m.metricsDao.GetMonthlyStats(ctx, *collection.ReportBucket, *collection.ReportPath, rollup)
	if err != nil {
		return nil, err
	}
	if monthlyStats == nil {
		return nil, errors.New("unable to fetch monthly stats")
	}

	var keywords []dao.Keywords
	if *collection.Platform == string(constants.InstagramPlatform) {
		keywords, err = m.metricsDao.GetKeywordsFromTitle(ctx, *collection.ReportBucket, *collection.ReportPath)
	} else if *collection.Platform == string(constants.YoutubePlatform) {
		keywords, err = m.metricsDao.GetKeywords(ctx, *collection.ReportBucket, *collection.ReportPath)
	}
	if err != nil {
		return nil, err
	}

	categories, err := m.metricsDao.GetCategories(ctx, *collection.ReportBucket, *collection.ReportPath)
	if err != nil {
		return nil, err
	}

	languages, err := m.metricsDao.GetLanguages(ctx, *collection.ReportBucket, *collection.ReportPath)
	if err != nil {
		return nil, err
	}
	if languages == nil {
		return nil, errors.New("unable to fetch languages")
	}

	languagesMonthly, err := m.metricsDao.GetLanguagesMonthly(ctx, *collection.ReportBucket, *collection.ReportPath)
	if err != nil {
		return nil, err
	}
	if languagesMonthly == nil {
		return nil, errors.New("unable to fetch monthly languages")
	}

	demography, err := m.metricsDao.GetDemography(ctx, *collection.ReportBucket, *collection.DemographyReportPath)
	if err != nil {
		return nil, err
	}
	if demography == nil {
		return nil, errors.New("unable to fetch demography")
	}

	appCtx := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)
	if collection != nil && *appCtx.PartnerId != -1 && *collection.PartnerId != *appCtx.PartnerId {
		return nil, errors.New("You are not authorized to view this collection")
	}
	topKeywords := make(map[string]float64)
	for i := range keywords {
		topKeywords[keywords[i].Keyword] = keywords[i].Score
	}
	runningSum := 0.0
	topCategories := make(map[string]float64)
	for i := range categories {
		topCategories[categories[i].Category] = categories[i].PostsPerc
		runningSum += categories[i].PostsPerc
		if i == 4 {
			topCategories["Others"] = 100 - runningSum
		}
	}
	topLanguages := make(map[string]float64)
	runningSum = 0.0
	for i := range languages {
		if i > 5 {
			continue
		}
		topLanguages[transformLanguageCode(languages[i].Language)] = languages[i].PostsPerc
		runningSum += languages[i].PostsPerc
		if i == 5 {
			topLanguages["Others"] = 100 - runningSum
		}
	}
	topLanguagesMonthly := make(map[string]map[string]float64)
	currentMonth := ""
	j := 0
	runningSum = 0.0
	for i := range languagesMonthly {
		if languagesMonthly[i].Month != currentMonth {
			j = 0
			runningSum = 0.0
		}
		currentMonth = languagesMonthly[i].Month
		j += 1
		if j > 5 {
			continue
		}
		_, isPresent := topLanguagesMonthly[languagesMonthly[i].Month]
		if !isPresent {
			topLanguagesMonthly[languagesMonthly[i].Month] = make(map[string]float64)
		}
		topLanguagesMonthly[languagesMonthly[i].Month][transformLanguageCode(languagesMonthly[i].Language)] = languagesMonthly[i].PostsPerc
		runningSum += languagesMonthly[i].PostsPerc
		if j == 5 {
			topLanguagesMonthly[languagesMonthly[i].Month]["Others"] = 100 - runningSum
		}
	}
	monthlyViews := make(map[string]int64)
	monthlyUploads := make(map[string]int64)
	for i := range monthlyStats {
		monthlyViews[monthlyStats[i].Month] = monthlyStats[i].Views
		monthlyUploads[monthlyStats[i].Month] = monthlyStats[i].Uploads
	}
	var demoAge map[string]*float64
	var demoGender map[string]*float64
	var demoAgeGender map[string]map[string]*float64
	ageStr := strings.ReplaceAll(demography.AudienceAge, "'", "\"")
	genderStr := strings.ReplaceAll(demography.AudienceGender, "'", "\"")
	ageGenderStr := strings.ReplaceAll(demography.AudienceAgeGender, "'", "\"")
	json.Unmarshal([]byte(ageStr), &demoAge)
	json.Unmarshal([]byte(genderStr), &demoGender)
	json.Unmarshal([]byte(ageGenderStr), &demoAgeGender)
	male := make(map[string]*float64)
	male = demoAgeGender["male"]
	female := make(map[string]*float64)
	female = demoAgeGender["female"]
	gender := make(map[string]*float64)
	gender["male"] = demoGender["male_per"]
	gender["female"] = demoGender["female_per"]
	topPosts := m.transformPostEntities(posts, collection)
	topProfiles := m.transformProfileEntities(ctx, profiles, collection)
	collection.Report = &domain.KeywordCollectionReport{
		TotalPosts:          summary.Posts,
		TotalReach:          summary.Reach,
		TotalLikes:          summary.Likes,
		TotalComments:       summary.Comments,
		Engagement:          summary.Engagement,
		EMP:                 float64(summary.Emp),
		EngagementRate:      summary.EngagementRate,
		TotalViews:          summary.Reach,
		TopKeywords:         topKeywords,
		TopCategories:       topCategories,
		TopLanguages:        topLanguages,
		TopProfiles:         topProfiles,
		TopPosts:            topPosts,
		TopLanguagesMonthly: topLanguagesMonthly,
		MonthlyViews:        monthlyViews,
		MonthlyUploads:      monthlyUploads,
		AudienceAgeGender: discoverydomain.AudienceAgeGender{
			Male:   male,
			Female: female,
		},
		AudienceGender: gender,
	}
	return collection, nil
}

func (m *KeywordCollectionManager) GetCollectionReportPartial(ctx context.Context, collection *domain.KeywordCollectionEntry) (*domain.KeywordCollectionEntry, error) {
	if collection.ReportBucket == nil || collection.ReportPath == nil || collection.DemographyReportPath == nil {
		return nil, errors.New("report not generated yet")
	}
	collection.ReportStatus = "COMPLETE"
	summary, err := m.metricsDao.GetSummary(ctx, *collection.ReportBucket, *collection.ReportPath)
	if err != nil {
		return nil, err
	}
	if summary == nil {
		return nil, errors.New("unable to fetch summary")
	}
	posts, err := m.metricsDao.GetPosts(ctx, *collection.ReportBucket, *collection.ReportPath, "reach", "DESC")
	if err != nil {
		return nil, err
	}
	if posts == nil {
		return nil, errors.New("unable to fetch posts")
	}

	profiles, err := m.metricsDao.GetProfiles(ctx, *collection.ReportBucket, *collection.ReportPath, "reach", "DESC")
	if err != nil {
		return nil, err
	}
	if profiles == nil {
		return nil, errors.New("unable to fetch profiles")
	}

	appCtx := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)
	if collection != nil && *appCtx.PartnerId != -1 && *collection.PartnerId != *appCtx.PartnerId {
		return nil, errors.New("You are not authorized to view this collection")
	}
	topPosts := m.transformPostEntities(posts, collection)
	topProfiles := m.transformProfileEntities(ctx, profiles, collection)
	collection.Report = &domain.KeywordCollectionReport{
		TotalPosts:     summary.Posts,
		TotalReach:     summary.Reach,
		TotalLikes:     summary.Likes,
		TotalComments:  summary.Comments,
		Engagement:     summary.Engagement,
		EMP:            float64(summary.Emp),
		EngagementRate: float64(summary.EngagementRate),
		TotalViews:     summary.Reach,
		TopProfiles:    topProfiles,
		TopPosts:       topPosts,
	}
	return collection, nil
}

func transformLanguageCode(language string) string {
	if strings.Contains(language, "-") {
		language = strings.Split(language, "-")[0]
	}
	languageName, present := constants.LanguageNamesExhaustive[language]
	if present {
		return languageName
	}
	return language
}

func (m *KeywordCollectionManager) transformProfileEntities(ctx context.Context, profiles []dao.Profile, collection *domain.KeywordCollectionEntry) []domain.KeywordCollectionProfile {
	var topProfiles []domain.KeywordCollectionProfile
	var profileIds []string
	for i := range profiles {
		profileIds = append(profileIds, profiles[i].ProfileId)
	}
	socialAccounts := m.discoveryManager.FindNotableFollowersDetails(ctx, profileIds, *collection.Platform)
	profileMap := make(map[string]coredomain.SocialAccount)
	if socialAccounts != nil {
		for i := range *socialAccounts {
			account := (*socialAccounts)[i]
			if account.Handle == nil {
				continue
			}
			profileMap[*account.Handle] = account
		}
	}
	for i := range profiles {
		handle := ""
		username := ""
		if *collection.Platform == string(constants.InstagramPlatform) {
			handle = profiles[i].Handle
			username = handle
		} else if *collection.Platform == string(constants.YoutubePlatform) {
			handle = profiles[i].ProfileId
			username = profiles[i].Handle
		}
		profileCode := ""
		account, present := profileMap[handle]
		if present && account.ProfileCode != nil {
			profileCode = *account.ProfileCode
		}
		topProfiles = append(topProfiles, domain.KeywordCollectionProfile{
			AccountInfo: coredomain.SocialAccount{
				Name:                     &profiles[i].Name,
				Platform:                 *collection.Platform,
				Handle:                   &handle,
				ProfileCode:              &profileCode,
				Username:                 &username,
				PlatformThumbnail:        &profiles[i].Thumbnail,
				EngagementRatePercenatge: profiles[i].EngagementRate,
				Followers:                &profiles[i].Followers,
				Following:                Int64(0), // TODO
				ProfileLink:              String(getProfileLink(*collection.Platform, handle)),
				Category:                 &profiles[i].Category,
			},
			Views:      &profiles[i].Views,
			Plays:      &profiles[i].Views,
			Likes:      &profiles[i].Likes,
			Comments:   &profiles[i].Comments,
			Engagement: Int64(profiles[i].Comments + profiles[i].Likes),
			Uploads:    &profiles[i].Uploads,
		})
	}
	return topProfiles
}

func (m *KeywordCollectionManager) transformPostEntities(posts []dao.Post, collection *domain.KeywordCollectionEntry) []domain.KeywordCollectionPost {
	var topPosts []domain.KeywordCollectionPost
	for i := range posts {
		handle := ""
		username := ""
		if *collection.Platform == string(constants.InstagramPlatform) {
			handle = posts[i].Handle
			username = handle
		} else if *collection.Platform == string(constants.YoutubePlatform) {
			handle = posts[i].ProfileId
			username = posts[i].Handle
		}
		topPosts = append(topPosts, domain.KeywordCollectionPost{
			AccountInfo: coredomain.SocialAccount{
				Name:                     &posts[i].Name,
				Platform:                 *collection.Platform,
				Handle:                   &handle,
				Username:                 &username,
				PlatformThumbnail:        &posts[i].ProfileThumbnailUrl,
				EngagementRatePercenatge: 0, // TODO
				Followers:                &posts[i].Followers,
				Following:                nil, //TODO
				ProfileLink:              String(getProfileLink(*collection.Platform, handle)),
				Category:                 &posts[i].Category,
			},
			Link:           getPostLink(*collection.Platform, posts[i].Shortcode),
			Title:          posts[i].PostTitle,
			Platform:       *collection.Platform,
			Thumbnail:      posts[i].PostThumbnailUrl,
			Views:          &posts[i].Reach,
			PublishedAt:    Int64(posts[i].PostPublishTime.Unix()),
			Plays:          &posts[i].Plays,
			Likes:          &posts[i].Likes,
			PostType:       &posts[i].PostType,
			Comments:       &posts[i].Comments,
			Engagement:     Int64(posts[i].Comments + posts[i].Likes),
			EngagementRate: &posts[i].EngagementRate,
		})
	}
	return topPosts
}

func getProfileLink(platform string, handle string) string {
	if platform == string(constants.InstagramPlatform) {
		return "https://instagram.com/" + handle
	} else if platform == string(constants.YoutubePlatform) {
		return "https://youtube.com/channel/" + handle
	}
	return ""
}

func getPostLink(platform string, shortcode string) string {
	if platform == string(constants.InstagramPlatform) {
		return "https://instagram.com/p/" + shortcode
	} else if platform == string(constants.YoutubePlatform) {
		return "https://youtube.com/v/" + shortcode
	}
	return ""
}

func String(s string) *string {
	return &s
}

func Int64(i int64) *int64 {
	return &i
}

func Float64(f float64) *float64 {
	return &f
}

func (m *KeywordCollectionManager) GetCollectionPosts(ctx context.Context, searchQuery coredomain.SearchQuery, sortBy string, sortDir string, page int, size int) ([]domain.KeywordCollectionPost, int64, error) {
	collectionId := rest.GetFilterValueForKey(searchQuery.Filters, "collectionId")
	collection, err := m.Manager.FindById(ctx, *collectionId)
	if err != nil {
		return nil, 0, err
	}
	if collection.ReportBucket == nil || collection.ReportPath == nil || collection.DemographyReportPath == nil {
		return nil, 0, errors.New("report not generated yet")
	}
	appCtx := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)
	if collection != nil && *appCtx.PartnerId != -1 && *collection.PartnerId != *appCtx.PartnerId {
		return nil, 0, errors.New("You are not authorized to view this collection")
	}
	postsSortBy := "reach"
	if sortBy == "recency" {
		postsSortBy = "post_publish_time"
	}
	if sortBy == "likes" {
		postsSortBy = "likes"
	}
	if sortBy == "comments" {
		postsSortBy = "comments"
	}
	posts, err := m.metricsDao.GetPosts(ctx, *collection.ReportBucket, *collection.ReportPath, postsSortBy, "DESC")
	if err != nil {
		return nil, 0, err
	}
	if posts == nil {
		return nil, 0, errors.New("unable to fetch posts")
	}
	topPosts := m.transformPostEntities(posts, collection)
	return topPosts, int64(len(topPosts)), nil
}

func (m *KeywordCollectionManager) GetCollectionProfiles(ctx context.Context, searchQuery coredomain.SearchQuery, sortBy string, sortDir string, page int, size int) ([]domain.KeywordCollectionProfile, int64, error) {
	collectionId := rest.GetFilterValueForKey(searchQuery.Filters, "collectionId")
	collection, err := m.Manager.FindById(ctx, *collectionId)
	if err != nil {
		return nil, 0, err
	}
	if collection.ReportBucket == nil || collection.ReportPath == nil || collection.DemographyReportPath == nil {
		return nil, 0, errors.New("report not generated yet")
	}
	appCtx := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)
	if collection != nil && *appCtx.PartnerId != -1 && *collection.PartnerId != *appCtx.PartnerId {
		return nil, 0, errors.New("You are not authorized to view this collection")
	}
	profilesSortBy := "reach"
	if sortBy == "recency" {
		profilesSortBy = "last_posted_on"
	}
	if sortBy == "likes" {
		profilesSortBy = "likes"
	}
	if sortBy == "comments" {
		profilesSortBy = "comments"
	}
	if sortBy == "views" {
		profilesSortBy = "reach"
	}
	profiles, err := m.metricsDao.GetProfiles(ctx, *collection.ReportBucket, *collection.ReportPath, profilesSortBy, "DESC")
	if err != nil {
		return nil, 0, err
	}
	if profiles == nil {
		return nil, 0, errors.New("unable to fetch profiles")
	}
	topProfiles := m.transformProfileEntities(ctx, profiles, collection)
	return topProfiles, int64(len(topProfiles)), nil
}

func (m *KeywordCollectionManager) FindByIdInternal(ctx context.Context, collectionId string) (*domain.KeywordCollectionEntry, error) {
	e, err := m.Manager.FindById(ctx, collectionId)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (m *KeywordCollectionManager) FindByShareId(ctx context.Context, shareId string) (*domain.KeywordCollectionEntry, error) {
	e, err := m.dao.FindByShareId(ctx, shareId)
	if err != nil {
		return nil, err
	}
	return m.ToEntry(e)
}

func (m *KeywordCollectionManager) CountCollectionsForPartner(ctx context.Context, partnerId int64) (int64, error) {
	count, err := m.dao.CountCollectionsForPartner(ctx, partnerId)
	if err != nil {
		return 0, err
	}
	return count, nil
}

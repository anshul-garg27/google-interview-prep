package leaderboard

import (
	discoverymanager "coffee/app/discovery/manager"
	"coffee/constants"

	coredomain "coffee/core/domain"
	"coffee/helpers"
	"context"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

type Manager struct {
	leaderboardDao             *LeaderboardDao
	socialProfileTimeSeriesDao *SocialProfileTimeSeriesDao
	discoveryManager           *discoverymanager.SearchManager
}

func CreateManager(ctx context.Context) *Manager {

	leaderboardDao := createLeaderboardDao(ctx)
	socialProfileTimeSeriesDao := createSocialProfileTimeSeriesDao(ctx)
	discoveryManager := discoverymanager.CreateSearchManager(ctx)
	manager := &Manager{
		leaderboardDao:             leaderboardDao,
		socialProfileTimeSeriesDao: socialProfileTimeSeriesDao,
		discoveryManager:           discoveryManager,
	}
	return manager
}
func getDateRanges(query coredomain.SearchQuery) (time.Time, time.Time) {
	var startDate, oneMonthBack, endDate time.Time
	var err error
	for _, filter := range query.Filters {
		if filter.Field == "month" {
			value := filter.Value
			startDate, err = time.Parse("2006-01-02", value) // parse input date string
			if err != nil {
				return oneMonthBack, startDate
			}
			firstDayOfMonth := time.Date(startDate.Year(), startDate.Month(), 1, 0, 0, 0, 0, startDate.Location())
			endDate = firstDayOfMonth.AddDate(0, 1, -1)

		}
	}
	return startDate, endDate
}

func (m *Manager) LeaderBoardByPlatform(ctx context.Context, platform string, query coredomain.SearchQuery, sortBy string, sortDir string, page int, size int, source string) (*[]LeaderboardProfile, int64, error) {
	var leaderboard []LeaderboardProfile
	if platform == string(constants.InstagramPlatform) || platform == string(constants.YoutubePlatform) {
		return m.getInstagramYoutubeLeaderBoard(ctx, platform, query, sortBy, sortDir, page, size, source)
	} else if platform == "ALL" {
		return m.getCrossPlatformLeaderBoard(ctx, platform, query, sortBy, sortDir, page, size, source)
	}
	return &leaderboard, 0, nil

}

func (m *Manager) getInstagramYoutubeLeaderBoard(ctx context.Context, platform string, query coredomain.SearchQuery, sortBy string, sortDir string, page int, size int, source string) (*[]LeaderboardProfile, int64, error) {
	var leaderboard []LeaderboardProfile
	monthFlag, searchFlag := false, false
	for _, filter := range query.Filters {
		if filter.Field == "month" {
			monthFlag = true
		}
		if filter.Field == "search_phrase" {
			searchFlag = true
		}
	}
	if !monthFlag {
		lastMonthStart := time.Date(time.Now().AddDate(0, -1, 0).Year(), time.Now().AddDate(0, -1, 0).Month(), 1, 0, 0, 0, 0, time.Now().AddDate(0, -1, 0).Location()).Format("2006-01-02")
		query.Filters = append(query.Filters, coredomain.SearchFilter{FilterType: "EQ", Field: "month", Value: lastMonthStart})
	}
	query.Filters = append(query.Filters, coredomain.SearchFilter{FilterType: "EQ", Field: "platform", Value: platform})
	sortBy = getSortByKey(query, sortBy)
	if sortBy != "ASC" {
		sortDir = "ASC"
	}
	entities, filteredCount, err := m.leaderboardDao.Search(ctx, query, sortBy, sortDir, page, size)
	if err != nil {
		return &leaderboard, 0, err
	}
	// Last Month & Influencer Basic Data
	var platformIdArray []int64
	for _, entity := range entities {
		platformIdArray = append(platformIdArray, entity.PlatformProfileId)
	}
	platformProfileData, err := m.discoveryManager.FindByPlatformProfileIds(ctx, platform, platformIdArray, source)
	var platformProfileMap = make(map[int64]coredomain.Profile)
	if err == nil {
		for _, insta := range *platformProfileData {
			code, _ := strconv.ParseInt(insta.PlatformCode, 10, 64)
			platformProfileMap[code] = insta
		}
	}
	var timeSeriesData *ProfilesTimeSeries
	timeseriesStartDate, timeseriesEndDate := getDateRanges(query)
	timeSeriesData, err = m.socialProfileTimeSeriesDao.GetMultipleTimeSeriesData(ctx, timeseriesStartDate, timeseriesEndDate, platform, platformIdArray)
	if err != nil {
		log.Error(err)
	}
	for i := range entities {
		rank := i + (page-1)*size
		leaderboardProfile := GetLeaderboardProfileFromEntity(entities[i], platformProfileMap, timeSeriesData, timeseriesStartDate.Format("2006-01-02"), timeseriesEndDate.Format("2006-01-02"), sortBy, rank+1, query, searchFlag)
		leaderboard = append(leaderboard, *leaderboardProfile)
	}
	return &leaderboard, filteredCount, err
}

func (m *Manager) getCrossPlatformLeaderBoard(ctx context.Context, platform string, query coredomain.SearchQuery, sortBy string, sortDir string, page int, size int, source string) (*[]LeaderboardProfile, int64, error) {
	var leaderboard []LeaderboardProfile
	monthFlag, searchFlag := false, false
	for _, filter := range query.Filters {
		if filter.Field == "month" {
			monthFlag = true
		}
		if filter.Field == "search_phrase" {
			searchFlag = true
		}
	}
	if !monthFlag {
		lastMonthStart := time.Date(time.Now().AddDate(0, -1, 0).Year(), time.Now().AddDate(0, -1, 0).Month(), 1, 0, 0, 0, 0, time.Now().AddDate(0, -1, 0).Location()).Format("2006-01-02")
		query.Filters = append(query.Filters, coredomain.SearchFilter{FilterType: "EQ", Field: "month", Value: lastMonthStart})
	}
	query.Filters = append(query.Filters, coredomain.SearchFilter{FilterType: "EQ", Field: "platform", Value: platform})

	sortBy = getSortByKey(query, sortBy)
	if sortBy != "ASC" {
		sortDir = "ASC"
	}
	entities, filteredCount, err := m.leaderboardDao.Search(ctx, query, sortBy, sortDir, page, size)
	if err != nil {
		return &leaderboard, 0, err
	}
	// Last Month & Influencer Basic Data
	var instaPlatformIdArray, ytPlatformIdArray []int64
	// var platformIdArray []string
	for _, entity := range entities {
		if *entity.ProfilePlatform == "YA" {
			ytPlatformIdArray = append(ytPlatformIdArray, entity.PlatformProfileId)
		} else if *entity.ProfilePlatform == "IA" {
			instaPlatformIdArray = append(instaPlatformIdArray, entity.PlatformProfileId)
		}
	}
	instaPlatformProfileData, instaErr := m.discoveryManager.FindByPlatformProfileIds(ctx, string(constants.InstagramPlatform), instaPlatformIdArray, source)
	ytPlatformProfileData, ytErr := m.discoveryManager.FindByPlatformProfileIds(ctx, string(constants.YoutubePlatform), ytPlatformIdArray, source)

	var instaPlatformProfileMap = make(map[int64]coredomain.Profile)
	var ytPlatformProfileMap = make(map[int64]coredomain.Profile)

	if instaErr == nil {
		for _, insta := range *instaPlatformProfileData {
			platformCode, _ := strconv.ParseInt(insta.PlatformCode, 10, 64)
			instaPlatformProfileMap[platformCode] = insta
		}
	}
	if ytErr == nil {
		for _, yt := range *ytPlatformProfileData {
			platformCode, _ := strconv.ParseInt(yt.PlatformCode, 10, 64)
			ytPlatformProfileMap[platformCode] = yt
		}
	}
	ia_profile_id_map := make(map[string]int64)
	ya_profile_id_map := make(map[string]int64)
	profile_ids := make(map[string]bool)

	// append LinkedIds To Get Linked Profile TimeSeries
	if instaPlatformProfileData != nil {
		for _, insta := range *instaPlatformProfileData {
			platform_profile_id, _ := strconv.ParseInt(insta.PlatformCode, 10, 64)
			profile_id := *insta.Code
			profile_ids[profile_id] = true
			ia_profile_id_map[profile_id] = platform_profile_id
			for _, linked := range insta.LinkedSocials {
				if linked.Platform != string(constants.InstagramPlatform) {
					ytPlatformIdArray = append(ytPlatformIdArray, linked.Code)
					ya_profile_id_map[profile_id] = linked.Code
				}
			}
		}
	}
	if ytPlatformProfileData != nil {
		for _, yt := range *ytPlatformProfileData {
			platform_profile_id, _ := strconv.ParseInt(yt.PlatformCode, 10, 64)
			profile_id := *yt.Code
			profile_ids[profile_id] = true
			ya_profile_id_map[profile_id] = platform_profile_id
			for _, linked := range yt.LinkedSocials {
				if linked.Platform != string(constants.YoutubePlatform) {
					instaPlatformIdArray = append(instaPlatformIdArray, linked.Code)
					ia_profile_id_map[profile_id] = linked.Code
				}

			}
		}
	}
	timeseriesStartDate, timeseriesEndDate := getDateRanges(query)
	instaTimeSeriesData, inerr := m.socialProfileTimeSeriesDao.GetMultipleTimeSeriesData(ctx, timeseriesStartDate, timeseriesEndDate, string(constants.InstagramPlatform), instaPlatformIdArray)
	if inerr != nil {
		log.Error(inerr)
	}
	ytTimeSeriesData, yterr := m.socialProfileTimeSeriesDao.GetMultipleTimeSeriesData(ctx, timeseriesStartDate, timeseriesEndDate, string(constants.YoutubePlatform), ytPlatformIdArray)
	if yterr != nil {
		log.Error(yterr)
	}
	finalStructIA := ProfilesTimeSeries{
		FollowerGraph:  make(map[int64]MetricTimeSeries),
		TotalViewGraph: make(map[int64]MetricTimeSeries),
	}
	finalStructYA := ProfilesTimeSeries{
		FollowerGraph:  make(map[int64]MetricTimeSeries),
		TotalViewGraph: make(map[int64]MetricTimeSeries),
	}
	for profile_id := range profile_ids {
		ia_id := ia_profile_id_map[profile_id]
		ya_id := ya_profile_id_map[profile_id]

		ia_follower_ts := instaTimeSeriesData.FollowerGraph[ia_id]
		ya_follower_ts := ytTimeSeriesData.FollowerGraph[ya_id]
		profile_follower_ts := merge(ia_follower_ts, ya_follower_ts)

		ia_total_views_ts := instaTimeSeriesData.TotalPlaysGraph[ia_id]
		ya_total_views_ts := ytTimeSeriesData.TotalViewGraph[ya_id]
		profile_total_views_ts := merge(ia_total_views_ts, ya_total_views_ts)

		finalStructIA.FollowerGraph[ia_id] = profile_follower_ts
		finalStructYA.FollowerGraph[ya_id] = profile_follower_ts

		finalStructIA.TotalViewGraph[ia_id] = profile_total_views_ts
		finalStructYA.TotalViewGraph[ya_id] = profile_total_views_ts
	}

	for i := range entities {
		rank := i + (page-1)*size
		leaderboardProfile := GetCrossLeaderboardProfileFromEntity(entities[i], instaPlatformProfileMap, ytPlatformProfileMap, finalStructIA, finalStructYA, timeseriesStartDate.Format("2006-01-02"), timeseriesEndDate.Format("2006-01-02"), sortBy, rank+1, query, searchFlag)
		leaderboard = append(leaderboard, *leaderboardProfile)
	}
	return &leaderboard, filteredCount, err
}
func merge(iaTimeSeries MetricTimeSeries, ytTimeSeries MetricTimeSeries) MetricTimeSeries {
	final := MetricTimeSeries{}
	for key, value := range iaTimeSeries {
		final[key] = value
	}
	for key, value := range ytTimeSeries {
		if final[key] != nil {
			final[key] = helpers.ToInt64(*final[key] + *value)
		} else {
			final[key] = value
		}

	}
	return final
}

func getSortByKey(query coredomain.SearchQuery, sortBy string) string {
	//TO DO Will CHange To map in Refactoring
	if sortBy == "plays" {
		sortBy = "plays_rank"
	} else if sortBy == "views" {
		sortBy = "views_rank"
	} else if sortBy == "followers" {
		sortBy = "followers_change_rank"
	} else if sortBy == "followers_lifetime" {
		sortBy = "followers_rank"
	} else {
		sortBy = "followers_change_rank"
	}
	return sortBy
}

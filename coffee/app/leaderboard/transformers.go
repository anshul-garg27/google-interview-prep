package leaderboard

import (
	"coffee/constants"
	coredomain "coffee/core/domain"
	"coffee/helpers"
	"encoding/json"
	"strconv"
)

func ToLeaderboardEntry(entity LeaderboardEntity) (*LeaderboardEntry, error) {
	lastMonthRankMap := make(map[string]int64)
	if entity.LastMonthRanks != nil {
		json.Unmarshal([]byte(*entity.LastMonthRanks), &lastMonthRankMap)
	}

	currentMonthRankMap := make(map[string]int64)
	if entity.CurrentMonthRanks != nil {
		json.Unmarshal([]byte(*entity.CurrentMonthRanks), &currentMonthRankMap)
	}

	return &LeaderboardEntry{
		ID:                  entity.ID,
		Month:               entity.Month.Format("2006-01-02"),
		Language:            entity.Language,
		Category:            entity.Category,
		Platform:            entity.Platform,
		PlatformProfileId:   entity.PlatformProfileId,
		Handle:              entity.Handle,
		Followers:           entity.Followers,
		EngagementRate:      entity.EngagementRate,
		AvgLikes:            entity.AvgLikes,
		AvgComments:         entity.AvgComments,
		AvgViews:            entity.AvgViews,
		AvgPlays:            entity.AvgPlays,
		Enabled:             entity.Enabled,
		FollowersRank:       entity.FollowersRank,
		FollowersChangeRank: entity.FollowersChangeRank,
		ViewsRank:           entity.ViewsRank,
		Country:             entity.Country,
		Views:               entity.Views,
		PrevFollowers:       entity.PrevFollowers,
		FollowersChange:     entity.FollowersChange,
		Likes:               entity.Likes,
		Comments:            entity.Comments,
		Plays:               entity.Plays,
		Engagement:          entity.Engagement,
		PlaysRank:           entity.PlaysRank,
		ProfilePlatform:     entity.ProfilePlatform,
		YtViews:             entity.YtViews,
		IaViews:             entity.IaViews,
		LastMonthRanks:      lastMonthRankMap,
		PrevPlays:           entity.PrevPlays,
		PrevViews:           entity.PrevViews,
		Uploads:             entity.Uploads,
		SearchPhrase:        entity.SearchPhrase,
		CurrentMonthRanks:   currentMonthRankMap,
	}, nil
}

func GetLeaderboardProfileFromEntity(entity LeaderboardEntity, platformProfileMap map[int64]coredomain.Profile, timeSeriesData *ProfilesTimeSeries, startDate string, endDate string, sortBy string, rankInt int, query coredomain.SearchQuery, searchFlag bool) *LeaderboardProfile {
	entry, _ := ToLeaderboardEntry(entity)
	rank := int64(rankInt)
	rankChange, rank := getRankChange(*entry, query, sortBy, rank, searchFlag)
	platformProfileCode := entry.PlatformProfileId
	linkedSocials := platformProfileMap[platformProfileCode].LinkedSocials

	var followerGrowthPercentage, playsGrowthPercentage, viewsGrowthPercentage float64
	var viewsChange, playsChange int64
	if entry.PrevFollowers > 0 {
		followerGrowthPercentage = float64(entry.FollowersChange) / float64(entry.PrevFollowers) * 100
	}
	if entry.PrevPlays > 0 {
		playsChange = entry.Plays
		playsGrowthPercentage = float64(playsChange) / float64(entry.PrevPlays) * 100
	}
	if entry.PrevViews > 0 {
		viewsChange = entry.Views
		viewsGrowthPercentage = float64(viewsChange) / float64(entry.PrevViews) * 100
	}
	leaderboardProfile := &LeaderboardProfile{
		Id:           entry.ID,
		Rank:         rank,
		RankChange:   helpers.ToInt64(rankChange),
		Platform:     entry.Platform,
		PlatformCode: strconv.FormatInt(platformProfileCode, 10),
		Month:        entry.Month,
		Category:     entry.Category,
		Language:     entry.Language,
		Country:      entry.Country,
		AccountInfo: AccountInfo{
			Name:                     platformProfileMap[platformProfileCode].Name,
			Handle:                   platformProfileMap[platformProfileCode].Handle,
			Username:                 platformProfileMap[platformProfileCode].Username,
			Thumbnail:                platformProfileMap[platformProfileCode].Thumbnail,
			Code:                     platformProfileMap[platformProfileCode].Code,
			PlatformCode:             platformProfileMap[platformProfileCode].PlatformCode,
			Followers:                platformProfileMap[platformProfileCode].Metrics.Followers,
			Views:                    platformProfileMap[platformProfileCode].Metrics.Views,
			Categories:               platformProfileMap[platformProfileCode].Categories,
			EngagementRatePercenatge: platformProfileMap[platformProfileCode].Metrics.EngagementRatePercenatge,
			FollowersGrowth7d:        platformProfileMap[platformProfileCode].Metrics.FollowersGrowth7d,
			LinkedSocials:            linkedSocials,
			SearchPhrase:             entry.SearchPhrase,
		},
		LeaderboardMonthMetrics: LeaderboardMonthMetrics{
			Metrics: LeaderboardMetrics{
				Followers:                entry.Followers,
				Views:                    entry.Views,
				PrevViews:                entry.PrevViews,
				ViewsChange:              viewsChange,
				Likes:                    entry.Likes,
				Comments:                 entry.Comments,
				Plays:                    entry.Plays,
				PrevPlays:                entry.PrevPlays,
				PlaysChange:              playsChange,
				Engagement:               entry.Engagement,
				PrevFollowers:            entry.PrevFollowers,
				FollowersChange:          entry.FollowersChange,
				EngagementRate:           entry.EngagementRate,
				EngagementRatePercentage: entry.EngagementRate * 100,
				AvgLikes:                 entry.AvgLikes,
				AvgComments:              entry.AvgComments,
				AvgViews:                 entry.AvgViews,
				AvgPlays:                 entry.AvgPlays,
				Uploads:                  entry.Uploads,
				FollowerGrowthPercentage: followerGrowthPercentage,
				PlaysGrowthPercentage:    playsGrowthPercentage,
				ViewsGrowthPercentage:    viewsGrowthPercentage,
			},
		},
	}
	if entry.Platform == string(constants.InstagramPlatform) && leaderboardProfile != nil {
		leaderboardProfile.IsVerified = platformProfileMap[platformProfileCode].IsVerified
		if val, ok := timeSeriesData.FollowerGraph[platformProfileCode]; ok {
			val[startDate] = &entity.PrevFollowers
			val[endDate] = &entity.Followers
			leaderboardProfile.FollowerGraph = val
		} else {
			leaderboardProfile.FollowerGraph = MetricTimeSeries{}
		}

		if val, ok := timeSeriesData.TotalPlaysGraph[platformProfileCode]; ok {
			val[startDate] = &entity.PrevPlays
			currentPlays := entity.PrevPlays + entity.Plays
			val[endDate] = &currentPlays
			leaderboardProfile.PlaysGraph = val
		} else {
			leaderboardProfile.PlaysGraph = MetricTimeSeries{}
		}
	}
	if entry.Platform == string(constants.YoutubePlatform) && leaderboardProfile != nil {
		if val, ok := timeSeriesData.FollowerGraph[platformProfileCode]; ok {
			val[startDate] = &entity.PrevFollowers
			val[endDate] = &entity.Followers
			leaderboardProfile.FollowerGraph = val
		} else {
			leaderboardProfile.FollowerGraph = MetricTimeSeries{}
		}
		if val, ok := timeSeriesData.TotalViewGraph[platformProfileCode]; ok {
			val[startDate] = &entity.PrevViews
			currentViews := entity.PrevViews + entity.Views
			val[endDate] = &currentViews
			leaderboardProfile.ViewGraph = val
		} else {
			leaderboardProfile.ViewGraph = MetricTimeSeries{}
		}
	}
	return leaderboardProfile
}

func GetCrossLeaderboardProfileFromEntity(entity LeaderboardEntity, instaPlatformProfileMap map[int64]coredomain.Profile, ytPlatformProfileMap map[int64]coredomain.Profile, finalStructIA ProfilesTimeSeries, finalStructYA ProfilesTimeSeries, startDate string, endDate string, sortBy string, rankInt int, query coredomain.SearchQuery, searchFlag bool) *LeaderboardProfile {
	entry, _ := ToLeaderboardEntry(entity)
	rank := int64(rankInt)
	rankChange, rank := getRankChange(*entry, query, sortBy, rank, searchFlag)
	platformProfileCode := entry.PlatformProfileId
	var followerMap, totalViewMap map[string]*int64
	if val, ok := finalStructIA.FollowerGraph[platformProfileCode]; ok {

		val[startDate] = &entity.PrevFollowers
		val[endDate] = &entity.Followers
		followerMap = val
	} else if val, ok := finalStructYA.FollowerGraph[platformProfileCode]; ok {
		val[startDate] = &entity.PrevFollowers
		val[endDate] = &entity.Followers
		followerMap = val
	}
	if val, ok := finalStructIA.TotalViewGraph[platformProfileCode]; ok {
		val[startDate] = &entity.PrevViews
		currentViews := entity.PrevViews + entity.Views
		val[endDate] = &currentViews
		totalViewMap = val
	} else if val, ok := finalStructYA.TotalViewGraph[platformProfileCode]; ok {
		val[startDate] = &entity.PrevViews
		currentViews := entity.PrevViews + entity.Views
		val[endDate] = &currentViews
		totalViewMap = val
	}

	var platformProfileMap = make(map[int64]coredomain.Profile)
	var linkedSocials []*coredomain.SocialAccount
	if *entity.ProfilePlatform == "YA" {
		platformProfileMap = ytPlatformProfileMap
		linkedSocials = platformProfileMap[platformProfileCode].LinkedSocials

	} else {
		platformProfileMap = instaPlatformProfileMap
		linkedSocials = platformProfileMap[platformProfileCode].LinkedSocials
	}
	var followerGrowthPercentage, playsGrowthPercentage, viewsGrowthPercentage float64
	var viewsChange, playsChange int64
	if entry.PrevFollowers > 0 {
		followerGrowthPercentage = float64(entry.FollowersChange) / float64(entry.PrevFollowers) * 100
	}
	if entry.PrevPlays > 0 {
		playsChange = entry.Plays
		playsGrowthPercentage = float64(playsChange) / float64(entry.PrevPlays) * 100
	}
	if entry.PrevViews > 0 {
		viewsChange = entry.Views
		viewsGrowthPercentage = float64(viewsChange) / float64(entry.PrevViews) * 100
	}
	return &LeaderboardProfile{
		Id:           entry.ID,
		Rank:         rank,
		RankChange:   helpers.ToInt64(rankChange),
		Platform:     entry.Platform,
		PlatformCode: strconv.FormatInt(platformProfileCode, 10),
		Month:        entry.Month,
		Category:     entry.Category,
		Language:     entry.Language,
		Country:      entry.Country,
		AccountInfo: AccountInfo{
			Name:                     platformProfileMap[platformProfileCode].Name,
			Handle:                   platformProfileMap[platformProfileCode].Handle,
			Username:                 platformProfileMap[platformProfileCode].Username,
			Thumbnail:                platformProfileMap[platformProfileCode].Thumbnail,
			Code:                     platformProfileMap[platformProfileCode].Code,
			PlatformCode:             platformProfileMap[platformProfileCode].PlatformCode,
			Followers:                platformProfileMap[platformProfileCode].Metrics.Followers,
			Views:                    platformProfileMap[platformProfileCode].Metrics.Views,
			Categories:               platformProfileMap[platformProfileCode].Categories,
			EngagementRatePercenatge: platformProfileMap[platformProfileCode].Metrics.EngagementRatePercenatge,
			FollowersGrowth7d:        platformProfileMap[platformProfileCode].Metrics.FollowersGrowth7d,
			LinkedSocials:            linkedSocials,
			SearchPhrase:             entry.SearchPhrase,
		},
		FollowerGraph: followerMap,
		ViewGraph:     totalViewMap,
		LeaderboardMonthMetrics: LeaderboardMonthMetrics{
			Metrics: LeaderboardMetrics{
				Followers:                entry.Followers,
				Views:                    entry.Views,
				PrevViews:                entry.PrevViews,
				ViewsChange:              viewsChange,
				YtViews:                  entry.YtViews,
				IaViews:                  entry.IaViews,
				Likes:                    entry.Likes,
				Comments:                 entry.Comments,
				Plays:                    entry.Plays,
				Engagement:               entry.Engagement,
				PrevFollowers:            entry.PrevFollowers,
				FollowersChange:          entry.FollowersChange,
				EngagementRate:           entry.EngagementRate,
				EngagementRatePercentage: entry.EngagementRate * 100,
				AvgLikes:                 entry.AvgLikes,
				AvgComments:              entry.AvgComments,
				AvgViews:                 entry.AvgViews,
				AvgPlays:                 entry.AvgPlays,
				Uploads:                  entry.Uploads,
				FollowerGrowthPercentage: followerGrowthPercentage,
				PlaysGrowthPercentage:    playsGrowthPercentage,
				ViewsGrowthPercentage:    viewsGrowthPercentage,
			},
		},
	}
}
func getRankChange(entry LeaderboardEntry, query coredomain.SearchQuery, sortBy string, currentMonthRank int64, searchFlag bool) (int64, int64) {
	sortMap := make(map[string]bool)
	for _, filter := range query.Filters {
		if filter.Field != "month" {
			if filter.Value != "" {
				sortMap[filter.Field] = true
			}
		}
	}
	_, categoryOk := sortMap["category"]
	_, languageOk := sortMap["language"]
	_, profileTypeOk := sortMap["profile_type"]
	var lastMonthRankkey string
	if sortBy == "plays_rank" {
		if categoryOk && languageOk && profileTypeOk {
			lastMonthRankkey = "plays_rank_by_cat_lang_profile"
		} else if categoryOk && languageOk {
			lastMonthRankkey = "plays_rank_by_cat_lang"
		} else if categoryOk && profileTypeOk {
			lastMonthRankkey = "plays_rank_by_cat_profile"
		} else if languageOk && profileTypeOk {
			lastMonthRankkey = "plays_rank_by_lang_profile"
		} else if categoryOk {
			lastMonthRankkey = "plays_rank_by_cat"
		} else if languageOk {
			lastMonthRankkey = "plays_rank_by_lang"
		} else if profileTypeOk {
			lastMonthRankkey = "plays_rank_by_profile"
		} else {
			lastMonthRankkey = "plays_rank"
		}
	} else if sortBy == "views_rank" {
		if categoryOk && languageOk && profileTypeOk {
			lastMonthRankkey = "views_rank_by_cat_lang_profile"
		} else if categoryOk && languageOk {
			lastMonthRankkey = "views_rank_by_cat_lang"
		} else if categoryOk && profileTypeOk {
			lastMonthRankkey = "views_rank_by_cat_profile"
		} else if languageOk && profileTypeOk {
			lastMonthRankkey = "views_rank_by_lang_profile"
		} else if categoryOk {
			lastMonthRankkey = "views_rank_by_cat"
		} else if languageOk {
			lastMonthRankkey = "views_rank_by_lang"
		} else if profileTypeOk {
			lastMonthRankkey = "views_rank_by_profile"
		} else {
			lastMonthRankkey = "views_rank"
		}
	} else if sortBy == "followers_change_rank" {
		if categoryOk && languageOk && profileTypeOk {
			lastMonthRankkey = "followers_change_rank_by_cat_lang_profile"
		} else if categoryOk && languageOk {
			lastMonthRankkey = "followers_change_rank_by_cat_lang"
		} else if categoryOk && profileTypeOk {
			lastMonthRankkey = "followers_change_rank_by_cat_profile"
		} else if languageOk && profileTypeOk {
			lastMonthRankkey = "followers_change_rank_by_lang_profile"
		} else if categoryOk {
			lastMonthRankkey = "followers_change_rank_by_cat"
		} else if languageOk {
			lastMonthRankkey = "followers_change_rank_by_lang"
		} else if profileTypeOk {
			lastMonthRankkey = "followers_change_rank_by_profile"
		} else {
			lastMonthRankkey = "followers_change_rank"
		}
	} else if sortBy == "followers_rank" {
		if categoryOk && languageOk && profileTypeOk {
			lastMonthRankkey = "followers_rank_by_cat_lang_profile"
		} else if categoryOk && languageOk {
			lastMonthRankkey = "followers_rank_by_cat_lang"
		} else if categoryOk && profileTypeOk {
			lastMonthRankkey = "followers_rank_by_cat_profile"
		} else if languageOk && profileTypeOk {
			lastMonthRankkey = "followers_rank_by_lang_profile"
		} else if categoryOk {
			lastMonthRankkey = "followers_rank_by_cat"
		} else if languageOk {
			lastMonthRankkey = "followers_rank_by_lang"
		} else if profileTypeOk {
			lastMonthRankkey = "followers_rank_by_profile"
		} else {
			lastMonthRankkey = "followers_rank"
		}
	} else {
		lastMonthRankkey = "followers_change_rank"
	}
	lastMonthRank := entry.LastMonthRanks[lastMonthRankkey]
	if searchFlag {
		currentMonthRank = entry.CurrentMonthRanks[lastMonthRankkey]
	}
	var rankChange int64
	if lastMonthRank > 0 {
		rankChange = lastMonthRank - currentMonthRank
	}
	return rankChange, currentMonthRank
}

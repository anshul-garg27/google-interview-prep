package manager

import (
	"coffee/app/content/dao"
	"coffee/app/content/domain"
	discoverymanager "coffee/app/discovery/manager"
	"coffee/constants"

	coredomain "coffee/core/domain"
	"coffee/core/rest"
	"context"
	"strconv"
)

type TrendingManager struct {
	*rest.Manager[domain.TrendingContentEntry, dao.TrendingContentEntity, int64]
	trendingContentDao *dao.TrendingContentDao
	discoveryManager   *discoverymanager.SearchManager
}

func CreateTrendingManager(ctx context.Context) *TrendingManager {
	trendingContentDao := dao.CreateTrendingContentDao(ctx)
	trendingRestManager := rest.NewManager(ctx, trendingContentDao.Dao, dao.ToTrendingContentEntry, dao.ToTrendingContentEntity)
	discoveryManager := discoverymanager.CreateSearchManager(ctx)

	manager := &TrendingManager{
		Manager:            trendingRestManager,
		trendingContentDao: trendingContentDao,
		discoveryManager:   discoveryManager,
	}
	return manager
}

func (m *TrendingManager) SearchTrendingContent(ctx context.Context, query coredomain.SearchQuery, sortBy string, sortDir string, page int, size int, source string) ([]domain.TrendingContentEntry, int64, error) {
	var contentEntries []domain.TrendingContentEntry

	contentEntities, filteredCount, err := m.trendingContentDao.Search(ctx, query, sortBy, sortDir, page, size)
	if err != nil {
		return contentEntries, filteredCount, err
	}
	if len(contentEntities) > 0 {
		var inArray, ytArray []int64
		for _, contentEntity := range contentEntities {
			contentEntry, _ := dao.ToTrendingContentEntry(&contentEntity)
			contentEntries = append(contentEntries, *contentEntry)
			if contentEntry.Platform == string(constants.InstagramPlatform) {
				if contentEntry.PlatformAccountId != nil {
					inArray = append(inArray, *contentEntry.PlatformAccountId)
				}
			} else if contentEntry.Platform == string(constants.YoutubePlatform) {
				if contentEntry.PlatformAccountId != nil {
					ytArray = append(ytArray, *contentEntry.PlatformAccountId)
				}
			}
		}
		if len(inArray) > 0 {
			contentEntries = m.FindInstagramAccountInfo(ctx, contentEntries, inArray, source)
		}
		if len(ytArray) > 0 {
			contentEntries = m.FindYoutubeAccountInfo(ctx, contentEntries, ytArray, source)
		}
	}
	return contentEntries, filteredCount, nil
}

func (m *TrendingManager) FindYoutubeAccountInfo(ctx context.Context, contentEntries []domain.TrendingContentEntry, ytArray []int64, source string) []domain.TrendingContentEntry {
	var ytProfileMap = make(map[int64]coredomain.Profile)

	ytProfiles, err := m.discoveryManager.FindByPlatformProfileIds(ctx, string(constants.YoutubePlatform), ytArray, source)
	if err == nil {
		for _, profile := range *ytProfiles {
			platformCodeStr := profile.PlatformCode
			platformCode, _ := strconv.ParseInt(platformCodeStr, 10, 64)
			ytProfileMap[platformCode] = profile
		}
		for key := range contentEntries {
			if contentEntries[key].Platform == string(constants.YoutubePlatform) {
				var profile coredomain.Profile
				if contentEntries[key].PlatformAccountId != nil {
					platformCode := *contentEntries[key].PlatformAccountId
					if val, ok := ytProfileMap[platformCode]; ok {
						profile = val
						contentEntries[key].AccountInfo.Code = platformCode
						contentEntries[key].AccountInfo.Name = profile.Name
						contentEntries[key].AccountInfo.Platform = string(constants.YoutubePlatform)
						contentEntries[key].AccountInfo.Handle = profile.Handle
						contentEntries[key].AccountInfo.PlatformThumbnail = profile.Thumbnail
						contentEntries[key].AccountInfo.Followers = profile.Metrics.Followers
						contentEntries[key].AudienceAge = make(map[string]*float64)
						contentEntries[key].AudienceGender = make(map[string]*float64)
					}
				}

			}

		}
	}

	return contentEntries
}

func (m *TrendingManager) FindInstagramAccountInfo(ctx context.Context, contentEntries []domain.TrendingContentEntry, inArray []int64, source string) []domain.TrendingContentEntry {
	var inProfileMap = make(map[int64]coredomain.Profile)

	inProfiles, err := m.discoveryManager.FindByPlatformProfileIds(ctx, string(constants.InstagramPlatform), inArray, source)
	if err == nil {
		for _, profile := range *inProfiles {
			platformCodeStr := profile.PlatformCode
			platformCode, _ := strconv.ParseInt(platformCodeStr, 10, 64)
			inProfileMap[platformCode] = profile
		}
		for key := range contentEntries {
			if contentEntries[key].Platform == string(constants.InstagramPlatform) {
				var profile coredomain.Profile
				platformCode := *contentEntries[key].PlatformAccountId
				if val, ok := inProfileMap[platformCode]; ok {
					profile = val
				}
				contentEntries[key].AccountInfo.Code = platformCode
				contentEntries[key].AccountInfo.Name = profile.Name
				contentEntries[key].AccountInfo.Platform = string(constants.InstagramPlatform)
				contentEntries[key].AccountInfo.Handle = profile.Handle
				contentEntries[key].AccountInfo.PlatformThumbnail = profile.Thumbnail
				contentEntries[key].AccountInfo.Followers = profile.Metrics.Followers
				contentEntries[key].AudienceAge = make(map[string]*float64)
				contentEntries[key].AudienceGender = make(map[string]*float64)
			}

		}
	}
	return contentEntries
}

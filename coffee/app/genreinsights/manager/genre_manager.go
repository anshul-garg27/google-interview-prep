package manager

import (
	discoverydomain "coffee/app/discovery/domain"
	"coffee/app/genreinsights/dao"
	"coffee/app/genreinsights/domain"
	"coffee/constants"
	coredomain "coffee/core/domain"
	"coffee/core/rest"
	"context"
	"fmt"
	"strings"
)

type GenreManager struct {
	*rest.Manager[domain.GenreOverviewEntry, dao.GenreOverviewEntity, int64]
	genreInsightsDao *dao.GenreOverviewDao
}

func CreateGenreManager(ctx context.Context) *GenreManager {
	genreInsightsDao := dao.CreateGenreInsightsDao(ctx)
	genreRestManager := rest.NewManager(ctx, genreInsightsDao.Dao, dao.ToGenreOverviewEntry, dao.ToGenreOverviewEntity)

	manager := &GenreManager{
		Manager:          genreRestManager,
		genreInsightsDao: genreInsightsDao,
	}
	return manager
}

func (m *GenreManager) LanguageSearch(ctx context.Context, platform string, searchQuery coredomain.SearchQuery, sortBy string, sortDir string, page int, size int) (domain.LanguageMap, int64, error) {
	var entities []dao.GenreOverviewEntity
	var filteredCount int64
	var err error

	searchQuery.Filters = append(searchQuery.Filters, coredomain.SearchFilter{
		FilterType: "EQ",
		Field:      "platform",
		Value:      platform,
	})
	entities, filteredCount, err = m.genreInsightsDao.Search(ctx, searchQuery, sortBy, sortDir, page, size)
	if err != nil {
		return domain.LanguageMap{}, 0, err
	}

	if platform == "ALL" {
		var languages []string
		var languagesString string

		for _, entity := range entities {
			languages = append(languages, entity.Language)
		}
		languagesString = strings.Join(languages, ",")
		fmt.Println(languagesString)
		for key := range searchQuery.Filters {
			filterField := searchQuery.Filters[key].Field
			if filterField == "language" {
				searchQuery.Filters[key].Value = languagesString
				searchQuery.Filters[key].FilterType = "IN"
			}
			if filterField == "platform" {
				searchQuery.Filters[key].Value = "INSTAGRAM,YOUTUBE"
				searchQuery.Filters[key].FilterType = "IN"
			}
		}
		size = size * 2
		entities, filteredCount, err = m.genreInsightsDao.Search(ctx, searchQuery, sortBy, sortDir, page, size)
		if err != nil {
			return domain.LanguageMap{}, 0, err
		}
	}
	var languageMap domain.LanguageMap
	var ytentries []domain.GenreOverviewEntry
	var inentries []domain.GenreOverviewEntry

	for _, entity := range entities {
		entry, _ := dao.ToGenreOverviewEntry(&entity)
		entry.AudienceAgeGender = discoverydomain.AudienceAgeGender{}
		entry.AudienceGender = make(map[string]*float64)
		platform := entry.Platform
		if platform == string(constants.InstagramPlatform) {
			inentries = append(inentries, *entry)
		}
		if platform == string(constants.YoutubePlatform) {
			ytentries = append(ytentries, *entry)
		}
	}
	languageMap.Instagram = inentries
	languageMap.YOUTUBE = ytentries

	return languageMap, filteredCount, nil

}

package api

import (
	"coffee/app/genreinsights/dao"
	"coffee/app/genreinsights/domain"
	"coffee/app/genreinsights/manager"
	coredomain "coffee/core/domain"
	"strconv"

	"coffee/core/rest"
	"context"
)

type Service struct {
	*rest.Service[domain.GenreOverviewResponse, domain.GenreOverviewEntry, dao.GenreOverviewEntity, int64]
	GenreManager manager.GenreManager
}

func NewGenreInsightsService(ctx context.Context) *Service {
	gmanager := manager.CreateGenreManager(ctx)
	return &Service{
		GenreManager: *gmanager,
	}
}

func (r *Service) SearchGenreOverview(ctx context.Context, searchQuery coredomain.SearchQuery, sortBy string, sortDir string, page int, size int) domain.GenreOverviewResponse {
	searchQuery = rest.ParseGenderAgeFilters(searchQuery)
	entries, filteredCount, err := r.GenreManager.Search(ctx, searchQuery, sortBy, sortDir, page, size)
	if err != nil {
		return domain.CreateGenreErrorResponse(err, 0)
	}
	var nextCursor string
	if len(entries) >= size {
		nextCursor = strconv.Itoa(page + 1)
	} else {
		nextCursor = ""
	}
	return domain.CreateGenreResponse(entries, filteredCount, nextCursor, "Record(s) Retrieved Successfully")
}

func (r *Service) LanguageSearch(ctx context.Context, platform string, searchQuery coredomain.SearchQuery, sortBy string, sortDir string, page int, size int) domain.LanguageSplitResponse {
	searchQuery = rest.ParseGenderAgeFilters(searchQuery)
	languageSplit, filteredCount, err := r.GenreManager.LanguageSearch(ctx, platform, searchQuery, sortBy, sortDir, page, size)

	return domain.CreateLanguageSplitResponse(err, languageSplit, filteredCount, "", "Record(s) Retrieved Successfully")
}

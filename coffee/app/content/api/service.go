package api

import (
	"coffee/app/content/dao"
	"coffee/app/content/domain"
	"coffee/app/content/manager"
	coredomain "coffee/core/domain"
	"strconv"

	"coffee/core/rest"
	"context"
)

type Service struct {
	*rest.Service[domain.TrendingContentResponse, domain.TrendingContentEntry, dao.TrendingContentEntity, int64]
	TrendingManager manager.TrendingManager
}

func NewContentService(ctx context.Context) *Service {
	tmanager := manager.CreateTrendingManager(ctx)
	return &Service{
		TrendingManager: *tmanager,
	}
}

func (r *Service) SearchTrendingContent(ctx context.Context, searchQuery coredomain.SearchQuery, sortBy string, sortDir string, page int, size int, source string) domain.TrendingContentResponse {
	searchQuery = rest.ParseGenderAgeFilters(searchQuery)
	entries, filteredCount, err := r.TrendingManager.SearchTrendingContent(ctx, searchQuery, sortBy, sortDir, page, size, source)
	if err != nil {
		return domain.CreateTrendingErrorResponse(err, 0)
	}
	var nextCursor string
	if len(entries) >= size {
		nextCursor = strconv.Itoa(page + 1)
	} else {
		nextCursor = ""
	}
	return domain.CreateTrendingResponse(entries, filteredCount, nextCursor, "Record(s) Retrieved Successfully")
}

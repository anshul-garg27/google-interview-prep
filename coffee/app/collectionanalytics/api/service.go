package api

import (
	"coffee/app/collectionanalytics/domain"
	"coffee/app/collectionanalytics/manager"
	coredomain "coffee/core/domain"
	"coffee/core/rest"
	"context"
	"strconv"
)

type CollectionAnalyticsService struct {
	*rest.Service[coredomain.BaseResponse, coredomain.Entry, coredomain.BaseEntity, string]
	m *manager.CollectionAnalyticsManager
}

func (s *CollectionAnalyticsService) Init(ctx context.Context) {
	cmanager := manager.CreateCollectionAnalyticsManager(ctx)
	s.m = cmanager
}

func CreateCollectionAnalyticsService(ctx context.Context) *CollectionAnalyticsService {
	service := &CollectionAnalyticsService{}
	service.Init(ctx)
	return service
}

func (s *CollectionAnalyticsService) createCollectionTimeSeriesResponse(timeSeries []domain.CollectionTimeSeries) domain.CollectionTimeSeriesResponse {
	return domain.CollectionTimeSeriesResponse{
		Status: coredomain.Status{
			Status:        "SUCCESS",
			Message:       "Successfully retrieve collection timeseries",
			FilteredCount: 1,
			HasNextPage:   false,
			NextCursor:    "",
		},
		TimeSeriesData: &domain.CollectionTimeSeriesEntry{TimeSeries: timeSeries},
	}
}

func (s *CollectionAnalyticsService) createCollectionTimeSeriesErrorResponse(err error, code int) domain.CollectionTimeSeriesResponse {
	return domain.CollectionTimeSeriesResponse{
		Status: coredomain.Status{
			Status:        "ERROR",
			Message:       err.Error(),
			FilteredCount: 0,
			HasNextPage:   false,
			NextCursor:    "",
		},
		TimeSeriesData: nil,
	}
}

func (s *CollectionAnalyticsService) createCollectionMetricsSummaryResponse(metricsSummary *domain.CollectionMetricsSummaryEntry) domain.CollectionMetricsSummaryResponse {
	return domain.CollectionMetricsSummaryResponse{
		Status: coredomain.Status{
			Status:        "SUCCESS",
			Message:       "Successfully retrieve collection metrics",
			FilteredCount: 1,
			HasNextPage:   false,
			NextCursor:    "",
		},
		MetricsSummaryData: metricsSummary,
	}
}

func (s *CollectionAnalyticsService) createCollectionMetricsSummaryErrorResponse(err error, code int) domain.CollectionMetricsSummaryResponse {
	return domain.CollectionMetricsSummaryResponse{
		Status: coredomain.Status{
			Status:        "ERROR",
			Message:       err.Error(),
			FilteredCount: 0,
			HasNextPage:   false,
			NextCursor:    "",
		},
		MetricsSummaryData: nil,
	}
}

func (s *CollectionAnalyticsService) createCollectionProfileMetricsSummaryResponse(profiles []domain.CollectionProfileMetricsSummaryEntry, page int, size int) domain.CollectionProfileMetricsSummaryResponse {
	return domain.CollectionProfileMetricsSummaryResponse{
		Status: coredomain.Status{
			Status:        "SUCCESS",
			Message:       "Successfully retrieved collection profile metrics summary",
			FilteredCount: 1,
			HasNextPage:   len(profiles) >= size,
			NextCursor:    strconv.Itoa(page + 1),
		},
		Profiles: profiles,
	}
}

func (s *CollectionAnalyticsService) createCollectionProfileMetricsSummaryErrorResponse(err error, code int) domain.CollectionProfileMetricsSummaryResponse {
	return domain.CollectionProfileMetricsSummaryResponse{
		Status: coredomain.Status{
			Status:        "ERROR",
			Message:       err.Error(),
			FilteredCount: 0,
			HasNextPage:   false,
			NextCursor:    "",
		},
		Profiles: []domain.CollectionProfileMetricsSummaryEntry{},
	}
}

func (s *CollectionAnalyticsService) createCollectionPostMetricsSummaryResponse(posts []domain.CollectionPostMetricsSummaryEntry, page int, size int) domain.CollectionPostMetricsSummaryResponse {
	return domain.CollectionPostMetricsSummaryResponse{
		Status: coredomain.Status{
			Status:        "SUCCESS",
			Message:       "Successfully retrieved collection posts metrics summary",
			FilteredCount: 1,
			HasNextPage:   len(posts) >= size,
			NextCursor:    strconv.Itoa(page + 1),
		},
		Posts: posts,
	}
}

func (s *CollectionAnalyticsService) createCollectionPostMetricsSummaryErrorResponse(err error, code int) domain.CollectionPostMetricsSummaryResponse {
	return domain.CollectionPostMetricsSummaryResponse{
		Status: coredomain.Status{
			Status:        "ERROR",
			Message:       err.Error(),
			FilteredCount: 0,
			HasNextPage:   false,
			NextCursor:    "",
		},
		Posts: []domain.CollectionPostMetricsSummaryEntry{},
	}
}

func (s *CollectionAnalyticsService) createCollectionHashtagResponse(hashtags []domain.CollectionHashtagEntry) domain.CollectionHashtagResponse {
	return domain.CollectionHashtagResponse{
		Status: coredomain.Status{
			Status:        "SUCCESS",
			Message:       "Hashtags Retrieved Successfully",
			FilteredCount: 0,
			HasNextPage:   false,
			NextCursor:    "",
		},
		Hashtags: hashtags,
	}
}

func (s *CollectionAnalyticsService) createCollectionHashtagErrorResponse(err error, code int) domain.CollectionHashtagResponse {
	return domain.CollectionHashtagResponse{
		Status: coredomain.Status{
			Status:        "ERROR",
			Message:       err.Error(),
			FilteredCount: 0,
			HasNextPage:   false,
			NextCursor:    "",
		},
		Hashtags: []domain.CollectionHashtagEntry{},
	}
}

func (s *CollectionAnalyticsService) createCollectionKeywordResponse(keywords []domain.CollectionKeywordEntry) domain.CollectionKeywordResponse {
	return domain.CollectionKeywordResponse{
		Status: coredomain.Status{
			Status:        "SUCCESS",
			Message:       "Hashtags Retrieved Successfully",
			FilteredCount: 0,
			HasNextPage:   false,
			NextCursor:    "",
		},
		Keywords: keywords,
	}
}

func (s *CollectionAnalyticsService) createCollectionKeywordErrorResponse(err error, code int) domain.CollectionKeywordResponse {
	return domain.CollectionKeywordResponse{
		Status: coredomain.Status{
			Status:        "ERROR",
			Message:       err.Error(),
			FilteredCount: 0,
			HasNextPage:   false,
			NextCursor:    "",
		},
		Keywords: []domain.CollectionKeywordEntry{},
	}
}

func (s *CollectionAnalyticsService) createCollectionSentimentErrorResponse(err error, code int) domain.CollectionSentimentResponse {
	return domain.CollectionSentimentResponse{
		Status: coredomain.Status{
			Status:        "ERROR",
			Message:       err.Error(),
			FilteredCount: 0,
			HasNextPage:   false,
			NextCursor:    "",
		},
		SentimentData: nil,
	}
}

func (s *CollectionAnalyticsService) createCollectionSentimentResponse(sentiments domain.CollectionSentimentEntry) domain.CollectionSentimentResponse {
	return domain.CollectionSentimentResponse{
		Status: coredomain.Status{
			Status:        "SUCCESS",
			Message:       "Successfully retrieve collection timeseries",
			FilteredCount: 1,
			HasNextPage:   false,
			NextCursor:    "",
		},
		SentimentData: &sentiments,
	}
}

func (s *CollectionAnalyticsService) FetchCollectionMetricsSummary(ctx context.Context, query coredomain.SearchQuery) domain.CollectionMetricsSummaryResponse {
	metricsSummary, err := s.m.FetchCollectionMetricsSummary(ctx, query)
	if err != nil {
		return s.createCollectionMetricsSummaryErrorResponse(err, 0)
	}
	return s.createCollectionMetricsSummaryResponse(metricsSummary)
}

func (s *CollectionAnalyticsService) FetchCollectionTimeSeries(ctx context.Context, query coredomain.SearchQuery) domain.CollectionTimeSeriesResponse {
	timeSeries, err := s.m.FetchCollectionTimeSeries(ctx, query)
	if err != nil {
		return s.createCollectionTimeSeriesErrorResponse(err, 0)
	}
	return s.createCollectionTimeSeriesResponse(timeSeries)
}

func (s *CollectionAnalyticsService) FetchCollectionClicksTimeSeries(ctx context.Context, query coredomain.SearchQuery) domain.CollectionTimeSeriesResponse {
	timeSeries, err := s.m.FetchCollectionClicksTimeSeries(ctx, query)
	if err != nil {
		return s.createCollectionTimeSeriesErrorResponse(err, 0)
	}
	return s.createCollectionTimeSeriesResponse(timeSeries)
}

func (s *CollectionAnalyticsService) FetchCollectionProfilesWithMetricsSummary(ctx context.Context, query coredomain.SearchQuery, sortBy string, sortDir string, page int, size int) domain.CollectionProfileMetricsSummaryResponse {
	profiles, err := s.m.FetchCollectionProfilesWithMetricsSummary(ctx, query, sortBy, sortDir, page, size)
	if err != nil {
		return s.createCollectionProfileMetricsSummaryErrorResponse(err, 0)
	}
	return s.createCollectionProfileMetricsSummaryResponse(profiles, page, size)
}

func (s *CollectionAnalyticsService) FetchCollectionPostsWithMetricsSummary(ctx context.Context, query coredomain.SearchQuery, sortBy string, sortDir string, page int, size int) domain.CollectionPostMetricsSummaryResponse {
	posts, err := s.m.FetchCollectionPostsWithMetricsSummary(ctx, query, sortBy, sortDir, page, size)
	if err != nil {
		return s.createCollectionPostMetricsSummaryErrorResponse(err, 0)
	}
	return s.createCollectionPostMetricsSummaryResponse(posts, page, size)
}

func (s *CollectionAnalyticsService) FetchCollectionHashtags(ctx context.Context, collectionId string, collectionType string) domain.CollectionHashtagResponse {
	hashtags, err := s.m.FetchHashtagsForCollections(ctx, []string{collectionId}, collectionType)
	if err != nil {
		return s.createCollectionHashtagErrorResponse(err, 0)
	}
	return s.createCollectionHashtagResponse(hashtags)
}

func (s *CollectionAnalyticsService) FetchCollectionKeywords(ctx context.Context, collectionId string, collectionType string) domain.CollectionKeywordResponse {
	keywords, err := s.m.FetchKeywordsForCollections(ctx, []string{collectionId}, collectionType)
	if err != nil {
		return s.createCollectionKeywordErrorResponse(err, 0)
	}
	return s.createCollectionKeywordResponse(keywords)
}

func (s *CollectionAnalyticsService) FetchCollectionHashtagsShared(ctx context.Context, shareId string, collectionType string) domain.CollectionHashtagResponse {
	hashtags, err := s.m.FetchHashtagsForCollectionShared(ctx, shareId, collectionType)
	if err != nil {
		return s.createCollectionHashtagErrorResponse(err, 0)
	}
	return s.createCollectionHashtagResponse(hashtags)
}

func (s *CollectionAnalyticsService) FetchCollectionKeywordsShared(ctx context.Context, shareId string, collectionType string) domain.CollectionKeywordResponse {
	keywords, err := s.m.FetchKeywordsForCollectionShared(ctx, shareId, collectionType)
	if err != nil {
		return s.createCollectionKeywordErrorResponse(err, 0)
	}
	return s.createCollectionKeywordResponse(keywords)
}

func (s *CollectionAnalyticsService) FetchCollectionGroupHashtags(ctx context.Context, collectionGroupId int64, collectionType string) domain.CollectionHashtagResponse {
	collectionIds, err := s.m.GetCollectionIdFromGroupId(ctx, collectionGroupId)
	if err != nil {
		return s.createCollectionHashtagErrorResponse(err, 0)
	}
	hashtags, err := s.m.FetchHashtagsForCollections(ctx, collectionIds, collectionType)
	if err != nil {
		return s.createCollectionHashtagErrorResponse(err, 0)
	}
	return s.createCollectionHashtagResponse(hashtags)
}

func (s *CollectionAnalyticsService) FetchCollectionGroupKeywords(ctx context.Context, collectionGroupId int64, collectionType string) domain.CollectionKeywordResponse {
	collectionIds, err := s.m.GetCollectionIdFromGroupId(ctx, collectionGroupId)
	if err != nil {
		return s.createCollectionKeywordErrorResponse(err, 0)
	}
	keywords, err := s.m.FetchKeywordsForCollections(ctx, collectionIds, collectionType)
	if err != nil {
		return s.createCollectionKeywordErrorResponse(err, 0)
	}
	return s.createCollectionKeywordResponse(keywords)
}

func (s *CollectionAnalyticsService) FetchCollectionGroupHashtagsShared(ctx context.Context, shareId string, collectionType string) domain.CollectionHashtagResponse {
	collectionIds, err := s.m.GetCollectionIdFromGroupShareId(ctx, shareId)
	if err != nil {
		return s.createCollectionHashtagErrorResponse(err, 0)
	}
	hashtags, err := s.m.FetchHashtagsForCollections(ctx, collectionIds, collectionType)
	if err != nil {
		return s.createCollectionHashtagErrorResponse(err, 0)
	}
	return s.createCollectionHashtagResponse(hashtags)
}

func (s *CollectionAnalyticsService) FetchCollectionGroupKeywordsShared(ctx context.Context, shareId string, collectionType string) domain.CollectionKeywordResponse {
	collectionIds, err := s.m.GetCollectionIdFromGroupShareId(ctx, shareId)
	if err != nil {
		return s.createCollectionKeywordErrorResponse(err, 0)
	}
	keywords, err := s.m.FetchKeywordsForCollections(ctx, collectionIds, collectionType)
	if err != nil {
		return s.createCollectionKeywordErrorResponse(err, 0)
	}
	return s.createCollectionKeywordResponse(keywords)
}

func (s *CollectionAnalyticsService) FetchCollectionSentiment(ctx context.Context, query coredomain.SearchQuery) domain.CollectionSentimentResponse {
	sentiments, err := s.m.FetchCollectionSentiment(ctx, query)
	if err != nil {
		return s.createCollectionSentimentErrorResponse(err, 0)
	}
	return s.createCollectionSentimentResponse(*sentiments)
}

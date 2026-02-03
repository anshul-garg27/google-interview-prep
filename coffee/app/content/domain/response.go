package domain

import (
	coredomain "coffee/core/domain"
)

type TrendingContentResponse struct {
	Status          coredomain.Status       `json:"status"`
	TrendingContent *[]TrendingContentEntry `json:"trending"`
}

func CreateTrendingResponse(trending []TrendingContentEntry, filterCount int64, nextCursor string, message string) TrendingContentResponse {
	HasNextPage := false
	if nextCursor != "" {
		HasNextPage = true
	}
	return TrendingContentResponse{
		Status: coredomain.Status{
			Status:        "SUCCESS",
			Message:       message,
			FilteredCount: filterCount,
			HasNextPage:   HasNextPage,
			NextCursor:    nextCursor,
		},
		TrendingContent: &trending,
	}
}

func CreateTrendingErrorResponse(err error, code int) TrendingContentResponse {
	return TrendingContentResponse{
		Status:          coredomain.Status{Status: "500", Message: err.Error(), FilteredCount: 0},
		TrendingContent: &[]TrendingContentEntry{},
	}
}

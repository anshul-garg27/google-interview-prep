package domain

import (
	coredomain "coffee/core/domain"
)

type CollectionGroupResponse struct {
	Status coredomain.Status       `json:"status"`
	Groups *[]CollectionGroupEntry `json:"groups,omitempty"`
}

func CreateResponse(groups []CollectionGroupEntry, filterCount int64, nextCursor string, message string) CollectionGroupResponse {
	return CollectionGroupResponse{
		Status: coredomain.Status{
			Status:        "SUCCESS",
			Message:       message,
			FilteredCount: filterCount,
			HasNextPage:   !(nextCursor == ""),
			NextCursor:    nextCursor,
		},
		Groups: &groups,
	}
}

func CreateErrorResponse(err error, code int) CollectionGroupResponse {
	return CollectionGroupResponse{
		Status: coredomain.Status{
			Status:        "ERROR",
			Message:       err.Error(),
			FilteredCount: 0,
			HasNextPage:   false,
			NextCursor:    "",
		},
		Groups: nil,
	}
}

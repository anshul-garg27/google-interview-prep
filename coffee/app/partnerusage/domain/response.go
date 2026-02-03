package domain

import (
	coredomain "coffee/core/domain"
)

type PartnerResponse struct {
	Status  coredomain.Status `json:"status"`
	Partner *[]PartnerInput   `json:"partner"`
}

func CreatePartnerResponse(partner []PartnerInput, filterCount int64, nextCursor string, message string) PartnerResponse {
	HasNextPage := false
	if nextCursor != "" {
		HasNextPage = true
	}
	return PartnerResponse{
		Status: coredomain.Status{
			Status:        "SUCCESS",
			Message:       message,
			FilteredCount: filterCount,
			HasNextPage:   HasNextPage,
			NextCursor:    nextCursor,
		},
		Partner: &partner,
	}
}

func CreatePartnerErrorResponse(err error, code int) PartnerResponse {
	return PartnerResponse{
		Status:  coredomain.Status{Status: "ERROR", Message: err.Error(), FilteredCount: 0},
		Partner: nil,
	}
}

type ActivityTrackerResponse struct {
	Status     coredomain.Status `json:"status"`
	Activities *[]ActivityInput  `json:"activities"`
}

func CreateActivityTrackerResponse(activities []ActivityInput, filterCount int64, nextCursor string, message string) ActivityTrackerResponse {
	HasNextPage := false
	if nextCursor != "" {
		HasNextPage = true
	}
	return ActivityTrackerResponse{
		Status: coredomain.Status{
			Status:        "SUCCESS",
			Message:       message,
			FilteredCount: filterCount,
			HasNextPage:   HasNextPage,
			NextCursor:    nextCursor,
		},
		Activities: &activities,
	}
}

func CreateActivityTrackerErrorResponse(err error, code int) ActivityTrackerResponse {
	return ActivityTrackerResponse{
		Status:     coredomain.Status{Status: "ERROR", Message: err.Error(), FilteredCount: 0},
		Activities: nil,
	}
}

func CreateUnlockProfileResponse(err error) coredomain.Status {
	if err != nil {
		return coredomain.Status{
			Status:  "ERROR",
			Message: err.Error(),
		}
	}
	return coredomain.Status{
		Status:  "SUCCESS",
		Message: "Profile Unlocked Successfully",
	}

}

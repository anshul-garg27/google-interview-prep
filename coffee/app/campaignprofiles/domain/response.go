package domain

import (
	coredomain "coffee/core/domain"
)

type CampaignProfileResponse struct {
	Status   coredomain.Status                  `json:"status"`
	Profiles *[]coredomain.CampaignProfileEntry `json:"profiles"`
}

func CreateCampaignProfileResponse(Campaigns []coredomain.CampaignProfileEntry, filterCount int64, nextCursor string, message string) CampaignProfileResponse {
	HasNextPage := false
	if nextCursor != "" {
		HasNextPage = true
	}

	return CampaignProfileResponse{
		Status: coredomain.Status{
			Status:        "SUCCESS",
			Message:       message,
			FilteredCount: filterCount,
			HasNextPage:   HasNextPage,
			NextCursor:    nextCursor,
		},
		Profiles: &Campaigns,
	}
}

func CreateCampaignProfileErrorResponse(err error, code int) CampaignProfileResponse {
	return CampaignProfileResponse{
		Status: coredomain.Status{Status: "500", Message: err.Error(), FilteredCount: 0},
	}
}

type CampaignProfileAudienceInsightsResponse struct {
	Status   coredomain.Status                     `json:"status"`
	Profiles *CampaignProfileAudienceInsightsEntry `json:"profile"`
}

func CreateCampaignProfileAudienceInsightsResponse(audienceData *CampaignProfileAudienceInsightsEntry, filterCount int64, nextCursor string, message string) CampaignProfileAudienceInsightsResponse {
	HasNextPage := false
	if nextCursor != "" {
		HasNextPage = true
	}

	return CampaignProfileAudienceInsightsResponse{
		Status: coredomain.Status{
			Status:        "SUCCESS",
			Message:       message,
			FilteredCount: filterCount,
			HasNextPage:   HasNextPage,
			NextCursor:    nextCursor,
		},
		Profiles: audienceData,
	}
}

func CreateCampaignProfileAudienceInsightsErrorResponse(err error, code int) CampaignProfileAudienceInsightsResponse {
	return CampaignProfileAudienceInsightsResponse{
		Status: coredomain.Status{Status: "500", Message: err.Error(), FilteredCount: 0},
	}
}

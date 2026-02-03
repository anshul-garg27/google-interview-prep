package domain

import (
	coredomain "coffee/core/domain"
	"strconv"
)

type SocialProfileResponse struct {
	Status   coredomain.Status     `json:"status"`
	Profiles *[]coredomain.Profile `json:"profiles"`
}
type ProfileAudienceRespnse struct {
	Status   coredomain.Status                 `json:"status"`
	Audience *[]SocialProfileAudienceInfoEntry `json:"audienceData"`
}
type ProfileGrowthRespnse struct {
	Status coredomain.Status `json:"status"`
	Growth *[]Growth         `json:"growth,omitempty"`
}
type ProfileContentResponse struct {
	Status  coredomain.Status `json:"status"`
	Content *[]Content        `json:"content,omitempty"`
}
type ProfileHashtagResponse struct {
	Status   coredomain.Status              `json:"status"`
	Hashtags *[]SocialProfileHashatagsEntry `json:"hashtags"`
}
type LocationsResponse struct {
	Status    coredomain.Status `json:"status"`
	Locations *[]LocationsEntry `json:"locations"`
}
type GccLocationsResponse struct {
	Status    coredomain.Status `json:"status"`
	Locations []string          `json:"locations"`
}

func CreateSoialResponse(profiles *[]coredomain.Profile, filteredCount int64, page int, size int, err error, activityId int64) SocialProfileResponse {
	if err != nil {
		return SocialProfileResponse{
			Status:   coredomain.Status{Status: "500", Message: err.Error(), FilteredCount: filteredCount, ActivityId: activityId},
			Profiles: &[]coredomain.Profile{},
		}
	}
	if profiles == nil || len(*profiles) == 0 {
		message := "record not found"
		return SocialProfileResponse{
			Profiles: &[]coredomain.Profile{},
			Status:   coredomain.Status{Status: "SUCCESS", Message: message, FilteredCount: filteredCount, ActivityId: activityId},
		}

	}

	var nextCursor string
	hasNextPage := false
	if len(*profiles) >= size {
		nextCursor = strconv.Itoa(page + 1)
		hasNextPage = true
	} else {
		nextCursor = ""
	}
	return SocialProfileResponse{
		Status: coredomain.Status{
			Status:        "SUCCESS",
			Message:       "Profiles successfully retrieved",
			FilteredCount: filteredCount,
			HasNextPage:   hasNextPage,
			NextCursor:    nextCursor,
			ActivityId:    activityId,
		},
		Profiles: profiles,
	}

}

func CreateProfileAudienceResponse(err error, audienceData *[]SocialProfileAudienceInfoEntry, filterCount int64, nextCursor string) ProfileAudienceRespnse {
	HasNextPage := false
	if nextCursor != "" {
		HasNextPage = true
	}
	if err != nil {
		return ProfileAudienceRespnse{
			Status:   coredomain.Status{Status: "500", Message: err.Error(), FilteredCount: filterCount},
			Audience: &[]SocialProfileAudienceInfoEntry{},
		}
	}
	if audienceData == nil {
		return ProfileAudienceRespnse{
			Status: coredomain.Status{Status: "500", Message: err.Error(), FilteredCount: filterCount},
		}

	}
	return ProfileAudienceRespnse{
		Status: coredomain.Status{
			Status:        "SUCCESS",
			Message:       "Profile Audience successfully retrieved",
			FilteredCount: filterCount,
			HasNextPage:   HasNextPage,
			NextCursor:    nextCursor,
		},
		Audience: audienceData,
	}
}

func CreateGrowthResponse(err error, growthData *[]Growth, filterCount int64, nextCursor string) ProfileGrowthRespnse {
	HasNextPage := false
	if nextCursor != "" {
		HasNextPage = true
	}
	if err != nil {
		return ProfileGrowthRespnse{
			Status: coredomain.Status{Status: "500", Message: err.Error(), FilteredCount: filterCount},
			Growth: &[]Growth{},
		}
	}
	if len(*growthData) == 0 {
		return ProfileGrowthRespnse{
			Status: coredomain.Status{Status: "500", Message: "record not found", FilteredCount: filterCount},
		}

	}
	return ProfileGrowthRespnse{
		Status: coredomain.Status{
			Status:        "SUCCESS",
			Message:       "Profile Growth successfully retrieved",
			FilteredCount: filterCount,
			HasNextPage:   HasNextPage,
			NextCursor:    nextCursor,
		},
		Growth: growthData,
	}
}

func CreateContentResponse(err error, contentData *[]Content, filterCount int64, nextCursor string) ProfileContentResponse {
	HasNextPage := false
	if nextCursor != "" {
		HasNextPage = true
	}
	if err != nil {
		return ProfileContentResponse{
			Status:  coredomain.Status{Status: "500", Message: err.Error(), FilteredCount: filterCount},
			Content: &[]Content{},
		}
	}
	if len(*contentData) == 0 {
		return ProfileContentResponse{
			Status: coredomain.Status{Status: "500", Message: "record not found", FilteredCount: filterCount},
		}

	}
	return ProfileContentResponse{
		Status: coredomain.Status{
			Status:        "SUCCESS",
			Message:       "Profile Content successfully retrieved",
			FilteredCount: filterCount,
			HasNextPage:   HasNextPage,
			NextCursor:    nextCursor,
		},
		Content: contentData,
	}
}

func CreateHashtagResponse(err error, hashtagsData *[]SocialProfileHashatagsEntry, filterCount int64, nextCursor string) ProfileHashtagResponse {
	HasNextPage := false
	if nextCursor != "" {
		HasNextPage = true
	}
	if err != nil {
		return ProfileHashtagResponse{
			Status:   coredomain.Status{Status: "500", Message: err.Error(), FilteredCount: filterCount},
			Hashtags: &[]SocialProfileHashatagsEntry{},
		}
	}
	if len(*hashtagsData) == 0 {
		return ProfileHashtagResponse{
			Status: coredomain.Status{Status: "500", Message: err.Error(), FilteredCount: filterCount},
		}

	}
	return ProfileHashtagResponse{
		Status: coredomain.Status{
			Status:        "SUCCESS",
			Message:       "Profile Hashtags successfully retrieved",
			FilteredCount: filterCount,
			HasNextPage:   HasNextPage,
			NextCursor:    nextCursor,
		},
		Hashtags: hashtagsData,
	}
}

func CreateLocationsResponse(locationsData *[]LocationsEntry, filteredCount int64, page int, size int, err error) LocationsResponse {
	if err != nil {
		return LocationsResponse{
			Status:    coredomain.Status{Status: "500", Message: err.Error(), FilteredCount: filteredCount},
			Locations: &[]LocationsEntry{},
		}
	}
	if locationsData == nil || len(*locationsData) == 0 {
		message := "record not found"
		return LocationsResponse{
			Locations: &[]LocationsEntry{},
			Status:    coredomain.Status{Status: "SUCCESS", Message: message, FilteredCount: filteredCount},
		}

	}

	var nextCursor string
	hasNextPage := false
	if len(*locationsData) >= size {
		nextCursor = strconv.Itoa(page + 1)
		hasNextPage = true
	} else {
		nextCursor = ""
	}
	return LocationsResponse{
		Status: coredomain.Status{
			Status:        "SUCCESS",
			Message:       "Profile Locations successfully retrieved",
			FilteredCount: filteredCount,
			HasNextPage:   hasNextPage,
			NextCursor:    nextCursor,
		},
		Locations: locationsData,
	}

}

func CreateGccLocationsResponse(err error, locationsData []string) GccLocationsResponse {

	if err != nil {
		return GccLocationsResponse{
			Status:    coredomain.Status{Status: "500", Message: err.Error(), FilteredCount: 0},
			Locations: []string{},
		}
	}
	if len(locationsData) == 0 {
		return GccLocationsResponse{
			Status: coredomain.Status{Status: "500", Message: "record not found", FilteredCount: 0},
		}

	}
	return GccLocationsResponse{
		Status: coredomain.Status{
			Status:  "SUCCESS",
			Message: "Locations Data successfully retrieved",
		},
		Locations: locationsData,
	}
}

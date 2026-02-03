package domain

import (
	coredomain "coffee/core/domain"
)

func CreateGenreResponse(genreData []GenreOverviewEntry, filterCount int64, nextCursor string, message string) GenreOverviewResponse {
	HasNextPage := false
	if nextCursor != "" {
		HasNextPage = true
	}
	return GenreOverviewResponse{
		Status: coredomain.Status{
			Status:        "SUCCESS",
			Message:       message,
			FilteredCount: filterCount,
			HasNextPage:   HasNextPage,
			NextCursor:    nextCursor,
		},
		Genre: &genreData,
	}
}

func CreateGenreErrorResponse(err error, code int) GenreOverviewResponse {
	return GenreOverviewResponse{
		Status: coredomain.Status{Status: "500", Message: err.Error(), FilteredCount: 0},
		Genre:  &[]GenreOverviewEntry{},
	}
}

type GenreOverviewResponse struct {
	Status coredomain.Status     `json:"status"`
	Genre  *[]GenreOverviewEntry `json:"genre"`
}
type LanguageSplitResponse struct {
	Status      coredomain.Status `json:"status"`
	Languagemap LanguageMap       `json:"genre"`
}

func CreateLanguageSplitResponse(err error, languageSplit LanguageMap, filterCount int64, nextCursor string, message string) LanguageSplitResponse {
	HasNextPage := false
	if nextCursor != "" {
		HasNextPage = true
	}
	if err != nil {
		return LanguageSplitResponse{
			Status:      coredomain.Status{Status: "500", Message: err.Error(), FilteredCount: filterCount},
			Languagemap: LanguageMap{},
		}
	}
	// if profiles == nil {
	// 	return SocialProfileResponse{
	// 		Status: domain.Status{Status: "500", Message: err.Error(), FilteredCount: filterCount},
	// 	}

	// }
	return LanguageSplitResponse{
		Status: coredomain.Status{
			Status:        "SUCCESS",
			Message:       message,
			FilteredCount: filterCount,
			HasNextPage:   HasNextPage,
			NextCursor:    nextCursor,
		},
		Languagemap: languageSplit,
	}
}

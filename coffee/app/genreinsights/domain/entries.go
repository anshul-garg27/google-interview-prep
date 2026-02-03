package domain

import (
	discoverydomain "coffee/app/discovery/domain"
)

type GenreOverviewEntry struct {
	ID                int64                             `json:"id,omitempty"`
	Category          string                            `json:"category,omitempty"`
	Month             string                            `json:"month,omitempty"`
	Platform          string                            `json:"platform,omitempty"`
	Country           string                            `json:"country,omitempty"`
	Language          string                            `json:"language,omitempty"`
	ProfileType       string                            `json:"profileType,omitempty"`
	Creators          *int64                            `json:"creators,omitempty"`
	Uploads           *int64                            `json:"uploads,omitempty"`
	Views             *int64                            `json:"views,omitempty"`
	Plays             *int64                            `json:"plays,omitempty"`
	Followers         *int64                            `json:"followers,omitempty"`
	Likes             *int64                            `json:"likes,omitempty"`
	Comments          *int64                            `json:"comments,omitempty"`
	Engagement        *int64                            `json:"engagement"`
	EngagementRate    *float64                          `json:"engagementRate"`
	AudienceAgeGender discoverydomain.AudienceAgeGender `json:"audienceAgeGender,omitempty"`
	AudienceGender    map[string]*float64               `json:"audienceGender,omitempty"`
}

type LanguageMap struct {
	Instagram []GenreOverviewEntry `json:"instagram,omitempty"`
	YOUTUBE   []GenreOverviewEntry `json:"youtube,omitempty"`
}

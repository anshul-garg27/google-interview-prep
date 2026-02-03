package dao

import (
	discoverydomain "coffee/app/discovery/domain"
	"coffee/app/genreinsights/domain"
	"coffee/helpers"
	"encoding/json"
)

func ToGenreOverviewEntry(entity *GenreOverviewEntity) (*domain.GenreOverviewEntry, error) {
	audienceAgeGender := discoverydomain.AudienceAgeGender{}
	if entity.AudienceAgeGender != nil {
		json.Unmarshal([]byte(*entity.AudienceAgeGender), &audienceAgeGender)
	}
	var audienceGenderMap map[string]*float64
	if entity.AudienceGender != nil {
		json.Unmarshal([]byte(*entity.AudienceGender), &audienceGenderMap)
	}
	var engagementRate *float64
	if entity.EngagementRate != nil {
		engagementRate = helpers.ToFloat64(*entity.EngagementRate * 100)
	}
	return &domain.GenreOverviewEntry{
		ID:                entity.ID,
		Category:          entity.Category,
		Month:             entity.Month,
		Platform:          entity.Platform,
		Country:           entity.Country,
		Language:          entity.Language,
		ProfileType:       entity.ProfileType,
		Creators:          entity.Creators,
		Uploads:           entity.Uploads,
		Views:             entity.Views,
		Plays:             entity.Plays,
		Followers:         entity.Followers,
		Likes:             entity.Likes,
		Comments:          entity.Comments,
		Engagement:        entity.Engagement,
		EngagementRate:    engagementRate,
		AudienceAgeGender: audienceAgeGender,
		AudienceGender:    audienceGenderMap,
	}, nil
}

func ToGenreOverviewEntity(entry *domain.GenreOverviewEntry) (*GenreOverviewEntity, error) {
	return &GenreOverviewEntity{
		Category:       entry.Category,
		Month:          entry.Month,
		Platform:       entry.Platform,
		Country:        entry.Country,
		Language:       entry.Language,
		ProfileType:    entry.ProfileType,
		Creators:       entry.Creators,
		Uploads:        entry.Uploads,
		Views:          entry.Views,
		Plays:          entry.Plays,
		Followers:      entry.Followers,
		Likes:          entry.Likes,
		Comments:       entry.Comments,
		Engagement:     entry.Engagement,
		EngagementRate: entry.EngagementRate,
		// AudienceAgeGender: entry.AudienceAgeGender,
		Enabled: true,
	}, nil
}

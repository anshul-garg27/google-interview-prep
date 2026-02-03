package dao

import (
	"coffee/app/content/domain"
	"coffee/helpers"
	"encoding/json"
	"time"
)

func ToTrendingContentEntry(entity *TrendingContentEntity) (*domain.TrendingContentEntry, error) {
	var audienceGenderMap map[string]*float64
	if entity.AudienceGender != nil {
		json.Unmarshal([]byte(*entity.AudienceGender), &audienceGenderMap)
	}
	var audienceAgeMap map[string]*float64
	if entity.AudienceAge != nil {
		json.Unmarshal([]byte(*entity.AudienceAge), &audienceAgeMap)
	}
	var engagementRate *float64
	if entity.EngagementRate != nil {
		engagementRate = helpers.ToFloat64(*entity.EngagementRate * 100)
	}
	return &domain.TrendingContentEntry{
		ID:                entity.ID,
		Category:          entity.Category,
		Platform:          entity.Platform,
		Language:          entity.Language,
		Views:             entity.Views,
		Plays:             entity.Plays,
		Likes:             entity.Likes,
		Comments:          entity.Comments,
		Duration:          entity.Duration,
		PublishedAt:       entity.PublishedAt.Unix(),
		PostId:            entity.PostId,
		PostType:          entity.PostType,
		PostLink:          entity.PostLink,
		Thumbnail:         entity.Thumbnail,
		PlatformAccountId: entity.PlatformAccountId,
		Month:             entity.Month,
		Title:             entity.Title,
		Description:       entity.Description,
		Engagement:        entity.Engagement,
		EngagementRate:    engagementRate,
		AudienceAge:       audienceAgeMap,
		AudienceGender:    audienceGenderMap,
	}, nil
}

func ToTrendingContentEntity(entry *domain.TrendingContentEntry) (*TrendingContentEntity, error) {
	return &TrendingContentEntity{
		ID:                entry.ID,
		Category:          entry.Category,
		Platform:          entry.Platform,
		Language:          entry.Language,
		Views:             entry.Views,
		Plays:             entry.Plays,
		Likes:             entry.Likes,
		Comments:          entry.Comments,
		Duration:          entry.Duration,
		PublishedAt:       time.Unix(entry.PublishedAt, 0),
		PostId:            entry.PostId,
		PostType:          entry.PostType,
		PostLink:          entry.PostLink,
		Thumbnail:         entry.Thumbnail,
		PlatformAccountId: entry.PlatformAccountId,
		Month:             entry.Month,
		Title:             entry.Title,
		Description:       entry.Description,
		Engagement:        entry.Engagement,
		EngagementRate:    entry.EngagementRate,
	}, nil
}

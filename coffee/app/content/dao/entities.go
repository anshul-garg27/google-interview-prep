package dao

import "time"

type TrendingContentEntity struct {
	ID                int64 `gorm:"primaryKey;column:id"`
	Category          string
	Platform          string
	Language          string
	Views             *int64
	Plays             *int64
	Likes             *int64
	Comments          *int64
	Duration          *int64
	PublishedAt       time.Time
	PostId            string
	PostType          string
	PostLink          string
	Thumbnail         string
	PlatformAccountId *int64
	Month             string
	Title             *string
	Description       *string
	Engagement        *int64
	EngagementRate    *float64
	AudienceAge       *string
	AudienceGender    *string
	Enabled           bool
}

func (TrendingContentEntity) TableName() string {
	return "trending_content"
}

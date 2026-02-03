package leaderboard

import (
	"time"
)

type LeaderboardEntity struct {
	ID                  int64
	Month               time.Time
	Language            string
	Category            string
	Platform            string
	PlatformProfileId   int64
	Handle              string
	Followers           int64
	EngagementRate      float64
	AvgLikes            float64
	AvgComments         float64
	AvgViews            float64
	AvgPlays            float64
	Enabled             bool
	FollowersRank       int64
	FollowersChangeRank int64
	ViewsRank           int64
	Country             string
	Views               int64
	PrevFollowers       int64
	FollowersChange     int64
	Likes               int64
	Comments            int64
	Plays               int64
	Engagement          int64
	PlaysRank           int64
	ProfilePlatform     *string
	YtViews             *int64
	IaViews             *int64
	LastMonthRanks      *string
	PrevPlays           int64
	PrevViews           int64
	Uploads             *int64
	CurrentMonthRanks   *string
	SearchPhrase        *string
	CreatedAt           time.Time `gorm:"default:current_timestamp"`
	UpdatedAt           time.Time `gorm:"default:current_timestamp"`
}

func (LeaderboardEntity) TableName() string {
	return "leaderboard"
}

type SocialProfileTimeSeriesEntity struct {
	ID                int64 `gorm:"primaryKey;column:id"`
	PlatformProfileId int64
	Platform          *string
	Followers         *int64
	Following         *int64
	Views             *int64
	Plays             *int64
	ViewsTotal        *int64
	PlaysTotal        *int64
	Date              time.Time
	CreatedAt         time.Time `gorm:"default:current_timestamp"`
	UpdatedAt         time.Time `gorm:"default:current_timestamp"`
}

func (SocialProfileTimeSeriesEntity) TableName() string {
	return "social_profile_time_series"
}

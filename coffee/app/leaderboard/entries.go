package leaderboard

import (
	coredomain "coffee/core/domain"
)

type LeaderboardEntry struct {
	ID                  int64            `json:"ID"`
	Month               string           `json:"month"`
	Language            string           `json:"language"`
	Category            string           `json:"category"`
	Platform            string           `json:"platform"`
	PlatformProfileId   int64            `json:"platformProfileId"`
	Handle              string           `json:"handle"`
	Followers           int64            `json:"followers"`
	EngagementRate      float64          `json:"engagementRate"`
	AvgLikes            float64          `json:"avgLikes"`
	AvgComments         float64          `json:"avgComments"`
	AvgViews            float64          `json:"avgViews"`
	AvgPlays            float64          `json:"avgPlays"`
	Enabled             bool             `json:"enabled"`
	FollowersRank       int64            `json:"followersRank"`
	FollowersChangeRank int64            `json:"followersChangeRank"`
	ViewsRank           int64            `json:"viewsRank"`
	Country             string           `json:"country"`
	Views               int64            `json:"views"`
	PrevFollowers       int64            `json:"prevFollowers"`
	FollowersChange     int64            `json:"followersChange"`
	Likes               int64            `json:"likes"`
	Comments            int64            `json:"comments"`
	Plays               int64            `json:"plays"`
	Engagement          int64            `json:"engagement"`
	PlaysRank           int64            `json:"playsRank"`
	ProfilePlatform     *string          `json:"profilePlatform"`
	YtViews             *int64           `json:"ytViews"`
	IaViews             *int64           `json:"iaViews"`
	LastMonthRanks      map[string]int64 `json:"lastMonthRanks"`
	PrevPlays           int64            `json:"prevPlays"`
	PrevViews           int64            `json:"prevViews"`
	Uploads             *int64           `json:"uploads"`
	SearchPhrase        *string          `json:"searchPhrase"`
	CurrentMonthRanks   map[string]int64 `json:"currentMonthRanks"`
}
type AccountInfo struct {
	Name                     *string                     `json:"name,omitempty"`
	Handle                   *string                     `json:"handle,omitempty"`
	Username                 *string                     `json:"username,omitempty"`
	Thumbnail                *string                     `json:"thumbnail,omitempty"`
	Code                     *string                     `json:"code,omitempty"`
	PlatformCode             string                      `json:"platformCode,omitempty"`
	Followers                *int64                      `json:"followers,omitempty"`
	Views                    *int64                      `json:"views"`
	Categories               *[]string                   `json:"categories,omitempty"`
	EngagementRatePercenatge float64                     `json:"engagementRatePercentage,omitempty"`
	FollowersGrowth7d        *float64                    `json:"followersGrowth7d,omitempty"`
	LinkedSocials            []*coredomain.SocialAccount `json:"linkedSocials,omitempty"`
	SearchPhrase             *string                     `json:"searchPhrase"`
}
type LeaderboardProfile struct {
	Id                      int64                   `json:"id"`
	Rank                    int64                   `json:"rank"`
	RankChange              *int64                  `json:"rankchange,omitempty"`
	Platform                string                  `json:"platform"`
	PlatformCode            string                  `json:"platformCode"`
	Month                   string                  `json:"month"`
	Category                string                  `json:"category"`
	Language                string                  `json:"language"`
	Country                 string                  `json:"country"`
	IsVerified              *bool                   `json:"isVerified,omitempty"`
	AccountInfo             AccountInfo             `json:"accountInfo"`
	FollowerGraph           MetricTimeSeries        `json:"followerGraph,omitempty"`
	ViewGraph               MetricTimeSeries        `json:"viewGraph,omitempty"`
	PlaysGraph              MetricTimeSeries        `json:"playsGraph,omitempty"`
	LeaderboardMonthMetrics LeaderboardMonthMetrics `json:"leaderboardMetrics"`
}
type LeaderboardMonthMetrics struct {
	Metrics LeaderboardMetrics `json:"metrics"`
}
type LeaderboardMetrics struct {
	Followers                int64   `json:"followers,omitempty"`
	Views                    int64   `json:"views,omitempty"`
	PrevViews                int64   `json:"prevViews,omitempty"`
	ViewsChange              int64   `json:"viewsChange,omitempty"`
	YtViews                  *int64  `json:"youtubeViews,omitempty"`
	IaViews                  *int64  `json:"instagramViews,omitempty"`
	Likes                    int64   `json:"likes,omitempty"`
	Comments                 int64   `json:"comments,omitempty"`
	Plays                    int64   `json:"plays,omitempty"`
	PrevPlays                int64   `json:"prevPlays,omitempty"`
	PlaysChange              int64   `json:"playsChange,omitempty"`
	Engagement               int64   `json:"engagement,omitempty"`
	PrevFollowers            int64   `json:"prevFollowers,omitempty"`
	FollowersChange          int64   `json:"followersChange,omitempty"`
	EngagementRate           float64 `json:"engagementRate,omitempty"`
	EngagementRatePercentage float64 `json:"engagementRatePercentage,omitempty"`
	AvgLikes                 float64 `json:"avgLikes,omitempty"`
	AvgComments              float64 `json:"avgComments,omitempty"`
	AvgViews                 float64 `json:"avgViews,omitempty"`
	AvgPlays                 float64 `json:"avgPlays,omitempty"`
	FollowerGrowthPercentage float64 `json:"followerGrowthPercentage,omitempty"`
	PlaysGrowthPercentage    float64 `json:"playsGrowthPercentage,omitempty"`
	ViewsGrowthPercentage    float64 `json:"viewsGrowthPercentage,omitempty"`
	Uploads                  *int64  `json:"uploads"`
}
type MetricTimeSeries map[string]*int64

type ProfilesTimeSeries struct {
	FollowerGraph   map[int64]MetricTimeSeries
	ViewsGraph      map[int64]MetricTimeSeries
	PlaysGraph      map[int64]MetricTimeSeries
	TotalViewGraph  map[int64]MetricTimeSeries
	TotalPlaysGraph map[int64]MetricTimeSeries
}

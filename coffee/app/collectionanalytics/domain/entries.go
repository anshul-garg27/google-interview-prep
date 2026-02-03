package domain

import (
	"coffee/app/collectiongroup/domain"
	"time"
)

type CollectionMetricsSummaryEntry struct {
	CollectionId   string `json:"collectionId,omitempty"`
	CollectionType string `json:"collectionType,omitempty"`
	Notice         string `json:"notice,omitempty"`
	NoticeType     string `json:"noticeType,omitempty"`
	TotalFollowers int64  `json:"totalFollowers,omitempty"`
	TotalProfiles  int64  `json:"totalProfiles,omitempty"`
	TotalPosts     int64  `json:"totalPosts,omitempty"`
	Reach          int64  `json:"reach,omitempty"`
	Plays          int64  `json:"plays,omitempty"`
	Impressions    int64  `json:"impressions,omitempty"`
	Views          int64  `json:"views,omitempty"`
	Likes          int64  `json:"likes,omitempty"`
	Comments       int64  `json:"comments,omitempty"`
	SwipeUps       int64  `json:"swipeUps,omitempty"`
	StickerTaps    int64  `json:"stickerTaps,omitempty"`
	Shares         int64  `json:"shares,omitempty"`
	Saves          int64  `json:"saves,omitempty"`
	Mentions       int64  `json:"mentions,omitempty"`
	StoryExits     int64  `json:"storyExits,omitempty"`
	StoryBackTaps  int64  `json:"storyBackTaps,omitempty"`
	StoryFwdTaps   int64  `json:"storyFwdTaps,omitempty"`

	TotalEngagement int64  `json:"totalEngagement,omitempty"`
	EngagementRate  string `json:"engagementRate,omitempty"`
	ReachPerc       string `json:"reachPerc,omitempty"`
	ViewsPerc       string `json:"viewsPerc,omitempty"`
	ImpressionsPerc string `json:"impressionsPerc,omitempty"`

	ClickThroughRate     string `json:"clickThroughRate,omitempty"`
	ClickPerEngagement   string `json:"clickPerEngagement,omitempty"`
	CostPer1KImpressions string `json:"costPer1KImpressions,omitempty"`
	CostPerEngagement    string `json:"costPerEngagement,omitempty"`
	CostPerView          string `json:"costPerView,omitempty"`
	CostPerClick         string `json:"costPerClick,omitempty"`
	CostPerReach         string `json:"costPerReach,omitempty"`
	CostPerOrder         string `json:"costPerOrder,omitempty"`

	LinkClicks                 int64 `json:"linkClicks,omitempty"`
	Orders                     int64 `json:"orders,omitempty"`
	DeliveredOrders            int64 `json:"deliveredOrders,omitempty"`
	CompletedOrders            int64 `json:"completedOrders,omitempty"`
	LeaderboardOverallOrders   int64 `json:"leaderboardOverallOrders,omitempty"`
	LeaderboardDeliveredOrders int64 `json:"leaderboardDeliveredOrders,omitempty"`
	LeaderboardCompletedOrders int64 `json:"leaderboardCompletedOrders,omitempty"`

	Budget int `json:"budget,omitempty"`

	PostWiseCounts PostCountsSummary      `json:"postWiseCounts,omitempty"`
	WeeklyMetrics  []CollectionTimeSeries `json:"weeklyMetrics,omitempty"`

	RefreshedAt     string                         `json:"refreshedAt,omitempty"`
	DisabledMetrics []string                       `json:"disabledMetrics,omitempty"`
	Metadata        domain.CollectionGroupMetadata `json:"metadata,omitempty"`
}

type PostCountsSummary struct {
	Images    int `json:"images,omitempty"`
	Stories   int `json:"stories,omitempty"`
	Videos    int `json:"videos,omitempty"`
	Shorts    int `json:"shorts,omitempty"`
	Reels     int `json:"reels,omitempty"`
	Carousels int `json:"carousels,omitempty"`
}

type CollectionTimeSeriesEntry struct {
	TimeSeries []CollectionTimeSeries `json:"timeSeries,omitempty"`
}

type CollectionTimeSeries struct {
	StartOfWeek     string `json:"startOfWeek,omitempty"`
	Date            string `json:"date,omitempty"`
	TotalPosts      int64  `json:"totalPosts,omitempty"`
	Views           int64  `json:"views,omitempty"`
	Likes           int64  `json:"likes,omitempty"`
	Comments        int64  `json:"comments,omitempty"`
	Reach           int64  `json:"reach,omitempty"`
	Clicks          int64  `json:"clicks,omitempty"`
	Impressions     int64  `json:"impressions,omitempty"`
	TotalEngagement int64  `json:"totalEngagement,omitempty"`
}

// profiles: Int
// hashtagReach: Int
// followers: Int

type CollectionProfileMetricsSummaryEntry struct {
	Platform      string  `json:"platform,omitempty"`
	ProfileHandle string  `json:"profileHandle,omitempty"`
	ProfileCode   string  `json:"profileCode,omitempty"`
	ProfileName   *string `json:"profileName,omitempty"`
	ProfilePic    *string `json:"profilePic,omitempty"`
	Followers     int64   `json:"followers,omitempty"`
	Cost          int     `json:"cost,omitempty"`

	PostWiseCounts PostCountsSummary `json:"postWiseCounts,omitempty"`

	TotalPosts      int64 `json:"totalPosts,omitempty"`
	Reach           int64 `json:"reach,omitempty"`
	Plays           int64 `json:"plays,omitempty"`
	Impressions     int64 `json:"impressions,omitempty"`
	Views           int64 `json:"views,omitempty"`
	Likes           int64 `json:"likes,omitempty"`
	Comments        int64 `json:"comments,omitempty"`
	SwipeUps        int64 `json:"swipeUps,omitempty"`
	StickerTaps     int64 `json:"stickerTaps,omitempty"`
	Shares          int64 `json:"shares,omitempty"`
	Saves           int64 `json:"saves,omitempty"`
	Mentions        int64 `json:"mentions,omitempty"`
	StoryExits      int64 `json:"storyExits,omitempty"`
	StoryBackTaps   int64 `json:"storyBackTaps,omitempty"`
	StoryFwdTaps    int64 `json:"storyFwdTaps,omitempty"`
	TotalEngagement int64 `json:"totalEngagement,omitempty"`

	EngagementRate  string `json:"engagementRate,omitempty"`
	ReachPerc       string `json:"reachPerc,omitempty"`
	ImpressionsPerc string `json:"impressionsPerc,omitempty"`
	ViewsPerc       string `json:"viewsPerc,omitempty"`

	ClickThroughRate     string `json:"clickThroughRate,omitempty"`
	ClickPerEngagement   string `json:"clickPerEngagement,omitempty"`
	CostPer1KImpressions string `json:"costPer1KImpressions,omitempty"`
	CostPerEngagement    string `json:"costPerEngagement,omitempty"`
	CostPerView          string `json:"costPerView,omitempty"`
	CostPerClick         string `json:"costPerClick,omitempty"`
	CostPerReach         string `json:"costPerReach,omitempty"`

	LinkClicks                 int64 `json:"linkClicks,omitempty"`
	Orders                     int64 `json:"orders,omitempty"`
	DeliveredOrders            int64 `json:"deliveredOrders,omitempty"`
	CompletedOrders            int64 `json:"completedOrders,omitempty"`
	LeaderboardOverallOrders   int64 `json:"leaderboardOverallOrders,omitempty"`
	LeaderboardDeliveredOrders int64 `json:"leaderboardDeliveredOrders,omitempty"`
	LeaderboardCompletedOrders int64 `json:"leaderboardCompletedOrders,omitempty"`
}

type CollectionPostMetricsSummaryEntry struct {
	Platform             string    `json:"platform,omitempty"`
	PostShortCode        string    `json:"postShortCode,omitempty"`
	PostType             string    `json:"postType,omitempty"`
	PostTitle            string    `json:"postTitle,omitempty"`
	PostLink             string    `json:"postLink,omitempty"`
	PostThumbnail        *string   `json:"postThumbnail,omitempty"`
	Hashtags             []string  `json:"hashtags,omitempty"`
	PublishedAt          time.Time `json:"publishedAt,omitempty"`
	PostCollectionItemId *int64    `json:"postCollectionItemId,omitempty"`
	Bookmarked           *bool     `json:"bookmarked,omitempty"`

	ProfileHandle string  `json:"profileHandle,omitempty"`
	ProfileName   *string `json:"profileName,omitempty"`
	ProfilePic    *string `json:"profilePic,omitempty"`
	Followers     int64   `json:"followers,omitempty"`
	Cost          int     `json:"cost,omitempty"`

	Reach           int64 `json:"reach,omitempty"`
	Plays           int64 `json:"plays,omitempty"`
	Impressions     int64 `json:"impressions,omitempty"`
	Views           int64 `json:"views,omitempty"`
	Likes           int64 `json:"likes,omitempty"`
	Comments        int64 `json:"comments,omitempty"`
	SwipeUps        int64 `json:"swipeUps,omitempty"`
	StickerTaps     int64 `json:"stickerTaps,omitempty"`
	Shares          int64 `json:"shares,omitempty"`
	Saves           int64 `json:"saves,omitempty"`
	Mentions        int64 `json:"mentions,omitempty"`
	StoryExits      int64 `json:"storyExits,omitempty"`
	StoryBackTaps   int64 `json:"storyBackTaps,omitempty"`
	StoryFwdTaps    int64 `json:"storyFwdTaps,omitempty"`
	TotalEngagement int64 `json:"totalEngagement,omitempty"`

	EngagementRate  string `json:"engagementRate,omitempty"`
	ReachPerc       string `json:"reachPerc,omitempty"`
	ImpressionsPerc string `json:"impressionsPerc,omitempty"`
	ViewsPerc       string `json:"viewsPerc,omitempty"`

	ClickThroughRate     string `json:"clickThroughRate,omitempty"`
	ClickPerEngagement   string `json:"clickPerEngagement,omitempty"`
	CostPer1KImpressions string `json:"costPer1KImpressions,omitempty"`
	CostPerEngagement    string `json:"costPerEngagement,omitempty"`
	CostPerView          string `json:"costPerView,omitempty"`
	CostPerClick         string `json:"costPerClick,omitempty"`
	CostPerReach         string `json:"costPerReach,omitempty"`

	LinkClicks                 int64 `json:"linkClicks,omitempty"`
	Orders                     int64 `json:"orders,omitempty"`
	DeliveredOrders            int64 `json:"deliveredOrders,omitempty"`
	CompletedOrders            int64 `json:"completedOrders,omitempty"`
	LeaderboardOverallOrders   int64 `json:"leaderboardOverallOrders,omitempty"`
	LeaderboardDeliveredOrders int64 `json:"leaderboardDeliveredOrders,omitempty"`
	LeaderboardCompletedOrders int64 `json:"leaderboardCompletedOrders,omitempty"`
}

type CollectionHashtagEntry struct {
	HashtagName    string `json:"hashtagName,omitempty"`
	TaggedCount    int64  `json:"taggedCount,omitempty"`
	UGCTaggedCount int64  `json:"ugcTaggedCount,omitempty"`
}

type CollectionKeywordEntry struct {
	KeywordName string `json:"keywordName,omitempty"`
	TaggedCount int64  `json:"taggedCount,omitempty"`
}

type CollectionSentimentEntry struct {
	Sentiment []Sentiment `json:"sentiments,omitempty"`
}

type Sentiment struct {
	Type           string  `json:"type,omitempty"`
	PostPercentage float64 `json:"postPercentage,omitempty"`
}

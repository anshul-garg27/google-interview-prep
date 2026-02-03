package dao

import (
	"time"

	"github.com/jackc/pgtype"
)

type CollectionKeywordEntity struct {
	CollectionId      string `gorm:"column:collection_id"`
	CollectionType    string `gorm:"column:collection_type"`
	CollectionShareId string `gorm:"column:collection_share_id"`

	KeywordName string `gorm:"column:keyword_name"`
	TaggedCount int64  `gorm:"column:tagged_count"`

	UpdatedAt time.Time `gorm:"column:updated_at;type:time"`
}

type CollectionHashtagEntity struct {
	CollectionId      string `gorm:"column:collection_id"`
	CollectionType    string `gorm:"column:collection_type"`
	CollectionShareId string `gorm:"column:collection_share_id"`

	HashtagName    string `gorm:"column:hashtag_name"`
	TaggedCount    int64  `gorm:"column:tagged_count"`
	UGCTaggedCount int64  `gorm:"column:ugc_tagged_count"`

	UpdatedAt time.Time `gorm:"column:updated_at;type:time"`
}

type CollectionPostMetricsSummaryEntity struct {
	Platform      string       `gorm:"column:platform"`
	PostShortCode string       `gorm:"column:post_short_code"`
	PostType      string       `gorm:"column:post_type"`
	PostTitle     string       `gorm:"column:post_title"`
	PostLink      string       `gorm:"column:post_link"`
	PostThumbnail *string      `gorm:"column:post_thumbnail"`
	Hashtags      pgtype.JSONB `gorm:"column:hashtags;default:'[]'"`
	PublishedAt   time.Time    `gorm:"column:published_at"`

	CollectionId         string `gorm:"column:collection_id"`
	CollectionType       string `gorm:"column:collection_type"`
	CollectionShareId    string `gorm:"column:collection_share_id"`
	PostCollectionItemId *int64 `gorm:"column:post_collection_item_id"`
	Bookmarked           *bool  `gorm:"column:bookmarked"`

	ProfileHandle string    `gorm:"column:profile_handle"`
	ProfileName   *string   `gorm:"column:profile_name"`
	ProfilePic    *string   `gorm:"column:profile_pic"`
	Followers     int64     `gorm:"column:followers"`
	Cost          int       `gorm:"column:cost"`
	UpdatedAt     time.Time `gorm:"column:updated_at;type:time"`

	LinkClicks                 int64   `gorm:"column:link_clicks"`
	Orders                     int64   `gorm:"column:orders"`
	DeliveredOrders            int64   `gorm:"column:delivered_orders"`
	CompletedOrders            int64   `gorm:"column:completed_orders"`
	LeaderboardOverallOrders   int64   `gorm:"column:leaderboard_overall_orders"`
	LeaderboardDeliveredOrders int64   `gorm:"column:leaderboard_delivered_orders"`
	LeaderboardCompletedOrders int64   `gorm:"column:leaderboard_completed_orders"`
	Views                      int64   `gorm:"column:views"`
	Likes                      int64   `gorm:"column:likes"`
	Comments                   int64   `gorm:"column:comments"`
	Impressions                int64   `gorm:"column:impressions"`
	Saves                      int64   `gorm:"column:saves"`
	Plays                      int64   `gorm:"column:plays"`
	Reach                      int64   `gorm:"column:reach"`
	SwipeUps                   int64   `gorm:"column:swipe_ups"`
	Mentions                   int64   `gorm:"column:mentions"`
	StickerTaps                int64   `gorm:"column:sticker_taps"`
	Shares                     int64   `gorm:"column:shares"`
	StoryExits                 int64   `gorm:"column:story_exits"`
	StoryBackTaps              int64   `gorm:"column:story_back_taps"`
	StoryFwdTaps               int64   `gorm:"column:story_forward_taps"`
	ER                         float64 `gorm:"<-;column:er"`
}

type CollectionPostMetricsTSEntity struct {
	Platform      string    `gorm:"column:platform"`
	PostShortCode string    `gorm:"column:post_short_code"`
	PostType      string    `gorm:"column:post_type"`
	PostTitle     string    `gorm:"column:post_title"`
	PostLink      string    `gorm:"column:post_link"`
	PostThumbnail string    `gorm:"column:post_thumbnail"`
	PublishedAt   time.Time `gorm:"column:published_at"`

	CollectionId         string `gorm:"column:collection_id"`
	CollectionType       string `gorm:"column:collection_type"`
	CollectionShareId    string `gorm:"column:collection_share_id"`
	PostCollectionItemId *int64 `gorm:"column:post_collection_item_id"`

	ProfileHandle string    `gorm:"column:profile_handle"`
	ProfileName   *string   `gorm:"column:profile_name"`
	ProfilePic    *string   `gorm:"column:profile_pic"`
	Followers     int64     `gorm:"column:followers"`
	Cost          int       `gorm:"column:cost"`
	UpdatedAt     time.Time `gorm:"column:updated_at;type:time"`

	StatsDate     time.Time `gorm:"column:stats_date"`
	InsightSource string    `gorm:"column:insight_source"`
	Views         int64     `gorm:"column:views"`
	Likes         int64     `gorm:"column:likes"`
	Comments      int64     `gorm:"column:comments"`
	Impressions   int64     `gorm:"column:impressions"`
	Saves         int64     `gorm:"column:saves"`
	Plays         int64     `gorm:"column:plays"`
	Reach         int64     `gorm:"column:reach"`
	SwipeUps      int64     `gorm:"column:swipe_ups"`
	Mentions      int64     `gorm:"column:mentions"`
	StickerTaps   int64     `gorm:"column:sticker_taps"`
	Shares        int64     `gorm:"column:shares"`
	StoryExits    int64     `gorm:"column:story_exits"`
	StoryBackTaps int64     `gorm:"column:story_back_taps"`
	StoryFwdTaps  int64     `gorm:"column:story_forward_taps"`
	ER            float64   `gorm:"<-;column:er"`
}

func (CollectionPostMetricsSummaryEntity) TableName() string {
	return "collection_post_metrics_summary"
}

func (CollectionPostMetricsTSEntity) TableName() string {
	return "collection_post_metrics_ts"
}

func (CollectionHashtagEntity) TableName() string {
	return "collection_hashtag"
}

func (CollectionKeywordEntity) TableName() string {
	return "collection_keyword"
}

type MetricsSummaryEntity struct {
	Totalposts                  int64
	Reach                       int64
	Plays                       int64
	Impressions                 int64
	Views                       int64
	Likes                       int64
	Comments                    int64
	SwipeUps                    int64
	StickerTaps                 int64
	Shares                      int64
	Saves                       int64
	Mentions                    int64
	Totalengagement             int64
	Storyexits                  int64
	Storybacktaps               int64
	Storyfwdtaps                int64
	Linkclicks                  int64
	Orders                      int64
	Deliveredorders             int64
	Completedorders             int64
	Leaderboardoverallorders    int64
	LeaderboarddeliveredOorders int64
	Leaderboardcompletedorders  int64
	Followers                   int64
	Profilescount               int64
	Er                          float64
	Lastrefreshtime             time.Time
}

type WeeklyMetricsSummary struct {
	Startofweek time.Time
	Totalposts  int64
	Views       float64
	Likes       int64
	Comments    int64
	Clicks      int64
	Reach       float64
	Impressions float64
}

type PostWiseCountEntity struct {
	Posttype string
	Posts    int
}

type ProfilePostWiseCountEntity struct {
	Platform      string
	Profilehandle string
	Posttype      string
	Posts         int
}

type ProfileMetricsSummaryEntity struct {
	Platform                    string
	Handle                      string
	Name                        *string
	Profilepic                  *string
	Followers                   int64
	Cost                        int
	Totalposts                  int64
	Reach                       int64
	Plays                       int64
	Impressions                 int64
	Views                       int64
	Likes                       int64
	Comments                    int64
	SwipeUps                    int64
	StickerTaps                 int64
	Shares                      int64
	Saves                       int64
	Mentions                    int64
	Totalengagement             int64
	Er                          float64
	Storyexits                  int64
	Storybacktaps               int64
	Storyfwdtaps                int64
	Linkclicks                  int64
	Orders                      int64
	Deliveredorders             int64
	Completedorders             int64
	Leaderboardoverallorders    int64
	LeaderboarddeliveredOorders int64
	Leaderboardcompletedorders  int64
	ProfileSocialId             *string

	Lastrefreshtime time.Time
}

type Total struct {
	Total int64
}

type SentimentCount struct {
	Count     int64
	Sentiment string
}

type KeywordScore struct {
	Keyword string
	Score   float64
}

package dao

import (
	campaignprofilesdao "coffee/app/campaignprofiles/dao"

	"time"

	"github.com/lib/pq"
)

type InstagramAccountEntity struct {
	ID                        int64 `gorm:"<-:create,primaryKey,id"`
	Name                      *string
	Handle                    *string
	BusinessID                *string `gorm:"column:business_id"`
	IgID                      *string `gorm:"column:ig_id"`
	Thumbnail                 *string
	DOB                       *time.Time
	Bio                       *string
	Categories                pq.StringArray `gorm:"column:categories;type:text[]"`
	Keywords                  pq.StringArray `gorm:"column:keywords;type:text[]"`
	KeywordsAdmin             pq.StringArray `gorm:"column:keywords_admin;type:text[]"`
	Label                     *string
	Gender                    *string
	City                      *string
	State                     *string
	Country                   *string
	IsBlacklisted             *bool   `gorm:"column:is_blacklisted"`
	ProfileId                 *string `gorm:"column:profile_id"`
	Followers                 *int64
	Following                 *int64
	PostCount                 *int64
	Ffratio                   *float64
	AvgViews                  *float64
	AvgLikes                  *float64
	AvgReach                  *float64
	AvgReelsPlayCount         *float64
	EngagementRate            float64
	FollowersGrowth7d         *float64
	ReelsReach                *int64
	StoryReach                *int64
	ImageReach                *int64
	AvgComments               *float64
	ErGrade                   *string
	LinkedSocials             *string `gorm:"column:linked_socials"`
	LinkedChannelId           *string
	AudienceGender            *string
	AudienceLocation          *string
	EstPostPrice              *string
	EstReach                  *string
	EstImpressions            *string
	FlagContactInfoAvailable  bool
	IsVerified                bool `gorm:"column:is_verified"`
	PostFrequencyWeek         *float64
	CountryRank               *int64
	CategoryRank              *string
	AuthenticEngagement       *int64
	CommentRatePercentage     *float64
	LikesSpreadPercentage     *float64
	GroupKey                  *string
	SearchPhrase              *string
	SearchPhraseAdmin         *string
	Phone                     *string
	Email                     *string
	AvgCommentsGrade          *string
	AvgLikesGrade             *string
	CommentsRateGrade         *string
	FollowersGrade            *string
	EngagementRateGrade       *string
	ReelsReachGrade           *string
	LikesToCommentRatioGrade  *string
	FollowersGrowth7dGrade    *string
	FollowersGrowth30dGrade   *string
	FollowersGrowth90dGrade   *string
	AudienceReachabilityGrade *string
	AudienceAuthencityGrade   *string
	AudienceQualityGrade      *string
	PostCountGrade            *string
	LikesSpreadGrade          *string
	FollowersGrowth30d        *string
	SimilarProfileGroupData   *string
	ImageReachGrade           *string
	StoryReachGrade           *string
	AvgReelsPlay30d           *int64 `gorm:"column:avg_reels_play_30d"`
	Plays30d                  *int64 `gorm:"column:plays_30d"`
	Uploads30d                *int64 `gorm:"column:uploads_30d"`
	AccountType               *string
	Enabled                   bool                                       `gorm:"default:true"`
	EnabledForSaas            bool                                       `gorm:"default:false"`
	CampaignProfile           *campaignprofilesdao.CampaignProfileEntity `gorm:"foreignKey:PlatformAccountId;references:id"`
	GccLinkedChannelId        *string
	GccProfileId              *string
	Location                  pq.StringArray `gorm:"column:location;type:text[]"`
	LocationList              pq.StringArray `gorm:"column:location_list;type:text[]"`
	Deleted                   *bool          `gorm:"default:false"`
	Impressions30d            *float64       `gorm:"column:impressions_30d"`
	ImageImpressionsGrade     *string
	ReelsImpressionsGrade     *string
	Reach30d                  *float64 `gorm:"column:reach_30d"`
	IsPrivate                 *bool
	Languages                 pq.StringArray `gorm:"column:languages;type:text[]"`
	CreatedAt                 time.Time      `gorm:"default:current_timestamp"`
	UpdatedAt                 time.Time      `gorm:"default:current_timestamp"`
}

func (InstagramAccountEntity) TableName() string {
	return "instagram_account"
}

type YoutubeAccountEntity struct {
	ID                       int64   `gorm:"<-create:primaryKey"`
	ChannelId                *string `gorm:"column:channel_id"`
	Title                    *string `gorm:"column:title"`
	Username                 *string
	Gender                   *string
	ProfileId                *string `gorm:"column:profile_id"`
	Description              *string
	Dob                      *string
	Thumbnail                *string
	Categories               pq.StringArray `gorm:"column:categories;type:text[]"`
	Keywords                 pq.StringArray `gorm:"column:keywords;type:text[]"`
	KeywordsAdmin            pq.StringArray `gorm:"column:keywords_admin;type:text[]"`
	SearchPhrase             *string
	SearchPhraseAdmin        *string
	Label                    *string
	IsBlacklisted            *bool `gorm:"column:is_blacklisted"`
	City                     *string
	State                    *string
	Country                  *string
	UploadsCount             *int64 `gorm:"column:uploads_count"`
	Followers                *int64 `gorm:"column:followers"`
	FollowersGrowth7d        *float64
	ViewsCount               *int64   `gorm:"column:views_count"`
	AvgViews                 *float64 `gorm:"column:avg_views"`
	VideoReach               *int64   `gorm:"column:video_reach"`
	ShortsReach              *int64   `gorm:"column:shorts_reach"`
	LinkedSocials            *string  `gorm:"column:linked_socials"`
	LinkedHandle             *string  `gorm:"column:linked_instagram_handle"`
	AudienceGender           *string
	AudienceLocation         *string
	EstPostPrice             *string
	EstReach                 *string
	EstImpressions           *string
	FlagContactInfoAvailable bool
	CountryRank              *int64
	CategoryRank             *string
	AuthenticEngagement      *int64
	CommentRatePercentage    *float64
	LikesSpreadPercentage    *float64
	VideoViewsLast30         *int64 `gorm:"column:video_views_last30"`
	GroupKey                 *string
	Phone                    *string
	Email                    *string
	CommentsRateGrade        *string
	FollowersGrade           *string
	FollowersGrowth7dGrade   *string
	ReactionRate             *int64
	CommentsRate             *int64
	Cpm                      *int64
	AvgPostsPerWeek          *int64
	AvgPostsPerWeekGrade     *string
	FollowersGrowth1y        *int64
	FollowersGrowth1yGrade   *string
	AvgShortsViews30d        *int64 `gorm:"column:avg_shorts_views_30d"`
	AvgVideoViews30d         *int64 `gorm:"column:avg_video_views_30d"`
	LatestVideoPublishTime   *string
	SimilarProfileGroupData  *string
	CampaignProfile          *campaignprofilesdao.CampaignProfileEntity `gorm:"foreignKey:PlatformAccountId;references:id"`
	Views30dGrade            *string                                    `gorm:"column:views_30d_grade"`
	Views30d                 *int64                                     `gorm:"column:views_30d"`
	Uploads30d               *int64                                     `gorm:"column:uploads_30d"`
	Enabled                  bool                                       `gorm:"default:true"`
	EnabledForSaas           bool                                       `gorm:"default:false"`
	GccLinkedInstagramHandle *string
	GccProfileId             *string
	Location                 pq.StringArray `gorm:"column:location;type:text[]"`
	LocationList             pq.StringArray `gorm:"column:location_list;type:text[]"`
	Deleted                  *bool          `gorm:"default:false"`
	Languages                pq.StringArray `gorm:"column:languages;type:text[]"`
	CreatedAt                time.Time      `gorm:"default:current_timestamp"`
	UpdatedAt                time.Time      `gorm:"default:current_timestamp"`
}

func (YoutubeAccountEntity) TableName() string {
	return "youtube_account"
}

type SocialProfileTimeSeriesEntity struct {
	ID                int64 `gorm:"primaryKey;column:id"`
	PlatformProfileId int64
	Platform          string
	Followers         *int64
	Following         *int64
	Views             *int64
	Plays             *int64
	Uploads           *int64
	EngagementRate    *float64
	FollowersChange   *int64
	ViewsTotal        *int64
	PlaysTotal        *int64
	Date              time.Time
	PrevDate          time.Time
	MonthlyStats      bool
	CreatedAt         time.Time `gorm:"default:current_timestamp"`
	UpdatedAt         time.Time `gorm:"default:current_timestamp"`
}

func (SocialProfileTimeSeriesEntity) TableName() string {
	return "social_profile_time_series"
}

type SocialProfileHashatagsEntity struct {
	ID                int64 `gorm:"primaryKey;column:id"`
	PlatformProfileId int64
	Platform          string
	Hashtags          pq.StringArray `gorm:"column:hashtags;type:text[]"`
	HashtagsCounts    pq.Int64Array  `gorm:"column:hashtags_counts;type:integer[]"`
	CreatedAt         time.Time      `gorm:"default:current_timestamp"`
	UpdatedAt         time.Time      `gorm:"default:current_timestamp"`
}

func (SocialProfileHashatagsEntity) TableName() string {
	return "social_profile_hashtags"
}

type LocationsEntity struct {
	ID        int64 `gorm:"primaryKey;column:id"`
	Name      string
	FullName  string
	Type      string
	CreatedAt time.Time `gorm:"default:current_timestamp"`
	UpdatedAt time.Time `gorm:"default:current_timestamp"`
}

func (LocationsEntity) TableName() string {
	return "locations"
}

type SocialProfileAudienceInfoEntity struct {
	ID                             int64 `gorm:"primaryKey;column:id"`
	PlatformProfileId              int64
	Platform                       string
	AudienceGenderSplit            *string
	AudienceAgeGenderSplit         *string
	AudienceLanguage               *string
	NotableFollowers               pq.StringArray `gorm:"column:notable_followers;type:text[]"`
	AudienceReachabilityPercentage *float64
	AudienceAuthenticityPercentage *float64
	CommentRatePercentage          *float64
	QualityAudiencePercentage      *float64
	QualityAudienceScore           *float64
	QualityScoreGrade              *string
	QualityScoreBreakup            *string
	AudienceLocationSplit          *string
	CreatedAt                      time.Time `gorm:"default:current_timestamp"`
	UpdatedAt                      time.Time `gorm:"default:current_timestamp"`
}

func (SocialProfileAudienceInfoEntity) TableName() string {
	return "social_profile_audience_info"
}

type GroupMetricsEntity struct {
	ID                           int64 `gorm:"primaryKey;column:id"`
	GroupKey                     string
	Profiles                     *int64
	BinStart                     *string
	BinEnd                       *string
	BinHeight                    *string
	GroupAvgLikes                *float64
	GroupAvgComments             *float64
	GroupAvgCommentsRate         *float64
	GroupAvgFollowers            *float64
	GroupAvgEngagementRate       *float64
	GroupAvgReelsReach           *string
	GroupAvgLikesToCommentRatio  *float64
	GroupAvgFollowersGrowth7d    *float64
	GroupAvgFollowersGrowth30d   *float64
	GroupAvgFollowersGrowth90d   *float64
	GroupAvgAudienceReachability *float64
	GroupAvgAudienceAuthencity   *float64
	GroupAvgAudienceQuality      *float64
	GroupAvgPostCount            *float64
	GroupAvgReactionRate         *float64
	GroupAvgLikesSpread          *float64
	GroupAvgFollowersGrowth1y    *float64
	GroupAvgPostsPerWeek         *float64
	GroupAvgImageReach           *float64
	GroupAvgStoryReach           *float64
	GroupAvgVideoReach           *float64
	GroupAvgShortsReach          *float64
	Enabled                      bool
}

func (GroupMetricsEntity) TableName() string {
	return "group_metrics"
}

package domain

type Profile struct {
	Code                    *string               `json:"code"`
	Locked                  bool                  `json:"locked"`
	GccCode                 *string               `json:"gccCode"`
	PlatformCode            string                `json:"platformCode,omitempty"`
	Name                    *string               `json:"name"`
	Handle                  *string               `json:"handle,omitempty"`
	Username                *string               `json:"username,omitempty"`
	Phone                   *string               `json:"phone,omitempty"`
	Email                   *string               `json:"email,omitempty"`
	IgId                    *string               `json:"igId,omitempty"`
	Description             *string               `json:"description,omitempty"`
	Platform                string                `json:"platform,omitempty"`
	ProfileType             *string               `json:"profileType,omitempty"`
	LinkedSocials           []*SocialAccount      `json:"linkedSocials,omitempty"`
	LinkedChannelId         *string               `json:"linkedChannelId,omitempty"`
	GccLinkedChannelId      *string               `json:"gccLinkedChannelId,omitempty"`
	LinkedHandle            *string               `json:"linkedHandle,omitempty"`
	GccLinkedHandle         *string               `json:"gccLinkedHandle,omitempty"`
	Thumbnail               *string               `json:"thumbnail,omitempty"`
	CategoryRank            map[string]*string    `json:"categoryRank,omitempty"`
	EstPostPrice            map[string]*float64   `json:"estPostPrice,omitempty"`
	Categories              *[]string             `json:"categories"`
	Metrics                 ProfileMetrics        `json:"metrics,omitempty"`
	Grades                  *Grades               `json:"grades,omitempty"`
	UpdatedAt               int64                 `json:"updatedAt"`
	IsVerified              *bool                 `json:"isVerified,omitempty"`
	GroupData               map[string]*float64   `json:"groupData,omitempty"`
	SimilarProfileGroupData *GroupMetricsEntry    `json:"similarProfileGroupData,omitempty"` //make This Separta Api for Er Histogram Only
	ProfileCollectionItemId int64                 `json:"profileCollectionItemId,omitempty"`
	CampaignProfile         *CampaignProfileEntry `json:"campaignProfile,omitempty"`
	Gender                  *string               `json:"gender,omitempty"`
	LocationList            *[]string             `json:"locationList,omitempty"`
	Languages               *[]string             `json:"languages,omitempty"`
}

type ProfileMetrics struct {
	Followers                *int64   `json:"followers,omitempty"`
	Views                    *int64   `json:"views,omitempty"`
	EngagementRatePercenatge float64  `json:"engagementRatePercentage,omitempty"`
	FollowersGrowth7d        *float64 `json:"followersGrowth7d,omitempty"`
	AuthenticEngagement      *int64   `json:"authenticEngagement,omitempty"`
	CommentRatePercentage    *float64 `json:"commentRatePercentage,omitempty"`
	LikesSpreadPercentage    *float64 `json:"likesSpreadPercentage,omitempty"`
	QualityAudienceCount     *int64   `json:"qualityAudienceCount,omitempty"`
	CountryRank              *int64   `json:"countryRank,omitempty"`
	VideoViewsLast30         *int64   `json:"videoViewsLast30,omitempty"`
	StoryReach               *int64   `json:"storyReach,omitempty"`
	ImageReach               *int64   `json:"imageReach,omitempty"`
	ReelsReach               *int64   `json:"reelsReach,omitempty"`
	AvgReach                 *float64 `json:"avgReach,omitempty"`
	AvgReelsPlayCount        *float64 `json:"avgReelsPlayCount,omitempty"`
	AvgLikes                 *float64 `json:"avgLikes,omitempty"`
	AvgComments              *float64 `json:"avgComments,omitempty"`
	AvgViews                 *float64 `json:"avgViews,omitempty"`
	VideoReach               *int64   `json:"videoReach,omitempty"`
	ShortsReach              *int64   `json:"shortsReach,omitempty"`
	ReactionRate             *int64   `json:"reactionRate,omitempty"`
	CommentsRate             *int64   `json:"commentsRate,omitempty"`
	Cpm                      *int64   `json:"cpm,omitempty"`
	AvgPostsPerWeek          *int64   `json:"avgPostsPerWeek,omitempty"`
	FollowersGrowth1y        *int64   `json:"followersGrowth1y,omitempty"`
	AvgShortsViews30d        *int64   `json:"avgShortsViews30d,omitempty"`
	AvgVideoViews30d         *int64   `json:"avgVideoViews30d,omitempty"`
	LatestVideoPublishTime   int64    `json:"latestVideoPublishTime,omitempty"`
	AvgReelsPlay30d          *int64   `json:"avgReelsPlay30d,omitempty"`
	Following                *int64   `json:"following,omitempty"`
	Uploads                  *int64   `json:"uploads,omitempty"`
	Ffratio                  *float64 `json:"ffratio,omitempty"`
	Plays30d                 *int64   `json:"plays30d,omitempty"`
	Views30d                 *int64   `json:"views30d,omitempty"`
	Uploads30d               *int64   `json:"uploads30d,omitempty"`
	InstagramAccountType     *string  `json:"instagramAccountType,omitempty"`
	Impressions30d           *float64 `json:"impressions30d"`
	ImageImpressions         *float64 `json:"imageImpressions,omitempty"`
	ReelsImpressions         *float64 `json:"reelsImpressions,omitempty"`
	FollowersGrowth30d       *string  `json:"followersGrowth30d,omitempty"`
}

type Grades struct {
	ErGrade                   *string `json:"erGrade,omitempty"`
	AvgCommentsGrade          *string `json:"avgCommentsGrade,omitempty"`
	AvgLikesGrade             *string `json:"avgLikesGrade,omitempty"`
	CommentsRateGrade         *string `json:"commentsRateGrade,omitempty"`
	FollowersGrade            *string `json:"followersGrade,omitempty"`
	EngagementRateGrade       *string `json:"engagementRateGrade,omitempty"`
	LikesToCommentRatioGrade  *string `json:"likesToCommentRatioGrade,omitempty"`
	FollowersGrowth7dGrade    *string `json:"followersGrowth7dGrade,omitempty"`
	FollowersGrowth30dGrade   *string `json:"followersGrowth30dGrade,omitempty"`
	FollowersGrowth90dGrade   *string `json:"followersGrowth90dGrade,omitempty"`
	AudienceReachabilityGrade *string `json:"audienceReachabilityGrade,omitempty"`
	AudienceAuthencityGrade   *string `json:"audienceAuthencityGrade,omitempty"`
	AudienceQualityGrade      *string `json:"audienceQualityGrade,omitempty"`
	PostCountGrade            *string `json:"postCountGrade,omitempty"`
	LikesSpreadGrade          *string `json:"likesSpreadGrade,omitempty"`
	AvgPostsPerWeekGrade      *string `json:"avgPostsPerWeekGrade,omitempty"`
	FollowersGrowth1yGrade    *string `json:"followersGrowth1yGrade,omitempty"`
	ReelsReachGrade           *string `json:"reelsReachGrade,omitempty"`
	ImageReachGrade           *string `json:"imageReachGrade,omitempty"`
	StoryReachGrade           *string `json:"storyReachGrade,omitempty"`
	Views30dGrade             *string `json:"views30dGrade,omitempty"`
	ImageImpressionsGrade     *string `json:"imageImpressionsGrade,omitempty"`
	ReelsImpressionsGrade     *string `json:"reelsImpressionsGrade,omitempty"`
}

type SocialAccount struct {
	Code                          int64                     `gorm:"code" json:"code"`
	ProfileCode                   *string                   `json:"profileCode"`
	GCCProfileCode                *string                   `json:"gccProfileCode"`
	Name                          *string                   `gorm:"name" json:"name"`
	Platform                      string                    `json:"platform,omitempty"`
	Handle                        *string                   `json:"handle"`
	Username                      *string                   `json:"username,omitempty"`
	IgId                          *string                   `json:"igId,omitempty"`
	PlatformThumbnail             *string                   `gorm:"thumbnail" json:"thumbnail"`
	EngagementRatePercenatge      float64                   `gorm:"engagement_rate" json:"engagementRate,omitempty"`
	Followers                     *int64                    `gorm:"followers" json:"followers,omitempty"`
	Following                     *int64                    `json:"following,omitempty"`
	AvgLikes                      *float64                  `json:"avgLikes,omitempty"`
	AvgViews                      *float64                  `json:"avgViews,omitempty"`
	AvgComments                   *float64                  `json:"avgComments,omitempty"`
	Uploads                       *int64                    `json:"uploads,omitempty"`
	AvgVideoViews30d              *int64                    `json:"avgVideoViews30d,omitempty"`
	AvgReelsPlay30d               *int64                    `json:"avgReelsPlay30d,omitempty"`
	AvgReach                      *float64                  `json:"avgReach,omitempty"` //TDB in YA
	StoryReach                    *int64                    `json:"storyReach,omitempty"`
	ImageReach                    *int64                    `json:"imageReach,omitempty"`
	ReelsReach                    *int64                    `json:"reelsReach,omitempty"`
	Views                         *int64                    `json:"views,omitempty"`
	IsVerified                    *bool                     `json:"isVerified,omitempty"`
	Source                        string                    `json:"source,omitempty"`
	AccountType                   *string                   `json:"accountType,omitempty"`
	ProfileLink                   *string                   `json:"profileLink,omitempty"`
	Category                      *string                   `json:"category,omitempty"`
	CampaignProfileId             *int64                    `json:"id,omitempty"`
	ImageEngagementRatePercenatge float64                   `json:"imageEngagementRate,omitempty"`
	ReelsEngagementRatePercenatge float64                   `json:"reelsEngagementRate,omitempty"`
	ImageAvgLikes                 *float64                  `json:"imageAvgLikes,omitempty"`
	ReelsAvgLikes                 *float64                  `json:"reelsAvgLikes,omitempty"`
	ImageAvgComments              *float64                  `json:"imageAvgComments,omitempty"`
	ReelsAvgComments              *float64                  `json:"reelsAvgComments,omitempty"`
	RecentPosts                   []SocialProfilePostsEntry `json:"recentPosts,omitempty"`
	UpdatedAt                     string                    `json:"updatedAt,omitempty"`
	AvgReelsPlayCount             *float64                  `json:"avgReelsPlayCount,omitempty"`
	AudienceCityAvailable         *bool                     `json:"audienceCityAvailable,omitempty"`
	AudienceGenderAvailable       *bool                     `json:"audienceGenderAvailable,omitempty"`
	IsPrivate                     *bool                     `json:"isPrivate,omitempty"`
	FollowersGrowth7d             *float64                  `json:"followersGrowth7d,omitempty"`
	FollowersGrowth30d            *string                   `json:"followersGrowth30d,omitempty"`
}

type GroupMetricsEntry struct {
	GroupKey                     string        `json:"groupKey,omitempty"`
	Profiles                     *int64        `json:"profiles,omitempty"`
	GroupAvgLikes                *float64      `json:"groupAvgLikes,omitempty"`
	GroupAvgComments             *float64      `json:"groupAvgComments,omitempty"`
	GroupAvgCommentsRate         *float64      `json:"groupAvgCommentsRate,omitempty"`
	GroupAvgFollowers            *float64      `json:"groupAvgFollowers,omitempty"`
	GroupAvgEngagementRate       *float64      `json:"groupAvgEngagementRate,omitempty"`
	GroupAvgLikesToCommentRatio  *float64      `json:"groupAvgLikesToCommentRatio,omitempty"`
	GroupAvgFollowersGrowth7d    *float64      `json:"groupAvgFollowersGrowth7d,omitempty"`
	GroupAvgFollowersGrowth30d   *float64      `json:"groupAvgFollowersGrowth30d,omitempty"`
	GroupAvgFollowersGrowth90d   *float64      `json:"groupAvgFollowersGrowth90d,omitempty"`
	GroupAvgAudienceReachability *float64      `json:"groupAvgAudienceReachability,omitempty"`
	GroupAvgAudienceAuthencity   *float64      `json:"groupAvgAudienceAuthencity,omitempty"`
	GroupAvgAudienceQuality      *float64      `json:"groupAvgAudienceQuality,omitempty"`
	GroupAvgPostCount            *float64      `json:"groupAvgPostCount,omitempty"`
	GroupAvgReactionRate         *float64      `json:"groupAvgReactionRate,omitempty"`
	GroupAvgLikesSpread          *float64      `json:"groupAvgLikesSpread,omitempty"`
	GroupAvgFollowersGrowth1y    *float64      `json:"groupAvgFollowersGrowth1y,omitempty"`
	GroupAvgPostsPerWeek         *float64      `json:"groupAvgPostsPerWeek,omitempty"`
	GroupAvgReelsReach           *string       `json:"groupAvgReelsReach,omitempty"`
	GroupAvgImageReach           *float64      `json:"groupAvgImageReach,omitempty"`
	GroupAvgStoryReach           *float64      `json:"groupAvgStoryReach,omitempty"`
	GroupAvgVideoReach           *float64      `json:"groupAvgVideoReach,omitempty"`
	GroupAvgShortsReach          *float64      `json:"groupAvgShortsReach,omitempty"`
	ErGraph                      []ERHistogram `json:"erHistogram,omitempty"`
}

type ERHistogram struct {
	Start float64 `json:"start"`
	End   float64 `json:"end"`
	Value float64 `json:"value"`
}

type AdminDetails struct {
	Name                *string           `json:"name,omitempty"`
	Email               *string           `json:"email,omitempty"`
	Phone               *string           `json:"phone,omitempty"`
	SecondaryPhone      *string           `json:"secondaryPhone,omitempty"`
	Gender              *string           `json:"gender,omitempty"`
	Languages           *[]string         `json:"languages,omitempty"`
	City                *string           `json:"city,omitempty"`
	State               *string           `json:"state,omitempty"`
	Country             *string           `json:"country,omitempty"`
	Dob                 *string           `json:"dob,omitempty"`
	Bio                 *string           `json:"bio,omitempty"`
	CampaignCategoryIds *[]int64          `json:"campaignCategoryIds,omitempty"`
	CreatorPrograms     *[]CreatorProgram `json:"creatorPrograms,omitempty"`
	IsBlacklisted       *bool             `json:"isBlacklisted,omitempty"`
	BlacklistedBy       *string           `json:"blacklistedBy,omitempty"`
	BlacklistedReason   *string           `json:"blacklistedReason,omitempty"`
	Location            *[]string         `json:"location,omitempty"`
	WhatsappOptIn       *bool             `json:"whatsappOptIn"`
	CountryCode         *string           `json:"countryCode,omitempty"`
	CreatorCohorts      *[]CreatorCohort  `json:"creatorCohorts,omitempty"`
	// ContentCategories     *[]ContentCategory `json:"contentCategories,omitempty"`
}
type CreatorProgram struct {
	Level string `json:"level,omitempty"`
	Tag   string `json:"tag,omitempty"`
}

type ContentCategory struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type CreatorCohort struct {
	Level string `json:"level,omitempty"`
	Tag   string `json:"tag,omitempty"`
}

type UserDetails struct {
	Email                       *string   `json:"email,omitempty"`
	Phone                       *string   `json:"phone,omitempty"`
	SecondaryPhone              *string   `json:"secondaryPhone,omitempty"`
	Name                        *string   `json:"name,omitempty"`
	Dob                         *string   `json:"dob,omitempty"`
	Gender                      *string   `json:"gender,omitempty"`
	WhatsappOptIn               *bool     `json:"whatsappOptIn"`
	CampaignCategoryIds         *[]int64  `json:"campaignCategoryIds,omitempty"`
	Languages                   *[]string `json:"languages,omitempty"`
	Location                    *[]string `json:"location,omitempty"`
	NotificationToken           *string   `json:"notificationToken,omitempty"`
	Bio                         *string   `json:"bio,omitempty"`
	MemberId                    *string   `json:"memberId,omitempty"`
	ReferenceCode               *string   `json:"referenceCode,omitempty"`
	EKYCPending                 *bool     `json:"ekycPending"`
	WebEngageUserId             *string   `json:"webengageUserId,omitempty"`
	InstantGratificationInvited *bool     `json:"instantGratificationInvited"`
	AmazonStoreLinkVerified     *bool     `json:"amazonStoreLinkVerified"`
}
type CampaignProfileEntry struct {
	Id                int64          `json:"id"`
	Platform          string         `json:"platform"`
	PlatformAccountId int64          `json:"platformAccountId"`
	Handle            *string        `json:"handle,omitempty"`
	OnGCC             *bool          `json:"onGcc"`
	OnGCCApp          *bool          `json:"onGccApp"`
	UserDetails       *UserDetails   `json:"userDetails,omitempty"`
	AdminDetails      *AdminDetails  `json:"adminDetails,omitempty"`
	UpdatedBy         *string        `json:"updatedBy,omitempty"`
	AccountId         *int64         `json:"accountId"`
	SocialDetails     *SocialAccount `json:"socialDetails,omitempty"`
	HasEmail          *bool          `json:"hasEmail"`
	HasPhone          *bool          `json:"hasPhone"`
}

type CampaignProfileInput struct {
	Id                int64            `json:"id"`
	Platform          string           `json:"platform"`
	PlatformAccountId int64            `json:"platformAccountId"`
	Handle            *string          `json:"handle,omitempty"`
	OnGCC             *bool            `json:"onGcc"`
	OnGCCApp          *bool            `json:"onGccApp"`
	UserDetails       *UserDetails     `json:"userDetails,omitempty"`
	AdminDetails      *AdminDetails    `json:"adminDetails,omitempty"`
	Socials           []*SocialAccount `json:"socials,omitempty"`
	UpdatedBy         *string          `json:"updatedBy,omitempty"`
	AccountId         *int64           `json:"accountId"`
	SocialDetails     *SocialAccount   `json:"socialDetails,omitempty"`
	HasEmail          *bool            `json:"hasEmail"`
	HasPhone          *bool            `json:"hasPhone"`
}

type CreatorStructs struct {
	CreatorProgram      *[]string `json:"creatorPrograms,omitempty"`
	CreatorCohorts      *[]string `json:"creatorCohorts,omitempty"`
	CampaignCategoryIds *[]string `json:"campaignCategoryIds,omitempty"`
}
type InstagramImpressions struct {
	ImagePosts *float64 `json:"image_posts,omitempty"`
	ReelPosts  *float64 `json:"reel_posts,omitempty"`
}

type SocialProfilePostsEntry struct {
	Platform       string  `json:"platform"`
	Handle         string  `json:"handle"`
	PostId         string  `json:"postId"`
	PostLink       string  `json:"postLink"`
	PostType       string  `json:"postType"`
	Thumbnail      string  `json:"thumbnail"`
	LikesCount     int64   `json:"likesCount"`
	CommentsCount  int64   `json:"commentsCount"`
	ViewsCount     int64   `json:"viewsCount"`
	EngagementRate float64 `json:"engagementRate"`
	PublishedAt    int64   `json:"publishedAt"`
	Caption        string  `json:"caption"`
	PlayCount      int64   `json:"playCount"`
}

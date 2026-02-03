package entry

type Client struct {
	Id      string `json:"id,omitempty"`
	AppName string `json:"appName,omitempty"`
	AppType string `json:"appType,omitempty"`
	Enabled bool   `json:"enabled,omitempty"`
}

type Gender string

const (
	MALE   Gender = "MALE"
	FEMALE Gender = "FEMALE"
)

type User struct {
	Id           int           `json:"id,omitempty"`
	Status       string        `json:"status,omitempty"`
	Uidx         string        `json:"uidx,omitempty"`
	Dob          *int64        `json:"dob,omitempty"`
	Emails       []UserEmail   `json:"emails,omitempty"`
	Phones       []UserPhone   `json:"phones,omitempty"`
	Gender       *Gender       `json:"gender,omitempty"`
	ReferralCode *string       `json:"referralCode,omitempty"`
	SocialLinks  []SocialLink  `json:"socialLinks,omitempty"`
	KycDocuments []KYCDocument `json:"kycDocuments,omitempty"`
}

type UserClientAccountStatus string

const (
	UserClientAccountStatusOnboarding UserClientAccountStatus = "ONBOARDING"
	UserClientAccountStatusActive     UserClientAccountStatus = "ACTIVE"
	UserClientAccountStatusInactive   UserClientAccountStatus = "INACTIVE"
	UserClientAccountStatusBlocked    UserClientAccountStatus = "BLOCKED"
)

type UserClientAccountOnboardingStep string

const (
	UserClientAccountOnboardingStepStart           UserClientAccountOnboardingStep = "START"
	UserClientAccountOnboardingStepReview_pending  UserClientAccountOnboardingStep = "REVIEW_PENDING"
	UserClientAccountOnboardingStepIncomplete_info UserClientAccountOnboardingStep = "INCOMPLETE_INFO"
	UserClientAccountOnboardingStepRejected        UserClientAccountOnboardingStep = "REJECTED"
	UserClientAccountOnboardingStepAccepted        UserClientAccountOnboardingStep = "ACCEPTED"
)

type UserClientAccount struct {
	Id             int                             `json:"id,omitempty"`
	ClientId       string                          `json:"clientId,omitempty"`
	ClientAppType  string                          `json:"clientAppType,omitempty"`
	UserId         int                             `json:"userId,omitempty"`
	Status         UserClientAccountStatus         `json:"status,omitempty"`
	OnboardingStep UserClientAccountOnboardingStep `json:"onboardingStep,omitempty"`
	Remark         *string                         `json:"remark,omitempty"`
	Permissions    []string                        `json:"permissions,omitempty"`
	Source         *string                         `json:"source,omitempty"`

	Name            *string `json:"name,omitempty"`
	Bio             *string `json:"bio,omitempty"`
	ProfileImageId  *int    `json:"profileImageId,omitempty"`
	CoverImageId    *int    `json:"coverImageId,omitempty"`
	PublicProfileId *string `json:"publicProfileId,omitempty"`
	StoreLink       *string `json:"storeLink,omitempty"`

	ReferredViaCode     *string `json:"referredViaCode,omitempty"`
	ReferredViaClientId *string `json:"referredViaClientId,omitempty"`

	LastUsedAppVersion     *int    `json:"lastUsedAppVersion,omitempty"`
	LastUsedAppVersionName *string `json:"lastUsedAppVersionName,omitempty"`
	LastLoginTime          *int64  `json:"lastLoginTime,omitempty"`

	User *User `json:"user,omitempty"`

	HostProfile     *HostProfile     `json:"hostProfile,omitempty"`
	CustomerProfile *CustomerProfile `json:"customerProfile,omitempty"`

	ProfileImage *AssetInfo `json:"profileImage,omitempty"`
	CoverImage   *AssetInfo `json:"coverImage,omitempty"`

	WebengageUserId *string `json:"webengageUserId,omitempty"`
}

type UserEmail struct {
	Email    string `json:"email,omitempty`
	UserID   int    `json:"userId,omitempty"`
	Verified bool   `json:"verified,omitempty"`
	Enabled  bool   `json:"enabled,omitempty"`
}

type UserPhone struct {
	CountryCode string `json:"countryCode,omitempty"`
	Phone       string `json:"phone,omitempty"`
	UserId      int    `json:"userId,omitempty"`
	Verified    bool   `json:"verified,omitempty"`
	Enabled     bool   `json:"enabled,omitempty"`
}

type SocialLink struct {
	Id         int    `json:"id,omitempty"`
	SocialId   string `json:"socialId,omitempty"`
	ProfileUrl string `json:"profileUrl,omitempty"`
	Type       string `json:"type,omitempty"`
	Verified   bool   `json:"verified,omitempty"`
	Enabled    bool   `json:"enabled,omitempty"`
	ClientId   string `json:"clientId,omitempty"`
	IconUrl    string `json:"iconUrl,omitempty"`
}

type KYCDocument struct {
	Id             int        `json:"id,omitempty"`
	DocumentNumber string     `json:"documentNumber,omitempty"`
	Type           string     `json:"type,omitempty"`
	FrontImage     *AssetInfo `json:"frontImage,omitempty"`
	BackImage      *AssetInfo `json:"backImage,omitempty"`
	Verified       bool       `json:"verified,omitempty"`
	Enabled        bool       `json:"enabled,omitempty"`
}

type SocialNetworkType string

const (
	GOOGLE   SocialNetworkType = "GOOGLE"
	FB       SocialNetworkType = "FB"
	TWITTER  SocialNetworkType = "TWITTER"
	IG       SocialNetworkType = "IG"
	TIKTOK   SocialNetworkType = "TIKTOK"
	YOUTUBE  SocialNetworkType = "YOUTUBE"
	WHATSAPP SocialNetworkType = "WHATSAPP"
	OTHERS   SocialNetworkType = "OTHERS"
	_NONE   SocialNetworkType = "NONE"
	SNAPCHAT    SocialNetworkType = "SNAPCHAT"
	PINTEREST    SocialNetworkType = "PINTEREST"
	LINKEDIN    SocialNetworkType = "LINKEDIN"
	DISCORD    SocialNetworkType = "DISCORD"
	TWITCH    SocialNetworkType = "TWITCH"
	PATREON    SocialNetworkType = "PATREON"
)

type AttributeKey string

const (
	FOLLOWING_COUNT                    AttributeKey = "FOLLOWING_COUNT"
	FOLLOWERS_COUNT                    AttributeKey = "FOLLOWERS_COUNT"
	SCHEDULED_VIDEO_COUNT              AttributeKey = "SCHEDULED_VIDEO_COUNT"
	CURRENT_LIVE_STREAM                AttributeKey = "CURRENT_LIVE_STREAM"
	FREE_ORDERS_COUNT                  AttributeKey = "FREE_ORDERS_COUNT"
	TOTAL_ORDERS_COUNT                 AttributeKey = "TOTAL_ORDERS_COUNT"
	LIKED_STREAMS_COUNT                AttributeKey = "LIKED_STREAMS_COUNT"
	VIEWED_STREAMS_COUNT               AttributeKey = "VIEWED_STREAMS_COUNT"
	DAILY_CHECKIN_TASKS_COMPLETED      AttributeKey = "DAILY_CHECKIN_TASKS_COMPLETED"
	DAILY_CHECKIN_DONE_TODAY           AttributeKey = "DAILY_CHECKIN_DONE_TODAY"
	ENABLED_FOR_GROUP                  AttributeKey = "ENABLED_FOR_GROUP"
	CART_ITEM_COUNT                    AttributeKey = "CART_ITEM_COUNT"
	SPOKEN_LANGUAGES                   AttributeKey = "SPOKEN_LANGUAGES"
	INTERESTED_CATEGORIES              AttributeKey = "INTERESTED_CATEGORIES"
	AUDITION_VIDEO_ASSET_IDS           AttributeKey = "AUDITION_VIDEO_ASSET_IDS"
	AUDITION_VIDEO_ASSETS              AttributeKey = "AUDITION_VIDEO_ASSETS"
	TOTAL_LIVE_COUNT                   AttributeKey = "TOTAL_LIVE_COUNT"
	LIVE_STREAM_ENABLED                AttributeKey = "LIVE_STREAM_ENABLED"
	STATIC_STREAM_ENABLED              AttributeKey = "STATIC_STREAM_ENABLED"
	STREAM_SHARE_MSG_TEMPLATE          AttributeKey = "STREAM_SHARE_MSG_TEMPLATE"
	EXPERT_TOPICS                      AttributeKey = "EXPERT_TOPICS"
	PROFILE_INTRO_VIDEO_ID             AttributeKey = "PROFILE_INTRO_VIDEO_ID"
	PROFILE_PROMO_IMAGE_ID             AttributeKey = "PROFILE_PROMO_IMAGE_ID"
	CUSTOMER_APP_PROFILE_SHARE_LINK_ID AttributeKey = "CUSTOMER_APP_PROFILE_SHARE_LINK_ID"
	NOTIFICATION_UNREAD_COUNT          AttributeKey = "NOTIFICATION_UNREAD_COUNT"
	PROFILE_INTRO_VIDEO                AttributeKey = "PROFILE_INTRO_VIDEO"
	PROFILE_PROMO_IMAGE                AttributeKey = "PROFILE_PROMO_IMAGE"
	CUSTOMER_APP_PROFILE_SHARE_LINK    AttributeKey = "CUSTOMER_APP_PROFILE_SHARE_LINK"
	CUSTOMER_APP_PROFILE_SHARE_MESSAGE AttributeKey = "CUSTOMER_APP_PROFILE_SHARE_MESSAGE"
)

type UserDeviceAttribute struct {
	Id             int          `json:"id,omitempty"`
	Key            AttributeKey `json:"key,omitempty"`
	Value          interface{}  `json:"value,omitempty"`
	UserId         int          `json:"userId,omitempty"`
	DeviceId       *string      `json:"deviceId,omitempty"`
	TargetClientId string       `json:"targetClientId,omitempty"`
}

type HostProfile struct {
	Id                  int `json:"id,omitempty"`
	UserClientAccountId int `json:"userClientAccountId,omitempty"`

	FollowerCount           int                       `json:"followerCount,omitempty"`
	ScheduledVideoCount     int                       `json:"scheduledVideoCount,omitempty"`
	TotalLiveCount          int                       `json:"totalLiveCount,omitempty"`
	CartItemCount           int                       `json:"cartItemCount,omitempty"`
	WishlistItemCount       int                       `json:"wishlistItemCount,omitempty"`
	UnreadNotificationCount int                       `json:"unreadNotificationCount,omitempty"`
	CurrentLiveStreamId     *int                      `json:"currentLiveStreamId,omitempty"`
	SpokenLanguages         []string                  `json:"spokenLanguages,omitempty"`
	Categories              []int                     `json:"categories,omitempty"`
	AuditionVideoIds        []int                     `json:"auditionVideoIds,omitempty"`
	ExpertTopics            []string                  `json:"expertTopics,omitempty"`
	IntroVideoId            int                       `json:"introVideoId,omitempty"`
	PromoImageId            int                       `json:"promoImageId,omitempty"`
	CustomerAppShareLinkIds map[SocialNetworkType]int `json:"customerAppShareLinkIds,omitempty"`
	HostAppShareLinkIds     map[SocialNetworkType]int `json:"hostAppShareLinkIds,omitempty"`
	StreamShareMsgTemplate  string                    `json:"streamShareMsgTemplate,omitempty"`

	IntroVideo            *AssetInfo                   `json:"introVideo,omitempty"`
	PromoImage            *AssetInfo                   `json:"promoImage,omitempty"`
	AuditionVideos        []AssetInfo                  `json:"auditionVideos,omitempty"`
	CustomerAppShareLinks map[SocialNetworkType]string `json:"customerAppShareLinks,omitempty"`
	HostAppShareLinks     map[SocialNetworkType]string `json:"hostAppShareLinks,omitempty"`
	ProfileShareMessage   *string                      `json:"profileShareMessage,omitempty"`
	ShareIcons            []ShareIcon                  `json:"shareIcons,omitempty"`
	ProfileStatusInfo     HostProfileStatusInfo        `json:"profileStatusInfo,omitempty"`

	ProductOfferingTypes  []string                     `json:"productOfferingTypes,omitempty"`
	StoreTheme            string                       `json:"storeTheme,omitempty"`
	WhatsappChatNumber    *string                       `json:"whatsappChatNumber,omitempty"`
	WhatsappChatIntroMessage    *string                 `json:"whatsappChatIntroMessage,omitempty"`
}

type ShareIcon struct {
	SocialNetworkType SocialNetworkType `json:"socialNetworkType,omitempty"`
	DisplayName       *string           `json:"displayName,omitempty"`
	IconImage         *string           `json:"iconImage,omitempty"`
	Link              *string           `json:"link,omitempty"`
}

type HostProfileStatusInfo struct {
	Title         *string `json:"title,omitempty"`
	Description   *string `json:"description,omitempty"`
	CtaText       *string `json:"ctaText,omitempty"`
	CtaActionLink *string `json:"ctaActionLink,omitempty"`

	ProfileCompletionPercent int `json:"profileCompletionPercent,omitempty"`

	ProfileCompletionStatus *string `json:"profileCompletionStatus,omitempty"`
}

type CustomerProfile struct {
	UnreadNotificationCount int  `json:"unreadNotificationCount,omitempty"`
	LikeStreamCount         int  `json:"likeStreamCount,omitempty"`
	ViewStreamCount         int  `json:"viewStreamCount,omitempty"`
	FollowingHostCount      int  `json:"followingHostCount,omitempty"`
	CartItemCount           int  `json:"cartItemCount,omitempty"`
	WishlistItemCount       int  `json:"wishlistItemCount,omitempty"`
	Credit                  int  `json:"credit,omitempty"`
	EnabledForGroup         bool `json:"enabledForGroup,omitempty"`

	PreferredPaymentMode *string `json:"preferredPaymentMode,omitempty"`
	FreeOrdersCount      int     `json:"freeOrdersCount,omitempty"`
	TotalOrdersCount     int     `json:"totalOrdersCount,omitempty"`
}
type OTPRequestOperation string

const (
	OTPRequestOperationLink   OTPRequestOperation = "LINK"
	OTPRequestOperationVerify OTPRequestOperation = "VERIFY"
	OTPRequestOperationLogin  OTPRequestOperation = "LOGIN"
	OTPRequestOperationSignup OTPRequestOperation = "SIGNUP"
)

type TrueCallerRequestStatus string

const (
	TrueCallerRequestStatusOpen          TrueCallerRequestStatus = "OPEN"
	TrueCallerRequestStatusTokenReceived TrueCallerRequestStatus = "TOKEN_RECEIVED"
	TrueCallerRequestStatusVerified      TrueCallerRequestStatus = "VERIFIED"
	TrueCallerRequestStatusUserRejected  TrueCallerRequestStatus = "USER_REJECTED"
)

type TrueCallerRequest struct {
	Id         string                  `json:"id"`
	Status     TrueCallerRequestStatus `json:"status,omitempty"`
	ExpiryTime int                     `json:"expiryTime,omitempty"`
}

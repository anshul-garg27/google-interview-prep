package entry

type SocialNetworkType string

const (
	GOOGLE    SocialNetworkType = "GOOGLE"
	FB        SocialNetworkType = "FB"
	TWITTER   SocialNetworkType = "TWITTER"
	TIKTOK    SocialNetworkType = "TIKTOK"
	YOUTUBE   SocialNetworkType = "YOUTUBE"
	INSTAGRAM SocialNetworkType = "INSTAGRAM"
)

type LeadboardStatus string

const (
	ENDED SocialNetworkType = "ENDED"
	LIVE  SocialNetworkType = "LIVE"
)

type Gender string

const (
	MALE   Gender = "MALE"
	FEMALE Gender = "FEMALE"
	OTHERS Gender = "OTHERS"
)

type Client struct {
	Id      string `json:"id,omitempty"`
	AppName string `json:"appName,omitempty"`
	AppType string `json:"appType,omitempty"`
	Enabled bool   `json:"enabled,omitempty"`
}

type User struct {
	Id           int         `json:"id,omitempty"`
	Uidx         string      `json:"uidx,omitempty"`
	Dob          *int64      `json:"dob,omitempty"`
	Emails       []UserEmail `json:"emails,omitempty"`
	Phones       []UserPhone `json:"phones,omitempty"`
	Gender       *Gender     `json:"gender,omitempty"`
	ReferralCode *string     `json:"referralCode,omitempty"`
}

type UserEmail struct {
	Email    string `json:"email,omitempty"`
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

type UserClientAccountStatus string

const (
	UserClientAccountStatusOnboarding UserClientAccountStatus = "ONBOARDING"
	UserClientAccountStatusActive     UserClientAccountStatus = "ACTIVE"
	UserClientAccountStatusInactive   UserClientAccountStatus = "INACTIVE"
	UserClientAccountStatusBlocked    UserClientAccountStatus = "BLOCKED"
)

type Portal string

const (
	PortalGraphy Portal = "GRAPHY"
)

type InfluencerTag string

const (
	InfluencerTagCeleb InfluencerTag = "CELEB"
	InfluencerTagMacro InfluencerTag = "MACRO_INFLUENCER"
	InfluencerTagMicro InfluencerTag = "MICRO_INFLUENCER"
	InfluencerTagNano  InfluencerTag = "NANO_INFLUENCER"
	InfluencerTagPage  InfluencerTag = "CURATION_PAGE"
	InfluencerTagBrand InfluencerTag = "BRAND"
	InfluencerTagSpam  InfluencerTag = "SPAM"
)

type UserClientAccount struct {
	Id            int                     `json:"id,omitempty"`
	ClientId      string                  `json:"clientId,omitempty"`
	ClientAppType string                  `json:"clientAppType,omitempty"`
	UserId        int                     `json:"userId,omitempty"`
	Status        UserClientAccountStatus `json:"status,omitempty"`

	Name            *string    `json:"name,omitempty"`
	Bio             *string    `json:"bio,omitempty"`
	Email           *string    `json:"email,omitempty"`
	ProfileImageId  *int       `json:"profileImageId,omitempty"`
	ProfileImage    *AssetInfo `json:"profileImage,omitempty"`
	CoverImageId    *int       `json:"coverImageId,omitempty"`
	CoverImage      *AssetInfo `json:"coverImage,omitempty"`
	PublicProfileId *string    `json:"publicProfileId,omitempty"`
	StoreLink       *string    `json:"storeLink,omitempty"`

	ReferredViaCode     *string `json:"referredViaCode,omitempty"`
	ReferredViaClientId *string `json:"referredViaClientId,omitempty"`

	LastUsedAppVersion     *int    `json:"lastUsedAppVersion,omitempty"`
	LastUsedAppVersionName *string `json:"lastUsedAppVersionName,omitempty"`
	LastLoginTime          *int64  `json:"lastLoginTime,omitempty"`
	FirstLoginTime         *int64  `json:"firstLoginTime,omitempty"`

	Permissions []string `json:"permissions,omitempty"`

	User            *User               `json:"user,omitempty"`
	HostProfile     *HostProfile        `json:"hostProfile,omitempty"`
	CustomerProfile *CustomerProfile    `json:"customerProfile,omitempty"`
	PartnerProfile  *PartnerUserProfile `json:"partnerProfile,omitempty"`

	SocialAccounts []SocialAccount `json:"socialAccounts,omitempty"`

	FirebaseToken *string `json:"firebaseToken,omitempty"`
	PortalId      *Portal `json:"portalId,omitempty"`

	PhoneNumber    *string `json:"phoneNumber,omitempty"`
	WhatsappNumber *string `json:"whatsappNumber,omitempty"`
	WhatsappOptin  *bool   `json:"whatsappOptin,omitempty"`

	KycData *UserKycData `json:"kycData,omitempty"`
}

type SocialAccount struct {
	Platform           SocialNetworkType     `json:"platform,omitempty"`
	Handle             string                `json:"handle,omitempty"`
	ExistingHandle     *bool                 `json:"existingHandle,omitempty"`
	Message            *string               `json:"message,omitempty"`
	GraphAccessToken   *string               `json:"graphAccessToken,omitempty"`
	AccessToken        *string               `json:"accessToken,omitempty"`
	AuthCode           *string               `json:"authCode,omitempty"`
	SocialUserId       *string               `json:"socialUserId,omitempty"`
	AccessTokenExpired *bool                 `json:"accessTokenExpired,omitempty"`
	Verified           bool                  `json:"verified,omitempty"`
	Name               *string               `json:"name,omitempty"`
	IgInfo             *InstantGratification `json:"igInfo,omitempty"`
	Notice             string                `json:"notice,omitempty"`
	Email              *string               `json:"email,omitempty"`
	Metrics            *SocialAccountMetrics `json:"metrics,omitempty"`
	ProfileImage       *string               `json:"profileImage,omitempty"`
}

type SocialAccountMetrics struct {
	Followers *int `json:"followers,omitempty"`
}

type HostProfile struct {
	Id                  int `json:"id,omitempty"`
	UserClientAccountId int `json:"userClientAccountId,omitempty"`

	MyGlammMemberId      *string `json:"myGlammMemberId,omitempty"`
	MyglammReferenceCode *string `json:"myglammReferenceCode,omitempty"`
	PlixxoId             *string `json:"plixxoId,omitempty"`
	IsPlixxoUser         bool    `json:"isPlixxoUser,omitempty"`
	CampaignDataMigrated bool    `json:"campaignDataMigrated,omitempty"`
	IsGccMigrated        bool    `json:"isGccMigrated,omitempty"`
	WebengageUserId      *string `json:"webengageUserId,omitempty"`

	ShopOnboarded *bool   `json:"shopOnboarded,omitempty"`
	ShopBio       *string `json:"shopBio,omitempty"`

	PlusMember                *bool         `json:"plusMember,omitempty"`
	BarterAllowed             *bool         `json:"barterAllowed,omitempty"`
	ShortlistedInPaidCampaign *bool         `json:"shortlistedInPaidCampaign,omitempty"`
	Label                     InfluencerTag `json:"label,omitempty"`

	Location             *string  `json:"location,omitempty"`
	Languages            []string `json:"languages,omitempty"`
	NotificationDeviceId *string  `json:"notificationDeviceId,omitempty"`

	ProfileStatusInfo *HostProfileStatusInfo `json:"profileStatusInfo,omitempty"`
}

type InstantGratification struct {
	CampaignIds []int `json:"campaignIds,omitempty"`
	IsInvited   bool  `json:"isInvited,omitempty"`
}

type HostProfileStatusInfo struct {
	Title         *string `json:"title,omitempty"`
	Description   *string `json:"description,omitempty"`
	CtaText       *string `json:"ctaText,omitempty"`
	CtaActionLink *string `json:"ctaActionLink,omitempty"`

	ProfileCompletionPercent    int `json:"profileCompletionPercent,omitempty"`
	AppProfileCompletionPercent int `json:"appProfileCompletionPercent,omitempty"`

	ProfileCompletionStatus string `json:"profileCompletionStatus,omitempty"`

	EditProfileToast *string `json:"editProfileToast,omitempty"`
}

type CustomerProfile struct {
	Id                      int `json:"id,omitempty"`
	UnreadNotificationCount int `json:"unreadNotificationCount,omitempty"`
	LikeStreamCount         int `json:"likeStreamCount,omitempty"`
	ViewStreamCount         int `json:"viewStreamCount,omitempty"`
	FollowingHostCount      int `json:"followingHostCount,omitempty"`
	CartItemCount           int `json:"cartItemCount,omitempty"`
	WishlistItemCount       int `json:"wishlistItemCount,omitempty"`
	Credit                  int `json:"credit,omitempty"`

	PreferredPaymentMode *string `json:"preferredPaymentMode,omitempty"`
	FreeOrdersCount      int     `json:"freeOrdersCount,omitempty"`
	TotalOrdersCount     int     `json:"totalOrdersCount,omitempty"`
}

type PartnerUserProfile struct {
	PartnerId           int64       `json:"partnerId,omitempty"`
	Phone               *string     `json:"phone,omitempty"`
	PlanType            string      `json:"planType,omitempty"`
	VisitedDemo         bool        `json:"visitedDemo,omitempty"`
	AssignedCampaignIds []int       `json:"assignedCampaignIds,omitempty"`
	WinklBrand          *WinklBrand `json:"brand,omitempty"`
}

type WinklBrand struct {
	Id                 int     `json:"id,omitempty"`
	CompanyName        string  `json:"companyName,omitempty"`
	IsOnboarded        bool    `json:"isOnboarded,omitempty"`
	LogoUrl            *string `json:"logoUrl,omitempty"`
	WhitelistedToolIds []int   `json:"whitelistedToolIds,omitempty"`
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

type OTPRequestEntry struct {
	Id          string `json:"id"`
	IdentityKey string `json:"identityKey"`
	Operation   string `json:"operation"`
	TargetMode  string `json:"targetMode"`
	Otp         string `json:"otp"`
}

type UserKycData struct {
	Aadhaar *KycDocument `json:"aadhaar,omitempty"`
	Pan     *KycDocument `json:"pan,omitempty"`
	Bank    *KycDocument `json:"bank,omitempty"`
	Gst     *KycDocument `json:"gstIn,omitempty"`

	IsComplete bool `json:"complete"`
	IsVerified bool `json:"verified"`
}

type KycDocument struct {
	DocumentId               string      `json:"documentIdentifier"`
	ManualVerificationOnHold bool        `json:"manualVerificationOnhold"`
	OnholdReason             *string     `json:"onholdReason,omitempty"`
	BankDetails              *BankDetail `json:"bankDetail,omitempty"`
	PanDetails               *PanDetail  `json:"panDetail,omitempty"`
}

type PanDetail struct {
	Number *string `json:"panCardNumber,omitempty"`
	Dob    *string `json:"dob,omitempty"`
	Name   *string `json:"panName,omitempty"`
}

type BankDetail struct {
	AccountNumber     *string `json:"accountNumber,omitempty"`
	BankName          *string `json:"bankName,omitempty"`
	AccountHolderName *string `json:"accountHolderName,omitempty"`
	Ifsc              *string `json:"ifsc,omitempty"`
}

type MemberInvite struct {
	Id                     int                `json:"id,omitempty"`
	Code                   string             `json:"code,omitempty"`
	Status                 MemberInviteStatus `json:"status,omitempty"`
	MemberEmail            string             `json:"memberEmail,omitempty"`
	PartnerId              int                `json:"partnerId,omitempty"`
	PartnerClientId        string             `json:"partnerClientId,omitempty"`
	Permissions            []string           `json:"permissions,omitempty"`
	ExpiresOn              int64              `json:"expiresOn,omitempty"`
	IsAdmin                bool               `json:"isAdmin,omitempty"`
	AcceptedOn             *int64             `json:"acceptedOn,omitempty"`
	InvitedByUserAccountId int                `json:"invitedByUserAccountId,omitempty"`
}

type MemberInviteStatus string

const (
	MemberInviteStatusOpen     MemberInviteStatus = "OPEN"
	MemberInviteStatusAccepted MemberInviteStatus = "ACCEPTED"
	MemberInviteStatusCanceled MemberInviteStatus = "CANCELLED"
)

type ResetPasswordRequest struct {
	Id        int                        `json:"id,omitempty"`
	Code      string                     `json:"code,omitempty"`
	Status    ResetPasswordRequestStatus `json:"status,omitempty"`
	Email     string                     `json:"email,omitempty"`
	ClientId  string                     `json:"clientId,omitempty"`
	ExpiresOn int64                      `json:"expiresOn,omitempty"`
}

type ResetPasswordRequestStatus string

const (
	ResetPasswordRequestStatusOpen     ResetPasswordRequestStatus = "OPEN"
	ResetPasswordRequestStatusComplete ResetPasswordRequestStatus = "COMPLETE"
)

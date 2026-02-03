package winkl

import (
	"fmt"
	"io"
	"strconv"
)

type WinklUpdateInfluencerInput struct {
	// Id                *string `json:"id,omitempty"`
	IsBlacklisted     *bool   `json:"isBlacklisted,omitempty"`
	BlacklistedBy     *string `json:"blacklistedBy,omitempty"`
	BlacklistedReason *string `json:"blacklistedReason,omitempty"`
	Name              *string `json:"name,omitempty"`
	Email             *string `json:"email,omitempty"`
	PhoneNumber       *string `json:"phoneNumber,omitempty"`
	City              *string `json:"city,omitempty"`
	State             *string `json:"state,omitempty"`
	Country           *string `json:"country,omitempty"`
	OnGCC             *bool   `json:"onGcc,omitempty"`
	// Barter            *bool                  `json:"barter,omitempty"`
	UserCategories              *[]int64               `json:"userCategories,omitempty"`
	Languages                   *[]string              `json:"languages,omitempty"`
	Dob                         *string                `json:"dob,omitempty"`
	Bio                         *string                `json:"bioOnApp,omitempty"`
	Label                       *string                `json:"label,omitempty"`
	Gender                      *string                `json:"gender,omitempty"`
	SocialAccounts              []SocialAccountInput   `json:"socialAccounts,omitempty"`
	CreatorPrograms             *[]CreatorProgramInput `json:"creatorPrograms"`
	CreatorCohorts              *[]CreatorCohortInput  `json:"creatorCohorts"`
	AccountId                   *string                `json:"accountId,omitempty"`
	WhatsappOptin               *bool                  `json:"whatsappOptIn,omitempty"`
	InstantGratificationInvited *bool                  `json:"instantGratificationInvited,omitempty"`
	AmazonStoreLinkVerified     *bool                  `json:"amazonStoreLinkVerified,omitempty"`
	SecondaryPhone              *string                `json:"secondaryPhone,omitempty"`
	NotificationToken           *string                `json:"notificationToken,omitempty"`
	MemberId                    *string                `json:"memberId,omitempty"`
	ReferenceCode               *string                `json:"referenceCode,omitempty"`
	EKYCPending                 *bool                  `json:"ekycPending,omitempty"`
	WebEngageUserId             *string                `json:"webengageUserId,omitempty"`
}

type CreatorProgramInput struct {
	Id    int    `json:"id,omitempty"`
	Level string `json:"level,omitempty"`
	Tag   string `json:"tag,omitempty"`
}

type CreatorCohortInput struct {
	Id    int    `json:"id,omitempty"`
	Level string `json:"level,omitempty"`
	Tag   string `json:"tag,omitempty"`
}

type SocialAccountInput struct {
	Handle   string            `json:"handle"`
	Platform SocialNetworkType `json:"platform"`
	CpId     string            `json:"cpId"`
}

type CreatorProgramTagInput struct {
	ID    int    `json:"id"`
	Level string `json:"level"`
	Tag   string `json:"tag"`
}

type CreatorCohortTagInput struct {
	ID    int    `json:"id"`
	Level string `json:"level"`
	Tag   string `json:"tag"`
}

type SocialNetworkType string

const (
	SocialNetworkTypeGoogle    SocialNetworkType = "GOOGLE"
	SocialNetworkTypeFb        SocialNetworkType = "FB"
	SocialNetworkTypeTwitter   SocialNetworkType = "TWITTER"
	SocialNetworkTypeTiktok    SocialNetworkType = "TIKTOK"
	SocialNetworkTypeInstagram SocialNetworkType = "INSTAGRAM"
	SocialNetworkTypeYoutube   SocialNetworkType = "YOUTUBE"
	SocialNetworkTypeWhatsapp  SocialNetworkType = "WHATSAPP"
	SocialNetworkTypeOthers    SocialNetworkType = "OTHERS"
	SocialNetworkTypeSnapchat  SocialNetworkType = "SNAPCHAT"
	SocialNetworkTypePinterest SocialNetworkType = "PINTEREST"
	SocialNetworkTypeLinkedin  SocialNetworkType = "LINKEDIN"
	SocialNetworkTypeBulbul    SocialNetworkType = "BULBUL"
	SocialNetworkTypeDiscord   SocialNetworkType = "DISCORD"
	SocialNetworkTypeTwitch    SocialNetworkType = "TWITCH"
	SocialNetworkTypePatreon   SocialNetworkType = "PATREON"
	SocialNetworkTypeNone      SocialNetworkType = "NONE"
	SocialNetworkTypeAmazon    SocialNetworkType = "AMAZON"
	SocialNetworkTypeGcc       SocialNetworkType = "GCC"
)

var AllSocialNetworkType = []SocialNetworkType{
	SocialNetworkTypeGoogle,
	SocialNetworkTypeFb,
	SocialNetworkTypeTwitter,
	SocialNetworkTypeTiktok,
	SocialNetworkTypeInstagram,
	SocialNetworkTypeYoutube,
	SocialNetworkTypeWhatsapp,
	SocialNetworkTypeOthers,
	SocialNetworkTypeSnapchat,
	SocialNetworkTypePinterest,
	SocialNetworkTypeLinkedin,
	SocialNetworkTypeBulbul,
	SocialNetworkTypeDiscord,
	SocialNetworkTypeTwitch,
	SocialNetworkTypePatreon,
	SocialNetworkTypeNone,
	SocialNetworkTypeAmazon,
	SocialNetworkTypeGcc,
}

func (e SocialNetworkType) IsValid() bool {
	switch e {
	case SocialNetworkTypeGoogle, SocialNetworkTypeFb, SocialNetworkTypeTwitter, SocialNetworkTypeTiktok, SocialNetworkTypeInstagram, SocialNetworkTypeYoutube, SocialNetworkTypeWhatsapp, SocialNetworkTypeOthers, SocialNetworkTypeSnapchat, SocialNetworkTypePinterest, SocialNetworkTypeLinkedin, SocialNetworkTypeBulbul, SocialNetworkTypeDiscord, SocialNetworkTypeTwitch, SocialNetworkTypePatreon, SocialNetworkTypeNone, SocialNetworkTypeAmazon, SocialNetworkTypeGcc:
		return true
	}
	return false
}

func (e SocialNetworkType) String() string {
	return string(e)
}

func (e *SocialNetworkType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = SocialNetworkType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid SocialNetworkType", str)
	}
	return nil
}

func (e SocialNetworkType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type SyncWinklApiEntry struct {
	Id *int64 `json:"campaignId,omitempty"`
}

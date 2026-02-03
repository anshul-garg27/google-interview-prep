package dao

import (
	"time"

	"github.com/lib/pq"
)

type CampaignProfileEntity struct {
	ID                int64     `gorm:"<-create:primaryKey,id"`
	GCCUserAccountID  *int64    `gorm:"column:gcc_user_account_id"`
	Platform          string    `gorm:"column:platform"`
	PlatformAccountId int64     `gorm:"column:platform_account_id"`
	OnGCC             *bool     `gorm:"column:on_gcc;default:false"`
	OnGCCApp          *bool     `gorm:"column:on_gcc_app;default:false"`
	AdminDetails      *string   `gorm:"column:admin_details;type:jsonb"`
	UserDetails       *string   `gorm:"column:user_details;type:jsonb"`
	Enabled           bool      `gorm:"column:enabled"`
	CreatedAt         time.Time `gorm:"default:current_timestamp"`
	UpdatedAt         time.Time `gorm:"default:current_timestamp"`
	UpdatedBy         *string   `gorm:"column:updated_by"`
	HasEmail          *bool     `gorm:"column:has_email;default:false"`
	HasPhone          *bool     `gorm:"column:has_phone;default:false"`
	Gender            *string
	Phone             pq.StringArray `gorm:"column:phone;type:text[]"`
	Dob               *string
	Location          pq.StringArray `gorm:"column:location;type:text[]"`
	Languages         pq.StringArray `gorm:"column:languages;type:text[]"`
}

func (CampaignProfileEntity) TableName() string {
	return "campaign_profiles"
}

type AdminDetailsEntity struct {
	Name                *string   `json:"name,omitempty"`
	Email               *string   `json:"email,omitempty"`
	Phone               *string   `json:"phone,omitempty"`
	SecondaryPhone      *string   `json:"secondaryPhone,omitempty"`
	Gender              *string   `json:"gender,omitempty"`
	Languages           *[]string `json:"languages,omitempty"`
	City                *string   `json:"city,omitempty"`
	State               *string   `json:"state,omitempty"`
	Country             *string   `json:"country,omitempty"`
	Dob                 *string   `json:"dob,omitempty"`
	Bio                 *string   `json:"bio,omitempty"`
	CampaignCategoryIds *[]string `json:"campaignCategoryIds,omitempty"`
	CreatorPrograms     *[]string `json:"creatorPrograms,omitempty"`
	IsBlacklisted       bool      `json:"isBlacklisted"`
	BlacklistedBy       *string   `json:"blacklistedBy,omitempty"`
	BlacklistedReason   *string   `json:"blacklistedReason,omitempty"`
	Location            *[]string `json:"location,omitempty"`
	WhatsappOptIn       bool      `json:"whatsappOptIn"`
	CountryCode         *string   `json:"countryCode,omitempty"`
	CreatorCohorts      *[]string `json:"creatorCohorts,omitempty"`
}

type UserDetailsEntity struct {
	Email                       *string   `json:"email,omitempty"`
	Phone                       *string   `json:"phone,omitempty"`
	SecondaryPhone              *string   `json:"secondaryPhone,omitempty"`
	Name                        *string   `json:"name,omitempty"`
	Dob                         *string   `json:"dob,omitempty"`
	Gender                      *string   `json:"gender,omitempty"`
	WhatsappOptIn               bool      `json:"whatsappOptIn"`
	CampaignCategoryIds         *[]string `json:"campaignCategoryIds,omitempty"`
	Languages                   *[]string `json:"languages,omitempty"`
	Location                    *[]string `json:"location,omitempty"`
	NotificationToken           *string   `json:"notificationToken,omitempty"`
	Bio                         *string   `json:"bio,omitempty"`
	MemberId                    *string   `json:"memberId,omitempty"`
	ReferenceCode               *string   `json:"referenceCode,omitempty"`
	EKYCPending                 bool      `json:"ekycPending"`
	WebEngageUserId             *string   `json:"webengageUserId,omitempty"`
	InstantGratificationInvited bool      `json:"instantGratificationInvited"`
	AmazonStoreLinkVerified     bool      `json:"amazonStoreLinkVerified"`
}

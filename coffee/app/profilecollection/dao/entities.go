package dao

import (
	"time"

	"github.com/jackc/pgtype"

	"github.com/lib/pq"
)

type ProfileCollectionEntity struct {
	Id               int64          `gorm:"<-:create;primaryKey"`
	Name             string         `gorm:"column:name"`
	PartnerId        int64          `gorm:"<-:create;column:partner_id"`
	Description      *string        `gorm:"column:description"`
	ShareId          string         `gorm:"column:share_id"`
	Source           string         `gorm:"<-:create;column:source"`
	SourceId         string         `gorm:"<-:create;column:source_id"`
	Categories       pq.StringArray `gorm:"column:categories;type:text[]"`
	CategoryIds      pgtype.JSONB   `gorm:"column:category_ids;type:jsonb"`
	Tags             pq.StringArray `gorm:"column:tags;type:text[]"`
	Featured         *bool          `gorm:"column:featured"`
	Enabled          *bool          `gorm:"column:enabled"`
	AnalyticsEnabled *bool          `gorm:"column:analytics_enabled"`
	DisabledMetrics  pgtype.JSONB   `gorm:"column:disabled_metrics;type:jsonb"`
	CampaignId       *int64         `gorm:"<-:create;column:campaign_id"`
	CampaignPlatform *string        `gorm:"<-:create;column:campaign_platform"`
	JobId            *int64         `gorm:"column:job_id"`
	CustomColumns    pgtype.JSONB   `gorm:"column:custom_columns;type:jsonb"`
	OrderedColumns   pgtype.JSONB   `gorm:"column:ordered_columns;type:jsonb"`
	ShareEnabled     *bool          `gorm:"column:share_enabled"`

	CreatedBy string    `gorm:"<-:create;column:created_by"`
	CreatedAt time.Time `gorm:"<-:create;column:created_at;type:time"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:time"`
}

type ProfileCollectionItemEntity struct {
	Id                  int64     `gorm:"<-:create;primaryKey;column:id"`
	Platform            string    `gorm:"<-:create;column:platform"`
	PlatformAccountCode int64     `gorm:"<-:create;column:platform_account_code"`
	CampaignProfileId   *int64    `gorm:"column:campaign_profile_id"`
	ProfileSocialId     string    `gorm:"<-create;column:profile_social_id"`
	ProfileCollectionId int64     `gorm:"<-:create;column:profile_collection_id"`
	PartnerId           int64     `gorm:"<-:create;column:partner_id"`
	Enabled             *bool     `gorm:"column:enabled"`
	ShortlistingStatus  *string   `gorm:"column:shortlisting_status"`
	ShortlistId         *string   `gorm:"column:shortlist_id"`
	CreatedBy           string    `gorm:"<-:create;column:created_by"`
	CreatedAt           time.Time `gorm:"<-:create;column:created_at;type:time"`
	UpdatedAt           time.Time `gorm:"column:updated_at;type:time"`

	Hidden                *bool   `gorm:"column:hidden"`
	Rank                  int64   `gorm:"column:rank"`
	BrandSelectionStatus  *string `gorm:"column:brand_selection_status"`
	BrandSelectionRemarks *string `gorm:"column:brand_selection_remarks"`
	InternalCommercials   *string `gorm:"column:internal_commercials"`
}

type ProfileCollectionItemCustomColumnEntity struct {
	Id                      int64  `gorm:"column:id"`
	ProfileCollectionItemId int64  `gorm:"<-:create;column:profile_collection_item_id"`
	Key                     string `gorm:"<-:create;column:key"`
	Value                   string `gorm:"column:value"`
}

func (ProfileCollectionEntity) TableName() string {
	return "profile_collection"
}

func (ProfileCollectionItemEntity) TableName() string {
	return "profile_collection_item"
}

func (ProfileCollectionItemCustomColumnEntity) TableName() string {
	return "profile_collection_item_cc_data"
}

type WinklCollectionInfoEntity struct {
	WinklCollectionId        int64  `gorm:"<-;winkl_collection_id"`
	WinklCollectionShareId   string `gorm:"<-;winkl_collection_share_id"`
	WinklCampaignId          *int64 `gorm:"<-;winkl_campaign_id"`
	WinklShortlistId         *int64 `gorm:"<-;winkl_shortlist_id"`
	ProfileCollectionId      int64  `gorm:"profile_collection_id"`
	ProfileCollectionShareId string `gorm:"profile_collection_share_id"`
}

func (WinklCollectionInfoEntity) TableName() string {
	return "winkl_collection_migration_info"
}

type CollectionItemCountsEntity struct {
	Collectionid int64
	Platform     string
	Allitems     int64
	Emailitems   int64
	Phoneitems   int64
}

type CollectionItemPlatformCustomColumnSum struct {
	Platform string
	Key      string
	Items    int64
	Sum      float64
}

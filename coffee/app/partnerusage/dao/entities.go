package dao

import (
	"time"
)

type PartnerUsageEntity struct {
	ID          int64 `gorm:"<-create:primaryKey;column:id"`
	PartnerId   int64
	Namespace   string
	Key         string
	Limit       int64
	Consumed    *int64
	NextResetOn time.Time
	Frequency   string
	StartDate   time.Time
	EndDate     time.Time
	Enabled     *bool
	CreatedAt   time.Time `gorm:"default:current_timestamp"`
	UpdatedAt   time.Time `gorm:"default:current_timestamp"`
}

func (PartnerUsageEntity) TableName() string {
	return "partner_usage"
}

type PartnerProfileTrackEntity struct {
	ID        int64 `gorm:"<-create:primaryKey;column:id"`
	PartnerId int64
	Key       string
	Enabled   *bool
	CreatedAt time.Time `gorm:"default:current_timestamp"`
	UpdatedAt time.Time `gorm:"default:current_timestamp"`
}

func (PartnerProfileTrackEntity) TableName() string {
	return "partner_profile_page_track"
}

type AcitivityTrackerEntity struct {
	Id        int64   `gorm:"<-create:primaryKey,id"`
	PartnerId int64   `gorm:"partner_id"`
	AccountId int64   `gorm:"account_id"`
	Type      string  `gorm:"type"`
	Meta      *string `gorm:"meta"`
	Enabled   *bool
	CreatedAt time.Time `gorm:"default:current_timestamp"`
	UpdatedAt time.Time `gorm:"default:current_timestamp"`
}

func (AcitivityTrackerEntity) TableName() string {
	return "activity_tracker"
}

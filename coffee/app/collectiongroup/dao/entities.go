package dao

import (
	"time"

	"github.com/jackc/pgtype"
	"github.com/lib/pq"
)

type CollectionGroupEntity struct {
	Id            int64          `gorm:"<-:create;primaryKey"`
	Name          string         `gorm:"column:name"`
	PartnerId     int64          `gorm:"<-:create;column:partner_id"`
	ShareId       string         `gorm:"column:share_id"`
	Objective     string         `gorm:"column:objective"`
	Source        string         `gorm:"<-:create;column:source"`
	SourceId      string         `gorm:"<-:create;column:source_id"`
	CollectionIds pq.StringArray `gorm:"column:collection_ids;type:jsonb"`
	Enabled       *bool          `gorm:"column:enabled"`
	Metadata      pgtype.JSONB   `gorm:"column:metadata;type:jsonb"`
	CreatedBy     string         `gorm:"<-:create;column:created_by"`
	CreatedAt     time.Time      `gorm:"<-:create;column:created_at;type:time"`
	UpdatedAt     time.Time      `gorm:"column:updated_at;type:time"`
}

func (CollectionGroupEntity) TableName() string {
	return "collection_group"
}

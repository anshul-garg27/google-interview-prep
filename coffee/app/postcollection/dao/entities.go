package dao

import (
	"time"

	"github.com/jackc/pgtype"
)

type PostCollectionEntity struct {
	Id                    string       `gorm:"<-:create;primaryKey;column:id"`
	ShareId               string       `gorm:"column:share_id"`
	Name                  string       `gorm:"column:name"`
	PartnerId             int64        `gorm:"<-:create;column:partner_id"`
	Enabled               *bool        `gorm:"column:enabled"`
	Source                string       `gorm:"<-:create;column:source"`
	SourceId              string       `gorm:"<-:create;column:source_id"`
	Budget                int          `gorm:"column:budget"`
	StartTime             *time.Time   `gorm:"column:start_time;type:time"`
	EndTime               *time.Time   `gorm:"column:end_time;type:time"`
	CreatedBy             string       `gorm:"<-:create;column:created_by"`
	CreatedAt             time.Time    `gorm:"<-:create;column:created_at;type:time"`
	UpdatedAt             time.Time    `gorm:"column:updated_at;type:time"`
	DisabledMetrics       pgtype.JSONB `gorm:"column:disabled_metrics;type:jsonb"`
	ExpectedMetricValues  pgtype.JSONB `gorm:"column:expected_metric_values;type:jsonb"`
	MetricsIngestionFreq  string       `gorm:"column:metrics_ingestion_freq"`
	Comments              string       `gorm:"column:comments"`
	JobId                 *int64       `gorm:"column:job_id"`
	SentimentReportPath   *string      `gorm:"column:sentiment_report_path"`
	SentimentReportBucket *string      `gorm:"column:sentiment_report_bucket"`
}

type PostCollectionItemEntity struct {
	Id                   int64        `gorm:"<-:create;primaryKey;column:id"`
	Platform             string       `gorm:"<-:create;column:platform"`
	ShortCode            string       `gorm:"<-:create;column:short_code"`
	PostCollectionId     string       `gorm:"<-:create;column:post_collection_id"`
	PostType             string       `gorm:"<-:create;column:post_type"`
	Cost                 *int64       `gorm:"column:cost"`
	Enabled              *bool        `gorm:"column:enabled"`
	Bookmarked           *bool        `gorm:"column:bookmarked"`
	PostedByHandle       *string      `gorm:"column:posted_by_handle"`
	SponsorLinks         pgtype.JSONB `gorm:"column:sponsor_links;type:jsonb"`
	PostTitle            *string      `gorm:"column:post_title"`
	PostLink             *string      `gorm:"column:post_link"`
	PostThumbnail        *string      `gorm:"column:post_thumbnail"`
	CreatedBy            string       `gorm:"<-:create;column:created_by"`
	MetricsIngestionFreq string       `gorm:"column:metrics_ingestion_freq"`
	RetrieveData         *bool        `gorm:"column:retrieve_data;default:true"`
	ShowInReport         *bool        `gorm:"column:show_in_report;default:true"`
	PostedByCpId         *int64       `gorm:"column:posted_by_cp_id"`
	CreatedAt            time.Time    `gorm:"<-:create;column:created_at;type:time"`
	UpdatedAt            time.Time    `gorm:"column:updated_at;type:time"`
}

func (PostCollectionEntity) TableName() string {
	return "post_collection"
}

func (PostCollectionItemEntity) TableName() string {
	return "post_collection_item"
}

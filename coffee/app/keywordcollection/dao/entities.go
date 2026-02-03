package dao

import (
	"time"

	"github.com/lib/pq"
)

type KeywordCollectionEntity struct {
	Id                   string         `gorm:"<-:create;primaryKey;column:id"`
	ShareId              string         `gorm:"column:share_id"`
	Name                 string         `gorm:"column:name"`
	PartnerId            int64          `gorm:"<-:create;column:partner_id"`
	Platform             string         `gorm:"<-:create;column:platform"`
	Keywords             pq.StringArray `gorm:"column:keywords;type:text[]"`
	Enabled              *bool          `gorm:"column:enabled"`
	StartTime            *time.Time     `gorm:"column:start_time;type:time"`
	EndTime              *time.Time     `gorm:"column:end_time;type:time"`
	CreatedBy            string         `gorm:"<-:create;column:created_by"`
	CreatedAt            time.Time      `gorm:"<-:create;column:created_at;type:time"`
	UpdatedAt            time.Time      `gorm:"column:updated_at;type:time"`
	JobId                *int64         `gorm:"column:job_id"`
	ReportBucket         *string        `gorm:"column:report_bucket"`
	ReportPath           *string        `gorm:"column:report_path"`
	DemographyReportPath *string        `gorm:"column:demography_report_path"`
}

func (KeywordCollectionEntity) TableName() string {
	return "keyword_collection"
}

type Summary struct {
	Posts          int64
	Reach          int64
	Likes          int64
	Comments       int64
	Engagement     int64
	Emp            int64
	EngagementRate float64
}

type Demography struct {
	AudienceAge       string
	AudienceGender    string
	AudienceAgeGender string
}

type Post struct {
	Shortcode           string
	PostThumbnailUrl    string
	ProfileThumbnailUrl string
	Reach               int64
	PostPublishTime     time.Time
	Views               int64
	Plays               int64
	Likes               int64
	PostType            string
	Comments            int64
	Followers           int64
	Category            string
	EngagementRate      float64
	Handle              string
	ProfileId           string
	Name                string
	PostTitle           string
}

type Profile struct {
	ProfileId      string
	Handle         string
	Thumbnail      string
	Followers      int64
	Category       string
	Name           string
	Reach          int64
	Views          int64
	Likes          int64
	Comments       int64
	Uploads        int64
	LastPostedOn   time.Time
	EngagementRate float64
}

type Category struct {
	Category  string
	Reach     int64
	Posts     int64
	PostsPerc float64
}

type Language struct {
	Month     string
	Language  string
	Reach     int64
	Posts     int64
	PostsPerc float64
}

type Keywords struct {
	Keyword   string
	Reach     int64
	Posts     int64
	PostsPerc float64
	Score     float64
}

type MonthlyStat struct {
	Month   string
	Views   int64
	Uploads int64
}

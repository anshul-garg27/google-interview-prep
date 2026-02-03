package domain

import (
	discoverydomain "coffee/app/discovery/domain"
	coredomain "coffee/core/domain"
)

// Profile Collection Entry

type KeywordCollectionEntry struct {
	Id                   *string                  `json:"id,omitempty"`
	ShareId              *string                  `json:"shareId,omitempty"`
	Name                 *string                  `json:"name,omitempty"`
	PartnerId            *int64                   `json:"partnerId,omitempty"`
	Platform             *string                  `json:"platform,omitempty"`
	Enabled              *bool                    `json:"enabled,omitempty"`
	StartTime            *int64                   `json:"startTime,omitempty"`
	EndTime              *int64                   `json:"endTime,omitempty"`
	JobId                *int64                   `json:"jobId,omitempty"`
	Keywords             *[]string                `json:"keywords,omitempty"`
	CreatedBy            *string                  `json:"createdBy,omitempty"`
	CreatedAt            *int64                   `json:"createdAt,omitempty"`
	UpdatedAt            *int64                   `json:"updatedAt,omitempty"`
	ReportBucket         *string                  `json:"report_bucket"`
	ReportPath           *string                  `json:"report_path"`
	DemographyReportPath *string                  `json:"demography_report_path"`
	Report               *KeywordCollectionReport `json:"report,omitempty"`
	ReportStatus         string                   `json:"reportStatus,omitempty"`
}

type KeywordCollectionReport struct {
	TotalPosts          int64                             `json:"totalPosts"`
	TotalReach          int64                             `json:"totalReach"`
	TotalLikes          int64                             `json:"totalLikes"`
	TotalComments       int64                             `json:"totalComments"`
	Engagement          int64                             `json:"engagement,omitempty"`
	EMP                 float64                           `json:"emp"`
	EngagementRate      float64                           `json:"engagementRate"`
	TotalViews          int64                             `json:"totalViews"`
	TopKeywords         map[string]float64                `json:"topKeywords"`
	TopCategories       map[string]float64                `json:"topCategories"`
	TopLanguages        map[string]float64                `json:"topLanguages"`
	TopProfiles         []KeywordCollectionProfile        `json:"topProfiles"`
	TopPosts            []KeywordCollectionPost           `json:"topPosts"`
	TopLanguagesMonthly map[string]map[string]float64     `json:"topLanguagesMonthly"`
	MonthlyViews        map[string]int64                  `json:"monthlyViews"`
	MonthlyUploads      map[string]int64                  `json:"monthlyUploads"`
	AudienceAgeGender   discoverydomain.AudienceAgeGender `json:"audienceAgeGender,omitempty"`
	AudienceGender      map[string]*float64               `json:"audienceGender,omitempty"`
}

type KeywordCollectionPost struct {
	AccountInfo    coredomain.SocialAccount `json:"accountInfo"`
	Link           string                   `json:"postLink,omitempty"`
	Title          string                   `json:"title,omitempty"`
	Platform       string                   `json:"platform,omitempty"`
	Thumbnail      string                   `json:"thumbnail,omitempty"`
	Views          *int64                   `json:"views,omitempty"`
	PublishedAt    *int64                   `json:"publishedAt,omitempty"`
	Plays          *int64                   `json:"plays,omitempty"`
	Likes          *int64                   `json:"likes,omitempty"`
	PostType       *string                  `json:"postType,omitempty"`
	Comments       *int64                   `json:"comments,omitempty"`
	Engagement     *int64                   `json:"engagement,omitempty"`
	EngagementRate *float64                 `json:"engagementRate,omitempty"`
}

type KeywordCollectionProfile struct {
	AccountInfo coredomain.SocialAccount `json:"accountInfo"`
	Views       *int64                   `json:"views,omitempty"`
	Plays       *int64                   `json:"plays,omitempty"`
	Likes       *int64                   `json:"likes,omitempty"`
	Comments    *int64                   `json:"comments,omitempty"`
	Engagement  *int64                   `json:"engagement,omitempty"`
	Uploads     *int64                   `json:"uploads,omitempty"`
}

type TopicResearchActivityMeta struct {
	Id       string              `json:"id,omitempty"`
	Platform string              `json:"platform,omitempty"`
	Result   TopicResearchResult `json:"result,omitempty"`
}

type TopicResearchResult struct {
	Name      string    `json:"name,omitempty"`
	Reach     int64     `json:"reach,omitempty"`
	PostCount int64     `json:"postCount,omitempty"`
	Keywords  *[]string `json:"keywords,omitempty"`
}

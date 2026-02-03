package domain

// Profile Collection Entry

type PostCollectionEntry struct {
	Id                    *string                   `json:"id,omitempty"`
	ShareId               *string                   `json:"shareId,omitempty"`
	Name                  *string                   `json:"name,omitempty"`
	PartnerId             *int64                    `json:"partnerId,omitempty"`
	Source                *string                   `json:"source,omitempty"`
	SourceId              *string                   `json:"sourceId,omitempty"`
	Budget                *int                      `json:"budget,omitempty"`
	Enabled               *bool                     `json:"enabled,omitempty"`
	StartTime             *int64                    `json:"startTime,omitempty"`
	EndTime               *int64                    `json:"endTime,omitempty"`
	JobId                 *int64                    `json:"jobId,omitempty"`
	ItemsAssetId          *int64                    `json:"itemsAssetId,omitempty"`
	Items                 []PostCollectionItemEntry `json:"items,omitempty"`
	CreatedBy             *string                   `json:"createdBy,omitempty"`
	CreatedAt             *int64                    `json:"createdAt,omitempty"`
	UpdatedAt             *int64                    `json:"updatedAt,omitempty"`
	DisabledMetrics       []string                  `json:"disabledMetrics,omitempty"`
	MetricsIngestionFreq  *string                   `json:"metricsIngestionFreq,omitempty"`
	SentimentReportPath   *string                   `json:"sentimentReportPath,omitempty"`
	SentimentReportBucket *string                   `json:"sentimentReportBucket,omitempty"`
	ExpectedMetricValues  map[string]*float64       `json:"expectedMetricValues,omitempty"`
	Comments              *string                   `json:"comments,omitempty"`
	ReportReady           *bool                     `json:"reportReady,omitempty"`
	Reach                 int64                     `json:"reach,omitempty"`
	CostPerReach          string                    `json:"costPerReach,omitempty"`
	TotalPosts            int64                     `json:"totalPosts,omitempty"`
}

// Item Entry

type PostCollectionItemEntry struct {
	Id                      *int64   `json:"id,omitempty"`
	Platform                *string  `json:"platform,omitempty"`
	ShortCode               *string  `json:"shortCode,omitempty"`
	PostCollectionId        *string  `json:"postCollectionId,omitempty"`
	Cost                    *int64   `json:"cost,omitempty"`
	Bookmarked              *bool    `json:"bookmarked,omitempty"`
	Enabled                 *bool    `json:"enabled,omitempty"`
	PostType                *string  `json:"postType,omitempty"`
	PostedByHandle          *string  `json:"postedByHandle,omitempty"`
	HandleProfilePic        *string  `json:"handleProfilePic,omitempty"`
	HandleProfilePicAssetId *int64   `json:"handleProfilePicAssetId,omitempty"`
	SponsorLinks            []string `json:"sponsorLinks,omitempty"`
	PostTitle               *string  `json:"postTitle,omitempty"`
	PostLink                *string  `json:"postLink,omitempty"`
	PostThumbnail           *string  `json:"postThumbnail,omitempty"`
	MetricsIngestionFreq    *string  `json:"metricsIngestionFreq,omitempty"`
	CreatedBy               *string  `json:"createdBy,omitempty"`
	RetrieveData            *bool    `json:"retrieveData,omitempty"`
	ShowInReport            *bool    `json:"showInReport,omitempty"`
	PostedByCpId            *int64   `json:"postedByCpId,omitempty"`
	CreatedAt               *int64   `json:"createdAt,omitempty"`
	UpdatedAt               *int64   `json:"updatedAt,omitempty"`
}

type CampaignReportActivityMeta struct {
	Id       string               `json:"id,omitempty"`
	Platform string               `json:"platform,omitempty"`
	Result   CampaignReportResult `json:"result,omitempty"`
}
type CampaignReportResult struct {
	Name      string `json:"name,omitempty"`
	Reach     int64  `json:"reach,omitempty"`
	PostCount int64  `json:"postCount,omitempty"`
}

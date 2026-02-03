package domain

import (
	coredomain "coffee/core/domain"
)

// Profile Collection Entry

type ProfileCollectionEntry struct {
	Id                       *int64                              `json:"id,omitempty"`
	Name                     *string                             `json:"name,omitempty"`
	PartnerId                *int64                              `json:"partnerId,omitempty"`
	PartnerName              *string                             `json:"partnerName,omitempty"`
	Source                   *string                             `json:"source,omitempty"`
	SourceId                 *string                             `json:"sourceId,omitempty"`
	Platform                 *string                             `json:"platform,omitempty"`
	Description              *string                             `json:"description,omitempty"`
	Featured                 *bool                               `json:"featured,omitempty"`
	Enabled                  *bool                               `json:"enabled,omitempty"`
	AnalyticsEnabled         *bool                               `json:"analyticsEnabled"`
	CreatedBy                *string                             `json:"createdBy,omitempty"`
	ShareId                  *string                             `json:"shareId,omitempty"`
	Categories               *[]string                           `json:"categories,omitempty"`
	CategoryIds              []int64                             `json:"categoryIds,omitempty"`
	Tags                     *[]string                           `json:"tags,omitempty"`
	OrderedColumns           []PlatformOrderedColumnData         `json:"orderedColumns,omitempty"`
	ShareEnabled             *bool                               `json:"shareEnabled"`
	CampaignId               *int64                              `json:"campaignId,omitempty"`
	CampaignPlatform         *string                             `json:"campaignPlatform,omitempty"`
	JobId                    *int64                              `json:"jobId,omitempty"`
	ItemsAssetId             *int64                              `json:"itemsAssetId,omitempty"`
	Items                    []ProfileCollectionItemEntry        `json:"items,omitempty"`
	ImportCollectionId       *int64                              `json:"importCollectionId,omitempty"`
	CreatedAt                *int64                              `json:"createdAt,omitempty"`
	UpdatedAt                *int64                              `json:"updatedAt,omitempty"`
	DisabledMetrics          []string                            `json:"disabledMetrics,omitempty"`
	ReportReady              *bool                               `json:"reportReady,omitempty"`
	Reach                    int64                               `json:"reach,omitempty"`
	CostPerReach             string                              `json:"costPerReach,omitempty"`
	TotalPosts               int64                               `json:"totalPosts,omitempty"`
	EngagementRate           string                              `json:"engagementRate,omitempty"`
	PlatformWiseItemsSummary map[string]PlatformWiseItemsSummary `json:"platformWiseItemsSummary,omitempty"`
	CustomColumns            []PlatformCustomColumnMeta          `json:"customColumns,omitempty"`
}

type PlatformWiseItemsSummary struct {
	ProfileCounter  ProfileCounter         `json:"profileCount,omitempty"`
	TopFiveAccounts []CompactSocialAccount `json:"topFiveAccounts,omitempty"`
}

type CompactSocialAccount struct {
	Code      *string `json:"code"`
	Name      *string `json:"name"`
	Thumbnail *string `json:"thumbnail,omitempty"`
	Link      *string `json:"link,omitempty"`
}

type ProfileCounter struct {
	Count      int64 `json:"count"`
	EmailCount int64 `json:"emailCount"`
	PhoneCount int64 `json:"phoneCount"`
}

// Item Entry

type ProfileCollectionItemEntry struct {
	Id                  *int64  `json:"id"`
	Platform            *string `json:"platform"`
	PlatformAccountCode *int64  `json:"platformAccountCode"`
	CampaignProfileId   *int64  `json:"campaignProfileId"`
	ProfileSocialId     *string `json:"profileSocialId"`
	ProfileCollectionId *int64  `json:"profileCollectionId"`
	PartnerId           *int64  `json:"partnerId"`
	ShortlistingStatus  *string `json:"shortlistingStatus"`
	Enabled             *bool   `json:"enabled"`
	CreatedBy           *string `json:"createdBy"`
	ShortlistId         *string `json:"shortlistId"`
	CreatedAt           *int64  `json:"createdAt,omitempty"`
	UpdatedAt           *int64  `json:"updatedAt,omitempty"`

	Hidden                *bool          `json:"hidden,omitempty"`
	Rank                  *int64         `json:"rank,omitempty"`
	CustomColumnsData     []CustomColumn `json:"customColumnsData,omitempty"`
	BrandSelectionStatus  *string        `json:"brandSelectionStatus,omitempty"`
	BrandSelectionRemarks *string        `json:"brandSelectionRemarks,omitempty"`
	Commercials           *string        `json:"commercials,omitempty"`

	Profile *coredomain.Profile `json:"profile,omitempty"`
}

type CustomColumn struct {
	Id    *int64 `json:"id"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

type PlatformCustomColumnMeta struct {
	Platform string             `json:"platform"`
	Columns  []CustomColumnMeta `json:"columns"`
}

type CustomColumnMeta struct {
	Key           *string `json:"key"`
	Title         string  `json:"title"`
	Value         *string `json:"value"`
	Hidden        *bool   `json:"hidden,omitempty"`
	AllowEdit     *bool   `json:"allowEdit,omitempty"`
	Compute       *string `json:"compute,omitempty"`
	ComputedValue string  `json:"computedValue,omitempty"`
	ComputedError *string `json:"computeError,omitempty"`
}

type PlatformOrderedColumnData struct {
	Platform string   `json:"platform"`
	Columns  []string `json:"columns"`
}

type WinklCollectionInfoEntry struct {
	WinklCollectionId        int64  `json:"winklCollectionId"`
	WinklCollectionShareId   string `json:"winklCollectionShareId"`
	WinklCampaignId          *int64 `json:"winklCampaignId"`
	WinklShortlistId         *int64 `json:"winklShortlistId"`
	ProfileCollectionId      int64  `json:"profileCollectionId"`
	ProfileCollectionShareId string `json:"profileCollectionShareId"`
}

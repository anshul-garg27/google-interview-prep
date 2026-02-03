package domain

type CollectionGroupMetadata struct {
	ExpectedMetricValues map[string]*float64 `json:"expectedMetricValues,omitempty"`
	DisabledMetrics      []string            `json:"disabledMetrics,omitempty"`
	Comments             *string             `json:"comments,omitempty"`
}

type CollectionGroupEntry struct {
	Id            *int64                   `json:"id,omitempty"`
	Name          *string                  `json:"name,omitempty"`
	CollectionIds *[]string                `json:"collectionIds,omitempty"`
	PartnerId     *int64                   `json:"partnerId,omitempty"`
	Objective     *string                  `json:"objective,omitempty"`
	PartnerName   *string                  `json:"partnerName,omitempty"`
	Source        *string                  `json:"source,omitempty"`
	SourceId      *string                  `json:"sourceId,omitempty"`
	Enabled       *bool                    `json:"enabled,omitempty"`
	CreatedBy     *string                  `json:"createdBy,omitempty"`
	ShareId       *string                  `json:"shareId,omitempty"`
	Metadata      *CollectionGroupMetadata `json:"metadata,omitempty"`
}

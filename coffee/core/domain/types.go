package domain

import (
	"gorm.io/plugin/optimisticlock"
)

type DataStore string

const (
	Postgres DataStore = "POSTGRES"
)

type ID interface{ ~int64 | ~string }

type Entry interface{}

type Entity interface{}

type BaseEntity struct {
	Version optimisticlock.Version
}

type Response interface{}

type BaseResponse struct {
	Status Status `json:"status"`
}

type SearchQuery struct {
	Filters []SearchFilter `json:"filters"`
}

type SearchFilter struct {
	FilterType string `json:"filterType"`
	Field      string `json:"field"`
	Value      string `json:"value"`
}

type Status struct {
	Status        string `json:"status,omitempty"`
	Message       string `json:"message,omitempty"`
	FilteredCount int64  `json:"filteredCount,omitempty"`
	HasNextPage   bool   `json:"hasNextPage,omitempty"`
	NextCursor    string `json:"nextCursor,omitempty"`
	ActivityId    int64  `json:"activityId,omitempty"`
}

type Posts struct {
	Posts []Post `json:"posts"`
}

type Post struct {
	Platform  *string `json:"platform"`
	ShortCode *string `json:"shortCode"`
}

type AssetInfoResponse struct {
	Status Status      `json:"status,omitempty"`
	List   []AssetInfo `json:"list,omitempty"`
}

type AssetInfo struct {
	Id       int64   `json:"id,omitempty"`
	Url      *string `json:"url,omitempty"`
	Width    *int    `json:"width,omitempty"`
	Height   *int    `json:"height,omitempty"`
	Duration *int    `json:"duration,omitempty"`
}
type JoinClauses struct {
	PreloadVariable  string
	PreloadCondition string
	PreloadValue     string
	JoinClause       string
}
type ErrorLog struct {
	Type        string
	Message     string
	RequestType string
	RequestUrl  string
	RequestBody string
	Headers     map[string][]string
}

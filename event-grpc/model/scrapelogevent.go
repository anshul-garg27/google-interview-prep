package model

import (
	"time"
)

type PostLogEvent struct {
	EventId         string
	Source          string
	Platform        string
	ProfileId       string
	Handle          string
	Shortcode       string
	PublishTime     time.Time
	Metrics         JSONB `sql:"type:jsonb"`
	Dimensions      JSONB `sql:"type:jsonb"`
	EventTimestamp  time.Time
	InsertTimestamp time.Time
}

type ProfileLogEvent struct {
	EventId         string
	Source          string
	Platform        string
	ProfileId       string
	Handle          string
	Metrics         JSONB `sql:"type:jsonb"`
	Dimensions      JSONB `sql:"type:jsonb"`
	EventTimestamp  time.Time
	InsertTimestamp time.Time
}

type ProfileRelationshipLogEvent struct {
	EventId          string
	Source           string
	RelationshipType string
	Platform         string
	SourceProfileId  string
	TargetProfileId  string
	SourceMetrics    JSONB `sql:"type:jsonb"`
	TargetMetrics    JSONB `sql:"type:jsonb"`
	SourceDimensions JSONB `sql:"type:jsonb"`
	TargetDimensions JSONB `sql:"type:jsonb"`
	EventTimestamp   time.Time
	InsertTimestamp  time.Time
}

type ScrapeRequestLogEvent struct {
	EventId        string
	SclId          uint64
	Flow           string
	Platform       string
	Params         JSONB `sql:"type:jsonb"`
	Status         string
	Priority       *int32
	Reason         string
	PickedAt       *time.Time
	ExpiresAt      *time.Time
	EventTimestamp time.Time
}

type OrderLogEvent struct {
	EventId         string
	Source          string
	Platform        string
	PlatformOrderId string
	Store           string
	Metrics         JSONB `sql:"type:jsonb"`
	Dimensions      JSONB `sql:"type:jsonb"`
	EventTimestamp  time.Time
	InsertTimestamp time.Time
}

type SentimentLogEvent struct {
	EventId         string
	Source          string
	Platform        string
	Shortcode       string
	CommentId       string
	Comment         string
	Sentiment       string
	Score           float64
	Metrics         JSONB `sql:"type:jsonb"`
	Dimensions      JSONB `sql:"type:jsonb"`
	EventTimestamp  time.Time
	InsertTimestamp time.Time
}

type PostActivityLogEvent struct {
	EventId         string
	Source          string
	ActivityType    string
	Platform        string
	ActorProfileId  string
	Shortcode       string
	PublishTime     time.Time
	Metrics         JSONB `sql:"type:jsonb"`
	Dimensions      JSONB `sql:"type:jsonb"`
	EventTimestamp  time.Time
	InsertTimestamp time.Time
}

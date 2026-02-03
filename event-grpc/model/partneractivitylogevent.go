package model

import "time"

type PartnerActivityLogEvent struct {
	EventId         string
	ModuleName      string
	ActivityName    string
	IncrementCount  int64
	PartnerId       int64
	AccountId       int64
	ActivityId      *int64
	ActivityMeta    JSONB `sql:"type:jsonb"`
	EventTimestamp  time.Time
	InsertTimestamp time.Time
}

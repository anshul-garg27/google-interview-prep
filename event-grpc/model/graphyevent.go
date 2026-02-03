package model

import "time"

type GraphyEvent struct {
	Id               string
	ReceiveTimestamp time.Time
	InsertTimestamp  time.Time
	EventName        string
	EventData        JSONB `sql:"type:jsonb"`
	Headers          JSONB `sql:"type:jsonb"`
}

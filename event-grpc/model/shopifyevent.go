package model

import (
	"time"
)

type ShopifyEvent struct {
	Id              string
	EventName       string
	StoreUrl        string
	EventTimestamp  time.Time
	InsertTimestamp time.Time
	Data            JSONB
}

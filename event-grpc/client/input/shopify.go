package input

import "time"

type ShopifyEvent struct {
	StoreUrl        string                 `json:"storeUrl,omitempty"`
	EventName       string                 `json:"eventName,omitempty"`
	EventTime       int64                  `json:"eventTime,omitempty"`
	EventData       map[string]interface{} `json:"eventData,omitempty"`
	InsertTimestamp time.Time
}

package model

import (
	"time"
)

type WebengageEvent struct {
	Id          string
	EventName   string
	AccountType string
	Timestamp   time.Time
	UserId      string
	AnonymousId string
	Data        JSONB
}

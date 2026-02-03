package model

import "time"

type AffiliateOrderEvent struct {
	OrderID        string
	Platform       string
	ReferralCode   string
	Meta           JSONB `sql:"type:jsonb"`
	CreatedOn      time.Time
	LastModifiedOn time.Time
	Amount         int64
	Status         string
	DeliveredOn    time.Time
}

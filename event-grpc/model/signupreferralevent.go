package model

import (
	"time"
)

type SignupReferralEvent struct {
	Id                 string
	SignupTimestamp    time.Time
	LaunchTimestamp    time.Time
	UserId             string  `json:"userId,omitempty"`
	DeviceId           string  `json:"deviceId,omitempty"`
	BbDeviceId         string  `json:"bbDeviceId,omitempty"`
	AccountId          string  `json:"accountId,omitempty"`
	ReferralIdentifier string  `json:"referralIdentifier,omitempty"`
	DynamicLink        string  `json:"dynamicLink,omitempty"`
	DeepLink           *string `json:"deepLink,omitempty"`
}

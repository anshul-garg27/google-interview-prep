package model

import (
	"time"
)

type AppInitEvent struct {
	Id              string
	EventTimestamp  time.Time
	InsertTimestamp time.Time
	EventData       JSONB   `sql:"type:jsonb"`
	ClientIp        string  `json:"clientIp,omitempty"`
	UserAgent       string  `json:"userAgent,omitempty"`
	UserId          *string `json:"userId,omitempty"`
	DeviceId        string  `json:"deviceId,omitempty"`
	BbDeviceId      string  `json:"bbDeviceId,omitempty"`
	UserAccountId   *string `json:"userAccountId,omitempty"`
	ClientId        string  `json:"clientId,omitempty"`
	Mozilla         string  `json:"mozilla"`
	Platform        string  `json:"platform"`
	Os              string  `json:"os"`
	OsFullName      string  `json:"osFullName"`
	OsName          string  `json:"osName"`
	OsVersion       string  `json:"osVersion"`
	Localization    string  `json:"localization"`
	Model           string  `json:"model"`
	BrowserName     string  `json:"browserName"`
	BrowserVersion  string  `json:"browserVersion"`
	Bot             int     `json:"bot"`
	Mobile          int     `json:"mobile"`
}

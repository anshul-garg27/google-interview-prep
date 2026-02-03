package model

import (
	"time"
)

type StgClickEvent struct {
	EventId        string    `json:"eventId,omitempty"`
	EventTimestamp time.Time `json:"eventTimestamp,omitempty"`
	CurrentUrl     string    `json:"currentUrl,omitempty"`
	DeepLink       string    `json:"deepLink,omitempty"`
	ClientIp       string    `json:"clientIp,omitempty"`
	BaseUrl        string    `json:"baseUrl,omitempty"`
	UserAgent      string    `json:"userAgent,omitempty"`
	Bot            int       `json:"bot,omitempty"`
	BrowserName    string    `json:"browserName,omitempty"`
	BrowserVersion string    `json:"browserVersion,omitempty"`
	Localization   string    `json:"localization,omitempty"`
	Mobile         int       `json:"mobile,omitempty"`
	Model          string    `json:"model,omitempty"`
	Mozilla        string    `json:"mozilla,omitempty"`
	Os             string    `json:"os,omitempty"`
	FullName       string    `json:"fullName,omitempty"`
	Name           string    `json:"name,omitempty"`
	Version        string    `json:"version,omitempty"`
	Platform       string    `json:"platform,omitempty"`
	Headers        string    `json:"headers,omitempty"`
	Referer        string    `json:"referer,omitempty"`
}

package model

import "time"

type BranchEvent struct {
	Id string
	Timestamp time.Time
	InsertTimestamp time.Time
	EventTimestamp time.Time
	Datasource string
	Origin string
	DeepLinked bool
	Name string
	Attributed bool
	ExistingUser bool
	FirstEventForUser bool
	DiMatchClickToken int
	SecondsFromInstallToEvent int
	SecondsFromLastAttributedTouchToEvent int
	MinutesFromLastAttributedTouchToEvent int
	HoursFromLastAttributedTouchToEvent int
	DaysFromLastAttributedTouchToEvent int
	LastAttributedTouchTimestamp time.Time
	LastAttributedTouchType string
	TransactionId string
	Revenue float64
	EventData JSONB `sql:"type:jsonb"`
	Platform string
	Aaid string
	GeoRegionEn string
	GeoCityEn string
	GeoCountryEn string
	UserData JSONB `sql:"type:jsonb"`
	Campaign string
	CampaignCode string
	Channel string
	Medium string
	MarketingTitle string
	Url string
	AndroidDeeplinkPath string
	AndroidUrl string
	IosUrl string
	LastAttributedTouchData JSONB `sql:"type:jsonb"`
	InstallActivity JSONB `sql:"type:jsonb"`
	CustomData JSONB `sql:"type:jsonb"`
	ContentItems JSONB `sql:"type:jsonb"`
	ReengagementActivity JSONB `sql:"type:jsonb"`
	CompleteEvent JSONB `sql:"type:jsonb"`
}
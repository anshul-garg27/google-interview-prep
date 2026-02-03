package model

import "time"

type Event struct {
	EventId string
	EventName string
	EventTimestamp time.Time
	InsertTimestamp time.Time

	SessionId string
	BbDeviceId int64
	UserId int64
	DeviceId string
	AndroidAdvertisingId string
	ClientId string
	Channel string
	Os string
	ClientType string
	AppLanguage string
	AppVersion string
	AppVersionCode string
	CurrentURL string
	MerchantId string
	UtmReferrer string
	UtmPlatform string
	UtmSource string
	UtmMedium string
	UtmCampaign string
	PpId string


	EventParams JSONB `sql:"type:jsonb"`
}
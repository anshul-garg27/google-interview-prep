package sinker

import (
	"encoding/json"
	"github.com/mssola/useragent"
	"github.com/streadway/amqp"
	"init.bulbul.tv/bulbul-backend/event-grpc/clickhouse"
	"init.bulbul.tv/bulbul-backend/event-grpc/config"
	"init.bulbul.tv/bulbul-backend/event-grpc/model"
	"log"
	"strconv"
	"time"
)

type LaunchReferEvent struct {
	EventId            string              `json:"eventId,omitempty"`
	Headers            map[string][]string `json:"headers,omitempty"`
	ClientIp           string              `json:"clientIp,omitempty"`
	UserAgent          string              `json:"userAgent,omitempty"`
	Timestamp          int64               `json:"timestamp,omitempty"`
	UserId             *int                `json:"userId,omitempty"`
	DeviceId           string              `json:"deviceId,omitempty"`
	BbDeviceId         string              `json:"bbDeviceId,omitempty"`
	UserAccountId      *int                `json:"userAccountId,omitempty"`
	ClientId           string              `json:"clientId,omitempty"`
	ReferralIdentifier string              `json:"referralIdentifier,omitempty"`
	DynamicLink        string              `json:"dynamicLink,omitempty"`
	DeepLink           *string             `json:"deepLink,omitempty"`
}

func SinkLaunchReferEventToClickhouse(delivery amqp.Delivery) bool {
	launchReferEvent := &LaunchReferEvent{}
	if err := json.Unmarshal(delivery.Body, launchReferEvent); err == nil {

		event := TransformLaunchReferEventToEventModel(launchReferEvent)

		if result := clickhouse.Clickhouse(config.New(), nil).Create(event); result != nil && result.Error == nil {
			log.Printf("Published successfully: %v", result.RowsAffected)
			return true
		} else if result == nil {
			log.Printf("Published failed: %v %v", result, event)
		} else {
			log.Printf("Published failed: %v %v", result.Error, event)
		}
	}

	return false
}

func TransformLaunchReferEventToEventModel(launchReferEvent *LaunchReferEvent) model.LaunchReferEvent {

	event := model.LaunchReferEvent{
		Id:                 launchReferEvent.EventId,
		EventTimestamp:     time.Unix(launchReferEvent.Timestamp, 0),
		InsertTimestamp:    time.Now(),
		ClientIp:           launchReferEvent.ClientIp,
		UserAgent:          launchReferEvent.UserAgent,
		DeviceId:           launchReferEvent.DeviceId,
		BbDeviceId:         launchReferEvent.BbDeviceId,
		ClientId:           launchReferEvent.ClientId,
		ReferralIdentifier: launchReferEvent.ReferralIdentifier,
		DynamicLink:        launchReferEvent.DynamicLink,
		DeepLink:           launchReferEvent.DeepLink,
	}

	if launchReferEvent.UserId != nil {
		userId := strconv.Itoa(*launchReferEvent.UserId)
		event.UserId = &userId
	}

	if launchReferEvent.UserAccountId != nil {
		userAccountId := strconv.Itoa(*launchReferEvent.UserAccountId)
		event.UserAccountId = &userAccountId
	}

	if ua := useragent.New(launchReferEvent.UserAgent); ua != nil {
		event.Mozilla = ua.Mozilla()
		event.Platform = ua.Platform()
		event.Os = ua.OS()
		event.Localization = ua.Localization()
		event.Model = ua.Model()
		if ua.Bot() {
			event.Bot = 1
		}
		if ua.Mobile() {
			event.Mobile = 1
		}
		event.OsName = ua.OSInfo().Name
		event.OsFullName = ua.OSInfo().FullName
		event.OsVersion = ua.OSInfo().Version
		event.BrowserName, event.BrowserVersion = ua.Browser()
	}

	return event
}

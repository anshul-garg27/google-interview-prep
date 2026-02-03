package sinker

import (
	"encoding/json"
	"github.com/mssola/useragent"
	"github.com/streadway/amqp"
	"init.bulbul.tv/bulbul-backend/event-grpc/clickhouse"
	"init.bulbul.tv/bulbul-backend/event-grpc/config"
	"init.bulbul.tv/bulbul-backend/event-grpc/model"
	"log"
	"time"
)

type AppInitEvent struct {
	EventId       string              `json:"eventId,omitempty"`
	Headers       map[string][]string `json:"headers,omitempty"`
	ClientIp      string              `json:"clientIp,omitempty"`
	UserAgent     string              `json:"userAgent,omitempty"`
	Timestamp     int64               `json:"timestamp,omitempty"`
	UserId        *string             `json:"userId,omitempty"`
	DeviceId      string              `json:"deviceId,omitempty"`
	BbDeviceId    string              `json:"bbDeviceId,omitempty"`
	UserAccountId *string             `json:"userAccountId,omitempty"`
	ClientId      string              `json:"clientId,omitempty"`
}

func SinkAppInitEventToClickhouse(delivery amqp.Delivery) bool {
	appInitEvent := &AppInitEvent{}
	if err := json.Unmarshal(delivery.Body, appInitEvent); err == nil {

		event := TransformAppInitEventToEventModel(appInitEvent)

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

func TransformAppInitEventToEventModel(appInitEvent *AppInitEvent) model.AppInitEvent {

	event := model.AppInitEvent{
		Id:              appInitEvent.EventId,
		EventTimestamp:  time.Unix(appInitEvent.Timestamp, 0),
		InsertTimestamp: time.Now(),
		ClientIp:        appInitEvent.ClientIp,
		UserAgent:       appInitEvent.UserAgent,
		UserId:          appInitEvent.UserId,
		DeviceId:        appInitEvent.DeviceId,
		BbDeviceId:      appInitEvent.BbDeviceId,
		UserAccountId:   appInitEvent.UserAccountId,
		ClientId:        appInitEvent.ClientId,
	}

	if ua := useragent.New(appInitEvent.UserAgent); ua != nil {
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

	data := map[string]interface{}{}
	if js, err := json.Marshal(appInitEvent); err == nil {
		if err := json.Unmarshal(js, &data); err == nil {
			event.EventData = data
		}
	}
	return event
}

package sinker

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
	"init.bulbul.tv/bulbul-backend/event-grpc/clickhouse"
	"init.bulbul.tv/bulbul-backend/event-grpc/config"
	"init.bulbul.tv/bulbul-backend/event-grpc/header"
	"init.bulbul.tv/bulbul-backend/event-grpc/model"
	"init.bulbul.tv/bulbul-backend/event-grpc/rabbit"
	"log"
	"strconv"
	"time"
)

type UserAccountCreationEvent struct {
	UserAccount UserAccount `json:"newEntry,omitempty"`
}

type UserAccount struct {
	Id             int             `json:"id"`
	CreatedOn      int64           `json:"createdOn"`
	ClientId       string          `json:"clientId"`
	ClientAppType  string          `json:"clientAppType"`
	UserId         int             `json:"userId"`
	SocialAccounts []SocialAccount `json:"socialAccounts"`
}

func ProcessUserAccountCreationEventForReferral(delivery amqp.Delivery) bool {
	if pc, err := rabbit.PopulateContextFromRabbitHeaderTable(delivery.Headers); err == nil {
		userAccountCreationEvent := &UserAccountCreationEvent{}
		if err := json.Unmarshal(delivery.Body, userAccountCreationEvent); err == nil {

			// fixme add time logic and ordering
			var launchReferEvents []model.LaunchReferEvent
			clickhouse.Clickhouse(config.New(), nil).
				Where("device_id = ? OR bb_device_id = ?", pc.XHeader.Get(header.Device), pc.XHeader.Get(header.BBDevice)).
				First(&launchReferEvents)

			if launchReferEvents != nil && len(launchReferEvents) > 0 {
				event := &model.SignupReferralEvent{
					Id:                 uuid.New().String(),
					SignupTimestamp:    time.Unix(userAccountCreationEvent.UserAccount.CreatedOn/1000, 0),
					LaunchTimestamp:    launchReferEvents[0].EventTimestamp,
					UserId:             strconv.Itoa(userAccountCreationEvent.UserAccount.UserId),
					DeviceId:           pc.XHeader.Get(header.Device),
					BbDeviceId:         pc.XHeader.Get(header.BBDevice),
					AccountId:          strconv.Itoa(userAccountCreationEvent.UserAccount.Id),
					ReferralIdentifier: launchReferEvents[0].ReferralIdentifier,
					DynamicLink:        launchReferEvents[0].DynamicLink,
					DeepLink:           launchReferEvents[0].DeepLink,
				}

				if result := clickhouse.Clickhouse(config.New(), nil).Create(event); result != nil && result.Error == nil {
					log.Printf("Published successfully: %v", result.RowsAffected)
					return true
				} else if result == nil {
					log.Printf("Published failed: %v %v", result, event)
				} else {
					log.Printf("Published failed: %v %v", result.Error, event)
				}
			}

			return true
		}
	}

	return false
}

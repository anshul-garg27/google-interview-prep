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

type UserAccountProfileCompletionEvent struct {
	UserAccount UserAccount `json:"newEntry,omitempty"`
}

type SocialAccount struct {
	Platform string                `json:"platform,omitempty"`
	Metrics  *SocialAccountMetrics `json:"metrics,omitempty"`
}

type SocialAccountMetrics struct {
	Followers int `json:"followers,omitempty"`
}

func ProcessUserAccountProfileCompletionEventForReferral(delivery amqp.Delivery) bool {
	if pc, err := rabbit.PopulateContextFromRabbitHeaderTable(delivery.Headers); err == nil {
		userAccountProfileCompletionEvent := &UserAccountProfileCompletionEvent{}
		if err := json.Unmarshal(delivery.Body, userAccountProfileCompletionEvent); err == nil {

			var launchReferEvents []model.LaunchReferEvent
			clickhouse.Clickhouse(config.New(), nil).
				Where("device_id = ? OR bb_device_id = ?", pc.XHeader.Get(header.Device), pc.XHeader.Get(header.BBDevice)).
				First(&launchReferEvents)

			var instaFollowers *int
			if userAccountProfileCompletionEvent.UserAccount.SocialAccounts != nil && len(userAccountProfileCompletionEvent.UserAccount.SocialAccounts) > 0 {
				socialAccounts := userAccountProfileCompletionEvent.UserAccount.SocialAccounts
				for i := 0; i < len(socialAccounts); i++ {
					socialAccount := socialAccounts[i]
					if socialAccount.Platform == "INSTAGRAM" && socialAccount.Metrics != nil {
						instaFollowers = &socialAccount.Metrics.Followers
					}
				}

			}
			if launchReferEvents != nil && len(launchReferEvents) > 0 {
				event := &model.ProfileCompletionReferralEvent{
					Id:                         uuid.New().String(),
					ProfileCompletionTimestamp: time.Now(),
					LaunchTimestamp:            launchReferEvents[0].EventTimestamp,
					UserId:                     strconv.Itoa(userAccountProfileCompletionEvent.UserAccount.UserId),
					DeviceId:                   pc.XHeader.Get(header.Device),
					BbDeviceId:                 pc.XHeader.Get(header.BBDevice),
					AccountId:                  strconv.Itoa(userAccountProfileCompletionEvent.UserAccount.Id),
					ReferralIdentifier:         launchReferEvents[0].ReferralIdentifier,
					DynamicLink:                launchReferEvents[0].DynamicLink,
					DeepLink:                   launchReferEvents[0].DeepLink,
					InstaFollowers:             instaFollowers,
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

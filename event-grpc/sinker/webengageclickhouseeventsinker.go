package sinker

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
	"init.bulbul.tv/bulbul-backend/event-grpc/clickhouse"
	"init.bulbul.tv/bulbul-backend/event-grpc/client/input"
	"init.bulbul.tv/bulbul-backend/event-grpc/config"
	"init.bulbul.tv/bulbul-backend/event-grpc/model"
	"log"
	"time"
)

func SinkWebengageEventToClickhouse(delivery amqp.Delivery) bool {

	webengageEvent := input.WebengageEvent{}
	if err := json.Unmarshal(delivery.Body, &webengageEvent); err == nil {

		t, err := time.Parse("2006-01-02T15:04:05Z0700", webengageEvent.EventTime)
		if err != nil {
			fmt.Errorf("Some error publishing webengage event to ch, invalid time: %v\n", err)
			return false
		}

		event := &model.WebengageEvent{
			Id:          uuid.New().String(),
			EventName:   webengageEvent.EventName,
			AccountType: "BRAND",
			Timestamp:   t,
			Data:        webengageEvent.EventData,
		}

		if webengageEvent.UserId != nil {
			event.UserId = *webengageEvent.UserId
		}

		if webengageEvent.AnonymousId != nil {
			event.AnonymousId = *webengageEvent.AnonymousId
		}

		if result := clickhouse.Clickhouse(config.New(), nil).Create(event); result != nil && result.Error == nil {
			log.Printf("Created successfully: %v", result.RowsAffected)
			return true
		} else if result == nil {
			log.Printf("Creation failed: %v", result)
		} else {
			log.Printf("Creation failed: %v", result.Error)
		}
	}

	return false
}

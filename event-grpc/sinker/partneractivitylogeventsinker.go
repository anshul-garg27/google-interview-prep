package sinker

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/streadway/amqp"
	"init.bulbul.tv/bulbul-backend/event-grpc/clickhouse"
	"init.bulbul.tv/bulbul-backend/event-grpc/config"
	"init.bulbul.tv/bulbul-backend/event-grpc/model"
	"init.bulbul.tv/bulbul-backend/event-grpc/sinker/parser"
)

func BufferPartnerActivityLogEvents(delivery amqp.Delivery, eventBuffer chan interface{}) bool {
	event := parser.GetPartnerActivityLogEvent(delivery.Body)
	if event != nil {
		eventBuffer <- *event
		return true
	}
	return false
}

func PartnerActivityLogEventsSinker(channel chan interface{}) {
	var buffer []model.PartnerActivityLogEvent
	ticker := time.NewTicker(1 * time.Minute)
	BufferLimit, _ := strconv.Atoi(os.Getenv("SCRAPE_LOG_BUFFER_LIMIT"))
	for {
		select {
		case v := <-channel:
			buffer = append(buffer, v.(model.PartnerActivityLogEvent))
			if len(buffer) >= BufferLimit {
				FlushPartnerActivityLogEvents(buffer)
				buffer = []model.PartnerActivityLogEvent{}
			}
		case _ = <-ticker.C:
			FlushPartnerActivityLogEvents(buffer)
			buffer = []model.PartnerActivityLogEvent{}
		}
	}
}

func FlushPartnerActivityLogEvents(events []model.PartnerActivityLogEvent) {
	log.Printf("Creating %v partner-activity log events in batches", len(events))
	if events == nil || len(events) == 0 {
		log.Printf("partner-activity log events empty")
		return
	}
	result := clickhouse.Clickhouse(config.New(), nil).Create(events)
	if result != nil && result.Error == nil {
		log.Printf("Created partner-activity log event successfully: %v", result.RowsAffected)
	} else if result == nil {
		log.Printf("partner-activity log event creation failed: %v", result)
	} else {
		log.Printf("partner-activity log event creation failed: %v", result.Error)
	}
}

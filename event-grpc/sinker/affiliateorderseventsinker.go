package sinker

import (
	"log"
	"time"

	"github.com/streadway/amqp"
	"init.bulbul.tv/bulbul-backend/event-grpc/clickhouse"
	"init.bulbul.tv/bulbul-backend/event-grpc/config"
	"init.bulbul.tv/bulbul-backend/event-grpc/model"
	"init.bulbul.tv/bulbul-backend/event-grpc/sinker/parser"
)

func BufferAffiliateOrderEvents(delivery amqp.Delivery, eventBuffer chan interface{}) bool {
	event := parser.GetAffiliateOrderEvent(delivery.Body)
	if event != nil {
		eventBuffer <- *event
		return true
	}
	return false
}

func AffiliateOrderEventsSinker(affiliateOrdersChannel chan interface{}) {
	var buffer []model.AffiliateOrderEvent
	ticker := time.NewTicker(5 * time.Minute)
	BufferLimit := 100
	for {
		select {
		case v := <-affiliateOrdersChannel:
			buffer = append(buffer, v.(model.AffiliateOrderEvent))
			if len(buffer) >= BufferLimit {
				FlushAffiliateOrderEvents(buffer)
				buffer = []model.AffiliateOrderEvent{}
			}
		case _ = <-ticker.C:
			FlushAffiliateOrderEvents(buffer)
			buffer = []model.AffiliateOrderEvent{}
		}
	}
}

func FlushAffiliateOrderEvents(events []model.AffiliateOrderEvent) {
	log.Printf("Creating %v affiliate order events in batches", len(events))
	if events == nil || len(events) == 0 {
		return
	}
	result := clickhouse.Clickhouse(config.New(), nil).Create(events)
	if result != nil && result.Error == nil {
		log.Printf("Created affiliate order events successfully: %v", result.RowsAffected)
	} else if result == nil {
		log.Printf("Affiliate order events creation failed: %v", result)
	} else {
		log.Printf("Affiliate order events creation failed: %v", result.Error)
	}
}

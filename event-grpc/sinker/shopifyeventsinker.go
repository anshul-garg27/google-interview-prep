package sinker

import (
	"github.com/google/uuid"
	"init.bulbul.tv/bulbul-backend/event-grpc/client/input"
	"init.bulbul.tv/bulbul-backend/event-grpc/model"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/streadway/amqp"
	"init.bulbul.tv/bulbul-backend/event-grpc/clickhouse"
	"init.bulbul.tv/bulbul-backend/event-grpc/config"
	"init.bulbul.tv/bulbul-backend/event-grpc/sinker/parser"
)

func BufferShopifyEvents(delivery amqp.Delivery, eventBuffer chan interface{}) bool {
	event := parser.GetShopifyEvent(delivery.Body)
	if event != nil {
		eventBuffer <- *event
		return true
	}
	return true
}

func ShopifyEventsSinker(channel chan interface{}) {
	var buffer []input.ShopifyEvent
	ticker := time.NewTicker(1 * time.Minute)
	BufferLimit, _ := strconv.Atoi(os.Getenv("SCRAPE_LOG_BUFFER_LIMIT"))
	for {
		select {
		case v := <-channel:
			buffer = append(buffer, v.(input.ShopifyEvent))
			if len(buffer) >= BufferLimit {
				FlushShopifyEvents(buffer)
				buffer = []input.ShopifyEvent{}
			}
		case _ = <-ticker.C:
			FlushShopifyEvents(buffer)
			buffer = []input.ShopifyEvent{}
		}
	}
}

func FlushShopifyEvents(events []input.ShopifyEvent) {
	log.Printf("Creating %v shopify log events in batches", len(events))
	if events == nil || len(events) == 0 {
		log.Printf("Shopify log events empty")
		return
	}
	var records []model.ShopifyEvent
	for i, _ := range events {
		record := model.ShopifyEvent{
			Id:              uuid.New().String(),
			EventName:       events[i].EventName,
			StoreUrl:        events[i].StoreUrl,
			EventTimestamp:  time.Unix(0, events[i].EventTime*int64(time.Millisecond)),
			InsertTimestamp: time.Now(),
			Data:            events[i].EventData,
		}
		records = append(records, record)
	}
	result := clickhouse.Clickhouse(config.New(), nil).Create(records)
	if result != nil && result.Error == nil {
		log.Printf("Created Shopfy event successfully: %v", result.RowsAffected)
	} else if result == nil {
		log.Printf("Shopfy event creation failed: %v", result)
	} else {
		log.Printf("Shopfy event creation failed: %v", result.Error)
	}
}

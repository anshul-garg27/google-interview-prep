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

func BufferPostLogEvents(delivery amqp.Delivery, eventBuffer chan interface{}) bool {
	event := parser.GetPostLogEvent(delivery.Body)
	if event != nil {
		eventBuffer <- *event
		return true
	}
	return true
}

func BufferSentimentLogEvents(delivery amqp.Delivery, eventBuffer chan interface{}) bool {
	event := parser.GetSentimentLogEvent(delivery.Body)
	if event != nil {
		eventBuffer <- *event
		return true
	}
	return false
}

func BufferPostActivityLogEvents(delivery amqp.Delivery, eventBuffer chan interface{}) bool {
	event := parser.GetPostActivityLogEvent(delivery.Body)
	if event != nil {
		eventBuffer <- *event
		return true
	}
	return false
}

func BufferProfileLogEvents(delivery amqp.Delivery, eventBuffer chan interface{}) bool {
	event := parser.GetProfileLogEvent(delivery.Body)
	if event != nil {
		eventBuffer <- *event
		return true
	}
	return false
}

func BufferProfileRelationshipLogEvents(delivery amqp.Delivery, eventBuffer chan interface{}) bool {
	event := parser.GetProfileRelationshipLogEvent(delivery.Body)
	if event != nil {
		eventBuffer <- *event
		return true
	}
	return false
}

func BufferScrapeRequestLogEvents(delivery amqp.Delivery, eventBuffer chan interface{}) bool {
	event := parser.GetScrapeRequestLogEvent(delivery.Body)
	if event != nil {
		eventBuffer <- *event
		return true
	}
	return false
}

func BufferOrderLogEvents(delivery amqp.Delivery, eventBuffer chan interface{}) bool {
	event := parser.GetOrderLogEvent(delivery.Body)
	if event != nil {
		eventBuffer <- *event
		return true
	}
	return false
}

func PostLogEventsSinker(channel chan interface{}) {
	var buffer []model.PostLogEvent
	ticker := time.NewTicker(5 * time.Minute)
	BufferLimit, _ := strconv.Atoi(os.Getenv("SCRAPE_LOG_BUFFER_LIMIT"))
	for {
		select {
		case v := <-channel:
			buffer = append(buffer, v.(model.PostLogEvent))
			if len(buffer) >= BufferLimit {
				FlushPostLogEvents(buffer)
				buffer = []model.PostLogEvent{}
			}
		case _ = <-ticker.C:
			FlushPostLogEvents(buffer)
			buffer = []model.PostLogEvent{}
		}
	}
}

func PostActivityLogEventsSinker(channel chan interface{}) {
	var buffer []model.PostActivityLogEvent
	ticker := time.NewTicker(1 * time.Minute)
	BufferLimit, _ := strconv.Atoi(os.Getenv("SCRAPE_LOG_BUFFER_LIMIT"))
	for {
		select {
		case v := <-channel:
			buffer = append(buffer, v.(model.PostActivityLogEvent))
			if len(buffer) >= BufferLimit {
				FlushPostActivityLogEvents(buffer)
				buffer = []model.PostActivityLogEvent{}
			}
		case _ = <-ticker.C:
			FlushPostActivityLogEvents(buffer)
			buffer = []model.PostActivityLogEvent{}
		}
	}
}

func SentimentLogEventsSinker(channel chan interface{}) {
	var buffer []model.SentimentLogEvent
	ticker := time.NewTicker(1 * time.Minute)
	BufferLimit, _ := strconv.Atoi(os.Getenv("SCRAPE_LOG_BUFFER_LIMIT"))
	log.Println(BufferLimit)
	for {
		select {
		case v := <-channel:
			buffer = append(buffer, v.(model.SentimentLogEvent))
			log.Println(len(buffer))
			if len(buffer) >= BufferLimit {
				log.Println(buffer)
				FlushSentimentLogEvents(buffer)
				buffer = []model.SentimentLogEvent{}
			}
		case _ = <-ticker.C:
			log.Println(buffer)
			FlushSentimentLogEvents(buffer)
			buffer = []model.SentimentLogEvent{}
		}
	}
}

func ProfileLogEventsSinker(channel chan interface{}) {
	var buffer []model.ProfileLogEvent
	ticker := time.NewTicker(1 * time.Minute)
	BufferLimit, _ := strconv.Atoi(os.Getenv("SCRAPE_LOG_BUFFER_LIMIT"))
	for {
		select {
		case v := <-channel:
			buffer = append(buffer, v.(model.ProfileLogEvent))
			if len(buffer) >= BufferLimit {
				FlushProfileLogEvents(buffer)
				buffer = []model.ProfileLogEvent{}
			}
		case _ = <-ticker.C:
			FlushProfileLogEvents(buffer)
			buffer = []model.ProfileLogEvent{}
		}
	}
}

func ProfileRelationshipLogEventsSinker(channel chan interface{}) {
	var buffer []model.ProfileRelationshipLogEvent
	ticker := time.NewTicker(1 * time.Minute)
	BufferLimit, _ := strconv.Atoi(os.Getenv("SCRAPE_LOG_BUFFER_LIMIT"))
	for {
		select {
		case v := <-channel:
			buffer = append(buffer, v.(model.ProfileRelationshipLogEvent))
			if len(buffer) >= BufferLimit {
				FlushProfileRelationshipLogEvents(buffer)
				buffer = []model.ProfileRelationshipLogEvent{}
			}
		case _ = <-ticker.C:
			FlushProfileRelationshipLogEvents(buffer)
			buffer = []model.ProfileRelationshipLogEvent{}
		}
	}
}

func ScrapeRequestLogEventsSinker(channel chan interface{}) {
	var buffer []model.ScrapeRequestLogEvent
	ticker := time.NewTicker(1 * time.Minute)
	BufferLimit, _ := strconv.Atoi(os.Getenv("SCRAPE_LOG_BUFFER_LIMIT"))
	for {
		select {
		case v := <-channel:
			buffer = append(buffer, v.(model.ScrapeRequestLogEvent))
			if len(buffer) >= BufferLimit {
				FlushScrapeRequestEvents(buffer)
				buffer = []model.ScrapeRequestLogEvent{}
			}
		case _ = <-ticker.C:
			FlushScrapeRequestEvents(buffer)
			buffer = []model.ScrapeRequestLogEvent{}
		}
	}
}

func OrderLogEventsSinker(channel chan interface{}) {
	var buffer []model.OrderLogEvent
	ticker := time.NewTicker(1 * time.Minute)
	BufferLimit, _ := strconv.Atoi(os.Getenv("SCRAPE_LOG_BUFFER_LIMIT"))
	for {
		select {
		case v := <-channel:
			buffer = append(buffer, v.(model.OrderLogEvent))
			if len(buffer) >= BufferLimit {
				FlushOrderLogEvents(buffer)
				buffer = []model.OrderLogEvent{}
			}
		case _ = <-ticker.C:
			FlushOrderLogEvents(buffer)
			buffer = []model.OrderLogEvent{}
		}
	}
}

func FlushPostLogEvents(events []model.PostLogEvent) {
	log.Printf("Creating %v post log events in batches", len(events))
	if events == nil || len(events) == 0 {
		log.Printf("Post log events empty")
		return
	}
	result := clickhouse.Clickhouse(config.New(), nil).Create(events)
	if result != nil && result.Error == nil {
		log.Printf("Created post log event successfully: %v", result.RowsAffected)
	} else if result == nil {
		log.Printf("Post log creation failed: %v", result)
	} else {
		log.Printf("Post log creation failed: %v", result.Error)
	}
}

func FlushPostActivityLogEvents(events []model.PostActivityLogEvent) {
	log.Printf("Creating %v post activity log events in batches", len(events))
	if events == nil || len(events) == 0 {
		log.Printf("Post activity log events empty")
		return
	}
	result := clickhouse.Clickhouse(config.New(), nil).Create(events)
	if result != nil && result.Error == nil {
		log.Printf("Created post activity log event successfully: %v", result.RowsAffected)
	} else if result == nil {
		log.Printf("Post activity log creation failed: %v", result)
	} else {
		log.Printf("Post activity log creation failed: %v", result.Error)
	}
}

func FlushSentimentLogEvents(events []model.SentimentLogEvent) {
	log.Printf("Creating %v sentiment log events in batches", len(events))
	if events == nil || len(events) == 0 {
		log.Printf("Sentiment log events empty")
		return
	}
	result := clickhouse.Clickhouse(config.New(), nil).Create(events)
	log.Println(result)
	if result != nil && result.Error == nil {
		log.Printf("Created sentiment log event successfully: %v", result.RowsAffected)
	} else if result == nil {
		log.Printf("Sentiment log creation failed: %v", result)
	} else {
		log.Printf("Sentiment log creation failed: %v", result.Error)
	}
}

func FlushProfileLogEvents(events []model.ProfileLogEvent) {
	log.Printf("Creating %v profile log events in batches", len(events))
	if events == nil || len(events) == 0 {
		log.Printf("Profile log events empty")
		return
	}
	result := clickhouse.Clickhouse(config.New(), nil).Create(events)
	if result != nil && result.Error == nil {
		log.Printf("Created profile log event successfully: %v", result.RowsAffected)
	} else if result == nil {
		log.Printf("Profile log event creation failed: %v", result)
	} else {
		log.Printf("Profile log event creation failed: %v", result.Error)
	}
}

func FlushProfileRelationshipLogEvents(events []model.ProfileRelationshipLogEvent) {
	log.Printf("Creating %v profile Relationship log events in batches", len(events))
	if events == nil || len(events) == 0 {
		log.Printf("Profile Relationship log events empty")
		return
	}
	result := clickhouse.Clickhouse(config.New(), nil).Create(events)
	if result != nil && result.Error == nil {
		log.Printf("Created Profile Relationship log event successfully: %v", result.RowsAffected)
	} else if result == nil {
		log.Printf("Profile Relationship log event creation failed: %v", result)
	} else {
		log.Printf("Profile Relationship log event creation failed: %v", result.Error)
	}
}

func FlushScrapeRequestEvents(events []model.ScrapeRequestLogEvent) {
	log.Printf("Creating %v scrape request log events in batches", len(events))
	if events == nil || len(events) == 0 {
		log.Printf("Scrape log events empty")
		return
	}
	result := clickhouse.Clickhouse(config.New(), nil).Create(events)
	if result != nil && result.Error == nil {
		log.Printf("Created scrape request log event successfully: %v", result.RowsAffected)
	} else if result == nil {
		log.Printf("Scrape request log event creation failed: %v", result)
	} else {
		log.Printf("Scrape request log event creation failed: %v", result.Error)
	}
}

func FlushOrderLogEvents(events []model.OrderLogEvent) {
	log.Printf("Creating %v order log events in batches", len(events))
	if events == nil || len(events) == 0 {
		log.Printf("Order log events empty")
		return
	}
	result := clickhouse.Clickhouse(config.New(), nil).Create(events)
	if result != nil && result.Error == nil {
		log.Printf("Created order log event successfully: %v", result.RowsAffected)
	} else if result == nil {
		log.Printf("Order log event creation failed: %v", result)
	} else {
		log.Printf("Order log event creation failed: %v", result.Error)
	}
}

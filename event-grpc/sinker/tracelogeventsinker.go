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

func BufferTraceLogEvent(delivery amqp.Delivery, eventBuffer chan interface{}) bool {
	event := parser.GetTraceLogEvent(delivery.Body)
	if event != nil {
		eventBuffer <- *event
		return true
	}
	return false
}

func TraceLogEventsSinker(traceLogsChannel chan interface{}) {
	var buffer []model.TraceLogEvent
	ticker := time.NewTicker(5 * time.Minute)
	BufferLimit := 100
	for {
		select {
		case v := <-traceLogsChannel:
			buffer = append(buffer, v.(model.TraceLogEvent))
			if len(buffer) >= BufferLimit {
				FlushTraceLogEvents(buffer)
				buffer = []model.TraceLogEvent{}
			}
		case _ = <-ticker.C:
			FlushTraceLogEvents(buffer)
			buffer = []model.TraceLogEvent{}
		}
	}
}

func FlushTraceLogEvents(events []model.TraceLogEvent) {
	log.Printf("Creating %v trace log events in batches", len(events))
	if events == nil || len(events) == 0 {
		log.Printf("Trace log events empty")
		return
	}
	result := clickhouse.Clickhouse(config.New(), nil).Create(events)
	if result != nil && result.Error == nil {
		log.Printf("Created trace log events successfully: %v", result.RowsAffected)
	} else if result == nil {
		log.Printf("Trace log events creation failed: %v", result)
	} else {
		log.Printf("Trace log events creation failed: %v", result.Error)
	}
}

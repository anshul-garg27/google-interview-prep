package sinker

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"google.golang.org/protobuf/encoding/protojson"
	"init.bulbul.tv/bulbul-backend/event-grpc/config"
	"init.bulbul.tv/bulbul-backend/event-grpc/pg"
	"init.bulbul.tv/bulbul-backend/event-grpc/proto/go/bulbulgrpc"
	"log"
)

func SinkClickEventToPG(delivery amqp.Delivery) bool {
	clickEvent := &ClickEvent{}
	if err := json.Unmarshal(delivery.Body, clickEvent); err == nil {

		event := TransformClickEventToEventModel(clickEvent)

		if result := pg.PG(config.New()).Create(event); result != nil && result.Error == nil {
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

func SinkEventToPG(delivery amqp.Delivery) bool {

	grpcEvent := &bulbulgrpc.Events{}
	if err := protojson.Unmarshal(delivery.Body, grpcEvent); err == nil {
		// iterate through all events and sink each to event

		if grpcEvent.Events == nil || len(grpcEvent.Events) == 0 {
			log.Println("No event in events: ", grpcEvent)
			return false
		}

		log.Printf("Received in q: %v", grpcEvent)

		for _, e := range grpcEvent.Events {
			event, err := TransformToEventModel(grpcEvent, e)
			if err != nil {
				return false
			}

			if result := pg.PG(config.New()).Create(event); result != nil && result.Error == nil {
				log.Printf("Published successfully: %v", result.RowsAffected)
			} else if result == nil {
				log.Printf("Published failed: %v %v", result, event)
			} else {
				log.Printf("Published failed: %v %v", result.Error, event)
			}
		}

	} else {
		log.Println("Some error processing event %v", err)
	}

	return true
}

func SinkErrorEventToPG(delivery amqp.Delivery) bool {


	grpcEvent := &bulbulgrpc.Events{}
	if err := protojson.Unmarshal(delivery.Body, grpcEvent); err == nil {
		// iterate through all events and sink each to event

		if grpcEvent.Events == nil || len(grpcEvent.Events) == 0 {
			log.Println("No event in events: ", grpcEvent)
			return false
		}

		for _, e := range grpcEvent.Events {

			errorEvent, err := TransformToErrorEventModel(grpcEvent, e)
			if err != nil {
				return false
			}

			if result := pg.PG(config.New()).Create(errorEvent); result != nil && result.Error == nil {
				log.Printf("Published successfully: %v", result.RowsAffected)
				return true
			} else if result == nil {
				log.Printf("Published failed: %v %v", result, errorEvent)
			} else {
				log.Printf("Published failed: %v %v", result.Error, errorEvent)
			}
		}
	} else {
		log.Println("Some error processing event %v", err)
	}

	return false
}
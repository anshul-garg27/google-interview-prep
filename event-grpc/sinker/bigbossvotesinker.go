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

func BufferBigBossVotes(delivery amqp.Delivery, eventBuffer chan interface{}) bool {
	event := parser.GetBigBossVoteLog(delivery.Body)
	if event != nil {
		eventBuffer <- *event
		return true
	}
	return false
}

func BigBossVotesSinker(bigBossVotesChannel chan interface{}) {
	var buffer []model.BigBossVoteLog
	ticker := time.NewTicker(5 * time.Minute)
	BufferLimit := 100
	for {
		select {
		case v := <-bigBossVotesChannel:
			buffer = append(buffer, v.(model.BigBossVoteLog))
			if len(buffer) >= BufferLimit {
				FlushBigBossVotes(buffer)
				buffer = []model.BigBossVoteLog{}
			}
		case _ = <-ticker.C:
			FlushBigBossVotes(buffer)
			buffer = []model.BigBossVoteLog{}
		}
	}
}

func FlushBigBossVotes(events []model.BigBossVoteLog) {
	log.Printf("Creating %v big boss votes events in batches", len(events))
	result := clickhouse.Clickhouse(config.New(), nil).Create(events)
	if result != nil && result.Error == nil {
		log.Printf("Created big boss votes events successfully: %v", result.RowsAffected)
	} else if result == nil {
		log.Printf("Big boss votes events creation failed: %v", result)
	} else {
		log.Printf("Big boss votes events creation failed: %v", result.Error)
	}
}

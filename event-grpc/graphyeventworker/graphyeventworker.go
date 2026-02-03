package graphyeventworker

import (
	"encoding/json"
	"init.bulbul.tv/bulbul-backend/event-grpc/config"
	"init.bulbul.tv/bulbul-backend/event-grpc/rabbit"
	"init.bulbul.tv/bulbul-backend/event-grpc/safego"
	"log"
	"strings"
	"sync"
)

var (
	eventWrapperChannel chan map[string]interface{}
	channelInit         sync.Once
)

func GetChannel(config config.Config) chan map[string]interface{}{
	channelInit.Do(func() {
		eventWrapperChannel = make(chan map[string]interface{}, config.GRAPHY_EVENT_WORKER_POOL_CONFIG.EVENT_BUFFERED_CHANNEL_SIZE)
		initWorkerPool(config, config.GRAPHY_EVENT_WORKER_POOL_CONFIG.EVENT_WORKER_POOL_SIZE, eventWrapperChannel, config.GRAPHY_EVENT_SINK_CONFIG)
	})
	return eventWrapperChannel
}

func initWorkerPool(config config.Config, workerPoolSize int, eventChannel <-chan map[string]interface{}, eventSinkConfig *config.EVENT_SINK_CONFIG) {
	for i := 0; i < workerPoolSize; i++ {
		safego.GoNoCtx(func() {
			worker(config, eventChannel)
		})

	}
}

func worker(config config.Config, eventChannel <-chan map[string]interface{}) {
	for e := range eventChannel {
		rabbitConn := rabbit.Rabbit(config)

		if b, err := json.Marshal(e); err != nil || b == nil || e["name"] == nil {
			log.Printf("Some error publishing message, empty event: %v", err)
		} else {
			err := rabbitConn.Publish(config.GRAPHY_EVENT_SINK_CONFIG.EVENT_EXCHANGE, strings.ToUpper(e["name"].(string)), b, map[string]interface{}{})
			if err != nil {
				log.Printf("Some error publishing message: %v", err)
			}
		}
	}
}
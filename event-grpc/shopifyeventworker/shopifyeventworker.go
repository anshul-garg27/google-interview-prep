package shopifyeventworker

import (
	"encoding/json"
	"init.bulbul.tv/bulbul-backend/event-grpc/client/input"
	"init.bulbul.tv/bulbul-backend/event-grpc/config"
	"init.bulbul.tv/bulbul-backend/event-grpc/rabbit"
	"init.bulbul.tv/bulbul-backend/event-grpc/safego"
	"log"
	"sync"
)

var (
	eventWrapperChannel chan input.ShopifyEvent
	channelInit         sync.Once
)

func GetChannel(config config.Config) chan input.ShopifyEvent {
	channelInit.Do(func() {
		eventWrapperChannel = make(chan input.ShopifyEvent, config.SHOPIFY_EVENT_WORKER_POOL_CONFIG.EVENT_BUFFERED_CHANNEL_SIZE)
		initWorkerPool(config, config.SHOPIFY_EVENT_WORKER_POOL_CONFIG.EVENT_WORKER_POOL_SIZE, eventWrapperChannel, config.SHOPIFY_EVENT_SINK_CONFIG)
	})
	return eventWrapperChannel
}

func initWorkerPool(config config.Config, workerPoolSize int, eventChannel <-chan input.ShopifyEvent, eventSinkConfig *config.EVENT_SINK_CONFIG) {
	for i := 0; i < workerPoolSize; i++ {
		safego.GoNoCtx(func() {
			worker(config, eventChannel)
		})
	}
}

func worker(config config.Config, eventChannel <-chan input.ShopifyEvent) {
	for e := range eventChannel {
		rabbitConn := rabbit.Rabbit(config)

		if b, err := json.Marshal(e); err != nil || b == nil {
			log.Printf("Some error publishing message, empty event: %v", err)
		} else {
			err := rabbitConn.Publish(config.SHOPIFY_EVENT_SINK_CONFIG.EVENT_EXCHANGE, "shopify_events", b, map[string]interface{}{})
			if err != nil {
				log.Printf("Some error publishing message: %v", err)
			}
		}
	}
}

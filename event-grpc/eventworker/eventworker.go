package eventworker

import (
	"fmt"
	"google.golang.org/protobuf/encoding/protojson"
	"init.bulbul.tv/bulbul-backend/event-grpc/config"
	"init.bulbul.tv/bulbul-backend/event-grpc/proto/go/bulbulgrpc"
	"init.bulbul.tv/bulbul-backend/event-grpc/rabbit"
	"init.bulbul.tv/bulbul-backend/event-grpc/safego"
	"init.bulbul.tv/bulbul-backend/event-grpc/transformer"
	"log"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	eventWrapperChannel chan bulbulgrpc.Events
	channelInit         sync.Once
)

func GetChannel(config config.Config) chan bulbulgrpc.Events{
	channelInit.Do(func() {
		eventWrapperChannel = make(chan bulbulgrpc.Events, config.EVENT_WORKER_POOL_CONFIG.EVENT_BUFFERED_CHANNEL_SIZE)
		initWorkerPool(config, config.EVENT_WORKER_POOL_CONFIG.EVENT_WORKER_POOL_SIZE, eventWrapperChannel, config.EVENT_SINK_CONFIG)
	})
	return eventWrapperChannel
}

func initWorkerPool(config config.Config, workerPoolSize int, eventChannel <-chan bulbulgrpc.Events, eventSinkConfig *config.EVENT_SINK_CONFIG) {
	for i := 0; i < workerPoolSize; i++ {
		safego.GoNoCtx(func() {
			worker(config, eventChannel)
		})

	}
}

func worker(config config.Config, eventChannel <-chan bulbulgrpc.Events) {
	for e := range eventChannel {
		rabbitConn := rabbit.Rabbit(config)
		header := e.Header
		eventsGrouped := make(map[string][]*bulbulgrpc.Event)
		for i, _ := range e.Events {
			eventName := fmt.Sprintf("%v", reflect.TypeOf(e.Events[i].GetEventOf()))
			eventName = strings.ReplaceAll(eventName, "*bulbulgrpc.Event_", "")
			et := time.Unix(transformer.InterfaceToInt64(e.Events[i].Timestamp)/1000, 0)
			timeNow := time.Now()
			if et.Unix() > timeNow.Unix() {
				e.Events[i].Timestamp = strconv.FormatInt(timeNow.Unix()*1000, 10)
			}
			eventsGrouped[eventName] = append(eventsGrouped[eventName], e.Events[i])
		}
		for eventName, events := range eventsGrouped {
			e := bulbulgrpc.Events {
				Header: header,
				Events: events,
			}
			if b, err := protojson.Marshal(&e); err != nil || b == nil {
				log.Printf("Some error publishing message, empty event: %v", err)
			} else {
				err := rabbitConn.Publish(config.EVENT_SINK_CONFIG.EVENT_EXCHANGE, eventName, b, map[string]interface{}{})
				if err != nil {
					log.Printf("Some error publishing message: %v", err)
				}
			}
		}
	}
}

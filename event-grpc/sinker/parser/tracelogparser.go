package parser

import (
	"encoding/json"
	"github.com/google/uuid"
	"init.bulbul.tv/bulbul-backend/event-grpc/model"
	"log"
	"time"
)

func GetTraceLogEvent(input []byte) *model.TraceLogEvent {
	rawTraceLogEvent := map[string]interface{}{}
	if err := json.Unmarshal(input, &rawTraceLogEvent); err == nil {

		if rawTraceLogEvent == nil {
			log.Printf("Incomplete event data: %v", rawTraceLogEvent)
			return nil
		}

		event := &model.TraceLogEvent{
			Id:             uuid.New().String(),
			HostName:       getString(rawTraceLogEvent["hostName"]),
			ServiceName:    getString(rawTraceLogEvent["serviceName"]),
			Timestamp:      time.Now(),
			TimeTaken:      int64(getFloat64(rawTraceLogEvent["timeTaken"])),
			Method:         getString(rawTraceLogEvent["method"]),
			URI:            getString(rawTraceLogEvent["uri"]),
			Headers:        getMap(rawTraceLogEvent["headers"]),
			ResponseStatus: int64(getFloat64(rawTraceLogEvent["responseStatus"])),
			RemoteAddress:  getString(rawTraceLogEvent["remoteAddress"]),
			OnlyRequest:    getBool(rawTraceLogEvent["onlyRequest"]),
			RequestBody:    getString(rawTraceLogEvent["requestBody"]),
			ResponseBody:   getString(rawTraceLogEvent["responseBody"]),
		}
		return event
	}
	return nil
}

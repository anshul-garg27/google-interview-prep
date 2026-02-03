package parser

import (
	"encoding/json"
	"log"
	"time"

	"init.bulbul.tv/bulbul-backend/event-grpc/model"
)

func GetPartnerActivityLogEvent(input []byte) *model.PartnerActivityLogEvent {
	rawEvent := map[string]interface{}{}
	event := &model.PartnerActivityLogEvent{}
	requiredKeys := []string{"eventId", "activityName", "incrementCount", "moduleName", "partnerId", "accountId", "activityId",
		"activityMeta", "eventTimestamp"}

	err := json.Unmarshal(input, &rawEvent)
	if err == nil {

		if rawEvent == nil {
			log.Printf("[ERROR] PartnerActivityLogEvent: Incomplete event data: %v", rawEvent)
			return nil
		}

		for _, key := range requiredKeys {
			if value, ok := rawEvent[key]; ok {
				switch key {
				case "eventId":
					event.EventId = getString(value)
				case "activityName":
					event.ActivityName = getString(value)
				case "moduleName":
					event.ModuleName = getString(value)
				case "incrementCount":
					event.IncrementCount = int64(value.(float64))
				case "partnerId":
					event.PartnerId = int64(value.(float64))
				case "accountId":
					event.AccountId = int64(value.(float64))
				case "activityId":
					activityId, ok := value.(float64)
					if !ok {
						event.ActivityId = nil
					} else {
						activityIdInt := int64(activityId)
						event.ActivityId = &activityIdInt
					}
				case "activityMeta":
					activityMetaJSON, ok := value.(string)
					if !ok {
						event.ActivityMeta = make(map[string]interface{})
					} else {
						var activityMeta map[string]interface{}
						err := json.Unmarshal([]byte(activityMetaJSON), &activityMeta)
						if err != nil {
							event.ActivityMeta = make(map[string]interface{})
						} else {
							event.ActivityMeta = activityMeta
						}
					}
				case "eventTimestamp":
					location, _ := time.LoadLocation("Asia/Kolkata")
					eventTimestamp, err := time.ParseInLocation("2006-01-02 15:04:05", getString(value), location)
					if err == nil {
						event.EventTimestamp = eventTimestamp
					}
				}
			} else {
				log.Printf("[ERROR] ProfileLogEvent: Incomplete event data, missing key: %s", key)
				return nil
			}
		}
		event.InsertTimestamp = time.Now()
		return event
	}
	return nil
}

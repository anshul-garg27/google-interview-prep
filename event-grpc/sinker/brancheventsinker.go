package sinker

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"init.bulbul.tv/bulbul-backend/event-grpc/clickhouse"
	"init.bulbul.tv/bulbul-backend/event-grpc/config"
	"init.bulbul.tv/bulbul-backend/event-grpc/model"
	"init.bulbul.tv/bulbul-backend/event-grpc/transformer"
	"log"
	"strings"
	"time"
)

func SinkBranchEventToClickhouse(delivery amqp.Delivery) bool {

	rawBranchEvent := map[string]interface{}{}
	if err := json.Unmarshal(delivery.Body, &rawBranchEvent); err == nil {

		if rawBranchEvent == nil || rawBranchEvent["id"] == nil || rawBranchEvent["name"] == nil || rawBranchEvent["timestamp"] == nil {
			log.Printf("Incomplete event data: %v", rawBranchEvent)
			return false
		}

		event := model.BranchEvent{
			Id:              rawBranchEvent["id"].(string),
			Name:            rawBranchEvent["name"].(string),
			Timestamp:       time.Unix(transformer.InterfaceToInt64(rawBranchEvent["timestamp"].(float64))/1000, 0),
			InsertTimestamp: time.Now(),
			CompleteEvent:   rawBranchEvent,
		}

		if val, ok := rawBranchEvent["event_timestamp"]; ok && val != nil {
			t := time.Unix(transformer.InterfaceToInt64(rawBranchEvent["timestamp"].(float64))/1000, 0)
			event.EventTimestamp = t
		}

		if val, ok := rawBranchEvent["datasource"]; ok && val != nil {
			datasource := fmt.Sprintf("%v", val)
			event.Datasource = datasource
		}

		if val, ok := rawBranchEvent["origin"]; ok && val != nil {
			origin := fmt.Sprintf("%v", val)
			event.Origin = origin
		}

		if val, ok := rawBranchEvent["deep_linked"]; ok && val != nil {
			if deepLinked, ok := val.(bool); ok {
				event.DeepLinked = deepLinked
			}
		}

		if val, ok := rawBranchEvent["attributed"]; ok && val != nil {
			if attributed, ok := val.(bool); ok {
				event.Attributed = attributed
			}
		}

		if val, ok := rawBranchEvent["existing_user"]; ok && val != nil {
			if existingUser, ok := val.(bool); ok {
				event.ExistingUser = existingUser
			}
		}

		if val, ok := rawBranchEvent["first_event_for_user"]; ok && val != nil {
			if firstEventForUser, ok := val.(bool); ok {
				event.FirstEventForUser = firstEventForUser
			}
		}

		if val, ok := rawBranchEvent["di_match_click_token"]; ok && val != nil {
			if diMatchClickToken := transformer.InterfaceToInt(val); diMatchClickToken != 0 {
				event.DiMatchClickToken = diMatchClickToken
			}
		}

		if val, ok := rawBranchEvent["seconds_from_install_to_event"]; ok && val != nil {
			if intVal := transformer.InterfaceToInt(val); intVal != 0 {
				event.SecondsFromInstallToEvent = intVal
			}
		}

		if val, ok := rawBranchEvent["seconds_from_last_attributed_touch_to_event"]; ok && val != nil {
			if intVal := transformer.InterfaceToInt(val); intVal != 0 {
				event.SecondsFromLastAttributedTouchToEvent = intVal
			}
		}

		if val, ok := rawBranchEvent["minutes_from_last_attributed_touch_to_event"]; ok && val != nil {
			if intVal := transformer.InterfaceToInt(val); intVal != 0 {
				event.MinutesFromLastAttributedTouchToEvent = intVal
			}
		}

		if val, ok := rawBranchEvent["hours_from_last_attributed_touch_to_event"]; ok && val != nil {
			if intVal := transformer.InterfaceToInt(val); intVal != 0 {
				event.HoursFromLastAttributedTouchToEvent = intVal
			}
		}

		if val, ok := rawBranchEvent["days_from_last_attributed_touch_to_event"]; ok && val != nil {
			if intVal := transformer.InterfaceToInt(val); intVal != 0 {
				event.DaysFromLastAttributedTouchToEvent = intVal
			}
		}

		if val, ok := rawBranchEvent["last_attributed_touch_timestamp"]; ok && val != nil {
			if intVal := transformer.InterfaceToInt(val); intVal != 0 {
				ts := time.Unix(transformer.InterfaceToInt64(intVal)/1000, 0)
				event.LastAttributedTouchTimestamp = ts
			}
		}

		if val, ok := rawBranchEvent["last_attributed_touch_type"]; ok && val != nil {
			strVal := fmt.Sprintf("%v", val)
			event.LastAttributedTouchType = strVal
		}

		if val, ok := rawBranchEvent["event_data"]; ok && val != nil {
			if mapData, ok := val.(map[string]interface{}); ok {

				if val, ok := mapData["transaction_id"]; ok && val != nil {
					strVal := fmt.Sprintf("%v", val)
					event.TransactionId = strVal
				}

				if val, ok := mapData["revenue"]; ok && val != nil {
					if flatVal, ok := val.(float64); ok {
						event.Revenue = flatVal
					}
				}

				event.EventData = mapData
			}
		}

		if val, ok := rawBranchEvent["user_data"]; ok && val != nil {
			if mapData, ok := val.(map[string]interface{}); ok {

				if val, ok := mapData["platform"]; ok && val != nil {
					strVal := fmt.Sprintf("%v", val)
					event.Platform = strVal
				}

				if val, ok := mapData["aaid"]; ok && val != nil {
					strVal := fmt.Sprintf("%v", val)
					event.Aaid = strVal
				}

				if val, ok := mapData["geo_region_en"]; ok && val != nil {
					strVal := fmt.Sprintf("%v", val)
					event.GeoRegionEn = strVal
				}

				if val, ok := mapData["geo_city_en"]; ok && val != nil {
					strVal := fmt.Sprintf("%v", val)
					event.GeoCityEn = strVal
				}

				if val, ok := mapData["geo_country_en"]; ok && val != nil {
					strVal := fmt.Sprintf("%v", val)
					event.GeoCountryEn = strVal
				}

				event.UserData = mapData
			}
		}

		if val, ok := rawBranchEvent["last_attributed_touch_data"]; ok && val != nil {
			if mapData, ok := val.(map[string]interface{}); ok {

				if val, ok := mapData["~campaign"]; ok && val != nil {
					strVal := fmt.Sprintf("%v", val)
					event.Campaign = strVal
					campaignCode := strings.ToLower(strings.Replace(strVal, " ", "", -1))
					event.CampaignCode = campaignCode
				}

				if val, ok := mapData["~channel"]; ok && val != nil {
					strVal := fmt.Sprintf("%v", val)
					event.Channel = strVal
				}

				if val, ok := mapData["~feature"]; ok && val != nil {
					strVal := fmt.Sprintf("%v", val)
					event.Medium = strVal
				}

				if val, ok := mapData["$marketing_title"]; ok && val != nil {
					strVal := fmt.Sprintf("%v", val)
					event.MarketingTitle = strVal
				}

				if val, ok := mapData["+url"]; ok && val != nil {
					strVal := fmt.Sprintf("%v", val)
					event.Url = strVal
				}

				if val, ok := mapData["$android_deeplink_path"]; ok && val != nil {
					strVal := fmt.Sprintf("%v", val)
					event.AndroidDeeplinkPath = strVal
				}

				if val, ok := mapData["$android_url"]; ok && val != nil {
					strVal := fmt.Sprintf("%v", val)
					event.AndroidUrl = strVal
				}

				if val, ok := mapData["$ios_url"]; ok && val != nil {
					strVal := fmt.Sprintf("%v", val)
					event.IosUrl = strVal
				}

				event.LastAttributedTouchData = mapData

			}
		}

		if val, ok := rawBranchEvent["install_activity"]; ok && val != nil {
			if mapData, ok := val.(map[string]interface{}); ok {
				event.InstallActivity = mapData
			}
		}

		if val, ok := rawBranchEvent["custom_data"]; ok && val != nil {
			if mapData, ok := val.(map[string]interface{}); ok {
				event.CustomData = mapData
			}
		}

		if val, ok := rawBranchEvent["content_items"]; ok && val != nil {
			if mapData, ok := val.(map[string]interface{}); ok {
				event.ContentItems = mapData
			}
		}

		if val, ok := rawBranchEvent["reengagement_activity"]; ok && val != nil {
			if mapData, ok := val.(map[string]interface{}); ok {
				event.ReengagementActivity = mapData
			}
		}

		if result := clickhouse.Clickhouse(config.New(), nil).Create(event); result != nil && result.Error == nil {
			log.Printf("Created successfully: %v", result.RowsAffected)
			return true
		} else if result == nil {
			log.Printf("Creation failed: %v", result)
		} else {
			log.Printf("Creation failed: %v", result.Error)
		}
	}

	return false
}

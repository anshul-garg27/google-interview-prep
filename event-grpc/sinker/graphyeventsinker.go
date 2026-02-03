package sinker

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
	"init.bulbul.tv/bulbul-backend/event-grpc/clickhouse"
	"init.bulbul.tv/bulbul-backend/event-grpc/client/identity"
	"init.bulbul.tv/bulbul-backend/event-grpc/client/input"
	"init.bulbul.tv/bulbul-backend/event-grpc/client/webengage"
	"init.bulbul.tv/bulbul-backend/event-grpc/config"
	"init.bulbul.tv/bulbul-backend/event-grpc/context"
	"init.bulbul.tv/bulbul-backend/event-grpc/model"
	"init.bulbul.tv/bulbul-backend/event-grpc/transformer"
	"log"
	"time"
)

func SinkGraphyEventToClickhouse(delivery amqp.Delivery) bool {

	rawGraphyEvent := map[string]interface{}{}
	if err := json.Unmarshal(delivery.Body, &rawGraphyEvent); err == nil {

		if rawGraphyEvent == nil {
			log.Printf("Incomplete event data: %v", rawGraphyEvent)
			return false
		}

		email := ""

		if emailInterface, ok := rawGraphyEvent["Email"]; ok {
			email = emailInterface.(string)
		} else if emailInterface, ok := rawGraphyEvent["Learner Email"]; ok {
			email = emailInterface.(string)
		}

		if email != "" {
			gatewayContext := context.New(&gin.Context{}, config.New())
			gatewayContext.GenerateClientAuthorizationWithClientId(gatewayContext.Config.GCC_HOST_APP_CLIENT_ID, "HOST")
			gatewayContext.DeviceAuthorization = "Basic YTpi"

			searchQuery := input.SearchQuery{
				Filters: []input.SearchFilter{
					{
						FilterType: "EQ",
						Field:      "clientId",
						Value:      gatewayContext.Client.Id,
					},
					{
						FilterType: "EQ",
						Field:      "email",
						Value:      email,
					},
				},
			}

			resp, _ := identity.New(gatewayContext).SearchUserClientAccount(&searchQuery)

			if resp != nil && len(resp.Accounts) > 0 && resp.Accounts[0].WebengageUserId != nil {

				ge := map[string]interface{}{}

				for k, v := range rawGraphyEvent {
					if vstr, ok := v.(string); ok {
						ge[k] = vstr
					}
				}

				webEngageEvent := input.WebengageEvent{
					UserId:    resp.Accounts[0].WebengageUserId,
					EventName: rawGraphyEvent["name"].(string),
					EventTime: time.Unix(transformer.InterfaceToInt64(rawGraphyEvent["receive-timestamp"].(float64)), 0).Format("2006-01-02T15:04:05Z0700"),
					EventData: ge,
				}

				resp, err := webengage.New(gatewayContext).DispatchWebengageEvent(&webEngageEvent)

				if err != nil {
					fmt.Errorf("Some error publishing to webengage: %v\n", err)
				} else {
					fmt.Printf("Published to webengage: %v\n", resp)
				}

			} else {
				fmt.Errorf("User not found for webengage: %v\n", rawGraphyEvent)
			}
		}

		event := model.GraphyEvent{
			Id:               rawGraphyEvent["x-eg-id"].(string),
			ReceiveTimestamp: time.Unix(transformer.InterfaceToInt64(rawGraphyEvent["receive-timestamp"].(float64)), 0),
			InsertTimestamp:  time.Now(),
			EventName:        rawGraphyEvent["name"].(string),
			EventData:        rawGraphyEvent,
			Headers:          rawGraphyEvent["x-eg-headers"].(map[string]interface{}),
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

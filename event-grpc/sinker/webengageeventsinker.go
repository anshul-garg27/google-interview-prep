package sinker

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
	"init.bulbul.tv/bulbul-backend/event-grpc/client/input"
	"init.bulbul.tv/bulbul-backend/event-grpc/client/webengage"
	"init.bulbul.tv/bulbul-backend/event-grpc/config"
	"init.bulbul.tv/bulbul-backend/event-grpc/context"
)

func SinkWebengageEventToWebengage(delivery amqp.Delivery) bool {

	webengageEvent := input.WebengageEvent{}
	if err := json.Unmarshal(delivery.Body, &webengageEvent); err == nil {

		gatewayContext := context.New(&gin.Context{}, config.New())
		gatewayContext.GenerateClientAuthorizationWithClientId(gatewayContext.Config.GCC_HOST_APP_CLIENT_ID, "HOST")
		gatewayContext.DeviceAuthorization = "Basic YTpi"

		// remove nil fields before event
		data := map[string]interface{}{}
		for k, v := range webengageEvent.EventData {
			if v != nil {
				data[k] = v
			}
		}
		webengageEvent.EventData = data

		resp, err := webengage.New(gatewayContext).DispatchWebengageEvent(&webengageEvent)

		if err != nil || resp == nil || resp.Response.Status == "error" {
			fmt.Errorf("Some error publishing to webengage: %v %v\n", err, resp)
		} else {
			fmt.Printf("Published to webengage: %v\n", resp)
			return true
		}
	}

	return false
}

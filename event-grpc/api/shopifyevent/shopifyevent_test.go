package shopifyevent

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"init.bulbul.tv/bulbul-backend/event-grpc/client/input"
	"init.bulbul.tv/bulbul-backend/event-grpc/config"
	"init.bulbul.tv/bulbul-backend/event-grpc/context"
	"log"
	"testing"
)

func TestDispatchEvent(t *testing.T) {
	godotenv.Load("../../.env")

	sinkContext := context.New(&gin.Context{}, config.New())
	sinkContext.GenerateClientAuthorizationWithClientId(sinkContext.Config.CUSTOMER_APP_CLIENT_ID, "CUSTOMER")
	sinkContext.DeviceAuthorization = "Basic YTpi"

	err := DispatchEvent(sinkContext, input.ShopifyEvent{
		UserId:      nil,
		AnonymousId: nil,
		EventName:   "",
		EventTime:   "",
		EventData:   nil,
	})

	if err == nil {
		log.Println("Published")
	} else {
		log.Fatalf("Error publishing: %v", err)
	}

}

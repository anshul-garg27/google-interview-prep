package branchevent

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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

	err := DispatchEvent(sinkContext, map[string]interface{}{
		"hello": "world",
	})

	if err == nil {
		log.Println("Published")
	} else {
		log.Fatalf("Error publishing: %v", err)
	}

}

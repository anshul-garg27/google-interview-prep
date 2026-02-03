package identity

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"init.bulbul.tv/bulbul-backend/event-grpc/client/input"
	"init.bulbul.tv/bulbul-backend/event-grpc/config"
	"init.bulbul.tv/bulbul-backend/event-grpc/context"
	"testing"
)

func TestClient_SearchUserClientAccount(t *testing.T) {
	err := godotenv.Load("../../.env")

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
				Value:      "ekansh.gupta@myglamm.com",
			},
		},
	}

	userClientAccountResponse, err := New(gatewayContext).SearchUserClientAccount(&searchQuery)

	if err != nil {
		gatewayContext.Logger.Fatal().Msg(fmt.Sprintf("failed: %+v", userClientAccountResponse))
	} else {
		gatewayContext.Logger.Log().Msg(fmt.Sprintf("response: %+v", userClientAccountResponse))
	}
}

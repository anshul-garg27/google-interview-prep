package socialstream

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"init.bulbul.tv/bulbul-backend/event-grpc/client/input"
	"init.bulbul.tv/bulbul-backend/event-grpc/config"
	"init.bulbul.tv/bulbul-backend/event-grpc/context"
	"testing"
)

func TestClient_SearchInfluencerCampaigns(t *testing.T) {
	gatewayContext := context.New(&gin.Context{}, config.New())
	gatewayContext.GenerateClientAuthorizationWithClientId(gatewayContext.Config.CUSTOMER_APP_CLIENT_ID, "CUSTOMER")

	gatewayContext.DeviceAuthorization = "Basic YTpi"

	query := input.SearchQuery{
		Filters: []input.SearchFilter{
			{
				FilterType: "EQ",
				Field:      "campaignCode",
				Value:      "anjana-v1-jan-youtube",
			},
			{
				FilterType: "EQ",
				Field:      "status",
				Value:      "LIVE",
			},
			{
				FilterType: "EQ",
				Field:      "isExternal",
				Value:      "false",
			},
		},
	}

	influencerCampaigns, err := New(gatewayContext).SearchInfluencerCampaigns(query, 1, 10, true)

	if err != nil {
		gatewayContext.Logger.Fatal().Msg(fmt.Sprintf("failed: %+v", influencerCampaigns))
	} else {
		gatewayContext.Logger.Log().Msg(fmt.Sprintf("response: %+v", influencerCampaigns))
	}
}
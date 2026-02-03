package partner

import (
	"fmt"
	"testing"

	"github.com/gin-gonic/gin"
	"init.bulbul.tv/bulbul-backend/saas-gateway/client/entry"
	"init.bulbul.tv/bulbul-backend/saas-gateway/config"
	"init.bulbul.tv/bulbul-backend/saas-gateway/context"
)

func TestClient_FindPartnerById(t *testing.T) {
	gatewayContext := context.New(&gin.Context{}, config.New())
	gatewayContext.GenerateClientAuthorizationWithClientId(gatewayContext.Config.ERP_CLIENT_ID, "ERP")
	gatewayContext.SetUserClientAccount(&entry.UserClientAccount{UserId: 1014})
	partnerDetailResponse, err := New(gatewayContext).FindPartnerById("21783")

	if err != nil {
		gatewayContext.Logger.Fatal().Msgf("Authentication failed: %+v", partnerDetailResponse)
	} else {
		gatewayContext.Logger.Log().Msg(fmt.Sprintf("success: %+v", partnerDetailResponse))
	}
}

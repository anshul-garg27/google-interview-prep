package identity

import (
	"fmt"
	"testing"

	"github.com/gin-gonic/gin"
	"init.bulbul.tv/bulbul-backend/saas-gateway/client/input"
	"init.bulbul.tv/bulbul-backend/saas-gateway/client/response"
	"init.bulbul.tv/bulbul-backend/saas-gateway/config"
	"init.bulbul.tv/bulbul-backend/saas-gateway/context"
)

func userAuthForTest() (*context.Context, *response.UserAuthResponse, error) {
	gatewayContext := context.New(&gin.Context{}, config.New())
	gatewayContext.GenerateClientAuthorizationWithClientId(gatewayContext.Config.ERP_CLIENT_ID, "ERP")

	userAuthResponse, err := New(gatewayContext).UserAuth(&input.UserAuth{
		Id:     "gateway@bulbul.tv",
		Secret: "bulbul@123",
		Mode:   "EMAIL",
	})
	return gatewayContext, userAuthResponse, err
}

func TestClient_VerifyToken(t *testing.T) {
	gatewayContext, userAuthResponse, _ := userAuthForTest()

	response, err := New(gatewayContext).VerifyToken(&input.Token{
		Token: userAuthResponse.Token,
	}, false)

	if err != nil {
		gatewayContext.Logger.Fatal().Msg(fmt.Sprintf("Verify token failed: %+v", err))
	} else {
		gatewayContext.Logger.Log().Msg(fmt.Sprintf("Verify token: %+v", response))
	}
}

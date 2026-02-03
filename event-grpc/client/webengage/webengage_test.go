package webengage

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"init.bulbul.tv/bulbul-backend/event-grpc/client/input"
	"init.bulbul.tv/bulbul-backend/event-grpc/config"
	"init.bulbul.tv/bulbul-backend/event-grpc/context"
	"testing"
	"time"
)

func TestClient_DispatchWebengageEvent(t *testing.T) {
	err := godotenv.Load("../../.env")

	gatewayContext := context.New(&gin.Context{}, config.New())
	gatewayContext.GenerateClientAuthorizationWithClientId(gatewayContext.Config.GCC_HOST_APP_CLIENT_ID, "HOST")

	userId := "12345"
	appType := input.WEBENGAGE_APP_TYPE_USER
	webengageEvent := input.WebengageEvent{
		AppType:   &appType,
		UserId:    &userId,
		EventName: "Test Event 1234",
		EventTime: time.Now().Format("2006-01-02T15:04:05Z0700"),
		EventData: map[string]interface{}{
			"Test":         "Me",
			"Also me":      123,
			"try arr":      []int{12, 34},
			"fdwsstry arr": []string{"12", "34"},
			"Nice":         true,
		},
	}

	webengageResp, err := New(gatewayContext).DispatchWebengageEvent(&webengageEvent)

	if err != nil {
		gatewayContext.Logger.Fatal().Msg(fmt.Sprintf("failed: %+v %+v", webengageResp, err))
	} else {
		gatewayContext.Logger.Log().Msg(fmt.Sprintf("response: %+v", webengageResp))
	}
}

func TestClient_DispatchWebengageUserEvent(t *testing.T) {
	err := godotenv.Load("../../.env")

	gatewayContext := context.New(&gin.Context{}, config.New())
	gatewayContext.GenerateClientAuthorizationWithClientId(gatewayContext.Config.GCC_HOST_APP_CLIENT_ID, "HOST")

	userId := "12345"
	appType := input.WEBENGAGE_APP_TYPE_USER
	phone := "9852147890"
	email := "jhbvcjhs@gmail.com"
	fname := "bjhbj"
	webengageEvent := input.WebengageUserEvent{
		AppType:   &appType,
		UserId:    &userId,
		Phone:     &phone,
		Email:     &email,
		FirstName: &fname,
		Attributes: map[string]interface{}{
			"Test":         "Me",
			"Also me":      123,
			"try arr":      []int{12, 34},
			"fdwsstry arr": []string{"12", "34"},
			"Nice":         true,
		},
	}

	webengageResp, err := New(gatewayContext).DispatchWebengageUserEvent(&webengageEvent)

	if err != nil {
		gatewayContext.Logger.Fatal().Msg(fmt.Sprintf("failed: %+v %+v", webengageResp, err))
	} else {
		gatewayContext.Logger.Log().Msg(fmt.Sprintf("response: %+v", webengageResp))
	}
}

func TestClient_DispatchWebengageEvent2(t *testing.T) {
	fmt.Println(time.Now().Format("2006-01-02T15:04:05Z0700"))
}

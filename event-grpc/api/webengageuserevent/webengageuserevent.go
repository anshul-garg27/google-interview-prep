package webengageuserevent

import (
	"errors"
	"github.com/gin-gonic/gin"
	"init.bulbul.tv/bulbul-backend/event-grpc/client/input"
	"init.bulbul.tv/bulbul-backend/event-grpc/config"
	"init.bulbul.tv/bulbul-backend/event-grpc/context"
	"init.bulbul.tv/bulbul-backend/event-grpc/util"
	"init.bulbul.tv/bulbul-backend/event-grpc/webengageusereventworker"
	"net/http"
	"time"
)

func HandleWebengageUserEvent(config config.Config) func(c *gin.Context) {
	return func(c *gin.Context) {
		gCtx := util.GatewayContextFromGinContext(c, config)

		webengageEvent := input.WebengageUserEvent{}
		if err := c.ShouldBindJSON(&webengageEvent); err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		if webengageEvent.UserId == nil {
			c.Status(http.StatusBadRequest)
			return
		}

		err := DispatchEvent(gCtx, webengageEvent)
		if err == nil {
			gCtx.Status(http.StatusOK)
		} else {
			gCtx.Status(http.StatusInternalServerError)
		}
	}
}

func DispatchEvent(gCtx *context.Context, webengageEvent input.WebengageUserEvent) error {
	if eventChannel := webengageusereventworker.GetChannel(gCtx.Config); eventChannel == nil {
		gCtx.Logger.Error().Msgf("Event channel nil: %v", webengageEvent)
	} else {
		select {
		case <-time.After(100 * time.Second):
			gCtx.Logger.Log().Msgf("Error publishing event: %v", webengageEvent)
		case eventChannel <- webengageEvent:
			return nil
		}
	}
	return errors.New("Error publishing event")
}

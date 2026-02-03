package shopifyevent

import (
	"errors"
	"github.com/gin-gonic/gin"
	"init.bulbul.tv/bulbul-backend/event-grpc/client/input"
	"init.bulbul.tv/bulbul-backend/event-grpc/config"
	"init.bulbul.tv/bulbul-backend/event-grpc/context"
	"init.bulbul.tv/bulbul-backend/event-grpc/shopifyeventworker"
	"init.bulbul.tv/bulbul-backend/event-grpc/util"
	"net/http"
	"time"
)

func HandleShopifyEvent(config config.Config) func(c *gin.Context) {
	return func(c *gin.Context) {
		gCtx := util.GatewayContextFromGinContext(c, config)

		event := input.ShopifyEvent{}
		if err := c.ShouldBindJSON(&event); err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		err := DispatchEvent(gCtx, event)
		if err == nil {
			gCtx.Status(http.StatusOK)
		} else {
			gCtx.Status(http.StatusInternalServerError)
		}
	}
}

func DispatchEvent(gCtx *context.Context, event input.ShopifyEvent) error {
	if eventChannel := shopifyeventworker.GetChannel(gCtx.Config); eventChannel == nil {
		gCtx.Logger.Error().Msgf("Event channel nil: %v", event)
	} else {
		select {
		case <-time.After(1 * time.Second):
			gCtx.Logger.Log().Msgf("Error publishing event: %v", event)
		case eventChannel <- event:
			return nil
		}
	}
	return errors.New("Error publishing event")
}

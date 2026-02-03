package vidoolyevent

import (
	"errors"
	"github.com/gin-gonic/gin"
	"init.bulbul.tv/bulbul-backend/event-grpc/config"
	"init.bulbul.tv/bulbul-backend/event-grpc/context"
	"init.bulbul.tv/bulbul-backend/event-grpc/util"
	"init.bulbul.tv/bulbul-backend/event-grpc/vidoolyeventworker"
	"net/http"
	"time"
)

func HandleVidoolyEvent(config config.Config) func(c *gin.Context) {
	return func(c *gin.Context) {
		gCtx := util.GatewayContextFromGinContext(c, config)

		vidoolyEvent := map[string]interface{}{}
		if err := c.ShouldBindJSON(&vidoolyEvent); err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		err := DispatchEvent(gCtx, vidoolyEvent)
		if err == nil {
			gCtx.Status(http.StatusOK)
		} else {
			gCtx.Status(http.StatusInternalServerError)
		}
	}
}

func DispatchEvent(gCtx *context.Context, vidoolyEvent map[string]interface{}) error {
	if eventChannel := vidoolyeventworker.GetChannel(gCtx.Config); eventChannel == nil {
		gCtx.Logger.Error().Msgf("Event channel nil: %v", vidoolyEvent)
	} else {
		select {
		case <-time.After(100 * time.Second):
			gCtx.Logger.Log().Msgf("Error publishing event: %v", vidoolyEvent)
		case eventChannel <- vidoolyEvent:
			return nil
		}
	}
	return errors.New("Error publishing event")
}

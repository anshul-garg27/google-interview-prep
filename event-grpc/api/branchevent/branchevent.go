package branchevent

import (
	"errors"
	"github.com/gin-gonic/gin"
	"init.bulbul.tv/bulbul-backend/event-grpc/brancheventworker"
	"init.bulbul.tv/bulbul-backend/event-grpc/config"
	"init.bulbul.tv/bulbul-backend/event-grpc/context"
	"init.bulbul.tv/bulbul-backend/event-grpc/util"
	"net/http"
	"time"
)

func HandleBranchEvent(config config.Config) func(c *gin.Context) {
	return func(c *gin.Context) {
		gCtx := util.GatewayContextFromGinContext(c, config)

		branchEvent := map[string]interface{}{}
		if err := c.ShouldBindJSON(&branchEvent); err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		err := DispatchEvent(gCtx, branchEvent)
		if err == nil {
			gCtx.Status(http.StatusOK)
		} else {
			gCtx.Status(http.StatusInternalServerError)
		}
	}
}

func DispatchEvent(gCtx *context.Context, branchEvent map[string]interface{}) error {
	if eventChannel := brancheventworker.GetChannel(gCtx.Config); eventChannel == nil {
		gCtx.Logger.Error().Msgf("Event channel nil: %v", branchEvent)
	} else {
		select {
		case <-time.After(100 * time.Second):
			gCtx.Logger.Log().Msgf("Error publishing event: %v", branchEvent)
		case eventChannel <- branchEvent:
			return nil
		}
	}
	return errors.New("Error publishing event")
}

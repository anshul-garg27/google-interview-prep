package graphyevent

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"init.bulbul.tv/bulbul-backend/event-grpc/config"
	"init.bulbul.tv/bulbul-backend/event-grpc/context"
	"init.bulbul.tv/bulbul-backend/event-grpc/graphyeventworker"
	"init.bulbul.tv/bulbul-backend/event-grpc/util"
	"net/http"
	"time"
)

func HandleGraphyEvent(config config.Config) func(c *gin.Context) {
	return func(c *gin.Context) {
		gCtx := util.GatewayContextFromGinContext(c, config)

		graphyEvent := map[string]interface{}{}
		if err := c.ShouldBindJSON(&graphyEvent); err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		headerMap := gCtx.Request.Header

		eventName := "NO_EVENT_NAME"
		if val, ok := headerMap["Spayee-Event"]; ok && len(val) > 0 {
			eventName = val[0]
		}

		graphyEvent["x-eg-headers"] = headerMap
		graphyEvent["name"] = eventName
		graphyEvent["receive-timestamp"] = time.Now().Unix()

		fmt.Println(graphyEvent)
		graphyEvent["x-eg-id"] = uuid.New().String()

		err := DispatchEvent(gCtx, graphyEvent)
		if err == nil {
			gCtx.Status(http.StatusOK)
		} else {
			gCtx.Status(http.StatusInternalServerError)
		}
	}
}

func DispatchEvent(gCtx *context.Context, graphyEvent map[string]interface{}) error {
	if eventChannel := graphyeventworker.GetChannel(gCtx.Config); eventChannel == nil {
		gCtx.Logger.Error().Msgf("Event channel nil: %v", graphyEvent)
	} else {
		select {
		case <-time.After(100 * time.Second):
			gCtx.Logger.Log().Msgf("Error publishing event: %v", graphyEvent)
		case eventChannel <- graphyEvent:
			return nil
		}
	}
	return errors.New("Error publishing event")
}

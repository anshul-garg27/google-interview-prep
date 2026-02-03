package heartbeat

import (
	"github.com/gin-gonic/gin"
	"init.bulbul.tv/bulbul-backend/event-grpc/config"
	"net/http"
	"strconv"
)

var BEAT = false

func Beat(config config.Config) func(c *gin.Context) {
	return func(c *gin.Context) {
		if BEAT {
			c.Status(http.StatusOK)
		} else {
			c.Status(http.StatusGone)
		}
	}
}

func ModifyBeat(config config.Config) func(c *gin.Context) {
	return func(c *gin.Context) {
		modifyBeat, err := strconv.ParseBool(c.Query("beat"))

		if err != nil {
			c.Status(http.StatusBadRequest)
		} else {
			BEAT = modifyBeat
			c.Status(http.StatusOK)
		}
	}
}

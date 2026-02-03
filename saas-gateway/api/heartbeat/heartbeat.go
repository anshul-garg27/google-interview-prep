package heartbeat

import (
	"github.com/gin-gonic/gin"
	"init.bulbul.tv/bulbul-backend/saas-gateway/config"
	"net/http"
	"strconv"
)

var beat = false

func Beat(config config.Config) func(c *gin.Context) {
	return func(c *gin.Context) {
		if beat {
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
			beat = modifyBeat
			c.Status(http.StatusOK)
		}
	}
}

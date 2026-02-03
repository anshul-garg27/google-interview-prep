package route

import (
	"github.com/gin-gonic/gin"
	"init.bulbul.tv/bulbul-backend/saas-gateway/config"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc func(config config.Config) func(c *gin.Context)
}

type Routes []Route

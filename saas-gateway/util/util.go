package util

import (
	"github.com/gin-gonic/gin"
	"init.bulbul.tv/bulbul-backend/saas-gateway/config"
	gatewayContext "init.bulbul.tv/bulbul-backend/saas-gateway/context"
)

func GatewayContextFromGinContext(ctx *gin.Context, config config.Config) *gatewayContext.Context {
	if ctx == nil || ctx.Request == nil || ctx.Request.Context() == nil || ctx.Request.Context().Value("GatewayContext") == nil {
		return gatewayContext.New(ctx, config)
	}
	return ctx.Request.Context().Value("GatewayContext").(*gatewayContext.Context)
}

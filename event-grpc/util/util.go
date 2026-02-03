package util

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"init.bulbul.tv/bulbul-backend/event-grpc/config"
	gatewayContext "init.bulbul.tv/bulbul-backend/event-grpc/context"
	"strconv"
	"time"
)

func WrapErrors(errs []error) error {
	var err error
	for _, error := range errs {
		err = errors.WithStack(error)
	}
	return err
}

func DayStartTimestamp() string {
	year, month, day := time.Now().Date()
	return strconv.FormatInt(time.Date(year, month, day, 0, 0, 0, 0, time.Now().Location()).Unix()*1000, 10)
}

func GatewayContextFromContext(ctx context.Context) (*gatewayContext.Context, error) {

	gatewayCtx := ctx.Value("GatewayContext")
	if gatewayCtx == nil {
		err := fmt.Errorf("could not retrieve gin.Context")
		return nil, err
	}

	gc, ok := gatewayCtx.(*gatewayContext.Context)
	if !ok {
		err := fmt.Errorf("gin.Context has wrong type")
		return nil, err
	}
	return gc, nil
}

func GatewayContextFromGinContext(ctx *gin.Context, config config.Config) *gatewayContext.Context {

	if ctx == nil || ctx.Request == nil || ctx.Request.Context() == nil || ctx.Request.Context().Value("GatewayContext") == nil {
		return gatewayContext.New(ctx, config)
	}

	return ctx.Request.Context().Value("GatewayContext").(*gatewayContext.Context)
}
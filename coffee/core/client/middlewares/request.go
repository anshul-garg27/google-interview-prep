package middlewares

import (
	"coffee/constants"
	"coffee/core/appcontext"
	"github.com/go-resty/resty/v2"
)

func PopulateHeadersRequestMiddleware(c *resty.Client, req *resty.Request) error {
	ctx := req.Context()
	val := ctx.Value(constants.AppContextKey)
	if val != nil {
		appCtx := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)
		appCtx.PopulateHeadersForHttpReq(req)
	}
	return nil
}

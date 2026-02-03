package client

import (
	"coffee/core/client/middlewares"
	"context"
	"time"

	"github.com/spf13/viper"

	"github.com/go-resty/resty/v2"
)

type BaseClient struct {
	*resty.Client
	Context context.Context
	BaseUrl string
}

func New(ctx context.Context, baseUrlKey string) *BaseClient {
	client := &BaseClient{
		resty.New().
			SetDebug(true).
			SetTimeout(5 * time.Second),
		ctx,
		viper.GetString(baseUrlKey),
	}

	client.OnBeforeRequest(middlewares.PopulateHeadersRequestMiddleware)
	return client
}

func (c *BaseClient) GenerateRequest() *resty.Request {
	c.SetAllowGetMethodPayload(true)
	req := c.NewRequest()
	req.SetContext(c.Context)
	return req
}

func (c *BaseClient) GenerateAuthorizedRequest(authorization string) *resty.Request {
	c.SetAllowGetMethodPayload(true)
	req := c.NewRequest()
	req.SetHeader("Authorization", authorization)

	// patch X-Headers from Context to request for flow to upstream
	// if c.Context.XHeader != nil {
	// 	for key, vals := range c.Context.XHeader {
	// 		for _, val := range vals {
	// 			req.Header.Add(key, val)
	// 		}
	// 	}
	// }

	return req
}

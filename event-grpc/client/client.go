package client

import (
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog"
	gatewayContext "init.bulbul.tv/bulbul-backend/event-grpc/context"
	"time"
)

type BaseClient struct {
	*resty.Client
	Context *gatewayContext.Context
}

func New(ctx *gatewayContext.Context) *BaseClient {
	return &BaseClient{
		resty.New().
			SetDebug(true).
			SetTimeout(5 * time.Second).
			SetLogger(NewLogger(ctx.Logger)).
			// disable response body logging if log level >= 3
			OnRequestLog(func(r *resty.RequestLog) error {
				if ctx.Config.LOG_LEVEL >= 3 {
					r.Body = ""
				}
				return nil
			}).
			OnResponseLog(func(r *resty.ResponseLog) error {
				if ctx.Config.LOG_LEVEL >= 3 {
					r.Body = ""
				}
				return nil
			}),
		ctx,
	}
}

func (c *BaseClient) GenerateRequest(authorization string) *resty.Request {
	c.SetAllowGetMethodPayload(true)
	req := c.NewRequest()
	req.SetHeader("Authorization", authorization)

	// patch X-Headers from Context to request for flow to upstream
	if c.Context.XHeader != nil {
		for key, vals := range c.Context.XHeader {
			for _, val := range vals {
				req.Header.Add(key, val)
			}
		}
	}

	return req
}

// custom resty logger interface impl
type logger struct {
	zerolog.Logger
}

func NewLogger(zlog zerolog.Logger) logger {
	return logger{zlog}
}

func (l logger) Errorf(format string, v ...interface{}) {
	l.Error().Msgf(format, v...)
}

func (l logger) Warnf(format string, v ...interface{}) {
	l.Warn().Msgf(format, v...)
}

// debug level for resty would be info level for us
func (l logger) Debugf(format string, v ...interface{}) {
	l.Error().Msgf(format, v...)
}

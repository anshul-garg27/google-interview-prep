package middlewares

import (
	"coffee/constants"
	"coffee/core/appcontext"
	"coffee/core/rest"
	"coffee/helpers"
	"context"
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

/**
TODO: Refactor into 2 separate middlewares
Had to do it in a single one because chis default timeout middleware was not playing nicely with our transactional semantics
*/

func ServiceSessionMiddlewares(container rest.ApplicationContainer) func(handler http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), time.Duration(viper.GetInt("SERVER_TIMEOUT"))*time.Second)
			defer cancel()
			r = r.WithContext(ctx)
			if r.URL.Path == "/metrics" {
				handler.ServeHTTP(w, r)
			}
			for service, api := range container.Services {
				if strings.HasPrefix(r.URL.String(), api.GetPrefix()) {
					log.Debug("Attaching ServiceSessionMiddleware : ", r.URL, "  ", api.GetPrefix())
					sessionCtx := r.Context()
					defer func(service rest.ServiceWrapper, sessionCtx context.Context) {
						cancel()
						if ctx.Err() == context.DeadlineExceeded {
							customWriter := w.(*helpers.ResponseWriterWithContent)
							helpers.GetSentryEvent(customWriter, r, customWriter.InputBodyCopy)
							rollbackSessions(sessionCtx)
							w.WriteHeader(http.StatusGatewayTimeout)
						} else {
							closeSessions(sessionCtx)
						}
					}(service, sessionCtx)
					handler.ServeHTTP(w, r.WithContext(sessionCtx))
				}
			}

		}
		return http.HandlerFunc(fn)
	}
}

func closeSessions(ctx context.Context) {
	session := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext).GetSession()
	if session != nil {
		session.Close(ctx)
	}
	chSession := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext).GetCHSession()
	if chSession != nil {
		chSession.Close(ctx)
	}
}

func rollbackSessions(ctx context.Context) {
	session := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext).GetSession()
	if session != nil {
		session.Rollback(ctx)
	}
	chSession := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext).GetCHSession()
	if chSession != nil {
		chSession.Rollback(ctx)
	}
}

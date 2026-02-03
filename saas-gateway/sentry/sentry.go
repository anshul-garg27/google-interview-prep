package sentry

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/99designs/gqlgen/graphql"
	"github.com/getsentry/sentry-go"
	"init.bulbul.tv/bulbul-backend/saas-gateway/config"
	"init.bulbul.tv/bulbul-backend/saas-gateway/context"
	"init.bulbul.tv/bulbul-backend/saas-gateway/header"
)

func Setup(config config.Config) {
	// sentry init
	if config.Env != "local" {
		if err := sentry.Init(sentry.ClientOptions{
			Dsn:              config.SENTRY_DSN,
			Environment:      config.Env,
			AttachStacktrace: true,
		}); err != nil {
			log.Println("Sentry init failed")
		}
	}
}

func PatchErrorToSentry(gc *context.Context, e error, err interface{}) {

	gcInterface := make(map[string]interface{})

	if gc != nil {
		gcInterface["headers"] = gc.XHeader

		if gc.UserClientAccount != nil {
			gcInterface["user"] = gc.UserClientAccount
		}

		if gc.Client != nil {
			gcInterface["client"] = gc.Client
		}

		gcInterface["path"] = graphql.GetFieldContext(gc.Context).Path()

		sentry.ConfigureScope(func(scope *sentry.Scope) {
			tags := map[string]string{
				header.RequestID:    gc.XHeader.Get(header.RequestID),
				header.ApolloOpName: gc.XHeader.Get(header.ApolloOpName),
				header.Version:      gc.XHeader.Get(header.Version),
				header.ClientID:     gc.XHeader.Get(header.ClientID),
				header.Channel:      gc.XHeader.Get(header.Channel),
			}

			if gc.UserClientAccount != nil {
				tags["userId"] = strconv.Itoa(gc.UserClientAccount.UserId)
			}

			scope.SetTags(tags)
		})

		sentry.AddBreadcrumb(&sentry.Breadcrumb{
			Category:  "context",
			Data:      gcInterface,
			Message:   "Context",
			Timestamp: time.Time{},
		})
	}

	if e != nil {
		sentry.CaptureException(e)
	}

	if err != nil {
		sentry.CaptureException(errors.New(fmt.Sprintf("%+v", err)))
	}
}

package router

import (
	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	cors "github.com/rs/cors/wrapper/gin"
	"init.bulbul.tv/bulbul-backend/event-grpc/config"
	"init.bulbul.tv/bulbul-backend/event-grpc/middleware"
	"init.bulbul.tv/bulbul-backend/event-grpc/route"
	"log"
	"net/http"
)

func SetupRouter(config config.Config) *gin.Engine {
	// sentry init
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              "http://3162a2a0ef954d3cb3b61dc9a9db26f0@172.31.14.149:9000/26",
		Environment:      config.Env,
		AttachStacktrace: true,
	}); err != nil {
		log.Println("Sentry init failed")
	}

	router := gin.New()
	router.Use(sentrygin.New(sentrygin.Options{Repanic: true}))
	router.Use(gin.Recovery())
	router.Use(cors.New(cors.Options{
		AllowedOrigins: []string{
			"*",
		},
		AllowedMethods: []string{
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}))
	router.Use(middleware.EventGrpcContextToContextMiddleware(config), middleware.RequestIdMiddleware(), middleware.RequestLogger(config))

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	heartbeatRouter := router.Group("/heartbeat")
	groupRoutes(config, heartbeatRouter, route.HeartbeatRoutes)

	vidoolyEventRouter := router.Group("/vidooly/event")
	groupRoutes(config, vidoolyEventRouter, route.VidoolyEventRoutes)

	branchEventRouter := router.Group("/branch/event")
	groupRoutes(config, branchEventRouter, route.BranchEventRoutes)

	webengageRouter := router.Group("/webengage/event")
	groupRoutes(config, webengageRouter, route.WebengageEventRoutes)

	shopifyRouter := router.Group("/shopify/event")
	groupRoutes(config, shopifyRouter, route.ShopifyEventRoutes)

	graphyEventRouter := router.Group("/graphy/event")
	groupRoutes(config, graphyEventRouter, route.GraphyEventRoutes)

	return router
}

func groupRoutes(config config.Config, router *gin.RouterGroup, routes route.Routes) {
	{
		for _, route := range routes {
			router.Handle(route.Method, route.Pattern, route.HandlerFunc(config))
		}
	}
}

package router

import (
	"context"
	"net/http"

	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	cors "github.com/rs/cors/wrapper/gin"
	"init.bulbul.tv/bulbul-backend/saas-gateway/config"
	gatewayContext "init.bulbul.tv/bulbul-backend/saas-gateway/context"
	"init.bulbul.tv/bulbul-backend/saas-gateway/handler/saas"
	"init.bulbul.tv/bulbul-backend/saas-gateway/middleware"
	"init.bulbul.tv/bulbul-backend/saas-gateway/route"
)

func SetupRouter(config config.Config) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(sentrygin.New(sentrygin.Options{Repanic: true}))
	router.Use(cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://staging.app.vidooly.com",
			"https://stage.cf-provider.goodcreator.co",
			"https://cf-provider.goodcreator.co",
			"https://stage.cf-services.goodcreator.co",
			"https://cf-services.goodcreator.co",
			"https://www.instagram.com",
			"https://www.youtube.com",
			"https://stage.services.goodcreator.co",
			"http://stage.services.goodcreator.co",
			"https://stage.provider.goodcreator.co",
			"http://stage.provider.goodcreator.co",
			"https://provider.goodcreator.co",
			"http://provider.goodcreator.co",
			"http://app.vidooly.com",
			"https://staging.app.vidooly.com",
			"https://suite.goodcreator.co",
			"https://beta.suite.goodcreator.co",
			"http://beta.suite.goodcreator.co",
			"https://momentum.goodcreator.co",
			"http://suite.goodcreator.co",
			"https://stage.suite.goodcreator.co",
			"https://stage.momentum.goodcreator.co",
			"http://stage.suite.goodcreator.co",
			"https://goodcreator.co",
			"https://staging.goodcreator.co",
			"https://app.vidooly.com",
			"http://localhost:3000",
			"http://localhost:3001",
			"http://localhost:8298",
			"http://staging.campaignmanager.goodcreator.co",
			"http://campaignmanager.goodcreator.co",
			"https://staging.campaignmanager.goodcreator.co",
			"https://campaignmanager.goodcreator.co",
			"https://beta.provider.goodcreator.co",
			"http://beta.provider.goodcreator.co",
			"https://beta.services.goodcreator.co",
			"http://beta.services.goodcreator.co",
		},
		AllowedMethods: []string{
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodOptions,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}))
	router.Use(GatewayContextToContextMiddleware(config), middleware.RequestIdMiddleware(), middleware.RequestLogger(config))
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	heartbeatRouter := router.Group("/heartbeat")
	groupRoutes(config, heartbeatRouter, route.HeartbeatRoutes)
	router.Any("/discovery-service/*any", middleware.AppAuth(config), saas.ReverseProxy("discovery-service", config))
	router.Any("/leaderboard-service/*any", middleware.AppAuth(config), saas.ReverseProxy("leaderboard-service", config))
	router.Any("/profile-collection-service/*any", middleware.AppAuth(config), saas.ReverseProxy("profile-collection-service", config))
	router.Any("/activity-service/*any", middleware.AppAuth(config), saas.ReverseProxy("activity-service", config))
	router.Any("/collection-analytics-service/*any", middleware.AppAuth(config), saas.ReverseProxy("collection-analytics-service", config))
	router.Any("/post-collection-service/*any", middleware.AppAuth(config), saas.ReverseProxy("post-collection-service", config))
	router.Any("/genre-insights-service/*any", middleware.AppAuth(config), saas.ReverseProxy("genre-insights-service", config))
	router.Any("/content-service/*any", middleware.AppAuth(config), saas.ReverseProxy("content-service", config))
	router.Any("/social-profile-service/*any", middleware.AppAuth(config), saas.ReverseProxy("social-profile-service", config))
	router.Any("/collection-group-service/*any", middleware.AppAuth(config), saas.ReverseProxy("collection-group-service", config))
	router.Any("/keyword-collection-service/*any", middleware.AppAuth(config), saas.ReverseProxy("keyword-collection-service", config))
	router.Any("/partner-usage-service/*any", middleware.AppAuth(config), saas.ReverseProxy("partner-usage-service", config))
	router.Any("/campaign-profile-service/*any", middleware.AppAuth(config), saas.ReverseProxy("campaign-profile-service", config))

	return router
}

func groupRoutes(config config.Config, router *gin.RouterGroup, routes route.Routes) {
	{
		for _, route := range routes {
			router.Handle(route.Method, route.Pattern, route.HandlerFunc(config))
		}
	}
}

func GatewayContextToContextMiddleware(config config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), "GatewayContext", gatewayContext.New(c, config)))
		c.Next()
	}
}

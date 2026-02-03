package server

import (
	"coffee/app"
	"coffee/constants"
	"coffee/core/persistence/clickhouse"
	"coffee/core/persistence/postgres"
	"coffee/core/rest"
	"coffee/listeners"
	"coffee/routes"
	"coffee/server/middlewares"
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	chiprometheus "github.com/apblbl/chi-prometheus"
	"github.com/getsentry/sentry-go"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Setup() {
	//log.Error("sample error in server.go")
	setupDB()
	setupClickHouseDB()
	listeners.SetupListeners()
	r := setupRoutes()
	server := createServer(":"+viper.GetString(constants.ServerPortConfig), r)
	if server != nil {
		setupSignals(server)
		err := startListening(server)
		if err != nil {
			sentry.CaptureException(err)
			log.Fatal(err)
		}
	} else {
		sentry.CaptureMessage("Failed to initialize server instance")
		log.Fatal("Failed to initialize server instance")
	}
}

func setupDB() {
	postgres.DB() // Setup DB Connection Pool
}

func setupClickHouseDB() {
	clickhouse.ClickhouseDB()
}

func createServer(address string, router http.Handler) *http.Server {
	/**
	Setup base server with sane defaults
	See also:
		https://adam-p.ca/blog/2022/01/golang-http-server-timeouts/
		https://blog.cloudflare.com/exposing-go-on-the-internet/#timeouts
	*/
	return &http.Server{
		Addr:              address,
		Handler:           router,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      100 * time.Second,
		IdleTimeout:       15 * time.Second,
	}
}

func setupSignals(server *http.Server) {
	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	go func() {
		<-quit
		shutdownServer(server)
		close(done)
	}()
}

func shutdownServer(server *http.Server) {
	log.Debug("Server is shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	server.SetKeepAlivesEnabled(false)
	defer sentry.Flush(time.Second * 5)
	if err := server.Shutdown(ctx); err != nil {
		sentry.CaptureException(err)
		log.Fatal("Something went wrong while shutting down server", err)
	}
}

func setupRoutes() *chi.Mux {
	r := chi.NewRouter()
	container := app.SetupContainer()
	attachMiddlewares(r, container)
	routes.AttachMetricsRoutes(r)
	routes.AttachHeartbeatRoutes(r)
	routes.AttachServiceRoutes(r, container)
	log.Debug("Completed initialization of routes")
	return r
}

func startListening(server *http.Server) error {
	return server.ListenAndServe()
}

func attachMiddlewares(r *chi.Mux, container rest.ApplicationContainer) {

	//This Middleware Needs to be Call first as it makes Custom Response Writer
	r.Use(middlewares.RequestInterceptor)
	r.Use(middlewares.SlowQueryLogger)
	r.Use(middlewares.SentryErrorLoggingMiddleware)
	r.Use(middlewares.ApplicationContext)
	r.Use(middlewares.ServiceSessionMiddlewares(container))
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(chiprometheus.NewPatternMiddleware(constants.ApplicationName, 300, 500, 1000, 5000, 10000, 20000, 30000))
	r.Use(middlewares.Recoverer)

}

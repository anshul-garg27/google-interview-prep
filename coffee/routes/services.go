package routes

import (
	"coffee/app/heartbeat"
	"coffee/core/rest"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

func AttachServiceRoutes(r *chi.Mux, container rest.ApplicationContainer) {
	log.Debug("Attaching services endpoints")
	for _, api := range container.Services {
		api.AttachRoutes(r)
	}
}

func AttachHeartbeatRoutes(r *chi.Mux) {
	log.Debug("Attaching heartbeat endpoints")
	r.Get("/heartbeat/", heartbeat.Beat)
	r.Put("/heartbeat/", heartbeat.ModifyBeat)
}

func AttachMetricsRoutes(r *chi.Mux) {
	log.Debug("Attaching metrics endpoints")
	r.Handle("/metrics", promhttp.Handler())
}

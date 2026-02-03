package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"init.bulbul.tv/bulbul-backend/saas-gateway/config"
	"init.bulbul.tv/bulbul-backend/saas-gateway/context"
	"init.bulbul.tv/bulbul-backend/saas-gateway/header"
	"strconv"
	"strings"
	"sync"
)

var (
	collectorInit                sync.Once
	LoadPageEmptyResponseCounter *prometheus.CounterVec
	WinklRequestCounter          *prometheus.CounterVec
	LoadPageInvalidSlugCounter   *prometheus.CounterVec
	LoadPageErrorCounter         *prometheus.CounterVec
)

func SetupCollectors(config config.Config) {
	collectorInit.Do(func() {
		WinklRequestCounter = promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "requests_winkl",
			Help: "Track Requests via Winkl Proxy",
		}, []string{"method", "path", "channel", "clientId", "responseStatus"})

		LoadPageEmptyResponseCounter = promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "zero_edges_loadPage",
			Help: "Track # of loadPage calls with 0 edges",
		}, []string{"slug", "version", "channel", "pageNo"})

		LoadPageInvalidSlugCounter = promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "invalid_slug_loadPage",
			Help: "Track # of loadPage calls with invalid slug",
		}, []string{"slug", "version", "channel", "pageNo"})

		LoadPageErrorCounter = promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "errors_loadPage",
			Help: "Track # of loadPage errors",
		}, []string{"slug", "version", "channel", "pageNo"})
	})
}

func LoadPageEmptyResponseInc(gc *context.Context, slug string, pageNo int) {
	if LoadPageEmptyResponseCounter != nil {
		slugParts := strings.Split(slug, "?")
		LoadPageEmptyResponseCounter.With(prometheus.Labels{
			"slug":    slugParts[0],
			"channel": gc.XHeader.Get(header.Channel),
			"version": gc.XHeader.Get(header.VersionName),
			"pageNo":  strconv.Itoa(pageNo),
		}).Inc()
	}
}

func WinklRequestInc(gc *context.Context, method string, path string, responseStatus string) {
	if WinklRequestCounter != nil {
		WinklRequestCounter.With(prometheus.Labels{
			"method":         method,
			"path":           path,
			"responseStatus": responseStatus,
			"channel":        gc.XHeader.Get(header.Channel),
			"clientId":       gc.XHeader.Get(header.ClientID),
		}).Inc()
	}
}

func LoadPageInvalidSlugInc(gc *context.Context, slug string, pageNo int) {
	if LoadPageInvalidSlugCounter != nil {
		slugParts := strings.Split(slug, "?")
		LoadPageInvalidSlugCounter.With(prometheus.Labels{
			"slug":    slugParts[0],
			"channel": gc.XHeader.Get(header.Channel),
			"version": gc.XHeader.Get(header.VersionName),
			"pageNo":  strconv.Itoa(pageNo),
		}).Inc()
	}
}

func LoadPageErrorInc(gc *context.Context, slug string, pageNo int) {
	if LoadPageErrorCounter != nil {
		slugParts := strings.Split(slug, "?")
		LoadPageErrorCounter.With(prometheus.Labels{
			"slug":    slugParts[0],
			"channel": gc.XHeader.Get(header.Channel),
			"version": gc.XHeader.Get(header.VersionName),
			"pageNo":  strconv.Itoa(pageNo),
		}).Inc()
	}
}

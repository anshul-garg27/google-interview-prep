package rest

import "github.com/go-chi/chi/v5"

type ApiWrapper interface {
	AttachRoutes(r *chi.Mux)
	GetPrefix() string
}

package api

import (
	"coffee/core/rest"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type genreInsightsImpl struct {
	service *Service
}

func NewGenreInsightsImplApi(service *Service) *genreInsightsImpl {
	apiImpl := genreInsightsImpl{service: service}
	return &apiImpl
}

func (s genreInsightsImpl) GetPrefix() string {
	return "/genre-insights-service"
}

func (s genreInsightsImpl) AttachRoutes(r *chi.Mux) {
	genreBaseUrl := s.GetPrefix() + "/api/genre"
	r.Post(genreBaseUrl+"/search", s.SearchGenreOverview)
	r.Post(genreBaseUrl+"/{platform}/language_split", s.SearchLanguage)

}
func (s genreInsightsImpl) SearchGenreOverview(w http.ResponseWriter, r *http.Request) {
	sortBy, sortDir, page, size, query, done := rest.ParseSearchParameters(w, r)
	if done {
		return
	}
	resp := s.service.SearchGenreOverview(r.Context(), query, sortBy, sortDir, page, size)
	render.JSON(w, r, resp)
}

func (s genreInsightsImpl) SearchLanguage(w http.ResponseWriter, r *http.Request) {
	platform := chi.URLParam(r, "platform")
	sortBy, sortDir, page, size, query, done := rest.ParseSearchParameters(w, r)
	if done {
		return
	}
	resp := s.service.LanguageSearch(r.Context(), platform, query, sortBy, sortDir, page, size)
	render.JSON(w, r, resp)
}

package api

import (
	"coffee/core/rest"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type contentModuleImpl struct {
	service *Service
}

func NewContentModuleImplApi(service *Service) *contentModuleImpl {
	apiImpl := contentModuleImpl{service: service}
	return &apiImpl
}

func (s contentModuleImpl) GetPrefix() string {
	return "/content-service"
}

func (s contentModuleImpl) AttachRoutes(r *chi.Mux) {
	contentBaseUrl := s.GetPrefix() + "/api/content"
	r.Post(contentBaseUrl+"/trending_content", s.SearchTrendingContent)

}

func (s contentModuleImpl) SearchTrendingContent(w http.ResponseWriter, r *http.Request) {
	sortBy, sortDir, page, size, query, done := rest.ParseSearchParameters(w, r)
	if done {
		return
	}
	source := r.URL.Query().Get("source")
	resp := s.service.SearchTrendingContent(r.Context(), query, sortBy, sortDir, page, size, source)
	render.JSON(w, r, resp)
}

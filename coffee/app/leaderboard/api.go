package leaderboard

import (
	"coffee/core/rest"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type leaderBopardImpl struct {
	service *Service
}

func NewLeaderboardApi(service *Service) *leaderBopardImpl {
	apiImpl := leaderBopardImpl{
		service: service,
	}
	return &apiImpl
}

func (s leaderBopardImpl) GetPrefix() string {
	return "/leaderboard-service"
}

func (s leaderBopardImpl) AttachRoutes(r *chi.Mux) {
	profileBaseUrl := s.GetPrefix() + "/api/leaderboard"
	// SaaS
	r.Post(profileBaseUrl+"/{platform}", s.LeaderBoardByPlatform)

}

func (s leaderBopardImpl) LeaderBoardByPlatform(w http.ResponseWriter, r *http.Request) {
	platform := chi.URLParam(r, "platform")
	sortBy, sortDir, page, size, query, done := rest.ParseSearchParameters(w, r)
	if done {
		return
	}
	source := r.URL.Query().Get("source")

	resp := s.service.LeaderBoardByPlatform(r.Context(), platform, query, sortBy, sortDir, page, size, source)
	render.JSON(w, r, resp)
}

package api

import (
	"coffee/constants"
	"coffee/core/rest"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type socialDiscoveryImpl struct {
	service *Service
}

func NewSocialDiscoveryApi(service *Service) *socialDiscoveryImpl {
	apiImpl := socialDiscoveryImpl{service: service}
	return &apiImpl
}

func (s socialDiscoveryImpl) GetPrefix() string {
	return "/discovery-service"
}

func (s socialDiscoveryImpl) AttachRoutes(r *chi.Mux) {
	profileBaseUrl := s.GetPrefix() + "/api/profile"
	// SaaS
	r.Get(profileBaseUrl+"/{profileId}", s.FindByProfileId)
	r.Get(profileBaseUrl+"/{platform}/{platformProfileId}", s.FindByPlatformProfileId)
	r.Get(profileBaseUrl+"/{platform}/{platformProfileId}/audience", s.FindAudienceByPlatformProfileId)
	r.Post(profileBaseUrl+"/{platform}/{platformProfileId}/timeseries", s.FindTimeSeriesByPlatformProfileId)
	r.Post(profileBaseUrl+"/{platform}/{platformProfileId}/content", s.FindContentByPlatformProfileId)
	r.Get(profileBaseUrl+"/{platform}/{platformProfileId}/hashtags", s.FindHashTagsByPlatformProfileId)
	r.Post(profileBaseUrl+"/{platform}/{platformProfileId}/similar_accounts", s.FindSimilarAccountsByPlatformProfileId)
	r.Post(profileBaseUrl+"/{platform}/search", s.SearchByPlatform)
	r.Post(profileBaseUrl+"/locations", s.SearchLocations)
	r.Post(profileBaseUrl+"/locations_new", s.SearchLocationsNew)
	r.Get(profileBaseUrl+"/gcc_locations", s.SearchGccLocations)

	// GCC
	r.Get(profileBaseUrl+"/{platform}/byhandle/{handle}", s.FindByPlatformHandle)

}

func (s socialDiscoveryImpl) FindByProfileId(w http.ResponseWriter, r *http.Request) {
	profileId := chi.URLParam(r, "profileId")
	linkedSource := r.URL.Query().Get("source")
	resp := s.service.FindByProfileId(r.Context(), profileId, linkedSource)
	render.JSON(w, r, resp)
}

func (s socialDiscoveryImpl) FindByPlatformProfileId(w http.ResponseWriter, r *http.Request) {
	platform := chi.URLParam(r, "platform")
	platformProfileIdStr := chi.URLParam(r, "platformProfileId")
	linkedSource := r.URL.Query().Get("source")
	platformProfileId, err := strconv.ParseInt(platformProfileIdStr, 10, 64)
	if err != nil {
		rest.RenderError(w, r, err)
		return
	}
	campaignProfileJoinStr := r.URL.Query().Get("campaignProfileJoin")
	campaignProfileJoin, err := strconv.ParseBool(campaignProfileJoinStr)
	if err != nil {
		campaignProfileJoin = false
	}
	resp := s.service.FindByPlatformProfileId(r.Context(), platform, platformProfileId, linkedSource, campaignProfileJoin)
	render.JSON(w, r, resp)
}

func (s socialDiscoveryImpl) FindAudienceByPlatformProfileId(w http.ResponseWriter, r *http.Request) {
	platform := chi.URLParam(r, "platform")
	platformProfileIdStr := chi.URLParam(r, "platformProfileId")
	platformProfileId, err := strconv.ParseInt(platformProfileIdStr, 10, 64)
	if err != nil {
		rest.RenderError(w, r, err)
		return
	}
	resp := s.service.FindAudienceByPlatformProfileId(r.Context(), platform, platformProfileId)
	render.JSON(w, r, resp)
}

func (s socialDiscoveryImpl) FindTimeSeriesByPlatformProfileId(w http.ResponseWriter, r *http.Request) {
	platform := chi.URLParam(r, "platform")
	platformProfileIdStr := chi.URLParam(r, "platformProfileId")
	platformProfileId, err := strconv.ParseInt(platformProfileIdStr, 10, 64)

	sortBy, sortDir, page, size, query, done := rest.ParseSearchParameters(w, r)
	if done {
		return
	}
	if err != nil {
		rest.RenderError(w, r, err)
		return
	}
	resp := s.service.FindTimeSeriesByPlatformProfileId(r.Context(), platform, platformProfileId, query, sortBy, sortDir, page, size)
	render.JSON(w, r, resp)
}

func (s socialDiscoveryImpl) FindContentByPlatformProfileId(w http.ResponseWriter, r *http.Request) {
	platform := chi.URLParam(r, "platform")
	platformProfileIdStr := chi.URLParam(r, "platformProfileId")
	platformProfileId, err := strconv.ParseInt(platformProfileIdStr, 10, 64)

	sortBy, sortDir, page, size, query, done := rest.ParseSearchParameters(w, r)
	if done {
		return
	}
	if err != nil {
		rest.RenderError(w, r, err)
		return
	}
	linkedSource := r.URL.Query().Get("source")

	resp := s.service.FindContentByPlatformProfileId(r.Context(), platform, platformProfileId, query, sortBy, sortDir, page, size, linkedSource)
	render.JSON(w, r, resp)
}

func (s socialDiscoveryImpl) FindHashTagsByPlatformProfileId(w http.ResponseWriter, r *http.Request) {
	platform := chi.URLParam(r, "platform")
	platformProfileIdStr := chi.URLParam(r, "platformProfileId")
	platformProfileId, err := strconv.ParseInt(platformProfileIdStr, 10, 64)
	if err != nil {
		rest.RenderError(w, r, err)
		return
	}
	resp := s.service.FindHashTagsByPlatformProfileId(r.Context(), platform, platformProfileId)
	render.JSON(w, r, resp)
}

func (s socialDiscoveryImpl) FindSimilarAccountsByPlatformProfileId(w http.ResponseWriter, r *http.Request) {
	platform := chi.URLParam(r, "platform")
	platformProfileIdStr := chi.URLParam(r, "platformProfileId")
	platformProfileId, err := strconv.ParseInt(platformProfileIdStr, 10, 64)
	sortBy, sortDir, page, size, query, done := rest.ParseSearchParameters(w, r)
	if done {
		return
	}
	if err != nil {
		rest.RenderError(w, r, err)
		return
	}
	linkedSource := r.URL.Query().Get("source")
	resp := s.service.FindSimilarAccountsByPlatformProfileId(r.Context(), platform, platformProfileId, query, sortBy, sortDir, page, size, linkedSource)
	render.JSON(w, r, resp)
}

func (s socialDiscoveryImpl) SearchByPlatform(w http.ResponseWriter, r *http.Request) {
	var computeCounts bool
	platform := chi.URLParam(r, "platform")
	linkedSource := r.URL.Query().Get("source")
	computeCountsStr := r.URL.Query().Get("computeCounts")
	if computeCountsStr != "" {
		parsedValue, err := strconv.ParseBool(computeCountsStr)
		if err == nil {
			computeCounts = parsedValue
		} else {
			computeCounts = true
		}
	} else {
		computeCounts = true
	}

	campaignProfileJoinStr := r.URL.Query().Get("campaignProfileJoin")
	campaignProfileJoin, _ := strconv.ParseBool(campaignProfileJoinStr)

	if linkedSource == string(constants.GCCPlatform) {
		campaignProfileJoin = true
	}
	activityId := r.URL.Query().Get("activityId")
	_, err := strconv.Atoi(activityId)
	if err != nil {
		activityId = ""
	}
	sortBy, sortDir, page, size, query, done := rest.ParseSearchParameters(w, r)
	if done {
		return
	}

	if sortBy == "followers_lifetime" {
		sortBy = "followers"
	}
	resp := s.service.SearchByPlatform(r.Context(), platform, query, sortBy, sortDir, page, size, linkedSource, campaignProfileJoin, activityId, computeCounts)
	render.JSON(w, r, resp)
}

func (s socialDiscoveryImpl) SearchLocations(w http.ResponseWriter, r *http.Request) {
	sortBy, sortDir, page, size, query, done := rest.ParseSearchParameters(w, r)
	if done {
		return
	}
	resp := s.service.SearchLocations(r.Context(), query, sortBy, sortDir, page, size)
	render.JSON(w, r, resp)
}

func (s socialDiscoveryImpl) FindByPlatformHandle(w http.ResponseWriter, r *http.Request) {
	platform := chi.URLParam(r, "platform")
	handle := chi.URLParam(r, "handle")
	lnikedSource := r.URL.Query().Get("source")
	campaignProfileJoinStr := r.URL.Query().Get("campaignProfileJoin")
	campaignProfileJoin, err := strconv.ParseBool(campaignProfileJoinStr)
	if err != nil {
		campaignProfileJoin = false
	}
	resp := s.service.FindByPlatformHandle(r.Context(), platform, handle, lnikedSource, campaignProfileJoin)
	render.JSON(w, r, resp)
}

func (s socialDiscoveryImpl) SearchGccLocations(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	resp := s.service.SearchGccLocations(r.Context(), name)
	render.JSON(w, r, resp)
}

func (s socialDiscoveryImpl) SearchLocationsNew(w http.ResponseWriter, r *http.Request) {
	_, _, _, size, query, done := rest.ParseSearchParameters(w, r)
	if done {
		return
	}
	resp := s.service.SearchLocationsNew(r.Context(), query, size)
	render.JSON(w, r, resp)
}

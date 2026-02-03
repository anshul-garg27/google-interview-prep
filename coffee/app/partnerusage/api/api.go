package api

import (
	"coffee/app/partnerusage/domain"
	"coffee/constants"
	"coffee/core/appcontext"
	"coffee/core/rest"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type partnerUsageApiImpl struct {
	service *Service
}

func (s partnerUsageApiImpl) GetPrefix() string {
	return "/partner-usage-service"
}

func (s partnerUsageApiImpl) AttachRoutes(r *chi.Mux) {
	contractBaseUrl := s.GetPrefix() + "/api/partner"
	r.Get(contractBaseUrl+"/{id}", s.FindByPartnerId)
	r.Put(contractBaseUrl+"/{id}", s.Update)
	r.Post(contractBaseUrl+"/search", s.Search)
	r.Post(contractBaseUrl+"/profile/{PlatformCode}/{platform}/unlock", s.UnlockProfile)
}

func NewPartnerUsageApi(service *Service) *partnerUsageApiImpl {
	apiImpl := partnerUsageApiImpl{service: service}
	return &apiImpl
}

func (s partnerUsageApiImpl) FindByPartnerId(w http.ResponseWriter, r *http.Request) {
	partnerIdStr := chi.URLParam(r, "id")
	partnerId, _ := strconv.ParseInt(partnerIdStr, 10, 64)
	resp := s.service.FindByPartnerId(r.Context(), partnerId)

	render.JSON(w, r, resp)
}
func (s partnerUsageApiImpl) Update(w http.ResponseWriter, r *http.Request) {
	var partnerContract domain.PartnerInput
	err := json.NewDecoder(r.Body).Decode(&partnerContract)
	if err != nil {
		rest.RenderError(w, r, err)
		return
	}
	partnerIdStr := chi.URLParam(r, "id")
	partnerId, _ := strconv.ParseInt(partnerIdStr, 10, 64)

	resp := s.service.UpdatePartnerUsage(r.Context(), partnerId, &partnerContract)
	render.JSON(w, r, resp)

}

func (s partnerUsageApiImpl) Search(w http.ResponseWriter, r *http.Request) {
	sortBy, sortDir, page, size, query, done := rest.ParseSearchParameters(w, r)
	if done {
		return
	}
	// var searchQuery coredomain.SearchQuery
	resp := s.service.Search(r.Context(), query, sortBy, sortDir, page, size)
	render.JSON(w, r, resp)
}

func (s partnerUsageApiImpl) UnlockProfile(w http.ResponseWriter, r *http.Request) {
	PlatformCode := chi.URLParam(r, "PlatformCode")
	platform := chi.URLParam(r, "platform")
	ctx := r.Context()
	appCtx := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)
	if appCtx.PartnerId == nil {
		rest.RenderError(w, r, errors.New("you are not authorized to access this resource"))
		return
	}
	resp := s.service.UnlockProfile(r.Context(), appCtx.PartnerId, PlatformCode, platform)
	render.JSON(w, r, resp)

}

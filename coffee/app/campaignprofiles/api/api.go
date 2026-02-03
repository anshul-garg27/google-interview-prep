package api

import (
	"coffee/constants"
	"coffee/core/rest"
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strconv"

	log "github.com/sirupsen/logrus"

	coredomain "coffee/core/domain"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type CampaignProfileImpl struct {
	service *Service
}

func (s CampaignProfileImpl) GetPrefix() string {
	return "/campaign-profile-service"
}

func (s CampaignProfileImpl) AttachRoutes(r *chi.Mux) {
	campaignBaseUrl := s.GetPrefix() + "/api/campaign-profile"
	r.Get(campaignBaseUrl+"/{id}", s.FindById)
	r.Get(campaignBaseUrl+"/handle/{platform}/{handle}", s.FindCampaignProfileByPlatformHandle)
	r.Post(campaignBaseUrl+"/upsert", s.UpsertByHandle)
	r.Post(campaignBaseUrl+"/upsert/{profileCode}", s.UpsertUsingProfileCode)
	r.Post(campaignBaseUrl+"/search", s.Search)
	r.Post(campaignBaseUrl+"/{id}/update", s.UpdateById)
	r.Put(campaignBaseUrl+"/{id}/refresh", s.RefreshSocialProfileById)
	r.Get(campaignBaseUrl+"/{id}/{platform}/insights", s.CreatorInsightsById)
	r.Get(campaignBaseUrl+"/{id}/{platform}/audienceinsights", s.CreatorAudienceInsightsDataById)
}

func NewCampaignApi(service *Service) *CampaignProfileImpl {
	apiImpl := CampaignProfileImpl{service: service}
	return &apiImpl
}

func (s CampaignProfileImpl) UpsertUsingProfileCode(w http.ResponseWriter, r *http.Request) {
	var input coredomain.CampaignProfileInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		rest.RenderError(w, r, err)
		return
	}
	requestBodyJSON, _ := json.Marshal(input)
	requestHeadersJSON, _ := json.Marshal(r.Header)
	log.Debug("Request header: ", string(requestHeadersJSON))
	log.Debug("Request URL: ", r.URL.String())
	log.Debug("Request to Upsert: ", string(requestBodyJSON))
	linkedSource := r.URL.Query().Get("source")
	profileCode := chi.URLParam(r, "profileCode")
	if linkedSource == "" {
		linkedSource = string(constants.GCCPlatform)
	}

	resp := s.service.UpsertUsingProfileCode(r.Context(), profileCode, input, linkedSource)
	render.JSON(w, r, resp)
}

func (s CampaignProfileImpl) FindCampaignProfileByPlatformHandle(w http.ResponseWriter, r *http.Request) {
	log.Debug("FindCampaignProfileByPlatformHandle")
	platform := chi.URLParam(r, "platform")
	platformHandle := chi.URLParam(r, "handle")
	if !isValidHandle(platformHandle, platform) {
		err := errors.New("handle invalid")
		rest.RenderError(w, r, err)
		return
	}
	linkedSource := r.URL.Query().Get("source")
	fullRefresh := r.URL.Query().Get("fullRefresh")
	fullRefreshFlag := false
	if fullRefresh == "1" {
		fullRefreshFlag = true
	}
	resp := s.service.FindCampaignProfileByPlatformHandle(r.Context(), platform, platformHandle, linkedSource, fullRefreshFlag)
	render.JSON(w, r, resp)
}

func (s CampaignProfileImpl) FindById(w http.ResponseWriter, r *http.Request) {
	fetchSocials := GetSocialFlagFromReq(r)
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	resp := s.service.FindById(r.Context(), id, fetchSocials)
	render.JSON(w, r, resp)
}

func (s CampaignProfileImpl) Search(w http.ResponseWriter, r *http.Request) {
	sortBy, sortDir, page, size, query, done := rest.ParseSearchParameters(w, r)
	if done {
		return
	}
	fetchSocials := GetSocialFlagFromReq(r)
	resp := s.service.Search(r.Context(), query, sortBy, sortDir, page, size, fetchSocials)
	render.JSON(w, r, resp)
}

func (s CampaignProfileImpl) UpdateById(w http.ResponseWriter, r *http.Request) {
	var cpInput coredomain.CampaignProfileInput
	err := json.NewDecoder(r.Body).Decode(&cpInput)
	if err != nil {
		rest.RenderError(w, r, err)
		return
	}
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	linkedSource := r.URL.Query().Get("source")
	resp := s.service.UpdateById(r.Context(), id, cpInput, linkedSource)
	render.JSON(w, r, resp)
}

func (s CampaignProfileImpl) UpsertByHandle(w http.ResponseWriter, r *http.Request) {
	var input coredomain.CampaignProfileInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		rest.RenderError(w, r, err)
		return
	}
	requestBodyJSON, _ := json.Marshal(input)
	requestHeadersJSON, _ := json.Marshal(r.Header)
	log.Debug("Request header: ", string(requestHeadersJSON))
	log.Debug("Request URL: ", r.URL.String())
	log.Debug("Request to Upsert: ", string(requestBodyJSON))
	linkedSource := r.URL.Query().Get("source")
	if linkedSource == "" {
		linkedSource = string(constants.GCCPlatform)
	}

	resp := s.service.UpsertByHandle(r.Context(), input, linkedSource)
	render.JSON(w, r, resp)
}

func GetSocialFlagFromReq(r *http.Request) bool {
	valueFetchSocials := r.URL.Query().Get("fetchSocials")
	fetchSocials, _ := strconv.ParseBool(valueFetchSocials)
	if !fetchSocials {
		value := r.URL.Query().Get("socials")
		fetchSocials, _ = strconv.ParseBool(value)
	}
	return fetchSocials
}

func (s CampaignProfileImpl) RefreshSocialProfileById(w http.ResponseWriter, r *http.Request) {

	idStr := chi.URLParam(r, "id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	resp := s.service.RefreshSocialProfileById(r.Context(), id)
	render.JSON(w, r, resp)
}

func (s CampaignProfileImpl) CreatorInsightsById(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	platform := chi.URLParam(r, "platform")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	token := r.URL.Query().Get("token")
	userId := r.URL.Query().Get("userId")
	resp := s.service.CreatorInsightsById(r.Context(), id, token, platform, userId)
	render.JSON(w, r, resp)
}

func (s CampaignProfileImpl) CreatorAudienceInsightsDataById(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	platform := chi.URLParam(r, "platform")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	token := r.URL.Query().Get("token")
	userId := r.URL.Query().Get("userId")
	resp := s.service.CreatorAudienceInsightsDataById(r.Context(), id, token, platform, userId)
	render.JSON(w, r, resp)
}

func isValidHandle(handle string, platform string) bool {

	pattern := ""
	if platform == "INSTAGRAM" {
		pattern = "^[a-zA-Z0-9_.]+$"
	} else {
		pattern = "^[a-zA-Z0-9_.-]+$"
	}
	re := regexp.MustCompile(pattern)
	return re.MatchString(handle) && handle != "null" && handle != "NULL"
}

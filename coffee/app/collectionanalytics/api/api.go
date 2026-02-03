package api

import (
	coredomain "coffee/core/domain"
	"coffee/core/rest"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type CollectionAnalyticsApiImpl struct {
	Service *CollectionAnalyticsService
}

func NewPostCollectionApi(service *CollectionAnalyticsService) *CollectionAnalyticsApiImpl {
	apiImpl := CollectionAnalyticsApiImpl{
		Service: service,
	}
	return &apiImpl
}

func (s CollectionAnalyticsApiImpl) GetPrefix() string {
	return "/collection-analytics-service"
}

func (api CollectionAnalyticsApiImpl) AttachRoutes(r *chi.Mux) {
	baseUrl := api.GetPrefix() + "/api/collection"
	r.Post(baseUrl+"/summary", api.FetchCollectionMetricsSummary)
	r.Post(baseUrl+"/profiles", api.FetchCollectionProfilesWithMetricsSummary)
	r.Post(baseUrl+"/posts", api.FetchCollectionPostsWithMetricsSummary)
	r.Post(baseUrl+"/timeseries", api.FetchCollectionTimeSeries)
	r.Post(baseUrl+"/timeseries/clicks", api.FetchCollectionClicksTimeSeries)
	r.Post(baseUrl+"/sentiment", api.FetchCollectionSentiment)

	r.Get(baseUrl+"/{collectionId}/{collectionType}/hashtags", api.FetchCollectionHashtags)
	r.Get(baseUrl+"/{collectionId}/{collectionType}/keywords", api.FetchCollectionKeywords)
	r.Get(baseUrl+"/byshareid/{shareId}/{collectionType}/hashtags", api.FetchCollectionHashtagsShared)
	r.Get(baseUrl+"/byshareid/{shareId}/{collectionType}/keywords", api.FetchCollectionKeywordsShared)

	r.Get(baseUrl+"/group/{collectionGroupId}/{collectionType}/hashtags", api.FetchCollectionGroupHashtags)
	r.Get(baseUrl+"/group/{collectionGroupId}/{collectionType}/keywords", api.FetchCollectionGroupKeywords)
	r.Get(baseUrl+"/group/byshareid/{collectionGroupShareId}/{collectionType}/hashtags", api.FetchCollectionGroupHashtagsShared)
	r.Get(baseUrl+"/group/byshareid/{collectionGroupShareId}/{collectionType}/keywords", api.FetchCollectionGroupKeywordsShared)
}

func (api CollectionAnalyticsApiImpl) FetchCollectionMetricsSummary(w http.ResponseWriter, r *http.Request) {
	if r.Body == http.NoBody {
		rest.RenderError(w, r, errors.New("Bad Request"))
	}
	var query coredomain.SearchQuery
	err := json.NewDecoder(r.Body).Decode(&query)
	if err != nil {
		rest.RenderError(w, r, err)
	}
	resp := api.Service.FetchCollectionMetricsSummary(r.Context(), query)
	render.JSON(w, r, resp)
}

func (api CollectionAnalyticsApiImpl) FetchCollectionTimeSeries(w http.ResponseWriter, r *http.Request) {
	if r.Body == http.NoBody {
		rest.RenderError(w, r, errors.New("Bad Request"))
	}
	var query coredomain.SearchQuery
	err := json.NewDecoder(r.Body).Decode(&query)
	if err != nil {
		rest.RenderError(w, r, err)
	}
	resp := api.Service.FetchCollectionTimeSeries(r.Context(), query)
	render.JSON(w, r, resp)
}

func (api CollectionAnalyticsApiImpl) FetchCollectionClicksTimeSeries(w http.ResponseWriter, r *http.Request) {
	if r.Body == http.NoBody {
		rest.RenderError(w, r, errors.New("Bad Request"))
	}
	var query coredomain.SearchQuery
	err := json.NewDecoder(r.Body).Decode(&query)
	if err != nil {
		rest.RenderError(w, r, err)
	}
	resp := api.Service.FetchCollectionClicksTimeSeries(r.Context(), query)
	render.JSON(w, r, resp)
}

func (api CollectionAnalyticsApiImpl) FetchCollectionProfilesWithMetricsSummary(w http.ResponseWriter, r *http.Request) {
	sortBy, sortDir, page, size, query, done := rest.ParseSearchParameters(w, r)
	if done {
		return
	}
	if sortBy == "id" {
		sortBy = "reach"
	}
	if sortBy != "views" && sortBy != "reach" && sortBy != "plays" && sortBy != "er" && sortBy != "likes" {
		sortBy = "reach"
	}
	resp := api.Service.FetchCollectionProfilesWithMetricsSummary(r.Context(), query, sortBy, sortDir, page, size)
	render.JSON(w, r, resp)
}

func (api CollectionAnalyticsApiImpl) FetchCollectionPostsWithMetricsSummary(w http.ResponseWriter, r *http.Request) {
	sortBy, sortDir, page, size, query, done := rest.ParseSearchParameters(w, r)
	if done {
		return
	}
	if sortBy == "id" {
		sortBy = "reach"
	}
	if sortBy != "views" && sortBy != "reach" && sortBy != "plays" && sortBy != "er" && sortBy != "likes" {
		sortBy = "reach"
	}
	resp := api.Service.FetchCollectionPostsWithMetricsSummary(r.Context(), query, sortBy, sortDir, page, size)
	render.JSON(w, r, resp)
}

func (api CollectionAnalyticsApiImpl) FetchCollectionHashtags(w http.ResponseWriter, r *http.Request) {
	collectionId := chi.URLParam(r, "collectionId")
	collectionType := chi.URLParam(r, "collectionType")
	resp := api.Service.FetchCollectionHashtags(r.Context(), collectionId, collectionType)
	render.JSON(w, r, resp)
}

func (api CollectionAnalyticsApiImpl) FetchCollectionKeywords(w http.ResponseWriter, r *http.Request) {
	collectionId := chi.URLParam(r, "collectionId")
	collectionType := chi.URLParam(r, "collectionType")
	resp := api.Service.FetchCollectionKeywords(r.Context(), collectionId, collectionType)
	render.JSON(w, r, resp)
}

func (api CollectionAnalyticsApiImpl) FetchCollectionHashtagsShared(w http.ResponseWriter, r *http.Request) {
	shareId := chi.URLParam(r, "shareId")
	collectionType := chi.URLParam(r, "collectionType")
	resp := api.Service.FetchCollectionHashtagsShared(r.Context(), shareId, collectionType)
	render.JSON(w, r, resp)
}

func (api CollectionAnalyticsApiImpl) FetchCollectionKeywordsShared(w http.ResponseWriter, r *http.Request) {
	shareId := chi.URLParam(r, "shareId")
	collectionType := chi.URLParam(r, "collectionType")
	resp := api.Service.FetchCollectionKeywordsShared(r.Context(), shareId, collectionType)
	render.JSON(w, r, resp)
}

func (api CollectionAnalyticsApiImpl) FetchCollectionGroupHashtags(w http.ResponseWriter, r *http.Request) {
	collectionGroupIdStr := chi.URLParam(r, "collectionGroupId")
	collectionType := chi.URLParam(r, "collectionType")
	collectionGroupId, _ := strconv.ParseInt(collectionGroupIdStr, 10, 64)
	resp := api.Service.FetchCollectionGroupHashtags(r.Context(), collectionGroupId, collectionType)
	render.JSON(w, r, resp)
}

func (api CollectionAnalyticsApiImpl) FetchCollectionGroupKeywords(w http.ResponseWriter, r *http.Request) {
	collectionGroupIdStr := chi.URLParam(r, "collectionGroupId")
	collectionType := chi.URLParam(r, "collectionType")
	collectionGroupId, _ := strconv.ParseInt(collectionGroupIdStr, 10, 64)
	resp := api.Service.FetchCollectionGroupKeywords(r.Context(), collectionGroupId, collectionType)
	render.JSON(w, r, resp)
}

func (api CollectionAnalyticsApiImpl) FetchCollectionGroupHashtagsShared(w http.ResponseWriter, r *http.Request) {
	shareId := chi.URLParam(r, "shareId")
	collectionType := chi.URLParam(r, "collectionType")
	resp := api.Service.FetchCollectionGroupHashtagsShared(r.Context(), shareId, collectionType)
	render.JSON(w, r, resp)
}

func (api CollectionAnalyticsApiImpl) FetchCollectionGroupKeywordsShared(w http.ResponseWriter, r *http.Request) {
	shareId := chi.URLParam(r, "shareId")
	collectionType := chi.URLParam(r, "collectionType")
	resp := api.Service.FetchCollectionGroupKeywordsShared(r.Context(), shareId, collectionType)
	render.JSON(w, r, resp)
}

func (api CollectionAnalyticsApiImpl) FetchCollectionSentiment(w http.ResponseWriter, r *http.Request) {
	if r.Body == http.NoBody {
		rest.RenderError(w, r, errors.New("Bad Request"))
	}
	var query coredomain.SearchQuery
	err := json.NewDecoder(r.Body).Decode(&query)
	if err != nil {
		rest.RenderError(w, r, err)
	}
	resp := api.Service.FetchCollectionSentiment(r.Context(), query)
	render.JSON(w, r, resp)
}

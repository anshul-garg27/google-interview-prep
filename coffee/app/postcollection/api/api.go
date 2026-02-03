package api

import (
	"coffee/app/postcollection/domain"
	"coffee/constants"
	"coffee/core/appcontext"
	coredomain "coffee/core/domain"
	"coffee/core/rest"
	"coffee/helpers"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"go.step.sm/crypto/randutil"
)

type PostCollectionApiImpl struct {
	Service *PostCollectionService
}

func NewPostCollectionApi(service *PostCollectionService) *PostCollectionApiImpl {
	apiImpl := PostCollectionApiImpl{
		Service: service,
	}
	return &apiImpl
}

func (s PostCollectionApiImpl) GetPrefix() string {
	return "/post-collection-service"
}

func (api PostCollectionApiImpl) AttachRoutes(r *chi.Mux) {
	baseUrl := api.GetPrefix() + "/api/collection"
	// post collection
	r.Post(baseUrl+"/", api.CreateCollection)
	r.Put(baseUrl+"/{id}", api.UpdateCollection)
	r.Get(baseUrl+"/{id}", api.FindCollectionById)
	r.Get(baseUrl+"/byShareId/{shareId}", api.FindCollectionByShareId)
	r.Put(baseUrl+"/{id}/renewShareId", api.RenewCollectionShareId)
	r.Delete(baseUrl+"/{id}", api.DisableCollection)
	r.Post(baseUrl+"/search", api.SearchCollections)

	// post collection items
	r.Post(baseUrl+"/{collectionId}/item/", api.AddItemsToCollection)
	r.Delete(baseUrl+"/{collectionId}/item", api.DeleteItemsFromCollection)
	r.Put(baseUrl+"/{collectionId}/item/{itemId}", api.UpdatePostCollectionItemById)
	r.Put(baseUrl+"/{collectionId}/item/byShortCode/{platform}/{shortCode}", api.UpdatePostCollectionItemByPlatformAndShortCode)
	r.Delete(baseUrl+"/{collectionId}/item/byShortCode", api.DeleteItemsFromCollectionByShortCode)
	r.Post(baseUrl+"/{collectionId}/item/search", api.SearchPostCollectionItems)

}

func (api PostCollectionApiImpl) CreateCollection(w http.ResponseWriter, r *http.Request) {
	var input domain.PostCollectionEntry
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		rest.RenderError(w, r, err)
		return
	}

	ctx := r.Context()
	appCtx := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)
	if appCtx.PartnerId == nil || appCtx.AccountId == nil {
		rest.RenderError(w, r, errors.New("You are not authorized to perform this action"))
		return
	}

	input.PartnerId = appCtx.PartnerId
	accountIdStr := strconv.FormatInt(*appCtx.AccountId, 10)
	input.CreatedBy = &accountIdStr
	input.Enabled = helpers.ToBool(true)
	if input.ShareId == nil {
		shareId, _ := randutil.Alphanumeric(7)
		input.ShareId = &shareId
	}
	if input.MetricsIngestionFreq == nil {
		defaultFreq := "default"
		input.MetricsIngestionFreq = &defaultFreq
	}
	input.DisabledMetrics = []string{}
	input.ExpectedMetricValues = make(map[string]*float64)
	resp := api.Service.Create(r.Context(), &input)
	render.JSON(w, r, resp)
}

func (api PostCollectionApiImpl) UpdateCollection(w http.ResponseWriter, r *http.Request) {
	var input domain.PostCollectionEntry
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		rest.RenderError(w, r, err)
		return
	}

	ctx := r.Context()
	appCtx := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)
	if appCtx.PartnerId == nil || appCtx.AccountId == nil {
		rest.RenderError(w, r, errors.New("You are not authorized to update this collection"))
		return
	}

	id := chi.URLParam(r, "id")
	resp := api.Service.Update(r.Context(), id, &input)
	render.JSON(w, r, resp)
}

func (api PostCollectionApiImpl) FindCollectionById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	appCtx := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)
	if appCtx.PartnerId == nil || appCtx.AccountId == nil {
		rest.RenderError(w, r, errors.New("You are not authorized to update this collection"))
		return
	}

	id := chi.URLParam(r, "id")
	resp := api.Service.FindById(r.Context(), id)
	render.JSON(w, r, resp)
}

func (api PostCollectionApiImpl) FindCollectionByShareId(w http.ResponseWriter, r *http.Request) {
	shareId := chi.URLParam(r, "shareId")
	resp := api.Service.FindByShareId(r.Context(), shareId)
	render.JSON(w, r, resp)
}

func (api PostCollectionApiImpl) RenewCollectionShareId(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	appCtx := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)
	if appCtx.PartnerId == nil || appCtx.AccountId == nil {
		rest.RenderError(w, r, errors.New("You are not authorized to update this collection"))
		return
	}

	id := chi.URLParam(r, "id")
	shareId, _ := randutil.Alphanumeric(7)
	input := domain.PostCollectionEntry{ShareId: &shareId}
	resp := api.Service.Service.Update(r.Context(), id, &input)
	render.JSON(w, r, resp)
}

func (api PostCollectionApiImpl) DisableCollection(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	appCtx := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)
	if appCtx.PartnerId == nil || appCtx.AccountId == nil {
		rest.RenderError(w, r, errors.New("You are not authorized to update this collection"))
		return
	}

	id := chi.URLParam(r, "id")
	input := domain.PostCollectionEntry{Enabled: helpers.ToBool(false)}
	resp := api.Service.Service.Update(r.Context(), id, &input)
	render.JSON(w, r, resp)
}

func (api PostCollectionApiImpl) SearchCollections(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	appCtx := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)
	if appCtx.PartnerId == nil || appCtx.AccountId == nil {
		rest.RenderError(w, r, errors.New("You are not authorized to perform this action"))
		return
	}
	sortBy, sortDir, page, size, query, done := rest.ParseSearchParameters(w, r)
	if done {
		return
	}
	query.Filters = append(query.Filters, coredomain.SearchFilter{
		FilterType: "EQ",
		Field:      "partner_id",
		Value:      strconv.FormatInt(*appCtx.PartnerId, 10),
	})
	resp := api.Service.SearchCollections(r.Context(), query, sortBy, sortDir, page, size)
	render.JSON(w, r, resp)
}

// Item Level Apis

func (api PostCollectionApiImpl) AddItemsToCollection(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	appCtx := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)
	if appCtx.PartnerId == nil || appCtx.AccountId == nil {
		rest.RenderError(w, r, errors.New("You are not authorized to perform this action"))
		return
	}
	var input []domain.PostCollectionItemEntry
	collectionId := chi.URLParam(r, "collectionId")
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		rest.RenderError(w, r, err)
		return
	}
	resp := api.Service.AddItemsToCollection(r.Context(), collectionId, input)
	render.JSON(w, r, resp)
}

func (api PostCollectionApiImpl) DeleteItemsFromCollection(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	appCtx := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)
	if appCtx.PartnerId == nil || appCtx.AccountId == nil {
		rest.RenderError(w, r, errors.New("You are not authorized to perform this action"))
		return
	}
	var input []domain.PostCollectionItemEntry
	collectionId := chi.URLParam(r, "collectionId")
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		rest.RenderError(w, r, err)
		return
	}
	resp := api.Service.DeleteItemsFromCollection(r.Context(), collectionId, input)
	render.JSON(w, r, resp)
}

func (api PostCollectionApiImpl) UpdatePostCollectionItemById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	appCtx := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)
	if appCtx.PartnerId == nil || appCtx.AccountId == nil {
		rest.RenderError(w, r, errors.New("You are not authorized to perform this action"))
		return
	}
	var input domain.PostCollectionItemEntry
	collectionId := chi.URLParam(r, "collectionId")
	itemIdStr := chi.URLParam(r, "itemId")
	itemId, err := strconv.ParseInt(itemIdStr, 10, 64)
	if err != nil {
		rest.RenderError(w, r, err)
		return
	}
	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		rest.RenderError(w, r, err)
		return
	}
	resp := api.Service.UpdatePostCollectionItemById(r.Context(), collectionId, itemId, input)
	render.JSON(w, r, resp)
}

func (api PostCollectionApiImpl) UpdatePostCollectionItemByPlatformAndShortCode(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	appCtx := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)
	if appCtx.PartnerId == nil || appCtx.AccountId == nil {
		rest.RenderError(w, r, errors.New("You are not authorized to perform this action"))
		return
	}
	var input domain.PostCollectionItemEntry
	collectionId := chi.URLParam(r, "collectionId")
	platform := chi.URLParam(r, "platform")
	shortCode := chi.URLParam(r, "shortCode")
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		rest.RenderError(w, r, err)
		return
	}
	resp := api.Service.UpdatePostCollectionItemByShortCode(r.Context(), collectionId, platform, shortCode, input)
	render.JSON(w, r, resp)
}

func (api PostCollectionApiImpl) DeleteItemsFromCollectionByShortCode(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	appCtx := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)
	if appCtx.PartnerId == nil || appCtx.AccountId == nil {
		rest.RenderError(w, r, errors.New("You are not authorized to perform this action"))
		return
	}
	var input []domain.PostCollectionItemEntry
	collectionId := chi.URLParam(r, "collectionId")
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		rest.RenderError(w, r, err)
		return
	}
	resp := api.Service.DeleteItemsFromCollectionByShortCode(r.Context(), collectionId, input)
	render.JSON(w, r, resp)
}

func (api PostCollectionApiImpl) SearchPostCollectionItems(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	appCtx := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)
	if appCtx.PartnerId == nil || appCtx.AccountId == nil {
		rest.RenderError(w, r, errors.New("You are not authorized to perform this action"))
		return
	}
	collectionId := chi.URLParam(r, "collectionId")
	sortBy, sortDir, page, size, query, done := rest.ParseSearchParameters(w, r)
	if done {
		return
	}
	resp := api.Service.SearchPostCollectionItems(r.Context(), collectionId, query, sortBy, sortDir, page, size)
	render.JSON(w, r, resp)
}

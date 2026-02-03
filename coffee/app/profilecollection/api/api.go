package api

import (
	"coffee/app/profilecollection/domain"
	"coffee/constants"
	"coffee/core/appcontext"
	"coffee/core/rest"
	"coffee/helpers"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"go.step.sm/crypto/randutil"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type ProfileCollectionApiImpl struct {
	Service *Service
}

func NewProfileCollectionApi(service *Service) *ProfileCollectionApiImpl {
	apiImpl := ProfileCollectionApiImpl{
		Service: service,
	}
	return &apiImpl
}

func (api ProfileCollectionApiImpl) GetPrefix() string {
	return "/profile-collection-service"
}

func (api ProfileCollectionApiImpl) AttachRoutes(r *chi.Mux) {
	baseUrl := api.GetPrefix() + "/api/collection"
	// profile collection
	r.Post(baseUrl+"/", api.Create)
	r.Get(baseUrl+"/{id}", api.FindById)
	r.Put(baseUrl+"/{id}", api.Update)
	r.Delete(baseUrl+"/{id}", api.Delete)
	r.Post(baseUrl+"/search", api.Search)
	r.Get(baseUrl+"/byshareid/{shareId}", api.FindByShareId)
	r.Post(baseUrl+"/{id}/link/renew", api.GetNewShareId)
	r.Get(baseUrl+"/recent", api.FetchRecentCollectionsForPartner)

	// profile collection items
	r.Post(baseUrl+"/{collectionId}/item/", api.AddItemsToCollection)
	r.Delete(baseUrl+"/{collectionId}/item", api.DeleteItemsFromCollection)
	r.Put(baseUrl+"/{collectionId}/item/bulk", api.UpdateCollectionItemsInBulk)
	r.Put(baseUrl+"/byshareid/{shareId}/item/bulk", api.UpdateCollectionItemsInBulkByShareId)
	r.Post(baseUrl+"/item/search", api.SearchItemsInCollection)

	r.Get(baseUrl+"/winkl/{id}/migrationinfo", api.FetchWinklCollectionInfoById)
	r.Get(baseUrl+"/winklshare/{shareId}/migrationinfo", api.FetchWinklCollectionInfoByShareId)
}

// -------------------------------
// profile collection service
// -------------------------------

func (api ProfileCollectionApiImpl) Create(w http.ResponseWriter, r *http.Request) {
	var input domain.ProfileCollectionEntry
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		rest.RenderError(w, r, err)
		return
	}

	ctx := r.Context()
	appCtx := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)
	if appCtx.PartnerId == nil || appCtx.AccountId == nil {
		rest.RenderError(w, r, errors.New("you are not authorized to access this resource"))
		return
	}

	if input.Source == nil || input.Name == nil {
		rest.RenderError(w, r, errors.New("bad Request. Missing parameter source/name in input"))
		return
	}
	if strings.ToUpper(*input.Source) == string(constants.SaasCollection) || *input.Source == string(constants.SaasAtCollection) {
		if input.ShareId == nil {
			shareId, _ := randutil.Alphanumeric(7)
			input.ShareId = &shareId
		}
		if input.SourceId == nil {
			input.SourceId = input.ShareId
		}
	} else if (*input.Source == string(constants.CampaignCollection) || *input.Source == string(constants.GccCampaignCollection)) && input.SourceId == nil {
		rest.RenderError(w, r, errors.New("bad Request. Missing parameter sourceId in input"))
		return
	}

	input.PartnerId = appCtx.PartnerId
	accountIdStr := strconv.FormatInt(*appCtx.AccountId, 10)
	input.CreatedBy = &accountIdStr
	input.Enabled = helpers.ToBool(true)
	input.Featured = helpers.ToBool(false)
	input.AnalyticsEnabled = helpers.ToBool(false)
	if *input.Source == string(constants.SaasAtCollection) {
		input.AnalyticsEnabled = helpers.ToBool(true)
	}
	if input.ShareEnabled == nil {
		input.ShareEnabled = helpers.ToBool(true)
	}
	input.DisabledMetrics = []string{}
	if input.CustomColumns == nil {
		input.CustomColumns = []domain.PlatformCustomColumnMeta{}
	}
	if input.OrderedColumns == nil {
		input.OrderedColumns = []domain.PlatformOrderedColumnData{}
	}
	if input.CategoryIds == nil {
		input.CategoryIds = []int64{}
	}

	resp := api.Service.Create(r.Context(), &input)
	render.JSON(w, r, resp)
}

func (api ProfileCollectionApiImpl) FindById(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	resp := api.Service.FindById(r.Context(), id)
	render.JSON(w, r, resp)
}

func (api ProfileCollectionApiImpl) FindByShareId(w http.ResponseWriter, r *http.Request) {
	shareId := chi.URLParam(r, "shareId")
	resp := api.Service.FindByShareId(r.Context(), shareId)
	render.JSON(w, r, resp)
}

func (api ProfileCollectionApiImpl) GetNewShareId(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	appCtx := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)
	if appCtx.PartnerId == nil || appCtx.AccountId == nil {
		rest.RenderError(w, r, errors.New("you are not authorized to access this resource"))
		return
	}
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	resp := api.Service.RenewShareId(r.Context(), id)
	render.JSON(w, r, resp)
}

func (api ProfileCollectionApiImpl) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	appCtx := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)
	if appCtx.PartnerId == nil || appCtx.AccountId == nil {
		rest.RenderError(w, r, errors.New("you are not authorized to access this resource"))
		return
	}
	var input domain.ProfileCollectionEntry
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		rest.RenderError(w, r, err)
		return
	}
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	resp := api.Service.Update(r.Context(), id, &input)
	render.JSON(w, r, resp)
}

func (api ProfileCollectionApiImpl) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	appCtx := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)
	if appCtx.PartnerId == nil || appCtx.AccountId == nil {
		rest.RenderError(w, r, errors.New("you are not authorized to access this resource"))
		return
	}
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	resp := api.Service.Delete(r.Context(), id)
	render.JSON(w, r, resp)
}

func (api ProfileCollectionApiImpl) Search(w http.ResponseWriter, r *http.Request) {
	sortBy, sortDir, page, size, query, done := rest.ParseSearchParameters(w, r)
	if done {
		return
	}
	resp := api.Service.Search(r.Context(), query, sortBy, sortDir, page, size)
	render.JSON(w, r, resp)
}

func (api ProfileCollectionApiImpl) FetchRecentCollectionsForPartner(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	appCtx := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)
	if appCtx.PartnerId == nil || appCtx.AccountId == nil {
		rest.RenderError(w, r, errors.New("you are not authorized to access this resource"))
		return
	}
	resp := api.Service.FetchRecentCollectionsForPartner(r.Context())
	render.JSON(w, r, resp)
}

// -------------------------------
// profile collection item apis
// -------------------------------

func (api ProfileCollectionApiImpl) AddItemsToCollection(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	appCtx := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)
	if appCtx.PartnerId == nil || appCtx.AccountId == nil {
		rest.RenderError(w, r, errors.New("you are not authorized to access this resource"))
		return
	}
	var input []domain.ProfileCollectionItemEntry
	idStr := chi.URLParam(r, "collectionId")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		rest.RenderError(w, r, err)
		return
	}
	resp := api.Service.AddItemsToCollection(r.Context(), id, input)
	render.JSON(w, r, resp)
}

func (api ProfileCollectionApiImpl) DeleteItemsFromCollection(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	appCtx := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)
	if appCtx.PartnerId == nil || appCtx.AccountId == nil {
		rest.RenderError(w, r, errors.New("you are not authorized to access this resource"))
		return
	}
	var input []domain.ProfileCollectionItemEntry
	idStr := chi.URLParam(r, "collectionId")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		rest.RenderError(w, r, err)
		return
	}
	resp := api.Service.DeleteItemsFromCollection(r.Context(), id, input)
	render.JSON(w, r, resp)
}

func (api ProfileCollectionApiImpl) UpdateCollectionItemsInBulk(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	appCtx := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)
	if appCtx.PartnerId == nil || appCtx.AccountId == nil {
		rest.RenderError(w, r, errors.New("you are not authorized to access this resource"))
		return
	}
	var input []domain.ProfileCollectionItemEntry
	idStr := chi.URLParam(r, "collectionId")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		rest.RenderError(w, r, err)
		return
	}
	resp := api.Service.UpdateCollectionItemsInBulk(r.Context(), id, input)
	render.JSON(w, r, resp)
}

func (api ProfileCollectionApiImpl) UpdateCollectionItemsInBulkByShareId(w http.ResponseWriter, r *http.Request) {
	var input []domain.ProfileCollectionItemEntry
	shareId := chi.URLParam(r, "shareId")
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		rest.RenderError(w, r, err)
		return
	}
	resp := api.Service.UpdateCollectionItemsInBulkByShareId(r.Context(), shareId, input)
	render.JSON(w, r, resp)
}

func (api ProfileCollectionApiImpl) SearchItemsInCollection(w http.ResponseWriter, r *http.Request) {
	sortBy, sortDir, page, size, query, done := rest.ParseSearchParameters(w, r)
	if done {
		return
	}
	itemResponse := api.Service.SearchItemsInCollection(r.Context(), query, sortBy, sortDir, page, size)
	render.JSON(w, r, itemResponse)
}

//---------------------------------
// Winkl Collection Info apis
//---------------------------------

func (api ProfileCollectionApiImpl) FetchWinklCollectionInfoById(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	resp := api.Service.FetchWinklCollectionInfoById(r.Context(), id)
	render.JSON(w, r, resp)
}

func (api ProfileCollectionApiImpl) FetchWinklCollectionInfoByShareId(w http.ResponseWriter, r *http.Request) {
	shareId := chi.URLParam(r, "shareId")
	resp := api.Service.FetchWinklCollectionInfoByShareId(r.Context(), shareId)
	render.JSON(w, r, resp)
}

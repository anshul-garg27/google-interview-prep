package api

import (
	"coffee/app/keywordcollection/domain"
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

type KeywordCollectionApiImpl struct {
	Service *KeywordCollectionService
}

func NewKeywordCollectionApi(service *KeywordCollectionService) *KeywordCollectionApiImpl {
	apiImpl := KeywordCollectionApiImpl{
		Service: service,
	}
	return &apiImpl
}

func (s KeywordCollectionApiImpl) GetPrefix() string {
	return "/keyword-collection-service"
}

func (api KeywordCollectionApiImpl) AttachRoutes(r *chi.Mux) {
	baseUrl := api.GetPrefix() + "/api/collection"
	r.Post(baseUrl+"/", api.CreateCollection)
	r.Put(baseUrl+"/{id}", api.UpdateCollection)
	r.Get(baseUrl+"/{id}", api.FindCollectionById)
	r.Get(baseUrl+"/{id}/createReport", api.CreateReport)
	r.Get(baseUrl+"/byShareId/{shareId}", api.FindCollectionByShareId)
	r.Put(baseUrl+"/{id}/renewShareId", api.RenewCollectionShareId)
	r.Delete(baseUrl+"/{id}", api.DisableCollection)
	r.Post(baseUrl+"/search", api.SearchCollections)

	r.Post(baseUrl+"/report", api.GetCollectionReport)
	r.Post(baseUrl+"/posts", api.GetCollectionPosts)
	r.Post(baseUrl+"/profiles", api.GetCollectionProfiles)
}

func (api KeywordCollectionApiImpl) CreateCollection(w http.ResponseWriter, r *http.Request) {
	var input domain.KeywordCollectionEntry
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
	resp := api.Service.Create(r.Context(), &input)
	render.JSON(w, r, resp)
}

func (api KeywordCollectionApiImpl) CreateReport(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	appCtx := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)
	if appCtx.PartnerId == nil || appCtx.AccountId == nil {
		rest.RenderError(w, r, errors.New("You are not authorized to update this collection"))
		return
	}

	id := chi.URLParam(r, "id")
	resp := api.Service.CreateReport(r.Context(), id)
	render.JSON(w, r, resp)
}

func (api KeywordCollectionApiImpl) UpdateCollection(w http.ResponseWriter, r *http.Request) {
	var input domain.KeywordCollectionEntry
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

func (api KeywordCollectionApiImpl) FindCollectionById(w http.ResponseWriter, r *http.Request) {
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

func (api KeywordCollectionApiImpl) FindCollectionByShareId(w http.ResponseWriter, r *http.Request) {
	shareId := chi.URLParam(r, "shareId")
	resp := api.Service.FindByShareId(r.Context(), shareId)
	render.JSON(w, r, resp)
}

func (api KeywordCollectionApiImpl) RenewCollectionShareId(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	appCtx := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)
	if appCtx.PartnerId == nil || appCtx.AccountId == nil {
		rest.RenderError(w, r, errors.New("You are not authorized to update this collection"))
		return
	}

	id := chi.URLParam(r, "id")
	shareId, _ := randutil.Alphanumeric(7)
	input := domain.KeywordCollectionEntry{ShareId: &shareId}
	resp := api.Service.Service.Update(r.Context(), id, &input)
	render.JSON(w, r, resp)
}

func (api KeywordCollectionApiImpl) DisableCollection(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	appCtx := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)
	if appCtx.PartnerId == nil || appCtx.AccountId == nil {
		rest.RenderError(w, r, errors.New("You are not authorized to update this collection"))
		return
	}

	id := chi.URLParam(r, "id")
	input := domain.KeywordCollectionEntry{Enabled: helpers.ToBool(false)}
	resp := api.Service.Service.Update(r.Context(), id, &input)
	render.JSON(w, r, resp)
}

func (api KeywordCollectionApiImpl) SearchCollections(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	appCtx := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)
	if appCtx.PartnerId == nil || appCtx.AccountId == nil {
		rest.RenderError(w, r, errors.New("You are not authorized to perform this action"))
		return
	}
	sortBy, sortDir, page, size, query, done := rest.ParseSearchParameters(w, r)
	sortBy = "created_at"
	sortDir = "DESC"
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

func (api KeywordCollectionApiImpl) GetCollectionReport(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	appCtx := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)
	if appCtx.PartnerId == nil || appCtx.AccountId == nil {
		rest.RenderError(w, r, errors.New("You are not authorized to perform this action"))
		return
	}
	_, _, _, _, query, done := rest.ParseSearchParameters(w, r)
	if done {
		return
	}
	query.Filters = append(query.Filters, coredomain.SearchFilter{
		FilterType: "EQ",
		Field:      "partner_id",
		Value:      strconv.FormatInt(*appCtx.PartnerId, 10),
	})
	resp := api.Service.GetCollectionReport(r.Context(), query)
	render.JSON(w, r, resp)
}

func (api KeywordCollectionApiImpl) GetCollectionPosts(w http.ResponseWriter, r *http.Request) {
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
	resp := api.Service.GetCollectionPosts(r.Context(), query, sortBy, sortDir, page, size)
	render.JSON(w, r, resp)
}

func (api KeywordCollectionApiImpl) GetCollectionProfiles(w http.ResponseWriter, r *http.Request) {
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
	resp := api.Service.GetCollectionProfiles(r.Context(), query, sortBy, sortDir, page, size)
	render.JSON(w, r, resp)
}

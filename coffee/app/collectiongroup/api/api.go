package api

import (
	"coffee/app/collectiongroup/domain"
	"coffee/constants"
	"coffee/core/appcontext"
	"coffee/core/rest"
	"coffee/helpers"
	"encoding/json"
	"errors"
	"go.step.sm/crypto/randutil"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type CollectionGroupApiImpl struct {
	Service *Service
}

func NewCollectionGroupApi(service *Service) *CollectionGroupApiImpl {
	apiImpl := CollectionGroupApiImpl{
		Service: service,
	}
	return &apiImpl
}

func (api CollectionGroupApiImpl) GetPrefix() string {
	return "/collection-group-service"
}

func (api CollectionGroupApiImpl) AttachRoutes(r *chi.Mux) {
	baseUrl := api.GetPrefix() + "/api/group"
	// profile collection
	r.Post(baseUrl+"/", api.Create)
	r.Get(baseUrl+"/{id}", api.FindById)
	r.Put(baseUrl+"/{id}", api.Update)
	r.Get(baseUrl+"/byshareid/{shareId}", api.FindByShareId)
}

func (api CollectionGroupApiImpl) Create(w http.ResponseWriter, r *http.Request) {
	var input domain.CollectionGroupEntry
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

	if input.Source == nil || input.Name == nil || input.SourceId == nil {
		rest.RenderError(w, r, errors.New("Bad Request. Missing parameter source/name/sourceId in input"))
		return
	}
	shareId, _ := randutil.Alphanumeric(7)
	input.ShareId = &shareId
	input.PartnerId = appCtx.PartnerId
	accountIdStr := strconv.FormatInt(*appCtx.AccountId, 10)
	input.CreatedBy = &accountIdStr
	input.Enabled = helpers.ToBool(true)
	input.Metadata = &domain.CollectionGroupMetadata{
		ExpectedMetricValues: make(map[string]*float64),
		DisabledMetrics:      []string{},
		Comments:             nil,
	}
	resp := api.Service.Create(r.Context(), &input)
	render.JSON(w, r, resp)
}

func (api CollectionGroupApiImpl) FindById(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	resp := api.Service.FindById(r.Context(), id)
	render.JSON(w, r, resp)
}

func (api CollectionGroupApiImpl) FindByShareId(w http.ResponseWriter, r *http.Request) {
	shareId := chi.URLParam(r, "shareId")
	resp := api.Service.FindByShareId(r.Context(), shareId)
	render.JSON(w, r, resp)
}

func (api CollectionGroupApiImpl) GetNewShareId(w http.ResponseWriter, r *http.Request) {
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

func (api CollectionGroupApiImpl) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	appCtx := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)
	if appCtx.PartnerId == nil || appCtx.AccountId == nil {
		rest.RenderError(w, r, errors.New("you are not authorized to access this resource"))
		return
	}
	var input domain.CollectionGroupEntry
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		rest.RenderError(w, r, err)
		return
	}
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	resp := api.Service.Service.Update(r.Context(), id, &input)
	render.JSON(w, r, resp)
}

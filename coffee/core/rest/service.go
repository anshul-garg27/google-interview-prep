package rest

import (
	"coffee/constants"
	"coffee/core/appcontext"
	"coffee/core/domain"
	"coffee/core/persistence"
	"coffee/core/persistence/postgres"
	"context"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type ServiceWrapper interface {
}

type Service[RES domain.Response, EX domain.Entry, EN domain.Entity, I domain.ID] struct {
	ctx                 context.Context
	manager             *Manager[EX, EN, I]
	createResponse      func([]EX, int64, string, string) RES
	createErrorResponse func(error, int) RES
	getPrimaryStore     func() domain.DataStore
}

func NewService[RES domain.Response, EX domain.Entry, EN domain.Entity, I domain.ID](
	ctx context.Context,
	manager *Manager[EX, EN, I],
	getPrimaryStore func() domain.DataStore,
	createResponse func([]EX, int64, string, string) RES,
	createErrorResponse func(error, int) RES) *Service[RES, EX, EN, I] {
	manager.SetCtx(ctx)
	return &Service[RES, EX, EN, I]{
		ctx:                 ctx,
		manager:             manager,
		getPrimaryStore:     getPrimaryStore,
		createResponse:      createResponse,
		createErrorResponse: createErrorResponse,
	}
}

func (r *Service[RES, EX, EN, I]) Create(ctx context.Context, entry *EX) RES {
	createdEntry, err := r.manager.Create(ctx, entry)
	if err != nil {
		return r.createErrorResponse(err, 0)
	}
	return r.createResponse([]EX{*createdEntry}, 1, "", "Record Created Successfully")
}
func (r *Service[RES, EX, EN, I]) FindById(ctx context.Context, id I) RES {
	entry, err := r.manager.FindById(ctx, id)
	if err != nil {
		return r.createErrorResponse(err, 0)
	}
	return r.createResponse([]EX{*entry}, 1, "", "Record Retrieved Successfully")
}
func (r *Service[RES, EX, EN, I]) FindByIds(ctx context.Context, ids []I) RES {
	entries, err := r.manager.FindByIds(ctx, ids)
	if err != nil {
		return r.createErrorResponse(err, 0)
	}
	return r.createResponse(entries, int64(len(entries)), "", "Record(s) Retrieved Successfully")
}
func (r *Service[RES, EX, EN, I]) Update(ctx context.Context, id I, entry *EX) RES {
	_, err := r.manager.Update(ctx, id, entry)
	if err != nil {
		return r.createErrorResponse(err, 0)
	}
	return r.createResponse(nil, 1, "", "Record Updated Successfully")
}

func (r *Service[RES, EX, En, I]) Search(ctx context.Context, searchQuery domain.SearchQuery, sortBy string, sortDir string, page int, size int) RES {
	entries, filteredCount, err := r.manager.Search(ctx, searchQuery, sortBy, sortDir, page, size)
	if err != nil {
		return r.createErrorResponse(err, 0)
	}
	var nextCursor string
	if len(entries) >= size {
		nextCursor = strconv.Itoa(page + 1)
	} else {
		nextCursor = ""
	}
	return r.createResponse(entries, filteredCount, nextCursor, "Record(s) Retrieved Successfully")
}

func (r *Service[RES, EX, EN, I]) InitializeSession(ctx context.Context) persistence.Session {
	log.Debug("Initializing Session")
	if r.getPrimaryStore() == domain.Postgres {
		return postgres.InitializePgSession(ctx)
	}
	return nil
}

func (r *Service[RES, EX, EN, I]) Close(ctx context.Context) {
	session := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext).GetSession()
	if session != nil {
		session.Close(ctx)
	}
	//ctx.Value(constants.AppContextKey).(*appcontext.RequestContext).GetSession().Close(ctx)
}

func (r *Service[RES, EX, EN, I]) Rollback(ctx context.Context) {
	session := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext).GetSession()
	if session != nil {
		session.Rollback(ctx)
	}
}

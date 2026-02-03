package manager

import (
	"coffee/app/collectiongroup/dao"
	"coffee/app/collectiongroup/domain"
	"coffee/constants"
	"coffee/core/appcontext"
	"coffee/core/rest"
	"coffee/helpers"
	"context"
	"errors"
	"go.step.sm/crypto/randutil"
	"strings"
)

// -------------------------------
// profile collection service
// -------------------------------

type Manager struct {
	*rest.Manager[domain.CollectionGroupEntry, dao.CollectionGroupEntity, int64]
	dao *dao.Dao
}

func (m *Manager) init(ctx context.Context) {
	pcdao := dao.CreateDao(ctx)
	m.Manager = rest.NewManager(ctx, pcdao.Dao, ToEntry, ToEntity)
	m.dao = pcdao
}

func CreateManager(ctx context.Context) *Manager {
	manager := &Manager{}
	manager.init(ctx)
	return manager
}

func (m *Manager) Create(ctx context.Context, input *domain.CollectionGroupEntry) (*domain.CollectionGroupEntry, error) {
	createdEntry, err := m.Manager.Create(ctx, input)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") { // Cast to pgConn error not working for some reason, TODO: FIX
			appCtx := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)
			appCtx.Session.Close(ctx) // TODO: Figure out a better way to reset session. Need to do this because previous session has been marked as rollback-only
			appCtx.Session = nil
			entry, err := m.FindBySource(ctx, *input.Source, *input.SourceId)
			if err != nil {
				return nil, err
			}
			return entry, nil
		}
		return nil, err
	}
	if createdEntry == nil || createdEntry.Id == nil {
		return nil, errors.New("Failed to create profile collection")
	}
	return createdEntry, err
}

func (m *Manager) FindById(ctx context.Context, id int64) (*domain.CollectionGroupEntry, error) {
	entity, err := m.dao.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	appCtx := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)
	if appCtx.PartnerId == nil || appCtx.AccountId == nil {
		// logged out user trying to access non-featured collections
		return nil, errors.New("you are not authorized to access this collection group")
	}
	if entity.PartnerId != *appCtx.PartnerId && *appCtx.PartnerId != -1 {
		return nil, errors.New("you are not authorized to access this collection group")
	}

	entry, err := ToEntry(entity)

	if err != nil {
		return nil, err
	}

	return entry, nil
}

func (m *Manager) FindByIdInternal(ctx context.Context, id int64) (*domain.CollectionGroupEntry, error) {
	entity, err := m.dao.FindById(ctx, id)
	if err != nil {
		return nil, err
	}
	entry, err := ToEntry(entity)
	if err != nil {
		return nil, err
	}
	return entry, nil
}

func (m *Manager) FindBySource(ctx context.Context, source string, sourceId string) (*domain.CollectionGroupEntry, error) {
	entity, err := m.dao.FindBySource(ctx, source, sourceId)
	if err != nil {
		return nil, err
	}
	entry, err := ToEntry(entity)
	if err != nil {
		return nil, err
	}
	return entry, nil
}

func (m *Manager) FindByShareId(ctx context.Context, shareId string) (*domain.CollectionGroupEntry, error) {
	entity, err := m.dao.FindByShareId(ctx, shareId)
	if err != nil {
		return nil, err
	}
	entry, _ := ToEntry(entity)
	if err != nil {
		return nil, err
	}

	return entry, nil
}

func (m *Manager) RenewShareId(ctx context.Context, id int64) (*domain.CollectionGroupEntry, error) {
	shareId, _ := randutil.Alphanumeric(7)
	entry := domain.CollectionGroupEntry{
		ShareId: &shareId,
	}

	updatedEntry, err := m.Update(ctx, id, &entry)
	if err != nil {
		return nil, err
	}

	return updatedEntry, nil
}

func (m *Manager) Delete(ctx context.Context, id int64) (*domain.CollectionGroupEntry, error) {
	entry := domain.CollectionGroupEntry{
		Enabled: helpers.ToBool(false),
	}
	updatedEntry, err := m.Update(ctx, id, &entry)
	if err != nil {
		return nil, err
	}
	return updatedEntry, nil
}

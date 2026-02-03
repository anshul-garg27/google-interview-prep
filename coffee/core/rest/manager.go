package rest

import (
	"coffee/core/domain"
	"context"
	"errors"
)

type Manager[EX domain.Entry, EN domain.Entity, I domain.ID] struct {
	dao      *Dao[EN, I]
	toEntity func(*EX) (*EN, error)
	toEntry  func(*EN) (*EX, error)
}

func NewManager[EX domain.Entry, EN domain.Entity, I domain.ID](ctx context.Context,
	dao *Dao[EN, I], toEntry func(*EN) (*EX, error), toEntity func(*EX) (*EN, error)) *Manager[EX, EN, I] {
	return &Manager[EX, EN, I]{
		dao:      dao,
		toEntity: toEntity,
		toEntry:  toEntry,
	}
}

func (r *Manager[EX, EN, I]) SetCtx(ctx context.Context) {
	r.dao.SetCtx(ctx)
}

func (r *Manager[EX, EN, I]) Create(ctx context.Context, entry *EX) (*EX, error) {
	if entry == nil {
		return nil, errors.New("domain is null")
	}
	entity, err := r.toEntity(entry)
	if err != nil {
		return nil, err
	}
	createdEntity, err := r.dao.Create(ctx, entity)
	if err != nil {
		return nil, err
	}
	return r.toEntry(createdEntity)
}

func (r *Manager[EX, EN, I]) FindById(ctx context.Context, id I) (*EX, error) {
	entity, err := r.dao.FindById(ctx, id)
	if err != nil {
		return nil, err
	}
	// if entity == nil {
	// 	return nil, errors.New(fmt.Sprintf("entity not found for id - %v", id))
	// }
	return r.toEntry(entity)
}

func (r *Manager[EX, EN, I]) FindByIds(ctx context.Context, ids []I) ([]EX, error) {
	entities, err := r.dao.FindByIds(ctx, ids)
	if err != nil {
		return nil, err
	}
	entries := []EX{}
	for _, entity := range entities {
		entry, _ := r.toEntry(&entity)
		entries = append(entries, *entry)
	}
	return entries, nil
}

func (r *Manager[EX, EN, I]) Update(ctx context.Context, id I, entry *EX) (*EX, error) {
	entity, err := r.toEntity(entry)
	if err != nil {
		return nil, err
	}
	updatedEntity, err := r.dao.Update(ctx, id, entity)
	if err != nil {
		return nil, err
	}
	return r.toEntry(updatedEntity)
}

func (m *Manager[EX, EN, I]) Search(
	ctx context.Context,
	query domain.SearchQuery,
	sortBy string,
	sortDir string,
	page int,
	size int) ([]EX, int64, error) {
	entities, filteredCount, err := m.dao.Search(ctx, query, sortBy, sortDir, page, size)
	if err != nil {
		return nil, 0, err
	}
	var entries []EX
	for i := range entities {
		entry, err := m.toEntry(&entities[i])
		if err == nil && entry != nil {
			if entry != nil {
				entries = append(entries, *entry)
			}
		}
	}
	return entries, filteredCount, err
}

package rest

import (
	"coffee/core/domain"
	"context"
)

type Dao[EN domain.Entity, I domain.ID] struct {
	DaoProvider[EN, I]
	ctx context.Context
}

func NewDao[EN domain.Entity, I domain.ID](ctx context.Context) *Dao[EN, I] {
	return &Dao[EN, I]{
		ctx: ctx,
	}
}

func (r *Dao[EN, I]) SetCtx(ctx context.Context) {
	r.ctx = ctx
}

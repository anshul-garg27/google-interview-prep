package dao

import (
	"coffee/core/persistence/postgres"
	"coffee/core/rest"
	"context"
)

type TrendingContentDao struct {
	*rest.Dao[TrendingContentEntity, int64]
	ctx context.Context
}

func CreateTrendingContentDao(ctx context.Context) *TrendingContentDao {
	dao := &TrendingContentDao{
		ctx: ctx,
		Dao: rest.NewDao[TrendingContentEntity, int64](ctx),
	}
	dao.DaoProvider = postgres.NewPgDao[TrendingContentEntity, int64](ctx)
	return dao
}

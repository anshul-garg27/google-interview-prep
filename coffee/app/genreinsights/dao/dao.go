package dao

import (
	"coffee/core/persistence/postgres"
	"coffee/core/rest"
	"context"
)

type GenreOverviewDao struct {
	*rest.Dao[GenreOverviewEntity, int64]
	ctx context.Context
}

func CreateGenreInsightsDao(ctx context.Context) *GenreOverviewDao {
	dao := &GenreOverviewDao{
		ctx: ctx,
		Dao: rest.NewDao[GenreOverviewEntity, int64](ctx),
	}
	dao.DaoProvider = postgres.NewPgDao[GenreOverviewEntity, int64](ctx)
	return dao
}

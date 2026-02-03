package rest

import (
	"coffee/core/domain"
	"context"

	"gorm.io/gorm"
)

type DaoProvider[EN domain.Entity, I domain.ID] interface {
	FindById(ctx context.Context, id I) (*EN, error)
	FindByIds(ctx context.Context, ids []I) ([]EN, error)
	Create(ctx context.Context, entity *EN) (*EN, error)
	Update(ctx context.Context, id I, entity *EN) (*EN, error)
	Search(ctx context.Context, query domain.SearchQuery, sortBy string, sortDir string, page int, size int) ([]EN, int64, error)
	SearchJoins(ctx context.Context, query domain.SearchQuery, sortBy string, sortDir string, page int, size int, joinTables []domain.JoinClauses) ([]EN, int64, error)
	GetSession(ctx context.Context) *gorm.DB
	AddPredicateForSearchJoins(req *gorm.DB, filter domain.SearchFilter) *gorm.DB
}

package dao

import (
	"coffee/core/persistence/postgres"
	"coffee/core/rest"
	"context"
)

// -------------------------------
// collection post metrics summary dao
// -------------------------------

type CollectionKeywordDAO struct {
	*rest.Dao[CollectionKeywordEntity, int64]
	ctx context.Context
}

func (d *CollectionKeywordDAO) init(ctx context.Context) {
	d.Dao = rest.NewDao[CollectionKeywordEntity, int64](ctx)
	d.Dao.DaoProvider = postgres.NewPgDao[CollectionKeywordEntity, int64](ctx)
	d.ctx = ctx
}

func CreateCollectionKeywordDAO(ctx context.Context) *CollectionKeywordDAO {
	dao := &CollectionKeywordDAO{}
	dao.init(ctx)
	return dao
}

func (d *CollectionKeywordDAO) FindByCollectionIdAndType(ctx context.Context, collectionIds []string, collectionType string) ([]CollectionKeywordEntity, error) {
	dbSession := d.DaoProvider.GetSession(ctx)
	var entities []CollectionKeywordEntity
	tx := dbSession.WithContext(ctx).Model(&CollectionKeywordEntity{}).Where("collection_id IN ? AND collection_type = ?", collectionIds, collectionType).Scan(&entities)
	if tx.Error != nil && tx.Error.Error() == "record not found" {
		return nil, nil
	}
	return entities, tx.Error
}

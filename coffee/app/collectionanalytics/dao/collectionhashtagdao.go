package dao

import (
	"coffee/core/persistence/postgres"
	"coffee/core/rest"
	"context"
)

// -------------------------------
// collection post metrics summary dao
// -------------------------------

type CollectionHashTagDAO struct {
	*rest.Dao[CollectionHashtagEntity, int64]
	ctx context.Context
}

func (d *CollectionHashTagDAO) init(ctx context.Context) {
	d.Dao = rest.NewDao[CollectionHashtagEntity, int64](ctx)
	d.Dao.DaoProvider = postgres.NewPgDao[CollectionHashtagEntity, int64](ctx)
	d.ctx = ctx
}

func CreateCollectionHashTagDAO(ctx context.Context) *CollectionHashTagDAO {
	dao := &CollectionHashTagDAO{}
	dao.init(ctx)
	return dao
}

func (d *CollectionHashTagDAO) FindByCollectionIdAndType(ctx context.Context, collectionIds []string, collectionType string) ([]CollectionHashtagEntity, error) {
	dbSession := d.DaoProvider.GetSession(ctx)
	var entities []CollectionHashtagEntity
	tx := dbSession.WithContext(ctx).Model(&CollectionHashtagEntity{}).Where("collection_id IN ? AND collection_type = ?", collectionIds, collectionType).Scan(&entities)
	if tx.Error != nil && tx.Error.Error() == "record not found" {
		return nil, nil
	}
	return entities, tx.Error
}

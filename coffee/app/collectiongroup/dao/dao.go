package dao

import (
	"coffee/core/persistence/postgres"
	"coffee/core/rest"
	"context"
)

type Dao struct {
	*rest.Dao[CollectionGroupEntity, int64]
	ctx context.Context
}

func (d *Dao) init(ctx context.Context) {
	d.Dao = rest.NewDao[CollectionGroupEntity, int64](ctx)
	d.Dao.DaoProvider = postgres.NewPgDao[CollectionGroupEntity, int64](ctx)
	d.ctx = ctx
}

func CreateDao(ctx context.Context) *Dao {
	dao := &Dao{}
	dao.init(ctx)
	return dao
}

func (d *Dao) CountCollectionsForPartner(ctx context.Context, partnerId int64) (int64, error) {
	dbSession := d.DaoProvider.GetSession(ctx)
	count := int64(0)
	tx := dbSession.WithContext(ctx).Model(&CollectionGroupEntity{}).Where(CollectionGroupEntity{PartnerId: partnerId}).Count(&count)
	if tx.Error != nil {
		return count, tx.Error
	}
	return count, nil
}

func (d *Dao) FindByShareId(ctx context.Context, shareId string) (*CollectionGroupEntity, error) {
	dbSession := d.DaoProvider.GetSession(ctx)
	var entity CollectionGroupEntity
	tx := dbSession.WithContext(ctx).Model(&entity).Where(CollectionGroupEntity{ShareId: shareId}).First(&entity)
	if tx.Error != nil && tx.Error.Error() == "record not found" {
		return nil, tx.Error
	}
	return &entity, tx.Error
}

func (d *Dao) FindBySource(ctx context.Context, source string, sourceId string) (*CollectionGroupEntity, error) {
	dbSession := d.DaoProvider.GetSession(ctx)
	var entity CollectionGroupEntity
	tx := dbSession.WithContext(ctx).Model(&entity).Where(CollectionGroupEntity{Source: source, SourceId: sourceId}).First(&entity)
	if tx.Error != nil && tx.Error.Error() == "record not found" {
		return nil, tx.Error
	}
	return &entity, tx.Error
}

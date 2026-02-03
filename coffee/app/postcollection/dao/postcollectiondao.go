package dao

import (
	"coffee/core/persistence/postgres"
	"coffee/core/rest"
	"context"
)

// -------------------------------
// profile collection service
// -------------------------------

type PostCollectionDAO struct {
	*rest.Dao[PostCollectionEntity, string]
	ctx context.Context
}

func (d *PostCollectionDAO) init(ctx context.Context) {
	d.Dao = rest.NewDao[PostCollectionEntity, string](ctx)
	d.Dao.DaoProvider = postgres.NewPgDao[PostCollectionEntity, string](ctx)
	d.ctx = ctx
}

func CreatePostCollectionDAO(ctx context.Context) *PostCollectionDAO {
	dao := &PostCollectionDAO{}
	dao.init(ctx)
	return dao
}

func (d *PostCollectionDAO) FindByShareId(ctx context.Context, shareId string) (*PostCollectionEntity, error) {
	dbSession := d.DaoProvider.GetSession(ctx)
	var entity PostCollectionEntity
	tx := dbSession.WithContext(ctx).Model(&entity).Where(PostCollectionEntity{ShareId: shareId}).First(&entity)
	if tx.Error != nil && tx.Error.Error() == "record not found" {
		return nil, tx.Error
	}
	return &entity, tx.Error
}

func (d *PostCollectionDAO) FindBySourceAndSourceId(ctx context.Context, source string, sourceId string) (*PostCollectionEntity, error) {
	dbSession := d.DaoProvider.GetSession(ctx)
	var entity PostCollectionEntity
	tx := dbSession.WithContext(ctx).Model(&entity).Where(PostCollectionEntity{Source: source, SourceId: sourceId}).First(&entity)
	if tx.Error != nil && tx.Error.Error() == "record not found" {
		return nil, tx.Error
	}
	return &entity, tx.Error
}

func (d *PostCollectionDAO) CountCollectionsForPartner(ctx context.Context, partnerId int64) (int64, error) {
	dbSession := d.DaoProvider.GetSession(ctx)
	count := int64(0)
	tx := dbSession.WithContext(ctx).Model(&PostCollectionEntity{}).Where(PostCollectionEntity{PartnerId: partnerId}).Count(&count)
	if tx.Error != nil {
		return count, tx.Error
	}
	return count, nil
}

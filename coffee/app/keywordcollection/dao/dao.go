package dao

import (
	"coffee/core/persistence/postgres"
	"coffee/core/rest"
	"context"
)

type KeywordCollectionDAO struct {
	*rest.Dao[KeywordCollectionEntity, string]
	ctx context.Context
}

func (d *KeywordCollectionDAO) init(ctx context.Context) {
	d.Dao = rest.NewDao[KeywordCollectionEntity, string](ctx)
	d.Dao.DaoProvider = postgres.NewPgDao[KeywordCollectionEntity, string](ctx)
	d.ctx = ctx
}

func CreateKeywordCollectionDAO(ctx context.Context) *KeywordCollectionDAO {
	dao := &KeywordCollectionDAO{}
	dao.init(ctx)
	return dao
}

func (d *KeywordCollectionDAO) FindByShareId(ctx context.Context, shareId string) (*KeywordCollectionEntity, error) {
	dbSession := d.DaoProvider.GetSession(ctx)
	var entity KeywordCollectionEntity
	tx := dbSession.WithContext(ctx).Model(&entity).Where(KeywordCollectionEntity{ShareId: shareId}).First(&entity)
	if tx.Error != nil && tx.Error.Error() == "record not found" {
		return nil, tx.Error
	}
	return &entity, tx.Error
}

func (d *KeywordCollectionDAO) CountCollectionsForPartner(ctx context.Context, partnerId int64) (int64, error) {
	dbSession := d.DaoProvider.GetSession(ctx)
	count := int64(0)
	tx := dbSession.WithContext(ctx).Model(&KeywordCollectionEntity{}).Where(KeywordCollectionEntity{PartnerId: partnerId}).Count(&count)
	if tx.Error != nil {
		return count, tx.Error
	}
	return count, nil
}

package dao

import (
	"coffee/core/persistence/postgres"
	"coffee/core/rest"
	"context"
	"time"
)

type PostCollectionItemDAO struct {
	*rest.Dao[PostCollectionItemEntity, int64]
	ctx context.Context
}

func (d *PostCollectionItemDAO) init(ctx context.Context) {
	d.Dao = rest.NewDao[PostCollectionItemEntity, int64](ctx)
	pgDaoProvider := postgres.NewPgDao[PostCollectionItemEntity, int64](ctx)
	d.Dao.DaoProvider = pgDaoProvider
	d.ctx = ctx
}

func CreatePostCollectionItemDAO(ctx context.Context) *PostCollectionItemDAO {
	dao := &PostCollectionItemDAO{}
	dao.init(ctx)
	return dao
}

func (d *PostCollectionItemDAO) FindCollectionItemByPlatformAndShortCode(ctx context.Context, collectionId string, platform string, shortCode string) (*PostCollectionItemEntity, error) {
	dbSession := d.DaoProvider.GetSession(ctx)
	var entity *PostCollectionItemEntity
	tx := dbSession.WithContext(ctx).Model(&entity).Where(PostCollectionItemEntity{PostCollectionId: collectionId, Platform: platform, ShortCode: shortCode}).First(&entity)
	if tx.Error != nil && tx.Error.Error() == "record not found" {
		return nil, nil
	}
	return entity, tx.Error
}

func (d *PostCollectionItemDAO) CountCollectionItemsInACollection(ctx context.Context, collectionId string) (int64, error) {
	dbSession := d.DaoProvider.GetSession(ctx)
	var count int64
	tx := dbSession.WithContext(ctx).Model(&PostCollectionItemEntity{}).Where(map[string]interface{}{"post_collection_id": collectionId, "enabled": true}).Count(&count)
	if tx.Error != nil {
		return 0, tx.Error
	}
	return count, nil
}

func (d *PostCollectionItemDAO) GetLatestItemAddedOn(ctx context.Context, collectionId string) (*time.Time, error) {
	dbSession := d.DaoProvider.GetSession(ctx)
	var timestamp *time.Time
	tx := dbSession.WithContext(ctx).Model(&PostCollectionItemEntity{}).Where(map[string]interface{}{"post_collection_id": collectionId, "enabled": true}).Select("max(created_at)").Find(timestamp)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return timestamp, nil
}

func (d *PostCollectionItemDAO) CountCollectionItemsInCollections(ctx context.Context, collectionIds []string) (map[string]int64, error) {
	dbSession := d.DaoProvider.GetSession(ctx)
	var counts map[string]int64
	tx := dbSession.WithContext(ctx).Model(&PostCollectionItemEntity{}).Select("post_collection_id as cid, count(*) as items").Where(map[string]interface{}{"post_collection_id": collectionIds, "enabled": true}).Find(&counts)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return counts, nil
}

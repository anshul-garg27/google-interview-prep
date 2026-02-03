package dao

import (
	"coffee/core/persistence/postgres"
	"coffee/core/rest"
	"context"
)

type PartnerUsageDao struct {
	*rest.Dao[PartnerUsageEntity, int64]
	ctx context.Context
}

func CreatePartnerUsageDao(ctx context.Context) *PartnerUsageDao {
	dao := &PartnerUsageDao{
		Dao: rest.NewDao[PartnerUsageEntity, int64](ctx),
		ctx: ctx,
	}
	dao.DaoProvider = postgres.NewPgDao[PartnerUsageEntity, int64](ctx)
	return dao
}

type PartnerProfileTrackDao struct {
	*rest.Dao[PartnerProfileTrackEntity, int64]
	ctx context.Context
}

func CreatePartnerProfileTrackDao(ctx context.Context) *PartnerProfileTrackDao {
	dao := &PartnerProfileTrackDao{
		Dao: rest.NewDao[PartnerProfileTrackEntity, int64](ctx),
		ctx: ctx,
	}
	dao.DaoProvider = postgres.NewPgDao[PartnerProfileTrackEntity, int64](ctx)
	return dao
}

type ActivityTrackerDao struct {
	*rest.Dao[AcitivityTrackerEntity, int64]
	ctx context.Context
}

func CreateActivityTrackerDao(ctx context.Context) *ActivityTrackerDao {
	dao := &ActivityTrackerDao{
		Dao: rest.NewDao[AcitivityTrackerEntity, int64](ctx),
		ctx: ctx,
	}
	dao.DaoProvider = postgres.NewPgDao[AcitivityTrackerEntity, int64](ctx)
	return dao
}

func (P *PartnerUsageDao) FindConsuptionByPartnerId(ctx context.Context, partnerId int64) (*[]PartnerUsageEntity, error) {

	dbSession := P.GetSession(ctx)
	var entities []PartnerUsageEntity
	var model PartnerUsageEntity
	tx := dbSession.WithContext(ctx).Model(&model).Where("partner_id = ? and key = ?", partnerId, "USAGE").Find(&entities)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &entities, tx.Error
}
func (P *PartnerUsageDao) CheckUsage(ctx context.Context, partnerId int64, moduleName string) (*PartnerUsageEntity, error) {
	dbSession := P.GetSession(ctx)
	var entity PartnerUsageEntity
	tx := dbSession.WithContext(ctx).Model(&entity).Where(PartnerUsageEntity{PartnerId: partnerId, Key: "USAGE", Namespace: moduleName}).First(&entity)

	if tx.Error != nil && tx.Error.Error() == "record not found" {
		return nil, tx.Error
	}
	return &entity, tx.Error
}

func (PT *PartnerProfileTrackDao) FindByProfileKey(ctx context.Context, partnerId int64, key string) (*PartnerProfileTrackEntity, error) {
	dbSession := PT.GetSession(ctx)
	var entity PartnerProfileTrackEntity
	tx := dbSession.WithContext(ctx).Model(&entity).Where(PartnerProfileTrackEntity{PartnerId: partnerId, Key: key}).First(&entity)

	if tx.Error != nil && tx.Error.Error() == "record not found" {
		return nil, tx.Error
	}
	return &entity, tx.Error
}
func (P *PartnerUsageDao) DeletePartnerUsage(ctx context.Context, partnerId int64) error {
	dbSession := P.GetSession(ctx)
	var entity PartnerUsageEntity
	deleteResult := dbSession.WithContext(ctx).Where("partner_id = ?", partnerId).Delete(&entity)

	if deleteResult.Error != nil {
		return deleteResult.Error
	}
	return nil
}

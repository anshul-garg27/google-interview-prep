package dao

import (
	"coffee/core/persistence/postgres"
	"coffee/core/rest"
	"context"
	"strconv"
)

type CampaignProfileDao struct {
	*rest.Dao[CampaignProfileEntity, int64]
	ctx context.Context
}

func CreateCampaignProfileDao(ctx context.Context) *CampaignProfileDao {
	dao := &CampaignProfileDao{
		ctx: ctx,
		Dao: rest.NewDao[CampaignProfileEntity, int64](ctx),
	}
	dao.DaoProvider = postgres.NewPgDao[CampaignProfileEntity, int64](ctx)
	return dao
}

func (d *CampaignProfileDao) FindByPlatformAndPlatformId(ctx context.Context, platform string, platformCode string) (*CampaignProfileEntity, error) {
	dbSession := d.DaoProvider.GetSession(ctx)
	var entity CampaignProfileEntity
	platformAccountId, _ := strconv.ParseInt(platformCode, 10, 64)
	tx := dbSession.WithContext(ctx).Model(&entity).Where(CampaignProfileEntity{Platform: platform, PlatformAccountId: platformAccountId}).First(&entity)
	if tx.Error != nil && tx.Error.Error() == "record not found" {
		return nil, nil
	}
	return &entity, tx.Error
}

func (d *CampaignProfileDao) Update(ctx context.Context, entity CampaignProfileEntity, id int64) (*CampaignProfileEntity, error) {
	dbSession := d.DaoProvider.GetSession(ctx)
	tx := dbSession.WithContext(ctx).Model(&entity).Select("*").Where(CampaignProfileEntity{ID: id}).Updates(entity)
	if tx.Error != nil && tx.Error.Error() == "record not found" {
		return nil, tx.Error
	}
	return &entity, nil
}

func (d *CampaignProfileDao) FindById(ctx context.Context, id int64) (*CampaignProfileEntity, error) {
	campaignEntity, err := d.DaoProvider.FindById(ctx, id)
	if err != nil {
		return nil, err
	}
	return campaignEntity, err
}

func (d *CampaignProfileDao) FindCPByAccountId(ctx context.Context, platform string, accountId *int64) (*CampaignProfileEntity, error) {
	dbSession := d.DaoProvider.GetSession(ctx)
	var entity CampaignProfileEntity
	tx := dbSession.WithContext(ctx).Model(&entity).Where(CampaignProfileEntity{Platform: platform, GCCUserAccountID: accountId}).First(&entity)
	if tx.Error != nil && tx.Error.Error() == "record not found" {
		return nil, tx.Error
	}
	return &entity, tx.Error
}

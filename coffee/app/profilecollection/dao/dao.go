package dao

import (
	"coffee/app/profilecollection/domain"
	"coffee/constants"
	"coffee/core/persistence/postgres"
	"coffee/core/rest"
	"coffee/helpers"
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"gorm.io/gorm/clause"
)

// -------------------------------
// profile collection service
// -------------------------------

type Dao struct {
	*rest.Dao[ProfileCollectionEntity, int64]
	ctx context.Context
}

func (d *Dao) init(ctx context.Context) {
	d.Dao = rest.NewDao[ProfileCollectionEntity, int64](ctx)
	d.Dao.DaoProvider = postgres.NewPgDao[ProfileCollectionEntity, int64](ctx)
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
	tx := dbSession.WithContext(ctx).Model(&ProfileCollectionEntity{}).Where(ProfileCollectionEntity{PartnerId: partnerId}).Count(&count)
	if tx.Error != nil {
		return count, tx.Error
	}
	return count, nil
}

func (d *Dao) FindByShareId(ctx context.Context, shareId string) (*ProfileCollectionEntity, error) {
	dbSession := d.DaoProvider.GetSession(ctx)
	var entity ProfileCollectionEntity
	tx := dbSession.WithContext(ctx).Model(&entity).Where(ProfileCollectionEntity{ShareId: shareId}).First(&entity)
	if tx.Error != nil && tx.Error.Error() == "record not found" {
		return nil, tx.Error
	}
	return &entity, tx.Error
}

func (d *Dao) FindBySourceAndSourceId(ctx context.Context, source string, sourceId string) (*ProfileCollectionEntity, error) {
	dbSession := d.DaoProvider.GetSession(ctx)
	var entity ProfileCollectionEntity
	tx := dbSession.WithContext(ctx).Model(&entity).Where(ProfileCollectionEntity{Source: source, SourceId: sourceId}).First(&entity)
	if tx.Error != nil && tx.Error.Error() == "record not found" {
		return nil, tx.Error
	}
	return &entity, tx.Error
}

func (d *Dao) FindCollectionByCampaignId(ctx context.Context, campaignId string) (*ProfileCollectionEntity, error) {
	dbSession := d.DaoProvider.GetSession(ctx)
	var entity ProfileCollectionEntity
	cId, _ := strconv.ParseInt(campaignId, 10, 64)
	tx := dbSession.WithContext(ctx).Model(&entity).Where(ProfileCollectionEntity{CampaignId: &cId}).First(&entity)
	if tx.Error != nil && tx.Error.Error() == "record not found" {
		return nil, tx.Error
	}
	return &entity, tx.Error
}

func (d *Dao) FindCollectionByShortlistId(ctx context.Context, shortlistId string) (*ProfileCollectionEntity, error) {
	dbSession := d.DaoProvider.GetSession(ctx)
	var entity ProfileCollectionEntity
	tx := dbSession.WithContext(ctx).Model(&entity).Joins("join profile_collection_item on profile_collection_item.profile_collection_id = profile_collection.id").Where("profile_collection_item.shortlist_id = ?", shortlistId).First(&entity)
	if tx.Error != nil && tx.Error.Error() == "record not found" {
		return nil, nil
	}
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &entity, tx.Error
}

func (d *Dao) FindMostRecentAddedCollection(ctx context.Context, partnerId int64) (*ProfileCollectionEntity, error) {
	dbSession := d.DaoProvider.GetSession(ctx)
	var entity ProfileCollectionEntity
	tx := dbSession.WithContext(ctx).Model(&entity).Where(ProfileCollectionEntity{PartnerId: partnerId, Enabled: helpers.ToBool(true), Source: string(constants.SaasCollection)}).Order(fmt.Sprintf("%s %s", "created_at", "desc")).First(&entity)
	if tx.Error != nil && tx.Error.Error() == "record not found" {
		return nil, nil
	}
	return &entity, tx.Error
}

// -------------------------------
// profile collection item dao
// -------------------------------

type ItemDao struct {
	*rest.Dao[ProfileCollectionItemEntity, int64]
	ctx context.Context
}

func (d *ItemDao) itemInit(ctx context.Context) {
	d.Dao = rest.NewDao[ProfileCollectionItemEntity, int64](ctx)
	pgDaoProvider := postgres.NewPgDao[ProfileCollectionItemEntity, int64](ctx)
	d.Dao.DaoProvider = pgDaoProvider
	d.ctx = ctx
}

func CreateItemDao(ctx context.Context) *ItemDao {
	itemDao := &ItemDao{}
	itemDao.itemInit(ctx)
	return itemDao
}

func (d *ItemDao) CountCollectionItems(ctx context.Context, collectionId int64) []CollectionItemCountsEntity {
	dbSession := d.DaoProvider.GetSession(ctx)
	var itemCountsRows []CollectionItemCountsEntity
	tx := dbSession.WithContext(ctx).Table("profile_collection_item").Select("profile_collection_id as collectionId, profile_collection_item.platform as platform, count(profile_collection_item.id) as allItems, count(campaign_profiles.id) FILTER (WHERE campaign_profiles.has_email = true) AS emailItems, count(campaign_profiles.id) FILTER (WHERE campaign_profiles.has_phone = true) AS phoneItems").Joins("JOIN campaign_profiles on campaign_profiles.platform_account_id = profile_collection_item.platform_account_code and campaign_profiles.platform = profile_collection_item.platform").Where(map[string]interface{}{"profile_collection_item.profile_collection_id": collectionId, "profile_collection_item.enabled": "true"}).Group("profile_collection_item.profile_collection_id, profile_collection_item.platform").Find(&itemCountsRows)
	if tx.Error != nil && tx.Error.Error() == "record not found" {
		return itemCountsRows
	}
	return itemCountsRows
}

/*
func (d *ItemDao) FetchCollectionItemsCountSummary(ctx context.Context, collectionId int64) {
	dbSession := d.DaoProvider.GetSession(ctx)

	profileCountsQuery := "select platform, count(distinct platform_account_code) as count from profile_collection_item where profile_collection_id = @collectionId  group by platform"
	dbSession.Raw(profileCountsQuery, map[string]interface{}{"collectionId": collectionId}).Find(&profileCounters)

	return profileCounters
}
*/

func (d *ItemDao) FindCollectionItem(ctx context.Context, collectionId int64, platform string, platformAccountCode int64) *ProfileCollectionItemEntity {
	dbSession := d.DaoProvider.GetSession(ctx)
	var entity ProfileCollectionItemEntity
	tx := dbSession.WithContext(ctx).Model(&entity).Where(ProfileCollectionItemEntity{ProfileCollectionId: collectionId, Platform: platform, PlatformAccountCode: platformAccountCode}).First(&entity)
	if tx.Error != nil && tx.Error.Error() == "record not found" {
		return nil
	}
	return &entity
}

func (d *ItemDao) FindMaxItemRankForCollectionByPlatform(ctx context.Context, collectionId int64) (int64, int64) {
	dbSession := d.DaoProvider.GetSession(ctx)
	type Counts struct {
		Platform string
		Maxrank  int64
	}
	var entity ProfileCollectionItemEntity
	counters := []Counts{}
	tx := dbSession.WithContext(ctx).Model(&entity).Select("platform, max(rank) maxrank").Where(ProfileCollectionItemEntity{ProfileCollectionId: collectionId}).Group("platform").Scan(&counters)

	rankInstagram := int64(1)
	rankYoutube := int64(1)
	if tx.Error == nil && len(counters) > 0 {
		for _, row := range counters {
			if row.Platform == string(constants.InstagramPlatform) {
				rankInstagram = row.Maxrank
			} else if row.Platform == string(constants.YoutubePlatform) {
				rankYoutube = row.Maxrank
			}
		}
	}
	return rankInstagram, rankYoutube
}

func (d *ItemDao) FindMostRecentCollectionByItemUpdate(ctx context.Context, partnerId int64) []ProfileCollectionItemEntity {
	dbSession := d.DaoProvider.GetSession(ctx)
	// fetch last 2 recently edited collections
	var entities []ProfileCollectionItemEntity
	tx := dbSession.WithContext(ctx).Model(&ProfileCollectionItemEntity{}).Joins("join profile_collection on profile_collection.id = profile_collection_item.profile_collection_id").Where("profile_collection.partner_id = ?", partnerId).Where("profile_collection.enabled = ?", "true").Order(fmt.Sprintf("%s %s", "profile_collection_item.updated_at", "desc")).Limit(2).Find(&entities)
	if tx.Error != nil && tx.Error.Error() == "record not found" {
		return nil
	}
	return entities
}

func (d *ItemDao) UpdateByCollectionItemCampaignProfileId(ctx context.Context, collectionId int64, campaignProfileId int64, entity ProfileCollectionItemEntity) (*ProfileCollectionItemEntity, error) {
	db := d.DaoProvider.GetSession(ctx)
	// fetch last 2 recently edited collections
	db = db.Model(&entity).Clauses(clause.Returning{}).Where("campaign_profile_id = ?", campaignProfileId).Where("profile_collection_id = ?", collectionId).Updates(&entity)
	if db.Error != nil {
		return nil, db.Error
	} else if db.RowsAffected == 0 {
		return nil, fmt.Errorf("invalid Campaign Profile Id - %v", campaignProfileId)
	}
	return &entity, nil
}

func (d *ItemDao) UpdateCollectionItemByPlatformAccountCode(ctx context.Context, collectionId int64, platformAccountCode int64, entity ProfileCollectionItemEntity) (*ProfileCollectionItemEntity, error) {
	db := d.DaoProvider.GetSession(ctx)
	// fetch last 2 recently edited collections
	db = db.Model(&entity).Clauses(clause.Returning{}).Where("platform_account_code = ?", platformAccountCode).Where("profile_collection_id = ?", collectionId).Updates(&entity)
	if db.Error != nil {
		return nil, db.Error
	} else if db.RowsAffected == 0 {
		return nil, fmt.Errorf("invalid Platform Account Code - %v", platformAccountCode)

	}
	return &entity, nil
}

func (d *ItemDao) FetchCollectionItemsByPlatformAccountCode(ctx context.Context, collectionId int64, platformAccountCodes []int64, platform string) ([]ProfileCollectionItemEntity, error) {
	db := d.DaoProvider.GetSession(ctx)
	// fetch last 2 recently edited collections
	var entities []ProfileCollectionItemEntity
	db = db.Model(&ProfileCollectionItemEntity{}).Where("profile_collection_id = ? ", collectionId).Where("platform_account_code IN (?) ", platformAccountCodes).Where("platform = ? ", platform).Find(&entities)
	if db.Error != nil {
		return nil, db.Error
	}
	return entities, nil
}

func (d *ItemDao) UpdateItemCustomColumnsData(ctx context.Context, collectionItemId int64, customColumns []domain.CustomColumn) error {
	existingCCRows, err := d.FetchExistingCustomColumnRowsForItem(ctx, collectionItemId)
	if err != nil {
		return err
	}

	columnKeyWiseRowMap := map[string]ProfileCollectionItemCustomColumnEntity{}
	for _, ccRow := range existingCCRows {
		columnKeyWiseRowMap[ccRow.Key] = ccRow
	}

	for i := range customColumns {
		column := customColumns[i]
		if column.Key == "" {
			return errors.New("custom column key cannot be empty")
		}
		db := d.DaoProvider.GetSession(ctx)
		existingRow, ok := columnKeyWiseRowMap[column.Key]
		if ok {
			db = db.Model(&ProfileCollectionItemCustomColumnEntity{}).Where("id = ?", existingRow.Id).Update("value", column.Value)
			if db.Error != nil {
				return db.Error
			}
		} else {
			newEntity := ProfileCollectionItemCustomColumnEntity{
				ProfileCollectionItemId: collectionItemId,
				Key:                     column.Key,
				Value:                   column.Value,
			}
			db = db.Model(&ProfileCollectionItemCustomColumnEntity{}).Create(&newEntity)
			if db.Error != nil {
				return err
			}
		}
	}
	return nil
}

func (d *ItemDao) FetchExistingCustomColumnRowsForItem(ctx context.Context, collectionItemId int64) ([]ProfileCollectionItemCustomColumnEntity, error) {
	db := d.DaoProvider.GetSession(ctx)
	var entities []ProfileCollectionItemCustomColumnEntity
	db = db.Model(&ProfileCollectionItemCustomColumnEntity{}).Where("profile_collection_item_id = ? ", collectionItemId).Find(&entities)
	if db.Error != nil {
		return nil, db.Error
	}
	return entities, nil
}

func (d *ItemDao) FetchColumnSumValuesForCollectionItems(ctx context.Context, collectionId int64) ([]CollectionItemPlatformCustomColumnSum, *string) {
	db := d.DaoProvider.GetSession(ctx)
	var entities []CollectionItemPlatformCustomColumnSum
	query := "select i.platform, ccd.key, count(ccd.profile_collection_item_id) items, sum(cast (CASE WHEN ccd.value is null or ccd.value = '' then '0' else ccd.value end as INTEGER)) as sum from profile_collection_item i join profile_collection_item_cc_data ccd on i.id = ccd.profile_collection_item_id where i.profile_collection_id = @collectionId and i.enabled = true group by i.platform, ccd.key"
	whereArgs := map[string]interface{}{"collectionId": collectionId}
	res := db.Raw(query, whereArgs).Find(&entities)
	if res.Error != nil {
		err := res.Error.Error()
		return entities, &err
	}
	return entities, nil
}

func (d *ItemDao) FetchCollectionIdsContainingPlatformAccountCode(ctx context.Context, platformAccountCode int64, partnerId int64) ([]string, error) {
	db := d.DaoProvider.GetSession(ctx)
	type Rows struct {
		Collectionid int64
	}
	var collectionIds []string
	var rows []Rows
	db = db.Model(&ProfileCollectionItemEntity{}).Select("profile_collection_id as collectionId").Where(ProfileCollectionItemEntity{PlatformAccountCode: platformAccountCode, Enabled: helpers.ToBool(true), PartnerId: partnerId}).Find(&rows)
	if db.Error != nil {
		return collectionIds, db.Error
	}

	for i := range rows {
		collectionIds = append(collectionIds, strconv.FormatInt(rows[i].Collectionid, 10))
	}
	return collectionIds, nil
}

// -------------------------------
// winkl profile collection migration info dao
// -------------------------------

type WinklCollectionInfoDao struct {
	*rest.Dao[WinklCollectionInfoEntity, int64]
	ctx context.Context
}

func (d *WinklCollectionInfoDao) winklCollectionInfoDaoInit(ctx context.Context) {
	d.Dao = rest.NewDao[WinklCollectionInfoEntity, int64](ctx)
	pgDaoProvider := postgres.NewPgDao[WinklCollectionInfoEntity, int64](ctx)
	d.Dao.DaoProvider = pgDaoProvider
	d.ctx = ctx
}

func CreateWinklCollectionInfoDao(ctx context.Context) *WinklCollectionInfoDao {
	dao := &WinklCollectionInfoDao{}
	dao.winklCollectionInfoDaoInit(ctx)
	return dao
}

func (d *WinklCollectionInfoDao) FetchWinklCollectionInfoById(ctx context.Context, collectionId int64) (*WinklCollectionInfoEntity, error) {
	db := d.DaoProvider.GetSession(ctx)
	var entity WinklCollectionInfoEntity
	db = db.Model(&WinklCollectionInfoEntity{}).Where(WinklCollectionInfoEntity{WinklCollectionId: collectionId}).First(&entity)
	if db.Error != nil && strings.ToLower(db.Error.Error()) == "record not found" {
		return nil, nil
	}
	if db.Error != nil {
		return nil, db.Error
	}

	return &entity, nil
}

func (d *WinklCollectionInfoDao) FetchWinklCollectionInfoByShareId(ctx context.Context, collectionShareId string) (*WinklCollectionInfoEntity, error) {
	db := d.DaoProvider.GetSession(ctx)
	var entity WinklCollectionInfoEntity
	db = db.Model(&WinklCollectionInfoEntity{}).Where(WinklCollectionInfoEntity{WinklCollectionShareId: collectionShareId}).First(&entity)
	if db.Error != nil && strings.ToLower(db.Error.Error()) == "record not found" {
		return nil, nil
	}
	if db.Error != nil {
		return nil, db.Error
	}

	return &entity, nil
}

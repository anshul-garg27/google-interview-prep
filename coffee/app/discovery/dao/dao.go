package dao

import (
	"coffee/constants"
	coredoamin "coffee/core/domain"
	"coffee/core/persistence/postgres"
	"coffee/core/rest"
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"gorm.io/gorm"
)

type InstagramAccountDao struct {
	*rest.Dao[InstagramAccountEntity, int64]
	ctx context.Context
}

func CreateInstagramAccountDao(ctx context.Context) *InstagramAccountDao {
	dao := &InstagramAccountDao{
		Dao: rest.NewDao[InstagramAccountEntity, int64](ctx),
		ctx: ctx,
	}
	dao.DaoProvider = postgres.NewPgDao[InstagramAccountEntity, int64](ctx)
	return dao
}

type YoutubeAccountDao struct {
	*rest.Dao[YoutubeAccountEntity, int64]
	ctx context.Context
}

func CreateYoutubeAccountDao(ctx context.Context) *YoutubeAccountDao {
	dao := &YoutubeAccountDao{
		ctx: ctx,
		Dao: rest.NewDao[YoutubeAccountEntity, int64](ctx),
	}
	dao.DaoProvider = postgres.NewPgDao[YoutubeAccountEntity, int64](ctx)
	return dao
}

type SocialProfileTimeSeriesDao struct {
	*rest.Dao[SocialProfileTimeSeriesEntity, int64]
	ctx context.Context
}

func CreateSocialProfileTimeSeriesDao(ctx context.Context) *SocialProfileTimeSeriesDao {
	dao := &SocialProfileTimeSeriesDao{
		ctx: ctx,
		Dao: rest.NewDao[SocialProfileTimeSeriesEntity, int64](ctx),
	}
	dao.DaoProvider = postgres.NewPgDao[SocialProfileTimeSeriesEntity, int64](ctx)
	return dao
}

type SocialProfileHashtagsDao struct {
	*rest.Dao[SocialProfileHashatagsEntity, int64]
	ctx context.Context
}

func CreateSocialProfileHashtagssDao(ctx context.Context) *SocialProfileHashtagsDao {
	dao := &SocialProfileHashtagsDao{
		ctx: ctx,
		Dao: rest.NewDao[SocialProfileHashatagsEntity, int64](ctx),
	}
	dao.DaoProvider = postgres.NewPgDao[SocialProfileHashatagsEntity, int64](ctx)
	return dao
}

type LocationsDao struct {
	*rest.Dao[LocationsEntity, int64]
	ctx context.Context
}

func CreateLocationsDao(ctx context.Context) *LocationsDao {
	dao := &LocationsDao{
		ctx: ctx,
		Dao: rest.NewDao[LocationsEntity, int64](ctx),
	}
	dao.DaoProvider = postgres.NewPgDao[LocationsEntity, int64](ctx)
	return dao
}

type SocialProfileAudienceInfoDao struct {
	*rest.Dao[SocialProfileAudienceInfoEntity, int64]
	ctx context.Context
}

func CreateSocialProfileAudienceInfoDao(ctx context.Context) *SocialProfileAudienceInfoDao {
	dao := &SocialProfileAudienceInfoDao{
		ctx: ctx,
		Dao: rest.NewDao[SocialProfileAudienceInfoEntity, int64](ctx),
	}
	dao.DaoProvider = postgres.NewPgDao[SocialProfileAudienceInfoEntity, int64](ctx)
	return dao
}

type GroupMetricsDao struct {
	*rest.Dao[GroupMetricsEntity, int64]
	ctx context.Context
}

func CreateGroupMetricsDao(ctx context.Context) *GroupMetricsDao {
	dao := &GroupMetricsDao{
		ctx: ctx,
		Dao: rest.NewDao[GroupMetricsEntity, int64](ctx),
	}
	dao.DaoProvider = postgres.NewPgDao[GroupMetricsEntity, int64](ctx)
	return dao
}

func (I *InstagramAccountDao) FindByProfileId(ctx context.Context, id string, isAdmin bool) (*InstagramAccountEntity, error) {
	dbSession := I.DaoProvider.GetSession(ctx)
	var entity InstagramAccountEntity
	var tx *gorm.DB
	if isAdmin {
		tx = dbSession.WithContext(ctx).Preload("CampaignProfile", "platform = ?", string(constants.InstagramPlatform)).Joins("LEFT JOIN campaign_profiles ON instagram_account.id = campaign_profiles.platform_account_id").Model(&entity).Where(InstagramAccountEntity{ProfileId: &id}).First(&entity)
	} else {
		tx = dbSession.WithContext(ctx).Model(&entity).Where(InstagramAccountEntity{ProfileId: &id}).First(&entity)
	}
	if tx.Error != nil && tx.Error.Error() == "record not found" {
		return nil, tx.Error
	}
	return &entity, tx.Error
}

func (Y *YoutubeAccountDao) FindByProfileId(ctx context.Context, id string, isAdmin bool) (*YoutubeAccountEntity, error) {
	dbSession := Y.GetSession(ctx)
	var entity YoutubeAccountEntity
	var tx *gorm.DB
	if isAdmin {
		tx = dbSession.WithContext(ctx).Preload("CampaignProfile", "platform = ?", string(constants.YoutubePlatform)).Joins("LEFT JOIN campaign_profiles ON youtube_account.id = campaign_profiles.platform_account_id").Model(&entity).Where(YoutubeAccountEntity{ProfileId: &id}).First(&entity)
	} else {
		tx = dbSession.WithContext(ctx).Model(&entity).Where(YoutubeAccountEntity{ProfileId: &id}).First(&entity)
	}

	if tx.Error != nil && tx.Error.Error() == "record not found" {
		return nil, tx.Error
	}
	return &entity, tx.Error
}

func (I *InstagramAccountDao) FindByGccProfileId(ctx context.Context, id string, isAdmin bool) (*InstagramAccountEntity, error) {
	dbSession := I.DaoProvider.GetSession(ctx)
	var entity InstagramAccountEntity
	var tx *gorm.DB
	if isAdmin {
		tx = dbSession.WithContext(ctx).Preload("CampaignProfile", "platform = ?", string(constants.InstagramPlatform)).Joins("LEFT JOIN campaign_profiles ON instagram_account.id = campaign_profiles.platform_account_id").Model(&entity).Where(InstagramAccountEntity{GccProfileId: &id}).First(&entity)
	} else {
		tx = dbSession.WithContext(ctx).Model(&entity).Where(InstagramAccountEntity{GccProfileId: &id}).First(&entity)
	}
	if tx.Error != nil && tx.Error.Error() == "record not found" {
		return nil, tx.Error
	}
	return &entity, tx.Error
}

func (Y *YoutubeAccountDao) FindByGccProfileId(ctx context.Context, id string, isAdmin bool) (*YoutubeAccountEntity, error) {
	dbSession := Y.GetSession(ctx)
	var entity YoutubeAccountEntity
	var tx *gorm.DB
	if isAdmin {
		tx = dbSession.WithContext(ctx).Preload("CampaignProfile", "platform = ?", string(constants.YoutubePlatform)).Joins("LEFT JOIN campaign_profiles ON youtube_account.id = campaign_profiles.platform_account_id").Model(&entity).Where(YoutubeAccountEntity{GccProfileId: &id}).First(&entity)
	} else {
		tx = dbSession.WithContext(ctx).Model(&entity).Where(YoutubeAccountEntity{GccProfileId: &id}).First(&entity)
	}

	if tx.Error != nil && tx.Error.Error() == "record not found" {
		return nil, tx.Error
	}
	return &entity, tx.Error
}

func (Y *YoutubeAccountDao) FindAllYtAccountsByGccProfileId(ctx context.Context, gccProfileId string) (*[]YoutubeAccountEntity, error) {
	dbSession := Y.GetSession(ctx)
	var entities []YoutubeAccountEntity
	tx := dbSession.WithContext(ctx).Model(&YoutubeAccountEntity{}).Where(YoutubeAccountEntity{GccProfileId: &gccProfileId}).Find(&entities)
	if tx.Error != nil && tx.Error.Error() == "record not found" {
		return nil, tx.Error
	}
	return &entities, tx.Error
}

func (I *InstagramAccountDao) FindByHandle(ctx context.Context, handle string, campaignProfileJoin bool) (*InstagramAccountEntity, error) {
	dbSession := I.GetSession(ctx)
	var entity InstagramAccountEntity
	var tx *gorm.DB
	if campaignProfileJoin {
		tx = dbSession.WithContext(ctx).Preload("CampaignProfile").Joins("LEFT JOIN campaign_profiles ON instagram_account.id = campaign_profiles.platform_account_id AND campaign_profiles.platform = 'INSTAGRAM'").Model(&entity).Where(InstagramAccountEntity{Handle: &handle}).First(&entity)
	} else {
		tx = dbSession.WithContext(ctx).Model(&entity).Where(InstagramAccountEntity{Handle: &handle}).First(&entity)
	}
	if tx.Error != nil && tx.Error.Error() == "record not found" {
		return nil, tx.Error
	}
	return &entity, tx.Error
}

func (Y *YoutubeAccountDao) FindByHandle(ctx context.Context, handle string, campaignProfileJoin bool) (*YoutubeAccountEntity, error) {
	dbSession := Y.GetSession(ctx)
	var entity YoutubeAccountEntity
	var tx *gorm.DB

	if campaignProfileJoin {
		tx = dbSession.WithContext(ctx).Preload("CampaignProfile").Joins("LEFT JOIN campaign_profiles ON youtube_account.id = campaign_profiles.platform_account_id AND campaign_profiles.platform = 'YOUTUBE'").Model(&entity).Where(YoutubeAccountEntity{ChannelId: &handle}).First(&entity)
	} else {
		tx = dbSession.WithContext(ctx).Model(&entity).Where(YoutubeAccountEntity{ChannelId: &handle}).First(&entity)
	}
	if tx.Error != nil && tx.Error.Error() == "record not found" {
		return nil, tx.Error
	}
	return &entity, tx.Error
}

func (I *InstagramAccountDao) FindByPlatformProfileIds(ctx context.Context, platformProfileId []int64) (*[]InstagramAccountEntity, error) {
	dbSession := I.GetSession(ctx)
	var entities []InstagramAccountEntity
	var model InstagramAccountEntity

	tx := dbSession.WithContext(ctx).Model(&model).Where("id IN ?", platformProfileId).Find(&entities)
	if tx.Error != nil && tx.Error.Error() == "record not found" {
		return nil, tx.Error
	}
	return &entities, tx.Error
}
func (Y *YoutubeAccountDao) FindByPlatformProfileIds(ctx context.Context, platformProfileId []int64) (*[]YoutubeAccountEntity, error) {
	dbSession := Y.GetSession(ctx)
	var entities []YoutubeAccountEntity
	var model YoutubeAccountEntity
	tx := dbSession.WithContext(ctx).Model(&model).Where("id IN ?", platformProfileId).Find(&entities)
	if tx.Error != nil && tx.Error.Error() == "record not found" {
		return nil, tx.Error
	}
	return &entities, tx.Error
}
func (I *InstagramAccountDao) FindByIgIds(ctx context.Context, igIds []string) (*[]InstagramAccountEntity, error) {
	dbSession := I.GetSession(ctx)
	var entities []InstagramAccountEntity
	var model InstagramAccountEntity

	tx := dbSession.WithContext(ctx).Model(&model).Where("ig_id IN ?", igIds).Find(&entities)
	if tx.Error != nil && tx.Error.Error() == "record not found" {
		return nil, tx.Error
	}
	return &entities, tx.Error
}
func (Y *YoutubeAccountDao) FindByChanneldIds(ctx context.Context, channelIds []string) (*[]YoutubeAccountEntity, error) {
	dbSession := Y.GetSession(ctx)
	var entities []YoutubeAccountEntity
	var model YoutubeAccountEntity
	tx := dbSession.WithContext(ctx).Model(&model).Where("channel_id IN ?", channelIds).Find(&entities)
	if tx.Error != nil && tx.Error.Error() == "record not found" {
		return nil, tx.Error
	}
	return &entities, tx.Error
}

func (S *SocialProfileTimeSeriesDao) GetTimeSeriesData(ctx context.Context, statDate time.Time, endDate time.Time, platform string, platformProfileId int64) (*[]SocialProfileTimeSeriesEntity, error) {
	dbSession := S.GetSession(ctx)
	var entities []SocialProfileTimeSeriesEntity
	var model SocialProfileTimeSeriesEntity
	tx := dbSession.WithContext(ctx).Model(&model).Where("platform = ? and platform_profile_id = ? and date between ? and ?", platform, platformProfileId, statDate, endDate).Find(&entities)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &entities, tx.Error
}

func (S *SocialProfileHashtagsDao) FindByPlatformProfileId(ctx context.Context, platform string, platformProfileId int64) (*SocialProfileHashatagsEntity, error) {
	dbSession := S.GetSession(ctx)
	var entity SocialProfileHashatagsEntity
	tx := dbSession.WithContext(ctx).Model(&entity).Where(SocialProfileHashatagsEntity{Platform: platform, PlatformProfileId: platformProfileId}).First(&entity)
	if tx.Error != nil && tx.Error.Error() == "record not found" {
		return nil, tx.Error
	}
	return &entity, nil
}

func (S *SocialProfileAudienceInfoDao) FindByPlatformProfileId(ctx context.Context, platform string, platformProfileId int64) (*SocialProfileAudienceInfoEntity, error) {
	dbSession := S.GetSession(ctx)
	var entity SocialProfileAudienceInfoEntity
	tx := dbSession.WithContext(ctx).Model(&entity).Where(SocialProfileAudienceInfoEntity{Platform: platform, PlatformProfileId: platformProfileId}).First(&entity)
	if tx.Error != nil && tx.Error.Error() == "record not found" {
		return nil, tx.Error
	}
	return &entity, nil
}
func (S *GroupMetricsDao) FindByGroupKey(ctx context.Context, group_key string) (*GroupMetricsEntity, error) {
	dbSession := S.GetSession(ctx)
	var entity GroupMetricsEntity
	tx := dbSession.WithContext(ctx).Model(&entity).Where(GroupMetricsEntity{GroupKey: group_key, Enabled: true}).First(&entity)
	if tx.Error != nil && tx.Error.Error() == "record not found" {
		return nil, tx.Error
	}
	return &entity, nil

}

func (Y YoutubeAccountDao) GetYoutubeDetalsByChannelIds(ctx context.Context, channelIdArray []*string) (*[]YoutubeAccountEntity, error) {
	dbSession := Y.GetSession(ctx)
	var entities []YoutubeAccountEntity
	var model YoutubeAccountEntity
	tx := dbSession.WithContext(ctx).Model(&model).Where("channel_id IN ?", channelIdArray).Find(&entities)
	if tx.Error != nil && tx.Error.Error() == "record not found" {
		return nil, tx.Error
	}
	return &entities, tx.Error
}
func (I InstagramAccountDao) GetInstagramDetalsByHandle(ctx context.Context, handleArray []*string) (*[]InstagramAccountEntity, error) {
	dbSession := I.GetSession(ctx)
	var entities []InstagramAccountEntity
	var model InstagramAccountEntity
	tx := dbSession.WithContext(ctx).Model(&model).Where("handle IN ?", handleArray).Find(&entities)
	if tx.Error != nil && tx.Error.Error() == "record not found" {
		return nil, tx.Error
	}
	return &entities, tx.Error
}

func (I *InstagramAccountDao) FindByIgId(ctx context.Context, igId string) (*InstagramAccountEntity, error) {
	dbSession := I.GetSession(ctx)
	var entity InstagramAccountEntity
	tx := dbSession.WithContext(ctx).Model(&entity).Where(InstagramAccountEntity{IgID: &igId}).First(&entity)
	if tx.Error != nil && tx.Error.Error() == "record not found" {
		return nil, tx.Error
	}
	return &entity, tx.Error
}

func (I *LocationsDao) FindLocationList(ctx context.Context, name string, table_name string, size int) ([]string, error) {
	dbSession := I.GetSession(ctx)
	var distinctLocations []string

	tx := dbSession.Table(table_name).
		Select("location").
		Where("location ILIKE ?", "%"+name+"%").
		Limit(size).
		Find(&distinctLocations)
	if tx.Error == nil && len(distinctLocations) == 0 {
		tx.Error = errors.New("record not found")
	}
	return distinctLocations, tx.Error
}

func (I *InstagramAccountDao) FindById(ctx context.Context, id int64, campaignProfileJoin bool) (*InstagramAccountEntity, error) {
	dbSession := I.GetSession(ctx)
	var entity InstagramAccountEntity
	var tx *gorm.DB
	if campaignProfileJoin {
		tx = dbSession.WithContext(ctx).Preload("CampaignProfile").Joins("LEFT JOIN campaign_profiles ON instagram_account.id = campaign_profiles.platform_account_id AND campaign_profiles.platform = 'INSTAGRAM'").Model(&entity).Where(InstagramAccountEntity{ID: id}).First(&entity)
	} else {
		tx = dbSession.WithContext(ctx).Model(&entity).Where(InstagramAccountEntity{ID: id}).First(&entity)
	}
	if tx.Error != nil && tx.Error.Error() == "record not found" {
		return nil, tx.Error
	}
	return &entity, tx.Error
}

func (Y *YoutubeAccountDao) FindById(ctx context.Context, id int64, campaignProfileJoin bool) (*YoutubeAccountEntity, error) {
	dbSession := Y.GetSession(ctx)
	var entity YoutubeAccountEntity
	var tx *gorm.DB

	if campaignProfileJoin {
		tx = dbSession.WithContext(ctx).Preload("CampaignProfile").Joins("LEFT JOIN campaign_profiles ON youtube_account.id = campaign_profiles.platform_account_id AND campaign_profiles.platform = 'YOUTUBE'").Model(&entity).Where(YoutubeAccountEntity{ID: id}).First(&entity)
	} else {
		tx = dbSession.WithContext(ctx).Model(&entity).Where(YoutubeAccountEntity{ID: id}).First(&entity)
	}
	if tx.Error != nil && tx.Error.Error() == "record not found" {
		return nil, tx.Error
	}
	return &entity, tx.Error
}

func (I *InstagramAccountDao) SearchInstagramProfiles(ctx context.Context, query coredoamin.SearchQuery, sortByArray []string, sortDir string, page int, size int, joinTables []coredoamin.JoinClauses, computeCounts bool) ([]InstagramAccountEntity, int64, error) {
	db := I.GetSession(ctx)
	var entities []InstagramAccountEntity
	filteredReq := db.Model(&entities)
	// computeCounts := shouldComputeCounts(query, page)

	for key := range joinTables {
		if joinTables[key].PreloadVariable != "" {
			preloadVariable := joinTables[key].PreloadVariable

			preloadCondition := joinTables[key].PreloadCondition
			preloadValue := joinTables[key].PreloadValue

			if preloadCondition == "" || preloadValue == "" {
				filteredReq = filteredReq.Preload(preloadVariable)
			} else {
				filteredReq = filteredReq.Preload(preloadVariable, preloadCondition, preloadValue)
			}
		}
		joinClause := joinTables[key].JoinClause
		filteredReq = filteredReq.Joins(joinClause)
	}
	var CampaignIdNotInClause bool
	var subquery *gorm.DB
	for _, filter := range query.Filters {
		if filter.Field == "profile_collection.campaign_id" && filter.FilterType == "NE" {
			CampaignIdNotInClause = true
			subquery = db.Table("profile_collection_item").
				Select("DISTINCT platform_account_code").Joins("JOIN profile_collection on profile_collection.id = profile_collection_item.profile_collection_id")
		}
	}
	for i := range query.Filters {
		filter := query.Filters[i]
		if CampaignIdNotInClause && (strings.Contains(filter.Field, "profile_collection_item.") || strings.Contains(filter.Field, "profile_collection.")) {
			I.makeSubqueryForProfileCollectionItems(subquery, filter)
		} else if filter.Value != "" && !(filter.Field == "profile_collection.campaign_id" && filter.FilterType == "NE") {
			I.Dao.AddPredicateForSearchJoins(filteredReq, filter)
		}
	}
	subqueryWhereClause := "instagram_account.id"
	if CampaignIdNotInClause {
		filteredReq.Where(subqueryWhereClause+" NOT IN (?)", subquery)
	}

	orderClauses := []string{}
	for i, column := range sortByArray {
		var orderClause string
		if i == 0 {
			orderClause = fmt.Sprintf("%s %s NULLS LAST", column, sortDir)
		} else {
			orderClause = fmt.Sprintf("%s %s", column, sortDir)
		}
		orderClauses = append(orderClauses, orderClause)
	}

	orderBy := strings.Join(orderClauses, ",")
	filteredReq.Order(orderBy)

	var filteredCount int64
	filteredCount = 0
	if computeCounts {
		filteredReq.Count(&filteredCount)
	} else {
		filteredCount = -1
	}

	paginatedReq := filteredReq.Limit(size)
	paginatedReq = paginatedReq.Offset((page - 1) * size)
	paginatedReq.Find(&entities)
	if len(entities) < size {
		filteredCount = int64(len(entities))
	}
	if filteredReq.Error != nil {
		return nil, 0, filteredReq.Error
	}
	return entities, filteredCount, nil
}

func (I InstagramAccountDao) makeSubqueryForProfileCollectionItems(req *gorm.DB, filter coredoamin.SearchFilter) *gorm.DB {
	if filter.Field == "profile_collection.campaign_id" && filter.FilterType == "NE" {
		filter.FilterType = "EQ"
	}
	req = I.Dao.AddPredicateForSearchJoins(req, filter)
	return req
}

func (Y *YoutubeAccountDao) SearchYoutubeProfiles(ctx context.Context, query coredoamin.SearchQuery, sortByArray []string, sortDir string, page int, size int, joinTables []coredoamin.JoinClauses, computeCounts bool) ([]YoutubeAccountEntity, int64, error) {
	db := Y.GetSession(ctx)
	var entities []YoutubeAccountEntity
	// computeCounts := shouldComputeCounts(query, page)
	filteredReq := db.Model(&entities)
	for key := range joinTables {
		if joinTables[key].PreloadVariable != "" {
			preloadVariable := joinTables[key].PreloadVariable

			preloadCondition := joinTables[key].PreloadCondition
			preloadValue := joinTables[key].PreloadValue

			if preloadCondition == "" || preloadValue == "" {
				filteredReq = filteredReq.Preload(preloadVariable)
			} else {
				filteredReq = filteredReq.Preload(preloadVariable, preloadCondition, preloadValue)
			}

		}

		joinClause := joinTables[key].JoinClause
		filteredReq = filteredReq.Joins(joinClause)
	}
	var CampaignIdNotInClause bool
	var subquery *gorm.DB
	for _, filter := range query.Filters {
		if filter.Field == "profile_collection.campaign_id" && filter.FilterType == "NE" {
			CampaignIdNotInClause = true
			subquery = db.Table("profile_collection_item").
				Select("DISTINCT platform_account_code").Joins("JOIN profile_collection on profile_collection.id = profile_collection_item.profile_collection_id")
		}
	}
	for i := range query.Filters {
		filter := query.Filters[i]
		if CampaignIdNotInClause && (strings.Contains(filter.Field, "profile_collection_item.") || strings.Contains(filter.Field, "profile_collection.")) {
			Y.makeSubqueryForProfileCollectionItems(subquery, filter)
		} else if filter.Value != "" && !(filter.Field == "profile_collection.campaign_id" && filter.FilterType == "NE") {
			Y.Dao.AddPredicateForSearchJoins(filteredReq, filter)
		}
	}
	subqueryWhereClause := "youtube_account.id"
	if CampaignIdNotInClause {
		filteredReq.Where(subqueryWhereClause+" NOT IN (?)", subquery)
	}

	orderClauses := []string{}
	for i, column := range sortByArray {
		var orderClause string
		if i == 0 {
			orderClause = fmt.Sprintf("%s %s NULLS LAST", column, sortDir)
		} else {
			orderClause = fmt.Sprintf("%s %s", column, sortDir)
		}
		orderClauses = append(orderClauses, orderClause)
	}

	orderBy := strings.Join(orderClauses, ",")
	filteredReq.Order(orderBy)

	var filteredCount int64
	rand.Seed(time.Now().UnixNano())
	//filteredCount = int64(rand.Intn(200000-10+1) + 10)
	filteredCount = 0
	if computeCounts {
		filteredReq.Count(&filteredCount)
	} else {
		filteredCount = -1
	}

	paginatedReq := filteredReq.Limit(size)
	paginatedReq = paginatedReq.Offset((page - 1) * size)
	paginatedReq.Find(&entities)
	if len(entities) < size {
		filteredCount = int64(len(entities))
	}
	if filteredReq.Error != nil {
		return nil, 0, filteredReq.Error
	}
	return entities, filteredCount, nil
}

// func shouldComputeCounts(query coredoamin.SearchQuery, page int) bool {
// 	computeCounts := false
// 	campaignIdFilter := rest.GetFilterForKey(query.Filters, "profile_collection.campaign_id")
// 	collectionIdFilter := rest.GetFilterForKey(query.Filters, "profile_collection.id")
// 	if campaignIdFilter != nil || collectionIdFilter != nil { // Campaign Related Pages
// 		computeCounts = true
// 	}
// 	if collectionIdFilter != nil && collectionIdFilter.FilterType == "NE" { // Disable for eligible influencers
// 		computeCounts = false
// 	}
// 	if campaignIdFilter != nil && campaignIdFilter.FilterType == "NE" { // Disable for eligible influencers
// 		computeCounts = false
// 	}
// 	coreFilters := []string{"profile_collection.campaign_id", "profile_collection.id", "profile_collection_item.platform", "profile_collection_item.enabled", "campaign_profiles.on_gcc", "instagram_account.platform", "youtube_account.followers", "instagram_account.is_blacklisted,campaign_profiles.admin_details.isBlacklisted", "youtube_account.is_blacklisted,campaign_profiles.admin_details.isBlacklisted", "instagram_account.is_blacklisted", "youtube_account.is_blacklisted"}
// 	nonCoreFilterCount := 0
// 	for i := range query.Filters {
// 		if !slices.Contains(coreFilters, query.Filters[i].Field) {
// 			nonCoreFilterCount += 1
// 		}
// 	}
// 	if nonCoreFilterCount > 0 {
// 		computeCounts = true
// 	}
// 	if page >= 2 {
// 		computeCounts = false
// 	}
// 	return computeCounts
// }

func (Y YoutubeAccountDao) makeSubqueryForProfileCollectionItems(req *gorm.DB, filter coredoamin.SearchFilter) *gorm.DB {
	if filter.Field == "profile_collection.campaign_id" && filter.FilterType == "NE" {
		filter.FilterType = "EQ"
	}
	req = Y.Dao.AddPredicateForSearchJoins(req, filter)
	return req
}

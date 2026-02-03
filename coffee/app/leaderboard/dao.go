package leaderboard

import (
	"coffee/core/persistence/postgres"
	"coffee/core/rest"
	"context"
	"time"
)

type LeaderboardDao struct {
	*rest.Dao[LeaderboardEntity, int64]
	ctx context.Context
}

func createLeaderboardDao(ctx context.Context) *LeaderboardDao {
	dao := &LeaderboardDao{
		ctx: ctx,
		Dao: rest.NewDao[LeaderboardEntity, int64](ctx),
	}
	dao.DaoProvider = postgres.NewPgDao[LeaderboardEntity, int64](ctx)
	return dao
}

type SocialProfileTimeSeriesDao struct {
	*rest.Dao[SocialProfileTimeSeriesEntity, int64]
	ctx context.Context
}

func createSocialProfileTimeSeriesDao(ctx context.Context) *SocialProfileTimeSeriesDao {
	dao := &SocialProfileTimeSeriesDao{
		ctx: ctx,
		Dao: rest.NewDao[SocialProfileTimeSeriesEntity, int64](ctx),
	}
	dao.DaoProvider = postgres.NewPgDao[SocialProfileTimeSeriesEntity, int64](ctx)
	return dao
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

func (S *SocialProfileTimeSeriesDao) GetMultipleTimeSeriesData(ctx context.Context, statDate time.Time, endDate time.Time, platform string, platformProfileIdArray []int64) (*ProfilesTimeSeries, error) {
	dbSession := S.GetSession(ctx)
	var timeSeries ProfilesTimeSeries
	var entities []SocialProfileTimeSeriesEntity
	var model SocialProfileTimeSeriesEntity
	tx := dbSession.WithContext(ctx).Model(&model).Where("monthly_stats = ? and platform = ? and platform_profile_id IN ? and date between ? and ?", false, platform, platformProfileIdArray, statDate, endDate).Find(&entities)
	if tx.Error != nil {
		return &timeSeries, tx.Error
	}
	var followersMap = make(map[int64]MetricTimeSeries)
	var viewMap = make(map[int64]MetricTimeSeries)
	var playsMap = make(map[int64]MetricTimeSeries)
	var totalViewMap = make(map[int64]MetricTimeSeries)
	var totalPlaysMap = make(map[int64]MetricTimeSeries)

	for _, entity := range entities {
		followersMap[entity.PlatformProfileId] = make(map[string]*int64)
		viewMap[entity.PlatformProfileId] = make(map[string]*int64)
		playsMap[entity.PlatformProfileId] = make(map[string]*int64)
		totalViewMap[entity.PlatformProfileId] = make(map[string]*int64)
		totalPlaysMap[entity.PlatformProfileId] = make(map[string]*int64)
	}

	for _, entity := range entities {
		date := entity.Date.Format("2006-01-02")
		followers := entity.Followers
		views := entity.Views
		plays := entity.Plays
		totalViews := entity.ViewsTotal
		totalPlays := entity.PlaysTotal

		followersMap[entity.PlatformProfileId][date] = followers
		viewMap[entity.PlatformProfileId][date] = views
		playsMap[entity.PlatformProfileId][date] = plays
		totalViewMap[entity.PlatformProfileId][date] = totalViews
		totalPlaysMap[entity.PlatformProfileId][date] = totalPlays

	}
	timeSeries.FollowerGraph = followersMap
	timeSeries.ViewsGraph = viewMap
	timeSeries.PlaysGraph = playsMap
	timeSeries.TotalViewGraph = totalViewMap
	timeSeries.TotalPlaysGraph = totalPlaysMap
	return &timeSeries, tx.Error
}

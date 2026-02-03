package manager

import (
	"coffee/app/discovery/dao"
	"coffee/app/discovery/domain"
	coredomain "coffee/core/domain"
	"context"
	"strconv"
	"time"
)

type TimeSeriesManager struct {
	socialProfileTimeSeriesDao *dao.SocialProfileTimeSeriesDao
}

func CreateManagerTimeSeries(ctx context.Context) *TimeSeriesManager {
	socialProfileTimeSeriesDao := dao.CreateSocialProfileTimeSeriesDao(ctx)
	managerTimeSeries := &TimeSeriesManager{
		socialProfileTimeSeriesDao: socialProfileTimeSeriesDao,
	}
	return managerTimeSeries
}

func (m *TimeSeriesManager) FindTimeSeriesByPlatformProfileId(ctx context.Context, platform string, platformProfileId int64, query coredomain.SearchQuery, sortBy string, sortDir string, page int, size int) (*domain.Growth, int64, error) {
	var growth *domain.Growth
	var monthly_stats bool
	for _, filter := range query.Filters {
		if filter.Field == "monthly_stats" && filter.Value == "true" {
			monthly_stats = true
		}
	}
	query.Filters = append(query.Filters, coredomain.SearchFilter{
		FilterType: "EQ",
		Field:      "platform",
		Value:      platform,
	})
	query.Filters = append(query.Filters, coredomain.SearchFilter{
		FilterType: "EQ",
		Field:      "platform_profile_id",
		Value:      strconv.FormatInt(platformProfileId, 10),
	})
	for i := range query.Filters {
		// if query.Filters[i].FilterType == "GTE" && query.Filters[i].Field == "date" {
		// 	dateStr := query.Filters[i].Value
		// 	date, err := time.Parse("2006-01-02", dateStr)
		// 	if err == nil {
		// 		weekday := date.Weekday()
		// 		daysUntilPreviousSunday := int(weekday) // Sunday is 0 in Weekday enumeration
		// 		previousSunday := date.AddDate(0, 0, -daysUntilPreviousSunday)
		// 		query.Filters[i].Value = previousSunday.Format("2006-01-02")
		// 	}
		// }
		if query.Filters[i].FilterType == "LTE" && query.Filters[i].Field == "date" {
			dateStr := query.Filters[i].Value
			date, err := time.Parse("2006-01-02", dateStr)
			if err == nil {
				weekday := date.Weekday()
				daysUntilNextSunday := 7 - int(weekday)
				nextSunday := date.AddDate(0, 0, daysUntilNextSunday)
				query.Filters[i].Value = nextSunday.Format("2006-01-02")
			}
		}
	}

	entities, filteredCount, err := m.socialProfileTimeSeriesDao.Search(ctx, query, sortBy, sortDir, page, size)
	if err != nil {
		return nil, filteredCount, err
	}
	if len(entities) > 0 {
		growth = GetGrowthFromSocialProfileTimeSeriesEntity(entities, monthly_stats)
	}
	return growth, filteredCount, err
}

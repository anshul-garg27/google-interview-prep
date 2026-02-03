package dao

import (
	coredomain "coffee/core/domain"
	"coffee/core/persistence/clickhouse"
	"coffee/core/rest"
	"context"

	"github.com/spf13/viper"

	"strings"
	"time"
)

// -------------------------------
// collection post metrics summary dao
// -------------------------------

type CollectionPostMetricsTSDAO struct {
	ctx         context.Context
	DaoProvider *clickhouse.ClickhouseDao
}

func (d *CollectionPostMetricsTSDAO) init(ctx context.Context) {
	d.DaoProvider = clickhouse.NewClickhouseDao(ctx)
	d.ctx = ctx
}

func CreateCollectionPostMetricsTSDAO(ctx context.Context) *CollectionPostMetricsTSDAO {
	dao := &CollectionPostMetricsTSDAO{}
	dao.init(ctx)
	return dao
}

func (d *CollectionPostMetricsTSDAO) ConstructWhereClause(collectionIds []string, collectionType string, query coredomain.SearchQuery) (string, map[string]interface{}) {

	postTypeFilter := rest.GetFilterForKey(query.Filters, "postType")
	hashtagsStr := rest.GetFilterValueForKey(query.Filters, "hashtags")
	handleFilter := rest.GetFilterForKey(query.Filters, "handle")
	shortcodeFilter := rest.GetFilterForKey(query.Filters, "shortcode")
	publishedAtAfter := rest.GetFilterValueForKeyOperation(query.Filters, "publishedAt", "GTE")
	publishedAtBefore := rest.GetFilterValueForKeyOperation(query.Filters, "publishedAt", "LTE")

	whereQuery := "collection_id IN @collectionId AND collection_type = @collectionType"
	whereArgs := map[string]interface{}{"collectionId": collectionIds, "collectionType": collectionType}
	if handleFilter != nil {
		if handleFilter.FilterType == "EQ" {
			whereQuery = whereQuery + " AND profile_handle = @handle"
			whereArgs["handle"] = handleFilter.Value
		} else {
			whereQuery = whereQuery + " AND profile_handle IN @handles"
			whereArgs["handles"] = strings.Split(handleFilter.Value, ",")
		}
	}
	if shortcodeFilter != nil {
		if shortcodeFilter.FilterType == "EQ" {
			whereQuery = whereQuery + " AND post_short_code = @shortcode"
			whereArgs["shortcode"] = shortcodeFilter.Value
		} else {
			whereQuery = whereQuery + " AND post_short_code IN @shortcodes"
			whereArgs["shortcodes"] = strings.Split(shortcodeFilter.Value, ",")
		}
	}
	if postTypeFilter != nil {
		if postTypeFilter.FilterType == "EQ" {
			whereQuery = whereQuery + " AND post_type = @postType"
			whereArgs["postType"] = postTypeFilter.Value
		} else {
			whereQuery = whereQuery + " AND post_type IN @postType"
			whereArgs["postType"] = strings.Split(postTypeFilter.Value, ",")
		}
	}
	if publishedAtAfter != nil {
		whereQuery = whereQuery + " AND published_at >= @PAfterTime"
		publishedAtAfter1, _ := time.Parse("2006-01-02", *publishedAtAfter)
		whereArgs["PAfterTime"] = publishedAtAfter1
	}
	if publishedAtBefore != nil {
		whereQuery = whereQuery + " AND published_at <= @PBeforeTime"
		publishedAtBefore1, _ := time.Parse("2006-01-02", *publishedAtBefore)
		whereArgs["PBeforeTime"] = publishedAtBefore1
	}
	if hashtagsStr != nil {
		hashtags := strings.Split(strings.ToLower(*hashtagsStr), ",")
		whereQuery = whereQuery + " AND hasAny(hashtags, arrayMap(x -> tupleElement(x,2), tupleToNameValuePairs(@HT)))"
		whereArgs["HT"] = hashtags
	}

	return whereQuery, whereArgs
}

func (d *CollectionPostMetricsTSDAO) FetchCollectionWeeklyMetricsSummary(ctx context.Context, collectionIds []string, collectionType string, query coredomain.SearchQuery, startTime *time.Time, endTime *time.Time, grouping string) []WeeklyMetricsSummary {
	groupingFunc := "toStartOfWeek"
	if grouping == "daily" {
		groupingFunc = "toStartOfDay"
	}
	db := d.DaoProvider.GetSession(ctx)
	if db == nil {
		return []WeeklyMetricsSummary{}
	}
	whereQuery, whereArgs := d.ConstructWhereClause(collectionIds, collectionType, query)

	rightNow := time.Now()
	if endTime == nil {
		endTime = &rightNow
	}
	if startTime == nil {
		s := (*endTime).AddDate(0, 0, -90)
		startTime = &s
	}
	whereQuery = whereQuery + " AND stats_date >= @SSD"
	whereArgs["SSD"] = time.Date((*startTime).Year(), (*startTime).Month(), (*startTime).Day(), 0, 0, 0, 0, time.UTC)

	whereQuery = whereQuery + " AND stats_date <= @SED"
	whereArgs["SED"] = time.Date((*endTime).Year(), (*endTime).Month(), (*endTime).Day(), 0, 0, 0, 0, time.UTC)

	env := viper.GetString("ENV")

	tableName := "mart_staging_collection_post_ts"
	if env == "PROD" {
		tableName = "mart_collection_post_ts"
	}
	var weeklyMetricsSummaryEntity []WeeklyMetricsSummary
	weeklyMetricsSummaryQuery := "select toDateTime(startOfWeek) as startofweek, count(distinct post_short_code) as totalposts, sum(views) as views, sum(likes) as likes, sum(comments) as comments, sum(impressions) as impressions, sum(reach) as reach from (select " + groupingFunc + "(stats_date) as startOfWeek, post_short_code, argMax(views, stats_date) as views, argMax(likes, stats_date) as likes, argMax(comments, stats_date) as comments, argMax(impressions, stats_date) as impressions, argMax(reach, stats_date) as reach from " + tableName + " where " + whereQuery + " group by startOfWeek, post_short_code) a group by startOfWeek order by startofweek asc"
	db.Raw(weeklyMetricsSummaryQuery, whereArgs).Scan(&weeklyMetricsSummaryEntity)

	return weeklyMetricsSummaryEntity
}

func (d *CollectionPostMetricsTSDAO) FetchCollectionClicksTS(ctx context.Context, collectionIds []string, collectionType string, query coredomain.SearchQuery, startTime *time.Time, endTime *time.Time, grouping string) []WeeklyMetricsSummary {
	groupingFunc := "toStartOfWeek"
	if grouping == "daily" {
		groupingFunc = "toStartOfDay"
	}
	db := d.DaoProvider.GetSession(ctx)
	if db == nil {
		return []WeeklyMetricsSummary{}
	}
	whereQuery, whereArgs := d.ConstructWhereClause(collectionIds, collectionType, query)

	rightNow := time.Now()
	if endTime == nil {
		endTime = &rightNow
	}
	if startTime == nil {
		s := (*endTime).AddDate(0, 0, -90)
		startTime = &s
	}
	whereQuery = whereQuery + " AND stats_date >= @SSD"
	whereArgs["SSD"] = time.Date((*startTime).Year(), (*startTime).Month(), (*startTime).Day(), 0, 0, 0, 0, time.UTC)

	whereQuery = whereQuery + " AND stats_date <= @SED"
	whereArgs["SED"] = time.Date((*endTime).Year(), (*endTime).Month(), (*endTime).Day(), 0, 0, 0, 0, time.UTC)

	env := viper.GetString("ENV")

	tableName := "mart_collection_clicks_ts"
	if env == "PROD" {
		tableName = "mart_collection_clicks_ts"
	} else {
		whereArgs["collectionId"] = []string{"6cefbae1-1a21-446d-a791-999c67f3ab62"}
	}
	var weeklyMetricsSummaryEntity []WeeklyMetricsSummary
	weeklyMetricsSummaryQuery := "select toDateTime(startOfWeek) as startofweek, sum(clicks) as clicks from (select " + groupingFunc + "(stats_date) as startOfWeek, post_short_code, max(clicks) as clicks from " + tableName + " where " + whereQuery + " group by startOfWeek, post_short_code) a group by startOfWeek order by startofweek asc"
	db.Raw(weeklyMetricsSummaryQuery, whereArgs).Scan(&weeklyMetricsSummaryEntity)

	return weeklyMetricsSummaryEntity
}

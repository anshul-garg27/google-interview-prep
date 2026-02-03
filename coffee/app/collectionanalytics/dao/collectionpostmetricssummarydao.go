package dao

import (
	coredomain "coffee/core/domain"
	"coffee/core/persistence/postgres"
	"coffee/core/rest"
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgtype"
)

// -------------------------------
// collection post metrics summary dao
// -------------------------------

type CollectionPostMetricsSummaryDAO struct {
	*rest.Dao[CollectionPostMetricsSummaryEntity, int64]
	ctx context.Context
}

func (d *CollectionPostMetricsSummaryDAO) init(ctx context.Context) {
	d.Dao = rest.NewDao[CollectionPostMetricsSummaryEntity, int64](ctx)
	d.Dao.DaoProvider = postgres.NewPgDao[CollectionPostMetricsSummaryEntity, int64](ctx)
	d.ctx = ctx
}

func CreateCollectionPostMetricsSummaryDAO(ctx context.Context) *CollectionPostMetricsSummaryDAO {
	dao := &CollectionPostMetricsSummaryDAO{}
	dao.init(ctx)
	return dao
}

func (d *CollectionPostMetricsSummaryDAO) ConstructWhereClause(collectionIds []string, collectionType string, query coredomain.SearchQuery) (string, map[string]interface{}) {

	postTypeFilter := rest.GetFilterForKey(query.Filters, "postType")
	hashtagsStr := rest.GetFilterValueForKey(query.Filters, "hashtags")
	publishedAtAfter := rest.GetFilterValueForKeyOperation(query.Filters, "publishedAt", "GTE")
	publishedAtBefore := rest.GetFilterValueForKeyOperation(query.Filters, "publishedAt", "LTE")

	whereQuery := "collection_id IN @collectionId AND collection_type = @collectionType"
	whereArgs := map[string]interface{}{"collectionId": collectionIds, "collectionType": collectionType}
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
	if hashtagsStr != nil && strings.Trim(*hashtagsStr, " ") != "" {
		hashtags := strings.Split(*hashtagsStr, ",")
		var htFilterQuery []string
		for i, hashtag := range hashtags {
			var ht pgtype.JSONB
			err := ht.Set([]string{strings.ToLower(hashtag)})
			if err == nil {
				varName := "HT" + strconv.Itoa(i)
				htFilterQuery = append(htFilterQuery, "hashtags @> @"+varName)
				whereArgs[varName] = ht
			}
		}
		if len(htFilterQuery) > 0 {
			whereQuery = whereQuery + " AND (" + strings.Join(htFilterQuery, " OR ") + ")"
		}
	}

	return whereQuery, whereArgs
}

func (d *CollectionPostMetricsSummaryDAO) FetchCollectionMetricsSummary(ctx context.Context, collectionIds []string, collectionType string, query coredomain.SearchQuery) MetricsSummaryEntity {

	db := d.DaoProvider.GetSession(ctx)
	whereQuery, whereArgs := d.ConstructWhereClause(collectionIds, collectionType, query)

	var metricsSummaryEntity MetricsSummaryEntity
	metricsSummaryQuery := "select count(distinct post_short_code) as totalPosts, sum(views) as views, sum(likes) as likes, sum(comments) as comments, sum(impressions) as impressions, sum(saves) as saves, sum(plays) as plays, sum(reach) as reach, sum(swipe_ups) as swipeUps, sum(mentions) as mentions, sum(sticker_taps) as stickerTaps, sum(shares) as shares, sum(story_exits) as storyExits, sum(story_back_taps) as storyBackTaps, sum(story_forward_taps) as storyForwardTaps, sum(link_clicks) as linkClicks, sum(orders) as orders, sum(delivered_orders) as deliveredOrders, sum(completed_orders) as completedOrders, sum(leaderboard_overall_orders) as leaderboardOverallOrders, sum(leaderboard_delivered_orders) as leaderboardDeliveredOrders, sum(leaderboard_completed_orders) as leaderboardCompletedOrders, max(updated_at) as lastRefreshTime, sum(likes+comments) as totalEngagement, count(distinct profile_handle) as profilesCount, case when sum(reach) > 0 then sum(likes + comments)/sum(reach) else 0 end as er from collection_post_metrics_summary where " + whereQuery
	db.Raw(metricsSummaryQuery, whereArgs).Find(&metricsSummaryEntity)

	type FollowerCount struct {
		Followers int64
	}
	entity := FollowerCount{}
	totalFollowersQuery := "select sum(followers) as followers from (select platform as Platform, profile_handle as profileHandle, max(followers) as followers from collection_post_metrics_summary where " + whereQuery + " group by platform, profile_handle) a"
	tx := db.Raw(totalFollowersQuery, whereArgs).Scan(&entity)
	if tx.Error == nil {
		metricsSummaryEntity.Followers = entity.Followers
	}
	return metricsSummaryEntity
}

func (d *CollectionPostMetricsSummaryDAO) FetchCollectionPostWiseCounts(ctx context.Context, collectionIds []string, collectionType string, query coredomain.SearchQuery) []PostWiseCountEntity {

	db := d.DaoProvider.GetSession(ctx)
	whereQuery, whereArgs := d.ConstructWhereClause(collectionIds, collectionType, query)

	var postWiseCounts []PostWiseCountEntity
	postTypeWiseCountQuery := "select post_type as Posttype, count(distinct post_short_code) as posts from collection_post_metrics_summary where " + whereQuery + " group by post_type"
	db.Raw(postTypeWiseCountQuery, whereArgs).Find(&postWiseCounts)

	return postWiseCounts
}

func (d *CollectionPostMetricsSummaryDAO) FetchCollectionProfilesPostWiseCounts(ctx context.Context, collectionIds []string, collectionType string, query coredomain.SearchQuery) map[string][]PostWiseCountEntity {

	db := d.DaoProvider.GetSession(ctx)
	whereQuery, whereArgs := d.ConstructWhereClause(collectionIds, collectionType, query)

	var profilePostCountsEntities []ProfilePostWiseCountEntity
	postTypeWiseCountQuery := "select platform, profile_handle as profileHandle, post_type as postType, count(distinct post_short_code) as posts from collection_post_metrics_summary where " + whereQuery + " group by platform, profile_handle, post_type"
	db.Raw(postTypeWiseCountQuery, whereArgs).Find(&profilePostCountsEntities)

	profilePostWiseCounts := map[string][]PostWiseCountEntity{}
	for _, row := range profilePostCountsEntities {
		key := strings.ToUpper(row.Platform) + "::" + row.Profilehandle
		existingProfilePostWiseCounts, ok := profilePostWiseCounts[key]
		if !ok {
			existingProfilePostWiseCounts = []PostWiseCountEntity{}
		}
		existingProfilePostWiseCounts = append(existingProfilePostWiseCounts, PostWiseCountEntity{Posttype: row.Posttype, Posts: row.Posts})
		profilePostWiseCounts[key] = existingProfilePostWiseCounts
	}
	return profilePostWiseCounts
}

func (d *CollectionPostMetricsSummaryDAO) FetchCollectionProfiles(ctx context.Context, collectionIds []string, collectionType string, query coredomain.SearchQuery, sortBy string, sortOrder string, page int, size int) []ProfileMetricsSummaryEntity {

	db := d.DaoProvider.GetSession(ctx)
	whereQuery, whereArgs := d.ConstructWhereClause(collectionIds, collectionType, query)

	offset := (page - 1) * size
	orderLimitClause := "order by " + sortBy + " " + strings.ToUpper(sortOrder) + " offset " + strconv.Itoa(offset) + " limit " + strconv.Itoa(size)

	var profileMetrics []ProfileMetricsSummaryEntity
	metricsSummaryQuery := "select platform, profile_handle as handle, profile_social_id, max(profile_name) as name, max(profile_pic) as profilePic, max(followers) as followers, max(cost) as cost, count(distinct post_short_code) as totalPosts, sum(views) as views, sum(likes) as likes, sum(comments) as comments, sum(impressions) as impressions, sum(saves) as saves, sum(plays) as plays, sum(reach) as reach, sum(swipe_ups) as swipeUps, sum(mentions) as mentions, sum(sticker_taps) as stickerTaps, sum(shares) as shares, sum(story_exits) as storyExits, sum(story_back_taps) as storyBackTaps, sum(story_forward_taps) as storyForwardTaps, sum(link_clicks) as linkClicks, sum(orders) as orders, sum(delivered_orders) as deliveredOrders, sum(completed_orders) as completedOrders, sum(leaderboard_overall_orders) as leaderboardOverallOrders, sum(leaderboard_delivered_orders) as leaderboardDeliveredOrders, sum(leaderboard_completed_orders) as leaderboardCompletedOrders, max(updated_at) as lastRefreshTime, sum(likes+comments) as totalEngagement, case when sum(reach) > 0 then sum(likes + comments)/sum(reach) else 0 end as er from collection_post_metrics_summary where " + whereQuery + " group by platform, profile_handle, profile_social_id " + orderLimitClause
	db.Raw(metricsSummaryQuery, whereArgs).Scan(&profileMetrics)

	return profileMetrics
}

func (d *CollectionPostMetricsSummaryDAO) FetchCollectionPosts(ctx context.Context, collectionIds []string, collectionType string, query coredomain.SearchQuery, sortBy string, sortOrder string, page int, size int) []CollectionPostMetricsSummaryEntity {

	db := d.DaoProvider.GetSession(ctx)
	whereQuery, whereArgs := d.ConstructWhereClauseForPosts(collectionIds, collectionType, query)

	offset := (page - 1) * size
	orderLimitClause := "order by " + sortBy + " " + strings.ToUpper(sortOrder) + " offset " + strconv.Itoa(offset) + " limit " + strconv.Itoa(size)

	var postMetrics []CollectionPostMetricsSummaryEntity
	if collectionType == "POST" {
		metricsSummaryQuery := "select cms.*, pci.bookmarked as bookmarked from collection_post_metrics_summary cms left join post_collection_item pci on cms.post_collection_item_id = pci.id where " + whereQuery + " " + orderLimitClause
		db.Raw(metricsSummaryQuery, whereArgs).Scan(&postMetrics)
	} else {
		metricsSummaryQuery := "select cms.*, 'false' as bookmarked from collection_post_metrics_summary cms where " + whereQuery + " " + orderLimitClause
		db.Raw(metricsSummaryQuery, whereArgs).Scan(&postMetrics)
	}

	return postMetrics
}

func (d *CollectionPostMetricsSummaryDAO) ConstructWhereClauseForPosts(collectionIds []string, collectionType string, query coredomain.SearchQuery) (string, map[string]interface{}) {

	postTypeFilter := rest.GetFilterForKey(query.Filters, "postType")
	publishedAtAfter := rest.GetFilterValueForKeyOperation(query.Filters, "publishedAt", "GTE")
	publishedAtBefore := rest.GetFilterValueForKeyOperation(query.Filters, "publishedAt", "LTE")
	bookmarked := rest.GetFilterValueForKey(query.Filters, "bookmarked")
	hashtagsStr := rest.GetFilterValueForKey(query.Filters, "hashtags")

	whereQuery := "cms.collection_id IN @collectionId AND cms.collection_type = @collectionType"
	whereArgs := map[string]interface{}{"collectionId": collectionIds, "collectionType": collectionType}
	if postTypeFilter != nil {
		if postTypeFilter.FilterType == "EQ" {
			whereQuery = whereQuery + " AND cms.post_type = @postType"
			whereArgs["postType"] = postTypeFilter.Value
		} else {
			whereQuery = whereQuery + " AND cms.post_type IN @postType"
			whereArgs["postType"] = strings.Split(postTypeFilter.Value, ",")
		}
	}
	if publishedAtAfter != nil {
		whereQuery = whereQuery + " AND cms.published_at >= @PAfterTime"
		publishedAtAfter1, _ := time.Parse("2006-01-02", *publishedAtAfter)
		whereArgs["PAfterTime"] = publishedAtAfter1
	}
	if publishedAtBefore != nil {
		whereQuery = whereQuery + " AND cms.published_at <= @PBeforeTime"
		publishedAtBefore1, _ := time.Parse("2006-01-02", *publishedAtBefore)
		whereArgs["PBeforeTime"] = publishedAtBefore1
	}
	if collectionType == "POST" && bookmarked != nil {
		whereQuery = whereQuery + " AND pci.bookmarked = @Bookmarked"
		whereArgs["Bookmarked"] = bookmarked
	}
	if hashtagsStr != nil && strings.Trim(*hashtagsStr, " ") != "" {
		hashtags := strings.Split(*hashtagsStr, ",")
		var htFilterQuery []string
		for i, hashtag := range hashtags {
			var ht pgtype.JSONB
			err := ht.Set([]string{strings.ToLower(hashtag)})
			if err == nil {
				varName := "HT" + strconv.Itoa(i)
				htFilterQuery = append(htFilterQuery, "hashtags @> @"+varName)
				whereArgs[varName] = ht
			}
		}
		if len(htFilterQuery) > 0 {
			whereQuery = whereQuery + " AND (" + strings.Join(htFilterQuery, " OR ") + ")"
		}
	}

	return whereQuery, whereArgs
}

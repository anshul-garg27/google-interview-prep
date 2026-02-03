package manager

import (
	"coffee/app/collectionanalytics/dao"
	"coffee/app/collectionanalytics/domain"
	domain2 "coffee/app/collectiongroup/domain"
	collectiongroup "coffee/app/collectiongroup/manager"
	discoverymanager "coffee/app/discovery/manager"
	postcollectionmanager "coffee/app/postcollection/manager"
	profilecollectionmanager "coffee/app/profilecollection/manager"
	"coffee/constants"
	coredomain "coffee/core/domain"
	"coffee/core/rest"
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type CollectionAnalyticsManager struct {
	summaryDAO                *dao.CollectionPostMetricsSummaryDAO
	tsDAO                     *dao.CollectionPostMetricsTSDAO
	hashtagDAO                *dao.CollectionHashTagDAO
	keywordDAO                *dao.CollectionKeywordDAO
	postCollectionManager     *postcollectionmanager.PostCollectionManager
	postCollectionItemManager *postcollectionmanager.PostCollectionItemManager
	profileCollectionManager  *profilecollectionmanager.Manager
	collectionGroupManager    *collectiongroup.Manager
	metricsDao                *dao.MetricsDao
	discoveryManager          *discoverymanager.SearchManager
}

func (m *CollectionAnalyticsManager) init(ctx context.Context) {
	m.summaryDAO = dao.CreateCollectionPostMetricsSummaryDAO(ctx)
	m.tsDAO = dao.CreateCollectionPostMetricsTSDAO(ctx)
	m.hashtagDAO = dao.CreateCollectionHashTagDAO(ctx)
	m.keywordDAO = dao.CreateCollectionKeywordDAO(ctx)
	m.postCollectionManager = postcollectionmanager.CreatePostCollectionManager(ctx)
	m.postCollectionItemManager = postcollectionmanager.CreatePostCollectionItemManager(ctx)
	m.profileCollectionManager = profilecollectionmanager.CreateManager(ctx)
	m.collectionGroupManager = collectiongroup.CreateManager(ctx)
	m.metricsDao = dao.CreatePostCollectionMetricsDao(ctx)
	m.discoveryManager = discoverymanager.CreateSearchManager(ctx)

}

func CreateCollectionAnalyticsManager(ctx context.Context) *CollectionAnalyticsManager {
	manager := &CollectionAnalyticsManager{}
	manager.init(ctx)
	return manager
}

func (m *CollectionAnalyticsManager) ConvertToPostCountsSummary(postWiseCountsData []dao.PostWiseCountEntity) domain.PostCountsSummary {
	postWiseCounts := domain.PostCountsSummary{}
	for _, postTypeCountRow := range postWiseCountsData {
		switch strings.ToLower(postTypeCountRow.Posttype) {
		case "image":
			postWiseCounts.Images = postTypeCountRow.Posts
		case "reels":
			postWiseCounts.Reels = postTypeCountRow.Posts
		case "carousel":
			postWiseCounts.Carousels = postTypeCountRow.Posts
		case "story":
			postWiseCounts.Stories = postTypeCountRow.Posts
		case "short":
			postWiseCounts.Shorts = postTypeCountRow.Posts
		case "video":
			postWiseCounts.Videos = postTypeCountRow.Posts
		}
	}
	return postWiseCounts
}

func (m *CollectionAnalyticsManager) FetchCollectionTimeSeries(ctx context.Context, query coredomain.SearchQuery) ([]domain.CollectionTimeSeries, error) {
	collectionIdFilter := rest.GetFilterValueForKey(query.Filters, "collectionId")
	collectionType := rest.GetFilterValueForKey(query.Filters, "collectionType")
	shareId := rest.GetFilterValueForKey(query.Filters, "shareId")
	rollupFilter := rest.GetFilterValueForKey(query.Filters, "rollup")
	collectionGroupId := rest.GetFilterValueForKey(query.Filters, "collectionGroupId")
	collectionGroupShareId := rest.GetFilterValueForKey(query.Filters, "collectionGroupShareId")
	rollup := "daily"
	if rollupFilter != nil && *rollupFilter == "weekly" {
		rollup = "weekly"
	}
	var collectionIds []string
	if collectionType == nil {
		return nil, errors.New("invalid collection type")
	}
	var err error
	if shareId != nil && collectionIdFilter == nil {
		collectionIdFilter, err = m.getCollectionIdFromShareId(ctx, *shareId, *collectionType)
		if err != nil {
			return nil, err
		}
	}
	if collectionGroupId != nil && collectionIdFilter == nil {
		groupId, _ := strconv.ParseInt(*collectionGroupId, 10, 64)
		collectionIds, err = m.GetCollectionIdFromGroupId(ctx, groupId)
		if err != nil {
			return nil, err
		}
	}
	if collectionGroupId != nil && collectionIdFilter == nil {
		groupId, _ := strconv.ParseInt(*collectionGroupId, 10, 64)
		collectionIds, err = m.GetCollectionIdFromGroupId(ctx, groupId)
		if err != nil {
			return nil, err
		}
	}
	if collectionGroupShareId != nil && collectionIdFilter == nil {
		collectionIds, err = m.GetCollectionIdFromGroupShareId(ctx, *collectionGroupShareId)
		if err != nil {
			return nil, err
		}
	}
	if collectionIdFilter != nil {
		collectionIds = append(collectionIds, *collectionIdFilter)
	}

	if collectionType == nil || collectionIds == nil || len(collectionIds) == 0 {
		return nil, errors.New("invalid collection")
	}

	var startTime *time.Time
	var endTime *time.Time
	if *collectionType == "POST" {
		postCollection, _ := m.postCollectionManager.FindByIdInternal(ctx, collectionIds[0])
		if (*postCollection).StartTime != nil {
			st := time.Unix(*postCollection.StartTime, 0)
			startTime = &st
		}
		if (*postCollection).EndTime != nil {
			et := time.Unix(*postCollection.EndTime, 0)
			endTime = &et
		}
	}
	var weeklyMetrics []domain.CollectionTimeSeries
	timeSeriesData := m.tsDAO.FetchCollectionWeeklyMetricsSummary(ctx, collectionIds, *collectionType, query, startTime, endTime, rollup)
	var firstWeek string
	var currentReach float64
	var lastReach float64
	var currentViews float64
	var lastViews float64
	var currentEngagement int64
	var lastEngagement int64
	lastReach = 0.0
	currentReach = 0.0
	lastViews = 0.0
	currentViews = 0.0
	lastEngagement = 0.0
	currentEngagement = 0.0
	for _, weeklyRow := range timeSeriesData {
		currentReach = weeklyRow.Reach
		currentViews = weeklyRow.Views
		currentEngagement = weeklyRow.Likes + weeklyRow.Comments
		if currentEngagement < lastEngagement {
			currentEngagement = lastEngagement
		} else {
			lastEngagement = currentEngagement
		}
		if currentReach < lastReach {
			currentReach = lastReach
		} else {
			lastReach = currentReach
		}
		if currentViews < lastViews {
			currentViews = lastViews
		} else {
			lastViews = currentViews
		}
		weeklyMetrics = append(weeklyMetrics, domain.CollectionTimeSeries{
			StartOfWeek:     weeklyRow.Startofweek.Format("2006-01-02"),
			Date:            weeklyRow.Startofweek.Format("2006-01-02"),
			TotalPosts:      weeklyRow.Totalposts,
			Views:           int64(currentViews),
			Likes:           weeklyRow.Likes,
			Comments:        weeklyRow.Comments,
			Reach:           int64(currentReach),
			TotalEngagement: currentEngagement,
			Impressions:     int64(weeklyRow.Impressions),
		})
		if firstWeek == "" {
			rollupHours := 7 * 24 * time.Hour
			if rollup == "daily" {
				rollupHours = 1 * 24 * time.Hour
			}
			firstWeek = weeklyRow.Startofweek.Add(rollupHours * -1).Format("2006-01-02")
		}
	}
	if firstWeek != "" {
		weeklyMetrics = append([]domain.CollectionTimeSeries{{
			StartOfWeek:     firstWeek,
			Date:            firstWeek,
			TotalPosts:      0,
			Views:           0,
			Likes:           0,
			Comments:        0,
			Reach:           0,
			Impressions:     0,
			TotalEngagement: 0,
		}}, weeklyMetrics...)
	}
	return weeklyMetrics, nil
}

func (m *CollectionAnalyticsManager) FetchCollectionClicksTimeSeries(ctx context.Context, query coredomain.SearchQuery) ([]domain.CollectionTimeSeries, error) {
	collectionIdFilter := rest.GetFilterValueForKey(query.Filters, "collectionId")
	collectionType := rest.GetFilterValueForKey(query.Filters, "collectionType")
	shareId := rest.GetFilterValueForKey(query.Filters, "shareId")
	rollupFilter := rest.GetFilterValueForKey(query.Filters, "rollup")
	collectionGroupId := rest.GetFilterValueForKey(query.Filters, "collectionGroupId")
	collectionGroupShareId := rest.GetFilterValueForKey(query.Filters, "collectionGroupShareId")
	rollup := "daily"
	if rollupFilter != nil && *rollupFilter == "weekly" {
		rollup = "weekly"
	}
	var collectionIds []string
	if collectionType == nil {
		return nil, errors.New("invalid collection type")
	}
	var err error
	if shareId != nil && collectionIdFilter == nil {
		collectionIdFilter, err = m.getCollectionIdFromShareId(ctx, *shareId, *collectionType)
		if err != nil {
			return nil, err
		}
	}
	if collectionGroupId != nil && collectionIdFilter == nil {
		groupId, _ := strconv.ParseInt(*collectionGroupId, 10, 64)
		collectionIds, err = m.GetCollectionIdFromGroupId(ctx, groupId)
		if err != nil {
			return nil, err
		}
	}
	if collectionGroupId != nil && collectionIdFilter == nil {
		groupId, _ := strconv.ParseInt(*collectionGroupId, 10, 64)
		collectionIds, err = m.GetCollectionIdFromGroupId(ctx, groupId)
		if err != nil {
			return nil, err
		}
	}
	if collectionGroupShareId != nil && collectionIdFilter == nil {
		collectionIds, err = m.GetCollectionIdFromGroupShareId(ctx, *collectionGroupShareId)
		if err != nil {
			return nil, err
		}
	}
	if collectionIdFilter != nil {
		collectionIds = append(collectionIds, *collectionIdFilter)
	}

	if collectionType == nil || collectionIds == nil || len(collectionIds) == 0 {
		return nil, errors.New("invalid collection")
	}

	var startTime *time.Time
	var endTime *time.Time
	var weeklyMetrics []domain.CollectionTimeSeries
	timeSeriesData := m.tsDAO.FetchCollectionClicksTS(ctx, collectionIds, *collectionType, query, startTime, endTime, rollup)
	var firstWeek string
	var currentClicks int64
	var lastClicks int64
	lastClicks = 0
	currentClicks = 0
	for _, weeklyRow := range timeSeriesData {
		currentClicks = weeklyRow.Clicks
		if currentClicks < lastClicks {
			currentClicks = lastClicks
		} else {
			lastClicks = currentClicks
		}
		weeklyMetrics = append(weeklyMetrics, domain.CollectionTimeSeries{
			StartOfWeek: weeklyRow.Startofweek.Format("2006-01-02"),
			Date:        weeklyRow.Startofweek.Format("2006-01-02"),
			Clicks:      weeklyRow.Clicks,
		})
		if firstWeek == "" {
			rollupHours := 7 * 24 * time.Hour
			if rollup == "daily" {
				rollupHours = 1 * 24 * time.Hour
			}
			firstWeek = weeklyRow.Startofweek.Add(-1 * rollupHours).Format("2006-01-02")
		}
	}
	if firstWeek != "" {
		weeklyMetrics = append([]domain.CollectionTimeSeries{{
			StartOfWeek: firstWeek,
			Date:        firstWeek,
			Clicks:      0,
		}}, weeklyMetrics...)
	}
	return weeklyMetrics, nil
}

func (m *CollectionAnalyticsManager) FetchCollectionMetricsSummary(ctx context.Context, query coredomain.SearchQuery) (*domain.CollectionMetricsSummaryEntry, error) {
	collectionIdFilter := rest.GetFilterValueForKey(query.Filters, "collectionId")
	collectionType := rest.GetFilterValueForKey(query.Filters, "collectionType")
	shareId := rest.GetFilterValueForKey(query.Filters, "shareId")
	collectionGroupId := rest.GetFilterValueForKey(query.Filters, "collectionGroupId")
	collectionGroupShareId := rest.GetFilterValueForKey(query.Filters, "collectionGroupShareId")
	isGroupCollection := false
	var collectionIds []string
	if collectionType == nil {
		return nil, errors.New("invalid collection type")
	}
	metadata := domain2.CollectionGroupMetadata{}
	var err error
	if shareId != nil && collectionIdFilter == nil {
		collectionIdFilter, err = m.getCollectionIdFromShareId(ctx, *shareId, *collectionType)
		if err != nil {
			return nil, err
		}
	}
	if collectionGroupId != nil && collectionIdFilter == nil {
		groupId, _ := strconv.ParseInt(*collectionGroupId, 10, 64)
		collectionGroup, err := m.collectionGroupManager.FindById(ctx, groupId)
		if err != nil {
			return nil, err
		}
		if collectionGroup == nil || collectionGroup.CollectionIds == nil {
			return nil, errors.New("collection group is empty")
		}
		isGroupCollection = true
		if collectionGroup.Metadata != nil {
			metadata = *collectionGroup.Metadata
		}
		collectionIds = *collectionGroup.CollectionIds
	}
	if collectionGroupId != nil && collectionIdFilter == nil {
		groupId, _ := strconv.ParseInt(*collectionGroupId, 10, 64)
		collectionIds, err = m.GetCollectionIdFromGroupId(ctx, groupId)
		if err != nil {
			return nil, err
		}
	}
	if collectionGroupShareId != nil && collectionIdFilter == nil {
		collectionIds, err = m.GetCollectionIdFromGroupShareId(ctx, *collectionGroupShareId)
		if err != nil {
			return nil, err
		}
	}
	if collectionIdFilter != nil {
		collectionIds = append(collectionIds, *collectionIdFilter)
	}

	if collectionType == nil || collectionIds == nil || len(collectionIds) == 0 {
		return nil, errors.New("invalid collection")
	}

	var startTime *time.Time
	var endTime *time.Time
	disabledMetrics := []string{}
	budget := 0
	var itemCount int64
	var latestItemAddedOn *time.Time
	itemCount = -1
	if *collectionType == "POST" {
		postCollection, _ := m.postCollectionManager.FindByIdInternal(ctx, collectionIds[0])

		if postCollection != nil {
			itemCount, _ = m.postCollectionItemManager.CountCollectionItemsInACollection(ctx, *postCollection.Id)
			if itemCount > 0 {
				latestItemAddedOn, _ = m.postCollectionItemManager.GetLatestItemAddedOn(ctx, *postCollection.Id)
			}
			pcbud := (*postCollection).Budget
			budget = *pcbud
			if len(collectionIds) > 0 {
				budget = 0
				for _, collectionId := range collectionIds {
					coll, _ := m.postCollectionManager.FindByIdInternal(ctx, collectionId)
					if coll != nil {
						_budget := (*coll).Budget
						if _budget != nil {
							budget += *_budget
						}
					}
				}
			}

			disabledMetrics = (*postCollection).DisabledMetrics // Backward Compat
			if !isGroupCollection {                             // Derive this info from group's meta
				metadata.DisabledMetrics = (*postCollection).DisabledMetrics
				metadata.ExpectedMetricValues = (*postCollection).ExpectedMetricValues
				metadata.Comments = (*postCollection).Comments
			}
			if (*postCollection).StartTime != nil {
				st := time.Unix(*postCollection.StartTime, 0)
				startTime = &st
			}
			if (*postCollection).EndTime != nil {
				et := time.Unix(*postCollection.EndTime, 0)
				endTime = &et
			}
		}
	} else if *collectionType == "PROFILE" {
		id, _ := strconv.ParseInt(collectionIds[0], 10, 64)
		profileCollection, _ := m.profileCollectionManager.FindByIdInternal(ctx, id)
		if profileCollection != nil {
			disabledMetrics = profileCollection.DisabledMetrics // Keep around for backward compat
			if !isGroupCollection {                             // Derive this info from group's meta
				metadata.DisabledMetrics = profileCollection.DisabledMetrics
			}
		}
	} else {
		return nil, errors.New("invalid collection type")
	}

	metricsSummary := m.summaryDAO.FetchCollectionMetricsSummary(ctx, collectionIds, *collectionType, query)
	postWiseCountsData := m.summaryDAO.FetchCollectionPostWiseCounts(ctx, collectionIds, *collectionType, query)
	metricsSummaryOverall := m.summaryDAO.FetchCollectionMetricsSummary(ctx, collectionIds, *collectionType, coredomain.SearchQuery{})
	notice := ""
	noticeType := ""
	if itemCount > 0 {
		diff := itemCount - metricsSummaryOverall.Totalposts
		if diff > 0 {
			if latestItemAddedOn != nil && latestItemAddedOn.Add(time.Hour*2).Before(time.Now()) {
				notice = "Coult not update data for posts added"
				noticeType = "ERROR"
			} else {
				noticeType = "IN_PROGRESS"
				if diff == 1 {
					notice = fmt.Sprintf("%d post is Pending....", diff)
				} else if diff > 1 {
					notice = fmt.Sprintf("%d posts are Pending....", diff)
				}
			}
		}
	}
	summaryEntry := domain.CollectionMetricsSummaryEntry{
		CollectionId:               collectionIds[0],
		CollectionType:             *collectionType,
		Notice:                     notice,
		NoticeType:                 noticeType,
		TotalFollowers:             metricsSummary.Followers,
		TotalProfiles:              metricsSummary.Profilescount,
		TotalPosts:                 metricsSummary.Totalposts,
		Reach:                      metricsSummary.Reach,
		Plays:                      metricsSummary.Plays,
		Impressions:                metricsSummary.Impressions,
		Views:                      metricsSummary.Views,
		Likes:                      metricsSummary.Likes,
		Comments:                   metricsSummary.Comments,
		SwipeUps:                   metricsSummary.SwipeUps,
		StickerTaps:                metricsSummary.StickerTaps,
		Shares:                     metricsSummary.Shares,
		Saves:                      metricsSummary.Saves,
		Mentions:                   metricsSummary.Mentions,
		TotalEngagement:            metricsSummary.Totalengagement,
		StoryExits:                 metricsSummary.Storyexits,
		StoryBackTaps:              metricsSummary.Storybacktaps,
		StoryFwdTaps:               metricsSummary.Storyfwdtaps,
		EngagementRate:             fmt.Sprintf("%.2f", metricsSummary.Er*100),
		ReachPerc:                  "0.0",
		ViewsPerc:                  "0.0",
		ImpressionsPerc:            "0.0",
		LinkClicks:                 metricsSummary.Linkclicks,
		Orders:                     metricsSummary.Orders,
		DeliveredOrders:            metricsSummary.Deliveredorders,
		CompletedOrders:            metricsSummary.Completedorders,
		LeaderboardOverallOrders:   metricsSummary.Leaderboardoverallorders,
		LeaderboardDeliveredOrders: metricsSummary.LeaderboarddeliveredOorders,
		LeaderboardCompletedOrders: metricsSummary.Leaderboardcompletedorders,
		Budget:                     budget,
		RefreshedAt:                metricsSummary.Lastrefreshtime.Format("2006-01-02 12:00"),
		DisabledMetrics:            disabledMetrics,
		Metadata:                   metadata,
	}

	if metricsSummary.Followers > 0 {
		summaryEntry.ReachPerc = fmt.Sprintf("%.2f", float64(summaryEntry.Reach)*100/float64(metricsSummary.Followers))
		summaryEntry.ViewsPerc = fmt.Sprintf("%.2f", float64(summaryEntry.Views)*100/float64(metricsSummary.Followers))
		summaryEntry.ImpressionsPerc = fmt.Sprintf("%.2f", float64(summaryEntry.Impressions)*100/float64(metricsSummary.Followers))
	}

	if *collectionType == "POST" {
		// do everything based on budget and cost to calculate derived metrics
		if summaryEntry.Reach > 0 {
			summaryEntry.ClickThroughRate = fmt.Sprintf("%.2f", float64(summaryEntry.LinkClicks*100)/float64(summaryEntry.Reach))
			summaryEntry.CostPerReach = fmt.Sprintf("%.2f", float64(summaryEntry.Budget)/float64(summaryEntry.Reach))
		}
		if summaryEntry.Likes > 0 || summaryEntry.Comments > 0 {
			summaryEntry.ClickPerEngagement = fmt.Sprintf("%.2f", float64(summaryEntry.LinkClicks)/float64(summaryEntry.Likes+summaryEntry.Comments))
			summaryEntry.CostPerEngagement = fmt.Sprintf("%.2f", float64(summaryEntry.Budget)/float64(summaryEntry.Likes+summaryEntry.Comments))
		}
		if summaryEntry.Impressions > 0 {
			summaryEntry.CostPer1KImpressions = fmt.Sprintf("%.2f", float64(summaryEntry.Budget*1000)/float64(summaryEntry.Impressions))
		}
		if summaryEntry.Views > 0 {
			summaryEntry.CostPerView = fmt.Sprintf("%.2f", float64(summaryEntry.Budget)/float64(summaryEntry.Views))
		}
		if summaryEntry.LinkClicks > 0 {
			summaryEntry.CostPerClick = fmt.Sprintf("%.2f", float64(summaryEntry.Budget)/float64(summaryEntry.LinkClicks))
		}
		if summaryEntry.Orders > 0 {
			summaryEntry.CostPerOrder = fmt.Sprintf("%.2f", float64(summaryEntry.Budget)/float64(summaryEntry.Orders))
		}
	}

	summaryEntry.PostWiseCounts = m.ConvertToPostCountsSummary(postWiseCountsData)

	metricsChan := make(chan []domain.CollectionTimeSeries)
	go m.fetchWeeklyMetricsAsync(ctx, query, collectionIds, collectionType, startTime, endTime, metricsChan)
	select {
	case weeklyMetrics := <-metricsChan:
		summaryEntry.WeeklyMetrics = weeklyMetrics
	case <-time.After(time.Second * 15):
		log.Error("Time Series Metric Fetching Timed Out")
	}
	return &summaryEntry, nil
}

func (m *CollectionAnalyticsManager) fetchWeeklyMetricsAsync(ctx context.Context, query coredomain.SearchQuery, collectionIds []string, collectionType *string, startTime *time.Time, endTime *time.Time, metricsChan chan<- []domain.CollectionTimeSeries) error {
	metricsChan <- m.fetchWeeklyMetrics(ctx, query, collectionIds, collectionType, startTime, endTime)
	return nil
}

func (m *CollectionAnalyticsManager) fetchWeeklyMetrics(ctx context.Context, query coredomain.SearchQuery, collectionIds []string, collectionType *string, startTime *time.Time, endTime *time.Time) []domain.CollectionTimeSeries {
	var weeklyMetrics []domain.CollectionTimeSeries
	weeklyMetricsData := m.tsDAO.FetchCollectionWeeklyMetricsSummary(ctx, collectionIds, *collectionType, query, startTime, endTime, "weekly")
	var firstWeek string
	for _, weeklyRow := range weeklyMetricsData {
		weeklyMetrics = append(weeklyMetrics, domain.CollectionTimeSeries{
			StartOfWeek:     weeklyRow.Startofweek.Format("2006-01-02"),
			TotalPosts:      weeklyRow.Totalposts,
			Views:           int64(weeklyRow.Views),
			Likes:           weeklyRow.Likes,
			Comments:        weeklyRow.Comments,
			Reach:           int64(weeklyRow.Reach),
			TotalEngagement: weeklyRow.Likes + weeklyRow.Comments,
			Impressions:     int64(weeklyRow.Impressions),
		})
		if firstWeek == "" {
			firstWeek = weeklyRow.Startofweek.Add(time.Hour * 24 * 7 * -1).Format("2006-01-02")
		}
	}
	if firstWeek != "" {
		weeklyMetrics = append([]domain.CollectionTimeSeries{{
			StartOfWeek:     firstWeek,
			TotalPosts:      0,
			Views:           0,
			Likes:           0,
			Comments:        0,
			Reach:           0,
			Impressions:     0,
			TotalEngagement: 0,
		}}, weeklyMetrics...)
	}
	return weeklyMetrics
}

func (m *CollectionAnalyticsManager) FetchCollectionProfilesWithMetricsSummary(ctx context.Context, query coredomain.SearchQuery, sortBy string, sortDir string, page int, size int) ([]domain.CollectionProfileMetricsSummaryEntry, error) {

	collectionIdFilter := rest.GetFilterValueForKey(query.Filters, "collectionId")
	collectionType := rest.GetFilterValueForKey(query.Filters, "collectionType")
	shareId := rest.GetFilterValueForKey(query.Filters, "shareId")
	collectionGroupId := rest.GetFilterValueForKey(query.Filters, "collectionGroupId")
	collectionGroupShareId := rest.GetFilterValueForKey(query.Filters, "collectionGroupShareId")
	var collectionIds []string
	if collectionType == nil {
		return nil, errors.New("invalid collection type")
	}
	var err error
	if shareId != nil && collectionIdFilter == nil {
		collectionIdFilter, err = m.getCollectionIdFromShareId(ctx, *shareId, *collectionType)
		if err != nil {
			return nil, err
		}
	}
	if collectionGroupId != nil && collectionIdFilter == nil {
		groupId, _ := strconv.ParseInt(*collectionGroupId, 10, 64)
		collectionIds, err = m.GetCollectionIdFromGroupId(ctx, groupId)
		if err != nil {
			return nil, err
		}
	}
	if collectionGroupId != nil && collectionIdFilter == nil {
		groupId, _ := strconv.ParseInt(*collectionGroupId, 10, 64)
		collectionIds, err = m.GetCollectionIdFromGroupId(ctx, groupId)
		if err != nil {
			return nil, err
		}
	}
	if collectionGroupShareId != nil && collectionIdFilter == nil {
		collectionIds, err = m.GetCollectionIdFromGroupShareId(ctx, *collectionGroupShareId)
		if err != nil {
			return nil, err
		}
	}
	if collectionIdFilter != nil {
		collectionIds = append(collectionIds, *collectionIdFilter)
	}

	if collectionType == nil || collectionIds == nil || len(collectionIds) == 0 {
		return nil, errors.New("invalid collection")
	}

	profileMetrics := m.summaryDAO.FetchCollectionProfiles(ctx, collectionIds, *collectionType, query, sortBy, sortDir, page, size)
	profilePostWiseCountsData := m.summaryDAO.FetchCollectionProfilesPostWiseCounts(ctx, collectionIds, *collectionType, query)
	var instagramHandles []string
	var ytChannelIds []string

	for i := range profileMetrics {
		if profileMetrics[i].Platform == string(constants.InstagramPlatform) && profileMetrics[i].ProfileSocialId != nil && *profileMetrics[i].ProfileSocialId != "" && *profileMetrics[i].ProfileSocialId != "MISSING" {
			instagramHandles = append(instagramHandles, *profileMetrics[i].ProfileSocialId)
		}
		if profileMetrics[i].Platform == string(constants.YoutubePlatform) && profileMetrics[i].Handle != "" {
			ytChannelIds = append(ytChannelIds, profileMetrics[i].Handle)
		}
	}

	var instaProfiles, ytProfiles map[string]coredomain.SocialAccount

	if len(instagramHandles) > 0 {
		instaProfiles, err = m.discoveryManager.GetInstagramDetailsByIgIds(ctx, instagramHandles)
		if err != nil {
			log.Error(err)
		}
	}
	if len(ytChannelIds) > 0 {
		ytProfiles, err = m.discoveryManager.GetYoutubeDetailsByChannelIds(ctx, ytChannelIds)
		if err != nil {
			log.Error(err)
		}
	}

	profiles := []domain.CollectionProfileMetricsSummaryEntry{}
	for _, profileMetricsRow := range profileMetrics {
		var profileCode string
		if strings.ToUpper(profileMetricsRow.Platform) == string(constants.InstagramPlatform) {
			profile, ok := instaProfiles[profileMetricsRow.Handle]
			if ok && profile.ProfileCode != nil {
				profileCode = *profile.ProfileCode
			}
		} else if strings.ToUpper(profileMetricsRow.Platform) == string(constants.YoutubePlatform) {
			profile, ok := ytProfiles[profileMetricsRow.Handle]
			if ok && profile.ProfileCode != nil {
				profileCode = *profile.ProfileCode
			}
		}
		entry := domain.CollectionProfileMetricsSummaryEntry{
			Platform:                   strings.ToUpper(profileMetricsRow.Platform),
			ProfileHandle:              profileMetricsRow.Handle,
			ProfileCode:                profileCode,
			ProfileName:                profileMetricsRow.Name,
			ProfilePic:                 profileMetricsRow.Profilepic,
			Followers:                  profileMetricsRow.Followers,
			Cost:                       profileMetricsRow.Cost,
			TotalPosts:                 profileMetricsRow.Totalposts,
			Reach:                      profileMetricsRow.Reach,
			Plays:                      profileMetricsRow.Plays,
			Impressions:                profileMetricsRow.Impressions,
			Views:                      profileMetricsRow.Views,
			Likes:                      profileMetricsRow.Likes,
			Comments:                   profileMetricsRow.Comments,
			SwipeUps:                   profileMetricsRow.SwipeUps,
			StickerTaps:                profileMetricsRow.StickerTaps,
			Shares:                     profileMetricsRow.Shares,
			Saves:                      profileMetricsRow.Saves,
			Mentions:                   profileMetricsRow.Mentions,
			StoryExits:                 profileMetricsRow.Storyexits,
			StoryBackTaps:              profileMetricsRow.Storybacktaps,
			StoryFwdTaps:               profileMetricsRow.Storyfwdtaps,
			TotalEngagement:            profileMetricsRow.Likes + profileMetricsRow.Comments,
			LinkClicks:                 profileMetricsRow.Linkclicks,
			Orders:                     profileMetricsRow.Orders,
			DeliveredOrders:            profileMetricsRow.Deliveredorders,
			CompletedOrders:            profileMetricsRow.Completedorders,
			LeaderboardOverallOrders:   profileMetricsRow.Leaderboardoverallorders,
			LeaderboardDeliveredOrders: profileMetricsRow.LeaderboarddeliveredOorders,
			LeaderboardCompletedOrders: profileMetricsRow.Leaderboardcompletedorders,
			EngagementRate:             fmt.Sprintf("%.2f", profileMetricsRow.Er*100),
		}
		profileKey := entry.Platform + "::" + entry.ProfileHandle
		entry.PostWiseCounts = m.ConvertToPostCountsSummary(profilePostWiseCountsData[profileKey])
		/*profileId, _ := strconv.ParseInt(entry.PlatformAccountCode, 10, 64)
		if entry.Platform == string(constants.InstagramPlatform) {
			igProfileIds = append(igProfileIds, profileId)
		} else if entry.Platform == string(constants.YoutubePlatform) {
			ytProfileIds = append(ytProfileIds, profileId)
		}*/
		if entry.Followers > 0 {
			entry.ReachPerc = fmt.Sprintf("%.2f", float64(entry.Reach)*100/float64(entry.Followers))
			entry.ViewsPerc = fmt.Sprintf("%.2f", float64(entry.Views)*100/float64(entry.Followers))
			entry.ImpressionsPerc = fmt.Sprintf("%.2f", float64(entry.Impressions)*100/float64(entry.Followers))
		}

		if *collectionType == "POST" {
			// do everything based on budget and cost to calculate derived metrics
			if entry.Reach > 0 {
				entry.ClickThroughRate = fmt.Sprintf("%.2f", float64(entry.LinkClicks*100)/float64(entry.Reach))
				entry.CostPerReach = fmt.Sprintf("%.2f", float64(entry.Cost)/float64(entry.Reach))
			}
			if entry.Likes > 0 || entry.Comments > 0 {
				entry.ClickPerEngagement = fmt.Sprintf("%.2f", float64(entry.LinkClicks)/float64(entry.Likes+entry.Comments))
				entry.CostPerEngagement = fmt.Sprintf("%.2f", float64(entry.Cost)/float64(entry.Likes+entry.Comments))
			}
			if entry.Impressions > 0 {
				entry.CostPer1KImpressions = fmt.Sprintf("%.2f", float64(entry.Cost*1000)/float64(entry.Impressions))
			}
			if entry.Views > 0 {
				entry.CostPerView = fmt.Sprintf("%.2f", float64(entry.Cost)/float64(entry.Views))
			}
			if entry.LinkClicks > 0 {
				entry.CostPerClick = fmt.Sprintf("%.2f", float64(entry.Cost)/float64(entry.LinkClicks))
			}
		}
		profiles = append(profiles, entry)
	}
	return profiles, nil
}

func (m *CollectionAnalyticsManager) FetchCollectionPostsWithMetricsSummary(ctx context.Context, query coredomain.SearchQuery, sortBy string, sortDir string, page int, size int) ([]domain.CollectionPostMetricsSummaryEntry, error) {
	collectionIdFilter := rest.GetFilterValueForKey(query.Filters, "collectionId")
	collectionType := rest.GetFilterValueForKey(query.Filters, "collectionType")
	shareId := rest.GetFilterValueForKey(query.Filters, "shareId")
	collectionGroupId := rest.GetFilterValueForKey(query.Filters, "collectionGroupId")
	collectionGroupShareId := rest.GetFilterValueForKey(query.Filters, "collectionGroupShareId")
	var collectionIds []string
	if collectionType == nil {
		return nil, errors.New("invalid collection type")
	}
	var err error
	if shareId != nil && collectionIdFilter == nil {
		collectionIdFilter, err = m.getCollectionIdFromShareId(ctx, *shareId, *collectionType)
		if err != nil {
			return nil, err
		}
	}
	if collectionGroupId != nil && collectionIdFilter == nil {
		groupId, _ := strconv.ParseInt(*collectionGroupId, 10, 64)
		collectionIds, err = m.GetCollectionIdFromGroupId(ctx, groupId)
		if err != nil {
			return nil, err
		}
	}
	if collectionGroupId != nil && collectionIdFilter == nil {
		groupId, _ := strconv.ParseInt(*collectionGroupId, 10, 64)
		collectionIds, err = m.GetCollectionIdFromGroupId(ctx, groupId)
		if err != nil {
			return nil, err
		}
	}
	if collectionGroupShareId != nil && collectionIdFilter == nil {
		collectionIds, err = m.GetCollectionIdFromGroupShareId(ctx, *collectionGroupShareId)
		if err != nil {
			return nil, err
		}
	}
	if collectionIdFilter != nil {
		collectionIds = append(collectionIds, *collectionIdFilter)
	}

	if collectionType == nil || collectionIds == nil || len(collectionIds) == 0 {
		return nil, errors.New("invalid collection")
	}

	postMetrics := m.summaryDAO.FetchCollectionPosts(ctx, collectionIds, *collectionType, query, sortBy, sortDir, page, size)
	posts := []domain.CollectionPostMetricsSummaryEntry{}
	for _, postMetricsRow := range postMetrics {
		entry := domain.CollectionPostMetricsSummaryEntry{
			Platform:                   strings.ToUpper(postMetricsRow.Platform),
			PostShortCode:              postMetricsRow.PostShortCode,
			PostType:                   postMetricsRow.PostType,
			PostTitle:                  postMetricsRow.PostTitle,
			PostLink:                   postMetricsRow.PostLink,
			PostThumbnail:              postMetricsRow.PostThumbnail,
			PostCollectionItemId:       postMetricsRow.PostCollectionItemId,
			PublishedAt:                postMetricsRow.PublishedAt,
			ProfileHandle:              postMetricsRow.ProfileHandle,
			ProfileName:                postMetricsRow.ProfileName,
			ProfilePic:                 postMetricsRow.ProfilePic,
			Followers:                  postMetricsRow.Followers,
			Cost:                       postMetricsRow.Cost,
			Reach:                      postMetricsRow.Reach,
			Plays:                      postMetricsRow.Plays,
			Impressions:                postMetricsRow.Impressions,
			Views:                      postMetricsRow.Views,
			Likes:                      postMetricsRow.Likes,
			Comments:                   postMetricsRow.Comments,
			SwipeUps:                   postMetricsRow.SwipeUps,
			StickerTaps:                postMetricsRow.StickerTaps,
			Shares:                     postMetricsRow.Shares,
			Saves:                      postMetricsRow.Saves,
			Mentions:                   postMetricsRow.Mentions,
			StoryExits:                 postMetricsRow.StoryExits,
			StoryBackTaps:              postMetricsRow.StoryBackTaps,
			StoryFwdTaps:               postMetricsRow.StoryFwdTaps,
			TotalEngagement:            postMetricsRow.Likes + postMetricsRow.Comments,
			EngagementRate:             fmt.Sprintf("%.2f", postMetricsRow.ER*100),
			LinkClicks:                 postMetricsRow.LinkClicks,
			Orders:                     postMetricsRow.Orders,
			DeliveredOrders:            postMetricsRow.DeliveredOrders,
			CompletedOrders:            postMetricsRow.CompletedOrders,
			LeaderboardOverallOrders:   postMetricsRow.LeaderboardOverallOrders,
			LeaderboardDeliveredOrders: postMetricsRow.LeaderboardDeliveredOrders,
			LeaderboardCompletedOrders: postMetricsRow.LeaderboardCompletedOrders,
			Bookmarked:                 postMetricsRow.Bookmarked,
		}
		postMetricsRow.Hashtags.AssignTo(&entry.Hashtags)
		/*profileId, _ := strconv.ParseInt(entry.PlatformAccountCode, 10, 64)
		if entry.Platform == string(constants.InstagramPlatform) {
			igProfileIds = append(igProfileIds, profileId)
		} else if entry.Platform == string(constants.YoutubePlatform) {
			ytProfileIds = append(ytProfileIds, profileId)
		}*/
		if entry.Followers > 0 {
			entry.ReachPerc = fmt.Sprintf("%.2f", float64(entry.Reach)*100/float64(entry.Followers))
			entry.ViewsPerc = fmt.Sprintf("%.2f", float64(entry.Views)*100/float64(entry.Followers))
			entry.ImpressionsPerc = fmt.Sprintf("%.2f", float64(entry.Impressions)*100/float64(entry.Followers))
		}

		if *collectionType == "POST" {
			// do everything based on budget and cost to calculate derived metrics
			if entry.Reach > 0 {
				entry.ClickThroughRate = fmt.Sprintf("%.2f", float64(entry.LinkClicks)/float64(entry.Reach))
				entry.CostPerReach = fmt.Sprintf("%.2f", float64(entry.Cost)/float64(entry.Reach))
			}
			if entry.Likes > 0 || entry.Comments > 0 {
				entry.ClickPerEngagement = fmt.Sprintf("%.2f", float64(entry.LinkClicks)/float64(entry.Likes+entry.Comments))
				entry.CostPerEngagement = fmt.Sprintf("%.2f", float64(entry.Cost)/float64(entry.Likes+entry.Comments))
			}
			if entry.Impressions > 0 {
				entry.CostPer1KImpressions = fmt.Sprintf("%.2f", float64(entry.Cost*1000)/float64(entry.Impressions))
			}
			if entry.Views > 0 {
				entry.CostPerView = fmt.Sprintf("%.2f", float64(entry.Cost)/float64(entry.Views))
			}
			if entry.LinkClicks > 0 {
				entry.CostPerClick = fmt.Sprintf("%.2f", float64(entry.Cost)/float64(entry.LinkClicks))
			}
		}
		posts = append(posts, entry)
	}
	return posts, nil
}

func (m *CollectionAnalyticsManager) FetchHashtagsForCollections(ctx context.Context, collectionIds []string, collectionType string) ([]domain.CollectionHashtagEntry, error) {
	hashtagEntries := []domain.CollectionHashtagEntry{}
	entities, err := m.hashtagDAO.FindByCollectionIdAndType(ctx, collectionIds, collectionType)
	if err != nil {
		return hashtagEntries, err
	}

	for _, entity := range entities {
		hashtagEntry := domain.CollectionHashtagEntry{
			HashtagName:    entity.HashtagName,
			TaggedCount:    entity.TaggedCount,
			UGCTaggedCount: entity.UGCTaggedCount,
		}
		hashtagEntries = append(hashtagEntries, hashtagEntry)
	}

	return hashtagEntries, nil
}

func (m *CollectionAnalyticsManager) FetchKeywordsForCollections(ctx context.Context, collectionIds []string, collectionType string) ([]domain.CollectionKeywordEntry, error) {
	keywordEntries := []domain.CollectionKeywordEntry{}
	entities, err := m.keywordDAO.FindByCollectionIdAndType(ctx, collectionIds, collectionType)
	if err != nil {
		return keywordEntries, err
	}

	for _, entity := range entities {
		hashtagEntry := domain.CollectionKeywordEntry{
			KeywordName: entity.KeywordName,
			TaggedCount: entity.TaggedCount,
		}
		keywordEntries = append(keywordEntries, hashtagEntry)
	}

	return keywordEntries, nil
}

func (m *CollectionAnalyticsManager) FetchHashtagsForCollectionShared(ctx context.Context, shareId string, collectionType string) ([]domain.CollectionHashtagEntry, error) {
	collectionId, err := m.getCollectionIdFromShareId(ctx, shareId, collectionType)
	if err != nil {
		return nil, err
	}
	if collectionId == nil {
		return nil, errors.New("collection not found")
	}
	return m.FetchHashtagsForCollections(ctx, []string{*collectionId}, collectionType)
}

func (m *CollectionAnalyticsManager) FetchKeywordsForCollectionShared(ctx context.Context, shareId string, collectionType string) ([]domain.CollectionKeywordEntry, error) {
	collectionId, err := m.getCollectionIdFromShareId(ctx, shareId, collectionType)
	if err != nil {
		return nil, err
	}
	if collectionId == nil {
		return nil, errors.New("collection not found")
	}
	return m.FetchKeywordsForCollections(ctx, []string{*collectionId}, collectionType)
}

func (m *CollectionAnalyticsManager) getCollectionIdFromShareId(ctx context.Context, id string, collectionType string) (*string, error) {
	if collectionType == "POST" {
		collection, err := m.postCollectionManager.FindByShareId(ctx, id)
		if err != nil {
			return nil, err
		}
		if collection == nil {
			return nil, errors.New("collection not found")
		}
		return collection.Id, nil
	} else if collectionType == "PROFILE" {
		collection, err := m.profileCollectionManager.FindByShareId(ctx, id)
		if err != nil {
			return nil, err
		}
		if collection == nil {
			return nil, errors.New("collection not found")
		}
		idStr := strconv.FormatInt(*collection.Id, 10)
		return &idStr, nil
	} else {
		return nil, errors.New("invalid collection type")
	}

}

func (m *CollectionAnalyticsManager) GetCollectionIdFromGroupId(ctx context.Context, collectionGroupId int64) ([]string, error) {
	group, err := m.collectionGroupManager.FindById(ctx, collectionGroupId)
	if err != nil {
		return nil, err
	}
	cids := group.CollectionIds
	if cids != nil {
		return *cids, nil
	}
	return nil, errors.New("empty collection")
}

func (m *CollectionAnalyticsManager) GetCollectionIdFromGroupShareId(ctx context.Context, collectionGroupShareId string) ([]string, error) {
	group, err := m.collectionGroupManager.FindByShareId(ctx, collectionGroupShareId)
	if err != nil {
		return nil, err
	}
	cids := group.CollectionIds
	if cids != nil {
		return *cids, nil
	}
	return nil, errors.New("empty collection")
}

func (m *CollectionAnalyticsManager) FetchCollectionSentiment(ctx context.Context, query coredomain.SearchQuery) (*domain.CollectionSentimentEntry, error) {
	var err error
	var totalCount int64

	collectionId := rest.GetFilterValueForKey(query.Filters, "collectionId")
	collectionType := rest.GetFilterValueForKey(query.Filters, "collectionType")
	if collectionType == nil {
		return nil, errors.New("invalid collection type")
	}
	shareId := rest.GetFilterValueForKey(query.Filters, "shareId")

	if shareId != nil && collectionId == nil {
		collectionId, err = m.getCollectionIdFromShareId(ctx, *shareId, *collectionType)
		if err != nil {
			return nil, err
		}
	}

	if collectionId == nil {
		return nil, errors.New("invalid Collection")
	}

	entries, err := m.postCollectionManager.FindById(ctx, *collectionId, false)
	if err != nil {
		log.Error(err)
	}
	if entries == nil {
		return nil, errors.New("collection not found for the given id")
	}
	if entries.SentimentReportPath == nil || entries.SentimentReportBucket == nil {
		return nil, errors.New("sentiment report not available")
	}

	individualSentimentCounts, err := m.metricsDao.GetSentimentCount(ctx, *entries.SentimentReportBucket, *entries.SentimentReportPath)
	if err != nil {
		return nil, err
	}

	for _, sentimentCount := range individualSentimentCounts {
		totalCount += sentimentCount.Count
	}

	sentimentEntry := &domain.CollectionSentimentEntry{}

	for _, sentimentCount := range individualSentimentCounts {
		sentiment := &domain.Sentiment{
			Type:           sentimentCount.Sentiment,
			PostPercentage: (float64(sentimentCount.Count) / float64(totalCount)) * 100,
		}
		sentimentEntry.Sentiment = append(sentimentEntry.Sentiment, *sentiment)
	}

	return sentimentEntry, nil
}

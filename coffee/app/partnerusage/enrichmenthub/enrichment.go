package enrichmenthub

import (
	collectionanalyticsmanager "coffee/app/collectionanalytics/manager"
	discoverydomain "coffee/app/discovery/domain"
	keywordcollectiondomain "coffee/app/keywordcollection/domain"
	keywordcollectionmanager "coffee/app/keywordcollection/manager"

	"coffee/app/partnerusage/dao"
	postcollectiondomain "coffee/app/postcollection/domain"
	coredomain "coffee/core/domain"

	"encoding/json"

	"context"
)

type PartnerUsageEnricher struct {
	collectionAnalyticsManager *collectionanalyticsmanager.CollectionAnalyticsManager
	keywordcollectionmanager   *keywordcollectionmanager.KeywordCollectionManager
}

func CreatePartnerUsageEnricher(ctx context.Context) *PartnerUsageEnricher {
	collectionanalyticsmanager := collectionanalyticsmanager.CreateCollectionAnalyticsManager(ctx)
	keywordcollectionmanager := keywordcollectionmanager.CreateKeywordCollectionManager(ctx)
	manager := &PartnerUsageEnricher{
		collectionAnalyticsManager: collectionanalyticsmanager,
		keywordcollectionmanager:   keywordcollectionmanager,
	}
	return manager
}

func (e *PartnerUsageEnricher) EnrichCampaignReportResult(ctx context.Context, entity dao.AcitivityTrackerEntity) *postcollectiondomain.CampaignReportActivityMeta {
	var activityMetaEntry postcollectiondomain.CampaignReportActivityMeta
	if entity.Meta != nil {
		json.Unmarshal([]byte(*entity.Meta), &activityMetaEntry)
	}
	if entity.Meta != nil {
		json.Unmarshal([]byte(*entity.Meta), &activityMetaEntry)
		query := coredomain.SearchQuery{Filters: []coredomain.SearchFilter{}}
		query.Filters = append(query.Filters, coredomain.SearchFilter{FilterType: "EQ", Field: "collectionId", Value: activityMetaEntry.Id})
		query.Filters = append(query.Filters, coredomain.SearchFilter{FilterType: "EQ", Field: "collectionType", Value: "POST"})
		metricsSummary, _ := e.collectionAnalyticsManager.FetchCollectionMetricsSummary(ctx, query)
		if metricsSummary != nil {
			activityMetaEntry.Result.PostCount = metricsSummary.TotalPosts
			activityMetaEntry.Result.Reach = metricsSummary.Reach
		}
	}
	return &activityMetaEntry
}

func (e *PartnerUsageEnricher) EnrichTopicResearchResult(ctx context.Context, entity dao.AcitivityTrackerEntity) *keywordcollectiondomain.TopicResearchActivityMeta {
	var activityMetaEntry keywordcollectiondomain.TopicResearchActivityMeta
	if entity.Meta != nil {
		json.Unmarshal([]byte(*entity.Meta), &activityMetaEntry)
	}
	if entity.Meta != nil {
		json.Unmarshal([]byte(*entity.Meta), &activityMetaEntry)
		entry, _ := e.keywordcollectionmanager.FindById(ctx, activityMetaEntry.Id, false)
		if entry != nil {
			entry, _ = e.keywordcollectionmanager.GetCollectionReportPartial(ctx, entry)
			if entry != nil && entry.Report != nil {
				activityMetaEntry.Result.PostCount = entry.Report.TotalPosts
				activityMetaEntry.Result.Reach = entry.Report.TotalReach
			}
		}
	}
	return &activityMetaEntry
}

func (e *PartnerUsageEnricher) EnrichDiscoveryResult(ctx context.Context, entity dao.AcitivityTrackerEntity) *discoverydomain.DiscoveryActivityMeta {
	var activityMetaEntry discoverydomain.DiscoveryActivityMeta
	if entity.Meta != nil {
		json.Unmarshal([]byte(*entity.Meta), &activityMetaEntry)
	}
	return &activityMetaEntry
}

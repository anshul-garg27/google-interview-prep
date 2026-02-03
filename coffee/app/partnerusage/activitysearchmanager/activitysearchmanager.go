package manager2

import (
	"coffee/app/partnerusage/dao"
	"coffee/app/partnerusage/domain"
	"coffee/app/partnerusage/enrichmenthub"
	"coffee/app/partnerusage/manager"

	"coffee/constants"
	coredomain "coffee/core/domain"
	"context"
)

type PartnerUsageManager struct {
	partnerUsageDao        *dao.PartnerUsageDao
	partnerProfileTrackDao *dao.PartnerProfileTrackDao
	activityTrackerDao     *dao.ActivityTrackerDao
	partnerenricher        *enrichmenthub.PartnerUsageEnricher
}

func CreatePartnerUsageManager(ctx context.Context) *PartnerUsageManager {
	partnerUsageDao := dao.CreatePartnerUsageDao(ctx)
	partnerProfileTrackDao := dao.CreatePartnerProfileTrackDao(ctx)
	activityTrackerDao := dao.CreateActivityTrackerDao(ctx)
	partnerenricher := enrichmenthub.CreatePartnerUsageEnricher(ctx)

	manager := &PartnerUsageManager{
		partnerUsageDao:        partnerUsageDao,
		partnerProfileTrackDao: partnerProfileTrackDao,
		activityTrackerDao:     activityTrackerDao,
		partnerenricher:        partnerenricher,
	}
	return manager
}

func (m *PartnerUsageManager) Search(ctx context.Context, query coredomain.SearchQuery, sortBy string, sortDir string, page int, size int) (*[]domain.ActivityInput, int64, error) {
	entities, filteredCount, err := m.activityTrackerDao.Search(ctx, query, sortBy, sortDir, page, size)
	var activities []domain.ActivityInput
	for i := range entities {
		activityType := entities[i].Type

		switch activityType {
		case string(constants.Discovery):
			activityMetaEntry := m.partnerenricher.EnrichDiscoveryResult(ctx, entities[i])
			activityEntry, _ := manager.ToActivityEntry(ctx, entities[i])
			activityEntry.Meta = activityMetaEntry
			activities = append(activities, *activityEntry)

		case string(constants.CampaignReport):
			activityMetaEntry := m.partnerenricher.EnrichCampaignReportResult(ctx, entities[i])
			activityEntry, _ := manager.ToActivityEntry(ctx, entities[i])
			activityEntry.Meta = activityMetaEntry
			activities = append(activities, *activityEntry)

		case string(constants.TopicResearch):
			activityMetaEntry := m.partnerenricher.EnrichTopicResearchResult(ctx, entities[i])
			activityEntry, _ := manager.ToActivityEntry(ctx, entities[i])
			activityEntry.Meta = activityMetaEntry
			activities = append(activities, *activityEntry)
		}
	}
	return &activities, filteredCount, err
}

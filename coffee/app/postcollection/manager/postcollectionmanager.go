package manager

import (
	"coffee/app/postcollection/dao"
	"coffee/app/postcollection/domain"
	"coffee/client/jobtracker"
	"coffee/constants"
	"coffee/core/appcontext"
	coredomain "coffee/core/domain"
	"coffee/core/rest"
	"coffee/helpers"
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type PostCollectionManager struct {
	*rest.Manager[domain.PostCollectionEntry, dao.PostCollectionEntity, string]
	postCollectionItemManager *PostCollectionItemManager
	dao                       *dao.PostCollectionDAO
}

func (m *PostCollectionManager) init(ctx context.Context) {
	newDAO := dao.CreatePostCollectionDAO(ctx)
	m.Manager = rest.NewManager(ctx, newDAO.Dao, m.ToEntry, m.ToEntity)
	m.dao = newDAO
}

func CreatePostCollectionManager(ctx context.Context) *PostCollectionManager {
	manager := &PostCollectionManager{}
	manager.init(ctx)
	manager.postCollectionItemManager = CreatePostCollectionItemManager(ctx)
	return manager
}

func (m *PostCollectionManager) ToEntry(entity *dao.PostCollectionEntity) (*domain.PostCollectionEntry, error) {
	entry := &domain.PostCollectionEntry{
		Id:                    &entity.Id,
		ShareId:               &entity.ShareId,
		Name:                  &entity.Name,
		PartnerId:             &entity.PartnerId,
		Source:                &entity.Source,
		SourceId:              &entity.SourceId,
		Budget:                &entity.Budget,
		Enabled:               entity.Enabled,
		CreatedBy:             &entity.CreatedBy,
		CreatedAt:             helpers.ToInt64(entity.CreatedAt.Unix()),
		UpdatedAt:             helpers.ToInt64(entity.UpdatedAt.Unix()),
		MetricsIngestionFreq:  &entity.MetricsIngestionFreq,
		JobId:                 entity.JobId,
		SentimentReportPath:   entity.SentimentReportPath,
		SentimentReportBucket: entity.SentimentReportBucket,
	}
	if entity.StartTime != nil {
		ts := (*entity.StartTime).Unix()
		entry.StartTime = &ts
	}
	if entity.EndTime != nil {
		ts := (*entity.EndTime).Unix()
		entry.EndTime = &ts
	}
	entry.Comments = &entity.Comments
	entity.DisabledMetrics.AssignTo(&entry.DisabledMetrics)
	entity.ExpectedMetricValues.AssignTo(&entry.ExpectedMetricValues)
	return entry, nil
}

func (m *PostCollectionManager) ToEntity(entry *domain.PostCollectionEntry) (*dao.PostCollectionEntity, error) {
	entity := &dao.PostCollectionEntity{}
	if entry.Id != nil {
		entity.Id = *entry.Id
	}
	if entry.ShareId != nil {
		entity.ShareId = *entry.ShareId
	}
	if entry.PartnerId != nil {
		entity.PartnerId = *entry.PartnerId
	}
	if entry.Source != nil {
		entity.Source = *entry.Source
	}
	if entry.SourceId != nil {
		entity.SourceId = *entry.SourceId
	}
	if entry.Name != nil {
		entity.Name = *entry.Name
	}
	if entry.Enabled != nil {
		entity.Enabled = entry.Enabled
	}
	if entry.StartTime != nil {
		st := time.Unix(*entry.StartTime, 0)
		entity.StartTime = &st
	}
	if entry.EndTime != nil {
		et := time.Unix(*entry.EndTime, 0)
		entity.EndTime = &et
	}
	if entry.Budget != nil {
		entity.Budget = *entry.Budget
	}
	if entry.CreatedBy != nil {
		entity.CreatedBy = *entry.CreatedBy
	}
	if entry.Comments != nil {
		entity.Comments = *entry.Comments
	}
	if entry.DisabledMetrics != nil {
		entity.DisabledMetrics.Set(entry.DisabledMetrics)
	}
	if entry.ExpectedMetricValues != nil {
		entity.ExpectedMetricValues.Set(entry.ExpectedMetricValues)
	}
	if entry.MetricsIngestionFreq != nil {
		entity.MetricsIngestionFreq = *entry.MetricsIngestionFreq
	}
	if entry.JobId != nil {
		entity.JobId = entry.JobId
	}
	if entry.SentimentReportPath != nil {
		entity.SentimentReportPath = entry.SentimentReportPath
	}
	if entry.SentimentReportBucket != nil {
		entity.SentimentReportBucket = entry.SentimentReportBucket
	}
	return entity, nil
}

func (m *PostCollectionManager) Create(ctx context.Context, input *domain.PostCollectionEntry) (*domain.PostCollectionEntry, error) {

	if input.Name == nil || input.Source == nil {
		return nil, errors.New("Bad Request - Missing mandatory fields")
	}

	uuid := uuid.New().String()
	input.Id = &uuid
	if strings.ToUpper(*input.Source) == string(constants.SaasCollection) {
		input.SourceId = &uuid
	}

	createdEntry, err := m.Manager.Create(ctx, input)
	if err != nil {
		// record exists
		appCtx := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)
		appCtx.Session.Close(ctx) // TODO: Figure out a better way to reset session. Need to do this because previous session has been marked as rollback-only
		appCtx.Session = nil
		errMsg := err.Error()
		var entry *domain.PostCollectionEntry
		if strings.Contains(errMsg, "post_collection_share_id_key") {
			entry, _ = m.FindByShareId(ctx, *input.ShareId)
		} else if strings.Contains(errMsg, "post_collection_source_key") {
			entry, _ = m.FindBySourceAndSourceId(ctx, *input.Source, *input.SourceId)
		}
		if entry != nil {
			return entry, nil
		}
		return nil, err
	}
	if createdEntry == nil || createdEntry.Id == nil {
		return nil, errors.New("Failed to create post collection")
	}

	if *input.Source == string(constants.SaasCollection) {
		if input.ItemsAssetId != nil {
			entry, err2 := m.addItemsViaCSV(ctx, input, createdEntry)
			if err2 != nil {
				return entry, err2
			}
		} else if input.Items != nil && len(input.Items) > 0 {
			_, err := m.postCollectionItemManager.AddItemsToCollection(ctx, *createdEntry, input.Items)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, errors.New("No Items to add to collection")
		}
	}
	return createdEntry, err
}

func (m *PostCollectionManager) addItemsViaCSV(ctx context.Context, input *domain.PostCollectionEntry, createdEntry *domain.PostCollectionEntry) (*domain.PostCollectionEntry, error) {
	jobEntry := jobtracker.JobEntry{}
	jobName := "Add Items Collection - " + *input.Name
	jobType := "ADD_ITEM_TO_POST_COLLECTION"
	jobEntry.JobType = &jobType
	jobEntry.JobName = &jobName
	jobEntry.Input = map[string]string{}
	jobEntry.Input["collectionId"] = *createdEntry.Id
	jobEntry.InputAssetInformation = &coredomain.AssetInfo{Id: *input.ItemsAssetId}

	client := jobtracker.New(ctx)
	jobResponse, err := client.CreateJob(jobEntry)
	if err != nil {
		return nil, err
	}
	if jobResponse == nil || jobResponse.Status.Type != "SUCCESS" || len(jobResponse.Data) == 0 || (jobResponse.Data)[0].Id == nil {
		return nil, errors.New("Failed to create job for bulk item import")
	}
	jobId := (jobResponse.Data)[0].Id

	updateCollectionEntry := domain.PostCollectionEntry{JobId: jobId}
	m.Manager.Update(ctx, *createdEntry.Id, &updateCollectionEntry)
	createdEntry.JobId = jobId
	return nil, nil
}

func (m *PostCollectionManager) Update(ctx context.Context, id string, input *domain.PostCollectionEntry) (*domain.PostCollectionEntry, error) {
	updatedEntry, err := m.Manager.Update(ctx, id, input)
	if err != nil {
		return nil, err
	}
	if updatedEntry == nil || updatedEntry.Id == nil {
		return nil, errors.New("Failed to create profile collection")
	}
	entry, err := m.Manager.FindById(ctx, id)
	if err == nil && entry != nil && *entry.Source == string(constants.SaasCollection) {
		if input.ItemsAssetId != nil {
			input.Name = entry.Name
			entry, err2 := m.addItemsViaCSV(ctx, input, entry)
			if err2 != nil {
				return entry, err2
			}
		}
	}
	return updatedEntry, err
}

func (m *PostCollectionManager) FindById(ctx context.Context, collectionId string, fetchTotalItems bool) (*domain.PostCollectionEntry, error) {
	e, err := m.Manager.FindById(ctx, collectionId)
	if err != nil {
		return nil, err
	}
	appCtx := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)
	if e != nil && appCtx != nil && appCtx.PartnerId != nil && *appCtx.PartnerId != -1 && *e.PartnerId != *appCtx.PartnerId {
		return nil, errors.New("You are not authorized to view this collection")
	}
	if fetchTotalItems {
		count, _ := m.postCollectionItemManager.CountCollectionItemsInACollection(ctx, collectionId)
		e.TotalPosts = count
	}
	return e, nil
}

func (m *PostCollectionManager) FindByIdInternal(ctx context.Context, collectionId string) (*domain.PostCollectionEntry, error) {
	e, err := m.Manager.FindById(ctx, collectionId)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (m *PostCollectionManager) FindByShareId(ctx context.Context, shareId string) (*domain.PostCollectionEntry, error) {
	e, err := m.dao.FindByShareId(ctx, shareId)
	if err != nil {
		return nil, err
	}
	return m.ToEntry(e)
}

func (m *PostCollectionManager) FindBySourceAndSourceId(ctx context.Context, source string, sourceId string) (*domain.PostCollectionEntry, error) {
	e, err := m.dao.FindBySourceAndSourceId(ctx, source, source)
	if err != nil {
		return nil, err
	}
	return m.ToEntry(e)
}

func (m *PostCollectionManager) CountCollectionsForPartner(ctx context.Context, partnerId int64) (int64, error) {
	count, err := m.dao.CountCollectionsForPartner(ctx, partnerId)
	if err != nil {
		return 0, err
	}
	return count, nil
}

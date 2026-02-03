package manager

import (
	"coffee/app/profilecollection/dao"
	"coffee/app/profilecollection/domain"
	"coffee/helpers"
	"time"
)

// Entity Methods

func ToEntity(entry *domain.ProfileCollectionEntry) (*dao.ProfileCollectionEntity, error) {
	if entry == nil {
		return nil, nil
	}
	entity := &dao.ProfileCollectionEntity{}
	if entry.Id != nil {
		entity.Id = *entry.Id
	}
	if entry.PartnerId != nil {
		entity.PartnerId = *entry.PartnerId
	}
	if entry.ShareId != nil {
		entity.ShareId = *entry.ShareId
	}
	if entry.Name != nil {
		entity.Name = *entry.Name
	}
	if entry.Source != nil {
		entity.Source = *entry.Source
	}
	if entry.SourceId != nil {
		entity.SourceId = *entry.SourceId
	}
	if entry.Description != nil {
		entity.Description = entry.Description
	}
	if entry.Enabled != nil {
		entity.Enabled = entry.Enabled
	}
	if entry.Featured != nil {
		entity.Featured = entry.Featured
	}
	if entry.AnalyticsEnabled != nil {
		entity.AnalyticsEnabled = entry.AnalyticsEnabled
	}
	if entry.Categories != nil {
		entity.Categories = *entry.Categories
	}
	if entry.Tags != nil {
		entity.Tags = *entry.Tags
	}
	if entry.CreatedBy != nil {
		entity.CreatedBy = *entry.CreatedBy
	}
	if entry.DisabledMetrics != nil {
		entity.DisabledMetrics.Set(entry.DisabledMetrics)
	}
	if entry.CategoryIds != nil {
		entity.CategoryIds.Set(entry.CategoryIds)
	}
	if entry.JobId != nil {
		entity.JobId = entry.JobId
	}
	if entry.CampaignId != nil {
		entity.CampaignId = entry.CampaignId
	}
	if entry.CampaignPlatform != nil {
		entity.CampaignPlatform = entry.CampaignPlatform
	}
	if entry.CustomColumns != nil {
		entity.CustomColumns.Set(entry.CustomColumns)
	}
	if entry.OrderedColumns != nil {
		entity.OrderedColumns.Set(entry.OrderedColumns)
	}
	if entry.ShareEnabled != nil {
		entity.ShareEnabled = entry.ShareEnabled
	}
	if entry.CreatedAt != nil {
		entity.CreatedAt = time.Unix(*entry.CreatedAt, 0)
	}
	if entry.UpdatedAt != nil {
		entity.UpdatedAt = time.Unix(*entry.UpdatedAt, 0)
	}
	return entity, nil
}

func ToItemEntity(entry *domain.ProfileCollectionItemEntry) (*dao.ProfileCollectionItemEntity, error) {
	if entry == nil {
		return nil, nil
	}
	entity := &dao.ProfileCollectionItemEntity{}
	if entry.Platform != nil {
		entity.Platform = *entry.Platform
	}
	if entry.PlatformAccountCode != nil {
		entity.PlatformAccountCode = *entry.PlatformAccountCode
	}
	if entry.CampaignProfileId != nil {
		entity.CampaignProfileId = entry.CampaignProfileId
	}
	if entry.ShortlistingStatus != nil {
		entity.ShortlistingStatus = entry.ShortlistingStatus
	}
	if entry.ShortlistId != nil {
		entity.ShortlistId = entry.ShortlistId
	}
	if entry.ProfileSocialId != nil {
		entity.ProfileSocialId = *entry.ProfileSocialId
	}
	if entry.ProfileCollectionId != nil {
		entity.ProfileCollectionId = *entry.ProfileCollectionId
	}
	if entry.PartnerId != nil {
		entity.PartnerId = *entry.PartnerId
	}
	if entry.Enabled != nil {
		entity.Enabled = entry.Enabled
	}

	if entry.Hidden != nil {
		entity.Hidden = entry.Hidden
	}
	if entry.Rank != nil {
		entity.Rank = *entry.Rank
	}
	if entry.BrandSelectionStatus != nil {
		entity.BrandSelectionStatus = entry.BrandSelectionStatus
	}
	if entry.BrandSelectionRemarks != nil {
		entity.BrandSelectionRemarks = entry.BrandSelectionRemarks
	}
	if entry.CreatedBy != nil {
		entity.CreatedBy = *entry.CreatedBy
	}
	if entry.Commercials != nil {
		entity.InternalCommercials = entry.Commercials
	}
	if entry.CreatedAt != nil {
		entity.CreatedAt = time.Unix(*entry.CreatedAt, 0)
	}
	if entry.UpdatedAt != nil {
		entity.UpdatedAt = time.Unix(*entry.UpdatedAt, 0)
	}

	return entity, nil
}

// Entry Methods

func ToEntry(entity *dao.ProfileCollectionEntity) (*domain.ProfileCollectionEntry, error) {
	if entity == nil {
		return nil, nil
	}
	entry := &domain.ProfileCollectionEntry{
		Id:               &entity.Id,
		Name:             &entity.Name,
		PartnerId:        &entity.PartnerId,
		Source:           &entity.Source,
		SourceId:         &entity.SourceId,
		Description:      entity.Description,
		Featured:         entity.Featured,
		Enabled:          entity.Enabled,
		AnalyticsEnabled: entity.AnalyticsEnabled,
		ShareEnabled:     entity.ShareEnabled,
		CreatedBy:        &entity.CreatedBy,
		ShareId:          &entity.ShareId,
		Categories:       (*[]string)(&entity.Categories),
		Tags:             (*[]string)(&entity.Tags),
		CampaignId:       entity.CampaignId,
		CampaignPlatform: entity.CampaignPlatform,
		JobId:            entity.JobId,
		CreatedAt:        helpers.ToInt64(entity.CreatedAt.Unix()),
		UpdatedAt:        helpers.ToInt64(entity.UpdatedAt.Unix()),
	}
	entity.DisabledMetrics.AssignTo(&entry.DisabledMetrics)
	entity.CustomColumns.AssignTo(&entry.CustomColumns)
	entity.OrderedColumns.AssignTo(&entry.OrderedColumns)
	entity.CategoryIds.AssignTo(&entry.CategoryIds)

	return entry, nil
}

func ToItemEntry(entity *dao.ProfileCollectionItemEntity) (*domain.ProfileCollectionItemEntry, error) {
	if entity == nil {
		return nil, nil
	}
	entry := &domain.ProfileCollectionItemEntry{
		Id:                    &entity.Id,
		Platform:              &entity.Platform,
		PlatformAccountCode:   &entity.PlatformAccountCode,
		CampaignProfileId:     entity.CampaignProfileId,
		ProfileSocialId:       &entity.ProfileSocialId,
		ProfileCollectionId:   &entity.ProfileCollectionId,
		PartnerId:             &entity.PartnerId,
		ShortlistingStatus:    entity.ShortlistingStatus,
		BrandSelectionStatus:  entity.BrandSelectionStatus,
		BrandSelectionRemarks: entity.BrandSelectionRemarks,
		Enabled:               entity.Enabled,
		Hidden:                entity.Hidden,
		Rank:                  &entity.Rank,
		CreatedBy:             &entity.CreatedBy,
		ShortlistId:           entity.ShortlistId,
		Commercials:           entity.InternalCommercials,
		CreatedAt:             helpers.ToInt64(entity.CreatedAt.Unix()),
		UpdatedAt:             helpers.ToInt64(entity.UpdatedAt.Unix()),
	}
	return entry, nil
}

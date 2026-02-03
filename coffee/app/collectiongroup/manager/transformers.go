package manager

import (
	"coffee/app/collectiongroup/dao"
	"coffee/app/collectiongroup/domain"
)

func ToEntity(entry *domain.CollectionGroupEntry) (*dao.CollectionGroupEntity, error) {
	if entry == nil {
		return nil, nil
	}
	entity := &dao.CollectionGroupEntity{}
	if entry.PartnerId != nil {
		entity.PartnerId = *entry.PartnerId
	}
	if entry.ShareId != nil {
		entity.ShareId = *entry.ShareId
	}
	if entry.Objective != nil {
		entity.Objective = *entry.Objective
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
	if entry.Enabled != nil {
		entity.Enabled = entry.Enabled
	}
	if entry.CollectionIds != nil {
		entity.CollectionIds = *entry.CollectionIds
	}
	if entry.CreatedBy != nil {
		entity.CreatedBy = *entry.CreatedBy
	}
	if entry.Metadata != nil {
		entity.Metadata.Set(entry.Metadata)
	}
	return entity, nil
}

func ToEntry(entity *dao.CollectionGroupEntity) (*domain.CollectionGroupEntry, error) {
	if entity == nil {
		return nil, nil
	}
	entry := &domain.CollectionGroupEntry{
		Id:            &entity.Id,
		Name:          &entity.Name,
		PartnerId:     &entity.PartnerId,
		Source:        &entity.Source,
		SourceId:      &entity.SourceId,
		CollectionIds: (*[]string)(&entity.CollectionIds),
		Enabled:       entity.Enabled,
		Objective:     &entity.Objective,
		CreatedBy:     &entity.CreatedBy,
		ShareId:       &entity.ShareId,
	}
	entity.Metadata.AssignTo(&entry.Metadata)
	return entry, nil
}

package api

import (
	"coffee/app/collectiongroup/dao"
	"coffee/app/collectiongroup/domain"
	"coffee/app/collectiongroup/manager"
	coredomain "coffee/core/domain"
	"coffee/core/rest"
	"context"
)

// common helper methods

func getPrimaryDataStore() coredomain.DataStore {
	return coredomain.Postgres
}

// -------------------------------
// profile collection service
// -------------------------------

type Service struct {
	*rest.Service[domain.CollectionGroupResponse, domain.CollectionGroupEntry, dao.CollectionGroupEntity, int64]
	Manager *manager.Manager
}

func (s *Service) Init(ctx context.Context) {
	manager := manager.CreateManager(ctx)
	service := rest.NewService(ctx, manager.Manager, getPrimaryDataStore, domain.CreateResponse, domain.CreateErrorResponse)
	s.Service = service
	s.Manager = manager
}

func CreateService(ctx context.Context) *Service {
	service := &Service{}
	service.Init(ctx)
	return service
}

func NewCollectionGroupService(ctx context.Context) *Service {
	service := CreateService(ctx)
	return service
}

func (s *Service) Create(ctx context.Context, entry *domain.CollectionGroupEntry) domain.CollectionGroupResponse {
	e, err := s.Manager.Create(ctx, entry)
	if err != nil {
		return domain.CreateErrorResponse(err, 0)
	}
	return domain.CreateResponse([]domain.CollectionGroupEntry{*e}, 1, "", "Record Created Successfully")
}

func (s *Service) FindById(ctx context.Context, id int64) domain.CollectionGroupResponse {
	entry, err := s.Manager.FindById(ctx, id)
	if err != nil {
		return domain.CreateErrorResponse(err, 0)
	}
	return domain.CreateResponse([]domain.CollectionGroupEntry{*entry}, 1, "", "Record Retrieved Successfully")
}

func (s *Service) FindByShareId(ctx context.Context, shareId string) domain.CollectionGroupResponse {
	collection, err := s.Manager.FindByShareId(ctx, shareId)
	count := int64(1)
	collections := []domain.CollectionGroupEntry{}
	if collection != nil {
		collections = append(collections, *collection)
	}
	if err != nil {
		return domain.CreateErrorResponse(err, 0)
	}
	return domain.CreateResponse(collections, count, "", "Record Retrieved Successfully")
}

func (s *Service) RenewShareId(ctx context.Context, id int64) domain.CollectionGroupResponse {
	updatedEntry, err := s.Manager.RenewShareId(ctx, id)
	if err != nil {
		return domain.CreateErrorResponse(err, 0)
	}
	return domain.CreateResponse([]domain.CollectionGroupEntry{*updatedEntry}, 1, "", "Record Updated Successfully")
}

func (s *Service) Delete(ctx context.Context, id int64) domain.CollectionGroupResponse {
	_, err := s.Manager.Delete(ctx, id)
	if err != nil {
		return domain.CreateErrorResponse(err, 0)
	}
	return domain.CreateResponse(nil, 1, "", "Record Updated Successfully")
}

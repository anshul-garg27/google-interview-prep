package api

import (
	activitysearchmanager "coffee/app/partnerusage/activitysearchmanager"
	"coffee/app/partnerusage/domain"
	"coffee/app/partnerusage/manager"
	"coffee/constants"
	"coffee/core/appcontext"
	coredomain "coffee/core/domain"
	"context"
	"errors"
	"strconv"
)

type Service struct {
	Manager               *manager.PartnerUsageManager
	ActivitySearchManager *activitysearchmanager.PartnerUsageManager
}

func NewPartnerUsageService() *Service {
	service := CreateService(context.Background())
	return service
}

func (s *Service) init(ctx context.Context) {
	manager := manager.CreatePartnerUsageManager(ctx)
	activitysearchmanager := activitysearchmanager.CreatePartnerUsageManager(ctx)
	s.Manager = manager
	s.ActivitySearchManager = activitysearchmanager
}

func CreateService(ctx context.Context) *Service {
	service := &Service{}
	service.init(ctx)
	return service
}

func (s *Service) FindByPartnerId(ctx context.Context, partnerId int64) domain.PartnerResponse {
	partnerusage, err := s.Manager.FindByPartnerId(ctx, partnerId)
	if err != nil {
		return domain.CreatePartnerErrorResponse(err, 500)
	}
	return domain.CreatePartnerResponse([]domain.PartnerInput{*partnerusage}, int64(1), "", "Partner Contract successfully retrieved")
}

func (s *Service) UpdatePartnerUsage(ctx context.Context, partnerId int64, partnerContract *domain.PartnerInput) domain.PartnerResponse {
	var err error

	if len(partnerContract.Contracts) == 0 {
		return domain.CreatePartnerErrorResponse(errors.New("contract not found"), 500)
	}
	if partnerContract.Contracts[0].Plan == string(constants.PaidPlan) || partnerContract.Contracts[0].Plan == string(constants.SaasPlan) {
		err = s.Manager.UpdatePartnerUsage(ctx, partnerId, partnerContract)

	} else if partnerContract.Contracts[0].Plan == string(constants.FreePlan) {
		err = s.Manager.DeletePartnerContract(ctx, partnerId)
	}
	if err != nil {
		return domain.CreatePartnerErrorResponse(err, 500)
	}
	return domain.CreatePartnerResponse(nil, int64(1), "", "Partner Contract successfully Updated")
}

func (s *Service) Search(ctx context.Context, query coredomain.SearchQuery, sortBy string, sortDir string, page int, size int) domain.ActivityTrackerResponse {

	appCtx := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)
	accountId := ""
	if appCtx.AccountId != nil {
		accountId = strconv.FormatInt(*appCtx.AccountId, 10)
	}
	query.Filters = append([]coredomain.SearchFilter{{
		FilterType: "EQ",
		Field:      "account_id",
		Value:      accountId,
	}}, query.Filters...)

	entries, filteredCount, err := s.ActivitySearchManager.Search(ctx, query, sortBy, sortDir, page, size)

	if err != nil {
		return domain.CreateActivityTrackerErrorResponse(err, int(filteredCount))
	}
	var nextCursor string
	if len(*entries) >= size {
		nextCursor = strconv.Itoa(page + 1)
	} else {
		nextCursor = ""
	}
	return domain.CreateActivityTrackerResponse(*entries, filteredCount, nextCursor, "Record(s) Retrieved Successfully")

}

func (s *Service) UnlockProfile(ctx context.Context, partnerID *int64, profileCode string, platform string) coredomain.Status {
	err := s.Manager.UnlockProfile(ctx, partnerID, profileCode, platform)
	return domain.CreateUnlockProfileResponse(err)
}

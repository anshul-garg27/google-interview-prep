package api

import (
	"coffee/app/campaignprofiles/domain"
	"coffee/app/campaignprofiles/manager"
	"coffee/constants"
	"coffee/core/appcontext"
	coredomain "coffee/core/domain"
	"context"
	"strconv"
)

type Service struct {
	Manager *manager.CampaignProfileManager
}

func NewCampaignService(ctx context.Context) *Service {
	service := CreateService(ctx)
	return service
}

func (s *Service) init(ctx context.Context) {
	manager := manager.CreateManager(ctx)
	s.Manager = manager

}

func CreateService(ctx context.Context) *Service {
	service := &Service{}
	service.init(ctx)
	return service
}

func (s *Service) UpsertUsingProfileCode(ctx context.Context, profileCode string, cpInput coredomain.CampaignProfileInput, linkedSource string) domain.CampaignProfileResponse {
	appCtx := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)
	accountIdString := ""
	if appCtx.AccountId != nil {
		accountIdString = strconv.Itoa(int(*appCtx.AccountId))
		cpInput.UpdatedBy = &accountIdString
	} else {
		system := "SYSTEM"
		cpInput.UpdatedBy = &system
	}
	campaigns, err := s.Manager.UpsertProfileAdminDetails(ctx, profileCode, cpInput, linkedSource)
	if err != nil {
		return domain.CreateCampaignProfileErrorResponse(err, 0)
	}
	return domain.CreateCampaignProfileResponse(campaigns, int64(1), "", "retrieved successfully")
}

func (s *Service) FindCampaignProfileByPlatformHandle(ctx context.Context, platform string, platformHandle string, linkedSource string, fullRefresh bool) domain.CampaignProfileResponse {
	campaigns, err := s.Manager.FindCampaignProfileByPlatformHandle(ctx, platform, platformHandle, linkedSource, fullRefresh)
	if err != nil {
		return domain.CreateCampaignProfileErrorResponse(err, 0)
	}
	return domain.CreateCampaignProfileResponse([]coredomain.CampaignProfileEntry{*campaigns}, int64(1), "", "retrieved successfully")
}

func (s *Service) FindById(ctx context.Context, id int64, socials bool) domain.CampaignProfileResponse {
	campaigns, err := s.Manager.FindById(ctx, id, socials)
	if err != nil {
		return domain.CreateCampaignProfileErrorResponse(err, 0)
	}
	return domain.CreateCampaignProfileResponse([]coredomain.CampaignProfileEntry{*campaigns}, int64(1), "", "retrieved successfully")
}

func (s *Service) UpdateById(ctx context.Context, id int64, input coredomain.CampaignProfileInput, linkedSource string) domain.CampaignProfileResponse {
	appCtx := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)
	accountIdString := ""
	if appCtx.AccountId != nil {
		accountIdString = strconv.Itoa(int(*appCtx.AccountId))
		input.UpdatedBy = &accountIdString
	} else {
		system := "SYSTEM"
		input.UpdatedBy = &system
	}
	campaigns, err := s.Manager.UpdateById(ctx, id, input, linkedSource)
	if err != nil {
		return domain.CreateCampaignProfileErrorResponse(err, 0)
	}
	return domain.CreateCampaignProfileResponse(campaigns, int64(1), "", "retrieved successfully")
}

func (s *Service) Search(ctx context.Context, query coredomain.SearchQuery, sortBy string, sortDir string, page int, size int, fetchSocials bool) domain.CampaignProfileResponse {
	entries, filteredCount, err := s.Manager.Search(ctx, query, sortBy, sortDir, page, size, fetchSocials)

	if err != nil {
		return domain.CreateCampaignProfileErrorResponse(err, int(filteredCount))
	}
	var nextCursor string
	if len(entries) >= size {
		nextCursor = strconv.Itoa(page + 1)
	} else {
		nextCursor = ""
	}
	return domain.CreateCampaignProfileResponse(entries, filteredCount, nextCursor, "Record(s) Retrieved Successfully")

}

func (s *Service) UpsertByHandle(ctx context.Context, input coredomain.CampaignProfileInput, linkedSource string) domain.CampaignProfileResponse {
	appCtx := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)
	accountIdString := ""
	if appCtx.AccountId != nil {
		accountIdString = strconv.Itoa(int(*appCtx.AccountId))
		input.UpdatedBy = &accountIdString
	} else {
		system := "SYSTEM"
		input.UpdatedBy = &system
	}
	campaigns, err := s.Manager.UpsertProfileAdminDetails(ctx, "", input, linkedSource)
	if err != nil {
		return domain.CreateCampaignProfileErrorResponse(err, 0)
	}
	return domain.CreateCampaignProfileResponse(campaigns, int64(1), "", "retrieved successfully")
}
func (s Service) RefreshSocialProfileById(ctx context.Context, id int64) domain.CampaignProfileResponse {
	campaignProfile, err := s.Manager.RefreshSocialProfileById(ctx, id)
	if err != nil || campaignProfile == nil {
		return domain.CreateCampaignProfileErrorResponse(err, 0)
	}
	return domain.CreateCampaignProfileResponse([]coredomain.CampaignProfileEntry{*campaignProfile}, int64(1), "", "retrieved successfully")

}

func (s Service) CreatorInsightsById(ctx context.Context, id int64, token string, platform string, userId string) domain.CampaignProfileResponse {
	campaignProfile, err := s.Manager.CreatorInsightsById(ctx, id, token, platform, userId)
	if err != nil || campaignProfile == nil {
		return domain.CreateCampaignProfileErrorResponse(err, 0)
	}
	return domain.CreateCampaignProfileResponse([]coredomain.CampaignProfileEntry{*campaignProfile}, int64(1), "", "retrieved successfully")

}

func (s Service) CreatorAudienceInsightsDataById(ctx context.Context, id int64, token string, platform string, userId string) domain.CampaignProfileAudienceInsightsResponse {
	campaignProfile, err := s.Manager.CreatorAudienceInsightsDataById(ctx, id, token, platform, userId)
	if err != nil || campaignProfile == nil {
		return domain.CreateCampaignProfileAudienceInsightsErrorResponse(err, 0)
	}
	return domain.CreateCampaignProfileAudienceInsightsResponse(campaignProfile, int64(1), "", "retrieved successfully")

}

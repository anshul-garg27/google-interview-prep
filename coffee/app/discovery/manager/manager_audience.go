package manager

import (
	"coffee/app/discovery/dao"
	"coffee/app/discovery/domain"
	"context"
)

type AudienceManager struct {
	socialProfileAudienceInfoDao *dao.SocialProfileAudienceInfoDao
	searchManager                *SearchManager
}

func CreateAudienceManager(ctx context.Context) *AudienceManager {
	socialProfileAudienceInfoDao := dao.CreateSocialProfileAudienceInfoDao(ctx)
	searchManager := CreateSearchManager(ctx)
	managerAudience := &AudienceManager{
		socialProfileAudienceInfoDao: socialProfileAudienceInfoDao,
		searchManager:                searchManager,
	}
	return managerAudience
}

func (m *AudienceManager) FindAudienceByPlatformProfileId(ctx context.Context, platform string, platformProfileId int64) (*domain.SocialProfileAudienceInfoEntry, error) {
	var audienceEntry *domain.SocialProfileAudienceInfoEntry
	audienceEntity, err := m.socialProfileAudienceInfoDao.FindByPlatformProfileId(ctx, platform, platformProfileId)
	if err != nil {
		return audienceEntry, err
	}
	idArray := ([]string)(audienceEntity.NotableFollowers)
	notableFollowers := m.searchManager.FindNotableFollowersDetails(ctx, idArray, platform)
	audienceEntry, _ = ToSocialProfileAudienceInfoEntry(*audienceEntity)
	if notableFollowers != nil {
		audienceEntry.NotableFollowers = *notableFollowers
	}

	return audienceEntry, nil
}

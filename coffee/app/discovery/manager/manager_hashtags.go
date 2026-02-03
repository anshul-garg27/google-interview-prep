package manager

import (
	"coffee/app/discovery/dao"
	"coffee/app/discovery/domain"
	"context"
)

type HashtagsManager struct {
	socialProfileHashtagsDao *dao.SocialProfileHashtagsDao
}

func CreateHashtagsManager(ctx context.Context) *HashtagsManager {
	socialProfileHashtagsDao := dao.CreateSocialProfileHashtagssDao(ctx)
	managerHashtags := &HashtagsManager{
		socialProfileHashtagsDao: socialProfileHashtagsDao,
	}
	return managerHashtags
}

func (m *HashtagsManager) FindHashTagsByPlatformProfileId(ctx context.Context, platform string, platformProfileId int64) (*domain.SocialProfileHashatagsEntry, error) {
	var err error
	hashtags, err := m.socialProfileHashtagsDao.FindByPlatformProfileId(ctx, platform, platformProfileId)
	if err != nil {
		return nil, err
	}
	hashtag, _ := ToSocialProfileHashtagsEntry(*hashtags)
	return hashtag, nil
}

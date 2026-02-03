package leaderboard

import (
	coredomain "coffee/core/domain"
	"coffee/core/rest"
	"context"
	"strconv"
)

type Service struct {
	Manager *Manager
}

func NewLeaderboardService() *Service {
	ctx := context.TODO()
	manager := CreateManager(ctx)
	return &Service{
		Manager: manager,
	}
}

func (s *Service) LeaderBoardByPlatform(ctx context.Context, platform string, searchQuery coredomain.SearchQuery, sortBy string, sortDir string, page int, size int, source string) LeaderboardResponse {

	var profiles *[]LeaderboardProfile
	var filteredCount int64
	var err error
	var nextCursor string

	err = validateRequest(ctx, page, size)
	if err != nil {
		return s.CreateLeaderboardResponse(err, profiles, filteredCount, nextCursor)
	}
	searchQuery = rest.ParseGenderAgeFilters(searchQuery)
	profiles, filteredCount, err = s.Manager.LeaderBoardByPlatform(ctx, platform, searchQuery, sortBy, sortDir, page, size, source)

	if filteredCount > int64(size)*int64(page) {
		nextCursor = strconv.Itoa(page + 1)
	} else {
		nextCursor = ""
	}
	return s.CreateLeaderboardResponse(err, profiles, filteredCount, nextCursor)
}

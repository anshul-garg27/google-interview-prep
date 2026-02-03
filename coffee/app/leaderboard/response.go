package leaderboard

import (
	coredomain "coffee/core/domain"
)

type LeaderboardResponse struct {
	Status   coredomain.Status     `json:"status"`
	Profiles *[]LeaderboardProfile `json:"profiles"`
}

func (s *Service) CreateLeaderboardResponse(err error, profiles *[]LeaderboardProfile, filterCount int64, nextCursor string) LeaderboardResponse {
	HasNextPage := false
	if nextCursor != "" {
		HasNextPage = true
	}
	if err != nil {
		return LeaderboardResponse{
			Status:   coredomain.Status{Status: "500", Message: err.Error(), FilteredCount: filterCount},
			Profiles: &[]LeaderboardProfile{},
		}
	}
	if profiles == nil {
		return LeaderboardResponse{
			Status: coredomain.Status{Status: "500", Message: err.Error(), FilteredCount: filterCount},
		}

	}
	return LeaderboardResponse{
		Status: coredomain.Status{
			Status:        "SUCCESS",
			Message:       "Leaderbaord Retrieved",
			FilteredCount: filterCount,
			HasNextPage:   HasNextPage,
			NextCursor:    nextCursor,
		},
		Profiles: profiles,
	}
}

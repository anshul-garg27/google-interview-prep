package leaderboard

import (
	"bytes"
	"coffee/client/dam"
	"coffee/client/jobtracker"
	"coffee/constants"
	"coffee/core/appcontext"
	coredomain "coffee/core/domain"
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
	log "github.com/sirupsen/logrus"
)

func SetupListeners(router *message.Router, subscriber *amqp.Subscriber) {

	router.AddNoPublisherHandler(
		"download_leaderboard",
		"jobtracker.dx___download_leaderboard_q",
		subscriber,
		PerformJobDownloadLeaderboard,
	)
}

func PerformJobDownloadLeaderboard(message *message.Message) error {
	// This method is used to listen to request for
	// downloading a collection
	ctx := message.Context()
	appCtx := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)
	jobBatch := &jobtracker.JobBatchEntry{}
	if err := json.Unmarshal(message.Payload, jobBatch); err == nil {
		log.Println(fmt.Sprintf("Process Download Leaderboard Job Batch %+v", jobBatch))
		jobBatch.Status = "FAILED"
		jobBatch.Remarks = "Incomplete request"
		if jobBatch != nil && jobBatch.Id > 0 && jobBatch.JobBatchItems != nil && len(jobBatch.JobBatchItems) > 0 {
			jobBatchItem := jobBatch.JobBatchItems[0]
			jobBatchItem.Status = "FAILURE"
			if false && (appCtx.PartnerId == nil || appCtx.AccountId == nil) {
				jobBatch.Remarks = "Unauthorized Access"
			} else {
				jobBatch.Remarks = "Input metadata missing - collectionId/platform"
				if jobBatchItem.InputNode != nil {
					platform, ok1 := jobBatchItem.InputNode["platform"]
					sortBy, ok2 := jobBatchItem.InputNode["sortBy"]
					sortDir, ok3 := jobBatchItem.InputNode["sortDir"]
					if !ok2 {
						sortBy = "followers"
					}
					if !ok3 {
						sortDir = "ASC"
					}
					if ok1 {
						var query coredomain.SearchQuery
						queryInput, ok3 := jobBatchItem.InputNode["query"]
						if ok3 {
							err = json.Unmarshal([]byte(queryInput), &query)
							if err != nil {
								log.Error(err)
							}
						}
						month := time.Date(time.Now().AddDate(0, -1, 0).Year(), time.Now().AddDate(0, -1, 0).Month(), 1, 0, 0, 0, 0, time.Now().AddDate(0, -1, 0).Location()).Format("2006-01-02")
						for i := range query.Filters {
							if query.Filters[i].Field == "month" {
								month = query.Filters[i].Value
							}
						}
						fileName := fmt.Sprintf("Leaderboard-%s-%s.csv", platform, month)
						data, err := DownloadLeaderboard(message.Context(), platform, query, sortBy, sortDir)
						if err == nil {
							// Upload to DAM and update Job with appropriate Asset ID and URL
							damClient := dam.New(message.Context())
							input := &dam.GenerateUrlReqEntry{
								Usage:     "JOB_TRACKER",
								MediaType: "CSV"}
							response, err := damClient.DamUpload(input, data, fileName)
							if err == nil {
								asset := response.List[0]
								if asset.Url != nil {
									jobBatch.Status = "COMPLETED"
									jobBatchItem.Status = "SUCCESS"
									jobBatch.Job.OutputFileAssetInformation = &asset
								}
							} else {
								jobBatch.Remarks = err.Error()
							}
						} else {
							jobBatch.Remarks = err.Error()
						}
					}
				} else {
					jobBatch.Remarks = "Input metadata is missing"
				}
			}

			jobBatch.JobBatchItems[0] = jobBatchItem
			client := jobtracker.New(message.Context())
			client.UpdateJobBatch(*jobBatch)
		}
	}
	return nil
}

func DownloadLeaderboard(ctx context.Context, platform string, query coredomain.SearchQuery, sortBy string, sortDir string) ([]byte, error) {
	manager := CreateManager(ctx)

	leaderboard, _, _ := manager.LeaderBoardByPlatform(ctx, platform, query, sortBy, sortDir, 1, 1000, string(constants.SaasLinkedProfile)) // TODO: Batch this
	if leaderboard == nil {
		return nil, errors.New("invalid leaderboard")
	}
	data := make([][]string, len(*leaderboard)+1)
	allLeaderboardColumns := []string{
		"Rank",
		"Rank Change",
		"Name",
		"Instagram Account",
		"Youtube Channel",
		"Category",
		"Total Followers",
		"Instagram Followers",
		"Youtube Subscribers",
		"Follower Growth",
		"Follower Growth Percentage",
		"Views (30 days)",
	}
	instaLeaderboardColumns := []string{
		"Rank",
		"Rank Change",
		"Name",
		"Instagram Account",
		"Category",
		"Total Followers",
		"Instagram Followers",
		"Follower Growth",
		"Follower Growth Percentage",
		"Engagement Rate",
		"Views (30 days)",
	}
	ytLeaderboardColumns := []string{
		"Rank",
		"Rank Change",
		"Name",
		"Youtube Account",
		"Category",
		"Total Followers",
		"Youtube Followers",
		"Follower Growth",
		"Follower Growth Percentage",
		"Views (30 days)",
		"Lifetime Views",
	}
	if platform == "ALL" {
		data[0] = allLeaderboardColumns
	} else if platform == string(constants.InstagramPlatform) {
		data[0] = instaLeaderboardColumns
	} else if platform == string(constants.YoutubePlatform) {
		data[0] = ytLeaderboardColumns
	}
	for j := range *leaderboard {
		row := (*leaderboard)[j]
		rank := fmt.Sprintf("%d", row.Rank)
		rankChange := "NEW"
		if row.RankChange != nil {
			rankChange = fmt.Sprintf("%d", *row.RankChange)
		}
		name := ""
		instaHandle := ""
		ytUsername := ""
		if row.AccountInfo.Name != nil {
			name = *row.AccountInfo.Name
		}
		category := row.Category
		totalFollowers := int64(0)
		instaFollowers := int64(0)
		ytFollowers := int64(0)
		followerGrowth := row.LeaderboardMonthMetrics.Metrics.FollowersChange
		followerGrowthPercentage := row.LeaderboardMonthMetrics.Metrics.FollowerGrowthPercentage
		views30d := row.LeaderboardMonthMetrics.Metrics.Views
		viewsPtr := row.AccountInfo.Views
		views := int64(0)
		if viewsPtr != nil {
			views = *viewsPtr
		}
		plays30d := row.LeaderboardMonthMetrics.Metrics.Plays
		engagementRate := row.LeaderboardMonthMetrics.Metrics.EngagementRatePercentage
		for i := range row.AccountInfo.LinkedSocials {
			if row.AccountInfo.LinkedSocials[i].Platform == string(constants.InstagramPlatform) {
				instaProfile := *row.AccountInfo.LinkedSocials[i]
				if instaProfile.Handle != nil {
					instaHandle = *instaProfile.Handle
				}
				if instaProfile.Followers != nil {
					instaFollowers = *instaProfile.Followers
					totalFollowers += *instaProfile.Followers
				}
			}
			if row.AccountInfo.LinkedSocials[i].Platform == string(constants.YoutubePlatform) {
				ytProfile := *row.AccountInfo.LinkedSocials[i]
				if ytProfile.Username != nil {
					ytUsername = *ytProfile.Username
				}
				if ytUsername == "" && ytProfile.Name != nil {
					ytUsername = *ytProfile.Name
				}
				if ytProfile.Followers != nil {
					ytFollowers = *ytProfile.Followers
					totalFollowers += *ytProfile.Followers
				}
			}
		}
		instaFollowersStr := ""
		ytFollowersStr := ""
		if instaFollowers > 0 {
			instaFollowersStr = fmt.Sprintf("%d", instaFollowers)
		}
		if ytFollowers > 0 {
			ytFollowersStr = fmt.Sprintf("%d", ytFollowers)
		}
		if platform == "ALL" {
			data[j+1] = []string{
				rank,
				rankChange,
				name,
				instaHandle,
				ytUsername,
				category,
				fmt.Sprintf("%d", totalFollowers),
				instaFollowersStr,
				ytFollowersStr,
				fmt.Sprintf("%d", followerGrowth),
				fmt.Sprintf("%f", followerGrowthPercentage),
				fmt.Sprintf("%d", views30d),
			}
		} else if platform == string(constants.InstagramPlatform) {
			data[j+1] = []string{
				rank,
				rankChange,
				name,
				instaHandle,
				category,
				fmt.Sprintf("%d", totalFollowers),
				instaFollowersStr,
				fmt.Sprintf("%d", followerGrowth),
				fmt.Sprintf("%f", followerGrowthPercentage),
				fmt.Sprintf("%f", engagementRate),
				fmt.Sprintf("%d", plays30d),
			}
		} else if platform == string(constants.YoutubePlatform) {
			data[j+1] = []string{
				rank,
				rankChange,
				name,
				ytUsername,
				category,
				fmt.Sprintf("%d", totalFollowers),
				ytFollowersStr,
				fmt.Sprintf("%d", followerGrowth),
				fmt.Sprintf("%f", followerGrowthPercentage),
				fmt.Sprintf("%d", views30d),
				fmt.Sprintf("%d", views),
			}
		}
	}
	csv, err := WriteAll(data)
	if err != nil {
		return nil, err
	}
	return csv, nil
}

func WriteAll(records [][]string) ([]byte, error) {
	if len(records) == 0 {
		return nil, errors.New("records cannot be nil or empty")
	}
	var buf bytes.Buffer
	csvWriter := csv.NewWriter(&buf)
	err := csvWriter.WriteAll(records)
	if err != nil {
		return nil, err
	}
	csvWriter.Flush()
	if err := csvWriter.Error(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

package action

import (
	"bytes"
	"coffee/app/profilecollection/domain"
	"coffee/app/profilecollection/manager"
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
	"strconv"

	log "github.com/sirupsen/logrus"
)

func PerformDownloadCollectionAction(ctx context.Context, jobBatch jobtracker.JobBatchEntry) {
	appCtx := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)
	log.Printf("Process Download Collection Job Batch %+v", jobBatch)
	jobBatch.Status = "FAILED"
	jobBatch.Remarks = "Incomplete request"
	if jobBatch.Id > 0 && jobBatch.JobBatchItems != nil && len(jobBatch.JobBatchItems) > 0 {
		jobBatchItem := jobBatch.JobBatchItems[0]
		jobBatchItem.Status = "FAILURE"
		if appCtx.PartnerId == nil || appCtx.AccountId == nil {
			jobBatch.Remarks = "Unauthorized Access"
		} else {
			jobBatch.Remarks = "Input metadata missing - collectionId/platform"
			if jobBatchItem.InputNode != nil {
				collectionId, ok1 := jobBatchItem.InputNode["collectionId"]
				shortlistId, ok4 := jobBatchItem.InputNode["shortlistId"]
				platform, ok2 := jobBatchItem.InputNode["platform"]

				if (ok1 && ok2) || (ok1 && ok2 && ok4) {
					var query coredomain.SearchQuery
					queryInput, ok3 := jobBatchItem.InputNode["query"]
					if ok3 {
						err3 := json.Unmarshal([]byte(queryInput), &query)
						log.Error(err3)
					}
					fileName := fmt.Sprintf("collection_%s_%s.csv", collectionId, platform)
					cId, _ := strconv.ParseInt(collectionId, 10, 64)
					var sId *int64
					if ok4 {
						sIdValue, _ := strconv.ParseInt(shortlistId, 10, 64)
						sId = &sIdValue
					} else {
						sId = nil
					}

					data, err := downloadCollection(ctx, cId, sId, platform, query)
					if err == nil {
						// Upload to DAM and update Job with appropriate Asset ID and URL
						damClient := dam.New(ctx)
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
		client := jobtracker.New(ctx)
		client.UpdateJobBatch(jobBatch)
	}
}

func downloadCollection(ctx context.Context, collectionId int64, shortlistId *int64, platform string, query coredomain.SearchQuery) ([]byte, error) {

	manager := manager.CreateManager(ctx)
	// shortlistid := *shortlistId
	collection, _ := manager.FindById(ctx, collectionId, false)
	if collection == nil {
		return nil, errors.New("invalid collection")
	}

	appCtx := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)
	if appCtx.PartnerId == nil || (*appCtx.PartnerId > 0 && *collection.PartnerId != *appCtx.PartnerId) {
		return nil, errors.New("you are not authorized to download this collection")
	}

	var customColumnsMeta []domain.CustomColumnMeta
	var customColumnHeaders []string
	if collection.CustomColumns != nil {
		for _, platformCCMeta := range collection.CustomColumns {
			if platformCCMeta.Platform == platform {
				customColumnsMeta = platformCCMeta.Columns
				for _, ccMeta := range platformCCMeta.Columns {
					customColumnHeaders = append(customColumnHeaders, ccMeta.Title)
				}
			}
		}
	}

	filter1 := coredomain.SearchFilter{FilterType: "EQ", Field: "collectionId", Value: strconv.Itoa(int(collectionId))}
	filter2 := coredomain.SearchFilter{FilterType: "EQ", Field: "platform", Value: platform}
	filter3 := coredomain.SearchFilter{FilterType: "EQ", Field: "enabled", Value: "true"}
	query.Filters = append(query.Filters, filter1, filter2, filter3)
	if shortlistId != nil {
		sId := *shortlistId
		filter4 := coredomain.SearchFilter{FilterType: "EQ", Field: "shortlistId", Value: strconv.Itoa(int(sId))}
		query.Filters = append(query.Filters, filter4)
	}

	page := 1
	size := 50
	filteredItemsCount := size
	igRows := []InstagramCSVRows{}
	ytRows := []YoutubeCSVRows{}
	var items []domain.ProfileCollectionItemEntry
	var err error
	for hasMoreItems := true; hasMoreItems; hasMoreItems = filteredItemsCount >= size {
		items, _, err = manager.SearchItemsInCollection(ctx, query, "id", "DESC", page, size)
		if err != nil {
			return nil, err
		}

		page = page + 1
		filteredItemsCount = len(items)

		if platform == string(constants.InstagramPlatform) && len(items) > 0 {
			for _, item := range items {
				if item.Profile == nil {
					continue
				}
				profile := *item.Profile
				row := InstagramCSVRows{}
				if profile.Name != nil {
					row.Name = *profile.Name
				} else {
					row.Name = *profile.Handle
				}
				row.URL = ""
				if profile.Handle != nil {
					url := "https://www.instagram.com/" + *profile.Handle
					row.URL = url
				}
				row.Followers = ""
				if profile.Metrics.Followers != nil {
					row.Followers = fmt.Sprint(*profile.Metrics.Followers)
				}
				row.ER = fmt.Sprintf("%.2f", profile.Metrics.EngagementRatePercenatge) + "%"
				row.AvgLikes = ""
				if profile.Metrics.AvgLikes != nil {
					row.AvgLikes = fmt.Sprint(*profile.Metrics.AvgLikes)
				}
				if profile.Metrics.AvgReelsPlayCount != nil {
					row.AvgReelsPlayCount = fmt.Sprintf("%.2f", *profile.Metrics.AvgReelsPlayCount)
				}
				if profile.Metrics.AvgReach != nil {
					row.EstimatedReach = fmt.Sprintf("%.2f", *profile.Metrics.AvgReach)
				}
				if profile.Metrics.ReelsReach != nil {
					row.ReelReach = fmt.Sprintf("%d", *profile.Metrics.ReelsReach)
				}
				if profile.Metrics.StoryReach != nil {
					row.StoryReach = fmt.Sprintf("%d", *profile.Metrics.StoryReach)
				}
				if profile.Metrics.ImageReach != nil {
					row.ImageReach = fmt.Sprintf("%d", *profile.Metrics.ImageReach)
				}
				row.CustomColumns = getCustomColumnValuesForItem(item, customColumnsMeta)
				igRows = append(igRows, row)
			}
		} else if platform == string(constants.YoutubePlatform) && len(items) > 0 {
			for _, item := range items {
				if item.Profile == nil {
					continue
				}
				profile := *item.Profile
				row := YoutubeCSVRows{}
				if profile.Name != nil {
					row.Name = *profile.Name
				} else {
					row.Name = *profile.Handle
				}
				row.URL = ""
				if profile.Handle != nil {
					url := "https://www.youtube.com/channel/" + *profile.Handle
					row.URL = url
				}
				row.Subscribers = ""
				if profile.Metrics.Followers != nil {
					row.Subscribers = fmt.Sprint(*profile.Metrics.Followers)
				}
				row.AvgViews = ""
				if profile.Metrics.AvgViews != nil {
					row.AvgViews = fmt.Sprintf("%.2f", *profile.Metrics.AvgViews)

				}
				row.CustomColumns = getCustomColumnValuesForItem(item, customColumnsMeta)
				ytRows = append(ytRows, row)
			}
		}
	}

	if len(igRows) > 0 || len(ytRows) > 0 {
		if platform == string(constants.InstagramPlatform) {
			headers := []string{"Name", "Url", "Followers", "EnagagementRate", "Avg Likes", "Avg Reels Plays", "Reels Reach", "Image Reach", "Story Reach"}
			headers = append(headers, customColumnHeaders...)
			data := make([][]string, len(igRows))
			for i, row := range igRows {
				defaultColumns := []string{
					row.Name,
					row.URL,
					row.Followers,
					row.ER,
					row.AvgLikes,
					row.AvgReelsPlayCount,
					row.ReelReach,
					row.ImageReach,
					row.StoryReach,
				}
				data[i] = append(defaultColumns, row.CustomColumns...)
			}
			csv, err := writeAllRows(headers, data)
			if err != nil {
				return nil, err
			}
			return csv, nil
		} else if platform == string(constants.YoutubePlatform) {
			headers := []string{"Name", "Url", "Subscribers", "Avg Views"}
			headers = append(headers, customColumnHeaders...)
			data := make([][]string, len(ytRows))
			for i, row := range ytRows {
				defaultColumns := []string{
					row.Name,
					row.URL,
					row.Subscribers,
					row.AvgViews,
				}
				data[i] = append(defaultColumns, row.CustomColumns...)
			}
			csv, err := writeAllRows(headers, data)
			if err != nil {
				return nil, err
			}
			return csv, nil
		}
	}
	return nil, errors.New("empty collection")
}

func writeAllRows(headers []string, records [][]string) ([]byte, error) {
	if len(records) == 0 {
		return nil, errors.New("records cannot be nil or empty")
	}
	var buf bytes.Buffer
	csvWriter := csv.NewWriter(&buf)
	csvWriter.Write(headers)
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

func getCustomColumnValuesForItem(item domain.ProfileCollectionItemEntry, customColumnsMeta []domain.CustomColumnMeta) []string {
	itemCustomColumnValuesMap := map[string]string{}
	if item.CustomColumnsData != nil {
		for _, col := range item.CustomColumnsData {
			itemCustomColumnValuesMap[col.Key] = col.Value
		}
	}
	var customColumnValues []string
	for _, columnMeta := range customColumnsMeta {
		itemVal, ok := itemCustomColumnValuesMap[*columnMeta.Key]
		if (!ok || itemVal == "") && columnMeta.Value != nil && *columnMeta.Value != "" {
			itemVal = *columnMeta.Value
		}
		customColumnValues = append(customColumnValues, itemVal)
	}
	return customColumnValues
}

// Name,Email ID,Instagram URL,Insta Followers,Insta Engagement,Insta Avg Likes,Avg Reach,Avg Video Views,Avg Reel Play Count,Avg Reel views
// Name,Email ID,Youtube URL,Youtube Subscribers,Youtube Avg Views,Youtube Video Count

type InstagramCSVRows struct {
	Name              string
	URL               string
	Followers         string
	ER                string
	AvgLikes          string
	AvgReelsPlayCount string
	EstimatedReach    string
	ReelReach         string
	ImageReach        string
	StoryReach        string
	CustomColumns     []string
}

type YoutubeCSVRows struct {
	Name          string
	URL           string
	Subscribers   string
	AvgViews      string
	CustomColumns []string
}

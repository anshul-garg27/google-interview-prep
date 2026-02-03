package listeners

import (
	"bytes"
	"coffee/app/collectionanalytics/manager"
	postcollectiondomain "coffee/app/postcollection/domain"
	postcollectionmanager "coffee/app/postcollection/manager"
	"coffee/client/coffeeclient"
	"coffee/client/dam"
	"coffee/client/jobtracker"
	"coffee/constants"
	"coffee/core/appcontext"
	"coffee/core/domain"
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
	log "github.com/sirupsen/logrus"
)

func SetupListeners(router *message.Router, subscriber *amqp.Subscriber) {
	router.AddNoPublisherHandler(
		"download_collection_creator_report",
		"jobtracker.dx___download_collection_creator_report_q",
		subscriber,
		PerformJobCollectionInfluencerReport,
	)

	router.AddNoPublisherHandler(
		"sentiment_collection_report_out",
		"coffee.dx___sentiment_collection_report_out_q",
		subscriber,
		PerformJobStoreCollectionReportPath,
	)
}

func PerformJobCollectionInfluencerReport(message *message.Message) error {
	// This method is used to listen to request for
	// downloading a collection
	ctx := message.Context()
	appCtx := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)
	jobBatch := &jobtracker.JobBatchEntry{}
	if err := json.Unmarshal(message.Payload, jobBatch); err == nil {
		log.Printf("Process Download Collection Job Batch %+v", jobBatch)
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
					collectionId, ok2 := jobBatchItem.InputNode["collectionId"]
					collectionType, ok3 := jobBatchItem.InputNode["collectionType"]

					if ok1 && ok2 && ok3 {
						fileName := fmt.Sprintf("collection-%s-%s.csv", platform, collectionId)
						data, err := DownloadCollectionInfluencerReport(message.Context(), platform, collectionId, collectionType)
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

func DownloadCollectionInfluencerReport(ctx context.Context, platform string, collectionId string, collectionType string) ([]byte, error) {
	manager := manager.CreateCollectionAnalyticsManager(ctx)
	collectionIdFilter := domain.SearchFilter{
		FilterType: "EQ",
		Field:      "collectionId",
		Value:      collectionId,
	}
	collectionTypeFilter := domain.SearchFilter{
		FilterType: "EQ",
		Field:      "collectionType",
		Value:      collectionType,
	}
	filters := []domain.SearchFilter{collectionIdFilter, collectionTypeFilter}
	query := domain.SearchQuery{Filters: filters}
	postData, err := manager.FetchCollectionPostsWithMetricsSummary(ctx, query, "reach", "DESC", 1, 1000)
	if err != nil {
		return nil, err
	}
	data := make([][]string, len(postData)+1)
	columns := []string{
		"Creator Name",
		"Video Link",
		"Posted On",
		"Total Views",
		"Cost Per View",
		"Total Likes",
		"Total Comments",
		"Engagement Rate",
		"Spends",
		"Views - Day 1",
		"Views - Day 2",
		"Views - Day 3",
		"Views - Day 4",
		"Views - Day 5",
		"Views - Day 6",
		"Views - Day 7",
		"Views - Day 14",
		"Views - Day 21",
		"Views - Day 28",
		"Views - Day 28",
		"Views - Day 35",
		"Views - Day 42",
		"Day 1 Views/Total Views",
		"Week 1 Views/Total Views",
	}
	data[0] = columns
	for j, _ := range postData {
		post := postData[j]
		name := ""
		link := ""
		postedOn := ""
		totalViews := ""
		costPerView := ""
		totalLikes := ""
		totalComments := ""
		engRate := ""
		spends := ""
		D1Views := ""
		D2Views := ""
		D3Views := ""
		D4Views := ""
		D5Views := ""
		D6Views := ""
		D7Views := ""
		D14Views := ""
		D21Views := ""
		D28Views := ""
		D35Views := ""
		D42Views := ""
		D1ViewRatio := ""
		W1ViewRatio := ""
		if post.ProfileName != nil {
			name = *post.ProfileName
		}
		link = post.PostLink
		postedOn = fmt.Sprintf("%s", post.PublishedAt.Format("2006-01-02"))
		totalViews = fmt.Sprintf("%d", post.Views)
		costPerView = post.CostPerView
		totalLikes = fmt.Sprintf("%d", post.Likes)
		totalComments = fmt.Sprintf("%d", post.Comments)
		engRate = post.EngagementRate
		spends = fmt.Sprintf("%d", post.Cost)
		postFilter := domain.SearchFilter{
			FilterType: "EQ",
			Field:      "shortcode",
			Value:      post.PostShortCode,
		}
		postFilters := []domain.SearchFilter{collectionIdFilter, collectionTypeFilter, postFilter}
		postQuery := domain.SearchQuery{Filters: postFilters}
		postTimeSeries, _ := manager.FetchCollectionTimeSeries(ctx, postQuery)
		for i, _ := range postTimeSeries {
			if i == 1 {
				D1Views = fmt.Sprintf("%d", postTimeSeries[i].Views)
				if post.Views > 0 {
					D1ViewRatio = fmt.Sprintf("%f", float64(postTimeSeries[i].Views)/float64(post.Views))
				}
			}
			if i == 2 {
				D2Views = fmt.Sprintf("%d", postTimeSeries[i].Views)
			}
			if i == 3 {
				D3Views = fmt.Sprintf("%d", postTimeSeries[i].Views)
			}
			if i == 4 {
				D4Views = fmt.Sprintf("%d", postTimeSeries[i].Views)
			}
			if i == 5 {
				D5Views = fmt.Sprintf("%d", postTimeSeries[i].Views)
			}
			if i == 6 {
				D6Views = fmt.Sprintf("%d", postTimeSeries[i].Views)
			}
			if i == 7 {
				D7Views = fmt.Sprintf("%d", postTimeSeries[i].Views)
				if post.Views > 0 {
					W1ViewRatio = fmt.Sprintf("%f", float64(postTimeSeries[i].Views)/float64(post.Views))
				}
			}
			if i == 14 {
				D14Views = fmt.Sprintf("%d", postTimeSeries[i].Views)
			}
			if i == 21 {
				D21Views = fmt.Sprintf("%d", postTimeSeries[i].Views)
			}
			if i == 28 {
				D28Views = fmt.Sprintf("%d", postTimeSeries[i].Views)
			}
			if i == 35 {
				D35Views = fmt.Sprintf("%d", postTimeSeries[i].Views)
			}
			if i == 42 {
				D42Views = fmt.Sprintf("%d", postTimeSeries[i].Views)
			}
		}
		data[j+1] = []string{
			name,
			link,
			postedOn,
			totalViews,
			costPerView,
			totalLikes,
			totalComments,
			engRate,
			spends,
			D1Views,
			D2Views,
			D3Views,
			D4Views,
			D5Views,
			D6Views,
			D7Views,
			D14Views,
			D21Views,
			D28Views,
			D35Views,
			D42Views,
			D1ViewRatio,
			W1ViewRatio,
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

func PerformJobStoreCollectionReportPath(message *message.Message) error {
	ctx := message.Context()
	sentimentReportEntry := &coffeeclient.SentimentReportEntry{}
	var error error
	if err := json.Unmarshal(message.Payload, sentimentReportEntry); err == nil {
		postcollectionmanager := postcollectionmanager.CreatePostCollectionManager(ctx)
		postCollectionEntry := &postcollectiondomain.PostCollectionEntry{
			Id:                    &sentimentReportEntry.CollectionId,
			SentimentReportPath:   &sentimentReportEntry.SentimentReportPath,
			SentimentReportBucket: &sentimentReportEntry.SentimentReportBucket,
		}
		postCollectionEntry, error = postcollectionmanager.Update(ctx, sentimentReportEntry.CollectionId, postCollectionEntry)
		if error != nil {
			errMessage := fmt.Sprintf("could not update sentiment report path for collectionId = %s", sentimentReportEntry.CollectionId)
			log.Error(errMessage)
			return nil
		}
	}
	return nil
}

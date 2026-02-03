package listeners

import (
	"coffee/app/postcollection/domain"
	postcollectionmanager "coffee/app/postcollection/manager"
	beatservice "coffee/client/beat"
	"coffee/client/jobtracker"
	"coffee/constants"
	"coffee/core/appcontext"
	"coffee/helpers"
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
)

func SetupListeners(router *message.Router, subscriber *amqp.Subscriber) {
	router.AddNoPublisherHandler(
		"add_item_to_post_collection",
		"jobtracker.dx___add_item_post_collection_q",
		subscriber,
		PerformJobAddItemToCollection,
	)
}

func PerformJobAddItemToCollection(message *message.Message) error {
	// This method is used to listen to request for
	// creating a new collection from a collection
	ctx := message.Context()
	appCtx := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)
	jobBatch := &jobtracker.JobBatchEntry{}
	err := json.Unmarshal(message.Payload, jobBatch)
	if err == nil {
		if jobBatch != nil && jobBatch.Id > 0 && jobBatch.JobBatchItems != nil && len(jobBatch.JobBatchItems) > 0 {
			log.Printf("Process AddItemToCollection Job Batch %+v", jobBatch)
			errorMessage := ""
			var c domain.PostCollectionEntry
			if appCtx.PartnerId != nil && appCtx.AccountId != nil {
				jobBatchItem := jobBatch.JobBatchItems[0]
				if jobBatchItem.InputNode != nil {
					collectionId, found := jobBatchItem.InputNode["collectionId"]
					if found {
						postcollectionManager := postcollectionmanager.CreatePostCollectionManager(ctx)
						collection, _ := postcollectionManager.FindById(ctx, collectionId, true)
						if collection == nil {
							errorMessage = "Invalid Collection"
						} else {
							c = *collection
						}
					}
				}
			} else {
				errorMessage = "Unauthorized Access"
			}

			for i, jobBatchItem := range jobBatch.JobBatchItems {
				outputNode := map[string]string{}
				itemRemarks := errorMessage
				itemStatus := "FAILURE"
				if errorMessage == "" {
					if (appCtx.PlanType == nil || *appCtx.PlanType == constants.FreePlan) && *c.Source != string(constants.GccCampaignCollection) && *c.PartnerId != -1 && c.TotalPosts >= 10 {
						itemRemarks = "You are not allowed to add more than 10 posts in one collection. Please upgrade your account."
					} else {
						itemEntry, err := TransformJobItemToCollectionItem(*c.Id, jobBatchItem)
						if err == nil && itemEntry != nil {
							newCollectionItemId, err := ProcessAddItemToCollectionJobItem(message.Context(), c, *itemEntry)
							if err != nil {
								itemRemarks = err.Error()
							} else {
								itemStatus = "SUCCESS"
								itemRemarks = "Item Added Successfully"
								outputNode["collectionItemId"] = strconv.FormatInt(newCollectionItemId, 10)
								outputNode["shortCode"] = *itemEntry.ShortCode
								outputNode["platform"] = *itemEntry.Platform
								c.TotalPosts++
							}
						} else {
							if err != nil {
								itemRemarks = err.Error()
							} else {
								itemRemarks = "something went wrong"
							}
						}
					}
				}
				outputNode["STATUS"] = itemStatus
				outputNode["REMARK"] = itemRemarks
				jobBatchItem.Status = itemStatus
				jobBatchItem.OutputNode = outputNode
				jobBatch.JobBatchItems[i] = jobBatchItem
			}
			jobBatch.Status = "COMPLETED"
			jobBatch.Remarks = "Items Processed"
			client := jobtracker.New(message.Context())
			client.UpdateJobBatch(*jobBatch)
		}
	} else {
		return err
	}
	return err
}

func ProcessAddItemToCollectionJobItem(ctx context.Context, collection domain.PostCollectionEntry, itemEntry domain.PostCollectionItemEntry) (int64, error) {
	postcollectionItemManager := postcollectionmanager.CreatePostCollectionItemManager(ctx)
	items, err := postcollectionItemManager.AddItemsToCollection(ctx, collection, []domain.PostCollectionItemEntry{itemEntry})
	if err == nil && items != nil && len(*items) > 0 {
		item := (*items)[0]
		return *item.Id, nil
	}
	if err != nil {
		return 0, err
	}
	return 0, errors.New("failed to add item to collection")
}

func TransformJobItemToCollectionItem(collectionId string, jobItem jobtracker.JobBatchItem) (*domain.PostCollectionItemEntry, error) {
	platformPtr := GetStringValueFromMap(jobItem.InputNode, "Platform")
	platform := ""
	if platformPtr != nil {
		platform = strings.ToUpper(*platformPtr)
	}
	if platform == "" || (platform != string(constants.InstagramPlatform) && platform != string(constants.YoutubePlatform)) {
		return nil, errors.New("invalid platform")
	}
	postTypePtr := GetStringValueFromMap(jobItem.InputNode, "Post Type")
	if postTypePtr == nil {
		return nil, errors.New("invalid post type")
	}
	postType := strings.ToLower(*postTypePtr)
	if postType != "reels" && postType != "image" && postType != "story" && postType != "video" && postType != "short" && postType != "carousel" {
		return nil, errors.New("invalid post type")
	}
	postedBy := GetStringValueFromMap(jobItem.InputNode, "Posted By (only for story posts)")
	postLink := GetStringValueFromMap(jobItem.InputNode, "Post Link")
	postThumbnailPtr := GetStringValueFromMap(jobItem.InputNode, "Post Thumbnail (optional)")
	postThumbnail := ""
	if postThumbnailPtr != nil {
		postThumbnail = *postThumbnailPtr
	}
	if !strings.HasPrefix("http", postThumbnail) {
		postThumbnail = ""
	}
	cost := GetInt64ValueFromMap(jobItem.InputNode, "Cost")
	postTitle := ""
	postedByHandle := ""
	shortCode := ""
	if postLink != nil {
		platformEx, postId, err := helpers.ExtractPlatformPostIdentifier(*postLink)
		if err != nil {
			if postType != "" && postType != "story" {
				return nil, err
			}
		}
		if platform != "" && platform != platformEx {
			if postType != "" && postType != "story" {
				return nil, errors.New("invalid platform link")
			}
		}
		if postId != "" {
			shortCode = postId
		} else {
			if postType != "" && postType != "story" {
				return nil, errors.New("invalid post link")
			}
		}
	}
	if postedBy != nil {
		platformEx, idType, id, err := helpers.ExtractPlatformProfileIdentifier(*postedBy)
		if err != nil {
			if postType == "story" { // Only mandatory for stories
				return nil, err
			}
		}
		if platform != "" && platform != platformEx {
			if postType == "story" { // Only mandatory for stories
				return nil, errors.New("invalid platform link")
			}
		}
		if platformEx == string(constants.YoutubePlatform) && idType == "HANDLE" {
			beatClient := beatservice.New(context.TODO())
			ytProfile, err := beatClient.FindYoutubeChannelIdByHandle(id)
			if err != nil {
				return nil, err
			}
			if ytProfile != nil && ytProfile.YoutubeChannel.ChannelId != "" {
				postedByHandle = ytProfile.YoutubeChannel.ChannelId
			} else {
				log.Error("Invalid channel link - " + *postedBy)
				postedByHandle = "" // Ignore the failure case
			}
		} else { // For all other cases, handle == id (instagram handle, or youtube channel id)
			postedByHandle = id
		}
	}
	if postType == "story" && shortCode == "" {
		shortCode = helpers.GenerateStoryDummyShortCode(postedByHandle)
	}
	return &domain.PostCollectionItemEntry{
		Platform:         &platform,
		ShortCode:        &shortCode,
		PostCollectionId: &collectionId,
		Cost:             cost,
		Enabled:          helpers.ToBool(true),
		PostType:         &postType,
		PostedByHandle:   &postedByHandle,
		SponsorLinks:     []string{},
		PostTitle:        &postTitle,
		PostLink:         postLink,
		PostThumbnail:    &postThumbnail,
	}, nil
}

func GetStringValueFromMap(inputNode map[string]string, key string) *string {
	value, ok := inputNode[key]
	if ok {
		return &value
	}
	return nil
}

func GetInt64ValueFromMap(inputNode map[string]string, key string) *int64 {
	value, ok := inputNode[key]
	if ok {
		val, _ := strconv.ParseInt(value, 10, 64)
		return &val
	}
	return nil
}

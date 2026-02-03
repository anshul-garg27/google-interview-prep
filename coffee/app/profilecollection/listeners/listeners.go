package listeners

import (
	domain2 "coffee/app/profilecollection/domain"
	"coffee/app/profilecollection/listeners/action"
	"coffee/app/profilecollection/manager"
	"coffee/client/jobtracker"
	"coffee/constants"
	"coffee/core/appcontext"
	coredomain "coffee/core/domain"
	"coffee/helpers"
	"context"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
	log "github.com/sirupsen/logrus"
	"go.step.sm/crypto/randutil"
)

func SetupListeners(router *message.Router, subscriber *amqp.Subscriber) {
	router.AddNoPublisherHandler(
		"duplicate_collection",
		"jobtracker.dx___duplicate_collection_q",
		subscriber,
		PerformJobDuplicateCollection,
	)

	router.AddNoPublisherHandler(
		"download_collection",
		"jobtracker.dx___download_collection_q",
		subscriber,
		PerformJobDownloadCollection,
	)

	router.AddNoPublisherHandler(
		"import_from_profile_collection",
		"jobtracker.dx___import_from_profile_collection_q",
		subscriber,
		PerformJobImportFromCollection,
	)

	router.AddNoPublisherHandler(
		"add_item_to_profile_collection",
		"jobtracker.dx___add_item_profile_collection_q",
		subscriber,
		PerformJobAddItemToCollection,
	)
}

func PerformJobDuplicateCollection(message *message.Message) error {
	// This method is used to listen to request for
	// creating a new collection from a collection
	ctx := message.Context()
	appCtx := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)
	jobBatch := &jobtracker.JobBatchEntry{}
	err := json.Unmarshal(message.Payload, jobBatch)
	if err == nil && jobBatch != nil && jobBatch.Id > 0 {
		log.Printf("Process Duplicate Collection Job Batch %+v", jobBatch)
		errorMessage := ""
		outputNode := map[string]string{}
		if appCtx.PartnerId == nil || appCtx.AccountId == nil {
			errorMessage = "Unauthorized Access"
		}

		if jobBatch.JobBatchItems != nil && len(jobBatch.JobBatchItems) > 0 {
			jobBatchItem := jobBatch.JobBatchItems[0]
			if errorMessage == "" {
				errorMessage = "Collection Id missing in request"
				if jobBatchItem.InputNode != nil {
					collectionId, ok := jobBatchItem.InputNode["collectionId"]
					if ok {
						fromCId, _ := strconv.ParseInt(collectionId, 10, 64)
						manager := manager.CreateManager(ctx)
						fromCollection, newCollection, err := CreateCollection(message.Context(), *appCtx.PartnerId, fromCId, manager)
						if err != nil {
							errorMessage = err.Error()
						} else {
							outputNode["collectionId"] = strconv.FormatInt(*newCollection.Id, 10)
							err = DuplicateCollectionItemsForPlatform(ctx, manager, *fromCollection, *newCollection, string(constants.InstagramPlatform))
							if err != nil {
								errorMessage = err.Error()
							}
							err = DuplicateCollectionItemsForPlatform(ctx, manager, *fromCollection, *newCollection, string(constants.YoutubePlatform))
							if err != nil {
								errorMessage = err.Error()
							}
						}
					}
				}
			}
			if errorMessage != "" {
				outputNode["STATUS"] = "FAILURE"
				outputNode["REMARK"] = errorMessage
				jobBatchItem.OutputNode = outputNode
			} else {
				outputNode["STATUS"] = "SUCCESS"
				outputNode["REMARK"] = "Collection Duplicated Successfully"
				jobBatchItem.OutputNode = outputNode
				jobBatch.JobBatchItems[0] = jobBatchItem
			}
		}
		jobBatch.Status = "COMPLETED"
		jobBatch.Remarks = "Processing Complete"
		client := jobtracker.New(message.Context())
		client.UpdateJobBatch(*jobBatch)
	} else {
		return err
	}
	return err
}

func PerformJobDownloadCollection(message *message.Message) error {
	// This method is used to listen to request for
	// downloading a collection
	ctx := message.Context()
	jobBatch := &jobtracker.JobBatchEntry{}
	if err := json.Unmarshal(message.Payload, jobBatch); err == nil && jobBatch != nil {
		action.PerformDownloadCollectionAction(ctx, *jobBatch)
	}
	return nil
}

func CreateCollection(ctx context.Context, partnerId int64, collectionId int64, manager *manager.Manager) (*domain2.ProfileCollectionEntry, *domain2.ProfileCollectionEntry, error) {
	appCtx := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)
	accountIdStr := strconv.FormatInt(*appCtx.AccountId, 10)

	collectionEntry, err := manager.FindById(ctx, collectionId, false)
	if err != nil {
		return nil, nil, err
	}
	if collectionEntry == nil {
		return nil, nil, errors.New("invalid collection Id to duplicate. No collection found")
	}

	// ************
	shareId, _ := randutil.Alphanumeric(7)
	name := "Copy Of " + *collectionEntry.Name
	newCollectionEntry := &domain2.ProfileCollectionEntry{
		Name:        &name,
		PartnerId:   &partnerId,
		Description: collectionEntry.Description,
		Featured:    helpers.ToBool(false),
		Enabled:     helpers.ToBool(true),
		CreatedBy:   &accountIdStr,
		ShareId:     &shareId,
		Categories:  collectionEntry.Categories,
		Tags:        collectionEntry.Tags,
	}
	createdCollectionEntry, err := manager.Create(ctx, newCollectionEntry)
	if err != nil {
		return nil, nil, err
	}
	if createdCollectionEntry == nil {
		return nil, nil, errors.New("failed to duplicate collection. Try again later")
	}
	// ************
	return collectionEntry, createdCollectionEntry, nil
}

func DuplicateCollectionItemsForPlatform(ctx context.Context, manager *manager.Manager, oldCollection domain2.ProfileCollectionEntry, newCollection domain2.ProfileCollectionEntry, platform string) error {
	// Find all collection items and replicate them with new collection ID
	filter := coredomain.SearchFilter{FilterType: "EQ", Field: "collectionId", Value: strconv.Itoa(int(*oldCollection.Id))}
	filter1 := coredomain.SearchFilter{FilterType: "EQ", Field: "platform", Value: platform}
	filters := []coredomain.SearchFilter{filter, filter1}
	query := coredomain.SearchQuery{
		Filters: filters,
	}

	page := 1
	size := 50
	filteredCount := int64(50)
	items := []domain2.ProfileCollectionItemEntry{}
	var err error
	for hasNextPage := true; hasNextPage; hasNextPage = filteredCount >= int64(size) {
		itemEntries := []domain2.ProfileCollectionItemEntry{}
		items, filteredCount, err = manager.SearchItemsInCollection(ctx, query, "id", "DESC", page, size)
		if err != nil {
			return err
		}
		page = page + 1
		if len(items) > 0 {
			for _, item := range items {
				itemEntry := &domain2.ProfileCollectionItemEntry{
					PlatformAccountCode: item.PlatformAccountCode,
					Platform:            item.Platform,
				}
				itemEntries = append(itemEntries, *itemEntry)
			}
			_, err, _ := manager.AddItemsToCollection(ctx, newCollection, itemEntries)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func PerformJobImportFromCollection(message *message.Message) error {
	// This method is used to listen to request for
	// creating a new collection from a collection
	ctx := message.Context()
	appCtx := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)
	jobBatch := &jobtracker.JobBatchEntry{}
	manager := manager.CreateManager(ctx)
	err := json.Unmarshal(message.Payload, jobBatch)
	if err == nil && jobBatch != nil && jobBatch.Id > 0 {
		log.Printf("Process Import From Collection Job Batch %+v", jobBatch)
		errorMessage := ""
		outputNode := map[string]string{}
		if appCtx.PartnerId == nil || appCtx.AccountId == nil {
			errorMessage = "Unauthorized Access"
		}

		if jobBatch.JobBatchItems != nil && len(jobBatch.JobBatchItems) > 0 {
			jobBatchItem := jobBatch.JobBatchItems[0]
			if errorMessage == "" {
				if jobBatchItem.InputNode != nil {
					collectionId, found := jobBatchItem.InputNode["fromCollectionId"]
					if found {
						fromCId, _ := strconv.ParseInt(collectionId, 10, 64)
						toCollectionId, found := jobBatchItem.InputNode["toCollectionId"]
						platform, foundPlatform := jobBatchItem.InputNode["platform"]
						if !found {
							errorMessage = "To CollectionId is missing"
						} else {
							toCId, _ := strconv.ParseInt(toCollectionId, 10, 64)
							fromCollection, err1 := manager.FindById(ctx, fromCId, false)
							if err1 != nil {
								errorMessage = "Invalid From Collection"
							}
							toCollection, err2 := manager.FindById(ctx, toCId, false)
							if err2 != nil {
								errorMessage = "Invalid To Collection"
							}
							if err1 == nil && err2 == nil {
								importInstagram := true
								importYoutube := true
								if foundPlatform {
									if platform == string(constants.InstagramPlatform) {
										importYoutube = false
									} else if platform == string(constants.YoutubePlatform) {
										importInstagram = false
									}
								}
								if importInstagram {
									err = DuplicateCollectionItemsForPlatform(ctx, manager, *fromCollection, *toCollection, string(constants.InstagramPlatform))
									if err != nil {
										errorMessage = err.Error()
									}
								}
								if importYoutube {
									err = DuplicateCollectionItemsForPlatform(ctx, manager, *fromCollection, *toCollection, string(constants.YoutubePlatform))
									if err != nil {
										errorMessage = err.Error()
									}
								}
							}
						}
					} else {
						errorMessage = "From Collection Id missing in request"
					}
				} else {
					errorMessage = "From Collection Id missing in request"
				}
			}
			if errorMessage != "" {
				jobBatchItem.Status = "FAILURE"
				outputNode["STATUS"] = "FAILURE"
				outputNode["REMARK"] = errorMessage
				jobBatchItem.OutputNode = outputNode
			} else {
				jobBatchItem.Status = "SUCCESS"
				outputNode["STATUS"] = "SUCCESS"
				outputNode["REMARK"] = "Items Imported Successfully"
				jobBatchItem.OutputNode = outputNode
			}
			jobBatch.JobBatchItems[0] = jobBatchItem
		}
		jobBatch.Status = "COMPLETED"
		jobBatch.Remarks = "Processing Complete"
		client := jobtracker.New(message.Context())
		client.UpdateJobBatch(*jobBatch)
	} else {
		return err
	}
	return err
}

func PerformJobAddItemToCollection(message *message.Message) error {
	// This method is used to listen to request for
	// creating a new collection from a collection
	ctx := message.Context()
	manager := manager.CreateManager(ctx)
	appCtx := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)
	jobBatch := &jobtracker.JobBatchEntry{}
	err := json.Unmarshal(message.Payload, jobBatch)
	collectionPlatform := ""
	if err == nil {
		if jobBatch != nil && jobBatch.Id > 0 && jobBatch.JobBatchItems != nil && len(jobBatch.JobBatchItems) > 0 {
			log.Printf("Process AddItemToCollection Job Batch %+v", jobBatch)
			errorMessage := ""
			var c domain2.ProfileCollectionEntry
			if appCtx.PartnerId != nil && appCtx.AccountId != nil {
				jobBatchItem := jobBatch.JobBatchItems[0]
				if jobBatchItem.InputNode != nil {
					collectionId, found := jobBatchItem.InputNode["collectionId"]
					_platform, foundPlatform := jobBatchItem.InputNode["collectionPlatform"]
					if foundPlatform {
						collectionPlatform = _platform
					}
					if found {
						cId, _ := strconv.ParseInt(collectionId, 10, 64)
						collection, err := manager.FindById(ctx, cId, true)
						if err != nil {
							errorMessage = err.Error()
						} else if collection == nil {
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
					totalItems := int64(0)
					for _, summary := range c.PlatformWiseItemsSummary {
						totalItems += summary.ProfileCounter.Count
					}
					if (appCtx.PlanType == nil || *appCtx.PlanType == constants.FreePlan) && *c.Source == string(constants.SaasAtCollection) && *c.PartnerId != -1 && totalItems >= 5 {
						itemRemarks = "You are not allowed to add more than 5 profiles in one collection. Please upgrade your account."
					} else {
						itemEntry, err := TransformJobItemToCollectionItem(*c.Id, jobBatchItem)
						if err != nil {
							log.Error(err)
						}
						if itemEntry.Platform == nil || *itemEntry.Platform != collectionPlatform {
							itemRemarks = "Invalid Platform"
							itemStatus = "FAILURE"
						} else {
							newCollectionItemId, err := ProcessAddItemToCollectionJobItem(message.Context(), c, itemEntry)
							if err != nil {
								itemRemarks = err.Error()
							} else {
								itemStatus = "SUCCESS"
								itemRemarks = "Item Added Successfully"
								outputNode["collectionItemId"] = strconv.FormatInt(newCollectionItemId, 10)
								outputNode["handle"] = *itemEntry.Profile.Handle
								outputNode["platform"] = *itemEntry.Platform
								totalItems++
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

func ProcessAddItemToCollectionJobItem(ctx context.Context, collection domain2.ProfileCollectionEntry, itemEntry domain2.ProfileCollectionItemEntry) (int64, error) {
	collectionManager := manager.CreateManager(ctx)
	items, err, _ := collectionManager.AddItemsToCollection(ctx, collection, []domain2.ProfileCollectionItemEntry{itemEntry})
	if err == nil && items != nil && len(items) > 0 {
		item := (items)[0]
		return *item.Id, nil
	}
	if err != nil {
		return 0, err
	}
	return 0, errors.New("failed to add item to collection")
}

func TransformJobItemToCollectionItem(collectionId int64, jobItem jobtracker.JobBatchItem) (domain2.ProfileCollectionItemEntry, error) {
	platform := GetStringValueFromMap(jobItem.InputNode, "platform")
	handle := GetStringValueFromMap(jobItem.InputNode, "handle")
	return domain2.ProfileCollectionItemEntry{
		Platform:            platform,
		ProfileCollectionId: &collectionId,
		Profile: &coredomain.Profile{
			Handle: handle,
		},
	}, nil
}

func GetStringValueFromMap(inputNode map[string]string, key string) *string {
	value, ok := inputNode[key]
	if ok {
		return &value
	}
	return nil
}

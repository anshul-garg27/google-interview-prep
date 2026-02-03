package manager

import (
	discoverymanager "coffee/app/discovery/manager"
	"coffee/app/profilecollection/dao"
	"coffee/app/profilecollection/domain"
	"coffee/constants"
	"coffee/core/appcontext"
	coredomain "coffee/core/domain"
	"context"
	"fmt"
	"strconv"
	"strings"
)

func (m *Manager) EnrichCollectionWithItemsSummary(ctx context.Context, entry domain.ProfileCollectionEntry, fetchTopAccounts bool) *domain.ProfileCollectionEntry {
	itemCountsRows := m.itemDAO.CountCollectionItems(ctx, *entry.Id)
	if len(itemCountsRows) == 0 {
		return &entry
	}

	platformWiseItemsSummaryMap := map[string]domain.PlatformWiseItemsSummary{}
	for _, row := range itemCountsRows {
		platformItemsSummary := domain.PlatformWiseItemsSummary{
			ProfileCounter: domain.ProfileCounter{
				Count:      row.Allitems,
				EmailCount: row.Emailitems,
				PhoneCount: row.Phoneitems,
			},
		}
		if fetchTopAccounts && row.Allitems > 0 {
			platformItemsSummary.TopFiveAccounts = m.GetTopFiveItemsForPlatform(ctx, *entry.Id, row.Platform)
		}
		platformWiseItemsSummaryMap[row.Platform] = platformItemsSummary
	}
	entry.PlatformWiseItemsSummary = platformWiseItemsSummaryMap
	return &entry
}

func (m *Manager) GetTopFiveItemsForPlatform(ctx context.Context, collectionId int64, platform string) []domain.CompactSocialAccount {
	filter1 := coredomain.SearchFilter{FilterType: "EQ", Field: "platform", Value: platform}
	filter2 := coredomain.SearchFilter{FilterType: "EQ", Field: "enabled", Value: "true"}
	filter3 := coredomain.SearchFilter{FilterType: "EQ", Field: "profile_collection_id", Value: strconv.FormatInt(collectionId, 10)}
	filters := []coredomain.SearchFilter{filter1, filter2, filter3}
	query := coredomain.SearchQuery{Filters: filters}

	var topFive []domain.CompactSocialAccount

	socialDiscoveryManager := discoverymanager.CreateSearchManager(ctx)
	itemEntities, _, _ := m.itemDAO.Search(ctx, query, "created_at", "DESC", 1, 5)
	if len(itemEntities) > 0 {
		platformAccountId2SocialProfileMap := map[int64]coredomain.Profile{}
		var platformAccountIds []string
		for i := range itemEntities {
			platformAccountIds = append(platformAccountIds, strconv.FormatInt(itemEntities[i].PlatformAccountCode, 10))
		}
		profileSearchQuery := coredomain.SearchQuery{
			Filters: []coredomain.SearchFilter{},
		}
		profileSearchQuery.Filters = append(profileSearchQuery.Filters, coredomain.SearchFilter{
			FilterType: "IN",
			Field:      "id",
			Value:      strings.Join(platformAccountIds, ","),
		})
		profiles, _, _ := socialDiscoveryManager.SearchByPlatform(ctx, platform, profileSearchQuery, "followers", "DESC", 1, 5, string(constants.SaasCollection), true, true)
		if profiles != nil && len(*profiles) > 0 {
			for i := range *profiles {
				socialProfile := (*profiles)[i]
				platformId, _ := strconv.ParseInt(socialProfile.PlatformCode, 10, 64)
				platformAccountId2SocialProfileMap[platformId] = socialProfile
			}
		}
		var items []domain.ProfileCollectionItemEntry
		for _, entity := range itemEntities {
			itemEntry, _ := ToItemEntry(&entity)
			if itemEntities != nil {
				sp, ok := platformAccountId2SocialProfileMap[*itemEntry.PlatformAccountCode]
				if ok {
					itemEntry.Profile = &sp
					items = append(items, *itemEntry)
				}
			}
		}

		for _, item := range items {
			if platform == string(constants.InstagramPlatform) {
				compactInstagramAccount := domain.CompactSocialAccount{
					Code:      item.Profile.Code,
					Name:      item.Profile.Name,
					Thumbnail: item.Profile.Thumbnail,
					Link:      getInstagramLink(*item.Profile),
				}
				topFive = append(topFive, compactInstagramAccount)
			} else if platform == string(constants.YoutubePlatform) {
				compactYTAccount := domain.CompactSocialAccount{
					Code:      item.Profile.Code,
					Name:      item.Profile.Name,
					Thumbnail: item.Profile.Thumbnail,
					Link:      getYoutubeLink(*item.Profile),
				}
				topFive = append(topFive, compactYTAccount)
			}
		}
	}
	return topFive
}

func getInstagramLink(profile coredomain.Profile) *string {
	var instaLink string
	if profile.Handle != nil {
		instaLink = "https://instagram.com/" + *profile.Handle
	}
	return &instaLink
}

func getYoutubeLink(profile coredomain.Profile) *string {
	var ytLink string
	if profile.Handle != nil {
		ytLink = "https://youtube.com/channel/" + *profile.Handle
	}
	return &ytLink
}

func (m *Manager) EnrichCollectionWithCustomColumnComputedValues(ctx context.Context, entry domain.ProfileCollectionEntry, filterHiddenColumns bool) *domain.ProfileCollectionEntry {
	platformColumnKeySums, err := m.itemDAO.FetchColumnSumValuesForCollectionItems(ctx, *entry.Id)
	if err == nil {
		platformColumnKeyToSumMap := map[string]dao.CollectionItemPlatformCustomColumnSum{}
		for i := range platformColumnKeySums {
			key := platformColumnKeySums[i].Platform + ":" + platformColumnKeySums[i].Key
			platformColumnKeyToSumMap[key] = platformColumnKeySums[i]
		}

		for i := range entry.CustomColumns {
			platform := entry.CustomColumns[i].Platform
			itemCountsSumary := entry.PlatformWiseItemsSummary[platform]
			// if collection has > 0 items for this platform.. then compute value for each column
			if itemCountsSumary.ProfileCounter.Count > 0 {
				var finalPlatformCustomColumnsMeta []domain.CustomColumnMeta
				for j := range entry.CustomColumns[i].Columns {
					if filterHiddenColumns && entry.CustomColumns[i].Columns[j].Hidden != nil && *entry.CustomColumns[i].Columns[j].Hidden {
						continue
					}
					if entry.CustomColumns[i].Columns[j].Compute != nil {
						totalItems := itemCountsSumary.ProfileCounter.Count
						remainingItems := totalItems
						allItemValuesSum := float64(0.0)
						itemsColumnDataSumEntity, ok := platformColumnKeyToSumMap[platform+":"+*entry.CustomColumns[i].Columns[j].Key]
						if ok {
							remainingItems = totalItems - itemsColumnDataSumEntity.Items
							allItemValuesSum = itemsColumnDataSumEntity.Sum
						}
						if strings.Trim(*entry.CustomColumns[i].Columns[j].Value, "") == "" {
							defVal := "0"
							entry.CustomColumns[i].Columns[j].Value = &defVal
						}
						defaultFillValue, err1 := strconv.ParseInt(*entry.CustomColumns[i].Columns[j].Value, 10, 64)
						if err1 != nil {
							errMsg := "Please check all values must be number"
							entry.CustomColumns[i].Columns[j].ComputedError = &errMsg
						} else {
							emptyErr := ""
							entry.CustomColumns[i].Columns[j].ComputedError = &emptyErr
							allItemValuesSum += float64(remainingItems * defaultFillValue)
							if strings.ToLower(*entry.CustomColumns[i].Columns[j].Compute) == "sum" {
								entry.CustomColumns[i].Columns[j].ComputedValue = fmt.Sprintf("%.2f", allItemValuesSum)
							} else if strings.ToLower(*entry.CustomColumns[i].Columns[j].Compute) == "average" {
								if defaultFillValue > 0 {
									avgVal := allItemValuesSum / float64(totalItems)
									entry.CustomColumns[i].Columns[j].ComputedValue = fmt.Sprintf("%.2f", avgVal)
								} else {
									avgVal := itemsColumnDataSumEntity.Sum / float64(itemsColumnDataSumEntity.Items)
									entry.CustomColumns[i].Columns[j].ComputedValue = fmt.Sprintf("%.2f", avgVal)
								}
							}
						}
					}
					finalPlatformCustomColumnsMeta = append(finalPlatformCustomColumnsMeta, entry.CustomColumns[i].Columns[j])
				}
				entry.CustomColumns[i].Columns = finalPlatformCustomColumnsMeta
			}
		}
	} else {
		for i := range entry.CustomColumns {
			var finalCustomColumns []domain.CustomColumnMeta
			for j := range entry.CustomColumns[i].Columns {
				if filterHiddenColumns && entry.CustomColumns[i].Columns[j].Hidden != nil && *entry.CustomColumns[i].Columns[j].Hidden {
					continue
				}
				errMsg := "Please check all values must be number"
				entry.CustomColumns[i].Columns[j].ComputedError = &errMsg
				finalCustomColumns = append(finalCustomColumns, entry.CustomColumns[i].Columns[j])
			}
			entry.CustomColumns[i].Columns = finalCustomColumns
		}
		appCtx := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)
		appCtx.Session.Rollback(ctx) // TODO: Figure out a better way to reset session. Need to do this because previous session has been marked as rollback-only
		appCtx.Session = nil
	}
	return &entry
}

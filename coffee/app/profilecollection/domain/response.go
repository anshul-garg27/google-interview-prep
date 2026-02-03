package domain

import (
	coredomain "coffee/core/domain"
	"strconv"
)

// -------------------------------
// profile collection service
// -------------------------------

type ProfileCollectionResponse struct {
	Status      coredomain.Status         `json:"status"`
	Collections *[]ProfileCollectionEntry `json:"collections,omitempty"`
}

func CreateResponse(collections []ProfileCollectionEntry, filterCount int64, nextCursor string, message string) ProfileCollectionResponse {
	return ProfileCollectionResponse{
		Status: coredomain.Status{
			Status:        "SUCCESS",
			Message:       message,
			FilteredCount: filterCount,
			HasNextPage:   !(nextCursor == ""),
			NextCursor:    nextCursor,
		},
		Collections: &collections,
	}
}
func CreateErrorResponse(err error, code int) ProfileCollectionResponse {
	return ProfileCollectionResponse{
		Status: coredomain.Status{
			Status:        "ERROR",
			Message:       err.Error(),
			FilteredCount: 0,
			HasNextPage:   false,
			NextCursor:    "",
		},
		Collections: nil,
	}
}

// -------------------------------
// profile collection item service
// -------------------------------

type ProfileCollectionItemResponse struct {
	Status                coredomain.Status             `json:"status"`
	CollectionItems       *[]ProfileCollectionItemEntry `json:"collectionItems"`
	FailedCollectionItems []ProfileCollectionItemEntry  `json:"failedItems"`
}

func CreateItemResponse(collectionItems []ProfileCollectionItemEntry, filteredCount int64, page int, size int, err error) ProfileCollectionItemResponse {
	if err == nil {
		var nextCursor string
		hasNextPage := false
		if len(collectionItems) >= size {
			nextCursor = strconv.Itoa(page + 1)
			hasNextPage = true
		} else {
			nextCursor = ""
		}
		return ProfileCollectionItemResponse{
			Status: coredomain.Status{
				Status:        "SUCCESS",
				Message:       "Collection Items Retrieved Successfully ",
				FilteredCount: filteredCount,
				HasNextPage:   hasNextPage,
				NextCursor:    nextCursor,
			},
			CollectionItems: &collectionItems,
		}
	} else {
		return CreateItemErrorResponse(err, 1001)
	}
}

func CreateItemErrorResponse(err error, code int) ProfileCollectionItemResponse {
	return ProfileCollectionItemResponse{
		Status: coredomain.Status{
			Status:        "500",
			Message:       err.Error(),
			FilteredCount: 0,
			HasNextPage:   false,
			NextCursor:    "",
		},
		CollectionItems: nil,
	}
}

func CreateItemDeleteResponse(err error) ProfileCollectionItemResponse {
	if err != nil {
		return CreateItemErrorResponse(err, 500)
	}
	return ProfileCollectionItemResponse{
		Status: coredomain.Status{
			Status:        "SUCCESS",
			Message:       "Collection Items Deleted Successfully ",
			FilteredCount: 0,
			HasNextPage:   false,
			NextCursor:    "",
		},
		CollectionItems: nil,
	}
}

func CreateItemAddResponse(items *[]ProfileCollectionItemEntry, err error, failedItems []ProfileCollectionItemEntry) ProfileCollectionItemResponse {
	if err != nil {
		return ProfileCollectionItemResponse{
			Status:          coredomain.Status{Status: "500", Message: err.Error(), FilteredCount: 0},
			CollectionItems: nil,
		}
	}
	return ProfileCollectionItemResponse{
		Status: coredomain.Status{
			Status:        "SUCCESS",
			Message:       "Collection Items Added Successfully ",
			FilteredCount: 0,
			HasNextPage:   false,
			NextCursor:    "",
		},
		CollectionItems:       items,
		FailedCollectionItems: failedItems,
	}
}

func CreateItemUpdateResponse(err error) ProfileCollectionItemResponse {
	if err != nil {
		return ProfileCollectionItemResponse{
			Status:          coredomain.Status{Status: "500", Message: err.Error(), FilteredCount: 0},
			CollectionItems: nil,
		}
	}
	return ProfileCollectionItemResponse{
		Status: coredomain.Status{
			Status:        "SUCCESS",
			Message:       "Collection Items Updated Successfully ",
			FilteredCount: 0,
			HasNextPage:   false,
			NextCursor:    "",
		},
		CollectionItems: nil,
	}
}

// ----- Wink collection info response
type WinklCollectionInfoResponse struct {
	Status         coredomain.Status         `json:"status"`
	CollectionInfo *WinklCollectionInfoEntry `json:"collectionInfo"`
}

func CreateCollectionMigrationInfoResponse(collectionInfo *WinklCollectionInfoEntry, err error) WinklCollectionInfoResponse {
	if err == nil {
		return WinklCollectionInfoResponse{
			Status: coredomain.Status{
				Status:        "SUCCESS",
				Message:       "Retrieve info successfully",
				FilteredCount: 0,
				HasNextPage:   false,
				NextCursor:    "",
			},
			CollectionInfo: collectionInfo,
		}
	} else {
		return WinklCollectionInfoResponse{
			Status: coredomain.Status{
				Status:        "500",
				Message:       "Failed to retrieve info",
				FilteredCount: 0,
				HasNextPage:   false,
				NextCursor:    "",
			},
			CollectionInfo: nil,
		}
	}
}

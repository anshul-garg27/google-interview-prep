package domain

import (
	coredomain "coffee/core/domain"
)

// -------------------------------
// post collection service
// -------------------------------

type PostCollectionResponse struct {
	Status      coredomain.Status      `json:"status"`
	Collections *[]PostCollectionEntry `json:"collections"`
}

// -------------------------------
// post collection item service
// -------------------------------

type PostCollectionItemResponse struct {
	Status          coredomain.Status          `json:"status"`
	CollectionItems *[]PostCollectionItemEntry `json:"collectionItems"`
}

package domain

import (
	coredomain "coffee/core/domain"
)

type KeywordCollectionResponse struct {
	Status      coredomain.Status         `json:"status"`
	Collections *[]KeywordCollectionEntry `json:"collections"`
}

type KeywordCollectionPostResponse struct {
	Status coredomain.Status        `json:"status"`
	Posts  *[]KeywordCollectionPost `json:"posts"`
}

type KeywordCollectionProfileResponse struct {
	Status   coredomain.Status           `json:"status"`
	Profiles *[]KeywordCollectionProfile `json:"profiles"`
}

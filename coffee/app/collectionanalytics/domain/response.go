package domain

import coredomain "coffee/core/domain"

// ------------------------------------------------------
//         Collection Metrics Summary Response
// ------------------------------------------------------

type CollectionMetricsSummaryResponse struct {
	Status             coredomain.Status              `json:"status"`
	MetricsSummaryData *CollectionMetricsSummaryEntry `json:"data"`
}

type CollectionTimeSeriesResponse struct {
	Status         coredomain.Status          `json:"status"`
	TimeSeriesData *CollectionTimeSeriesEntry `json:"data"`
}

// ------------------------------------------------------
//         Collection Profile List Metrics Summary Response
// ------------------------------------------------------

type CollectionProfileMetricsSummaryResponse struct {
	Status   coredomain.Status                      `json:"status"`
	Profiles []CollectionProfileMetricsSummaryEntry `json:"profiles"`
}

// ------------------------------------------------------
//         Collection Profile List Metrics Summary Response
// ------------------------------------------------------

type CollectionPostMetricsSummaryResponse struct {
	Status coredomain.Status                   `json:"status"`
	Posts  []CollectionPostMetricsSummaryEntry `json:"posts"`
}

// ------------------------------------------------------
//         Collection Hashtags Response
// ------------------------------------------------------

type CollectionHashtagResponse struct {
	Status   coredomain.Status        `json:"status"`
	Hashtags []CollectionHashtagEntry `json:"list"`
}

// ------------------------------------------------------
//         Collection Keywords Response
// ------------------------------------------------------

type CollectionKeywordResponse struct {
	Status   coredomain.Status        `json:"status"`
	Keywords []CollectionKeywordEntry `json:"list"`
}

// ------------------------------------------------------
//
//	Collection Sentiment Response
//
// ------------------------------------------------------
type CollectionSentimentResponse struct {
	Status        coredomain.Status         `json:"status"`
	SentimentData *CollectionSentimentEntry `json:"data"`
}

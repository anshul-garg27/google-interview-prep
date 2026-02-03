package coffeeclient

type SentimentReportEntry struct {
	CollectionType        string `json:"collectionType,omitempty"`
	CollectionId          string `json:"collectionId,omitempty"`
	SentimentReportPath   string `json:"sentimentReportPath,omitempty"`
	SentimentReportBucket string `json:"sentimentReportBucket,omitempty"`
}

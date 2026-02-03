package dao

import (
	"coffee/core/persistence/clickhouse"
	"context"
	"fmt"
)

type MetricsDao struct {
	ctx         context.Context
	DaoProvider *clickhouse.ClickhouseDao
}

func (d *MetricsDao) init(ctx context.Context) {
	d.DaoProvider = clickhouse.NewClickhouseDao(ctx)
	d.ctx = ctx
}

func CreatePostCollectionMetricsDao(ctx context.Context) *MetricsDao {
	dao := &MetricsDao{}
	dao.init(ctx)
	return dao
}

func (d *MetricsDao) getS3Url(path string, bucket string) string {
	path = fmt.Sprintf("https://%s.s3.ap-south-1.amazonaws.com/%s", bucket, path)
	return path
}

func (d *MetricsDao) GetSentimentCount(ctx context.Context, bucket string, path string) ([]SentimentCount, error) {
	path = d.getS3Url(path, bucket)
	access_key := "AKIAXGXUCIERSOUELY73"
	access_secret := "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"
	db := d.DaoProvider.GetSession(ctx)
	var sentimentCount []SentimentCount
	whereArgs := map[string]interface{}{}
	query := fmt.Sprintf(`
			SELECT
				count(*) count, sentiment sentiment
			FROM s3('%s', '%s', '%s', 'Parquet')
			GROUP BY sentiment
	`, path, access_key, access_secret)
	db.Raw(query, whereArgs).Scan(&sentimentCount)
	return sentimentCount, nil
}

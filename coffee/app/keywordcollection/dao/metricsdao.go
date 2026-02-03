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

func CreateKeywordCollectionMetricsDao(ctx context.Context) *MetricsDao {
	dao := &MetricsDao{}
	dao.init(ctx)
	return dao
}

func (d *MetricsDao) GetSummary(ctx context.Context, bucket string, path string) (*Summary, error) {
	path = d.getS3Url(path, bucket)
	access_key := "AKIAXGXUCIER5DX7YI4F"
	access_secret := "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"
	db := d.DaoProvider.GetSession(ctx)
	summary := &Summary{}
	whereArgs := map[string]interface{}{}
	query := fmt.Sprintf(`
			SELECT
				uniqExact(shortcode) posts,
				toInt64(sumIf(post_reach, post_reach > 0 and not isInfinite(post_reach))) as reach,
				toInt64(sum(post_likes)) as likes,
				toInt64(sum(post_comments)) as comments,
				toInt64(likes + comments) as engagement,
				toInt64(round(engagement/(reach/1000))) as emp,
				round(engagement/reach,2) as engagement_rate
			FROM s3('%s', '%s', '%s', 'Parquet', 'shortcode String, post_reach Nullable(Float64), post_likes Nullable(Int64), post_comments Nullable(Int64)')
	`, path, access_key, access_secret)
	db.Raw(query, whereArgs).Scan(summary)
	return summary, nil
}

func (d *MetricsDao) GetDemography(ctx context.Context, bucket string, path string) (*Demography, error) {
	path = d.getS3Url(path, bucket)
	access_key := "AKIAXGXUCIER5DX7YI4F"
	access_secret := "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"
	db := d.DaoProvider.GetSession(ctx)
	demography := &Demography{}
	whereArgs := map[string]interface{}{}
	query := fmt.Sprintf(`
			SELECT
				toString(audience_age) audience_age,
				toString(audience_gender) audience_gender,
				toString(audience_age_gender) audience_age_gender
			FROM s3('%s', '%s', '%s', 'Parquet', 'audience_age_gender String, audience_age String, audience_gender String')
	`, path, access_key, access_secret)
	db.Raw(query, whereArgs).Scan(demography)
	return demography, nil
}

func (d *MetricsDao) GetPosts(ctx context.Context, bucket string, path string, sortBy string, sortDesc string) ([]Post, error) {
	path = d.getS3Url(path, bucket)
	access_key := "AKIAXGXUCIER5DX7YI4F"
	access_secret := "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"
	db := d.DaoProvider.GetSession(ctx)
	var posts []Post
	whereArgs := map[string]interface{}{}
	query := fmt.Sprintf(`
			SELECT
				shortcode,
				max(post_thumbnail_url) post_thumbnail_url,
				max(profile_pic_url) as profile_thumbnail_url,
				toInt64(round(sumIf(post_reach, post_reach > 0 and not isInfinite(post_reach)))) reach,
				fromUnixTimestamp(max(post_publish_time)) post_publish_time,
				toInt64(round(sum(post_views))) views,
				toInt64(sum(post_likes)) likes,
				max(post_type) post_type,
				toInt64(sum(post_comments)) comments,
				toInt64(max(profile_followers)) followers,
				max(category) category,
				round((comments + likes) / followers,2) engagement_rate,
				max(handle) handle,
				max(profile_id) profile_id,
				max(profile_title) name,
				max(post_title) post_title
			FROM s3('%s', '%s', '%s', 'Parquet', 'shortcode Nullable(String), post_reach Nullable(Float64), post_likes Nullable(Int64), post_comments Nullable(Int64), post_title Nullable(String), profile_title Nullable(String), post_thumbnail_url Nullable(String), post_views Nullable(Int64), category Nullable(String), post_publish_time Nullable(Int64), profile_pic_url Nullable(String), handle Nullable(String), profile_followers Nullable(Int64), post_type Nullable(String),profile_id Nullable(String)')
			GROUP BY shortcode
			ORDER BY %s %s
			LIMIT 12
	`, path, access_key, access_secret, sortBy, sortDesc)
	db.Raw(query, whereArgs).Scan(&posts)
	return posts, nil
}

func (d *MetricsDao) GetProfiles(ctx context.Context, bucket string, path string, sortBy string, sortDesc string) ([]Profile, error) {
	path = d.getS3Url(path, bucket)
	access_key := "AKIAXGXUCIER5DX7YI4F"
	access_secret := "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"
	db := d.DaoProvider.GetSession(ctx)
	var profiles []Profile
	whereArgs := map[string]interface{}{}
	query := fmt.Sprintf(`
			SELECT
					profile_id, 
        			max(handle) as handle,
					max(profile_title) as name,
					max(profile_pic_url) as thumbnail,
					max(profile_followers) as followers,
					max(category) as category,
					toInt64(uniqExact(shortcode)) uploads,
					toInt64(round(sumIf(post_reach, post_reach > 0 and not isInfinite(post_reach)))) as reach,
					toInt64(round(sum(post_views))) as views,
					toInt64(sum(post_likes)) as likes,
					toInt64(sum(post_comments)) as comments,
					fromUnixTimestamp(max(post_publish_time)) as last_posted_on,
					max(engagement_rate) engagement_rate
			FROM s3('%s', '%s', '%s', 'Parquet', 'shortcode Nullable(String), post_reach Nullable(Float64), post_likes Nullable(Int64), post_comments Nullable(Int64), profile_title Nullable(String), post_views Nullable(Int64), category Nullable(String), post_publish_time Nullable(Int64), profile_pic_url Nullable(String),handle Nullable(String), profile_followers Nullable(Int64), profile_id Nullable(String) , engagement_rate	Nullable(Float64)')
			GROUP BY profile_id
			ORDER BY %s %s
			LIMIT 12
	`, path, access_key, access_secret, sortBy, sortDesc)
	db.Raw(query, whereArgs).Scan(&profiles)
	return profiles, nil
}

func (d *MetricsDao) getS3Url(path string, bucket string) string {
	path = fmt.Sprintf("https://%s.s3.ap-south-1.amazonaws.com/%s", bucket, path)
	return path
}

func (d *MetricsDao) GetCategories(ctx context.Context, bucket string, path string) ([]Category, error) {
	path = d.getS3Url(path, bucket)
	access_key := "AKIAXGXUCIER5DX7YI4F"
	access_secret := "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"
	db := d.DaoProvider.GetSession(ctx)
	var categories []Category
	whereArgs := map[string]interface{}{}
	query := fmt.Sprintf(`
			SELECT
				category,
				toInt64(round(sumIf(post_reach, post_reach > 0 and not isInfinite(post_reach)))) reach,
				uniqExact(shortcode) posts,
				round(100*(posts/sum(posts) over ()), 2) posts_perc
			FROM s3('%s', '%s', '%s', 'Parquet', 'post_reach Nullable(Float64), shortcode Nullable(String)')
			where category != ''
			group by category
			order by posts desc, category desc
			LIMIT 5
	`, path, access_key, access_secret)
	db.Raw(query, whereArgs).Scan(&categories)
	return categories, nil
}

func (d *MetricsDao) GetMonthlyStats(ctx context.Context, bucket string, path string, rollup string) ([]MonthlyStat, error) {
	path = d.getS3Url(path, bucket)
	access_key := "AKIAXGXUCIER5DX7YI4F"
	access_secret := "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"
	db := d.DaoProvider.GetSession(ctx)
	var monthlyStats []MonthlyStat
	rollupQueryFn := "toStartOfMonth"
	if rollup == "weekly" {
		rollupQueryFn = "toStartOfWeek"
	}
	whereArgs := map[string]interface{}{}
	query := fmt.Sprintf(`
			SELECT
				formatDateTime(%s(fromUnixTimestamp(post_publish_time)), '%%F') month,
				toInt64(round(sumIf(post_reach, post_reach > 0 and not isInfinite(post_reach)))) views,
				toInt64(uniqExact(shortcode)) uploads,
				round(100*(uploads/sum(uploads) over ()), 2) uploads_perc
			FROM s3('%s', '%s', '%s', 'Parquet', 'post_publish_time Nullable(Int64), post_reach Nullable(Float64), shortcode Nullable(String)')
			group by month
			order by month asc
	`, rollupQueryFn, path, access_key, access_secret)
	db.Raw(query, whereArgs).Scan(&monthlyStats)
	return monthlyStats, nil
}

func (d *MetricsDao) GetKeywords(ctx context.Context, bucket string, path string) ([]Keywords, error) {
	path = d.getS3Url(path, bucket)
	access_key := "AKIAXGXUCIER5DX7YI4F"
	access_secret := "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"
	db := d.DaoProvider.GetSession(ctx)
	var keywords []Keywords
	whereArgs := map[string]interface{}{}
	query := fmt.Sprintf(`
			with data as (
				SELECT
					lower(replaceAll(arrayJoin(JSONExtractArrayRaw(replaceAll(ifNull(post_keywords, ''), '\'', '"'))), '"', '')) keyword,
					toInt64(round(sumIf(post_reach, post_reach > 0 and not isInfinite(post_reach)))) reach,
					toInt64(uniqExact(shortcode)) posts,
					round(100*(posts/sum(posts) over ()), 2) posts_perc
				FROM s3('%s', '%s', '%s', 'Parquet', 'post_keywords Nullable(String), post_reach Nullable(Float64), shortcode Nullable(String)')
				where keyword != ''
				group by keyword
				order by posts desc, keyword desc
			),
			stats as (
				select 
					uniqExact(keyword) total_keywords,
					max(posts) total_posts
				from data
			),
			final as (
				select
					keyword,
					reach,
					posts,
					posts_perc,
					posts / stats.total_keywords tf,
					if (posts = stats.total_posts, log(stats.total_posts / posts+0.5), log(stats.total_posts / posts)) idf,
					100 - 1/(tf*idf) score
				from data,stats
			)
			select * from final order by score desc limit 50
	`, path, access_key, access_secret)
	db.Raw(query, whereArgs).Scan(&keywords)
	return keywords, nil
}
func (d *MetricsDao) GetKeywordsFromTitle(ctx context.Context, bucket string, path string) ([]Keywords, error) {
	path = d.getS3Url(path, bucket)
	access_key := "AKIAXGXUCIER5DX7YI4F"
	access_secret := "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"
	db := d.DaoProvider.GetSession(ctx)
	var keywords []Keywords
	whereArgs := map[string]interface{}{}
	query := fmt.Sprintf(`
			with data as (
					SELECT
						lower(arrayJoin(splitByWhitespace(ifNull(post_title, '')))) keyword,
						toInt64(round(sumIf(post_reach, post_reach > 0 and not isInfinite(post_reach)))) reach,
						toInt64(uniqExact(shortcode)) posts,
						round(100*(posts/sum(posts) over ()), 2) posts_perc
					FROM s3('%s', '%s', '%s', 'Parquet', 'post_title Nullable(String), post_reach Nullable(Float64), shortcode Nullable(String)')
					where keyword != '' and 
					length(keyword) > 4 and 
					length(keyword) <= 30 and 
					lower(keyword) not like 'http%%' and
					lower(keyword) not IN 
							('your', 'its', 'it\'s', 'from', 'this', 'stop', 'with', 'that', 'have', 'https', 'using', 
							'also', 'which', 'these', 'will', 'what', 'this', 'about',
								'more', 'their', 'they', 'used', 'what', 'some', 'just', 'like', 'don\'t', 'dont')
					group by keyword
					having posts >= 3
					order by posts desc, keyword desc
					LIMIT 1000
			),
			stats as (
				select 
					uniqExact(keyword) total_keywords,
					max(posts) total_posts
				from data
			),
			final as (
				select
					keyword,
					reach,
					posts,
					posts_perc,
					posts / stats.total_keywords tf,
					log(stats.total_posts / posts) idf,
					tf*idf tfidf,
					max(tfidf) over () max_tfidf,
					min(tfidf) over () min_tfidf,
					(tfidf - min_tfidf) / (max_tfidf - min_tfidf) tfidf_n,
					100 - 1/(tfidf_n) score_old,
					((tfidf - min_tfidf) / (max_tfidf - min_tfidf))  * (95 - 40) + 40 score
				from data,stats
			)
			select * from final order by score desc limit 50
	`, path, access_key, access_secret)
	db.Raw(query, whereArgs).Scan(&keywords)
	return keywords, nil
}

func (d *MetricsDao) GetLanguages(ctx context.Context, bucket string, path string) ([]Language, error) {
	path = d.getS3Url(path, bucket)
	access_key := "AKIAXGXUCIER5DX7YI4F"
	access_secret := "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"
	db := d.DaoProvider.GetSession(ctx)
	var languages []Language
	whereArgs := map[string]interface{}{}
	query := fmt.Sprintf(`
			SELECT
				splitByChar('-', ifNull(post_language, ''))[1] language,
				toInt64(round(sumIf(post_reach, post_reach > 0 and not isInfinite(post_reach)))) reach,
				toInt64(uniqExact(shortcode)) posts,
				round(100*(posts/sum(posts) over ()), 2) posts_perc
			FROM s3('%s', '%s', '%s', 'Parquet', 'post_language Nullable(String), post_reach Nullable(Float64), shortcode Nullable(String)')
			where language != ''
			group by language
			order by posts desc, language desc
	`, path, access_key, access_secret)
	db.Raw(query, whereArgs).Scan(&languages)
	return languages, nil
}

func (d *MetricsDao) GetLanguagesMonthly(ctx context.Context, bucket string, path string) ([]Language, error) {
	path = d.getS3Url(path, bucket)
	access_key := "AKIAXGXUCIER5DX7YI4F"
	access_secret := "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"
	db := d.DaoProvider.GetSession(ctx)
	var languages []Language
	whereArgs := map[string]interface{}{}
	query := fmt.Sprintf(`
			SELECT
				formatDateTime(toStartOfMonth(fromUnixTimestamp(post_publish_time)), '%%F') month,
				splitByChar('-', ifNull(post_language, ''))[1] language,
				toInt64(round(sumIf(post_reach, post_reach > 0 and not isInfinite(post_reach)))) reach,
				toInt64(uniqExact(shortcode)) posts,
				round(100*(posts/sum(posts) over (partition by month)), 2) posts_perc
			FROM s3('%s', '%s', '%s', 'Parquet', 'post_publish_time Nullable(Int64), post_language Nullable(String), post_reach Nullable(Float64), shortcode Nullable(String)')
			where language != ''
			group by month, language
			order by month desc, posts desc, language desc
	`, path, access_key, access_secret)
	db.Raw(query, whereArgs).Scan(&languages)
	return languages, nil
}

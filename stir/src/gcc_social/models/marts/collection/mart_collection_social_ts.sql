{{ config(materialized = 'table', tags=["collections"], order_by='collection_id, stats_date') }}
with
social_ts as (
	select
		platform,
		ifNull(collection_id, 'MISSING') collection_id,
		collection_type,
		ifNull(stats_date, date('2000-01-01')) stats_date,
		uniqExact(post_short_code) posts,
		sum(views) views,
		sum(likes) likes,
		sum(comments) comments,
		sum(impressions) impressions,
		sum(reach) reach,
		sum(total_engagement) total_engagement
    from {{ref('mart_collection_post_ts')}}
	group by platform, collection_id, collection_type, stats_date
	order by platform asc, collection_id asc, stats_date asc
)
select * from social_ts
SETTINGS max_bytes_before_external_group_by = 20000000000, max_bytes_before_external_sort = 20000000000, join_algorithm='partial_merge'

{{ config(materialized = 'table', tags=["collections"], order_by='collection_id, stats_date') }}
with
collection_posts as (
    select
    	platform,
    	post_collection_id collection_id,
    	'POST' collection_type,
    	posted_by_handle profile_handle,
        short_code post_short_code,
    	replace(JSONExtractArrayRaw(ifNull(sponsor_links, '[]'))[1], '"', '') link
    from {{ref('stg_coffee_post_collection_item')}}
    where enabled = True and show_in_report = True and link is not null and link != ''
),
links as (
    select link from collection_posts
    group by link
),
link_clicks as (
    select
    	concat('https://', url) link,
    	toStartOfDay(d) stats_date,
    	sum(unique_clicks) clicks
    from dbt.mart_link_clicks
    where link in links
    group by link, stats_date
),
clicks_ts as (
	select
		cp.platform platform,
		ifNull(cp.collection_id, 'MISSING') collection_id,
		cp.collection_type collection_type,
		cp.profile_handle profile_handle,
		cp.post_short_code post_short_code,
		ifNull(lc.stats_date, date(2000-01-01)) stats_date,
		sum(lc.clicks) clicks,
		sum(clicks) over (partition by platform, collection_id order by stats_date asc) clicks_cumul
	from collection_posts cp
	left join link_clicks lc on lc.link = cp.link
	where stats_date >= '2020-01-01' and stats_date <= now()
	group by platform, collection_id, collection_type, profile_handle, post_short_code, stats_date
	order by platform asc, collection_id asc, stats_date asc
)
select * from clicks_ts
SETTINGS max_bytes_before_external_group_by = 20000000000, max_bytes_before_external_sort = 20000000000, join_algorithm='partial_merge'

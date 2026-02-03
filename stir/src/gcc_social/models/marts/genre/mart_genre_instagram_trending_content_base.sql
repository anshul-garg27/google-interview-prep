{{ config(materialized = 'table', tags=["daily", "genre"], order_by='month, short_code') }}
with
	base as (
		select
			max(ip.profile_id) profile_id,
			ifNull(ip.shortcode, 'MISSING') short_code,
			max(caption) title,
			max(caption) description,
			max(ip.post_type) post_type,
			max(COALESCE(mia.profile_type, 'CREATOR')) profile_type,
			max(ip.likes) likes,
			max(ip.comments) comments,
			max(ip.views) views,
			max(ip.plays) plays,
			max(ip.reach) reach,
			mia.primary_category category,
			ifNull(toStartOfMonth(ip.publish_time), date('2000-01-01')) month
		from dbt.stg_beat_instagram_post ip
		left join dbt.mart_instagram_account mia on ip.profile_id = mia.ig_id
		where mia.country = 'IN' and ip.post_type = 'reels' and ip.plays >0 and ip.likes >0
		group by month, short_code, category
	),
	base_with_rank as (
		select 
			*,
			rank() over(partition by month, category order by plays desc) plays_rank
		from base post
	),
	sorted as (
		SELECT
			*
		from base_with_rank
		WHERE (
			plays_rank < 501
		)
		ORDER BY plays_rank
	)
select * from sorted
SETTINGS allow_experimental_nlp_functions = 1, max_bytes_before_external_group_by = 40000000000, max_bytes_before_external_sort = 40000000000, join_algorithm='partial_merge'
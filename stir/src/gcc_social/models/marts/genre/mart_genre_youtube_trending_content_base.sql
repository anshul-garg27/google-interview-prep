{{ config(materialized = 'table', tags=["daily", "genre"], order_by='month, short_code') }}

with
	base as (
		select
            max(ifNull(yp.channel_id, 'MISSING')) profile_id,
            ifNull(yp.shortcode, 'MISSING') short_code,
			max(title) title,
			max(description) description,
            max(yp.post_type) post_type,
            max(COALESCE(mya.profile_type, 'CREATOR')) profile_type,
            max(yp.likes) likes,
            max(yp.comments) comments,
            max(yp.views) views,
			mya.category category,
            ifNull(toStartOfMonth(yp.published_at), date('2000-01-01')) month
        from dbt.stg_beat_youtube_post yp
        left join dbt.mart_youtube_account mya on yp.channel_id = mya.channel_id
        where mya.country = 'IN' and yp.views>0
        group by month, short_code, category
	),
	base_with_rank as (
		select 
            *,
            rank() over(partition by month, category order by views desc) views_rank,
            0 plays_rank
        from base post
	),
	sorted as (
		SELECT
			*
		from base_with_rank
		WHERE (
			views_rank < 501
		)
		ORDER BY views_rank
	)
select * from sorted
SETTINGS allow_experimental_nlp_functions = 1, max_bytes_before_external_group_by = 40000000000, max_bytes_before_external_sort = 40000000000, join_algorithm='partial_merge'

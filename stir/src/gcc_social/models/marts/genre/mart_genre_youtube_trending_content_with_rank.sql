{{ config(materialized = 'table', tags=["daily", "genre"], order_by='month, short_code') }}
WITH
short_code as (
    select short_code
    from {{ref('mart_genre_youtube_trending_content_base')}}
),
posts as (
	select 
		sbip.shortcode ,
		max(sbip.post_type) post_type ,
		max(sbip.thumbnail) thumbnail,
		max(sbip.published_at) published_at
	from dbt.stg_beat_youtube_post sbip 
	where sbip.shortcode in short_code
    group by sbip.shortcode
),
profile_ids as (
    select profile_id
    from {{ref('mart_genre_youtube_trending_content_base')}}
),
profiles as (
    select
        channelid profile_id,
        category category,
        language
        from dbt.mart_youtube_account mya
    where mya.channelid in profile_ids
),
final as (
    select 
        'YOUTUBE' platform,
        base.month month,
        base.profile_id profile_id,
        base.short_code short_code,
        base.title title,
        base.description description,
        base.views views,
        base.likes likes,
        base.comments comments,
        likes + comments engagement,
        (engagement) / (1.0*views) engagement_rate,
        0 plays,
        base.views_rank views_rank,
        base.plays_rank plays_rank,
        post.post_type post_type,
        base.profile_type profile_type,
        if(post_type = 'short', concat('https://www.youtube.com/shorts/' , short_code), concat('https://www.youtube.com/watch?v=', short_code)) post_url,
        COALESCE(
            NULLIF(REPLACE(post.thumbnail, '/default.jpg', '/hqdefault.jpg'), post.thumbnail),
            NULLIF(REPLACE(post.thumbnail, '/default_live.jpg', '/hqdefault_live.jpg'), post.thumbnail),
            post.thumbnail
        ) AS thumbnail_url,
        post.published_at publish_time,
        IF(profile.language='', NULL, profile.language) language,
        IF(base.category='', NULL, base.category) category
    from {{ref('mart_genre_youtube_trending_content_base')}} base
    left join posts post on base.short_code = post.shortcode
    left join profiles profile on base.profile_id = profile.profile_id
)
select * from final
SETTINGS allow_experimental_nlp_functions = 1, max_bytes_before_external_group_by = 40000000000, max_bytes_before_external_sort = 40000000000, join_algorithm='partial_merge'
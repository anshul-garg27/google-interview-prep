{{ config(materialized = 'table', tags=["daily", "genre"], order_by='month, short_code') }}
WITH
short_code as (
    select short_code
    from {{ref('mart_genre_instagram_trending_content_base')}}
),
posts as (
	select 
		sbip.shortcode ,
		max(sbip.post_type) post_type ,
		max(COALESCE(concat('https://d24w28i6lzk071.cloudfront.net/', sbal.asset_url) , sbip.thumbnail_url)) thumbnail_url,
		max(sbip.publish_time) publish_time
	from dbt.stg_beat_instagram_post sbip
	left join dbt.stg_beat_asset_log sbal on sbip.shortcode = sbal.entity_id
	where sbip.shortcode in short_code
	group by sbip.shortcode 
),
profile_ids as (
	select profile_id
	from {{ref('mart_genre_instagram_trending_content_base')}}
),
profiles as (
	select
		ig_id profile_id,
		followers,
		primary_language language,
		primary_category category
		from dbt.mart_instagram_account mia 
	WHERE mia.ig_id in profile_ids and followers>0
),
final as (
    select 
        'INSTAGRAM' platform,
        base.month month,
        base.profile_id profile_id,
        base.short_code short_code,
        base.title title,
        base.description description,
        base.likes likes,
        base.comments comments,
        likes + comments engagement,
        toInt64(round(plays * (0.94 - (log2(profile.followers) * 0.001)))) _reach_reels,
        toInt64(round((7.6 - (log10(likes) * 0.7)) * 0.85*likes)) _reach_non_reels,
        if (post_type = 'reels', max2(_reach_reels,0), max2(_reach_non_reels, 0)) _reach,
        if (base.reach is not null and base.reach > 0, toInt64(round(base.reach)), toInt64(round(_reach))) reach,
        reach views,
        1.0*engagement/reach engagement_rate,
        base.plays plays,
        0 views_rank,
        base.plays_rank plays_rank,
        post.post_type post_type,
        base.profile_type profile_type,
        concat('https://www.instagram.com/reel/', short_code) post_url,
        post.thumbnail_url thumbnail_url,
        post.publish_time publish_time,
        IF(profile.language='', NULL, profile.language) language,
        IF(profile.category='', NULL, profile.category) category
    from {{ref('mart_genre_instagram_trending_content_base')}} base
    left join posts post on base.short_code = post.shortcode 
    left join profiles profile on base.profile_id = profile.profile_id
    where profile.followers > 0
)
select * from final
SETTINGS allow_experimental_nlp_functions = 1, max_bytes_before_external_group_by = 40000000000, max_bytes_before_external_sort = 40000000000, join_algorithm='partial_merge'
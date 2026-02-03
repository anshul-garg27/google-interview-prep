{{ config(
        materialized = 'table',
        tags=["post_ranker_partial"],
        order_by='post_shortcode'
    )
}}
with
latest_handles as (
    select handle from
        {{ ref('stg_beat_instagram_account') }}
    where updated_at >= now() - INTERVAL 24 HOUR
),
latest_pids as (
    select profile_id from
        {{ ref('stg_beat_instagram_account') }}
    where updated_at >= now() - INTERVAL 24 HOUR
),
old_posts as (
    select
        profile_id profile_id,
        handle handle,
        ifNull(ip.shortcode, 'MISSING') post_shortcode,
        ip.post_type post_type,
        if(ip.post_type = 'reels', 'reels', 'static') post_class,
        ip.likes likes,
        ip.comments comments,
        ip.views views,
        ip.plays plays,
        likes + comments engagement,
        ifNull(ip.publish_time, date('2000-01-01')) - INTERVAL 330 MINUTE publish_time
    from vidooly.instagram_posts_historical_cleaned_v2 ip
    where handle in latest_handles or profile_id in latest_pids
),
new_posts as (
    select
        profile_id profile_id,
        handle handle,
        ifNull(ip.shortcode, 'MISSING') post_shortcode,
        ip.post_type post_type,
        if(ip.post_type = 'reels', 'reels', 'static') post_class,
        max(ip.likes) likes,
        max(ip.comments) comments,
        max(ip.views) views,
        max(ip.plays) plays,
        likes + comments engagement,
        max(ifNull(ip.publish_time, date('2000-01-01'))) publish_time
    from {{ ref('stg_beat_instagram_post') }} ip
    where post_type != 'story'
    and (handle in latest_handles or profile_id in latest_pids)
    group by profile_id, handle, post_shortcode, post_type, post_class
),
all_posts_merged as (
    select * from old_posts union all select * from new_posts
)
select * from all_posts_merged order by profile_id desc, publish_time desc
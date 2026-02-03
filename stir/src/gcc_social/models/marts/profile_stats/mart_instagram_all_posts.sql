{{ config(
        materialized = 'table',
        tags=["post_ranker"],
        order_by='post_shortcode'
    )
}}
with
old_posts as (
    select
        profile_id profile_id,
        handle handle,
        ifNull(ip.shortcode, 'MISSING') post_shortcode,
        if(ip.post_type = 'reels', 'reels', 'static') post_class,
        ip.post_type post_type,
        ip.likes likes,
        ip.comments comments,
        ip.views views,
        ip.plays plays,
        likes + comments engagement,
        ifNull(ip.publish_time, date('2000-01-01')) - INTERVAL 330 MINUTE publish_time,
        0 reach,
        0 impressions
    from vidooly.instagram_posts_historical_cleaned_v2 ip
),
new_posts as (
    select
        profile_id profile_id,
        handle handle,
        ifNull(ip.shortcode, 'MISSING') post_shortcode,
        if(ip.post_type = 'reels', 'reels', 'static') post_class,
        ip.post_type post_type,
        max(ip.likes) likes,
        max(ip.comments) comments,
        max(ip.views) views,
        max(ip.plays) plays,
        likes + comments engagement,
        max(ifNull(ip.publish_time, date('2000-01-01'))) publish_time,
        plays * (0.94 - (log2(max(a.followers)) * 0.001)) _reach_reels,
        toFloat64(plays) _impressions_reels,
        (7.6 - (log10(likes) * 0.7)) * 0.85*likes _reach_non_reels,
        (7.6-(log10(likes) * 0.75))*0.95*likes _impressions_non_reels,
        if (post_class = 'reels', max2(_reach_reels,0), max2(_reach_non_reels,0)) reach,
        if (post_class = 'reels', max2(_impressions_reels,0), max2(_impressions_non_reels,0)) impressions
    from {{ ref('stg_beat_instagram_post') }} ip
    left join {{ ref('stg_beat_instagram_account') }} a on a.profile_id = ip.profile_id
    where post_type != 'story'
    group by profile_id, handle, post_shortcode, post_class, post_type
),
all_posts_merged as (
    select
            profile_id,
            handle,
            post_shortcode,
            post_class,
            post_type,
            likes,
            comments,
            views,
            plays,
            engagement,
            publish_time,
            reach,
            impressions
    from old_posts
        union all
    select profile_id,
            handle,
            post_shortcode,
            post_class,
            post_type,
            likes,
            comments,
            views,
            plays,
            engagement,
            publish_time,
            reach,
            impressions
    from new_posts
)
select * from all_posts_merged order by profile_id desc, publish_time desc
SETTINGS max_bytes_before_external_group_by = 40000000000, max_bytes_before_external_sort = 40000000000
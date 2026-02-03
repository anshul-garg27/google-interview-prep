{{ config(materialized = 'table', tags=["staging_collections"], order_by='collection_id,post_short_code,stats_date') }}
with
collection_posts as (
    select id, platform, short_code, post_type, post_collection_id
    from {{ref('stg_coffee_staging_post_collection_item')}}
    where platform = 'YOUTUBE'
),
shortcodes as (
    select short_code from collection_posts
),
all_profiles as (
    select profile_id
    from dbt.stg_beat_youtube_post
    where shortcode in shortcodes
    group by profile_id
),
posts as (
    select
        channel_id profile_id,
        shortcode,
        channel_id handle,
        'video' post_type,
        description post_title,
        extractAll(ifNull(description, ''), '#[a-zA-Z0-9]+') hashtags,
        extractAll(ifNull(description, ''), '[a-zA-Z0-9]{4,}') keywords,
        published_at,
        concat('https://youtube.com/v/', shortcode) post_link,
        thumbnail post_thumbnail
    from dbt.stg_beat_youtube_post where shortcode in shortcodes
),
base_ts_fake as (
    select
        platform,
        post_short_code short_code,
        stats_date,
        toInt64(views) views,
        likes,
        comments,
        platform,
        stats_date,
        now() updated_at,
        post_short_code short_code,
        likes + comments total_engagement
    from {{ref('mart_staging_fake_events')}}
    where short_code in shortcodes and platform = 'YOUTUBE'
),
base_ts_raw as (
    select
        platform,
        toStartOfDay(pts.event_timestamp) stats_date,
        pts.event_timestamp updated_at,
        shortcode short_code,

        JSONExtractInt(metrics, 'comments') _comments,
        JSONExtractInt(metrics, 'likes') _likes,
        JSONExtractInt(metrics, 'view_count') _views,
        JSONExtractInt(metrics, 'play_count') _plays,

        JSONExtractString(metrics, 'comments') str_comments,
        JSONExtractString(metrics, 'likes') str_likes,
        JSONExtractString(metrics, 'view_count') str_views,
        JSONExtractString(metrics, 'play_count') str_plays,

        if(_comments = 0, toInt64OrZero(str_comments), _comments) comments,
        if(_likes = 0, toInt64OrZero(str_likes), _likes) likes,
        if(_views = 0, toInt64OrZero(str_views), _views) __views,
        if(_plays = 0, toInt64OrZero(str_plays), _plays) __plays,
        if (__views == 0, __plays, __views) views,
        likes + comments total_engagement
    from _e.post_log_events pts
    where shortcode in shortcodes and platform = 'YOUTUBE'
),
base_ts_all as (
    select
        platform,
        short_code,
        stats_date,
        updated_at,
        comments,
        likes,
        views,
        total_engagement
    from base_ts_raw
        union all
    select
        platform,
        short_code,
        stats_date,
        updated_at,
        comments,
        likes,
        views,
        total_engagement
    from base_ts_fake
),
base_ts as (
    select
        platform,
        short_code,
        stats_date,
        max(updated_at) updated_at,
        max(comments) comments,
        max(likes) likes,
        max(views) views,
        max(total_engagement) total_engagement
    from base_ts_all
    group by platform, short_code, stats_date
),
with_metrics as (
    select
        pts.short_code short_code,
        pts.stats_date stats_date,
        p.handle handle,
        p.post_type post_type,
        p.profile_id profile_id,
        p.post_title post_title,
        p.post_link post_link,
        p.post_thumbnail post_thumbnail,
        p.hashtags hashtags,
        p.keywords keywords,
        p.published_at published_at,
        pts.comments comments,
        0 mentions,
        pts.updated_at updated_at,
        0 followers,
        pts.total_engagement total_engagement,
        0 saves,
        0 shares,
        0 exits,
        0 sticker_taps,
        0 taps_forward,
        0 taps_back,
        0 swipe_ups,
        pts.likes likes,
        0 plays,
        1.0*views reach,
        1.0*views impressions,
        total_engagement / reach engagement_rate
    from base_ts pts
    left join posts p on p.shortcode = pts.short_code
),
all_posts_ts as (
    select
        stats_date,
        short_code short_code,
        max(handle) handle,
        max(post_title) post_title,
        max(post_link) post_link,
        max(post_thumbnail) post_thumbnail,
        max(hashtags) hashtags,
        max(keywords) keywords,
        max(published_at) published_at,
        max(profile_id) profile_id,
        max(followers) followers,
        max(post_type) post_type,
        max(comments) comments,
        max(likes) likes,
        max(total_engagement) total_engagement,
        max(plays) plays,
        max(impressions) impressions,
        max(reach) reach,
        maxIf(engagement_rate, engagement_rate < 100000) engagement_rate,
        max(saves) saves,
        max(swipe_ups) swipe_ups,
        max(mentions) mentions,
        max(sticker_taps) sticker_taps,
        max(shares) shares,
        max(exits) story_exits,
        max(taps_back) story_back_taps,
        max(taps_forward) story_forward_taps,
        max(updated_at) updated_at,
        'pipeline' insight_source
    from with_metrics
    group by stats_date, short_code
),
collection_posts_final as (
    select
        collection_posts.post_collection_id post_collection_id,
        collection_posts.id post_collection_item_id,
        collection_posts.short_code short_code,
        all_posts_ts.profile_id post_profile_id,
        all_posts_ts.handle post_handle,
        all_posts_ts.*
    from
        all_posts_ts
    left join collection_posts on all_posts_ts.short_code = collection_posts.short_code
),
channel_ids as (
    select post_profile_id from collection_posts_final group by post_profile_id
),
channel_names as (
    select channel_id, max(title) title, max(username) username
    from dbt.mart_youtube_account
    where channel_id in channel_ids
    group by channel_id
)
select
    'YOUTUBE' platform,
    ifNull(short_code, 'MISSING') post_short_code,
    ifNull(post_type, '') post_type,
    post_collection_item_id,
    post_title,
    post_link,
    post_thumbnail,
    hashtags,
    keywords,
    published_at,
    profile_id profile_social_id,
    handle profile_handle,
    post_handle,
    if(cn.username = '' or cn.username is null,
        ifNull(cn.title, 'MISSING'),
        cn.username) profile_name,
    'POST' collection_type,
    ifNull(toString(post_collection_id), 'MISSING') collection_id,
    followers,
    0 link_clicks,
    0 orders,
    0 delivered_orders,
    0 completed_orders,
    0 leaderboard_overall_orders,
    0 leaderboard_delivered_orders,
    0 leaderboard_completed_orders,
    ifNull(stats_date, date('2000-01-01')) stats_date,
    reach views,
    likes,
    comments,
    impressions,
    saves,
    plays,
    reach,
    swipe_ups,
    mentions,
    sticker_taps,
    shares,
    story_exits,
    story_back_taps,
    story_forward_taps,
    total_engagement,
    engagement_rate,
    0 real_reach,
    0 real_impressions,
    updated_at
from collection_posts_final
left join channel_names cn on cn.channel_id = collection_posts_final.profile_id
order by platform desc, profile_social_id desc, profile_handle desc
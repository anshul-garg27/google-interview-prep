{{ config(materialized = 'table', tags=["staging_collections"], order_by='collection_id,post_short_code,stats_date') }}
with
all_collection_profiles as (
    select id, platform, profile_social_id profile_id, profile_collection_id
    from {{ref('stg_coffee_staging_profile_collection_item')}} pci
    where platform = 'INSTAGRAM' and shortlist_id is null
),
all_profiles as (
    select profile_id from all_collection_profiles group by profile_id
),
posts as (
    select
        profile_id,
        shortcode,
        handle,
        post_type,
        caption post_title,
        extractAll(ifNull(caption, ''), '#[a-zA-Z0-9]+') hashtags,
        extractAll(ifNull(caption, ''), '[a-zA-Z0-9]{4,}') keywords,
        publish_time published_at,
        concat('https://instagram.com/p/', shortcode) post_link,
        display_url post_thumbnail
    from dbt.stg_beat_instagram_post where profile_id in all_profiles
),
shortcodes as (
    select shortcode from posts
),
thumbnails as (
    select
        entity_id shortcode,
        concat('https://d24w28i6lzk071.cloudfront.net/', asset_url) post_thumbnail
    from
        dbt.stg_beat_asset_log
        where entity_type = 'POST' and entity_id IN shortcodes
),
followers as (
    select profile_id, followers from dbt.stg_beat_instagram_account ia
    where ia.profile_id in all_profiles
),
base_ts_fake as (
    select
        platform,
        post_short_code short_code,
        stats_date,
        views,
        likes,
        comments,
        platform,
        stats_date,
        now() updated_at,
        post_short_code short_code,
        comments,
        likes,
        plays,
        0 saves,
        0 mentions,
        0 shares,
        0 taps_back,
        0 sticker_taps,
        0 exits,
        0 taps_forward,
        0 swipe_ups,
        1.0*real_reach real_reach,
        1.0*real_impressions real_impressions,
        likes + comments total_engagement
    from {{ref('mart_staging_fake_events')}}
    where short_code in shortcodes and platform = 'INSTAGRAM'
),
base_ts_raw as (
    select
        platform,
        toStartOfDay(pts.event_timestamp) stats_date,
        pts.event_timestamp updated_at,
        shortcode short_code,
        JSONExtractInt(metrics, 'comments') _comments,
        JSONExtractInt(metrics, 'likes') _likes,
        JSONExtractInt(metrics, 'play_count') _plays,
        JSONExtractInt(metrics, 'saves') _saves,
        JSONExtractInt(metrics, 'mentions') _mentions,
        JSONExtractInt(metrics, 'shares') _shares,
        JSONExtractInt(metrics, 'taps_back') _taps_back,
        JSONExtractInt(metrics, 'sticker_taps') _sticker_taps,
        JSONExtractInt(metrics, 'exits') _exits,
        JSONExtractInt(metrics, 'taps_forward') _taps_forward,
        JSONExtractInt(metrics, 'swipe_ups') _swipe_ups,

        JSONExtractString(metrics, 'comments') str_comments,
        JSONExtractString(metrics, 'likes') str_likes,
        JSONExtractString(metrics, 'play_count') str_plays,
        JSONExtractString(metrics, 'saves') str_saves,
        JSONExtractString(metrics, 'mentions') str_mentions,
        JSONExtractString(metrics, 'shares') str_shares,
        JSONExtractString(metrics, 'taps_back') str_taps_back,
        JSONExtractString(metrics, 'sticker_taps') str_sticker_taps,
        JSONExtractString(metrics, 'exits') str_exits,
        JSONExtractString(metrics, 'taps_forward') str_taps_forward,
        JSONExtractString(metrics, 'swipe_ups') str_swipe_ups,

        if(_comments = 0, toInt64OrZero(str_comments), _comments) comments,
        if(_likes = 0, toInt64OrZero(str_likes), _likes) likes,
        if(_plays = 0, toInt64OrZero(str_plays), _plays) plays,
        if(_saves = 0, toInt64OrZero(str_saves), _saves) saves,
        if(_mentions = 0, toInt64OrZero(str_mentions), _mentions) mentions,
        if(_shares = 0, toInt64OrZero(str_shares), _shares) shares,
        if(_taps_back = 0, toInt64OrZero(str_taps_back), _taps_back) taps_back,
        if(_sticker_taps = 0, toInt64OrZero(str_sticker_taps), _sticker_taps) sticker_taps,
        if(_exits = 0, toInt64OrZero(str_exits), _exits) exits,
        if(_taps_forward = 0, toInt64OrZero(str_taps_forward), _taps_forward) taps_forward,
        if(_swipe_ups = 0, toInt64OrZero(str_swipe_ups), _swipe_ups) swipe_ups,

        1.0*JSONExtractInt(metrics, 'reach') real_reach,
        1.0*JSONExtractInt(metrics, 'impressions') real_impressions,

        likes + comments total_engagement
    from _e.post_log_events pts
    where shortcode in shortcodes and platform = 'INSTAGRAM'
),
base_ts_all as (
    select
        platform,
        short_code,
        stats_date,
        updated_at,
        comments,
        likes,
        plays,
        saves,
        mentions,
        shares,
        taps_back,
        sticker_taps,
        exits,
        taps_forward,
        swipe_ups,
        total_engagement,
        real_reach,
        real_impressions
    from base_ts_raw
        union all
    select
        platform,
        short_code,
        stats_date,
        updated_at,
        comments,
        likes,
        plays,
        saves,
        mentions,
        shares,
        taps_back,
        sticker_taps,
        exits,
        taps_forward,
        swipe_ups,
        total_engagement,
        real_reach,
        real_impressions
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
        max(plays) plays,
        max(saves) saves,
        max(mentions) mentions,
        max(shares) shares,
        max(taps_back) taps_back,
        max(sticker_taps) sticker_taps,
        max(exits) exits,
        max(taps_forward) taps_forward,
        max(swipe_ups) swipe_ups,
        max(total_engagement) total_engagement,
        max(real_reach) real_reach,
        max(real_impressions) real_impressions
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
        p.hashtags hashtags,
        p.keywords keywords,
        p.post_title post_title,
        p.post_link post_link,
        coalesce(t.post_thumbnail, p.post_thumbnail) post_thumbnail,
        p.published_at published_at,
        pts.comments comments,
        pts.mentions mentions,
        pts.updated_at updated_at,
        f.followers followers,
        pts.total_engagement total_engagement,
        pts.saves saves,
        pts.shares shares,
        pts.exits exits,
        pts.sticker_taps sticker_taps,
        pts.taps_forward taps_forward,
        pts.taps_back taps_back,
        pts.swipe_ups swipe_ups,
        pts.likes likes,
        pts.plays plays,
        1.0*pts.real_reach real_reach,
        1.0*pts.real_impressions real_impressions,
        plays * (0.94 - (log2(followers) * 0.001)) _reach_reels,
        toFloat64(plays) _impressions_reels,
        (7.6 - (log10(likes) * 0.7)) * 0.85*likes _reach_non_reels,
        (7.6-(log10(likes) * 0.75))*0.95*likes _impressions_non_reels,
        if (post_type = 'reels', max2(_reach_reels,0), max2(_reach_non_reels,0)) _fake_reach,
        if (post_type = 'reels', max2(_impressions_reels,0), max2(_impressions_non_reels,0)) _fake_impressions,
        if(real_reach > 0, real_reach, _fake_reach) reach,
        if(real_impressions > 0, real_impressions, _fake_impressions) impressions,
        total_engagement / reach engagement_rate
    from base_ts pts
    left join posts p on p.shortcode = pts.short_code
    left join thumbnails t on t.shortcode = pts.short_code
    left join followers f on f.profile_id = profile_id
),
all_posts_ts as (
    select
        stats_date,
        short_code short_code,
        max(handle) handle,
        max(post_title) post_title,
        max(post_link) post_link,
        max(post_thumbnail) post_thumbnail,
        max(published_at) published_at,
        max(profile_id) profile_id,
        max(hashtags) hashtags,
        max(keywords) keywords,
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
        max(real_reach) real_reach,
        max(real_impressions) real_impressions,
        max(updated_at) updated_at,
        'pipeline' insight_source
    from with_metrics
    group by stats_date, short_code
),
collection_posts as (
    select
        all_collection_profiles.profile_collection_id profile_collection_id,
        all_collection_profiles.profile_id profile_id,
        all_posts_ts.profile_id post_profile_id,
        all_posts_ts.handle post_handle,
        all_posts_ts.*
    from
        all_posts_ts
    left join all_collection_profiles on all_posts_ts.profile_id = all_collection_profiles.profile_id
),
profile_ids as (
  select post_profile_id from collection_posts group by post_profile_id
),
profile_handles as (
    select profile_id, max(handle) handle, max(full_name) name
    from dbt.stg_beat_instagram_account where profile_id in profile_ids
    group by profile_id
)
select
    'INSTAGRAM' platform,
    ifNull(collection_posts.short_code, 'MISSING') post_short_code,
    ifNull(post_type, '') post_type,
    0 post_collection_item_id,
    post_title,
    post_link,
    post_thumbnail,
    hashtags,
    keywords,
    published_at,
    collection_posts.profile_id profile_social_id,
    ifNull(ph.handle, 'MISSING') profile_handle,
    post_handle,
    ifNull(ph.name, 'MISSING') profile_name,
    'PROFILE' collection_type,
    ifNull(toString(profile_collection_id), 'MISSING') collection_id,
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
    real_reach,
    real_impressions,
    updated_at
from collection_posts
left join profile_handles ph on ph.profile_id = collection_posts.profile_id
order by platform desc, profile_social_id desc, profile_handle desc

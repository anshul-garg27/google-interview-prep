{{ config(materialized = 'table', tags=["collections"], order_by='collection_id, post_short_code') }}
with
_from_collection_raw as (
    select
        platform,
        post_collection_id collection_id,
        mia.profile_id profile_social_id,
        item.short_code post_short_code,
        item.post_type post_type,
        '' post_title,
        post_link,
        post_thumbnail,
        [] hashtags,
        [] keywords,
        item.id post_collection_item_id,
        item.created_at published_at,
        'POST' collection_type,
        posted_by_handle profile_handle,
        posted_by_handle profile_name,
        mia.followers followers,
        0 views,
        0 likes,
        0 comments,
        0 saves,
        0 plays,
        0 reach,
        0 swipe_ups,
        0 mentions,
        0 sticker_taps,
        0 shares,
        0 story_exits,
        0 story_back_taps,
        0 story_forward_taps,
        0 link_clicks,
        0 orders,
        0 impressions,
        0 total_engagement,
        0 engagement_rate,
        item.created_at created_at,
        item.enabled enabled,
        replace(JSONExtractArrayRaw(ifNull(item.sponsor_links, '[]'))[1], '"', '') link,
        if (collection_id in (select arrayJoin(['731c9bf8-5428-4b6f-ac65-1d1b5c0d97d2']) as id), false, true) flag
    from {{ref('stg_coffee_post_collection_item')}} item
    left join dbt.stg_beat_instagram_account mia on mia.handle = item.posted_by_handle and item.posted_by_handle != ''
    where item.platform = 'INSTAGRAM' and item.enabled = True
),
items_shortcode as (
    select post_short_code
    from _from_collection_raw
    where post_type='story'
),
story_data as (
    select shortcode,
           argMax(profile_follower, updated_at) profile_follower,
           argMax(profile_er, updated_at) profile_er
    from dbt.stg_beat_instagram_post
    where shortcode in items_shortcode
    group by shortcode
),
collection_date as (
    select
            id,
            created_at
    from {{ref('stg_coffee_post_collection')}}
),
freeze_post_shortcode as (
    select post_short_code
    from _from_collection_raw
    where flag = false
),
freeze_post_values as (
    select shortcode, argMax(reach, updated_at) reach, argMax(impressions, updated_at) impressions
                from stg_beat_instagram_post
                where  shortcode in freeze_post_shortcode
                group by shortcode
),
_from_collection as (
    select
        item.platform platform,
        item.collection_id collection_id,
        item.profile_social_id profile_social_id,
        item.post_short_code post_short_code,
        item.post_type post_type,
        item.post_title post_title,
        item.post_link post_link,
        item.post_thumbnail post_thumbnail,
        item.hashtags hashtags,
        item.keywords keywords,
        item.post_collection_item_id post_collection_item_id,
        item.published_at published_at,
        item.collection_type collection_type,
        item.profile_handle profile_handle,
        item.profile_name profile_name,
        item.followers followers,
        item.views views,
        item.likes likes,
        item.comments comments,
        item.saves saves,
        item.plays plays,
        CASE
            WHEN item.created_at < '2023-11-22' THEN round((-0.000025017) * sd.profile_follower + (1.11 * (sd.profile_follower * abs(log2(sd.profile_er)) * 2 / 100)))
            ELSE round((-0.000025017) * sd.profile_follower + (1.11 * (sd.profile_follower * abs(log2((sd.profile_er + 2))) * 2 / 100)))
        END AS _story_reach,
        if(item.post_type = 'story', if(item.flag = false, toFloat64(fpv.reach),_story_reach), 0) reach,
        item.swipe_ups swipe_ups,
        item.mentions mentions,
        item.sticker_taps sticker_taps,
        item.shares shares,
        item.story_exits story_exits,
        item.story_back_taps story_back_taps,
        item.story_forward_taps story_forward_taps,
        item.link_clicks link_clicks,
        item.orders orders,
        IFNULL(
            CASE
                WHEN item.post_type = 'story' AND (cd.created_at >= '2023-11-23' or cd.id IN ('47bf0165-60fa-40fd-a5e0-705c5b3b9b07','293f6f80-6e8b-423f-bf4e-90e77b57cfd7','0a2d79fe-579b-43ec-8cba-166a9945fef9')) THEN if(item.flag=false, toFloat64(fpv.impressions), round(reach*(1.02 + ((reach % 10)/100))))
                ELSE item.impressions
            END, 0
        ) AS impressions,
        item.total_engagement total_engagement,
        item.total_engagement / reach engagement_rate,
        item.link link
    from _from_collection_raw item
    left join story_data sd on sd.shortcode = item.post_short_code
    left join collection_date cd on cd.id = item.collection_id
    left join freeze_post_values fpv on fpv.shortcode = item.post_short_code
    where item.platform = 'INSTAGRAM' and item.enabled = True
),
urls as (
    select link from _from_collection
    group by link
),
post_short_code as (
    select post_short_code from _from_collection
    group by post_short_code
),
link_clicks as (
    select concat('https://', url) link, sum(unique_clicks) clicks
    from dbt.mart_link_clicks
    where concat('https://', url) in urls
    group by link
),
manual_link_clicks as (
    select short_code, post_collection_id, fc.link, max(Clicks) clicks
    from dbt.mart_collection_post_clicks
    left join _from_collection fc on fc.post_short_code = short_code
    where short_code in post_short_code
    group by short_code, post_collection_id, fc.link
),
link_orders as (
    select concat('https://', link) url,
           sum(overall_orders) all_orders,
           sum(delivered_orders) delivered_orders,
           sum(completed_orders) completed_orders,
            sum(leaderboard_overall_orders) leaderboard_overall_orders,
           sum(leaderboard_delivered_orders) leaderboard_delivered_orders,
           sum(leaderboard_completed_orders) leaderboard_completed_orders
    from dbt.mart_link_clicks_orders
    where concat('https://', link)  in urls
    group by url
),
from_collection as (
    select
        platform,
        collection_id,
        profile_social_id,
        post_short_code,
        post_type,
        post_title,
        post_link,
        post_thumbnail,
        hashtags,
        keywords,
        post_collection_item_id,
        published_at,
        collection_type,
        profile_handle,
        profile_name,
        followers,
        views,
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
        if(lc.clicks>ifNull(mlc.clicks, 0), lc.clicks, ifNull(mlc.clicks,0)) link_clicks,
        lo.all_orders orders,
        lo.delivered_orders delivered_orders,
        lo.completed_orders completed_orders,
        lo.leaderboard_overall_orders leaderboard_overall_orders,
        lo.leaderboard_delivered_orders leaderboard_delivered_orders,
        lo.leaderboard_completed_orders leaderboard_completed_orders,
        total_engagement,
        engagement_rate
    from _from_collection base
    left join link_orders lo on lo.url = base.link
    left join link_clicks lc on lc.link = base.link
    left join manual_link_clicks mlc on mlc.short_code = base.post_short_code
),
from_ts as (
    select
        platform,
        collection_id,
        profile_social_id,
        post_short_code,
        max(post_type) post_type,
        max(post_title) post_title,
        max(post_link) post_link,
        max(post_thumbnail) post_thumbnail,
        max(hashtags) hashtags,
        max(keywords) keywords,
        max(post_collection_item_id) post_collection_item_id,
        max(published_at) published_at,
        max(collection_type) collection_type,
        max(profile_handle) profile_handle,
        max(profile_name) profile_name,
        max(followers) followers,
        argMax(views, stats_date) views,
        argMax(likes, stats_date) likes,
        argMax(comments, stats_date) comments,
        argMax(impressions, stats_date) impressions,
        max(saves) saves,
        argMax(plays, stats_date) plays,
        argMax(reach, stats_date) reach,
        max(swipe_ups) swipe_ups,
        max(mentions) mentions,
        max(sticker_taps) sticker_taps,
        max(shares) shares,
        max(story_exits) story_exits,
        max(story_back_taps) story_back_taps,
        max(story_forward_taps) story_forward_taps,
        max(link_clicks) link_clicks,
        max(orders) orders,
        max(delivered_orders) delivered_orders,
        max(completed_orders) completed_orders,
        max(leaderboard_overall_orders) leaderboard_overall_orders,
        max(leaderboard_delivered_orders) leaderboard_delivered_orders,
        max(leaderboard_completed_orders) leaderboard_completed_orders,
        max(total_engagement) total_engagement,
        total_engagement / reach engagement_rate
    from {{ref('mart_collection_post_ts')}}
    group by platform, collection_id, profile_social_id, post_short_code
),
all as (
    select
        platform,
        collection_id,
        profile_social_id,
        post_short_code,
        post_type,
        post_title,
        post_link,
        post_thumbnail,
        hashtags,
        keywords,
        post_collection_item_id,
        published_at,
        collection_type,
        profile_handle,
        profile_name,
        followers,
        views,
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
        link_clicks,
        orders,
        delivered_orders,
        completed_orders,
        leaderboard_overall_orders,
        leaderboard_delivered_orders,
        leaderboard_completed_orders,
        total_engagement,
        if(isInfinite(if(isNaN(engagement_rate), engagement_rate, 0)), engagement_rate, 0) engagement_rate
    from from_collection
    union all
    select
        platform,
        collection_id,
        profile_social_id,
        post_short_code,
        post_type,
        post_title,
        post_link,
        post_thumbnail,
        hashtags,
        keywords,
        post_collection_item_id,
        published_at,
        collection_type,
        profile_handle,
        profile_name,
        followers,
        views,
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
        link_clicks,
        orders,
        delivered_orders,
        completed_orders,
        leaderboard_overall_orders,
        leaderboard_delivered_orders,
        leaderboard_completed_orders,
        total_engagement,
        engagement_rate
    from from_ts
),
final as (
    select
        platform,
        collection_id,
        post_short_code,
        max(profile_social_id) profile_social_id,
        max(post_type) post_type,
        max(post_title) post_title,
        max(post_link) post_link,
        max(post_thumbnail) post_thumbnail,
        max(hashtags) hashtags,
        max(keywords) keywords,
        max(post_collection_item_id) post_collection_item_id,
        max(published_at) published_at,
        max(collection_type) collection_type,
        max(profile_handle) profile_handle,
        max(profile_name) profile_name,
        max(followers) followers,
        max(views) views,
        max(likes) likes,
        max(comments) comments,
        max(impressions) impressions,
        max(saves) saves,
        max(plays) plays,
        max(reach) reach,
        max(swipe_ups) swipe_ups,
        max(mentions) mentions,
        max(sticker_taps) sticker_taps,
        max(shares) shares,
        max(story_exits) story_exits,
        max(story_back_taps) story_back_taps,
        max(story_forward_taps) story_forward_taps,
        max(link_clicks) link_clicks,
        max(orders) orders,
        max(delivered_orders) delivered_orders,
        max(completed_orders) completed_orders,
        max(leaderboard_overall_orders) leaderboard_overall_orders,
        max(leaderboard_delivered_orders) leaderboard_delivered_orders,
        max(leaderboard_completed_orders) leaderboard_completed_orders,
        max(total_engagement) total_engagement,
        max(engagement_rate) engagement_rate
    from all
    group by platform, collection_id, post_short_code
)
select * from final
SETTINGS max_bytes_before_external_group_by = 20000000000, max_bytes_before_external_sort = 20000000000, join_algorithm='partial_merge'

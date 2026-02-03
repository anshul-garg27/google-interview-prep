{{ config(materialized = 'table', tags=["daily", "leaderboard"], order_by="date, ifNull(platform_id, 'MISSING')") }}
with
channels as (
    select channelid from dbt.mart_youtube_account where country = 'IN'
),
profile_ids as (
    select ig_id from dbt.mart_instagram_account where country = 'IN'
),
yt_monthly as (
    select
        date,
        p_date,
        platform,
        platform_id,
        toInt64(if(base.views < 0, 0, base.views)) views_total,
        toInt64(if(views_change < 0, 0, views_change)) views,
        likes,
        comments,
        engagement,
        plays,
        0 plays_total,
        followers,
        following,
        avg_views,
        avg_likes,
        avg_comments,
        engagement_rate,
        avg_plays,
        toInt64(if(uploads_change < 0, 0, uploads_change)) uploads,
        followers_change,
        monthly_stats
    from dbt.mart_yt_account_monthly base
    where platform_id in channels
),
yt_weekly as (
    select
        date,
        p_date,
        platform,
        platform_id,
        toInt64(if(base.views < 0, 0, base.views)) views_total,
        toInt64(if(views_change < 0, 0, views_change)) views,
        likes,
        comments,
        engagement,
        plays,
        0 plays_total,
        followers,
        following,
        avg_views,
        avg_likes,
        avg_comments,
        engagement_rate,
        avg_plays,
        toInt64(if(uploads_change < 0, 0, uploads_change)) uploads,
        followers_change,
        false monthly_stats
    from dbt.mart_yt_account_weekly base
    where platform_id in channels
),
insta_weekly as (
    select
        date,
        p_date,
        platform,
        platform_id,
        0 views_total,
        0 views,
        likes,
        comments,
        engagement,
        plays,
        plays_total,
        followers,
        following,
        avg_views,
        avg_likes,
        avg_comments,
        engagement_rate,
        avg_plays,
        toInt64(uploads) uploads,
        followers_change,
        monthly_stats
    from dbt.mart_insta_account_weekly
    where platform_id in profile_ids
),
insta_monthly as (
    select
        date,
        p_date,
        platform,
        platform_id,
        0 views_total,
        0 views,
        likes,
        comments,
        engagement,
        plays,
        plays_total,
        followers,
        following,
        avg_views,
        avg_likes,
        avg_comments,
        engagement_rate,
        avg_plays,
        toInt64(uploads) uploads,
        followers_change,
        monthly_stats
    from dbt.mart_insta_account_monthly
    where platform_id in profile_ids
)
select * from yt_weekly
    union all
select * from insta_weekly
    union all
select * from insta_monthly
    union all
select * from yt_monthly
SETTINGS max_bytes_before_external_group_by = 20000000000, max_bytes_before_external_sort = 20000000000, join_algorithm='partial_merge'

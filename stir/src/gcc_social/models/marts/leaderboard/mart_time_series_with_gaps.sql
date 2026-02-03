{{ config(materialized = 'table', tags=["disabled"], order_by='date, platform_id') }}
with
channels as (
    select channelid from dbt.mart_youtube_account where country = 'IN'
),
profile_ids as (
    select ig_id from dbt.mart_instagram_account where country = 'IN'
),
yt_weekly as (
    select
        date,
        platform,
        platform_id,
        views,
        likes,
        comments,
        engagement,
        plays,
        followers,
        following,
        avg_views,
        avg_likes,
        avg_comments,
        engagement_rate,
        avg_plays,
        uploads,
        followers_change
    from dbt.mart_yt_account_weekly
    where platform_id in channels
),
insta_weekly as (
    select
        date,
        platform,
        platform_id,
        views,
        likes,
        comments,
        engagement,
        plays,
        followers,
        following,
        avg_views,
        avg_likes,
        avg_comments,
        engagement_rate,
        avg_plays,
        uploads,
        followers_change
    from dbt.mart_insta_account_weekly
    where platform_id in profile_ids
)
select * from yt_weekly
    union all
select * from insta_weekly
SETTINGS max_bytes_before_external_group_by = 20000000000, max_bytes_before_external_sort = 20000000000, join_algorithm='partial_merge'

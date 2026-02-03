{{ config(materialized = 'table', tags=["daily", "leaderboard"], order_by='channelid') }}
with
lasty as (
  select
            id channelid,
            max(toInt64OrZero(lifetime_subscribers)) followers,
            max(toInt64OrZero(lifetime_uploads)) uploads,
            max(toInt64OrZero(lifetime_views)) views
        from vidooly.channels_monthly_stats_20230219
        where parseDateTimeBestEffortOrZero(month) <= now() - INTERVAL 1 YEAR
    group by channelid
),
lastw as (
    select
        platform_id channelid,
        argMax(followers, date) followers,
        argMax(uploads, date) uploads,
        argMax(views, date) views
    from {{ref('mart_yt_account_weekly')}} lastw
    where lastw.date <= now() - INTERVAL 1 WEEK
    group by channelid
),
lastm as (
    select
        platform_id channelid,
        argMax(followers, date) followers,
        argMax(uploads, date) uploads,
        argMax(views, date) views
    from {{ref('mart_yt_account_weekly')}} lastm
    where lastm.date <= now() - INTERVAL 1 MONTH
    group by channelid
),
currentw as (
    select
        platform_id channelid,
        argMax(followers, date) followers,
        argMax(uploads, date) uploads,
        argMax(views, date) views
    from {{ref('mart_yt_account_weekly')}} lastw
    group by channelid
),
weekly_growth as (
    select
        channelid,
        currentw.followers - lastw.followers follower_7d,
        currentw.uploads - lastw.uploads uploads_7d,
        currentw.views - lastw.views views_7d
    from currentw
    left join lastw on currentw.channelid = lastw.channelid
),
monthly_growth as (
    select
        channelid,
        currentw.followers - lastm.followers follower_30d,
        currentw.uploads - lastm.uploads uploads_30d,
        currentw.views - lastm.views views_30d
    from currentw
    left join lastm on currentw.channelid = lastm.channelid
),
yearly_growth as (
    select
        channelid,
        currentw.followers - lasty.followers follower_1y,
        currentw.uploads - lasty.uploads uploads_1y,
        currentw.views - lasty.views views_1y
    from currentw
    left join lasty on currentw.channelid = lasty.channelid
)
select
    ifNull(wg.channelid, 'MISSING') channelid,
    follower_7d,
    uploads_7d,
    views_7d,
    follower_30d,
    uploads_30d,
    views_30d,
    follower_1y,
    uploads_1y,
    views_1y
from weekly_growth wg
left join monthly_growth mg on mg.channelid = wg.channelid
left join yearly_growth yg on yg.channelid = wg.channelid
SETTINGS max_bytes_before_external_group_by = 20000000000, max_bytes_before_external_sort = 20000000000, join_algorithm='partial_merge'

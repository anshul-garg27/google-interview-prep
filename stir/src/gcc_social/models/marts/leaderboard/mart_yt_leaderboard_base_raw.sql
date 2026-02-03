{{ config(materialized = 'table', tags=["daily", "leaderboard"], order_by="ifNull(month, 'MISSING'), ifNull(channel_id, 'MISSING')") }}
with
old_stats as (
    select
        toStartOfMonth(parseDateTimeBestEffortOrZero(month)) month,
        id channel_id,
        toInt64OrZero(lifetime_views) lifetime_views,
        toInt64OrZero(lifetime_subscribers) lifetime_followers,
        toInt64OrZero(lifetime_uploads) lifetime_uploads,
        1 priority
    from vidooly.channels_monthly_stats_20230219
),
new_stats_vidooly as (
    select
        toStartOfMonth(parseDateTimeBestEffortOrZero(month)) month,
        id channel_id,
        lifetime_views lifetime_views,
        lifetime_subscribers lifetime_followers,
        lifetime_uploads lifetime_uploads,
        1 priority
    from vidooly.channels_monthly_stats_20230709
),
new_stats_gcc as (
    select 
        toStartOfMonth(event_timestamp) month,
        handle channel_id,
        argMax(JSONExtractInt(metrics, 'view_count'), event_timestamp) lifetime_views,
        argMax(JSONExtractInt(metrics, 'subscriber_count'), event_timestamp) lifetime_followers,
        argMax(JSONExtractInt(metrics, 'video_count'), event_timestamp) lifetime_uploads,
        0 priority
    from _e.profile_log_events 
    where 
        platform = 'YOUTUBE' and event_timestamp >= date('2023-06-01')
    group by month, channel_id
    order by month desc, channel_id desc
),
monthly_stats as (
    select * from (select * from old_stats union all select * from new_stats_vidooly union all select * from new_stats_gcc)
    order by channel_id desc, month desc
),
final as (
    select 
        month, 
        channel_id, 
        argMax(lifetime_views, priority) lifetime_views,
        argMax(lifetime_followers, priority) lifetime_followers,
        argMax(lifetime_uploads, priority) lifetime_uploads
    from monthly_stats
    group by month, channel_id
)
select * from final where month >= date('2022-01-01')
SETTINGS max_bytes_before_external_group_by = 20000000000, max_bytes_before_external_sort = 20000000000, join_algorithm='partial_merge'


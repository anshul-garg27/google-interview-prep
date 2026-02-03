{{ config(materialized = 'table', tags=["daily", "leaderboard"], order_by='month, channel_id') }}
with monthly_stats_with_prev as (
    select
        ifNull(month, 'MISSING') month,
        ifNull(channel_id, 'MISSING') channel_id,
        lifetime_followers followers,
        lifetime_views,
        lifetime_uploads,
        groupArray(month) OVER (PARTITION BY channel_id ORDER BY month ASC
            ROWS BETWEEN 1 PRECEDING AND 0 FOLLOWING) AS month_values,
        month_values[1] p_month,
        abs(p_month - month) <= 31 last_month_available,

        groupArray(followers) OVER (PARTITION BY channel_id ORDER BY month ASC
        Rows BETWEEN 1 PRECEDING AND 0 FOLLOWING) AS followers_value,
        followers_value[1] prev_followers,
        if(last_month_available = 1, followers - prev_followers, NULL) followers_change,
        groupArray(lifetime_views) OVER (PARTITION BY channel_id ORDER BY month ASC
        Rows BETWEEN 1 PRECEDING AND 0 FOLLOWING) AS views_value,
        views_value[1] prev_views,
        if(last_month_available = 1, lifetime_views - prev_views, NULL) views,
        groupArray(lifetime_uploads) OVER (PARTITION BY channel_id ORDER BY month ASC
        Rows BETWEEN 1 PRECEDING AND 0 FOLLOWING) AS uploads_value,
        uploads_value[1] prev_uploads,
        if(last_month_available = 1, lifetime_uploads - prev_uploads, NULL) uploads
    from {{ref('mart_yt_leaderboard_base_raw')}}
)
select * from monthly_stats_with_prev where month >= date('2022-01-01')
SETTINGS max_bytes_before_external_group_by = 20000000000, max_bytes_before_external_sort = 20000000000, join_algorithm='partial_merge'


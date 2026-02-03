{{ config(materialized = 'table', tags=["daily", "leaderboard"], order_by="date, ifNull(platform_id, 'MISSING')") }}
with
monthly_stats_raw as (
    select
        month date,
        'YOUTUBE' platform,
        channel_id platform_id,
        max(lifetime_views) views,
        max(lifetime_followers)  followers,
        max(lifetime_uploads) uploads,
        0 engagement_rate,
        0 likes,
        0 plays,
        0 comments,
        0 engagement,
        0 following,
        0 avg_views,
        0 avg_likes,
        0 avg_comments,
        0 avg_plays,
        true monthly_stats
    from dbt.mart_yt_leaderboard_base_raw
    group by date, platform, platform_id
),
monthly_stats_final_raw as (
    select * from monthly_stats_raw where date >= date('2020-01-01')
),
monthly_stats_final as (
    select
        date,
        platform,
        platform_id,
        views,
        followers,
        toUInt64(uploads) uploads,
        groupArray(date) OVER (PARTITION BY platform_id ORDER BY date ASC
            ROWS BETWEEN 1 PRECEDING AND 0 FOLLOWING) AS date_values,
        date_values[1] p_date,
        abs(p_date - date) <= 31 last_month_available,
        groupArray(followers) OVER (PARTITION BY platform_id ORDER BY date ASC
        Rows BETWEEN 1 PRECEDING AND 0 FOLLOWING) AS followers_values,
        groupArray(views) OVER (PARTITION BY platform_id ORDER BY date ASC
        Rows BETWEEN 1 PRECEDING AND 0 FOLLOWING) AS views_values,
        groupArray(uploads) OVER (PARTITION BY platform_id ORDER BY date ASC
        Rows BETWEEN 1 PRECEDING AND 0 FOLLOWING) AS uploads_values,
        followers_values[1] followers_prev,
        views_values[1] views_prev,
        uploads_values[1] uploads_prev,
        if(last_month_available = 1, followers - followers_prev, NULL) followers_change,
        if(last_month_available = 1, views - views_prev, NULL) views_change,
        if(last_month_available = 1, uploads - uploads_prev, NULL) uploads_change,
        engagement_rate,
        likes,
        plays,
        comments,
        engagement,
        following,
        avg_views,
        avg_likes,
        avg_comments,
        avg_plays,
        monthly_stats
    from monthly_stats_final_raw
),
final as (
    select
        date,
        p_date,
        platform,
        platform_id,
        views,
        followers,
        toUInt64(uploads) uploads,
        if(followers_change is NULL or followers_change<0 , 0, followers_change) followers_change,
        if(views_change is NULL or views_change<0 , 0, views_change) views_change,
        if(uploads_change is NULL or uploads_change<0 , 0, uploads_change) uploads_change,
        engagement_rate,
        likes,
        plays,
        comments,
        engagement,
        following,
        avg_views,
        avg_likes,
        avg_comments,
        avg_plays,
        monthly_stats
    from monthly_stats_final
)
select * from final
SETTINGS max_bytes_before_external_group_by = 20000000000, max_bytes_before_external_sort = 20000000000, join_algorithm='partial_merge'

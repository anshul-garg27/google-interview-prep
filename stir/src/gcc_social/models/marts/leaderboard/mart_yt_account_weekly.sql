{{ config(materialized = 'table', tags=["daily", "leaderboard"], order_by="date, ifNull(platform_id, 'MISSING')") }}
with
weekly_stats_old as (
	select
        toStartOfMonth(parseDateTimeBestEffortOrZero(month)) date,
        'YOUTUBE' platform,
        id platform_id,
        argMin(toInt64OrZero(lifetime_views), parseDateTimeBestEffortOrZero(month)) views,
        argMin(toInt64OrZero(lifetime_subscribers), parseDateTimeBestEffortOrZero(month))  followers,
        argMin(toUInt64OrZero(lifetime_uploads), parseDateTimeBestEffortOrZero(month)) uploads,
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
        false monthly_stats
    from vidooly.channels_monthly_stats_20230219
    where date < date('2023-02-01')
    group by date, platform_id
),
ts_vidooly as (
    select 
        toStartOfWeek(parseDateTimeBestEffortOrZero(splitByChar('.', splitByChar('_', _path)[-1])[1])) date,
        'YOUTUBE' platform,
        id platform_id,
        toInt64(argMin(toInt64OrZero(view_count), parseDateTimeBestEffortOrZero(splitByChar('.', splitByChar('_', _path)[-1])[1]))) views,
        toInt64(argMin(toInt64OrZero(subscriber_count), parseDateTimeBestEffortOrZero(splitByChar('.', splitByChar('_', _path)[-1])[1])))  followers,
        toInt64(argMin(toUInt64OrZero(video_count), parseDateTimeBestEffortOrZero(splitByChar('.', splitByChar('_', _path)[-1])[1]))) uploads
    from vidooly.youtube_channels
    group by date, platform, platform_id
    having date < '2023-08-01' and date >= date('2023-02-01')
),
ts_gcc as (
    select 
        toStartOfWeek(event_timestamp) date,
        'YOUTUBE' platform,
        handle platform_id,
        toInt64(argMin(JSONExtractInt(metrics, 'view_count'), event_timestamp)) views,
        toInt64(argMin(JSONExtractInt(metrics, 'subscriber_count'), event_timestamp))  followers,
        toInt64(argMin(JSONExtractInt(metrics, 'video_count'), event_timestamp)) uploads
    from _e.profile_log_events
    where platform = 'YOUTUBE'
    group by date, platform, platform_id
    having date >= '2023-08-01'
),
ts_all as (
    select * from ts_vidooly union all select * from ts_gcc
),
weekly_stats_new as (
    select
        date,
        platform,
        platform_id,
        max(views) views,
        max(followers) followers,
        max(uploads) uploads,
        groupArray(date) OVER (PARTITION BY platform_id ORDER BY date ASC
            ROWS BETWEEN 1 PRECEDING AND 0 FOLLOWING) AS date_values,
        date_values[1] p_date,
        abs(p_date - date) <= 7 last_month_available,
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
        false monthly_stats
    from ts_all
    group by date, platform, platform_id
),
weekly_stats_old_final as (
	select
        date,
        p_date,
        platform,
        platform_id,
        toInt64(views) views,
        toInt64(followers) followers,
        toInt64(if(followers_change<0 , 0, followers_change)) followers_change,
        toInt64(if(views_change<0 , 0, views_change)) views_change,
        toInt64(if(uploads_change<0 , 0, uploads_change)) uploads_change,
        toInt64(uploads) uploads,
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
    from weekly_stats_old
),
weekly_stats_new_final as (
    select
        date,
        p_date,
        platform,
        platform_id,
        views,
        followers,
        if(followers_change<0 , 0, followers_change) followers_change,
        if(views_change<0 , 0, views_change) views_change,
        if(uploads_change<0 , 0, uploads_change) uploads_change,
        uploads,
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
    from weekly_stats_new 
),
weekly_stats_final as (
	select * from weekly_stats_old_final
	union all
	select * from weekly_stats_new_final
)
select * from weekly_stats_final
SETTINGS max_bytes_before_external_group_by = 20000000000, max_bytes_before_external_sort = 20000000000, join_algorithm='partial_merge'

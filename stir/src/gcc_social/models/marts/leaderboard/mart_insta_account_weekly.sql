{{ config(materialized = 'table', tags=["daily", "leaderboard"], order_by='date, platform_id') }}
with
views_by_week as (
    select
            toStartOfWeek(publish_time) week,
            profile_id,
            uniqExact(post_shortcode) posts,
            max(handle) handle,
            sum(views) views,
            sum(likes) likes,
            sum(comments) comments,
            sum(plays) plays,
            sum(reach) reach,
            sum(impressions) impressions,
            sum(plays) over (partition by profile_id order by week asc) plays_total,
            likes + comments engagement,
            views / posts avg_views,
            likes / posts avg_likes,
            comments / posts avg_comments,
            plays / posts avg_plays,
            engagement / posts avg_engagement
    from dbt.mart_instagram_all_posts_with_ranks
    group by week, profile_id
    order by week desc, views desc
),
followers_by_week_old as (
    select
        toStartOfWeek(parseDateTimeBestEffortOrZero(concat(month, '01'))) week,
        page_id profile_id,
        toInt64OrZero(prev_followers) + toInt64OrZero(followers) followers,
        0 following
    from
    vidooly.instagram_page_stats stats
    order by week desc
),
dec_followers as (
   select
        toStartOfWeek(date('2022-12-01')) week,
        toString(profile_id) profile_id,
        toInt64OrZero(followers) followers,
        0 following
    from
        dbt.stg_vidooly_instagram_profiles profiles
    order by week desc
),
followers_by_week_new as (
    select toStartOfWeek(insert_timestamp)           week,
        profile_id,
        max(JSONExtractInt(metrics, 'followers')) followers,
        max(JSONExtractInt(metrics, 'following')) following
    from _e.profile_log_events
    where platform = 'INSTAGRAM'
    and JSONHas(metrics, 'followers')
    group by week, profile_id
),
followers_by_week_all as (
    select * from
        followers_by_week_old
    union all
        select * from
        followers_by_week_new
    union all
        select * from
        dec_followers
),
followers_by_week as (
        select
            week,
            profile_id,
            max(followers) followers,
            max(following) following,
            groupArray(week) OVER (PARTITION BY profile_id ORDER BY week ASC
            ROWS BETWEEN 1 PRECEDING AND 0 FOLLOWING) AS week_values,
            week_values[1] p_week,
            abs(p_week - week) <= 7 last_week_available,
            groupArray(followers) OVER (PARTITION BY profile_id ORDER BY week ASC
            Rows BETWEEN 1 PRECEDING AND 0 FOLLOWING) AS followers_values,
            followers_values[1] followers_prev,
            if(last_week_available = 1, followers - followers_prev, NULL) followers_change
        from followers_by_week_all
        group by week, profile_id
        order by profile_id asc, week desc
),
final as (
    select
        ifNull(coalesce(followers.week, views.week), toStartOfWeek(date('2000-01-01'))) date,
        followers.p_week p_date,
        'INSTAGRAM' platform,
        ifNull(coalesce(followers.profile_id, views.profile_id), '') platform_id,
        views.views views,
        views.likes likes,
        views.comments comments,
        views.engagement engagement,
        views.plays plays,
        views.reach reach,
        views.impressions impressions,
        views.plays_total plays_total,
        followers.followers followers,
        followers.following following,
        followers.followers_change followers_change,
        views.avg_views avg_views,
        views.avg_likes avg_likes,
        views.avg_comments avg_comments,
        views.avg_engagement / followers engagement_rate,
        views.avg_plays avg_plays,
        views.posts uploads,
        1 enabled,
        false monthly_stats
    from
        followers_by_week followers
        left join views_by_week views on followers.profile_id = views.profile_id and views.week = followers.week
    order by date desc
)
select * from final
SETTINGS max_bytes_before_external_group_by = 20000000000, max_bytes_before_external_sort = 20000000000, join_algorithm='partial_merge'

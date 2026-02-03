{{ config(materialized = 'table', tags=["daily", "leaderboard"], order_by='date, platform_id') }}
with
views_by_month as (
    select
            toStartOfMonth(publish_time) month,
            profile_id,
            uniqExact(post_shortcode) posts,
            max(handle) handle,
            sum(views) views,
            sum(likes) likes,
            sum(comments) comments,
            sum(plays) plays,
            sum(plays) over (partition by profile_id order by month asc) plays_total,
            likes + comments engagement,
            views / posts avg_views,
            likes / posts avg_likes,
            comments / posts avg_comments,
            plays / posts avg_plays,
            engagement / posts avg_engagement
    from dbt.mart_instagram_all_posts_with_ranks
    group by month, profile_id
    order by month desc, views desc
),
followers_by_month_old as (
    select
        toStartOfMonth(parseDateTimeBestEffortOrZero(concat(month, '01'))) month,
        page_id profile_id,
        toInt64OrZero(prev_followers) + toInt64OrZero(followers) followers,
        0 following
    from
    vidooly.instagram_page_stats stats
    order by month desc
),
dec_followers as (
   select
        toStartOfMonth(date('2022-12-01')) month,
        toString(profile_id) profile_id,
        toInt64OrZero(followers) followers,
        0 following
    from
        dbt.stg_vidooly_instagram_profiles profiles
    order by month desc
),
followers_by_month_new as (
    select 
       toStartOfMonth(insert_timestamp)                               month,
       profile_id,
       argMax(JSONExtractInt(metrics, 'followers'), insert_timestamp) followers,
       argMax(JSONExtractInt(metrics, 'following'), insert_timestamp) following
    from _e.profile_log_events
    where platform = 'INSTAGRAM' and JSONHas(metrics, 'followers')
    group by month, profile_id
),
followers_by_month_all as (
    select * from
        followers_by_month_old
    union all
        select * from
        followers_by_month_new
    union all
        select * from
        dec_followers
),
followers_by_month as (
        select
            month,
            profile_id,
            argMax(followers, month) followers,
            argMax(following, month) following,
            groupArray(month) OVER (PARTITION BY profile_id ORDER BY month ASC
            ROWS BETWEEN 1 PRECEDING AND 0 FOLLOWING) AS month_values,
            month_values[1] p_month,
            abs(p_month - month) <= 31 last_month_available,
            groupArray(followers) OVER (PARTITION BY profile_id ORDER BY month ASC
            Rows BETWEEN 1 PRECEDING AND 0 FOLLOWING) AS followers_values,
            followers_values[1] followers_prev,
            if(last_month_available = 1, followers - followers_prev, NULL) followers_change
        from followers_by_month_all
        group by month, profile_id
        order by profile_id asc, month desc
),
final as (
    select
        ifNull(coalesce(followers.month, views.month), toStartOfMonth(date('2000-01-01'))) date,
        followers.p_month p_date,
        'INSTAGRAM' platform,
        ifNull(coalesce(followers.profile_id, views.profile_id), '') platform_id,
        views.views views,
        views.likes likes,
        views.comments comments,
        views.engagement engagement,
        views.plays plays,
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
        true monthly_stats
    from
        followers_by_month followers
        left join views_by_month views on followers.profile_id = views.profile_id and views.month = followers.month
    order by date desc
)
select * from final
SETTINGS max_bytes_before_external_group_by = 20000000000, max_bytes_before_external_sort = 20000000000, join_algorithm='partial_merge'

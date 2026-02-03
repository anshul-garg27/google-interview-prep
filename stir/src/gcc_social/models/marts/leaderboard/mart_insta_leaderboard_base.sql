{{ config(materialized = 'table', tags=["daily", "leaderboard"], order_by='month, profile_id') }}
with
views_by_month as (
    select
            toStartOfMonth(publish_time) month,
            profile_id,
            toInt64(uniqExact(post_shortcode)) posts,
            max(handle) handle,
            sum(views) views,
            sum(likes) likes,
            sum(comments) comments,
            sum(plays) plays,
            sum(plays) OVER (partition by profile_id order by month asc) total_plays,
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
        toInt64OrZero(prev_followers) + toInt64OrZero(followers) followers
    from
    vidooly.instagram_page_stats stats
    order by month desc
),
dec_followers as (
   select
        toStartOfMonth(date('2022-12-01')) month,
        toString(profile_id) profile_id,
        toInt64OrZero(followers) followers
    from
        dbt.stg_vidooly_instagram_profiles profiles
    order by month desc
),
followers_by_month_new as (
    select toStartOfMonth(toStartOfWeek(insert_timestamp))           month,
        profile_id,
        max(JSONExtractInt(metrics, 'followers')) followers
    from _e.profile_log_events
    where platform = 'INSTAGRAM'
    and JSONHas(metrics, 'followers')
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
            max(followers) followers
        from followers_by_month_all
        group by month, profile_id
        order by profile_id asc, month desc
),
followers_by_month_with_prev as (
    select
        month,
        profile_id,
        followers,
        groupArray(month) OVER (PARTITION BY profile_id ORDER BY month ASC
            ROWS BETWEEN 1 PRECEDING AND 0 FOLLOWING) AS month_values,
        month_values[1] p_month,
        abs(p_month - month) <= 31 last_month_available,

        groupArray(followers) OVER (PARTITION BY profile_id ORDER BY month ASC
        Rows BETWEEN 1 PRECEDING AND 0 FOLLOWING) AS followers_value,
        followers_value[1] prev_followers,
        if(last_month_available = 1 , followers - prev_followers, NULL) followers_change
    from followers_by_month
),
final as (
    select
        ifNull(coalesce(if(followers.month <= date('2000-01-01'), null, followers.month), views.month), toStartOfMonth(date('2000-01-01'))) month,
        ifNull(coalesce(views.profile_id, followers.profile_id), '') profile_id,
        views.views views,
        views.likes likes,
        views.comments comments,
        views.engagement engagement,
        views.plays plays,
        views.total_plays total_plays,
        total_plays - plays prev_plays,            
        followers.followers followers,
        followers.prev_followers prev_followers,
        followers.followers_change followers_change,
        views.avg_views avg_views,
        views.posts uploads,
        views.avg_likes avg_likes,
        views.avg_comments avg_comments,
        views.avg_engagement / followers avg_er,
        views.avg_plays avg_plays
    from
        followers_by_month_with_prev followers
        left join views_by_month views on followers.profile_id = views.profile_id and views.month = followers.month
    order by month desc
)
select * from final
SETTINGS max_bytes_before_external_group_by = 20000000000, max_bytes_before_external_sort = 20000000000, join_algorithm='partial_merge'

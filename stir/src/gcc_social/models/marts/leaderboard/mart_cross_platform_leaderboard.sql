{{ config(materialized = 'table', tags=["daily", "leaderboard"], order_by='month, platform_id') }}
with
insta_with_links as (
    select
        il.month month,
        il.category category,
        il.language language,
        il.country country,
        il.profile_type profile_type,
        il.plays views,
        il.prev_plays prev_views,
        il.uploads uploads,
        il.plays ia_views,
        0 yt_views,
        il.followers followers,
        il.followers_change followers_change,
        il.prev_followers prev_followers,
        il.likes likes,
        il.comments comments,
        'IA' profile_platform,
        il.platform_id platform_id,
        1 priority
    from
        {{ref('mart_instagram_leaderboard')}} il
),
yt_with_links as (
    select
        yl.month month,
        yl.category category,
        yl.language language,
        yl.country country,
        yl.profile_type profile_type,
        yl.views views,
        yl.prev_views prev_views,
        yl.uploads uploads,
        0 ia_views,
        yl.views yt_views,
        yl.followers followers,
        yl.followers_change followers_change,
        yl.prev_followers prev_followers,
        yl.likes likes,
        yl.comments comments,
        if(link.profile_id != '', 'IA', 'YA') profile_platform,
        if(link.profile_id != '', link.profile_id, yl.platform_id) platform_id,
        2 priority
    from
        {{ref('mart_youtube_leaderboard')}} yl
    left join dbt.mart_linked_socials link on link.channel_id = yl.platform_id
),
combined as (
    select * from insta_with_links
             union all
    select * from yt_with_links
),
base as (
    select
        month,
        profile_platform,
        platform_id,
        argMax(category, priority) category,
        max(profile_type) profile_type,
        argMax(language, priority) language,
        sum(views) views,
        sum(prev_views) prev_views,
        sum(uploads) uploads,
        sum(ia_views) ia_views,
        sum(yt_views) yt_views,
        sum(followers) followers,
        sum(prev_followers) prev_followers,
        sum(followers_change) followers_change,
        sum(likes) likes,
        sum(comments) comments,

        rank() over (partition by month order by followers desc, platform_id desc) followers_rank,
        rank() over (partition by month order by followers_change desc, platform_id desc) followers_change_rank,
        rank() over (partition by month order by views desc, platform_id desc) views_rank,

        rank() over (partition by month, category order by followers desc, platform_id desc) followers_rank_by_cat,
        rank() over (partition by month, category order by followers_change desc, platform_id desc) followers_change_rank_by_cat,
        rank() over (partition by month, category order by views desc, platform_id desc) views_rank_by_cat,


        rank() over (partition by month, language order by followers desc, platform_id desc) followers_rank_by_lang,
        rank() over (partition by month, language order by followers_change desc, platform_id desc) followers_change_rank_by_lang,
        rank() over (partition by month, language order by views desc, platform_id desc) views_rank_by_lang,


        rank() over (partition by month, category, language order by followers desc, platform_id desc) followers_rank_by_cat_lang,
        rank() over (partition by month, category, language order by followers_change desc, platform_id desc) followers_change_rank_by_cat_lang,
        rank() over (partition by month, category, language order by views desc, platform_id desc) views_rank_by_cat_lang,

        rank() over (partition by month, profile_type order by followers desc, platform_id desc) followers_rank_by_profile,
        rank() over (partition by month, profile_type order by followers_change desc, platform_id desc) followers_change_rank_by_profile,
        rank() over (partition by month, profile_type order by views desc, platform_id desc) views_rank_by_profile,


        rank() over (partition by month, category, profile_type order by followers desc, platform_id desc) followers_rank_by_cat_profile,
        rank() over (partition by month, category, profile_type order by followers_change desc, platform_id desc) followers_change_rank_by_cat_profile,
        rank() over (partition by month, category, profile_type order by views desc, platform_id desc) views_rank_by_cat_profile,


        rank() over (partition by month, language, profile_type order by followers desc, platform_id desc) followers_rank_by_lang_profile,
        rank() over (partition by month, language, profile_type order by followers_change desc, platform_id desc) followers_change_rank_by_lang_profile,
        rank() over (partition by month, language, profile_type order by views desc, platform_id desc) views_rank_by_lang_profile,


        rank() over (partition by month, category, language, profile_type order by followers desc, platform_id desc) followers_rank_by_cat_lang_profile,
        rank() over (partition by month, category, language, profile_type order by followers_change desc, platform_id desc) followers_change_rank_by_cat_lang_profile,
        rank() over (partition by month, category, language, profile_type order by views desc, platform_id desc) views_rank_by_cat_lang_profile
    from
        combined
    group by month, profile_platform, platform_id
),
sorted as (
    select
        *
    from base
    where (
             followers_rank_by_cat < 1000 or
             followers_change_rank_by_cat < 1000 or
             views_rank_by_cat < 1000 or
             followers_rank_by_lang < 1000 or
             followers_change_rank_by_lang < 1000 or
             views_rank_by_lang < 1000 or
             followers_rank_by_cat_lang < 1000 or
             followers_change_rank_by_cat_lang < 1000 or
             views_rank_by_cat_lang < 1000 or
             views_rank_by_cat_profile < 1000 or
             followers_change_rank_by_cat_profile < 1000 or
             followers_rank_by_cat_profile < 1000 or
             views_rank_by_lang_profile < 1000 or
             followers_change_rank_by_lang_profile < 1000 or
             followers_rank_by_lang_profile < 1000 or
             views_rank_by_profile < 1000 or
             followers_change_rank_by_profile < 1000 or
             followers_rank_by_profile < 1000 or
             views_rank_by_cat_lang_profile < 1000 or
             followers_change_rank_by_cat_lang_profile < 1000 or
             followers_rank_by_cat_lang_profile < 1000
        )
    order by profile_platform desc, platform_id desc, month desc
),
with_prev_month_ranks as (
    select
        *,
        'ALL' platform,
        month,
        0 avg_likes,
        0 avg_views,
        0 engagement_rate,
        'IN' country,
        0 engagement,
        0 avg_comments,
        groupArray(month) OVER (PARTITION BY platform_id ORDER BY month ASC
            ROWS BETWEEN 1 PRECEDING AND 0 FOLLOWING) AS month_values,
        month_values[1] p_month,
        abs(p_month - month) <= 31 and length(month_values) = 2 last_month_available,

        groupArray(followers_rank) OVER (PARTITION BY platform_id ORDER BY month ASC
        Rows BETWEEN 1 PRECEDING AND 0 FOLLOWING) AS followers_rank_value,
        if(last_month_available = 1, followers_rank_value[1], NULL) followers_rank_prev,

        groupArray(views_rank) OVER (PARTITION BY platform_id ORDER BY month ASC
        Rows BETWEEN 1 PRECEDING AND 0 FOLLOWING) AS views_rank_value,
        if(last_month_available = 1, views_rank_value[1], NULL) views_rank_prev,

        groupArray(followers_change_rank) OVER (PARTITION BY platform_id ORDER BY month ASC
        Rows BETWEEN 1 PRECEDING AND 0 FOLLOWING) AS followers_change_rank_value,
        if(last_month_available = 1, followers_change_rank_value[1], NULL) followers_change_rank_prev,

        groupArray(followers_rank_by_cat) OVER (PARTITION BY platform_id ORDER BY month ASC
        Rows BETWEEN 1 PRECEDING AND 0 FOLLOWING) AS followers_rank_by_cat_value,
        if(last_month_available = 1, followers_rank_by_cat_value[1], NULL) followers_rank_by_cat_prev,

        groupArray(views_rank_by_cat) OVER (PARTITION BY platform_id ORDER BY month ASC
        Rows BETWEEN 1 PRECEDING AND 0 FOLLOWING) AS views_rank_by_cat_value,
        if(last_month_available = 1, views_rank_by_cat_value[1], NULL) views_rank_by_cat_prev,

        groupArray(followers_change_rank_by_cat) OVER (PARTITION BY platform_id ORDER BY month ASC
        Rows BETWEEN 1 PRECEDING AND 0 FOLLOWING) AS followers_change_rank_by_cat_value,
        if(last_month_available = 1, followers_change_rank_by_cat_value[1], NULL) followers_change_rank_by_cat_prev,

        groupArray(followers_rank_by_lang) OVER (PARTITION BY platform_id ORDER BY month ASC
        Rows BETWEEN 1 PRECEDING AND 0 FOLLOWING) AS followers_rank_by_lang_value,
        if(last_month_available = 1, followers_rank_by_lang_value[1], NULL) followers_rank_by_lang_prev,

        groupArray(views_rank_by_lang) OVER (PARTITION BY platform_id ORDER BY month ASC
        Rows BETWEEN 1 PRECEDING AND 0 FOLLOWING) AS views_rank_by_lang_value,
        if(last_month_available = 1, views_rank_by_lang_value[1], NULL) views_rank_by_lang_prev,

        groupArray(followers_change_rank_by_lang) OVER (PARTITION BY platform_id ORDER BY month ASC
        Rows BETWEEN 1 PRECEDING AND 0 FOLLOWING) AS followers_change_rank_by_lang_value,
        if(last_month_available = 1, followers_change_rank_by_lang_value[1], NULL) followers_change_rank_by_lang_prev,

        groupArray(followers_rank_by_cat_lang) OVER (PARTITION BY platform_id ORDER BY month ASC
        Rows BETWEEN 1 PRECEDING AND 0 FOLLOWING) AS followers_rank_by_cat_lang_value,
        if(last_month_available = 1, followers_rank_by_cat_lang_value[1], NULL) followers_rank_by_cat_lang_prev,

        groupArray(views_rank_by_cat_lang) OVER (PARTITION BY platform_id ORDER BY month ASC
        Rows BETWEEN 1 PRECEDING AND 0 FOLLOWING) AS views_rank_by_cat_lang_value,
        if(last_month_available = 1, views_rank_by_cat_lang_value[1], NULL) views_rank_by_cat_lang_prev,

        groupArray(followers_change_rank_by_cat_lang) OVER (PARTITION BY platform_id ORDER BY month ASC
        Rows BETWEEN 1 PRECEDING AND 0 FOLLOWING) AS followers_change_rank_by_cat_lang_value,
        if(last_month_available = 1, followers_change_rank_by_cat_lang_value[1], NULL) followers_change_rank_by_cat_lang_prev,

        groupArray(followers_rank_by_profile) OVER (PARTITION BY platform_id ORDER BY month ASC
        Rows BETWEEN 1 PRECEDING AND 0 FOLLOWING) AS followers_rank_by_profile_value,
        if(last_month_available = 1, followers_rank_by_profile_value[1], NULL) followers_rank_by_profile_prev,

        groupArray(followers_change_rank_by_profile) OVER (PARTITION BY platform_id ORDER BY month ASC
        Rows BETWEEN 1 PRECEDING AND 0 FOLLOWING) AS followers_change_rank_by_profile_value,
        if(last_month_available = 1, followers_change_rank_by_profile_value[1], NULL) followers_change_rank_by_profile_prev,

        groupArray(views_rank_by_profile) OVER (PARTITION BY platform_id ORDER BY month ASC
        Rows BETWEEN 1 PRECEDING AND 0 FOLLOWING) AS views_rank_by_profile_value,
        if(last_month_available = 1, views_rank_by_profile_value[1], NULL) views_rank_by_profile_prev,

        groupArray(followers_rank_by_cat_profile) OVER (PARTITION BY platform_id ORDER BY month ASC
        Rows BETWEEN 1 PRECEDING AND 0 FOLLOWING) AS followers_rank_by_cat_profile_value,
        if(last_month_available = 1, followers_rank_by_cat_profile_value[1], NULL) followers_rank_by_cat_profile_prev,

        groupArray(followers_change_rank_by_cat_profile) OVER (PARTITION BY platform_id ORDER BY month ASC
        Rows BETWEEN 1 PRECEDING AND 0 FOLLOWING) AS followers_change_rank_by_cat_profile_value,
        if(last_month_available = 1, followers_change_rank_by_cat_profile_value[1], NULL) followers_change_rank_by_cat_profile_prev,

        groupArray(views_rank_by_cat_profile) OVER (PARTITION BY platform_id ORDER BY month ASC
        Rows BETWEEN 1 PRECEDING AND 0 FOLLOWING) AS views_rank_by_cat_profile_value,
        if(last_month_available = 1, views_rank_by_cat_profile_value[1], NULL) views_rank_by_cat_profile_prev,

        groupArray(followers_rank_by_lang_profile) OVER (PARTITION BY platform_id ORDER BY month ASC
        Rows BETWEEN 1 PRECEDING AND 0 FOLLOWING) AS followers_rank_by_lang_profile_value,
        if(last_month_available = 1, followers_rank_by_lang_profile_value[1], NULL) followers_rank_by_lang_profile_prev,

        groupArray(followers_change_rank_by_lang_profile) OVER (PARTITION BY platform_id ORDER BY month ASC
        Rows BETWEEN 1 PRECEDING AND 0 FOLLOWING) AS followers_change_rank_by_lang_profile_value,
        if(last_month_available = 1, followers_change_rank_by_lang_profile_value[1], NULL) followers_change_rank_by_lang_profile_prev,

        groupArray(views_rank_by_lang_profile) OVER (PARTITION BY platform_id ORDER BY month ASC
        Rows BETWEEN 1 PRECEDING AND 0 FOLLOWING) AS views_rank_by_lang_profile_value,
        if(last_month_available = 1, views_rank_by_lang_profile_value[1], NULL) views_rank_by_lang_profile_prev,

        groupArray(followers_rank_by_cat_lang_profile) OVER (PARTITION BY platform_id ORDER BY month ASC
        Rows BETWEEN 1 PRECEDING AND 0 FOLLOWING) AS followers_rank_by_cat_lang_profile_value,
        if(last_month_available = 1, followers_rank_by_cat_lang_profile_value[1], NULL) followers_rank_by_cat_lang_profile_prev,

        groupArray(followers_change_rank_by_cat_lang_profile) OVER (PARTITION BY platform_id ORDER BY month ASC
        Rows BETWEEN 1 PRECEDING AND 0 FOLLOWING) AS followers_change_rank_by_cat_lang_profile_value,
        if(last_month_available = 1, followers_change_rank_by_cat_lang_profile_value[1], NULL) followers_change_rank_by_cat_lang_profile_prev,

        groupArray(views_rank_by_cat_lang_profile) OVER (PARTITION BY platform_id ORDER BY month ASC
        Rows BETWEEN 1 PRECEDING AND 0 FOLLOWING) AS views_rank_by_cat_lang_profile_value,
        if(last_month_available = 1, views_rank_by_cat_lang_profile_value[1], NULL) views_rank_by_cat_lang_profile_prev

    from sorted
)
select * from with_prev_month_ranks
SETTINGS max_bytes_before_external_group_by = 20000000000, max_bytes_before_external_sort = 20000000000, join_algorithm='partial_merge'

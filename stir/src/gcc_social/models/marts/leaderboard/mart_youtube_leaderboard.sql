{{ config(materialized = 'table', tags=["daily", "leaderboard"], order_by='month, platform_id') }}
with in_base as (
    select
        'YOUTUBE' platform,
        ifNull(month, date('2000-01-01')) month,
        ifNull(base.channel_id, 'MISSING') platform_id,
        ya.category category,
        ya.language language,
        ya.country country,
        COALESCE(ya.profile_type, 'CREATOR') profile_type,
        base.views views,
        base.prev_views prev_views,
        0 ia_views,
        base.views yt_views,
        base.uploads uploads,
        base.followers followers,
        base.prev_followers prev_followers,
        base.followers_change followers_change,
        1 enabled
    from
        {{ref('mart_yt_leaderboard_base')}} base
        left join {{ ref('mart_youtube_account') }} ya on ya.channelid = base.channel_id
    where country = 'IN' and prev_followers != 0 and views > 0
),
base as (
    select
        platform,
        month,
        platform_id,
        category,
        language,
        country,
        profile_type,
        views,
        prev_views,
        ia_views,
        yt_views,
        uploads,
        followers,
        prev_followers,
        followers_change,
        0 engagement_rate,
        0 avg_likes,
        0 avg_comments,
        0 avg_views,
        0 likes,
        0 comments,
        0 engagement,
        rank() over (partition by month order by followers desc, followers_change desc, views desc, platform_id desc) followers_rank,
        rank() over (partition by month order by followers_change desc, followers desc, views desc, platform_id desc) followers_change_rank,
        rank() over (partition by month order by views desc, followers desc, followers_change desc) views_rank,

        rank() over (partition by month, category order by followers desc, followers_change desc, views desc, platform_id desc) followers_rank_by_cat,
        rank() over (partition by month, category order by followers_change desc, followers desc, views desc, platform_id desc) followers_change_rank_by_cat,
        rank() over (partition by month, category order by views desc, followers desc, followers_change desc, platform_id desc) views_rank_by_cat,

        rank() over (partition by month, language order by followers desc, followers_change desc, views desc, platform_id desc) followers_rank_by_lang,
        rank() over (partition by month, language order by followers_change desc, followers desc, views desc, platform_id desc) followers_change_rank_by_lang,
        rank() over (partition by month, language order by views desc, followers desc, followers_change desc, platform_id desc) views_rank_by_lang,

        rank() over (partition by month, category, language order by followers desc, followers_change desc, views desc, platform_id desc) followers_rank_by_cat_lang,
        rank() over (partition by month, category, language order by followers_change desc, followers desc, views desc, platform_id desc) followers_change_rank_by_cat_lang,
        rank() over (partition by month, category, language order by views desc, followers desc, followers_change desc, platform_id desc) views_rank_by_cat_lang,

        rank() over (partition by month, profile_type order by followers desc, followers_change desc, views desc, platform_id desc) followers_rank_by_profile,
        rank() over (partition by month, profile_type order by followers_change desc, followers desc, views desc, platform_id desc) followers_change_rank_by_profile,
        rank() over (partition by month, profile_type order by views desc, followers desc, followers_change desc, platform_id desc) views_rank_by_profile,


        rank() over (partition by month, category, profile_type order by followers desc, followers_change desc, views desc, platform_id desc) followers_rank_by_cat_profile,
        rank() over (partition by month, category, profile_type order by followers_change desc, followers desc, views desc, platform_id desc) followers_change_rank_by_cat_profile,
        rank() over (partition by month, category, profile_type order by views desc, followers desc, followers_change desc, platform_id desc) views_rank_by_cat_profile,


        rank() over (partition by month, language, profile_type order by followers desc, followers_change desc, views desc, platform_id desc) followers_rank_by_lang_profile,
        rank() over (partition by month, language, profile_type order by followers_change desc, followers desc, views desc, platform_id desc) followers_change_rank_by_lang_profile,
        rank() over (partition by month, language, profile_type order by views desc, followers desc, followers_change desc, platform_id desc) views_rank_by_lang_profile,


        rank() over (partition by month, category, language, profile_type order by followers desc, followers_change desc, views desc, platform_id desc) followers_rank_by_cat_lang_profile,
        rank() over (partition by month, category, language, profile_type order by followers_change desc, followers desc, views desc, platform_id desc) followers_change_rank_by_cat_lang_profile,
        rank() over (partition by month, category, language, profile_type order by views desc, followers desc, followers_change desc, platform_id desc) views_rank_by_cat_lang_profile
    from in_base
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
    order by platform desc, platform_id desc, month desc
),
with_prev_month_ranks as (
    select
        *,
        month,
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
        if(last_month_available = 1 ,followers_change_rank_by_cat_lang_value[1], NULL) followers_change_rank_by_cat_lang_prev,

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

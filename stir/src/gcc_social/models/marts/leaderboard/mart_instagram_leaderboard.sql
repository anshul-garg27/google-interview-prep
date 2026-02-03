{{ config(materialized = 'table', tags=["daily", "leaderboard"], order_by='month, platform_id') }}
with
handle_info as (
    select ig_id profile_id,
           argMax(handle, updated_at) handle,
           argMax(primary_category, updated_at) primary_category,
           argMax(primary_language, updated_at) primary_language,
           argMax(country, updated_at) country,
           COALESCE(argMax(profile_type, updated_at), 'CREATOR') profile_type
    from dbt.mart_instagram_account
    group by ig_id
),
base as (
    select
        'INSTAGRAM' platform,
        ifNull(month, date('2000-01-01')) month,
        ifNull(base.profile_id, 'MISSING') platform_id,
        ia.handle handle,
        ia.primary_category category,
        ia.primary_language language,
        ia.country country,
        ia.profile_type profile_type,
        base.views views,
        base.views ia_views,
        0 yt_views,
        base.followers followers,
        base.prev_followers prev_followers,
        base.followers_change followers_change,
        base.views views,
        base.uploads uploads,
        base.likes likes,
        base.comments comments,
        base.plays plays,
        base.prev_plays prev_plays,
        base.engagement engagement,
        base.avg_views avg_views,
        base.avg_likes avg_likes,
        base.avg_plays avg_plays,
        base.avg_comments avg_comments,
        base.avg_er engagement_rate,
        1 enabled,

        rank() over (partition by month order by base.followers desc) followers_rank,
        rank() over (partition by month order by base.followers_change desc) followers_change_rank,
        rank() over (partition by month order by base.views desc) views_rank,
        rank() over (partition by month order by base.plays desc) plays_rank,


        rank() over (partition by month, category order by base.followers desc) followers_rank_by_cat,
        rank() over (partition by month, category order by base.followers_change desc) followers_change_rank_by_cat,
        rank() over (partition by month, category order by base.views desc) views_rank_by_cat,
        rank() over (partition by month, category order by base.plays desc) plays_rank_by_cat,


        rank() over (partition by month, language order by base.followers desc) followers_rank_by_lang,
        rank() over (partition by month, language order by base.followers_change desc) followers_change_rank_by_lang,
        rank() over (partition by month, language order by base.views desc) views_rank_by_lang,
        rank() over (partition by month, language order by base.plays desc) plays_rank_by_lang,


        rank() over (partition by month, category, language order by base.followers desc) followers_rank_by_cat_lang,
        rank() over (partition by month, category, language order by base.followers_change desc) followers_change_rank_by_cat_lang,
        rank() over (partition by month, category, language order by base.views desc) views_rank_by_cat_lang,
        rank() over (partition by month, category, language order by base.plays desc) plays_rank_by_cat_lang,

        rank() over (partition by month, profile_type order by base.followers desc) followers_rank_by_profile,
        rank() over (partition by month, profile_type order by base.followers_change desc) followers_change_rank_by_profile,
        rank() over (partition by month, profile_type order by base.views desc) views_rank_by_profile,
        rank() over (partition by month, profile_type order by base.plays desc) plays_rank_by_profile,


        rank() over (partition by month, category, profile_type order by base.followers desc) followers_rank_by_cat_profile,
        rank() over (partition by month, category, profile_type order by base.followers_change desc) followers_change_rank_by_cat_profile,
        rank() over (partition by month, category, profile_type order by base.views desc) views_rank_by_cat_profile,
        rank() over (partition by month, category, profile_type order by base.plays desc) plays_rank_by_cat_profile,


        rank() over (partition by month, language, profile_type order by base.followers desc) followers_rank_by_lang_profile,
        rank() over (partition by month, language, profile_type order by base.followers_change desc) followers_change_rank_by_lang_profile,
        rank() over (partition by month, language, profile_type order by base.views desc) views_rank_by_lang_profile,
        rank() over (partition by month, language, profile_type order by base.plays desc) plays_rank_by_lang_profile,


        rank() over (partition by month, category, language, profile_type order by base.followers desc) followers_rank_by_cat_lang_profile,
        rank() over (partition by month, category, language, profile_type order by base.followers_change desc) followers_change_rank_by_cat_lang_profile,
        rank() over (partition by month, category, language, profile_type order by base.views desc) views_rank_by_cat_lang_profile,
        rank() over (partition by month, category, language, profile_type order by base.plays desc) plays_rank_by_cat_lang_profile

    from
        {{ ref('mart_insta_leaderboard_base') }} base
        left join handle_info ia on ia.profile_id = base.profile_id
        left join dbt.stg_vidooly_instagram_profiles via on via.profile_id = base.profile_id
        where ia.handle != '' and ia.country = 'IN' and prev_followers != 0 and plays > 0
        order by month desc, followers_rank asc
),
sorted as (
    select
        *
    from base
    where (
             followers_rank_by_cat < 1000 or
             followers_change_rank_by_cat < 1000 or
             views_rank_by_cat < 1000 or
             plays_rank_by_cat < 1000 or
             followers_rank_by_lang < 1000 or
             followers_change_rank_by_lang < 1000 or
             views_rank_by_lang < 1000 or
             plays_rank_by_lang < 1000 or
             followers_rank_by_cat_lang < 1000 or
             followers_change_rank_by_cat_lang < 1000 or
             views_rank_by_cat_lang < 1000 or
             plays_rank_by_cat_lang < 1000 or
             plays_rank_by_cat_profile < 1000 or
             views_rank_by_cat_profile < 1000 or
             followers_change_rank_by_cat_profile < 1000 or
             followers_rank_by_cat_profile < 1000 or
             plays_rank_by_lang_profile < 1000 or
             views_rank_by_lang_profile < 1000 or
             followers_change_rank_by_lang_profile < 1000 or
             followers_rank_by_lang_profile < 1000 or
             plays_rank_by_profile < 1000 or
             views_rank_by_profile < 1000 or
             followers_change_rank_by_profile < 1000 or
             followers_rank_by_profile < 1000 or
             plays_rank_by_cat_lang_profile < 1000 or
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

        groupArray(plays_rank) OVER (PARTITION BY platform_id ORDER BY month ASC
        Rows BETWEEN 1 PRECEDING AND 0 FOLLOWING) AS plays_rank_value,
        if(last_month_available = 1, plays_rank_value[1], NULL) plays_rank_prev,

        groupArray(followers_rank_by_cat) OVER (PARTITION BY platform_id ORDER BY month ASC
        Rows BETWEEN 1 PRECEDING AND 0 FOLLOWING) AS followers_rank_by_cat_value,
        if(last_month_available = 1, followers_rank_by_cat_value[1], NULL) followers_rank_by_cat_prev,

        groupArray(views_rank_by_cat) OVER (PARTITION BY platform_id ORDER BY month ASC
        Rows BETWEEN 1 PRECEDING AND 0 FOLLOWING) AS views_rank_by_cat_value,
        if(last_month_available = 1, views_rank_by_cat_value[1], NULL) views_rank_by_cat_prev,

        groupArray(followers_change_rank_by_cat) OVER (PARTITION BY platform_id ORDER BY month ASC
        Rows BETWEEN 1 PRECEDING AND 0 FOLLOWING) AS followers_change_rank_by_cat_value,
        if(last_month_available = 1, followers_change_rank_by_cat_value[1], NULL) followers_change_rank_by_cat_prev,

        groupArray(plays_rank_by_cat) OVER (PARTITION BY platform_id ORDER BY month ASC
        Rows BETWEEN 1 PRECEDING AND 0 FOLLOWING) AS plays_rank_by_cat_value,
        if(last_month_available = 1, plays_rank_by_cat_value[1], NULL) plays_rank_by_cat_prev,

        groupArray(followers_rank_by_lang) OVER (PARTITION BY platform_id ORDER BY month ASC
        Rows BETWEEN 1 PRECEDING AND 0 FOLLOWING) AS followers_rank_by_lang_value,
        if(last_month_available = 1, followers_rank_by_lang_value[1], NULL) followers_rank_by_lang_prev,

        groupArray(views_rank_by_lang) OVER (PARTITION BY platform_id ORDER BY month ASC
        Rows BETWEEN 1 PRECEDING AND 0 FOLLOWING) AS views_rank_by_lang_value,
        if(last_month_available = 1, views_rank_by_lang_value[1], NULL) views_rank_by_lang_prev,

        groupArray(followers_change_rank_by_lang) OVER (PARTITION BY platform_id ORDER BY month ASC
        Rows BETWEEN 1 PRECEDING AND 0 FOLLOWING) AS followers_change_rank_by_lang_value,
        if(last_month_available = 1, followers_change_rank_by_lang_value[1], NULL) followers_change_rank_by_lang_prev,

        groupArray(plays_rank_by_lang) OVER (PARTITION BY platform_id ORDER BY month ASC
        Rows BETWEEN 1 PRECEDING AND 0 FOLLOWING) AS plays_rank_by_lang_value,
        if(last_month_available = 1, plays_rank_by_lang_value[1], NULL) plays_rank_by_lang_prev,

        groupArray(followers_rank_by_cat_lang) OVER (PARTITION BY platform_id ORDER BY month ASC
        Rows BETWEEN 1 PRECEDING AND 0 FOLLOWING) AS followers_rank_by_cat_lang_value,
        if(last_month_available = 1, followers_rank_by_cat_lang_value[1], NULL) followers_rank_by_cat_lang_prev,

        groupArray(views_rank_by_cat_lang) OVER (PARTITION BY platform_id ORDER BY month ASC
        Rows BETWEEN 1 PRECEDING AND 0 FOLLOWING) AS views_rank_by_cat_lang_value,
        if(last_month_available = 1, views_rank_by_cat_lang_value[1], NULL) views_rank_by_cat_lang_prev,

        groupArray(followers_change_rank_by_cat_lang) OVER (PARTITION BY platform_id ORDER BY month ASC
        Rows BETWEEN 1 PRECEDING AND 0 FOLLOWING) AS followers_change_rank_by_cat_lang_value,
        if(last_month_available = 1, followers_change_rank_by_cat_lang_value[1], NULL) followers_change_rank_by_cat_lang_prev,

        groupArray(plays_rank_by_cat_lang) OVER (PARTITION BY platform_id ORDER BY month ASC
        Rows BETWEEN 1 PRECEDING AND 0 FOLLOWING) AS plays_rank_by_cat_lang_value,
        if(last_month_available = 1, plays_rank_by_cat_lang_value[1], NULL) plays_rank_by_cat_lang_prev,

        groupArray(followers_rank_by_profile) OVER (PARTITION BY platform_id ORDER BY month ASC
        Rows BETWEEN 1 PRECEDING AND 0 FOLLOWING) AS followers_rank_by_profile_value,
        if(last_month_available = 1, followers_rank_by_profile_value[1], NULL) followers_rank_by_profile_prev,

        groupArray(followers_change_rank_by_profile) OVER (PARTITION BY platform_id ORDER BY month ASC
        Rows BETWEEN 1 PRECEDING AND 0 FOLLOWING) AS followers_change_rank_by_profile_value,
        if(last_month_available = 1, followers_change_rank_by_profile_value[1], NULL) followers_change_rank_by_profile_prev,

        groupArray(views_rank_by_profile) OVER (PARTITION BY platform_id ORDER BY month ASC
        Rows BETWEEN 1 PRECEDING AND 0 FOLLOWING) AS views_rank_by_profile_value,
        if(last_month_available = 1, views_rank_by_profile_value[1], NULL) views_rank_by_profile_prev,

        groupArray(plays_rank_by_profile) OVER (PARTITION BY platform_id ORDER BY month ASC
        Rows BETWEEN 1 PRECEDING AND 0 FOLLOWING) AS plays_rank_by_profile_value,
        if(last_month_available = 1, plays_rank_by_profile_value[1], NULL) plays_rank_by_profile_prev,

        groupArray(followers_rank_by_cat_profile) OVER (PARTITION BY platform_id ORDER BY month ASC
        Rows BETWEEN 1 PRECEDING AND 0 FOLLOWING) AS followers_rank_by_cat_profile_value,
        if(last_month_available = 1, followers_rank_by_cat_profile_value[1], NULL) followers_rank_by_cat_profile_prev,

        groupArray(followers_change_rank_by_cat_profile) OVER (PARTITION BY platform_id ORDER BY month ASC
        Rows BETWEEN 1 PRECEDING AND 0 FOLLOWING) AS followers_change_rank_by_cat_profile_value,
        if(last_month_available = 1, followers_change_rank_by_cat_profile_value[1], NULL) followers_change_rank_by_cat_profile_prev,

        groupArray(views_rank_by_cat_profile) OVER (PARTITION BY platform_id ORDER BY month ASC
        Rows BETWEEN 1 PRECEDING AND 0 FOLLOWING) AS views_rank_by_cat_profile_value,
        if(last_month_available = 1, views_rank_by_cat_profile_value[1], NULL) views_rank_by_cat_profile_prev,

        groupArray(plays_rank_by_cat_profile) OVER (PARTITION BY platform_id ORDER BY month ASC
        Rows BETWEEN 1 PRECEDING AND 0 FOLLOWING) AS plays_rank_by_cat_profile_value,
        if(last_month_available = 1, plays_rank_by_cat_profile_value[1], NULL) plays_rank_by_cat_profile_prev,

        groupArray(followers_rank_by_lang_profile) OVER (PARTITION BY platform_id ORDER BY month ASC
        Rows BETWEEN 1 PRECEDING AND 0 FOLLOWING) AS followers_rank_by_lang_profile_value,
        if(last_month_available = 1, followers_rank_by_lang_profile_value[1], NULL) followers_rank_by_lang_profile_prev,

        groupArray(followers_change_rank_by_lang_profile) OVER (PARTITION BY platform_id ORDER BY month ASC
        Rows BETWEEN 1 PRECEDING AND 0 FOLLOWING) AS followers_change_rank_by_lang_profile_value,
        if(last_month_available = 1, followers_change_rank_by_lang_profile_value[1], NULL) followers_change_rank_by_lang_profile_prev,

        groupArray(views_rank_by_lang_profile) OVER (PARTITION BY platform_id ORDER BY month ASC
        Rows BETWEEN 1 PRECEDING AND 0 FOLLOWING) AS views_rank_by_lang_profile_value,
        if(last_month_available = 1, views_rank_by_lang_profile_value[1], NULL) views_rank_by_lang_profile_prev,

        groupArray(plays_rank_by_lang_profile) OVER (PARTITION BY platform_id ORDER BY month ASC
        Rows BETWEEN 1 PRECEDING AND 0 FOLLOWING) AS plays_rank_by_lang_profile_value,
        if(last_month_available = 1, plays_rank_by_lang_profile_value[1], NULL) plays_rank_by_lang_profile_prev,

        groupArray(followers_rank_by_cat_lang_profile) OVER (PARTITION BY platform_id ORDER BY month ASC
        Rows BETWEEN 1 PRECEDING AND 0 FOLLOWING) AS followers_rank_by_cat_lang_profile_value,
        if(last_month_available = 1, followers_rank_by_cat_lang_profile_value[1], NULL) followers_rank_by_cat_lang_profile_prev,

        groupArray(followers_change_rank_by_cat_lang_profile) OVER (PARTITION BY platform_id ORDER BY month ASC
        Rows BETWEEN 1 PRECEDING AND 0 FOLLOWING) AS followers_change_rank_by_cat_lang_profile_value,
        if(last_month_available = 1, followers_change_rank_by_cat_lang_profile_value[1], NULL) followers_change_rank_by_cat_lang_profile_prev,

        groupArray(views_rank_by_cat_lang_profile) OVER (PARTITION BY platform_id ORDER BY month ASC
        Rows BETWEEN 1 PRECEDING AND 0 FOLLOWING) AS views_rank_by_cat_lang_profile_value,
        if(last_month_available = 1, views_rank_by_cat_lang_profile_value[1], NULL) views_rank_by_cat_lang_profile_prev,

        groupArray(plays_rank_by_cat_lang_profile) OVER (PARTITION BY platform_id ORDER BY month ASC
        Rows BETWEEN 1 PRECEDING AND 0 FOLLOWING) AS plays_rank_by_cat_lang_profile_value,
        if(last_month_available = 1, plays_rank_by_cat_lang_profile_value[1], NULL) plays_rank_by_cat_lang_profile_prev
        
    from sorted
)
select * from with_prev_month_ranks
SETTINGS max_bytes_before_external_group_by = 20000000000, max_bytes_before_external_sort = 20000000000, join_algorithm='partial_merge'

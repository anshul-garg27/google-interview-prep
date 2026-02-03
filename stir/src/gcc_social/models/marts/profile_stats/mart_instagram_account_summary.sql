{{ config(materialized = 'table', tags=["post_ranker"], order_by='profile_id') }}
with
recent_posts_by_class as (
    select
        profile_id,
        followers,
        following,
        post_shortcode,
        post_type,
        likes,
        comments,
        views,
        plays,
        engagement,
        impressions,
        reach,
        er,
        post_class,
        publish_time,
        post_rank,
        post_rank_by_class,
        rank() over (partition by profile_id, post_class order by likes desc, comments desc, publish_time desc) likes_rank_by_class,
        rank() over (partition by profile_id, post_class order by comments desc, likes desc, publish_time desc) comments_rank_by_class,
        rank() over (partition by profile_id, post_class order by plays desc, likes desc, comments desc, publish_time desc) plays_rank_by_class,
        rank() over (partition by profile_id, post_class order by views desc, likes desc, comments desc, publish_time desc) views_rank_by_class
    from {{ref('mart_instagram_recent_post_stats')}}
    where post_rank_by_class <= 12
),
recent_posts as (
    select
        post_shortcode,
        rank() over (partition by profile_id order by likes desc, comments desc, publish_time desc) likes_rank,
        rank() over (partition by profile_id order by comments desc, likes desc, publish_time desc) comments_rank
    from {{ref('mart_instagram_recent_post_stats')}}
    where post_rank <= 12
),
all_recent_posts as (
    select
        post_shortcode,
        max(profile_id) profile_id,
        max(followers) followers,
        max(following) following,
        max(post_class) post_class,
        max(likes) likes,
        max(views) views,
        max(plays) plays,
        max(comments) comments,
        max(engagement) engagement,
        max(impressions) impressions,
        max(reach) reach,
        max(post_rank) post_rank,
        max(likes_rank) likes_rank,
        max(comments_rank) comments_rank,
        max(likes_rank_by_class) likes_rank_by_class,
        max(comments_rank_by_class) comments_rank_by_class,
        max(plays_rank_by_class) plays_rank_by_class,
        max(views_rank_by_class) views_rank_by_class,
        max(post_rank_by_class) post_rank_by_class
    from recent_posts_by_class
    left join recent_posts on recent_posts.post_shortcode = recent_posts_by_class.post_shortcode
    group by post_shortcode
    order by profile_id asc, post_class asc, post_rank asc
),
profile_summary as (
    select
        profile_id,

        max(followers) AS followers,
        max(following) AS following,

        countIf(post_shortcode, post_rank <= 12) recent_posts,
        countIf(post_shortcode, post_class='static') recent_statics,
        countIf(post_shortcode, post_class='reels') recent_reels,

        sumIf(engagement, post_rank <= 12) / followers recent_er,
        sumIf(impressions, post_rank <= 12) recent_impressions,
        sumIf(reach, post_rank <= 12) recent_reach,
        sumIf(likes, post_rank <= 12) recent_likes,
        sumIf(comments, post_rank <= 12) recent_comments,
        sumIf(reach, post_rank <= 12) recent_reach,

        recent_likes / recent_posts avg_likes,
        recent_comments / recent_posts avg_comments,
        recent_impressions / recent_posts avg_impressions,
        recent_reach / recent_posts avg_reach,
        recent_er / recent_posts avg_er,
        recent_reach / recent_posts avg_reach,
        recent_impressions / recent_posts avg_impressions,

        (100 * maxIf(likes, post_rank <= 12)) / avg_likes likes_spread,

        sumIf(comments, post_class = 'reels') recent_reels_comments,
        sumIf(likes, post_class = 'reels') recent_reels_likes,
        sumIf(views, post_class = 'reels') recent_reels_views,
        sumIf(plays, post_class = 'reels') recent_reels_plays,
        sumIf(reach, post_class = 'reels') recent_reels_reach,
        sumIf(impressions, post_class = 'reels') recent_reels_impressions,
        (recent_reels_comments +  recent_reels_plays) / followers recent_reels_er,

        recent_reels_comments / recent_reels avg_reels_comments,
        recent_reels_likes / recent_reels avg_reels_likes,
        recent_reels_views / recent_reels avg_reels_views,
        recent_reels_plays / recent_reels avg_reels_plays,
        recent_reels_er / recent_reels avg_reels_er,
        recent_reels_impressions / recent_reels avg_reels_impressions,
        recent_reels_reach / recent_reels avg_reels_reach,

        sumIf(comments, post_class = 'static') recent_static_comments,
        sumIf(likes, post_class = 'static') recent_static_likes,
        sumIf(engagement, post_class='static') / followers recent_static_er,
        sumIf(reach, post_class='static') recent_static_reach,
        sumIf(impressions, post_class='static') recent_static_impressions,

        recent_static_comments / recent_statics avg_static_comments,
        recent_static_likes / recent_statics avg_static_likes,
        recent_static_er / recent_statics avg_static_er,
        recent_static_reach / recent_statics avg_static_reach,
        recent_static_impressions / recent_statics avg_static_impressions

    from all_recent_posts
    group by profile_id
),
profile_summary_with_outliers as (
    select
        profile_id,
        max(followers) followers,
        max(following) following,

        countIf(post_shortcode, likes_rank >= 3 and likes_rank <= ps.recent_posts - 2) non_outliers,

        sumIf(comments, likes_rank >= 3 and likes_rank <= ps.recent_posts - 2) comments_non_outliers,
        sumIf(likes, likes_rank >= 3 and likes_rank <= ps.recent_posts - 2) likes_non_outliers,
        sumIf(reach, likes_rank >= 3 and likes_rank <= ps.recent_posts - 2) reach_non_outliers,
        sumIf(impressions, likes_rank >= 3 and likes_rank <= ps.recent_posts - 2) impressions_non_outliers,
        (comments_non_outliers + likes_non_outliers) / followers er_non_outliers,

        comments_non_outliers / non_outliers avg_comments_non_outliers,
        likes_non_outliers / non_outliers avg_likes_non_outliers,
        er_non_outliers / non_outliers avg_er_non_outliers,
        impressions_non_outliers / non_outliers avg_impressions_non_outliers,
        reach_non_outliers / non_outliers avg_reach_non_outliers,

        countIf(post_shortcode, post_class='reels' and plays_rank_by_class >= 3 and plays_rank_by_class <= ps.recent_reels - 2) reels_non_outliers,
        sumIf(comments, post_class = 'reels' and plays_rank_by_class >= 3 and plays_rank_by_class <= ps.recent_reels - 2) reels_comments_non_outliers,
        sumIf(likes, post_class = 'reels' and plays_rank_by_class >= 3 and plays_rank_by_class <= ps.recent_reels - 2) reels_likes_non_outliers,
        sumIf(plays, post_class = 'reels' and plays_rank_by_class >= 3 and plays_rank_by_class <= ps.recent_reels - 2) reels_plays_non_outliers,
        sumIf(views, post_class = 'reels' and plays_rank_by_class >= 3 and plays_rank_by_class <= ps.recent_reels - 2) reels_views_non_outliers,
        sumIf(reach, post_class = 'reels' and plays_rank_by_class >= 3 and plays_rank_by_class <= ps.recent_reels - 2) reels_reach_non_outliers,
        sumIf(impressions, post_class = 'reels' and plays_rank_by_class >= 3 and plays_rank_by_class <= ps.recent_reels - 2) reels_impressions_non_outliers,

        (reels_comments_non_outliers + reels_likes_non_outliers) / followers reels_er_non_outliers,

        reels_comments_non_outliers / reels_non_outliers avg_reels_comments_non_outliers,
        reels_likes_non_outliers / reels_non_outliers avg_reels_likes_non_outliers,
        reels_views_non_outliers / reels_non_outliers avg_reels_views_non_outliers,
        reels_plays_non_outliers / reels_non_outliers avg_reels_plays_non_outliers,
        reels_er_non_outliers / reels_non_outliers avg_reels_er_non_outliers,
        reels_reach_non_outliers / reels_non_outliers avg_reels_reach_non_outliers,
        reels_impressions_non_outliers / reels_non_outliers avg_reels_impressions_non_outliers,

        countIf(post_shortcode, post_class='static' and plays_rank_by_class >= 3 and plays_rank_by_class <= ps.recent_statics - 2) static_non_outliers,
        sumIf(comments, post_class = 'static' and plays_rank_by_class >= 3 and plays_rank_by_class <= ps.recent_statics - 2) static_comments_non_outliers,
        sumIf(likes, post_class = 'static' and plays_rank_by_class >= 3 and plays_rank_by_class <= ps.recent_statics - 2) static_likes_non_outliers,
        sumIf(plays, post_class = 'static' and plays_rank_by_class >= 3 and plays_rank_by_class <= ps.recent_statics - 2) static_plays_non_outliers,
        sumIf(views, post_class = 'static' and plays_rank_by_class >= 3 and plays_rank_by_class <= ps.recent_statics - 2) static_views_non_outliers,
        sumIf(reach, post_class = 'static' and plays_rank_by_class >= 3 and plays_rank_by_class <= ps.recent_statics - 2) static_reach_non_outliers,
        sumIf(impressions, post_class = 'static' and plays_rank_by_class >= 3 and plays_rank_by_class <= ps.recent_statics - 2) static_impressions_non_outliers,

        (static_comments_non_outliers + static_likes_non_outliers) / followers static_er_non_outliers,
        static_comments_non_outliers / static_non_outliers avg_static_comments_non_outliers,
        static_likes_non_outliers / static_non_outliers avg_static_likes_non_outliers,
        static_er_non_outliers / static_non_outliers avg_static_er_non_outliers,
        static_reach_non_outliers / static_non_outliers avg_static_reach_non_outliers,
        static_impressions_non_outliers / static_non_outliers avg_static_impressions_non_outliers

    from all_recent_posts r
    left join profile_summary ps on ps.profile_id = r.profile_id
    group by profile_id
),
final_summary as (
    select
        summary.profile_id profile_id,
        summary.followers followers,
        summary.following following,
        if (summary.following > 0 , summary.followers / summary.following, toFloat64(summary.followers)) ffratio,

        summary.recent_reels recent_reels,
        summary.recent_statics recent_statics,
        summary.recent_posts recent_posts,

        summary.avg_comments avg_comments,
        summary.avg_likes avg_likes,
        summary.avg_reach avg_reach,
        summary.avg_impressions avg_impressions,
        summary.avg_er avg_er,
        summary.likes_spread likes_spread_percentage,

        summary.avg_reels_comments avg_reels_comments,
        summary.avg_reels_likes avg_reels_likes,
        summary.avg_reels_views avg_reels_views,
        summary.avg_reels_reach avg_reels_reach,
        summary.avg_reels_impressions avg_reels_impressions,
        summary.avg_reels_plays avg_reels_plays,
        summary.avg_reels_er reels_er,

        summary.avg_static_comments avg_statics_comments,
        summary.avg_static_likes avg_statics_likes,
        summary.avg_static_reach avg_static_reach,
        summary.avg_static_impressions avg_static_impressions,
        summary.avg_static_er avg_static_er,


        so.non_outliers posts_non_outliers,
        so.static_non_outliers statics_non_outliers,
        so.reels_non_outliers reels_non_outliers,

        so.reach_non_outliers reach_non_outliers,
        so.likes_non_outliers likes_non_outliers,
        so.comments_non_outliers comments_non_outliers,
        so.impressions_non_outliers impressions_non_outliers,

        so.avg_reach_non_outliers avg_reach_non_outliers,
        so.avg_er_non_outliers avg_er_non_outliers,
        so.avg_likes_non_outliers avg_likes_non_outliers,
        so.avg_comments_non_outliers avg_comments_non_outliers,
        so.avg_impressions_non_outliers avg_impresions_non_outliers,

        so.avg_reels_comments_non_outliers avg_reels_comments_non_outliers,
        so.avg_reels_likes_non_outliers avg_reels_likes_non_outliers,
        so.avg_reels_plays_non_outliers avg_reels_plays_non_outliers,
        so.avg_reels_views_non_outliers avg_reels_views_non_outliers,
        so.avg_reels_reach_non_outliers avg_reels_reach_non_outliers,
        so.avg_reels_impressions_non_outliers avg_reels_impressions_non_outliers,
        so.reels_er_non_outliers reels_er_non_outliers,

        so.avg_static_comments_non_outliers avg_statics_comments_non_outliers,
        so.avg_static_likes_non_outliers avg_statics_likes_non_outliers,
        so.avg_static_er_non_outliers static_er_non_outliers,
        so.avg_static_reach_non_outliers static_reach_non_outliers,
        so.avg_static_impressions_non_outliers static_impressions_non_outliers

    from profile_summary summary
    left join profile_summary_with_outliers  so on so.profile_id = summary.profile_id
),
final as (
    select
        ifNull(final_summary.profile_id, 'MISSING') profile_id,
        final_summary.*,
        ia.media_count total_posts,
        round((-0.000025017) * ia.followers + (1.11 * (ia.followers * abs(log2(final_summary.avg_er * 100)) * 2 / 100))) story_reach
    from final_summary
    left join {{ ref('stg_beat_instagram_account') }} ia on ia.profile_id = final_summary.profile_id
    order by followers desc
)
select * from final SETTINGS max_bytes_before_external_group_by = 40000000000, max_bytes_before_external_sort = 40000000000


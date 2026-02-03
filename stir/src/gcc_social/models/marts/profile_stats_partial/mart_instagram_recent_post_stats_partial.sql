{{ config(
        materialized = 'table',
        tags=["post_ranker_partial"],
        order_by='post_shortcode'
    )
}}
with profile_post_ranking as (
    select
        ifNull(ip.post_shortcode, 'MISSING') post_shortcode,
        max(profile_id) profile_id,
        max(a.handle) handle,
        max(a.followers) followers,
        max(a.following) following,
        max(ip.post_type) post_type,
        max(ip.likes) likes,
        max(ip.comments) comments,
        max(ip.views) views,
        max(ip.plays) plays,
        max(ip.post_type) post_type,
        max(ip.publish_time) publish_time,
        max(post_rank) post_rank,
        max(post_rank_by_class) post_rank_by_class,
        if(post_type = 'reels', 'reels', 'static') post_class,
        likes + comments engagement,
        plays * (0.94 - (log2(followers) * 0.001)) _reach_reels,
        toFloat64(plays) _impressions_reels,
        (7.6 - (log10(likes) * 0.7)) * 0.85*likes _reach_non_reels,
        (7.6-(log10(likes) * 0.75))*0.95*likes _impressions_non_reels,
        if (post_class = 'reels', max2(_reach_reels,0), max2(_reach_non_reels,0)) reach,
        if (post_class = 'reels', max2(_impressions_reels,0), max2(_impressions_non_reels,0)) impressions,
        engagement / followers er
    from
        {{ ref('mart_instagram_all_posts_with_ranks_partial') }} ip
        left join {{ ref('stg_beat_instagram_account') }} a on a.profile_id = ip.profile_id
    where ip.post_rank_by_class <= 12
    group by post_shortcode
    order by handle asc, post_class desc, publish_time desc
)
select * from profile_post_ranking
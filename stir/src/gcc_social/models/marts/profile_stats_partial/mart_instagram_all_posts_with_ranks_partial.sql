{{ config(
        materialized = 'table',
        tags=["post_ranker_partial"],
        order_by='post_shortcode'
    )
}}
select
    r.post_shortcode post_shortcode,
    r.post_rank post_rank,
    r.post_rank_by_class post_rank_by_class,
    p.profile_id profile_id,
    p.handle handle,
    p.post_type post_type,
    p.likes likes,
    p.comments comments,
    p.views views,
    p.plays plays,
    p.engagement engagement,
    p.post_class post_class,
    p.publish_time publish_time
from
    {{ref('mart_instagram_post_ranks_partial')}} r
    left join {{ref('mart_instagram_all_posts_partial')}} p on  r.post_shortcode = p.post_shortcode

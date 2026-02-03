{{ config(
        materialized = 'table',
        tags=["post_ranker"],
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
    p.publish_time publish_time,
    p.reach reach,
    p.impressions impressions
from
    {{ref('mart_instagram_post_ranks')}} r
    left join {{ref('mart_instagram_all_posts')}} p on  r.post_shortcode = p.post_shortcode
    SETTINGS max_bytes_before_external_group_by = 40000000000, max_bytes_before_external_sort = 40000000000, join_algorithm='partial_merge'

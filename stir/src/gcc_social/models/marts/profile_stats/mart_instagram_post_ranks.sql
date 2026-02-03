{{ config(
        materialized = 'table',
        tags=["post_ranker"],
        order_by='post_shortcode'
    )
}}
with
all_posts as (
    select
        ifNull(post_shortcode, 'MISSING') post_shortcode,
        max(profile_id) profile_id,
        max(post_class) post_class,
        max(publish_time) publish_time
    from {{ ref('mart_instagram_all_posts') }}
    group by post_shortcode
    order by profile_id desc, publish_time desc
),
all_posts_with_ranks as (
    select
        post_shortcode,
        profile_id,
        post_class,
        publish_time,
        rank() over(partition by profile_id order by publish_time desc) post_rank,
        rank() over(partition by profile_id, post_class order by publish_time desc) post_rank_by_class
    from all_posts
)
select * from all_posts_with_ranks
    SETTINGS max_bytes_before_external_group_by = 30000000000, max_bytes_before_external_sort = 30000000000, join_algorithm='partial_merge'

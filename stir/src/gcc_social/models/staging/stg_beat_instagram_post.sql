{{ config(
            materialized = 'incremental',
            engine='ReplacingMergeTree',
            tags=["post_ranker_partial"],
            order_by='ifNull(id, 0)'
    )
}}
with
data as (
     select 
        id,
        profile_id,
        handle,
        post_id,
        post_type,
        caption,
        likes,
        views,
        plays,
        comments,
        reach,
        saved,
        shares,
        impressions,
        interactions,
        engagement,
        shortcode,
        display_url,
        thumbnail_url,
        updated_at,
        created_at,
        publish_time,
        content_type,
        taps_forward,
        taps_back,
        exits,
        replies,
        navigation,
        profile_activity,
        profile_visits,
        follows,
        video_views,
        automatic_forward,
        swipe_back,
        swipe_down,
        swipe_forward,
        swipe_up,
        profile_follower,
        profile_er,
        profile_static_er,
        profile_reels_er
     from beat_replica.instagram_post_simple
    {% if is_incremental() %}
        where updated_at > (select max(updated_at) - INTERVAL 4 HOUR from {{ this }})
    {% endif %}
)
select * from data

{{ config(materialized = 'table', tags=["daily", "genre"], order_by='month, short_code') }}
SELECT    
    platform,
    month,
    profile_id,
    short_code,
    title,
    description,
    views,
    likes,
    comments,
    engagement,
    engagement_rate,
    plays,
    views_rank,
    plays_rank,
    post_type,
    profile_type,
    post_url,
    thumbnail_url,
    publish_time,
    language,
    category
FROM {{ref('mart_genre_instagram_trending_content_with_rank')}}

UNION ALL

SELECT
    platform,
    month,
    profile_id,
    short_code,
    title,
    description,
    views,
    likes,
    comments,
    engagement,
    engagement_rate,
    plays,
    views_rank,
    plays_rank,
    post_type,
    profile_type,
    post_url,
    thumbnail_url,
    publish_time,
    language,
    category
FROM {{ref('mart_genre_youtube_trending_content_with_rank')}}
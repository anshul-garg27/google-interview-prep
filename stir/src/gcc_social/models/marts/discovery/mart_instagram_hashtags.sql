{{ config(materialized = 'table', tags=["weekly"], order_by='platform_id') }}
with hashtags as (
    select
        profile_id,
        arrayJoin(hashtags) hashtag,
        count(*) count,
        row_number() over (partition by profile_id order by count desc) rank
    from dbt.mart_instagram_post_hashtags
    group by profile_id, hashtag
    order by profile_id desc, rank asc
),
final as (
    select
        'INSTAGRAM' platform,
        ifNull(profile_id, 'MISSING') platform_id,
        groupArray(hashtag) hashtags,
        groupArray(count) hashtags_counts
    from hashtags
    where rank <= 10
    group by platform_id
)
select * from final
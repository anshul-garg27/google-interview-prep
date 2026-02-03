{{ config(materialized = 'table', tags=["weekly"], order_by='shortcode') }}
with base as (
    select
        profile_id,
        handle,
        shortcode,
        max(caption) caption
    from dbt.stg_beat_instagram_post
    group by profile_id, handle, shortcode
),
final as (
    select
        profile_id,
        handle,
        ifNull(shortcode, 'MISSING') shortcode,
        extractAll(ifNull(caption, ''), '#[a-zA-Z0-9]+') hashtags,
        arrayMap(x -> lower(x), hashtags) hashtags_l
    from
        base
)
select * from final
SETTINGS max_bytes_before_external_group_by = 20000000000, max_bytes_before_external_sort = 20000000000, join_algorithm='partial_merge'

{{ config(materialized = 'table', tags=["weekly"], order_by='profile_id') }}
with local as (
    SELECT
        profile_id,
        detectLanguage(ifNull(caption, '')) lang,
        uniqExact(shortcode) AS posts
    FROM dbt.stg_beat_instagram_post
    where caption != '' and lang IN ('bn', 'ta', 'te', 'mr', 'ka', 'te', 'hi')
    group by profile_id, lang
),
final as (
    select
        ifNull(profile_id, 'MISSING') profile_id,
        argMax(lang, posts) lang
    from local
    group by profile_id
)
select * from final
SETTINGS allow_experimental_nlp_functions = 1, max_bytes_before_external_group_by = 40000000000, max_bytes_before_external_sort = 40000000000, join_algorithm='partial_merge'

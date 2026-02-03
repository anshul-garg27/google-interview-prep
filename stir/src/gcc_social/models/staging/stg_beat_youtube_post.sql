{{ config(
            materialized = 'incremental',
            engine='ReplacingMergeTree',
            tags=["post_ranker_partial"],
            order_by='ifNull(id, 0)'
    )
}}
with
data as (
     select * from beat_replica.youtube_post_simple
    {% if is_incremental() %}
        where updated_at > (select max(updated_at) - INTERVAL 4 HOUR from {{ this }})
    {% endif %}
)
select * from data

{{ config(
            materialized = 'incremental',
            engine='ReplacingMergeTree',
            partition_by='date(ifNull(timestamp, date(\'2020-01-01\')))',
            tags=["disabled"],
            order_by='id',
            post_hook = 'OPTIMIZE TABLE dbt.stg_beat_post_log'
    )
}}
with
data as (
     select * from beat_replica.post_log
    {% if is_incremental() %}
        where timestamp > (select max(timestamp) - INTERVAL 4 HOUR from {{ this }})
    {% endif %}
)
select * from data

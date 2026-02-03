{{ config(
            materialized = 'incremental',
            engine='ReplacingMergeTree',
            partition_by='date(ifNull(created_at, date(\'2020-01-01\')))',
            tags=["hourly"],
            order_by='id',
            post_hook = 'OPTIMIZE TABLE dbt.stg_beat_asset_log'
    )
}}
with
data as (
     select * from beat_replica.asset_log
    {% if is_incremental() %}
        where created_at > (select max(created_at) - INTERVAL 4 HOUR from {{ this }})
    {% endif %}
)
select * from data

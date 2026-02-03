{{ config(
            materialized = 'incremental',
            engine='ReplacingMergeTree',
            partition_by='date(ifNull(created_at, date(\'2020-01-01\')))',
            tags=["hourly"],
            order_by='id'
    )
}}
with
data as (
     select * from beat_replica.recent_scrape_request_log
    {% if is_incremental() %}
        where
            created_at > (select max(created_at) - INTERVAL 4 HOUR from {{ this }}) or
            scraped_at > (select max(scraped_at) - INTERVAL 4 HOUR from {{ this }})
    {% endif %}
)
select * from data

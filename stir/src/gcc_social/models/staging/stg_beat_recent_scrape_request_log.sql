{{ config(materialized = 'table', tags=["recent_scl"], order_by='id') }}
with
data as (
     select ifNull(id, 0) id, source, platform, data, scraped_at, flow, idempotency_key, status, priority, params, created_at from beat_replica.recent_scrape_request_log
)
select * from data

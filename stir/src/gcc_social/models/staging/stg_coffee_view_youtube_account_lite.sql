{{ config(materialized = 'table', tags=["hourly"], order_by='ifNull(id, 0)') }}
with
data as (
     select * from coffee_replica.view_youtube_account_lite
)
select * from data

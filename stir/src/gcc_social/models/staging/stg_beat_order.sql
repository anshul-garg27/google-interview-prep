{{ config(materialized = 'table', tags=["hourly"], order_by='id') }}
with
data as (
     select * from beat_replica.order
)
select * from data

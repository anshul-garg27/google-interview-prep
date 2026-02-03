{{ config(materialized = 'table', tags=["hourly"], order_by='ifNull(id, 0)') }}
with
data as (
     select * from coffee_replica.view_campaign_profile_lite
)
select * from data

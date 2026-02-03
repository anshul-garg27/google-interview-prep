{{ config(materialized = 'table', tags=["hourly"], order_by='id') }}
with
data as (
    select * from beat_replica.youtube_profile_insights
)
select * from data

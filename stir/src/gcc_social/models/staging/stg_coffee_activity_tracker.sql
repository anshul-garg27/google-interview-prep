{{ config(materialized = 'table', tags=["collections"], order_by='id') }}
with
data as (
     select * from coffee_replica.activity_tracker
)
select * from data

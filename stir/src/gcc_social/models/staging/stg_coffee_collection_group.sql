{{ config(materialized = 'table', tags=["collections"], order_by='id') }}
with
data as (
     select * from coffee_replica.collection_group
)
select * from data

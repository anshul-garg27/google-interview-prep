{{ config(materialized = 'table', tags=["collections"], order_by='id') }}
with
data as (
     select * from coffee_replica.post_collection_item
)
select * from data

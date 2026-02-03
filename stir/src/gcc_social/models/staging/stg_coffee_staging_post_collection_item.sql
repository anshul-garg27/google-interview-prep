{{ config(materialized = 'table', tags=["staging_collections"], order_by='id') }}
with
data as (
     select * from coffee_stage.post_collection_item
)
select * from data

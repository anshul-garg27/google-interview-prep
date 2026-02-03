{{ config(materialized = 'table', tags=["disabled"], order_by='ifNull(id, 0)') }}
with
data as (
     select * from coffee_stage.view_instagram_account_lite
)
select * from data

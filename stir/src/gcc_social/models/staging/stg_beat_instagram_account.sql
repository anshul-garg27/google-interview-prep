{{ config(materialized = 'table', tags=["post_ranker_partial"], order_by='id') }}
with
data as (
     select * from beat_replica.instagram_account
)
select * from data

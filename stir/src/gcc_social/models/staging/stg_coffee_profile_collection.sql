{{ config(materialized = 'table', tags=["collections"], order_by='id') }}
with data as (
      select 
          id, 
          share_id, 
          name, 
          partner_id, 
          created_at, 
          source 
      from 
      coffee_replica.profile_collection
)
select * from data
{{ config(materialized = 'table', tags=["collections"], order_by='id') }}
with
data as (
     select 
          id,
          platform,
          platform_account_code,
          profile_collection_id,
          hidden,
          enabled,
          '' custom_columns,
          rank,
          created_by,
          created_at,
          updated_at,
          partner_id,
          profile_social_id,
          campaign_profile_id,
          shortlisting_status,
          shortlist_id
 from coffee_replica.profile_collection_item
)
select * from data

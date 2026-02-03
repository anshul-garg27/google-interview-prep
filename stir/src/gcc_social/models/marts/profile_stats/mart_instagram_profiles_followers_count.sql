{{ config(materialized = 'table', tags=["daily"], order_by='followers_in_db') }}
with follower_data as (select log.target_profile_id as                                       target_profile_id,
                              JSONExtractString(source_dimensions, 'handle')    follower_handle,
                              JSONExtractString(source_dimensions, 'full_name') follower_full_name
                       from dbt.stg_beat_profile_relationship_log log
                       where follower_handle is not null
                         and follower_full_name is not null
                         and follower_handle != ''
                         and follower_full_name != '')
select target_profile_id, uniqExact(follower_handle) followers_in_db
from follower_data
group by target_profile_id
order by followers_in_db desc
SETTINGS max_bytes_before_external_group_by = 40000000000, max_bytes_before_external_sort = 40000000000
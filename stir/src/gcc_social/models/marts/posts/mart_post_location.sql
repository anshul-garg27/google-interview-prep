{{ config(materialized = 'table', tags=["weekly"], order_by='short_code') }}
select
    source,
    ifNull(shortcode, 'MISSING') short_code,
    profile_id,
    handle,
    JSONExtractString(dimensions, 'location') location,
    JSONExtractString(dimensions, 'location_lat') location_lat,
    JSONExtractString(dimensions, 'location_long') location_long
from _e.post_log_events
where source != 'graphapi' and location != '[]' and location != ''
SETTINGS max_bytes_before_external_group_by = 20000000000, max_bytes_before_external_sort = 20000000000, join_algorithm='partial_merge'

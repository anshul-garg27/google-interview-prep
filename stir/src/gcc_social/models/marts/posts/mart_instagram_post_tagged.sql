{{ config(materialized = 'table', tags=["weekly"], order_by='shortcode') }}
select
       ifNull(shortcode, 'MISSING') shortcode,
       max(arrayMap(x ->
            JSONExtractString(JSONExtractString(x, 'user'), 'username'),
            JSONExtractArrayRaw(JSONExtractString(dimensions, 'user_tags')))) tagged
from _e.post_log_events
   where JSONExtractString(dimensions, 'user_tags') is not null and
   JSONExtractString(dimensions, 'user_tags') != '' and
   JSONExtractString(dimensions, 'user_tags') != '[]' and
   JSONExtractString(dimensions, 'user_tags') like '[%'
group by shortcode
SETTINGS max_bytes_before_external_group_by = 20000000000, max_bytes_before_external_sort = 20000000000, join_algorithm='partial_merge'

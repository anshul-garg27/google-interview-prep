{{ config(materialized = 'table', tags=["daily"], order_by='handle') }}
with
youtube_profile_relationship as (
    select
        JSONExtractString(ifNull(params, ''), 'channel_id') as handle,
        max(scraped_at) last_activities_attempt
    from
        dbt.stg_beat_scrape_request_log
    where
        flow = 'refresh_yt_profile_relationship'
    group by handle
    order by last_activities_attempt desc
)

select
    ypr.handle handle,
    ypr.last_activities_attempt last_activities_attempt
from youtube_profile_relationship ypr
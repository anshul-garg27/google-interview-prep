{{ config(materialized = 'table', tags=["account_tracker_stats"], order_by='handle') }}
with 
handles_attempted_raw as (
    select 
        JSONExtractArrayRaw(ifNull(params, ''), 'channel_ids') as handles,
        max(scraped_at) last_attempt
    from 
        dbt.stg_beat_scrape_request_log
    where 
        status IN ('FAILED', 'COMPLETE') 
        and ((data NOT LIKE '(204%'
        and data NOT LIKE '(409%'
        and data NOT LIKE '(429%')
        or (data is null))
        and flow = 'refresh_yt_profiles'
    group by handles
    order by last_attempt desc
),
handles_attempted as (
    select 
        trim(BOTH '"' FROM arrayJoin(handles)) handle,
        max(last_attempt) last_attempt
    from handles_attempted_raw
    group by handle
),
handles_videos_attempted as (
    select 
        JSONExtractString(params, 'channel_id') as handle,
        max(scraped_at) last_attempt
    from 
        dbt.stg_beat_scrape_request_log
    where 
        status IN ('FAILED', 'COMPLETE') 
        and (
                (
                    data NOT LIKE '(204%'
                    and data NOT LIKE '(409%'
                    and data NOT LIKE '(429%'
                ) 
                or 
                (data is null)
            )
        and flow = 'refresh_yt_posts_by_channel_id'
    group by handle
    order by last_attempt desc
),
saas_collection_handles as (
    select 
        ya.channel_id handle, 'saas_collection' source, 1 priority, 1 frequency
    from dbt.stg_coffee_profile_collection_item pci
    left join dbt.stg_coffee_profile_collection pc on pc.id = pci.profile_collection_id
    left join dbt.stg_coffee_view_youtube_account_lite ya on ya.id = pci.platform_account_code
    where handle != '' and platform = 'YOUTUBE' and pc.source = 'SAAS-AT'
    group by handle
),
saas_profile_collection_handles as (
    select 
        ya.channel_id handle, 'saas_profile_collection' source, 2 priority, 1 frequency
    from dbt.stg_coffee_profile_collection_item pci
    left join dbt.stg_coffee_profile_collection pc on pc.id = pci.profile_collection_id
    left join dbt.stg_coffee_view_youtube_account_lite ya on ya.id = pci.platform_account_code
    where handle != '' and platform = 'YOUTUBE' and pc.source = 'SAAS'
    group by handle
),
leaderboard_handles as (
    select 
        l.platform_id handle, 'leaderboard' source, 3 priority, 1 frequency
    from dbt.mart_youtube_leaderboard l
    where 
        month >= date('2023-02-01')
        and 
            (
            followers_rank < 1000 or 
            followers_change_rank < 1000 or 
            views_rank < 1000 or 

            followers_rank_by_cat < 1000 or 
            followers_change_rank_by_cat < 1000 or 
            views_rank_by_cat < 1000 or 

            followers_rank_by_lang < 1000 or 
            views_rank_by_lang < 1000 or 
            followers_change_rank_by_lang < 1000 or

            views_rank_by_cat_lang < 1000 or 
            followers_rank_by_cat_lang < 1000 or 
            followers_change_rank_by_cat_lang < 1000
            )
),
rest_10k as (
    select 
        bia.channel_id handle, 'rest_beat_10k' source, 4 priority, 1 frequency
    from 
        dbt.stg_beat_youtube_account bia
    where
        bia.subscribers > 10000 
    group by handle
),
rest as (
    select 
        bia.channel_id handle, 'rest_beat' source, 5 priority, 7 frequency
    from 
        dbt.stg_beat_youtube_account bia
    group by handle
),
all_handles as (
    select handle, source, priority p, frequency f from saas_collection_handles
    union all
    select handle, source, priority p, frequency f from saas_profile_collection_handles
    union all
    select handle, source, priority p, frequency f from leaderboard_handles
    union all
    select handle, source, priority p, frequency f from rest_10k
    union all
    select handle, source, priority p, frequency f from rest
),
combined as (
    select 
        ifNull(handle, 'MISSING') handle, 
        argMin(source, p) source, 
        min(p) priority,
        min(f) frequency
    from all_handles
    group by handle
)
select 
    combined.handle handle, 
    source, 
    priority, 
    frequency,
    ha.last_attempt last_attempt,
    hav.last_attempt last_videos_attempt,
    if(
        last_attempt is not null and last_attempt >= '2020-01-01', 
            last_attempt + interval frequency day, 
            now()
    ) next_attempt,
    if(
        last_videos_attempt is not null and last_videos_attempt >= '2020-01-01', 
            last_videos_attempt + interval frequency day, 
            now()
    ) next_videos_attempt
from combined
left join handles_attempted ha on ha.handle = combined.handle
left join handles_videos_attempted hav on hav.handle = combined.handle

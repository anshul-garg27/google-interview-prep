{{ config(materialized = 'table', tags=["hourly"], order_by='channel_id') }}
with base as (
    select
        ifNull(channel_id, 'MISSING') channel_id,
        multiIf(audio_language like 'en-%', 'en', audio_language like 'hi-%', 'hi', audio_language) audio_language,
        --detectLanguage(description) audio_language,
        count(*) total_videos,
        sum(total_videos) over (partition by channel_id) total,
        1.0*total_videos / total perc
    from dbt.stg_beat_youtube_post
    where audio_language is not null and audio_language != 'zxx'
    group by channel_id, audio_language
    order by channel_id asc
)
select channel_id, audio_language language from base
where perc > 0.5 and total_videos >= 5
SETTINGS allow_experimental_nlp_functions = 1
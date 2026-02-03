{{ config(
            materialized = 'incremental',
            engine='ReplacingMergeTree',
            unique_key='(subscriber_id, channel_id)',
            order_by='(subscriber_id, channel_id)',
    )
}}
with events as (
    select *
    from _stage_e.yt_profile_relationship_log_events
),
semi_final as (
    select
          source_channel_id                 as subscriber_id,
          target_channel_id                 as channel_id,
          toDateTime(subscribed_on)            timestamp,
          toDateTime(yprle.event_timestamp) as event_timestamp
   from events yprle
),
final as (
    select channel_id,
           subscriber_id,
           timestamp,
           event_timestamp,
           dictGet('default.mart_yt_dict', 'category', channel_id) channel_gcc_category,
           dictGet('default.mart_yt_dict', 'yt_category', channel_id) channel_yt_category,
           dictGet('default.mart_yt_dict', 'country', channel_id) channel_country,
           dictGet('default.mart_yt_dict', 'category', subscriber_id) subscriber_gcc_category,
           dictGet('default.mart_yt_dict', 'yt_category', subscriber_id) subscriber_yt_category,
           dictGet('default.mart_yt_dict', 'country', subscriber_id) subscriber_country
    from semi_final
)
select * from final
{% if is_incremental() %}
where event_timestamp >= now() - INTERVAL 4 DAY
{% endif %}
{{ config(materialized = 'table', tags=["weekly"], order_by='profile_id') }}
with phone_from_bio as (
    SELECT
        ifNull(profile_id, 'MISSING') profile_id,
        handle,
        extractAllGroupsHorizontal(ifNull(biography, ''), '(?i)(.*\\W|^)([1-9][0-9]{9})(\\W.*|$)')[2][1] phone
    from dbt.stg_beat_instagram_account
    where phone != ''
),
phones_from_contact as (
    select
        profile_id,
        handle,
        toString(JSONExtractString(dimensions, 'contact_phone')) phone
    from
    _e.profile_log_events
    where phone != ''
)
select * from phone_from_bio union all select * from phones_from_contact

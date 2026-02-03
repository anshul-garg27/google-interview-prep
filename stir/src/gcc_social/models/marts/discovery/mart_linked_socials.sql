{{ config(materialized = 'table', tags=["hourly"], order_by='channel_id, profile_id') }}
with
base as (
    select
        youtube_userid channel_id,
        instagram_userid profile_id
    from
    vidooly.social_profiles_new_20230217
    where youtube_userid != '' and instagram_userid != ''
),
ia_linked as (
    select profile_id from base
),
ya_linked as (
    select channel_id from base
),
ia as (
    select *
    from
    dbt.mart_instagram_account ia
    where ig_id in ia_linked and country = 'IN'
),
ya as (
    select * from
    dbt.mart_youtube_account y
    where y.channelid in ya_linked and country = 'IN'
),
final as (
    select
        base.channel_id youtube_url,
        ia.ig_id profile_id,
        ia.name name,
        ya.followers yt_followers,
        ia.followers ia_followers,
        ya.category yt_category,
        ia.primary_category_original ia_category,
        ia.primary_category ia_category_final
    from
        base
    left join ia on ia.ig_id = base.profile_id
    left join ya on ya.channelid = base.channel_id
    where ia.handle != '' and  yt_followers > 0 and ia_followers > 0
    order by ya.followers + ia.followers desc
),
semi_final as (
    select
        profile_id,
        argMax(youtube_url, yt_followers) channel_id,
        max(yt_followers) channel_followers,
        max(ia_followers) handle_followers,
        max(yt_category) channel_category,
        max(ia_category) handle_category,
        max(ia_category_final) handle_category_final
    from final
    group by profile_id
),
handle_maps as (
    select
        ig_id profile_id,
        handle
    from dbt.mart_instagram_account
),
links as (
    select
        ifNull(channel_id, 'MISSING') channel_id,
        ifNull(argMax(profile_id, handle_followers), 'MISSING') profile_id,
        max(channel_followers) ya_followers,
        max(channel_category) ya_category,
        max(handle_followers) ia_followers,
        max(handle_category) ia_category,
        max(handle_category_final) ia_category_final,
        ya_followers + ia_followers total_followers
    from semi_final
    group by channel_id
    order by total_followers desc
)
select
    links.channel_id channel_id,
    links.profile_id profile_id,
    hm.handle handle,
    ya_category,
    ia_category,
    ia_category_final,
    coalesce(if(ya_category = '', NULL, ya_category), ia_category) category,
    ya_followers,
    ia_followers,
    total_followers
from links
left join handle_maps hm on hm.profile_id = links.profile_id

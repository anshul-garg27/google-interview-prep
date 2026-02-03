{{ config(materialized = 'table', tags=["hourly"], order_by='platform_id') }}
with
handles_and_targets as (
    select
        JSONExtractString(source_dimensions, 'handle') handle,
        target_profile_id
    from
        dbt.stg_beat_profile_relationship_log pl
        where target_profile_id IN ('2094200507', '8446808', '2254659037', '3229110093', '526577491','3986772984','47117696314','2198834499','13820762279','3861104115',
        '2094200507', '8446808', '2254659037', '3229110093', '526577491',
        '3861104115','3160696766','3423799830','7882454627','4019288322','1461335285','649029116','25818229759','4954426873','44675201741','2136138769','1552440413','3986772984','47117696314','2198834499','13820762279','3861104115')
),
handles as (
    select handle
    from handles_and_targets
    group by handle
),
targets as (
    select target_profile_id
    from handles_and_targets
    group by target_profile_id
),
target_handles as (
    select
        ig_id profile_id,
        handle,
        100*engagement_rate engagement_rate,
        likes_to_comment_ratio
    from
    dbt.mart_instagram_account
    where profile_id in targets
),
base as (
    select
        target_profile_id,
        ia.profile_id profile_id,
        iab.handle target_handle,
        h.handle handle,
        followers,
        following,
        if(following > 0, followers/following, 1.0*followers) ffratio,
        media_count,
        biography
    from
    dbt.stg_beat_instagram_account ia
    left join handles_and_targets h on h.handle = ia.handle
    left join target_handles iab on iab.profile_id = target_profile_id
    where ia.handle in handles
),
follower_ids as (
    select profile_id from base group by profile_id
),
follower_ia as (
    select * from dbt.mart_instagram_account where ig_id in follower_ids
),
follower_details as (
    select
        *,
        if(
            gender_male > gender_female and gender_male > 0.7, 'M',
            if(gender_female > 0.7, 'F', NULL)) gender_from_name,
        coalesce(if(gender_from_img_model = '', NULL, gender_from_img_model),
            gender_from_name) gender,
        multiIf(age >= 13 and age < 18, '13-17',
                age >= 18 and age < 25, '18-24',
                age >= 25 and age < 35, '25-34',
                age >= 35 and age < 45, '35-44',
                age >= 45 and age < 55, '45-54',
                age >= 55 and age < 65, '55-64',
                age >= 65, '65+',
                'NA'
                ) age_bucket,
        concat(gender, '_', age_bucket) age_group

    from
    ds.instagram_prediction
),
followers as (
    select * from base
    left join follower_details fd on toString(fd.profile_id) = base.profile_id
),
age_groups as (
    select target_profile_id,
            age_group,
            uniqExact(handle) handles,
            sum(handles) over (partition by target_profile_id) total, handles/total perc
    from followers
    where age_group != 'NA' and age_group != 'F_NA' and age_group != 'M_NA'
    group by target_profile_id, age_group
    order by target_profile_id desc, perc desc, age_group desc
),
age_group_summary as (
    select
        target_profile_id,
        max(handles) age_group_handles,
        groupArray(age_group) groups,
        groupArray(round(100*perc,2)) group_percs
    from age_groups
    group by target_profile_id
),
cities as (
    select target_profile_id,
            top_cities city,
            if(uniqExact(handle) > 0, uniqExact(handle)*100 + round(randUniform(10,100)), 0) handles,
            sum(handles) over (partition by target_profile_id) total,
            handles/total perc
    from followers
    where city != ''
    group by target_profile_id, city
    order by target_profile_id desc, perc desc, city desc
),
city_summary as (
    select
        target_profile_id,
        sum(handles) city_handles,
        groupArray(city) city_list,
        groupArray(round(100*perc,2)) city_percs
    from cities
    group by target_profile_id
),
states as (
    select target_profile_id,
            state,
            uniqExact(handle) handles,
            sum(handles) over (partition by target_profile_id) total, handles/total perc
    from followers
    where state != ''
    group by target_profile_id, state
    order by target_profile_id desc, perc desc, state desc
),
state_summary as (
    select
        target_profile_id,
        sum(handles) state_handles,
        groupArray(state) state_list,
        groupArray(round(100*perc,2)) state_percs
    from states
    group by target_profile_id
),
countries as (
    select target_profile_id,
            country,
            uniqExact(handle) handles,
            sum(handles) over (partition by target_profile_id) total, handles/total perc
    from followers
    where country != ''
    group by target_profile_id, country
    order by target_profile_id desc, perc desc, country desc
),
countries_summary as (
    select
        target_profile_id,
        sum(handles) country_handles,
        groupArray(country) country_list,
        groupArray(round(100*perc,2)) country_percs
    from countries
    group by target_profile_id
),
notable_followers_base as (
    select
        target_profile_id, target_handle, ia.ig_id platform_id, ia.handle handle, ia.followers followers
    from base
    left join follower_ia ia on ia.ig_id = base.profile_id
    where ia.ig_id != ''
),
notable_followers as (
    select
        target_profile_id,
        groupArrayIf(platform_id, followers > 100000) notable_followers
    from notable_followers_base
    group by target_profile_id
    having length(notable_followers) > 0
),
final as (
    select
        f.target_profile_id target_profile_id,
        target_handle,
        uniqExact(handle) total_followers,
        uniqExactIf(handle, following < 1500) / uniqExact(handle) reachable_followers,
        uniqExactIf(handle, ffratio < 0.1) / uniqExact(handle) ff_ratio_lt_01,
        uniqExactIf(handle, following > 5000) / uniqExact(handle) high_following,
        uniqExactIf(handle, followers < 10) / uniqExact(handle) low_followers,
        uniqExactIf(handle, media_count = 0) / uniqExact(handle) zero_post_followers,
        uniqExactIf(handle, gender = 'M' or gender = 'F') followers_with_gender,
        uniqExactIf(handle, gender = 'M') / followers_with_gender male_followers_perc,
        uniqExactIf(handle, gender = 'F') / followers_with_gender female_followers_perc
    from followers f
    group by target_profile_id, target_handle
),
data as (
    select
        'INSTAGRAM' platform,
        ifNull(final.target_profile_id, 'MISSING') platform_id,
        final.target_handle platform_handle,
        round(100*final.reachable_followers, 2) reachable_followers_perc,
        round(100*(1 - ff_ratio_lt_01),2) audience_authenticity,
        round(100*final.low_followers,2) breakup_low_followers,
        round(100*final.high_following) breakup_high_following,
        round(100*1 - final.ff_ratio_lt_01) breakup_audience_authenticity,
        multiIf(audience_authenticity >= 90, 'Excellent',
                audience_authenticity >= 75, 'Very Good',
                audience_authenticity >= 60, 'Good',
                audience_authenticity >= 40, 'Average',
                'Poor'
            ) breakup_audience_authenticity_grade,
        multiIf((100-breakup_high_following) >= 90, 'Excellent',
                (100-breakup_high_following) >= 75, 'Very Good',
                (100-breakup_high_following) >= 60, 'Good',
                (100-breakup_high_following) >= 40, 'Average',
                'Poor'
            ) breakup_high_following_grade,
        multiIf((100 - breakup_low_followers) >= 90, 'Excellent',
                (100 - breakup_low_followers) >= 75, 'Very Good',
                (100 - breakup_low_followers) >= 60, 'Good',
                (100 - breakup_low_followers) >= 40, 'Average',
                'Poor'
            ) breakup_low_followers_grade,
        'Poor' breakup_followers_growth30d_grade,
        'Poor' breakup_views_30d_grade,
        'Poor' breakup_uploads_30d_grade,
        100*(1-(0.25 * final.low_followers) - (0.25*final.high_following) - (0.5* ff_ratio_lt_01)) audience_quality_score,
        multiIf(audience_quality_score >= 90, 'Excellent',
                audience_quality_score >= 75, 'Very Good',
                audience_quality_score >= 60, 'Good',
                audience_quality_score >= 40, 'Average',
                'Poor'
            ) audience_quality_grade,
        round(100*final.male_followers_perc,2) male_followers_perc,
        round(100*final.female_followers_perc,2) female_followers_perc,
        nf.notable_followers notable_followers,
        ags.groups age_groups,
        ags.group_percs age_group_percs,
        cs.city_list cities,
        cs.city_percs city_percs,
        cns.country_list countries,
        cns.country_percs country_percs
    from final
    left join age_group_summary ags on ags.target_profile_id = final.target_profile_id
    left join city_summary cs on cs.target_profile_id = final.target_profile_id
    left join state_summary ss on ss.target_profile_id = final.target_profile_id
    left join countries_summary cns on cns.target_profile_id = final.target_profile_id
    left join notable_followers nf on nf.target_profile_id = final.target_profile_id
),
mart_chans as (
    select channelid from dbt.mart_youtube_account where country = 'IN'
),
yt_info as (
    select
        'YOUTUBE' platform,
        channelid platform_id,
        channelid platform_handle,
        0 reachable_followers_perc,
        0 audience_authenticity,
        0 breakup_low_followers,
        0 breakup_high_following,
        0 breakup_audience_authenticity,
        'Poor' breakup_audience_authenticity_grade,
        'Poor' breakup_high_following_grade,
        'Poor' breakup_low_followers_grade,
        followers_growth30d_grade breakup_followers_growth30d_grade,
        views_30d_grade breakup_views_30d_grade,
        uploads_30d_grade breakup_uploads_30d_grade,
        channel_quality_score audience_quality_score,
        channel_quality_grade audience_quality_grade,
        youtube_demographic_m_all male_followers_perc,
        youtube_demographic_f_all female_followers_perc,
        [] notable_followers,
        ['M_13-17', 'M_18-24', 'M_25-34', 'M_35-44', 'M_45-54', 'M_55-64', 'M_65+', 'F_13-17', 'F_18-24', 'F_25-34', 'F_35-44', 'F_45-54', 'F_55-64', 'F_65+'] age_groups,
        [
         ifNull(youtube_demographic_m_13_17, 0),
         ifNull(youtube_demographic_m_18_24, 0),
         ifNull(youtube_demographic_m_25_34, 0),
         ifNull(youtube_demographic_m_35_44, 0),
         ifNull(youtube_demographic_m_45_54, 0),
         ifNull(youtube_demographic_m_55_64, 0),
         ifNull(youtube_demographic_m_65, 0),
         ifNull(youtube_demographic_f_13_17, 0),
         ifNull(youtube_demographic_f_18_24, 0),
         ifNull(youtube_demographic_f_25_34, 0),
         ifNull(youtube_demographic_f_35_44, 0),
         ifNull(youtube_demographic_f_45_54, 0),
         ifNull(youtube_demographic_f_55_64, 0),
         ifNull(youtube_demographic_f_65, 0)
        ] age_groups_percs,
        [] cities,
        [] city_percs,
        ['IN'] countries,
        [100.0] country_percs
    from dbt.mart_youtube_account ya
    left join vidooly.social_profiles_new_20230217 sp on sp.youtube_userid = ya.channelid
    where ya.country = 'IN'
)
select * from data union all select * from yt_info

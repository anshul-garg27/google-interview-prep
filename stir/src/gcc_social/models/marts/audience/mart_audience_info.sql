{{ config(materialized = 'table', tags=["hourly"], order_by='platform_id') }}
with
private_user as (
    select
        maip.platform platform,
        maip.platform_id platform_id,
        maip.platform_handle platform_handle,
        maip.reachable_followers_perc reachable_followers_perc,
        maip.audience_authenticity audience_authenticity,
        maip.breakup_low_followers breakup_low_followers,
        maip.breakup_high_following breakup_high_following,
        maip.breakup_audience_authenticity breakup_audience_authenticity,
        maip.breakup_audience_authenticity_grade breakup_audience_authenticity_grade,
        maip.breakup_high_following_grade breakup_high_following_grade,
        maip.breakup_low_followers_grade breakup_low_followers_grade,
        maip.breakup_followers_growth30d_grade breakup_followers_growth30d_grade,
        maip.breakup_views_30d_grade breakup_views_30d_grade,
        maip.breakup_uploads_30d_grade breakup_uploads_30d_grade,
        maip.audience_quality_score audience_quality_score,
        maip.audience_quality_grade audience_quality_grade,
        maip.male_followers_perc male_followers_perc,
        maip.female_followers_perc female_followers_perc,
        maip.notable_followers notable_followers,
        maip.age_groups age_groups,
        maip.age_group_percs age_group_percs,
        maip.cities cities,
        maip.city_percs city_percs,
        maip.countries countries,
        maip.country_percs country_percs
    from dbt.mart_audience_info_private maip
),
follower_users_not_in_private as (
    select
        maif.platform platform,
        maif.platform_id platform_id,
        maif.platform_handle platform_handle,
        maif.reachable_followers_perc reachable_followers_perc,
        maif.audience_authenticity audience_authenticity,
        maif.breakup_low_followers breakup_low_followers,
        maif.breakup_high_following breakup_high_following,
        maif.breakup_audience_authenticity breakup_audience_authenticity,
        maif.breakup_audience_authenticity_grade breakup_audience_authenticity_grade,
        maif.breakup_high_following_grade breakup_high_following_grade,
        maif.breakup_low_followers_grade breakup_low_followers_grade,
        maif.breakup_followers_growth30d_grade breakup_followers_growth30d_grade,
        maif.breakup_views_30d_grade breakup_views_30d_grade,
        maif.breakup_uploads_30d_grade breakup_uploads_30d_grade,
        maif.audience_quality_score audience_quality_score,
        maif.audience_quality_grade audience_quality_grade,
        maif.male_followers_perc male_followers_perc,
        maif.female_followers_perc female_followers_perc,
        maif.notable_followers notable_followers,
        maif.age_groups age_groups,
        maif.age_group_percs age_group_percs,
        maif.cities cities,
        maif.city_percs city_percs,
        maif.countries countries,
        maif.country_percs country_percs
    from dbt.mart_audience_info_follower maif
    left join private_user private on private.platform_handle = maif.platform_handle
    where private.platform_handle is NULL
),
gpt_users_not_in_private_and_follower as (
    select
        maig.platform as platform,
        maig.platform_id platform_id,
        maig.platform_handle platform_handle,
        maig.reachable_followers_perc reachable_followers_perc,
        maig.audience_authenticity audience_authenticity,
        maig.breakup_low_followers breakup_low_followers,
        maig.breakup_high_following breakup_high_following,
        maig.breakup_audience_authenticity breakup_audience_authenticity,
        maig.breakup_audience_authenticity_grade breakup_audience_authenticity_grade,
        maig.breakup_high_following_grade breakup_high_following_grade,
        maig.breakup_low_followers_grade breakup_low_followers_grade,
        maig.breakup_followers_growth30d_grade breakup_followers_growth30d_grade,
        maig.breakup_views_30d_grade breakup_views_30d_grade,
        maig.breakup_uploads_30d_grade breakup_uploads_30d_grade,
        maig.audience_quality_score audience_quality_score,
        maig.audience_quality_grade audience_quality_grade,
        maig.male_followers_perc male_followers_perc,
        maig.female_followers_perc female_followers_perc,
        maig.notable_followers notable_followers,
        maig.age_groups age_groups,
        maig.age_group_percs age_group_percs,
        maig.cities cities,
        maig.city_percs city_percs,
        maig.countries countries,
        maig.country_percs country_percs
    from dbt.mart_audience_info_gpt maig
    left join private_user private on maig.platform_handle = private.platform_handle
    left join follower_users_not_in_private follower on follower.platform_handle = maig.platform_handle
    where private.platform_handle IS NULL AND follower.platform_handle = ''
)
select * from private_user 
         union all
select * from follower_users_not_in_private
         union all
select * from gpt_users_not_in_private_and_follower
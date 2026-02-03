{{ config(materialized = 'table', tags=["hourly"], order_by='platform_id') }}
with
unique_insagram_age_gender as (
    select
        profile_id,
        argMax(dimensions, event_timestamp) dimension
    from _e.profile_log_events
    where source = 'openapi_profile_info_v0.11' and JSONHas(dimensions, 'audience_gender_age')
    group by profile_id
),
unique_insagram_cities as (
    select
        profile_id,
        argMax(dimensions, event_timestamp) dimension
    from _e.profile_log_events
    where source = 'openapi_profile_info_v0.11' and JSONHas(dimensions, 'audience_city')
    group by profile_id
),
instagram_age_gender_groups as (
    select
        uiag.profile_id as handle,
        mia.ig_id as profile_id,
        JSONExtractString(uiag.dimension, 'audience_gender_age') audience_gender_age,
        JSONExtractFloat(audience_gender_age, '13-17 female') f_13_17,
        JSONExtractFloat(audience_gender_age, '13-17 male') m_13_17,
        JSONExtractFloat(audience_gender_age, '18-24 female') f_18_24,
        JSONExtractFloat(audience_gender_age, '18-24 male') m_18_24,
        JSONExtractFloat(audience_gender_age, '25-34 female') f_25_34,
        JSONExtractFloat(audience_gender_age, '25-34 male') m_25_34,
        JSONExtractFloat(audience_gender_age, '35-44 female') f_35_44,
        JSONExtractFloat(audience_gender_age, '35-44 male') m_35_44,
        JSONExtractFloat(audience_gender_age, '45-54 female') f_45_54,
        JSONExtractFloat(audience_gender_age, '45-54 male') m_45_54,
        JSONExtractFloat(audience_gender_age, '55-64 female') f_55_64,
        JSONExtractFloat(audience_gender_age, '55-64 male') m_55_64,
        JSONExtractFloat(audience_gender_age, '65+ female') f_65,
        JSONExtractFloat(audience_gender_age, '65+ male') m_65,
        m_13_17 + m_18_24 + m_25_34 + m_35_44 + m_45_54 + m_55_64 + m_65 male_total,
        f_13_17 + f_18_24 + f_25_34 + f_35_44 + f_45_54 + f_55_64 + f_65 female_total,
        male_total + female_total total,
        round((male_total / total) * 100, 3) male_followers_perc,
        round((female_total / total) * 100, 3) female_followers_perc,
        round((m_13_17 / total) * 100, 3) m_13_17_perc,
        round((m_18_24 / total) * 100, 3) m_18_24_perc,
        round((m_25_34 / total) * 100, 3) m_25_34_perc,
        round((m_35_44 / total) * 100, 3) m_35_44_perc,
        round((m_45_54 / total) * 100, 3) m_45_54_perc,
        round((m_55_64 / total) * 100, 3) m_55_64_perc,
        round((m_65 / total) * 100, 3) m_65_perc,
        round((f_13_17 / total) * 100, 3) f_13_17_perc,
        round((f_18_24 / total) * 100, 3) f_18_24_perc,
        round((f_25_34 / total) * 100, 3) f_25_34_perc,
        round((f_35_44 / total) * 100, 3) f_35_44_perc,
        round((f_45_54 / total) * 100, 3) f_45_54_perc,
        round((f_55_64 / total) * 100, 3) f_55_64_perc,
        round((f_65 / total) * 100, 3) f_65_perc
    from unique_insagram_age_gender uiag
    left join dbt.mart_instagram_account mia on uiag.profile_id = mia.handle
),
instagram_age_gender_data as (
    select
        handle,
        profile_id,
        if(isNaN(male_followers_perc), 0, male_followers_perc) AS male_followers_perc,
        if(isNaN(female_followers_perc), 0, female_followers_perc) AS female_followers_perc,
        ['M_13-17', 'M_18-24', 'M_25-34', 'M_35-44', 'M_45-54', 'M_55-64', 'M_65+',
        'F_13-17', 'F_18-24', 'F_25-34', 'F_35-44', 'F_45-54', 'F_55-64', 'F_65+'] age_groups,
        arrayMap(x -> if(isNaN(ifNull(x, 0)), 0, ifNull(x, 0)),
        [m_13_17_perc, m_18_24_perc, m_25_34_perc, m_35_44_perc, m_45_54_perc, m_55_64_perc, m_65_perc,
        f_13_17_perc, f_18_24_perc, f_25_34_perc, f_35_44_perc, f_45_54_perc, f_55_64_perc, f_65_perc]) AS age_groups_percs
    from instagram_age_gender_groups
),
instagram_city_data AS (
    SELECT
        uic.profile_id as handle,
        mia.ig_id as profile_id,
        JSONExtractString(dimension, 'audience_city') AS ac,
        JSONLength(dimension, 'audience_city') as total_cities,
        JSONExtractKeysAndValues(ac, 'Float64') AS keys_values,
        arrayStringConcat(keys_values, ',') as joined,
        if (total_cities<5, [] ,arrayMap(x -> lowerUTF8(x.1),
        keys_values)) AS city_list,
        arrayMap(x -> x.2, keys_values) AS values_array,
        arraySum(arrayMap(x -> x.2, keys_values)) AS total_value_sum,
        if (total_cities<5, [], arrayMap(x -> round((x / total_value_sum) * 100, 2), arrayMap(x -> x.2, keys_values))) AS city_percs
    FROM unique_insagram_cities uic
    left join dbt.mart_instagram_account mia on uic.profile_id = mia.handle
),
instagram_data as (
    select
        'INSTAGRAM'   platform,
        ifNull(iagd.profile_id, 'MISSING') platform_id,
        iagd.handle   platform_handle,
        0     reachable_followers_perc,
        0     audience_authenticity,
        0     breakup_low_followers,
        0     breakup_high_following,
        0     breakup_audience_authenticity,
        ''    breakup_audience_authenticity_grade,
        ''    breakup_high_following_grade,
        ''    breakup_low_followers_grade,
        ''    breakup_followers_growth30d_grade,
        ''    breakup_views_30d_grade,
        ''    breakup_uploads_30d_grade,
        0     audience_quality_score,
        ''    audience_quality_grade,
        iagd.male_followers_perc   as male_followers_perc,
        iagd.female_followers_perc as female_followers_perc,
        ['0'] notable_followers,
        iagd.age_groups as age_groups,
        iagd.age_groups_percs age_group_percs,
        icd.city_list cities,
        icd.city_percs  as city_percs,
        ['IN'] countries,
        [1.0] country_percs
    from instagram_age_gender_data iagd
    inner join instagram_city_data icd on icd.profile_id = iagd.profile_id
)
select * from instagram_data
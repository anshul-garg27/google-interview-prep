{{ config(materialized = 'table', tags=["hourly"], order_by='group_key') }}
{%
    set metrics = [
        "avg_likes",
        "avg_comments",
        "comments_rate",
        "avg_posts_per_week",
        "followers",
        "engagement_rate",
        "video_reach",
        "uploads_30d",
        "views_30d",
        "reels_reach",
        "image_reach",
        "story_reach",
        "shorts_reach",
        "video_reach",
        "likes_to_comment_ratio",
        "followers_growth7d",
        "followers_growth30d",
        "followers_growth90d",
        "followers_growth1y",
        "audience_reachability",
        "audience_authencity",
        "audience_quality",
        "post_count",
        "likes_spread",
        "reels_impressions",
        "image_impressions"
    ]
%}
with primary_metrics_ia_base as (
    select
        primary_group,
        quantileIf(0.95)(engagement_rate, engagement_rate > 0)  p95er,
        max(engagement_rate) max_engagement_rate
    from dbt.mart_instagram_account
    where engagement_rate >= 0 and engagement_rate <= 1000
    group by primary_group
),
secondary_metrics_ia_base as (
    select
        secondary_group,
        quantileIf(0.95)(engagement_rate, engagement_rate > 0)  p95er,
        max(engagement_rate) max_engagement_rate
    from dbt.mart_instagram_account
    where engagement_rate >= 0 and engagement_rate <= 1000
    group by secondary_group
),
primary_metrics_ia as (
    select
        ifNull(primary_group, 'MISSING') group_key,
        {% for metric in metrics %}
            quantileIf(0.40)({{metric}}, {{metric}} > 0) p40_{{metric}},
            quantileIf(0.60)({{metric}}, {{metric}} > 0) p60_{{metric}},
            quantileIf(0.75)({{metric}}, {{metric}} > 0) p75_{{metric}},
            quantileIf(0.90)({{metric}}, {{metric}} > 0) p90_{{metric}},
            avgIf({{metric}}, {{metric}} > 0 and {{metric}} < 10000000000) group_avg_{{metric}},
        {% endfor %}
        uniqExact(ig_id) profiles,
        histogramIf(7)(engagement_rate, engagement_rate > 0 and engagement_rate < p95er) er_histogram,
        0 group_avg_reaction_rate,
        arrayConcat(arrayMap(x -> tupleElement(x,1), er_histogram), [max(p95er)]) bin_start,
        arrayConcat(arrayMap(x -> tupleElement(x,2), er_histogram), [max(max_engagement_rate)]) bin_end,
        arrayConcat(
            arrayMap(x -> tupleElement(x,3), er_histogram),
            [arraySum(arrayMap(x -> tupleElement(x,3), er_histogram))/9.5]) bin_height
    from
         {{ ref('mart_instagram_account') }} a
    left join primary_metrics_ia_base b on b.primary_group = a.primary_group
    group by group_key
    having profiles > 100
    order by profiles desc
),
secondary_metrics_ia as (
     select
        ifNull(secondary_group, 'MISSING') group_key,
        {% for metric in metrics %}
            quantileIf(0.40)({{metric}}, {{metric}} > 0) p40_{{metric}},
            quantileIf(0.60)({{metric}}, {{metric}} > 0) p60_{{metric}},
            quantileIf(0.75)({{metric}}, {{metric}} > 0) p75_{{metric}},
            quantileIf(0.90)({{metric}}, {{metric}} > 0) p90_{{metric}},
            avgIf({{metric}}, {{metric}} > 0 and {{metric}} < 10000000000) group_avg_{{metric}},
        {% endfor %}
        uniqExact(ig_id) profiles,
        histogramIf(7)(engagement_rate, engagement_rate > 0 and engagement_rate < p95er) er_histogram,
        0 group_avg_reaction_rate,
        arrayConcat(arrayMap(x -> tupleElement(x,1), er_histogram), [max(p95er)]) bin_start,
        arrayConcat(arrayMap(x -> tupleElement(x,2), er_histogram), [max(max_engagement_rate)]) bin_end,
        arrayConcat(
            arrayMap(x -> tupleElement(x,3), er_histogram),
            [arraySum(arrayMap(x -> tupleElement(x,3), er_histogram))/9.5]) bin_height
    from
         {{ ref('mart_instagram_account') }} a
    left join secondary_metrics_ia_base b on b.secondary_group = a.secondary_group
    group by group_key
    having profiles > 100
    order by profiles desc
),
primary_metrics_ya as (
    select
        ifNull(primary_group, 'MISSING') group_key,
        {% for metric in metrics %}
            quantileIf(0.40)({{metric}}, {{metric}} > 0) p40_{{metric}},
            quantileIf(0.60)({{metric}}, {{metric}} > 0) p60_{{metric}},
            quantileIf(0.75)({{metric}}, {{metric}} > 0) p75_{{metric}},
            quantileIf(0.90)({{metric}}, {{metric}} > 0) p90_{{metric}},
            avgIf({{metric}}, {{metric}} > 0 and {{metric}} < 10000000000) group_avg_{{metric}},
        {% endfor %}
        uniqExact(channelid) profiles,
        histogramIf(7)(engagement_rate, engagement_rate > 0 and engagement_rate < 1000) er_histogram,
        0 group_avg_reaction_rate,
        arrayMap(x -> tupleElement(x,1), er_histogram) bin_start,
        arrayMap(x -> tupleElement(x,2), er_histogram) bin_end,
        arrayMap(x -> tupleElement(x,3), er_histogram) bin_height
    from
         {{ ref('mart_youtube_account') }}
    group by group_key
    having profiles > 100
    order by profiles desc
),
secondary_metrics_ya as (
     select
        ifNull(secondary_group, 'MISSING') group_key,
        {% for metric in metrics %}
            quantileIf(0.40)({{metric}}, {{metric}} > 0) p40_{{metric}},
            quantileIf(0.60)({{metric}}, {{metric}} > 0) p60_{{metric}},
            quantileIf(0.75)({{metric}}, {{metric}} > 0) p75_{{metric}},
            quantileIf(0.90)({{metric}}, {{metric}} > 0) p90_{{metric}},
            avgIf({{metric}}, {{metric}} > 0 and {{metric}} < 10000000000) group_avg_{{metric}},
        {% endfor %}
        uniqExact(channelid) profiles,
        histogramIf(7)(engagement_rate, engagement_rate > 0 and engagement_rate < 1000) er_histogram,
        0 group_avg_reaction_rate,
        arrayMap(x -> tupleElement(x,1), er_histogram) bin_start,
        arrayMap(x -> tupleElement(x,2), er_histogram) bin_end,
        arrayMap(x -> tupleElement(x,3), er_histogram) bin_height
    from
        {{ ref('mart_youtube_account') }}
    group by group_key
    order by profiles desc
)
select * from primary_metrics_ia
         union all
select * from secondary_metrics_ia
        union all
select * from primary_metrics_ya
        union all
select * from secondary_metrics_ya
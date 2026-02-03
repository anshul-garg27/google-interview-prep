{{ config(materialized = 'table', tags=["staging_collections_fake_events"], order_by='platform, post_short_code, stats_date') }}
with
partial_time_series as (
    select * from dbt.mart_staging_collection_post_ts
),
posts as (
    select * from dbt.mart_staging_collection_post
),
collection_minmax as (
    select
        collection_id,
        toDate(min(stats_date)) start_date,
        toDate(if(max(stats_date) > toDate(now()), toDate(now()), max(stats_date))) end_date
    from partial_time_series
    group by collection_id
),
minmax as (
    select
            min(stats_date) start_date,
            if(max(stats_date) > toDate(now()), toDate(now()), max(stats_date)) end_date
    from partial_time_series
),
tseries as (
    SELECT arrayJoin(arrayMap(x -> toDate(x), range(toUInt32(start_date), toUInt32(end_date + INTERVAL 1 DAY), 3600*24))) as stats_date
    from minmax
),
collection_ts as (
    select *
    from collection_minmax cm,  tseries ts
    where ts.stats_date >= cm.start_date and ts.stats_date <= cm.end_date
),
collection_post_ts as (
    select
        collection_id,
        stats_date,
        post_short_code
    from collection_ts cts
    left join posts cp on cp.collection_id = cts.collection_id
    order by collection_id desc, post_short_code desc, stats_date asc
),
with_gaps as (
    select
            *,

            max(views) OVER (PARTITION BY collection_id, post_short_code ORDER BY stats_date ASC) _views,
            coalesce(if(views = 0, NULL, views), _views) m_views,

            max(plays) OVER (PARTITION BY collection_id, post_short_code ORDER BY stats_date ASC) _plays,
            coalesce(if(plays = 0, NULL, plays), _plays) m_plays,

            max(likes) OVER (PARTITION BY collection_id, post_short_code ORDER BY stats_date ASC) _likes,
            coalesce(if(likes = 0, NULL, likes), _likes) m_likes,

            max(comments) OVER (PARTITION BY collection_id, post_short_code ORDER BY stats_date ASC) _comments,
            coalesce(if(comments = 0, NULL, comments), _comments) m_comments,

            minIf(real_reach, real_reach > 0) OVER (PARTITION BY collection_id, post_short_code) first_reach,

            max(real_reach) OVER (PARTITION BY collection_id, post_short_code ORDER BY stats_date ASC) __real_reach,

            if(plays > 0 and __real_reach = 0, first_reach, __real_reach) _real_reach,

            coalesce(if(real_reach = 0, NULL, real_reach), _real_reach) m_real_reach,

            max(real_impressions) OVER (PARTITION BY collection_id, post_short_code ORDER BY stats_date ASC) _real_impressions,
            coalesce(if(real_impressions = 0, NULL, real_impressions), _real_impressions) m_real_impressions,

            max(platform) OVER (PARTITION BY collection_id, post_short_code ORDER BY stats_date ASC) _platform,
            coalesce(if(platform = '', NULL, platform), _platform) d_platform

    from collection_post_ts cpts
    left join partial_time_series on partial_time_series.collection_id = cpts.collection_id and partial_time_series.stats_date = cpts.stats_date and partial_time_series.post_short_code = cpts.post_short_code
),
gaps as (
    select
        d_platform _platform,
        post_short_code,
        stats_date,
        m_views views,
        m_likes likes ,
        m_comments comments,
        m_plays plays,
        m_real_impressions real_impressions,
        m_real_reach real_reach
    from with_gaps
    where
    (_platform != '' and platform = '') or
    (m_views > with_gaps.views) or
    (m_likes > with_gaps.likes) or
    (m_comments > with_gaps.comments) or
    (m_real_reach > with_gaps.real_reach) or
    (m_real_impressions > with_gaps.m_real_impressions) or
    (m_plays > with_gaps.plays)
),
final as (
    select
        ifNull(_platform, 'MISSING') platform,
        ifNull(post_short_code, 'MISSING') post_short_code,
        stats_date,
        views,
        likes,
        comments,
        plays,
        real_reach,
        real_impressions
    from gaps
)
select * from final
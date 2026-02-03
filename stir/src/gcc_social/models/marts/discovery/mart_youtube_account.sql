{{ config(materialized = 'table', tags=["hourly"], order_by='channelid') }}
{%
    set metrics = [
        "avg_likes",
        "avg_comments",
        "comments_rate",
        "followers",
        "engagement_rate",
        "likes_to_comment_ratio",
        "followers_growth7d",
        "followers_growth30d",
        "uploads_30d",
        "views_30d",
        "followers_growth90d",
        "followers_growth1y",
        "audience_reachability",
        "audience_authencity",
        "audience_quality",
        "post_count",
        "likes_spread",
        "avg_posts_per_week",
        "shorts_reach",
        "video_reach"
    ]
%}
with
blacklisted as (
    select 'UCExnfw93m4mkllxFvT326wg' as channelid
    union all
    select 'UCTSUSqcwLzTZueXIxRFCbZQ' as channelid
    union all
    select 'UCyFvRMas4fsQFq07aS8PPxA' as channelid
),
base_vidooly as (
    select
        id channel_id,
        max(parseDateTimeBestEffortOrZero(last_updated_at)) _updated_at,
        argMax(country, last_updated_at) country,
        argMax(custom_url, last_updated_at) custom_url,
        argMax(description, last_updated_at) description,
        argMax(title, last_updated_at) title,
        argMax(thumbnail, last_updated_at) thumbnail,
        argMax(toInt64OrZero(video_count), last_updated_at) uploads,
        argMax(toInt64OrZero(view_count), last_updated_at) views,
        argMax(toInt64OrZero(subscriber_count), last_updated_at) subscribers
    from
        vidooly.youtube_channels
    where id not in blacklisted
        group by channel_id
),
base_gcc as (
    select
        channel_id,
        max(updated_at) _updated_at,
        argMax(coalesce(
            if (tagged_country = '', NULL, tagged_country),
            if (country = 'India', 'IN', country)
        ), updated_at) country,
        argMax(custom_url, updated_at) custom_url,
        argMax(description, updated_at) description,
        argMax(title, updated_at) title,
        argMax(thumbnail, updated_at) thumbnail,
        argMax(uploads, updated_at) uploads,
        argMax(views, updated_at) views,
        argMax(subscribers, updated_at) subscribers
    from dbt.stg_beat_youtube_account
    where channel_id not in blacklisted
    group by channel_id
),
base_all as (
    select * from base_vidooly union all select * from base_gcc
),
collection_profiles as (
    select
        profile_social_id channel_id
    from dbt.mart_collection_post
    where platform = 'YOUTUBE'
    group by profile_social_id
),
base as (
    select
        channel_id,
        max(_updated_at) updated_at,
        argMax(country, _updated_at) country,
        argMax(custom_url, _updated_at) custom_url,
        argMax(description, _updated_at) description,
        argMax(title, _updated_at) title,
        argMax(thumbnail, _updated_at) thumbnail,
        argMax(uploads, _updated_at) video_count,
        argMax(views, _updated_at) view_count,
        argMax(subscribers, _updated_at) subscriber_count
    from base_all
    group by channel_id
),
links as (
    select
        channel_id, max(handle) handle
    from
        dbt.mart_linked_socials
    group by channel_id
),
semi_final as (
    select
        ifNull(base.channel_id, 'MISSING') channelid,
        channelid channel_id,
        base.custom_url username,
        coalesce(
            if (bya.tagged_country = '', NULL, bya.tagged_country),
            base.country
        ) country,
        base.updated_at updated_at,
        links.handle linked_instagram_handle,
        multiIf(
            yc.category = 'Film & Animation', 'Films',
            yc.category = 'Autos & Vehicles', 'Autos & Vehicles',
            yc.category = 'Music', 'Music',
            yc.category = 'Pets & Animals', 'Pets & Animals',
            yc.category = 'Sports', 'Sports',
            yc.category = 'Short Movies', 'Movie & Shows',
            yc.category = 'Travel & Events', 'Events',
            yc.category = 'Gaming', 'Gaming',
            yc.category = 'Videoblogging', 'Vloging',
            yc.category = 'People & Blogs', 'People & Culture',
            yc.category = 'Comedy', 'Comedy',
            yc.category = 'Entertainment', 'Entertainment',
            yc.category = 'News & Politics', 'News & Politics',
            yc.category = 'Howto & Style', 'DIY',
            yc.category = 'Education', 'Education',
            yc.category = 'Science & Technology', 'Science & Technology',
            yc.category = 'Nonprofits & Activism', 'Non profits',
            yc.category = 'Movies', 'Films',
            yc.category = 'Anime/Animation', 'Kids & Animation',
            yc.category = 'Action/Adventure', 'Movie & Shows',
            yc.category = 'Classics', 'Films',
            yc.category = 'Documentary', 'Infotainment',
            yc.category = 'Drama', 'Films',
            yc.category = 'Family', 'Family & Parenting',
            yc.category = 'Foreign', 'Movie & Shows',
            yc.category = 'Horror', 'Movie & Shows',
            yc.category = 'Sci-Fi/Fantasy', 'Movie & Shows',
            yc.category = 'Thriller', 'Movie & Shows',
            yc.category = 'Shorts', 'Movie & Shows',
            yc.category = 'Shows', 'Movie & Shows',
            yc.category = 'Trailers', 'Movie & Shows',
            NULL
        ) yt_category,
        coalesce(
            if(myt.category = '', NULL, myt.category),
            if(cc.category = '', NULL, cc.category),
            yt_category) category,
        [category] categories,
        base.description bio,
        replaceOne(base.thumbnail, '=s88-c', '=s512-c') profile_pic_url,
        profile_pic_url thumbnail,
        multiIf(
                followers < 5000, '0-5k',
                followers < 10000, '5k-10k',
                followers < 50000, '10k-50k',
                followers < 100000, '50k-100k',
                followers < 500000, '100k-500k',
                followers < 1000000, '500k-1M',
                '1M+'
            ) follower_group,
        concat('YT', '_', country, '_', category, '_', follower_group) primary_group,
        concat('YT', '_', country, '_', follower_group) secondary_group,
        if(base.video_count > 0, base.view_count / base.video_count, 0) avg_video_views,
        if(base.video_count > 0, base.view_count / base.video_count, 0) avg_views,
        base.subscriber_count followers,
        base.video_count uploads,
        base.video_count uploads_count,
        base.view_count views,
        base.view_count views_count,
        base.title name,
        base.title title,
        base.description description,
        COALESCE (
                bya.account_type,
                if (myt.profile_type = '', NULL, myt.profile_type),
                'CREATOR') profile_type,
        multiIf(
                followers < 1000, 'Nano',
                followers < 75000, 'Micro',
                followers < 1000000, 'Macro',
                'Mega'
            ) label,
        multiIf(
                bya.language is not null, bya.language,
                myt.language is not null and myt.language != '', myt.language,
                myl.language is not null and myl.language != '', myl.language,
                sp.language like 'en-%', 'en',
                sp.language like 'hi-%', 'en',
                sp.language = 'hi', 'hi',
                sp.language = 'ta', 'ta',
                sp.language = 'bn', 'bn',
                sp.language = 'ml', 'ml',
                sp.language = 'mr', 'mr',
                sp.language = 'gu', 'gu',
                sp.language = 'pa', 'pa',
                sp.language = 'or', 'or',
                sp.language = 'ka', 'ka',
                sp.language) language,
        [language] languages,
        yg.follower_7d followers_growth7d,
        yg.uploads_7d avg_posts_per_week,
        yg.views_7d views_7d,
        yg.follower_30d followers_growth30d,
        yg.uploads_30d uploads_30d,
        yg.views_30d views_30d,
        yg.follower_1y followers_growth1y,
        yg.views_1y views_1y,
        splitByWhitespace(ifNull(concat(lower(name), ' ', lower(username)), '')) keywords,
        concat(lower(name), ' ', lower(username),' ', channel_id) search_phrase,
        0 is_blacklisted,
        0 comments,
        round(if(base.video_count > 0, base.view_count / base.video_count, 0)) video_reach,
        0 shorts_reach,
        0 authentic_engagement,
        0 comment_rate_percentage,
        0 est_reach,
        0 est_impressions,
        0 shorts_impressions,
        0 video_impressions,
        views_30d video_views_last30,
        0 reaction_rate,
        '' dob,
        myt.gender gender,
        myt.city city,
        myt.state state,
        0 comments_rate,
        0 avg_comments,
        0 likes_spread,
        0 followers_growth90d,
        0 audience_reachability,
        0 audience_quality,
        0 likes_to_comment_ratio,
        base.video_count post_count,
        0 engagement_rate,
        0 avg_likes,
        0 audience_authencity,
        0 cpm,
        coalesce(if(sp.contacts_email = '', NULL, sp.contacts_email), myt.email) email,
        sp.contacts_phone phone,
        rank() over (partition by country order by followers desc) country_rank,
        rank() over (partition by country, category order by followers desc) category_rank,
        if (cp.channel_id is not null and cp.channel_id != '', 1, 0) is_collection_profile
    from
        base
        left join {{ ref('stg_beat_youtube_account') }} bya on bya.channel_id = base.channel_id
        left join vidooly.social_profiles_new_20230217 sp on sp.youtube_userid = base.channel_id
        left join vidooly.youtube_category yc on toString(yc.id) = sp.youtube_category
        left join vidooly.category cc on toString(cc.id) = ccatid
        left join dbt.mart_yt_growth yg on yg.channelid = sp.youtube_userid
        left join links on links.channel_id = base.channel_id
        left join collection_profiles cp on cp.channel_id = base.channel_id
        left join dbt.mart_youtube_account_language myl on myl.channel_id = base.channel_id
        left join dbt.mart_manual_data_yt myt on myt.channel_id = base.channel_id
),
final_without_group_metrics as (
    select
        channelid,
        channel_id,
        updated_at,
        username,
        country,
        yt_category,
        linked_instagram_handle,
        category,
        categories,
        bio,
        profile_pic_url,
        thumbnail,
        follower_group,
        primary_group,
        secondary_group,
        if(avg_video_views<0, 0, avg_video_views) avg_video_views,
        if(avg_views<0, 0, avg_views) avg_views,
        followers,
        if(uploads<0, 0, uploads) uploads,
        if(uploads_count<0, 0, uploads_count) uploads_count,
        if(views<0, 0, views) views,
        if(views_count<0, 0, views_count) views_count,
        name,
        title,
        description,
        profile_type,
        label,
        language,
        languages,
        followers_growth7d,
        if(avg_posts_per_week<0, 0, avg_posts_per_week) avg_posts_per_week,
        if(views_7d<0, 0, views_7d) views_7d,
        followers_growth30d,
        if(uploads_30d<0, 0, uploads_30d) uploads_30d,
        if(views_30d<0, 0, views_30d) views_30d,
        followers_growth1y,
        views_1y,
        is_blacklisted,
        comments,
        video_reach,
        shorts_reach,
        authentic_engagement,
        comment_rate_percentage,
        est_reach,
        est_impressions,
        shorts_impressions,
        video_impressions,
        video_views_last30,
        reaction_rate,
        dob,
        gender,
        keywords,
        search_phrase,
        city,
        state,
        comments_rate,
        group_key,
        if(avg_comments<0, 0, avg_comments) avg_comments,
        likes_spread,
        followers_growth90d,
        audience_reachability,
        audience_quality,
        likes_to_comment_ratio,
        if(post_count<0, 0, post_count) post_count,
        engagement_rate,
        avg_likes,
        audience_authencity,
        cpm,
        email,
        phone,
        country_rank,
        category_rank,
        if (pgm.profiles > 0, primary_group, secondary_group) group_key,
        is_collection_profile
    from
    semi_final
    left join dbt.mart_primary_group_metrics pgm on pgm.group_key = semi_final.primary_group
    left join dbt.mart_primary_group_metrics sgm on sgm.group_key = semi_final.secondary_group
),
final as (
    select
        channelid,
        channel_id,
        updated_at,
        username,
        country,
        yt_category,
        linked_instagram_handle,
        category,
        categories,
        bio,
        profile_pic_url,
        thumbnail,
        follower_group,
        primary_group,
        secondary_group,
        avg_video_views,
        avg_views,
        followers,
        uploads,
        uploads_count,
        views,
        views_count,
        name,
        title,
        description,
        profile_type,
        label,
        language,
        languages,
        followers_growth7d,
        avg_posts_per_week,
        views_7d,
        followers_growth30d,
        uploads_30d,
        views_30d,
        followers_growth1y,
        views_1y,
        is_blacklisted,
        comments,
        video_reach,
        shorts_reach,
        authentic_engagement,
        comment_rate_percentage,
        est_reach,
        est_impressions,
        shorts_impressions,
        video_impressions,
        video_views_last30,
        reaction_rate,
        dob,
        gender,
        keywords,
        search_phrase,
        city,
        state,
        comments_rate,
        group_key,
        avg_comments,
        followers_growth30d,
        likes_spread,
        followers_growth90d,
        audience_reachability,
        audience_quality,
        likes_to_comment_ratio,
        post_count,
        engagement_rate,
        avg_likes,
        audience_authencity,
        cpm,
        email,
        phone,
        country_rank,
        category_rank,
        0 reels_reach,
        0 story_reach,
        0 image_reach,
        {% for metric in metrics %}
            multiIf(
                    final_without_group_metrics.{{metric}} is null or final_without_group_metrics.{{metric}} == 0 or isInfinite(final_without_group_metrics.{{metric}}) or isNaN(final_without_group_metrics.{{metric}}), NULL,
                    final_without_group_metrics.{{metric}} > gm.p90_{{metric}}, 'Excellent',
                    final_without_group_metrics.{{metric}} > gm.p75_{{metric}}, 'Very Good',
                    final_without_group_metrics.{{metric}} > gm.p60_{{metric}}, 'Good',
                    final_without_group_metrics.{{metric}} > gm.p40_{{metric}}, 'Average',
                    'Poor'
                ) {{metric}}_grade,
            final_without_group_metrics.{{metric}} {{metric}}_score,
            gm.group_avg_{{metric}} group_avg_{{metric}},
        {% endfor %}
        group_key,
        '' latest_video_publish_time,
        0 avg_shorts_views_30d,
        video_views_last30 avg_video_views_30d,
        transform(followers_growth30d_grade, ['Poor', 'Average', 'Good', 'Very Good', 'Excellent'], [1, 2, 3, 4, 5], 0) AS followers_growth30d_rating,
        transform(views_30d_grade, ['Poor', 'Average', 'Good', 'Very Good', 'Excellent'], [1, 2, 3, 4, 5], 0) AS views_30d_rating,
        transform(uploads_30d_grade, ['Poor', 'Average', 'Good', 'Very Good', 'Excellent'], [1, 2, 3, 4, 5], 0) AS uploads_30d_rating,
        19*(0.5*views_30d_rating + 0.25*followers_growth30d_rating + 0.25*uploads_30d_rating) channel_quality_score,
        transform(
            round(0.5*views_30d_rating + 0.25*followers_growth30d_rating + 0.25*uploads_30d_rating),
            [1, 2, 3, 4, 5],
            ['Poor', 'Average', 'Good', 'Very Good', 'Excellent'],
            'Poor'
        ) channel_quality_grade,
        is_collection_profile,
        0 reels_impressions,
        0 image_impressions
    from
        final_without_group_metrics
    left join dbt.mart_primary_group_metrics gm on gm.group_key = final_without_group_metrics.group_key
)
select * from final SETTINGS max_bytes_before_external_group_by = 20000000000, max_bytes_before_external_sort = 20000000000, join_algorithm='partial_merge'

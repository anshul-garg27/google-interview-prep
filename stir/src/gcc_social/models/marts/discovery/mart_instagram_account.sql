{{ config(materialized = 'table', tags=["hourly"], order_by='ig_id') }}
{%
    set metrics = [
        "avg_likes",
        "avg_comments",
        "comments_rate",
        "followers",
        "engagement_rate",
        "reels_reach",
        "image_reach",
        "story_reach",
        "reels_impressions",
        "image_impressions",
        "likes_to_comment_ratio",
        "followers_growth7d",
        "followers_growth30d",
        "followers_growth90d",
        "audience_reachability",
        "audience_authencity",
        "post_count",
        "likes_spread",
        "avg_posts_per_week"
    ]
%}
with
links as (
    select
        handle, max(channel_id) channel_id, max(category) category
    from
        dbt.mart_linked_socials
    group by handle
),
profile_pics as (
    select
        entity_id profile_id,
        concat('https://d24w28i6lzk071.cloudfront.net/', asset_url) profile_pic_url
    from
    dbt.stg_beat_asset_log
    where entity_type = 'PROFILE'
),
weekly_posts_base as (
    select
        toStartOfWeek(publish_time) week,
        profile_id,
        uniqExact(shortcode) posts
    from {{ref('stg_beat_instagram_post')}}
    group by week, profile_id
),
weekly_posts as (
    select
        profile_id,
        ia.handle handle,
        avg(posts) avg_posts_per_week
    from weekly_posts_base
    left join {{ref('stg_beat_instagram_account')}} ia on ia.profile_id = weekly_posts_base.profile_id
    group by profile_id, handle
),
overtime_growth as (
    select
        platform_id profile_id,
        argMax(followers, date) cur_followers,
        argMaxIf(followers, date, date < toStartOfWeek(now()) - INTERVAL 7 DAY) followers_7d,
        argMaxIf(followers, date, date < toStartOfWeek(now()) - INTERVAL 30 DAY) followers_30d,
        argMaxIf(followers, date, date < toStartOfWeek(now()) - INTERVAL 90 DAY) followers_90d,
        argMaxIf(followers, date, date < toStartOfWeek(now()) - INTERVAL 1 YEAR) followers_1y,
        4*sum(plays)/uniqExact(date) avg_reels_play_count_30d,
        sumIf(uploads, date < toStartOfWeek(now()) - INTERVAL 30 DAY) uploads_30d,
        sumIf(plays, date < toStartOfWeek(now()) - INTERVAL 30 DAY) plays_30d,
        sumIf(reach, date < toStartOfWeek(now()) - INTERVAL 30 DAY) reach_30d,
        sumIf(impressions, date < toStartOfWeek(now()) - INTERVAL 30 DAY) impressions_30d,
        cur_followers - followers_7d followers_growth_7d,
        cur_followers - followers_30d followers_growth_30d,
        cur_followers - followers_90d followers_growth_90d,
        cur_followers - followers_1y followers_growth_1y
    from dbt.mart_insta_account_weekly
    group by profile_id
),
indians as (
    select
        profile_id,
        'IN' country,
        uniqExact(short_code) posts,
        uniqExactIf(short_code, location like '%India') indian_posts
    from
        dbt.mart_post_location
    group by profile_id
    having indian_posts / posts  > 0.5
),
indians3 as (
    select
        profile_social_id profile_id,
        'IN' country
    from dbt.mart_collection_post
    where platform = 'INSTAGRAM'
    group by profile_id
),
insta_data as (
    select
        instagram_username handle,
        max(country) country,
        max(ccatid) ccatid
    from
        vidooly.social_profiles_new_20230217
    where instagram_userid != ''
    group by handle
    having country = 'IN'
),
vidooly_data as (
    select
        s.profile_name handle,
        max(concat('https://cdn.vidooly.com/images', s.saved_thumbnail)) profile_pic_url,
        max(if(country = 'NULL', NULL, country)) country,
        max(if(contact_phone_number = '' or contact_phone_number = 'NULL', NULL, contact_phone_number)) contact_phone_number,
        max(if(public_phone_number = '' or public_phone_number = 'NULL', NULL, public_phone_number)) public_phone_number,
        max(if(email = '' or email = 'NULL', NULL, email)) email,
        max(if(city = '' or city = 'NULL', NULL, splitByChar(',', city)[1])) city,
        max(multiIf(gender  = 'M', 'MALE', gender = 'F', 'FEMALE', NULL)) gender,
        coalesce(public_phone_number, contact_phone_number) phone
    from
        dbt.stg_vidooly_instagram_profiles s
    group by handle
),
semi_final as (
    select
        ifNull(ia.profile_id, 'MISSING') ig_id,
        ifNull(ia.handle, 'MISSING') handle,
        ia.updated_at updated_at,
        links.channel_id linked_channel_id,
        coalesce(
                if(links.category = '', NULL, links.category),
                if(cc.category = '', NULL, cc.category),
                predc.category) category,
        coalesce(
                if(links.category = '' or links.category is NULL, NULL, 'link_social'),
                if(cc.category = '' or cc.category is NULL, NULL, 'vidooly'),
                'predicted') category_source,
        coalesce(
                if(cc.category = '', NULL, cc.category),
                predc.category) category_original,

        coalesce(
                if(ia.language = '', NULL, ia.language),
                if(migbd.language = '', NULL, migbd.language),
                if(predl.lang = '', NULL, predl.lang),
                'en') language,
        coalesce(
                if(ia.language = '' or ia.language is NULL, NULL, 'beat'),
                if(migbd.language = '' or migbd.language is NULL, NULL, 'gpt'),
                if(predl.lang = '' or predl.lang is NULL, NULL, 'predicted')
                ) language_source,
        pc.phone phone,
        pc.email email,
        ia.biography bio,
        ia.fbid business_id,
        ia.full_name name,
        ia.followers followers,
        ia.following following,
        COALESCE (ia.account_type, 'CREATOR') profile_type,
        multiIf(
            ia.profile_type = 'personal', 'PERSONAL',
            ia.profile_type = 'business', 'BUSINESS',
            ia.profile_type = 'creator', 'CREATOR',
            ia.is_business_or_creator = 1, 'BUSINESS_OR_CREATOR',
            'UNKNOWN'
        ) account_type,
        multiIf(
                followers < 5000, '0-5k',
                followers < 10000, '5k-10k',
                followers < 50000, '10k-50k',
                followers < 100000, '50k-100k',
                followers < 500000, '100k-500k',
                followers < 1000000, '500k-1M',
                '1M+'
            ) follower_group,
        ifNull(concat('INSTA', '_', coalesce(country, 'MISSING'), '_', coalesce(category, 'MISSING'), '_', coalesce(follower_group, 'MISSING')), 'MISSING') primary_group,
        ifNull(concat('INSTA', '_', coalesce(country, 'MISSING'), '_', coalesce(follower_group, 'MISSING')), 'MISSING') secondary_group,
        coalesce(
                if(ia.country = '', NULL, ia.country),
                if(c.country = '', NULL, c.country),
                if(pc.country = '', NULL, pc.country),
                if(migbd.country = '', NULL, migbd.country),
                if(indians.country = '', NULL, indians.country)
        ) country,
        coalesce(
                if(ia.country = '' or ia.country is NULL, NULL, 'beat'),
                if(c.country = '' or c.country is NULL, NULL, 'vidooly_sp'),
                if(pc.country = '' or pc.country is NULL, NULL, 'vidooly_ia'),
                if(migbd.country = '' or migbd.country is NULL, NULL, 'gpt'),
                if(indians.country = '' or indians.country is NULL, NULL, 'posts')
        ) country_source,
        coalesce(
                if(pc.city = '', NULL, pc.city),
                if(migbd.city = '', NULL, migbd.city)
        ) city,
        coalesce(
                if(pc.city = '' or pc.city is NULL, NULL, 'vidooly_ia'),
                if(migbd.city = '' or migbd.city is NULL, NULL, 'gpt')
        ) city_source,
        migbd.state state,
        multiIf(
                followers < 1000, 'Nano',
                followers < 75000, 'Micro',
                followers < 1000000, 'Macro',
                'Mega'
            ) label,
        coalesce(
                if(pc.gender = '', NULL, pc.gender),
                if(migbd.gender = '', NULL, migbd.gender)
        ) gender,
        coalesce(
                if(pc.gender = '' or pc.gender is NULL, NULL, 'vidooly_ia'),
                if(migbd.gender = '' or migbd.gender is NULL, NULL, 'gpt')
        ) gender_source,
        '' dob,
        0 is_blacklisted,
        0 authentic_engagement,
        og.uploads_30d uploads_30d,
        og.plays_30d plays_30d,
        og.reach_30d reach_30d,
        og.impressions_30d impressions_30d,
        og.followers_growth_7d followers_growth7d,
        og.followers_growth_30d followers_growth30d,
        og.followers_growth_90d followers_growth90d,
        og.followers_growth_1y followers_growth1y,
        ifNull(og.avg_reels_play_count_30d, 0.0) avg_reels_play_30d,
        ifNull(aud.reachable_followers_perc, 0) audience_reachability,
        ifNull(aud.audience_authenticity, 0) audience_authencity,
        ifNull(aud.audience_quality_score, 0) audience_quality,
        aud.audience_quality_grade audience_quality_grade,
        0 post_frequency_week,
        0 is_verified,
        ia.is_private is_private,
        stats.likes_spread_percentage likes_spread_percentage,
        round(100 * avg_comments / followers) comment_rate,
        stats.ffratio ffratio,
        coalesce(pcn.profile_pic_url, pc.profile_pic_url) profile_pic_url,
        weekly_posts.avg_posts_per_week avg_posts_per_week,
        stats.total_posts post_count,
        post_count num_of_posts,
        if(isInfinite(if(num_of_posts = 0, 0, stats."avg_reach")), 0,if(num_of_posts = 0, 0, stats."avg_reach"))  avg_reach_all,
        if(num_of_posts = 0, 0, stats."avg_reach") avg_reach_non_outliers,
        if(isInfinite(coalesce(if(avg_reach_non_outliers = 0, NULL, avg_reach_non_outliers), avg_reach_all)), 0,coalesce(if(avg_reach_non_outliers = 0, NULL, avg_reach_non_outliers), avg_reach_all))  avg_reach,
    
        if(num_of_posts = 0, 0, stats."avg_likes") avg_likes_all,
        if(num_of_posts = 0, 0, stats."avg_likes_non_outliers") avg_likes_non_outliers,
        coalesce(if(avg_likes_non_outliers = 0, NULL, avg_likes_non_outliers), avg_likes_all) avg_likes,
        

        if(num_of_posts = 0, 0, stats."avg_reels_plays") avg_reels_plays_all,
        if(num_of_posts = 0, 0, stats."avg_reels_plays_non_outliers") avg_reels_plays_non_outliers,
        coalesce(if(avg_reels_plays_non_outliers = 0, NULL, avg_reels_plays_non_outliers), avg_reels_plays_all) avg_reels_play_count,


        if(num_of_posts = 0, 0, stats."avg_er") avg_er_all,
        if(num_of_posts = 0, 0, stats."avg_er_non_outliers") avg_er_non_outliers,
        if(isInfinite(coalesce(if(avg_er_non_outliers = 0, NULL, avg_er_non_outliers), avg_er_all)), 0, coalesce(if(avg_er_non_outliers = 0, NULL, avg_er_non_outliers), avg_er_all)) engagement_rate,


        if(num_of_posts = 0, 0, stats."avg_reels_reach") avg_reels_reach_all,
        if(num_of_posts = 0, 0, stats."avg_reels_reach_non_outliers") avg_reels_reach_non_outliers,
        coalesce(if(avg_reels_reach_non_outliers = 0, NULL, avg_reels_reach_non_outliers), avg_reels_reach_all) reels_reach,
        
        stats.story_reach story_reach,

        if(num_of_posts = 0, 0, stats."avg_static_reach") avg_static_reach_all,
        if(num_of_posts = 0, 0, stats."static_reach_non_outliers") avg_static_reach_non_outliers,
        coalesce(if(avg_static_reach_non_outliers = 0, NULL, avg_static_reach_non_outliers), avg_static_reach_all) image_reach,


        if(num_of_posts = 0, 0, stats."avg_comments") avg_comments_all,
        if(num_of_posts = 0, 0, stats."avg_comments_non_outliers") avg_comments_non_outliers,
        coalesce(if(avg_comments_non_outliers = 0, NULL, avg_comments_non_outliers), avg_comments_all) avg_comments,

        
        avg_reels_play_count avg_views,
        stats.likes_spread_percentage likes_spread_percentage,
        stats.avg_reels_impressions reels_impressions,
        stats.avg_static_impressions image_impressions,
        if (indians3.profile_id is not null and indians3.profile_id != '', 1, 0) is_collection_profile
    from
        {{ref('mart_instagram_account_summary')}} stats
        left join {{ref('stg_beat_instagram_account')}} ia on ia.profile_id = stats.profile_id
        left join {{ref('mart_audience_info')}} aud on aud.platform_id = ia.profile_id and aud.platform = 'INSTAGRAM'
        left join insta_data c on c.handle = ia.handle
        left join vidooly_data pc on pc.handle = ia.handle
        left join links on links.handle = ia.handle
        left join vidooly.category cc on toString(cc.id) = insta_data.ccatid
        left join weekly_posts on weekly_posts.handle = ia.handle
        left join indians on indians.profile_id = ia.profile_id
        left join indians3 on indians3.profile_id = ia.profile_id
        left join dbt.mart_insta_predicted_categories predc on predc.profile_id = ia.profile_id
        left join dbt.mart_insta_predicted_lang predl on predl.profile_id = ia.profile_id
        left join overtime_growth og on og.profile_id = ia.profile_id
        left join profile_pics pcn on pcn.profile_id = ia.profile_id
        left join dbt.mart_instagram_gpt_basic_data migbd on migbd.handle = ia.handle
),
final_without_group_metrics_with_dups as (
    select
        handle,
        name,
        business_id,
        ig_id,
        updated_at,
        profile_pic_url thumbnail,
        account_type,
        concat(lower(name), ' ', lower(handle)) search_phrase,
        linked_channel_id,
        dob,
        city,
        city_source,
        state,
        label,
        gender,
        gender_source,
        is_blacklisted,
        if(avg_posts_per_week<0, 0, avg_posts_per_week) avg_posts_per_week,
        splitByWhitespace(ifNull(concat(lower(name), ' ', lower(handle)), '')) keywords,
        phone,
        email,
        bio,
        reels_impressions,
        image_impressions,
        [category] categories,
        [language] languages,
        language primary_language,
        language_source,
        category primary_category,
        category_original primary_category_original,
        category_source,
        country,
        country_source,
        followers,
        following,
        profile_type,
        ffratio,
        uploads_30d,
        plays_30d,
        reach_30d,
        impressions_30d,
        followers_growth7d,
        followers_growth30d,
        followers_growth90d,
        followers_growth1y,
        if(avg_reels_play_30d<0, 0, avg_reels_play_30d) avg_reels_play_30d,
        audience_reachability,
        audience_authencity,
        audience_quality,
        audience_quality_grade,
        is_verified,
        is_private,
        authentic_engagement,
        post_frequency_week,
        likes_spread_percentage likes_spread,
        likes_spread_percentage,
        100 * avg_comments / followers comments_rate,
        100 * avg_comments / followers comments_rate_percentage,
        if(post_count<0, 0, post_count) post_count,
        if(avg_views<0, 0, avg_views) avg_views,
        if(avg_reach<0, 0, avg_reach) avg_reach,
        if(avg_likes<0, 0, avg_likes) avg_likes,
        avg_reels_play_count,
        engagement_rate,
        if(reels_reach<0, 0, reels_reach) reels_reach,
        story_reach,
        image_reach,
        if(avg_comments<0, 0, avg_comments) avg_comments,
        search_phrase,
        avg_likes / avg_comments likes_to_comment_ratio,
        rank() over (partition by country order by followers desc) country_rank,
        rank() over (partition by country, category order by followers desc) category_rank,
        comment_rate comment_rate_percentage,
        primary_group,
        secondary_group,
        if (pgm.profiles > 0, primary_group, secondary_group) group_key,
        is_collection_profile
    from semi_final
    left join dbt.mart_primary_group_metrics pgm on pgm.group_key = semi_final.primary_group
    left join dbt.mart_primary_group_metrics sgm on sgm.group_key = semi_final.secondary_group
),
final_without_group_metrics as (
    select
        handle,
        max(updated_at) _updated_at,
        argMax(name, updated_at) name,
        argMax(business_id, updated_at) business_id,
        argMax(ig_id, updated_at) ig_id,
        argMax(thumbnail, updated_at) thumbnail,
        argMax(account_type, updated_at) account_type,
        argMax(search_phrase, updated_at) search_phrase,
        argMax(linked_channel_id, updated_at) linked_channel_id,
        argMax(dob, updated_at) dob,
        argMax(city, updated_at) city,
        argMax(city_source, updated_at) city_source,
        argMax(state, updated_at) state,
        argMax(label, updated_at) label,
        argMax(gender, updated_at) gender,
        argMax(gender_source, updated_at) gender_source,
        argMax(is_blacklisted, updated_at) is_blacklisted,
        argMax(avg_posts_per_week, updated_at) avg_posts_per_week,
        argMax(keywords, updated_at) keywords,
        argMax(phone, updated_at) phone,
        argMax(email, updated_at) email,
        argMax(bio, updated_at) bio,
        argMax(reels_impressions, updated_at) reels_impressions,
        argMax(image_impressions, updated_at) image_impressions,
        argMax(categories, updated_at) categories,
        argMax(languages, updated_at) languages,
        argMax(primary_language, updated_at) primary_language,
        argMax(language_source, updated_at) language_source,
        argMax(primary_category, updated_at) primary_category,
        argMax(category_source, updated_at) category_source,
        argMax(primary_category_original, updated_at) primary_category_original,
        argMax(country, updated_at) country,
        argMax(country_source, updated_at) country_source,
        argMax(followers, updated_at) followers,
        argMax(following, updated_at) following,
        argMax(profile_type, updated_at) profile_type,
        argMax(ffratio, updated_at) ffratio,
        argMax(uploads_30d, updated_at) uploads_30d,
        argMax(plays_30d, updated_at) plays_30d,
        argMax(reach_30d, updated_at) reach_30d,
        argMax(impressions_30d, updated_at) impressions_30d,
        argMax(followers_growth7d, updated_at) followers_growth7d,
        argMax(followers_growth30d, updated_at) followers_growth30d,
        argMax(followers_growth90d, updated_at) followers_growth90d,
        argMax(followers_growth1y, updated_at) followers_growth1y,
        argMax(avg_reels_play_30d, updated_at) avg_reels_play_30d,
        argMax(audience_reachability, updated_at) audience_reachability,
        argMax(audience_authencity, updated_at) audience_authencity,
        argMax(audience_quality, updated_at) audience_quality,
        argMax(audience_quality_grade, updated_at) audience_quality_grade,
        argMax(is_verified, updated_at) is_verified,
        argMax(is_private, updated_at) is_private,
        argMax(authentic_engagement, updated_at) authentic_engagement,
        argMax(post_frequency_week, updated_at) post_frequency_week,
        argMax(likes_spread, updated_at) likes_spread,
        argMax(likes_spread_percentage, updated_at) likes_spread_percentage,
        argMax(comments_rate, updated_at) comments_rate,
        argMax(comments_rate_percentage, updated_at) comments_rate_percentage,
        argMax(post_count, updated_at) post_count,
        argMax(avg_views, updated_at) avg_views,
        argMax(avg_reach, updated_at) avg_reach,
        argMax(avg_likes, updated_at) avg_likes,
        argMax(avg_reels_play_count, updated_at) avg_reels_play_count,
        argMax(engagement_rate, updated_at) engagement_rate,
        argMax(reels_reach, updated_at) reels_reach,
        argMax(story_reach, updated_at) story_reach,
        argMax(image_reach, updated_at) image_reach,
        argMax(avg_comments, updated_at) avg_comments,
        argMax(search_phrase, updated_at) search_phrase,
        argMax(likes_to_comment_ratio, updated_at) likes_to_comment_ratio,
        argMax(country_rank, updated_at) country_rank,
        argMax(category_rank, updated_at) category_rank,
        argMax(comment_rate_percentage, updated_at) comment_rate_percentage,
        argMax(primary_group, updated_at) primary_group,
        argMax(secondary_group, updated_at) secondary_group,
        argMax(group_key, updated_at) group_key,
        argMax(is_collection_profile, updated_at) is_collection_profile
    from final_without_group_metrics_with_dups
    group by handle
),
final as (
    select
        name,
        handle,
        business_id,
        _updated_at updated_at,
        ifNull(ig_id, 'MISSING') ig_id,
        thumbnail,
        account_type,
        keywords,
        bio,
        categories,
        linked_channel_id,
        dob,
        city,
        city_source,
        state,
        label,
        gender,
        gender_source,
        is_blacklisted,
        languages,
        primary_language,
        language_source,
        primary_category,
        category_source,
        primary_category_original,
        country,
        country_source,
        followers,
        following,
        profile_type,
        ffratio,
        followers_growth7d,
        followers_growth30d,
        followers_growth90d,
        followers_growth1y,
        avg_reels_play_30d,
        audience_reachability,
        audience_authencity,
        audience_quality,
        audience_quality_grade,
        likes_spread,
        avg_posts_per_week,
        likes_spread_percentage,
        post_count,
        avg_views,
        avg_reach,
        avg_likes,
        likes_to_comment_ratio,
        avg_reels_play_count,
        engagement_rate,
        comments_rate,
        comment_rate_percentage,
        reels_reach,
        if(isNaN(story_reach), NULL, story_reach) story_reach,
        image_reach,
        engagement_rate_grade er_grade,
        image_impressions,
        reels_impressions,
        is_verified,
        is_private,
        post_frequency_week,
        authentic_engagement,
        avg_comments,
        country_rank,
        category_rank,
        primary_group,
        secondary_group,
        phone,
        email,
        search_phrase,
        0 shorts_reach,
        0 video_reach,
        uploads_30d,
        plays_30d,
        reach_30d,
        impressions_30d,
        plays_30d views_30d,
        {% for metric in metrics %}
            multiIf(
                    {{metric}} is null or {{metric}} == 0 or isInfinite({{metric}}) or isNaN({{metric}}), NULL,
                    {{metric}} > gm.p90_{{metric}}, 'Excellent',
                    {{metric}} > gm.p75_{{metric}}, 'Very Good',
                    {{metric}} > gm.p60_{{metric}}, 'Good',
                    {{metric}} > gm.p40_{{metric}}, 'Average',
                    'Poor'
                ) {{metric}}_grade,
            gm.group_avg_{{metric}} group_avg_{{metric}},
        {% endfor %}
        group_key,
        is_collection_profile
    from final_without_group_metrics
    left join dbt.mart_primary_group_metrics gm
        on gm.group_key = final_without_group_metrics.group_key
)
select * from final
order by followers desc

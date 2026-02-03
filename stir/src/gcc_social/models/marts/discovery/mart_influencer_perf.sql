{{ config(materialized = 'table', tags=["hourly"], order_by='handle') }}
with
campaign_types as (
    select
        collection_id collection_id,
        if(uniqExactIf(post_short_code, link_clicks > 0) > 0, 1, 0) click_campaign,
        if(uniqExactIf(post_short_code, orders > 0) > 0, 1, 0) order_campaign
    from mart_collection_post
    where collection_type = 'POST'
    group by collection_id
),
campaign_costs as (
    select
        post_collection_id collection_id,
        posted_by_handle handle,
        max(cost) cost
    from dbt.stg_coffee_post_collection_item cpi
    group by collection_id, handle
),
campaign_data_collection_level as (
	select
	    mpc.platform platform,
	    mpc.collection_id collection_id,
	    max(ct.click_campaign) click_campaign,
	    max(ct.order_campaign) order_campaign,
	    uniqExact(short_code) posts,
		ifNull(profile_handle, 'MISSING') handle,
		sum(orders) orders,
		sum(link_clicks) clicks,
		sum(views) views,
		sum(likes) likes,
		sum(mpc.comments) comments,
		sum(reach) reach,
		likes + comments total_engagement,
		sum(impressions) impressions,
        max(cc.cost) cost
	from mart_collection_post mpc
	left join dbt.stg_coffee_post_collection cp on cp.id = mpc.collection_id
	left join dbt.stg_coffee_post_collection_item cpi on cpi.post_collection_id = mpc.collection_id and cpi.short_code = mpc.post_short_code
	left join campaign_costs cc on cc.collection_id = mpc.collection_id and cc.handle = mpc.profile_handle
	left join campaign_types ct on ct.collection_id = mpc.collection_id
	where platform = 'INSTAGRAM' and cp.source = 'GCC_CAMPAIGN'
	group by platform, collection_id, handle
),
campaign_data as (
    select
        platform,
        handle,
        uniqExact(collection_id) campaigns,
        uniqExactIf(collection_id, click_campaign = 1) click_campaigns,
        uniqExactIf(collection_id, order_campaign = 1) order_campaigns,
        sum(posts) posts,
        sum(orders) orders,
        sum(clicks) clicks,
        sum(views) views,
        sum(likes) likes,
        sum(mpc.comments) comments,
        sum(reach) reach,
        sum(total_engagement) total_engagement,
        sum(impressions) impressions,
        sum(cost) total_cost,
        clicks / orders cto,
        round(clicks / reach, 6) ctr,
        sum(cost) total_cost,
        total_cost / views cpv,
		total_cost / clicks cpc,
		total_cost / orders cpo,
		total_cost / reach cpr,
		total_cost / (impressions*1000) cpm,
		total_cost / total_engagement cpe
    from
        campaign_data_collection_level mpc
    group by platform, handle
),
trinity_data_insta as (
	select
	    'INSTAGRAM' platform,
		ifNull(handle, 'MISSING') handle,
		"g3.Name" "g3_Name",
        "g3.City" "g3_City",
        "g3.State" "g3_State",
        "g3.Tier" "g3_Tier",
        "g3.Brands Interested" "g3_Brands_Interested",
        "g3.Brands Purchased From" "g3_Brands_Purchased_From",
        "g3.Categories Interested" "g3_Categories_Interested",
        "g3.Concerns" "g3_Concerns",
        "g3.Ingredients Interested" "g3_Ingredients_Interested",
        "g3.Price preference" "g3_Price_preference",
        "g3.Discount preference" "g3_Discount_preference",
        "g3.Preferred channel" "g3_Preferred_channel",
        "g3.Preferred day" "g3_Preferred_day",
        "g3.Preferred hour" "g3_Preferred_hour",
        "g3.Persona" "g3_Persona",
        "g3.Insta Handle" "g3_Insta_Handle",
        "profiles_i_follow" "profiles_i_follow",
        "my_top_influencers" "my_top_influencers",
        "my_top_brands" "my_top_brands",
        "brand_mentions" "brand_mentions",
        "time_slots" "time_slots",
        "age" "age"
	from
		dbt.mart_trinity_phase_two
),
social_data_insta as (
	select
	    'INSTAGRAM' platform,
		ifNull(handle, 'MISSING') handle,
		country,
		primary_category,
		label,
		gender,
		primary_language,
		followers,
		following,
		ffratio,
		avg_views,
		avg_reach,
		avg_likes,
		avg_reels_play_count,
		engagement_rate,
		reels_reach,
		story_reach,
		image_reach,
		country_rank,
		category_rank
	from
		dbt.mart_instagram_account
),
app_data as (
    select
        sa.platform,
        sa.handle,
        hp.user_categories uc,
        sa.account_id account_id,
        hp.webengage_user_id webengage_user_id,
        arrayMap(x -> JSONExtractString(x, 'name'),
            JSONExtractArrayRaw(ifNull(hp.user_categories, ''))) user_categories
    from dbt.stg_identity_social_account sa
    left join dbt.stg_identity_host_profile hp on hp.user_client_account_id = sa.account_id
    where sa.platform = 'INSTAGRAM' and sa.account_id is not null and sa.account_id > 0
),
audience_data as (
    select
        handle handle,
        audience_gender_age,
        audience_locale,
        audience_city,
        audience_country
    from
        dbt.stg_beat_instagram_profile_insights i
),
final as (
	select
	    ifNull(app_data.handle, 'MISSING') handle,
	    campaign_data.campaigns campaigns,
	    campaign_data.click_campaigns click_campaigns,
	    campaign_data.order_campaigns order_campaigns,
        campaign_data.posts posts,
        campaign_data.orders orders,
        campaign_data.clicks clicks,
        campaign_data.views views,
        campaign_data.likes likes,
        campaign_data.comments comments,
        campaign_data.reach reach,
        campaign_data.total_engagement total_engagement,
        campaign_data.impressions impressions,
        campaign_data.total_cost total_cost,
        campaign_data.cto cto,
        campaign_data.ctr ctr,
        campaign_data.total_cost total_cost,
        campaign_data.cpv cpv,
        campaign_data.cpc cpc,
        campaign_data.cpo cpo,
        campaign_data.cpr cpr,
        campaign_data.cpm cpm,
        campaign_data.cpe cpe,
	    trinity_data.g3_Name g3_Name,
        trinity_data.g3_City g3_City,
        trinity_data.g3_State g3_State,
        trinity_data.g3_Tier g3_Tier,
        trinity_data.g3_Brands_Interested g3_Brands_Interested,
        trinity_data.g3_Brands_Purchased_From g3_Brands_Purchased_From,
        trinity_data.g3_Categories_Interested g3_Categories_Interested,
        trinity_data.g3_Concerns g3_Concerns,
        trinity_data.g3_Ingredients_Interested g3_Ingredients_Interested,
        trinity_data.g3_Price_preference g3_Price_preference,
        trinity_data.g3_Discount_preference g3_Discount_preference,
        trinity_data.g3_Preferred_channel g3_Preferred_channel,
        trinity_data.g3_Preferred_day g3_Preferred_day,
        trinity_data.g3_Preferred_hour g3_Preferred_hour,
        trinity_data.g3_Persona g3_Persona,
        trinity_data.g3_Insta_Handle g3_Insta_Handle,
        trinity_data.profiles_i_follow profiles_i_follow,
        trinity_data.my_top_influencers my_top_influencers,
        trinity_data.my_top_brands my_top_brands,
        trinity_data.brand_mentions brand_mentions,
        trinity_data.time_slots time_slots,
        trinity_data.age age,
		social_data.country country,
        social_data.primary_category primary_category,
        social_data.label label,
        social_data.gender gender,
        social_data.primary_language primary_language,
        social_data.followers followers,
        social_data.following following,
        social_data.ffratio ffratio,
        social_data.avg_views avg_views,
        social_data.avg_reach avg_reach,
        social_data.avg_likes avg_likes,
        social_data.avg_reels_play_count avg_reels_play_count,
        social_data.engagement_rate engagement_rate,
        social_data.reels_reach reels_reach,
        social_data.story_reach story_reach,
        social_data.image_reach image_reach,
        social_data.country_rank country_rank,
        social_data.category_rank category_rank,
        app_data.user_categories user_categories,
        app_data.account_id account_id,
        app_data.webengage_user_id webengage_user_id,
        audience_data.audience_city audience_city,
        audience_data.audience_country audience_country,
        audience_data.audience_gender_age audience_gender_age,
        audience_data.audience_locale audience_locale
	from app_data
	left join trinity_data_insta trinity_data on trinity_data.handle = app_data.handle
	left join social_data_insta social_data on social_data.handle = app_data.handle
	left join campaign_data campaign_data on campaign_data.handle = app_data.handle
	left join audience_data audience_data on audience_data.handle = app_data.handle
)
select * from final order by orders desc

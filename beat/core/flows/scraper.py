from loguru import logger
from core.client.upload_assets import asset_upload_flow, asset_upload_flow_stories
from core.helpers.session import sessionize
from core.models.models import ScrapeRequestLog
from gpt.flows.fetch_gpt_data import refresh_instagram_gpt_data_audience_age_gender, \
    refresh_instagram_gpt_data_audience_cities, refresh_instagram_gpt_data_base_location, \
    refresh_instagram_gpt_data_base_gender, refresh_instagram_gpt_data_base_categ_lang_topics, \
    refresh_instagram_gpt_data_base_gender, refresh_instagram_gpt_data_base_categ_lang_topics, \
    refresh_instagram_gpt_data_gender_location_lang
from instagram.flows.profile_extra import fetch_profile_followers, fetch_profile_following, fetch_post_likes, \
    fetch_post_comments
from instagram.flows.profile_extra import fetch_hashtag_posts
from instagram.flows.refresh_profile import refresh_post_by_shortcode, refresh_stories_posts, refresh_story_insights, \
    refresh_profile_basic, refresh_profile_custom, refresh_story_posts_by_profile_id
from instagram.flows.refresh_profile import refresh_profile_by_handle
from instagram.flows.refresh_profile import refresh_profile_by_profile_id
from instagram.flows.refresh_profile import refresh_post_insights
from instagram.flows.refresh_profile import refresh_profile_insights
from instagram.flows.refresh_profile import refresh_tagged_posts_by_profile_id
from youtube.flows.csv_jobs import (fetch_ad_vs_organic_csv, fetch_channel_audience_geography_csv, fetch_channel_daily_stats_csv, fetch_channel_demographic_csv,
                                    fetch_channel_engagement_csv, fetch_channel_monthly_stats_csv, fetch_channels_csv, fetch_channel_videos_csv,
                                    fetch_keyword_videos_csv, fetch_playlist_videos_csv, fetch_video_details_csv)
from youtube.flows.profile_extra import fetch_yt_post_comments
from youtube.flows.refresh_profile import refresh_yt_post_type, refresh_yt_profile_relationship, refresh_yt_profiles
from youtube.flows.refresh_profile import refresh_yt_posts
from youtube.flows.refresh_profile import refresh_yt_profile_insights
from youtube.flows.refresh_profile import refresh_yt_posts_by_channel_id
from youtube.flows.refresh_profile import refresh_yt_posts_by_playlist_id
from youtube.flows.refresh_profile import refresh_yt_posts_by_genre
from youtube.flows.refresh_profile import refresh_yt_posts_by_search, refresh_yt_activities
from shopify.flows.refresh_orders import refresh_orders_by_store

_whitelist = [
              refresh_profile_custom,
              refresh_profile_by_profile_id,
              refresh_profile_by_handle,
              refresh_profile_basic,
              refresh_post_by_shortcode,
              refresh_post_insights,
              refresh_profile_insights,
              refresh_tagged_posts_by_profile_id,
              fetch_profile_followers,
              fetch_hashtag_posts,
              asset_upload_flow,
              refresh_yt_profiles,
              refresh_yt_posts,
              refresh_yt_profile_insights,
              fetch_profile_followers,
              asset_upload_flow_stories,
              fetch_hashtag_posts,
              refresh_stories_posts,
              refresh_story_insights,
              refresh_yt_posts_by_channel_id,
              refresh_yt_posts_by_genre,
              refresh_yt_posts_by_playlist_id,
              refresh_yt_posts_by_search,
              refresh_yt_post_type,
              refresh_orders_by_store,
              fetch_profile_following,
              fetch_post_likes,
              fetch_post_comments,
              fetch_yt_post_comments,
              refresh_instagram_gpt_data_base_gender,
              refresh_instagram_gpt_data_base_location,
              refresh_instagram_gpt_data_base_categ_lang_topics,
              refresh_instagram_gpt_data_audience_cities,
              refresh_instagram_gpt_data_audience_age_gender,
              refresh_story_posts_by_profile_id,
              fetch_channels_csv,
              fetch_channel_videos_csv,
              fetch_playlist_videos_csv,
              fetch_keyword_videos_csv,
              fetch_video_details_csv,
              fetch_channel_daily_stats_csv,
              fetch_channel_monthly_stats_csv,
              fetch_channel_demographic_csv,
              refresh_yt_activities,
              fetch_channel_audience_geography_csv,
              fetch_channel_engagement_csv,
              refresh_instagram_gpt_data_gender_location_lang,
              fetch_ad_vs_organic_csv,
              refresh_yt_profile_relationship
              ]

csv_flows = ["fetch_keyword_videos_csv", "fetch_playlist_videos_csv", "fetch_channels_csv", "fetch_channel_videos_csv",
             "fetch_video_details_csv","fetch_channel_demographic_csv", "fetch_channel_daily_stats_csv",
             "fetch_channel_monthly_stats_csv","fetch_channel_audience_geography_csv", "fetch_channel_engagement_csv",
             "fetch_ad_vs_organic_csv"]
@sessionize
async def perform_scrape_task(scrape_log: ScrapeRequestLog, session=None) -> dict | None:
    if not scrape_log or not scrape_log.flow or scrape_log.flow not in globals():
        raise Exception("Missing Flow - %s" % scrape_log.flow)
    _flow = globals()[scrape_log.flow]
    if _flow.__name__ in csv_flows:
        scrape_log.params['scrape_id'] = scrape_log.id
    return await _flow(**scrape_log.params, session=session)

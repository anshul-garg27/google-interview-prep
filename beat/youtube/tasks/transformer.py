"""
    Map scrapper keys to entity keys
"""
from youtube.entities.entities import YoutubeAccount, YoutubePost, YoutubeAccountTimeSeries, \
                                      YoutubePostTimeSeries, YoutubeProfileInsights
from youtube.metric_dim_store import KEYWORDS, TITLE, CUSTOM_URL, VIDEO_COUNT, VIEW_COUNT, SUBSCRIBER_COUNT, \
    HIDDEN_SUBSCRIBER_COUNT, DESCRIPTION, PUBLISH_TIME, UPLOADS_PLAYLIST, COUNTRY, THUMBNAIL, THUMBNAILS, \
    LIKES, FAVOURITE_COUNT, TAGS, CATEGORY_ID, AUDIO_LANGUAGE, IS_LICENSED, CHANNEL_TITLE, DURATION, POST_TYPE, \
    CONTENT_RATING, PLAYLIST_ID, PRIVACY_STATUS, CITY, GENDER_AGE, COMMENTS, CATEGORY, EXTERNAL_URLS, IS_VERIFIED

profile_map = {
    TITLE: YoutubeAccount.title.name,
    CUSTOM_URL: YoutubeAccount.custom_url.name,
    VIDEO_COUNT: YoutubeAccount.uploads.name,
    VIEW_COUNT: YoutubeAccount.views.name,
    EXTERNAL_URLS: YoutubeAccount.external_urls.name,
    IS_VERIFIED: YoutubeAccount.is_verified.name,
    SUBSCRIBER_COUNT: YoutubeAccount.subscribers.name,
    HIDDEN_SUBSCRIBER_COUNT: YoutubeAccount.is_subscriber_count_hidden.name,
    DESCRIPTION: YoutubeAccount.description.name,
    PUBLISH_TIME: YoutubeAccount.published_at.name,
    UPLOADS_PLAYLIST: YoutubeAccount.uploads_playlist_id.name,
    COUNTRY: YoutubeAccount.country.name,
    THUMBNAIL: YoutubeAccount.thumbnail.name,
    THUMBNAILS: YoutubeAccount.thumbnails.name,
}

profile_map_ts = {
    TITLE: YoutubeAccountTimeSeries.title.name,
    CUSTOM_URL: YoutubeAccountTimeSeries.custom_url.name,
    VIDEO_COUNT: YoutubeAccountTimeSeries.uploads.name,
    VIEW_COUNT: YoutubeAccountTimeSeries.views.name,
    SUBSCRIBER_COUNT: YoutubeAccountTimeSeries.subscribers.name,
    HIDDEN_SUBSCRIBER_COUNT: YoutubeAccountTimeSeries.is_subscriber_count_hidden.name,
    DESCRIPTION: YoutubeAccountTimeSeries.description.name,
    PUBLISH_TIME: YoutubeAccountTimeSeries.published_at.name,
    UPLOADS_PLAYLIST: YoutubeAccountTimeSeries.uploads_playlist_id.name,
    COUNTRY: YoutubeAccountTimeSeries.country.name,
    THUMBNAIL: YoutubeAccountTimeSeries.thumbnail.name,
    THUMBNAILS: YoutubeAccountTimeSeries.thumbnails.name,
}

post_map = {
    TITLE: YoutubePost.title.name,
    VIEW_COUNT: YoutubePost.views.name,
    LIKES: YoutubePost.likes.name,
    COMMENTS: YoutubePost.comments.name,
    FAVOURITE_COUNT: YoutubePost.favourites.name,
    DESCRIPTION: YoutubePost.description.name,
    PUBLISH_TIME: YoutubePost.published_at.name,
    THUMBNAIL: YoutubePost.thumbnail.name,
    THUMBNAILS: YoutubePost.thumbnails.name,
    TAGS: YoutubePost.tags.name,
    CATEGORY_ID: YoutubePost.category_id.name,
    CATEGORY: YoutubePost.category.name,
    AUDIO_LANGUAGE: YoutubePost.audio_language.name,
    IS_LICENSED: YoutubePost.is_licensed.name,
    CHANNEL_TITLE: YoutubePost.channel_title.name,
    DURATION: YoutubePost.duration.name,
    CONTENT_RATING: YoutubePost.content_rating.name,
    PLAYLIST_ID: YoutubePost.playlist_id.name,
    PRIVACY_STATUS: YoutubePost.privacy_status.name,
    POST_TYPE: YoutubePost.post_type.name,
    KEYWORDS: YoutubePost.keywords.name
}

# post_map_ts = {
#     TITLE: YoutubePostTimeSeries.title.name,
#     VIEW_COUNT: YoutubePostTimeSeries.views.name,
#     LIKE_COUNT: YoutubePostTimeSeries.likes.name,
#     FAVOURITE_COUNT: YoutubePostTimeSeries.favourites.name,
#     DESCRIPTION: YoutubePostTimeSeries.description.name,
#     PUBLISH_TIME: YoutubePostTimeSeries.published_at.name,
#     THUMBNAIL: YoutubePostTimeSeries.thumbnail.name,
#     THUMBNAILS: YoutubePostTimeSeries.thumbnails.name,
#     TAGS: YoutubePostTimeSeries.tags.name,
#     CATEGORY_ID: YoutubePostTimeSeries.category_id.name,
#     AUDIO_LANGUAGE: YoutubePostTimeSeries.audio_language.name,
#     IS_LICENSED: YoutubePostTimeSeries.is_licensed.name,
#     CHANNEL_TITLE: YoutubePostTimeSeries.channel_title.name,
#     DURATION: YoutubePostTimeSeries.duration.name,
#     CONTENT_RATING: YoutubePostTimeSeries.content_rating.name,
#     PLAYLIST_ID: YoutubePostTimeSeries.playlist_id.name,
#     PRIVACY_STATUS: YoutubePostTimeSeries.privacy_status.name
# }

profile_insights_map = {
    CITY: YoutubeProfileInsights.city.name,
    COUNTRY: YoutubeProfileInsights.country.name,
    GENDER_AGE: YoutubeProfileInsights.gender_age.name
}

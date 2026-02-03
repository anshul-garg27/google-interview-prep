from datetime import date, datetime
from typing import List

from dateutil.relativedelta import relativedelta

from loguru import logger
from core.models.models import RequestContext
from credentials.manager import CredentialManager
from utils.exceptions import PostScrapingFailed, QuotaExceeded, ProfileScrapingFailed, SubscriptionNotAllowedError, \
    NoSubscriptionsFoundError
from utils.exceptions import TokenValidationFalied, StartDateAfterEndDateError, InsightScrapingFailed
from utils.request import make_request
from utils.getter import safe_get
from youtube.functions.retriever.interface import YoutubeCrawlerInterface
from youtube.functions.retriever.ytapi.ytapi_parser import transform_channel_for_post_ids, transform_profile, \
    transform_post, \
    transform_profile_insights, transform_playlist_for_post_ids, transform_post_ids_by_genre, transform_post_comments, \
    transform_channel_ids_by_search, transform_post_log_by_search, transform_post_ids_by_search, \
    transform_activities_for_channel_id, transform_yt_profile_relationship_for_channel_id
from youtube.models.models import YoutubeProfileLog, YoutubePostLog, YoutubePostActivityLog, \
    YoutubeActivityLog, YoutubeProfileRelationshipLog

cred_provider = CredentialManager()


class YoutubeDataApi(YoutubeCrawlerInterface):

    @staticmethod
    async def fetch_profiles_by_channel_ids(ctx: RequestContext, channel_ids: list[str]) -> dict:
        key = ctx.credential.credentials["key"]
        parts = "snippet,statistics,contentDetails"
        base_uri = "https://youtube.googleapis.com/youtube/v3/channels"
        max_results = 50
        ids = ",".join(channel_ids)
        url = f"{base_uri}"
        querystring = {
            "part": parts,
            "key": key,
            "max_results": max_results,
            "id": ids
        }
        resp = await make_request("GET", url, params=querystring)

        if resp.status_code == 200 and "items" in resp.json():
            return resp.json()
        elif resp.status_code == 403:
            raise QuotaExceeded(resp.status_code, "Quota exceeded")
        else:
            raise ProfileScrapingFailed(resp.status_code, "Something went wrong")

    @staticmethod
    async def fetch_profile_insights(ctx: RequestContext) -> dict:
        refresh_token = ctx.credential.credentials['refresh_token']
        channel_id = ctx.credential.credentials['channel_id']
        client_id = ctx.credential.credentials['client_id']
        client_secret = ctx.credential.credentials['client_secret']
        base_uri = 'https://oauth2.googleapis.com/token'
        headers = {'Content-Type': 'application/x-www-form-urlencoded'}

        data = {'client_secret': client_secret,
                'grant_type': 'refresh_token',
                'refresh_token': refresh_token,
                'client_id': client_id}

        url = f"{base_uri}"
        resp = await make_request("POST", url, headers=headers, data=data)
        access_token = resp.json()['access_token']

        response = {}
        insights = [
            ('city', 'city', 'views', 50, '-views'),
            ('country', 'country', 'views', 10, '-views'),
            ('gender_age', 'gender,ageGroup', 'viewerPercentage', 100, '-viewerPercentage')
        ]
        for name, dimensions, metrics, max_results, sort in insights:
            resp = await YoutubeDataApi.fetch_profile_insights_by_dimension(
                access_token, dimensions, metrics, max_results, sort
            )
            response[name] = resp
        return response

    @staticmethod
    async def fetch_profile_insights_by_dimension(
            access_token: str,
            dimensions: str,
            metrics: str,
            max_results: int,
            sort: str
    ) -> dict:
        base_uri = 'https://youtubeanalytics.googleapis.com/v2/reports'
        ids = "channel==MINE"
        start_date = "2000-01-01"
        end_date = str(date.today())
        filters = ""
        if dimensions == "gender,ageGroup":
            filters = "isCurated==1"
        querystring = {
            'ids': ids,
            'dimensions': dimensions,
            'metrics': metrics,
            'maxResults': max_results,
            'sort': sort,
            'startDate': start_date,
            'endDate': end_date,
            'filters': filters
        }
        authorization = "Bearer " + access_token
        accept = "application/json"
        headers = {
            'Authorization': authorization,
            'Accept': accept
        }
        url = f"{base_uri}"
        resp = await make_request("GET", url, headers=headers, params=querystring)
        if resp.status_code == 200:
            return resp.json()

        if resp.status_code == 401:
            resp_dict = resp.json()
            if "error" in resp_dict and resp_dict["error"]["code"] == 401:
                logger.error(resp)
                raise TokenValidationFalied(resp.status_code, "Token is invalid")
        if resp.status_code == 412:
            if "error" in resp_dict and resp_dict["error"]["code"] == 412:
                logger.error(resp)
                raise StartDateAfterEndDateError(resp.status_code,
                                                 "start-date must be before or the same as the end-date.")
        logger.error(resp)
        raise InsightScrapingFailed(resp.status_code, "Something went wrong")

    @staticmethod
    async def fetch_posts_by_post_ids(ctx: RequestContext, post_ids: list[str]) -> dict:
        key = ctx.credential.credentials["key"]
        parts = "snippet,statistics,contentDetails"
        base_uri = "https://youtube.googleapis.com/youtube/v3/videos"
        ids = ",".join(post_ids)
        url = f"{base_uri}"
        querystring = {
            "part": parts,
            "key": key,
            "id": ids
        }
        resp = await make_request("GET", url, params=querystring)

        if resp.status_code == 200 and "items" in resp.json():
            return resp.json()
        elif resp.status_code == 403:
            raise QuotaExceeded(resp.status_code, "Quota exceeded")
        else:
            raise ProfileScrapingFailed(resp.status_code, "Something went wrong")

    @staticmethod
    async def fetch_post_ids_by_playlist_id(ctx: RequestContext, playlist_id: str, cursor=None) -> dict:
        key = ctx.credential.credentials["key"]
        parts = "snippet,contentDetails,status,id"
        max_results = 50
        base_uri = "https://www.googleapis.com/youtube/v3/playlistItems"
        url = f"{base_uri}"
        querystring = {
            "part": parts,
            "key": key,
            "playlistId": playlist_id,
            "maxResults": max_results

        }
        if cursor:
            querystring["pageToken"] = cursor
        resp = await make_request("GET", url, params=querystring)

        if resp.status_code == 200:
            resp_dict = resp.json()
            next_cursor = safe_get(resp_dict, 'nextPageToken', type=str)
            if next_cursor:
                has_next_page = True
            else:
                has_next_page = False
            return resp_dict, has_next_page, next_cursor
        elif resp.status_code == 403:
            raise QuotaExceeded(resp.status_code, "Quota exceeded")
        else:
            raise ProfileScrapingFailed(resp.status_code, "Something went wrong")

    @staticmethod
    async def fetch_post_ids_by_genre(ctx: RequestContext, category=None, language=None, cursor=None) -> (
            dict, bool, str):
        key = ctx.credential.credentials["key"]
        d = datetime.now() + relativedelta(months=-3)
        base_uri = "https://youtube.googleapis.com/youtube/v3/search"
        url = f"{base_uri}"
        querystring = {
            "key": key,
            "part": "snippet",
            "regionCode": "In",
            "order": "viewCount",
            "maxResults": 50,
            "publishedAfter": d.strftime("%Y-%m-%dT%H:%M:%SZ")
        }

        if category:
            querystring["q"] = category
        if language:
            querystring["relevanceLanguage"] = language
        if cursor:
            querystring["pageToken"] = cursor

        resp = await make_request("GET", url, params=querystring)

        if resp.status_code == 200:
            resp_dict = resp.json()
            next_cursor = safe_get(resp_dict, 'nextPageToken', type=str)
            if next_cursor:
                has_next_page = True
            else:
                has_next_page = False
            return resp_dict, has_next_page, next_cursor
        elif resp.status_code == 403:
            raise QuotaExceeded(resp.status_code, "Quota exceeded")
        else:
            raise ProfileScrapingFailed(resp.status_code, "Something went wrong")

    @staticmethod
    async def fetch_post_type(post_id: str) -> str:
        shorts_url = "https://www.youtube.com/shorts/%s" % post_id

        response = await make_request("GET", shorts_url, follow_redirects=False)

        if response.status_code == 200:
            post_type = "short"
        elif response.status_code == 303:
            post_type = "video"
        else:
            raise PostScrapingFailed(f"Something went wrong - {response.status_code}")
        return post_type

    @staticmethod
    async def fetch_post_comments(ctx: RequestContext, video_id: str, cursor: str = None) -> (dict, bool, str):
        key = ctx.credential.credentials["key"]
        parts = "snippet,id,replies"
        url = "https://www.googleapis.com/youtube/v3/commentThreads"
        querystring = {
            "part": parts,
            "videoId": video_id,
            "maxResults": "50",
            "order": "time",
            "textFormat": "html",
            "key": key
        }
        if cursor:
            querystring["pageToken"] = cursor
        response = await make_request("GET", url, params=querystring)
        if response.status_code == 200:
            data = response.json()
            next_cursor = safe_get(data, 'nextPageToken', type=str)
            if next_cursor:
                has_next_page = True
            else:
                has_next_page = False
            return data, has_next_page, next_cursor
        logger.error(f"Error getting comments for post - {video_id} - {response.json()}")
        raise ProfileScrapingFailed(response.status_code, "Something went wrong")

    @staticmethod
    async def fetch_posts_by_search(ctx: RequestContext, keyword: str, start_date: str, end_date: str, cursor=None) -> (
            dict, bool, str):
        key = ctx.credential.credentials["key"]
        url = "https://www.googleapis.com/youtube/v3/search"
        querystring = {
            "q": keyword,
            "part": "snippet",
            "regionCode": "IN",
            "type": "video",
            "order": "viewCount",
            "upload_date": "month",
            "publishedAfter": start_date,
            "publishedBefore": end_date,
            "max_results": 50,
            "key": key
        }
        if cursor:
            querystring["pageToken"] = cursor
        resp = await make_request("GET", url, params=querystring)
        if not resp.json():
            logger.debug("getting none in response")
            resp = await make_request("GET", url, params=querystring)
        if resp.status_code == 200:
            resp_dict = resp.json()
            next_cursor = safe_get(resp_dict, 'nextPageToken', type=str)
            logger.debug(next_cursor)
            if next_cursor:
                has_next_page = True
            else:
                has_next_page = False
            return resp_dict, has_next_page, next_cursor
        else:
            logger.debug(resp)
            raise ProfileScrapingFailed(resp.status_code, "Something went wrong")

    @staticmethod
    async def fetch_post_ids_by_channel_id(ctx: RequestContext, channel_id: str, filter=None, cursor=None) -> (
            dict, bool, str):
        key = ctx.credential.credentials["key"]
        base_uri = "https://www.googleapis.com/youtube/v3/search"
        url = f"{base_uri}"
        querystring = {
            "channelId": channel_id,
            "part": "id",
            "order": "date",
            "maxResults": 50,
            "key": key
        }
        if filter is not None and filter['start_date'] is not None:
            querystring['publishedAfter'] = filter['start_date'] + "T00:00:00Z"
        if filter is not None and filter['end_date'] is not None:
            querystring['publishedBefore'] = filter['end_date'] + "T23:59:59Z"

        if cursor:
            querystring["pageToken"] = cursor

        resp = await make_request("GET", url, params=querystring)

        if resp.status_code == 200:
            resp_dict = resp.json()
            next_cursor = safe_get(resp_dict, 'nextPageToken', type=str)
            if next_cursor:
                has_next_page = True
            else:
                has_next_page = False
            return resp_dict, has_next_page, next_cursor
        elif resp.status_code == 403:
            raise QuotaExceeded(resp.status_code, "Quota exceeded")
        else:
            logger.debug(resp)
            raise ProfileScrapingFailed(resp.status_code, "Something went wrong")

    @staticmethod
    async def fetch_activities_by_channel_id(ctx: RequestContext, channel_id: str, filter=None, cursor=None) -> (
            dict, bool, str):
        key = ctx.credential.credentials["key"]
        base_uri = "https://youtube.googleapis.com/youtube/v3/activities"
        url = f"{base_uri}"
        querystring = {
            "channelId": channel_id,
            "part": "contentDetails,id,snippet",
            "order": "date",
            "maxResults": 50,
            "key": key
        }

        if cursor:
            querystring["pageToken"] = cursor

        resp = await make_request("GET", url, params=querystring)

        if resp.status_code == 200:
            resp_dict = resp.json()
            next_cursor = safe_get(resp_dict, 'nextPageToken', type=str)
            if next_cursor:
                has_next_page = True
            else:
                has_next_page = False
            return resp_dict, has_next_page, next_cursor
        elif resp.status_code == 403:
            raise QuotaExceeded(resp.status_code, "Quota exceeded")
        else:
            logger.debug(resp)
            raise ProfileScrapingFailed(resp.status_code, "Something went wrong")

    @staticmethod
    async def fetch_yt_profile_relationship_by_channel_id(ctx: RequestContext, channel_id: str, filter=None, cursor=None) -> (
            dict, bool, str):
        key = ctx.credential.credentials["key"]
        base_uri = "https://youtube.googleapis.com/youtube/v3/subscriptions"
        url = f"{base_uri}"
        querystring = {
            "channelId": channel_id,
            "part": "contentDetails,id,subscriberSnippet,snippet",
            "order": "unread",
            "maxResults": 50,
            "key": key
        }

        if cursor:
            querystring["pageToken"] = cursor

        resp = await make_request("GET", url, params=querystring)
        logger.debug(resp.json())
        if resp.status_code == 200:
            resp_dict = resp.json()
            if 'pageInfo' in resp_dict and 'totalResults' in resp_dict['pageInfo'] and resp_dict['pageInfo'][
                'totalResults'] == 0:
                raise NoSubscriptionsFoundError(resp.status_code,
                                                f"channel_id {channel_id} don't have any subscriptions")
            next_cursor = safe_get(resp_dict, 'nextPageToken', type=str)
            if next_cursor:
                has_next_page = True
            else:
                has_next_page = False
            return resp_dict, has_next_page, next_cursor
        elif resp.status_code == 403:
            resp_dict = resp.json()
            if 'error' in resp_dict and 'errors' in resp_dict['error'] and 'reason' in resp_dict['error']['errors'][0] and \
                    resp_dict['error']['errors'][0]['reason'] == 'subscriptionForbidden':
                raise SubscriptionNotAllowedError(resp.status_code,
                                                  f"Not allowed to access subscription for channel_id {channel_id}")
            else:
                raise QuotaExceeded(resp.status_code, "Quota exceeded")
        else:
            logger.debug(resp)
            raise ProfileScrapingFailed(resp.status_code, "Something went wrong")

    @staticmethod
    def parse_profiles_data(resp_dict: dict) -> list[YoutubeProfileLog]:
        return [transform_profile(item) for item in resp_dict["items"]]

    @staticmethod
    def parse_posts_data(resp_dict: dict, category=None, topic=None) -> list[YoutubePostLog]:
        return [transform_post(item, category=category, topic=topic) for item in resp_dict["items"]]

    @staticmethod
    def parse_profile_insights_data(resp_dict: dict) -> YoutubeProfileLog:
        return transform_profile_insights(resp_dict)

    @staticmethod
    def parse_post_ids_by_playlist_id(resp_dict: dict) -> list[str]:
        return [transform_playlist_for_post_ids(item) for item in resp_dict["items"]]

    @staticmethod
    def parse_post_ids_by_genre(resp_dict: dict) -> list[str]:
        return [transform_post_ids_by_genre(item) for item in resp_dict["items"]]

    @staticmethod
    def parse_post_comments(shortcode: str, resp_dict: dict) -> List[YoutubePostActivityLog]:
        resp_data = resp_dict["items"]
        return [transform_post_comments(shortcode, data) for data in resp_data]

    @staticmethod
    def parse_post_ids_by_search(resp_dict: dict) -> tuple[list[str], list[str], list[YoutubePostLog]]:
        resp_data = resp_dict["items"]
        return [transform_post_ids_by_search(data) for data in resp_data], [transform_channel_ids_by_search(data) for
                                                                            data in resp_data], [
            transform_post_log_by_search(data) for data in resp_data]

    @staticmethod
    def parse_activities_by_channel_id(resp_dict: dict) -> list[YoutubeActivityLog]:
        return [transform_activities_for_channel_id(item) for item in resp_dict["items"]]

    @staticmethod
    def parse_yt_profile_relationship_by_channel_id(resp_dict: dict) -> list[YoutubeProfileRelationshipLog]:
        return [transform_yt_profile_relationship_for_channel_id(item) for item in resp_dict["items"]]

    @staticmethod
    def parse_post_ids_by_channel_id(resp_dict: dict, filter=None) -> list[str]:
        post_ids = []
        for item in resp_dict["items"]:
            try:
                post_ids.append(transform_channel_for_post_ids(item))
            except PostScrapingFailed as e:
                logger.error(e)
        return post_ids

import asyncio
from datetime import datetime, timedelta

from typing import List

from loguru import logger

from core.models.models import RequestContext
from credentials.manager import CredentialManager
from instagram.functions.retriever.graphapi.graphapi_parser import transform, transform_post, transform_post_insights, \
    transform_profile_insights, transform_story_insights
from instagram.functions.retriever.interface import InstagramCrawlerInterface
from instagram.models.models import InstagramPostLog, InstagramProfileLog
from utils.exceptions import ProfileScrapingFailed, InsightScrapingFailed, ProfileMissingError, NotSupportedError, \
    ProfileInsightsMissingError, TokenValidationFalied, DataAccessExpired
from utils.request import make_request

cred_provider = CredentialManager()


class GraphApi(InstagramCrawlerInterface):

    @staticmethod
    async def is_token_valid(ctx: RequestContext, handle: str) -> bool:
        try:
            _, profile_id = await GraphApi.fetch_profile_by_handle(ctx, handle)
        except ProfileMissingError:
            return True
        except ProfileScrapingFailed:
            return False
        else:
            if profile_id:
                return True
            return False
    
    @staticmethod
    async def debug_token(token: str) -> dict:
        resp = {}
        url = "https://graph.facebook.com/v15.0/debug_token"
        querystring = {"input_token": token, "access_token": "308225188075192|2HyQEub5Ma4ltNvxpE7jmMRn0zA"}
        response = await make_request("GET", url, params=querystring)
        resp['data_access_expired'] = False
        resp['scopes'] = []
        if response.status_code == 200:
            resp_dict = response.json()
            if "data" in resp_dict and "data_access_expires_at" in resp_dict["data"]:
                if datetime.fromtimestamp(resp_dict["data"]["data_access_expires_at"]) <= datetime.now() - timedelta(days=7):
                    resp['data_access_expired'] = True
            scopes = set(resp_dict["data"]["scopes"])
            if "data" in resp_dict and "scopes" in resp_dict["data"]:
                resp['scopes'] = scopes
        return resp

    @staticmethod
    async def fetch_profile_by_handle(ctx: RequestContext, handle: str) -> (dict, str):
        token = ctx.credential.credentials['token']
        user_id = ctx.credential.credentials['user_id']
        fields = "business_discovery.username(%s){id,ig_id,name,website,followers_count,follows_count,profile_picture_url,username,biography,media_count,media{id,user_id,username,comments_count,like_count,media_type,media_url,caption,timestamp,permalink}}" % handle
        url = "https://graph.facebook.com/v12.0/%s" % user_id
        querystring = {"access_token": token, "fields": fields}
        response = await make_request("GET", url, params=querystring)
        resp_dict = response.json()
        if response.status_code == 200:
            if 'business_discovery' not in resp_dict or 'ig_id' not in resp_dict['business_discovery']:
                raise ProfileScrapingFailed(response.status_code, "GraphAPI - business discovery or ig_id not available")
            profile_id = str(resp_dict['business_discovery']['ig_id'])
            if profile_id == "":
                raise ProfileScrapingFailed(response.status_code, "GraphAPI - profile_id is empty")
            return resp_dict, profile_id
        if response.status_code == 400:
            if "error" in resp_dict and resp_dict["error"]["code"] == 110:
                logger.error(response)
                raise ProfileMissingError(response.status_code, "Invalid username")
            elif "error" in resp_dict and resp_dict["error"]["code"] == 10:
                logger.error(response)
                raise ProfileScrapingFailed(response.status_code, resp_dict["error"]["message"])
            else:
                logger.error(response)
                raise ProfileScrapingFailed(response.status_code, resp_dict["error"]["message"])
        if response.status_code == 403:
            logger.error(response)
            raise ProfileScrapingFailed(response.status_code, "Forbidden exception")
        logger.error(response)
        raise ProfileScrapingFailed(response.status_code, resp_dict["error"]["message"])

    @staticmethod
    async def fetch_profile_posts_by_handle(ctx: RequestContext, handle: str, cursor: str = None) -> (
            List[InstagramPostLog], bool, str):
        token = ctx.credential.credentials['token']
        user_id = ctx.credential.credentials['user_id']
        fields = "business_discovery.username(%s){media_count,media{id,ig_id,name,comments_count,like_count,media_type,media_url,caption,timestamp,permalink}}" % handle
        url = "https://graph.facebook.com/v12.0/%s" % user_id
        querystring = {"access_token": token, "fields": fields}
        if cursor:
            querystring["after"] = cursor
        response = await make_request("GET", url, params=querystring)
        if response.status_code == 200:
            resp_dict = response.json()
            page_info = resp_dict['business_discovery']['media']['paging']
            has_next_page = True if bool(page_info['cursors']['after']) else False
            return resp_dict, has_next_page, page_info['cursors']['after']
        if response.status_code == 400:
            resp_dict = response.json()
            if "error" in resp_dict and resp_dict["error"]["code"] == 110:
                logger.error(response)
                raise ProfileMissingError(response.status_code, "Invalid username")
            elif "error" in resp_dict and resp_dict["error"]["code"] == 10:
                logger.error(response)
                raise ProfileScrapingFailed(response.status_code, "Something went wrong")
            else:
                logger.error(response)
                raise ProfileScrapingFailed(response.status_code, "Something went wrong")
        if response.status_code == 403:
            logger.error(response)
            raise ProfileScrapingFailed(response.status_code, "Forbidden exception")
        logger.error(response)
        raise ProfileScrapingFailed(response.status_code, "Something went wrong")

    @staticmethod
    async def fetch_stories_posts(ctx: RequestContext, handle: str, cursor: str = None) -> (dict, str):
        token = ctx.credential.credentials['token']
        user_id = ctx.credential.credentials['user_id']
        fields = "business_discovery.username(%s){ig_id,username,media_count,media{id,username,name,comments_count,like_count,media_type,media_product_type,media_url,caption,timestamp,permalink},stories{id,media_type,media_product_type,media_url,owner,permalink,timestamp,username,caption}}" % handle
        url = "https://graph.facebook.com/v12.0/%s" % user_id
        querystring = {"access_token": token, "fields": fields}
        if cursor:
            fields = "business_discovery.username(%s){ig_id,username,media_count,media.after(%s){id,username,name,comments_count,like_count,media_type,media_product_type,media_url,caption,timestamp,permalink},stories{id,media_type,media_product_type,media_url,owner,permalink,timestamp,username,caption}}" % (handle, cursor)
            querystring = {"access_token": token, "fields": fields}
        response = await make_request("GET", url, params=querystring)
        if response.status_code == 200:
            resp_dict = response.json()
            if 'paging' not in resp_dict['business_discovery']['media']:
                return resp_dict, False, None
            post_page_info = resp_dict['business_discovery']['media']['paging']
            if 'after' not in post_page_info['cursors']:
                return resp_dict, False, None
            has_next_page = True if bool(post_page_info['cursors']['after']) else False
            return resp_dict, has_next_page, post_page_info['cursors']['after']

        if response.status_code == 400:
            resp_dict = response.json()
            if "error" in resp_dict and resp_dict["error"]["code"] == 110:
                logger.error(response)
                raise ProfileMissingError(response.status_code, "Invalid username")
            elif "error" in resp_dict and resp_dict["error"]["code"] == 100:
                logger.error(response)
                raise ProfileScrapingFailed(response.status_code,
                                            "Unsupported get request. Object with ID does not exist")
            elif "error" in resp_dict and resp_dict["error"]["code"] == 190:
                logger.error(response)
                raise ProfileScrapingFailed(response.status_code, "Malformed access token.")
            elif "error" in resp_dict and resp_dict["error"]["code"] == 10:
                logger.error(response)
                raise ProfileScrapingFailed(response.status_code, "Application does not have permission for this action")
            else:
                logger.error(response)
                raise ProfileScrapingFailed(response.status_code, "Something went wrong")
        if response.status_code == 403:
            logger.error(response)
            raise ProfileScrapingFailed(response.status_code, "Forbidden exception")
        logger.error(response)
        raise ProfileScrapingFailed(response.status_code, "Something went wrong")

    @staticmethod
    async def fetch_post_insights(ctx: RequestContext, handle: str, post_id: str, post_type: str) -> dict:
        metric = None
        token = ctx.credential.credentials['token']
        period = "lifetime"
        if post_type == "image":
            metric = "likes,reach,impressions,comments,engagement,saved,shares"
        elif post_type == "reels":
            metric = "plays,likes,comments,reach,total_interactions,saved,shares"
        elif post_type == "carousel":
            metric = "reach,impressions,engagement,saved,video_views"
        if '_' in str(post_id):
            post_id = str(post_id).split('_')[0]
        url = "https://graph.facebook.com/v12.0/%s/insights" % post_id
        querystring = {"access_token": token, "period": period, "metric": metric}
        response = await make_request("GET", url, params=querystring)
        if response.status_code == 200:
            resp_dict = response.json()
            return resp_dict
        if response.status_code == 400:
            resp_dict = response.json()
            if resp_dict["error"]["code"] == 100 and "error_subcode" in resp_dict["error"] and resp_dict["error"]["error_subcode"] == 33:
                logger.error(response)
                token_debug_info = await GraphApi.debug_token(token)
                scopes = token_debug_info['scopes']
                scopes_needed = set(
                    ["instagram_basic",
                    "instagram_manage_insights",
                    "pages_read_engagement",
                    "pages_show_list"
                    ]
                )
                if len(scopes_needed - scopes) > 0:
                    raise TokenValidationFalied(401, "Token is invalid")
                elif token_debug_info['data_access_expired']:
                    raise DataAccessExpired(401, "Data access expired")
                else:
                    raise InsightScrapingFailed(404, "Post missing")
            elif resp_dict["error"]["code"] == 100 and "error_subcode" in resp_dict["error"] and resp_dict["error"]["error_subcode"] == 2108006:
                logger.error(response)
                raise InsightScrapingFailed(response.status_code,
                                            "Media was posted before account was converted to business account")
            elif resp_dict["error"]["code"] == 190:
                logger.error(response)
                raise TokenValidationFalied(401, "Invalid access token")
        if response.status_code == 403:
            logger.error(response)
            raise InsightScrapingFailed(response.status_code, "Forbidden exception")
        logger.error(response)
        raise InsightScrapingFailed(response.status_code, "Something went wrong")


    @staticmethod
    async def get_insights(token: str, user_id: str, metric: str, breakdown: str):
        period = "lifetime"
        timeframe = "last_90_days"
        metric_type = "total_value"
        url = "https://graph.facebook.com/v17.0/%s/insights" % user_id
        query_string = {"access_token": token, "period": period, "metric": metric, "metric_type": metric_type,
                        "timeframe": timeframe, "breakdown": breakdown}
        response = await make_request("GET", url, params=query_string)
        return response

    @staticmethod
    async def fetch_profile_insights(ctx: RequestContext, handle: str, api_call: bool) -> dict:
        token = ctx.credential.credentials['token']
        user_id = ctx.credential.credentials['user_id']
        
        if api_call:
            metrics = ["follower_demographics"]
            breakdowns = ["city", "gender, age"]
        else:
            metrics = ["follower_demographics", "reached_audience_demographics", "engaged_audience_demographics"]
            breakdowns = ["country", "city", "gender, age"]
    
        tasks = [asyncio.create_task(GraphApi.get_insights(token, user_id, metric, breakdown)) for metric in metrics for breakdown in breakdowns]
        responses = await asyncio.gather(*tasks)
        
        final_response = {}
        response_error = False
        response_follower_demographics = ""
        for response in responses:
            response_follower_demographics = response
            if response.status_code == 200:
                response_error = True
            data = response.json()
            if 'data' in data and data['data']:
                element = data['data'][0]
                if 'name' in element and 'total_value' in element:
                    name = element['name']
                    total_value = element['total_value']
                    if 'breakdowns' in total_value and total_value['breakdowns']:
                        breakdown = total_value['breakdowns'][0]
                        if name in final_response:
                            final_response[name].append(breakdown)
                        else:
                            final_response[name] = [breakdown]
                    
        if response_error:
            return final_response
        else:
            response = response_follower_demographics
            if response.status_code == 400:
                resp_dict = response.json()
                if "error" in resp_dict and resp_dict["error"]["code"] == 110:
                    raise ProfileMissingError(response.status_code, "Invalid username")
                elif resp_dict["error"]["code"] == 100 and "error_subcode" in resp_dict["error"] and resp_dict["error"]["error_subcode"] == 33:
                    raise InsightScrapingFailed(response.status_code, "User ID doesn't exist")
                elif resp_dict["error"]["code"] == 190:
                    raise TokenValidationFalied(401, "Invalid access token")
                else:
                    raise InsightScrapingFailed(response.status_code, "Something went wrong")
            if response.status_code == 403:
                raise InsightScrapingFailed(response.status_code, "Forbidden exception")
            logger.error(response)
            raise InsightScrapingFailed(response.status_code, "Something went wrong")
            
    @staticmethod
    async def fetch_story_insights(ctx: RequestContext, handle: str, story_id: str) -> dict:
        token = ctx.credential.credentials['token']
        period = "lifetime"
        metric = "impressions,reach,replies,saved,shares,total_interactions,follows,profile_visits,profile_activity"
        url = "https://graph.facebook.com/v15.0/%s/insights" % story_id
        query_string = {"access_token": token, "period": period, "metric": metric}
        response = await make_request("GET", url, params=query_string)
        query_string["metric"] = "navigation"
        query_string["breakdown"] = "story_navigation_action_type"
        response_navigation_breakdown = await make_request("GET", url, params=query_string)
        if response.status_code == 200:
            if not response.json()["data"]:
                logger.error(response)
                raise ProfileInsightsMissingError(response.status_code, "Profile Insights Missing")
            else:
                resp_dict = response.json()
                resp_dict_navigation_breakdown = response_navigation_breakdown.json()["data"][0]
                resp_dict['data'].append(resp_dict_navigation_breakdown)
                return resp_dict
        if response.status_code == 400:
            resp_dict = response.json()
            if resp_dict["error"]["code"] == 100 and "error_subcode" in resp_dict["error"] and resp_dict["error"]["error_subcode"] == 33:
                logger.error(response)
                raise InsightScrapingFailed(response.status_code, "Story ID doesn't exist")
            elif resp_dict["error"]["code"] == 190:
                logger.error(response)
                raise InsightScrapingFailed(response.status_code, "Invalid access token")
            elif resp_dict["error"]["code"] == 10:
                raise InsightScrapingFailed(response.status_code, "Not enough viewers for the media to show insights")
            else:
                logger.error(response)
                raise InsightScrapingFailed(response.status_code, "Something went wrong")
        if response.status_code == 403:
            logger.error(response)
            raise InsightScrapingFailed(response.status_code, "Forbidden exception")
        logger.error(response)
        raise InsightScrapingFailed(response.status_code, "Something went wrong")

    @staticmethod
    async def fetch_profile_by_id(ctx: RequestContext, profile_id: str) -> dict:
        raise NotSupportedError("Not Supported")

    @staticmethod
    async def fetch_profile_posts_by_id(ctx: RequestContext, profile_id: str, cursor: str = None) -> (
            List[InstagramPostLog], bool, str):
        raise NotSupportedError("Not Supported")

    @staticmethod
    async def fetch_post_by_shortcode(ctx: RequestContext, shortcode: str) -> dict:
        raise NotSupportedError("Not Supported")

    @staticmethod
    async def fetch_reels_posts(ctx: RequestContext, profile_id: str) -> (List[InstagramPostLog], bool, str):
        raise NotSupportedError("Not Supported")

    @staticmethod
    def parse_profile_data(resp_dict: dict) -> InstagramProfileLog:
        return transform(resp_dict)

    @staticmethod
    def parse_profile_post_data(posts: list[dict]) -> List[InstagramPostLog]:
        posts = [transform_post(p) for p in posts]
        return posts

    @staticmethod
    def parse_profiles_post_data(posts: list[dict]) -> List[InstagramPostLog]:
        posts = [transform_post(p) for p in posts]
        return posts

    @staticmethod
    def parse_post_insights_data(shortcode: str, resp_dict: dict) -> InstagramPostLog:
        insights = transform_post_insights(shortcode, resp_dict)
        return insights

    @staticmethod
    def parse_profile_insights_data(resp_dict: dict) -> InstagramProfileLog:
        insights = transform_profile_insights(resp_dict)
        return insights

    @staticmethod
    def parse_story_insights_data(story_id: str, resp_dict: dict) -> InstagramPostLog:
        insights = transform_story_insights(story_id, resp_dict)
        return insights


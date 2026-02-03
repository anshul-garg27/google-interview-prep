from loguru import logger

from instagram.functions.retriever.interface import InstagramCrawlerInterface
from instagram.functions.retriever.rapidapi.jotucker.jotucker import JoTucker


def test_crawler(crawler: InstagramCrawlerInterface) -> None:
    profile_by_id_resp = crawler.fetch_profile_by_id("2947839736")
    profile_by_handle_resp = crawler.fetch_profile_by_handle("winklapp")
    profile_posts_resp = crawler.fetch_profile_posts_by_id("2947839736")
    post_resp = crawler.fetch_post_by_shortcode("CmdjXcCPl3Y")
    try:
        parsed_profile_by_id_resp = crawler.parse_profile_data(profile_by_id_resp)
        logger.debug(parsed_profile_by_id_resp)
    except Exception as e:
        logger.error("Unable to parse profile by id resp - %s", e)
    try:
        parsed_profile_by_handle_resp = crawler.parse_profile_data(profile_by_handle_resp[0])
        logger.debug(parsed_profile_by_handle_resp)
    except Exception as e:
        logger.error("Unable to parse profile by handle resp - %s", e)
    try:
        parsed_profile_posts_resp = crawler.parse_profile_post_data(profile_posts_resp[0])
        logger.debug(parsed_profile_posts_resp)
    except Exception as e:
        logger.error("Unable to parse profile posts resp - %s", e)
    try:
        parsed_post_resp = crawler.parse_post_by_shortcode(post_resp)
        logger.debug(parsed_post_resp)
    except Exception as e:
        logger.error("Unable to parse profile posts resp - %s", e)


if __name__ == "__main__":
    test_crawler(JoTucker)

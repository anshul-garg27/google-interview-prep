from instagram.models.models import InstagramProfileLog


class GptCrawlerInterface:

    @staticmethod
    async def fetch_instagram_gpt_data_base_gender(name: str, handle: str, bio: str) -> dict:
        pass
    
    @staticmethod
    async def fetch_instagram_gpt_data_base_location(name: str, handle: str, bio: str) -> dict:
        pass
    
    @staticmethod
    async def fetch_instagram_gpt_data_base_categ_lang_topics(name: str, handle: str, bio: str) -> dict:
        pass

    @staticmethod
    async def fetch_instagram_gpt_data_audience_age_gender(name: str, handle: str, bio: str) -> dict:
        pass

    @staticmethod
    async def fetch_instagram_gpt_data_audience_cities(name: str, handle: str, bio: str) -> dict:
        pass

    @staticmethod
    def parse_instagram_gpt_data_base_gender(resp_dict: dict) -> InstagramProfileLog:
        pass
    
    @staticmethod
    def parse_instagram_gpt_data_base_location(resp_dict: dict) -> InstagramProfileLog:
        pass
    
    @staticmethod
    def parse_instagram_gpt_data_base_categ_lang_topics(resp_dict: dict) -> InstagramProfileLog:
        pass
    
    @staticmethod
    def parse_instagram_gpt_data_audience_age_gender(resp_dict: dict) -> InstagramProfileLog:
        pass
    
    @staticmethod
    def parse_instagram_gpt_data_audience_cities(resp_dict: dict) -> InstagramProfileLog:
        pass

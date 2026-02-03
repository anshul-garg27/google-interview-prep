import random
from typing import Optional

from loguru import logger

from core.helpers.session import sessionize
from core.models.models import PaginationContext
from credentials.manager import CredentialManager
from credentials.models import CredentialModel


class Crawler:
        
    @sessionize
    async def fetch_enabled_sources(self, fetch_type: str, exclusions: list = None, session=None, pagination: PaginationContext = None) -> Optional[
        CredentialModel]:
        if not exclusions:
            exclusions = []
        if pagination and pagination.source is not None:
            source_cred = await CredentialManager.get_enabled_cred(str(pagination.source), session=session)
            if source_cred:
                return source_cred
        else:
            if fetch_type in ["fetch_followers"]:
                weights = {"rapidapi-jotucker": 0.3, "rocketapi": 0.7}
                source = random.choices(list(weights.keys()), weights=list(weights.values()), k=1)[0]
                source_cred = await CredentialManager.get_enabled_cred(source, session=session)
                if source_cred:
                    return source_cred
            else:
                for source in self.available_sources[fetch_type]:
                    if source in exclusions:
                        continue
                    source_cred = await CredentialManager.get_enabled_cred(source, session=session)
                    if source_cred:
                        return source_cred
        return None

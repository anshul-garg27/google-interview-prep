import os
from typing import Optional
from dataclasses import dataclass
from utils.request import make_request

@dataclass
class IdentitySocialAccount:
    socialUserId: str
    graphAccessToken: str
    handle: str
    platform: str
    accessTokenExpired: bool


class IdentityClient:

    @staticmethod
    async def get_handle_details(handle:str) -> Optional[IdentitySocialAccount]:
        params = ["socialUserId", "graphAccessToken", "handle", "platform"]
        base_url = os.environ["IDENTITY_URL"]
        url= base_url + "/identity-service/api/socialaccount/platform/INSTAGRAM/handle/%s" % handle
        headers = {"Authorization": "Basic YTpi"}
        response = await make_request("GET", url, headers=headers)
        if response.status_code == 200:
            resp_dict = response.json()
            if "status" in resp_dict and "key" in resp_dict["status"] and resp_dict["status"]["key"] == "SUCCESS":
                socialAccount = resp_dict["socialAccounts"][0]
                if socialAccount and set(params).issubset(set(list((socialAccount.keys())))):
                    identity = IdentitySocialAccount(
                        socialUserId = socialAccount["socialUserId"],
                        graphAccessToken = socialAccount["graphAccessToken"],
                        handle = socialAccount["handle"],
                        platform = socialAccount["platform"],
                        accessTokenExpired = False
                    )
                    return identity
                return None
            return None
        return None

import re
from utils.request import make_request


async def fetch_yt_channel_id_from_handle(handle:str):
    if handle.find('@') == -1:
        handle = '@' + handle

    url = 'https://www.youtube.com/' + handle
    response = await make_request("GET", url)

    if response.status_code == 200:
        found = re.findall('<meta itemprop="channelId" content="([^"]*)"', response.text)
        if not found:
            found = re.findall('<meta itemprop="identifier" content="([^"]*)"', response.text)
        return found[0]
    else:
        return False



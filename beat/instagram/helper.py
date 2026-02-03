import math

from instagram.entities.entities import InstagramPost, InstagramAccount


def get_engagement_rate(entity: InstagramPost, account: InstagramAccount):
    followers = account.followers or 0
    if followers == 0:
        return 0
    likes = entity.likes or 0
    comments = entity.comments or 0
    return round(((likes + comments) / followers) * 100, 2)


def get_reach(entity: InstagramPost, account: InstagramAccount):
    plays = entity.plays or 0
    likes = entity.likes or 0
    followers = account.followers or 0

    reach = entity.reach

    if ((plays == 0 or followers == 0) and entity.post_type == "reels") or (likes == 0 and entity.post_type in ("image", "carousel")):
        return 0

    if reach is None or reach == 0:
        if entity.post_type == 'reels':
            reach = int(plays * (0.94 - (math.log2(followers) * 0.001)))
        else:
            reach = int((7.6 - (math.log10(likes) * 0.7)) * 0.85 * likes)
    return reach

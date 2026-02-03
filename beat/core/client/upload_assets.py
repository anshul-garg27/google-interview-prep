import os

from datetime import datetime
import aioboto3
import asks

from core.entities.entities import AssetLog
from core.helpers.session import sessionize
from utils.db import get_or_create
from utils.exceptions import UrlSignatureExpiredError


async def put(entity_id: str, entity_type: str, platform: str, profile_url: str, extension: str) -> str:
    response = await asks.get(profile_url)
    if response.text == "URL signature expired":
        raise UrlSignatureExpiredError("URL signature expired")

    region_name = 'ap-south-1'
    session = aioboto3.Session(
        aws_access_key_id='AKIAXGXUCIERUOPVUUPW',
        aws_secret_access_key='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',
        region_name=region_name
    )
    content_type = ""
    if extension == "jpg" or extension == "png"  or extension == "webp":
        content_type = "image/" + extension
    elif extension == "mp4":
        content_type = "video/mp4"

    blob_s3_key = f"assets/{entity_type.lower()}s/{platform.lower()}/{entity_id}.{extension}"
    async with session.client('s3') as s3:
        bucket = os.environ["AWS_BUCKET"]
        resp = await s3.put_object(Body=response.content, Bucket=bucket, Key=blob_s3_key,
                                   ContentType=content_type,
                                   ContentDisposition='inline')
        await s3.close()
        s3_object_url = f"{blob_s3_key}"
        return s3_object_url


@sessionize
async def upsert_save_profile_url(query_string: dict, session=None) -> None:
    instance: AssetLog = (await get_or_create(session, AssetLog, entity_id=query_string['entity_id']))[0]

    for k, v in query_string.items():
        setattr(instance, k, v)


async def asset_upload_flow(entity_type: str, entity_id: str, asset_type: str, original_url: str, platform: str,
                            hash: str, session=None):
    extension = ""
    if asset_type.upper() == "IMAGE":
        if "jpg" in original_url or "jpeg" in original_url:
            extension = "jpg"
        elif "webp" in original_url:
            extension = "webp"
        else:
            extension = "png"
    elif asset_type.upper() == "VIDEO":
        extension = "mp4"

    asset_url = await put(entity_id, entity_type, platform, original_url, extension)
    query_string = {
        'entity_type': entity_type,
        'entity_id': entity_id,
        'platform': platform,
        'asset_type': asset_type,
        'asset_url': asset_url,
        'original_url': original_url,
        'updated_at': datetime.now()
    }
    await upsert_save_profile_url(query_string, session=session)
    
    
async def asset_upload_flow_stories(entity_type: str, entity_id: str, asset_type: str, original_url: str, platform: str,
                            hash: str, session=None):
    await asset_upload_flow(entity_type, entity_id, asset_type, original_url, platform, hash, session=session)

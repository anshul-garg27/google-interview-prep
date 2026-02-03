import os

import asyncio
import time

from core.amqp.amqp import publish
from keyword_collection.generate_instagram_report import intermediate_report_generation, categorize_data, \
    final_report_generation as instagram_final_report_generation, demographic_report_generation as instagram_demographic_report_generation
from keyword_collection.generate_youtube_report import fetch_data, final_report_generation as youtube_final_report_generation, \
    demographic_report_generation as youtube_demographic_report_generation
from keyword_collection.keyword_collection_logs import keyword_collection_logs


async def generate_report(payload, session=None):
    platform = payload['platform']
    job_id = payload['job_id']
    if platform == 'YOUTUBE':
        start = time.perf_counter()
        await fetch_data(payload, session=session)
        await asyncio.sleep(60)
        parquet_file_path = await youtube_final_report_generation(payload)
        keyword_collection_logs(parquet_file_path, job_id)
        demographic_parquet_file_path = await youtube_demographic_report_generation(payload, parquet_file_path)
        payload = {
            "job_id": job_id,
            "report_asset": {
                "path": parquet_file_path,
                "bucket": os.environ["AWS_BUCKET"]
            },
            "demographic_report_asset": {
                "path": demographic_parquet_file_path,
                "bucket": os.environ["AWS_BUCKET"]
            }
        }
        publish(payload, "coffee.dx", "keyword_collection_report_completion")
        stop = time.perf_counter()
        message = f"The total time taken to search for YOUTUBE with job Id as {job_id} :- {stop - start}"
        keyword_collection_logs(message, job_id)

    elif platform == 'INSTAGRAM':
        start = time.perf_counter()
        parquet_file_path = await intermediate_report_generation(payload)
        keyword_collection_logs(parquet_file_path, job_id)
        await categorize_data(parquet_file_path, job_id)
        await asyncio.sleep(60)
        await instagram_final_report_generation(parquet_file_path)
        demographic_parquet_file_path = await instagram_demographic_report_generation(payload, parquet_file_path)
        keyword_collection_logs(demographic_parquet_file_path, job_id)
        payload = {
            "job_id": job_id,
            "report_asset": {
                "path": parquet_file_path,
                "bucket": os.environ["AWS_BUCKET"]
            },
            "demographic_report_asset": {
                "path": demographic_parquet_file_path,
                "bucket": os.environ["AWS_BUCKET"]
            }
        }
        publish(payload, "coffee.dx", "keyword_collection_report_completion")
        stop = time.perf_counter()
        message = f"The total time taken to search for INSTAGRAM with job Id as {job_id } :- {stop-start}"
        keyword_collection_logs(message, job_id)


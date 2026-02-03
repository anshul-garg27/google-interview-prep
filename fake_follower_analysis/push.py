import itertools
import json
import multiprocessing
from datetime import datetime
import time
import ijson

import boto3
import clickhouse_connect

# Create high-level resources for AWS services
sqs = boto3.resource('sqs', region_name='ap-south-1')
s3 = boto3.resource('s3', aws_access_key_id='AKIAXGXUCIERSOUELY73',
                    aws_secret_access_key='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')
kinesis_client = boto3.client('kinesis', region_name='ap-south-1')

# Create a connection to the ClickHouse database
client = clickhouse_connect.get_client(host="ec2-52-66-200-31.ap-south-1.compute.amazonaws.com",
                                  user="cube",
                                  password="CuB31Z@3995ddOpT")

current_date = datetime.today()
file_name = f"{current_date.year}_{current_date.month}_{current_date.day}_creator_followers.json"

query = f"""
            insert into function s3('https://gcc-social-data.s3.ap-south-1.amazonaws.com/temp/{file_name}', 
            'AKIAXGXUCIERSOUELY73', 'dHb+F8Nj9o0imlDFq0gq/pCnWzSL/F5kkrywtSxO', 'JSONEachRow')
            with 
                handles as (select *
                            from s3('https://gcc-social-data.s3.ap-south-1.amazonaws.com/temp/creators_handles.csv',
                                    'AKIAXGXUCIERSOUELY73', 'dHb+F8Nj9o0imlDFq0gq/pCnWzSL/F5kkrywtSxO', 'CSVWithNames')),
                profile_ids as (select ig_id from dbt.mart_instagram_account mia where handle in handles),
                follower_data as (select log.target_profile_id as                          target_profile_id,
                                        JSONExtractString(source_dimensions, 'handle')    follower_handle,
                                        JSONExtractString(source_dimensions, 'full_name') follower_full_name
                                from dbt.stg_beat_profile_relationship_log log
                                where target_profile_id in profile_ids
                                    and follower_handle is not null
                                    and follower_full_name is not null
                                    and follower_handle != ''
                                    and follower_full_name != ''),
                follower_events_data as (select log.target_profile_id as                          target_profile_id,
                                                JSONExtractString(source_dimensions, 'handle')    follower_handle,
                                                JSONExtractString(source_dimensions, 'full_name') follower_full_name
                                        from _e.profile_relationship_log_events log
                                        where target_profile_id in profile_ids
                                            and follower_handle is not null
                                            and follower_full_name is not null
                                            and follower_handle != ''
                                            and follower_full_name != ''),
                data as (select *
                        from follower_data
                        UNION ALL
                        select *
                        from follower_events_data)

            select mia.handle handle, follower_handle, follower_full_name
            from data
                    inner join dbt.mart_instagram_account mia on mia.ig_id = target_profile_id
            group by handle, follower_handle, follower_full_name
            SETTINGS s3_truncate_on_insert=1
            """

result = client.query(query)

# # Download the file from S3
s3.Bucket('gcc-social-data').download_file(f"temp/{file_name}", file_name)

# # Create the queue and get its URL
queue = sqs.create_queue(QueueName='creator_follower_in',
                         Attributes={
                             # Specify the maximum message size in bytes
                             'MaximumMessageSize': '262144',
                             # Specify the message retention period in seconds
                             'MessageRetentionPeriod': '345600',
                             # Specify the visibility timeout in seconds
                             'VisibilityTimeout': '30'
                         })
queue_url = queue.url
print('Created queue:', queue_url)

# Create the stream and get its status
try:
    stream = kinesis_client.create_stream(StreamName='creator_out',
                                   StreamModeDetails={'StreamMode': 'ON_DEMAND'})
    print('Created stream: creator_out')
except Exception as e:
    if e.__class__.__name__ == 'ResourceInUseException':
        print("Stream already created")
        print(e)
    else:
        exit
        raise e

def final(messages):
    # print(messages)
    # Send data to SQS queue in batches
    response = queue.send_messages(Entries=messages)
    # Check for any failed messages
    failed = response.get('Failed')
    if failed:
        print(failed)


def divide_into_buckets(lst, num_buckets):
    buckets = [[] for _ in range(num_buckets)]
    for index, item in enumerate(lst):
        bucket_index = index % num_buckets
        buckets[bucket_index].append(item)
    return buckets


i = 0
start_time = time.perf_counter()
# Load the data from the file as a list of JSON objects
data = []
with open(file_name, 'r') as f:
    while True:
        batch_lines = list(itertools.islice(f, 10000))
        if not batch_lines:
            break
        messages = [{
                'Id': str(i),
                'MessageBody': row
            } for i, row in enumerate(batch_lines)]

        # Divide the messages into buckets
        # print(messages)
        num_buckets = 8
        buckets = divide_into_buckets(messages, num_buckets)

        # Create a pool of workers
        pool = multiprocessing.Pool(processes=num_buckets)

        # Process the buckets in parallel
        pool.map(final, buckets)

        # Close the pool
        pool.close()
        pool.join()

end_time = time.perf_counter()
print(f"total time :- {end_time - start_time}")

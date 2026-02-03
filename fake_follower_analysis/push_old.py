import itertools
import multiprocessing
from datetime import datetime
import time

import boto3
import clickhouse_connect

client = clickhouse_connect.get_client(host="ec2-52-66-200-31.ap-south-1.compute.amazonaws.com",
                                       username="cube",
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

session = boto3.Session(
    aws_access_key_id='AKIAXGXUCIERSOUELY73',
    aws_secret_access_key='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',
)

sqs = boto3.client('sqs', region_name='ap-south-1')

s3_client = boto3.client('s3', aws_access_key_id='AKIAXGXUCIERSOUELY73',
                         aws_secret_access_key='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')

kinesis_client = boto3.client('kinesis', region_name='ap-south-1')

s3_client.download_file('gcc-social-data', f"temp/{file_name}", file_name)

# Specify the queue name and attributes
queue_name = 'creator_follower_in'
attributes = {
    # Specify the maximum message size in bytes
    'MaximumMessageSize': '262144',
    # Specify the message retention period in seconds
    'MessageRetentionPeriod': '345600',
    # Specify the visibility timeout in seconds
    'VisibilityTimeout': '30'
}
# Create the queue and get its URL
response = sqs.create_queue(QueueName=queue_name, Attributes=attributes)
queue_url = response['QueueUrl']
print('Created queue:', queue_url)

# Specify the stream name and the number of shards
stream_name = 'creator_out'
# Create the stream and get its status
try:
    response = kinesis_client.create_stream(StreamName=stream_name,
                                            StreamModeDetails={'StreamMode': 'ON_DEMAND'})
    print('Created stream:', stream_name)
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
    for i in range(0, len(messages), 10):
        # print(messages[i:i+10])
        # print(messages[0])
        try:
            response = sqs.send_message_batch(
                QueueUrl=queue_url,
                Entries=messages[i:i + 10]
            )
        except Exception as e:
            print(e)
            print(messages)


def divide_into_buckets(lst, num_buckets):
    buckets = [[] for _ in range(num_buckets)]
    for index, item in enumerate(lst):
        bucket_index = index % num_buckets
        buckets[bucket_index].append(item)
    return buckets


i = 0
start_time = time.perf_counter()
with open(file_name, 'r') as f:
    while True:
        batch_lines = list(itertools.islice(f, 10000))
        if not batch_lines:
            break
        # print(batch_lines)
        # Prepare the messages
        messages = [{
            'Id': str(i),
            'MessageBody': row
        } for i, row in enumerate(batch_lines)]

        # print(messages)
        num_buckets = 8
        buckets = divide_into_buckets(messages, num_buckets)
        workers = []

        for bucket in buckets:
            # print(bucket)

            w = multiprocessing.Process(target=final, args=(bucket,))
            w.start()
            workers.append(w)

        for w in workers:
            w.join()
end_time = time.perf_counter()
print(f"total time :- {end_time - start_time}")



























# # Send data to SQS queue in batches
# for i in range(0, len(messages), 10):
#     response = sqs.send_message_batch(
#         QueueUrl=queue_url,
#         Entries=messages[i:i+10]
#     )

# # Send data to SQS queue
# for row in result.result_rows:
#     message = {
#         "handle": row[0],
#         "follower_handle": row[1],
#         "follower_full_name": row[2]
#     }
#     sqs = session.resource('sqs', region_name='ap-south-1')
#     queue = sqs.get_queue_by_name(QueueName='creator_follower_in')
#     queue.send_message(MessageBody=json.dumps(message))
#     queue.send_message_batch


# messages = [
#     {
#         "handle": row[0],
#         "follower_handle": row[1],
#         "follower_full_name": row[2]
#     }
#     for row in result.result_rows]

# print(messages)
# # Send data to SQS queue in batches
# response = sqs.send_message_batch(
#     QueueUrl=queue_url,
#     Entries=messages
# )

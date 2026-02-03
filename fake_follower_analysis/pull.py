# from datetime import datetime
# import json
# import time
# import boto3

# column_name = ['handle', 'follower_handle', 'follower_full_name', 'symbolic_name', 'fake_real_based_on_lang',
#                'transliterated_follower_name', 'decoded_name', 'cleaned_handle', 'cleaned_name', 'chhitij_logic',
#                'number_handle', 'number_more_than_4_handle', 'similarity_score', 'fake_real_based_on_fuzzy_score_90',
#                'process1_', 'final_', 'numeric_handle', 'indian_name_score', 'score_80']

# kinesis = boto3.client('kinesis')

# shards = kinesis.list_shards(
#     StreamName='creator_out',
#     MaxResults=123,
#     StreamARN='arn:aws:kinesis:ap-south-1:495506833699:stream/creator_out'
# )

# # print(shards)

# current_date = datetime.today()
# file_name = f"anshul_{current_date.year}_{current_date.month}_{current_date.day}_creator_followers_final_fake_analysis.json"

# for shard in shards['Shards']:
#     print(shard)
#     shard_id = shard['ShardId']
#     print(shard_id)
#     StartingSequenceNumber = shard['SequenceNumberRange']['StartingSequenceNumber']
#     print(StartingSequenceNumber)
# # TRIM_HORIZON
#     shard_iterator = kinesis.get_shard_iterator(
#         StreamName='creator_out',
#         ShardId=shard_id,
#         ShardIteratorType='AFTER_SEQUENCE_NUMBER',
#         # Timestamp=datetime(2023,10,16, 10, 26, 25),
#         StartingSequenceNumber=StartingSequenceNumber,
#     )['ShardIterator']

#     print(datetime(2023, 10, 16, 10, 26, 25))
#     print(shard_iterator)
#     start = time.perf_counter()
#     messages = []
#     total_length = 0
#     while True:
#         response = kinesis.get_records(
#             ShardIterator=shard_iterator,
#             Limit=10000,
#             StreamARN='arn:aws:kinesis:ap-south-1:495506833699:stream/creator_out'
#         )
#         if len(response['Records']) < 1:
#             break
#         shard_iterator = response['NextShardIterator']
#         response = response['Records']
#         total_length += len(response)
#         with open(file_name, 'a') as f:
#             print("-----------------------")
#             for message in response:
#                 response = json.loads(message['Data'])
#                 response = list(response.values())
#                 row_dict = dict(zip(column_name, response))
#                 try:
#                     f.write(json.dumps(row_dict) + '\n')
#                 except Exception as e:
#                     print(e)
#     print(total_length)
#     print(f"fetch from sqs queue {time.perf_counter() - start}")


from datetime import datetime
import json
import time
import boto3
import logging
from multiprocessing import Pool

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

column_name = ['handle', 'follower_handle', 'follower_full_name', 'symbolic_name', 'fake_real_based_on_lang',
               'transliterated_follower_name', 'decoded_name', 'cleaned_handle', 'cleaned_name', 'chhitij_logic',
               'number_handle', 'number_more_than_4_handle', 'similarity_score', 'fake_real_based_on_fuzzy_score_90',
               'process1_', 'final_', 'numeric_handle', 'indian_name_score', 'score_80']

client = boto3.client('kinesis')
stream_name = 'creator_out'

current_date = datetime.today()
file_name = f"{current_date.year}_{current_date.month}_{current_date.day}_creator_followers_final_fake_analysis.json"

def process_shard(shard_id, starting_sequence_number):
    logger.info(f"Processing shard {shard_id}")
    start = time.perf_counter()
    total_length = 0
    shard_iterator = client.get_shard_iterator(
        StreamName=stream_name,
        ShardId=shard_id,
        ShardIteratorType='AFTER_SEQUENCE_NUMBER',
        StartingSequenceNumber=starting_sequence_number
    )['ShardIterator']
    while True:
        messages = []
        response = client.get_records(
            ShardIterator=shard_iterator,
            Limit=10000
        )
        if len(response['Records']) < 1:
            break
        shard_iterator = response['NextShardIterator']
        response = response['Records']
        total_length += len(response)
        for record in response:
            response = json.loads(record['Data'])
            response = list(response.values())
            row_dict = dict(zip(column_name, response))
            messages.append(row_dict)
        with open(file_name, 'a') as f:
            for message in messages:
                try:
                    f.write(json.dumps(message) + '\n')
                except Exception as e:
                    logger.error(e)
    logger.info(f"Processed {total_length} records from shard {shard_id} in {time.perf_counter() - start} seconds")

if __name__ == '__main__':
    pool = Pool()
    shards = client.list_shards(StreamName=stream_name)['Shards']
    args = [(shard['ShardId'], shard['SequenceNumberRange']['StartingSequenceNumber']) for shard in shards]
    pool.starmap(process_shard, args)
    pool.close()
    pool.join()


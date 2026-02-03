import json
import time
import boto3


  
kinesis = boto3.client('kinesis')

start = time.perf_counter()
messages = []
total_length = 0

data = {"handle":"zym1.0","follower_handle":"rahul_prasad27","follower_full_name":"Rahul Prasad"}

print(type(data))
response = kinesis.put_record(
    StreamName='creator_out',
    Data=json.dumps(data),
    PartitionKey='follower_handle',
    StreamARN='arn:aws:kinesis:ap-south-1:495506833699:stream/creator_out'
)
print(response)















# while True:    
#     response = kinesis.get_records(
#         ShardIterator=shard_iterator,
#         Limit=10000,
#         StreamARN='arn:aws:kinesis:ap-south-1:495506833699:stream/creator_out'
#     )
#     # break
#     if len(response['Records'])<1:
#         break
#     shard_iterator = response['NextShardIterator']
#     response = response['Records']
#     total_length += len(response)
#     with open('fake_output.json', 'a') as f:
#         for message in response:
#             response = json.loads(message['Data'])
#             response = list(response.values())
#             row_dict = dict(zip(column_name, response))
#             try:
#                 f.write(json.dumps(row_dict) + '\n')
#             except Exception as e:
#                 print(e)
# print(total_length)
# print(f"fetch from sqs queue {time.perf_counter() - start}")
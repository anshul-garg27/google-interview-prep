import os
import pandas as pd
import requests
import psycopg2
from psycopg2 import extras
import csv

def readCsvFileinArray(filename):
    data = []
    current_dir = os.path.dirname(os.path.abspath(__file__))
    filepath = os.path.join(current_dir, filename)
    with open(filepath, 'r') as csv_file:
        csv_reader = csv.reader(csv_file)
        for row in csv_reader:
            data.append(row[0])
    return data


def getDataFromBeat(row_dict):

    data = {}
    timeout_seconds = 15
    try:
        if row_dict['platform'] == 'instagram':
            handle = row_dict['instagram_handle']
        elif row_dict['platform'] == 'youtube':
            handle = row_dict['youtube_channel_url']
        url = "http://beat.goodcreator.co/profiles/"+ row_dict['platform'] +"/byhandle/"+ handle
        response = requests.get(url, timeout=timeout_seconds)
        response.raise_for_status()
        data = response.json()
        return data
    except requests.Timeout:
        return data
    except requests.RequestException as e:
        return data


def UpdateInstagramUsingBeat(profile):
    insertDict = {}
    if profile['full_name'] is not None:
        name =  profile['full_name']
    else:
        name = ""
    insertDict['handle'] = profile['handle']
    insertDict['business_id'] = profile['fbid']
    insertDict['thumbnail'] = profile['profile_pic_url']
    insertDict['ig_id'] = profile['profile_id']
    insertDict['bio'] = profile['biography']
    insertDict['following'] = profile['following']
    insertDict['followers'] = profile['followers']
    insertDict['name'] = name
    insertDict['search_phrase'] = name + " " + profile['handle']
   
    update_query = 'UPDATE instagram_account SET business_id = %s, thumbnail = %s, ig_id = %s, bio = %s, following = %s, followers = %s, name = %s, search_phrase = %s, deleted = %s WHERE handle = %s'
    
    coffee_postgres_cursor.execute(update_query, (insertDict['business_id'], insertDict['thumbnail'],insertDict['ig_id'] ,insertDict['bio'],insertDict['following'],insertDict['followers'],insertDict['name'],insertDict['search_phrase'],False, insertDict['handle']))
    coffee_postgres_conn.commit()

def UpdateYoutubeUsingBeat(profile,channel_id):
    insertDict = {}
    insertDict['channel_id'] = profile['channel_id']
    insertDict['title'] = profile['title']
    insertDict['followers'] = profile['subscribers']
    insertDict['uploads_count'] = profile['uploads']
    insertDict['views_count'] = profile['views']
    insertDict['thumbnail'] = profile['thumbnail']
    insertDict['search_phrase'] =  profile['title']
   
    update_query = 'UPDATE youtube_account SET title = %s, followers = %s, uploads_count = %s, views_count = %s, thumbnail = %s, search_phrase = %s,deleted = %s WHERE channel_id = %s'
    
    coffee_postgres_cursor.execute(update_query, (insertDict['title'], insertDict['followers'],insertDict['uploads_count'] ,insertDict['views_count'],insertDict['thumbnail'],insertDict['search_phrase'],False, channel_id))
    coffee_postgres_conn.commit()


def updateYoutubeRows():
    query = 'SELECT channel_id from youtube_account where deleted = True and processed = False'
    coffee_postgres_cursor.execute(query)
    youtubeDict = coffee_postgres_cursor.fetchall()
    count = 0
    deleted_count = 0
    invalid_count = 0
    beat_call_count = 0
    updated_beat_call_count = 0
    try:
        for item in youtubeDict:
            count = count +1
            channel_id = item[0]
            if channel_id in deletedHandles:
                deleted_count = deleted_count + 1
                update_query = 'UPDATE youtube_account SET processed = %s WHERE channel_id = %s'
                coffee_postgres_cursor.execute(update_query, (True, channel_id))
                coffee_postgres_conn.commit()
            elif channel_id in invalidHandles:
                invalid_count = invalid_count + 1
                update_query = 'UPDATE youtube_account SET deleted = %s, processed = %s WHERE channel_id = %s'
                coffee_postgres_cursor.execute(update_query, (False,True, channel_id))
                coffee_postgres_conn.commit()
            else:
                beat_call_count = beat_call_count + 1
                row_dict = {}
                row_dict['platform'] = 'youtube'
                row_dict['youtube_channel_url'] = channel_id
                beatData = getDataFromBeat(row_dict)
                if 'profile' in beatData:
                    profile = beatData['profile']
                    if 'channel_id' in profile:
                        print(profile['channel_id'])
                        updated_beat_call_count = updated_beat_call_count + 1
                        UpdateYoutubeUsingBeat(profile,channel_id)
                    else:
                        deleted_count = deleted_count + 1
                        update_query = 'UPDATE youtube_account SET processed = %s WHERE channel_id = %s'
                        coffee_postgres_cursor.execute(update_query, (True, channel_id))
                        coffee_postgres_conn.commit()
                else:
                    deleted_count = deleted_count + 1
                    update_query = 'UPDATE youtube_account SET processed = %s WHERE channel_id = %s'
                    coffee_postgres_cursor.execute(update_query, (True, channel_id))
                    coffee_postgres_conn.commit()
                        

            if count > 0 and (count % chunk_size == 0 or count == len(youtubeDict)):
                print(str(count)+ " Youtube Completed with " + str(deleted_count) + " Deleted and " + str(invalid_count)+ " Invalid and " + str(beat_call_count) + " Beat Calls " + str(updated_beat_call_count) + " Updated By beat")
    except Exception as err:
        print(err)
        deleted_count = deleted_count + 1
        update_query = 'UPDATE youtube_account SET processed = %s WHERE channel_id = %s'
        coffee_postgres_cursor.execute(update_query, (True, channel_id))
        coffee_postgres_conn.commit()
        
# coffee_pg_host = 'localhost'
# coffee_pg_port = 5432
# coffee_pg_database = 'coffee'
# coffee_pg_user = 'root'
# coffee_pg_password = 'Glamm@123'

coffee_pg_host = '172.31.2.21'
coffee_pg_port = 5432
coffee_pg_database = 'coffee'
coffee_pg_user = 'gccuser'
coffee_pg_password = 'dbBeat123UseRpr0d'

coffee_postgres_conn = psycopg2.connect(
    host=coffee_pg_host,
    port=coffee_pg_port,
    database=coffee_pg_database,
    user=coffee_pg_user,
    password=coffee_pg_password
)
coffee_postgres_cursor = coffee_postgres_conn.cursor(cursor_factory=extras.DictCursor)

print("Reding CSV file...RELAX")
deletedHandles = readCsvFileinArray("deleted_handle.csv")
print("Deleted Handle CSV Loaded...")

invalidHandles = readCsvFileinArray("useless_handle.csv")
print("Useless Handle CSV Loaded...")

chunk_size = 100
# updateInstagramRows()
updateYoutubeRows()

import psycopg2
from psycopg2 import extras
import time
import requests
import math
from datetime import datetime,timedelta


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

def postgresConnection():
    coffee_postgres_conn = psycopg2.connect(
        host=coffee_pg_host,
        port=coffee_pg_port,
        database=coffee_pg_database,
        user=coffee_pg_user,
        password=coffee_pg_password
    )
    coffee_postgres_cursor = coffee_postgres_conn.cursor(cursor_factory=extras.DictCursor)
    return coffee_postgres_conn,coffee_postgres_cursor

coffee_postgres_conn,coffee_postgres_cursor = postgresConnection()

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
        # response = requests.get(url)
        response.raise_for_status()
        data = response.json()
        return data
    except requests.Timeout:
        return data
    except requests.RequestException as e:
        return data

def UpdateInstagramUsingBeat(profile,instagram_handle):
    try:
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
        insertDict['last_synced_at'] = datetime.fromtimestamp(time.time())
    
        update_query = 'UPDATE instagram_account SET business_id = %s, thumbnail = %s, ig_id = %s, bio = %s, following = %s, followers = %s, name = %s, search_phrase = %s, deleted = %s, last_synced_at = %s WHERE handle = %s'
        
        coffee_postgres_cursor.execute(update_query, (insertDict['business_id'], insertDict['thumbnail'], insertDict['ig_id'], insertDict['bio'], insertDict['following'], insertDict['followers'], insertDict['name'], insertDict['search_phrase'], False, insertDict['last_synced_at'], instagram_handle))
        coffee_postgres_conn.commit()
        return True
    except Exception as err:
        return False

def updateInstagramRows():
    query = 'SELECT handle from instagram_account where ig_id LIKE \'MISSING%\''
    coffee_postgres_cursor.execute(query)
    instagramDict = coffee_postgres_cursor.fetchall()
    
    count = 0
    beat_call_count = 0
    updated_beat_call_count = 0
    start_time = time.time()
    try:
        print("Total Items That Needs To be Updated is : "+ str(len(instagramDict)))
        for item in instagramDict:
            count = count + 1
            instagram_handle = item[0].lower()
            beat_call_count = beat_call_count + 1
            # id = item[1]
            row_dict = {}
            row_dict['platform'] = 'instagram'
            row_dict['instagram_handle'] = instagram_handle
            beatData = getDataFromBeat(row_dict)
            if 'profile' in beatData:
                profile = beatData['profile']
                if 'handle' in profile:
                    status = UpdateInstagramUsingBeat(profile,instagram_handle)
                    if status:
                        updated_beat_call_count = updated_beat_call_count + 1

            if count > 0 and (count % chunk_size == 0 or count == len(instagramDict)):
                end_time = time.time()
                exec_time = math.ceil(end_time - start_time)
                print(str(count)+ " Instagram Completed with " + str(beat_call_count) + " Beat Calls " + str(updated_beat_call_count) + " Updated By beat in " + str(exec_time)+ " sec")
                start_time = time.time()
    except Exception as err:
        print(err)
        


chunk_size = 10
updateInstagramRows()

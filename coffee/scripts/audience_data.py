from clickhouse_driver import Client
import psycopg2
import os
import boto3
import time
import json
import math

def emptyErrorFile():
    # MAKE FILE EMPTY ON START
    current_dir = os.path.dirname(os.path.abspath(__file__))
    file_path = os.path.join(current_dir, "audience_data_error.txt")
    with open(file_path, 'w') as file:
        file.truncate(0)

def uploadJsonTos3():
    print("File start dumped To s3")
    clickhouse_query = 'INSERT INTO FUNCTION s3(\'https://gcc-social-data.s3.ap-south-1.amazonaws.com/temp/coffee/audience/audience_data.json\', \'AKIAXGXUCIERSOUELY73\', \'dHb+F8Nj9o0imlDFq0gq/pCnWzSL/F5kkrywtSxO\', \'JSONEachRow\') SELECT * from dbt.mart_audience_info SETTINGS s3_truncate_on_insert=1, output_format_json_quote_64bit_integers=0'
    clickhouse_client.execute(clickhouse_query)

def downloadCSV():
    print("File Start Downloading")
    ACCESS_KEY_ID = "AKIAXGXUCIERSOUELY73"
    SECRET_ACCESS_KEY = "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"
    session = boto3.Session(
        aws_access_key_id=ACCESS_KEY_ID,
        aws_secret_access_key=SECRET_ACCESS_KEY,
    )
    s3 = session.client("s3")
    BUCKET_NAME = "gcc-social-data"

   
    FILE_KEY = "temp/coffee/audience/audience_data.json"
    filename = "audience_data.json"

    current_dir = os.path.dirname(os.path.abspath(__file__))
    file_path = os.path.join(current_dir, filename)
    s3.download_file(BUCKET_NAME, FILE_KEY, file_path)
    print("Downloading Complete")

def read_large_json(file_path):
    try:
        with open(file_path, 'r', encoding='utf-8', errors='ignore') as json_file:
            for line in json_file:
                try:
                    yield json.loads(line)
                except Exception as err:
                    yield {}  # Return an empty dictionary
    except Exception as err:
        current_dir = os.path.dirname(os.path.abspath(__file__))
        file_path = os.path.join(current_dir, "audience_data_error.txt")
        with open(file_path, "a") as f:
            f.write(err)
            f.write("\n")
            f.close()
        yield {}

def readJsonFile(count,errorcount,start_time):
    filename = "audience_data.json"    
    current_dir = os.path.dirname(os.path.abspath(__file__))
    file_path = os.path.join(current_dir, filename)
    for row in read_large_json(file_path):
        if row:
            errorcount = InsertUpdate(row,errorcount)
        else: errorcount = errorcount + 1
        count = count + 1
        if count > 0 and (count % chunk_size == 0):
            end_time = time.time()
            exec_time = math.ceil(end_time - start_time)
            print(str(count) + " Audience Data Completed with "+ str(errorcount) + " errors in " + str(exec_time)+ " seconds")
            start_time = time.time()

def formatData(row):
    city_array ={}
    for i, val in enumerate(row['cities']):
        city_name = row['cities'][i].capitalize()
        city_array[city_name] = row['city_percs'][i]
    country_array ={}
    for i, val in enumerate(row['countries']):
        country_array[row['countries'][i]] = row['country_percs'][i]
    # city_array = dict(sorted(city_array.items(), key=lambda x: x[1], reverse=True))

    male_per = female_per = 0
    maleJson = {}
    femaleJson = {}
    ageJson = {}
    for i, val in enumerate(row['age_groups']):
        bucket = row['age_groups'][i]
        if 'age_group_percs' in row and i < len(row['age_group_percs']):
            bucket_value = row['age_group_percs'][i]
        else:
            bucket_value = 0
        if bucket.startswith('M_'):
            male_per = male_per + bucket_value
            bucket_name = bucket.replace('M_', '')
            maleJson[bucket_name] = bucket_value
        elif bucket.startswith('F_'):
            female_per = female_per + bucket_value
            bucket_name = bucket.replace('F_', '')
            femaleJson[bucket_name] = bucket_value
        
        if bucket_name in ageJson:
            ageJson[bucket_name] = ageJson[bucket_name] + bucket_value
        else:
            ageJson[bucket_name] = bucket_value

    quality_breakup = {}
    if row['platform'] == 'INSTAGRAM':
        quality_breakup['AUDIENCE_AUTHENTICITY'] = row['breakup_audience_authenticity_grade']
        quality_breakup['HIGH_FOLLOWING'] = row['breakup_high_following_grade']
        quality_breakup['LOW_FOLLOWERS'] = row['breakup_low_followers_grade']
    elif row['platform'] == 'YOUTUBE':
        quality_breakup['FOLLOWERS'] = row['breakup_followers_growth30d_grade']
        quality_breakup['VIEWS'] = row['breakup_views_30d_grade']
        quality_breakup['UPLOADS'] = row['breakup_uploads_30d_grade']

        
    insert_row = {}
    insert_row['platform'] = row['platform']
    if(len(city_array)>0):
        insert_row['city_wise_audience_location'] = json.dumps(city_array)
    else:
        insert_row['city_wise_audience_location'] = None
    if(len(country_array)>0):
        insert_row['country_wise_audience_location'] = json.dumps(country_array)
    else:
        insert_row['country_wise_audience_location'] = None
    insert_row['audience_gender_split'] = json.dumps({'male_per':male_per,'female_per':female_per})
    insert_row['audience_age_gender_split'] = json.dumps({'male':maleJson,'female':femaleJson})
    insert_row['audience_age_split'] = json.dumps(ageJson)
    insert_row['audience_language'] = None
    if len(row['notable_followers']):
        insert_row['notable_followers'] = row['notable_followers']  #need To change This in apis
    else:
        insert_row['notable_followers'] = None
    insert_row['audience_reachability_percentage'] = row['reachable_followers_perc']
    insert_row['audience_authenticity_percentage'] = row['audience_authenticity']
    insert_row['quality_audience_score'] = row['audience_quality_score']
    insert_row['quality_score_grade'] = row['audience_quality_grade']
    insert_row['quality_score_breakup'] = json.dumps(quality_breakup)
    insert_row['enabled'] = True
    insert_row['quality_audience_percentage'] = row['audience_authenticity']
    insert_row['platform_id'] = row['platform_id']
    insert_row['audience_location_split'] = json.dumps({'city':city_array,'country':country_array})
    return insert_row

def InsertUpdate(row,errorcount):
    try:
        row = formatData(row)
        check_query = 'SELECT * FROM social_profile_audience_info WHERE platform_id = %s'
        postgres_cursor.execute(check_query, (row['platform_id'],))
        existing_row = postgres_cursor.fetchone()
       
        if existing_row:
            postgres_cursor.execute(postgres_update_query, (row['platform'], row['audience_gender_split'], row['audience_age_gender_split'], row['audience_age_split'], row['audience_language'], row['notable_followers'], row['audience_reachability_percentage'], row['audience_authenticity_percentage'], row['quality_audience_score'], row['quality_score_grade'], row['quality_score_breakup'], row['enabled'], row['quality_audience_percentage'], row['audience_location_split'], row['platform_id']))
        else:
            postgres_cursor.execute(postgres_insert_query, row)
            
    except Exception as err:
        current_dir = os.path.dirname(os.path.abspath(__file__))
        file_path = os.path.join(current_dir, "audience_data_error.txt")

        with open(file_path, "a") as f:
            f.write(row['platform_id']+ "_" +row['platform'])
            f.write("\n")
            f.close()
        errorcount = errorcount+1
              
    postgres_conn.commit()
    return errorcount


############################### MAIN BLOCK STARTS #############################################################
ch_host = '172.31.28.68'
ch_port = 9000
ch_database = 'dbt'
ch_user = 'coffee'
ch_password = '0M7sooN_WQW'
clickhouse_client = Client(ch_host, port=ch_port, database=ch_database, user=ch_user, password=ch_password)

# Define connection parameters for PostgreSQL

# pg_host = 'localhost'
# pg_port = 5432
# pg_database = 'coffee'
# pg_user = 'root'
# pg_password = 'Glamm@123'

pg_host = '172.31.2.21'
pg_port = 5432
pg_database = 'coffee'
pg_user = 'gccuser'
pg_password = 'dbBeat123UseRpr0d'


# Create a connection object for PostgreSQL
postgres_conn = psycopg2.connect(
    host=pg_host,
    port=pg_port,
    database=pg_database,
    user=pg_user,
    password=pg_password
)

postgres_cursor = postgres_conn.cursor()

postgres_insert_query = 'INSERT INTO social_profile_audience_info (platform,audience_gender_split,audience_age_gender_split,audience_age_split,audience_language,notable_followers,audience_reachability_percentage,audience_authenticity_percentage,quality_audience_score,quality_score_grade,quality_score_breakup,enabled,quality_audience_percentage,audience_location_split,platform_id) VALUES (%(platform)s, %(audience_gender_split)s, %(audience_age_gender_split)s, %(audience_age_split)s, %(audience_language)s, %(notable_followers)s, %(audience_reachability_percentage)s, %(audience_authenticity_percentage)s, %(quality_audience_score)s, %(quality_score_grade)s, %(quality_score_breakup)s, %(enabled)s, %(quality_audience_percentage)s, %(audience_location_split)s, %(platform_id)s)'

postgres_update_query = 'UPDATE social_profile_audience_info SET platform = %s, audience_gender_split = %s, audience_age_gender_split = %s, audience_age_split = %s, audience_language = %s, notable_followers = %s, audience_reachability_percentage = %s, audience_authenticity_percentage = %s, quality_audience_score = %s, quality_score_grade = %s, quality_score_breakup = %s, enabled = %s, quality_audience_percentage = %s, audience_location_split = %s WHERE platform_id = %s'
clickhouse_query = 'SELECT * FROM mart_audience_info'


errorcount = 0
count = 0
chunk_size = 5000
emptyErrorFile()
uploadJsonTos3()
downloadCSV()
start_time = time.time()
readJsonFile(count,errorcount,start_time)


print("social_profile_audience_info completed. Check audience_data_error.txt for the platform_id's")

# To Update Instagram Platform Id in social_profile_audience info
chunk_size = 10000 
offset = 0
while True:
    query = f"""
    WITH cte AS (
        SELECT spi.id AS spi_id, ia.id AS ia_id
        FROM social_profile_audience_info spi
        JOIN instagram_account ia ON spi.platform_id = ia.ig_id
        WHERE spi.platform = 'INSTAGRAM'
        LIMIT {chunk_size}
        OFFSET {offset}
    )
    UPDATE social_profile_audience_info spi
    SET platform_profile_id = cte.ia_id
    FROM cte
    WHERE spi.id = cte.spi_id;
    """
    
    postgres_cursor.execute(query)
    postgres_conn.commit()
    
    # Check if there are more rows to update
    if postgres_cursor.rowcount == 0:
        break
    
    offset += chunk_size
    print(str(offset)+ " INSTAGRAM platform_profile_id Get Updated")

print("platform_profile_id column of social_profile_audience_info TABLE for INSTAGRAM GOT UPDATED")

# To Update Youtube Platform Id in social_profile_audience info
chunk_size = 10000 
offset = 0
while True:
    query = f"""
    WITH cte AS (
        SELECT spi.id AS spi_id, ya.id AS ya_id
        FROM social_profile_audience_info spi
        JOIN youtube_account ya ON spi.platform_id = ya.channel_id
        WHERE spi.platform = 'YOUTUBE'
        LIMIT {chunk_size}
        OFFSET {offset}
    )
    UPDATE social_profile_audience_info spi
    SET platform_profile_id = cte.ya_id
    FROM cte
    WHERE spi.id = cte.spi_id;
    """
    
    postgres_cursor.execute(query)
    postgres_conn.commit()
    
    # Check if there are more rows to update
    if postgres_cursor.rowcount == 0:
        break
    
    offset += chunk_size
    print(str(offset)+ " YOUTUBE platform_profile_id Get Updated")

print("platform_profile_id column of social_profile_audience_info TABLE for YOUTUBE GOT UPDATED")

# To Update Audience Data IN IA
offset = 0
chunk_size = 10000
while True:
    query1 = f"""
    WITH cte AS (
        SELECT spi.audience_gender_split, spi.audience_age_split, spi.audience_location_split, yt.channel_id
        FROM social_profile_audience_info spi
        JOIN youtube_account yt ON spi.platform_id = yt.channel_id
        LIMIT {chunk_size}
        OFFSET {offset}
    )
    UPDATE youtube_account yt
    SET audience_gender = cte.audience_gender_split,audience_age = cte.audience_age_split, audience_location = cte.audience_location_split
    FROM cte
    WHERE yt.channel_id = cte.channel_id;
    """
    postgres_cursor.execute(query1)
    postgres_conn.commit()
    if postgres_cursor.rowcount == 0:
        break
    
    offset += chunk_size
    print(str(offset)+ " youtube_account audience_data Get Updated")
print("youtube_account TABLE audience columns GOT UPDATED")

# To Update Audience Data IN YA
offset = 0
chunk_size = 10000
while True:
    query1 = f"""
    WITH cte AS (
        SELECT spi.audience_gender_split, spi.audience_age_split, spi.audience_location_split, ia.ig_id
        FROM social_profile_audience_info spi
        JOIN instagram_account ia ON spi.platform_id = ia.ig_id
        LIMIT {chunk_size}
        OFFSET {offset}
    )
    UPDATE instagram_account ia
    SET audience_gender = cte.audience_gender_split, audience_age = cte.audience_age_split, audience_location = cte.audience_location_split
    FROM cte
    WHERE ia.ig_id = cte.ig_id;
    """
    postgres_cursor.execute(query1)
    postgres_conn.commit()
    if postgres_cursor.rowcount == 0:
        break
    
    offset += chunk_size
    print(str(offset)+ " instagram_account audience_data Get Updated")
print("instagram_account TABLE audience data columns GOT UPDATED")

postgres_cursor.close()
postgres_conn.close()
clickhouse_client.disconnect()



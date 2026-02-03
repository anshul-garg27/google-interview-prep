import psycopg2
from psycopg2 import extras
import os
import pandas as pd
import csv
from datetime import datetime
import json
import time
import math
import traceback
import boto3
import sys


from clickhouse_driver import Client

coffee_pg_host = '172.31.2.21'
coffee_pg_port = 5432
coffee_pg_database = 'coffee'
coffee_pg_user = 'gccuser'
coffee_pg_password = 'dbBeat123UseRpr0d'

# coffee_pg_host = 'localhost'
# coffee_pg_port = 5432
# coffee_pg_database = 'coffee'
# coffee_pg_user = 'root'
# coffee_pg_password = 'Glamm@123'

ch_host = '172.31.28.68'
ch_port = 9000
ch_database = 'dbt'
ch_user = 'cube'
ch_password = 'CuB31Z@3995ddOpT'

clickhouse_client = Client(ch_host, port=ch_port, database=ch_database, user=ch_user, password=ch_password)

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

def writeInErrorFile(data_array):
    current_dir = os.path.dirname(os.path.abspath(__file__))
    file_path = os.path.join(current_dir, "updatesearchPhraseError1.csv")
    with open(file_path, "a", newline='') as f:
        writer = csv.writer(f)
        writer.writerow(data_array)

def read_large_json(file_path):
    with open(file_path, 'r') as json_file:
        for line in json_file:
            yield json.loads(line)

def processJsonFIle(filename):
    count = 0
    error_count = 0
    updated_count = 0
    current_dir = os.path.dirname(os.path.abspath(__file__))
    file_path = os.path.join(current_dir, filename)
    start_time = time.time()
    for data_row in read_large_json(file_path):
        
        try:
            search_phrase_admin = ''
            keywords_admin_set = set()
            if  data_row['admin_details_name'] is not None and data_row['admin_details_name'] != "":
                name_parts = data_row['admin_details_name'].split()  # Split the name by spaces
                lowercase_name_parts = [part.lower() for part in name_parts]  # Convert parts to lowercase
                keywords_admin_set.update(lowercase_name_parts)  # Use update() to add elements to the set
                if search_phrase_admin is not None:
                    search_phrase_admin = search_phrase_admin + " " + data_row['admin_details_name'].lower()
                else:
                    search_phrase_admin = data_row['admin_details_name'].lower()

            if  data_row['admin_details_email'] is not None and data_row['admin_details_email'] != "":
                email = data_row['admin_details_email'].lower()
                keywords_admin_set.add(email)  # Use add() to add elements to the set
                if search_phrase_admin is not None:
                    search_phrase_admin = search_phrase_admin + " " + email
                else:
                    search_phrase_admin = email

            if  data_row['admin_details_phone'] is not None and data_row['admin_details_phone'] != "":
                phone = data_row['admin_details_phone'].lower()
                keywords_admin_set.add(phone)  # Use add() to add elements to the set
                if search_phrase_admin is not None:
                    search_phrase_admin = search_phrase_admin + " " + phone
                else:
                    search_phrase_admin = phone

            if  data_row['admin_details_secondaryPhone'] is not None and data_row['admin_details_secondaryPhone'] != "":
                secondaryPhone = data_row['admin_details_secondaryPhone'].lower()
                keywords_admin_set.add(secondaryPhone)  # Use add() to add elements to the set
                if search_phrase_admin is not None:
                    search_phrase_admin = search_phrase_admin + " " + secondaryPhone
                else:
                    search_phrase_admin = secondaryPhone

            if  data_row['user_details_name'] is not None and data_row['user_details_name'] != "":
                name_parts = data_row['user_details_name'].split()  # Split the name by spaces
                lowercase_name_parts = [part.lower() for part in name_parts]  # Convert parts to lowercase
                keywords_admin_set.update(lowercase_name_parts)
                if search_phrase_admin is not None:
                    search_phrase_admin = search_phrase_admin + " " + data_row['user_details_name'].lower()
                else:
                    search_phrase_admin = data_row['user_details_name'].lower()

            if  data_row['user_details_email'] is not None and data_row['user_details_email'] != "":
                email = data_row['user_details_email'].lower()
                keywords_admin_set.add(email)
                if search_phrase_admin is not None:
                    search_phrase_admin = search_phrase_admin + " " + email
                else:
                    search_phrase_admin = email

            if  data_row['user_details_phone'] is not None and data_row['user_details_phone'] != "":
                phone = data_row['user_details_phone'].lower()
                keywords_admin_set.add(phone)
                if search_phrase_admin is not None:
                    search_phrase_admin = search_phrase_admin + " " + phone
                else:
                    search_phrase_admin = phone

            if  data_row['user_details_secondaryPhone'] is not None and data_row['user_details_secondaryPhone'] != "":
                secondaryPhone = data_row['user_details_secondaryPhone'].lower()
                keywords_admin_set.add(secondaryPhone)
                if search_phrase_admin is not None:
                    search_phrase_admin = search_phrase_admin + " " + secondaryPhone
                else:
                    search_phrase_admin = secondaryPhone

            # Convert the set back to a list if needed

            platform_account_id = data_row['platform_account_id']
            platform = data_row['platform']
            
            if platform == "INSTAGRAM":
                check_mapping_query = "Select * from instagram_account where id = %s"
                coffee_postgres_cursor.execute(check_mapping_query, (platform_account_id,))
                social_profile = coffee_postgres_cursor.fetchone()
                if social_profile is not None:
                    if social_profile['name'] is not None:
                        name_parts = social_profile['name'].split()  
                        lowercase_name_parts = [part.lower() for part in name_parts]
                        keywords_admin_set.update(lowercase_name_parts)
                        if search_phrase_admin is not None:
                            search_phrase_admin = search_phrase_admin + " " + social_profile['name'].lower()
                        else:
                            search_phrase_admin = social_profile['name'].lower()

                    if social_profile['handle'] is not None:
                        keywords_admin_set.add(social_profile['handle'].lower())
                        if search_phrase_admin is not None:
                            search_phrase_admin = search_phrase_admin + " " + social_profile['handle'].lower()
                        else:
                            search_phrase_admin = social_profile['handle'].lower()
                
                    keywords_admin = list(keywords_admin_set)
                    # search_phrase_admin = ', '.join(keywords_admin)

                    update_query = 'UPDATE instagram_account SET keywords_admin = %s, search_phrase_admin = %s WHERE id = %s'
                    coffee_postgres_cursor.execute(update_query, (keywords_admin, search_phrase_admin, platform_account_id))
                    coffee_postgres_conn.commit()
                    updated_count = updated_count + 1



            elif platform == "YOUTUBE":
                check_mapping_query = "Select * from youtube_account where id = %s"
                coffee_postgres_cursor.execute(check_mapping_query, (platform_account_id,))
                social_profile = coffee_postgres_cursor.fetchone()
                if social_profile is not None:
                    if social_profile['title'] is not None:
                        name_parts = social_profile['title'].split()  
                        lowercase_name_parts = [part.lower() for part in name_parts]
                        keywords_admin_set.update(lowercase_name_parts)
                        if search_phrase_admin is not None:
                            search_phrase_admin = search_phrase_admin + " " + social_profile['title'].lower()
                        else:
                            search_phrase_admin = social_profile['title'].lower()

                    if social_profile['username'] is not None:
                        username = social_profile['username'].lower()
                        keywords_admin_set.add(username) 
                        if search_phrase_admin is not None: 
                            search_phrase_admin = search_phrase_admin + " " + username
                        else:
                            search_phrase_admin = username 
                
                    keywords_admin = list(keywords_admin_set)
                    # search_phrase_admin = ', '.join(keywords_admin)

                    update_query = 'UPDATE youtube_account SET keywords_admin = %s, search_phrase_admin = %s WHERE id = %s'
                    coffee_postgres_cursor.execute(update_query, (keywords_admin, search_phrase_admin, platform_account_id))
                    coffee_postgres_conn.commit()
                    updated_count = updated_count + 1


            count = count + 1
            if count > 0 and (count % print_chunk_size == 0):
                end_time = time.time()
                exec_time = math.ceil(end_time - start_time)
                start_time = time.time()

                print(str(count)+" Processed in " + str(updated_count) + " updated with " + str(error_count) + " errors in " + str(exec_time) + " seconds")
        except Exception as err:
            error_dict=[]
            error_dict.append(str(data_row['id']))
            writeInErrorFile(error_dict)        
            traceback.print_exc()
            error_count = error_count + 1

def downloadCSV(filename):
    print("File Start Downloading")
    ACCESS_KEY_ID = "AKIAXGXUCIERSOUELY73"
    SECRET_ACCESS_KEY = "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"
    session = boto3.Session(
        aws_access_key_id=ACCESS_KEY_ID,
        aws_secret_access_key=SECRET_ACCESS_KEY,
    )
    s3 = session.client("s3")
    BUCKET_NAME = "gcc-social-data"
    FILE_KEY = "temp/coffee/cp/campaign_profiles.json"

    current_dir = os.path.dirname(os.path.abspath(__file__))
    file_path = os.path.join(current_dir, filename)
    s3.download_file(BUCKET_NAME, FILE_KEY, file_path)
    print("Downloading Complete")

def uploadCamapaignProfileJsonTos3(query):
    print("File start dumped To s3")
    clickhouse_query = 'INSERT INTO FUNCTION s3(\'https://gcc-social-data.s3.ap-south-1.amazonaws.com/temp/coffee/cp/campaign_profiles.json\', \'AKIAXGXUCIERSOUELY73\', \'dHb+F8Nj9o0imlDFq0gq/pCnWzSL/F5kkrywtSxO\', \'JSONEachRow\') ' + query + ' SETTINGS s3_truncate_on_insert=1, output_format_json_quote_64bit_integers=0'
    clickhouse_client.execute(clickhouse_query)

def emptyErrorFile(error_file_name):
    current_dir = os.path.dirname(os.path.abspath(__file__))
    error_file_path = os.path.join(current_dir, error_file_name)
    with open(error_file_path, 'w') as file:
        file.truncate(0)

def readErrorFileInChunk(filename, chunk_size):
    current_dir = os.path.dirname(os.path.abspath(__file__))
    file_path = os.path.join(current_dir, filename)
    
    csv_reader = pd.read_csv(file_path, chunksize=chunk_size)
    for chunk in csv_reader:
        yield','.join(map(str, chunk.values.flatten().tolist()))

def main():
    global coffee_postgres_conn,coffee_postgres_cursor,download_chunk_size,print_chunk_size
    error_file_name = "updatesearchPhraseError.csv"
    download_chunk_size = 10000
    print_chunk_size = 1000
    coffee_postgres_conn,coffee_postgres_cursor = postgresConnection()
    
    cron_type = ''
    if len(sys.argv) > 1:
        cron_type = sys.argv[1] # partial or full
    else:
        cron_type = '' 
    if cron_type == "deleted":
        for csv_string in readErrorFileInChunk(error_file_name,download_chunk_size):
            dataQuery = "SELECT JSONExtractString(admin_details, 'name') admin_details_name, JSONExtractString(admin_details, 'email') admin_details_email, JSONExtractString(admin_details, 'phone') admin_details_phone, JSONExtractString(admin_details, 'secondaryPhone') admin_details_secondaryPhone, JSONExtractString(user_details, 'name') user_details_name, JSONExtractString(user_details, 'email') user_details_email, JSONExtractString(user_details, 'phone') user_details_phone, JSONExtractString(user_details, 'secondaryPhone') user_details_secondaryPhone, platform_account_id, platform, id from dbt.stg_coffee_campaign_profiles where id IN ("+csv_string+")"
            
            uploadCamapaignProfileJsonTos3(dataQuery)
            downloadCSV("campaign_profiles.json")
            processJsonFIle("campaign_profiles.json")

    else:
        dataQuery = "SELECT JSONExtractString(admin_details, 'name') admin_details_name, JSONExtractString(admin_details, 'email') admin_details_email, JSONExtractString(admin_details, 'phone') admin_details_phone, JSONExtractString(admin_details, 'secondaryPhone') admin_details_secondaryPhone, JSONExtractString(user_details, 'name') user_details_name, JSONExtractString(user_details, 'email') user_details_email, JSONExtractString(user_details, 'phone') user_details_phone, JSONExtractString(user_details, 'secondaryPhone') user_details_secondaryPhone, platform_account_id, platform, id from dbt.stg_coffee_campaign_profiles where platform IN (\'INSTAGRAM\',\'YOUTUBE\') and (admin_details IS NOT NULL OR user_details IS NOT NULL)"
        # emptyErrorFile(error_file_name)
        uploadCamapaignProfileJsonTos3(dataQuery)
        downloadCSV("campaign_profiles.json")
        processJsonFIle("campaign_profiles.json")



if __name__ == "__main__":
    main()
import psycopg2
from psycopg2 import extras
import time
import requests
import math
import os
import csv


# coffee_pg_host = 'gcc-pg-stage'
# coffee_pg_port = 5432
# coffee_pg_database = 'coffee'
# coffee_pg_user = 'gccuser'
# coffee_pg_password = 'GCc123'

coffee_pg_host = '172.31.2.21'
coffee_pg_port = 5432
coffee_pg_database = 'coffee'
coffee_pg_user = 'gccuser'
coffee_pg_password = 'dbBeat123UseRpr0d'

chunk_size = 10
# MAKE FILE EMPTY ON START
current_dir = os.path.dirname(os.path.abspath(__file__))
file_path = os.path.join(current_dir, "check_duplicates.csv")
with open(file_path, 'w') as file:
    file.truncate(0)

def writeInFile(data_array):
    current_dir = os.path.dirname(os.path.abspath(__file__))
    file_path = os.path.join(current_dir, "check_duplicates.csv")
    
    with open(file_path, "a", newline='') as f:
        writer = csv.writer(f)
        writer.writerow(data_array)

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


query = 'Select ig_id from instagram_account where ig_id != \'\' group by ig_id having count(ig_id) > 1'
coffee_postgres_cursor.execute(query)
instagramDict = coffee_postgres_cursor.fetchall()

count = 0
start_time = time.time()
print("Duplicate Data Fetched Cron Started")
for item in instagramDict:
    ig_id = item[0]   
    check_query = 'Select id from instagram_account where ig_id = %s'
    coffee_postgres_cursor.execute(check_query, (ig_id,))
    instaIds = coffee_postgres_cursor.fetchall()
    instaIds_tuple = tuple(row[0] for row in instaIds)
    # leaderboard_query = 'Select DISTINCT platform_profile_id from leaderboard where platform_profile_id IN %s'
    # coffee_postgres_cursor.execute(leaderboard_query, (instaIds_tuple,))
    # result = coffee_postgres_cursor.fetchall()


    cp_query = 'Select id,updated_at from campaign_profiles where platform = \'INSTAGRAM\' and platform_account_id IN %s order by updated_at desc'
    coffee_postgres_cursor.execute(cp_query, (instaIds_tuple,))
    cp_result = coffee_postgres_cursor.fetchall()

    writer = []
    # for id_value in instaIds_tuple:
    #     writer.append(str(id_value))
    # for row in result:
    #     writer.append(str(row[0]))
    if len(cp_result)>1:
        for row in cp_result:
            writer.append(str(row[0]))
        writeInFile(writer)
    count = count + 1
    if count > 0 and (count % chunk_size == 0 or count == len(instagramDict)):
        end_time = time.time()
        exec_time = math.ceil(end_time - start_time)
        print(str(count)+ " Instagram Completed with " + str(exec_time)+ " sec")
        start_time = time.time()

print("completed Check csv check_duplicates.csv")
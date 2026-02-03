import os
import csv
import psycopg2
from psycopg2 import extras
import pandas as pd


def writeInFile(data_array):
    current_dir = os.path.dirname(os.path.abspath(__file__))
    file_path = os.path.join(current_dir, "check_cp_left.csv")
    
    with open(file_path, "a", newline='') as f:
        writer = csv.writer(f)
        writer.writerow(data_array)

def readCsvFile(filename,key):

    current_dir = os.path.dirname(os.path.abspath(__file__))
    filepath = os.path.join(current_dir, filename)
    data = pd.read_csv(filepath, low_memory=False)
    dict = {}
    for index, row in data.iterrows():
        handle = row[key]
        row_dict = row.to_dict()
        # Filter out NaN values from the row dictionary
        row_dict = {k: v for k, v in row_dict.items() if not pd.isna(v)}
        dict[handle] = row_dict
    return dict

def verifyIACPS(winklCps):
    count = 0
    left_count = 0
    for key, value in winklCps.items():
        check_query = 'Select platform_account_id from campaign_profiles where id = %s'
        coffee_postgres_cursor.execute(check_query, (value['id'],))
        existing_row = coffee_postgres_cursor.fetchone()
        if existing_row is None:
            writeInFile([str(value['id'])])
            left_count = left_count + 1
        else:
            if 'instagram_handle' in value and  value['instagram_handle'] is not None:
                value['instagram_handle'] = value['instagram_handle'].lower()
                if 'channel_id' in value and value['channel_id'] is not None:
                    check_ia = "Select id from instagram_account where id = %s and handle = %s and gcc_linked_channel_id = %s"
                    coffee_postgres_cursor.execute(check_ia, (existing_row['platform_account_id'],value['instagram_handle'],value['channel_id'],))
                else:
                    check_ia = "Select id from instagram_account where id = %s and handle = %s"
                    coffee_postgres_cursor.execute(check_ia, (existing_row['platform_account_id'],value['instagram_handle'],))

                existing_ia = coffee_postgres_cursor.fetchone()
                if existing_ia is None:
                    channel_id = ''
                    # if 'channel_id' in value:
                    #     channel_id = value['channel_id']
                    writeInFile([str(value['id'])])
                    left_count = left_count + 1

            
            elif 'youtube_channel_url' in value and  value['youtube_channel_url'] is not None:
                check_ya = "Select id from youtube_account where id = %s and channel_id = %s"
                coffee_postgres_cursor.execute(check_ya, (existing_row['platform_account_id'],value['youtube_channel_url'],))
                existing_ya = coffee_postgres_cursor.fetchone()
                if existing_ya is None:
                    writeInFile([str(value['id'])])
                    left_count = left_count + 1

        
        count = count + 1
        if count > 0 and (count % 1000 == 0 or count == len(winklCps)):
            print(str(count)+ " CP Completed and " + str(left_count) +" cps left in this ")
            # left_count = 0

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

# MAKE FILE EMPTY ON START
current_dir = os.path.dirname(os.path.abspath(__file__))
file_path = os.path.join(current_dir, "check_cp_left.csv")
with open(file_path, 'w') as file:
    file.truncate(0)

print("Reading Instagram CP CSV files....")
winklCps = readCsvFile("check_cp_ia.csv","instagram_handle")
writeInFile(["INSTAGRAM"])
verifyIACPS(winklCps)

print("Reading Youtube CP CSV files....")
winklCps = readCsvFile("check_cp_yt.csv","youtube_channel_url")
writeInFile(["YOUTUBE"])
verifyIACPS(winklCps)

print("Reading GCC CP CSV files....")
winklCps = readCsvFile("check_cp_gcc.csv","bulbul_user_account_id")
writeInFile(["GCC"])
verifyIACPS(winklCps)


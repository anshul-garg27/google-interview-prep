import psycopg2
from psycopg2 import extras
import os
import pandas as pd
import csv
import requests
import time
import math
from datetime import datetime
import urllib.parse
import itertools
import concurrent.futures


coffee_pg_host = '172.31.2.21'
coffee_pg_port = 5432
coffee_pg_database = 'coffee'
coffee_pg_user = 'gccuser'
coffee_pg_password = 'dbBeat123UseRpr0d'
COFFEE_URL = "http://coffee.goodcreator.co"
batch_size = 1000
max_concurrency = 8

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
    file_path = os.path.join(current_dir, "create_missing_profiles_IA_error.csv")
    with open(file_path, "a", newline='') as f:
        writer = csv.writer(f)
        writer.writerow(data_array)

def writeInFile(data_array):
    current_dir = os.path.dirname(os.path.abspath(__file__))
    file_path = os.path.join(current_dir, "create_missing_profiles_IA.csv")
    with open(file_path, "a", newline='') as f:
        writer = csv.writer(f)
        writer.writerow(data_array)

def generateCsvForInstagramHandles():
    instagram_query = 'SELECT handle FROM instagram_account order by followers desc nulls last'
    coffee_postgres_cursor.execute(instagram_query)
    result_set = coffee_postgres_cursor.fetchall()
    coffee_postgres_conn.commit()
    try:
        output_file='instagram_handles.csv'
        with open(output_file, 'w', newline='') as csvfile:
            csv_writer = csv.writer(csvfile)
            for row in result_set:
                csv_writer.writerow(row.values())
    except IOError as e:
        error_dict=[]
        error_dict.append(e)
        writeInErrorFile(error_dict)
    return os.path.abspath(output_file)

def createMissingCPForInstagram(instagramEntry):
    handle = instagramEntry[0]
    # encoded_handle = urllib.parse.quote(handle)
    url = COFFEE_URL+"/campaign-profile-service/api/campaign-profile/handle/INSTAGRAM/"+handle+"?source=GCC"
    response = requests.get(url)
    try:
        json_data = response.json()
        status_data = json_data['status']['status']
        if status_data != "SUCCESS":
            error_dict=[]
            error_dict.append(json_data)
            error_dict.append(url)
            writeInErrorFile(error_dict)
            # print(json_data)
        elif status_data == "SUCCESS":
            dict=[]
            dict.append(json_data)
            dict.append(url)
            writeInFile(dict)
    except Exception as e:
        error_dict=[]
        error_dict.append(e)
        error_dict.append(url)
        writeInErrorFile(error_dict)

def process_function():
    
    #instagram processing
    writeInErrorFile(["INSTAGRAM"])
    instagram_input_file = generateCsvForInstagramHandles()
    print("instagram_path = " + instagram_input_file)
    start_time = time.time()
    completed_size=0
    with open(instagram_input_file, 'r', newline='') as csv_file:
        while True:
            csv_reader = csv.reader(csv_file)
            batch_lines = list(itertools.islice(csv_reader, batch_size))
            if not batch_lines:
                break
            with concurrent.futures.ThreadPoolExecutor(max_workers=max_concurrency) as executor:
                batch_start_time = time.time()
                futures = [executor.submit(createMissingCPForInstagram, line) for line in batch_lines]
                concurrent.futures.wait(futures)
                completed_size += batch_size
                batch_end_time = time.time()
                batch_exec_time = math.ceil(batch_end_time - batch_start_time)
                # print("batch exec time = " + str(batch_exec_time))
                print(str(completed_size) + " completed in " + str(batch_exec_time))

    end_time = time.time()
    exec_time = math.ceil(end_time - start_time)
    print("time for creating instagram CPs = " + str(exec_time))

coffee_postgres_conn,coffee_postgres_cursor = postgresConnection()

current_dir = os.path.dirname(os.path.abspath(__file__))
error_file_path = os.path.join(current_dir, "create_missing_profiles_IA_error.csv")
with open(error_file_path, 'w') as file:
    file.truncate(0)

file_path = os.path.join(current_dir, "create_missing_profiles_IA.csv")
with open(file_path, 'w') as file:
    file.truncate(0)

process_function()
import psycopg2
from psycopg2 import extras
import pandas as pd
from datetime import datetime
import concurrent.futures
import csv
import os
import time
import itertools
import math

coffee_pg_host = '172.31.2.21'
coffee_pg_port = 5432
coffee_pg_database = 'coffee'
coffee_pg_user = 'gccuser'
coffee_pg_password = 'dbBeat123UseRpr0d'
batch_size=10000
max_concurrency = 8

def postgresConnection():
    coffee_postgres_conn = psycopg2.connect(
        host=coffee_pg_host,
        port=coffee_pg_port,
        database=coffee_pg_database,
        user=coffee_pg_user,
        password=coffee_pg_password
    )
    coffee_postgres_conn.autocommit = True
    coffee_postgres_cursor = coffee_postgres_conn.cursor(cursor_factory=extras.DictCursor)
    return coffee_postgres_conn,coffee_postgres_cursor

def writeInFile(data_array):
    current_dir = os.path.dirname(os.path.abspath(__file__))
    file_path = os.path.join(current_dir, "update_cp_error.csv")
    
    with open(file_path, "a", newline='') as f:
        writer = csv.writer(f)
        writer.writerow(data_array)

def retrieveInstagramCPDataHasEmail():
    print("creating csv file for instagram has_email")
    query = '''
            select 
                cp.id
            from
                campaign_profiles cp
            left join instagram_account ia on ia.id = cp.platform_account_id and cp.platform = 'INSTAGRAM'
            where ((cp.admin_details->>'email' IS NOT NULL AND cp.admin_details->>'email' != '' AND cp.admin_details->>'email' != 'None') OR
                (cp.user_details->>'email' IS NOT NULL AND cp.user_details->>'email' != '' AND cp.user_details->>'email' != 'None') OR
                (ia.email IS NOT NULL AND ia.email != '')) = True and ia.id>0
            '''
    coffee_postgres_cursor.execute(query)
    result_set = coffee_postgres_cursor.fetchall()
    coffee_postgres_conn.commit()
    try:
        output_file='instagram_data_has_email.csv'
        with open(output_file, 'w', newline='') as csvfile:
            csv_writer = csv.writer(csvfile)
            for row in result_set:
                csv_writer.writerow(row.values())
    except IOError as e:
        error_dict=[]
        error_dict.append(e)
        writeInFile(error_dict)
    return os.path.abspath(output_file)

def retrieveInstagramCPDataHasPhone():
    print("creating csv file for instagram has_phone")
    query = '''
            select 
                cp.id
            from
                campaign_profiles cp
            left join instagram_account ia on ia.id = cp.platform_account_id and cp.platform = 'INSTAGRAM'
            where ((cp.admin_details->>'phone' IS NOT NULL AND cp.admin_details->>'phone' != '' AND cp.admin_details->>'phone' != 'None') OR
                (cp.admin_details->>'secondaryPhone' IS NOT NULL AND cp.admin_details->>'secondaryPhone' != '' AND cp.admin_details->>'secondaryPhone' != 'None') OR
                (cp.user_details->>'phone' IS NOT NULL AND cp.user_details->>'phone' != '' AND cp.user_details->>'phone' != 'None') or
                (cp.user_details->>'secondaryPhone' IS NOT NULL AND cp.user_details->>'secondaryPhone' != '' AND cp.user_details->>'secondaryPhone' != 'None') OR
                (ia.phone IS NOT NULL AND ia.phone != '')) = True and ia.id>0
            '''
    coffee_postgres_cursor.execute(query)
    result_set = coffee_postgres_cursor.fetchall()
    coffee_postgres_conn.commit()
    try:
        output_file='instagram_data_has_phone.csv'
        with open(output_file, 'w', newline='') as csvfile:
            csv_writer = csv.writer(csvfile)
            for row in result_set:
                csv_writer.writerow(row.values())
    except IOError as e:
        error_dict=[]
        error_dict.append(e)
        writeInFile(error_dict)
    return os.path.abspath(output_file)

def retrieveYoutubeCPDataHasEmail():
    print("creating csv file for youtube has_email")
    query = '''
            select 
                cp.id
            from
                campaign_profiles cp
            left join youtube_account ya on ya.id = cp.platform_account_id and cp.platform = 'YOUTUBE'
            where ((cp.admin_details->>'email' IS NOT NULL AND cp.admin_details->>'email' != '' AND cp.admin_details->>'email' != 'None') OR
                (cp.user_details->>'email' IS NOT NULL AND cp.user_details->>'email' != '' AND cp.user_details->>'email' != 'None') OR
                (ya.email IS NOT NULL AND ya.email != '')) = True and ya.id>0
            '''
    coffee_postgres_cursor.execute(query)
    result_set = coffee_postgres_cursor.fetchall()
    coffee_postgres_conn.commit()
    try:
        output_file='youtube_data_has_email.csv'
        with open(output_file, 'w', newline='') as csvfile:
            csv_writer = csv.writer(csvfile)
            for row in result_set:
                csv_writer.writerow(row.values())
    except IOError as e:
        error_dict=[]
        error_dict.append(e)
        writeInFile(error_dict)
    return os.path.abspath(output_file)

def retrieveYoutubeCPDataHasPhone():
    print("creating csv file for youtube has_phone")
    query = '''
            select 
                cp.id
            from
                campaign_profiles cp
            left join youtube_account ya on ya.id = cp.platform_account_id and cp.platform = 'YOUTUBE'
            where ((cp.admin_details->>'phone' IS NOT NULL AND cp.admin_details->>'phone' != '' AND cp.admin_details->>'phone' != 'None') OR
                (cp.admin_details->>'secondaryPhone' IS NOT NULL AND cp.admin_details->>'secondaryPhone' != '' AND cp.admin_details->>'secondaryPhone' != 'None') OR
                (cp.user_details->>'phone' IS NOT NULL AND cp.user_details->>'phone' != '' AND cp.user_details->>'phone' != 'None') or
                (cp.user_details->>'secondaryPhone' IS NOT NULL AND cp.user_details->>'secondaryPhone' != '' AND cp.user_details->>'secondaryPhone' != 'None') OR
                (ya.phone IS NOT NULL AND ya.phone != '')) = True and ya.id >0
            '''
    coffee_postgres_cursor.execute(query)
    result_set = coffee_postgres_cursor.fetchall()
    coffee_postgres_conn.commit()
    try:
        output_file='youtube_data_has_phone.csv'
        with open(output_file, 'w', newline='') as csvfile:
            csv_writer = csv.writer(csvfile)
            for row in result_set:
                csv_writer.writerow(row.values())
    except IOError as e:
        error_dict=[]
        error_dict.append(e)
        writeInFile(error_dict)
    return os.path.abspath(output_file)

def retrieveGccCPDataHasEmail():
    print("creating csv file for GCC has_email")
    query = '''        
            select 
                cp.id
            from
                campaign_profiles cp
            where ((cp.admin_details->>'email' IS NOT NULL AND cp.admin_details->>'email' != '' AND cp.admin_details->>'email' != 'None') OR
                (cp.user_details->>'email' IS NOT NULL AND cp.user_details->>'email' != '' AND cp.user_details->>'email' != 'None')) = True and cp.platform = 'GCC'
            '''
    coffee_postgres_cursor.execute(query)
    result_set = coffee_postgres_cursor.fetchall()
    coffee_postgres_conn.commit()
    try:
        output_file='gcc_has_email.csv'
        with open(output_file, 'w', newline='') as csvfile:
            csv_writer = csv.writer(csvfile)
            for row in result_set:
                csv_writer.writerow(row.values())
    except IOError as e:
        error_dict=[]
        error_dict.append(e)
        writeInFile(error_dict)
    return os.path.abspath(output_file)

def retrieveGccCPDataHasPhone():
    print("creating csv file for gcc has_phone")
    query = '''
            select 
                cp.id
            from
                campaign_profiles cp
            where ((cp.admin_details->>'phone' IS NOT NULL AND cp.admin_details->>'phone' != '' AND cp.admin_details->>'phone' != 'None') OR
                (cp.admin_details->>'secondaryPhone' IS NOT NULL AND cp.admin_details->>'secondaryPhone' != '' AND cp.admin_details->>'secondaryPhone' != 'None') OR
                (cp.user_details->>'phone' IS NOT NULL AND cp.user_details->>'phone' != '' AND cp.user_details->>'phone' != 'None') or
                (cp.user_details->>'secondaryPhone' IS NOT NULL AND cp.user_details->>'secondaryPhone' != '' AND cp.user_details->>'secondaryPhone' != 'None')) = true and cp.platform='GCC'
            '''
    coffee_postgres_cursor.execute(query)
    result_set = coffee_postgres_cursor.fetchall()
    coffee_postgres_conn.commit()
    try:
        output_file='gcc_has_phone.csv'
        with open(output_file, 'w', newline='') as csvfile:
            csv_writer = csv.writer(csvfile)
            for row in result_set:
                csv_writer.writerow(row.values())
    except IOError as e:
        error_dict=[]
        error_dict.append(e)
        writeInFile(error_dict)
    return os.path.abspath(output_file)

def updateCPHasEmail(batch):
    ids = tuple([int(x) for x in flatten(batch)])
    try:
        query = 'UPDATE campaign_profiles SET has_email = True WHERE id IN %s'
        # print(coffee_postgres_cursor.mogrify(query, (ids,)))
        coffee_postgres_cursor.execute(query, (ids,))
    except Exception as e:
        error_dict=[]
        error_dict.append(ids)
        writeInFile(error_dict)

def updateCPHasPhone(batch):
    ids = tuple([int(x) for x in flatten(batch)])
    try:
        query = 'UPDATE campaign_profiles SET has_phone = True WHERE id IN %s'
        # print(coffee_postgres_cursor.mogrify(query, (ids,)))
        coffee_postgres_cursor.execute(query, (ids,))
    except Exception as e:
        error_dict=[]
        error_dict.append(ids)
        writeInFile(error_dict)
    
def flatten(test_list):
    if isinstance(test_list, list):
        temp = []
        for ele in test_list:
            temp.extend(flatten(ele))
        return temp
    else:
        return [test_list]

def split(a, n):
    k, m = divmod(len(a), n)
    return (a[i*k+min(i, m):(i+1)*k+min(i+1, m)] for i in range(n))

def batchProcessingForHasEmail(input_file, max_concurrency):
    completed_size=0
    with open(input_file, 'r', newline='') as csv_file:
        while True:
            csv_reader = csv.reader(csv_file)
            batch_lines = list(itertools.islice(csv_reader, batch_size))
            if not batch_lines:
                break
            batches = split(batch_lines, max_concurrency)
            with concurrent.futures.ThreadPoolExecutor(max_workers=max_concurrency) as executor:
                batch_start_time = time.time()
                futures = [executor.submit(updateCPHasEmail, batch) for batch in batches]
                concurrent.futures.wait(futures)
                batch_end_time = time.time()
                batch_exec_time = math.ceil(batch_end_time - batch_start_time)
                # print("batch exec time = " + str(batch_exec_time))
                print(str(completed_size) + " completed in " + str(batch_exec_time))

def batchProcessingForHasPhone(input_file, max_concurrency):
    completed_size=0
    with open(input_file, 'r', newline='') as csv_file:
        while True:
            csv_reader = csv.reader(csv_file)
            batch_lines = list(itertools.islice(csv_reader, batch_size))
            if not batch_lines:
                break
            batches = split(batch_lines, max_concurrency)
            with concurrent.futures.ThreadPoolExecutor(max_workers=max_concurrency) as executor:
                batch_start_time = time.time()
                futures = [executor.submit(updateCPHasPhone, batch) for batch in batches]
                concurrent.futures.wait(futures)
                completed_size += batch_size
                batch_end_time = time.time()
                batch_exec_time = math.ceil(batch_end_time - batch_start_time)
                # print("batch exec time = " + str(batch_exec_time))
                print(str(completed_size) + " completed in " + str(batch_exec_time))

def process_function():
    start_time = time.time()
    #instagram - has_email
    instagram_has_email_input_file = retrieveInstagramCPDataHasEmail()
    print("instagram_has_email_path = " + instagram_has_email_input_file)
    batchProcessingForHasEmail(instagram_has_email_input_file, max_concurrency)
    # #instagram - has_phone
    instagram_has_phone_input_file = retrieveInstagramCPDataHasPhone()
    print("instagram_has_phone_path = " + instagram_has_phone_input_file)
    batchProcessingForHasPhone(instagram_has_phone_input_file, max_concurrency)
    end_time = time.time()
    exec_time = math.ceil(end_time - start_time)
    print("time for updating instagram CPs = " + str(exec_time))

    start_time = time.time()
    #youtube - has_email
    youtube_has_email_input_file = retrieveYoutubeCPDataHasEmail()
    print("youtube_has_email_path = " + youtube_has_email_input_file)
    batchProcessingForHasEmail(youtube_has_email_input_file, max_concurrency)
    #youtube - has_phone
    youtube_has_phone_input_file = retrieveYoutubeCPDataHasPhone()
    print("youtube_has_phone_path = " + youtube_has_phone_input_file)
    batchProcessingForHasPhone(youtube_has_phone_input_file, max_concurrency)
    end_time = time.time()
    exec_time = math.ceil(end_time - start_time)
    print("time for updating youtube CPs = " + str(exec_time))

    start_time = time.time()
    #gcc - has_email
    gcc_has_email_input_file = retrieveGccCPDataHasEmail()
    print("gcc_has_email_path = " + gcc_has_email_input_file)
    batchProcessingForHasEmail(gcc_has_email_input_file, max_concurrency)
    #gcc - has_phone
    gcc_has_phone_input_file = retrieveGccCPDataHasPhone()
    print("gcc_has_phone_path = " + gcc_has_phone_input_file)
    batchProcessingForHasPhone(gcc_has_phone_input_file, max_concurrency)
    end_time = time.time()
    exec_time = math.ceil(end_time - start_time)
    print("time for updating GCC CPs = " + str(exec_time))

coffee_postgres_conn,coffee_postgres_cursor = postgresConnection()

current_dir = os.path.dirname(os.path.abspath(__file__))
file_path = os.path.join(current_dir, "update_cp_error.csv")
with open(file_path, 'w') as file:
    file.truncate(0)

process_function()
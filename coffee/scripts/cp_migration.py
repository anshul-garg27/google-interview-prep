import psycopg2
from psycopg2 import extras
import os
import traceback
import json
import pandas as pd
import csv
import requests
import time
import math
from datetime import datetime

language_dict = {"english": "en", "assamese": "as", "bengali": "bn", "bodo": "brx", "dogri": "doi", "gujarati": "gu", "hindi": "hi", "kannada": "kn", "kashmiri": "ks", "konkani": "kok", "maithili": "mai", "malayalam": "ml", "manipuri": "mni", "marathi": "mr", "nepali": "ne", "odia": "or", "punjabi": "pa", "sanskrit": "sa", "santali": "sat", "sindhi": "sd", "tamil": "ta", "telugu": "te", "telgu": "te", "urdu": "ur"}

country_dict = {"india":"IN","united states":"US","indonesia":"ID","singapore":"SG","united arab emirates":"AE","spain":"ES","China":"CN","thailand":"TH","pakistan":"PK","germany":"DE","saudi arabia":"SA","canada":"CA"}

def checkIAinCoffee(handle):
    check_instagram_query = 'SELECT * FROM instagram_account WHERE handle = %s'
    coffee_postgres_cursor.execute(check_instagram_query, (handle,))
    existing_row = coffee_postgres_cursor.fetchone()
    return existing_row

def checkYAinCoffee(channelId):
    check_youtube_query = 'Select id from youtube_account where channel_id = %s'
    coffee_postgres_cursor.execute(check_youtube_query, (channelId,))
    existing_row = coffee_postgres_cursor.fetchone()
    return existing_row
    
def insertbeatinInstagramAccount(profile):
    inserted_id = 0
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

    postgres_insert_query = 'INSERT INTO instagram_account (handle, business_id, thumbnail, ig_id, bio, following, followers, name, search_phrase) VALUES (%(handle)s, %(business_id)s, %(thumbnail)s, %(ig_id)s, %(bio)s, %(following)s, %(followers)s, %(name)s, %(search_phrase)s) RETURNING id'
    coffee_postgres_cursor.execute(postgres_insert_query, insertDict)
    inserted_id = coffee_postgres_cursor.fetchone()[0]
    coffee_postgres_conn.commit()
    return inserted_id

def insertbeatinYoutubeAccount(profile):
    insertDict = {}
    insertDict['channel_id'] = profile['channel_id']
    insertDict['title'] = profile['title']
    insertDict['followers'] = profile['subscribers']
    insertDict['uploads_count'] = profile['uploads']
    insertDict['views_count'] = profile['views']
    insertDict['thumbnail'] = profile['thumbnail']
    insertDict['search_phrase'] =  profile['title']

    postgres_insert_query = 'INSERT INTO youtube_account (channel_id, title, followers, uploads_count, views_count, thumbnail, search_phrase) VALUES (%(channel_id)s, %(title)s, %(followers)s, %(uploads_count)s, %(views_count)s, %(thumbnail)s, %(search_phrase)s)  RETURNING id' 
    coffee_postgres_cursor.execute(postgres_insert_query, insertDict)
    inserted_id = coffee_postgres_cursor.fetchone()[0]
    coffee_postgres_conn.commit()
    return inserted_id

def addInstagramDeletedProfile(handle,deleted):
    insert_dict = {}
    insert_dict['handle'] = handle
    insert_dict['ig_id'] = 'MISSING'
    insert_dict['deleted'] = deleted
    postgres_insert_query = 'INSERT INTO instagram_account (handle, ig_id, deleted) VALUES (%(handle)s, %(ig_id)s, %(deleted)s) RETURNING id'
    coffee_postgres_cursor.execute(postgres_insert_query, insert_dict)
    inserted_id = coffee_postgres_cursor.fetchone()[0]
    coffee_postgres_conn.commit()
    return inserted_id

def addYoutubeDeletedProfile(handle,deleted):
    insert_dict = {}
    insert_dict['channel_id'] = handle
    insert_dict['deleted'] = deleted
    postgres_insert_query = 'INSERT INTO youtube_account (channel_id, deleted) VALUES (%(channel_id)s, %(deleted)s) RETURNING id'
    coffee_postgres_cursor.execute(postgres_insert_query, insert_dict)
    inserted_id = coffee_postgres_cursor.fetchone()[0]
    coffee_postgres_conn.commit()
    return inserted_id

def getDataFromBeat(row_dict):    
    data = {}
    if row_dict['platform'] == 'instagram':
        handle = row_dict['instagram_handle']
    elif row_dict['platform'] == 'youtube':
        handle = row_dict['youtube_channel_url']
    url = "http://beat.goodcreator.co/profiles/"+ row_dict['platform'] +"/byhandle/"+ handle
    response = requests.get(url)
    response.raise_for_status()
    data = response.json()
    return data
    
def readCsvFileinArray(filename):
    data = []
    current_dir = os.path.dirname(os.path.abspath(__file__))
    filepath = os.path.join(current_dir, filename)
    with open(filepath, 'r') as csv_file:
        csv_reader = csv.reader(csv_file)
        for row in csv_reader:
            data.append(row[0])
    return data

def getValuesBykey(key,cpRow,identityRow):
    
    if key in identityRow:
        value = identityRow[key]
    elif key in cpRow:
        value = cpRow[key]
    else:
        value = None
    return value

def checkCampaignProfileinCoffee(id):
    
    check_youtube_query = 'Select id from campaign_profiles where id = %s'
    coffee_postgres_cursor.execute(check_youtube_query, (id,))
    existing_row = coffee_postgres_cursor.fetchone()
    return existing_row

def readYoutubeCPS(filename,key):

    current_dir = os.path.dirname(os.path.abspath(__file__))
    filepath = os.path.join(current_dir, filename)
    data = pd.read_csv(filepath, low_memory=False)
    youtubeDict = {}
    youtubeGroupDict = {}
    for  index,row in data.iterrows():
       
        handle = row[key]
        group_id = row['group_id']
        row_dict = row.to_dict()
        # Filter out NaN values from the row dictionary
        row_dict = {k: v for k, v in row_dict.items() if not pd.isna(v)}
        youtubeDict[handle] = row_dict
        youtubeGroupDict[group_id] = handle
    return youtubeDict,youtubeGroupDict

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

# def writeInFile(data_array):
#     current_dir = os.path.dirname(os.path.abspath(__file__))
#     file_path = os.path.join(current_dir, "cp_migration_error.txt")
#     with open(file_path, "a") as f:
#         line = ",".join(data_array)
#         f.write(line + "\n")

def writeInFile(data_array):
    current_dir = os.path.dirname(os.path.abspath(__file__))
    file_path = os.path.join(current_dir, "cp_migration_error.csv")
    
    with open(file_path, "a", newline='') as f:
        writer = csv.writer(f)
        writer.writerow(data_array)

def makeCpEntry(cpRow,platformRow,identityDict):
    identityRow = {}
    platform = ""
    if 'youtube_channel_url' in cpRow:
        platform = "YOUTUBE"
        key = cpRow['youtube_channel_url']
        if key in identityDict:
            identityRow = identityDict[key]
    elif 'instagram_handle' in cpRow:
        platform = "INSTAGRAM"
        key = cpRow['instagram_handle']
        if key in identityDict:
            identityRow = identityDict[key]
    admin_details = {}
    if  getValuesBykey('name',cpRow,identityRow) is not None:
        admin_details['name'] =  getValuesBykey('name',cpRow,identityRow)

    if  getValuesBykey('email',cpRow,identityRow) is not None:
        admin_details['email'] = getValuesBykey('email',cpRow,identityRow)

    if  getValuesBykey('bio',cpRow,identityRow) is not None:
        admin_details['bio'] = getValuesBykey('bio',cpRow,identityRow)
    
    if  getValuesBykey('dob',cpRow,identityRow) is not None:
        admin_details['dob'] = getValuesBykey('dob',cpRow,identityRow)
    
    if  getValuesBykey('city',cpRow,identityRow) is not None:
        admin_details['city'] = getValuesBykey('city',cpRow,identityRow)

    if  getValuesBykey('state',cpRow,identityRow) is not None:
        admin_details['state'] = getValuesBykey('state',cpRow,identityRow)

    if  getValuesBykey('country',cpRow,identityRow) is not None:
        admin_details['country'] = getValuesBykey('country',cpRow,identityRow)
        country = getValuesBykey('country',cpRow,identityRow)
        country = country.lower()
        if country in country_dict:
            admin_details['countryCode'] = country_dict[country]

    if  getValuesBykey('location',cpRow,identityRow) is not None:
        location = getValuesBykey('location',cpRow,identityRow)
        locationArray = [x.strip() for x in location.split(",")]
        admin_details['location'] = locationArray

    if  getValuesBykey('primary_phone',cpRow,identityRow) is not None:
        phone = getValuesBykey('primary_phone',cpRow,identityRow)
        admin_details['phone'] = str(phone)
    
    if  getValuesBykey('secondaryPhone',cpRow,identityRow) is not None:
        secondaryPhone = getValuesBykey('secondary_phone',cpRow,identityRow)
        admin_details['secondaryPhone'] = secondaryPhone

    if  getValuesBykey('gender',cpRow,identityRow) is not None:
        admin_details['gender'] = getValuesBykey('gender',cpRow,identityRow)

    if  getValuesBykey('languages',cpRow,identityRow) is not None:
        languageString = getValuesBykey('languages',cpRow,identityRow)
        languages = languageString.strip("{}").split(",")
        language_codes = [language_dict.get(lang.lower()) for lang in languages if lang.lower() in language_dict]
        if len(language_codes) >0:
            admin_details['languages'] = language_codes

    if  getValuesBykey('is_blacklisted',cpRow,identityRow) is not None:
        admin_details['isBlacklisted'] = getValuesBykey('is_blacklisted',cpRow,identityRow)

    if  getValuesBykey('blacklisted_by',cpRow,identityRow) is not None:
        admin_details['blacklistedBy'] = getValuesBykey('blacklisted_by',cpRow,identityRow)
    
    if  getValuesBykey('blacklisted_reason',cpRow,identityRow) is not None:
        admin_details['blacklistedReason'] = getValuesBykey('blacklisted_reason',cpRow,identityRow)

    if 'category_tags' in cpRow and cpRow['category_tags'] is not None:
        admin_details['campaignCategories'] = json.loads(cpRow['category_tags'])

    if  getValuesBykey('category_ids_list',cpRow,identityRow) is not None:
        set = getValuesBykey('category_ids_list',cpRow,identityRow)
        set = set.strip("{}")
        categoryIdList = [x.strip() for x in set.split(",")]
        admin_details['campaignCategoryIds'] = categoryIdList
    
    if 'creator_program_list' in identityRow and identityRow['creator_program_list'] is not None:
        creatorPrograms = identityRow['creator_program_list']
        creatorProgramArray = json.loads(creatorPrograms)
        newCreatorProgramArray = []
        for item in creatorProgramArray:
            tag = item["tag"]
            level = item["level"]
            if tag == "BEAUTY":
                final = str(1)+"_"+level.upper()
            elif tag == "GOOD_PARENTING":
                final = str(2)+"_"+level.upper()
            elif tag == "GOOD_LIFE":
                final = str(3)+"_"+level.upper()
            newCreatorProgramArray.append(final)
        if len(newCreatorProgramArray) > 0:
            admin_details['creatorPrograms'] = newCreatorProgramArray
    elif 'creator_program_list' in cpRow and cpRow['creator_program_list'] is not None and cpRow['creator_program_list'] != '{}':
        string_data = cpRow['creator_program_list']
        string_data = string_data.strip("{}")
        newCreatorProgramArray = [x.strip() for x in string_data.split(",")]
        admin_details['creatorPrograms'] = newCreatorProgramArray
    
    if 'creator_cohort_list' in cpRow and cpRow['creator_cohort_list'] is not None and cpRow['creator_cohort_list'] != '{}':
        string_data = cpRow['creator_cohort_list']
        string_data = string_data.strip("{}")
        newCreatorProgramArray = [x.strip() for x in string_data.split(",")]
        admin_details['creatorCohorts'] = newCreatorProgramArray

    if 'whatsapp_enabled' in identityRow:
        if identityRow['whatsapp_enabled'] == 1:
            admin_details['WhatsappOptIn'] = True
        else:
            admin_details['WhatsappOptIn'] = False
    elif 'whatsapp_approved' in cpRow:
        admin_details['WhatsappOptIn'] = cpRow['whatsapp_approved']


    user_details = {}
    if  getValuesBykey('email',cpRow,identityRow) is not None:
        user_details['email'] =  getValuesBykey('email',cpRow,identityRow)
    
    if  getValuesBykey('primary_phone',cpRow,identityRow) is not None:
        phone = getValuesBykey('primary_phone',cpRow,identityRow)
        user_details['phone'] = str(phone)

    if  getValuesBykey('secondary_phone',cpRow,identityRow) is not None:
        secondaryPhone = getValuesBykey('secondary_phone',cpRow,identityRow)
        user_details['secondaryPhone'] = secondaryPhone

    if  getValuesBykey('name',cpRow,identityRow) is not None:
        user_details['name'] =  getValuesBykey('name',cpRow,identityRow)

    if  getValuesBykey('dob',cpRow,identityRow) is not None:
        user_details['dob'] = getValuesBykey('dob',cpRow,identityRow)
    
    if  getValuesBykey('gender',cpRow,identityRow) is not None:
        user_details['gender'] = getValuesBykey('gender',cpRow,identityRow)
    
    if 'whatsapp_enabled' in identityRow:
        if identityRow['whatsapp_enabled'] == 1:
            user_details['WhatsappOptIn'] = True
        else:
            user_details['WhatsappOptIn'] = False
    elif 'whatsapp_approved' in cpRow:
        user_details['WhatsappOptIn'] = cpRow['whatsapp_approved']

    if  getValuesBykey('category_ids_list',cpRow,identityRow) is not None:
        set = getValuesBykey('category_ids_list',cpRow,identityRow)
        set = set.strip("{}")
        categoryIdList = [x.strip() for x in set.split(",")]
        user_details['campaignCategoryIds'] = categoryIdList
    
    if  getValuesBykey('languages',cpRow,identityRow) is not None:
        languageString = getValuesBykey('languages',cpRow,identityRow)
        languages = languageString.strip("{}").split(",")
        language_codes = [language_dict.get(lang.lower()) for lang in languages if lang.lower() in language_dict]
        if len(language_codes) >0:
            user_details['languages'] = language_codes
    
    if  getValuesBykey('location',cpRow,identityRow) is not None:
        location = getValuesBykey('location',cpRow,identityRow)
        locationArray = [x.strip() for x in location.split(",")]
        user_details['location'] = locationArray
    
    if  getValuesBykey('bio',cpRow,identityRow) is not None:
        user_details['bio'] = getValuesBykey('bio',cpRow,identityRow)

    if  getValuesBykey('member_id',cpRow,identityRow) is not None:
        user_details['memberId'] = getValuesBykey('member_id',cpRow,identityRow)

    if  getValuesBykey('reference_code',cpRow,identityRow) is not None:
        user_details['referenceCode'] = getValuesBykey('reference_code',cpRow,identityRow)
    
    if  getValuesBykey('ekyc_pending',cpRow,identityRow) is not None:
        user_details['ekycPending'] = getValuesBykey('ekyc_pending',cpRow,identityRow)

    if  getValuesBykey('webengage_user_id',cpRow,identityRow) is not None:
        user_details['webengageUserId'] = str(getValuesBykey('webengage_user_id',cpRow,identityRow))
    
    if  getValuesBykey('notification_device_id',cpRow,identityRow) is not None:
        user_details['notificationToken'] = getValuesBykey('notification_device_id',cpRow,identityRow)
    
    if  getValuesBykey('instant_gratification_invited',cpRow,identityRow) is not None:
        user_details['instantGratificationInvited'] = getValuesBykey('instant_gratification_invited',cpRow,identityRow)

    if  getValuesBykey('amazon_store_link_verified',cpRow,identityRow) is not None:
        user_details['amazonStoreLinkVerified'] = getValuesBykey('amazon_store_link_verified',cpRow,identityRow)
    
    insert_dict = {}
    insert_dict["id"] = cpRow['id']
    insert_dict["platform"] = platform
    insert_dict["platform_account_id"] = platformRow['id']
    insert_dict["admin_details"] =json.dumps(admin_details) 
    if bool(user_details):
        insert_dict['user_details'] = json.dumps(user_details)
    else:
        insert_dict['user_details'] = None
    insert_dict['on_gcc'] = cpRow['on_gcc'] if 'on_gcc' in cpRow else None
    insert_dict['on_gcc_app'] = cpRow['on_gcc_app'] if 'on_gcc_app' in cpRow else None
    insert_dict['gcc_user_account_id'] = identityRow['account_id'] if 'account_id' in identityRow else None
    
    postgres_insert_query = 'INSERT INTO campaign_profiles (id, platform, platform_account_id, admin_details, user_details, on_gcc, on_gcc_app, gcc_user_account_id) VALUES (%(id)s, %(platform)s, %(platform_account_id)s, %(admin_details)s, %(user_details)s, %(on_gcc)s, %(on_gcc_app)s, %(gcc_user_account_id)s)'
    coffee_postgres_cursor.execute(postgres_insert_query, insert_dict)
    coffee_postgres_conn.commit()
        

def insertIA(instagram_handle,ia_beat_call_cp_count):
    insertId = 0
    existing_row = checkIAinCoffee(instagram_handle)
    if existing_row is not None:
        insertId = existing_row['id']
    elif instagram_handle in deletedHandles:
        insertId = addInstagramDeletedProfile(instagram_handle,True)
    elif instagram_handle in invalidHandles:
        insertId = addInstagramDeletedProfile(instagram_handle,False)
    else:
        ia_beat_call_cp_count = ia_beat_call_cp_count + 1
        row_dict = {}
        row_dict['platform'] = 'instagram'
        row_dict['instagram_handle'] = instagram_handle
        beatData = getDataFromBeat(row_dict)
        if bool(beatData) and 'profile' in beatData:
            profile = beatData['profile']
            if 'handle' in profile:
                insertId = insertbeatinInstagramAccount(profile)
            else:
                insertId = addInstagramDeletedProfile(row_dict['instagram_handle'],True)
        else:
            insertId = addInstagramDeletedProfile(row_dict['instagram_handle'],True)
    return insertId,ia_beat_call_cp_count

def updateYAProfileId(gcc_profile_id,gcc_linked_instagram_handle,channelId):
    updated_at = datetime.fromtimestamp(time.time())
    check_mapping_query = "Select id from youtube_account where gcc_linked_instagram_handle = %s"
    coffee_postgres_cursor.execute(check_mapping_query, (gcc_linked_instagram_handle,))
    existing_row = coffee_postgres_cursor.fetchone()

    if existing_row is not None:
        youtube_profile_id = 'YA_' + str(existing_row['id'])
        update_query = 'UPDATE youtube_account SET gcc_profile_id = %s, gcc_linked_instagram_handle = %s, updated_at = %s WHERE id = %s'
        coffee_postgres_cursor.execute(update_query, (youtube_profile_id, None, updated_at, existing_row['id']))
        coffee_postgres_conn.commit()

        query = 'UPDATE youtube_account SET gcc_profile_id = %s, updated_at = %s, gcc_linked_instagram_handle = %s WHERE channel_id = %s'
        coffee_postgres_cursor.execute(query, (gcc_profile_id, updated_at, gcc_linked_instagram_handle, channelId))
        coffee_postgres_conn.commit()

    else:
        query = 'UPDATE youtube_account SET gcc_profile_id = %s, updated_at = %s, gcc_linked_instagram_handle = %s WHERE channel_id = %s'
        coffee_postgres_cursor.execute(query, (gcc_profile_id, updated_at, gcc_linked_instagram_handle, channelId))
        coffee_postgres_conn.commit()

def instagramCps():
    insta_left = 0
    insta_total = 0
    ia_beat_call_cp_count = 0
    
    start_time = time.time()
    for key, row_dict in instagramDict.items():
        insertId = 0
        channelID = ''
        try:
            if 'instagram_handle' in row_dict and row_dict['instagram_handle'] is not None:
                row_dict['instagram_handle'] = row_dict['instagram_handle'].lower()
                check = checkCampaignProfileinCoffee(row_dict['id'])
                if check is None:
                    insertId,ia_beat_call_cp_count = insertIA(row_dict['instagram_handle'],ia_beat_call_cp_count)
                    if insertId:
                        if 'group_id' in row_dict and row_dict['group_id'] is not None and row_dict['group_id'] in youtubeGroupDict:
                            channelID = youtubeGroupDict[row_dict['group_id']]
                            _,ia_beat_call_cp_count= insertYA(channelID,ia_beat_call_cp_count)
                            profile_id = 'IA_'+ str(insertId)
                            updateYAProfileId(profile_id,row_dict['instagram_handle'],channelID)
                        coffee_ia_row = {'id': insertId}
                        makeCpEntry(row_dict,coffee_ia_row,identityDict)
                    else:
                        insta_left = insta_left + 1
                        
                        writer = []
                        writer.append(str(row_dict['id']))
                        # if 'instagram_handle' in row_dict:
                        #     row_dict['instagram_handle'] = row_dict['instagram_handle'][:50]
                        #     writer.append(row_dict['instagram_handle'])
                        # if 'channelID' in locals() and channelID is not None:
                        #     writer.append(channelID)
                        writeInFile(writer)
            else:
                insta_left = insta_left + 1
                
                writer = []
                writer.append(str(row_dict['id']))
                # if 'instagram_handle' in row_dict:
                #     row_dict['instagram_handle'] = row_dict['instagram_handle'][:50]
                #     writer.append(row_dict['instagram_handle'])
                # if 'channelID' in locals() and channelID is not None:
                #     writer.append(channelID)
                writeInFile(writer)

        except Exception as err:
            traceback.print_exc() 
            coffee_postgres_conn.commit()
            insta_left = insta_left + 1
            writer = []
            writer.append(str(row_dict['id']))
            # if 'instagram_handle' in row_dict:
            #     row_dict['instagram_handle'] = row_dict['instagram_handle'][:50]
            #     writer.append(row_dict['instagram_handle'])
            # if 'channelID' in locals() and channelID is not None:
            #     writer.append(channelID)
            writeInFile(writer)

        insta_total = insta_total+ 1
        if insta_total > 0 and (insta_total % chunk_size == 0 or insta_total == len(instagramDict)):
            end_time = time.time()
            exec_time = math.ceil(end_time - start_time)
            print(str(insta_total)+ " instagram CP Completed and " + str(insta_left) +" cps left in this "+ str(chunk_size)+" chunk in " + str(exec_time) + " seconds")
            print(str(ia_beat_call_cp_count)+ " BeaT calls")
            # insta_left = 0
            start_time = time.time()

def insertYA(channelId,yt_beat_call_cp_count):
    insertId = 0
    existing_row = checkYAinCoffee(channelId)
    if existing_row is not None:
        insertId = existing_row['id']
    elif channelId in deletedHandles:
        insertId = addYoutubeDeletedProfile(channelId,True)
    elif channelId in invalidHandles:
        insertId = addYoutubeDeletedProfile(channelId,False)
    else:
        yt_beat_call_cp_count = yt_beat_call_cp_count + 1
        row_dict = {}
        row_dict['youtube_channel_url'] = channelId
        row_dict['platform'] = 'youtube'
        beatData = getDataFromBeat(row_dict)
        if bool(beatData) and 'profile' in beatData:
            profile = beatData['profile']
            if 'channel_id' in profile:
                insertId = insertbeatinYoutubeAccount(profile)
            else:
                insertId = addYoutubeDeletedProfile(channelId,True)
        else:
            insertId = addYoutubeDeletedProfile(channelId,True)
    return insertId,yt_beat_call_cp_count

def youtubeCps():
    yt_left = 0
    yt_total = 0
    yt_beat_call_cp_count = 0
    start_time = time.time()
    for key, value in youtubeDict.items():
        insertId = 0
        row_dict = value
        try:
            if 'youtube_channel_url' in row_dict and row_dict['youtube_channel_url'] is not None:
                check = checkCampaignProfileinCoffee(row_dict['id'])
                if check is None:
                    insertId,yt_beat_call_cp_count = insertYA(row_dict['youtube_channel_url'],yt_beat_call_cp_count)
                    if insertId:
                        coffee_ia_row = {'id': insertId}
                        makeCpEntry(row_dict,coffee_ia_row,identityDict)
                    else:
                        insta_left = insta_left + 1

                        writer = []
                        writer.append(str(row_dict['id']))
                        # if 'youtube_channel_url' in row_dict:
                        #     writer.append(row_dict['youtube_channel_url'])
                        writeInFile(writer)
            else:
                yt_left = yt_left + 1
                writer = []
                writer.append(str(row_dict['id']))
                # if 'youtube_channel_url' in row_dict:
                #     writer.append(row_dict['youtube_channel_url'])
                writeInFile(writer)
        except Exception as err:
            traceback.print_exc() 
            coffee_postgres_conn.commit()
            yt_left = yt_left + 1
            writer = []
            writer.append(str(row_dict['id']))
            # if 'youtube_channel_url' in row_dict:
            #     writer.append(row_dict['youtube_channel_url'])
            writeInFile(writer)

        yt_total = yt_total+ 1
        if  yt_total >0 and  (yt_total % chunk_size == 0 or yt_total == len(youtubeDict)):
            end_time = time.time()
            exec_time = math.ceil(end_time - start_time)
            print(str(yt_total)+ " youtube CP Completed and " + str(yt_left) +" cps left in this "+str(chunk_size)+" chunk in " + str(exec_time) + " seconds")
            print(str(yt_beat_call_cp_count)+ " BeaT calls")
            # yt_left = 0
            start_time = time.time()

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
print("Reding CSV file...RELAX")
deletedHandles = readCsvFileinArray("deleted_handle.csv")
print("Deleted Handle CSV Loaded...")

invalidHandles = readCsvFileinArray("useless_handle.csv")
print("Useless Handle CSV Loaded...")

identityDict = readCsvFile("identity.csv","handle")
print("Identity CSV Loaded...")

instagramDict = readCsvFile("winkl_ia.csv","instagram_handle")
print("Instagram CP CSV Loaded...")

youtubeDict,youtubeGroupDict = readYoutubeCPS("winkl_yt.csv","youtube_channel_url")
print("Youtube CP CSV Loaded...")


chunk_size = 1000
# MAKE FILE EMPTY ON START
current_dir = os.path.dirname(os.path.abspath(__file__))
file_path = os.path.join(current_dir, "cp_migration_error.csv")
with open(file_path, 'w') as file:
    file.truncate(0)
print("Running For Instgram CPS")
writeInFile(["INSTAGRAM"])
instagramCps()
print("Running For Youtube CPS")
writeInFile(["YOUTUBE"])
youtubeCps()

query = 'UPDATE instagram_account SET gcc_profile_id = Concat(\'IA_\',id) where gcc_profile_id IS NULL'
coffee_postgres_cursor.execute(query)
print("instagram_account TABLE column gcc_profile_id GOT UPDATED")
coffee_postgres_conn.commit()

query = 'UPDATE youtube_account SET gcc_profile_id = Concat(\'YA_\',id) where gcc_profile_id IS NULL'
coffee_postgres_cursor.execute(query)
print("youtube_account TABLE column gcc_profile_id GOT UPDATED")
coffee_postgres_conn.commit()

query = 'UPDATE instagram_account ia SET gcc_linked_channel_id = ya.channel_id FROM youtube_account ya WHERE ia.gcc_profile_id = ya.gcc_profile_id'
coffee_postgres_cursor.execute(query)
print("instagram_account TABLE column linked_channel_id GOT UPDATED")
coffee_postgres_conn.commit()
coffee_postgres_conn.close()

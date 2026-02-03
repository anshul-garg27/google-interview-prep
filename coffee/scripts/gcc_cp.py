import os
import psycopg2
from psycopg2 import extras
import pandas as pd
import json
import traceback




language_dict = {"english": "en", "assamese": "as", "bengali": "bn", "bodo": "brx", "dogri": "doi", "gujarati": "gu", "hindi": "hi", "kannada": "kn", "kashmiri": "ks", "konkani": "kok", "maithili": "mai", "malayalam": "ml", "manipuri": "mni", "marathi": "mr", "nepali": "ne", "odia": "or", "punjabi": "pa", "sanskrit": "sa", "santali": "sat", "sindhi": "sd", "tamil": "ta", "telugu": "te", "telgu": "te", "urdu": "ur"}
country_dict = {"india":"IN","united states":"US","indonesia":"ID","singapore":"SG","united arab emirates":"AE","spain":"ES","China":"CN","thailand":"TH","pakistan":"PK","germany":"DE","saudi arabia":"SA","canada":"CA"}

def readCsvFile(filename,key):
    result_dict = {}
    current_dir = os.path.dirname(os.path.abspath(__file__))
    filepath = os.path.join(current_dir, filename)
    data = pd.read_csv(filepath, low_memory=False)
    for index, row in data.iterrows():
        handle = row[key]
        row_dict = row.dropna().to_dict()
        result_dict.update({handle: row_dict})
    return result_dict

def gccCps():
    gcc_left = 0
    gcc_total = 0
    for key, cpRow in gccCpDict.items():
        try:
            admin_details = {}
            if 'name' in cpRow is not None:
                admin_details['name'] =  cpRow['name']

            if 'email' in cpRow is not None:
                admin_details['email'] = cpRow['email']

            if 'bio' in cpRow is not None:
                admin_details['bio'] = cpRow['bio']
            
            if  'dob' in cpRow is not None:
                admin_details['dob'] = cpRow['dob']
            
            if  'city' in cpRow is not None:
                admin_details['city'] = cpRow['city']

            if  'state' in cpRow is not None:
                admin_details['state'] = cpRow['state']

            if  'country' in cpRow is not None:
                admin_details['country'] = cpRow['country']
                country = cpRow['country']
                country = country.lower()
                if country in country_dict:
                    admin_details['countryCode'] = country_dict[country]
            
            if  'location' in cpRow is not None:
                location = cpRow['location']
                locationArray = [x.strip() for x in location.split(",")]
                admin_details['location'] = locationArray

            if  'primary_phone' in cpRow is not None:
                admin_details['phone'] = cpRow['primary_phone']
            
            if  'secondaryPhone' in cpRow is not None:
                admin_details['secondaryPhone'] = cpRow['secondary_phone']

            if  'gender' in cpRow is not None:
                admin_details['gender'] = cpRow['gender']

            if  'languages' in cpRow is not None:
                languageString = cpRow['languages']
                languages = languageString.strip("{}").split(",")
                language_codes = [language_dict.get(lang.lower()) for lang in languages if lang.lower() in language_dict]
                if len(language_codes) >0:
                    admin_details['languages'] = language_codes

            if  'is_blacklisted' in cpRow is not None:
                admin_details['isBlacklisted'] = cpRow['is_blacklisted']
            
            if 'whatsapp_approved' in cpRow is not None:
                admin_details['whatsappOptin'] = cpRow['whatsapp_approved']

            if 'category_tags' in cpRow is not None:
                admin_details['campaignCategories'] = json.loads(cpRow['category_tags'])

            if  'category_ids_list' in cpRow is not None:
                set = cpRow['category_ids_list']
                set = set.strip("{}")
                if set != '':
                    categoryIdList = [x.strip() for x in set.split(",")]
                    admin_details['campaignCategoryIds'] = categoryIdList
            
            if 'creator_program_list' in cpRow and cpRow['creator_program_list'] is not None and cpRow['creator_program_list'] != '{}':
                string_data = cpRow['creator_program_list']
                string_data = string_data.strip("{}")
                newCreatorProgramArray = [x.strip() for x in string_data.split(",")]
                admin_details['creatorPrograms'] = newCreatorProgramArray
            
            user_details = {}

            if  'email' in cpRow is not None:
                user_details['email'] =  cpRow['email']
        
            if  'primary_phone' in cpRow is not None:
                user_details['phone'] =  cpRow['primary_phone']

            if  'secondary_phone' in cpRow is not None:
                user_details['secondaryPhone'] =  cpRow['secondary_phone']

            if  'name' in cpRow is not None:
                user_details['name'] =  cpRow['name']

            if  'dob' in cpRow is not None:
                user_details['dob'] = cpRow['dob']
            
            if  'gender' in cpRow is not None:
                user_details['gender'] = cpRow['gender']
            
            if 'whatsapp_approved' in cpRow:
                user_details['WhatsappOptIn'] = cpRow['whatsapp_approved']
            
            if  'category_ids_list' in cpRow is not None:
                set = cpRow['category_ids_list']
                set = set.strip("{}")
                categoryIdList = [x.strip() for x in set.split(",")]
                user_details['campaignCategoryIds'] = categoryIdList

            if  'languages' in cpRow  is not None:
                languageString = cpRow['languages']
                languages = languageString.strip("{}").split(",")
                language_codes = [language_dict.get(lang.lower()) for lang in languages if lang.lower() in language_dict]
                if len(language_codes) >0:
                    user_details['languages'] = language_codes
        
            if  'location' in cpRow is not None:
                location = cpRow['location']
                locationArray = [x.strip() for x in location.split(",")]
                user_details['location'] = locationArray
            
            if  'bio' in cpRow is not None:
                user_details['bio'] = cpRow['bio']

            if  'member_id' in cpRow is not None:
                user_details['memberId'] = cpRow['member_id']

            if  'reference_code' in cpRow is not None:
                user_details['referenceCode'] = cpRow['reference_code']
            
            if  'ekyc_pending' in cpRow is not None:
                user_details['ekycPending'] = cpRow['ekyc_pending']

            if  'webengage_user_id' in cpRow is not None:
                user_details['webengageUserId'] = cpRow['webengage_user_id']
            
            if  'notification_device_id' in cpRow is not None:
                user_details['notificationToken'] = cpRow['notification_device_id']
            
            if 'instant_gratification_invited' in cpRow is not None:
                user_details['instantGratificationInvited'] = cpRow['instant_gratification_invited']

            if  'amazon_store_link_verified' in cpRow is not None:
                user_details['amazonStoreLinkVerified'] = cpRow['amazon_store_link_verified']

            check_youtube_query = 'Select id from campaign_profiles where id = %s'
            coffee_postgres_cursor.execute(check_youtube_query, (cpRow['id'],))
            existing_row = coffee_postgres_cursor.fetchone()
            if existing_row is None:
                insert_dict = {}
                insert_dict["id"] = cpRow['id']
                insert_dict["platform"] = 'GCC'
                insert_dict["platform_account_id"] = cpRow['bulbul_user_account_id']
                insert_dict["admin_details"] =json.dumps(admin_details) 
                insert_dict['user_details'] = json.dumps(user_details)
                insert_dict['on_gcc'] = False
                insert_dict['on_gcc_app'] = False
                insert_dict['gcc_user_account_id'] = cpRow['bulbul_user_account_id']
                
                postgres_insert_query = 'INSERT INTO campaign_profiles (id, platform, platform_account_id, admin_details, user_details, on_gcc, on_gcc_app, gcc_user_account_id) VALUES (%(id)s, %(platform)s, %(platform_account_id)s, %(admin_details)s, %(user_details)s, %(on_gcc)s, %(on_gcc_app)s,%(gcc_user_account_id)s)'
                coffee_postgres_cursor.execute(postgres_insert_query, insert_dict)
                coffee_postgres_conn.commit()
        except Exception as err:
            coffee_postgres_conn.commit()
            gcc_left = gcc_left + 1
            print(err)
            # traceback.print_exc()
            # print(cpRow)
            # exit()

        gcc_total = gcc_total+ 1
        if  gcc_total >0 and  (gcc_total % chunk_size == 0 or gcc_total == len(gccCpDict)):
            print(str(gcc_total)+ " GCC CP Completed and " + str(gcc_left) +" cps left in this "+str(chunk_size)+" chunk")
            # gcc_left = 0

print("Reading GCC CSV File")
gccCpDict = readCsvFile("gcc_cp.csv","bulbul_user_account_id")

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

chunk_size = 1000

print("Running For GCC CPS")
gccCps()
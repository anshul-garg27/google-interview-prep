from clickhouse_driver import Client
import boto3
import traceback
import json
import os
from datetime import datetime,timedelta
import time
import psycopg2
from psycopg2 import extras
import math
import sys

indian_states = ["Andhra Pradesh","Arunachal Pradesh","Assam","Bihar","Chhattisgarh","Goa","Gujarat","Haryana","Himachal Pradesh","Jharkhand","Karnataka","Kerala","Madhya Pradesh","Maharashtra","Manipur","Meghalaya","Mizoram","Nagaland","Odisha","Punjab","Rajasthan","Sikkim","Tamil Nadu","Telangana","Tripura","Uttar Pradesh","Uttarakhand","West Bengal","Jammu & Kashmir"]

country_data = {"AF":"Afghanistan","AX":"Åland Islands","AL":"Albania","DZ":"Algeria","AS":"American Samoa","AD":"Andorra","AO":"Angola","AI":"Anguilla","AQ":"Antarctica","AG":"Antigua and Barbuda","AR":"Argentina","AM":"Armenia","AW":"Aruba","AU":"Australia","AT":"Austria","AZ":"Azerbaijan","BS":"Bahamas","BH":"Bahrain","BD":"Bangladesh","BB":"Barbados","BY":"Belarus","BE":"Belgium","BZ":"Belize","BJ":"Benin","BM":"Bermuda","BT":"Bhutan","BO":"Bolivia","BQ":"Bonaire, Sint Eustatius and Saba","BA":"Bosnia and Herzegovina","BW":"Botswana","BV":"Bouvet Island","BR":"Brazil","IO":"British Indian Ocean Territory","BN":"Brunei Darussalam","BG":"Bulgaria","BF":"Burkina Faso","BI":"Burundi","CV":"Cabo Verde","KH":"Cambodia","CM":"Cameroon","CA":"Canada","KY":"Cayman Islands","CF":"Central African Republic","TD":"Chad","CL":"Chile","CN":"China","CX":"Christmas Island","CC":"Cocos (Keeling) Islands","CO":"Colombia","KM":"Comoros","CG":"Congo","CD":"Congo, Democratic Republic of the","CK":"Cook Islands","CR":"Costa Rica","CI":"Côte d'Ivoire","HR":"Croatia","CU":"Cuba","CW":"Curaçao","CY":"Cyprus","CZ":"Czech Republic","DK":"Denmark","DJ":"Djibouti","DM":"Dominica","DO":"Dominican Republic","EC":"Ecuador","EG":"Egypt","SV":"El Salvador","GQ":"Equatorial Guinea","ER":"Eritrea","EE":"Estonia","SZ":"Eswatini","ET":"Ethiopia","FK":"Falkland Islands (Malvinas)","FO":"Faroe Islands","FJ":"Fiji","FI":"Finland","FR":"France","GF":"French Guiana","PF":"French Polynesia","TF":"French Southern Territories","GA":"Gabon","GM":"Gambia","GE":"Georgia","DE":"Germany","GH":"Ghana","GI":"Gibraltar","GR":"Greece","GL":"Greenland","GD":"Grenada","GP":"Guadeloupe","GU":"Guam","GT":"Guatemala","GG":"Guernsey","GN":"Guinea","GW":"Guinea-Bissau","GY":"Guyana","HT":"Haiti","HM":"Heard Island and McDonald Islands","VA":"Holy See","HN":"Honduras","HK":"Hong Kong","HU":"Hungary","IS":"Iceland","IN":"India","ID":"Indonesia","IR":"Iran, Islamic Republic of","IQ":"Iraq","IE":"Ireland","IM":"Isle of Man","IL":"Israel","IT":"Italy","JM":"Jamaica","JP":"Japan","JE":"Jersey","JO":"Jordan","KZ":"Kazakhstan","KE":"Kenya","KI":"Kiribati","KP":"Korea, Democratic People's Republic of","KR":"Korea, Republic of","KW":"Kuwait","KG":"Kyrgyzstan","LA":"Lao People's Democratic Republic","LV":"Latvia","LB":"Lebanon","LS":"Lesotho","LR":"Liberia","LY":"Libya","LI":"Liechtenstein","LT":"Lithuania","LU":"Luxembourg","MO":"Macao","MG":"Madagascar","MW":"Malawi","MY":"Malaysia","MV":"Maldives","ML":"Mali","MT":"Malta","MH":"Marshall Islands","MQ":"Martinique","MR":"Mauritania","MU":"Mauritius","YT":"Mayotte","MX":"Mexico","FM":"Micronesia, Federated States of","MD":"Moldova, Republic of","MC":"Monaco","MN":"Mongolia","ME":"Montenegro","MS":"Montserrat","MA":"Morocco","MZ":"Mozambique","MM":"Myanmar","NA":"Namibia","NR":"Nauru","NP":"Nepal","NL":"Netherlands","NC":"New Caledonia","NZ":"New Zealand","NI":"Nicaragua","NE":"Niger","NG":"Nigeria","NU":"Niue","NF":"Norfolk Island","MK":"North Macedonia","MP":"Northern Mariana Islands","NO":"Norway","OM":"Oman","PK":"Pakistan","PW":"Palau","PS":"Palestine, State of","PA":"Panama","PG":"Papua New Guinea","PY":"Paraguay","PE":"Peru","PH":"Philippines","PN":"Pitcairn","PL":"Poland","PT":"Portugal","PR":"Puerto Rico","QA":"Qatar","RE":"Réunion","RO":"Romania","RU":"Russian Federation","RW":"Rwanda","BL":"Saint Barthélemy","SH":"Saint Helena, Ascension and Tristan da Cunha","KN":"Saint Kitts and Nevis","LC":"Saint Lucia","MF":"Saint Martin (French part)","PM":"Saint Pierre and Miquelon","VC":"Saint Vincent and the Grenadines","WS":"Samoa","SM":"San Marino","ST":"Sao Tome and Principe","SA":"Saudi Arabia","SN":"Senegal","RS":"Serbia","SC":"Seychelles","SL":"Sierra Leone","SG":"Singapore","SX":"Sint Maarten (Dutch part)","SK":"Slovakia","SI":"Slovenia","SB":"Solomon Islands","SO":"Somalia","ZA":"South Africa","GS":"South Georgia and the South Sandwich Islands","SS":"South Sudan","ES":"Spain","LK":"Sri Lanka","SD":"Sudan","SR":"Suriname","SJ":"Svalbard and Jan Mayen","SE":"Sweden","CH":"Switzerland","SY":"Syrian Arab Republic","TW":"Taiwan, Province of China","TJ":"Tajikistan","TZ":"Tanzania, United Republic of","TH":"Thailand","TL":"Timor-Leste","TG":"Togo","TK":"Tokelau","TO":"Tonga","TT":"Trinidad and Tobago","TN":"Tunisia","TR":"Turkey","TM":"Turkmenistan","TC":"Turks and Caicos Islands","TV":"Tuvalu","UG":"Uganda","UA":"Ukraine","AE":"United Arab Emirates","GB":"United Kingdom","US":"United States","UM":"United States Minor Outlying Islands","UY":"Uruguay","UZ":"Uzbekistan","VU":"Vanuatu","VE":"Venezuela, Bolivarian Republic of","VN":"Viet Nam","VG":"Virgin Islands, British","VI":"Virgin Islands, U.S.","WF":"Wallis and Futuna","EH":"Western Sahara","YE":"Yemen","ZM":"Zambia","ZW":"Zimbabwe"}

def emptyErrorFile():
    # MAKE FILE EMPTY ON START
    current_dir = os.path.dirname(os.path.abspath(__file__))
    file_path = os.path.join(current_dir, "ya_data_error.txt")
    with open(file_path, 'w') as file:
        file.truncate(0)

def is_convertible_to_int(s):
    try:
        int(s)
        return True
    except ValueError:
        return False

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
        file_path = os.path.join(current_dir, "ya_data_error.txt")
        with open(file_path, "a") as f:
            f.write(err)
            f.write("\n")
            f.close()
        yield {}

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

    if cron_type == "full":
        FILE_KEY = "temp/coffee/discovery/sync_ya.json"
        filename = "sync_ya.json"
    else:
        FILE_KEY = "temp/coffee/discovery/partial_sync_ya.json"
        filename = "partial_sync_ya.json"

    current_dir = os.path.dirname(os.path.abspath(__file__))
    file_path = os.path.join(current_dir, filename)
    s3.download_file(BUCKET_NAME, FILE_KEY, file_path)
    print("Downloading Complete")

def uploadJsonTos3():
    print("File start dumped To s3")
    if cron_type == "full":
        clickhouse_query = 'INSERT INTO FUNCTION s3(\'https://gcc-social-data.s3.ap-south-1.amazonaws.com/temp/coffee/discovery/sync_ya.json\', \'AKIAXGXUCIERSOUELY73\', \'dHb+F8Nj9o0imlDFq0gq/pCnWzSL/F5kkrywtSxO\', \'JSONEachRow\') SELECT * from dbt.mart_youtube_account SETTINGS s3_truncate_on_insert=1, output_format_json_quote_64bit_integers=0'
    else:
        if is_convertible_to_int(cron_type):
            timeInMinutes = int(cron_type)
            updated_at_value = datetime.utcnow() - timedelta(minutes=timeInMinutes)
        else:
            updated_at_value = datetime.utcnow() - timedelta(minutes=120)
        clickhouse_query = f'INSERT INTO FUNCTION s3(\'https://gcc-social-data.s3.ap-south-1.amazonaws.com/temp/coffee/discovery/partial_sync_ya.json\', \'AKIAXGXUCIERSOUELY73\', \'dHb+F8Nj9o0imlDFq0gq/pCnWzSL/F5kkrywtSxO\', \'JSONEachRow\') SELECT * from dbt.mart_youtube_account where updated_at > \'{updated_at_value}\' SETTINGS s3_truncate_on_insert=1, output_format_json_quote_64bit_integers=0'
    clickhouse_client.execute(clickhouse_query)

def getLocationByRow(row):
    location = []
    location_list = []
    if row['city'] and row['city'].strip() != '':
        splitted_city = row['city'].split()
        capitalized_splitted_city = [city.capitalize() for city in splitted_city]
        row['city'] = " ".join(capitalized_splitted_city)
        location.append("city_"+row['city'])
        location_list.append(row['city'])
    if row['state']:
        splitted_state = row['state'].split()
        capitalized_splitted_state = [state.capitalize() for state in splitted_state]
        row['state'] = " ".join(capitalized_splitted_state)
        if row['country'] == 'IN' and row['state'] in indian_states:
            location.append("state_"+row['state'])
            location_list.append(row['state'])
        elif row['country'] == 'IN' and row['state'] not in indian_states:
            row['state'] = None
    if row['country']:
        location.append("country_"+row['country'])
        if row['country'] in country_data:
            country_name = country_data[row['country']]
            location_list.append(country_name)
    
    return location,location_list

def youtube_prediction(subscribers,avg_views):
   
    video_rate = (0.42044539*subscribers)+(0.09215903*avg_views)
    if 100000<subscribers <=300000:
        video_rate = (0.42044539*100000)+(0.09215903*(avg_views/subscribers)*100000)
        video_rate = video_rate*(1+((subscribers/100000)-1)*0.1)
    if 300000<subscribers<=500000:
        video_rate = (0.42044539*100000)+(0.09215903*(avg_views/subscribers)*100000)
        video_rate = video_rate*(1+((subscribers/100000)-1)*0.6)
    if 500000<subscribers<=2000000:
        video_rate = (0.42044539*100000)+(0.09215903*(avg_views/subscribers)*100000)
        video_rate = video_rate*(1+((subscribers/100000)-1)*0.6)
    if subscribers>2000000:
        video_rate = (0.42044539*500000)+(0.09215903*(avg_views/subscribers)*500000)
        video_rate = video_rate*(1+((subscribers/100000)-5)*0.1)
    if (avg_views/subscribers) < 0.11:
        video_rate = video_rate*0.5
    est_price = json.dumps({"video_posts": video_rate})
    return est_price

def makeYARow(row):
    if row['languages'] is None or (len(row['languages']) == 1 and row['languages'][0] is None) or all(elem == "" for elem in row['languages']):
        row['languages'] = []

    if row['categories'] is None or (len(row['categories']) == 1 and row['categories'][0] is None) or all(elem == "" for elem in row['categories']):
        row['categories'] = []

    if len(row['categories'])>0:
        row['category_rank']=json.dumps(json.loads('{"' + str(row['categories'][0]) + '": "' + str(row['category_rank']) + '"}'))
    else:
        row['category_rank'] = None

    if row['linked_instagram_handle'] != "":
        row['linked_socials'] = json.dumps({"instagram_handle": row['linked_instagram_handle']})
    else:
        row['linked_socials'] = None
    
    
    for key in row:
        if key not in ["is_blacklisted"] and (row[key] is None or (isinstance(row[key], str) and row[key]=='') or (isinstance(row[key], (int, float)) and (math.isinf(row[key]) or math.isnan(row[key])))):
            row[key] = None
        # Story Reach Beautification
    est_reach_obj = {}
    if row['shorts_reach']:
        est_reach_obj['short_posts'] = row['shorts_reach']
    else:
        est_reach_obj['short_posts'] = 0

    if row['video_reach']:
        est_reach_obj['video_posts'] = row['video_reach']
    else:
        est_reach_obj['video_posts'] = 0
    
    row['est_reach'] = json.dumps(est_reach_obj)
    
    # Est Impressions Beautification

    est_impressions_obj = {}
    if row['shorts_impressions']:
        est_impressions_obj['short_posts'] = row['shorts_impressions']
    else:
        est_impressions_obj['short_posts'] = 0

    if row['video_impressions']:
        est_impressions_obj['video_posts'] = row['video_impressions']
    else:
        est_impressions_obj['video_posts'] = 0
    row['est_impressions'] = json.dumps(est_impressions_obj)
        
    # Similar Profile Data Beautification
    similar_group_obj = {}
    if row['group_avg_shorts_reach']:
        similar_group_obj['groupAvgShortsReach'] = row['group_avg_shorts_reach']
    else:
        similar_group_obj['groupAvgShortsReach'] = 0

    if row['group_avg_video_reach']:
        similar_group_obj['groupAvgVideoReach'] = row['group_avg_video_reach']
    else:
        similar_group_obj['groupAvgVideoReach'] = 0
    row['similar_profile_group_data'] = json.dumps(similar_group_obj)
    row['flag_contact_info_available'] = False
    if row['email'] or row['phone']:
        row['flag_contact_info_available'] = True
    
    if row['followers'] and row['avg_views']:
        row['est_post_price'] = youtube_prediction(row['followers'],row['avg_views'])
    else:
        row['est_post_price'] = None
    if row['updated_at'] is None:
        row['updated_at'] = datetime.fromtimestamp(time.time())
    
    if row['country'] == 'IN' and row['followers'] and row['views_count'] and row['avg_views']:
        row['enabled_for_saas'] = True
    else:
        row['enabled_for_saas'] = False
    
    is_blacklisted = row['is_blacklisted']
    if is_blacklisted == 1:
        row['is_blacklisted'] = True
    elif is_blacklisted == 0:
        row['is_blacklisted'] = False

    row['last_synced_at'] = datetime.fromtimestamp(time.time())

    if row['comments_rate'] is not None:
        row['comments_rate_grade'] = None
    
    if row['followers'] is not None:
        row['followers_grade'] = None

    if row['followers_growth7d'] is not None:
        row['followers_growth7d_grade'] = None

    if row['avg_posts_per_week'] is not None:
        row['avg_posts_per_week_grade'] = None
    
    if row['followers_growth1y'] is not None:
        row['followers_growth1y_grade'] = None
    
    if row['views_30d'] is not None:
        row['views_30d_grade'] = None
    
    if row['followers_growth30d'] is not None:
        row['followers_growth30d_grade'] = None
    return row

def makeCpRow(row):
    cpRow={}
    cpRow['platform']='YOUTUBE'
    cpRow['updated_by']='SYSTEM'
    cpRow['has_email']=False
    cpRow['has_phone']=False
    if row['email']:
        cpRow['has_email']=True
    if row['phone']:
        cpRow['has_phone']=True
    adminDetails = {}
    adminDetails['isBlacklisted']=False
    adminDetails['whatsappOptIn']=False
    userDetails = {}
    userDetails['ekycPending']=False
    userDetails['whatsappOptIn']=False
    userDetails['amazonStoreLinkVerified']=False
    userDetails['instantGratificationInvited']=False
    cpRow['admin_details']=json.dumps(adminDetails)
    cpRow['user_details']=json.dumps(userDetails)
    cpRow['gender'] = None
    cpRow['dob'] = None
    cpRow['phone'] = []
    cpRow['location'] = []
    cpRow['languages'] = []
    return cpRow

def updateCpColumns(existing_cp_row,row):
    adminDetails = {}
    userDetails = {}
    existing_cp_row['phone'] = []
    existing_cp_row['location'] = []
    existing_cp_row['languages'] = []

    if 'admin_details' in existing_cp_row and existing_cp_row['admin_details'] is not None:
        adminDetails = existing_cp_row['admin_details']

    if 'user_details' in existing_cp_row and existing_cp_row['user_details'] is not None:
        userDetails = existing_cp_row['user_details']

    if 'gender' in userDetails and userDetails['gender'] is not None:
        existing_cp_row['gender'] = userDetails['gender']
    elif 'gender' in adminDetails and adminDetails['gender'] is not None:
        existing_cp_row['gender'] = adminDetails['gender']
    elif 'gender' in row and row['gender'] is not None:
        existing_cp_row['gender'] = row['gender']
    
    if 'phone' in userDetails and userDetails['phone'] is not None:
        existing_cp_row['phone'].append(userDetails['phone'])
    if 'phone' in adminDetails and adminDetails['phone'] is not None and adminDetails['phone'] not in existing_cp_row['phone']:
        existing_cp_row['phone'].append(adminDetails['phone'])
    
    if 'dob' in userDetails and userDetails['dob'] is not None:
        existing_cp_row['dob'] = userDetails['dob']
    elif 'dob' in adminDetails and adminDetails['dob'] is not None:
        existing_cp_row['dob'] = adminDetails['dob']
    
    if 'location' in userDetails and userDetails['location'] is not None:
        merged_location_set = set(existing_cp_row['location']) | set(userDetails['location'])
        existing_cp_row['location'] = list(merged_location_set)

    if 'location' in adminDetails and adminDetails['location'] is not None:
        merged_location_set = set(existing_cp_row['location']) | set(adminDetails['location'])
        existing_cp_row['location'] = list(merged_location_set)

    if 'location_list' in row and row['location_list'] is not None:
        merged_location_set = set(existing_cp_row['location']) | set(row['location_list'])
        existing_cp_row['location'] = list(merged_location_set)
    
    if 'languages' in userDetails and userDetails['languages'] is not None:
        merged_langauges_set = set(existing_cp_row['languages']) | set(userDetails['languages'])
        existing_cp_row['languages'] = list(merged_langauges_set)

    if 'languages' in adminDetails and adminDetails['languages'] is not None:
        merged_langauges_set = set(existing_cp_row['languages']) | set(adminDetails['languages'])
        existing_cp_row['languages'] = list(merged_langauges_set)

    if 'languages' in row and row['languages'] is not None:
        merged_langauges_set = set(existing_cp_row['languages']) | set(row['languages'])
        existing_cp_row['languages'] = list(merged_langauges_set)
    for key, value in existing_cp_row.items():
        if value == '' or value == []:
            existing_cp_row[key] = None
    return existing_cp_row

def makeSearchPhraseAdminAndKeywordsAdmin(row,cpRow):
    adminDetails = {}
    userDetails = {}

    if 'admin_details' in cpRow and cpRow['admin_details'] is not None:
        adminDetails = cpRow['admin_details']

    if 'user_details' in cpRow and cpRow['user_details'] is not None:
        userDetails = cpRow['user_details']

    row['search_phrase_admin'] = None
    row['keywords_admin'] = []

    if 'channel_id' in row and row['channel_id'] is not None:
        row['keywords_admin'].append(row['channel_id'])
        if row['search_phrase_admin'] is not None:
            row['search_phrase_admin'] += " " + row['channel_id']
        else:
            row['search_phrase_admin'] = row['channel_id']

    if 'username' in row and row['username'] is not None:
        row['keywords_admin'].append(row['username'])
        if row['search_phrase_admin'] is not None:
            row['search_phrase_admin'] += " " + row['username']
        else:
            row['search_phrase_admin'] = row['username']
    
    if 'title' in row and row['title'] is not None:
        row['keywords_admin'].extend(word.lower() for word in row['title'].split())
        if row['search_phrase_admin'] is not None:
            row['search_phrase_admin'] += " " + row['title'].lower()
        else:
            row['search_phrase_admin'] = row['title'].lower()

    if 'name' in adminDetails and adminDetails['name'] is not None:
        row['keywords_admin'].extend(word.lower() for word in adminDetails['name'].split())
        if row['search_phrase_admin'] is not None:
            row['search_phrase_admin'] += " " + adminDetails['name'].lower()
        else:
            row['search_phrase_admin'] = adminDetails['name'].lower()

    if 'email' in adminDetails and adminDetails['email'] is not None:
        row['keywords_admin'].append(adminDetails['email'].lower())
        if row['search_phrase_admin'] is not None:
            row['search_phrase_admin'] += " " + adminDetails['email'].lower()
        else:
            row['search_phrase_admin'] = adminDetails['email'].lower()

    if 'phone' in adminDetails and adminDetails['phone'] is not None:
        row['keywords_admin'].append(adminDetails['phone'].lower())
        if row['search_phrase_admin'] is not None:
            row['search_phrase_admin'] += " " + adminDetails['phone'].lower()
        else:
            row['search_phrase_admin'] = adminDetails['phone'].lower()

    if 'secondaryPhone' in adminDetails and adminDetails['secondaryPhone'] is not None:
        row['keywords_admin'].append(adminDetails['secondaryPhone'].lower())
        if row['search_phrase_admin'] is not None:
            row['search_phrase_admin'] += " " + adminDetails['secondaryPhone'].lower()
        else:
            row['search_phrase_admin'] = adminDetails['secondaryPhone'].lower()

    if 'name' in userDetails and userDetails['name'] is not None:
        row['keywords_admin'].extend(word.lower() for word in userDetails['name'].split())
        if row['search_phrase_admin'] is not None:
            row['search_phrase_admin'] += " " + userDetails['name'].lower()
        else:
            row['search_phrase_admin'] = userDetails['name'].lower()

    if 'email' in userDetails and userDetails['email'] is not None:
        row['keywords_admin'].append(userDetails['email'].lower())
        if row['search_phrase_admin'] is not None:
            row['search_phrase_admin'] += " " + userDetails['email'].lower()
        else:
            row['search_phrase_admin'] = userDetails['email'].lower()

    if 'phone' in userDetails and userDetails['phone'] is not None:
        row['keywords_admin'].append(userDetails['phone'].lower())
        if row['search_phrase_admin'] is not None:
            row['search_phrase_admin'] += " " + userDetails['phone'].lower()
        else:
            row['search_phrase_admin'] = userDetails['phone'].lower()

    if 'secondaryPhone' in userDetails and userDetails['secondaryPhone'] is not None:
        row['keywords_admin'].append(userDetails['secondaryPhone'].lower())
        if row['search_phrase_admin'] is not None:
            row['search_phrase_admin'] += " " + userDetails['secondaryPhone'].lower()
        else:
            row['search_phrase_admin'] = userDetails['secondaryPhone'].lower()
    return row
    
def InsertUpdate(row,errorcount):
    try:
        row['location'],row['location_list'] = getLocationByRow(row)
        row = makeYARow(row)
        check_query = 'SELECT * FROM youtube_account WHERE channel_id = %s'
        postgres_cursor.execute(check_query, (row['channel_id'],))
        existing_row = postgres_cursor.fetchone()

        if existing_row:            
            check_cp_query = 'SELECT * FROM campaign_profiles WHERE campaign_profiles.platform = \'YOUTUBE\' AND platform_account_id = %s'
            postgres_cursor.execute(check_cp_query, (existing_row['id'],))
            existing_cp_row = postgres_cursor.fetchone()
            
            if existing_cp_row:
                existing_cp_row = updateCpColumns(existing_cp_row,row)
            
            postgres_cursor.execute(postgres_cp_update_query, (existing_cp_row['has_email'], existing_cp_row['has_phone'], existing_cp_row['gender'], existing_cp_row['phone'], existing_cp_row['dob'], existing_cp_row['location'], existing_cp_row['languages'], existing_cp_row['platform_account_id'], existing_cp_row['platform']))

            row = makeSearchPhraseAdminAndKeywordsAdmin(row,existing_cp_row)
            
            if row['linked_instagram_handle'] is None or row['linked_instagram_handle'] == '':
                row['profile_id'] = 'YA_' + str(existing_row['id'])
            else:
                row['profile_id'] = existing_row['profile_id']
           
            postgres_cursor.execute(postgres_ya_update_query, (row['title'], row['username'], row['gender'], row['description'], row['dob'], row['categories'], row['languages'], row['is_blacklisted'], row['city'], row['state'], row['country'], row['keywords'], row['uploads_count'], row['followers'], row['views_count'], row['avg_views'], row['video_reach'], row['shorts_reach'], row['thumbnail'], row['linked_socials'], row['linked_instagram_handle'], row['est_post_price'], row['est_reach'], row['est_impressions'], row['followers_growth7d'], row['country_rank'], row['category_rank'], row['authentic_engagement'], row['comment_rate_percentage'], row['video_views_last30'], row['search_phrase'], row['location'], row['label'], row['phone'], row['email'], row['reaction_rate'], row['comments_rate_grade'], row['followers_grade'], row['comments_rate'], row['followers_growth7d_grade'], row['group_key'], row['cpm'], row['avg_posts_per_week'], row['avg_posts_per_week_grade'], row['followers_growth1y'], row['followers_growth1y_grade'], row['avg_shorts_views_30d'], row['avg_video_views_30d'], row['latest_video_publish_time'], row['similar_profile_group_data'],row['views_30d_grade'], row['updated_at'], row['followers_growth30d'], row['followers_growth30d_grade'], row['profile_type'], row['enabled_for_saas'], row['views_30d'], row['uploads_30d'], row['location_list'], row['profile_id'], row['last_synced_at'], row['search_phrase_admin'], row['keywords_admin'], row['channel_id']))

        else:
            cpRow = makeCpRow(row)
            row = makeSearchPhraseAdminAndKeywordsAdmin(row,cpRow)
            postgres_cursor.execute(postgres_ya_insert_query, row)
            inserted_id = postgres_cursor.fetchone()[0]
            cpRow['platform_account_id']=inserted_id
            postgres_cursor.execute(postgres_cp_insert_query, cpRow)
    except Exception as err:
        # traceback.print_exc()
        errorcount = errorcount+1
        current_dir = os.path.dirname(os.path.abspath(__file__))
        file_path = os.path.join(current_dir, "ya_data_error.txt")

        with open(file_path, "a") as f:
            f.write(row['channel_id'])
            f.write("\n")
            f.close()
              
    postgres_conn.commit()
    return errorcount

def readJsonFile(count,errorcount,start_time):
    if cron_type == "full":
        filename = "sync_ya.json"
    else:
        filename = "partial_sync_ya.json"
        
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
            print(str(count) + " Youtube Accounts Completed with "+ str(errorcount) + " errors in " + str(exec_time)+ " seconds")
            start_time = time.time()


def afterScriptOperations():
    # Sync Profile Id
        try:
            print("Syncing Profile Id")
            postgres_cursor.execute("BEGIN;")
            query = 'UPDATE youtube_account SET profile_id = Concat(\'YA_\',id) where profile_id IS NULL'
            postgres_cursor.execute(query)

            query = 'UPDATE youtube_account SET profile_id = instagram_account.profile_id FROM instagram_account WHERE  youtube_account.linked_instagram_handle != \'\' and youtube_account.linked_instagram_handle = instagram_account.handle'
            postgres_cursor.execute(query)

            query = 'UPDATE youtube_account SET profile_id = instagram_account.profile_id FROM instagram_account WHERE instagram_account.linked_channel_id != \'\' and youtube_account.channel_id = instagram_account.linked_channel_id'
            postgres_cursor.execute(query)
            postgres_cursor.execute("COMMIT;")

        except psycopg2.Error as e:
            # Rollback the transaction if any error occurs
            postgres_cursor.execute("ROLLBACK;")
            print("Error:", e)
        finally:
            # Close the cursor and connection
            postgres_cursor.close()
            postgres_conn.close()
        print("youtube_account TABLE column profile_id GOT Linked")
        print("youtube_account completed. Check ya_data_error.txt for the channel_id's")
        clickhouse_client.disconnect()

def main():
    global clickhouse_client, postgres_cursor, postgres_conn, postgres_ya_insert_query, postgres_cp_insert_query, postgres_ya_update_query, postgres_cp_update_query, errorcount, count, chunk_size, cron_type

    if len(sys.argv) > 1:
        cron_type = sys.argv[1] # partial or full
    else:
        cron_type = '' 

    errorcount = 0
    count = 0
    chunk_size = 1000
    # Define connection parameters

    # ch_host = '52.66.200.31'
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
    postgres_cursor = postgres_conn.cursor(cursor_factory=extras.DictCursor)


    postgres_ya_insert_query = 'INSERT INTO youtube_account (title, username, gender, description, dob, categories, languages, is_blacklisted, city, state, country, keywords, uploads_count, followers, views_count, avg_views, video_reach, shorts_reach, thumbnail, linked_socials, linked_instagram_handle, est_post_price, est_reach, est_impressions, followers_growth7d, country_rank, category_rank, authentic_engagement, comment_rate_percentage, video_views_last30, search_phrase, location, label, phone, email, reaction_rate, comments_rate_grade, followers_grade, comments_rate, followers_growth7d_grade, group_key, cpm, avg_posts_per_week, avg_posts_per_week_grade, followers_growth1y, followers_growth1y_grade, avg_shorts_views_30d, avg_video_views_30d, latest_video_publish_time, similar_profile_group_data,views_30d_grade, updated_at, followers_growth30d, followers_growth30d_grade, profile_type, enabled_for_saas, views_30d, uploads_30d, location_list, keywords_admin, search_phrase_admin, last_synced_at, channel_id) VALUES (%(title)s, %(username)s, %(gender)s, %(description)s, %(dob)s, %(categories)s, %(languages)s, %(is_blacklisted)s, %(city)s, %(state)s, %(country)s, %(keywords)s, %(uploads_count)s, %(followers)s, %(views_count)s, %(avg_views)s, %(video_reach)s, %(shorts_reach)s, %(thumbnail)s, %(linked_socials)s, %(linked_instagram_handle)s, %(est_post_price)s, %(est_reach)s, %(est_impressions)s, %(followers_growth7d)s, %(country_rank)s, %(category_rank)s, %(authentic_engagement)s, %(comment_rate_percentage)s, %(video_views_last30)s, %(search_phrase)s, %(location)s, %(label)s, %(phone)s, %(email)s, %(reaction_rate)s, %(comments_rate_grade)s, %(followers_grade)s, %(comments_rate)s, %(followers_growth7d_grade)s, %(group_key)s, %(cpm)s, %(avg_posts_per_week)s, %(avg_posts_per_week_grade)s, %(followers_growth1y)s, %(followers_growth1y_grade)s, %(avg_shorts_views_30d)s, %(avg_video_views_30d)s, %(latest_video_publish_time)s, %(similar_profile_group_data)s, %(views_30d_grade)s, %(updated_at)s, %(followers_growth30d)s, %(followers_growth30d_grade)s, %(profile_type)s, %(enabled_for_saas)s, %(views_30d)s, %(uploads_30d)s, %(location_list)s, %(keywords_admin)s, %(search_phrase_admin)s, %(last_synced_at)s, %(channel_id)s) RETURNING id'

    postgres_cp_insert_query = 'INSERT INTO campaign_profiles (platform, platform_account_id, admin_details, user_details, on_gcc, on_gcc_app, gcc_user_account_id, updated_by, enabled, has_email, has_phone, gender, phone, dob, location, languages) VALUES (%(platform)s, %(platform_account_id)s, %(admin_details)s, %(user_details)s, false, false, NULL, %(updated_by)s, true, %(has_email)s, %(has_phone)s, %(gender)s, %(phone)s, %(dob)s, %(location)s, %(languages)s)'

    postgres_ya_update_query = 'UPDATE youtube_account SET title = %s, username = %s, gender = %s, description = %s, dob = %s, categories = %s, languages = %s, is_blacklisted = %s, city = %s, state = %s, country = %s, keywords = %s, uploads_count = %s, followers = %s, views_count = %s, avg_views = %s, video_reach = %s, shorts_reach = %s, thumbnail = %s, linked_socials = %s, linked_instagram_handle = %s, est_post_price = %s, est_reach = %s, est_impressions = %s, followers_growth7d = %s, country_rank = %s, category_rank = %s, authentic_engagement = %s, comment_rate_percentage = %s, video_views_last30 = %s, search_phrase = %s, location = %s, label = %s, phone = %s, email = %s, reaction_rate = %s, comments_rate_grade = %s, followers_grade = %s, comments_rate = %s, followers_growth7d_grade = %s, group_key = %s, cpm = %s, avg_posts_per_week = %s, avg_posts_per_week_grade = %s, followers_growth1y = %s, followers_growth1y_grade = %s, avg_shorts_views_30d = %s, avg_video_views_30d = %s, latest_video_publish_time = %s, similar_profile_group_data = %s, views_30d_grade = %s, updated_at = %s, followers_growth30d = %s, followers_growth30d_grade = %s, profile_type = %s, enabled_for_saas = %s, views_30d = %s, uploads_30d = %s,location_list = %s, profile_id = %s, last_synced_at = %s, search_phrase_admin = %s, keywords_admin = %s WHERE channel_id = %s RETURNING id'

    postgres_cp_update_query = 'UPDATE campaign_profiles SET has_email = %s, has_phone = %s, gender = %s, phone = %s, dob = %s, location = %s, languages = %s where platform_account_id = %s and platform = %s'

    try:
        emptyErrorFile()
        uploadJsonTos3()
        downloadCSV()
        start_time = time.time()
        readJsonFile(count,errorcount,start_time)
        afterScriptOperations()

    except Exception as err:
        traceback.print_exc()

if __name__ == "__main__":
    main()
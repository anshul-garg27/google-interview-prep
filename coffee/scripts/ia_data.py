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

indian_states = ["Andhra Pradesh","Arunachal Pradesh","Assam","Bihar","Chhattisgarh","Goa","Gujarat","Haryana","Himachal Pradesh","Jharkhand","Karnataka","Kerala","Madhya Pradesh","Maharashtra","Manipur","Meghalaya","Mizoram","Nagaland","Odisha","Punjab","Rajasthan","Sikkim","Tamil Nadu","Telangana","Tripura","Uttar Pradesh","Uttarakhand","West Bengal","Jammu and Kashmir","Delhi","Puducherry","chandigarh","Ladakh","Lakshadweep","Andaman and Nicobar Islands"]

country_data = {"AF":"Afghanistan","AX":"Åland Islands","AL":"Albania","DZ":"Algeria","AS":"American Samoa","AD":"Andorra","AO":"Angola","AI":"Anguilla","AQ":"Antarctica","AG":"Antigua and Barbuda","AR":"Argentina","AM":"Armenia","AW":"Aruba","AU":"Australia","AT":"Austria","AZ":"Azerbaijan","BS":"Bahamas","BH":"Bahrain","BD":"Bangladesh","BB":"Barbados","BY":"Belarus","BE":"Belgium","BZ":"Belize","BJ":"Benin","BM":"Bermuda","BT":"Bhutan","BO":"Bolivia","BQ":"Bonaire, Sint Eustatius and Saba","BA":"Bosnia and Herzegovina","BW":"Botswana","BV":"Bouvet Island","BR":"Brazil","IO":"British Indian Ocean Territory","BN":"Brunei Darussalam","BG":"Bulgaria","BF":"Burkina Faso","BI":"Burundi","CV":"Cabo Verde","KH":"Cambodia","CM":"Cameroon","CA":"Canada","KY":"Cayman Islands","CF":"Central African Republic","TD":"Chad","CL":"Chile","CN":"China","CX":"Christmas Island","CC":"Cocos (Keeling) Islands","CO":"Colombia","KM":"Comoros","CG":"Congo","CD":"Congo, Democratic Republic of the","CK":"Cook Islands","CR":"Costa Rica","CI":"Côte d'Ivoire","HR":"Croatia","CU":"Cuba","CW":"Curaçao","CY":"Cyprus","CZ":"Czech Republic","DK":"Denmark","DJ":"Djibouti","DM":"Dominica","DO":"Dominican Republic","EC":"Ecuador","EG":"Egypt","SV":"El Salvador","GQ":"Equatorial Guinea","ER":"Eritrea","EE":"Estonia","SZ":"Eswatini","ET":"Ethiopia","FK":"Falkland Islands (Malvinas)","FO":"Faroe Islands","FJ":"Fiji","FI":"Finland","FR":"France","GF":"French Guiana","PF":"French Polynesia","TF":"French Southern Territories","GA":"Gabon","GM":"Gambia","GE":"Georgia","DE":"Germany","GH":"Ghana","GI":"Gibraltar","GR":"Greece","GL":"Greenland","GD":"Grenada","GP":"Guadeloupe","GU":"Guam","GT":"Guatemala","GG":"Guernsey","GN":"Guinea","GW":"Guinea-Bissau","GY":"Guyana","HT":"Haiti","HM":"Heard Island and McDonald Islands","VA":"Holy See","HN":"Honduras","HK":"Hong Kong","HU":"Hungary","IS":"Iceland","IN":"India","ID":"Indonesia","IR":"Iran, Islamic Republic of","IQ":"Iraq","IE":"Ireland","IM":"Isle of Man","IL":"Israel","IT":"Italy","JM":"Jamaica","JP":"Japan","JE":"Jersey","JO":"Jordan","KZ":"Kazakhstan","KE":"Kenya","KI":"Kiribati","KP":"Korea, Democratic People's Republic of","KR":"Korea, Republic of","KW":"Kuwait","KG":"Kyrgyzstan","LA":"Lao People's Democratic Republic","LV":"Latvia","LB":"Lebanon","LS":"Lesotho","LR":"Liberia","LY":"Libya","LI":"Liechtenstein","LT":"Lithuania","LU":"Luxembourg","MO":"Macao","MG":"Madagascar","MW":"Malawi","MY":"Malaysia","MV":"Maldives","ML":"Mali","MT":"Malta","MH":"Marshall Islands","MQ":"Martinique","MR":"Mauritania","MU":"Mauritius","YT":"Mayotte","MX":"Mexico","FM":"Micronesia, Federated States of","MD":"Moldova, Republic of","MC":"Monaco","MN":"Mongolia","ME":"Montenegro","MS":"Montserrat","MA":"Morocco","MZ":"Mozambique","MM":"Myanmar","NA":"Namibia","NR":"Nauru","NP":"Nepal","NL":"Netherlands","NC":"New Caledonia","NZ":"New Zealand","NI":"Nicaragua","NE":"Niger","NG":"Nigeria","NU":"Niue","NF":"Norfolk Island","MK":"North Macedonia","MP":"Northern Mariana Islands","NO":"Norway","OM":"Oman","PK":"Pakistan","PW":"Palau","PS":"Palestine, State of","PA":"Panama","PG":"Papua New Guinea","PY":"Paraguay","PE":"Peru","PH":"Philippines","PN":"Pitcairn","PL":"Poland","PT":"Portugal","PR":"Puerto Rico","QA":"Qatar","RE":"Réunion","RO":"Romania","RU":"Russian Federation","RW":"Rwanda","BL":"Saint Barthélemy","SH":"Saint Helena, Ascension and Tristan da Cunha","KN":"Saint Kitts and Nevis","LC":"Saint Lucia","MF":"Saint Martin (French part)","PM":"Saint Pierre and Miquelon","VC":"Saint Vincent and the Grenadines","WS":"Samoa","SM":"San Marino","ST":"Sao Tome and Principe","SA":"Saudi Arabia","SN":"Senegal","RS":"Serbia","SC":"Seychelles","SL":"Sierra Leone","SG":"Singapore","SX":"Sint Maarten (Dutch part)","SK":"Slovakia","SI":"Slovenia","SB":"Solomon Islands","SO":"Somalia","ZA":"South Africa","GS":"South Georgia and the South Sandwich Islands","SS":"South Sudan","ES":"Spain","LK":"Sri Lanka","SD":"Sudan","SR":"Suriname","SJ":"Svalbard and Jan Mayen","SE":"Sweden","CH":"Switzerland","SY":"Syrian Arab Republic","TW":"Taiwan, Province of China","TJ":"Tajikistan","TZ":"Tanzania, United Republic of","TH":"Thailand","TL":"Timor-Leste","TG":"Togo","TK":"Tokelau","TO":"Tonga","TT":"Trinidad and Tobago","TN":"Tunisia","TR":"Turkey","TM":"Turkmenistan","TC":"Turks and Caicos Islands","TV":"Tuvalu","UG":"Uganda","UA":"Ukraine","AE":"United Arab Emirates","GB":"United Kingdom","US":"United States","UM":"United States Minor Outlying Islands","UY":"Uruguay","UZ":"Uzbekistan","VU":"Vanuatu","VE":"Venezuela, Bolivarian Republic of","VN":"Viet Nam","VG":"Virgin Islands, British","VI":"Virgin Islands, U.S.","WF":"Wallis and Futuna","EH":"Western Sahara","YE":"Yemen","ZM":"Zambia","ZW":"Zimbabwe"}


def emptyErrorFile():
    # MAKE FILE EMPTY ON START
    
    current_dir = os.path.dirname(os.path.abspath(__file__))
    file_path = os.path.join(current_dir, "ia_data_error.txt")
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
        file_path = os.path.join(current_dir, "ia_data_error.txt")
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
        FILE_KEY = "temp/coffee/discovery/sync_ia.json"
        filename = "sync_ia.json"
    else:
        FILE_KEY = "temp/coffee/discovery/partial_sync_ia.json"
        filename = "partial_sync_ia.json"

    current_dir = os.path.dirname(os.path.abspath(__file__))
    file_path = os.path.join(current_dir, filename)
    s3.download_file(BUCKET_NAME, FILE_KEY, file_path)
    print("Downloading Complete")

def uploadJsonTos3():
    print("File start dumped To s3")
    if cron_type == "full":
        clickhouse_query = 'INSERT INTO FUNCTION s3(\'https://gcc-social-data.s3.ap-south-1.amazonaws.com/temp/coffee/discovery/sync_ia.json\', \'AKIAXGXUCIERSOUELY73\', \'dHb+F8Nj9o0imlDFq0gq/pCnWzSL/F5kkrywtSxO\', \'JSONEachRow\') SELECT * from dbt.mart_instagram_account where handle != \'MISSING\' SETTINGS s3_truncate_on_insert=1, output_format_json_quote_64bit_integers=0'
    else:
        if is_convertible_to_int(cron_type):
            timeInMinutes = int(cron_type)
            updated_at_value = datetime.utcnow() - timedelta(minutes=timeInMinutes)
        else:
            updated_at_value = datetime.utcnow() - timedelta(minutes=120)
        
        clickhouse_query = f'INSERT INTO FUNCTION s3(\'https://gcc-social-data.s3.ap-south-1.amazonaws.com/temp/coffee/discovery/partial_sync_ia.json\', \'AKIAXGXUCIERSOUELY73\', \'dHb+F8Nj9o0imlDFq0gq/pCnWzSL/F5kkrywtSxO\', \'JSONEachRow\') SELECT * from dbt.mart_instagram_account where handle != \'MISSING\' and updated_at > \'{updated_at_value}\' SETTINGS s3_truncate_on_insert=1, output_format_json_quote_64bit_integers=0'
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
        if row['state']:
            location.append("state_"+row['state'])
            location_list.append(row['state'])
    if row['country']:
        location.append("country_"+row['country'])
        if row['country'] in country_data:
            country_name = country_data[row['country']]
            location_list.append(country_name)
    
    return location,location_list

def instagram_price_v2(followers,eng_rate):
    eng_num = followers*eng_rate/100
    # num_post = story+image+video+reel+carousel
    if followers <= 10000:
        image_rate = (0.08580218 * followers + 0.00947688 * eng_num) * 2.3
        story_rate = image_rate * 0.6
        video_rate = image_rate * 2
        reel_rate = image_rate * 2.8
    if 10000 < followers <= 20000:
        image_rate = (0.024 * followers + 0.00947688 * eng_num) * 7.5
        story_rate = image_rate * 0.6
        video_rate = image_rate * 2
        reel_rate = image_rate * 2.8
    if 20000 < followers <= 50000:
        image_rate_1 = (0.03517652 * followers + 0.01662693 * eng_num) * 3.5
        image_rate_2 = (0.024 * 20000 + 0.00947688 * eng_num) * 7.5
        if image_rate_2 > image_rate_1:
            image_rate = image_rate_2
        else: 
            image_rate = image_rate_1
        story_rate = image_rate * 0.6
        video_rate = image_rate * 2
        reel_rate = image_rate * 2.8
    if 50000 < followers <= 100000:
        image_rate_1 = (0.02752566 * followers + 0.00921316 * eng_num) * 2.9
        image_rate_2 = (0.03517652 * 50000 + 0.01662693 * eng_num) * 3.5
        if image_rate_2 > image_rate_1:
            image_rate = image_rate_2 
        else:
            image_rate = image_rate_1
        story_rate = image_rate * 0.6
        video_rate = image_rate * 2
        reel_rate = image_rate * 2.8
    if 100000 < followers <=300000:
        image_rate_1 = (0.01687308 * followers + 0.10067394 * eng_num) * 2.8
        image_rate_2 = (0.02752566 * 100000 + 0.00921316 * eng_num) * 2.9
        if image_rate_2 > image_rate_1:
            image_rate = image_rate_2 
        else:
            image_rate = image_rate_1
        story_rate = image_rate * 0.6
        video_rate = image_rate * 2
        reel_rate = image_rate * 2.8
    if 300000 <followers <= 800000 :
        image_rate_1 = (0.01595681 * followers + 0.08665652 * eng_num) * 3.0
        image_rate_2 = (0.01687308 * followers + 0.10067394 * eng_num) * 2.8
        if image_rate_2 > image_rate_1:
            image_rate = image_rate_2 
        else:
            image_rate = image_rate_1
        story_rate = image_rate * 0.6
        video_rate = image_rate * 2
        reel_rate = image_rate * 2.8
    if 800000 <followers <= 1000000  :
        image_rate = (0.01595681 * followers + 0.08665652 * eng_num) * 2.5
        story_rate = image_rate * 0.6
        video_rate = image_rate * 2
        reel_rate = image_rate * 2.8
    if 1000000 <followers <= 3000000  :
        image_rate = (0.01595681 * followers + 0.08665652 * eng_num) * 2.5
        story_rate = image_rate * 0.6
        video_rate = image_rate * 2
        reel_rate = image_rate * 2.8
    if 3000000 < followers <= 5000000 :
        image_rate = (0.01595681 * followers + 0.08665652 * eng_num) * 2.5
        story_rate = image_rate * 0.6
        video_rate = image_rate * 2
        reel_rate = image_rate * 2.8
    if followers>5000000:
        video_rate = 0.24583049*5000000 + 0.01875384*eng_num
        story_rate = video_rate*0.03
        image_rate = video_rate*0.40
        reel_rate = video_rate*0.8
        video_rate = video_rate*0.9
    est_price = json.dumps({"reel_posts": reel_rate, "image_posts": image_rate, "story_posts": story_rate})
    return est_price

def makeIARow(row):
    if row['categories'] is None or (len(row['categories']) == 1 and row['categories'][0] is None) or all(elem == "" for elem in row['categories']):
        row['categories'] = []
    if len(row['categories'])>0:
        row['category_rank']=json.dumps(json.loads('{"' + str(row['categories'][0]) + '": "' + str(row['category_rank']) + '"}'))
    else:
        row['category_rank'] = None
    if row['linked_channel_id'] != "":
        row['linked_socials'] = json.dumps({"channel_id": row['linked_channel_id']})
    else:
        row['linked_socials'] = None

    for key in row:
        if key not in ["is_private", "is_blacklisted", "is_verified"] and (not row[key] or (isinstance(row[key], (int, float)) and (math.isinf(row[key]) or math.isnan(row[key])))):
            row[key] = None
    
    # Story Reach Beautification
    est_reach_obj = {}
    if row['reels_reach']:
        est_reach_obj['reel_posts'] = row['reels_reach']
    else:
        est_reach_obj['reel_posts'] = 0

    if row['image_reach']:
        est_reach_obj['image_posts'] = row['image_reach']
    else:
        est_reach_obj['image_posts'] = 0
    if row['story_reach']:
        est_reach_obj['story_posts'] = row['story_reach']
    else:
        est_reach_obj['story_posts'] = 0
    row['est_reach'] = json.dumps(est_reach_obj)

    # Est Impressions Beautification

    est_impressions_obj = {}
    if row['reels_impressions']:
        est_impressions_obj['reel_posts'] = row['reels_impressions']
    else:
        est_impressions_obj['reel_posts'] = 0

    if row['image_impressions']:
        est_impressions_obj['image_posts'] = row['image_impressions']
    else:
        est_impressions_obj['image_posts'] = 0
    row['est_impressions'] = json.dumps(est_impressions_obj)

    # Similar Profile Data Beautification
    similar_group_obj = {}
    if row['group_avg_reels_reach']:
        similar_group_obj['groupAvgReelsReach'] = row['group_avg_reels_reach']
    else:
        similar_group_obj['groupAvgReelsReach'] = 0

    if row['group_avg_image_reach']:
        similar_group_obj['groupAvgImageReach'] = row['group_avg_image_reach']
    else:
        similar_group_obj['groupAvgImageReach'] = 0

    if row['group_avg_story_reach']:
        similar_group_obj['groupAvgStoryReach'] = row['group_avg_story_reach']
    else:
        similar_group_obj['groupAvgStoryReach'] = 0
    row['similar_profile_group_data'] = json.dumps(similar_group_obj)
    
    if row['followers'] and row['engagement_rate']:
        row['est_post_price'] = instagram_price_v2(row['followers'],row['engagement_rate'])
    else:
        row['est_post_price'] = None

    row['flag_contact_info_available'] = False
    if row['email'] or row['phone']:
            row['flag_contact_info_available'] = True

    if row['updated_at'] is None:
        row['updated_at'] = datetime.fromtimestamp(time.time())

    if row['country'] == 'IN' and row['followers'] and row['engagement_rate'] and row['avg_reach'] and row['avg_likes'] and row['avg_comments']:
        row['enabled_for_saas'] = True
    else:
        row['enabled_for_saas'] = False
    
    is_private = row['is_private']
    if is_private == 1:
        row['is_private'] = True
    elif is_private == 0:
        row['is_private'] = False

    is_blacklisted = row['is_blacklisted']
    if is_blacklisted == 1:
        row['is_blacklisted'] = True
    elif is_blacklisted == 0:
        row['is_blacklisted'] = False

    is_verified = row['is_verified']
    if is_verified == 1:
        row['is_verified'] = True
    elif is_verified == 0:
        row['is_verified'] = False
    
    row['last_synced_at'] = datetime.fromtimestamp(time.time())

    if row['avg_comments'] is None:
        row['avg_comments_grade'] = None

    if row['avg_likes'] is None:
        row['avg_likes_grade'] = None

    if row['comments_rate'] is None:
        row['comments_rate_grade'] = None
    
    if row['followers'] is None:
        row['followers_grade'] = None

    if row['engagement_rate'] is None:
        row['engagement_rate_grade'] = None
        row['er_grade'] = None

    if row['reels_reach'] is not None:
        row['reels_reach_grade'] = None
    
    if row['likes_to_comment_ratio'] is not None:
        row['likes_to_comment_ratio_grade'] = None
    
    if row['followers_growth7d'] is not None:
        row['followers_growth7d_grade'] = None
        
    if row['followers_growth30d'] is not None:
        row['followers_growth30d_grade'] = None
        
    if row['followers_growth90d'] is not None:
        row['followers_growth90d_grade'] = None

    if row['audience_reachability'] is not None:
        row['audience_reachability_grade'] = None
        
    if row['audience_authencity'] is not None:
        row['audience_authencity_grade'] = None

    if row['audience_quality'] is not None:
        row['audience_quality_grade'] = None
    
    if row['post_count'] is not None:
        row['post_count_grade'] = None

    if row['likes_spread'] is not None:
        row['likes_spread_grade'] = None

    if row['image_reach'] is not None:
        row['image_reach_grade'] = None
    
    if row['story_reach'] is not None:
        row['story_reach_grade'] = None
    
    if row['image_impressions'] is not None:
        row['image_impressions_grade'] = None
    
    if row['reels_impressions'] is not None:
        row['reels_impressions'] = None
    return row

def makeCpRow(row):
    cpRow={}
    cpRow['platform']='INSTAGRAM'
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
    cpRow['phone'] = None
    cpRow['location'] = None
    cpRow['languages'] = None
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

    if 'handle' in row and row['handle'] is not None:
        row['keywords_admin'].append(row['handle'].lower())
        if row['search_phrase_admin'] is not None:
            row['search_phrase_admin'] += " " + row['handle'].lower()
        else:
            row['search_phrase_admin'] = row['handle'].lower()
    
    if 'name' in row and row['name'] is not None:
        row['keywords_admin'].extend(word.lower() for word in row['name'].split())
        if row['search_phrase_admin'] is not None:
            row['search_phrase_admin'] += " " + row['name'].lower()
        else:
            row['search_phrase_admin'] = row['name'].lower()

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
        row = makeIARow(row)
        check_query = 'SELECT * FROM instagram_account WHERE ig_id = %s'
        postgres_cursor.execute(check_query, (row['ig_id'],))
        existing_row = postgres_cursor.fetchone()
        
        if existing_row:
            check_cp_query = 'SELECT * FROM campaign_profiles WHERE campaign_profiles.platform = \'INSTAGRAM\' AND platform_account_id = %s'
            postgres_cursor.execute(check_cp_query, (existing_row['id'],))
            existing_cp_row = postgres_cursor.fetchone()
            if existing_cp_row:
                existing_cp_row = updateCpColumns(existing_cp_row,row)
            
            postgres_cursor.execute(postgres_cp_update_query, (existing_cp_row['has_email'], existing_cp_row['has_phone'], existing_cp_row['gender'], existing_cp_row['phone'], existing_cp_row['dob'], existing_cp_row['location'], existing_cp_row['languages'], existing_cp_row['platform_account_id'], existing_cp_row['platform']))
            
            row = makeSearchPhraseAdminAndKeywordsAdmin(row,existing_cp_row)
            postgres_cursor.execute(postgres_ia_update_query, (row['name'], row['handle'], row['business_id'], row['thumbnail'], row['dob'], row['bio'], row['location'], row['categories'], row['label'], row['languages'], row['gender'], row['city'], row['state'], row['country'], row['keywords'], row['is_blacklisted'], row['followers'], row['following'], row['post_count'], row['ffratio'], row['avg_views'], row['avg_likes'], row['avg_reach'], row['avg_reels_play_count'], row['engagement_rate'], row['followers_growth7d'], row['reels_reach'], row['story_reach'], row['image_reach'], row['avg_comments'], row['er_grade'], row['linked_socials'], row['linked_channel_id'], row['est_post_price'], row['est_reach'], row['est_impressions'], row['flag_contact_info_available'], row['is_verified'], row['post_frequency_week'], row['country_rank'], row['category_rank'], row['authentic_engagement'], row['comment_rate_percentage'], row['likes_spread_percentage'], row['group_key'], row['search_phrase'], row['phone'], row['email'], row['avg_comments_grade'], row['avg_likes_grade'], row['comments_rate_grade'], row['followers_grade'], row['engagement_rate_grade'], row['reels_reach_grade'], row['likes_to_comment_ratio_grade'], row['followers_growth7d_grade'], row['followers_growth30d_grade'], row['followers_growth90d_grade'], row['audience_reachability_grade'], row['audience_authencity_grade'], row['audience_quality_grade'], row['post_count_grade'], row['likes_spread_grade'], row['followers_growth30d'], row['similar_profile_group_data'], row['image_reach_grade'], row['story_reach_grade'], row['avg_reels_play_30d'], row['updated_at'], row['profile_type'], row['enabled_for_saas'],row['plays_30d'], row['uploads_30d'], row['account_type'],row['location_list'],row['impressions_30d'], row['reach_30d'], row['image_impressions_grade'], row['reels_impressions_grade'], row['is_private'], row['last_synced_at'], row['search_phrase_admin'], row['keywords_admin'], row['ig_id']))
        else:
            cpRow = makeCpRow(row)
            row = makeSearchPhraseAdminAndKeywordsAdmin(row,cpRow)
            postgres_cursor.execute(postgres_ia_insert_query, row)
            inserted_id = postgres_cursor.fetchone()[0]
            cpRow['platform_account_id']=inserted_id
            postgres_cursor.execute(postgres_cp_insert_query, cpRow)
            
    except Exception as err:
        traceback.print_exc()
        current_dir = os.path.dirname(os.path.abspath(__file__))
        file_path = os.path.join(current_dir, "ia_data_error.txt")

        with open(file_path, "a") as f:
            f.write(row['ig_id'])
            f.write("\n")
            f.close()
        errorcount = errorcount+1

    postgres_conn.commit()
    return errorcount

def readJsonFile(count,errorcount,start_time):
    if cron_type == "full":
        filename = "sync_ia.json"
    else:
        filename = "partial_sync_ia.json"
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
            print(str(count) + " Instagram Accounts Completed with "+ str(errorcount) + " errors in " + str(exec_time)+ " seconds")
            start_time = time.time()

def createClickhouseConnection():
    # Define connection parameters

    # ch_host = '52.66.200.31'
    ch_host = '172.31.28.68'
    ch_port = 9000
    ch_database = 'dbt'
    ch_user = 'coffee'
    ch_password = '0M7sooN_WQW'
    clickhouse_client = Client(ch_host, port=ch_port, database=ch_database, user=ch_user, password=ch_password)
    return clickhouse_client

def createPostgresConnection():
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
    return postgres_conn,postgres_cursor

def afterScriptOperations():
    query = 'REFRESH MATERIALIZED VIEW CONCURRENTLY mv_location_master_tokens'
    postgres_cursor.execute(query)

    query = 'REFRESH MATERIALIZED VIEW CONCURRENTLY mv_location_master'
    postgres_cursor.execute(query) 
    print("MV views For Location Updated")

    # Sync Profile Id
    query = 'UPDATE instagram_account SET profile_id = Concat(\'IA_\',id) where profile_id IS NULL'
    postgres_cursor.execute(query)
    print("instagram_account TABLE column profile_id GOT UPDATED")

    postgres_conn.commit()
    print("instagram_account completed. Check ia_data_error.txt for the ig_id's")
    postgres_cursor.close()
    postgres_conn.close()
    clickhouse_client.disconnect()

def main():
    global clickhouse_client, postgres_cursor, postgres_conn, postgres_ia_insert_query, postgres_cp_insert_query, postgres_ia_update_query, postgres_cp_update_query, errorcount, count, chunk_size, cron_type

    if len(sys.argv) > 1:
        cron_type = sys.argv[1] # partial or full
    else:
        cron_type = ''

    errorcount = 0
    count = 0
    chunk_size = 1000
    clickhouse_client = createClickhouseConnection()
    postgres_conn, postgres_cursor = createPostgresConnection()

    postgres_ia_insert_query = 'INSERT INTO instagram_account (name, handle, business_id, thumbnail, dob, bio, location, categories, label, languages, gender, city, state, country, keywords, is_blacklisted, followers, following, post_count, ffratio, avg_views, avg_likes, avg_reach, avg_reels_play_count, engagement_rate, followers_growth7d, reels_reach, story_reach, image_reach, avg_comments, er_grade, linked_socials, linked_channel_id,est_post_price, est_reach, est_impressions, flag_contact_info_available, is_verified, post_frequency_week, country_rank, category_rank, authentic_engagement, comment_rate_percentage, likes_spread_percentage, group_key, search_phrase, phone, email, avg_comments_grade, avg_likes_grade, comments_rate_grade, followers_grade, engagement_rate_grade, reels_reach_grade, likes_to_comment_ratio_grade, followers_growth7d_grade, followers_growth30d_grade, followers_growth90d_grade, audience_reachability_grade, audience_authencity_grade, audience_quality_grade, post_count_grade, likes_spread_grade, followers_growth30d, similar_profile_group_data, image_reach_grade, story_reach_grade, avg_reels_play_30d, updated_at, profile_type, enabled_for_saas, plays_30d, uploads_30d, account_type, location_list, impressions_30d, reach_30d, image_impressions_grade, reels_impressions_grade, keywords_admin, search_phrase_admin, is_private, last_synced_at, ig_id) VALUES (%(name)s, %(handle)s, %(business_id)s, %(thumbnail)s, %(dob)s, %(bio)s, %(location)s, %(categories)s, %(label)s, %(languages)s, %(gender)s, %(city)s, %(state)s, %(country)s, %(keywords)s, %(is_blacklisted)s, %(followers)s, %(following)s, %(post_count)s, %(ffratio)s, %(avg_views)s, %(avg_likes)s, %(avg_reach)s, %(avg_reels_play_count)s, %(engagement_rate)s, %(followers_growth7d)s, %(reels_reach)s, %(story_reach)s, %(image_reach)s, %(avg_comments)s, %(er_grade)s, %(linked_socials)s, %(linked_channel_id)s, %(est_post_price)s, %(est_reach)s, %(est_impressions)s, %(flag_contact_info_available)s, %(is_verified)s, %(post_frequency_week)s, %(country_rank)s, %(category_rank)s, %(authentic_engagement)s, %(comment_rate_percentage)s, %(likes_spread_percentage)s, %(group_key)s, %(search_phrase)s, %(phone)s, %(email)s, %(avg_comments_grade)s, %(avg_likes_grade)s, %(comments_rate_grade)s, %(followers_grade)s, %(engagement_rate_grade)s, %(reels_reach_grade)s, %(likes_to_comment_ratio_grade)s, %(followers_growth7d_grade)s, %(followers_growth30d_grade)s, %(followers_growth90d_grade)s, %(audience_reachability_grade)s, %(audience_authencity_grade)s, %(audience_quality_grade)s, %(post_count_grade)s, %(likes_spread_grade)s, %(followers_growth30d)s, %(similar_profile_group_data)s, %(image_reach_grade)s, %(story_reach_grade)s, %(avg_reels_play_30d)s, %(updated_at)s, %(profile_type)s, %(enabled_for_saas)s, %(plays_30d)s, %(uploads_30d)s, %(account_type)s, %(location_list)s, %(impressions_30d)s, %(reach_30d)s, %(image_impressions_grade)s, %(reels_impressions_grade)s, %(keywords_admin)s, %(search_phrase_admin)s, %(is_private)s, %(last_synced_at)s, %(ig_id)s) RETURNING id'

    postgres_cp_insert_query = 'INSERT INTO campaign_profiles (platform, platform_account_id, admin_details, user_details, on_gcc, on_gcc_app, gcc_user_account_id, updated_by, enabled, has_email, has_phone, gender, phone, dob, location, languages) VALUES (%(platform)s, %(platform_account_id)s, %(admin_details)s, %(user_details)s, false, false, NULL, %(updated_by)s, true, %(has_email)s, %(has_phone)s, %(gender)s, %(phone)s, %(dob)s, %(location)s, %(languages)s)'

    postgres_ia_update_query = 'UPDATE instagram_account SET name = %s, handle = %s, business_id = %s, thumbnail = %s, dob = %s, bio = %s, location = %s, categories = %s, label = %s, languages = %s, gender = %s, city = %s, state = %s, country = %s, keywords = %s, is_blacklisted = %s, followers = %s, following = %s, post_count = %s, ffratio = %s, avg_views = %s, avg_likes = %s, avg_reach = %s, avg_reels_play_count = %s, engagement_rate = %s, followers_growth7d = %s, reels_reach = %s, story_reach = %s, image_reach = %s, avg_comments = %s, er_grade = %s, linked_socials = %s, linked_channel_id = %s, est_post_price = %s, est_reach = %s, est_impressions = %s, flag_contact_info_available = %s, is_verified = %s, post_frequency_week = %s, country_rank = %s, category_rank = %s, authentic_engagement = %s, comment_rate_percentage = %s, likes_spread_percentage = %s, group_key = %s, search_phrase = %s, phone = %s, email = %s, avg_comments_grade = %s, avg_likes_grade = %s, comments_rate_grade = %s, followers_grade = %s, engagement_rate_grade = %s, reels_reach_grade = %s, likes_to_comment_ratio_grade = %s, followers_growth7d_grade = %s, followers_growth30d_grade = %s, followers_growth90d_grade = %s, audience_reachability_grade = %s, audience_authencity_grade = %s, audience_quality_grade = %s, post_count_grade = %s, likes_spread_grade = %s, followers_growth30d = %s, similar_profile_group_data = %s, image_reach_grade = %s, story_reach_grade = %s, avg_reels_play_30d = %s, updated_at = %s, profile_type = %s, enabled_for_saas = %s, plays_30d = %s, uploads_30d = %s, account_type = %s, location_list = %s, impressions_30d = %s, reach_30d = %s, image_impressions_grade = %s, reels_impressions_grade = %s, is_private = %s, last_synced_at = %s, search_phrase_admin = %s, keywords_admin = %s WHERE ig_id = %s RETURNING id'


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
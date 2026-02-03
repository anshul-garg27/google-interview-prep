import psycopg2
import random
from collections import defaultdict
from faker import Faker
from faker.providers import DynamicProvider
from pprint import pprint
import json
import datetime
import string

conn = psycopg2.connect(database="coffee",
                        host="localhost",
                        user="root",
                        password="Glamm@123",
                        port="5432")

# conn = psycopg2.connect(database="coffee",
#                         host="gcc-pg-stage",
#                         user="gccuser",
#                         password="GCc123",
#                         port="5432")

cursor = conn.cursor()

def create_custom_generator(provider_name, elements):
    provider = DynamicProvider(
     provider_name = provider_name,
     elements = elements,
    )
    return provider

# Text data
suffix_provider = create_custom_generator("suffix", ["Group", "Org", "Pvt Limited", "Exports"])

# Array data
categories_provider = create_custom_generator("categories", ["ART", "FASHION", "BEAUTY", "GAMING", "ENTERTAINMENT", "DANCE", "TECHNOLOGY"])
tags_provider = create_custom_generator("tags", ["Attack", "Defense", "Support", "Tank"])

labels_provider = create_custom_generator("labels", ["Brand", "Creator", "Celeb"])
gender_provider = create_custom_generator("gender", ["male", "female", "other"])
er_grade_provider = create_custom_generator("er_grade", ["good", "bad", "very good"])

countries_provider = create_custom_generator("countries", ["INDIA", "USA", "UK", "JAPAN", "CANADA", "NEPAL", "SINGAPORE"])
states_provider = create_custom_generator("states", ["DELHI", "TEXAS", "KARNATAKA", "CALIFORNIA", "KANSAS", "UTAH", "RAJASTHAN", "PUNJAB"])
cities_provider = create_custom_generator("cities", ["DELHI", "BANGALORE", "TOKOYO", "DALLAS", "SANTRA", "MUMBAI", "BARODA"])
language_provider = create_custom_generator("language", ["hi", "en", "gj","ma","ta","te"])
group_key_provider = create_custom_generator("group_key",["INSTAGRAM_FASHION_1000_200000","YOUTUBE_FASHION_1000_200000","INSTAGRAM_FASHION","YOUTUBE_FASHION","INSTAGRAM_1000_200000","YOUTUBE_1000_200000","INSTAGRAM_BEAUTY_1000_200000","YOUTUBE_BEAUTY_1000_200000","INSTAGRAM_BEAUTY","YOUTUBE_BEAUTY","INSTAGRAM_1000_200000","YOUTUBE_1000_200000"])
post_type_provider = create_custom_generator("post_type",["shorts","video","reel","image"])

def create_fake():
    fake = Faker()
    fake.add_provider(suffix_provider)
    fake.add_provider(categories_provider)
    fake.add_provider(tags_provider)
    fake.add_provider(countries_provider)
    fake.add_provider(states_provider)
    fake.add_provider(cities_provider)
    fake.add_provider(labels_provider)
    fake.add_provider(gender_provider)
    fake.add_provider(er_grade_provider)
    fake.add_provider(language_provider)
    fake.add_provider(group_key_provider)
    fake.add_provider(post_type_provider)

    return fake

data_array = []
fake = create_fake()
end_date = datetime.date(year=2015, month=1, day=1)

social_timeseries_start_date = datetime.date(year=2023, month=1, day=1)
social_timeseries_end_date = datetime.date(year=2023, month=1, day=31)

social_timeseries_end_date = fake.date_between(start_date=social_timeseries_start_date, end_date=social_timeseries_end_date)
for i in range(2):
    fake = create_fake()
    value_dict = {
        "name": fake.name() + " " + fake.suffix(),
        "handle": "@"+fake.name()+ fake.suffix(),
        "business_id": random.randint(7734454502, 7834454502),
        "ig_id": random.randint(7734454502, 7834454502),
        "thumbnail":"https://cdn.vidooly.com/images/instagram/profiles/normal/9668474879.jpeg",
        "dob": fake.date_between(start_date='-30y', end_date=end_date),
        "bio": fake.text(),
        "categories": list(set([fake.categories(), fake.categories(), fake.categories()])),
        "label": fake.labels(),
        "languages": list(set([fake.language(), fake.language(), fake.language()])),
        "description": fake.text(),
        "gender": fake.gender(),
        "city": fake.cities(),
        "state":fake.states(),
        "country":fake.countries(),
        "is_blacklisted":False,
        "profile_id":"IA_"+str(i+1),
        "followers":random.randint(7734454502, 7834454502),
        "following":random.randint(7734454502, 7834454502),
        "post_count":random.randint(7734454502, 7834454502),
        "likes_count":random.randint(7734454502, 7834454502),
        "comments_count":random.randint(7734454502, 7834454502),
        "ffratio":random.uniform(10.5, 75.5),
        "avg_views":random.uniform(10.5, 75.5),
        "avg_likes":random.uniform(10.5, 75.5),
        "avg_reach":random.uniform(10.5, 75.5),
        "avg_reels_play_count":random.uniform(10.5, 75.5),
        "engagement_rate":random.uniform(10.5, 75.5),
        "followers_growth7d":random.uniform(10.5, 75.5),
        "reels_reach":random.randint(7734454502, 7834454502),
        "story_reach":random.randint(7734454502, 7834454502),
        "image_reach":random.randint(7734454502, 7834454502),
        "avg_comments":random.uniform(10.5, 75.5),
        "er_grade":fake.er_grade(),
        "profile_admin_details": json.dumps({
            "gender": fake.gender(),
            "phone": str(random.randint(9_000_000_00, 9_900_000_00))+"4",
            "name": fake.name()
        }),
        "profile_user_details":json.dumps({
            "gender": fake.gender(),
            "phone": str(random.randint(9_000_000_00, 9_900_000_00))+"4",
            "name": fake.name()
        }),
        "linked_socials":json.dumps({
            "code": random.randint(7734454502, 7834454502),
            "platform": random.choice(["YOUTUBE","INSTAGRAM"]),
            "thumbnail":"https://cdn.vidooly.com/images/instagram/profiles/normal/9668474879.jpeg",
            "followers":random.randint(7734454502, 7834454502),
            "username": "@"+fake.name(),
            "engagementRate":random.uniform(1.5, 7.5),
        }),
        "audience_gender":json.dumps( {"male_per": random.randint(10, 50), "female_per": random.randint(10, 40)}),
        "audience_location":json.dumps({
        "city": {
                fake.city(): random.randint(10, 50),
                fake.city(): random.randint(10, 50),
                fake.city(): random.randint(10, 50),
            },
            "country": {
                fake.country(): random.randint(10, 50),
                fake.country(): random.randint(10, 50),
                fake.country(): random.randint(10, 50),
            }
        }),
        "est_post_price":json.dumps({"reel_posts": random.randint(5000, 10000), "image_posts": random.randint(5000, 10000), "story_posts": random.randint(5000, 10000)}),
        "est_reach":json.dumps({"reel_posts": random.randint(5000, 10000), "image_posts": random.randint(5000, 10000), "story_posts": random.randint(5000, 10000)}),
        "est_impressions":json.dumps({"reel_posts": random.randint(5000, 10000), "image_posts": random.randint(5000, 10000), "story_posts": random.randint(5000, 10000)}),
        "est_post_price_youtube":json.dumps({"video_posts": random.randint(5000, 10000), "short_posts": random.randint(5000, 10000)}),
        "est_reach_youtube":json.dumps({"video_posts": random.randint(5000, 10000), "short_posts": random.randint(5000, 10000)}),
        "est_impressions_youtube":json.dumps({"video_posts": random.randint(5000, 10000), "short_posts": random.randint(5000, 10000)}),
        "channel_id": "UC"+str(''.join(random.choices(string.ascii_letters, k=20))),
        "title": str(''.join(random.choices(string.ascii_letters, k=20))),
        "username":"@"+fake.name(),
        "uploads_count": random.randint(100,2000),
        "views_count": random.randint(100,2000),
        "video_reach": random.randint(1000,20000),
        "shorts_reach": random.randint(1000,20000),
        "partner_id": random.randint(1000,200000),
        "activity_type":"search",
        "params":json.dumps({"size": "12", "cursor": "2", "sortBy": "id", "sortDir": "DESC"}),
        "filters":json.dumps([{"field": "audience_location.country", "value": "INDIA", "filter_type": "EQ"},{"field": "audience_age", "value": "1995-01-01", "filter_type": "GTE"},{"field": "audience_age", "value": "2021-01-01", "filter_type": "LTE"},{"field": "audience_gender", "value": "MALE", "filter_type": "EQ"},{"field": "categories", "value": "BEAUTY,HOW_TO_STYLE", "filter_type": "IN"},{"field": "keyword", "value": "virat", "filter_type": "LIKE"},{"field": "influencer_location", "value": "INDIA", "filter_type": "EQ"},{"field": "influencer_gender", "value": "MALE", "filter_type": "EQ"},{"field": "influencer_language", "value": "english", "filter_type": "EQ"},{"field": "estd_post_price", "value": "34234", "filter_type": "GTE"},{"field": "estd_post_price", "value": "98094", "filter_type": "LTE"},{"field": "avg_likes", "value": "123", "filter_type": "GTE"},{"field": "avg_likes", "value": "23233", "filter_type": "LTE"},{"field": "estd_reach", "value": "123", "filter_type": "GTE"},{"field": "estd_reach", "value": "23233", "filter_type": "LTE"},{"field": "estd_impression", "value": "123", "filter_type": "GTE"},{"field": "estd_impression", "value": "23233", "filter_type": "LTE"}]),
        "rank":i+1,
        "month":datetime.datetime.strptime('2023/01/01', '%Y/%m/%d'),
        "language":fake.language(),
        "category": fake.categories(),
        "platform":random.choice(["YOUTUBE","INSTAGRAM"]),
        "platform_profile_id": i+1,
        "post_id":random.randint(1000,20000),
        "post_link":"https://cdn.vidooly.com/images/instagram/profiles/normal/9668474879.jpeg",
        "published_at": datetime.datetime.strptime('2021/02/01 12:21:34', '%Y/%m/%d %H:%M:%S'),
        "hashtags":json.dumps({str(''.join(random.choices(string.ascii_letters, k=10))): random.randint(10,50), str(''.join(random.choices(string.ascii_letters, k=10))): random.randint(10,50),str(''.join(random.choices(string.ascii_letters, k=10))): random.randint(10,50), str(''.join(random.choices(string.ascii_letters, k=10))): random.randint(10,50)}),
        "tags": [fake.tags(),fake.tags()],
        "keywords": [fake.name],
        "audience_age": json.dumps({"13-17":random.randint(5,15),"18-24": random.randint(5,15),"25-34":random.randint(5,15),"35-44":random.randint(5,15),"45-54":random.randint(5,15),"55-64":random.randint(5,15),"65+":random.randint(5,15)}),
        "created_by": fake.name(),
        "social_timeseries_date": fake.date_between(start_date=social_timeseries_start_date, end_date=social_timeseries_end_date),
        "flag_featured": False,
        "flag_contact_info_available": True,
        "result_count": random.randint(10000,200000),
        "quality_audience_count": random.randint(10000,200000),
        "post_frequencey_week":random.uniform(1.5, 7.5),
        "country_rank":random.randint(1, 1000),
        "category_rank":json.dumps({"FASHION": random.randint(1,1000), "BEAUTY": random.randint(1,1000), "ART": random.randint(1,1000)}),
        "authenticate_engagement" :random.randint(10000,200000),
        "comment_rate_percentage": random.uniform(1.5, 7.5),
        "audience_reachability_percentage": random.uniform(1.5, 7.5),
        "audience_authenticity_percentage": random.uniform(1.5, 7.5),
        "engagement_rate":random.uniform(1.5, 7.5),
        "engagement_rate_distribution":json.dumps({"0.0-1.6":random.uniform(1.5, 7.5),"1.6-5.2":random.uniform(1.5, 7.5),"5.2-15.2":random.uniform(1.5, 7.5),"15.2-25.2":random.uniform(1.5, 7.5),"25.2-35.2":random.uniform(1.5, 7.5),"35.2-45.2":random.uniform(1.5, 7.5)}),
        "comment_rate_percentage":random.uniform(1.5, 7.5),
        "like_comment_percentage":random.uniform(1.5, 7.5),
        "likes_spread_percentage":random.uniform(1.5, 7.5),
        "city_wise_audience_location":json.dumps({
                fake.city(): random.randint(10, 50),
                fake.city(): random.randint(10, 50),
                fake.city(): random.randint(10, 50),
            }),
        "country_wise_audience_location": json.dumps({
                fake.country(): random.randint(10, 50),
                fake.country(): random.randint(10, 50),
                fake.country(): random.randint(10, 50),
            }),
        "audience_language":json.dumps({fake.language():random.uniform(1.5, 7.5),fake.language():random.uniform(1.5, 7.5),fake.language():random.uniform(1.5, 7.5)}),
        "notable_followers":json.dumps([{"code": 1,"name": "SAmple Name","handle": "@dssdssf","platform": "INSTAGRAM","followers": 2323,"thumbnail": "https://imgur.com","engagement_rate": 23.3}]),
        "quality_audience_score":random.randint(60, 100),
        "quality_score_grade":fake.er_grade(),
        "insta_quality_score_breakup":json.dumps({"AUDIENCE_AUTHENTICITY": 40, "LIKES_ACTIVITY": 80, "COMMENTS_ACTIVITY": 100}),
        "yt_quality_score_breakup":json.dumps({"CREATOR": 30, "AUDIENCE": 60, "CREDIBILITY": 90, "ENGAGEMENT": 100}),
        "group_key":fake.group_key(),
        "share_id":str(''.join(random.choices(string.ascii_letters, k=10))),
        "profile_collection_id":2,
        "custom_columns":None,
        "ordering":random.randint(1, 100),
        "video_views_last30":random.randint(1, 10000),
        "audience_age_gender_split":json.dumps({"male":{"13-17":random.randint(5,15),"18-24": random.randint(5,15),"25-34":random.randint(5,15),"35-44":random.randint(5,15),"45-54":random.randint(5,15),"55-64":random.randint(5,15),"65+":random.randint(5,15)},"female":{"13-17":random.randint(5,15),"18-24": random.randint(5,15),"25-34":random.randint(5,15),"35-44":random.randint(5,15),"45-54":random.randint(5,15),"55-64":random.randint(5,15),"65+":random.randint(5,15)}}),
        "post_type":fake.post_type(),
        "instagram_platform":"INSTAGRAM",
        "youtube_platform":"YOUTUBE",

    }
    value_dict["location"] = list(set(["city_"+value_dict["city"], "state_"+value_dict["state"],"country_"+value_dict["country"]])),
    data_array.append(value_dict)

def insertInstragramAccountData():
    query = """ INSERT INTO instagram_account (
                                                    name,
                                                    handle,
                                                    business_id,
                                                    ig_id,
                                                    thumbnail,
                                                    dob,
                                                    location,
                                                    bio,
                                                    categories,
                                                    label,
                                                    languages,
                                                    description,
                                                    gender,
                                                    city,
                                                    state,
                                                    country,
                                                    keywords,
                                                    is_blacklisted,
                                                    profile_id,
                                                    followers,
                                                    following,
                                                    post_count,
                                                    likes_count,
                                                    comments_count,
                                                    ffratio,
                                                    avg_views,
                                                    avg_likes,
                                                    avg_reach,
                                                    avg_reels_play_count,
                                                    engagement_rate,
                                                    followers_growth7d,
                                                    reels_reach,
                                                    story_reach,
                                                    image_reach,
                                                    avg_comments,
                                                    er_grade,
                                                    profile_admin_details,
                                                    profile_user_details,
                                                    linked_socials,
                                                    audience_gender,
                                                    audience_age,
                                                    audience_location,
                                                    est_post_price,
                                                    est_reach,
                                                    est_impressions,
                                                    flag_contact_info_available,
                                                    post_frequencey_week,
                                                    country_rank,
                                                    category_rank,
                                                    authenticate_engagement,
                                                    comment_rate_percentage,
                                                    likes_spread_percentage,
                                                    primary_group_key) VALUES 
    ( %(name)s, %(handle)s, %(business_id)s, %(ig_id)s, %(thumbnail)s, %(dob)s, %(location)s, %(bio)s, %(categories)s, %(label)s, %(languages)s, %(description)s, %(gender)s, %(city)s, %(state)s, %(country)s, %(keywords)s, %(is_blacklisted)s, %(profile_id)s, %(followers)s, %(following)s, %(post_count)s, %(likes_count)s, %(comments_count)s, %(ffratio)s, %(avg_views)s, %(avg_likes)s, %(avg_reach)s, %(avg_reels_play_count)s, %(engagement_rate)s, %(followers_growth7d)s, %(reels_reach)s, %(story_reach)s, %(image_reach)s, %(avg_comments)s, %(er_grade)s, %(profile_admin_details)s, %(profile_user_details)s, %(linked_socials)s, %(audience_gender)s, %(audience_age)s ,%(audience_location)s, %(est_post_price)s, %(est_reach)s, %(est_impressions)s, %(flag_contact_info_available)s, %(post_frequencey_week)s, %(country_rank)s, %(category_rank)s, %(authenticate_engagement)s, %(comment_rate_percentage)s, %(likes_spread_percentage)s, %(group_key)s)"""
    cursor.execute(query, data)
    conn.commit()

def insertYoutubeAccountData():
    query = """ INSERT INTO youtube_account (
                                                    channel_id,
                                                    title,
                                                    username,
                                                    gender,
                                                    profile_id,
                                                    description,
                                                    dob,
                                                    location,
                                                    categories,
                                                    label,
                                                    languages,
                                                    is_blacklisted,
                                                    city,
                                                    state,
                                                    country,
                                                    keywords,
                                                    likes_count,
                                                    uploads_count,
                                                    comments_count,
                                                    followers,
                                                    views_count,
                                                    avg_views,
                                                    video_reach,
                                                    shorts_reach,
                                                    thumbnail,
                                                    profile_admin_details,
                                                    profile_user_details,
                                                    linked_socials,
                                                    audience_gender,
                                                    audience_age,
                                                    audience_location,
                                                    est_post_price,
                                                    est_reach,
                                                    est_impressions,
                                                    flag_contact_info_available,
                                                    followers_growth7d,
                                                    post_frequencey_week,
                                                    country_rank,
                                                    category_rank,
                                                    authenticate_engagement,
                                                    comment_rate_percentage,
                                                    likes_spread_percentage,
                                                    video_views_last30,
                                                    primary_group_key) VALUES
    ( %(channel_id)s, %(title)s, %(username)s, %(gender)s, %(profile_id)s, %(description)s, %(dob)s, %(location)s, %(categories)s, %(label)s, %(languages)s, %(is_blacklisted)s, %(city)s, %(state)s, %(country)s, %(keywords)s, %(likes_count)s, %(uploads_count)s, %(comments_count)s, %(followers)s, %(views_count)s, %(avg_views)s, %(video_reach)s, %(shorts_reach)s, %(thumbnail)s, %(profile_admin_details)s, %(profile_user_details)s, %(linked_socials)s, %(audience_gender)s, %(audience_age)s, %(audience_location)s, %(est_post_price_youtube)s, %(est_reach_youtube)s, %(est_impressions_youtube)s, %(flag_contact_info_available)s, %(followers_growth7d)s, %(post_frequencey_week)s, %(country_rank)s, %(category_rank)s, %(authenticate_engagement)s, %(comment_rate_percentage)s, %(likes_spread_percentage)s, %(video_views_last30)s, %(group_key)s)"""
    cursor.execute(query, data)
    conn.commit()

def insertActivityData():
    query = """ INSERT INTO activity (
                                        partner_id,
                                        title,
                                        activity_type,
                                        params,
                                        filters,
                                        result_count) VALUES
        ( %(partner_id)s, %(title)s, %(activity_type)s, %(params)s,%(filters)s,%(result_count)s)"""
    cursor.execute(query, data)
    conn.commit()

def insertLeaderboardData():
    
    query = """ INSERT INTO leaderboard (
                                        rank,
                                        month,
                                        language,
                                        category,
                                        platform,
                                        platform_profile_id,
                                        profile_id,
                                        followers,
                                        engagement_rate,
                                        avg_likes,
                                        avg_comments,
                                        avg_views,
                                        er_grade) VALUES
        ( %(rank)s, %(month)s, %(language)s, %(category)s, %(platform)s, %(platform_profile_id)s, %(profile_id)s, %(followers)s, %(engagement_rate)s, %(avg_likes)s, %(avg_comments)s, %(avg_views)s, %(er_grade)s)"""
    cursor.execute(query, data)
    conn.commit()

def insertSocialProfileTimeSeriesData():
    query = """ INSERT INTO social_profile_time_series (
                                        platform_profile_id,
                                        platform,
                                        followers,
                                        following,
                                        engagement_rate,
                                        date) VALUES
        ( %(platform_profile_id)s, %(platform)s, %(followers)s, %(following)s, %(engagement_rate)s, %(social_timeseries_date)s)"""
    cursor.execute(query, data)
    conn.commit()

def insertSocialProfilePosts():
    query = """ INSERT INTO social_profile_posts (
                                        platform_profile_id,
                                        platform,
                                        post_id,
                                        post_link,
                                        post_type,
                                        thumbnail,
                                        likes_count,
                                        comments_count,
                                        views_count,
                                        engagement_rate,
                                        published_at) VALUES
        ( %(platform_profile_id)s, %(platform)s, %(post_id)s, %(post_link)s, %(post_type)s, %(thumbnail)s, %(likes_count)s, %(comments_count)s, %(views_count)s, %(engagement_rate)s, %(published_at)s)"""
    cursor.execute(query, data)
    conn.commit()

def insertSocialProfileHashtags():
    query = """ INSERT INTO social_profile_hashtags (
                                        platform_profile_id,
                                        platform,
                                        hashtags) VALUES
        ( %(platform_profile_id)s, %(platform)s, %(hashtags)s)"""
    cursor.execute(query, data)
    conn.commit()

def insertSocialProfileAudienceInfoInstagram():
    query = """ INSERT INTO social_profile_audience_info (
                                        platform_profile_id,
                                        platform,
                                        city_wise_audience_location,
                                        country_wise_audience_location,
                                        audience_gender_split,
                                        audience_age_gender_split,
                                        audience_language,
                                        notable_followers,
                                        audience_reachability_percentage,
                                        audience_authenticity_percentage,
                                        quality_audience_score,
                                        quality_score_grade,
                                        quality_score_breakup) VALUES
        ( %(platform_profile_id)s, %(instagram_platform)s, %(city_wise_audience_location)s, %(country_wise_audience_location)s,  %(audience_gender)s, %(audience_age_gender_split)s, %(audience_language)s, %(notable_followers)s, %(audience_reachability_percentage)s, %(audience_authenticity_percentage)s, %(quality_audience_score)s, %(quality_score_grade)s, %(insta_quality_score_breakup)s)"""
    cursor.execute(query, data)
    conn.commit()

def insertSocialProfileAudienceInfoYoutube():
    query = """ INSERT INTO social_profile_audience_info (
                                        platform_profile_id,
                                        platform,
                                        city_wise_audience_location,
                                        country_wise_audience_location,
                                        audience_gender_split,
                                        audience_age_gender_split,
                                        audience_language,
                                        notable_followers,
                                        audience_reachability_percentage,
                                        audience_authenticity_percentage,
                                        quality_audience_score,
                                        quality_score_grade,
                                        quality_score_breakup) VALUES
        ( %(platform_profile_id)s, %(youtube_platform)s, %(city_wise_audience_location)s, %(country_wise_audience_location)s,  %(audience_gender)s, %(audience_age_gender_split)s, %(audience_language)s, %(notable_followers)s, %(audience_reachability_percentage)s, %(audience_authenticity_percentage)s, %(quality_audience_score)s, %(quality_score_grade)s, %(yt_quality_score_breakup)s)"""
    cursor.execute(query, data)
    conn.commit()

def insertSocialProfileGroupData():
    query = """ INSERT INTO social_profile_group_data (
                                    group_key, 
                                    audience_reachability_percentage,
                                    audience_authenticity_percentage,
                                    engagement_rate,
                                    engagement_rate_distribution,
                                    comment_rate_percentage,
                                    like_comment_percentage,
                                    likes_spread_percentage) VALUES
        ( %(group_key)s, %(audience_reachability_percentage)s, %(audience_authenticity_percentage)s, %(engagement_rate)s, %(engagement_rate_distribution)s, %(comment_rate_percentage)s, %(like_comment_percentage)s, %(likes_spread_percentage)s)"""
    cursor.execute(query, data)
    conn.commit()

def insertProfileCollectionData():
    query = """ INSERT INTO profile_collection (
                                    name,
                                    partner_id,
                                    description,
                                    categories,
                                    tags,
                                    featured,
                                    created_by,
                                    share_id) VALUES
        ( %(name)s, %(partner_id)s, %(description)s, %(categories)s, %(tags)s, %(flag_featured)s, %(created_by)s, %(share_id)s)"""
    cursor.execute(query, data)
    conn.commit()

def insertProfileCollectionItem():
    query = """ INSERT INTO profile_collection_item (
                                        platform,
                                        platform_profile_id,
                                        profile_collection_id,
                                        custom_columns,
                                        ordering,
                                        created_by) VALUES
        ( %(platform)s, %(platform_profile_id)s, %(profile_collection_id)s, %(custom_columns)s, %(ordering)s, %(created_by)s)"""
    cursor.execute(query, data)
    conn.commit()

for data in data_array:
    # insertInstragramAccountData()
    # insertYoutubeAccountData()
    # insertActivityData()
    # insertLeaderboardData()
    #  insertSocialProfileTimeSeriesData()
    # insertSocialProfilePosts()
    # insertSocialProfileHashtags()
    insertSocialProfileAudienceInfoInstagram()
    insertSocialProfileAudienceInfoYoutube()
    # insertSocialProfileGroupData()
    # insertProfileCollectionData()
    # insertProfileCollectionItem()
{{ config(materialized = 'table', tags=["hourly"], order_by='channel_id') }}
with
tagged_channels as (
     select 0 "priority", channelid, "Creator/Brand", "language", "gcc_category", "gender", "city", "state", "linked_instagram_handle", "email"
        FROM s3('https://gcc-social-data.s3.ap-south-1.amazonaws.com/temp/manual_data/overrides.csv', 'AKIAXGXUCIER5DX7YI4F', 'WnVvyEs+3NFFcNWnEGJr3Ctt/KuMUmaRYG9xbPWj', 'CSVWithNames')
    UNION ALL
    select 1 "priority", channelid, "Creator/Brand", "language", "gcc_category", "gender", "city", "state", "linked_instagram_handle", "email"
        FROM s3('https://gcc-social-data.s3.ap-south-1.amazonaws.com/temp/manual_data/mar.csv', 'AKIAXGXUCIER5DX7YI4F', 'WnVvyEs+3NFFcNWnEGJr3Ctt/KuMUmaRYG9xbPWj', 'CSVWithNames')
    UNION ALL
    select 2 "priority", channelid, "Creator/Brand", "language", "gcc_category", "gender", "city", "state", "linked_instagram_handle", "email"
        FROM s3('https://gcc-social-data.s3.ap-south-1.amazonaws.com/temp/manual_data/april.csv', 'AKIAXGXUCIER5DX7YI4F', 'WnVvyEs+3NFFcNWnEGJr3Ctt/KuMUmaRYG9xbPWj', 'CSVWithNames')
    UNION ALL
    select 3 "priority", channelid, "Creator/Brand", "language", "gcc_category", "gender", "city", "state", "linked_instagram_handle", "email"
        FROM s3('https://gcc-social-data.s3.ap-south-1.amazonaws.com/temp/manual_data/may.csv', 'AKIAXGXUCIER5DX7YI4F', 'WnVvyEs+3NFFcNWnEGJr3Ctt/KuMUmaRYG9xbPWj', 'CSVWithNames')
    UNION ALL
    select 4 "priority", channelid, "Creator/Brand", "language", "gcc_category", "gender", "city", "state", "linked_instagram_handle", "email"
        FROM s3('https://gcc-social-data.s3.ap-south-1.amazonaws.com/temp/manual_data/june_news.csv', 'AKIAXGXUCIER5DX7YI4F', 'WnVvyEs+3NFFcNWnEGJr3Ctt/KuMUmaRYG9xbPWj', 'CSVWithNames')
    UNION ALL
    select 5 "priority", channelid, "Creator/Brand", "language", "gcc_category", "gender", "city", "state", "linked_instagram_handle", "email"
        FROM s3('https://gcc-social-data.s3.ap-south-1.amazonaws.com/temp/manual_data/jul1.csv', 'AKIAXGXUCIER5DX7YI4F', 'WnVvyEs+3NFFcNWnEGJr3Ctt/KuMUmaRYG9xbPWj', 'CSVWithNames')
    UNION ALL
    select 6 "priority", channelid, "Creator/Brand", "language", "gcc_category", "gender", "city", "state", "linked_instagram_handle", "email"
        FROM s3('https://gcc-social-data.s3.ap-south-1.amazonaws.com/temp/manual_data/jul2.csv', 'AKIAXGXUCIER5DX7YI4F', 'WnVvyEs+3NFFcNWnEGJr3Ctt/KuMUmaRYG9xbPWj', 'CSVWithNames')
    UNION ALL
    select 7 "priority", channelid, "Creator/Brand", "language", "gcc_category", "gender", "city", "state", "linked_instagram_handle", "email"
        FROM s3('https://gcc-social-data.s3.ap-south-1.amazonaws.com/temp/manual_data/jul3.csv', 'AKIAXGXUCIER5DX7YI4F', 'WnVvyEs+3NFFcNWnEGJr3Ctt/KuMUmaRYG9xbPWj', 'CSVWithNames')
    UNION ALL
    select 8 "priority", channelid, "Creator/Brand", "language", "gcc_category", "gender", "city", "state", "linked_instagram_handle", "email"
        FROM s3('https://gcc-social-data.s3.ap-south-1.amazonaws.com/temp/manual_data/aug1.csv', 'AKIAXGXUCIER5DX7YI4F', 'WnVvyEs+3NFFcNWnEGJr3Ctt/KuMUmaRYG9xbPWj', 'CSVWithNames')
    UNION ALL
    select 9 "priority", channelid, "Creator/Brand", "language", "gcc_category", "gender", "city", "state", "linked_instagram_handle", "email"
        FROM s3('https://gcc-social-data.s3.ap-south-1.amazonaws.com/temp/manual_data/aug2.csv', 'AKIAXGXUCIER5DX7YI4F', 'WnVvyEs+3NFFcNWnEGJr3Ctt/KuMUmaRYG9xbPWj', 'CSVWithNames')
    UNION ALL
    select 10 "priority", channelid, "Creator/Brand", "language", "gcc_category", "gender", "city", "state", "linked_instagram_handle", "email"
        FROM s3('https://gcc-social-data.s3.ap-south-1.amazonaws.com/temp/manual_data/sep.csv', 'AKIAXGXUCIER5DX7YI4F', 'WnVvyEs+3NFFcNWnEGJr3Ctt/KuMUmaRYG9xbPWj', 'CSVWithNames')
),
valid_states as (
    select state from dbt.mart_youtube_account group by state having count(*) > 2
),
raw as (
    select
        ifNull(channelid, 'MISSING') channel_id,
        multiIf (`Creator/Brand` = 'Creator', 'CREATOR', `Creator/Brand` = 'Brand', 'BRAND', NULL) profile_type,
        multiIf(lower(language) = 'english', 'en', lower(language) = 'assamese', 'as', lower(language) = 'manipuri', 'mni', lower(language) = 'bodo', 'brx', lower(language) = 'nagamese', 'nag', lower(language) = 'bhojpuri', 'bho', lower(language) = 'rajasthani', 'raj', lower(language) = 'nepali', 'ne', lower(language) = 'khasi', 'kha', lower(language) = 'mizo', 'lus',lower(language) = 'urdu', 'ur', lower(language) = 'hindi', 'hi', lower(language) = 'bengali', 'bn', lower(language) = 'tamil', 'ta', lower(language) = 'telugu', 'te', lower(language) = 'malayalam', 'ml', lower(language) = 'marathi', 'mr', lower(language) = 'kannada', 'ka', lower(language) = 'gujarati', 'gu', lower(language) = 'punjabi', 'pa', lower(language) = 'odia', 'or', NULL) language,
        multiIf(
        gcc_category = 'Music', 'Music',
        gcc_category =     'Entertainment', 'Entertainment',
        gcc_category =     'Comedy', 'Comedy',
        gcc_category =     'Gaming', 'Gaming',
        gcc_category =     'miscellaneous', 'Miscellaneous',
        gcc_category =     'Vlogging', 'Vloging',
        gcc_category =     'Infotainment', 'Infotainment',
        gcc_category =     'Education', 'Education',
        gcc_category =     'Religious', 'Religious Content',
        gcc_category =     'Kids & Animation', 'Kids & Animation',
        gcc_category =     'News & Politics', 'News & Politics',
        gcc_category =     'Science & Technology', 'Science & Technology',
        gcc_category =     'Food & Recipe', 'Food & Drinks',
        gcc_category =     'Motivational', 'Motivational',
        gcc_category =     'Health & Fitness', 'Health & Fitness',
        gcc_category =     'Business & Finance', 'Finance',
        gcc_category =     'Travel & Leisure', 'Travel & Leisure',
        gcc_category =     'Other', 'Miscellaneous',
        gcc_category =     'Sports', 'Sports',
        gcc_category =     'DIY & Home Decor', 'DIY',
        gcc_category =     'Photography & Editing', 'Photography & Editing',
        gcc_category =     'Fashion & Style', 'Fashion & Style',
        gcc_category =     'Auto & Vehicles', 'Autos & Vehicles',
        gcc_category =     'Agriculture & Allied Sectors', 'Agriculture & Allied Sectors',
        gcc_category =     'Astrology', 'Astrology',
        gcc_category =     'Pets & Animals', 'Pets & Animals',
        gcc_category =     'Beauty', 'Beauty',
        gcc_category =     'Adult', 'Adult',
        gcc_category =     'Supernatural', 'Supernatural',
        gcc_category =     'Real Estate', 'Real Estate',
        gcc_category =     'Defence', 'Defence',
        gcc_category =     'People & Culture', 'People & Culture',
        gcc_category =     'IT & ITES', 'IT & ITES',
        gcc_category =     'Family & Parenting', 'Family & Parenting',
        gcc_category =     'Arts_Entertainment', 'Arts & Craft',
        gcc_category =     'Government', 'Government',
        gcc_category =     'Non profit & activism', 'Non profits',
        gcc_category =     'Business', 'Finance',
        gcc_category =     'Fashion_Style', 'Fashion & Style',
        gcc_category =     'Careers', 'Miscellaneous',
        gcc_category =     'Personal_Finance', 'Finance',
        gcc_category =     'Retail & Shopping', 'Miscellaneous',
        gcc_category =     'Technology_Computing', 'Science & Technology',
        gcc_category =     'Personal Finance', 'Finance',
        gcc_category =     'FMCG', 'Miscellaneous',
        gcc_category =     'Real_Estate', 'Real Estate',
        gcc_category =     'Automotive', 'Autos & Vehicles',
        gcc_category =     'Food & Drink', 'Food & Drinks',
        gcc_category =     'Health_Fitness', 'Health & Fitness',
        gcc_category =     'Home_Garden', 'Home & Garden',
        gcc_category =     'comedy', 'Comedy',
        gcc_category =     'Home & Garden', 'Home & Garden',
        gcc_category =     'Food_Drink', 'Food & Drinks',
        gcc_category =     'NGO', 'Non profits',
        NULL
        ) category,
        multiIf(gender = 'Male' or gender = 'M' or gender = 'MALE', 'MALE', gender = 'Female' or gender = 'F' or gender = 'FEMALE', 'FEMALE', NULL) gender,
        city,
        multiIf(vs.state = 'Tamilnadu', 'Tamil Nadu', vs.state = 'Madhy Pradesh', 'Madhya Pradesh', vs.state = 'UP', 'Uttar Pradesh', vs.state) state, 
        splitByChar('/', ifNull(linked_instagram_handle, ''))[-1] linked_instagram_handle,
        email,
        priority
    from tagged_channels
    left join valid_states vs on vs.state = tagged_channels.state
),
final as (
    select
        channel_id,
        max(profile_type) profile_type,
        argMax(language, priority) language,
        argMax(category, priority) category,
        max(gender) gender,
        max(city) city,
        max(state) state,
        max(linked_instagram_handle) linked_instagram_handle,
        max(email) email
    from raw
    group by channel_id
)
select * from final

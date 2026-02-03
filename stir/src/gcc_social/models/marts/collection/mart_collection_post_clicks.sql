{{ config(materialized = 'table', tags=["hourly"], order_by="handle, post_collection_id") }}
with raw_csv_data as (
    SELECT 
     *
    FROM s3('https://gcc-social-data.s3.ap-south-1.amazonaws.com/temp/Campaign_post_links_clicks.csv', 'AKIAXGXUCIER3E4AHEPR', 'fMM+VR+d5D43ydFHkbHInwuGyM8/t+dJRWSMW0yE', 'CSVWithNames')
),
processed_output as (
    select 
        ifNull(handle, 'MISSING') handle,
        toString(short_code) short_code,
        post_link,
        toUInt64(Clicks) Clicks,
        ifNull(post_collection_id, 'MISSING') post_collection_id
    from raw_csv_data
)
select * from processed_output

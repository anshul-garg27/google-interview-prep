{{ config(materialized = 'table', tags=["hourly"], order_by="ifNull(Handle, 'MISSING')") }}
with raw_csv_data as (
    SELECT 
     *
    FROM s3('https://gcc-social-data.s3.ap-south-1.amazonaws.com/manual_tagging/manual_tagging_9.csv', 'AKIAXGXUCIER3E4AHEPR', 'fMM+VR+d5D43ydFHkbHInwuGyM8/t+dJRWSMW0yE', 'CSVWithNames')
),
processed_output as (
    select 
        *,
        arrayFilter(x -> x is not null, array(
            if ("C1 Y/N" = 'Yes', "Category 1", NULL), 
            if ("C2 Y/N" = 'Yes', "Category 2", NULL),
            if ("C3 Y/N" = 'Yes', "Category 3", NULL),
            if ("Additional Category" != '', "Additional Category", NULL)
        )) categories,
        arrayFilter(x -> x is not null, array(
            if ("Language 1" != '', "Language 1", NULL), 
            if ("Language 2" != '', "Language 2", NULL),
            if ("Language 3" != '', "Language 3", NULL)
        )) languages
    from raw_csv_data
)
select * from processed_output

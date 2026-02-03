{{ config(materialized = 'table', tags=["gcc_orders"], order_by='platform, platform_order_id') }}
with final as (
    select
        ifNull(platform,'MISSING') platform,
        store,
        JSONExtractString(utm_params, 'source') utm_source,
        JSONExtractString(utm_params, 'medium') utm_medium,
        JSONExtractString(utm_params, 'campaign') utm_campaign,
        ifNull(platform_order_id, 'MISSING') platform_order_id,
        JSONExtractString(order_data, 'name') store_order_id,
        JSONExtractString(order_data, 'total_price') total_price,
        JSONExtractString(JSONExtractArrayRaw(ifNull(order_data, ''), 'fulfillments')[1], 'status') fulfillment_status,
        JSONExtractString(order_data, 'tags') tags,
        parseDateTimeBestEffortOrNull(JSONExtractString(order_data, 'created_at')) created_at,
        parseDateTimeBestEffortOrNull(JSONExtractString(order_data, 'updated_at')) updated_at,
        parseDateTimeBestEffortOrNull(JSONExtractString(order_data, 'cancelled_at')) cancelled_at,
        0 is_returned,
        parseDateTimeBestEffortOrNull(JSONExtractString(JSONExtractArrayRaw(ifNull(order_data, ''), 'fulfillments')[1], 'created_at')) confirmed_at,
        if(now() < confirmed_at + INTERVAL 7 DAY, NULL, confirmed_at + INTERVAL 7 DAY) delivered_at,
        multiIf(
            cancelled_at is not null, 'CANCELLED',
            fulfillment_status = 'success' and delivered_at is null, 'CONFIRMED',
            fulfillment_status = 'success' and delivered_at is not null, 'DELIVERED',
            fulfillment_status != 'success', 'PENDING',
            is_returned, 'RETURNED',
            'UNKNOWN'
        ) status,
        order_data
    from dbt.stg_beat_order
    where utm_source = 'gcc'
)
select * from final
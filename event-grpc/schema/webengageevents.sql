CREATE TABLE _e.webengage_events
(
    `id` String,
    `event_name` String,
    `account_type` String,
    `timestamp` DateTime,
    `user_id` String,
    `anonymous_id` String,
    `data` String
) ENGINE = ReplacingMergeTree(timestamp)
PARTITION BY toYYYYMMDD(timestamp)
ORDER BY id
SETTINGS index_granularity = 8192;
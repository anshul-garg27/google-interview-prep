CREATE TABLE _e.events (
	event_id varchar(55) NOT NULL,
	event_name varchar(55) NOT NULL,
	event_timestamp timestamp without time zone NOT NULL,
	insert_timestamp timestamp without time zone NOT NULL,
    session_id varchar(50),
    bb_device_id bigint,
    user_id bigint,
    device_id varchar(50),
    android_advertising_id varchar(50),
    client_id varchar(30),
    channel varchar(20),
    os varchar(30),
    client_type varchar(20),
    app_language varchar(20),
    app_version varchar(20),
    app_version_code varchar(20),
    current_url varchar,
    merchant_id varchar(20),
    utm_referrer varchar,
    utm_platform varchar(50),
    utm_source varchar(80),
    utm_medium varchar(80),
    utm_campaign varchar(80),
	event_params jsonb) PARTITION BY RANGE (event_timestamp);

CREATE TABLE _e.events_20210220 PARTITION OF _e.events FOR VALUES FROM ('2021-02-20  00:00:00') TO ('2021-02-21  00:00:00');
CREATE TABLE _e.events_20210221 PARTITION OF _e.events FOR VALUES FROM ('2021-02-21  00:00:00') TO ('2021-02-22  00:00:00');
CREATE TABLE _e.events_20210222 PARTITION OF _e.events FOR VALUES FROM ('2021-02-22  00:00:00') TO ('2021-02-23  00:00:00');
CREATE TABLE _e.events_20210223 PARTITION OF _e.events FOR VALUES FROM ('2021-02-23  00:00:00') TO ('2021-02-24  00:00:00');
CREATE TABLE _e.events_20210224 PARTITION OF _e.events FOR VALUES FROM ('2021-02-24  00:00:00') TO ('2021-02-25  00:00:00');
CREATE TABLE _e.events_20210225 PARTITION OF _e.events FOR VALUES FROM ('2021-02-25  00:00:00') TO ('2021-02-26  00:00:00');
CREATE TABLE _e.events_20210226 PARTITION OF _e.events FOR VALUES FROM ('2021-02-26  00:00:00') TO ('2021-02-27  00:00:00');
CREATE TABLE _e.events_20210227 PARTITION OF _e.events FOR VALUES FROM ('2021-02-27  00:00:00') TO ('2021-02-28  00:00:00');
CREATE TABLE _e.events_20210228 PARTITION OF _e.events FOR VALUES FROM ('2021-02-28  00:00:00') TO ('2021-03-01  00:00:00');
CREATE TABLE _e.events_20210301 PARTITION OF _e.events FOR VALUES FROM ('2021-03-01  00:00:00') TO ('2021-03-02  00:00:00');
CREATE TABLE _e.events_20210302 PARTITION OF _e.events FOR VALUES FROM ('2021-03-02  00:00:00') TO ('2021-03-03  00:00:00');
CREATE TABLE _e.events_20210303 PARTITION OF _e.events FOR VALUES FROM ('2021-03-03  00:00:00') TO ('2021-03-04  00:00:00');
CREATE TABLE _e.events_20210304 PARTITION OF _e.events FOR VALUES FROM ('2021-03-04  00:00:00') TO ('2021-03-05  00:00:00');
CREATE TABLE _e.events_20210305 PARTITION OF _e.events FOR VALUES FROM ('2021-03-05  00:00:00') TO ('2021-03-06  00:00:00');
CREATE TABLE _e.events_20210306 PARTITION OF _e.events FOR VALUES FROM ('2021-03-06  00:00:00') TO ('2021-03-07  00:00:00');
CREATE TABLE _e.events_20210307 PARTITION OF _e.events FOR VALUES FROM ('2021-03-07  00:00:00') TO ('2021-03-08  00:00:00');
CREATE TABLE _e.events_20210308 PARTITION OF _e.events FOR VALUES FROM ('2021-03-08  00:00:00') TO ('2021-03-09  00:00:00');
CREATE TABLE _e.events_20210309 PARTITION OF _e.events FOR VALUES FROM ('2021-03-09  00:00:00') TO ('2021-03-10  00:00:00');
CREATE TABLE _e.events_20210310 PARTITION OF _e.events FOR VALUES FROM ('2021-03-10  00:00:00') TO ('2021-03-11  00:00:00');
CREATE TABLE _e.events_20210311 PARTITION OF _e.events FOR VALUES FROM ('2021-03-11  00:00:00') TO ('2021-03-12  00:00:00');
CREATE TABLE _e.events_20210312 PARTITION OF _e.events FOR VALUES FROM ('2021-03-12  00:00:00') TO ('2021-03-13  00:00:00');
CREATE TABLE _e.events_20210313 PARTITION OF _e.events FOR VALUES FROM ('2021-03-13  00:00:00') TO ('2021-03-14  00:00:00');
CREATE TABLE _e.events_20210314 PARTITION OF _e.events FOR VALUES FROM ('2021-03-14  00:00:00') TO ('2021-03-15  00:00:00');
CREATE TABLE _e.events_20210315 PARTITION OF _e.events FOR VALUES FROM ('2021-03-15  00:00:00') TO ('2021-03-16  00:00:00');
CREATE TABLE _e.events_20210316 PARTITION OF _e.events FOR VALUES FROM ('2021-03-16  00:00:00') TO ('2021-03-17  00:00:00');
CREATE TABLE _e.events_20210317 PARTITION OF _e.events FOR VALUES FROM ('2021-03-17  00:00:00') TO ('2021-03-18  00:00:00');
CREATE TABLE _e.events_20210318 PARTITION OF _e.events FOR VALUES FROM ('2021-03-18  00:00:00') TO ('2021-03-19  00:00:00');
CREATE TABLE _e.events_20210319 PARTITION OF _e.events FOR VALUES FROM ('2021-03-19  00:00:00') TO ('2021-03-20  00:00:00');
CREATE TABLE _e.events_20210320 PARTITION OF _e.events FOR VALUES FROM ('2021-03-20  00:00:00') TO ('2021-03-21  00:00:00');
CREATE TABLE _e.events_20210321 PARTITION OF _e.events FOR VALUES FROM ('2021-03-21  00:00:00') TO ('2021-03-22  00:00:00');
CREATE TABLE _e.events_20210322 PARTITION OF _e.events FOR VALUES FROM ('2021-03-22  00:00:00') TO ('2021-03-23  00:00:00');

CREATE TABLE _e.trace_log_events
(
    `id` String,
    `host_name` String,
    `service_name` String,
    `timestamp` DateTime,
    `time_taken` Int,
    `method` String,
    `uri` String,
    `headers` String,
    `response_status` Int,
    `remote_address` String,
    `only_request` UInt8,
    `request_body` String,
    `response_body` String
)
    ENGINE = ReplacingMergeTree
ORDER BY id
SETTINGS index_granularity = 8192;


CREATE TABLE _e.affiliate_order_events
(
    `order_id` String,
    `platform` String,
    `referral_code` String,
    `meta` String,
    `created_on` DateTime,
    `last_modified_on` DateTime,
    `amount` Int,
    `status` String
)
    ENGINE = ReplacingMergeTree
ORDER BY (order_id, platform)
SETTINGS index_granularity = 8192;


CREATE TABLE _stage_e.trace_log_events
(
    `id` String,
    `host_name` String,
    `service_name` String,
    `timestamp` DateTime,
    `time_taken` Int,
    `method` String,
    `uri` String,
    `headers` String,
    `response_status` Int,
    `remote_address` String,
    `only_request` UInt8,
    `request_body` String,
    `response_body` String
)
    ENGINE = ReplacingMergeTree
ORDER BY id
SETTINGS index_granularity = 8192;


CREATE TABLE _stage_e.affiliate_order_events
(
    `order_id` String,
    `platform` String,
    `referral_code` String,
    `meta` String,
    `created_on` DateTime,
    `last_modified_on` DateTime,
    `amount` Int,
    `status` String
)
    ENGINE = ReplacingMergeTree
ORDER BY (order_id, platform)
SETTINGS index_granularity = 8192;


CREATE TABLE _stage_e.bigboss_vote_log
(
    `id` String,
    `key` String,
    `vendorcode` String,
    `country` String,
    `meta` String,
    `contestantid` String,
    `identifier` String,
    `languages` String,
    `sdc_received_at` String,
    `sdc_batched_at` String,
    `sdc_extracted_at` String,
    `createdat` DateTime,
    `updatedat` DateTime,
    `sdc_table_version` Int,
    `sdc_sequence` Int
)
    ENGINE = ReplacingMergeTree
ORDER BY id
SETTINGS index_granularity = 8192;

CREATE TABLE _e.bigboss_vote_log
(
    `id` String,
    `key` String,
    `vendorcode` String,
    `country` String,
    `meta` String,
    `contestantid` String,
    `identifier` String,
    `languages` String,
    `sdc_received_at` String,
    `sdc_batched_at` String,
    `sdc_extracted_at` String,
    `createdat` DateTime,
    `updatedat` DateTime,
    `sdc_table_version` Int,
    `sdc_sequence` Int
)
    ENGINE = ReplacingMergeTree
ORDER BY id
SETTINGS index_granularity = 8192;


CREATE TABLE _stage_e.post_log_events
(
    `event_id` String,
    `source` String,
    `platform` String,
    `profile_id` String,
    `handle` Nullable(String),
    `shortcode` String,
    `publish_time` Nullable(DateTime),
    `event_timestamp` DateTime,
    `insert_timestamp` DateTime,
    `metrics` String,
    `dimensions` String
)
ENGINE = MergeTree
PARTITION BY toYYYYMM(event_timestamp)
ORDER BY (platform, profile_id, shortcode, event_timestamp)
SETTINGS index_granularity = 8192;



CREATE TABLE _e.post_log_events
(
    `event_id` String,
    `source` String,
    `platform` String,
    `profile_id` String,
    `handle` Nullable(String),
    `shortcode` String,
    `publish_time` Nullable(DateTime),
    `event_timestamp` DateTime,
    `insert_timestamp` DateTime,
    `metrics` String,
    `dimensions` String
)
ENGINE = MergeTree
PARTITION BY toYYYYMM(event_timestamp)
ORDER BY (platform, profile_id, shortcode, event_timestamp)
SETTINGS index_granularity = 8192;


CREATE TABLE _stage_e.profile_log_events
(
    `event_id` String,
    `source` String,
    `platform` String,
    `profile_id` String,
    `handle` Nullable(String),
    `event_timestamp` DateTime,
    `insert_timestamp` DateTime,
    `metrics` String,
    `dimensions` String
)
    ENGINE = MergeTree
PARTITION BY toYYYYMM(event_timestamp)
ORDER BY (platform, profile_id, event_timestamp)
SETTINGS index_granularity = 8192;



CREATE TABLE _e.profile_log_events
(
    `event_id` String,
    `source` String,
    `platform` String,
    `profile_id` String,
    `handle` Nullable(String),
    `event_timestamp` DateTime,
    `insert_timestamp` DateTime,
    `metrics` String,
    `dimensions` String
)
    ENGINE = MergeTree
PARTITION BY toYYYYMM(event_timestamp)
ORDER BY (platform, profile_id, event_timestamp)
SETTINGS index_granularity = 8192;


CREATE TABLE _stage_e.scrape_request_log_events
(
    `event_id` String,
    `scl_id` UInt64,
    `platform` String,
    `flow` String,
    `params` String,
    `status` String,
    `priority` Nullable(Int32),
    `reason` Nullable(String),
    `picked_at` Nullable(DateTime),
    `expires_at` Nullable(DateTime),
    `event_timestamp` DateTime,
    `insert_timestamp` DateTime
)
    ENGINE = MergeTree
PARTITION BY (toYYYYMMDD(event_timestamp), flow)
ORDER BY (platform, flow, event_timestamp)
SETTINGS index_granularity = 8192;


CREATE TABLE _e.scrape_request_log_events
(
    `event_id` String,
    `scl_id` UInt64,
    `platform` String,
    `flow` String,
    `params` String,
    `status` String,
    `priority` Nullable(Int32),
    `reason` Nullable(String),
    `picked_at` Nullable(DateTime),
    `expires_at` Nullable(DateTime),
    `event_timestamp` DateTime,
    `insert_timestamp` DateTime
)
    ENGINE = MergeTree
PARTITION BY (toYYYYMMDD(event_timestamp), flow)
ORDER BY (platform, flow, event_timestamp)
SETTINGS index_granularity = 8192;


CREATE TABLE _stage_e.order_log_events
(
    `event_id` String,
    `source` String,
    `platform` String,
    `platform_order_id` String,
    `store` String,
    `event_timestamp` DateTime,
    `insert_timestamp` DateTime,
    `metrics` String,
    `dimensions` String
)
    ENGINE = MergeTree
PARTITION BY toYYYYMM(event_timestamp)
ORDER BY (platform, store, platform_order_id, event_timestamp)
SETTINGS index_granularity = 8192;


CREATE TABLE _e.order_log_events
(
    `event_id` String,
    `source` String,
    `platform` String,
    `platform_order_id` String,
    `store` String,
    `event_timestamp` DateTime,
    `insert_timestamp` DateTime,
    `metrics` String,
    `dimensions` String
)
    ENGINE = MergeTree
PARTITION BY toYYYYMM(event_timestamp)
ORDER BY (platform, store, platform_order_id, event_timestamp)
SETTINGS index_granularity = 8192;
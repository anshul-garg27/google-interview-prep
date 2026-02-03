CREATE TABLE _stage_events.events (
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

CREATE TABLE _stage_events.events_20210220 PARTITION OF _stage_events.events FOR VALUES FROM ('2021-02-20  00:00:00') TO ('2021-02-21  00:00:00');
CREATE TABLE _stage_events.events_20210221 PARTITION OF _stage_events.events FOR VALUES FROM ('2021-02-21  00:00:00') TO ('2021-02-22  00:00:00');
CREATE TABLE _stage_events.events_20210222 PARTITION OF _stage_events.events FOR VALUES FROM ('2021-02-22  00:00:00') TO ('2021-02-23  00:00:00');
CREATE TABLE _stage_events.events_20210223 PARTITION OF _stage_events.events FOR VALUES FROM ('2021-02-23  00:00:00') TO ('2021-02-24  00:00:00');
CREATE TABLE _stage_events.events_20210224 PARTITION OF _stage_events.events FOR VALUES FROM ('2021-02-24  00:00:00') TO ('2021-02-25  00:00:00');
CREATE TABLE _stage_events.events_20210225 PARTITION OF _stage_events.events FOR VALUES FROM ('2021-02-25  00:00:00') TO ('2021-02-26  00:00:00');
CREATE TABLE _stage_events.events_20210226 PARTITION OF _stage_events.events FOR VALUES FROM ('2021-02-26  00:00:00') TO ('2021-02-27  00:00:00');
CREATE TABLE _stage_events.events_20210227 PARTITION OF _stage_events.events FOR VALUES FROM ('2021-02-27  00:00:00') TO ('2021-02-28  00:00:00');
CREATE TABLE _stage_events.events_20210228 PARTITION OF _stage_events.events FOR VALUES FROM ('2021-02-28  00:00:00') TO ('2021-03-01  00:00:00');
CREATE TABLE _stage_events.events_20210301 PARTITION OF _stage_events.events FOR VALUES FROM ('2021-03-01  00:00:00') TO ('2021-03-02  00:00:00');
CREATE TABLE _stage_events.events_20210302 PARTITION OF _stage_events.events FOR VALUES FROM ('2021-03-02  00:00:00') TO ('2021-03-03  00:00:00');
CREATE TABLE _stage_events.events_20210303 PARTITION OF _stage_events.events FOR VALUES FROM ('2021-03-03  00:00:00') TO ('2021-03-04  00:00:00');
CREATE TABLE _stage_events.events_20210304 PARTITION OF _stage_events.events FOR VALUES FROM ('2021-03-04  00:00:00') TO ('2021-03-05  00:00:00');
CREATE TABLE _stage_events.events_20210305 PARTITION OF _stage_events.events FOR VALUES FROM ('2021-03-05  00:00:00') TO ('2021-03-06  00:00:00');
CREATE TABLE _stage_events.events_20210306 PARTITION OF _stage_events.events FOR VALUES FROM ('2021-03-06  00:00:00') TO ('2021-03-07  00:00:00');
CREATE TABLE _stage_events.events_20210307 PARTITION OF _stage_events.events FOR VALUES FROM ('2021-03-07  00:00:00') TO ('2021-03-08  00:00:00');
CREATE TABLE _stage_events.events_20210308 PARTITION OF _stage_events.events FOR VALUES FROM ('2021-03-08  00:00:00') TO ('2021-03-09  00:00:00');
CREATE TABLE _stage_events.events_20210309 PARTITION OF _stage_events.events FOR VALUES FROM ('2021-03-09  00:00:00') TO ('2021-03-10  00:00:00');
CREATE TABLE _stage_events.events_20210310 PARTITION OF _stage_events.events FOR VALUES FROM ('2021-03-10  00:00:00') TO ('2021-03-11  00:00:00');
CREATE TABLE _stage_events.events_20210311 PARTITION OF _stage_events.events FOR VALUES FROM ('2021-03-11  00:00:00') TO ('2021-03-12  00:00:00');
CREATE TABLE _stage_events.events_20210312 PARTITION OF _stage_events.events FOR VALUES FROM ('2021-03-12  00:00:00') TO ('2021-03-13  00:00:00');
CREATE TABLE _stage_events.events_20210313 PARTITION OF _stage_events.events FOR VALUES FROM ('2021-03-13  00:00:00') TO ('2021-03-14  00:00:00');
CREATE TABLE _stage_events.events_20210314 PARTITION OF _stage_events.events FOR VALUES FROM ('2021-03-14  00:00:00') TO ('2021-03-15  00:00:00');
CREATE TABLE _stage_events.events_20210315 PARTITION OF _stage_events.events FOR VALUES FROM ('2021-03-15  00:00:00') TO ('2021-03-16  00:00:00');
CREATE TABLE _stage_events.events_20210316 PARTITION OF _stage_events.events FOR VALUES FROM ('2021-03-16  00:00:00') TO ('2021-03-17  00:00:00');
CREATE TABLE _stage_events.events_20210317 PARTITION OF _stage_events.events FOR VALUES FROM ('2021-03-17  00:00:00') TO ('2021-03-18  00:00:00');
CREATE TABLE _stage_events.events_20210318 PARTITION OF _stage_events.events FOR VALUES FROM ('2021-03-18  00:00:00') TO ('2021-03-19  00:00:00');
CREATE TABLE _stage_events.events_20210319 PARTITION OF _stage_events.events FOR VALUES FROM ('2021-03-19  00:00:00') TO ('2021-03-20  00:00:00');
CREATE TABLE _stage_events.events_20210320 PARTITION OF _stage_events.events FOR VALUES FROM ('2021-03-20  00:00:00') TO ('2021-03-21  00:00:00');
CREATE TABLE _stage_events.events_20210321 PARTITION OF _stage_events.events FOR VALUES FROM ('2021-03-21  00:00:00') TO ('2021-03-22  00:00:00');
CREATE TABLE _stage_events.events_20210322 PARTITION OF _stage_events.events FOR VALUES FROM ('2021-03-22  00:00:00') TO ('2021-03-23  00:00:00');

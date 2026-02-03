CREATE DATABASE beat;

\c beat;


CREATE SEQUENCE public.crawl_dump_id_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 9223372036854775807
	START 1
	CACHE 1
	NO CYCLE;



CREATE SEQUENCE public.postlog_id_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 9223372036854775807
	START 1
	CACHE 1
	NO CYCLE;





CREATE SEQUENCE public.profile_log_id_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 9223372036854775807
	START 1
	CACHE 1
	NO CYCLE;



CREATE SEQUENCE public.social_post_id_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 9223372036854775807
	START 1
	CACHE 1
	NO CYCLE;




CREATE SEQUENCE public.social_profile_id_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 9223372036854775807
	START 1
	CACHE 1
	NO CYCLE;


CREATE TABLE instagram_account (
    id bigint NOT NULL DEFAULT nextval('social_profile_id_seq'::regclass),
    profile_id character varying UNIQUE,
    handle text,
    full_name text,
    profile_type text,
    biography text,
    followers bigint,
    following bigint,
    profile_pic_url text,
    category text,
    category_id text,
    fbid text,
    media_count bigint,
    is_business_or_creator boolean,
    is_creator boolean,
    is_in_gcc_collection boolean,
    is_private boolean,
    is_verified boolean,
    on_gcc boolean,
    created_at timestamp without time zone,
    refreshed_at timestamp without time zone,
    updated_at timestamp without time zone,
    avg_comments double precision,
    avg_engagement double precision,
    avg_impressions double precision,
    avg_likes double precision,
    avg_reach double precision,
    avg_reels_plays double precision
);


CREATE TABLE credential (
    id integer GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    idempotency_key text UNIQUE,
    source text,
    credentials jsonb,
    handle text,
    created_at timestamp without time zone,
    updated_at timestamp without time zone,
    disabled_till timestamp without time zone,
    enabled boolean,
    data_access_expired boolean
);

-- Indices -------------------------------------------------------

CREATE UNIQUE INDEX credential_pkey ON credential(id int4_ops);





-- Indices -------------------------------------------------------

CREATE UNIQUE INDEX instagram_account_profile_id_key ON instagram_account(profile_id text_ops);


CREATE TABLE instagram_account_ts (
    id bigint NOT NULL DEFAULT nextval('social_profile_id_seq'::regclass),
    profile_id character varying,
    handle text,
    full_name text,
    profile_type text,
    biography text,
    followers bigint,
    following bigint,
    profile_pic_url text,
    category text,
    category_id text,
    fbid text,
    is_business_or_creator boolean,
    is_private boolean,
    is_verified boolean,
    media_count bigint,
    created_at timestamp without time zone
);

CREATE TABLE instagram_post (
    id bigint NOT NULL DEFAULT nextval('social_post_id_seq'::regclass),
    profile_id character varying,
    handle text,
    post_id text,
    post_type text,   
    caption text,
    likes bigint,
    views bigint,
    plays bigint,
    comments bigint,
    reach bigint,
    saved bigint,
    shares bigint,
    impressions bigint,
    interactions bigint,
    engagement double precision,
    shortcode character varying UNIQUE,
    display_url text,
    thumbnail_url text,
    updated_at timestamp without time zone,
    created_at timestamp without time zone,  
    publish_time timestamp without time zone
);

-- Indices -------------------------------------------------------

CREATE UNIQUE INDEX instagram_post_shortcode_key ON instagram_post(shortcode text_ops);

CREATE TABLE instagram_post_ts (
    id bigint NOT NULL DEFAULT nextval('social_post_id_seq'::regclass),
    profile_id character varying,
    handle text,
    post_id text,
    post_type text,   
    caption text,
    likes bigint,
    views bigint,
    plays bigint,
    comments bigint,
    reach bigint,
    saved bigint,
    shares bigint,
    impressions bigint,
    interactions bigint,
    engagement double precision,
    shortcode character varying,
    display_url text,
    thumbnail_url text,
    created_at timestamp without time zone,  
    publish_time timestamp without time zone
);

CREATE TABLE post_log (
    id bigint DEFAULT nextval('postlog_id_seq'::regclass) PRIMARY KEY,
    platform text,
    platform_post_id text,
    metrics jsonb,
    dimensions jsonb,
    profile_id text,
    source text,
    timestamp timestamp without time zone
);

-- Indices -------------------------------------------------------

CREATE UNIQUE INDEX postlog_pkey ON post_log(id int8_ops);

CREATE TABLE profile_log (
    id BIGSERIAL PRIMARY KEY,
    profile_id text,
    platform text,
    dimensions jsonb,
    metrics jsonb,
    source text,
    timestamp timestamp without time zone
);

-- Indices -------------------------------------------------------

CREATE UNIQUE INDEX profile_log_pkey ON profile_log(id int8_ops);


CREATE TABLE scrape_request_log (
    id bigint DEFAULT nextval('crawl_dump_id_seq'::regclass) PRIMARY KEY,
    source text,
    platform text,
    data text,
    scraped_at timestamp without time zone,
    flow text,
    idempotency_key text UNIQUE,
    status text,
    priority integer,
    params jsonb,
    created_at timestamp without time zone,
    expires_at timestamp without time zone
);

-- Indices -------------------------------------------------------

CREATE UNIQUE INDEX crawl_dump_pkey ON scrape_request_log(id int8_ops);
CREATE UNIQUE INDEX scrape_request_log_idempotency_key_key ON scrape_request_log(idempotency_key text_ops);

CREATE TABLE public.instagram_profile_insights (
    id BIGSERIAL PRIMARY KEY,
    profile_id character varying UNIQUE,
    handle text,
    created_at timestamp without time zone,
    updated_at timestamp without time zone,
    audience_country jsonb,
    audience_city jsonb,
    audience_gender_age jsonb,
    audience_locale jsonb
);

ALTER TABLE instagram_profile_insights ADD COLUMN reached_audience_country jsonb;
ALTER TABLE instagram_profile_insights ADD COLUMN reached_audience_city jsonb;
ALTER TABLE instagram_profile_insights ADD COLUMN reached_audience_gender_age jsonb;
ALTER TABLE instagram_profile_insights ADD COLUMN engaged_audience_country jsonb;
ALTER TABLE instagram_profile_insights ADD COLUMN engaged_audience_city jsonb;
ALTER TABLE instagram_profile_insights ADD COLUMN engaged_audience_gender_age jsonb;

CREATE TABLE public.youtube_account
(
    id                         bigserial NOT NULL,
    channel_id                 text      NULL,
    description                text      NULL,
    created_at                 timestamp NULL,
    updated_at                 timestamp NULL,
    custom_url                 text      NULL,
    title                      text      NULL,
    published_at               timestamp NULL,
    thumbnails                 text      NULL,
    thumbnail                  text      NULL,
    is_subscriber_count_hidden bool      NULL,
    subscribers                int8      NULL,
    uploads                    int8      NULL,
    "views"                    int8      NULL,
    shorts                     int8      NULL,
    username                   text      NULL,
    gplsuid                    text      NULL,
    country                    text      NULL,
    yt_category_id             text      NULL,
    yt_category                text      NULL,
    locale                     text      NULL,
    er12p                      float8    NULL,
    uploads_playlist_id        varchar   NULL,
    CONSTRAINT youtube_account_pkey PRIMARY KEY (id)
);

CREATE TABLE public.youtube_profile_insights
(
    id                          bigserial PRIMARY KEY,
    channel_id                  character varying UNIQUE,
    created_at                  timestamp without time zone,
    updated_at                  timestamp without time zone,
    country                     jsonb,
    city                        jsonb,
    gender_age                  jsonb
);

CREATE TABLE public.asset_log
(
    id bigserial        PRIMARY KEY ,
    entity_id           character varying UNIQUE,
    entity_type         text,
    platform            text,
    asset_type          text,
    asset_url           text,
    original_url        text,
    asset_last_updated_at TIMESTAMP without time zone
)
ALTER TABLE beat_replica.asset_log ADD COLUMN created_at TIMESTAMP WITHOUT TIME ZONE;
ALTER TABLE beat_replica.asset_log RENAME COLUMN asset_last_updated_at TO updated_at;


CREATE TABLE youtube_post (
    id bigserial PRIMARY KEY,
    description text,
    created_at timestamp without time zone,
    updated_at timestamp without time zone,
    published_at timestamp without time zone,
    channel_id text,
    title text,
    views bigint,
    thumbnails text,
    thumbnail text,
    shortcode text UNIQUE,
    likes bigint,
    favourites bigint,
    tags text,
    channel_title text,
    duration text,
    is_licensed boolean,
    category_id text,
    category text,
    audio_language text,
    content_rating text,
    playlist_id text,
    privacy_status text
);

-- Indices -------------------------------------------------------

CREATE UNIQUE INDEX youtube_post_pkey ON youtube_post(id int4_ops);

CREATE TABLE youtube_account_ts (
    id bigserial PRIMARY KEY,
    description text,
    created_at timestamp without time zone,
    published_at timestamp without time zone,
    channel_id text,
    title text,
    custom_url text,
    subscribers bigint,
    uploads bigint,
    views bigint,
    country text,
    is_subscriber_count_hidden boolean,
    uploads_playlist_id text,
    thumbnails text,
    thumbnail text
);

-- Indices -------------------------------------------------------

CREATE UNIQUE INDEX youtube_account_ts_pkey ON youtube_account_ts(id int4_ops);

CREATE TABLE youtube_post_ts (
    id bigserial PRIMARY KEY,
    description text,
    created_at timestamp without time zone,
    published_at timestamp without time zone,
    channel_id text,
    title text,
    views bigint,
    thumbnails text,
    thumbnail text,
    shortcode text,
    likes bigint,
    favourites bigint,
    tags text,
    channel_title text,
    duration text,
    is_licensed boolean,
    category_id text,
    audio_language text,
    content_rating text,
    playlist_id text,
    privacy_status text
);

CREATE TABLE "order" (
    id bigint GENERATED ALWAYS AS IDENTITY UNIQUE,
    status text,
    name text,
    landing_site text,
    utm_params text,
    created_at timestamp without time zone,
    updated_at timestamp without time zone,
    store text,
    platform_order_id text UNIQUE,
    order_data jsonb,
    platform text,
    store_name text
);

-- Indices -------------------------------------------------------

CREATE UNIQUE INDEX order_id_key ON "order"(id int8_ops);
CREATE UNIQUE INDEX order_platform_order_id_key ON "order"(platform_order_id text_ops);

-- Indices -------------------------------------------------------

CREATE UNIQUE INDEX youtube_post_ts_pkey ON youtube_post_ts(id int4_ops);
ALTER TABLE asset_log ADD COLUMN created_at TIMESTAMP WITHOUT TIME ZONE;
ALTER TABLE asset_log RENAME COLUMN asset_last_updated_at TO updated_at;

ALTER TABLE instagram_post
    ADD COLUMN content_type text,
    ADD COLUMN taps_forward bigint,
    ADD COLUMN taps_back bigint,
    ADD COLUMN exits bigint,
    ADD COLUMN replies bigint,
    ADD COLUMN navigation bigint,
    ADD COLUMN profile_activity bigint,
    ADD COLUMN profile_visits bigint,
    ADD COLUMN follows bigint,
    ADD COLUMN video_views bigint;


ALTER TABLE instagram_post_ts
    ADD COLUMN content_type text,
    ADD COLUMN taps_forward bigint,
    ADD COLUMN taps_back bigint,
    ADD COLUMN exits bigint,
    ADD COLUMN replies bigint,
    ADD COLUMN navigation bigint,
    ADD COLUMN profile_activity bigint,
    ADD COLUMN profile_visits bigint,
    ADD COLUMN follows bigint,
    ADD COLUMN video_views bigint;


ALTER TABLE instagram_post
    ADD COLUMN automatic_forward bigint,
    ADD COLUMN swipe_back bigint,
    ADD COLUMN swipe_down bigint,
    ADD COLUMN swipe_forward bigint,
    ADD COLUMN swipe_up bigint


-- ALTER TABLE instagram_post_ts
--     ADD COLUMN automatic_forward bigint,
--     ADD COLUMN swipe_back bigint,
--     ADD COLUMN swipe_down bigint,
--     ADD COLUMN swipe_forward bigint,
--     ADD COLUMN swipe_up bigint


ALTER TABLE instagram_post ADD COLUMN keywords text;


ALTER TABLE youtube_post ADD COLUMN keywords text;

create table public.sentiment_log
(
    id               bigserial
        primary key,
    platform         text,
    platform_post_id text,
    comment          text,
    sentiment        text,
    score            text,
    dimensions       text,
    metrics          text,
    source           text,
    timestamp        timestamp without time zone,
    comment_id       text
);


create table public.post_activity_log
(
    id               bigserial
        primary key,
    activity_type    text,
    platform         text,
    platform_post_id text,
    actor_profile_id text,
    dimensions       text,
    metrics          text,
    source           text,
    timestamp        timestamp without time zone
);
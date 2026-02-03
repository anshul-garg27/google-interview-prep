-- public.instagram_account definition

-- Drop table

-- DROP TABLE public.instagram_account;
CREATE TABLE public.instagram_account (
                                        id bigserial NOT NULL,
                                        "name" varchar(255) NULL,
                                        handle varchar(255) NOT NULL,
                                        business_id varchar(255) NULL,
                                        ig_id varchar(255) NOT NULL,
                                        thumbnail text NULL,
                                        dob date NULL,
                                        bio text NULL,
                                        "location" _text NULL,
                                        categories _text NULL,
                                        "label" varchar(255) NULL,
                                        languages _text NULL,
                                        gender varchar(255) NULL,
                                        city varchar(255) NULL,
                                        state varchar(255) NULL,
                                        country varchar(255) NULL,
                                        keywords _text NULL,
                                        search_phrase text NULL,
                                        is_blacklisted bool NULL,
                                        profile_id varchar(255) NULL,
                                        followers int8 NULL,
                                        "following" int8 NULL,
                                        post_count int8 NULL,
                                        ffratio float8 NULL,
                                        avg_views float8 NULL,
                                        avg_likes float8 NULL,
                                        avg_reach float8 NULL,
                                        avg_reels_play_count float8 NULL,
                                        engagement_rate float8 NULL,
                                        followers_growth7d float8 NULL,
                                        reels_reach int8 NULL,
                                        story_reach int8 NULL,
                                        image_reach int8 NULL,
                                        avg_comments float8 NULL,
                                        er_grade varchar(255) NULL,
                                        profile_admin_details jsonb NULL,
                                        profile_user_details jsonb NULL,
                                        linked_socials jsonb NULL,
                                        linked_channel_id varchar(255) NULL,
                                        audience_gender jsonb NULL,
                                        audience_age jsonb NULL,
                                        audience_location jsonb NULL,
                                        est_post_price jsonb NULL,
                                        est_reach jsonb NULL,
                                        est_impressions jsonb NULL,
                                        flag_contact_info_available bool NULL,
                                        is_verified bool DEFAULT FALSE,
                                        post_frequency_week float8 NULL,
                                        country_rank int8 NULL,
                                        category_rank jsonb NULL,
                                        authentic_engagement int8 NULL,
                                        comment_rate_percentage float8 NULL,
                                        likes_spread_percentage float8 NULL,
                                        group_key varchar(255) NULL,
                                        phone varchar(255) NULL,
                                        email varchar(255) NULL,
                                        avg_comments_grade varchar(255) NULL,
                                        avg_likes_grade varchar(255) NULL,
                                        comments_rate_grade varchar(255) NULL,
                                        followers_grade varchar(255) NULL,
                                        engagement_rate_grade varchar(255) NULL,
                                        reels_reach_grade varchar(255) NULL,
                                        likes_to_comment_ratio_grade varchar(255) NULL,
                                        followers_growth7d_grade varchar(255) NULL,
                                        followers_growth30d_grade varchar(255) NULL,
                                        followers_growth90d_grade varchar(255) NULL,
                                        audience_reachability_grade varchar(255) NULL,
                                        audience_authencity_grade varchar(255) NULL,
                                        audience_quality_grade varchar(255) NULL,
                                        post_count_grade varchar(255) NULL,
                                        likes_spread_grade varchar(255) NULL,
                                        followers_growth30d varchar(255) NULL,
                                        similar_profile_group_data jsonb NULL,
                                        image_reach_grade varchar(255) NULL,
                                        story_reach_grade varchar(255) NULL,
                                        avg_reels_play_30d int8 NULL,
                                        profile_type varchar(255) NULL,
                                        keywords_admin _text NULL,
                                        search_phrase_admin text NULL,
                                        enabled_for_saas bool DEFAULT FALSE,
                                        deleted bool DEFAULT FALSE,
                                        gcc_linked_channel_id varchar(255) NULL,
                                        gcc_profile_id varchar(255) NULL,
                                        plays_30d int8 NULL,
                                        uploads_30d int8 NULL,
                                        location_list _text NULL,
                                        impressions_30d float8 NULL,
                                        enabled bool DEFAULT TRUE,
                                        created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                        updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                        CONSTRAINT ig_id_idx UNIQUE (ig_id),
                                        CONSTRAINT instagram_account_pkey PRIMARY KEY (id)
);
CREATE INDEX CONCURRENTLY idx_youtube_acount_linked_handle ON youtube_account(linked_instagram_handle);

-- public.youtube_account definition

-- Drop table

-- DROP TABLE public.youtube_account;

CREATE TABLE public.youtube_account (
                                        id bigserial NOT NULL,
                                        channel_id varchar(255) NOT NULL,
                                        title varchar(255) NULL,
                                        username varchar(255) NULL,
                                        gender varchar(255) NULL,
                                        profile_id varchar(255) NULL,
                                        description text NULL,
                                        dob date NULL,
                                        "location" _text NULL,
                                        categories _text NULL,
                                        "label" varchar(255) NULL,
                                        languages _text NULL,
                                        is_blacklisted bool NULL,
                                        city varchar(255) NULL,
                                        state varchar(255) NULL,
                                        country varchar(255) NULL,
                                        keywords _text NULL,
                                        search_phrase text NULL,
                                        uploads_count int8 NULL,
                                        followers int8 NULL,
                                        views_count int8 NULL,
                                        avg_views  float8 NULL,
                                        video_reach int8 NULL,
                                        shorts_reach int8 NULL,
                                        thumbnail text NULL,
                                        profile_admin_details jsonb NUll,
                                        profile_user_details jsonb NULL,
                                        linked_socials jsonb NULL,
                                        linked_instagram_handle varchar(255) NULL,
                                        audience_gender jsonb NULL,
                                        audience_age jsonb NULL,
                                        audience_location jsonb NULL,
                                        est_post_price jsonb NULL,
                                        est_reach jsonb NULL,
                                        est_impressions jsonb NULL,
                                        flag_contact_info_available bool NULL,
                                        followers_growth7d float8 NULL,
                                        country_rank int8 NULL,
                                        category_rank jsonb NULL,
                                        authentic_engagement int8 NULL,
                                        comment_rate_percentage float8 NULL,
                                        video_views_last30 int8 NULL,
                                        likes_spread_percentage float8 NULL,
                                        group_key varchar(255) NULL,
                                        phone varchar(255) NULL,
                                        email varchar(255) NULL,
                                        reaction_rate int8 NULL,
                                        comments_rate_grade varchar(255) NULL,
                                        followers_grade varchar(255) NULL,
                                        comments_rate int8 NULL,
                                        followers_growth7d_grade varchar(255) NULL,
                                        cpm int8 NULL,
                                        avg_posts_per_week int8 NULL,
                                        avg_posts_per_week_grade varchar(255) NULL,
                                        followers_growth1y int8 NULL,
                                        followers_growth1y_grade varchar(255) NULL,
                                        avg_shorts_views_30d int8 NULL,
                                        avg_video_views_30d int8 NULL,
                                        latest_video_publish_time timestamp NULL,
                                        similar_profile_group_data jsonb NULL,
                                        views_30d_grade varchar(255) NULL,
                                        followers_growth30d int8 NULL,
                                        followers_growth30d_grade varchar(255) NULL,
                                        profile_type varchar(255) NULL,
                                        keywords_admin _text NULL,
                                        search_phrase_admin text NULL,
                                        enabled_for_saas bool DEFAULT FALSE,
                                        deleted bool DEFAULT FALSE,
                                        gcc_linked_instagram_handle varchar(255) NULL,
                                        gcc_profile_id varchar(255) NULL,
                                        views_30d int8 NULL,
                                        uploads_30d int8 NULL,
                                        location_list _text NULL,
                                        enabled bool DEFAULT TRUE,
                                        created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                        updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                        CONSTRAINT youtube_account_channel_id_key UNIQUE (channel_id),
                                        CONSTRAINT youtube_account_pkey PRIMARY KEY (id)
);
CREATE INDEX CONCURRENTLY idx_youtube_acount_gcc_linked_handle ON youtube_account (gcc_linked_instagram_handle);


-- public.profile_admin_details definition

-- Drop table

-- DROP TABLE public.profile_admin_details;

CREATE TABLE public.profile_admin_details (
                                            id bigserial NOT NULL,
                                            "name" varchar(255) NULL,
                                            email varchar(255) NULL,
                                            gender varchar(1) NULL,
                                            categories _text NULL,
                                            "label" varchar(255) NULL,
                                            languages _text NULL,
                                            country varchar(255) NULL,
                                            state varchar(255) NULL,
                                            city varchar(255) NULL,
                                            dob date NULL,
                                            phone varchar(20) NULL,
                                            whatsapp_enabled bool NULL,
                                            created_by varchar(255) NULL,
                                            platform varchar(40) NOT NULL,
                                            platform_profile_id int8 NOT NULL,
                                            is_blacklisted bool NULL,
                                            blacklisted_by varchar NULL,
                                            blacklisting_reason varchar NULL,
                                            "enabled" bool DEFAULT TRUE,
                                            created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                            updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                            CONSTRAINT profile_admin_details_pkey PRIMARY KEY (id)

);


-- public.profile_user_details definition

-- Drop table

-- DROP TABLE public.profile_user_details;

CREATE TABLE public.profile_user_details (
                                            id bigserial NOT NULL,
                                            "name" varchar(255) NULL,
                                            categories _text NULL,
                                            languages _text NULL,
                                            gender varchar(1) NULL,
                                            dob date NULL,
                                            "location" varchar(255) NULL,
                                            created_at timestamp NULL,
                                            updated_at timestamp NULL,
                                            created_by varchar(255) NULL,
                                            platform varchar(40) NOT NULL,
                                            platform_profile_id int8 NOT NULL,
                                            "enabled" bool DEFAULT TRUE,
                                            created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                            updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                            CONSTRAINT profile_user_details_pkey PRIMARY KEY (id)
);


-- public.activity definition

-- Drop table

-- DROP TABLE public.activity;

CREATE TABLE public.activity (
                                id bigserial NOT NULL,
                                partner_id int8 NOT NULL,
                                title text NOT NULL,
                                activity_type varchar(255) NOT NULL,
                                params jsonb NULL,
                                filters jsonb NULL,
                                result_count int8 NULL,
                                platform varchar(255) NOT NULL,
                                "enabled" bool DEFAULT TRUE,
                                created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                CONSTRAINT activity_pkey PRIMARY KEY (id)
);

-- public.leaderboard definition

-- Drop table

-- DROP TABLE public.leaderboard;

CREATE TABLE public.leaderboard (
                                id bigserial NOT NULL,
                                month date NULL,
                                language varchar(255) NULL,
                                category varchar(255) NULL,
                                platform varchar(255) NULL,
                                platform_profile_id int8 NULL,
                                handle varchar(255) NULL,
                                followers int8 NULL,
                                engagement_rate float8 NULL,
                                avg_likes float8 NULL,
                                avg_comments float8 NULL,
                                avg_views float8 NULL,
                                avg_plays float8 NULL,
                                followers_rank int8 NULL,
                                followers_change_rank int8 NULL,
                                views_rank int8 NULL,
                                followers_rank_by_cat int8 NULL,
                                views_rank_by_cat int8 NULL,
                                followers_rank_by_lang int8 NULL,
                                views_rank_by_lang int8 NULL,
                                followers_rank_by_cat_lang int8 NULL,
                                views_rank_by_cat_lang int8 NULL,
                                platform_id varchar NULL,
                                country varchar(50) NULL,
                                views int8 NULL,
                                prev_followers int8 NULL,
                                followers_change int8 NULL,
                                likes int8 NULL,
                                comments int8 NULL,
                                plays int8 NULL,
                                engagement int8 NULL,
                                followers_change_rank_by_cat_lang int4 NULL,
                                plays_rank int8 NULL,
                                plays_rank_by_cat int8 NULL,
                                plays_rank_by_lang int8 NULL,
                                plays_rank_by_cat_lang int8 NULL,
                                followers_change_rank_by_cat int8 NULL,
                                followers_change_rank_by_lang int8 NULL,
                                profile_platform varchar(50) NULL,
                                yt_views int8 NULL,
                                ia_views int8 NULL,
                                last_month_ranks jsonb NULL,
                                profile_type varchar(255) NULL,
                                enabled bool DEFAULT TRUE,
                                created_at timestamp NULL DEFAULT CURRENT_TIMESTAMP,
                                updated_at timestamp NULL DEFAULT CURRENT_TIMESTAMP,
                                CONSTRAINT leaderboard_pkey PRIMARY KEY (id)
);
CREATE INDEX leaderboard_tmp_month_platform_platform_profile_id_idx ON public.leaderboard USING btree (month, platform, platform_profile_id);
CREATE INDEX leaderboard_tmp_month_platform_followers_change_rank_id_idx ON public.leaderboard USING btree (month, platform, followers_change_rank, id);
CREATE INDEX leaderboard_tmp_month_platform_plays_rank_id_idx ON public.leaderboard USING btree (month, platform, plays_rank, id);
CREATE INDEX leaderboard_tmp_month_platform_views_rank_id_idx ON public.leaderboard USING btree (month, platform, views_rank, id);
CREATE INDEX idx_language_month_platform_followers_change_rank_by_lang_id ON public.leaderboard USING btree (language, month, platform,followers_change_rank_by_lang, id);
CREATE INDEX idx_category_month_platform_followers_change_rank_by_cat_id ON public.leaderboard USING btree (category, month, platform, followers_change_rank_by_cat, id);
CREATE INDEX idx_category_language_month_platform_followers_change_rank_by_c ON public.leaderboard USING btree (category, language, month, platform, followers_change_rank_by_cat_lang, id);
CREATE INDEX idx_language_month_platform_plays_rank_by_lang_id ON public.leaderboard USING btree (language, month, platform, plays_rank_by_lang, id);
CREATE INDEX idx_category_month_platform_plays_rank_by_cat_id ON public.leaderboard USING btree (category, month, platform, plays_rank_by_cat, id);
CREATE INDEX idx_cat_lang_month_platform_plays_rank_by_cat_lang_id ON public.leaderboard USING btree (category, language, month, platform, plays_rank_by_cat_lang, id);
CREATE INDEX idx_language_month_platform_views_rank_by_lang_id ON public.leaderboard USING btree (language, month, platform, views_rank_by_lang, id);
CREATE INDEX idx_cat_lang_month_platform_views_rank_by_cat_lang_id ON public.leaderboard USING btree (category, language, month, platform, views_rank_by_cat_lang, id);
-- public.social_profile_time_series definition

-- Drop table

-- DROP TABLE public.social_profile_time_series;
CREATE TABLE public.social_profile_time_series (
                                    id bigserial NOT NULL,
                                    platform_profile_id int8 NULL,
                                    platform varchar(255) NOT NULL,
                                    followers int8 NULL,
                                    following int8 NULL,
                                    views int8 NULL,
                                    plays int8 NULL,
                                    uploads int8 NULL,
                                    engagement_rate float8 NULL,
                                    "date" DATE NOT NULL,
                                    platform_id varchar(255) NOT NULL,
                                    followers_change int8 NULL,
                                    views_total int8 NULL,
                                    plays_total int8 NULL,
                                    enabled bool DEFAULT TRUE,
                                    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                    CONSTRAINT social_profile_time_series_pkey PRIMARY KEY (id)
);
CREATE index idx_platform_platform_profile_id_date ON social_profile_time_series (platform,platform_profile_id,date);


-- public.social_profile_hashtags definition

-- Drop table

-- DROP TABLE public.social_profile_hashtags;


CREATE TABLE public.social_profile_hashtags (
                                        id bigserial NOT NULL,
                                        platform_profile_id int8 NULL,
                                        platform varchar(255) NOT NULL,
                                        hashtags _text NOT NULL,
                                        hashtags_counts _text NOT NULL,
                                        platform_id varchar(255) NOT NULL,
                                        "enabled" bool DEFAULT TRUE,
                                        created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                        updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                        CONSTRAINT social_profile_hashtags_pkey PRIMARY KEY (id)
);

-- public.locations definition

-- Drop table

-- DROP TABLE public.locations;

CREATE TABLE public.locations(
                                id bigserial NOT NULL,
                                name varchar(255) NOT NULL,
                                full_name varchar(255) NOT NULL,
                                type varchar(255) NOT NULL,
                                "enabled" bool DEFAULT TRUE,
                                created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                CONSTRAINT locations_pkey PRIMARY KEY (id)
);


-- public.social_profile_audience_info definition

-- Drop table

-- DROP TABLE public.social_profile_audience_info;

CREATE TABLE public.social_profile_audience_info(
                                id bigserial NOT NULL,
                                platform_profile_id int8 NULL,
                                platform varchar(255) NOT NULL,
                                audience_location_split jsonb NULL,
                                audience_gender_split jsonb NULL,
                                audience_age_gender_split jsonb NULL,
                                audience_age_split jsonb NULL,
                                audience_language jsonb NULL,
                                notable_followers _text NULL,
                                audience_reachability_percentage float8 NULL,
                                audience_authenticity_percentage float8 NULL,
                                quality_audience_percentage float8 NULL,
                                quality_audience_score float8 NULL,
                                quality_score_grade varchar(255) NULL,
                                quality_score_breakup jsonb NULL,
                                platform_id varchar(255) NULL,
                                "enabled" bool DEFAULT TRUE,
                                created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                CONSTRAINT social_profile_audience_info_pkey PRIMARY KEY (id)
);
CREATE index idx_platform_id ON social_profile_audience_info (platform_id);

-- public.group_metrics definition

-- Drop table

-- DROP TABLE public.group_metrics;

CREATE TABLE public.group_metrics (
                                id int8 NOT NULL DEFAULT nextval('group_metrics_v2_id_seq'::regclass),
                                group_key varchar(50) NULL,
                                profiles int4 NULL,
                                bin_start _text NULL,
                                bin_end _text NULL,
                                bin_height _text NULL,
                                enabled bool NULL DEFAULT true,
                                created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                group_avg_likes float4 NULL,
                                group_avg_comments float4 NULL,
                                group_avg_comments_rate float4 NULL,
                                group_avg_followers float4 NULL,
                                group_avg_engagement_rate float4 NULL,
                                group_avg_reels_reach varchar(50) NULL,
                                group_avg_likes_to_comment_ratio float4 NULL,
                                group_avg_followers_growth7d float4 NULL,
                                group_avg_followers_growth30d float4 NULL,
                                group_avg_followers_growth90d float4 NULL,
                                group_avg_audience_reachability float4 NULL,
                                group_avg_audience_authencity float4 NULL,
                                group_avg_audience_quality float4 NULL,
                                group_avg_post_count float4 NULL,
                                group_avg_reaction_rate int4 NULL,
                                group_avg_likes_spread float4 NULL,
                                group_avg_followers_growth1y float4 NULL,
                                group_avg_posts_per_week float4 NULL,
                                group_avg_image_reach float4 NULL,
                                group_avg_story_reach float4 NULL,
                                group_avg_video_reach float4 NULL,
                                group_avg_shorts_reach float4 NULL,
                                CONSTRAINT group_metrics_v2_pk PRIMARY KEY (id)
);

-- public.profile_collection definition

-- Drop table

-- DROP TABLE public.profile_collection;

CREATE TABLE public.profile_collection (
                                id bigserial NOT NULL,
                                share_id varchar NULL,
                                partner_id int8 NOT NULL,
                                "name" varchar(255) NOT NULL,
                                description text NULL,
                                categories _text NULL,
                                tags _text NULL,
                                enabled bool NOT NULL DEFAULT true,
                                featured bool NOT NULL DEFAULT false,
                                created_by varchar NOT NULL,
                                created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                analytics_enabled bool NOT NULL DEFAULT false,
                                disabled_metrics jsonb NOT NULL DEFAULT '[]'::jsonb,
                                "source" varchar(45) NULL,
                                source_id varchar(255) NULL,
                                campaign_id int8 NULL,
                                campaign_platform varchar(45) NULL,
                                job_id int8 NULL,
                                category_ids jsonb NULL,
                                category jsonb NULL,
                                custom_columns jsonb NOT NULL default '[]'::jsonb,
                                ordered_columns jsonb NOT NULL default '[]'::jsonb,
                                share_enabled boolean NOT NULL default true,
                                CONSTRAINT profile_collection_pkey PRIMARY KEY (id),
                                CONSTRAINT profile_collection_share_id_key UNIQUE (share_id)
);
CREATE INDEX profile_collection_name_key ON public.profile_collection USING btree (name);
CREATE INDEX profile_collection_partnerid_key ON public.profile_collection USING btree (partner_id);
CREATE UNIQUE INDEX profile_collection_source_key ON public.profile_collection USING btree (source, source_id);
alter table profile_collection ADD CONSTRAINT profile_collection_campaign_idx UNIQUE NULLS NOT DISTINCT (campaign_id);
-- older versions
CREATE UNIQUE INDEX profile_collection_campaign_idx ON profile_collection (campaign_id) WHERE campaign_id IS NULL;
-- public.profile_collection_item definition

-- Drop table

-- DROP TABLE public.profile_collection_item;

CREATE TABLE public.profile_collection_item (
                                id bigserial NOT NULL,
                                platform varchar(40) NOT NULL,
                                platform_account_code varchar(255) NOT NULL,
                                profile_collection_id int8 NOT NULL,
                                hidden bool NOT NULL DEFAULT false,
                                enabled bool NULL DEFAULT true,
                                rank int8 NOT NULL,
                                created_by varchar NOT NULL,
                                created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                partner_id int8 NULL,
                                profile_social_id varchar(255) NULL,
                                campaign_profile_id int8 NULL,
                                shortlisting_status varchar(45) NULL,
                                shortlist_id varchar(255) NULL,
                                brand_selection_status varchar(45) NULL,
                                brand_selection_remarks text NULL,
                                custom_columns_data jsonb NOT NULL DEFAULT '[]'::jsonb,
                                internal_commercials varchar(127) DEFAULT NULL,
                                CONSTRAINT profile_collection_item_pkey PRIMARY KEY (id),
                                CONSTRAINT profile_collection_item_profile_collection_id_platform_acco_key UNIQUE (profile_collection_id, platform_account_code, platform)
);
CREATE INDEX platform_account_code_idx ON public.profile_collection_item USING btree (platform_account_code);
CREATE INDEX platform_idx ON public.profile_collection_item USING btree (platform);
CREATE INDEX profile_collection_idx ON public.profile_collection_item USING btree (profile_collection_id);
CREATE INDEX profile_social_idx ON public.profile_collection_item USING btree (profile_social_id);
CREATE INDEX profile_collection_item_shortlist_idx ON public.profile_collection_item USING btree (shortlist_id);
-- Post Collection Tables

CREATE TABLE public.post_collection (
                                           id varchar(255) NOT NULL,
                                           share_id varchar NULL UNIQUE,
                                           partner_id int8 NOT NULL,
                                           name varchar(255) NOT NULL,
                                           source varchar(45) NOT NULL,
                                           source_id varchar(255) NOT NULL,
                                           budget int8 NULL,
                                           enabled bool NOT NULL DEFAULT true,
                                           start_time timestamp NULL,
                                           end_time timestamp NULL,
                                           created_by varchar NOT NULL,
                                           created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                           updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                           disabled_metrics jsonb NULL,
                                           metrics_ingestion_freq varchar(45) NOT NULL default 'default',
                                           job_id int8 NULL,
                                           CONSTRAINT post_collection_pkey PRIMARY KEY (id)
);

CREATE INDEX post_collection_name_key ON public.post_collection (name);
CREATE INDEX post_collection_partnerid_key ON public.post_collection (partner_id);
CREATE UNIQUE INDEX post_collection_source_key ON public.post_collection (source, source_id);
ALTER TABLE post_collection ADD CONSTRAINT post_collection_source_key UNIQUE NULLS NOT DISTINCT (source, source_id)


CREATE TABLE public.post_collection_item (
                                            id bigserial NOT NULL,
                                            platform varchar(40) NOT NULL,
                                            short_code text NOT NULL,
                                            post_type varchar(45) NOT NULL,
                                            post_collection_id varchar(255) NOT NULL,
                                            bookmarked bool NOT NULL DEFAULT FALSE,
                                            enabled bool NOT NULL DEFAULT TRUE,
                                            cost int8 NULL,
                                            posted_by_handle varchar(255) NULL,
                                            sponsor_links jsonb NULL,
                                            post_title text NULL,
                                            post_link varchar(255) NULL,
                                            post_thumbnail text NULL,
                                            metrics_ingestion_freq varchar(45) NOT NULL default 'default',
                                            retrieve_data boolean default false,
                                            show_in_report boolean default false,
                                            posted_by_cp_id int8 default NULL,
                                            created_by varchar NOT NULL,
                                            created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                            updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,

                                            CONSTRAINT post_collection_item_pkey PRIMARY KEY (id)
);

CREATE UNIQUE INDEX post_collection_item_u_key ON public.post_collection_item (post_collection_id, platform, short_code);
CREATE INDEX post_collection_item_pcid_key ON public.post_collection_item (post_collection_id);
CREATE INDEX post_collection_item_scode_key ON public.post_collection_item (short_code);


-- Collection Analytics Metrics Summary AND TimeSeries Tables

CREATE TABLE public.collection_keyword (
                                            collection_id varchar(255) NOT NULL,
                                            collection_type varchar(45) NOT NULL,
                                            collection_share_id varchar(255) NOT NULL,
                                            updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                            keyword_name varchar(255) NOT NULL,
                                            tagged_count int8 NOT NULL DEFAULT 0
);

CREATE UNIQUE INDEX collection_kw_u_key ON public.collection_keyword (collection_type, collection_id);
CREATE UNIQUE INDEX collection_kw_u_pp_key ON public.collection_keyword (collection_type, collection_share_id);

CREATE TABLE public.collection_hashtag (
                                           collection_id varchar(255) NOT NULL,
                                           collection_type varchar(45) NOT NULL,
                                           collection_share_id varchar(255) NOT NULL,
                                           updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                           hashtag_name varchar(255) NOT NULL,
                                           tagged_count int8 NOT NULL DEFAULT 0,
                                           ugc_tagged_count int8 NOT NULL DEFAULT 0
);

CREATE UNIQUE INDEX collection_ht_u_key ON public.collection_hashtag (collection_type, collection_id);
CREATE UNIQUE INDEX collection_ht_u_pp_key ON public.collection_hashtag (collection_type, collection_share_id);

CREATE TABLE public.collection_post_metrics_summary (
                                                    id bigserial NOT NULL,
                                                    platform varchar(45) NOT NULL,
                                                    post_short_code varchar(255) NOT NULL,
                                                    post_type varchar(45) NOT NULL,
                                                    post_title text NOT NULL,
                                                    post_link varchar(255) NOT NULL,
                                                    post_thumbnail text NULL,
                                                    hashtags jsonb NULL,
                                                    published_at date NOT NULL,

                                                    platform_account_code varchar(45) NULL,
                                                    profile_social_id varchar(255) NULL,
                                                    profile_handle varchar(255) NOT NULL,
                                                    profile_name varchar(255) NOT NULL,
                                                    profile_pic text NULL,
                                                    followers int8 NOT NULL DEFAULT 0,
                                                    cost int8 NULL,

                                                    collection_id varchar(255) NOT NULL,
                                                    collection_type varchar(45) NOT NULL,
                                                    collection_share_id varchar(255) NOT NULL,
                                                    post_collection_item_id int8 NULL,

                                                    views int8 NOT NULL DEFAULT 0,
                                                    likes int8 NOT NULL DEFAULT 0,
                                                    comments int8 NOT NULL DEFAULT 0,
                                                    impressions int8 NOT NULL DEFAULT 0,
                                                    saves int8 NOT NULL DEFAULT 0,
                                                    plays int8 NOT NULL DEFAULT 0,
                                                    reach int8 NOT NULL DEFAULT 0,
                                                    swipe_ups int8 NOT NULL DEFAULT 0,
                                                    mentions int8 NOT NULL DEFAULT 0,
                                                    sticker_taps int8 NOT NULL DEFAULT 0,
                                                    shares int8 NOT NULL DEFAULT 0,
                                                    story_exits int8 NOT NULL DEFAULT 0,
                                                    story_back_taps int8 NOT NULL DEFAULT 0,
                                                    story_forward_taps int8 NOT NULL DEFAULT 0,
                                                    link_clicks int8 NOT NULL DEFAULT 0,
                                                    orders int8 NOT NULL DEFAULT 0,
                                                    delivered_orders int8 NOT NULL DEFAULT 0,
                                                    completed_orders int8 NOT NULL DEFAULT 0,
                                                    leaderboard_overall_orders int8 NOT NULL DEFAULT 0,
                                                    leaderboard_delivered_orders int8 NOT NULL DEFAULT 0,
                                                    leaderboard_completed_orders int8 NOT NULL DEFAULT 0,
                                                    er float8 NOT NULL DEFAULT 0.0,

                                                    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                                    CONSTRAINT collection_post_metrics_summary_pkey PRIMARY KEY (id)
);

CREATE UNIQUE INDEX collection_post_metrics_summary_u_key ON public.collection_post_metrics_summary (collection_id, collection_type, platform, post_short_code);
CREATE INDEX collection_post_metrics_summary_scode_key ON public.collection_post_metrics_summary (post_short_code);
CREATE INDEX collection_post_metrics_summary_cid_key ON public.collection_post_metrics_summary (collection_id);
CREATE INDEX collection_post_metrics_summary_ctype_key ON public.collection_post_metrics_summary (collection_type);
CREATE INDEX collection_post_metrics_summary_csid_key ON public.collection_post_metrics_summary (collection_share_id);


CREATE TABLE public.collection_post_metrics_ts (
                                                  id bigserial NOT NULL,
                                                  platform varchar(45) NOT NULL,
                                                  post_short_code varchar(255) NOT NULL,
                                                  post_type varchar(45) NOT NULL,
                                                  post_title text NOT NULL,
                                                  post_link varchar(255) NOT NULL,
                                                  post_thumbnail text NULL,
                                                  published_at date NOT NULL,

                                                  platform_account_code varchar(45) NOT NULL,
                                                  profile_social_id varchar(255) NOT NULL,
                                                  profile_handle varchar(255) NOT NULL,
                                                  profile_name varchar(255) NOT NULL,
                                                  profile_pic text NULL,
                                                  followers int8 NOT NULL DEFAULT 0,
                                                  cost int8 NOT NULL DEFAULT 0,

                                                  collection_id varchar(255) NOT NULL,
                                                  collection_type varchar(45) NOT NULL,
                                                  collection_share_id varchar(255) NOT NULL,
                                                  post_collection_item_id int8 NULL,

                                                  stats_date date NOT NULL,
                                                  insight_source varchar(45) NOT NULL DEFAULT 'pipeline',
                                                  views int8 NOT NULL DEFAULT 0,
                                                  likes int8 NOT NULL DEFAULT 0,
                                                  comments int8 NOT NULL DEFAULT 0,
                                                  impressions int8 NOT NULL DEFAULT 0,
                                                  saves int8 NOT NULL DEFAULT 0,
                                                  plays int8 NOT NULL DEFAULT 0,
                                                  reach int8 NOT NULL DEFAULT 0,
                                                  swipe_ups int8 NOT NULL DEFAULT 0,
                                                  mentions int8 NOT NULL DEFAULT 0,
                                                  sticker_taps int8 NOT NULL DEFAULT 0,
                                                  shares int8 NOT NULL DEFAULT 0,
                                                  story_exits int8 NOT NULL DEFAULT 0,
                                                  story_back_taps int8 NOT NULL DEFAULT 0,
                                                  story_forward_taps int8 NOT NULL DEFAULT 0,
                                                  link_clicks int8 NOT NULL DEFAULT 0,
                                                  orders int8 NOT NULL DEFAULT 0,
                                                  er float8 NOT NULL DEFAULT 0.0,

                                                  updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                                  CONSTRAINT collection_post_metrics_ts_pkey PRIMARY KEY (id)
);


CREATE UNIQUE INDEX collection_post_metrics_post_u_key ON public.collection_post_metrics_ts (collection_id, collection_type, platform, post_short_code, stats_date);
CREATE INDEX collection_post_metrics_post_scode_key ON public.collection_post_metrics_ts (post_short_code);
CREATE INDEX collection_post_metrics_ctype_key ON public.collection_post_metrics_ts (collection_type);
CREATE INDEX collection_post_metrics_cid_key ON public.collection_post_metrics_ts (collection_id);
CREATE INDEX collection_post_metrics_csid_key ON public.collection_post_metrics_ts (collection_share_id);


CREATE TABLE public.genre_overview (
                                        id bigserial NOT NULL,
                                        category varchar(255) NOT NULL,
                                        month varchar(255) NOT NULL,
                                        platform varchar(255) NOT NULL,
                                        language varchar(255) NOT NULL,
                                        profile_type varchar(255) NOT NULL,
                                        creators int8 NULL,
                                        uploads int8 NULL,
                                        views int8 NULL,
                                        followers int8 NULL,
                                        likes int8 NULL,
                                        comments int8 NULL,
                                        audience_age_gender jsonB NULL,
                                        audience_age jsonB NULL,
                                        audience_gender jsonB NULL,
                                        engagement int8 NULL,
                                        engagement_rate int8 NULL,
                                        enabled bool DEFAULT TRUE,
                                        created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                        updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                        CONSTRAINT genre_overview_pkey PRIMARY KEY (id),
                                        UNIQUE (category, month, platform,language,profile_type)
);

CREATE TABLE public.trending_content (
                                        id bigserial NOT NULL,
                                        category varchar(255) NULL,
                                        platform varchar(255) NOT NULL,
                                        language varchar(255) NULL,
                                        profile_type varchar(255) NOT NULL,
                                        views int8 NULL,
                                        plays int8 NULL,
                                        likes int8 NULL,
                                        comments int8 NULL,
                                        published_at timestamp,
                                        post_id varchar(255) NULL,
                                        post_type varchar(255) NULL,
                                        post_link varchar(255) NULL,
                                        thumbnail text NULL,
                                        platform_account_id int8 NULL,
                                        profile_social_id varchar(255) NULL,
                                        audience_age_gender jsonB NULL,
                                        month varchar(255) NULL,
                                        plays_rank int8 NULL,
                                        views_rank int8 NULL,
                                        title text NULL,
                                        description text NULL,
                                        enabled bool DEFAULT TRUE,
                                        created_at timestamp NULL DEFAULT CURRENT_TIMESTAMP,
                                        updated_at timestamp NULL DEFAULT CURRENT_TIMESTAMP,
                                        CONSTRAINT trending_content_pkey PRIMARY KEY (id)
);


CREATE TABLE public.campaign_profiles (
                                        id bigserial NOT NULL,
                                        platform character varying(255) NOT NULL,
                                        platform_account_id bigint NOT NULL,
                                        admin_details jsonb,
                                        user_details jsonb,
                                        on_gcc boolean,
                                        on_gcc_app boolean,
                                        gcc_user_account_id bigint,
                                        updated_by character varying(255),
                                        enabled boolean DEFAULT true,
                                        created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
                                        updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);
CREATE UNIQUE INDEX idx_cp_id ON public.campaign_profiles USING btree (id);
CREATE UNIQUE INDEX cp_platfom_id ON public.campaign_profiles (platform, platform_account_id);

CREATE TABLE public.collection_group (
                                         id bigserial NOT NULL,
                                         share_id varchar NULL,
                                         partner_id int8 NOT NULL,
                                         name varchar(255) NOT NULL,
                                         "source" varchar(45) NOT NULL,
                                         source_id varchar(255) NOT NULL,
                                         collection_ids _text NULL,
                                         enabled bool NOT NULL DEFAULT true,
                                         disabled_metrics jsonb NULL,
                                         created_by varchar NOT NULL,
                                         created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                         updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                         metadata jsonb NULL,
                                         objective varchar(40) NULL,
                                         CONSTRAINT collection_group_pkey PRIMARY KEY (id),
                                         CONSTRAINT collection_group_share_id_key UNIQUE (share_id)
);
CREATE INDEX collection_group_partnerid_key ON public.collection_group USING btree (partner_id);
CREATE UNIQUE INDEX collection_group_source_key ON public.collection_group USING btree (source, source_id);

-- Permissions

ALTER TABLE public.collection_group OWNER TO gccuser;
GRANT ALL ON TABLE public.collection_group TO gccuser;

CREATE TABLE public.profile_collection_item_cc_data
(
    id                         bigserial    NOT NULL,
    profile_collection_item_id int8         NOT NULL,
    key                      varchar(127) NOT NULL,
    value                      varchar(255)  NOT NULL,

    CONSTRAINT pKey PRIMARY KEY (id),
    CONSTRAINT pci_custom_column_uniq UNIQUE (profile_collection_item_id, key)
);

CREATE TABLE public.partner_usage (
                                    id bigserial NOT NULL,
                                    partner_id int8 NOT NULL,
                                    namespace varchar(255) NOT NULL,
                                    key varchar(255) NOT NULL,
                                    "limit" int8 NOT NULL,
                                    consumed int8 NULL,
                                    start_date timestamp NOT NULL,
                                    end_date timestamp NOT NULL,
                                    next_reset_on timestamp NOT NULL,
                                    frequency varchar(255) NOT NULL,
                                    enabled bool NULL DEFAULT true,
                                    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                    CONSTRAINT partner_usage_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_partner_usage_key ON public.partner_usage USING btree (partner_id, key);
CREATE INDEX idx_partner_usage_key_module ON public.partner_usage USING btree (partner_id, key, namespace);

CREATE TABLE public.partner_profile_page_track(
                                        id bigserial NOT NULL,
                                        partner_id int8 NOT NULL,
                                        key varchar(255) NOT NULL,
                                        enabled bool NULL DEFAULT true,
                                        created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                        updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                        CONSTRAINT partner_profile_page_track_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_partner_profile_key ON public.partner_profile_page_track USING btree (partner_id, key);

CREATE TABLE winkl_collection_migration_info
(
    winkl_collection_id         int8    NOT NULL,
    winkl_collection_share_id   varchar(255) NOT NULL,
    profile_collection_id       int8 NOT NULL,
    winkl_campaign_id           int8 DEFAULT NULL,
    winkl_shortlist_id          int8 DEFAULT NULL,
    profile_collection_share_id varchar(255) NOT NULL,

    CONSTRAINT winkl_collection_id_uniq UNIQUE (winkl_collection_id),
    CONSTRAINT winkl_collection_share_id_uniq UNIQUE (winkl_collection_share_id),
    CONSTRAINT saas_collection_shortlist_id_uniq UNIQUE NULLS NOT DISTINCT (profile_collection_id, winkl_shortlist_id)
);

-- public.activity_tracker definition

-- Drop table

-- DROP TABLE public.activity_tracker;

CREATE TABLE public.activity_tracker (
	id bigserial NOT NULL,
	"type" varchar(255) NOT NULL,
	meta jsonb NULL,
	enabled bool NULL,
	created_at timestamp NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at timestamp NULL DEFAULT CURRENT_TIMESTAMP,
	account_id int8 NOT NULL,
	partner_id int8 NOT NULL,
	CONSTRAINT activity_tracker_pkey PRIMARY KEY (id)
);

CREATE UNIQUE INDEX activity_tracker_pkey ON public.activity_tracker USING btree (id);
# ULTRA-DEEP ANALYSIS: STIR DATA PLATFORM

## PROJECT OVERVIEW

| Attribute | Value |
|-----------|-------|
| **Project Name** | Stir |
| **Purpose** | Enterprise Data Platform for Social Media Analytics - Influencer Discovery, Leaderboards, Collections |
| **Architecture** | Modern Data Stack (ELT) with Apache Airflow + dbt + ClickHouse |
| **Git Commits** | 1,476 (mature production project) |
| **Total Lines of Code** | ~17,500+ |
| **Total DAGs** | 76 |
| **Total dbt Models** | 112 (29 staging + 83 marts) |

---

## 1. COMPLETE DIRECTORY STRUCTURE

```
/stir/
├── .git/                              # Git repository (1,476 commits)
├── .gitignore                         # Git ignore rules
├── .gitlab-ci.yml                     # CI/CD pipeline configuration
│
├── dags/                              # Airflow DAGs (76 files)
│   ├── __pycache__/                   # Python cache
│   │
│   │ # dbt Orchestration DAGs (11)
│   ├── dbt_core.py                    # Core models (*/15 min)
│   ├── dbt_hourly.py                  # Hourly transforms (every 2h)
│   ├── dbt_daily.py                   # Daily batch (19:00 UTC)
│   ├── dbt_weekly.py                  # Weekly aggregates
│   ├── dbt_collections.py             # Collection processing (*/30 min)
│   ├── dbt_staging_collections.py     # Staging collections (*/30 min)
│   ├── dbt_recent_scl.py              # Recent scrape logging (*/5 min)
│   ├── dbt_gcc_orders.py              # Order processing (daily)
│   ├── dbt_refresh_account_tracker_stats.py  # Account stats (hourly)
│   ├── post_ranker.py                 # Post ranking (*/5 hours)
│   ├── post_ranker_partial.py         # Partial ranking (*/15 min)
│   │
│   │ # Instagram Sync DAGs (17)
│   ├── sync_insta_collection_posts.py
│   ├── sync_insta_collection_stories.py
│   ├── sync_insta_post_comments.py
│   ├── sync_insta_post_insights.py
│   ├── sync_insta_profile_followers.py
│   ├── sync_insta_profile_following.py
│   ├── sync_insta_profile_insights.py
│   ├── sync_insta_profiles_by_handle.py
│   ├── sync_insta_stories.py
│   ├── sync_insta_stories_explicitly.py
│   ├── sync_insta_story_insights.py
│   ├── sync_instagram_gpt_data_audience_age_gender copy.py
│   ├── sync_instagram_gpt_data_audience_cities.py
│   ├── sync_instagram_gpt_data_base_categ_lang_topics.py
│   ├── sync_instagram_gpt_data_base_gender.py
│   ├── sync_instagram_gpt_data_base_location.py
│   ├── sync_instagram_gpt_data_gender_location_lang.py
│   │
│   │ # YouTube Sync DAGs (12)
│   ├── sync_yt_channels.py
│   ├── sync_yt_collection_posts.py
│   ├── sync_yt_critical_daily.py      # 127 hardcoded channels
│   ├── sync_yt_genre_videos.py
│   ├── sync_yt_post_comments.py
│   ├── sync_yt_post_type.py
│   ├── sync_yt_profile_insights.py
│   ├── sync_yt_profile_relationship_by_channel_id.py
│   ├── sync_yt_profiles_by_handle.py
│   ├── sync_yt_profiles_videos_by_handle.py
│   ├── sync_vidooly_es_youtube_channels.py
│   ├── retry_yt_scrape_events.py
│   │
│   │ # Collection & Leaderboard Sync DAGs (15)
│   ├── sync_collection_post_summary_prod.py
│   ├── sync_collection_post_summary_staging.py
│   ├── sync_collection_post_ts_prod.py
│   ├── sync_collection_post_ts_staging.py
│   ├── sync_collection_hashtags_prod.py
│   ├── sync_collection_keywords_prod.py
│   ├── sync_leaderboard_prod.py
│   ├── sync_leaderboard_staging.py
│   ├── sync_time_series_prod.py
│   ├── sync_time_series_staging.py
│   ├── sync_trending_content_prod.py
│   ├── sync_trending_content_staging.py
│   ├── sync_genre_overview_prod.py
│   ├── sync_genre_overview_staging.py
│   ├── sync_hashtags_prod.py
│   │
│   │ # Operational & Verification DAGs (9)
│   ├── sync_g3_collection_posts.py
│   ├── sync_group_metrics_prod.py
│   ├── sync_final_post_submitted_for_cp.py
│   ├── sync_missing_journey.py
│   ├── sync_post_with_saas_again.py
│   ├── sync_post_collection_sentiment_report_path.py
│   ├── sync_keyword_collection_report.py
│   ├── sync_shopify_orders.py
│   ├── mark_campaign_completed_on_refer_launch.py
│   │
│   │ # Asset Upload DAGs (7)
│   ├── upload_post_asset.py
│   ├── upload_post_asset_stories.py
│   ├── upload_insta_profile_asset.py
│   ├── upload_profile_relationship_asset.py
│   ├── upload_content_verification.py
│   ├── upload_handle_verification.py
│   ├── uca_su_saas_item_sync.py
│   │
│   │ # Utility DAGs (4)
│   ├── crm_recommendation_invitation.py
│   ├── create_payout_for_active_su_if_not_present.py
│   ├── track_hashtags.py
│   └── slack_connection.py            # Slack notification helper
│
├── src/gcc_social/                    # dbt Project Directory
│   ├── dbt_project.yml                # dbt configuration
│   ├── dbt_packages/                  # External packages
│   ├── logs/                          # Execution logs
│   ├── target/                        # Compiled artifacts
│   │
│   └── models/                        # dbt Models (112 total)
│       │
│       ├── staging/                   # Staging Layer (29 models)
│       │   │
│       │   ├── beat/                  # Beat Source (13 models)
│       │   │   ├── stg_beat_asset_log.sql
│       │   │   ├── stg_beat_credential.sql
│       │   │   ├── stg_beat_instagram_account.sql
│       │   │   ├── stg_beat_instagram_post.sql
│       │   │   ├── stg_beat_instagram_profile_insights.sql
│       │   │   ├── stg_beat_order.sql
│       │   │   ├── stg_beat_post_activity_log.sql
│       │   │   ├── stg_beat_post_log.sql
│       │   │   ├── stg_beat_profile_log.sql
│       │   │   ├── stg_beat_profile_relationship_log.sql
│       │   │   ├── stg_beat_recent_scrape_request_log.sql
│       │   │   ├── stg_beat_scrape_request_log.sql
│       │   │   ├── stg_beat_youtube_account.sql
│       │   │   ├── stg_beat_youtube_post.sql
│       │   │   └── stg_beat_youtube_profile_insights.sql
│       │   │
│       │   └── coffee/                # Coffee Source (16 models)
│       │       ├── stg_coffee_activity_tracker.sql
│       │       ├── stg_coffee_campaign_profiles.sql
│       │       ├── stg_coffee_collection_group.sql
│       │       ├── stg_coffee_keyword_collection.sql
│       │       ├── stg_coffee_post_collection.sql
│       │       ├── stg_coffee_post_collection_item.sql
│       │       ├── stg_coffee_profile_collection.sql
│       │       ├── stg_coffee_profile_collection_item.sql
│       │       ├── stg_coffee_stage_view_instagram_account_lite.sql
│       │       ├── stg_coffee_stage_view_youtube_account_lite.sql
│       │       ├── stg_coffee_staging_post_collection_item.sql
│       │       ├── stg_coffee_staging_profile_collection_item.sql
│       │       └── stg_coffee_view_*_account_lite.sql
│       │
│       └── marts/                     # Mart Layer (83 models)
│           │
│           ├── audience/              # Audience Analytics (4)
│           │   ├── mart_audience_info.sql
│           │   ├── mart_audience_info_follower.sql
│           │   ├── mart_audience_info_gpt.sql
│           │   └── mart_audience_info_private.sql
│           │
│           ├── collection/            # Collection Management (13)
│           │   ├── mart_collection_clicks_ts.sql
│           │   ├── mart_collection_post.sql
│           │   ├── mart_collection_post_clicks.sql
│           │   ├── mart_collection_post_ts.sql
│           │   ├── mart_collection_social_ts.sql
│           │   ├── mart_fake_events.sql
│           │   ├── mart_post_collection_*_post_ts.sql
│           │   └── mart_profile_collection_*_post_ts.sql
│           │
│           ├── discovery/             # Discovery & Analytics (16)
│           │   ├── mart_influencer_perf.sql
│           │   ├── mart_instagram_account.sql
│           │   ├── mart_instagram_gpt_basic_data.sql
│           │   ├── mart_instagram_hashtags.sql
│           │   ├── mart_instagram_phone.sql
│           │   ├── mart_instagram_tracked_profiles.sql
│           │   ├── mart_insta_predicted_*.sql
│           │   ├── mart_linked_socials.sql
│           │   ├── mart_manual_data_*.sql
│           │   ├── mart_primary_group_metrics.sql
│           │   ├── mart_youtube_account.sql
│           │   ├── mart_youtube_account_language.sql
│           │   ├── mart_youtube_profile_relationship.sql
│           │   └── mart_youtube_tracked_profiles.sql
│           │
│           ├── genre/                 # Genre Analysis (7)
│           │   ├── mart_genre_instagram_trending_content_*.sql
│           │   ├── mart_genre_overview.sql
│           │   ├── mart_genre_overview_all.sql
│           │   ├── mart_genre_trending_content.sql
│           │   └── mart_genre_youtube_trending_content_*.sql
│           │
│           ├── leaderboard/           # Rankings (14)
│           │   ├── mart_cross_platform_leaderboard.sql
│           │   ├── mart_insta_account_monthly.sql
│           │   ├── mart_insta_account_weekly.sql
│           │   ├── mart_insta_leaderboard_base.sql
│           │   ├── mart_instagram_leaderboard.sql
│           │   ├── mart_leaderboard.sql
│           │   ├── mart_time_series.sql
│           │   ├── mart_time_series_with_gaps.sql
│           │   ├── mart_youtube_leaderboard.sql
│           │   ├── mart_yt_account_*.sql
│           │   ├── mart_yt_growth.sql
│           │   └── mart_yt_leaderboard_base*.sql
│           │
│           ├── orders/                # Order Processing (1)
│           │   └── mart_gcc_orders.sql
│           │
│           ├── posts/                 # Post Analytics (3)
│           │   ├── mart_instagram_post_hashtags.sql
│           │   ├── mart_instagram_post_tagged.sql
│           │   └── mart_post_location.sql
│           │
│           ├── profile_stats_full/    # Full Profile Stats (8)
│           │   ├── mart_instagram_account_summary.sql
│           │   ├── mart_instagram_account_summary_via_handles.sql
│           │   ├── mart_instagram_all_posts.sql
│           │   ├── mart_instagram_all_posts_with_ranks.sql
│           │   ├── mart_instagram_creators_followers_fake_analysis.sql
│           │   ├── mart_instagram_post_ranks.sql
│           │   ├── mart_instagram_profiles_followers_count.sql
│           │   └── mart_instagram_recent_post_stats.sql
│           │
│           ├── profile_stats_partial/ # Partial Updates (6)
│           │   └── mart_instagram_*_partial.sql
│           │
│           └── staging_collection/    # Staging Collections (9)
│               └── mart_staging_*.sql
│
├── .dbt/                              # dbt Configuration
│   ├── profiles.yml                   # Database connections
│   ├── dbt_project.yml
│   ├── dbt_packages/
│   ├── logs/
│   └── target/
│
├── scripts/
│   └── start.sh                       # Deployment startup
│
├── requirements.txt                   # Python dependencies (350+)
├── start.sh                           # Project startup
└── README.md                          # Documentation
```

---

## 2. TECHNOLOGY STACK

### Core Technologies
| Category | Technology | Version | Purpose |
|----------|-----------|---------|---------|
| **Orchestration** | Apache Airflow | 2.6.3 | Workflow scheduling & monitoring |
| **Data Transformation** | dbt-core | 1.3.1 | SQL-based transformations |
| **Analytics DB** | ClickHouse | via dbt-clickhouse 1.3.2 | OLAP queries, real-time analytics |
| **Transactional DB** | PostgreSQL | 5.4.0 | Operational data storage |
| **Cloud Storage** | AWS S3 | gcc-social-data bucket | Data staging & backups |
| **CI/CD** | GitLab CI | -- | Deployment pipeline |

### Python Dependencies (350+)
```
# Airflow & Extensions
apache-airflow==2.6.3
airflow-dbt-python==0.15.2
airflow-clickhouse-plugin==1.0.0

# Database Drivers
clickhouse-connect==0.5.12
clickhouse-driver==0.2.5
psycopg2-binary==2.9.5
SQLAlchemy==1.4.45

# dbt
dbt-core==1.3.1
dbt-clickhouse==1.3.2
dbt-postgres==1.3.1

# ML/Data Libraries
tensorflow==2.11.0
torch==2.0.1
scikit-learn==1.0.2
transformers==4.29.2
pandas==1.3.5
numpy==1.21.6

# Messaging & Monitoring
pika==1.3.2  # RabbitMQ
slack-sdk==3.21.3
loguru==0.7.0

# HTTP & Web
requests
paramiko  # SSH
```

---

## 3. APACHE AIRFLOW DAGs (76 TOTAL)

### DAG Category Breakdown

| Category | Count | Description |
|----------|-------|-------------|
| dbt Orchestration | 11 | Transformation pipelines |
| Instagram Sync | 17 | Instagram data collection |
| YouTube Sync | 12 | YouTube data collection |
| Collection/Leaderboard Sync | 15 | Analytics sync |
| Operational/Verification | 9 | Data quality & operations |
| Asset Upload | 7 | Media asset processing |
| Utility | 5 | One-off & helper DAGs |

### Scheduling Frequencies

| Frequency | DAGs | Purpose |
|-----------|------|---------|
| `*/5 * * * *` | 2 | Real-time: dbt_recent_scl, post_ranker |
| `*/10 * * * *` | 5 | Near real-time: Profile lookups |
| `*/15 * * * *` | 2 | Core: dbt_core, post_ranker_partial |
| `*/30 * * * *` | 2 | Collections: dbt_collections, staging |
| `0 * * * *` | 12 | Hourly: Most sync operations |
| `0 */3 * * *` | 8 | Every 3 hours: Heavy syncs |
| `0 0 * * *` | 15 | Daily midnight: Full refreshes |
| `0 19 * * *` | 1 | Daily 19:00: dbt_daily |
| `15 20 * * *` | 2 | Daily 20:15: Leaderboards |
| `0 6 */7 * *` | 1 | Weekly: dbt_weekly |

### Operator Distribution

| Operator | Count | Usage |
|----------|-------|-------|
| **PythonOperator** | 46 | Data fetching, API calls, processing |
| **PostgresOperator** | 20 | Data loading, table operations |
| **ClickHouseOperator** | 19 | Export queries, analytics |
| **SSHOperator** | 18 | File transfer, remote commands |
| **DbtRunOperator** | 11 | dbt model execution |

### Connection IDs Used
```python
connections = {
    "clickhouse_gcc": "Primary ClickHouse cluster",
    "prod_pg": "Production PostgreSQL",
    "stage_pg": "Staging PostgreSQL",
    "ssh_prod_pg": "SSH to production server",
    "ssh_stage_pg": "SSH to staging server",
    "beat": "Beat API service",
    "slack_failure_conn": "Slack failure notifications",
    "slack_success_conn": "Slack success notifications"
}
```

---

## 4. dbt MODELS (112 TOTAL)

### Model Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                     dbt MODEL LAYERS                             │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  SOURCES                                                         │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│  │ beat_replica │  │   vidooly    │  │    coffee    │          │
│  │ Instagram/YT │  │ Cross-plat   │  │ Campaigns    │          │
│  └──────────────┘  └──────────────┘  └──────────────┘          │
│         ↓                 ↓                 ↓                    │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │              STAGING LAYER (29 models)                    │  │
│  │  stg_beat_*  (13)  │  stg_coffee_*  (16)                 │  │
│  │  - Raw data extraction with minimal transformation        │  │
│  │  - Type casting, NULL handling                            │  │
│  └──────────────────────────────────────────────────────────┘  │
│         ↓                                                        │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │               MART LAYER (83 models)                      │  │
│  │                                                           │  │
│  │  ┌─────────────────┐  ┌─────────────────┐               │  │
│  │  │ Audience (4)    │  │ Collection (13) │               │  │
│  │  │ mart_audience_* │  │ mart_collection*│               │  │
│  │  └─────────────────┘  └─────────────────┘               │  │
│  │                                                           │  │
│  │  ┌─────────────────┐  ┌─────────────────┐               │  │
│  │  │ Discovery (16)  │  │ Genre (7)       │               │  │
│  │  │ mart_instagram_ │  │ mart_genre_*    │               │  │
│  │  │ mart_youtube_*  │  │                 │               │  │
│  │  └─────────────────┘  └─────────────────┘               │  │
│  │                                                           │  │
│  │  ┌─────────────────┐  ┌─────────────────┐               │  │
│  │  │ Leaderboard(14) │  │ Profile Stats   │               │  │
│  │  │ mart_leaderboard│  │ Full (8)        │               │  │
│  │  │ mart_time_series│  │ Partial (6)     │               │  │
│  │  └─────────────────┘  └─────────────────┘               │  │
│  │                                                           │  │
│  │  ┌─────────────────┐  ┌─────────────────┐               │  │
│  │  │ Posts (3)       │  │ Orders (1)      │               │  │
│  │  │ Staging Col (9) │  │ mart_gcc_orders │               │  │
│  │  └─────────────────┘  └─────────────────┘               │  │
│  └──────────────────────────────────────────────────────────┘  │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### dbt Tags Classification

| Tag | Count | Schedule | Purpose |
|-----|-------|----------|---------|
| **core** | 11 | */15 min | Core business metrics |
| **deprecated** | 12 | Excluded | Legacy/unused models |
| **post_ranker** | 10 | */5 hours | Post ranking algorithms |
| **collections** | 2 | */30 min | Collection management |
| **daily** | 1 | 19:00 UTC | Daily batch processing |
| **hourly** | 1 | Every 2h | Hourly updates |
| **weekly** | 1 | Weekly | Weekly aggregates |
| **staging_collections** | 1 | */30 min | Staging area |
| **gcc_orders** | 1 | Daily | Order processing |
| **account_tracker_stats** | 1 | Hourly | Account tracking |

### Materialization Strategies

```sql
-- Table (most common) - Full refresh
{{ config(materialized='table') }}

-- Incremental - For time-series data
{{ config(
    materialized='incremental',
    engine='ReplacingMergeTree()',
    order_by='(profile_id, date)',
    incremental_strategy='append'
) }}

-- View - For simple transformations (rare)
{{ config(materialized='view') }}
```

### Staging Models (29)

**Beat Source (13 models):**
```sql
-- stg_beat_instagram_account
SELECT
    profile_id,
    handle,
    full_name,
    biography,
    followers_count,
    following_count,
    posts_count,
    is_verified,
    is_business,
    category,
    external_url,
    profile_pic_url,
    created_at,
    updated_at
FROM beat_replica.instagram_account

-- stg_beat_instagram_post
SELECT
    post_id,
    profile_id,
    short_code,
    post_type,
    caption,
    likes_count,
    comments_count,
    views_count,
    timestamp,
    location_id,
    hashtags,
    mentions
FROM beat_replica.instagram_post
```

**Coffee Source (16 models):**
```sql
-- stg_coffee_post_collection
SELECT
    collection_id,
    name,
    partner_id,
    is_active,
    show_in_report,
    created_at
FROM coffee.post_collection

-- stg_coffee_campaign_profiles
SELECT
    campaign_id,
    profile_id,
    status,
    deliverables_count,
    completed_deliverables
FROM coffee.campaign_profiles
```

### Key Mart Models (83)

**mart_instagram_account (Core Discovery Model):**
```sql
{{ config(
    materialized='table',
    tags=['core', 'hourly']
) }}

WITH base_accounts AS (
    SELECT * FROM {{ ref('stg_beat_instagram_account') }}
),

post_stats AS (
    SELECT
        profile_id,
        COUNT(*) as total_posts,
        AVG(likes_count) as avg_likes,
        AVG(comments_count) as avg_comments,
        SUM(likes_count) as total_likes,
        SUM(comments_count) as total_comments
    FROM {{ ref('stg_beat_instagram_post') }}
    WHERE timestamp > now() - INTERVAL 30 DAY
    GROUP BY profile_id
),

engagement AS (
    SELECT
        profile_id,
        (avg_likes + avg_comments) / NULLIF(followers_count, 0) * 100 as engagement_rate
    FROM base_accounts
    JOIN post_stats USING (profile_id)
)

SELECT
    a.*,
    ps.total_posts,
    ps.avg_likes,
    ps.avg_comments,
    e.engagement_rate,
    -- Rankings
    row_number() OVER (ORDER BY followers_count DESC) as followers_rank,
    row_number() OVER (PARTITION BY category ORDER BY followers_count DESC) as followers_rank_by_cat,
    row_number() OVER (PARTITION BY language ORDER BY followers_count DESC) as followers_rank_by_lang
FROM base_accounts a
LEFT JOIN post_stats ps USING (profile_id)
LEFT JOIN engagement e USING (profile_id)
```

**mart_leaderboard (Multi-dimensional Rankings):**
```sql
{{ config(
    materialized='table',
    tags=['daily']
) }}

SELECT
    profile_id,
    handle,
    platform,
    followers_count,
    engagement_rate,
    category,
    language,
    country,

    -- Global ranks
    followers_rank,
    engagement_rank,

    -- Category ranks
    followers_rank_by_cat,
    engagement_rank_by_cat,

    -- Language ranks
    followers_rank_by_lang,
    engagement_rank_by_lang,

    -- Combined ranks
    followers_rank_by_cat_lang,
    engagement_rank_by_cat_lang,

    -- Rank changes (vs last month)
    followers_rank - lag(followers_rank) OVER (
        PARTITION BY profile_id ORDER BY snapshot_date
    ) as rank_change,

    snapshot_date
FROM {{ ref('mart_instagram_account') }}
```

**mart_time_series (Incremental):**
```sql
{{ config(
    materialized='incremental',
    engine='ReplacingMergeTree()',
    order_by='(profile_id, date)',
    unique_key='(profile_id, date)'
) }}

SELECT
    profile_id,
    toDate(created_at) as date,
    argMax(followers_count, created_at) as followers,
    argMax(following_count, created_at) as following,
    argMax(posts_count, created_at) as posts,
    max(created_at) as last_updated
FROM {{ ref('stg_beat_instagram_account') }}

{% if is_incremental() %}
WHERE created_at > (SELECT max(last_updated) - INTERVAL 4 HOUR FROM {{ this }})
{% endif %}

GROUP BY profile_id, date
```

---

## 5. DATA FLOW ARCHITECTURE

### Three-Layer Data Flow Pattern
```
┌─────────────────────────────────────────────────────────────────┐
│                    DATA FLOW ARCHITECTURE                        │
└─────────────────────────────────────────────────────────────────┘

LAYER 1: CLICKHOUSE (Analytics Engine)
┌─────────────────────────────────────────────────────────────────┐
│  ClickHouse Database (172.31.28.68:9000)                        │
│  ├── dbt.stg_* (29 staging tables)                              │
│  └── dbt.mart_* (83 mart tables)                                │
│                                                                  │
│  Features:                                                       │
│  - ReplacingMergeTree for upserts                               │
│  - Partitioning by date for query optimization                  │
│  - ArrayMap for hashtag normalization                           │
│  - argMax for latest value extraction                           │
└─────────────────────────────────────────────────────────────────┘
                            ↓
                 (ClickHouseOperator)
                 INSERT INTO FUNCTION s3(...)
                            ↓
LAYER 2: AWS S3 (Staging)
┌─────────────────────────────────────────────────────────────────┐
│  S3 Bucket: gcc-social-data                                     │
│  Path: /data-pipeline/tmp/*.json                                │
│  Format: JSONEachRow                                            │
│                                                                  │
│  Settings:                                                       │
│  - s3_truncate_on_insert=1                                      │
│  - Automatic compression                                        │
└─────────────────────────────────────────────────────────────────┘
                            ↓
                 (SSHOperator)
                 aws s3 cp s3://... /tmp/
                            ↓
LAYER 3: POSTGRESQL (Operational)
┌─────────────────────────────────────────────────────────────────┐
│  PostgreSQL Database (172.31.2.21:5432)                         │
│  Database: beat                                                 │
│                                                                  │
│  Tables:                                                        │
│  - instagram_account                                            │
│  - youtube_account                                              │
│  - collection_post_metrics_summary                              │
│  - leaderboard_*                                                │
│                                                                  │
│  Operations:                                                    │
│  - COPY from /tmp/*.json                                        │
│  - JSONB parsing with type casting                              │
│  - Atomic table swap (RENAME)                                   │
└─────────────────────────────────────────────────────────────────┘
```

### Detailed Data Flow Example (dbt_collections DAG)

```python
# Step 1: dbt transformation in ClickHouse
dbt_run_task = DbtRunOperator(
    task_id='dbt_run_collections',
    models='tag:collections',
    profiles_dir='/Users/.dbt',
    target='gcc_warehouse'  # ClickHouse
)

# Step 2: Export to S3
clickhouse_export = ClickHouseOperator(
    task_id='export_to_s3',
    clickhouse_conn_id='clickhouse_gcc',
    sql="""
        INSERT INTO FUNCTION s3(
            's3://gcc-social-data/data-pipeline/tmp/collection_post.json',
            'AWS_KEY', 'AWS_SECRET',
            'JSONEachRow'
        )
        SELECT * FROM dbt.mart_collection_post
        SETTINGS s3_truncate_on_insert=1
    """
)

# Step 3: Download via SSH
ssh_download = SSHOperator(
    task_id='download_from_s3',
    ssh_conn_id='ssh_prod_pg',
    command='aws s3 cp s3://gcc-social-data/data-pipeline/tmp/collection_post.json /tmp/'
)

# Step 4: Load into PostgreSQL temp table
pg_load = PostgresOperator(
    task_id='load_temp_table',
    postgres_conn_id='prod_pg',
    sql="""
        CREATE TEMP TABLE tmp_collection_post (data JSONB);
        COPY tmp_collection_post FROM '/tmp/collection_post.json';
    """
)

# Step 5: Transform and insert
pg_transform = PostgresOperator(
    task_id='transform_insert',
    postgres_conn_id='prod_pg',
    sql="""
        INSERT INTO collection_post_metrics_summary_new
        SELECT
            (data->>'collection_id')::bigint,
            (data->>'post_short_code')::text,
            (data->>'likes_count')::bigint,
            (data->>'comments_count')::bigint,
            (data->>'engagement_rate')::float
        FROM tmp_collection_post
    """
)

# Step 6: Atomic table swap
pg_swap = PostgresOperator(
    task_id='atomic_swap',
    postgres_conn_id='prod_pg',
    sql="""
        ALTER TABLE collection_post_metrics_summary
            RENAME TO collection_post_metrics_summary_old_bkp;
        ALTER TABLE collection_post_metrics_summary_new
            RENAME TO collection_post_metrics_summary;
    """
)

# Task dependencies
dbt_run_task >> clickhouse_export >> ssh_download >> pg_load >> pg_transform >> pg_swap
```

---

## 6. CLICKHOUSE INTEGRATION

### Connection Configuration
```python
# Connection details
connection = {
    "host": "172.31.28.68",
    "port": 9000,
    "database": "dbt",
    "user": "airflow",
    "threads": 3,
    "send_receive_timeout": 3600
}
```

### Key ClickHouse Features Used

**1. ReplacingMergeTree Engine:**
```sql
-- Efficient upsert operations
ENGINE = ReplacingMergeTree(updated_at)
ORDER BY (profile_id, date)
PARTITION BY toYYYYMM(date)
```

**2. S3 Integration:**
```sql
-- Export to S3
INSERT INTO FUNCTION s3(
    's3://bucket/path/file.json',
    'KEY', 'SECRET',
    'JSONEachRow'
)
SELECT * FROM dbt.mart_table
SETTINGS s3_truncate_on_insert=1

-- Import from S3
SELECT * FROM s3(
    's3://bucket/path/file.csv',
    'KEY', 'SECRET',
    'CSV'
)
```

**3. Window Functions:**
```sql
-- Ranking
row_number() OVER (ORDER BY followers_count DESC)
row_number() OVER (PARTITION BY category ORDER BY followers_count DESC)

-- Latest value
argMax(followers_count, created_at)

-- Cardinality
uniqExact(profile_id)
```

**4. Array Operations:**
```sql
-- Hashtag extraction
arrayMap(x -> lower(trim(x)), splitByChar(',', hashtags))

-- Array aggregation
groupArray(hashtag)
```

### Query Patterns

**Incremental Processing:**
```sql
{% if is_incremental() %}
WHERE created_at > (
    SELECT max(created_at) - INTERVAL 4 HOUR
    FROM {{ this }}
)
{% endif %}
```

**Partition Pruning:**
```sql
WHERE toYYYYMM(date) >= toYYYYMM(now() - INTERVAL 30 DAY)
```

---

## 7. BUSINESS LOGIC & METRICS

### Core Business Problems Solved

**1. Influencer Discovery & Ranking:**
- Multi-dimensional leaderboards (followers, engagement, growth)
- Segmentation by category, language, country
- Monthly snapshots with rank change tracking
- Cross-platform comparisons (Instagram vs YouTube)

**2. Collection Analytics:**
- Curated collections of posts/profiles
- Aggregated engagement metrics
- Time-series analysis for trending content
- Sentiment tracking

**3. Campaign Management:**
- Order tracking from campaigns
- Content verification for deliverables
- Referral journey tracking
- Payout processing

**4. Cross-Platform Linking:**
- Match Instagram profiles with YouTube channels
- Unified audience insights
- Combined reach calculations

### Key Metrics Computed

**Profile-Level Metrics:**
```sql
-- Engagement calculations
engagement_rate = (avg_likes + avg_comments) / followers_count * 100

-- Growth metrics
followers_change = current_followers - previous_month_followers
growth_rate = followers_change / previous_month_followers * 100

-- Activity metrics
avg_posts_per_week = total_posts / weeks_active
avg_likes = total_likes / total_posts
avg_comments = total_comments / total_posts
avg_views = total_views / total_videos  -- YouTube
```

**Ranking Dimensions:**
```sql
-- Global ranks
followers_rank              -- All profiles
engagement_rank             -- By engagement rate

-- Category ranks
followers_rank_by_cat       -- Within category
engagement_rank_by_cat      -- Within category

-- Language ranks
followers_rank_by_lang      -- Within language
engagement_rank_by_lang     -- Within language

-- Combined ranks
followers_rank_by_cat_lang  -- Category + Language
engagement_rank_by_cat_lang -- Category + Language
```

---

## 8. DAG CONFIGURATIONS

### dbt DAGs Configuration

**dbt_core (Critical Path):**
```python
dag = DAG(
    dag_id='dbt_core',
    default_args={
        'owner': 'airflow',
        'retries': 1,
        'retry_delay': timedelta(minutes=5),
        'on_failure_callback': SlackNotifier.slack_fail_alert
    },
    schedule_interval='*/15 * * * *',  # Every 15 minutes
    start_date=datetime(2023, 1, 1),
    catchup=False,
    max_active_runs=1,
    concurrency=1,
    dagrun_timeout=timedelta(minutes=60)
)

dbt_run = DbtRunOperator(
    task_id='dbt_run_core',
    models='tag:core',
    profiles_dir='/Users/.dbt',
    target='gcc_warehouse',
    dag=dag
)
```

**dbt_daily (Full Refresh):**
```python
dag = DAG(
    dag_id='dbt_daily',
    schedule_interval='0 19 * * *',  # 19:00 UTC daily
    dagrun_timeout=timedelta(minutes=360),
    max_active_runs=1
)

dbt_run = DbtRunOperator(
    task_id='dbt_run_daily',
    models='tag:daily',
    full_refresh=True
)
```

### Sync DAGs Configuration

**sync_insta_collection_posts:**
```python
dag = DAG(
    dag_id='sync_insta_collection_posts',
    schedule_interval='*/10 * * * *',  # Every 10 minutes
    max_active_runs=1,
    concurrency=1
)

def create_scrape_requests(**context):
    """Creates scrape requests for collection posts"""
    connection = BaseHook.get_connection("clickhouse_gcc")
    client = clickhouse_connect.get_client(
        host=connection.host,
        password=connection.password,
        username=connection.login
    )

    # Query posts needing refresh
    sql = """
        SELECT post_id, short_code
        FROM dbt.mart_collection_post
        WHERE last_scraped < now() - INTERVAL 1 HOUR
        LIMIT 1000
    """
    result = client.query(sql)

    # Create scrape requests via Beat API
    for row in result:
        requests.post(
            f'{BEAT_URL}/scrape_request_log/flow/instagram_post',
            json={'short_code': row['short_code']}
        )

create_requests_task = PythonOperator(
    task_id='create_scrape_requests',
    python_callable=create_scrape_requests,
    dag=dag
)
```

**sync_leaderboard_prod (ClickHouse → S3 → PostgreSQL):**
```python
dag = DAG(
    dag_id='sync_leaderboard_prod',
    schedule_interval='15 20 * * *',  # Daily at 20:15 UTC
    max_active_runs=1
)

# Task 1: Export from ClickHouse to S3
export_task = ClickHouseOperator(
    task_id='export_leaderboard',
    clickhouse_conn_id='clickhouse_gcc',
    sql="""
        INSERT INTO FUNCTION s3(
            's3://gcc-social-data/data-pipeline/tmp/leaderboard.json',
            'KEY', 'SECRET', 'JSONEachRow'
        )
        SELECT * FROM dbt.mart_leaderboard
        SETTINGS s3_truncate_on_insert=1
    """
)

# Task 2: Download via SSH
download_task = SSHOperator(
    task_id='download_leaderboard',
    ssh_conn_id='ssh_prod_pg',
    command='aws s3 cp s3://gcc-social-data/data-pipeline/tmp/leaderboard.json /tmp/'
)

# Task 3: Load into PostgreSQL
load_task = PostgresOperator(
    task_id='load_leaderboard',
    postgres_conn_id='prod_pg',
    sql="""
        -- Create temp table
        CREATE TEMP TABLE tmp_leaderboard (data JSONB);

        -- Load JSON
        COPY tmp_leaderboard FROM '/tmp/leaderboard.json';

        -- Insert with transformation
        INSERT INTO leaderboard_new
        SELECT
            (data->>'profile_id')::bigint,
            (data->>'handle')::text,
            (data->>'followers_rank')::int,
            (data->>'engagement_rank')::int,
            now()
        FROM tmp_leaderboard;
    """,
    execution_timeout=timedelta(seconds=14400)
)

# Task 4: Atomic swap
swap_task = PostgresOperator(
    task_id='swap_tables',
    postgres_conn_id='prod_pg',
    sql="""
        ALTER TABLE leaderboard RENAME TO leaderboard_old;
        ALTER TABLE leaderboard_new RENAME TO leaderboard;
        DROP TABLE IF EXISTS leaderboard_old;
    """
)

export_task >> download_task >> load_task >> swap_task
```

---

## 9. CI/CD & DEPLOYMENT

### GitLab CI Configuration
```yaml
stages:
  - deploy_prod

deploy_prod:
  stage: deploy_prod
  script:
    - /bin/bash scripts/start.sh
  tags:
    - beat-deployer-1
  environment:
    name: prod
  when: manual  # Manual trigger required
  only:
    - master
```

### Deployment Script (start.sh)
```bash
#!/bin/bash

# Copy local configuration
cp airflow.cfg.local airflow.cfg

# Initialize Airflow database
airflow db init

# Start Airflow in standalone mode
airflow standalone
```

### dbt Profiles Configuration
```yaml
# .dbt/profiles.yml

gcc_social:
  target: dev
  outputs:
    dev:
      type: postgres
      host: 172.31.2.21
      port: 5432
      user: airflow
      password: "{{ env_var('PG_PASSWORD') }}"
      database: beat
      schema: public
      threads: 6

gcc_warehouse:
  target: prod
  outputs:
    prod:
      type: clickhouse
      host: 172.31.28.68
      port: 9000
      user: airflow
      password: "{{ env_var('CH_PASSWORD') }}"
      database: dbt
      schema: dbt
      threads: 3
```

---

## 10. KEY METRICS SUMMARY

| Metric | Value |
|--------|-------|
| **Total DAG Files** | 76 |
| **Lines of Python (DAGs)** | 7,406 |
| **Staging Models** | 29 |
| **Mart Models** | 83 |
| **Total dbt Models** | 112 |
| **Git Commits** | 1,476 |
| **Python Dependencies** | 350+ |
| **ClickHouse Operators** | 19 |
| **PostgreSQL Operators** | 20 |
| **SSH Operators** | 18 |
| **Python Operators** | 46 |
| **dbt Run Operators** | 11 |

---

## 11. DATA SOURCES INTEGRATED

| Source | Type | Data |
|--------|------|------|
| **Instagram API** | Social | Posts, profiles, insights, stories, comments, followers |
| **YouTube API** | Social | Videos, channels, profiles, insights, comments |
| **Beat Database** | Internal | Scrape requests, credentials, asset logs |
| **Coffee Database** | Internal | Campaigns, profiles, collections, deliverables |
| **Shopify** | E-commerce | Orders |
| **Vidooly** | External | YouTube channel data |
| **S3** | Storage | Temporary data staging |

---

## 12. SKILLS DEMONSTRATED

### Data Engineering
- **ETL/ELT Architecture**: Modern data stack with Airflow + dbt
- **Workflow Orchestration**: 76 production DAGs with complex dependencies
- **SQL Optimization**: ClickHouse-specific tuning, window functions
- **Incremental Processing**: Smart backfill with 4-hour windows
- **Data Quality**: Validation, atomic operations, backup strategies

### Analytics & Data Modeling
- **Star Schema Design**: Fact and dimension tables
- **dbt Mastery**: Staging/mart layers, incremental materializations
- **Multi-dimensional Analysis**: Rankings across categories, languages, regions
- **Time-Series Processing**: Trend detection, growth tracking

### Database Administration
- **ClickHouse**: ReplacingMergeTree, partitioning, S3 integration
- **PostgreSQL**: JSONB parsing, atomic table swaps, copy operations
- **Cross-Database Sync**: ClickHouse → S3 → PostgreSQL patterns

### DevOps & Infrastructure
- **CI/CD**: GitLab pipeline with manual production deploys
- **Cloud Integration**: AWS S3, SSH tunneling
- **Monitoring**: Slack notifications, execution timeouts
- **Configuration Management**: dbt profiles, Airflow connections

---

## 13. INTERVIEW TALKING POINTS

### 1. "Tell me about a data platform you built"
- **Scale**: 76 DAGs processing billions of records daily
- **Architecture**: Modern data stack with Airflow + dbt + ClickHouse
- **Dual Database**: ClickHouse for OLAP, PostgreSQL for OLTP
- **Outcome**: Real-time influencer discovery and campaign analytics

### 2. "Describe your experience with dbt"
- **Models**: Built 112 models (29 staging + 83 marts)
- **Incremental**: ReplacingMergeTree for efficient upserts
- **Tags**: Organized execution by core, daily, hourly, collections
- **Testing**: Data quality validation framework

### 3. "How do you handle large-scale data processing?"
- **ClickHouse**: OLAP engine for billion-record analytics
- **Partitioning**: Date-based for query optimization
- **Incremental**: 4-hour lookback windows
- **Parallelism**: 25 worker threads, concurrent DAGs

### 4. "Explain a complex data pipeline you built"
- **Flow**: ClickHouse → S3 → SSH → PostgreSQL → Atomic swap
- **Transformation**: JSONB parsing with type casting
- **Reliability**: Backup tables, atomic operations
- **Monitoring**: Slack alerts, execution timeouts

### 5. "How do you ensure data quality?"
- **Atomic Operations**: Table rename prevents partial updates
- **Validation**: Type casting, NULL handling
- **Monitoring**: Slack notifications on failures
- **Backups**: Old table preserved before swap

---

## 14. NOTABLE OBSERVATIONS

### Production Hardening
- `max_active_runs=1` prevents concurrent execution issues
- `catchup=False` avoids backfill on schedule changes
- `dagrun_timeout` prevents runaway jobs
- Slack notifications for all failures

### Performance Optimizations
- `ORDER BY` clauses for ClickHouse query efficiency
- `PARTITION BY` for date-based pruning
- `argMax` for efficient latest value extraction
- `uniqExact` for accurate cardinality

### Data Pipeline Patterns
- Three-layer flow: ClickHouse → S3 → PostgreSQL
- Atomic table swaps for zero-downtime updates
- 4-hour incremental windows for freshness vs performance

### Security Notes
- Credentials in profiles.yml (should use secrets manager)
- AWS keys in DAG code (should use IAM roles)
- Hardcoded channel lists (127 YouTube channels)

---

*Analysis covers 17,500+ lines of code across 76 DAGs, 112 dbt models, and complete data infrastructure spanning ClickHouse, PostgreSQL, S3, and multiple APIs.*

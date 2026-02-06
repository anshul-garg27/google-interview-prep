# GCC (Good Creator Co.) - Complete System Architecture

> **MASTER DOCUMENT**: How all 6 projects work together as a unified influencer marketing analytics platform.
> Built by **Anshul Garg** at Good Creator Co. -- ~54K+ lines of production code across Python and Go.

---

## TABLE OF CONTENTS

1. [Complete System Architecture Diagram](#1-complete-system-architecture-diagram)
2. [Data Flow Diagrams by Use Case](#2-data-flow-diagrams-by-use-case)
3. [Technology Stack Matrix](#3-technology-stack-matrix)
4. [Inter-Project Communication](#4-inter-project-communication)
5. [Database Architecture](#5-database-architecture)
6. [Deployment Architecture](#6-deployment-architecture)
7. [Scale Numbers](#7-scale-numbers-across-all-projects)
8. [Interview: 60-Second Combined Pitch](#8-how-to-talk-about-this-in-interviews)

---

## 1. COMPLETE SYSTEM ARCHITECTURE DIAGRAM

```
============================================================================================
                     GCC PLATFORM - COMPLETE SYSTEM ARCHITECTURE
                     Built by Anshul Garg | ~54K+ LOC | 6 Projects
============================================================================================

                                INTERNET / CLIENTS
                    (Mobile Apps, Web Apps, Brand Dashboards)
                                      |
                                      | HTTPS
                                      v
+-------------------------------------------------------------------------------------+
|                                                                                     |
|                        [1] SAAS GATEWAY  (Go, 2.2K LOC)                             |
|                            Port 8009 | Gin Framework                                |
|                                                                                     |
|   +---------------------------+  +-------------------+  +------------------------+  |
|   | 7-Layer Middleware         |  | JWT + Redis Auth  |  | Reverse Proxy          |  |
|   | Pipeline                  |  | HMAC-SHA256       |  | 13 Microservices       |  |
|   |  1. Panic Recovery        |  | Session Lookup    |  | Proxied                |  |
|   |  2. Sentry Error Track    |  | Plan Validation   |  |                        |  |
|   |  3. CORS (41 origins)     |  |                   |  | discovery-service      |  |
|   |  4. Gateway Context       |  | Two-Layer Cache:  |  | leaderboard-service    |  |
|   |  5. Request ID (UUID)     |  | L1: Ristretto     |  | profile-collection-svc |  |
|   |  6. Request Logger        |  |     10M keys, 1GB |  | post-collection-svc    |  |
|   |  7. AppAuth (JWT)         |  | L2: Redis Cluster  |  | genre-insights-svc     |  |
|   +---------------------------+  |     3-6 nodes     |  | collection-analytics   |  |
|                                  +-------------------+  | keyword-collection     |  |
|   Headers Enriched:                                     | campaign-profile-svc   |  |
|   x-bb-partner-id, x-bb-plan-type,                     | partner-usage-svc      |  |
|   x-bb-requestid, Authorization                        | content-svc            |  |
|                                                         | collection-group-svc   |  |
|                                                         | activity-svc           |  |
|                                                         | social-profile-svc     |  |
+-----------------------------+-------------------------------+-----------------------+
                              |                               |
                              | (Enriched HTTP)               | (Enriched HTTP)
                              v                               v
+-----------------------------+-------------------------------+-----------------------+
|                                                                                     |
|                         [4] COFFEE  (Go, 8.5K LOC)                                  |
|                         Port 7179 | Chi Router | Go 1.18                            |
|                                                                                     |
|   +---------------------------+  +-------------------+  +------------------------+  |
|   | 4-Layer REST Architecture |  | 12 Business       |  | Dual Database          |  |
|   |                           |  | Modules:          |  | Strategy:              |  |
|   | API Layer (Handler)       |  |                   |  |                        |  |
|   |   |                       |  | discovery         |  | PostgreSQL (OLTP)      |  |
|   | Service Layer             |  | profilecollection |  |   27+ tables           |  |
|   |   |                       |  | postcollection    |  |   Collections, Profiles|  |
|   | Manager Layer             |  | leaderboard       |  |   Campaigns, Partners  |  |
|   |   |                       |  | collectionanalyt. |  |                        |  |
|   | DAO Layer (GORM)          |  | genreinsights     |  | ClickHouse (OLAP)      |  |
|   |                           |  | keywordcollection |  |   Time-series queries  |  |
|   | Go Generics for           |  | campaignprofiles  |  |   Analytics aggregation|  |
|   | type-safe CRUD            |  | collectiongroup   |  |   Leaderboard data     |  |
|   |                           |  | partnerusage      |  |                        |  |
|   | 10-Layer Middleware:       |  | content           |  | Redis Cluster          |  |
|   | Prometheus, Sentry,       |  | activitytracker   |  |   Session cache        |  |
|   | Timeout, Session Mgmt     |  |                   |  |   Hot data             |  |
|   +---------------------------+  +-------------------+  +------------------------+  |
|                                                                                     |
|   Multi-Tenant: Partner/Account/User hierarchy | Plan-Based Feature Gating          |
|   50+ REST Endpoints | Watermill AMQP Event Publishing                              |
+-----------------------------+---+---------------------------+-----------------------+
                              |   |                           |
           (Publishes events) |   | (Reads from)              | (Reads/Writes)
                              v   |                           v
+-----------------------------+   |   +-------------------------------------------------+
|                             |   |   |                                                 |
|     RabbitMQ Message        |   |   |              SHARED DATABASES                   |
|     Broker Cluster          |   |   |                                                 |
|                             |   |   |  +------------------+  +---------------------+  |
|  11 Exchanges:              |   |   |  |   PostgreSQL     |  |    ClickHouse       |  |
|  - grpc_event.tx            |   |   |  |   (172.31.2.21)  |  |    (172.31.28.68)   |  |
|  - grpc_event.dx            |   |   |  |                  |  |                     |  |
|  - grpc_event_error.dx      |   |   |  |  Beat Tables:    |  |  Event Tables:      |  |
|  - branch_event.tx          |   |   |  |  instagram_acct  |  |  event              |  |
|  - webengage_event.dx       |   |   |  |  instagram_post  |  |  click_event        |  |
|  - identity.dx              |   |   |  |  youtube_account |  |  branch_event       |  |
|  - beat.dx                  |   |   |  |  youtube_post    |  |  webengage_event    |  |
|  - coffee.dx                |   |   |  |  scrape_req_log  |  |  trace_log          |  |
|  - shopify_event.dx         |   |   |  |  credential      |  |  post_log_event     |  |
|  - affiliate.dx             |   |   |  |  profile_log     |  |  profile_log        |  |
|  - bigboss                  |   |   |  |  post_log        |  |  sentiment_log      |  |
|                             |   |   |  |  asset_log       |  |  shopify_event      |  |
|  26+ Queues:                |   |   |  |                  |  |  affiliate_order    |  |
|  (see Event-gRPC section)   |   |   |  |  Coffee Tables:  |  |  graphy_event       |  |
|                             |   |   |  |  profile_collctn |  |  app_init_event     |  |
|  Durable queues             |   |   |  |  post_collection |  |  entity_metric      |  |
|  Dead letter routing        |   |   |  |  leaderboard     |  |  ab_assignment      |  |
|  Retry logic (max 2)        |   |   |  |  time_series     |  |                     |  |
|  Prefetch QoS               |   |   |  |  partner_usage   |  |  dbt Tables:        |  |
|                             |   |   |  |  activity_tracker|  |  29 stg_* models    |  |
+---------+------+------------+   |   |  |  collection_grp  |  |  83 mart_* models   |  |
          |      |                |   |  +------------------+  +---------------------+  |
          |      |                |   |                                                 |
          |      |                |   |  +------------------+  +---------------------+  |
          |      |                |   |  |   Redis Cluster  |  |    AWS S3            |  |
          |      |                |   |  |   3-6 nodes      |  |  gcc-social-data     |  |
          |      |                |   |  |                  |  |                     |  |
          |      |                |   |  |  Rate limit state|  |  Media assets (CDN) |  |
          |      |                |   |  |  Session data    |  |  Data pipeline tmp  |  |
          |      |                |   |  |  Cache layer     |  |  ClickHouse exports |  |
          |      |                |   |  +------------------+  +---------------------+  |
          |      |                |   +-------------------------------------------------+
          |      |                |
          v      v                v
+---------+------+----------------+-------------------------------------------------+
|                                                                                   |
|                    [2] EVENT-GRPC  (Go, 10K LOC)                                  |
|                    gRPC Port 8017 | HTTP Port 8019 | Gin + gRPC                   |
|                                                                                   |
|  +---------------------------+  +---------------------+  +---------------------+  |
|  | gRPC Server               |  | 6 Worker Pools      |  | 20+ Sinkers         |  |
|  |                           |  |                     |  |                     |  |
|  | EventService.dispatch()   |  | Event Worker Pool   |  | DIRECT Sinkers:     |  |
|  | 60+ event types via       |  | Branch Worker Pool  |  | SinkEventToCH       |  |
|  | protobuf (688 lines)      |  | WebEngage Worker    |  | SinkClickEventToCH  |  |
|  |                           |  | Vidooly Worker      |  | SinkBranchEvent     |  |
|  | HTTP Endpoints:           |  | Graphy Worker       |  | SinkGraphyEvent     |  |
|  | /vidooly/event            |  | Shopify Worker      |  | SinkWebengageToCH   |  |
|  | /branch/event             |  |                     |  | SinkWebengageToAPI  |  |
|  | /webengage/event          |  | Buffered channels   |  |                     |  |
|  | /shopify/event            |  | (1000+ capacity)    |  | BUFFERED Sinkers:   |  |
|  | /heartbeat                |  | Safe goroutines     |  | TraceLogSinker      |  |
|  | /metrics                  |  | Panic recovery      |  | PostLogSinker       |  |
|  |                           |  |                     |  | ProfileLogSinker    |  |
|  | Events grouped by type    |  | Timestamp           |  | ScrapeLogSinker     |  |
|  | Published to RabbitMQ     |  | correction          |  | AffiliateOrderSink  |  |
|  | (grpc_event.tx exchange)  |  |                     |  | ActivityTrackerSink |  |
|  +---------------------------+  +---------------------+  +---------------------+  |
|                                                                                   |
|  26 RabbitMQ Consumers | 90+ Workers | Batch flush: 1000 records or 5 seconds    |
|  Multi-DB: ClickHouse (analytics) + PostgreSQL (transactional)                    |
|  Auto-reconnect: 1-second cron for ClickHouse + RabbitMQ                          |
+-----------------------------------------------------------------------------------+
          ^                    ^                              ^
          |                    |                              |
          | (Events via AMQP)  | (gRPC dispatch)             | (HTTP webhooks)
          |                    |                              |
+---------+--------------------+-----+          +------------+-----------+
|                                    |          |                        |
|   [1] BEAT  (Python, 15K LOC)      |          |  External Platforms    |
|   Port 8000 | FastAPI + uvloop     |          |                        |
|                                    |          |  Instagram Graph API   |
|  +--------------------+            |          |  YouTube Data API v3   |
|  | Worker Pool System |            |          |  RapidAPI (6 sources)  |
|  |                    |            |          |  Shopify API           |
|  | 75+ Scraping Flows |            |          |  Vidooly API           |
|  | 150+ Workers       |            |          |  Branch.io             |
|  | Multiprocessing    |            |          |  WebEngage             |
|  |                    |            |          |  OpenAI GPT            |
|  | Instagram:         |            |          |  Lama API (fallback)   |
|  |  25+ flows         |            |          +------------------------+
|  |  Profile, Posts    |
|  |  Stories, Comments |            +-------------------------------------------+
|  |  Followers, Likes  |            |                                           |
|  |                    |            |   [3] STIR  (Python, 17.5K LOC)            |
|  | YouTube:           |            |   Apache Airflow + dbt + ClickHouse       |
|  |  20+ flows         |            |                                           |
|  |  Channels, Videos  |            |  +---------------------+                  |
|  |  Comments, Playlts |            |  | 76 Airflow DAGs     |                  |
|  |                    |            |  |                     |                  |
|  | GPT Enrichment:    |            |  | 11 dbt orchestration|                  |
|  |  6 flows           |            |  | 17 Instagram sync   |                  |
|  |  Gender, Location  |            |  | 12 YouTube sync     |                  |
|  |  Category, Topics  |            |  | 15 Leaderboard sync |                  |
|  |                    |            |  | 7 Asset upload      |                  |
|  | Shopify:           |            |  | 14 Operational      |                  |
|  |  Orders refresh    |            |  +---------------------+                  |
|  |                    |            |                                           |
|  | Asset Upload:      |            |  +---------------------+                  |
|  |  15 workers        |            |  | 112 dbt Models      |                  |
|  |  S3/CDN delivery   |            |  |                     |                  |
|  +--------------------+            |  | 29 staging (stg_*)  |                  |
|                                    |  |   13 beat source    |                  |
|  +--------------------+            |  |   16 coffee source  |                  |
|  | FastAPI REST API   |            |  |                     |                  |
|  | Rate Limiting:     |            |  | 83 marts (mart_*)   |                  |
|  |  Redis-backed      |            |  |   16 discovery      |                  |
|  |  3-level stacking  |            |  |   14 leaderboard    |                  |
|  |  20K/day global    |            |  |   13 collection     |                  |
|  |  Source-specific   |            |  |   7 genre           |                  |
|  +--------------------+            |  |   8 profile stats   |                  |
|                                    |  |   4 audience        |                  |
|  +--------------------+            |  |   3 posts           |                  |
|  | 5 AMQP Listeners  |            |  |   9 staging coll    |                  |
|  | credentials_valid  |            |  |   1 orders          |                  |
|  | identity_token     |            |  |   8 partial updates |                  |
|  | keyword_collection |            |  +---------------------+                  |
|  | sentiment_analysis |            |                                           |
|  | sentiment_report   |            |  Data Flow:                               |
|  +--------------------+            |  ClickHouse --> S3 --> SSH --> PostgreSQL  |
|                                    |  Atomic table swap (zero-downtime)        |
|  15+ API Integrations              |                                           |
|  SQL-based task queue              |  Schedules: */5min to weekly              |
|  FOR UPDATE SKIP LOCKED            |  Slack alerts on failure                  |
+------------------------------------+-------------------------------------------+


+-----------------------------------------------------------------------------------+
|                                                                                   |
|           [6] FAKE FOLLOWER ANALYSIS  (Python, 1K LOC)                            |
|           AWS Lambda + ECR | Serverless                                           |
|                                                                                   |
|  +-------------------------------+  +------------------------------------------+  |
|  | ML Detection Pipeline         |  | AWS Infrastructure                       |  |
|  |                               |  |                                          |  |
|  | 5-Feature Ensemble:           |  | ClickHouse --> S3 --> SQS --> Lambda     |  |
|  |  1. Non-Indic script detect   |  |                        |                 |  |
|  |  2. Digit count (>4 = fake)   |  |                        v                 |  |
|  |  3. Handle-name correlation   |  |               ECR Container (Py 3.10)    |  |
|  |  4. Fuzzy similarity (0-100)  |  |               fake.handler()             |  |
|  |  5. Indian name DB (35K+)     |  |                        |                 |  |
|  |                               |  |                        v                 |  |
|  | 10 Indic Scripts:             |  |               Kinesis --> Output JSON     |  |
|  |  Hindi, Bengali, Gujarati,    |  |                                          |  |
|  |  Kannada, Malayalam, Odia,    |  | SQS: 8 parallel workers                  |  |
|  |  Punjabi, Tamil, Telugu, Urdu |  | Kinesis: ON_DEMAND auto-scaling          |  |
|  |                               |  | Lambda: 10-20 records/sec                |  |
|  | HMM-based transliteration     |  |                                          |  |
|  | RapidFuzz fuzzy matching      |  | Scores: 0.0 (real), 0.33 (weak fake),   |  |
|  | 13 Unicode symbol variants    |  |         1.0 (definitely fake)            |  |
|  +-------------------------------+  +------------------------------------------+  |
+-----------------------------------------------------------------------------------+
```

---

## 2. DATA FLOW DIAGRAMS BY USE CASE

### 2.1 Influencer Profile Crawl --> Event --> Analytics --> API Serving

This is the primary end-to-end data flow. A brand searches for an influencer on the SaaS dashboard, and data flows through all 6 projects.

```
USE CASE: Brand searches for influencer "virat.kohli" on goodcreator.co

STEP 1: API REQUEST
===================
  Brand Dashboard (Web)
       |
       | GET /discovery-service/api/profile/INSTAGRAM/virat.kohli
       v
  [SaaS Gateway] Port 8009
       |
       | JWT validation --> Redis session lookup --> Plan check
       | Enrich headers: x-bb-partner-id=123, x-bb-plan-type=PRO
       v
  [Coffee] Port 7179
       |
       | discovery/api/api.go --> service/service.go --> manager/searchmanager.go
       | DAO queries ClickHouse for analytics data + PostgreSQL for profile data
       |
       | CACHE HIT? --> Return cached profile from Redis
       | CACHE MISS? --> Trigger scrape via Beat
       v

STEP 2: DATA COLLECTION (if profile stale or missing)
======================================================
  [Coffee] publishes to RabbitMQ
       |
       | coffee.dx exchange --> beat.dx --> scrape_request queue
       v
  [Beat] Worker Pool picks up task
       |
       | poll() --> SELECT ... FROM scrape_request_log
       |            WHERE flow='refresh_profile_by_handle'
       |            FOR UPDATE SKIP LOCKED
       v
  Instagram Crawler (Strategy Pattern)
       |
       | Try GraphAPI (primary) --> 6 RapidAPI sources (fallback) --> Lama (last resort)
       | Rate limited: Redis-backed, 20K/day global, 60/min, 1/sec per handle
       v
  Raw API Response
       |
       | Stage 1: RETRIEVAL  (instagram/tasks/retrieval.py)
       | Stage 2: PARSING    (instagram/tasks/ingestion.py)  --> InstagramProfileLog
       | Stage 3: PROCESSING (instagram/tasks/processing.py) --> Upsert to PostgreSQL
       v

STEP 3: EVENT INGESTION
========================
  [Beat] publishes profile_log event to RabbitMQ
       |
       | beat.dx exchange --> profile_log_events_rk routing key
       v
  [Event-gRPC] Consumer: profile_log_events_q (2 workers)
       |
       | BufferProfileLogEvent() --> adds to channel (capacity 10,000)
       | ProfileLogEventsSinker() --> batches 1000 records OR flushes every 5 seconds
       v
  ClickHouse: _e.profile_log table (columnar, compressed)

STEP 4: DATA TRANSFORMATION
============================
  [Stir] Airflow DAG: dbt_core (runs every 15 minutes)
       |
       | DbtRunOperator --> dbt run --models tag:core
       |
       | stg_beat_instagram_account  (staging: raw data extraction)
       |         |
       |         v
       | mart_instagram_account      (mart: engagement rates, ranks, grades)
       |         |
       |         v
       | mart_instagram_leaderboard  (mart: multi-dimensional rankings)
       v

STEP 5: SYNC TO SERVING LAYER
===============================
  [Stir] Airflow DAG: sync_leaderboard_prod (daily at 20:15 UTC)
       |
       | Task 1: ClickHouseOperator --> INSERT INTO FUNCTION s3(...) SELECT * FROM dbt.mart_leaderboard
       |         Exports to S3: gcc-social-data/data-pipeline/tmp/leaderboard.json
       |
       | Task 2: SSHOperator --> aws s3 cp s3://gcc-social-data/.../leaderboard.json /tmp/
       |         Downloads to PostgreSQL server
       |
       | Task 3: PostgresOperator --> CREATE TEMP TABLE, COPY, INSERT with JSONB parsing
       |
       | Task 4: PostgresOperator --> ALTER TABLE leaderboard RENAME TO leaderboard_old;
       |                              ALTER TABLE leaderboard_new RENAME TO leaderboard;
       |         ATOMIC TABLE SWAP (zero downtime)
       v

STEP 6: API RESPONSE
=====================
  [Coffee] serves from PostgreSQL (fresh data) + ClickHouse (analytics)
       |
       | StandardResponse: {data: {...profile, engagement_rate, rank}, status: "SUCCESS"}
       v
  [SaaS Gateway] forwards response to Brand Dashboard
       |
       | Plan-based gating: FREE users get limited fields
       | PRO users get full audience demographics
       v
  Brand Dashboard displays influencer profile with analytics
```

### 2.2 Leaderboard Generation Pipeline

```
USE CASE: Daily leaderboard computation for 500K+ influencers

  [Beat] continuously scrapes profiles
       |
       | 150+ workers crawling Instagram/YouTube profiles 24/7
       | Each profile refresh publishes profile_log event
       v
  [Event-gRPC] ingests profile_log events into ClickHouse
       |
       | Buffered sinker: batches of 1000 records
       v
  ClickHouse: profile_log table (time-series data)
       |
       v
  [Stir] dbt_core DAG (every 15 minutes)
       |
       | +------------------------------------------------------------+
       | | stg_beat_instagram_account (staging)                       |
       | |   SELECT profile_id, handle, followers, following,         |
       | |          avg_likes, avg_comments, engagement_rate           |
       | |   FROM beat_replica.instagram_account                      |
       | +------------------------------------------------------------+
       |         |
       |         v
       | +------------------------------------------------------------+
       | | mart_instagram_account (core mart)                         |
       | |   - Calculates: engagement_rate, avg_likes, avg_comments   |
       | |   - Rankings: followers_rank, engagement_rank               |
       | |   - Rankings by category: followers_rank_by_cat             |
       | |   - Rankings by language: followers_rank_by_lang            |
       | |   - Combined: followers_rank_by_cat_lang                   |
       | +------------------------------------------------------------+
       |         |
       |         v
       | +------------------------------------------------------------+
       | | mart_leaderboard (daily aggregation)                       |
       | |   - Cross-platform rankings (Instagram + YouTube)          |
       | |   - Rank change tracking (vs. previous month)              |
       | |   - Segmented by category, language, country               |
       | +------------------------------------------------------------+
       |         |
       |         v
       | +------------------------------------------------------------+
       | | mart_time_series (incremental, ReplacingMergeTree)         |
       | |   - Daily snapshots of follower counts                     |
       | |   - Growth rate calculations                               |
       | |   - argMax(followers_count, created_at) per day            |
       | +------------------------------------------------------------+
       v
  [Stir] sync_leaderboard_prod DAG (daily at 20:15 UTC)
       |
       | ClickHouse --> S3 (JSONEachRow) --> SSH download --> PostgreSQL
       | COPY from /tmp/*.json --> JSONB parsing --> Atomic table swap
       v
  [Coffee] leaderboard-service endpoints
       |
       | POST /leaderboard-service/api/leaderboard/platform/INSTAGRAM
       | POST /leaderboard-service/api/leaderboard/platform/INSTAGRAM/category/fashion
       | POST /leaderboard-service/api/leaderboard/cross-platform
       v
  [SaaS Gateway] serves to brand dashboards
```

### 2.3 Asset Upload Pipeline

```
USE CASE: 8M+ images per day (profile pics, post thumbnails, story media)

  [Stir] Airflow DAGs trigger asset upload tasks:
       |
       | upload_post_asset.py          (every 3 hours)
       | upload_post_asset_stories.py  (every 3 hours)
       | upload_insta_profile_asset.py (hourly)
       |
       | Creates scrape_request_log entries in PostgreSQL
       | flow = 'asset_upload_flow', status = 'PENDING'
       v
  [Beat] Asset Worker Pool
       |
       | main_assets.py: 15 workers, 5 concurrency each = 75 concurrent uploads
       |
       | poll() --> SELECT ... FROM scrape_request_log
       |            WHERE flow='asset_upload_flow'
       |            FOR UPDATE SKIP LOCKED
       v
  Download from Platform CDN
       |
       | Instagram/YouTube media URL --> HTTP download (aiohttp)
       | Rate limited per source
       v
  Upload to AWS S3
       |
       | Bucket: gcc-social-data
       | Path: /media/{platform}/{profile_id}/{asset_type}/{hash}.{ext}
       | CloudFront CDN distribution
       v
  Update PostgreSQL
       |
       | UPDATE instagram_post SET thumbnail_url = 'cdn.goodcreator.co/...'
       | INSERT INTO asset_log (profile_id, asset_type, s3_url, cdn_url)
       v
  Publish event to RabbitMQ
       |
       | beat.dx --> asset_log_events_rk
       v
  [Event-gRPC] sinks asset_log event to ClickHouse for tracking
```

### 2.4 Fake Follower Detection Pipeline

```
USE CASE: Analyze follower authenticity for a creator's audience

  [Stir] / Manual Trigger
       |
       | ClickHouse query: Extract follower data
       | SELECT handle, follower_handle, follower_full_name
       | FROM dbt.stg_beat_profile_relationship_log
       | UNION ALL
       | SELECT ... FROM _e.profile_relationship_log_events
       v
  Export to S3
       |
       | INSERT INTO FUNCTION s3('gcc-social-data/temp/{date}_creator_followers.json')
       v
  [Fake Follower Analysis] push.py
       |
       | Download from S3 --> Read 10,000 lines at a time
       | Divide into 8 buckets (round-robin)
       | 8 parallel workers (multiprocessing.Pool)
       | Batch send to SQS (10 messages per API call)
       v
  AWS SQS Queue: creator_follower_in
       |
       | MaximumMessageSize: 256KB
       | MessageRetentionPeriod: 4 days
       v
  AWS Lambda (ECR Container)
       |
       | fake.handler(event, context)
       |
       | +------------------------------------------------------+
       | | STEP 1: Symbol conversion (13 Unicode variants)      |
       | |   "Fancy Text" --> "Alice"                          |
       | |                                                      |
       | | STEP 2: Language detection (non-Indic = fake signal) |
       | |   Greek, Armenian, Chinese, Korean --> flag           |
       | |                                                      |
       | | STEP 3: Indic transliteration (10 scripts)           |
       | |   "rahul" (Hindi) --> "Rahul" (English)              |
       | |   Uses HMM-based ML models per language              |
       | |                                                      |
       | | STEP 4: Feature extraction (5 features)              |
       | |   1. Non-Indic language flag (0/1)                   |
       | |   2. Digit count > 4 (0/1)                           |
       | |   3. Handle-name correlation (0/1/2)                 |
       | |   4. Fuzzy similarity score (0-100)                  |
       | |   5. Indian name DB match (0-100)                    |
       | |                                                      |
       | | STEP 5: Ensemble scoring                             |
       | |   0.0 = real, 0.33 = weak fake, 1.0 = definitely    |
       | +------------------------------------------------------+
       v
  AWS Kinesis Stream: creator_out (ON_DEMAND auto-scaling)
       |
       | PartitionKey: follower_handle
       v
  [Fake Follower Analysis] pull.py
       |
       | Multi-shard parallel read from Kinesis
       | Output: {date}_creator_followers_final_fake_analysis.json
       v
  Results loaded back to ClickHouse / consumed by Stir dbt models
       |
       | mart_instagram_creators_followers_fake_analysis
       v
  [Coffee] surfaces fake follower % on profile pages
```

### 2.5 Real-Time Event Tracking Pipeline

```
USE CASE: Track user interactions on goodcreator.co (clicks, page views, purchases)

  Mobile App / Web App
       |
       | gRPC (protobuf): EventService.dispatch(Events)
       | 60+ event types defined in eventservice.proto (688 lines)
       |
       | Event types:
       |   LaunchEvent, PageOpenedEvent, WidgetViewEvent, WidgetCtaClickedEvent,
       |   AddToCartEvent, InitiatePurchaseEvent, CompletePurchaseEvent,
       |   StreamEnterEvent, PhoneVerificationInitiateEvent, etc.
       v
  [Event-gRPC] gRPC Server (Port 8017)
       |
       | Event Worker Pool
       |   - Buffered channel (1000 capacity)
       |   - Events grouped by type (reflect.TypeOf)
       |   - Timestamp correction (future timestamps --> now)
       |   - Protobuf --> JSON serialization
       v
  RabbitMQ Exchange: grpc_event.tx
       |
       | Routing key = event type name
       | Fanout to multiple consumers
       |
       +------> grpc_clickhouse_event_q (2 workers)
       |              |
       |              v
       |        SinkEventToClickhouse()
       |              |
       |              v
       |        ClickHouse: event table
       |        (event_id, event_name, event_params JSONB, session_id, user_id, ...)
       |
       +------> clickhouse_click_event_q (2 workers)
       |              |
       |              v
       |        SinkClickEventToClickhouse() -- 600+ lines of click tracking logic
       |              |
       |              v
       |        ClickHouse: click_event table
       |
       +------> app_init_event_q (2 workers)
       |              |
       |              v
       |        SinkAppInitEvent()
       |              |
       |              v
       |        ClickHouse: app_init_event table
       |
       +------> webengage_event_q (3 workers)
                      |
                      v
                SinkWebengageToAPI() -- forwards to WebEngage REST API
                      |
                      v
                WebEngage platform (marketing automation)

  ALSO: HTTP webhook ingestion (same pipeline)
       |
       | POST /branch/event   --> Branch.io attribution data --> ClickHouse
       | POST /vidooly/event  --> Vidooly analytics --> ClickHouse
       | POST /shopify/event  --> Shopify commerce events --> ClickHouse
       | POST /webengage/event --> WebEngage marketing events --> ClickHouse + API
```

---

## 3. TECHNOLOGY STACK MATRIX

| Component | Beat | Event-gRPC | Stir | Coffee | SaaS Gateway | Fake Follower |
|-----------|------|-----------|------|--------|-------------|---------------|
| **Language** | Python 3.11 | Go 1.14 | Python 3.x | Go 1.18 | Go 1.12 | Python 3.10 |
| **Framework** | FastAPI | Gin + gRPC | Airflow 2.6.3 | Chi Router | Gin v1.8.1 | AWS Lambda |
| **Lines of Code** | 15,000+ | 10,000+ | 17,500+ | 8,500+ | 2,200+ | 955+ |
| **Primary DB** | PostgreSQL | ClickHouse | ClickHouse | PostgreSQL + ClickHouse | -- (proxy) | ClickHouse |
| **Message Queue** | RabbitMQ (aio-pika) | RabbitMQ (amqp) | RabbitMQ (pika) | RabbitMQ (Watermill) | -- | AWS SQS |
| **Cache** | Redis (rate limiting) | Redis + Ristretto | -- | Redis (sessions) | Redis + Ristretto | -- |
| **ORM/Driver** | SQLAlchemy (async) | GORM | dbt + clickhouse-connect | GORM | -- | clickhouse_connect |
| **Event Loop** | uvloop (async) | goroutines | Airflow scheduler | goroutines | goroutines | Lambda runtime |
| **Serialization** | JSON | Protobuf + JSON | JSONEachRow | JSON | JSON | JSON |
| **Monitoring** | -- | Prometheus + Sentry | Slack alerts | Prometheus + Sentry | Prometheus + Sentry | CloudWatch |
| **CI/CD** | GitLab CI | GitLab CI | GitLab CI (manual) | GitLab CI | GitLab CI | ECR + Lambda deploy |
| **Port** | 8000 | 8017 (gRPC), 8019 (HTTP) | Airflow UI | 7179 | 8009 | Lambda (serverless) |
| **Concurrency Model** | Multiprocessing + asyncio Semaphore | Goroutine pools + channels | Airflow workers + dbt threads | Goroutine per request | Goroutine per request | Lambda auto-scale |
| **Prod Nodes** | 2 (beat-deployer-1, beat-deployer-2) | 1 | 1 (beat-deployer-1) | 2 (cb1-1, cb2-1) | 2 (cb1-1, cb2-1) | Serverless |

---

## 4. INTER-PROJECT COMMUNICATION

### 4.1 Communication Channels Summary

```
+------------------+     RabbitMQ (beat.dx)      +------------------+
|                  | --------------------------> |                  |
|      BEAT        |     RabbitMQ (coffee.dx)     |    EVENT-GRPC   |
|   (Scraper)      | <-------------------------- |   (Ingestion)    |
|                  |     HTTP (Beat API)          |                  |
+------------------+ <-------------------------- +------------------+
        |                                                |
        | PostgreSQL (shared DB)                         | ClickHouse (writes)
        | RabbitMQ (beat.dx)                             |
        v                                                v
+------------------+     ClickHouse (reads)       +------------------+
|                  | --------------------------> |                  |
|      STIR        |     S3 (data staging)       |     COFFEE       |
|  (Data Platform) | --------------------------> |   (SaaS API)     |
|                  |     PostgreSQL (writes)      |                  |
+------------------+ --------------------------> +------------------+
        |                                                |
        | ClickHouse (exports)                           |
        v                                                v
+------------------+                              +------------------+
|  FAKE FOLLOWER   |                              |   SAAS GATEWAY   |
|   ANALYSIS       |                              |   (API Gateway)  |
|  (ML Detection)  |                              |  Proxies Coffee  |
+------------------+                              +------------------+
```

### 4.2 RabbitMQ Exchanges and Routing

| Exchange | Type | Publishers | Consumers | Purpose |
|----------|------|-----------|-----------|---------|
| **beat.dx** | direct | Beat | Event-gRPC, Beat | Scraping events, credentials, keyword collections |
| **coffee.dx** | topic | Coffee | Coffee, Event-gRPC | Collection events, profile updates |
| **grpc_event.tx** | topic | Event-gRPC | Event-gRPC | App events --> ClickHouse sinkers |
| **grpc_event.dx** | direct | Event-gRPC | Event-gRPC | Direct event routing |
| **grpc_event_error.dx** | direct | Event-gRPC | Event-gRPC | Dead letter / error events |
| **branch_event.tx** | topic | Event-gRPC | Event-gRPC | Branch.io attribution events |
| **webengage_event.dx** | direct | Event-gRPC | Event-gRPC | WebEngage marketing events |
| **identity.dx** | direct | Identity Service | Beat, Event-gRPC | Auth tokens, trace logs |
| **shopify_event.dx** | direct | Event-gRPC | Event-gRPC | Shopify commerce events |
| **affiliate.dx** | direct | External | Event-gRPC | Affiliate order tracking |
| **bigboss** | direct | External | Event-gRPC | BigBoss voting events |

### 4.3 Shared Database Access

| Database | Beat | Event-gRPC | Stir | Coffee | SaaS Gateway | Fake Follower |
|----------|------|-----------|------|--------|-------------|---------------|
| **PostgreSQL (beat DB)** | READ/WRITE | WRITE (referrals) | READ (beat_replica source) + WRITE (sync targets) | READ/WRITE | -- | -- |
| **ClickHouse (dbt DB)** | -- | WRITE (18+ tables) | READ/WRITE (112 dbt models) | READ (analytics) | -- | READ (exports) |
| **ClickHouse (_e DB)** | -- | WRITE (events) | READ (real-time events) | -- | -- | READ (follower events) |
| **Redis Cluster** | READ/WRITE (rate limits) | READ/WRITE (cache) | -- | READ/WRITE (sessions) | READ/WRITE (auth cache) | -- |
| **AWS S3** | WRITE (media assets) | -- | READ/WRITE (data pipeline staging) | -- | -- | READ/WRITE (data pipeline) |

### 4.4 HTTP/API Communication

| Caller | Callee | Endpoint | Purpose |
|--------|--------|----------|---------|
| SaaS Gateway | Coffee | `http://coffee.goodcreator.co:7179/*` | All SaaS API calls proxied |
| SaaS Gateway | Identity Service | `/auth/verify/token` | JWT token validation |
| SaaS Gateway | Partner Service | `/partner/{id}` | Plan/contract lookup |
| Stir (Airflow) | Beat | `POST /scrape_request_log/flow/{flow}` | Trigger scrape tasks |
| Coffee | Beat | `POST /scrape_request_log/flow/{flow}` | On-demand profile refresh |
| Beat | Instagram Graph API | `graph.facebook.com/v15.0/*` | Profile/post data |
| Beat | YouTube Data API v3 | `youtube.googleapis.com/v3/*` | Channel/video data |
| Beat | 6 RapidAPI sources | Various endpoints | Fallback scraping |
| Beat | OpenAI (Azure) | GPT-3.5-turbo | Data enrichment |
| Event-gRPC | WebEngage API | REST endpoint | Forward marketing events |
| Event-gRPC | Identity Service | REST endpoint | User account operations |

---

## 5. DATABASE ARCHITECTURE

### 5.1 PostgreSQL Tables (Shared by Beat + Coffee + Stir)

```
+=============================================================================+
|                     POSTGRESQL (172.31.2.21:5432)                            |
|                     Database: beat | Schema: public                          |
+=============================================================================+
|                                                                             |
|  BEAT-OWNED TABLES (Core Data):                                             |
|  +----------------------------+  +----------------------------+             |
|  | instagram_account          |  | youtube_account            |             |
|  | - id (PK)                  |  | - id (PK)                  |             |
|  | - profile_id (unique)      |  | - channel_id (unique)      |             |
|  | - handle (indexed)         |  | - title                    |             |
|  | - followers, following     |  | - subscribers              |             |
|  | - avg_likes, avg_comments  |  | - total_views, avg_views   |             |
|  | - engagement_rate          |  | - plays, shorts_reach      |             |
|  | - profile_type             |  | - category, language       |             |
|  | - is_verified, is_business |  | - created_at, updated_at   |             |
|  | - audience_gender (JSONB)  |  +----------------------------+             |
|  | - audience_age (JSONB)     |                                             |
|  | - audience_location (JSONB)|  +----------------------------+             |
|  +----------------------------+  | instagram_post             |             |
|                                  | - post_id (unique)         |             |
|  +----------------------------+  | - shortcode (indexed)      |             |
|  | scrape_request_log         |  | - profile_id (FK)          |             |
|  | - id (PK)                  |  | - post_type, caption       |             |
|  | - flow (75+ types)         |  | - likes, comments, plays   |             |
|  | - status (PENDING/PROC/OK) |  | - reach, shares, saved     |             |
|  | - params (JSONB)           |  | - publish_time             |             |
|  | - priority, retry_count    |  +----------------------------+             |
|  | - FOR UPDATE SKIP LOCKED   |                                             |
|  +----------------------------+  +----------------------------+             |
|                                  | youtube_post               |             |
|  +----------------------------+  | - video_id (unique)        |             |
|  | credential                 |  | - channel_id (FK)          |             |
|  | - source (graphapi, ytapi) |  | - title, description       |             |
|  | - credentials (JSONB)      |  | - views, likes, comments   |             |
|  | - enabled, disabled_till   |  | - publish_time             |             |
|  +----------------------------+  +----------------------------+             |
|                                                                             |
|  +----------------------------+  +----------------------------+             |
|  | profile_log (audit)        |  | post_log (audit)           |             |
|  | - dimensions (JSONB)       |  | - dimensions (JSONB)       |             |
|  | - metrics (JSONB)          |  | - metrics (JSONB)          |             |
|  | - source, timestamp        |  | - source, timestamp        |             |
|  +----------------------------+  +----------------------------+             |
|                                                                             |
|  +----------------------------+  +----------------------------+             |
|  | asset_log                  |  | sentiment_log              |             |
|  | - profile_id, asset_type   |  | - comment_id, score        |             |
|  | - s3_url, cdn_url          |  | - positive, negative       |             |
|  +----------------------------+  +----------------------------+             |
|                                                                             |
|  COFFEE-OWNED TABLES (SaaS Data):                                           |
|  +----------------------------+  +----------------------------+             |
|  | profile_collection         |  | post_collection            |             |
|  | - name, partner_id         |  | - name, partner_id         |             |
|  | - share_id (unique)        |  | - ingestion_frequency      |             |
|  | - source (SAAS, GCC_CAMP.) |  | - show_in_report           |             |
|  +----------------------------+  +----------------------------+             |
|                                                                             |
|  +----------------------------+  +----------------------------+             |
|  | profile_collection_item    |  | post_collection_item       |             |
|  | - collection_id (FK)       |  | - collection_id (FK)       |             |
|  | - platform_profile_id      |  | - post_shortcode           |             |
|  | - custom_data (JSONB)      |  | - platform                 |             |
|  +----------------------------+  +----------------------------+             |
|                                                                             |
|  +----------------------------+  +----------------------------+             |
|  | leaderboard                |  | social_profile_time_series |             |
|  | - platform, profile_id     |  | - platform, profile_id     |             |
|  | - month (date)             |  | - date, followers           |             |
|  | - followers_rank           |  | - following, engagement     |             |
|  | - engagement_rate_rank     |  +----------------------------+             |
|  | - *_rank_by_cat            |                                             |
|  | - *_rank_by_lang           |  +----------------------------+             |
|  +----------------------------+  | collection_post_metrics    |             |
|                                  | - collection_id             |             |
|  +----------------------------+  | - post_id, likes, comments |             |
|  | partner_usage              |  | - engagement_rate           |             |
|  | - partner_id, module       |  +----------------------------+             |
|  | - usage_count, limit_count |                                             |
|  | - plan_type, valid_from/to |  +----------------------------+             |
|  +----------------------------+  | activity_tracker           |             |
|                                  | - partner_id, user_id      |             |
|  +----------------------------+  | - activity_type, metadata  |             |
|  | partner_profile_page_track |  +----------------------------+             |
|  | - partner_id, platform     |                                             |
|  | - access_count             |  +----------------------------+             |
|  +----------------------------+  | collection_group           |             |
|                                  | - name, partner_id         |             |
|  Total: 27+ tables              | - collection_ids (JSONB)   |             |
|                                  +----------------------------+             |
+=============================================================================+
```

### 5.2 ClickHouse Tables (Shared by Event-gRPC + Stir + Coffee)

```
+=============================================================================+
|                    CLICKHOUSE (172.31.28.68:9000)                            |
|                    Multiple Databases                                        |
+=============================================================================+
|                                                                             |
|  DATABASE: _e (Event-gRPC writes)                                           |
|  +----------------------------+  +----------------------------+             |
|  | event                      |  | click_event               |             |
|  | - event_id, event_name     |  | - 600+ lines of click     |             |
|  | - event_timestamp          |  |   tracking columns         |             |
|  | - session_id, user_id      |  | - page, element, action    |             |
|  | - device_id, channel, os   |  +----------------------------+             |
|  | - utm_* (5 params)         |                                             |
|  | - event_params (String/JSON)|  +----------------------------+             |
|  | ENGINE: MergeTree()        |  | branch_event (49 fields)  |             |
|  +----------------------------+  | - attribution tracking     |             |
|                                  | - geo data, campaign data  |             |
|  +----------------------------+  | - event_data (JSONB)       |             |
|  | error_event                |  +----------------------------+             |
|  | - all event fields + error |                                             |
|  +----------------------------+  +----------------------------+             |
|                                  | webengage_event            |             |
|  +----------------------------+  | - userId, eventData        |             |
|  | trace_log                  |  +----------------------------+             |
|  | - hostName, serviceName    |                                             |
|  | - timeTaken                |  +----------------------------+             |
|  +----------------------------+  | graphy_event               |             |
|                                  +----------------------------+             |
|  +----------------------------+                                             |
|  | app_init_event             |  +----------------------------+             |
|  | - device, session metrics  |  | shopify_event              |             |
|  +----------------------------+  +----------------------------+             |
|                                                                             |
|  +----------------------------+  +----------------------------+             |
|  | post_log_event             |  | profile_log                |             |
|  | - post activity tracking   |  | - profile metrics over time|             |
|  +----------------------------+  +----------------------------+             |
|                                                                             |
|  +----------------------------+  +----------------------------+             |
|  | sentiment_log              |  | profile_relationship_log   |             |
|  | - comment sentiment scores |  | - follower/following data  |             |
|  +----------------------------+  +----------------------------+             |
|                                                                             |
|  +----------------------------+  +----------------------------+             |
|  | order_log                  |  | affiliate_order            |             |
|  +----------------------------+  +----------------------------+             |
|                                                                             |
|  +----------------------------+  +----------------------------+             |
|  | ab_assignment              |  | entity_metric              |             |
|  | - A/B test variant data    |  | - generic metric storage   |             |
|  +----------------------------+  +----------------------------+             |
|                                                                             |
|  +----------------------------+                                             |
|  | activity_tracker           |  Total: 18+ tables                          |
|  +----------------------------+                                             |
|                                                                             |
|  DATABASE: dbt (Stir writes, Coffee reads)                                  |
|  +----------------------------+  +----------------------------+             |
|  | STAGING LAYER (29 models)  |  | MART LAYER (83 models)    |             |
|  |                            |  |                            |             |
|  | Beat Sources (13):         |  | Audience (4):              |             |
|  | stg_beat_instagram_account |  | mart_audience_info         |             |
|  | stg_beat_instagram_post    |  | mart_audience_info_gpt     |             |
|  | stg_beat_youtube_account   |  |                            |             |
|  | stg_beat_youtube_post      |  | Collection (13):           |             |
|  | stg_beat_profile_log       |  | mart_collection_post       |             |
|  | stg_beat_post_log          |  | mart_collection_post_ts    |             |
|  | stg_beat_scrape_request_log|  | mart_fake_events           |             |
|  | stg_beat_asset_log         |  |                            |             |
|  | stg_beat_order             |  | Discovery (16):            |             |
|  | stg_beat_profile_insights  |  | mart_instagram_account     |             |
|  | stg_beat_yt_profile_insght |  | mart_youtube_account       |             |
|  | stg_beat_post_activity_log |  | mart_instagram_gpt_data    |             |
|  | stg_beat_profile_rel_log   |  | mart_linked_socials        |             |
|  |                            |  |                            |             |
|  | Coffee Sources (16):       |  | Leaderboard (14):          |             |
|  | stg_coffee_profile_collctn |  | mart_leaderboard           |             |
|  | stg_coffee_post_collection |  | mart_cross_platform_ldrbd  |             |
|  | stg_coffee_campaign_prfls  |  | mart_time_series           |             |
|  | stg_coffee_collection_grp  |  | mart_insta_leaderboard     |             |
|  | stg_coffee_keyword_collctn |  | mart_youtube_leaderboard   |             |
|  | stg_coffee_activity_tracker|  | mart_insta_account_monthly |             |
|  | stg_coffee_*_account_lite  |  |                            |             |
|  | stg_coffee_*_collection_itm|  | Genre (7):                 |             |
|  +----------------------------+  | mart_genre_overview        |             |
|                                  | mart_genre_trending_content|             |
|                                  |                            |             |
|                                  | Profile Stats (14):        |             |
|                                  | mart_insta_account_summary |             |
|                                  | mart_insta_all_posts       |             |
|                                  | mart_insta_creators_fake   |             |
|                                  | mart_insta_recent_stats    |             |
|                                  |                            |             |
|                                  | Posts (3):                 |             |
|                                  | mart_instagram_post_hashtag|             |
|                                  | mart_post_location         |             |
|                                  |                            |             |
|                                  | Orders (1):                |             |
|                                  | mart_gcc_orders            |             |
|                                  |                            |             |
|                                  | Staging Collection (9):    |             |
|                                  | mart_staging_*             |             |
|                                  +----------------------------+             |
|                                                                             |
|  ClickHouse-Specific Features Used:                                         |
|  - ReplacingMergeTree (upserts for time-series)                             |
|  - MergeTree with ORDER BY for query optimization                           |
|  - PARTITION BY toYYYYMM(date) for pruning                                  |
|  - argMax(value, timestamp) for latest value extraction                     |
|  - arrayMap/splitByChar for hashtag normalization                           |
|  - INSERT INTO FUNCTION s3(...) for S3 export                              |
|  - JSONExtractString for JSONB field extraction                            |
+=============================================================================+
```

### 5.3 Redis Usage Across All Projects

```
+=============================================================================+
|                    REDIS CLUSTER (3-6 nodes)                                |
|                    Production: bulbul-redis-prod-1:6379, -2:6380, -3:6381   |
+=============================================================================+
|                                                                             |
|  SAAS GATEWAY:                                                              |
|  +----------------------------+                                             |
|  | session:{sessionId}        | JWT session validation                      |
|  | Pool: 100 connections      | EXISTS check for auth                       |
|  +----------------------------+                                             |
|                                                                             |
|  BEAT:                                                                      |
|  +----------------------------+                                             |
|  | beat_server_*              | Multi-level rate limiting                    |
|  | - refresh_profile_daily    | Global: 20K/day                             |
|  | - refresh_profile_minute   | Per-minute: 60/min                          |
|  | - refresh_profile_{handle} | Per-handle: 1/sec                           |
|  | - source-specific limits   | 2-850 requests per period                   |
|  +----------------------------+                                             |
|                                                                             |
|  EVENT-GRPC:                                                                |
|  +----------------------------+                                             |
|  | Ristretto (in-memory, L1)  | Per-instance LFU cache (10M keys, 1GB)     |
|  | Redis Cluster (L2)         | Shared cache across instances               |
|  +----------------------------+                                             |
|                                                                             |
|  COFFEE:                                                                    |
|  +----------------------------+                                             |
|  | Session data               | Request-scoped sessions                     |
|  | Hot profile data           | Frequently accessed profiles                |
|  +----------------------------+                                             |
+=============================================================================+
```

---

## 6. DEPLOYMENT ARCHITECTURE

### 6.1 Server Topology

```
+=============================================================================+
|                         PRODUCTION INFRASTRUCTURE                           |
+=============================================================================+
|                                                                             |
|  LOAD BALANCER                                                              |
|       |                                                                     |
|       +--------+--------+                                                   |
|                |        |                                                   |
|  +-------------+--+  +--+-------------+                                     |
|  | Node: cb1-1    |  | Node: cb2-1    |  (2 application nodes)              |
|  |                |  |                |                                     |
|  | SaaS Gateway   |  | SaaS Gateway   |  Port 8009                         |
|  | Coffee         |  | Coffee         |  Port 7179                         |
|  +----------------+  +----------------+                                     |
|                                                                             |
|  +---------------------------------------+                                  |
|  | Node: beat-deployer-1                 |  (Beat primary)                  |
|  |                                       |                                  |
|  | Beat workers   (main.py)              |  150+ worker processes           |
|  | Beat API       (server.py)            |  Port 8000                       |
|  | Beat assets    (main_assets.py)       |  75 concurrent uploads           |
|  | Stir/Airflow   (standalone)           |  76 DAGs                         |
|  +---------------------------------------+                                  |
|                                                                             |
|  +---------------------------------------+                                  |
|  | Node: beat-deployer-2                 |  (Beat secondary)                |
|  |                                       |                                  |
|  | Beat workers   (main.py)              |  Additional worker capacity      |
|  | Beat API       (server.py)            |  Port 8000                       |
|  +---------------------------------------+                                  |
|                                                                             |
|  +---------------------------------------+                                  |
|  | Node: deployer-stage-search           |  (Event-gRPC)                    |
|  |                                       |                                  |
|  | Event-gRPC server                     |  gRPC: 8017, HTTP: 8019          |
|  | 26 RabbitMQ consumers                 |  90+ workers                     |
|  +---------------------------------------+                                  |
|                                                                             |
|  +---------------------------------------+                                  |
|  | Database Nodes                        |                                  |
|  |                                       |                                  |
|  | PostgreSQL     (172.31.2.21:5432)     |  27+ tables                      |
|  | ClickHouse     (172.31.28.68:9000)    |  100+ tables (events + dbt)      |
|  | Redis Cluster  (3 nodes, ports 6379+) |  Session + rate limit + cache    |
|  | RabbitMQ       (cluster)              |  11 exchanges, 26+ queues        |
|  +---------------------------------------+                                  |
|                                                                             |
|  +---------------------------------------+                                  |
|  | AWS (Serverless)                      |                                  |
|  |                                       |                                  |
|  | Lambda (ECR)    Fake follower ML      |  Auto-scaling                    |
|  | S3              Media + data pipeline |  gcc-social-data bucket          |
|  | CloudFront      CDN for media assets  |                                  |
|  | SQS             Batch job queuing     |  creator_follower_in             |
|  | Kinesis         Result streaming      |  creator_out (ON_DEMAND)         |
|  +---------------------------------------+                                  |
+=============================================================================+
```

### 6.2 GitLab CI/CD Pipeline

All 5 server-deployed projects follow the same pattern:

```
+=============================================================================+
|                           CI/CD PIPELINE (All Projects)                     |
+=============================================================================+
|                                                                             |
|  STAGES:  test --> build --> deploy_staging --> deploy_prod                  |
|                                                                             |
|  BUILD STAGE (Go projects):                                                 |
|    GOOS=linux GOARCH=amd64 go build                                         |
|    Artifacts: binary + .env.* + scripts/start.sh                            |
|    Expire: 1 week                                                           |
|                                                                             |
|  BUILD STAGE (Python projects):                                             |
|    python3.11 -m venv venv && pip install -r requirements.txt               |
|                                                                             |
|  DEPLOY STAGING:                                                            |
|    Trigger: automatic on master/dev push                                    |
|    Runner tags: deployer-stage-search                                       |
|    Script: /bin/bash scripts/start.sh                                       |
|                                                                             |
|  DEPLOY PRODUCTION:                                                         |
|    Trigger: MANUAL (requires click in GitLab UI)                            |
|    Runner tags: beat-deployer-1, beat-deployer-2, cb1-1, cb2-1              |
|    Script: /bin/bash scripts/start.sh                                       |
|    Rolling deploy: deploy to node 1, verify, then node 2                    |
|                                                                             |
|  GRACEFUL DEPLOYMENT SCRIPT (scripts/start.sh):                             |
|  +-----------------------------------------------------------------------+  |
|  |  1. ulimit -n 100000                    # Increase file descriptors    |  |
|  |  2. curl -XPUT /heartbeat/?beat=false   # Remove from load balancer   |  |
|  |  3. sleep 15                            # Drain in-flight requests     |  |
|  |  4. kill -9 $PID                        # Kill old process             |  |
|  |  5. sleep 10                            # Wait for cleanup             |  |
|  |  6. ENV=$ENV ./binary >> logs/out.log & # Start new process            |  |
|  |  7. sleep 20                            # Wait for startup             |  |
|  |  8. curl -XPUT /heartbeat/?beat=true    # Add back to load balancer   |  |
|  +-----------------------------------------------------------------------+  |
|                                                                             |
|  SPECIAL: Event-gRPC also publishes client libraries:                       |
|    publish_web:  protoc + npm publish (JavaScript/TypeScript)               |
|    publish_java: mvn clean install + mvn deploy (Java/Maven)               |
|                                                                             |
|  SPECIAL: Fake Follower uses Docker + ECR:                                  |
|    docker build -t fake-follower .                                          |
|    docker push {account}.dkr.ecr.ap-south-1.amazonaws.com/fake-follower    |
|    Update Lambda function to use new ECR image                              |
+=============================================================================+
```

---

## 7. SCALE NUMBERS ACROSS ALL PROJECTS

### Combined Platform Metrics

| Metric | Value | Source Project(s) |
|--------|-------|-------------------|
| **Total Lines of Code** | ~54,000+ | All 6 projects combined |
| **Daily Data Points Processed** | 10M+ | Beat + Event-gRPC |
| **Events Ingested Per Second** | 10K+ | Event-gRPC gRPC server |
| **Images Uploaded Per Day** | 8M+ | Beat asset pipeline |
| **Concurrent Scraping Workers** | 150+ | Beat worker pool |
| **RabbitMQ Consumer Workers** | 90+ | Event-gRPC (26 queues) |
| **Airflow DAGs** | 76 | Stir |
| **dbt Models** | 112 (29 staging + 83 marts) | Stir |
| **Scraping Flows** | 75+ | Beat |
| **API Endpoints** | 50+ | Coffee |
| **Microservices Proxied** | 13 | SaaS Gateway |
| **External API Integrations** | 15+ | Beat |
| **RabbitMQ Exchanges** | 11 | Event-gRPC |
| **RabbitMQ Queues** | 26+ | Event-gRPC |
| **ClickHouse Tables** | 100+ | Event-gRPC (18) + Stir (112 dbt) |
| **PostgreSQL Tables** | 27+ | Beat + Coffee |
| **Event Types** | 60+ | Event-gRPC protobuf |
| **Indian Name Database** | 35,183 names | Fake Follower |
| **Indic Scripts Supported** | 10 | Fake Follower |
| **Buffer Channel Capacity** | 100K+ combined | Event-gRPC |
| **Rate Limit Rules** | 7+ source-specific | Beat |
| **CORS Whitelisted Origins** | 41 | SaaS Gateway |
| **Middleware Layers (Gateway)** | 7 | SaaS Gateway |
| **Middleware Layers (Coffee)** | 10 | Coffee |
| **Cache: Ristretto Capacity** | 10M keys, 1GB | SaaS Gateway |
| **Cache: Redis Cluster Nodes** | 3-6 | Shared |
| **Production Server Nodes** | 6+ | All projects |
| **Python Dependencies (Stir)** | 350+ | Stir |
| **Go Dependencies** | 30+ per project | Coffee, Event-gRPC, Gateway |
| **Git Commits (Stir)** | 1,476 | Stir (most mature) |

### Performance Benchmarks

| Metric | Before | After | Improvement |
|--------|--------|-------|-------------|
| Analytics query latency | 30s (PostgreSQL) | 2s (ClickHouse) | 15x faster |
| Profile lookup latency | 200ms | 5ms (Redis cache) | 40x faster |
| Log retrieval times | -- | -- | 2.5x faster |
| API response times | -- | -- | 25% faster |
| Data freshness (core metrics) | 24 hours | 15 minutes | 96x fresher |
| Infrastructure costs | -- | -- | 30% reduction |
| Data latency | -- | -- | 50% reduction |
| Batch write overhead | 1 write/event | 1 write/1000 events | 99% I/O reduction |

---

## 8. HOW TO TALK ABOUT THIS IN INTERVIEWS

### 60-Second Combined System Pitch

> "At Good Creator Co., I single-handedly built the complete data infrastructure for a social media analytics SaaS platform -- six production microservices, about 54,000 lines of code across Python and Go.
>
> The system works like this: a Python-based scraping engine with 150+ concurrent workers crawls Instagram and YouTube profiles through 15+ API integrations with fallback strategies and Redis-backed rate limiting. Scraped data flows through RabbitMQ to a Go gRPC service that ingests 10K events per second, routing them through 26 consumer queues with buffered sinkers that batch-write to ClickHouse -- reducing database I/O by 99%.
>
> An Apache Airflow data platform runs 76 DAGs orchestrating 112 dbt models that transform raw data into analytics marts -- leaderboards, time-series, audience insights, genre rankings. These marts sync from ClickHouse through S3 to PostgreSQL using atomic table swaps for zero-downtime updates.
>
> A Go REST API with a 4-layer generic architecture serves 50+ endpoints with dual-database strategy -- PostgreSQL for transactional data, ClickHouse for analytics. An API gateway with two-layer caching (Ristretto + Redis) and JWT authentication proxies 13 microservices. And a serverless ML pipeline on AWS Lambda detects fake followers across 10 Indian languages.
>
> The result: analytics queries went from 30 seconds to 2 seconds, profile lookups from 200ms to 5ms, data freshness from 24 hours to 15 minutes, and infrastructure costs dropped 30%."

### Key Architecture Decisions to Highlight

**1. Why dual database (PostgreSQL + ClickHouse)?**
> "PostgreSQL excels at transactional CRUD for collections, profiles, and user data. But analytics queries -- aggregations across millions of rows, time-series analysis, ranking computations -- were killing PostgreSQL at 30+ seconds per query. ClickHouse's columnar storage and vectorized execution handles those same queries in 2 seconds. The trade-off: we needed a sync pipeline (Stir) to move computed results back to PostgreSQL for the API layer."

**2. Why RabbitMQ instead of Kafka?**
> "At our scale (10K events/sec, not 1M/sec), RabbitMQ gave us what we needed: durable queues, per-consumer acknowledgment, dead letter routing, and simple retry logic. The exchange-based routing let us fan out events to different sinkers without consumer coordination. Kafka would have added operational overhead we didn't need yet."

**3. Why SQL-based task queue instead of Celery?**
> "Beat uses PostgreSQL's FOR UPDATE SKIP LOCKED as a task queue. This gives us priority ordering, task expiration, status tracking, and query-based monitoring -- all in the same database where task results are stored. No separate infrastructure to manage. At 150 workers, the contention is manageable because SKIP LOCKED ensures workers never block each other."

**4. Why buffered sinkers instead of direct writes?**
> "Event-gRPC processes 10K events/second. Writing each event individually to ClickHouse would be 10K database round-trips per second. The buffered sinker pattern batches 1000 records and flushes every 5 seconds -- whichever comes first. This reduces I/O by 99% and ClickHouse handles bulk inserts much more efficiently than individual rows."

**5. Why ClickHouse --> S3 --> PostgreSQL sync pattern?**
> "ClickHouse has no native way to push data to PostgreSQL. We export to S3 using ClickHouse's built-in S3 function (INSERT INTO FUNCTION s3), download via SSH to the PostgreSQL server, COPY the JSON into a temp table, parse JSONB into typed columns, and do an atomic table swap (RENAME) for zero-downtime updates. The S3 intermediate step decouples the two databases and provides a backup."

---

*This master document maps the complete GCC platform architecture. All 6 projects were built by Anshul Garg as the sole backend engineer at Good Creator Co., totaling ~54K+ lines of production code in Python and Go.*

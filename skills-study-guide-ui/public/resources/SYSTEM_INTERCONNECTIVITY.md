# COMPLETE SYSTEM INTERCONNECTIVITY MAP
## How beat, event-grpc, stir, coffee, and ClickHouse Work Together

---

# VISUAL ARCHITECTURE

```
┌─────────────────────────────────────────────────────────────────────────────────────────────┐
│                              COMPLETE SYSTEM ARCHITECTURE                                    │
└─────────────────────────────────────────────────────────────────────────────────────────────┘

                                    ┌─────────────────┐
                                    │   EXTERNAL      │
                                    │   SCRAPING      │
                                    │   APIs          │
                                    │ (Instagram,     │
                                    │  YouTube, etc.) │
                                    └────────┬────────┘
                                             │
                                             ▼
┌─────────────────────────────────────────────────────────────────────────────────────────────┐
│                                       BEAT (Python)                                          │
│  ┌─────────────────────────────────────────────────────────────────────────────────────┐   │
│  │  • 150+ worker processes                                                              │   │
│  │  • Scrapes Instagram/YouTube profiles and posts                                       │   │
│  │  • Publishes events to RabbitMQ                                                       │   │
│  │  • Stores transactional data in PostgreSQL                                            │   │
│  │  • Rate limiting with Redis                                                           │   │
│  └─────────────────────────────────────────────────────────────────────────────────────┘   │
│                                             │                                               │
│         ┌───────────────────────────────────┼───────────────────────────────────┐          │
│         │                                   │                                   │          │
│         ▼                                   ▼                                   ▼          │
│  ┌─────────────┐                    ┌─────────────┐                    ┌─────────────┐    │
│  │ PostgreSQL  │                    │  RabbitMQ   │                    │    Redis    │    │
│  │ (beat DB)   │                    │  (beat.dx)  │                    │ (rate limit)│    │
│  └─────────────┘                    └──────┬──────┘                    └─────────────┘    │
└─────────────────────────────────────────────┼───────────────────────────────────────────────┘
                                              │
                    ┌─────────────────────────┼─────────────────────────┐
                    │                         │                         │
                    ▼                         ▼                         ▼
           ┌───────────────┐         ┌───────────────┐         ┌───────────────┐
           │post_log_events│         │profile_log_   │         │sentiment_log_ │
           │_q (20 workers)│         │events_q       │         │events_q       │
           └───────┬───────┘         └───────┬───────┘         └───────┬───────┘
                   │                         │                         │
                   └─────────────────────────┼─────────────────────────┘
                                             │
                                             ▼
┌─────────────────────────────────────────────────────────────────────────────────────────────┐
│                                    EVENT-GRPC (Go)                                           │
│  ┌─────────────────────────────────────────────────────────────────────────────────────┐   │
│  │  • 26 RabbitMQ consumer queues                                                        │   │
│  │  • Buffered sinkers (1000 events/batch, 5-sec flush)                                  │   │
│  │  • Writes to ClickHouse (_e.* tables)                                                 │   │
│  │  • 70+ concurrent workers                                                             │   │
│  └─────────────────────────────────────────────────────────────────────────────────────┘   │
│                                             │                                               │
│                                             ▼                                               │
│                                    ┌───────────────┐                                       │
│                                    │  ClickHouse   │                                       │
│                                    │ (_e database) │                                       │
│                                    │ 21+ tables    │                                       │
│                                    └───────────────┘                                       │
└─────────────────────────────────────────────────────────────────────────────────────────────┘
                                             │
                                             ▼
┌─────────────────────────────────────────────────────────────────────────────────────────────┐
│                                    STIR (Airflow + dbt)                                      │
│  ┌─────────────────────────────────────────────────────────────────────────────────────┐   │
│  │  • 76 Airflow DAGs                                                                    │   │
│  │  • 112 dbt models (29 staging + 83 marts)                                             │   │
│  │  • Reads from beat_replica + coffee_replica + _e (events)                             │   │
│  │  • Transforms in ClickHouse (dbt schema)                                              │   │
│  │  • Syncs to PostgreSQL via S3                                                         │   │
│  └─────────────────────────────────────────────────────────────────────────────────────┘   │
│         │                                                                    │              │
│         │ SOURCE                                                      SINK   │              │
│         ▼                                                                    ▼              │
│  ┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐    ┌─────────────┐     │
│  │ beat_replica    │    │ coffee_replica  │    │  S3 (staging)   │    │ PostgreSQL  │     │
│  │ (ClickHouse)    │    │ (ClickHouse)    │    │  JSON files     │    │ (beat DB)   │     │
│  └─────────────────┘    └─────────────────┘    └─────────────────┘    └─────────────┘     │
└─────────────────────────────────────────────────────────────────────────────────────────────┘
                                             │
                                             │ Publishes upsert events
                                             ▼
┌─────────────────────────────────────────────────────────────────────────────────────────────┐
│                                      COFFEE (Go)                                             │
│  ┌─────────────────────────────────────────────────────────────────────────────────────┐   │
│  │  • REST API for SaaS platform                                                         │   │
│  │  • Multi-tenant architecture                                                          │   │
│  │  • Calls Beat API for real-time data                                                  │   │
│  │  • Consumes upsert events from stir                                                   │   │
│  │  • Manages collections, discovery, campaigns                                          │   │
│  └─────────────────────────────────────────────────────────────────────────────────────┘   │
│         │                         │                         │                              │
│         ▼                         ▼                         ▼                              │
│  ┌─────────────┐          ┌─────────────┐          ┌─────────────┐                        │
│  │ PostgreSQL  │          │   Redis     │          │  RabbitMQ   │                        │
│  │ (coffee DB) │          │  (cache)    │          │ (events)    │                        │
│  └─────────────┘          └─────────────┘          └─────────────┘                        │
└─────────────────────────────────────────────────────────────────────────────────────────────┘
                                             │
                                             ▼
                                    ┌───────────────┐
                                    │   SaaS UI     │
                                    │   (Frontend)  │
                                    └───────────────┘
```

---

# DATA FLOWS

## FLOW 1: Profile/Post Data Ingestion

```
┌──────────────────────────────────────────────────────────────────────────────────────┐
│                        PROFILE/POST DATA INGESTION FLOW                               │
└──────────────────────────────────────────────────────────────────────────────────────┘

Step 1: BEAT scrapes from external APIs
─────────────────────────────────────────
Instagram Graph API ─┐
Instagram RapidAPI ──┼──► beat workers ──► PostgreSQL (instagram_account, instagram_post)
YouTube Data API ────┘                           │
                                                 │ make_scrape_log_event()
                                                 ▼
Step 2: BEAT publishes events to RabbitMQ
─────────────────────────────────────────
beat.dx exchange
├── profile_log_events (routing key)
├── post_log_events
├── post_activity_log_events
├── sentiment_log_events
├── profile_relationship_log_events
└── scrape_request_log_events

Step 3: EVENT-GRPC consumes and flushes to ClickHouse
─────────────────────────────────────────────────────
RabbitMQ queues                    ClickHouse tables
├── profile_log_events_q ─────────► _e.profile_log_events
├── post_log_events_q ────────────► _e.post_log_events
├── sentiment_log_events_q ───────► _e.sentiment_log_events
└── post_activity_log_events_q ───► _e.post_activity_log_events

Buffered sinker: 1000 events/batch OR 5-second flush

Step 4: STIR transforms data with dbt
─────────────────────────────────────
Sources (ClickHouse):
├── beat_replica.instagram_account
├── beat_replica.instagram_post
├── _e.profile_log_events
└── _e.post_log_events
        │
        │ dbt run --models tag:core
        ▼
Marts (ClickHouse dbt schema):
├── mart_instagram_account
├── mart_youtube_account
├── mart_leaderboard
├── mart_time_series
└── mart_genre_overview

Step 5: STIR syncs to PostgreSQL for COFFEE
───────────────────────────────────────────
ClickHouse (dbt.mart_leaderboard)
        │
        │ INSERT INTO FUNCTION s3(...)
        ▼
S3 (gcc-social-data/data-pipeline/tmp/leaderboard.json)
        │
        │ SSH download to pg server
        ▼
/tmp/leaderboard.json
        │
        │ COPY + atomic table swap
        ▼
PostgreSQL (leaderboard table)

Step 6: COFFEE serves data via REST API
───────────────────────────────────────
SaaS UI ──► coffee API ──► PostgreSQL (leaderboard, instagram_account, etc.)
```

---

## FLOW 2: Real-Time Profile Lookup

```
┌──────────────────────────────────────────────────────────────────────────────────────┐
│                        REAL-TIME PROFILE LOOKUP FLOW                                  │
└──────────────────────────────────────────────────────────────────────────────────────┘

User searches for "@virat.kohli" in SaaS UI
        │
        ▼
Coffee API: GET /discovery/instagram/byhandle/virat.kohli
        │
        │ Check PostgreSQL first
        ▼
┌─────────────────────────────────────────────────────┐
│ SELECT * FROM instagram_account WHERE handle = ?    │
└─────────────────────────────────────────────────────┘
        │
        │ If NOT FOUND:
        ▼
Coffee calls Beat API
        │
        │ GET http://beat.goodcreator.co/profiles/INSTAGRAM/byhandle/virat.kohli
        ▼
┌─────────────────────────────────────────────────────┐
│ Beat:                                               │
│ 1. Check rate limits (Redis)                        │
│ 2. Get credential (PostgreSQL)                      │
│ 3. Call Instagram Graph API                         │
│ 4. Parse response                                   │
│ 5. Save to PostgreSQL                               │
│ 6. Publish event to RabbitMQ                        │
│ 7. Return response to Coffee                        │
└─────────────────────────────────────────────────────┘
        │
        ▼
Coffee receives profile data
        │
        │ Transform and save
        ▼
┌─────────────────────────────────────────────────────┐
│ INSERT INTO instagram_account (...) VALUES (...)    │
│ + Create campaign_profile                           │
│ + Enrich with keywords, location                    │
└─────────────────────────────────────────────────────┘
        │
        │ Return to UI
        ▼
User sees profile details
```

---

## FLOW 3: Time-Series Analytics

```
┌──────────────────────────────────────────────────────────────────────────────────────┐
│                        TIME-SERIES ANALYTICS FLOW                                     │
└──────────────────────────────────────────────────────────────────────────────────────┘

Every profile crawl generates a snapshot
        │
        ▼
Beat: make_scrape_log_event("profile_log", {
    profile_id: "123",
    followers: 10000000,
    following: 500,
    timestamp: "2024-01-15 10:30:00"
})
        │
        │ RabbitMQ
        ▼
event-grpc: BufferProfileLogEvents()
        │
        │ Batch insert
        ▼
ClickHouse: _e.profile_log_events
┌─────────────────────────────────────────────────────────────────────────────────┐
│ profile_id │ followers │ following │ event_timestamp      │ insert_timestamp    │
├────────────┼───────────┼───────────┼──────────────────────┼─────────────────────┤
│ 123        │ 10000000  │ 500       │ 2024-01-15 10:30:00  │ 2024-01-15 10:30:05 │
│ 123        │ 10050000  │ 502       │ 2024-01-16 10:30:00  │ 2024-01-16 10:30:05 │
│ 123        │ 10100000  │ 505       │ 2024-01-17 10:30:00  │ 2024-01-17 10:30:05 │
└─────────────────────────────────────────────────────────────────────────────────┘
        │
        │ dbt model (stir)
        ▼
dbt.mart_time_series
┌─────────────────────────────────────────────────────┐
│ SELECT                                              │
│     profile_id,                                     │
│     toDate(event_timestamp) as date,                │
│     argMax(followers, event_timestamp) as followers │
│ FROM _e.profile_log_events                          │
│ GROUP BY profile_id, date                           │
└─────────────────────────────────────────────────────┘
        │
        │ Sync to PostgreSQL
        ▼
Coffee API: GET /analytics/timeseries/{profile_id}
        │
        │ Query PostgreSQL
        ▼
User sees follower growth chart
```

---

# DETAILED CONNECTION MAPS

## BEAT → Everything Else

```
┌─────────────────────────────────────────────────────────────────────────────────────────────┐
│                                    BEAT CONNECTIONS                                          │
└─────────────────────────────────────────────────────────────────────────────────────────────┘

BEAT PUBLISHES TO:
─────────────────
RabbitMQ Exchanges:
├── beat.dx (main exchange)
│   ├── profile_log_events ──────────────────────► event-grpc → ClickHouse
│   ├── post_log_events ─────────────────────────► event-grpc → ClickHouse
│   ├── post_activity_log_events ────────────────► event-grpc → ClickHouse
│   ├── sentiment_log_events ────────────────────► event-grpc → ClickHouse
│   ├── profile_relationship_log_events ─────────► event-grpc → ClickHouse
│   ├── scrape_request_log_events ───────────────► event-grpc → ClickHouse
│   ├── order_log_events ────────────────────────► event-grpc → ClickHouse
│   └── keyword_collection_rk ───────────────────► beat (internal job processing)
│
├── coffee.dx
│   ├── keyword_collection_report_completion ────► coffee (report ready notification)
│   └── sentiment_collection_report_out_rk ──────► coffee (sentiment report ready)
│
└── identity.dx
    ├── trace_log ───────────────────────────────► event-grpc → ClickHouse
    └── access_token_expired_rk ─────────────────► identity service


BEAT CONSUMES FROM:
──────────────────
RabbitMQ:
├── identity.dx / new_access_token_rk ───────────► Update credentials
├── beat.dx / credentials_validate_rk ───────────► Validate tokens
├── beat.dx / keyword_collection_rk ─────────────► Generate keyword reports
├── beat.dx / post_activity_log_bulk ────────────► Sentiment extraction
└── beat.dx / sentiment_collection_report_in_rk ─► Generate sentiment reports


BEAT READS FROM:
───────────────
PostgreSQL (beat database):
├── credential
├── scrape_request_log (task queue)
├── instagram_account
├── instagram_post
├── youtube_account
└── youtube_post

ClickHouse (dbt schema):
├── dbt.stg_coffee_post_collection_item
├── dbt.stg_coffee_post_collection
└── dbt.mart_genre_overview


BEAT WRITES TO:
──────────────
PostgreSQL (beat database):
├── All tables above (upserts)
└── profile_log, post_log (audit tables)

S3 (gcc-social-data bucket):
├── keyword_collections/{date}/{job_id}.parquet
├── sentiment_reports/{date}/{job_id}.parquet
└── assets/{entity_type}/{platform}/{id}.jpg


BEAT CALLS APIs:
───────────────
├── Instagram Graph API (graph.facebook.com)
├── YouTube Data API (googleapis.com)
├── RapidAPI (multiple providers)
├── Identity Service (identityservice.bulbul.tv)
├── RAY ML Service (ray.goodcreator.co)
│   ├── CATEGORIZER model
│   └── SENTIMENT model
└── Azure OpenAI (gcc-openai.openai.azure.com)
```

---

## EVENT-GRPC → Everything Else

```
┌─────────────────────────────────────────────────────────────────────────────────────────────┐
│                                  EVENT-GRPC CONNECTIONS                                      │
└─────────────────────────────────────────────────────────────────────────────────────────────┘

EVENT-GRPC CONSUMES FROM:
────────────────────────
RabbitMQ (26 queues total):

From beat.dx:
├── post_log_events_q (20 workers) ────────────► PostLogEventsSinker
├── profile_log_events_q (2 workers) ──────────► ProfileLogEventsSinker
├── sentiment_log_events_q (2 workers) ────────► SentimentLogEventsSinker
├── post_activity_log_events_q (2 workers) ────► PostActivityLogEventsSinker
├── profile_relationship_log_events_q (2) ─────► ProfileRelationshipLogEventsSinker
├── scrape_request_log_events_q (2 workers) ───► ScrapeRequestLogEventsSinker
└── order_log_events_q (2 workers) ────────────► OrderLogEventsSinker

From identity.dx:
└── trace_log ─────────────────────────────────► TraceLogEventsSinker

From coffee.dx:
└── activity_tracker_q ────────────────────────► PartnerActivityLogEventsSinker

From other exchanges:
├── grpc_event.tx / grpc_clickhouse_event_q ───► SinkEventToClickhouse
├── ab.dx / ab_assignments ────────────────────► SinkABAssignmentsToClickhouse
├── branch_event.tx / branch_event_q ──────────► SinkBranchEventToClickhouse
├── webengage_event.dx / webengage_ch_event_q ─► SinkWebengageEventToClickhouse
└── shopify_event.dx / shopify_events_q ───────► BufferShopifyEvents


EVENT-GRPC WRITES TO:
────────────────────
ClickHouse (_e database - 21+ tables):
├── profile_log_events
├── post_log_events
├── sentiment_log_events
├── post_activity_log_events
├── profile_relationship_log_events
├── scrape_request_log_events
├── order_log_events
├── trace_log_events
├── partner_activity_log_events
├── event
├── error_event
├── ab_assignment
├── branch_event
├── affiliate_order_event
├── bigboss_vote_log
├── shopify_event
└── webengage_event


BUFFERED SINKER PATTERN:
───────────────────────
┌─────────────────────────────────────────────────────┐
│ Channel buffer: 10,000 messages                     │
│ Batch size: 1,000 events                            │
│ Flush interval: 5 seconds                           │
│ Flush condition: batch full OR timer tick           │
└─────────────────────────────────────────────────────┘
```

---

## STIR → Everything Else

```
┌─────────────────────────────────────────────────────────────────────────────────────────────┐
│                                    STIR CONNECTIONS                                          │
└─────────────────────────────────────────────────────────────────────────────────────────────┘

STIR READS FROM (ClickHouse):
────────────────────────────
beat_replica schema (PostgreSQL replica):
├── instagram_account
├── instagram_post_simple
├── youtube_account
├── youtube_post_simple
├── profile_log
├── post_activity_log
├── credential
├── scrape_request_log
├── asset_log
├── instagram_profile_insights
└── youtube_profile_insights

coffee_replica schema (PostgreSQL replica):
├── post_collection
├── post_collection_item
├── profile_collection
├── profile_collection_item
├── keyword_collection
├── collection_group
├── activity_tracker
├── view_instagram_account_lite
└── view_youtube_account_lite

_e schema (event-grpc writes):
├── profile_log_events
├── post_log_events
├── sentiment_log_events
└── post_activity_log_events


STIR TRANSFORMS IN (ClickHouse dbt schema):
──────────────────────────────────────────
Staging models (29):
├── stg_beat_instagram_account
├── stg_beat_instagram_post
├── stg_beat_youtube_account
├── stg_beat_profile_log
├── stg_coffee_post_collection
├── stg_coffee_post_collection_item
└── ...

Mart models (83):
├── mart_instagram_account
├── mart_youtube_account
├── mart_leaderboard
├── mart_time_series
├── mart_genre_overview
├── mart_trending_content
├── mart_collection_post
├── mart_instagram_hashtags
└── ...


STIR SYNCS TO (PostgreSQL beat database):
────────────────────────────────────────
Three-layer sync pattern:
ClickHouse dbt.mart_* → S3 JSON → PostgreSQL

Target tables:
├── leaderboard
├── time_series
├── genre_overview
├── trending_content
├── collection_post_metrics_summary
├── collection_post_metrics_ts
├── social_profile_hashtags
├── collection_hashtag
├── collection_keyword
└── group_metrics


STIR PUBLISHES TO (RabbitMQ):
────────────────────────────
upserttracker.dx exchange:
├── upsert_instagram_account_rk ───► coffee (profile updates)
└── upsert_youtube_account_rk ─────► coffee (profile updates)
```

---

## COFFEE → Everything Else

```
┌─────────────────────────────────────────────────────────────────────────────────────────────┐
│                                   COFFEE CONNECTIONS                                         │
└─────────────────────────────────────────────────────────────────────────────────────────────┘

COFFEE CALLS APIs:
─────────────────
Beat API (http://beat.goodcreator.co):
├── GET /profiles/{platform}/byhandle/{handle}
├── GET /profiles/{platform}/byid/{id}
├── GET /recent/posts/{platform}/byprofileid/{id}
├── GET /profiles/INSTAGRAM/byhandle/{handle}/insights
├── GET /profiles/INSTAGRAM/byhandle/{handle}/audienceinsights
└── GET /youtube/channel/byhandle/{handle}

Other services:
├── JobTracker (jobtrackerservice.bulbul.tv) - Async job management
├── DAM (damservice.bulbul.tv) - Digital asset management
├── Partner Service (productservice.bulbul.tv) - Partner contracts
└── Identity Service - Authentication


COFFEE READS FROM:
─────────────────
PostgreSQL (coffee database):
├── instagram_account
├── youtube_account
├── campaign_profiles
├── profile_collection
├── profile_collection_item
├── post_collection
├── post_collection_item
├── keyword_collection
├── collection_analytics
└── activity_tracker

ClickHouse (dbt schema):
└── Analytics queries for time-series data


COFFEE CONSUMES FROM (RabbitMQ):
───────────────────────────────
upserttracker.dx (from stir):
├── upsert_instagram_account_q ──► PerformUpsertInstagramAccount()
└── upsert_youtube_account_q ────► PerformUpsertYoutubeAccount()

jobtracker.dx:
├── duplicate_collection_q
├── download_collection_q
├── import_from_profile_collection_q
├── add_item_profile_collection_q
└── add_item_post_collection_q


COFFEE PUBLISHES TO (RabbitMQ):
─────────────────────────────
beat.dx:
└── keyword_collection_rk ──────► beat (trigger keyword report)

coffee.dx:
└── activity_tracker_rk ────────► event-grpc → ClickHouse


COFFEE CACHES IN (Redis):
────────────────────────
├── partnercontract-{partnerId} (12h TTL)
└── Session data
```

---

# DATABASE CONNECTIONS SUMMARY

```
┌─────────────────────────────────────────────────────────────────────────────────────────────┐
│                              DATABASE CONNECTIONS MAP                                        │
└─────────────────────────────────────────────────────────────────────────────────────────────┘

PostgreSQL (172.31.2.21:5432)
├── beat database
│   ├── Written by: beat, stir (sync)
│   └── Read by: beat, coffee
│
└── coffee database
    ├── Written by: coffee
    └── Read by: coffee


ClickHouse (172.31.28.68:9000)
├── _e database (events)
│   ├── Written by: event-grpc
│   └── Read by: stir
│
├── beat_replica schema
│   ├── Written by: ClickHouse replication
│   └── Read by: stir (dbt sources)
│
├── coffee_replica schema
│   ├── Written by: ClickHouse replication
│   └── Read by: stir (dbt sources)
│
└── dbt schema (transformations)
    ├── Written by: stir (dbt run)
    └── Read by: stir (for sync), beat (for reports)


Redis Cluster
├── beat: Rate limiting, credential state
└── coffee: Partner contract cache


S3 (gcc-social-data bucket)
├── Written by: beat (reports, assets), stir (sync files)
└── Read by: stir (sync to PostgreSQL), CDN (assets)
```

---

# RABBITMQ EXCHANGE/QUEUE MAP

```
┌─────────────────────────────────────────────────────────────────────────────────────────────┐
│                                   RABBITMQ TOPOLOGY                                          │
└─────────────────────────────────────────────────────────────────────────────────────────────┘

beat.dx (exchange)
├── profile_log_events ────────► profile_log_events_q ────► event-grpc
├── post_log_events ───────────► post_log_events_q ────────► event-grpc
├── sentiment_log_events ──────► sentiment_log_events_q ───► event-grpc
├── post_activity_log_events ──► post_activity_log_events_q ► event-grpc
├── scrape_request_log_events ─► scrape_request_log_events_q ► event-grpc
├── keyword_collection_rk ─────► keyword_collection_q ─────► beat
└── credentials_validate_rk ───► credentials_validate_q ───► beat

identity.dx (exchange)
├── trace_log ─────────────────► trace_log ────────────────► event-grpc
├── new_access_token_rk ───────► identity_token_q ─────────► beat
└── access_token_expired_rk ───► identity service

coffee.dx (exchange)
├── activity_tracker_rk ───────► activity_tracker_q ───────► event-grpc
└── keyword_collection_report_completion ──────────────────► coffee

upserttracker.dx (exchange)
├── upsert_instagram_account ──► upsert_instagram_account_q ► coffee
└── upsert_youtube_account ────► upsert_youtube_account_q ──► coffee
```

---

# INTERVIEW EXPLANATION

**"Explain how your systems work together"**

> "We built a microservices architecture for social media analytics:
>
> **The Data Flow:**
>
> 1. **beat** (Python) scrapes Instagram/YouTube using 150+ workers with rate limiting
>
> 2. Instead of direct database writes for time-series data, beat publishes events to **RabbitMQ**
>
> 3. **event-grpc** (Go) consumes these events with buffered sinkers (1000 events/batch) and writes to **ClickHouse**
>
> 4. **stir** (Airflow + dbt) transforms the raw data in ClickHouse into analytics-ready marts, then syncs to PostgreSQL via S3
>
> 5. **coffee** (Go) serves the REST API, reading from PostgreSQL for transactional queries and calling beat for real-time lookups
>
> **Why this architecture?**
>
> - **Separation of concerns**: OLTP (PostgreSQL) vs OLAP (ClickHouse)
> - **Event-driven**: Decouples producers from consumers
> - **Scalability**: Each component can scale independently
> - **Reliability**: Buffered sinkers prevent data loss
> - **Performance**: ClickHouse for fast analytics, PostgreSQL for transactional consistency"

---

*This document shows the complete interconnectivity between all 5 systems in your work experience.*

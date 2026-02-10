# GCC (Good Creator Co.) + PayU - Interview Prep

> **You have ACTUAL SOURCE CODE for 6 GCC projects in google-interview-prep/**
> This file maps resume bullets to real code, architecture, and interview answers.

---

# GOOD CREATOR CO. - THE BIG PICTURE

## What Was GCC?

A SaaS platform for **influencer marketing analytics**. Brands (like Nike, Pepsi) use it to discover influencers, track their performance, and manage campaigns.

## Your 6 Projects (ACTUAL CODE EXISTS)

```
┌──────────────┐    ┌──────────────┐    ┌──────────────┐    ┌──────────────┐
│    BEAT      │    │  EVENT-GRPC  │    │    STIR      │    │   COFFEE    │
│  (Python)    │    │    (Go)      │    │ (Airflow+dbt)│    │    (Go)     │
│              │    │              │    │              │    │              │
│ Scrape social│───►│ Ingest events│───►│ Transform &  │───►│ Serve via   │
│ media data   │    │ to ClickHouse│    │ build marts  │    │ REST API    │
│ 75+ flows    │    │ 10K events/s │    │ 76 DAGs      │    │ 40+ endpoints│
│ 15K+ LOC     │    │ 10K+ LOC     │    │ 17.5K+ LOC   │    │ 8.5K+ LOC   │
└──────────────┘    └──────────────┘    └──────────────┘    └──────────────┘

+ saas-gateway (Go API Gateway, 13 microservices)
+ fake_follower_analysis (Python ML, AWS Lambda)
```

**Total: ~60,000+ lines of code across 6 projects.**

---

# RESUME BULLET → CODE MAPPING

## Bullet 1: "Optimized API response times by 25% and reduced operational costs by 30%"

**Which Project:** Coffee (Go REST API) + Stir (Data Platform)

**What You Actually Did:**
- **Coffee:** Added Redis caching for frequently queried influencer profiles. ClickHouse for analytics queries instead of PostgreSQL.
- **Stir:** Built incremental dbt models instead of full refreshes. Reduced query costs by only processing new data.

**Architecture Decision:**
```
BEFORE: All data in PostgreSQL (row-based, slow for analytics)
AFTER:  PostgreSQL (transactional) + ClickHouse (analytics) + Redis (cache)

Result: Analytics queries: 30s → 2s (ClickHouse columnar)
        Profile lookups: 200ms → 5ms (Redis cache)
        Infra cost: 30% reduction (ClickHouse compression + incremental processing)
```

**Interview Answer:**
> "I implemented a dual-database strategy. PostgreSQL for transactional data, ClickHouse for analytics. Analytics queries that took 30 seconds in PostgreSQL ran in 2 seconds in ClickHouse because of columnar storage. Added Redis caching for hot data - profile lookups went from 200ms to 5ms. On the data platform side, switched from full-refresh dbt models to incremental, cutting compute costs by 30%."

---

## Bullet 2: "Designed an asynchronous data processing system handling 10M+ daily data points"

**Which Project:** Beat (Python) + Event-gRPC (Go)

**What You Actually Did:**
- **Beat:** Worker pool system with 150+ concurrent workers processing social media scraping tasks. FastAPI + uvloop for async I/O. Semaphore-based concurrency control. 75+ scraping flows.
- **Event-gRPC:** gRPC server ingesting events at 10K/sec. 26 RabbitMQ consumer queues with 90+ concurrent workers. Buffered sinkers batch-writing to ClickHouse.

**Architecture:**
```
Beat (Python):
├── 75+ scraping flows (Instagram, YouTube, Shopify)
├── Worker pool: 150+ concurrent workers
├── Rate limiting: Redis-backed distributed limiter
├── Output: RabbitMQ events

Event-gRPC (Go):
├── gRPC server: 10K+ events/sec ingestion
├── 26 RabbitMQ consumer queues
├── 90+ concurrent workers
├── Buffered sinkers: batch write to ClickHouse (1000 records/flush)
└── Multi-destination: ClickHouse, WebEngage, Shopify, Branch
```

**Interview Answer:**
> "I built a distributed data processing pipeline. The scraping service (Python) had 150+ concurrent workers crawling Instagram and YouTube profiles using async I/O. Events were published to RabbitMQ. The event ingestion service (Go) consumed from 26 RabbitMQ queues with 90+ workers and used buffered sinkers to batch-write to ClickHouse at 10K events per second. The buffering was key - instead of one DB write per event, we batched 1000 records per flush, reducing I/O overhead by 99%."

**Technical Deep Dive (If They Ask):**
> "The worker pool used semaphore-based concurrency control. Each scraping flow had a configurable concurrency limit. For Instagram, we limited to 50 concurrent requests because of rate limiting. For YouTube (using official API), we could go to 100. Redis-backed distributed rate limiter ensured we didn't exceed API quotas even across multiple instances."

---

## Bullet 3: "Built high-performance logging system with RabbitMQ, transitioning to ClickHouse, 2.5x faster retrieval"

**Which Project:** Event-gRPC (Go) + Stir (Airflow)

**What You Actually Did:**
```
BEFORE:
- All event logs stored in PostgreSQL
- 10M+ logs/day overwhelmed PostgreSQL
- Write latency: 5ms → 500ms (100x degradation)
- Analytics queries: 30+ seconds
- Storage bloat: billions of rows

AFTER:
- Events flow through RabbitMQ → Event-gRPC → ClickHouse
- Buffered writes: 1000 records per batch
- ClickHouse columnar storage: 5x compression
- Analytics queries: 30s → 12s → finally ~5s (2.5x improvement)
- Stir (Airflow + dbt) transforms raw logs into analytics marts
```

**Interview Answer:**
> "PostgreSQL was choking on 10M+ daily log writes - write latency went from 5ms to 500ms. I designed an event-driven architecture: scraping service publishes to RabbitMQ, Go-based event-gRPC service consumes and batch-writes to ClickHouse. ClickHouse's columnar storage gives 5x compression and much faster aggregation queries. The transition reduced log retrieval times by 2.5x. I used a dual-write strategy during migration - writing to both PostgreSQL and ClickHouse simultaneously, then switching reads once ClickHouse was validated."

---

## Bullet 4: "ETL data pipelines (Apache Airflow) cutting data latency by 50%"

**Which Project:** Stir (76 DAGs, 112 dbt models)

**What You Actually Did:**
```
Stir Data Platform:
├── 76 Airflow DAGs
│   ├── 11 dbt orchestration DAGs (core, hourly, daily, weekly)
│   ├── 17 Instagram sync DAGs
│   ├── 12 YouTube sync DAGs
│   ├── 15 Collection/Leaderboard sync DAGs
│   └── 21 Other sync DAGs
├── 112 dbt models
│   ├── 29 staging models (raw → cleaned)
│   └── 83 mart models (cleaned → business-ready)
└── ClickHouse → S3 → PostgreSQL sync pipeline
```

**Interview Answer:**
> "I built the entire data platform on Apache Airflow with dbt. 76 DAGs orchestrating everything from social media data ingestion to analytics mart generation. The key innovation was frequency-based scheduling: critical metrics run every 15 minutes (dbt_core), hourly transforms every 2 hours, and heavy aggregations daily. Before this, all processing was batch-daily, so metrics were 24 hours stale. After: core metrics refresh in 15 minutes. That's the 50% latency reduction - we cut average data freshness from 24 hours to under 1 hour."

---

## Bullet 5: "AWS S3-based asset upload system processing 8M images daily"

**Which Project:** Beat (Python - main_assets.py)

**What You Actually Did:**
```
Asset Upload Pipeline:
├── Beat crawls influencer profiles + posts
├── Extracts image/video URLs
├── Worker pool downloads media
├── Uploads to AWS S3 (CDN-backed)
├── Updates database with CDN URLs
└── 8M+ images/day across Instagram + YouTube
```

**Interview Answer:**
> "When we crawled influencer profiles, we needed to store profile pictures, post thumbnails, and video previews. I built an asset processing pipeline: dedicated worker pool in Beat downloads media from social platforms, uploads to S3 with lifecycle policies (move to Glacier after 90 days), and updates the database with CDN URLs. At 8 million images daily, the key optimization was batch processing - downloading and uploading in parallel with configurable concurrency, and using S3 multipart upload for large files."

---

## Bullet 6: "Real-time social media insights modules"

**Which Project:** Coffee (Go) - Genre Insights + Keyword Analytics modules

**What You Actually Did:**
```
Coffee modules:
├── genreinsights/        # Genre-based creator analytics
├── keywordcollection/    # Keyword tracking
├── leaderboard/          # Creator rankings
├── collectionanalytics/  # Collection performance
├── discovery/            # Influencer search
│   ├── searchmanager.go
│   ├── timeseriesmanager.go
│   ├── hashtagsmanager.go
│   └── audiencemanager.go
└── partnerusage/         # Usage tracking for billing
```

**Interview Answer:**
> "I built the analytics modules in our Go SaaS platform. Genre Insights let brands see top creators by category (fashion, tech, food). Keyword Analytics tracked trending hashtags and topics. Both were powered by ClickHouse time-series queries aggregated by dbt marts. The 10% engagement growth came from brands discovering relevant creators they wouldn't have found manually."

---

## Bullet 7: "Automated content filtering and 50% faster processing"

**Which Project:** Beat + Fake Follower Analysis (ML)

**What You Actually Did:**
```
Fake Follower Detection:
├── Ensemble ML model with 5 features
│   ├── Username pattern analysis (bot-like patterns)
│   ├── Profile completeness score
│   ├── Engagement ratio anomaly detection
│   ├── Posting frequency analysis
│   └── Follower/following ratio
├── Multi-language NLP (10 Indic scripts)
├── AWS Lambda + SQS + Kinesis serverless pipeline
└── RapidFuzz fuzzy matching for name detection
```

**Interview Answer:**
> "I built a fake follower detection system. ML ensemble model analyzing 5 features: username patterns, profile completeness, engagement ratios, posting frequency, and follower/following ratio. Deployed on AWS Lambda for serverless scaling. Used RapidFuzz for fuzzy name matching across 35K name database and transliteration for 10 Indic scripts. The 50% processing speed improvement came from parallelizing the detection across Lambda functions instead of sequential processing."

---

# PAYU - INTERVIEW PREP

## What Was PayU?

India's leading digital payments company. API Lending division - providing loans through partner banks via APIs.

## Bullet 1: "Scaled business operations by 40%, reduced failure rate 93%"

**Interview Answer:**
> "After loan disbursal, users wanted to access additional loans from other sources. I integrated new lending partner APIs into the platform, enabling cross-lending. This opened 40% more business volume. The 93% failure reduction (4.6% → 0.3%) came from fixing race conditions in the disbursal flow and adding proper retry logic with idempotency keys for partner API calls."

## Bullet 2: "Decreased TAT by 66% (3.2 min → 1.1 min)"

**Interview Answer:**
> "The loan disbursal journey had multiple synchronous API calls in sequence. I profiled each call, identified the slowest ones (KYC verification, bank account validation), and parallelized independent calls using CompletableFuture. Also added caching for repeated KYC lookups. The combination of parallelization and caching reduced TAT from 3.2 minutes to 1.1 minutes."

## Intern Bullets: "83% test coverage, SonarQube, Flyway"

**Interview Answer:**
> "I refactored the Loan Origination System to be testable. The existing code had tight coupling - business logic in controllers, database calls mixed with validation. I extracted services, added interfaces for dependency injection, and wrote comprehensive JUnit tests. Coverage went from 30% to 83%. Then I integrated SonarQube into GitHub Actions for automated quality checks and Flyway for database migration versioning - deployment errors dropped 90% because schema changes were tracked and reversible."

---

# HOW TO SPEAK ABOUT GCC + PAYU

## 60-Second GCC Pitch

> "At Good Creator Co., I built the data infrastructure for a social media analytics SaaS platform. Six production microservices: a Python-based scraping engine with 150+ concurrent workers crawling Instagram and YouTube, a Go gRPC service ingesting 10K events per second into ClickHouse, an Apache Airflow data platform with 76 DAGs and 112 dbt models, a Go REST API serving 40+ endpoints, an API gateway proxying 13 microservices, and an ML-based fake follower detection system on AWS Lambda. The platform processed 10 million data points daily with 8 million image uploads."

## 30-Second PayU Pitch

> "At PayU's API Lending division, I optimized the loan disbursal journey. Reduced failure rates by 93% through partner API integration fixes and retry logic. Cut disbursal time by 66% by parallelizing independent API calls. As an intern, I brought test coverage from 30% to 83% and implemented Flyway DB migrations that reduced deployment errors by 90%."

## Pivot Technique: GCC → Walmart

> "At GCC I built data pipelines and event processing systems. At Walmart, I applied similar patterns at a different scale - Kafka instead of RabbitMQ, BigQuery instead of ClickHouse, but the same principles: async processing, schema enforcement, and decoupled architectures."

---

# KEY NUMBERS TO REMEMBER

## GCC Numbers
| Metric | Value |
|--------|-------|
| Daily data points | 10M+ |
| Events/sec (gRPC) | 10K+ |
| Images/day (S3) | 8M |
| Airflow DAGs | 76 |
| dbt models | 112 |
| Concurrent workers (Beat) | 150+ |
| RabbitMQ queues | 26 |
| Log retrieval improvement | 2.5x |
| API response improvement | 25% |
| Cost reduction | 30% |
| Data latency reduction | 50% |
| Total LOC across 6 projects | ~60,000+ |

## PayU Numbers
| Metric | Value |
|--------|-------|
| Failure rate reduction | 93% (4.6% → 0.3%) |
| TAT reduction | 66% (3.2 min → 1.1 min) |
| Business scaling | 40% |
| Test coverage | 30% → 83% |
| Deployment error reduction | 90% |
| Release frequency | 2x improvement |

---

*GCC and PayU are SUPPORTING stories. Lead with Walmart, use these when asked about previous experience or to show range.*

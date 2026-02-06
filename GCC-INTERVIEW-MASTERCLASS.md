# GCC INTERVIEW MASTERCLASS
## The Ultimate Interview Guide for Good Creator Co. Work

> **This guide synthesizes insights from ALL 6 GCC projects (Beat, Event-gRPC, Stir, Coffee, SaaS Gateway, Fake Follower Analysis) into a unified, interview-ready strategy. You have actual source code for every project in google-interview-prep/.**

---

# 1. THE BIG PICTURE - What Was GCC?

## Company Context

Good Creator Co. (GCC) was a **SaaS platform for influencer marketing analytics**. Brands like Nike, Pepsi, and other enterprise clients used it to discover influencers, track their social media performance, and manage marketing campaigns. The platform aggregated data from Instagram, YouTube, and Shopify, then surfaced real-time analytics through a multi-tenant web application at goodcreator.co.

## Your Role

**Software Engineer I** (SWE-I), Feb 2023 - May 2024. You were one of the **core backend engineers** on a small startup team. You built and maintained 6 production microservices spanning two languages (Go and Python), an ML pipeline, and an enterprise data platform. This was a startup environment -- you wore many hats, owned services end-to-end, and shipped to production frequently.

## The 6 Projects - At a Glance

```
+------------------+    +------------------+    +------------------+    +------------------+
|      BEAT        |    |   EVENT-GRPC     |    |      STIR        |    |     COFFEE       |
|    (Python)      |    |      (Go)        |    | (Airflow + dbt)  |    |      (Go)        |
|                  |    |                  |    |                  |    |                  |
| Scrape social    |--->| Ingest events    |--->| Transform &      |--->| Serve via        |
| media data       |    | to ClickHouse    |    | build marts      |    | REST API         |
| 75+ flows        |    | 10K events/s     |    | 76 DAGs          |    | 40+ endpoints    |
| 15K+ LOC         |    | 10K+ LOC         |    | 17.5K+ LOC       |    | 8.5K+ LOC        |
+------------------+    +------------------+    +------------------+    +------------------+

                  + SaaS Gateway (Go, 13 microservices proxied)
                  + Fake Follower Analysis (Python ML, AWS Lambda)
```

**Total: approximately 60,000+ lines of production code across 6 projects.**

## System Architecture Summary

```
External APIs (Instagram Graph API, YouTube Data API, RapidAPIs)
        |
        v
   BEAT (Python) -- scrapes, processes, publishes to RabbitMQ
        |
        +---> PostgreSQL (transactional data: profiles, posts, credentials)
        +---> RabbitMQ (event publishing: beat.dx exchange)
        +---> S3 (media assets: 8M images/day)
        |
        v
   EVENT-GRPC (Go) -- consumes from 26 RabbitMQ queues, batch-writes
        |
        +---> ClickHouse (_e database: 21+ event tables)
        +---> External APIs (WebEngage, Branch.io, Shopify)
        |
        v
   STIR (Airflow + dbt) -- transforms raw data into analytics marts
        |
        +---> ClickHouse (dbt schema: 29 staging + 83 mart models)
        +---> S3 (staging for sync)
        +---> PostgreSQL (synced marts for Coffee)
        |
        v
   COFFEE (Go) -- serves analytics to frontend via REST API
        |
        +---> PostgreSQL (reads synced marts)
        +---> ClickHouse (direct analytics queries)
        +---> Redis (caching)
        |
        v
   SAAS GATEWAY (Go) -- API gateway, auth, caching, reverse proxy
        |
        +---> 13 downstream microservices
        +---> Identity Service (JWT validation)
        +---> Ristretto (in-memory) + Redis (distributed) caching
        |
        v
   SaaS UI (goodcreator.co frontend)

   FAKE FOLLOWER ANALYSIS (Python ML) -- runs separately on AWS Lambda
        +---> SQS/Kinesis pipeline
        +---> 10 Indic language transliteration models
```

---

# 2. YOUR 60-SECOND GCC PITCH

> "At Good Creator Co., I built the backend data infrastructure for a social media analytics SaaS platform -- think of it as a tool that helps brands discover the right influencers and measure their campaigns.
>
> I owned six production microservices end-to-end. First, a Python-based scraping engine called Beat with 150+ concurrent workers crawling Instagram and YouTube data using 15+ API integrations, with 3-level stacked rate limiting. That published events to RabbitMQ, which were consumed by Event-gRPC -- a Go service I built that ingested 10K events per second using buffered sinkers that batch-wrote to ClickHouse, reducing per-event I/O by 99%.
>
> Then I built the data platform, Stir, on Apache Airflow with 76 DAGs and 112 dbt models transforming raw data into analytics marts. This cut our data freshness from 24 hours to under 1 hour -- a 50% latency reduction.
>
> On the serving side, I built Coffee, a Go REST API with 40+ endpoints using a 4-layer architecture with dual-database strategy -- PostgreSQL for transactional data and ClickHouse for analytics. I also built the API gateway that proxied 13 microservices with JWT auth, two-tier caching with Ristretto and Redis, and a 7-layer middleware pipeline.
>
> Finally, I built a fake follower detection system using ML -- an ensemble model with HMM transliteration for 10 Indic scripts, deployed on AWS Lambda. All together, the platform processed 10 million data points daily."

---

# 3. YOUR 30-SECOND GCC PITCH

> "Before Walmart, I was at Good Creator Co., a social media analytics startup. I built six production microservices: a Python scraping engine, a Go event ingestion service handling 10K events per second into ClickHouse, an Airflow data platform with 76 DAGs and 112 dbt models, a Go REST API, an API gateway for 13 services, and an ML-based fake follower detection system. The platform processed 10 million data points daily and served enterprise brands for influencer marketing analytics. About 60K lines of production code across Go and Python."

---

# 4. STAR STORIES FROM GCC

---

## Story 1: "Building the Data Pipeline" (Beat -> Event-gRPC -> ClickHouse)

**When to use:** System design questions, scaling questions, "Tell me about a complex system you built," "How did you handle high throughput?"

### Situation
Our SaaS platform was growing fast, and PostgreSQL was choking on 10 million+ daily event logs. Write latency degraded from 5ms to 500ms -- a 100x spike. Analytics queries over billions of rows in PostgreSQL were taking 30+ seconds. The database was becoming the bottleneck for the entire platform, and clients were seeing stale data.

### Task
I needed to design an event-driven architecture that could handle 10K events per second ingestion, provide fast analytical queries, and decouple our scraping service from the analytics database. The key constraint: zero data loss during the transition, and we had to keep the existing PostgreSQL-based services running in parallel.

### Action
I designed a three-tier pipeline:

1. **Beat (Python scraper)** publishes events to RabbitMQ instead of writing directly to PostgreSQL. Each profile/post scrape generates a log event (profile_log, post_log, sentiment_log, etc.) published to the `beat.dx` exchange with routing keys.

2. **Event-gRPC (Go)** consumes from 26 RabbitMQ queues with 90+ concurrent workers. The critical innovation was the **buffered sinker pattern**: instead of one DB write per event, I created Go channels with 10,000 message capacity. A separate goroutine batches these into groups of 1,000 or flushes every 5 seconds (whichever comes first), then does a single batch INSERT into ClickHouse. This reduced I/O overhead by 99%.

3. I implemented a **dual-write migration strategy** -- writing to both PostgreSQL and ClickHouse simultaneously during the transition, then switching reads to ClickHouse once validated.

The gRPC service defined 60+ event types in Protocol Buffers, grouped them by type before publishing to RabbitMQ, and handled timestamp correction (rejecting future timestamps). Each consumer had retry logic with dead letter routing (max 2 retries before moving to error queue).

### Result
- Log retrieval queries went from 30s to approximately 12s, a **2.5x improvement**
- ClickHouse columnar storage provided **5x compression** over PostgreSQL row-based storage
- Event ingestion throughput: **10K+ events/sec** sustained
- The buffered sinker pattern produced a **33x write throughput improvement** over individual inserts
- Zero data loss during migration

**Code evidence:** `event-grpc/sinker/eventsinker.go` (buffered sinker), `event-grpc/main.go` (26 consumer configs), `event-grpc/proto/eventservice.proto` (60+ event types)

---

## Story 2: "Multi-Language NLP Challenge" (Fake Follower Analysis)

**When to use:** ML questions, NLP questions, "Tell me about a creative solution," "How did you handle ambiguity?"

### Situation
Brands on our platform needed to know if an influencer's followers were real or fake. The problem was uniquely Indian -- followers had names in 10+ Indic scripts (Hindi, Bengali, Tamil, Telugu, Gujarati, Kannada, Malayalam, Odia, Punjabi, Urdu) plus English. Existing fake detection tools only worked for English. We had no labeled training data, and fake accounts often used Unicode tricks (fancy fonts like "Alicce") or transliterated names ("raahul_27").

### Task
Build an ML-powered fake follower detection system that could handle multi-language names, run serverlessly at scale, and integrate with our existing ClickHouse data pipeline.

### Action
I built a 7-step detection pipeline:

1. **Symbol conversion** -- 13 Unicode symbol variants (mathematical, circled, fullwidth letters) normalized to ASCII
2. **Language detection** -- Regex classifiers for Greek, Armenian, Georgian, Chinese, Korean scripts (non-Indic = suspicious)
3. **Indic transliteration** -- HMM-based (Hidden Markov Model) transliteration using the `indic-trans` library with pre-trained models for all 10 languages. "raahul" in Devanagari becomes "Rahul" in English
4. **Unicode normalization** using `unidecode`
5. **Handle cleaning** -- multi-stage: special chars to spaces, digits removed, lowercased

Then a **5-feature ensemble model**:
- Feature 1: Non-Indic script detection (Greek/Chinese/Korean names on Indian platform = fake)
- Feature 2: Digit count in handle (more than 4 digits = likely bot)
- Feature 3: Handle-name correlation analysis using special character logic
- Feature 4: RapidFuzz similarity scoring with weighted partial_ratio, token_sort_ratio, and token_set_ratio across all name permutations
- Feature 5: Indian name database matching against 35,183 baby names with fuzzy matching

Deployed on **AWS Lambda** with ECR containers. The pipeline: ClickHouse -> S3 -> SQS -> Lambda (processing) -> Kinesis (output stream).

### Result
- Automated fake follower detection across **10 Indic languages** where none existed before
- **50% faster processing** vs sequential approach through Lambda parallelization
- System processed followers in batches from ClickHouse, wrote results back for analytics
- Used by brands to validate influencer authenticity before signing contracts

**Code evidence:** `fake_follower_analysis/fake.py` (385 lines, core ML algorithm), `fake_follower_analysis/lambda_ecr_files/indic-trans-master/` (HMM models for 10 languages)

---

## Story 3: "The Rate Limiting Problem" (Beat)

**When to use:** API design, resilience, "Tell me about a constraint you worked around," distributed systems

### Situation
We were scraping Instagram and YouTube data from 15+ different API providers (Facebook Graph API, 6 RapidAPI providers, YouTube Data API v3, etc.). Each had different rate limits, and exceeding them meant temporary bans, wasted API credits, or permanent key revocation. With 150+ concurrent workers, we were hitting limits constantly and losing expensive API credentials.

### Task
Design a rate limiting system that respects per-provider limits, per-handle limits, and global daily quotas -- all while maximizing throughput across 73 different scraping flows.

### Action
I built a **3-level stacked rate limiting system** using Redis-backed distributed limiters (`asyncio-redis-rate-limiter`):

```
Level 1: Global daily limit    -- 20,000 requests/day per provider
Level 2: Global per-minute     -- 60 requests/minute
Level 3: Per-handle            -- 1 request/second per handle
```

Each API provider had custom limits:
- `youtube138`: 850 requests per 60 seconds
- `insta-best-performance`: 2 requests per 1 second
- `arraybobo`: 100 requests per 30 seconds
- `rocketapi`: 100 requests per 30 seconds

The limiters were **stacked** -- a request had to pass all three levels. I used Redis as the backend so limits were enforced across all worker processes and instances.

Additionally, I built a **credential rotation system**. Beat maintained multiple API credentials per provider in PostgreSQL. When a credential approached its limit, the system rotated to the next available credential. An AMQP listener (`credentials/listener.py`) consumed token refresh events from the Identity service, automatically updating expired tokens.

The worker pool itself used **semaphore-based concurrency control** -- each of the 73 flows had configurable workers and concurrency limits (e.g., Instagram profile refresh: 10 workers x 5 concurrency = 50 parallel requests).

### Result
- **25% faster API response times** by avoiding rate limit backoff penalties
- **30% cost reduction** from eliminated wasted API calls and credential bans
- Zero credential revocations after implementation
- System handled 10M+ data points daily without exceeding any provider's limits

**Code evidence:** `beat/server.py` (stacked rate limiters), `beat/utils/request.py` (HTTP rate limiting), `beat/credentials/manager.py` (credential rotation), `beat/main.py` (73 flow configurations with per-flow concurrency)

---

## Story 4: "Building the Data Platform" (Stir)

**When to use:** Data engineering questions, Airflow/dbt questions, "Tell me about a time you improved a process," ETL/ELT questions

### Situation
Metrics on the platform were 24 hours stale. Everything was batch-daily -- social media profiles, leaderboards, genre analytics, collection summaries. Brands were making campaign decisions on yesterday's data. A post could go viral and we would not reflect that until the next day's batch run.

### Task
Build a real-time data platform that could transform raw scraped data and event logs into business-ready analytics marts, with configurable freshness from 5 minutes to daily depending on the metric.

### Action
I built the entire data platform on **Apache Airflow 2.6.3 + dbt-core 1.3.1 + ClickHouse**:

**76 Airflow DAGs** organized by function:
- 11 dbt orchestration DAGs (core runs every 15 min, hourly every 2h, daily, weekly)
- 17 Instagram sync DAGs
- 12 YouTube sync DAGs
- 15 collection/leaderboard sync DAGs
- 7 asset upload DAGs
- 9 operational/verification DAGs
- 5 utility DAGs

**112 dbt models** in a two-layer architecture:
- 29 staging models (`stg_beat_*`, `stg_coffee_*`) -- raw data extraction with type casting and NULL handling
- 83 mart models organized into: audience (4), collection (13), discovery (16), genre (7), leaderboard (14), orders (1), posts (3), profile stats full (8), profile stats partial (6), staging collections (9)

**Frequency-based scheduling** was the key innovation:
- Every 5 min: `dbt_recent_scl` (recent scrape logs), `post_ranker`
- Every 15 min: `dbt_core` (critical metrics)
- Every 30 min: `dbt_collections`, `dbt_staging_collections`
- Hourly: most sync operations (12 DAGs)
- Daily: `dbt_daily`, full refresh operations (15 DAGs)
- Weekly: `dbt_weekly` (heavy aggregations)

**ClickHouse -> S3 -> PostgreSQL sync pipeline**: Stir reads from 3 ClickHouse schemas (`beat_replica`, `coffee_replica`, `_e`), transforms in dbt, then syncs results to PostgreSQL for Coffee to serve. The sync used `INSERT INTO FUNCTION s3(...)` to export, SSH download to PG server, then `COPY + atomic table swap` to minimize downtime.

### Result
- **50% latency reduction**: average data freshness went from 24 hours to under 1 hour
- Core metrics refreshed every **15 minutes** (previously 24h)
- 76 DAGs running across 5 scheduling tiers with 4 Airflow operators (PythonOperator, PostgresOperator, ClickHouseOperator, SSHOperator)
- Slack notifications for pipeline failures via `slack_failure_conn`

**Code evidence:** `stir/dags/` (76 DAG files), `stir/src/gcc_social/models/` (112 dbt models), `stir/src/gcc_social/dbt_project.yml` (dbt configuration)

---

## Story 5: "The API Gateway" (SaaS Gateway + Coffee)

**When to use:** System design, API gateway design, caching, authentication, "How did you handle cross-cutting concerns?"

### Situation
The GCC platform had grown to 13 microservices, each with its own authentication, logging, and error handling. The frontend was making direct calls to each service. There was no centralized auth, no caching layer, no request tracing, and CORS was configured inconsistently across services.

### Task
Build a unified API gateway that provides authentication, caching, logging, CORS, and reverse proxying for all 13 services. Additionally, build the main SaaS REST API (Coffee) that serves the analytics data to the frontend.

### Action

**SaaS Gateway (Go/Gin)** -- 7-layer middleware pipeline:

1. `gin.Recovery()` -- panic recovery preventing server crashes
2. `sentrygin.New()` -- error tracking with Sentry (repanic: true)
3. `cors.New()` -- CORS for 41 whitelisted origins (production, staging, localhost)
4. `GatewayContextMiddleware` -- request context injection (UUID generation, bot blocking, header extraction)
5. `RequestIdMiddleware` -- UUID v4 generation propagated in request/response headers
6. `RequestLogger` -- structured request/response logging with zerolog
7. `AppAuth()` -- JWT validation + Redis session lookup

**Two-tier caching:**
- Layer 1: **Ristretto** (in-memory LFU cache) -- 10M key capacity, 1GB max, nanosecond lookups, per-instance
- Layer 2: **Redis Cluster** -- 3-6 nodes, 100 pool size, shared across gateway instances, `session:{id}` keys

**JWT Authentication flow:**
1. Extract Bearer token from Authorization header
2. Validate JWT signature and claims
3. Check Redis for active session
4. Cache session in Ristretto for subsequent requests
5. Inject user context (userId, merchantId, clientId) into request headers for downstream services

**Coffee (Go/Chi)** -- 4-layer REST architecture:
- API Layer (handlers) -> Service Layer (transactions) -> Manager Layer (business logic) -> DAO Layer (database)
- Go generics: `Service[RES Response, EX Entry, EN Entity, I ID]`
- 10 modules: discovery, profile collection, post collection, leaderboard, collection analytics, genre insights, keyword collection, campaign profiles, collection group, partner usage
- Watermill event listeners for AMQP-based async processing
- Dual database: PostgreSQL (GORM) + ClickHouse for analytics queries

### Result
- Unified auth, logging, and CORS for **13 microservices** through a single gateway
- Two-tier caching reduced auth lookups by orders of magnitude
- Profile lookups: 200ms -> 5ms with Redis caching
- Analytics queries: 30s -> 2s using ClickHouse columnar storage
- **25% faster API responses** overall, **30% cost reduction** from caching and incremental dbt models

**Code evidence:** `saas-gateway/router/router.go` (7-layer middleware), `saas-gateway/middleware/auth.go` (JWT auth), `saas-gateway/cache/ristretto.go` + `redis.go` (two-tier caching), `coffee/core/rest/` (4-layer generic architecture), `coffee/listeners/listeners.go` (Watermill AMQP)

---

# 5. TECHNICAL DEEP-DIVE QUESTIONS (20 Questions with Answers)

---

### Q1: How did the buffered sinker pattern work in Event-gRPC?

**A:** Each RabbitMQ consumer pushes messages to a Go buffered channel (capacity: 10,000). A separate goroutine reads from this channel. It accumulates events into a slice. When the slice reaches 1,000 events OR 5 seconds elapse (whichever first), it does a single batch INSERT into ClickHouse. This reduces I/O overhead by 99% compared to per-event writes. The pattern uses `sync.Once` for lazy initialization and `safego.GoNoCtx` for panic-safe goroutine creation. If ClickHouse is temporarily down, the channel buffers messages preventing data loss until the buffer is full.

---

### Q2: Walk me through the gRPC event flow from client to ClickHouse.

**A:** Client sends a `dispatch` RPC call with a batch of Events (defined in `eventservice.proto` with 60+ event types). The gRPC server deserializes the Protobuf, pushes events to a buffered Go channel (`eventWrapperChannel`). A pool of worker goroutines reads from this channel, groups events by type using reflection, corrects future timestamps, then publishes grouped events to RabbitMQ exchange `grpc_event.tx` with the event type as routing key. On the consumer side, 26 queues are configured in `main.go`, each with dedicated consumers. The consumers either write directly to ClickHouse (for standard events) or route through buffered sinkers (for high-volume log events like post_log_events with 20 workers). Failed messages get retried twice, then routed to error queues via dead letter exchange.

---

### Q3: How did you handle Go concurrency in Event-gRPC?

**A:** Multiple patterns:
1. **sync.Once** for singleton initialization (RabbitMQ connection, worker pools, channels)
2. **Buffered channels** (1000 capacity) for backpressure between gRPC handlers and worker pools
3. **Worker pool pattern** -- configurable pool sizes, goroutines read from channels
4. **safego.GoNoCtx** -- wrapper that recovers panics in goroutines, logs them to Sentry, prevents a single panic from crashing the service
5. **sync.Mutex** for shared state in sinkers (batch accumulation)
6. **select with time.After** for flush timeouts in buffered sinkers

---

### Q4: How did async Python work in Beat?

**A:** Beat used FastAPI with **uvloop** (a libuv-based event loop that is 2-4x faster than the default asyncio loop). The worker pool system used `multiprocessing.Process` to spawn separate OS processes per flow, and within each process, `asyncio.Semaphore` controlled concurrency. HTTP calls used `aiohttp` sessions with connection pooling. The `aio-pika` library handled async RabbitMQ connections. Rate limiting used `asyncio-redis-rate-limiter` backed by Redis. SQLAlchemy async was used for database operations. Each worker process ran its own event loop, avoiding the GIL limitation for CPU-bound tasks.

---

### Q5: Explain the RabbitMQ topology across the GCC system.

**A:** Multiple exchanges with different purposes:
- `beat.dx` (direct exchange) -- Beat publishes log events (profile_log, post_log, sentiment_log, etc.)
- `coffee.dx` -- Coffee publishes activity tracking events
- `identity.dx` -- Identity service publishes token refresh events and trace logs
- `grpc_event.tx` (topic exchange) -- Event-gRPC publishes client-side events grouped by type
- `branch_event.tx` -- Branch.io attribution events
- `webengage_event.dx` -- WebEngage marketing events
- `shopify_event.dx` -- Shopify e-commerce events

Total: 26 consumer queues in Event-gRPC alone. Each queue has configurable consumer count (2-20 workers), durable queues, retry logic (max 2 retries), error queues for dead letters, and prefetch QoS of 1 for fair dispatch.

---

### Q6: Why ClickHouse? How does it differ from PostgreSQL in your system?

**A:** **PostgreSQL** was the transactional database for Beat and Coffee -- profile records, post records, credentials, collections, campaign data. Row-based storage, ACID transactions, foreign keys.

**ClickHouse** was the analytics database. Columnar storage means it only reads columns needed for a query (instead of full rows). For aggregations over millions of rows (e.g., "average engagement rate for fashion influencers last 30 days"), ClickHouse is 10-100x faster. It also provides 5x better compression because similar data types in columns compress well. In our system, ClickHouse stored event logs (_e schema), PostgreSQL replicas (beat_replica, coffee_replica for dbt), and dbt transformation results. The tradeoff: ClickHouse does not support UPDATE/DELETE well, so transactional data stays in PostgreSQL.

---

### Q7: How did dbt work with ClickHouse?

**A:** We used `dbt-clickhouse 1.3.2`. dbt models were SQL SELECT statements that materialized as tables or views in ClickHouse. The `dbt_project.yml` configured the ClickHouse connection. dbt models referenced three source schemas: `beat_replica`, `coffee_replica`, and `_e`. The staging layer (29 models) did type casting and NULL handling. The mart layer (83 models) did business logic -- joins, aggregations, window functions. dbt handled dependency ordering automatically. Incremental models only processed new data since last run, reducing compute costs. Airflow orchestrated dbt runs via `DbtRunOperator` with tag-based selection (`dbt run --models tag:core`).

---

### Q8: Explain the ClickHouse -> S3 -> PostgreSQL sync pipeline.

**A:** This was how analytics marts became available to Coffee's REST API:
1. dbt transforms data in ClickHouse (e.g., `dbt.mart_leaderboard`)
2. Airflow ClickHouseOperator runs: `INSERT INTO FUNCTION s3('gcc-social-data/tmp/leaderboard.json', 'JSONEachRow')`
3. Airflow SSHOperator downloads the file from S3 to the PostgreSQL server
4. Airflow PostgresOperator does an atomic table swap: create temp table, COPY data in, rename tables, drop old
5. Coffee reads from the fresh PostgreSQL table

The atomic swap ensured zero downtime -- readers always see either the old or new data, never partial.

---

### Q9: How did multi-tenancy work in Coffee?

**A:** Every request passed through `ApplicationContext` middleware that extracted `x-bb-*` headers:
- `x-bb-userid` -- the authenticated user
- `x-bb-merchantid` -- the tenant (brand/agency)
- `x-bb-clientid` -- the application client
- `x-bb-ppid` -- partner platform ID

These populated a `RequestContext` struct available throughout the request lifecycle. The DAO layer appended tenant filters to all queries automatically. Collections, profiles, and analytics were all scoped to the merchant. The gateway injected these headers after JWT validation, so downstream services trusted them.

---

### Q10: Describe the fake follower ML ensemble model.

**A:** Five independent features scored 0 or 1 (fake/real), then combined:

1. **Non-Indic script detection**: regex for Greek, Armenian, Georgian, Chinese, Korean characters. Real Indian followers rarely use foreign scripts.
2. **Digit count**: more than 4 digits in handle = likely bot (e.g., `user_999999`)
3. **Handle-name correlation**: if handle has special chars (_, -, .) but the cleaned name does not match the handle (similarity < 80), likely fake
4. **Fuzzy similarity**: RapidFuzz weighted scoring -- (2 * partial_ratio + token_sort_ratio + token_set_ratio) / 4, tested across all name permutations
5. **Indian name database**: 35,183 baby names, fuzzy matched against follower name parts

The `final()` function combined these into a weighted score: 0.0 (definitely real), 0.33 (uncertain), 1.0 (definitely fake). Output was a 19-field JSON response per follower.

---

### Q11: How did the Indic transliteration work?

**A:** We used Hidden Markov Models (HMMs) from the `indic-trans` library. Each language pair (hin-eng, ben-eng, tam-eng, etc.) had pre-trained model files: coefficient matrices (`coef_.npy`), character class mappings (`classes.npy`), initial/transition/final state probabilities, and a feature vocabulary (`sparse.vec`). The Viterbi algorithm (implemented in Cython for speed) found the most likely character-by-character transliteration. For example, "raahul" in Devanagari -> the HMM outputs "Rahul" in English. I also built a Hindi-specific rule-based system (`createDict.py`) using vowel (svar.csv, 24 mappings) and consonant (vyanjan.csv, 42 mappings) tables.

---

### Q12: How did the Airflow DAG scheduling work?

**A:** Five scheduling tiers:
- **Every 5 min**: recent scrape logs, post ranking (near real-time user-facing data)
- **Every 15 min**: `dbt_core` -- core metrics like account summaries, leaderboard bases
- **Every 30 min**: collection processing (user-created influencer lists)
- **Hourly**: most Instagram/YouTube sync operations
- **Daily/Weekly**: full refresh aggregations, heavy computations

Each DAG used appropriate operators: `PythonOperator` (46 DAGs) for API calls and processing, `PostgresOperator` (20) for data loading, `ClickHouseOperator` (19) for export queries, `SSHOperator` (18) for file transfers, `DbtRunOperator` (11) for dbt execution. Connections: `clickhouse_gcc`, `prod_pg`, `stage_pg`, `ssh_prod_pg`, `beat`.

---

### Q13: Explain the SaaS Gateway's JWT authentication flow.

**A:** The `AppAuth()` middleware:
1. Extract Bearer token from `Authorization` header
2. Decode JWT using `golang-jwt/jwt` -- validate signature, check expiration
3. Extract claims: userId, sessionId, clientType
4. Check Ristretto (in-memory) cache for session data -- nanosecond lookup
5. If miss, check Redis cluster for `session:{sessionId}` key
6. If Redis has it, cache in Ristretto for future requests
7. If not in Redis, call Identity Service `/auth/verify/token` API
8. Cache result in both Ristretto and Redis
9. Inject user context into request headers for downstream services
10. If any step fails, return 401 Unauthorized

The two-tier cache meant most auth checks completed in nanoseconds (Ristretto hit) rather than milliseconds (Redis) or hundred-milliseconds (Identity API call).

---

### Q14: How did you handle AWS Lambda deployment for the ML model?

**A:** The system used ECR (Elastic Container Registry) with a custom Docker image:
- Base: `public.ecr.aws/lambda/python:3.10`
- System dependencies: `gcc-c++` for Cython compilation of Viterbi algorithm
- Installed `indic-trans` from source (custom `setup.py`)
- Copied pre-trained HMM models for 10 languages into site-packages
- Included `baby_names_.csv` (35,183 names), `svar.csv` (vowels), `vyanjan.csv` (consonants)
- Entry point: `CMD ["fake.handler"]`

The pipeline: `push.py` queries ClickHouse for follower data, uploads to S3, sends messages to SQS (`creator_follower_in`). Lambda triggers on SQS, processes each batch, writes results to Kinesis (`creator_out`). This serverless approach auto-scaled with load.

---

### Q15: How did Coffee's 4-layer generic architecture work?

**A:** Go 1.18 generics enabled type-safe layers:

- **API Layer** (`ApiWrapper` interface): parses HTTP requests, calls Service, returns JSON
- **Service Layer** (`Service[RES, EX, EN, I]`): manages transactions (`InitializeSession`, `Close`, `Rollback`), orchestrates Manager calls
- **Manager Layer** (`Manager[EX, EN, I]`): business logic, `toEntity()`/`toEntry()` transformations
- **DAO Layer** (`DaoProvider[EN, I]`): GORM-based database operations, `FindById`, `Create`, `Search`, `SearchJoins`

Each module (discovery, leaderboard, etc.) implemented these interfaces with its own types. The generic Service handled transaction management identically across all modules. This eliminated boilerplate -- adding a new entity type meant implementing the 4 interfaces, and all CRUD + search + pagination worked automatically.

---

### Q16: What was the Watermill integration in Coffee?

**A:** Watermill is a Go library for event-driven architectures. Coffee used it for AMQP (RabbitMQ) event consumption. The `listeners/listeners.go` file registered handlers for events from Stir (upsert notifications when data was synced to PostgreSQL) and from Beat (real-time profile updates). When Stir completed a sync (e.g., new leaderboard data), it published an event to `coffee.dx`, and Coffee's Watermill listener processed it -- invalidating Redis cache entries so the next API request would read fresh data. This was how real-time data propagated from scraping to the API.

---

### Q17: How did you handle database connections in Go services?

**A:** Three connection patterns:

1. **PostgreSQL (GORM)**: connection pool with configurable max open/idle connections, `gorm.Open()` with `postgres.Open(dsn)`. Used `GetSession(ctx)` to get request-scoped sessions for transaction management.

2. **ClickHouse**: multi-database connection with auto-reconnect. `clickhouse-go` driver with batch insert support. Configured pool sizes (10 idle, 20 max open).

3. **Redis**: both go-redis (cluster client, 100 pool size, 3-6 nodes) and Ristretto (in-memory LFU, 10M keys, 1GB). Singleton pattern with `sync.Once` for initialization.

All services used `scripts/start.sh` for graceful deployment -- drain connections, wait for in-flight requests, then shut down.

---

### Q18: How did Beat's worker pool system work?

**A:** `main.py` defined 73 flow configurations with workers and concurrency:

```
'refresh_profile_by_handle': {'no_of_workers': 10, 'no_of_concurrency': 5}
```

This spawns 10 OS processes (via `multiprocessing.Process`), each running an async event loop. Within each process, `asyncio.Semaphore(5)` limits concurrent API calls. Each worker runs a `looper()` function that:
1. Polls `scrape_request_log` table for pending tasks (SQL-based task queue)
2. Acquires semaphore
3. Calls the appropriate flow function (e.g., `refresh_profile.py`)
4. Updates task status in PostgreSQL
5. Publishes events to RabbitMQ

Total capacity: 73 flows x variable workers x variable concurrency = 150+ concurrent operations.

---

### Q19: How did you monitor the system?

**A:** Multi-layer observability:
- **Prometheus metrics** in both Go services (SaaS Gateway and Coffee via `chi-prometheus`, Event-gRPC via `prometheus/client_golang`)
- **Sentry** for error tracking (all services: `getsentry/sentry-go` for Go, `sentry-sdk` for Python)
- **Structured logging**: zerolog in Gateway, logrus in Coffee/Event-gRPC, loguru in Stir
- **Airflow UI** for DAG monitoring, success/failure rates, SLA tracking
- **Slack notifications** for pipeline failures (`slack_failure_conn`, `slack_success_conn` in Airflow)
- **Health check endpoints**: `/heartbeat` in all services, toggleable via PUT for graceful shutdown during deployments

---

### Q20: What was the complete tech stack?

**A:**
| Category | Technologies |
|----------|-------------|
| Languages | Go 1.14-1.18, Python 3.11, SQL, Protobuf |
| Go Frameworks | Gin, Chi, GORM, Watermill, go-resty |
| Python Frameworks | FastAPI, uvloop, aio-pika, SQLAlchemy async |
| Databases | PostgreSQL (transactional), ClickHouse (analytics), Redis (cache) |
| Message Broker | RabbitMQ (AMQP) |
| Data Platform | Apache Airflow 2.6.3, dbt-core 1.3.1, dbt-clickhouse 1.3.2 |
| ML/NLP | RapidFuzz, indic-trans (HMM), pandas, numpy, unidecode |
| AWS | Lambda, ECR, S3, SQS, Kinesis, CloudFront |
| Caching | Ristretto (in-memory LFU), Redis Cluster |
| CI/CD | GitLab CI/CD |
| Monitoring | Prometheus, Sentry, Slack, zerolog/logrus/loguru |
| API Design | gRPC + Protocol Buffers, REST (OpenAPI) |

---

# 6. BEHAVIORAL QUESTIONS MAPPED TO GCC

---

### Q1: "Tell me about a time you dealt with ambiguity."

> "When I started the fake follower detection project, there was no clear spec -- just 'detect fake followers for Indian influencers.' Nobody had done multi-language fake detection before. I broke the ambiguity down: first I studied what makes a follower fake (bot-like usernames, foreign scripts, no correlation between handle and display name). Then I identified the multi-language problem as the hardest piece. I researched HMM-based transliteration, found the indic-trans library, and built a prototype with just Hindi first. Once that worked, I extended to all 10 Indic languages. Instead of waiting for perfect requirements, I shipped iteratively -- first English-only detection, then Indic transliteration, then the full 5-feature ensemble."

---

### Q2: "Tell me about a time you had to learn quickly."

> "When I joined GCC, I had no Go experience. The first project assigned was Event-gRPC -- a high-throughput Go service handling 10K events per second. I learned Go's concurrency model (goroutines, channels, sync primitives) in the first two weeks by reading the Event-gRPC codebase and Go documentation. Within a month, I was writing production Go code. The buffered sinker pattern I built -- using channels as bounded queues with flush timers -- came from deeply understanding Go's concurrency primitives. I then applied the same Go knowledge to Coffee and SaaS Gateway."

---

### Q3: "Tell me about your biggest technical challenge."

> "The PostgreSQL-to-ClickHouse migration for event logging. We had 10M+ daily events overwhelming PostgreSQL with 100x write latency degradation. The challenge was not just picking ClickHouse -- it was designing the migration path. I could not do a big-bang cutover because the platform was live. I designed a dual-write strategy: Beat continued writing to PostgreSQL AND published to RabbitMQ. Event-gRPC consumed from RabbitMQ and wrote to ClickHouse. For weeks, both databases had the same data. I built validation queries comparing counts and checksums. Once confidence was high, I switched reads to ClickHouse. The hardest part was the buffered sinker design -- getting the batch size (1000) and flush interval (5 seconds) right required load testing to balance latency vs throughput."

---

### Q4: "Tell me about a time you disagreed with a teammate."

> "Our team lead wanted to use Kafka for event streaming between services. I advocated for RabbitMQ. My reasoning: we were a small startup with 4 engineers. Kafka requires ZooKeeper, partition management, and operational overhead that our team could not support. RabbitMQ gave us routing keys, dead letter exchanges, and per-queue consumer counts with much simpler operations. We were processing 10K events/sec, well within RabbitMQ's capacity. I presented a comparison: RabbitMQ would take 1 day to set up vs a week for Kafka, and our throughput needs did not justify Kafka's complexity. The team lead agreed. It was the right call -- we never hit RabbitMQ's limits, and the operational simplicity let us focus on features."

---

### Q5: "Tell me about a time you had to prioritize competing demands."

> "I was simultaneously building Stir (data platform) and maintaining Beat (scraper). A critical bug appeared in Beat -- Instagram credentials were being revoked because our rate limiter was not accounting for a new API provider's limits. At the same time, the data platform was supposed to launch that week. I triaged: the Beat bug was losing us money (API credits) every hour. I fixed it in half a day by adding the new provider's rate limits to the stacked limiter config and deploying. Then I shifted back to Stir. I deprioritized the nice-to-have weekly aggregation DAGs and launched with the core 5-minute and 15-minute DAGs. We hit the launch date with the most impactful features, then backfilled the weekly DAGs the following week."

---

### Q6: "Tell me about a time you improved a process."

> "Data syncing from ClickHouse to PostgreSQL was manual. Engineers ran SQL queries, exported CSVs, and loaded them into PostgreSQL with downtime. I automated the entire pipeline in Stir: ClickHouse INSERT INTO FUNCTION s3() exports, SSHOperator downloads, PostgreSQL COPY with atomic table swaps. What used to take 30 minutes of manual work (and caused downtime) became a fully automated DAG running on schedule with Slack notifications on failure. I built it for one table (leaderboard), then templated it for all 15+ sync targets. This freed up hours of engineering time per week."

---

### Q7: "Tell me about a time you showed ownership."

> "The SaaS Gateway was originally a thin proxy with no caching. I noticed that every API request triggered a full JWT verification against the Identity Service, adding 200ms+ latency. Nobody asked me to fix this. I added the two-tier caching layer: Ristretto for in-memory nanosecond lookups, Redis for shared state across gateway instances. I also added the full middleware pipeline -- request logging, request ID generation, CORS configuration, Sentry error tracking. The gateway went from a 50-line proxy to a proper API gateway with observability. API response times improved by 25% across the entire platform."

---

### Q8: "Tell me about how you handle technical debt."

> "Beat's server.py grew to 43KB -- a massive file handling REST endpoints, rate limiting, credential management, and flow orchestration all in one place. While I could not do a full refactor during sprint, I incrementally extracted modules: credentials got their own package (`credentials/manager.py`, `validator.py`, `listener.py`), rate limiting logic moved to `utils/request.py`, and flow configurations moved to a separate dict in `main.py`. Each extraction was a self-contained PR that did not change behavior. Over 3 months, the code became much more maintainable without any big-bang refactor."

---

### Q9: "Tell me about a time you mentored someone or shared knowledge."

> "When a new engineer joined the team, the data pipeline was opaque -- events flowed from Beat through RabbitMQ through Event-gRPC into ClickHouse through Stir into Coffee, and nobody had documented it. I created the System Interconnectivity document mapping every service-to-service connection, every RabbitMQ exchange and queue, every database table relationship. I walked the new engineer through the entire data flow using real code references. This document became the team's architecture bible -- everyone referenced it when debugging cross-service issues."

---

### Q10: "Why did you leave GCC / What did you learn?"

> "GCC was a startup, and startups are a great learning environment. In 15 months, I went from zero Go experience to building production microservices. I learned how to design event-driven architectures, build data platforms, and handle the constraints of a small team (no dedicated SRE, no dedicated DBA -- you own everything). I left because I wanted to work at scale. Walmart processes 2M+ transactions daily across 1,200+ suppliers in 3 markets -- a different magnitude of complexity. But the fundamentals I learned at GCC -- async processing, schema enforcement, decoupled architectures, dual-database strategies -- translate directly."

---

# 7. "WHY X NOT Y?" CROSS-PROJECT DECISIONS

---

### Decision 1: Why RabbitMQ across all services vs Kafka?

> "Three reasons. First, **operational simplicity** -- we were a 4-person engineering team. Kafka requires ZooKeeper (or KRaft), partition management, consumer group rebalancing, and offset management. RabbitMQ is a single binary with a management UI. Second, **routing flexibility** -- RabbitMQ's exchange-binding-queue model with routing keys let us fan out events to multiple consumers easily (one profile_log event goes to ClickHouse sinker AND activity tracker). Kafka's topic model would require consumer groups and more complex routing. Third, **throughput fit** -- our peak was 10K events/sec. RabbitMQ handles this comfortably. Kafka shines at 100K+/sec. We would be paying complexity tax for capacity we did not need."

---

### Decision 2: Why ClickHouse + PostgreSQL vs a single database?

> "Different access patterns demand different storage engines. PostgreSQL excels at: point lookups by ID, transactional updates (profile edits, collection mutations), foreign key constraints, and ACID guarantees. ClickHouse excels at: aggregating millions of rows (SUM, AVG, COUNT GROUP BY), time-series queries, columnar compression (5x better than PostgreSQL for our data), and append-only event storage. Our analytics queries (e.g., 'average engagement rate for 10K influencers over 30 days') ran 10-100x faster in ClickHouse. But ClickHouse does not support efficient UPDATE/DELETE or row-level transactions. So we used both: PostgreSQL for CRUD operations, ClickHouse for analytics. dbt bridged them by transforming ClickHouse data and syncing aggregated results back to PostgreSQL."

---

### Decision 3: Why Go for services + Python for data vs a single language?

> "Each language excelled at its task. **Go** for Event-gRPC, Coffee, and SaaS Gateway because: goroutines give excellent concurrency for high-throughput request handling, static typing catches bugs at compile time, single binary deployment is simple, and the standard library is production-quality. **Python** for Beat (scraping) because: the richest ecosystem of API client libraries (aiohttp, rapidfuzz, aio-pika), data manipulation (pandas, numpy), and ML libraries (scikit-learn, indic-trans). Python for Stir because: Airflow and dbt are Python-native tools. Using Go for the data pipeline would mean fighting the ecosystem."

---

### Decision 4: Why GitLab CI vs GitHub Actions?

> "GCC's codebase was hosted on GitLab, so GitLab CI was the natural choice. It was already integrated -- no additional setup, no webhook configuration, no secrets management bridge. Each project had a `.gitlab-ci.yml` file. The pipeline was simple: lint, test, build Docker image, deploy to staging/production. For a small team, minimizing DevOps overhead was critical."

---

### Decision 5: Why self-hosted servers vs Kubernetes?

> "Cost and operational simplicity. Kubernetes adds a significant overhead: cluster management, node scaling, YAML proliferation, networking complexity. With 6 microservices and a 4-person team, a few self-hosted servers (cb1-1, cb2-1 seen in configs) with systemd services were sufficient. Each service had a `scripts/start.sh` that handled graceful deployment (drain connections, health check toggle, process restart). We used a simple load balancer in front. If the team had grown to 20+ engineers or if we needed auto-scaling, Kubernetes would have been the right call."

---

### Decision 6: Why async Python (uvloop) vs Go for scraping (Beat)?

> "Scraping is I/O-bound, not CPU-bound. The workers spend 95% of their time waiting for API responses. Python's async/await with uvloop gives excellent I/O concurrency with simpler code than Go for this use case. More importantly, the Python ecosystem for scraping is unmatched: `aiohttp` for async HTTP, `rapidfuzz` for fuzzy matching, `aio-pika` for async RabbitMQ, `asyncio-redis-rate-limiter` for distributed rate limiting, YAKE for keyword extraction, and 128 dependencies in `requirements.txt` -- most of which have no Go equivalent. Go would have been faster for CPU-bound work, but for I/O-bound scraping with heavy library needs, Python was the pragmatic choice."

---

### Decision 7: Why SQL-based task queue (scrape_request_log) vs a dedicated queue service?

> "Beat used PostgreSQL's `scrape_request_log` table as a task queue. Workers poll for tasks with `status=PENDING`, update to `IN_PROGRESS`, process, then update to `COMPLETED` or `FAILED`. Why not Redis Queue (RQ) or Celery? Simplicity and visibility. SQL gives us: full task history (when was it queued, started, completed), easy debugging (`SELECT * FROM scrape_request_log WHERE status='FAILED'`), natural retry logic (`UPDATE ... SET status='PENDING' WHERE status='FAILED' AND retry_count < 3`), and priority queuing (`ORDER BY priority ASC`). A dedicated queue service would add another infrastructure dependency. For our throughput (thousands of tasks, not millions), PostgreSQL was sufficient."

---

### Decision 8: Why Watermill + AMQP for Coffee vs direct RabbitMQ for Beat/Event-gRPC?

> "Different maturity levels and needs. Beat used `aio-pika` directly because Python's AMQP libraries are straightforward. Event-gRPC used a custom `rabbit/rabbit.go` wrapper with singleton pattern and auto-reconnect because it needed fine-grained control (26 consumer configurations, buffered channels, retry logic). Coffee used **Watermill** because it provides a higher-level abstraction: message routers, middleware (retry, throttle, poison queue), and a clean subscriber/publisher interface. Watermill let Coffee's developers focus on business logic (what to do when a profile is updated) without worrying about connection management or error handling. It was the right abstraction for a service that consumes events but is not itself an event infrastructure component."

---

### Decision 9: Why separate Beat and Event-gRPC vs a single service?

> "Separation of concerns and language optimization. Beat is a **producer** -- it scrapes external APIs and generates events. Event-gRPC is a **consumer** -- it ingests events and writes to databases. They have fundamentally different scaling profiles: Beat scales with the number of external APIs (I/O-bound, Python excels), Event-gRPC scales with event ingestion throughput (CPU-bound batching, Go excels). They also have different failure modes: if scraping fails, events stop being produced but existing events are still processed. If ingestion fails, events buffer in RabbitMQ (durable queues) until the service recovers. Combining them would mean a single failure takes down both data collection and data ingestion."

---

### Decision 10: Why serverless Lambda for ML vs running on Beat?

> "Three reasons. First, **resource isolation** -- the fake follower ML model loads 10 language models, a 35K name database, and runs fuzzy matching on every follower. This is CPU-intensive and would starve Beat's I/O-bound scraping workers. Second, **auto-scaling** -- some influencers have 10M followers, others have 1K. Lambda scales to zero when there is no work and scales up to hundreds of concurrent executions for large analyses. Beat has fixed worker counts. Third, **deployment independence** -- the ML model changes frequently (new features, tuned thresholds) while Beat's scraping logic is stable. Lambda/ECR lets us deploy model updates without touching the scraping infrastructure."

---

# 8. NUMBERS CHEAT SHEET

## All GCC Numbers in One Table

| Category | Metric | Value | Source Project |
|----------|--------|-------|----------------|
| **Scale** | Daily data points processed | 10M+ | Beat + Event-gRPC |
| | Events/sec ingestion | 10K+ | Event-gRPC |
| | Images uploaded/day | 8M | Beat (main_assets.py) |
| | Concurrent workers (scraping) | 150+ | Beat |
| | RabbitMQ consumer queues | 26 | Event-gRPC |
| | Event-gRPC concurrent workers | 90+ | Event-gRPC |
| | Scraping flows configured | 73 (75+) | Beat |
| | External API integrations | 15+ | Beat |
| | Airflow DAGs | 76 | Stir |
| | dbt models (total) | 112 | Stir |
| | dbt staging models | 29 | Stir |
| | dbt mart models | 83 | Stir |
| | Coffee REST endpoints | 40+ | Coffee |
| | Coffee modules | 10 | Coffee |
| | Downstream services proxied | 13 | SaaS Gateway |
| | CORS whitelisted origins | 41 | SaaS Gateway |
| | Indic languages supported | 10 | Fake Follower |
| | Baby name database size | 35,183 | Fake Follower |
| | Protobuf event types | 60+ | Event-gRPC |
| | ClickHouse event tables | 21+ | Event-gRPC |
| **Performance** | Log retrieval improvement | 2.5x faster | Event-gRPC + ClickHouse |
| | API response improvement | 25% faster | Coffee + Gateway |
| | Cost reduction | 30% | Caching + incremental dbt |
| | Data latency reduction | 50% (24h -> <1h) | Stir |
| | ClickHouse compression | 5x vs PostgreSQL | Event-gRPC |
| | Analytics query speedup | 30s -> 2s | ClickHouse vs PostgreSQL |
| | Profile lookup speedup | 200ms -> 5ms | Redis cache |
| | Write I/O reduction | 99% (via batching) | Event-gRPC buffered sinker |
| | Batch write improvement | 33x throughput | Event-gRPC buffered sinker |
| | Processing speed improvement | 50% faster | Fake Follower (Lambda) |
| **Code** | Total LOC (all projects) | ~60,000+ | All 6 projects |
| | Beat LOC | ~15,000+ | Beat |
| | Stir LOC | ~17,500+ | Stir |
| | Event-gRPC LOC | ~10,000+ | Event-gRPC |
| | Coffee LOC | ~8,500+ | Coffee |
| | Fake Follower LOC | ~955+ | Fake Follower |
| | SaaS Gateway LOC | included in total | SaaS Gateway |
| | Python dependencies (Beat) | 128 | Beat |
| | Python dependencies (Stir) | 350+ | Stir |
| **Infrastructure** | Buffered channel capacity | 10,000 messages | Event-gRPC |
| | Batch insert size | 1,000 events | Event-gRPC |
| | Flush interval | 5 seconds | Event-gRPC |
| | Ristretto cache capacity | 10M keys, 1GB | SaaS Gateway |
| | Redis cluster nodes | 3-6 | SaaS Gateway |
| | Redis pool size | 100 connections | SaaS Gateway |
| | Middleware layers (Gateway) | 7 | SaaS Gateway |
| | Middleware layers (Coffee) | 10 | Coffee |
| | Coffee schema tables | 27+ | Coffee |
| | Stir Git commits | 1,476 | Stir |
| | Retry max attempts | 2 | Event-gRPC |
| | Rate limit (daily global) | 20,000 req/day | Beat |
| | Rate limit (per-minute) | 60 req/min | Beat |
| | Rate limit (per-handle) | 1 req/sec | Beat |

---

# 9. HOW TO PIVOT FROM GCC TO WALMART

## The Pivot Framework

When asked about GCC experience, always bridge to how it prepared you for Walmart. Use this mapping:

| GCC Experience | Walmart Parallel | Pivot Statement |
|---------------|-----------------|-----------------|
| RabbitMQ event streaming | Kafka event streaming | "At GCC I used RabbitMQ for 10K events/sec. At Walmart, same pattern at larger scale with Kafka handling 2M+ daily audit events with Avro schemas and Connect sink connectors." |
| ClickHouse analytics | BigQuery / PostgreSQL analytics | "GCC taught me to separate transactional and analytics workloads. At Walmart I apply the same principle with PostgreSQL for API serving and GCS/BigQuery for analytical queries." |
| Python async (uvloop) | Java CompletableFuture | "GCC's async scraping with 150+ workers translates directly to Walmart's parallel processing with CompletableFuture for bulk inventory queries." |
| dbt data models | Spring Boot data services | "Building 112 dbt models taught me data modeling discipline. At Walmart, I apply that to designing API response schemas and database partitioning strategies." |
| Multi-tenancy (merchantId) | Multi-tenancy (supplier DUNS) | "GCC's multi-tenant SaaS with merchantId scoping became Walmart's supplier authorization framework with Consumer -> DUNS -> GTIN -> Store hierarchy." |
| SaaS Gateway (Gin + JWT) | Spring Boot + Strati AF | "Building an API gateway from scratch at GCC gave me deep understanding of middleware patterns, auth flows, and caching. At Walmart, I work within the Strati Application Framework but understand exactly what it does under the hood." |
| GitLab CI/CD | KITT + Looper + Canary | "GCC's simple CI/CD pipeline evolved into Walmart's enterprise-grade deployment: KITT for orchestration, Looper for pipelines, Flagger for canary deployments, Istio for service mesh." |
| Self-hosted servers | Kubernetes (WCNP) | "Going from GCC's self-hosted servers to Walmart's Kubernetes clusters (WCNP) was a natural progression. Same services, much more sophisticated orchestration." |
| Small team (4 engineers) | Large org (12+ teams) | "At GCC I owned everything end-to-end. At Walmart, I learned to build reusable libraries adopted by 12+ teams -- different skill, same engineering principles." |
| Sentry + Prometheus | OpenTelemetry + Dynatrace + Grafana | "GCC's monitoring was Sentry + Prometheus. Walmart taught me enterprise observability: distributed tracing with OpenTelemetry, APM with Dynatrace, custom Grafana dashboards." |

## Sample Pivot Answers

**"How does your previous experience apply to Walmart?"**

> "At GCC, I built event-driven data pipelines and multi-tenant microservices at startup scale. The patterns translate directly: async processing, schema enforcement, decoupled architectures, dual-database strategies. At Walmart, the scale is different -- millions of transactions instead of thousands, Kafka instead of RabbitMQ, Kubernetes instead of bare metal -- but the architectural principles are the same. GCC gave me the foundation; Walmart gave me enterprise scale."

**"What is the biggest difference between GCC and Walmart?"**

> "Operational rigor. At GCC, if a service went down at 2 AM, I would SSH in and restart it. At Walmart, we have canary deployments, automated rollbacks, SLO tracking, PagerDuty rotations, and 99.9% uptime requirements. The code quality bar is also higher -- contract testing, SonarQube, 80%+ code coverage. GCC was about shipping fast; Walmart is about shipping reliably at scale."

---

# 10. DANGER ZONES - What NOT to Say

---

### Danger Zone 1: Do NOT Claim You Built Everything Alone

**Wrong:** "I built the entire platform."
**Right:** "I was a core contributor on a small team. I owned these six microservices end-to-end, but I worked with the frontend team, a product manager, and other backend engineers."

Why: Interviewers will probe this. If you claim sole ownership, they will ask "Did you design the database schema?" for every piece, and inconsistencies will erode credibility. The truth -- that you were a key SWE-I on a small startup team who owned 6 backend services -- is already incredibly impressive.

---

### Danger Zone 2: Do NOT Exaggerate Team Size

**Wrong:** "Our team of 20 engineers..."
**Right:** "We were a small startup team -- about 4-5 backend engineers. That is why I owned multiple services."

Why: The small team is a feature, not a bug. It shows you can handle ambiguity, own systems end-to-end, and are not a cog in a large machine. Inflating the team size actually undermines your story.

---

### Danger Zone 3: Do NOT Mix Up ClickHouse vs PostgreSQL Roles

**Wrong:** "We stored everything in ClickHouse."
**Right:** "PostgreSQL for transactional data (profiles, posts, collections, credentials). ClickHouse for analytics (event logs, time-series, aggregations). dbt transformed ClickHouse data and synced results to PostgreSQL for the API."

Why: The dual-database strategy is one of your strongest talking points. Confusing which database serves which purpose signals shallow understanding.

---

### Danger Zone 4: Do NOT Forget It Was a Startup

**Wrong:** Describing GCC with enterprise-level processes (formal code reviews, extensive testing, etc.)
**Right:** "It was a startup, so we prioritized shipping speed. We had GitLab CI/CD for automated deployments, but formal processes like code coverage targets and contract testing came later at Walmart."

Why: Interviewers from large companies understand startup tradeoffs. Pretending it was enterprise-grade makes them suspicious. Acknowledging startup constraints and then showing how you grew into enterprise practices at Walmart is a much better narrative.

---

### Danger Zone 5: Do NOT Overstate ML Expertise

**Wrong:** "I trained the HMM models from scratch."
**Right:** "I used pre-trained HMM transliteration models from the indic-trans library and built the ensemble detection pipeline around them. The ML was more about feature engineering and integration than model training."

Why: If you claim to have trained HMMs, an ML-focused interviewer will ask about training data, loss functions, hyperparameter tuning, and cross-validation. You will not have good answers. Your real strength is the end-to-end pipeline: data extraction, feature engineering, multi-language handling, and serverless deployment.

---

### Danger Zone 6: Do NOT Use "Celery" When Describing Beat

**Wrong:** "Beat used Celery workers."
**Right:** "Beat used a custom worker pool with multiprocessing.Process and asyncio.Semaphore for concurrency control, with a SQL-based task queue."

Why: The resume says "Celery workers" but the actual code uses `multiprocessing.Process` with a SQL task queue. If asked to explain Celery internals, you will be caught. The actual implementation is arguably more interesting to discuss.

---

### Danger Zone 7: Do NOT Confuse Event-gRPC's gRPC Server with Its Consumer Role

The service has two roles:
1. **gRPC server** (port 8017): receives events from client apps via Protobuf, publishes to RabbitMQ
2. **RabbitMQ consumer**: consumes from 26 queues, writes to ClickHouse

Most of the complexity (and your talking points) are in role #2 (the consumer/sinker side). Do not spend all your time explaining the gRPC server when the buffered sinker pattern is the real engineering innovation.

---

### Danger Zone 8: Do NOT Say "Kafka" When You Mean "RabbitMQ"

GCC used **RabbitMQ**. Walmart uses **Kafka**. Do not accidentally swap them. If asked "Did you use Kafka at GCC?", the answer is "No, we used RabbitMQ. At Walmart I use Kafka." Then explain the tradeoff (Decision 1 above).

---

### Danger Zone 9: Do NOT Forget the Timeline

**Feb 2023 - May 2024** (15 months) at GCC. This is before Walmart (Jun 2024 - Present). If you accidentally say you built these services last year, it conflicts with your Walmart timeline.

---

### Danger Zone 10: Keep Numbers Consistent

These are the canonical numbers. Do not improvise different numbers in the interview:

- Daily data points: **10M+** (not "hundreds of millions")
- Events/sec: **10K+** (not "millions")
- Workers: **150+** (not "thousands")
- DAGs: **76** (not "hundreds")
- dbt models: **112** (not "200+")
- LOC: **~60K** (not "100K")
- Services: **6** (not "dozen")
- Queues: **26** (not "50+")

Inflating numbers is the fastest way to lose credibility. The real numbers are already impressive.

---

*This is the master guide. Read it before every GCC-related interview. The 5 STAR stories cover system design, ML/NLP, API resilience, data engineering, and API gateway design -- the full breadth of your GCC work. The 20 technical deep-dives prepare you for follow-up questions. The 10 behavioral mappings show GCC experience through a behavioral lens. The 10 "Why X not Y" decisions show architectural judgment. The danger zones keep you from shooting yourself in the foot.*

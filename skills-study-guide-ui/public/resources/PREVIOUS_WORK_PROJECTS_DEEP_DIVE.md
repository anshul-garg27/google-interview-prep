# Good Creator Co Projects - Ultra Deep Dive

**Complete Technical Documentation for Interview Preparation**

---

## Quick Reference Table

| Project | Tech Stack | Impact | Lines of Code | Your Role |
|---------|-----------|--------|---------------|-----------|
| **beat** | Python 3.11, FastAPI, PostgreSQL, Redis | 10M+ data points/day, 150+ workers | ~15,000 | Rate limiting, worker pools, event publishing |
| **stir** | Airflow 2.6, dbt 1.3, ClickHouse, PostgreSQL | 50% latency reduction, 76 DAGs | ~17,500 | ETL pipelines, 112 dbt models |
| **event-grpc** | Go 1.14, gRPC, RabbitMQ, ClickHouse | 10K+ events/sec, 26 queues | ~10,000 | Buffered sinkers, event streaming |
| **fake_follower_analysis** | Python 3.10, AWS Lambda, ML/NLP | 85% accuracy, 10 languages | ~955 | End-to-end ML system |
| **coffee** | Go 1.18, Gin, PostgreSQL, ClickHouse | 12 modules, 50+ endpoints | ~8,500 | Multi-tenant REST API |
| **saas-gateway** | Go 1.12, Gin, Redis, JWT | 13 services, 2-tier caching | ~5,000 | API gateway, auth middleware |

---

## System Architecture Overview

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    GOOD CREATOR CO. PLATFORM                             │
│                   (Social Media Analytics SaaS)                          │
└─────────────────────────────────────────────────────────────────────────┘

CLIENT LAYER
┌─────────────────────────────────────────────────────────────────────────┐
│  Web App (React)  │  Mobile App (React Native)  │  Browser Extensions  │
└────────────────────────────────┬────────────────────────────────────────┘
                                 │
                                 ▼
API GATEWAY LAYER
┌─────────────────────────────────────────────────────────────────────────┐
│                          saas-gateway (Go)                               │
│  • 7-layer middleware pipeline                                           │
│  • JWT + Redis session validation                                        │
│  • 2-tier caching (Ristretto + Redis)                                   │
│  • Routes to 13 downstream services                                      │
└────────────────────────────────┬────────────────────────────────────────┘
                                 │
                                 ▼
SERVICE LAYER
┌─────────────────────────────────────────────────────────────────────────┐
│                            coffee (Go)                                   │
│  • 12 business modules (discovery, collections, leaderboards)           │
│  • Multi-tenant with plan-based feature gating                          │
│  • 4-layer architecture: API → Service → Manager → DAO                  │
└────────────────────────────────┬────────────────────────────────────────┘
                                 │
        ┌────────────────────────┼────────────────────────┐
        ▼                        ▼                        ▼
DATA COLLECTION         EVENT STREAMING           DATA PIPELINE
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│  beat (Python)  │    │event-grpc (Go)  │    │  stir (Airflow) │
├─────────────────┤    ├─────────────────┤    ├─────────────────┤
│ • 150+ workers  │───▶│ • 26 RabbitMQ   │───▶│ • 76 DAGs       │
│ • 15+ APIs      │    │   queues        │    │ • 112 dbt models│
│ • Rate limiting │    │ • Buffered sink │    │ • CH → PG sync  │
│ • 75+ flows     │    │ • ClickHouse    │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
        │                       │                       │
        └───────────────────────┴───────────────────────┘
                                │
                                ▼
DATA LAYER
┌─────────────────────────────────────────────────────────────────────────┐
│  PostgreSQL (OLTP)  │  ClickHouse (OLAP)  │  Redis Cluster  │  S3     │
│  • Transactional    │  • Time-series      │  • Sessions     │ • Assets│
│  • 27+ tables       │  • Analytics        │  • Rate limits  │ • Logs  │
└─────────────────────────────────────────────────────────────────────────┘

ML/ANALYTICS
┌─────────────────────────────────────────────────────────────────────────┐
│  fake_follower_analysis (AWS Lambda)                                    │
│  • 10 Indic language transliteration                                    │
│  • 5-feature ensemble model                                             │
│  • SQS → Lambda → Kinesis pipeline                                      │
└─────────────────────────────────────────────────────────────────────────┘
```

---

# 1. beat - Backend API Platform

## Problem Statement

**Business Need**: Aggregate social media data (Instagram, YouTube, Shopify) for 50,000+ influencers in real-time to power analytics dashboards, leaderboards, and campaign management.

**Technical Challenges**:
1. 15+ different API providers with varying rate limits (2-850 requests/period)
2. 10M+ data points per day requiring processing
3. API credentials rotating/expiring requiring automatic fallback
4. Time-series data (follower growth, engagement trends) overwhelming PostgreSQL
5. 75+ different data collection workflows with different priorities

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────────────┐
│                      BEAT SERVICE ARCHITECTURE                           │
└─────────────────────────────────────────────────────────────────────────┘

ENTRY POINTS (3 processes)
┌──────────────────┬──────────────────┬──────────────────────────────────┐
│   server.py      │     main.py      │      main_assets.py              │
│   (FastAPI)      │   (Worker Pool)  │      (Asset Workers)             │
│   Port: 8000     │   150+ workers   │      50 workers × 100 conc       │
└────────┬─────────┴─────────┬────────┴──────────────┬───────────────────┘
         │                   │                       │
         │                   ▼                       │
         │     ┌────────────────────────────────┐   │
         │     │   WORKER POOL SYSTEM           │   │
         │     │   25 flows × 1-50 workers      │   │
         │     │   Multiprocessing + Asyncio    │   │
         │     └────────────┬───────────────────┘   │
         │                  │                       │
         │                  ▼                       │
         │     ┌────────────────────────────────┐   │
         │     │   RATE LIMITING LAYER          │   │
         │     │   3-level: Daily/Min/Handle    │   │
         │     │   Redis-backed                 │   │
         │     └────────────┬───────────────────┘   │
         │                  │                       │
         ▼                  ▼                       ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                   API INTEGRATIONS (15+ APIs)                            │
│  Instagram: Graph API, RapidAPI (6 providers)                           │
│  YouTube: Data API v3, RapidAPI (3 providers)                           │
│  Other: OpenAI, Shopify, Identity, S3                                   │
└────────────────────────────────┬────────────────────────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                        MESSAGE QUEUE (RabbitMQ)                          │
│  Event Publishing: profile_log, post_log, sentiment, relationships      │
└────────────────────────────────┬────────────────────────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                   DATABASE LAYER                                         │
│  PostgreSQL: Main entities (accounts, posts, credentials)               │
│  Redis: Rate limit state, session cache                                 │
│  S3/CloudFront: Media assets (8M images/day)                            │
└─────────────────────────────────────────────────────────────────────────┘
```

## Code Deep Dive

### 1. Worker Pool Architecture

**File**: `beat/main.py` (421 lines)

**The Problem**: Python's GIL limits CPU parallelism, but we need to maximize I/O concurrency for API calls.

**The Solution**: Hybrid architecture combining multiprocessing (bypass GIL) + asyncio (I/O concurrency).

```python
# Configuration: 25 flows with different worker counts
_whitelist = {
    # High-volume flows get more workers
    'refresh_profile_by_handle': {
        'no_of_workers': 10,        # 10 processes
        'no_of_concurrency': 5      # 5 concurrent tasks per process
    },                               # = 50 total concurrent operations

    'refresh_yt_profiles': {'no_of_workers': 10, 'no_of_concurrency': 5},
    'asset_upload_flow': {'no_of_workers': 15, 'no_of_concurrency': 5},

    # Low-volume flows get fewer workers
    'fetch_post_comments': {'no_of_workers': 1, 'no_of_concurrency': 5},

    # ... 21 more flows
}

def main():
    """Main entry point - spawns 150+ worker processes"""
    workers = []

    for flow_name, config in _whitelist.items():
        for i in range(config['no_of_workers']):
            # Each worker is a separate process
            process = multiprocessing.Process(
                target=looper,
                args=(flow_name, config['no_of_concurrency'])
            )
            process.daemon = False
            process.start()
            workers.append(process)

    # Block until all workers exit
    for worker in workers:
        worker.join()
```

**Why This Works**:
- **Multiprocessing**: Each worker is a separate OS process, bypassing Python GIL
- **Asyncio inside each process**: Event loop handles I/O concurrency within the process
- **Semaphore control**: Prevents overwhelming external APIs

### 2. Async Polling Loop

```python
def looper(flow_name: str, concurrency: int):
    """Worker process entry point"""
    # Install high-performance event loop
    uvloop.install()  # 2-4x faster than default asyncio

    asyncio.run(poller(flow_name, concurrency))

async def poller(flow_name: str, concurrency: int):
    """Async polling loop with semaphore-based concurrency control"""
    # Limit concurrent tasks per worker
    semaphore = asyncio.Semaphore(concurrency)

    while True:
        # Poll for next task (non-blocking)
        task = await poll(flow_name)

        if task:
            # Create async task (don't await - fire and forget)
            asyncio.create_task(
                perform_task(task, semaphore)
            )

        # Yield control to event loop
        await asyncio.sleep(0.1)

async def poll(flow_name: str) -> Optional[ScrapeRequestLog]:
    """Pick task from database with FOR UPDATE SKIP LOCKED"""
    query = """
        UPDATE scrape_request_log
        SET status = 'PROCESSING', picked_at = NOW()
        WHERE id = (
            SELECT id FROM scrape_request_log
            WHERE flow = :flow
              AND status = 'PENDING'
              AND (expires_at IS NULL OR expires_at > NOW())
            ORDER BY priority DESC, created_at ASC
            FOR UPDATE SKIP LOCKED  -- Critical: Skip locked rows
            LIMIT 1
        )
        RETURNING *
    """
    return await session.execute(query, {'flow': flow_name})
```

**FOR UPDATE SKIP LOCKED Explained**:
- Multiple workers poll the same queue simultaneously
- `FOR UPDATE` locks the selected row
- `SKIP LOCKED` makes other workers skip it and pick the next available row
- **Result**: No task is picked twice, no deadlocks

### 3. Three-Level Rate Limiting

**File**: `beat/server.py` (lines 312-338)

**Problem**: External APIs have different rate limits:
- Instagram Graph API: 200 calls/hour per token
- RapidAPI providers: 2-850 calls/minute
- Our own limit: 20,000 calls/day total

**Solution**: Stacked rate limiters (all must pass)

```python
from asyncio_redis_rate_limit import RateLimiter, RateSpec
from redis.asyncio import Redis

redis = Redis.from_url(REDIS_URL)

# Level 1: Daily global limit
global_limit_day = RateSpec(requests=20000, seconds=86400)

# Level 2: Per-minute global limit
global_limit_minute = RateSpec(requests=60, seconds=60)

# Level 3: Per-handle limit (prevent hammering same profile)
handle_limit = RateSpec(requests=1, seconds=1)

@app.get("/profiles/{platform}/byhandle/{handle}")
async def refresh_profile(platform: str, handle: str):
    # All three limiters must allow the request
    async with RateLimiter(
        unique_key=f"beat_daily_{platform}",
        backend=redis,
        rate_spec=global_limit_day
    ):
        async with RateLimiter(
            unique_key=f"beat_minute_{platform}",
            backend=redis,
            rate_spec=global_limit_minute
        ):
            async with RateLimiter(
                unique_key=f"beat_handle_{handle}",
                backend=redis,
                rate_spec=handle_limit
            ):
                # Execute API call
                result = await make_api_call(handle)
                return result
```

**How It Works**:
1. Request comes in for handle "cristiano"
2. Check daily limit (19,500/20,000 used) → Pass
3. Check minute limit (55/60 used) → Pass
4. Check handle limit (cristiano last called 0.5s ago) → **Block 0.5s**
5. After wait, execute API call

**Redis State**:
```
beat_daily_instagram: 19500 (TTL: 12h remaining)
beat_minute_instagram: 55 (TTL: 15s remaining)
beat_handle_cristiano: 1 (TTL: 0s remaining)
```

### 4. Credential Rotation System

**File**: `beat/credentials/manager.py`

**Problem**: API tokens expire, get rate-limited, or banned. Need automatic rotation.

```python
class CredentialManager:
    async def get_enabled_cred(self, source: str) -> Optional[Credential]:
        """Get random enabled credential that's not currently disabled"""
        query = select(Credential).where(
            Credential.source == source,
            Credential.enabled == True,
            or_(
                Credential.disabled_till.is_(None),
                Credential.disabled_till < func.now()  # TTL expired
            )
        )

        creds = await session.execute(query)
        enabled_creds = creds.scalars().all()

        if not enabled_creds:
            raise NoCredentialsAvailable(f"No credentials for {source}")

        # Random selection for load distribution
        return random.choice(enabled_creds)

    async def disable_creds(self, cred_id: int, duration: int = 3600):
        """Disable credential with TTL (auto re-enables after duration)"""
        await session.execute(
            update(Credential)
            .where(Credential.id == cred_id)
            .values(
                enabled=False,
                disabled_till=func.now() + timedelta(seconds=duration)
            )
        )
```

**Usage in API call**:
```python
async def fetch_profile(handle: str):
    max_retries = 3

    for attempt in range(max_retries):
        # Get random working credential
        cred = await credential_manager.get_enabled_cred('graphapi')

        try:
            result = await api_call(handle, cred.credentials['token'])
            return result

        except RateLimitError as e:
            # Disable this credential for 1 hour
            await credential_manager.disable_creds(cred.id, duration=3600)

            if attempt < max_retries - 1:
                continue  # Retry with different credential
            raise
```

### 5. Event Publishing System

**File**: `beat/utils/request.py` (lines 217-238)

**Context**: This is the architectural change that enabled 2.5x faster log retrieval.

**Before (Direct DB writes)**:
```python
# OLD CODE (commented out in production)
@sessionize
async def upsert_profile(profile_id, profile_log, session=None):
    # Direct write to PostgreSQL
    ts_entry = InstagramAccountTS(
        profile_id=profile_id,
        date=datetime.now().date(),
        followers=profile_log.get_metric('followers'),
        following=profile_log.get_metric('following')
    )
    session.add(ts_entry)  # ❌ Slow for high-volume time-series
```

**After (Event-driven)**:
```python
@sessionize
async def upsert_profile(profile_id, profile_log, recent_posts_log, session=None):
    # Update main account table (low volume)
    account = await upsert_insta_account(profile_log, profile_id, session)

    # Publish event for time-series data (high volume)
    await emit_profile_log_event(ProfileLog(
        platform='INSTAGRAM',
        profile_id=profile_id,
        metrics=[
            Metric('followers', profile_log.followers),
            Metric('following', profile_log.following),
            Metric('engagement_rate', profile_log.engagement_rate)
        ],
        dimensions=[
            Dimension('handle', profile_log.handle),
            Dimension('bio', profile_log.bio),
            Dimension('category', profile_log.category)
        ],
        source=profile_log.source,
        timestamp=datetime.now()
    ))

    return account

async def emit_profile_log_event(log: ProfileLog):
    """Publish event to RabbitMQ instead of DB write"""
    payload = {
        "event_id": str(uuid.uuid4()),
        "platform": log.platform,
        "profile_id": log.profile_id,
        "event_timestamp": log.timestamp.isoformat(),
        "metrics": {m.key: m.value for m in log.metrics},
        "dimensions": {d.key: d.value for d in log.dimensions}
    }

    # Publish to RabbitMQ (async, non-blocking)
    await rabbitmq.publish(
        exchange="beat.dx",
        routing_key="profile_log_events",
        body=json.dumps(payload)
    )
```

**Why This is Faster**:
```
Direct DB Write                    Event Publishing
─────────────────                  ────────────────
1. Parse API response              1. Parse API response
2. Build ProfileLog                2. Build ProfileLog
3. Begin transaction               3. Serialize to JSON
4. INSERT INTO time_series         4. Publish to queue (1-2ms)
5. Wait for disk write (50ms)      5. Done ✓
6. Commit transaction
7. Wait for fsync (20ms)
─────────────────────
Total: ~100ms                      Total: ~5ms (20x faster!)

Plus:
- No table locks
- No connection pool exhaustion
- Batch processing downstream (1000 events → 1 ClickHouse INSERT)
```

## Technical Decisions & Trade-offs

### Decision 1: Multiprocessing vs Threading

**Chosen**: Multiprocessing
**Alternative**: ThreadPoolExecutor

**Reasoning**:
- Python GIL limits CPU parallelism in threads
- API parsing (JSON decoding, data transformation) is CPU-intensive
- Separate processes = true parallelism on multi-core machines
- Trade-off: Higher memory usage (50MB per process × 150 = 7.5GB)

### Decision 2: SQL-based Task Queue vs Redis/RabbitMQ

**Chosen**: PostgreSQL with `FOR UPDATE SKIP LOCKED`
**Alternative**: Redis list or RabbitMQ queue

**Reasoning**:
- Task definitions already in PostgreSQL (scrape_request_log table)
- Need complex queries (priority, expiry, flow filtering)
- `SKIP LOCKED` prevents double-processing without distributed locks
- Trade-off: Slightly higher latency than Redis (5ms vs 1ms)

### Decision 3: Event-driven vs Direct DB Writes

**Chosen**: Event-driven (RabbitMQ → event-grpc → ClickHouse)
**Alternative**: Direct writes to PostgreSQL

**Reasoning**:
- PostgreSQL couldn't handle 10M+ writes/day for time-series
- ClickHouse columnar storage = 5x compression
- Buffered batch inserts (1000 events) amortize network cost
- Trade-off: Eventual consistency (1-5 second delay)

## Metrics & Impact

| Metric | Before | After | Improvement |
|--------|--------|-------|-------------|
| API response time | 250ms | 187ms | **25% faster** |
| Rate limit violations | 50/day | 0/day | **100% reduction** |
| Time-series query (30 days) | 30s | 12s | **2.5x faster** |
| Storage (1B logs) | 500 GB | 100 GB | **5x compression** |
| Operational costs | $1000/mo | $700/mo | **30% reduction** |
| Daily data points processed | 8M | 10M+ | **25% increase** |

## Interview Talking Points

**Q: "Tell me about the most complex system you've built"**

> "At Good Creator Co, I built beat - a distributed data collection system processing 10M+ social media data points daily. The complexity came from managing 15+ different API providers with varying rate limits, credential rotation, and needing to support 75 different data collection workflows.
>
> **Architecture**: I designed a hybrid system using multiprocessing (150+ worker processes) combined with asyncio for I/O concurrency. Each worker polls a SQL-based task queue using `FOR UPDATE SKIP LOCKED` for distributed coordination.
>
> **Rate Limiting**: I implemented a 3-level stacked rate limiting system backed by Redis - global daily limits, per-minute limits, and per-resource limits. All three must pass before making an API call.
>
> **Key Innovation**: I transitioned from direct PostgreSQL writes to an event-driven architecture. Instead of writing time-series data directly, we publish events to RabbitMQ, which are consumed by event-grpc (Go service) and batch-inserted into ClickHouse. This achieved 2.5x faster query performance and 5x storage reduction.
>
> **Results**: 25% faster API responses, 30% cost reduction, and zero rate limit violations."

---


# 2. stir - Modern Data Pipeline

## Problem Statement

**Business Need**: Transform raw social media data from beat into analytics-ready datasets for dashboards, leaderboards, and reports. Data must be fresh (< 15 minutes latency), accurate, and queryable at scale.

**Technical Challenges**:
1. Dual database strategy: ClickHouse (OLAP analytics) + PostgreSQL (OLTP operations)
2. 112 interdependent data transformations needing orchestration
3. Real-time requirements (*/15 min) vs batch requirements (daily full refresh)
4. Cross-database sync with zero-downtime (atomic table swaps)
5. Data quality validation across 10+ data sources

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    STIR DATA PLATFORM                                    │
│              Modern Data Stack: Airflow + dbt + ClickHouse              │
└─────────────────────────────────────────────────────────────────────────┘

ORCHESTRATION LAYER (Apache Airflow 2.6.3)
┌─────────────────────────────────────────────────────────────────────────┐
│  76 DAGs with Different Schedules:                                      │
│                                                                          │
│  */5 min:  dbt_recent_scl          (real-time scrape logging)          │
│  */15 min: dbt_core                (core business metrics)             │
│  */30 min: dbt_collections         (collection analytics)              │
│  Every 2h: dbt_hourly              (hourly aggregations)               │
│  Daily:    dbt_daily               (full refresh at 19:00 UTC)         │
│  Weekly:   dbt_weekly              (weekly reports)                     │
│                                                                          │
│  Operators Used:                                                         │
│  - DbtRunOperator (11): Execute dbt models                              │
│  - PythonOperator (46): Data fetching, API calls                        │
│  - ClickHouseOperator (19): Export queries                              │
│  - PostgresOperator (20): Data loading                                  │
│  - SSHOperator (18): File transfer                                      │
└────────────────────────────────┬────────────────────────────────────────┘
                                 │
                                 ▼
TRANSFORMATION LAYER (dbt 1.3.1)
┌─────────────────────────────────────────────────────────────────────────┐
│  112 dbt Models (2-layer architecture):                                 │
│                                                                          │
│  STAGING LAYER (29 models)                                              │
│  ├── beat source (13):   stg_beat_instagram_account                    │
│  │                       stg_beat_youtube_account                       │
│  │                       stg_beat_profile_log                           │
│  └── coffee source (16): stg_coffee_profile_collection                 │
│                          stg_coffee_campaign_profiles                   │
│                                                                          │
│  MART LAYER (83 models organized by domain):                            │
│  ├── Discovery (16):     mart_instagram_account                        │
│  │                       mart_youtube_account                           │
│  ├── Leaderboard (14):   mart_leaderboard                              │
│  │                       mart_time_series                               │
│  ├── Collection (13):    mart_collection_post                          │
│  ├── Genre (7):          mart_genre_trending_content                   │
│  ├── Audience (4):       mart_audience_info                            │
│  └── Profile Stats (14): mart_instagram_account_summary                │
└────────────────────────────────┬────────────────────────────────────────┘
                                 │
                                 ▼
DATA FLOW (3-Layer Sync)
┌─────────────────────────────────────────────────────────────────────────┐
│                                                                          │
│  LAYER 1: ClickHouse (172.31.28.68:9000)                               │
│  ┌────────────────────────────────────────────────────────────────┐    │
│  │  • dbt transforms data in-place (OLAP)                         │    │
│  │  • ReplacingMergeTree for efficient upserts                    │    │
│  │  • Partitioning by date (toYYYYMM) for pruning                │    │
│  │  • 112 models materialized as tables/incremental               │    │
│  └────────────────────────────────────────────────────────────────┘    │
│                            ↓                                             │
│                  INSERT INTO FUNCTION s3(...)                            │
│                            ↓                                             │
│  LAYER 2: AWS S3 (gcc-social-data bucket)                              │
│  ┌────────────────────────────────────────────────────────────────┐    │
│  │  • Path: /data-pipeline/tmp/{model}.json                       │    │
│  │  • Format: JSONEachRow (one JSON object per line)             │    │
│  │  • Decouples ClickHouse from PostgreSQL                       │    │
│  └────────────────────────────────────────────────────────────────┘    │
│                            ↓                                             │
│                  aws s3 cp s3://... /tmp/                               │
│                            ↓                                             │
│  LAYER 3: PostgreSQL (172.31.2.21:5432)                                │
│  ┌────────────────────────────────────────────────────────────────┐    │
│  │  1. COPY FROM /tmp/{model}.json (JSONB)                        │    │
│  │  2. Transform: JSONB → typed columns                           │    │
│  │  3. INSERT INTO {model}_new                                    │    │
│  │  4. ALTER TABLE {model} RENAME TO {model}_old                  │    │
│  │  5. ALTER TABLE {model}_new RENAME TO {model}                  │    │
│  │  6. DROP TABLE {model}_old                                     │    │
│  │  → Atomic swap, zero downtime                                  │    │
│  └────────────────────────────────────────────────────────────────┘    │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘
```

## Code Deep Dive

### 1. dbt Core DAG (Critical Path)

**File**: `stir/dags/dbt_core.py`

**Runs**: Every 15 minutes
**Purpose**: Update core business metrics that power real-time dashboards

```python
from airflow import DAG
from airflow_dbt_python.operators.dbt import DbtRunOperator
from datetime import datetime, timedelta

default_args = {
    'owner': 'data-engineering',
    'retries': 1,
    'retry_delay': timedelta(minutes=5),
    'execution_timeout': timedelta(minutes=60),
    'on_failure_callback': slack_alert_failure  # Slack notification
}

dag = DAG(
    dag_id='dbt_core',
    default_args=default_args,
    schedule_interval='*/15 * * * *',  # Every 15 minutes
    start_date=datetime(2023, 1, 1),
    catchup=False,                      # Don't backfill
    max_active_runs=1,                  # Prevent overlapping runs
    concurrency=1,
    dagrun_timeout=timedelta(minutes=60)
)

# Execute all models tagged with 'core'
dbt_run = DbtRunOperator(
    task_id='dbt_run_core_models',
    models='tag:core',                  # Run only 'core' tagged models
    profiles_dir='/home/airflow/.dbt',
    target='gcc_warehouse',             # ClickHouse connection
    full_refresh=False,                 # Incremental mode
    dag=dag
)
```

**What This Does**:
1. Every 15 minutes, Airflow triggers this DAG
2. dbt reads models tagged with `core` (11 models)
3. Each model executes SQL in ClickHouse
4. Models with `materialized='incremental'` only process new data
5. Slack notification if any model fails

### 2. Three-Layer Data Sync Pattern

**File**: `stir/dags/sync_leaderboard_prod.py`

**Purpose**: Sync leaderboard rankings from ClickHouse to PostgreSQL for API consumption

```python
from airflow import DAG
from airflow.providers.common.sql.operators.sql import SQLExecuteQueryOperator
from airflow.providers.ssh.operators.ssh import SSHOperator

dag = DAG('sync_leaderboard_prod', schedule_interval='15 20 * * *')  # Daily 20:15 UTC

# STEP 1: Export from ClickHouse to S3
clickhouse_export = SQLExecuteQueryOperator(
    task_id='export_leaderboard_to_s3',
    conn_id='clickhouse_gcc',
    sql="""
        INSERT INTO FUNCTION s3(
            's3://gcc-social-data/data-pipeline/tmp/leaderboard.json',
            '{aws_key}', '{aws_secret}',
            'JSONEachRow'
        )
        SELECT
            platform,
            profile_id,
            followers_rank,
            followers_rank_by_cat,
            followers_rank_by_lang,
            engagement_rate_rank,
            month
        FROM dbt.mart_leaderboard
        SETTINGS s3_truncate_on_insert=1  -- Overwrite file
    """
)

# STEP 2: Download from S3 to PostgreSQL server via SSH
ssh_download = SSHOperator(
    task_id='download_from_s3',
    ssh_conn_id='ssh_prod_pg',
    command='aws s3 cp s3://gcc-social-data/data-pipeline/tmp/leaderboard.json /tmp/'
)

# STEP 3: Load JSON into temp PostgreSQL table
pg_load_temp = SQLExecuteQueryOperator(
    task_id='load_temp_table',
    conn_id='prod_pg',
    sql="""
        -- Create temp table with JSONB column
        CREATE TEMP TABLE tmp_leaderboard (data JSONB);

        -- Load JSON file
        COPY tmp_leaderboard FROM '/tmp/leaderboard.json';
    """,
    execution_timeout=timedelta(seconds=14400)  # 4 hour timeout
)

# STEP 4: Transform JSONB to typed columns
pg_transform = SQLExecuteQueryOperator(
    task_id='transform_and_insert',
    conn_id='prod_pg',
    sql="""
        -- Insert into new table with type casting
        INSERT INTO leaderboard_new (
            platform,
            profile_id,
            followers_rank,
            followers_rank_by_cat,
            followers_rank_by_lang,
            engagement_rate_rank,
            month,
            updated_at
        )
        SELECT
            (data->>'platform')::text,
            (data->>'profile_id')::text,
            (data->>'followers_rank')::integer,
            (data->>'followers_rank_by_cat')::integer,
            (data->>'followers_rank_by_lang')::integer,
            (data->>'engagement_rate_rank')::integer,
            (data->>'month')::date,
            NOW()
        FROM tmp_leaderboard;
    """
)

# STEP 5: Atomic table swap (zero-downtime)
pg_atomic_swap = SQLExecuteQueryOperator(
    task_id='atomic_table_swap',
    conn_id='prod_pg',
    sql="""
        BEGIN;

        -- Rename current table to backup
        ALTER TABLE leaderboard
            RENAME TO leaderboard_old_bkp;

        -- Rename new table to production
        ALTER TABLE leaderboard_new
            RENAME TO leaderboard;

        -- Drop old backup (optional)
        DROP TABLE IF EXISTS leaderboard_old_bkp;

        COMMIT;
    """
)

# Define dependencies
clickhouse_export >> ssh_download >> pg_load_temp >> pg_transform >> pg_atomic_swap
```

**Why This Works**:
```
Traditional Approach (Direct PostgreSQL INSERT)
───────────────────────────────────────────────
1. SELECT FROM clickhouse
2. For each row: INSERT INTO postgres  (slow!)
3. Users see incomplete data during load
4. Long transaction = table locks

Problem: 1M rows × 50ms/insert = 13 hours!


Our Approach (Atomic Swap)
───────────────────────────────────────────────
1. Export to S3 (5 minutes)
2. Load to _new table (10 minutes)
3. Swap tables (1 second)
4. Users see complete, consistent data

Total: 15 minutes + instant cutover
Zero downtime, no partial state visible
```

### 3. dbt Incremental Model

**File**: `stir/src/gcc_social/models/marts/leaderboard/mart_time_series.sql`

**Purpose**: Daily follower/following counts per profile (time-series data)

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

    -- Get latest value per day (handles duplicates)
    argMax(followers_count, created_at) as followers,
    argMax(following_count, created_at) as following,
    argMax(posts_count, created_at) as posts,

    max(created_at) as last_updated

FROM {{ ref('stg_beat_instagram_account') }}

{% if is_incremental() %}
    -- Only process recent data (4-hour lookback for safety)
    WHERE created_at > (
        SELECT max(last_updated) - INTERVAL 4 HOUR
        FROM {{ this }}
    )
{% endif %}

GROUP BY profile_id, date
```

**How Incremental Works**:
```
First Run (No table exists)
─────────────────────────────
is_incremental() = False
→ SELECT all data from stg_beat_instagram_account
→ Create mart_time_series table


Subsequent Runs (Every 15 minutes)
──────────────────────────────────
is_incremental() = True
→ SELECT only data created_at > (max - 4 hours)
→ INSERT new rows into mart_time_series
→ ReplacingMergeTree deduplicates by (profile_id, date)


Benefits:
1. First run: Process all historical data (slow, once)
2. Later runs: Process only new data (fast, frequent)
3. 4-hour lookback: Safety buffer for late-arriving data
```

**ClickHouse Optimizations**:
```sql
-- Engine configuration
engine='ReplacingMergeTree()'  -- Deduplicates by unique_key
order_by='(profile_id, date)'  -- Clustered index for fast lookups

-- Partition by month (query optimization)
PARTITION BY toYYYYMM(date)

-- Example query that benefits:
SELECT date, followers
FROM mart_time_series
WHERE profile_id = 'cristiano'
  AND date > '2024-01-01'

-- ClickHouse optimizations applied:
1. Partition pruning: Only scan 2024-01 and later
2. Index skip: Jump directly to 'cristiano' rows
3. Columnar scan: Read only (date, followers) columns
→ Result: 100x faster than full table scan
```

### 4. Complex dbt Model with CTEs

**File**: `stir/src/gcc_social/models/marts/discovery/mart_instagram_account.sql`

**Purpose**: Enrich Instagram profiles with calculated metrics, rankings, and cross-references

```sql
{{ config(
    materialized='table',
    tags=['core', 'hourly']
) }}

-- CTE 1: Base profile data
WITH base_accounts AS (
    SELECT * FROM {{ ref('stg_beat_instagram_account') }}
    WHERE profile_id IS NOT NULL
),

-- CTE 2: Calculate post statistics (last 30 days)
post_stats AS (
    SELECT
        profile_id,
        COUNT(*) as total_posts_30d,
        AVG(likes_count) as avg_likes,
        AVG(comments_count) as avg_comments,
        AVG(views_count) as avg_views,
        SUM(likes_count) as total_likes,
        SUM(comments_count) as total_comments
    FROM {{ ref('stg_beat_instagram_post') }}
    WHERE timestamp > now() - INTERVAL 30 DAY
    GROUP BY profile_id
),

-- CTE 3: Calculate engagement rate
engagement AS (
    SELECT
        a.profile_id,
        CASE
            WHEN a.followers_count = 0 THEN 0
            ELSE ((ps.avg_likes + ps.avg_comments) / a.followers_count) * 100
        END as engagement_rate
    FROM base_accounts a
    LEFT JOIN post_stats ps USING (profile_id)
),

-- CTE 4: Assign engagement grade (A+ to F)
engagement_grades AS (
    SELECT
        profile_id,
        engagement_rate,
        CASE
            WHEN engagement_rate >= 10 THEN 'A+'
            WHEN engagement_rate >= 5 THEN 'A'
            WHEN engagement_rate >= 3 THEN 'B'
            WHEN engagement_rate >= 1 THEN 'C'
            ELSE 'F'
        END as engagement_rate_grade
    FROM engagement
),

-- CTE 5: Calculate rankings (global and by category)
rankings AS (
    SELECT
        a.profile_id,
        row_number() OVER (ORDER BY a.followers_count DESC) as followers_rank,
        row_number() OVER (PARTITION BY a.category ORDER BY a.followers_count DESC) as followers_rank_by_cat,
        row_number() OVER (PARTITION BY a.language ORDER BY a.followers_count DESC) as followers_rank_by_lang,
        row_number() OVER (ORDER BY e.engagement_rate DESC) as engagement_rank
    FROM base_accounts a
    JOIN engagement e USING (profile_id)
)

-- Final SELECT: Join everything together
SELECT
    a.*,                                    -- All base fields
    ps.total_posts_30d,
    ps.avg_likes,
    ps.avg_comments,
    ps.avg_views,
    eg.engagement_rate,
    eg.engagement_rate_grade,
    r.followers_rank,
    r.followers_rank_by_cat,
    r.followers_rank_by_lang,
    r.engagement_rank,
    NOW() as last_updated
FROM base_accounts a
LEFT JOIN post_stats ps USING (profile_id)
LEFT JOIN engagement_grades eg USING (profile_id)
LEFT JOIN rankings r USING (profile_id)
```

**Execution Plan** (ClickHouse):
```
1. Materialize base_accounts (5M rows)
2. Aggregate post_stats (scan 50M posts, group by profile_id)
3. Calculate engagement (5M rows)
4. Assign grades (5M rows)
5. Compute window functions for rankings (sort 5M rows 4 times)
6. Final join (5M rows)
7. Write to mart_instagram_account table

Total time: ~2 minutes for 5M profiles
```

## Technical Decisions & Trade-offs

### Decision 1: ClickHouse + PostgreSQL Dual Database

**Chosen**: ClickHouse for OLAP, PostgreSQL for OLTP
**Alternative**: PostgreSQL for everything

**Reasoning**:
```
PostgreSQL (OLTP - Row-based)          ClickHouse (OLAP - Columnar)
────────────────────────────────────   ────────────────────────────────────
✓ Transactional (ACID)                 ✗ Eventually consistent
✓ Strong consistency                   ✗ Weak consistency
✓ Fast single-row lookups              ✗ Slow single-row lookups
✗ Slow aggregations (billions of rows) ✓ Fast aggregations (vectorized)
✗ 1x compression                       ✓ 5-10x compression
✗ Expensive storage                    ✓ Cheap storage

Use Cases:
- User accounts, API tokens            - Time-series analytics
- Collection metadata                  - Leaderboard calculations
- Campaign configurations              - Trend detection
- Transactional operations             - Dashboard queries
```

**Our Solution**: Use both
- dbt transforms in ClickHouse (fast)
- Sync final results to PostgreSQL (API layer needs ACID)

### Decision 2: Incremental vs Full Refresh

**Chosen**: Mix of both
- Core models: Incremental (*/15 min)
- Daily models: Full refresh (daily)

**Reasoning**:
```
Incremental Materialization          Full Refresh
───────────────────────────────      ────────────────────
✓ Fast (only new data)                ✗ Slow (all data)
✓ Low compute cost                    ✗ High compute cost
✗ Can accumulate errors               ✓ Fresh start, no errors
✗ Requires unique_key logic           ✓ Simple logic

When to use incremental:
- Large tables (100M+ rows)
- Frequent updates (*/15 min)
- Append-only data

When to use full refresh:
- Small tables (< 10M rows)
- Complex joins requiring full context
- Daily reconciliation
```

### Decision 3: Atomic Table Swap vs Direct UPDATE

**Chosen**: Atomic swap with _new tables
**Alternative**: Direct UPDATE/INSERT

**Reasoning**:
```
Direct UPDATE/INSERT                 Atomic Swap
────────────────────────────────────  ────────────────────────────────────
✗ Long transaction = table locks      ✓ Short transaction (1 second)
✗ Users see partial state             ✓ All-or-nothing consistency
✗ Rollback difficult                  ✓ Easy rollback (rename back)
✗ Can fail mid-update                 ✓ Fails before visible to users

Steps:
1. Load data into leaderboard_new
2. Validate: SELECT count(*) FROM leaderboard_new
3. If good: RENAME tables (atomic)
4. If bad: DROP leaderboard_new, retry
```

## Metrics & Impact

| Metric | Before | After | Improvement |
|--------|--------|-------|-------------|
| Data latency (core metrics) | 1 hour | 15 minutes | **4x faster** |
| Data latency (leaderboards) | 4 hours | 20 minutes | **12x faster** |
| Query time (time-series 30 days) | 30s | 12s | **2.5x faster** |
| Query time (leaderboard rankings) | 15s | 2s | **7.5x faster** |
| Storage cost (1B logs) | $2000/mo | $400/mo | **80% reduction** |
| dbt run time (incremental) | N/A | 2 min | Enabled real-time |
| Data freshness | 4 hours | 15 min | **94% reduction** |
| Data quality incidents | 5/month | 0/month | **100% reduction** |

## Interview Talking Points

**Q: "Tell me about your experience with data pipelines"**

> "At Good Creator Co, I built stir - a modern data platform using Apache Airflow + dbt + ClickHouse. It processes social media data from 50,000+ influencers into analytics-ready datasets.
>
> **Architecture**: I implemented a dual-database strategy - ClickHouse for OLAP (fast aggregations on billions of rows) and PostgreSQL for OLTP (transactional API layer). dbt transforms data in ClickHouse, then we sync final results to PostgreSQL using atomic table swaps for zero-downtime updates.
>
> **Orchestration**: 76 Airflow DAGs with different schedules - every 15 minutes for core metrics, hourly for aggregations, daily for full refreshes. I used dbt tags to group models by update frequency.
>
> **Key Innovation**: Three-layer sync pattern (ClickHouse → S3 → PostgreSQL). This decouples the databases, enables JSONB parsing for flexible schemas, and atomic swaps guarantee users never see partial state.
>
> **Impact**: Reduced data latency by 50% (4 hours to 15 minutes), achieved 2.5x faster queries, and cut storage costs by 80% through columnar compression."

---


# 3. event-grpc - Event Streaming Platform

## Problem Statement

**Business Need**: Ingest and distribute real-time events from mobile apps and web applications to multiple data sinks for analytics, monitoring, and business intelligence.

**Technical Challenges**:
1. Handle 10,000+ events per second with sub-second latency
2. Route events to 26 different queues based on event type
3. Batch inserts to ClickHouse (1000 events/batch) without blocking
4. Ensure exactly-once processing with retry logic
5. Support multiple event sources (gRPC, HTTP webhooks)

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    EVENT-GRPC ARCHITECTURE                               │
└─────────────────────────────────────────────────────────────────────────┘

CLIENT APPLICATIONS
┌────────────────────────────────────────────────────────────────────────┐
│  Mobile (iOS/Android)  │  Web (React)  │  Backend Services            │
│  - gRPC protobuf       │  - gRPC-Web   │  - Server-to-server         │
└────────────────┬───────────────────────┬──────────────────────────────┘
                 │                       │
                 ▼                       ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                    EVENT-GRPC SERVICE                                    │
│  ┌──────────────────┐         ┌───────────────────────┐                │
│  │ gRPC Server      │         │ HTTP/Gin Server       │                │
│  │ Port: 8017       │         │ Port: 8019            │                │
│  │ - Dispatch()     │         │ - /heartbeat          │                │
│  │ - HealthCheck()  │         │ - /metrics            │                │
│  └────────┬─────────┘         │ - /vidooly/event      │                │
│           │                   │ - /branch/event       │                │
│           │                   │ - /webengage/event    │                │
│           │                   └───────────┬───────────┘                │
│           │                               │                            │
│           └───────────────┬───────────────┘                            │
│                           ▼                                            │
│  ┌──────────────────────────────────────────────────────────────────┐ │
│  │              WORKER POOL SYSTEM (6 Pools)                        │ │
│  │  ┌────────────────┐  ┌────────────────┐  ┌────────────────┐    │ │
│  │  │ Event Worker   │  │ Branch Worker  │  │WebEngage Worker│    │ │
│  │  │ Chan: 1000     │  │ Chan: 1000     │  │ Chan: 1000     │    │ │
│  │  │ Goroutines: N  │  │ Goroutines: N  │  │ Goroutines: N  │    │ │
│  │  └────────────────┘  └────────────────┘  └────────────────┘    │ │
│  │  ┌────────────────┐  ┌────────────────┐  ┌────────────────┐    │ │
│  │  │Vidooly Worker  │  │ Graphy Worker  │  │Shopify Worker  │    │ │
│  │  │ Chan: 1000     │  │ Chan: 1000     │  │ Chan: 1000     │    │ │
│  │  └────────────────┘  └────────────────┘  └────────────────┘    │ │
│  └──────────────────────────────────────────────────────────────────┘ │
└────────────────────────────┬─────────────────────────────────────────┘
                             │ Publish (JSON)
                             ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                    RABBITMQ MESSAGE BROKER                               │
│  ┌──────────────────────────────────────────────────────────────────┐  │
│  │ Exchanges (11):                                                  │  │
│  │  grpc_event.tx, grpc_event.dx, grpc_event_error.dx              │  │
│  │  branch_event.tx, webengage_event.dx, identity.dx               │  │
│  │  beat.dx, coffee.dx, shopify_event.dx, affiliate.dx, bigboss    │  │
│  └──────────────────────────────────────────────────────────────────┘  │
│                                                                          │
│  ┌──────────────────────────────────────────────────────────────────┐  │
│  │ Queues (26 Configured):                                          │  │
│  │  - grpc_clickhouse_event_q          (2 workers)                  │  │
│  │  - clickhouse_click_event_q         (2 workers)                  │  │
│  │  - post_log_events_q                (20 workers!)                │  │
│  │  - profile_log_events_q             (2 workers)                  │  │
│  │  - webengage_ch_event_q             (5 workers)                  │  │
│  │  - trace_log                        (2 workers, buffered)        │  │
│  │  - affiliate_orders_event_q         (2 workers, buffered)        │  │
│  │  ... and 19 more queues                                          │  │
│  └──────────────────────────────────────────────────────────────────┘  │
│                                                                          │
│  Features: Durable queues, Retry (max 2), Dead letter routing          │
└────────────────────────────┬─────────────────────────────────────────┘
                             │ Consume
                             ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                    SINKER LAYER (20+ Sinkers)                           │
│  ┌────────────────────────────────────────────────────────────────┐    │
│  │ Direct Sinkers:                  Buffered Sinkers:             │    │
│  │  - SinkEventToClickhouse         - TraceLogEventsSinker        │    │
│  │  - SinkClickEventToClickhouse    - PostLogEventsSinker         │    │
│  │  - SinkBranchEvent               - ProfileLogEventsSinker      │    │
│  │  - SinkWebengageToAPI            - SentimentLogEventsSinker    │    │
│  │  - SinkGraphyEvent               - ActivityTrackerSinker       │    │
│  └────────────────────────────────────────────────────────────────┘    │
└────────────────────────────┬─────────────────────────────────────────┘
                             │
        ┌────────────────────┼────────────────────┐
        ▼                    ▼                    ▼
┌──────────────────┐  ┌──────────────────┐  ┌──────────────────┐
│   ClickHouse     │  │   PostgreSQL     │  │  External APIs   │
│  (Analytics)     │  │ (Transactional)  │  │                  │
│  • event         │  │ • referral_event │  │ • WebEngage API  │
│  • branch_event  │  │ • user_account   │  │ • Identity Svc   │
│  • click_event   │  │                  │  │                  │
│  • post_log      │  │ Pool: 10/20      │  │ Resty client     │
│  • profile_log   │  │ Transactions     │  │ Retry logic      │
│  • trace_log     │  │                  │  │                  │
│  • + 12 more     │  │                  │  │                  │
└──────────────────┘  └──────────────────┘  └──────────────────┘
```

## Code Deep Dive

### 1. Protocol Buffers Definition

**File**: `proto/eventservice.proto` (688 lines)

**60+ Event Types Defined**:

```protobuf
syntax = "proto3";
option go_package = "proto/go/bulbulgrpc";

service EventService {
    rpc dispatch(Events) returns (Response) {}
}

message Events {
    Header header = 1;
    repeated Event events = 2;
}

message Header {
    string sessionId = 1;
    string bbDeviceId = 2;
    string userId = 3;
    string deviceId = 4;
    string androidAdvertisingId = 5;
    string clientId = 6;
    string channel = 7;
    string os = 8;
    string clientType = 9;
    string appLanguage = 10;
    string merchantId = 11;
    string ppId = 12;
    string appVersion = 13;
    string currentURL = 14;
    string utmReferrer = 15;
    string utmSource = 16;
    string utmMedium = 17;
    string utmCampaign = 18;
    // ... 6 more fields (24 total)
}

message Event {
    string eventId = 1;
    string timestamp = 2;
    
    oneof event_of {
        // User Flow (5 events)
        LaunchEvent launchEvent = 3;
        LaunchReferEvent launchReferEvent = 4;
        PageOpenedEvent pageOpenedEvent = 5;
        
        // Widget Events (10 events)
        WidgetViewEvent widgetViewEvent = 10;
        WidgetCtaClickedEvent widgetCtaClickedEvent = 11;
        
        // Commerce Events (10 events)
        AddToCartEvent addToCartEvent = 20;
        InitiatePurchaseEvent initiatePurchaseEvent = 21;
        CompletePurchaseEvent completePurchaseEvent = 22;
        PurchaseFailedEvent purchaseFailedEvent = 23;
        
        // Streaming Events (5 events)
        StreamEnterEvent streamEnterEvent = 30;
        SocialStreamEnterEvent socialStreamEnterEvent = 31;
        
        // ... 30+ more event types
    }
}

// Example event structure
message AddToCartEvent {
    string productId = 1;
    string variantId = 2;
    double price = 3;
    string currency = 4;
    int32 quantity = 5;
    string merchantId = 6;
}
```

### 2. Worker Pool Pattern

**File**: `eventworker/eventworker.go`

```go
var (
    eventWrapperChannel chan bulbulgrpc.Events
    channelInit         sync.Once
)

func GetChannel(config config.Config) chan bulbulgrpc.Events {
    channelInit.Do(func() {
        // Create buffered channel (1000 capacity)
        eventWrapperChannel = make(chan bulbulgrpc.Events,
            config.EVENT_WORKER_POOL_CONFIG.EVENT_BUFFERED_CHANNEL_SIZE)

        // Initialize worker pool
        initWorkerPool(config,
            config.EVENT_WORKER_POOL_CONFIG.EVENT_WORKER_POOL_SIZE,
            eventWrapperChannel)
    })
    return eventWrapperChannel
}

func initWorkerPool(config config.Config, poolSize int,
    eventChannel <-chan bulbulgrpc.Events) {
    
    for i := 0; i < poolSize; i++ {
        // safego wraps goroutines with panic recovery
        safego.GoNoCtx(func() {
            worker(config, eventChannel)
        })
    }
}

func worker(config config.Config, eventChannel <-chan bulbulgrpc.Events) {
    for e := range eventChannel {
        rabbitConn := rabbit.Rabbit(config)
        
        // Group events by type (reflection-based)
        eventsGrouped := make(map[string][]*bulbulgrpc.Event)
        for _, evt := range e.Events {
            eventName := fmt.Sprintf("%v", reflect.TypeOf(evt.GetEventOf()))
            eventName = strings.ReplaceAll(eventName, "*bulbulgrpc.Event_", "")
            
            // Timestamp correction (future timestamps → now)
            et := time.Unix(transformerToInt64(evt.Timestamp)/1000, 0)
            if et.Unix() > time.Now().Unix() {
                evt.Timestamp = strconv.FormatInt(time.Now().Unix()*1000, 10)
            }
            
            eventsGrouped[eventName] = append(eventsGrouped[eventName], evt)
        }
        
        // Publish grouped events to RabbitMQ
        for eventName, events := range eventsGrouped {
            e := bulbulgrpc.Events{Header: e.Header, Events: events}
            if b, err := protojson.Marshal(&e); err == nil {
                rabbitConn.Publish("grpc_event.tx", eventName, b, map[string]interface{}{})
            }
        }
    }
}
```

**Why This Pattern Works**:
```
Single Channel (1000 capacity)
    │
    ├─→ Worker 1 (goroutine) ──→ Process events
    ├─→ Worker 2 (goroutine) ──→ Process events
    ├─→ Worker 3 (goroutine) ──→ Process events
    └─→ ... N workers

Benefits:
1. Buffered channel prevents blocking gRPC handlers
2. Multiple workers process events concurrently
3. Panic recovery prevents crash from bad events
4. Event grouping reduces RabbitMQ calls
```

### 3. RabbitMQ Consumer Setup

**File**: `main.go` (lines 382-400)

```go
// Example: Buffered consumer for high-volume logs
traceLogChan := make(chan interface{}, 10000)  // 10K buffer

traceLogEventConsumerCfg := rabbit.RabbitConsumerConfig{
    QueueName:            "trace_log",
    Exchange:             "identity.dx",
    RoutingKey:           "trace_log",
    RetryOnError:         true,
    ErrorExchange:        &errorExchange,
    ErrorRoutingKey:      &errorRoutingKey,
    ConsumerCount:        2,                              // 2 concurrent consumers
    BufferChan:           traceLogChan,                   // Buffered channel
    BufferedConsumerFunc: sinker.BufferTraceLogEvent,    // Buffer function
}

rabbit.Rabbit(config).InitConsumer(traceLogEventConsumerCfg)

// Start background sinker for batch processing
go sinker.TraceLogEventsSinker(traceLogChan)
```

**All 26 Consumer Configurations**:

```go
func main() {
    config := config.New()
    
    // High-volume: 20 workers for post logs
    postLogChan := make(chan interface{}, 10000)
    postLogCfg := rabbit.RabbitConsumerConfig{
        QueueName:            "post_log_events_q",
        Exchange:             "beat.dx",
        RoutingKey:           "post_log_events",
        RetryOnError:         true,
        ConsumerCount:        20,  // Highest concurrency!
        BufferChan:           postLogChan,
        BufferedConsumerFunc: sinker.BufferPostLogEvents,
    }
    rabbit.Rabbit(config).InitConsumer(postLogCfg)
    go sinker.PostLogEventsSinker(postLogChan)
    
    // Medium-volume: 5 workers for WebEngage
    webEngageChan := make(chan interface{}, 10000)
    webEngageCfg := rabbit.RabbitConsumerConfig{
        QueueName:     "webengage_ch_event_q",
        Exchange:      "webengage_event.dx",
        RoutingKey:    "webengage_ch_event",
        ConsumerCount: 5,
        BufferChan:    webEngageChan,
        BufferedConsumerFunc: sinker.BufferWebengageEvents,
    }
    rabbit.Rabbit(config).InitConsumer(webEngageCfg)
    go sinker.WebengageEventsSinker(webEngageChan)
    
    // Low-volume: 2 workers for profile logs
    profileLogChan := make(chan interface{}, 10000)
    profileLogCfg := rabbit.RabbitConsumerConfig{
        QueueName:     "profile_log_events_q",
        Exchange:      "beat.dx",
        RoutingKey:    "profile_log_events",
        ConsumerCount: 2,
        BufferChan:    profileLogChan,
        BufferedConsumerFunc: sinker.BufferProfileLogEvents,
    }
    rabbit.Rabbit(config).InitConsumer(profileLogCfg)
    go sinker.ProfileLogEventsSinker(profileLogChan)
    
    // ... 23 more consumer configurations
}
```

### 4. Buffered Sinker Pattern

**File**: `sinker/profile_log_sinker.go`

```go
func BufferProfileLogEvents(delivery amqp.Delivery, c chan interface{}) bool {
    """Buffer function - adds event to channel without blocking"""
    var profileLog map[string]interface{}
    if err := json.Unmarshal(delivery.Body, &profileLog); err == nil {
        c <- profileLog  // Non-blocking send to buffered channel
        return true
    }
    return false
}

func ProfileLogEventsSinker(c chan interface{}) {
    """Batch processor - runs in separate goroutine"""
    ticker := time.NewTicker(5 * time.Second)  // Flush every 5 seconds
    batch := []model.ProfileLogEvent{}
    
    for {
        select {
        case event := <-c:
            // Parse and add to batch
            profileLog := parseProfileLog(event)
            batch = append(batch, profileLog)
            
            // Flush if batch full (1000 events)
            if len(batch) >= 1000 {
                flushBatch(batch)
                batch = []model.ProfileLogEvent{}
            }
            
        case <-ticker.C:
            // Periodic flush (even if batch not full)
            if len(batch) > 0 {
                flushBatch(batch)
                batch = []model.ProfileLogEvent{}
            }
        }
    }
}

func flushBatch(batch []model.ProfileLogEvent) {
    """Single INSERT for entire batch"""
    result := clickhouse.Clickhouse(config.New(), nil).Create(batch)
    if result.Error != nil {
        log.Printf("Batch insert failed: %v", result.Error)
    }
}
```

**Why Buffering Works**:
```
WITHOUT BUFFERING               WITH BUFFERING
──────────────────               ──────────────
Event 1 → INSERT (5ms)          Event 1   ↘
Event 2 → INSERT (5ms)          Event 2    ├─ Buffer (0.1ms each)
Event 3 → INSERT (5ms)          Event 3    │
...                             ...        │
Event 1000 → INSERT (5ms)       Event 1000↗
                                    ↓
Total: 5000ms                   Flush 1000 → INSERT (50ms)
                                Total: 100ms + 50ms = 150ms
                                
33x faster!
```

### 5. RabbitMQ Retry Logic

**File**: `rabbit/rabbit.go`

```go
func (rabbit *RabbitConnection) Consume(cfg RabbitConsumerConfig) error {
    msgs, _ := channel.Consume(
        cfg.QueueName,
        "",    // consumer tag (auto-generated)
        false, // auto-ack (we handle manually)
        false, // exclusive
        false, // no-local
        false, // no-wait
        nil,   // args
    )
    
    for msg := range msgs {
        // Call consumer function
        listenerResponse := consumerFunc(msg)
        
        if listenerResponse {
            // Success: acknowledge
            msg.Ack(false)
            
        } else if !cfg.RetryOnError {
            // Failure (no retry): send to error queue
            publishToErrorQueue(msg, cfg.ErrorExchange, cfg.ErrorRoutingKey)
            msg.Ack(false)
            
        } else {
            // Failure (with retry): check retry count
            retryCount := getRetryCount(msg.Headers)
            
            if retryCount >= 2 {
                // Max retries exceeded: dead letter
                publishToErrorQueue(msg, cfg.ErrorExchange, cfg.ErrorRoutingKey)
                msg.Ack(false)
            } else {
                // Republish with incremented retry count
                msg.Headers["x-retry-count"] = retryCount + 1
                republish(msg, cfg.QueueName)
                msg.Ack(false)
            }
        }
    }
}

func getRetryCount(headers map[string]interface{}) int {
    if val, ok := headers["x-retry-count"]; ok {
        if count, ok := val.(int); ok {
            return count
        }
    }
    return 0
}
```

**Retry Flow**:
```
Event Processing Attempt
    │
    ├─ SUCCESS → ACK → Done
    │
    └─ FAILURE
        │
        ├─ No Retry Config → Error Queue → ACK
        │
        └─ Retry Enabled
            │
            ├─ Attempt 1 (count=0) → Republish (count=1)
            │
            ├─ Attempt 2 (count=1) → Republish (count=2)
            │
            └─ Attempt 3 (count=2) → Error Queue → ACK
```

## Technical Decisions & Trade-offs

### Decision 1: gRPC vs REST

**Chosen**: gRPC with Protocol Buffers
**Alternative**: REST with JSON

**Reasoning**:
```
gRPC + Protobuf                    REST + JSON
──────────────────────────────    ────────────────────────────
✓ 5-10x smaller payload            ✗ Larger payload size
✓ Faster serialization             ✗ Slower JSON parsing
✓ Strongly typed schema            ✗ Runtime type checking
✓ Bi-directional streaming         ✗ Request-response only
✗ Complex tooling setup            ✓ Simple curl/postman

Use Case Fit:
- High-volume events (10K+/sec)
- Mobile apps (bandwidth sensitive)
- Strong contract (60+ event types)
```

### Decision 2: Buffered vs Direct Sinkers

**Chosen**: Mix of both
- Direct: Low-volume, real-time (branch events, error logs)
- Buffered: High-volume, batch-friendly (post logs, profile logs)

**Reasoning**:
```
Direct Sinkers                    Buffered Sinkers
──────────────────────────────    ────────────────────────────
✓ Immediate visibility             ✗ 1-5 second delay
✓ Simpler logic                    ✗ Complex batching logic
✗ High DB connection usage         ✓ Low DB connections
✗ Network overhead per event       ✓ Amortized network cost

Decision Criteria:
- Volume > 1000 events/hour → Buffered
- Real-time required → Direct
- ClickHouse destination → Buffered (loves batches)
```

### Decision 3: 26 Separate Queues vs Single Queue

**Chosen**: 26 separate queues (one per event type)
**Alternative**: Single queue with routing logic

**Reasoning**:
```
Separate Queues                   Single Queue
──────────────────────────────    ────────────────────────────
✓ Independent scaling              ✗ Single bottleneck
✓ Isolated failures                ✗ One failure blocks all
✓ Per-queue monitoring             ✗ Aggregated metrics
✗ 26 consumer processes            ✓ Single consumer
✗ More complex setup               ✓ Simpler architecture

Our Scale:
- post_log_events_q: 20 workers (highest volume)
- profile_log_events_q: 2 workers
- Independent scaling matches workload
```

## Metrics & Impact

| Metric | Value |
|--------|-------|
| **Event Throughput** | 10,000+ events/sec |
| **Latency (P99)** | < 100ms (dispatch to ClickHouse) |
| **Worker Pools** | 6 (configurable) |
| **RabbitMQ Queues** | 26 |
| **Total Consumer Workers** | 70+ |
| **Batch Size** | 1000 events |
| **Flush Interval** | 5 seconds |
| **Retry Attempts** | 2 max |
| **ClickHouse Tables** | 18+ |

## Interview Talking Points

**Q: "Tell me about a high-throughput event streaming system you built"**

> "At Good Creator Co, I built event-grpc - a gRPC-based event streaming platform handling 10,000+ events per second. The system ingests events from mobile apps and web applications, then distributes them to 26 RabbitMQ queues for persistence in ClickHouse and PostgreSQL.
>
> **Architecture**: 6 worker pools with buffered channels (1000 capacity each). Events are grouped by type using reflection, then published to RabbitMQ. I implemented 26 consumer configurations with different worker counts based on volume - post logs get 20 workers, profile logs get 2.
>
> **Buffered Sinker Pattern**: For high-volume events, I used a two-stage pattern: buffer function adds to channel (non-blocking), then batch processor flushes every 5 seconds or 1000 events. This achieved 33x better performance than individual inserts.
>
> **Resilience**: Retry logic with max 2 attempts, dead letter queues for failed events, panic recovery wrappers on all goroutines, and auto-reconnect for database and queue connections.
>
> **Results**: Sub-100ms P99 latency, supports billions of events, zero data loss with retry mechanism."

---


# 4. fake_follower_analysis - ML-Powered Fake Detection

## Problem Statement

**Business Need**: Detect fake/bot followers on Instagram to provide accurate influencer analytics for brand partnerships.

**Technical Challenges**:
1. 10 Indic scripts (Hindi, Bengali, Tamil, Telugu, etc.) require transliteration
2. 35,183 Indian baby names database for validation
3. 100K+ followers per influencer need analysis
4. Serverless Lambda with cold start constraints
5. Multi-feature ensemble model for 85% accuracy

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────────────┐
│              FAKE FOLLOWER ANALYSIS PIPELINE                             │
└─────────────────────────────────────────────────────────────────────────┘

DATA EXTRACTION
┌─────────────────────────────────────────────────────────────────────────┐
│  ClickHouse → Complex CTE Query                                          │
│  ├── handles AS (Load creator handles from S3 CSV)                      │
│  ├── profile_ids AS (Map handles → profile_id)                          │
│  ├── follower_data AS (Historical: stg_beat_profile_relationship_log)   │
│  ├── follower_events_data AS (Real-time: profile_relationship_log_events)│
│  └── data AS (UNION ALL historical + real-time)                         │
│                                                                          │
│  Output: S3 → gcc-social-data/temp/{date}_creator_followers.json        │
│  Format: JSONEachRow (one follower per line)                            │
└────────────────────────────────┬────────────────────────────────────────┘
                                 │
                                 ▼
MESSAGE DISTRIBUTION
┌─────────────────────────────────────────────────────────────────────────┐
│  push.py - Multiprocessing Pipeline                                     │
│  ├── Download S3 file                                                   │
│  ├── Batch processing (10,000 lines/chunk)                              │
│  ├── 8 parallel workers (multiprocessing.Pool)                          │
│  └── SQS batch send (10 messages/call)                                  │
│                                                                          │
│  SQS Queue: creator_follower_in (eu-north-1)                            │
│  ├── MaximumMessageSize: 256 KB                                         │
│  ├── MessageRetentionPeriod: 4 days                                     │
│  └── VisibilityTimeout: 30 seconds                                      │
└────────────────────────────────┬────────────────────────────────────────┘
                                 │ Event trigger
                                 ▼
LAMBDA PROCESSING
┌─────────────────────────────────────────────────────────────────────────┐
│  AWS Lambda (ECR Container - Python 3.10)                               │
│  ├── Handler: fake.handler(event, context)                              │
│  ├── ML Models: 10 Indic language HMMs                                  │
│  ├── Name Database: 35,183 Indian names                                 │
│  ├── Transliteration: indictrans library                                │
│  └── Processing: 50-100ms per follower                                  │
│                                                                          │
│  5-Feature Ensemble Model:                                              │
│  ├── 1. Non-Indic language detection (0/1)                              │
│  ├── 2. >4 digits in handle (0/1)                                       │
│  ├── 3. Special char mismatch logic (0/1/2)                             │
│  ├── 4. Handle-name similarity (0-100)                                  │
│  └── 5. Indian name database match (0-100)                              │
│                                                                          │
│  Output: 19 fields per follower (features + final score)                │
└────────────────────────────────┬────────────────────────────────────────┘
                                 │ Kinesis put_record
                                 ▼
RESULTS STREAMING
┌─────────────────────────────────────────────────────────────────────────┐
│  Kinesis Stream: creator_out (ON_DEMAND)                                │
│  ├── PartitionKey: follower_handle                                      │
│  ├── Region: ap-south-1                                                 │
│  └── Auto-scaling shards                                                │
│                                                                          │
│  pull.py - Multi-shard parallel read                                    │
│  └── Output: {date}_creator_followers_final_fake_analysis.json          │
└─────────────────────────────────────────────────────────────────────────┘
```

## Code Deep Dive

### 1. ML Detection Pipeline

**File**: `fake.py` (385 lines)

```python
def handler(event, context):
    """AWS Lambda entry point"""
    # Parse SQS message
    for record in event['Records']:
        body = json.loads(record['body'])
        follower_handle = body['follower_handle']
        follower_full_name = body['follower_full_name']
        
        # Execute ML pipeline
        result = model(follower_handle, follower_full_name)
        
        # Stream to Kinesis
        kinesis.put_record(
            StreamName='creator_out',
            Data=json.dumps(result),
            PartitionKey=follower_handle
        )

def model(follower_handle, follower_full_name):
    """5-stage ML pipeline with ensemble scoring"""
    
    # STAGE 1: Unicode symbol normalization
    symbolic_name = symbol_name_convert(follower_full_name)
    # Example: "𝓐𝓵𝓲𝓬𝓮" → "Alice"
    
    # STAGE 2: Language detection
    lang_fake = check_lang_other_than_indic(symbolic_name)
    # Pattern: r'[Α-Ωα-ωԱ-Ֆა-ჰ一-鿿가-힣]+'
    # Detects: Greek, Armenian, Georgian, Chinese, Korean
    
    # STAGE 3: Indic script transliteration
    if is_indic_script(symbolic_name):
        language = detect_language(symbolic_name)  # hin, ben, tam, etc.
        transliterated_name = transliterate(symbolic_name, language)
        # Example: "राहुल" → "Rahul" (Hindi → English)
    else:
        transliterated_name = symbolic_name
    
    # STAGE 4: Final decoding
    decoded_name = uni_decode(transliterated_name)
    # Example: "Ràhul" → "Rahul" (remove diacritics)
    
    # STAGE 5: Clean handle and name
    cleaned_handle = clean_handle(follower_handle)
    cleaned_name = clean_handle(decoded_name)
    # Rules: [_\-.] → space, remove digits, lowercase
    
    # FEATURE EXTRACTION
    
    # Feature 1: Non-Indic language (1 = FAKE)
    fake_real_based_on_lang = lang_fake
    
    # Feature 2: Digit count (>4 = FAKE)
    number_handle = count_numerical_digits(follower_handle)
    number_more_than_4_handle = 1 if number_handle > 4 else 0
    
    # Feature 3: Special character logic
    chhitij_logic = process(follower_handle, cleaned_handle, cleaned_name)
    # 0 = REAL (good match with special chars)
    # 1 = FAKE (special chars but poor match)
    # 2 = INCONCLUSIVE (no special chars)
    
    # Feature 4: Fuzzy similarity (RapidFuzz)
    similarity_score = generate_similarity_score(cleaned_handle, cleaned_name)
    # Range: 0-100 (higher = more similar)
    
    # Feature 5: Indian name database match
    indian_name_score = check_indian_names(cleaned_name)
    # Range: 0-100 (match against 35,183 names)
    
    # ENSEMBLE SCORING
    
    # Binary classification
    process1_ = process1(fake_real_based_on_lang, 
                         number_more_than_4_handle, 
                         chhitij_logic)
    # 0 = REAL, 1 = FAKE, 2 = INCONCLUSIVE
    
    # Final weighted score
    final_ = final(fake_real_based_on_lang, 
                   similarity_score,
                   number_more_than_4_handle, 
                   chhitij_logic)
    # 0.0 = Definitely REAL
    # 0.33 = Weak FAKE indicator
    # 1.0 = Definitely FAKE
    
    return {
        # Processed names
        "symbolic_name": symbolic_name,
        "transliterated_follower_name": transliterated_name,
        "decoded_name": decoded_name,
        "cleaned_handle": cleaned_handle,
        "cleaned_name": cleaned_name,
        
        # Binary features
        "fake_real_based_on_lang": fake_real_based_on_lang,
        "chhitij_logic": chhitij_logic,
        "number_handle": number_handle,
        "number_more_than_4_handle": number_more_than_4_handle,
        
        # Similarity scores
        "similarity_score": similarity_score,
        "indian_name_score": indian_name_score,
        
        # Ensemble outputs
        "process1_": process1_,
        "final_": final_
    }
```

### 2. Fuzzy String Matching

```python
def generate_similarity_score(handle, name):
    """
    RapidFuzz-based weighted ensemble
    
    Algorithm:
    1. Generate name permutations (max 24 for 4 words)
    2. For each permutation, calculate 3 metrics:
       - partial_ratio: Substring matching (weight: 2x)
       - token_sort_ratio: Order-invariant matching
       - token_set_ratio: Subset matching
    3. Weighted average: (2×partial + sort + set) / 4
    4. Return max score across all permutations
    """
    from itertools import permutations
    from rapidfuzz import fuzz
    
    name_parts = name.split()
    
    if len(name_parts) <= 4:
        # Generate all permutations for better matching
        # "John Doe" → ["John Doe", "Doe John"]
        name_permutations = [' '.join(p) for p in permutations(name_parts)]
    else:
        # Too many permutations (120 for 5 words), use as-is
        name_permutations = [name]
    
    similarity_score = -1
    
    for name_variant in name_permutations:
        # Substring matching (most important)
        partial_ratio = fuzz.partial_ratio(handle, name_variant)
        
        # Order-invariant: "john doe" == "doe john"
        token_sort_ratio = fuzz.token_sort_ratio(handle, name_variant)
        
        # Subset matching: ignores extra tokens
        token_set_ratio = fuzz.token_set_ratio(handle, name_variant)
        
        # Weighted average (partial_ratio weighted 2x)
        score = (2 * partial_ratio + token_sort_ratio + token_set_ratio) / 4
        
        similarity_score = max(similarity_score, score)
    
    return similarity_score

# Example Results:
# handle="john_doe", name="John Doe" → ~95
# handle="johndoe123", name="John Doe" → ~85
# handle="xyz_random", name="John Doe" → ~20
```

### 3. Indian Name Database Matching

```python
# Global: Load 35,183 names at module import (Lambda warm start)
namess = pd.read_csv('baby_names_.csv')['Baby Names'].str.lower().tolist()

def check_indian_names(name):
    """
    Fuzzy match against Indian baby names database
    
    Process:
    1. Split name into first_name + last_name
    2. For each part, scan entire database (35,183 names)
    3. Use same weighted formula as similarity scoring
    4. Return maximum score found
    
    Special Cases:
    - Name < 2 chars → Return 1 (too short, likely fake)
    - Last name < 2 chars → Ignore last name
    """
    if len(name) < 2:
        return 1  # Too short
    
    name_parts = name.split()
    first_name = name_parts[0]
    last_name = name_parts[1] if len(name_parts) >= 2 else None
    
    similarity_score = 0
    
    # Match first name against all 35,183 names
    for db_name in namess:
        score = (
            2 * fuzz.ratio(db_name, first_name) +
            fuzz.token_sort_ratio(db_name, first_name) +
            fuzz.token_set_ratio(db_name, first_name)
        ) / 4
        similarity_score = max(similarity_score, score)
    
    # Match last name if valid
    if last_name and len(last_name) >= 2:
        for db_name in namess:
            score = (
                2 * fuzz.ratio(db_name, last_name) +
                fuzz.token_sort_ratio(db_name, last_name) +
                fuzz.token_set_ratio(db_name, last_name)
            ) / 4
            similarity_score = max(similarity_score, score)
    
    return similarity_score

# Example Results:
# "Rahul Kumar" → ~95 (both match database)
# "Alex Smith" → ~60 (partial match)
# "xyz abc" → ~10 (no match, likely fake)
```

### 4. Indic Language Transliteration

```python
from indictrans import Transliterator

# Supported: Hindi, Bengali, Tamil, Telugu, Kannada, Malayalam,
#            Gujarati, Punjabi, Odia, Urdu (10 languages)

def detect_language(word):
    """
    Character-by-character language identification
    
    Uses Unicode ranges:
    - Hindi (Devanagari): \u0900-\u097F
    - Bengali: \u0980-\u09FF
    - Tamil: \u0B80-\u0BFF
    - Telugu: \u0C00-\u0C7F
    - ... (10 languages total)
    """
    char_to_lang = {
        'अ': 'hin', 'आ': 'hin', 'इ': 'hin',  # Hindi (77 chars)
        'ਅ': 'pan', 'ਆ': 'pan', 'ਇ': 'pan',  # Punjabi (61 chars)
        'অ': 'ben', 'আ': 'ben', 'ই': 'ben',  # Bengali (65 chars)
        'அ': 'tam', 'ஆ': 'tam', 'இ': 'tam',  # Tamil (62 chars)
        # ... 400+ character mappings
    }
    
    for char in word:
        if char in char_to_lang:
            return char_to_lang[char]
    
    return None

def transliterate(word, language):
    """
    ML-based transliteration using pre-trained HMM models
    
    Model Structure (per language):
    - coef_.npy: HMM coefficient matrix
    - classes.npy: Output character mapping
    - intercept_init_.npy: Initial state probabilities
    - intercept_trans_.npy: Transition probabilities
    - intercept_final_.npy: Final state probabilities
    - sparse.vec: Feature vocabulary
    """
    if language == 'hin':
        # Hindi: Use custom vowel/consonant mappings
        return process_word_hindi(word)
    else:
        # Other languages: Use indictrans ML models
        trn = Transliterator(source=language, target='eng', decode='viterbi')
        return trn.transform(word)

# Example Transliterations:
# "राहुल" (Hindi) → "Rahul"
# "মুকেশ" (Bengali) → "Mukesh"
# "ராகுல்" (Tamil) → "Ragul"
# "రాహుల్" (Telugu) → "Rahul"
```

### 5. Ensemble Scoring Logic

```python
def process1(fake_real_based_on_lang, number_more_than_4_handle, chhitij_logic):
    """
    Binary feature combination classifier
    
    Decision Tree:
    ├── Non-Indic language detected? → 1 (FAKE)
    ├── >4 digits in handle? → 1 (FAKE)
    ├── Special char mismatch (chhitij=1)? → 1 (FAKE)
    ├── No special chars (chhitij=2)? → 2 (INCONCLUSIVE)
    └── Otherwise → 0 (REAL)
    """
    if fake_real_based_on_lang:
        return 1  # Non-Indic script = FAKE
    if number_more_than_4_handle:
        return 1  # Too many digits = FAKE
    if chhitij_logic == 1:
        return 1  # Special chars but poor match = FAKE
    elif chhitij_logic == 2:
        return 2  # No special chars = INCONCLUSIVE
    return 0  # Default: REAL

def final(fake_real_based_on_lang, similarity_score,
          number_more_than_4_handle, chhitij_logic):
    """
    Weighted final score (confidence level)
    
    Scoring Rules:
    ├── Non-Indic language? → 1.0 (100% FAKE)
    ├── Similarity 0-40? → 0.33 (33% confidence FAKE)
    ├── >4 digits in handle? → 1.0 (100% FAKE)
    ├── Special char mismatch? → 1.0 (100% FAKE)
    ├── No special chars? → 0.0 (REAL)
    └── Otherwise → 0.0 (REAL)
    
    Output Interpretation:
    - 0.0  = Definitely REAL (high confidence)
    - 0.33 = Weak FAKE signal (low confidence)
    - 1.0  = Definitely FAKE (high confidence)
    """
    # Strong FAKE indicators
    if fake_real_based_on_lang:
        return 1.0  # Foreign script = bot
    if number_more_than_4_handle:
        return 1.0  # Random digits = bot
    if chhitij_logic == 1:
        return 1.0  # Intentional separators but fake name = bot
    
    # Weak FAKE indicator
    if 0 < similarity_score <= 40:
        return 0.33  # Poor handle-name match = suspicious
    
    # REAL indicators
    if chhitij_logic == 2:
        return 0.0  # No separators = likely real
    
    return 0.0  # Default: real
```

## Technical Decisions & Trade-offs

### Decision 1: Lambda vs EC2

**Chosen**: AWS Lambda with ECR container
**Alternative**: EC2 instance

**Reasoning**:
```
Lambda + ECR                      EC2
──────────────────────────────    ────────────────────────────
✓ Auto-scaling (0 to 1000s)       ✗ Manual scaling
✓ Pay per execution               ✗ 24/7 cost
✓ No infrastructure management    ✗ OS patching, monitoring
✗ Cold start (1-3 seconds)        ✓ Always warm
✗ 10GB memory limit               ✓ Any instance size

Use Case Fit:
- Batch processing (not real-time)
- Variable load (daily jobs)
- Container includes 10 ML models (500MB)
```

### Decision 2: Ensemble vs Single Model

**Chosen**: 5-feature ensemble with rule-based logic
**Alternative**: Single ML classifier (Random Forest, XGBoost)

**Reasoning**:
```
Ensemble (Rule-Based)             Single ML Model
──────────────────────────────    ────────────────────────────
✓ Explainable (each feature)      ✗ Black box
✓ No training data needed         ✗ Requires labeled data
✓ Domain knowledge encoded        ✗ Learns from data
✓ Easy to debug/tune              ✗ Hard to interpret
✗ May miss complex patterns       ✓ Learns patterns

Our Context:
- Limited labeled data (no training set)
- Need explainability for clients
- Domain expertise available (Indic scripts, Indian names)
- 85% accuracy good enough
```

### Decision 3: 35K Name Database vs API

**Chosen**: Local CSV (35,183 names)
**Alternative**: External name validation API

**Reasoning**:
```
Local Database                    External API
──────────────────────────────    ────────────────────────────
✓ No network latency              ✗ 50-200ms per call
✓ No rate limits                  ✗ Cost per API call
✓ Works offline                   ✗ Depends on uptime
✗ Static (no updates)             ✓ Dynamic updates
✗ 287KB container size            ✓ No container bloat

Performance:
- Local: 10-50ms to scan 35K names
- API: 50-200ms + network overhead
- For 100K followers: 1-2 hours vs 5-10 hours
```

## Metrics & Impact

| Metric | Value |
|--------|-------|
| **Accuracy** | 85% |
| **Languages Supported** | 10 Indic scripts |
| **Name Database** | 35,183 Indian names |
| **Processing Time** | 50-100ms per follower |
| **Throughput** | 10-20 followers/sec (single Lambda) |
| **Confidence Levels** | 3 (0.0, 0.33, 1.0) |
| **Features** | 5 independent heuristics |
| **Output Fields** | 19 per follower |

## Interview Talking Points

**Q: "Tell me about an ML system you built"**

> "At Good Creator Co, I built a fake follower detection system for Instagram influencers using an ensemble ML approach. The challenge was handling 10 Indic languages and 100K+ followers per influencer.
>
> **Architecture**: AWS Lambda with ECR container packaging 10 pre-trained HMM models for script transliteration. The pipeline extracts followers from ClickHouse, distributes via SQS (8 parallel workers), processes in Lambda, and streams results via Kinesis.
>
> **ML Approach**: 5-feature ensemble combining:
> 1. Language detection (Greek/Chinese/Korean scripts = fake)
> 2. Digit count in handle (>4 = fake)
> 3. Special character logic (separators should match name)
> 4. Fuzzy string matching (RapidFuzz with weighted scoring)
> 5. Indian name database (35,183 names for validation)
>
> **Key Innovation**: Used indictrans library with HMM-based transliteration for 10 Indic scripts. Example: 'राहुल' (Hindi) → 'Rahul' (English), then fuzzy match against handle and name database.
>
> **Results**: 85% accuracy, processes 100K followers in 1-2 hours, provides 3 confidence levels (0.0/0.33/1.0) for client interpretation."

---

# 5. coffee - Multi-Tenant SaaS Platform

## Problem Statement

**Business Need**: Build a unified REST API for influencer discovery, profile collections, leaderboards, and analytics serving 1000+ B2B clients.

**Technical Challenges**:
1. Multi-tenancy with partner/account isolation
2. Dual database strategy (PostgreSQL OLTP + ClickHouse OLAP)
3. 12 business modules with consistent architecture
4. Plan-based feature gating (FREE, SAAS, PAID)
5. Transaction management across HTTP requests

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    COFFEE SERVICE ARCHITECTURE                           │
└─────────────────────────────────────────────────────────────────────────┘

HTTP REQUEST
    │
    ▼
┌─────────────────────────────────────────────────────────────────────────┐
│              10-LAYER MIDDLEWARE PIPELINE                                │
│  1. gin.Recovery()           ← Panic recovery                           │
│  2. sentrygin.New()          ← Error tracking                           │
│  3. cors.New()               ← CORS headers                             │
│  4. ApplicationContext       ← Header extraction (x-bb-*)               │
│  5. ServiceSession           ← Transaction management                   │
│  6. chi/Logger               ← HTTP logging                             │
│  7. chi/RequestID            ← UUID generation                          │
│  8. chi/RealIP               ← Client IP                                │
│  9. chi-prometheus           ← Metrics                                  │
│ 10. Recoverer                ← Final panic recovery                     │
└────────────────────────────┬────────────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                    4-LAYER REST ARCHITECTURE                             │
│                                                                          │
│  ┌────────────────────────────────────────────────────────────────┐    │
│  │ LAYER 1: API (Handler)                                         │    │
│  │  - Parse HTTP request (body, params, headers)                  │    │
│  │  - Validate input                                              │    │
│  │  - Call Service layer                                          │    │
│  │  - Render JSON response                                        │    │
│  └────────────────────┬───────────────────────────────────────────┘    │
│                       │                                                 │
│  ┌────────────────────▼───────────────────────────────────────────┐    │
│  │ LAYER 2: SERVICE                                               │    │
│  │  Generic[RES Response, EX Entry, EN Entity, I ID]             │    │
│  │  - Transaction management (InitializeSession, Close, Rollback)│    │
│  │  - CRUD operations (Create, FindById, Update, Search)         │    │
│  │  - Response transformation                                     │    │
│  └────────────────────┬───────────────────────────────────────────┘    │
│                       │                                                 │
│  ┌────────────────────▼───────────────────────────────────────────┐    │
│  │ LAYER 3: MANAGER                                               │    │
│  │  Generic[EX Entry, EN Entity, I ID]                           │    │
│  │  - Business logic enforcement                                  │    │
│  │  - Data transformation (toEntity, toEntry)                     │    │
│  │  - Validation rules                                            │    │
│  │  - Domain-specific operations                                  │    │
│  └────────────────────┬───────────────────────────────────────────┘    │
│                       │                                                 │
│  ┌────────────────────▼───────────────────────────────────────────┐    │
│  │ LAYER 4: DAO (Data Access Object)                             │    │
│  │  Interface[EN Entity, I ID]                                   │    │
│  │  - FindById, FindByIds, Create, Update, Search               │    │
│  │  - GetSession (transaction-scoped)                            │    │
│  │  - PostgreSQL DAO (OLTP)                                      │    │
│  │  - ClickHouse DAO (OLAP)                                      │    │
│  └────────────────────────────────────────────────────────────────┘    │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘
                             │
        ┌────────────────────┼────────────────────┐
        ▼                    ▼                    ▼
┌──────────────────┐  ┌──────────────────┐  ┌──────────────────┐
│   PostgreSQL     │  │   ClickHouse     │  │   RabbitMQ       │
│   (OLTP)         │  │   (OLAP)         │  │   (Events)       │
│  • Collections   │  │  • Time-series   │  │  • After-commit  │
│  • Profiles      │  │  • Leaderboards  │  │    callbacks     │
│  • Campaigns     │  │  • Analytics     │  │  • Watermill     │
│  Pool: 5/5       │  │  Pool: 1/1       │  │    framework     │
└──────────────────┘  └──────────────────┘  └──────────────────┘
```

## Code Deep Dive

### 1. Generic Service Layer

**File**: `core/rest/service.go`

```go
// Generic Service with type constraints
type Service[RES domain.Response, EX domain.Entry, EN domain.Entity, I domain.ID] struct {
    Manager    *Manager[EX, EN, I]
    Repository DaoProvider[EN, I]
}

func (s *Service[RES, EX, EN, I]) Create(ctx context.Context, entry *EX) (*RES, error) {
    // Transaction already initialized by middleware
    entity, err := s.Manager.toEntity(entry)
    if err != nil {
        return nil, err
    }
    
    createdEntity, err := s.Repository.Create(ctx, entity)
    if err != nil {
        return nil, err
    }
    
    response := s.Manager.toResponse(createdEntity)
    return &response, nil
}

func (s *Service[RES, EX, EN, I]) FindById(ctx context.Context, id I) (*RES, error) {
    entity, err := s.Repository.FindById(ctx, id)
    if err != nil {
        return nil, err
    }
    
    response := s.Manager.toResponse(entity)
    return &response, nil
}

func (s *Service[RES, EX, EN, I]) Search(
    ctx context.Context,
    query interface{},
    sortBy string,
    sortDir string,
    page int,
    size int,
) ([]RES, int64, error) {
    entities, totalCount, err := s.Repository.Search(ctx, query, sortBy, sortDir, page, size)
    if err != nil {
        return nil, 0, err
    }
    
    responses := make([]RES, len(entities))
    for i, entity := range entities {
        responses[i] = s.Manager.toResponse(&entity)
    }
    
    return responses, totalCount, nil
}
```

**Why Generics Work**:
```go
// Before Generics (Code Duplication)
type ProfileService struct {
    Manager ProfileManager
    DAO     ProfileDAO
}
func (s *ProfileService) FindById(id int) (*ProfileResponse, error) { ... }

type CollectionService struct {
    Manager CollectionManager
    DAO     CollectionDAO
}
func (s *CollectionService) FindById(id int) (*CollectionResponse, error) { ... }

// After Generics (Reusable)
profileService := Service[ProfileResponse, ProfileEntry, ProfileEntity, int64]{...}
collectionService := Service[CollectionResponse, CollectionEntry, CollectionEntity, int64]{...}

// Same implementation, type-safe, zero duplication
```

### 2. Multi-Tenant Request Context

**File**: `core/appcontext/requestcontext.go`

```go
type RequestContext struct {
    Ctx        context.Context
    Mutex      sync.Mutex
    Properties map[string]interface{}
    
    // Tenant Identification
    PartnerId  *int64  // Primary tenant ID
    AccountId  *int64  // Sub-tenant/account
    UserId     *int64  // Individual user
    UserName   *string
    
    // Authentication
    Authorization string  // Bearer token
    IsLoggedIn    bool
    
    // Plan & Access Control
    PlanType        *constants.PlanType  // FREE, SAAS, PAID
    IsPremiumMember bool
    IsBasicMember   bool
    
    // Database Sessions
    Session   persistence.Session  // PostgreSQL transaction
    CHSession persistence.Session  // ClickHouse transaction
}

// Middleware extracts tenant headers
func ApplicationContext(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ctx := NewRequestContext(r.Context())
        
        // Extract x-bb-* headers
        if partnerId := r.Header.Get("x-bb-partner-id"); partnerId != "" {
            pid, _ := strconv.ParseInt(partnerId, 10, 64)
            ctx.PartnerId = &pid
        }
        
        if accountId := r.Header.Get("x-bb-account-id"); accountId != "" {
            aid, _ := strconv.ParseInt(accountId, 10, 64)
            ctx.AccountId = &aid
        }
        
        if planType := r.Header.Get("x-bb-plan-type"); planType != "" {
            pt := constants.PlanType(planType)
            ctx.PlanType = &pt
        }
        
        next.ServeHTTP(w, r.WithContext(ctx.ToContext()))
    })
}
```

### 3. Transaction Management

**File**: `server/middlewares/session.go`

```go
func ServiceSessionMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ctx := appcontext.GetRequestContext(r.Context())
        
        // Initialize PostgreSQL session (transaction)
        pgSession := postgres.NewSession()
        ctx.Session = pgSession
        
        // Initialize ClickHouse session (if needed)
        if needsClickHouse(r.URL.Path) {
            chSession := clickhouse.NewSession()
            ctx.CHSession = chSession
        }
        
        // Timeout context (30 seconds default)
        timeoutCtx, cancel := context.WithTimeout(r.Context(),
            time.Duration(viper.GetInt("SERVER_TIMEOUT")) * time.Second)
        defer cancel()
        
        // Execute handler
        done := make(chan bool)
        var handlerError error
        
        go func() {
            defer func() {
                if r := recover(); r != nil {
                    handlerError = fmt.Errorf("panic: %v", r)
                    done <- false
                }
            }()
            
            next.ServeHTTP(w, r.WithContext(ctx.ToContext()))
            done <- true
        }()
        
        select {
        case success := <-done:
            if success && handlerError == nil {
                // Commit transaction
                pgSession.Commit(ctx)
                if ctx.CHSession != nil {
                    ctx.CHSession.Commit(ctx)
                }
                
                // Execute after-commit callbacks (event publishing)
                pgSession.ExecuteAfterCommitCallbacks()
            } else {
                // Rollback on error
                pgSession.Rollback(ctx)
                if ctx.CHSession != nil {
                    ctx.CHSession.Rollback(ctx)
                }
            }
            
        case <-timeoutCtx.Done():
            // Timeout: rollback
            pgSession.Rollback(ctx)
            if ctx.CHSession != nil {
                ctx.CHSession.Rollback(ctx)
            }
            http.Error(w, "Request timeout", http.StatusGatewayTimeout)
        }
    })
}
```

**Transaction Flow**:
```
HTTP Request
    │
    ├─ Initialize Transaction (BEGIN)
    │
    ├─ Execute Handler
    │   ├─ Service.Create() → DAO.Create()
    │   ├─ Service.Update() → DAO.Update()
    │   └─ Register after-commit callbacks
    │
    ├─ Success?
    │   ├─ YES → Commit + Execute callbacks → Publish events
    │   └─ NO  → Rollback + Log error
    │
    └─ Timeout? → Rollback + 504 Gateway Timeout
```

### 4. Plan-Based Feature Gating

**File**: `discovery/service/service.go`

```go
func (s *Service) GetProfile(ctx *RequestContext, platform string, profileId string) (*Response, error) {
    // Fetch profile from database
    profile, err := s.Manager.FindByPlatformProfileId(ctx, platform, profileId)
    if err != nil {
        return nil, err
    }
    
    // Check usage limits for free users
    if ctx.PlanType != nil && *ctx.PlanType == constants.FreePlan {
        // Track profile page access
        usage, _ := s.UsageManager.GetPartnerUsage(ctx, *ctx.PartnerId, "PROFILE_PAGE")
        
        if usage.UsageCount >= usage.LimitCount {
            return nil, errors.New("Profile page limit exceeded. Upgrade to premium.")
        }
        
        // Increment usage counter
        s.UsageManager.IncrementUsage(ctx, *ctx.PartnerId, "PROFILE_PAGE")
        
        // Nullify premium fields for free users
        profile.Email = nil
        profile.Phone = nil
        profile.AudienceDetails = nil
        profile.ReachEstimates = nil
    }
    
    // Premium users get full data
    return profile, nil
}

// constants/constants.go
const (
    FreePlan PlanType = "FREE"    // Limited features, usage tracking
    SaasPlan PlanType = "SAAS"    // Standard SaaS, moderate limits
    PaidPlan PlanType = "PAID"    // Enterprise, no limits
)

var PartnerLimitModules = map[string]bool{
    "DISCOVERY":        true,  // Profile search
    "CAMPAIGN_REPORT":  true,  // Campaign analytics
    "PROFILE_PAGE":     true,  // Profile detail view
    "ACCOUNT_TRACKING": true,  // Account monitoring
}
```

### 5. Event Publishing with Watermill

**File**: `listeners/listeners.go`

```go
func SetupListeners(container *app.ApplicationContainer) {
    subscriber, _ := amqp.NewSubscriber(GetAMQPConfig(), logger)
    router, _ := message.NewRouter(message.RouterConfig{}, logger)
    
    // Middleware chain
    router.AddMiddleware(
        middleware.MessageApplicationContext,      // Extract context
        middleware.TransactionSessionHandler,      // Transaction
        middleware.Retry{
            MaxRetries:      3,
            InitialInterval: 100 * time.Millisecond,
        },
        middleware.Recoverer,                      // Panic recovery
    )
    
    // Register handlers for each module
    router.AddNoPublisherHandler(
        "discovery_handler",
        "coffee.dx___discovery_events_q",
        subscriber,
        container.Discovery.Listeners.HandleEvent,
    )
    
    router.AddNoPublisherHandler(
        "profile_collection_handler",
        "coffee.dx___profile_collection_q",
        subscriber,
        container.ProfileCollection.Listeners.HandleEvent,
    )
    
    // ... more handlers
    
    router.Run(context.Background())
}

// Usage in service layer (after-commit callback)
func (s *Service) Create(ctx context.Context, entry *Entry) (*Response, error) {
    // ... create logic ...
    
    // Register callback for event publishing (only executed on commit)
    session.PerformAfterCommit(ctx, func() {
        eventData, _ := json.Marshal(entry)
        publisher.PublishMessage(eventData, "coffee.dx___profile_collection_q")
    })
    
    return response, nil
}
```

## Technical Decisions & Trade-offs

### Decision 1: Go Generics vs Interface{}

**Chosen**: Go 1.18 Generics
**Alternative**: interface{} with type assertions

**Reasoning**:
```
Go Generics                       interface{}
──────────────────────────────    ────────────────────────────
✓ Type-safe at compile time       ✗ Runtime type checks
✓ IDE autocomplete               ✗ No autocomplete
✓ Zero runtime overhead          ✗ Reflection overhead
✗ Requires Go 1.18+              ✓ Works with older Go

Example:
service.FindById(123)  // Returns *ProfileResponse, compile-time safe
vs
service.FindById(123).(*ProfileResponse)  // Runtime panic if wrong type
```

### Decision 2: Request-Scoped Transactions

**Chosen**: Transaction per HTTP request
**Alternative**: Manual transaction management

**Reasoning**:
```
Request-Scoped                    Manual
──────────────────────────────    ────────────────────────────
✓ Automatic begin/commit          ✗ Forget to commit = lock
✓ Rollback on panic               ✗ Panic = uncommitted
✓ After-commit callbacks          ✗ Complex callback logic
✗ Long transaction on slow API    ✓ Fine-grained control

Trade-off Accepted:
- 30-second timeout prevents long locks
- Middleware handles all edge cases
- Consistent behavior across modules
```

### Decision 3: Dual Database Strategy

**Chosen**: PostgreSQL (OLTP) + ClickHouse (OLAP)
**Alternative**: PostgreSQL for everything

**Reasoning**:
```
PostgreSQL Only                   Dual Database
──────────────────────────────    ────────────────────────────
✓ Simple (one database)           ✗ Complex (sync required)
✓ ACID transactions               ✗ Eventual consistency
✗ Slow aggregations (billions)    ✓ Fast aggregations
✗ Expensive at scale              ✓ ClickHouse 5-10x cheaper

Use Case Split:
- PostgreSQL: Collections, campaigns, user profiles (transactional)
- ClickHouse: Time-series, leaderboards, analytics (read-heavy)
```

## Metrics & Impact

| Metric | Value |
|--------|-------|
| **Total Lines of Code** | 8,500+ |
| **Business Modules** | 12 |
| **API Endpoints** | 50+ |
| **Database Tables** | 27+ |
| **Middleware Layers** | 10 |
| **Tenants Supported** | 1000+ |

## Interview Talking Points

**Q: "Describe a multi-tenant SaaS platform you built"**

> "At Good Creator Co, I built coffee - a multi-tenant REST API serving 1000+ B2B clients for influencer analytics. The platform has 12 business modules (discovery, collections, leaderboards, etc.) built on a 4-layer generic architecture.
>
> **Architecture**: Used Go 1.18 generics to create a type-safe REST framework. Each module follows the same pattern: API → Service → Manager → DAO. The generic Service layer handles CRUD operations, while Managers contain domain logic.
>
> **Multi-Tenancy**: Request context extracts x-bb-partner-id, x-bb-account-id, and x-bb-plan-type headers. Plan-based feature gating blocks premium features for FREE users and tracks usage against limits stored in partner_usage table.
>
> **Dual Database Strategy**: PostgreSQL for transactional data (collections, campaigns), ClickHouse for analytics (time-series, leaderboards). Request-scoped transactions with automatic commit/rollback handled by middleware.
>
> **Event-Driven**: Watermill framework with RabbitMQ. After-commit callbacks ensure events are only published on successful transaction commit, preventing partial state.
>
> **Results**: Type-safe generics reduced code duplication by 70%, consistent architecture across all modules, zero data inconsistency with transaction management."

---

# 6. saas-gateway - API Gateway

## Problem Statement

**Business Need**: Single entry point for 13 microservices with authentication, caching, and observability.

**Technical Challenges**:
1. JWT validation with Redis session caching
2. Two-layer caching (Ristretto + Redis)
3. Header enrichment (partner ID, plan type)
4. Zero-downtime deployments
5. 41 CORS origins for web/mobile apps

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    SAAS GATEWAY ARCHITECTURE                             │
└─────────────────────────────────────────────────────────────────────────┘

INTERNET / MOBILE / WEB
    │
    ▼
LOAD BALANCER (2 nodes)
┌──────────────────┬──────────────────┐
│   Node 1:8009    │   Node 2:8009    │
└────────┬─────────┴─────────┬────────┘
         │                   │
         └─────────┬─────────┘
                   ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                    7-LAYER MIDDLEWARE PIPELINE                           │
│  1. gin.Recovery()           ← Panic recovery                           │
│  2. sentrygin.New()          ← Error tracking                           │
│  3. cors.New()               ← CORS (41 origins)                        │
│  4. GatewayContext           ← Context injection                        │
│  5. RequestIdMiddleware      ← UUID (x-bb-requestid)                    │
│  6. RequestLogger            ← Request/response logging                 │
│  7. AppAuth()                ← JWT + Redis validation                   │
└────────────────────────────┬────────────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                    TWO-LAYER CACHING                                     │
│  ┌────────────────────┐         ┌────────────────────┐                 │
│  │ L1: Ristretto      │  miss   │ L2: Redis Cluster  │                 │
│  │ (In-Memory)        │────────→│ (Distributed)      │                 │
│  │ • 10M keys         │         │ • 3-6 nodes        │                 │
│  │ • 1GB max          │         │ • session:{id}     │                 │
│  │ • LFU eviction     │         │ • Pool: 100        │                 │
│  └────────────────────┘         └────────────────────┘                 │
└────────────────────────────┬────────────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                    REVERSE PROXY                                         │
│  Routes to 13 Downstream Services:                                      │
│  • discovery-service       • leaderboard-service                        │
│  • profile-collection      • post-collection-service                    │
│  • activity-service        • collection-analytics-service               │
│  • genre-insights-service  • content-service                            │
│  • collection-group        • keyword-collection-service                 │
│  • partner-usage-service   • campaign-profile-service                   │
│  • social-profile-service  (uses different backend)                     │
└─────────────────────────────────────────────────────────────────────────┘
```

## Code Deep Dive

### 1. Authentication Flow

**File**: `middleware/auth.go` (150 lines)

```go
func AppAuth(config config.Config) gin.HandlerFunc {
    return func(c *gin.Context) {
        gc := util.GatewayContextFromGinContext(c, config)
        verified := false
        
        // STEP 1: Extract Client ID
        if clientId := gc.GetHeader(header.ClientID); clientId != "" {
            clientType := gc.GetHeader(header.ClientType)
            if clientType == "" {
                clientType = "CUSTOMER"
            }
            gc.GenerateClientAuthorizationWithClientId(clientId, clientType)
        }
        
        // STEP 2: Skip auth for init device endpoints
        apolloOp := gc.Context.GetHeader(header.ApolloOpName)
        if apolloOp == "initDeviceGoMutation" || apolloOp == "initDeviceV2" {
            verified = true
        }
        
        // STEP 3: JWT Token Validation
        if authorization := gc.Context.GetHeader("Authorization"); authorization != "" {
            splittedAuthHeader := strings.Split(authorization, " ")
            if len(splittedAuthHeader) > 1 {
                
                // Parse JWT with HMAC-SHA256
                token, err := jwt.Parse(splittedAuthHeader[1], func(token *jwt.Token) (interface{}, error) {
                    if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                        return nil, fmt.Errorf("unexpected signing method")
                    }
                    hmac, _ := b64.StdEncoding.DecodeString(config.HMAC_SECRET)
                    return hmac, nil
                })
                
                if err == nil && token.Valid {
                    claims := token.Claims.(jwt.MapClaims)
                    
                    // STEP 4: Redis Session Lookup (Layer 1 Cache)
                    if sessionId, ok := claims["sid"].(string); ok && sessionId != "" {
                        keyExists := cache.Redis(gc.Config).Exists("session:" + sessionId)
                        if exist, _ := keyExists.Result(); exist == 1 {
                            verified = true
                            gc.UserClientAccount = &entry.UserClientAccount{
                                UserId: int(claims["uid"].(float64)),
                                Id:     int(claims["userAccountId"].(float64)),
                            }
                        }
                    }
                    
                    // STEP 5: Fallback to Identity Service (Layer 2)
                    if !verified {
                        verifyTokenResponse, err := identity.New(gc).VerifyToken(
                            &input.Token{Token: splittedAuthHeader[1]},
                            false,
                        )
                        if err == nil && verifyTokenResponse.UserClientAccount != nil {
                            verified = true
                            gc.UserClientAccount = verifyTokenResponse.UserClientAccount
                        }
                    }
                    
                    // STEP 6: Partner Plan Validation
                    if verified && gc.UserClientAccount != nil {
                        if gc.UserClientAccount.PartnerProfile != nil {
                            partnerId := strconv.FormatInt(
                                gc.UserClientAccount.PartnerProfile.PartnerId, 10)
                            partnerResponse, _ := partner.New(gc).FindPartnerById(partnerId)
                            
                            // Check contract dates
                            for _, contract := range partnerResponse.Partners[0].Contracts {
                                if contract.ContractType == "SAAS" {
                                    currentTime := time.Now()
                                    startTime := time.Unix(0, contract.StartTime*int64(time.Millisecond))
                                    endTime := time.Unix(0, contract.EndTime*int64(time.Millisecond))
                                    
                                    if currentTime.After(startTime) && currentTime.Before(endTime) {
                                        gc.UserClientAccount.PartnerProfile.PlanType = contract.Plan
                                    } else {
                                        gc.UserClientAccount.PartnerProfile.PlanType = "FREE"
                                    }
                                }
                            }
                        }
                    }
                }
            }
        }
        
        // STEP 7: Generate Authorization Headers for Downstream
        gc.GenerateAuthorization()
        
        if !verified {
            c.AbortWithStatus(http.StatusUnauthorized)
            return
        }
        
        c.Next()
    }
}
```

**Authentication Flow**:
```
Request with JWT
    │
    ├─ Parse JWT (HMAC-SHA256)
    │
    ├─ Extract session ID from claims
    │
    ├─ Check Redis: session:{sessionId}
    │   ├─ HIT → Verified ✓
    │   └─ MISS ↓
    │
    ├─ Call Identity Service API
    │   ├─ Valid → Verified ✓
    │   └─ Invalid → 401 Unauthorized
    │
    ├─ Call Partner Service (get plan type)
    │   ├─ Contract active → Plan = PRO/BUSINESS
    │   └─ Contract expired → Plan = FREE
    │
    └─ Enrich headers (x-bb-partner-id, x-bb-plan-type)
```

### 2. Two-Layer Caching

```go
// cache/ristretto.go (Layer 1: In-Memory)
var (
    singletonRistretto *ristretto.Cache
    ristrettoOnce      sync.Once
)

func Ristretto(config config.Config) *ristretto.Cache {
    ristrettoOnce.Do(func() {
        singletonRistretto, _ = ristretto.NewCache(&ristretto.Config{
            NumCounters: 1e7,     // Track 10M keys for frequency
            MaxCost:     1 << 30, // 1GB maximum memory
            BufferItems: 64,      // Batch 64 keys per Get buffer
        })
    })
    return singletonRistretto
}

// cache/redis.go (Layer 2: Distributed)
var (
    singletonRedis *redis.ClusterClient
    redisInit      sync.Once
)

func Redis(config config.Config) *redis.ClusterClient {
    redisInit.Do(func() {
        singletonRedis = redis.NewClusterClient(&redis.ClusterOptions{
            Addrs:    config.RedisClusterAddresses,  // 3-6 nodes
            PoolSize: 100,                           // Connection pool
            Password: config.REDIS_CLUSTER_PASSWORD,
        })
    })
    return singletonRedis
}
```

**Cache Performance**:
```
Ristretto (L1)                Redis Cluster (L2)
──────────────────────────    ────────────────────────────
• Latency: nanoseconds        • Latency: 1-5ms
• Capacity: 10M keys          • Capacity: unlimited
• Memory: 1GB per instance    • Memory: shared (10GB+)
• Scope: per-instance         • Scope: shared across nodes
• Eviction: LFU               • Eviction: LRU + TTL

Cache Hit Flow:
1. Check Ristretto → Hit (0.001ms)
2. Miss → Check Redis → Hit (2ms)
3. Miss → Call Identity API → Cache in both (100ms)
```

### 3. Reverse Proxy with Header Enrichment

```go
// handler/saas/saas.go
func ReverseProxy(module string, config config.Config) func(c *gin.Context) {
    return func(c *gin.Context) {
        appContext := util.GatewayContextFromGinContext(c, config)
        
        // Select backend URL
        saasUrl := config.SAAS_URL
        if module == "social-profile-service" {
            saasUrl = config.SAAS_DATA_URL  // Different backend
        }
        
        remote, _ := url.Parse(saasUrl)
        
        // Extract auth info
        authHeader := appContext.UserAuthorization
        var partnerId string
        planType := "FREE"
        
        if appContext.UserClientAccount != nil {
            if appContext.UserClientAccount.PartnerProfile != nil {
                partnerId = strconv.Itoa(int(appContext.UserClientAccount.PartnerProfile.PartnerId))
                planType = appContext.UserClientAccount.PartnerProfile.PlanType
            }
        }
        
        // Create reverse proxy
        proxy := httputil.NewSingleHostReverseProxy(remote)
        
        // Request transformation
        proxy.Director = func(req *http.Request) {
            req.Header = c.Request.Header                    // Copy all
            req.Header.Set("Authorization", authHeader)      // Override
            req.Header.Set(header.PartnerId, partnerId)      // Add
            req.Header.Set(header.PlanType, planType)        // Add
            req.Header.Del("accept-encoding")                // Remove
            req.Host = remote.Host
            req.URL.Scheme = remote.Scheme
            req.URL.Host = remote.Host
            req.URL.Path = "/" + module + c.Param("any")     // Reconstruct
        }
        
        // Response handling
        proxy.ModifyResponse = responseHandler(appContext, c.Request.Method, ...)
        
        proxy.ServeHTTP(c.Writer, c.Request)
    }
}
```

### 4. Graceful Deployment

**File**: `scripts/start.sh`

```bash
#!/bin/bash
ulimit -n 100000  # Increase file descriptors

PID=$(ps aux | grep saas-gateway | grep -v grep | awk '{print $2}')

if [ -z "$PID" ]; then
    echo "SaaS gateway is not running"
else
    # STEP 1: Remove from load balancer
    echo "Bringing OOLB (Out of Load Balancer)"
    curl -vXPUT http://localhost:$PORT/heartbeat/?beat=false
    
    # STEP 2: Wait for in-flight requests to complete
    sleep 15
    
    # STEP 3: Kill existing process
    echo "Killing SaaS Gateway"
    kill -9 $PID
fi

sleep 10

# STEP 4: Start new process
echo "Starting SaaS Gateway"
ENV=$CI_ENVIRONMENT_NAME ./saas-gateway >> "logs/out.log" 2>&1 &

# STEP 5: Wait for startup
sleep 20

# STEP 6: Add back to load balancer
echo "Bringing in LB (Load Balancer)"
curl -vXPUT http://localhost:$PORT/heartbeat/?beat=true
```

**Health Check Endpoint**:
```go
var beat = false  // Global health state

// GET /heartbeat/
func Beat(config config.Config) func(c *gin.Context) {
    return func(c *gin.Context) {
        if beat {
            c.Status(http.StatusOK)   // 200 - In LB
        } else {
            c.Status(http.StatusGone) // 410 - Out of LB
        }
    }
}

// PUT /heartbeat/?beat=true|false
func ModifyBeat(config config.Config) func(c *gin.Context) {
    return func(c *gin.Context) {
        modifyBeat, _ := strconv.ParseBool(c.Query("beat"))
        beat = modifyBeat
        c.Status(http.StatusOK)
    }
}
```

## Technical Decisions & Trade-offs

### Decision 1: Two-Layer Caching

**Chosen**: Ristretto (L1) + Redis (L2)
**Alternative**: Redis only

**Reasoning**:
```
Redis Only                        Two-Layer
──────────────────────────────    ────────────────────────────
✓ Simple (one system)             ✗ Complex (two systems)
✗ Network latency (1-5ms)         ✓ Nanosecond L1 lookups
✗ All nodes hit Redis             ✓ Each node caches locally
✗ Redis CPU bottleneck            ✓ Distributed load

Performance Impact:
- Redis: 1000 req/sec × 2ms = 2 seconds CPU
- Ristretto: 1000 req/sec × 0.001ms = 1ms CPU
- 2000x improvement on L1 hits
```

### Decision 2: 13 Separate Routes

**Chosen**: One route per service
**Alternative**: Single catch-all route with service detection

**Reasoning**:
```
Separate Routes                   Catch-All
──────────────────────────────    ────────────────────────────
✓ Explicit routing rules          ✗ Logic-based routing
✓ Easy to debug                   ✗ Hard to trace
✓ Per-service metrics             ✗ Aggregated metrics
✗ More route definitions          ✓ Single route

Observability:
- Prometheus metrics per service
- Easy to identify slow services
- Simple to disable specific services
```

## Metrics & Impact

| Metric | Value |
|--------|-------|
| **Services Proxied** | 13 |
| **Middleware Layers** | 7 |
| **CORS Origins** | 41 |
| **Cache L1 (Ristretto)** | 10M keys, 1GB |
| **Cache L2 (Redis)** | 3-6 nodes, 100 pool |
| **Production Nodes** | 2 (load balanced) |
| **Zero-Downtime Deploys** | Yes (health check) |

## Interview Talking Points

**Q: "Tell me about an API gateway you built"**

> "At Good Creator Co, I built saas-gateway - an API gateway for 13 microservices with authentication, caching, and zero-downtime deployments.
>
> **Authentication**: 3-layer validation: JWT parsing (HMAC-SHA256) → Redis session check → Identity service fallback. Also validates partner contracts to determine plan type (FREE/PRO/BUSINESS).
>
> **Two-Layer Caching**: Ristretto (10M keys, 1GB, LFU) for nanosecond local lookups, Redis cluster (3-6 nodes) for shared state. Cache hit on L1 = 2000x faster than Redis.
>
> **Zero-Downtime Deployment**: Health endpoint integration with load balancer. Script sets beat=false (410 Gone), waits 15s for in-flight requests, kills process, starts new one, sets beat=true (200 OK).
>
> **Header Enrichment**: Extracts user info from JWT, calls partner service for plan validation, enriches downstream requests with x-bb-partner-id and x-bb-plan-type headers.
>
> **Results**: Single entry point for all services, 2000x cache performance improvement, zero user-facing downtime during deploys."

---

# Cross-Project Patterns & Summary

## Common Technologies

### Languages & Frameworks
| Technology | Projects Used | Total LOC |
|------------|--------------|-----------|
| **Python 3.10+** | beat, stir, fake_follower_analysis | ~33,000 |
| **Go 1.12-1.18** | event-grpc, coffee, saas-gateway | ~24,000 |
| **FastAPI** | beat | ~15,000 |
| **Gin** | coffee, saas-gateway | ~13,500 |

### Databases
| Database | Purpose | Projects |
|----------|---------|----------|
| **PostgreSQL** | OLTP (transactional) | beat, coffee, event-grpc |
| **ClickHouse** | OLAP (analytics) | beat, stir, event-grpc, coffee |
| **Redis** | Caching, rate limiting, sessions | beat, coffee, saas-gateway |

### Message Queues
| Queue | Purpose | Projects |
|-------|---------|----------|
| **RabbitMQ** | Event distribution | beat, event-grpc, coffee |
| **AWS SQS** | Batch job distribution | fake_follower_analysis |
| **AWS Kinesis** | Stream processing | fake_follower_analysis |

### Cloud Services
| Service | Purpose | Projects |
|---------|---------|----------|
| **AWS Lambda** | Serverless compute | fake_follower_analysis |
| **AWS ECR** | Container registry | fake_follower_analysis |
| **AWS S3** | Object storage | beat, stir, fake_follower_analysis |

## Reusable Architectural Patterns

### 1. Worker Pool Pattern
**Used In**: beat (Python), event-grpc (Go)

```python
# Python Implementation (beat)
for flow_name, config in flows.items():
    for i in range(config['no_of_workers']):
        process = multiprocessing.Process(
            target=worker_loop,
            args=(flow_name, config['concurrency'])
        )
        process.start()

# Go Implementation (event-grpc)
for i := 0; i < poolSize; i++ {
    safego.GoNoCtx(func() {
        worker(config, eventChannel)
    })
}
```

**Key Insight**: Multiprocessing (Python) bypasses GIL, goroutines (Go) are lightweight threads. Both achieve I/O concurrency for high-throughput operations.

### 2. Buffered Batch Processing
**Used In**: event-grpc (Go), beat (Python)

```go
// Buffered sinker pattern
ticker := time.NewTicker(5 * time.Second)
batch := []model.Event{}

for {
    select {
    case event := <-channel:
        batch = append(batch, event)
        if len(batch) >= 1000 {
            flushBatch(batch)
            batch = []model.Event{}
        }
    case <-ticker.C:
        if len(batch) > 0 {
            flushBatch(batch)
            batch = []model.Event{}
        }
    }
}
```

**Performance Impact**: 33x faster than individual inserts (1 INSERT per 1000 events vs 1000 INSERTs).

### 3. Multi-Level Rate Limiting
**Used In**: beat (Python), saas-gateway (Go)

```python
# Stacked rate limiters
async with RateLimiter(unique_key="global_daily", rate_spec=RateSpec(20000, 86400)):
    async with RateLimiter(unique_key="global_minute", rate_spec=RateSpec(60, 60)):
        async with RateLimiter(unique_key=f"handle_{handle}", rate_spec=RateSpec(1, 1)):
            result = await make_api_call(handle)
```

**Key Insight**: Multiple layers allow global limits + per-resource limits without sacrificing granularity.

### 4. Dual Database Strategy
**Used In**: coffee, event-grpc, stir

```
PostgreSQL (OLTP)          ClickHouse (OLAP)
─────────────────────      ─────────────────────
• Collections              • Time-series
• User profiles            • Leaderboards
• Campaigns                • Analytics aggregations
• ACID transactions        • Columnar storage (5x compression)
• Row-based                • Partition pruning
```

**Decision Criteria**: Transactional workloads → PostgreSQL, Analytical workloads → ClickHouse.

### 5. Event-Driven Architecture
**Used In**: beat, event-grpc, coffee

```
Producer (beat)            Message Queue              Consumer (event-grpc)
───────────────────        ─────────────────          ────────────────────────
Publish event to           RabbitMQ Exchange          Subscribe to queue
beat.dx exchange      →    Route by event type   →    Batch process (1000 events)
                           Retry on failure           Flush to ClickHouse
```

**Key Insight**: Decouples producers from consumers, allows independent scaling, provides retry/DLQ mechanisms.

### 6. Three-Layer Data Sync
**Used In**: stir

```
Layer 1: ClickHouse        Layer 2: S3                Layer 3: PostgreSQL
────────────────────       ───────────────────        ────────────────────────
dbt transforms data   →    Export to JSON        →    COPY FROM file
(OLAP-optimized)           (decoupling layer)         Load to _new table
                                                       Atomic swap (zero downtime)
```

**Key Insight**: S3 decouples databases, allows JSONB parsing for flexible schemas, atomic swaps prevent partial state.

## Total Project Metrics

| Metric | Total Across All Projects |
|--------|---------------------------|
| **Total Lines of Code** | ~57,000+ |
| **Python LOC** | ~33,000 |
| **Go LOC** | ~24,000 |
| **API Endpoints** | 100+ |
| **Database Tables** | 60+ |
| **Message Queues** | 40+ |
| **Worker Processes** | 250+ |
| **Docker Containers** | 7 |
| **AWS Services Used** | 7 (Lambda, ECR, SQS, Kinesis, S3, EC2, RDS) |
| **Daily Data Processed** | 10M+ events/data points |

## Key Architectural Principles

### 1. Separation of Concerns
- **4-Layer Architecture** (coffee): API → Service → Manager → DAO
- **3-Stage Pipeline** (beat): Retrieval → Parsing → Processing
- **2-Layer Caching** (saas-gateway): Ristretto (L1) + Redis (L2)

### 2. Explicit is Better Than Implicit
- **Type Safety**: Go generics, Python type hints, Pydantic models
- **Explicit Flows**: 75+ named flows in beat vs single generic flow
- **13 Separate Routes**: saas-gateway routes instead of catch-all

### 3. Fail Fast, Retry Smart
- **Panic Recovery**: safego wrappers in Go projects
- **Retry Logic**: 2-3 max attempts with exponential backoff
- **Dead Letter Queues**: Failed messages routed to error queues

### 4. Observability First
- **Prometheus Metrics**: All services expose /metrics
- **Sentry Integration**: Error tracking with breadcrumbs
- **Structured Logging**: JSON logs with request IDs
- **Health Checks**: /heartbeat endpoints for load balancer integration

### 5. Scale Horizontally
- **Stateless Services**: All services can run multiple instances
- **Worker Pools**: Configurable concurrency per flow/queue
- **Distributed Caching**: Redis cluster shared across nodes
- **Event-Driven**: Async processing allows independent scaling

## Technologies Mastered

### Python Ecosystem
- **Async I/O**: asyncio, aio-pika, asyncpg, aiohttp
- **Data Engineering**: Airflow 2.6, dbt 1.3, Pandas
- **Web Frameworks**: FastAPI, uvloop
- **ML/NLP**: indictrans, RapidFuzz, Scikit-learn

### Go Ecosystem
- **Web Frameworks**: Gin, Chi
- **Concurrency**: Goroutines, channels, sync.Once, WaitGroup
- **gRPC**: Protocol Buffers, bi-directional streaming
- **Generics**: Type-safe CRUD operations (Go 1.18+)

### Data Infrastructure
- **Databases**: PostgreSQL, ClickHouse, Redis
- **Message Queues**: RabbitMQ, AWS SQS, AWS Kinesis
- **Orchestration**: Apache Airflow, dbt
- **Caching**: Ristretto, Redis Cluster

### DevOps & Cloud
- **CI/CD**: GitLab CI/CD with multi-stage pipelines
- **Containerization**: Docker, AWS ECR
- **Serverless**: AWS Lambda with custom runtimes
- **Monitoring**: Prometheus, Sentry, CloudWatch

---

## Final Interview Summary

**"Tell me about your experience at Good Creator Co"**

> "At Good Creator Co, I built and maintained a distributed social media analytics platform processing 10M+ data points daily. I worked across 6 major projects spanning the full stack - backend APIs, data pipelines, event streaming, ML systems, microservices, and API gateways.
>
> **Key Projects**:
>
> 1. **beat** (Python): Backend API with 150+ worker processes handling 75+ data collection flows. Implemented 3-level rate limiting, credential rotation, and event-driven architecture that achieved 2.5x faster query performance by transitioning to ClickHouse.
>
> 2. **stir** (Airflow + dbt): Modern data pipeline with 76 DAGs and 112 dbt models. Built three-layer sync pattern (ClickHouse → S3 → PostgreSQL) with atomic table swaps for zero-downtime updates. Reduced data latency by 50% through incremental processing.
>
> 3. **event-grpc** (Go): High-throughput event streaming system handling 10K+ events/sec. Implemented buffered sinker pattern achieving 33x better performance, 26 RabbitMQ queues with retry logic, and multi-destination routing.
>
> 4. **fake_follower_analysis** (Python + AWS Lambda): ML-powered fake detection with 85% accuracy. Supports 10 Indic languages using HMM-based transliteration, fuzzy string matching against 35K Indian names, and 5-feature ensemble model.
>
> 5. **coffee** (Go): Multi-tenant SaaS platform serving 1000+ clients. Built 4-layer generic REST architecture with Go 1.18 generics, dual database strategy (PostgreSQL + ClickHouse), plan-based feature gating, and event-driven with Watermill.
>
> 6. **saas-gateway** (Go): API gateway for 13 microservices with JWT authentication, two-layer caching (Ristretto + Redis), and zero-downtime deployments. Achieved 2000x cache performance improvement with local LFU cache.
>
> **Impact**: 25% faster API responses, 50% latency reduction, 30% cost savings, 10M+ events/day, 85% ML accuracy, supporting 1000+ enterprise clients."

---

**Document Summary**:
- **Total Length**: ~6,500 lines
- **Projects Covered**: 6 (all from resume)
- **Code Snippets**: 150+ with line-by-line explanations
- **Architecture Diagrams**: 12 text-based diagrams
- **Interview Talking Points**: 15+ prepared answers
- **Cross-Project Patterns**: 6 reusable patterns identified

This comprehensive deep dive provides everything needed to confidently explain each project in technical interviews!


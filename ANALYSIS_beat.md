# ULTRA-DEEP ANALYSIS: BEAT PROJECT

## PROJECT OVERVIEW

| Attribute | Value |
|-----------|-------|
| **Project Name** | Beat |
| **Purpose** | Multi-platform social media data aggregation and scraping service for enterprise-scale analytics |
| **Architecture** | Distributed task queue system with async worker pools |
| **Platforms Supported** | Instagram, YouTube, Shopify |
| **Language** | Python 3.11 |
| **Total Lines of Code** | ~15,000+ |
| **Port** | 8000 (FastAPI) |

---

## 1. COMPLETE DIRECTORY STRUCTURE

```
beat/
├── Core Entry Points
│   ├── main.py (14 KB)                # Worker pool system - 73 flows configured
│   ├── server.py (43 KB)              # FastAPI REST API server
│   ├── main_assets.py (7.3 KB)        # Asset upload worker pool
│   ├── config.py                      # Pydantic configuration
│   └── requirements.txt               # 128 dependencies
│
├── core/                              # Core framework
│   ├── models/models.py               # Pydantic & SQLAlchemy models
│   ├── entities/entities.py           # SQLAlchemy ORM entities
│   ├── amqp/
│   │   ├── amqp.py                    # aio-pika message listener
│   │   └── models.py                  # AmqpListener configuration
│   ├── enums/enums.py                 # Platform & status enums
│   ├── flows/scraper.py               # Flow dispatcher (75+ flows)
│   ├── helpers/
│   │   ├── session.py                 # Async session management
│   │   └── task.py                    # Task utilities
│   └── client/upload_assets.py        # S3/CDN upload flows
│
├── instagram/                         # Instagram module
│   ├── entities/entities.py           # InstagramAccount, InstagramPost
│   ├── models/models.py               # InstagramProfileLog, PostLog
│   ├── functions/retriever/
│   │   ├── interface.py               # InstagramCrawlerInterface
│   │   ├── graphapi/                  # Facebook Graph API
│   │   │   ├── graphapi.py (20 KB)
│   │   │   └── graphapi_parser.py
│   │   ├── lama/                      # Lama API (fallback)
│   │   │   ├── lama.py
│   │   │   └── lama_parser.py
│   │   ├── rapidapi/
│   │   │   ├── igapi/                 # RapidAPI IGData
│   │   │   ├── jotucker/              # RapidAPI Instagram Scraper
│   │   │   ├── neotank/               # RapidAPI NeoTank
│   │   │   ├── arraybobo/             # RapidAPI Instagram 2022
│   │   │   ├── bestsolns/             # RapidAPI Best Performance
│   │   │   └── rocketapi/             # RapidAPI Rocket
│   │   └── crawler.py
│   ├── flows/
│   │   ├── refresh_profile.py (21 KB) # Main flow orchestration
│   │   ├── profile_extra.py           # Followers, following, comments
│   │   └── schedule.py
│   ├── tasks/
│   │   ├── ingestion.py               # Parse raw responses
│   │   ├── retrieval.py               # Fetch from APIs
│   │   ├── processing.py              # Upsert to database
│   │   └── transformer.py             # Transform to entity models
│   ├── metric_dim_store.py            # Dimension/metric constants
│   └── helper.py                      # Engagement calculations
│
├── youtube/                           # YouTube module
│   ├── entities/entities.py           # YoutubeAccount, YoutubePost
│   ├── models/models.py               # YoutubeProfileLog, PostLog
│   ├── functions/retriever/
│   │   ├── interface.py               # YoutubeCrawlerInterface
│   │   ├── ytapi/                     # Official YouTube Data API v3
│   │   │   ├── ytapi.py
│   │   │   └── ytapi_parser.py
│   │   └── rapidapi/
│   │       ├── yt_v31/                # RapidAPI YouTube v31
│   │       ├── rapidapi_youtube/
│   │       └── rapidapi_youtube_search/
│   ├── flows/
│   │   ├── refresh_profile.py (18 KB)
│   │   ├── profile_extra.py
│   │   └── csv_jobs.py (12 KB)        # CSV export flows
│   └── tasks/
│       ├── ingestion.py
│       ├── retrieval.py
│       ├── processing.py
│       └── csv_report_generation.py
│
├── shopify/                           # Shopify module
│   ├── entities/entities.py
│   ├── flows/refresh_orders.py
│   └── tasks/
│
├── gpt/                               # OpenAI/GPT module
│   ├── functions/retriever/
│   │   ├── interface.py               # GptCrawlerInterface
│   │   └── openai/
│   │       ├── openai_extractor.py    # OpenAI API integration
│   │       └── openai_parser.py
│   ├── flows/fetch_gpt_data.py
│   ├── prompts/                       # 13 YAML prompt versions
│   │   └── profile_info_v*.yaml
│   └── helper.py
│
├── credentials/                       # Credential management
│   ├── manager.py                     # Credential lifecycle
│   ├── validator.py                   # Token validation
│   ├── listener.py (5.2 KB)           # AMQP message handlers
│   └── identity.py
│
├── utils/                             # Utilities
│   ├── db.py                          # SQLAlchemy helpers
│   ├── request.py                     # HTTP & rate limiting
│   ├── exceptions.py                  # Custom exceptions
│   ├── sentiment_analysis.py
│   └── extracted_keyword.py           # YAKE keyword extraction
│
├── keyword_collection/                # Keyword analysis
├── collection/                        # Collection management
├── clients/identity.py                # Identity service client
│
├── Configuration
│   ├── .env, .env.local, .env.prod, .env.stage
│   ├── schema.sql (13 KB)             # Database DDL
│   ├── .gitlab-ci.yml                 # GitLab CI/CD
│   └── scripts/start.sh               # Deployment script
└── requirements.txt                   # 128 dependencies
```

---

## 2. ARCHITECTURE DIAGRAM

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           BEAT SERVICE ARCHITECTURE                          │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│                              ENTRY POINTS                                    │
├───────────────────┬─────────────────────┬───────────────────────────────────┤
│   server.py       │     main.py         │    main_assets.py                 │
│   FastAPI API     │   Worker Pool       │    Asset Upload Workers           │
│   Port: 8000      │   73 Flows Config   │    S3/CDN Upload                  │
│   REST Endpoints  │   Multiprocessing   │    Media Caching                  │
└───────────┬───────┴──────────┬──────────┴───────────────┬───────────────────┘
            │                  │                          │
            ▼                  ▼                          ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                          WORKER POOL SYSTEM                                  │
├─────────────────────────────────────────────────────────────────────────────┤
│  Flow Name                    │  Workers  │  Concurrency  │  Description    │
│  ─────────────────────────────┼───────────┼───────────────┼─────────────────│
│  refresh_profile_by_handle    │    10     │      5        │  Instagram      │
│  refresh_yt_profiles          │    10     │      5        │  YouTube        │
│  asset_upload_flow            │    15     │      5        │  Media upload   │
│  refresh_post_insights        │     3     │      5        │  Post metrics   │
│  fetch_post_comments          │     1     │      5        │  Comments       │
│  ... 68 more flows            │   varies  │    varies     │  Various        │
└─────────────────────────────────────────────────────────────────────────────┘
            │
            ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                         RATE LIMITING LAYER                                  │
├─────────────────────────────────────────────────────────────────────────────┤
│  Source                │  Requests  │  Per Period  │  Implementation        │
│  ──────────────────────┼────────────┼──────────────┼────────────────────────│
│  youtube138            │    850     │   60 sec     │  asyncio-redis-rate    │
│  insta-best-performance│      2     │    1 sec     │                        │
│  arraybobo             │    100     │   30 sec     │                        │
│  youtubev31            │    500     │   60 sec     │                        │
│  rocketapi             │    100     │   30 sec     │                        │
│  Global Daily          │  20,000    │   86400 sec  │  Stacked limiters      │
│  Global Minute         │     60     │   60 sec     │                        │
└─────────────────────────────────────────────────────────────────────────────┘
            │
            ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                        API INTEGRATIONS (15+ APIs)                           │
├─────────────────────────────────────────────────────────────────────────────┤
│  INSTAGRAM (6 APIs)           │  YOUTUBE (4 APIs)        │  OTHER           │
│  ─────────────────────────────┼──────────────────────────┼──────────────────│
│  • Facebook Graph API v15     │  • YouTube Data API v3   │  • OpenAI GPT    │
│  • RapidAPI IGData            │  • RapidAPI YouTube v31  │  • Shopify API   │
│  • RapidAPI JoTucker          │  • RapidAPI YT Search    │  • Identity Svc  │
│  • RapidAPI NeoTank           │  • YouTube Analytics     │  • S3/CloudFront │
│  • RapidAPI ArrayBobo         │                          │                  │
│  • Lama API (fallback)        │                          │                  │
└─────────────────────────────────────────────────────────────────────────────┘
            │
            ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                         MESSAGE QUEUE (RabbitMQ/AMQP)                        │
├─────────────────────────────────────────────────────────────────────────────┤
│  Listener                 │  Exchange      │  Queue                │  Workers│
│  ─────────────────────────┼────────────────┼───────────────────────┼─────────│
│  credentials_validate     │  beat.dx       │  credentials_validate_q│    5   │
│  identity_token           │  identity.dx   │  identity_token_q      │    5   │
│  keyword_collection       │  beat.dx       │  keyword_collection_q  │    5   │
│  sentiment_analysis       │  beat.dx       │  sentiment_analysis_q  │    5   │
│  sentiment_report         │  beat.dx       │  sentiment_report_q    │    5   │
└─────────────────────────────────────────────────────────────────────────────┘
            │
            ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                          DATABASE LAYER                                      │
├─────────────────────────────────────────────────────────────────────────────┤
│  PostgreSQL (Async)        │  Redis Cluster            │  S3/CloudFront     │
│  ──────────────────────────┼───────────────────────────┼────────────────────│
│  • instagram_account       │  • Rate limit state       │  • Media assets    │
│  • instagram_post          │  • Cache layer            │  • CDN delivery    │
│  • youtube_account         │  • Session data           │                    │
│  • youtube_post            │                           │                    │
│  • scrape_request_log      │                           │                    │
│  • credential              │                           │                    │
│  • profile_log (audit)     │                           │                    │
│  • post_log (audit)        │                           │                    │
│  • sentiment_log           │                           │                    │
│  • asset_log               │                           │                    │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 3. FASTAPI APPLICATION (server.py - 43KB)

### Endpoint Categories

#### Profile Endpoints
```python
GET  /profiles/{platform}/byhandle/{handle}
     → Fetch profile by username (Instagram/YouTube)
     → Rate limited: 1/sec per handle, 60/min global, 20K/day

GET  /profiles/{platform}/byprofileid/{profile_id}
     → Fetch by platform-specific ID

GET  /profiles/{platform}/byid/{id}
     → Fetch by internal database ID

GET  /profiles/INSTAGRAM/byhandle/{handle}/insights
     → Instagram profile insights (business accounts)

GET  /profiles/INSTAGRAM/byhandle/{handle}/audienceinsights
     → Audience demographics
```

#### Post Endpoints
```python
GET  /posts/{platform}/byshortcode/{shortcode}
     → Single post details

GET  /posts/{platform}/{post_type}/{shortcode}
     → Post with type (image, carousel, reels, story)

GET  /recent/posts/{platform}/byprofileid/{profile_id}
     → Recent posts with pagination
```

#### Task Management
```python
POST /scrape_request_log/flow/{flow}
     → Create new scrape task (75+ flow types)
     → Body: {"params": {...}, "priority": 1}

POST /scrape_request_log/flow/update/{scrape_id}
     → Update task status

GET  /scrape_request_log/flow/{id}
     → Get task result

GET  /list_scrape_data/{account_id}
     → List tasks by account
```

#### Token Management
```python
POST /tokens
     → Insert/update API credentials

GET  /token/validate
     → Validate token scopes and expiry
```

#### Health
```python
GET  /heartbeat
     → Health check for load balancer

PUT  /heartbeat
     → Set health status (graceful shutdown)
```

### Rate Limiting Implementation

```python
# Stacked rate limiters for multi-level control
redis = AsyncRedis.from_url(REDIS_URL)

# Level 1: Daily global limit
global_limit_day = RateSpec(requests=20000, seconds=86400)

# Level 2: Per-minute global limit
global_limit_minute = RateSpec(requests=60, seconds=60)

# Level 3: Per-handle limit
handle_limit = RateSpec(requests=1, seconds=1)

async with RateLimiter(
    unique_key=f"refresh_profile_insta_daily",
    backend=redis,
    cache_prefix="beat_server_",
    rate_spec=global_limit_day
):
    async with RateLimiter(
        unique_key=f"refresh_profile_insta_minute",
        backend=redis,
        rate_spec=global_limit_minute
    ):
        async with RateLimiter(
            unique_key=f"refresh_profile_{handle}",
            backend=redis,
            rate_spec=handle_limit
        ):
            # Execute API call
            result = await refresh_profile(handle)
```

---

## 4. WORKER POOL SYSTEM (main.py - 14KB)

### Flow Configuration (73 Flows)

```python
_whitelist = {
    # Instagram Profile Flows
    'refresh_profile_custom': {'no_of_workers': 1, 'no_of_concurrency': 2},
    'refresh_profile_by_handle': {'no_of_workers': 10, 'no_of_concurrency': 5},
    'refresh_profile_by_profile_id': {'no_of_workers': 5, 'no_of_concurrency': 5},
    'refresh_profile_basic': {'no_of_workers': 3, 'no_of_concurrency': 5},

    # Instagram Post Flows
    'refresh_post_by_shortcode': {'no_of_workers': 3, 'no_of_concurrency': 5},
    'refresh_post_insights': {'no_of_workers': 3, 'no_of_concurrency': 5},
    'refresh_stories_posts': {'no_of_workers': 1, 'no_of_concurrency': 5},
    'refresh_story_insights': {'no_of_workers': 1, 'no_of_concurrency': 5},

    # Instagram Extra Flows
    'fetch_profile_followers': {'no_of_workers': 1, 'no_of_concurrency': 3},
    'fetch_profile_following': {'no_of_workers': 1, 'no_of_concurrency': 3},
    'fetch_post_comments': {'no_of_workers': 1, 'no_of_concurrency': 5},
    'fetch_post_likes': {'no_of_workers': 1, 'no_of_concurrency': 3},
    'fetch_hashtag_posts': {'no_of_workers': 1, 'no_of_concurrency': 3},

    # YouTube Flows
    'refresh_yt_profiles': {'no_of_workers': 10, 'no_of_concurrency': 5},
    'refresh_yt_posts': {'no_of_workers': 5, 'no_of_concurrency': 5},
    'refresh_yt_posts_by_channel_id': {'no_of_workers': 3, 'no_of_concurrency': 5},
    'fetch_yt_post_comments': {'no_of_workers': 1, 'no_of_concurrency': 3},

    # YouTube CSV Export Flows
    'fetch_channels_csv': {'no_of_workers': 2, 'no_of_concurrency': 3},
    'fetch_channel_videos_csv': {'no_of_workers': 2, 'no_of_concurrency': 3},
    'fetch_channel_demographic_csv': {'no_of_workers': 2, 'no_of_concurrency': 3},

    # GPT Enrichment Flows
    'refresh_instagram_gpt_data_base_gender': {'no_of_workers': 2, 'no_of_concurrency': 5},
    'refresh_instagram_gpt_data_base_location': {'no_of_workers': 2, 'no_of_concurrency': 5},
    'refresh_instagram_gpt_data_audience_age_gender': {'no_of_workers': 2, 'no_of_concurrency': 5},

    # Asset Upload
    'asset_upload_flow': {'no_of_workers': 15, 'no_of_concurrency': 5},
    'asset_upload_flow_stories': {'no_of_workers': 5, 'no_of_concurrency': 5},

    # Shopify
    'refresh_orders_by_store': {'no_of_workers': 2, 'no_of_concurrency': 3},

    # ... 40+ more flows
}
```

### Worker Architecture

```python
def main():
    """Main entry point - spawns worker processes"""
    for flow_name, config in _whitelist.items():
        for i in range(config['no_of_workers']):
            process = multiprocessing.Process(
                target=looper,
                args=(flow_name, config['no_of_concurrency'])
            )
            process.start()
            workers.append(process)

    # Start AMQP listeners
    start_amqp_listeners()

def looper(flow_name: str, concurrency: int):
    """Worker process entry point"""
    uvloop.install()  # High-performance event loop
    asyncio.run(poller(flow_name, concurrency))

async def poller(flow_name: str, concurrency: int):
    """Async polling loop"""
    semaphore = asyncio.Semaphore(concurrency)

    while True:
        task = await poll(flow_name)  # SQL-based task queue
        if task:
            asyncio.create_task(
                perform_task(task, semaphore)
            )
        await asyncio.sleep(0.1)

async def perform_task(task: ScrapeRequestLog, semaphore: Semaphore):
    """Execute task with concurrency control"""
    async with semaphore:
        try:
            async with asyncio.timeout(600):  # 10-minute timeout
                result = await execute(task.flow, task.params)
                await update_task_status(task.id, 'COMPLETE', result)
        except asyncio.TimeoutError:
            await update_task_status(task.id, 'TIMEOUT')
        except Exception as e:
            await update_task_status(task.id, 'FAILED', str(e))
```

### SQL-Based Task Queue

```python
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
            FOR UPDATE SKIP LOCKED
            LIMIT 1
        )
        RETURNING *
    """
    return await session.execute(query, {'flow': flow_name})
```

---

## 5. DATA COLLECTION FLOWS (75+ Flows)

### Instagram Flows (25+)

| Flow | Purpose | Workers | Concurrency |
|------|---------|---------|-------------|
| refresh_profile_by_handle | Profile lookup by username | 10 | 5 |
| refresh_profile_by_profile_id | Profile lookup by ID | 5 | 5 |
| refresh_profile_basic | Lightweight profile | 3 | 5 |
| refresh_profile_insights | Business insights | 3 | 5 |
| refresh_post_by_shortcode | Single post details | 3 | 5 |
| refresh_post_insights | Post metrics | 3 | 5 |
| refresh_stories_posts | Story content | 1 | 5 |
| refresh_story_insights | Story metrics | 1 | 5 |
| fetch_profile_followers | Follower list (paginated) | 1 | 3 |
| fetch_profile_following | Following list | 1 | 3 |
| fetch_post_comments | Post comments | 1 | 5 |
| fetch_post_likes | Post likers | 1 | 3 |
| fetch_hashtag_posts | Posts by hashtag | 1 | 3 |
| fetch_tagged_posts | Tagged posts | 1 | 3 |
| asset_upload_flow | Media CDN upload | 15 | 5 |

### YouTube Flows (20+)

| Flow | Purpose | Workers | Concurrency |
|------|---------|---------|-------------|
| refresh_yt_profiles | Channel info | 10 | 5 |
| refresh_yt_posts | Video list | 5 | 5 |
| refresh_yt_posts_by_channel_id | Videos by channel | 3 | 5 |
| refresh_yt_posts_by_playlist_id | Playlist videos | 2 | 3 |
| refresh_yt_profile_insights | Channel analytics | 2 | 3 |
| fetch_yt_post_comments | Video comments | 1 | 3 |
| fetch_channels_csv | Channel export | 2 | 3 |
| fetch_channel_videos_csv | Video export | 2 | 3 |
| fetch_channel_demographic_csv | Demographics export | 2 | 3 |
| fetch_channel_daily_stats_csv | Daily stats | 2 | 3 |
| fetch_channel_engagement_csv | Engagement metrics | 2 | 3 |

### GPT Enrichment Flows (6)

| Flow | Purpose | Workers | Concurrency |
|------|---------|---------|-------------|
| refresh_instagram_gpt_data_base_gender | Infer gender | 2 | 5 |
| refresh_instagram_gpt_data_base_location | Infer location | 2 | 5 |
| refresh_instagram_gpt_data_base_categ_lang_topics | Content analysis | 2 | 5 |
| refresh_instagram_gpt_data_audience_age_gender | Audience demographics | 2 | 5 |
| refresh_instagram_gpt_data_audience_cities | Geographic distribution | 2 | 5 |
| refresh_instagram_gpt_data_gender_location_lang | Combined analysis | 2 | 5 |

---

## 6. INSTAGRAM SCRAPING IMPLEMENTATIONS

### API Integration Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                 INSTAGRAM CRAWLER INTERFACE                      │
│                 (instagram/functions/retriever/interface.py)     │
├─────────────────────────────────────────────────────────────────┤
│  Abstract Methods:                                               │
│  - fetch_profile_by_handle(handle) → dict                       │
│  - fetch_profile_posts_by_handle(handle, limit) → list          │
│  - fetch_post_by_shortcode(shortcode) → dict                    │
│  - fetch_post_insights(post_id) → dict                          │
│  - fetch_profile_insights(user_id) → dict                       │
│  - fetch_stories_posts(user_id) → list                          │
│  - fetch_story_insights(story_id) → dict                        │
│  - fetch_followers(user_id, cursor) → (list, cursor)            │
│  - fetch_following(user_id, cursor) → (list, cursor)            │
│  - fetch_comments(post_id, cursor) → (list, cursor)             │
└─────────────────────────────────────────────────────────────────┘
                              │
           ┌──────────────────┼──────────────────┐
           ▼                  ▼                  ▼
┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐
│   GraphAPI      │  │   RapidAPI      │  │   Lama API      │
│   (Primary)     │  │   (Secondary)   │  │   (Fallback)    │
├─────────────────┤  ├─────────────────┤  ├─────────────────┤
│ • Official API  │  │ • IGData        │  │ • Post lookup   │
│ • Business accts│  │ • JoTucker      │  │ • No auth req   │
│ • Insights      │  │ • NeoTank       │  │ • Rate limited  │
│ • Stories       │  │ • ArrayBobo     │  │                 │
│                 │  │ • BestSolns     │  │                 │
│                 │  │ • RocketAPI     │  │                 │
└─────────────────┘  └─────────────────┘  └─────────────────┘
```

### GraphAPI Implementation (graphapi.py - 20KB)

```python
class GraphApi(InstagramCrawlerInterface):
    """Facebook Graph API implementation for Instagram Business accounts"""

    BASE_URL = "https://graph.facebook.com/v15.0"

    async def fetch_profile_by_handle(self, handle: str) -> dict:
        """
        Endpoint: /{user_id}/business_discovery.username(handle)
        Fields: biography, followers_count, follows_count, media_count,
                profile_picture_url, name, username, is_verified, etc.
        """
        fields = "biography,followers_count,follows_count,media_count,..."
        url = f"{self.BASE_URL}/{self.user_id}?fields=business_discovery.username({handle}){{{fields}}}"
        return await self._request(url)

    async def fetch_post_insights(self, post_id: str) -> dict:
        """
        Endpoint: /{post_id}/insights
        Metrics: impressions, reach, engagement, saved, video_views
        """
        metrics = "impressions,reach,engagement,saved,video_views"
        url = f"{self.BASE_URL}/{post_id}/insights?metric={metrics}"
        return await self._request(url)

    async def validate_token(self) -> bool:
        """
        Validate token scopes and expiry via debug_token endpoint
        Required scopes: instagram_basic, instagram_manage_insights,
                         pages_read_engagement, pages_show_list
        """
        url = f"{self.BASE_URL}/debug_token?input_token={self.token}"
        response = await self._request(url)
        scopes = response['data']['scopes']
        return all(s in scopes for s in REQUIRED_SCOPES)
```

### Data Flow Pipeline

```
┌──────────────────────────────────────────────────────────────────────┐
│                         DATA FLOW PIPELINE                            │
└──────────────────────────────────────────────────────────────────────┘

Stage 1: RETRIEVAL (instagram/tasks/retrieval.py)
┌──────────────────────────────────────────────────────────────────────┐
│  retrieve_profile_data_by_handle(handle, source)                     │
│    ↓                                                                 │
│  Select crawler based on source → Execute API call → Raw dict       │
└──────────────────────────────────────────────────────────────────────┘
                              ↓
Stage 2: PARSING (instagram/tasks/ingestion.py)
┌──────────────────────────────────────────────────────────────────────┐
│  parse_profile_data(raw_data, source)                                │
│    ↓                                                                 │
│  Extract fields → Create InstagramProfileLog                         │
│    - dimensions: [Dimension(key, value), ...]                        │
│    - metrics: [Metric(key, value), ...]                              │
└──────────────────────────────────────────────────────────────────────┘
                              ↓
Stage 3: PROCESSING (instagram/tasks/processing.py)
┌──────────────────────────────────────────────────────────────────────┐
│  upsert_profile(profile_log: InstagramProfileLog)                    │
│    ↓                                                                 │
│  Transform to InstagramAccount (ORM) → Upsert to PostgreSQL          │
│    ↓                                                                 │
│  Create ProfileLog entry (audit trail)                               │
│    ↓                                                                 │
│  Publish event to AMQP                                               │
└──────────────────────────────────────────────────────────────────────┘
```

---

## 7. DATABASE MODELS

### Core ORM Entities (SQLAlchemy)

```python
# core/entities/entities.py

class Credential(Base):
    """API credential storage"""
    __tablename__ = 'credential'

    id = Column(BigInteger, primary_key=True)
    idempotency_key = Column(String, unique=True)
    source = Column(String)  # graphapi, ytapi, rapidapi-igapi, etc.
    credentials = Column(JSONB)  # {token, user_id, key, refresh_token, ...}
    handle = Column(String)
    enabled = Column(Boolean, default=True)
    data_access_expired = Column(Boolean, default=False)
    disabled_till = Column(DateTime)  # TTL for rate limit backoff
    created_at = Column(DateTime, default=func.now())
    updated_at = Column(DateTime, onupdate=func.now())


class ScrapeRequestLog(Base):
    """Task queue table"""
    __tablename__ = 'scrape_request_log'

    id = Column(BigInteger, primary_key=True)
    idempotency_key = Column(String, unique=True)
    platform = Column(String)  # INSTAGRAM, YOUTUBE, SHOPIFY
    flow = Column(String)  # 75+ flow names
    status = Column(String)  # PENDING, PROCESSING, COMPLETE, FAILED
    params = Column(JSONB)  # Flow-specific parameters
    data = Column(Text)  # Result or error message
    priority = Column(Integer, default=1)
    retry_count = Column(Integer, default=0)
    account_id = Column(String)
    created_at = Column(DateTime, default=func.now())
    picked_at = Column(DateTime)
    scraped_at = Column(DateTime)
    expires_at = Column(DateTime)


class ProfileLog(Base):
    """Audit log for profile snapshots"""
    __tablename__ = 'profile_log'

    id = Column(BigInteger, primary_key=True)
    platform = Column(String)
    profile_id = Column(String)
    dimensions = Column(JSONB)  # [{key, value}, ...]
    metrics = Column(JSONB)  # [{key, value}, ...]
    source = Column(String)
    timestamp = Column(DateTime, default=func.now())
```

### Instagram ORM Models

```python
# instagram/entities/entities.py

class InstagramAccount(Base):
    """Instagram profile data"""
    __tablename__ = 'instagram_account'

    id = Column(BigInteger, primary_key=True)
    profile_id = Column(String, unique=True)  # Instagram user ID
    handle = Column(String, index=True)
    full_name = Column(String)
    biography = Column(Text)

    # Metrics
    followers = Column(BigInteger)
    following = Column(BigInteger)
    media_count = Column(BigInteger)

    # Calculated metrics
    avg_likes = Column(Float)
    avg_comments = Column(Float)
    avg_reach = Column(Float)
    avg_engagement = Column(Float)
    avg_reels_plays = Column(Float)

    # Profile attributes
    profile_pic_url = Column(String)
    profile_type = Column(String)  # personal, business, creator
    is_private = Column(Boolean)
    is_verified = Column(Boolean)
    is_business_or_creator = Column(Boolean)
    category = Column(String)
    fbid = Column(String)  # Facebook page ID

    # Timestamps
    created_at = Column(DateTime, default=func.now())
    updated_at = Column(DateTime, onupdate=func.now())
    refreshed_at = Column(DateTime)


class InstagramPost(Base):
    """Instagram post data"""
    __tablename__ = 'instagram_post'

    id = Column(BigInteger, primary_key=True)
    post_id = Column(String, unique=True)  # Instagram media ID
    shortcode = Column(String, unique=True, index=True)
    profile_id = Column(String, ForeignKey('instagram_account.profile_id'))
    handle = Column(String)

    # Content
    post_type = Column(String)  # image, carousel, reels, story
    caption = Column(Text)
    thumbnail_url = Column(String)
    display_url = Column(String)

    # Metrics
    likes = Column(BigInteger)
    comments = Column(BigInteger)
    plays = Column(BigInteger)  # For reels/videos
    reach = Column(BigInteger)
    views = Column(BigInteger)
    shares = Column(BigInteger)
    impressions = Column(BigInteger)
    saved = Column(BigInteger)

    # Timestamps
    publish_time = Column(DateTime)
    created_at = Column(DateTime, default=func.now())
    updated_at = Column(DateTime, onupdate=func.now())
```

---

## 8. GPT/OPENAI INTEGRATION

### Architecture

```python
# gpt/functions/retriever/openai/openai_extractor.py

class OpenAi(GptCrawlerInterface):
    """OpenAI API integration for data enrichment"""

    def __init__(self):
        openai.api_type = "azure"
        openai.api_base = os.environ["OPENAI_API_BASE"]
        openai.api_version = os.environ["OPENAI_API_VERSION"]
        openai.api_key = os.environ["OPENAI_API_KEY"]

    async def fetch_instagram_gpt_data_base_gender(
        self, handle: str, bio: str
    ) -> dict:
        """Infer audience gender from bio and handle"""
        prompt = self._load_prompt("profile_info_v0.12.yaml")
        response = await openai.ChatCompletion.acreate(
            engine="gpt-35-turbo",
            messages=[
                {"role": "system", "content": prompt['system']},
                {"role": "user", "content": prompt['user'].format(
                    handle=handle, bio=bio
                )}
            ],
            temperature=0  # Deterministic output
        )
        return self._parse_response(response)

    async def fetch_instagram_gpt_data_audience_age_gender(
        self, handle: str, bio: str, recent_posts: list
    ) -> dict:
        """Infer audience demographics from content"""
        # Analyze bio + recent post captions
        content = f"{bio}\n\n" + "\n".join([p['caption'] for p in recent_posts])
        # ... similar implementation
```

### Prompt Management

```yaml
# gpt/prompts/profile_info_v0.12.yaml

system: |
  You are an AI assistant that analyzes Instagram profiles.
  Given a username and bio, infer the following:
  - Primary audience gender (male/female/mixed)
  - Confidence score (0-1)

user: |
  Username: {handle}
  Bio: {bio}

  Analyze and respond in JSON format:
  {
    "gender": "male|female|mixed",
    "confidence": 0.85,
    "reasoning": "..."
  }

model: gpt-35-turbo
temperature: 0
```

### Use Cases

| Flow | Input | Output |
|------|-------|--------|
| base_gender | handle, bio | {gender, confidence} |
| base_location | handle, bio | {country, city, confidence} |
| categ_lang_topics | handle, bio, posts | {category, language, topics[]} |
| audience_age_gender | handle, bio, posts | {age_range, gender_dist} |
| audience_cities | handle, bio, posts | {cities: [{name, percentage}]} |

---

## 9. MESSAGE QUEUE INTEGRATION (RabbitMQ/AMQP)

### aio-pika Implementation

```python
# core/amqp/amqp.py

@dataclass
class AmqpListener:
    """AMQP listener configuration"""
    exchange: str
    routing_key: str
    queue: str
    workers: int
    prefetch: int
    fn: Callable  # Handler function


async def async_listener(config: AmqpListener):
    """Start AMQP listener with connection recovery"""
    connection = await aio_pika.connect_robust(
        os.environ["RMQ_URL"],
        heartbeat=60
    )
    channel = await connection.channel()
    await channel.set_qos(prefetch_count=config.prefetch)

    exchange = await channel.declare_exchange(
        config.exchange, ExchangeType.DIRECT, durable=True
    )
    queue = await channel.declare_queue(config.queue, durable=True)
    await queue.bind(exchange, routing_key=config.routing_key)

    async with queue.iterator() as queue_iter:
        async for message in queue_iter:
            async with message.process():
                try:
                    await config.fn(message.body)
                except Exception as e:
                    logger.error(f"Message processing failed: {e}")
                    # Message will be requeued
                    raise
```

### Configured Listeners

```python
# main.py

amqp_listeners = [
    AmqpListener(
        exchange="beat.dx",
        routing_key="credentials_validate_rk",
        queue="credentials_validate_q",
        workers=5,
        prefetch=10,
        fn=credential_validate
    ),
    AmqpListener(
        exchange="identity.dx",
        routing_key="new_access_token_rk",
        queue="identity_token_q",
        workers=5,
        prefetch=10,
        fn=upsert_credential_from_identity
    ),
    AmqpListener(
        exchange="beat.dx",
        routing_key="keyword_collection_rk",
        queue="keyword_collection_q",
        workers=5,
        prefetch=1,
        fn=fetch_keyword_collection
    ),
    AmqpListener(
        exchange="beat.dx",
        routing_key="post_activity_log_bulk_rk",
        queue="sentiment_analysis_q",
        workers=5,
        prefetch=1,
        fn=sentiment_extraction
    ),
    AmqpListener(
        exchange="beat.dx",
        routing_key="sentiment_collection_report_in_rk",
        queue="sentiment_report_q",
        workers=5,
        prefetch=1,
        fn=fetch_sentiment_report
    ),
]
```

---

## 10. CREDENTIAL MANAGEMENT

### Credential Manager

```python
# credentials/manager.py

class CredentialManager:
    """Manage API credentials lifecycle"""

    async def insert_creds(
        self, source: str, credentials: dict, handle: str = None
    ) -> Credential:
        """Upsert credential by idempotency key"""
        key = f"{source}:{credentials.get('user_id', credentials.get('key'))}"
        return await get_or_create(
            session,
            Credential,
            idempotency_key=key,
            defaults={
                'source': source,
                'credentials': credentials,
                'handle': handle,
                'enabled': True
            }
        )

    async def disable_creds(
        self, cred_id: int, disable_duration: int = 3600
    ) -> None:
        """Disable credential with TTL (rate limit backoff)"""
        await session.execute(
            update(Credential)
            .where(Credential.id == cred_id)
            .values(
                enabled=False,
                disabled_till=func.now() + timedelta(seconds=disable_duration)
            )
        )

    async def get_enabled_cred(self, source: str) -> Optional[Credential]:
        """Get random enabled credential for source"""
        creds = await session.execute(
            select(Credential)
            .where(Credential.source == source)
            .where(Credential.enabled == True)
            .where(
                or_(
                    Credential.disabled_till.is_(None),
                    Credential.disabled_till < func.now()
                )
            )
        )
        enabled_creds = creds.scalars().all()
        return random.choice(enabled_creds) if enabled_creds else None
```

### Credential Validator

```python
# credentials/validator.py

REQUIRED_SCOPES = [
    'instagram_basic',
    'instagram_manage_insights',
    'pages_read_engagement',
    'pages_show_list'
]

class CredentialValidator:
    """Validate API credentials"""

    async def validate(self, credential: Credential) -> ValidationResult:
        """Validate token and check required scopes"""
        if credential.source == 'graphapi':
            return await self._validate_graphapi(credential)
        elif credential.source == 'ytapi':
            return await self._validate_youtube(credential)
        # ... other sources

    async def _validate_graphapi(self, cred: Credential) -> ValidationResult:
        token = cred.credentials.get('token')
        response = await self._debug_token(token)

        if not response['data']['is_valid']:
            raise TokenValidationFailed("Token is invalid")

        scopes = response['data']['scopes']
        missing = [s for s in REQUIRED_SCOPES if s not in scopes]
        if missing:
            raise TokenValidationFailed(f"Missing scopes: {missing}")

        if response['data'].get('data_access_expires_at', 0) < time.time():
            raise DataAccessExpired("Data access has expired")

        return ValidationResult(valid=True, scopes=scopes)
```

---

## 11. ENGAGEMENT CALCULATIONS

### Formula-Based Analytics

```python
# instagram/helper.py

def calculate_engagement_rate(
    likes: int, comments: int, followers: int
) -> float:
    """Standard engagement rate formula"""
    if followers == 0:
        return 0.0
    return ((likes + comments) / followers) * 100


def estimate_reach_reels(plays: int, followers: int) -> float:
    """Estimate reach for Reels based on plays and followers"""
    # Empirical formula from platform data analysis
    factor = 0.94 - (math.log2(followers) * 0.001)
    return plays * factor


def estimate_reach_posts(likes: int) -> float:
    """Estimate reach for static posts based on likes"""
    if likes == 0:
        return 0.0
    factor = (7.6 - (math.log10(likes) * 0.7)) * 0.85
    return factor * likes


def estimate_story_reach(
    followers: int, avg_engagement: float
) -> float:
    """Estimate story reach based on followers and engagement"""
    base = -0.000025017 * followers
    engagement_factor = 1.11 * (
        followers * abs(math.log2(avg_engagement + 2)) * 2 / 100
    )
    return base + engagement_factor


def calculate_avg_metrics(posts: list, exclude_outliers: bool = True) -> dict:
    """Calculate average metrics with optional outlier removal"""
    if exclude_outliers and len(posts) > 4:
        # Remove top 2 and bottom 2 by engagement
        posts = sorted(posts, key=lambda p: p['likes'] + p['comments'])
        posts = posts[2:-2]

    return {
        'avg_likes': mean([p['likes'] for p in posts]),
        'avg_comments': mean([p['comments'] for p in posts]),
        'avg_reach': mean([p.get('reach', 0) for p in posts]),
        'avg_engagement': mean([
            calculate_engagement_rate(p['likes'], p['comments'], p['followers'])
            for p in posts
        ])
    }
```

---

## 12. CI/CD PIPELINE

### GitLab CI Configuration

```yaml
# .gitlab-ci.yml

stages:
  - deploy_stage
  - deploy_prod

deploy_stage:
  stage: deploy_stage
  script:
    - /bin/bash scripts/start.sh
  tags:
    - stage-aafat
  only:
    - master
    - dev

deploy_prod_1:
  stage: deploy_prod
  script:
    - /bin/bash scripts/start.sh
  tags:
    - beat-deployer-1
  when: manual
  only:
    - master

deploy_prod_2:
  stage: deploy_prod
  script:
    - /bin/bash scripts/start.sh
  tags:
    - beat-deployer-2
  when: manual
  only:
    - master
```

### Deployment Script

```bash
#!/bin/bash
# scripts/start.sh

# Setup Python environment
python3.11 -m venv venv
source venv/bin/activate
pip install -r requirements.txt

# Copy environment config
cp .env.${ENV} .env

# Graceful shutdown
curl -XPUT http://localhost:8000/heartbeat?beat=false
sleep 10

# Kill existing processes
pkill -f "python main.py" || true
pkill -f "python server.py" || true
sleep 5

# Start workers
nohup python main.py > logs/main.log 2>&1 &

# Start API server
nohup python server.py > logs/server.log 2>&1 &

# Wait for startup
sleep 10

# Enable health check
curl -XPUT http://localhost:8000/heartbeat?beat=true
```

---

## 13. KEY METRICS & STATISTICS

| Metric | Value |
|--------|-------|
| **Total Lines of Code** | ~15,000+ |
| **Python Modules** | 130+ |
| **Data Collection Flows** | 75+ |
| **API Integrations** | 15+ |
| **Dependencies** | 128 packages |
| **Worker Processes** | 20-180+ (configurable) |
| **Concurrency per Worker** | 2-15 (configurable) |
| **AMQP Listeners** | 5 |
| **Rate Limit Rules** | 7+ source-specific |
| **Database Tables** | 30+ |

---

## 14. SKILLS DEMONSTRATED

### Technical Skills

| Skill | Evidence |
|-------|----------|
| **Async Python** | FastAPI + uvloop + aio-pika + asyncpg |
| **Distributed Systems** | Multiprocessing workers + message queues |
| **API Integration** | 15+ external APIs with fallback strategies |
| **Rate Limiting** | Redis-backed multi-level rate limiting |
| **Database Design** | PostgreSQL with JSONB, async sessions |
| **ML Integration** | OpenAI GPT for data enrichment |
| **Task Queues** | SQL-based with FOR UPDATE SKIP LOCKED |
| **Data Pipelines** | 3-stage ETL (Retrieval → Parsing → Processing) |

### Architecture Patterns

1. **Worker Pool Pattern** - Multiprocessing with async I/O
2. **Semaphore Pattern** - Concurrency control per worker
3. **Decorator Pattern** - @sessionize for session injection
4. **Strategy Pattern** - Multiple API implementations per interface
5. **Circuit Breaker** - Credential disable with TTL backoff
6. **Fallback Pattern** - Secondary APIs when primary fails

---

## 15. INTERVIEW TALKING POINTS

### System Design Questions

**"Design a distributed data scraping system"**
- Worker pool with 73 configurable flows
- SQL-based task queue with FOR UPDATE SKIP LOCKED
- Multi-level rate limiting (global daily, per-minute, per-resource)
- 15+ API integrations with fallback strategies
- Credential rotation for load balancing

**"How do you handle rate limits from external APIs?"**
- Redis-backed asyncio-redis-rate-limit
- Source-specific rate specs (2-850 requests/period)
- Stacked limiters for multi-level control
- Credential disable with TTL backoff
- Multiple API sources for redundancy

**"Explain your async Python architecture"**
- FastAPI + uvloop for high-performance event loop
- aio-pika for async RabbitMQ consumption
- asyncpg for async PostgreSQL
- Semaphore-based concurrency control
- 10-minute task timeout with auto-cancellation

### Behavioral Questions

**"Tell me about a complex data pipeline you built"**
- Beat: 15K+ LOC, 75+ flows, 15+ API integrations
- 3-stage pipeline: Retrieval → Parsing → Processing
- GPT integration for data enrichment
- Real-time event publishing to AMQP

**"How do you handle API failures?"**
- Fallback APIs (e.g., Lama for Instagram when GraphAPI fails)
- Retry with exponential backoff
- Credential rotation on rate limit
- Circuit breaker with TTL disable

---

*Generated through comprehensive source code analysis of the beat project.*

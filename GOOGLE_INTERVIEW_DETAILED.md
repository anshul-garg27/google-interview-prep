# GOOGLE INTERVIEW - DETAILED COMPONENT-WISE BREAKDOWN

---

# PROJECT 1: BEAT (Data Aggregation Service)

## BEAT has 12 Major Components You Built:

| # | Component | Lines of Code | Complexity | Interview Story Potential |
|---|-----------|---------------|------------|---------------------------|
| 1 | Worker Pool System | main.py (14KB) | High | System Design |
| 2 | SQL-Based Task Queue | core/flows/ | Medium | Distributed Systems |
| 3 | Multi-Level Rate Limiting | server.py, utils/request.py | High | Scalability |
| 4 | API Integration Framework | instagram/functions/retriever/ | High | Design Patterns |
| 5 | Credential Management | credentials/ | Medium | Security |
| 6 | 3-Stage Data Pipeline | tasks/ (retrieval→parsing→processing) | High | Data Engineering |
| 7 | GPT/OpenAI Integration | gpt/ | Medium | AI/ML |
| 8 | RabbitMQ/AMQP Listeners | core/amqp/ | Medium | Event-Driven |
| 9 | Asset Upload System | main_assets.py, client/ | Medium | CDN/Storage |
| 10 | Engagement Calculations | instagram/helper.py | Low | Analytics |
| 11 | FastAPI REST API | server.py (43KB) | Medium | API Design |
| 12 | Graceful Deployment | scripts/start.sh | Low | DevOps |

---

## COMPONENT 1: Worker Pool System (main.py)

### What You Built
A distributed worker pool system that spawns multiple processes, each running async event loops with semaphore-based concurrency control.

### Technical Deep Dive

```python
# Architecture: Multiprocessing + Asyncio + Semaphore

def main():
    """Entry point - spawns 150+ workers across 73 flows"""
    for flow_name, config in _whitelist.items():
        for i in range(config['no_of_workers']):
            process = multiprocessing.Process(
                target=looper,
                args=(flow_name, config['no_of_concurrency'])
            )
            process.start()
            workers.append(process)

def looper(flow_name: str, concurrency: int):
    """Each process has its own async event loop"""
    uvloop.install()  # 2-4x faster than default asyncio
    asyncio.run(poller(flow_name, concurrency))

async def poller(flow_name: str, concurrency: int):
    """Async polling with semaphore-based concurrency"""
    semaphore = asyncio.Semaphore(concurrency)

    while True:
        task = await poll(flow_name)  # SQL-based queue
        if task:
            asyncio.create_task(perform_task(task, semaphore))
        await asyncio.sleep(0.1)  # Prevent busy-waiting

async def perform_task(task, semaphore):
    """Execute with concurrency control + timeout"""
    async with semaphore:
        try:
            async with asyncio.timeout(600):  # 10-min timeout
                result = await execute(task.flow, task.params)
                await update_task_status(task.id, 'COMPLETE', result)
        except asyncio.TimeoutError:
            await update_task_status(task.id, 'TIMEOUT')
        except Exception as e:
            await update_task_status(task.id, 'FAILED', str(e))
```

### Why This Design?

| Decision | Why | Alternative Considered |
|----------|-----|----------------------|
| **Multiprocessing** | Python GIL limits CPU-bound parallelism | Threads (blocked by GIL) |
| **Asyncio inside each process** | I/O-bound work (API calls, DB) benefits from async | Sync requests (slow) |
| **Semaphore** | Prevent overwhelming external APIs | No limit (429 errors) |
| **uvloop** | 2-4x faster event loop | Default asyncio (slower) |
| **SQL polling** | Simple, no extra infrastructure | Celery (complex setup) |

### Interview STAR Story

**Situation**: We needed to collect data from 10M+ social media profiles daily, but external APIs had strict rate limits.

**Task**: Design a system that maximizes throughput while respecting rate limits and handling failures gracefully.

**Action**:
1. Designed hybrid architecture: Multiprocessing for parallelism + Asyncio for I/O concurrency
2. Used semaphores to control per-flow concurrency (e.g., 5 concurrent API calls per worker)
3. Implemented 10-minute timeout to prevent stuck tasks
4. Added graceful error handling with automatic retry via task queue

**Result**:
- **150+ workers** running concurrently
- **10M+ daily data points** processed
- **99.9% uptime** with automatic recovery
- **25% faster** than previous sync implementation

### Questions They Might Ask

**Q: Why not use Celery?**
A: Celery adds complexity (Redis/RabbitMQ broker, beat scheduler, multiple processes). For our use case, a simple SQL-based queue was sufficient and easier to debug. We already had PostgreSQL, so no new infrastructure needed.

**Q: How do you handle worker crashes?**
A: Tasks remain in `PROCESSING` status. A separate cleanup cron resets tasks stuck in `PROCESSING` for >15 minutes back to `PENDING`. The task's `retry_count` is incremented.

**Q: Why 10-minute timeout?**
A: Some flows (like fetching 1000 followers with pagination) legitimately take 5-8 minutes. 10 minutes gives buffer while catching truly stuck tasks.

---

## COMPONENT 2: SQL-Based Task Queue

### What You Built
A distributed task queue using PostgreSQL with `FOR UPDATE SKIP LOCKED` for concurrent worker coordination.

### Technical Deep Dive

```python
async def poll(flow_name: str) -> Optional[ScrapeRequestLog]:
    """
    Atomic task pickup with row-level locking

    Key SQL features:
    - FOR UPDATE: Lock the row to prevent double-pickup
    - SKIP LOCKED: Don't wait, skip to next available row
    - Priority ordering: High-priority tasks first
    - Expiry check: Skip expired tasks
    """
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

### Schema Design

```sql
CREATE TABLE scrape_request_log (
    id BIGSERIAL PRIMARY KEY,
    idempotency_key VARCHAR(255) UNIQUE,  -- Prevent duplicate tasks
    platform VARCHAR(50),                  -- INSTAGRAM, YOUTUBE, SHOPIFY
    flow VARCHAR(100),                     -- 73 flow types
    status VARCHAR(20) DEFAULT 'PENDING',  -- PENDING, PROCESSING, COMPLETE, FAILED
    params JSONB,                          -- Flow-specific parameters
    data TEXT,                             -- Result or error message
    priority INTEGER DEFAULT 1,            -- Higher = processed first
    retry_count INTEGER DEFAULT 0,
    account_id VARCHAR(100),               -- For grouping/filtering
    created_at TIMESTAMP DEFAULT NOW(),
    picked_at TIMESTAMP,
    scraped_at TIMESTAMP,
    expires_at TIMESTAMP                   -- Auto-skip if expired
);

-- Indexes for performance
CREATE INDEX idx_flow_status ON scrape_request_log(flow, status);
CREATE INDEX idx_priority_created ON scrape_request_log(priority DESC, created_at ASC);
```

### Why This Design?

| Feature | Purpose |
|---------|---------|
| **FOR UPDATE SKIP LOCKED** | Multiple workers can poll simultaneously without blocking |
| **idempotency_key** | Prevent duplicate task creation (e.g., same profile scraped twice) |
| **priority** | Business-critical profiles processed first |
| **expires_at** | Don't process stale requests |
| **JSONB params** | Flexible flow-specific parameters |

### Interview STAR Story

**Situation**: Needed a task queue for 73 different flows with 150+ workers, but Celery felt too heavy.

**Task**: Build a lightweight, reliable task queue using existing PostgreSQL.

**Action**:
1. Designed schema with proper indexes for fast polling
2. Used `FOR UPDATE SKIP LOCKED` for concurrent-safe task pickup
3. Added idempotency keys to prevent duplicate tasks
4. Implemented priority-based ordering for business-critical tasks
5. Created cleanup cron for stuck tasks

**Result**:
- **Sub-millisecond** task pickup latency
- **Zero duplicate** task processing
- **No additional infrastructure** needed
- **Easy debugging** - just query the table

### Questions They Might Ask

**Q: What's the throughput of this queue?**
A: With proper indexes, we achieved ~1000 task pickups/second. The bottleneck was API rate limits, not the queue.

**Q: How do you handle task failures?**
A: Failed tasks get `status='FAILED'` with error in `data` column. A separate process can retry failed tasks or alert on-call.

**Q: Why not Redis-based queue?**
A: PostgreSQL was already our primary datastore. Adding Redis would mean:
- Another service to manage
- Data consistency issues (task in Redis, result in Postgres)
- For our scale (~1000 tasks/sec), Postgres was sufficient

---

## COMPONENT 3: Multi-Level Rate Limiting

### What You Built
A stacked rate limiting system using Redis that enforces limits at multiple levels simultaneously.

### Technical Deep Dive

```python
# 3-Level Stacked Rate Limiting

from asyncio_redis_rate_limit import RateLimiter, RateSpec

# Configuration per API source
RATE_LIMITS = {
    'graphapi': RateSpec(requests=200, seconds=3600),      # 200/hour
    'youtube138': RateSpec(requests=850, seconds=60),      # 850/minute
    'insta-best-performance': RateSpec(requests=2, seconds=1),  # 2/second
    'arraybobo': RateSpec(requests=100, seconds=30),       # 100/30sec
    'rocketapi': RateSpec(requests=100, seconds=30),
}

# Global limits
GLOBAL_DAILY = RateSpec(requests=20000, seconds=86400)
GLOBAL_MINUTE = RateSpec(requests=60, seconds=60)

async def rate_limited_request(handle: str, source: str):
    """
    Stacked limiters - ALL must pass before request proceeds

    Level 1: Global daily (20K/day) - prevent runaway costs
    Level 2: Global per-minute (60/min) - smooth traffic
    Level 3: Per-handle (1/sec) - prevent hammering same profile
    Level 4: Per-source (varies) - respect API-specific limits
    """
    redis = AsyncRedis.from_url(REDIS_URL)

    async with RateLimiter(
        unique_key="beat_global_daily",
        backend=redis,
        rate_spec=GLOBAL_DAILY
    ):
        async with RateLimiter(
            unique_key="beat_global_minute",
            backend=redis,
            rate_spec=GLOBAL_MINUTE
        ):
            async with RateLimiter(
                unique_key=f"beat_handle_{handle}",
                backend=redis,
                rate_spec=RateSpec(requests=1, seconds=1)
            ):
                async with RateLimiter(
                    unique_key=f"beat_source_{source}",
                    backend=redis,
                    rate_spec=RATE_LIMITS[source]
                ):
                    return await make_api_call(handle, source)
```

### Redis Data Structure

```
# Sliding window counter pattern
Key: beat_server_beat_global_daily
Value: {
    "count": 15234,
    "window_start": 1706745600
}

Key: beat_server_beat_handle_virat.kohli
Value: {
    "count": 1,
    "window_start": 1706832000
}
```

### Why This Design?

| Level | Purpose | Limit |
|-------|---------|-------|
| **Global Daily** | Cost control - don't exceed API budget | 20K/day |
| **Global Minute** | Traffic smoothing - prevent bursts | 60/min |
| **Per-Handle** | Prevent hammering same profile | 1/sec |
| **Per-Source** | Respect each API's specific limits | Varies |

### Interview STAR Story

**Situation**: External APIs (Instagram, YouTube) have strict rate limits. Exceeding them results in 429 errors, temporary bans, or permanent API key revocation.

**Task**: Implement rate limiting that respects all API limits while maximizing throughput.

**Action**:
1. Analyzed each API's rate limit documentation
2. Implemented stacked limiters - request must pass ALL levels
3. Used Redis for distributed state (multiple workers share limits)
4. Added per-source configuration for easy adjustment
5. Implemented automatic backoff on 429 responses

**Result**:
- **Zero API bans** since implementation
- **30% cost reduction** by staying within quota
- **Smooth traffic** distribution throughout the day
- **Easy tuning** - just update config dict

### Questions They Might Ask

**Q: What happens when rate limit is exceeded?**
A: The `RateLimiter` context manager blocks (async wait) until the window resets. Alternatively, we can raise an exception and retry later.

**Q: How do you handle different time windows?**
A: Redis sliding window algorithm. Each key stores count + window_start. On request:
1. If current_time > window_start + window_size: reset count
2. If count < limit: increment and proceed
3. Else: wait or reject

**Q: What if Redis goes down?**
A: Graceful degradation - we fall back to in-memory rate limiting per worker. Less accurate, but prevents complete failure.

---

## COMPONENT 4: API Integration Framework (Strategy Pattern)

### What You Built
A pluggable API integration framework using the Strategy pattern, allowing easy addition of new data sources.

### Technical Deep Dive

```python
# Interface Definition
class InstagramCrawlerInterface(ABC):
    """Abstract interface for Instagram data retrieval"""

    @abstractmethod
    async def fetch_profile_by_handle(self, handle: str) -> dict:
        pass

    @abstractmethod
    async def fetch_profile_posts_by_handle(self, handle: str, limit: int) -> list:
        pass

    @abstractmethod
    async def fetch_post_by_shortcode(self, shortcode: str) -> dict:
        pass

    @abstractmethod
    async def fetch_post_insights(self, post_id: str) -> dict:
        pass

    @abstractmethod
    async def fetch_followers(self, user_id: str, cursor: str) -> Tuple[list, str]:
        pass


# Implementation 1: Facebook Graph API (Official)
class GraphApi(InstagramCrawlerInterface):
    BASE_URL = "https://graph.facebook.com/v15.0"

    async def fetch_profile_by_handle(self, handle: str) -> dict:
        fields = "biography,followers_count,follows_count,media_count,..."
        url = f"{self.BASE_URL}/{self.user_id}?fields=business_discovery.username({handle}){{{fields}}}"
        return await self._request(url)


# Implementation 2: RapidAPI IGData
class RapidApiIGData(InstagramCrawlerInterface):
    BASE_URL = "https://instagram-data1.p.rapidapi.com"

    async def fetch_profile_by_handle(self, handle: str) -> dict:
        url = f"{self.BASE_URL}/user/info?username={handle}"
        headers = {"X-RapidAPI-Key": self.api_key}
        return await self._request(url, headers)


# Implementation 3: Lama API (Fallback)
class LamaApi(InstagramCrawlerInterface):
    # ... minimal implementation for fallback


# Factory/Selector
def get_crawler(source: str) -> InstagramCrawlerInterface:
    crawlers = {
        'graphapi': GraphApi,
        'rapidapi-igdata': RapidApiIGData,
        'rapidapi-jotucker': RapidApiJoTucker,
        'rapidapi-neotank': RapidApiNeoTank,
        'rapidapi-arraybobo': RapidApiArrayBobo,
        'lama': LamaApi,
    }
    return crawlers[source]()
```

### Fallback Strategy

```python
async def fetch_profile_with_fallback(handle: str) -> dict:
    """
    Try sources in order of reliability/cost
    1. GraphAPI (official, best data quality)
    2. RapidAPI options (paid, good quality)
    3. Lama (free, limited data)
    """
    sources = ['graphapi', 'rapidapi-igdata', 'rapidapi-jotucker', 'lama']

    for source in sources:
        try:
            crawler = get_crawler(source)
            cred = await credential_manager.get_enabled_cred(source)
            if not cred:
                continue  # No available credentials

            crawler.set_credentials(cred)
            result = await crawler.fetch_profile_by_handle(handle)

            if result:
                return result

        except RateLimitError:
            # Disable this credential temporarily
            await credential_manager.disable_creds(cred.id, 3600)
            continue

        except Exception as e:
            logger.error(f"Source {source} failed: {e}")
            continue

    raise AllSourcesFailedError(f"Could not fetch {handle}")
```

### Why This Design?

| Pattern | Benefit |
|---------|---------|
| **Strategy Pattern** | Easy to add new APIs without changing core logic |
| **Interface** | All crawlers have same method signatures |
| **Factory** | Single point of crawler instantiation |
| **Fallback Chain** | Reliability - if one fails, try next |

### Interview STAR Story

**Situation**: We needed to fetch Instagram data, but no single API was reliable enough. GraphAPI requires business account connection, RapidAPI has rate limits, etc.

**Task**: Design a system that can use multiple data sources with easy fallback.

**Action**:
1. Defined abstract interface with all required methods
2. Implemented 6 different API integrations following the interface
3. Created factory function for crawler selection
4. Built fallback chain that tries sources in priority order
5. Integrated with credential manager for API key rotation

**Result**:
- **6 Instagram APIs** integrated seamlessly
- **99.5% success rate** with fallback chain
- **New API integration** takes ~2 hours (just implement interface)
- **Clean separation** between API logic and business logic

### Questions They Might Ask

**Q: How do you decide source priority?**
A: Based on data quality, cost, and reliability:
1. GraphAPI - Best quality, free (but limited to business accounts)
2. RapidAPI IGData - Good quality, $50/month
3. RapidAPI JoTucker - Good, cheaper
4. Lama - Free but limited fields

**Q: What if all sources fail?**
A: Task marked as FAILED. A separate alerting system notifies on-call if failure rate exceeds threshold (>5% in 1 hour).

**Q: How do you handle different response formats?**
A: Each crawler has a `_parser.py` file that normalizes responses to our internal schema. The interface guarantees output format.

---

## COMPONENT 5: Credential Management System

### What You Built
A credential lifecycle management system with automatic rotation, validation, and TTL-based backoff.

### Technical Deep Dive

```python
# credentials/manager.py

class CredentialManager:
    """
    Manages API credentials across multiple sources

    Features:
    - Upsert with idempotency
    - TTL-based disable (rate limit backoff)
    - Random selection for load balancing
    - Automatic re-enable after TTL
    """

    async def insert_creds(self, source: str, credentials: dict,
                          handle: str = None) -> Credential:
        """Upsert credential with idempotency"""
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

    async def disable_creds(self, cred_id: int,
                           disable_duration: int = 3600) -> None:
        """
        Disable credential with TTL
        Used when API returns 429 (rate limit) or 401 (token expired)
        """
        await session.execute(
            update(Credential)
            .where(Credential.id == cred_id)
            .values(
                enabled=False,
                disabled_till=func.now() + timedelta(seconds=disable_duration)
            )
        )

    async def get_enabled_cred(self, source: str) -> Optional[Credential]:
        """
        Get random enabled credential for load balancing

        Checks:
        1. Source matches
        2. enabled=True
        3. Either no disabled_till OR disabled_till has passed
        """
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


# credentials/validator.py

REQUIRED_SCOPES = [
    'instagram_basic',
    'instagram_manage_insights',
    'pages_read_engagement',
    'pages_show_list'
]

class CredentialValidator:
    """Validates API credentials before use"""

    async def validate(self, cred: Credential) -> ValidationResult:
        if cred.source == 'graphapi':
            return await self._validate_graphapi(cred)
        elif cred.source == 'ytapi':
            return await self._validate_youtube(cred)
        # ... other sources

    async def _validate_graphapi(self, cred: Credential) -> ValidationResult:
        """
        Validate Facebook Graph API token

        Checks:
        1. Token is valid (not expired)
        2. Has required scopes
        3. Data access hasn't expired
        """
        token = cred.credentials.get('token')
        url = f"https://graph.facebook.com/debug_token?input_token={token}"
        response = await self._request(url)

        if not response['data']['is_valid']:
            raise TokenInvalidError()

        scopes = response['data']['scopes']
        missing = [s for s in REQUIRED_SCOPES if s not in scopes]
        if missing:
            raise MissingScopesError(missing)

        if response['data'].get('data_access_expires_at', 0) < time.time():
            raise DataAccessExpiredError()

        return ValidationResult(valid=True, scopes=scopes)
```

### Schema

```sql
CREATE TABLE credential (
    id BIGSERIAL PRIMARY KEY,
    idempotency_key VARCHAR(255) UNIQUE,
    source VARCHAR(50),              -- graphapi, ytapi, rapidapi-*
    credentials JSONB,               -- {token, user_id, api_key, ...}
    handle VARCHAR(100),             -- Associated Instagram handle
    enabled BOOLEAN DEFAULT TRUE,
    data_access_expired BOOLEAN DEFAULT FALSE,
    disabled_till TIMESTAMP,         -- TTL for temporary disable
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP
);
```

### Interview STAR Story

**Situation**: We had 20+ API credentials across 6 sources. When one hit rate limit, we needed to automatically use another. Manual rotation was error-prone.

**Task**: Build a credential management system with automatic rotation and TTL-based backoff.

**Action**:
1. Designed schema with enable/disable flags and TTL
2. Implemented random selection for load balancing across credentials
3. Added automatic TTL-based re-enable (credential auto-recovers after 1 hour)
4. Built validator to check token validity before use
5. Integrated with AMQP listener for real-time credential updates from Identity service

**Result**:
- **20+ credentials** managed automatically
- **Zero manual intervention** for rate limit handling
- **Automatic recovery** after TTL expires
- **Load balanced** across credentials

---

## COMPONENT 6: 3-Stage Data Pipeline

### What You Built
A clean 3-stage pipeline separating data retrieval, parsing, and processing.

### Technical Deep Dive

```
STAGE 1: RETRIEVAL (instagram/tasks/retrieval.py)
─────────────────────────────────────────────────
Input: handle, source
Output: Raw API response (dict)

async def retrieve_profile_data(handle: str, source: str) -> dict:
    crawler = get_crawler(source)
    cred = await credential_manager.get_enabled_cred(source)
    crawler.set_credentials(cred)
    return await crawler.fetch_profile_by_handle(handle)


STAGE 2: PARSING (instagram/tasks/ingestion.py)
─────────────────────────────────────────────────
Input: Raw API response
Output: Normalized ProfileLog with dimensions/metrics

def parse_profile_data(raw_data: dict, source: str) -> InstagramProfileLog:
    """
    Normalize different API responses to common schema

    Dimensions: Static attributes (handle, name, bio, category)
    Metrics: Numeric values (followers, following, posts)
    """
    parser = get_parser(source)  # Source-specific parser

    return InstagramProfileLog(
        dimensions=[
            Dimension('handle', parser.get_handle(raw_data)),
            Dimension('full_name', parser.get_name(raw_data)),
            Dimension('biography', parser.get_bio(raw_data)),
            Dimension('category', parser.get_category(raw_data)),
            Dimension('is_verified', parser.get_verified(raw_data)),
        ],
        metrics=[
            Metric('followers', parser.get_followers(raw_data)),
            Metric('following', parser.get_following(raw_data)),
            Metric('media_count', parser.get_media_count(raw_data)),
        ]
    )


STAGE 3: PROCESSING (instagram/tasks/processing.py)
─────────────────────────────────────────────────
Input: Normalized ProfileLog
Output: Database upsert + Event publish

async def upsert_profile(profile_log: InstagramProfileLog):
    """
    1. Convert to ORM model
    2. Upsert to PostgreSQL
    3. Create audit log entry
    4. Publish event to AMQP
    """
    # 1. Convert to ORM
    account = InstagramAccount(
        profile_id=profile_log.get_dimension('profile_id'),
        handle=profile_log.get_dimension('handle'),
        followers=profile_log.get_metric('followers'),
        # ... other fields
    )

    # 2. Upsert (insert or update on conflict)
    await session.execute(
        insert(InstagramAccount)
        .values(account.to_dict())
        .on_conflict_do_update(
            index_elements=['profile_id'],
            set_=account.to_dict()
        )
    )

    # 3. Create audit log (for time-series analytics)
    await session.execute(
        insert(ProfileLog).values(
            profile_id=account.profile_id,
            dimensions=profile_log.dimensions_json,
            metrics=profile_log.metrics_json,
            source=profile_log.source
        )
    )

    # 4. Publish event for downstream consumers
    await amqp.publish(
        exchange='beat.dx',
        routing_key='profile_log_events',
        body=profile_log.to_json()
    )
```

### Why This Design?

| Stage | Responsibility | Benefit |
|-------|---------------|---------|
| **Retrieval** | API communication | Easy to add new sources |
| **Parsing** | Response normalization | Isolates API quirks |
| **Processing** | Business logic | Clean, testable |

### Interview STAR Story

**Situation**: Different APIs return data in different formats. GraphAPI uses `followers_count`, RapidAPI uses `follower_count`, some return strings, some integers.

**Task**: Build a pipeline that handles all API variations and produces consistent output.

**Action**:
1. Separated concerns into 3 stages
2. Created source-specific parsers that normalize to common schema
3. Used Dimension/Metric pattern for flexibility
4. Added audit logging for time-series analysis
5. Published events for downstream consumers

**Result**:
- **New API integration** only requires new parser (1 file)
- **Consistent data format** regardless of source
- **Full audit trail** for debugging
- **Event-driven downstream** processing

---

## COMPONENT 7: GPT/OpenAI Integration

### What You Built
AI-powered data enrichment using OpenAI GPT for inferring demographics, categories, and topics from profile bios.

### Technical Deep Dive

```python
# gpt/functions/retriever/openai/openai_extractor.py

class OpenAi(GptCrawlerInterface):
    """Azure OpenAI integration for profile enrichment"""

    def __init__(self):
        openai.api_type = "azure"
        openai.api_base = os.environ["OPENAI_API_BASE"]
        openai.api_version = "2023-05-15"
        openai.api_key = os.environ["OPENAI_API_KEY"]

    async def fetch_instagram_gpt_data_base_gender(
        self, handle: str, bio: str
    ) -> dict:
        """Infer creator's audience gender from bio"""
        prompt = self._load_prompt("profile_info_v0.12.yaml")

        response = await openai.ChatCompletion.acreate(
            engine="gpt-35-turbo",
            messages=[
                {"role": "system", "content": prompt['system']},
                {"role": "user", "content": prompt['user'].format(
                    handle=handle, bio=bio
                )}
            ],
            temperature=0,  # Deterministic
            max_tokens=200
        )

        return self._parse_json_response(response)

    async def fetch_instagram_gpt_data_audience_age_gender(
        self, handle: str, bio: str, recent_posts: list
    ) -> dict:
        """Infer audience demographics from bio + recent posts"""
        content = f"Bio: {bio}\n\nRecent posts:\n"
        content += "\n".join([p['caption'][:200] for p in recent_posts[:5]])

        # Similar implementation with different prompt
```

### Prompt Engineering (13 versions!)

```yaml
# gpt/prompts/profile_info_v0.12.yaml

system: |
  You are an AI assistant that analyzes Instagram creator profiles.
  Based on the username and bio, infer:
  1. Primary audience gender (male/female/mixed)
  2. Confidence score (0.0 to 1.0)

  Consider:
  - Gendered words in bio
  - Content category implications
  - Handle patterns

  Respond ONLY in JSON format.

user: |
  Username: {handle}
  Bio: {bio}

  Output:
  {
    "gender": "male|female|mixed",
    "confidence": 0.85,
    "reasoning": "Brief explanation"
  }

model: gpt-35-turbo
temperature: 0
max_tokens: 200
```

### Use Cases Built

| Flow | Input | Output | Use Case |
|------|-------|--------|----------|
| `base_gender` | handle, bio | gender, confidence | Audience targeting |
| `base_location` | handle, bio | country, city | Geo targeting |
| `categ_lang_topics` | handle, bio, posts | category, language, topics[] | Content classification |
| `audience_age_gender` | handle, bio, posts | age_range, gender_dist | Demographics |
| `audience_cities` | handle, bio, posts | cities with % | Geographic reach |

### Interview STAR Story

**Situation**: Brands wanted to know creator demographics (audience gender, age, location) but Instagram API doesn't provide this for non-business accounts.

**Task**: Build an AI-powered system to infer demographics from publicly available data.

**Action**:
1. Integrated Azure OpenAI with async support
2. Designed prompts through 13 iterations (v0.01 to v0.12)
3. Added temperature=0 for consistent outputs
4. Built JSON parsing with error handling
5. Created separate flows for different enrichment types

**Result**:
- **5 enrichment types** available
- **~85% accuracy** on gender inference (validated against known accounts)
- **13 prompt versions** - continuous improvement
- **Async processing** - doesn't block main flow

### Questions They Might Ask

**Q: How did you validate accuracy?**
A: Compared against 500 creator accounts where we knew actual demographics. Achieved 85% accuracy on gender, 70% on location.

**Q: How do you handle GPT rate limits?**
A: Separate worker pool with low concurrency (2 workers × 5 concurrency). Also implemented exponential backoff on 429.

**Q: What about hallucinations?**
A: Used temperature=0 for deterministic outputs. Also validated JSON schema and rejected malformed responses.

---

## COMPONENT 8: Engagement Calculations

### What You Built
Analytical formulas for calculating engagement metrics from raw data.

### Technical Deep Dive

```python
# instagram/helper.py

def calculate_engagement_rate(likes: int, comments: int,
                              followers: int) -> float:
    """
    Standard engagement rate formula

    Formula: (likes + comments) / followers × 100

    Industry benchmarks:
    - 1-3%: Low engagement
    - 3-6%: Good engagement
    - 6%+: Excellent engagement
    """
    if followers == 0:
        return 0.0
    return ((likes + comments) / followers) * 100


def estimate_reach_reels(plays: int, followers: int) -> float:
    """
    Estimate reach for Reels based on plays

    Empirical formula derived from 10K+ data points:
    factor = 0.94 - (log2(followers) × 0.001)

    Larger accounts have lower reach/follower ratio
    """
    import math
    factor = 0.94 - (math.log2(followers) * 0.001)
    return plays * factor


def estimate_reach_posts(likes: int) -> float:
    """
    Estimate reach for static posts based on likes

    Empirical formula:
    factor = (7.6 - (log10(likes) × 0.7)) × 0.85

    Based on typical like-to-reach ratio
    """
    if likes == 0:
        return 0.0
    import math
    factor = (7.6 - (math.log10(likes) * 0.7)) * 0.85
    return factor * likes


def calculate_avg_metrics(posts: list, exclude_outliers: bool = True) -> dict:
    """
    Calculate average metrics with outlier removal

    Why exclude outliers?
    - Viral posts skew averages
    - Remove top 2 and bottom 2 posts
    - Gives more realistic "typical" performance
    """
    if exclude_outliers and len(posts) > 4:
        # Sort by engagement and remove extremes
        posts = sorted(posts, key=lambda p: p['likes'] + p['comments'])
        posts = posts[2:-2]  # Remove top 2, bottom 2

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

### Interview Point

**Story**: We needed to estimate reach for posts without Instagram Insights access. Derived empirical formulas from 10K+ data points where we had both likes and actual reach.

---

# PROJECT 2: STIR (Data Platform)

## STIR has 8 Major Components:

| # | Component | Files | Complexity |
|---|-----------|-------|------------|
| 1 | Airflow DAG Architecture | 76 DAG files | High |
| 2 | dbt Transformation Layer | 112 models | High |
| 3 | Three-Layer Data Flow | ClickHouse→S3→PostgreSQL | High |
| 4 | Incremental Processing | dbt incremental models | Medium |
| 5 | Multi-Dimensional Rankings | mart_leaderboard | Medium |
| 6 | Time-Series Processing | mart_time_series | Medium |
| 7 | Collection Analytics | mart_collection_* | Medium |
| 8 | Cross-Database Sync | PostgresOperator + SSHOperator | Medium |

---

## COMPONENT 1: Airflow DAG Architecture

### What You Built

```python
# 76 DAGs organized by function

DAG_CATEGORIES = {
    'dbt_orchestration': 11,    # dbt model execution
    'instagram_sync': 17,       # Instagram data triggers
    'youtube_sync': 12,         # YouTube data triggers
    'collection_sync': 15,      # Collection analytics
    'operational': 9,           # Data quality, verification
    'asset_upload': 7,          # Media processing
    'utility': 5                # One-off, helpers
}

# Scheduling Strategy
SCHEDULES = {
    '*/5 * * * *': ['dbt_recent_scl', 'post_ranker_partial'],      # Real-time
    '*/15 * * * *': ['dbt_core'],                                   # Core metrics
    '*/30 * * * *': ['dbt_collections', 'dbt_staging_collections'], # Collections
    '0 * * * *': 12,  # Hourly syncs
    '0 */3 * * *': 8,  # Every 3 hours
    '0 19 * * *': ['dbt_daily'],  # Daily full refresh
    '0 6 */7 * *': ['dbt_weekly']  # Weekly
}
```

### DAG Pattern

```python
from airflow import DAG
from airflow.operators.python import PythonOperator
from airflow_dbt_python.operators.dbt import DbtRunOperator

dag = DAG(
    dag_id='dbt_core',
    default_args={
        'owner': 'airflow',
        'retries': 1,
        'retry_delay': timedelta(minutes=5),
        'on_failure_callback': slack_alert
    },
    schedule_interval='*/15 * * * *',
    start_date=datetime(2023, 1, 1),
    catchup=False,              # Don't backfill
    max_active_runs=1,          # Prevent overlap
    concurrency=1,
    dagrun_timeout=timedelta(minutes=60)
)

dbt_run = DbtRunOperator(
    task_id='dbt_run_core',
    models='tag:core',          # Only models tagged 'core'
    profiles_dir='/Users/.dbt',
    target='gcc_warehouse',     # ClickHouse
    dag=dag
)
```

---

## COMPONENT 2: dbt Transformation Layer

### Model Organization

```
models/
├── staging/                    # 29 models - Raw data cleanup
│   ├── beat/                   # From beat service
│   │   ├── stg_beat_instagram_account.sql
│   │   ├── stg_beat_instagram_post.sql
│   │   └── stg_beat_profile_log.sql
│   └── coffee/                 # From coffee service
│       ├── stg_coffee_campaign_profiles.sql
│       └── stg_coffee_collection.sql
│
└── marts/                      # 83 models - Business logic
    ├── discovery/              # Influencer discovery
    │   ├── mart_instagram_account.sql
    │   └── mart_youtube_account.sql
    ├── leaderboard/            # Rankings
    │   ├── mart_leaderboard.sql
    │   └── mart_time_series.sql
    └── collection/             # Campaign analytics
        └── mart_collection_post.sql
```

### Staging Model Example

```sql
-- models/staging/beat/stg_beat_instagram_account.sql

{{ config(
    materialized='table',
    tags=['staging']
) }}

SELECT
    profile_id,
    handle,
    full_name,
    biography,

    -- Type casting
    CAST(followers_count AS Int64) as followers,
    CAST(following_count AS Int64) as following,
    CAST(posts_count AS Int64) as media_count,

    -- Boolean conversion
    is_verified = 1 as is_verified,
    is_business_account = 1 as is_business,

    -- NULL handling
    COALESCE(category, 'Unknown') as category,

    -- Timestamps
    created_at,
    updated_at

FROM {{ source('beat_replica', 'instagram_account') }}
WHERE handle IS NOT NULL
  AND handle != ''
```

### Mart Model Example

```sql
-- models/marts/discovery/mart_instagram_account.sql

{{ config(
    materialized='table',
    tags=['core', 'hourly']
) }}

WITH base AS (
    SELECT * FROM {{ ref('stg_beat_instagram_account') }}
),

post_stats AS (
    SELECT
        profile_id,
        COUNT(*) as total_posts,
        AVG(likes_count) as avg_likes,
        AVG(comments_count) as avg_comments,
        SUM(likes_count) as total_likes
    FROM {{ ref('stg_beat_instagram_post') }}
    WHERE publish_time > now() - INTERVAL 30 DAY
    GROUP BY profile_id
),

engagement AS (
    SELECT
        b.profile_id,
        (ps.avg_likes + ps.avg_comments) / NULLIF(b.followers, 0) * 100
            as engagement_rate
    FROM base b
    LEFT JOIN post_stats ps USING (profile_id)
)

SELECT
    b.*,
    ps.total_posts,
    ps.avg_likes,
    ps.avg_comments,
    e.engagement_rate,

    -- Rankings
    row_number() OVER (ORDER BY b.followers DESC) as followers_rank,
    row_number() OVER (PARTITION BY b.category
                       ORDER BY b.followers DESC) as followers_rank_by_category,
    row_number() OVER (PARTITION BY b.language
                       ORDER BY b.followers DESC) as followers_rank_by_language

FROM base b
LEFT JOIN post_stats ps USING (profile_id)
LEFT JOIN engagement e USING (profile_id)
```

---

## COMPONENT 3: Three-Layer Data Flow

### Architecture

```
LAYER 1: ClickHouse (OLAP)
         ↓ INSERT INTO FUNCTION s3(...)
LAYER 2: S3 (Staging)
         ↓ aws s3 cp + COPY
LAYER 3: PostgreSQL (Application)
```

### Implementation

```python
# DAG: sync_leaderboard_prod.py

# Task 1: Export from ClickHouse to S3
export_task = ClickHouseOperator(
    task_id='export_to_s3',
    clickhouse_conn_id='clickhouse_gcc',
    sql="""
        INSERT INTO FUNCTION s3(
            's3://gcc-social-data/data-pipeline/tmp/leaderboard.json',
            'AWS_KEY', 'AWS_SECRET',
            'JSONEachRow'
        )
        SELECT
            profile_id,
            handle,
            followers_rank,
            engagement_rank,
            category,
            language
        FROM dbt.mart_leaderboard
        SETTINGS s3_truncate_on_insert=1
    """
)

# Task 2: Download via SSH
download_task = SSHOperator(
    task_id='download_from_s3',
    ssh_conn_id='ssh_prod_pg',
    command='aws s3 cp s3://gcc-social-data/data-pipeline/tmp/leaderboard.json /tmp/'
)

# Task 3: Load into PostgreSQL with atomic swap
load_task = PostgresOperator(
    task_id='load_to_postgres',
    postgres_conn_id='prod_pg',
    sql="""
        -- Create temp table
        CREATE TEMP TABLE tmp_lb (data JSONB);

        -- Load JSON
        COPY tmp_lb FROM '/tmp/leaderboard.json';

        -- Insert with type casting
        INSERT INTO leaderboard_new
        SELECT
            (data->>'profile_id')::bigint,
            (data->>'handle')::text,
            (data->>'followers_rank')::int,
            (data->>'engagement_rank')::int
        FROM tmp_lb;

        -- Atomic swap
        ALTER TABLE leaderboard RENAME TO leaderboard_old;
        ALTER TABLE leaderboard_new RENAME TO leaderboard;
        DROP TABLE IF EXISTS leaderboard_old;
    """
)

export_task >> download_task >> load_task
```

### Why This Pattern?

| Benefit | Explanation |
|---------|-------------|
| **Zero downtime** | Atomic rename, not delete+insert |
| **Decoupled systems** | S3 as intermediate, no direct connection |
| **Debuggable** | Can inspect S3 files |
| **Recoverable** | If load fails, re-download from S3 |

---

## COMPONENT 4: Incremental Processing

### ClickHouse Incremental Model

```sql
-- models/marts/leaderboard/mart_time_series.sql

{{ config(
    materialized='incremental',
    engine='ReplacingMergeTree(updated_at)',
    order_by='(profile_id, date)',
    unique_key='(profile_id, date)'
) }}

SELECT
    profile_id,
    toDate(created_at) as date,
    argMax(followers, created_at) as followers,
    argMax(following, created_at) as following,
    max(created_at) as updated_at

FROM {{ ref('stg_beat_instagram_account') }}

{% if is_incremental() %}
-- Only process new data (with 4-hour safety buffer)
WHERE created_at > (
    SELECT max(updated_at) - INTERVAL 4 HOUR
    FROM {{ this }}
)
{% endif %}

GROUP BY profile_id, date
```

### Why 4-Hour Buffer?

| Reason | Explanation |
|--------|-------------|
| **Late-arriving data** | Some events arrive delayed |
| **Failed retries** | Tasks retried may have old timestamps |
| **Clock drift** | Different servers may have slight time differences |

---

# PROJECT 3: EVENT-GRPC (ClickHouse Sinker)

## Your Specific Work: Consumer → ClickHouse Pipeline

### COMPONENT 1: Buffered Sinker Pattern

```go
// sinker/trace_log_sinker.go

func TraceLogEventsSinker(c chan interface{}) {
    ticker := time.NewTicker(5 * time.Second)  // Flush every 5 sec
    batch := []model.TraceLogEvent{}

    for {
        select {
        case event := <-c:
            // Parse and add to batch
            traceLog := parseTraceLog(event)
            batch = append(batch, traceLog)

            // Flush if batch is full
            if len(batch) >= 1000 {
                flushBatch(batch)
                batch = []model.TraceLogEvent{}
            }

        case <-ticker.C:
            // Periodic flush (even if batch not full)
            if len(batch) > 0 {
                flushBatch(batch)
                batch = []model.TraceLogEvent{}
            }
        }
    }
}

func flushBatch(batch []model.TraceLogEvent) {
    db := clickhouse.Clickhouse(config.New(), nil)
    result := db.Create(&batch)
    if result.Error != nil {
        log.Error("Batch insert failed", result.Error)
        // Could implement retry or dead-letter queue here
    }
}
```

### Why Buffered Sinker?

| Without Buffering | With Buffering |
|-------------------|----------------|
| 1 INSERT per event | 1 INSERT per 1000 events |
| 10,000 DB calls/sec | 10 DB calls/sec |
| High latency | Lower latency |
| DB connection exhaustion | Stable connections |

### COMPONENT 2: Consumer Configuration

```go
// main.go - Your configuration for 26 consumers

// High-volume buffered consumer
traceLogChan := make(chan interface{}, 10000)  // 10K buffer

traceLogConfig := rabbit.RabbitConsumerConfig{
    QueueName:            "trace_log",
    Exchange:             "identity.dx",
    RoutingKey:           "trace_log",
    RetryOnError:         true,
    ErrorExchange:        &errorExchange,
    ErrorRoutingKey:      &errorRoutingKey,
    ConsumerCount:        2,
    BufferChan:           traceLogChan,
    BufferedConsumerFunc: sinker.BufferTraceLogEvent,
}

rabbit.Rabbit(config).InitConsumer(traceLogConfig)
go sinker.TraceLogEventsSinker(traceLogChan)  // Background batch processor
```

### COMPONENT 3: Connection Auto-Recovery

```go
// clickhouse/clickhouse.go

func clickhouseConnectionCron(config config.Config) {
    ticker := time.NewTicker(1 * time.Second)

    for range ticker.C {
        for dbName, db := range singletonClickhouseMap {
            if db == nil {
                reconnect(dbName)
                continue
            }

            // Health check
            sqlDB, _ := db.DB()
            if err := sqlDB.Ping(); err != nil {
                log.Warn("ClickHouse connection lost, reconnecting...")
                reconnect(dbName)
            }
        }
    }
}
```

---

# PROJECT 4: FAKE_FOLLOWER_ANALYSIS

## 7 Components You Built:

| # | Component | Code |
|---|-----------|------|
| 1 | Symbol Normalization | 13 Unicode variant handling |
| 2 | Language Detection | Character-to-language mapping |
| 3 | Indic Transliteration | HMM models for 10 languages |
| 4 | Fuzzy Matching | RapidFuzz weighted ensemble |
| 5 | Indian Name Database | 35,183 names matching |
| 6 | Ensemble Scoring | 5-feature classification |
| 7 | AWS Serverless Pipeline | SQS → Lambda → Kinesis |

### COMPONENT 1: Symbol Normalization

```python
def symbol_name_convert(name):
    """
    Convert fancy Unicode to ASCII

    Handles 13 Unicode variant sets:
    1. 🅐🅑🅒 → ABC (Circled)
    2. 𝐀𝐁𝐂 → ABC (Mathematical Bold)
    3. 𝐴𝐵𝐶 → ABC (Mathematical Italic)
    ... 10 more sets
    """
    # Character mapping tables for each variant
    mappings = {
        'circled': {...},
        'math_bold': {...},
        # ...
    }

    for variant, mapping in mappings.items():
        for fancy, normal in mapping.items():
            name = name.replace(fancy, normal)

    return name
```

### COMPONENT 4: Fuzzy Matching Algorithm

```python
def generate_similarity_score(handle, name):
    """
    Weighted fuzzy matching using RapidFuzz

    Formula: (2×partial + sort + set) / 4

    Why weighted?
    - partial_ratio: Best for substring matching
    - token_sort_ratio: Handles word reordering
    - token_set_ratio: Handles extra/missing words
    """
    from itertools import permutations
    from rapidfuzz import fuzz

    # Generate name permutations
    name_parts = name.split()
    if len(name_parts) <= 4:
        perms = [' '.join(p) for p in permutations(name_parts)]
    else:
        perms = [name]  # Too many permutations

    best_score = 0
    for perm in perms:
        partial = fuzz.partial_ratio(handle, perm)
        sort = fuzz.token_sort_ratio(handle, perm)
        set_ratio = fuzz.token_set_ratio(handle, perm)

        score = (2 * partial + sort + set_ratio) / 4
        best_score = max(best_score, score)

    return best_score
```

### COMPONENT 6: Ensemble Scoring

```python
def final_score(lang_flag, similarity, digit_count, special_char_flag):
    """
    5-feature ensemble → 3 confidence levels

    Returns:
    - 1.0: Definitely FAKE
    - 0.33: Weak FAKE signal
    - 0.0: Likely REAL
    """
    # Strong FAKE indicators
    if lang_flag:  # Non-Indic script
        return 1.0

    if digit_count > 4:  # Too many numbers
        return 1.0

    if special_char_flag == 1:  # Has _ but name doesn't match
        return 1.0

    # Weak FAKE indicator
    if 0 < similarity <= 40:
        return 0.33

    # Default: REAL
    return 0.0
```

---

# SUMMARY: Components Per Project

| Project | Total Components | Key Technical Skills |
|---------|-----------------|---------------------|
| **beat** | 12 | Worker pools, Rate limiting, API integration, Async Python |
| **stir** | 8 | Airflow, dbt, ClickHouse, Data modeling |
| **event-grpc** | 3 | Go channels, Batch processing, Message queues |
| **fake_follower** | 7 | NLP, ML ensemble, AWS Lambda, Transliteration |

---

# HOW TO USE IN INTERVIEW

**When they ask about one project, go deeper:**

Interviewer: "Tell me about beat"

You: "Beat had 12 major components. Which would you like me to dive into?
1. Worker pool architecture
2. Rate limiting system
3. API integration framework
4. Credential management
5. GPT integration
..."

**This shows:**
- You understand the system holistically
- You can go deep on any component
- You have ownership and expertise

---

*Remember: Pick 2-3 components per project you're most confident about and prepare deep-dive stories for those.*

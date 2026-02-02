# GOOGLE INTERVIEW PREPARATION
## Googleyness & Hiring Manager Rounds

---

## YOUR PROJECT OWNERSHIP SUMMARY

| Project | Your Role | Impact |
|---------|-----------|--------|
| **beat** | Core Developer | Built 73 data flows, 150+ workers, 15+ API integrations |
| **stir** | Core Developer | Built 76 Airflow DAGs, 112 dbt models, ClickHouse pipelines |
| **event-grpc** (ClickHouse sinker) | Implemented | Consumer→ClickHouse flush, buffered batch processing |
| **fake_follower_analysis** | Solo Developer | End-to-end ML system, 10 Indic languages, AWS Lambda |

---

# PART 1: GOOGLEYNESS STORIES (STAR Format)

Google evaluates: **Collaboration, Navigating Ambiguity, Pushing Back Respectfully, Helping Others, Learning from Failures, Bias to Action**

---

## STORY 1: Building a Highly Reliable Data Scraping Platform (beat)

### Situation
At Good Creator Co., we needed a system to aggregate social media data (Instagram, YouTube) at scale for our influencer analytics platform. The challenge was handling 15+ external APIs with different rate limits, frequent failures, and the need to process 10M+ data points daily.

### Task
I was responsible for designing and building the entire data collection service from scratch - including the worker pool system, rate limiting, API integrations, and data processing flows.

### Action
1. **Designed worker pool architecture**: Created 73 configurable flows with multiprocessing workers + async concurrency
   ```python
   # Each flow had configurable workers and concurrency
   'refresh_profile_by_handle': {'no_of_workers': 10, 'no_of_concurrency': 5}
   ```

2. **Built multi-level rate limiting**: Implemented Redis-backed stacked limiters (global daily 20K, per-minute 60, per-handle 1/sec)
   ```python
   async with RateLimiter(rate_spec=global_limit_day):
       async with RateLimiter(rate_spec=global_limit_minute):
           async with RateLimiter(rate_spec=handle_limit):
               result = await refresh_profile(handle)
   ```

3. **Designed fallback strategy**: When primary APIs failed, system automatically rotated to secondary sources (6 Instagram APIs, 4 YouTube APIs)

4. **Implemented credential rotation**: Built manager to disable credentials with TTL backoff when rate-limited

### Result
- **10M+ daily data points** processed reliably
- **25% improvement** in API response times through optimization
- **30% cost reduction** through intelligent rate limiting and caching
- **150+ concurrent workers** running smoothly
- System ran in production for 15+ months with minimal incidents

### Googleyness Signals
- **Bias to Action**: Built from scratch rather than waiting for perfect spec
- **Navigating Ambiguity**: External APIs changed frequently, designed for adaptability
- **Helping Others**: The platform enabled entire analytics team to deliver insights

---

## STORY 2: Designing the Data Platform (stir)

### Situation
Our analytics platform needed to compute influencer rankings, engagement metrics, and time-series data across billions of records. The existing manual SQL queries were slow, error-prone, and couldn't scale.

### Task
Build an enterprise data platform that could:
- Transform raw data into analytics-ready models
- Sync data between ClickHouse (analytics) and PostgreSQL (application)
- Run reliably with proper monitoring and error handling

### Action
1. **Chose Modern Data Stack**: Selected Airflow + dbt + ClickHouse after evaluating alternatives
   - Airflow for orchestration (vs Prefect) - better community, more operators
   - dbt for transformation (vs stored procedures) - version control, testing
   - ClickHouse (vs BigQuery) - self-hosted, no egress costs

2. **Designed 3-layer data flow**:
   ```
   ClickHouse (analytics) → S3 (staging) → PostgreSQL (application)
   ```
   This atomic swap pattern ensured zero-downtime updates

3. **Built 76 DAGs with proper scheduling**:
   - `*/5 min`: Real-time (dbt_recent_scl)
   - `*/15 min`: Core metrics (dbt_core)
   - `Daily 19:00`: Full refresh (dbt_daily)

4. **Implemented incremental processing**:
   ```sql
   {% if is_incremental() %}
   WHERE created_at > (SELECT max(created_at) - INTERVAL 4 HOUR FROM {{ this }})
   {% endif %}
   ```

### Result
- **50% reduction** in data latency
- **76 production DAGs** running reliably
- **112 dbt models** powering all analytics
- **Billions of records** processed daily
- Enabled multi-dimensional leaderboards (category, language, country rankings)

### Googleyness Signals
- **Collaboration**: Worked with data analysts to understand requirements
- **Learning**: Learned dbt and ClickHouse specifically for this project
- **Pushing Back**: Advocated for dbt over raw SQL despite initial resistance

---

## STORY 3: Real-Time Event Processing to ClickHouse (event-grpc)

### Situation
Our mobile and web apps generated thousands of events per second (user actions, clicks, purchases). These needed to be reliably stored in ClickHouse for real-time analytics.

### Task
I specifically owned the **consumer → ClickHouse** pipeline - the part that consumes messages from RabbitMQ and flushes them to ClickHouse reliably.

### Action
1. **Designed buffered sinker pattern** for high-volume events:
   ```go
   func TraceLogEventsSinker(c chan interface{}) {
       ticker := time.NewTicker(5 * time.Second)
       batch := []model.TraceLogEvent{}

       for {
           select {
           case event := <-c:
               batch = append(batch, parseEvent(event))
               if len(batch) >= 1000 {
                   flushBatch(batch)
                   batch = []model.TraceLogEvent{}
               }
           case <-ticker.C:
               if len(batch) > 0 {
                   flushBatch(batch)
                   batch = []model.TraceLogEvent{}
               }
           }
       }
   }
   ```

2. **Implemented retry logic with dead letter queue**:
   - Max 2 retries before routing to error queue
   - Preserved message metadata across retries

3. **Built connection auto-recovery**:
   ```go
   // 1-second cron to check and reconnect ClickHouse
   func clickhouseConnectionCron(config config.Config) {
       ticker := time.NewTicker(1 * time.Second)
       for range ticker.C {
           for dbName, db := range singletonClickhouseMap {
               if db == nil || db.Error != nil {
                   reconnect(dbName)
               }
           }
       }
   }
   ```

4. **Created 26 different consumer configurations** for various event types

### Result
- **10,000+ events/second** processed reliably
- **70+ concurrent workers** handling different event types
- **Zero data loss** through buffered writes + retry logic
- **18+ ClickHouse tables** populated with real-time data

### Googleyness Signals
- **Technical Excellence**: Designed for reliability and scale
- **Ownership**: Took full responsibility for critical data path
- **Bias to Action**: Proactively added monitoring and alerting

---

## STORY 4: ML-Powered Fake Follower Detection (fake_follower_analysis)

### Situation
Influencer marketing campaigns were being affected by fake followers. Brands needed to verify that creators had genuine audiences. The challenge: detecting fakes among followers who could use any of 10+ Indian languages.

### Task
Build an end-to-end ML system that could:
- Detect fake followers with high accuracy
- Handle 10 Indic scripts (Hindi, Bengali, Tamil, etc.)
- Scale to millions of followers
- Run cost-effectively on serverless infrastructure

### Action
1. **Designed ensemble detection model** with 5 independent features:
   ```python
   # Feature 1: Non-Indic language detection (Greek, Chinese, Korean = FAKE)
   # Feature 2: Digit count in handle (>4 digits = FAKE)
   # Feature 3: Handle-name correlation (special chars but no match = FAKE)
   # Feature 4: Fuzzy similarity score (RapidFuzz weighted)
   # Feature 5: Indian name database match (35,183 names)
   ```

2. **Built multi-language transliteration pipeline**:
   - Integrated indictrans library with HMM models
   - Custom Hindi processing with 24 vowel + 42 consonant mappings
   - Symbol normalization for 13 Unicode variants

3. **Designed AWS serverless architecture**:
   ```
   ClickHouse → S3 → SQS → Lambda (ECR) → Kinesis → Output
   ```

4. **Optimized for cost and performance**:
   - Batch processing with 10,000 messages/batch
   - 8 parallel workers using multiprocessing
   - ON_DEMAND Kinesis for auto-scaling

### Result
- **Processes entire follower base** in minutes
- **10 Indic scripts** supported seamlessly
- **3 confidence levels** (0.0, 0.33, 1.0) for nuanced decisions
- **35,183 name database** for validation
- Enabled brands to make data-driven influencer selections

### Googleyness Signals
- **Innovation**: Combined NLP, ML, and cloud architecture creatively
- **User Focus**: Understood brand needs and delivered actionable scores
- **Technical Depth**: Deep dive into linguistics, Unicode, ML models

---

## STORY 5: Navigating Ambiguity - API Rate Limit Crisis (beat)

### Situation
One day, Instagram Graph API started returning 429 (rate limit) errors at 10x the normal rate. Our data collection stopped. No warning from Facebook, no documentation about changes.

### Task
Quickly diagnose and fix the issue while maintaining data freshness for customers.

### Action
1. **Immediate investigation**: Analyzed patterns - found Facebook had silently reduced rate limits
2. **Quick mitigation**: Reduced concurrent workers from 50 to 20 immediately
3. **Medium-term fix**: Implemented credential rotation across multiple Facebook accounts
4. **Long-term solution**: Built adaptive rate limiting that learns from 429 responses
   ```python
   async def disable_creds(cred_id, disable_duration=3600):
       # Disable credential with TTL backoff
       await session.execute(
           update(Credential)
           .where(Credential.id == cred_id)
           .values(
               enabled=False,
               disabled_till=func.now() + timedelta(seconds=disable_duration)
           )
       )
   ```

### Result
- **Restored data collection** within 2 hours
- **Built resilient system** that handles future rate limit changes automatically
- **Documented incident** and created runbook for team

### Googleyness Signals
- **Bias to Action**: Didn't wait for perfect info, acted immediately
- **Learning**: Turned crisis into opportunity to build better system
- **Helping Others**: Created documentation for future incidents

---

## STORY 6: Pushing Back Respectfully - Technology Choice (stir)

### Situation
When building the data platform, the team wanted to use a commercial ETL tool (Fivetran) for transformations. I believed dbt would be a better choice for our use case.

### Task
Convince stakeholders that open-source dbt was the right choice without creating conflict.

### Action
1. **Built a POC**: Created 5 example models in dbt showing the workflow
2. **Presented trade-offs objectively**:
   - Fivetran: Easier setup, but $500+/month, limited customization
   - dbt: Learning curve, but free, full control, version-controlled
3. **Addressed concerns**:
   - "Learning curve?" → I'll create documentation and train the team
   - "Support?" → Active community, extensive documentation
4. **Offered compromise**: "Let's try dbt for 2 weeks. If it doesn't work, we switch to Fivetran"

### Result
- **dbt adopted** as primary transformation tool
- **$6,000+/year saved** on licensing
- **Team upskilled** in modern data stack
- **112 models** built successfully using dbt

### Googleyness Signals
- **Pushing Back Respectfully**: Data-driven argument, not emotional
- **Collaboration**: Offered to train team, shared responsibility
- **User Focus**: Chose what was best for long-term success

---

## STORY 7: Helping Others - Mentoring Junior Engineer

### Situation
A junior engineer was struggling to debug a complex Airflow DAG failure. The error messages were cryptic, and they had been stuck for 2 days.

### Task
Help them while teaching them how to debug such issues in the future.

### Action
1. **Didn't just fix it**: Sat with them and walked through the debugging process
2. **Taught systematic approach**:
   - Check Airflow logs (scheduler, worker)
   - Identify which task failed
   - Check task's Python code and dependencies
   - Verify database connections and permissions
3. **Found root cause together**: A ClickHouse connection timeout due to missing retry logic
4. **Implemented fix together**: Added connection retry with exponential backoff
5. **Created documentation**: Wrote a "DAG Debugging Guide" for the team

### Result
- **Junior engineer** solved future issues independently
- **Debugging guide** used by entire team
- **Reduced escalations** by 40%

### Googleyness Signals
- **Helping Others**: Invested time in teaching, not just doing
- **Collaboration**: Made it a learning experience
- **Documentation**: Created lasting value for team

---

# PART 2: HIRING MANAGER DEEP DIVE QUESTIONS

---

## Technical Deep Dive: beat Project

### Q: "Walk me through the architecture of beat"
**Answer:**
```
Client (coffee API) → FastAPI REST API → SQL-based task queue
                                              ↓
                         Worker Pool (73 flows × N workers each)
                                              ↓
                         Rate Limiter (Redis-backed, stacked limits)
                                              ↓
                         15+ APIs (Instagram Graph, RapidAPI, YouTube)
                                              ↓
                         3-Stage Pipeline: Retrieval → Parsing → Processing
                                              ↓
                         PostgreSQL (transactional) + RabbitMQ (events)
```

Key design decisions:
1. **SQL-based task queue** instead of Celery - simpler, no additional infrastructure
2. **Stacked rate limiters** - per-handle, per-minute, per-day for fine-grained control
3. **Interface-based API integrations** - easy to add/swap API providers

### Q: "How did you handle 10M+ daily data points?"
**Answer:**
1. **Parallel processing**: 150+ workers with semaphore-based concurrency
2. **Async I/O**: uvloop + aio-pika + asyncpg for non-blocking operations
3. **Batch operations**: Grouped similar tasks, batch database writes
4. **Intelligent caching**: Redis for rate limit state, API responses
5. **Priority queuing**: Important profiles processed first

### Q: "What would you do differently if building it again?"
**Answer:**
1. **Replace SQL task queue with Redis Streams** - better for high-throughput
2. **Add distributed tracing** - easier debugging across services
3. **Use structured logging** - better for alerting and analysis
4. **Add circuit breakers per API** - isolate failures better

---

## Technical Deep Dive: stir Project

### Q: "Why did you choose Airflow + dbt + ClickHouse?"
**Answer:**

| Technology | Why Chosen | Alternatives Considered |
|------------|------------|------------------------|
| **Airflow** | Python-native, huge ecosystem, 50+ operators | Prefect (newer, smaller community) |
| **dbt** | SQL transformations, version control, testing | Stored procedures (no versioning) |
| **ClickHouse** | OLAP performance, columnar storage, free | BigQuery (egress costs), Snowflake (expensive) |

The combination allowed:
- **Airflow**: Scheduling, monitoring, retries, dependencies
- **dbt**: Modular SQL, incremental processing, documentation
- **ClickHouse**: Sub-second queries on billions of rows

### Q: "Explain the 3-layer data flow"
**Answer:**
```
LAYER 1: ClickHouse (Analytics)
   - dbt models run here
   - Fast OLAP queries
   - ReplacingMergeTree for efficient upserts
         ↓
   INSERT INTO FUNCTION s3('bucket/file.json')
         ↓
LAYER 2: S3 (Staging)
   - Intermediate storage
   - Decouples systems
   - Allows retry without re-processing
         ↓
   aws s3 cp + COPY command
         ↓
LAYER 3: PostgreSQL (Application)
   - Powers REST APIs
   - JSONB parsing for flexibility
   - Atomic table swap (RENAME) for zero-downtime
```

Why this pattern?
1. **No direct connection needed** between ClickHouse and PostgreSQL
2. **Atomic updates** - application sees old data until new is ready
3. **Easy debugging** - S3 files can be inspected

### Q: "How did you handle incremental processing?"
**Answer:**
```sql
{{ config(
    materialized='incremental',
    engine='ReplacingMergeTree(updated_at)',
    order_by='(profile_id, date)'
) }}

{% if is_incremental() %}
WHERE created_at > (
    SELECT max(created_at) - INTERVAL 4 HOUR  -- Safety buffer
    FROM {{ this }}
)
{% endif %}
```

The 4-hour lookback handles:
- Late-arriving data
- Failed task retries
- Clock drift between systems

---

## Technical Deep Dive: fake_follower_analysis

### Q: "How does your ML model detect fake followers?"
**Answer:**
5-feature ensemble:

| Feature | Logic | Weight |
|---------|-------|--------|
| **Non-Indic Script** | Greek/Chinese/Korean = FAKE | 1.0 (definite) |
| **Digit Count** | >4 digits in handle = FAKE | 1.0 (definite) |
| **Handle-Name Correlation** | Special chars but no match = FAKE | 1.0 (definite) |
| **Fuzzy Similarity** | 0-40% similarity = weak FAKE | 0.33 (weak) |
| **Indian Name Match** | <80% match to 35K names = suspicious | supplementary |

Decision tree:
```python
if non_indic_language: return 1.0  # FAKE
if digits > 4: return 1.0  # FAKE
if special_chars and similarity < 80: return 1.0  # FAKE
if similarity < 40: return 0.33  # Weak signal
return 0.0  # REAL
```

### Q: "How did you handle 10 different Indian languages?"
**Answer:**
1. **Character-to-language mapping**: Each script has unique Unicode ranges
   - Hindi: 0900-097F
   - Bengali: 0980-09FF
   - Tamil: 0B80-0BFF

2. **HMM-based transliteration**: Pre-trained models for each language pair
   ```python
   trn = Transliterator(source='hin', target='eng', decode='viterbi')
   english_name = trn.transform("राहुल")  # "Rahul"
   ```

3. **Custom Hindi processing**: 24 vowel + 42 consonant manual mappings for accuracy

4. **Fallback chain**: Indic script → ML transliteration → unidecode → ASCII

### Q: "Why AWS Lambda over EC2?"
**Answer:**
| Factor | Lambda | EC2 |
|--------|--------|-----|
| **Cost** | Pay per invocation | Pay always |
| **Scaling** | Automatic | Manual setup |
| **Maintenance** | None | OS, security patches |
| **Cold start** | ~2s (acceptable for batch) | None |

For batch processing of followers (not real-time), Lambda was:
- **90% cheaper** than running EC2 24/7
- **Zero operational overhead**
- **Auto-scales** with SQS queue depth

---

# PART 3: BEHAVIORAL QUESTIONS

## "Tell me about a time you failed"

**Story**: GPT Integration Timeout Issue

**Situation**: Added OpenAI GPT integration to beat for data enrichment. In testing, it worked great.

**What went wrong**: In production, 30% of requests timed out. GPT API had variable latency that I didn't account for.

**What I learned**:
1. Always load-test with realistic conditions
2. Implement timeouts and fallbacks for external services
3. Make features degradable - system should work without optional enrichments

**What I did**:
1. Added 30-second timeout with retry
2. Made GPT enrichment asynchronous (separate worker)
3. System works without GPT data, enriches later

---

## "Describe a conflict with a teammate"

**Story**: Database Choice Disagreement

**Situation**: A senior engineer wanted to use MongoDB for the analytics platform. I believed ClickHouse was better for our OLAP workload.

**How I handled it**:
1. **Listened first**: Understood their reasons (familiar with Mongo, document flexibility)
2. **Proposed experiment**: "Let's benchmark both with our actual queries"
3. **Shared results objectively**: ClickHouse was 50x faster for aggregation queries
4. **Acknowledged trade-offs**: "You're right about flexibility, but performance wins here"

**Outcome**: We went with ClickHouse. Senior engineer became an advocate after seeing performance.

---

## "How do you prioritize tasks?"

**Framework I use**:
1. **Impact vs Effort matrix**: High impact, low effort first
2. **Dependencies**: Unblock others before personal tasks
3. **Deadlines**: Customer-facing deadlines are non-negotiable
4. **Technical debt**: Allocate 20% time to pay down debt

**Example from stir**:
- Had to choose between new leaderboard feature vs. DAG reliability improvements
- Chose reliability first because failing DAGs blocked entire team
- Delivered leaderboard 1 week later, but with stable foundation

---

## "Why Google?"

**Honest answer**:
1. **Scale**: Want to work on systems serving billions of users
2. **Learning**: Google's engineering culture is legendary
3. **Impact**: Build tools used by developers worldwide
4. **Growth**: Learn from the best engineers in the industry

**What I bring**:
1. End-to-end ownership experience (built complete systems)
2. Both Python and Go expertise
3. Data engineering + backend + ML breadth
4. Startup scrappiness (bias to action, do more with less)

---

# PART 4: QUESTIONS TO ASK

## For Hiring Manager
1. "What does success look like in the first 6 months?"
2. "What are the biggest technical challenges the team is facing?"
3. "How do you balance feature work vs. technical debt?"
4. "What's the team's approach to on-call and incident response?"

## For Googleyness
1. "Can you share an example of how the team navigated ambiguity recently?"
2. "How does the team handle disagreements on technical decisions?"
3. "What opportunities are there for cross-team collaboration?"
4. "How does Google support continuous learning?"

---

# PART 5: METRICS CHEAT SHEET

## beat
- **10M+ daily data points** processed
- **73 flows**, **150+ workers**
- **15+ API integrations**
- **25% faster** API response times
- **30% cost reduction**

## stir
- **76 Airflow DAGs**
- **112 dbt models** (29 staging + 83 marts)
- **50% data latency reduction**
- **Billions of records** processed
- **1,476 git commits** (mature project)

## event-grpc (your part)
- **10,000+ events/second**
- **26 consumer queues**
- **70+ concurrent workers**
- **18+ ClickHouse tables**
- **Buffered batch writes** (1000 events/batch)

## fake_follower_analysis
- **10 Indic scripts** supported
- **35,183 name database**
- **5-feature ensemble model**
- **AWS Lambda** serverless
- **3 confidence levels** (0.0, 0.33, 1.0)

---

# PART 6: TECHNICAL TERMS TO KNOW

| Term | What You Built | How to Explain |
|------|----------------|----------------|
| **Worker Pool** | beat main.py | Multiprocessing + async semaphores for concurrency control |
| **Rate Limiting** | beat server.py | Redis-backed stacked limiters (daily, per-min, per-resource) |
| **ELT** | stir | Extract-Load-Transform - load raw data first, transform in warehouse |
| **Incremental Processing** | stir dbt models | Only process new/changed data, not full table |
| **Buffered Sinker** | event-grpc | Batch events in memory, flush periodically for efficiency |
| **HMM Transliteration** | fake_follower | Hidden Markov Model for character sequence conversion |
| **Ensemble Model** | fake_follower | Combine multiple weak classifiers for robust prediction |

---

# PART 7: RESUME ↔ PROJECT MAPPING

Your resume says → Proof from projects:

| Resume Bullet | Project Evidence |
|--------------|------------------|
| "Built a high-performance logging system with RabbitMQ, Python and Golang, transitioning to ClickHouse, achieving a 2.5x reduction in log retrieval times" | **event-grpc**: Consumer→ClickHouse sinker, buffered batch writes |
| "Crafted and streamlined ETL data pipelines (Apache Airflow) for batch data ingestion" | **stir**: 76 DAGs, dbt transformations, ClickHouse→S3→PostgreSQL flow |
| "Designed an asynchronous data processing system handling 10M+ daily data points" | **beat**: 73 flows, 150+ workers, async Python with uvloop |
| "Optimized API response times by 25%" | **beat**: Rate limiting, caching, credential rotation |
| "Reduced operational costs by 30%" | **beat**: Intelligent rate limiting reduced API costs |

---

*Good luck with your Google interview! Remember: Be specific, use numbers, and show ownership.*

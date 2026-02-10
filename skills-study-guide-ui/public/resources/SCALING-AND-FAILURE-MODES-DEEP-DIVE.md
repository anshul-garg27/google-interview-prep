# Scaling Strategies & Failure Modes — Complete Deep Dive

> **What this file covers**: Everything NOT in the other docs. Scaling scenarios (10x/100x/1000x), Kafka-specific failure modes, system-wide issues across ALL projects (Walmart + GCC), "what could go wrong" thinking, and exactly how to answer when an interviewer says "what if..."

---

## Table of Contents

1. [Kafka Audit System — Scaling Roadmap (10x / 100x / 1000x)](#1-kafka-audit-system--scaling-roadmap)
2. [Kafka Deep Problems — Every Issue a Senior Engineer Should Know](#2-kafka-deep-problems)
3. [Cross-System Failure Modes — All Walmart Projects](#3-cross-system-failure-modes--walmart)
4. [Cross-System Failure Modes — All GCC Projects](#4-cross-system-failure-modes--gcc)
5. [Database Scaling Issues (PostgreSQL, ClickHouse, Cosmos, Redis)](#5-database-scaling-issues)
6. [Distributed Systems Failure Patterns](#6-distributed-systems-failure-patterns)
7. [Interview "What If" Questions — Complete Answer Bank](#7-interview-what-if-questions)

---

## 1. Kafka Audit System — Scaling Roadmap

### Current State (Baseline)

| Metric | Value |
|--------|-------|
| Events/day | 2M (~23/sec avg, ~100/sec peak) |
| Partitions | 12 |
| Brokers | 3 (shared cluster) |
| GCS Sink workers | 1 (tasks.max=1) |
| Avg event size | ~2KB (Avro) |
| Daily data volume | ~4GB raw, ~400MB Parquet |
| Latency impact | <5ms P99 |

---

### Scaling to 10x (20M events/day, ~1,000/sec peak)

**What breaks first**: Tier 2 HTTP endpoint (connection exhaustion) and Tier 3 single consumer task.

| Component | Problem | Fix | Effort |
|-----------|---------|-----|--------|
| **Tier 1 ThreadPool** | Queue (100) fills at sustained 1K/sec, audits start dropping | Increase to `core=20, max=40, queue=500`. Monitor queue utilization at 80% | Config change |
| **Tier 1 → Tier 2 HTTP** | 1K HTTP POST/sec per pod. Connection pool exhaustion, TCP overhead | Switch from `RestTemplate` to `WebClient` (non-blocking). Connection pool: `maxConnections=200, pendingAcquireTimeout=5s` | Code change (1 week) |
| **Tier 2 Pods** | Single pod handles ~500 req/sec max (Spring Boot) | HPA: scale to 3-5 pods. `targetCPUUtilizationPercentage: 60` | Config change |
| **Kafka Partitions** | 12 partitions still fine for throughput, but consumer parallelism limited | Increase to **24 partitions**. Rebalance during low-traffic window | Ops (1 day) |
| **Tier 3 Tasks** | `tasks.max=1` is the bottleneck. Single thread processes 50 records per poll | Increase `tasks.max=6` (one task per 4 partitions). Scale Kafka Connect workers to 2 | Config change |
| **GCS** | Handles 10x easily. No changes needed | None | — |
| **Avro/Schema Registry** | Caches schema. Read-heavy, handles 10x trivially | None | — |
| **Monitoring** | Prometheus cardinality grows but manageable | Add aggregated dashboards per-service (not per-request) | Config change |

**Total effort**: ~1-2 weeks. Mostly config changes + one WebClient migration.

**Key insight for interview**: *"At 10x, the architecture holds. We tune thread pools, scale horizontally, increase partition count, and add more Kafka Connect tasks. No fundamental redesign."*

---

### Scaling to 100x (200M events/day, ~10,000/sec peak)

**What breaks**: The HTTP hop between Tier 1 and Tier 2 becomes the bottleneck. Shared Kafka cluster can't handle the load. Storage costs explode.

| Component | Problem | Architectural Fix |
|-----------|---------|-------------------|
| **Tier 1 → Tier 2 HTTP** | 10K HTTP calls/sec = unsustainable. Each call: DNS lookup + TCP handshake + TLS + serialize + deserialize + response wait = ~50ms overhead per event | **ELIMINATE TIER 2.** Embed a lightweight Kafka producer directly in the common library JAR. Async, batched (`linger.ms=50`, `batch.size=64KB`). One network hop: App → Kafka directly |
| **Kafka Cluster** | Shared cluster. 200M messages/day with RF=3 = 600M total writes/day. Other business topics get starved | **Dedicated Kafka cluster** for audit. 6-9 brokers. NVMe SSDs for hot tier. Enable Kafka tiered storage (KIP-405) for cold data on object storage |
| **Partitions** | 24 partitions → max 24 parallel consumers. Not enough | **100+ partitions.** New key: `service_name\|endpoint_name\|hash(request_id) % shard_count` for even distribution |
| **Kafka Connect** | Single cluster, few workers | **Distributed mode.** 5+ workers per region. `tasks.max=20` per connector. Separate clusters for US/CA/MX |
| **Serialization** | Avro is fine. Schema Registry caching handles volume | Switch Kafka compression from `snappy` to **`lz4`** — 2x faster compression/decompression at similar ratio. Matters at high throughput |
| **Storage** | 200M × 2KB = 400GB/day raw, ~40GB Parquet. **14.6TB/year** | **Tiered GCS lifecycle**: Standard (30 days) → Nearline (1 year) → Coldline (7 years). Cost: $430/year Coldline vs $15K/year Standard |
| **BigQuery** | External tables over 14.6TB get slow (full file scans) | **Native BigQuery tables** with daily batch loads. Partition by `_PARTITIONDATE`, cluster by `service_name`. BI Engine for dashboards |
| **Monitoring** | 10K events/sec generates extreme metric cardinality | **Aggregated metrics only** — per-service per-minute buckets. Drop per-request tracing for audit events. Use sampling (1 in 100) for distributed traces |

**New architecture at 100x:**

```
12+ Services (Tier 1 — embedded Kafka producer, no HTTP hop)
       │ Avro Binary (batched, lz4 compressed)
       ▼
Dedicated Kafka Cluster (6-9 brokers, 100+ partitions, tiered storage)
       │
  ┌────┼────┐
  ▼    ▼    ▼
KC-US KC-CA KC-MX  (5 workers each, 20 tasks each)
  │    │    │
  ▼    ▼    ▼
GCS (Lifecycle: Standard→Nearline→Coldline)
       │
       ▼
BigQuery (Native tables, partitioned, clustered, BI Engine)
```

**Total effort**: 4-6 weeks. Requires library JAR v2 (embedded producer), Kafka cluster provisioning, storage migration.

**Key insight for interview**: *"At 100x, the architecture changes fundamentally. I'd eliminate the HTTP middle tier, give Kafka its own dedicated cluster, and move to tiered storage everywhere. The key principle: reduce network hops, increase batching, and use lifecycle policies to manage cost."*

---

### Scaling to 1000x (2B events/day, ~100,000/sec peak)

**This is Walmart-scale territory** (think Black Friday across all systems).

| Component | Approach |
|-----------|----------|
| **Ingestion** | Kafka producer in library with `linger.ms=200`, `batch.size=256KB`, `buffer.memory=128MB`. Accept slight latency for throughput. Client-side batching is critical |
| **Kafka** | Multi-cluster with MirrorMaker 2 for cross-DC replication. 500+ partitions across multiple topics (shard by service). Dedicated broker pools for producers vs consumers |
| **Processing** | Replace Kafka Connect with **Apache Flink** — better at backpressure, exactly-once, windowing. Flink checkpoints for exactly-once GCS writes |
| **Storage** | **Delta Lake** or **Apache Iceberg** on GCS instead of raw Parquet. ACID transactions, schema evolution, time-travel queries. Optimized for petabyte-scale |
| **Query** | **Trino/Presto** federated query across Iceberg tables. Or **BigQuery Omni** for multi-cloud |
| **Cost** | ~146TB/year in Parquet. Coldline: $4,300/year. Must auto-archive aggressively. Consider dropping response bodies after 90 days (keep metadata only) |
| **Compliance** | Data residency: Separate clusters per country. GDPR: Right-to-deletion on audit logs = complex (need to identify PII fields, build deletion pipeline) |

**Key insight for interview**: *"At 1000x, Kafka Connect can't keep up — I'd move to Flink for processing. Storage shifts from raw Parquet to Iceberg for ACID guarantees. And we need to start thinking about data lifecycle — do we really need response bodies after 90 days? Reducing payload size is the cheapest way to scale."*

---

## 2. Kafka Deep Problems — Every Issue a Senior Engineer Should Know

These are the problems interviewers test for. Knowing these = senior-level Kafka understanding.

---

### 2.1 Consumer Rebalancing Storm

**What**: When consumers join/leave a group, Kafka rebalances partition assignments. During rebalance, ALL consumers in the group STOP processing.

**Why it's dangerous**: If rebalances happen frequently (KEDA scaling, rolling deploys, health check failures), you get a "rebalance storm" — consumers spend more time rebalancing than consuming.

**Your real example**: KEDA autoscaling on consumer lag → more consumers → rebalance → no consumption → more lag → more scaling → infinite loop.

**Fixes**:
- `session.timeout.ms=45000` (higher = fewer false rebalances)
- `max.poll.interval.ms=300000` (long enough for slow batches)
- `group.instance.id=static-member-N` (static membership — no rebalance on restart if same ID)
- Use **cooperative rebalancing** (`partition.assignment.strategy=CooperativeStickyAssignor`) — only revokes partitions that need to move, not all
- NEVER scale Kafka consumers on lag. Scale on CPU/memory only

**Interview answer**: *"I hit this exact issue. KEDA was scaling on consumer lag which created a feedback loop. The fix was removing KEDA, switching to CPU-based HPA, and using static group membership so restarts don't trigger full rebalances."*

---

### 2.2 Exactly-Once Semantics (EOS)

**What**: Can you guarantee each audit event is written to GCS exactly once? Not duplicated, not lost?

**Current state**: At-least-once delivery. Duplicates are possible if:
- Kafka Connect processes a batch, writes to GCS, then crashes before committing offset → re-processes same batch → duplicate files in GCS

**Why it matters**: For audit logs, duplicates are tolerable (query with DISTINCT). But if someone asks...

**How to achieve exactly-once**:
1. **Idempotent producer**: `enable.idempotence=true` (already default in Kafka 3.0+). Prevents producer duplicates.
2. **Transactional producer**: `transactional.id=audit-producer-1`. Wraps produce + offset commit in single transaction.
3. **Kafka Connect EOS**: Set `exactly.once.source.support=enabled` (for source connectors). For sink connectors, rely on idempotent writes.
4. **Idempotent GCS writes**: Use deterministic file names (`topic-partition-offset.parquet`). Re-writing same file = idempotent.

**Interview answer**: *"Currently we're at-least-once, which is fine for audit logs — duplicates are filtered at query time. If we needed exactly-once, I'd enable idempotent producers, use transactional semantics, and make GCS writes idempotent with deterministic file naming."*

---

### 2.3 Message Ordering Guarantees

**What**: Are audit logs written in order?

**Current**: Messages are ordered **within a partition** only. Key = `service_name|endpoint_name` ensures all events for the same service+endpoint go to the same partition → ordered for that service.

**Problem at scale**: If you increase partitions from 12 → 100, existing keys get redistributed → temporary ordering break during partition increase.

**Fixes**:
- **Don't rekey during partition increase** — old data stays in old partitions, new data goes to new ones. Ordering within partition is maintained.
- If strict global ordering needed (we don't): single partition (kills parallelism) or use **sequence numbers** in the payload and sort at query time.

**Interview answer**: *"We have per-service ordering via partition key. For audit logs, that's sufficient — we don't need global ordering. If we did, I'd add a sequence number to the Avro schema and sort at query time in BigQuery."*

---

### 2.4 Poison Pill Messages

**What**: A malformed message that crashes the consumer every time it tries to process it. Consumer restarts, re-reads same message, crashes again → infinite loop.

**Your system's handling**:
- `errors.tolerance=all` in Kafka Connect → skips bad messages
- DLQ (`errors.deadletterqueue.topic.name`) captures them for later analysis
- SMT filter has try-catch with fallback (PR #35)

**What could still go wrong**:
- DLQ fills up and nobody monitors it → data loss goes unnoticed
- A schema change in Avro makes ALL messages look like poison pills to the deserializer

**Fix**:
- **Alert on DLQ depth** > 0 (any message in DLQ = something is wrong)
- **Schema Registry compatibility checks**: Set `BACKWARD` compatibility mode → new schema must be able to read old data
- **Schema version pinning**: Consumer explicitly requests schema version, doesn't auto-upgrade

---

### 2.5 Kafka Broker Disk Full

**What**: If brokers run out of disk, they stop accepting writes. All producers get errors.

**Why relevant**: At 2M events/day with RF=3 and 7-day retention, you're storing ~84GB on each broker. Fine now. At 100x, that's 8.4TB per broker.

**Fixes**:
- **Tiered storage** (KIP-405): Hot data on broker disk (24h), cold data on object storage (S3/GCS)
- **Reduce retention**: 7 days → 3 days on broker (full data in GCS anyway)
- **Monitor disk usage** with alerts at 70% and 85%
- **Log compaction OFF** for audit topic (append-only, delete policy only)
- **Topic-level retention override**: `retention.ms=259200000` (3 days)

---

### 2.6 Schema Evolution & Breaking Changes

**What**: You add a new field to the Avro schema. Old consumers can't read new messages. Or you remove a field that consumers depend on.

**Compatibility modes**:
- **BACKWARD** (default): New schema can read old data. Safe to ADD optional fields with defaults.
- **FORWARD**: Old schema can read new data. Safe to REMOVE optional fields.
- **FULL**: Both backward and forward compatible. Safest.
- **NONE**: No compatibility check. Dangerous.

**Your system**: Should use `BACKWARD` — you'll only ADD fields (new audit metadata), never remove.

**Dangerous scenario**: Someone changes a field type (string → int). Schema Registry rejects it if compatibility is enforced. If someone sets compatibility to NONE to "fix" it → all downstream consumers break.

**Fix**: Lock Schema Registry compatibility to BACKWARD for the audit topic. Add CI/CD check that validates schema changes before deployment.

---

### 2.7 Consumer Lag Cascading

**What**: Consumer falls behind → lag grows → broker keeps more data on disk → disk fills → broker slows down → producers slow down → more lag → cascade.

**Your system**: GCS sink had this during the KEDA incident. Lag hit 500K+.

**Fixes**:
- **Alert at lag thresholds**: Warning at 10K, Critical at 50K, Page at 100K
- **Emergency drain procedure**: Temporarily increase `max.poll.records` from 50 → 500, add more tasks
- **Backfill strategy**: If lag is unrecoverable, reset consumer offset to latest (accept data gap), then backfill from Kafka to GCS using a batch job
- **Log retention as safety net**: Keep Kafka retention at 7 days → 7-day window to catch up before data is lost from Kafka

---

### 2.8 Network Partition Between Kafka and Producers

**What**: Network issue between your services and Kafka brokers. Producers can't send.

**Your system**: Dual-region failover (EUS2 → SCUS). But what if BOTH regions are unreachable?

**Fixes**:
- **Local buffer**: `buffer.memory=64MB` in Kafka producer → buffers locally for ~30 seconds of events
- **Circuit breaker**: After 3 consecutive failures, stop trying for 30 seconds → prevents thread exhaustion
- **Fallback to disk**: Write to local file, background thread retries when Kafka is back
- **Accept data loss**: For audit logs, if Kafka is down for 10+ minutes, it's acceptable to drop events (best-effort system). Alert and document the gap.

**Interview answer**: *"We have dual-region failover. If both are down, Kafka producer's internal buffer holds ~30s of events. Beyond that, we'd accept data loss for audit logs — they're best-effort. But I'd add a circuit breaker to prevent thread exhaustion and alert immediately so we can investigate."*

---

### 2.9 Kafka Connect Worker Crash During Write

**What**: Kafka Connect writes a Parquet file to GCS, then crashes before committing the offset back to Kafka. On restart, it re-reads the same messages and writes them AGAIN.

**Result**: Duplicate files in GCS.

**Fixes**:
- **Deterministic file naming**: `topic_partition_startOffset_endOffset.parquet` → re-writing same file is idempotent (GCS overwrites)
- **BigQuery deduplication**: `SELECT DISTINCT request_id, ...` or use `ROW_NUMBER() OVER (PARTITION BY request_id ORDER BY created_ts) = 1`
- **Exactly-once mode** in Kafka Connect (if connector supports it): Wraps read + write + commit in a transaction

---

### 2.10 Head-of-Line Blocking in Partitions

**What**: One partition has a slow consumer (GCS write timeout). It blocks progress on that partition while other partitions are caught up. Overall lag looks fine but one partition has growing lag.

**Diagnosis**: Check per-partition lag, not just total lag.

**Fix**:
- Monitor **per-partition consumer lag** (not just aggregate)
- Set `max.poll.interval.ms` appropriately so slow partitions trigger rebalance
- Investigate GCS write latency per bucket (maybe one region's GCS is slow)

---

## 3. Cross-System Failure Modes — Walmart Projects

These are issues that could hit ANY of your Walmart services (audit-api-logs-srv, cp-nrti-apis, inventory-status-srv, inventory-events-srv).

---

### 3.1 Database Connection Pool Exhaustion

**Affects**: cp-nrti-apis, inventory-status-srv, inventory-events-srv (all use PostgreSQL)

**What**: Under high load, all DB connections are in use. New requests wait → timeouts → 503 errors cascade.

**Root cause**: Default HikariCP pool = 10 connections. If each query takes 50ms, max throughput = 200 queries/sec. Beyond that, requests queue.

**Symptoms**:
- `HikariPool-1 - Connection is not available, request timed out after 30000ms`
- Response times spike from 50ms → 30,000ms (pool wait timeout)
- 5xx error rate spikes

**Fixes**:
- **Right-size pool**: `maximumPoolSize = (CPU cores × 2) + effective_spindle_count`. For 4-core pod: 10-12 connections is optimal
- **Connection leak detection**: `leakDetectionThreshold=60000` (60s) → logs stack trace of any connection held >60s
- **Query optimization**: Slow queries hold connections longer. Add indexes, optimize JOINs
- **Read replicas**: Route read-heavy queries (inventory search) to read replicas
- **Circuit breaker on DB**: Resilience4j circuit breaker — if DB fails 5 times in 10s, open circuit → fast-fail instead of blocking

**Interview answer**: *"I'd monitor HikariCP metrics (active, idle, waiting) in Grafana. If connections are always maxed, either queries are too slow (optimize) or traffic is too high (scale pods or add read replicas). I'd also add leak detection to catch unclosed connections."*

---

### 3.2 Cascading Timeout Failures

**Affects**: The entire Tier 1 → Tier 2 audit flow, any service-to-service call

**What**: Service A calls Service B with 30s timeout. Service B is slow (DB issue). Service A's threads all block waiting for B. Service A's own callers timeout. Cascade propagates upstream.

**Example in your system**:
- Supplier calls `inventory-status-srv` → it calls internal DB + audit-api-logs-srv
- If audit-api-logs-srv is slow, inventory-status-srv threads block (even though audit is non-critical)
- Supplier gets 504 Gateway Timeout for an inventory query because the AUDIT system is slow

**Fixes**:
- **Bulkhead pattern**: Isolate audit calls in their own thread pool (already done! Your @Async with dedicated ThreadPool is exactly this)
- **Aggressive timeouts**: Audit HTTP call = 2s max (not 30s). If audit is slow, drop it
- **Circuit breaker**: If audit fails 3 times → stop calling for 30s → auto-recover
- **Async fire-and-forget** (your approach): Best fix. The API thread never blocks on audit

**Interview answer**: *"This is exactly why I designed the audit library as fire-and-forget with @Async. The API thread returns immediately. Audit runs in a separate thread pool. If audit is slow or down, the API response is unaffected. This is the bulkhead pattern — isolating non-critical dependencies."*

---

### 3.3 Multi-Tenant Data Leakage

**Affects**: cp-nrti-apis (US/CA/MX), inventory-status-srv, inventory-events-srv

**What**: A US supplier accidentally sees Canadian supplier data because the `wm-site-id` header is missing or wrong.

**Your system**: SiteConfigFactory pattern routes requests to the correct site-specific configuration. But if the header is missing...

**Fixes**:
- **Strict validation**: Reject requests without `wm-site-id` header (don't default to US)
- **Row-level security**: PostgreSQL RLS policies → even if code has a bug, DB enforces tenant isolation
- **Request context propagation**: Store site-id in `ThreadLocal` / `MDC` at the filter level → available everywhere without passing as parameter
- **Audit trail**: Your audit logging system captures `wm-site-id` header → post-hoc detection of cross-tenant access
- **Integration tests**: Test each API with US/CA/MX headers, verify data isolation

---

### 3.4 Secrets Rotation Failure

**Affects**: All services (Kafka credentials, DB passwords, API keys, GCS service accounts)

**What**: Secrets (Kafka broker password, DB credentials) are rotated by the security team. Your service still uses old credentials → authentication failures → all requests fail.

**Symptoms**: Sudden 100% error rate. Logs show `AuthenticationException` or `SASL authentication failed`.

**Fixes**:
- **Dynamic secret loading**: Use Vault/CCM with auto-refresh (not hardcoded in config)
- **Graceful rotation**: Support TWO active credentials simultaneously during rotation window
- **Health checks**: Proactive health endpoint that validates DB + Kafka connectivity
- **Alert on auth failures**: Any `AuthenticationException` = P1 alert

---

### 3.5 API Proxy / Gateway Limits

**Affects**: All external-facing APIs (supplier-facing endpoints)

**What**: Walmart API Gateway (Kong/Azure Front Door) has hidden limits — payload size, rate limits, header size, connection count. You hit them at scale.

**Your real example**: 413 Request Entity Too Large (PRs #49-51). Default 1MB limit vs. large inventory responses.

**Other hidden limits to watch**:
- **Header size**: 8KB default. If you forward too many headers → 431
- **Rate limit per consumer**: If one supplier sends 1000 req/sec → throttled
- **Connection timeout**: Gateway closes idle connections after 60s → pooled connections die silently
- **Request body buffering**: Gateway might buffer entire request body before forwarding → adds latency for large payloads

**Fix**: Document ALL infrastructure limits in a runbook. Test with realistic payload sizes. Monitor 4xx responses grouped by status code.

---

### 3.6 Rolling Deployment → Momentary Inconsistency

**Affects**: All services during deployment

**What**: During rolling deployment, some pods run v1, others run v2. If v2 changes response format or adds required fields, clients get inconsistent responses depending on which pod they hit.

**Your approach**: Canary with Flagger — 10% traffic increments, automatic rollback on error rate spike.

**Still dangerous**:
- **Database schema change**: v2 adds a column, v1 doesn't know about it → if column is NOT NULL, v1 inserts fail
- **Kafka schema change**: v2 produces new Avro schema, v1 consumers can't deserialize → poison pill

**Fixes**:
- **Backward-compatible DB migrations** (your Flyway approach): Always ADD nullable columns first, backfill, then add constraints in next release
- **Schema Registry BACKWARD mode**: New schema can read old data, old schema can read new data (if FULL compatibility)
- **Feature flags**: Deploy code behind flag, enable flag after ALL pods are on v2

---

## 4. Cross-System Failure Modes — GCC Projects

These affect your GCC systems (Beat scraping, Event-gRPC, Stir ETL, Coffee API, Fake Follower ML).

---

### 4.1 Rate Limiting Cascade (Beat Scraping Engine)

**Affects**: 08-gcc-beat-scraping (Python async scraper)

**What**: Instagram/YouTube APIs rate-limit your scraper. You have 3-level rate limiting (daily/minute/handle). But what if:
- Instagram changes rate limits without notice → your daily limit is wrong → account gets banned
- Multiple workers use the same credential simultaneously → exceed per-handle limit

**Fixes**:
- **Adaptive rate limiting**: Track 429 responses. If you get 3 in a row, halve the rate for that API. Double back after 10 minutes.
- **Credential lease system**: Each worker "checks out" a credential from Redis. Lock it. Return when done. No two workers use the same credential.
- **Backoff with jitter**: `sleep(base_delay * 2^attempt + random(0, 1000ms))` — prevents thundering herd when rate limit lifts
- **Graceful degradation**: If Instagram is rate-limited, process YouTube jobs instead. Don't stop the whole pipeline.

---

### 4.2 ClickHouse Write Amplification (Event-gRPC)

**Affects**: 09-gcc-event-grpc (Go gRPC ingestion)

**What**: You use buffered sinker pattern — batch 1000 events, flush to ClickHouse. But ClickHouse MergeTree engine does background merges. Too many small inserts → too many parts → merge overload → "Too many parts" error → inserts rejected.

**Root cause**: If flush interval is too frequent (every 100ms) or batch size too small (<1000 rows), ClickHouse creates too many parts.

**Fixes**:
- **Increase batch size**: Flush every 10,000 events OR every 10 seconds (whichever first)
- **Use Buffer table engine**: ClickHouse's built-in buffer that auto-flushes to MergeTree
- **Monitor parts count**: `SELECT count() FROM system.parts WHERE active AND table = 'events'` → alert if >300
- **Async inserts**: ClickHouse 22.x+ supports `async_insert=1` — server-side batching

**Interview answer**: *"ClickHouse is optimized for large batch inserts, not individual rows. I used a buffered sinker pattern in Go — accumulate events in memory, flush every 10 seconds or 10K events. This keeps part count low and merge overhead manageable."*

---

### 4.3 Airflow DAG Pile-up (Stir Data Platform)

**Affects**: 10-gcc-stir-data-platform (76 DAGs)

**What**: DAG scheduled every 5 minutes. If a run takes >5 minutes, next run queues. If runs keep taking >5 min, queue grows infinitely → Airflow scheduler slows down → ALL DAGs are delayed.

**Symptoms**:
- `Running` DAG runs pile up
- Scheduler heartbeat >5s (normally <1s)
- Other DAGs start missing their schedules

**Fixes**:
- **`max_active_runs_per_dag=1`**: Prevents pile-up. If current run isn't done, skip next scheduled run
- **`dagrun_timeout=timedelta(minutes=15)`**: Kill runs that take too long
- **`catchup=False`**: Don't backfill missed runs (for real-time DAGs)
- **Separate pools**: Critical DAGs in their own pool → not affected by other DAGs' pile-up
- **SLA alerts**: `sla=timedelta(minutes=10)` → alert if DAG hasn't completed in 10 min

---

### 4.4 Dual-Database Consistency (Coffee API)

**Affects**: 11-gcc-coffee-saas-api (PostgreSQL OLTP + ClickHouse OLAP)

**What**: Write goes to PostgreSQL (source of truth). ETL copies to ClickHouse. Read API queries ClickHouse. If ETL is delayed, API returns stale data.

**Scenario**: User creates a campaign at 10:00. ETL runs at 10:05. Between 10:00-10:05, the SaaS dashboard shows old data.

**Fixes**:
- **Write-through**: Write to BOTH databases synchronously. Slower writes but consistent reads.
- **Read from PostgreSQL for recent data**: If `created_at > now() - 5 minutes`, query PostgreSQL. Otherwise ClickHouse.
- **Event-driven sync**: Instead of scheduled ETL, use CDC (Change Data Capture) → near-real-time sync
- **UI indicator**: Show "Last updated: 10:00" in the dashboard → user knows data freshness

---

### 4.5 Serverless Cold Start (Fake Follower ML)

**Affects**: 12-gcc-fake-follower-ml (AWS Lambda)

**What**: Lambda cold start = 3-5 seconds for Python + ML model loading. If traffic is bursty (batch of 10K followers to analyze), first 100 Lambda invocations all cold-start simultaneously.

**Fixes**:
- **Provisioned concurrency**: Keep 10 Lambda instances warm ($$$)
- **Warm-up pings**: CloudWatch scheduled event every 5 min → keeps at least 1 instance warm
- **Batching**: Instead of 1 Lambda per follower, process 100 followers per invocation → fewer cold starts
- **Move to container**: If cold starts are unacceptable, move to ECS Fargate → always-on, no cold start

---

## 5. Database Scaling Issues

### PostgreSQL

| Issue | When It Hits | Fix |
|-------|-------------|-----|
| **Lock contention** | Concurrent UPDATEs on same rows (e.g., inventory count updates) | Use `SELECT ... FOR UPDATE SKIP LOCKED` (your Beat scraping already does this!) |
| **Bloat from dead tuples** | Heavy UPDATE/DELETE workload → VACUUM can't keep up | Tune `autovacuum_vacuum_cost_delay=2`, increase `autovacuum_max_workers=6` |
| **Connection exhaustion** | Many pods × HikariCP pool size > max_connections | Use PgBouncer connection pooler between services and PostgreSQL |
| **Slow sequential scans** | Missing indexes on frequently queried columns | `EXPLAIN ANALYZE` → add indexes. Partial indexes for WHERE clauses |
| **Replication lag** | Read replicas fall behind during write spikes | Monitor `pg_stat_replication.replay_lag`. If >5s, route reads to primary |
| **Table size > RAM** | Large tables can't fit in shared_buffers → disk I/O spikes | Partition by date (`PARTITION BY RANGE (created_at)`), drop old partitions |

### ClickHouse

| Issue | When It Hits | Fix |
|-------|-------------|-----|
| **Too many parts** | Frequent small inserts (< 1000 rows) | Batch inserts: 10K+ rows per INSERT. Use Buffer table engine |
| **Merge overload** | Parts accumulate faster than background merges | Reduce insert frequency. Monitor `system.merges` |
| **Memory OOM** | Complex GROUP BY on high-cardinality columns | Set `max_memory_usage=10G` per query. Use `optimize_aggregation_in_order=1` |
| **Distributed query failures** | Cross-shard queries timeout | Use `distributed_group_by_no_merge=1` for simple aggregations |
| **ZooKeeper bottleneck** | Replicated tables with heavy writes → ZK can't keep up | Migrate to ClickHouse Keeper (built-in, replaces ZK) |

### Redis

| Issue | When It Hits | Fix |
|-------|-------------|-----|
| **Memory full** | Cache grows unbounded → OOM → crash | Set `maxmemory-policy allkeys-lru`. Always set TTL on keys |
| **Hot key** | One key read 100K times/sec → single shard overloaded | Replicate hot keys with random suffix: `key:1`, `key:2`, read from random |
| **Thundering herd** | Cache expires → 1000 requests hit DB simultaneously | **Cache stampede protection**: Lock + single-flight pattern. First request refreshes cache, others wait |
| **Split brain** | Redis Sentinel promotes wrong node → two masters → data divergence | Use Redis Cluster instead. Or set `min-slaves-to-write=1` |

---

## 6. Distributed Systems Failure Patterns

These are universal patterns that apply across ALL your projects.

### 6.1 Thundering Herd

**What**: A resource becomes available (cache refills, service recovers) and ALL waiting clients rush in simultaneously → resource immediately overwhelmed again.

**Your exposure**: Redis cache TTL expiry in Coffee SaaS API (Ristretto + Redis 2-layer cache). If Redis key expires, all 100 concurrent requests hit PostgreSQL.

**Fix**: **Cache stampede lock**. First request acquires lock, refreshes cache. Others get stale data until refresh completes. Or use probabilistic early expiration (refresh at 80% of TTL with some randomness).

### 6.2 Retry Amplification

**What**: Service A calls Service B. B is slow. A retries 3 times. But A has 100 instances → 300 retries → B is even more overwhelmed → cascade.

**Fix**:
- **Exponential backoff with jitter**: `delay = min(base * 2^attempt, max_delay) + random(0, base)`
- **Retry budget**: Max 10% of requests can be retries. If retry rate > 10%, stop retrying
- **Circuit breaker**: After N failures, stop calling entirely for cooldown period

### 6.3 Gray Failures

**What**: Service is not fully down, but degraded. Returns 200 but with wrong/stale data. Health check says "healthy" but responses are garbage.

**Example**: PostgreSQL read replica has replication lag of 5 minutes. Service returns 200 with stale data. Health check queries primary (healthy) but reads go to replica (stale).

**Fix**:
- **Deep health checks**: Health endpoint queries the SAME database the application uses (not just primary)
- **Consistency checks**: Compare response data freshness (e.g., `max(updated_at)` should be within 30s)
- **Canary queries**: Known-answer queries that validate data correctness, not just connectivity

### 6.4 Back-Pressure Not Propagated

**What**: Your system accepts requests faster than it can process them. No mechanism to tell callers to slow down. Memory grows → OOM → crash.

**Your exposure**:
- Tier 1 audit library: ThreadPool queue fills → RejectedExecutionException (backpressure via dropping)
- Beat scraper: 150 workers sending to Event-gRPC → if gRPC server is slow, workers pile up
- Airflow DAGs: 76 DAGs × multiple tasks → Airflow scheduler overwhelmed

**Fix**:
- **Bounded queues** everywhere (your ThreadPool queue=100 is correct)
- **HTTP 429 responses** with Retry-After header
- **Kafka producer `buffer.memory`** limit → blocks producer if buffer is full (natural backpressure)
- **gRPC flow control**: Built-in. Server can signal client to slow down

### 6.5 Clock Skew

**What**: Your Tier 1 records `request_ts` and `response_ts` using `Instant.now()`. If two pods have clock skew of 5 seconds, audit logs show impossible timelines (response before request on a different pod).

**Fix**:
- **NTP sync** on all pods (usually handled by cloud provider)
- **Use monotonic clocks** for duration measurement (`System.nanoTime()` not `System.currentTimeMillis()`)
- **Include pod ID** in audit log → post-hoc correction if clock skew detected

---

## 7. Interview "What If" Questions — Complete Answer Bank

### Kafka-Specific

**Q: What if a Kafka broker goes down?**
> "With replication factor 3, we can lose 2 brokers and still serve reads/writes (min.insync.replicas=2). Kafka automatically elects a new partition leader from the remaining replicas. Producers and consumers failover transparently. In our dual-region setup, if an entire AZ goes down, we failover to SCUS region."

**Q: What if Schema Registry goes down?**
> "Producers and consumers cache the schema locally after first fetch. So existing producers/consumers continue working. But NEW schema versions can't be registered, and new consumer instances can't start (no schema to deserialize). Fix: Schema Registry should be HA (3-node cluster). We'd also add a local schema cache with TTL as fallback."

**Q: What if message ordering matters?**
> "Ordering is guaranteed within a partition. Our key is service_name|endpoint_name, so all events for the same service+endpoint are ordered. If we need global ordering, we'd use a single partition (kills parallelism) or add sequence numbers and sort at query time."

**Q: What if Kafka is down for 1 hour?**
> "Tier 1 library's @Async pool fills up (100 items → ~200KB). Beyond that, audit events are dropped. When Kafka recovers, new events flow again. We'd have a 1-hour data gap in audit logs. For audit (best-effort), this is acceptable. We'd alert immediately and document the gap. If critical, we could add a local file buffer — write to disk when Kafka is down, drain file when Kafka recovers."

**Q: What if you need to replay all messages from the beginning?**
> "Reset consumer group offset to earliest: `kafka-consumer-groups --reset-offsets --to-earliest`. GCS writes are idempotent (deterministic file names), so replay = overwrite same files. BigQuery queries are unaffected. Replay 7 days of data (Kafka retention) takes ~2-3 hours with parallelism."

### System Design

**Q: What if we add a new region (Europe)?**
> "1) Deploy Tier 2 publisher in EU region. 2) Create new GCS bucket (audit-api-logs-eu-prod) in EU region (data residency). 3) Add new SMT filter (AuditLogSinkEUFilter) in GCS sink. 4) Add wm-site-id=EU header routing. 5) Consider GDPR implications — EU audit logs might need deletion capability (currently append-only). 6) Latency consideration: EU services should publish to EU Kafka cluster, not US. Needs MirrorMaker 2 for cross-region replication if global querying is needed."

**Q: How would you make this real-time queryable (not just batch)?**
> "Current: GCS → BigQuery (batch, minutes delay). For real-time: 1) Add a Kafka Streams application that reads the audit topic and writes to an Elasticsearch index (or ClickHouse). 2) Build a REST API on top for real-time queries. 3) Or use ksqlDB for streaming SQL queries directly on the Kafka topic. 4) Trade-off: real-time adds complexity + cost. Current batch (5-min delay) is sufficient for audit use cases."

**Q: How do you test this system end-to-end?**
> "1) Unit tests: MockKafkaProducer, WireMock for HTTP calls. 2) Integration tests: Testcontainers with embedded Kafka + PostgreSQL. 3) Contract tests: Testburst for API contracts between tiers. 4) End-to-end: Deploy to staging, send test requests, verify data appears in GCS within 5 minutes. 5) Chaos testing: Kill a Kafka broker in staging, verify failover works. 6) Load testing: Gatling/k6 to simulate 10x traffic and measure latency + data completeness."

**Q: What's the weakest part of this system?**
> "The HTTP hop between Tier 1 and Tier 2. It adds latency, creates a coupling point, and is the first thing that breaks at scale. If I were starting over, I'd embed a lightweight Kafka producer directly in the library. The reason we have the HTTP hop is historical — when we started, not all consuming services had Kafka dependencies, and we didn't want to force that. But with 12+ teams now on board, it's the obvious next optimization."

### Behavioral

**Q: What was the hardest problem you solved?**
> "The Silent Failure. 5 days of data loss. What made it hard wasn't any single bug — it was 4 bugs compounding simultaneously: NPE in SMT filter, KEDA rebalance storm, JVM heap exhaustion, and silent error tolerance. Each fix revealed the next issue. I learned that Kafka Connect's defaults are not production-ready, and that `errors.tolerance=all` is dangerous without DLQ monitoring."

**Q: What would you do differently?**
> "Two things: 1) I'd embed the Kafka producer directly in the library — the HTTP middle tier adds complexity without proportional value. 2) I'd set up comprehensive monitoring BEFORE production launch, not after. The Silent Failure would have been caught in minutes if we'd had consumer lag alerts from day one."

**Q: How did you get 12+ teams to adopt your library?**
> "Three things: 1) Made integration trivial — single POM dependency, zero code changes, auto-configures via Spring Boot starter. 2) Showed concrete value — before vs. after comparison (Splunk: $$$, 30-day retention vs. Kafka+GCS: 99% cost reduction, 7-year retention). 3) Hands-on support — paired with each team during integration, hosted brown-bag sessions to explain the architecture."

---

## Quick Reference: Numbers to Memorize

| System | Metric | Value |
|--------|--------|-------|
| Kafka Audit | Events/day | 2M+ |
| Kafka Audit | P99 latency impact | <5ms |
| Kafka Audit | Cost reduction vs Splunk | 99% |
| Kafka Audit | Parquet compression | 90% |
| Kafka Audit | Avro vs JSON size | 70% smaller |
| Kafka Audit | Teams adopted | 12+ |
| Kafka Audit | Integration time | 1 day (vs 2-3 weeks custom) |
| Kafka Audit | DR failover | 15 minutes |
| DSD Notification | Associates notified | 1,200+ |
| DSD Notification | Stores | 300+ |
| DSD Notification | Replenishment improvement | 35% |
| Spring Boot Migration | Downtime | Zero |
| OpenAPI | Integration overhead reduction | 30% |
| GCC Beat | Data points/day | 10M+ |
| GCC Beat | Images/day | 8M |
| GCC Beat | Concurrent workers | 150+ |
| GCC Event-gRPC | Events/sec | 10K+ |
| GCC Event-gRPC | Event types | 60+ |
| GCC ClickHouse | Log retrieval speedup | 2.5x |
| GCC Stir | Airflow DAGs | 76 |
| GCC Stir | dbt models | 112 |
| GCC Stir | Data latency reduction | 50% |
| GCC Coffee | API endpoints | 50+ |
| GCC Coffee | Response time improvement | 25% |
| GCC Coffee | Cost reduction | 30% |
| PayU | Business scale increase | 40% |
| PayU | Error rate reduction | 4.6% → 0.3% |
| PayU | TAT reduction | 3.2 min → 1.1 min |
| PayU | Unit test coverage | 30% → 83% |
| PayU | Deployment errors reduction | 90% |

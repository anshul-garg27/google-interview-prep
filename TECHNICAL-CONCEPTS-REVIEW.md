# 🔧 TECHNICAL CONCEPTS DEEP DIVE

> **Purpose:** Understand the "why" behind every technical decision for interview discussions

---

## TABLE OF CONTENTS

1. [Kafka Connect & SMT Filters](#1-kafka-connect--smt-filters)
2. [Async Processing Patterns](#2-async-processing-patterns)
3. [ContentCachingWrapper](#3-contentcachingwrapper)
4. [Avro Serialization](#4-avro-serialization)
5. [Parquet File Format](#5-parquet-file-format)
6. [Spring Boot 3 Migration Concepts](#6-spring-boot-3-migration-concepts)
7. [WebClient vs RestTemplate](#7-webclient-vs-resttemplate)
8. [CompletableFuture vs ListenableFuture](#8-completablefuture-vs-listenablefuture)
9. [Canary Deployments with Flagger](#9-canary-deployments-with-flagger)
10. [Multi-Region Architecture](#10-multi-region-architecture)

---

## 1. KAFKA CONNECT & SMT FILTERS

### What is Kafka Connect?

Kafka Connect is a **framework for streaming data between Kafka and external systems** without writing custom code.

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         KAFKA CONNECT ARCHITECTURE                           │
│                                                                              │
│  ┌─────────────────┐                          ┌─────────────────┐           │
│  │  Source         │    Kafka Topics          │  Sink           │           │
│  │  Connectors     │                          │  Connectors     │           │
│  │                 │     ┌─────────┐          │                 │           │
│  │  - JDBC Source  │ ──► │ topic-1 │ ──────► │  - GCS Sink     │           │
│  │  - Debezium     │     │ topic-2 │          │  - JDBC Sink    │           │
│  │  - File Source  │     │ topic-3 │          │  - S3 Sink      │           │
│  └─────────────────┘     └─────────┘          └─────────────────┘           │
│                                                                              │
│  Worker 1              Worker 2              Worker 3                        │
│  ├── Task 1            ├── Task 1            ├── Task 1                     │
│  └── Task 2            └── Task 2            └── Task 2                     │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Why Kafka Connect over Custom Consumer?

| Aspect | Custom Consumer | Kafka Connect |
|--------|-----------------|---------------|
| **Offset Management** | You implement it | Built-in (automatic) |
| **Scaling** | You implement it | Automatic task distribution |
| **Fault Tolerance** | You implement it | Automatic rebalancing |
| **Monitoring** | Custom metrics | Standard JMX metrics |
| **Error Handling** | You implement it | DLQ, retries built-in |
| **Connector Ecosystem** | N/A | 100+ pre-built connectors |

### What is SMT (Single Message Transform)?

SMT allows you to **modify messages inline** as they flow through Connect.

```java
// Your SMT Filter Implementation
public class AuditLogSinkUSFilter<R extends ConnectRecord<R>>
    extends BaseAuditLogSinkFilter<R> {

    @Override
    protected String getHeaderValue() {
        return "US";  // Target site ID
    }

    @Override
    public boolean verifyHeader(R record) {
        // US filter is PERMISSIVE:
        // Accept if wm-site-id = "US" OR if header is missing (legacy data)
        return hasUSHeader(record) || hasNoSiteIdHeader(record);
    }
}
```

### How SMT Filters Work in Your System

```
┌────────────────────────────────────────────────────────────────────────────┐
│                        SMT FILTER PIPELINE                                  │
│                                                                             │
│  Kafka Message arrives:                                                     │
│  {payload: {...}, headers: {wm-site-id: "US"}}                             │
│                                                                             │
│      │                                                                      │
│      ▼                                                                      │
│  ┌─────────────────┐                                                        │
│  │ Transform 1:    │  InsertRollingRecordTimestamp                         │
│  │ Add date header │  → Adds partition date for GCS path                   │
│  └────────┬────────┘                                                        │
│           │                                                                 │
│           ▼                                                                 │
│  ┌─────────────────┐                                                        │
│  │ Transform 2:    │  FilterUS / FilterCA / FilterMX                       │
│  │ Filter by site  │  → Returns record OR null (drops it)                  │
│  └────────┬────────┘                                                        │
│           │                                                                 │
│           ▼                                                                 │
│  Record passed to GCS Sink Connector                                        │
│  Written to: gs://audit-logs-us-prod/2025-02-05/...                        │
│                                                                             │
└────────────────────────────────────────────────────────────────────────────┘
```

### Interview Q&A

**Q: Why SMT over post-processing?**
> "SMT processes messages inline with zero external dependencies. No additional service to maintain. The filtering happens at consumption time, so only relevant records reach each GCS bucket. Lower latency, simpler architecture."

**Q: What happens when SMT returns null?**
> "Kafka Connect convention - returning null from apply() drops the message. It won't be written to the sink. That's how we route: US connector's filter returns null for non-US records."

---

## 2. ASYNC PROCESSING PATTERNS

### The Problem: Blocking vs Non-Blocking

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    SYNCHRONOUS (Blocking) - BAD!                             │
│                                                                              │
│  Client Request ─────► API ─────► Audit Service ─────► Response             │
│                              │                    │                          │
│                              │    BLOCKING WAIT   │                          │
│                              │◄───────────────────┤                          │
│                              │                                               │
│  Total Time = API Processing + Audit Time (300ms + 50ms = 350ms)            │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│                    ASYNCHRONOUS (Non-Blocking) - GOOD!                       │
│                                                                              │
│  Client Request ─────► API ─────► Response (returned immediately!)          │
│                              │                                               │
│                              └───► @Async ───► Audit Service (background)   │
│                                                                              │
│  Total Time = API Processing Only (300ms)                                   │
│  Audit happens in background, user doesn't wait!                            │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Your Implementation: @Async + ThreadPoolTaskExecutor

```java
@Configuration
@EnableAsync  // Enables @Async annotation processing
public class AuditLogAsyncConfig {

    @Bean
    public Executor taskExecutor() {
        ThreadPoolTaskExecutor executor = new ThreadPoolTaskExecutor();

        executor.setCorePoolSize(6);       // Always 6 threads running
        executor.setMaxPoolSize(10);       // Can grow to 10 under load
        executor.setQueueCapacity(100);    // Buffer 100 tasks before rejection
        executor.setThreadNamePrefix("Audit-log-executor-");

        return executor;
    }
}
```

### Why These Thread Pool Values?

| Setting | Value | Why |
|---------|-------|-----|
| **corePoolSize: 6** | Always-running threads | ~100 req/sec × 50ms/audit = 5 threads needed; +1 buffer |
| **maxPoolSize: 10** | Peak capacity | Handles traffic spikes without thread explosion |
| **queueCapacity: 100** | Task buffer | Absorbs bursts; prevents OOM from unbounded queue |

### What Happens Under Load?

```
Requests: 50/sec     →  6 core threads handle easily
Requests: 100/sec    →  6 threads busy, tasks queue (0-100 in queue)
Requests: 150/sec    →  Queue filling, scale to 10 threads
Requests: 200/sec    →  Queue full (100), scale to 10 maxed out
Requests: 250/sec    →  RejectedExecutionException! Tasks dropped!
```

### Fire-and-Forget Pattern

```java
@Service
public class AuditLogService {

    @Async  // Magic annotation - runs in thread pool
    public void sendAuditLogRequest(AuditLogPayload payload) {
        try {
            httpService.post(auditServiceUrl, payload);
        } catch (Exception e) {
            // LOG but don't THROW
            log.error("Audit failed for {}", payload.getRequestId(), e);
            // User experience is NOT affected!
        }
    }
}
```

### Interview Q&A

**Q: What's the downside of fire-and-forget?**
> "We can lose audit data if the queue is full or if the audit service is down. This is acceptable for audit logging - it's observability data, not critical business data. The API must never fail because of audit logging."

**Q: How do you know if you're losing audit data?**
> "Metrics. We track rejected_count from the executor and alerting on audit service error rates. If queue reaches 80%, we get a warning."

---

## 3. CONTENTCACHINGWRAPPER

### The Problem: HTTP Streams Can Only Be Read Once

```java
// HTTP request body is an InputStream
// Read once = EMPTY for everyone else!

// In Filter:
String body = request.getReader().lines().collect(joining());  // ✅ Works

// In Controller (called after Filter):
String body = request.getReader().lines().collect(joining());  // ❌ EMPTY!
```

### The Solution: ContentCachingWrapper

```java
// ContentCachingWrapper reads the stream ONCE and CACHES the bytes
// Everyone can read from the cache!

ContentCachingRequestWrapper wrappedRequest =
    new ContentCachingRequestWrapper(request);

// Let the filter chain (including controller) run
filterChain.doFilter(wrappedRequest, response);

// NOW we can read the body (from cache!)
byte[] body = wrappedRequest.getContentAsByteArray();  // ✅ Works!
```

### How It Works Internally

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    CONTENTCACHINGWRAPPER FLOW                                │
│                                                                              │
│  1. Original Request arrives                                                │
│     └── Body: InputStream (can only read once)                             │
│                                                                              │
│  2. Wrap with ContentCachingRequestWrapper                                  │
│     └── Internal: ByteArrayOutputStream cachedContent                       │
│                                                                              │
│  3. Controller calls request.getReader()                                    │
│     └── Wrapper: Reads InputStream → Writes to cachedContent               │
│     └── Controller: Gets the data ✅                                        │
│                                                                              │
│  4. Filter calls getContentAsByteArray()                                    │
│     └── Returns cachedContent.toByteArray() ✅                              │
│                                                                              │
│  Result: Both Controller AND Filter can read the body!                      │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Interview Q&A

**Q: What's the memory impact?**
> "The entire body is cached in memory. For large payloads, this could be problematic. We mitigate with gateway-level payload limits (2MB max) and option to skip body logging for specific endpoints."

**Q: Why not just log request parameters?**
> "For POST/PUT APIs with JSON bodies, the actual data IS the body. URL parameters alone don't capture what the supplier sent. We need full body for debugging 'why did this request fail?'"

---

## 4. AVRO SERIALIZATION

### What is Avro?

Apache Avro is a **binary serialization format with schema enforcement**.

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      JSON vs AVRO COMPARISON                                 │
│                                                                              │
│  JSON (Text):                                                               │
│  {"request_id":"abc-123","service_name":"nrti","response_code":200}        │
│  Size: ~100 bytes                                                           │
│  ❌ No schema enforcement                                                    │
│  ❌ Field names repeated in every record                                     │
│                                                                              │
│  AVRO (Binary):                                                             │
│  [magic byte][schema_id][binary_data]                                       │
│  Size: ~30 bytes (70% smaller!)                                             │
│  ✅ Schema enforced by Schema Registry                                      │
│  ✅ Field names defined once in schema                                      │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Schema Registry Integration

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    AVRO + SCHEMA REGISTRY FLOW                               │
│                                                                              │
│  PRODUCER:                                                                  │
│  1. Producer has schema: {type: record, name: LogEvent, fields: [...]}     │
│  2. Registers schema → Gets schema_id = 123                                 │
│  3. Serializes payload to binary                                            │
│  4. Sends: [0][123][binary_data] to Kafka                                  │
│                                                                              │
│  CONSUMER:                                                                  │
│  1. Receives: [0][123][binary_data]                                        │
│  2. Extracts schema_id = 123                                                │
│  3. Fetches schema from Registry (cached)                                   │
│  4. Deserializes binary_data using schema                                   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Schema Evolution (Why Avro is Powerful)

```java
// Version 1: Original schema
{
  "type": "record",
  "name": "LogEvent",
  "fields": [
    {"name": "request_id", "type": "string"},
    {"name": "service_name", "type": "string"}
  ]
}

// Version 2: Add new field with default (BACKWARD COMPATIBLE!)
{
  "type": "record",
  "name": "LogEvent",
  "fields": [
    {"name": "request_id", "type": "string"},
    {"name": "service_name", "type": "string"},
    {"name": "trace_id", "type": ["null", "string"], "default": null}  // NEW!
  ]
}

// Old consumers can read new data (ignore trace_id)
// New consumers can read old data (trace_id = null)
```

### Interview Q&A

**Q: Why Avro over Protobuf?**
> "Both are good. We chose Avro because it integrates naturally with Kafka ecosystem (Confluent Schema Registry) and has first-class support in Kafka Connect. Protobuf would work too."

**Q: What if Schema Registry is down?**
> "Producers cache schemas locally after first fetch. Existing schemas continue to work. New schemas can't be registered until Registry recovers. This is a known limitation - we monitor Registry health."

---

## 5. PARQUET FILE FORMAT

### What is Parquet?

Apache Parquet is a **columnar storage format** optimized for analytics.

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    ROW vs COLUMNAR STORAGE                                   │
│                                                                              │
│  ROW-BASED (JSON, CSV):                                                     │
│  ┌────────────────────────────────────────┐                                 │
│  │ {id:1, name:"A", status:200}           │  Each row stored together       │
│  │ {id:2, name:"B", status:400}           │  To read "status" column:       │
│  │ {id:3, name:"C", status:200}           │  Must read ENTIRE file!         │
│  └────────────────────────────────────────┘                                 │
│                                                                              │
│  COLUMNAR (Parquet):                                                        │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐                                    │
│  │ id       │ │ name     │ │ status   │  Each column stored together       │
│  │ 1        │ │ "A"      │ │ 200      │  To read "status" column:          │
│  │ 2        │ │ "B"      │ │ 400      │  Read ONLY status column!          │
│  │ 3        │ │ "C"      │ │ 200      │  90% less data scanned!            │
│  └──────────┘ └──────────┘ └──────────┘                                    │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Why Parquet for Audit Logs?

| Benefit | How It Helps |
|---------|-------------|
| **Columnar storage** | Queries like "all response_codes" read only that column |
| **Compression** | Similar values in column → excellent compression (90%!) |
| **BigQuery native** | BQ reads Parquet directly, no transformation |
| **Schema in file** | Self-describing, no external schema needed |
| **Predicate pushdown** | "WHERE status > 400" skips entire row groups |

### BigQuery External Table

```sql
-- BigQuery can query GCS Parquet files directly!
CREATE EXTERNAL TABLE `project.audit_logs.api_logs_us`
OPTIONS (
  format = 'PARQUET',
  uris = ['gs://audit-api-logs-us-prod/*.parquet']
);

-- Now suppliers can query their own data!
SELECT request_id, endpoint_name, response_code
FROM `project.audit_logs.api_logs_us`
WHERE DATE(request_ts) = '2025-02-05'
  AND consumer_id = 'pepsi-supplier-123'
  AND response_code >= 400;
```

### Interview Q&A

**Q: Why not store in BigQuery directly?**
> "Cost. BigQuery streaming inserts are expensive at 2M events/day. Writing to GCS first and using external tables is ~10x cheaper. We get the same query capability."

**Q: How do you handle GDPR delete requests?**
> "Parquet files are immutable. For deletion, we must: identify files containing that user, read-filter-rewrite each file, delete originals. It's expensive. We're considering Apache Iceberg for ACID transactions and easier deletes."

---

## 6. SPRING BOOT 3 MIGRATION CONCEPTS

### javax → jakarta Namespace Change

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    THE GREAT NAMESPACE MIGRATION                             │
│                                                                              │
│  Why did this happen?                                                       │
│  - Oracle owned Java EE (javax.*)                                           │
│  - Oracle donated Java EE to Eclipse Foundation                             │
│  - Eclipse couldn't use "javax" trademark                                   │
│  - Renamed to Jakarta EE (jakarta.*)                                        │
│                                                                              │
│  BEFORE (Spring Boot 2.x):                 AFTER (Spring Boot 3.x):         │
│  import javax.persistence.*;        →      import jakarta.persistence.*;    │
│  import javax.validation.*;         →      import jakarta.validation.*;     │
│  import javax.servlet.*;            →      import jakarta.servlet.*;        │
│                                                                              │
│  74 files changed in your migration!                                        │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Hibernate 6 Changes

```java
// Problem: PostgreSQL custom enum types
// Hibernate 5: Worked implicitly
// Hibernate 6: Needs explicit type mapping!

// BEFORE (Hibernate 5):
@Column(name = "status")
@Enumerated(EnumType.STRING)
private Status status;  // Worked fine

// AFTER (Hibernate 6):
@Column(name = "status", columnDefinition = "status_enum")
@Enumerated(EnumType.STRING)
@JdbcTypeCode(SqlTypes.NAMED_ENUM)  // NEW! Required for PostgreSQL enums
private Status status;
```

### Spring Security 6 Changes

```java
// BEFORE (Spring Security 5):
@Configuration
public class SecurityConfig extends WebSecurityConfigurerAdapter {  // DEPRECATED!
    @Override
    protected void configure(HttpSecurity http) throws Exception {
        http.authorizeRequests()
            .antMatchers("/actuator/**").permitAll();
    }
}

// AFTER (Spring Security 6):
@Configuration
public class SecurityConfig {
    @Bean
    public SecurityFilterChain filterChain(HttpSecurity http) throws Exception {
        http.authorizeHttpRequests(auth -> auth
            .requestMatchers("/actuator/**").permitAll()  // Changed method name!
        );
        return http.build();
    }
}
```

---

## 7. WEBCLIENT VS RESTTEMPLATE

### Why WebClient?

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    RESTTEMPLATE vs WEBCLIENT                                 │
│                                                                              │
│  RestTemplate (Deprecated in Spring 5.0):                                   │
│  - Synchronous/blocking                                                     │
│  - One thread per request (expensive under load)                            │
│  - Simple API                                                               │
│                                                                              │
│  WebClient (Modern):                                                        │
│  - Asynchronous/non-blocking by default                                     │
│  - Uses Reactor's event loop (efficient)                                    │
│  - Can be used blocking with .block()                                       │
│  - More complex API (builder pattern)                                       │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Your Migration Strategy: WebClient with .block()

```java
// You chose .block() for BACKWARDS COMPATIBILITY
// Full reactive would require rewriting all service classes

return webClient.get()
    .uri(url)
    .headers(h -> h.addAll(headers))
    .retrieve()
    .bodyToMono(StoreResponse.class)
    .block();  // ← Blocking! Behaves like RestTemplate

// WHY NOT FULL REACTIVE?
// 1. Scope control - framework upgrade, not architecture change
// 2. Team readiness - reactive has learning curve
// 3. Risk reduction - minimize business logic changes
```

### Interview Q&A

**Q: Why not go fully reactive?**
> "Scope control and risk management. A full reactive migration would touch every service class and fundamentally change error handling. For a framework upgrade, I wanted to minimize changes to business logic. Reactive would be a dedicated project with proper training."

---

## 8. COMPLETABLEFUTURE VS LISTENABLEFUTURE

### The Migration

```java
// BEFORE (Spring Kafka with ListenableFuture):
ListenableFuture<SendResult<String, Message>> future = kafkaTemplate.send(message);
future.addCallback(
    result -> log.info("Success"),
    ex -> log.error("Failed", ex)
);

// AFTER (Spring Kafka with CompletableFuture):
CompletableFuture<SendResult<String, Message>> future = kafkaTemplate.send(message);
future
    .thenAccept(result -> log.info("Success: {}", result))
    .exceptionally(ex -> {
        log.error("Failed", ex);
        handleFailure(message);  // Can chain operations!
        return null;
    })
    .join();  // Wait for completion if needed
```

### Why CompletableFuture is Better

| Aspect | ListenableFuture | CompletableFuture |
|--------|------------------|-------------------|
| **Chaining** | Callbacks only | .thenApply(), .thenCompose(), etc. |
| **Combining** | Manual | .allOf(), .anyOf() |
| **Error handling** | Single callback | .exceptionally(), .handle() |
| **Standard Java** | Spring-specific | java.util.concurrent (standard!) |

---

## 9. CANARY DEPLOYMENTS WITH FLAGGER

### What is Canary Deployment?

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    CANARY DEPLOYMENT STRATEGY                                │
│                                                                              │
│  Instead of: Deploy 100% at once (BIG BANG - risky!)                        │
│  Do: Gradually shift traffic while monitoring                               │
│                                                                              │
│  Hour 0:    100% ──► Old Version                                            │
│  Hour 1:     90% ──► Old Version                                            │
│              10% ──► NEW VERSION (canary)  ← Monitor!                       │
│                                                                              │
│  Hour 2:     75% ──► Old Version                                            │
│              25% ──► NEW VERSION           ← If errors < 1%, continue       │
│                                                                              │
│  Hour 4:     50% ──► Old Version                                            │
│              50% ──► NEW VERSION           ← If errors < 1%, continue       │
│                                                                              │
│  Hour 6:      0% ──► Old Version                                            │
│             100% ──► NEW VERSION           ← Full rollout!                  │
│                                                                              │
│  If errors > 1% at any point: AUTOMATIC ROLLBACK!                           │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Flagger + Istio

```yaml
# Flagger configuration in kitt.yml
flagger:
  enabled: true
  analysis:
    threshold: 5           # Max failed checks before rollback
    maxWeight: 50          # Max percentage to canary
    stepWeight: 10         # Increase by 10% each step
    interval: 1m           # Check every minute

    metrics:
    - name: request-success-rate
      threshold: 99        # Must have 99% success rate
    - name: request-duration
      threshold: 500       # P99 must be < 500ms
```

---

## 10. MULTI-REGION ARCHITECTURE

### Why Multi-Region?

1. **Disaster Recovery**: If EUS2 goes down, SCUS continues
2. **Latency**: Users closer to SCUS get faster responses
3. **Compliance**: Data residency requirements
4. **Load Distribution**: Split traffic across regions

### Your Active/Active Design

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    ACTIVE/ACTIVE KAFKA ARCHITECTURE                          │
│                                                                              │
│                     ┌───────────────────────┐                               │
│                     │   Global Load Balancer │                               │
│                     └───────────┬───────────┘                               │
│                                 │                                            │
│              ┌──────────────────┼──────────────────┐                        │
│              │                  │                  │                        │
│              ▼                  ▼                  ▼                        │
│    ┌─────────────────┐ ┌─────────────────┐ ┌─────────────────┐             │
│    │   EUS2 Region   │ │   SCUS Region   │ │   INTL (Future) │             │
│    │                 │ │                 │ │                 │             │
│    │  audit-api-srv  │ │  audit-api-srv  │ │                 │             │
│    │       │         │ │       │         │ │                 │             │
│    │       ▼         │ │       ▼         │ │                 │             │
│    │  Kafka Cluster  │◄┼─► Kafka Cluster │ │                 │             │
│    │       │         │ │       │         │ │                 │             │
│    │       ▼         │ │       ▼         │ │                 │             │
│    │   GCS Sink      │ │   GCS Sink      │ │                 │             │
│    └─────────────────┘ └─────────────────┘ └─────────────────┘             │
│              │                  │                                           │
│              └──────────┬───────┘                                           │
│                         ▼                                                   │
│              ┌─────────────────┐                                            │
│              │   GCS Buckets   │  (Shared across regions)                  │
│              │   - us-prod     │                                            │
│              │   - ca-prod     │                                            │
│              │   - mx-prod     │                                            │
│              └─────────────────┘                                            │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Dual-Region Kafka Publishing

```java
// Your implementation: Try primary, fallback to secondary
public void publishMessageToTopic(LoggingApiRequest request) {
    Message<LogEvent> message = prepareMessage(request);

    try {
        kafkaPrimaryTemplate.send(message).get();  // EUS2
    } catch (Exception ex) {
        log.warn("Primary Kafka failed, trying secondary: {}", ex.getMessage());
        kafkaSecondaryTemplate.send(message).get();  // SCUS fallback
    }
}
```

---

## QUICK REFERENCE: Technical Terms

| Term | Definition |
|------|------------|
| **SMT** | Single Message Transform - inline message processing in Kafka Connect |
| **@Async** | Spring annotation for non-blocking execution in thread pool |
| **ContentCachingWrapper** | Spring utility to read HTTP body multiple times |
| **Avro** | Binary serialization format with schema enforcement |
| **Parquet** | Columnar file format optimized for analytics |
| **Canary** | Gradual rollout strategy with automatic rollback |
| **Flagger** | Kubernetes operator for canary deployments |
| **Schema Registry** | Central service for Avro schema management |

---

*Last updated: February 2026*

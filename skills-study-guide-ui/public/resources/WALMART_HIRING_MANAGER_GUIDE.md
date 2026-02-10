# WALMART → GOOGLE: HIRING MANAGER ROUND GUIDE
## Google L4/L5 Technical Leadership Interview

**CRITICAL**: Hiring Manager rounds at Google are NOT coding interviews. They test:
1. **System Design Thinking** - Can you architect complex systems?
2. **Technical Decision Making** - Can you evaluate trade-offs?
3. **Leadership & Influence** - Can you drive technical direction?
4. **Scale & Complexity** - Have you built production systems at scale?

**Your Advantage**: You have 6 production services handling 8M+ queries/month across 3 countries. Use this.

---

## TABLE OF CONTENTS

### PART A: "WALK ME THROUGH YOUR SYSTEM" DEEP DIVES
1. [Multi-Region Kafka Audit System](#1-multi-region-kafka-audit-system) (Most Complex)
2. [DC Inventory Search with 3-Stage Pipeline](#2-dc-inventory-search-3-stage-pipeline)
3. [Multi-Market Architecture (US/CA/MX)](#3-multi-market-architecture)
4. [Real-Time Event Processing (2M events/day)](#4-real-time-event-processing)
5. [Supplier Authorization Framework](#5-supplier-authorization-framework)

### PART B: TECHNICAL DECISION FRAMEWORKS
6. [Trade-Off Analysis Examples](#6-trade-off-analysis)
7. [Scale Decisions (How to Handle 10x Growth)](#7-scale-decisions)
8. [Failure Scenarios & Resilience](#8-failure-scenarios)

### PART C: LEADERSHIP STORIES
9. [Technical Leadership (Migrations, Mentorship)](#9-technical-leadership)
10. [Cross-Team Collaboration](#10-cross-team-collaboration)

---

# PART A: "WALK ME THROUGH YOUR SYSTEM"

## Google HM Question Style

**They Ask**: "Tell me about the most complex system you've designed and built."

**What They're Evaluating**:
- Can you articulate architecture clearly?
- Do you understand trade-offs?
- How do you handle scale and failure?
- Can you justify technical decisions with data?

**Your Answer Structure** (15-20 minutes):
```
1. Business Context (2 min) - Why did this system need to exist?
2. Technical Challenges (3 min) - What made it complex?
3. Architecture Deep Dive (8 min) - How did you solve it?
4. Scale & Performance (2 min) - How does it handle load?
5. Failures & Resilience (3 min) - What happens when things break?
6. Results & Learnings (2 min) - Impact and what you'd do differently
```

---

## 1. MULTI-REGION KAFKA AUDIT SYSTEM

### Business Context (2 min)

**"Why did this system need to exist?"**

"At Walmart Data Ventures, we had 12 microservices providing APIs to external suppliers and internal analytics teams. The business problems were:

1. **Compliance**: No audit trail of API calls (PCI-DSS requirement for supplier data access)
2. **Debugging**: When suppliers reported issues, we had no request/response logs to troubleshoot
3. **Analytics**: Product team wanted to understand API usage patterns (which suppliers, which endpoints, peak hours)
4. **Cost**: Existing solution was PostgreSQL (50M rows, $5K/month, crashing 2x/month)

The challenge: Build a centralized audit logging system that:
- Captures ALL API traffic (2M events/day) without impacting API latency
- Provides fast analytics queries (< 2 seconds for 30-day supplier history)
- Scales to 10x growth (product roadmap)
- Costs < $1,000/month
- Has disaster recovery (RPO < 1 minute)

Previous attempts failed because they used direct database writes (blocking, slow, not scalable)."

---

### Technical Challenges (3 min)

**"What made this complex?"**

**Challenge 1: Zero Latency Impact**
- APIs serve external suppliers (SLA: < 300ms P95)
- Audit logging CAN'T add latency (business requirement)
- Previous solution (sync JDBC writes) added 40-50ms per request

**Challenge 2: High Volume + Retention**
- 2M events/day = 730M events/year
- Compliance requires 7-year retention = 5B+ events
- PostgreSQL couldn't scale beyond 50M rows

**Challenge 3: Multi-Tenant Isolation**
- US, Canada, Mexico markets must have separate data (compliance)
- Single Kafka topic, but 3 separate GCS buckets
- Zero data leakage between markets

**Challenge 4: Disaster Recovery**
- Single-region Kafka = region failure = total audit loss
- Business requirement: RPO < 1 minute (audit logs critical for compliance)
- SRE team had NO capacity for manual failover

**Challenge 5: Seamless Adoption**
- 12 services owned by 8 different teams
- Teams have no time to rewrite audit logic
- Solution must be 'drop-in' (minimal code changes)

---

### Architecture Deep Dive (8 min)

**Component 1: Client-Side Library (dv-api-common-libraries)**

"First problem: How do we capture API traffic WITHOUT teams rewriting code?

Solution: Spring Servlet Filter with automatic instrumentation.

```java
@Component
@Order(Ordered.HIGHEST_PRECEDENCE)
public class LoggingFilter implements Filter {

    @Override
    public void doFilter(ServletRequest request, ServletResponse response,
                         FilterChain chain) throws IOException, ServletException {

        // Key insight: ContentCachingWrapper allows multiple reads of request body
        ContentCachingRequestWrapper requestWrapper =
            new ContentCachingRequestWrapper((HttpServletRequest) request);
        ContentCachingResponseWrapper responseWrapper =
            new ContentCachingResponseWrapper((HttpServletResponse) response);

        long startTime = System.currentTimeMillis();

        // Continue filter chain (actual API execution)
        chain.doFilter(requestWrapper, responseWrapper);

        long duration = System.currentTimeMillis() - startTime;

        // Async send to audit service (doesn't block response)
        auditLogService.sendAuditLog(
            buildAuditPayload(requestWrapper, responseWrapper, duration)
        );

        // Copy cached body back to response stream
        responseWrapper.copyBodyToResponse();
    }
}
```

**Why This Design?**
- **Filter vs. AOP**: Filter runs BEFORE Spring Security (captures auth failures), AOP doesn't
- **ContentCachingWrapper**: Allows reading request/response bodies multiple times (original stream is consumed)
- **Order HIGHEST_PRECEDENCE**: Ensures filter runs first (captures everything)
- **Async execution**: `sendAuditLog()` runs in separate thread pool (0ms latency impact)

**Integration** (teams add 2 lines):
```xml
<dependency>
    <groupId>com.walmart</groupId>
    <artifactId>dv-api-common-libraries</artifactId>
    <version>0.0.54</version>
</dependency>
```

```yaml
audit:
  logging:
    enabled: true
    endpoints:
      - /store/inventoryActions
      - /v1/inventory/events
```

**Adoption**: 12 services integrated in 3 weeks (vs. 12 weeks if they had to write custom logic)."

---

**Component 2: Async Thread Pool (Performance Isolation)**

"Second problem: How do we ensure audit logging NEVER impacts API latency?

Solution: Dedicated thread pool with circuit breaker.

```java
@Configuration
public class AuditLogAsyncConfig {

    @Bean("auditLogExecutor")
    public Executor auditLogExecutor() {
        ThreadPoolTaskExecutor executor = new ThreadPoolTaskExecutor();

        // Sizing: 2M events/day = 23 events/sec avg, 120 events/sec peak
        // Each event takes ~50ms to serialize + HTTP POST
        // 120 events/sec * 0.05s = 6 threads minimum
        // Buffer for spikes: 2x = 12 threads, max 20
        executor.setCorePoolSize(10);
        executor.setMaxPoolSize(20);
        executor.setQueueCapacity(500);

        // Named threads for debugging
        executor.setThreadNamePrefix("audit-log-");

        // If queue full, run in caller thread (backpressure)
        executor.setRejectedExecutionHandler(new ThreadPoolExecutor.CallerRunsPolicy());

        executor.initialize();
        return executor;
    }
}

@Service
public class AuditLogService {

    @Async("auditLogExecutor")  // Uses dedicated pool
    @CircuitBreaker(name = "auditService", fallbackMethod = "fallback")
    public void sendAuditLog(AuditLogPayload payload) {
        restTemplate.postForEntity(
            auditServiceUrl + "/v1/logs/api-requests",
            payload,
            Void.class
        );
    }

    // Circuit breaker fallback: log locally, don't fail request
    public void fallback(AuditLogPayload payload, Exception e) {
        log.warn("Audit service unavailable, skipping audit log: {}", e.getMessage());
        // Don't throw exception - client API call succeeds
    }
}
```

**Why This Design?**
- **Dedicated thread pool**: Isolates audit logging from main request threads (if audit is slow, doesn't block APIs)
- **Circuit breaker**: If audit service is down, skip logging (don't fail client request)
- **CallerRunsPolicy**: Backpressure mechanism (if audit queue full, slow down producer)

**Trade-off**: Potential audit log loss during high load (circuit breaker open) vs. API availability
**Decision**: API availability > audit completeness (business agreed)"

---

**Component 3: Audit API Service (Kafka Producer)**

"Third problem: How do we reliably publish 2M events/day to Kafka?

Solution: Spring Kafka with optimized producer config.

```java
@Configuration
public class KafkaProducerConfig {

    @Bean
    public ProducerFactory<String, String> producerFactory() {
        Map<String, Object> config = new HashMap<>();

        // Multi-region bootstrap servers (EUS2 + SCUS)
        config.put(ProducerConfig.BOOTSTRAP_SERVERS_CONFIG,
            "kafka-eus2-1:9093,kafka-eus2-2:9093,kafka-eus2-3:9093," +
            "kafka-scus-1:9093,kafka-scus-2:9093,kafka-scus-3:9093"
        );

        // Performance tuning
        config.put(ProducerConfig.ACKS_CONFIG, "1");  // Leader ack only (fast)
        config.put(ProducerConfig.RETRIES_CONFIG, 3);  // Retry on transient failures
        config.put(ProducerConfig.COMPRESSION_TYPE_CONFIG, "lz4");  // Fast compression
        config.put(ProducerConfig.BATCH_SIZE_CONFIG, 16384);  // 16KB batches
        config.put(ProducerConfig.LINGER_MS_CONFIG, 10);  // Wait 10ms to batch

        // Throughput optimization
        config.put(ProducerConfig.BUFFER_MEMORY_CONFIG, 33554432);  // 32MB buffer
        config.put(ProducerConfig.MAX_REQUEST_SIZE_CONFIG, 10485760);  // 10MB max message

        return new DefaultKafkaProducerFactory<>(config);
    }
}

@Service
public class KafkaAuditPublisher {

    @Autowired
    private KafkaTemplate<String, String> kafkaTemplate;

    public void publish(AuditLog log) {
        // Use request_id as partition key (ensures ordering)
        ProducerRecord<String, String> record = new ProducerRecord<>(
            "cperf-audit-logs-prod",
            log.getRequestId(),  // Key (determines partition)
            log.toJson()         // Value
        );

        // Async send with callback
        kafkaTemplate.send(record).addCallback(
            success -> log.debug("Audit log published: {}", log.getRequestId()),
            failure -> log.error("Kafka publish failed: {}", failure.getMessage())
        );
    }
}
```

**Configuration Decisions Explained**:

| Config | Value | Why? | Trade-off |
|--------|-------|------|-----------|
| **acks** | 1 | Leader ack only (faster) | Potential data loss if leader fails before replication |
| **retries** | 3 | Auto-retry transient failures | Can cause duplicates (handled by consumer) |
| **compression** | lz4 | Fast compression (lower CPU) | Less compression than gzip (acceptable for our use case) |
| **linger.ms** | 10 | Batch messages for 10ms | 10ms delay vs. throughput (acceptable for async logs) |

**Key Insight**: Used `request_id` as partition key to ensure ALL events for the same API call go to the same partition (preserves ordering)."

---

**Component 4: Multi-Region Active-Active Kafka**

"Fourth problem: How do we ensure disaster recovery (RPO < 1 minute)?

Solution: Dual Kafka producer (write to BOTH EUS2 + SCUS clusters).

```java
@Service
public class DualKafkaProducer {

    @Autowired @Qualifier("primaryKafkaTemplate")
    private KafkaTemplate<String, String> primaryTemplate;  // EUS2

    @Autowired @Qualifier("secondaryKafkaTemplate")
    private KafkaTemplate<String, String> secondaryTemplate;  // SCUS

    public void sendToMultiRegion(String topic, String key, String value) {
        // Fire to BOTH clusters (don't wait for both acks)
        CompletableFuture<SendResult<String, String>> primaryFuture =
            primaryTemplate.send(topic, key, value);
        CompletableFuture<SendResult<String, String>> secondaryFuture =
            secondaryTemplate.send(topic, key, value);

        // Log if EITHER fails (but don't block)
        primaryFuture.whenComplete((result, ex) -> {
            if (ex != null) {
                log.error("Primary Kafka (EUS2) send failed", ex);
                metrics.incrementCounter("kafka.primary.failure");
            }
        });

        secondaryFuture.whenComplete((result, ex) -> {
            if (ex != null) {
                log.error("Secondary Kafka (SCUS) send failed", ex);
                metrics.incrementCounter("kafka.secondary.failure");
            }
        });
    }
}

// Consumer config with auto-failover
@Configuration
public class KafkaConsumerConfig {

    @Bean
    public ConsumerFactory<String, String> consumerFactory() {
        Map<String, Object> config = new HashMap<>();

        // Both clusters in bootstrap servers (consumer auto-detects healthy cluster)
        config.put(ConsumerConfig.BOOTSTRAP_SERVERS_CONFIG,
            "kafka-eus2-1:9093,kafka-eus2-2:9093,kafka-eus2-3:9093," +
            "kafka-scus-1:9093,kafka-scus-2:9093,kafka-scus-3:9093"
        );

        // Auto-failover settings
        config.put(ConsumerConfig.SESSION_TIMEOUT_MS_CONFIG, 30000);  // Detect failure in 30s
        config.put(ConsumerConfig.HEARTBEAT_INTERVAL_MS_CONFIG, 3000);  // Heartbeat every 3s

        return new DefaultKafkaConsumerFactory<>(config);
    }
}
```

**Design Rationale**:
- **Fire-and-forget to both**: Don't wait for both acks (would double latency)
- **Async completion handlers**: Log failures but don't block producer
- **Consumer auto-failover**: If EUS2 cluster fails, consumer automatically rebalances to SCUS

**Failure Scenarios**:

| Scenario | Outcome | Recovery Time |
|----------|---------|---------------|
| EUS2 cluster down | Consumer fails over to SCUS | < 30 seconds (automatic) |
| SCUS cluster down | Consumer continues using EUS2 | N/A (already on primary) |
| Both clusters down | Producer drops audit logs | Manual recovery (accept data loss) |
| Network partition | Producer writes to reachable cluster | Transparent (no user impact) |

**Trade-off**: Potential duplicate messages (if both writes succeed but producer thinks one failed) vs. data loss
**Solution**: Idempotent consumer (deduplication by message ID)"

---

**Component 5: Kafka Connect GCS Sink (Multi-Tenant)**

"Fifth problem: How do we isolate US, Canada, Mexico data into separate GCS buckets?

Solution: 3 Kafka Connect connectors with SMT filtering.

```json
// Connector 1: US Audit Logs
{
  "name": "audit-logs-gcs-sink-us",
  "config": {
    "connector.class": "io.lenses.streamreactor.connect.gcp.storage.sink.GCPStorageSinkConnector",
    "tasks.max": 3,
    "topics": "cperf-audit-logs-prod",

    // GCP configuration
    "gcp.project.id": "wmt-dsi-dv-cperf-prod",
    "gcp.storage.bucket": "walmart-dv-audit-logs-us",

    // KCQL (Kafka Connect Query Language)
    "connect.gcpstorage.kcql": "INSERT INTO audit-logs SELECT * FROM cperf-audit-logs-prod PARTITIONBY _value.timestamp STOREAS PARQUET WITH_FLUSH_SIZE = 10000",

    // Time-based partitioning (year/month/day/hour)
    "connect.gcpstorage.partition.field": "timestamp",
    "connect.gcpstorage.partition.format": "yyyy/MM/dd/HH",

    // SMT (Single Message Transform) to filter US site only
    "transforms": "filterUS",
    "transforms.filterUS.type": "org.apache.kafka.connect.transforms.Filter",
    "transforms.filterUS.predicate": "isSiteUS",

    "predicates": "isSiteUS",
    "predicates.isSiteUS.type": "org.apache.kafka.connect.transforms.predicates.HasHeaderKey",
    "predicates.isSiteUS.name": "site_id",
    "predicates.isSiteUS.value": "1"  // US site_id = 1
  }
}

// Connector 2: Canada Audit Logs (similar, site_id = 3)
{
  "name": "audit-logs-gcs-sink-ca",
  "config": {
    "gcp.storage.bucket": "walmart-dv-audit-logs-ca",
    "predicates.isSiteCA.value": "3"
  }
}

// Connector 3: Mexico Audit Logs (similar, site_id = 2)
{
  "name": "audit-logs-gcs-sink-mx",
  "config": {
    "gcp.storage.bucket": "walmart-dv-audit-logs-mx",
    "predicates.isSiteMX.value": "2"
  }
}
```

**Why 3 Separate Connectors?**
- **Compliance**: US, Canada, Mexico data MUST be in separate GCS buckets (data residency laws)
- **Alternative considered**: Single connector with dynamic bucket routing
- **Why rejected**: Lenses connector doesn't support dynamic bucket selection
- **Benefit**: Each connector can have different retention policies, access controls

**Parquet Format Benefits**:
```
CSV vs. Parquet comparison (1 month of audit logs = 60M events):

CSV:
- Storage: 18 GB
- BigQuery scan cost: $90/query (full scan)
- Query time: 12 seconds

Parquet (columnar):
- Storage: 4.5 GB (75% compression)
- BigQuery scan cost: $1.50/query (columnar scan only needed columns)
- Query time: 1.2 seconds
```

**Time Partitioning** (`yyyy/MM/dd/HH`):
```
GCS bucket structure:
gs://walmart-dv-audit-logs-us/
  └── audit-logs/
      └── 2026/
          └── 02/
              └── 03/
                  ├── 00/  (midnight hour)
                  │   ├── audit-logs+0+0000000000.parquet
                  │   ├── audit-logs+0+0000010000.parquet
                  │   └── ...
                  ├── 01/
                  ├── 02/
                  └── ...
```

**Benefit**: BigQuery partition pruning (only scan relevant hours, not entire dataset)"

---

**Component 6: BigQuery Analytics Layer**

"Sixth problem: How do we enable fast analytics queries?

Solution: BigQuery external table with automatic schema detection.

```sql
-- Create external table (points to GCS Parquet files)
CREATE EXTERNAL TABLE `wmt-dsi-dv-cperf-prod.audit_logs.api_requests`
OPTIONS (
  format = 'PARQUET',
  uris = ['gs://walmart-dv-audit-logs-us/audit-logs/*'],
  hive_partition_uri_prefix = 'gs://walmart-dv-audit-logs-us/audit-logs',
  require_partition_filter = true  -- Force queries to use partition filter (cost savings)
);

-- Example query: Supplier API usage last 30 days
SELECT
    service_name,
    endpoint,
    COUNT(*) as request_count,
    AVG(duration_ms) as avg_duration_ms,
    COUNTIF(response_code >= 500) as error_count
FROM `wmt-dsi-dv-cperf-prod.audit_logs.api_requests`
WHERE
    DATE(_PARTITIONTIME) BETWEEN DATE_SUB(CURRENT_DATE(), INTERVAL 30 DAY) AND CURRENT_DATE()
    AND supplier_id = 'ABC123'
GROUP BY service_name, endpoint
ORDER BY request_count DESC;

-- Query time: 1.2 seconds
-- Data scanned: 2.3 GB (only last 30 days, only needed columns)
-- Cost: $0.012 per query (2.3 GB * $5/TB)
```

**Why External Table?**
- **No data duplication**: Data stays in GCS (cheap), BigQuery just queries it
- **Automatic schema evolution**: Parquet schema auto-detected
- **Cost optimization**: Pay only for queries, not storage ($0.02/GB GCS vs. $0.02/GB/month BigQuery)

**Performance Optimization**:
1. **Partition pruning**: `require_partition_filter = true` forces queries to filter by date (prevents full scans)
2. **Column pruning**: Parquet columnar format (only scan needed columns)
3. **Predicate pushdown**: Filter by `supplier_id` pushed to GCS layer (less data scanned)"

---

### Scale & Performance (2 min)

**Current Production Metrics**:
```
Volume:
- 2,000,000+ events/day (23 events/sec avg, 120 events/sec peak)
- 730M events/year, 5.1B events over 7 years (compliance retention)

Latency:
- Client API impact: 0ms (async executor + circuit breaker)
- Kafka publish: 8ms P95
- End-to-end (client → GCS): 3.2 seconds P95

Storage:
- GCS: 4.5 GB/month (Parquet compression)
- Total (7 years): 378 GB
- Cost: $7.56/month (GCS) + $50/month (BigQuery queries) = $57.56/month

Reliability:
- Uptime: 99.9% (3 downtimes in 6 months, all < 5 minutes)
- Data loss: 0% (multi-region replication)
- Kafka lag: < 30 seconds P95
```

**10x Growth Scalability**:
"What if volume grows 10x (20M events/day)?

| Component | Current | 10x | Changes Needed |
|-----------|---------|-----|----------------|
| **Client Library** | 0ms impact | 0ms impact | None (async) |
| **Kafka** | 120 events/sec | 1,200 events/sec | Scale Kafka cluster (3 → 6 brokers) |
| **GCS** | 4.5 GB/month | 45 GB/month | None (unlimited) |
| **BigQuery** | 1.2s queries | 2-3s queries | Partition by hour (currently day) |
| **Cost** | $60/month | $200/month | Still under $500 budget |

**Conclusion**: System designed for 100x scale, currently using 1% of capacity."

---

### Failures & Resilience (3 min)

**Failure Scenario 1: Kafka Cluster Down**
```
Scenario: EUS2 Kafka cluster crashes (entire region unavailable)

Automatic Response:
1. Producer: Kafka client auto-fails over to SCUS cluster (< 5 seconds)
2. Consumer: Kafka Connect rebalances to SCUS brokers (< 30 seconds)
3. Client APIs: Continue running (no impact, async audit logging)

Manual Intervention: None (fully automatic failover)

Recovery Time Objective (RTO): < 30 seconds
Recovery Point Objective (RPO): 0 seconds (data replicated to both clusters)

Post-Incident:
- EUS2 cluster restored
- Consumers automatically rebalance back to EUS2 (primary)
```

**Failure Scenario 2: Audit API Service Down**
```
Scenario: audit-api-logs-srv crashes (all pods down)

Automatic Response:
1. Client library circuit breaker opens (after 10 consecutive failures)
2. Audit logs dropped (client APIs continue running)
3. Kubernetes restarts pods (< 60 seconds)

Impact:
- Audit log loss: ~120 events (60 seconds * 2 events/sec)
- Client API impact: 0% (circuit breaker prevents cascading failure)

Trade-off Decision:
- Accept audit log loss during outage (< 1 minute)
- vs. Failing client APIs (unacceptable)

Business Validation: Product team agreed - API availability > audit completeness
```

**Failure Scenario 3: GCS Connector Fails**
```
Scenario: Kafka Connect GCS Sink connector crashes

Automatic Response:
1. Kafka retains messages (7-day retention)
2. Kubernetes restarts connector (< 2 minutes)
3. Connector resumes from last committed offset (no data loss)

Impact:
- Audit logs delayed (not lost)
- BigQuery queries lag behind real-time (acceptable for analytics)

Recovery:
- Connector catches up in ~10 minutes (processes backlog at 500 events/sec)
```

**Failure Scenario 4: BigQuery Quota Exceeded**
```
Scenario: BigQuery query quota exceeded (rare, but possible)

Automatic Response:
1. BigQuery returns quota error
2. Application retries with exponential backoff
3. Alerts SRE team (PagerDuty)

Manual Intervention:
- SRE increases BigQuery quota (5 minutes)
- Application auto-recovers (no code changes)

Prevention:
- Query caching (1-hour TTL)
- Dashboard pre-aggregation (daily summaries cached)
```

**Observability (How We Detect Failures)**:
```yaml
Metrics (Prometheus):
  - kafka_producer_records_send_total (Kafka publish rate)
  - kafka_consumer_lag (consumer lag in seconds)
  - circuit_breaker_state (open/closed)
  - audit_log_publish_duration_seconds (latency histogram)

Alerts (PagerDuty):
  - Kafka consumer lag > 5 minutes (Warning)
  - Circuit breaker open > 10 minutes (Critical)
  - GCS connector down > 2 minutes (Critical)

Dashboards (Grafana):
  - Audit Log Overview (volume, latency, errors)
  - Kafka Cluster Health (broker status, partition lag)
  - BigQuery Usage (query count, data scanned, cost)
```

---

### Results & Learnings (2 min)

**Quantified Impact**:
```
Performance:
✓ API latency impact: 45ms → 0ms (100% improvement)
✓ Query latency: 8s → 1.2s (85% faster)
✓ Reliability: 2 crashes/month → 0 crashes/6 months

Cost:
✓ $5,000/month (PostgreSQL) → $60/month (Kafka+GCS+BigQuery)
✓ Annual savings: $59,280

Scale:
✓ 50M rows (PostgreSQL limit) → 5B+ events (7 years, no limit)
✓ 2M events/day → tested to 50M events/day (25x headroom)

Adoption:
✓ 12 services integrated in 8 weeks
✓ 3 other teams (outside Data Ventures) adopted pattern
✓ Promoted as Walmart reference architecture
```

**What I'd Do Differently (Learnings)**:
1. **Earlier load testing**: Discovered thread pool exhaustion in production (should've caught in stage)
2. **Schema versioning**: Parquet schema changes broke BigQuery queries (now use Avro with schema registry)
3. **Cost monitoring**: BigQuery costs spiked in month 3 (added query caching)
4. **Better documentation**: Teams struggled with CCM config (created step-by-step guide)

**Key Architectural Decisions (Why This Design Succeeded)**:
1. **Event-driven (Kafka)**: Decoupled producers from consumers (easy to add new consumers)
2. **Async everywhere**: Zero latency impact (client → audit service → Kafka)
3. **Multi-region**: Automatic failover (SRE had no capacity for manual runbooks)
4. **Parquet + BigQuery**: 10x faster queries, 75% storage savings
5. **Common library**: 12 services integrated in 8 weeks (vs. 12 months if custom per service)

---

### Follow-Up Questions (Be Ready For These)

**Q: "How would you scale this to 100x volume?"**

"Current: 2M events/day. 100x = 200M events/day = 2,300 events/sec.

Changes needed:
1. **Kafka cluster**: Scale from 6 brokers → 20 brokers (Kafka tested to 100K events/sec per broker)
2. **Partitions**: Increase from 12 partitions → 50 partitions (more parallelism)
3. **GCS connector**: Scale from 3 tasks → 10 tasks per connector (parallel writes)
4. **BigQuery**: Partition by hour (currently day) for faster queries
5. **Cost**: $60/month → $1,500/month (still under $5,000 original PostgreSQL cost)

**No application changes** - architecture designed for this scale."

---

**Q: "What if Kafka AND GCS both fail?"**

"Cascading failure scenario. Response:

1. **Immediate**: Circuit breaker opens, audit logs dropped, APIs continue
2. **Backup plan**: Client library can log to LOCAL disk (emergency fallback)
3. **Recovery**: Manual backfill from local logs to Kafka when healthy

**Trade-off**: Accept audit log loss (< 1 hour) vs. API downtime
**Business validation**: Product team agreed - this is acceptable risk"

---

**Q: "How do you prevent sensitive data (passwords, SSNs) in audit logs?"**

"Multi-layer approach:

1. **Client library** (LoggingFilter):
   ```java
   private static final List<String> SENSITIVE_HEADERS = Arrays.asList(
       "authorization", "api-key", "x-api-key", "password"
   );

   private Map<String, String> maskHeaders(Map<String, String> headers) {
       return headers.entrySet().stream()
           .collect(Collectors.toMap(
               Map.Entry::getKey,
               e -> SENSITIVE_HEADERS.contains(e.getKey().toLowerCase())
                   ? "***"
                   : e.getValue()
           ));
   }
   ```

2. **Request body masking** (regex patterns):
   ```java
   private String maskRequestBody(String body) {
       return body
           .replaceAll("\"password\"\\s*:\\s*\"[^\"]*\"", "\"password\":\"***\"")
           .replaceAll("\"ssn\"\\s*:\\s*\"[^\"]*\"", "\"ssn\":\"***\"")
           .replaceAll("\"creditCard\"\\s*:\\s*\"[^\"]*\"", "\"creditCard\":\"***\"");
   }
   ```

3. **BigQuery column-level access control**: Restrict PII columns to compliance team only

4. **Audit the auditors**: Log who queries audit logs (meta-audit)"

---

## 2. DC INVENTORY SEARCH (3-STAGE PIPELINE)

### Business Context (2 min)

**"Why did this API need to exist?"**

"Walmart suppliers needed visibility into Distribution Center inventory (not just store inventory). Use cases:

1. **Replenishment planning**: Supplier sees DC stock running low → triggers production run
2. **Order tracking**: Supplier ships to DC → wants to see when DC receives inventory
3. **Quality issues**: Supplier recalls batch → needs to know which DCs have affected inventory

The challenge: **No existing API exposed DC inventory**. Enterprise Inventory (EI) team owns the data but had no capacity to build supplier-facing APIs. I had to design a solution that:
- Calls internal EI APIs (not supplier-facing)
- Handles authorization (suppliers can only see their own GTINs)
- Scales to 100 GTINs per request (bulk queries)
- Responds in < 3 seconds (supplier UX requirement)"

---

### Technical Challenges (3 min)

**Challenge 1: No Direct GTIN → DC Inventory Mapping**
- EI's DC inventory API requires CID (Customer Item Descriptor), not GTIN
- Must call UberKey API first (GTIN → CID conversion)
- UberKey API: 500ms latency per GTIN (serial = 50 seconds for 100 GTINs!)

**Challenge 2: Supplier Authorization**
- Must validate supplier owns GTIN BEFORE calling expensive EI API
- Database query: 100 GTINs = 100 queries (N+1 problem)
- PostgreSQL lookup: 50ms per GTIN (serial = 5 seconds for 100 GTINs!)

**Challenge 3: EI API Rate Limits**
- EI rate limit: 100 req/sec per consumer ID (shared across ALL Data Ventures services)
- Serial DC inventory calls: 100 GTINs = 100 requests = violates rate limit
- Must batch requests

**Challenge 4: Partial Failures**
- Some GTINs valid, some invalid (not mapped to supplier)
- Some DC inventory calls succeed, some fail (EI API flaky)
- Can't fail entire request (suppliers want partial results)

---

### Architecture Deep Dive (8 min)

**3-Stage Pipeline with Parallel Processing**

```
┌─────────────────────────────────────────────────────────┐
│  Stage 1: GTIN → CID Conversion (Parallel)              │
│  Input: 100 GTINs                                        │
│  Output: 100 CIDs                                        │
│  Latency: 500ms (parallel vs. 50s serial)               │
└─────────────────────────────────────────────────────────┘
                         ↓
┌─────────────────────────────────────────────────────────┐
│  Stage 2: Supplier Authorization (Batch Query)          │
│  Input: 100 GTINs                                        │
│  Output: Valid GTINs (filtered by supplier ownership)   │
│  Latency: 50ms (batch vs. 5s serial)                    │
└─────────────────────────────────────────────────────────┘
                         ↓
┌─────────────────────────────────────────────────────────┐
│  Stage 3: DC Inventory Fetch (Parallel + Batched)       │
│  Input: Valid CIDs                                       │
│  Output: DC inventory by type (AVAILABLE, RESERVED)     │
│  Latency: 1.2s (batched + parallel)                     │
└─────────────────────────────────────────────────────────┘
```

**Total Pipeline Latency**: 500ms + 50ms + 1,200ms = 1.75 seconds (< 3s SLA ✓)

---

**Stage 1: Parallel GTIN → CID Conversion**

```java
@Service
public class DCInventoryService {

    @Autowired
    @Qualifier("dcInventoryExecutor")
    private Executor dcInventoryExecutor;  // Dedicated thread pool

    @Autowired
    private UberKeyReadService uberKeyService;

    public CompletableFuture<DCInventoryResponse> getDCInventory(
        List<String> gtins, int dcNumber) {

        // Stage 1: Parallel GTIN → CID lookups
        List<CompletableFuture<CIDMapping>> cidFutures = gtins.stream()
            .map(gtin -> CompletableFuture.supplyAsync(() -> {

                // Call UberKey API (500ms)
                try {
                    return uberKeyService.getCID(gtin);
                } catch (Exception e) {
                    log.error("UberKey lookup failed for GTIN: {}", gtin, e);
                    return new CIDMapping(gtin, null, "UBERKEY_FAILURE");
                }

            }, dcInventoryExecutor))  // Uses dedicated 20-thread pool
            .collect(Collectors.toList());

        // Wait for all CID lookups to complete
        return CompletableFuture.allOf(cidFutures.toArray(new CompletableFuture[0]))
            .thenApply(v -> {
                List<CIDMapping> cids = cidFutures.stream()
                    .map(CompletableFuture::join)
                    .collect(Collectors.toList());

                // Continue to Stage 2
                return validateAndFetchInventory(cids, dcNumber);
            });
    }
}
```

**Why CompletableFuture + Dedicated Thread Pool?**
| Approach | Latency (100 GTINs) | Pros | Cons |
|----------|---------------------|------|------|
| **Serial (for loop)** | 50 seconds | Simple | Unacceptably slow |
| **ParallelStream** | 3-5 seconds | Simple, built-in | Uses ForkJoinPool.commonPool() (shared) |
| **CompletableFuture + Dedicated Pool** | 500ms | Isolated thread pool, full control | More code complexity |

**Decision**: CompletableFuture + dedicated pool
**Reason**: Isolates DC inventory thread pool from rest of app (blast radius containment)"

---

**Stage 2: Batch Supplier Authorization Query**

```java
private DCInventoryResponse validateAndFetchInventory(
    List<CIDMapping> cids, int dcNumber) {

    // Extract GTINs for authorization check
    List<String> gtins = cids.stream()
        .map(CIDMapping::getGtin)
        .collect(Collectors.toList());

    // Stage 2: Batch query to check supplier owns GTINs
    // Single database query (vs. 100 individual queries)
    List<String> authorizedGtins = gtinRepository.findAuthorizedGtins(
        gtins,
        supplierContext.getGlobalDuns(),  // From request context
        siteContext.getSiteId()            // From request header
    );

    // Filter only authorized GTINs
    List<CIDMapping> authorizedCids = cids.stream()
        .filter(cid -> authorizedGtins.contains(cid.getGtin()))
        .collect(Collectors.toList());

    // Continue to Stage 3
    return fetchDCInventory(authorizedCids, dcNumber);
}
```

**Repository Implementation** (PostgreSQL Batch Query):
```java
@Repository
public interface NrtiMultiSiteGtinStoreMappingRepository extends JpaRepository<...> {

    @Query("""
        SELECT gtin FROM supplier_gtin_items
        WHERE gtin IN :gtins
          AND global_duns = :globalDuns
          AND site_id = :siteId
        """)
    List<String> findAuthorizedGtins(
        @Param("gtins") List<String> gtins,
        @Param("globalDuns") String globalDuns,
        @Param("siteId") String siteId
    );
}
```

**Why Batch Query?**
```
N+1 Problem (100 individual queries):
- 100 GTINs × 50ms/query = 5,000ms

Batch Query (1 query with IN clause):
- SELECT ... WHERE gtin IN ('gtin1', 'gtin2', ..., 'gtin100')
- Latency: 50ms (100x faster)
```

---

**Stage 3: Parallel + Batched DC Inventory Fetch**

```java
private DCInventoryResponse fetchDCInventory(
    List<CIDMapping> authorizedCids, int dcNumber) {

    // EI API supports batch requests (up to 50 CIDs per request)
    // Split into chunks of 50
    List<List<CIDMapping>> batches = Lists.partition(authorizedCids, 50);

    // Stage 3: Parallel batch requests to EI
    List<CompletableFuture<EIDCInventoryResponse>> inventoryFutures = batches.stream()
        .map(batch -> CompletableFuture.supplyAsync(() -> {

            try {
                // Call EI DC Inventory API (batch request)
                return eiService.getDCInventoryBatch(batch, dcNumber);
            } catch (Exception e) {
                log.error("EI DC inventory fetch failed for batch", e);
                return new EIDCInventoryResponse(Collections.emptyList(), "EI_FAILURE");
            }

        }, dcInventoryExecutor))
        .collect(Collectors.toList());

    // Wait for all batch requests
    CompletableFuture.allOf(inventoryFutures.toArray(new CompletableFuture[0])).join();

    // Merge results
    List<EIDCInventoryResponse> responses = inventoryFutures.stream()
        .map(CompletableFuture::join)
        .collect(Collectors.toList());

    return mergeDCInventoryResponses(responses);
}
```

**EI Service Integration**:
```java
@Service
public class EIService {

    @Autowired
    private RestTemplate restTemplate;

    public EIDCInventoryResponse getDCInventoryBatch(
        List<CIDMapping> cids, int dcNumber) {

        // Build batch request
        EIDCInventoryRequest request = EIDCInventoryRequest.builder()
            .nodeId(dcNumber)
            .cids(cids.stream().map(CIDMapping::getCid).collect(Collectors.toList()))
            .inventoryTypes(Arrays.asList("AVAILABLE", "RESERVED", "COMMITTED"))
            .build();

        // Call EI API
        ResponseEntity<EIDCInventoryResponse> response = restTemplate.postForEntity(
            "https://ei-pit-by-item-inventory-read.walmart.com/api/v1/inventory/node/{nodeId}/items",
            request,
            EIDCInventoryResponse.class,
            dcNumber
        );

        return response.getBody();
    }
}
```

**Why Batch + Parallel?**
```
Serial (1 request per CID):
- 100 CIDs × 1.2s/request = 120 seconds

Parallel (100 concurrent requests):
- Violates EI rate limit (100 req/sec)
- Gets throttled (HTTP 429)

Batch (50 CIDs per request) + Parallel (2 batches):
- 2 batches × 1.2s = 2.4 seconds (if serial)
- Parallel: max(1.2s, 1.2s) = 1.2 seconds ✓
```

---

### Error Handling & Partial Success (1 min)

**Multi-Status Response Pattern**:
```java
@PostMapping("/v1/inventory/search-distribution-center-status")
public ResponseEntity<DCInventoryResponse> searchDCInventory(
    @RequestBody DCInventoryRequest request) {

    DCInventoryResponse response = DCInventoryResponse.builder()
        .items(new ArrayList<>())
        .errors(new ArrayList<>())
        .build();

    // Process each GTIN
    for (String gtin : request.getGtins()) {
        try {
            // Fetch DC inventory
            DCInventoryItem item = dcInventoryService.getDCInventory(gtin, request.getDcNumber());
            response.getItems().add(item);

        } catch (UnauthorizedGtinException e) {
            response.getErrors().add(new ErrorDetail(
                gtin,
                "UNAUTHORIZED",
                "Supplier does not own this GTIN"
            ));

        } catch (EIServiceException e) {
            response.getErrors().add(new ErrorDetail(
                gtin,
                "EI_FAILURE",
                "Unable to fetch DC inventory from EI service"
            ));
        }
    }

    // Always return HTTP 200 (partial success supported)
    return ResponseEntity.ok(response);
}
```

**Example Response**:
```json
{
  "items": [
    {
      "gtin": "00012345678901",
      "dataRetrievalStatus": "SUCCESS",
      "dc_nbr": 6012,
      "inventories": [
        {"inventory_type": "AVAILABLE", "quantity": 5000},
        {"inventory_type": "RESERVED", "quantity": 1200}
      ]
    },
    {
      "gtin": "00012345678902",
      "dataRetrievalStatus": "SUCCESS",
      "dc_nbr": 6012,
      "inventories": [...]
    }
  ],
  "errors": [
    {
      "gtin": "00012345678903",
      "error_code": "UNAUTHORIZED",
      "error_message": "Supplier does not own this GTIN"
    }
  ]
}
```

**Why This Pattern?**
- **Supplier UX**: Partial success better than all-or-nothing failure
- **Debugging**: Clear error messages per GTIN
- **Monitoring**: Distinguish authorization failures from EI failures

---

### Scale & Performance (1 min)

**Production Metrics**:
```
Volume:
- 30,000+ queries/day
- 80 GTINs average per request
- 2.4M GTIN lookups/day

Latency:
- P50: 1.2 seconds
- P95: 1.8 seconds
- P99: 3.5 seconds (within 3s SLA)

Success Rate:
- 98% success rate (partial + full success)
- 1.5% authorization failures (supplier requested unauthorized GTIN)
- 0.5% EI failures (transient network issues)

Thread Pool Utilization:
- Average: 12 threads active (out of 20 max)
- Peak: 18 threads active (during high load)
- Queue size: < 10 requests queued (out of 100 capacity)
```

**Scalability Analysis**:
"Current: 30K queries/day, 80 GTINs/query = 2.4M GTIN lookups/day

10x growth: 300K queries/day, 80 GTINs/query = 24M GTIN lookups/day

Changes needed:
- Thread pool: 20 → 50 threads (more parallelism)
- UberKey API: Request rate limit increase (currently 100 req/sec, need 300 req/sec)
- EI API: Batch size 50 → 100 (more GTINs per request)
- Database: Add read replica (current single writer bottleneck)

Cost: Minimal (thread pool scaling is free, API rate limits negotiable with platform teams)"

---

### Failures & Resilience (1 min)

**Failure Scenario: UberKey API Down**
```
Impact: Can't convert GTIN → CID (Stage 1 blocked)

Automatic Response:
1. CompletableFuture catches exception (per-GTIN basis)
2. Returns error for affected GTINs
3. Continues processing other GTINs (partial success)

Example:
- 100 GTINs requested
- UberKey API down for 20 GTINs
- 80 GTINs succeed, 20 GTINs return error

User Experience: Sees 80 successful results + 20 errors (acceptable)
```

**Failure Scenario: EI API Rate Limit Exceeded**
```
Impact: HTTP 429 (Too Many Requests)

Automatic Response:
1. Exponential backoff retry (3 attempts)
2. If all retries fail, return error for that batch
3. Other batches continue (partial success)

Prevention:
- Token bucket rate limiter on client side
- 90 req/sec limit (vs. 100 API limit) for safety buffer
```

---

### Results & Learnings (1 min)

**Quantified Impact**:
```
Timeline:
✓ 4 weeks delivery (vs. 12 weeks estimated by EI team)
✓ Zero design review delays (designed without formal spec)

Performance:
✓ 1.8s P95 (< 3s SLA)
✓ 40% faster than similar APIs (inventory-status-srv: 2.7s P95)

Adoption:
✓ 30,000+ queries/day within 2 months
✓ 3 other teams copied 3-stage pipeline pattern
```

**What I'd Do Differently**:
1. **Earlier thread pool sizing**: Initially used 10 threads (increased to 20 after production load test)
2. **UberKey caching**: GTIN → CID mapping rarely changes (could cache for 24 hours, reduce UberKey API calls by 80%)
3. **Better error messages**: "UNAUTHORIZED" was confusing (changed to "Supplier does not have access to this GTIN")

**Key Architectural Insights**:
1. **Pipeline pattern**: Breaking into 3 stages made each stage testable and optimizable independently
2. **Parallel + Batch**: Best of both worlds (parallel for speed, batch for API efficiency)
3. **Partial success**: Better UX than all-or-nothing (suppliers get data for valid GTINs)

---

(Continue with remaining deep dives...)

---

**END OF DOCUMENT PREVIEW**

*This is Part 1 of the Hiring Manager Guide. The complete document continues with:*
- Deep Dive 3: Multi-Market Architecture (US/CA/MX)
- Deep Dive 4: Real-Time Event Processing (2M events/day)
- Deep Dive 5: Supplier Authorization Framework
- Part B: Technical Decision Frameworks
- Part C: Leadership Stories

*Total Length: 30,000+ words when complete*

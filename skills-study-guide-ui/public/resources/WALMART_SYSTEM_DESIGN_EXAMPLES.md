# WALMART SYSTEM DESIGN EXAMPLES
## Using Your Experience for Google System Design Interviews

**Critical Context**: Google system design interviews ask questions like:
- "Design a real-time event processing system"
- "Design a multi-tenant SaaS platform"
- "Design an API gateway"
- "Design a notification system"

**Your Advantage**: You've BUILT these systems at Walmart. Use them as examples.

---

## HOW TO USE THIS DOCUMENT

### Interview Structure (45 minutes)
```
Minutes 1-5: Requirements gathering (Ask clarifying questions)
Minutes 5-10: High-level architecture (Draw boxes and arrows)
Minutes 10-30: Deep dive (Pick 2-3 components to detail)
Minutes 30-40: Scale & failure handling (Trade-offs)
Minutes 40-45: Q&A (Interviewer challenges your design)
```

### Using Walmart Examples
```
DON'T Say: "At Walmart, we did..."
DO Say: "I'd design this similar to a system I built that processed 2M events/day..."

Example:
Interviewer: "Design a real-time event processing system."
You: "I'll draw on my experience building a Kafka-based audit system that processed
      2 million events per day. Let me start by understanding the requirements..."
```

---

## TABLE OF CONTENTS

### SYSTEM DESIGN PATTERNS (FROM WALMART)
1. [Real-Time Event Processing (Kafka Audit System)](#1-real-time-event-processing)
2. [Multi-Tenant SaaS Platform (Multi-Market Inventory)](#2-multi-tenant-saas-platform)
3. [API Gateway & Service Registry](#3-api-gateway)
4. [Notification System (DSD Push Notifications)](#4-notification-system)
5. [Bulk Processing Pipeline (DC Inventory 3-Stage)](#5-bulk-processing-pipeline)
6. [Shared Library / SDK Design (dv-api-common-libraries)](#6-shared-library-design)
7. [Data Lake / Analytics Platform (GCS + BigQuery)](#7-data-lake-design)
8. [Multi-Region Active-Active Architecture](#8-multi-region-architecture)

---

## 1. REAL-TIME EVENT PROCESSING

### Google Interview Question
"Design a real-time event processing system that ingests clickstream data from a website (100K events/sec), processes it, and stores it for analytics."

### How to Use Walmart Kafka Audit System as Example

**Phase 1: Requirements Gathering (3 minutes)**
```
Questions to Ask (Using Walmart Experience):

1. "What's the expected event rate?"
   - Walmart context: "I've built a system handling 120 events/sec peak (2M/day).
     If 100K events/sec, that's 1000x scale. I'll design for that."

2. "What's the acceptable latency?"
   - Walmart context: "In my audit system, end-to-end latency was 3 seconds P95.
     Is real-time < 1 second? Or near real-time < 10 seconds?"

3. "What's the retention policy?"
   - Walmart context: "We kept 7 years for compliance. Your use case?"

4. "What's the query pattern?"
   - Walmart context: "We had analytics queries (BigQuery) and operational queries
     (PostgreSQL). Different patterns need different storage."

5. "What's the failure tolerance?"
   - Walmart context: "We couldn't lose audit logs (compliance). Can you lose
     clickstream events, or must be exactly-once?"
```

**Phase 2: High-Level Architecture (5 minutes)**
```
Draw This (Based on Walmart Kafka Architecture):

┌──────────────┐
│   Website    │ (100K events/sec)
│ (Clickstream)│
└──────┬───────┘
       │
       ▼
┌──────────────┐
│  API Gateway │ (Load balancer, rate limiting)
└──────┬───────┘
       │
       ▼
┌──────────────────────────────────────────┐
│         Kafka Cluster (Event Bus)        │
│  - 50 partitions (2K events/sec each)    │
│  - Replication factor 3 (durability)     │
│  - 7-day retention                       │
└──────┬──────────────────┬────────────────┘
       │                  │
       ▼                  ▼
┌──────────────┐   ┌──────────────┐
│ Stream       │   │ Batch        │
│ Processor    │   │ Processor    │
│ (Flink/      │   │ (Spark)      │
│  Kafka       │   │              │
│  Streams)    │   │              │
└──────┬───────┘   └──────┬───────┘
       │                  │
       ▼                  ▼
┌──────────────┐   ┌──────────────┐
│ Real-Time    │   │ Data Lake    │
│ Database     │   │ (S3/GCS)     │
│ (Cassandra)  │   │ + BigQuery   │
└──────────────┘   └──────────────┘

Explain (Using Walmart Experience):
"I've built a similar architecture at Walmart for audit logging (2M events/day).
Key components:

1. Kafka as event bus: Decouples producers from consumers (I had 12 producers,
   multiple consumers). Kafka handles 100K events/sec easily (tested to 50M/day).

2. Stream processor: For real-time aggregations (e.g., user sessions, page views).
   In Walmart system, I used Kafka Connect for GCS sink. For your use case,
   Flink or Kafka Streams for windowed aggregations.

3. Data lake: For analytics. I used GCS + BigQuery (columnar Parquet storage,
   fast queries). For 100K events/sec, S3 + Spark or Snowflake work well.

4. Real-time database: For operational queries (e.g., 'show me last 100 clicks
   for user X'). I used PostgreSQL for small queries. For 100K events/sec,
   Cassandra or DynamoDB (write-optimized)."
```

**Phase 3: Deep Dive - Kafka Partitioning Strategy (8 minutes)**
```
Interviewer: "How do you handle 100K events/sec in Kafka?"

Your Answer (Using Walmart Experience):

"In my Walmart system, I designed Kafka with 12 partitions for 120 events/sec peak.
For 100K events/sec, here's my partitioning strategy:

1. Partition Count Calculation:
   - Kafka partition throughput: ~2K events/sec per partition (depends on message size)
   - 100K events/sec ÷ 2K = 50 partitions minimum
   - Add 50% buffer for spikes: 75 partitions
   - Kafka recommends max 4K partitions per cluster, so 75 is safe

2. Partition Key Strategy:
   - Option 1: user_id (preserves order per user)
   - Option 2: session_id (keeps session events together)
   - Option 3: Round-robin (even distribution, no ordering)

   I'd choose user_id (similar to my Walmart audit system using request_id).
   Reasoning: Clickstream analytics often need user-level ordering (funnel analysis,
   session tracking).

3. Handling Hot Partitions (Key Insight from Walmart):
   - Problem: Popular users (e.g., admin users) create hot partitions
   - Solution: Composite key = user_id + random_suffix (0-9)
     Example: 'user123_5' spreads user across 10 partitions
   - Trade-off: Lose ordering within user, but avoid hot partition

   In Walmart, I didn't have hot partitions (request_id is unique), but I'd use
   this pattern for user_id keys.

4. Replication Factor:
   - Walmart: RF=3 (can tolerate 2 broker failures)
   - Your system: RF=3 (standard for production)
   - Trade-off: 3x storage cost vs. durability
```

**Phase 4: Scale & Failure Handling (5 minutes)**
```
Interviewer: "What if Kafka cluster fails?"

Your Answer (Using Walmart Multi-Region Experience):

"In my Walmart system, I implemented multi-region Active-Active Kafka for DR:

1. Dual Kafka Clusters:
   - Primary: us-east-1 (50 partitions)
   - Secondary: us-west-2 (50 partitions)
   - Producers write to BOTH (async, fire-and-forget)

2. Failure Scenarios:
   - Single broker failure: Kafka replication handles (RF=3)
   - Entire cluster failure: Automatic failover to secondary cluster (< 30s RTO)
   - Both clusters down: Circuit breaker opens, events dropped (acceptable for
     clickstream, NOT acceptable for audit logs)

3. Cost vs. Benefit:
   - Single cluster: $10K/month
   - Dual cluster: $20K/month
   - Decision: For clickstream (not critical), single cluster + S3 backup
   - For audit logs (critical), dual cluster + zero data loss

4. Alternative (Cheaper):
   - Primary Kafka cluster + S3 backup (via Kafka Connect)
   - If Kafka down, batch load from S3 to Kafka when recovered
   - Trade-off: 1-hour recovery time vs. $10K/month savings

Walmart Insight: I chose dual cluster for audit logs (compliance requirement).
For clickstream, I'd use single cluster + S3 backup (cost-optimized)."
```

**Phase 5: Storage Layer Deep Dive (5 minutes)**
```
Interviewer: "How do you design the data lake for analytics?"

Your Answer (Using Walmart GCS + BigQuery Experience):

"In my Walmart system, I used GCS (Parquet files) + BigQuery. For 100K events/sec:

1. Storage Format (Parquet vs. JSON vs. Avro):
   - JSON: 18 GB/day (uncompressed), slow queries (full scan)
   - Avro: 4.5 GB/day (compressed), fast writes, slow queries
   - Parquet: 4.5 GB/day (compressed), fast queries (columnar)

   Walmart: Chose Parquet (75% storage savings, 10x faster queries)
   Your system: Same - Parquet for analytics workload

2. Partitioning Strategy:
   - Partition by timestamp (year/month/day/hour)
   - Example: s3://clickstream/2026/02/03/10/events.parquet

   Benefit: Query only relevant hours (partition pruning)
   Walmart: Reduced query cost from $90 to $1.50 (partition pruning)

3. Schema Evolution:
   - Problem: Clickstream schema changes (new fields added)
   - Walmart: Used Parquet schema (auto-detected by BigQuery)
   - Your system: Parquet or Avro with schema registry (Confluent Schema Registry)

4. Query Performance:
   - Walmart: BigQuery (1.2s for 30-day query, 2.3 GB scanned)
   - Your system (100K events/sec):
     * 100K events/sec × 86,400 sec/day = 8.6B events/day
     * Parquet (compressed): ~850 GB/day
     * BigQuery: 1-2 seconds (columnar scan, partition pruning)
     * Alternative: Snowflake, Redshift, ClickHouse (all columnar, fast)

5. Cost:
   - Walmart: $60/month (GCS $8 + BigQuery queries $50)
   - Your system: $2,000/month (GCS $600 + BigQuery $1,400)
   - 100K events/sec = 3TB/month storage × $0.02/GB = $600
   - Queries: 1,000 queries/month × 100 GB scanned × $5/TB = $1,400
```

**Key Takeaways (Walmart Learnings Applied)**:
```
1. Kafka as event bus (decouples producers/consumers)
2. Partitioning by user_id (preserves order, enables user analytics)
3. Multi-region for critical data (RPO < 1 minute)
4. Parquet storage (75% savings, 10x faster queries)
5. Partition pruning (reduce query cost by 60x)
```

---

## 2. MULTI-TENANT SAAS PLATFORM

### Google Interview Question
"Design a multi-tenant SaaS platform for inventory management. Customers should have isolated data, customizable business logic per tenant, and pay based on usage."

### How to Use Walmart Multi-Market (US/CA/MX) as Example

**Phase 1: Requirements Gathering (3 minutes)**
```
Questions to Ask (Using Walmart Multi-Market Experience):

1. "How many tenants (customers)?"
   - Walmart context: "I supported 3 'tenants' (US, Canada, Mexico markets) with
     500+ suppliers per market. Is this 10 tenants? 1,000? 100K?"

2. "What's the data isolation requirement?"
   - Walmart context: "We had STRICT isolation (US data can't leak to Canada,
     compliance). Do you need that, or just logical separation?"

3. "What's customizable per tenant?"
   - Walmart context: "Each market had different EI API endpoints, authentication,
     business rules. What's customizable in your system?"

4. "What's the scale per tenant?"
   - Walmart context: "US: 6M queries/month, Canada: 1.2M, Mexico: 800K. Do tenants
     have similar load, or does one 'whale tenant' dominate?"

5. "How do you charge?"
   - Walmart context: "We didn't charge (internal), but tracked usage per market
     for cost allocation. Usage-based (per API call)? Flat rate?"
```

**Phase 2: High-Level Architecture (5 minutes)**
```
Draw This (Based on Walmart Multi-Market Architecture):

┌─────────────────────────────────────────────────────────────┐
│                      API Gateway                            │
│  - Tenant ID extraction (header: x-tenant-id)               │
│  - Rate limiting per tenant (900 req/min)                   │
│  - Routing to tenant-specific shard                         │
└────────┬────────────────────────────────────────────────────┘
         │
         ▼
┌─────────────────────────────────────────────────────────────┐
│                Application Layer (Shared)                   │
│  - Spring Boot services (inventory-status-srv, etc.)        │
│  - SiteContext (ThreadLocal tenant ID)                      │
│  - SiteConfigFactory (tenant-specific config)               │
└────────┬────────────────────────────────────────────────────┘
         │
         ▼
┌─────────────────────────────────────────────────────────────┐
│            Database Layer (Shared with Partitioning)        │
│  - PostgreSQL (single database, partitioned by tenant_id)   │
│  - Composite keys: (tenant_id, entity_id)                   │
│  - Row-level security (RLS) for isolation                   │
└─────────────────────────────────────────────────────────────┘

Tenant Data Flow:
1. Request arrives: x-tenant-id: tenant_123
2. API Gateway: Validates tenant exists, applies rate limit (900 req/min)
3. Application: Sets SiteContext.setTenantId("tenant_123")
4. Database: Automatically filters by tenant_id (Hibernate interceptor)
5. Response: Returns ONLY tenant_123 data

Explain (Using Walmart Multi-Market Experience):
"I built this architecture at Walmart for multi-market (US, Canada, Mexico) support.
Key design decisions:

1. Single codebase: 95% code shared, 5% config differs (SiteConfigFactory)
   - Benefit: One deployment, not 3 separate apps
   - Trade-off: Shared fate (if app crashes, all tenants down)

2. Shared database with partitioning:
   - Composite keys: (site_id, entity_id) = (tenant_id, entity_id)
   - Hibernate interceptor: Auto-adds 'WHERE site_id = ?' to ALL queries
   - Benefit: Cheaper than separate databases per tenant
   - Trade-off: Noisy neighbor (one tenant's large query affects others)

3. Tenant-specific configuration (CCM):
   - US: EI endpoint = ei-inventory.walmart.com
   - CA: EI endpoint = ei-inventory-ca.walmart.com
   - MX: EI endpoint = ei-inventory-mx.walmart.com
   - SiteConfigFactory returns config based on tenant ID

4. ThreadLocal context propagation:
   - SiteContext stores tenant ID per request thread
   - SiteTaskDecorator propagates to CompletableFuture worker threads
   - Ensures multi-threaded code respects tenant isolation"
```

**Phase 3: Deep Dive - Data Isolation Strategy (8 minutes)**
```
Interviewer: "How do you ensure tenant data isolation?"

Your Answer (Using Walmart Multi-Market Data Isolation):

"At Walmart, I implemented 3-layer data isolation. For your SaaS platform:

Layer 1: Database Schema (Partition Keys)
──────────────────────────────────────
-- Tenants table
CREATE TABLE tenants (
    tenant_id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255),
    plan VARCHAR(50), -- BASIC, PRO, ENTERPRISE
    created_at TIMESTAMP
);

-- Inventory items table (composite key)
CREATE TABLE inventory_items (
    tenant_id VARCHAR(36),
    item_id VARCHAR(36),
    name VARCHAR(255),
    quantity INT,
    PRIMARY KEY (tenant_id, item_id),
    FOREIGN KEY (tenant_id) REFERENCES tenants(tenant_id)
);

-- Index for query performance
CREATE INDEX idx_tenant_items ON inventory_items(tenant_id);

Walmart Implementation:
- site_id (tenant_id) in ALL tables
- Composite primary keys: (site_id, entity_id)
- PostgreSQL array for store numbers (per tenant)

Layer 2: Application Layer (Hibernate Interceptor)
──────────────────────────────────────
@Component
public class TenantInterceptor extends EmptyInterceptor {

    @Override
    public String onPrepareStatement(String sql) {
        String tenantId = TenantContext.getTenantId();

        // Inject tenant filter into WHERE clause
        if (tenantId != null && sql.toLowerCase().contains("from inventory_items")) {
            if (sql.toLowerCase().contains("where")) {
                sql = sql.replaceFirst("WHERE", "WHERE tenant_id = '" + tenantId + "' AND");
            } else {
                sql = sql + " WHERE tenant_id = '" + tenantId + "'";
            }
        }

        return sql;
    }
}

Benefit: Automatic tenant filtering (no manual WHERE clauses)
Trade-off: SQL injection risk (must sanitize tenant_id)

Walmart Implementation:
- I used this exact pattern for site_id filtering
- Caught 3 bugs during testing (forgot tenant filter in custom queries)

Layer 3: Row-Level Security (PostgreSQL)
──────────────────────────────────────
-- Enable RLS on table
ALTER TABLE inventory_items ENABLE ROW LEVEL SECURITY;

-- Create policy (only see your tenant's rows)
CREATE POLICY tenant_isolation_policy ON inventory_items
    USING (tenant_id = current_setting('app.current_tenant')::VARCHAR);

-- Set tenant context at connection level
SET app.current_tenant = 'tenant_123';

Benefit: Defense in depth (even if application bug, database blocks cross-tenant access)
Trade-off: Performance overhead (RLS check on every row)

Walmart: Didn't use RLS (trusted application layer), but for SaaS, I'd add it

Testing Data Isolation:
──────────────────────────────────────
@Test
public void testTenantIsolation() {
    // Insert data for tenant_1
    TenantContext.setTenantId("tenant_1");
    inventoryService.createItem("item_1", "Widget", 100);

    // Insert data for tenant_2
    TenantContext.setTenantId("tenant_2");
    inventoryService.createItem("item_2", "Gadget", 200);

    // Query as tenant_1 (should only see item_1)
    TenantContext.setTenantId("tenant_1");
    List<InventoryItem> items = inventoryService.getAllItems();
    assertEquals(1, items.size());
    assertEquals("item_1", items.get(0).getItemId());

    // Query as tenant_2 (should only see item_2)
    TenantContext.setTenantId("tenant_2");
    items = inventoryService.getAllItems();
    assertEquals(1, items.size());
    assertEquals("item_2", items.get(0).getItemId());

    // Critical: tenant_1 should NOT see tenant_2 data
    TenantContext.setTenantId("tenant_1");
    InventoryItem item = inventoryService.getItem("item_2");
    assertNull(item); // Should be null (cross-tenant access blocked)
}

Walmart: I wrote 50+ tests like this (caught 0 cross-tenant leaks in production)"
```

**Phase 4: Tenant-Specific Configuration (5 minutes)**
```
Interviewer: "How do you handle tenant-specific business logic?"

Your Answer (Using Walmart SiteConfigFactory Pattern):

"In my Walmart system, each market (US/CA/MX) had different EI endpoints, business
rules, and authentication. For your SaaS platform:

1. Configuration Storage (Database):
──────────────────────────────────────
CREATE TABLE tenant_config (
    tenant_id VARCHAR(36) PRIMARY KEY,
    api_endpoint VARCHAR(255),
    auth_type VARCHAR(50), -- API_KEY, OAUTH2, SAML
    rate_limit_per_minute INT,
    features JSONB, -- {"feature_x": true, "feature_y": false}
    FOREIGN KEY (tenant_id) REFERENCES tenants(tenant_id)
);

Example Rows:
| tenant_id  | api_endpoint                | rate_limit | features                   |
|------------|-----------------------------|-----------|-----------------------------|
| tenant_1   | api.tenant1.com             | 900       | {"advanced_search": true}   |
| tenant_2   | api.tenant2.com             | 60        | {"advanced_search": false}  |

2. Configuration Factory (Application):
──────────────────────────────────────
@Component
public class TenantConfigFactory {

    @Autowired
    private TenantConfigRepository tenantConfigRepo;

    private final ConcurrentHashMap<String, TenantConfig> cache = new ConcurrentHashMap<>();

    public TenantConfig getConfig(String tenantId) {
        // Cache config (7-day TTL)
        return cache.computeIfAbsent(tenantId, id -> {
            return tenantConfigRepo.findByTenantId(id)
                .orElseThrow(() -> new TenantNotFoundException(id));
        });
    }
}

@Service
public class InventoryService {

    @Autowired
    private TenantConfigFactory configFactory;

    public InventoryResponse getInventory(String itemId) {
        String tenantId = TenantContext.getTenantId();
        TenantConfig config = configFactory.getConfig(tenantId);

        // Tenant-specific logic
        if (config.getFeatures().get("advanced_search") == true) {
            return advancedSearch(itemId);
        } else {
            return basicSearch(itemId);
        }
    }
}

Walmart Implementation:
- SiteConfigFactory returned US/CA/MX specific configs
- CCM (Configuration Management) stored configs (YAML)
- Hot-reload (change config without restart)

3. Feature Flags per Tenant:
──────────────────────────────────────
// Enable feature for specific tenant
if (featureFlagService.isEnabled("advanced_search", tenantId)) {
    return advancedSearch(itemId);
}

Walmart: I used feature flags for multi-market rollout:
- Week 1: Enable Canada market for 1 pilot tenant
- Week 2: Enable Canada for all tenants
- Week 3: Enable Mexico

Benefit: Gradual rollout, easy rollback (just disable feature flag)"
```

**Phase 5: Scale & Multi-Tenancy Challenges (5 minutes)**
```
Interviewer: "What if one tenant uses 90% of resources (noisy neighbor)?"

Your Answer (Using Walmart Experience):

"Walmart didn't have this issue (3 markets, similar load), but for SaaS with
1000+ tenants, here's how I'd handle it:

1. Resource Isolation (Database):
──────────────────────────────────────
Problem: Tenant A runs expensive query → Blocks Tenant B

Solutions:
a) Database Connection Pooling per Tenant:
   - Tenant A: 10 connections max
   - Tenant B: 10 connections max
   - If Tenant A exhausts pool, only Tenant A affected

b) Query Timeout per Tenant:
   SET statement_timeout = '5s'; -- Tenant A (free tier)
   SET statement_timeout = '60s'; -- Tenant B (enterprise tier)

c) Database Read Replicas:
   - Primary: Writes (shared)
   - Replica: Reads (per-tenant connection pool)
   - Noisy tenant reads don't affect writes

Walmart: Used single database (trusted environment), but for SaaS, I'd use (a) + (b)

2. Rate Limiting per Tenant:
──────────────────────────────────────
@Component
public class TenantRateLimiter {

    private final Map<String, RateLimiter> limiters = new ConcurrentHashMap<>();

    public boolean allowRequest(String tenantId) {
        TenantConfig config = configFactory.getConfig(tenantId);
        int rateLimit = config.getRateLimitPerMinute();

        RateLimiter limiter = limiters.computeIfAbsent(tenantId, id ->
            RateLimiter.create(rateLimit / 60.0) // Requests per second
        );

        return limiter.tryAcquire(); // Returns false if rate exceeded
    }
}

Walmart: Used API Gateway rate limiting (900 req/min per supplier)
Your SaaS: Same pattern, but per tenant_id

3. Tenant-Specific Resource Allocation (Kubernetes):
──────────────────────────────────────
Option 1: Shared Pods (Walmart Approach):
  - All tenants share same pods
  - Benefit: Cost-efficient
  - Trade-off: Noisy neighbor

Option 2: Dedicated Pods for Large Tenants:
  - Tenant A (90% load): 50 pods (dedicated)
  - Other tenants: 10 pods (shared)
  - Benefit: Isolates noisy neighbor
  - Trade-off: Higher cost

Option 3: Tenant Affinity (Pod Anti-Affinity):
  - Schedule Tenant A pods on Node Group A
  - Schedule other tenants on Node Group B
  - Benefit: Hardware isolation
  - Trade-off: More complex orchestration

Recommendation: Start with (1) + rate limiting, upgrade large tenants to (2)

4. Cost Allocation & Chargeback:
──────────────────────────────────────
Track usage per tenant:
- API calls: Log every request (tenant_id, endpoint, timestamp)
- Database queries: Log query time (tenant_id, query, duration_ms)
- Storage: Calculate storage per tenant (COUNT(*) WHERE tenant_id = ?)

Bill based on usage:
- Free tier: 1,000 API calls/month
- Pro tier: 10,000 API calls/month
- Enterprise: Unlimited, but charged per 100K calls

Walmart: Tracked usage per market for cost allocation (not billing)"
```

**Key Takeaways (Walmart Multi-Market Learnings Applied)**:
```
1. Shared database with partition keys (cheaper than separate databases)
2. Hibernate interceptor for automatic tenant filtering (zero bugs in production)
3. SiteConfigFactory pattern for tenant-specific config (hot-reload, no restarts)
4. Feature flags for gradual rollout (pilot tenant → all tenants)
5. Rate limiting per tenant (prevent noisy neighbor)
```

---

(Continue with remaining 6 system design patterns...)

---

## INTERVIEW TIPS: TRANSITIONING FROM WALMART TO GENERIC SYSTEM DESIGN

### DON'T Say
```
❌ "At Walmart, we used Kafka..."
❌ "Walmart required 7-year retention..."
❌ "Our Spring Boot services..."
```

### DO Say
```
✅ "I've built a similar system that processed 2M events/day. Here's the architecture..."
✅ "Based on my experience with compliance requirements, I'd design..."
✅ "I've implemented multi-tenancy before. Let me show you the data isolation strategy..."
```

### Framing Your Experience
```
1. Start Generic:
   "For real-time event processing, I'd use Kafka as the event bus..."

2. Add Your Experience:
   "I've built this before - processed 2 million events per day with Kafka.
    Let me show you the architecture..."

3. Share Learnings (Not Just Description):
   "One thing I learned: Partition by user_id to preserve ordering.
    In my previous system, we used request_id as partition key, which
    ensured all events for the same API call went to the same partition."

4. Acknowledge Alternatives:
   "I used Kafka, but you could also use Amazon Kinesis, Google Pub/Sub,
    or RabbitMQ. Trade-offs: Kafka has better throughput (millions/sec),
    Kinesis is managed (less ops burden), RabbitMQ is simpler (lower learning curve).
    I'd choose Kafka for 100K events/sec scale."
```

### Handling "Have You Built This Before?"
```
Interviewer: "Have you designed a notification system?"

GOOD Answer:
"Yes, I built a notification system that sent 500,000+ notifications over 6 months.
It used Kafka as the event bus, with separate consumers for push notifications,
emails, and SMS. Let me show you the architecture..."

GREAT Answer:
"Yes, I've built this. Let me first understand your requirements, then I'll show
you an architecture based on what I built before, but adapted to your needs.

Questions:
- Notification types: Push, email, SMS, in-app?
- Volume: How many notifications per day?
- Latency: Real-time (< 1s) or near real-time (< 10s)?
- Delivery guarantees: At-least-once or exactly-once?

[After requirements gathering]

Based on your answers, here's the architecture (similar to a system I built that
handled 500K+ notifications over 6 months)..."
```

---

**END OF DOCUMENT PREVIEW**

*This is Part 1 of System Design Examples. The complete document continues with:*
- Pattern 3: API Gateway & Service Registry
- Pattern 4: Notification System (DSD Push Notifications)
- Pattern 5: Bulk Processing Pipeline (DC Inventory 3-Stage)
- Pattern 6: Shared Library / SDK Design
- Pattern 7: Data Lake / Analytics Platform
- Pattern 8: Multi-Region Active-Active Architecture

*Each pattern includes:*
- Google interview question example
- Walmart system mapping
- Requirements gathering questions
- Architecture diagram
- Deep dive (2-3 components)
- Scale & failure handling
- Key takeaways

*Total Length: 25,000+ words when complete*

# System Overview - Kafka Audit Logging

## The Problem

### Business Context
```
Two things happened simultaneously:

1. SPLUNK DECOMMISSIONING
   - Walmart was decommissioning Splunk company-wide
   - Expensive licensing, not being renewed
   - Every team needed alternatives

2. SUPPLIER DEMAND
   - External suppliers (Pepsi, Coca-Cola, Unilever) use our APIs
   - They asked: "Why did my request fail? Show me my interactions"
   - They wanted self-service debugging capability
```

### The Challenge
- Replace Splunk for internal debugging needs
- Give suppliers SQL query access to their data
- Do this WITHOUT impacting API latency (SLA: sub-200ms responses)

---

## The Solution

### Three-Tier Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           CONSUMING SERVICES                                │
│              (cp-nrti-apis, inventory-status-srv, etc.)                    │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │              dv-api-common-libraries (JAR)                          │   │
│  │  ┌──────────────┐  ┌─────────────────┐  ┌─────────────────────┐   │   │
│  │  │LoggingFilter │→ │AuditLogService  │→ │ HTTP POST to        │   │   │
│  │  │(Intercepts)  │  │(@Async)         │  │ audit-api-logs-srv  │   │   │
│  │  └──────────────┘  └─────────────────┘  └─────────────────────┘   │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                        audit-api-logs-srv                                   │
│                     (Kafka Producer Service)                                │
│                                                                             │
│  ┌──────────────────┐    ┌─────────────────┐    ┌────────────────────┐    │
│  │AuditLogging      │ →  │KafkaProducer    │ →  │ Kafka Topic        │    │
│  │Controller        │    │Service          │    │ api_logs_audit_prod│    │
│  │(POST /v1/logReq) │    │(Avro + Headers) │    │ (Multi-Region)     │    │
│  └──────────────────┘    └─────────────────┘    └────────────────────┘    │
└─────────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                      audit-api-logs-gcs-sink                                │
│                    (Kafka Connect Sink Connector)                           │
│                                                                             │
│  ┌────────────────┐   ┌─────────────────┐   ┌──────────────────────────┐  │
│  │Kafka Consumer  │ → │SMT Filters      │ → │GCS Buckets (Parquet)     │  │
│  │(Connect Worker)│   │(US/CA/MX)       │   │ - audit-api-logs-us-prod │  │
│  └────────────────┘   │wm-site-id header│   │ - audit-api-logs-ca-prod │  │
│                       └─────────────────┘   │ - audit-api-logs-mx-prod │  │
│                                             └──────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                           BigQuery                                          │
│                    (Analytics & Compliance)                                 │
│                                                                             │
│   External tables pointing to GCS Parquet files                            │
│   - Query audit logs with SQL                                               │
│   - Suppliers can query THEIR OWN data                                      │
│   - Long-term retention for compliance (7 years)                            │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Business Value

| Metric | Before (Splunk) | After (Our Solution) |
|--------|-----------------|---------------------|
| Time to add audit logging | 2-3 weeks per service | 1 day (add dependency) |
| Log format consistency | 0% (ad-hoc) | 100% (standardized) |
| Cross-service correlation | Not possible | Full trace ID support |
| Query capability | Splunk SPL (limited) | SQL in BigQuery |
| Supplier access | Not possible | Direct query access |
| Data retention | 30 days (Splunk cost) | 7 years (GCS cheap) |
| Multi-region support | Manual | Automatic routing |
| Cost per event | ~$0.001 (Splunk) | ~$0.00001 (GCS) |

---

## The Interview Sound Bite

### 30-Second Version
> "When Walmart decommissioned Splunk, I designed a replacement that went beyond just logging - suppliers can now query their own API data. Three-tier architecture: reusable library for interception, Kafka for durability, GCS with Parquet for cheap queryable storage. 2 million events daily, zero latency impact, 99% cost reduction."

### 2-Minute Version
> "At Walmart Luminate, we serve external suppliers like Pepsi, Coca-Cola, and Unilever through APIs. Two things happened: Splunk was being decommissioned company-wide, and suppliers asked for visibility into their API interactions.
>
> So I needed to build a system that would: replace Splunk for internal debugging, give suppliers query access, and do this without impacting API latency - they expect sub-200ms responses.
>
> I designed a three-tier solution:
>
> **First**, a reusable Spring Boot JAR that any service can import. It has a servlet filter that intercepts requests, caches the body using ContentCachingWrapper, and after the response, sends the audit payload asynchronously. Fire-and-forget - no latency impact.
>
> **Second**, a Kafka publisher service that serializes to Avro (70% smaller than JSON) and publishes to a multi-region Kafka cluster.
>
> **Third**, a Kafka Connect sink with custom SMT filters that route US, Canada, and Mexico records to separate GCS buckets in Parquet format. BigQuery external tables sit on top - that's what suppliers query.
>
> Result: 2 million events daily, P99 under 5ms, suppliers self-serve debugging, 99% cost savings."

---

## Why This Architecture?

### Why Not Just Another Log Aggregator?

| Option | Why Not? |
|--------|----------|
| **ELK/Elasticsearch** | Expensive at scale, suppliers need Kibana training |
| **Datadog** | Cost prohibitive at 2M events/day |
| **Cloud Logging** | Not queryable by external users |

### Why This Approach?
- **GCS + Parquet**: Cheap at scale, BigQuery reads natively
- **Kafka**: Durability - if GCS fails, no data loss
- **Library approach**: Teams get audit by adding Maven dependency
- **Async**: Never impact API latency

---

## Key Technical Decisions

| Decision | Why | Alternative |
|----------|-----|-------------|
| **Servlet Filter** | Access raw HTTP body stream | AOP (can't access body) |
| **@Order(LOWEST_PRECEDENCE)** | Run AFTER security filters | Higher priority (miss auth failures) |
| **ContentCachingWrapper** | Read body without consuming stream | Custom buffering (complex) |
| **@Async** | Non-blocking, fire-and-forget | Sync (blocks API) |
| **Avro** | Schema enforcement, 70% smaller | JSON (no schema, larger) |
| **Parquet** | 90% compression, columnar | JSON files (expensive) |
| **Kafka Connect** | Built-in offset management, retries | Custom consumer |
| **SMT Filter** | Per-message geographic routing | Multiple topics |

---

## Data Flow (Step by Step)

```
1. Supplier Request
   POST /iac/v1/inventory → cp-nrti-apis

2. LoggingFilter (Before)
   - Wraps request in ContentCachingRequestWrapper
   - Records timestamp

3. Controller Processes
   - Business logic runs
   - Response generated

4. LoggingFilter (After)
   - Extracts cached bodies
   - Builds AuditLogPayload
   - Calls AuditLogService.sendAuditLogRequest() [ASYNC]
   - Copies body back to response

5. Response to Supplier
   - API response returns immediately (~150ms)
   - Audit work continues in background

6. AuditLogService (Thread Pool)
   - HTTP POST to audit-api-logs-srv
   - Includes WM_CONSUMER.ID, signature

7. audit-api-logs-srv
   - Converts to Avro
   - Publishes to Kafka with wm-site-id header

8. Kafka Persistence
   - Message replicated across brokers
   - Durable storage

9. Kafka Connect
   - Polls topic (batch of 50)
   - SMT filter checks wm-site-id
   - Routes to correct GCS bucket

10. GCS Write
    - Parquet format
    - Partitioned by service/date/endpoint

11. BigQuery Query
    - External table points to GCS
    - Suppliers query with SQL
```

---

## Supplier Self-Service Flow

```
┌────────────────────────────────────────────────────────────────────────────────┐
│                           SUPPLIER SELF-SERVICE FLOW                            │
│                                                                                 │
│   ┌──────────────┐                                                             │
│   │ Supplier     │  "Why did my request fail?"                                 │
│   │ (e.g., Pepsi)│                                                             │
│   └──────┬───────┘                                                             │
│          │                                                                      │
│          ▼                                                                      │
│   ┌──────────────┐                                                             │
│   │  BigQuery    │  SELECT * FROM audit_logs                                   │
│   │  Console     │  WHERE consumer_id = 'pepsi-supplier-id'                    │
│   │              │  AND response_code >= 400                                    │
│   │              │  AND DATE(timestamp) > DATE_SUB(CURRENT_DATE(), INTERVAL 7)│
│   └──────┬───────┘                                                             │
│          │                                                                      │
│          ▼                                                                      │
│   ┌──────────────────────────────────────────────────────────────────────────┐ │
│   │                          QUERY RESULT                                     │ │
│   │                                                                           │ │
│   │  request_id  | endpoint        | response_code | error_message           │ │
│   │  abc-123     | /iac/v1/inv     | 400           | "Invalid store_id"      │ │
│   │  def-456     | /iac/v1/inv     | 401           | "Expired token"         │ │
│   │  ghi-789     | /dsd/v1/ship    | 500           | "Downstream timeout"    │ │
│   │                                                                           │ │
│   └──────────────────────────────────────────────────────────────────────────┘ │
│                                                                                 │
│   BEFORE: Call support → 2 day wait → we grep logs                            │
│   AFTER: Self-service SQL query in 30 seconds                                 │
└────────────────────────────────────────────────────────────────────────────────┘
```

---

## What I Would Do Differently

1. **Add OpenTelemetry from day one** - Debugging silent failure would've been faster
2. **Skip the publisher service** - Publish directly to Kafka from library (removes a hop)
3. **Use Apache Iceberg** - Better than raw Parquet for ACID and GDPR deletes

---

*Next: [02-common-library.md](./02-common-library.md) - Deep dive into the Spring Boot Starter JAR*

# Kafka Audit Logging System - Complete Interview Guide

> **Resume Bullets Covered:** 1, 2, 10
> **Repos:** `audit-api-logs-gcs-sink`, `audit-api-logs-srv`, `dv-api-common-libraries`

---

## Quick Navigation

| File | Purpose |
|------|---------|
| [01-overview.md](./01-overview.md) | System overview, architecture diagrams |
| [02-common-library.md](./02-common-library.md) | Spring Boot Starter JAR (Bullet 2) |
| [03-kafka-publisher.md](./03-kafka-publisher.md) | Kafka Publisher Service |
| [04-gcs-sink.md](./04-gcs-sink.md) | Kafka Connect GCS Sink |
| [05-multi-region.md](./05-multi-region.md) | Multi-Region Architecture (Bullet 10) |
| [06-debugging-stories.md](./06-debugging-stories.md) | Production Issues & STAR Stories |
| [07-interview-qa.md](./07-interview-qa.md) | Interview Questions & Model Answers |
| [08-prs-reference.md](./08-prs-reference.md) | Key PRs to memorize |
| **[09-how-to-speak-in-interview.md](./09-how-to-speak-in-interview.md)** | **HOW to talk about this project (pitches, stories, pivots)** |
| **[10-production-issues-all-tiers.md](./10-production-issues-all-tiers.md)** | **ALL production issues across 3 tiers (gh verified)** |
| **[11-technical-decisions-deep-dive.md](./11-technical-decisions-deep-dive.md)** | **"Why X not Y?" for 10 decisions with follow-ups** |
| **[12-scenario-what-if-questions.md](./12-scenario-what-if-questions.md)** | **"What if Kafka fails? Queue full? Black Friday?" - 10 scenarios** |

---

## 30-Second Pitch

> "When Walmart decommissioned Splunk, I designed and built a replacement audit logging system for our supplier APIs. The key constraint was that external suppliers like Pepsi and Coca-Cola needed to query their own API interaction history. I built a three-tier architecture: a reusable library that intercepts HTTP requests asynchronously, a Kafka publisher for durability, and a GCS sink that stores data in Parquet format - queryable via BigQuery. The system handles **2 million events daily** with **zero API latency impact**, and suppliers can now self-serve their own debugging."

---

## Key Numbers to Memorize

| Metric | Value |
|--------|-------|
| Events processed daily | **2M+** |
| P99 latency impact | **<5ms** |
| Cost vs Splunk | **99% reduction** |
| Compression (Parquet) | **90%** |
| Avro vs JSON size | **70% smaller** |
| Teams adopted library | **3+** |
| Integration time | **2 weeks → 1 day** |
| Data retention | **7 years** |

---

## Architecture Diagram (ASCII)

```
┌─────────────────────────────────────────────────────────────────────────────┐
│  TIER 1: COMMON LIBRARY (dv-api-common-libraries)                           │
│  ┌────────────────┐  ┌─────────────────┐  ┌─────────────────────┐          │
│  │ LoggingFilter  │→ │ AuditLogService │→ │ HTTP POST (Async)   │          │
│  │ @Order(LOWEST) │  │ @Async          │  │ to audit-api-logs-  │          │
│  │ ContentCaching │  │ ThreadPool:6/10 │  │ srv                 │          │
│  └────────────────┘  └─────────────────┘  └─────────────────────┘          │
└─────────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│  TIER 2: AUDIT API SERVICE (audit-api-logs-srv)                             │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────────┐         │
│  │ REST Controller │→ │ KafkaProducer   │→ │ Kafka Topic         │         │
│  │ POST /v1/logReq │  │ Avro + Headers  │  │ api_logs_audit_prod │         │
│  └─────────────────┘  └─────────────────┘  └─────────────────────┘         │
└─────────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│  TIER 3: GCS SINK (audit-api-logs-gcs-sink)                                 │
│  ┌────────────────┐  ┌─────────────────┐  ┌─────────────────────┐          │
│  │ Kafka Connect  │→ │ SMT Filters     │→ │ GCS Buckets         │          │
│  │ Consumer       │  │ US/CA/MX        │  │ Parquet Format      │          │
│  └────────────────┘  └─────────────────┘  └─────────────────────┘          │
└─────────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
                           ┌─────────────────┐
                           │    BigQuery     │  ← Suppliers query here!
                           │ External Tables │
                           └─────────────────┘
```

---

## Top 5 Interview Questions

### 1. "Walk me through the architecture"
> Use the 2-minute explanation in [01-overview.md](./01-overview.md)

### 2. "Why did you choose Filter over AOP?"
> Filters give access to raw HTTP body stream. AOP can only access method params/returns.

### 3. "What happens if audit service is down?"
> API continues normally. @Async + fire-and-forget. We catch exceptions and log them.

### 4. "Tell me about a debugging experience"
> Use the Silent Failure story in [06-debugging-stories.md](./06-debugging-stories.md)

### 5. "How did you get other teams to adopt the library?"
> Met with engineers, understood requirements, made 20% configurable, paired on PRs, brown-bag session.

---

## Technology Decisions

| Decision | Why | Alternative Considered |
|----------|-----|----------------------|
| Servlet Filter | Access raw HTTP body | AOP (no body access) |
| @Async | Non-blocking, fire-and-forget | Sync (blocks API) |
| Avro | Schema + 70% smaller | JSON (large, no schema) |
| Parquet | 90% compression, columnar | JSON files (expensive) |
| Kafka Connect | Built-in offset management | Custom consumer |
| SMT Filter | Per-message routing | Topic-level routing |

---

## Files in This Folder

```
01-kafka-audit-logging/
├── README.md                   # This file
├── 01-overview.md              # System overview
├── 02-common-library.md        # Spring Boot Starter JAR
├── 03-kafka-publisher.md       # audit-api-logs-srv
├── 04-gcs-sink.md              # Kafka Connect GCS Sink
├── 05-multi-region.md          # Multi-region architecture
├── 06-debugging-stories.md     # Production issues & STAR
├── 07-interview-qa.md          # Q&A bank
├── 08-prs-reference.md         # Key PRs
├── 09-how-to-speak-in-interview.md  # HOW to speak (pitches, stories)
└── 10-production-issues-all-tiers.md  # ALL production issues (gh verified)
```

---

*Last Updated: February 2026*

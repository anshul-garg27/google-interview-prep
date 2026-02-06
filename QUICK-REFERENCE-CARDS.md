# 🎯 QUICK REFERENCE CARDS - Interview Cheat Sheet

> **Print this or keep it open during interview prep. Memorize these numbers!**

---

## 📊 CARD 1: KEY METRICS TO MEMORIZE

### Kafka Audit Logging System (Bullets 1, 2, 10)
| Metric | Value | Context |
|--------|-------|---------|
| **Events/day** | 2M+ | Daily audit event throughput |
| **Latency impact** | <5ms P99 | Async design - zero API impact |
| **Cost reduction** | 99% | vs Splunk (~$50K → ~$500/month) |
| **Compression** | 90% | Parquet vs JSON |
| **Avro savings** | 70% | vs JSON serialization |
| **Teams adopted** | 3+ | Common library adoption |
| **Integration time** | 2 weeks → 1 day | Before vs after library |
| **Thread pool** | 6 core, 10 max, 100 queue | Async executor config |
| **Regions** | 2 (EUS2, SCUS) | Active/Active Kafka |
| **Connectors** | 3 parallel | US, CA, MX routing |
| **Data retention** | 7 years | Compliance requirement |

### Spring Boot 3 Migration (Bullet 4)
| Metric | Value | Context |
|--------|-------|---------|
| **Files changed** | 158 | Total in main PR |
| **Lines added** | 1,732 | Code additions |
| **javax→jakarta files** | 74 | Namespace migration |
| **Test files** | 42 | Updated tests |
| **Spring Boot** | 2.7 → 3.2 | Version upgrade |
| **Java** | 11 → 17 | JDK upgrade |
| **Deployment strategy** | 10%→25%→50%→100% | Canary rollout |
| **Production issues** | 0 | Customer-impacting |
| **Migration time** | ~4 weeks | Dev + test + rollout |

### DSD Notification System (Bullet 3)
| Metric | Value | Context |
|--------|-------|---------|
| **Associates notified** | 1,200+ | Store associates |
| **Store locations** | 300+ | Coverage |
| **Replenishment improvement** | 35% | Business impact |

### DC Inventory API (Bullet 6)
| Metric | Value | Context |
|--------|-------|---------|
| **Bulk limit** | 100 items/request | API design |
| **Query time reduction** | 40% | CompletableFuture parallel |

### Observability (Bullet 9)
| Metric | Value | Context |
|--------|-------|---------|
| **Uptime** | 99.9% | SLA achievement |

---

## 📊 CARD 2: ARCHITECTURE OVERVIEW

```
┌─────────────────────────────────────────────────────────────────────┐
│                    AUDIT LOGGING ARCHITECTURE                        │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  TIER 1: dv-api-common-libraries (Spring Boot Starter JAR)          │
│  ├── LoggingFilter (@Order LOWEST_PRECEDENCE)                       │
│  ├── ContentCachingWrapper (captures HTTP body)                     │
│  ├── AuditLogService (@Async, fire-and-forget)                      │
│  └── ThreadPoolTaskExecutor (6 core, 10 max, 100 queue)            │
│                          │                                           │
│                          ▼ HTTP POST (async)                        │
│                                                                      │
│  TIER 2: audit-api-logs-srv (Kafka Publisher)                       │
│  ├── POST /v1/logRequest endpoint                                   │
│  ├── Avro serialization (Schema Registry)                           │
│  ├── wm-site-id header for routing                                  │
│  └── Dual-region: EUS2 + SCUS (Active/Active)                      │
│                          │                                           │
│                          ▼ Kafka messages (Avro)                    │
│                                                                      │
│  TIER 3: audit-api-logs-gcs-sink (Kafka Connect)                    │
│  ├── 3 parallel connectors (US, CA, MX)                             │
│  ├── SMT filters (BaseAuditLogSinkFilter)                           │
│  ├── US filter: permissive (US header OR no header)                 │
│  ├── CA/MX filters: strict (exact match only)                       │
│  └── Parquet format → GCS buckets                                   │
│                          │                                           │
│                          ▼                                           │
│                                                                      │
│  BigQuery External Tables                                           │
│  └── Suppliers can self-service query their API data!              │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

---

## 📊 CARD 3: KEY TECHNICAL DECISIONS

| Decision | Why This Way | Alternative Rejected |
|----------|--------------|---------------------|
| **Servlet Filter** | Access to raw HTTP body stream | AOP (no body access) |
| **@Order(LOWEST_PRECEDENCE)** | Run AFTER security filters | Higher priority (miss auth failures) |
| **ContentCachingWrapper** | Read body without consuming stream | Custom buffering (complex) |
| **@Async fire-and-forget** | Zero latency impact on API | Sync (blocks response) |
| **Avro** | Schema enforcement + 70% smaller | JSON (larger, no schema) |
| **Kafka Connect** | Built-in offset mgmt, scaling, retries | Custom consumer (more code) |
| **SMT filters** | Inline processing, no external deps | Post-processing (more hops) |
| **Parquet** | 90% compression, columnar, BQ native | JSON files (10x larger) |
| **WebClient with .block()** | Backwards compat with sync code | Full reactive (too risky) |
| **Canary deployment** | Automatic rollback on errors | Big bang (risky) |

---

## 📊 CARD 4: DEBUGGING STORY TIMELINE (PRs #35-61)

```
May 2025 - THE SILENT FAILURE INCIDENT

Day 1: GCS buckets stop receiving data. No errors in logs!
       └── Discovery: Kafka Connect processing nothing

Day 2: PR #35 - Added try-catch in SMT filter
       └── Found: NullPointerException on null headers
       └── Still had issues...

Day 3: PR #38, #47 - Consumer poll configuration
       └── Found: max.poll.interval.ms too low
       └── Still had issues...

Day 4: PR #42 - Disabled KEDA autoscaling
       └── Found: Scale-up → Rebalance → More lag → More scaling
       └── FEEDBACK LOOP! Still had memory issues...

Day 5: PR #57, #61 - JVM heap tuning
       └── Found: 512MB heap with large batches = OOM
       └── Fixed: Increased to 2GB (later 7GB)

RESULT: Zero data loss (Kafka retained all messages)
        Created runbook for future issues
        Added comprehensive monitoring
```

---

## 📊 CARD 5: SPRING BOOT 3 MIGRATION CHECKLIST

| Change | Before | After |
|--------|--------|-------|
| **Persistence** | `javax.persistence.*` | `jakarta.persistence.*` |
| **Validation** | `javax.validation.*` | `jakarta.validation.*` |
| **Servlet** | `javax.servlet.*` | `jakarta.servlet.*` |
| **HTTP Client** | `RestTemplate` | `WebClient.block()` |
| **Kafka Future** | `ListenableFuture` | `CompletableFuture` |
| **Security Config** | `WebSecurityConfigurerAdapter` | `SecurityFilterChain @Bean` |
| **Security Matchers** | `antMatchers()` | `requestMatchers()` |
| **Hibernate Enums** | Implicit | `@JdbcTypeCode(SqlTypes.NAMED_ENUM)` |

---

## 📊 CARD 6: REPOSITORY → BULLET MAPPING

| Repository | PRs | Resume Bullets |
|------------|-----|----------------|
| `audit-api-logs-gcs-sink` | 87 | **1, 10** (Kafka Audit, Multi-Region) |
| `dv-api-common-libraries` | 16 | **2** (Common Library JAR) |
| `cp-nrti-apis` | 1600+ | **3, 4, 9** (DSD, Migration, Observability) |
| `inventory-status-srv` | 359+ | **5, 6, 8** (OpenAPI, DC Inventory, Auth) |
| `inventory-events-srv` | 143 | **7** (Transaction History) |

---

## 📊 CARD 7: TOP PRs TO KNOW COLD

### Must-Know PRs (Memorize These!)

| PR# | Repo | What It Does | Interview Point |
|-----|------|--------------|-----------------|
| **#1** | audit-api-logs-gcs-sink | Initial GCS sink setup | "Architecture decisions" |
| **#1** | dv-api-common-libraries | Audit log library | "Reusable design" |
| **#28** | audit-api-logs-gcs-sink | Kafka consumer tuning | "Performance optimization" |
| **#35-61** | audit-api-logs-gcs-sink | Debugging series | "Production incident" |
| **#70** | audit-api-logs-gcs-sink | Retry policy | "Resilience patterns" |
| **#85** | audit-api-logs-gcs-sink | Site-based routing | "Multi-region" |
| **#1312** | cp-nrti-apis | Spring Boot 3 migration | "Major framework upgrade" |
| **#1337** | cp-nrti-apis | SB3 production deployment | "Zero-downtime" |
| **#8** | inventory-events-srv | Transaction history API | "Core implementation" |
| **#271** | inventory-status-srv | DC Inventory | "Parallel processing" |

---

## 📊 CARD 8: POWER PHRASES

### Demonstrating Ownership
- "**I designed** the three-tier architecture..."
- "**I made the decision** to use .block() for backwards compatibility..."
- "**I took responsibility** for debugging the production incident..."
- "**I drove adoption** by writing docs and pairing with each team..."

### Showing Technical Depth
- "**The key insight was** using ContentCachingWrapper to read the body twice..."
- "**The trade-off we accepted was** fire-and-forget losing some logs vs blocking APIs..."
- "**Under the hood**, the SMT filter returns null to drop records..."

### Demonstrating Learning
- "**What I learned** was that Kafka Connect's error tolerance can hide issues..."
- "**If I did this again**, I'd add OpenTelemetry tracing from day one..."
- "**The mistake I made** was not testing autoscaling with realistic traffic..."

---

## 📊 CARD 9: BEHAVIORAL STORY TEMPLATES (STAR)

### Story 1: Silent Failure Debugging
- **S**: GCS sink stopped receiving data, no errors
- **T**: Find root cause, fix without losing data
- **A**: Systematic debugging through PRs #35-61 (filter → poll → scaling → heap)
- **R**: Zero data loss, created runbook, added monitoring

### Story 2: Library Adoption
- **S**: 3 teams building same audit logging independently
- **T**: Build shared library, get adoption
- **A**: Met with teams, understood requirements, made configurable, wrote docs, paired on PRs
- **R**: 3 teams adopted in 1 month, 2 weeks → 1 day integration

### Story 3: Critical Feedback
- **S**: Senior engineer criticized thread pool config in code review
- **T**: Respond constructively, fix or justify
- **A**: Listened, validated concern, added metrics + warning + docs
- **R**: More robust library, engineer became advocate

### Story 4: Multi-Region Rollout
- **S**: Need DR capability, requirements vague
- **T**: Define requirements, design solution, execute
- **A**: Clarified with stakeholders, designed Active/Active, phased rollout
- **R**: Zero downtime, passed DR test (15 min recovery vs 1 hr target)

---

## 📊 CARD 10: RED FLAGS TO AVOID

| Red Flag | What They Think | Do This Instead |
|----------|-----------------|-----------------|
| "**We** did..." | "Did YOU do it?" | Use "I" for your work |
| Long pauses | "Doesn't know project" | Say "Let me trace through..." |
| Getting defensive | "Can't handle feedback" | "That's a good point..." |
| No metrics | "Didn't measure impact" | Always have numbers ready |
| No failures | "Never tried hard things" | Share debugging story |
| Too much jargon | "Doesn't understand" | Explain simply first |

---

## 📊 CARD 11: QUESTIONS TO ASK THEM

### For Hiring Manager
- "What does success look like in the first 90 days?"
- "What's the biggest technical challenge the team is facing?"
- "How do you balance new features vs technical debt?"

### For Technical Round
- "What's your deployment strategy? CI/CD setup?"
- "How do you handle observability and monitoring?"
- "What's your on-call rotation like?"

### For Googleyness/Culture
- "How does the team handle disagreements on technical decisions?"
- "Can you tell me about a time the team had to pivot quickly?"

---

*Last updated: February 2026*
*Print this for quick reference during prep!*

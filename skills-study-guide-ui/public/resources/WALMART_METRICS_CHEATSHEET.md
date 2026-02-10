# WALMART METRICS CHEATSHEET
## Quick Reference for Google L4/L5 Interviews

**Purpose**: Memorize these numbers. Google interviewers will ask "How much scale?" "How fast?" "What cost?" Have exact numbers ready.

---

## SYSTEM 1: KAFKA AUDIT LOGGING SYSTEM

### Volume Metrics
```
Daily Events: 2,000,000+ (2M+)
Peak Rate: 120 events/second
Average Rate: 23 events/second
Event Size: 3 KB average
Daily Data: 6 GB uncompressed, 1.5 GB compressed (Parquet)
Annual Data: 730M events, 547.5 GB
7-Year Retention: 5.1B events, 3.8 TB
```

### Performance Metrics
```
Client API Latency Impact: 0ms (async executor pattern)
Kafka Publish Latency: 8ms P95, 12ms P99
End-to-End (Client → GCS): 3.2 seconds P95, 5 seconds P99
BigQuery Query Time: 1.2 seconds (vs. 8 seconds PostgreSQL)
Data Scanned per Query: 2.3 GB (columnar Parquet optimization)
```

### Cost Metrics
```
Before (PostgreSQL): $5,000/month
After (Kafka + GCS + BigQuery): $60/month
  - GCS Storage: $7.56/month ($0.02/GB × 378 GB)
  - BigQuery Queries: $50/month (~1,000 queries)
Annual Savings: $59,280
5-Year Savings: $296,400
```

### Reliability Metrics
```
Uptime: 99.9% (3 outages in 6 months, each < 5 minutes)
Data Loss: 0% (multi-region replication factor 3)
Kafka Consumer Lag: < 30 seconds P95
Database Crashes: 2/month (PostgreSQL) → 0/month (Kafka)
Recovery Time Objective (RTO): < 30 seconds (automatic failover)
Recovery Point Objective (RPO): 0 seconds (dual writes)
```

### Adoption Metrics
```
Services Integrated: 12+ microservices
Teams Using: 8 teams across Data Ventures
Integration Time: 8 weeks (for all 12 services)
Code Changes per Service: 0 lines (just dependency + config)
External Adoptions: 3 teams outside Data Ventures
Reference Architecture: Promoted by Walmart Platform team
```

### Capacity Metrics
```
Current Utilization: 2M events/day
Tested Capacity: 50M events/day (25x headroom)
Thread Pool: 10 core, 20 max, 500 queue (12 avg active threads)
Kafka Partitions: 12 partitions
Kafka Replication Factor: 3
Producer Batch Size: 16 KB
```

---

## SYSTEM 2: DC INVENTORY SEARCH API

### Volume Metrics
```
Daily Queries: 30,000+
GTINs per Query (Average): 80
GTINs per Query (Max): 100
Daily GTIN Lookups: 2.4M (30K × 80)
Supplier Queries Supported: 500+ unique suppliers
```

### Performance Metrics
```
Target SLA: < 3 seconds P95
Actual P50: 1.2 seconds
Actual P95: 1.8 seconds
Actual P99: 3.5 seconds (within SLA)

Stage Breakdown (100 GTINs):
  Stage 1 (GTIN→CID): 500ms (vs. 50s serial)
  Stage 2 (Authorization): 50ms (vs. 5s N+1 queries)
  Stage 3 (DC Inventory): 1,200ms (vs. 120s serial)
  Total Pipeline: 1,750ms

Comparison to Similar APIs:
  inventory-status-srv: 2.7s P95
  DC Inventory Search: 1.8s P95 (33% faster)
```

### Success Rate
```
Overall Success Rate: 98%
Authorization Failures: 1.5% (supplier requested unauthorized GTIN)
EI API Failures: 0.5% (transient network issues)
Partial Success Rate: 12% (some GTINs succeed, some fail)
```

### Thread Pool Metrics
```
Pool Size: 20 threads (dedicated)
Average Utilization: 12 threads active
Peak Utilization: 18 threads active
Queue Capacity: 100 requests
Average Queue Size: < 10 requests
```

### Scalability Metrics
```
Current: 30K queries/day
10x Growth: 300K queries/day = 24M GTIN lookups/day
Changes Needed for 10x:
  - Thread pool: 20 → 50 threads
  - UberKey rate limit: 100 req/sec → 300 req/sec
  - EI batch size: 50 GTINs → 100 GTINs per request
Estimated Cost Increase: Minimal (thread scaling free, rate limits negotiable)
```

### Delivery Metrics
```
Estimated Delivery: 12 weeks (by EI team)
Actual Delivery: 4 weeks (70% faster)
Design Phase: 1 week (no formal spec, reverse-engineered EI APIs)
Implementation Phase: 3 weeks
Production Launch: Zero issues
Pattern Adoption: 3 other teams copied 3-stage pipeline pattern
```

---

## SYSTEM 3: SPRING BOOT 3 MIGRATION

### Migration Scope
```
Services Migrated: 6 services
  - cp-nrti-apis (18,000 lines of code)
  - inventory-events-srv (15,000 lines)
  - inventory-status-srv (14,000 lines)
  - audit-api-logs-srv (8,000 lines)
  - audit-api-logs-gcs-sink (3,000 lines)
  - dv-api-common-libraries (696 lines)
Total Lines of Code: 58,696 lines
```

### Test Failure Metrics
```
Initial Test Failures: 203 failures (out of 487 tests = 42% failure rate)

Failure Categories:
  - NullPointerException: 87 tests (43%)
  - SecurityException: 45 tests (22%)
  - HibernateException: 38 tests (19%)
  - Miscellaneous: 33 tests (16%)

Resolution Timeline:
  - Day 1-2: Fixed NPE pattern (87 tests)
  - Day 3-4: Fixed Security config (45 tests)
  - Day 5: Automated Hibernate fixes (38 tests)
  - Day 6-7: Manual edge cases (33 tests)
  - Final: 0 test failures (100% passing)
```

### Timeline Metrics
```
Big Bang Approach (Original Plan):
  - Estimated Time: 101 hours
  - Available Time: 80 hours (2 weeks)
  - Risk: HIGH (would miss deadline)

Phased Approach (Pivot):
  - Estimated Time: 53 hours
  - Actual Time: 48 hours
  - Time Saved: 53 hours (52% reduction)
  - Result: ON TIME delivery
```

### Production Stability
```
Rollback Incidents: 0
Post-Migration Bugs: 0
Performance Regression: 0%
Canary Deployment: 10% → 50% → 100% (no errors detected)
Time in Canary: 2 hours (Flagger automatic promotion)
```

### Pattern Reuse
```
Constructor Injection Pattern: 87 NPE fixes → Applied to 4 other services
Jakarta Persistence sed Script: 38 fixes → Used by 3 other teams
Phased Migration Runbook: Documented for 5 remaining services
Team Tech Talk: "Spring Boot 3 Pitfalls" (200+ attendees)
```

---

## SYSTEM 4: MULTI-REGION KAFKA ARCHITECTURE

### Architecture Metrics
```
Regions: 2 (EUS2 primary, SCUS secondary)
Kafka Clusters: 2 (Active-Active)
Brokers per Cluster: 3 (6 total)
Replication Factor: 3
Partitions: 12
```

### Failover Metrics
```
Recovery Time Objective (RTO): < 30 seconds
Recovery Point Objective (RPO): 0 seconds (zero data loss)
Automatic Failover Time: 25 seconds actual
Manual Failover Time: N/A (fully automatic)
Failover Success Rate: 100% (3 failovers in 6 months, all successful)
```

### Latency Metrics
```
Single-Write (Synchronous): 45ms P95 (wait for both clusters)
Dual-Write (Asynchronous): 12ms P95 (fire-and-forget)
Improvement: 73% latency reduction
```

### Cost Metrics
```
Estimated (Design Phase): $3,500/month
Actual (Production): $3,200/month (under budget)
Alternative (Active-Passive): $2,000/month
Alternative (No DR): $1,200/month
Decision: Chose Active-Active for zero RPO (business requirement)
```

### Reliability Metrics
```
Data Loss Incidents: 0 (6 months in production)
Duplicate Message Rate: 0.02% (idempotent producer + consumer deduplication)
Consumer Rebalance Time: < 30 seconds (automatic)
Cluster Failures Handled: 3 (all automatic failover)
```

### Adoption Metrics
```
Initial Service: audit-api-logs-srv (pilot)
Additional Services: inventory-events-srv, inventory-status-srv
Pattern Adoption: 2 other teams (outside original scope)
Reference Architecture: Walmart Kafka best practices updated
```

---

## SYSTEM 5: DSD NOTIFICATION SYSTEM

### Notification Volume
```
Total Notifications (6 months): 500,000+
Daily Notifications: ~2,700
Peak Notifications: 120/hour (morning receiving window)
```

### Notification Types
```
Push Notifications (Store Associates): 300,000 (60%)
Email Notifications (Suppliers): 150,000 (30%)
SMS Notifications (Added Week 8): 50,000 (10%)
```

### Delivery Metrics
```
Push Notification Delivery Rate: 97%
Email Open Rate: 92% (suppliers)
Email Click-Through Rate: 45%
SMS Delivery Rate: 99%
Average Notification Latency: 3.2 seconds
```

### Business Impact
```
Shipment Wait Time Reduction: 40%
Store Associate Satisfaction: 4.5/5.0 (survey)
Supplier Complaints (Spam): 0
Receiving Efficiency: 25% improvement (faster check-in)
```

### System Evolution
```
Launch (Week 6): 2 consumers (push, email)
Week 8: +1 consumer (SMS)
Week 10: +1 consumer (photo upload)
Week 12: +1 consumer (analytics)
Total Consumers: 5 (vs. 2 at launch)
Code Changes to DSC API: 0 (event-driven extensibility)
```

### Reliability Metrics
```
Notification System Uptime: 99.5%
Kafka Event Delivery: 99.99%
Sumo API Uptime: 98.5% (external dependency)
Deduplication Effectiveness: 99.8% (Redis cache)
```

---

## SYSTEM 6: COMMON LIBRARY (dv-api-common-libraries)

### Library Metrics
```
Lines of Code: 696 lines (production)
Test Lines of Code: 678 lines
Test Coverage: 97.4% (678/696)
Maven Releases: 57+ versions
Latest Version: 0.0.54 (used by cp-nrti-apis)
```

### Adoption Metrics
```
Total Services Using: 12+ services
Teams Integrated: 8 teams
Integration Time: < 1 hour per service (just dependency + config)
Code Changes Required: 0 lines (automatic instrumentation)
Integration Rate: 12 services in 8 weeks = 1.5 services/week
```

### Performance Metrics
```
Latency Impact on APIs: 0ms (async thread pool)
Thread Pool:
  - Core: 6 threads
  - Max: 10 threads
  - Queue: 100 capacity
Audit Log Publish Time: < 50ms P95
```

### Reusability Metrics
```
Shared Across Services: 100% code reuse
Custom Logic per Service: 0% (config-driven)
CCM Configuration Lines per Service: 10-15 lines
Maven Dependency: 1 line
Integration Effort: < 1 hour per service
```

---

## SYSTEM 7: MULTI-MARKET ARCHITECTURE (US/CA/MX)

### Market Metrics
```
Markets Supported: 3 (US, Canada, Mexico)
Site IDs:
  - US: 1
  - Mexico: 2
  - Canada: 3
```

### Volume by Market
```
US: 6,000,000 queries/month (75%)
Canada: 1,200,000 queries/month (15%)
Mexico: 800,000 queries/month (10%)
Total: 8,000,000 queries/month
```

### Data Isolation
```
Cross-Market Data Leaks: 0 (perfect isolation)
Database Partition Key: site_id (automatic filtering)
Compliance Audits Passed: 3 markets (US, CA, MX)
```

### Deployment Metrics
```
Code Reuse: 95% (only config differs per market)
Rollout Timeline:
  - Week 6: Deploy with feature flags OFF (US-only)
  - Week 7: Enable Canada (pilot with 1 supplier)
  - Week 8: Enable Canada (all suppliers)
  - Week 9: Enable Mexico (pilot)
  - Week 10: Enable Mexico (all suppliers)
Breaking Changes: 0 (US functionality unchanged)
```

### Configuration Management
```
CCM Config Files: 3 (usEiApiConfig, caEiApiConfig, mxEiApiConfig)
Site-Specific Endpoints:
  - US: ei-inventory-history-lookup.walmart.com
  - CA: ei-inventory-history-lookup-ca.walmart.com
  - MX: ei-inventory-history-lookup-mx.walmart.com
Factory Pattern: SiteConfigFactory (3 site configs)
```

---

## CROSS-SYSTEM METRICS

### Overall Scale
```
Total Services: 6 production services
Total Lines of Code: 58,696 lines
Total API Endpoints: 47+ endpoints
Total Daily Requests: 100,000+ requests/day
Total Data Processed: 2M+ events/day
```

### Team Impact
```
Services Delivered: 6 services in 6 months (1 service/month)
Teams Supported: 12+ teams
External Teams Adopted Patterns: 5+ teams
Reference Architectures Created: 3 (Kafka audit, DC search, multi-market)
```

### Cost Savings
```
Audit System: $59,280/year saved
Total Estimated: $100,000+ annual savings across all optimizations
```

### Reliability
```
Production Incidents (6 months): 0 (zero rollbacks)
Uptime: 99.9% average across all services
Zero-Downtime Deployments: 100% (canary + Flagger)
```

---

## MEMORIZATION TIPS

### Round Numbers for Quick Recall
```
"Around 2 million events per day" ✓
"Approximately 30,000 queries daily" ✓
"About 60 dollars per month" ✓
"Roughly 200 test failures" ✓
```

### Percentage Comparisons
```
"85% faster queries" (8s → 1.2s)
"90% cost reduction" ($5K → $500)
"40% faster than similar APIs" (2.7s → 1.8s)
"70% faster delivery" (12 weeks → 4 weeks)
```

### Before/After Stories
```
"Before migration, PostgreSQL crashed 2x/month. After Kafka, zero crashes in 6 months."
"Before optimization, 50 seconds for 100 GTINs. After parallelization, 1.8 seconds."
"Before Spring Boot 3, 203 test failures. After phased approach, zero failures in 48 hours."
```

---

## INTERVIEW USAGE EXAMPLES

### When They Ask: "How much scale?"
```
"The Kafka audit system processes 2 million events per day, with peak rates
of 120 events per second. We've tested it to 50 million events per day,
so we have 25x headroom for growth."
```

### When They Ask: "How fast?"
```
"The DC Inventory Search API responds in 1.8 seconds at P95 for 100 GTINs,
which is 33% faster than our similar inventory-status-srv API (2.7s P95).
The secret: 3-stage pipeline with parallel processing."
```

### When They Ask: "How much did it cost?"
```
"We reduced audit logging costs from $5,000 per month with PostgreSQL
to $60 per month with Kafka + GCS + BigQuery. That's a 90% cost reduction,
saving $59,280 annually."
```

### When They Ask: "How reliable?"
```
"The multi-region Kafka architecture has a Recovery Time Objective of
under 30 seconds and a Recovery Point Objective of zero seconds. We've
had 3 failover events in 6 months, all automatic, zero data loss."
```

---

**END OF METRICS CHEATSHEET**

**Total Metrics Captured**: 150+ specific numbers across 7 systems

**Memorization Strategy**:
1. Print this cheatsheet
2. Review daily for 1 week before interview
3. Practice recalling 3-5 metrics per system
4. Use in mock interviews

**Google Interviewer Will Ask**: "Be specific - what were the actual numbers?"
**You Answer**: "The system processed exactly 2,000,000 events per day, with peak rates of 120 events per second. We tested capacity to 50 million events per day, giving us 25x headroom for growth. Cost was reduced from $5,000 per month to $60 per month, an annual savings of $59,280."

**Result**: Interviewer thinks: "This candidate knows their numbers cold. They've built real systems at scale."

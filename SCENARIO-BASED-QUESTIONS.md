# SCENARIO-BASED QUESTIONS - "What If X Fails/Breaks?"

> **This is what hiring managers test to see if you TRULY understand your systems.**
> "Tell me what happens when..." type questions separate juniors from seniors.
> Covers: Kafka, SUMO, Downstream APIs, Database, Kubernetes, Caching, and more.

---

# SECTION 1: KAFKA FAILURE SCENARIOS

## Q: "What happens if Kafka is down?"

### In Audit Logging System (Tier 2: audit-api-logs-srv)

> "Two scenarios:
>
> **If primary Kafka (EUS2) is down:** The publisher tries the secondary Kafka (SCUS). We implemented dual-region failover in PR #65 using CompletableFuture chaining. If primary fails, `.exceptionally()` triggers secondary. Messages are not lost.
>
> **If BOTH Kafka clusters are down:** The publisher throws `AuditLoggingUnavailableException`. But this doesn't affect the API - the common library's @Async call catches all exceptions. The API response returns normally. Audit events are lost during the outage. We see this in our metrics (audit error rate spikes) and get alerted.
>
> **Key point:** Kafka being down NEVER impacts supplier API responses. The fire-and-forget design ensures the API is decoupled from Kafka health."

### In DSD Notification System

> "DSD publishes shipment events to Kafka for persistence. If Kafka is down, the event isn't persisted but the SUMO notification can still be sent (it's a separate HTTP call). We'd lose the audit trail but not the notification. If Kafka is critical, we'd need to either queue locally and retry, or fail the API call and let the supplier retry."

---

## Q: "What if a Kafka topic gets deleted accidentally?"

> "All committed messages are lost. But:
> - Kafka Connect's last committed offset becomes invalid
> - Connect worker will fail to start and alert fires
> - We'd need to recreate the topic and re-publish any messages from source
>
> For audit logging, the common library doesn't buffer messages locally, so anything not yet published is lost. This is why we have multi-region - if one region's topic is deleted, the other still has the data.
>
> **Prevention:** Topic deletion is restricted to admin roles. We have ACLs on Kafka topics."

---

## Q: "What if Kafka consumer lag keeps growing?"

> "We have two-tier alerting:
> - **Warning at 100 messages** for 5 min → Investigate
> - **Critical at 500 messages** for 5 min → Act immediately
>
> Causes could be:
> 1. **Slow GCS writes** → Increase batch size, check GCS health
> 2. **Consumer processing errors** → Check error logs, DLQ
> 3. **Insufficient consumers** → Scale Connect workers (but NOT on lag - learned from KEDA incident!)
> 4. **Large message backlog** → One-time event, will clear naturally
>
> **The KEDA lesson:** Don't autoscale on lag. It causes rebalance storms. Scale on CPU instead."

---

## Q: "What if Schema Registry is down?"

> "Producers cache schemas locally after first registration. So:
> - **Existing schemas:** Continue working (cached)
> - **New schemas:** Can't register → producer fails → alert fires
> - **Schema evolution:** Blocked until Registry recovers
>
> Impact: If all cached schemas are valid, zero impact on existing messages. Only new schema versions are blocked."

---

## Q: "What if messages are out of order?"

> "For audit logging, order doesn't matter - each audit event is independent with its own request_id and timestamp. BigQuery queries use timestamp for ordering.
>
> For DSD, order matters: ENROUTE should come before ARRIVED. If they arrive out of order, an associate might get 'arrived' before 'enroute'. We handle this by including `eventCreationTime` in the Kafka message. Consumers can sort by this timestamp. At the notification level, if ARRIVED comes first, the associate still goes to the dock - the outcome is the same."

---

# SECTION 2: DATABASE FAILURE SCENARIOS

## Q: "What if PostgreSQL is down?"

> "Different impact per service:
>
> **inventory-status-srv (DC Inventory):** API calls fail with 500. We'd return error responses for all items. The supplier sees the error and can retry. We'd get 5xx alerts within 30 seconds.
>
> **inventory-events-srv (Transaction History):** Same - 500 errors, alerts fire.
>
> **audit-api-logs-srv:** No PostgreSQL dependency - uses Kafka. Unaffected.
>
> **Mitigation:** WCNP (Walmart Cloud Native Platform) manages PostgreSQL with automated failover. Typical recovery: <1 minute. We also have connection pool retry logic."

---

## Q: "What if a database migration fails midway?"

> "This happened during the Cosmos to PostgreSQL migration (PR #80). The approach:
>
> 1. **Dual-write period:** Write to both Cosmos and PostgreSQL for a transition period
> 2. **Read from old (Cosmos)** until new (PostgreSQL) is validated
> 3. **Switch reads** to PostgreSQL once data parity is confirmed
> 4. **Decommission** Cosmos after stable period
>
> If migration fails midway: roll back to Cosmos reads. No data loss because old DB is still active. This is blue/green for databases."

---

## Q: "What if the authorization cache (megha-cache) goes down?"

> "Every API request checks supplier authorization. Without cache:
>
> - **Immediate:** All requests hit PostgreSQL directly
> - **Impact:** Latency increases (50ms cache → 200ms database)
> - **Risk:** If traffic is high, PostgreSQL could get overloaded
>
> **Mitigation:**
> - megha-cache has VM availability alerts (<90% triggers)
> - Cache has high availability across multiple nodes
> - If total failure, we'd see latency alerts and could temporarily increase DB connection pool
>
> **Long-term fix:** Add Caffeine as a local L1 cache in front of megha-cache (L2). L1 miss → L2 miss → Database."

---

# SECTION 3: DOWNSTREAM API FAILURE SCENARIOS

## Q: "What if the EI (Enterprise Inventory) API is slow or down?"

> "For DC Inventory API:
>
> - **Slow (>2s):** Our request times out. We return partial response - items that succeeded before timeout + errors for items that didn't.
> - **Down (5xx):** All items fail at Stage 3. We return all errors with `errorSource: EI_SERVICE`.
>
> **What we added after Spring Boot 3 migration (PR #1564):**
> - Explicit timeout per request
> - Retry with exponential backoff (2-3 retries)
> - `NrtiUnavailableException` when retries exhausted
>
> **What I'd add next:** Resilience4j circuit breaker. After 5 consecutive failures, stop calling EI for 30 seconds. Return cached data or 'service unavailable' immediately."

---

## Q: "What if SUMO push notification service is down?"

> "For DSD notifications:
>
> 1. `SignatureHandleException` is caught → logged, not propagated
> 2. DSD API response to supplier is NOT affected
> 3. DSD event is still published to Kafka (persistence intact)
> 4. We can disable SUMO via feature flag (`isSumoEnabled = false`)
>
> **What we lose:** Associates don't get real-time alerts. They fall back to periodic checks.
>
> **What we keep:** All event data in Kafka. Can replay notifications when SUMO recovers.
>
> **What I'd add:**
> - Circuit breaker to stop calling SUMO after N failures
> - Queue failed notifications for retry when SUMO recovers
> - SMS fallback for critical ARRIVED events (high-value shipments)"

---

## Q: "What if the UberKey service is down?"

> "UberKey converts WmItemNumbers to GTINs for DC Inventory.
>
> - **Impact:** Stage 1 of the 3-stage pipeline fails. ALL items get `errorSource: UBERKEY` errors.
> - **Response:** HTTP 200 with all items showing error status.
> - **Supplier sees:** 'Invalid WmItemNumber' for every item.
>
> **What I'd add:** Cache frequently used WmItemNbr-to-GTIN mappings. If UberKey is down, serve from cache for known items."

---

# SECTION 4: KUBERNETES / INFRASTRUCTURE SCENARIOS

## Q: "What happens during a pod restart?"

> "Depends on the service:
>
> **Stateless services (cp-nrti-apis, inventory-status-srv):** Zero impact. K8s routes traffic to other pods. Readiness probe marks the restarting pod as unavailable.
>
> **Kafka Connect (gcs-sink):** Consumer group rebalance. Other workers pick up partitions. Messages are NOT lost (Kafka retains them). Lag increases temporarily but clears when pod recovers.
>
> **Key:** This is why we have `startupProbe`, `readinessProbe`, and `livenessProbe` configured. Startup probe gives Spring Boot time to initialize. Readiness probe ensures traffic only goes to ready pods."

---

## Q: "What if a canary deployment detects an issue?"

> "Flagger handles this automatically:
>
> 1. Flagger checks metrics every 2 minutes (our config: `interval: 2m`)
> 2. If `request-success-rate < 99%` OR `request-duration P99 > 500ms`
> 3. Flagger increments failure counter
> 4. After 5 failures (`threshold: 5`): **AUTOMATIC ROLLBACK**
> 5. All traffic shifts back to old version instantly via Istio VirtualService
> 6. Alert fires to Slack
>
> **No human intervention needed.** The old version pods are still running. Rollback is instant."

---

## Q: "What if Kubernetes HPA scales up too aggressively?"

> "This is the KEDA lesson. Aggressive autoscaling on wrong metrics causes cascading failures:
>
> **What happened:** KEDA scaled Kafka Connect on consumer lag. More workers → rebalancing → more lag → more scaling → infinite loop.
>
> **What we do now:** Scale on CPU utilization (70% threshold), not business metrics. CPU is a stable, predictable metric. Business metrics (lag, queue depth) can create feedback loops.
>
> **General principle:** Scale on resource metrics (CPU, memory), not application metrics. Application metrics can create feedback loops that resource metrics can't."

---

# SECTION 5: DATA & COMPLIANCE SCENARIOS

## Q: "What if a supplier requests GDPR data deletion?"

> "Hard with immutable Parquet in GCS. Current process:
>
> 1. Identify all Parquet files containing supplier's data
> 2. Read each file, filter out supplier's records, write new file
> 3. Delete original file
> 4. Verify deletion in BigQuery
>
> **Cost:** Expensive (read-write entire files). Slow (millions of files over 7 years).
>
> **Better approach (roadmap):**
> - **Crypto-shredding:** Encrypt each record with supplier-specific key. For deletion, delete the key. Data remains but is unreadable.
> - **Apache Iceberg:** Supports row-level deletes with ACID transactions."

---

## Q: "What if audit data is inconsistent between regions?"

> "Active/Active means both regions process independently. Possible inconsistency:
>
> 1. **Same record in both regions:** Deduplication by `request_id` in BigQuery
> 2. **Record in only one region:** Geographic routing via `wm-site-id` minimizes this. US records go to EUS2, CA records go to SCUS.
> 3. **Regional outage during write:** Kafka retains messages. When region recovers, backlog processes.
>
> **Validation:** We run periodic data parity checks comparing record counts between regions."

---

## Q: "What if someone deploys a breaking schema change?"

> "Schema Registry prevents this:
>
> - **Backward compatibility mode:** New schema must be readable by old consumers
> - If incompatible schema is registered → Registry REJECTS it → Producer fails fast
> - Existing messages continue to work with cached schema
>
> **In practice:** We add fields with defaults (backward compatible). We never remove required fields. Schema evolution is part of our PR review checklist."

---

# SECTION 6: WALMART-SPECIFIC SCENARIOS

## Q: "What happens during Black Friday / peak traffic?"

> "Our audit system at current load: 2M events/day = ~23/sec. Peak could be 10-20x.
>
> **Preparation:**
> 1. **Pre-scale:** Increase Connect workers and Kafka partitions BEFORE peak
> 2. **Increase thread pool:** Audit library's core/max pool sizes for higher throughput
> 3. **Monitor lag:** Real-time consumer lag dashboards during peak
> 4. **Disable non-critical:** Can disable response body logging to reduce payload size
>
> **What survives:** The async design means even at 10x traffic, APIs aren't impacted. The question is whether audit pipeline can keep up. Kafka's retention (7 days) gives us buffer."

---

## Q: "What if WCNP (Walmart Cloud Native Platform) has a regional outage?"

> "This is exactly why we built Active/Active multi-region:
>
> 1. EUS2 goes down → SCUS handles all traffic
> 2. Global load balancer (Azure Front Door) routes around the outage
> 3. Kafka in SCUS has all messages (dual publishing)
> 4. GCS sink in SCUS continues writing to GCS
>
> **Tested in DR drill:** Recovered in 15 minutes vs 1-hour RTO target.
>
> **What we CAN'T survive:** Both regions down simultaneously (extremely unlikely with Azure's availability guarantees)."

---

## Q: "What if CCM (Cloud Configuration Management) is down?"

> "CCM is Walmart's runtime configuration service. If it's down:
>
> - **Cached configs continue working** - CCM clients cache locally
> - **New config changes can't be pushed** - feature flags stuck at current state
> - **Impact:** Can't toggle `isSumoEnabled`, can't update endpoint filtering, can't change commodity mappings
>
> **Mitigation:** CCM itself is highly available. Local caching means short outages have zero impact."

---

# SECTION 7: HOW TO ANSWER SCENARIO QUESTIONS

## The Framework

When they ask "What happens if X fails?", use this structure:

```
1. DETECT: "How would we know?" (alerts, metrics, logs)
2. IMPACT: "What's the blast radius?" (users affected, data at risk)
3. MITIGATE: "What happens automatically?" (retries, failover, circuit breaker)
4. RECOVER: "How do we fix it?" (manual steps, replay, rollback)
5. PREVENT: "How do we stop it happening again?" (monitoring, testing, architecture)
```

### Example:

**Q: "What if Kafka Connect worker crashes?"**

> **DETECT:** "Consumer lag alert fires within 5 minutes. K8s event shows pod restart."
>
> **IMPACT:** "Messages queue in Kafka. No data loss - Kafka retains messages. Other workers pick up orphaned partitions after rebalance."
>
> **MITIGATE:** "K8s automatically restarts the pod. Connect resumes from last committed offset. Backlog processes automatically."
>
> **RECOVER:** "Usually automatic - pod restarts, catches up. If OOM, we may need to increase heap (learned from PR #57)."
>
> **PREVENT:** "Proper heap sizing, startup probes, liveness probes. JVM memory alerts in Grafana."

---

## The Power Move

After answering any scenario question, add:

> "We actually faced something similar. [Tell a real story from your production issues]. That's when we added [specific monitoring/fix]."

This connects hypothetical scenarios to your REAL production experience. Much more impressive than theoretical answers.

---

*Practice these scenarios out loud. The framework (DETECT → IMPACT → MITIGATE → RECOVER → PREVENT) works for ANY "what if" question.*

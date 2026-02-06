# ⚠️ GAP ANALYSIS - Areas Needing More Preparation

> **Purpose:** Identify weak spots and prepare defenses
> **Last Updated:** February 2026

---

## 🔴 HIGH PRIORITY GAPS

### Gap 1: Bullet 3 - DSD Notification System (35% Improvement Claim)

**The Issue:**
Your resume claims "35% improvement in stock replenishment timing" but your documentation doesn't explain:
- How was this measured?
- What was the baseline?
- What's the direct causation?

**Questions You Might Get:**
- "How did you measure the 35% improvement?"
- "What was the replenishment time before vs after?"
- "How do you know the notification system caused the improvement?"

**Prepare This Answer:**
> "We measured replenishment timing as the gap between DSD shipment arrival notification and inventory update in the system. Before SUMO push notifications, store associates discovered deliveries through periodic checks - average lag was X hours. After implementing real-time push notifications, they responded within Y minutes. The 35% improvement is the reduction in this lag time. We tracked this through timestamps in our transaction events."

**Action Items:**
- [ ] Find the actual before/after metrics in your PRs or documentation
- [ ] Understand how SUMO push notifications work
- [ ] Be ready to explain the causation chain

---

### Gap 2: Bullets 5 & 6 - OpenAPI & DC Inventory (Less Documented)

**The Issue:**
Your documentation focuses heavily on Kafka audit (Bullets 1, 2, 10) and Spring Boot 3 (Bullet 4), but has less depth on:
- OpenAPI design-first approach (Bullet 5)
- DC Inventory Search API with CompletableFuture (Bullet 6)

**Questions You Might Get:**
- "Walk me through the DC Inventory API implementation"
- "How did CompletableFuture improve performance?"
- "What does 'design-first' mean in practice?"

**Prepare These Answers:**

**For DC Inventory:**
> "The DC Inventory API allows suppliers to query inventory across multiple distribution centers. The challenge was that querying 10+ DCs serially was too slow - 3+ seconds. I used CompletableFuture.allOf() to query all DCs in parallel, then aggregate results. This reduced response time from ~3 seconds to ~800ms - the slowest DC call plus aggregation overhead."

**For OpenAPI Design-First:**
> "Design-first means we write the OpenAPI spec before any code. The spec becomes the contract - it defines request/response schemas, validation rules, error codes. We use openapi-generator-maven-plugin to generate server stubs from the spec. The benefit is that consumers can start integration work using the spec while we build the implementation. We also integrated R2C (Request to Contract) testing to validate implementations match the spec."

**Action Items:**
- [ ] Review PRs #271, #260 in inventory-status-srv
- [ ] Understand CompletableFuture.allOf() pattern
- [ ] Review contract testing setup (R2C)

---

### Gap 3: Bullet 7 - Transaction Event History API

**The Issue:**
While documented, you should be able to explain:
- The multi-tenant architecture
- Cursor-based pagination (why cursor over offset?)
- The Cosmos to PostgreSQL migration

**Questions You Might Get:**
- "Why cursor pagination instead of offset?"
- "How do you handle multi-tenancy?"
- "Why migrate from Cosmos to PostgreSQL?"

**Prepare These Answers:**

**Cursor Pagination:**
> "Offset pagination (`LIMIT 10 OFFSET 1000`) is expensive for large datasets - the database has to scan all 1000 rows before returning the next 10. Cursor pagination uses a bookmark (like `WHERE created_ts > last_seen_ts`) which is efficient because it hits an index. For our transaction history with millions of events, cursor pagination keeps response times constant regardless of how deep you page."

**Multi-Tenancy:**
> "We use site-based partitioning. Each request includes a siteId (US, CA, MX). The PostgreSQL tables are partitioned by site, and queries include the partition key. This provides data isolation - Canadian suppliers can't accidentally query US data - and performance benefits from partition pruning."

**Cosmos to PostgreSQL:**
> "We migrated from Cosmos to PostgreSQL for cost and operational simplicity. Cosmos pricing is based on RUs (request units) which made costs unpredictable. PostgreSQL on WCNP (Walmart's managed database) has predictable costs and we have better operational expertise with it."

---

### Gap 4: Bullet 8 - Supplier Authorization Framework

**The Issue:**
The 4-level authorization hierarchy is mentioned but needs clearer explanation:
- Consumer → DUNS → GTIN → Store

**Questions You Might Get:**
- "Explain the authorization hierarchy"
- "How does Caffeine caching help?"
- "What happens on a cache miss?"

**Prepare This Answer:**
> "The authorization hierarchy works like a tree:
>
> 1. **Consumer**: The API consumer (e.g., 'PepsiCo Inc')
> 2. **DUNS**: Data Universal Numbering System - identifies legal entities (e.g., 'PepsiCo Frito-Lay Division')
> 3. **GTIN**: Global Trade Item Number - product identifier (e.g., '00012345678905' = Lay's Classic Chips)
> 4. **Store**: Individual store locations (e.g., 'Store #4236')
>
> Authorization is checked hierarchically: Does this consumer have access? → For this DUNS? → For this GTIN? → At this store? Any 'no' returns 403.
>
> We cache authorization data in Caffeine with 7-day TTL because this data changes rarely but is queried on every request. On cache miss, we hit PostgreSQL and populate the cache."

---

## 🟡 MEDIUM PRIORITY GAPS

### Gap 5: System Design Questions

**The Issue:**
Be ready for "how would you design X" questions that test broader thinking.

**Possible Questions:**
- "How would you design the audit system to handle 100x traffic?"
- "What if you needed real-time querying instead of batch?"
- "How would you add a new region (e.g., Europe)?"

**Your 100x Traffic Answer:**
> "Current: 2M/day = 23/sec. Target: 200M/day = 2300/sec.
>
> 1. **Library**: Replace HTTP-to-publisher with direct Kafka publishing
> 2. **Kafka**: Increase partitions from 12 to 48+ for parallelism
> 3. **Connect**: Scale workers horizontally, 2-3 more workers
> 4. **GCS**: Can handle it; may need to batch smaller files
>
> Cost would scale ~10x, mostly Kafka partition storage."

---

### Gap 6: Failure/Mistake Questions

**The Issue:**
Be ready to discuss failures vulnerably but constructively.

**Possible Questions:**
- "Tell me about a time you failed"
- "What's your biggest mistake at Walmart?"
- "What would you do differently?"

**Your Answer:**
> "The KEDA autoscaling incident was my mistake. I configured autoscaling based on consumer lag without understanding the implications - scaling triggers rebalancing, which increases lag, which triggers more scaling. The system got stuck in a feedback loop.
>
> I should have tested autoscaling with production-like traffic patterns, not just unit tests. Now I always test scaling scenarios with realistic load before production."

---

### Gap 7: Why Are You Leaving Walmart?

**The Issue:**
This isn't technical but is guaranteed. Have a positive, brief answer.

**Your Answer:**
> "I've had a great experience at Walmart - I've grown from building my first production system to leading major migrations. I'm looking for my next challenge - specifically, I'm excited about [company's mission/product] because [specific reason]. I want to continue growing as an engineer and have impact at [scale/stage] of a company."

---

## 🟢 NICE TO HAVE

### Gap 8: Observability Details (Bullet 9)

If asked about specific tools:
- **Dynatrace**: APM (Application Performance Monitoring) - traces, metrics
- **Flagger**: Canary deployment controller for Kubernetes
- **Prometheus**: Metrics collection and alerting
- **Grafana**: Dashboards and visualization

### Gap 9: Kubernetes/Infrastructure Knowledge

If asked about deployment:
- Know basics of: Pods, Deployments, Services, Ingress
- Know Istio basics: Service mesh, traffic splitting, mTLS
- Know WCNP (Walmart Cloud Native Platform) is Kubernetes-based

---

## 📝 PREPARATION CHECKLIST

### Before Your Interview

- [ ] Review DC Inventory PRs (#271, #260)
- [ ] Practice explaining CompletableFuture parallel execution
- [ ] Review Transaction History PRs (#8, #38)
- [ ] Practice explaining cursor pagination
- [ ] Prepare concrete metrics for DSD notification improvement
- [ ] Review SUMO push notification integration
- [ ] Practice "why are you leaving" answer
- [ ] Review your failure story (KEDA incident)

### Numbers You MUST Know

| Bullet | Metric | Value |
|--------|--------|-------|
| 1 | Events/day | 2M+ |
| 1 | Latency impact | <5ms P99 |
| 2 | Teams adopted | 3+ |
| 2 | Integration time | 2 weeks → 1 day |
| 3 | Associates | 1,200+ |
| 3 | Stores | 300+ |
| 3 | Improvement | 35% |
| 4 | Files changed | 136 |
| 4 | javax→jakarta | 74 files |
| 5 | Integration reduction | 30% |
| 6 | Query time reduction | 40% |

---

## 🚨 RED FLAG QUESTIONS (BE CAREFUL!)

### "Why did you choose X over Y?"

**Trap:** They want to see if you blindly chose technology or made informed decisions.

**Good Answer Format:**
> "We considered [alternatives]. We chose [X] because [specific reasons]. The trade-off was [downside], which we mitigated by [action]."

### "What are the trade-offs of your design?"

**Trap:** No design is perfect. If you say "none," you seem inexperienced.

**Good Answer:**
> "The main trade-off of async fire-and-forget is we can lose audit data if the queue fills up. We mitigated this with metrics and alerting, but it's a conscious trade-off - audit data is valuable but not worth blocking API responses."

### "What would you do differently?"

**Trap:** Saying "nothing" seems like you don't learn. Saying too much seems like regret.

**Good Answer:**
> "Two things: Add OpenTelemetry tracing from day one for better observability. And investigate automated migration tools for the namespace changes - we did it semi-manually which was time-consuming."

---

*Address these gaps before your interview!*

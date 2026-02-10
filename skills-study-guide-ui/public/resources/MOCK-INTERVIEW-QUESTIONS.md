# 🎤 MOCK INTERVIEW QUESTIONS - Complete Q&A Bank

> **Purpose:** Practice answering these questions out loud until they feel natural
> **Last Updated:** February 2026

---

## HOW TO USE THIS DOCUMENT

1. **Cover the answer** - Read only the question first
2. **Answer out loud** - Don't just think it, speak it
3. **Compare** - Check against the model answer
4. **Time yourself** - Most answers should be 1-3 minutes
5. **Practice variations** - The same question can be asked many ways

---

# SECTION 1: TECHNICAL DEEP DIVE QUESTIONS

---

## Q1: Walk me through the Kafka Audit Logging architecture

**What they're testing:** System design, component understanding, trade-offs

### Model Answer (2-3 minutes)

> "I designed a three-tier architecture for audit logging that processes over 2 million events per day with less than 5ms P99 latency impact on the source APIs.
>
> **Tier 1** is a Spring Boot Starter library that lives in `dv-api-common-libraries`. It uses a Servlet Filter with `@Order(LOWEST_PRECEDENCE)` so it runs last - after security filters. This is important because we capture the final HTTP response, including auth failures. The filter uses `ContentCachingRequestWrapper` to read the HTTP body without consuming it, then sends the audit payload asynchronously using `@Async` with a dedicated thread pool.
>
> **Tier 2** is the Kafka publisher service `audit-api-logs-srv`. It receives the HTTP payload, serializes it to Avro using Schema Registry for type safety and 70% size reduction, then publishes to Kafka. We run this Active/Active in two regions - East US 2 and South Central US - for disaster recovery.
>
> **Tier 3** is Kafka Connect with the GCS Sink Connector. I use SMT filters to route messages by geography based on the `wm-site-id` header - separate buckets for US, Canada, and Mexico. Records are written as Parquet files, which gives us 90% compression and native BigQuery compatibility.
>
> The key trade-off is fire-and-forget: we don't block the API response waiting for audit confirmation. If the queue fills up, we could lose audit data. I mitigated this with Prometheus metrics for queue depth, rejected task counters, and alerting at 80% capacity."

### Follow-up Questions to Prepare

- "Why Servlet Filter instead of AOP?" → AOP can't access raw HTTP body
- "Why LOWEST_PRECEDENCE?" → Run after all other filters to capture final state
- "What if Kafka is down?" → Messages queue locally, alerts fire, API continues working
- "How do you handle schema evolution?" → Avro backward compatibility, Schema Registry

---

## Q2: Explain the Spring Boot 3 migration you led

**What they're testing:** Technical leadership, risk management, attention to detail

### Model Answer (2-3 minutes)

> "I led the migration of `cp-nrti-apis` from Spring Boot 2.7 to 3.2, which also required upgrading from Java 11 to 17. This was driven by end-of-life security concerns - Snyk was flagging vulnerabilities we couldn't patch without upgrading.
>
> The migration had three major challenges:
>
> **First**, the javax-to-jakarta namespace change. Java EE became Jakarta EE, so every import of `javax.persistence`, `javax.validation`, `javax.servlet` had to change to `jakarta.*`. I updated 74 files across entities, controllers, validators, and filters.
>
> **Second**, RestTemplate deprecation. Spring Boot 3 pushes toward WebClient, which is reactive. Our codebase was synchronous. I made a strategic decision to use `WebClient` with `.block()` calls for backwards compatibility rather than rewriting everything to be reactive - that would've been a separate, larger initiative.
>
> **Third**, Hibernate 6 compatibility. PostgreSQL enum types broke because Hibernate changed how it maps them. I had to add explicit `@JdbcTypeCode(SqlTypes.NAMED_ENUM)` annotations to all enum fields.
>
> For deployment, I used Flagger canary releases: 10% traffic initially, monitoring error rates and latency, then gradual increase to 25%, 50%, 100% over 24 hours. Automatic rollback was configured if error rate exceeded 1%.
>
> The result was zero customer-impacting issues. 158 files changed, 1,732 lines added. Three minor post-migration fixes - all caught proactively through monitoring before users noticed."

### Follow-up Questions to Prepare

- "Why `.block()` instead of full reactive?" → Risk management - too many changes at once
- "What was the hardest part?" → WebClient test mocking - method chain requires complex mocking
- "What would you do differently?" → Investigate OpenRewrite for automated namespace migration

---

## Q3: Tell me about the debugging incident with Kafka Connect

**What they're testing:** Debugging skills, systematic thinking, learning from failure

### Model Answer (3-4 minutes)

> "Two weeks after launching the Kafka GCS sink in production, I noticed GCS buckets stopped receiving data. This was May 2025. We were losing compliance-critical audit data, but there were no errors in our dashboards. The system was failing silently.
>
> I took a systematic approach over five days:
>
> **Day 1:** I checked the obvious - Kafka Connect running? Yes. Messages in topic? Yes, millions backing up. So the issue was between consumption and GCS write.
>
> **Day 2:** I enabled DEBUG logging and found NullPointerException in our SMT filter. Legacy data didn't have the `wm-site-id` header we expected. I added try-catch with graceful fallback - default to US bucket if header missing.
>
> **Day 3:** Problem persisted. I noticed consumer poll timeouts. The default `max.poll.interval.ms` was 5 minutes, but our GCS writes for large batches took longer. I tuned this and related configs.
>
> **Day 4:** Still issues. I correlated with Kubernetes events and discovered KEDA autoscaling was causing instability. When lag increased, KEDA scaled up workers. But scaling triggered consumer group rebalancing, which caused more lag, which triggered more scaling - a feedback loop! I disabled KEDA and switched to CPU-based HPA.
>
> **Day 5:** Traced remaining issues to JVM heap exhaustion. Default 512MB wasn't enough for large batch Avro deserialization. Increased to 2GB.
>
> **Result:** Zero data loss because Kafka retained all messages. Backlog cleared in 4 hours. I created a troubleshooting runbook that's been used twice by other teams."

### Follow-up Questions to Prepare

- "What was the root cause?" → Multiple issues compounding - null headers, poll timeout, KEDA feedback loop, heap size
- "Why didn't you catch this in testing?" → Didn't test autoscaling with production-like traffic patterns
- "What monitoring did you add?" → Consumer lag alerts, JVM memory metrics, dropped record counters

---

## Q4: How does your async audit logging work technically?

**What they're testing:** Deep technical understanding, trade-off awareness

### Model Answer (2 minutes)

> "The async pattern uses Spring's `@Async` annotation with a custom `ThreadPoolTaskExecutor`. The configuration is: 6 core threads, 10 maximum, queue capacity of 100.
>
> When an HTTP request completes, the `LoggingFilter` captures request/response data using `ContentCachingRequestWrapper` and `ContentCachingResponseWrapper`. These wrappers let us read the HTTP body without consuming the input stream, which would break downstream processing.
>
> The filter then calls `auditLogService.publishAuditData()`, which is marked `@Async`. Spring proxies this call and submits it to the thread pool instead of running inline. The HTTP response returns immediately to the client.
>
> The trade-off is fire-and-forget: if the thread pool's queue fills up, `RejectedExecutionHandler` kicks in. By default, it throws an exception, but since we wrap everything in try-catch, the audit would be silently dropped.
>
> I added three safeguards: a Prometheus counter for rejected tasks, a WARN log when queue exceeds 80% capacity, and documentation explaining the trade-off. In production, we've hit the 80% warning once during a downstream slowdown - we caught it early and scaled the thread pool."

### Follow-up Questions to Prepare

- "Why 6/10/100?" → Based on expected throughput, tuned after load testing
- "What happens during high load?" → Queue fills, warning fires, worst case drops logs (but never blocks API)
- "Why not use message queue directly from filter?" → Added latency, complexity; HTTP to internal service is simpler

---

## Q5: Explain the multi-region Kafka architecture

**What they're testing:** Distributed systems knowledge, disaster recovery thinking

### Model Answer (2 minutes)

> "We run Active/Active Kafka across two Azure regions: East US 2 and South Central US. Each region has its own Kafka cluster.
>
> The publisher service uses a `wm-site-id` header to determine routing. US traffic primarily routes to East US 2, with failover to South Central US. The routing is geography-based, not round-robin, to minimize cross-region latency.
>
> Each region has its own Kafka Connect cluster with GCS Sink connectors. The SMT filter checks the header and writes to region-appropriate buckets. The US filter is permissive - it processes records with US header OR no header (legacy data). Canada and Mexico filters are strict - exact header match only.
>
> For disaster recovery, if one region fails, the other handles all traffic. RTO target was 1 hour; we achieved 15 minutes in DR testing. RPO target was 4 hours; we achieved zero data loss because Kafka retains messages.
>
> The challenge was exactly-once processing across regions. Kafka's idempotent producer helps, but we also added deduplication in the sink based on `request_id` - if a record was already processed by that region, the SMT filter drops it."

### Follow-up Questions to Prepare

- "Why Active/Active vs Active/Passive?" → Audit is write-heavy, can't tolerate 30-minute failover
- "How do you handle split-brain?" → Each region processes its geography; deduplication handles overlap
- "Cost implications?" → ~2x infrastructure, justified by compliance requirements

---

# SECTION 2: SYSTEM DESIGN QUESTIONS

---

## Q6: How would you scale the audit system to 100x traffic?

**What they're testing:** Scalability thinking, bottleneck identification

### Model Answer (2 minutes)

> "Current state: 2 million events per day is about 23 events per second. 100x means 200 million per day, or 2,300 per second.
>
> I'd address each tier:
>
> **Library tier:** Replace HTTP-to-publisher with direct Kafka publishing from the library. This eliminates the network hop and the publisher service bottleneck.
>
> **Kafka tier:** Increase partitions from 12 to 48 or more. More partitions means more parallel consumers. Also increase replication factor if we haven't hit durability requirements.
>
> **Connect tier:** Scale workers horizontally - add 2-3 more workers. Kafka Connect scales well because it uses consumer groups. Also tune batch sizes - larger batches mean fewer GCS writes.
>
> **GCS tier:** GCS can handle it; the limit is more on the Connect side. May need to batch smaller files to avoid memory pressure.
>
> **Monitoring:** Add more granular metrics - per-partition lag, per-connector throughput. Set up alerts for lag growth rate, not just absolute values.
>
> Cost would scale roughly 10x, not 100x, because of batching efficiency. Most cost is Kafka partition storage and Connect compute."

### Follow-up Questions to Prepare

- "What's the first bottleneck?" → Publisher service - HTTP overhead
- "Why not go serverless?" → Cold start latency, harder to manage state, Kafka Connect is proven
- "How would you test this?" → Load testing environment, synthetic traffic, chaos engineering

---

## Q7: Design an audit system from scratch

**What they're testing:** System design fundamentals, requirements gathering

### Model Answer (Framework, then fill in)

> "First, let me clarify requirements:
>
> - **What data?** HTTP requests/responses, or also business events?
> - **Latency tolerance?** Can we block the API, or must be async?
> - **Retention?** How long must we keep data?
> - **Query patterns?** Real-time or batch analytics?
> - **Scale?** Events per second, data volume?
>
> Assuming: HTTP audit, async required, 7-year retention, batch analytics, 100 events/sec.
>
> **Capture layer:** Servlet Filter with ContentCachingWrapper. Async dispatch to avoid blocking.
>
> **Transport:** Kafka for durability and decoupling. Avro for schema enforcement and compression.
>
> **Storage:** Object storage (GCS/S3) for cost at scale. Parquet for compression and analytics compatibility.
>
> **Query:** BigQuery external tables or Athena for SQL access. No ETL needed.
>
> **Key decisions:**
> - Filter vs AOP → Filter for HTTP body access
> - Avro vs JSON → Avro for schema + compression
> - Kafka vs direct-to-storage → Kafka for durability + buffer
>
> **Trade-offs:**
> - Async means possible data loss → Mitigate with monitoring
> - Parquet means no real-time → Acceptable for compliance use case"

---

## Q8: How would you add real-time querying to your audit system?

**What they're testing:** Evolving systems, real-time vs batch understanding

### Model Answer (2 minutes)

> "Currently, our system optimizes for cost and batch analytics - Parquet files in GCS queried through BigQuery. Real-time would require architectural changes.
>
> **Option 1: Add a streaming consumer**
> Fork the Kafka topic to a new consumer that writes to Elasticsearch or OpenSearch. Keeps existing batch path intact. Downside: two systems to maintain, data in two places.
>
> **Option 2: Replace GCS with streaming database**
> Write to something like ClickHouse or Apache Druid that handles both real-time ingestion and analytics. Downside: more expensive, operational complexity.
>
> **Option 3: Kafka KSQL/Flink**
> Use stream processing for real-time aggregations and alerts. Keep GCS for historical. Good for 'alert me when...' use cases without full query capability.
>
> **My recommendation:** Option 1 for our use case. Suppliers mainly need debugging capability for recent requests - Elasticsearch with 30-day retention covers 95% of queries. Keep Parquet for long-term compliance. The complexity increase is manageable."

---

# SECTION 3: BEHAVIORAL QUESTIONS

---

## Q9: Tell me about a time you failed

**What they're testing:** Self-awareness, learning ability, honesty

### Model Answer (2 minutes)

> "The KEDA autoscaling incident was my mistake.
>
> I configured Kubernetes Event-Driven Autoscaling based on Kafka consumer lag. The logic seemed sound: more lag means more work, scale up workers to catch up.
>
> What I didn't anticipate was the feedback loop. When workers scaled up, Kafka triggered a consumer group rebalance. During rebalancing, no messages are processed, so lag increases. Increased lag triggers more scaling, which triggers more rebalancing. The system got stuck in this cycle.
>
> The root cause was that I tested autoscaling in isolation but not with production traffic patterns. My unit tests verified the scaling trigger worked, but I didn't simulate the rebalancing behavior.
>
> I fixed it by disabling KEDA and switching to CPU-based HPA, which is more predictable. The learning: autoscaling distributed systems requires understanding second-order effects. Now I always test scaling behavior with realistic load patterns, including failure modes."

### Follow-up Questions to Prepare

- "What did you learn?" → Test scaling with realistic traffic, understand second-order effects
- "How did you prevent this from happening again?" → Added to runbook, shared with team
- "Did anyone else make this mistake?" → Shared the pattern; one team avoided similar issue

---

## Q10: Tell me about a time you influenced without authority

**What they're testing:** Collaboration, communication, leadership

### Model Answer (2 minutes)

> "When building the audit logging library, I noticed three teams were building the same functionality independently. I had no authority over those teams, but I saw an opportunity to save everyone time.
>
> First, I scheduled meetings with the lead engineers. Critically, I didn't come with a solution - I came with questions. 'What are your requirements? What are your constraints?' This made them collaborators, not critics.
>
> There was resistance. 'Our needs are different.' One team wanted response body logging, another didn't. I addressed this by showing the common 80% and making the different 20% configurable. Response body logging became a config flag.
>
> I wrote comprehensive documentation - not just API docs, but a migration guide. Then I personally paired with each team on their integration PRs. When they found issues, I fixed them quickly, which built trust.
>
> Result: Three teams adopted in one month. Integration time went from two weeks of custom development to one day with the library. One engineer I helped later became an advocate and helped onboard a fourth team."

---

## Q11: Describe a time you received difficult feedback

**What they're testing:** Humility, growth mindset, emotional intelligence

### Model Answer (90 seconds)

> "During code review for the audit library, a senior engineer criticized my thread pool configuration. He pointed out that my queue size of 100 was arbitrary and could cause silent data loss when the queue filled up.
>
> My first instinct was defensive - I had thought about this! But I paused before responding.
>
> I ran the numbers. Each audit payload is about 2KB. Queue of 100 is 200KB - not a memory issue. But what about when it fills? The default `RejectedExecutionHandler` throws an exception, and since we catch all exceptions in the async method, the audit would be silently dropped.
>
> He was right. We could lose data with no indication.
>
> I added three things: a Prometheus metric for rejected tasks, a WARN log when queue exceeds 80%, and documentation explaining the trade-off. I thanked him on the PR and asked if I'd missed anything else.
>
> The library is more robust because of that review. We've actually had the 80% warning trigger once - caught it before it became critical."

---

## Q12: Why are you leaving Walmart?

**What they're testing:** Professionalism, motivation, red flags

### Model Answer (30-45 seconds)

> "I've had a great experience at Walmart. I've grown from building my first production system to leading major migrations and designing systems that serve millions of requests.
>
> I'm looking for my next challenge. Specifically, I'm excited about [COMPANY] because [SPECIFIC REASON - their scale, their product, their technical challenges].
>
> I want to continue growing as an engineer, and [COMPANY]'s [specific technical challenge or product area] is exactly the kind of problem I want to work on."

**Customize for:**
- **Rippling:** "excited about building HR/finance infrastructure that companies depend on daily"
- **Google:** "excited about the scale and impact - building systems that billions of people use"
- **Startup:** "excited about the pace and ownership - building from zero with direct customer impact"

---

## Q13: Tell me about a project you're most proud of

**What they're testing:** Passion, depth of contribution, impact awareness

### Model Answer (2-3 minutes)

> "The Kafka Audit Logging system, because it combined technical challenge with real business impact.
>
> **The challenge:** Splunk was being decommissioned. We needed a replacement that was cheaper, didn't impact API latency, and served both internal compliance needs and external supplier debugging.
>
> **My contribution:** I designed the three-tier architecture, built the core library, debugged a complex production incident over 5 days, and led the multi-region expansion for disaster recovery.
>
> **The technical depth:** I learned Kafka Connect internals, SMT filter development, Avro schema evolution, Parquet optimization. The debugging incident taught me systematic troubleshooting and the importance of observability.
>
> **The impact:** 2 million events per day, 99% cost reduction versus Splunk, less than 5ms latency impact. But what I'm most proud of is the supplier self-service feature - they can now debug their own API issues in 30 seconds instead of waiting 2 days for support.
>
> **The growth:** Three other teams adopted the library. It became the template for audit logging across the organization. I went from being the person asking questions to the person others ask."

---

# SECTION 4: GOOGLEYNESS / CULTURE FIT QUESTIONS

---

## Q14: Tell me about a time you helped someone else succeed

**What they're testing:** Collaboration, mentorship, team orientation

### Model Answer (90 seconds)

> "When the audit logging library was ready, I didn't just publish it and wait for adoption. I actively helped each team integrate.
>
> For one team, their lead engineer was skeptical - 'we've been burned by shared libraries before.' I spent an afternoon pairing with him on their PR. I showed him how to configure the endpoint filtering, helped debug a classloader issue, and wrote a custom extension for their specific use case.
>
> By the end, he understood the code deeply enough to maintain it himself. He later became an advocate and helped onboard a fourth team I hadn't even talked to.
>
> The lesson: adoption isn't just about building good software. It's about reducing friction and building trust through personal investment in others' success."

---

## Q15: Describe a time you had to make a decision with incomplete information

**What they're testing:** Decision-making, risk management, action bias

### Model Answer (2 minutes)

> "The multi-region rollout had vague requirements. Leadership said 'make it resilient' without specifying RTO, RPO, budget, or timeline. The only hard requirement was 'don't lose audit data.'
>
> Instead of waiting for clarity, I proactively defined requirements through stakeholder conversations. I asked: 'How much data loss is acceptable? How quickly must we recover?' I learned compliance needed max 4-hour RPO and 1-hour RTO.
>
> I designed three options with different trade-offs and presented them to the team. We chose Active/Active because audit is write-heavy and compliance couldn't tolerate 30-minute failover delays.
>
> The key decision under uncertainty was routing strategy. Without latency data, I had to choose between geography-based routing and round-robin. I chose geography-based because it's more deterministic and reduces cross-region traffic. This turned out to be correct.
>
> I documented my assumptions explicitly - 'I'm assuming RPO of 4 hours because...' - and shared with stakeholders. That document became the de facto requirements spec."

---

## Q16: How do you handle disagreements with teammates?

**What they're testing:** Conflict resolution, collaboration, ego management

### Model Answer (90 seconds)

> "I focus on understanding before advocating.
>
> During the audit library design, another engineer wanted to use AOP instead of Servlet Filters. My initial reaction was that AOP wouldn't work - it can't access the raw HTTP body. But instead of arguing, I asked: 'What's the benefit you're going for with AOP?'
>
> He explained: cleaner code, no filter chain complexity. Valid points. So I walked through the technical constraint - we need ContentCachingWrapper, which only works at the Servlet level. I showed him the code, not just the argument.
>
> We ended up with a hybrid: Servlet Filter for capture, but extracted the processing logic into a separate service class that was easier to test - addressing his maintainability concern.
>
> The lesson: disagreements are often about different priorities, not right vs wrong. Understanding their 'why' usually leads to a better solution than either original position."

---

## Q17: Tell me about a time you took initiative beyond your job description

**What they're testing:** Ownership, proactivity, impact orientation

### Model Answer (90 seconds)

> "The supplier self-service feature wasn't in my requirements. My task was to replace Splunk for internal logging. But while building it, I realized we could give suppliers direct access to their own data.
>
> I made technical choices that enabled this: Parquet format because BigQuery can query it directly, schema design that captured everything suppliers would need to debug, partition structure that supported row-level security.
>
> Then I worked with our data team to set up BigQuery external tables with security policies - each supplier only sees their own records.
>
> This wasn't asked for, but it had huge impact. Support ticket volume dropped. One supplier's engineer told me this was the first time they had visibility into their API interactions with any vendor.
>
> The lesson: when building infrastructure, always ask 'who else might benefit?' The primary use case drives requirements, but asking that question often reveals high-impact additions with minimal extra work."

---

# SECTION 5: ROLE-SPECIFIC QUESTIONS

---

## Q18: What's your approach to code reviews?

**What they're testing:** Quality standards, collaboration style, mentorship

### Model Answer (90 seconds)

> "I focus on three things: correctness, maintainability, and teaching.
>
> **Correctness:** Does it work? Are there edge cases? I look for error handling, null checks, concurrency issues.
>
> **Maintainability:** Will someone understand this in 6 months? I look for clear naming, appropriate abstractions, tests that document behavior.
>
> **Teaching:** Code review is a learning opportunity. When I see something I'd do differently, I explain why, not just what. 'Consider using X because Y' is better than 'Use X.'
>
> I also pick my battles. Not every style preference is worth a comment. If it works and it's readable, I approve even if I'd have written it differently.
>
> The feedback I received on my thread pool configuration taught me this: critical feedback is a gift. Now when I give feedback, I try to be as specific and constructive as that senior engineer was with me."

---

## Q19: How do you approach debugging production issues?

**What they're testing:** Systematic thinking, tooling knowledge, composure

### Model Answer (2 minutes)

> "Systematic isolation, starting from the highest level.
>
> **Step 1: Understand the symptom.** What exactly is broken? For the GCS sink issue, it was 'no new data in buckets.' Not 'Kafka down' or 'API failing' - specifically the sink output.
>
> **Step 2: Identify the boundary.** Where does the system work vs fail? I checked: Kafka has messages? Yes. Connect running? Yes. So the issue is between consumption and write.
>
> **Step 3: Increase observability.** Enable DEBUG logging, add metrics, trace a single message through. This revealed the NullPointerException in SMT filter.
>
> **Step 4: Hypothesis and test.** Form a hypothesis, make a single change, observe. The mistake I made was multiple changes at once early on - couldn't tell what helped.
>
> **Step 5: Correlate with changes.** What changed recently? Kubernetes events showed KEDA autoscaling - correlated with the issue timeline.
>
> **Step 6: Document and prevent.** Every issue becomes a runbook entry and potentially a new alert. The troubleshooting doc I wrote has been used twice since."

---

## Q20: How do you prioritize technical debt vs features?

**What they're testing:** Product sense, pragmatism, communication

### Model Answer (90 seconds)

> "I frame technical debt in terms of business impact, not engineering purity.
>
> 'This code is ugly' is not a compelling argument. 'This code causes 2 hours of debugging per week and will slow feature delivery by 30%' gets prioritized.
>
> For the Spring Boot 3 migration, the argument wasn't 'newer is better.' It was 'Snyk is flagging CVEs we cannot patch without upgrading, and we'll fail security audit in 3 months.'
>
> I also look for opportunities to pay down debt incrementally. When I'm in a file for a feature, I'll clean up related debt - better naming, adding tests, fixing warnings. Small improvements compound.
>
> But I'm pragmatic. Some debt is fine. If code works, is tested, and isn't changing, I don't refactor it just because I'd write it differently today. The question is always: 'What's the cost of leaving it vs the cost of fixing it?'"

---

# SECTION 6: RAPID FIRE QUESTIONS

Quick answers (30 seconds each)

---

## Q21: What's your biggest strength?

> "Systematic debugging. When something breaks, I stay calm and methodically isolate the issue. The 5-day Kafka Connect debugging is a good example - multiple root causes, but I traced each one systematically."

---

## Q22: What's your biggest weakness?

> "I can over-engineer solutions. I want to handle every edge case from day one. I've learned to start simpler and iterate - the audit library's initial version was more complex than necessary. Feedback from adopters helped me simplify."

---

## Q23: Where do you see yourself in 5 years?

> "Leading technical initiatives with broader scope - whether that's architecting systems used by multiple teams, mentoring engineers, or defining technical direction. I want to multiply my impact through others, not just my own code."

---

## Q24: What questions do you have for me?

**Always have 2-3 ready:**

> "What does success look like in the first 90 days for this role?"

> "What's the biggest technical challenge the team is facing right now?"

> "How does the team balance new features versus technical debt?"

> "Can you tell me about a recent technical decision the team made and how you reached it?"

---

# PRACTICE CHECKLIST

Before your interview:

- [ ] Can you explain the Kafka architecture in 3 minutes?
- [ ] Can you walk through the debugging incident step by step?
- [ ] Do you have specific metrics memorized (2M events, <5ms, 158 files)?
- [ ] Can you explain trade-offs for each major decision?
- [ ] Do you have 3 questions ready to ask them?
- [ ] Have you practiced answers out loud, not just read them?
- [ ] Can you explain technical concepts without jargon first?
- [ ] Do you use "I" not "we" for your contributions?

---

*Practice these out loud until they feel natural. Record yourself and listen back!*


# INTERVIEW MASTERCLASS: Kafka Audit Logging System & Spring Boot 3 Migration
## For Hiring Manager (Rippling) & Googleyness (Google) Rounds

> **Your Projects:**
> - Bullet 1, 2, 10 - Kafka-Based Audit Logging, Common Library JAR, Multi-Region Architecture
> - Bullet 4 - Spring Boot 3 & Java 17 Migration

---

# TABLE OF CONTENTS

1. [Understanding The Interview Rounds](#part-1-understanding-the-interview-rounds)
2. [The Onion Layer Framework](#part-2-the-onion-layer-framework)
3. [Project Explanation - Multiple Depths](#part-3-project-explanation---multiple-depths)
4. [Complete Q&A Bank with Follow-ups](#part-4-complete-qa-bank-with-follow-ups)
5. [Technical Deep Dive Trees](#part-5-technical-deep-dive-trees)
6. [Behavioral Questions & Stories](#part-6-behavioral-questions--stories)
7. [Curveball & What-If Scenarios](#part-7-curveball--what-if-scenarios)
8. [Mock Interview Scripts](#part-8-mock-interview-scripts)
9. [Red Flags & How to Redirect](#part-9-red-flags--how-to-redirect)
10. [Confidence Builders & Power Phrases](#part-10-confidence-builders--power-phrases)

---

# PART 1: UNDERSTANDING THE INTERVIEW ROUNDS

## Hiring Manager Round (Rippling)

### What They're Evaluating:

| Criteria | What They Look For | How to Demonstrate |
|----------|-------------------|-------------------|
| **Technical Depth** | Can you go beyond surface level? | Explain WHY decisions were made, not just WHAT |
| **Ownership** | Did you drive the project or just participate? | Use "I" for your contributions, "we" for team efforts |
| **Impact** | Did it matter to the business? | Quantify: events/day, latency reduction, cost savings |
| **Collaboration** | How did you work with others? | Mention stakeholders, cross-team dependencies |
| **Growth Mindset** | Did you learn and improve? | Share mistakes and what you learned |
| **Communication** | Can you explain complex things simply? | Start high-level, go deep when asked |

### Hiring Manager's Mental Model:

```
"Would I want this person on my team?"
     │
     ├── Can they do the job technically? (60%)
     │     ├── Relevant experience
     │     ├── Problem-solving ability
     │     └── Technical depth
     │
     ├── Will they fit the team? (25%)
     │     ├── Communication style
     │     ├── Collaboration signals
     │     └── Humility + confidence balance
     │
     └── Will they grow? (15%)
           ├── Learning from mistakes
           ├── Curiosity signals
           └── Self-awareness
```

---

## Googleyness Round (Google)

### What "Googleyness" Actually Means:

| Trait | Definition | How to Show It |
|-------|------------|----------------|
| **Doing the right thing** | Ethics, user focus | "We chose X because it was better for data privacy" |
| **Thriving in ambiguity** | Comfort with uncertainty | "Requirements were unclear, so I..." |
| **Valuing feedback** | Open to criticism | "My colleague pointed out a flaw, and I..." |
| **Challenging the status quo** | Questioning norms | "Everyone was doing X, but I proposed Y because..." |
| **Bringing others along** | Inclusive leadership | "I made sure junior devs understood by..." |
| **Putting users first** | User empathy | "Suppliers needed fast responses, so audit logging couldn't block..." |

### Googleyness Question Pattern:

```
"Tell me about a time when..."
     │
     ├── You disagreed with your team/manager
     ├── You had to make a decision with incomplete information
     ├── You received critical feedback
     ├── You helped someone else succeed
     ├── You had to prioritize competing demands
     └── You took a risk that didn't pay off
```

---

# PART 2: THE ONION LAYER FRAMEWORK

## How Interviewers Drill Down

Every topic has 5 layers. Interviewers start at Layer 1 and go deeper based on your answers.

```
┌─────────────────────────────────────────────────────────────────────┐
│  LAYER 1: WHAT (Surface)                                            │
│  "What did you build?"                                              │
│  → Audit logging system that captures API requests and stores in GCS│
├─────────────────────────────────────────────────────────────────────┤
│  LAYER 2: WHY (Purpose)                                             │
│  "Why did you need this?"                                           │
│  → Compliance requirements, debugging, analytics for 2M+ events/day │
├─────────────────────────────────────────────────────────────────────┤
│  LAYER 3: HOW (Implementation)                                      │
│  "How does it work?"                                                │
│  → Filter captures req/res → Async to Kafka → Connect writes to GCS │
├─────────────────────────────────────────────────────────────────────┤
│  LAYER 4: WHY THIS WAY (Decisions)                                  │
│  "Why did you choose this approach?"                                │
│  → Filter over AOP for body access, Async for non-blocking, etc.    │
├─────────────────────────────────────────────────────────────────────┤
│  LAYER 5: TRADE-OFFS (Mastery)                                      │
│  "What are the trade-offs? What would you do differently?"          │
│  → Memory overhead of caching, fire-and-forget loses some logs      │
└─────────────────────────────────────────────────────────────────────┘
```

## The Golden Rule

> **Answer at the level asked, but signal you can go deeper.**

Example:
- **Q:** "What did you build?"
- **A:** "I built a Kafka-based audit logging system that captures 2M+ API events daily and stores them in GCS for compliance. *Would you like me to walk through the architecture?*"

The italicized part signals depth without overwhelming.

---

# PART 3: PROJECT EXPLANATION - MULTIPLE DEPTHS

## 30-Second Pitch (Elevator)

> "When Walmart decommissioned Splunk, I designed and built a replacement audit logging system for our supplier APIs. The key constraint was that external suppliers like Pepsi and Coca-Cola needed to query their own API interaction history. I built a three-tier architecture: a reusable library that intercepts HTTP requests asynchronously, a Kafka publisher for durability, and a GCS sink that stores data in Parquet format - queryable via BigQuery. The system handles 2 million events daily with zero API latency impact, and suppliers can now self-serve their own debugging."

**Use when:** Initial introduction, time-constrained situations

**Key phrase to remember:** *"Splunk was being decommissioned, and suppliers needed query access to their data."*

---

## 2-Minute Explanation (Overview)

> "At Walmart Luminate, we serve external suppliers like Pepsi, Coca-Cola, and Unilever through APIs for inventory data and transaction history. Two things happened simultaneously: Walmart started decommissioning Splunk company-wide, and our suppliers asked for visibility into their API interactions - they wanted to see why their requests failed, debug issues themselves.
>
> So I needed to build a system that would:
> 1. Replace Splunk for our team's logging needs
> 2. Give suppliers query access to their own data
> 3. Do this without impacting API latency - suppliers expect sub-200ms responses
>
> I designed a three-tier solution:
>
> **First**, a reusable Spring Boot JAR that any service can import. It has a servlet filter that intercepts requests, caches the body using Spring's ContentCachingWrapper, and after the response, sends the audit payload asynchronously to our publisher service. Fire-and-forget - no latency impact.
>
> **Second**, a Kafka publisher service that receives these payloads, serializes them to Avro (70% smaller than JSON), and publishes to a multi-region Kafka cluster with routing headers.
>
> **Third**, a Kafka Connect sink with custom SMT filters that route US, Canada, and Mexico records to separate GCS buckets in Parquet format. BigQuery external tables sit on top - and that's what suppliers query.
>
> The result: 2 million events daily, P99 audit latency under 5ms, suppliers self-serve debugging, and we saved 99% compared to Splunk costs."

**Use when:** Project overview questions, "tell me about your work"

**Key differentiator:** *"Not just replacing Splunk - giving suppliers data access they never had before."*

---

## 5-Minute Deep Dive (Technical)

> "Let me walk you through the system end-to-end.
>
> **The Problem:**
> Two forces converged: Walmart was decommissioning Splunk enterprise-wide - it was expensive and the license was ending. Simultaneously, our external suppliers (Pepsi, Coca-Cola, Unilever) requested visibility into their API calls. They wanted to self-debug: 'Why did this request fail? Show me my last week's interactions.'
>
> We couldn't just use another log aggregator because suppliers needed SQL query access, not log grep. And we couldn't slow down APIs - SLAs require sub-200ms responses.
>
> **Component 1: Common Library (dv-api-common-libraries)**
>
> I built this as a Spring Boot starter JAR. The core is a servlet filter called LoggingFilter, annotated with @Order(LOWEST_PRECEDENCE) so it runs last in the filter chain - after auth, after validation. This ensures we capture the final state.
>
> The filter wraps the request and response in ContentCachingWrapper. Why? Because HTTP bodies are streams - you can only read them once. Without caching, either the filter or the controller could read the body, but not both.
>
> After the controller processes the request, the filter extracts the cached bodies, builds an AuditLogPayload with request ID, timestamps, headers, bodies, and response code. Then it calls AuditLogService.sendAuditLogRequest() - which is annotated with @Async.
>
> The @Async annotation is crucial. It means the audit call happens in a separate thread pool - 6 core threads, max 10, queue of 100. The API response returns immediately to the user. If audit logging fails, we log the error but don't impact the user.
>
> **Component 2: Publisher Service (audit-api-logs-srv)**
>
> This is a simple Spring Boot service with one endpoint: POST /v1/logRequest. It receives the payload, converts it to Avro using a schema in Confluent Schema Registry, and publishes to Kafka.
>
> Why Avro? Three reasons: schema enforcement prevents bad data, binary format is 70% smaller than JSON, and schema evolution lets us add fields without breaking consumers.
>
> The key detail: we forward the wm-site-id header to the Kafka message. This enables geographic routing downstream.
>
> **Component 3: GCS Sink (audit-api-logs-gcs-sink)**
>
> This is a Kafka Connect worker running three sink connectors in parallel - one each for US, Canada, and Mexico.
>
> Each connector has a custom Single Message Transform that filters based on the wm-site-id header. The US filter is special - it's permissive, accepting records with US header OR records with no header at all. This handles legacy data. Canada and Mexico filters are strict.
>
> The connector writes to GCS in Parquet format, partitioned by service name, date, and endpoint. We chose Parquet because it compresses to 10% of JSON size and BigQuery reads it natively for analytics.
>
> **The Result:**
> - 2M+ events/day processed
> - P99 audit latency: <5ms
> - Zero impact on API response time
> - Full data residency compliance
> - 90% storage cost reduction (Parquet vs JSON)"

**Use when:** Technical deep dive requests, architecture discussions

---

## 10-Minute Full Story (With Challenges)

*[Include everything above, plus:]*

> **Challenges I Faced:**
>
> **Challenge 1: The Silent Failure**
>
> Two weeks after launch, we noticed GCS buckets weren't receiving new data. No errors in logs. This was terrifying - we were losing audit data.
>
> I spent two days debugging. First, I added try-catch blocks in the SMT filter - found NullPointerException on records with null headers. Fixed that, but still had issues.
>
> Next, I noticed consumer poll timeouts. Kafka Connect default is 5 minutes, but our processing was slow. I tuned max.poll.interval.ms and related configs.
>
> Then I discovered KEDA autoscaling was causing problems. Every scale-up triggered consumer group rebalancing, which caused more timeouts, which triggered more scaling. A feedback loop. I disabled autoscaling temporarily and added proper scaling thresholds.
>
> Finally, I found heap exhaustion. Large batches were causing OOM. Increased JVM heap from 512MB to 2GB.
>
> This debugging story taught me that Kafka Connect's error tolerance can hide problems, and that distributed systems fail in unexpected ways.
>
> **Challenge 2: Multi-Region Rollout**
>
> We needed to expand from single region to dual-region for disaster recovery. But we couldn't have downtime - this is compliance-critical data.
>
> I designed a phased approach:
> 1. First, modified the publisher to write to both Kafka clusters
> 2. Added the wm-site-id header for routing
> 3. Deployed GCS sink to second region
> 4. Validated data parity between regions
> 5. Cut over traffic gradually
>
> Zero downtime, zero data loss.
>
> **What I Learned:**
>
> 1. Async processing is powerful but debugging is hard - invest in observability upfront
> 2. Default configurations are rarely production-ready - always tune
> 3. Multi-region is complex - plan for data consistency
> 4. Reusable libraries save enormous time - other teams adopted our JAR in days"

**Use when:** Behavioral questions about challenges, full project narratives

---

# PART 3.5: THE BUSINESS VALUE STORY (NEW!)

## Why This Project Is Interview Gold

This project has THREE compelling narratives:

### 1. **Crisis Response** (Splunk Decommissioning)
> "Walmart was removing Splunk. We had X weeks to replace it or lose all logging capability."

This shows: **Urgency handling, working under pressure, time constraints**

### 2. **Customer-Centric Innovation** (Supplier Access)
> "Suppliers couldn't debug their own issues. They'd call our support team. I gave them direct access to their data."

This shows: **User empathy, solving real problems, thinking beyond the obvious**

### 3. **Technical Excellence** (Zero Latency Impact)
> "I designed it so 2M events/day are processed with zero impact on API response time."

This shows: **Performance awareness, async patterns, production thinking**

## Supplier Access Architecture

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
│   Supplier sees: "Oh, my token expired and I had an invalid store_id"          │
│   Before: They would call support, wait 2 days, we'd grep logs                │
│   After: Self-service in 30 seconds                                           │
└────────────────────────────────────────────────────────────────────────────────┘
```

## Questions About Supplier Access

**Q: "How do suppliers access their data?"**
> "BigQuery external tables point to GCS Parquet files. Suppliers get read-only access scoped to their consumer_id. They can run SQL queries directly - 'show me all failed requests from last week' - without involving our team."

**Q: "How do you handle data isolation?"**
> "Row-level security in BigQuery. Each supplier's access policy filters by their WM_CONSUMER.ID header. Pepsi can only see Pepsi's data, even though it's all in the same bucket."

**Q: "What queries do suppliers typically run?"**
> "Three main patterns:
> 1. Debug failed requests: `WHERE response_code >= 400`
> 2. Latency analysis: `AVG(response_ts - request_ts)`
> 3. Usage patterns: `COUNT(*) GROUP BY endpoint, date`"

**Q: "How does this differ from just giving them logs?"**
> "Logs require grep and regex. This is structured data with SQL. A supplier can answer 'what percentage of my requests failed yesterday' in one query. In Splunk or ELK, that's a complex SPL/KQL query. In BigQuery, it's `SELECT COUNT(*) FILTER(WHERE response_code >= 400) / COUNT(*)`."

---

# PART 3.6: BULLET 2 - COMMON LIBRARY JAR (Dedicated Section)

## Why This Deserves Its Own Story

The **dv-api-common-libraries** JAR is a separate resume bullet because it demonstrates:
- **Reusability thinking** - Built once, used by many teams
- **API design skills** - Zero-code integration for consumers
- **Spring Boot internals** - Filters, auto-configuration, async processing

## The 30-Second Pitch for the Library Alone

> "I built a reusable Spring Boot starter JAR that gives any service audit logging by just adding a Maven dependency. No code changes needed. It uses a servlet filter with ContentCachingWrapper to capture HTTP bodies, and @Async processing so it doesn't impact API latency. Three teams adopted it within a month, reducing integration time from 2 weeks to 1 day."

## Technical Deep Dive: How It Works

### The Magic: Zero-Code Integration

```xml
<!-- That's ALL a consuming service needs to add -->
<dependency>
    <groupId>com.walmart</groupId>
    <artifactId>dv-api-common-libraries</artifactId>
    <version>0.0.45</version>
</dependency>
```

**How does it work without code changes?**

```
┌────────────────────────────────────────────────────────────────────────────────┐
│                    SPRING BOOT AUTO-CONFIGURATION MAGIC                         │
│                                                                                 │
│  1. @ComponentScan picks up LoggingFilter (it's a @Component)                  │
│                                                                                 │
│  2. @EnableAsync + @Bean registers ThreadPoolTaskExecutor                      │
│                                                                                 │
│  3. @ManagedConfiguration loads CCM config (feature flags, endpoints)          │
│                                                                                 │
│  4. Spring's filter chain auto-registers LoggingFilter                         │
│                                                                                 │
│  Result: Service starts → Filter active → Audit logging works!                 │
└────────────────────────────────────────────────────────────────────────────────┘
```

### Key Design Decisions

| Decision | Why | Alternative Considered |
|----------|-----|----------------------|
| **Servlet Filter** | Access to raw HTTP stream | AOP (no body access) |
| **@Order(LOWEST_PRECEDENCE)** | Run AFTER security filters | Higher priority (miss auth failures) |
| **ContentCachingWrapper** | Read body without consuming stream | Custom buffering (complex) |
| **@Async** | Non-blocking, fire-and-forget | Sync (blocks API response) |
| **CCM config for endpoints** | Runtime enable/disable without deploy | Hardcoded list (inflexible) |
| **ThreadPool (6 core, 10 max, 100 queue)** | Handle bursts, bounded memory | Unbounded (OOM risk) |

### The ContentCachingWrapper Problem & Solution

**The Problem:**
```java
// HTTP body is a stream - can only read ONCE
String body = request.getReader().lines().collect(joining());  // Works
String body2 = request.getReader().lines().collect(joining()); // EMPTY!
```

**The Solution:**
```java
// Wrap the request - now body is cached
ContentCachingRequestWrapper wrapped = new ContentCachingRequestWrapper(request);
filterChain.doFilter(wrapped, response);  // Controller reads body
byte[] body = wrapped.getContentAsByteArray();  // Filter reads SAME body
```

### Interview Questions About the Library

**Q: "Why a library instead of a sidecar?"**
> "Sidecars like Envoy work at the network layer - they can't see application context. We needed to capture which endpoint was called semantically (not just the URL), what the error message meant, and which supplier made the request. That context only exists inside the application."

**Q: "How did you get other teams to adopt it?"**
> "Three steps: First, I met with engineers from other teams to understand their requirements and incorporated their needs. Second, I wrote comprehensive documentation and a migration guide. Third, I did a brown-bag session and personally helped each team integrate - spent an afternoon pairing on their PRs. Three teams adopted within a month."

**Q: "What was the hardest design decision?"**
> "The thread pool sizing. A senior engineer challenged my queue size of 100 during code review - said it could cause silent data loss. He was right. I added metrics for rejected tasks and a warning when queue is 80% full. That monitoring has actually caught downstream slowdowns twice."

**Q: "What happens if the audit service is down?"**
> "API continues normally. We catch all exceptions in AuditLogService and log them, but don't propagate. The thread pool queue absorbs brief outages. For extended outages, we lose audit data but users are unaffected. This is intentional - audit is best-effort, not critical path."

**Q: "How do teams configure which endpoints to audit?"**
> "CCM (Cloud Configuration Management). Each team has their own CCM config with a list of endpoint patterns. They can add/remove endpoints without code changes or redeploys. We also support regex patterns for flexible matching."

### Code Snippets to Know

**LoggingFilter - The Core:**
```java
@Component
@Order(Ordered.LOWEST_PRECEDENCE)  // Run LAST
public class LoggingFilter extends OncePerRequestFilter {

    @ManagedConfiguration
    FeatureFlagCCMConfig featureFlagCCMConfig;  // Feature flag

    @Autowired
    AuditLogService auditLogService;  // Async sender

    @Override
    protected void doFilterInternal(HttpServletRequest request,
                                    HttpServletResponse response,
                                    FilterChain chain) {
        if (!featureFlagCCMConfig.isAuditLogEnabled()) {
            chain.doFilter(request, response);
            return;
        }

        // Wrap to capture body
        ContentCachingRequestWrapper reqWrapper =
            new ContentCachingRequestWrapper(request);
        ContentCachingResponseWrapper resWrapper =
            new ContentCachingResponseWrapper(response);

        long startTime = Instant.now().getEpochSecond();
        chain.doFilter(reqWrapper, resWrapper);  // Controller runs
        long endTime = Instant.now().getEpochSecond();

        // Build and send audit (async)
        AuditLogPayload payload = buildPayload(reqWrapper, resWrapper,
                                               startTime, endTime);
        auditLogService.sendAuditLogRequest(payload);  // @Async!

        resWrapper.copyBodyToResponse();  // Don't forget this!
    }
}
```

**AuditLogService - Async Sender:**
```java
@Service
public class AuditLogService {

    @Async  // Runs in ThreadPoolTaskExecutor
    public void sendAuditLogRequest(AuditLogPayload payload) {
        try {
            httpService.post(auditServiceUrl, payload);
        } catch (Exception e) {
            log.error("Audit failed for {}", payload.getRequestId(), e);
            // Don't throw - audit failure shouldn't break API
        }
    }
}
```

**Thread Pool Config:**
```java
@Configuration
@EnableAsync
public class AuditLogAsyncConfig {
    @Bean
    public Executor taskExecutor() {
        ThreadPoolTaskExecutor executor = new ThreadPoolTaskExecutor();
        executor.setCorePoolSize(6);      // Always running
        executor.setMaxPoolSize(10);      // Can scale to
        executor.setQueueCapacity(100);   // Buffer
        executor.setThreadNamePrefix("Audit-");
        return executor;
    }
}
```

### The Adoption Story (Behavioral)

**STAR Format:**

**Situation:** "Our team needed audit logging, but I noticed two other teams building the same thing independently."

**Task:** "I proposed building a shared library and took responsibility for getting buy-in and driving adoption."

**Action:**
- Met with engineers from other teams to understand requirements
- Made the 20% different requirements configurable (response body logging toggle, endpoint filtering)
- Wrote documentation and migration guide
- Did a brown-bag session demo
- Personally helped each team integrate (paired programming)

**Result:** "Three teams adopted within a month. Integration time dropped from 2 weeks to 1 day. The library became the standard for all new services."

### Numbers to Remember

| Metric | Value |
|--------|-------|
| Lines of code | ~500 (entire library) |
| Integration time | 1 day (vs 2 weeks custom) |
| Teams using it | 3+ |
| PRs in library | 16 |
| Thread pool | 6 core, 10 max, 100 queue |
| Maven version | 0.0.45 |

---

# PART 3.7: BULLET 4 - SPRING BOOT 3 & JAVA 17 MIGRATION (Dedicated Section)

## Why This Migration Is Interview Gold

The Spring Boot 3 migration demonstrates:
- **Technical leadership** - Led a major framework upgrade
- **Risk management** - Zero-downtime deployment strategy
- **Deep framework knowledge** - javax→jakarta, Hibernate 6, WebClient
- **Modern Java skills** - Java 17 features, CompletableFuture

## The 30-Second Pitch

> "I led the migration of cp-nrti-apis from Spring Boot 2.7 to 3.2 and Java 11 to 17. The main challenges were the javax→jakarta namespace change across 74 files, migrating from RestTemplate to WebClient, and adapting to Hibernate 6's stricter type handling. We deployed to production with zero downtime using Flagger canary deployments - 10% traffic initially, gradually increasing over 24 hours with automatic rollback if error rate exceeded 1%."

## The 2-Minute Explanation

> "Spring Boot 2.7 was approaching end-of-life, and Snyk was flagging security vulnerabilities in our dependencies. I led the migration for our main API service.
>
> Three main technical challenges:
>
> **First**, the javax→jakarta namespace change. Java EE became Jakarta EE, so every import of `javax.persistence`, `javax.validation`, `javax.servlet` had to change. 74 files across entities, controllers, validators, and filters.
>
> **Second**, RestTemplate deprecation. Spring Boot 3 pushes you toward WebClient, which is reactive. Our codebase was synchronous, so I decided to use `.block()` for backwards compatibility rather than rewriting everything to be reactive - that would be a separate initiative.
>
> **Third**, Hibernate 6 compatibility. PostgreSQL enum types needed explicit `@JdbcTypeCode` annotations. Without them, Hibernate tried to use VARCHAR and PostgreSQL rejected it.
>
> For deployment, we used Flagger canary strategy - 10% traffic to the new version, monitored for 4 hours, then gradually increased. Zero customer-impacting issues."

## Key Technical Decisions

| Decision | Why | Alternative Considered |
|----------|-----|----------------------|
| **Use .block() with WebClient** | Maintain sync behavior, minimize scope | Full reactive rewrite (too risky) |
| **Java 17 not 21** | LTS version, Walmart standard | Java 21 was too new |
| **Canary deployment** | Automatic rollback on errors | Big bang (too risky) |
| **Manual migration over OpenRewrite** | Team wanted to understand each change | OpenRewrite automation |

## Interview Questions

**Q: "Why not go fully reactive with WebClient?"**
> "Scope control and risk management. A full reactive migration would change every service class and fundamentally change error handling. For a framework upgrade, I wanted to minimize changes to business logic. Reactive would be a dedicated project with proper training."

**Q: "What was the hardest part?"**
> "The WebClient test mocking. RestTemplate is simple - you mock `exchange()` and you're done. WebClient uses method chaining: `webClient.get().uri().headers().retrieve().bodyToMono()`. Each method in the chain needs to be mocked. Test files doubled in complexity."

**Q: "Any production issues?"**
> "No customer-impacting issues. Three minor post-migration fixes: Snyk vulnerabilities in transitive dependencies, one Sonar code smell around mutable collections, and a logging format adjustment. All caught in monitoring and fixed proactively."

**Q: "How did you handle the javax→jakarta change?"**
> "IDE refactoring for bulk changes, but manual verification for servlet-related code. Some classes changed behavior, not just package names. We added comprehensive tests to ensure filter ordering was preserved."

## Code Patterns to Know

**WebClient with .block():**
```java
return webClient.get()
    .uri(url)
    .headers(h -> h.addAll(headers))
    .retrieve()
    .bodyToMono(StoreResponse.class)
    .block();  // Blocking for sync compatibility
```

**Hibernate 6 Enum Handling:**
```java
@Column(name = "status", columnDefinition = "status_enum")
@Enumerated(EnumType.STRING)
@JdbcTypeCode(SqlTypes.NAMED_ENUM)  // Required for PostgreSQL enums
private Status status;
```

**CompletableFuture for Kafka:**
```java
kafkaTemplate.send(message)
    .thenAccept(result -> log.info("Success: {}", result))
    .exceptionally(ex -> {
        handleFailure(message);
        return null;
    })
    .join();
```

## Numbers to Remember

| Metric | Value |
|--------|-------|
| Files changed | 136 |
| javax→jakarta files | 74 |
| Lines added | 1,732 |
| Spring Boot version | 2.7 → 3.2 |
| Java version | 11 → 17 |
| Production issues | 0 customer-impacting |
| Deployment strategy | Canary: 10%→25%→50%→100% |

## The Behavioral Story (STAR)

**Situation:** "Spring Boot 2.7 was approaching end-of-life. Security scans were flagging vulnerabilities. We needed to migrate before EOL."

**Task:** "I volunteered to lead the migration for our main API service, cp-nrti-apis, which handles supplier API requests."

**Action:**
- Analyzed migration path, identified 3 main challenges (namespace, WebClient, Hibernate)
- Made strategic decision to use .block() instead of full reactive rewrite
- Updated 158 files, wrote comprehensive tests
- Implemented canary deployment with automatic rollback
- Staged for 1 week, then 24-hour gradual production rollout

**Result:** "Zero customer-impacting issues. Migration completed in 4 weeks. Three minor post-migration fixes, all caught proactively. Set the template for other team migrations."

---

# PART 4: COMPLETE Q&A BANK WITH FOLLOW-UPS

## Opening Questions

### Q1: "Tell me about this audit logging system you built."

**Level 1 Answer:**
> "When Walmart decommissioned Splunk, I built a replacement audit logging system. But I went beyond just logging - I designed it so external suppliers could query their own API interaction data. It handles 2 million events daily with zero API latency impact."

**Expected Follow-up:** "Why was Splunk being decommissioned?"

**Level 1.5 Answer:**
> "Cost and licensing. Splunk is expensive at scale, and Walmart decided not to renew the enterprise license. Every team needed to find alternatives. For us, it was an opportunity to build something better."

**Expected Follow-up:** "Can you walk me through the architecture?"

**Level 2 Answer:**
> "Sure. It's a three-component system. First, a reusable Spring Boot library that other services import - it has a filter that intercepts HTTP requests. Second, a Kafka publisher service that serializes to Avro and publishes to a multi-region Kafka cluster. Third, a Kafka Connect sink with custom filters that route records to geographic buckets in GCS. BigQuery sits on top for SQL queries - that's how suppliers access their data."

**Expected Follow-up:** "Why did you choose this architecture?"

**Level 3 Answer:**
> "The key constraint was supplier access. They needed SQL queries, not log grep. That ruled out traditional log aggregators. I chose:
> - GCS with Parquet because it's cheap at scale and BigQuery reads it natively
> - Kafka for durability - if GCS sink fails, no data loss
> - A library approach so teams get audit logging by adding a Maven dependency
> - Async processing because we couldn't impact API latency"

**Expected Follow-up:** "What were the alternatives you considered?"

**Level 4 Answer:**
> "For storage: Elasticsearch - great for search but expensive and suppliers would need Kibana training. Datadog - cost prohibitive at 2M events/day. Cloud Logging - not queryable by external users.
>
> For interception: AOP - doesn't give clean access to HTTP body streams. Sidecar proxy - doesn't have application context like which endpoint was called.
>
> For pipeline: Custom Kafka consumer - but Kafka Connect handles offset management and retries out of the box."

**Expected Follow-up:** "What would you do differently today?"

**Level 5 Answer:**
> "Three things. First, I'd add structured logging with OpenTelemetry from day one - debugging the silent failure would have been faster. Second, I'd skip the publisher service and publish directly to Kafka from the library - removes a hop. Third, I'd use Apache Iceberg instead of raw Parquet - it gives ACID transactions and easier GDPR deletes."

---

### Q2: "Why did you build this as a library instead of a sidecar?"

**Answer:**
> "Good question. Sidecars like Envoy can do request logging, but they operate at the network layer - they don't have application context. We needed to capture business-level information like which endpoint was called, what the semantic error was, and which supplier made the request. That information is only available inside the application.
>
> Also, our services already had CCM (configuration management) integration. A library could leverage that for feature flags and endpoint filtering. A sidecar would need its own configuration mechanism."

**Follow-up:** "What about performance impact of running in-process?"

**Answer:**
> "The async design minimizes impact. The filter adds maybe 1ms to capture the body. The actual HTTP call to the publisher happens in a separate thread pool after the response is sent. We measured P99 impact at under 5ms. And if the publisher is slow or down, the API still responds normally - audit is fire-and-forget."

**Follow-up:** "Fire-and-forget means you can lose audit logs?"

**Answer:**
> "Yes, that's a trade-off we made consciously. The alternative - making audit synchronous - would mean a publisher outage causes API failures. For compliance, we mitigate this with: thread pool queue (buffers 100 requests), publisher retry logic, and Kafka durability. In practice, we've lost less than 0.01% of events."

---

### Q3: "How did you handle the HTTP body stream problem?"

**Answer:**
> "HTTP request and response bodies are input/output streams in Java - you can only read them once. If the filter reads the body, the controller gets empty input.
>
> I used Spring's ContentCachingRequestWrapper and ContentCachingResponseWrapper. These classes read the stream once and cache the bytes. The filter lets the chain proceed normally, then reads from the cache after the controller is done.
>
> The trade-off is memory - we're keeping the entire body in memory. For large payloads, this could be problematic. We mitigate by having payload size limits at the gateway level (2MB max) and an option to skip response body logging for certain endpoints."

**Follow-up:** "What if someone uploads a 2GB file?"

**Answer:**
> "Gateway rejects it before it reaches our service. But even if it got through, ContentCachingWrapper has a configurable limit. Beyond that limit, it stops caching. We'd audit the request metadata without the body, which is still useful for compliance."

---

## Technical Deep Dive Questions

### Q4: "Explain the Kafka Connect SMT filter design."

**Answer:**
> "SMT stands for Single Message Transform - it's Kafka Connect's plugin system for processing messages inline.
>
> I created an abstract base class called BaseAuditLogSinkFilter that implements the Transformation interface. It has one abstract method: getHeaderValue(), which subclasses override to return their target site ID.
>
> The apply() method checks if the record's wm-site-id header matches. If yes, return the record. If no, return null - which tells Kafka Connect to drop it.
>
> The clever part is the US filter. Legacy records don't have the wm-site-id header. So the US filter is permissive - it accepts records WITH the US header OR records with NO header. Canada and Mexico filters are strict - exact match only.
>
> This design means all records get processed, nothing is lost, and data ends up in the right geographic bucket."

**Follow-up:** "Why not filter at the producer side?"

**Answer:**
> "Two reasons. First, separation of concerns - the producer's job is to publish reliably, not to know about downstream routing. Second, flexibility - if we add a new country, we just deploy a new connector. No changes to producers."

**Follow-up:** "What happens if a record has an unknown site ID?"

**Answer:**
> "Currently it would be dropped by all three connectors. In practice, we validate site ID at the API layer before publishing. But you're right - we should have a catch-all connector that routes unknowns to a 'review' bucket. That's on our roadmap."

---

### Q5: "How does the async processing work in detail?"

**Answer:**
> "Spring's @Async annotation is the key. When you annotate a method with @Async, Spring creates a proxy that submits the method call to an Executor instead of running it directly.
>
> I configured a ThreadPoolTaskExecutor with:
> - corePoolSize: 6 (threads always running)
> - maxPoolSize: 10 (can scale to this)
> - queueCapacity: 100 (buffer before rejection)
>
> When LoggingFilter calls auditLogService.sendAuditLogRequest(), the proxy submits the call to the executor. The filter returns immediately. The executor runs the method when a thread is available.
>
> If all 10 threads are busy and the queue is full (100 items), the default RejectedExecutionHandler throws an exception. But our sendAuditLogRequest catches all exceptions, so the request is silently dropped. This is intentional - we never want audit to fail the API."

**Follow-up:** "Why those specific thread pool sizes?"

**Answer:**
> "Based on traffic analysis. Our APIs handle about 100 requests per second per pod. Each audit call takes about 50ms. So we need 5 threads to keep up (100 * 0.05). I set 6 core for headroom. Max 10 handles spikes. Queue of 100 handles burst traffic without creating too many threads."

**Follow-up:** "What if you get 1000 requests per second?"

**Answer:**
> "The queue would fill up and requests would be rejected. We'd see this in our metrics (rejected_count) and would scale up pods. In practice, we haven't hit this because Kubernetes HPA scales pods before queue pressure builds."

---

### Q6: "Walk me through what happens when a request comes in."

**Complete Flow Answer:**

> "Let me trace a single request through the entire system.
>
> **Step 1: Request Arrives**
> A supplier sends POST /iac/v1/inventory with a JSON body. The request hits our Kubernetes ingress, goes through Istio sidecar for mTLS, and reaches the cp-nrti-apis pod.
>
> **Step 2: Filter Chain Begins**
> Spring's DispatcherServlet starts the filter chain. Filters run in order: CORS, Security, then our LoggingFilter (at LOWEST_PRECEDENCE, so it's last).
>
> **Step 3: LoggingFilter - Before**
> LoggingFilter.doFilterInternal() is called. It checks the feature flag (CCM config) - if audit is disabled, it just passes through. If enabled, it wraps the request and response in ContentCachingWrapper and records the timestamp.
>
> **Step 4: Controller Executes**
> filterChain.doFilter() is called. The controller processes the request - validates input, calls downstream services, builds response. Controller returns ResponseEntity.
>
> **Step 5: LoggingFilter - After**
> Control returns to LoggingFilter. It records the response timestamp. It extracts the cached request body (contentCachingRequestWrapper.getContentAsByteArray()) and response body. It builds AuditLogPayload with all the data.
>
> **Step 6: Async Handoff**
> auditLogService.sendAuditLogRequest(payload) is called. Because of @Async, this immediately returns. The actual work is queued to the thread pool. LoggingFilter calls contentCachingResponseWrapper.copyBodyToResponse() to send the response to the client.
>
> **Step 7: Client Gets Response**
> The response goes back through Istio, through ingress, to the supplier. Total time: ~150ms. Audit work is still happening in background.
>
> **Step 8: Audit HTTP Call**
> In the thread pool, AuditLogService makes an HTTP POST to audit-api-logs-srv with the payload. Headers include WM_CONSUMER.ID and a signature generated from our private key.
>
> **Step 9: Kafka Publishing**
> audit-api-logs-srv receives the payload. AuditLoggingController calls LoggingRequestService which calls KafkaProducerService. The payload is converted to Avro, headers are set (including wm-site-id), and kafkaTemplate.send() publishes to api_logs_audit_prod topic.
>
> **Step 10: Kafka Persistence**
> Kafka broker receives the message, replicates to followers, acks to producer. Message is now durable.
>
> **Step 11: Connect Consumption**
> Kafka Connect worker polls the topic. It receives our message along with others in a batch (max 50 per poll).
>
> **Step 12: SMT Filtering**
> For each message, the configured SMT chain runs. InsertRollingRecordTimestamp adds a date header. Then FilterUS/CA/MX checks wm-site-id. Our message has US, so FilterUS passes it through, FilterCA and FilterMX return null.
>
> **Step 13: GCS Write**
> The US connector's GCPStorageSinkConnector batches records and writes to gs://audit-api-logs-us-prod/. The path includes service name, date, and endpoint for partitioning. Format is Parquet.
>
> **Step 14: Commit**
> After successful write, connector commits offsets to Kafka. If we restart, we resume from here.
>
> **Step 15: BigQuery Access**
> An external table in BigQuery points to the GCS bucket. Analysts can now query with SQL:
> ```sql
> SELECT * FROM audit_logs.api_logs_us
> WHERE DATE(created_ts) = CURRENT_DATE()
> ```
>
> **Total latency:** Client sees 150ms. Full audit pipeline completes in ~5 seconds."

---

# PART 5: TECHNICAL DEEP DIVE TREES

## Question Tree: Kafka Connect

```
"Why Kafka Connect?"
│
├─► "What alternatives did you consider?"
│   │
│   ├─► Custom Consumer
│   │   └─► "Why not that?" → Offset management, scaling, error handling all manual
│   │
│   ├─► Apache Flink
│   │   └─► "Why not that?" → Overkill for simple sink, operational complexity
│   │
│   └─► Cloud Functions (trigger on Kafka)
│       └─► "Why not that?" → Cold starts, cost at scale, no built-in batching
│
├─► "How does Connect handle failures?"
│   │
│   ├─► "What if GCS is down?"
│   │   └─► RETRY policy, exponential backoff, 5 retries, then DLQ
│   │
│   ├─► "What if message is malformed?"
│   │   └─► errors.tolerance: all → routes to DLQ with context headers
│   │
│   └─► "What if Connect worker crashes?"
│       └─► K8s restarts pod, Connect resumes from last committed offset
│
├─► "How does Connect scale?"
│   │
│   ├─► "Horizontally?"
│   │   └─► Add more workers, partitions distributed automatically
│   │
│   ├─► "Vertically?"
│   │   └─► Increase tasks.max, but we keep it at 1 for ordering
│   │
│   └─► "Auto-scaling?"
│       └─► We tried KEDA on lag, caused rebalance storms, now manual
│
└─► "What's the performance?"
    │
    ├─► "Throughput?"
    │   └─► ~10K messages/sec per task, we do 23/sec average
    │
    ├─► "Latency?"
    │   └─► Kafka to GCS: ~2 seconds in batches, not real-time critical
    │
    └─► "Cost?"
        └─► Connect is open source, only GCS storage cost (~$500/month)
```

## Question Tree: Avro

```
"Why Avro over JSON?"
│
├─► "What's the size difference?"
│   └─► ~70% smaller, saves Kafka storage and network
│
├─► "What about schema evolution?"
│   │
│   ├─► "Can you add fields?"
│   │   └─► Yes, with default values, backward compatible
│   │
│   ├─► "Can you remove fields?"
│   │   └─► Yes, if optional, forward compatible
│   │
│   └─► "What if schema is incompatible?"
│       └─► Schema Registry rejects, producer fails fast
│
├─► "What's Schema Registry?"
│   │
│   ├─► "How does it work?"
│   │   └─► Central service, stores schemas by subject, assigns IDs
│   │
│   ├─► "Message format?"
│   │   └─► [Magic Byte][4-byte Schema ID][Avro Payload]
│   │
│   └─► "What if Registry is down?"
│       └─► Producer caches schemas, continues for cached, fails for new
│
└─► "Alternatives to Avro?"
    │
    ├─► "Protobuf?"
    │   └─► Similar benefits, we chose Avro for Kafka ecosystem fit
    │
    ├─► "JSON Schema?"
    │   └─► Less compact, validation only, not serialization
    │
    └─► "Parquet directly?"
        └─► Not for streaming, Parquet is columnar for batch
```

## Question Tree: GCS + Parquet

```
"Why GCS with Parquet?"
│
├─► "Why not just JSON files?"
│   │
│   ├─► "Size?"
│   │   └─► Parquet is 10x smaller with compression
│   │
│   └─► "Query performance?"
│       └─► Parquet is columnar, only reads needed columns
│
├─► "Why not a database?"
│   │
│   ├─► "PostgreSQL?"
│   │   └─► Expensive at 2M rows/day, limited retention
│   │
│   ├─► "BigQuery directly?"
│   │   └─► Streaming inserts are expensive, batched is cheaper
│   │
│   └─► "Elasticsearch?"
│       └─► Great for search, expensive for compliance archival
│
├─► "How does BigQuery read GCS?"
│   │
│   ├─► "External tables?"
│   │   └─► Yes, points to GCS path, no data copy
│   │
│   ├─► "Performance?"
│   │   └─► Scans files on query, partitioning by date helps
│   │
│   └─► "Cost?"
│       └─► Pay for bytes scanned, Parquet reduces this
│
└─► "Data lifecycle?"
    │
    ├─► "How long do you keep data?"
    │   └─► 7 years for compliance
    │
    ├─► "Do you tier storage?"
    │   └─► After 30 days → Nearline, after 1 year → Coldline
    │
    └─► "Can you delete specific records?"
        └─► Yes but expensive, GDPR requests require rewriting files
```

---

# PART 6: BEHAVIORAL QUESTIONS & STORIES

## STAR++ Format

For each story, prepare:
- **S**ituation: Context and challenge
- **T**ask: Your specific responsibility
- **A**ction: What YOU did (not the team)
- **R**esult: Quantified outcome
- **+Technical**: The technical depth
- **+Learning**: What you'd do differently

---

## Story 1: The Silent Failure (Debugging)

### "Tell me about a time you debugged a complex production issue."

**Situation:**
> "Two weeks after launching the GCS sink in production, I noticed the GCS buckets had stopped receiving new data. This was in May 2025. We were losing compliance-critical audit data, but there were no error alerts."

**Task:**
> "As the system owner, I needed to identify the root cause and fix it without losing more data. Time was critical because we have SLAs around audit data availability."

**Action:**
> "I took a systematic approach:
>
> First, I checked the obvious - was Kafka Connect running? Yes. Were there messages in the topic? Yes, millions backing up. So the issue was between consumption and GCS write.
>
> I enabled DEBUG logging and found an exception in our SMT filter - NullPointerException when a record had null headers. I added try-catch blocks and deployed. Problem persisted.
>
> Next, I noticed consumer poll timeouts in logs. The default max.poll.interval.ms was 5 minutes, but GCS writes were taking longer for large batches. I tuned this to 5 minutes with smaller batches.
>
> Still had issues. Then I correlated with Kubernetes events and found KEDA autoscaling was triggering frequently. Each scale-up caused consumer group rebalancing, which caused more timeouts, which triggered more scaling. A vicious cycle. I disabled autoscaling temporarily.
>
> Finally, after stabilizing, I found memory exhaustion on large batches. GC logs showed constant full GCs. I increased heap from 512MB to 2GB.
>
> Throughout this, I kept a document of each hypothesis and result. This became our troubleshooting runbook."

**Result:**
> "Fixed the issue with zero data loss - Kafka retained all messages. Backlog cleared in 4 hours once fixed. We implemented better monitoring including consumer lag alerts and JVM memory metrics. The runbook has been used twice since by other teams."

**+Technical:**
> "The core issue was Kafka Connect's error tolerance setting. We had errors.tolerance: all, which silently drops bad records. Combined with lack of header validation in SMT, records were being filtered incorrectly without any alert."

**+Learning:**
> "I should have implemented comprehensive observability before launch - not just 'is it running' but 'is it processing correctly'. Now I advocate for 'silent failure' metrics in all my systems."

---

## Story 2: Multi-Region Rollout (Ambiguity)

### "Tell me about a time you had to make decisions with incomplete information."

**Situation:**
> "We needed to expand from single-region to multi-region for disaster recovery. Requirements were vague - leadership said 'make it resilient' without specifying RTO/RPO targets, budget, or timeline."

**Task:**
> "I needed to define the requirements, design the solution, and execute the rollout - all while continuing to serve production traffic."

**Action:**
> "First, I clarified requirements through stakeholder conversations. I learned compliance needed max 4-hour data gap (RPO) and 1-hour recovery time (RTO). Budget was flexible but we needed to justify costs.
>
> I designed three options:
> 1. Active-Passive: Simple but slow failover
> 2. Active-Active: Complex but immediate failover
> 3. Hybrid: Active-Active for write, single read endpoint
>
> I presented trade-offs to the team. We chose Active-Active because audit is write-heavy and compliance couldn't tolerate delays.
>
> For implementation, I phased it:
> - Week 1: Deploy publisher to second region, write to both Kafkas
> - Week 2: Deploy GCS sink to second region
> - Week 3: Validate data parity between regions
> - Week 4: Update routing to use both regions
>
> The key decision under uncertainty was whether to route by geography or round-robin. With incomplete latency data, I chose geography-based routing using the wm-site-id header. This turned out correct - it reduced cross-region traffic."

**Result:**
> "Achieved Active-Active in 4 weeks. Zero downtime during migration. Passed the disaster recovery test - we failed over a region and recovered in 15 minutes, well under the 1-hour target."

**+Technical:**
> "The tricky part was ensuring exactly-once processing across regions. Kafka's idempotent producer helped, but we also added deduplication in the sink based on request_id."

**+Learning:**
> "I learned to explicitly document assumptions when requirements are vague. My assumptions doc became the de facto requirements spec after stakeholders reviewed it."

---

## Story 3: Building the Reusable Library (Collaboration)

### "Tell me about a time you built something that other teams adopted."

**Situation:**
> "Our team needed audit logging, but I noticed two other teams were building similar functionality. We were duplicating effort and would end up with inconsistent implementations."

**Task:**
> "I proposed building a shared library. My task was to get buy-in from other teams, understand their requirements, and build something that worked for everyone."

**Action:**
> "First, I met with engineers from the other teams. I asked about their requirements - what endpoints they needed to audit, what format they expected, their latency constraints. I documented the union of requirements.
>
> There was initial resistance - 'our needs are different.' I addressed this by showing the common 80% and making the different 20% configurable. For example, one team wanted response body logging, another didn't - I made it a config flag.
>
> I wrote the library with extensibility in mind - the filter checks CCM config for enabled endpoints, so each team configures their own list. I also wrote comprehensive documentation and a migration guide.
>
> To drive adoption, I did a brown-bag session demonstrating the library. Then I personally helped each team integrate - spent an afternoon with each pairing on their PRs."

**Result:**
> "Three teams adopted the library within a month. Integration time went from '2 weeks of custom development' to '1 day with the shared library'. The library is now the standard for all new services in our organization."

**+Technical:**
> "The key design decision was using Spring's OncePerRequestFilter with configurable endpoint matching. This gave flexibility without requiring code changes per service."

**+Learning:**
> "Building for reusability takes longer upfront but pays off exponentially. I now estimate 1.5x time for 'make it reusable' and pitch that to management as an investment."

---

## Story 4: Handling Feedback (Growth)

### "Tell me about a time you received critical feedback."

**Situation:**
> "During code review for the audit library, a senior engineer criticized my thread pool configuration. He said my queue size of 100 was arbitrary and could cause memory issues or silent data loss."

**Task:**
> "I needed to respond to the feedback, validate his concerns, and either fix the issue or justify my design."

**Action:**
> "My first instinct was defensive - I had thought about this! But I paused and really considered his point.
>
> I ran the numbers: each audit payload is ~2KB. Queue of 100 = 200KB. Not a memory issue. But what about at queue rejection? We'd lose audit events silently.
>
> He was right about the silent loss. I added three things:
> 1. A metric for rejected tasks (so we'd see it in dashboards)
> 2. A log at WARN level when queue exceeds 80%
> 3. Documentation explaining the trade-off
>
> I also went back and validated the queue size with actual traffic data. 100 was reasonable for our 100 req/sec, but I added it to the config so teams could adjust.
>
> I thanked him in the PR comments and explained what I changed based on his feedback."

**Result:**
> "The library is more robust. We've had one instance where the queue warning triggered - it helped us identify a downstream slowdown before it became critical. The senior engineer became an advocate for the library."

**+Learning:**
> "I learned to separate ego from code. Critical feedback is a gift - it's someone spending their time to make my work better."

---

# PART 7: CURVEBALL & WHAT-IF SCENARIOS

## "What if" Questions

### Q: "What if you needed 100x the throughput?"

**Answer:**
> "Good question. Let me think through the bottlenecks:
>
> **Current: 2M events/day = ~23 events/sec**
> **Target: 200M events/day = ~2300 events/sec**
>
> **Component 1 (Library):** Should be fine. The async thread pool can be scaled. Real bottleneck is the HTTP call to publisher.
>
> **Change needed:** Switch from HTTP to direct Kafka publishing from the library. Skip the publisher service entirely. This removes a network hop and the publisher becomes a bottleneck point.
>
> **Component 2 (Kafka):** Need to increase partitions. Currently 12 partitions at 23 events/sec = negligible. At 2300/sec, still fine per partition. But for parallelism in consumers, I'd increase to 24-48 partitions.
>
> **Component 3 (Connect):** Need to scale horizontally. Increase tasks.max to match partitions. Add more Connect workers. Each worker can handle ~10K/sec, so we'd need 1-2 more workers.
>
> **GCS:** Can handle virtually unlimited writes. Might need to increase flush frequency to avoid too many small files.
>
> **Cost:** Would go from ~$500/month to ~$5000/month. Mostly Kafka partition storage."

---

### Q: "What if you needed real-time querying instead of batch?"

**Answer:**
> "The current design optimizes for cost and compliance, not real-time. For real-time, I'd redesign:
>
> **Option 1: Stream to Elasticsearch**
> Keep Kafka but add an Elasticsearch sink connector. Gives sub-second query latency. Expensive at scale but great for debugging.
>
> **Option 2: ksqlDB**
> Use Kafka Streams (ksqlDB) to create materialized views. Can answer 'show me last 5 minutes of requests for endpoint X' without external store.
>
> **Option 3: Druid/Clickhouse**
> Real-time OLAP database. Kafka ingests directly. Sub-second aggregation queries at scale.
>
> **My recommendation:** Add Elasticsearch for recent data (last 7 days) alongside GCS for archival. Query router directs recent queries to ES, historical to BigQuery. This balances cost and latency."

---

### Q: "What if a regulator asked for all data for a specific supplier to be deleted?"

**Answer:**
> "GDPR/CCPA right to deletion scenario. This is hard with immutable storage like GCS.
>
> **Current capability:** We can identify records by supplier ID (in headers). But Parquet files are immutable - we can't delete specific rows.
>
> **Deletion process:**
> 1. Identify all Parquet files containing that supplier's data
> 2. For each file, read, filter out supplier's records, write new file
> 3. Delete original file
> 4. Update BigQuery external table if needed
>
> **Challenges:**
> - Expensive (read/write entire files)
> - Slow (could be millions of files over 7 years)
> - Leaves temporary state where data exists in both files
>
> **Better design for future:**
> - Encrypt each record with supplier-specific key
> - For deletion, delete the key
> - Data remains but is unreadable (crypto-shredding)
>
> This is on our roadmap but not yet implemented."

---

### Q: "What happens if Kafka loses messages?"

**Answer:**
> "Kafka is designed to not lose messages, but let's trace the scenarios:
>
> **Scenario 1: Before ack**
> Producer sends message, Kafka broker crashes before ack.
> - Our producer has retries enabled
> - Message gets resent
> - Could cause duplicate (handled by deduplication on request_id)
>
> **Scenario 2: After ack, before replication**
> Message acked, leader crashes before replication.
> - With acks=all, this can't happen (leader waits for replica ack)
> - We use acks=all
>
> **Scenario 3: All replicas fail**
> Three brokers all die simultaneously.
> - Extremely unlikely with proper DC design
> - If it happens, messages in that partition are lost
> - We'd detect via consumer lag jumping, alert would fire
>
> **Scenario 4: Consumer commits before processing**
> Connect commits offset, then crashes before writing to GCS.
> - Kafka Connect uses at-least-once by default
> - It commits only after successful sink write
> - This can cause duplicates (handled by idempotent sink writes based on Kafka offset)
>
> In practice, with proper configuration (acks=all, replication=3, min.insync.replicas=2), Kafka message loss is effectively zero."

---

### Q: "If you had to rebuild this from scratch, what would you change?"

**Answer:**
> "Three significant changes:
>
> **1. Skip the publisher service**
> The publisher (audit-api-logs-srv) adds an HTTP hop. Instead, the library should publish directly to Kafka. Benefits:
> - Lower latency
> - Fewer failure points
> - Simpler architecture
>
> I'd make the library include a Kafka producer. Teams would configure Kafka credentials in their CCM.
>
> **2. Add OpenTelemetry from day one**
> Debugging the silent failure was hard because I didn't have good tracing. I'd add:
> - Trace ID from request through Kafka to GCS
> - Metrics at every hop (filter, producer, consumer, sink)
> - Log correlation across components
>
> **3. Use Apache Iceberg instead of raw Parquet**
> Iceberg gives:
> - ACID transactions
> - Time travel (query historical data)
> - Efficient deletes (for GDPR)
> - Better BigQuery/Athena integration
>
> GCS with Parquet works, but Iceberg would make operations much easier."

---

# PART 8: MOCK INTERVIEW SCRIPTS

## Script 1: Hiring Manager Round (15 minutes)

**Interviewer:** "Tell me about the audit logging system you mentioned."

**You:** "I designed and built an end-to-end audit logging system for Walmart's supplier APIs when Splunk was being decommissioned. It captures 2 million events daily across multiple regions with zero impact on API latency - and the key differentiator is that suppliers can now query their own data. Would you like me to start with the business problem or dive into the technical architecture?"

**Interviewer:** "Start with the business problem."

**You:** "Two things happened at once. First, Walmart was decommissioning Splunk company-wide - it was expensive and the license wasn't being renewed. Second, our external suppliers - Pepsi, Coca-Cola, Unilever - they use our APIs to get inventory data and started asking 'why did my request fail last week?' They wanted visibility into their own interactions.

So I had dual requirements: replace Splunk for our internal debugging needs, AND give suppliers self-service query access to their data. Traditional log aggregators like ELK couldn't solve the second problem - suppliers needed SQL, not grep.

The challenge was our SLA. Suppliers expect sub-200ms responses. Whatever I built couldn't add latency to the critical path."

**Interviewer:** "How did you solve that?"

**You:** "Three components. First, a reusable Spring Boot library that intercepts HTTP requests using a servlet filter. The key insight was using ContentCachingWrapper to capture request and response bodies without blocking, and @Async annotation to offload the logging work to a separate thread pool.

Second, a publisher service that receives these events and publishes to Kafka. We use Avro serialization for schema enforcement and 70% size reduction.

Third, a Kafka Connect sink with custom SMT filters that routes events to GCS buckets by geography. We store in Parquet format for 90% compression and native BigQuery compatibility."

**Interviewer:** "What was the hardest part?"

**You:** "Production debugging. Two weeks after launch, the GCS buckets stopped receiving data. No alerts - the system was failing silently. I spent two days tracing through Kafka Connect logs, finding issues layer by layer - first a null pointer in our SMT filter, then consumer poll timeouts, then KEDA autoscaling causing rebalance storms, finally JVM heap exhaustion.

The experience taught me that distributed systems fail in unexpected combinations. I now build observability as a first-class requirement, not an afterthought."

**Interviewer:** "What would you do differently?"

**You:** "Three things. Add OpenTelemetry tracing from day one - every message should be traceable from request to GCS. Skip the publisher service and publish directly to Kafka from the library - removes a failure point. Use Apache Iceberg instead of raw Parquet for easier operations and GDPR deletes."

---

## Script 2: Googleyness Round (10 minutes)

**Interviewer:** "Tell me about a time you received feedback that was hard to hear."

**You:** "During code review for the audit library, a senior engineer criticized my thread pool configuration. He said the queue size was arbitrary and could cause silent data loss. My first instinct was defensive - I had thought about this!

But I paused and really listened. He was right about silent loss. When the queue fills, tasks are rejected with no indication. For compliance, this is a problem.

I added three things based on his feedback: a metric for rejected tasks, a warning log when queue is 80% full, and documentation of the trade-off. I also thanked him in the PR comments and explained what I changed.

The library is more robust because of that review. We've actually had the queue warning trigger once, which helped us catch a downstream slowdown early."

**Interviewer:** "How do you handle disagreements with your manager?"

**You:** "Early in this project, my manager wanted to ship faster and skip the library approach - just build logging for our service. I disagreed because I saw two other teams building the same thing.

I didn't argue in the meeting. Instead, I gathered data: I showed the cost of three separate implementations, the inconsistency it would create, and estimated the time savings from a shared approach.

I scheduled a 1:1, presented the numbers, and proposed a compromise: invest two extra weeks now to make it reusable, with my personal commitment to drive adoption. If no other team adopted it in a month, I'd own that failure.

He agreed. Three teams adopted within a month. The key was approaching the disagreement with data, not opinions, and taking ownership of the risk."

---

# PART 9: RED FLAGS & HOW TO REDIRECT

## Red Flags That Kill Interviews

| Red Flag | What They Think | How to Redirect |
|----------|-----------------|-----------------|
| **"We did..."** | "Did you do it or did you watch?" | Always say "I" for your work, "we" for team context |
| **Long pauses** | "Doesn't know their own project" | If you need to think, say "Let me trace through that..." |
| **Getting defensive** | "Can't handle feedback" | Say "That's a good point, I hadn't considered..." |
| **Bashing old tech** | "Blames tools, not problems" | Say "X was the right choice then, but I'd choose Y now because..." |
| **No metrics** | "Didn't measure impact" | Always have: events/day, latency, cost, adoption |
| **No failures** | "Either lying or never tried hard things" | Share the debugging story vulnerably |
| **Too much jargon** | "Doesn't actually understand" | Explain like you're teaching a junior |

## Graceful Redirects

**When you don't know something:**
> "I'm not certain about the exact number, but my understanding is [X]. I'd want to verify that with the actual data."

**When you made a mistake:**
> "That's actually something I'd do differently now. At the time I chose X, but in retrospect Y would have been better because..."

**When the question is unclear:**
> "I want to make sure I answer what you're asking. Are you interested in [interpretation A] or [interpretation B]?"

**When you're going too deep:**
> "I'm realizing I'm getting into the weeds. Should I continue at this depth or zoom back out?"

---

# PART 10: CONFIDENCE BUILDERS & POWER PHRASES

## Power Phrases to Use

### Demonstrating Ownership
- "I designed..."
- "I made the decision to..."
- "I took responsibility for..."
- "I drove the adoption of..."

### Showing Technical Depth
- "The key insight was..."
- "The trade-off we accepted was..."
- "The reason we chose X over Y was..."
- "Under the hood, what happens is..."

### Demonstrating Learning
- "What I learned from this was..."
- "If I did this again, I would..."
- "The mistake I made was..., which taught me..."
- "My assumption was X, but we discovered Y..."

### Showing Collaboration
- "I partnered with..."
- "I got buy-in by..."
- "I brought the team along by..."
- "I helped [junior person] understand by..."

## Confidence Anchors

Before the interview, remind yourself:

1. **You solved a real business problem** - Splunk decommissioning + supplier access need
2. **You built something at scale** - 2M events/day is not trivial
3. **You debugged hard problems** - The silent failure story shows resilience
4. **You designed for users** - Suppliers can now self-serve their debugging
5. **You made others successful** - The library was adopted by three teams
6. **You saved significant cost** - 99% cheaper than Splunk at same scale

## The Before vs After Story (POWERFUL!)

| Aspect | Before (Splunk Era) | After (Your Solution) |
|--------|---------------------|----------------------|
| **Supplier debugging** | Call support → 2 day wait → grep logs | Self-service SQL in 30 seconds |
| **Cost per month** | ~$50,000 (Splunk license) | ~$500 (GCS + BigQuery) |
| **Data retention** | 30 days (cost limited) | 7 years (compliance ready) |
| **Add new service** | 2 weeks custom code | 1 day - add Maven dependency |
| **Query capability** | Splunk SPL (complex) | Standard SQL (everyone knows) |
| **API latency impact** | N/A (separate system) | Zero (async design) |

### The Sound Bite:
> *"I turned a Splunk decommissioning crisis into an opportunity to give suppliers something they never had - direct access to their own API data."*

## Body Language Reminders

- Pause before answering (shows thoughtfulness)
- Maintain eye contact during key points
- Use hands to "draw" architecture
- Lean in slightly when they ask follow-ups
- Smile when talking about wins

---

# FINAL CHECKLIST

## Before the Interview

- [ ] Re-read this document
- [ ] Practice the 30-sec, 2-min, 5-min versions out loud
- [ ] Know your numbers: 2M events, <5ms latency, 90% compression, 3 teams adopted
- [ ] Prepare 3 stories: Debugging, Multi-region, Library adoption
- [ ] Have questions ready for them

## During the Interview

- [ ] Listen fully before answering
- [ ] Start answers high-level, offer to go deeper
- [ ] Use "I" for your contributions
- [ ] Include numbers and specifics
- [ ] Acknowledge trade-offs and failures
- [ ] Connect technical work to business impact

## After Each Answer

- [ ] Did I answer what was asked?
- [ ] Did I demonstrate ownership?
- [ ] Did I show technical depth?
- [ ] Did I mention impact?
- [ ] Should I offer to elaborate?

---

*You've got this. You built something real. Now go tell the story.*

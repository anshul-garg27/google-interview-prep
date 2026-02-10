# 🎭 BEHAVIORAL STORIES - STAR FORMAT

> **Purpose:** Practice telling your stories concisely and with impact
> **Format:** Situation → Task → Action → Result (+ Technical + Learning)

---

## HOW TO USE THIS DOCUMENT

1. **Read each story out loud** - Practice speaking, not just reading
2. **Time yourself** - Each story should be 2-3 minutes max
3. **Vary your depth** - Start brief, go deeper when asked
4. **Own your actions** - Use "I" not "we" for your contributions

---

## STORY 1: THE SILENT FAILURE (Debugging)

### The Question Triggers
- "Tell me about a time you debugged a complex production issue"
- "Describe a challenging technical problem you solved"
- "Tell me about a time you worked under pressure"

### THE STORY

**SITUATION** (30 seconds)
> "Two weeks after launching our Kafka-based audit logging system in production, I noticed the GCS buckets had stopped receiving new data. This was in May 2025. We were losing compliance-critical audit data - suppliers need this for debugging their API interactions - but there were no errors in our monitoring dashboards. The system was failing silently."

**TASK** (15 seconds)
> "As the system owner, I needed to identify the root cause and fix it before we lost more data. Time was critical because we have SLAs around audit data availability, and the backlog was growing - messages were piling up in Kafka."

**ACTION** (90 seconds)
> "I took a systematic debugging approach over five days:
>
> First, I checked the obvious - was Kafka Connect running? Yes. Were there messages in the topic? Yes, millions backing up. So the issue was between consumption and GCS write.
>
> I enabled DEBUG logging and found the first issue - NullPointerException in our SMT filter when processing records with null headers. Legacy data didn't have the wm-site-id header we expected. I added try-catch blocks with graceful fallback - if the header is missing, default to the US bucket.
>
> Deployed that fix, but problem persisted. Next I noticed consumer poll timeouts in logs. Kafka Connect's default max.poll.interval.ms was 5 minutes, but our GCS writes were taking longer for large batches. I tuned this and related configs.
>
> Still had issues. Then I correlated with Kubernetes events and discovered KEDA autoscaling was causing instability. Every time lag increased, KEDA scaled up workers. But scaling triggered consumer group rebalancing, which caused more lag, which triggered more scaling - a feedback loop! I disabled KEDA and switched to CPU-based HPA.
>
> Finally, after stabilizing, I traced the remaining issue to JVM heap exhaustion. We were using the default 512MB, but our large batch sizes with Avro deserialization needed more. I increased heap to 2GB.
>
> Throughout this, I documented each hypothesis and result. That document became our troubleshooting runbook."

**RESULT** (30 seconds)
> "Zero data loss - Kafka retained all messages during the entire debugging period, and the backlog cleared in 4 hours once fixed. I implemented comprehensive monitoring including consumer lag alerts, JVM memory metrics, and error rate tracking. The runbook I created has been used twice since by other teams facing similar issues."

**+TECHNICAL** (if asked to go deeper)
> "The core issue was Kafka Connect's error tolerance setting. We had `errors.tolerance: all`, which silently drops bad records without alerting. Combined with lack of header validation in our SMT, records were being filtered incorrectly. The lesson was: default error tolerance can hide problems - always add explicit monitoring for dropped records."

**+LEARNING** (if asked what you'd do differently)
> "I should have implemented comprehensive observability before launch - not just 'is it running' but 'is it processing correctly'. Now I advocate for 'silent failure' metrics in all my systems - counters for dropped messages, filtered records, anything that could be lost silently."

---

## STORY 2: LIBRARY ADOPTION (Collaboration/Influence)

### The Question Triggers
- "Tell me about a time you influenced others without authority"
- "Describe a time you built something other teams adopted"
- "Tell me about working with cross-functional teams"

### THE STORY

**SITUATION** (30 seconds)
> "When Splunk was being decommissioned at Walmart, our team needed audit logging for our supplier APIs. But I noticed two other teams - Inventory Status and Transaction Events - were building similar functionality independently. We were about to have three different, inconsistent implementations."

**TASK** (15 seconds)
> "I proposed building a shared library instead. My task was to get buy-in from the other teams, understand their requirements, and build something that worked for everyone - without having any authority over those teams."

**ACTION** (90 seconds)
> "First, I scheduled meetings with the lead engineers from both teams. I didn't come with a solution - I came with questions. 'What endpoints do you need to audit? What format do you expect? What are your latency constraints?' I documented the union of all requirements.
>
> There was initial resistance - 'our needs are different.' One team wanted response body logging, another didn't. One had strict latency requirements, another was more flexible. I addressed this by showing the common 80% and making the different 20% configurable. Response body logging became a config flag. Endpoint filtering used regex patterns from CCM.
>
> I built the library with extensibility in mind. The filter checks CCM config for enabled endpoints, so each team configures their own list without code changes. I also wrote comprehensive documentation - not just API docs, but a migration guide with step-by-step instructions.
>
> To drive adoption, I did a brown-bag session demonstrating the library. Then I personally helped each team integrate - I spent an afternoon with each team pairing on their PRs. I fixed issues they found and released updates quickly, which built trust."

**RESULT** (30 seconds)
> "Three teams adopted the library within a month. Integration time went from 'two weeks of custom development' to 'one day with the shared library.' The library is now the standard for all new services in our organization. One of the engineers I helped later became an advocate and helped onboard a fourth team."

**+TECHNICAL** (if asked about the design)
> "The key design decision was using Spring's OncePerRequestFilter with @Order(LOWEST_PRECEDENCE) so it runs last in the filter chain. This ensures we capture the final state after all security filters. ContentCachingWrapper handles the HTTP body stream problem, and @Async with a bounded thread pool makes it non-blocking. Zero code changes needed in consuming services - just add the Maven dependency."

**+LEARNING** (what you'd do differently)
> "Building for reusability takes longer upfront but pays off exponentially. I now estimate 1.5x time for 'make it reusable' and pitch that to management as an investment. Also, getting feedback early is crucial - the first version I built wasn't configurable enough. User feedback shaped the final design."

---

## STORY 3: HANDLING CRITICAL FEEDBACK (Growth/Humility)

### The Question Triggers
- "Tell me about a time you received feedback that was hard to hear"
- "Describe a situation where you were wrong"
- "How do you handle disagreement?"

### THE STORY

**SITUATION** (30 seconds)
> "During code review for the audit logging library, a senior engineer named [redacted] criticized my thread pool configuration. He said the queue size of 100 was arbitrary and could cause silent data loss. His comment was public on the PR, and honestly, my first instinct was defensive - I had thought about this!"

**TASK** (15 seconds)
> "I needed to respond constructively, validate whether his concern was legitimate, and either fix the issue or justify my design decision with data."

**ACTION** (60 seconds)
> "I paused before responding. Instead of defending immediately, I really thought about his point.
>
> I ran the numbers: each audit payload is ~2KB. Queue of 100 = 200KB. Not a memory issue. But what about when the queue fills up? The default RejectedExecutionHandler throws an exception, and since we catch all exceptions in the async method, the audit log would be silently dropped.
>
> He was right. We could lose data silently with no indication.
>
> I added three things based on his feedback: First, a Prometheus metric for rejected tasks so we'd see it in dashboards. Second, a log at WARN level when the queue exceeds 80% capacity - early warning. Third, documentation explaining the trade-off and how to monitor for it.
>
> I replied on the PR thanking him for the catch, explained what I was adding, and asked if there was anything else I'd missed."

**RESULT** (30 seconds)
> "The library is more robust because of that review. We've actually had the queue warning trigger once - during a downstream slowdown. We caught it before it became critical and scaled up the thread pool. The senior engineer later became an advocate for the library and helped promote it to other teams."

**+LEARNING**
> "I learned to separate ego from code. Critical feedback is a gift - someone is spending their time to make my work better. Now when I get defensive about feedback, I pause and ask myself: 'What if they're right?' Usually, there's at least a kernel of truth to address."

---

## STORY 4: MULTI-REGION ROLLOUT (Ambiguity/Decision Making)

### The Question Triggers
- "Tell me about a time you had to make decisions with incomplete information"
- "Describe a project with unclear requirements"
- "How do you handle ambiguity?"

### THE STORY

**SITUATION** (30 seconds)
> "We needed to expand our audit logging system from single-region to multi-region for disaster recovery. But requirements were vague - leadership said 'make it resilient' without specifying RTO/RPO targets, budget, or timeline. The only hard requirement was 'don't lose audit data.'"

**TASK** (15 seconds)
> "I needed to define the requirements myself, design a solution, and execute the rollout - all while continuing to serve production traffic with zero downtime."

**ACTION** (90 seconds)
> "First, I clarified requirements through stakeholder conversations. I asked questions like: 'How much data can we afford to lose in a disaster?' 'How quickly do we need to recover?' 'What's the budget?' I learned compliance needed max 4-hour data gap (RPO) and 1-hour recovery time (RTO).
>
> I designed three options and presented trade-offs to the team:
> - **Active-Passive**: Simple, but slow failover (~30 min)
> - **Active-Active**: Complex, but immediate failover
> - **Hybrid**: Active-Active for writes, single read endpoint
>
> We chose Active-Active because audit is write-heavy and compliance couldn't tolerate 30-minute failover delays.
>
> For implementation, I phased it:
> - Week 1: Deploy publisher service to second region, write to both Kafka clusters
> - Week 2: Deploy GCS sink to second region
> - Week 3: Validate data parity between regions
> - Week 4: Update routing to use both regions, test failover
>
> The key decision under uncertainty was routing strategy - geography-based using the wm-site-id header versus round-robin. Without latency data, I chose geography-based because it's more deterministic and reduces cross-region traffic. This turned out correct."

**RESULT** (30 seconds)
> "Achieved Active-Active in 4 weeks. Zero downtime during migration, zero data loss. We passed the disaster recovery test - I intentionally failed over a region and recovered in 15 minutes, well under the 1-hour target. The solution became the reference architecture for other teams doing multi-region deployments."

**+TECHNICAL**
> "The tricky part was ensuring exactly-once processing across regions. Kafka's idempotent producer helped, but we also added deduplication in the sink based on request_id. The SMT filter checks if a record was already processed by that region before writing."

**+LEARNING**
> "I learned to explicitly document assumptions when requirements are vague. I wrote an 'assumptions document' - 'I'm assuming RPO of 4 hours because...' - and shared it with stakeholders. It became the de facto requirements spec after they reviewed and approved it."

---

## STORY 5: SPRING BOOT 3 MIGRATION (Technical Leadership)

### The Question Triggers
- "Tell me about a major technical initiative you led"
- "Describe a time you upgraded a critical system"
- "How do you handle technical debt?"

### THE STORY

**SITUATION** (30 seconds)
> "Spring Boot 2.7 was approaching end-of-life, and Snyk was flagging security vulnerabilities in our dependencies that couldn't be fixed without upgrading. Our main API service, cp-nrti-apis, handles supplier API requests - it's critical infrastructure. We needed to migrate to Spring Boot 3 and Java 17."

**TASK** (15 seconds)
> "I volunteered to lead the migration. My responsibility was to plan the approach, execute the changes, and ensure zero customer-impacting issues in production."

**ACTION** (90 seconds)
> "I analyzed the migration path and identified three main challenges:
>
> First, the javax→jakarta namespace change. Java EE became Jakarta EE, so every import of javax.persistence, javax.validation, javax.servlet had to change. 74 files across entities, controllers, validators, and filters.
>
> Second, RestTemplate deprecation. Spring Boot 3 pushes you toward WebClient, which is reactive. Our codebase was synchronous. I made a strategic decision to use `.block()` for backwards compatibility rather than rewriting everything to be reactive - that would be a separate initiative.
>
> Third, Hibernate 6 compatibility. PostgreSQL enum types needed explicit @JdbcTypeCode annotations. Without them, Hibernate used VARCHAR and PostgreSQL rejected it.
>
> For deployment, I planned a staged approach:
> - Stage environment for 1 week with full regression testing
> - Production with Flagger canary: 10% traffic initially
> - Monitor error rates, latency, memory usage
> - Gradual increase: 25% → 50% → 100% over 24 hours
> - Automatic rollback if error rate exceeded 1%
>
> I updated 158 files total, wrote comprehensive tests for all changes, and documented every breaking change for other teams who would follow."

**RESULT** (30 seconds)
> "Zero customer-impacting issues. Migration completed in 4 weeks. Three minor post-migration fixes - all caught in monitoring and fixed proactively: Snyk vulnerabilities in transitive dependencies, a mutable collection issue, and a logging configuration adjustment. The migration became the template for other team migrations."

**+TECHNICAL** (if asked about specific challenges)
> "The trickiest part was WebClient test mocking. RestTemplate is simple - you mock `exchange()` and you're done. WebClient uses method chaining: `webClient.get().uri().headers().retrieve().bodyToMono()`. Each method in the chain needs to be mocked. Test files basically doubled in complexity."

**+LEARNING**
> "I'd do two things differently. First, start the migration earlier - we waited until close to EOL which added pressure. Second, investigate automated migration tools like OpenRewrite. The javax→jakarta change is deterministic and could be scripted. We did it semi-manually which was time-consuming."

---

## STORY 6: SUPPLIER SELF-SERVICE (User-Centric Thinking)

### The Question Triggers
- "Tell me about a time you improved the user experience"
- "Describe solving a customer problem"
- "How do you think about users?"

### THE STORY

**SITUATION** (30 seconds)
> "Our external suppliers - Pepsi, Coca-Cola, Unilever - use our APIs to get inventory data. When their requests failed, they had no visibility into why. They'd call our support team, open a ticket, wait 1-2 days, and we'd grep through logs to find the issue. This was frustrating for everyone."

**TASK** (15 seconds)
> "When building the audit logging system, I saw an opportunity. Beyond just replacing Splunk for internal logging, I could give suppliers direct access to their own API interaction data."

**ACTION** (60 seconds)
> "I designed the system with supplier self-service in mind:
>
> First, I stored audit data in GCS with Parquet format specifically because BigQuery can query it directly. This wasn't just a cost optimization - it was enabling SQL access.
>
> Second, I ensured the schema captured everything a supplier would need to debug: request body, response body, error messages, timestamps, HTTP status codes.
>
> Third, I worked with our data team to set up BigQuery external tables with row-level security. Each supplier can only see records where the consumer_id matches their identity.
>
> Fourth, I created documentation and sample queries - 'show me all failed requests from last week', 'what's my average response time', 'why did this specific request fail'."

**RESULT** (30 seconds)
> "Suppliers can now self-service debug in 30 seconds instead of waiting 2 days for support. Support ticket volume for 'why did my request fail' dropped significantly. One supplier's engineer told me this was the first time they had visibility into their API interactions with any vendor."

**+LEARNING**
> "The lesson was: when building infrastructure, think about who else might benefit. The primary use case was internal logging, but asking 'who else needs this data?' led to the supplier access feature. It wasn't much extra work but had huge impact."

---

## STORY PRACTICE CHECKLIST

Before your interview, practice each story:

- [ ] **Silent Failure** - Can you tell it in 3 minutes?
- [ ] **Library Adoption** - Do you emphasize YOUR actions?
- [ ] **Critical Feedback** - Do you sound genuinely humble?
- [ ] **Multi-Region** - Can you explain the technical decisions?
- [ ] **Spring Boot 3** - Do you know the specific challenges?
- [ ] **Supplier Self-Service** - Does it show user empathy?

---

## QUICK REFERENCE: Story Selection

| Question About... | Tell This Story |
|------------------|-----------------|
| Debugging | Silent Failure (#1) |
| Collaboration | Library Adoption (#2) |
| Receiving Feedback | Critical Feedback (#3) |
| Ambiguity | Multi-Region (#4) |
| Technical Leadership | Spring Boot 3 (#5) |
| User Focus | Supplier Self-Service (#6) |
| Failure | Silent Failure (early part) |
| Success | Any - but quantify results! |

---

*Practice these out loud until they feel natural!*

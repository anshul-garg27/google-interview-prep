# WALMART → GOOGLE: GOOGLEYNESS & LEADERSHIP QUESTIONS
## Google L4/L5 Interview Preparation

**Critical Context**: This is NOT a generic behavioral guide. Every answer maps to your Walmart Data Ventures work with DEEP technical details. Google Googleyness interviewers will dig 3-4 levels deep. Be ready.

---

## TABLE OF CONTENTS

### GOOGLEYNESS ATTRIBUTES (Google's 6 Core Values)
1. [Thriving in Ambiguity](#1-thriving-in-ambiguity) - 12 questions
2. [Valuing Feedback](#2-valuing-feedback) - 10 questions
3. [Challenging Status Quo](#3-challenging-status-quo) - 11 questions
4. [Putting User First](#4-putting-user-first) - 10 questions
5. [Doing the Right Thing](#5-doing-the-right-thing) - 9 questions
6. [Caring About Team](#6-caring-about-team) - 10 questions

### GOOGLE LEADERSHIP PRINCIPLES
7. [Ownership](#7-ownership-google-style) - 8 questions
8. [Dive Deep](#8-dive-deep-technical-depth) - 6 questions
9. [Bias for Action](#9-bias-for-action) - 5 questions
10. [Deliver Results](#10-deliver-results) - 5 questions

**Total**: 86 questions with STAR answers + follow-up handling

---

# GOOGLEYNESS ATTRIBUTES

## 1. THRIVING IN AMBIGUITY

### Google Definition
"Comfort with uncertainty. Able to make progress when requirements are unclear, specifications change, or the path forward is uncertain."

### Walmart Work Mapping
- DC Inventory Search (no API existed, designed from scratch)
- Kafka Connect custom SMT (no documentation, reverse engineered)
- Multi-region architecture (no precedent in team)

---

### Q1.1: "Tell me about a time you had to solve a problem with incomplete information."

**SITUATION** (Business Context):
"At Walmart Data Ventures, suppliers requested a Distribution Center inventory search API. The challenge: no existing API exposed DC inventory data, and the upstream Enterprise Inventory team had no capacity to build one. I had to design a solution without a clear specification."

**TASK** (Your Role):
"As the tech lead for inventory APIs, I owned delivering this capability. The ambiguity: I didn't know which EI endpoints to call, what data format they returned, or how to map DC numbers to internal node IDs."

**ACTION** (Detailed Technical Decisions):
**Phase 1 - Discovery (Week 1)**:
1. Reverse-engineered EI's internal APIs using Postman
2. Found 3 potential endpoints:
   - `ei-pit-by-item-inventory-read` (Point-in-Time inventory)
   - `ei-onhand-inventory-read` (Store inventory - didn't work for DCs)
   - `ei-dc-inventory-status` (Undocumented, found in network traces)
3. Tested with production traffic captures from Charles Proxy
4. Discovered DC inventory required CID (Customer Item Descriptor), not GTIN

**Phase 2 - Architecture (Week 2)**:
Designed 3-stage pipeline:
```
Stage 1: GTIN → CID conversion (UberKey API)
Stage 2: Validate supplier owns GTIN (PostgreSQL query)
Stage 3: Fetch DC inventory (EI DC API)
```

**Phase 3 - Implementation (Week 3-4)**:
```java
// CompletableFuture for parallel processing
public CompletableFuture<DCInventoryResponse> getDCInventory(
    List<String> gtins, int dcNumber) {

    // Stage 1: Parallel GTIN→CID lookups
    List<CompletableFuture<CIDMapping>> cidFutures = gtins.stream()
        .map(gtin -> CompletableFuture.supplyAsync(
            () -> uberKeyService.getCID(gtin),
            dcInventoryExecutor))
        .collect(Collectors.toList());

    return CompletableFuture.allOf(cidFutures.toArray(new CompletableFuture[0]))
        .thenApply(v -> {
            // Stage 2: Validate supplier authorization
            List<CIDMapping> cids = cidFutures.stream()
                .map(CompletableFuture::join)
                .filter(cid -> validateSupplierGTIN(cid.getGtin()))
                .collect(Collectors.toList());

            // Stage 3: Batch fetch DC inventory
            return eiService.getDCInventoryBatch(cids, dcNumber);
        });
}
```

**Key Design Decisions**:
- CompletableFuture over ParallelStream: Isolated thread pool (20 threads), not shared ForkJoinPool
- Fail-fast validation: Check supplier authorization before expensive EI calls
- Partial success: Return inventory for valid GTINs, errors for invalid ones

**RESULT** (Quantified Impact):
✓ Delivered in 4 weeks (supplier expected 12 weeks)
✓ Zero production incidents after launch
✓ API response time: 1.2s P95 for 100 GTINs (40% faster than similar APIs)
✓ 30,000+ queries/day within 2 months
✓ 3 other teams adopted the pattern for their APIs

**LEARNING** (Growth Mindset):
"I learned to embrace ambiguity by breaking it down: first understand the data (reverse engineering), then design the architecture (3-stage pipeline), finally implement with fallbacks (partial success). When specs are unclear, ship a v1 with conservative assumptions, then iterate based on real usage."

---

### Q1.2: "Describe a situation where requirements changed mid-project."

**SITUATION**:
"I was implementing a Kafka-based audit logging system. Initial requirement: audit 2 endpoints. Mid-project (week 3 of 6), product team said: 'Actually, we need to audit ALL 47 endpoints across 6 services, including request/response bodies for compliance.'"

**TASK**:
"Deadline unchanged. Architecture needed to scale from 2 endpoints to 47 WITHOUT rewriting code in 6 different services."

**ACTION**:
**Original Architecture** (Endpoint-Specific):
```java
// Each service had custom audit logic
@PostMapping("/inventoryActions")
public ResponseEntity<?> inventoryActions(@RequestBody Request req) {
    auditLogService.logInventoryAction(req); // Hardcoded
    return processInventoryAction(req);
}
```

**New Architecture** (Generic Filter):
**1. Created dv-api-common-libraries JAR**:
```java
@Component
@Order(Ordered.HIGHEST_PRECEDENCE)
public class LoggingFilter implements Filter {

    @Autowired
    private AuditLoggingConfig config; // CCM-driven

    @Override
    public void doFilter(ServletRequest request, ServletResponse response,
                         FilterChain chain) {

        // Check if endpoint should be audited
        String requestURI = ((HttpServletRequest) request).getRequestURI();
        if (!config.getEnabledEndpoints().contains(requestURI)) {
            chain.doFilter(request, response);
            return;
        }

        // Wrap to cache bodies
        ContentCachingRequestWrapper wrappedRequest =
            new ContentCachingRequestWrapper((HttpServletRequest) request);
        ContentCachingResponseWrapper wrappedResponse =
            new ContentCachingResponseWrapper((HttpServletResponse) response);

        long startTime = System.currentTimeMillis();
        chain.doFilter(wrappedRequest, wrappedResponse);
        long duration = System.currentTimeMillis() - startTime;

        // Async audit (doesn't block response)
        auditLogService.sendAuditLog(
            buildPayload(wrappedRequest, wrappedResponse, duration)
        );

        wrappedResponse.copyBodyToResponse();
    }
}
```

**2. CCM Configuration for Flexibility**:
```yaml
# NON-PROD-1.0-ccm.yml
auditLoggingConfig:
  enabledEndpoints:
    - /store/inventoryActions
    - /store/directshipment
    - /v1/inventory/events
    - /v1/inventory/search-items
    # ... 43 more endpoints
  isResponseLoggingEnabled: true
  maxRequestSizeBytes: 10240
  maxResponseSizeBytes: 10240
```

**3. Zero-Code Integration for Services**:
```xml
<!-- Just add dependency -->
<dependency>
    <groupId>com.walmart</groupId>
    <artifactId>dv-api-common-libraries</artifactId>
    <version>0.0.54</version>
</dependency>
```

**Trade-offs Evaluated**:

| Approach | Pros | Cons | Decision |
|----------|------|------|----------|
| **AOP (Aspect)** | Annotation-based, fine-grained | Each service needs code changes | ❌ Rejected |
| **Servlet Filter** | Zero code changes, automatic | All requests intercepted (overhead) | ✅ **Chosen** |
| **Spring Interceptor** | Spring-native | Doesn't capture servlet errors | ❌ Rejected |
| **Manual Instrumentation** | Full control | 47 endpoints × 6 services = 282 changes | ❌ Rejected |

**Why Filter Won**:
- Zero code changes in consuming services
- CCM-driven endpoint configuration (hot reload)
- Captures full request/response cycle including errors
- Filter ordering ensures it runs AFTER security filters

**RESULT**:
✓ **Scope**: 2 endpoints → 47 endpoints (2350% increase)
✓ **Timeline**: Delivered on time (week 6)
✓ **Adoption**: 12 services adopted in 2 months (vs. original 6)
✓ **Performance**: 0ms latency impact (async executor pattern)
✓ **Code changes**: 0 lines in consuming services

**LEARNING**:
"When requirements change mid-flight, I ask: 'What's the root need?' Original need wasn't '2 endpoints', it was 'compliance audit trail'. Solution: build for extensibility (filter pattern), not specific endpoints. This pattern now audits 2M+ events/day across 12+ services."

---

### Q1.3: "How do you approach a problem where the solution isn't obvious?"

**SITUATION**:
"During Spring Boot 3 migration, cp-nrti-apis had 200+ failing tests. Root cause unclear - could be Spring Security changes, Hibernate 6 breaking changes, or Jakarta EE namespace issues. 48 hours to fix before deployment deadline."

**TASK**:
"As migration lead, I needed a systematic approach to debug 200+ test failures across 18,000 lines of code."

**ACTION**:
**Phase 1 - Triage & Categorize (Hour 1-4)**:
```bash
# Automated failure analysis script
#!/bin/bash
mvn test > test-output.log 2>&1

# Extract error patterns
grep -A 3 "FAILED" test-output.log | \
  sed 's/.*Exception: \(.*\)/\1/' | \
  sort | uniq -c | sort -rn > failure-patterns.txt

# Output:
# 87 NullPointerException at NrtiStoreServiceImpl.java:142
# 45 SecurityException: Failed to authenticate
# 38 HibernateException: Unknown entity mapping
# 30 Various other errors
```

**Failure Categories**:
1. **NPE (87 tests)**: `NrtiStoreServiceImpl` line 142
2. **SecurityException (45 tests)**: Spring Security 6 breaking changes
3. **HibernateException (38 tests)**: Jakarta Persistence API changes
4. **Misc (30 tests)**: Various issues

**Phase 2 - Root Cause Analysis (Hour 5-12)**:
**NPE Investigation**:
```java
// Line 142 (Spring Boot 2.7)
@Autowired
private TransactionMarkingManager txnManager;

public void getInventory() {
    var txn = txnManager.currentTransaction()  // NPE here
        .addChildTransaction("EI_CALL", "GET")
        .start();
}
```

**Root Cause**: Spring Boot 3 changed `@Autowired` field injection behavior. If bean not found, sets to null instead of failing fast.

**Fix**:
```java
// Changed to constructor injection (Spring Boot 3 best practice)
private final TransactionMarkingManager txnManager;

@Autowired
public NrtiStoreServiceImpl(TransactionMarkingManager txnManager) {
    this.txnManager = Objects.requireNonNull(txnManager,
        "TransactionMarkingManager must not be null");
}
```

**SecurityException Investigation**:
```java
// Spring Security 5 (Boot 2.7)
@Configuration
public class SecurityConfig extends WebSecurityConfigurerAdapter {
    @Override
    protected void configure(HttpSecurity http) throws Exception {
        http.authorizeRequests()
            .antMatchers("/actuator/**").permitAll();
    }
}
```

**Problem**: `WebSecurityConfigurerAdapter` deprecated in Spring Security 6.

**Fix**:
```java
// Spring Security 6 (Boot 3)
@Configuration
public class SecurityConfig {
    @Bean
    public SecurityFilterChain filterChain(HttpSecurity http) throws Exception {
        return http
            .authorizeHttpRequests(auth -> auth
                .requestMatchers("/actuator/**").permitAll()
                .anyRequest().authenticated())
            .build();
    }
}
```

**HibernateException Investigation**:
```java
// javax.persistence (Boot 2.7)
import javax.persistence.Entity;
import javax.persistence.Table;

// jakarta.persistence (Boot 3)
import jakarta.persistence.Entity;
import jakarta.persistence.Table;
```

**Solution**: Global find-replace with validation:
```bash
# Find all javax.persistence imports
find src -name "*.java" -exec grep -l "javax.persistence" {} \; > affected-files.txt

# Replace with jakarta.persistence
for file in $(cat affected-files.txt); do
    sed -i 's/javax\.persistence/jakarta.persistence/g' "$file"
done

# Verify no remaining javax imports
grep -r "javax.persistence" src/ && echo "FAILED" || echo "SUCCESS"
```

**Phase 3 - Batch Fix & Validate (Hour 13-24)**:
```bash
# Fix priority 1 (NPE) - 87 tests
git checkout -b fix/npe-constructor-injection
# Apply constructor injection pattern to 12 service classes
mvn test -Dtest="*ServiceImplTest" # 87 → 0 failures ✓

# Fix priority 2 (Security) - 45 tests
git checkout -b fix/spring-security-6
# Migrate SecurityConfig
mvn test -Dtest="*SecurityTest" # 45 → 0 failures ✓

# Fix priority 3 (Hibernate) - 38 tests
git checkout -b fix/jakarta-persistence
# Run sed script
mvn test -Dtest="*RepositoryTest" # 38 → 3 failures (manual fixes)

# Fix priority 4 (Misc) - 30 tests
# Addressed case-by-case
```

**Decision Framework**:
1. **Automate what you can**: sed script for imports (38 tests fixed in 10 minutes)
2. **Pattern-fix similar issues**: Constructor injection pattern (87 tests)
3. **Manual fix edge cases**: 3 repository tests with custom HQL

**RESULT**:
✓ **Timeline**: 200+ failures → 0 failures in 24 hours
✓ **Pattern reuse**: Constructor injection pattern documented for 5 other services
✓ **Automation**: sed script used by 3 other teams
✓ **Deployment**: On-time deployment, zero rollback

**LEARNING**:
"When the solution isn't obvious, I don't immediately start coding. I spend 20% of time on triage (categorize failures), 40% on root cause analysis (deep dive on top 3 categories), 40% on systematic fixes (automate where possible). This gave me 3 distinct fixes that covered 170+ failures, instead of 200 one-off fixes."

---

### Q1.4: "Tell me about a time you had to make a decision without all the data you wanted."

**SITUATION**:
"During multi-region Kafka architecture design, I had to choose between Active-Active vs Active-Passive replication for audit logs. Product said: 'We need disaster recovery, but I don't know our RPO/RTO requirements.' Finance said: 'Minimize costs.' Engineering said: 'No SRE capacity for complex DR.'"

**TASK**:
"Make architecture decision in 1 week to unblock 3 teams waiting for audit log infrastructure."

**ACTION**:
**Data Gathering (Incomplete)**:
```
Available Data:
✓ Current volume: 50 events/sec average, 120 events/sec peak
✓ Message size: 3KB average
✓ Walmart Kafka SLA: 99.99% (4 nines)
✗ Business RPO/RTO (product didn't know)
✗ Cost budget (finance couldn't provide)
✗ SRE runbook capacity (SRE team said "figure it out")
```

**Options Evaluated**:

**Option 1: Active-Passive (MirrorMaker 2)**
```yaml
Architecture:
  Primary: EUS2 cluster (3 brokers)
  Secondary: SCUS cluster (3 brokers)
  Replication: MirrorMaker 2 (async replication)

Trade-offs:
  Pros:
    - Simple failover (change bootstrap servers in producer config)
    - Lower cost (secondary cluster can be smaller)
    - Walmart pattern (other teams use this)
  Cons:
    - Manual failover (1-5 minute RPO)
    - Data loss window (async replication)
    - Secondary cluster idle (wasted capacity)

Estimated Cost: $2,000/month (asymmetric cluster sizing)
Estimated RPO: 1-5 minutes
Estimated RTO: 5-10 minutes (manual failover)
```

**Option 2: Active-Active (Dual Producer)**
```yaml
Architecture:
  Primary: EUS2 cluster (3 brokers)
  Secondary: SCUS cluster (3 brokers)
  Producers: Write to BOTH clusters
  Consumers: Read from PRIMARY, failover to SECONDARY

Trade-offs:
  Pros:
    - Zero RPO (both clusters have all data)
    - Fast failover (automatic client failover)
    - No data loss
  Cons:
    - Higher cost (both clusters must handle full load)
    - Producer complexity (dual writes)
    - Deduplication needed (consumers see duplicates)

Estimated Cost: $3,500/month (symmetric cluster sizing)
Estimated RPO: 0 seconds (synchronous writes)
Estimated RTO: < 30 seconds (automatic client failover)
```

**Option 3: Single Cluster (No DR)**
```yaml
Architecture:
  Single: EUS2 cluster (3 brokers, RF=3)

Trade-offs:
  Pros:
    - Simplest (no replication logic)
    - Lowest cost
    - Kafka's internal replication (RF=3) = good durability
  Cons:
    - Region failure = total outage
    - No DR (violates compliance?)

Estimated Cost: $1,200/month
Estimated RPO: ∞ (region failure = data loss)
Estimated RTO: ∞ (region failure = manual recovery)
```

**Decision Framework (Without Complete Data)**:
```python
# Assign weights based on known constraints
priorities = {
    "zero_data_loss": 9,      # Compliance audit logs (inferred from "audit")
    "automatic_failover": 6,  # SRE has no capacity (given constraint)
    "low_cost": 4,            # Finance concern, but no hard budget
    "simplicity": 7           # Team velocity concern
}

# Score options
scores = {
    "active_passive": (
        priorities["zero_data_loss"] * 0.3 +      # 1-5 min data loss
        priorities["automatic_failover"] * 0.0 +   # Manual failover
        priorities["low_cost"] * 0.8 +            # Medium cost
        priorities["simplicity"] * 0.6            # Medium complexity
    ),  # Score: 7.1

    "active_active": (
        priorities["zero_data_loss"] * 1.0 +      # Zero data loss
        priorities["automatic_failover"] * 1.0 +   # Auto failover
        priorities["low_cost"] * 0.4 +            # Higher cost
        priorities["simplicity"] * 0.3            # More complex
    ),  # Score: 17.7

    "single_cluster": (
        priorities["zero_data_loss"] * 0.0 +      # Region failure = loss
        priorities["automatic_failover"] * 0.0 +   # No DR
        priorities["low_cost"] * 1.0 +            # Cheapest
        priorities["simplicity"] * 1.0            # Simplest
    )   # Score: 11.0
}

# Decision: Active-Active (highest score)
```

**Key Assumptions Documented**:
```markdown
## DECISION RECORD: Multi-Region Kafka Architecture

**Decision**: Active-Active Dual Producer

**Assumptions** (to validate later):
1. RPO requirement < 5 minutes (inferred from "audit logs" = compliance)
2. Cost budget > $3,000/month (typical Walmart infra spend)
3. SRE can't manually failover during incidents (given constraint)

**Validation Gates**:
- Week 1: Confirm RPO with compliance team
- Week 2: Get finance approval for $3,500/month
- Week 3: Load test dual-write pattern (latency impact)

**Rollback Plan**:
- If cost rejected: Downgrade to Active-Passive (config change only)
- If latency > 50ms: Async dual writes with queue
```

**Implementation** (Mitigating Complexity):
```java
// Dual producer with built-in fallback
@Service
public class DualKafkaProducer {

    @Autowired @Qualifier("primaryKafkaTemplate")
    private KafkaTemplate<String, String> primaryTemplate;

    @Autowired @Qualifier("secondaryKafkaTemplate")
    private KafkaTemplate<String, String> secondaryTemplate;

    public void send(String topic, String key, String value) {
        // Fire both, don't wait for both
        CompletableFuture<SendResult<String, String>> primary =
            primaryTemplate.send(topic, key, value);
        CompletableFuture<SendResult<String, String>> secondary =
            secondaryTemplate.send(topic, key, value);

        // Log if either fails (but don't block)
        primary.whenComplete((result, ex) -> {
            if (ex != null) log.error("Primary cluster send failed", ex);
        });
        secondary.whenComplete((result, ex) -> {
            if (ex != null) log.error("Secondary cluster send failed", ex);
        });
    }
}
```

**RESULT**:
✓ **Decision validated**:
  - Week 2: Compliance confirmed RPO must be < 1 minute ✓
  - Week 3: Finance approved $3,500/month ✓
  - Week 4: Load test showed 12ms P95 latency (< 50ms target) ✓

✓ **Production metrics** (6 months later):
  - Zero data loss incidents
  - 3 automatic failovers during EUS2 maintenance (users didn't notice)
  - Actual cost: $3,200/month (under budget)

✓ **Pattern adopted**: 2 other teams (inventory-events, inventory-status) copied architecture

**LEARNING**:
"When I don't have all the data, I make assumptions EXPLICIT and create validation gates. I chose Active-Active based on 'compliance audit logs likely need low RPO' (assumption), then validated in Week 2. If I'd been wrong, config-only rollback to Active-Passive. Google calls this 'bias for action with reversible decisions' - make the call, but design for changeability."

---

### Q1.5: "Describe a time you navigated ambiguous stakeholder requirements."

**SITUATION**:
"Product team requested: 'Build supplier notification system for DSD shipments.' When I asked for specs:
- Product: 'Suppliers need to know when shipments arrive'
- Operations: 'Store associates need notifications'
- Compliance: 'Must audit all notifications'
They couldn't agree on: who gets notified, how, and when."

**TASK**:
"Deliver notification system in 6 weeks despite conflicting stakeholder requirements."

**ACTION**:
**Week 1 - Requirements Workshop**:
Instead of waiting for consensus, I ran a structured workshop:

**Exercise 1: User Story Mapping**:
```
Supplier Journey:
1. Supplier creates shipment in WMS → 🔔 "Shipment planned"
2. Shipment leaves warehouse → 🔔 "Shipment departed"
3. Shipment arrives at store → 🔔 "Shipment arrived"
4. Store receives shipment → 🔔 "Shipment received"

Store Associate Journey:
1. Shipment 2 hours away → 🔔 "Prepare receiving dock"
2. Shipment arrives → 🔔 "Trailer at door"
3. Receiving complete → 🔔 "Close PO"
```

**Key Insight**: Stakeholders wanted DIFFERENT notifications at DIFFERENT stages.

**Exercise 2: Priority Matrix**:
```
| Stakeholder      | Must Have (P0)           | Nice to Have (P1)       |
|------------------|--------------------------|-------------------------|
| **Supplier**     | Arrival confirmation     | Real-time tracking      |
| **Store Assoc**  | 2-hour advance notice    | Receiving instructions  |
| **Compliance**   | Audit trail              | Notification analytics  |
```

**Consensus**: Focus on P0 for v1, P1 for v2.

**Week 2 - Architecture Design**:
Designed **Event-Driven Multi-Channel Architecture**:

```
┌─────────────────┐
│  DSC Event API  │ (Direct Shipment Capture)
└────────┬────────┘
         │ POST /store/directshipment
         ▼
┌─────────────────┐
│  Kafka Producer │
└────────┬────────┘
         │ Publish to: cperf-nrt-prod-dsc
         ▼
    ┌────────────────────┐
    │  Kafka Topic       │
    │  (cperf-nrt-prod-dsc) │
    └────┬───────────┬───┘
         │           │
         │           └──────────────┐
         ▼                          ▼
┌──────────────────┐       ┌──────────────────┐
│ Notification     │       │ Audit Consumer   │
│ Consumer         │       │                  │
└────────┬─────────┘       └──────────────────┘
         │
         ├─────► Push Notification (Sumo API)
         │        → Store Associates
         │
         └─────► Email Notification (SMTP)
                  → Suppliers
```

**Key Design Decisions**:
1. **Kafka Event Bus**: Decouples notification channels from DSC API
2. **Multi-Consumer Pattern**: Each stakeholder gets their own consumer
3. **Event Schema**: Single event, different consumers extract what they need

**Event Payload**:
```json
{
  "message_id": "uuid",
  "event_type": "PLANNED|ARRIVED|RECEIVED",
  "vendor_id": "544528",
  "destinations": [
    {
      "store_nbr": 3067,
      "planned_eta_at": "2026-02-05T14:00:00Z",
      "actual_arrival_time_at": "2026-02-05T14:23:00Z",
      "loads": [
        {
          "asn": "12345875886",
          "pallet_qty": 28,
          "case_qty": 1324
        }
      ]
    }
  ]
}
```

**Week 3-4 - Phased Implementation**:
**Phase 1: Push Notifications (Store Associates)**:
```java
@Service
public class SumoNotificationConsumer {

    @KafkaListener(topics = "cperf-nrt-prod-dsc", groupId = "sumo-notif-group")
    public void consume(DSCEvent event) {

        // Only notify on ARRIVED events (2-hour window)
        if (!"ARRIVED".equals(event.getEventType())) {
            return;
        }

        for (Destination dest : event.getDestinations()) {
            // Calculate if within 2-hour window
            LocalDateTime arrivalTime = dest.getPlannedEtaAt();
            LocalDateTime now = LocalDateTime.now();
            long hoursUntilArrival = Duration.between(now, arrivalTime).toHours();

            if (hoursUntilArrival <= 2 && hoursUntilArrival >= 0) {
                // Send push notification to store associates
                sumoService.sendPushNotification(
                    PushNotification.builder()
                        .storeNumber(dest.getStoreNbr())
                        .role("US_STORE_ASSET_PROT_DSD")  // Asset Protection - DSD role
                        .title("Shipment Arriving Soon")
                        .body(String.format(
                            "Vendor %s shipment with %d pallets arriving in %d hours",
                            event.getVendorId(),
                            dest.getLoads().get(0).getPalletQty(),
                            hoursUntilArrival
                        ))
                        .actionUrl(String.format(
                            "walmart://store/receiving?asn=%s",
                            dest.getLoads().get(0).getAsn()
                        ))
                        .build()
                );
            }
        }
    }
}
```

**Phase 2: Email Notifications (Suppliers)**:
```java
@Service
public class SupplierEmailConsumer {

    @KafkaListener(topics = "cperf-nrt-prod-dsc", groupId = "supplier-email-group")
    public void consume(DSCEvent event) {

        // Notify supplier on RECEIVED events (confirmation)
        if ("RECEIVED".equals(event.getEventType())) {

            // Lookup supplier email from vendor ID
            String supplierEmail = supplierMappingService
                .getSupplierEmail(event.getVendorId());

            emailService.send(
                Email.builder()
                    .to(supplierEmail)
                    .subject("Shipment Received Confirmation")
                    .body(String.format(
                        "Your shipment (ASN: %s) was received at store %d on %s",
                        event.getDestinations().get(0).getLoads().get(0).getAsn(),
                        event.getDestinations().get(0).getStoreNbr(),
                        event.getDestinations().get(0).getArrivalTimeAt()
                    ))
                    .build()
            );
        }
    }
}
```

**Week 5-6 - Stakeholder Validation**:
Instead of "big bang" launch, I did **staggered rollouts**:

**Week 5: Pilot with 1 store + 1 supplier**:
- Store 3067 (pilot store)
- Vendor 544528 (Core-Mark International)
- **Result**: 10 shipments, 8 successful notifications, 2 edge cases found

**Edge Case 1**: ETA changed after initial notification
**Fix**: Added event deduplication:
```java
@Service
public class NotificationDeduplicator {

    private final RedisTemplate<String, String> redis;

    public boolean shouldNotify(DSCEvent event, String storeNbr) {
        String key = String.format("notif:%s:%s", event.getMessageId(), storeNbr);

        // Set key with 24-hour expiry
        Boolean isNew = redis.opsForValue().setIfAbsent(key, "sent", 24, TimeUnit.HOURS);

        return Boolean.TRUE.equals(isNew); // Only notify if first time
    }
}
```

**Edge Case 2**: Store associate not on shift during arrival
**Fix**: Added configurable notification window in CCM:
```yaml
sumoConfig:
  notificationWindowHours: 2-8  # Only notify 2-8 hours before arrival
  rolesTargeted:
    - US_STORE_ASSET_PROT_DSD
    - US_STORE_RECEIVING_CLERK
```

**Week 6: Production Launch**:
- 50 stores, 10 vendors
- **Metrics**: 95% notification delivery rate, 3-second P95 latency

**Handling Ongoing Ambiguity**:
Even after launch, requirements kept changing:
- Week 8: "Can suppliers get SMS notifications?" → Added SMS consumer
- Week 10: "Can we notify on cancellations?" → Added CANCELLED event type
- Week 12: "Can we include trailer photos?" → Added photo URL to event payload

**Architectural Resilience**:
Because I used **event-driven architecture**, each change was a NEW consumer, not a change to existing code:
```
Original:
- DSC API → Kafka → 2 consumers (push, email)

After changes:
- DSC API → Kafka → 5 consumers (push, email, SMS, photo-upload, analytics)
```

**RESULT**:
✓ **Stakeholder Satisfaction**:
  - Product: "Exactly what we needed, and we can extend it"
  - Operations: "Reduced shipment wait time by 40%"
  - Compliance: "Full audit trail via Kafka"

✓ **Production Metrics** (6 months):
  - 500,000+ notifications sent
  - 97% delivery rate (push notifications)
  - 92% email open rate (suppliers)
  - Zero complaints about spam (smart deduplication)

✓ **Extensibility**:
  - 5 consumers added after launch (vs. 2 at launch)
  - Zero changes to DSC API
  - 3 other teams (returns, recalls, quality) copied pattern

**LEARNING**:
"When requirements are ambiguous, I don't wait for perfect clarity. I:
1. Run structured workshops to extract hidden priorities (user journey mapping)
2. Design for extensibility (event-driven architecture = easy to add consumers)
3. Validate early with pilots (1 store, 1 supplier = found 2 edge cases)
4. Embrace change as a feature (5 consumers added after launch = proof of good architecture)

Google's Googleyness interviews care about THIS: did you make progress despite uncertainty? My answer: YES - shipped v1 in 6 weeks with 80% clarity, then iterated based on real usage."

---

### Q1.6: "How do you handle situations where you don't have precedent to follow?"

**SITUATION**:
"Walmart's Channel Performance team had never built a multi-region, multi-market (US, Canada, Mexico) architecture. Previous services were US-only. Product team asked: 'Can we support international suppliers?' No documentation, no reference architecture, no team expertise."

**TASK**:
"Design multi-market architecture for inventory-status-srv WITHOUT breaking existing US functionality."

**ACTION**:
**Phase 1: Research & Discovery (Week 1)**
Since no internal precedent, I researched externally:

**Netflix Multi-Region Pattern**:
- Edge caching layer per region
- Central database, regional read replicas
- Routing based on geo-IP

**Uber Multi-Market Pattern**:
- Separate databases per market
- Market-specific business logic
- Shared core services

**Stripe Multi-Currency Pattern**:
- Single API, market parameter
- Market-specific configuration
- Centralized billing, localized payments

**Decision**: Hybrid approach - Stripe's "market parameter" + Uber's "market-specific config"

**Phase 2: Architecture Design (Week 2)**
**Challenge**: Existing code was tightly coupled to US:
```java
// Before: US-only hardcoded
@Service
public class InventoryService {

    private static final String EI_URL = "https://ei-inventory-history-lookup.walmart.com";

    public InventoryResponse getInventory(String gtin) {
        return restTemplate.getForObject(
            EI_URL + "/v1/historyForInventoryState/countryCode/us/nodeId/{nodeId}/gtin/{gtin}",
            InventoryResponse.class,
            nodeId, gtin
        );
    }
}
```

**Problem**: Can't just change EI_URL, because Canada/Mexico use DIFFERENT endpoints, DIFFERENT data formats, DIFFERENT authorization.

**Solution: Site-Based Partitioning with Factory Pattern**:

**Step 1: Introduce SiteContext (Thread-Local)**:
```java
@Component
public class SiteContext {

    private static final ThreadLocal<Long> siteIdThreadLocal = new ThreadLocal<>();

    public static void setSiteId(Long siteId) {
        siteIdThreadLocal.set(siteId);
    }

    public static Long getSiteId() {
        Long siteId = siteIdThreadLocal.get();
        if (siteId == null) {
            return 1L; // Default to US
        }
        return siteId;
    }

    public static void clear() {
        siteIdThreadLocal.remove();
    }
}
```

**Step 2: Extract Site-Specific Config to Factory**:
```java
@Component
public class SiteConfigFactory {

    private final Map<String, SiteConfig> configMap;

    @Autowired
    public SiteConfigFactory(
        @ManagedConfiguration("usEiApiConfig") USEiApiCCMConfig usConfig,
        @ManagedConfiguration("caEiApiConfig") CAEiApiCCMConfig caConfig,
        @ManagedConfiguration("mxEiApiConfig") MXEiApiCCMConfig mxConfig
    ) {
        this.configMap = Map.of(
            "1", new SiteConfig(usConfig),  // US
            "3", new SiteConfig(caConfig),  // Canada
            "2", new SiteConfig(mxConfig)   // Mexico
        );
    }

    public SiteConfig getConfig(Long siteId) {
        return configMap.get(String.valueOf(siteId));
    }
}

public class SiteConfig {
    private final String eiApiUrl;
    private final String countryCode;
    private final String authHeader;

    // From CCM config
    public SiteConfig(EiApiCCMConfig ccmConfig) {
        this.eiApiUrl = ccmConfig.getEiApiUrl();
        this.countryCode = ccmConfig.getCountryCode();
        this.authHeader = ccmConfig.getAuthHeader();
    }
}
```

**Step 3: Refactor Service to Use Factory**:
```java
@Service
public class InventoryService {

    private final SiteConfigFactory siteConfigFactory;
    private final SiteContext siteContext;

    public InventoryResponse getInventory(String gtin) {
        Long siteId = siteContext.getSiteId();
        SiteConfig config = siteConfigFactory.getConfig(siteId);

        // Now dynamic based on market
        return restTemplate.getForObject(
            config.getEiApiUrl() + "/v1/historyForInventoryState/countryCode/{country}/nodeId/{nodeId}/gtin/{gtin}",
            InventoryResponse.class,
            config.getCountryCode(),  // "us", "ca", or "mx"
            nodeId,
            gtin
        );
    }
}
```

**Step 4: Populate SiteContext from Request Header**:
```java
@Component
@Order(1)
public class SiteContextFilter implements Filter {

    @Override
    public void doFilter(ServletRequest request, ServletResponse response, FilterChain chain) {
        HttpServletRequest httpRequest = (HttpServletRequest) request;

        // Extract site ID from header
        String siteIdHeader = httpRequest.getHeader("WM-Site-Id");
        Long siteId = (siteIdHeader != null) ? Long.parseLong(siteIdHeader) : 1L;

        SiteContext.setSiteId(siteId);

        try {
            chain.doFilter(request, response);
        } finally {
            SiteContext.clear(); // Prevent thread pool pollution
        }
    }
}
```

**CCM Configuration (Per Market)**:
```yaml
# NON-PROD-1.0-ccm.yml

usEiApiConfig:
  eiApiUrl: "https://ei-inventory-history-lookup.walmart.com"
  countryCode: "us"
  authHeader: "Bearer ${US_EI_TOKEN}"

caEiApiConfig:
  eiApiUrl: "https://ei-inventory-history-lookup-ca.walmart.com"
  countryCode: "ca"
  authHeader: "Bearer ${CA_EI_TOKEN}"

mxEiApiConfig:
  eiApiUrl: "https://ei-inventory-history-lookup-mx.walmart.com"
  countryCode: "mx"
  authHeader: "Bearer ${MX_EI_TOKEN}"
```

**Phase 3: Database Multi-Tenancy (Week 3)**
**Challenge**: Single PostgreSQL database, but data must be isolated per market (compliance).

**Solution: Partition Keys + Composite Primary Keys**:
```java
@Entity
@Table(name = "supplier_gtin_items")
public class NrtiMultiSiteGtinStoreMapping {

    @EmbeddedId
    private NrtiMultiSiteGtinStoreMappingKey primaryKey;

    @PartitionKey  // Hibernate hint for sharding
    @Column(name = "site_id")
    private String siteId;

    @Column(name = "gtin")
    private String gtin;

    @Column(name = "global_duns")
    private String globalDuns;

    @Column(name = "store_nbr", columnDefinition = "integer[]")
    private Integer[] storeNumber;
}

// Composite key includes site_id
@Embeddable
public class NrtiMultiSiteGtinStoreMappingKey implements Serializable {
    private String siteId;
    private String gtin;
    private String globalDuns;
}
```

**Repository with Site Filtering**:
```java
@Repository
public interface NrtiMultiSiteGtinStoreMappingRepository
    extends JpaRepository<NrtiMultiSiteGtinStoreMapping, NrtiMultiSiteGtinStoreMappingKey> {

    // Hibernate automatically adds site_id filter from @PartitionKey
    Optional<NrtiMultiSiteGtinStoreMapping> findByGtinAndGlobalDuns(
        String gtin, String globalDuns
    );
}
```

**Hibernate Interceptor for Automatic Site Filtering**:
```java
@Component
public class MultiTenantInterceptor extends EmptyInterceptor {

    @Override
    public String onPrepareStatement(String sql) {
        Long siteId = SiteContext.getSiteId();
        if (siteId != null && sql.toLowerCase().contains("from supplier_gtin_items")) {
            // Inject site_id filter into WHERE clause
            sql = sql.replace("WHERE", "WHERE site_id = " + siteId + " AND");
        }
        return sql;
    }
}
```

**Result**: Every database query automatically filtered by market, NO code changes in repositories.

**Phase 4: Testing Multi-Market Scenarios (Week 4)**:
```java
@Test
public void testMultiMarketIsolation() {
    // Insert US data
    SiteContext.setSiteId(1L);
    gtinRepo.save(new NrtiMultiSiteGtinStoreMapping("1", "00012345678901", "US_DUNS"));

    // Insert Canada data
    SiteContext.setSiteId(3L);
    gtinRepo.save(new NrtiMultiSiteGtinStoreMapping("3", "00012345678901", "CA_DUNS"));

    // Query from US context
    SiteContext.setSiteId(1L);
    List<NrtiMultiSiteGtinStoreMapping> usResults = gtinRepo.findByGtin("00012345678901");
    assertEquals(1, usResults.size());
    assertEquals("US_DUNS", usResults.get(0).getGlobalDuns());

    // Query from Canada context
    SiteContext.setSiteId(3L);
    List<NrtiMultiSiteGtinStoreMapping> caResults = gtinRepo.findByGtin("00012345678901");
    assertEquals(1, caResults.size());
    assertEquals("CA_DUNS", caResults.get(0).getGlobalDuns());

    // CRITICAL: US context should NOT see Canada data
    SiteContext.setSiteId(1L);
    List<NrtiMultiSiteGtinStoreMapping> allResults = gtinRepo.findAll();
    assertEquals(1, allResults.size());  // Only US data visible
}
```

**Phase 5: Deployment Strategy (Week 5-6)**:
**Risk**: Multi-market changes could break existing US functionality.

**Solution: Feature Flag Rollout**:
```yaml
# CCM feature flags
featureFlags:
  enableCanadaMarket: false  # Start with US-only
  enableMexicoMarket: false
```

**Code**:
```java
@Service
public class InventoryService {

    @Autowired
    private FeatureFlagCCMConfig featureFlags;

    public InventoryResponse getInventory(String gtin) {
        Long siteId = siteContext.getSiteId();

        // Feature flag check
        if (siteId == 3L && !featureFlags.isEnableCanadaMarket()) {
            throw new MarketNotSupportedException("Canada market not yet enabled");
        }

        // Proceed with multi-market logic
        ...
    }
}
```

**Rollout Timeline**:
- **Week 6**: Deploy to prod, feature flags OFF (US-only behavior preserved)
- **Week 7**: Enable Canada market for 1 pilot supplier
- **Week 8**: Enable Canada market for all suppliers
- **Week 9**: Enable Mexico market for 1 pilot supplier
- **Week 10**: Full production rollout

**RESULT**:
✓ **Zero breaking changes**: US functionality unchanged during rollout
✓ **Market expansion**: US (6M queries/month) → Canada (1.2M) → Mexico (800K)
✓ **Code reuse**: 95% of code shared across markets (only config differs)
✓ **Compliance**: Perfect data isolation, passed audit
✓ **Pattern adoption**: 2 other services (inventory-events, cp-nrti-apis) copied multi-market pattern

**Metrics (6 months post-launch)**:
- US: 6,000,000 queries/month
- Canada: 1,200,000 queries/month
- Mexico: 800,000 queries/month
- Cross-market data leak incidents: 0
- Rollback incidents: 0

**LEARNING**:
"When there's no precedent, I:
1. **Research external patterns** (Netflix, Uber, Stripe) for inspiration
2. **Design for extensibility** (SiteContext + Factory pattern = easy to add markets)
3. **Test data isolation rigorously** (multi-market unit tests caught 3 bugs)
4. **De-risk with feature flags** (enabled markets one-by-one, not big bang)

Google values 'thriving in ambiguity'. This project had maximum ambiguity (no team experience, no docs, no reference), but I created a pattern that now supports 8M+ queries/month across 3 markets. That's thriving, not just surviving."

---

### Q1.7: "Describe a time you had to pivot your approach mid-execution."

**SITUATION**:
"I was implementing the Spring Boot 3 migration for cp-nrti-apis. Initial plan: 'Big bang' migration - upgrade all dependencies at once, fix tests, deploy. Week 2 of 4, I discovered Spring Security 6 + Hibernate 6 + Jakarta EE namespace changes were creating 200+ test failures. If I continued the 'big bang' approach, I'd miss the deadline."

**TASK**:
"Deliver Spring Boot 3 migration on time (2 weeks remaining) despite 200+ cascading failures."

**ACTION**:
**Original Plan (Weeks 1-4)**:
```
Week 1: Upgrade Spring Boot parent POM
Week 2: Fix compilation errors
Week 3: Fix test failures
Week 4: Deployment + validation
```

**Week 2 Reality Check**:
```bash
$ mvn test
[INFO] Tests run: 487, Failures: 203, Errors: 0, Skipped: 0

# Failure categories:
- 87 NullPointerException (Spring dependency injection changes)
- 45 SecurityException (Spring Security 6 breaking changes)
- 38 HibernateException (Jakarta Persistence namespace)
- 33 Various other failures
```

**Decision Point**: Continue with 'big bang' or pivot?

**Risk Analysis**:
```python
big_bang_approach = {
    "remaining_time": 2 weeks,
    "failure_rate": 203 / 487,  # 42% test failure
    "estimated_fix_time": 203 * 30 minutes,  # 6,090 minutes = 101 hours
    "available_hours": 2 weeks * 40 hours = 80 hours,
    "risk": "HIGH - Will miss deadline"
}

phased_approach = {
    "remaining_time": 2 weeks,
    "phase_1_failures": 87 (NPE),  # Fix constructor injection pattern
    "phase_2_failures": 45 (Security),  # Migrate SecurityConfig
    "phase_3_failures": 38 (Hibernate),  # Sed script for imports
    "estimated_fix_time": (87 * 10) + (45 * 20) + (38 * 5) + (33 * 30),  # 3,180 minutes = 53 hours
    "available_hours": 80 hours,
    "risk": "MEDIUM - Tight but feasible"
}

# Decision: PIVOT to phased approach
```

**New Plan (Pivoted)**:
```
Phase 1 (Days 1-2): Fix NPE pattern (87 tests) → 87% of tests passing
Phase 2 (Days 3-4): Fix Security config (45 tests) → 97% of tests passing
Phase 3 (Day 5): Automate Hibernate fixes (38 tests) → 99% of tests passing
Phase 4 (Days 6-7): Manual fixes for remaining (33 tests) → 100% passing
Phase 5 (Days 8-10): Deployment + validation
```

**Phase 1: NPE Pattern Fix (Days 1-2)**
**Root Cause Analysis**:
```java
// Spring Boot 2.7: @Autowired field injection fails fast if bean not found
@Autowired
private TransactionMarkingManager txnManager;  // NPE if bean missing

// Spring Boot 3: @Autowired field injection sets to NULL if bean not found (breaking change)
```

**Pattern-Based Fix**:
```bash
# Step 1: Find all @Autowired field injections
grep -r "@Autowired" src/main/java | grep "private" > autowired-fields.txt

# Step 2: Convert to constructor injection (IDE refactoring)
# Example:
# Before:
@Service
public class NrtiStoreServiceImpl {
    @Autowired private TransactionMarkingManager txnManager;
    @Autowired private SupplierMappingService supplierService;
}

# After:
@Service
public class NrtiStoreServiceImpl {
    private final TransactionMarkingManager txnManager;
    private final SupplierMappingService supplierService;

    @Autowired
    public NrtiStoreServiceImpl(
        TransactionMarkingManager txnManager,
        SupplierMappingService supplierService
    ) {
        this.txnManager = Objects.requireNonNull(txnManager);
        this.supplierService = Objects.requireNonNull(supplierService);
    }
}
```

**Automation with IntelliJ Refactoring**:
```
1. Select all @Service classes
2. Refactor → Convert to Constructor Injection
3. Run tests: 487 → 400 failures (87 fixed!)
```

**Phase 2: Security Config Migration (Days 3-4)**
**Root Cause**: Spring Security 6 deprecated `WebSecurityConfigurerAdapter`

**Before**:
```java
@Configuration
public class SecurityConfig extends WebSecurityConfigurerAdapter {
    @Override
    protected void configure(HttpSecurity http) throws Exception {
        http.authorizeRequests()
            .antMatchers("/actuator/**").permitAll()
            .anyRequest().authenticated();
    }
}
```

**After**:
```java
@Configuration
public class SecurityConfig {
    @Bean
    public SecurityFilterChain filterChain(HttpSecurity http) throws Exception {
        return http
            .authorizeHttpRequests(auth -> auth
                .requestMatchers("/actuator/**").permitAll()
                .anyRequest().authenticated())
            .build();
    }
}
```

**Result**: 400 → 355 failures (45 fixed!)

**Phase 3: Hibernate Namespace Changes (Day 5)**
**Root Cause**: Jakarta EE rebranding - `javax.persistence` → `jakarta.persistence`

**Automated Fix**:
```bash
#!/bin/bash
# find-replace-jakarta.sh

# Find all files with javax.persistence imports
find src -name "*.java" -exec grep -l "javax.persistence" {} \; > affected-files.txt

# Replace javax.persistence with jakarta.persistence
for file in $(cat affected-files.txt); do
    sed -i 's/import javax\.persistence/import jakarta.persistence/g' "$file"
    sed -i 's/import javax\.validation/import jakarta.validation/g' "$file"
done

# Verify no remaining javax imports
if grep -r "javax.persistence" src/; then
    echo "ERROR: Still have javax.persistence imports"
    exit 1
else
    echo "SUCCESS: All imports converted to jakarta"
fi

# Run tests
mvn test

# Output: 355 → 320 failures (35 fixed!)
# 3 failures needed manual fixes (custom HQL with javax.persistence.Query)
```

**Phase 4: Manual Edge Cases (Days 6-7)**
**Remaining 33 failures**: Various edge cases requiring manual fixes
- Custom HQL queries with `javax.persistence.Query`
- Mockito mocks expecting old method signatures
- Integration tests with hardcoded URLs

**Example Manual Fix**:
```java
// Before:
import javax.persistence.EntityManager;
import javax.persistence.Query;

public List<String> customQuery() {
    Query query = entityManager.createQuery("SELECT g.gtin FROM Gtin g");
    return query.getResultList();
}

// After:
import jakarta.persistence.EntityManager;
import jakarta.persistence.Query;

public List<String> customQuery() {
    Query query = entityManager.createQuery("SELECT g.gtin FROM Gtin g");
    return query.getResultList();
}
```

**Result**: 320 → 0 failures (33 fixed!)

**Phase 5: Deployment + Validation (Days 8-10)**
```bash
# Day 8: Deploy to dev environment
mvn clean install
kitt deploy --env dev --service cp-nrti-apis

# Day 9: Deploy to stage environment
mvn clean install -Pstage
kitt deploy --env stage --service cp-nrti-apis
# Run R2C contract tests
concord run r2c-tests --service cp-nrti-apis --env stage
# Result: 80% pass rate (meets threshold)

# Day 10: Deploy to production (canary)
mvn clean install -Pprod
kitt deploy --env prod --service cp-nrti-apis --canary
# Flagger canary analysis:
# - 10% traffic for 2 minutes
# - Monitor 5XX errors (< 1% threshold)
# - Promote to 100% traffic
```

**Key Pivot Decisions**:
1. **Pattern over One-Off**: Fix 87 NPEs with constructor injection PATTERN, not 87 individual fixes
2. **Automate Repetitive**: sed script for 38 namespace changes (10 minutes vs. 3 hours manual)
3. **Prioritize by Impact**: Fix 87 NPEs first (highest count), not alphabetically
4. **Accept Manual for Edge Cases**: 33 remaining failures worth manual fixes (high variance)

**RESULT**:
✓ **Timeline**: Delivered on time (Day 10 of 10 remaining)
✓ **Test Coverage**: 0 test failures (100% passing)
✓ **Production Stability**: Zero rollback incidents
✓ **Deployment**: Canary promoted to 100% (no issues detected)
✓ **Pattern Documentation**: Created runbook for other 5 services

**Metrics**:
- Original estimate: 101 hours (would miss deadline)
- Actual time spent: 48 hours (pivot saved 53 hours)
- Test fixes: 203 failures → 0 failures
- Pattern reuse: 4 other services used constructor injection pattern

**LEARNING**:
"When I hit the wall in Week 2, I didn't push through with the original plan. I:
1. **Stopped and analyzed**: Categorized 203 failures into 4 patterns
2. **Ran the math**: Big bang = 101 hours (impossible), phased = 53 hours (tight but feasible)
3. **Communicated the pivot**: Told team 'changing approach, here's why' (transparency)
4. **Automated aggressively**: sed script for 38 fixes in 10 minutes (vs. hours manually)

Google's Googleyness is about ADAPTING. I didn't stick to a failing plan - I pivoted based on data (203 failures → 4 patterns), and shipped on time with zero production issues."

---

## 2. VALUING FEEDBACK

### Google Definition
"Seeks out and incorporates feedback from others. Actively listens to diverse perspectives. Adjusts approach based on input. Gives constructive feedback to help others grow."

### Walmart Work Mapping
- Code reviews (Spring Boot 3 migration peer reviews)
- Post-mortem analysis (Kafka downtime incident)
- User feedback incorporation (supplier API UX improvements)

---

### Q2.1: "Tell me about a time you received critical feedback and how you responded."

**SITUATION**:
"After launching the DC Inventory Search API, our senior architect (John) reviewed my code and said: 'Your CompletableFuture implementation has a memory leak. Under load, you'll exhaust the thread pool.' I was defensive at first - 'It passed load tests!' But he showed me production metrics proving the issue."

**TASK**:
"Fix the memory leak without downtime, and understand WHY I missed it."

**ACTION**:
**Initial Reaction** (Defensive):
"I was confident my code was fine because:
- Load tests passed (1000 concurrent requests)
- Dev/stage environments stable
- Code review approved by 2 engineers

So when John said 'memory leak', I thought 'impossible'."

**Phase 1: Understanding the Feedback (Hour 1-4)**
John showed me production Grafana dashboard:
```
Metric: jvm_threads_live_threads
Value: Steadily increasing from 50 → 200 → 500 → OOM crash (every 6 hours)

Metric: hikaricp_connections_pending
Value: Spiking during DC inventory queries

Correlation: Thread leak during DC inventory API calls
```

**My Code (Problematic)**:
```java
// DC Inventory Service (v1 - LEAKED THREADS)
@Service
public class DCInventoryService {

    public CompletableFuture<DCInventoryResponse> getDCInventory(List<String> gtins) {

        // Phase 1: Convert GTINs to CIDs (parallel)
        List<CompletableFuture<CID>> cidFutures = gtins.stream()
            .map(gtin -> CompletableFuture.supplyAsync(() -> {
                return uberKeyService.getCID(gtin);  // External API call (2-5s)
            }))  // ❌ PROBLEM: Uses ForkJoinPool.commonPool() (shared thread pool)
            .collect(Collectors.toList());

        // Wait for all CID lookups
        CompletableFuture.allOf(cidFutures.toArray(new CompletableFuture[0])).join();

        // Phase 2: Fetch DC inventory (parallel)
        List<CID> cids = cidFutures.stream()
            .map(CompletableFuture::join)
            .collect(Collectors.toList());

        return eiService.getDCInventory(cids);
    }
}
```

**John's Explanation**:
"Your `CompletableFuture.supplyAsync()` uses `ForkJoinPool.commonPool()` by default. This pool is SHARED across the entire JVM - Spring uses it for @Async methods, parallel streams, etc.

When you call `getCID(gtin)` (external API, 2-5 seconds), you're BLOCKING a ForkJoinPool thread. If 100 requests hit your DC API simultaneously, you exhaust the common pool (default size = CPU cores = 8 threads). Now OTHER parts of the app can't use ForkJoinPool."

**Proof** (John's Test):
```java
// Reproduce the issue
@Test
public void testThreadPoolExhaustion() {
    // Simulate 100 concurrent DC inventory requests
    List<CompletableFuture<Void>> futures = IntStream.range(0, 100)
        .mapToObj(i -> CompletableFuture.supplyAsync(() -> {
            // Simulate 5-second external API call
            Thread.sleep(5000);
            return null;
        }))
        .collect(Collectors.toList());

    // Check ForkJoinPool size
    ForkJoinPool commonPool = ForkJoinPool.commonPool();
    System.out.println("Pool size: " + commonPool.getPoolSize());  // 8 (maxed out)
    System.out.println("Active threads: " + commonPool.getActiveThreadCount());  // 8
    System.out.println("Queued submissions: " + commonPool.getQueuedSubmissionCount());  // 92 (waiting!)

    // Now try to use common pool elsewhere
    List<Integer> numbers = IntStream.range(0, 1000).boxed().collect(Collectors.toList());
    long start = System.currentTimeMillis();
    numbers.parallelStream().forEach(n -> {
        // This is BLOCKED because common pool exhausted
    });
    long duration = System.currentTimeMillis() - start;
    System.out.println("Duration: " + duration + "ms");  // 5000ms (should be ~100ms)
}
```

**My Mistake**: I didn't understand that `ForkJoinPool.commonPool()` is SHARED. I thought each CompletableFuture got its own thread.

**Phase 2: Implementing the Fix (Day 2)**
**Solution: Dedicated Thread Pool**:
```java
// Step 1: Create dedicated thread pool
@Configuration
public class AsyncConfig {

    @Bean("dcInventoryExecutor")
    public Executor dcInventoryExecutor() {
        ThreadPoolTaskExecutor executor = new ThreadPoolTaskExecutor();
        executor.setCorePoolSize(20);      // 20 threads (vs. 8 in common pool)
        executor.setMaxPoolSize(50);       // Burst capacity
        executor.setQueueCapacity(100);    // Queue up to 100 requests
        executor.setThreadNamePrefix("dc-inventory-");  // Easy to identify in logs
        executor.setRejectedExecutionHandler(new ThreadPoolExecutor.CallerRunsPolicy());
        executor.initialize();
        return executor;
    }
}

// Step 2: Use dedicated thread pool in CompletableFuture
@Service
public class DCInventoryService {

    @Autowired
    @Qualifier("dcInventoryExecutor")
    private Executor dcInventoryExecutor;

    public CompletableFuture<DCInventoryResponse> getDCInventory(List<String> gtins) {

        // Phase 1: Convert GTINs to CIDs (parallel with DEDICATED pool)
        List<CompletableFuture<CID>> cidFutures = gtins.stream()
            .map(gtin -> CompletableFuture.supplyAsync(() -> {
                return uberKeyService.getCID(gtin);
            }, dcInventoryExecutor))  // ✅ FIX: Use dedicated thread pool
            .collect(Collectors.toList());

        CompletableFuture.allOf(cidFutures.toArray(new CompletableFuture[0])).join();

        List<CID> cids = cidFutures.stream()
            .map(CompletableFuture::join)
            .collect(Collectors.toList());

        return eiService.getDCInventory(cids);
    }
}
```

**Phase 3: Validating the Fix (Day 3)**
**Load Test (Before Fix)**:
```bash
# JMeter: 100 concurrent users, 60 seconds
$ jmeter -n -t dc-inventory-load-test.jmx

Results:
- Average response time: 5.2 seconds
- P95 response time: 12.4 seconds (❌ SLA is 3s)
- Error rate: 12% (timeout errors)
- JVM threads: 200+ (growing continuously)
```

**Load Test (After Fix)**:
```bash
$ jmeter -n -t dc-inventory-load-test.jmx

Results:
- Average response time: 1.8 seconds
- P95 response time: 2.7 seconds (✅ Under 3s SLA)
- Error rate: 0.2% (transient network errors only)
- JVM threads: Stable at 70 (20 for DC inventory, 50 for rest of app)
```

**Phase 4: Documenting Learnings (Day 4)**
I created a team wiki page: "CompletableFuture Best Practices"

**Key Learnings**:
1. **Never block ForkJoinPool.commonPool()**:
   - Common pool is SHARED across entire JVM
   - Blocking = exhausts pool = app-wide performance degradation

2. **When to use dedicated thread pool**:
   - External API calls (> 100ms latency)
   - Database queries (if not using reactive)
   - Any I/O-bound operation

3. **When ForkJoinPool.commonPool() is OK**:
   - CPU-bound operations (computations, parsing)
   - Short-lived tasks (< 10ms)
   - Operations already non-blocking (CompletableFuture chains)

4. **Monitoring**:
   - Added Grafana dashboard: `jvm_threads_live_threads` per executor
   - Alert if thread count > 80% of max pool size
   - Thread dump analysis for leaks

**Phase 5: Proactive Code Review (Week 2)**
I reviewed OTHER services for the same pattern:
```bash
# Find all CompletableFuture.supplyAsync() calls without explicit executor
grep -r "supplyAsync(" src/ | grep -v "executor" > potential-issues.txt

# Found 3 other services with same issue:
- inventory-events-srv (GTIN lookup)
- inventory-status-srv (store inbound queries)
- audit-api-logs-srv (Kafka publish)
```

Created PRs for all 3 services with same fix.

**RESULT**:
✓ **Production Stability**:
  - OOM crashes: 6/week → 0/week
  - P95 latency: 12.4s → 2.7s (54% improvement)
  - Error rate: 12% → 0.2%

✓ **Pattern Adoption**:
  - 3 other services fixed (proactive)
  - Team wiki page viewed 200+ times
  - Added to code review checklist: "Check CompletableFuture uses dedicated executor"

✓ **Personal Growth**:
  - Understood concurrency primitives deeply
  - John became my mentor (weekly 1:1s)
  - Presented "CompletableFuture Pitfalls" at team tech talk

**LEARNING**:
"When John gave me critical feedback, I was initially defensive because I thought I'd done everything right (tests passed, code review approved). But I:
1. **Listened to the data**: Production metrics showed clear thread leak
2. **Asked 'why did I miss this?'**: Load tests didn't catch it because I tested in isolation (no other services sharing common pool)
3. **Shared the learning**: Created wiki page, fixed 3 other services, presented to team

Google's Googleyness is about VALUING feedback, not defending your approach. I valued John's feedback, learned deeply, and paid it forward by helping 3 other teams avoid the same mistake."

---

### Q2.2: "Describe a time you asked for feedback and how you incorporated it."

**SITUATION**:
"I designed the multi-region Kafka architecture (Active-Active dual producer pattern). Before implementation, I asked 3 people for feedback: John (senior architect), Sarah (SRE team lead), and Mark (Kafka platform owner)."

**TASK**:
"Get critical feedback on architecture BEFORE implementation to avoid costly rework."

**ACTION**:
**Phase 1: Structured Feedback Request (Week 1)**
I didn't just say "thoughts on my design?" I structured the feedback request:

**Email Template**:
```
Subject: RFC: Multi-Region Kafka Architecture for Audit Logs

Hi [Name],

I'm designing multi-region Kafka architecture for audit logging (2M events/day).
I'd value your feedback on [SPECIFIC AREA] because of your expertise in [REASON].

**Context**:
- Current: Single-region Kafka (EUS2)
- Goal: Disaster recovery (RPO < 1 minute, RTO < 30 seconds)
- Constraints: SRE has no capacity for manual failover

**Proposed Architecture**:
[Diagram attached]
- Active-Active: Dual producer writes to both EUS2 + SCUS clusters
- Consumers read from primary (EUS2), failover to secondary (SCUS)
- Kafka client auto-failover on cluster failure

**Specific Feedback Request**:
- [For John]: Architecture patterns - is Active-Active overkill? Should I consider Active-Passive?
- [For Sarah]: Operational complexity - can SRE support this with zero manual failover?
- [For Mark]: Kafka platform - does this violate any Walmart Kafka best practices?

**Decision Timeline**:
- Week 1: Gather feedback
- Week 2: Finalize architecture
- Week 3-6: Implementation

Thanks,
Anshul
```

**Why This Works**:
1. **Specific ask**: Not "thoughts?" but "is Active-Active overkill?"
2. **Context provided**: They understand WHY I'm asking
3. **Respect their time**: Clear decision timeline (1 week for feedback)

**Phase 2: Incorporating Feedback (Week 1-2)**
**Feedback from John (Senior Architect)**:
> "Active-Active adds complexity for marginal benefit. You're writing to 2 Kafka clusters on EVERY produce call. What if one cluster is slower? Do you wait for both acks, or just one? If just one, you don't have true Active-Active (data might be missing from secondary)."

**My Initial Reaction**: Defensive - "But I need zero data loss!"

**After Thinking**:
"John's right. Let me map out the scenarios":

| Scenario | Active-Active (both acks) | Active-Active (one ack) | Active-Passive (MirrorMaker) |
|----------|---------------------------|-------------------------|------------------------------|
| **Normal operation** | 2x write latency (wait for both) | 1x write latency (wait for faster) | 1x write latency |
| **One cluster slow** | Timeout errors (producer waits) | Fine (use faster cluster) | Fine (no secondary write) |
| **One cluster down** | Timeout errors (producer fails) | Fine (use healthy cluster) | Manual failover (1-5 min) |
| **Both clusters down** | Fail | Fail | Fail |

**Decision**: John's feedback changed my approach. Instead of "both acks", I implemented "fire-and-forget to both":

**Revised Architecture**:
```java
@Service
public class DualKafkaProducer {

    @Autowired @Qualifier("primaryKafkaTemplate")
    private KafkaTemplate<String, String> primaryTemplate;

    @Autowired @Qualifier("secondaryKafkaTemplate")
    private KafkaTemplate<String, String> secondaryTemplate;

    public void send(String topic, String key, String value) {
        // Fire to BOTH clusters without waiting for both acks
        CompletableFuture<SendResult<String, String>> primaryFuture =
            primaryTemplate.send(topic, key, value);
        CompletableFuture<SendResult<String, String>> secondaryFuture =
            secondaryTemplate.send(topic, key, value);

        // Log if EITHER fails, but don't block producer
        primaryFuture.whenComplete((result, ex) -> {
            if (ex != null) {
                log.error("Primary Kafka send failed (EUS2)", ex);
                metrics.incrementCounter("kafka.primary.failure");
            }
        });

        secondaryFuture.whenComplete((result, ex) -> {
            if (ex != null) {
                log.error("Secondary Kafka send failed (SCUS)", ex);
                metrics.incrementCounter("kafka.secondary.failure");
            }
        });

        // Return immediately (don't wait for both)
    }
}
```

**Trade-off**: Occasional data loss (if one cluster fails BEFORE async replication completes), but 0ms latency impact on producer.

**Feedback from Sarah (SRE)**:
> "We can't manually failover. If you design this with manual runbook, we'll never execute it during incidents (we're oncall for 50+ services). Failover MUST be automatic."

**My Response**: "Got it. Kafka client auto-failover is already automatic (bootstrap servers list includes both clusters). But what if consumers need to EXPLICITLY switch?"

**Sarah's Recommendation**: "Use Kafka consumer group rebalancing. If primary cluster is unhealthy, consumers will automatically rebalance to secondary cluster."

**Implementation**:
```java
// Consumer config with auto-failover
@Configuration
public class KafkaConsumerConfig {

    @Bean
    public ConsumerFactory<String, String> consumerFactory() {
        Map<String, Object> config = new HashMap<>();

        // Both clusters in bootstrap servers (comma-separated)
        config.put(ConsumerConfig.BOOTSTRAP_SERVERS_CONFIG,
            "kafka-eus2-1:9093,kafka-eus2-2:9093,kafka-eus2-3:9093," +  // Primary (EUS2)
            "kafka-scus-1:9093,kafka-scus-2:9093,kafka-scus-3:9093");   // Secondary (SCUS)

        config.put(ConsumerConfig.GROUP_ID_CONFIG, "audit-log-consumer");

        // Auto-failover settings
        config.put(ConsumerConfig.SESSION_TIMEOUT_MS_CONFIG, 30000);  // 30s timeout
        config.put(ConsumerConfig.HEARTBEAT_INTERVAL_MS_CONFIG, 3000);  // 3s heartbeat
        config.put(ConsumerConfig.MAX_POLL_INTERVAL_MS_CONFIG, 300000);  // 5 min max poll

        return new DefaultKafkaConsumerFactory<>(config);
    }
}
```

**Sarah's Validation**: "This works. If EUS2 cluster dies, consumers detect heartbeat failure in 30 seconds, rebalance to SCUS brokers. No manual intervention."

**Feedback from Mark (Kafka Platform Owner)**:
> "Dual producer violates Walmart's 'single source of truth' principle. You'll have duplicate messages in both clusters. How do you handle deduplication?"

**My Initial Thought**: "Crap, I didn't think about deduplication."

**Mark's Recommendation**: "Use Kafka idempotent producer + message key for deduplication."

**Implementation**:
```java
// Producer config with idempotence
@Configuration
public class KafkaProducerConfig {

    @Bean
    public ProducerFactory<String, String> producerFactory() {
        Map<String, Object> config = new HashMap<>();

        config.put(ProducerConfig.BOOTSTRAP_SERVERS_CONFIG, "...");

        // Idempotent producer (prevents duplicates within single producer session)
        config.put(ProducerConfig.ENABLE_IDEMPOTENCE_CONFIG, true);
        config.put(ProducerConfig.ACKS_CONFIG, "all");
        config.put(ProducerConfig.RETRIES_CONFIG, Integer.MAX_VALUE);
        config.put(ProducerConfig.MAX_IN_FLIGHT_REQUESTS_PER_CONNECTION, 5);

        return new DefaultKafkaProducerFactory<>(config);
    }
}

// Consumer-side deduplication
@Service
public class AuditLogConsumer {

    private final ConcurrentHashMap<String, Long> processedMessageIds = new ConcurrentHashMap<>();

    @KafkaListener(topics = "cperf-audit-logs-prod")
    public void consume(ConsumerRecord<String, String> record) {
        String messageId = record.key();  // Use message key as dedup ID

        // Check if already processed (within last 5 minutes)
        Long processedTimestamp = processedMessageIds.get(messageId);
        if (processedTimestamp != null &&
            System.currentTimeMillis() - processedTimestamp < 300000) {
            log.debug("Duplicate message detected: {}", messageId);
            return;  // Skip duplicate
        }

        // Process message
        processAuditLog(record.value());

        // Mark as processed
        processedMessageIds.put(messageId, System.currentTimeMillis());

        // Cleanup old entries (prevent memory leak)
        if (processedMessageIds.size() > 10000) {
            processedMessageIds.entrySet().removeIf(entry ->
                System.currentTimeMillis() - entry.getValue() > 300000
            );
        }
    }
}
```

**Mark's Validation**: "This handles duplicates. Idempotent producer prevents duplicates from retries, consumer-side cache handles duplicates from dual writes."

**Phase 3: Final Architecture (Post-Feedback)**
After incorporating ALL feedback, my architecture changed significantly:

**Before Feedback** (v1):
```
- Active-Active with synchronous dual writes (wait for both acks)
- Manual failover runbook for SRE
- No deduplication
```

**After Feedback** (v2):
```
- Active-Active with asynchronous dual writes (fire-and-forget)
- Automatic consumer failover (Kafka client auto-rebalance)
- Idempotent producer + consumer-side deduplication cache
```

**Changes Summary**:
| Aspect | Before | After | Reason |
|--------|--------|-------|--------|
| **Write latency** | 2x (wait for both) | 1x (fire-and-forget) | John's feedback |
| **Failover** | Manual runbook | Automatic | Sarah's feedback |
| **Deduplication** | None | Idempotent + cache | Mark's feedback |

**RESULT**:
✓ **Production Metrics** (6 months):
  - Write latency: 12ms P95 (vs. 45ms with synchronous dual writes)
  - Failover time: < 30 seconds automatic (vs. 5 minutes manual)
  - Duplicate message rate: 0.02% (vs. 12% without deduplication)

✓ **Architectural Validation**:
  - John approved final design: "Much simpler, achieves same goals"
  - Sarah's SRE team had ZERO manual failover incidents
  - Mark's Kafka platform team promoted design as reference architecture

✓ **Pattern Adoption**:
  - 2 other teams (inventory-events, inventory-status) copied architecture
  - Walmart Kafka best practices updated with dual producer pattern

**LEARNING**:
"I ASKED for feedback upfront, not after implementation. This saved me:
1. **2 weeks of rework**: John's feedback changed async dual-write approach (vs. synchronous)
2. **Operational burden**: Sarah's feedback ensured automatic failover (vs. manual runbook)
3. **Data quality issues**: Mark's feedback added deduplication (vs. 12% duplicate messages)

Google's Googleyness is about SEEKING feedback proactively. I didn't wait for code review - I asked BEFORE implementation. Result: 3 major design changes that saved weeks of rework and prevented production issues."

---

## 3. CHALLENGING STATUS QUO

### Google Definition
"Questions assumptions and proposes better ways of doing things. Doesn't accept 'that's how we've always done it.' Brings fresh perspective and innovative solutions."

### Walmart Work Mapping
- Kafka Connect instead of direct GCS writes
- CompletableFuture instead of ParallelStream
- Event-driven architecture for notifications

---

### Q3.1: "Tell me about a time you challenged the way things were done."

**SITUATION**:
"When I joined the Channel Performance team, audit logging worked like this: Each service manually wrote audit logs to a PostgreSQL database using JDBC. The database had 50M+ rows, queries were slow (5-10 seconds), and the DB crashed twice/month from write load."

**TASK**:
"The team's solution: 'Scale up the database (add more RAM).' I challenged this: 'Why are we using a database for append-only logs?'"

**ACTION**:
**Phase 1: Questioning the Status Quo**
**Current Architecture** (What everyone accepted):
```
Service 1 → JDBC → PostgreSQL audit_logs table (50M rows)
Service 2 → JDBC → PostgreSQL audit_logs table
...
Service 12 → JDBC → PostgreSQL audit_logs table

Queries:
- "Show me all API calls for supplier X in last 30 days"
- Query time: 5-10 seconds (no indexes on supplier_id, slow full table scans)
- Database: Frequent OOM crashes (high write load)
```

**Team's Proposed Solution**: "Scale up PostgreSQL (32GB → 128GB RAM), add read replicas."

**My Challenge**:
"Wait - audit logs are:
1. **Append-only** (never updated)
2. **Time-series** (queried by date range)
3. **High volume** (2M writes/day)
4. **Rarely queried** (analytics, not operational)

Why are we using a TRANSACTIONAL database (ACID guarantees, write-ahead log, B-tree indexes) for this workload? That's like using a sledgehammer to crack a nut."

**Phase 2: Proposing Alternative (Research)**
I researched 3 alternatives:

**Option 1: Keep PostgreSQL (Status Quo)**
```
Pros:
- Familiar (team knows SQL)
- Existing queries work
- ACID guarantees

Cons:
- Slow queries (5-10s)
- Write bottleneck (2M writes/day)
- High cost (128GB RAM = $5,000/month)
- OOM crashes (2x/month)
```

**Option 2: Move to Elasticsearch**
```
Pros:
- Fast full-text search
- Built for logs
- Scalable (horizontal sharding)

Cons:
- Team has zero Elasticsearch expertise
- Operational complexity (cluster management)
- Cost (5-node cluster = $8,000/month)
- Data modeling (need to learn Elasticsearch query DSL)
```

**Option 3: Kafka → Google Cloud Storage (GCS) → BigQuery**
```
Architecture:
Service → Kafka (stream) → Kafka Connect GCS Sink → GCS (Parquet) → BigQuery (analytics)

Pros:
- Kafka: High throughput (millions/sec), durable (replication)
- GCS: Cheap storage ($0.02/GB/month)
- BigQuery: Fast analytics (columnar, petabyte-scale)
- Decouples writes (Kafka) from queries (BigQuery)

Cons:
- New technology (team has minimal Kafka experience)
- Migration effort (rewrite 12 services)
- Operational complexity (Kafka Connect setup)
```

**Decision Framework**:
```python
# Scoring (out of 10)
scoring = {
    "PostgreSQL": {
        "performance": 2,      # Slow queries (5-10s)
        "scalability": 3,      # Vertical scaling only
        "cost": 4,             # $5,000/month
        "reliability": 5,      # OOM crashes 2x/month
        "team_expertise": 10,  # Team knows SQL well
        "total": 24
    },
    "Elasticsearch": {
        "performance": 9,      # Fast full-text search
        "scalability": 8,      # Horizontal scaling
        "cost": 2,             # $8,000/month
        "reliability": 7,      # Mature, but complex ops
        "team_expertise": 2,   # Zero team experience
        "total": 28
    },
    "Kafka_GCS_BigQuery": {
        "performance": 10,     # BigQuery sub-second queries
        "scalability": 10,     # Infinite (GCS/BigQuery)
        "cost": 9,             # $500/month (GCS cheap)
        "reliability": 9,      # Kafka 99.99% SLA
        "team_expertise": 5,   # Some Kafka experience
        "total": 43
    }
}

# Winner: Kafka + GCS + BigQuery (43 points)
```

**Phase 3: Overcoming Resistance**
**Team's Objections**:

**Objection 1**: "We don't know Kafka well enough."
**My Response**:
"Walmart already runs Kafka at scale (500+ topics, 1TB/day). We don't need to OPERATE Kafka, just USE it. I'll handle Kafka Connect setup, team just needs to publish messages (simple Spring Kafka)."

**Proof of Concept** (Week 1):
```java
// Before: Direct database write (complex)
@Service
public class AuditLogService {
    @Autowired private JdbcTemplate jdbcTemplate;

    public void logAuditEvent(AuditLog log) {
        jdbcTemplate.update(
            "INSERT INTO audit_logs (request_id, service_name, endpoint, timestamp, ...) VALUES (?, ?, ?, ?, ...)",
            log.getRequestId(), log.getServiceName(), log.getEndpoint(), log.getTimestamp(), ...
        );
    }
}

// After: Kafka publish (simpler)
@Service
public class AuditLogService {
    @Autowired private KafkaTemplate<String, String> kafkaTemplate;

    public void logAuditEvent(AuditLog log) {
        kafkaTemplate.send("cperf-audit-logs-prod", log.toJson());
    }
}
```

**Result**: "Actually simpler than JDBC!"

**Objection 2**: "What about queries? How do we query GCS files?"
**My Response**:
"We don't query GCS directly. GCS is storage layer. Queries go to BigQuery, which is SQL-compatible."

**Demo** (Week 2):
```sql
-- BigQuery query (exact same SQL as PostgreSQL)
SELECT
    service_name,
    COUNT(*) as request_count,
    AVG(duration_ms) as avg_duration
FROM `wmt-dsi-dv-cperf-prod.audit_logs.api_requests_*`
WHERE
    DATE(_PARTITIONTIME) BETWEEN '2026-01-01' AND '2026-01-31'
    AND supplier_id = 'ABC123'
GROUP BY service_name;

-- Query time: 1.2 seconds (vs. 8 seconds in PostgreSQL)
-- Data scanned: 2.3 GB (columnar Parquet = only scan relevant columns)
```

**Team Reaction**: "Wait, this is FASTER and uses the SAME SQL?"

**Objection 3**: "What if Kafka goes down? Do we lose audit logs?"
**My Response**:
"Good question. Let's design for failure."

**Solution: Multi-Layer Resilience**:
```
Layer 1 (Producer): If Kafka is down, client request still succeeds (async executor + circuit breaker)
Layer 2 (Kafka): Replication factor 3 (data replicated to 3 brokers)
Layer 3 (Kafka Connect): Auto-restart on failure (Kubernetes liveness probes)
Layer 4 (GCS): 11 nines durability (99.999999999%)
```

**Proof** (Chaos Engineering Test):
```bash
# Kill Kafka broker 1
kubectl delete pod kafka-broker-1 -n kafka

# Producer: Continues writing to broker 2 & 3 (auto-failover)
# Result: Zero audit log loss

# Kill ALL Kafka brokers (extreme scenario)
kubectl scale deployment kafka-broker --replicas=0 -n kafka

# Producer: Circuit breaker opens, skips audit logging (client request succeeds)
# Result: Some audit logs lost, but NO CLIENT IMPACT

# Restore Kafka
kubectl scale deployment kafka-broker --replicas=3 -n kafka

# Result: Audit logging resumes in 30 seconds
```

**Phase 4: Pilot & Validation**
Instead of migrating all 12 services, I ran a pilot:

**Week 3: Pilot with 1 service (cp-nrti-apis)**
- Deployed Kafka producer changes
- Set up Kafka Connect GCS Sink
- Configured BigQuery external table

**Metrics (Pilot)**:
| Metric | PostgreSQL (Before) | Kafka+GCS+BigQuery (After) |
|--------|---------------------|----------------------------|
| **Write latency** | 45ms P95 (blocking JDBC) | 2ms P95 (async Kafka) |
| **Query latency** | 8 seconds | 1.2 seconds |
| **Storage cost** | $5,000/month (128GB RAM) | $150/month (GCS) |
| **Reliability** | 2 crashes/month | 0 crashes/6 months |
| **Data retention** | 90 days (disk space limits) | 7 years (GCS cheap) |

**Week 4-8: Rollout to remaining 11 services**
- Created dv-api-common-libraries JAR (automatic Kafka integration)
- Services just added Maven dependency (zero code changes)
- Completed rollout in 4 weeks

**RESULT**:
✓ **Performance**:
  - Write latency: 45ms → 2ms (95% improvement)
  - Query latency: 8s → 1.2s (85% improvement)
  - Database crashes: 2/month → 0 (100% improvement)

✓ **Cost**:
  - $5,000/month (PostgreSQL) → $500/month (GCS + BigQuery)
  - **Savings**: $4,500/month = $54,000/year

✓ **Scalability**:
  - PostgreSQL: 2M writes/day (at limit)
  - Kafka: Tested to 50M writes/day (25x headroom)

✓ **Data Retention**:
  - PostgreSQL: 90 days (disk constraints)
  - GCS: 7 years (compliance requirement met)

✓ **Team Adoption**:
  - 12 services migrated in 8 weeks
  - 3 other teams (not in Channel Performance) adopted the pattern
  - Walmart Data Ventures promoted as reference architecture

**LEARNING**:
"When I challenged 'scale up the database', the team's initial reaction was 'we've always used PostgreSQL for logs.' But I:
1. **Questioned the assumption**: 'Why transactional DB for append-only logs?'
2. **Proposed data-driven alternatives**: Scored 3 options, Kafka+GCS+BigQuery won (43 vs. 24 points)
3. **Addressed objections with proof**: POC showed simpler code, chaos test proved resilience
4. **De-risked with pilot**: 1 service first, then 11 more

Google's Googleyness is about CHALLENGING status quo with BETTER alternatives. I didn't just complain about PostgreSQL - I built a solution that was 95% faster, 90% cheaper, and infinitely more scalable."

---

(Continuing with remaining Googleyness questions and Leadership sections...)

---

# QUICK REFERENCE: METRICS CHEATSHEET

## Kafka Audit Logging System
- **Volume**: 2M+ events/day, 50 events/sec avg, 120 events/sec peak
- **Latency**: 0ms client impact (async), <10ms P95 Kafka publish
- **Cost**: $5,000/mo → $500/mo (90% reduction)
- **Adoption**: 12+ services, 8 teams

## DC Inventory Search API
- **Performance**: 1.2s P95 for 100 GTINs (40% faster than similar APIs)
- **Volume**: 30,000+ queries/day
- **Delivery**: 4 weeks (vs. 12 weeks estimated)
- **Pattern Reuse**: 3 teams adopted 3-stage pipeline

## Spring Boot 3 Migration
- **Test Failures**: 203 → 0 in 24 hours
- **Services Migrated**: 6 services
- **Timeline**: On-time delivery (48 hours actual vs. 101 hours big-bang)
- **Zero Production Incidents**: 0 rollbacks

## Multi-Region Architecture
- **Failover Time**: <30 seconds (automatic)
- **RPO**: 0 seconds (zero data loss)
- **Cost**: $3,200/month (under $3,500 budget)
- **Adoption**: 2 other teams copied architecture

## DSD Notification System
- **Notifications**: 500,000+ sent (6 months)
- **Delivery Rate**: 97% (push), 92% email open rate
- **Extensibility**: 5 consumers added post-launch (vs. 2 at launch)

## Common Library (dv-api-common-libraries)
- **Adoption**: 12+ services
- **Code Changes**: 0 lines in consuming services
- **Performance**: 0ms latency impact
- **Releases**: 57+ versions

---

**END OF PART 1 - GOOGLEYNESS ATTRIBUTES**

*Note: This document continues with Q3.2-Q10.5 covering remaining Googleyness attributes and Leadership Principles. Total 86 questions.*

*Word Count: ~15,000 words (target: 30,000+ for complete document)*

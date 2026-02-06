# DEEP DIVE: Bullet 4 - Spring Boot 3 & Java 17 Migration

> **🎯 RESUME BULLET 4:** "Led Spring Boot 3 & Java 17 migration for cp-nrti-apis, handling javax→jakarta namespace changes, Hibernate 6 compatibility, and zero-downtime production deployment"

---

## TABLE OF CONTENTS

1. [Why This Migration Matters](#1-why-this-migration-matters)
2. [The Migration Scope](#2-the-migration-scope)
3. [Key Technical Challenges](#3-key-technical-challenges)
4. [Code Changes Deep Dive](#4-code-changes-deep-dive)
5. [PR History & Timeline](#5-pr-history--timeline)
6. [Production Deployment Strategy](#6-production-deployment-strategy)
7. [Interview Q&A Bank](#7-interview-qa-bank)
8. [Technical Decision Tree](#8-technical-decision-tree)
9. [Debugging Stories](#9-debugging-stories)
10. [Numbers to Remember](#10-numbers-to-remember)

---

## 1. WHY THIS MIGRATION MATTERS

### The Business Context

```
┌────────────────────────────────────────────────────────────────────────────────┐
│                        WHY SPRING BOOT 3 MIGRATION?                             │
│                                                                                 │
│  1. END OF SUPPORT                                                              │
│     • Spring Boot 2.7 EOL: November 2023                                       │
│     • Java 11 extended support ending                                          │
│     • No security patches for older versions                                   │
│                                                                                 │
│  2. PERFORMANCE IMPROVEMENTS                                                    │
│     • GraalVM native image support                                             │
│     • Faster startup times                                                     │
│     • Better memory efficiency with Java 17                                    │
│                                                                                 │
│  3. SECURITY REQUIREMENTS                                                       │
│     • Snyk/Sonar flagging vulnerabilities in old versions                      │
│     • Compliance requirements for latest stable versions                       │
│                                                                                 │
│  4. ECOSYSTEM COMPATIBILITY                                                     │
│     • New GTP BOM versions require Spring Boot 3                               │
│     • Walmart internal libraries moving to Jakarta EE                          │
└────────────────────────────────────────────────────────────────────────────────┘
```

### The Interview Sound Bite

> "Spring Boot 2.7 was approaching end-of-life, and we were getting security vulnerabilities flagged by Snyk. I led the migration to Spring Boot 3.2 and Java 17, which involved updating 74 files for the javax→jakarta namespace change, migrating from RestTemplate to WebClient, and adapting to Hibernate 6's new query semantics. Deployed to production with zero downtime using Flagger canary deployments."

---

## 2. THE MIGRATION SCOPE

### Version Changes

| Component | Before | After |
|-----------|--------|-------|
| **Spring Boot** | 2.7.x | 3.5.7 |
| **Java** | 11 | 17 |
| **Hibernate** | 5.x | 6.x |
| **Jakarta EE** | javax.* | jakarta.* |
| **GTP BOM** | 1.x | 2.2.4 |
| **Spring Security** | 5.x | 6.x |

### Files Changed

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         MIGRATION IMPACT ANALYSIS                            │
│                                                                              │
│  PR #1312 Statistics:                                                        │
│  ├── Files changed: 136                                                      │
│  ├── Lines added: 1,732                                                      │
│  ├── Lines removed: 1,287                                                    │
│  └── Net change: +445 lines                                                  │
│                                                                              │
│  Breakdown by category:                                                      │
│  ├── Entity classes (javax→jakarta): 15 files                               │
│  ├── Controllers (validation): 12 files                                     │
│  ├── Services (RestTemplate→WebClient): 8 files                             │
│  ├── Tests (updated mocking): 34 files                                      │
│  ├── Kafka (ListenableFuture→CompletableFuture): 4 files                   │
│  ├── Configuration (pom.xml, kitt.yml): 3 files                             │
│  └── Exception handlers: 2 files                                            │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 3. KEY TECHNICAL CHALLENGES

### Challenge 1: javax → jakarta Namespace Migration

**The Problem:**
Java EE was transferred from Oracle to Eclipse Foundation and renamed to Jakarta EE. All `javax.*` packages became `jakarta.*`.

**Affected Areas:**
```java
// BEFORE (Spring Boot 2.x)
import javax.persistence.Entity;
import javax.persistence.Column;
import javax.validation.Valid;
import javax.validation.constraints.NotNull;
import javax.servlet.http.HttpServletRequest;

// AFTER (Spring Boot 3.x)
import jakarta.persistence.Entity;
import jakarta.persistence.Column;
import jakarta.validation.Valid;
import jakarta.validation.constraints.NotNull;
import jakarta.servlet.http.HttpServletRequest;
```

**Files Changed:** 74 files across entities, controllers, validators, filters, and tests

**Interview Point:**
> "The javax→jakarta migration wasn't just a find-replace. Some classes changed behavior, and we had to carefully test each component. For example, servlet filters needed updated method signatures."

---

### Challenge 2: RestTemplate → WebClient Migration

**The Problem:**
`RestTemplate` is deprecated in Spring Boot 3. Must migrate to `WebClient` (reactive, non-blocking).

**Before (RestTemplate):**
```java
@Service
public class NrtiStoreServiceImpl {

    @Autowired
    private RestTemplate restTemplate;

    public StoreResponse getStoreData(String storeId) {
        ResponseEntity<StoreResponse> response = restTemplate.exchange(
            url,
            HttpMethod.GET,
            new HttpEntity<>(headers),
            StoreResponse.class
        );
        return response.getBody();
    }
}
```

**After (WebClient):**
```java
@Service
public class NrtiStoreServiceImpl {

    @Autowired
    private WebClient webClient;

    public StoreResponse getStoreData(String storeId) {
        return webClient.get()
            .uri(url)
            .headers(h -> h.addAll(headers))
            .retrieve()
            .bodyToMono(StoreResponse.class)
            .block();  // Blocking for compatibility with sync code
    }
}
```

**Key Considerations:**
- WebClient is reactive by default; used `.block()` to maintain sync behavior
- Error handling changed from exceptions to status handlers
- Had to mock `WebClient` chain in tests (more complex than RestTemplate mocking)

**Interview Point:**
> "The WebClient migration was the trickiest part. We had to decide between going fully reactive or using `.block()` for backwards compatibility. We chose `.block()` because our entire codebase was synchronous, and a full reactive migration would be a separate initiative."

---

### Challenge 3: Kafka ListenableFuture → CompletableFuture

**The Problem:**
Spring Kafka deprecated `ListenableFuture` in favor of `CompletableFuture`.

**Before:**
```java
ListenableFuture<SendResult<String, Message>> future = kafkaTemplate.send(message);
future.addCallback(
    result -> log.info("Success: {}", result),
    ex -> log.error("Failed: {}", ex.getMessage())
);
```

**After:**
```java
CompletableFuture<SendResult<String, Message<IacKafkaPayload>>> future =
    kafkaTemplate.send(iacKafkaMessage);

future
    .thenAccept(result -> {
        RecordMetadata metadata = result.getRecordMetadata();
        log.info("Sent successfully to partition: {}, offset: {}",
                 metadata.partition(), metadata.offset());
    })
    .exceptionally(ex -> {
        log.warn("Failed in primary region, trying secondary: {}", ex.getMessage());
        handleFailure(topicName, message, messageId).join();
        return null;
    })
    .join();
```

**Interview Point:**
> "The CompletableFuture migration improved our error handling. With the old callback style, exceptions could be swallowed. With CompletableFuture, we use `.exceptionally()` and `.join()` to ensure failures are properly handled and we can fall back to the secondary Kafka region."

---

### Challenge 4: Hibernate 6 Compatibility

**The Problem:**
Hibernate 6 changed several default behaviors and deprecated methods.

**Key Changes:**
```java
// 1. Enum Type Handling - NEW in Hibernate 6
@Column(name = "status", columnDefinition = "status_enum")
@Enumerated(EnumType.STRING)
@JdbcTypeCode(SqlTypes.NAMED_ENUM)  // NEW: Required for PostgreSQL enums
private Status status;

// 2. Query Changes
// BEFORE: Positional parameters started at 0
// AFTER: Positional parameters start at 1 (JPA spec)

// 3. Type mappings changed
// org.hibernate.type.* → org.hibernate.type.SqlTypes.*
```

**Interview Point:**
> "Hibernate 6 introduced stricter JPA compliance. We had to add `@JdbcTypeCode` annotations for PostgreSQL enum types and update some queries that used non-standard Hibernate features."

---

### Challenge 5: Spring Security 6 Changes

**The Problem:**
Spring Security 6 changed the configuration style from `WebSecurityConfigurerAdapter` to component-based configuration.

**Before (Spring Security 5):**
```java
@Configuration
@EnableWebSecurity
public class SecurityConfig extends WebSecurityConfigurerAdapter {

    @Override
    protected void configure(HttpSecurity http) throws Exception {
        http
            .csrf().disable()
            .authorizeRequests()
            .antMatchers("/actuator/**").permitAll()
            .anyRequest().authenticated();
    }
}
```

**After (Spring Security 6):**
```java
@Configuration
@EnableWebSecurity
public class SecurityConfig {

    @Bean
    public SecurityFilterChain filterChain(HttpSecurity http) throws Exception {
        http
            .csrf(csrf -> csrf.disable())
            .authorizeHttpRequests(auth -> auth
                .requestMatchers("/actuator/**").permitAll()
                .anyRequest().authenticated()
            );
        return http.build();
    }
}
```

**Key Changes:**
- No more `WebSecurityConfigurerAdapter` (deprecated)
- Lambda-based DSL configuration
- `antMatchers()` → `requestMatchers()`
- `authorizeRequests()` → `authorizeHttpRequests()`

---

### Challenge 6: Exception Handler Updates

**The Problem:**
Spring MVC 6 introduced `NoResourceFoundException` for 404 errors (previously no specific exception).

**Code Added:**
```java
@ExceptionHandler(NoResourceFoundException.class)
@ResponseStatus(NOT_FOUND)
public ResponseEntity<NoResourceErrorBean> handleNoResourceFoundException(
    NoResourceFoundException ex, HttpServletRequest request) {

    String timestamp = ZonedDateTime.now(ZoneOffset.UTC)
        .format(DateTimeFormatter.ISO_INSTANT);

    NoResourceErrorBean errorBean = NoResourceErrorBean.builder()
        .timestamp(timestamp)
        .status(NOT_FOUND.value())
        .error(NOT_FOUND.getReasonPhrase())
        .path(request.getRequestURI())
        .build();

    return new ResponseEntity<>(errorBean, NOT_FOUND);
}
```

**Also Changed:**
- `HttpStatus` → `HttpStatusCode` in method signatures
- Added custom error response beans for better API consistency

---

## 4. CODE CHANGES DEEP DIVE

### Entity Migration Example

**ParentCompanyMapping.java:**
```java
package com.walmart.cpnrti.entity;

// CHANGED: javax.persistence → jakarta.persistence
import jakarta.persistence.Column;
import jakarta.persistence.EmbeddedId;
import jakarta.persistence.Entity;
import jakarta.persistence.EnumType;
import jakarta.persistence.Enumerated;
import jakarta.persistence.Index;
import jakarta.persistence.Table;
import jakarta.persistence.UniqueConstraint;

// NEW: Hibernate 6 type handling
import org.hibernate.annotations.JdbcTypeCode;
import org.hibernate.type.SqlTypes;

@Entity
@Table(
    name = NRT_CONSUMERS_TABLE,
    uniqueConstraints = @UniqueConstraint(
        name = "uk_nrt_consumers",
        columnNames = {"consumer_id", "site_id"}
    ),
    indexes = @Index(name = "consumer_id", columnList = "consumer_id")
)
public class ParentCompanyMapping extends BaseEntity {

    // NEW: JdbcTypeCode for PostgreSQL enum compatibility
    @Column(name = "status", columnDefinition = "status_enum", nullable = false)
    @Enumerated(EnumType.STRING)
    @JdbcTypeCode(SqlTypes.NAMED_ENUM)  // ← Hibernate 6 requirement
    private Status status;
}
```

### Controller Migration Example

**NrtiStoreControllerV1.java:**
```java
// CHANGED: javax.validation → jakarta.validation
import jakarta.validation.Valid;
import jakarta.validation.constraints.NotBlank;

@RestController
@RequestMapping("/api/v1/store")
public class NrtiStoreControllerV1 {

    @PostMapping("/inventory")
    public ResponseEntity<InventoryResponse> getInventory(
            @Valid @RequestBody InventoryRequest request) {  // jakarta.validation.Valid
        return ResponseEntity.ok(service.getInventory(request));
    }
}
```

### Test Migration Example

**HttpServiceImplTest.java:**
```java
// BEFORE: Mocking RestTemplate
@Mock
private RestTemplate restTemplate;

@Test
void testGetRequest() {
    when(restTemplate.exchange(any(), eq(HttpMethod.GET), any(), eq(String.class)))
        .thenReturn(ResponseEntity.ok("response"));

    String result = httpService.get(url, headers, String.class);
    assertEquals("response", result);
}

// AFTER: Mocking WebClient (more complex)
@Mock
private WebClient webClient;
@Mock
private WebClient.RequestHeadersUriSpec requestHeadersUriSpec;
@Mock
private WebClient.RequestHeadersSpec requestHeadersSpec;
@Mock
private WebClient.ResponseSpec responseSpec;

@Test
void testGetRequest() {
    // Chain mocking required for WebClient
    when(webClient.get()).thenReturn(requestHeadersUriSpec);
    when(requestHeadersUriSpec.uri(anyString())).thenReturn(requestHeadersSpec);
    when(requestHeadersSpec.headers(any())).thenReturn(requestHeadersSpec);
    when(requestHeadersSpec.retrieve()).thenReturn(responseSpec);
    when(responseSpec.bodyToMono(String.class)).thenReturn(Mono.just("response"));

    String result = httpService.get(url, headers, String.class);
    assertEquals("response", result);
}
```

---

## 5. PR HISTORY & TIMELINE

### Core Migration PRs

| PR# | Title | State | Key Changes | Interview Point |
|-----|-------|-------|-------------|-----------------|
| **#1312** | Feature/springboot 3 migration | MERGED | **THE MAIN PR** - Complete migration | "136 files, 1700+ lines changed" |
| **#1337** | Feature/springboot 3 migration (Release) | MERGED | Production deployment | "Zero-downtime deployment" |
| #1332 | Snyk issues fixed + mutable list | MERGED | Security + immutable collections | "Post-migration hardening" |
| #1334 | Snyk issues merged to stage | MERGED | Stage deployment | "Staged rollout" |
| #1335 | Major sonar issue fixed | MERGED | Code quality | "Addressed static analysis" |
| #1176 | Spring boot 3 migration final | MERGED | Earlier attempt | "Iterative approach" |
| #1153 | Logging issue fix | MERGED | Log4j2 compatibility | "Framework compatibility" |

### Timeline

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        MIGRATION TIMELINE                                    │
│                                                                              │
│  Dec 2024  │ #1153: Initial logging fixes for SB3 compatibility            │
│            │                                                                 │
│  Mar 2025  │ #1176: First complete migration attempt                        │
│            │ #1263: Snyk security fixes                                     │
│            │                                                                 │
│  Apr 2025  │ #1312: THE MAIN PR - Complete migration ✅                     │
│            │ #1332-1335: Post-migration fixes                               │
│            │ #1337: Production deployment (Apr 30) 🚀                       │
│            │                                                                 │
│  May 2025  │ Stable in production, monitoring for issues                    │
│            │                                                                 │
│  Aug 2025  │ #1424-1425: GTP BOM updates for continued compatibility       │
│            │                                                                 │
│  Sep 2025  │ #1498-1501: GTP BOM issue fixes                                │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 6. PRODUCTION DEPLOYMENT STRATEGY

### Zero-Downtime Deployment with Flagger

```
┌────────────────────────────────────────────────────────────────────────────────┐
│                      CANARY DEPLOYMENT STRATEGY                                 │
│                                                                                 │
│   Phase 1: Deploy to Stage (1 week)                                            │
│   ├── Run full regression suite                                                │
│   ├── Performance testing (load tests)                                         │
│   └── Validate all endpoints working                                           │
│                                                                                 │
│   Phase 2: Canary to Production (traffic split)                                │
│   ├── 10% traffic → New version (Spring Boot 3)                               │
│   ├── 90% traffic → Old version (Spring Boot 2.7)                             │
│   ├── Monitor error rates, latency, CPU/memory                                │
│   └── Duration: 4 hours                                                        │
│                                                                                 │
│   Phase 3: Gradual Rollout                                                     │
│   ├── 25% → 50% → 75% → 100%                                                  │
│   ├── Automatic rollback if error rate > 1%                                   │
│   └── Duration: 24 hours total                                                 │
│                                                                                 │
│   Phase 4: Full Production                                                     │
│   ├── 100% traffic on Spring Boot 3                                           │
│   ├── Old version kept for 48 hours (rollback option)                         │
│   └── Decommission old version                                                 │
└────────────────────────────────────────────────────────────────────────────────┘
```

### kitt.yml Changes for Deployment

```yaml
# Changes made in PR #1312
metadata:
  name: cp-nrti-apis

spec:
  java:
    version: "17"  # CHANGED from 11

  container:
    image: docker.io/library/openjdk:17-slim  # CHANGED

  deployment:
    replicas: 3
    strategy:
      type: RollingUpdate
      rollingUpdate:
        maxUnavailable: 0
        maxSurge: 1

  flagger:
    enabled: true
    analysis:
      threshold: 5
      maxWeight: 50
      stepWeight: 10
```

---

## 7. INTERVIEW Q&A BANK

### Q1: "Tell me about the Spring Boot 3 migration you led."

**Level 1 Answer:**
> "I led the migration of cp-nrti-apis from Spring Boot 2.7 to Spring Boot 3.2, along with Java 11 to Java 17. The main challenges were the javax→jakarta namespace change, migrating from RestTemplate to WebClient, and adapting to Hibernate 6. We deployed to production with zero downtime using canary deployments."

**Level 2 Follow-up:** "What was the biggest challenge?"

> "The WebClient migration was trickiest. RestTemplate has been the standard for years, and the team was comfortable with it. WebClient is reactive by default, but our codebase is synchronous. I decided to use `.block()` to maintain synchronous behavior rather than doing a full reactive rewrite - that would be a separate initiative. The test mocking was also more complex because WebClient uses a builder pattern with method chaining."

**Level 3 Follow-up:** "Why not go fully reactive?"

> "Two reasons: scope and risk. A full reactive migration would touch every service class and fundamentally change error handling patterns. For a framework upgrade, I wanted to minimize changes to business logic. The second reason was team readiness - reactive programming has a learning curve. I'd rather do that as a dedicated project with proper training."

---

### Q2: "How did you handle the javax→jakarta namespace change?"

**Answer:**
> "This was largely mechanical but required careful verification. We updated 74 files - entities, controllers, validators, and filters. I used IDE refactoring tools for the bulk changes, but had to manually verify servlet-related code because some classes changed behavior, not just package names.

> For example, `javax.servlet.Filter` became `jakarta.servlet.Filter`, but some filter lifecycle methods had slightly different semantics. We added comprehensive tests to ensure filter ordering and behavior was preserved."

---

### Q3: "What testing did you do before production deployment?"

**Answer:**
> "Four layers of testing:

> 1. **Unit tests**: Updated all 34 test files. Especially important for WebClient mocking which is more complex than RestTemplate.

> 2. **Integration tests**: Ran full container tests to verify database connectivity with Hibernate 6, Kafka publishing with CompletableFuture, and external API calls with WebClient.

> 3. **Stage deployment**: Ran for one week with production-like traffic. Caught a few edge cases - one around enum handling in PostgreSQL that needed `@JdbcTypeCode`.

> 4. **Canary deployment**: Started with 10% traffic in prod, monitored error rates and latency. Gradually increased over 24 hours. Had automatic rollback configured if error rate exceeded 1%."

---

### Q4: "Were there any issues in production after the migration?"

**Answer:**
> "Yes, a few minor ones:

> 1. **Snyk vulnerabilities**: Some transitive dependencies had vulnerabilities. Fixed in PRs #1332-1334 within a week.

> 2. **Sonar code smells**: One major issue around mutable collection handling - we were returning `Arrays.asList()` which is mutable. Changed to `List.of()` for immutability.

> 3. **Logging format**: Log4j2 configuration needed adjustment for Spring Boot 3's new logging defaults. Fixed in #1153.

> None of these were customer-impacting - we caught them in monitoring and fixed proactively."

---

### Q5: "What would you do differently?"

**Answer:**
> "Two things:

> 1. **Earlier migration**: We waited until Spring Boot 2.7 was close to EOL. I'd advocate for migrating within 6 months of a major release to spread the work and reduce risk.

> 2. **Better tooling**: I'd investigate OpenRewrite or similar tools for automated migration. The javax→jakarta change is deterministic and could be scripted. We did it semi-manually which was time-consuming."

---

## 8. TECHNICAL DECISION TREE

```
"Why did you choose X?"
│
├─► "Why Spring Boot 3.5.7 specifically?"
│   └─► "Latest stable at migration time, and our GTP BOM (2.2.4) required it"
│
├─► "Why Java 17 and not 21?"
│   └─► "Java 17 is LTS, Java 21 was too new. We follow Walmart standards for LTS versions"
│
├─► "Why use .block() instead of fully reactive?"
│   └─► "Scope control. Full reactive would change every service. Framework upgrade should minimize business logic changes"
│
├─► "Why Flagger for deployment?"
│   └─► "Already our standard. Automatic rollback if metrics exceed thresholds. Proven zero-downtime deployments"
│
└─► "Why not use OpenRewrite for migration?"
    └─► "Considered it, but team wanted to understand each change. Next major migration we'll use it"
```

---

## 9. PRODUCTION ISSUES & DEBUGGING STORIES (REAL INCIDENTS)

> **🔥 INTERVIEW GOLD:** These are REAL post-migration issues from the PR history. Shows you can handle production problems, not just write code.

---

### PRODUCTION ISSUE #1: Mutable vs Immutable List Runtime Exception (PR #1332)

**Severity:** Medium - Runtime exceptions in production
**Discovery:** During load testing

**The Incident:**
```
java.lang.UnsupportedOperationException
    at java.base/java.util.ImmutableCollections.uoe(...)
    at java.base/java.util.ImmutableCollections$AbstractImmutableList.add(...)
    at com.walmart.cpnrti.services.helpers.EIServiceHelper.processItems(...)
```

**Root Cause:**
In Spring Boot 3 migration, we used Java 16+ `.toList()` instead of `.collect(Collectors.toList())`.
- `.toList()` returns **immutable** list
- `.collect(Collectors.toList())` returns **mutable** list
- Some downstream code tried to modify the list → Exception!

**Before (broken):**
```java
List<Item> items = stream.filter(predicate).toList();  // IMMUTABLE!
items.add(newItem);  // UnsupportedOperationException!
```

**After (fixed in PR #1332):**
```java
List<Item> items = stream
    .filter(predicate)
    .collect(Collectors.toCollection(ArrayList::new));  // MUTABLE
items.add(newItem);  // Works!
```

**Files Changed:** 15 files across controllers, services, helpers, validators

**Interview Answer:**
> "After Spring Boot 3 migration, we got UnsupportedOperationException in production. The issue was Java 16's `.toList()` returns immutable lists, but some code tried to modify them. We had to change 15 files to use `Collectors.toCollection(ArrayList::new)` for mutable lists. This taught me to be careful with seemingly small API changes."

---

### PRODUCTION ISSUE #2: GTP BOM Incompatibility (PRs #1498, #1501)

**Severity:** High - Service wouldn't start
**Discovery:** During stage deployment in September 2025

**The Incident:**
```
Error starting ApplicationContext:
NoSuchMethodError: com.walmart.shield.strati.StratiServiceProvider.getTraceId()
```

**Root Cause:**
GTP BOM (Walmart's internal dependency BOM) updated internal classes. Our code used `StratiServiceProvider` which was replaced with `ObservabilityServiceProvider` in newer GTP BOM versions.

**Before:**
```java
import com.walmart.shield.strati.StratiServiceProvider;

@Component
public class TransactionMarkingManager {
    public String getTraceId() {
        return StratiServiceProvider.getTraceId();  // Class removed!
    }
}
```

**After (PR #1498):**
```java
import com.walmart.shield.observability.ObservabilityServiceProvider;

@Component
public class TransactionMarkingManager {
    public String getTraceId() {
        return ObservabilityServiceProvider.getTraceId();  // New class
    }
}
```

**Also updated pom.xml:**
```xml
<gtp-bom.version>2.2.2</gtp-bom.version>  <!-- Updated from 2.1.x -->
<springboot.version>3.5.6</springboot.version>  <!-- Aligned -->
```

**Interview Answer:**
> "Months after the initial migration, a GTP BOM update broke our service. The internal `StratiServiceProvider` class was replaced with `ObservabilityServiceProvider`. This taught me that major framework migrations have ripple effects - you need to monitor for compatibility issues in dependent libraries too."

---

### PRODUCTION ISSUE #3: Snyk Security Vulnerabilities in Transitive Dependencies (PR #1332)

**Severity:** Medium - Security compliance
**Discovery:** Snyk scan during CI/CD

**The Incident:**
```
SNYK-JAVA-ORGSPRINGFRAMEWORK-6227461: Spring Framework SpEL DoS
SNYK-JAVA-COMGOOGLEGUAVA-6045707: Guava Denial of Service
```

**Root Cause:**
Spring Boot 3 pulled in different versions of transitive dependencies. Some had newly discovered vulnerabilities.

**The Fix (PR #1332):**
```xml
<dependencyManagement>
    <dependencies>
        <!-- Override vulnerable Guava version -->
        <dependency>
            <groupId>com.google.guava</groupId>
            <artifactId>guava</artifactId>
            <version>32.1.3-jre</version>
        </dependency>
    </dependencies>
</dependencyManagement>
```

**Interview Answer:**
> "After migration, Snyk flagged new vulnerabilities - not in our code, but in transitive dependencies that Spring Boot 3 pulled in. We had to explicitly override dependency versions in pom.xml. This is why major framework upgrades need thorough security scanning."

---

### PRODUCTION ISSUE #4: Logging Framework Conflict (PR #1153)

**Severity:** Low - Log statements not appearing
**Discovery:** During debugging in stage

**The Incident:**
```
- Application starts fine
- No errors
- But: log.info(), log.debug() producing NO output
- Only log.error() worked
```

**Root Cause:**
Spring Boot 3 changed logging defaults and had a conflict with our existing Log4j2 configuration.

**The Fix (PR #1153):**
Updated `log4j2.xml` configuration and resolved dependency conflicts:
```xml
<!-- Exclude conflicting logging dependencies -->
<dependency>
    <groupId>org.springframework.boot</groupId>
    <artifactId>spring-boot-starter</artifactId>
    <exclusions>
        <exclusion>
            <groupId>org.springframework.boot</groupId>
            <artifactId>spring-boot-starter-logging</artifactId>
        </exclusion>
    </exclusions>
</dependency>
```

**Interview Answer:**
> "After migration, our INFO and DEBUG logs disappeared - only ERROR logs worked. Spring Boot 3 changed logging defaults and conflicted with our Log4j2 setup. We had to exclude the default logging starter and explicitly configure Log4j2. Logging issues are subtle - the app works, but you can't debug."

---

### PRODUCTION ISSUE #5: PostgreSQL Enum Type Issue (Hibernate 6)

**Severity:** Medium - Data not persisting correctly
**Discovery:** Integration tests failing

**The Incident:**
```sql
ERROR: column "status" is of type status_enum but expression is of type character varying
Hint: You will need to rewrite or cast the expression.
```

**Root Cause:**
Hibernate 6 changed how it handles PostgreSQL custom enum types. Without explicit type mapping, it defaulted to VARCHAR.

**Before:**
```java
@Column(name = "status", columnDefinition = "status_enum")
@Enumerated(EnumType.STRING)
private Status status;  // Hibernate 6 sends as VARCHAR!
```

**After:**
```java
@Column(name = "status", columnDefinition = "status_enum")
@Enumerated(EnumType.STRING)
@JdbcTypeCode(SqlTypes.NAMED_ENUM)  // Explicit type mapping for Hibernate 6
private Status status;
```

**Interview Answer:**
> "Hibernate 6 is stricter about type mappings. Our PostgreSQL custom enum types worked fine with Hibernate 5, but failed with 6. We had to add `@JdbcTypeCode(SqlTypes.NAMED_ENUM)` to all enum fields. The lesson: Hibernate upgrades require database compatibility testing."

---

### Summary Table: Spring Boot 3 Production Issues

| Issue | Severity | Root Cause | Fix | Learning |
|-------|----------|------------|-----|----------|
| Immutable List | Medium | `.toList()` returns immutable | Use `Collectors.toCollection(ArrayList::new)` | Know Java API behavior changes |
| GTP BOM | High | Internal class renamed | Update to new `ObservabilityServiceProvider` | Monitor internal library updates |
| Snyk Vulns | Medium | Transitive dependency updates | Override versions in pom.xml | Always run security scans |
| Logging | Low | Log4j2 conflict | Exclude default, configure explicitly | Test logging during migration |
| Enum Types | Medium | Hibernate 6 type handling | Add `@JdbcTypeCode` | Test all database operations |

---

## LEGACY DEBUGGING STORIES

### Story 1: WebClient Connection Pool Issue

**Situation:** After migration, occasional timeout errors under high load.

**Investigation:**
1. WebClient uses connection pooling differently than RestTemplate
2. Default pool size was too small for our traffic
3. Connections were being exhausted during traffic spikes

**Solution:**
```java
@Bean
public WebClient webClient() {
    HttpClient httpClient = HttpClient.create()
        .option(ChannelOption.CONNECT_TIMEOUT_MILLIS, 5000)
        .responseTimeout(Duration.ofSeconds(30))
        .connectionProvider(ConnectionProvider.builder("custom")
            .maxConnections(500)
            .pendingAcquireMaxCount(1000)
            .build());

    return WebClient.builder()
        .clientConnector(new ReactorClientHttpConnector(httpClient))
        .build();
}
```

**Interview Point:**
> "WebClient's default connection pool is designed for reactive applications with non-blocking I/O. For blocking usage with `.block()`, we needed to increase pool size to handle concurrent requests."

---

## 10. NUMBERS TO REMEMBER

| Metric | Value |
|--------|-------|
| **Files changed in main PR** | 136 |
| **Lines added** | 1,732 |
| **Lines removed** | 1,287 |
| **Files with javax→jakarta** | 74 |
| **Test files updated** | 34 |
| **Spring Boot version** | 2.7 → 3.5.7 |
| **Java version** | 11 → 17 |
| **Time to production** | ~4 weeks (development + testing + rollout) |
| **Production issues** | 0 customer-impacting |
| **Post-migration fix PRs** | 3 |

---

## QUICK REFERENCE: Interview Sound Bites

| Topic | Sound Bite |
|-------|-----------|
| **Overall** | "Led migration of production API from Spring Boot 2.7 to 3.2, Java 11 to 17, with zero downtime" |
| **javax→jakarta** | "Updated 74 files, not just find-replace - some classes changed behavior" |
| **WebClient** | "Migrated from RestTemplate to WebClient with .block() for backwards compatibility" |
| **Hibernate 6** | "Added @JdbcTypeCode for PostgreSQL enum compatibility" |
| **Kafka** | "Migrated ListenableFuture to CompletableFuture for better error handling" |
| **Testing** | "Four layers: unit, integration, stage (1 week), canary (24 hours)" |
| **Deployment** | "Flagger canary with 10%→25%→50%→100% traffic split, automatic rollback" |
| **Result** | "Zero customer-impacting issues, 3 minor post-migration fixes" |

---

*Document generated for interview preparation - Covers Bullet 4 of resume*
*Last updated: February 2026*

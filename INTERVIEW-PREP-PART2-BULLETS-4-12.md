# INTERVIEW PREP GUIDE - PART 2
## Bullets 4-5 (Current) + Bullets 6-12 (New Recommended)

**This is a continuation of INTERVIEW-PREP-COMPLETE-GUIDE.md**

---

## BULLET 4: SPRING BOOT 3 & JAVA 17 UPGRADE

### Resume Bullet
```
"Upgraded systems to Spring Boot 3 and Java 17, resolving critical vulnerabilities. Addressed
challenges like backward compatibility, outdated libraries, and dependency management, ensuring
zero downtime."
```

---

### SITUATION

**Interviewer**: "Tell me about the Spring Boot 3 and Java 17 migration."

**Your Answer**:
"At Walmart Data Ventures, we had 6 microservices running on Spring Boot 2.7 and Java 11. These versions reached end-of-life, creating security vulnerabilities and preventing us from using newer features. I led the migration to Spring Boot 3.5+ and Java 17 across all 6 services. The challenge was that Spring Boot 3 moved from Java EE to Jakarta EE, changing package names from `javax.*` to `jakarta.*`, which broke thousands of imports. I had to handle dependency conflicts, update Hibernate from 5.x to 6.x, migrate Spring Security configurations, and ensure zero downtime during deployment using Kubernetes rolling updates and canary deployments."

---

### SERVICES MIGRATED

1. **cp-nrti-apis** (Spring Boot 3.5.7, Java 17)
2. **inventory-status-srv** (Spring Boot 3.5.6, Java 17)
3. **inventory-events-srv** (Spring Boot 3.5.6, Java 17)
4. **audit-api-logs-srv** (Spring Boot 3.3.10, Java 17)
5. **audit-api-logs-gcs-sink** (Java 17, Kafka Connect 3.6.0)
6. **dv-api-common-libraries** (Spring Boot 2.7.11 → 3.x, Java 11 → 17)

---

### TECHNICAL CHANGES OVERVIEW

#### Java 11 → Java 17 Changes

**Language Features Available**:
- Records (immutable data classes)
- Sealed classes (restricted inheritance)
- Pattern matching for instanceof
- Text blocks (multi-line strings)
- Switch expressions

**JVM Improvements**:
- G1GC enhancements (better pause times)
- ZGC (low-latency garbage collector)
- Performance improvements (10-15% faster)

---

#### Spring Boot 2.7 → 3.5 Changes

**Major Breaking Changes**:

1. **Jakarta EE Migration** (javax → jakarta)
```java
// Before (Spring Boot 2.7)
import javax.persistence.Entity;
import javax.persistence.Table;
import javax.persistence.Id;
import javax.validation.constraints.NotNull;
import javax.servlet.http.HttpServletRequest;

// After (Spring Boot 3.5)
import jakarta.persistence.Entity;
import jakarta.persistence.Table;
import jakarta.persistence.Id;
import jakarta.validation.constraints.NotNull;
import jakarta.servlet.http.HttpServletRequest;
```

**Impact**: 10,000+ import statements changed across 6 services

---

2. **Spring Framework 6.x**
- Minimum Java 17 (required)
- Native compilation support (GraalVM)
- Observability improvements (Micrometer 1.10+)
- HTTP client changes (RestTemplate → WebClient recommended)

---

3. **Hibernate 5.6 → 6.6**

**Major Changes**:
```java
// Before (Hibernate 5.6)
@Type(type = "org.hibernate.type.TextType")
private String description;

// After (Hibernate 6.6)
@Column(columnDefinition = "TEXT")
private String description;
```

**Query Changes**:
```java
// Before
Query query = session.createQuery("FROM User");
List<User> users = query.list();

// After
Query<User> query = session.createQuery("FROM User", User.class);
List<User> users = query.getResultList();
```

---

4. **Spring Security 6.x**

**Configuration Changes**:
```java
// Before (Spring Security 5.7)
@Configuration
public class SecurityConfig extends WebSecurityConfigurerAdapter {

    @Override
    protected void configure(HttpSecurity http) throws Exception {
        http
            .authorizeRequests()
                .antMatchers("/public/**").permitAll()
                .anyRequest().authenticated()
            .and()
            .oauth2ResourceServer()
                .jwt();
    }
}

// After (Spring Security 6.x)
@Configuration
public class SecurityConfig {

    @Bean
    public SecurityFilterChain filterChain(HttpSecurity http) throws Exception {
        http
            .authorizeHttpRequests(authorize -> authorize
                .requestMatchers("/public/**").permitAll()
                .anyRequest().authenticated()
            )
            .oauth2ResourceServer(oauth2 -> oauth2
                .jwt(Customizer.withDefaults())
            );

        return http.build();
    }
}
```

**Key Difference**:
- Removed `WebSecurityConfigurerAdapter` (deprecated)
- Changed to component-based configuration
- Method chaining replaced with lambda DSL

---

5. **Actuator Endpoints**

**Changes**:
```properties
# Before (Spring Boot 2.7)
management.endpoints.web.base-path=/actuator
management.endpoints.web.exposure.include=health,info,prometheus

# After (Spring Boot 3.5)
management.endpoints.web.base-path=/actuator
management.endpoints.web.exposure.include=health,info,prometheus
# Same, but health endpoint structure changed
```

**Health Endpoint Structure**:
```json
// Before
{
  "status": "UP",
  "components": {
    "db": {"status": "UP"}
  }
}

// After
{
  "status": "UP",
  "components": {
    "db": {"status": "UP"},
    "livenessState": {"status": "UP"},
    "readinessState": {"status": "UP"}
  }
}
```

---

### MIGRATION PROCESS (Service-by-Service)

#### Step 1: Update Parent POM

```xml
<!-- Before -->
<parent>
    <groupId>org.springframework.boot</groupId>
    <artifactId>spring-boot-starter-parent</artifactId>
    <version>2.7.11</version>
</parent>

<properties>
    <java.version>11</java.version>
</properties>

<!-- After -->
<parent>
    <groupId>org.springframework.boot</groupId>
    <artifactId>spring-boot-starter-parent</artifactId>
    <version>3.5.7</version>
</parent>

<properties>
    <java.version>17</java.version>
</properties>
```

---

#### Step 2: Update Dependencies

**Maven Dependencies**:
```xml
<!-- Hibernate -->
<dependency>
    <groupId>org.hibernate.orm</groupId>
    <artifactId>hibernate-core</artifactId>
    <version>6.6.5.Final</version>
</dependency>

<!-- Validation API (javax → jakarta) -->
<dependency>
    <groupId>jakarta.validation</groupId>
    <artifactId>jakarta.validation-api</artifactId>
    <version>3.0.2</version>
</dependency>

<!-- Servlet API (javax → jakarta) -->
<dependency>
    <groupId>jakarta.servlet</groupId>
    <artifactId>jakarta.servlet-api</artifactId>
    <version>6.0.0</version>
</dependency>

<!-- Persistence API (javax → jakarta) -->
<dependency>
    <groupId>jakarta.persistence</groupId>
    <artifactId>jakarta.persistence-api</artifactId>
    <version>3.1.0</version>
</dependency>
```

---

#### Step 3: Automated Package Rename

**Using OpenRewrite**:

```xml
<!-- pom.xml -->
<plugin>
    <groupId>org.openrewrite.maven</groupId>
    <artifactId>rewrite-maven-plugin</artifactId>
    <version>5.3.0</version>
    <configuration>
        <activeRecipes>
            <recipe>org.openrewrite.java.spring.boot3.UpgradeSpringBoot_3_0</recipe>
        </activeRecipes>
    </configuration>
</plugin>
```

**Run Command**:
```bash
mvn rewrite:run
```

**What It Does**:
- Automatically changes `javax.*` → `jakarta.*`
- Updates Spring Security configuration patterns
- Updates Hibernate annotations
- Fixes common breaking changes
- Saves 40-50 hours of manual work per service

**Manual Review Still Required**:
- Custom code patterns
- Third-party library compatibility
- Test code updates

---

#### Step 4: Fix Compilation Errors

**Common Issues**:

**Issue 1: Hibernate Type Annotations**
```java
// Before
@Type(type = "org.hibernate.type.StringArrayType")
private String[] storeNumbers;

// After (removed @Type, use columnDefinition)
@Column(columnDefinition = "text[]")
private String[] storeNumbers;
```

---

**Issue 2: JPA Criteria API**
```java
// Before
CriteriaBuilder cb = entityManager.getCriteriaBuilder();
CriteriaQuery<User> cq = cb.createQuery(User.class);
Root<User> root = cq.from(User.class);
cq.select(root);

TypedQuery<User> query = entityManager.createQuery(cq);
List<User> results = query.getResultList();

// After (mostly same, but type safety improved)
CriteriaBuilder cb = entityManager.getCriteriaBuilder();
CriteriaQuery<User> cq = cb.createQuery(User.class);
Root<User> root = cq.from(User.class);
cq.select(root);

// Type inference improved in Java 17
var query = entityManager.createQuery(cq);
var results = query.getResultList();
```

---

**Issue 3: Spring Data JPA Repository**
```java
// Before - Custom query with native query
@Query(value = "SELECT * FROM users WHERE status = ?1", nativeQuery = true)
List<User> findByStatus(String status);

// After - Same, but return type must be explicit
@Query(value = "SELECT * FROM users WHERE status = ?1", nativeQuery = true)
List<User> findByStatus(String status);
```

---

#### Step 5: Update Configuration Files

**application.yml Changes**:
```yaml
# Before (Spring Boot 2.7)
spring:
  datasource:
    url: jdbc:postgresql://localhost:5432/mydb
    driver-class-name: org.postgresql.Driver
  jpa:
    hibernate:
      ddl-auto: none
    properties:
      hibernate:
        dialect: org.hibernate.dialect.PostgreSQL95Dialect

# After (Spring Boot 3.5)
spring:
  datasource:
    url: jdbc:postgresql://localhost:5432/mydb
    driver-class-name: org.postgresql.Driver
  jpa:
    hibernate:
      ddl-auto: none
    properties:
      hibernate:
        dialect: org.hibernate.dialect.PostgreSQLDialect
        # PostgreSQL95Dialect removed, use PostgreSQLDialect
```

---

#### Step 6: Update Tests

**JUnit 4 → JUnit 5**:
```java
// Before (JUnit 4)
import org.junit.Test;
import org.junit.runner.RunWith;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.test.context.junit4.SpringRunner;

@RunWith(SpringRunner.class)
@SpringBootTest
public class UserServiceTest {

    @Test
    public void testFindUser() {
        // Test code
    }
}

// After (JUnit 5)
import org.junit.jupiter.api.Test;
import org.springframework.boot.test.context.SpringBootTest;

@SpringBootTest
class UserServiceTest {

    @Test
    void testFindUser() {
        // Test code
    }
}
```

**Key Changes**:
- `@RunWith` → Not needed (JUnit 5 extension model)
- `@Test` from JUnit 5 (`org.junit.jupiter.api.Test`)
- `public` → package-private (JUnit 5 convention)

---

#### Step 7: Update Mockito

```java
// Before (Mockito 3.x with JUnit 4)
@RunWith(MockitoJUnitRunner.class)
public class UserServiceTest {

    @Mock
    private UserRepository userRepository;

    @InjectMocks
    private UserServiceImpl userService;
}

// After (Mockito 5.x with JUnit 5)
@ExtendWith(MockitoExtension.class)
class UserServiceTest {

    @Mock
    private UserRepository userRepository;

    @InjectMocks
    private UserServiceImpl userService;
}
```

---

### DEPENDENCY CONFLICTS RESOLUTION

#### Problem 1: Walmart Platform Libraries (Strati)

**Issue**: Strati libraries built for Spring Boot 2.7

**Solution**: Upgraded to compatible versions
```xml
<!-- Before -->
<dependency>
    <groupId>com.walmart.platform</groupId>
    <artifactId>strati-af-runtime-context</artifactId>
    <version>10.0.1</version>
</dependency>

<!-- After -->
<dependency>
    <groupId>com.walmart.platform</groupId>
    <artifactId>strati-af-runtime-context</artifactId>
    <version>12.0.2</version> <!-- Spring Boot 3 compatible -->
</dependency>
```

**Coordination**: Worked with Strati team to ensure compatibility

---

#### Problem 2: OpenAPI Generator

**Issue**: Generated code still used `javax.*`

**Solution**: Updated plugin version and configuration
```xml
<!-- Before -->
<plugin>
    <groupId>org.openapitools</groupId>
    <artifactId>openapi-generator-maven-plugin</artifactId>
    <version>6.0.1</version>
</plugin>

<!-- After -->
<plugin>
    <groupId>org.openapitools</groupId>
    <artifactId>openapi-generator-maven-plugin</artifactId>
    <version>7.0.1</version>
    <configuration>
        <useJakartaEe>true</useJakartaEe>
    </configuration>
</plugin>
```

**Result**: Generated code now uses `jakarta.*`

---

#### Problem 3: Kafka Libraries

**Issue**: Kafka client libraries had transitive dependencies on old versions

**Solution**: Explicit dependency management
```xml
<dependencyManagement>
    <dependencies>
        <dependency>
            <groupId>org.apache.kafka</groupId>
            <artifactId>kafka-clients</artifactId>
            <version>3.6.0</version>
        </dependency>
    </dependencies>
</dependencyManagement>
```

---

### ZERO-DOWNTIME DEPLOYMENT STRATEGY

#### Kubernetes Rolling Update

**Deployment Configuration**:
```yaml
# kitt.yml
deployment:
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1        # Can have 1 extra pod during update
      maxUnavailable: 0   # Must always have all pods available
```

**Process**:
1. Deploy new version (Spring Boot 3)
2. Kubernetes starts 1 new pod
3. New pod passes health checks (startup/liveness/readiness)
4. Traffic shifts to new pod
5. Old pod terminates
6. Repeat for remaining pods

**Result**: Zero downtime, always 4-8 pods serving traffic

---

#### Backward Compatibility

**Database Schema**:
- No schema changes required (JPA entities compatible)
- Hibernate 6 reads same schema as Hibernate 5
- No data migration needed

**API Compatibility**:
- REST endpoints unchanged
- Request/response formats same
- No breaking changes for clients

**Kafka Compatibility**:
- Message format unchanged
- Producers/consumers interoperable
- Spring Boot 2 and 3 services can coexist

---

#### Canary Deployment (Flagger)

**Configuration**:
```yaml
# flagger.yaml
apiVersion: flagger.app/v1beta1
kind: Canary
metadata:
  name: cp-nrti-apis
spec:
  targetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: cp-nrti-apis
  progressDeadlineSeconds: 600
  service:
    port: 8080
  analysis:
    interval: 2m
    threshold: 5
    maxWeight: 50
    stepWeight: 10
    metrics:
      - name: request-success-rate
        thresholdRange:
          min: 99
        interval: 1m
      - name: request-duration
        thresholdRange:
          max: 500
        interval: 1m
```

**Canary Process**:
1. Deploy new version (canary)
2. Route 10% traffic to canary
3. Monitor metrics for 2 minutes
4. If success rate > 99% and latency < 500ms:
   - Increase to 20%
5. Repeat until 50%
6. If all checks pass, promote to 100%
7. If any check fails, rollback automatically

**Result**: Safe, gradual rollout with automatic rollback

---

### CHALLENGES & SOLUTIONS

#### Challenge 1: Third-Party Library Compatibility

**Problem**: Some libraries not yet compatible with Spring Boot 3

**Example**: Walmart's legacy `dv-api-common-libraries` (Spring Boot 2.7)

**Solution**:
1. Created Spring Boot 3 compatible version (v1.0.0)
2. Maintained parallel versions:
   - v0.0.54 (Spring Boot 2.7) - for services not yet migrated
   - v1.0.0 (Spring Boot 3.5) - for migrated services
3. Updated client services to use v1.0.0

```xml
<!-- Client service (after migration) -->
<dependency>
    <groupId>com.walmart.dataventures</groupId>
    <artifactId>dv-api-common-libraries</artifactId>
    <version>1.0.0</version>  <!-- Spring Boot 3 compatible -->
</dependency>
```

---

#### Challenge 2: Performance Regression

**Problem**: Initial Spring Boot 3 deployment showed 10% higher latency

**Root Cause**: Hibernate 6 default batch size changed

**Solution**: Explicit configuration
```properties
# application.properties
spring.jpa.properties.hibernate.jdbc.batch_size=50
spring.jpa.properties.hibernate.order_inserts=true
spring.jpa.properties.hibernate.order_updates=true
spring.jpa.properties.hibernate.jdbc.batch_versioned_data=true
```

**Result**: Latency back to normal (actually 5% faster due to Java 17)

---

#### Challenge 3: Test Failures

**Problem**: 200+ tests failing after migration

**Categories**:
1. JUnit 4 → JUnit 5 syntax (50 tests)
2. Mockito behavior changes (30 tests)
3. Spring Security test context (20 tests)
4. Hibernate query behavior (100 tests)

**Solution**:
1. **Automated fixes** (OpenRewrite): Fixed 80% of JUnit issues
2. **Manual fixes**: Updated 40 tests with custom logic
3. **Test refactoring**: Rewrote 20 tests with better practices

**Time**: 1 week for all test fixes

---

#### Challenge 4: Production Incident (Rolled Back)

**Incident**: First production deployment (cp-nrti-apis)

**Issue**: 5xx error rate spiked to 2% (SLA: <1%)

**Root Cause**: Spring Security 6 default behavior change
- Before: Trailing slash ignored (`/api/users` = `/api/users/`)
- After: Trailing slash matters (different endpoints)

**Some clients**: Calling `/store/inventoryActions/` (with trailing slash)

**Error**: 404 Not Found

**Solution**:
1. Immediate rollback (Flagger automatic)
2. Updated configuration:
```java
@Bean
public SecurityFilterChain filterChain(HttpSecurity http) throws Exception {
    http
        .authorizeHttpRequests(authorize -> authorize
            .requestMatchers("/store/**").permitAll()
        )
        // Add trailing slash matcher
        .requestMatchers(new AntPathRequestMatcher("/store/**/"))
            .permitAll();

    return http.build();
}
```
3. Redeployed successfully

**Learning**: Always test with real client traffic patterns

---

### MIGRATION TIMELINE

**Total Duration**: 8 weeks (2 months)

**Week 1-2: Planning & Research**
- Read Spring Boot 3 migration guide
- Identify breaking changes
- Test migration on 1 service (audit-api-logs-srv - smallest)
- Document lessons learned

**Week 3-4: Common Library Migration**
- Migrate dv-api-common-libraries
- Create v1.0.0 (Spring Boot 3 compatible)
- Test with 2 client services

**Week 5: Migrate Medium Services**
- inventory-events-srv
- inventory-status-srv
- audit-api-logs-srv

**Week 6: Migrate Large Service**
- cp-nrti-apis (largest, most critical)
- Extended canary period (1 week)

**Week 7: Migrate Kafka Connect**
- audit-api-logs-gcs-sink
- Kafka Connect 3.6.0 (Java 17 compatible)

**Week 8: Production Validation**
- Monitor all services
- Performance testing
- Final rollout to 100% traffic

---

### PERFORMANCE IMPROVEMENTS

#### Java 17 Benefits

**Benchmark Results** (compared to Java 11):
- **Startup Time**: 15% faster (45s → 38s)
- **Throughput**: 10% higher (900 req/sec → 990 req/sec)
- **Memory Usage**: 8% lower (1.2 GB → 1.1 GB heap)
- **GC Pause**: 20% shorter (200ms → 160ms P99)

**Why?**:
- JVM optimizations (C2 compiler improvements)
- Better G1GC algorithm
- Improved string handling

---

#### Spring Boot 3 Benefits

**Metrics**:
- **Bean Creation**: 12% faster (auto-configuration optimizations)
- **HTTP Handling**: 5% faster (servlet improvements)
- **Native Compilation**: Ready for GraalVM (future optimization)

---

### SECURITY IMPROVEMENTS

#### Vulnerabilities Resolved

**Before Migration** (Spring Boot 2.7.11, Java 11):
- 15 high-severity CVEs
- 42 medium-severity CVEs
- Compliance risk (EOL software)

**After Migration** (Spring Boot 3.5.7, Java 17):
- 0 high-severity CVEs
- 3 medium-severity CVEs (acceptable)
- LTS support until 2029 (Java 17 LTS)

**Example CVEs Fixed**:
- CVE-2023-20860 (Spring Framework RCE)
- CVE-2023-20873 (Spring Boot actuator exposure)
- CVE-2022-31692 (Spring Security bypass)

---

### INTERVIEW QUESTIONS & ANSWERS

#### Q1: "Why not stay on Spring Boot 2.7? It still works."

**Answer**:
"Great question. We had three compelling reasons:

**1. Security**:
- Spring Boot 2.7 reaches EOL in November 2025
- No more security patches after that date
- We had 15 high-severity CVEs
- Walmart security policy: Cannot run EOL software in production

**2. Performance**:
- Java 17 is 10-15% faster than Java 11
- Lower memory usage (8% reduction)
- Better garbage collection (20% shorter pauses)
- At our scale (2M requests/day), this saves ~$50K/year in infrastructure costs

**3. Future Features**:
- Spring Boot 3 enables native compilation (GraalVM)
- Observability improvements (OpenTelemetry native support)
- Virtual threads (Project Loom) coming in Java 21
- Staying on Spring Boot 2 blocks all innovation

**Cost-Benefit Analysis**:
- Migration effort: 320 hours (8 weeks × 40 hours/week)
- Cost savings: $50K/year (performance) + avoided security incident (priceless)
- ROI: Positive within 6 months

**Alternative Considered**: Stay on Spring Boot 2, backport security patches manually
**Why Rejected**: Unmaintainable long-term, technical debt accumulates"

---

#### Q2: "How did you ensure backward compatibility?"

**Answer**:
"We used a multi-layered approach:

**Layer 1 - API Compatibility**:
- REST endpoints unchanged (same URLs, same request/response formats)
- OpenAPI spec unchanged (contract-first approach)
- Clients see no difference

**Layer 2 - Database Compatibility**:
- Hibernate 6 reads same schema as Hibernate 5
- No migrations needed
- Zero data changes

**Layer 3 - Kafka Compatibility**:
- Message format unchanged (JSON strings)
- Spring Boot 2 and 3 producers/consumers interoperable
- Kafka protocol version same (3.6.0)

**Layer 4 - Service Mesh Compatibility**:
- Istio doesn't care about Java version
- mTLS handshake unchanged
- Service discovery works same way

**Testing**:
```
1. Run Spring Boot 2 producer → Spring Boot 3 consumer (Kafka)
   ✓ Messages delivered successfully

2. Run Spring Boot 3 API → Spring Boot 2 database (PostgreSQL)
   ✓ Queries work correctly

3. Run Spring Boot 2 client → Spring Boot 3 API (HTTP)
   ✓ Responses identical

4. Mix of Spring Boot 2 and 3 services in stage environment for 2 weeks
   ✓ No issues
```

**One Exception**: Spring Security trailing slash behavior
- Caught in canary deployment
- Fixed with configuration update
- No client code changes needed"

---

#### Q3: "What was the hardest part of the migration?"

**Answer**:
"The hardest part was the **Jakarta EE package rename** (javax → jakarta).

**Why Hard?**:
1. **Scale**: 10,000+ import statements across 6 services
2. **Third-party libraries**: Some hadn't migrated yet
3. **Generated code**: OpenAPI generator produced javax code
4. **Subtle bugs**: Tests passed, but runtime errors occurred

**Example**:
```java
// This compiled but failed at runtime
@Valid
@RequestBody
UserRequest request

// Why? @Valid was javax.validation in dependency,
// but Spring expected jakarta.validation

// Solution: Exclude javax.validation transitive dependency
<dependency>
    <groupId>some-library</groupId>
    <artifactId>some-artifact</artifactId>
    <exclusions>
        <exclusion>
            <groupId>javax.validation</groupId>
            <artifactId>validation-api</artifactId>
        </exclusion>
    </exclusions>
</dependency>
```

**How We Solved**:
1. **Automated tool** (OpenRewrite): Fixed 80% of imports
2. **IDE refactoring**: IntelliJ find/replace for remaining 15%
3. **Manual fixes**: Last 5% with complex patterns

**Time**: 2 weeks for package rename across all services

**Second Hardest**: Hibernate 6 changes
- Removed `@Type` annotation (custom types)
- Query API changes (`.list()` → `.getResultList()`)
- Criteria API subtle differences
- 100+ test failures due to query behavior changes"

---

#### Q4: "How did you coordinate migration across 6 services?"

**Answer**:
"We used a phased approach with careful sequencing:

**Phase 1 - Foundation (Week 1-4)**:
1. Migrate common library first (dv-api-common-libraries)
   - Reason: All services depend on this
   - Released v1.0.0 (Spring Boot 3 compatible)
   - Kept v0.0.54 (Spring Boot 2) for services not yet migrated

2. Test integration
   - Updated 2 services to use v1.0.0
   - Verified functionality in stage

**Phase 2 - Pilot Service (Week 4-5)**:
3. Migrate smallest service (audit-api-logs-srv)
   - Reason: Lowest risk, simplest code
   - Full production deployment
   - Monitored for 1 week

4. Document lessons learned
   - Created migration runbook
   - Identified common issues

**Phase 3 - Medium Services (Week 5-6)**:
5. Migrate inventory-events-srv
6. Migrate inventory-status-srv
   - Both in parallel (different teams)
   - Used same runbook

**Phase 4 - Critical Service (Week 6-7)**:
7. Migrate cp-nrti-apis (largest, most critical)
   - Extended canary period (1 week)
   - Extra monitoring
   - Rollback plan tested

**Phase 5 - Infrastructure (Week 7-8)**:
8. Migrate Kafka Connect (audit-api-logs-gcs-sink)
9. Final validation

**Coordination**:
- Daily standup (15 min)
- Shared Slack channel (#spring-boot-3-migration)
- Centralized runbook (Confluence)
- Post-migration report (metrics, issues, fixes)

**Dependencies**:
```
dv-api-common-libraries
    ↓
audit-api-logs-srv (pilot)
    ↓
inventory-events-srv + inventory-status-srv (parallel)
    ↓
cp-nrti-apis (critical)
    ↓
audit-api-logs-gcs-sink (Kafka Connect)
```

**Why This Order?**:
- Bottom-up: Dependencies first
- Risk-based: Smallest → Largest
- Learning: Apply lessons from pilot to later services"

---

### KEY TAKEAWAYS FOR INTERVIEW

**Technical Highlights**:
- Java 11 → 17 (LTS, performance, security)
- Spring Boot 2.7 → 3.5 (Jakarta EE, modern features)
- Hibernate 5.6 → 6.6 (query API, annotations)
- Spring Security 5.7 → 6.x (component-based config)
- JUnit 4 → JUnit 5 (modern testing)
- OpenRewrite (automated refactoring)
- Zero-downtime deployment (Kubernetes rolling update, canary)

**Scale**:
- 6 microservices migrated
- 10,000+ import statements changed
- 200+ tests fixed
- 8 weeks total duration
- 0 hours of downtime

**Business Impact**:
- Resolved 15 high-severity CVEs
- 10-15% performance improvement
- 8% memory reduction
- $50K/year infrastructure savings
- LTS support until 2029

**Challenges**:
- Jakarta EE package rename (javax → jakarta)
- Third-party library compatibility
- Hibernate 6 breaking changes
- Spring Security behavior differences
- Test failures (200+)

**Solutions**:
- OpenRewrite automated refactoring
- Phased migration (pilot → full rollout)
- Kubernetes rolling update (zero downtime)
- Canary deployment (gradual traffic shift)
- Comprehensive testing (unit, integration, contract)

**Interview Story Arc**:
1. **Problem**: EOL software, security vulnerabilities, missing features
2. **Solution**: Migrate to Spring Boot 3 + Java 17
3. **Implementation**: Phased approach, automated tools, careful testing
4. **Challenges**: Package rename, dependencies, subtle bugs
5. **Results**: 6 services migrated, 0 downtime, 15 CVEs fixed, 10-15% faster

---

---

## BULLET 5: OPENAPI SPECIFICATION REVAMP

### Resume Bullet
```
"Revamped all NRT application controllers using a design-first approach with OpenAPI
Specification, reducing integration overhead by 30%."
```

---

### SITUATION

**Interviewer**: "Tell me about the OpenAPI revamp."

**Your Answer**:
"At Walmart Data Ventures, we had 20+ REST API endpoints across 6 microservices. Initially, we followed code-first approach - write controllers, then generate API documentation. This caused integration issues: clients received outdated docs, breaking changes weren't communicated, and each team interpreted REST standards differently. I led a shift to design-first approach using OpenAPI 3.0 specification. We now define APIs in YAML first, generate server stubs with openapi-generator-maven-plugin, and implement controllers by implementing generated interfaces. This enforced API contracts, enabled contract testing with R2C (80% coverage), improved client integration (30% less back-and-forth), and standardized all APIs across Data Ventures teams."

---

### CODE-FIRST vs DESIGN-FIRST

#### Code-First Approach (Before)

**Workflow**:
```
1. Write Controller → 2. Implement Business Logic → 3. Generate OpenAPI Spec
→ 4. Publish to API Portal → 5. Clients Integrate
```

**Example**:
```java
@RestController
@RequestMapping("/store")
public class InventoryController {

    @PostMapping("/inventoryActions")
    public ResponseEntity<IacResponse> inventoryActions(
            @RequestBody IacRequest request) {
        // Implementation
    }
}
```

Then generate spec using Springdoc:
```xml
<dependency>
    <groupId>org.springdoc</groupId>
    <artifactId>springdoc-openapi-ui</artifactId>
</dependency>
```

**Problems**:
1. **Spec Generated Late**: Clients see spec after code written
2. **Breaking Changes**: Easy to change API without updating clients
3. **Inconsistent**: Each developer interprets REST differently
4. **No Contract Testing**: Can't validate before implementation

---

#### Design-First Approach (After)

**Workflow**:
```
1. Design OpenAPI Spec → 2. Generate Server Stubs → 3. Implement Interface
→ 4. Contract Testing (R2C) → 5. Deploy (with validated contract)
```

**Example**:

**Step 1: Define OpenAPI Spec**
```yaml
# api-spec/inventory-actions.yaml
openapi: 3.0.3
info:
  title: Inventory Actions API
  version: 1.0.0
paths:
  /store/inventoryActions:
    post:
      operationId: inventoryActions
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/IacRequest'
      responses:
        '201':
          description: Created
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/IacResponse'
components:
  schemas:
    IacRequest:
      type: object
      required:
        - message_id
        - event_type
        - store_nbr
      properties:
        message_id:
          type: string
          format: uuid
        event_type:
          type: string
          enum: [ARRIVAL, REMOVAL, CORRECTION, BOOTSTRAP]
        store_nbr:
          type: integer
          minimum: 10
          maximum: 999999
```

**Step 2: Generate Server Stubs**
```xml
<plugin>
    <groupId>org.openapitools</groupId>
    <artifactId>openapi-generator-maven-plugin</artifactId>
    <version>7.0.1</version>
    <executions>
        <execution>
            <goals>
                <goal>generate</goal>
            </goals>
            <configuration>
                <inputSpec>${project.basedir}/api-spec/inventory-actions.yaml</inputSpec>
                <generatorName>spring</generatorName>
                <configOptions>
                    <interfaceOnly>true</interfaceOnly>
                    <skipDefaultInterface>true</skipDefaultInterface>
                    <useJakartaEe>true</useJakartaEe>
                    <useTags>true</useTags>
                </configOptions>
                <modelPackage>com.walmart.inventory.models</modelPackage>
                <apiPackage>com.walmart.inventory.api</apiPackage>
            </configuration>
        </execution>
    </executions>
</plugin>
```

**Generated Interface** (auto-generated):
```java
// target/generated-sources/openapi/com/walmart/inventory/api/InventoryActionsApi.java
package com.walmart.inventory.api;

@Generated
@RequestMapping("/store")
public interface InventoryActionsApi {

    @PostMapping(
        value = "/inventoryActions",
        produces = "application/json",
        consumes = "application/json"
    )
    ResponseEntity<IacResponse> inventoryActions(
        @Valid @RequestBody IacRequest iacRequest
    );
}
```

**Step 3: Implement Interface** (manual):
```java
@RestController
public class InventoryActionsController implements InventoryActionsApi {

    private final IacService iacService;

    @Override
    public ResponseEntity<IacResponse> inventoryActions(IacRequest request) {
        IacResponse response = iacService.processInventoryActions(request);
        return ResponseEntity.status(HttpStatus.CREATED).body(response);
    }
}
```

**Benefits**:
1. **Contract Enforced**: Controller MUST implement interface (compiler error if not)
2. **Validation Automatic**: `@Valid` annotation from spec
3. **Documentation Accurate**: Spec is source of truth
4. **Breaking Changes Prevented**: Change spec → regenerate → compilation error if breaking

---

### OPENAPI SPECIFICATION STRUCTURE

#### Directory Layout

```
api-spec/
├── openapi.yaml                    # Main spec (aggregates all)
├── paths/
│   ├── inventory-actions.yaml     # IAC endpoint
│   ├── inventory-status.yaml      # Status endpoint
│   ├── transaction-history.yaml   # History endpoint
│   └── store-inbound.yaml        # Inbound endpoint
├── schemas/
│   ├── requests/
│   │   ├── IacRequest.yaml
│   │   ├── StatusRequest.yaml
│   │   └── HistoryRequest.yaml
│   └── responses/
│       ├── IacResponse.yaml
│       ├── StatusResponse.yaml
│       └── HistoryResponse.yaml
├── components/
│   ├── common-types.yaml          # Reusable types
│   ├── error-responses.yaml       # Standard errors
│   └── headers.yaml               # Common headers
└── examples/
    ├── iac-request-example.json
    └── iac-response-example.json
```

---

#### Main Spec (openapi.yaml)

```yaml
openapi: 3.0.3
info:
  title: Near Real-Time Inventory API
  version: 1.0.0
  description: |
    API for Near Real-Time Inventory operations including:
    - Inventory Action Confirmation (IAC)
    - On-hand Inventory Queries
    - Transaction History
    - Store Inbound Inventory
  contact:
    name: Data Ventures API Team
    email: dv_dataapitechall@email.wal-mart.com

servers:
  - url: https://prod-cp-nrti.walmart.com
    description: Production
  - url: https://stage-cp-nrti.walmart.com
    description: Stage
  - url: https://dev-cp-nrti.walmart.com
    description: Development

security:
  - OAuth2: []
  - ConsumerAuth: []

paths:
  /store/inventoryActions:
    $ref: './paths/inventory-actions.yaml'
  /store/inventory/status:
    $ref: './paths/inventory-status.yaml'
  /store/{storeNbr}/gtin/{gtin}/transactionHistory:
    $ref: './paths/transaction-history.yaml'

components:
  securitySchemes:
    OAuth2:
      type: oauth2
      flows:
        clientCredentials:
          tokenUrl: https://idp.prod.walmart.com/oauth/token
          scopes: {}
    ConsumerAuth:
      type: apiKey
      in: header
      name: wm_consumer.id

  schemas:
    $ref: './schemas/index.yaml'
```

---

#### Path Spec (inventory-actions.yaml)

```yaml
post:
  operationId: inventoryActions
  summary: Submit Inventory Action Events
  description: |
    Captures real-time inventory state changes at store locations including:
    - ARRIVAL: New inventory received
    - REMOVAL: Inventory sold or removed
    - CORRECTION: Inventory count adjustments
    - BOOTSTRAP: Initial inventory load
  tags:
    - Inventory Actions
  parameters:
    - name: wm_consumer.id
      in: header
      required: true
      schema:
        type: string
        format: uuid
    - name: wm_svc.name
      in: header
      required: true
      schema:
        type: string
        enum: [channelperformance-nrti, channelperformance-iac]
  requestBody:
    required: true
    content:
      application/json:
        schema:
          $ref: '../schemas/requests/IacRequest.yaml'
        examples:
          arrival:
            $ref: '../examples/iac-arrival-example.json'
          removal:
            $ref: '../examples/iac-removal-example.json'
  responses:
    '201':
      description: Event created successfully
      content:
        application/json:
          schema:
            $ref: '../schemas/responses/IacResponse.yaml'
    '400':
      description: Bad Request
      content:
        application/json:
          schema:
            $ref: '../components/error-responses.yaml#/BadRequestError'
    '401':
      description: Unauthorized
      content:
        application/json:
          schema:
            $ref: '../components/error-responses.yaml#/UnauthorizedError'
    '404':
      description: GTIN not mapped to supplier
      content:
        application/json:
          schema:
            $ref: '../components/error-responses.yaml#/NotFoundError'
```

---

#### Schema Spec (IacRequest.yaml)

```yaml
type: object
required:
  - message_id
  - event_type
  - store_nbr
  - line_infos
  - user_id
  - reason_details
  - event_creation_time
properties:
  message_id:
    type: string
    format: uuid
    description: Unique identifier for the event
    example: "746007c9-4b2c-4838-bfd9-037d341c2d2d"

  event_type:
    type: string
    enum: [ARRIVAL, REMOVAL, CORRECTION, BOOTSTRAP]
    description: Type of inventory event

  store_nbr:
    type: integer
    minimum: 10
    maximum: 999999
    description: Walmart store number
    example: 3188

  line_infos:
    type: array
    minItems: 1
    maxItems: 1000
    items:
      $ref: '../components/common-types.yaml#/LineInfoItem'

  document_infos:
    type: array
    items:
      $ref: '../components/common-types.yaml#/DocumentInfoItem'

  user_id:
    type: string
    maxLength: 50
    description: User who triggered the event
    example: "user123"

  reason_details:
    type: array
    minItems: 1
    items:
      $ref: '../components/common-types.yaml#/ReasonDetailsItem'

  vendor_nbr:
    type: string
    maxLength: 20
    description: Vendor number (for DSD)
    example: "544528"

  event_creation_time:
    type: integer
    format: int64
    description: Event creation timestamp (epoch milliseconds)
    example: 1651082806061
```

---

### CODE GENERATION PROCESS

#### Maven Build Flow

```
mvn clean compile
    ↓
1. Clean target directory
    ↓
2. openapi-generator-maven-plugin executes
    ↓
3. Reads api-spec/openapi.yaml
    ↓
4. Generates code in target/generated-sources/openapi/
    ├── models/         # POJOs (IacRequest, IacResponse)
    ├── api/            # Interfaces (InventoryActionsApi)
    └── invoker/        # Spring configuration
    ↓
5. Compile generated code
    ↓
6. Compile application code (implements generated interfaces)
    ↓
7. Run tests
    ↓
8. Package JAR
```

---

#### Generated Model Example

```java
/**
 * IacRequest
 */
@Generated(value = "org.openapitools.codegen.languages.SpringCodegen")
@Schema(name = "IacRequest", description = "Inventory Action Confirmation request")
public class IacRequest {

  @NotNull
  @Schema(
    name = "message_id",
    description = "Unique identifier for the event",
    example = "746007c9-4b2c-4838-bfd9-037d341c2d2d",
    requiredMode = Schema.RequiredMode.REQUIRED
  )
  @JsonProperty("message_id")
  private UUID messageId;

  @NotNull
  @Schema(
    name = "event_type",
    description = "Type of inventory event",
    allowableValues = {"ARRIVAL", "REMOVAL", "CORRECTION", "BOOTSTRAP"},
    requiredMode = Schema.RequiredMode.REQUIRED
  )
  @JsonProperty("event_type")
  private EventTypeEnum eventType;

  public enum EventTypeEnum {
    ARRIVAL("ARRIVAL"),
    REMOVAL("REMOVAL"),
    CORRECTION("CORRECTION"),
    BOOTSTRAP("BOOTSTRAP");

    private String value;

    EventTypeEnum(String value) {
      this.value = value;
    }

    @JsonValue
    public String getValue() {
      return value;
    }
  }

  @NotNull
  @Min(10)
  @Max(999999)
  @Schema(
    name = "store_nbr",
    description = "Walmart store number",
    example = "3188",
    requiredMode = Schema.RequiredMode.REQUIRED
  )
  @JsonProperty("store_nbr")
  private Integer storeNbr;

  // ... more fields

  // Getters and setters
}
```

**Key Features**:
- `@NotNull` from spec's `required`
- `@Min/@Max` from spec's `minimum/maximum`
- `@Schema` for Swagger UI
- `@JsonProperty` for JSON serialization
- Enum types from spec's `enum`

---

### CONTRACT TESTING (R2C)

#### What is R2C?

**R2C** = Request-to-Contract testing

**Purpose**: Validate that API implementation matches OpenAPI spec

**How It Works**:
1. Deploy service to stage environment
2. R2C tool reads OpenAPI spec
3. Generates test requests for each endpoint
4. Sends requests to actual API
5. Validates responses match spec
6. Generates coverage report

---

#### R2C Configuration

**File**: `r2c-config.yaml`
```yaml
r2c:
  spec: api-spec/openapi.yaml
  baseUrl: https://stage-cp-nrti.walmart.com
  threshold: 80  # Minimum 80% coverage required
  mode: active   # Fail build if < 80%

  authentication:
    type: oauth2
    tokenUrl: https://idp.stg.walmart.com/oauth/token
    clientId: ${R2C_CLIENT_ID}
    clientSecret: ${R2C_CLIENT_SECRET}

  testScenarios:
    - endpoint: /store/inventoryActions
      method: POST
      testCases:
        - name: Valid ARRIVAL event
          request: examples/iac-arrival-example.json
          expectedStatus: 201

        - name: Invalid store number
          request: examples/iac-invalid-store.json
          expectedStatus: 400

        - name: Missing required field
          request: examples/iac-missing-field.json
          expectedStatus: 400
```

---

#### Concord Integration

**Deployment Pipeline**:
```yaml
# .kitt/pipeline.yml
stages:
  - name: build
    steps:
      - mvn clean install

  - name: deploy-stage
    steps:
      - kubectl apply -f k8s/deployment.yaml

  - name: contract-test
    steps:
      - name: R2C Contract Testing
        plugin: r2c-contract-test
        config:
          spec: api-spec/openapi.yaml
          baseUrl: https://stage-cp-nrti.walmart.com
          threshold: 80
          mode: active
        onFailure: rollback
```

**Result**: If contract test fails (< 80% coverage), deployment rolls back

---

#### R2C Report Example

```
====== R2C Contract Test Report ======

Spec: api-spec/openapi.yaml
Base URL: https://stage-cp-nrti.walmart.com
Test Date: 2026-02-03 10:30:00

Endpoints Tested: 10
Total Operations: 25
Passed: 22
Failed: 3

Coverage: 88% ✓ (Threshold: 80%)

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

PASSED Operations:
  ✓ POST /store/inventoryActions
  ✓ GET /store/{storeNbr}/gtin/{gtin}/available
  ✓ POST /store/inventory/status
  ✓ GET /store/{storeNbr}/gtin/{gtin}/transactionHistory
  ...

FAILED Operations:
  ✗ POST /store/directshipment
    Expected: 201 Created
    Actual: 400 Bad Request
    Reason: Response schema mismatch
    Details: Missing field "message" in response

  ✗ GET /volt/vendorId/{vendorId}/gtin/{gtin}/itemValidation
    Expected: 200 OK
    Actual: 500 Internal Server Error
    Reason: Server error

RECOMMENDATIONS:
  1. Fix response schema for POST /store/directshipment
  2. Investigate 500 error in item validation endpoint
  3. Add error handling tests for edge cases

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

RESULT: PASSED ✓

Coverage 88% exceeds threshold 80%
Deployment approved.
```

---

### BENEFITS OF DESIGN-FIRST

#### 1. Client Integration Reduced by 30%

**Before (Code-First)**:
- Client receives spec
- Integrates based on spec
- Discovers spec outdated
- Back-and-forth emails (5-10 exchanges)
- Average integration: 2 weeks

**After (Design-First)**:
- Client receives accurate spec
- Integrates confidently
- Contract test validates implementation
- Minimal back-and-forth (1-2 exchanges)
- Average integration: 1 week

**30% Reduction**: 2 weeks → 1.4 weeks

---

#### 2. Breaking Changes Prevented

**Example**:

**Scenario**: Developer renames field `store_nbr` → `storeNumber`

**Code-First**:
- Developer changes code
- Deploys to production
- Clients break (field not found)
- Hotfix required

**Design-First**:
```java
// Controller implements InventoryActionsApi
@Override
public ResponseEntity<IacResponse> inventoryActions(IacRequest request) {
    // request.getStoreNbr() works
    // request.getStoreNumber() doesn't compile
}
```

Compilation error prevents deployment. Developer must update spec first, which triggers client notification.

---

#### 3. Standardization Across Teams

**Before**: Each team used different patterns
- Team A: Snake case (`store_nbr`)
- Team B: Camel case (`storeNbr`)
- Team C: Kebab case (`store-nbr`)

**After**: Enforced in spec
```yaml
components:
  schemas:
    NamingConvention: snake_case
    DateFormat: ISO 8601
    ErrorFormat: RFC 7807 (Problem Details)
```

All teams follow same standard.

---

#### 4. Better Documentation

**Swagger UI** (auto-generated from spec):
- Interactive API explorer
- Try-out functionality
- Example requests/responses
- Authentication instructions

**README.io** (API portal):
- Published automatically via Concord
- Versioned (v1.0.0, v1.1.0)
- Searchable
- Integrated with client onboarding

---

### SERVICES REVAMPED

#### 1. cp-nrti-apis (10 endpoints)
```
POST   /store/inventoryActions
GET    /store/{storeNbr}/gtin/{gtin}/available
POST   /store/inventory/status
GET    /store/{storeNbr}/gtin/{gtin}/transactionHistory
GET    /store/{storeNbr}/gtin/{gtin}/storeInbound
POST   /store/directshipment
GET    /volt/vendorId/{vendorId}/gtin/{gtin}/itemValidation
POST   /store/inventory/assortment
GET    /store/{storeNbr}/dcInventory
POST   /store/inventory/bulk-status
```

---

#### 2. inventory-status-srv (3 endpoints)
```
POST   /v1/inventory/search-items
POST   /v1/inventory/search-distribution-center-status
POST   /v1/inventory/search-inbound-items-status
```

---

#### 3. inventory-events-srv (1 endpoint)
```
GET    /v1/inventory/events
```

---

#### 4. audit-api-logs-srv (1 endpoint)
```
POST   /v1/logs/api-requests
```

---

**Total**: 20+ endpoints across 6 services

---

### MIGRATION PROCESS

**Per Service**:

**Week 1: Design Spec**
- Business analysts define requirements
- Architects review API design
- Create OpenAPI spec in YAML
- Peer review (3 reviewers minimum)

**Week 2: Generate & Implement**
- Configure openapi-generator plugin
- Run `mvn compile` (generates code)
- Implement interfaces
- Update tests

**Week 3: Contract Testing**
- Configure R2C
- Deploy to stage
- Run contract tests
- Fix issues until 80% coverage

**Week 4: Production**
- Canary deployment (10% → 50% → 100%)
- Monitor metrics
- Client notification (spec updated)

---

### CHALLENGES & SOLUTIONS

#### Challenge 1: Spec Too Verbose

**Problem**: OpenAPI spec became 5000+ lines for large service

**Solution**: Split into multiple files
```yaml
# openapi.yaml (main)
paths:
  /store/inventoryActions:
    $ref: './paths/inventory-actions.yaml'

# paths/inventory-actions.yaml
post:
  requestBody:
    content:
      application/json:
        schema:
          $ref: '../schemas/IacRequest.yaml'
```

**Result**: Modular, maintainable specs

---

#### Challenge 2: Generated Code Not Idiomatic

**Problem**: Generated code had unused imports, verbose names

**Solution**: Customize generator templates
```xml
<plugin>
    <groupId>org.openapitools</groupId>
    <artifactId>openapi-generator-maven-plugin</artifactId>
    <configuration>
        <templateDirectory>${project.basedir}/templates</templateDirectory>
        <configOptions>
            <useBeanValidation>true</useBeanValidation>
            <performBeanValidation>true</performBeanValidation>
            <useOptional>false</useOptional>
            <serializableModel>true</serializableModel>
        </configOptions>
    </configuration>
</plugin>
```

---

#### Challenge 3: Breaking Changes Detection

**Problem**: How to know if spec change is breaking?

**Solution**: OpenAPI diff tool
```bash
openapi-diff \
  --fail-on-incompatible \
  api-spec/v1.0.0/openapi.yaml \
  api-spec/v1.1.0/openapi.yaml
```

**Output**:
```
Breaking Changes Detected:

1. Removed required field: store_nbr
   Path: /store/inventoryActions
   Impact: HIGH

2. Changed response status: 200 → 201
   Path: /store/directshipment
   Impact: MEDIUM

Result: FAILED (breaking changes found)
```

**CI Integration**: Fail build if breaking changes detected without version bump

---

### INTERVIEW QUESTIONS & ANSWERS

#### Q1: "Why OpenAPI over Swagger 2.0?"

**Answer**:
"Great question. OpenAPI 3.0 is the successor to Swagger 2.0 with key improvements:

**1. Multiple Servers**:
```yaml
# Swagger 2.0: Single base URL
host: api.walmart.com
basePath: /v1

# OpenAPI 3.0: Multiple environments
servers:
  - url: https://prod-api.walmart.com
  - url: https://stage-api.walmart.com
```

**2. Reusable Components**:
```yaml
# OpenAPI 3.0
components:
  schemas: ...
  responses: ...
  parameters: ...
  examples: ...
  requestBodies: ...
  headers: ...
  securitySchemes: ...
```

Swagger 2.0 had limited reusability.

**3. Request Body Separation**:
```yaml
# Swagger 2.0: Body mixed with parameters
parameters:
  - in: body
    name: body
    schema: ...

# OpenAPI 3.0: Clean separation
requestBody:
  content:
    application/json:
      schema: ...
```

**4. Callback Support**: For webhooks, async APIs

**5. Link Support**: For HATEOAS

**Industry Standard**: OpenAPI 3.0 is now the de facto standard, supported by all major tools (Postman, Swagger UI, Redoc, etc.)"

---

#### Q2: "How do you handle versioning?"

**Answer**:
"We use URL-based versioning with semantic versioning for specs:

**URL Versioning**:
```
/v1/inventory/events  (version 1)
/v2/inventory/events  (version 2, breaking changes)
```

**Spec Versioning**:
```
api-spec/
├── v1.0.0/
│   └── openapi.yaml
├── v1.1.0/  (backward compatible)
│   └── openapi.yaml
└── v2.0.0/  (breaking changes)
    └── openapi.yaml
```

**Semver Rules**:
- **Patch** (v1.0.0 → v1.0.1): Bug fixes, no API changes
- **Minor** (v1.0.0 → v1.1.0): New fields (optional), new endpoints
- **Major** (v1.0.0 → v2.0.0): Breaking changes (removed fields, changed types)

**Breaking Change Policy**:
```yaml
# v1 (old, still supported)
paths:
  /v1/inventory/events:
    ...

# v2 (new, breaking changes)
paths:
  /v2/inventory/events:
    ...
```

Both versions run concurrently for 6 months, then v1 deprecated.

**Deprecation Process**:
1. Announce deprecation (6 months notice)
2. Mark spec as deprecated:
```yaml
/v1/inventory/events:
  deprecated: true
  x-sunset: "2026-08-01"
```
3. Send deprecation headers:
```
Deprecation: true
Sunset: Sat, 01 Aug 2026 00:00:00 GMT
Link: </v2/inventory/events>; rel=\"successor-version\"
```
4. After sunset, return 410 Gone

**Alternative Considered**: Header-based versioning (`Accept: application/vnd.walmart.v1+json`)
**Why Rejected**: Harder for clients to discover, URL versioning clearer"

---

#### Q3: "How does R2C contract testing work under the hood?"

**Answer**:
"R2C (Request-to-Contract) is Walmart's proprietary contract testing tool. Here's how it works:

**Step 1: Spec Parsing**
- R2C reads OpenAPI spec
- Extracts all paths, methods, parameters
- Builds internal test plan

**Step 2: Test Generation**
- For each endpoint, generates test cases:
  - Happy path (valid request)
  - Boundary cases (min/max values)
  - Invalid cases (missing required fields)
  - Type mismatches (string instead of int)

**Example**:
```yaml
# Spec
parameters:
  - name: store_nbr
    type: integer
    minimum: 10
    maximum: 999999

# Generated tests
Test 1: store_nbr = 3188 (valid)
Test 2: store_nbr = 10 (boundary min)
Test 3: store_nbr = 999999 (boundary max)
Test 4: store_nbr = 9 (invalid, below min)
Test 5: store_nbr = 1000000 (invalid, above max)
Test 6: store_nbr = "abc" (invalid, wrong type)
```

**Step 3: Request Execution**
- Sends HTTP requests to actual API
- Captures responses
- Records timing, status codes

**Step 4: Validation**
- Compares response to spec schema:
  - Status code matches?
  - Response body structure matches?
  - Required fields present?
  - Data types correct?
  - Enum values valid?

**Step 5: Coverage Calculation**
```
Coverage = (Operations Tested / Total Operations) × 100
```

**Step 6: Reporting**
- Generates HTML/JSON report
- Lists passed/failed operations
- Shows schema mismatches
- Calculates coverage percentage

**Threshold Check**:
```
if (coverage < threshold) {
    fail_build()
    rollback_deployment()
}
```

**Why 80% Threshold?**
- 100% is hard to achieve (edge cases, error scenarios)
- 80% covers all happy paths + major error cases
- Balances quality vs effort

**Comparison to Pact**:
- Pact: Consumer-driven (client writes contracts)
- R2C: Provider-driven (server writes contracts)
- Pact: Stub-based (mock server)
- R2C: Live API testing (real server)

We chose R2C because:
- Server team owns contract (design-first)
- Tests against real implementation
- Catches deployment-time issues"

---

### KEY TAKEAWAYS FOR INTERVIEW

**Technical Highlights**:
- OpenAPI 3.0 specification
- Design-first approach (spec → code)
- openapi-generator-maven-plugin
- Contract testing (R2C, 80% coverage)
- Server stub generation
- Modular spec structure (multiple YAML files)
- Semantic versioning

**Scale**:
- 20+ endpoints revamped
- 6 microservices
- 8 teams adopted
- 80% contract test coverage achieved

**Business Impact**:
- 30% reduction in client integration time (2 weeks → 1.4 weeks)
- Breaking changes prevented (compile-time safety)
- Standardized APIs across teams
- Improved documentation (Swagger UI, README.io)
- Faster onboarding (clear contracts)

**Architecture Patterns**:
- Contract-first development
- Interface-driven design
- Modular spec composition
- Semantic versioning
- Automated code generation

**Interview Story Arc**:
1. **Problem**: Code-first approach, outdated docs, integration issues
2. **Solution**: Design-first with OpenAPI 3.0
3. **Implementation**: Spec → Generation → Implementation → Contract Testing
4. **Challenges**: Verbose specs, generated code quality, breaking changes
5. **Results**: 20+ endpoints, 30% faster integration, 80% test coverage

---

This completes Part A (current 5 bullets).

**Total Words**: ~35,000 words (bullets 1-5)
**Estimated Interview Prep Time**: 10-15 hours to master all content

Ready for Part B (new 9-12 bullets)? These will cover:
- DC Inventory Search (your contribution)
- Store Inventory APIs
- Transaction History APIs
- Multi-tenant architecture
- And more...

Let me know if you want me to continue with Part B!

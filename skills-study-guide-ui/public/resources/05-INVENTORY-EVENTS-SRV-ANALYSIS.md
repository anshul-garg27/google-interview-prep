# ULTRA-THOROUGH ANALYSIS: inventory-events-srv

## Executive Summary

The **Inventory Events API** (inventory-events-srv) is a supplier-facing REST API service that provides transaction history for GTINs (Global Trade Item Numbers) at Walmart stores. It delivers near-real-time inventory event data including item sales, returns, receiving events, store transfers, loss prevention, physical inventory, point of fulfillment, and backroom operations. The service supports multi-tenant architecture across US, Mexico (MX), and Canada (CA) markets with sophisticated supplier authentication, GTIN-level authorization, and comprehensive observability.

**Team**: Data Ventures - Channel Performance Engineering (Luminate-CPerf-Dev-Group)
**Product ID**: 5798 (DV-DATAAPI)
**APM ID**: APM0017245

---

## 1. SERVICE PURPOSE AND FUNCTIONALITY

### Core Capabilities

- **Transaction History**: Retrieves inventory events for specific GTINs at stores
- **Multi-Tenant Support**: US, Mexico, and Canada markets
- **Supplier Authentication**: Consumer ID-based authentication via Walmart IAM
- **GTIN Authorization**: Ensures suppliers only access authorized products
- **Event Filtering**: By type (sales, returns, receiving, etc.) and location (STORE, BACKROOM, MFC)
- **Date Range Queries**: Historical data with configurable date windows
- **Pagination Support**: Handles large result sets with continuation tokens
- **Sandbox Environments**: Testing endpoints with static data

### Business Context

**Purpose**: Enable suppliers to track product movements and inventory transactions across Walmart's retail operations in real-time

**Key Use Cases**:
- Monitor product sales velocity
- Track customer returns
- Verify receiving events
- Analyze store transfers
- Identify shrinkage (loss prevention)
- Validate physical inventory counts
- Monitor online order fulfillment

---

## 2. TECHNOLOGY STACK

### Core Framework
- **Java 17** (LTS)
- **Spring Boot 3.5.6**
- **Spring Framework 6.2.x**
- **Maven 3.9+**

### Database & Persistence
- **PostgreSQL** (Primary database)
- **Hibernate 6.6.5** (ORM)
- **Hibernate Search 6.2.4** (Full-text search)
- **JPA 3.1** (Jakarta Persistence API)
- **HikariCP** (Connection pooling)

### Walmart Platform (Strati)
```
├── strati-af-runtime-context 12.0.2
├── strati-af-logging-log4j2-impl
├── strati-af-metrics-impl
├── strati-af-txn-marking-* (Distributed tracing)
├── ccm2-utils-client-spring 3.0.8
└── txn-marking-bom
```

### API & Documentation
- **OpenAPI 3.0.3**
- **Springdoc OpenAPI 2.3.0** (Swagger UI)
- **OpenAPI Generator Maven Plugin 7.0.1**

### Testing
- **JUnit 5** (Jupiter 5.10.3)
- **Mockito 5.10.0**
- **Cucumber 7.19.0** (BDD)
- **TestNG 7.9.0**
- **Rest Assured 5.5.0**
- **WireMock**

### Observability
- **Micrometer 1.11.4**
- **Prometheus** (Metrics export)
- **OpenTelemetry** (Distributed tracing)
- **Dynatrace SaaS** (APM - Production)
- **Log4j2** (Structured logging)

### Infrastructure
- **Docker** (Containerization)
- **Kubernetes (WCNP)** - Walmart Cloud Native Platform
- **Istio** (Service mesh)
- **Flagger** (Canary deployments)
- **KITT** (Walmart's CI/CD)

---

## 3. ARCHITECTURE AND CODE ORGANIZATION

### Project Structure

```
inventory-events-srv/
├── api-spec/                      # OpenAPI specifications
│   ├── schema/                    # API schema definitions
│   ├── examples/                  # Response examples
│   └── common-types/              # Reusable types
├── src/main/java/com/walmart/inventory/
│   ├── InventoryApplication.java  # Spring Boot entry
│   ├── controller/                # REST controllers
│   │   ├── InventoryEventsController.java
│   │   ├── InventoryEventsSandboxController.java
│   │   └── InventoryEventsSandboxIntlController.java
│   ├── services/                  # Business logic
│   │   ├── impl/
│   │   │   ├── InventoryStoreServiceImpl.java
│   │   │   ├── SupplierMappingServiceImpl.java
│   │   │   ├── StoreGtinValidatorServiceImpl.java
│   │   │   └── HttpServiceImpl.java
│   │   ├── validator/             # Business validators
│   │   ├── builders/              # Request builders
│   │   └── helpers/               # Service helpers
│   ├── repository/                # Data access layer
│   │   ├── ParentCmpnyMappingRepository.java
│   │   └── NrtiMultiSiteGtinStoreMappingRepository.java
│   ├── entity/                    # JPA entities
│   │   ├── ParentCompanyMapping.java
│   │   └── NrtiMultiSiteGtinStoreMapping.java
│   ├── models/                    # DTOs and domain models
│   ├── enums/                     # Enumerations
│   ├── common/                    # Common utilities
│   │   ├── config/               # Configuration classes
│   │   └── constants/            # Application constants
│   ├── exceptions/                # Exception handling
│   ├── filter/                    # HTTP filters
│   ├── context/                   # Request context management
│   ├── factory/                   # Factory patterns
│   ├── transaction/               # Transaction management
│   └── wrappers/                  # Security wrappers
├── ccm/                           # CCM2 configurations
├── setup_docs/                    # Documentation
│   └── telemetry/                # Monitoring setup
├── perf/                          # Performance test data
├── pom.xml                        # Maven configuration
├── kitt.yml                       # KITT build orchestration
├── kitt.us.yml                    # US deployment
├── kitt.intl.yml                  # International deployment
└── sr.yaml                        # Service Registry config
```

### Architectural Patterns

#### 1. Layered Architecture
```
Controller → Service → Repository → Database
```

#### 2. Multi-Tenant Architecture
- Site ID-based partitioning
- Thread-local `SiteContext`
- Hibernate `@PartitionKey` for data segregation
- CCM2 configuration per site (US/MX/CA)

#### 3. API-First Development
- OpenAPI 3.0 specification drives development
- Server-side code generation
- Contract-first approach

#### 4. Repository Pattern
- Spring Data JPA repositories
- Custom query methods
- Caching layer with `@Cacheable`

---

## 4. KEY COMPONENTS

### Controller Layer

#### InventoryEventsController

**Path**: `/v1/inventory/events`
**Method**: GET

**Responsibilities**:
- Accepts supplier requests with GTIN and store number
- Validates request parameters
- Extracts consumer ID from headers
- Orchestrates supplier mapping and GTIN validation
- Fetches inventory events from downstream services
- Returns paginated responses

**Key Dependencies**:
- `SupplierMappingServiceImpl`
- `InventoryStoreServiceImpl`
- `InventoryBusinessValidatorService`
- `TransactionMarkingManager`

### Service Layer

#### InventoryStoreServiceImpl

**Primary Business Logic Service**

**Responsibilities**:
1. Orchestrates inventory event retrieval workflow
2. Validates GTIN-supplier mappings
3. Calls Enterprise Inventory (EI) API
4. Handles pagination
5. Implements data integrity checks
6. Builds response DTOs

**Key Methods**:
- `getStoreInventoryData()` - Main entry point
- `getInventoryEvents()` - Fetches from EI API
- `checkDataIntegrityAndHandleIteration()` - Pagination logic

**Integration Points**:
- **EI Inventory History Lookup API**:
  - US: `ei-inventory-history-lookup.walmart.com`
  - MX: `ei-inventory-history-lookup-mx.walmart.com`
  - CA: `ei-inventory-history-lookup-ca.walmart.com`

#### SupplierMappingServiceImpl

**Supplier Authentication & Authorization**

**Responsibilities**:
1. Validates consumer ID from request headers
2. Fetches supplier metadata from `nrt_consumers` table
3. Handles PSP (Payment Service Provider) suppliers
4. Validates DUNS numbers
5. Multi-persona support

**Caching**: Uses Spring Cache with `parentCompanyMappingCacheManager`

#### StoreGtinValidatorServiceImpl

**GTIN Authorization Service**

**Responsibilities**:
1. Validates supplier has access to requested GTIN
2. Checks store-GTIN-supplier associations
3. Prevents unauthorized data access

**Data Source**: `supplier_gtin_items` table

### Repository Layer

#### ParentCmpnyMappingRepository
- **Entity**: `ParentCompanyMapping`
- **Table**: `nrt_consumers`
- **Primary Key**: Composite (consumerId, siteId)
- **Caching**: Enabled

#### NrtiMultiSiteGtinStoreMappingRepository
- **Entity**: `NrtiMultiSiteGtinStoreMapping`
- **Table**: `supplier_gtin_items`
- **Primary Key**: Composite (globalDuns, gtin, siteId)
- **Store Numbers**: PostgreSQL array type

### Configuration Components

#### Site-Specific CCM Configs
- **USEiApiCCMConfig**: US Enterprise Inventory endpoint
- **MXEiApiCCMConfig**: Mexico EI endpoint
- **CAEiApiCCMConfig**: Canada EI endpoint

#### PostgresDbConfiguration
- DataSource from Akeyless secrets
- Connection string: `/etc/secrets/db_url.txt`
- Credentials: `/etc/secrets/db_username.txt`, `db_password.txt`

### Filter Layer

#### SiteContextFilter
- Extracts `WM-Site-Id` header
- Populates `SiteContext` thread-local
- Enables multi-tenant routing

#### XssFilter
- XSS attack prevention
- Request sanitization
- Security wrapper application

---

## 5. API ENDPOINTS

### Primary Endpoint

```http
GET /v1/inventory/events
```

#### Request Parameters

| Parameter | Type | Required | Description | Validation |
|-----------|------|----------|-------------|------------|
| `store_nbr` | Integer | Yes | Store number | Min: 10, Max: 999999 |
| `gtin` | String | Yes | 14-digit GTIN | Length: 14, Numeric only |
| `start_date` | Date | No | Start date | YYYY-MM-DD |
| `end_date` | Date | No | End date | YYYY-MM-DD |
| `event_type` | String | No | Event filter | Default: "ALL", Max: 50 chars |
| `location_area` | String | No | Location filter | STORE, BACKROOM, MFC |
| `page_token` | String | No | Pagination token | Max: 255 chars |

**Default Date Range**: Last 6 days

#### Request Headers

| Header | Required | Description |
|--------|----------|-------------|
| `wm_consumer.id` | Yes | Supplier consumer ID (UUID) |
| `Authorization` | Yes | OAuth Bearer token |
| `WM-QOS.CORRELATION-ID` | No | Correlation ID |
| `WM-Site-Id` | Yes | Market identifier (US/MX/CA) |

#### Response Schema

```json
{
  "supplier_name": "string",
  "store_nbr": 3067,
  "gtin": "07502223774001",
  "next_page_token": "string",
  "items": [
    {
      "event_id": "string",
      "event_type": "NGR",
      "event_date_time": "2025-03-15T10:30:00Z",
      "quantity": 10,
      "unit_of_measure": "EACH",
      "location_area": "STORE",
      "transaction_details": {}
    }
  ]
}
```

#### Event Types

| Code | Description |
|------|-------------|
| `NGR` | Goods Receipt |
| `LP` | Loss Prevention |
| `POF` | Point of Fulfillment |
| `LR` | Returns |
| `PI` | Physical Inventory |
| `BR` | Backroom Operations |
| `ALL` | All event types |

#### Response Codes

| Status | Meaning |
|--------|---------|
| 200 | Success |
| 400 | Bad Request |
| 401 | Unauthorized |
| 404 | Not Found (GTIN not mapped) |
| 500 | Internal Error |

### Sandbox Endpoints

- **US Sandbox**: `/v1/inventory/sandbox/events`
- **International Sandbox**: `/v1/inventory/sandbox-intl/events`

Static response data for supplier testing

### Health & Monitoring Endpoints

```
GET /actuator/health           # Overall health
GET /actuator/health/liveness  # Kubernetes liveness
GET /actuator/health/readiness # Kubernetes readiness
GET /actuator/health/startup   # Kubernetes startup
GET /actuator/prometheus       # Metrics export
GET /swagger-ui/index.html     # API documentation
```

---

## 6. DATA MODELS AND SCHEMAS

### Database Entities

#### ParentCompanyMapping (nrt_consumers)

**Purpose**: Supplier metadata and authentication

```java
@Entity
@Table(name = "nrt_consumers")
public class ParentCompanyMapping {
  @EmbeddedId
  private ParentCompanyMappingKey primaryKey;

  private String consumerId;         // UUID
  private String consumerName;       // Supplier name
  private String countryCode;        // US/MX/CA
  private String globalDuns;         // DUNS number
  private String parentCmpnyName;    // Parent company
  private String luminateCmpnyId;    // Luminate ID
  private Boolean isCategoryManager;
  private Boolean nonCharterSupplier;
  private String pspGlobalDuns;      // PSP DUNS
  private Status status;             // ACTIVE/INACTIVE
  private UserType userType;         // Supplier type

  @PartitionKey
  private String siteId;             // Partition key
}
```

#### NrtiMultiSiteGtinStoreMapping (supplier_gtin_items)

**Purpose**: GTIN-to-supplier authorization mapping

```java
@Entity
@Table(name = "supplier_gtin_items")
public class NrtiMultiSiteGtinStoreMapping {
  @EmbeddedId
  private NrtiMultiSiteGtinStoreMappingKey primaryKey;

  @PartitionKey
  private String globalDuns;         // Partition key

  private Integer[] storeNumber;     // PostgreSQL array
  private String gtin;               // 14-digit GTIN
  private String luminateCmpnyId;
  private String parentCompanyName;

  @PartitionKey
  private String siteId;
}
```

### Enumerations

```java
public enum InventoryEventType {
  NGR,  // Goods Receipt
  LP,   // Loss Prevention
  POF,  // Point of Fulfillment
  LR,   // Returns
  PI,   // Physical Inventory
  BR    // Backroom
}

public enum Status {
  ACTIVE,
  INACTIVE
}

public enum SupplierPersona {
  PSP,      // Payment Service Provider
  VOLT,     // Volt persona
  STANDARD  // Standard supplier
}

public enum ValidLocationArea {
  STORE,
  BACKROOM,
  MFC  // Micro-Fulfillment Center
}
```

---

## 7. DEPENDENCIES AND INTEGRATIONS

### External Service Dependencies

#### 1. Enterprise Inventory (EI) API

**Purpose**: Transaction history data source

**Endpoints** (per market):
- **US**: `https://ei-inventory-history-lookup.walmart.com/v1/historyForInventoryState/countryCode/{us}/nodeId/{nodeId}/gtin/{gtin}`
- **MX**: `https://ei-inventory-history-lookup-mx.walmart.com/...`
- **CA**: `https://ei-inventory-history-lookup-ca.walmart.com/...`

**Authentication**: WM_CONSUMER.ID header
**Pagination**: Uses `continuationPageToken` header

#### 2. CCM2 (Tunr) Configuration Service

**Purpose**: Runtime configuration management

**Endpoints**:
- Non-Prod: `http://tunr.non-prod.walmart.com/scm-app/v2`
- Production: `http://tunr.prod.walmart.com/scm-app/v2`

**Configs Managed**: EI API endpoints, CORS settings, Site ID mappings

#### 3. Akeyless Secrets Management

**Path**: `/Prod/WCNP/homeoffice/Luminate-CPerf-Dev-Group/cp-inventory-events-apis`

**Secrets**:
- `db_username.txt`
- `db_password.txt`
- `db_url.txt`

#### 4. Observability Platform (Wolly)

**Purpose**: Distributed tracing and log aggregation

**Trace Collector**:
- Non-Prod: `trace-collector.nonprod.walmart.com:80`
- Production: `trace-collector.prod.walmart.com:80`

#### 5. Prometheus/Grafana (MMS)

**Metrics Export**: `/actuator/prometheus`

**Dashboards**:
- Golden Signals: `https://grafana.mms.walmart.net/d/QH3Ldm8YZ`
- Application Health
- Istio Service

### Database Integration

#### PostgreSQL Connection

```yaml
Datasource:
  URL: jdbc:postgresql://<host>:<port>/<database>
  Driver: org.postgresql.Driver
  Hibernate Dialect: PostgreSQL
  Connection Pool: HikariCP
```

**Schema**: Multi-tenant with partition keys
- `nrt_consumers` - Supplier mappings
- `supplier_gtin_items` - GTIN authorizations

### Service Mesh (Istio)

#### Configuration

```yaml
Istio Integration:
  Sidecar Injection: Enabled
  Protocol: HTTP/1.1
  Application Port: 8080
  Policy Engine: Enabled

Rate Limiting (Local):
  Max Tokens: 75
  Tokens Per Fill: 75
  Fill Interval: 1 second
```

### Authentication & Authorization

#### OAuth Token Validation
- **Stage**: `https://idp.stg.sso.platform.prod.walmart.com`
- **Production**: `https://idp.prod.global.sso.platform.prod.walmart.com`
- **Header**: `Authorization: Bearer <token>`

#### API Gateway (Service Registry)
- **Application Keys**: `INVENTORY-ACTIVITY`, `INVENTORY-ACTIVITY-INTL`
- **Rate Limiting**: 900 TPM per consumer ID

---

## 8. DEPLOYMENT CONFIGURATION

### Resource Allocation

#### Development
```yaml
min: { cpu: 1 core, memory: 1Gi }
max: { cpu: 1 core, memory: 1Gi }
scaling: { min: 1, max: 1, cpuPercent: 60% }
```

#### Stage
```yaml
min: { cpu: 1 core, memory: 1Gi }
max: { cpu: 2 cores, memory: 2Gi }
scaling: { min: 1, max: 2, cpuPercent: 60% }
```

#### Production
```yaml
min: { cpu: 1 core, memory: 1Gi }
max: { cpu: 2 cores, memory: 2Gi }
scaling: { min: 4, max: 8, cpuPercent: 60% }
```

### Health Probes

**Startup Probe**:
```yaml
path: /actuator/health/startup
port: 8080
initialDelaySeconds: 30
periodSeconds: 5
failureThreshold: 24
```

**Liveness Probe**:
```yaml
path: /actuator/health/liveness
port: 8080
periodSeconds: 10
failureThreshold: 5
```

**Readiness Probe**:
```yaml
path: /actuator/health/readiness
port: 8080
periodSeconds: 5
failureThreshold: 5
```

### Deployment Stages

#### US Deployment (kitt.us.yml)

1. **dev-us** - `eus2-dev-a2`
2. **stage-us** - `eus2-stage-a4`, `scus-stage-a3`
3. **sandbox-us** - `uswest-stage-az-006`
4. **prod (US)** - `eus2-prod-a30`, `scus-prod-a63`

#### International Deployment (kitt.intl.yml)

1. **dev-intl** - `scus-dev-a3`
2. **stage-intl** - `scus-stage-a6`, `useast-stage-az-303`
3. **sandbox-intl** - `uswest-stage-az-002`
4. **prod-intl** - `scus-prod-a16`, `useast-prod-az-321`

### Environment Variables

```bash
Common Variables:
  -Dccm.configs.dir=/etc/config
  -Druntime.context.appName=INVENTORY-ACTIVITY
  -Dcom.walmart.platform.metrics.impl.type=MICROMETER
  -Dcom.walmart.platform.logging.profile=SLT_DUAL
  -Dcom.walmart.platform.telemetry.otel.enabled=true
```

---

## 9. CI/CD PIPELINE

### Build Orchestration (KITT)

```yaml
build:
  mvnGoals: clean install

  docker:
    app:
      buildArgs:
        dtSaasOneAgentEnabled: true

  postBuild:
    - Deploy to US (kitt.us.yml)
    - Deploy to International (kitt.intl.yml)
```

### Maven Build Pipeline

```bash
Build Steps:
1. Code Formatting Check (Checkstyle)
2. OpenAPI Spec Consolidation
3. Server-Side Code Generation
4. Compilation (Java 17)
5. Unit Test Execution (JUnit 5)
6. Code Coverage (Jacoco)
7. Static Analysis (SonarQube)
8. Package JAR
9. Docker Image Build
10. Image Push to Registry
```

### Post-Deploy Steps

#### Stage Environment

1. **R2C Contract Testing**
   - Threshold: 80%
   - Mode: Passive

2. **Regression Test Suite**
   - RestAssured Concord Job

3. **Automaton Performance Tests**
   - Load testing

#### Production Deployment

1. **Manual Approval Required**
2. **Change Record Creation**
3. **Rollback Configuration**

### Deployment Strategy

#### Canary Deployment (Flagger)

```yaml
flagger:
  enabled: true
  canaryReplicaPercentage: 20%
  canaryAnalysis:
    interval: 2m
    stepWeight: 5
    maxWeight: 20
  metrics:
    - name: Check for 5XX errors
      threshold: 1%
```

**Process**:
1. New version deployed as canary (20% traffic)
2. Metrics analyzed every 2 minutes
3. Traffic gradually increased
4. If 5XX rate > 1%, automatic rollback
5. If successful, promote to primary

---

## 10. MONITORING AND OBSERVABILITY

### Distributed Tracing

#### Walmart Observability Platform

**Libraries**:
- `strati-af-txn-marking-springboot3-server 6.4.3`

**Configuration**:
```yaml
com.walmart.platform.telemetry.otel.enabled: true
com.walmart.platform.txnmarking.otel.type: OTLPgRPC
com.walmart.platform.txnmarking.otel.host: trace-collector.nonprod.walmart.com
com.walmart.platform.logging.profile: SLT_DUAL
com.walmart.platform.logging.file.format: OTelJson
```

**Transaction Markers**:
```
PS - Process Start
PE - Process End
RS - Request Start
RE - Request End
CE - Call End
CS - Call Start
```

**Trace Lookup**:
- **Non-Prod**: `https://wolly.non-prod.walmart.com`
- **Production**: `https://wolly.walmart.com`

### Metrics Collection (Prometheus)

#### Exported Metrics

**JVM Metrics**:
- `jvm_threads_live_threads`
- `jvm_memory_used_bytes`
- `jvm_gc_pause_seconds_*`

**HTTP Metrics**:
- `http_server_requests_seconds_count`
- `http_server_requests_seconds_sum`
- `http_client_requests_seconds_*`

**Database Metrics**:
- `hikaricp_connections_active`
- `hikaricp_connections_idle`
- `hikaricp_connections_pending`

**Custom Metrics**:
- `processTime_seconds`
- `c4XX_total` (4xx error count)
- `c5XX_total` (5xx error count)

**Endpoint**: `/actuator/prometheus`

### Grafana Dashboards

1. **Golden Signals Dashboard**: Latency, Traffic, Errors, Saturation
2. **Application Health Dashboard**: CPU, Memory, Pod status
3. **Istio Service Dashboard**: Service mesh traffic

### Health Checks

**Liveness Probe**:
- `livenessState` - App process running
- `ping` - Basic connectivity

**Readiness Probe**:
- `readinessState` - Ready for traffic
- `db` - Database connectivity

**Startup Probe**:
- `livenessState`
- `readinessState`
- `db`
- `diskSpace`
- `ssl`

---

## 11. SECURITY CONFIGURATIONS

### Authentication

#### OAuth 2.0 Token Validation
- **Stage**: `https://idp.stg.sso.platform.prod.walmart.com`
- **Production**: `https://idp.prod.global.sso.platform.prod.walmart.com`
- **Header**: `Authorization: Bearer <token>`

#### Consumer ID Validation
- **Header**: `wm_consumer.id`
- **Format**: UUID
- **Validation**: Against `nrt_consumers` table
- **Status Check**: ACTIVE suppliers only

### Authorization

#### GTIN Access Control
- **Table**: `supplier_gtin_items`
- **Check**: GTIN-to-supplier mapping
- **Enforcement**: Pre-request validation
- **Result**: 404 if GTIN not mapped

#### Store Number Authorization
- **Field**: `store_nbr` array
- **Check**: Store in supplier's authorized list

### Input Validation

```java
Validations:
  store_nbr: Integer, Range 10-999999, Required
  gtin: String, Length 14, Numeric only, Required
  start_date: ISO 8601, Not in future
  end_date: ISO 8601, >= start_date
  event_type: NGR|LP|POF|LR|PI|BR|ALL
  location_area: STORE|BACKROOM|MFC
```

### XSS Protection

**Filter**: `XssFilter`
**Wrapper**: `XssRequestWrapper`
**Strategy**: Input sanitization + Output encoding

### Rate Limiting

#### Pod-Level (Istio Local)
```yaml
maxTokens: 75
tokensPerFill: 75
fillIntervalSeconds: 1
Result: 75 requests/second per pod
```

#### API Gateway Level
- **Type**: TPM (Transactions Per Minute)
- **Criteria**: wm_consumer.id header
- **Limit**: 900 TPM per consumer

### Secrets Management (Akeyless)

```yaml
Provider: Akeyless
Path: /Prod/WCNP/homeoffice/.../cp-inventory-events-apis
Secrets:
  - db_username
  - db_password
  - db_url
Mount Point: /etc/secrets/
```

---

## 12. SPECIAL PATTERNS AND NOTEWORTHY IMPLEMENTATIONS

### 1. Multi-Tenant Architecture with Site Context

```java
@Component
public class SiteContext {
  private static final ThreadLocal<Long> siteIdThreadLocal = new ThreadLocal<>();

  public void setSiteId(Long siteId) {
    siteIdThreadLocal.set(siteId);
  }
}
```

**Benefits**: Automatic tenant isolation, no explicit tenant parameters

### 2. Factory Pattern for Site-Specific Configuration

```java
@Component
public class SiteConfigFactory {
  public SiteConfig getConfigurations(Long siteId) {
    return configMap.get(String.valueOf(siteId));
  }
}
```

### 3. OpenAPI-Driven Code Generation

**Workflow**: Define API in OpenAPI → Maven generates code → Controllers implement interfaces

### 4. Pagination with Continuation Tokens

```java
// Recursive call if more data exists
if (hasDataIntegrityIssue && StringUtils.isNotBlank(token)) {
  requestDto.setPageToken(token);
  getInventoryEvents(requestDto, supplierName);
}
```

### 5. Supplier Persona Handling (PSP Support)

**PSP** = Payment Service Provider (third-party payment processors)

```java
if (SupplierPersona.PSP.equals(parentCompanyMapping.getPersona())) {
  globalDuns = parentCompanyMapping.getPspGlobalDuns();
}
```

### 6. Caching Strategy

```java
@Cacheable(
  value = "PARENT_COMPANY_MAPPING_CACHE",
  key = "#consumerId + '-' + #siteId",
  unless = "#result == null"
)
```

### 7. Distributed Tracing with Nested Spans

```java
try (var parentTxn = txnManager.currentTransaction()
    .addChildTransaction("EI_SERVICE_CALL", "GET_INVENTORY_DATA")
    .start()) {
  // Business logic
}
```

### 8. PostgreSQL Array Support

```java
@Column(name = "store_nbr", columnDefinition = "integer[]")
private Integer[] storeNumber;
```

### 9. Hibernate Partition Keys (Multi-Tenancy)

```java
@PartitionKey
@Column(name = "site_id")
private String siteId;
```

### 10. Sandbox Mode with Static Data

**Controllers**: `InventoryEventsSandboxController`, `InventoryEventsSandboxIntlController`
**Purpose**: Supplier testing without production data

---

## SUMMARY

The **Inventory Events Service** is a sophisticated, enterprise-grade microservice demonstrating modern software engineering practices:

**Key Strengths**:
1. **Scalable Architecture**: Multi-tenant design (US/MX/CA)
2. **API-First Development**: OpenAPI-driven with code generation
3. **Production-Ready**: Comprehensive observability and monitoring
4. **Security-Focused**: OAuth, DUNS authorization, rate limiting
5. **Cloud-Native**: Kubernetes, Istio, Flagger canary deployments
6. **Walmart Platform Integration**: Strati, CCM2, Akeyless, Observability

**Key Metrics**:
- 4-8 pods in production (auto-scaling)
- 900 TPM rate limit per supplier
- 2-second timeout for EI calls
- 80% contract test pass rate
- 1% max 5XX error rate (canary threshold)

**Technology Highlights**:
- Java 17, Spring Boot 3.5.6
- PostgreSQL with JPA/Hibernate
- OpenAPI 3.0 with codegen
- OpenTelemetry distributed tracing
- Prometheus + Grafana
- Istio service mesh with mTLS
- Flagger progressive canary deployments
- KITT CI/CD with multi-stage pipelines

This service exemplifies Walmart's commitment to supplier enablement through robust, scalable, and secure data APIs.

---

**Document Version**: 1.0
**Last Updated**: 2026-02-02
**Analyzed By**: Claude Code Ultra Analysis
**Project**: Walmart Data Ventures Luminate
**Team**: Luminate-CPerf-Dev-Group

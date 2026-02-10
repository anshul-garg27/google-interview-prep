# ULTRA-THOROUGH ANALYSIS: inventory-status-srv

## Executive Summary

**inventory-status-srv** is a query/read service that provides near real-time visibility of inventory status for suppliers at Walmart stores and distribution centers. Unlike inventory-events-srv (which processes inventory change events), this service focuses on querying current inventory state through REST APIs. It supports multi-tenant architecture across US, Canada, and Mexico markets with sophisticated supplier authentication, GTIN-level authorization, and bulk query capabilities (up to 100 items per request).

**Team**: Data Ventures - Channel Performance Engineering
**Product ID**: 5798 (DV-DATAAPI)
**APM ID**: APM0017245

---

## 1. SERVICE PURPOSE AND FUNCTIONALITY

### Key Difference from inventory-events-srv

| Aspect | inventory-status-srv | inventory-events-srv |
|--------|---------------------|---------------------|
| **Purpose** | Query current inventory status | Process inventory change events |
| **Pattern** | Request-Response (Synchronous) | Event-Driven (Asynchronous) |
| **Data Flow** | Reads from EI/Database → Returns | Consumes Kafka → Processes → Publishes |
| **Primary Tech** | Spring Web, PostgreSQL, REST | Kafka Streams, Event sourcing |
| **Users** | External suppliers (querying) | Internal systems (processing) |
| **Operations** | GET-like queries (via POST) | Event processing, aggregation |

### Core Business Capabilities

1. **Store Inventory Search** (`/v1/inventory/search-items`)
   - Query current inventory status at Walmart stores
   - Supports GTIN or WM Item Number lookups
   - Up to 100 items per request
   - Near real-time visibility for suppliers

2. **Distribution Center Inventory Search** (`/v1/inventory/search-distribution-center-status`)
   - Query inventory status at distribution centers
   - WM Item Number based lookups
   - Up to 100 items per request
   - Inventory by type and state

3. **Inbound Inventory Tracking** (`/v1/inventory/search-inbound-items-status`)
   - Track items in transit from DC to stores
   - Store inbound items in transit
   - Expected arrival dates
   - Location and state tracking

### Business Value

- **Supplier Empowerment**: Real-time inventory visibility
- **Operational Efficiency**: Batch queries reduce API calls
- **Data Quality**: Validation ensures accurate lookups
- **Multi-Market Support**: US, Canada, Mexico

---

## 2. TECHNOLOGY STACK

### Core Framework
- **Java 17** (LTS)
- **Spring Boot 3.5.6** (Jakarta EE)
- **Spring Framework 6.2.x**
- **Maven** (Build tool)

### Database & Persistence
- **PostgreSQL** (walmart-postgresql driver)
- **Spring Data JPA** with Hibernate 6.6.5
- **Repository pattern**
- **Multi-tenant with partition keys**

### Walmart Platform Integration (Strati)
```
├── strati-af-runtime-context (Runtime environment)
├── strati-af-logging-log4j2-impl (Structured logging)
├── strati-af-metrics-impl (Metrics collection)
├── strati-af-txn-marking-* (Distributed tracing)
└── ccm2-utils-client-spring (Configuration management)
```

### API & Documentation
- **OpenAPI 3.0.3** specification
- **SpringDoc OpenAPI 2.3.0** (Swagger UI)
- **OpenAPI Generator Maven Plugin** (Code generation)

### Testing
- **JUnit 5** (Jupiter)
- **Mockito** (Mocking)
- **TestContainers** (Integration tests)
- **Cucumber** (BDD via Test Genie)
- **R2C** (Contract testing)
- **JMeter** (Performance testing)

### Observability
- **Micrometer** (Metrics facade)
- **Prometheus** (Metrics export)
- **Grafana** (Dashboards)
- **OpenTelemetry** (Distributed tracing)
- **Log4j2** (Structured logging)

### Infrastructure
- **Docker** (Containerization)
- **Kubernetes (WCNP)** - Walmart Cloud Native Platform
- **Istio** (Service mesh)
- **KITT** (CI/CD)

---

## 3. ARCHITECTURE AND CODE ORGANIZATION

### Project Structure

```
com.walmart.inventory/
├── InventoryApplication.java          # Spring Boot entry point
│
├── common/
│   ├── config/                         # Configuration classes
│   │   ├── ApplicationConfig.java      # CSV processing on startup
│   │   ├── EiApiCCMConfig.java         # EI API configuration
│   │   └── postgres/                   # PostgreSQL configuration
│   └── constants/                      # Constants and messages
│
├── controller/                         # REST API controllers
│   ├── InventoryItemsController.java                      # Store inventory
│   ├── InventorySearchDistributionCenterController.java   # DC inventory
│   ├── *SandboxController.java                            # Sandbox endpoints
│   └── *SandboxIntlController.java                        # International sandbox
│
├── services/                           # Business logic services
│   ├── impl/                           # Service implementations
│   │   ├── InventoryStoreServiceImpl.java
│   │   ├── InventorySearchDistributionCenterServiceImpl.java
│   │   ├── InventoryStoreInboundItemsServiceImpl.java
│   │   ├── UberKeyReadService.java
│   │   ├── RequestProcessor.java
│   │   └── ...
│   ├── helpers/                        # Helper utilities
│   ├── builders/                       # Request builders
│   └── validator/                      # Business validators
│
├── repository/                         # Data access layer
│   ├── ParentCmpnyMappingRepository.java
│   └── NrtiMultiSiteGtinStoreMappingRepository.java
│
├── entity/                             # JPA entities
│   ├── ParentCompanyMapping.java
│   └── NrtiMultiSiteGtinStoreMapping.java
│
├── models/                             # Domain models
│   ├── ei/                             # EI service models
│   └── association/                    # Association models
│
├── filter/                             # Request/response filters
│   ├── EndpointAccessValidationFilter.java  # Feature flag control
│   ├── SiteContextFilter.java              # Multi-site support
│   ├── RequestFilter.java                  # Request logging
│   └── XssFilter.java                      # XSS protection
│
├── context/                            # Thread-local context
│   ├── SiteContext.java                # Site-specific context
│   └── SiteTaskDecorator.java          # Thread pool decorator
│
├── factory/                            # Configuration factories
│   ├── SiteConfigFactory.java          # Site-specific configs
│   └── impl/                           # US, CA, MX configs
│
├── enums/                              # Enumerations
├── exceptions/                         # Exception handling
├── transaction/                        # Transaction marking
├── utils/                              # Utility classes
└── wrappers/                           # Request wrappers
```

### Architectural Patterns

#### 1. Multi-Site Architecture
- Supports US, CA, MX markets
- Site context propagation through filters
- Site-specific configurations via factory pattern
- CCM-driven per-market configuration

#### 2. Multi-Tenancy
- Site-based partitioning (`@PartitionKey` on site_id)
- Site context in thread-local storage
- Site-aware database queries
- Composite keys include site_id

#### 3. Service Layer Pattern
```
Controller → Service → Repository → Database
          → External API (EI, UberKey)
```

#### 4. API-First Development
- OpenAPI spec defines contract
- Code generation for models and API interfaces
- Controller delegates implement generated interfaces
- Contract testing validates implementation

#### 5. Bulk Processing Pattern
- RequestProcessor for batch validation
- CompletableFuture for parallel processing
- List partitioning for batch operations
- Error collection without stopping

#### 6. Error Handling Strategy
- Partial success responses (207 Multi-Status pattern)
- Error collection during processing
- Combined success/error responses
- Always return HTTP 200

---

## 4. KEY COMPONENTS

### Controllers

#### InventoryItemsController

**Endpoints**:
- `POST /v1/inventory/search-items` - Store inventory query
- `POST /v1/inventory/search-inbound-items-status` - Inbound tracking

**Responsibilities**:
- Validates request parameters
- Extracts consumer ID and site ID from headers
- Delegates to InventoryStoreService
- Returns paginated responses

#### InventorySearchDistributionCenterController

**Endpoint**: `POST /v1/inventory/search-distribution-center-status`

**Responsibilities**:
- DC inventory queries
- WM Item Number only (no GTIN support)
- Delegates to InventorySearchDistributionCenterService

#### Sandbox Controllers

**Purpose**: Testing endpoints with mock data
- InventorySandboxController (US market)
- InventorySandboxIntlController (International)
- Static responses from CSV files loaded at startup

### Core Services

#### InventoryStoreServiceImpl

**Primary Business Logic Service**

**Responsibilities**:
1. Store inventory retrieval
2. GTIN validation via StoreGtinValidatorService
3. UberKey integration for GTIN↔WM Item Number mapping
4. EI service calls for inventory data
5. Parallel processing with CompletableFuture

**Key Methods**:
- `getStoreInventoryData()` - Main entry point
- `getAllInventoryData()` - Batch processing
- `getInventoryForGtin()` - Single GTIN lookup

#### InventorySearchDistributionCenterServiceImpl

**DC Inventory Service**

**Processing Stages**:
1. **Stage 1**: WmItemNbr → GTIN conversion (UberKey)
2. **Stage 2**: Supplier validation
3. **Stage 3**: EI data fetch
4. Comprehensive error handling at each stage

#### InventoryStoreInboundItemsServiceImpl

**Inbound Inventory Tracking**

**Responsibilities**:
- Track items in transit from DC to store
- Expected arrival dates
- Location and state tracking

#### RequestProcessor

**Generic Bulk Validation Processor**

**Features**:
- Reusable across different request types
- Error collection and aggregation
- Validation without stopping processing

#### UberKeyReadService

**GTIN ↔ WM Item Number Mapping**

**Key Methods**:
- `getWmItemNbr()` - GTIN to WM Item Number
- `getCidDetails()` - Get CID (Consumer Item ID)
- Batch processing support
- CompletableFuture for parallel calls

### Validation & Business Logic

#### InventoryBusinessValidatorService

**Responsibilities**:
- Request parameter validation
- Business rule enforcement
- Error message preparation
- Max item count validation (100)

#### StoreGtinValidatorService

**GTIN Authorization**

**Responsibilities**:
- GTIN-supplier-store mapping validation
- Database queries for authorization
- Prevents unauthorized data access

#### SupplierMappingService

**Supplier Authentication**

**Responsibilities**:
- Consumer ID to supplier mapping
- Database lookups
- Global DUNS validation

### Filters (Security & Processing)

#### 1. EndpointAccessValidationFilter (Order 0)
- Feature flag-based endpoint control
- CCM-driven configuration
- Fail-open strategy (allow on CCM error)
- Endpoint enable/disable without deployment

#### 2. SiteContextFilter
- Multi-site support
- Extracts WM-Site-Id header
- Populates SiteContext thread-local

#### 3. RequestFilter
- Request logging
- Header extraction
- Correlation ID tracking

#### 4. XssFilter
- XSS attack prevention
- Input sanitization
- Request wrapping

#### 5. InventoryCorsFilter
- CORS policy enforcement
- Allowed origins configuration

---

## 5. API ENDPOINTS

### Production Endpoints

#### 1. Store Inventory Search

```http
POST /v1/inventory/search-items
```

**Request Body**:
```json
{
  "item_type": "gtin",           // "gtin" or "wm_item_nbr"
  "store_nbr": 3188,             // 10-999999
  "item_type_values": [          // 1-100 items
    "00012345678901",
    "00012345678902"
  ]
}
```

**Response**:
```json
{
  "items": [
    {
      "gtin": "00012345678901",
      "wm_item_nbr": 123456789,
      "dataRetrievalStatus": "SUCCESS",
      "store_nbr": 3188,
      "inventories": [
        {
          "location_area": "STORE",
          "state": "ON_HAND",
          "quantity": 150
        }
      ]
    }
  ]
}
```

**Features**:
- Supports both GTIN and WM Item Number
- Up to 100 items per request
- Multi-status responses (partial success)
- Cross-reference details

---

#### 2. Distribution Center Inventory

```http
POST /v1/inventory/search-distribution-center-status
```

**Request Body**:
```json
{
  "distribution_center_nbr": 6012,
  "wm_item_nbrs": [123456789, 987654321]  // 1-100 items
}
```

**Response**:
```json
{
  "items": [
    {
      "wm_item_nbr": 123456789,
      "gtin": "00012345678901",
      "dataRetrievalStatus": "SUCCESS",
      "dc_nbr": 6012,
      "inventories": [
        {
          "inventory_type": "AVAILABLE",
          "quantity": 5000
        }
      ]
    }
  ]
}
```

**Constraints**:
- WM Item Number only (no GTIN)
- DC number required

---

#### 3. Inbound Inventory Status

```http
POST /v1/inventory/search-inbound-items-status
```

**Request Body**:
```json
{
  "store_nbr": 3188,
  "gtins": ["00012345678901"]
}
```

**Response**:
```json
{
  "items": [
    {
      "gtin": "00012345678901",
      "store_nbr": 3188,
      "dataRetrievalStatus": "SUCCESS",
      "inbound_items": [
        {
          "expected_arrival_date": "2026-02-05",
          "quantity": 100,
          "location_area": "STORE",
          "state": "IN_TRANSIT"
        }
      ]
    }
  ]
}
```

---

### Sandbox Endpoints

Same paths with mock/test data:
- US Sandbox: Separate controller
- International Sandbox: Separate controller
- Data loaded from CSV files at startup

### Health & Monitoring

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
  private ParentCompanyMappingKey primaryKey; // (consumerId, siteId)

  private String consumerId;         // UUID
  private String consumerName;
  private String countryCode;        // US/MX/CA
  private String globalDuns;         // DUNS number
  private String parentCmpnyName;
  private String luminateCmpnyId;
  private Status status;             // ACTIVE/INACTIVE
  private UserType userType;

  @PartitionKey
  private String siteId;             // Multi-tenant partition
}
```

#### NrtiMultiSiteGtinStoreMapping (supplier_gtin_items)

**Purpose**: GTIN-to-supplier authorization mapping

```java
@Entity
@Table(name = "supplier_gtin_items")
public class NrtiMultiSiteGtinStoreMapping {
  @EmbeddedId
  private NrtiMultiSiteGtinStoreMappingKey primaryKey; // (gtin, globalDuns, siteId)

  @PartitionKey
  private String globalDuns;

  private Integer[] storeNumber;     // PostgreSQL array
  private String gtin;               // 14-digit GTIN
  private String luminateCmpnyId;
  private String parentCompanyName;

  @PartitionKey
  private String siteId;
}
```

### External Service Models

#### EI (Enterprise Inventory) Models
- `SingleGtinInventoryResponse` - Store inventory
- `EIDCInventoryResponse` - DC inventory
- `EIInboundInventoryResponse` - Inbound inventory
- Nested structures for location, state, quantity

#### API Models (Auto-generated from OpenAPI)
- `InventorySearchItemsRequest/Response`
- `InventorySearchDistributionCenterStatusRequest/Response`
- `InventorySearchInboundItemsStatusRequest/Response`
- `InventoryItemDetails`
- `CrossReferenceDetails`
- `InventoryLocationDetails`

### Processing Models
- `RequestProcessingResult` - Generic result container
- `RequestValidationErrors` - Error details
- `SiteConfigMapper` - Site-specific configuration

---

## 7. DEPENDENCIES AND INTEGRATIONS

### Internal Walmart Services

#### 1. Enterprise Inventory (EI) Service

**Purpose**: Source of truth for inventory data

**Integration**: REST API (HTTP GET)

**Endpoints**:
- Store inventory
- DC inventory
- Inbound inventory

**Data**: Real-time inventory by location, state

#### 2. UberKey Service

**Purpose**: GTIN ↔ WM Item Number mapping

**Integration**: REST API

**Features**:
- Item number translation
- CID retrieval
- Batch support (multiple GTINs per request)

#### 3. Service Registry

**Purpose**: API catalog and discovery

**Integration**: sr.yaml configuration

**Features**:
- Consumer management
- API versioning
- Policy enforcement

### External Dependencies

#### PostgreSQL Database
- Multi-site GTIN-store mappings
- Supplier-consumer mappings
- Partition key: site_id
- Connection pooling via HikariCP

#### CCM2 (Configuration Management)
- Centralized configuration
- Hot-reload support via @ManagedConfiguration
- Configs: EI endpoints, feature flags, database

#### Akeyless (Secrets Management)
- API keys for external services
- Database credentials
- Integration via Concord workflows

### Platform Services

1. **Kitt** - Kubernetes deployment
2. **MMS** - Prometheus/Grafana monitoring
3. **Looper** - CI/CD pipeline
4. **Testburst** - Test result reporting
5. **Concord** - Workflow orchestration

---

## 8. DEPLOYMENT CONFIGURATION

### KITT Deployment

**Configuration Files**:
- `kitt.yml` - Build orchestration
- `kitt.us.yml` - US market deployment
- `kitt.intl.yml` - International deployment

**Multi-Market Support**:
- US Market: kitt.us.yml
- International: kitt.intl.yml (CA, MX)
- Separate deployments, shared codebase

**Build Configuration**:
```yaml
setup:
  releaseRefs: ["main", "stage", "develop"]

profiles:
  - git://dsi-dataventures-luminate:inventory-status-srv:main:kitt-common
  - git://dsi-dataventures-luminate:dv-build-assets:main:snyk-pr-check

build:
  attributes:
    mvnGoals: "clean install -DconcordProcessId=... -DbuildUrl=..."
  docker:
    app:
      buildArgs:
        dtSaasOneAgentEnabled: "true"  # Dynatrace agent injection

  postBuild:
    - deployApp to kitt.intl.yml
    - deployApp to kitt.us.yml
```

### CCM Configuration

**Files**:
- **NON-PROD-1.0-ccm.yml** - Dev/Stage environments
- **PROD-1.0-ccm.yml** - Production

**Configuration Types**:
- Database connection strings
- EI API endpoints (US/CA/MX)
- Feature flags (endpoint enable/disable)
- CORS settings
- Logging levels

### Service Registry

**File**: `sr.yaml`

**Key Configuration**:
- Application key: "INVENTORY-STATUS"
- Business criticality: MAJOR
- SOA Integration: Enabled
- Environments: stg, prod
- Consumer management with IDs
- Policies: URITransformPolicy, RateLimitPolicy

### Managed Namespace

**APM ID**: APM0017245
**Owner Group**: Luminate-CPerf-Dev-Group

---

## 9. CI/CD PIPELINE

### Build Pipeline

#### 1. Pre-Build
- **Snyk security scanning** (PR comments)
- **Pre-commit hooks** (.pre-commit-config.yaml)
- **Checkmarx scans** (.cxignore)

#### 2. Build
- Maven goals: `clean install`
- JUnit 5 tests
- Jacoco code coverage
- OpenAPI code generation
- Checkstyle validation (Google style)
- Shield enforcer rules

#### 3. Post-Build
- Docker image build
- Dynatrace OneAgent injection
- Parallel deployment to US and INTL

#### 4. Quality Gates
- **SonarQube analysis**
- **Code coverage thresholds**
- **Code smell detection**
- **Vulnerability scanning**

### Testing Pipeline

1. **Unit Tests** - JUnit 5 + Mockito
2. **Integration Tests** - TestContainers + PostgreSQL
3. **Contract Tests** - R2C (Provider contract testing)
4. **Functional Tests** - Test Genie (Cucumber/BDD)
5. **Performance Tests** - JMeter (perf/ directory)

### Continuous Deployment

- **Non-Prod**: Automatic on merge to develop/stage
- **Production**: Manual approval required
- **Canary Support**: Configured in setup_docs/

### Concord Workflows

**File**: `concord.yml`

**Workflows**:
1. API spec publishing to README.io
2. Multi-market spec updates
3. Akeyless secret retrieval
4. Looper job triggers

---

## 10. MONITORING AND OBSERVABILITY

### Metrics (Prometheus)

**Enabled Metrics**:
- JVM metrics (memory, threads, GC)
- HTTP metrics (request/response)
- Custom business metrics
- Database connection pool metrics

**Endpoint**: `/actuator/prometheus`

### Health Checks

**Probes**:
- **Liveness**: `/actuator/health/liveness` - Process alive
- **Readiness**: `/actuator/health/readiness` - DB connectivity
- **Startup**: `/actuator/health/startup` - Full initialization

**Health Indicator Groups**:
```properties
management.endpoint.health.group.liveness.include=livenessState,ping
management.endpoint.health.group.readiness.include=readinessState,db
management.endpoint.health.group.startup.include=livenessState,readinessState,db,diskSpace,ssl
```

### Grafana Dashboards

1. **Golden Signals Dashboard**: Latency, traffic, errors, saturation
2. **Application Health Dashboard**: Service metrics, JVM stats
3. **Istio Service Dashboard**: Service mesh metrics
4. **Torbit Dashboard**: Origin VIP health

### Alerting (MMS)

**Custom alert rules** in `telemetry/rules/`:

| Alert | Trigger | Severity |
|-------|---------|----------|
| service-5xx-alerts.yaml | 5xx errors > 1/30s | Critical |
| service-4xx-alerts.yaml | 4xx error patterns | Warning |
| service-latency-alerts.yaml | Latency thresholds | Warning |
| kafka-consumer-lag-alerts.yaml | Consumer lag | Warning |
| cosmos-alerts.yaml | Cosmos DB issues | Critical |
| azure-sql-alerts.yaml | SQL performance | Critical |
| elastic-search-alerts.yaml | ES health | Warning |
| megha-cache-alerts.yaml | Cache issues | Warning |

**Alert Routing**:
- **Slack**: #data-ventures-cperf-dev-ops
- **Email**: ChannelPerformanceEn@wal-mart.com
- **XMatters**: di-wcnp-training group
- **Teams**: di-wcnp-training

### Transaction Marking (Strati)

**Features**:
- Distributed tracing via TransactionMarkingManager
- Child transactions for external calls
- Trace ID propagation in responses
- Headers: WM_QOS.CORRELATION_ID, WM_SVC.NAME

**Usage**:
```java
try (var txn = transactionMarkingManager.currentTransaction()
    .addChildTransaction("EI_SERVICE_CALL", "GET_INVENTORY")
    .start()) {
  // Business logic
}
```

### Logging

**Framework**: Strati AF logging (log4j2)

**Features**:
- Structured logging with context
- LoggingAspect for method entry/exit
- SLF4J facade
- Trace ID in all logs

---

## 11. SECURITY

### Authentication

- **API Key** (wm_consumer.id header)
- **Authorization** header
- **Consumer ID validation** via database

### Authorization

**Multi-Level Authorization**:
1. **Consumer Level**: Supplier → Consumer mapping (ParentCompanyMapping)
2. **GTIN Level**: GTIN authorization (NrtiMultiSiteGtinStoreMapping)
3. **Store Level**: Store-level permissions
4. **Global DUNS**: DUNS number validation

### Security Filters

#### 1. XssFilter
- **XSS attack prevention**
- Wraps HttpServletRequest
- Sanitizes input parameters

#### 2. EndpointAccessValidationFilter
- **Feature flag-based access control**
- Endpoint enable/disable via CCM
- Hot-reload without restart

#### 3. SiteContextFilter
- **Site-based access control**
- Multi-tenant isolation

#### 4. InventoryCorsFilter
- **CORS policy enforcement**
- Configurable allowed origins

### Security Scanning

- **Checkmarx (SAST)** - .cxignore
- **Snyk (SCA)** - PR vulnerability scanning
- **Sentinel Policy** enforcement
- **Shield enforcer** rules

### Secrets Management

- **Akeyless** for API keys, DB credentials
- **CCM** for configuration secrets
- **No hardcoded secrets**

### Input Validation

**Validation Strategy**:
- JSON schema validation
- Bean validation (@Valid)
- Business rule validation
- Max items: 100 per request
- Pattern matching for item types

---

## 12. SPECIAL PATTERNS AND FEATURES

### 1. Multi-Status Response Pattern

**Always return HTTP 200** with status per item:

```json
{
  "items": [
    {"gtin": "123", "dataRetrievalStatus": "SUCCESS", ...},
    {"gtin": "456", "dataRetrievalStatus": "ERROR", "reason": "..."}
  ]
}
```

**Benefits**:
- Partial success supported
- Batch operations don't fail entirely
- Clear error messages per item

### 2. Bulk Request Processing

**Implementation**:
- RequestProcessor with generic validation
- Parallel CompletableFuture for UberKey calls
- List partitioning for batch size limits
- Error collection without stopping processing

### 3. Site Context Propagation

```java
@Component
public class SiteTaskDecorator implements TaskDecorator {
  @Override
  public Runnable decorate(Runnable runnable) {
    Long siteId = siteContext.getSiteId();
    return () -> {
      siteContext.setSiteId(siteId);
      try {
        runnable.run();
      } finally {
        siteContext.clear();
      }
    };
  }
}
```

**Purpose**: Propagates site context to worker threads for CompletableFuture operations

### 4. Multi-Market Configuration

```java
@Component
public class SiteConfigFactory {
  private Map<String, SiteConfigMapper> configMap;

  public SiteConfigMapper getConfigurations(Long siteId) {
    return configMap.get(String.valueOf(siteId));
  }
}
```

**Factory returns** US/CA/MX specific configs (different EI endpoints per market)

### 5. Feature Flag Control

**EndpointAccessValidationFilter**:
- Endpoint-level enable/disable
- CCM-driven configuration
- Hot-reload without restart
- Fail-open strategy (allow on error)

### 6. API-First Development

**Workflow**:
1. Define API in OpenAPI spec
2. Maven generates server-side code
3. Controllers implement generated interfaces
4. R2C validates implementation

### 7. Sandbox Endpoints

**Implementation**:
- Separate controllers for testing
- Mock data from CSV files
- ApplicationConfig loads CSV on startup
- US and International variants

### 8. Comprehensive Error Handling

```java
// Error collection during processing
RequestExceptionCollector.collect(error);
// Errors added to response, not thrown
// Partial success supported
```

### 9. Strati AF Integration

**Features**:
- Transaction marking for observability
- Child transactions for downstream calls
- Metric collection
- Logging framework
- Runtime context

### 10. Database Partitioning

**Multi-tenant with partition keys**:
- `@PartitionKey` on site_id
- Composite keys with site_id
- Site-aware queries via SiteContext

### 11. Cross-Reference Support

**Features**:
- GTIN variants tracked
- Cross-reference details in response
- Traversed GTIN tracking to avoid duplicates

### 12. Parallel Processing with CompletableFuture

```java
List<CompletableFuture<UberKeyResponse>> futures = gtins.stream()
  .map(gtin -> CompletableFuture.supplyAsync(
    () -> uberKeyService.getWmItemNbr(gtin),
    taskExecutor
  ))
  .collect(Collectors.toList());

CompletableFuture.allOf(futures.toArray(new CompletableFuture[0])).join();
```

---

## SUMMARY

**inventory-status-srv** is a sophisticated query service providing real-time inventory visibility to Walmart suppliers. Key strengths include:

### Technical Highlights

1. **Scalable Architecture**: Multi-tenant (US/CA/MX), bulk processing (100 items)
2. **API-First**: OpenAPI-driven with code generation
3. **Production-Ready**: Comprehensive observability, monitoring, alerting
4. **Security-Focused**: Multi-level authorization, XSS protection, feature flags
5. **Cloud-Native**: Kubernetes, Istio, Dynatrace, multi-market deployment
6. **Walmart Platform**: Deep Strati integration, CCM2, Akeyless

### Key Capabilities

- **Store Inventory**: Query current stock by GTIN/WM Item Number
- **DC Inventory**: Distribution center inventory status
- **Inbound Tracking**: Items in transit with ETA
- **Bulk Operations**: Up to 100 items per request
- **Partial Success**: Multi-status response pattern
- **Sandbox Support**: Testing endpoints for suppliers

### Technology Stack

- Java 17, Spring Boot 3.5.6
- PostgreSQL with JPA/Hibernate
- OpenAPI 3.0 with codegen
- Prometheus + Grafana + Dynatrace
- OpenTelemetry distributed tracing
- KITT CI/CD with multi-stage pipelines

### Comparison to inventory-events-srv

| Feature | inventory-status-srv | inventory-events-srv |
|---------|---------------------|---------------------|
| **Purpose** | Query current state | Process change events |
| **Pattern** | Synchronous query | Asynchronous events |
| **Users** | External suppliers | Internal systems |
| **API** | REST endpoints | Kafka consumers |
| **Operations** | Read queries | Write/process events |

---

**Document Version**: 1.0
**Last Updated**: 2026-02-02
**Analyzed By**: Claude Code Ultra Analysis
**Project**: Walmart Data Ventures Luminate
**Team**: Luminate-CPerf-Dev-Group

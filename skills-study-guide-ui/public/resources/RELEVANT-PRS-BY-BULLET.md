# Relevant PRs - Deep Analysis for Interview Preparation

> **Purpose:** Only the technically interesting PRs mapped to each resume bullet
> **Filter Criteria:** Excluded dummy commits, CCM-only changes, onboarding PRs, GBI updates
> **Focus:** Architecture decisions, core implementations, bug fixes, production deployments

---

## Bullet 1 & 10: Kafka-Based Audit Logging System + Multi-Region Architecture
**Repo:** `audit-api-logs-gcs-sink`

### Core Implementation PRs (MUST KNOW)

| PR# | Title | Why It's Important | Interview Talking Point |
|-----|-------|-------------------|------------------------|
| **#1** | Updating readme file, kafka topic and gcs bucket name | **Initial Setup** - Architecture decisions for Kafka Connect sink | "I designed the initial GCS sink architecture, choosing Parquet format for efficient columnar storage and BigQuery compatibility" |
| **#3** | Deployment changes main | **Deployment Architecture** - K8s deployment, helm charts | "Set up the Kubernetes deployment with proper resource allocation and health checks" |
| **#9** | Prod config, renaming class, updating field | **Production Config** - Class design, field mapping | "Configured production-grade settings with proper error handling and field transformations" |
| **#12** | [CRQ:CHG3041101] GCS Sink changes for US records | **Multi-Region Deployment** - First US production deployment | "Led the US region deployment with CRQ process, ensuring zero data loss during migration" |
| **#28** | CHG3180819 - Updating flush count, max.poll.interval.ms and offset assignor | **Performance Tuning** - Critical Kafka consumer tuning | "Tuned Kafka consumer configs to handle 2M+ events/day - adjusted flush count and poll intervals" |
| **#29** | Adding autoscaling profile and related changes | **Autoscaling** - KEDA-based scaling on Kafka lag | "Implemented autoscaling based on Kafka consumer lag to handle traffic spikes automatically" |

### Production Debugging Story (GREAT FOR BEHAVIORAL)

| PR# | Title | What Happened |
|-----|-------|---------------|
| **#35** | Adding try catch block in filter | Bug fix - unhandled exception in SMT filter |
| **#38** | Adding consumer poll config | Consumer configuration issues |
| **#42** | Removing scaling on lag. Updating the heartbeat | Autoscaling causing instability |
| **#47** | Increasing the poll timeout | Consumer rebalance issues |
| **#52** | Setting error tolerance to none | Changed from tolerant to strict mode |
| **#57** | Overriding KAFKA_HEAP_OPTS | Memory issues - heap tuning |
| **#59** | Adding the filter and custom header converter | Header parsing issues |
| **#61** | Increasing heap max | Final memory fix |

> **Interview Story:** "We had a production incident where the Kafka Connect sink was failing silently. Through PRs #35-61, I systematically debugged: first identified exception handling gaps (#35), then consumer polling issues (#38, #47), then autoscaling instability (#42), and finally heap memory constraints (#57, #61). The debugging process taught me Kafka Connect internals deeply."

### Multi-Region & Site-Based Routing (BULLET 10)

| PR# | Title | Why It's Important |
|-----|-------|-------------------|
| **#30** | adding prod-eus2 env | East US 2 region setup |
| **#70** | [CHG3235845] : Upgrade connector version. Remove custom header. Include retry policy | **Key PR** - Retry policy for resilience |
| **#85** | GCS Sink Changes to support different site ids | **Site-Based Routing** - SMT filters for multi-region |

---

## Bullet 2: Common Library JAR (Spring Boot Starter)
**Repo:** `dv-api-common-libraries`

### Core Implementation PRs (MUST KNOW)

| PR# | Title | Why It's Important | Interview Talking Point |
|-----|-------|-------------------|------------------------|
| **#1** | Audit log service | **THE MAIN PR** - Complete audit logging library | "Built a reusable Spring Boot starter with HTTP interceptor, async Kafka producer, and configurable endpoint filtering" |
| **#3** | adding webclient in jar | **WebClient Integration** - Reactive HTTP client for downstream calls | "Added WebClient support for tracing outbound HTTP calls in the audit trail" |
| **#6** | Develop | **Core Functionality** - Additional features merged | "Enhanced the library with configuration options and filtering capabilities" |
| **#11** | enable endpoints with regex | **Regex Filtering** - Flexible endpoint matching | "Added regex-based endpoint filtering so teams could selectively audit specific APIs" |

### Architecture Components (Know These):
```
dv-api-common-libraries/
├── AuditFilter.java           # HTTP request/response interceptor
├── KafkaAuditProducer.java    # Async Kafka producer
├── AsyncConfig.java           # Thread pool configuration
├── AuditLogProperties.java    # Spring Boot configuration properties
└── WebClientConfig.java       # Outbound HTTP tracing
```

---

## Bullet 3: DSD (Direct Store Delivery) Notification System
**Repo:** `cp-nrti-apis`

### FOUNDATIONAL PRs (The Beginning - IMPORTANT!)

| PR# | Title | Why It's Important | Interview Talking Point |
|-----|-------|-------------------|------------------------|
| **#3** | DSIDVL-2827: DSD API Implementation | **THE ORIGINAL DSD API** - First implementation | "Built the initial DSD (Direct Store Delivery) API from scratch" |
| **#441** | Dsidvl 17188 dsd shipment endpoint v1 | **DSD Shipment V1** - First shipment capture endpoint | "Designed the DSD Shipment endpoint for capturing delivery data" |
| **#448** | DSD shipment sandbox endpoint stage to main | **Production Launch** | "Launched DSD Shipment to production" |
| **#451** | DSIDVL-17605 Dsd shipments sandbox request audit into cosmos db | **Audit Trail** - Cosmos DB integration | "Added audit logging for DSD shipments to CosmosDB" |

### SUMO Push Notifications (KEY FEATURE!)

| PR# | Title | Why It's Important | Interview Talking Point |
|-----|-------|-------------------|------------------------|
| **#651** | Dsidvl 26408 sumo integration with common api | **SUMO Integration** - Push notification setup | "Integrated SUMO push notification service for real-time store alerts" |
| **#683** | DSD-Shipment---Change-Text-of-Sumo-Notification | **Notification Content** - Message customization | "Customized SUMO notification content for store associates" |
| **#704** | DSIDVL-27529 FeatureFlag-IsSumoEnabled | **Feature Flag** - Controlled rollout | "Added feature flag for controlled SUMO rollout" |
| **#1085** | DSIDVL-35619 DSC - Sumo roles changes | **Roles Integration** - Store-level routing | "Implemented store-level notification routing based on roles" |
| **#1093** | Dsidvl 35620 Update sumo api store roles changes | **Production [CRQ]** | "Production deployment of SUMO roles" |

### Kafka Event Publishing (CORE!)

| PR# | Title | Why It's Important | Interview Talking Point |
|-----|-------|-------------------|------------------------|
| **#151** | NRTIN-41,NRTIN-96 Kafka topic configurations | **Kafka Setup** - Initial topic config | "Set up Kafka topics and configurations for event streaming" |
| **#275** | DSIDVL-10230 Kafka upgrade changes | **Kafka Upgrade** - Version upgrade | "Led Kafka client upgrade with zero downtime" |
| **#699** | Dsidvl 27528 Dsc kafka payload changes | **Kafka Payload** - Event schema | "Designed DSC Kafka event schema" |
| **#1098** | DSC, IAC kafka connection properties streamline | **Kafka Optimization** - Connection pooling | "Streamlined Kafka connection properties for better performance" |
| **#1122** | Dsidvl 36644 kafka configs stream line [CRQ] | **Production** | "Production deployment of Kafka optimizations" |

### Recent Feature Additions

| PR# | Title | Why It's Important | Interview Talking Point |
|-----|-------|-------------------|------------------------|
| **#1397** | DSIDVL-44428 adding driver id to DSC request and kafka payload | **Driver ID Feature** - New field in Kafka events | "Added driver identification to DSC events for delivery tracking integration" |
| **#1400** | DSIDVL-44428 DSC Driver ID Main | Production deployment of driver ID | "Deployed to production with proper validation and backward compatibility" |
| **#1405** | adding packNbrQty to kafka payload | **Pack Quantity Feature** - Enhanced event data | "Extended Kafka payload with pack quantity for inventory reconciliation" |
| **#1417** | adding packNbrQty to kafka payload (Main) | Production deployment | "Production release with Kafka schema evolution" |
| **#1129** | DSIDVL-36599-DSC-PalletQty and CaseQty change | **Pallet/Case Qty** - Inventory tracking | "Added pallet and case quantity fields for inventory accuracy" |
| **#1179** | Dsc pallet qty change totalling pallet quantity | **Pallet Totalling** - Aggregation logic | "Implemented pallet quantity totalling" |
| **#1433** | DSIDVL-47611 isDockRequired request field changes | **Dock Required Feature** - Store delivery requirement | "Added dock requirement flag for store delivery scheduling optimization" |

### Bug Fixes (Good Debugging Stories)

| PR# | Title | What Happened |
|-----|-------|---------------|
| **#682** | Validation and exceptions for DSC request scenario | Validation improvements |
| **#701** | bug fix for validation | DSC validation edge cases |
| **#1144** | Errors fixed for DSC and IAC regression suite | Regression test failures |
| **#1516** | trip id fix and dock required bug fix | Trip ID validation issue |
| **#1543** | Modified order for NumberFormatException | **5xx Error Fix** - Vendor ID parsing order issue |

### Deprecation Work

| PR# | Title | Why It's Important |
|-----|-------|-------------------|
| **#1389-1391** | DSIDVL-43873 removing inventory snapshot code | Clean deprecation of legacy feature |

---

## Bullet 4: Spring Boot 3 & Java 17 Migration
**Repo:** `cp-nrti-apis`

### Core Migration PRs (MUST KNOW)

| PR# | Title | Why It's Important | Interview Talking Point |
|-----|-------|-------------------|------------------------|
| **#1312** | Feature/springboot 3 migration | **THE MAIN PR** - Complete Spring Boot 3 migration | "Led the migration from Spring Boot 2.7 to 3.2, handling javax→jakarta namespace changes, Hibernate 6 compatibility, and security configuration updates" |
| **#1337** | Feature/springboot 3 migration (Release) | **Production Deployment** | "Deployed to production with zero downtime using Flagger canary deployments" |

### Post-Migration Fixes (Shows Thoroughness)

| PR# | Title | What It Fixed |
|-----|-------|--------------|
| **#1332** | snyk issues fixed and mutable list issue fixed | Security vulnerabilities + immutable collections |
| **#1334** | snyk issues and mutable list (merged) | Stage deployment |
| **#1335** | one major sonar issue fixed | Code quality issue |

### Key Migration Challenges:
1. **javax → jakarta** namespace migration
2. **Hibernate 6** - Query changes, deprecated methods
3. **Spring Security 6** - New configuration style
4. **WebFlux compatibility** - Reactor updates

---

## Bullet 5: OpenAPI Design-First Approach
**Repo:** `inventory-status-srv`

### Core Implementation PRs (MUST KNOW)

| PR# | Title | Why It's Important | Interview Talking Point |
|-----|-------|-------------------|------------------------|
| **#249** | DSIDVL-44382-api spec development (Store Inbound) | **Store Inbound API Spec** - OpenAPI contract | "Designed the Store Inbound API spec using design-first approach, ensuring R2C compatibility" |
| **#260** | Feature/dc inventory api spec | **DC Inventory API Spec** - New API contract | "Created DC Inventory API specification with comprehensive schema definitions" |
| **#219** | Contract Testing Integration | **R2C Setup** - Contract testing framework | "Integrated R2C contract testing into CI/CD pipeline for API validation" |
| **#233** | Contract Testing Integration (Merged) | Production deployment of contract tests | "Deployed contract testing to all environments" |
| **#262** | Feature/dsidvl 46057 api publishing automation | **API Publishing Automation** - Developer portal publishing | "Automated API spec publishing to Walmart's developer portal" |
| **#331** | Update Contract test gate for us | **Contract Test Gate** - CI gate for API changes | "Added contract test as a deployment gate to prevent breaking changes" |

### Supporting Infrastructure

| PR# | Title | Why It's Important |
|-----|-------|-------------------|
| **#259** | Container Tests | Integration testing setup |
| **#266** | Container Tests (Merged) | Production deployment |
| **#338** | Container test for dc inventory from main | DC Inventory specific tests |

---

## Bullet 6: DC Inventory Search API + Store Inventory Status
**Repos:** `inventory-status-srv`, `cp-nrti-apis`

### DC Inventory - Core Implementation PRs (MUST KNOW)

| PR# | Repo | Title | Why It's Important | Interview Talking Point |
|-----|------|-------|-------------------|------------------------|
| **#271** | inventory-status-srv | Feature/dc inventory development | **THE MAIN PR** - Complete DC Inventory implementation | "Built the DC Inventory Search API with CompletableFuture for parallel DC queries, reducing response time from 3s to 800ms for bulk requests" |
| **#312** | inventory-status-srv | fix dc inventory api spec and remove gtin in response | **API Spec Fix** - Response optimization | "Optimized response payload by removing redundant GTIN field" |
| **#322** | inventory-status-srv | V2 DSIDVL-48592 error handling dc inv | **Error Handling** - Comprehensive error responses | "Implemented standardized error handling with proper HTTP status codes and error messages" |
| **#330** | inventory-status-srv | Feature/dc inv error fix | Production error fixes | "Fixed production issues with proper exception handling" |

### Store Inventory Status - cp-nrti-apis (FOUNDATIONAL!)

| PR# | Title | Why It's Important | Interview Talking Point |
|-----|-------|-------------------|------------------------|
| **#513** | DSIDVL_20025 Dc Inventory Sandbox with Walmart Item Nbr | **DC Inventory Sandbox** - Initial implementation | "Built DC Inventory sandbox endpoint with WmItemNbr support" |
| **#535** | DSIDVL_19593 Dc inventory status wmItemNbr prod endpoint | **DC Inventory Production** | "Deployed DC Inventory Status to production" |
| **#554** | Dsidvl 21671 DC Inventory Status Error Response Enhancement | **Error Handling** | "Enhanced error responses for DC Inventory" |
| **#573** | DSIDVL-21992 Adding WmItemNbr to Store Inventory Response | **Store Inventory Enhancement** | "Added Walmart Item Number to Store Inventory response" |
| **#623** | Data Not Found Fix for Store & DC Inventory Status | **Bug Fix** | "Fixed data not found scenarios for both APIs" |

### Multi-GTIN & UberKey (KEY FEATURE!)

| PR# | Title | Why It's Important | Interview Talking Point |
|-----|-------|-------------------|------------------------|
| **#528** | DSIDVL-20624 Multi GTIN with WmItemNbr Enhancement - Sandbox | **Multi-GTIN** | "Implemented multi-GTIN query support for bulk inventory lookup" |
| **#531** | DSIDVL-20624 Store and DC Inventory Status Item identifier to prod | **Production** | "Production deployment of multi-GTIN feature" |
| **#797** | DSIDVL-29882 UberKey size limitation changes | **UberKey Optimization** | "Optimized UberKey handling for large inventory queries" |
| **#832** | DSIDVL_30280 Uber key size limit issue [CRQ] | **Production Fix** | "Fixed UberKey size limitation in production" |
| **#987** | DSIDVL-33578 Expand multi gtin mfc response | **MFC Response** | "Enhanced multi-GTIN response for MFC inventory" |

### Multi-Tenant Support

| PR# | Repo | Title | Why It's Important |
|-----|------|-------|-------------------|
| **#320** | inventory-status-srv | Adding CCM-ANY value and siteId as tenantId | **Multi-Tenant Architecture** - Site-based isolation |

### Performance & Testing

| PR# | Repo | Title | Why It's Important |
|-----|------|-------|-------------------|
| **#241** | inventory-status-srv | Load Testing | Performance validation |
| **#337** | inventory-status-srv | Added resiliency test | Chaos engineering |

---

## Bullet 7: Transaction Event History API
**Repo:** `inventory-events-srv`

### Core Implementation PRs (MUST KNOW)

| PR# | Title | Why It's Important | Interview Talking Point |
|-----|-------|-------------------|------------------------|
| **#5** | API spec development | **Initial Design** - OpenAPI specification | "Designed the Transaction Event History API with support for 10+ event types and cursor-based pagination" |
| **#8** | Dsidvl 40282 transaction history | **THE MAIN PR** - Complete implementation | "Built the Transaction History API with multi-tenant support, handling 5M+ transactions daily with efficient cursor-based pagination" |
| **#10** | DSIDVL-40848 Junit Sonar Snyk Codegate | **Testing & Quality** - Comprehensive test coverage | "Achieved 85%+ code coverage with unit and integration tests" |
| **#22** | Dsidvl 40282 transaction history (Merged) | Dev to Stage deployment | - |
| **#38** | CRQ CHG3095678 TransactionHistoryApi Changes Stage to Main | **Production Deployment** [CRQ] | "Deployed to production with CRQ process" |

### Contract Testing

| PR# | Title | Why It's Important |
|-----|-------|-------------------|
| **#25** | Feature/contract testing | R2C integration |

### Database Migration

| PR# | Title | Why It's Important |
|-----|-------|-------------------|
| **#80** | cosmos to postgres migration | **Database Migration** - CosmosDB to PostgreSQL |
| **#96** | Inventory Events CA Launch Cosmos-Postgres CCM Snyk [CRQ] | **Canada Launch** with all changes |

### IAC (Inventory Activity) Integration

| PR# | Title | Why It's Important |
|-----|-------|-------------------|
| **#123** | Feature/develop add iac endpoint contract ext | IAC V2 endpoint integration |
| **#125** | Feature/iac v2 to stage | Stage deployment |

---

## Bullet 8: Supplier Authorization Framework
**Repos:** `inventory-status-srv`, `cp-nrti-apis`

### Store Inbound Inventory - cp-nrti-apis (FOUNDATIONAL!)

| PR# | Title | Why It's Important | Interview Talking Point |
|-----|-------|-------------------|------------------------|
| **#374** | Added code for Store Inbound for Sandbox | **Initial Implementation** | "Built Store Inbound Inventory sandbox endpoint" |
| **#393** | [DSIDVL-14562] Added code for Store Inbound for Prod | **Production Code** | "Implemented Store Inbound for production" |
| **#477** | [DSIDVL-14562] Store Inbound Inventory code promotion [CRQ] | **Production Launch** | "Launched Store Inbound Inventory to production with CRQ" |
| **#505** | DSIDVL-XXXX Store Inbound Inventory End Date Enhancement | **End Date Feature** | "Added end date handling for store inbound inventory" |

### Core Implementation PRs - inventory-status-srv (MUST KNOW)

| PR# | Title | Why It's Important | Interview Talking Point |
|-----|-------|-------------------|------------------------|
| **#161** | inventory-status store-GTIN enhancements [CRQ:CHG3014193] | **Authorization Enhancements** - GTIN-level permissions | "Enhanced the 4-level authorization hierarchy (Consumer→DUNS→GTIN→Store) with new validation rules" |
| **#276** | Sandbox code for StoreInbound API | **Store Inbound Authorization** - New API with auth | "Implemented Store Inbound API with proper supplier authorization checks" |
| **#281** | Update error handling for store inbound | **Error Handling** - Authorization error messages | "Added clear authorization error messages for debugging" |

### Supplier Validation & GTIN Mapping - cp-nrti-apis

| PR# | Title | Why It's Important | Interview Talking Point |
|-----|-------|-------------------|------------------------|
| **#14** | DSIDVL-2828 supplier mapping service and cache | **Supplier Cache** | "Built supplier mapping service with Caffeine cache for fast lookups" |
| **#356** | DSIDVL-11658 Vendor GTIN Item Validation SANDBOX Endpoint | **GTIN Validation** | "Implemented GTIN-level supplier validation endpoint" |
| **#425** | DSIDVL-15852 Volt-IAC GTIN_NOT_MAPPED error message | **Error Messages** | "Added clear error messages for GTIN mapping failures" |
| **#460** | DSIDVL-18586 Enhanced store Grain validation | **Store Validation** | "Enhanced supplier validation at store grain level" |
| **#657** | PSP supplier context - NRT app | **PSP Integration** | "Integrated Partner Service Platform for supplier context" |
| **#722** | clean up for nrti_gtin_parent_company | **Data Cleanup** | "Cleaned up GTIN parent company data" |

### Database & Performance

| PR# | Repo | Title | Why It's Important |
|-----|------|-------|-------------------|
| **#197** | inventory-status-srv | Feature/cosmos to postgres migration | CosmosDB to PostgreSQL migration |
| **#216** | inventory-status-srv | Feature/cosmos to postgres migration (Merged) | Production deployment |

---

## Bullet 9: Observability Stack
**Repos:** `cp-nrti-apis`, `audit-api-logs-gcs-sink`, `inventory-status-srv`

### Dynatrace & APM (MUST KNOW)

| PR# | Repo | Title | Why It's Important | Interview Talking Point |
|-----|------|-------|-------------------|------------------------|
| **#105** | cp-nrti-apis | DSIDVL-3189: Enable Dynatrace monitoring | **Initial Dynatrace Setup** | "Set up Dynatrace monitoring for distributed tracing from the beginning" |
| **#106** | cp-nrti-apis | DSIDVL-3189: Enable Dynatrace in main | **Production Deployment** | "Deployed Dynatrace to production" |
| **#1545** | cp-nrti-apis | Added Dynatrace Saas | **Dynatrace SaaS Migration** | "Migrated from on-prem to Dynatrace SaaS for better scalability" |

### Correlation ID & Distributed Tracing

| PR# | Repo | Title | Why It's Important | Interview Talking Point |
|-----|------|-------|-------------------|------------------------|
| **#349** | cp-nrti-apis | Added changes for traceID | **TraceID Implementation** | "Implemented trace ID propagation across all services" |
| **#396** | cp-nrti-apis | correlation_id and trace id changes | **Correlation ID** | "Added correlation ID for end-to-end request tracking" |
| **#408** | cp-nrti-apis | traceId and ccm platform logging changes | **Logging Enhancement** | "Enhanced logging with trace context" |
| **#1527** | cp-nrti-apis | Correlation Id Issue Fix | **Bug Fix** | "Fixed correlation ID propagation edge cases" |

### Canary Deployments (Flagger)

| PR# | Repo | Title | Why It's Important | Interview Talking Point |
|-----|------|-------|-------------------|------------------------|
| **#1113** | cp-nrti-apis | DSIDVL-18013: Enabling canary deployment with Custom metrics | **Canary Setup** | "Implemented Flagger canary deployments with custom metrics gates" |
| **#1118** | cp-nrti-apis | Enabling canary deployment with Custom metrics (Stage) | Stage deployment | - |
| **#1208** | cp-nrti-apis | Enabling canary deployment with Custom metrics (Prod) | **Production Canary** | "Deployed canary deployment strategy to production" |

### Prometheus Metrics

| PR# | Repo | Title | Why It's Important | Interview Talking Point |
|-----|------|-------|-------------------|------------------------|
| **#222** | cp-nrti-apis | Enabling goldensignal metrics for NRT | **Golden Signals** - RED metrics | "Implemented Rate, Error, Duration (RED) metrics" |
| **#327** | inventory-status-srv | Added metrics in Kitt | **Custom Metrics** | "Added custom Prometheus metrics for business KPIs" |
| **#1271** | cp-nrti-apis | enable product metrics | **Product Metrics** | "Added product-specific metrics for monitoring" |

### Audit Logging Integration (Cross-Repo)

| PR# | Repo | Title | Component |
|-----|------|-------|-----------|
| **#155** | cp-nrti-apis | Nrtin 48 audit service changes | Initial audit service |
| **#1254** | cp-nrti-apis | integrate audit log service | Full integration |
| **#1269** | cp-nrti-apis | updated new version for audit logging | Version update |
| **#1289** | cp-nrti-apis | enable all endpoints for audit logging | Full coverage |
| **#1303** | cp-nrti-apis | audit log jar integration | JAR integration |
| **#1363** | cp-nrti-apis | Update ccm.yml | CCM configuration |
| **#1368** | cp-nrti-apis | update ccm for audit log endpoints | Endpoint configuration |
| **#1373** | cp-nrti-apis | Update audit log jar version | Version update |
| **#1489** | cp-nrti-apis | audit log config changes fix | Config fix |

---

## BONUS: IAC (Inventory Activity Capture) - Related Work
**Repo:** `cp-nrti-apis`

> IAC is the core inventory event capture system that powers many of the resume bullets

### Foundational PRs

| PR# | Title | Why It's Important | Interview Talking Point |
|-----|-------|-------------------|------------------------|
| **#142** | NRTIN-35: Inventory actions Sandbox and Prod endpoint | **Initial IAC Endpoint** | "Built the Inventory Activity Capture endpoint from scratch" |
| **#147** | NRTIN-77: Inventory actions capture Sandbox and Prod | **Full Implementation** | "Implemented full IAC flow with validation and Kafka publishing" |
| **#190** | NRTN-96 iac kafka header change | **Kafka Headers** | "Added custom headers to IAC Kafka events for routing" |
| **#295** | IAC: eventSubType in payload header | **Event Types** | "Implemented event subtypes for granular event categorization" |

### Kafka Idempotency & Performance

| PR# | Title | Why It's Important | Interview Talking Point |
|-----|-------|-------------------|------------------------|
| **#318** | DSIDVL-11154 Iac kafka idempotency config | **Idempotent Producer** | "Configured Kafka idempotent producer for exactly-once semantics" |
| **#333** | IAC Split && kafka Idempotent && Supplier mapping | **IAC Optimization** | "Split IAC processing and optimized Kafka publishing" |
| **#335** | IAC Split && kafka Idempotent stage to main [CRQ] | **Production** | "Production deployment of IAC optimizations" |

### Multi-Region Kafka

| PR# | Title | Why It's Important | Interview Talking Point |
|-----|-------|-------------------|------------------------|
| **#1045** | Iac-Kafka-cluster-changes-publish-in-both-regions [CRQ] | **Multi-Region Kafka** | "Implemented dual-region Kafka publishing for disaster recovery" |

### Supplier Validation

| PR# | Title | Why It's Important | Interview Talking Point |
|-----|-------|-------------------|------------------------|
| **#14** | DSIDVL-2828 supplier mapping service and cache | **Supplier Cache** | "Built supplier mapping service with Caffeine cache" |
| **#356** | DSIDVL-11658 Vendor GTIN Item Validation SANDBOX Endpoint | **GTIN Validation** | "Implemented GTIN-level supplier validation" |
| **#460** | Enhanced store Grain validation | **Store Validation** | "Enhanced supplier validation at store grain level" |
| **#657** | PSP supplier context - NRT app | **PSP Integration** | "Integrated Partner Service Platform for supplier context" |

---

## Quick Reference: Top 35 PRs for Interview (FINAL)

| Priority | Bullet | PR# | Repo | Title | Why |
|----------|--------|-----|------|-------|-----|
| 1 | 1 | #1 | audit-api-logs-gcs-sink | Initial setup | Architecture decisions |
| 2 | 2 | #1 | dv-api-common-libraries | Audit log service | Complete library |
| 3 | 4 | #1312 | cp-nrti-apis | Spring Boot 3 migration | Major migration |
| 4 | 7 | #8 | inventory-events-srv | Transaction history | Main implementation |
| 5 | 6 | #271 | inventory-status-srv | DC inventory development | Core feature |
| 6 | 3 | #3 | cp-nrti-apis | DSD API Implementation | **Original DSD!** |
| 7 | 3 | #651 | cp-nrti-apis | SUMO integration | **Push Notifications!** |
| 8 | 8 | #477 | cp-nrti-apis | Store Inbound Inventory [CRQ] | **Store Inbound Launch!** |
| 9 | 6 | #528 | cp-nrti-apis | Multi GTIN Enhancement | **Multi-GTIN Feature!** |
| 10 | 1 | #28 | audit-api-logs-gcs-sink | Performance tuning | Kafka expertise |
| 11 | 5 | #219 | inventory-status-srv | Contract Testing | R2C setup |
| 12 | 3 | #441 | cp-nrti-apis | DSD Shipment V1 | Original shipment capture |
| 13 | 9 | #105 | cp-nrti-apis | Enable Dynatrace | **Initial Observability** |
| 14 | IAC | #147 | cp-nrti-apis | IAC Sandbox and Prod | Core IAC implementation |
| 15 | 1 | #29 | audit-api-logs-gcs-sink | Autoscaling | K8s scaling |
| 16 | 10 | #85 | audit-api-logs-gcs-sink | Site-based routing | Multi-region |
| 17 | 8 | #161 | inventory-status-srv | GTIN enhancements | Authorization |
| 18 | 6 | #797 | cp-nrti-apis | UberKey size limitation | Performance fix |
| 19 | 9 | #1113 | cp-nrti-apis | Canary deployment | Flagger setup |
| 20 | 4 | #1337 | cp-nrti-apis | Spring Boot 3 (Prod) | Production deployment |
| 21 | 7 | #38 | inventory-events-srv | Transaction History (Prod) | CRQ deployment |
| 22 | 1 | #12 | audit-api-logs-gcs-sink | US records [CRQ] | Multi-region |
| 23 | 5 | #260 | inventory-status-srv | DC inventory API spec | Design-first |
| 24 | 6 | #322 | inventory-status-srv | Error handling DC inv | Production quality |
| 25 | 3 | #1397 | cp-nrti-apis | Driver ID to DSC | Kafka events |
| 26 | 3 | #1433 | cp-nrti-apis | isDockRequired | Feature addition |
| 27 | IAC | #318 | cp-nrti-apis | Kafka idempotency | Exactly-once semantics |
| 28 | 2 | #11 | dv-api-common-libraries | Regex endpoints | Flexibility |
| 29 | 7 | #80 | inventory-events-srv | Cosmos to Postgres | Database migration |
| 30 | 9 | #349 | cp-nrti-apis | TraceID implementation | Distributed tracing |
| 31 | IAC | #1045 | cp-nrti-apis | Multi-region Kafka | Dual-region publishing |
| 32 | 3 | #1085 | cp-nrti-apis | SUMO roles | Store-level routing |
| 33 | 8 | #14 | cp-nrti-apis | Supplier mapping cache | Caffeine cache |
| 34 | 8 | #356 | cp-nrti-apis | GTIN Item Validation | GTIN validation endpoint |
| 35 | 6 | #987 | cp-nrti-apis | Multi GTIN MFC response | MFC enhancement |

---

## Commands to View PR Details

```bash
# Set host for Walmart GitHub
export GH_HOST=gecgithub01.walmart.com

# View any PR
gh pr view <PR#> --repo dsi-dataventures-luminate/<repo-name>

# View PR diff (code changes)
gh pr diff <PR#> --repo dsi-dataventures-luminate/<repo-name>

# Examples for top PRs:
gh pr view 1 --repo dsi-dataventures-luminate/audit-api-logs-gcs-sink
gh pr view 1 --repo dsi-dataventures-luminate/dv-api-common-libraries
gh pr view 1312 --repo dsi-dataventures-luminate/cp-nrti-apis
gh pr view 8 --repo dsi-dataventures-luminate/inventory-events-srv
gh pr view 271 --repo dsi-dataventures-luminate/inventory-status-srv
```

---

## Interview Preparation Checklist

### Technical Deep Dive Questions:
- [ ] Explain Kafka Connect architecture and how sink connectors work
- [ ] How does CompletableFuture improve API performance?
- [ ] What are the key changes in Spring Boot 3 migration?
- [ ] How does R2C contract testing work?
- [ ] Explain multi-tenant authorization hierarchy

### Behavioral Stories (STAR Format):
- [ ] **Debugging Story:** PRs #35-61 in audit-api-logs-gcs-sink
- [ ] **Performance Story:** PR #28 - Kafka consumer tuning
- [ ] **Migration Story:** PR #1312 - Spring Boot 3 migration
- [ ] **Architecture Story:** PR #1 in dv-api-common-libraries

### System Design Knowledge:
- [ ] Draw audit logging architecture (Library → Kafka → GCS)
- [ ] Explain multi-region deployment strategy
- [ ] Describe authorization hierarchy (Consumer→DUNS→GTIN→Store)
- [ ] Explain cursor-based pagination implementation

---

*Generated for focused interview preparation - Only relevant PRs included*

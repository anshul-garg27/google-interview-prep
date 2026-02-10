# PR to Project Mapping - Interview Deep Dive Guide

> **Purpose:** Map all your merged PRs to resume bullets for interview preparation
> **Last Updated:** February 2026

---

## Table of Contents
1. [Repos Overview](#repos-overview)
2. [Bullet 1 & 10: Kafka Audit Logging & Multi-Region Architecture](#bullet-1--10-kafka-audit-logging--multi-region-architecture)
3. [Bullet 2: Common Library (JAR)](#bullet-2-common-library-jar)
4. [Bullet 3: DSD Notification System](#bullet-3-dsd-notification-system)
5. [Bullet 4: Spring Boot 3 & Java 17 Migration](#bullet-4-spring-boot-3--java-17-migration)
6. [Bullet 5: OpenAPI Design-First](#bullet-5-openapi-design-first)
7. [Bullet 6: DC Inventory Search API](#bullet-6-dc-inventory-search-api)
8. [Bullet 7: Transaction Event History API](#bullet-7-transaction-event-history-api)
9. [Bullet 8: Supplier Authorization Framework](#bullet-8-supplier-authorization-framework)
10. [Bullet 9: Observability Stack](#bullet-9-observability-stack)
11. [Quick Reference: Top PRs per Bullet](#quick-reference-top-prs-per-bullet)

---

## Repos Overview

| Repo | Total PRs | Primary Projects |
|------|-----------|------------------|
| `dsi-dataventures-luminate/inventory-status-srv` | 270+ | Bullet 6, 8, 9 (DC Inventory, Auth, Observability) |
| `dsi-dataventures-luminate/cp-nrti-apis` | 200+ | Bullet 3, 4, 5, 6, 7 (DSD, Migration, OpenAPI, Core APIs) |
| `dsi-dataventures-luminate/dv-api-common-libraries` | 15 | **Bullet 2** (Common Library JAR) |
| `dsi-dataventures-luminate/audit-api-logs-gcs-sink` | 87 | **Bullet 1, 10** (Kafka Audit, Multi-Region) |
| `dsi-dataventures-luminate/inventory-events-srv` | 143 | **Bullet 7** (Transaction Event History) |

---

## Bullet 1 & 10: Kafka Audit Logging & Multi-Region Architecture

### Resume Bullets
> **Bullet 1:** Engineered a high-throughput, Kafka-based audit logging system processing over 2 million events daily, adopted by 12+ teams to enable real-time API tracking.

> **Bullet 10 (Potential):** Engineered multi-region active/active Kafka architecture with Connect sink connectors streaming audit events to GCS with Parquet partitioning, featuring custom SMT filters and site-based routing across EUS2/SCUS clusters.

### Repository: `audit-api-logs-gcs-sink`

#### Key PRs (Interview-Worthy)

| PR# | Title | Date | Why It's Important |
|-----|-------|------|-------------------|
| **#1** | Updating readme file, kafka topic and gcs bucket name | 2025-01-07 | **Initial setup** - Architecture decisions, topic naming conventions |
| **#3** | Deployment changes main | 2025-01-22 | First production deployment |
| **#5** | Updating the repo name - audit-api-logs-gcs-sink | 2025-02-24 | Repo organization |
| **#9** | Prod config, renaming class, updating field | 2025-02-26 | Production readiness, class design |
| **#12** | GCS Sink changes for US records [CRQ:CHG3041101] | 2025-03-28 | **Multi-region US deployment** |
| **#28** | Updating flush count, max.poll.interval.ms [CHG3180819] | 2025-05-02 | **Performance tuning** - Talk about trade-offs |
| **#29** | Adding autoscaling profile and related changes | 2025-05-05 | **Scalability decisions** |
| **#35-59** | Hotfix series (config updates) | 2025-05 | **Troubleshooting story** - Debugging in prod |
| **#70** | Upgrade connector version, retry policy [CHG3235845] | 2025-06-02 | **Reliability improvements**, retry mechanisms |
| **#85** | GCS Sink Changes to support different site ids | 2025-11-19 | **Multi-tenant/site-based routing** |
| **#87** | Vulnerability update and Start up probe changes | 2026-01-27 | Security & health checks |

#### All PRs (Chronological)

```
#1   - Updating readme file, kafka topic and gcs bucket name (2025-01-07)
#3   - Deployment changes main (2025-01-22)
#4   - Updating the repo name as audit-api-logs-gcs-sink (2025-02-24)
#5   - Updating the repo name - audit-api-logs-gcs-sink (2025-02-24)
#9   - Prod config, renaming class, updating field (2025-02-26)
#10  - Storing in new test bucket (2025-03-27)
#11  - Updating the project id of new buckets (2025-03-28)
#12  - GCS Sink changes for US records [CRQ:CHG3041101] (2025-03-28)
#13  - Removing ephemeral storage (2025-04-01)
#14  - Updating the schema name folder (2025-04-02)
#15  - Updating the notification email (2025-04-02)
#18  - Dummy commit to reflect non prod ccm and sr changes (2025-04-02)
#19  - [CRQ:CHG3041101] Dummy commit (2025-04-02)
#20  - Updating the topic names (2025-04-03)
#21  - Updating the topic names (2025-04-03)
#22  - Updating KCQL of prod SCUS (2025-04-03)
#23  - Updating KCQL of prod SCUS (2025-04-03)
#24  - Dummy commit (2025-04-03)
#25  - CRQ : CHG3041101 Dummy commit (2025-04-03)
#26  - Adding insights.yml (2025-04-03)
#27  - Updating flush count, max.poll.interval.ms (2025-04-28)
#28  - CHG3180819 - Updating flush count, max.poll.interval.ms (2025-05-02)
#29  - Adding autoscaling profile and related changes (2025-05-05)
#30  - adding prod-eus2 env (2025-05-12)
#32  - Adding metadata as parent of labels in helm values (2025-05-12)
#34  - Increasing the min and max replicas (2025-05-12)
#35  - Adding try catch block in filter (2025-05-12)
#37  - Adding loggers (2025-05-13)
#38  - Adding consumer poll config (2025-05-13)
#39  - Changing to standard consumer configs (2025-05-13)
#40  - Updating CPU and ephemeral storage (2025-05-13)
#41  - Updating the heartbeat (2025-05-14)
#42  - Removing scaling on lag. Updating the heartbeat (2025-05-14)
#43  - Increasing the session time out (2025-05-14)
#46  - Decreasing the session time out (2025-05-14)
#47  - Increasing the poll timeout (2025-05-14)
#49  - Hotfix/update configs 1 (2025-05-14)
#50  - Reducing the tasks to 1 (2025-05-14)
#51  - Setting error tolerance to none (2025-05-14)
#52  - Setting error tolerance to none (2025-05-14)
#53  - update kc config (2025-05-14)
#55  - update latest version (2025-05-15)
#57  - Overriding KAFKA_HEAP_OPTS (2025-05-15)
#58  - Removing unit test changes and filter class (2025-05-15)
#59  - Adding the filter and custom header converter (2025-05-15)
#60  - Adding header converter at worker (2025-05-15)
#61  - Increasing heap max (2025-05-15)
#63  - Renaming stage as release/stage (2025-05-19)
#64  - intl env changes in non prod ccm (2025-05-22)
#65  - Adding gcs sink intl app key in sr.yaml (2025-05-26)
#66  - Adding names for config overrides (2025-05-27)
#68  - Adding tunr names in CCM (2025-05-27)
#69  - Upgrade connector version, retry policy (2025-05-30)
#70  - [CHG3235845] : Upgrade connector version, retry policy (2025-06-02)
#72  - Updating README.md with local set up (2025-06-09)
#73  - Adding all worker and connector property in CCM (2025-06-11)
#75  - Updating the data type as DECIMAL (2025-06-11)
#77  - Updating resolution path (2025-06-13)
#80  - Adding dummy commit (2025-06-17)
#81  - Updating the connector version (2025-06-18)
#83  - CCM Changes for country code to site id (2025-11-03)
#85  - GCS Sink Changes to support different site ids (2025-11-19)
#86  - Adding dummy commit (2026-01-14)
#87  - Vulnerability update and Start up probe changes (2026-01-27)
```

#### Interview Talking Points
- **System design**: Explain the pipeline architecture (API → Kafka → Connect → GCS)
- **Performance tuning**: The hotfix series (#35-59) shows real debugging - talk about flush.size, poll.interval, session.timeout trade-offs
- **Multi-region**: How you handled EUS2/SCUS active/active
- **Exactly-once vs at-least-once**: What semantics did you choose and why?
- **Schema evolution**: How Parquet format helps with schema changes

---

## Bullet 2: Common Library (JAR)

### Resume Bullet
> Delivered a JAR adopted by 12+ teams, enabling applications to monitor requests and responses through a simple POM dependency and configuration.

### Repository: `dv-api-common-libraries`

#### Key PRs (Interview-Worthy)

| PR# | Title | Date | Why It's Important |
|-----|-------|------|-------------------|
| **#1** | Audit log service | 2025-03-07 | **THE MAIN PR** - Initial library creation |
| **#3** | adding webclient in jar | 2025-03-20 | HTTP client integration |
| **#6** | Develop | 2025-03-20 | Core functionality merge |
| **#11** | enable endpoints with regex | 2025-04-08 | Flexibility feature - regex-based endpoint filtering |
| **#15** | Version update for gtp bom | 2025-09-19 | Java 17 compatibility |

#### All PRs (Chronological)

```
#1  - Audit log service (2025-03-07) *** MAIN PR ***
#3  - adding webclient in jar (2025-03-20)
#4  - snkyfix (2025-03-20)
#5  - Update pom.xml (2025-03-20)
#6  - Develop (2025-03-20)
#11 - enable endpoints with regex (2025-04-08)
#12 - Create CODEOWNERS (2025-04-11)
#15 - Version update for gtp bom (2025-09-19)
```

#### Interview Talking Points
- **Library design principles**: Spring Boot starter pattern, auto-configuration
- **HTTP Interceptor**: How you capture request/response automatically
- **Async Processing**: ThreadPoolExecutor for non-blocking audit
- **Why library vs sidecar?**: Trade-offs discussion
- **Sensitive data masking**: How did you handle PII?
- **Backward compatibility**: How did you ensure consumers don't break?
- **Thread pool sizing**: How did you decide on pool sizes?

---

## Bullet 3: DSD Notification System

### Resume Bullet
> Developed a notification system for Direct-Shipment-Delivery (DSD) suppliers, automating alerts to over 1,200 Walmart associates across 300+ store locations with a 35% improvement in stock replenishment timing.

### Repository: `cp-nrti-apis`

#### Key PRs (Interview-Worthy)

| PR# | Title | Date | Why It's Important |
|-----|-------|------|-------------------|
| **#1098** | DSC, IAC kafka connection properties streamline | 2024-11-11 | Kafka configuration |
| **#1129** | DSIDVL-36599-DSC-PalletQty and CaseQty change | 2024-12-02 | **Core DSC logic** |
| **#1145-1146** | DSC Pallet Qty change Totalling | 2024-12-12 | Business logic - totalling |
| **#1179-1181** | DSC Pallet Qty Totalling | 2025-01-10/13 | Pallet quantity calculations |
| **#1397** | adding driver id to DSC request and kafka payload | 2025-05-27 | **DSC Kafka integration** |
| **#1405-1406** | adding packNbrQty to kafka payload | 2025-06-25/26 | Payload enhancements |
| **#1433** | DSIDVL-47611 isDockRequired request field | 2025-08-22 | Feature additions |
| **#1516-1518** | trip id fix and dock required bug fix | 2025-09-25 | Bug fixes |

#### All DSC-Related PRs

```
#1098 - DSC, IAC kafka connection properties streamline (2024-11-11)
#1122 - Kafka configs streamline changes [CRQ: CHG2913306] (2024-11-26)
#1129 - DSIDVL-36599-DSC-PalletQty and CaseQty change (2024-12-02)
#1130 - DSC-PalletQty and CaseQty change (#1129) (2024-12-04)
#1140 - DSIDVL-36599-DSC-PalletQty stg-to-main (2024-12-10)
#1144 - Errors fixed for DSC and IAC regression suite (2024-12-11)
#1145 - DSC Pallet Qty change Totalling (2024-12-12)
#1146 - DSC Pallet Qty change dev to stg (2024-12-12)
#1179 - DSC Pallet Qty Totalling (2025-01-10)
#1181 - DSC Pallet Qty dev to stg (2025-01-13)
#1397 - adding driver id to DSC request and kafka payload (2025-05-27)
#1399 - DSC Driver ID Stage (2025-06-02)
#1400 - DSC Driver ID main (2025-06-03)
#1405 - adding packNbrQty to kafka payload (2025-06-25)
#1406 - packNbrQty to kafka payload (#1405) (2025-06-26)
#1417 - packNbrQty kafka payload main (2025-07-09)
#1433 - DSIDVL-47611 isDockRequired request field (2025-08-22)
#1446 - DSIDVL-47611 DockRequired request field changes (#1433) (2025-09-05)
#1514 - dock required changes to main (2025-09-24)
#1516 - trip id fix and dock required bug fix (2025-09-25)
#1517 - trip id fix stage (2025-09-25)
#1518 - trip id fix main (2025-09-25)
```

#### Interview Talking Points
- **Event-driven architecture**: Kafka event → Process → Sumo Push
- **Notification delivery guarantees**: How did you handle failures?
- **Rate limiting**: How did you prevent notification flooding?
- **Measuring 35% improvement**: How did you measure business impact?
- **Store-associate mapping**: How did you route to correct associates?

---

## Bullet 4: Spring Boot 3 & Java 17 Migration

### Resume Bullet
> Upgraded systems to Spring Boot 3 and Java 17, resolving critical vulnerabilities. Addressed challenges like backward compatibility, outdated libraries, and dependency management, ensuring zero downtime.

### Repository: `cp-nrti-apis`

#### Key PRs (Interview-Worthy)

| PR# | Title | Date | Why It's Important |
|-----|-------|------|-------------------|
| **#1174-1175** | pom dependencies upgraded for snyk | 2025-01-06/07 | Dependency preparation |
| **#1176** | SpringBoot 3 migration for NRT Final | 2025-01-07 | **THE MAIN PR** - Full migration |
| **#1153** | Testing fix for logging issue in spring Boot 3 | 2024-12-19 | Post-migration fixes |
| **#1312** | Feature/springboot 3 migration | 2025-04-17 | Development branch |
| **#1337** | Spring Boot 3 migration (#1312) to main | 2025-04-24 | **Production deployment** |
| **#1332-1335** | snyk issues fixed and mutable list issue | 2025-04-23 | Post-migration fixes |

#### All Migration-Related PRs

```
#1153 - Testing fix for logging issue in spring Boot 3 migration (2024-12-19)
#1161 - Application rollback (2024-12-20)
#1163 - Empty-Commit rollback (2024-12-21)
#1174 - pom dependencies upgraded (2025-01-06)
#1175 - pom dependencies upgraded for snyk issues (#1174) (2025-01-07)
#1176 - SpringBoot 3 migration for NRT Final (2025-01-07) *** MAIN PR ***
#1304 - fixing walmart-cosmos and azure-cosmos version for TLS (2025-04-14)
#1312 - Feature/springboot 3 migration (2025-04-17)
#1332 - snyk issues fixed and mutable list issue (2025-04-23)
#1334 - snyk issues fixed (#1332) (2025-04-23)
#1335 - one major sonar issue fixed (2025-04-23)
#1337 - Feature/springboot 3 migration (#1312) to main (2025-04-24) *** PROD ***
```

#### Interview Talking Points
- **Jakarta EE Migration**: javax.* → jakarta.* namespace changes
- **Hibernate 6.x**: Query syntax changes, type handling
- **Zero Downtime Strategy**: Kubernetes rolling updates, Istio traffic shifting
- **Rollback plan**: Notice #1161, #1163 - you had issues and rolled back
- **Dependency hell**: How did you resolve conflicting versions?
- **Testing strategy**: How did you ensure compatibility?

---

## Bullet 5: OpenAPI Design-First

### Resume Bullet
> Revamped all NRT application controllers using a design-first approach with OpenAPI Specification, reducing integration overhead by 30%.

### Repositories: Multiple

#### Key PRs (Interview-Worthy)

| Repo | PR# | Title | Date | Why It's Important |
|------|-----|-------|------|-------------------|
| `inventory-status-srv` | **#249** | DSIDVL-44382-Store-Inbound-Api-contract | 2025-06-06 | Store Inbound API spec |
| `inventory-status-srv` | **#260** | Feature/dc-inventory-api-spec | 2025-07-08 | DC Inventory API spec |
| `inventory-status-srv` | **#219, #233** | Contract Testing Integration | 2025-04-23/28 | R2C testing setup |
| `inventory-status-srv` | **#262** | API publishing automation | 2025-07-09 | Automation pipeline |
| `inventory-events-srv` | **#5** | API spec development | 2025-02-12 | Events API spec |
| `inventory-events-srv` | **#25** | Contract testing | 2025-03-21 | R2C for events |
| `cp-nrti-apis` | **#1183-1184** | ReadMe API Spec | 2025-01-14 | Documentation specs |

#### All OpenAPI/Contract Testing PRs

**inventory-status-srv:**
```
#219 - Contract Testing Integration (2025-04-23)
#233 - Contract Testing Integration (#219) (2025-04-28)
#234 - Dev Deployment for Contract Testing (2025-04-29)
#235 - Stage Deployment for Contract Test (2025-04-29)
#249 - DSIDVL-44382-Store-Inbound-Api-contract (2025-06-06)
#260 - Feature/dc-inventory-api-spec (2025-07-08)
#262 - Feature/dsidvl 46057 api publishing automation (2025-07-09)
#277 - minor fix in wm_item_nbr variable name (2025-08-21)
#312 - fix dc inventory api spec (2025-11-17)
#328 - Modified Contract tests to Passive mode (2025-12-02)
#331 - Update Contract test gate for us (2025-12-15)
#334 - Fixed Linter Issues (2026-01-06)
```

**inventory-events-srv:**
```
#5  - API spec development (2025-02-12)
#25 - Feature/contract testing (2025-03-21)
#34 - fixed stage contract testing (2025-03-25)
#35 - fixed stage contract testing (#34) (2025-03-25)
```

**cp-nrti-apis:**
```
#1169 - Added ReadMe Spec file (2024-12-23)
#1183 - Added API Spec file for NRTI Readme project (2025-01-14)
#1184 - Added API Spec file for NRTI Readme Project [skip-cli] (2025-01-14)
#1185 - Added ReadMe prod domain origin in ccm (2025-01-14)
```

#### Interview Talking Points
- **Design-first vs Code-first**: Why you chose design-first
- **openapi-generator-maven-plugin**: How it generates server stubs
- **Contract Testing (R2C)**: Request to Contract testing approach
- **Spec versioning**: How did you handle breaking changes?
- **Consumer collaboration**: How did you work with API consumers?
- **Measuring 30% reduction**: What specifically was reduced?

---

## Bullet 6: DC Inventory Search API

### Resume Bullet (Potential)
> Built DC inventory search and store inventory query APIs supporting bulk operations (100 items/request) with CompletableFuture parallel processing, reducing supplier query time by 40%.

### Repository: `inventory-status-srv`

#### Key PRs (Interview-Worthy)

| PR# | Title | Date | Why It's Important |
|-----|-------|------|-------------------|
| **#45** | /inventory/search-items API DEV Changes | 2025-01-03 | Initial API development |
| **#58** | Inventory search-items API develop to stage | 2025-01-21 | Stage deployment |
| **#271** | Feature/dc inventory development | 2025-07-24 | **CORE PR** - Main implementation |
| **#260** | Feature/dc-inventory-api-spec | 2025-07-08 | API design/spec |
| **#322** | V2 DSIDVL-48592 error handling dc inv | 2025-11-21 | Error handling improvements |
| **#330** | Feature/dc inv error fix | 2025-12-03 | Bug fixes |
| **#338** | container test for dc inventory from main | 2026-01-13 | Testing |

#### All DC Inventory PRs

```
#45  - /inventory/search-items API DEV Changes (2025-01-03)
#48  - testing /inventory/search-items to /search-items (2025-01-09)
#50  - Uri transform policy [skip-ci] (2025-01-10)
#51  - Uri transform policy (2025-01-13)
#58  - Inventory search-items API develop to stage (2025-01-21)
#260 - Feature/dc-inventory-api-spec (2025-07-08)
#271 - Feature/dc inventory development (2025-07-24) *** CORE PR ***
#277 - minor fix in wm_item_nbr variable name (2025-08-21)
#312 - fix dc inventory api spec (2025-11-17)
#322 - V2 DSIDVL-48592 error handling dc inv (2025-11-21)
#330 - Feature/dc inv error fix (2025-12-03)
#338 - container test for dc inventory from main (2026-01-13)
```

#### Interview Talking Points
- **Bulk API design**: Why 100 item limit per request?
- **CompletableFuture**: How you implemented parallel processing
- **Partial failure handling**: MultiStatusResponse pattern
- **UberKey integration**: Walmart's item identification system
- **Performance optimization**: How you achieved 40% reduction

---

## Bullet 7: Transaction Event History API

### Resume Bullet (Potential)
> Developed supplier-facing transaction event history API with multi-tenant architecture, GTIN-level authorization, and pagination support, serving 10+ event types across 300+ store locations.

### Repository: `inventory-events-srv`

#### Key PRs (Interview-Worthy)

| PR# | Title | Date | Why It's Important |
|-----|-------|------|-------------------|
| **#5** | API spec development | 2025-02-12 | API design |
| **#8** | DSIDVL-40282 transaction history | 2025-02-25 | **CORE PR** - Initial implementation |
| **#10** | Junit, Sonar, Snyk, Codegate | 2025-02-28 | Testing & quality |
| **#22** | Transaction history develop to stage | 2025-03-21 | Stage deployment |
| **#38** | TransactionHistoryApi Changes Stage to Main | 2025-03-26 | **Production deployment** |
| **#80** | Cosmos to Postgres migration | 2025-05-14 | Database migration |
| **#96** | CA Launch, Cosmos-Postgres, CCM | 2025-06-10 | Multi-market support |

#### All Transaction History PRs

```
#5  - API spec development (2025-02-12)
#7  - Sandbox development (2025-02-24)
#8  - DSIDVL-40282 transaction history (2025-02-25) *** CORE PR ***
#9  - Update sr.yaml (2025-02-28)
#10 - DSIDVL-40848 Junit Sonar Snyk Codegate (2025-02-28)
#13 - updated kitt config (2025-03-07)
#14 - sr.yaml configuration (2025-03-10)
#15 - Add transaction in uri transform policy (2025-03-13)
#17 - Reverted transactions in uri-mapping-policy (2025-03-13)
#22 - Transaction history develop to stage (#8) (2025-03-21)
#24 - Fixed sonar issues (2025-03-21)
#25 - Feature/contract testing (2025-03-21)
#36 - 405 error code fix (2025-03-25)
#37 - 405 error code fix (#36) (2025-03-25)
#38 - TransactionHistoryApi Changes Stage to Main (2025-03-26) *** PROD ***
#80 - Cosmos to Postgres migration (2025-05-14)
#96 - CA Launch, Cosmos-Postgres, CCM [CRQ:CHG3238236] (2025-06-10)
```

#### Interview Talking Points
- **Multi-tenant architecture**: Site-based partitioning (US/CA/MX)
- **Authorization hierarchy**: Consumer → DUNS → GTIN → Store
- **Cursor-based pagination**: Why cursor vs offset?
- **Event types**: 10+ distinct transaction types (Sales, Returns, etc.)
- **Cosmos to Postgres migration**: Why and how?

---

## Bullet 8: Supplier Authorization Framework

### Resume Bullet (Potential)
> Implemented multi-level supplier authorization framework (Consumer→DUNS→GTIN→Store) with PostgreSQL partition keys, Caffeine caching (7-day TTL), and site-aware queries, securing access to 10,000+ GTINs.

### Repositories: `inventory-status-srv` & `cp-nrti-apis`

#### Key PRs (Interview-Worthy)

| Repo | PR# | Title | Date | Why It's Important |
|------|-----|-------|------|-------------------|
| `inventory-status-srv` | **#161** | inventory-status store-GTIN enhancements | 2025-02-19 | **GTIN authorization** |
| `inventory-status-srv` | **#197, #216** | Cosmos to Postgres migration | 2025-04-09/22 | Auth data storage |
| `inventory-status-srv` | **#320** | Adding CCM-ANY value and siteId as tenantId | 2025-11-20 | Multi-tenant |
| `cp-nrti-apis` | **#1091** | Item validation Endpoint bug fix | 2024-10-25 | Auth validation |
| `cp-nrti-apis` | **#1120** | Item validation Endpoint bug fix (#1091) | 2024-11-26 | Auth validation |
| `cp-nrti-apis` | **#1478** | Cosmos to Postgres Migration | 2025-09-15 | Multi-tenant auth |

#### All Authorization-Related PRs

**inventory-status-srv:**
```
#142 - Added storeNbr lookup w.r.t siteId (2025-02-06)
#146 - Added storeNbr lookup check w.r.t siteId (#142) (2025-02-11)
#161 - inventory-status store-GTIN enhancements [CRQ:CHG3014193] (2025-02-19) *** KEY ***
#197 - Feature/cosmos to postgres migration (2025-04-09)
#216 - Feature/cosmos to postgres migration (#197) (2025-04-22)
#320 - Adding CCM-ANY value and siteId as tenantId (2025-11-20)
```

**cp-nrti-apis:**
```
#1091 - DSIDVL-35561 Item validation Endpoint bug fix (2024-10-25)
#1120 - Item validation Endpoint bug fix (#1091) (2024-11-26)
#1478 - Cosmos to Postgres Migration (#1427) (2025-09-15)
#1487 - Removed Supplier Persona and Modified access type (2025-09-16)
#1491 - Removed Supplier Persona (#1487) (2025-09-17)
```

#### Interview Talking Points
- **4-level authorization hierarchy**: Why this design?
- **Caffeine caching**: 7-day TTL trade-offs
- **Cache invalidation strategy**: How did you handle auth changes?
- **PostgreSQL partition keys**: Site isolation
- **Performance impact**: Auth checks on every request

---

## Bullet 9: Observability Stack

### Resume Bullet (Potential)
> Implemented comprehensive observability stack with OpenTelemetry distributed tracing, Prometheus metrics, Dynatrace APM, and custom Grafana dashboards, achieving 99.9% uptime with automated canary deployments.

### Repositories: `cp-nrti-apis` & `inventory-status-srv`

#### Key PRs (Interview-Worthy)

| Repo | PR# | Title | Date | Why It's Important |
|------|-----|-------|------|-------------------|
| `inventory-status-srv` | **#103, #105** | Enabling Dynatrace Saas | 2025-01-29 | **APM integration** |
| `inventory-status-srv` | **#136-137** | Added metrics in Kitt, different spanIds | 2025-02-05 | **Distributed tracing** |
| `inventory-status-srv` | **#327** | Added metrics in Kitt | 2025-12-02 | Prometheus metrics |
| `cp-nrti-apis` | **#1073** | Logging fix main | 2024-10-17 | Observability logging |
| `cp-nrti-apis` | **#1113, #1118** | Enabling canary deployment with Custom metrics | 2024-11-20/25 | **Flagger/Canary** |
| `cp-nrti-apis` | **#1208** | Enabling canary deployment with Custom metrics | 2025-01-21 | Prod canary |
| `cp-nrti-apis` | **#1545** | Added Dynatrace Saas | 2025-10-17 | APM |

#### All Observability PRs

**inventory-status-srv:**
```
#77  - inventory-status-service : Fixing TraceId issue (2025-01-24)
#83  - Fixing TraceId issue (#77) (2025-01-27)
#103 - inventory-status-service : Enabling Dynatrace Saas (2025-01-29)
#105 - Enabling Dynatrace (#103) (2025-01-29)
#136 - Added metrics in Kitt, Added different spanIds (2025-02-05)
#137 - Added metrics in Kitt, Added different spanIds (#136) (2025-02-05)
#163 - Pod profiler (2025-02-24)
#327 - Added metrics in Kitt (2025-12-02)
#335 - Added probe config (2026-01-07)
#337 - Added resiliency test (2026-01-09)
```

**cp-nrti-apis:**
```
#1053 - logs context added before exception is raised (#993) (2024-10-16)
#1056 - log updated (2024-10-16)
#1069 - log updated (2024-10-17)
#1073 - DSIDVL_33564_Logging_fix_main (2024-10-17)
#1113 - DSIDVL-18013: Enabling canary deployment with Custom metrics (2024-11-20)
#1118 - Canary deployment dev-to-stage (#1113) (2024-11-25)
#1208 - Enabling canary deployment with Custom metrics (2025-01-21)
#1220 - triggering canary deployment (2025-01-22)
#1271 - enable product metrics (2025-03-28)
#1280 - [CRQ:CHG3024893] prod deployment (2025-04-03)
#1469 - logs and ccm updated (2025-09-12)
#1474 - logs and ccm updated (#1469) (2025-09-12)
#1527 - Correlation Id Issue Fix (2025-10-06)
#1545 - Added Dynatrace Saas (2025-10-17)
```

#### Interview Talking Points
- **OpenTelemetry vs vendor-specific**: Why OpenTelemetry?
- **SLIs/SLOs**: What did you define?
- **Alert fatigue**: How did you manage it?
- **Canary deployments**: Flagger + Istio pattern
- **Custom metrics**: What business metrics did you track?

---

## Quick Reference: Top PRs per Bullet

| Bullet | Must-Review PRs | Repo |
|--------|-----------------|------|
| **1. Kafka Audit** | #1, #28, #70 | `audit-api-logs-gcs-sink` |
| **2. Common Library** | #1, #6, #11 | `dv-api-common-libraries` |
| **3. DSD Notifications** | #1129, #1397, #1433 | `cp-nrti-apis` |
| **4. Spring Boot 3** | #1176, #1312, #1337 | `cp-nrti-apis` |
| **5. OpenAPI** | #249, #260, #219 | `inventory-status-srv` |
| **6. DC Inventory** | #271, #322, #338 | `inventory-status-srv` |
| **7. Transaction History** | #8, #38, #80 | `inventory-events-srv` |
| **8. Auth Framework** | #161, #197 | `inventory-status-srv` |
| **9. Observability** | #1113, #103, #137 | `cp-nrti-apis` / `inventory-status-srv` |
| **10. Multi-Region Kafka** | #12, #85, #70 | `audit-api-logs-gcs-sink` |

---

## How to Deep Dive into a PR

### Using GitHub CLI

```bash
# View PR details
GH_HOST=gecgithub01.walmart.com gh pr view <PR#> --repo dsi-dataventures-luminate/<repo-name>

# View PR diff
GH_HOST=gecgithub01.walmart.com gh pr diff <PR#> --repo dsi-dataventures-luminate/<repo-name>

# View PR files changed
GH_HOST=gecgithub01.walmart.com gh pr view <PR#> --repo dsi-dataventures-luminate/<repo-name> --json files
```

### Example for Bullet 1 (Kafka Audit):
```bash
# View the main setup PR
GH_HOST=gecgithub01.walmart.com gh pr view 1 --repo dsi-dataventures-luminate/audit-api-logs-gcs-sink

# View the performance tuning PR
GH_HOST=gecgithub01.walmart.com gh pr view 28 --repo dsi-dataventures-luminate/audit-api-logs-gcs-sink
```

---

## Interview Preparation Checklist

For each bullet point, be ready to discuss:

- [ ] **What**: High-level description of what you built
- [ ] **Why**: Business problem it solved
- [ ] **How**: Technical implementation details
- [ ] **Trade-offs**: Decisions you made and alternatives considered
- [ ] **Metrics**: Quantifiable impact (events/day, teams adopted, % improvement)
- [ ] **Challenges**: Obstacles you faced and how you overcame them
- [ ] **Learnings**: What you would do differently

---

*Generated for interview preparation - Good luck!*

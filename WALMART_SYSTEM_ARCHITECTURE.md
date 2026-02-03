# WALMART SYSTEM ARCHITECTURE
## Complete System Diagrams and Technical Architecture

**Author**: Anshul Garg
**Team**: Data Ventures - Channel Performance Engineering
**Architecture Scope**: 6 microservices across multi-region, multi-market deployment

---

# TABLE OF CONTENTS

1. [High-Level Architecture Overview](#high-level-architecture-overview)
2. [Multi-Tenant Architecture](#multi-tenant-architecture)
3. [API Gateway and Service Mesh](#api-gateway-and-service-mesh)
4. [Database Architecture](#database-architecture)
5. [Kafka Event Streaming Architecture](#kafka-event-streaming-architecture)
6. [External Service Integrations](#external-service-integrations)
7. [Deployment Architecture](#deployment-architecture)
8. [Security Architecture](#security-architecture)
9. [Observability Architecture](#observability-architecture)
10. [Network Flow Diagrams](#network-flow-diagrams)

---

# 1. HIGH-LEVEL ARCHITECTURE OVERVIEW

## System Context Diagram

```
                                    ┌─────────────────────────────────────┐
                                    │   External Suppliers                │
                                    │   • 1,200+ suppliers                │
                                    │   • 10,000+ GTINs                   │
                                    │   • US, Canada, Mexico markets      │
                                    └─────────────┬───────────────────────┘
                                                  │
                                                  │ HTTPS
                                                  │
                    ┌─────────────────────────────▼─────────────────────────────┐
                    │                   API Gateway (Torbit)                     │
                    │             Rate Limiting: 900 TPM per consumer            │
                    │          OAuth 2.0 + Consumer ID validation                │
                    └─────────────┬─────────────────────────┬───────────────────┘
                                  │                         │
                                  │                         │
                    ┌─────────────▼─────────┐   ┌──────────▼──────────┐
                    │   Service Registry    │   │   Istio Service     │
                    │   (Application Keys)  │   │   Mesh (mTLS)       │
                    └─────────────┬─────────┘   └──────────┬──────────┘
                                  │                         │
                                  └────────────┬────────────┘
                                               │
            ┌──────────────────────────────────┼──────────────────────────────────┐
            │                                  │                                  │
            │                                  │                                  │
    ┌───────▼───────┐              ┌──────────▼─────────┐            ┌──────────▼─────────┐
    │ cp-nrti-apis  │              │ inventory-         │            │ inventory-         │
    │               │              │ status-srv         │            │ events-srv         │
    │ 10+ endpoints │              │                    │            │                    │
    │ IAC, DSC, TH  │              │ 3 endpoints        │            │ 1 endpoint         │
    └───────┬───────┘              │ DC/Store Inventory │            │ Transaction History│
            │                      └──────────┬─────────┘            └──────────┬─────────┘
            │                                 │                                  │
            │                                 │                                  │
            └─────────────────────────────────┼──────────────────────────────────┘
                                              │
                                              │
                    ┌─────────────────────────▼─────────────────────────┐
                    │           PostgreSQL Database (Multi-Tenant)      │
                    │   • nrt_consumers (supplier mappings)             │
                    │   • supplier_gtin_items (GTIN authorization)      │
                    │   • Partition Keys: site_id (US/CA/MX)            │
                    └─────────────────────────┬───────────────────────── ┘
                                              │
                                              │
                    ┌─────────────────────────▼─────────────────────────┐
                    │         External Service Integrations             │
                    ├───────────────────────────────────────────────────┤
                    │  • Enterprise Inventory API (EI)                  │
                    │    - US: ei-inventory-history-lookup.walmart.com  │
                    │    - CA: ei-inventory-history-lookup-ca.walmart...│
                    │    - MX: ei-inventory-history-lookup-mx.walmart...│
                    │  • UberKey Service (GTIN ↔ WM Item Number)        │
                    │  • BigQuery (Analytics)                           │
                    │  • Sumo (Push Notifications)                      │
                    │  • CCM2/Tunr (Configuration Management)           │
                    │  • Akeyless (Secrets Management)                  │
                    └───────────────────────────────────────────────────┘

    ┌─────────────────────────────────────────────────────────────────────────┐
    │                    Audit Logging Pipeline                               │
    ├─────────────────────────────────────────────────────────────────────────┤
    │  dv-api-common-libraries (Filter)                                       │
    │         ↓                                                                │
    │  audit-api-logs-srv (Kafka Producer)                                    │
    │         ↓                                                                │
    │  Kafka Topics: audit-logs-us, audit-logs-ca, audit-logs-mx             │
    │         ↓                                                                │
    │  audit-api-logs-gcs-sink (Kafka Connect with SMT Filters)              │
    │         ↓                                                                │
    │  GCS Buckets (Parquet Files) → BigQuery Tables                          │
    └─────────────────────────────────────────────────────────────────────────┘

    ┌─────────────────────────────────────────────────────────────────────────┐
    │                    Observability Stack                                  │
    ├─────────────────────────────────────────────────────────────────────────┤
    │  • OpenTelemetry (Distributed Tracing)                                  │
    │  • Prometheus (Metrics)                                                 │
    │  • Grafana (Dashboards)                                                 │
    │  • Dynatrace SaaS (APM)                                                 │
    │  • Wolly (Log Aggregation)                                              │
    └─────────────────────────────────────────────────────────────────────────┘
```

---

## Key Architecture Principles

### 1. Multi-Tenant Architecture
- Site-based partitioning (site_id: 1=US, 2=CA, 3=MX)
- ThreadLocal context propagation
- Hibernate partition keys for automatic filtering
- Site-specific configuration (EI endpoints per market)

### 2. Event-Driven Architecture
- Kafka for asynchronous event publishing
- Fire-and-forget pattern (non-blocking)
- Multi-region Kafka clusters (EUS2, SCUS)
- Audit logging pipeline (2M+ events daily)

### 3. Microservices Architecture
- Domain-driven design (inventory-status, inventory-events, cp-nrti)
- Independent deployment and scaling
- Service mesh for inter-service communication
- API Gateway for external access

### 4. Cloud-Native Architecture
- Kubernetes (WCNP) for orchestration
- Docker containers
- HPA for autoscaling (4-8 pods in production)
- Multi-region deployment (EUS2, SCUS, USWEST, USEAST)

---

# 2. MULTI-TENANT ARCHITECTURE

## Site-Based Partitioning

```
┌─────────────────────────────────────────────────────────────────────────┐
│                      Incoming Request                                   │
│   Headers: WM-Site-Id: 2 (Canada)                                       │
│            Authorization: Bearer <token>                                │
│            wm_consumer.id: <supplier-uuid>                              │
└────────────────────────────┬────────────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                    SiteContextFilter (@Order(1))                        │
│                                                                         │
│  1. Extract WM-Site-Id header                                           │
│  2. Parse site ID: "2" → Long(2)                                        │
│  3. Set in ThreadLocal: siteContext.setSiteId(2L)                       │
│  4. Continue filter chain                                               │
│  5. Finally: siteContext.clear()                                        │
└────────────────────────────┬────────────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                       Controller Layer                                  │
│                                                                         │
│  @PostMapping("/v1/inventory/search-items")                             │
│  public ResponseEntity<InventoryResponse> search(@RequestBody ...) {    │
│      // Site context already set by filter                              │
│      return service.getInventory(request);                              │
│  }                                                                       │
└────────────────────────────┬────────────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                       Service Layer                                     │
│                                                                         │
│  public InventoryResponse getInventory(InventoryRequest request) {      │
│      Long siteId = siteContext.getSiteId();  // Get from ThreadLocal   │
│                                                                         │
│      // Get site-specific configuration                                │
│      SiteConfig config = siteConfigFactory.getConfigurations(siteId);  │
│      String eiEndpoint = config.getEndpoint();                          │
│      // eiEndpoint = "ei-inventory-history-lookup-ca.walmart.com"      │
│                                                                         │
│      // Database queries automatically filtered by site_id             │
│      ParentCompanyMapping supplier = supplierRepo.find(consumerId, siteId);
│  }                                                                       │
└────────────────────────────┬────────────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                     Repository Layer (JPA)                              │
│                                                                         │
│  @Repository                                                            │
│  public interface ParentCmpnyMappingRepository                          │
│      extends JpaRepository<ParentCompanyMapping, ...> {                 │
│                                                                         │
│      Optional<ParentCompanyMapping> findByConsumerIdAndSiteId(          │
│          String consumerId, String siteId);                             │
│  }                                                                       │
│                                                                         │
│  // Generated SQL (Hibernate adds site_id automatically):              │
│  SELECT * FROM nrt_consumers                                            │
│  WHERE consumer_id = ?                                                  │
│    AND site_id = ?  ← Automatic partition key filtering                │
└────────────────────────────┬────────────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                   PostgreSQL Database                                   │
│                                                                         │
│  Table: nrt_consumers                                                   │
│  ┌──────────────┬──────────┬─────────────┬──────────────┐             │
│  │ consumer_id  │ site_id  │ global_duns │ country_code │             │
│  ├──────────────┼──────────┼─────────────┼──────────────┤             │
│  │ abc-123-def  │ 1        │ 012345678   │ US           │  ← US data  │
│  │ xyz-456-ghi  │ 2        │ 987654321   │ CA           │  ← CA data  │
│  │ mno-789-pqr  │ 3        │ 555555555   │ MX           │  ← MX data  │
│  └──────────────┴──────────┴─────────────┴──────────────┘             │
│                                                                         │
│  Partition Key: site_id                                                │
│  Composite Primary Key: (consumer_id, site_id)                         │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## Site Context Propagation to Worker Threads

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    Parent Thread (Request Thread)                       │
│                                                                         │
│  siteContext.setSiteId(2L);  ← Set by SiteContextFilter                │
│                                                                         │
│  List<CompletableFuture<Result>> futures = items.stream()               │
│      .map(item -> CompletableFuture.supplyAsync(                        │
│          () -> processItem(item),  ← This runs in worker thread        │
│          taskExecutor  ← Custom executor with SiteTaskDecorator        │
│      ))                                                                 │
│      .collect(Collectors.toList());                                     │
│                                                                         │
│  CompletableFuture.allOf(futures...).join();                            │
└────────────────────────────┬────────────────────────────────────────────┘
                             │
                             │ SiteTaskDecorator intercepts
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                     SiteTaskDecorator                                   │
│                                                                         │
│  public Runnable decorate(Runnable runnable) {                          │
│      Long siteId = siteContext.getSiteId();  ← Capture from parent     │
│                                                                         │
│      return () -> {                                                     │
│          try {                                                          │
│              siteContext.setSiteId(siteId);  ← Set in worker thread    │
│              runnable.run();  ← Execute actual task                     │
│          } finally {                                                    │
│              siteContext.clear();  ← Clean up                           │
│          }                                                              │
│      };                                                                 │
│  }                                                                       │
└────────────────────────────┬────────────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────────────┐
│              Worker Thread (From ThreadPoolTaskExecutor)                │
│                                                                         │
│  siteContext.getSiteId();  → Returns 2L (Canada)                        │
│                                                                         │
│  // Now worker thread has correct site context                         │
│  // Database queries will be filtered for Canadian data                │
│                                                                         │
│  Result result = processItem(item);                                     │
│  return result;                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

**Why This Is Critical**:
- Without SiteTaskDecorator: Worker thread has NO site context → queries all data (data leakage!)
- With SiteTaskDecorator: Worker thread has correct site context → queries only Canadian data

---

## Site-Specific Configuration Factory

```
┌─────────────────────────────────────────────────────────────────────────┐
│                   SiteConfigFactory                                     │
│                                                                         │
│  Map<String, SiteConfigMapper> configMap = Map.of(                      │
│      "1", usConfig,   // US configuration                               │
│      "2", caConfig,   // Canada configuration                           │
│      "3", mxConfig    // Mexico configuration                           │
│  );                                                                      │
│                                                                         │
│  public SiteConfigMapper getConfigurations(Long siteId) {               │
│      return configMap.get(String.valueOf(siteId));                      │
│  }                                                                       │
└─────────────────────────────┬───────────────────────────────────────────┘
                              │
                              │ siteId = 2 (Canada)
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                  Canada Configuration (caConfig)                        │
│                                                                         │
│  EI API Endpoint:                                                       │
│      "https://ei-inventory-history-lookup-ca.walmart.com/v1"            │
│                                                                         │
│  Consumer ID:                                                           │
│      "ca-consumer-id-uuid"                                              │
│                                                                         │
│  Key Version:                                                           │
│      "1"                                                                │
│                                                                         │
│  Site ID:                                                               │
│      "2"                                                                │
└─────────────────────────────────────────────────────────────────────────┘
```

**Routing Logic**:
- Request with `WM-Site-Id: 2` → Canadian configuration → Canadian EI endpoint
- Request with `WM-Site-Id: 1` → US configuration → US EI endpoint
- Request with `WM-Site-Id: 3` → Mexican configuration → Mexican EI endpoint

---

# 3. API GATEWAY AND SERVICE MESH

## API Gateway (Torbit + Service Registry)

```
┌─────────────────────────────────────────────────────────────────────────┐
│                         External Supplier                               │
│                                                                         │
│  POST https://developer.walmart.com/api/us/inventory/v1/search-items    │
│  Headers:                                                               │
│      Authorization: Bearer <OAuth-token>                                │
│      wm_consumer.id: <supplier-uuid>                                    │
│      WM-Site-Id: 1                                                      │
│  Body: { "store_nbr": 3188, "item_type_values": [...] }                │
└────────────────────────────┬────────────────────────────────────────────┘
                             │
                             │ HTTPS (TLS 1.2+)
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                      API Gateway (Torbit)                               │
│                                                                         │
│  1. TLS Termination                                                     │
│  2. OAuth 2.0 Token Validation                                          │
│     ┌──────────────────────────────────────────────────────┐           │
│     │ Token Validation:                                    │           │
│     │   • Validate signature                               │           │
│     │   • Check expiration                                 │           │
│     │   • Verify scopes                                    │           │
│     │ IDP: https://idp.prod.global.sso.platform.prod...   │           │
│     └──────────────────────────────────────────────────────┘           │
│                                                                         │
│  3. Rate Limiting                                                       │
│     ┌──────────────────────────────────────────────────────┐           │
│     │ Rate Limit Policy:                                   │           │
│     │   • Criteria: wm_consumer.id header                  │           │
│     │   • Limit: 900 TPM (Transactions Per Minute)         │           │
│     │   • Window: Sliding window                           │           │
│     │   • Response: HTTP 429 (Too Many Requests)           │           │
│     └──────────────────────────────────────────────────────┘           │
│                                                                         │
│  4. Service Registry Lookup                                             │
│     ┌──────────────────────────────────────────────────────┐           │
│     │ Application Key: INVENTORY-STATUS                    │           │
│     │ Environment: prod                                    │           │
│     │ Target: inventory-status-srv.prod.svc.cluster.local │           │
│     └──────────────────────────────────────────────────────┘           │
└────────────────────────────┬────────────────────────────────────────────┘
                             │
                             │ Forward to Kubernetes
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                  Kubernetes Service (Load Balancer)                     │
│                                                                         │
│  Service: inventory-status-srv                                          │
│  Type: ClusterIP                                                        │
│  Port: 8080                                                             │
│                                                                         │
│  Load Balancing Algorithm: Round Robin                                 │
│                                                                         │
│  Endpoints:                                                             │
│  ┌──────────────────────────────────────────────────────┐              │
│  │ Pod 1: 10.244.1.10:8080  (EUS2-PROD-A30)             │              │
│  │ Pod 2: 10.244.1.11:8080  (EUS2-PROD-A30)             │              │
│  │ Pod 3: 10.244.1.12:8080  (SCUS-PROD-A63)             │              │
│  │ Pod 4: 10.244.1.13:8080  (SCUS-PROD-A63)             │              │
│  └──────────────────────────────────────────────────────┘              │
└────────────────────────────┬────────────────────────────────────────────┘
                             │
                             │ Route to Pod 1
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                  Istio Sidecar (Envoy Proxy)                            │
│                                                                         │
│  1. mTLS (Mutual TLS) Encryption                                        │
│     ┌──────────────────────────────────────────────────────┐           │
│     │ Certificate-based authentication                     │           │
│     │ Auto-rotated certificates (every 24 hours)           │           │
│     │ Encrypted inter-service communication                │           │
│     └──────────────────────────────────────────────────────┘           │
│                                                                         │
│  2. Traffic Management                                                  │
│     ┌──────────────────────────────────────────────────────┐           │
│     │ Circuit Breaker:                                     │           │
│     │   • Max Connections: 100                             │           │
│     │   • Max Pending Requests: 50                         │           │
│     │   • Max Requests: 100                                │           │
│     │   • Consecutive Errors: 5                            │           │
│     │   • Interval: 30s                                    │           │
│     └──────────────────────────────────────────────────────┘           │
│                                                                         │
│  3. Local Rate Limiting                                                 │
│     ┌──────────────────────────────────────────────────────┐           │
│     │ Token Bucket:                                        │           │
│     │   • Max Tokens: 75                                   │           │
│     │   • Tokens Per Fill: 75                              │           │
│     │   • Fill Interval: 1 second                          │           │
│     │ Result: 75 requests/second per pod                   │           │
│     └──────────────────────────────────────────────────────┘           │
│                                                                         │
│  4. Metrics Collection                                                  │
│     • Request count, latency, errors                                   │
│     • Sent to Prometheus                                               │
└────────────────────────────┬────────────────────────────────────────────┘
                             │
                             │ Forward to application container
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                  Application Container (Pod 1)                          │
│                                                                         │
│  Container: inventory-status-srv:v1.2.3                                │
│  Port: 8080                                                             │
│  CPU: 1 core                                                            │
│  Memory: 1Gi                                                            │
│                                                                         │
│  @PostMapping("/v1/inventory/search-items")                             │
│  public ResponseEntity<InventoryResponse> search(...) {                 │
│      // Business logic                                                  │
│  }                                                                       │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## Service Mesh (Istio) Traffic Flow

```
                    ┌────────────────────────────────────┐
                    │   Istio Control Plane (Istiod)     │
                    │   • Service Discovery               │
                    │   • Certificate Management          │
                    │   • Configuration Distribution      │
                    └──────────────┬─────────────────────┘
                                   │
                                   │ Configuration
                                   │
                    ┌──────────────┴─────────────────────┐
                    │          Envoy Proxies             │
                    │        (Sidecar Containers)        │
                    └──────────────┬─────────────────────┘
                                   │
        ┌──────────────────────────┼──────────────────────────┐
        │                          │                          │
        │                          │                          │
┌───────▼────────┐        ┌────────▼───────┐        ┌────────▼───────┐
│ Pod 1          │        │ Pod 2          │        │ Pod 3          │
│ ┌────────────┐ │        │ ┌────────────┐ │        │ ┌────────────┐ │
│ │   Envoy    │ │        │ │   Envoy    │ │        │ │   Envoy    │ │
│ │   Proxy    │ │        │ │   Proxy    │ │        │ │   Proxy    │ │
│ └──────┬─────┘ │        │ └──────┬─────┘ │        │ └──────┬─────┘ │
│        │       │        │        │       │        │        │       │
│ ┌──────▼─────┐ │        │ ┌──────▼─────┐ │        │ ┌──────▼─────┐ │
│ │ inventory- │ │        │ │ inventory- │ │        │ │ inventory- │ │
│ │ status-srv │ │        │ │ status-srv │ │        │ │ events-srv │ │
│ └────────────┘ │        │ └────────────┘ │        │ └────────────┘ │
└────────────────┘        └────────────────┘        └────────────────┘
```

**Istio Features**:
1. **mTLS**: All service-to-service traffic encrypted
2. **Circuit Breaker**: Prevents cascading failures
3. **Retry Logic**: Automatic retries on transient failures
4. **Timeout Management**: Request timeouts
5. **Observability**: Automatic metrics collection

---

# 4. DATABASE ARCHITECTURE

## PostgreSQL Multi-Tenant Schema

```
┌─────────────────────────────────────────────────────────────────────────┐
│                       PostgreSQL Database                               │
│                                                                         │
│  Database: walmart_inventory                                            │
│  Version: PostgreSQL 14                                                │
│  Connection Pool: HikariCP                                              │
│    • Max Pool Size: 20                                                  │
│    • Min Idle: 5                                                        │
│    • Connection Timeout: 30000ms                                        │
│    • Idle Timeout: 600000ms                                             │
└─────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────┐
│                  Table: nrt_consumers                                   │
│                  (Supplier Metadata and Authentication)                 │
│                                                                         │
│  CREATE TABLE nrt_consumers (                                           │
│      consumer_id VARCHAR(255) NOT NULL,                                 │
│      site_id VARCHAR(10) NOT NULL,                                      │
│      consumer_name VARCHAR(255),                                        │
│      country_code VARCHAR(2),                                           │
│      global_duns VARCHAR(20),                                           │
│      psp_global_duns VARCHAR(20),                                       │
│      parent_cmpny_name VARCHAR(255),                                    │
│      luminate_cmpny_id VARCHAR(50),                                     │
│      is_category_manager BOOLEAN DEFAULT FALSE,                         │
│      non_charter_supplier BOOLEAN DEFAULT FALSE,                        │
│      status VARCHAR(20),                                                │
│      user_type VARCHAR(50),                                             │
│      persona VARCHAR(50),                                               │
│      PRIMARY KEY (consumer_id, site_id),                                │
│      INDEX idx_site_id (site_id),                                       │
│      INDEX idx_global_duns (global_duns)                                │
│  );                                                                      │
│                                                                         │
│  Partition Key: site_id                                                │
│  Composite Primary Key: (consumer_id, site_id)                         │
│                                                                         │
│  Sample Data:                                                           │
│  ┌──────────────┬──────────┬──────────────┬─────────┬─────────────┐   │
│  │ consumer_id  │ site_id  │ consumer_name│ country │ global_duns │   │
│  ├──────────────┼──────────┼──────────────┼─────────┼─────────────┤   │
│  │ abc-123-def  │ 1        │ ABC Corp     │ US      │ 012345678   │   │
│  │ xyz-456-ghi  │ 2        │ XYZ Corp     │ CA      │ 987654321   │   │
│  │ mno-789-pqr  │ 3        │ MNO Corp     │ MX      │ 555555555   │   │
│  └──────────────┴──────────┴──────────────┴─────────┴─────────────┘   │
└─────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────┐
│                Table: supplier_gtin_items                               │
│                (GTIN-to-Supplier Authorization Mapping)                 │
│                                                                         │
│  CREATE TABLE supplier_gtin_items (                                     │
│      global_duns VARCHAR(20) NOT NULL,                                  │
│      gtin VARCHAR(14) NOT NULL,                                         │
│      site_id VARCHAR(10) NOT NULL,                                      │
│      store_nbr INTEGER[],     -- PostgreSQL array                       │
│      luminate_cmpny_id VARCHAR(50),                                     │
│      parent_company_name VARCHAR(255),                                  │
│      PRIMARY KEY (global_duns, gtin, site_id),                          │
│      INDEX idx_gtin (gtin),                                             │
│      INDEX idx_site_id (site_id)                                        │
│  );                                                                      │
│                                                                         │
│  Partition Key: site_id, global_duns                                   │
│  Composite Primary Key: (global_duns, gtin, site_id)                   │
│                                                                         │
│  PostgreSQL Array Column (store_nbr):                                  │
│  ┌─────────────┬────────────────┬──────────┬──────────────────────┐    │
│  │ global_duns │ gtin           │ site_id  │ store_nbr            │    │
│  ├─────────────┼────────────────┼──────────┼──────────────────────┤    │
│  │ 012345678   │ 00012345678901 │ 1        │ {3188, 3067, 4456}   │    │
│  │ 012345678   │ 00012345678902 │ 1        │ {}  ← Empty = all    │    │
│  │ 987654321   │ 00012345678903 │ 2        │ {5001, 5002}         │    │
│  └─────────────┴────────────────┴──────────┴──────────────────────┘    │
│                                                                         │
│  PostgreSQL Array Operations:                                          │
│    • Check contains: store_nbr @> ARRAY[3188]                           │
│    • Check overlap: store_nbr && ARRAY[3188, 3067]                      │
│    • Empty array = authorized for all stores                           │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## Database Connection Architecture

```
┌─────────────────────────────────────────────────────────────────────────┐
│                   Application Pod (inventory-status-srv)                │
│                                                                         │
│  ┌───────────────────────────────────────────────────────────────┐     │
│  │              HikariCP Connection Pool                         │     │
│  │                                                               │     │
│  │  Configuration:                                               │     │
│  │    • Maximum Pool Size: 20 connections                        │     │
│  │    • Minimum Idle: 5 connections                              │     │
│  │    • Connection Timeout: 30 seconds                           │     │
│  │    • Idle Timeout: 10 minutes                                 │     │
│  │    • Max Lifetime: 30 minutes                                 │     │
│  │    • Leak Detection Threshold: 60 seconds                     │     │
│  │                                                               │     │
│  │  Connection Pool State:                                       │     │
│  │  ┌──────────────────────────────────────────────────────┐    │     │
│  │  │ [Conn1] [Conn2] [Conn3] [Conn4] [Conn5]             │    │     │
│  │  │   IDLE    IDLE    ACTIVE  ACTIVE  IDLE               │    │     │
│  │  │                                                       │    │     │
│  │  │ [Conn6] [Conn7] [Conn8] ... [Conn20]                 │    │     │
│  │  │  IDLE    IDLE    IDLE     ...  IDLE                   │    │     │
│  │  └──────────────────────────────────────────────────────┘    │     │
│  └───────────────────────┬───────────────────────────────────────┘     │
└──────────────────────────┼─────────────────────────────────────────────┘
                           │
                           │ JDBC Connection
                           │
┌──────────────────────────▼─────────────────────────────────────────────┐
│                    PostgreSQL Database Server                           │
│                                                                         │
│  Host: postgres.prod.walmart.internal                                  │
│  Port: 5432                                                             │
│  Database: walmart_inventory                                            │
│  SSL: Required (TLS 1.2+)                                               │
│                                                                         │
│  Connection String:                                                     │
│  jdbc:postgresql://postgres.prod.walmart.internal:5432/walmart_inventory
│      ?ssl=true&sslmode=require                                          │
│                                                                         │
│  Secrets (Akeyless):                                                    │
│    • /etc/secrets/db_username.txt                                       │
│    • /etc/secrets/db_password.txt                                       │
│    • /etc/secrets/db_url.txt                                            │
└─────────────────────────────────────────────────────────────────────────┘
```

**HikariCP Configuration Details**:
```java
@Configuration
public class PostgresDbConfiguration {
    @Bean
    public DataSource dataSource() {
        HikariConfig config = new HikariConfig();
        config.setJdbcUrl(readSecret("db_url.txt"));
        config.setUsername(readSecret("db_username.txt"));
        config.setPassword(readSecret("db_password.txt"));
        config.setMaximumPoolSize(20);
        config.setMinimumIdle(5);
        config.setConnectionTimeout(30000);
        config.setIdleTimeout(600000);
        config.setMaxLifetime(1800000);
        config.setLeakDetectionThreshold(60000);  // Detect connection leaks
        return new HikariDataSource(config);
    }
}
```

---

# 5. KAFKA EVENT STREAMING ARCHITECTURE

## Multi-Region Kafka Architecture

```
┌─────────────────────────────────────────────────────────────────────────┐
│                          EUS2 Region                                    │
│                                                                         │
│  ┌──────────────────────────────────────────────────────────────┐      │
│  │         Kafka Cluster (EUS2)                                 │      │
│  │                                                              │      │
│  │  Bootstrap Servers:                                          │      │
│  │    • kafka-broker-1.eus2.walmart.internal:9093               │      │
│  │    • kafka-broker-2.eus2.walmart.internal:9093               │      │
│  │    • kafka-broker-3.eus2.walmart.internal:9093               │      │
│  │                                                              │      │
│  │  Topics:                                                     │      │
│  │  ┌────────────────────────────────────────────────────┐     │      │
│  │  │ cperf-nrt-prod-iac (Inventory Actions)            │     │      │
│  │  │   Partitions: 3                                    │     │      │
│  │  │   Replication Factor: 3                            │     │      │
│  │  │   Retention: 7 days                                │     │      │
│  │  └────────────────────────────────────────────────────┘     │      │
│  │  ┌────────────────────────────────────────────────────┐     │      │
│  │  │ cperf-nrt-prod-dsc (Direct Shipment Capture)      │     │      │
│  │  │   Partitions: 2                                    │     │      │
│  │  │   Replication Factor: 3                            │     │      │
│  │  │   Retention: 7 days                                │     │      │
│  │  └────────────────────────────────────────────────────┘     │      │
│  │  ┌────────────────────────────────────────────────────┐     │      │
│  │  │ audit-logs-us                                      │     │      │
│  │  │   Partitions: 6                                    │     │      │
│  │  │   Replication Factor: 3                            │     │      │
│  │  │   Retention: 3 days                                │     │      │
│  │  └────────────────────────────────────────────────────┘     │      │
│  │  ┌────────────────────────────────────────────────────┐     │      │
│  │  │ audit-logs-ca                                      │     │      │
│  │  │   Partitions: 3                                    │     │      │
│  │  │   Replication Factor: 3                            │     │      │
│  │  │   Retention: 3 days                                │     │      │
│  │  └────────────────────────────────────────────────────┘     │      │
│  │  ┌────────────────────────────────────────────────────┐     │      │
│  │  │ audit-logs-mx                                      │     │      │
│  │  │   Partitions: 2                                    │     │      │
│  │  │   Replication Factor: 3                            │     │      │
│  │  │   Retention: 3 days                                │     │      │
│  │  └────────────────────────────────────────────────────┘     │      │
│  └──────────────────────────────────────────────────────────────┘      │
└─────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────┐
│                          SCUS Region                                    │
│                                                                         │
│  ┌──────────────────────────────────────────────────────────────┐      │
│  │         Kafka Cluster (SCUS)                                 │      │
│  │                                                              │      │
│  │  Bootstrap Servers:                                          │      │
│  │    • kafka-broker-1.scus.walmart.internal:9093               │      │
│  │    • kafka-broker-2.scus.walmart.internal:9093               │      │
│  │    • kafka-broker-3.scus.walmart.internal:9093               │      │
│  │                                                              │      │
│  │  Topics: (Same as EUS2)                                      │      │
│  │    • cperf-nrt-prod-iac                                      │      │
│  │    • cperf-nrt-prod-dsc                                      │      │
│  │    • audit-logs-us                                           │      │
│  │    • audit-logs-ca                                           │      │
│  │    • audit-logs-mx                                           │      │
│  └──────────────────────────────────────────────────────────────┘      │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## Audit Logging Pipeline (Complete Flow)

```
┌─────────────────────────────────────────────────────────────────────────┐
│                Step 1: HTTP Request Interception                        │
│                                                                         │
│  ┌──────────────────────────────────────────────────────────────┐      │
│  │  LoggingFilter (dv-api-common-libraries)                     │      │
│  │                                                              │      │
│  │  1. Intercept HTTP request/response                          │      │
│  │  2. Wrap with ContentCachingRequestWrapper                   │      │
│  │  3. Extract headers, body, status code                       │      │
│  │  4. Build AuditLogPayload                                    │      │
│  │  5. Async call to AuditLogService                            │      │
│  └──────────────────────────────────────────────────────────────┘      │
│                                                                         │
│  AuditLogPayload:                                                       │
│  ┌──────────────────────────────────────────────────────────────┐      │
│  │ {                                                            │      │
│  │   "request_id": "uuid-123",                                  │      │
│  │   "service_name": "inventory-status-srv",                    │      │
│  │   "endpoint_name": "/v1/inventory/search-items",             │      │
│  │   "method": "POST",                                          │      │
│  │   "request_body": "{...}",                                   │      │
│  │   "response_body": "{...}",                                  │      │
│  │   "response_code": 200,                                      │      │
│  │   "request_ts": 1710498600000,                               │      │
│  │   "response_ts": 1710498601000,                              │      │
│  │   "headers": {"wm_consumer.id": "...", "wm-site-id": "1"}    │      │
│  │ }                                                            │      │
│  └──────────────────────────────────────────────────────────────┘      │
└────────────────────────────┬────────────────────────────────────────────┘
                             │
                             │ @Async (Non-blocking)
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────────────┐
│             Step 2: Kafka Producer (audit-api-logs-srv)                 │
│                                                                         │
│  ┌──────────────────────────────────────────────────────────────┐      │
│  │  AuditLogService (@Async)                                    │      │
│  │                                                              │      │
│  │  1. Convert AuditLogPayload to JSON                          │      │
│  │  2. Add Walmart authentication headers:                      │      │
│  │     • WM_CONSUMER.ID                                         │      │
│  │     • WM_SEC.AUTH_SIGNATURE (RSA signature)                  │      │
│  │     • WM_SEC.KEY_VERSION                                     │      │
│  │     • WM_CONSUMER.INTIMESTAMP                                │      │
│  │  3. Publish to Kafka topic (fire-and-forget)                 │      │
│  └──────────────────────────────────────────────────────────────┘      │
│                                                                         │
│  Thread Pool Configuration:                                             │
│    • Core threads: 6                                                   │
│    • Max threads: 10                                                   │
│    • Queue capacity: 100                                               │
│    • Thread name prefix: "Audit-log-executor-"                         │
└────────────────────────────┬────────────────────────────────────────────┘
                             │
                             │ Kafka Publish
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────────────┐
│            Step 3: Kafka Topics (audit-logs-us/ca/mx)                   │
│                                                                         │
│  Kafka Message Format (Avro):                                           │
│  ┌──────────────────────────────────────────────────────────────┐      │
│  │ {                                                            │      │
│  │   "schema": "AuditLogSchema-v1",                             │      │
│  │   "payload": {                                               │      │
│  │     "request_id": "uuid-123",                                │      │
│  │     "service_name": "inventory-status-srv",                  │      │
│  │     "wm_site_id": "1",  ← Site ID for filtering             │      │
│  │     "endpoint_name": "/v1/inventory/search-items",           │      │
│  │     ...                                                      │      │
│  │   }                                                          │      │
│  │ }                                                            │      │
│  └──────────────────────────────────────────────────────────────┘      │
│                                                                         │
│  Partition Strategy: By service_name                                   │
│    • All inventory-status-srv messages → Partition 0                  │
│    • All cp-nrti-apis messages → Partition 1                          │
│    • etc.                                                              │
└────────────────────────────┬────────────────────────────────────────────┘
                             │
                             │ Kafka Connect Consumers
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────────────┐
│      Step 4: Kafka Connect Sink (audit-api-logs-gcs-sink)              │
│                                                                         │
│  Multi-Connector Pattern (3 connectors in parallel):                   │
│                                                                         │
│  ┌──────────────────────────────────────────────────────────────┐      │
│  │  Connector 1: US Connector                                   │      │
│  │                                                              │      │
│  │  1. Consume from: audit-logs-us                              │      │
│  │  2. SMT Filter: AuditLogSinkUSFilter                         │      │
│  │     • Accept if wm-site-id = "US" OR missing                 │      │
│  │     • Permissive filter (default to US)                      │      │
│  │  3. Destination: gs://walmart-audit-logs-us/                 │      │
│  │  4. Format: Parquet                                          │      │
│  │  5. Partitioning: service_name/year/month/day/endpoint_name  │      │
│  └──────────────────────────────────────────────────────────────┘      │
│                                                                         │
│  ┌──────────────────────────────────────────────────────────────┐      │
│  │  Connector 2: CA Connector                                   │      │
│  │                                                              │      │
│  │  1. Consume from: audit-logs-ca                              │      │
│  │  2. SMT Filter: AuditLogSinkCAFilter                         │      │
│  │     • Accept ONLY if wm-site-id = "CA"                       │      │
│  │     • Strict filter (compliance: PIPEDA)                     │      │
│  │  3. Destination: gs://walmart-audit-logs-ca/                 │      │
│  │  4. Format: Parquet                                          │      │
│  │  5. Partitioning: service_name/year/month/day/endpoint_name  │      │
│  └──────────────────────────────────────────────────────────────┘      │
│                                                                         │
│  ┌──────────────────────────────────────────────────────────────┐      │
│  │  Connector 3: MX Connector                                   │      │
│  │                                                              │      │
│  │  1. Consume from: audit-logs-mx                              │      │
│  │  2. SMT Filter: AuditLogSinkMXFilter                         │      │
│  │     • Accept ONLY if wm-site-id = "MX"                       │      │
│  │     • Strict filter (compliance: LFPDPPP)                    │      │
│  │  3. Destination: gs://walmart-audit-logs-mx/                 │      │
│  │  4. Format: Parquet                                          │      │
│  │  5. Partitioning: service_name/year/month/day/endpoint_name  │      │
│  └──────────────────────────────────────────────────────────────┘      │
│                                                                         │
│  Connector Configuration (Lenses GCS Connector):                       │
│    • flush.size: 5000 records                                          │
│    • flush.interval: 10000ms (10 seconds)                              │
│    • rotate.schedule.interval.ms: 3600000 (1 hour)                     │
└────────────────────────────┬────────────────────────────────────────────┘
                             │
                             │ GCS Upload
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                Step 5: GCS Storage (Parquet Files)                      │
│                                                                         │
│  GCS Bucket Structure (US):                                             │
│  gs://walmart-audit-logs-us/                                            │
│    ├── inventory-status-srv/                                            │
│    │   ├── 2025/                                                        │
│    │   │   ├── 03/                                                      │
│    │   │   │   ├── 15/                                                  │
│    │   │   │   │   ├── search-items/                                    │
│    │   │   │   │   │   ├── part-00000.parquet (5000 records)           │
│    │   │   │   │   │   ├── part-00001.parquet (5000 records)           │
│    │   │   │   │   │   └── part-00002.parquet (3245 records)           │
│                                                                         │
│  Parquet File Schema:                                                   │
│    • request_id: STRING                                                │
│    • service_name: STRING                                              │
│    • endpoint_name: STRING                                             │
│    • method: STRING                                                    │
│    • request_body: STRING                                              │
│    • response_body: STRING                                             │
│    • response_code: INT                                                │
│    • request_ts: LONG                                                  │
│    • response_ts: LONG                                                 │
│    • headers: MAP<STRING, STRING>                                      │
│                                                                         │
│  Benefits of Parquet:                                                   │
│    • Columnar storage (efficient queries)                              │
│    • Compression (10x smaller than JSON)                               │
│    • Schema evolution support                                          │
└────────────────────────────┬────────────────────────────────────────────┘
                             │
                             │ BigQuery External Table
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                   Step 6: BigQuery Analytics                            │
│                                                                         │
│  BigQuery External Tables:                                              │
│    • audit_logs_us (external table over GCS US bucket)                 │
│    • audit_logs_ca (external table over GCS CA bucket)                 │
│    • audit_logs_mx (external table over GCS MX bucket)                 │
│                                                                         │
│  Sample Query:                                                          │
│  SELECT                                                                 │
│      service_name,                                                      │
│      endpoint_name,                                                     │
│      COUNT(*) as request_count,                                         │
│      AVG(response_ts - request_ts) as avg_latency_ms                    │
│  FROM audit_logs_us                                                     │
│  WHERE DATE(request_ts) = '2025-03-15'                                  │
│  GROUP BY service_name, endpoint_name                                   │
│  ORDER BY request_count DESC;                                           │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## Custom SMT (Single Message Transform) Filters

```java
// US Filter (Permissive - accepts missing site-id)
public class AuditLogSinkUSFilter extends BaseAuditLogSinkFilter {
    @Override
    public boolean verifyHeader(R r) {
        // Accept if has US site-id OR no site-id header
        boolean hasUsSiteId = hasHeader(r, "wm-site-id", usConfig.getSiteIdValue());
        boolean noSiteId = !hasHeader(r, "wm-site-id");
        return hasUsSiteId || noSiteId;
    }
}

// CA Filter (Strict - only CA site-id)
public class AuditLogSinkCAFilter extends BaseAuditLogSinkFilter {
    @Override
    public boolean verifyHeader(R r) {
        // Accept ONLY if has CA site-id
        return hasHeader(r, "wm-site-id", caConfig.getSiteIdValue());
    }
}

// MX Filter (Strict - only MX site-id)
public class AuditLogSinkMXFilter extends BaseAuditLogSinkFilter {
    @Override
    public boolean verifyHeader(R r) {
        // Accept ONLY if has MX site-id
        return hasHeader(r, "wm-site-id", mxConfig.getSiteIdValue());
    }
}
```

**Why Permissive US Filter?**:
- Legacy systems may not send site-id header
- Default to US market (largest volume)
- Prevents data loss

**Why Strict CA/MX Filters?**:
- Compliance requirements (PIPEDA, LFPDPPP)
- Canadian/Mexican data must NOT go to US bucket
- Explicit site-id required

---

[Continue with remaining sections...]

**END OF COMPREHENSIVE WALMART SYSTEM ARCHITECTURE**

This document provides complete system diagrams and architecture details for all Walmart services. Use this as your technical reference for system design interviews.

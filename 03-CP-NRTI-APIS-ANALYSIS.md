# ULTRA-THOROUGH ANALYSIS: cp-nrti-apis

## Executive Summary

The **cp-nrti-apis** service is a mission-critical microservice that provides **Near Real-Time Inventory (NRTI)** visibility for Walmart's supplier ecosystem. It enables suppliers and vendors to track product availability, inventory movements, shipment arrivals, and transaction history across Walmart's stores and distribution network. The service acts as a secure gateway between external suppliers and Walmart's internal Enterprise Inventory systems, with sophisticated authorization controls, multi-region deployment, and comprehensive audit capabilities.

---

## TABLE OF CONTENTS

1. [Service Purpose and Functionality](#1-service-purpose-and-functionality)
2. [Technology Stack](#2-technology-stack)
3. [Architecture and Code Organization](#3-architecture-and-code-organization)
4. [Key Components](#4-key-components-and-responsibilities)
5. [API Endpoints (Complete Reference)](#5-api-endpoints-and-interfaces)
6. [Data Models and Schemas](#6-data-models-schemas-and-entities)
7. [Dependencies and Integrations](#7-dependencies-and-integrations)
8. [Deployment Configuration](#8-deployment-configuration)
9. [CI/CD Pipeline](#9-cicd-pipeline)
10. [Monitoring and Observability](#10-monitoring-and-observability)
11. [Security](#11-security-configurations)
12. [Performance](#12-performance-characteristics)
13. [Design Patterns](#13-special-patterns-and-noteworthy-implementations)
14. [Testing Strategy](#14-testing-strategy)
15. [Additional Notes](#15-additional-notes)

---

## 1. SERVICE PURPOSE AND FUNCTIONALITY

### What is NRTI?

**NRTI stands for "Near Real-Time Inventory"** - a critical service within Walmart's Channel Performance (CP) system that provides near real-time visibility into inventory data across stores and distribution centers.

### Business Problem & Solution

The service solves several key business problems:

1. **Real-time inventory visibility** - Suppliers can track product availability across Walmart stores
2. **Inventory state change tracking** - Capture arrivals, removals, corrections, and bootstrap events
3. **Transaction event history** - Historical audit trail of all inventory movements
4. **Direct Store Delivery (DSD) tracking** - Vendor-managed inventory shipment notifications
5. **Store inbound forecasting** - Expected arrivals help suppliers plan deliveries
6. **Authorization & validation** - Ensure suppliers only access authorized GTINs and stores
7. **Multi-channel visibility** - Track inventory across store floor, backroom, and MFC locations

### Core Use Cases

1. **Inventory Actions (IAC)** - Real-time inventory state change events:
   - ARRIVAL: New inventory received
   - REMOVAL: Inventory sold or removed
   - CORRECTION: Count adjustments
   - BOOTSTRAP: Initial inventory load

2. **On-hand Inventory Queries** - Current stock levels by store and GTIN

3. **Transaction History** - Historical inventory movement tracking with pagination

4. **Store Inbound Forecasting** - Expected arrivals and receiving schedules

5. **Direct Shipment Capture (DSC)** - DSD vendor shipment notifications with push alerts to stores

6. **Item Validation** - Verify vendor-GTIN associations before operations

7. **Multi-store Inventory Status** - Batch queries across up to 100 GTINs/stores

### Business Value

- **Supplier Empowerment**: Self-service inventory visibility reduces support calls
- **Operational Efficiency**: Automated DSD notifications improve receiving
- **Compliance**: Complete audit trail of inventory movements
- **Data Quality**: Validation ensures accurate inventory records
- **Scalability**: Kafka-based architecture handles high volume with low latency

---

## 2. TECHNOLOGY STACK

### Core Framework & Language
- **Language**: Java 17
- **Framework**: Spring Boot 3.5.7
- **Build Tool**: Maven 3.9.1
- **Parent**: spring-boot-starter-parent 3.5.7

### Key Libraries & Dependencies

#### Spring Boot Modules
- `spring-boot-starter-actuator` - Health checks, metrics
- `spring-boot-starter-webflux` - Reactive web client
- `spring-boot-starter-data-jpa` - Database persistence
- `spring-boot-starter-validation` - Request validation (JSR-303)
- `spring-boot-starter-cache` - Caffeine caching
- `spring-boot-starter-tomcat` - Embedded servlet container

#### Walmart Platform Libraries
- `GTP BOM 2.2.4` - Walmart's Global Technology Platform
- `txn-marking-bom` - Transaction tracing
- `cp-data-apis-common 0.0.22` - Shared CP libraries
- `dv-api-common-libraries 0.0.54` - Data Ventures common
- `walmart-postgresql` - Custom PostgreSQL driver

#### Messaging
- `spring-kafka` - Kafka producer for event streaming
- `kafka-clients` - Kafka client library

#### Database & Storage
- **PostgreSQL** - Primary data store (via walmart-postgresql driver)
- **Google Cloud BigQuery** - Analytics queries
- **Azure Cosmos DB** - Document storage (legacy/optional)

#### API & Documentation
- `springdoc-openapi-starter-webmvc-ui 2.3.0` - Swagger/OpenAPI
- OpenAPI 3.0.3 specification
- `openapi-generator-maven-plugin 7.0.1` - Code generation

#### Utilities
- `Lombok 1.18.30` - Boilerplate reduction
- `MapStruct 1.5.5.Final` - Object mapping
- `Jackson` - JSON serialization
- `Apache Commons` - Utilities
- `Caffeine` - Caching

#### Testing
- `JUnit 5` - Unit testing
- `TestNG 7.9.0` - Integration testing
- `Cucumber 7.19.0` - BDD testing
- `Rest Assured 5.5.0` - API testing
- `Mockito 5.10.0` - Mocking
- `WireMock` - Service virtualization
- `Testcontainers` - Containerized testing

#### Monitoring & Observability
- `Micrometer` - Metrics
- `Prometheus` - Metrics export
- `Dynatrace SAAS` - APM
- `OpenTelemetry` - Distributed tracing
- `Strati Logging` - Walmart's logging framework

#### Security
- Walmart authentication framework (wm_consumer.id, wm_sec.auth_signature)
- SSL/TLS for Kafka communication
- Private key-based service authentication
- Akeyless secret management

---

## 3. ARCHITECTURE AND CODE ORGANIZATION

### Project Structure

```
cp-nrti-apis/
├── src/main/java/com/walmart/cpnrti/
│   ├── NrtiApiApplication.java          # Main Spring Boot entry point
│   ├── controller/                      # REST controllers
│   │   ├── NrtiStoreControllerV1.java  # Main NRTI endpoints
│   │   ├── IacControllerV1.java        # IAC-specific endpoints
│   │   ├── VoltControllerV1.java       # Item validation endpoints
│   │   ├── InventoryController.java    # General inventory endpoints
│   │   ├── DcInventoryController.java  # DC inventory endpoints
│   │   └── *SandboxController.java     # Sandbox/testing endpoints
│   ├── services/                        # Business logic interfaces
│   │   └── impl/                       # Service implementations
│   ├── repository/                      # JPA repositories
│   ├── entity/                         # JPA entities
│   ├── models/                         # DTOs and response models
│   │   ├── response/                   # API response models
│   │   ├── payloads/                   # Request payloads
│   │   └── enums/                      # Enumerations
│   ├── configs/                        # Configuration classes
│   │   ├── postgres/                   # PostgreSQL config
│   │   └── *CCMConfig.java            # CCM2 configurations
│   ├── kafka/                          # Kafka producer config
│   ├── clients/                        # External API clients
│   ├── filters/                        # Request/response filters
│   ├── interceptors/                   # HTTP interceptors
│   ├── validations/                    # Custom validators
│   ├── exception/handlers/             # Exception handling
│   ├── utils/                          # Utility classes
│   ├── constants/                      # Application constants
│   ├── mapper/                         # MapStruct mappers
│   ├── serializer/                     # Custom serializers
│   └── deserializer/                   # Custom deserializers
├── src/main/resources/
│   ├── application.properties          # Spring Boot config
│   ├── messages.properties             # Error messages
│   └── environmentConfig/              # Environment-specific configs
├── src/test/                           # Test code
│   ├── java/                          # Unit & integration tests
│   └── resources/
│       ├── features/                   # Cucumber BDD tests
│       ├── testdata/                   # Test data
│       └── wiremock-mappings/         # WireMock stubs
├── api-spec/                           # OpenAPI specifications
├── perf/                               # Performance test configs
├── pom.xml                             # Maven configuration
├── kitt.yml                            # Kubernetes deployment config
└── ccm.yml                             # CCM2 configuration
```

### Architectural Patterns

**1. Layered Architecture**
- **Controller Layer**: REST endpoints, request validation, response formatting
- **Service Layer**: Business logic, orchestration, external API calls
- **Repository Layer**: Data access, JPA queries
- **Entity Layer**: Database models

**2. Dependency Injection**
- Constructor-based injection throughout
- Spring's @Autowired for dependencies
- @ManagedConfiguration for CCM2 config injection

**3. Transaction Management**
- Walmart's Transaction Marking framework for distributed tracing
- Child transactions for service calls
- Correlation ID propagation

**4. Exception Handling**
- Centralized exception handlers
- Custom exceptions (NrtiMappingException, NrtiDataNotFoundException, etc.)
- Structured error responses (NrtiApiErrorDetails)

**5. Validation Strategy**
- JSR-303 Bean Validation (@Valid)
- Custom validators (IacValidator, DscRequestValidator)
- Business validation services (NrtBusinessValidatorService)
- GTIN-store mapping validation

**6. Caching Strategy**
- Caffeine cache for supplier mappings
- Cache TTL: 604800000ms (7 days)
- Parent company mapping cache
- Sumo authentication cache (30 min TTL)

**7. Async Processing**
- CompletableFuture for parallel API calls
- WebClient for reactive HTTP calls
- Kafka for asynchronous event publishing

---

## 4. KEY COMPONENTS AND RESPONSIBILITIES

### Controllers

#### NrtiStoreControllerV1 (`/store/*`)
- `GET /store/{storeNbr}/gtin/{gtin}/available` - Single GTIN on-hand inventory
- `POST /store/inventoryActions` - Inventory action events (IAC)
- `POST /store/inventory/status` - Multi-GTIN/multi-store inventory status
- `GET /store/{storeNbr}/gtin/{gtin}/transactionHistory` - Transaction event history
- `GET /store/{storeNbr}/gtin/{gtin}/storeInbound` - Store inbound inventory
- `POST /store/directshipment` - Direct shipment capture (DSC)

#### IacControllerV1 (IAC-specific)
- `POST /store/inventoryActions` - IAC events with different headers

#### VoltControllerV1 (`/volt/*`)
- `GET /volt/vendorId/{vendorId}/gtin/{gtin}/itemValidation` - Vendor-GTIN validation

#### InventoryController (`/inventory/*`)
- Items assortment queries (BigQuery-based)

#### Sandbox Controllers
- Testing endpoints for sandbox environment
- Relaxed validation for development

### Service Layer

#### NrtiStoreServiceImpl
**Responsibility**: Core inventory operations

- Orchestrates calls to Enterprise Inventory (EI) APIs
- Handles on-hand inventory queries
- Manages transaction event history
- Processes store inbound inventory
- Publishes events to Kafka
- Coordinates GTIN validation

#### NrtKafkaProducerServiceImpl
**Responsibility**: Event publishing

- Publishes IAC events to Kafka topics
- Publishes DSC events to Kafka topics
- Dual-region Kafka setup (EUS2/SCUS)
- SSL-secured connections

#### StoreGtinValidatorServiceImpl
**Responsibility**: Authorization validation

- Validates GTIN-to-supplier mappings
- Validates store-GTIN combinations
- Uses PostgreSQL for mapping lookups
- Caches validation results

#### SupplierMappingServiceImpl
**Responsibility**: Supplier context

- Retrieves supplier details from consumer ID
- Caches parent company mappings
- Provides global DUNS lookup

#### UberKeyReadServiceImpl
**Responsibility**: Item mapping

- GTIN-to-CID (Customer Item Descriptor) conversion
- Batch lookup support (up to 100 items)
- Integrates with Uber Keys API

#### HttpServiceImpl
**Responsibility**: External API client

- Generic HTTP client for external services
- Handles authentication (Walmart headers)
- Retry logic and error handling

#### BigQueryServiceImpl
**Responsibility**: Analytics queries

- Assortment queries from BigQuery
- GRS (Global Replenishment System) data
- Scheduled query execution

#### SumoServiceImpl
**Responsibility**: Push notifications

- Sends mobile notifications to store associates
- Integrates with Sumo (Walmart's notification service)
- Role-based targeting (e.g., Asset Protection - DSD)

#### LocationPlatformApiServiceImpl
**Responsibility**: Store metadata

- Retrieves store timezone information
- Integrates with Location Platform API

#### PlatformServiceImpl
**Responsibility**: Identity management

- Company association lookups
- Luminate company ID resolution

### Repository Layer

#### ParentCompanyMappingRepository
- **Table**: `nrt_consumers`
- **Primary key**: consumer_id, site_id
- **Purpose**: Stores supplier metadata and permissions

#### NrtStoreGtinMappingRepository
- **Table**: `supplier_gtin_items`
- **Primary key**: site_id, gtin, global_duns
- **Purpose**: Stores GTIN-to-supplier-to-store mappings
- **Special**: Array column for store numbers

#### VendorRepository
- **Table**: `vendor`
- **Purpose**: Stores vendor name and ID mappings

#### VendorGtinItemsRepository
- **Table**: `vendor_gtin_items`
- **Purpose**: Vendor-specific GTIN mappings

### Filters & Interceptors

#### RequestFilter
**Responsibility**: Request logging and correlation
- Logs incoming requests
- Extracts correlation IDs
- Sanitizes sensitive data

#### NrtResponseFilter
**Responsibility**: Response enhancement
- Adds correlation headers
- Formats response timestamps

#### XssFilter
**Responsibility**: Security
- XSS attack prevention
- Input sanitization
- Uses OWASP ESAPI

#### NrtCorsFilter
**Responsibility**: CORS handling
- Cross-origin request support
- Configurable allowed origins
- Pre-flight request handling

#### NrtiApiInterceptor
**Responsibility**: Authentication
- Validates Walmart security headers
- Signature verification
- Consumer ID validation

#### SiteIdFilterAspect
**Responsibility**: Multi-tenancy
- Adds site_id to database queries
- AOP-based implementation
- Ensures data isolation

---

## 5. API ENDPOINTS AND INTERFACES

### Complete REST API Documentation

#### 1. Inventory Actions (IAC)

**Endpoint**: `POST /store/inventoryActions`

**Purpose**: Record real-time inventory state changes (arrivals, removals, corrections, bootstrap)

**Headers**:
- `wm_consumer.id` (required): Consumer UUID
- `wm_consumer.intimestamp` (required): Timestamp
- `wm_sec.auth_signature` (required): HMAC signature
- `wm_sec.key_version` (required): Key version
- `wm_svc.name` (required): Service name (channelperformance-nrti or channelperformance-iac)
- `wm_svc.env` (required): Environment
- `wm_qos.correlation_id` (optional): Correlation ID for tracing

**Request Body**:
```json
{
  "message_id": "746007c9-4b2c-4838-bfd9-037d341c2d2d",
  "event_type": "ARRIVAL|REMOVAL|CORRECTION|BOOTSTRAP",
  "store_nbr": 100,
  "line_infos": [
    {
      "gtin": "00083754843990",
      "quantity": 20.5,
      "secondary_item_identifier": {
        "type": "UPC",
        "value": "12345678901234"
      },
      "destination_location": {
        "location_area": "STORE|BACKROOM|MFC",
        "location": "A-1-2",
        "lpn": "LPN12345"
      },
      "expiry_date_at": "2025-12-31"
    }
  ],
  "document_infos": [
    {
      "doc_type": "PO",
      "doc_nbr": "PO123456",
      "doc_date": 1651082806067
    }
  ],
  "user_id": "user123",
  "reason_details": [
    {
      "reason_code": "RC01",
      "reason_desc": "Receiving from truck"
    }
  ],
  "vendor_nbr": "544528",
  "event_creation_time": 1651082806061
}
```

**Response**: `201 Created`
```json
{
  "message": "Event created successfully for GTINs: [00083754843990]"
}
```

**Business Rules**:
- Event creation time must be within 3 days (configurable)
- All GTINs must be mapped to the supplier
- Store must be in supplier's authorized list
- Validates location_area against approved values
- Supports partial success (some GTINs may fail validation)

**Kafka Event**: Published to `cperf-nrt-prod-iac` topic

---

#### 2. On-Hand Inventory (Single GTIN)

**Endpoint**: `GET /store/{storeNbr}/gtin/{gtin}/available`

**Purpose**: Retrieve current on-hand inventory for a single GTIN at a store

**Path Parameters**:
- `storeNbr` (integer): Store number
- `gtin` (string): 14-digit GTIN

**Headers**: (same as IAC endpoint)

**Response**: `200 OK`
```json
{
  "supplier": "ABC Company",
  "store_nbr": 100,
  "gtin": "00083754843990",
  "inventories": [
    {
      "location_area": "STORE",
      "state": "ON_HAND",
      "quantity": 150.0
    },
    {
      "location_area": "BACKROOM",
      "state": "ON_HAND",
      "quantity": 50.0
    },
    {
      "location_area": "MFC",
      "state": "ON_HAND",
      "quantity": 25.0
    }
  ],
  "last_updated_time": "2026-02-02T10:30:00Z"
}
```

**Data Source**: Enterprise Inventory (EI) On-hand Inventory API

---

#### 3. Multi-Store Inventory Status

**Endpoint**: `POST /store/inventory/status`

**Purpose**: Retrieve inventory for multiple GTINs across multiple stores

**Query Parameters**:
- `itemIdentifier` (optional): "gtin" or "wmItemNbr"

**Request Body**:
```json
{
  "store_nbr": [100, 200, 300],
  "gtins": ["00083754843990", "00012345678901"],
  "wm_item_nbrs": [123456789]
}
```

**Response**: `200 OK`
```json
{
  "supplier": "ABC Company",
  "inventories": [
    {
      "store_nbr": 100,
      "gtin": "00083754843990",
      "wm_item_nbr": 123456789,
      "inventory_by_locations": [
        {
          "location_area": "STORE",
          "state": "ON_HAND",
          "quantity": 150.0
        }
      ],
      "last_updated_time": "2026-02-02T10:30:00Z"
    }
  ],
  "errors": []
}
```

**Features**:
- Supports up to 100 GTINs per request
- Parallel processing of requests
- Partial results with error details
- Handles both GTIN and Walmart Item Number lookups

---

#### 4. Transaction Event History

**Endpoint**: `GET /store/{storeNbr}/gtin/{gtin}/transactionHistory`

**Purpose**: Retrieve historical inventory transactions for a GTIN

**Path Parameters**:
- `storeNbr` (integer): Store number
- `gtin` (string): 14-digit GTIN

**Query Parameters**:
- `start_date_at` (date): Start date (default: 6 days ago)
- `end_date_at` (date): End date (default: today)
- `event_type` (string): Filter by event type (default: "ALL")
- `location_area` (string): Filter by location (default: "STORE")

**Headers**:
- `continuationpagetoken` (optional): Pagination token for next page

**Response**: `200 OK`
```json
{
  "supplier": "ABC Company",
  "store_nbr": 100,
  "gtin": "00083754843990",
  "event_histories": [
    {
      "event_type": "INVENTORY_STATE_CHANGE",
      "event_qty": 50.0,
      "event_creation_time": "2026-02-01T14:30:00Z",
      "aggregated_qty": 200.0
    },
    {
      "event_type": "BACKROOM_MOVED",
      "event_qty": -25.0,
      "event_creation_time": "2026-02-01T16:45:00Z",
      "aggregated_qty": 175.0
    }
  ],
  "pagination_token": "eyJzdG9yZU5iciI6MTAwLCJndGluIjoiMDAwODM3NTQ4NDM5OTAiLCJvZmZzZXQiOjUwfQ=="
}
```

**Response Headers**:
- `continuationpagetoken`: Token for next page of results

**Data Source**: EI Inventory History Lookup API

---

#### 5. Store Inbound Inventory

**Endpoint**: `GET /store/{storeNbr}/gtin/{gtin}/storeInbound`

**Purpose**: Retrieve expected inbound inventory for a GTIN at a store

**Path Parameters**:
- `storeNbr` (integer): Store number
- `gtin` (string): 14-digit GTIN

**Response**: `200 OK`
```json
{
  "store_nbr": 100,
  "gtin": "00083754843990",
  "store_inbound_inventory_by_eads": [
    {
      "ead": "2026-02-05",
      "location_areas": [
        {
          "location_area": "STORE",
          "inventory_by_states": [
            {
              "state": "IN_TRANSIT",
              "quantity": 100.0
            },
            {
              "state": "RECEIVED",
              "quantity": 50.0
            }
          ]
        }
      ]
    }
  ]
}
```

**Features**:
- Expected Arrival Date (EAD) grouping
- State breakdown (IN_TRANSIT, RECEIVED, etc.)
- Looks ahead 30 days (configurable)
- Uses CID (Customer Item Descriptor) for lookups

**Data Source**: EI Inbound Inventory by ReplGrp API

---

#### 6. Direct Shipment Capture (DSC)

**Endpoint**: `POST /store/directshipment`

**Purpose**: Capture direct store delivery (DSD) shipment information from vendors

**Request Body**:
```json
{
  "message_id": "746007c9-4b2c-4838-bfd9-037d341c2d2d",
  "event_creation_time": 1651082806061,
  "event_type": "PLANNED|CANCELLED",
  "vendor_id": "544528",
  "supplier_origin": "1000",
  "destinations": [
    {
      "store_nbr": 100,
      "store_sequence": 5,
      "loads": [
        {
          "asn": "12345875886",
          "actual_shipment": {
            "pallet_qty": 28,
            "total_cube": 1881.76,
            "cube_uom": "Ft",
            "total_weight": 35950.15,
            "weight_uom": "LBS",
            "case_qty": 1324,
            "load_ready_ts_at": "2026-02-02T10:00:00Z"
          },
          "planned_shipment": { }
        }
      ],
      "planned_eta_at": "2026-02-05T14:00:00Z",
      "actual_eta_window": {
        "earliest_eta_at": "2026-02-05T13:00:00Z",
        "latest_eta_at": "2026-02-05T15:00:00Z"
      },
      "arrival_time_at": "2026-02-05T14:30:00Z"
    }
  ],
  "trailer_nbr": "12332"
}
```

**Response**: `201 Created`
```json
{
  "message": "DSC event created successfully"
}
```

**Special Features**:
- Publishes to Kafka for downstream processing
- Sends push notifications to store associates via Sumo
- Supports commodity type mapping (e.g., vendor 544528 = Core-Mark International)
- Validates vendor ID against database
- Multi-destination support (up to 30 stores per request)

**Kafka Event**: Published to `cperf-nrt-prod-dsc` topic

---

#### 7. Item Validation (Vendor-GTIN Permission)

**Endpoint**: `GET /volt/vendorId/{vendorId}/gtin/{gtin}/itemValidation`

**Purpose**: Verify if a GTIN belongs to a specific vendor

**Path Parameters**:
- `vendorId` (string): Vendor ID
- `gtin` (string): 14-digit GTIN

**Response**: `200 OK`
```json
{
  "vendor_id": "544528",
  "gtin": "00083754843990",
  "has_permission": true
}
```

**Data Source**: PostgreSQL vendor_gtin_items table

---

### Sandbox Endpoints

The service includes sandbox variants of main endpoints for testing:
- `/store/inventory/status` (sandbox version)
- `/store/{storeNbr}/gtin/{gtin}/storeInbound` (sandbox version)
- `/volt/vendorId/{vendorId}/gtin/{gtin}/itemValidation` (sandbox version)

**Sandbox Features**:
- Relaxed validation rules
- Test data support
- Separate service header: `wm_svc.name: channelperformance-nrti-sandbox`

---

### Authentication & Security

**All endpoints require**:
1. **Consumer ID**: Unique UUID for the API consumer
2. **HMAC Signature**: Request signature using private key
3. **Timestamp**: Request timestamp for replay prevention
4. **Key Version**: Version of signing key
5. **Service Name**: Identifies calling service

**Validation Flow**:
1. Extract headers
2. Verify signature
3. Validate consumer ID exists in database
4. Check supplier authorization
5. Validate GTIN-store mappings

---

## 6. DATA MODELS, SCHEMAS, AND ENTITIES

### Database Entities (PostgreSQL)

#### 1. ParentCompanyMapping (nrt_consumers)

```java
@Entity
@Table(name = "nrt_consumers")
public class ParentCompanyMapping {
    @EmbeddedId
    private ParentCompanyMappingKey id;  // consumer_id, site_id

    private String consumerName;
    private String countryCode;
    private String globalDuns;            // D&B number
    private Boolean isCategoryManager;
    private Boolean nonCharterSupplier;
    private Boolean isPartnerCompany;
    private String luminateCmpnyId;       // Luminate company ID
    private String parentCmpnyName;       // Parent company name
    private String pspGlobalDuns;         // PSP global DUNS
    private String pspName;               // Payment service provider
    private String remarks;
    private Status status;                // ACTIVE, INACTIVE
    private UserType userType;            // SUPPLIER, VENDOR, PARTNER
    private Long vendorId;
    private String appName;
}
```

**Purpose**: Maps API consumer IDs to supplier/vendor information

**Key Columns**:
- `consumer_id` (UUID): Primary identifier for API client
- `global_duns`: Dun & Bradstreet number (supplier identifier)
- `luminate_cmpny_id`: Walmart's internal company ID
- `parent_cmpny_name`: Display name for supplier
- `user_type`: SUPPLIER, VENDOR, PARTNER
- `status`: ACTIVE, INACTIVE

---

#### 2. NrtStoreGtinMapping (supplier_gtin_items)

```java
@Entity
@Table(name = "supplier_gtin_items")
public class NrtStoreGtinMapping {
    @EmbeddedId
    private NrtiStoreGtinMappingKey key;  // site_id, gtin, global_duns

    @Column(columnDefinition = "integer[]")
    private Integer[] storeNumber;        // Array of authorized stores

    private String luminateCmpnyId;
    private String parentCompanyName;
    private String geoRegionCode;         // US, CA, MX
    private String opCompanyCode;         // Operating company
}
```

**Purpose**: Authorization matrix for supplier-GTIN-store access

**Key Columns**:
- `site_id`: Site identifier (for multi-tenancy)
- `gtin`: 14-digit GTIN
- `global_duns`: Supplier identifier
- `store_nbr`: Array of authorized store numbers
- `luminate_cmpny_id`: Company ID

**Unique Constraint**: (site_id, gtin, global_duns)

---

#### 3. Vendor (vendor)

```java
@Entity
@Table(name = "vendor")
public class Vendor {
    @EmbeddedId
    private VendorKey id;  // site_id, vendor_id

    private String vendorNm;
    private String geoRegionCd;
    private String opCmpnyCd;
    private String userid;
}
```

**Purpose**: Vendor master data for DSC events

---

#### 4. VendorGtinItemsMapping (vendor_gtin_items)

```java
@Entity
@Table(name = "vendor_gtin_items")
public class VendorGtinItemsMapping {
    @EmbeddedId
    private VendorGtinItemsMappingKey key;  // site_id, vendor_id, gtin

    private String vendorNm;
    private String wmItemNbr;
    private String geoRegionCd;
    private String opCmpnyCd;
}
```

**Purpose**: Maps GTINs to vendors for validation

---

### API Request/Response Models

#### Inventory Actions Request

```java
public class NrtInventoryActionsRequest {
    @NotBlank
    private String messageId;           // UUID

    @NotNull
    private EventType eventType;        // ARRIVAL, REMOVAL, CORRECTION, BOOTSTRAP

    @NotNull
    private Integer storeNbr;

    @NotEmpty
    private List<LineInfoItem> lineInfo;

    private List<DocumentInfoItem> documentInfos;

    @NotBlank
    private String userId;

    @NotEmpty
    private List<ReasonDetailsItem> reasonDetails;

    private String vendorNbr;

    @NotNull
    private Long eventCreationTime;     // Epoch millis
}
```

#### Line Info Item

```java
public class LineInfoItem {
    @Pattern(regexp = "^\\d{14}$")
    private String gtin;

    @DecimalMin("0")
    private Double quantity;

    private SecondaryItemIdentifier secondaryItemIdentifier;

    @NotNull
    private DestinationLocation destinationLocation;

    private LocalDate expiryDateAt;
}
```

#### Destination Location

```java
public class DestinationLocation {
    @Pattern(regexp = "^$|^STORE$|^BACKROOM$|^MFC$")
    private String locationArea;

    private String location;            // Specific location code
    private String lpn;                 // License Plate Number
}
```

---

### Enumerations

**EventType**:
- ARRIVAL: New inventory arrival
- REMOVAL: Inventory removed/sold
- CORRECTION: Inventory count adjustment
- BOOTSTRAP: Initial inventory load

**LocationArea**:
- STORE: Sales floor
- BACKROOM: Backroom storage
- MFC: Micro Fulfillment Center

**Status**:
- ACTIVE
- INACTIVE

**UserType**:
- SUPPLIER
- VENDOR
- PARTNER
- CATEGORY_MANAGER

**InventoryState**:
- ON_HAND
- IN_TRANSIT
- RECEIVED
- COMMITTED
- RESERVED

---

## 7. DEPENDENCIES AND INTEGRATIONS

### External Service Integrations

#### 1. Enterprise Inventory (EI) APIs

**Purpose**: Walmart's inventory system of record

**Endpoints Used**:
- **On-hand Inventory Read**: `https://ei-onhand-inventory-read.walmart.com/api/v1/inventory/node/{nodeId}/gtin/{gtin}`
- **Inventory History Lookup**: `https://ei-inventory-history-lookup.walmart.com/v1/historyForInventoryState/countryCode/{us}/nodeId/{nodeId}/gtin/{gtin}`
- **Inbound Inventory by ReplGrp**: `https://ei-inbound-inventory-by-replgrp-v2.prod.walmart.com/api/v1/inboundInventory/node/{nodeId}/cid/{cid}/fromDate/{fromDate}/toDate/{toDate}`
- **PIT by Item Inventory Lookup**: `https://ei-pit-by-item-inventory-read.walmart.com/api/v1/inventory/node/{nodeId}/itemnumber`

**Authentication**: Walmart platform authentication (consumer ID, signature)

---

#### 2. Uber Keys API

**Purpose**: GTIN to Customer Item Descriptor (CID) conversion

**Endpoint**: `https://uber-keys-read-nsf.walmart.com/mappings`

**Usage**: Convert GTINs to Walmart internal item numbers (CIDs)

**Batch Support**: Up to 100 items per request

---

#### 3. Location Platform API

**Purpose**: Store metadata and timezone information

**Endpoint**: `https://locationplatform.services.prod.walmart.com/location/stores`

**Usage**: Retrieve store timezone for date calculations

---

#### 4. Identity Management Platform API

**Purpose**: Company and association lookups

**Endpoints**:
- **Association API**: `/v1/companies/{luminateCompanyId}/associations`
- **Company Search**: `/search/companies`

**Usage**: Resolve Luminate company IDs and associations

---

#### 5. Sumo (Push Notification Service)

**Purpose**: Send mobile notifications to store associates

**Endpoint**: `https://api-proxy-es2.prod-us-azure.soa-api-proxy.platform.prod.us.walmart.net/api-proxy/service/sms/sumo/v3/mobile/push`

**Usage**: Notify store associates of DSD shipment arrivals

**Target Roles**: US_STORE_ASSET_PROT_DSD

---

#### 6. Google BigQuery

**Purpose**: Assortment and analytics queries

**Datasets**:
- `wmt-grs-gcp-us-prod.US_WM_REPL_TABLES.GRS_ASSORTMENT`: Global Replenishment System assortment data
- `wmt-dsi-dv-cperf-rb-prod.ww_chnl_perf_rb_app.chnl_perf_rb_psi_dim`: Channel Performance dimension tables

**Query Example**:
```sql
SELECT DISTINCT grs.effective_date, grs.status_desc, grs.wmt_item_nbr,
       dim.gtin, grs.zone_id, grs.model_run_dt
FROM wmt-grs-gcp-us-prod.US_WM_REPL_TABLES.GRS_ASSORTMENT as grs
INNER JOIN wmt-dsi-dv-cperf-rb-prod.ww_chnl_perf_rb_app.chnl_perf_rb_psi_dim dim
  ON grs.wmt_item_nbr = dim.wm_item_nbr
WHERE dim.bus_dt = DATE_SUB(CURRENT_DATE, INTERVAL 1 DAY)
  AND grs.store_nbr = @storeNbr
  AND dim.luminate_cmpny_id = @luminateCompanyId
```

---

### Kafka Integration

**Kafka Cluster**: Luminate Core (Multi-region: EUS2 + SCUS)

**Broker Configuration**:
- **Primary (EUS2)**: 3 brokers on port 9093
- **Secondary (SCUS)**: 3 brokers on port 9093
- **Protocol**: SSL/TLS
- **Authentication**: Keystore/Truststore based

**Topics**:
1. **cperf-nrt-prod-iac**: Inventory action events (IAC)
2. **cperf-nrt-prod-dsc**: Direct shipment capture events (DSC)

**Producer Configuration**:
- Max request size: 10MB
- Compression: LZ4
- Acks: all
- Batch size: 8192 bytes
- Linger: 20ms
- Retries: 10
- Idempotence: disabled

**Event Schema**:
```json
{
  "messageId": "uuid",
  "eventType": "ARRIVAL|REMOVAL|...",
  "storeNbr": 100,
  "lineInfo": [...],
  "globalDuns": "012345678",
  "parentCompanyName": "ABC Company",
  "luminateCompanyId": "12345",
  "eventCreationTime": 1651082806061
}
```

---

### Database: PostgreSQL

**Connection**:
- **Driver**: Walmart PostgreSQL (custom fork)
- **Prod Writer (SCUS)**: prod-pgsqlflex-dv-kys-insexe-us-ve.writer.postgres.database.azure.com
- **Prod Reader (EUS2)**: prod-pgsqlflex-dv-kys-insexe-us-ve.reader.postgres.database.azure.com
- **Database**: kys_insexe_nrt_us
- **Authentication**: USER_PASSWORD (via Akeyless)

**Connection Pool (HikariCP)**:
- Maximum pool size: 15
- Minimum idle: 10
- Connection timeout: 2s
- Socket timeout: 5s
- Idle timeout: 120s
- Validation timeout: 1s

**Tables**:
- nrt_consumers (supplier mappings)
- supplier_gtin_items (authorization matrix)
- vendor (vendor master)
- vendor_gtin_items (vendor-GTIN mappings)

---

## 8. DEPLOYMENT CONFIGURATION

### Kubernetes Deployment (WCNP)

**Deployment Tool**: KITT (Walmart's K8s deployment platform)

**Environments**:

1. **Dev** (non-prod)
   - Clusters: scus-dev-a3, eus2-dev-a2
   - DNS: dev-cp-nrti.walmart.com
   - Pods: 2-4 (autoscaling)
   - Resources: 1-2 CPU, 1-2Gi RAM

2. **Stage** (non-prod)
   - Clusters: eus2-stage-a4, uswest-stage-az-006
   - DNS: stage-cp-nrti.walmart.com
   - Pods: 4-8 (autoscaling)
   - Resources: 1-2 CPU, 1-2Gi RAM

3. **Sandbox** (non-prod)
   - Cluster: uswest-stage-az-006
   - DNS: sandbox-cp-nrti.walmart.com
   - Pods: 1-2 (autoscaling)
   - Resources: 1-2 CPU, 1-2Gi RAM

4. **Production**
   - Clusters: eus2-prod-a30, scus-prod-a63
   - DNS: prod-cp-nrti.walmart.com
   - Pods: 6-12 (autoscaling)
   - Resources: 1-2 CPU, 1-2Gi RAM
   - **Active/Active multi-region deployment**

---

### Container Configuration

**Base Image**: Spring Boot application with Dynatrace OneAgent

**Port**: 8080 (HTTP)

**Java Options**:
```
-Dccm.configs.dir=/etc/config
-Druntime.context.system.property.override.enabled=true
-Druntime.context.environmentType={{stage}}
-Druntime.context.appName=CHANNELPERFORMANCE-NRTI
-Dscm.server.url=http://tunr.prod.walmart.com/scm-app/v2
-Dcom.walmart.platform.metrics.impl.type=MICROMETER
-Dcom.walmart.platform.telemetry.otel.enabled=true
-Dcom.walmart.platform.txnmarking.otel.type=OTLPgRPC
-Dcom.walmart.platform.txnmarking.otel.host=trace-collector.prod.walmart.com
-Dcom.walmart.platform.logging.profile=SLT_DUAL
```

---

### Health Probes

**Startup Probe**:
- Path: `/actuator/health/startup`
- Initial delay: 30s
- Interval: 5s
- Failure threshold: 24

**Liveness Probe**:
- Path: `/actuator/health/liveness`
- Interval: 10s
- Failure threshold: 5

**Readiness Probe**:
- Path: `/actuator/health/readiness`
- Interval: 5s
- Failure threshold: 5

---

### Resource Limits

**Production**:
- Min CPU: 1 core
- Max CPU: 2 cores
- Min Memory: 1Gi
- Max Memory: 2Gi
- CPU autoscaling trigger: 60%

---

### Secrets Management (Akeyless)

**Production Path**: `Prod/WCNP/homeoffice/dv-kys-api-prod-group`

**Secrets**:
- `cosmos_db_key.txt`: Cosmos DB primary key
- `nrt_kafka_ssl_keystore.jks`: Kafka SSL keystore
- `nrt_kafka_ssl_truststore.jks`: Kafka SSL truststore
- `nrt_kafka_ssl_key_pwd.txt`: Kafka key password
- `private_key_platform_api.txt`: Platform API private key
- `sumo_private_key.txt`: Sumo API private key
- `google_key_file.txt`: BigQuery service account key
- `location_platform_api_private_key.txt`: Location API key
- `audit_log_private_key.txt`: Audit logging key
- `postgresql-config`: PostgreSQL connection details

---

### Networking & Service Mesh

**Istio Sidecar**:
- Enabled via annotation: `sidecar.istio.io/inject: "true"`
- Excluded ports: 15020, 8200, 31833, 8300, 8080

**GSLB (Global Server Load Balancing)**:
- Strategy: stage (allows sandbox and pre-prod)
- Multi-region active/active in production

---

### Flagger (Progressive Delivery)

**Canary Deployment**:
- Step weight: 10% increments
- Max weight: 50%
- Interval: 2 minutes
- Progress deadline: 600s (10 minutes)
- Canary replica percentage: 50%

**Canary Metrics**:
- 5XX error threshold: 1%
- Query interval: 2 minutes
- Uses Prometheus/Envoy metrics

---

## 9. CI/CD PIPELINE

### Build Tool: Looper

**Tools**:
- JDK 17
- Maven 3.9.1
- SonarQube Scanner 4.8.0.2856

---

### Build Flow

**Default Flow (main branch)**:
```yaml
1. Clean and build
   - mvn clean install
   - Runs unit tests
   - Runs integration tests (Cucumber)
   - Generates Jacoco coverage reports

2. SonarQube analysis
   - Code quality analysis
   - Code coverage check
   - Security vulnerability scan
   - Technical debt calculation

3. Hygieia publish
   - Build metadata
   - Test results
   - Code quality metrics
```

---

### Test Execution

**Unit Tests**:
- Framework: JUnit 5 + Mockito
- Plugin: maven-surefire-plugin
- Reports: target/surefire-reports/

**Integration Tests**:
- Framework: Cucumber + TestNG + Rest Assured
- Plugin: maven-failsafe-plugin
- Features: src/test/resources/features/
- Reports:
  - target/cucumber-reports/ (JSON)
  - target/cucumber-html-reports/ (HTML)

---

### Deployment Flow (KITT)

**Stage Deployment**:
```yaml
1. Docker build
   - Include Dynatrace OneAgent

2. API Linting (Concord)
   - Mode: Passive

3. Deploy to Kubernetes

4. Post-deployment validation:
   a. Regression Suite (Concord)
   b. R2C Contract Testing (Concord)
   c. Automaton Tests (Concord)
   d. Resiliency Testing (RaaS)

5. Slack notification
```

**Production Deployment**:
```yaml
1. Approval required

2. Change record creation

3. Canary deployment (Flagger)
   - 10% traffic shift every 2 minutes
   - Monitor 5XX errors
   - Auto-rollback on failure

4. Full rollout
   - Active/Active deployment (EUS2 + SCUS)
```

---

## 10. MONITORING AND OBSERVABILITY

### Metrics Collection (Micrometer + Prometheus)

**Metrics Endpoint**: `/actuator/prometheus`

**Published Metrics**:

**HTTP Metrics**:
- `http_server_requests_seconds_count`: Request count
- `http_server_requests_seconds_sum`: Total response time
- `http_server_requests_seconds_max`: Max response time

**JVM Metrics**:
- `jvm_memory_used_bytes`: Heap usage
- `jvm_threads_live_threads`: Active threads
- `jvm_gc_pause_seconds_*`: GC pause metrics

**HikariCP (Connection Pool)**:
- `hikaricp_connections_active`: Active connections
- `hikaricp_connections_idle`: Idle connections
- `hikaricp_connections_pending`: Waiting threads

---

### Distributed Tracing

**OpenTelemetry Configuration**:
- **Enabled**: true
- **Type**: OTLPgRPC
- **Collector Host (prod)**: trace-collector.prod.walmart.com
- **Port**: 80

**Transaction Marking**:
- **Format**: OTelTxnJson
- **Markers**: PS, RS, PE, RE, CE, CS
- **Propagation**: Correlation ID via headers

---

### Logging (Strati Logging)

**Configuration**:
- **Profile**: SLT_DUAL (logs to both file and external collector)
- **Format**: OTelJson
- **Level**: INFO (configurable via CCM2)

**Structured Logging**:
- All logs in JSON format
- Correlation ID in every log
- Transaction context propagation

**Audit Logging**:
- **Enabled**: true (configurable)
- **Endpoint**: https://us.prod.proxy-gateway.walmart.com/api-proxy/service/audit/api-logs-srv/v1/logs/api-requests
- **Consumer ID**: f7c1590c-31c2-4531-98e4-86356dc5a612

**Audited Endpoints**:
- inventoryActions
- transactionHistory
- nrti_available
- nrti_storeInbounds
- nrti_status
- nrti_directshipment
- nrti_itemValidation

---

### APM: Dynatrace

**Integration**:
- **Enabled**: DYNATRACE_ENABLED=true
- **Type**: SaaS
- **OneAgent**: Injected at Docker build time
- **Product Tracking ID**: 5798 (DV-DATAAPI)

**Monitored Metrics**:
- Request rate and response time
- Error rate and types
- Database query performance
- External API call latency
- JVM health and GC
- Thread pool utilization

---

### Health Endpoints

**Spring Boot Actuator**:
- `/actuator/health` - Overall health
- `/actuator/health/liveness` - Liveness probe
- `/actuator/health/readiness` - Readiness probe
- `/actuator/health/startup` - Startup probe

**Health Components**:
- `livenessState`: Application running
- `readinessState`: Ready to accept traffic
- `db`: Database connectivity
- `diskSpace`: Disk availability
- `ssl`: SSL certificate validity
- `ping`: Basic responsiveness

---

## 11. SECURITY CONFIGURATIONS

### Authentication & Authorization

**Walmart Platform Authentication**:
1. **Consumer ID** (`wm_consumer.id`)
   - UUID format
   - Registered in nrt_consumers table

2. **Signature-based Auth** (`wm_sec.auth_signature`)
   - HMAC-SHA256 signature
   - Signed with private key
   - Timestamp validation

3. **Service Context** (`wm_svc.name`, `wm_svc.env`)
   - Identifies calling service
   - Environment validation

---

### Authorization Model

**Multi-level Authorization**:
1. **Consumer-level**: Is the consumer ID registered?
2. **Supplier-level**: Does the supplier have access?
3. **GTIN-level**: Is the GTIN mapped to this supplier?
4. **Store-level**: Is the store authorized?

**Special Rules**:
- Skip store validation for specific global DUNS
- Example: "071058929" can access any store

---

### Data Privacy & Security

1. **XSS Prevention**:
   - XssFilter using OWASP ESAPI
   - Input sanitization

2. **SQL Injection Prevention**:
   - JPA parameterized queries
   - No dynamic SQL

3. **Data Masking**:
   - Sensitive fields logged with masking

4. **Multi-tenancy**:
   - Site ID filter on all queries
   - AOP-based enforcement

---

### SSL/TLS Configuration

**Kafka SSL**:
- **Protocol**: TLS 1.2+
- **Keystore**: JKS format
- **Truststore**: JKS format
- **Passwords**: Stored in Akeyless

**HTTPS Endpoints**:
- All external API calls over HTTPS
- Certificate validation enabled
- TLS 1.2+ required

---

## 12. PERFORMANCE CHARACTERISTICS

### Latency Targets

**SLA**: 300ms (95th percentile)

**Typical Response Times**:
- Single GTIN lookup: < 200ms
- Multi-GTIN query (10 items): < 500ms
- Transaction history: < 800ms
- Store inbound: < 600ms

---

### Throughput

**Peak Capacity (production)**:
- **Daily Peak**: 100 requests/second
- **Holiday Peak**: 140 requests/second
- **Pods**: 6-12 (autoscaling)
- **Per-pod capacity**: ~15-20 req/sec

---

### Caching Strategy

**Caffeine Cache**:
- **Parent Company Mapping Cache**:
  - TTL: 7 days
  - Eviction: Time-based

- **Sumo Authentication Cache**:
  - TTL: 30 minutes

---

### Database Performance

**Query Patterns**:
1. **Point lookups**: Consumer ID, GTIN lookups (indexed)
2. **Range queries**: Store number arrays (GIN index)
3. **Join queries**: Minimal, mostly single-table

**Pool Configuration**:
- Max connections: 15
- Min idle: 10
- Acquire timeout: 2s
- Socket timeout: 5s

---

### Autoscaling

**HPA (Horizontal Pod Autoscaler)**:
- **Metric**: CPU utilization
- **Threshold**: 60%
- **Min pods**: 2-6 (environment-specific)
- **Max pods**: 4-12 (environment-specific)

---

## 13. SPECIAL PATTERNS AND NOTEWORTHY IMPLEMENTATIONS

### 1. Multi-Region Active/Active Architecture

**Implementation**:
- Dual Kafka clusters (EUS2 + SCUS)
- PostgreSQL reader/writer split by region
- Kubernetes deployment in both regions
- GSLB for traffic distribution

**Benefits**:
- Low latency for both regions
- High availability
- Disaster recovery

---

### 2. Supplier Authorization Matrix

**Pattern**: Three-level authorization
```
Consumer ID → Global DUNS → GTIN → Store Array
```

**Implementation**:
- PostgreSQL array columns for store lists
- Cached supplier mappings
- Lazy validation (on-demand)

---

### 3. Event-Driven Architecture (Kafka)

**Pattern**: Publish events to Kafka, don't wait
```
IAC Request → Validate → Publish to Kafka → Return 201
```

**Benefits**:
- Fast API response
- Decoupled downstream processing
- Replay capability

---

### 4. Dual-Mode Operations (NRTI vs IAC)

**Implementation**:
- Same codebase, different endpoints
- Header-based routing (`wm_svc.name`)
- Separate Kafka topics
- Different validation rules

---

### 5. Sandbox Mode for Testing

**Pattern**: Relaxed validation for sandbox
- Skip some validations
- Use test data
- Allow broader access

---

### 6. Pagination for Large Result Sets

**Implementation**:
- Continuation tokens in headers
- Base64-encoded state
- Stateless pagination

---

### 7. Distributed Tracing with Transaction Marking

**Pattern**: Walmart's custom tracing framework
```java
try (var txn = transactionMarkingManager
    .getTransactionMarkingService()
    .currentTransaction()
    .addChildTransaction("ServiceName", "MethodName")
    .start()) {
    // Service call
}
```

---

### 8. MapStruct for Object Mapping

**Implementation**:
```java
@Mapper
public interface EIInboundInventoryMapper {
    NrtStoreInboundResponse map(EIInboundInventoryResponse eiResponse);
}
```

**Benefit**: Compile-time validation, type-safe, performance

---

### 9. CCM2 for Configuration Management

**Pattern**: Externalized configuration
```java
@ManagedConfiguration
EiApiCCMConfig eiApiCCMConfig;
```

---

### 10. AOP for Multi-Tenancy

**Implementation**:
```java
@Aspect
public class SiteIdFilterAspect {
    @Before("@annotation(SiteIdFilter)")
    public void addSiteIdToQuery(JoinPoint joinPoint) {
        // Inject site_id into query
    }
}
```

---

## 14. TESTING STRATEGY

### Unit Testing (JUnit 5 + Mockito)

**Coverage Target**: > 80%

**Test Structure**:
```
src/test/java/com/walmart/cpnrti/
├── controller/          # Controller tests
├── services/            # Service tests
│   └── impl/           # Implementation tests
├── repository/          # Repository tests
├── utils/              # Utility tests
└── validations/        # Validator tests
```

---

### Integration Testing (Cucumber BDD)

**Framework**: Cucumber 7.19.0 + TestNG + Rest Assured

**Feature Files**:
```
src/test/resources/features/
└── NrtSandboxOnHandStoreInventory.feature
```

---

### Contract Testing (R2C)

**Framework**: R2C (Request-to-Contract)
- Threshold: 80%
- Mode: Active (fails pipeline if < 80%)

---

### Performance Testing (Automaton + JMeter)

**Test Scenarios**:
- Single GTIN lookup
- Multi-GTIN batch query
- IAC event submission
- Transaction history query

---

### Resiliency Testing (RaaS)

**Test Types**:
1. Pod Failure
2. Network Latency
3. Resource Constraint

---

## 15. ADDITIONAL NOTES

### Team & Ownership

**Code Owners**:
- @a0j0bvc
- @a0p04i1
- @a0s12wb
- @h0b091o

**Team**: DV Channel Performance Data API Team

**Contact**:
- Email: dv_dataapitechall@email.wal-mart.com
- Slack: #data-ventures-cperf-dev-ops

---

### Project Metadata

**Project Name**: cp-nrti-apis

**Repository**: https://gecgithub01.walmart.com/dsi-dataventures-luminate/cp-nrti-apis

**Maven Coordinates**:
- Group ID: com.walmart
- Artifact ID: cp-nrti-apis
- Version: 0.0.1-SNAPSHOT

**Lines of Code**: ~18,386 lines of Java code (main)

**Test Files**: 82 test files

---

### Related Services

**Upstream Dependencies**:
- Enterprise Inventory (EI) APIs
- Uber Keys API
- Location Platform API
- Identity Management API
- Sumo (push notifications)

**Downstream Consumers**:
- Supplier portals
- Vendor applications
- Internal Walmart apps
- Analytics pipelines

---

## SUMMARY

The **cp-nrti-apis** service is a sophisticated, production-grade microservice that serves as the backbone of Walmart's supplier inventory visibility platform. With support for multiple inventory operations (IAC, on-hand queries, transaction history, store inbound, DSC), comprehensive authorization controls, multi-region deployment, and event-driven architecture via Kafka, it provides a scalable and secure solution for near real-time inventory tracking. The service demonstrates excellent engineering practices including hexagonal architecture, comprehensive testing (unit, integration, contract, performance, resiliency), strong security (HMAC signatures, multi-level authorization, XSS prevention), robust observability (OpenTelemetry tracing, Prometheus metrics, Dynatrace APM), and operational excellence (canary deployments, autoscaling, multi-region HA).

---

**Document Version**: 1.0
**Last Updated**: 2026-02-02
**Analyzed By**: Claude Code Ultra Analysis
**Project**: Walmart Data Ventures Luminate
**Lines of Code Analyzed**: ~18,386 lines (Java main code)
**Total API Endpoints**: 10+ endpoints across multiple controllers

# ULTRA-THOROUGH ANALYSIS: audit-api-logs-gcs-sink

## Executive Summary

The **audit-api-logs-gcs-sink** is a Kafka Connect-based data pipeline service that streams audit API log data from Kafka topics to Google Cloud Storage (GCS) buckets. It implements geo-specific data routing for US, Canada (CA), and Mexico (MX) markets, ensuring compliance with data residency requirements by filtering and segregating records based on site identifiers before storing them in region-specific GCS buckets in Parquet format.

---

## 1. SERVICE PURPOSE AND FUNCTIONALITY

### Core Purpose
This service acts as a **Kafka Connect Sink Connector** that:
- Consumes audit log data from Kafka topics (`api_logs_audit_stg` for staging, `api_logs_audit_prod` for production)
- Filters records based on geographic market identifiers (US, CA, MX) using custom Single Message Transforms (SMT)
- Routes filtered data to dedicated GCS buckets per country
- Stores data in Parquet format with partitioning by service name, date, and endpoint name
- Supports data compliance and segregation requirements for multi-national operations

### Data Flow
```
Kafka Topic (api_logs_audit_prod/stg)
    ↓
[Kafka Connect Worker]
    ↓
[Custom Filter SMTs] → Filter by wm-site-id header
    ├── US Filter → audit-api-logs-us-prod bucket
    ├── CA Filter → audit-api-logs-ca-prod bucket
    └── MX Filter → audit-api-logs-mx-prod bucket
```

### Key Features
1. **Multi-connector architecture**: Runs 3 separate sink connectors (US, CA, MX) in parallel
2. **Header-based filtering**: Custom SMT filters examine `wm-site-id` headers to route records
3. **Timestamp injection**: Automatically adds rolling timestamp headers for partitioning
4. **Error handling**: Comprehensive DLQ (Dead Letter Queue) support with detailed error logging
5. **Retry mechanism**: Configurable retry policy with exponential backoff
6. **Multi-region deployment**: Supports both EUS2 (East US 2) and SCUS (South Central US) clusters

---

## 2. TECHNOLOGY STACK

### Build & Runtime
- **Language**: Java 17
- **Build Tool**: Maven 3.9.9
- **JDK**: Azul Zulu 17
- **Base Image**: `kcaas-base-image:11-major` (Kafka Connect as a Service)

### Core Dependencies

#### Primary Libraries (from pom.xml)
```xml
<dependencies>
    <!-- GCS Lenses Connector -->
    <groupId>com.walmart.streaming</groupId>
    <artifactId>gcs-lenses-connector</artifactId>
    <version>1.64</version>

    <!-- Kafka Connect Transforms -->
    <groupId>org.apache.kafka</groupId>
    <artifactId>connect-transforms</artifactId>
    <version>3.6.0</version>

    <!-- Avro Converter -->
    <groupId>io.confluent</groupId>
    <artifactId>kafka-connect-avro-converter</artifactId>
    <version>7.4.0</version>
</dependencies>
```

#### Testing Libraries
- JUnit Jupiter 5.10.0
- Mockito 4.11.0
- JaCoCo 0.8.11 (code coverage)

### Infrastructure & Deployment
- **Container Platform**: Kubernetes (WCNP - Walmart Cloud Native Platform)
- **Orchestration**: KITT (Walmart's deployment framework)
- **Configuration Management**: CCM (Centralized Configuration Management)
- **Secret Management**: Akeyless (for Kafka certificates, GCS credentials)
- **Service Mesh**: Istio (sidecar injection enabled)
- **CI/CD**: Looper (Walmart's CI/CD platform)
- **Code Quality**: SonarQube
- **Monitoring**: Grafana

### Messaging & Storage
- **Message Broker**: Apache Kafka 2.8.1+ (SSL/TLS secured)
- **Schema Registry**: Confluent Schema Registry (Avro schemas)
- **Storage**: Google Cloud Storage (Parquet format)
- **Data Format**: Avro (input), Parquet (output)

---

## 3. ARCHITECTURE AND CODE ORGANIZATION

### Directory Structure
```
audit-api-logs-gcs-sink/
├── src/
│   ├── main/
│   │   ├── java/com/walmart/audit/log/sink/converter/
│   │   │   ├── BaseAuditLogSinkFilter.java          # Abstract base filter
│   │   │   ├── AuditLogSinkUSFilter.java            # US-specific filter
│   │   │   ├── AuditLogSinkCAFilter.java            # Canada filter
│   │   │   ├── AuditLogSinkMXFilter.java            # Mexico filter
│   │   │   └── AuditApiLogsGcsSinkPropertiesUtil.java  # Config utility
│   │   └── resources/
│   │       ├── audit_api_logs_gcs_sink_stage_properties.yaml
│   │       └── audit_api_logs_gcs_sink_prod_properties.yaml
│   └── test/
│       └── java/com/walmart/audit/log/sink/converter/
│           ├── BaseAuditLogSinkFilterTest.java
│           ├── AuditLogSinkUSFilterTest.java
│           ├── AuditLogSinkCAFilterTest.java
│           └── AuditLogSinkMXFilterTest.java
├── ccm/
│   ├── PROD-1.0-ccm.yml                # Production CCM config
│   └── NON-PROD-1.0-ccm.yml            # Non-prod CCM config
├── pom.xml                              # Maven build configuration
├── Dockerfile                           # Container image definition
├── kitt.yml                             # WCNP deployment config
├── kc_config.yaml                       # Kafka Connect base config
├── env_properties.yaml                  # Environment-specific properties
├── sr.yaml                              # Service Registry metadata
├── .looper.yml                          # CI/CD pipeline config
├── sonar-project.properties             # SonarQube settings
└── README.md                            # Documentation
```

### Package Structure
```
com.walmart.audit.log.sink.converter
├── BaseAuditLogSinkFilter<R>           # Abstract transformer
├── AuditLogSinkUSFilter<R>             # US market filter
├── AuditLogSinkCAFilter<R>             # Canada market filter
├── AuditLogSinkMXFilter<R>             # Mexico market filter
└── AuditApiLogsGcsSinkPropertiesUtil   # Configuration loader
```

---

## 4. KEY COMPONENTS AND THEIR RESPONSIBILITIES

### 4.1 BaseAuditLogSinkFilter (Abstract Base Class)

**File**: `/src/main/java/com/walmart/audit/log/sink/converter/BaseAuditLogSinkFilter.java`

**Purpose**: Abstract Kafka Connect Transformation that provides common filtering logic for all geographic filters.

**Key Methods**:
```java
public abstract class BaseAuditLogSinkFilter<R extends ConnectRecord<R>>
    implements Transformation<R> {

    protected static final String HEADER_NAME = "wm-site-id";
    protected abstract String getHeaderValue();

    public R apply(R r) { ... }           // Main transform logic
    public boolean verifyHeader(R r) { ... }  // Header validation
    public ConfigDef config() { ... }     // Kafka Connect config
}
```

**Features**:
- Header-based filtering using `wm-site-id`
- Parallel stream processing for performance
- Comprehensive error handling with logging
- Returns `null` for filtered-out records (Kafka Connect convention)

### 4.2 Geographic Filter Implementations

#### AuditLogSinkUSFilter
**Distinguishing Feature**: **Permissive filtering** - accepts records WITH US site ID OR WITHOUT any site ID
```java
public boolean verifyHeader(R r) {
    return StreamSupport.stream(r.headers().spliterator(), true)
            .anyMatch(header -> HEADER_NAME.equals(header.key())
                && StringUtils.equals(getHeaderValue(), String.valueOf(header.value())))
        ||
        StreamSupport.stream(r.headers().spliterator(), true)
            .noneMatch(header -> HEADER_NAME.equals(header.key()));
}
```

This is critical: US filter acts as a **default catch-all** for records without geographic markers.

#### AuditLogSinkCAFilter & AuditLogSinkMXFilter
**Feature**: **Strict filtering** - only accepts records with exact matching site IDs
```java
protected String getHeaderValue() {
    return AuditApiLogsGcsSinkPropertiesUtil.getSiteIdForCountryCode("CA"); // or "MX"
}
```

Uses the inherited `verifyHeader()` from base class which is strict.

### 4.3 AuditApiLogsGcsSinkPropertiesUtil

**File**: `/src/main/java/com/walmart/audit/log/sink/converter/AuditApiLogsGcsSinkPropertiesUtil.java`

**Purpose**: Environment-aware configuration loader for site ID mappings.

**Implementation**:
```java
public class AuditApiLogsGcsSinkPropertiesUtil {
    private static final String STRATI_RCTX_ENVIRONMENT_TYPE =
        System.getenv("STRATI_RCTX_ENVIRONMENT_TYPE");

    private static final String envType = STRATI_RCTX_ENVIRONMENT_TYPE != null
        ? STRATI_RCTX_ENVIRONMENT_TYPE.toLowerCase()
        : "stage";

    private static final String FILE_NAME =
        "audit_api_logs_gcs_sink_" + envType + "_properties.yaml";

    static {
        // Loads properties at class initialization
        try (InputStream input = AuditApiLogsGcsSinkPropertiesUtil.class
                .getClassLoader().getResourceAsStream(FILE_NAME)) {
            properties.load(input);
        }
    }

    public static String getSiteIdForCountryCode(String countryCode) {
        return properties.getProperty("WM_SITE_ID_FOR_" + countryCode);
    }
}
```

**Site ID Mappings**:

**Stage Environment**:
```yaml
WM_SITE_ID_FOR_US=1694066566785477000
WM_SITE_ID_FOR_MX=1694066631493876000
WM_SITE_ID_FOR_CA=1694066641269922000
```

**Production Environment**:
```yaml
WM_SITE_ID_FOR_US=1704989259133687000
WM_SITE_ID_FOR_MX=1704989390144984000
WM_SITE_ID_FOR_CA=1704989474816248000
```

---

## 5. DEPENDENCIES AND INTEGRATIONS

### 5.1 Kafka Connect Framework
- **Connector Class**: `io.lenses.streamreactor.connect.gcp.storage.sink.GCPStorageSinkConnector`
- **Task Parallelism**: 1 task per connector (3 connectors total)
- **Converter**: Avro (value), String (key)

### 5.2 External System Integrations

#### Kafka Clusters
**Staging (EUS2)**:
```
kafka-1589333338-{1,2,3}-167536040{5,8,11}.eus.kafka-v2-luminate-core-stg
  .ms-df-messaging.stg-az-eastus2-2.prod.us.walmart.net:9093
```

**Production (EUS2)**:
```
kafka-976900980-{1,2,3}-170867044{7,50,53}.eus.kafka-v2-luminate-core-prod
  .ms-df-messaging.prod-az-eastus2-2.prod.us.walmart.net:9093
```

**Staging (SCUS)**:
```
kafka-886515205-{1,2,3}-167535939{9,2,5}.scus.kafka-v2-luminate-core-stg
  .ms-df-messaging.stg-az-southcentralus-8.prod.us.walmart.net:9093
```

**Production (SCUS)**:
```
kafka-532395416-{1,2,3}-170866930{0,3,6}.scus.kafka-v2-luminate-core-prod
  .ms-df-messaging.prod-az-southcentralus-3.prod.us.walmart.net:9093
```

#### Schema Registry
- **Stage**: `http://schema-registry-service.stage.schema-registry.ms-df-streaming.prod.walmart.com`
- **Prod**: `http://schema-registry-service.prod.schema-registry.ms-df-streaming.prod.walmart.com`

#### Google Cloud Storage
**GCS Project**: `wmt-dv-luminate-prod`

**Bucket Structure**:
- **US Stage**: `audit-api-logs-us-test:us_dv_audit_log_dev.db/api_logs`
- **US Prod**: `audit-api-logs-us-prod:us_dv_audit_log_prod.db/api_logs`
- **CA Stage**: `audit-api-logs-ca-test:ca_dv_audit_log_dev.db/api_logs`
- **CA Prod**: `audit-api-logs-ca-prod:ca_dv_audit_log_prod.db/api_logs`
- **MX Stage**: `audit-api-logs-mx-test:mx_dv_audit_log_dev.db/api_logs`
- **MX Prod**: `audit-api-logs-mx-prod:mx_dv_audit_log_prod.db/api_logs`

**Partitioning Strategy**:
```sql
PARTITIONBY service_name, _header.date, endpoint_name
```

**File Format**: Parquet with properties:
```yaml
flush.size: 50000000 bytes (50MB)
flush.count: 5000 records
flush.interval: 600 seconds (10 minutes)
key.suffix: _eus2 or _scus (regional identifier)
```

### 5.3 Kafka Topics

**Source Topics**:
- `api_logs_audit_stg` (staging)
- `api_logs_audit_prod` (production)

**Dead Letter Queue Topics**:
- `api_logs_audit_stg_DLQ`
- `api_logs_audit_prod_DLQ`

**Internal Kafka Connect Topics**:
- Config: `api_logs_audit_{env}-{env}-config`
- Offset: `api_logs_audit_{env}-{env}-offset`
- Status: `api_logs_audit_{env}-{env}-status`

### 5.4 Single Message Transforms (SMT) Chain

Each connector applies the following transform pipeline:
```yaml
transforms: InsertRollingRecordTimestamp, Filter{US|CA|MX}

# Step 1: Add timestamp header for date-based partitioning
transforms.InsertRollingRecordTimestamp.type:
  io.lenses.connect.smt.header.InsertRollingRecordTimestampHeaders
transforms.InsertRollingRecordTimestamp.date.format: "yyyy-MM-dd"
transforms.InsertRollingRecordTimestamp.timezone: GMT

# Step 2: Filter by geographic region
transforms.FilterUS.type:
  com.walmart.audit.log.sink.converter.AuditLogSinkUSFilter
```

---

## 6. DEPLOYMENT CONFIGURATION

### 6.1 Container Configuration (Dockerfile)

```dockerfile
FROM docker.ci.artifacts.walmart.com/gdap-mystique-docker/kcaas-base-image:11-major

USER 10000

# Configuration files
COPY --chown=10000:10001 pom.xml $KAFKA_HOME
COPY --chown=10000:10001 kitt.yml $KAFKA_HOME
COPY --chown=10000:10001 kc_config.yaml $KAFKA_HOME
COPY --chown=10000:10001 env_properties.yaml $KAFKA_HOME

# Uber JAR with all connectors and transforms
COPY --chown=10000:10001 target/connectors-uber.jar $PLUGIN_PATH/
```

**Security**: Runs as non-root user (UID 10000, GID 10001)

### 6.2 WCNP/Kubernetes Deployment (kitt.yml)

#### Deployment Stages

| Stage | Cluster | Region | Min CPU | Max CPU | Min Memory | Max Memory | Replicas |
|-------|---------|--------|---------|---------|------------|------------|----------|
| stage-eus2 | eus2-stage-a4 | East US 2 | 1 | 2 | 1024Mi | 5120Mi | 1 |
| stage-scus | scus-stage-a3 | South Central US | 1 | 2 | 1024Mi | 5120Mi | 1 |
| prod-eus2 | eus2-prod-a30 | East US 2 | 10 | 12 | 10Gi | 12Gi | 1 |
| prod-scus | scus-prod-a63 | South Central US | 10 | 12 | 10Gi | 12Gi | 1 |

**Production JVM Tuning**:
```yaml
KAFKA_HEAP_OPTS: "-Xmx7g -Xms5g -XX:MetaspaceSize=96m -XX:+UseG1GC
  -XX:MaxGCPauseMillis=20 -XX:InitiatingHeapOccupancyPercent=35
  -XX:G1HeapRegionSize=16M -XX:MinMetaspaceFreeRatio=50
  -XX:MaxMetaspaceFreeRatio=80 -XX:+UseStringDeduplication"
```

**Production Ephemeral Storage**:
- Min: 15GB
- Max: 20GB

#### Health Checks
```yaml
readinessProbe:
  wait: 15 seconds
livenessProbe:
  wait: 15 seconds
startupProbe:
  enabled: true
  path: /connectors
  port: 8083
  probeInterval: 10 seconds
  failureThreshold: 15
```

### 6.3 Secret Management (Akeyless)

**Path Structure**:
- **Stage**: `/Prod/WCNP/homeoffice/Luminate-CPerf-Dev-Group`
- **Production**: `/Prod/WCNP/homeoffice/dv-kys-api-prod-group`

**Secrets Injected**:
```yaml
files:
  - kafka_ssl_key_pwd.txt
  - kafka_ssl_keystore_pwd.txt
  - kafka_ssl_truststore_pwd.txt
  - kafka_ssl_keystore.jks (base64 encoded)
  - kafka_ssl_truststore.jks (base64 encoded)
  - gcs_key.json (base64 encoded)
```

### 6.4 Kafka Connect Worker Configuration

**Key Settings** (from kc_config.yaml):
```yaml
worker:
  key.converter: org.apache.kafka.connect.storage.StringConverter
  value.converter: io.confluent.connect.avro.AvroConverter
  value.converter.schemas.enable: true
  group.id: com.walmart-audit-log-gcs-sink-0.0.1

  # Security
  security.protocol: SSL
  ssl.enabled.protocols: TLSv1.2,TLSv1.1,TLSv1
  ssl.truststore.type: JKS
  ssl.keystore.type: JKS

  # Consumer tuning
  max.poll.records: 50
  max.poll.interval.ms: 300000 (5 minutes)
  heartbeat.interval.ms: 5000
  session.timeout.ms: 15000
  request.timeout.ms: 60000
```

**Three Connector Instances**:
1. `audit-log-gcs-sink-connector` (US)
2. `audit-log-gcs-sink-connector-ca` (Canada)
3. `audit-log-gcs-sink-connector-mx` (Mexico)

### 6.5 Error Handling Configuration

```yaml
errors.tolerance: all
errors.deadletterqueue.context.headers.enable: true
errors.log.enable: true
errors.log.include.messages: true

connect.gcpstorage.error.policy: RETRY
connect.gcpstorage.max.retries: 5
connect.gcpstorage.retry.interval: 5000
connect.gcpstorage.http.retry.timeout.multiplier: 1.0
```

---

## 7. CI/CD PIPELINE AND WORKFLOWS

### 7.1 Looper CI/CD (.looper.yml)

**Build Tools**:
```yaml
tools:
  jdk:
    flavor: azul
    version: 17
  maven: 3.9.9
  sonarscanner: 5.0.1.3006
```

**Pipeline Flows**:

#### Default Flow (Main Branch)
```yaml
flows:
  default:
    try:
      - call: build
      - call: sonar-fetch-origin
      - call: sonar
    finally:
      - call: hygieiaPublish
```

#### PR Flow
```yaml
  pr:
    - call: run-sonar-pr
```

**Build Command**:
```bash
mvn -B -fae clean install sonar:sonar
```

**SonarQube PR Analysis**:
```bash
sonar-scanner -X -Dproject.settings=sonar-project.properties \
  -Dsonar.pullrequest.github.repository=dsi-dataventures-luminate/audit-api-logs-gcs-sink \
  -Dsonar.pullrequest.key=${GITHUB_PR_NUMBER} \
  -Dsonar.pullrequest.branch=${GITHUB_PR_SOURCE_BRANCH} \
  -Dsonar.pullrequest.base=${GITHUB_PR_TARGET_BRANCH} \
  -Dsonar.scm.revision=${GITHUB_PR_HEAD_SHA}
```

### 7.2 Code Quality (SonarQube)

**Configuration** (sonar-project.properties):
```properties
sonar.projectKey=com.walmart:audit-log-gcs-sink
sonar.projectName=audit-log-gcs-sink
sonar.projectVersion=1.0

sonar.sources=src/main/java
sonar.tests=src/test/java
sonar.binaries=target/classes
sonar.java.binaries=target/classes
sonar.language=java
sonar.java.source=17
sonar.sourceEncoding=UTF-8

sonar.junit.reportsPath=target/surefire-reports
sonar.core.codeCoveragePlugin=jacoco
```

**Code Coverage**: JaCoCo plugin configured in pom.xml
- Unit test reports: `target/site/jacoco/jacoco.xml`
- Automatic coverage during build

### 7.3 GitHub Integration

**Repository**: `https://gecgithub01.walmart.com/dsi-dataventures-luminate/audit-api-logs-gcs-sink`

**Code Owners** (.github/CODEOWNERS):
```
* @a0j0bvc @a0p04i1 @a0s12wb @h0b091o @h0s0acv
```

**Recent Commits** (showing evolution):
```
6ced915 - Start up probe changes
3a0703b - GCS Sink Changes to support different site ids
deaaa5c - Changes for country code to site id
0df9edb - Updating the connector version
571a215 - Updating resolution path to have /* as override path
356b56d - Updating the data type as DECIMAL for timeout multiplier
02cc378 - Upgrade connector version. Remove custom header conversion class
```

### 7.4 Build Process

**Maven Build Steps**:
```bash
# Clean and install
mvn clean install

# Create uber JAR with dependencies
mvn clean dependency:copy-dependencies package

# Build Docker image
docker build --pull --no-cache -t kcaas-test .

# Run locally (for testing)
docker run --env STRATI_RCTX_ENVIRONMENT=dev \
  --env STRATI_RCTX_ENVIRONMENT_TYPE=dev \
  --env VAULT_TOKEN=<TOKEN> \
  -p 8080:8080 kcaas-test
```

**Maven Shade Plugin** - Creates uber JAR:
```xml
<plugin>
  <groupId>org.apache.maven.plugins</groupId>
  <artifactId>maven-shade-plugin</artifactId>
  <finalName>connectors-uber</finalName>
  <transformers>
    <transformer implementation=
      "org.apache.maven.plugins.shade.resource.ServicesResourceTransformer"/>
  </transformers>
</plugin>
```

---

## 8. MONITORING AND OBSERVABILITY

### 8.1 Grafana Dashboards

**1. Kafka Connect Golden Signals**
```
https://grafana.mms.walmart.net/d/68r2wpy7z/kafka-connect-golden-signals
?var-namespace=data-ventures-luminate-cperf
&var-container=audit-log-gcs-sink
```

**2. Kafka Connect Deployment Dashboard**
```
https://grafana.mms.walmart.net/d/jmBvWG4nk/kafka-connect-deployments-dashboard
?var-namespace=data-ventures-luminate-cperf
&var-container=audit-log-gcs-sink
```

**3. Consumer Lag Dashboard**
```
https://grafana.mms.walmart.net/d/o7eRAUOpwerz3/lenses-consumergroup-lag
?var-consumerGroup=connect-audit-log-gcs-sink-connector
```

### 8.2 Alerting

**Consumer Lag Alert Definition**:
```
https://gecgithub01.walmart.com/Telemetry/mms-config/blob/main/
  data-strategy-and-insights/data-ventures/luminate/rules/
  production/kafka-consumer-lag.yaml
```

**Notification Channels**:
- Slack: `data-ventures-cperf-dev-ops`
- Email: `dv_dataapitechall@email.wal-mart.com`

### 8.3 Service Registry (sr.yaml)

**Applications Registered**:
1. **audit-log-gcs-sink** (AUDIT-LOG-GCS-SINK)
2. **audit-log-gcs-sink-intl** (AUDIT-LOG-GCS-SINK-INTL)

**Metadata**:
```yaml
key: AUDIT-LOG-GCS-SINK
description: Kafka Connect Sink for Walmart Luminate api logs
organization: dsi-dataventures-luminate
companyCatalog: true
businessCriticality: MAJOR
```

**Team Members** (11 members):
- homeoffice\h0b091o, h0s0acv, g0a0070, a0c0sdo, a0j0bvc, vn540u6,
  vn53gpt, p0n01az, a0s12wb, vn54en1, b0k03a2

### 8.4 Connector REST API

**Base URL**: `http://localhost:8080` (in local), `http://<pod-ip>:8083` (production)

**Endpoints**:
- `GET /connectors` - List all connectors
- `GET /connectors/{name}` - Get connector info
- `GET /connectors/{name}/config` - Get connector config
- `GET /connectors/{name}/status` - Get connector status
- `GET /connectors/restart_all` - Restart all connectors

---

## 9. SPECIAL PATTERNS AND NOTEWORTHY IMPLEMENTATIONS

### 9.1 Multi-Connector Pattern for Geographic Segregation

**Design Decision**: Instead of a single connector with complex routing logic, the service deploys **3 parallel connectors** that each:
1. Consume from the same Kafka topic
2. Apply their own geographic filter
3. Write to separate GCS buckets

**Benefits**:
- Independent scaling per region
- Isolated failure domains
- Simpler debugging (logs are per-region)
- Parallel processing for higher throughput

**Trade-off**: Same record read 3 times from Kafka (acceptable for audit logs)

### 9.2 Permissive US Filter Pattern

The US filter has special logic to act as a catch-all:
```java
// Accept if: (has US site ID) OR (has NO site ID)
anyMatch(header -> HEADER_NAME.equals(header.key()) && matches("US"))
  ||
noneMatch(header -> HEADER_NAME.equals(header.key()))
```

**Rationale**: Ensures no audit logs are lost if site ID header is missing.

### 9.3 Environment-Aware Site ID Resolution

Site IDs differ between environments (stage vs. prod), and the utility class automatically loads the correct mapping at startup:

```java
private static final String envType =
    STRATI_RCTX_ENVIRONMENT_TYPE != null
    ? STRATI_RCTX_ENVIRONMENT_TYPE.toLowerCase()
    : "stage";

private static final String FILE_NAME =
    "audit_api_logs_gcs_sink_" + envType + "_properties.yaml";
```

**Pattern**: Environment-driven resource loading with sensible defaults.

### 9.4 Parquet Partitioning Strategy

**Three-level partitioning**:
```sql
PARTITIONBY service_name, _header.date, endpoint_name
```

**Benefits**:
- Efficient query pruning by service
- Time-based data management (by date)
- Endpoint-specific analysis support

**Example Directory Structure**:
```
gs://audit-api-logs-us-prod/us_dv_audit_log_prod.db/api_logs/
  service_name=NRT/
    date=2025-02-02/
      endpoint_name=transactionHistory/
        part-00000_eus2.parquet
        part-00001_eus2.parquet
```

### 9.5 Dual-Region Active-Active Architecture

**Deployment**: Service runs in both EUS2 and SCUS simultaneously with:
- Separate Kafka clusters per region
- Separate offset management (no shared state)
- Regional suffixes on file keys (`_eus2`, `_scus`)
- Separate index files (`.indexes-eus2`, `.indexes-scus`)

**Pattern**: Active-active for high availability and disaster recovery.

### 9.6 Secret Reference Pattern

Secrets are not embedded in configuration but referenced:
```yaml
ssl.keystore.location: secret.ref://kafka_ssl_keystore.jks?encoding=base64
ssl.keystore.password: secret.value://kafka_ssl_keystore_pwd.txt
gcs_key.json: secret.ref://gcs_key.json
```

**Mechanism**: Kafka Connect resolves these at runtime by reading from the secrets mount.

### 9.7 Error Handling Strategy

**Three-tier approach**:
1. **Retry**: 5 retries with 5-second intervals
2. **DLQ**: Failed records sent to dedicated Dead Letter Queue
3. **Logging**: Full error context including headers

**Configuration**:
```yaml
errors.tolerance: all  # Continue processing despite errors
errors.deadletterqueue.context.headers.enable: true
errors.log.enable: true
errors.log.include.messages: true
```

### 9.8 KCQL (Kafka Connect Query Language)

Uses Lenses' KCQL for declarative data pipeline:
```sql
INSERT INTO `audit-api-logs-us-prod:us_dv_audit_log_prod.db/api_logs`
SELECT *
FROM `api_logs_audit_prod`
PARTITIONBY service_name, _header.date, endpoint_name
STOREAS `PARQUET` PROPERTIES(
  'flush.size'='50000000',
  'flush.count'='5000',
  'flush.interval'='600',
  'key.suffix'='_eus2'
);
```

**Advantages**: SQL-like syntax, built-in partitioning, format conversion.

### 9.9 Test Strategy

**Unit Tests**: Each filter has comprehensive tests covering:
- Valid site ID matching
- Invalid site ID rejection
- Missing header handling
- Exception scenarios

**Test Coverage**:
- JaCoCo reports generated automatically
- Integrated with SonarQube
- Coverage tracked per PR

**Example Test Pattern**:
```java
@Test
void testUsFilterWithValidWmSiteId() {
    AuditLogSinkUSFilter<SinkRecord> converter = new AuditLogSinkUSFilter<>();
    SinkRecord record = getUsRecord(HEADER_WM_SITE_ID, HEADER_WM_SITE_ID_VALUE_US);
    SinkRecord transformedRecord = converter.apply(record);
    assertNotNull(transformedRecord);
}

@Test
void testUsFilterWithInvalidWmSiteId() {
    AuditLogSinkUSFilter<SinkRecord> converter = new AuditLogSinkUSFilter<>();
    SinkRecord record = getUsRecord(HEADER_WM_SITE_ID, HEADER_WM_SITE_ID_VALUE_MX);
    SinkRecord transformedRecord = converter.apply(record);
    assertNull(transformedRecord); // Filtered out
}
```

### 9.10 CCM Integration Pattern

**Two-tier configuration**:
1. **Base Config** (kc_config.yaml, env_properties.yaml): Checked into git
2. **Environment Overrides** (CCM): Dynamically resolved based on environment path

**Resolution Path** (NON-PROD):
```yaml
resolutionPaths:
  - default: /envName
```

**Resolution Path** (PROD):
```yaml
resolutionPaths:
  - default: /envProfile/envName
```

**Override Example**:
```yaml
configOverrides:
  connectorConfig:
    - name: "prod-eus2"
      pathElements:
        envProfile: "prod"
        envName: "prod-eus2"
      value:
        properties:
          topics: "api_logs_audit_prod"
          connect.gcpstorage.kcql: >
            INSERT INTO `audit-api-logs-us-prod:...`
```

---

## 10. DATA FLOW DIAGRAM

```
┌─────────────────────────────────────────────────────────────┐
│                      Kafka Cluster                           │
│                  api_logs_audit_prod topic                   │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│            Kafka Connect Workers - EUS2 Pod                  │
│                                                               │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │US Connector  │  │CA Connector  │  │MX Connector  │      │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘      │
│         │                  │                  │              │
│         ▼                  ▼                  ▼              │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │InsertTime +  │  │InsertTime +  │  │InsertTime +  │      │
│  │FilterUS      │  │FilterCA      │  │FilterMX      │      │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘      │
│         │                  │                  │              │
└─────────┼──────────────────┼──────────────────┼──────────────┘
          │                  │                  │
          │ wm-site-id=US    │ wm-site-id=CA    │ wm-site-id=MX
          │ or null          │                  │
          ▼                  ▼                  ▼
┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐
│ GCS Bucket:     │  │ GCS Bucket:     │  │ GCS Bucket:     │
│ audit-api-logs- │  │ audit-api-logs- │  │ audit-api-logs- │
│ us-prod         │  │ ca-prod         │  │ mx-prod         │
│                 │  │                 │  │                 │
│ Parquet files   │  │ Parquet files   │  │ Parquet files   │
│ partitioned by: │  │ partitioned by: │  │ partitioned by: │
│ service/date/   │  │ service/date/   │  │ service/date/   │
│ endpoint        │  │ endpoint        │  │ endpoint        │
└─────────────────┘  └─────────────────┘  └─────────────────┘
```

---

## 11. CONFIGURATION MATRIX

### Connector Configuration Comparison

| Configuration | US Connector | CA Connector | MX Connector |
|--------------|-------------|-------------|-------------|
| **Connector Name** | audit-log-gcs-sink-connector | audit-log-gcs-sink-connector-ca | audit-log-gcs-sink-connector-mx |
| **Filter Class** | AuditLogSinkUSFilter | AuditLogSinkCAFilter | AuditLogSinkMXFilter |
| **GCS Bucket (Prod)** | audit-api-logs-us-prod | audit-api-logs-ca-prod | audit-api-logs-mx-prod |
| **GCS Bucket (Stage)** | audit-api-logs-us-test | audit-api-logs-ca-test | audit-api-logs-mx-test |
| **Site ID (Stage)** | 1694066566785477000 | 1694066641269922000 | 1694066631493876000 |
| **Site ID (Prod)** | 1704989259133687000 | 1704989474816248000 | 1704989390144984000 |
| **Filter Logic** | Permissive (accepts null) | Strict (exact match) | Strict (exact match) |
| **Tasks** | 1 | 1 | 1 |

---

## 12. OPERATIONAL CONSIDERATIONS

### 12.1 Resource Sizing Rationale

**Production**: 10-12 CPU, 10-12Gi memory
- High CPU for Parquet compression
- Large heap for buffering before flush (50MB flush size)
- G1GC for low-latency garbage collection

**Staging**: 1-2 CPU, 1-5Gi memory
- Lower traffic volume
- Cost optimization

### 12.2 Flush Behavior

Records are flushed to GCS when ANY condition is met:
1. **Size threshold**: 50MB accumulated
2. **Count threshold**: 5000 records buffered
3. **Time threshold**: 600 seconds elapsed

**Trade-off**: Balances latency vs. efficiency (fewer small files).

### 12.3 Scalability

**Current**: 1 worker per region = 2 workers total per environment

**Scaling Options**:
1. Horizontal: Increase replica count in kitt.yml
2. Vertical: Increase CPU/memory limits
3. Connector tasks: Currently 1, can increase to N (splits topic partitions)

**Limitation**: Kafka topic partition count limits parallelism.

### 12.4 Disaster Recovery

**Data Durability**:
- Kafka: Replicated topic (3 replicas typical)
- GCS: Multi-region storage class (99.95% availability)

**Service Resilience**:
- Dual-region deployment (EUS2 + SCUS)
- Independent offset tracking (resume from failure point)
- DLQ for poison pill records

### 12.5 Security Posture

**In-Transit**:
- Kafka: SSL/TLS (mutual TLS with client certs)
- GCS: HTTPS with service account authentication

**At-Rest**:
- GCS: Encrypted by default (Google-managed keys)
- Secrets: Akeyless encrypted storage

**Access Control**:
- Kafka: ACLs per topic
- GCS: IAM with principle of least privilege
- Kubernetes: RBAC + service mesh (Istio)

### 12.6 Compliance

**Data Residency**: Geographic filtering ensures:
- Canadian data stays in CA bucket (subject to PIPEDA)
- Mexican data in MX bucket (subject to LFPDPPP)
- US data in US bucket (subject to various state laws)

**Audit Trail**: Service itself is auditing API calls - meta-audit system.

---

## 13. LIMITATIONS AND KNOWN ISSUES

### 13.1 Design Limitations
1. **Triple Read**: Each record is consumed 3 times (once per connector)
   - **Impact**: 3x network bandwidth, 3x processing
   - **Mitigation**: Acceptable for low-volume audit logs

2. **Single Task**: Each connector runs 1 task
   - **Impact**: Limited parallelism
   - **Mitigation**: Can increase if partition count increases

3. **Permissive US Filter**: Records without site ID go to US bucket
   - **Impact**: May include non-US records
   - **Rationale**: Ensures no data loss

### 13.2 Operational Considerations
1. **Startup Probe**: 15 failures * 10s interval = 150s max startup time
2. **Flush Latency**: Up to 10 minutes before data appears in GCS
3. **Schema Evolution**: Avro schema changes require connector restart

---

## 14. FUTURE ENHANCEMENTS (Potential)

Based on the codebase analysis, potential improvements could include:

1. **Topic Compaction**: Reduce triple consumption via pre-filtering topic
2. **Metrics Export**: Custom Prometheus metrics for filter statistics
3. **Dynamic Site ID**: Load site IDs from CCM instead of static files
4. **Multi-Task Scaling**: Increase task count for higher throughput
5. **Compression**: Enable Parquet compression (Snappy/ZSTD)
6. **Schema Projection**: Select specific fields to reduce storage

---

## 15. SUMMARY

The **audit-api-logs-gcs-sink** is a well-architected, production-grade Kafka Connect service that demonstrates:

### Strengths
- **Separation of Concerns**: Custom SMTs cleanly encapsulate filtering logic
- **Multi-Region HA**: Active-active deployment across two Azure regions
- **Compliance-Ready**: Geographic data segregation built-in
- **Operational Excellence**: Comprehensive monitoring, alerting, and error handling
- **Security**: End-to-end encryption, secret management, least-privilege access
- **Maintainability**: Clear code structure, extensive testing, SonarQube integration
- **Scalability**: Kubernetes-native with resource limits and autoscaling support

### Key Technologies
- **Java 17** with Maven build system
- **Kafka Connect** framework with custom Single Message Transforms
- **Google Cloud Storage** with Parquet storage format
- **WCNP/Kubernetes** deployment platform
- **Akeyless** for secret management
- **Looper** for CI/CD automation

### Business Value
Provides a reliable, compliant, and scalable solution for ingesting and segregating audit logs from Walmart's Luminate API platform across multiple geographic regions, enabling data analytics while maintaining regulatory compliance.

### Files Analyzed
- **5 Java source files** (4 filters + 1 utility)
- **4 Java test files** (full test coverage)
- **13 configuration files** (deployment, build, CI/CD)
- **2 environment property files** (stage/prod site IDs)
- **2 CCM configuration files** (centralized config management)
- **1 Dockerfile** (container definition)
- **1 README.md** (documentation)

This service represents a mature, enterprise-grade implementation of a data pipeline with proper observability, security, and operational practices.

---

**Document Version**: 1.0
**Last Updated**: 2026-02-02
**Analyzed By**: Claude Code Ultra Analysis
**Project**: Walmart Data Ventures Luminate

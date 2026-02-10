# Kafka Publisher Service - audit-api-logs-srv

> **Resume Bullet:** 1 (Part of Audit Logging System)
> **Repo:** `audit-api-logs-srv`

---

## What This Service Does

The publisher service is **Tier 2** of the audit logging architecture. It receives audit payloads via HTTP from the common library and publishes them to Kafka with Avro serialization.

```
Common Library (Tier 1)
    │
    │ HTTP POST (async)
    ▼
┌─────────────────────────────────────────────────────────────────────┐
│                    audit-api-logs-srv (Tier 2)                       │
│                                                                      │
│  ┌──────────────────┐    ┌────────────────────┐    ┌─────────────┐ │
│  │AuditLogging      │ →  │KafkaProducer       │ →  │ Kafka Topic │ │
│  │Controller        │    │Service             │    │ api_logs_    │ │
│  │POST /v1/logReq   │    │- Avro serialization│    │ audit_prod  │ │
│  │                  │    │- Add wm-site-id    │    │             │ │
│  │                  │    │- Set Kafka headers  │    │ Multi-Region│ │
│  └──────────────────┘    └────────────────────┘    └─────────────┘ │
└─────────────────────────────────────────────────────────────────────┘
```

---

## Key Files (verified from repo)

```
audit-api-logs-srv/src/main/java/com/walmart/audit/
├── builder/
│   └── AuditKafkaPayloadBuilder.java      # Builds Kafka messages
├── controllers/
│   └── AuditLoggingController.java        # REST endpoint
├── kafka/
│   ├── KafkaProducerConfig.java           # Kafka configuration
│   └── KafkaProducerService.java          # Kafka publisher (primary + secondary)
├── models/
│   ├── AuditKafkaPayload.java             # Kafka payload
│   ├── AuditKafkaPayloadBody.java         # Message body
│   └── AuditKafkaPayloadKey.java          # Message key (service|endpoint)
├── services/
│   ├── ExecutorPoolService.java           # Thread pool for async publishing
│   └── LoggingRequestService.java         # Request processing service
└── common/config/
    └── AuditLogsKafkaCCMConfig.java       # CCM configuration
```

---

## Code Deep Dive

### 1. AuditLoggingController.java - REST Endpoint

```java
@Slf4j
@RestController
@RequestMapping(produces = MediaType.APPLICATION_JSON_VALUE)
public class AuditLoggingController implements AuditLogsApi {

  @Autowired
  LoggingRequestService loggingRequestService;

  @PostMapping("v1/logRequest")
  public ResponseEntity<Boolean> saveRequest(@RequestBody LoggingApiRequest request) {
    Boolean result = loggingRequestService.processLoggingRequest(request);
    return ResponseEntity.status(HttpStatus.NO_CONTENT).body(result);
  }

  @Override  // From OpenAPI generated interface
  public ResponseEntity<Void> saveApiLog(LoggingApiRequest request) {
    loggingRequestService.processLoggingRequest(request);
    return new ResponseEntity<>(HttpStatus.NO_CONTENT);
  }
}
```

**Why `NO_CONTENT` (204)?**
- No response body needed - caller doesn't wait for confirmation
- Faster response - saves bandwidth
- Standard pattern for fire-and-forget operations

---

### 2. KafkaProducerService.java - Kafka Publisher

```java
@Slf4j
@Service
public class KafkaProducerService implements TargetedResources {

  @ManagedConfiguration
  AuditLogsKafkaCCMConfig auditLogsKafkaCCMConfig;

  @Autowired
  private KafkaTemplate<String, Message<LogEvent>> kafkaPrimaryTemplate;

  @Autowired
  private KafkaTemplate<String, Message<LogEvent>> kafkaSecondaryTemplate;

  public void publishMessageToTopic(LoggingApiRequest request) {

    String topicName = auditLogsKafkaCCMConfig.getAuditLoggingKafkaTopicName();
    Message<LogEvent> kafkaMessage = prepareAuditLoggingKafkaMessage(request, topicName);

    try {
      kafkaPrimaryTemplate.send(kafkaMessage);  // EUS2 region
    } catch (Exception ex) {
      log.error("Primary Kafka failed, trying secondary", ex);
      kafkaSecondaryTemplate.send(kafkaMessage);  // SCUS fallback
    }
  }

  public static Message<LogEvent> prepareAuditLoggingKafkaMessage(
      LoggingApiRequest request, String topicName) {

    // 1. Convert to Avro
    LogEvent event = AvroUtils.getLogEvent(request);

    // 2. Create message key for partitioning
    AuditKafkaPayloadKey key = AuditKafkaPayloadKey.getAuditKafkaPayloadKey(
        request.getEndpointName(),
        request.getServiceName()
    );

    // 3. Build message with headers
    MessageBuilder<LogEvent> builder = MessageBuilder.withPayload(event);

    // 4. Set Kafka headers
    builder.setHeader(KafkaHeaders.TOPIC, topicName);
    builder.setHeader(KafkaHeaders.KEY, key.toString());

    // 5. Pass through important headers for routing
    Set<String> allowedHeaders = Set.of(
        "wm_consumer.id",
        "wm_qos.correlation_id",
        "wm_svc.name",
        "wm_svc.version",
        "wm_svc.env",
        "wm-site-id"  // CRITICAL for geo-routing
    );

    request.getHeaders().forEach((k, v) -> {
      if (allowedHeaders.contains(k.toLowerCase())) {
        builder.setHeader(k, v);
      }
    });

    return builder.build();
  }
}
```

---

## Why Avro Serialization?

| Aspect | JSON | Avro |
|--------|------|------|
| **Size** | ~100 bytes per record | ~30 bytes (70% smaller) |
| **Schema** | None - anything goes | Enforced via Schema Registry |
| **Evolution** | Breaking changes | Backward/forward compatible |
| **Speed** | Slow (text parsing) | Fast (binary) |

### How Avro + Schema Registry Works

```
PRODUCER:
1. Has schema: {type: record, name: LogEvent, fields: [...]}
2. Registers schema → Gets schema_id = 123
3. Serializes payload to binary
4. Sends: [0][123][binary_data] to Kafka

CONSUMER:
1. Receives: [0][123][binary_data]
2. Extracts schema_id = 123
3. Fetches schema from Registry (cached)
4. Deserializes binary_data using schema
```

---

## Key Design Decisions

| Decision | Why | Alternative |
|----------|-----|-------------|
| **Separate service** | Decouples library from Kafka | Direct Kafka from library (fewer hops but more coupling) |
| **Avro serialization** | Schema enforcement + 70% smaller | JSON (no schema, larger) |
| **Dual KafkaTemplate** | Primary + Secondary region failover | Single region (no DR) |
| **Header forwarding** | Enables geo-routing downstream | Multiple topics (complex) |
| **OpenAPI interface** | Contract-first design | Manual controller |

---

## Dual-Region Publishing

```java
// Try primary (EUS2), fallback to secondary (SCUS)
public void publishMessageToTopic(LoggingApiRequest request) {
    Message<LogEvent> message = prepareMessage(request);

    try {
        kafkaPrimaryTemplate.send(message).get();  // EUS2
    } catch (Exception ex) {
        log.warn("Primary Kafka failed, trying secondary: {}", ex.getMessage());
        kafkaSecondaryTemplate.send(message).get();  // SCUS fallback
    }
}
```

---

## Interview Questions

### Q: "Why a separate publisher service instead of publishing directly from the library?"
> "Separation of concerns. The library's job is to intercept and forward. The publisher handles Kafka connectivity, Avro serialization, and multi-region failover. This means services using the library don't need Kafka dependencies. If I rebuilt this today, I'd consider direct publishing to reduce latency, but the current design keeps the library lightweight."

### Q: "Why forward selected headers to Kafka?"
> "The `wm-site-id` header is critical - it tells the downstream GCS sink which geographic bucket to write to. We also forward correlation IDs for distributed tracing. We use an allowlist pattern to avoid forwarding sensitive headers."

### Q: "What happens if both Kafka regions are down?"
> "The publisher returns an error, which the common library catches and logs. The API response is not affected because the library call is async. We'd see increased error rates in our audit service metrics and get alerted."

---

## Key PRs

| PR# | Title | Why Important |
|-----|-------|---------------|
| **#4** | Initial Kafka service | First implementation |
| **#44** | [CRQ:CHG3024893] Production deployment | First prod release |
| **#49-51** | Payload size increase to 2MB | Fixed 413 errors |
| **#57** | International changes | CA/MX support |
| **#65** | Publish in both regions | Multi-region Kafka |
| **#67** | Allow selected headers | Geo-routing support |
| **#77** | Contract test changes | R2C testing |

---

*Next: [04-gcs-sink.md](./04-gcs-sink.md) - Kafka Connect GCS Sink*

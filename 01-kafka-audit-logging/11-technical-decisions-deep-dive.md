# Technical Decisions Deep Dive - "Why THIS and Not THAT?"

> **Every major decision in the Kafka Audit Logging system with alternatives, trade-offs, and follow-up answers.**
> This is what separates "good" from "excellent" in a hiring manager interview.

---

## Decision 1: Servlet Filter vs AOP vs Sidecar

### What's In The Code
```java
@Order(Ordered.LOWEST_PRECEDENCE)
@Component
public class LoggingFilter extends OncePerRequestFilter { ... }
```

### The Options

| Approach | Pros | Cons | Our Choice |
|----------|------|------|------------|
| **Servlet Filter** | Access raw HTTP body stream, runs in filter chain, captures auth failures | In-process overhead | **YES** |
| **AOP (@Around)** | Clean, annotation-based | Can't access raw HTTP body, only method params | No |
| **Sidecar (Envoy)** | Zero application changes | No application context (endpoint name, supplier ID, error details) | No |

### Interview Answer
> "AOP can only access method parameters and return values - not the raw HTTP body stream. For debugging, suppliers need the actual request body. A sidecar like Envoy works at the network layer - it can't see which endpoint was called semantically or which supplier made the request. Servlet filter gives us both: raw HTTP access AND application context."

### Follow-up: "What about performance impact?"
> "The filter adds ~1ms for ContentCachingWrapper. The actual audit is @Async so it doesn't block the response. P99 impact is under 5ms."

### Follow-up: "Why LOWEST_PRECEDENCE?"
> "So it runs LAST in the filter chain - after security filters. This means we capture the final state of the request, including auth failures. If we ran first, we'd miss what the security filter did."

---

## Decision 2: @Async Fire-and-Forget vs Synchronous vs Message Queue

### What's In The Code
```java
@Async
public void sendAuditLogRequest(AuditLogPayload payload, AuditLoggingConfig config) {
    try {
        httpService.sendHttpRequest(...);
    } catch (Exception e) {
        log.error("Audit failed for {}", payload.getRequestId(), e);
        // Don't throw - audit failure shouldn't break API
    }
}
```

### The Options

| Approach | Latency Impact | Data Safety | Complexity | Our Choice |
|----------|---------------|-------------|------------|------------|
| **@Async fire-and-forget** | Zero (non-blocking) | Can lose data if queue full | Simple | **YES** |
| **Synchronous** | +50ms per request | No loss | Simplest | No - blocks API |
| **In-process Kafka producer** | Low (~5ms) | Durable | Medium | Considered for v2 |
| **Reactive (WebFlux)** | Zero | Same as async | Complex | Overkill |

### Interview Answer
> "Synchronous would add 50ms to every API response - unacceptable for supplier-facing APIs with 200ms SLA. A direct Kafka producer from the library would be more durable but adds Kafka dependencies to every service. @Async with a bounded thread pool is the simplest approach that gives us zero latency impact. The trade-off is we can lose audit data if the queue fills - we mitigate with Prometheus metrics and 80% queue warning."

### Follow-up: "What if you lose audit data?"
> "It's a conscious trade-off. Audit data is valuable but not worth blocking API responses. We monitor rejected task count. In practice, we've lost less than 0.01% of events."

### Follow-up: "Why not publish directly to Kafka from the library?"
> "That would couple every consuming service to Kafka. They'd need Kafka dependencies, Schema Registry config, SSL certificates. The HTTP-based approach keeps the library lightweight - just add a Maven dependency. If I rebuilt this, I'd consider direct Kafka for the latency benefit, but the trade-off in coupling was the right call for adoption."

---

## Decision 3: ThreadPool 6 Core / 10 Max / 100 Queue

### What's In The Code
```java
executor.setCorePoolSize(6);
executor.setMaxPoolSize(10);
executor.setQueueCapacity(100);
```

### Why These Numbers?

| Setting | Value | Reasoning |
|---------|-------|-----------|
| **Core: 6** | ~100 req/sec × 50ms/audit = 5 threads needed. +1 headroom | Based on traffic analysis |
| **Max: 10** | Handles spikes without thread explosion | 2x normal is enough |
| **Queue: 100** | Buffer bursts. Each payload ~2KB = 200KB total | Won't cause OOM |

### Interview Answer
> "Based on traffic analysis. Our APIs handle about 100 requests per second per pod. Each audit HTTP call takes about 50ms. So we need 5 threads to keep up. I set 6 core for headroom. Max 10 handles traffic spikes. Queue of 100 handles burst traffic - at 2KB per payload, that's only 200KB of memory."

### Follow-up: "A senior engineer challenged the queue size..."
> "He was right. When the queue fills, tasks are silently rejected. I added three safeguards: Prometheus metric for rejected tasks, WARN log at 80% capacity, and documentation explaining the trade-off. That 80% warning has triggered once - caught a downstream slowdown."

---

## Decision 4: Avro vs JSON vs Protobuf

### What's In The Code
```java
LogEvent event = AvroUtils.getLogEvent(request);
// Serialized to binary with Schema Registry
```

### The Options

| Format | Size | Schema | Evolution | Kafka Ecosystem | Our Choice |
|--------|------|--------|-----------|-----------------|------------|
| **Avro** | 30% of JSON (~70% smaller) | Enforced via Schema Registry | Backward/forward compatible | Native (Confluent) | **YES** |
| **JSON** | 100% (baseline) | None | Breaking changes possible | Works but no schema | No |
| **Protobuf** | ~30% of JSON | .proto files | Compatible | Supported but less native | Considered |

### Interview Answer
> "Three reasons: schema enforcement prevents bad data from entering the pipeline, binary format is 70% smaller which matters at 2 million events per day - that's significant Kafka storage savings, and schema evolution lets us add fields without breaking consumers. We chose Avro over Protobuf because it integrates natively with the Confluent Kafka ecosystem and Schema Registry."

### Follow-up: "What if Schema Registry is down?"
> "Producers cache schemas locally after first fetch. Existing schemas continue working. New schema registrations fail until Registry recovers - we monitor Registry health separately."

---

## Decision 5: Kafka Connect vs Custom Consumer vs Flink

### The Options

| Approach | Dev Time | Offset Mgmt | Scaling | Error Handling | Our Choice |
|----------|----------|-------------|---------|----------------|------------|
| **Kafka Connect** | 2-3 days | Built-in | Automatic | DLQ + retries built-in | **YES** |
| **Custom Consumer** | 2-3 weeks | Manual | Manual | Manual | No |
| **Apache Flink** | 1-2 weeks | Built-in | Auto | Complex | Overkill |
| **Cloud Functions** | 1 week | Event-driven | Auto | Cold starts | No |

### Interview Answer
> "Kafka Connect gives us offset management, scaling, fault tolerance, and retry logic out of the box. A custom consumer would take 2-3 weeks to build the same features. Flink is powerful but overkill for a simple sink operation. Connect is battle-tested and let us focus on our SMT filter logic instead of infrastructure."

### Follow-up: "But you had so many production issues with Connect..."
> "True - the debugging incident exposed that Connect's defaults aren't production-ready. But those were configuration issues, not architectural ones. A custom consumer would have had different but equally challenging production issues - offset management bugs, rebalancing problems, OOM. The infrastructure Connect provides is worth the configuration effort."

---

## Decision 6: SMT Filter (Sink-Side) vs Producer-Side Routing

### What's In The Code
```java
public abstract class BaseAuditLogSinkFilter<R extends ConnectRecord<R>>
    implements Transformation<R> {
    public R apply(R record) {
        if (verifyHeader(record)) return record;   // Pass
        return null;                                 // Drop
    }
}
```

### The Options

| Approach | Where Filtering Happens | Adding New Country | Coupling | Our Choice |
|----------|------------------------|-------------------|----------|------------|
| **SMT Filter (sink-side)** | At consumption | Deploy new connector only | Loose | **YES** |
| **Producer-side routing** | At publishing | Change producer code | Tight | No |
| **Multiple topics** | By topic name | New topic + consumer | Medium | Considered |

### Interview Answer
> "Separation of concerns. The producer's job is to publish reliably, not to know about downstream routing. If we add a new country, we just deploy a new connector with a new filter - no changes to the publisher. Also, the US filter is permissive (accepts records with OR without headers) which handles legacy data gracefully. That logic belongs at the sink, not the producer."

---

## Decision 7: Parquet vs JSON vs ORC for Storage

### The Options

| Format | Compression | Query Speed | BigQuery | Cost at 2M/day | Our Choice |
|--------|------------|-------------|----------|-----------------|------------|
| **Parquet** | 90% | Fast (columnar) | Native external tables | ~$500/month | **YES** |
| **JSON files** | ~30% | Slow (scan all) | Needs transformation | ~$5,000/month | No |
| **ORC** | 85% | Fast (columnar) | Not native | N/A | No |
| **BigQuery direct** | N/A | Fastest | Native | ~$3,000/month (streaming) | No (cost) |

### Interview Answer
> "Parquet gives us 90% compression, columnar format so BigQuery reads only needed columns, and native BigQuery support via external tables. JSON files would be 10x larger and 10x more expensive. Direct BigQuery streaming inserts are fast but expensive at 2M events/day. Parquet in GCS is the sweet spot: cheap storage, powerful queries."

### Follow-up: "Why not Apache Iceberg?"
> "Iceberg would give us ACID transactions and easier GDPR deletes. It's on our roadmap. At the time, raw Parquet was simpler and met our requirements. If I rebuilt this today, I'd use Iceberg from day one."

---

## Decision 8: Active/Active vs Active/Passive for Multi-Region

### The Options

| Approach | Failover Time | Cost | Complexity | Our Choice |
|----------|--------------|------|------------|------------|
| **Active/Active** | ~15 min (achieved) | 2x infrastructure | Complex (dedup needed) | **YES** |
| **Active/Passive** | ~30 min | 1.5x (standby wasted) | Simple | No |
| **Single region** | N/A (no failover) | 1x | None | Not an option (compliance) |

### Interview Answer
> "Audit is write-heavy and compliance couldn't tolerate 30-minute failover delays. Active/Active gives immediate failover. The cost is 2x infrastructure, but compliance failure costs far more. We achieve 15-minute recovery, well under our 1-hour RTO target."

### Follow-up: "How do you handle duplicates?"
> "Geographic routing via wm-site-id header minimizes overlap. Kafka idempotent producer prevents duplicates within same cluster. Deduplication in sink based on request_id handles any remaining overlap. Worst case: duplicates, never data loss."

---

## Decision 9: ContentCachingWrapper vs Custom Buffering

### What's In The Code
```java
ContentCachingRequestWrapper contentCachingRequestWrapper =
    new ContentCachingRequestWrapper(request);
ContentCachingResponseWrapper contentCachingResponseWrapper =
    new ContentCachingResponseWrapper(response);
```

### Why Spring's Built-in Wrapper?

| Approach | Body Access | Memory | Edge Cases | Our Choice |
|----------|------------|--------|------------|------------|
| **ContentCachingWrapper** | Read twice (cached) | Entire body in memory | Handles charset, multipart | **YES** |
| **Custom InputStream wrapper** | Complex, error-prone | Same | Must handle all edge cases | No |
| **Read body, pass bytes** | Requires changing all controllers | Same | Controller changes everywhere | No |

### Interview Answer
> "HTTP request bodies are InputStreams - you can only read them once. Without caching, either the filter or controller reads the body, not both. ContentCachingWrapper is Spring's built-in solution - well-tested, handles charset encoding, multipart, and other edge cases. A custom wrapper would require re-implementing all of that. The trade-off is memory - the entire body is in memory - but our gateway enforces a 2MB payload limit."

---

## Decision 10: isResponseLoggingEnabled Config Flag

### What's In The Code
```java
if (Boolean.TRUE.equals(auditLoggingConfig.isResponseLoggingEnabled())) {
    responseBody = convertByteArrayToString(
        contentCachingResponseWrapper.getContentAsByteArray(), ...);
}
```

### Why A Config Flag?

> "One team wanted response body logging for debugging. Another team didn't - their responses were large (100KB+) and would increase storage costs 10x. Making it a CCM config flag means each team decides independently. No code changes, no redeployment - just a config toggle. This was the key to library adoption across 3 teams - the 20% that was different became configurable."

---

## Quick Reference: All 10 Decisions

| # | Decision | Chose | Over | Key Reason |
|---|----------|-------|------|------------|
| 1 | Servlet Filter | Raw HTTP access | AOP, Sidecar | Body stream + app context |
| 2 | @Async fire-and-forget | Zero latency | Sync, Direct Kafka | Keep library lightweight |
| 3 | ThreadPool 6/10/100 | Traffic-based | Unbounded | Bounded memory + monitoring |
| 4 | Avro | Schema + 70% smaller | JSON, Protobuf | Kafka ecosystem native |
| 5 | Kafka Connect | Built-in infra | Custom consumer, Flink | 2 days vs 2 weeks |
| 6 | SMT filter (sink-side) | Loose coupling | Producer routing | Add country = new connector |
| 7 | Parquet | 90% compression | JSON, ORC, BQ direct | Cheap storage + native BQ |
| 8 | Active/Active | 15-min failover | Active/Passive | Compliance requirement |
| 9 | ContentCachingWrapper | Well-tested | Custom buffering | Edge cases handled |
| 10 | isResponseLoggingEnabled | Per-team config | Hardcoded | Enabled 3-team adoption |

---

*When they ask "Why X?", answer: "We considered [alternatives]. We chose [X] because [reason]. The trade-off is [downside]. We mitigated by [action]."*

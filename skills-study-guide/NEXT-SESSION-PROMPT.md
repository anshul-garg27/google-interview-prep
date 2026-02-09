# Prompt for Next Session — Create Remaining Study Guides

Copy-paste this entire prompt into a new Claude Code session:

---

## Context

I'm Anshul Garg, a Backend Software Engineer at Walmart (SDE-III). My resume and interview prep materials are in `/Users/anshullkgarg/Desktop/work_ex/`. I already have 8 study guides in `/Users/anshullkgarg/Desktop/work_ex/skills-study-guide/`:

1. `01-kafka-avro-schema-registry.md` — Kafka, Avro, Schema Registry
2. `02-grpc-protobuf.md` — gRPC, Protobuf, HTTP/2
3. `03-system-design-distributed-systems.md` — System Design, Distributed Systems, EDA, Microservices
4. `04-istio-kubernetes-docker.md` — Kubernetes, Docker, Istio, Flagger
5. `05-aws-services.md` — S3, Lambda, SQS, Kinesis, EC2
6. `06-clickhouse-redis-rabbitmq.md` — ClickHouse, Redis, RabbitMQ
7. `07-go-language-essentials.md` — Go concurrency, goroutines, channels
8. `08-observability-data-pipeline.md` — Prometheus, Grafana, Dynatrace, ELK, Airflow, DBT, BigQuery, Parquet

## Task

Create these remaining study guides in the SAME folder (`/Users/anshullkgarg/Desktop/work_ex/skills-study-guide/`). Read my resume at `/Users/anshullkgarg/Desktop/work_ex/anshul_garg_resume_v3.tex` first to understand my experience. Also read the existing study guides (especially `01-kafka-avro-schema-registry.md`) to match the style and format.

### Documents to Create (run ALL in parallel using background agents):

**09 - Java + Java 17 Deep Dive** (`09-java-java17-deep-dive.md`)
- Java 17 features: records, sealed classes, pattern matching, text blocks, switch expressions
- CompletableFuture deep dive (used in Walmart Active/Active failover): thenApply, thenCompose, exceptionally, allOf, anyOf, timeout
- @Async + thread pool configuration in Spring
- Java concurrency: ExecutorService, ThreadPoolExecutor, ForkJoinPool, virtual threads (Java 21 preview)
- Java collections internals: HashMap, ConcurrentHashMap, ArrayList vs LinkedList
- Java memory model: heap, stack, GC algorithms (G1, ZGC, Shenandoah)
- Generics, type erasure, wildcards
- Streams API, Optional
- How Anshul used it: CompletableFuture failover at Walmart, Spring Boot starter JAR, ContentCachingWrapper
- Interview Q&A: 20+ questions
- Mermaid diagrams

**10 - Spring Boot 3 Deep Dive** (`10-spring-boot3-deep-dive.md`)
- Spring Boot auto-configuration: how it works, @Conditional annotations
- Creating custom Spring Boot starters (how Anshul built the audit JAR)
- Spring Boot 2.7 → 3.2 migration: Jakarta EE namespace change, removed APIs, property changes
- ContentCachingWrapper: how it works, request/response body capture
- @Async thread pool: configuration, TaskExecutor, rejection policies
- Spring Security basics
- Spring Data JPA / Hibernate: lazy loading, N+1 problem, batch fetching
- Actuator endpoints, health checks
- Spring profiles, configuration management
- How Anshul used it: audit logging starter JAR, ContentCachingWrapper, @Async, Spring Boot migration
- Interview Q&A: 20+ questions
- Mermaid diagrams

**11 - API Design + Kafka SMT + DR Concepts** (`11-api-design-smt-dr-concepts.md`)
- OpenAPI Specification: what it is, YAML/JSON format, design-first vs code-first
- How to write an OpenAPI spec (like Anshul's 898-line spec)
- API versioning strategies
- REST API best practices: naming, pagination, error handling, HATEOAS
- Kafka SMT (Single Message Transforms): what they are, built-in SMTs, custom SMTs
- SMT-based geographic routing (how Anshul routed US/CA/MX)
- Disaster Recovery concepts: RTO, RPO, DR strategies (Active/Active, Active/Passive, Pilot Light, Warm Standby)
- How Anshul achieved 15-min DR recovery
- P99/P95/P50 latency: what they mean, how to measure, how to optimize
- HTTP status codes: 2xx, 3xx, 4xx, 5xx with common codes
- GCS (Google Cloud Storage) basics
- Interview Q&A: 15+ questions
- Mermaid diagrams

**12 - SQL + PostgreSQL Internals** (`12-sql-postgresql-internals.md`)
- SQL fundamentals: JOINs, subqueries, CTEs, window functions, GROUP BY
- PostgreSQL architecture: processes, shared memory, WAL
- MVCC (Multi-Version Concurrency Control)
- Indexing: B-tree, Hash, GIN, GiST, BRIN — when to use which
- Query plans: EXPLAIN ANALYZE, sequential scan vs index scan
- Vacuum and autovacuum
- Connection pooling (PgBouncer, HikariCP)
- Partitioning (range, list, hash)
- ACID properties, transaction isolation levels
- PostgreSQL vs MySQL comparison
- How Anshul used it: GCC dual-database (PostgreSQL + ClickHouse), connection pooling, write saturation problem
- Interview Q&A: 20+ questions

**13 - Behavioral / Googleyness Interview Prep** (`13-behavioral-googleyness-prep.md`)
- Google's Googleyness criteria: what they evaluate
- STAR method (Situation, Task, Action, Result) for each story
- Map each resume bullet to a behavioral story:
  - Leadership: Led Spring Boot migration, Managed Cultural Fest
  - Collaboration: Audit JAR adopted by 12+ teams, cross-team notification system
  - Conflict resolution: Schema compatibility issues, dependency conflicts
  - Failure/learning: Production debugging, migration rollback scenarios
  - Mentoring: TA for 50+ students
  - Initiative: Built reusable JAR without being asked, design-first OpenAPI approach
- "Tell me about yourself" — 2-minute pitch using resume
- "Why Google?" — framework answer
- "Why leave Walmart?" — framework answer
- Common behavioral questions (20+) with STAR answers mapped to Anshul's experience
- Read `/Users/anshullkgarg/Desktop/work_ex/BEHAVIORAL-STORIES-STAR.md` and `/Users/anshullkgarg/Desktop/work_ex/INTERVIEW-PLAYBOOK-HOW-TO-SPEAK.md` for existing stories

**14 - DSA / Coding Interview Patterns** (`14-dsa-coding-patterns.md`)
- Top 15 coding patterns for Google interviews:
  1. Two Pointers
  2. Sliding Window
  3. Binary Search
  4. BFS/DFS (Trees + Graphs)
  5. Dynamic Programming (1D, 2D, knapsack)
  6. Backtracking
  7. Topological Sort
  8. Union Find
  9. Trie
  10. Monotonic Stack/Queue
  11. Heap/Priority Queue
  12. Intervals
  13. Linked List manipulation
  14. Prefix Sum
  15. Bit Manipulation
- For each pattern: explanation, template code (Python + Java), 3-5 example problems, time/space complexity
- Google-specific tips: communicate clearly, optimize step by step, test edge cases
- Complexity analysis cheat sheet: Big O for all data structures
- Common data structures: Array, LinkedList, Stack, Queue, HashMap, TreeMap, Heap, Trie, Graph
- 50 must-practice LeetCode problems (categorized by pattern)
- How to approach an unknown problem in 45 minutes

**15 - Python + FastAPI + AsyncIO** (`15-python-fastapi-asyncio.md`)
- Python async/await: event loop, coroutines, tasks
- AsyncIO internals: how the event loop works
- aiohttp, aioboto3, aio-pika (async libraries Anshul used)
- FastAPI: routing, dependency injection, Pydantic models, middleware
- FastAPI vs Django vs Flask comparison
- Python concurrency: threading vs multiprocessing vs asyncio
- uvloop (used in GCC Beat scraper)
- Connection pooling in async Python
- Rate limiting implementation
- How Anshul used it: Beat scraper (150+ async workers), FastAPI services, async S3 uploads
- Interview Q&A: 15+ questions

Each document should include:
- Deep explanations with examples
- Mermaid diagrams (architecture, flows, comparisons)
- Interview Q&A section
- "How Anshul Used It" section connecting to resume bullets
- Comparison tables where applicable
- Code examples where relevant
- No length limit — be as comprehensive as possible

Run all 7 agents in parallel using background mode. Use `fullstack-developer` subagent type with `bypassPermissions` mode.

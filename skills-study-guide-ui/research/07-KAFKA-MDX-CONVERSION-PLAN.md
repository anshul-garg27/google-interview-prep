# Kafka Guide MDX Conversion Plan

## Guide Stats
- **2,400 lines** | **14,015 words** | **53 code blocks** | **10 mermaid diagrams** | **25 Q&A** | **151 table rows**
- 5 major parts: Kafka Core → Avro → Schema Registry → Interview Q&A → How Anshul Used It

---

## Section-by-Section Plan

### PART 1: APACHE KAFKA (Lines 42-943)

#### 1.1 What Is Kafka (Lines 44-70)
- **Current:** Text explaining what Kafka is
- **Plan:** Keep as markdown prose. Add a `<KeyConcept>` callout for the "Key properties" section
- **Interactive:** None needed — this is foundational text
- **Format:** 95% markdown

#### 1.2 Core Architecture (Lines 71-160)
- **Current:** Mermaid diagram + text explaining Broker, Topic, Partition, Producer, Consumer
- **Plan:** Replace mermaid with `<KafkaClusterDiagram>`
  - Clickable brokers → tooltip shows which partitions it leads/follows
  - Clickable partitions → show message offset range
  - Color: Leaders in amber, Followers in gray
  - Animated message dots flowing from Producer → Partition → Consumer
- **Interactive:** HIGH — this is the anchor diagram of the entire guide
- **Format:** 40% markdown (explanations) + 60% interactive component

#### 1.3 How Messages Flow (Lines 162-252)
- **Current:** Mermaid sequence diagram (Producer → Broker → Consumer)
- **Plan:** Replace with `<MessageFlowAnimation>`
  - Step-by-step: click "Next" to see each step
  - Step 1: Producer calls send(topic, key, value)
  - Step 2: Partitioner determines partition (show hash(key) % N calculation)
  - Step 3: ProduceRequest to leader broker (animate arrow)
  - Step 4: Followers replicate (FetchRequest)
  - Step 5: ProduceResponse back
  - Step 6: Consumer JoinGroup + heartbeat
  - Step 7: Consumer FetchRequest
  - Step 8: Consumer processes + commits offset
- **Interactive:** HIGH — this is the most important flow to understand
- **Format:** 30% markdown + 70% step-by-step animation

#### 1.4 Kafka Internals (Lines 253-363)
- **Current:** Text + ASCII art for commit log, offsets, replication, ISR
- **Plan:**
  - `<CommitLogVisualizer>` — show partition as append-only log, messages adding
  - Replace ASCII art with proper SVG of offset numbering
  - ISR explanation: `<ISRDiagram>` showing leader + followers syncing
  - Keep the text explanations as markdown
- **Interactive:** MEDIUM — commit log visual is key, rest is text
- **Format:** 60% markdown + 40% interactive

#### 1.5 Partitioning Strategies (Lines 364-430)
- **Current:** Text + Java code for custom partitioner
- **Plan:** `<PartitioningVisualizer>`
  - Enter message keys (e.g., "user-1", "user-2", "user-3")
  - See which partition each lands in
  - Toggle between strategies: Hash / Round-Robin / Sticky / Custom
  - Show the actual hash calculation
- **Interactive:** HIGH — lets you "feel" how partitioning works
- **Format:** 40% markdown + 60% interactive
- **Code:** Keep Java custom partitioner with `<AnnotatedCodeBlock>` (hover explanations per line)

#### 1.6 Delivery Semantics (Lines 431-524)
- **Current:** Text with numbered step sequences for each guarantee
- **Plan:** `<DeliverySemanticsSimulator>`
  - Three columns: At-Most-Once | At-Least-Once | Exactly-Once
  - Each column shows: Producer → Broker → Consumer flow
  - Click "Inject Failure" at any step to see what happens
  - Show message loss (red X), duplication (yellow dup), or exactly-once (green check)
  - Toggle acks=0, acks=1, acks=all and see behavior change
- **Interactive:** VERY HIGH — this is the #1 interview question topic
- **Format:** 20% markdown + 80% interactive simulator

#### 1.7 Consumer Groups & Rebalancing (Lines 525-624)
- **Current:** Mermaid sequence diagram + text + table
- **Plan:** `<ConsumerGroupSimulator>` (already built in demo!)
  - Add/remove consumers
  - Toggle strategy (Range, RoundRobin, CooperativeSticky)
  - "Kill" a consumer to see rebalancing animate
  - Show the actual JoinGroup → SyncGroup protocol step-by-step
- **Interactive:** VERY HIGH — already proven in demo, students love it
- **Format:** 30% markdown + 70% interactive

#### 1.8 Kafka Connect & Streams (Lines 625-711)
- **Current:** Text + code blocks
- **Plan:** Keep mostly as markdown with `<AnnotatedCodeBlock>` for Kafka Streams code
- Add a `<ConnectorFlowDiagram>` — simple animated Source → Kafka → Sink flow
- **Interactive:** LOW — this section is reference material
- **Format:** 80% markdown + 20% interactive

#### 1.9 Multi-Region Active/Active (Lines 712-819)
- **Current:** Mermaid diagram + Java code + comparison table
- **Plan:** `<MultiRegionArchitecture>`
  - Two-region diagram with animated traffic flow
  - Click "Fail Region EUS2" → see traffic reroute to SCUS
  - Show the CompletableFuture failover code with `<AnnotatedCodeBlock>`
  - Decision table as interactive comparison
- **Interactive:** HIGH — this is Anshul's #1 architectural achievement
- **Format:** 30% markdown + 70% interactive

#### 1.10 CompletableFuture Failover (Lines 820-871)
- **Current:** Java code + explanation
- **Plan:** `<CodeWalkthrough>` — step-by-step through the failover code
  - Each step highlights lines + shows explanation panel
  - "What happens when primary fails?" → animate the exception → secondary call
- **Interactive:** MEDIUM
- **Format:** 40% markdown + 60% code walkthrough

#### 1.11 Performance Tuning (Lines 872-943)
- **Current:** Tables with config parameters
- **Plan:** `<KafkaTuningDashboard>`
  - Sliders for: batch.size, linger.ms, compression.type, acks
  - Live chart showing: Throughput (msgs/sec) and Latency (ms)
  - Presets: "Low Latency", "High Throughput", "Balanced", "Walmart Audit"
- **Interactive:** VERY HIGH — makes abstract configs tangible
- **Format:** 20% markdown + 80% interactive

---

### PART 2: APACHE AVRO (Lines 944-1419)

#### 2.1 What Is Avro + Comparison (Lines 946-989)
- **Current:** Text + comparison table (Avro vs JSON vs Protobuf vs Thrift)
- **Plan:** `<ComparisonToggle>` — interactive comparison
  - Click columns to highlight
  - Toggle "Best for..." filters
  - Show message size comparison with actual byte counts
- **Interactive:** MEDIUM
- **Format:** 60% markdown + 40% interactive table

#### 2.2 Schema Definition (Lines 990-1135)
- **Current:** JSON examples of Avro schemas + code
- **Plan:** `<AvroSchemaPlayground>`
  - Live Avro schema editor (left panel)
  - Generated sample data preview (right panel)
  - Validate button → shows if schema is valid
  - Show type system visually (primitive → complex)
- **Interactive:** HIGH — schema definition is confusing in text
- **Format:** 40% markdown + 60% interactive

#### 2.3 Schema Evolution (Lines 1136-1278)
- **Current:** Text with examples of adding/removing fields
- **Plan:** `<SchemaEvolutionDemo>`
  - Two-panel: Schema V1 (left) → Schema V2 (right)
  - Highlight what changed (green for added, red for removed)
  - Show: "Old reader reads new data" and "New reader reads old data"
  - Toggle compatibility mode and see if it passes/fails
- **Interactive:** HIGH — evolution rules are the hardest part of Avro
- **Format:** 30% markdown + 70% interactive

#### 2.4 Code Examples (Lines 1279-1419)
- **Current:** Java producer/consumer code with Avro
- **Plan:** `<AnnotatedCodeBlock>` with hover explanations
- **Format:** 70% code + 30% annotations

---

### PART 3: CONFLUENT SCHEMA REGISTRY (Lines 1420-1796)

#### 3.1 Architecture + How It Works (Lines 1422-1536)
- **Current:** Text + mermaid diagrams
- **Plan:** `<SchemaRegistryFlow>` — animated flow showing:
  - Producer → Schema Registry (register) → get schema ID → attach to message
  - Consumer → Schema Registry (lookup) → deserialize
  - Show the magic byte + schema ID wire format
- **Interactive:** MEDIUM
- **Format:** 50% markdown + 50% interactive

#### 3.2 REST API (Lines 1604-1670)
- **Current:** curl commands
- **Plan:** `<APIExplorer>` — interactive REST API playground
  - Show each endpoint with description
  - Click "Try it" → show the curl command + response
  - Not actual API call — just simulated request/response
- **Interactive:** MEDIUM
- **Format:** 30% markdown + 70% interactive explorer

#### 3.3 Compatibility Modes (Lines 1537-1603)
- **Current:** Text + table
- **Plan:** Keep as markdown with `<ComparisonToggle>` for the modes table
- **Format:** 70% markdown + 30% interactive

---

### PART 4: INTERVIEW Q&A (Lines 1797-2221) — 25 Questions

- **Plan:** `<FlashcardDeck>` component
  - All 25 questions as cards
  - Click to flip → reveal answer
  - Difficulty rating: Again / Hard / Good / Easy
  - Spaced repetition scheduling (SM-2 algorithm in localStorage)
  - Toggle: "Study Mode" (all visible) vs "Quiz Mode" (one at a time)
  - Progress: "15/25 reviewed"
  - Filter by subtopic: Kafka Core, Avro, Schema Registry, Anshul-specific
- **Interactive:** VERY HIGH — this is THE study tool feature
- **Format:** 10% markdown (section header) + 90% interactive flashcards

---

### PART 5: HOW ANSHUL USED IT AT WALMART (Lines 2222-2346)

#### 5.1 System Overview (Lines 2224-2345)
- **Plan:** `<WalmartAuditArchitecture>` — the crown jewel
  - Full interactive 3-tier architecture diagram
  - Click each tier to expand and see details
  - Animated data flow: Request → LoggingFilter → Kafka → GCS → BigQuery
  - Before/After comparison table with animated counters
  - Timeline component showing PRs and milestones
  - "Try the flow" — trace a single audit event through the system
- **Interactive:** VERY HIGH — this tells Anshul's story
- **Format:** 20% markdown + 80% interactive

#### 5.2 Quick Reference Card (Lines 2347-2400)
- **Plan:** `<CommandReference>` — copy-friendly command cards
  - Each command in its own card with copy button
  - Categorized: Topic, Consumer Group, Schema Registry
  - Search/filter commands
- **Interactive:** MEDIUM
- **Format:** 20% markdown + 80% interactive reference

---

## Interactive Components to Build

| # | Component | Used In | Effort | Priority |
|---|-----------|---------|--------|----------|
| 1 | `<KafkaClusterDiagram>` | 1.2 Core Architecture | 3-4 days | P0 |
| 2 | `<MessageFlowAnimation>` | 1.3 How Messages Flow | 3 days | P0 |
| 3 | `<ConsumerGroupSimulator>` | 1.7 Consumer Groups | 2 days (exists in demo) | P0 |
| 4 | `<DeliverySemanticsSimulator>` | 1.6 Delivery Semantics | 3 days | P0 |
| 5 | `<FlashcardDeck>` | Part 4 Q&A | 3 days | P0 |
| 6 | `<AnnotatedCodeBlock>` | 1.5, 1.9, 1.10, 2.4 | 2 days (exists in demo) | P0 |
| 7 | `<KafkaTuningDashboard>` | 1.11 Performance Tuning | 2-3 days | P1 |
| 8 | `<PartitioningVisualizer>` | 1.5 Partitioning Strategies | 2 days | P1 |
| 9 | `<SchemaEvolutionDemo>` | 2.3 Schema Evolution | 2-3 days | P1 |
| 10 | `<WalmartAuditArchitecture>` | Part 5 System Overview | 3-4 days | P1 |
| 11 | `<ComparisonToggle>` | 2.1, 3.3, various | 1 day | P1 |
| 12 | `<SchemaRegistryFlow>` | 3.1 Architecture | 2 days | P2 |
| 13 | `<AvroSchemaPlayground>` | 2.2 Schema Definition | 2-3 days | P2 |
| 14 | `<APIExplorer>` | 3.2 REST API | 2 days | P2 |
| 15 | `<CommitLogVisualizer>` | 1.4 Kafka Internals | 1-2 days | P2 |
| 16 | `<CommandReference>` | 5.2 Quick Reference | 1 day | P2 |
| 17 | `<CodeWalkthrough>` | 1.10 CompletableFuture | 1 day (exists in demo) | P2 |

**Total: 17 components, ~35-40 days for all**

---

## Phased Build Plan

### Phase 1: Foundation + Core (Week 1-2)
Set up MDX pipeline + build the 6 P0 components

1. Set up Vite MDX plugin (@mdx-js/rollup)
2. Convert kafka guide .md → .mdx
3. Build `<ConsumerGroupSimulator>` (port from demo)
4. Build `<AnnotatedCodeBlock>` (port from demo)
5. Build `<FlashcardDeck>` with SM-2 spaced repetition
6. Build `<DeliverySemanticsSimulator>`
7. Build `<MessageFlowAnimation>`
8. Build `<KafkaClusterDiagram>`

**Result:** Kafka guide already feels 5x better than raw markdown

### Phase 2: Polish + P1 Components (Week 3-4)
7. Build `<KafkaTuningDashboard>` with sliders
8. Build `<PartitioningVisualizer>`
9. Build `<SchemaEvolutionDemo>`
10. Build `<WalmartAuditArchitecture>`
11. Build `<ComparisonToggle>`

**Result:** Kafka guide is now world-class

### Phase 3: P2 + Refinement (Week 5-6)
12-17. Remaining components
- Mobile optimization
- Dark mode testing
- Performance optimization

**Result:** Complete Kafka interactive guide

---

## Content Split Estimate

| Section | Markdown % | Interactive % | Lines of React |
|---------|-----------|---------------|----------------|
| Part 1: Kafka Core | 35% | 65% | ~3,500 |
| Part 2: Avro | 50% | 50% | ~1,500 |
| Part 3: Schema Registry | 55% | 45% | ~1,200 |
| Part 4: Q&A (25 items) | 10% | 90% | ~800 |
| Part 5: Anshul's Story | 20% | 80% | ~1,500 |
| **Overall** | **~35%** | **~65%** | **~8,500** |

This is MORE interactive than the original "90% markdown" MDX idea — closer to 35/65 split.
But the key difference from Approach 4 (pure React) is: the markdown prose stays as markdown,
the interactive parts are surgically placed React components.

---

## What Stays as Markdown (35%)
- All paragraph explanations
- Bullet point lists
- Simple tables (rebalance triggers, technologies used)
- Inline code references
- Section headings

## What Becomes Interactive React (65%)
- Architecture diagrams → clickable SVG
- Sequence diagrams → step-by-step animations
- Code blocks → annotated with hover tooltips
- Consumer groups → drag/simulate
- Delivery semantics → failure injection simulator
- Config tuning → slider dashboard
- Q&A → flashcard mode
- Walmart story → interactive architecture walkthrough

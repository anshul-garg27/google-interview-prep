# 05 - Kafka Guide Interactive Prototype: Component-by-Component Design Specification

**Agent**: 5 - Prototype Architect
**Date**: 2026-02-08
**Source Guide**: `/skills-study-guide/01-kafka-avro-schema-registry.md` (2,400 lines, 5 Parts, 25 Q&A)
**Tech Stack**: React 19, Tailwind CSS 3, Vite 7, react-syntax-highlighter, lucide-react, mermaid (existing)
**Goal**: Transform every section of Guide 01 into a fully interactive, hand-crafted learning experience

---

## TABLE OF CONTENTS

1. [Guide Structure Analysis](#1-guide-structure-analysis)
2. [Section-by-Section Interactive Designs](#2-section-by-section-interactive-designs)
   - 2.1 Core Architecture Diagram
   - 2.2 How Messages Flow
   - 2.3 Kafka Internals (Commit Log + Offsets + ISR)
   - 2.4 Partitioning Strategies
   - 2.5 Delivery Semantics
   - 2.6 Consumer Groups and Rebalancing
   - 2.7 Kafka Connect and Kafka Streams
   - 2.8 Multi-Region Active/Active
   - 2.9 CompletableFuture Failover Pattern
   - 2.10 Performance Tuning Tables
   - 2.11 Avro Schema + Evolution
   - 2.12 Schema Registry Architecture
   - 2.13 Compatibility Matrix
   - 2.14 Avro Serialization Flow
   - 2.15 Code Examples (Producer/Consumer)
   - 2.16 Interview Q&A (25 Questions)
   - 2.17 Walmart System Overview
   - 2.18 Comparison Tables (Kafka vs RabbitMQ vs SQS, Avro vs JSON vs Protobuf)
3. [Reusable Component Library](#3-reusable-component-library)
4. [Effort Estimates](#4-effort-estimates)
5. [Dependency Map](#5-dependency-map)
6. [Implementation Priority](#6-implementation-priority)

---

## 1. GUIDE STRUCTURE ANALYSIS

The Kafka guide has these content types that need interactive treatment:

| Content Type | Count in Guide | Current Rendering | Interactive Treatment |
|-------------|---------------|-------------------|----------------------|
| Mermaid architecture diagrams | 8 | Static SVG via mermaid lib | Custom SVG with hover/click/animate |
| Mermaid sequence diagrams | 4 | Static SVG via mermaid lib | Step-by-step animated walkthrough |
| ASCII art diagrams | 4 | Monospace code block | Animated SVG with tooltips |
| Comparison tables | 6 | Pipe-delimited markdown | Sortable, hoverable, with decision wizard |
| Configuration tables | 3 | Pipe-delimited markdown | Interactive tuning sliders |
| Java code blocks | 12 | Syntax highlighted via Prism | Annotated, line-by-line walkthrough |
| Bash code blocks | 4 | Syntax highlighted | Copy-to-clipboard (already done) |
| Interview Q&A pairs | 25 | Flat text | Flashcard deck with spaced repetition |
| Bulleted explanations | ~30 | Rendered markdown | Progressive disclosure + tooltips |
| Inline diagrams (wire format) | 2 | ASCII in code block | Interactive byte inspector |

---

## 2. SECTION-BY-SECTION INTERACTIVE DESIGNS

---

### 2.1 CORE ARCHITECTURE DIAGRAM

**Current content** (lines 77-124): A mermaid `graph TB` diagram showing 3 brokers, each containing topic-partition replicas (leader/follower), 2 producers, 2 consumer groups, and ZooKeeper/KRaft.

**Interactive Version: `<KafkaClusterDiagram />`**

#### Visual Layout (ASCII Sketch)

```
+------------------------------------------------------------------+
|  KAFKA CLUSTER EXPLORER                          [Reset] [Zoom]  |
+------------------------------------------------------------------+
|                                                                  |
|  +-----------+    +-----------+    +-----------+                 |
|  | BROKER 1  |    | BROKER 2  |    | BROKER 3  |                 |
|  | id: 0     |    | id: 1     |    | id: 2     |                 |
|  |           |    |           |    |           |                 |
|  | [P0-L]    |    | [P0-F]    |    | [P0-F]    |  <-- Topic-A   |
|  | [P1-F]    |    | [P1-L]    |    | [P1-F]    |                 |
|  | [TB:P0-F] |    | [TB:P0-L] |    | [TB:P0-F] |  <-- Topic-B   |
|  +-----------+    +-----------+    +-----------+                 |
|        ^                ^                ^                       |
|        |    ............|................|....   Replication      |
|        |    :           |               |   :                    |
|   [Producer 1]    [Producer 2]              :                    |
|        |                |                   :                    |
|        v                v                   :                    |
|   +---------+    +---------+                :                    |
|   |Consumer |    |Consumer |    +---------+ :                    |
|   |Group A  |    |Group A  |    |Consumer | :                    |
|   | C1: P0  |    | C2: P1  |    |Group B  | :                    |
|   +---------+    +---------+    | C3: ALL | :                    |
|                                 +---------+ :                    |
|   ..............................................                 |
|   : ZooKeeper/KRaft Controller :                                 |
|   :............................:                                 |
|                                                                  |
|  [+ Add Broker]  [- Remove Broker]  [+ Add Partition]            |
|  [Kill Broker ▼]  [Trigger Rebalance]                            |
+------------------------------------------------------------------+
|  INFO PANEL (appears on hover/click)                             |
|  Broker 1: Leader for Topic-A/P0, Topic-B/P0(F). ISR: healthy   |
+------------------------------------------------------------------+
```

#### Hover Behaviors

- **Hover over a Broker box**:
  - Broker box gets a 2px amber glow border
  - Tooltip appears below showing:
    ```
    Broker 1 (id: 0)
    --------------------------------
    Leader for:  Topic-A/Partition-0
    Follower for: Topic-A/Partition-1, Topic-B/Partition-0
    ISR status: All partitions in-sync
    Disk: 2.3 GB used
    Network: 45 MB/s in, 120 MB/s out
    ```
  - All partition boxes on this broker get highlighted

- **Hover over a Partition box** (e.g., `[P0-L]`):
  - The partition box pulses with amber fill (leader) or gray fill (follower)
  - All replicas of this partition across other brokers get a dashed highlight
  - Arrow lines appear showing replication direction (leader -> followers)
  - Tooltip shows:
    ```
    Topic-A / Partition 0
    Role: LEADER
    Replicas: Broker 1 (leader), Broker 2, Broker 3
    ISR: {1, 2, 3}
    Log End Offset: 1,247,893
    High Watermark: 1,247,890
    Size: 412 MB
    ```

- **Hover over Producer**:
  - Animated dashed line from producer to the leader partition it targets
  - Line pulses with small dots moving in the send direction
  - Tooltip: `Producer 1 -> Topic-A/Partition-0 (key hash: murmur2)`

- **Hover over Consumer**:
  - Solid line from consumer to its assigned partition(s)
  - Tooltip: `Consumer 1 (Group A) -> Partition 0 | Lag: 23 messages | Offset: 1,247,870`

#### Click Behaviors

- **Click a Partition box**: Expands inline to show a mini commit log:
  ```
  +---------------------------------------------------+
  | Topic-A / Partition 0 -- Commit Log               |
  +---------------------------------------------------+
  | Offset | Key                        | Timestamp   |
  |--------|----------------------------|-------------|
  | 1247888| cp-nrti-apis|/iac/v1/inv   | 14:23:01.3  |
  | 1247889| supplier-svc|/v2/orders    | 14:23:01.5  |
  | 1247890| cp-nrti-apis|/iac/v1/inv   | 14:23:01.8  |  <-- HW
  | 1247891| cp-nrti-apis|/v1/items     | 14:23:02.1  |
  | 1247892| supplier-svc|/v2/orders    | 14:23:02.3  |
  | 1247893| cp-nrti-apis|/iac/v1/inv   | 14:23:02.5  |  <-- LEO
  +---------------------------------------------------+
  | [Produce Message]  [Scroll Log]  [Close]          |
  +---------------------------------------------------+
  ```

- **Click "Kill Broker" dropdown -> select Broker 2**:
  1. Broker 2 box turns red with an X overlay (500ms fade)
  2. Its leader partitions (Topic-A/P1, Topic-B/P0) flash red
  3. After 1 second pause, an animation shows:
     - Controller detects failure (ZK box flashes)
     - New leader election: Topic-A/P1 leader moves to Broker 3 (was follower)
     - Topic-B/P0 leader moves to Broker 1
     - Partition boxes update their L/F labels with a slide animation
     - ISR sets shrink (remove Broker 2) with tooltips updating
  4. Info panel narrates each step:
     "Step 1: Controller detects Broker 2 failure via session timeout..."
     "Step 2: Electing new leader for Topic-A/P1 from ISR {Broker 1, Broker 3}..."
     "Step 3: Broker 3 promoted to leader for Topic-A/P1..."

- **Click "+ Add Broker"**:
  1. New Broker 4 box slides in from right
  2. Currently holds no partitions (empty box)
  3. After 1 second, reassignment animation:
     - Some follower replicas move to Broker 4
     - Replication lines animate to show new topology
  4. Info: "Broker 4 added. Partition reassignment in progress..."

- **Click "+ Add Partition"**:
  1. Dialog: "Add partition to which topic? [Topic-A] [Topic-B]"
  2. New partition box appears in each broker (as followers) with one designated leader
  3. Consumer group assignments update (triggers rebalance animation)

#### Animation: Message Flow

A looping background animation (can be toggled on/off):
1. Small colored dots appear at Producer 1
2. Dots travel along the arrow to the leader partition
3. Dot enters the partition box (offset counter increments)
4. After brief pause, smaller dots travel from leader to follower replicas
5. Once replicated, a dot travels from the partition to the consumer
6. Cycle repeats every 3 seconds

#### Color Coding

| Element | Color | Tailwind Class |
|---------|-------|---------------|
| Leader partition | Amber/Gold fill | `bg-amber-100 border-amber-500` (light), `bg-amber-900/30 border-amber-400` (dark) |
| Follower partition | Gray fill | `bg-gray-100 border-gray-400` (light), `bg-gray-800 border-gray-600` (dark) |
| Healthy broker | Green left-border | `border-l-4 border-l-emerald-500` |
| Dead broker | Red overlay | `bg-red-500/20 border-red-500` |
| Producer | Blue | `bg-blue-100 border-blue-500` |
| Consumer | Purple | `bg-purple-100 border-purple-500` |
| ZK/KRaft | Dashed gray border | `border-dashed border-gray-400` |
| Replication arrows | Gray dashed | `stroke-dasharray="4 4" stroke-gray-400` |
| Message dots | Amber animated circles | SVG `<circle>` with CSS `@keyframes` |

#### Component Specification

```
<KafkaClusterDiagram>
  Props:
    initialBrokers: number (default 3)
    initialTopics: [{ name, partitions, replicationFactor }]
    showAnimation: boolean (default true)
    interactive: boolean (default true)

  State:
    brokers: [{ id, alive, partitions: [{ topic, partition, role, isr }] }]
    producers: [{ id, targetTopic }]
    consumerGroups: [{ id, consumers: [{ id, partitions }] }]
    selectedElement: { type, id } | null
    animationState: { messages: [{ id, from, to, progress }] }
    infoPanel: { title, content } | null

  Subcomponents:
    <BrokerBox broker={} onHover={} onClick={} />
    <PartitionChip partition={} isLeader={} isHighlighted={} />
    <PartitionDetailPanel partition={} messages={} />
    <ProducerNode producer={} />
    <ConsumerGroupBox group={} />
    <ReplicationArrow from={} to={} animated={} />
    <MessageDot position={} color={} />
    <InfoPanel title={} content={} />
    <ClusterControls onAddBroker={} onKillBroker={} onAddPartition={} />
```

**Lines of React code**: ~650
**Custom subcomponents**: 10
**Time to implement**: 3-4 days
**Dependencies**: None beyond existing (SVG rendering, CSS animations)

---

### 2.2 HOW MESSAGES FLOW (Step-by-Step Animation)

**Current content** (lines 166-250): A mermaid `sequenceDiagram` showing Producer -> Partitioner -> Leader -> Followers -> Consumer, plus detailed text steps 1-7.

**Interactive Version: `<MessageFlowWalkthrough />`**

#### Visual Layout

```
+------------------------------------------------------------------+
|  MESSAGE FLOW: End-to-End                                        |
|                                                                  |
|  Step [3] of 7      [< Prev]  [Next >]  [Auto-Play]  [Reset]    |
+------------------------------------------------------------------+
|                                                                  |
|  [PRODUCER]   [PARTITIONER]   [LEADER]   [F1]  [F2]   [CONSUMER]|
|     ( )          ( )            ( )       ( )   ( )      ( )     |
|      |            |              |         |     |        |      |
|      |--- send -->|              |         |     |        |      |
|      |            |-- hash(key) -|         |     |        |      |
|      |            |   = P0      [*]        |     |        |      |
|      |            |              |         |     |        |      |
|      |            |              |         |     |        |      |
|      |            |              |         |     |        |      |
|      |            |              |         |     |        |      |
|      |            |              |         |     |        |      |
|                                                                  |
|  CURRENT STEP HIGHLIGHTED (amber box around active elements)     |
|                                                                  |
+------------------------------------------------------------------+
|  STEP 3: ProduceRequest sent to Leader Broker                    |
|                                                                  |
|  The partitioner determined this message goes to Partition 0.    |
|  It sends a ProduceRequest to the leader broker for Partition 0  |
|  (Broker 1). The request includes the serialized message batch.  |
|                                                                  |
|  > The batch may include multiple messages if linger.ms > 0      |
|  > Batch is compressed with snappy/lz4 if configured             |
|                                                                  |
|  Code:                                                           |
|  ```java                                                         |
|  // Internally, the RecordAccumulator batches messages            |
|  // When batch.size reached or linger.ms expires:                 |
|  sender.sendProduceRequest(batch, leader);                        |
|  ```                                                             |
+------------------------------------------------------------------+
```

#### Step Definitions

| Step | Active Elements | Animation | Explanation |
|------|----------------|-----------|-------------|
| 1 | Producer highlighted | Producer glows, code snippet shows `producer.send(record, callback)` | "Producer calls send() with a ProducerRecord containing topic, key, value, and headers" |
| 2 | Partitioner highlighted | Hash animation: key text "cp-nrti-apis\|/iac/v1/inv" -> visual murmur2 hash wheel -> output "Partition 0" | "Partitioner computes hash(key) % numPartitions. Key 'cp-nrti-apis\|/iac/v1/inv' hashes to partition 0." |
| 3 | Arrow from Partitioner to Leader | Animated message packet (small rectangle with "ProduceRequest" label) slides from left to leader column | "Message batch sent to the leader broker for Partition 0" |
| 4 | Leader highlighted | Leader box fills with amber, internal animation shows message appended to a mini log (offset counter increments) | "Leader appends message to the partition log segment. Assigns offset 1,247,893. Written to OS page cache first." |
| 5 | Arrows from Leader to F1, F2 | Two parallel animated arrows slide from Leader to F1 and F2. Small "FetchRequest" labels on the arrows | "Follower brokers pull from the leader (pull-based replication). Each follower sends a FetchRequest." |
| 6 | F1 and F2 get green check marks | Both followers flash green. ISR indicator shows "{Leader, F1, F2}". "acks=all" badge appears. Arrow from Leader back to Producer with "ProduceResponse" label | "All ISR replicas have caught up. With acks=all, the leader sends ProduceResponse back to the producer with the assigned offset and timestamp." |
| 7 | Consumer highlighted | Consumer sends "FetchRequest" arrow to Leader. Leader returns "FetchResponse" with batch. Consumer processes and sends "OffsetCommit" | "Consumer fetches messages starting from its committed offset. After processing, commits the new offset to __consumer_offsets." |

#### Interactive Controls

- **Prev / Next buttons**: Navigate steps. Current step's active elements are highlighted with amber outline and subtle pulse animation. Previous steps' arrows/elements remain visible but dimmed (50% opacity).
- **Auto-Play toggle**: Steps advance automatically every 4 seconds with smooth transitions.
- **Speed slider**: 0.5x, 1x, 2x, 4x
- **Reset**: Returns to step 1 with all elements in initial state.

#### Configuration Toggles

Below the diagram, a row of toggle switches:
```
[acks: 0 | 1 | all]   [Compression: none | snappy | lz4]   [Idempotent: on | off]
```

Changing these affects the animation:
- `acks=0`: Step 6 is skipped (Producer gets no acknowledgement). Step 4 shows a warning: "Fire-and-forget. Message could be lost if broker crashes before writing."
- `acks=1`: Step 5 is skipped (no waiting for follower replication). Step 6 shows leader-only acknowledgement.
- `acks=all`: Full flow as described above.
- `Compression: snappy`: Step 3 shows a "compress" icon on the message packet, and size indicator shrinks visually.
- `Idempotent: on`: Step 3 shows "PID=7, Seq=1247893" label on the message.

**Lines of React code**: ~450
**Custom subcomponents**: 7 (StepNavigator, SequenceColumn, AnimatedArrow, MessagePacket, StepExplanation, ConfigToggleRow, MiniLogView)
**Time to implement**: 2-3 days
**Dependencies**: framer-motion (optional, for smooth step transitions -- could use CSS transitions instead)

---

### 2.3 KAFKA INTERNALS: COMMIT LOG + OFFSETS + ISR

**Current content** (lines 253-361): ASCII art commit log, offset range diagram, ISR dynamics text, leader election steps.

#### 2.3a Commit Log Visualizer: `<CommitLogExplorer />`

```
+------------------------------------------------------------------+
|  COMMIT LOG: Partition 0                                         |
+------------------------------------------------------------------+
|                                                                  |
|  Log Start    Consumer A    Consumer B    High        Log End    |
|  Offset       Committed     Committed     Watermark   Offset     |
|    |              |              |            |          |        |
|    v              v              v            v          v        |
|  +----+----+----+----+----+----+----+----+----+----+----+----+   |
|  | 0  | 1  | 2  | 3  | 4  | 5  | 6  | 7  | 8  | 9  | 10 | 11|  |
|  |msg |msg |msg |msg |msg |msg |msg |msg |msg |msg |msg |msg |  |
|  +----+----+----+----+----+----+----+----+----+----+----+----+   |
|  |<-- consumed -->|         |<- available ->|<- not visible ->|  |
|  | (can delete)   |         |  to consumers |  to consumers   |  |
|                                                                  |
|  [Produce Message]  [Consume (A)]  [Consume (B)]  [Compact]     |
+------------------------------------------------------------------+
|  OFFSET INSPECTOR (click any offset cell)                        |
|  Offset 5:                                                       |
|    Key: "cp-nrti-apis|/iac/v1/inventory"                         |
|    Value: {"request_id":"abc-123","service_name":"cp-nrti..."}   |
|    Timestamp: 2025-05-15T14:23:01.500Z                           |
|    Headers: {"wm-site-id":"US","wm_consumer.id":"supp-789"}     |
|    Size: 847 bytes                                               |
+------------------------------------------------------------------+
```

**Interactions**:
- **Click an offset cell**: Expands below to show the full message details (key, value, headers, timestamp).
- **"Produce Message" button**: Appends a new cell at the end with a slide-in animation. LEO pointer moves right.
- **"Consume (A)" button**: Consumer A's pointer advances. If it catches up to HW, a "waiting for new messages" indicator appears.
- **"Consume (B)" button**: Consumer B's pointer advances independently.
- **"Compact" button**: Demonstrates log compaction -- duplicate keys are removed, only latest value per key remains. Cells with old values fade out and slide together.
- **Drag the offset pointers**: Users can drag Consumer A, Consumer B, HW, and Log Start Offset markers to see how different positions affect the available range.

**Color coding for offset cells**:
- Green fill: Consumed by all consumers and before Log Start Offset (eligible for deletion)
- Amber fill: Available for consumption (between committed offset and HW)
- Gray fill: Beyond HW (not visible to consumers, pending replication)

**Lines of React code**: ~350
**Custom subcomponents**: 5 (OffsetCell, PointerMarker, MessageInspector, LogSegmentBar, RangeLabel)
**Time to implement**: 2 days

#### 2.3b ISR Dynamics Simulator: `<ISRSimulator />`

```
+------------------------------------------------------------------+
|  ISR DYNAMICS SIMULATOR                                          |
+------------------------------------------------------------------+
|  Time: [T0] [T1] [T2] [T3]      [Auto-Play] [Reset]            |
|                                                                  |
|  Broker 1 (Leader)    [=============================] LEO: 1000  |
|  Broker 2 (Follower)  [=============================] LEO: 1000  |
|  Broker 3 (Follower)  [=============================] LEO: 1000  |
|                                                                  |
|  ISR: {Broker 1, Broker 2, Broker 3}    min.insync.replicas: 2  |
|  Status: ALL HEALTHY                                             |
+------------------------------------------------------------------+
```

**Time T0**: All brokers caught up. All bars full and green.
**Time T1 (click)**: Broker 3 bar shrinks (falls behind). After `replica.lag.time.max.ms`, Broker 3 turns red and is removed from ISR. ISR display updates with animation.
**Time T2 (click)**: Broker 3 bar grows back (catching up). Turns green and re-enters ISR.
**Time T3 (click)**: Broker 1 (Leader) dies. Bar turns red with X. Controller animation: Broker 2 promoted to Leader. ISR = {Broker 2, Broker 3}.

**Interactive controls**: Users can click individual brokers to "kill" or "slow down" them and see the ISR response.

**Lines of React code**: ~250
**Subcomponents**: 4 (ReplicaBar, ISRBadge, TimelineSelector, StatusPanel)
**Time**: 1.5 days

---

### 2.4 PARTITIONING STRATEGIES

**Current content** (lines 364-427): Key-based, round-robin, sticky, custom partitioning with code examples.

**Interactive Version: `<PartitioningVisualizer />`**

```
+------------------------------------------------------------------+
|  PARTITIONING STRATEGY EXPLORER                                  |
|                                                                  |
|  Strategy: [Key-Based*] [Round-Robin] [Sticky] [Custom-Geo]     |
|  Partitions: [4] [6*] [12]                                      |
+------------------------------------------------------------------+
|                                                                  |
|  MESSAGE INPUT:                                                  |
|  Key: [cp-nrti-apis|/iac/v1/inventory  ]                        |
|  Value: [{ "request_id": "abc-123" }   ]                        |
|  [Send Message]  [Send 20 Random]                                |
|                                                                  |
|  HASH CALCULATION (visible for key-based):                       |
|  murmur2("cp-nrti-apis|/iac/v1/inventory") = 0x7A3F21B4         |
|  toPositive(0x7A3F21B4) = 2050531764                             |
|  2050531764 % 6 = 4                                              |
|  --> Partition 4                                                  |
|                                                                  |
+------------------------------------------------------------------+
|  PARTITION DISTRIBUTION:                                         |
|                                                                  |
|  P0: [||||||||          ] 8 msgs   35%  keys: svc-A|/orders     |
|  P1: [||                ] 2 msgs    9%  keys: svc-B|/items      |
|  P2: [||||||            ] 6 msgs   26%  keys: svc-A|/inv        |
|  P3: [|                 ] 1 msg     4%  keys: svc-C|/users      |
|  P4: [||||              ] 4 msgs   17%  keys: svc-A|/iac/inv    |
|  P5: [||                ] 2 msgs    9%  keys: svc-B|/orders     |
|                                                                  |
|  Distribution score: 72% even (ideal: 100%)                      |
|  Hot partition warning: P0 has 35% of messages                   |
+------------------------------------------------------------------+
```

**Interactions**:
- **Strategy tabs**: Switching strategy clears the partition view and replays the same messages with different assignment logic. Visual comparison shows how distribution changes.
- **Send Message**: The message animates from the input box, through the hash calculation area (numbers appear step by step), into the target partition bar. The bar grows.
- **Send 20 Random**: Quickly generates 20 messages with realistic keys and sends them one by one (fast animation, 100ms each). Shows distribution building up in real time.
- **Hash Calculation panel**: For key-based, shows the actual murmur2 computation step by step. For round-robin, shows a rotating pointer. For sticky, shows "current batch target: P2 (switching at batch boundary)".
- **Partition count selector**: Changing partition count re-hashes all existing messages and animates them to their new positions. Shows how adding partitions changes key-to-partition mapping.

**Lines of React code**: ~400
**Subcomponents**: 6 (StrategyTabs, MessageInputForm, HashCalculationPanel, PartitionBar, DistributionStats, HotPartitionWarning)
**Time**: 2 days

---

### 2.5 DELIVERY SEMANTICS

**Current content** (lines 431-521): Three delivery guarantees (at-most-once, at-least-once, exactly-once) with ASCII step sequences and code.

**Interactive Version: `<DeliverySemanticsSimulator />`**

```
+------------------------------------------------------------------+
|  DELIVERY SEMANTICS COMPARISON                                   |
|                                                                  |
|  +---- At-Most-Once ----+  +--- At-Least-Once ---+  +- Exactly-Once -+
|  |                      |  |                     |  |                |
|  |  1. Fetch [5,6,7]    |  |  1. Fetch [5,6,7]  |  |  1. Fetch      |
|  |  2. COMMIT offset 8  |  |  2. Process msg 5  |  |  2. Begin TX   |
|  |  3. Process msg 5    |  |  3. Process msg 6  |  |  3. Process    |
|  |  4. Process msg 6    |  |  4. Process msg 7  |  |  4. Commit TX  |
|  |  5. CRASH!           |  |  5. CRASH!         |  |  5. CRASH!     |
|  |  6. Restart at 8     |  |  6. Restart at 5   |  |  6. Rollback   |
|  |  7. Msgs 6,7 LOST    |  |  7. Msgs 5,6 DUPE  |  |  7. Replay     |
|  |                      |  |                     |  |                |
|  +----------------------+  +---------------------+  +----------------+
|                                                                  |
|  CLICK WHERE THE FAILURE OCCURS:                                 |
|  [Step 1] [Step 2] [Step 3] [Step 4] [Step 5] [Step 6] [Step 7] |
|                                                                  |
|  Result: 2 messages LOST (6, 7 never processed)                  |
+------------------------------------------------------------------+
```

**Three-Column Animated Comparison**:

Each column runs the same scenario simultaneously but with different commit strategies. The user watches all three in parallel.

**Interactive Failure Injection**:
- User clicks on any step in any column to inject a crash at that point
- The column animates: red explosion icon at crash point, then restart animation
- The outcome panel shows what happened:
  - "At-Most-Once: Messages 6 and 7 are LOST -- offset was already committed past them"
  - "At-Least-Once: Messages 5 and 6 are REPROCESSED -- offset was not yet committed"
  - "Exactly-Once: Transaction rolled back. All 3 messages reprocessed exactly once."
- A "messages processed" counter shows the final tally: some show < 3, some show > 3, EOS shows exactly 3

**Configuration Panel Below**:
```
Producer Settings:
  acks:  [0] [1] [all*]
  enable.idempotence: [on*] [off]
  retries: [0] [3] [MAX*]

Consumer Settings:
  enable.auto.commit: [true] [false*]
  isolation.level: [read_uncommitted] [read_committed*]
```

Changing these settings updates the animations and outcome descriptions.

**Lines of React code**: ~500
**Subcomponents**: 7 (SemanticColumn, StepBlock, CrashInjector, OutcomePanel, MessageCounter, ConfigPanel, ComparisonHeader)
**Time**: 2.5 days

---

### 2.6 CONSUMER GROUPS AND REBALANCING

**Current content** (lines 525-621): Consumer group basics, rebalancing sequence diagram, rebalance protocols (eager vs cooperative), triggers table.

**Interactive Version: `<ConsumerGroupSimulator />`**

```
+------------------------------------------------------------------+
|  CONSUMER GROUP SIMULATOR                                        |
|                                                                  |
|  Topic: audit-logs    Partitions: [6]                            |
|  Assignment Strategy: [Range] [RoundRobin*] [Sticky] [Cooperative]|
+------------------------------------------------------------------+
|                                                                  |
|  PARTITIONS:                                                     |
|  +----+ +----+ +----+ +----+ +----+ +----+                      |
|  | P0 | | P1 | | P2 | | P3 | | P4 | | P5 |                     |
|  +--+-+ +--+-+ +--+-+ +--+-+ +--+-+ +--+-+                      |
|     |      |      |      |      |      |                         |
|     v      v      |      |      v      v                         |
|  +----------+     |      |   +----------+                        |
|  |Consumer 1|     v      v   |Consumer 3|                        |
|  | P0, P1   |  +----------+  | P4, P5   |                        |
|  +----------+  |Consumer 2|  +----------+                        |
|                | P2, P3   |                                      |
|                +----------+                                      |
|                                                                  |
|  [+ Add Consumer]  [Drag consumers below to remove]              |
|                                                                  |
|  +----------------------------------------------------------+   |
|  |  REMOVE ZONE (drag consumer here to simulate failure)     |   |
|  +----------------------------------------------------------+   |
|                                                                  |
+------------------------------------------------------------------+
|  REBALANCE LOG:                                                  |
|  14:23:01 - Consumer 3 joined group. Triggering rebalance.       |
|  14:23:01 - Revoking: C1=[P0,P1,P2], C2=[P3,P4,P5]             |
|  14:23:02 - Running RoundRobin assignor...                       |
|  14:23:02 - Assigned: C1=[P0,P3], C2=[P1,P4], C3=[P2,P5]       |
|  14:23:03 - All consumers resumed processing.                    |
|  REBALANCE TIME: 1.8 seconds (eager) | 0.3 seconds (cooperative) |
+------------------------------------------------------------------+
```

**Interactions**:

- **Drag a consumer to the "Remove Zone"**:
  1. Consumer box turns red and fades out (simulates crash)
  2. Its assigned partitions flash red, then animate: lines disconnect, partitions float up
  3. Rebalance animation plays:
     - **Eager**: ALL partition lines disconnect, brief "stop-the-world" flash, then all reassign
     - **Cooperative**: ONLY the dead consumer's partitions disconnect and reassign; other consumers keep working (subtle green pulse to show continued processing)
  4. Rebalance log scrolls with real-time entries

- **"+ Add Consumer" button**:
  1. New consumer box slides in
  2. Rebalance triggers (based on selected strategy)
  3. Partitions redistribute with animation

- **Assignment Strategy tabs**:
  - Switching strategies triggers a re-assignment animation
  - Each strategy produces different partition-to-consumer mappings
  - Visual comparison label shows: "Range: uneven (C1 gets 2, C3 gets 2, C2 gets 2)" vs "RoundRobin: even distribution"
  - Side panel explains the algorithm for the selected strategy

- **Add consumers beyond partition count**: If 7th consumer added with 6 partitions, it appears grayed out with label "IDLE -- no partitions assigned" and a tooltip explaining why.

- **Partition count slider**: Adjusting from 6 to 12 partitions triggers reassignment animation.

**Lines of React code**: ~550
**Subcomponents**: 8 (PartitionChip, ConsumerBox, AssignmentLine, RemoveZone, StrategySelector, RebalanceLog, IdleConsumerIndicator, RebalanceTimer)
**Time**: 3 days

---

### 2.7 KAFKA CONNECT AND KAFKA STREAMS

**Current content** (lines 625-708): Kafka Connect concepts, Walmart's GCS sink YAML config, Kafka Streams code example, comparison table.

**Interactive Version: `<KafkaConnectPipeline />`**

```
+------------------------------------------------------------------+
|  KAFKA CONNECT PIPELINE BUILDER                                  |
|                                                                  |
|  SOURCE CONNECTORS:        KAFKA           SINK CONNECTORS:      |
|  +----------------+       TOPICS         +------------------+    |
|  | [Debezium CDC] |---->+--------+------>| [GCS Connector]  |    |
|  | [File Source]   |     | audit- |       | [ES Connector]   |    |
|  | [REST Source]   |     | logs   |       | [JDBC Connector]  |    |
|  +----------------+     +--------+       +------------------+    |
|                                                                  |
|  SMT Pipeline (click to configure):                              |
|  [Source] -> [InsertTimestamp] -> [FilterUS] -> [Route] -> [Sink] |
|                                                                  |
+------------------------------------------------------------------+
```

This is primarily an illustrative diagram with hover tooltips rather than deep simulation. Each connector box has a hover card explaining what it does. The SMT pipeline is clickable, expanding each transform to show its configuration YAML.

**Lines of React code**: ~200
**Subcomponents**: 4 (ConnectorCard, SMTStage, PipelineArrow, ConfigPanel)
**Time**: 1 day

---

### 2.8 MULTI-REGION ACTIVE/ACTIVE

**Current content** (lines 712-817): Walmart's dual-region architecture diagram, dual KafkaTemplate code, design decisions table.

**Interactive Version: `<MultiRegionSimulator />`**

```
+------------------------------------------------------------------+
|  MULTI-REGION KAFKA: ACTIVE/ACTIVE                               |
|                                                                  |
|  SCENARIO: [Normal Operation*] [EUS2 Kafka Down] [SCUS Region Down] [Both Down]
|                                                                  |
|  +--- Azure Front Door (Global LB) ---+                         |
|  |  Traffic: [======EUS2======|==SCUS==]  50/50 split            |
|  +-------------------------------------+                         |
|           |                    |                                  |
|           v                    v                                  |
|  +--- EUS2 Region ---+   +--- SCUS Region ---+                  |
|  | audit-api-logs-srv|   | audit-api-logs-srv |                  |
|  |    [4 pods]       |   |    [4 pods]        |                  |
|  |        |          |   |        |           |                  |
|  |        v          |   |        v           |                  |
|  |  Kafka Cluster    |   |  Kafka Cluster     |                  |
|  |  [12 part, RF=3]  |   |  [12 part, RF=3]   |                  |
|  |        |          |   |        |            |                  |
|  |  GCS Sink Connect |   |  GCS Sink Connect  |                  |
|  +--------+----------+   +--------+-----------+                  |
|           |                        |                              |
|           +--- GCS Buckets --------+                              |
|           |  [US] [CA] [MX]       |                              |
|           +--- BigQuery -----------+                              |
+------------------------------------------------------------------+
```

**Scenario Selector**:
- **Normal Operation**: Both regions green. Message flow arrows animated from both regions to GCS.
- **EUS2 Kafka Down**: EUS2 Kafka cluster turns red. Message animation shows: "Primary publish fails -> CompletableFuture exceptionally -> failover to SCUS Kafka." The srv in EUS2 draws a cross-region dashed line to SCUS Kafka.
- **SCUS Region Down**: Entire SCUS region box turns red. Azure Front Door animation shows traffic re-routing 100% to EUS2. EUS2 scales up (pod count animates from 4 to 8).
- **Both Down**: Everything red. Error state displayed. "DR recovery: manual intervention required."

**Click the "audit-api-logs-srv" pod**: Expands to show the CompletableFuture code with primary/secondary template highlighting.

**Click the GCS Sink**: Expands to show the 3-connector SMT routing diagram with US/CA/MX filter logic.

**Lines of React code**: ~400
**Subcomponents**: 6 (RegionBox, KafkaClusterMini, GCSBucketRow, TrafficBar, FailoverArrow, ScenarioSelector)
**Time**: 2 days

---

### 2.9 COMPLETABLEFUTURE FAILOVER PATTERN

**Current content** (lines 820-869): Java code for CompletableFuture-based failover + comparison table.

**Interactive Version: `<CompletableFutureVisualizer />`**

```
+------------------------------------------------------------------+
|  COMPLETABLEFUTURE FAILOVER PATTERN                              |
+------------------------------------------------------------------+
|                                                                  |
|  EXECUTION FLOW (click nodes to inject failures):                |
|                                                                  |
|  supplyAsync(primarySend)                                        |
|       |                                                          |
|       +-- [SUCCESS] --> return SendResult    (happy path)        |
|       |                                                          |
|       +-- [TIMEOUT 5s] -+                                       |
|       |                  |                                       |
|       +-- [EXCEPTION] --+--> exceptionallyCompose()              |
|                               |                                  |
|                          supplyAsync(secondarySend)              |
|                               |                                  |
|                               +-- [SUCCESS] --> return result    |
|                               |                                  |
|                               +-- [EXCEPTION] --> TOTAL FAILURE  |
|                                    metrics++                     |
|                                                                  |
|  [Simulate: Primary Success]                                     |
|  [Simulate: Primary Timeout -> Secondary Success]                |
|  [Simulate: Primary Exception -> Secondary Success]              |
|  [Simulate: Both Fail]                                           |
+------------------------------------------------------------------+
|  TIMELINE:                                                       |
|  0ms         1000ms       5000ms       6000ms       10000ms      |
|  |-- primary send --|                                            |
|                         |--- secondary send --|                  |
|  Result: SUCCESS via secondary (total: 6,243ms)                  |
+------------------------------------------------------------------+
```

The flow chart is an animated state machine. Clicking "Simulate: Primary Timeout" animates:
1. "supplyAsync" node lights up amber
2. Timer counts up to 5000ms
3. "TIMEOUT" path flashes red
4. "exceptionallyCompose" node lights up
5. Secondary "supplyAsync" lights up
6. Timer counts to ~1000ms
7. "SUCCESS" path lights green
8. Timeline bar below shows the full duration

**Lines of React code**: ~300
**Subcomponents**: 5 (FlowNode, FlowEdge, SimulationButton, TimelineBar, MetricsCounter)
**Time**: 1.5 days

---

### 2.10 PERFORMANCE TUNING TABLES

**Current content** (lines 872-941): Three large config tables (Producer, Consumer, Broker) with Parameter, Default, Recommended, Explanation columns.

**Interactive Version: `<TuningDashboard />`**

```
+------------------------------------------------------------------+
|  KAFKA PERFORMANCE TUNING                                        |
|  Tab: [Producer*] [Consumer] [Broker]                            |
+------------------------------------------------------------------+
|                                                                  |
|  VISUAL TUNING (drag sliders):                                   |
|                                                                  |
|  batch.size:     [16KB ----[====]--------- 256KB]  = 64KB       |
|  linger.ms:      [0 ----[====]-------------- 100]  = 5ms        |
|  acks:           [0]  [1]  [all*]                                |
|  compression:    [none]  [snappy]  [lz4*]  [zstd]               |
|  buffer.memory:  [16MB ----[====]--------- 128MB]  = 64MB       |
|                                                                  |
|  IMPACT PREVIEW:                                                 |
|  +------------------+------------------+------------------+      |
|  | Throughput       | Latency (p99)    | Durability       |      |
|  |  ~~~~~~~~        |  ~~~~~~~~        |  ~~~~~~~~        |      |
|  |  /        \      |    /\            |        ________  |      |
|  | /    85K   \     |   /  \ 12ms      |  HIGH /        | |      |
|  |/   msg/sec  \    |  /    \          | _____/          | |      |
|  +------------------+------------------+------------------+      |
|                                                                  |
|  PRESET PROFILES:                                                |
|  [Low Latency] [High Throughput*] [Balanced] [Walmart Audit]    |
|                                                                  |
+------------------------------------------------------------------+
|  DETAILED TABLE (hover any row for tooltip):                     |
|  Parameter           | Current | Default | Status               |
|  batch.size          | 64KB    | 16KB    | Modified (amber)     |
|  linger.ms           | 5       | 0       | Modified (amber)     |
|  acks                | all     | all     | Default (gray)       |
|  ...                 |         |         |                      |
+------------------------------------------------------------------+
```

**Interactions**:
- **Drag sliders**: Real-time update of the impact graphs (throughput, latency, durability) using approximate formulas.
- **Preset profiles**: Click "Walmart Audit" to set sliders to `batch.size=65536, linger.ms=20, acks=all, compression=snappy`. Shows actual values from the guide.
- **Hover any table row**: Rich tooltip explaining what the parameter does, trade-offs, and when to change it.
- **Click a parameter name**: Expands inline to show the relevant code snippet configuring this parameter.
- **"Walmart Audit" preset** shows a note: "These are the actual values used by the audit-api-logs-srv at Walmart processing 2M+ events/day."

**Lines of React code**: ~400
**Subcomponents**: 7 (TuningSlider, ImpactChart, PresetButton, TuningTable, ParameterTooltip, TabBar, ComparisonView)
**Time**: 2.5 days

---

### 2.11 AVRO SCHEMA + EVOLUTION

**Current content** (lines 944-1213): What is Avro, schema definition JSON, complex types, schema evolution (add/remove/rename fields).

**Interactive Version: `<SchemaEvolutionPlayground />`**

```
+------------------------------------------------------------------+
|  AVRO SCHEMA EVOLUTION PLAYGROUND                                |
+------------------------------------------------------------------+
|                                                                  |
|  SCHEMA EDITOR (editable JSON):           SCHEMA TIMELINE:       |
|  +---------------------------------+    +---+---+---+            |
|  | {                               |    |v1 |v2*|v3 |            |
|  |   "type": "record",            |    +---+---+---+            |
|  |   "name": "LogEvent",          |                              |
|  |   "fields": [                   |    DIFF VIEW:               |
|  |     { "name": "request_id",    |    + trace_id (added)       |
|  |       "type": "string" },      |    (green highlight)        |
|  |     { "name": "service_name",  |                              |
|  |       "type": "string" },      |    COMPATIBILITY:           |
|  |     { "name": "response_code", |    v2 vs v1:                |
|  |       "type": "int" },         |    [x] Backward             |
|  |     { "name": "trace_id",      |    [x] Forward              |
|  |       "type": ["null","string"]|    [x] Full                  |
|  |       "default": null }        |                              |
|  |   ]                            |                              |
|  | }                              |                              |
|  +---------------------------------+                              |
|                                                                  |
|  SIMULATION:                                                     |
|  Old Producer (v1) --> [message without trace_id]                |
|  New Consumer (v2) reads it --> trace_id = null (default)  [OK]  |
|                                                                  |
|  New Producer (v2) --> [message with trace_id = "abc"]           |
|  Old Consumer (v1) reads it --> trace_id ignored           [OK]  |
+------------------------------------------------------------------+
```

**Interactions**:
- **Schema Timeline**: Click v1, v2, v3 to load different schema versions in the editor. Diff view highlights additions (green), removals (red), modifications (amber).
- **Edit the schema directly**: Users can add/remove fields. On each keystroke, the compatibility checker runs and updates the Backward/Forward/Full checkboxes in real time.
- **"Add a required field without default" test**: If user adds `{"name": "region", "type": "string"}` (no default), the Backward checkbox turns red with X and a warning: "INCOMPATIBLE: Old data missing 'region' cannot be read by new consumer (no default to fill in)."
- **Simulation panel**: Shows what happens when old producer writes to new consumer and vice versa. Animated message flow with the field values visible.
- **"Try Evolution" presets**: Buttons for common operations: [Add Optional Field] [Add Required Field] [Remove Field] [Rename Field] [Change Type]. Each sets up the scenario and shows the result.

**Lines of React code**: ~500
**Subcomponents**: 8 (SchemaEditor, SchemaDiffView, CompatibilityChecker, VersionTimeline, EvolutionSimulation, FieldHighlighter, PresetButtons, WarningBanner)
**Time**: 3 days

---

### 2.12 SCHEMA REGISTRY ARCHITECTURE

**Current content** (lines 1420-1498): Schema Registry architecture diagram, how schemas are stored, leader election.

**Interactive Version: `<SchemaRegistryExplorer />`**

```
+------------------------------------------------------------------+
|  SCHEMA REGISTRY: REQUEST FLOW                                   |
|                                                                  |
|  Scenario: [Producer Registers Schema*] [Consumer Looks Up]      |
|            [Compatibility Check]  [Cache Hit]                    |
+------------------------------------------------------------------+
|                                                                  |
|  [Producer]                    [Schema Registry]   [_schemas]    |
|     |                            (Leader)           (Kafka)      |
|     |-- POST /subjects/         |                      |        |
|     |   topic-value/versions -->|                      |        |
|     |                           |-- Check compat ----->|        |
|     |                           |<-- OK ---------------|        |
|     |                           |-- Store schema ----->|        |
|     |<-- { "id": 42 } ---------|                      |        |
|     |                                                            |
|  [Local Cache: fingerprint -> 42]                                |
|                                                                  |
|  Next time: Cache hit, no Schema Registry call.                  |
+------------------------------------------------------------------+
```

Similar to the Message Flow walkthrough (2.2) -- step-by-step animated sequence with scenario selector.

**Lines of React code**: ~300
**Subcomponents**: Reuses StepByStepAnimation, SequenceColumn, AnimatedArrow from 2.2
**Time**: 1 day (reuses components)

---

### 2.13 COMPATIBILITY MATRIX

**Current content** (lines 1264-1276): Compatibility matrix table (Change type vs BACKWARD/FORWARD/FULL/NONE).

**Interactive Version: `<CompatibilityMatrix />`**

```
+------------------------------------------------------------------+
|  SCHEMA COMPATIBILITY MATRIX                                     |
+------------------------------------------------------------------+
|                                                                  |
|  Hover a cell to see explanation. Click to see example.          |
|                                                                  |
|                  | BACKWARD | FORWARD | FULL | NONE |           |
|  ----------------+----------+---------+------+------+           |
|  Add w/ default  |   YES    |   YES   | YES  | YES  |           |
|  Add w/o default |   NO     |   YES   |  NO  | YES  |           |
|  Remove w/ def   |   YES    |   YES   | YES  | YES  |           |
|  Remove w/o def  |   YES    |    NO   |  NO  | YES  |           |
|  Rename (alias)  |   YES    |   YES   | YES  | YES  |           |
|  Change type     |  MAYBE   |  MAYBE  |MAYBE | YES  |           |
|  Add enum symbol |   NO     |   YES   |  NO  | YES  |           |
|  Remove enum sym |   YES    |    NO   |  NO  | YES  |           |
|                                                                  |
|  "YES" cells: green background                                   |
|  "NO" cells:  red background                                     |
|  "MAYBE" cells: amber background                                 |
|                                                                  |
+------------------------------------------------------------------+
|  HOVER DETAIL (appears on hover over any cell):                  |
|  "Add field WITHOUT default + BACKWARD mode = INCOMPATIBLE       |
|   because: New consumer (reader) expects the field, but old      |
|   data (writer) does not have it and there is no default."       |
+------------------------------------------------------------------+
```

**Interactions**:
- **Hover a cell**: Tooltip explains WHY this combination is compatible or not, with a concrete example.
- **Click a cell**: Expands a panel below showing a before/after schema and what happens when producer/consumer versions mismatch.
- **Hover a column header (e.g., "BACKWARD")**: Highlights the entire column with a blue glow. Shows definition: "New reader can read old data."
- **Hover a row header**: Highlights the entire row.
- **"Decision Wizard" button**: Opens a dialog: "What change do you want to make?" -> "What compatibility mode is configured?" -> "Result: SAFE / UNSAFE with explanation."

**Lines of React code**: ~250
**Subcomponents**: 5 (MatrixCell, ColumnHighlighter, RowHighlighter, DetailTooltip, DecisionWizard)
**Time**: 1.5 days

---

### 2.14 AVRO SERIALIZATION FLOW (Wire Format)

**Current content** (lines 1502-1533): The 5-byte wire format diagram showing magic byte + schema ID + Avro bytes.

**Interactive Version: `<WireFormatInspector />`**

```
+------------------------------------------------------------------+
|  AVRO WIRE FORMAT INSPECTOR                                      |
+------------------------------------------------------------------+
|                                                                  |
|  RAW BYTES (hover each section):                                 |
|                                                                  |
|  +------+------+------+------+------+------+------+------+---   |
|  | 0x00 | 0x00 | 0x00 | 0x00 | 0x2A | 0x0C | 0x1A | 0x63 |.. |
|  +------+------+------+------+------+------+------+------+---   |
|    ^                              ^     ^                        |
|    |                              |     |                        |
|    Magic Byte                Schema ID   Avro Binary Data        |
|    (always 0x00)             (42)        (LogEvent fields)       |
|                                                                  |
|  DECODED VIEW:                                                   |
|  +---------------------------+                                   |
|  | Magic Byte: 0x00          |  "Identifies Schema Registry     |
|  | Schema ID:  42            |   wire format"                    |
|  | Avro Data:                |                                   |
|  |   request_id: "abc-123"   |  Schema ID 42 = LogEvent v2      |
|  |   service_name: "cp-nrti" |  Registered 2025-03-15            |
|  |   response_code: 200      |                                   |
|  |   trace_id: "trace-xyz"   |                                   |
|  +---------------------------+                                   |
|                                                                  |
|  [Encode a message]  [Decode raw bytes]                          |
+------------------------------------------------------------------+
```

**Interactions**:
- **Hover over each byte**: Highlights the byte and shows its role in a tooltip.
- **Hover over "Magic Byte" section (byte 0)**: "Always 0x00. Signals that this message uses Confluent Schema Registry wire format."
- **Hover over "Schema ID" section (bytes 1-4)**: "4-byte big-endian integer. Value: 42. Used to look up the schema from Schema Registry."
- **Hover over "Avro Data" section (bytes 5+)**: "Binary Avro-encoded LogEvent. No field names included -- consumer uses the schema to decode."
- **"Encode a message" button**: Opens a form where user enters field values; the inspector generates the byte representation in real time.
- **"Decode raw bytes" button**: Opens a hex input field; user pastes bytes and sees decoded fields.

**Lines of React code**: ~250
**Subcomponents**: 5 (ByteCell, SectionHighlight, DecodedView, EncodeForm, DecodeForm)
**Time**: 1.5 days

---

### 2.15 CODE EXAMPLES (Producer / Consumer)

**Current content** (lines 1300-1416): Full Java producer and consumer code with KafkaAvroSerializer, Avro record building, consumer loop with manual commit.

**Interactive Version: `<AnnotatedCodeBlock />`**

```
+------------------------------------------------------------------+
|  JAVA PRODUCER EXAMPLE                     [Java*] [Python] [Go] |
+------------------------------------------------------------------+
|                                                                  |
|  1  | Properties props = new Properties();                       |
|  2  | props.put("bootstrap.servers",    <-- Tooltip: Kafka      |
|     |     "kafka-eus2.walmart.com:9093");   cluster address      |
|  3  | props.put("key.serializer",                                |
|     |     "...StringSerializer");                                |
|  4  | props.put("value.serializer",       <-- Tooltip: Uses     |
|     |     "...KafkaAvroSerializer");          Avro, not JSON     |
|  5  | props.put("schema.registry.url",                           |
|     |     "https://schema-registry...");                          |
|  ...                                                             |
|  16 | LogEvent event = LogEvent.newBuilder()                      |
|  17 |     .setRequestId(UUID.randomUUID()   <-- Tooltip: Used   |
|     |         .toString())                      for dedup in BQ  |
|  ...                                                             |
|                                                                  |
|  MODE: [Read] [Step-Through*] [Run In Head]                     |
|                                                                  |
|  STEP-THROUGH STATE:                                             |
|  Line 17: event.requestId = "f47ac10b-58cc-..."                  |
|  Line 18: event.serviceName = "cp-nrti-apis"                     |
|  Variables: { event: LogEvent(requestId="f47ac...", ...) }       |
+------------------------------------------------------------------+
```

**Three Modes**:

1. **Read Mode**: Standard syntax-highlighted code (what we have now). Each line has a faint "?" icon on the right margin. Clicking "?" shows a tooltip explaining that line.

2. **Step-Through Mode**:
   - Current line is highlighted with amber background
   - "Next" / "Previous" buttons advance the highlight
   - Below the code, a "State Panel" shows all variables and their current values
   - When stepping through the builder pattern, each `.set*()` call updates the variable display
   - When reaching `producer.send()`, a mini animation shows the message being serialized (Avro bytes appear) and sent

3. **Run In Head Mode**:
   - All lines are grayed out
   - Lines are revealed one at a time with a 2-second delay
   - After each line is revealed, a quiz prompt appears: "What will the value of `event.serviceName` be?" with an input field
   - User answers, and correctness is shown before proceeding

**Language Toggle** (Java / Python / Go):
- Same logic implemented in three languages
- Each language tab has equivalent code with the same line-by-line annotations
- Python version uses `confluent-kafka-python`, Go uses `confluent-kafka-go`

**Click any class/type name** (e.g., `KafkaTemplate`, `KafkaAvroSerializer`, `ProducerRecord`):
- Hover card appears with:
  ```
  KafkaAvroSerializer
  Package: io.confluent.kafka.serializers
  Purpose: Serializes Java objects to Avro binary format
  with Schema Registry integration.
  Key behavior: Registers schema on first use, caches
  schema ID, prepends 5-byte header.
  ```

**Lines of React code**: ~500
**Subcomponents**: 8 (AnnotatedLine, LineTooltip, StepThroughController, StatePanel, QuizPrompt, LanguageToggle, ClassHoverCard, RunInHeadTimer)
**Time**: 3 days

---

### 2.16 INTERVIEW Q&A (25 Questions)

**Current content** (lines 1797-2219): 25 interview questions with detailed paragraph answers.

**Interactive Version: `<InterviewFlashcardDeck />`**

```
+------------------------------------------------------------------+
|  INTERVIEW PREP: KAFKA, AVRO & SCHEMA REGISTRY                  |
|                                                                  |
|  Mode: [Flashcards*] [Practice] [Quiz] [Review All]             |
|  Progress: 8/25 mastered  [=======           ] 32%              |
|  Session: 12 cards reviewed, 4 correct streak                    |
+------------------------------------------------------------------+
|                                                                  |
|  +------------------------------------------------------+       |
|  |                                                      |       |
|  |    Q3: What is a consumer group and why?             |       |
|  |                                                      |       |
|  |    Difficulty: [Easy] [Medium*] [Hard]               |       |
|  |    Topic: Consumer Groups | Rebalancing              |       |
|  |                                                      |       |
|  |              [ FLIP CARD ]                           |       |
|  |                                                      |       |
|  +------------------------------------------------------+       |
|                                                                  |
|  [< Previous]    Card 3 of 25    [Next >]   [Shuffle]           |
|                                                                  |
+------------------------------------------------------------------+
```

**After flipping** (card rotates with 3D CSS transform):

```
+------------------------------------------------------------------+
|  +------------------------------------------------------+       |
|  |  A3: Consumer Groups                                 |       |
|  |                                                      |       |
|  |  A consumer group is a set of consumer instances     |       |
|  |  that cooperate to consume a topic. Kafka assigns    |       |
|  |  each partition to exactly one consumer within the   |       |
|  |  group.                                              |       |
|  |                                                      |       |
|  |  KEY POINTS:                                         |       |
|  |  * Parallel processing (12 partitions = 12 max)     |       |
|  |  * Fault tolerance (dead consumer -> rebalance)     |       |
|  |  * Pub-sub via multiple groups                       |       |
|  |  * CooperativeStickyAssignor for incremental rebal  |       |
|  |                                                      |       |
|  |  WALMART CONNECTION:                                 |       |
|  |  Two groups: GCS sink + monitoring dashboard         |       |
|  |                                                      |       |
|  +------------------------------------------------------+       |
|                                                                  |
|  How well did you know this?                                     |
|  [Again (1min)] [Hard (10min)] [Good (1day)] [Easy (4days)]     |
+------------------------------------------------------------------+
```

**Modes**:

1. **Flashcards**: Card flip with self-rating (spaced repetition scheduling). Cards due for review are prioritized. Rating determines next review interval:
   - Again: 1 minute
   - Hard: 10 minutes
   - Good: 1 day
   - Easy: 4 days
   - Intervals multiply by 2.5x on consecutive "Good" ratings (Anki-like algorithm)

2. **Practice Mode**:
   - Timer starts (2 minutes per question)
   - Question displayed, no peeking at answer
   - User types their answer in a text area
   - After submitting, both their answer and the model answer are shown side by side
   - AI-assisted comparison (future feature -- for now, manual self-assessment)

3. **Quiz Mode**:
   - Random selection of 10 questions
   - Multiple choice: each question shows 4 answer excerpts (1 correct, 3 from other questions)
   - Scored at the end with percentage and breakdown

4. **Review All**: Traditional scrollable Q&A list with expand/collapse. Search/filter by topic tag.

**Data Model** (stored in localStorage):

```javascript
{
  "flashcards": {
    "kafka-q1": {
      "interval": 86400000,     // 1 day in ms
      "easeFactor": 2.5,
      "dueDate": "2026-02-09T14:00:00Z",
      "repetitions": 3,
      "lastRating": "good",
      "tags": ["architecture", "fundamentals"]
    },
    // ...25 cards
  },
  "stats": {
    "totalReviews": 47,
    "masteredCount": 8,
    "averageEase": 2.3,
    "streakDays": 3
  }
}
```

**Tag System**: Each question is tagged (e.g., "architecture", "consumer-groups", "avro", "schema-registry", "walmart-experience", "exactly-once", "multi-region"). Users can filter by tag to focus study.

**Lines of React code**: ~600
**Subcomponents**: 10 (FlashcardDeck, CardFront, CardBack, FlipAnimation, RatingButtons, ProgressBar, PracticeMode, QuizMode, ReviewList, SpacedRepetitionEngine)
**Time**: 3 days

---

### 2.17 WALMART SYSTEM OVERVIEW (Three-Tier Architecture)

**Current content** (lines 2222-2345): The three-tier architecture diagram, technologies table, metrics table, timeline.

**Interactive Version: `<WalmartArchitectureExplorer />`**

```
+------------------------------------------------------------------+
|  WALMART AUDIT LOGGING: ARCHITECTURE DEEP DIVE                   |
+------------------------------------------------------------------+
|  View: [Architecture Diagram*] [Data Flow Animation] [Timeline]  |
+------------------------------------------------------------------+
|                                                                  |
|  TIER 1: Common Library JAR                                      |
|  +-------------------------+                                     |
|  | cp-nrti-apis            |  Click any                         |
|  | inventory-status-srv    |  service to                        |
|  | (6+ consumer services)  |  see details                       |
|  |   |                     |                                     |
|  |   v                     |                                     |
|  | LoggingFilter           |  <-- click: shows @Order,           |
|  | (@Order LOWEST_PREC)    |      ContentCachingWrapper          |
|  |   |                     |                                     |
|  |   v                     |                                     |
|  | AuditLogService         |  <-- click: shows @Async config,    |
|  | (@Async, 6-10 threads)  |      thread pool details            |
|  +-------------------------+                                     |
|           |                                                       |
|           | HTTP POST (async)                                     |
|           v                                                       |
|  TIER 2: Kafka Publisher                                         |
|  +-------------------------+                                     |
|  | AuditLoggingController  |                                     |
|  | POST /v1/logRequest     |                                     |
|  |   |                     |                                     |
|  |   v                     |                                     |
|  | KafkaProducerService    |                                     |
|  | - Avro Serialization    |                                     |
|  | - Dual-Region Publish   |                                     |
|  +-------------------------+                                     |
|           |                                                       |
|           v                                                       |
|  TIER 3: Kafka Connect GCS Sink                                  |
|  +--------+--------+--------+                                    |
|  | US SMT | CA SMT | MX SMT |                                    |
|  +--------+--------+--------+                                    |
|           |                                                       |
|           v                                                       |
|  [GCS US] [GCS CA] [GCS MX] --> [BigQuery]                       |
|                                                                  |
+------------------------------------------------------------------+
|  METRICS COMPARISON (Before/After):                              |
|  Events/day:    500K --> 2M+     [animated counter]              |
|  Integration:   2-3 weeks --> 1 day                              |
|  Retention:     30 days --> 7 years                              |
|  DR Recovery:   Hours --> 15 minutes                             |
+------------------------------------------------------------------+
```

**Data Flow Animation view**: A single message (audit event) animates through all three tiers:
1. HTTP request arrives at cp-nrti-apis
2. LoggingFilter intercepts response
3. AuditLogService creates async task
4. HTTP POST to audit-api-logs-srv
5. Avro serialization + Kafka publish
6. GCS Sink reads from Kafka
7. SMT filter routes to US bucket
8. Parquet file written to GCS
9. BigQuery external table query

Each step pauses for 2 seconds with explanation text.

**Timeline view**: Interactive timeline showing key PRs and milestones. Click a milestone to see the PR details and what changed.

**Lines of React code**: ~450
**Subcomponents**: 7 (TierBox, ServiceNode, DataFlowAnimation, MetricsCounter, TimelineView, PRDetail, ComparisonTable)
**Time**: 2.5 days

---

### 2.18 COMPARISON TABLES

**Current content**: 6 comparison tables across the guide:
1. Kafka vs RabbitMQ vs SQS (line 58)
2. Kafka Connect vs Kafka Streams (line 700)
3. Avro vs JSON vs Protobuf vs Thrift (line 962)
4. Compatibility modes (line 1543)
5. Rebalance triggers (line 604)
6. Try-catch vs CompletableFuture (line 862)

**Interactive Version: `<ComparisonTable />` (reusable component)**

```
+------------------------------------------------------------------+
|  AVRO vs JSON vs PROTOBUF vs THRIFT                              |
|  Sort by: [Feature*] [Best for Kafka] [Best for Size]            |
+------------------------------------------------------------------+
|                                                                  |
|  FEATURE        | AVRO    | JSON      | PROTOBUF  | THRIFT     |
|  ===============+=========+===========+===========+============  |
|  Format         | Binary  | Text      | Binary    | Binary     |
|  Schema         | .avsc   | optional  | .proto    | .thrift    |
|  Msg Size       | ~30%    | 100%      | ~25%      | ~25%       |
|  Evolution      | [=====] | [       ] | [====]    | [====]     |
|  Kafka Support  | [=====] | [====]    | [===]     | [=]        |
|  Human Readable | [     ] | [=====]   | [     ]   | [     ]    |
|  Speed          | [====]  | [==]      | [=====]   | [====]     |
|                                                                  |
|  Hover: Column "AVRO" highlighted with blue glow                 |
|                                                                  |
+------------------------------------------------------------------+
|  DECISION WIZARD:                                                |
|  What matters most?                                              |
|  [x] Compact message size                                        |
|  [x] Schema evolution support                                    |
|  [ ] Human readability                                           |
|  [x] Kafka ecosystem integration                                 |
|  [ ] Fastest serialization speed                                 |
|                                                                  |
|  RECOMMENDATION: Apache Avro (score: 92%)                        |
|  "Best fit for Kafka pipelines with schema evolution needs.       |
|   Used by Walmart's audit logging pipeline."                     |
+------------------------------------------------------------------+
```

**Interactions**:
- **Hover column header**: Entire column gets blue/amber background highlight. Summary card appears: "Avro: Binary serialization with JSON schemas. Best for Kafka-centric architectures."
- **Hover any cell**: Tooltip with detailed explanation.
- **Click technology name**: Expands to a summary card with pros/cons and use cases.
- **Sortable columns**: Click column header to sort by that column (for numeric/boolean columns).
- **Decision Wizard**: Checkbox-based questionnaire that scores each technology and recommends the best fit. Weights are pre-configured based on common engineering priorities.
- **Visual bars**: Rating columns show filled bars instead of plain text (e.g., [=====] is a green bar, [==] is a short amber bar).

**Lines of React code**: ~350
**Subcomponents**: 7 (SortableHeader, HoverHighlightColumn, RatingBar, TechSummaryCard, DecisionWizard, CheckboxQuestion, RecommendationCard)
**Time**: 2 days

---

## 3. REUSABLE COMPONENT LIBRARY

These components would be shared across all 15 study guides:

### Core Interactive Components

| Component | Description | Used By Sections | Props |
|-----------|-------------|-----------------|-------|
| `<StepByStepAnimation>` | Animated sequence diagram with step navigation | 2.2, 2.12, 2.17 | `steps[], autoPlay, speed, configToggles[]` |
| `<InteractiveDiagram>` | SVG diagram with hover, click, and animation | 2.1, 2.8, 2.17 | `nodes[], edges[], onNodeClick, onNodeHover, animateMessages` |
| `<AnnotatedCodeBlock>` | Code with line tooltips, step-through, language toggle | 2.15 (all code in every guide) | `code, language, annotations[], languages[], mode` |
| `<ComparisonTable>` | Sortable, hoverable table with decision wizard | 2.18 (all tables in every guide) | `columns[], rows[], sortable, wizard, ratingBars` |
| `<FlashcardDeck>` | Flashcard system with spaced repetition | 2.16 (all Q&A in every guide) | `cards[], onRate, enableTimer, enableQuiz` |
| `<CommitLogVisualizer>` | Offset-based log with drag pointers | 2.3a (Kafka-specific but pattern reusable) | `offsets[], consumers[], pointers[]` |
| `<SimulatorShell>` | Container with controls, info panel, log | 2.1, 2.5, 2.6, 2.8 | `children, controls[], infoPanel, logEntries[]` |
| `<WireFormatInspector>` | Byte-level hex viewer with section highlighting | 2.14 (also useful for gRPC guide, protobuf guide) | `bytes[], sections[], decodeView` |

### Utility Components

| Component | Description | Used Everywhere | Props |
|-----------|-------------|----------------|-------|
| `<HoverTooltip>` | Rich tooltip on hover/click | Every interactive element | `content, position, maxWidth, delay` |
| `<AnimatedArrow>` | SVG arrow with message dots | Diagrams, sequence flows | `from, to, label, animated, dotColor` |
| `<ConfigToggle>` | Toggle switch or segmented control | Tuning, settings, mode | `options[], value, onChange, label` |
| `<ProgressTracker>` | Progress bar with session stats | Flashcards, step-through | `current, total, label` |
| `<InfoPanel>` | Collapsible explanation panel | Simulators, diagrams | `title, content, collapsible` |
| `<DragDropZone>` | Drag and drop target area | Consumer group sim | `onDrop, label, highlighted` |
| `<MetricsCounter>` | Animated number counter | Walmart metrics, tuning | `value, label, prefix, suffix, animationDuration` |
| `<TabBar>` | Tab navigation for modes/views | Multiple sections | `tabs[], activeTab, onChange` |
| `<SliderControl>` | Range slider with labels | Tuning dashboard | `min, max, value, step, label, formatValue` |
| `<SearchFilter>` | Search + tag filter | Flashcards, review mode | `items[], tags[], onFilter` |

### Layout Components

| Component | Description |
|-----------|-------------|
| `<InteractiveSection>` | Wrapper that detects when an interactive version exists for a markdown section and renders the interactive component instead of (or alongside) the static markdown |
| `<ToggleView>` | "View: [Static] [Interactive]" toggle that lets users switch between the original markdown and the interactive version |
| `<SplitView>` | Side-by-side layout: static content on left, interactive demo on right |

---

## 4. EFFORT ESTIMATES

### Per-Section Breakdown

| Section | Component | Lines of Code | Custom Sub-components | Time (days) | Complexity |
|---------|-----------|--------------|----------------------|-------------|------------|
| 2.1 | KafkaClusterDiagram | ~650 | 10 | 3-4 | HIGH |
| 2.2 | MessageFlowWalkthrough | ~450 | 7 | 2-3 | MEDIUM |
| 2.3a | CommitLogExplorer | ~350 | 5 | 2 | MEDIUM |
| 2.3b | ISRSimulator | ~250 | 4 | 1.5 | LOW-MED |
| 2.4 | PartitioningVisualizer | ~400 | 6 | 2 | MEDIUM |
| 2.5 | DeliverySemanticsSimulator | ~500 | 7 | 2.5 | HIGH |
| 2.6 | ConsumerGroupSimulator | ~550 | 8 | 3 | HIGH |
| 2.7 | KafkaConnectPipeline | ~200 | 4 | 1 | LOW |
| 2.8 | MultiRegionSimulator | ~400 | 6 | 2 | MEDIUM |
| 2.9 | CompletableFutureVisualizer | ~300 | 5 | 1.5 | MEDIUM |
| 2.10 | TuningDashboard | ~400 | 7 | 2.5 | MEDIUM |
| 2.11 | SchemaEvolutionPlayground | ~500 | 8 | 3 | HIGH |
| 2.12 | SchemaRegistryExplorer | ~300 | (reuses) | 1 | LOW |
| 2.13 | CompatibilityMatrix | ~250 | 5 | 1.5 | LOW-MED |
| 2.14 | WireFormatInspector | ~250 | 5 | 1.5 | MEDIUM |
| 2.15 | AnnotatedCodeBlock | ~500 | 8 | 3 | HIGH |
| 2.16 | InterviewFlashcardDeck | ~600 | 10 | 3 | HIGH |
| 2.17 | WalmartArchitectureExplorer | ~450 | 7 | 2.5 | MEDIUM |
| 2.18 | ComparisonTable (reusable) | ~350 | 7 | 2 | MEDIUM |
| --- | **Shared component library** | ~800 | 18 | 4 | HIGH |
| **TOTAL** | | **~7,450** | **~87 unique** | **~42-47 days** | |

### Rollup for Full Guide 01

- **Total lines of React code**: ~7,450 (across all components)
- **Total unique sub-components**: ~87
- **Total development time**: ~42-47 person-days (8-9 weeks for one developer)
- **Testing and polish**: Add ~30% = ~55-60 person-days total
- **For all 15 guides** (with reusable components): ~30 days for shared library + ~25 days per guide on average (some simpler, some harder) = ~30 + (15 x 25) = ~405 person-days = ~1.5 years for 1 developer, ~5 months for a team of 3

### Recommended MVP Scope (4 weeks)

Build these first -- they provide the most impact and create the reusable foundation:

| Priority | Component | Why | Time |
|----------|-----------|-----|------|
| P0 | Shared component library (core 8 components) | Foundation for everything | 4 days |
| P0 | AnnotatedCodeBlock | Used in ALL 15 guides (every code block) | 3 days |
| P0 | ComparisonTable | Used in ALL 15 guides (every table) | 2 days |
| P0 | InterviewFlashcardDeck | Used in ALL 15 guides (every Q&A section) | 3 days |
| P1 | KafkaClusterDiagram | Flagship interactive for guide 01 | 4 days |
| P1 | MessageFlowWalkthrough | Core learning for Kafka concepts | 3 days |
| P1 | ConsumerGroupSimulator | High-value hands-on learning | 3 days |
| P2 | DeliverySemanticsSimulator | Common interview topic | 2.5 days |

**MVP Total**: ~24.5 days = 5 weeks for 1 developer

---

## 5. DEPENDENCY MAP

### NPM Dependencies (new, beyond current package.json)

| Package | Purpose | Size | Required By |
|---------|---------|------|-------------|
| `framer-motion` | Smooth animations (optional -- CSS fallback) | ~150KB | Animations in all simulators |
| `@dnd-kit/core` + `@dnd-kit/sortable` | Drag-and-drop for consumer group sim | ~45KB | ConsumerGroupSimulator |
| `zustand` | Lightweight state management for simulators | ~8KB | All simulators sharing state |
| `recharts` | Small charts for tuning dashboard | ~200KB | TuningDashboard impact graphs |
| `monaco-editor` or `@uiw/react-codemirror` | Schema editor in evolution playground | ~500KB | SchemaEvolutionPlayground |

**Total new bundle size**: ~900KB (or ~400KB without monaco-editor and with CSS animations instead of framer-motion)

**Recommended minimal additions**: `zustand` (8KB) for state management and `@dnd-kit/core` (45KB) for drag-and-drop. Everything else can be achieved with CSS animations and custom SVG.

### Dependency-Free Alternatives

| Instead Of | Use |
|-----------|-----|
| framer-motion | CSS `@keyframes`, `transition`, `transform` with React state-driven class toggling |
| recharts | Custom SVG `<rect>` bars and `<line>` elements for simple bar charts |
| monaco-editor | `<textarea>` with syntax highlighting via existing `react-syntax-highlighter` for read-only preview |

---

## 6. IMPLEMENTATION PRIORITY

### Phase 1: Foundation (Week 1-2)

Build the reusable component library:
1. `HoverTooltip` -- simplest, used everywhere
2. `ConfigToggle` / `TabBar` -- used in every simulator
3. `AnimatedArrow` -- SVG arrow with CSS animation
4. `InfoPanel` -- collapsible explanation panel
5. `ProgressTracker` -- for flashcards and step-through
6. `SimulatorShell` -- standard container for all simulators

Then build the three cross-guide components:
7. `AnnotatedCodeBlock` (replaces current CodeBlock for interactive mode)
8. `ComparisonTable` (replaces current TableWrapper for interactive mode)
9. `InterviewFlashcardDeck` (new)

### Phase 2: Kafka Guide Flagships (Week 3-4)

Build the Kafka-specific interactive components:
10. `KafkaClusterDiagram` -- the hero component
11. `MessageFlowWalkthrough` -- step-by-step sequence
12. `ConsumerGroupSimulator` -- drag-and-drop rebalancing

### Phase 3: Kafka Guide Complete (Week 5-6)

Fill in the remaining Kafka-specific components:
13. `DeliverySemanticsSimulator`
14. `CommitLogExplorer`
15. `PartitioningVisualizer`
16. `TuningDashboard`
17. `SchemaEvolutionPlayground`

### Phase 4: Polish + Remaining Sections (Week 7-8)

18. `MultiRegionSimulator`
19. `CompletableFutureVisualizer`
20. `WireFormatInspector`
21. `SchemaRegistryExplorer`
22. `CompatibilityMatrix`
23. `WalmartArchitectureExplorer`
24. `KafkaConnectPipeline`

### Integration Strategy

Each interactive component should be opt-in alongside the existing markdown rendering. The `<ToggleView>` component provides:

```jsx
<ToggleView
  staticContent={<ReactMarkdown>{markdownSection}</ReactMarkdown>}
  interactiveContent={<KafkaClusterDiagram />}
  defaultView="interactive"
/>
```

This means:
- Existing markdown rendering continues to work unchanged
- Interactive components are layered on top, section by section
- Users can toggle between static and interactive views
- If a section does not have an interactive version yet, only static is shown
- Progressive enhancement: the guide is always functional, just increasingly rich

### File Structure

```
src/
  components/
    interactive/                    # NEW directory
      shared/                       # Reusable across all guides
        HoverTooltip.jsx
        AnimatedArrow.jsx
        ConfigToggle.jsx
        TabBar.jsx
        InfoPanel.jsx
        ProgressTracker.jsx
        SimulatorShell.jsx
        SliderControl.jsx
        DragDropZone.jsx
        MetricsCounter.jsx
        SearchFilter.jsx
      cross-guide/                  # Used in every guide
        AnnotatedCodeBlock.jsx
        ComparisonTable.jsx
        InterviewFlashcardDeck.jsx
        ToggleView.jsx
        SplitView.jsx
      kafka/                        # Kafka guide specific
        KafkaClusterDiagram.jsx
        MessageFlowWalkthrough.jsx
        ConsumerGroupSimulator.jsx
        DeliverySemanticsSimulator.jsx
        CommitLogExplorer.jsx
        ISRSimulator.jsx
        PartitioningVisualizer.jsx
        TuningDashboard.jsx
        SchemaEvolutionPlayground.jsx
        SchemaRegistryExplorer.jsx
        CompatibilityMatrix.jsx
        WireFormatInspector.jsx
        CompletableFutureVisualizer.jsx
        MultiRegionSimulator.jsx
        KafkaConnectPipeline.jsx
        WalmartArchitectureExplorer.jsx
    markdown/                       # EXISTING -- unchanged
      CodeBlock.jsx
      TableWrapper.jsx
      HeadingAnchor.jsx
      MermaidDiagram.jsx
  hooks/
    useSpacedRepetition.js          # NEW -- Anki-like scheduling
    useSimulatorState.js            # NEW -- shared simulator state
    useDragDrop.js                  # NEW -- drag-and-drop abstraction
    useReadingProgress.js           # EXISTING
  data/
    kafka-flashcards.json           # NEW -- 25 Q&A cards with tags
    kafka-code-annotations.json     # NEW -- line-by-line annotations
    kafka-tuning-presets.json       # NEW -- preset configurations
```

---

## APPENDIX: KEY DESIGN PRINCIPLES

1. **Progressive Disclosure**: Never overwhelm. Start with the simple view, let users drill down by clicking/hovering.

2. **Always Escapable**: Every interactive element has a clear "close" or "reset" action. Users are never trapped in a state.

3. **Markdown Fallback**: If JavaScript fails or is disabled, the original markdown renders. Interactive components enhance, they do not replace.

4. **Consistent Interaction Patterns**: Hover = tooltip. Click = expand/drill down. Drag = reorder/move. Toggle = switch mode. Same everywhere.

5. **Performant**: No component should cause jank. Use `React.memo`, `useMemo`, `useCallback`. SVG animations use CSS transforms (GPU-accelerated), not JavaScript-driven position changes.

6. **Accessible**: All interactive elements are keyboard-navigable. Tooltips work with focus, not just hover. Screen reader labels on SVG elements. Animations respect `prefers-reduced-motion`.

7. **Dark Mode Compatible**: All colors use CSS custom properties (`var(--bg-surface)`, `var(--text-primary)`, etc.) from the existing theme system in `ThemeContext.jsx`.

8. **Mobile Responsive**: Touch interactions replace hover (long-press = hover tooltip). Diagrams scale down or become scrollable. Flashcards are swipeable.

---

*End of prototype specification. This document contains sufficient detail for a developer to implement any of the described interactive components without further design input. Total estimated content: ~7,450 lines of React code across ~87 sub-components, implementable in ~55-60 person-days.*

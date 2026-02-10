# Interview Preparation - Master Index

> **Your Interviews:** Rippling Hiring Manager, Google Googleyness
> **Last Updated:** February 2026

---

## Quick Access by Topic

| Folder | Resume Bullets | Main PR | 30-Sec Pitch |
|--------|----------------|---------|--------------|
| [01-kafka-audit-logging](./01-kafka-audit-logging/) | 1, 2, 10 | #1 (dv-api-common-libraries) | "Built audit system for 2M events/day, zero latency impact, suppliers can query their own data" |
| [02-spring-boot-3-migration](./02-spring-boot-3-migration/) | 4 | #1312 (cp-nrti-apis) | "Led migration: 158 files, javax→jakarta, WebClient, zero-downtime canary deployment" |
| [03-dsd-notification-system](./03-dsd-notification-system/) | 3 | cp-nrti-apis | "Real-time SUMO push to 1,200+ associates, 300+ stores, 35% replenishment improvement" |
| [05-openapi-dc-inventory](./05-openapi-dc-inventory/) | 5, 6, 8 | inventory-status-srv | "Design-first API, DC Inventory with factory pattern, 8,000+ lines, 30% integration reduction" |
| [06-transaction-event-history](./06-transaction-event-history/) | 7 | inventory-events-srv | "Cursor pagination, Cosmos→Postgres migration, multi-tenant" |
| [07-observability](./07-observability/) | 9 | Cross-cutting | "Dynatrace tracing, Prometheus alerts, Flagger canary, 99.9% uptime" |

---

## Numbers to Memorize (CRITICAL!)

### Kafka Audit Logging
| Metric | Value |
|--------|-------|
| Events/day | **2M+** |
| P99 latency impact | **<5ms** |
| Cost vs Splunk | **99% reduction** |
| Teams adopted library | **3+** |
| Integration time | **2 weeks → 1 day** |

### Spring Boot 3 Migration
| Metric | Value |
|--------|-------|
| Files changed | **136** |
| Lines added | **1,732** |
| javax→jakarta files | **74** |
| Production issues | **0 customer-impacting** |
| Deployment | **Canary: 10%→100%** |

---

## Top 5 Stories to Tell

### 1. The Silent Failure (Debugging)
> **Location:** [01-kafka-audit-logging/06-debugging-stories.md](./01-kafka-audit-logging/06-debugging-stories.md)
>
> "Two weeks after launch, GCS buckets stopped receiving data. No alerts. Over 5 days, I found: null pointer in SMT filter, consumer poll timeouts, KEDA autoscaling causing rebalance storms, and JVM heap exhaustion. Zero data loss - Kafka retained everything."

### 2. Library Adoption (Collaboration)
> **Location:** [01-kafka-audit-logging/02-common-library.md](./01-kafka-audit-logging/02-common-library.md)
>
> "Saw 3 teams building same functionality. Proposed shared library, understood their requirements, made differences configurable. Three teams adopted in one month. Integration time: 2 weeks → 1 day."

### 3. Critical Feedback (Growth)
> **Location:** [01-kafka-audit-logging/07-interview-qa.md](./01-kafka-audit-logging/07-interview-qa.md)
>
> "Senior engineer criticized my thread pool queue size - could cause silent data loss. He was right. Added metrics for rejected tasks, warning at 80% capacity. That monitoring has caught issues twice."

### 4. Spring Boot 3 Migration (Technical Leadership)
> **Location:** [02-spring-boot-3-migration/06-interview-qa.md](./02-spring-boot-3-migration/06-interview-qa.md)
>
> "Led migration from Spring Boot 2.7 to 3.2. Strategic decision: used .block() instead of full reactive rewrite - scope control. Zero-downtime deployment using Flagger canary."

### 5. Multi-Region Rollout (Ambiguity)
> **Location:** [01-kafka-audit-logging/01-overview.md](./01-kafka-audit-logging/01-overview.md)
>
> "Requirements were vague - 'make it resilient.' I defined RTO/RPO through stakeholder conversations, designed three options, chose Active/Active. Zero downtime during migration."

---

## Interview Types

### Hiring Manager (Rippling)
**What they evaluate:**
- Technical depth (60%)
- Team fit (25%)
- Growth potential (15%)

**Focus on:**
- Use "I" for your contributions
- Quantify impact (events/day, latency, cost)
- Show trade-off awareness
- Share mistakes and learnings

### Googleyness (Google)
**What they evaluate:**
- Doing the right thing
- Thriving in ambiguity
- Valuing feedback
- Challenging status quo
- Bringing others along

**Focus on:**
- "Tell me about a time..." answers
- STAR format with specific examples
- How you handled disagreements
- How you helped others succeed

---

## START HERE THE NIGHT BEFORE

> **[INTERVIEW-ROUND-GUIDE.md](./INTERVIEW-ROUND-GUIDE.md)** - RESEARCH-BASED guide: exact structure of each round, what's evaluated, scoring rubric, kya karna hai, kya NAHI karna hai. Based on IGotAnOffer, Google re:Work, Glassdoor, Blind.
>
> **[INTERVIEW-PLAYBOOK-HOW-TO-SPEAK.md](./INTERVIEW-PLAYBOOK-HOW-TO-SPEAK.md)** - HOW to talk, lead the conversation, handle traps, and close strong.
>
> **[SCENARIO-BASED-QUESTIONS.md](./SCENARIO-BASED-QUESTIONS.md)** - "What if Kafka fails? DB down? SUMO down? Black Friday?" - 25+ failure scenarios with DETECT→IMPACT→MITIGATE→RECOVER→PREVENT framework.

---

## Before Interview Checklist

- [ ] Read the PLAYBOOK (how to speak, not what to say)
- [ ] Practice 90-second intro OUT LOUD
- [ ] Practice debugging story OUT LOUD (under 3 min)
- [ ] Review numbers (2M events, <5ms, 158 files)
- [ ] Know the 6 STAR stories and which question triggers each
- [ ] Have 3 questions ready to ask them
- [ ] Review trade-offs for each major decision

---

## Folder Structure

```
work-ex/
├── INTERVIEW-MASTER-INDEX.md          # This file
├── 01-kafka-audit-logging/            # Bullets 1, 2, 10
│   ├── README.md
│   ├── 01-overview.md
│   ├── 02-common-library.md
│   ├── 03-kafka-publisher.md
│   ├── 04-gcs-sink.md
│   ├── 05-multi-region.md
│   ├── 06-debugging-stories.md
│   ├── 07-interview-qa.md
│   └── 08-prs-reference.md
├── 02-spring-boot-3-migration/        # Bullet 4
│   ├── README.md
│   ├── 01-overview.md
│   ├── 02-technical-challenges.md
│   ├── 03-code-changes.md
│   ├── 04-deployment-strategy.md
│   ├── 05-production-issues.md
│   └── 06-interview-qa.md
├── 03-dsd-notification-system/        # Bullet 3
│   └── README.md
├── 04-common-library-jar/             # Bullet 2 (→ redirects to 01/02)
│   └── README.md
├── QUICK-REFERENCE-CARDS.md           # One-page cheat sheet
├── BEHAVIORAL-STORIES-STAR.md         # STAR format stories
├── MOCK-INTERVIEW-QUESTIONS.md        # 24 Q&A with answers
├── TECHNICAL-CONCEPTS-REVIEW.md       # 10 deep concepts
└── GAP-ANALYSIS.md                    # Areas to strengthen
```

---

## Old Files (Still Useful)

| File | Purpose |
|------|---------|
| QUICK-REFERENCE-CARDS.md | One-page metrics and diagrams |
| BEHAVIORAL-STORIES-STAR.md | 6 STAR-format stories |
| MOCK-INTERVIEW-QUESTIONS.md | 24 questions with model answers |
| GAP-ANALYSIS.md | Areas needing more preparation |
| INTERVIEW-MASTERCLASS-KAFKA-AUDIT.md | Comprehensive Kafka guide |
| DEEP-DIVE-BULLET-4-SPRINGBOOT3-JAVA17.md | Spring Boot 3 deep dive |

---

## Power Phrases to Use

### Demonstrating Ownership
- "I designed..."
- "I made the decision to..."
- "I took responsibility for..."

### Showing Technical Depth
- "The key insight was..."
- "The trade-off we accepted was..."
- "Under the hood, what happens is..."

### Demonstrating Learning
- "What I learned from this was..."
- "If I did this again, I would..."
- "The mistake I made was..., which taught me..."

---

## Red Flags to Avoid

| Don't Do | Do Instead |
|----------|------------|
| "We did..." | "I did..." (for your work) |
| Long pauses | "Let me trace through that..." |
| Getting defensive | "That's a good point, I hadn't considered..." |
| No metrics | Always have numbers ready |
| No failures | Share the debugging story |

---

*Good luck with your interviews!*

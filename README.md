# Google L4 Interview Preparation

Personal interview preparation materials for Google L4 Software Engineer role - Googleyness & Hiring Manager rounds.

---

## Quick Start - What to Read

| Priority | File | Time | Purpose |
|----------|------|------|---------|
| 1 | [GOOGLEYNESS_ALL_QUESTIONS.md](GOOGLEYNESS_ALL_QUESTIONS.md) | 30 min | **60+ questions with STAR answers** |
| 2 | [GOOGLE_INTERVIEW_SCRIPTS.md](GOOGLE_INTERVIEW_SCRIPTS.md) | 15 min | Word-by-word scripts |
| 3 | [GOOGLE_L4_FINAL_PREP.md](GOOGLE_L4_FINAL_PREP.md) | 20 min | Complete interview guide |

---

## All Files

### Interview Preparation
| File | Description |
|------|-------------|
| [GOOGLEYNESS_ALL_QUESTIONS.md](GOOGLEYNESS_ALL_QUESTIONS.md) | 60+ behavioral questions with detailed STAR answers |
| [GOOGLE_L4_FINAL_PREP.md](GOOGLE_L4_FINAL_PREP.md) | Complete L4 interview guide - structure, tips, metrics |
| [GOOGLE_INTERVIEW_SCRIPTS.md](GOOGLE_INTERVIEW_SCRIPTS.md) | Exact word-by-word scripts for top questions |
| [GOOGLE_INTERVIEW_PREP.md](GOOGLE_INTERVIEW_PREP.md) | Original STAR stories for Googleyness attributes |
| [GOOGLE_INTERVIEW_DETAILED.md](GOOGLE_INTERVIEW_DETAILED.md) | Component-wise project breakdown (30+ components) |

### Technical Documentation
| File | Description |
|------|-------------|
| [RESUME_TO_CODE_MAPPING.md](RESUME_TO_CODE_MAPPING.md) | Resume bullets mapped to actual code & architecture |
| [SYSTEM_INTERCONNECTIVITY.md](SYSTEM_INTERCONNECTIVITY.md) | Complete system architecture diagrams |
| [BEAT_ADVANCED_FEATURES.md](BEAT_ADVANCED_FEATURES.md) | ML/Statistics features - Gradient Descent, YAKE, etc. |

### Project Analysis
| File | Description |
|------|-------------|
| [MASTER_PORTFOLIO_SUMMARY.md](MASTER_PORTFOLIO_SUMMARY.md) | Overview of all 6 projects |
| [ANALYSIS_beat.md](ANALYSIS_beat.md) | Beat - Data aggregation service (Python) |
| [ANALYSIS_stir.md](ANALYSIS_stir.md) | Stir - Data platform (Airflow + dbt) |
| [ANALYSIS_event_grpc.md](ANALYSIS_event_grpc.md) | Event-grpc - Event consumer (Go) |
| [ANALYSIS_coffee.md](ANALYSIS_coffee.md) | Coffee - REST API service (Go) |
| [ANALYSIS_fake_follower_analysis.md](ANALYSIS_fake_follower_analysis.md) | Fake follower ML detection system |
| [ANALYSIS_saas_gateway.md](ANALYSIS_saas_gateway.md) | SaaS Gateway service |

---

## Key Numbers to Remember

| Metric | Value |
|--------|-------|
| Daily data points processed | 10M+ |
| Worker processes | 150+ |
| Airflow DAGs | 76 |
| dbt models | 112 |
| Events per second | 10,000+ |
| Log retrieval speedup | 2.5x |
| Data latency reduction | 50% |
| Indic scripts supported | 10 |

---

## Top 5 Stories

| Story | Use For |
|-------|---------|
| **Event-driven architecture** | Leadership, Technical decisions, Accomplishment |
| **GPT timeout failure** | Failure, Learning from mistakes |
| **MongoDB vs ClickHouse debate** | Conflict resolution, Data-driven decisions |
| **Fake follower detection** | Ambiguity, Innovation, Customer impact |
| **Junior engineer mentoring** | Helping others, Mentorship |

---

## 6 Googleyness Attributes

1. **Thriving in Ambiguity** - Fake follower detection with no training data
2. **Valuing Feedback** - Rate limiting code refactored after review
3. **Challenging Status Quo** - dbt over Fivetran
4. **Putting User First** - Data integrity over cheaper API
5. **Doing the Right Thing** - Reliability over new features
6. **Caring About Team** - Teaching junior engineer debugging

---

## System Architecture

```
┌──────────────┐     ┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│     BEAT     │────▶│   RabbitMQ   │────▶│  EVENT-GRPC  │────▶│  ClickHouse  │
│   (Python)   │     │   (Events)   │     │     (Go)     │     │   (OLAP)     │
│  150+ workers│     │              │     │  26 queues   │     │              │
└──────────────┘     └──────────────┘     └──────────────┘     └──────┬───────┘
                                                                       │
                                                                       ▼
┌──────────────┐     ┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│    COFFEE    │◀────│  PostgreSQL  │◀────│      S3      │◀────│     STIR     │
│  (Go API)    │     │   (OLTP)     │     │  (Staging)   │     │ (Airflow+dbt)│
└──────────────┘     └──────────────┘     └──────────────┘     └──────────────┘
```

---

## Projects Covered

| Project | Tech Stack | My Role |
|---------|------------|---------|
| **beat** | Python, FastAPI, asyncio, Redis | Core Developer |
| **stir** | Airflow, dbt, ClickHouse, PostgreSQL | Core Developer |
| **event-grpc** | Go, RabbitMQ, ClickHouse | ClickHouse sinker implementation |
| **fake_follower** | Python, AWS Lambda, ML | Solo Developer |
| **coffee** | Go, PostgreSQL, Redis | Related project |
| **saas-gateway** | Go | Related project |

---

*Best of luck for the interview!* 🚀

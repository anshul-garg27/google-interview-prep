# Walmart Googleyness Questions - Quick Reference Index

**Total Questions**: 10+ Googleyness questions with detailed STAR answers

**Source Files**:
- WALMART_GOOGLEYNESS_QUESTIONS.md
- WALMART_INTERVIEW_ALL_QUESTIONS.md
- WALMART_HIRING_MANAGER_GUIDE.md

---

## Table of Contents

1. [Thriving in Ambiguity (7 questions)](#1-thriving-in-ambiguity-7-questions)
2. [Valuing Feedback (2 questions)](#2-valuing-feedback-2-questions)
3. [Challenging Status Quo (1+ question)](#3-challenging-status-quo-1-question)

---

## 1. Thriving in Ambiguity (7 questions)

**Definition**: Comfort with uncertainty. Able to make progress when requirements are unclear, specifications change, or the path forward is uncertain.

**Walmart Work Examples**: DC Inventory Search (no API existed), Kafka Connect custom SMT (no documentation), Multi-region architecture (no precedent)

### Questions:

- [Q1.1: "Tell me about a time you had to solve a problem with incomplete information."](#q11-tell-me-about-a-time-you-had-to-solve-a-problem-with-incomplete-information)
  - **Story**: DC Inventory Search - no API existed, reverse-engineered EI endpoints
  - **Result**: 4 weeks delivery (vs 12 weeks estimated), 30K+ queries/day

- [Q1.2: "Describe a situation where requirements changed mid-project."](#q12-describe-a-situation-where-requirements-changed-mid-project)
  - **Story**: Audit logging scope: 2 endpoints → 47 endpoints mid-project
  - **Result**: Created filter pattern, zero code changes in consuming services

- [Q1.3: "How do you approach a problem where the solution isn't obvious?"](#q13-how-do-you-approach-a-problem-where-the-solution-isnt-obvious)
  - **Story**: Spring Boot 3 migration - 200+ failing tests, systematic triage
  - **Result**: 200+ failures → 0 failures in 24 hours

- [Q1.4: "Tell me about a time you had to make a decision without all the data you wanted."](#q14-tell-me-about-a-time-you-had-to-make-a-decision-without-all-the-data-you-wanted)
  - **Story**: Multi-region Kafka architecture - no RPO/RTO requirements known
  - **Result**: Active-Active pattern, 0 data loss, <30s failover

- [Q1.5: "Describe a time you navigated ambiguous stakeholder requirements."](#q15-describe-a-time-you-navigated-ambiguous-stakeholder-requirements)
  - **Story**: DSD notification system - conflicting stakeholder needs
  - **Result**: Event-driven architecture, 500K+ notifications, 5 consumers added post-launch

- [Q1.6: "How do you handle situations where you don't have precedent to follow?"](#q16-how-do-you-handle-situations-where-you-dont-have-precedent-to-follow)
  - **Story**: Multi-market architecture (US/CA/MX) - no team precedent
  - **Result**: Site-based partitioning, 8M+ queries/month across 3 markets

- [Q1.7: "Describe a time you had to pivot your approach mid-execution."](#q17-describe-a-time-you-had-to-pivot-your-approach-mid-execution)
  - **Story**: Spring Boot 3 migration - pivoted from big-bang to phased approach
  - **Result**: Delivered on time, zero production incidents

---

## 2. Valuing Feedback (2 questions)

**Definition**: Seeks out and incorporates feedback from others. Actively listens to diverse perspectives. Adjusts approach based on input.

**Walmart Work Examples**: Code review improvements (CompletableFuture thread pool), Multi-region architecture feedback from 3 experts

### Questions:

- [Q2.1: "Tell me about a time you received critical feedback and how you responded."](#q21-tell-me-about-a-time-you-received-critical-feedback-and-how-you-responded)
  - **Story**: CompletableFuture memory leak - architect's critical feedback
  - **Result**: Thread pool isolation, 0% error rate, pattern adopted by 3 teams

- [Q2.2: "Describe a time you asked for feedback and how you incorporated it."](#q22-describe-a-time-you-asked-for-feedback-and-how-you-incorporated-it)
  - **Story**: Multi-region Kafka - asked 3 experts (architect, SRE, platform owner)
  - **Result**: Async dual-write pattern, automatic failover, deduplication added

---

## 3. Challenging Status Quo (1+ question)

**Definition**: Questions assumptions and proposes better ways of doing things. Doesn't accept 'that's how we've always done it.'

**Walmart Work Examples**: Kafka vs PostgreSQL for audit logs, CompletableFuture vs ParallelStream

### Questions:

- [Q3.1: "Tell me about a time you challenged the way things were done."](#q31-tell-me-about-a-time-you-challenged-the-way-things-were-done)
  - **Story**: Challenged PostgreSQL for audit logs, proposed Kafka + GCS + BigQuery
  - **Result**: 95% faster writes, 90% cost reduction ($5K → $500/month), 25x scalability

---

## Links to Full Answers

All detailed STAR answers with technical depth, metrics, and learnings are available in:

### [WALMART_GOOGLEYNESS_QUESTIONS.md](WALMART_GOOGLEYNESS_QUESTIONS.md)
- Full STAR format answers (Situation, Task, Action, Result, Learning)
- Technical implementation details
- Code examples
- Architecture diagrams
- Follow-up question handling

### [WALMART_INTERVIEW_ALL_QUESTIONS.md](WALMART_INTERVIEW_ALL_QUESTIONS.md)
- Additional behavioral questions (60+ total)
- Cross-references to Googleyness stories
- Production debugging examples
- Team collaboration stories

### [WALMART_HIRING_MANAGER_GUIDE.md](WALMART_HIRING_MANAGER_GUIDE.md)
- System design deep dives
- Architecture walkthroughs
- Scale and performance analysis
- Failure scenarios and resilience

---

## Quick Reference: Key Metrics

| System | Metric | Value |
|--------|--------|-------|
| **DC Inventory Search** | Queries/day | 30,000+ |
| **DC Inventory Search** | Latency P95 | 1.8s |
| **Kafka Audit System** | Events/day | 2M+ |
| **Kafka Audit System** | Cost savings | $54K/year |
| **Multi-region Kafka** | Failover time | <30 seconds |
| **Multi-region Kafka** | RPO | 0 seconds |
| **Spring Boot 3 Migration** | Services migrated | 6 services |
| **Spring Boot 3 Migration** | Test failures fixed | 200+ |
| **DSD Notifications** | Notifications sent | 500K+ |
| **Multi-market Architecture** | Markets supported | US/CA/MX |
| **Multi-market Architecture** | Queries/month | 8M+ |

---

## How to Use This Index

1. **For Interview Prep**: Review questions by category, practice STAR answers
2. **For Quick Reference**: Use metrics table for impact numbers
3. **For Deep Dive**: Click links to full documents for technical details
4. **For Follow-ups**: Each question in main docs has follow-up handling

---

**Last Updated**: February 2026
**Questions Coverage**: Googleyness attributes 1-3 (Ambiguity, Feedback, Status Quo)
**Note**: Document continues with Q3.2-Q10.5 covering remaining Googleyness attributes and Leadership Principles

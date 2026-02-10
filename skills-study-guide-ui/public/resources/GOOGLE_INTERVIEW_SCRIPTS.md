# GOOGLE INTERVIEW - EXACT SCRIPTS
## Word-by-Word Kya Bolna Hai

---

# OPENING: "TELL ME ABOUT YOURSELF"

## 90-Second Script (Practice This!)

```
"Hi, I'm Anshul. I'm a software engineer with about 3 years of experience,
currently at Good Creator Co, which is India's largest influencer marketing platform.

My work has been primarily in three areas:

FIRST, I built 'beat', our data aggregation service. It scrapes Instagram and
YouTube data at scale - we process about 10 million data points daily. I designed
the entire worker pool architecture - 150 concurrent workers, multi-level rate
limiting, and integrations with 15+ external APIs.

SECOND, I built 'stir', our data platform using Airflow and dbt. 76 production DAGs,
112 dbt models, processing billions of records. This reduced our data latency by 50%.

THIRD, I built a fake follower detection system from scratch - an ML ensemble that
supports 10 Indian languages and runs on serverless architecture.

What excites me about Google is the scale of impact - building systems that serve
billions of users - and the engineering culture of learning from the best.

I'd love to hear more about the team and the challenges you're working on."
```

**Time**: 90 seconds exactly
**Practice**: Record yourself, time it

---

# GOOGLEYNESS QUESTIONS - SCRIPT BY SCRIPT

## Question 1: "Tell me about a time you led a team through a difficult situation"

### Script (2-3 minutes):

```
"Sure, let me tell you about redesigning our data pipeline at Good Creator Co.

SITUATION:
We were storing all our influencer data - profiles, posts, engagement metrics -
directly in PostgreSQL. As we scaled to 10 million daily data points, we started
hitting problems. Query times were increasing, we were losing time-series granularity,
and the database was becoming a bottleneck.

TASK:
I had to design a new architecture that could handle this scale. I also had to lead
the implementation across three different services with different tech stacks.

ACTION:
First, I proposed an event-driven architecture. Instead of direct database writes,
we'd publish events to RabbitMQ. I built the consumer service in Go - called
event-grpc - that batches 1000 events and flushes to ClickHouse every 5 seconds.

The challenging part was getting buy-in. Some team members were comfortable with
the existing PostgreSQL approach. So I created a proof-of-concept showing that
ClickHouse could run analytics queries 50x faster than PostgreSQL for our use case.

I also documented the migration path carefully - we did a phased rollout, starting
with non-critical events, so we could verify reliability before migrating everything.

I coordinated across three services:
- beat publishes events to RabbitMQ
- event-grpc consumes and writes to ClickHouse
- stir transforms the data using dbt

RESULT:
We achieved 2.5x faster log retrieval times. The system handles 10,000+ events per
second with zero data loss. And we enabled time-series analytics that weren't
possible before - like tracking follower growth over time.

The architecture has been running in production for over 15 months now with minimal
incidents."
```

---

## Question 2: "Tell me about a time you failed"

### Script (2 minutes):

```
"I'll tell you about a production incident I caused with our GPT integration.

SITUATION:
I had added OpenAI GPT integration to our data service for inferring creator
demographics from their profile bios. In testing, it worked perfectly - fast
responses, accurate results.

WHAT WENT WRONG:
When we deployed to production, about 30% of requests started timing out.
GPT API had variable latency - sometimes 2 seconds, sometimes 45 seconds.
I hadn't accounted for this variability. In my testing, I always got quick
responses, but production load was different.

THE FAILURE:
I should have load-tested with realistic conditions. I was excited about
the feature and rushed it to production without proper testing.

WHAT I LEARNED:
Three key lessons:
First, external APIs have unpredictable behavior - always test with realistic load.
Second, features should be degradable - the system should work without optional
components.
Third, always set explicit timeouts for every external call.

WHAT I DID TO FIX IT:
I added a 30-second timeout with a circuit breaker pattern. Then I made GPT
enrichment asynchronous - moved it to a separate worker queue. Now the main
system works without GPT data and enriches in the background.

The integration now runs with 95%+ success rate, and more importantly,
I apply these learnings to every external integration we build."
```

---

## Question 3: "Tell me about a time you had a conflict with a colleague"

### Script (2 minutes):

```
"I had a technical disagreement with a senior engineer about our database choice.

SITUATION:
We were designing the analytics platform. The senior engineer strongly advocated
for MongoDB because of his expertise with it. I believed ClickHouse was better
for our OLAP workload - aggregation queries over billions of rows.

HOW IT STARTED:
In a design review, he said 'MongoDB can handle this easily.' I said 'I think we
need columnar storage for these queries.' The discussion got a bit heated.

HOW I HANDLED IT:
First, I stepped back and asked to understand his perspective: 'Help me understand
why MongoDB fits here. What advantages do you see?'

He explained: familiarity, document flexibility, and faster development.

Then I proposed an experiment rather than arguing: 'What if we benchmark both with
our actual queries? Let the data decide.'

We tested with 100 million rows:
- MongoDB took 45 seconds for our typical aggregation query
- ClickHouse took 0.8 seconds for the same query

I presented results objectively, and I acknowledged his valid points: 'You're right
that MongoDB is more flexible for schema changes. But for our analytics use case,
the performance difference is 50x.'

RESULT:
We went with ClickHouse. The senior engineer actually became an advocate after
seeing the production performance. Our relationship improved because I respected
his input and let data make the decision."
```

---

## Question 4: "Tell me about navigating ambiguity"

### Script (2 minutes):

```
"Building the fake follower detection system was full of ambiguity.

SITUATION:
Brands wanted to know which influencer followers were fake. But there was no clear
definition of 'fake' - is it bots? Inactive accounts? Purchased followers? And
we had no labeled training data. Plus, followers could have names in 10 different
Indian languages - Hindi, Bengali, Tamil, and so on.

TASK:
Build an ML system with high accuracy despite having no ground truth to train on.

ACTION:
First, I decomposed the ambiguous problem into concrete, observable signals.

Instead of trying to define 'fake', I identified patterns:
- Non-Indic scripts like Greek or Chinese in an Indian influencer's followers = suspicious
- More than 4 digits in username = suspicious
- Username doesn't match display name = suspicious

For the multi-language challenge, I built a transliteration pipeline supporting
10 Indic scripts using HMM models.

I designed a scoring system with 3 confidence levels - 0.0, 0.33, 1.0 - instead
of binary fake or real. This acknowledged the inherent uncertainty.

To validate, I manually checked 500 accounts where I could verify authenticity
through other signals.

RESULT:
Achieved about 85% accuracy on the validated accounts. The system now processes
millions of followers and gives brands actionable insights.

The key was breaking down an ambiguous goal into concrete, measurable signals."
```

---

## Question 5: "Tell me about receiving critical feedback"

### Script (90 seconds):

```
"Early in building beat, a senior engineer reviewed my rate limiting code
and said it was 'too complex and hard to maintain.'

MY INITIAL REACTION:
Honestly, I felt defensive. I had worked hard on it.

WHAT I DID:
First, I took a day before responding. I re-read my code with fresh eyes.

And I realized they were right. My rate limiting had 5 different classes,
complex inheritance, and was hard to follow.

I asked them to pair with me and we refactored it together. The new version
used simple stacked context managers - much cleaner.

RESULT:
Code became simpler, easier to test, and easier for others to understand.

WHAT I LEARNED:
'Complex' isn't impressive. Simple is hard and valuable. Now I actively
seek code review feedback and specifically ask 'Is this too complex?'"
```

---

## Question 6: "Why Google?"

### Script (60 seconds):

```
"Three reasons:

FIRST, Scale. At Good Creator Co, I built systems handling 10 million daily
data points. At Google, I'd work on systems serving billions of users. That's
the scale of impact I want.

SECOND, Learning. Google's engineering culture is legendary. The opportunity
to learn from engineers who've built YouTube, Search, Cloud - that's incredible.

THIRD, Growth. I've been the one building systems from scratch. I want to be
in an environment where I can also learn from existing world-class systems.

What I bring is end-to-end ownership experience. I've built complete systems -
data pipelines, real-time processing, ML models. I know both Python and Go.
And I have the startup scrappiness - bias to action, doing more with less.

I'm excited about the team you're building and would love to contribute."
```

---

# SITUATIONAL QUESTIONS - SCRIPTS

## "How would you handle a disagreement with your manager?"

```
"I'd approach it with data and curiosity, not defensiveness.

First, I'd make sure I understand their perspective. Maybe they have context
I don't have.

Then, if I still disagree, I'd present my view with evidence. Not 'I think
this is wrong' but 'Based on these metrics, I believe option B would be better
because X, Y, Z.'

I'd propose an experiment if possible: 'Can we try my approach on a small
scale and measure the results?'

And ultimately, if they decide to go a different direction after hearing
my input, I'd commit to that decision fully. Disagree and commit.

I actually did this at Good Creator Co when advocating for dbt over Fivetran.
I presented the trade-offs, offered to prove it with a POC, and the team
agreed. But if they hadn't, I would have committed to Fivetran."
```

---

## "What would you do if you missed a deadline?"

```
"Three things:

FIRST, communicate early. As soon as I see risk, I tell stakeholders.
Not on the deadline day, but as soon as I know.

SECOND, explain why with specifics. Not 'it's taking longer' but 'I found
edge cases that need handling to avoid data corruption.'

THIRD, propose a plan. 'I need one more week, and here's exactly what I'll
deliver and why it's worth the wait.'

I actually experienced this with our ClickHouse sync pipeline. At 10 days
into a 2-week estimate, I found edge cases I hadn't anticipated. I
communicated immediately, explained the risks of not handling them, and
we agreed on 1 extra week. The result was a robust pipeline with proper
error handling and rollback capability."
```

---

## "How do you prioritize when everything is urgent?"

```
"I use a simple framework:

FIRST, what unblocks others? If my delay blocks the whole team, that's
highest priority.

SECOND, what has customer impact? External commitments over internal
improvements.

THIRD, what's the impact-to-effort ratio? Quick wins that deliver big
value come before big efforts with uncertain value.

For example, at Good Creator Co, I had to choose between a new leaderboard
feature and fixing DAG reliability. Product wanted the feature urgently.

But I realized: failing DAGs blocked the entire data team. New features
on an unstable foundation would just create more problems.

I advocated for reliability first. We fixed DAG failures in 2 weeks, then
shipped leaderboard. Product later acknowledged this was the right call."
```

---

# CLOSING QUESTIONS TO ASK

## When They Ask "Do you have questions for me?"

### Question 1:
```
"What does a typical project look like for an L4 engineer on your team?
I'm curious about the scope and the kind of ownership expected."
```

### Question 2:
```
"How does the team handle technical disagreements? Is there a culture of
written design docs, or more informal discussions?"
```

### Question 3:
```
"What's an example of how the team navigated ambiguity recently? I'm
interested in how decisions get made when requirements aren't clear."
```

### Question 4:
```
"What does success look like in the first 6 months for someone in this role?"
```

---

# EMERGENCY SCRIPTS

## If You Don't Know the Answer:

```
"That's a great question. I haven't faced that exact situation, but
let me think about how I'd approach it..."

[Then use a related experience or talk through your reasoning]
```

## If You Need Time to Think:

```
"That's an interesting question. Let me take a moment to think of the
best example..."

[5-second pause is fine!]
```

## If You're Rambling:

```
"Let me summarize - the key point is [one sentence summary].
Does that answer your question, or would you like me to go deeper?"
```

## If Interviewer Interrupts:

```
[Stop immediately]
"Of course, what would you like me to clarify?"
```

---

# TIMING GUIDE

| Question Type | Target Time |
|--------------|-------------|
| "Tell me about yourself" | 90 seconds |
| STAR story | 2-3 minutes |
| Follow-up question | 30-60 seconds |
| "Why Google?" | 60 seconds |
| Your questions | 5 minutes total |

---

# PRACTICE CHECKLIST

- [ ] Record yourself answering each question
- [ ] Time each answer
- [ ] Watch recording - check for filler words (um, uh, like)
- [ ] Practice with a friend for mock interview
- [ ] Do one full mock interview (all questions) before tomorrow

---

**REMEMBER: You have REAL experience. Just tell YOUR stories with confidence!**

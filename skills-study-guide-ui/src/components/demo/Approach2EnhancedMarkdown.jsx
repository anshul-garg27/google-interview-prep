import { useState } from 'react'
import ReactMarkdown from 'react-markdown'
import remarkGfm from 'remark-gfm'
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter'
import { oneDark, oneLight } from 'react-syntax-highlighter/dist/esm/styles/prism'
import { useTheme } from '../../context/ThemeContext'
import { ChevronRight, Copy, Check, Briefcase, Eye, EyeOff } from 'lucide-react'

// Same content but the RENDERER detects patterns and enhances them
const SAMPLE_MARKDOWN = `## Consumer Groups and Rebalancing

### Consumer Group Basics

A **consumer group** is a set of consumers that cooperate to consume messages from one or more topics. Kafka guarantees that each partition is consumed by exactly one consumer within a group.

Key rules:
- Each partition is assigned to exactly **one consumer** in the group
- A consumer can be assigned **multiple partitions**
- If consumers > partitions, some consumers will be **idle**
- If consumers < partitions, some consumers handle **multiple partitions**

\`\`\`java
@KafkaListener(topics = "audit-events", groupId = "audit-processor")
public void consume(ConsumerRecord<String, AuditEvent> record) {
    log.info("Partition: {}, Offset: {}, Key: {}",
        record.partition(), record.offset(), record.key());
    auditRepository.save(record.value());
}
\`\`\`

### Rebalance Protocols

| Protocol | Description | Use Case |
|----------|-------------|----------|
| **Eager** | All consumers stop, all partitions revoked, full reassignment | Legacy, simple |
| **Cooperative Sticky** | Only affected partitions move, minimal disruption | Production recommended |
| **Range** | Partitions assigned in ranges per topic | Co-partitioned topics |
| **RoundRobin** | Partitions distributed evenly across consumers | General purpose |
`

function EnhancedCodeBlock({ language, children }) {
  const [copied, setCopied] = useState(false)
  const { darkMode } = useTheme()
  const code = String(children).replace(/\n$/, '')

  const LABELS = { java: 'Java', python: 'Python', go: 'Go', sql: 'SQL' }

  return (
    <div className="rounded-lg border border-[var(--border)] overflow-hidden my-4 not-prose">
      <div className="flex items-center justify-between px-3 py-2 bg-[var(--bg-surface-alt)] border-b border-[var(--border)]">
        <span className="font-mono text-[10px] font-semibold uppercase tracking-[0.06em] text-[var(--text-tertiary)]">
          {LABELS[language] || language}
        </span>
        <button
          onClick={() => { navigator.clipboard.writeText(code); setCopied(true); setTimeout(() => setCopied(false), 2000) }}
          className="flex items-center gap-1 text-[10px] font-mono text-[var(--text-tertiary)] hover:text-[var(--text-primary)] transition-colors"
        >
          {copied ? <><Check size={12} className="text-emerald-500" /> Copied</> : <><Copy size={12} /> Copy</>}
        </button>
      </div>
      <SyntaxHighlighter style={darkMode ? oneDark : oneLight} language={language} PreTag="div"
        customStyle={{ margin: 0, padding: '1rem', fontSize: '13px', lineHeight: '1.7', background: 'var(--code-bg)' }}>
        {code}
      </SyntaxHighlighter>
    </div>
  )
}

function QAAccordionItem({ question, answer }) {
  const [open, setOpen] = useState(false)
  return (
    <div className="border border-[var(--border)] rounded-lg overflow-hidden not-prose">
      <button onClick={() => setOpen(!open)}
        className="w-full flex items-start gap-3 p-4 text-left hover:bg-[var(--bg-hover)] transition-colors">
        <ChevronRight size={16} className={`shrink-0 mt-0.5 text-[var(--accent)] transition-transform duration-200 ${open ? 'rotate-90' : ''}`} />
        <span className="font-body font-semibold text-[15px] leading-relaxed">{question}</span>
      </button>
      {open && (
        <div className="px-4 pb-4 pl-11 border-t border-[var(--border)] bg-[var(--bg-surface-alt)]/50 animate-fade-in">
          <p className="pt-3 text-[var(--text-secondary)] text-[14px] leading-relaxed">{answer}</p>
        </div>
      )}
    </div>
  )
}

function CalloutBox({ children }) {
  return (
    <div className="my-6 rounded-r-lg border-l-4 border-l-[var(--accent)] bg-[var(--accent-subtle)] p-4 not-prose">
      <div className="flex items-center gap-2 mb-2">
        <Briefcase size={14} className="text-[var(--accent)]" />
        <span className="font-mono text-[10px] font-semibold uppercase tracking-[0.06em] text-[var(--accent)]">
          Real-World Application
        </span>
      </div>
      <div className="text-[14px] leading-relaxed">{children}</div>
    </div>
  )
}

export default function Approach2EnhancedMarkdown() {
  const [allQAOpen, setAllQAOpen] = useState(false)

  return (
    <div>
      {/* Approach label */}
      <div className="mb-8 p-4 rounded-xl border border-[var(--border)] bg-[var(--bg-surface)]">
        <div className="flex items-center gap-3 mb-2">
          <span className="font-mono text-[10px] font-semibold uppercase tracking-[0.06em] px-2 py-1 rounded bg-blue-100 dark:bg-blue-900/30 text-blue-600 dark:text-blue-400">
            Approach 2
          </span>
          <span className="font-heading font-semibold">Enhanced Markdown</span>
        </div>
        <p className="text-[13px] text-[var(--text-secondary)] leading-relaxed">
          Same .md file but the <strong>renderer detects patterns</strong> and enhances them: Q&A becomes accordions,
          "How Anshul Used It" becomes callout boxes, code blocks get headers + copy buttons, tables get sticky columns.
          <strong> Zero content changes needed.</strong>
        </p>
        <div className="flex gap-4 mt-3 font-mono text-[10px] text-[var(--text-tertiary)]">
          <span>Effort: ~2 weeks</span>
          <span>Interactivity: Medium</span>
          <span>Content changes: None</span>
        </div>
      </div>

      {/* Rendered content */}
      <div className="prose prose-gray dark:prose-invert font-body">
        <ReactMarkdown
          remarkPlugins={[remarkGfm]}
          components={{
            code({ inline, className, children, ...props }) {
              const match = /language-(\w+)/.exec(className || '')
              if (inline || !match) {
                return <code className="font-mono text-[0.875em] px-1.5 py-0.5 rounded bg-[var(--bg-surface-alt)] border border-[var(--border)]" {...props}>{children}</code>
              }
              return <EnhancedCodeBlock language={match[1]}>{children}</EnhancedCodeBlock>
            },
            table({ children }) {
              return (
                <div className="overflow-x-auto rounded-lg border border-[var(--border)] my-4 not-prose">
                  <table className="w-full text-sm">{children}</table>
                </div>
              )
            },
            th({ children }) {
              return <th className="px-4 py-3 text-left bg-[var(--bg-surface-alt)] font-mono text-[11px] uppercase tracking-[0.06em] font-semibold text-[var(--text-tertiary)]">{children}</th>
            },
            td({ children }) {
              return <td className="px-4 py-3 border-t border-[var(--border)] text-[13px]">{children}</td>
            },
          }}
        >
          {SAMPLE_MARKDOWN}
        </ReactMarkdown>
      </div>

      {/* Callout — detected from "How Anshul Used It" heading */}
      <CalloutBox>
        <p>At Walmart, we configured consumer groups for the audit logging pipeline:</p>
        <ul className="list-disc pl-5 mt-2 space-y-1 text-[13px]">
          <li><strong>6 partitions</strong> per topic for parallel processing</li>
          <li><strong>Cooperative Sticky</strong> rebalancing to minimize disruption during deployments</li>
          <li>Consumer group ID pattern: <code className="font-mono text-xs px-1 py-0.5 rounded bg-[var(--bg-surface-alt)]">audit-processor-{'{region}'}</code> (US, CA, MX)</li>
          <li>Used <code className="font-mono text-xs px-1 py-0.5 rounded bg-[var(--bg-surface-alt)]">session.timeout.ms=30000</code> and <code className="font-mono text-xs px-1 py-0.5 rounded bg-[var(--bg-surface-alt)]">heartbeat.interval.ms=10000</code></li>
        </ul>
      </CalloutBox>

      {/* Q&A Section — detected from **Q1:** pattern */}
      <div className="mt-8">
        <div className="flex items-center justify-between mb-4">
          <span className="font-mono text-[11px] font-semibold uppercase tracking-[0.06em] text-[var(--accent)]">
            Interview Q&A
          </span>
          <button onClick={() => setAllQAOpen(!allQAOpen)}
            className="flex items-center gap-1.5 font-mono text-[11px] text-[var(--text-tertiary)] hover:text-[var(--text-secondary)] transition-colors">
            {allQAOpen ? <EyeOff size={12} /> : <Eye size={12} />}
            {allQAOpen ? 'Collapse All' : 'Expand All'}
          </button>
        </div>
        <div className="space-y-3">
          <QAAccordionItem
            question="Q1: What is a consumer group and why do we need it?"
            answer="A consumer group is a set of consumers that cooperate to consume data from topics. We need it for: (1) Parallel processing — each partition is consumed by one consumer, enabling horizontal scaling. (2) Fault tolerance — if a consumer dies, its partitions are reassigned to others. (3) Load balancing — partitions are evenly distributed across consumers."
          />
          <QAAccordionItem
            question="Q2: What happens when a consumer dies in the middle of processing?"
            answer="When a consumer stops sending heartbeats (exceeds session.timeout.ms), the group coordinator triggers a rebalance. The dead consumer's partitions are reassigned to surviving consumers. If auto-commit was used, some messages may be reprocessed (at-least-once). If manual commit was used, only uncommitted messages are reprocessed."
          />
        </div>
      </div>
    </div>
  )
}

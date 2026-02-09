import ReactMarkdown from 'react-markdown'
import remarkGfm from 'remark-gfm'
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter'
import { oneDark, oneLight } from 'react-syntax-highlighter/dist/esm/styles/prism'
import { useTheme } from '../../context/ThemeContext'

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

### How Anshul Used It at Walmart

At Walmart, we configured consumer groups for the audit logging pipeline:
- **6 partitions** per topic for parallel processing
- **Cooperative Sticky** rebalancing to minimize disruption during deployments
- Consumer group ID pattern: \`audit-processor-{region}\` (US, CA, MX)
- Used \`session.timeout.ms=30000\` and \`heartbeat.interval.ms=10000\`

**Q1: What is a consumer group and why do we need it?**

A consumer group is a set of consumers that cooperate to consume data from topics. We need it for: (1) Parallel processing — each partition is consumed by one consumer, enabling horizontal scaling. (2) Fault tolerance — if a consumer dies, its partitions are reassigned to others. (3) Load balancing — partitions are evenly distributed across consumers.

**Q2: What happens when a consumer dies in the middle of processing?**

When a consumer stops sending heartbeats (exceeds session.timeout.ms), the group coordinator triggers a rebalance. The dead consumer's partitions are reassigned to surviving consumers. If auto-commit was used, some messages may be reprocessed (at-least-once). If manual commit was used, only uncommitted messages are reprocessed.
`

export default function Approach1RawMarkdown() {
  const { darkMode } = useTheme()

  return (
    <div>
      {/* Approach label */}
      <div className="mb-8 p-4 rounded-xl border border-[var(--border)] bg-[var(--bg-surface)]">
        <div className="flex items-center gap-3 mb-2">
          <span className="font-mono text-[10px] font-semibold uppercase tracking-[0.06em] px-2 py-1 rounded bg-gray-200 dark:bg-gray-700 text-gray-600 dark:text-gray-300">
            Approach 1
          </span>
          <span className="font-heading font-semibold">Raw Markdown Rendering</span>
        </div>
        <p className="text-[13px] text-[var(--text-secondary)] leading-relaxed">
          Current approach: <code className="text-xs font-mono px-1.5 py-0.5 rounded bg-[var(--bg-surface-alt)] border border-[var(--border)]">react-markdown</code> + <code className="text-xs font-mono px-1.5 py-0.5 rounded bg-[var(--bg-surface-alt)] border border-[var(--border)]">remark-gfm</code>.
          The .md file is fetched and rendered as-is. No interactivity. No special treatment for Q&A, callouts, or diagrams beyond mermaid.
        </p>
        <div className="flex gap-4 mt-3 font-mono text-[10px] text-[var(--text-tertiary)]">
          <span>Effort: None (current)</span>
          <span>Interactivity: Minimal</span>
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
              return (
                <SyntaxHighlighter style={darkMode ? oneDark : oneLight} language={match[1]} PreTag="div"
                  customStyle={{ margin: 0, padding: '1rem', fontSize: '13px', lineHeight: '1.7', borderRadius: '0.5rem' }}>
                  {String(children).replace(/\n$/, '')}
                </SyntaxHighlighter>
              )
            }
          }}
        >
          {SAMPLE_MARKDOWN}
        </ReactMarkdown>
      </div>
    </div>
  )
}

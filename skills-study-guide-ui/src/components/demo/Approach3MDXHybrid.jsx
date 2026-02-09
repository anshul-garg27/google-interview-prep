import { useState, useMemo } from 'react'
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter'
import { oneDark, oneLight } from 'react-syntax-highlighter/dist/esm/styles/prism'
import { useTheme } from '../../context/ThemeContext'
import { Copy, Check, Briefcase, ChevronRight, Play, RotateCcw } from 'lucide-react'

// This simulates what MDX would look like:
// Markdown prose rendered normally, but with React components embedded inline

function AnnotatedCodeBlock({ language, annotations, children }) {
  const [copied, setCopied] = useState(false)
  const [hoveredLine, setHoveredLine] = useState(null)
  const { darkMode } = useTheme()
  const code = children.trim()
  const lines = code.split('\n')

  const annotationMap = useMemo(() => {
    const map = {}
    annotations?.forEach(a => { map[a.line] = a.text })
    return map
  }, [annotations])

  return (
    <div className="rounded-lg border border-[var(--border)] overflow-hidden my-4 not-prose">
      <div className="flex items-center justify-between px-3 py-2 bg-[var(--bg-surface-alt)] border-b border-[var(--border)]">
        <span className="font-mono text-[10px] font-semibold uppercase tracking-[0.06em] text-[var(--text-tertiary)]">
          {language}
        </span>
        <button
          onClick={() => { navigator.clipboard.writeText(code); setCopied(true); setTimeout(() => setCopied(false), 2000) }}
          className="flex items-center gap-1 text-[10px] font-mono text-[var(--text-tertiary)] hover:text-[var(--text-primary)] transition-colors"
        >
          {copied ? <><Check size={12} className="text-emerald-500" /> Copied</> : <><Copy size={12} /> Copy</>}
        </button>
      </div>
      <div className="relative">
        {lines.map((line, i) => (
          <div
            key={i}
            className={`flex items-start font-mono text-[13px] leading-[1.7] px-4 transition-colors cursor-pointer
              ${hoveredLine === i ? 'bg-[var(--accent-subtle)]' : ''}
              ${annotationMap[i] ? 'border-l-2 border-l-[var(--accent)]' : 'border-l-2 border-l-transparent'}`}
            onMouseEnter={() => setHoveredLine(i)}
            onMouseLeave={() => setHoveredLine(null)}
          >
            <span className="w-8 shrink-0 text-[var(--text-tertiary)] text-right mr-4 select-none text-[11px] leading-[1.85]">
              {i + 1}
            </span>
            <pre className="flex-1 overflow-x-auto" style={{ margin: 0, padding: 0, background: 'transparent' }}>
              <code>{line || ' '}</code>
            </pre>
          </div>
        ))}

        {/* Annotation tooltip */}
        {hoveredLine !== null && annotationMap[hoveredLine] && (
          <div className="absolute right-4 top-0 mt-2 max-w-[250px] p-3 rounded-lg bg-[var(--accent)] text-white text-[12px] leading-relaxed shadow-lg z-10 animate-fade-in"
            style={{ top: `${hoveredLine * 22.1 + 8}px` }}>
            {annotationMap[hoveredLine]}
          </div>
        )}
      </div>
    </div>
  )
}

function ConsumerGroupMiniViz() {
  const [consumers, setConsumers] = useState(3)
  const partitions = 6

  const assignment = useMemo(() => {
    const result = Array.from({ length: partitions }, (_, i) => ({
      id: i,
      consumer: consumers > 0 ? i % consumers : null,
    }))
    return result
  }, [consumers])

  const COLORS = ['#D97706', '#059669', '#2563EB', '#DC2626', '#7C3AED', '#0891B2']

  return (
    <div className="my-6 p-5 rounded-xl border border-[var(--border)] bg-[var(--bg-surface)] not-prose">
      <div className="flex items-center justify-between mb-4">
        <span className="font-mono text-[10px] font-semibold uppercase tracking-[0.06em] text-[var(--accent)]">
          Interactive: Consumer Group Visualizer
        </span>
        <button onClick={() => setConsumers(3)}
          className="flex items-center gap-1 font-mono text-[10px] text-[var(--text-tertiary)] hover:text-[var(--text-primary)] transition-colors">
          <RotateCcw size={10} /> Reset
        </button>
      </div>

      {/* Controls */}
      <div className="flex items-center gap-4 mb-5">
        <div className="flex items-center gap-2">
          <span className="font-mono text-[11px] text-[var(--text-secondary)]">Consumers:</span>
          <div className="flex items-center gap-1">
            <button onClick={() => setConsumers(Math.max(0, consumers - 1))}
              className="w-7 h-7 rounded border border-[var(--border)] flex items-center justify-center font-mono text-sm hover:bg-[var(--bg-hover)] transition-colors">
              -
            </button>
            <span className="w-6 text-center font-mono font-semibold text-[var(--accent)]">{consumers}</span>
            <button onClick={() => setConsumers(Math.min(8, consumers + 1))}
              className="w-7 h-7 rounded border border-[var(--border)] flex items-center justify-center font-mono text-sm hover:bg-[var(--bg-hover)] transition-colors">
              +
            </button>
          </div>
        </div>
        <span className="font-mono text-[11px] text-[var(--text-tertiary)]">
          {partitions} partitions &middot; {Math.max(0, consumers - partitions)} idle consumers
        </span>
      </div>

      {/* Visualization */}
      <div className="grid grid-cols-6 gap-2 mb-4">
        {assignment.map(p => (
          <div key={p.id}
            className="rounded-lg border-2 p-3 text-center transition-all duration-300"
            style={{
              borderColor: p.consumer !== null ? COLORS[p.consumer % COLORS.length] : 'var(--border)',
              backgroundColor: p.consumer !== null ? COLORS[p.consumer % COLORS.length] + '15' : 'transparent',
            }}>
            <div className="font-mono text-[10px] text-[var(--text-tertiary)]">P{p.id}</div>
            <div className="font-mono text-[11px] font-semibold mt-1"
              style={{ color: p.consumer !== null ? COLORS[p.consumer % COLORS.length] : 'var(--text-tertiary)' }}>
              {p.consumer !== null ? `C${p.consumer}` : '—'}
            </div>
          </div>
        ))}
      </div>

      {/* Consumer list */}
      <div className="flex flex-wrap gap-2">
        {Array.from({ length: consumers }, (_, i) => {
          const assignedPartitions = assignment.filter(p => p.consumer === i)
          return (
            <div key={i} className="flex items-center gap-1.5 px-2.5 py-1.5 rounded-lg border border-[var(--border)] bg-[var(--bg-surface-alt)]">
              <div className="w-2.5 h-2.5 rounded-full" style={{ backgroundColor: COLORS[i % COLORS.length] }} />
              <span className="font-mono text-[11px] font-semibold">C{i}</span>
              <span className="font-mono text-[10px] text-[var(--text-tertiary)]">
                ({assignedPartitions.length > 0 ? assignedPartitions.map(p => `P${p.id}`).join(', ') : 'idle'})
              </span>
            </div>
          )
        })}
      </div>

      <p className="mt-4 text-[12px] text-[var(--text-tertiary)] italic">
        Try adding more consumers than partitions (>6) to see idle consumers, or set to 0 to see unassigned partitions.
      </p>
    </div>
  )
}

export default function Approach3MDXHybrid() {
  return (
    <div>
      {/* Approach label */}
      <div className="mb-8 p-4 rounded-xl border border-[var(--border)] bg-[var(--bg-surface)]">
        <div className="flex items-center gap-3 mb-2">
          <span className="font-mono text-[10px] font-semibold uppercase tracking-[0.06em] px-2 py-1 rounded bg-emerald-100 dark:bg-emerald-900/30 text-emerald-600 dark:text-emerald-400">
            Approach 3
          </span>
          <span className="font-heading font-semibold">MDX Hybrid</span>
        </div>
        <p className="text-[13px] text-[var(--text-secondary)] leading-relaxed">
          Content authored in <code className="text-xs font-mono px-1.5 py-0.5 rounded bg-[var(--bg-surface-alt)] border border-[var(--border)]">.mdx</code> —
          90% markdown prose + 10% embedded React components. Code blocks have <strong>line-by-line hover annotations</strong>.
          Diagrams are <strong>interactive React components</strong> instead of static mermaid.
        </p>
        <div className="flex gap-4 mt-3 font-mono text-[10px] text-[var(--text-tertiary)]">
          <span>Effort: ~4 weeks</span>
          <span>Interactivity: High</span>
          <span>Content changes: Minimal (.md → .mdx)</span>
        </div>
      </div>

      {/* Rendered content — simulating MDX output */}
      <div className="prose prose-gray dark:prose-invert font-body">
        <h2 className="font-heading">Consumer Groups and Rebalancing</h2>
        <h3 className="font-heading">Consumer Group Basics</h3>
        <p>
          A <strong>consumer group</strong> is a set of consumers that cooperate to consume messages from one or more topics.
          Kafka guarantees that each partition is consumed by exactly one consumer within a group.
        </p>
        <p>Key rules:</p>
        <ul>
          <li>Each partition is assigned to exactly <strong>one consumer</strong> in the group</li>
          <li>A consumer can be assigned <strong>multiple partitions</strong></li>
          <li>If consumers &gt; partitions, some consumers will be <strong>idle</strong></li>
          <li>If consumers &lt; partitions, some consumers handle <strong>multiple partitions</strong></li>
        </ul>
      </div>

      {/* Interactive component embedded in MDX */}
      <ConsumerGroupMiniViz />

      {/* Annotated code block */}
      <AnnotatedCodeBlock
        language="Java"
        annotations={[
          { line: 0, text: '@KafkaListener — Spring annotation that creates a consumer and subscribes to the topic' },
          { line: 1, text: 'ConsumerRecord contains the key, value, partition, offset, and headers' },
          { line: 3, text: 'Always log the partition and offset — critical for debugging rebalancing issues' },
          { line: 4, text: 'Save to repository — this should be idempotent in case of reprocessing!' },
        ]}
      >
{`@KafkaListener(topics = "audit-events", groupId = "audit-processor")
public void consume(ConsumerRecord<String, AuditEvent> record) {
    log.info("Partition: {}, Offset: {}, Key: {}",
        record.partition(), record.offset(), record.key());
    auditRepository.save(record.value());
}`}
      </AnnotatedCodeBlock>

      <div className="prose prose-gray dark:prose-invert font-body mt-6">
        <h3 className="font-heading">Rebalance Protocols</h3>
      </div>

      {/* Enhanced table with hover highlighting */}
      <div className="overflow-x-auto rounded-lg border border-[var(--border)] my-4">
        <table className="w-full text-sm">
          <thead>
            <tr>
              <th className="px-4 py-3 text-left bg-[var(--bg-surface-alt)] font-mono text-[11px] uppercase tracking-[0.06em] font-semibold text-[var(--text-tertiary)]">Protocol</th>
              <th className="px-4 py-3 text-left bg-[var(--bg-surface-alt)] font-mono text-[11px] uppercase tracking-[0.06em] font-semibold text-[var(--text-tertiary)]">Description</th>
              <th className="px-4 py-3 text-left bg-[var(--bg-surface-alt)] font-mono text-[11px] uppercase tracking-[0.06em] font-semibold text-[var(--text-tertiary)]">Use Case</th>
            </tr>
          </thead>
          <tbody>
            {[
              ['Eager', 'All consumers stop, all partitions revoked, full reassignment', 'Legacy, simple'],
              ['Cooperative Sticky', 'Only affected partitions move, minimal disruption', 'Production recommended'],
              ['Range', 'Partitions assigned in ranges per topic', 'Co-partitioned topics'],
              ['RoundRobin', 'Partitions distributed evenly across consumers', 'General purpose'],
            ].map(([name, desc, use], i) => (
              <tr key={i} className="border-t border-[var(--border)] hover:bg-[var(--bg-hover)] transition-colors">
                <td className="px-4 py-3 font-semibold text-[13px]">{name}</td>
                <td className="px-4 py-3 text-[13px] text-[var(--text-secondary)]">{desc}</td>
                <td className="px-4 py-3 text-[13px] text-[var(--text-secondary)]">{use}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {/* Callout */}
      <div className="my-6 rounded-r-lg border-l-4 border-l-[var(--accent)] bg-[var(--accent-subtle)] p-4">
        <div className="flex items-center gap-2 mb-2">
          <Briefcase size={14} className="text-[var(--accent)]" />
          <span className="font-mono text-[10px] font-semibold uppercase tracking-[0.06em] text-[var(--accent)]">
            Real-World Application
          </span>
        </div>
        <p className="text-[14px] leading-relaxed">
          At Walmart, we configured <strong>6 partitions</strong> per topic with <strong>Cooperative Sticky</strong> rebalancing.
          Consumer group ID pattern: <code className="font-mono text-xs px-1 py-0.5 rounded bg-white/50 dark:bg-black/20">audit-processor-{'{region}'}</code>
        </p>
      </div>
    </div>
  )
}

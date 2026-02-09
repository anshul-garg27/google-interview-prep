import { useState, useEffect, useRef, useMemo } from 'react'
import { useTheme } from '../../context/ThemeContext'
import { Briefcase, ChevronRight, Play, Pause, SkipForward, RotateCcw, Zap, AlertTriangle, Check } from 'lucide-react'

// Full custom React — every element hand-crafted, interactive, animated

function KafkaConsumerGroupSimulator() {
  const [consumers, setConsumers] = useState(3)
  const [strategy, setStrategy] = useState('roundrobin')
  const [failedConsumer, setFailedConsumer] = useState(null)
  const [isAnimating, setIsAnimating] = useState(false)
  const partitions = 6

  const COLORS = ['#D97706', '#059669', '#2563EB', '#DC2626', '#7C3AED', '#0891B2']

  const assignment = useMemo(() => {
    const activeConsumers = Array.from({ length: consumers }, (_, i) => i).filter(i => i !== failedConsumer)
    const result = Array.from({ length: partitions }, (_, i) => ({ id: i, consumer: null }))

    if (activeConsumers.length === 0) return result

    if (strategy === 'roundrobin') {
      result.forEach((p, i) => { p.consumer = activeConsumers[i % activeConsumers.length] })
    } else if (strategy === 'range') {
      const perConsumer = Math.ceil(partitions / activeConsumers.length)
      result.forEach((p, i) => {
        const idx = Math.floor(i / perConsumer)
        p.consumer = activeConsumers[Math.min(idx, activeConsumers.length - 1)]
      })
    }
    return result
  }, [consumers, strategy, failedConsumer])

  const killConsumer = (id) => {
    setIsAnimating(true)
    setFailedConsumer(id)
    setTimeout(() => setIsAnimating(false), 600)
  }

  return (
    <div className="my-8 rounded-xl border-2 border-[var(--accent)]/30 bg-[var(--bg-surface)] overflow-hidden not-prose">
      {/* Header */}
      <div className="px-5 py-4 border-b border-[var(--border)] bg-[var(--accent-subtle)]">
        <div className="flex items-center justify-between">
          <div>
            <div className="flex items-center gap-2">
              <Zap size={16} className="text-[var(--accent)]" />
              <span className="font-heading font-semibold text-base">Kafka Consumer Group Simulator</span>
            </div>
            <p className="text-[12px] text-[var(--text-secondary)] mt-1">
              Add/remove consumers, kill them, change strategy — watch partitions redistribute
            </p>
          </div>
          <button onClick={() => { setConsumers(3); setFailedConsumer(null); setStrategy('roundrobin') }}
            className="flex items-center gap-1 px-3 py-1.5 rounded-lg border border-[var(--border)] text-[11px] font-mono hover:bg-[var(--bg-hover)] transition-colors">
            <RotateCcw size={11} /> Reset
          </button>
        </div>
      </div>

      <div className="p-5">
        {/* Controls */}
        <div className="flex flex-wrap items-center gap-6 mb-6">
          <div className="flex items-center gap-3">
            <span className="font-mono text-[11px] text-[var(--text-secondary)]">Consumers:</span>
            <div className="flex items-center gap-1">
              {[1, 2, 3, 4, 5, 6, 7, 8].map(n => (
                <button key={n} onClick={() => { setConsumers(n); setFailedConsumer(null) }}
                  className={`w-7 h-7 rounded-lg font-mono text-[12px] font-semibold transition-all
                    ${consumers === n
                      ? 'bg-[var(--accent)] text-white shadow-sm'
                      : 'border border-[var(--border)] hover:bg-[var(--bg-hover)]'}`}>
                  {n}
                </button>
              ))}
            </div>
          </div>

          <div className="flex items-center gap-2">
            <span className="font-mono text-[11px] text-[var(--text-secondary)]">Strategy:</span>
            {['roundrobin', 'range'].map(s => (
              <button key={s} onClick={() => setStrategy(s)}
                className={`px-3 py-1.5 rounded-lg font-mono text-[11px] transition-all
                  ${strategy === s
                    ? 'bg-[var(--accent)] text-white'
                    : 'border border-[var(--border)] hover:bg-[var(--bg-hover)]'}`}>
                {s === 'roundrobin' ? 'Round Robin' : 'Range'}
              </button>
            ))}
          </div>
        </div>

        {/* Topic label */}
        <div className="flex items-center gap-2 mb-3">
          <div className="h-px flex-1 bg-[var(--border)]" />
          <span className="font-mono text-[10px] text-[var(--text-tertiary)] uppercase tracking-wider">Topic: audit-events (6 partitions)</span>
          <div className="h-px flex-1 bg-[var(--border)]" />
        </div>

        {/* Partitions visualization */}
        <div className="grid grid-cols-6 gap-3 mb-6">
          {assignment.map(p => (
            <div key={p.id}
              className={`rounded-xl border-2 p-4 text-center transition-all duration-500 ${isAnimating ? 'scale-95' : 'scale-100'}`}
              style={{
                borderColor: p.consumer !== null ? COLORS[p.consumer % COLORS.length] : 'var(--border)',
                backgroundColor: p.consumer !== null ? COLORS[p.consumer % COLORS.length] + '12' : 'transparent',
                boxShadow: p.consumer !== null ? `0 2px 8px ${COLORS[p.consumer % COLORS.length]}20` : 'none',
              }}>
              <div className="font-mono text-[10px] text-[var(--text-tertiary)] mb-1">Partition</div>
              <div className="font-mono text-xl font-bold" style={{ color: p.consumer !== null ? COLORS[p.consumer % COLORS.length] : 'var(--text-tertiary)' }}>
                {p.id}
              </div>
              <div className="font-mono text-[10px] mt-1"
                style={{ color: p.consumer !== null ? COLORS[p.consumer % COLORS.length] : 'var(--text-tertiary)' }}>
                {p.consumer !== null ? `→ Consumer ${p.consumer}` : 'Unassigned'}
              </div>
            </div>
          ))}
        </div>

        {/* Consumers */}
        <div className="flex items-center gap-2 mb-3">
          <div className="h-px flex-1 bg-[var(--border)]" />
          <span className="font-mono text-[10px] text-[var(--text-tertiary)] uppercase tracking-wider">Consumer Group: audit-processor</span>
          <div className="h-px flex-1 bg-[var(--border)]" />
        </div>

        <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-3">
          {Array.from({ length: consumers }, (_, i) => {
            const isFailed = failedConsumer === i
            const myPartitions = assignment.filter(p => p.consumer === i)
            const isIdle = myPartitions.length === 0

            return (
              <div key={i}
                className={`rounded-xl border-2 p-3 transition-all duration-300
                  ${isFailed ? 'opacity-30 scale-95 border-red-400 bg-red-50 dark:bg-red-950/20' : ''}
                  ${isIdle && !isFailed ? 'border-dashed border-[var(--border)]' : ''}`}
                style={{
                  borderColor: isFailed ? '#DC2626' : isIdle ? undefined : COLORS[i % COLORS.length],
                  backgroundColor: isFailed ? undefined : isIdle ? 'transparent' : COLORS[i % COLORS.length] + '08',
                }}>
                <div className="flex items-center justify-between mb-2">
                  <div className="flex items-center gap-2">
                    <div className="w-3 h-3 rounded-full" style={{ backgroundColor: isFailed ? '#DC2626' : COLORS[i % COLORS.length] }} />
                    <span className="font-mono text-[12px] font-semibold">Consumer {i}</span>
                  </div>
                  {!isFailed && (
                    <button onClick={() => killConsumer(i)} title="Simulate failure"
                      className="p-1 rounded hover:bg-red-100 dark:hover:bg-red-900/30 transition-colors">
                      <AlertTriangle size={12} className="text-red-400" />
                    </button>
                  )}
                </div>
                <div className="font-mono text-[10px] text-[var(--text-tertiary)]">
                  {isFailed ? '☠ DEAD — partitions reassigned' :
                    isIdle ? 'Idle (no partitions)' :
                      `Partitions: ${myPartitions.map(p => `P${p.id}`).join(', ')}`}
                </div>
              </div>
            )
          })}
        </div>

        {/* Status bar */}
        <div className="mt-5 pt-4 border-t border-[var(--border)] flex flex-wrap gap-4 font-mono text-[11px] text-[var(--text-tertiary)]">
          <span>Active: {consumers - (failedConsumer !== null ? 1 : 0)}/{consumers}</span>
          <span>Idle: {Math.max(0, (consumers - (failedConsumer !== null ? 1 : 0)) - partitions)}</span>
          <span>Unassigned: {assignment.filter(p => p.consumer === null).length}</span>
          <span>Strategy: {strategy === 'roundrobin' ? 'Round Robin' : 'Range'}</span>
        </div>
      </div>
    </div>
  )
}

function InteractiveCodeWalkthrough() {
  const [step, setStep] = useState(0)
  const { darkMode } = useTheme()

  const steps = [
    { line: 0, highlight: [0], title: 'Annotation', desc: '@KafkaListener creates a MessageListenerContainer that polls messages from "audit-events" topic. The groupId determines which consumer group this consumer belongs to.' },
    { line: 1, highlight: [1], title: 'Method signature', desc: 'Spring deserializes the Kafka record into a ConsumerRecord<String, AuditEvent>. The String is the key type, AuditEvent is the value type (using Avro/JSON deserializer).' },
    { line: 2, highlight: [2, 3], title: 'Logging', desc: 'ALWAYS log partition + offset + key. This is critical for debugging. When a rebalance happens, you need to know exactly which messages were processed by which consumer.' },
    { line: 4, highlight: [4], title: 'Persistence', desc: 'Save the audit event. This MUST be idempotent! If the consumer crashes after saving but before committing the offset, this message will be redelivered.' },
  ]

  const code = [
    '@KafkaListener(topics = "audit-events", groupId = "audit-processor")',
    'public void consume(ConsumerRecord<String, AuditEvent> record) {',
    '    log.info("Partition: {}, Offset: {}, Key: {}",',
    '        record.partition(), record.offset(), record.key());',
    '    auditRepository.save(record.value());',
    '}',
  ]

  const currentStep = steps[step]

  return (
    <div className="my-8 rounded-xl border-2 border-blue-500/30 bg-[var(--bg-surface)] overflow-hidden not-prose">
      <div className="px-5 py-4 border-b border-[var(--border)] bg-blue-50 dark:bg-blue-950/20">
        <div className="flex items-center gap-2">
          <Play size={16} className="text-blue-500" />
          <span className="font-heading font-semibold text-base">Code Walkthrough</span>
          <span className="font-mono text-[10px] text-[var(--text-tertiary)] ml-2">Step {step + 1} of {steps.length}</span>
        </div>
      </div>

      <div className="flex flex-col md:flex-row">
        {/* Code */}
        <div className="flex-1 p-4 border-b md:border-b-0 md:border-r border-[var(--border)]">
          {code.map((line, i) => (
            <div key={i}
              className={`flex font-mono text-[13px] leading-[1.8] px-3 rounded transition-all duration-300
                ${currentStep.highlight.includes(i) ? 'bg-[var(--accent-subtle)] border-l-2 border-l-[var(--accent)]' : 'border-l-2 border-l-transparent'}`}>
              <span className="w-6 shrink-0 text-[var(--text-tertiary)] text-right mr-3 select-none text-[11px]">{i + 1}</span>
              <pre style={{ margin: 0, background: 'transparent' }}><code>{line}</code></pre>
            </div>
          ))}
        </div>

        {/* Explanation */}
        <div className="w-full md:w-[280px] p-4">
          <div className="font-mono text-[10px] text-[var(--accent)] uppercase tracking-wider mb-2">
            {currentStep.title}
          </div>
          <p className="text-[13px] leading-relaxed text-[var(--text-secondary)]">
            {currentStep.desc}
          </p>
        </div>
      </div>

      {/* Step controls */}
      <div className="px-5 py-3 border-t border-[var(--border)] flex items-center justify-between bg-[var(--bg-surface-alt)]">
        <div className="flex gap-1">
          {steps.map((_, i) => (
            <button key={i} onClick={() => setStep(i)}
              className={`w-2 h-2 rounded-full transition-all ${step === i ? 'bg-[var(--accent)] w-6' : 'bg-[var(--border)] hover:bg-[var(--text-tertiary)]'}`} />
          ))}
        </div>
        <div className="flex items-center gap-2">
          <button onClick={() => setStep(Math.max(0, step - 1))} disabled={step === 0}
            className="px-3 py-1.5 rounded-lg border border-[var(--border)] font-mono text-[11px] hover:bg-[var(--bg-hover)] transition-colors disabled:opacity-30">
            Prev
          </button>
          <button onClick={() => setStep(Math.min(steps.length - 1, step + 1))} disabled={step === steps.length - 1}
            className="px-3 py-1.5 rounded-lg bg-[var(--accent)] text-white font-mono text-[11px] hover:opacity-90 transition-colors disabled:opacity-30 flex items-center gap-1">
            Next <SkipForward size={11} />
          </button>
        </div>
      </div>
    </div>
  )
}

function FlashcardDemo() {
  const [flipped, setFlipped] = useState(false)
  const [cardIdx, setCardIdx] = useState(0)

  const cards = [
    { q: 'What is a consumer group and why do we need it?', a: 'A consumer group is a set of consumers that cooperate to consume data from topics. Needed for: (1) Parallel processing, (2) Fault tolerance, (3) Load balancing.' },
    { q: 'What happens when a consumer dies mid-processing?', a: 'Heartbeat stops → group coordinator triggers rebalance → partitions reassigned to survivors. With auto-commit: possible reprocessing. With manual commit: only uncommitted messages reprocessed.' },
  ]

  const card = cards[cardIdx]

  return (
    <div className="my-8 not-prose">
      <div className="flex items-center justify-between mb-3">
        <span className="font-mono text-[10px] font-semibold uppercase tracking-[0.06em] text-[var(--accent)]">
          Flashcard Mode — Q&A
        </span>
        <span className="font-mono text-[11px] text-[var(--text-tertiary)]">
          Card {cardIdx + 1} of {cards.length}
        </span>
      </div>

      {/* Card */}
      <div
        onClick={() => setFlipped(!flipped)}
        className="relative cursor-pointer rounded-xl border-2 border-[var(--border)] bg-[var(--bg-surface)] p-8 min-h-[200px] flex items-center justify-center transition-all hover:shadow-md hover:border-[var(--accent)]"
        style={{ perspective: '1000px' }}
      >
        <div className={`text-center transition-all duration-300 ${flipped ? 'opacity-0 scale-95' : 'opacity-100 scale-100'}`}
          style={{ display: flipped ? 'none' : 'block' }}>
          <div className="font-mono text-[10px] text-[var(--text-tertiary)] uppercase tracking-wider mb-4">Question</div>
          <p className="font-heading font-semibold text-xl leading-relaxed">{card.q}</p>
          <p className="text-[12px] text-[var(--text-tertiary)] mt-4 font-mono">Click to reveal answer</p>
        </div>

        <div className={`text-center transition-all duration-300 ${flipped ? 'opacity-100 scale-100' : 'opacity-0 scale-95'}`}
          style={{ display: flipped ? 'block' : 'none' }}>
          <div className="font-mono text-[10px] text-emerald-500 uppercase tracking-wider mb-4">Answer</div>
          <p className="text-[15px] leading-relaxed text-[var(--text-secondary)]">{card.a}</p>
        </div>
      </div>

      {/* Controls */}
      <div className="flex items-center justify-between mt-4">
        <div className="flex gap-2">
          {['Again', 'Hard', 'Good', 'Easy'].map((label, i) => {
            const colors = ['text-red-500 border-red-200 hover:bg-red-50 dark:hover:bg-red-950/20', 'text-orange-500 border-orange-200 hover:bg-orange-50 dark:hover:bg-orange-950/20', 'text-emerald-500 border-emerald-200 hover:bg-emerald-50 dark:hover:bg-emerald-950/20', 'text-blue-500 border-blue-200 hover:bg-blue-50 dark:hover:bg-blue-950/20']
            return (
              <button key={label}
                onClick={() => { setFlipped(false); setCardIdx((cardIdx + 1) % cards.length) }}
                className={`px-3 py-1.5 rounded-lg border font-mono text-[11px] font-semibold transition-colors ${colors[i]}`}>
                {label}
              </button>
            )
          })}
        </div>
        <button onClick={() => { setFlipped(false); setCardIdx((cardIdx + 1) % cards.length) }}
          className="flex items-center gap-1 px-3 py-1.5 rounded-lg bg-[var(--accent)] text-white font-mono text-[11px] hover:opacity-90 transition-colors">
          Next <SkipForward size={11} />
        </button>
      </div>
    </div>
  )
}

export default function Approach4CustomReact() {
  return (
    <div>
      {/* Approach label */}
      <div className="mb-8 p-4 rounded-xl border border-[var(--border)] bg-[var(--bg-surface)]">
        <div className="flex items-center gap-3 mb-2">
          <span className="font-mono text-[10px] font-semibold uppercase tracking-[0.06em] px-2 py-1 rounded bg-amber-100 dark:bg-amber-900/30 text-amber-600 dark:text-amber-400">
            Approach 4
          </span>
          <span className="font-heading font-semibold">Full Custom React</span>
        </div>
        <p className="text-[13px] text-[var(--text-secondary)] leading-relaxed">
          Every element hand-crafted as React components. <strong>Interactive SVG diagrams</strong> you can manipulate,
          <strong> code walkthroughs</strong> with step-by-step explanations, <strong>flashcard Q&A</strong> with spaced repetition,
          <strong> simulators</strong> where you kill consumers and watch rebalancing animate. Maximum "wow" factor.
        </p>
        <div className="flex gap-4 mt-3 font-mono text-[10px] text-[var(--text-tertiary)]">
          <span>Effort: ~10 weeks per guide</span>
          <span>Interactivity: Maximum</span>
          <span>Content changes: Full rewrite to JSX</span>
        </div>
      </div>

      {/* Custom rendered content */}
      <div className="prose prose-gray dark:prose-invert font-body">
        <h2 className="font-heading">Consumer Groups and Rebalancing</h2>
        <h3 className="font-heading">Consumer Group Basics</h3>
        <p>
          A <strong>consumer group</strong> is a set of consumers that cooperate to consume messages from one or more topics.
          Kafka guarantees that each partition is consumed by exactly one consumer within a group.
        </p>
      </div>

      {/* Interactive consumer group simulator */}
      <KafkaConsumerGroupSimulator />

      {/* Code walkthrough */}
      <InteractiveCodeWalkthrough />

      {/* Callout */}
      <div className="my-6 rounded-r-lg border-l-4 border-l-[var(--accent)] bg-[var(--accent-subtle)] p-4 not-prose">
        <div className="flex items-center gap-2 mb-2">
          <Briefcase size={14} className="text-[var(--accent)]" />
          <span className="font-mono text-[10px] font-semibold uppercase tracking-[0.06em] text-[var(--accent)]">
            Real-World Application — Walmart Audit Pipeline
          </span>
        </div>
        <p className="text-[14px] leading-relaxed">
          Try setting the simulator above to <strong>6 partitions, 3 consumers, Round Robin</strong> — that's exactly
          how we configured it at Walmart. Then kill Consumer 1 to see how the partitions redistribute — this is what
          happens during a deployment rollout.
        </p>
      </div>

      {/* Flashcard Q&A */}
      <FlashcardDemo />
    </div>
  )
}

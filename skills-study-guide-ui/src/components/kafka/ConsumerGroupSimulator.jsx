import { useState, useMemo, useCallback, useEffect, useRef } from 'react'
import { useTheme } from '../../context/ThemeContext'
import { Zap, RotateCcw, AlertTriangle, Users, Layers, ArrowRight, Check, Loader2, Radio, Shuffle } from 'lucide-react'

const COLORS = ['#D97706', '#059669', '#2563EB', '#DC2626', '#7C3AED', '#0891B2']
const PARTITION_OPTIONS = [3, 6, 9, 12]
const CONSUMER_OPTIONS = [1, 2, 3, 4, 5, 6, 7, 8]

const STRATEGIES = [
  { key: 'roundrobin', label: 'Round Robin' },
  { key: 'range', label: 'Range' },
  { key: 'cooperative-sticky', label: 'Cooperative Sticky' },
]

const REBALANCE_STEPS = [
  { label: 'Consumers stop processing', icon: 'pause' },
  { label: 'JoinGroup requests sent to coordinator', icon: 'send' },
  { label: 'Leader elected, runs assignor', icon: 'leader' },
  { label: 'SyncGroup with new assignments', icon: 'sync' },
  { label: 'Consumers resume processing', icon: 'resume' },
]

const COOPERATIVE_REBALANCE_STEPS = [
  { label: 'Revoke only migrating partitions', icon: 'pause' },
  { label: 'JoinGroup requests sent to coordinator', icon: 'send' },
  { label: 'Leader runs sticky assignor', icon: 'leader' },
  { label: 'SyncGroup — only moved partitions reassigned', icon: 'sync' },
  { label: 'Consumers resume (most never stopped)', icon: 'resume' },
]

/**
 * Compute partition assignment for round-robin strategy.
 */
function assignRoundRobin(partitionCount, activeConsumers) {
  return Array.from({ length: partitionCount }, (_, i) => ({
    id: i,
    consumer: activeConsumers.length > 0 ? activeConsumers[i % activeConsumers.length] : null,
  }))
}

/**
 * Compute partition assignment for range strategy.
 */
function assignRange(partitionCount, activeConsumers) {
  if (activeConsumers.length === 0) {
    return Array.from({ length: partitionCount }, (_, i) => ({ id: i, consumer: null }))
  }
  const perConsumer = Math.ceil(partitionCount / activeConsumers.length)
  return Array.from({ length: partitionCount }, (_, i) => {
    const idx = Math.floor(i / perConsumer)
    return {
      id: i,
      consumer: activeConsumers[Math.min(idx, activeConsumers.length - 1)],
    }
  })
}

/**
 * Cooperative sticky: try to keep partitions with their previous owner,
 * only moving partitions that must move (e.g. from a dead consumer or to
 * even out load).
 */
function assignCooperativeSticky(partitionCount, activeConsumers, prevAssignment) {
  if (activeConsumers.length === 0) {
    return Array.from({ length: partitionCount }, (_, i) => ({ id: i, consumer: null }))
  }

  const activeSet = new Set(activeConsumers)
  const target = Math.ceil(partitionCount / activeConsumers.length)

  // Start from previous assignment — keep valid ones
  const result = Array.from({ length: partitionCount }, (_, i) => {
    const prev = prevAssignment?.find(p => p.id === i)
    const owner = prev && prev.consumer !== null && activeSet.has(prev.consumer)
      ? prev.consumer
      : null
    return { id: i, consumer: owner }
  })

  // Count current load per consumer
  const load = new Map(activeConsumers.map(c => [c, 0]))
  result.forEach(p => {
    if (p.consumer !== null) load.set(p.consumer, (load.get(p.consumer) || 0) + 1)
  })

  // Strip over-assigned partitions (beyond target) from consumers with too many
  for (const p of result) {
    if (p.consumer !== null && load.get(p.consumer) > target) {
      load.set(p.consumer, load.get(p.consumer) - 1)
      p.consumer = null
    }
  }

  // Assign unassigned partitions to least-loaded consumers
  const unassigned = result.filter(p => p.consumer === null)
  for (const p of unassigned) {
    let minC = activeConsumers[0]
    let minLoad = load.get(minC) ?? 0
    for (const c of activeConsumers) {
      const cl = load.get(c) ?? 0
      if (cl < minLoad) { minC = c; minLoad = cl }
    }
    p.consumer = minC
    load.set(minC, (load.get(minC) || 0) + 1)
  }

  return result
}

/**
 * Determine which partitions moved between two assignments.
 */
function diffAssignments(prev, next) {
  if (!prev || prev.length !== next.length) return new Set(next.map(p => p.id))
  const moved = new Set()
  for (let i = 0; i < next.length; i++) {
    if (prev[i].consumer !== next[i].consumer) moved.add(next[i].id)
  }
  return moved
}

// ----------- Rebalance Protocol Visualization -----------

function RebalanceProtocol({ active, step, isCooperative }) {
  const steps = isCooperative ? COOPERATIVE_REBALANCE_STEPS : REBALANCE_STEPS

  if (!active && step === -1) return null

  return (
    <div className="mt-5 rounded-lg border border-[var(--border)] bg-[var(--bg-surface-alt)] overflow-hidden">
      <div className="px-4 py-2.5 border-b border-[var(--border)] flex items-center gap-2">
        <Radio size={13} className="text-[var(--accent)]" />
        <span className="font-mono text-[10px] font-semibold uppercase tracking-wider text-[var(--accent)]">
          {isCooperative ? 'Cooperative Incremental Rebalance' : 'Eager Rebalance Protocol (JoinGroup)'}
        </span>
      </div>
      <div className="px-4 py-3 flex flex-col gap-0">
        {steps.map((s, i) => {
          const isDone = i < step
          const isCurrent = i === step
          const isPending = i > step

          return (
            <div key={i} className="flex items-start gap-3">
              {/* Vertical connector */}
              <div className="flex flex-col items-center">
                <div
                  className={`w-6 h-6 rounded-full flex items-center justify-center shrink-0 transition-all duration-300
                    ${isDone ? 'bg-emerald-500 text-white' : ''}
                    ${isCurrent ? 'bg-[var(--accent)] text-white ring-2 ring-[var(--accent)]/30 scale-110' : ''}
                    ${isPending ? 'bg-[var(--bg-surface)] border-2 border-[var(--border)] text-[var(--text-tertiary)]' : ''}`}
                >
                  {isDone ? <Check size={12} /> : isCurrent ? <Loader2 size={12} className="animate-spin" /> : <span className="text-[10px] font-mono">{i + 1}</span>}
                </div>
                {i < steps.length - 1 && (
                  <div className={`w-0.5 h-5 transition-colors duration-300 ${isDone ? 'bg-emerald-400' : 'bg-[var(--border)]'}`} />
                )}
              </div>
              <span className={`font-mono text-[12px] pt-0.5 transition-all duration-300
                ${isDone ? 'text-emerald-600 dark:text-emerald-400 line-through opacity-70' : ''}
                ${isCurrent ? 'text-[var(--text-primary)] font-semibold' : ''}
                ${isPending ? 'text-[var(--text-tertiary)]' : ''}`}>
                {s.label}
              </span>
            </div>
          )
        })}
      </div>
    </div>
  )
}

// ----------- Main Simulator -----------

export default function ConsumerGroupSimulator() {
  const { darkMode } = useTheme()

  const [consumers, setConsumers] = useState(3)
  const [partitions, setPartitions] = useState(6)
  const [strategy, setStrategy] = useState('roundrobin')
  const [failedConsumer, setFailedConsumer] = useState(null)
  const [isAnimating, setIsAnimating] = useState(false)
  const [movedPartitions, setMovedPartitions] = useState(new Set())

  // Rebalance protocol animation state
  const [rebalanceActive, setRebalanceActive] = useState(false)
  const [rebalanceStep, setRebalanceStep] = useState(-1)
  const rebalanceTimerRef = useRef(null)

  // Keep track of previous assignment for cooperative-sticky diffing
  const prevAssignmentRef = useRef(null)

  const activeConsumers = useMemo(() => {
    return Array.from({ length: consumers }, (_, i) => i).filter(i => i !== failedConsumer)
  }, [consumers, failedConsumer])

  const assignment = useMemo(() => {
    let result
    if (strategy === 'roundrobin') {
      result = assignRoundRobin(partitions, activeConsumers)
    } else if (strategy === 'range') {
      result = assignRange(partitions, activeConsumers)
    } else {
      result = assignCooperativeSticky(partitions, activeConsumers, prevAssignmentRef.current)
    }
    return result
  }, [partitions, activeConsumers, strategy])

  // Compute moved partitions whenever assignment changes
  useEffect(() => {
    if (prevAssignmentRef.current) {
      const moved = diffAssignments(prevAssignmentRef.current, assignment)
      setMovedPartitions(moved)
    }
    prevAssignmentRef.current = assignment
  }, [assignment])

  // Clear moved-partition highlighting after animation
  useEffect(() => {
    if (movedPartitions.size > 0) {
      const timer = setTimeout(() => setMovedPartitions(new Set()), 1200)
      return () => clearTimeout(timer)
    }
  }, [movedPartitions])

  /**
   * Run the animated rebalance protocol steps.
   */
  const runRebalanceProtocol = useCallback(() => {
    // Clear any existing timer
    if (rebalanceTimerRef.current) clearTimeout(rebalanceTimerRef.current)

    setRebalanceActive(true)
    setRebalanceStep(0)

    const stepCount = 5
    let current = 0

    function advanceStep() {
      current++
      if (current < stepCount) {
        setRebalanceStep(current)
        rebalanceTimerRef.current = setTimeout(advanceStep, 600)
      } else {
        // Finish: keep last step briefly, then mark complete
        rebalanceTimerRef.current = setTimeout(() => {
          setRebalanceStep(stepCount)
          rebalanceTimerRef.current = setTimeout(() => {
            setRebalanceActive(false)
            setRebalanceStep(-1)
          }, 800)
        }, 500)
      }
    }

    rebalanceTimerRef.current = setTimeout(advanceStep, 600)
  }, [])

  // Cleanup timers on unmount
  useEffect(() => {
    return () => {
      if (rebalanceTimerRef.current) clearTimeout(rebalanceTimerRef.current)
    }
  }, [])

  const killConsumer = useCallback((id) => {
    setIsAnimating(true)
    setFailedConsumer(id)
    runRebalanceProtocol()
    setTimeout(() => setIsAnimating(false), 600)
  }, [runRebalanceProtocol])

  const handleConsumerChange = useCallback((n) => {
    setConsumers(n)
    setFailedConsumer(null)
    if (n !== consumers) runRebalanceProtocol()
  }, [consumers, runRebalanceProtocol])

  const handlePartitionChange = useCallback((n) => {
    setPartitions(n)
    setFailedConsumer(null)
    if (n !== partitions) runRebalanceProtocol()
  }, [partitions, runRebalanceProtocol])

  const handleStrategyChange = useCallback((s) => {
    // Clear previous ref when switching strategies so cooperative-sticky starts fresh
    if (s !== strategy) prevAssignmentRef.current = null
    setStrategy(s)
  }, [strategy])

  const handleReset = useCallback(() => {
    if (rebalanceTimerRef.current) clearTimeout(rebalanceTimerRef.current)
    prevAssignmentRef.current = null
    setConsumers(3)
    setPartitions(6)
    setFailedConsumer(null)
    setStrategy('roundrobin')
    setIsAnimating(false)
    setMovedPartitions(new Set())
    setRebalanceActive(false)
    setRebalanceStep(-1)
  }, [])

  const activeCount = consumers - (failedConsumer !== null ? 1 : 0)
  const idleCount = Math.max(0, activeCount - partitions)
  const unassignedCount = assignment.filter(p => p.consumer === null).length
  const isCooperative = strategy === 'cooperative-sticky'

  // Responsive grid columns for partitions
  const partitionCols = partitions <= 3 ? 'grid-cols-3'
    : partitions <= 6 ? 'grid-cols-3 sm:grid-cols-6'
    : partitions <= 9 ? 'grid-cols-3 sm:grid-cols-6 lg:grid-cols-9'
    : 'grid-cols-3 sm:grid-cols-6 lg:grid-cols-12'

  return (
    <div className="my-8 rounded-xl border-2 border-[var(--accent)]/30 bg-[var(--bg-surface)] overflow-hidden not-prose">
      {/* Header */}
      <div className="px-5 py-4 border-b border-[var(--border)] bg-[var(--accent-subtle)]">
        <div className="flex items-center justify-between flex-wrap gap-2">
          <div>
            <div className="flex items-center gap-2">
              <Zap size={16} className="text-[var(--accent)]" />
              <span className="font-heading font-semibold text-base">Kafka Consumer Group Simulator</span>
            </div>
            <p className="text-[12px] text-[var(--text-secondary)] mt-1">
              Configure partitions & consumers, kill consumers, change strategy — watch partitions redistribute live
            </p>
          </div>
          <button
            onClick={handleReset}
            className="flex items-center gap-1 px-3 py-1.5 rounded-lg border border-[var(--border)] text-[11px] font-mono hover:bg-[var(--bg-hover)] transition-colors"
          >
            <RotateCcw size={11} /> Reset
          </button>
        </div>
      </div>

      <div className="p-5">
        {/* Controls */}
        <div className="flex flex-col gap-4 mb-6">
          {/* Row 1: Partitions + Consumers */}
          <div className="flex flex-wrap items-center gap-6">
            {/* Partition count */}
            <div className="flex items-center gap-3">
              <span className="font-mono text-[11px] text-[var(--text-secondary)] flex items-center gap-1">
                <Layers size={12} /> Partitions:
              </span>
              <div className="flex items-center gap-1">
                {PARTITION_OPTIONS.map(n => (
                  <button
                    key={n}
                    onClick={() => handlePartitionChange(n)}
                    className={`w-8 h-7 rounded-lg font-mono text-[12px] font-semibold transition-all
                      ${partitions === n
                        ? 'bg-[var(--accent)] text-white shadow-sm'
                        : 'border border-[var(--border)] hover:bg-[var(--bg-hover)]'}`}
                  >
                    {n}
                  </button>
                ))}
              </div>
            </div>

            {/* Consumer count */}
            <div className="flex items-center gap-3">
              <span className="font-mono text-[11px] text-[var(--text-secondary)] flex items-center gap-1">
                <Users size={12} /> Consumers:
              </span>
              <div className="flex items-center gap-1">
                {CONSUMER_OPTIONS.map(n => (
                  <button
                    key={n}
                    onClick={() => handleConsumerChange(n)}
                    className={`w-7 h-7 rounded-lg font-mono text-[12px] font-semibold transition-all
                      ${consumers === n
                        ? 'bg-[var(--accent)] text-white shadow-sm'
                        : 'border border-[var(--border)] hover:bg-[var(--bg-hover)]'}`}
                  >
                    {n}
                  </button>
                ))}
              </div>
            </div>
          </div>

          {/* Row 2: Strategy */}
          <div className="flex items-center gap-2">
            <span className="font-mono text-[11px] text-[var(--text-secondary)] flex items-center gap-1">
              <Shuffle size={12} /> Strategy:
            </span>
            {STRATEGIES.map(s => (
              <button
                key={s.key}
                onClick={() => handleStrategyChange(s.key)}
                className={`px-3 py-1.5 rounded-lg font-mono text-[11px] transition-all
                  ${strategy === s.key
                    ? 'bg-[var(--accent)] text-white'
                    : 'border border-[var(--border)] hover:bg-[var(--bg-hover)]'}`}
              >
                {s.label}
              </button>
            ))}
          </div>
        </div>

        {/* Topic label */}
        <div className="flex items-center gap-2 mb-3">
          <div className="h-px flex-1 bg-[var(--border)]" />
          <span className="font-mono text-[10px] text-[var(--text-tertiary)] uppercase tracking-wider">
            Topic: audit-events ({partitions} partitions)
          </span>
          <div className="h-px flex-1 bg-[var(--border)]" />
        </div>

        {/* Partitions visualization */}
        <div className={`grid ${partitionCols} gap-3 mb-6`}>
          {assignment.map(p => {
            const wasMoved = movedPartitions.has(p.id)
            const color = p.consumer !== null ? COLORS[p.consumer % COLORS.length] : null
            const stayed = isCooperative && !wasMoved && p.consumer !== null

            return (
              <div
                key={p.id}
                className={`rounded-xl border-2 p-3 text-center transition-all duration-500
                  ${isAnimating ? 'scale-95' : 'scale-100'}
                  ${wasMoved && isCooperative ? 'animate-pulse ring-2 ring-amber-400/50' : ''}
                  ${stayed ? 'ring-1 ring-emerald-400/40' : ''}`}
                style={{
                  borderColor: color || 'var(--border)',
                  backgroundColor: color ? color + '12' : 'transparent',
                  boxShadow: color ? `0 2px 8px ${color}20` : 'none',
                }}
              >
                <div className="font-mono text-[10px] text-[var(--text-tertiary)] mb-0.5">Partition</div>
                <div
                  className="font-mono text-xl font-bold"
                  style={{ color: color || 'var(--text-tertiary)' }}
                >
                  {p.id}
                </div>
                <div
                  className="font-mono text-[10px] mt-0.5"
                  style={{ color: color || 'var(--text-tertiary)' }}
                >
                  {p.consumer !== null ? `C${p.consumer}` : 'Unassigned'}
                </div>
                {/* Cooperative sticky move indicator */}
                {isCooperative && wasMoved && p.consumer !== null && (
                  <div className="mt-1 font-mono text-[9px] text-amber-600 dark:text-amber-400 font-semibold">
                    MOVED
                  </div>
                )}
                {isCooperative && stayed && (
                  <div className="mt-1 font-mono text-[9px] text-emerald-600 dark:text-emerald-400">
                    stayed
                  </div>
                )}
              </div>
            )
          })}
        </div>

        {/* Arrow bridge */}
        <div className="flex justify-center mb-3">
          <ArrowRight size={16} className="text-[var(--text-tertiary)] rotate-90" />
        </div>

        {/* Consumer group label */}
        <div className="flex items-center gap-2 mb-3">
          <div className="h-px flex-1 bg-[var(--border)]" />
          <span className="font-mono text-[10px] text-[var(--text-tertiary)] uppercase tracking-wider">
            Consumer Group: audit-processor
          </span>
          <div className="h-px flex-1 bg-[var(--border)]" />
        </div>

        {/* Consumers */}
        <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-3">
          {Array.from({ length: consumers }, (_, i) => {
            const isFailed = failedConsumer === i
            const myPartitions = assignment.filter(p => p.consumer === i)
            const isIdle = myPartitions.length === 0
            const color = COLORS[i % COLORS.length]

            return (
              <div
                key={i}
                className={`rounded-xl border-2 p-3 transition-all duration-300
                  ${isFailed ? 'opacity-30 scale-95 border-red-400 bg-red-50 dark:bg-red-950/20' : ''}
                  ${isIdle && !isFailed ? 'border-dashed border-[var(--border)]' : ''}`}
                style={{
                  borderColor: isFailed ? '#DC2626' : isIdle ? undefined : color,
                  backgroundColor: isFailed ? undefined : isIdle ? 'transparent' : color + '08',
                }}
              >
                <div className="flex items-center justify-between mb-2">
                  <div className="flex items-center gap-2">
                    <div
                      className="w-3 h-3 rounded-full shrink-0"
                      style={{ backgroundColor: isFailed ? '#DC2626' : color }}
                    />
                    <span className="font-mono text-[12px] font-semibold">Consumer {i}</span>
                  </div>
                  {!isFailed && (
                    <button
                      onClick={() => killConsumer(i)}
                      title="Simulate consumer failure"
                      className="p-1 rounded hover:bg-red-100 dark:hover:bg-red-900/30 transition-colors"
                    >
                      <AlertTriangle size={12} className="text-red-400" />
                    </button>
                  )}
                </div>
                <div className="font-mono text-[10px] text-[var(--text-tertiary)]">
                  {isFailed
                    ? '\u2620 DEAD \u2014 partitions reassigned'
                    : isIdle
                      ? 'Idle (no partitions)'
                      : `Partitions: ${myPartitions.map(p => `P${p.id}`).join(', ')}`}
                </div>
              </div>
            )
          })}
        </div>

        {/* Rebalance protocol animation */}
        <RebalanceProtocol
          active={rebalanceActive}
          step={rebalanceStep}
          isCooperative={isCooperative}
        />

        {/* Status bar */}
        <div className="mt-5 pt-4 border-t border-[var(--border)] flex flex-wrap gap-4 font-mono text-[11px] text-[var(--text-tertiary)]">
          <span>Active: {activeCount}/{consumers}</span>
          <span>Idle: {idleCount}</span>
          <span>Unassigned: {unassignedCount}</span>
          <span>Strategy: {STRATEGIES.find(s => s.key === strategy)?.label}</span>
          <span>Partitions: {partitions}</span>
          {isCooperative && movedPartitions.size > 0 && (
            <span className="text-amber-600 dark:text-amber-400 font-semibold">
              Moved: {movedPartitions.size} / {partitions} partitions
            </span>
          )}
        </div>
      </div>
    </div>
  )
}

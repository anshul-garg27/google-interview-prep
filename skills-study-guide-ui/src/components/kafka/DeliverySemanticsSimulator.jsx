import { useState, useCallback, useEffect, useRef } from 'react'
import { useTheme } from '../../context/ThemeContext'
import {
  Play,
  RotateCcw,
  Zap,
  AlertTriangle,
  CheckCircle2,
  XCircle,
  Copy,
  ArrowDown,
  Radio,
  Server,
  Database,
  Users,
  Settings,
} from 'lucide-react'

const STEPS = [
  { id: 1, label: 'Producer sends message', icon: Radio, zone: 'producer' },
  { id: 2, label: 'Leader broker receives', icon: Server, zone: 'broker' },
  { id: 3, label: 'Leader commits to log', icon: Database, zone: 'broker' },
  { id: 4, label: 'Followers replicate', icon: Copy, zone: 'broker' },
  { id: 5, label: 'Ack sent to producer', icon: CheckCircle2, zone: 'broker' },
  { id: 6, label: 'Consumer fetches', icon: ArrowDown, zone: 'consumer' },
  { id: 7, label: 'Consumer processes', icon: Settings, zone: 'consumer' },
  { id: 8, label: 'Consumer commits offset', icon: Database, zone: 'consumer' },
]

const ACKS_OPTIONS = [
  { value: 0, label: 'acks=0', desc: 'Fire and forget' },
  { value: 1, label: 'acks=1', desc: 'Leader only' },
  { value: 'all', label: 'acks=all', desc: 'All ISR' },
]

const DELIVERY_MODES = [
  {
    key: 'at-most-once',
    label: 'At-Most-Once',
    color: '#EF4444',
    desc: 'Messages may be lost, never duplicated',
  },
  {
    key: 'at-least-once',
    label: 'At-Least-Once',
    color: '#F59E0B',
    desc: 'Messages never lost, may be duplicated',
  },
  {
    key: 'exactly-once',
    label: 'Exactly-Once',
    color: '#10B981',
    desc: 'Messages never lost, never duplicated',
  },
]

function getStepOutcome(mode, failStep, stepId, acks) {
  // If step hasn't been reached yet by the animation, it's pending
  // If step is before failure, it succeeded
  if (stepId < failStep) return 'success'
  if (stepId > failStep) {
    // Steps after failure depend on mode
    if (mode === 'at-most-once') return 'skipped'
    if (mode === 'at-least-once') {
      // Retries from beginning — steps re-run, some may duplicate
      if (stepId <= 5) return 'retry'
      if (stepId === 6 || stepId === 7) return 'duplicate'
      if (stepId === 8) return 'success'
    }
    if (mode === 'exactly-once') {
      if (stepId <= 5) return 'retry'
      return 'success'
    }
  }
  // stepId === failStep
  return 'failure'
}

function computeStats(mode, failStep) {
  if (!failStep) return { sent: 3, received: 3, lost: 0, duplicated: 0 }

  if (mode === 'at-most-once') {
    // Messages lost after failure point
    if (failStep <= 3) return { sent: 3, received: 2, lost: 1, duplicated: 0 }
    if (failStep <= 5) return { sent: 3, received: 2, lost: 1, duplicated: 0 }
    if (failStep <= 7) return { sent: 3, received: 2, lost: 1, duplicated: 0 }
    return { sent: 3, received: 3, lost: 0, duplicated: 0 }
  }
  if (mode === 'at-least-once') {
    // Retry causes duplication
    if (failStep <= 5) return { sent: 3, received: 3, lost: 0, duplicated: 1 }
    if (failStep <= 7) return { sent: 3, received: 4, lost: 0, duplicated: 1 }
    return { sent: 3, received: 3, lost: 0, duplicated: 0 }
  }
  // exactly-once — idempotent
  return { sent: 3, received: 3, lost: 0, duplicated: 0 }
}

function getFailureExplanation(mode, failStep) {
  const stepLabel = STEPS.find((s) => s.id === failStep)?.label || ''
  if (mode === 'at-most-once') {
    if (failStep <= 2)
      return `Failure at "${stepLabel}": Producer does not retry. Message is lost permanently.`
    if (failStep <= 5)
      return `Failure at "${stepLabel}": No retry, offset already advanced. Message lost.`
    if (failStep <= 7)
      return `Failure at "${stepLabel}": Consumer committed offset before processing. Message skipped.`
    return `Failure at "${stepLabel}": Message already processed. No impact.`
  }
  if (mode === 'at-least-once') {
    if (failStep <= 5)
      return `Failure at "${stepLabel}": Producer retries send. Broker may receive the message again, causing a duplicate.`
    if (failStep <= 7)
      return `Failure at "${stepLabel}": Consumer retries from last committed offset. Message processed again — duplicate delivery.`
    return `Failure at "${stepLabel}": Processing complete before offset commit failed. On rebalance, message re-delivered.`
  }
  // exactly-once
  if (failStep <= 5)
    return `Failure at "${stepLabel}": Idempotent producer retries with same sequence number. Broker deduplicates — no duplicate.`
  if (failStep <= 7)
    return `Failure at "${stepLabel}": Transactional consumer reads only committed offsets. On retry, exactly-once semantics preserved.`
  return `Failure at "${stepLabel}": Atomic offset commit within transaction. Failure triggers rollback, then clean retry.`
}

function StepRow({ step, outcome, isActive, onInjectFailure, canInject }) {
  const Icon = step.icon
  const bgMap = {
    success: 'bg-emerald-500/15 border-emerald-500/40',
    failure: 'bg-red-500/15 border-red-500/40',
    retry: 'bg-amber-500/15 border-amber-500/40',
    duplicate: 'bg-amber-500/15 border-amber-500/40',
    skipped: 'bg-[var(--bg-surface-alt)] border-[var(--border)] opacity-50',
    pending: 'bg-[var(--bg-surface)] border-[var(--border)]',
  }
  const iconColorMap = {
    success: 'text-emerald-500',
    failure: 'text-red-500',
    retry: 'text-amber-500',
    duplicate: 'text-amber-500',
    skipped: 'text-[var(--text-tertiary)]',
    pending: 'text-[var(--text-secondary)]',
  }
  const tagMap = {
    success: null,
    failure: { label: 'FAIL', cls: 'bg-red-500 text-white' },
    retry: { label: 'RETRY', cls: 'bg-amber-500 text-white' },
    duplicate: { label: 'DUP', cls: 'bg-amber-500 text-white' },
    skipped: { label: 'LOST', cls: 'bg-red-500/80 text-white' },
    pending: null,
  }

  const tag = tagMap[outcome]

  return (
    <div
      className={`group relative flex items-center gap-2 rounded-lg border px-3 py-2 transition-all duration-300 ${bgMap[outcome]} ${isActive ? 'ring-2 ring-[var(--accent)] scale-[1.02]' : ''}`}
    >
      <Icon className={`w-4 h-4 flex-shrink-0 ${iconColorMap[outcome]}`} />
      <span className="font-mono text-[11px] text-[var(--text-secondary)] flex-shrink-0 w-4">
        {step.id}
      </span>
      <span className="text-xs font-body text-[var(--text-primary)] flex-1 leading-tight">
        {step.label}
      </span>
      {tag && (
        <span
          className={`text-[9px] font-mono font-bold px-1.5 py-0.5 rounded ${tag.cls}`}
        >
          {tag.label}
        </span>
      )}
      {canInject && (
        <button
          onClick={() => onInjectFailure(step.id)}
          className="absolute -right-1 -top-1 opacity-0 group-hover:opacity-100 transition-opacity bg-red-500 hover:bg-red-600 text-white rounded-full w-5 h-5 flex items-center justify-center shadow-md"
          title={`Inject failure at step ${step.id}`}
        >
          <Zap className="w-3 h-3" />
        </button>
      )}
    </div>
  )
}

function ModeColumn({ mode, acks, failStep, activeStep, onInjectFailure, isAnimating }) {
  const stats = computeStats(mode.key, failStep)
  const explanation = failStep ? getFailureExplanation(mode.key, failStep) : null

  return (
    <div className="flex flex-col flex-1 min-w-[220px]">
      {/* Header */}
      <div className="text-center mb-3">
        <h4
          className="font-heading font-semibold text-sm"
          style={{ color: mode.color }}
        >
          {mode.label}
        </h4>
        <p className="text-[10px] text-[var(--text-tertiary)] font-body mt-0.5">
          {mode.desc}
        </p>
      </div>

      {/* Steps */}
      <div className="flex flex-col gap-1.5">
        {STEPS.map((step) => {
          // Skip follower replication step for acks=0 and acks=1
          if (step.id === 4 && acks !== 'all') {
            return (
              <div
                key={step.id}
                className="flex items-center gap-2 rounded-lg border border-dashed border-[var(--border)] px-3 py-2 opacity-40"
              >
                <step.icon className="w-4 h-4 text-[var(--text-tertiary)]" />
                <span className="font-mono text-[11px] text-[var(--text-tertiary)] w-4">
                  {step.id}
                </span>
                <span className="text-xs text-[var(--text-tertiary)] italic font-body flex-1">
                  Skipped (acks≠all)
                </span>
              </div>
            )
          }

          let outcome = 'pending'
          if (failStep) {
            outcome = getStepOutcome(mode.key, failStep, step.id, acks)
          } else if (activeStep >= step.id) {
            outcome = 'success'
          }

          return (
            <StepRow
              key={step.id}
              step={step}
              outcome={outcome}
              isActive={!failStep && activeStep === step.id}
              onInjectFailure={onInjectFailure}
              canInject={!isAnimating && !failStep}
            />
          )
        })}
      </div>

      {/* Stats bar */}
      <div className="mt-3 grid grid-cols-4 gap-1 text-center">
        <StatBadge label="Sent" value={stats.sent} color="text-[var(--text-primary)]" />
        <StatBadge
          label="Recv"
          value={stats.received}
          color={stats.received < stats.sent ? 'text-red-500' : stats.received > stats.sent ? 'text-amber-500' : 'text-emerald-500'}
        />
        <StatBadge
          label="Lost"
          value={stats.lost}
          color={stats.lost > 0 ? 'text-red-500' : 'text-emerald-500'}
        />
        <StatBadge
          label="Dups"
          value={stats.duplicated}
          color={stats.duplicated > 0 ? 'text-amber-500' : 'text-emerald-500'}
        />
      </div>

      {/* Explanation */}
      {explanation && (
        <div
          className="mt-2 rounded-lg border px-3 py-2 text-[11px] font-body leading-relaxed animate-fade-in"
          style={{
            borderColor: mode.color + '40',
            backgroundColor: mode.color + '10',
            color: 'var(--text-primary)',
          }}
        >
          {explanation}
        </div>
      )}
    </div>
  )
}

function StatBadge({ label, value, color }) {
  return (
    <div className="rounded bg-[var(--bg-surface-alt)] px-1 py-1">
      <div className={`font-mono text-sm font-bold ${color}`}>{value}</div>
      <div className="text-[9px] text-[var(--text-tertiary)] uppercase tracking-wider">
        {label}
      </div>
    </div>
  )
}

export default function DeliverySemanticsSimulator() {
  const { darkMode } = useTheme()
  const [acks, setAcks] = useState(1)
  const [failStep, setFailStep] = useState(null)
  const [activeStep, setActiveStep] = useState(0)
  const [isAnimating, setIsAnimating] = useState(false)
  const timerRef = useRef(null)

  const clearAnimation = useCallback(() => {
    if (timerRef.current) {
      clearInterval(timerRef.current)
      timerRef.current = null
    }
    setIsAnimating(false)
  }, [])

  const runAnimation = useCallback(() => {
    clearAnimation()
    setFailStep(null)
    setActiveStep(0)
    setIsAnimating(true)

    let step = 0
    timerRef.current = setInterval(() => {
      step += 1
      if (step > 8) {
        clearInterval(timerRef.current)
        timerRef.current = null
        setIsAnimating(false)
        return
      }
      // Skip step 4 if acks != all
      if (step === 4 && acks !== 'all') {
        step = 5
      }
      setActiveStep(step)
    }, 500)
  }, [acks, clearAnimation])

  const handleInjectFailure = useCallback(
    (stepId) => {
      clearAnimation()
      setFailStep(stepId)
      setActiveStep(0)
    },
    [clearAnimation]
  )

  const handleReset = useCallback(() => {
    clearAnimation()
    setFailStep(null)
    setActiveStep(0)
  }, [clearAnimation])

  useEffect(() => {
    return () => {
      if (timerRef.current) clearInterval(timerRef.current)
    }
  }, [])

  return (
    <div className="rounded-xl border border-[var(--border)] bg-[var(--bg-surface)] p-4 sm:p-6 my-6">
      {/* Title */}
      <div className="flex items-center gap-2 mb-4">
        <AlertTriangle className="w-5 h-5 text-[var(--accent)]" />
        <h3 className="font-heading text-lg font-semibold text-[var(--text-primary)]">
          Delivery Semantics Simulator
        </h3>
      </div>

      {/* Controls */}
      <div className="flex flex-wrap items-center gap-3 mb-5">
        {/* Acks toggle */}
        <div className="flex items-center rounded-lg border border-[var(--border)] overflow-hidden">
          {ACKS_OPTIONS.map((opt) => (
            <button
              key={opt.value}
              onClick={() => {
                setAcks(opt.value)
                handleReset()
              }}
              className={`px-3 py-1.5 text-xs font-mono transition-colors ${
                acks === opt.value
                  ? 'bg-[var(--accent)] text-white'
                  : 'bg-[var(--bg-surface)] text-[var(--text-secondary)] hover:bg-[var(--bg-hover)]'
              }`}
              title={opt.desc}
            >
              {opt.label}
            </button>
          ))}
        </div>

        {/* Run / Reset */}
        <button
          onClick={runAnimation}
          disabled={isAnimating}
          className="flex items-center gap-1.5 rounded-lg bg-emerald-500 hover:bg-emerald-600 disabled:opacity-50 text-white px-3 py-1.5 text-xs font-body transition-colors"
        >
          <Play className="w-3.5 h-3.5" />
          Run Flow
        </button>
        <button
          onClick={handleReset}
          className="flex items-center gap-1.5 rounded-lg border border-[var(--border)] bg-[var(--bg-surface)] hover:bg-[var(--bg-hover)] text-[var(--text-secondary)] px-3 py-1.5 text-xs font-body transition-colors"
        >
          <RotateCcw className="w-3.5 h-3.5" />
          Reset
        </button>

        {!failStep && !isAnimating && (
          <span className="text-[11px] text-[var(--text-tertiary)] font-body italic">
            Hover a step and click <Zap className="w-3 h-3 inline text-red-500" /> to inject a failure
          </span>
        )}
      </div>

      {/* Three columns */}
      <div className="flex flex-col lg:flex-row gap-4">
        {DELIVERY_MODES.map((mode) => (
          <ModeColumn
            key={mode.key}
            mode={mode}
            acks={acks}
            failStep={failStep}
            activeStep={activeStep}
            onInjectFailure={handleInjectFailure}
            isAnimating={isAnimating}
          />
        ))}
      </div>

      {/* Legend */}
      <div className="mt-4 flex flex-wrap gap-3 text-[10px] font-mono text-[var(--text-tertiary)] border-t border-[var(--border)] pt-3">
        <span className="flex items-center gap-1">
          <span className="w-2.5 h-2.5 rounded-sm bg-emerald-500/30 border border-emerald-500/50" />
          Success
        </span>
        <span className="flex items-center gap-1">
          <span className="w-2.5 h-2.5 rounded-sm bg-red-500/30 border border-red-500/50" />
          Failure
        </span>
        <span className="flex items-center gap-1">
          <span className="w-2.5 h-2.5 rounded-sm bg-amber-500/30 border border-amber-500/50" />
          Retry / Duplicate
        </span>
        <span className="flex items-center gap-1">
          <span className="w-2.5 h-2.5 rounded-sm bg-[var(--bg-surface-alt)] border border-[var(--border)]" />
          Skipped / Lost
        </span>
      </div>
    </div>
  )
}

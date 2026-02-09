import { useState, useCallback, useEffect, useRef } from 'react'
import { useTheme } from '../../context/ThemeContext'
import {
  ChevronLeft,
  ChevronRight,
  Play,
  Pause,
  RotateCcw,
  Radio,
  Hash,
  Server,
  Database,
  Copy,
  CheckCircle2,
  ArrowLeftRight,
  Users,
  Download,
  Package,
  Settings,
  Save,
} from 'lucide-react'

const FLOW_STEPS = [
  {
    id: 1,
    label: 'Producer sends record',
    icon: Radio,
    highlight: 'producer',
    detail: 'Producer calls send(topic, key, value, headers). The record is serialized and added to the internal batch buffer.',
    code: 'producer.send({ topic: "orders", key: "user-42", value: JSON.stringify(order), headers: { traceId: "abc" } })',
  },
  {
    id: 2,
    label: 'Partitioner selects partition',
    icon: Hash,
    highlight: 'partitioner',
    detail: 'The partitioner calculates the target partition using the key hash. With 6 partitions: murmur2("user-42") % 6 = partition 3.',
    code: 'partition = murmur2(key) % numPartitions  // murmur2("user-42") % 6 = 3',
  },
  {
    id: 3,
    label: 'ProduceRequest to leader',
    icon: ArrowLeftRight,
    highlight: 'arrow-to-leader',
    detail: 'The batched ProduceRequest is sent over TCP to the leader broker for partition 3. The request includes the topic, partition, and record batch.',
    code: 'ProduceRequest { topic: "orders", partition: 3, records: [RecordBatch] }',
  },
  {
    id: 4,
    label: 'Leader appends to log',
    icon: Database,
    highlight: 'leader',
    detail: 'The leader broker appends the record batch to the local log segment. Each record gets a monotonically increasing offset assigned.',
    code: 'orders-3/00000000000000042.log  →  offset 42 assigned',
  },
  {
    id: 5,
    label: 'Followers replicate',
    icon: Copy,
    highlight: 'followers',
    detail: 'Follower brokers send FetchRequests to the leader to replicate the new records. They append to their own local logs.',
    code: 'Follower-1 FetchRequest(offset=42) → Leader\nFollower-2 FetchRequest(offset=42) → Leader',
  },
  {
    id: 6,
    label: 'ISR fully caught up',
    icon: CheckCircle2,
    highlight: 'followers',
    detail: 'All In-Sync Replicas (ISR) have replicated up to offset 42. The high watermark advances. If acks=all, this triggers the ack.',
    code: 'ISR = [broker-0 (leader), broker-1, broker-2]\nHigh Watermark: 42  ✓',
  },
  {
    id: 7,
    label: 'ProduceResponse to producer',
    icon: ArrowLeftRight,
    highlight: 'arrow-to-producer',
    detail: 'The leader sends a ProduceResponse back to the producer with the assigned offset, confirming the write was successful.',
    code: 'ProduceResponse { topic: "orders", partition: 3, offset: 42, timestamp: 1707350400 }',
  },
  {
    id: 8,
    label: 'Consumer joins group',
    icon: Users,
    highlight: 'consumer',
    detail: 'Consumer sends JoinGroup + Heartbeat requests to the Group Coordinator. It receives its partition assignment (e.g., partitions 3, 4).',
    code: 'JoinGroupRequest  { groupId: "order-processors", memberId: "c-1" }\nAssignment: [orders-3, orders-4]',
  },
  {
    id: 9,
    label: 'Consumer sends FetchRequest',
    icon: Download,
    highlight: 'arrow-to-leader',
    detail: 'The consumer sends a FetchRequest to the leader of partition 3, starting from its last committed offset.',
    code: 'FetchRequest { topic: "orders", partition: 3, fetchOffset: 42, maxBytes: 1048576 }',
  },
  {
    id: 10,
    label: 'FetchResponse with batch',
    icon: Package,
    highlight: 'arrow-to-consumer',
    detail: 'The leader responds with a batch of records starting at offset 42. Records are delivered in order within the partition.',
    code: 'FetchResponse { records: [{ offset: 42, key: "user-42", value: {...} }], highWatermark: 42 }',
  },
  {
    id: 11,
    label: 'Consumer processes messages',
    icon: Settings,
    highlight: 'consumer',
    detail: 'The consumer deserializes and processes each record in the batch. Business logic is executed (e.g., update order status).',
    code: 'for (const record of batch) {\n  await processOrder(record.value)  // business logic\n}',
  },
  {
    id: 12,
    label: 'Consumer commits offset',
    icon: Save,
    highlight: 'consumer',
    detail: 'After successful processing, the consumer commits the offset to the __consumer_offsets topic, marking progress.',
    code: 'OffsetCommitRequest { group: "order-processors", topic: "orders", partition: 3, offset: 43 }',
  },
]

function DiagramNode({ label, active, x, y, width, height, icon: Icon, variant }) {
  const baseClasses =
    'absolute rounded-lg border-2 flex flex-col items-center justify-center transition-all duration-500 select-none'

  const activeClasses = active
    ? 'border-[var(--accent)] bg-[var(--accent)]/10 shadow-lg shadow-[var(--accent)]/10 scale-105'
    : 'border-[var(--border)] bg-[var(--bg-surface)]'

  return (
    <div
      className={`${baseClasses} ${activeClasses}`}
      style={{ left: x, top: y, width, height }}
    >
      {Icon && (
        <Icon
          className={`w-5 h-5 mb-1 transition-colors duration-300 ${
            active ? 'text-[var(--accent)]' : 'text-[var(--text-tertiary)]'
          }`}
        />
      )}
      <span
        className={`text-[10px] font-mono font-medium text-center leading-tight px-1 transition-colors duration-300 ${
          active ? 'text-[var(--accent)]' : 'text-[var(--text-secondary)]'
        }`}
      >
        {label}
      </span>
    </div>
  )
}

function AnimatedArrow({ active, x1, y1, x2, y2, direction }) {
  return (
    <svg
      className="absolute top-0 left-0 w-full h-full pointer-events-none"
      style={{ zIndex: 5 }}
    >
      <defs>
        <marker
          id={`arrowhead-${direction}`}
          markerWidth="8"
          markerHeight="6"
          refX="7"
          refY="3"
          orient="auto"
        >
          <polygon
            points="0 0, 8 3, 0 6"
            fill={active ? 'var(--accent)' : 'var(--border-strong)'}
            className="transition-all duration-300"
          />
        </marker>
      </defs>
      <line
        x1={x1}
        y1={y1}
        x2={x2}
        y2={y2}
        stroke={active ? 'var(--accent)' : 'var(--border-strong)'}
        strokeWidth={active ? 2.5 : 1.5}
        strokeDasharray={active ? 'none' : '4 3'}
        markerEnd={`url(#arrowhead-${direction})`}
        className="transition-all duration-500"
      />
      {active && (
        <circle r="3" fill="var(--accent)">
          <animateMotion
            dur="1s"
            repeatCount="indefinite"
            path={`M${x1},${y1} L${x2},${y2}`}
          />
        </circle>
      )}
    </svg>
  )
}

function FlowDiagram({ currentStep }) {
  const step = FLOW_STEPS[currentStep - 1]
  const hl = step?.highlight || ''

  return (
    <div className="relative w-full h-[260px] sm:h-[220px] overflow-hidden">
      {/* Nodes */}
      <DiagramNode
        label="Producer"
        active={hl === 'producer' || hl === 'arrow-to-producer'}
        x="2%"
        y="30%"
        width="72px"
        height="72px"
        icon={Radio}
      />
      <DiagramNode
        label="Partitioner"
        active={hl === 'partitioner'}
        x="18%"
        y="30%"
        width="72px"
        height="72px"
        icon={Hash}
      />
      <DiagramNode
        label="Leader"
        active={hl === 'leader' || hl === 'arrow-to-leader'}
        x="40%"
        y="12%"
        width="80px"
        height="72px"
        icon={Server}
      />
      <DiagramNode
        label="Follower 1"
        active={hl === 'followers'}
        x="62%"
        y="2%"
        width="72px"
        height="60px"
        icon={Database}
      />
      <DiagramNode
        label="Follower 2"
        active={hl === 'followers'}
        x="62%"
        y="52%"
        width="72px"
        height="60px"
        icon={Database}
      />
      <DiagramNode
        label="Consumer"
        active={hl === 'consumer' || hl === 'arrow-to-consumer'}
        x="84%"
        y="30%"
        width="72px"
        height="72px"
        icon={Users}
      />

      {/* Arrows */}
      <AnimatedArrow
        active={hl === 'producer' || hl === 'partitioner'}
        x1="11%"
        y1="50%"
        x2="18%"
        y2="50%"
        direction="right1"
      />
      <AnimatedArrow
        active={hl === 'arrow-to-leader' && currentStep <= 3}
        x1="26%"
        y1="48%"
        x2="40%"
        y2="42%"
        direction="right2"
      />
      <AnimatedArrow
        active={hl === 'followers'}
        x1="56%"
        y1="30%"
        x2="62%"
        y2="24%"
        direction="up"
      />
      <AnimatedArrow
        active={hl === 'followers'}
        x1="56%"
        y1="55%"
        x2="62%"
        y2="60%"
        direction="down"
      />
      <AnimatedArrow
        active={hl === 'arrow-to-producer'}
        x1="40%"
        y1="56%"
        x2="11%"
        y2="56%"
        direction="left"
      />
      <AnimatedArrow
        active={hl === 'arrow-to-leader' && currentStep >= 9}
        x1="84%"
        y1="52%"
        x2="56%"
        y2="42%"
        direction="left2"
      />
      <AnimatedArrow
        active={hl === 'arrow-to-consumer'}
        x1="56%"
        y1="42%"
        x2="84%"
        y2="50%"
        direction="right3"
      />
    </div>
  )
}

export default function MessageFlowAnimation() {
  const { darkMode } = useTheme()
  const [currentStep, setCurrentStep] = useState(1)
  const [autoPlay, setAutoPlay] = useState(false)
  const timerRef = useRef(null)

  const totalSteps = FLOW_STEPS.length

  const goNext = useCallback(() => {
    setCurrentStep((s) => (s < totalSteps ? s + 1 : s))
  }, [totalSteps])

  const goPrev = useCallback(() => {
    setCurrentStep((s) => (s > 1 ? s - 1 : s))
  }, [])

  const reset = useCallback(() => {
    setAutoPlay(false)
    setCurrentStep(1)
  }, [])

  const toggleAutoPlay = useCallback(() => {
    setAutoPlay((prev) => !prev)
  }, [])

  // Auto-play timer
  useEffect(() => {
    if (autoPlay) {
      timerRef.current = setInterval(() => {
        setCurrentStep((s) => {
          if (s >= totalSteps) {
            setAutoPlay(false)
            return s
          }
          return s + 1
        })
      }, 2000)
    } else {
      if (timerRef.current) clearInterval(timerRef.current)
    }
    return () => {
      if (timerRef.current) clearInterval(timerRef.current)
    }
  }, [autoPlay, totalSteps])

  // Keyboard navigation
  useEffect(() => {
    const handleKey = (e) => {
      if (e.key === 'ArrowRight') goNext()
      if (e.key === 'ArrowLeft') goPrev()
      if (e.key === ' ') {
        e.preventDefault()
        toggleAutoPlay()
      }
    }
    window.addEventListener('keydown', handleKey)
    return () => window.removeEventListener('keydown', handleKey)
  }, [goNext, goPrev, toggleAutoPlay])

  const step = FLOW_STEPS[currentStep - 1]

  return (
    <div className="rounded-xl border border-[var(--border)] bg-[var(--bg-surface)] p-4 sm:p-6 my-6">
      {/* Header */}
      <div className="flex items-center justify-between mb-4">
        <div className="flex items-center gap-2">
          <Package className="w-5 h-5 text-[var(--accent)]" />
          <h3 className="font-heading text-lg font-semibold text-[var(--text-primary)]">
            Message Flow Animation
          </h3>
        </div>
        <span className="font-mono text-xs text-[var(--text-tertiary)]">
          Step {currentStep} of {totalSteps}
        </span>
      </div>

      {/* Progress bar */}
      <div className="w-full h-1.5 rounded-full bg-[var(--bg-surface-alt)] mb-4 overflow-hidden">
        <div
          className="h-full rounded-full bg-[var(--accent)] transition-all duration-500 ease-out"
          style={{ width: `${(currentStep / totalSteps) * 100}%` }}
        />
      </div>

      {/* Step indicators */}
      <div className="flex gap-1 mb-4 overflow-x-auto pb-1">
        {FLOW_STEPS.map((s) => (
          <button
            key={s.id}
            onClick={() => {
              setAutoPlay(false)
              setCurrentStep(s.id)
            }}
            className={`flex-shrink-0 w-7 h-7 rounded-md text-[10px] font-mono font-bold transition-all duration-200 ${
              s.id === currentStep
                ? 'bg-[var(--accent)] text-white scale-110 shadow-md'
                : s.id < currentStep
                  ? 'bg-emerald-500/20 text-emerald-600 dark:text-emerald-400 border border-emerald-500/30'
                  : 'bg-[var(--bg-surface-alt)] text-[var(--text-tertiary)] border border-[var(--border)]'
            }`}
            title={s.label}
          >
            {s.id}
          </button>
        ))}
      </div>

      {/* Main content area */}
      <div className="flex flex-col lg:flex-row gap-4">
        {/* Diagram */}
        <div className="flex-1 rounded-lg border border-[var(--border)] bg-[var(--bg-surface-alt)] p-3 min-h-[260px] sm:min-h-[220px]">
          <FlowDiagram currentStep={currentStep} />
        </div>

        {/* Explanation panel */}
        <div className="lg:w-[320px] flex-shrink-0 flex flex-col gap-3">
          {/* Step title */}
          <div className="flex items-center gap-2 animate-fade-in" key={step.id}>
            <div className="flex items-center justify-center w-8 h-8 rounded-lg bg-[var(--accent)]/10 border border-[var(--accent)]/30">
              <step.icon className="w-4 h-4 text-[var(--accent)]" />
            </div>
            <div>
              <div className="font-mono text-[10px] text-[var(--text-tertiary)] uppercase tracking-wider">
                Step {step.id}
              </div>
              <div className="font-heading font-semibold text-sm text-[var(--text-primary)]">
                {step.label}
              </div>
            </div>
          </div>

          {/* Detail text */}
          <p
            className="text-xs font-body leading-relaxed text-[var(--text-secondary)] animate-fade-in"
            key={`detail-${step.id}`}
          >
            {step.detail}
          </p>

          {/* Code snippet */}
          <div
            className="rounded-lg border border-[var(--border)] bg-[var(--code-bg)] p-3 animate-fade-in"
            key={`code-${step.id}`}
          >
            <pre className="text-[11px] font-mono text-[var(--text-primary)] whitespace-pre-wrap leading-relaxed overflow-x-auto">
              {step.code}
            </pre>
          </div>
        </div>
      </div>

      {/* Controls */}
      <div className="flex items-center justify-center gap-2 mt-4 pt-3 border-t border-[var(--border)]">
        <button
          onClick={reset}
          className="flex items-center gap-1 px-3 py-1.5 rounded-lg border border-[var(--border)] bg-[var(--bg-surface)] hover:bg-[var(--bg-hover)] text-[var(--text-secondary)] text-xs font-body transition-colors"
          title="Reset to step 1"
        >
          <RotateCcw className="w-3.5 h-3.5" />
          Reset
        </button>
        <button
          onClick={goPrev}
          disabled={currentStep <= 1}
          className="flex items-center gap-1 px-3 py-1.5 rounded-lg border border-[var(--border)] bg-[var(--bg-surface)] hover:bg-[var(--bg-hover)] text-[var(--text-secondary)] text-xs font-body transition-colors disabled:opacity-40"
          title="Previous step"
        >
          <ChevronLeft className="w-3.5 h-3.5" />
          Prev
        </button>
        <button
          onClick={toggleAutoPlay}
          className={`flex items-center gap-1 px-4 py-1.5 rounded-lg text-xs font-body transition-colors ${
            autoPlay
              ? 'bg-amber-500 hover:bg-amber-600 text-white'
              : 'bg-emerald-500 hover:bg-emerald-600 text-white'
          }`}
          title={autoPlay ? 'Pause auto-play' : 'Start auto-play'}
        >
          {autoPlay ? (
            <>
              <Pause className="w-3.5 h-3.5" />
              Pause
            </>
          ) : (
            <>
              <Play className="w-3.5 h-3.5" />
              Auto-play
            </>
          )}
        </button>
        <button
          onClick={goNext}
          disabled={currentStep >= totalSteps}
          className="flex items-center gap-1 px-3 py-1.5 rounded-lg border border-[var(--border)] bg-[var(--bg-surface)] hover:bg-[var(--bg-hover)] text-[var(--text-secondary)] text-xs font-body transition-colors disabled:opacity-40"
          title="Next step"
        >
          Next
          <ChevronRight className="w-3.5 h-3.5" />
        </button>
      </div>

      {/* Keyboard hint */}
      <p className="text-center text-[10px] text-[var(--text-tertiary)] font-mono mt-2">
        Arrow keys to navigate, Space to toggle auto-play
      </p>
    </div>
  )
}

import { useState, useMemo, useCallback, useEffect, useRef } from 'react'
import { useTheme } from '../../context/ThemeContext'
import { RotateCcw, Skull, Plus, Server, Radio, Users } from 'lucide-react'

// Generate partition-to-broker assignment with round-robin leader placement
function generateAssignment(partitionCount, brokerCount, replicationFactor) {
  const rf = Math.min(replicationFactor, brokerCount)
  const partitions = []
  for (let p = 0; p < partitionCount; p++) {
    const leader = p % brokerCount
    const replicas = [leader]
    for (let r = 1; r < rf; r++) {
      replicas.push((leader + r) % brokerCount)
    }
    partitions.push({ id: p, leader, replicas })
  }
  return partitions
}

function PartitionCell({ partition, brokerId, isLeader, isHighlighted, isDead, onHover, onLeave }) {
  const bg = isDead
    ? 'bg-red-100 dark:bg-red-900/20 border-red-300 dark:border-red-700'
    : isLeader
      ? 'bg-[var(--accent-subtle)] border-[var(--accent)]'
      : isHighlighted
        ? 'bg-blue-50 dark:bg-blue-900/20 border-blue-300 dark:border-blue-600'
        : 'bg-[var(--bg-surface-alt)] border-[var(--border)]'

  return (
    <div
      className={`relative rounded border-2 px-2 py-1 text-center transition-all duration-200 cursor-default select-none ${bg}`}
      onMouseEnter={() => onHover(partition.id)}
      onMouseLeave={onLeave}
    >
      <div className="font-mono text-[10px] text-[var(--text-tertiary)]">P{partition.id}</div>
      {isLeader && !isDead && (
        <div className="font-mono text-[9px] font-bold text-[var(--accent)]" title="Leader">L</div>
      )}
      {!isLeader && !isDead && (
        <div className="font-mono text-[9px] text-[var(--text-tertiary)]" title="Follower">F</div>
      )}
      {isDead && (
        <div className="font-mono text-[9px] text-red-500 line-through">X</div>
      )}
    </div>
  )
}

function BrokerBox({ brokerId, partitions, deadBrokers, highlightedPartition, onPartitionHover, onPartitionLeave, onKill }) {
  const isDead = deadBrokers.has(brokerId)
  const myPartitions = partitions.filter(p => p.replicas.includes(brokerId))

  // Build tooltip: which partitions this broker leads / follows
  const leaderFor = myPartitions.filter(p => p.leader === brokerId).map(p => `P${p.id}`)
  const followerFor = myPartitions.filter(p => p.leader !== brokerId).map(p => `P${p.id}`)

  const tooltipText = isDead
    ? `Broker ${brokerId + 1} -- DOWN`
    : `Broker ${brokerId + 1} -- Leader: ${leaderFor.join(', ') || 'none'} | Follower: ${followerFor.join(', ') || 'none'}`

  return (
    <div className="relative group" title={tooltipText}>
      <div
        className={`rounded-lg border-2 p-3 transition-all duration-300 ${
          isDead
            ? 'border-red-400 dark:border-red-600 bg-red-50 dark:bg-red-900/10 opacity-60'
            : 'border-[var(--border-strong)] bg-[var(--bg-surface)]'
        }`}
      >
        {/* Broker header */}
        <div className="flex items-center justify-between mb-2">
          <div className="flex items-center gap-1.5">
            <Server size={12} className={isDead ? 'text-red-500' : 'text-[var(--text-tertiary)]'} />
            <span className={`font-mono text-[11px] font-semibold ${isDead ? 'text-red-500 line-through' : 'text-[var(--text-primary)]'}`}>
              Broker {brokerId + 1}
            </span>
          </div>
          {!isDead && (
            <button
              onClick={() => onKill(brokerId)}
              className="opacity-0 group-hover:opacity-100 transition-opacity p-0.5 rounded hover:bg-red-50 dark:hover:bg-red-900/20"
              title={`Kill Broker ${brokerId + 1}`}
            >
              <Skull size={11} className="text-red-400 hover:text-red-600" />
            </button>
          )}
        </div>

        {/* Partitions grid */}
        <div className="grid grid-cols-3 gap-1.5">
          {myPartitions.map(p => (
            <PartitionCell
              key={p.id}
              partition={p}
              brokerId={brokerId}
              isLeader={p.leader === brokerId}
              isHighlighted={highlightedPartition === p.id}
              isDead={isDead}
              onHover={onPartitionHover}
              onLeave={onPartitionLeave}
            />
          ))}
          {myPartitions.length === 0 && (
            <div className="col-span-3 text-center font-mono text-[10px] text-[var(--text-tertiary)] py-2">
              No partitions
            </div>
          )}
        </div>
      </div>

      {/* Dead overlay */}
      {isDead && (
        <div className="absolute inset-0 flex items-center justify-center pointer-events-none">
          <span className="font-mono text-[10px] font-bold text-red-500 bg-red-100 dark:bg-red-900/40 px-2 py-0.5 rounded">
            DOWN
          </span>
        </div>
      )}
    </div>
  )
}

function AnimatedDot({ active, direction }) {
  if (!active) return null
  const animClass = direction === 'right'
    ? 'animate-[flowRight_1.5s_ease-in-out_infinite]'
    : 'animate-[flowLeft_1.5s_ease-in-out_infinite]'

  return (
    <div className={`w-2 h-2 rounded-full bg-[var(--accent)] ${animClass}`} />
  )
}

function ProducerNode({ label, active }) {
  return (
    <div className="flex items-center gap-2">
      <div className={`rounded-lg border-2 px-3 py-2 text-center transition-all ${
        active ? 'border-[var(--accent)] bg-[var(--accent-subtle)]' : 'border-[var(--border)] bg-[var(--bg-surface)]'
      }`}>
        <Radio size={12} className={active ? 'text-[var(--accent)] mx-auto mb-0.5' : 'text-[var(--text-tertiary)] mx-auto mb-0.5'} />
        <div className="font-mono text-[10px] font-semibold">{label}</div>
      </div>
      <div className="flex items-center gap-1">
        <div className="w-6 h-[2px] bg-[var(--border-strong)]" />
        <AnimatedDot active={active} direction="right" />
        <div className="w-3 h-[2px] bg-[var(--border-strong)]" />
        <div className="w-0 h-0 border-t-[4px] border-t-transparent border-b-[4px] border-b-transparent border-l-[6px] border-l-[var(--border-strong)]" />
      </div>
    </div>
  )
}

function ConsumerNode({ label, active }) {
  return (
    <div className="flex items-center gap-2">
      <div className="flex items-center gap-1">
        <div className="w-0 h-0 border-t-[4px] border-t-transparent border-b-[4px] border-b-transparent border-l-[6px] border-l-[var(--border-strong)]" />
        <div className="w-3 h-[2px] bg-[var(--border-strong)]" />
        <AnimatedDot active={active} direction="right" />
        <div className="w-6 h-[2px] bg-[var(--border-strong)]" />
      </div>
      <div className={`rounded-lg border-2 px-3 py-2 text-center transition-all ${
        active ? 'border-emerald-500 bg-emerald-50 dark:bg-emerald-900/10' : 'border-[var(--border)] bg-[var(--bg-surface)]'
      }`}>
        <Users size={12} className={active ? 'text-emerald-500 mx-auto mb-0.5' : 'text-[var(--text-tertiary)] mx-auto mb-0.5'} />
        <div className="font-mono text-[10px] font-semibold">{label}</div>
      </div>
    </div>
  )
}

export default function KafkaClusterDiagram() {
  const { darkMode } = useTheme()
  const [brokerCount, setBrokerCount] = useState(3)
  const [partitionCount, setPartitionCount] = useState(6)
  const [replicationFactor, setReplicationFactor] = useState(3)
  const [deadBrokers, setDeadBrokers] = useState(new Set())
  const [highlightedPartition, setHighlightedPartition] = useState(null)
  const [animating, setAnimating] = useState(true)
  const [electionLog, setElectionLog] = useState([])

  const baseAssignment = useMemo(
    () => generateAssignment(partitionCount, brokerCount, replicationFactor),
    [partitionCount, brokerCount, replicationFactor]
  )

  // Apply leader re-election for dead brokers
  const assignment = useMemo(() => {
    if (deadBrokers.size === 0) return baseAssignment
    return baseAssignment.map(p => {
      if (deadBrokers.has(p.leader)) {
        // Find first alive replica to become leader
        const newLeader = p.replicas.find(b => !deadBrokers.has(b))
        return { ...p, leader: newLeader !== undefined ? newLeader : -1 }
      }
      return p
    })
  }, [baseAssignment, deadBrokers])

  const handleKillBroker = useCallback((brokerId) => {
    setDeadBrokers(prev => {
      const next = new Set(prev)
      next.add(brokerId)
      return next
    })

    // Log which partitions needed re-election
    const affected = baseAssignment.filter(p => p.leader === brokerId)
    if (affected.length > 0) {
      const logEntry = `Broker ${brokerId + 1} killed. Leader election for: ${affected.map(p => `P${p.id}`).join(', ')}`
      setElectionLog(prev => [...prev, logEntry])
    }
  }, [baseAssignment])

  const handleAddPartition = useCallback(() => {
    setPartitionCount(prev => prev + 1)
    setElectionLog(prev => [...prev, `Added partition P${partitionCount} with ${Math.min(replicationFactor, brokerCount)} replicas`])
  }, [partitionCount, replicationFactor, brokerCount])

  const handleReset = useCallback(() => {
    setBrokerCount(3)
    setPartitionCount(6)
    setReplicationFactor(3)
    setDeadBrokers(new Set())
    setHighlightedPartition(null)
    setElectionLog([])
  }, [])

  const aliveBrokerCount = brokerCount - deadBrokers.size

  return (
    <div className="my-6 rounded-xl border border-[var(--border)] bg-[var(--bg-surface)] not-prose overflow-hidden">
      {/* Header */}
      <div className="flex items-center justify-between px-5 py-3 border-b border-[var(--border)] bg-[var(--bg-surface-alt)]">
        <div className="flex items-center gap-2">
          <Server size={14} className="text-[var(--accent)]" />
          <span className="font-mono text-[10px] font-semibold uppercase tracking-[0.06em] text-[var(--accent)]">
            Interactive: Kafka Cluster Diagram
          </span>
        </div>
        <button
          onClick={handleReset}
          className="flex items-center gap-1 font-mono text-[10px] text-[var(--text-tertiary)] hover:text-[var(--text-primary)] transition-colors"
        >
          <RotateCcw size={10} /> Reset
        </button>
      </div>

      <div className="p-5">
        {/* Controls */}
        <div className="flex flex-wrap items-center gap-4 mb-5">
          {/* Broker count */}
          <div className="flex items-center gap-2">
            <span className="font-mono text-[11px] text-[var(--text-secondary)]">Brokers:</span>
            <div className="flex items-center gap-1">
              <button
                onClick={() => { setBrokerCount(Math.max(1, brokerCount - 1)); setDeadBrokers(new Set()); setElectionLog([]); }}
                className="w-7 h-7 rounded border border-[var(--border)] flex items-center justify-center font-mono text-sm hover:bg-[var(--bg-hover)] transition-colors"
              >-</button>
              <span className="w-6 text-center font-mono font-semibold text-[var(--accent)]">{brokerCount}</span>
              <button
                onClick={() => { setBrokerCount(Math.min(6, brokerCount + 1)); setDeadBrokers(new Set()); setElectionLog([]); }}
                className="w-7 h-7 rounded border border-[var(--border)] flex items-center justify-center font-mono text-sm hover:bg-[var(--bg-hover)] transition-colors"
              >+</button>
            </div>
          </div>

          {/* Replication factor */}
          <div className="flex items-center gap-2">
            <span className="font-mono text-[11px] text-[var(--text-secondary)]">RF:</span>
            <div className="flex items-center gap-1">
              <button
                onClick={() => setReplicationFactor(Math.max(1, replicationFactor - 1))}
                className="w-7 h-7 rounded border border-[var(--border)] flex items-center justify-center font-mono text-sm hover:bg-[var(--bg-hover)] transition-colors"
              >-</button>
              <span className="w-6 text-center font-mono font-semibold text-[var(--accent)]">{replicationFactor}</span>
              <button
                onClick={() => setReplicationFactor(Math.min(brokerCount, replicationFactor + 1))}
                className="w-7 h-7 rounded border border-[var(--border)] flex items-center justify-center font-mono text-sm hover:bg-[var(--bg-hover)] transition-colors"
              >+</button>
            </div>
          </div>

          {/* Action buttons */}
          <button
            onClick={handleAddPartition}
            className="flex items-center gap-1 px-2.5 py-1.5 rounded border border-[var(--border)] font-mono text-[10px] text-[var(--text-secondary)] hover:bg-[var(--bg-hover)] hover:text-[var(--text-primary)] transition-colors"
          >
            <Plus size={10} /> Add Partition
          </button>

          <span className="font-mono text-[10px] text-[var(--text-tertiary)]">
            Topic: audit-events &middot; {partitionCount} partitions &middot; RF={Math.min(replicationFactor, brokerCount)} &middot; {aliveBrokerCount}/{brokerCount} alive
          </span>
        </div>

        {/* Main diagram: Producers -> Brokers -> Consumers */}
        <div className="flex items-center gap-4 overflow-x-auto pb-2">
          {/* Producers */}
          <div className="flex flex-col gap-3 shrink-0">
            <ProducerNode label="Producer 1" active={animating && aliveBrokerCount > 0} />
            <ProducerNode label="Producer 2" active={animating && aliveBrokerCount > 0} />
          </div>

          {/* Brokers */}
          <div className={`flex-1 grid gap-3 min-w-0`} style={{ gridTemplateColumns: `repeat(${brokerCount}, minmax(120px, 1fr))` }}>
            {Array.from({ length: brokerCount }, (_, i) => (
              <BrokerBox
                key={i}
                brokerId={i}
                partitions={assignment}
                deadBrokers={deadBrokers}
                highlightedPartition={highlightedPartition}
                onPartitionHover={setHighlightedPartition}
                onPartitionLeave={() => setHighlightedPartition(null)}
                onKill={handleKillBroker}
              />
            ))}
          </div>

          {/* Consumers */}
          <div className="flex flex-col gap-3 shrink-0">
            <ConsumerNode label="Group A" active={animating && aliveBrokerCount > 0} />
            <ConsumerNode label="Group B" active={animating && aliveBrokerCount > 0} />
          </div>
        </div>

        {/* Legend */}
        <div className="flex flex-wrap items-center gap-4 mt-4 pt-3 border-t border-[var(--border)]">
          <div className="flex items-center gap-1.5">
            <div className="w-3 h-3 rounded border-2 border-[var(--accent)] bg-[var(--accent-subtle)]" />
            <span className="font-mono text-[10px] text-[var(--text-tertiary)]">Leader (L)</span>
          </div>
          <div className="flex items-center gap-1.5">
            <div className="w-3 h-3 rounded border-2 border-[var(--border)] bg-[var(--bg-surface-alt)]" />
            <span className="font-mono text-[10px] text-[var(--text-tertiary)]">Follower (F)</span>
          </div>
          <div className="flex items-center gap-1.5">
            <div className="w-3 h-3 rounded border-2 border-blue-300 bg-blue-50 dark:bg-blue-900/20" />
            <span className="font-mono text-[10px] text-[var(--text-tertiary)]">Highlighted replica</span>
          </div>
          <div className="flex items-center gap-1.5">
            <div className="w-3 h-3 rounded border-2 border-red-400 bg-red-100 dark:bg-red-900/20" />
            <span className="font-mono text-[10px] text-[var(--text-tertiary)]">Dead broker</span>
          </div>
        </div>

        {/* Election log */}
        {electionLog.length > 0 && (
          <div className="mt-4 pt-3 border-t border-[var(--border)]">
            <div className="font-mono text-[10px] font-semibold uppercase tracking-[0.06em] text-[var(--text-tertiary)] mb-2">
              Event Log
            </div>
            <div className="max-h-24 overflow-y-auto space-y-1">
              {electionLog.map((entry, i) => (
                <div key={i} className="font-mono text-[11px] text-[var(--text-secondary)] flex items-start gap-2">
                  <span className="text-[var(--text-tertiary)] shrink-0">[{i + 1}]</span>
                  <span>{entry}</span>
                </div>
              ))}
            </div>
          </div>
        )}

        <p className="mt-4 text-[12px] text-[var(--text-tertiary)] italic">
          Hover partitions to highlight all replicas. Hover brokers to see leader/follower summary. Click the skull icon to simulate a broker failure.
        </p>
      </div>

      {/* Keyframes for animation */}
      <style>{`
        @keyframes flowRight {
          0% { transform: translateX(-8px); opacity: 0; }
          20% { opacity: 1; }
          80% { opacity: 1; }
          100% { transform: translateX(8px); opacity: 0; }
        }
        @keyframes flowLeft {
          0% { transform: translateX(8px); opacity: 0; }
          20% { opacity: 1; }
          80% { opacity: 1; }
          100% { transform: translateX(-8px); opacity: 0; }
        }
      `}</style>
    </div>
  )
}

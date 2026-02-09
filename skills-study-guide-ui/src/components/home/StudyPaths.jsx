import { useNavigate } from 'react-router-dom'
import { Zap, Target, BookOpen, Clock } from 'lucide-react'

const PATHS = [
  {
    id: 'fast-track',
    name: 'Google Interview Fast Track',
    description: 'Core topics for technical + behavioral rounds',
    icon: Zap,
    time: '~6 hours',
    guides: ['14', '03', '09', '01', '13'],
  },
  {
    id: 'backend',
    name: 'Backend Engineer Foundation',
    description: 'Full stack of backend technologies',
    icon: Target,
    time: '~10 hours',
    guides: ['09', '10', '01', '12', '05', '04', '03', '11'],
  },
  {
    id: 'day-before',
    name: 'Day Before Interview',
    description: 'Key templates and story matrix',
    icon: Clock,
    time: '~3 hours',
    guides: ['14', '03', '13'],
  },
]

export default function StudyPaths() {
  const navigate = useNavigate()

  return (
    <section className="mb-12">
      <h2 className="font-mono text-xs font-semibold uppercase tracking-[0.08em] text-[var(--text-tertiary)] mb-4 px-1">
        Study Paths
      </h2>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        {PATHS.map(path => {
          const Icon = path.icon
          return (
            <button
              key={path.id}
              onClick={() => navigate(`/guide/${getFirstGuideSlug(path.guides[0])}`)}
              className="group text-left p-5 rounded-xl border border-[var(--border)] bg-[var(--bg-surface)] hover:border-[var(--accent)] hover:shadow-sm transition-all"
            >
              <div className="flex items-center gap-3 mb-3">
                <div className="w-8 h-8 rounded-lg bg-[var(--accent-subtle)] flex items-center justify-center">
                  <Icon size={16} className="text-[var(--accent)]" />
                </div>
                <span className="font-mono text-[10px] text-[var(--text-tertiary)] uppercase tracking-wide">
                  {path.time}
                </span>
              </div>

              <h3 className="font-heading font-semibold text-base mb-1 group-hover:text-[var(--accent)] transition-colors">
                {path.name}
              </h3>
              <p className="text-[12px] text-[var(--text-secondary)] mb-3 leading-relaxed">
                {path.description}
              </p>

              <div className="flex gap-1.5 flex-wrap">
                {path.guides.map(num => (
                  <span key={num} className="font-mono text-[10px] px-2 py-0.5 rounded bg-[var(--bg-surface-alt)] text-[var(--text-tertiary)] border border-[var(--border)]">
                    {num}
                  </span>
                ))}
              </div>
            </button>
          )
        })}
      </div>
    </section>
  )
}

function getFirstGuideSlug(num) {
  // Map guide number to slug
  const slugMap = {
    '01': '01-kafka-avro-schema-registry',
    '03': '03-system-design-distributed-systems',
    '09': '09-java-java17-deep-dive',
    '13': '13-behavioral-googleyness-prep',
    '14': '14-dsa-coding-patterns',
  }
  return slugMap[num] || `${num}-guide`
}

import { useMemo } from 'react'
import { useNavigate } from 'react-router-dom'
import { ArrowRight, BookOpen } from 'lucide-react'
import guideIndex from '../../guide-index.json'
import { getLastReadGuide } from '../../hooks/useProgress'

export default function ContinueReading() {
  const navigate = useNavigate()

  const lastRead = useMemo(() => {
    const data = getLastReadGuide()
    if (!data) return null

    const guide = guideIndex.find(g => g.slug === data.slug)
    if (!guide) return null

    return { ...guide, percent: data.percent, lastVisited: data.lastVisited }
  }, [])

  if (!lastRead) return null

  const timeAgo = getTimeAgo(lastRead.lastVisited)

  return (
    <section className="max-w-5xl mx-auto px-5 mb-10">
      <p className="font-mono text-[10px] font-semibold uppercase tracking-[0.08em] text-[var(--text-tertiary)] mb-3">
        Continue Reading
      </p>

      <button
        onClick={() => navigate(`/guide/${lastRead.slug}`)}
        className="group w-full text-left p-5 rounded-xl border border-[var(--accent)] bg-[var(--accent-subtle)] hover:shadow-md transition-all"
      >
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-4 min-w-0">
            <div className="shrink-0 w-10 h-10 rounded-lg bg-[var(--accent)] flex items-center justify-center">
              <BookOpen size={18} className="text-white" />
            </div>
            <div className="min-w-0">
              <p className="font-heading font-semibold text-lg truncate group-hover:text-[var(--accent)] transition-colors">
                {lastRead.title}
              </p>
              <p className="text-[12px] text-[var(--text-secondary)] font-mono">
                {lastRead.percent}% complete &middot; {timeAgo}
              </p>
            </div>
          </div>

          <ArrowRight size={18} className="shrink-0 text-[var(--accent)] group-hover:translate-x-1 transition-transform" />
        </div>

        {/* Progress bar */}
        <div className="mt-3 h-1.5 rounded-full bg-[var(--border)] overflow-hidden">
          <div
            className="h-full rounded-full bg-[var(--accent)] transition-all duration-500"
            style={{ width: `${lastRead.percent}%` }}
          />
        </div>
      </button>
    </section>
  )
}

function getTimeAgo(timestamp) {
  if (!timestamp) return ''
  const diff = Date.now() - timestamp
  const mins = Math.floor(diff / 60000)
  if (mins < 1) return 'Just now'
  if (mins < 60) return `${mins}m ago`
  const hours = Math.floor(mins / 60)
  if (hours < 24) return `${hours}h ago`
  const days = Math.floor(hours / 24)
  return `${days}d ago`
}

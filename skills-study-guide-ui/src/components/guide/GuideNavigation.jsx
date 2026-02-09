import { useNavigate } from 'react-router-dom'
import { ChevronLeft, ChevronRight } from 'lucide-react'

export default function GuideNavigation({ prev, next }) {
  const navigate = useNavigate()

  return (
    <div className="flex items-stretch gap-4 mt-16 pt-8 border-t border-[var(--border)]">
      {prev ? (
        <button
          onClick={() => navigate(`/guide/${prev.slug}`)}
          className="flex-1 text-left p-4 rounded-xl border border-[var(--border)] hover:border-[var(--border-strong)] hover:bg-[var(--bg-hover)] transition-all group"
        >
          <span className="flex items-center gap-1 text-[11px] font-mono text-[var(--text-tertiary)] uppercase tracking-wide mb-1">
            <ChevronLeft size={12} /> Previous
          </span>
          <span className="text-sm font-heading font-semibold group-hover:text-[var(--accent)] transition-colors">
            {prev.title}
          </span>
        </button>
      ) : <div className="flex-1" />}

      {next ? (
        <button
          onClick={() => navigate(`/guide/${next.slug}`)}
          className="flex-1 text-right p-4 rounded-xl border border-[var(--border)] hover:border-[var(--border-strong)] hover:bg-[var(--bg-hover)] transition-all group"
        >
          <span className="flex items-center justify-end gap-1 text-[11px] font-mono text-[var(--text-tertiary)] uppercase tracking-wide mb-1">
            Next <ChevronRight size={12} />
          </span>
          <span className="text-sm font-heading font-semibold group-hover:text-[var(--accent)] transition-colors">
            {next.title}
          </span>
        </button>
      ) : <div className="flex-1" />}
    </div>
  )
}

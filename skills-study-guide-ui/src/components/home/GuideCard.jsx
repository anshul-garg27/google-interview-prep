import { useNavigate } from 'react-router-dom'
import { Code2, GitBranch, HelpCircle } from 'lucide-react'

export default function GuideCard({ guide, progress }) {
  const navigate = useNavigate()
  const pct = progress?.[guide.slug]?.percent || 0

  return (
    <button
      onClick={() => navigate(`/guide/${guide.slug}`)}
      className="group text-left w-full p-5 rounded-xl border border-[var(--border)] bg-[var(--bg-surface)] hover:border-[var(--border-strong)] hover:-translate-y-0.5 hover:shadow-sm transition-all duration-200 relative overflow-hidden"
    >
      {/* Top row: number + reading time */}
      <div className="flex items-center justify-between mb-3">
        <span className="font-mono text-[13px] font-semibold text-[var(--text-tertiary)]">
          {guide.num}
        </span>
        <span className="font-mono text-[10px] text-[var(--text-tertiary)] uppercase tracking-wide">
          {guide.readingTime} min
        </span>
      </div>

      {/* Title */}
      <h3 className="font-heading font-semibold text-lg leading-snug mb-2 group-hover:text-[var(--accent)] transition-colors">
        {guide.title}
      </h3>

      {/* Description */}
      {guide.description && (
        <p className="text-[13px] text-[var(--text-secondary)] leading-relaxed line-clamp-2 mb-4">
          {guide.description}
        </p>
      )}

      {/* Stats row */}
      <div className="flex items-center gap-4 text-[var(--text-tertiary)]">
        {guide.stats.codeBlocks > 0 && (
          <span className="flex items-center gap-1 font-mono text-[10px]">
            <Code2 size={11} /> {guide.stats.codeBlocks}
          </span>
        )}
        {guide.stats.mermaidDiagrams > 0 && (
          <span className="flex items-center gap-1 font-mono text-[10px]">
            <GitBranch size={11} /> {guide.stats.mermaidDiagrams}
          </span>
        )}
        {guide.stats.qaItems > 0 && (
          <span className="flex items-center gap-1 font-mono text-[10px]">
            <HelpCircle size={11} /> {guide.stats.qaItems} Q&A
          </span>
        )}
      </div>

      {/* Tags */}
      {guide.tags.length > 0 && (
        <div className="flex flex-wrap gap-1.5 mt-3">
          {guide.tags.slice(0, 4).map(tag => (
            <span
              key={tag}
              className="font-mono text-[9px] px-2 py-0.5 rounded-full bg-[var(--bg-surface-alt)] text-[var(--text-tertiary)] border border-[var(--border)]"
            >
              {tag}
            </span>
          ))}
        </div>
      )}

      {/* Progress bar */}
      {pct > 0 && (
        <div className="absolute bottom-0 left-0 right-0 h-1 bg-[var(--border)]">
          <div
            className="h-full bg-[var(--accent)] transition-all duration-500"
            style={{ width: `${pct}%` }}
          />
        </div>
      )}
    </button>
  )
}

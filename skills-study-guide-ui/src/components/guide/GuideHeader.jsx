import { Code2, GitBranch, HelpCircle, Clock } from 'lucide-react'
import guideIndex from '../../guide-index.json'

export default function GuideHeader({ guide }) {
  return (
    <header className="mb-10 pb-8 border-b border-[var(--border)]">
      <div className="flex items-center gap-3 mb-3">
        <span className="font-mono text-xs font-semibold uppercase tracking-[0.06em] text-[var(--accent)] px-2.5 py-1 rounded-md bg-[var(--accent-subtle)]">
          {guide.category}
        </span>
        <span className="font-mono text-xs text-[var(--text-tertiary)]">
          Guide {guide.num} of {guideIndex.length}
        </span>
      </div>

      <h1 className="font-heading font-extrabold text-3xl md:text-4xl tracking-tight leading-[1.15] mb-4">
        {guide.title}
      </h1>

      {guide.description && (
        <p className="text-[var(--text-secondary)] text-lg leading-relaxed max-w-2xl mb-6">
          {guide.description}
        </p>
      )}

      <div className="flex flex-wrap items-center gap-5 text-[var(--text-tertiary)]">
        <span className="flex items-center gap-1.5 font-mono text-xs">
          <Clock size={13} /> {guide.readingTime} min read
        </span>
        {guide.stats.codeBlocks > 0 && (
          <span className="flex items-center gap-1.5 font-mono text-xs">
            <Code2 size={13} /> {guide.stats.codeBlocks} code blocks
          </span>
        )}
        {guide.stats.mermaidDiagrams > 0 && (
          <span className="flex items-center gap-1.5 font-mono text-xs">
            <GitBranch size={13} /> {guide.stats.mermaidDiagrams} diagrams
          </span>
        )}
        {guide.stats.qaItems > 0 && (
          <span className="flex items-center gap-1.5 font-mono text-xs">
            <HelpCircle size={13} /> {guide.stats.qaItems} Q&A
          </span>
        )}
      </div>
    </header>
  )
}

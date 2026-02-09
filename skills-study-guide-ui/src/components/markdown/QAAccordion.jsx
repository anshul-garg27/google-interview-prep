import { useState } from 'react'
import { ChevronRight, Eye, EyeOff } from 'lucide-react'

export default function QAAccordion({ question, children }) {
  const [open, setOpen] = useState(false)

  return (
    <div className="border border-[var(--border)] rounded-lg my-4 overflow-hidden">
      <button
        onClick={() => setOpen(!open)}
        className="w-full flex items-start gap-3 p-4 text-left hover:bg-[var(--bg-hover)] transition-colors"
      >
        <ChevronRight
          size={16}
          className={`shrink-0 mt-1 text-[var(--accent)] transition-transform duration-200 ${open ? 'rotate-90' : ''}`}
        />
        <span className="font-body font-semibold text-[15px] leading-relaxed flex-1">
          {question}
        </span>
      </button>

      {open && (
        <div className="px-4 pb-4 pl-11 border-t border-[var(--border)] bg-[var(--bg-surface-alt)]/50 animate-fade-in">
          <div className="pt-3 text-[var(--text-secondary)] text-[15px] leading-relaxed prose prose-sm dark:prose-invert max-w-none">
            {children}
          </div>
        </div>
      )}
    </div>
  )
}

export function QASection({ children }) {
  const [allOpen, setAllOpen] = useState(false)

  return (
    <div className="my-8">
      <div className="flex items-center justify-between mb-4">
        <span className="font-mono text-[11px] font-semibold uppercase tracking-[0.06em] text-[var(--accent)]">
          Interview Q&A
        </span>
        <button
          onClick={() => setAllOpen(!allOpen)}
          className="flex items-center gap-1.5 font-mono text-[11px] text-[var(--text-tertiary)] hover:text-[var(--text-secondary)] transition-colors"
        >
          {allOpen ? <EyeOff size={12} /> : <Eye size={12} />}
          {allOpen ? 'Collapse All' : 'Expand All'}
        </button>
      </div>
      {children}
    </div>
  )
}

import { Zap, Moon, Briefcase } from 'lucide-react'

// Usage:
// <InterviewHeader project={project} quickMode={false} onToggleMode={() => {}} />

export default function InterviewHeader({ project, quickMode, onToggleMode }) {
  return (
    <header className="mb-10 pb-8 border-b border-[var(--border)]">
      {/* Company & role badge */}
      <div className="flex flex-wrap items-center gap-3 mb-3">
        <span className="font-mono text-xs font-semibold uppercase tracking-[0.06em] text-[var(--accent)] px-2.5 py-1 rounded-md bg-[var(--accent-subtle)]">
          {project.company}
        </span>
        <span className="font-mono text-xs text-[var(--text-tertiary)] flex items-center gap-1.5">
          <Briefcase size={13} />
          {project.role}
        </span>
      </div>

      {/* Title */}
      <h1 className="font-heading font-extrabold text-3xl md:text-4xl tracking-tight leading-[1.15] mb-3">
        {project.title}
      </h1>

      {/* Subtitle */}
      {project.subtitle && (
        <p className="text-[var(--text-secondary)] text-lg leading-relaxed max-w-2xl mb-6 font-body">
          {project.subtitle}
        </p>
      )}

      {/* Key numbers grid */}
      {project.numbers?.length > 0 && (
        <div className="grid grid-cols-3 md:grid-cols-6 gap-3 mb-6">
          {project.numbers.map((num, i) => (
            <div
              key={i}
              className="text-center px-2 py-3 rounded-lg bg-[var(--bg-surface-alt)] border border-[var(--border)]"
            >
              <div className="font-heading font-bold text-lg md:text-xl text-[var(--accent)] leading-tight">
                {num.value}
              </div>
              <div className="font-mono text-[10px] md:text-[11px] uppercase tracking-[0.06em] text-[var(--text-tertiary)] mt-1 leading-tight">
                {num.label}
              </div>
            </div>
          ))}
        </div>
      )}

      {/* Quick Mode / Deep Mode toggle */}
      <button
        onClick={onToggleMode}
        className="
          inline-flex items-center gap-2 px-4 py-2 rounded-lg
          font-mono text-xs font-semibold uppercase tracking-[0.06em]
          transition-all duration-150 cursor-pointer
          border
          focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[var(--accent)]
        "
        style={{
          background: quickMode ? 'var(--accent)' : 'var(--bg-surface)',
          color: quickMode ? '#FFFFFF' : 'var(--text-secondary)',
          borderColor: quickMode ? 'var(--accent)' : 'var(--border)',
        }}
        aria-pressed={quickMode}
        title={quickMode ? 'Switch to Deep Mode (all sections)' : 'Switch to Quick Mode (key sections only)'}
      >
        {quickMode ? <Zap size={14} /> : <Moon size={14} />}
        {quickMode ? 'Quick Mode' : 'Deep Mode'}
      </button>
    </header>
  )
}

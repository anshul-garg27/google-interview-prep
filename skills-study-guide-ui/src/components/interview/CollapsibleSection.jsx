import { useRef, useCallback } from 'react'
import {
  Zap, Layers, Package, Send, Database, Globe, Bug,
  MessageCircle, Mic, AlertTriangle, GitBranch, HelpCircle,
  GitPullRequest, ChevronDown, Hash, Triangle, Heart,
  Volume2, FileText
} from 'lucide-react'

const ICONS = {
  Zap, Layers, Package, Send, Database, Globe, Bug,
  MessageCircle, Mic, AlertTriangle, GitBranch, HelpCircle,
  GitPullRequest, Hash, Triangle, Heart, Volume2, FileText
}

// Usage:
// <CollapsibleSection
//   id="system-design"
//   title="System Design"
//   icon="Layers"
//   isQuickMode={true}
//   expanded={true}
//   onToggle={() => toggle('system-design')}
//   quickModeActive={false}
//   index={1}
// >
//   <p>Section content here...</p>
// </CollapsibleSection>

export default function CollapsibleSection({
  id,
  title,
  icon,
  isQuickMode,
  expanded,
  onToggle,
  quickModeActive,
  children,
  index,
}) {
  const contentRef = useRef(null)
  const Icon = ICONS[icon] || HelpCircle

  const dimmed = quickModeActive && !isQuickMode
  const isExpanded = dimmed ? false : expanded

  const handleToggle = useCallback(() => {
    if (dimmed) return
    onToggle()
  }, [dimmed, onToggle])

  const handleKeyDown = useCallback((e) => {
    if (e.key === 'Enter' || e.key === ' ') {
      e.preventDefault()
      handleToggle()
    }
  }, [handleToggle])

  return (
    <section
      id={id}
      className="border-b border-[var(--border)]"
      style={{ scrollMarginTop: '5rem' }}
    >
      {/* Header */}
      <button
        onClick={handleToggle}
        onKeyDown={handleKeyDown}
        aria-expanded={isExpanded}
        aria-controls={`${id}-content`}
        disabled={dimmed}
        className={`w-full flex items-center gap-3 px-4 py-3.5 text-left transition-all duration-200
          ${dimmed
            ? 'opacity-50 cursor-default'
            : 'hover:bg-[var(--bg-hover)] cursor-pointer'
          }`}
      >
        {/* Section number */}
        <span className="font-mono text-[12px] font-bold text-[var(--accent)] shrink-0 w-5 text-right tabular-nums">
          {String(index).padStart(2, '0')}
        </span>

        {/* Icon */}
        <Icon
          size={16}
          className={`shrink-0 ${dimmed ? 'text-[var(--text-tertiary)]' : 'text-[var(--text-secondary)]'}`}
        />

        {/* Title */}
        <span className={`font-heading text-[15px] leading-snug flex-1 ${dimmed ? '' : 'font-semibold'}`}>
          {title}
        </span>

        {/* Quick badge */}
        {isQuickMode && (
          <span className="shrink-0 px-1.5 py-0.5 rounded text-[9px] font-mono font-bold uppercase tracking-[0.08em] bg-[var(--accent-subtle)] text-[var(--accent)]">
            Quick
          </span>
        )}

        {/* Dimmed indicator */}
        {dimmed && (
          <span className="shrink-0 font-mono text-[10px] text-[var(--text-tertiary)] italic">
            hidden in quick mode
          </span>
        )}

        {/* Chevron */}
        <ChevronDown
          size={16}
          className={`shrink-0 text-[var(--text-tertiary)] transition-transform duration-200
            ${isExpanded ? 'rotate-180' : 'rotate-0'}`}
        />
      </button>

      {/* Content - uses grid-rows trick for smooth expand/collapse */}
      <div
        id={`${id}-content`}
        role="region"
        aria-labelledby={id}
        className="grid transition-[grid-template-rows] duration-300 ease-in-out"
        style={{ gridTemplateRows: isExpanded ? '1fr' : '0fr' }}
      >
        <div ref={contentRef} className="overflow-hidden">
          <div className="px-2 pb-6 pt-2">
            {children}
          </div>
        </div>
      </div>
    </section>
  )
}

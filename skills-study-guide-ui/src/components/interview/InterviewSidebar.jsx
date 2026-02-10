import { useCallback } from 'react'
import {
  Zap, Layers, Package, Send, Database, Globe, Bug,
  MessageCircle, Mic, AlertTriangle, GitBranch, HelpCircle,
  GitPullRequest, Hash, Triangle, Heart, Volume2, FileText
} from 'lucide-react'

const ICONS = {
  Zap, Layers, Package, Send, Database, Globe, Bug,
  MessageCircle, Mic, AlertTriangle, GitBranch, HelpCircle,
  GitPullRequest, Hash, Triangle, Heart, Volume2, FileText
}

// Usage:
// <InterviewSidebar
//   sections={[
//     { id: 'system-design', title: 'System Design', icon: 'Layers', quickMode: true },
//     { id: 'debugging', title: 'Debugging', icon: 'Bug', quickMode: false },
//   ]}
//   activeSection="system-design"
//   quickMode={false}
// />

export default function InterviewSidebar({ sections, activeSection, quickMode }) {
  const handleClick = useCallback((id) => {
    const el = document.getElementById(id)
    if (el) {
      el.scrollIntoView({ behavior: 'smooth', block: 'start' })
    }
  }, [])

  return (
    <nav
      className="sticky top-20 py-8 pr-4 max-h-[calc(100vh-5rem)] overflow-y-auto"
      aria-label="Interview sections"
    >
      <p className="font-mono text-[10px] font-semibold uppercase tracking-[0.08em] text-[var(--text-tertiary)] mb-3">
        Sections
      </p>

      <ul className="space-y-0.5">
        {sections.map((section) => {
          const Icon = ICONS[section.icon] || HelpCircle
          const isActive = activeSection === section.id
          const dimmed = quickMode && !section.quickMode

          return (
            <li key={section.id}>
              <button
                onClick={() => handleClick(section.id)}
                title={section.title}
                className={`w-full text-left flex items-center gap-2 px-3 py-1.5 border-l-2 rounded-r transition-all font-mono text-[11px] leading-relaxed
                  ${isActive
                    ? 'border-[var(--accent)] text-[var(--accent)] font-semibold bg-[var(--accent-subtle)]'
                    : dimmed
                      ? 'border-transparent text-[var(--text-tertiary)] opacity-40'
                      : 'border-transparent text-[var(--text-tertiary)] hover:text-[var(--text-secondary)] hover:border-[var(--border-strong)]'
                  }`}
              >
                <Icon size={12} className="shrink-0" />
                <span className="truncate">
                  {section.title}
                </span>
              </button>
            </li>
          )
        })}
      </ul>
    </nav>
  )
}

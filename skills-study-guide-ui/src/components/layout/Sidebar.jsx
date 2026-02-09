import { useMemo } from 'react'
import { useLocation, useNavigate } from 'react-router-dom'
import { BookOpen, Database, Code2, Layers, BarChart3, GraduationCap } from 'lucide-react'
import guideIndex from '../../guide-index.json'
import categoryOrder from '../../categories.json'

const CATEGORY_ICONS = {
  'Infrastructure': Layers,
  'System Design': BarChart3,
  'Databases & Storage': Database,
  'Languages': Code2,
  'Data & Observability': BarChart3,
  'Interview Prep': GraduationCap,
}

export default function Sidebar({ open, onClose }) {
  const location = useLocation()
  const navigate = useNavigate()

  const grouped = useMemo(() => {
    const groups = {}
    for (const guide of guideIndex) {
      if (!groups[guide.category]) groups[guide.category] = []
      groups[guide.category].push(guide)
    }
    return categoryOrder.map(cat => ({
      name: cat,
      guides: groups[cat] || [],
    }))
  }, [])

  const currentSlug = location.pathname.startsWith('/guide/')
    ? location.pathname.replace('/guide/', '')
    : null

  return (
    <aside
      className={`fixed top-14 left-0 bottom-0 w-[260px] z-40 bg-[var(--bg-surface)] border-r border-[var(--border)] overflow-y-auto transition-transform duration-250 ease-out
        ${open ? 'translate-x-0' : '-translate-x-full'} lg:translate-x-0`}
    >
      <nav className="py-4 px-3">
        {grouped.map(group => {
          const Icon = CATEGORY_ICONS[group.name] || BookOpen
          return (
            <div key={group.name} className="mb-5">
              <div className="flex items-center gap-2 px-3 mb-2">
                <Icon size={12} className="text-[var(--text-tertiary)]" />
                <span className="font-mono text-[11px] font-semibold uppercase tracking-[0.06em] text-[var(--text-tertiary)]">
                  {group.name}
                </span>
              </div>

              {group.guides.map(guide => {
                const isActive = currentSlug === guide.slug
                return (
                  <button
                    key={guide.slug}
                    onClick={() => navigate(`/guide/${guide.slug}`)}
                    className={`w-full text-left px-3 py-2 rounded-md text-[13px] font-body flex items-start gap-2 transition-colors
                      ${isActive
                        ? 'sidebar-item-active font-semibold'
                        : 'text-[var(--text-secondary)] hover:bg-[var(--bg-hover)]'
                      }`}
                  >
                    <span className="font-mono text-[11px] text-[var(--text-tertiary)] mt-0.5 shrink-0 w-5">
                      {guide.num}
                    </span>
                    <span className="leading-snug">{guide.title}</span>
                  </button>
                )
              })}
            </div>
          )
        })}

        {/* Total stats */}
        <div className="mt-6 pt-4 border-t border-[var(--border)] px-3">
          <p className="font-mono text-[10px] uppercase tracking-widest text-[var(--text-tertiary)]">
            {guideIndex.length} guides &middot; {Math.round(guideIndex.reduce((sum, g) => sum + g.stats.wordCount, 0) / 1000)}K words
          </p>
        </div>
      </nav>
    </aside>
  )
}

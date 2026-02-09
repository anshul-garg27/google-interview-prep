import { useState, useMemo, useEffect, useRef } from 'react'
import { useNavigate } from 'react-router-dom'
import { Search, FileText, ArrowRight } from 'lucide-react'
import Fuse from 'fuse.js'
import guideIndex from '../../guide-index.json'

export default function SearchModal({ onClose }) {
  const [query, setQuery] = useState('')
  const [selectedIdx, setSelectedIdx] = useState(0)
  const inputRef = useRef(null)
  const navigate = useNavigate()

  // Build search index from guide sections
  const fuse = useMemo(() => {
    const searchData = guideIndex.flatMap(guide =>
      [
        {
          guideSlug: guide.slug,
          guideNum: guide.num,
          guideTitle: guide.title,
          sectionTitle: guide.title,
          sectionId: '',
          type: 'guide',
        },
        ...guide.sections.map(section => ({
          guideSlug: guide.slug,
          guideNum: guide.num,
          guideTitle: guide.title,
          sectionTitle: section.title,
          sectionId: section.id,
          type: 'section',
        })),
      ]
    )

    return new Fuse(searchData, {
      keys: [
        { name: 'guideTitle', weight: 3 },
        { name: 'sectionTitle', weight: 2 },
      ],
      threshold: 0.35,
      includeMatches: true,
    })
  }, [])

  const results = useMemo(() => {
    if (query.length < 2) {
      // Show all guides when no query
      return guideIndex.map(g => ({
        item: {
          guideSlug: g.slug,
          guideNum: g.num,
          guideTitle: g.title,
          sectionTitle: g.title,
          sectionId: '',
          type: 'guide',
        }
      }))
    }
    return fuse.search(query).slice(0, 12)
  }, [query, fuse])

  useEffect(() => {
    inputRef.current?.focus()
  }, [])

  useEffect(() => {
    setSelectedIdx(0)
  }, [query])

  const goToResult = (result) => {
    const item = result.item
    const path = `/guide/${item.guideSlug}${item.sectionId ? `#${item.sectionId}` : ''}`
    navigate(path)
    onClose()
  }

  const handleKeyDown = (e) => {
    if (e.key === 'ArrowDown') {
      e.preventDefault()
      setSelectedIdx(prev => Math.min(prev + 1, results.length - 1))
    } else if (e.key === 'ArrowUp') {
      e.preventDefault()
      setSelectedIdx(prev => Math.max(prev - 1, 0))
    } else if (e.key === 'Enter' && results[selectedIdx]) {
      goToResult(results[selectedIdx])
    } else if (e.key === 'Escape') {
      onClose()
    }
  }

  return (
    <div
      className="fixed inset-0 z-50 flex items-start justify-center pt-[15vh] px-4"
      role="dialog"
      aria-modal="true"
      aria-label="Search guides"
      onClick={onClose}
    >
      {/* Backdrop */}
      <div className="fixed inset-0 bg-black/50 backdrop-blur-sm" aria-hidden="true" />

      {/* Modal */}
      <div
        className="relative w-full max-w-[560px] bg-[var(--bg-surface)] rounded-xl border border-[var(--border)] shadow-2xl overflow-hidden animate-fade-in"
        onClick={e => e.stopPropagation()}
      >
        {/* Search input */}
        <div className="flex items-center gap-3 px-4 border-b border-[var(--border)]">
          <Search size={16} className="text-[var(--text-tertiary)] shrink-0" />
          <input
            ref={inputRef}
            type="text"
            value={query}
            onChange={e => setQuery(e.target.value)}
            onKeyDown={handleKeyDown}
            placeholder="Search guides and sections..."
            className="flex-1 py-3.5 text-base bg-transparent outline-none placeholder:text-[var(--text-tertiary)]"
          />
          <kbd className="text-[10px] font-mono px-1.5 py-0.5 rounded bg-[var(--bg-surface-alt)] border border-[var(--border)] text-[var(--text-tertiary)]">
            ESC
          </kbd>
        </div>

        {/* Results */}
        <div className="max-h-[50vh] overflow-y-auto py-2">
          {results.length === 0 ? (
            <div className="px-4 py-8 text-center">
              <p className="text-[var(--text-tertiary)] text-sm">No results found</p>
            </div>
          ) : (
            results.map((result, i) => {
              const item = result.item
              const isActive = i === selectedIdx
              const isSection = item.type === 'section'

              return (
                <button
                  key={`${item.guideSlug}-${item.sectionId || 'root'}-${i}`}
                  onClick={() => goToResult(result)}
                  onMouseEnter={() => setSelectedIdx(i)}
                  className={`w-full text-left px-4 py-2.5 flex items-center gap-3 transition-colors
                    ${isActive ? 'bg-[var(--accent-subtle)]' : 'hover:bg-[var(--bg-hover)]'}`}
                >
                  <div className={`shrink-0 w-8 h-8 rounded-lg flex items-center justify-center text-xs font-mono font-semibold
                    ${isActive ? 'bg-[var(--accent)] text-white' : 'bg-[var(--bg-surface-alt)] text-[var(--text-tertiary)]'}`}>
                    {isSection ? '#' : item.guideNum}
                  </div>

                  <div className="min-w-0 flex-1">
                    <p className="text-sm font-medium truncate">
                      {item.sectionTitle}
                    </p>
                    {isSection && (
                      <p className="text-[11px] text-[var(--text-tertiary)] truncate">
                        Guide {item.guideNum}: {item.guideTitle}
                      </p>
                    )}
                  </div>

                  {isActive && (
                    <ArrowRight size={14} className="text-[var(--accent)] shrink-0" />
                  )}
                </button>
              )
            })
          )}
        </div>

        {/* Footer hint */}
        <div className="px-4 py-2.5 border-t border-[var(--border)] flex items-center gap-4 text-[10px] font-mono text-[var(--text-tertiary)]">
          <span>↑↓ Navigate</span>
          <span>↵ Open</span>
          <span>ESC Close</span>
        </div>
      </div>
    </div>
  )
}

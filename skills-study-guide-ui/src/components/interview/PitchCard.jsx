import { useState, useId } from 'react'

// Usage:
// <PitchCard pitches={{ '30s': 'My 30-sec pitch...', '90s': 'My 90-sec pitch...', '2min': 'My 2-min pitch...' }} />

const TABS = [
  { key: '30s', label: '30 sec' },
  { key: '90s', label: '90 sec' },
  { key: '2min', label: '2 min' },
]

/**
 * Renders pitch text with **bold** markdown and paragraph splitting.
 * Splits on double newlines for paragraphs, converts single newlines to <br/>,
 * and wraps **text** in <strong> tags.
 */
function renderPitchText(text) {
  if (!text) return null
  return text.split('\n\n').map((paragraph, i) => (
    <p
      key={i}
      className="mb-3 last:mb-0"
      dangerouslySetInnerHTML={{
        __html: paragraph
          .replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>')
          .replace(/\n/g, '<br/>'),
      }}
    />
  ))
}

export default function PitchCard({ pitches }) {
  const [activeTab, setActiveTab] = useState('30s')
  const id = useId()
  const tablistId = `${id}-tablist`
  const panelId = `${id}-panel`

  if (!pitches) return null

  return (
    <div
      className="bg-[var(--bg-surface)] border border-[var(--border)] border-l-4 border-l-[var(--accent)] rounded-lg overflow-hidden"
      role="region"
      aria-label="Interview pitch"
    >
      {/* Header label */}
      <div className="px-5 pt-5 pb-0">
        <span className="font-mono text-[10px] font-semibold uppercase tracking-[0.12em] text-[var(--text-tertiary)]">
          Your Pitch
        </span>
      </div>

      {/* Tab buttons */}
      <div
        className="flex gap-2 px-5 pt-3 pb-4"
        role="tablist"
        id={tablistId}
        aria-label="Pitch duration"
      >
        {TABS.map((tab) => {
          const isActive = activeTab === tab.key
          const hasContent = !!pitches[tab.key]
          return (
            <button
              key={tab.key}
              role="tab"
              id={`${id}-tab-${tab.key}`}
              aria-selected={isActive}
              aria-controls={panelId}
              tabIndex={isActive ? 0 : -1}
              disabled={!hasContent}
              onClick={() => setActiveTab(tab.key)}
              onKeyDown={(e) => {
                const currentIndex = TABS.findIndex((t) => t.key === activeTab)
                let nextIndex = -1
                if (e.key === 'ArrowRight') {
                  nextIndex = (currentIndex + 1) % TABS.length
                } else if (e.key === 'ArrowLeft') {
                  nextIndex = (currentIndex - 1 + TABS.length) % TABS.length
                } else if (e.key === 'Home') {
                  nextIndex = 0
                } else if (e.key === 'End') {
                  nextIndex = TABS.length - 1
                }
                if (nextIndex >= 0 && pitches[TABS[nextIndex].key]) {
                  e.preventDefault()
                  setActiveTab(TABS[nextIndex].key)
                  document
                    .getElementById(`${id}-tab-${TABS[nextIndex].key}`)
                    ?.focus()
                }
              }}
              className={`px-4 py-1.5 rounded-full text-[13px] font-mono font-medium transition-all duration-200 focus-visible:outline-2 focus-visible:outline-[var(--accent)] focus-visible:outline-offset-2
                ${
                  isActive
                    ? 'bg-[var(--accent)] text-white shadow-sm'
                    : hasContent
                      ? 'border border-[var(--border)] text-[var(--text-secondary)] hover:bg-[var(--bg-hover)] hover:text-[var(--text-primary)]'
                      : 'border border-[var(--border)] text-[var(--text-tertiary)] opacity-40 cursor-not-allowed'
                }`}
            >
              {tab.label}
            </button>
          )
        })}
      </div>

      {/* Pitch content panel */}
      <div
        role="tabpanel"
        id={panelId}
        aria-labelledby={`${id}-tab-${activeTab}`}
        tabIndex={0}
        className="px-5 pb-5"
      >
        <blockquote className="border-l-[3px] border-l-[var(--accent)] bg-[var(--accent-subtle)] rounded-r-lg py-4 px-5">
          <div
            key={activeTab}
            className="font-body text-[15px] leading-relaxed text-[var(--text-primary)] animate-fade-in"
          >
            {renderPitchText(pitches[activeTab])}
          </div>
        </blockquote>
      </div>
    </div>
  )
}

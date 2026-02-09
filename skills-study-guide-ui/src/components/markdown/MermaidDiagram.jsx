import { useState, useEffect, useRef, memo } from 'react'
import { useTheme } from '../../context/ThemeContext'

const MermaidDiagram = memo(function MermaidDiagram({ chart }) {
  const { darkMode } = useTheme()
  const containerRef = useRef(null)
  const [svg, setSvg] = useState(null)
  const [error, setError] = useState(null)
  const [isVisible, setIsVisible] = useState(false)

  // Detect visibility with IntersectionObserver
  useEffect(() => {
    const el = containerRef.current
    if (!el) return

    const observer = new IntersectionObserver(
      ([entry]) => {
        if (entry.isIntersecting) {
          setIsVisible(true)
          observer.disconnect()
        }
      },
      { rootMargin: '200px' }
    )

    observer.observe(el)
    return () => observer.disconnect()
  }, [])

  // Lazy-load and render mermaid when visible
  useEffect(() => {
    if (!isVisible) return
    let cancelled = false

    setSvg(null) // Reset on theme change
    import('mermaid').then(({ default: mermaid }) => {
      if (cancelled) return

      mermaid.initialize({
        startOnLoad: false,
        theme: darkMode ? 'dark' : 'base',
        securityLevel: 'strict',
        fontFamily: "'JetBrains Mono', monospace",
        fontSize: 13,
        themeVariables: darkMode ? {
          primaryColor: '#3A3A3A',
          primaryBorderColor: '#F59E0B',
          primaryTextColor: '#E8E8E8',
          lineColor: '#666666',
          secondaryColor: '#252525',
          tertiaryColor: '#1E1E1E',
        } : {
          primaryColor: '#FEF3C7',
          primaryBorderColor: '#D97706',
          primaryTextColor: '#1A1A1A',
          lineColor: '#6B6B6B',
          secondaryColor: '#F5F4F0',
          tertiaryColor: '#FAFAF8',
        },
      })

      const id = `mermaid-${Math.random().toString(36).slice(2, 9)}`

      mermaid.render(id, chart.trim())
        .then(({ svg: renderedSvg }) => {
          if (!cancelled) setSvg(renderedSvg)
        })
        .catch(err => {
          if (!cancelled) setError(err.message || 'Diagram render failed')
        })
    }).catch(err => {
      if (!cancelled) setError('Failed to load diagram library')
    })

    return () => { cancelled = true }
  }, [isVisible, chart, darkMode])

  return (
    <div ref={containerRef} className="my-6">
      {!isVisible && (
        <div className="mermaid-container animate-pulse h-48 items-center justify-center">
          <span className="font-mono text-xs text-[var(--text-tertiary)]">Diagram</span>
        </div>
      )}

      {isVisible && !svg && !error && (
        <div className="mermaid-container h-48 items-center justify-center">
          <span className="font-mono text-xs text-[var(--text-tertiary)]">Rendering diagram...</span>
        </div>
      )}

      {svg && (
        <div className="my-6 rounded-xl border border-[var(--border)] bg-[var(--bg-surface)] overflow-hidden">
          <div className="flex items-center gap-2 px-4 py-2.5 border-b border-[var(--border)] bg-[var(--bg-surface-alt)]">
            <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" className="text-[var(--accent)]"><path d="M3 3h18v18H3z"/><path d="M12 3v18"/><path d="M3 12h18"/></svg>
            <span className="font-mono text-[10px] font-semibold uppercase tracking-[0.06em] text-[var(--accent)]">
              Diagram
            </span>
          </div>
          <div
            className="mermaid-container border-0 rounded-none m-0"
            dangerouslySetInnerHTML={{ __html: svg }}
          />
        </div>
      )}

      {error && (
        <div className="rounded-lg border border-red-200 dark:border-red-900/50 bg-red-50 dark:bg-red-950/20 p-4 my-6">
          <p className="font-mono text-[10px] uppercase tracking-wide text-red-500 mb-2">Diagram Error</p>
          <pre className="font-mono text-xs text-[var(--text-secondary)] overflow-x-auto whitespace-pre-wrap">{chart}</pre>
        </div>
      )}
    </div>
  )
})

export default MermaidDiagram

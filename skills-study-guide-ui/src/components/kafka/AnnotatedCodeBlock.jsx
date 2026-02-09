import { useState, useMemo, useCallback, useRef, useEffect } from 'react'
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter'
import { oneDark, oneLight } from 'react-syntax-highlighter/dist/esm/styles/prism'
import { useTheme } from '../../context/ThemeContext'
import { Copy, Check, ChevronRight, ChevronLeft, Eye, List, BookOpen, SkipForward, Code2 } from 'lucide-react'

const LANGUAGE_LABELS = {
  js: 'JavaScript', jsx: 'JSX', ts: 'TypeScript', tsx: 'TSX',
  java: 'Java', python: 'Python', py: 'Python', go: 'Go',
  sql: 'SQL', yaml: 'YAML', yml: 'YAML', json: 'JSON',
  bash: 'Bash', sh: 'Shell', xml: 'XML', html: 'HTML',
  css: 'CSS', protobuf: 'Protobuf', properties: 'Properties',
}

const MODES = [
  { id: 'read', label: 'Read', icon: BookOpen },
  { id: 'step', label: 'Step-Through', icon: ChevronRight },
  { id: 'all', label: 'Annotate All', icon: List },
]

/**
 * Normalize annotations to multi-line group format.
 * Supports two input shapes:
 *   Legacy:  { line: number, text: string }
 *   Group:   { lines: number[], title: string, desc: string }
 */
function normalizeAnnotations(raw) {
  if (!raw || raw.length === 0) return []
  // Detect format from first entry
  if (raw[0].lines) {
    // Already group format
    return raw.map((a, i) => ({
      lines: a.lines,
      title: a.title || `Step ${i + 1}`,
      desc: a.desc || a.text || '',
    }))
  }
  // Legacy single-line format — wrap each into a group
  return raw.map((a, i) => ({
    lines: [a.line],
    title: `Line ${a.line + 1}`,
    desc: a.text || '',
  }))
}

export default function AnnotatedCodeBlock({
  language = 'java',
  title,
  code,
  annotations: rawAnnotations = [],
  children,
}) {
  const { darkMode } = useTheme()
  const [copied, setCopied] = useState(false)
  const [mode, setMode] = useState(rawAnnotations.length > 0 ? 'step' : 'read')
  const [hoveredLine, setHoveredLine] = useState(null)
  const [stepIndex, setStepIndex] = useState(0)
  const containerRef = useRef(null)
  const codeRef = useRef(null)

  // Support both code prop and children
  const rawCode = (code || (typeof children === 'string' ? children : '')).trim()
  const lines = rawCode.split('\n')

  // Normalize annotations into group format
  const annotations = useMemo(() => normalizeAnnotations(rawAnnotations), [rawAnnotations])

  // Map: line number -> annotation group index (for all lines in all groups)
  const lineToGroupIdx = useMemo(() => {
    const m = new Map()
    annotations.forEach((a, idx) => {
      a.lines.forEach(l => m.set(l, idx))
    })
    return m
  }, [annotations])

  // Set of all annotated line numbers
  const allAnnotatedLines = useMemo(() => {
    const s = new Set()
    annotations.forEach(a => a.lines.forEach(l => s.add(l)))
    return s
  }, [annotations])

  const handleCopy = useCallback(async () => {
    try {
      await navigator.clipboard.writeText(rawCode)
    } catch {
      const ta = document.createElement('textarea')
      ta.value = rawCode
      document.body.appendChild(ta)
      ta.select()
      document.execCommand('copy')
      document.body.removeChild(ta)
    }
    setCopied(true)
    setTimeout(() => setCopied(false), 2000)
  }, [rawCode])

  // Current step annotation group
  const currentGroup = annotations[stepIndex] || null

  // Active highlighted line set for current mode
  const activeLines = useMemo(() => {
    if (mode === 'read') {
      if (hoveredLine !== null && allAnnotatedLines.has(hoveredLine)) {
        // Highlight the entire group that contains the hovered line
        const gIdx = lineToGroupIdx.get(hoveredLine)
        if (gIdx !== undefined) return new Set(annotations[gIdx].lines)
      }
      return new Set()
    }
    if (mode === 'step') return new Set(currentGroup?.lines || [])
    // 'all' mode
    return allAnnotatedLines
  }, [mode, hoveredLine, allAnnotatedLines, lineToGroupIdx, annotations, currentGroup])

  // Annotation text to show below a line (only at the LAST line of the group)
  const getInlineAnnotation = useCallback((lineNum) => {
    if (mode === 'read') {
      if (hoveredLine === null) return null
      const gIdx = lineToGroupIdx.get(hoveredLine)
      if (gIdx === undefined) return null
      const group = annotations[gIdx]
      const lastLine = group.lines[group.lines.length - 1]
      if (lineNum === lastLine) return group
      return null
    }
    if (mode === 'step') {
      if (!currentGroup) return null
      const lastLine = currentGroup.lines[currentGroup.lines.length - 1]
      if (lineNum === lastLine) return currentGroup
      return null
    }
    if (mode === 'all') {
      const gIdx = lineToGroupIdx.get(lineNum)
      if (gIdx === undefined) return null
      const group = annotations[gIdx]
      const lastLine = group.lines[group.lines.length - 1]
      if (lineNum === lastLine) return group
      return null
    }
    return null
  }, [mode, hoveredLine, lineToGroupIdx, annotations, currentGroup])

  // Reset step index when switching modes
  useEffect(() => {
    setStepIndex(0)
    setHoveredLine(null)
  }, [mode])

  // Scroll highlighted line into view in step mode
  useEffect(() => {
    if (mode === 'step' && currentGroup && codeRef.current) {
      const firstLine = currentGroup.lines[0]
      const el = codeRef.current.querySelector(`[data-line="${firstLine}"]`)
      if (el) el.scrollIntoView({ behavior: 'smooth', block: 'nearest' })
    }
  }, [mode, currentGroup])

  const canStepPrev = stepIndex > 0
  const canStepNext = stepIndex < annotations.length - 1
  const label = LANGUAGE_LABELS[language.toLowerCase()] || language

  return (
    <div className="my-6 rounded-xl border-2 border-blue-500/30 bg-[var(--bg-surface)] overflow-hidden not-prose" ref={containerRef}>
      {/* Header */}
      <div className="flex items-center justify-between px-4 py-2.5 bg-blue-50 dark:bg-blue-950/20 border-b border-[var(--border)]">
        <div className="flex items-center gap-2">
          <Code2 size={15} className="text-blue-500" />
          {title ? (
            <span className="font-heading font-semibold text-[14px]">{title}</span>
          ) : (
            <span className="font-heading font-semibold text-[14px]">Code Walkthrough</span>
          )}
          <span className="font-mono text-[10px] text-[var(--text-tertiary)] ml-1 px-1.5 py-0.5 rounded bg-[var(--bg-surface-alt)] border border-[var(--border)]">
            {label}
          </span>
        </div>
        <button
          onClick={handleCopy}
          className="flex items-center gap-1 px-2 py-1 rounded-md text-[10px] font-mono border border-[var(--border)] bg-[var(--bg-surface)] hover:bg-[var(--bg-hover)] transition-colors"
          title="Copy code"
        >
          {copied ? (
            <><Check size={11} className="text-emerald-500" /> Copied</>
          ) : (
            <><Copy size={11} /> Copy</>
          )}
        </button>
      </div>

      {/* Mode toggle + step controls */}
      {annotations.length > 0 && (
        <div className="flex items-center justify-between px-4 py-2 bg-[var(--bg-surface)] border-b border-[var(--border)]">
          <div className="flex items-center gap-1">
            {MODES.map(m => {
              const Icon = m.icon
              return (
                <button
                  key={m.id}
                  onClick={() => setMode(m.id)}
                  className={`flex items-center gap-1 px-2.5 py-1.5 rounded-lg font-mono text-[10px] transition-all ${
                    mode === m.id
                      ? 'bg-blue-500 text-white'
                      : 'border border-[var(--border)] hover:bg-[var(--bg-hover)] text-[var(--text-tertiary)]'
                  }`}
                >
                  <Icon size={10} />
                  {m.label}
                </button>
              )
            })}
          </div>

          {/* Step-through controls */}
          {mode === 'step' && annotations.length > 0 && (
            <div className="flex items-center gap-2">
              <div className="flex gap-1 mr-2">
                {annotations.map((_, i) => (
                  <button
                    key={i}
                    onClick={() => setStepIndex(i)}
                    className={`w-2 h-2 rounded-full transition-all ${
                      stepIndex === i ? 'bg-[var(--accent)] w-5' : 'bg-[var(--border)] hover:bg-[var(--text-tertiary)]'
                    }`}
                  />
                ))}
              </div>
              <button
                onClick={() => setStepIndex(i => Math.max(0, i - 1))}
                disabled={!canStepPrev}
                className="px-2 py-1 rounded-lg border border-[var(--border)] font-mono text-[11px] hover:bg-[var(--bg-hover)] transition-colors disabled:opacity-30 flex items-center gap-0.5"
              >
                <ChevronLeft size={11} /> Prev
              </button>
              <span className="font-mono text-[10px] text-[var(--text-tertiary)]">
                {stepIndex + 1}/{annotations.length}
              </span>
              <button
                onClick={() => setStepIndex(i => Math.min(annotations.length - 1, i + 1))}
                disabled={!canStepNext}
                className="px-2 py-1 rounded-lg bg-[var(--accent)] text-white font-mono text-[11px] hover:opacity-90 transition-colors disabled:opacity-30 flex items-center gap-0.5"
              >
                Next <SkipForward size={11} />
              </button>
            </div>
          )}
        </div>
      )}

      {/* Code lines */}
      <div className="relative" ref={codeRef}>
        {lines.map((line, i) => {
          const lineNum = i // annotations use 0-indexed line numbers for legacy, 1-indexed for group
          // Support both 0-indexed (legacy) and 1-indexed (group) by checking which is in the set
          const highlighted = activeLines.has(i) || activeLines.has(i + 1)
          const annotated = allAnnotatedLines.has(i) || allAnnotatedLines.has(i + 1)
          // Try both indexing schemes
          const inlineAnnotation = getInlineAnnotation(i) || getInlineAnnotation(i + 1)
          // But avoid double-showing: only show once
          const showAnnotation = inlineAnnotation && (
            inlineAnnotation.lines.includes(i) || inlineAnnotation.lines.includes(i + 1)
          )
          // Determine final highlight using whichever index is in the set
          const isHighlighted = activeLines.has(i)

          return (
            <div key={i} data-line={i}>
              <div
                className={`flex items-start font-mono text-[13px] leading-[1.7] transition-colors duration-200 ${
                  isHighlighted
                    ? 'bg-[var(--accent-subtle)]'
                    : mode === 'read' && hoveredLine === i
                      ? 'bg-[var(--bg-hover)]'
                      : ''
                } ${mode === 'read' && annotated ? 'cursor-pointer' : ''}`}
                onMouseEnter={mode === 'read' ? () => setHoveredLine(i) : undefined}
                onMouseLeave={mode === 'read' ? () => setHoveredLine(null) : undefined}
              >
                {/* Annotation indicator bar */}
                <div className={`w-[3px] shrink-0 transition-colors duration-150 ${
                  allAnnotatedLines.has(i) ? 'bg-[var(--accent)]' : 'bg-transparent'
                }`} />

                {/* Line number */}
                <span className="w-8 shrink-0 text-right mr-4 ml-2 select-none text-[11px] leading-[1.85] text-[var(--text-tertiary)]">
                  {i + 1}
                </span>

                {/* Code content */}
                <pre className="flex-1 overflow-x-auto" style={{ margin: 0, padding: 0, background: 'transparent' }}>
                  <SyntaxHighlighter
                    language={language.toLowerCase()}
                    style={darkMode ? oneDark : oneLight}
                    PreTag="span"
                    customStyle={{
                      margin: 0,
                      padding: 0,
                      background: 'transparent',
                      fontSize: '13px',
                      lineHeight: '1.7',
                      display: 'inline',
                    }}
                    codeTagProps={{
                      style: { fontFamily: "'JetBrains Mono', monospace" }
                    }}
                  >
                    {line || ' '}
                  </SyntaxHighlighter>
                </pre>
              </div>

              {/* Inline annotation — shown below the last line of the group */}
              {showAnnotation && (mode === 'step' || mode === 'all' || mode === 'read') && (
                <div className="flex items-start gap-2 px-4 py-2.5 bg-[var(--accent-subtle)] border-l-[3px] border-l-[var(--accent)] animate-fade-in">
                  <Eye size={12} className="text-[var(--accent)] shrink-0 mt-0.5" />
                  <div>
                    {inlineAnnotation.title && (
                      <div className="font-mono text-[10px] text-[var(--accent)] font-semibold uppercase tracking-wider mb-0.5">
                        {inlineAnnotation.title}
                      </div>
                    )}
                    <span className="text-[12px] leading-relaxed text-[var(--text-primary)] font-body">
                      {inlineAnnotation.desc}
                    </span>
                  </div>
                </div>
              )}
            </div>
          )
        })}
      </div>

      {/* Footer hint */}
      {annotations.length > 0 && mode === 'read' && (
        <div className="px-4 py-2 bg-[var(--bg-surface-alt)] border-t border-[var(--border)]">
          <span className="font-mono text-[10px] text-[var(--text-tertiary)]">
            Hover annotated lines (amber bar) to see explanations. {annotations.length} annotation{annotations.length !== 1 ? 's' : ''} available.
          </span>
        </div>
      )}
      {annotations.length > 0 && mode === 'step' && currentGroup && (
        <div className="px-4 py-2 bg-[var(--bg-surface-alt)] border-t border-[var(--border)]">
          <span className="font-mono text-[10px] text-[var(--text-tertiary)]">
            Step through each annotation with the Prev/Next buttons. Lines {currentGroup.lines.map(l => l + 1).join(', ')} highlighted.
          </span>
        </div>
      )}
    </div>
  )
}

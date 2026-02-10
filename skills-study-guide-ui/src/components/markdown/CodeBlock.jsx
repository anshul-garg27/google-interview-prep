import { useState, memo } from 'react'
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter'
import { oneDark, oneLight } from 'react-syntax-highlighter/dist/esm/styles/prism'
import { Copy, Check, BoxSelect, ChevronRight, Code2 } from 'lucide-react'
import { useTheme } from '../../context/ThemeContext'
import MermaidDiagram from './MermaidDiagram'

const COLLAPSE_THRESHOLD = 5 // Lines above which code blocks collapse

const LANGUAGE_LABELS = {
  js: 'JavaScript', jsx: 'JSX', ts: 'TypeScript', tsx: 'TSX',
  java: 'Java', python: 'Python', py: 'Python', go: 'Go',
  sql: 'SQL', yaml: 'YAML', yml: 'YAML', json: 'JSON',
  bash: 'Bash', sh: 'Shell', xml: 'XML', html: 'HTML',
  css: 'CSS', protobuf: 'Protobuf', proto: 'Protobuf',
  properties: 'Properties', ini: 'INI', dockerfile: 'Dockerfile',
  ruby: 'Ruby', http: 'HTTP', promql: 'PromQL', text: 'Text',
}

// Detect if a code block looks like an ASCII diagram
function isAsciiDiagram(code) {
  const lines = code.split('\n')
  if (lines.length < 3) return false

  const diagramChars = /[+\-|>←→↓↑▸▾┌┐└┘│─╔╗╚╝═║\^]/
  const boxChars = /[+\-|]/
  const arrowChars = /[→←↑↓>]|-->|--|==/

  let diagramLines = 0
  for (const line of lines) {
    if (diagramChars.test(line) || arrowChars.test(line)) diagramLines++
  }

  // If >30% of lines look like diagram, treat as diagram
  return diagramLines / lines.length > 0.3
}

// Beautiful ASCII diagram renderer
function AsciiDiagram({ code }) {
  return (
    <div className="my-6 not-prose">
      <div className="rounded-xl border border-[var(--border)] bg-[var(--bg-surface)] overflow-hidden">
        {/* Header */}
        <div className="flex items-center gap-2 px-4 py-2.5 border-b border-[var(--border)] bg-[var(--bg-surface-alt)]">
          <BoxSelect size={13} className="text-[var(--accent)]" />
          <span className="font-mono text-[10px] font-semibold uppercase tracking-[0.06em] text-[var(--accent)]">
            Diagram
          </span>
        </div>
        {/* Diagram content */}
        <div className="p-5 overflow-x-auto">
          <pre className="font-mono text-[13px] leading-[1.6] text-[var(--text-primary)] whitespace-pre" style={{ margin: 0 }}>
            {code}
          </pre>
        </div>
      </div>
    </div>
  )
}

const CodeBlock = memo(function CodeBlock({ inline, className, children, node, ...props }) {
  const [copied, setCopied] = useState(false)
  const { darkMode } = useTheme()

  const match = /language-(\w+)/.exec(className || '')
  const language = match?.[1] || ''
  const code = String(children).replace(/\n$/, '')

  // Inline code — detect by: explicit inline prop, OR no className + single line + short text
  const isInline = inline || (!className && !code.includes('\n') && code.length < 200)

  if (isInline) {
    return (
      <code
        className="font-mono text-[0.875em] font-medium px-1.5 py-0.5 rounded bg-[var(--bg-surface-alt)] border border-[var(--border)]"
        {...props}
      >
        {children}
      </code>
    )
  }

  // Mermaid diagrams
  if (language === 'mermaid') {
    return <MermaidDiagram chart={code} />
  }

  // No language specified — check if it's ASCII art diagram
  if (!language || language === 'text') {
    if (isAsciiDiagram(code)) {
      return <AsciiDiagram code={code} />
    }
    // Multi-line text without language — render as plain text block
    if (!language) {
      return <AsciiDiagram code={code} />
    }
  }

  const label = LANGUAGE_LABELS[language] || language
  const lineCount = code.split('\n').length
  const isCollapsible = lineCount > COLLAPSE_THRESHOLD

  const handleCopy = async () => {
    try {
      await navigator.clipboard.writeText(code)
    } catch {
      // Fallback
      const ta = document.createElement('textarea')
      ta.value = code
      document.body.appendChild(ta)
      ta.select()
      document.execCommand('copy')
      document.body.removeChild(ta)
    }
    setCopied(true)
    setTimeout(() => setCopied(false), 2000)
  }

  // Collapsible wrapper for long code blocks
  if (isCollapsible) {
    return <CollapsibleCode
      code={code}
      language={language}
      label={label}
      lineCount={lineCount}
      darkMode={darkMode}
      handleCopy={handleCopy}
      copied={copied}
    />
  }

  return (
    <div className="code-block-wrapper not-prose">
      <div className="code-block-header">
        <span>{label}</span>
        <button
          onClick={handleCopy}
          className="flex items-center gap-1 text-[10px] hover:text-[var(--text-primary)] transition-colors"
          aria-label="Copy code"
        >
          {copied ? (
            <>
              <Check size={12} className="text-emerald-500" />
              <span className="text-emerald-500">Copied</span>
            </>
          ) : (
            <>
              <Copy size={12} />
              <span>Copy</span>
            </>
          )}
        </button>
      </div>

      <SyntaxHighlighter
        style={darkMode ? oneDark : oneLight}
        language={language}
        PreTag="div"
        customStyle={{
          margin: 0,
          padding: '1rem 1.25rem',
          background: 'var(--code-bg)',
          fontSize: '13px',
          lineHeight: '1.7',
          borderRadius: 0,
        }}
        codeTagProps={{
          style: {
            fontFamily: "'JetBrains Mono', monospace",
          }
        }}
      >
        {code}
      </SyntaxHighlighter>
    </div>
  )
})

// Collapsible code block for long code snippets
function CollapsibleCode({ code, language, label, lineCount, darkMode, handleCopy, copied }) {
  const [expanded, setExpanded] = useState(false)
  // Show first 3 lines as preview
  const preview = code.split('\n').slice(0, 3).join('\n')

  return (
    <div className="code-block-wrapper not-prose">
      {/* Header — always visible, clickable to expand */}
      <div
        className="code-block-header cursor-pointer select-none"
        onClick={() => setExpanded(prev => !prev)}
        role="button"
        aria-expanded={expanded}
        tabIndex={0}
        onKeyDown={(e) => { if (e.key === 'Enter' || e.key === ' ') { e.preventDefault(); setExpanded(prev => !prev) } }}
      >
        <span className="flex items-center gap-2">
          <ChevronRight
            size={13}
            className={`transition-transform duration-200 ${expanded ? 'rotate-90' : ''}`}
          />
          <Code2 size={12} />
          <span>{label}</span>
          <span className="text-[var(--text-tertiary)] font-normal">{lineCount} lines</span>
        </span>
        <div className="flex items-center gap-3">
          {expanded && (
            <button
              onClick={(e) => { e.stopPropagation(); handleCopy() }}
              className="flex items-center gap-1 text-[10px] hover:text-[var(--text-primary)] transition-colors"
              aria-label="Copy code"
            >
              {copied ? (
                <><Check size={12} className="text-emerald-500" /><span className="text-emerald-500">Copied</span></>
              ) : (
                <><Copy size={12} /><span>Copy</span></>
              )}
            </button>
          )}
          <span className="text-[10px] text-[var(--text-tertiary)]">
            {expanded ? 'collapse' : 'expand'}
          </span>
        </div>
      </div>

      {/* Code content — collapsed or expanded */}
      <div
        className="grid transition-[grid-template-rows] duration-300 ease-in-out"
        style={{ gridTemplateRows: expanded ? '1fr' : '0fr' }}
      >
        <div className="overflow-hidden">
          <SyntaxHighlighter
            style={darkMode ? oneDark : oneLight}
            language={language}
            PreTag="div"
            customStyle={{
              margin: 0,
              padding: '1rem 1.25rem',
              background: 'var(--code-bg)',
              fontSize: '13px',
              lineHeight: '1.7',
              borderRadius: 0,
            }}
            codeTagProps={{
              style: { fontFamily: "'JetBrains Mono', monospace" }
            }}
          >
            {code}
          </SyntaxHighlighter>
        </div>
      </div>

      {/* Preview when collapsed — show first 3 lines faded */}
      {!expanded && (
        <div className="relative overflow-hidden" style={{ maxHeight: '4.5em' }}>
          <pre
            className="font-mono text-[13px] leading-[1.5] px-5 py-2 text-[var(--text-tertiary)]"
            style={{ margin: 0, background: 'var(--code-bg)' }}
          >
            {preview}
          </pre>
          <div
            className="absolute inset-x-0 bottom-0 h-8"
            style={{ background: 'linear-gradient(transparent, var(--code-bg))' }}
          />
        </div>
      )}
    </div>
  )
}

export default CodeBlock

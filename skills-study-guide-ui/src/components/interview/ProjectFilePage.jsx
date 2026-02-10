import { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { ArrowLeft } from 'lucide-react'
import MarkdownRenderer from '../markdown/MarkdownRenderer'

const cache = new Map()

export default function ProjectFilePage() {
  const { projectSlug, filename } = useParams()
  const navigate = useNavigate()
  const [content, setContent] = useState(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)

  const displayName = filename
    .replace(/^\d+-/, '')
    .replace(/[-_]/g, ' ')
    .replace(/\b\w/g, c => c.toUpperCase())

  useEffect(() => {
    let cancelled = false
    const key = `${projectSlug}/${filename}`

    if (cache.has(key)) {
      setContent(cache.get(key))
      setLoading(false)
      return
    }

    setLoading(true)
    setError(null)

    fetch(`${import.meta.env.BASE_URL}interview/${projectSlug}/${filename}.md`)
      .then(res => {
        if (!res.ok) throw new Error('File not found')
        return res.text()
      })
      .then(text => {
        if (cancelled) return
        cache.set(key, text)
        setContent(text)
        setLoading(false)
      })
      .catch(err => {
        if (cancelled) return
        setError(err.message)
        setLoading(false)
      })

    return () => { cancelled = true }
  }, [projectSlug, filename])

  useEffect(() => {
    document.title = `${displayName} — ${projectSlug}`
    window.scrollTo(0, 0)
    return () => { document.title = 'Study Guides — Anshul Garg' }
  }, [displayName, projectSlug])

  return (
    <div className="animate-fade-in">
      <div className="max-w-[960px] mx-auto px-5 md:px-10 py-8">
        {/* Back link */}
        <button
          onClick={() => navigate(`/interview/${projectSlug}`)}
          className="flex items-center gap-2 text-[var(--text-tertiary)] hover:text-[var(--accent)] text-sm font-mono mb-4 transition-colors"
        >
          <ArrowLeft size={14} />
          Back to {projectSlug.replace(/-/g, ' ')}
        </button>

        {/* Header */}
        <div className="mb-8 pb-6 border-b border-[var(--border)]">
          <span className="font-mono text-[10px] font-semibold uppercase tracking-[0.06em] text-[var(--accent)]">
            {projectSlug.replace(/-/g, ' ')}
          </span>
          <h1 className="font-heading font-extrabold text-2xl md:text-3xl mt-1">
            {displayName}
          </h1>
          <p className="font-mono text-[11px] text-[var(--text-tertiary)] mt-2">
            {filename}.md
          </p>
        </div>

        {/* Content */}
        {loading && (
          <div className="py-16 text-center text-[var(--text-tertiary)]">Loading...</div>
        )}
        {error && (
          <div className="py-16 text-center">
            <p className="text-[var(--text-secondary)] mb-4">{error}</p>
            <button
              onClick={() => navigate(`/interview/${projectSlug}`)}
              className="text-[var(--accent)] hover:underline text-sm"
            >
              Go back to project
            </button>
          </div>
        )}
        {content && (
          <div className="guide-content">
            <MarkdownRenderer content={content} />
          </div>
        )}
      </div>
    </div>
  )
}

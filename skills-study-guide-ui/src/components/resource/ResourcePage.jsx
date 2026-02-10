import { useState, useEffect, useMemo } from 'react'
import { useParams, Link } from 'react-router-dom'
import MarkdownRenderer from '../markdown/MarkdownRenderer'
import { getResource } from '../../data/standalone-resources'

const cache = new Map()

export default function ResourcePage() {
  const { filename } = useParams()
  const [content, setContent] = useState(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)

  const meta = useMemo(() => getResource(`${filename}.md`), [filename])

  // Fetch markdown from /resources/{filename}.md
  useEffect(() => {
    if (!filename) return
    let cancelled = false
    const key = filename

    if (cache.has(key)) {
      setContent(cache.get(key))
      setLoading(false)
      setError(null)
      return
    }

    setLoading(true)
    setError(null)

    fetch(`${import.meta.env.BASE_URL}resources/${filename}.md`)
      .then(res => {
        if (!res.ok) throw new Error(`Failed to load ${filename}.md`)
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
  }, [filename])

  // Set document title
  useEffect(() => {
    const title = meta?.label || formatFilename(filename)
    document.title = `${title} — Study Guide`
    return () => { document.title = 'Study Guides — Anshul Garg' }
  }, [filename, meta])

  // Scroll to top on filename change
  useEffect(() => { window.scrollTo(0, 0) }, [filename])

  const displayTitle = meta?.label || formatFilename(filename)
  const categoryName = meta?.category || null

  if (loading) {
    return (
      <div className="animate-fade-in">
        <ResourceHeader title={displayTitle} category={categoryName} />
        <div className="max-w-[960px] mx-auto px-5 py-16 text-center text-[var(--text-tertiary)]">
          Loading resource...
        </div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="animate-fade-in">
        <ResourceHeader title={displayTitle} category={categoryName} />
        <div className="max-w-[960px] mx-auto px-5 py-16 text-center">
          <p className="text-[var(--text-secondary)] mb-4">{error}</p>
          <Link
            to="/"
            className="text-[var(--accent)] hover:text-[var(--accent-hover)] underline underline-offset-4"
          >
            Back to Home
          </Link>
        </div>
      </div>
    )
  }

  return (
    <div className="animate-fade-in">
      <ResourceHeader title={displayTitle} category={categoryName} />
      <article className="max-w-[960px] mx-auto px-5 md:px-10 py-8">
        <MarkdownRenderer content={content} />
      </article>
    </div>
  )
}

function ResourceHeader({ title, category }) {
  return (
    <header className="border-b border-[var(--border)] bg-[var(--bg-surface)]">
      <div className="max-w-[960px] mx-auto px-5 md:px-10 py-6">
        {category && (
          <span className="inline-block text-xs font-mono uppercase tracking-wider text-[var(--accent)] mb-2">
            {category}
          </span>
        )}
        <h1 className="font-heading text-2xl md:text-3xl font-semibold text-[var(--text-primary)] leading-tight">
          {title}
        </h1>
      </div>
    </header>
  )
}

/**
 * Fallback title formatter when the file isn't in the resource lookup.
 * Removes .md, replaces _ and - with spaces, strips leading numbers, title-cases.
 */
function formatFilename(filename) {
  if (!filename) return 'Resource'
  return filename
    .replace(/\.md$/i, '')
    .replace(/^(\d{1,2}[-_])/, '')     // strip leading "01-" or "06_"
    .replace(/[_-]/g, ' ')             // underscores and hyphens to spaces
    .replace(/\b\w/g, c => c.toUpperCase()) // title case
    .trim()
}

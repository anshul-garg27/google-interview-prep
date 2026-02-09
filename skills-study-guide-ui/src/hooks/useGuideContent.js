import { useState, useEffect } from 'react'

const cache = new Map()

export function useGuideContent(slug) {
  const [content, setContent] = useState(cache.get(slug) || null)
  const [loading, setLoading] = useState(!cache.has(slug))
  const [error, setError] = useState(null)

  useEffect(() => {
    let cancelled = false

    if (cache.has(slug)) {
      setContent(cache.get(slug))
      setLoading(false)
      setError(null)
      return
    }

    const controller = new AbortController()
    setLoading(true)
    setError(null)

    fetch(`${import.meta.env.BASE_URL}guides/${slug}.md`, { signal: controller.signal })
      .then(res => {
        if (!res.ok) throw new Error(`Guide not found: ${slug}`)
        return res.text()
      })
      .then(text => {
        if (cancelled) return
        cache.set(slug, text)
        setContent(text)
        setLoading(false)
      })
      .catch(err => {
        if (cancelled) return
        if (err.name !== 'AbortError') {
          setError(err.message)
          setLoading(false)
        }
      })

    return () => {
      cancelled = true
      controller.abort()
    }
  }, [slug])

  return { content, loading, error }
}

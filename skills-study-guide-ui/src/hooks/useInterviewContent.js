import { useState, useEffect, useCallback } from 'react'

const cache = new Map()

export function useInterviewContent(projectSlug, sections) {
  const [contents, setContents] = useState({})
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)

  useEffect(() => {
    if (!projectSlug || !sections?.length) return
    let cancelled = false

    // Check if all sections are cached
    const allCached = sections.every(s => cache.has(`${projectSlug}/${s.file}`))
    if (allCached) {
      const cached = {}
      for (const s of sections) {
        cached[s.id] = cache.get(`${projectSlug}/${s.file}`)
      }
      setContents(cached)
      setLoading(false)
      return
    }

    setLoading(true)
    setError(null)

    const promises = sections.map(section => {
      const key = `${projectSlug}/${section.file}`
      if (cache.has(key)) {
        return Promise.resolve({ id: section.id, content: cache.get(key) })
      }
      return fetch(`${import.meta.env.BASE_URL}interview/${projectSlug}/${section.file}`)
        .then(res => {
          if (!res.ok) throw new Error(`Failed to load ${section.file}`)
          return res.text()
        })
        .then(text => {
          cache.set(key, text)
          return { id: section.id, content: text }
        })
    })

    Promise.all(promises)
      .then(results => {
        if (cancelled) return
        const map = {}
        for (const r of results) {
          map[r.id] = r.content
        }
        setContents(map)
        setLoading(false)
      })
      .catch(err => {
        if (cancelled) return
        setError(err.message)
        setLoading(false)
      })

    return () => { cancelled = true }
  }, [projectSlug, sections])

  const getContent = useCallback((sectionId) => contents[sectionId] || null, [contents])

  return { contents, getContent, loading, error }
}

// Fetch a single section lazily (for CollapsibleSection)
export function useLazySection(projectSlug, file) {
  const [content, setContent] = useState(null)
  const [loading, setLoading] = useState(false)

  const load = useCallback(() => {
    if (content) return
    const key = `${projectSlug}/${file}`
    if (cache.has(key)) {
      setContent(cache.get(key))
      return
    }

    setLoading(true)
    fetch(`${import.meta.env.BASE_URL}interview/${projectSlug}/${file}`)
      .then(res => res.ok ? res.text() : Promise.reject(new Error('Failed')))
      .then(text => {
        cache.set(key, text)
        setContent(text)
        setLoading(false)
      })
      .catch(() => setLoading(false))
  }, [projectSlug, file, content])

  return { content, loading, load }
}

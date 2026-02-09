import { useState, useEffect, useMemo } from 'react'
import { extractHeadings } from '../utils/extractHeadings'

export function useScrollSpy(content) {
  const [activeHeading, setActiveHeading] = useState('')
  const headings = useMemo(() => extractHeadings(content), [content])

  useEffect(() => {
    if (!headings.length) return

    let mounted = true

    const observer = new IntersectionObserver(
      (entries) => {
        if (!mounted) return
        const visible = entries
          .filter(e => e.isIntersecting)
          .sort((a, b) => a.boundingClientRect.top - b.boundingClientRect.top)

        if (visible.length > 0) {
          setActiveHeading(visible[0].target.id)
        }
      },
      { rootMargin: '-80px 0px -60% 0px', threshold: 0 }
    )

    // Delay to ensure DOM headings exist after markdown render
    const timer = setTimeout(() => {
      if (!mounted) return
      headings.forEach(h => {
        const el = document.getElementById(h.id)
        if (el) observer.observe(el)
      })
    }, 150)

    return () => {
      mounted = false
      clearTimeout(timer)
      observer.disconnect()
    }
  }, [headings])

  return { activeHeading, headings }
}

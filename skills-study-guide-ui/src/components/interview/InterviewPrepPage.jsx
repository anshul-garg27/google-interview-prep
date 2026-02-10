import { useState, useEffect, useRef, useCallback, useMemo, memo } from 'react'
import { useParams } from 'react-router-dom'
import { getProject, splitBySections } from '../../data/interview-projects'
import InterviewHeader from './InterviewHeader'
import CollapsibleSection from './CollapsibleSection'
import PitchCard from './PitchCard'
import InterviewSidebar from './InterviewSidebar'
import MarkdownRenderer from '../markdown/MarkdownRenderer'

const contentCache = new Map()

export default function InterviewPrepPage() {
  const { projectSlug } = useParams()
  const project = getProject(projectSlug)

  const [masterContent, setMasterContent] = useState(null)
  const [loading, setLoading] = useState(true)
  const [quickMode, setQuickMode] = useState(false)
  const [expandedSections, setExpandedSections] = useState(new Set())
  const [mountedSections, setMountedSections] = useState(new Set()) // Track which sections have been opened
  const [activeSection, setActiveSection] = useState(null)

  const sectionRefs = useRef({})
  const observerRef = useRef(null)

  // Fetch the single master document
  useEffect(() => {
    if (!project) return
    let cancelled = false
    const key = `${project.slug}/${project.masterFile}`

    if (contentCache.has(key)) {
      setMasterContent(contentCache.get(key))
      setLoading(false)
      return
    }

    setLoading(true)
    fetch(`${import.meta.env.BASE_URL}interview/${project.slug}/${project.masterFile}`)
      .then(res => res.ok ? res.text() : Promise.reject(new Error('Failed to load')))
      .then(text => {
        if (cancelled) return
        contentCache.set(key, text)
        setMasterContent(text)
        setLoading(false)
      })
      .catch(() => {
        if (cancelled) return
        setLoading(false)
      })

    return () => { cancelled = true }
  }, [project])

  // Split content into sections by H2 headings
  const sections = useMemo(() => {
    if (!masterContent || !project) return []
    const parsed = splitBySections(masterContent)
    return parsed.map(s => ({
      ...s,
      icon: project.sectionMeta[s.title]?.icon || 'FileText',
      quickMode: project.sectionMeta[s.title]?.quickMode ?? false,
    }))
  }, [masterContent, project])

  // Start with only first 3 quick-mode sections expanded (fast initial load)
  useEffect(() => {
    if (sections.length > 0 && expandedSections.size === 0) {
      const quickSections = sections.filter(s => s.quickMode && s.title !== 'Your Pitches')
      const initial = quickSections.slice(0, 3).map(s => s.id)
      setExpandedSections(new Set(initial))
      setMountedSections(new Set(initial))
    }
  }, [sections])

  // When a section is expanded, mark it as mounted (so it stays rendered even when collapsed)
  useEffect(() => {
    setMountedSections(prev => {
      const next = new Set(prev)
      for (const id of expandedSections) next.add(id)
      return next.size !== prev.size ? next : prev
    })
  }, [expandedSections])

  // Set document title
  useEffect(() => {
    if (project) document.title = `${project.title} — Interview Prep`
    return () => { document.title = 'Study Guides — Anshul Garg' }
  }, [project])

  // Scroll to top on project change
  useEffect(() => { window.scrollTo(0, 0) }, [projectSlug])

  // IntersectionObserver for active section tracking
  useEffect(() => {
    if (!sections.length) return

    observerRef.current = new IntersectionObserver(
      (entries) => {
        for (const entry of entries) {
          if (entry.isIntersecting) setActiveSection(entry.target.id)
        }
      },
      { rootMargin: '-80px 0px -60% 0px', threshold: 0 }
    )

    const timer = setTimeout(() => {
      for (const id of Object.keys(sectionRefs.current)) {
        const el = sectionRefs.current[id]
        if (el) observerRef.current.observe(el)
      }
    }, 200)

    return () => {
      clearTimeout(timer)
      if (observerRef.current) observerRef.current.disconnect()
    }
  }, [sections])

  const registerSectionRef = useCallback((id, el) => {
    sectionRefs.current[id] = el
    if (el && observerRef.current) observerRef.current.observe(el)
  }, [])

  const handleToggleMode = useCallback(() => {
    setQuickMode(prev => {
      const next = !prev
      if (next) {
        setExpandedSections(new Set(
          sections.filter(s => s.quickMode && s.title !== 'Your Pitches').map(s => s.id)
        ))
      } else {
        setExpandedSections(new Set(
          sections.filter(s => s.title !== 'Your Pitches').map(s => s.id)
        ))
      }
      return next
    })
  }, [sections])

  const toggleSection = useCallback((sectionId) => {
    setExpandedSections(prev => {
      const next = new Set(prev)
      if (next.has(sectionId)) next.delete(sectionId)
      else next.add(sectionId)
      return next
    })
  }, [])

  if (!project) {
    return (
      <div className="max-w-prose mx-auto px-6 py-16 text-center">
        <h2 className="font-heading font-semibold text-2xl mb-4">Project not found</h2>
        <p className="text-[var(--text-secondary)]">No interview prep project found for "{projectSlug}".</p>
      </div>
    )
  }

  if (loading) {
    return (
      <div className="animate-fade-in">
        <InterviewHeader project={project} quickMode={quickMode} onToggleMode={handleToggleMode} />
        <div className="max-w-[960px] mx-auto px-5 py-16 text-center text-[var(--text-tertiary)]">
          Loading interview prep...
        </div>
      </div>
    )
  }

  const filteredSections = sections.filter(s => s.title !== 'Your Pitches')
  const sidebarSections = filteredSections.map(s => ({
    id: s.id, title: s.title, icon: s.icon, quickMode: s.quickMode,
  }))

  return (
    <div className="animate-fade-in">
      <InterviewHeader
        project={project}
        quickMode={quickMode}
        onToggleMode={handleToggleMode}
      />

      <div className="flex justify-center">
        <article className="flex-1 min-w-0 max-w-[960px] px-5 md:px-10 py-8">
          {project.pitches && <PitchCard pitches={project.pitches} />}

          {filteredSections.map((section) => (
            <div
              key={section.id}
              id={section.id}
              ref={(el) => registerSectionRef(section.id, el)}
            >
              <CollapsibleSection
                id={section.id}
                title={section.title}
                icon={section.icon}
                isQuickMode={section.quickMode}
                expanded={expandedSections.has(section.id)}
                onToggle={() => toggleSection(section.id)}
                quickModeActive={quickMode}
                index={section.index}
              >
                {/* LAZY: Only render markdown once section has been expanded at least once */}
                {mountedSections.has(section.id) ? (
                  <MemoizedMarkdown content={section.content} />
                ) : (
                  <div className="py-6 text-center text-[var(--text-tertiary)] text-sm">
                    Click to expand
                  </div>
                )}
              </CollapsibleSection>
            </div>
          ))}
        </article>

        <div className="hidden xl:block w-[220px] shrink-0">
          <InterviewSidebar
            sections={sidebarSections}
            activeSection={activeSection}
            quickMode={quickMode}
          />
        </div>
      </div>
    </div>
  )
}

// Memoize MarkdownRenderer so it doesn't re-render on parent state changes
const MemoizedMarkdown = memo(function MemoizedMarkdown({ content }) {
  return <MarkdownRenderer content={content} />
})

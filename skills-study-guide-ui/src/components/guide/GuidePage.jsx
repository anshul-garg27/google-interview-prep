import { useEffect } from 'react'
import { useParams } from 'react-router-dom'
import { useGuideContent } from '../../hooks/useGuideContent'
import { useScrollSpy } from '../../hooks/useScrollSpy'
import { useReadingProgress } from '../../hooks/useReadingProgress'
import { useProgress } from '../../hooks/useProgress'
import guideIndex from '../../guide-index.json'
import GuideHeader from './GuideHeader'
import GuideSkeleton from './GuideSkeleton'
import GuideNavigation from './GuideNavigation'
import ReadingProgress from './ReadingProgress'
import TableOfContents from './TableOfContents'
import MarkdownRenderer from '../markdown/MarkdownRenderer'

export default function GuidePage() {
  const { slug } = useParams()
  const { content, loading, error } = useGuideContent(slug)
  const { activeHeading, headings } = useScrollSpy(content)
  const progress = useReadingProgress()
  const { updateProgress } = useProgress(slug)

  // Save reading progress to localStorage
  useEffect(() => {
    if (progress > 5) {
      updateProgress(progress)
    }
  }, [progress, updateProgress])

  const guide = guideIndex.find(g => g.slug === slug)
  const guideIdx = guideIndex.findIndex(g => g.slug === slug)
  const prevGuide = guideIdx > 0 ? guideIndex[guideIdx - 1] : null
  const nextGuide = guideIdx < guideIndex.length - 1 ? guideIndex[guideIdx + 1] : null

  // Scroll to top + set document title on guide change
  useEffect(() => {
    window.scrollTo(0, 0)
    if (guide) {
      document.title = `${guide.title} — Study Guides`
    }
    return () => { document.title = 'Study Guides — Anshul Garg' }
  }, [slug, guide])

  if (loading) return <GuideSkeleton />

  if (error) {
    return (
      <div className="max-w-prose mx-auto px-6 py-16 text-center">
        <h2 className="font-heading font-semibold text-2xl mb-4">Guide not found</h2>
        <p className="text-[var(--text-secondary)]">{error}</p>
      </div>
    )
  }

  return (
    <div className="animate-fade-in">
      <ReadingProgress progress={progress} />

      <div className="flex justify-center">
        {/* Content area */}
        <article className="flex-1 min-w-0 max-w-[780px] px-5 md:px-8 py-8">
          {guide && <GuideHeader guide={guide} />}

          <div className="guide-content">
            <MarkdownRenderer content={content} />
          </div>

          <GuideNavigation prev={prevGuide} next={nextGuide} />
        </article>

        {/* Right TOC */}
        <div className="hidden xl:block w-[200px] shrink-0">
          <TableOfContents
            headings={headings}
            activeHeading={activeHeading}
          />
        </div>
      </div>
    </div>
  )
}

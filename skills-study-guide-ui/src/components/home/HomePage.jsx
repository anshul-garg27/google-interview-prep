import guideIndex from '../../guide-index.json'
import categoryOrder from '../../categories.json'
import HeroSection from './HeroSection'
import GuideGrid from './GuideGrid'
import ContinueReading from './ContinueReading'
import StudyPaths from './StudyPaths'
import { useAllProgress } from '../../hooks/useProgress'

export default function HomePage() {
  const progress = useAllProgress()
  const totalWords = guideIndex.reduce((sum, g) => sum + g.stats.wordCount, 0)
  const totalQA = guideIndex.reduce((sum, g) => sum + g.stats.qaItems, 0)
  const totalDiagrams = guideIndex.reduce((sum, g) => sum + g.stats.mermaidDiagrams, 0)
  const totalReadingHours = Math.round(guideIndex.reduce((sum, g) => sum + g.readingTime, 0) / 60)
  const guidesStarted = Object.keys(progress).length

  return (
    <div className="animate-fade-in">
      <HeroSection
        guideCount={guideIndex.length}
        totalWords={totalWords}
        totalQA={totalQA}
        totalDiagrams={totalDiagrams}
        totalReadingHours={totalReadingHours}
        guidesStarted={guidesStarted}
      />

      <ContinueReading />

      <div className="max-w-5xl mx-auto px-5 pb-16">
        <StudyPaths />

        <GuideGrid
          guides={guideIndex}
          categories={categoryOrder}
          progress={progress}
        />
      </div>
    </div>
  )
}

import GuideCard from './GuideCard'

export default function GuideGrid({ guides, categories, progress }) {
  const grouped = {}
  for (const guide of guides) {
    if (!grouped[guide.category]) grouped[guide.category] = []
    grouped[guide.category].push(guide)
  }

  return (
    <div className="space-y-10">
      {categories.map(cat => {
        const catGuides = grouped[cat]
        if (!catGuides || catGuides.length === 0) return null

        return (
          <section key={cat}>
            <h2 className="font-mono text-xs font-semibold uppercase tracking-[0.08em] text-[var(--text-tertiary)] mb-4 px-1">
              {cat}
            </h2>
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
              {catGuides.map(guide => (
                <GuideCard key={guide.slug} guide={guide} progress={progress} />
              ))}
            </div>
          </section>
        )
      })}
    </div>
  )
}

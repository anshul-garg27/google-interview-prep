export default function HeroSection({ guideCount, totalWords, totalQA, totalDiagrams, totalReadingHours, guidesStarted }) {
  return (
    <section className="max-w-5xl mx-auto px-5 pt-12 pb-10">
      <h1 className="font-heading font-extrabold text-4xl md:text-5xl tracking-tight leading-[1.1] mb-4">
        Interview Study Guides
      </h1>
      <p className="text-[var(--text-secondary)] text-lg font-body leading-relaxed max-w-2xl mb-8">
        {guideCount} comprehensive guides covering backend systems, distributed architecture,
        languages, databases, and interview preparation. Built from real engineering experience.
      </p>

      <div className="flex flex-wrap gap-6 text-sm">
        <Stat value={`${Math.round(totalWords / 1000)}K`} label="words" />
        <Stat value={totalQA} label="interview Q&As" />
        <Stat value={totalDiagrams} label="diagrams" />
        <Stat value={`${totalReadingHours}h`} label="reading time" />
        {guidesStarted > 0 && (
          <Stat value={`${guidesStarted}/${guideCount}`} label="started" />
        )}
      </div>
    </section>
  )
}

function Stat({ value, label }) {
  return (
    <div className="flex items-baseline gap-1.5">
      <span className="font-mono font-semibold text-[var(--accent)] text-base">{value}</span>
      <span className="font-mono text-[11px] text-[var(--text-tertiary)] uppercase tracking-wide">{label}</span>
    </div>
  )
}

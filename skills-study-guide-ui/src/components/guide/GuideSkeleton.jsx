export default function GuideSkeleton() {
  return (
    <div className="max-w-[780px] mx-auto px-5 md:px-8 py-8 animate-pulse">
      {/* Category + number */}
      <div className="flex items-center gap-3 mb-4">
        <div className="h-6 w-28 rounded-md bg-[var(--bg-surface-alt)]" />
        <div className="h-4 w-20 rounded bg-[var(--bg-surface-alt)]" />
      </div>

      {/* Title */}
      <div className="h-10 w-3/4 rounded bg-[var(--bg-surface-alt)] mb-4" />

      {/* Description */}
      <div className="h-5 w-full rounded bg-[var(--bg-surface-alt)] mb-2" />
      <div className="h-5 w-2/3 rounded bg-[var(--bg-surface-alt)] mb-8" />

      {/* Stats */}
      <div className="flex gap-4 mb-10 pb-8 border-b border-[var(--border)]">
        {[1, 2, 3, 4].map(i => (
          <div key={i} className="h-4 w-24 rounded bg-[var(--bg-surface-alt)]" />
        ))}
      </div>

      {/* Content blocks */}
      {Array.from({ length: 6 }, (_, i) => (
        <div key={i} className="mb-8">
          <div className="h-7 w-1/2 rounded bg-[var(--bg-surface-alt)] mb-4" />
          <div className="space-y-2">
            <div className="h-4 w-full rounded bg-[var(--bg-surface-alt)]" />
            <div className="h-4 w-full rounded bg-[var(--bg-surface-alt)]" />
            <div className="h-4 w-5/6 rounded bg-[var(--bg-surface-alt)]" />
          </div>
        </div>
      ))}
    </div>
  )
}

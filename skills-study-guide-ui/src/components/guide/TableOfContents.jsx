export default function TableOfContents({ headings, activeHeading }) {
  if (!headings.length) return null

  return (
    <nav className="sticky top-20 py-8 pr-4 max-h-[calc(100vh-5rem)] overflow-y-auto" aria-label="Table of contents">
      <p className="font-mono text-[10px] font-semibold uppercase tracking-[0.08em] text-[var(--text-tertiary)] mb-3">
        On this page
      </p>

      <ul className="space-y-0.5">
        {headings.filter(h => h.level <= 3).map(heading => (
          <li key={heading.id}>
            <a
              href={`#${heading.id}`}
              title={heading.title.length > 40 ? heading.title : undefined}
              className={`block text-[12px] leading-relaxed py-1 border-l-2 transition-all
                ${heading.level === 3 ? 'pl-4' : 'pl-3'}
                ${activeHeading === heading.id
                  ? 'border-[var(--accent)] text-[var(--accent)] font-medium bg-[var(--accent-subtle)] rounded-r'
                  : 'border-transparent text-[var(--text-tertiary)] hover:text-[var(--text-secondary)] hover:border-[var(--border-strong)]'
                }`}
            >
              {heading.title.length > 40
                ? heading.title.slice(0, 37) + '...'
                : heading.title}
            </a>
          </li>
        ))}
      </ul>
    </nav>
  )
}

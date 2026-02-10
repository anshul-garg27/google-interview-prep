import { useState, useMemo } from 'react'
import { useLocation, useNavigate } from 'react-router-dom'
import {
  Building2,
  Briefcase,
  BookOpen,
  Search,
  FileCode,
  FileText,
  GraduationCap,
  Zap,
  ChevronDown,
  Database,
  Code2,
  Layers,
  BarChart3,
} from 'lucide-react'
import guideIndex from '../../guide-index.json'
import categoryOrder from '../../categories.json'
import { resourceCategories, pinnedResources } from '../../data/standalone-resources'
import projectFiles from '../../project-files.json'

// ── Sidebar project lists ────────────────────────────────────────

const WALMART_PROJECTS = [
  { slug: 'kafka-audit-logging', label: 'Kafka Audit Logging' },
  { slug: 'spring-boot-3-migration', label: 'Spring Boot 3 Migration' },
  { slug: 'dsd-notification-system', label: 'DSD Notifications' },
  { slug: 'common-library-jar', label: 'Common Library JAR' },
  { slug: 'openapi-dc-inventory', label: 'DC Inventory API' },
  { slug: 'transaction-event-history', label: 'Transaction Events' },
  { slug: 'observability', label: 'Observability' },
]

const GCC_PROJECTS = [
  { slug: 'gcc-beat-scraping', label: 'Beat Scraping Engine' },
  { slug: 'gcc-event-grpc', label: 'Event-gRPC Pipeline' },
  { slug: 'gcc-stir-data-platform', label: 'Stir Data Platform' },
  { slug: 'gcc-coffee-saas-api', label: 'Coffee + SaaS Gateway' },
  { slug: 'gcc-fake-follower-ml', label: 'Fake Follower ML' },
]

// ── Icon map for resource categories ─────────────────────────────

const RESOURCE_ICON_MAP = {
  Building2,
  Briefcase,
  BookOpen,
  Search,
  FileCode,
  FileText,
  GraduationCap,
  Zap,
}

const CATEGORY_ICONS = {
  'Infrastructure': Layers,
  'System Design': BarChart3,
  'Databases & Storage': Database,
  'Languages': Code2,
  'Data & Observability': BarChart3,
  'Interview Prep': GraduationCap,
}

// ── Helper: strip .md extension from filename ────────────────────

function stripMd(file) {
  return file.replace(/\.md$/i, '')
}

// ── Section header ───────────────────────────────────────────────

function SectionHeader({ icon: Icon, label }) {
  return (
    <div className="flex items-center gap-2 px-3 mb-2">
      <Icon size={12} className="text-[var(--text-tertiary)]" />
      <span className="font-mono text-[11px] font-semibold uppercase tracking-[0.06em] text-[var(--text-tertiary)]">
        {label}
      </span>
    </div>
  )
}

// ── Sidebar item button ──────────────────────────────────────────

function SidebarItem({ label, isActive, onClick, icon: Icon }) {
  return (
    <button
      onClick={onClick}
      className={`w-full text-left px-3 py-2 rounded-md text-[13px] font-body flex items-start gap-2 transition-colors
        ${isActive
          ? 'sidebar-item-active font-semibold'
          : 'text-[var(--text-secondary)] hover:bg-[var(--bg-hover)]'
        }`}
    >
      {Icon && <Icon size={13} className="mt-0.5 shrink-0 text-[var(--text-tertiary)]" />}
      <span className="leading-snug">{label}</span>
    </button>
  )
}

// ── Collapsible resource category group ──────────────────────────

function CollapsibleGroup({ name, children, defaultOpen = false }) {
  const [open, setOpen] = useState(defaultOpen)

  return (
    <div className="mb-1">
      <button
        onClick={() => setOpen(prev => !prev)}
        className="w-full flex items-center gap-1.5 px-3 py-1.5 rounded-md text-[12px] font-mono font-medium text-[var(--text-tertiary)] hover:bg-[var(--bg-hover)] transition-colors"
        aria-expanded={open}
      >
        <ChevronDown
          size={12}
          className={`shrink-0 transition-transform duration-200 ${open ? '' : '-rotate-90'}`}
        />
        <span className="truncate">{name}</span>
      </button>
      {open && (
        <div className="ml-1 mt-0.5">
          {children}
        </div>
      )}
    </div>
  )
}

// ── Project with expandable sub-files ─────────────────────────────

function ProjectWithFiles({ project, files, location, navigate }) {
  const isProjectActive = location.pathname === `/interview/${project.slug}`
  const isFileActive = location.pathname.startsWith(`/interview/${project.slug}/`)
  const activeFile = isFileActive ? location.pathname.split('/').pop() : null
  const [open, setOpen] = useState(isProjectActive || isFileActive)

  return (
    <div>
      <div className="flex items-center">
        <button
          onClick={() => navigate(`/interview/${project.slug}`)}
          className={`flex-1 text-left px-3 py-2 rounded-md text-[13px] font-body leading-snug transition-colors
            ${isProjectActive
              ? 'sidebar-item-active font-semibold'
              : 'text-[var(--text-secondary)] hover:bg-[var(--bg-hover)]'
            }`}
        >
          {project.label}
        </button>
        {files.length > 0 && (
          <button
            onClick={() => setOpen(prev => !prev)}
            className="p-1.5 rounded hover:bg-[var(--bg-hover)] text-[var(--text-tertiary)]"
            aria-label="Toggle files"
          >
            <ChevronDown size={12} className={`transition-transform duration-200 ${open ? '' : '-rotate-90'}`} />
          </button>
        )}
      </div>
      {open && files.length > 0 && (
        <div className="ml-3 mt-0.5 mb-1 border-l border-[var(--border)] pl-2">
          {files.map(f => {
            const fileSlug = f.file.replace(/\.md$/, '')
            const isActive = activeFile === fileSlug
            return (
              <button
                key={f.file}
                onClick={() => navigate(`/interview/${project.slug}/${fileSlug}`)}
                className={`w-full text-left px-2 py-1 rounded text-[11px] font-mono leading-snug transition-colors truncate
                  ${isActive
                    ? 'text-[var(--accent)] font-semibold bg-[var(--accent-subtle)]'
                    : 'text-[var(--text-tertiary)] hover:text-[var(--text-secondary)] hover:bg-[var(--bg-hover)]'
                  }`}
                title={f.label}
              >
                {f.label}
              </button>
            )
          })}
        </div>
      )}
    </div>
  )
}

// ── Main Sidebar component ───────────────────────────────────────

export default function Sidebar({ open, onClose }) {
  const location = useLocation()
  const navigate = useNavigate()

  // Study guide categories (for section 4)
  const grouped = useMemo(() => {
    const groups = {}
    for (const guide of guideIndex) {
      if (!groups[guide.category]) groups[guide.category] = []
      groups[guide.category].push(guide)
    }
    return categoryOrder.map(cat => ({
      name: cat,
      guides: groups[cat] || [],
    }))
  }, [])

  const currentSlug = location.pathname.startsWith('/guide/')
    ? location.pathname.replace('/guide/', '')
    : null

  return (
    <aside
      className={`fixed top-14 left-0 bottom-0 w-[260px] z-40 bg-[var(--bg-surface)] border-r border-[var(--border)] overflow-y-auto transition-transform duration-250 ease-out
        ${open ? 'translate-x-0' : '-translate-x-full'} lg:translate-x-0`}
    >
      <nav className="py-4 px-3">

        {/* ── Section 1: Walmart Projects ──────────────────── */}
        <div className="mb-5">
          <SectionHeader icon={Building2} label="Walmart Projects" />
          {WALMART_PROJECTS.map(project => (
            <ProjectWithFiles
              key={project.slug}
              project={project}
              files={projectFiles[project.slug] || []}
              location={location}
              navigate={navigate}
            />
          ))}
        </div>

        {/* ── Section 2: GCC Projects ─────────────────────── */}
        <div className="mb-5 pt-4 border-t border-[var(--border)]">
          <SectionHeader icon={Briefcase} label="GCC Projects" />
          {GCC_PROJECTS.map(project => (
            <ProjectWithFiles
              key={project.slug}
              project={project}
              files={projectFiles[project.slug] || []}
              location={location}
              navigate={navigate}
            />
          ))}
        </div>

        {/* ── Section 3: Interview Resources ──────────────── */}
        <div className="mb-5 pt-4 border-t border-[var(--border)]">
          <SectionHeader icon={BookOpen} label="Key Resources" />

          {/* Pinned files — always visible, no collapsing */}
          {pinnedResources.map(resource => {
            const slug = stripMd(resource.file)
            return (
              <SidebarItem
                key={resource.file}
                label={resource.label}
                isActive={location.pathname === `/resource/${slug}`}
                onClick={() => navigate(`/resource/${slug}`)}
              />
            )
          })}

          {/* Other resources — collapsible categories */}
          <div className="mt-3 pt-3 border-t border-[var(--border)]">
            <SectionHeader icon={FileText} label="All Resources" />
          </div>
          {resourceCategories.map(category => {
            const CatIcon = RESOURCE_ICON_MAP[category.icon] || FileText
            // Check if any resource in this category is active
            const hasActive = category.resources.some(
              r => location.pathname === `/resource/${stripMd(r.file)}`
            )
            return (
              <CollapsibleGroup
                key={category.name}
                name={category.name}
                defaultOpen={hasActive}
              >
                {category.resources.map(resource => {
                  const slug = stripMd(resource.file)
                  return (
                    <SidebarItem
                      key={resource.file}
                      label={resource.label}
                      isActive={location.pathname === `/resource/${slug}`}
                      onClick={() => navigate(`/resource/${slug}`)}
                    />
                  )
                })}
              </CollapsibleGroup>
            )
          })}
        </div>

        {/* ── Section 4: Study Guides ─────────────────────── */}
        <div className="pt-4 border-t border-[var(--border)]">
          {grouped.map(group => {
            const Icon = CATEGORY_ICONS[group.name] || BookOpen
            return (
              <div key={group.name} className="mb-5">
                <SectionHeader icon={Icon} label={group.name} />
                {group.guides.map(guide => {
                  const isActive = currentSlug === guide.slug
                  return (
                    <button
                      key={guide.slug}
                      onClick={() => navigate(`/guide/${guide.slug}`)}
                      className={`w-full text-left px-3 py-2 rounded-md text-[13px] font-body flex items-start gap-2 transition-colors
                        ${isActive
                          ? 'sidebar-item-active font-semibold'
                          : 'text-[var(--text-secondary)] hover:bg-[var(--bg-hover)]'
                        }`}
                    >
                      <span className="font-mono text-[11px] text-[var(--text-tertiary)] mt-0.5 shrink-0 w-5">
                        {guide.num}
                      </span>
                      <span className="leading-snug">{guide.title}</span>
                    </button>
                  )
                })}
              </div>
            )
          })}
        </div>

        {/* ── Stats ───────────────────────────────────────── */}
        <div className="mt-6 pt-4 border-t border-[var(--border)] px-3">
          <p className="font-mono text-[10px] uppercase tracking-widest text-[var(--text-tertiary)]">
            {guideIndex.length} guides &middot; {Math.round(guideIndex.reduce((sum, g) => sum + g.stats.wordCount, 0) / 1000)}K words
          </p>
        </div>
      </nav>
    </aside>
  )
}

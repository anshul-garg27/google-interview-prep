import { Briefcase, Lightbulb, AlertTriangle, Info } from 'lucide-react'

const CALLOUT_TYPES = {
  anshul: {
    icon: Briefcase,
    label: 'Real-World Application',
    borderColor: 'border-l-[var(--accent)]',
    bgColor: 'bg-[var(--accent-subtle)]',
    iconColor: 'text-[var(--accent)]',
  },
  tip: {
    icon: Lightbulb,
    label: 'Tip',
    borderColor: 'border-l-emerald-500',
    bgColor: 'bg-emerald-50 dark:bg-emerald-950/20',
    iconColor: 'text-emerald-600 dark:text-emerald-400',
  },
  warning: {
    icon: AlertTriangle,
    label: 'Warning',
    borderColor: 'border-l-amber-500',
    bgColor: 'bg-amber-50 dark:bg-amber-950/20',
    iconColor: 'text-amber-600 dark:text-amber-400',
  },
  info: {
    icon: Info,
    label: 'Note',
    borderColor: 'border-l-blue-500',
    bgColor: 'bg-blue-50 dark:bg-blue-950/20',
    iconColor: 'text-blue-600 dark:text-blue-400',
  },
}

export default function CalloutBox({ type = 'anshul', title, children }) {
  const config = CALLOUT_TYPES[type] || CALLOUT_TYPES.anshul
  const Icon = config.icon

  return (
    <div className={`my-6 rounded-r-lg border-l-4 ${config.borderColor} ${config.bgColor} p-4 not-prose`}>
      <div className="flex items-center gap-2 mb-2">
        <Icon size={14} className={config.iconColor} />
        <span className={`font-mono text-[10px] font-semibold uppercase tracking-[0.06em] ${config.iconColor}`}>
          {title || config.label}
        </span>
      </div>
      <div className="text-[14px] leading-relaxed text-[var(--text-primary)]">
        {children}
      </div>
    </div>
  )
}

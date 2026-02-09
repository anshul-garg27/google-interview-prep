import { useState } from 'react'
import { ArrowLeft } from 'lucide-react'
import { useNavigate } from 'react-router-dom'
import Approach1RawMarkdown from './Approach1RawMarkdown'
import Approach2EnhancedMarkdown from './Approach2EnhancedMarkdown'
import Approach3MDXHybrid from './Approach3MDXHybrid'
import Approach4CustomReact from './Approach4CustomReact'

const TABS = [
  { id: 'raw', label: 'Raw Markdown', sublabel: 'Current approach' },
  { id: 'enhanced', label: 'Enhanced MD', sublabel: 'Better plugins' },
  { id: 'mdx', label: 'MDX Hybrid', sublabel: 'MD + React' },
  { id: 'custom', label: 'Custom React', sublabel: 'Full interactive' },
]

export default function DemoPage() {
  const [activeTab, setActiveTab] = useState('raw')
  const navigate = useNavigate()

  return (
    <div className="animate-fade-in">
      {/* Header */}
      <div className="max-w-5xl mx-auto px-5 pt-8 pb-6">
        <button
          onClick={() => navigate('/')}
          className="flex items-center gap-2 text-[var(--text-tertiary)] hover:text-[var(--text-primary)] transition-colors mb-6 font-mono text-xs"
        >
          <ArrowLeft size={14} /> Back to guides
        </button>

        <h1 className="font-heading font-extrabold text-3xl md:text-4xl tracking-tight mb-3">
          Content Rendering: 4 Approaches
        </h1>
        <p className="text-[var(--text-secondary)] text-lg font-body leading-relaxed max-w-2xl mb-2">
          Same Kafka "Consumer Groups" content — rendered 4 different ways.
          Switch tabs to compare the reading experience.
        </p>
        <p className="text-[var(--text-tertiary)] text-sm font-mono">
          Demo section from Guide 01: Kafka, Avro & Schema Registry
        </p>
      </div>

      {/* Tab bar */}
      <div className="sticky top-14 z-30 bg-[var(--bg-primary)]/90 backdrop-blur-xl border-b border-[var(--border)]">
        <div className="max-w-5xl mx-auto px-5">
          <div className="flex gap-1 py-2 overflow-x-auto">
            {TABS.map((tab, i) => (
              <button
                key={tab.id}
                onClick={() => setActiveTab(tab.id)}
                className={`flex-shrink-0 px-4 py-2.5 rounded-lg font-body text-sm transition-all
                  ${activeTab === tab.id
                    ? 'bg-[var(--accent)] text-white shadow-sm'
                    : 'hover:bg-[var(--bg-hover)] text-[var(--text-secondary)]'
                  }`}
              >
                <div className="font-semibold">{tab.label}</div>
                <div className={`text-[10px] font-mono mt-0.5 ${activeTab === tab.id ? 'text-white/70' : 'text-[var(--text-tertiary)]'}`}>
                  {tab.sublabel}
                </div>
              </button>
            ))}
          </div>
        </div>
      </div>

      {/* Content area */}
      <div className="max-w-4xl mx-auto px-5 py-8">
        <div className="animate-fade-in" key={activeTab}>
          {activeTab === 'raw' && <Approach1RawMarkdown />}
          {activeTab === 'enhanced' && <Approach2EnhancedMarkdown />}
          {activeTab === 'mdx' && <Approach3MDXHybrid />}
          {activeTab === 'custom' && <Approach4CustomReact />}
        </div>
      </div>
    </div>
  )
}

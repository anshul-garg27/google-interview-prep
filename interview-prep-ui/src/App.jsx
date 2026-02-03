import { useState, useEffect, useMemo } from 'react'
import ReactMarkdown from 'react-markdown'
import remarkGfm from 'remark-gfm'
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter'
import { oneDark, oneLight } from 'react-syntax-highlighter/dist/esm/styles/prism'
import Fuse from 'fuse.js'
import {
  Menu, X, Search, Moon, Sun, ChevronRight, ChevronDown,
  BookOpen, FileText, Code, Brain, Network, Building2,
  Target, MessageSquare, Briefcase, Github, Zap, Layers,
  TrendingUp, Users, Award, Clock, ArrowRight, Home
} from 'lucide-react'
import { documents } from './data'

// Enhanced document metadata
const getDocumentMetadata = (doc) => {
  const company = doc.category.includes('walmart') ? 'Walmart' :
                 doc.category.includes('google') ? 'Good Creator Co' :
                 'General'

  const wordCount = doc.content.split(/\s+/).length
  const readingTime = Math.ceil(wordCount / 200) // 200 words per minute

  const techStack = doc.content.match(/(?:Spring Boot|Kafka|PostgreSQL|ClickHouse|React|Python|Go|Java|Redis|Kubernetes|Airflow|dbt|AWS Lambda|FastAPI)/gi) || []
  const uniqueTech = [...new Set(techStack.map(t => t.toLowerCase()))]

  return {
    company,
    readingTime,
    techStack: uniqueTech.slice(0, 5),
    wordCount,
    sections: (doc.content.match(/^#{2,3}\s+/gm) || []).length
  }
}

// Convert blockquote code to proper code blocks
const preprocessMarkdown = (content) => {
  const lines = content.split('\n')
  const processed = []
  let inBlockquote = false
  let inCodeFence = false
  let blockquoteBuffer = []
  let codeFenceLanguage = ''

  for (let i = 0; i < lines.length; i++) {
    const line = lines[i]
    const trimmedLine = line.replace(/^>\s?/, '') // Remove blockquote prefix

    // Check if we're starting/ending a blockquote
    if (line.startsWith('>') && !inBlockquote) {
      inBlockquote = true
      blockquoteBuffer = []
    }

    if (inBlockquote) {
      // Check for code fence markers inside blockquote
      const codeFenceMatch = trimmedLine.match(/^```(\w*)/)

      if (codeFenceMatch && !inCodeFence) {
        // Start of code fence inside blockquote
        inCodeFence = true
        codeFenceLanguage = codeFenceMatch[1] || 'java'
        blockquoteBuffer.push(trimmedLine)
      } else if (trimmedLine.match(/^```$/) && inCodeFence) {
        // End of code fence inside blockquote
        inCodeFence = false
        blockquoteBuffer.push(trimmedLine)
      } else if (line.startsWith('>')) {
        // Regular blockquote line
        blockquoteBuffer.push(trimmedLine)
      } else if (!line.startsWith('>') && line.trim() === '' && inBlockquote) {
        // Empty line might still be in blockquote
        blockquoteBuffer.push('')
      } else {
        // End of blockquote
        // Process the blockquote buffer
        processed.push(...processBlockquoteBuffer(blockquoteBuffer))
        processed.push(line)
        inBlockquote = false
        inCodeFence = false
        blockquoteBuffer = []
      }
    } else {
      processed.push(line)
    }
  }

  // Flush remaining blockquote
  if (blockquoteBuffer.length > 0) {
    processed.push(...processBlockquoteBuffer(blockquoteBuffer))
  }

  return processed.join('\n')
}

// Process blockquote buffer to handle code blocks inside
const processBlockquoteBuffer = (buffer) => {
  const result = []
  let inCodeBlock = false
  let codeBuffer = []
  let language = 'java'

  for (let i = 0; i < buffer.length; i++) {
    const line = buffer[i]

    // Check for code fence
    const startFence = line.match(/^```(\w*)/)
    const endFence = line.match(/^```$/)

    if (startFence && !inCodeBlock) {
      // Start of code block
      inCodeBlock = true
      language = startFence[1] || 'java'
      codeBuffer = []
    } else if (endFence && inCodeBlock) {
      // End of code block - output it properly
      result.push('```' + language)
      result.push(...codeBuffer)
      result.push('```')
      result.push('') // Empty line after code block
      inCodeBlock = false
      codeBuffer = []
    } else if (inCodeBlock) {
      // Inside code block - collect lines
      codeBuffer.push(line)
    } else {
      // Regular blockquote line - keep as blockquote
      result.push('> ' + line)
    }
  }

  // Flush any remaining code
  if (inCodeBlock && codeBuffer.length > 0) {
    result.push('```' + language)
    result.push(...codeBuffer)
    result.push('```')
  }

  return result
}

// Quick start paths
const quickStartPaths = [
  {
    id: 'essential',
    name: 'Essential Reading',
    description: 'Core interview prep in 2 hours',
    icon: Zap,
    color: 'amber',
    docs: ['GOOGLE_INTERVIEW_MASTER_GUIDE', 'GOOGLEYNESS_ALL_QUESTIONS', 'WALMART_INTERVIEW_ALL_QUESTIONS']
  },
  {
    id: 'walmart-deep',
    name: 'Walmart Deep Dive',
    description: 'Complete technical overview',
    icon: Building2,
    color: 'blue',
    docs: ['WALMART_SYSTEM_ARCHITECTURE', 'WALMART_HIRING_MANAGER_GUIDE', 'WALMART_RESUME_TO_CODE_MAPPING']
  },
  {
    id: 'behavioral',
    name: 'Behavioral Mastery',
    description: 'All STAR stories & questions',
    icon: Users,
    color: 'purple',
    docs: ['WALMART_GOOGLEYNESS_QUESTIONS', 'WALMART_LEADERSHIP_STORIES', 'GOOGLEYNESS_ALL_QUESTIONS']
  }
]

function App() {
  const [darkMode, setDarkMode] = useState(() => {
    if (typeof window !== 'undefined') {
      return localStorage.getItem('darkMode') === 'true' ||
             window.matchMedia('(prefers-color-scheme: dark)').matches
    }
    return false
  })
  const [sidebarOpen, setSidebarOpen] = useState(false)
  const [activeDoc, setActiveDoc] = useState(null)
  const [searchQuery, setSearchQuery] = useState('')
  const [searchResults, setSearchResults] = useState([])
  const [showSearch, setShowSearch] = useState(false)
  const [expandedCategories, setExpandedCategories] = useState(['master', 'walmart-interview'])
  const [tocExpanded, setTocExpanded] = useState(true)

  // Apply dark mode
  useEffect(() => {
    document.documentElement.classList.toggle('dark', darkMode)
    localStorage.setItem('darkMode', darkMode)
  }, [darkMode])

  // Enhanced search with context
  const fuse = useMemo(() => {
    const searchData = documents.map(doc => ({
      id: doc.id,
      title: doc.title,
      content: doc.content,
      category: doc.category,
      badge: doc.badge
    }))
    return new Fuse(searchData, {
      keys: [
        { name: 'title', weight: 2 },
        { name: 'content', weight: 1 }
      ],
      threshold: 0.3,
      includeMatches: true,
      minMatchCharLength: 3
    })
  }, [])

  useEffect(() => {
    if (searchQuery.length > 2) {
      const results = fuse.search(searchQuery).slice(0, 8)
      setSearchResults(results)
    } else {
      setSearchResults([])
    }
  }, [searchQuery, fuse])

  // Extract TOC from current document
  const currentDocTOC = useMemo(() => {
    if (!activeDoc) return []
    const doc = documents.find(d => d.id === activeDoc)
    if (!doc) return []

    const lines = doc.content.split('\n')
    const headings = []
    lines.forEach((line, index) => {
      const match = line.match(/^(#{2,3})\s+(.+)$/)
      if (match) {
        const level = match[1].length
        const text = match[2].replace(/[#*`]/g, '')
        const id = text.toLowerCase().replace(/\s+/g, '-').replace(/[^\w-]/g, '')
        headings.push({ level, text, id, line: index })
      }
    })
    return headings
  }, [activeDoc])

  const currentDoc = documents.find(d => d.id === activeDoc)
  const metadata = currentDoc ? getDocumentMetadata(currentDoc) : null

  // Organized categories with icons
  const categories = [
    {
      id: 'master',
      name: '🎯 Start Here',
      icon: Target,
      color: 'emerald',
      docs: documents.filter(d => d.category === 'master')
    },
    {
      id: 'walmart',
      name: 'Walmart (Current)',
      icon: Building2,
      color: 'blue',
      subcategories: [
        {
          id: 'walmart-interview',
          name: 'Interview Prep',
          icon: MessageSquare,
          docs: documents.filter(d => d.category === 'walmart-interview')
        },
        {
          id: 'walmart-detailed',
          name: 'Deep Dives',
          icon: Layers,
          docs: documents.filter(d => d.category === 'walmart-detailed')
        },
        {
          id: 'walmart-technical',
          name: 'Technical Docs',
          icon: Code,
          docs: documents.filter(d => d.category === 'walmart-technical')
        },
        {
          id: 'walmart-microservices',
          name: 'Microservices',
          icon: Network,
          docs: documents.filter(d => d.category === 'walmart-microservices')
        },
        {
          id: 'walmart-portfolio',
          name: 'Portfolio',
          icon: Briefcase,
          docs: documents.filter(d => d.category === 'walmart-portfolio')
        }
      ]
    },
    {
      id: 'google',
      name: 'Good Creator Co (Previous)',
      icon: Brain,
      color: 'purple',
      subcategories: [
        {
          id: 'google-interview',
          name: 'Interview Prep',
          icon: MessageSquare,
          docs: documents.filter(d => d.category === 'google-interview')
        },
        {
          id: 'google-technical',
          name: 'Technical Docs',
          icon: Code,
          docs: documents.filter(d => d.category === 'google-technical')
        },
        {
          id: 'google-analysis',
          name: 'Project Analysis',
          icon: FileText,
          docs: documents.filter(d => d.category === 'google-analysis')
        }
      ]
    }
  ]

  const handleDocSelect = (docId) => {
    setActiveDoc(docId)
    setSidebarOpen(false)
    setShowSearch(false)
    setSearchQuery('')
    setSearchResults([])
    window.scrollTo(0, 0)
  }

  const toggleCategory = (categoryId) => {
    setExpandedCategories(prev =>
      prev.includes(categoryId)
        ? prev.filter(id => id !== categoryId)
        : [...prev, categoryId]
    )
  }

  const getSearchContext = (content, matches) => {
    if (!matches || matches.length === 0) return ''
    const match = matches[0]
    const start = Math.max(0, match.indices[0][0] - 60)
    const end = Math.min(content.length, match.indices[0][1] + 60)
    let context = content.slice(start, end)
    if (start > 0) context = '...' + context
    if (end < content.length) context = context + '...'
    return context
  }

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-950 transition-colors duration-300">
      {/* Header */}
      <header className="fixed top-0 left-0 right-0 z-50 bg-white/80 dark:bg-gray-900/80 backdrop-blur-xl border-b border-gray-200 dark:border-gray-800">
        <div className="flex items-center justify-between px-4 lg:px-6 h-16">
          {/* Left side */}
          <div className="flex items-center gap-4">
            <button
              onClick={() => setSidebarOpen(!sidebarOpen)}
              className="lg:hidden p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-800 transition-colors"
            >
              {sidebarOpen ? <X size={20} /> : <Menu size={20} />}
            </button>

            <button
              onClick={() => setActiveDoc(null)}
              className="flex items-center gap-3 group"
            >
              <div className="flex items-center gap-1.5">
                <div className="w-2 h-2 rounded-full bg-blue-500 animate-pulse"></div>
                <div className="w-2 h-2 rounded-full bg-red-500 animate-pulse" style={{animationDelay: '0.2s'}}></div>
                <div className="w-2 h-2 rounded-full bg-amber-500 animate-pulse" style={{animationDelay: '0.4s'}}></div>
                <div className="w-2 h-2 rounded-full bg-emerald-500 animate-pulse" style={{animationDelay: '0.6s'}}></div>
              </div>
              <div>
                <h1 className="text-base font-mono font-semibold text-gray-900 dark:text-white tracking-tight">
                  INTERVIEW.PREP
                </h1>
                <p className="text-[10px] font-mono text-gray-500 dark:text-gray-400 -mt-0.5">
                  Google L4 • {documents.length} docs
                </p>
              </div>
            </button>
          </div>

          {/* Right side actions */}
          <div className="flex items-center gap-1">
            <button
              onClick={() => setShowSearch(!showSearch)}
              className="p-2.5 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-800 transition-colors"
              title="Search (⌘K)"
            >
              <Search size={18} className="text-gray-600 dark:text-gray-400" />
            </button>

            <button
              onClick={() => setDarkMode(!darkMode)}
              className="p-2.5 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-800 transition-colors"
            >
              {darkMode ? (
                <Sun size={18} className="text-amber-500" />
              ) : (
                <Moon size={18} className="text-gray-600" />
              )}
            </button>

            <a
              href="https://github.com/anshul-garg27/google-interview-prep"
              target="_blank"
              rel="noopener noreferrer"
              className="p-2.5 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-800 transition-colors"
            >
              <Github size={18} className="text-gray-600 dark:text-gray-400" />
            </a>
          </div>
        </div>

        {/* Enhanced Search bar */}
        {showSearch && (
          <div className="px-4 lg:px-6 pb-4 bg-white/80 dark:bg-gray-900/80 backdrop-blur-xl">
            <div className="max-w-2xl mx-auto">
              <div className="relative">
                <Search size={16} className="absolute left-4 top-1/2 -translate-y-1/2 text-gray-400" />
                <input
                  type="text"
                  placeholder="Search across 32 documents..."
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  className="w-full pl-12 pr-4 py-3 rounded-xl border border-gray-300 dark:border-gray-700 bg-white dark:bg-gray-900 text-gray-900 dark:text-white placeholder:text-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500/50 font-mono text-sm transition-all"
                  autoFocus
                />
                {searchQuery && (
                  <button
                    onClick={() => setSearchQuery('')}
                    className="absolute right-4 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600 dark:hover:text-gray-300"
                  >
                    <X size={16} />
                  </button>
                )}
              </div>

              {/* Enhanced Search results */}
              {searchResults.length > 0 && (
                <div className="mt-3 bg-white dark:bg-gray-900 rounded-xl shadow-2xl border border-gray-200 dark:border-gray-800 max-h-96 overflow-y-auto">
                  {searchResults.map((result, index) => {
                    const meta = getDocumentMetadata(result.item)
                    const context = getSearchContext(result.item.content, result.matches)
                    return (
                      <button
                        key={result.item.id}
                        onClick={() => handleDocSelect(result.item.id)}
                        className={`w-full text-left px-4 py-3 hover:bg-gray-50 dark:hover:bg-gray-800/50 transition-colors ${
                          index !== searchResults.length - 1 ? 'border-b border-gray-100 dark:border-gray-800' : ''
                        }`}
                      >
                        <div className="flex items-start justify-between gap-3 mb-1">
                          <div className="font-semibold text-gray-900 dark:text-white text-sm">
                            {result.item.title}
                          </div>
                          {result.item.badge && (
                            <span className="text-[10px] font-mono px-2 py-0.5 rounded bg-emerald-100 dark:bg-emerald-900/30 text-emerald-700 dark:text-emerald-400 shrink-0">
                              {result.item.badge}
                            </span>
                          )}
                        </div>
                        <div className="flex items-center gap-2 text-xs text-gray-500 dark:text-gray-400 mb-1">
                          <span className="font-mono">{meta.company}</span>
                          <span>•</span>
                          <span>{meta.readingTime} min read</span>
                        </div>
                        {context && (
                          <p className="text-xs text-gray-600 dark:text-gray-400 line-clamp-2 font-mono">
                            {context}
                          </p>
                        )}
                      </button>
                    )
                  })}
                </div>
              )}

              {searchQuery.length > 2 && searchResults.length === 0 && (
                <div className="mt-3 text-center py-8 text-gray-500 dark:text-gray-400 text-sm font-mono">
                  No results found for "{searchQuery}"
                </div>
              )}
            </div>
          </div>
        )}
      </header>

      {/* Sidebar overlay for mobile */}
      {sidebarOpen && (
        <div
          className="fixed inset-0 z-40 bg-black/50 backdrop-blur-sm lg:hidden"
          onClick={() => setSidebarOpen(false)}
        />
      )}

      {/* Sidebar */}
      <aside
        className={`fixed top-16 left-0 z-40 w-80 h-[calc(100vh-4rem)] bg-white dark:bg-gray-900 border-r border-gray-200 dark:border-gray-800 overflow-y-auto transition-transform duration-300 ${
          sidebarOpen ? 'translate-x-0' : '-translate-x-full'
        } lg:translate-x-0`}
      >
        <nav className="p-4">
          {categories.map((category) => (
            <div key={category.id} className="mb-6">
              {/* Category header */}
              <button
                onClick={() => toggleCategory(category.id)}
                className="w-full flex items-center justify-between px-3 py-2 text-xs font-mono font-bold text-gray-900 dark:text-white uppercase tracking-wider hover:bg-gray-100 dark:hover:bg-gray-800 rounded-lg transition-colors"
              >
                <div className="flex items-center gap-2">
                  <category.icon size={14} />
                  {category.name}
                </div>
                <ChevronDown
                  size={14}
                  className={`transition-transform ${
                    expandedCategories.includes(category.id) ? 'rotate-0' : '-rotate-90'
                  }`}
                />
              </button>

              {expandedCategories.includes(category.id) && (
                <>
                  {/* Direct docs */}
                  {category.docs && category.docs.length > 0 && (
                    <ul className="mt-2 space-y-1">
                      {category.docs.map((doc) => (
                        <li key={doc.id}>
                          <button
                            onClick={() => handleDocSelect(doc.id)}
                            className={`w-full text-left px-3 py-2.5 rounded-lg text-xs transition-all ${
                              activeDoc === doc.id
                                ? 'bg-blue-50 dark:bg-blue-950/50 text-blue-700 dark:text-blue-400 font-semibold border-l-2 border-blue-500'
                                : 'text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-800/50'
                            }`}
                          >
                            <div className="flex items-start justify-between gap-2">
                              <span className="flex-1">{doc.title}</span>
                              {doc.badge && (
                                <span className="text-[9px] font-mono px-1.5 py-0.5 rounded bg-emerald-100 dark:bg-emerald-900/30 text-emerald-700 dark:text-emerald-400 shrink-0">
                                  {doc.badge}
                                </span>
                              )}
                            </div>
                          </button>
                        </li>
                      ))}
                    </ul>
                  )}

                  {/* Subcategories */}
                  {category.subcategories && (
                    <div className="mt-3 space-y-4">
                      {category.subcategories.map((subcat) => (
                        <div key={subcat.id}>
                          <button
                            onClick={() => toggleCategory(subcat.id)}
                            className="w-full flex items-center justify-between px-3 py-1.5 text-[10px] font-mono font-semibold text-gray-600 dark:text-gray-400 uppercase tracking-wider hover:text-gray-900 dark:hover:text-white transition-colors"
                          >
                            <div className="flex items-center gap-2">
                              <subcat.icon size={12} />
                              {subcat.name}
                              <span className="text-[9px] text-gray-400">({subcat.docs.length})</span>
                            </div>
                            <ChevronDown
                              size={12}
                              className={`transition-transform ${
                                expandedCategories.includes(subcat.id) ? 'rotate-0' : '-rotate-90'
                              }`}
                            />
                          </button>

                          {expandedCategories.includes(subcat.id) && (
                            <ul className="mt-1 space-y-1">
                              {subcat.docs.map((doc) => (
                                <li key={doc.id}>
                                  <button
                                    onClick={() => handleDocSelect(doc.id)}
                                    className={`w-full text-left px-3 py-2 rounded-lg text-xs transition-all ${
                                      activeDoc === doc.id
                                        ? 'bg-blue-50 dark:bg-blue-950/50 text-blue-700 dark:text-blue-400 font-semibold border-l-2 border-blue-500'
                                        : 'text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-800/50'
                                    }`}
                                  >
                                    <div className="flex items-start justify-between gap-2">
                                      <span className="flex-1 leading-relaxed">{doc.title}</span>
                                      {doc.badge && (
                                        <span className="text-[9px] font-mono px-1.5 py-0.5 rounded bg-emerald-100 dark:bg-emerald-900/30 text-emerald-700 dark:text-emerald-400 shrink-0">
                                          {doc.badge}
                                        </span>
                                      )}
                                    </div>
                                  </button>
                                </li>
                              ))}
                            </ul>
                          )}
                        </div>
                      ))}
                    </div>
                  )}
                </>
              )}
            </div>
          ))}
        </nav>

        {/* Quick stats */}
        <div className="p-4 border-t border-gray-200 dark:border-gray-800">
          <div className="space-y-2">
            <div className="flex items-center justify-between px-3 py-2 bg-blue-50 dark:bg-blue-950/30 rounded-lg">
              <span className="text-xs font-mono text-gray-600 dark:text-gray-400">Total Docs</span>
              <span className="text-lg font-bold font-mono text-blue-600 dark:text-blue-400">{documents.length}</span>
            </div>
            <div className="flex items-center justify-between px-3 py-2 bg-emerald-50 dark:bg-emerald-950/30 rounded-lg">
              <span className="text-xs font-mono text-gray-600 dark:text-gray-400">Questions</span>
              <span className="text-lg font-bold font-mono text-emerald-600 dark:text-emerald-400">140+</span>
            </div>
          </div>
        </div>
      </aside>

      {/* Main content */}
      <main className="lg:ml-80 pt-16 min-h-screen">
        {!activeDoc ? (
          /* Dashboard / Home view */
          <div className="max-w-6xl mx-auto px-6 py-12">
            {/* Hero section */}
            <div className="mb-12">
              <h1 className="text-5xl font-bold font-mono text-gray-900 dark:text-white mb-4 tracking-tight">
                Google L4 Interview
                <br />
                <span className="text-blue-600 dark:text-blue-400">Preparation Hub</span>
              </h1>
              <p className="text-lg text-gray-600 dark:text-gray-400 max-w-2xl leading-relaxed">
                Comprehensive interview materials covering <span className="font-mono font-semibold">6 microservices</span>,
                {' '}<span className="font-mono font-semibold">140+ questions</span>, and
                {' '}<span className="font-mono font-semibold">2 companies</span> worth of experience.
              </p>
            </div>

            {/* Quick start paths */}
            <div className="mb-16">
              <h2 className="text-sm font-mono font-bold text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-6">
                Quick Start Paths
              </h2>
              <div className="grid md:grid-cols-3 gap-4">
                {quickStartPaths.map((path) => {
                  const PathIcon = path.icon
                  return (
                    <button
                      key={path.id}
                      onClick={() => handleDocSelect(path.docs[0])}
                      className="group relative text-left p-6 rounded-2xl bg-white dark:bg-gray-900 border-2 border-gray-200 dark:border-gray-800 hover:border-blue-500 dark:hover:border-blue-500 transition-all hover:shadow-xl"
                    >
                      <div className="flex items-start justify-between mb-4">
                        <div className={`p-3 rounded-xl bg-${path.color}-100 dark:bg-${path.color}-950/30`}>
                          <PathIcon size={24} className={`text-${path.color}-600 dark:text-${path.color}-400`} />
                        </div>
                        <ArrowRight size={20} className="text-gray-400 group-hover:text-blue-500 group-hover:translate-x-1 transition-all" />
                      </div>
                      <h3 className="text-lg font-bold text-gray-900 dark:text-white mb-2">
                        {path.name}
                      </h3>
                      <p className="text-sm text-gray-600 dark:text-gray-400 mb-4">
                        {path.description}
                      </p>
                      <div className="flex items-center gap-2 text-xs font-mono text-gray-500">
                        <FileText size={12} />
                        <span>{path.docs.length} documents</span>
                      </div>
                    </button>
                  )
                })}
              </div>
            </div>

            {/* Company overview */}
            <div className="grid md:grid-cols-2 gap-6 mb-16">
              {/* Walmart */}
              <div className="p-8 rounded-2xl bg-gradient-to-br from-blue-50 to-blue-100/50 dark:from-blue-950/30 dark:to-blue-900/20 border border-blue-200 dark:border-blue-900/50">
                <div className="flex items-center gap-3 mb-4">
                  <div className="p-2 rounded-lg bg-blue-600 dark:bg-blue-500">
                    <Building2 size={20} className="text-white" />
                  </div>
                  <div>
                    <h3 className="text-xl font-bold font-mono text-gray-900 dark:text-white">Walmart</h3>
                    <p className="text-xs font-mono text-gray-600 dark:text-gray-400">Current Company • 2024-Present</p>
                  </div>
                </div>
                <div className="space-y-2 mb-4">
                  <p className="text-sm text-gray-700 dark:text-gray-300">
                    6 microservices • Spring Boot 3 • Kafka • PostgreSQL
                  </p>
                  <p className="text-sm text-gray-600 dark:text-gray-400">
                    Multi-region inventory systems processing 2M+ events/day
                  </p>
                </div>
                <div className="grid grid-cols-2 gap-3">
                  <div className="p-3 rounded-lg bg-white/50 dark:bg-gray-900/50">
                    <div className="text-2xl font-bold font-mono text-blue-600 dark:text-blue-400">18</div>
                    <div className="text-xs font-mono text-gray-600 dark:text-gray-400">Documents</div>
                  </div>
                  <div className="p-3 rounded-lg bg-white/50 dark:bg-gray-900/50">
                    <div className="text-2xl font-bold font-mono text-blue-600 dark:text-blue-400">12</div>
                    <div className="text-xs font-mono text-gray-600 dark:text-gray-400">Bullets</div>
                  </div>
                </div>
              </div>

              {/* Good Creator Co */}
              <div className="p-8 rounded-2xl bg-gradient-to-br from-purple-50 to-purple-100/50 dark:from-purple-950/30 dark:to-purple-900/20 border border-purple-200 dark:border-purple-900/50">
                <div className="flex items-center gap-3 mb-4">
                  <div className="p-2 rounded-lg bg-purple-600 dark:bg-purple-500">
                    <Brain size={20} className="text-white" />
                  </div>
                  <div>
                    <h3 className="text-xl font-bold font-mono text-gray-900 dark:text-white">Good Creator Co</h3>
                    <p className="text-xs font-mono text-gray-600 dark:text-gray-400">Previous Company • 2021-2024</p>
                  </div>
                </div>
                <div className="space-y-2 mb-4">
                  <p className="text-sm text-gray-700 dark:text-gray-300">
                    6 projects • Python • Go • ClickHouse • ML
                  </p>
                  <p className="text-sm text-gray-600 dark:text-gray-400">
                    Distributed systems processing 10K+ events/sec
                  </p>
                </div>
                <div className="grid grid-cols-2 gap-3">
                  <div className="p-3 rounded-lg bg-white/50 dark:bg-gray-900/50">
                    <div className="text-2xl font-bold font-mono text-purple-600 dark:text-purple-400">14</div>
                    <div className="text-xs font-mono text-gray-600 dark:text-gray-400">Documents</div>
                  </div>
                  <div className="p-3 rounded-lg bg-white/50 dark:bg-gray-900/50">
                    <div className="text-2xl font-bold font-mono text-purple-600 dark:text-purple-400">6</div>
                    <div className="text-xs font-mono text-gray-600 dark:text-gray-400">Projects</div>
                  </div>
                </div>
              </div>
            </div>

            {/* Stats grid */}
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
              {[
                { label: 'Interview Questions', value: '140+', icon: MessageSquare },
                { label: 'Code Examples', value: '50+', icon: Code },
                { label: 'STAR Stories', value: '100+', icon: Award },
                { label: 'Reading Time', value: '12h', icon: Clock }
              ].map((stat) => {
                const StatIcon = stat.icon
                return (
                  <div key={stat.label} className="p-4 rounded-xl bg-white dark:bg-gray-900 border border-gray-200 dark:border-gray-800">
                    <StatIcon size={20} className="text-gray-400 mb-2" />
                    <div className="text-3xl font-bold font-mono text-gray-900 dark:text-white mb-1">
                      {stat.value}
                    </div>
                    <div className="text-xs font-mono text-gray-600 dark:text-gray-400">
                      {stat.label}
                    </div>
                  </div>
                )
              })}
            </div>
          </div>
        ) : (
          /* Document view */
          <div className="flex">
            {/* Main document */}
            <div className="flex-1 min-w-0">
              <div className="max-w-4xl mx-auto px-6 py-8">
                {/* Document header */}
                <div className="mb-8">
                  {/* Breadcrumb */}
                  <div className="flex items-center gap-2 text-xs font-mono text-gray-500 dark:text-gray-400 mb-4">
                    <button
                      onClick={() => setActiveDoc(null)}
                      className="hover:text-gray-900 dark:hover:text-white transition-colors"
                    >
                      Home
                    </button>
                    <ChevronRight size={12} />
                    <span>{currentDoc?.category}</span>
                    <ChevronRight size={12} />
                    <span className="text-gray-900 dark:text-white font-semibold">{currentDoc?.title}</span>
                  </div>

                  {/* Title */}
                  <h1 className="text-4xl md:text-5xl font-bold text-gray-900 dark:text-white mb-4 leading-tight">
                    {currentDoc?.title}
                  </h1>

                  {/* Metadata */}
                  {metadata && (
                    <div className="flex flex-wrap items-center gap-4 text-sm">
                      <div className="flex items-center gap-2 px-3 py-1.5 rounded-lg bg-gray-100 dark:bg-gray-800 font-mono">
                        <Building2 size={14} />
                        <span className="font-semibold">{metadata.company}</span>
                      </div>
                      <div className="flex items-center gap-2 text-gray-600 dark:text-gray-400 font-mono">
                        <Clock size={14} />
                        <span>{metadata.readingTime} min read</span>
                      </div>
                      <div className="flex items-center gap-2 text-gray-600 dark:text-gray-400 font-mono">
                        <FileText size={14} />
                        <span>{metadata.sections} sections</span>
                      </div>
                      {currentDoc?.badge && (
                        <span className="px-3 py-1.5 rounded-lg bg-emerald-100 dark:bg-emerald-900/30 text-emerald-700 dark:text-emerald-400 text-xs font-mono font-semibold">
                          {currentDoc.badge}
                        </span>
                      )}
                    </div>
                  )}

                  {/* Tech stack tags */}
                  {metadata && metadata.techStack.length > 0 && (
                    <div className="flex flex-wrap gap-2 mt-4">
                      {metadata.techStack.map((tech) => (
                        <span
                          key={tech}
                          className="px-2 py-1 rounded text-[10px] font-mono bg-blue-50 dark:bg-blue-950/30 text-blue-700 dark:text-blue-400 uppercase tracking-wider"
                        >
                          {tech}
                        </span>
                      ))}
                    </div>
                  )}
                </div>

                {/* Markdown content */}
                <article className="prose prose-gray dark:prose-invert max-w-none prose-headings:scroll-mt-24 prose-headings:font-bold prose-h1:text-3xl prose-h2:text-2xl prose-h2:mt-12 prose-h2:mb-4 prose-h2:pb-2 prose-h2:border-b prose-h2:border-gray-200 dark:prose-h2:border-gray-800 prose-h3:text-xl prose-h3:mt-8 prose-a:text-blue-600 dark:prose-a:text-blue-400 prose-a:no-underline hover:prose-a:underline prose-code:text-sm prose-code:bg-gray-100 dark:prose-code:bg-gray-800 prose-code:px-1.5 prose-code:py-0.5 prose-code:rounded prose-code:before:content-none prose-code:after:content-none prose-pre:bg-white dark:prose-pre:bg-gray-950 prose-pre:border prose-pre:border-gray-200 dark:prose-pre:border-gray-800 prose-pre:text-gray-800 dark:prose-pre:text-gray-200">
                  <ReactMarkdown
                    remarkPlugins={[remarkGfm]}
                    components={{
                      code({ node, inline, className, children, ...props }) {
                        const match = /language-(\w+)/.exec(className || '')
                        return !inline && match ? (
                          <SyntaxHighlighter
                            style={darkMode ? oneDark : oneLight}
                            language={match[1]}
                            PreTag="div"
                            showLineNumbers={false}
                            wrapLines={true}
                            customStyle={{
                              margin: 0,
                              borderRadius: '0.75rem',
                              fontSize: '0.875rem',
                              background: darkMode ? '#1f2937' : '#ffffff',
                              padding: '1.25rem',
                              border: darkMode ? 'none' : '1px solid #e5e7eb'
                            }}
                            codeTagProps={{
                              style: {
                                background: 'transparent',
                                fontFamily: 'var(--font-mono), monospace'
                              }
                            }}
                            {...props}
                          >
                            {String(children).replace(/\n$/, '')}
                          </SyntaxHighlighter>
                        ) : (
                          <code className={className} {...props}>
                            {children}
                          </code>
                        )
                      },
                      table({ children }) {
                        return (
                          <div className="overflow-x-auto my-6 rounded-lg border border-gray-200 dark:border-gray-800">
                            <table className="min-w-full">{children}</table>
                          </div>
                        )
                      },
                      h2({ children, ...props }) {
                        const id = String(children).toLowerCase().replace(/\s+/g, '-').replace(/[^\w-]/g, '')
                        return <h2 id={id} {...props}>{children}</h2>
                      },
                      h3({ children, ...props }) {
                        const id = String(children).toLowerCase().replace(/\s+/g, '-').replace(/[^\w-]/g, '')
                        return <h3 id={id} {...props}>{children}</h3>
                      }
                    }}
                  >
                    {preprocessMarkdown(currentDoc?.content || '')}
                  </ReactMarkdown>
                </article>

                {/* Navigation footer */}
                <div className="mt-16 pt-8 border-t border-gray-200 dark:border-gray-800">
                  <div className="flex items-center justify-between">
                    {(() => {
                      const allDocs = documents
                      const currentIndex = allDocs.findIndex(d => d.id === activeDoc)
                      const prevDoc = currentIndex > 0 ? allDocs[currentIndex - 1] : null
                      const nextDoc = currentIndex < allDocs.length - 1 ? allDocs[currentIndex + 1] : null

                      return (
                        <>
                          {prevDoc ? (
                            <button
                              onClick={() => handleDocSelect(prevDoc.id)}
                              className="group flex items-center gap-3 p-4 rounded-xl hover:bg-gray-100 dark:hover:bg-gray-800 transition-all"
                            >
                              <ChevronRight size={20} className="rotate-180 text-gray-400 group-hover:text-blue-500 transition-colors" />
                              <div className="text-left">
                                <div className="text-xs font-mono text-gray-500 dark:text-gray-400 mb-1">Previous</div>
                                <div className="text-sm font-semibold text-gray-900 dark:text-white group-hover:text-blue-600 dark:group-hover:text-blue-400 transition-colors">
                                  {prevDoc.title}
                                </div>
                              </div>
                            </button>
                          ) : <div />}
                          {nextDoc ? (
                            <button
                              onClick={() => handleDocSelect(nextDoc.id)}
                              className="group flex items-center gap-3 p-4 rounded-xl hover:bg-gray-100 dark:hover:bg-gray-800 transition-all"
                            >
                              <div className="text-right">
                                <div className="text-xs font-mono text-gray-500 dark:text-gray-400 mb-1">Next</div>
                                <div className="text-sm font-semibold text-gray-900 dark:text-white group-hover:text-blue-600 dark:group-hover:text-blue-400 transition-colors">
                                  {nextDoc.title}
                                </div>
                              </div>
                              <ChevronRight size={20} className="text-gray-400 group-hover:text-blue-500 transition-colors" />
                            </button>
                          ) : <div />}
                        </>
                      )
                    })()}
                  </div>
                </div>
              </div>
            </div>

            {/* TOC Sidebar (Desktop only) */}
            {currentDocTOC.length > 3 && (
              <div className="hidden xl:block w-64 shrink-0">
                <div className="sticky top-20 p-6">
                  <button
                    onClick={() => setTocExpanded(!tocExpanded)}
                    className="flex items-center justify-between w-full mb-4 text-xs font-mono font-bold text-gray-500 dark:text-gray-400 uppercase tracking-wider hover:text-gray-900 dark:hover:text-white transition-colors"
                  >
                    <span>On This Page</span>
                    <ChevronDown
                      size={14}
                      className={`transition-transform ${tocExpanded ? 'rotate-0' : '-rotate-90'}`}
                    />
                  </button>
                  {tocExpanded && (
                    <nav className="space-y-2">
                      {currentDocTOC.map((heading, index) => (
                        <a
                          key={index}
                          href={`#${heading.id}`}
                          className={`block text-xs hover:text-blue-600 dark:hover:text-blue-400 transition-colors ${
                            heading.level === 2
                              ? 'font-semibold text-gray-900 dark:text-white'
                              : 'pl-3 text-gray-600 dark:text-gray-400'
                          }`}
                          onClick={(e) => {
                            e.preventDefault()
                            document.getElementById(heading.id)?.scrollIntoView({
                              behavior: 'smooth',
                              block: 'start'
                            })
                          }}
                        >
                          {heading.text}
                        </a>
                      ))}
                    </nav>
                  )}
                </div>
              </div>
            )}
          </div>
        )}

        {/* Footer */}
        <footer className="border-t border-gray-200 dark:border-gray-800 mt-16 py-8 text-center">
          <p className="text-sm font-mono text-gray-600 dark:text-gray-400">
            Comprehensive Google L4 Interview Preparation
          </p>
          <p className="text-xs font-mono text-gray-500 dark:text-gray-500 mt-2">
            {documents.length} documents • 140+ questions • 2 companies • Good luck! 🚀
          </p>
        </footer>
      </main>
    </div>
  )
}

export default App

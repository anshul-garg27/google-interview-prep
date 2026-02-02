import { useState, useEffect, useMemo } from 'react'
import ReactMarkdown from 'react-markdown'
import remarkGfm from 'remark-gfm'
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter'
import { oneDark, oneLight } from 'react-syntax-highlighter/dist/esm/styles/prism'
import Fuse from 'fuse.js'
import {
  Menu, X, Search, Moon, Sun, ChevronRight,
  BookOpen, FileText, Code, Brain, Network,
  Target, MessageSquare, Briefcase, Github
} from 'lucide-react'
import { documents } from './data'

function App() {
  const [darkMode, setDarkMode] = useState(() => {
    if (typeof window !== 'undefined') {
      return localStorage.getItem('darkMode') === 'true' ||
             window.matchMedia('(prefers-color-scheme: dark)').matches
    }
    return false
  })
  const [sidebarOpen, setSidebarOpen] = useState(false)
  const [activeDoc, setActiveDoc] = useState('GOOGLEYNESS_ALL_QUESTIONS')
  const [searchQuery, setSearchQuery] = useState('')
  const [searchResults, setSearchResults] = useState([])
  const [showSearch, setShowSearch] = useState(false)

  // Apply dark mode
  useEffect(() => {
    document.documentElement.classList.toggle('dark', darkMode)
    localStorage.setItem('darkMode', darkMode)
  }, [darkMode])

  // Search functionality
  const fuse = useMemo(() => {
    const searchData = documents.map(doc => ({
      id: doc.id,
      title: doc.title,
      content: doc.content,
      category: doc.category
    }))
    return new Fuse(searchData, {
      keys: ['title', 'content'],
      threshold: 0.4,
      includeMatches: true
    })
  }, [])

  useEffect(() => {
    if (searchQuery.length > 2) {
      const results = fuse.search(searchQuery).slice(0, 10)
      setSearchResults(results)
    } else {
      setSearchResults([])
    }
  }, [searchQuery, fuse])

  const currentDoc = documents.find(d => d.id === activeDoc)

  const categories = [
    {
      name: 'Interview Prep',
      icon: Target,
      docs: documents.filter(d => d.category === 'interview')
    },
    {
      name: 'Technical Docs',
      icon: Code,
      docs: documents.filter(d => d.category === 'technical')
    },
    {
      name: 'Project Analysis',
      icon: FileText,
      docs: documents.filter(d => d.category === 'analysis')
    }
  ]

  const handleDocSelect = (docId) => {
    setActiveDoc(docId)
    setSidebarOpen(false)
    setShowSearch(false)
    setSearchQuery('')
    window.scrollTo(0, 0)
  }

  return (
    <div className={`min-h-screen bg-gray-50 dark:bg-gray-900 transition-colors duration-300`}>
      {/* Header */}
      <header className="fixed top-0 left-0 right-0 z-50 bg-white dark:bg-gray-800 shadow-sm border-b border-gray-200 dark:border-gray-700">
        <div className="flex items-center justify-between px-4 h-16">
          {/* Mobile menu button */}
          <button
            onClick={() => setSidebarOpen(!sidebarOpen)}
            className="lg:hidden p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700"
          >
            {sidebarOpen ? <X size={24} /> : <Menu size={24} />}
          </button>

          {/* Logo */}
          <div className="flex items-center gap-3">
            <div className="flex gap-1">
              <span className="w-3 h-3 rounded-full bg-blue-500"></span>
              <span className="w-3 h-3 rounded-full bg-red-500"></span>
              <span className="w-3 h-3 rounded-full bg-yellow-500"></span>
              <span className="w-3 h-3 rounded-full bg-green-500"></span>
            </div>
            <h1 className="text-xl font-bold text-gray-800 dark:text-white hidden sm:block">
              Google L4 Interview Prep
            </h1>
          </div>

          {/* Right side actions */}
          <div className="flex items-center gap-2">
            {/* Search button */}
            <button
              onClick={() => setShowSearch(!showSearch)}
              className="p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700"
            >
              <Search size={20} className="text-gray-600 dark:text-gray-300" />
            </button>

            {/* Theme toggle */}
            <button
              onClick={() => setDarkMode(!darkMode)}
              className="p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700"
            >
              {darkMode ? (
                <Sun size={20} className="text-yellow-500" />
              ) : (
                <Moon size={20} className="text-gray-600" />
              )}
            </button>

            {/* GitHub link */}
            <a
              href="https://github.com/anshul-garg27/google-interview-prep"
              target="_blank"
              rel="noopener noreferrer"
              className="p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700"
            >
              <Github size={20} className="text-gray-600 dark:text-gray-300" />
            </a>
          </div>
        </div>

        {/* Search bar */}
        {showSearch && (
          <div className="px-4 pb-4 bg-white dark:bg-gray-800">
            <div className="relative">
              <Search size={18} className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" />
              <input
                type="text"
                placeholder="Search questions, topics..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="w-full pl-10 pr-4 py-2 rounded-lg border border-gray-300 dark:border-gray-600 bg-gray-50 dark:bg-gray-700 text-gray-800 dark:text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
                autoFocus
              />
            </div>

            {/* Search results */}
            {searchResults.length > 0 && (
              <div className="mt-2 bg-white dark:bg-gray-700 rounded-lg shadow-lg border border-gray-200 dark:border-gray-600 max-h-64 overflow-y-auto">
                {searchResults.map((result) => (
                  <button
                    key={result.item.id}
                    onClick={() => handleDocSelect(result.item.id)}
                    className="w-full text-left px-4 py-3 hover:bg-gray-100 dark:hover:bg-gray-600 border-b border-gray-100 dark:border-gray-600 last:border-0"
                  >
                    <div className="font-medium text-gray-800 dark:text-white">
                      {result.item.title}
                    </div>
                    <div className="text-sm text-gray-500 dark:text-gray-400">
                      {result.item.category}
                    </div>
                  </button>
                ))}
              </div>
            )}
          </div>
        )}
      </header>

      {/* Sidebar overlay for mobile */}
      {sidebarOpen && (
        <div
          className="fixed inset-0 z-40 bg-black/50 lg:hidden"
          onClick={() => setSidebarOpen(false)}
        />
      )}

      {/* Sidebar */}
      <aside
        className={`fixed top-16 left-0 z-40 w-72 h-[calc(100vh-4rem)] bg-white dark:bg-gray-800 border-r border-gray-200 dark:border-gray-700 overflow-y-auto transition-transform duration-300 ${
          sidebarOpen ? 'translate-x-0' : '-translate-x-full'
        } lg:translate-x-0`}
      >
        <nav className="p-4">
          {categories.map((category) => (
            <div key={category.name} className="mb-6">
              <div className="flex items-center gap-2 px-3 py-2 text-sm font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                <category.icon size={16} />
                {category.name}
              </div>
              <ul className="mt-1 space-y-1">
                {category.docs.map((doc) => (
                  <li key={doc.id}>
                    <button
                      onClick={() => handleDocSelect(doc.id)}
                      className={`w-full text-left px-3 py-2 rounded-lg text-sm transition-colors ${
                        activeDoc === doc.id
                          ? 'bg-blue-50 dark:bg-blue-900/30 text-blue-600 dark:text-blue-400 font-medium'
                          : 'text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700'
                      }`}
                    >
                      <div className="flex items-center gap-2">
                        <ChevronRight size={14} className={activeDoc === doc.id ? 'text-blue-500' : 'text-gray-400'} />
                        <span className="truncate">{doc.title}</span>
                      </div>
                      {doc.badge && (
                        <span className="ml-6 text-xs bg-green-100 dark:bg-green-900 text-green-700 dark:text-green-300 px-2 py-0.5 rounded-full">
                          {doc.badge}
                        </span>
                      )}
                    </button>
                  </li>
                ))}
              </ul>
            </div>
          ))}
        </nav>

        {/* Quick stats */}
        <div className="p-4 border-t border-gray-200 dark:border-gray-700">
          <div className="grid grid-cols-2 gap-3">
            <div className="bg-blue-50 dark:bg-blue-900/30 rounded-lg p-3 text-center">
              <div className="text-2xl font-bold text-blue-600 dark:text-blue-400">60+</div>
              <div className="text-xs text-gray-600 dark:text-gray-400">Questions</div>
            </div>
            <div className="bg-green-50 dark:bg-green-900/30 rounded-lg p-3 text-center">
              <div className="text-2xl font-bold text-green-600 dark:text-green-400">6</div>
              <div className="text-xs text-gray-600 dark:text-gray-400">Projects</div>
            </div>
          </div>
        </div>
      </aside>

      {/* Main content */}
      <main className="lg:ml-72 pt-16 min-h-screen">
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
          {/* Breadcrumb */}
          <div className="flex items-center gap-2 text-sm text-gray-500 dark:text-gray-400 mb-6">
            <BookOpen size={16} />
            <span>{currentDoc?.category}</span>
            <ChevronRight size={14} />
            <span className="text-gray-800 dark:text-white font-medium">{currentDoc?.title}</span>
          </div>

          {/* Document title */}
          <h1 className="text-3xl sm:text-4xl font-bold text-gray-900 dark:text-white mb-8">
            {currentDoc?.title}
          </h1>

          {/* Markdown content */}
          <article className="prose prose-gray dark:prose-invert max-w-none prose-headings:scroll-mt-20 prose-a:text-blue-600 dark:prose-a:text-blue-400 prose-pre:bg-gray-900 dark:prose-pre:bg-gray-950">
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
                      {...props}
                    >
                      {String(children).replace(/\n$/, '')}
                    </SyntaxHighlighter>
                  ) : (
                    <code className={`${className} bg-gray-100 dark:bg-gray-800 px-1.5 py-0.5 rounded text-sm`} {...props}>
                      {children}
                    </code>
                  )
                },
                table({ children }) {
                  return (
                    <div className="overflow-x-auto my-4">
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
              {currentDoc?.content || ''}
            </ReactMarkdown>
          </article>

          {/* Navigation footer */}
          <div className="mt-12 pt-8 border-t border-gray-200 dark:border-gray-700">
            <div className="flex justify-between">
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
                        className="flex items-center gap-2 text-blue-600 dark:text-blue-400 hover:underline"
                      >
                        <ChevronRight size={16} className="rotate-180" />
                        <span className="hidden sm:inline">{prevDoc.title}</span>
                        <span className="sm:hidden">Previous</span>
                      </button>
                    ) : <div />}
                    {nextDoc ? (
                      <button
                        onClick={() => handleDocSelect(nextDoc.id)}
                        className="flex items-center gap-2 text-blue-600 dark:text-blue-400 hover:underline"
                      >
                        <span className="hidden sm:inline">{nextDoc.title}</span>
                        <span className="sm:hidden">Next</span>
                        <ChevronRight size={16} />
                      </button>
                    ) : <div />}
                  </>
                )
              })()}
            </div>
          </div>
        </div>

        {/* Footer */}
        <footer className="border-t border-gray-200 dark:border-gray-700 mt-12 py-8 text-center text-gray-500 dark:text-gray-400">
          <p>Google L4 Interview Preparation Materials</p>
          <p className="text-sm mt-2">Good luck with your interview! 🚀</p>
        </footer>
      </main>
    </div>
  )
}

export default App

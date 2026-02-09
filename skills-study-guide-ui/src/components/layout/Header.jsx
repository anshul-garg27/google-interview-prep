import { useState, useEffect, useCallback } from 'react'
import { useNavigate } from 'react-router-dom'
import { Menu, X, Search, Moon, Sun, BookOpen } from 'lucide-react'
import { useTheme } from '../../context/ThemeContext'
import SearchModal from '../search/SearchModal'

export default function Header({ onMenuToggle, sidebarOpen }) {
  const { darkMode, toggleDarkMode } = useTheme()
  const [searchOpen, setSearchOpen] = useState(false)
  const navigate = useNavigate()

  // Cmd+K shortcut
  const handleKeyDown = useCallback((e) => {
    if ((e.metaKey || e.ctrlKey) && e.key === 'k') {
      e.preventDefault()
      setSearchOpen(true)
    }
    if (e.key === 'Escape') {
      setSearchOpen(false)
    }
  }, [])

  useEffect(() => {
    document.addEventListener('keydown', handleKeyDown)
    return () => document.removeEventListener('keydown', handleKeyDown)
  }, [handleKeyDown])

  return (
    <>
      <header className="fixed top-0 left-0 right-0 h-14 z-40 bg-[var(--bg-surface)]/80 backdrop-blur-xl border-b border-[var(--border)]">
        <div className="flex items-center justify-between h-full px-4 lg:px-6">
          {/* Left: menu + logo */}
          <div className="flex items-center gap-3">
            <button
              onClick={onMenuToggle}
              className="lg:hidden p-2 rounded-lg hover:bg-[var(--bg-hover)] transition-colors"
              aria-label="Toggle menu"
            >
              {sidebarOpen ? <X size={18} /> : <Menu size={18} />}
            </button>

            <button
              onClick={() => navigate('/')}
              className="flex items-center gap-2 hover:opacity-80 transition-opacity"
            >
              <BookOpen size={18} className="text-[var(--accent)]" />
              <span className="font-heading font-semibold text-base tracking-tight">
                Study Guides
              </span>
            </button>
          </div>

          {/* Right: search + theme */}
          <div className="flex items-center gap-2">
            <button
              onClick={() => setSearchOpen(true)}
              className="flex items-center gap-2 px-3 py-1.5 rounded-lg border border-[var(--border)] hover:bg-[var(--bg-hover)] transition-colors text-[var(--text-tertiary)]"
            >
              <Search size={14} />
              <span className="hidden sm:inline text-xs font-mono">Search</span>
              <kbd className="hidden sm:inline text-[10px] font-mono px-1.5 py-0.5 rounded bg-[var(--bg-surface-alt)] border border-[var(--border)]">
                {/Mac|iPhone|iPad/.test(navigator.userAgent) ? '⌘' : 'Ctrl'}K
              </kbd>
            </button>

            <button
              onClick={toggleDarkMode}
              className="p-2 rounded-lg hover:bg-[var(--bg-hover)] transition-colors"
              aria-label="Toggle dark mode"
            >
              {darkMode ? <Sun size={16} /> : <Moon size={16} />}
            </button>
          </div>
        </div>
      </header>

      {searchOpen && <SearchModal onClose={() => setSearchOpen(false)} />}
    </>
  )
}

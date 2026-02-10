import { useState, useEffect, useRef, useCallback } from 'react'
import { useLocation } from 'react-router-dom'
import Header from './Header'
import Sidebar from './Sidebar'

const SWIPE_THRESHOLD = 50  // minimum px to trigger swipe
const EDGE_ZONE = 30        // px from left edge to start swipe-to-open

export default function Layout({ children }) {
  const [sidebarOpen, setSidebarOpen] = useState(false)
  const location = useLocation()
  const touchRef = useRef({ startX: 0, startY: 0, startedFromEdge: false })

  // Close sidebar on route change (mobile)
  useEffect(() => {
    setSidebarOpen(false)
  }, [location.pathname])

  // Close sidebar on Escape key
  useEffect(() => {
    if (!sidebarOpen) return
    const handleEsc = (e) => {
      if (e.key === 'Escape') setSidebarOpen(false)
    }
    document.addEventListener('keydown', handleEsc)
    return () => document.removeEventListener('keydown', handleEsc)
  }, [sidebarOpen])

  // Swipe gestures for mobile sidebar
  const handleTouchStart = useCallback((e) => {
    const touch = e.touches[0]
    touchRef.current = {
      startX: touch.clientX,
      startY: touch.clientY,
      startedFromEdge: touch.clientX < EDGE_ZONE,
    }
  }, [])

  const handleTouchEnd = useCallback((e) => {
    const touch = e.changedTouches[0]
    const { startX, startY, startedFromEdge } = touchRef.current
    const deltaX = touch.clientX - startX
    const deltaY = touch.clientY - startY

    // Only horizontal swipes (not vertical scrolling)
    if (Math.abs(deltaX) < SWIPE_THRESHOLD || Math.abs(deltaY) > Math.abs(deltaX)) return

    // Don't trigger on desktop (sidebar is always visible on lg+)
    if (window.innerWidth >= 1024) return

    if (deltaX > 0 && startedFromEdge && !sidebarOpen) {
      // Swipe RIGHT from left edge → open sidebar
      setSidebarOpen(true)
    } else if (deltaX < 0 && sidebarOpen) {
      // Swipe LEFT while sidebar open → close sidebar
      setSidebarOpen(false)
    }
  }, [sidebarOpen])

  useEffect(() => {
    document.addEventListener('touchstart', handleTouchStart, { passive: true })
    document.addEventListener('touchend', handleTouchEnd, { passive: true })
    return () => {
      document.removeEventListener('touchstart', handleTouchStart)
      document.removeEventListener('touchend', handleTouchEnd)
    }
  }, [handleTouchStart, handleTouchEnd])

  return (
    <div className="min-h-screen bg-[var(--bg-primary)]">
      <Header
        onMenuToggle={() => setSidebarOpen(prev => !prev)}
        sidebarOpen={sidebarOpen}
      />

      {/* Sidebar */}
      <Sidebar open={sidebarOpen} onClose={() => setSidebarOpen(false)} />

      {/* Mobile backdrop — also swipeable to close */}
      {sidebarOpen && (
        <div
          className="fixed inset-0 bg-black/40 backdrop-blur-sm z-30 lg:hidden"
          onClick={() => setSidebarOpen(false)}
          aria-hidden="true"
        />
      )}

      {/* Main content */}
      <main className="lg:pl-[260px] pt-14 min-h-screen">
        {children}
      </main>
    </div>
  )
}

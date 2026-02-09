import { useState, useEffect, useCallback, useMemo } from 'react'
import ReactMarkdown from 'react-markdown'
import remarkGfm from 'remark-gfm'
import { RotateCcw, Eye, EyeOff, ChevronRight, Shuffle, Brain, ChevronDown, ChevronsUpDown, Filter } from 'lucide-react'
import { useTheme } from '../../context/ThemeContext'

const STORAGE_KEY = 'sg-flashcards-kafka'

// ---------- SM-2 Algorithm ----------
function sm2(quality, repetition, efactor, interval) {
  if (quality >= 3) {
    if (repetition === 0) interval = 1
    else if (repetition === 1) interval = 6
    else interval = Math.round(interval * efactor)
    repetition++
  } else {
    repetition = 0
    interval = 1
  }
  efactor = Math.max(1.3, efactor + (0.1 - (5 - quality) * (0.08 + (5 - quality) * 0.02)))
  return { interval, repetition, efactor }
}

function getDefaultCardState() {
  return { interval: 0, repetition: 0, efactor: 2.5, nextReview: null, lastReviewed: null }
}

function loadCardStates() {
  try {
    const raw = localStorage.getItem(STORAGE_KEY)
    return raw ? JSON.parse(raw) : {}
  } catch {
    return {}
  }
}

function saveCardStates(states) {
  try {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(states))
  } catch { /* quota exceeded, silently fail */ }
}

function todayStr() {
  return new Date().toISOString().slice(0, 10)
}

function isDue(cardState) {
  if (!cardState || !cardState.nextReview) return true
  return cardState.nextReview <= todayStr()
}

function isMastered(cardState) {
  return cardState && cardState.repetition >= 3 && cardState.interval >= 6
}

function isNew(cardState) {
  return !cardState || cardState.lastReviewed === null
}

// ---------- Difficulty Buttons ----------
const DIFFICULTIES = [
  { label: 'Again', quality: 0, color: 'red' },
  { label: 'Hard', quality: 2, color: 'orange' },
  { label: 'Good', quality: 3, color: 'green' },
  { label: 'Easy', quality: 5, color: 'blue' },
]

const difficultyStyles = {
  red: {
    border: 'border-red-400 dark:border-red-500',
    hover: 'hover:bg-red-50 dark:hover:bg-red-500/10',
    text: 'text-red-600 dark:text-red-400',
    activeBg: 'bg-red-50 dark:bg-red-500/10',
  },
  orange: {
    border: 'border-orange-400 dark:border-orange-500',
    hover: 'hover:bg-orange-50 dark:hover:bg-orange-500/10',
    text: 'text-orange-600 dark:text-orange-400',
    activeBg: 'bg-orange-50 dark:bg-orange-500/10',
  },
  green: {
    border: 'border-green-400 dark:border-green-500',
    hover: 'hover:bg-green-50 dark:hover:bg-green-500/10',
    text: 'text-green-600 dark:text-green-400',
    activeBg: 'bg-green-50 dark:bg-green-500/10',
  },
  blue: {
    border: 'border-blue-400 dark:border-blue-500',
    hover: 'hover:bg-blue-50 dark:hover:bg-blue-500/10',
    text: 'text-blue-600 dark:text-blue-400',
    activeBg: 'bg-blue-50 dark:bg-blue-500/10',
  },
}

// ---------- Shuffle helper ----------
function shuffleArray(arr) {
  const a = [...arr]
  for (let i = a.length - 1; i > 0; i--) {
    const j = Math.floor(Math.random() * (i + 1))
    ;[a[i], a[j]] = [a[j], a[i]]
  }
  return a
}

// ---------- Main Component ----------
export default function FlashcardDeck({ items = [] }) {
  const { darkMode } = useTheme()
  const [mode, setMode] = useState('quiz') // 'quiz' | 'study'
  const [cardStates, setCardStates] = useState(loadCardStates)

  // Quiz mode state
  const [order, setOrder] = useState(() => items.map((_, i) => i))
  const [currentIdx, setCurrentIdx] = useState(0)
  const [flipped, setFlipped] = useState(false)
  const [rated, setRated] = useState(false)
  const [reviewedToday, setReviewedToday] = useState(0)

  // Study mode state
  const [expandedItems, setExpandedItems] = useState(new Set())
  const [studyFilter, setStudyFilter] = useState('all') // 'all' | 'due' | 'new' | 'mastered'

  // Persist card states whenever they change
  useEffect(() => {
    saveCardStates(cardStates)
  }, [cardStates])

  // Sort order: due cards first, then by next review date
  const sortByDue = useCallback((indices) => {
    return [...indices].sort((a, b) => {
      const stA = cardStates[a]
      const stB = cardStates[b]
      const dueA = isDue(stA)
      const dueB = isDue(stB)
      if (dueA && !dueB) return -1
      if (!dueA && dueB) return 1
      const dateA = stA?.nextReview || '0000-00-00'
      const dateB = stB?.nextReview || '0000-00-00'
      if (dateA < dateB) return -1
      if (dateA > dateB) return 1
      return 0
    })
  }, [cardStates])

  // Initialize order sorted by due date
  useEffect(() => {
    setOrder(sortByDue(items.map((_, i) => i)))
    setCurrentIdx(0)
    setFlipped(false)
    setRated(false)
  }, [items.length]) // eslint-disable-line react-hooks/exhaustive-deps

  const handleShuffle = () => {
    setOrder(shuffleArray(items.map((_, i) => i)))
    setCurrentIdx(0)
    setFlipped(false)
    setRated(false)
  }

  const handleFlip = () => {
    if (!flipped) setFlipped(true)
  }

  const handleRate = (quality) => {
    const cardIdx = order[currentIdx]
    const prev = cardStates[cardIdx] || getDefaultCardState()
    const result = sm2(quality, prev.repetition, prev.efactor, prev.interval)
    const today = todayStr()
    const nextDate = new Date()
    nextDate.setDate(nextDate.getDate() + result.interval)

    setCardStates(prev => ({
      ...prev,
      [cardIdx]: {
        interval: result.interval,
        repetition: result.repetition,
        efactor: result.efactor,
        nextReview: nextDate.toISOString().slice(0, 10),
        lastReviewed: today,
      }
    }))
    setRated(true)
    setReviewedToday(c => c + 1)
  }

  const handleNext = () => {
    if (currentIdx < order.length - 1) {
      setCurrentIdx(currentIdx + 1)
    } else {
      // Wrap around to start, re-sort by due
      setOrder(sortByDue(items.map((_, i) => i)))
      setCurrentIdx(0)
    }
    setFlipped(false)
    setRated(false)
  }

  const handleReset = () => {
    setCardStates({})
    setReviewedToday(0)
    setOrder(items.map((_, i) => i))
    setCurrentIdx(0)
    setFlipped(false)
    setRated(false)
  }

  // Study mode helpers
  const toggleExpand = (idx) => {
    setExpandedItems(prev => {
      const next = new Set(prev)
      if (next.has(idx)) next.delete(idx)
      else next.add(idx)
      return next
    })
  }

  const expandAll = () => setExpandedItems(new Set(items.map((_, i) => i)))
  const collapseAll = () => setExpandedItems(new Set())

  const filteredStudyItems = useMemo(() => {
    return items.map((item, idx) => ({ item, idx })).filter(({ idx }) => {
      const st = cardStates[idx]
      switch (studyFilter) {
        case 'due': return isDue(st)
        case 'new': return isNew(st)
        case 'mastered': return isMastered(st)
        default: return true
      }
    })
  }, [items, cardStates, studyFilter])

  // Stats
  const totalCards = items.length
  const dueCount = items.filter((_, i) => isDue(cardStates[i])).length
  const newCount = items.filter((_, i) => isNew(cardStates[i])).length
  const masteredCount = items.filter((_, i) => isMastered(cardStates[i])).length
  const progressPct = totalCards > 0 ? Math.round(((totalCards - newCount) / totalCards) * 100) : 0

  if (!items.length) return null

  const currentCardIdx = order[currentIdx]
  const currentItem = items[currentCardIdx]

  return (
    <div className="my-8">
      {/* Progress bar */}
      <div className="h-1.5 rounded-full bg-[var(--bg-surface-alt)] mb-6 overflow-hidden">
        <div
          className="h-full rounded-full bg-[var(--accent)] transition-all duration-500 ease-out"
          style={{ width: `${progressPct}%` }}
        />
      </div>

      {/* Header: Mode toggle + stats */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 mb-6">
        {/* Mode toggle */}
        <div className="inline-flex rounded-lg border border-[var(--border)] bg-[var(--bg-surface)] p-0.5">
          <button
            onClick={() => setMode('quiz')}
            className={`flex items-center gap-1.5 px-4 py-2 rounded-md font-mono text-xs font-semibold uppercase tracking-wide transition-colors ${
              mode === 'quiz'
                ? 'bg-[var(--accent)] text-white'
                : 'text-[var(--text-secondary)] hover:text-[var(--text-primary)]'
            }`}
          >
            <Brain size={14} />
            Quiz Mode
          </button>
          <button
            onClick={() => setMode('study')}
            className={`flex items-center gap-1.5 px-4 py-2 rounded-md font-mono text-xs font-semibold uppercase tracking-wide transition-colors ${
              mode === 'study'
                ? 'bg-[var(--accent)] text-white'
                : 'text-[var(--text-secondary)] hover:text-[var(--text-primary)]'
            }`}
          >
            <Eye size={14} />
            Study Mode
          </button>
        </div>

        {/* Stats pills */}
        <div className="flex items-center gap-3 font-mono text-xs text-[var(--text-tertiary)]">
          <span>{dueCount} due</span>
          <span className="w-px h-3 bg-[var(--border)]" />
          <span>{newCount} new</span>
          <span className="w-px h-3 bg-[var(--border)]" />
          <span>{masteredCount} mastered</span>
        </div>
      </div>

      {mode === 'quiz' ? (
        <QuizMode
          currentItem={currentItem}
          currentIdx={currentIdx}
          totalCards={totalCards}
          flipped={flipped}
          rated={rated}
          reviewedToday={reviewedToday}
          onFlip={handleFlip}
          onRate={handleRate}
          onNext={handleNext}
          onShuffle={handleShuffle}
          onReset={handleReset}
        />
      ) : (
        <StudyMode
          filteredItems={filteredStudyItems}
          expandedItems={expandedItems}
          studyFilter={studyFilter}
          onSetFilter={setStudyFilter}
          onToggleExpand={toggleExpand}
          onExpandAll={expandAll}
          onCollapseAll={collapseAll}
          cardStates={cardStates}
          dueCount={dueCount}
          newCount={newCount}
          masteredCount={masteredCount}
          totalCards={totalCards}
        />
      )}
    </div>
  )
}

// ---------- Quiz Mode ----------
function QuizMode({
  currentItem, currentIdx, totalCards,
  flipped, rated, reviewedToday,
  onFlip, onRate, onNext, onShuffle, onReset
}) {
  return (
    <div>
      {/* Toolbar */}
      <div className="flex items-center justify-between mb-4">
        <span className="font-mono text-xs text-[var(--text-tertiary)]">
          Card {currentIdx + 1} of {totalCards}
          <span className="ml-3 text-[var(--accent)]">{reviewedToday} reviewed today</span>
        </span>
        <div className="flex items-center gap-2">
          <button
            onClick={onShuffle}
            className="flex items-center gap-1.5 px-3 py-1.5 rounded-lg border border-[var(--border)] text-[var(--text-secondary)] hover:bg-[var(--bg-hover)] font-mono text-xs transition-colors"
            title="Shuffle cards"
          >
            <Shuffle size={13} /> Shuffle
          </button>
          <button
            onClick={onReset}
            className="flex items-center gap-1.5 px-3 py-1.5 rounded-lg border border-[var(--border)] text-[var(--text-secondary)] hover:bg-[var(--bg-hover)] font-mono text-xs transition-colors"
            title="Reset all progress"
          >
            <RotateCcw size={13} /> Reset
          </button>
        </div>
      </div>

      {/* Flashcard */}
      <div
        onClick={!flipped ? onFlip : undefined}
        className={`relative min-h-[200px] rounded-xl border border-[var(--border)] bg-[var(--bg-surface)] shadow-sm transition-all duration-300 ${
          !flipped ? 'cursor-pointer hover:border-[var(--accent)] hover:shadow-md' : ''
        }`}
      >
        {/* Question face */}
        <div
          className="p-8 flex flex-col items-center justify-center min-h-[200px] transition-opacity duration-300"
          style={{ opacity: flipped ? 0 : 1, position: flipped ? 'absolute' : 'relative', inset: 0 }}
        >
          <span className="font-mono text-[10px] font-bold uppercase tracking-[0.1em] text-[var(--accent)] mb-4">
            Question
          </span>
          <h3 className="font-heading text-lg md:text-xl font-bold text-center leading-relaxed text-[var(--text-primary)]">
            {currentItem?.question}
          </h3>
          {!flipped && (
            <button
              onClick={onFlip}
              className="mt-6 flex items-center gap-1.5 px-4 py-2 rounded-lg border border-[var(--border)] text-[var(--text-secondary)] hover:bg-[var(--bg-hover)] font-mono text-xs transition-colors"
            >
              <Eye size={13} /> Show Answer
            </button>
          )}
        </div>

        {/* Answer face */}
        <div
          className="p-8 flex flex-col items-center justify-center min-h-[200px] transition-opacity duration-300"
          style={{ opacity: flipped ? 1 : 0, position: flipped ? 'relative' : 'absolute', inset: 0 }}
        >
          <span className="font-mono text-[10px] font-bold uppercase tracking-[0.1em] text-green-600 dark:text-green-400 mb-4">
            Answer
          </span>
          <div className="w-full max-w-2xl prose prose-sm prose-gray dark:prose-invert font-body text-center [&>p]:text-[var(--text-primary)] [&>ul]:text-left [&>ol]:text-left">
            <ReactMarkdown remarkPlugins={[remarkGfm]}>
              {currentItem?.answer || ''}
            </ReactMarkdown>
          </div>
        </div>
      </div>

      {/* Difficulty buttons + Next */}
      {flipped && (
        <div className="mt-5 flex flex-col items-center gap-4">
          {!rated ? (
            <div className="flex items-center gap-3">
              <span className="font-mono text-[10px] uppercase tracking-wide text-[var(--text-tertiary)] mr-2">
                How well did you know this?
              </span>
              {DIFFICULTIES.map(d => {
                const s = difficultyStyles[d.color]
                return (
                  <button
                    key={d.label}
                    onClick={() => onRate(d.quality)}
                    className={`px-4 py-2 rounded-lg border-2 ${s.border} ${s.text} ${s.hover} font-mono text-xs font-semibold transition-colors`}
                  >
                    {d.label}
                  </button>
                )
              })}
            </div>
          ) : (
            <button
              onClick={onNext}
              className="flex items-center gap-1.5 px-6 py-2.5 rounded-lg bg-[var(--accent)] hover:bg-[var(--accent-hover)] text-white font-mono text-xs font-semibold uppercase tracking-wide transition-colors"
            >
              Next <ChevronRight size={14} />
            </button>
          )}
        </div>
      )}
    </div>
  )
}

// ---------- Study Mode ----------
function StudyMode({
  filteredItems, expandedItems, studyFilter,
  onSetFilter, onToggleExpand, onExpandAll, onCollapseAll,
  cardStates, dueCount, newCount, masteredCount, totalCards
}) {
  const filterOptions = [
    { key: 'all', label: 'All', count: totalCards },
    { key: 'due', label: 'Due for Review', count: dueCount },
    { key: 'new', label: 'New', count: newCount },
    { key: 'mastered', label: 'Mastered', count: masteredCount },
  ]

  return (
    <div>
      {/* Filter bar */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3 mb-5">
        <div className="flex items-center gap-1.5 flex-wrap">
          <Filter size={13} className="text-[var(--text-tertiary)]" />
          {filterOptions.map(f => (
            <button
              key={f.key}
              onClick={() => onSetFilter(f.key)}
              className={`px-3 py-1.5 rounded-md font-mono text-xs transition-colors ${
                studyFilter === f.key
                  ? 'bg-[var(--accent)] text-white'
                  : 'text-[var(--text-secondary)] hover:bg-[var(--bg-hover)] border border-[var(--border)]'
              }`}
            >
              {f.label} ({f.count})
            </button>
          ))}
        </div>
        <button
          onClick={expandedItems.size > 0 ? onCollapseAll : onExpandAll}
          className="flex items-center gap-1.5 px-3 py-1.5 rounded-lg border border-[var(--border)] text-[var(--text-secondary)] hover:bg-[var(--bg-hover)] font-mono text-xs transition-colors"
        >
          <ChevronsUpDown size={13} />
          {expandedItems.size > 0 ? 'Collapse All' : 'Expand All'}
        </button>
      </div>

      {/* Accordion list */}
      <div className="space-y-2">
        {filteredItems.length === 0 && (
          <div className="text-center py-12 text-[var(--text-tertiary)] font-body">
            No cards match this filter.
          </div>
        )}
        {filteredItems.map(({ item, idx }) => {
          const expanded = expandedItems.has(idx)
          const st = cardStates[idx]
          const statusLabel = isMastered(st) ? 'mastered' : isNew(st) ? 'new' : isDue(st) ? 'due' : 'reviewed'
          const statusColor = isMastered(st)
            ? 'text-green-600 dark:text-green-400'
            : isNew(st)
            ? 'text-blue-500 dark:text-blue-400'
            : isDue(st)
            ? 'text-orange-500 dark:text-orange-400'
            : 'text-[var(--text-tertiary)]'

          return (
            <div
              key={idx}
              className="rounded-lg border border-[var(--border)] bg-[var(--bg-surface)] overflow-hidden transition-colors hover:border-[var(--border-strong)]"
            >
              <button
                onClick={() => onToggleExpand(idx)}
                className="w-full flex items-start gap-3 px-5 py-4 text-left"
              >
                <ChevronDown
                  size={16}
                  className={`mt-0.5 shrink-0 text-[var(--text-tertiary)] transition-transform duration-200 ${expanded ? 'rotate-180' : ''}`}
                />
                <span className="flex-1 font-body text-sm font-medium text-[var(--text-primary)] leading-relaxed">
                  {item.question}
                </span>
                <span className={`shrink-0 font-mono text-[10px] uppercase tracking-wide ${statusColor}`}>
                  {statusLabel}
                </span>
              </button>
              {expanded && (
                <div className="px-5 pb-5 pt-0 ml-7 border-t border-[var(--border)]">
                  <div className="pt-4 prose prose-sm prose-gray dark:prose-invert font-body [&>p]:text-[var(--text-primary)]">
                    <ReactMarkdown remarkPlugins={[remarkGfm]}>
                      {item.answer}
                    </ReactMarkdown>
                  </div>
                </div>
              )}
            </div>
          )
        })}
      </div>
    </div>
  )
}

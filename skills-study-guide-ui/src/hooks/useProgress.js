import { useState, useEffect, useCallback } from 'react'

const STORAGE_KEY = 'sg-progress'

function loadProgress() {
  try {
    return JSON.parse(localStorage.getItem(STORAGE_KEY) || '{}')
  } catch {
    return {}
  }
}

function saveProgress(data) {
  localStorage.setItem(STORAGE_KEY, JSON.stringify(data))
}

export function useProgress(slug) {
  const [allProgress, setAllProgress] = useState(loadProgress)

  const guideProgress = allProgress[slug] || { percent: 0, lastVisited: null }

  const updateProgress = useCallback((percent) => {
    setAllProgress(prev => {
      const next = {
        ...prev,
        [slug]: {
          percent: Math.max(prev[slug]?.percent || 0, Math.round(percent)),
          lastVisited: Date.now(),
        }
      }
      saveProgress(next)
      return next
    })
  }, [slug])

  return { guideProgress, allProgress, updateProgress }
}

export function useAllProgress() {
  const [progress] = useState(loadProgress)
  return progress
}

export function getLastReadGuide() {
  const progress = loadProgress()
  let latest = null
  let latestTime = 0

  for (const [slug, data] of Object.entries(progress)) {
    if (data.lastVisited && data.lastVisited > latestTime && data.percent < 100) {
      latest = { slug, ...data }
      latestTime = data.lastVisited
    }
  }

  return latest
}

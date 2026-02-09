export function extractHeadings(content) {
  if (!content) return []

  const lines = content.split('\n')
  const headings = []
  const idCounts = new Map()

  for (const line of lines) {
    const match = line.match(/^(#{2,3})\s+(.+)$/)
    if (match) {
      const text = match[2].replace(/[*`[\]]/g, '').trim()
      let baseId = text.toLowerCase()
        .replace(/\s+/g, '-')
        .replace(/[^\w-]/g, '')

      // Handle ID collisions
      const count = idCounts.get(baseId) || 0
      const id = count > 0 ? `${baseId}-${count}` : baseId
      idCounts.set(baseId, count + 1)

      headings.push({ level: match[1].length, title: text, id })
    }
  }

  return headings
}

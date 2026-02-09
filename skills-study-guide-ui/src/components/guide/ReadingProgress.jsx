export default function ReadingProgress({ progress }) {
  if (progress <= 0) return null

  return (
    <div
      className="reading-progress no-print"
      style={{ width: `${progress}%` }}
    />
  )
}

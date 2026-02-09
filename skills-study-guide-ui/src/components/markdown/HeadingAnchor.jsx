import { Link2 } from 'lucide-react'

export default function HeadingAnchor({ level, children, ...props }) {
  const text = getTextContent(children)
  const id = text.toLowerCase().replace(/\s+/g, '-').replace(/[^\w-]/g, '')

  const Tag = `h${level}`

  // Skip rendering h1 — the GuideHeader already shows the title
  if (level === 1) return null

  return (
    <Tag id={id} {...props}>
      {children}
      <a
        href={`#${id}`}
        className="anchor-link inline-block align-middle"
        aria-label={`Link to ${text}`}
      >
        <Link2 size={14} />
      </a>
    </Tag>
  )
}

function getTextContent(children) {
  if (typeof children === 'string') return children
  if (Array.isArray(children)) return children.map(getTextContent).join('')
  if (children?.props?.children) return getTextContent(children.props.children)
  return ''
}

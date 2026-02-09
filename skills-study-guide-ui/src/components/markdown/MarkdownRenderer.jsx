import { useMemo } from 'react'
import ReactMarkdown from 'react-markdown'
import remarkGfm from 'remark-gfm'
import CodeBlock from './CodeBlock'
import TableWrapper from './TableWrapper'
import HeadingAnchor from './HeadingAnchor'
import CalloutBox from './CalloutBox'

export default function MarkdownRenderer({ content }) {
  const processedContent = useMemo(() => preprocessMarkdown(content), [content])

  return (
    <div className="prose prose-gray dark:prose-invert font-body">
      <ReactMarkdown
        remarkPlugins={[remarkGfm]}
        components={{
          code: CodeBlock,
          table: TableWrapper,
          h1: (props) => <HeadingAnchor level={1} {...props} />,
          h2: (props) => {
            const text = getChildText(props.children)
            // Detect "How Anshul Used It" section headings
            if (text.toLowerCase().includes('how anshul used')) {
              return (
                <>
                  <HeadingAnchor level={2} {...props} />
                  <CalloutBox type="anshul" title="Real-World Application" />
                </>
              )
            }
            return <HeadingAnchor level={2} {...props} />
          },
          h3: (props) => <HeadingAnchor level={3} {...props} />,
          h4: (props) => <HeadingAnchor level={4} {...props} />,
          // Style blockquotes as callouts
          blockquote: ({ children }) => (
            <div className="my-6 rounded-r-lg border-l-4 border-l-[var(--accent)] bg-[var(--accent-subtle)] p-4 not-prose">
              <div className="text-[14px] leading-relaxed text-[var(--text-primary)] [&>p]:m-0">
                {children}
              </div>
            </div>
          ),
        }}
      >
        {processedContent}
      </ReactMarkdown>
    </div>
  )
}

function getChildText(children) {
  if (typeof children === 'string') return children
  if (Array.isArray(children)) return children.map(getChildText).join('')
  if (children?.props?.children) return getChildText(children.props.children)
  return ''
}

// Preprocess: extract code blocks from blockquotes
function preprocessMarkdown(content) {
  if (!content) return ''

  const lines = content.split('\n')
  const result = []
  let inBlockquote = false
  let blockquoteBuffer = []

  for (let i = 0; i < lines.length; i++) {
    const line = lines[i]

    if (line.startsWith('>') && !inBlockquote) {
      inBlockquote = true
      blockquoteBuffer = [line.replace(/^>\s?/, '')]
    } else if (inBlockquote && line.startsWith('>')) {
      blockquoteBuffer.push(line.replace(/^>\s?/, ''))
    } else if (inBlockquote) {
      // End of blockquote — process buffer
      result.push(...processBlockquoteBuffer(blockquoteBuffer))
      result.push(line)
      inBlockquote = false
      blockquoteBuffer = []
    } else {
      result.push(line)
    }
  }

  if (blockquoteBuffer.length > 0) {
    result.push(...processBlockquoteBuffer(blockquoteBuffer))
  }

  return result.join('\n')
}

function processBlockquoteBuffer(buffer) {
  const result = []
  let inCode = false
  let codeBuffer = []
  let language = ''

  for (const line of buffer) {
    const startFence = line.match(/^```(\w*)/)
    const endFence = line.match(/^```$/)

    if (startFence && !inCode) {
      inCode = true
      language = startFence[1] || 'text'
      codeBuffer = []
    } else if (endFence && inCode) {
      result.push('```' + language)
      result.push(...codeBuffer)
      result.push('```')
      result.push('')
      inCode = false
      codeBuffer = []
    } else if (inCode) {
      codeBuffer.push(line)
    } else {
      result.push('> ' + line)
    }
  }

  if (inCode && codeBuffer.length > 0) {
    result.push('```' + language)
    result.push(...codeBuffer)
    result.push('```')
  }

  return result
}

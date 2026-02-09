import { useMemo } from 'react'
import ReactMarkdown from 'react-markdown'
import remarkGfm from 'remark-gfm'
import CodeBlock from '../markdown/CodeBlock'
import TableWrapper from '../markdown/TableWrapper'
import HeadingAnchor from '../markdown/HeadingAnchor'

export default function MarkdownSection({ content, className = '' }) {
  if (!content) return null

  return (
    <div className={`prose prose-gray dark:prose-invert font-body ${className}`}>
      <ReactMarkdown
        remarkPlugins={[remarkGfm]}
        components={{
          code: CodeBlock,
          table: TableWrapper,
          h1: (props) => <HeadingAnchor level={1} {...props} />,
          h2: (props) => <HeadingAnchor level={2} {...props} />,
          h3: (props) => <HeadingAnchor level={3} {...props} />,
          h4: (props) => <HeadingAnchor level={4} {...props} />,
          blockquote: ({ children }) => (
            <div className="my-6 rounded-r-lg border-l-4 border-l-[var(--accent)] bg-[var(--accent-subtle)] p-4 not-prose">
              <div className="text-[14px] leading-relaxed text-[var(--text-primary)] [&>p]:m-0">
                {children}
              </div>
            </div>
          ),
        }}
      >
        {content}
      </ReactMarkdown>
    </div>
  )
}

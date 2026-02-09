import { useEffect, lazy, Suspense } from 'react'
import { useScrollSpy } from '../../hooks/useScrollSpy'
import { useReadingProgress } from '../../hooks/useReadingProgress'
import { useProgress } from '../../hooks/useProgress'
import guideIndex from '../../guide-index.json'
import GuideHeader from '../guide/GuideHeader'
import GuideNavigation from '../guide/GuideNavigation'
import ReadingProgress from '../guide/ReadingProgress'
import TableOfContents from '../guide/TableOfContents'
import MarkdownSection from './MarkdownSection'
import { Briefcase } from 'lucide-react'

// Content sections
import {
  whatIsKafka, coreArchitecture, howMessagesFlow, kafkaInternals,
  partitioningStrategies, deliverySemantics, consumerGroups,
  kafkaConnectStreams, multiRegion, completableFuture, performanceTuning,
  whatIsAvro, avroComparison, avroSchema, schemaEvolution,
  compatibilityRules, avroCodeExample, whatIsSchemaRegistry,
  schemaIdMagicByte, compatibilityModes, producerConsumerIntegration,
  howAnshulUsedIt, quickReference, qaItems,
} from './kafka-content'

// Interactive components — lazy loaded
const ConsumerGroupSimulator = lazy(() => import('./ConsumerGroupSimulator'))
const DeliverySemanticsSimulator = lazy(() => import('./DeliverySemanticsSimulator'))
const MessageFlowAnimation = lazy(() => import('./MessageFlowAnimation'))
const FlashcardDeck = lazy(() => import('./FlashcardDeck'))
const KafkaClusterDiagram = lazy(() => import('./KafkaClusterDiagram'))
const AnnotatedCodeBlock = lazy(() => import('./AnnotatedCodeBlock'))

const SLUG = '01-kafka-avro-schema-registry'

// Combine all markdown for heading extraction
const fullContent = [
  whatIsKafka, coreArchitecture, howMessagesFlow, kafkaInternals,
  partitioningStrategies, deliverySemantics, consumerGroups,
  kafkaConnectStreams, multiRegion, completableFuture, performanceTuning,
  whatIsAvro, avroComparison, avroSchema, schemaEvolution,
  compatibilityRules, avroCodeExample, whatIsSchemaRegistry,
  schemaIdMagicByte, compatibilityModes, producerConsumerIntegration,
  howAnshulUsedIt, quickReference,
].join('\n')

function ComponentLoader({ children }) {
  return (
    <Suspense fallback={
      <div className="my-8 p-8 rounded-xl border border-[var(--border)] bg-[var(--bg-surface-alt)] flex items-center justify-center">
        <span className="font-mono text-xs text-[var(--text-tertiary)] animate-pulse">Loading interactive component...</span>
      </div>
    }>
      {children}
    </Suspense>
  )
}

function SectionDivider({ title }) {
  return (
    <div className="flex items-center gap-3 mt-16 mb-8">
      <div className="h-px flex-1 bg-[var(--border)]" />
      <span className="font-mono text-[10px] font-semibold uppercase tracking-[0.08em] text-[var(--accent)]">
        {title}
      </span>
      <div className="h-px flex-1 bg-[var(--border)]" />
    </div>
  )
}

function Callout({ type = 'anshul', title, children }) {
  const styles = {
    anshul: 'border-l-[var(--accent)] bg-[var(--accent-subtle)]',
    tip: 'border-l-emerald-500 bg-emerald-50 dark:bg-emerald-950/20',
  }
  return (
    <div className={`my-6 rounded-r-lg border-l-4 ${styles[type] || styles.anshul} p-4 not-prose`}>
      <div className="flex items-center gap-2 mb-2">
        <Briefcase size={14} className="text-[var(--accent)]" />
        <span className="font-mono text-[10px] font-semibold uppercase tracking-[0.06em] text-[var(--accent)]">
          {title || 'Real-World Application'}
        </span>
      </div>
      <div className="text-[14px] leading-relaxed">{children}</div>
    </div>
  )
}

export default function KafkaGuidePage() {
  const { activeHeading, headings } = useScrollSpy(fullContent)
  const progress = useReadingProgress()
  const { updateProgress } = useProgress(SLUG)

  const guide = guideIndex.find(g => g.slug === SLUG)
  const guideIdx = guideIndex.findIndex(g => g.slug === SLUG)
  const prevGuide = guideIdx > 0 ? guideIndex[guideIdx - 1] : null
  const nextGuide = guideIdx < guideIndex.length - 1 ? guideIndex[guideIdx + 1] : null

  useEffect(() => {
    window.scrollTo(0, 0)
    document.title = `${guide?.title || 'Kafka'} — Study Guides`
    return () => { document.title = 'Study Guides — Anshul Garg' }
  }, [guide])

  useEffect(() => {
    if (progress > 5) updateProgress(progress)
  }, [progress, updateProgress])

  return (
    <div className="animate-fade-in">
      <ReadingProgress progress={progress} />

      <div className="flex justify-center">
        <article className="flex-1 min-w-0 max-w-[780px] px-5 md:px-8 py-8">
          {guide && <GuideHeader guide={guide} />}

          <div className="guide-content">
            {/* ========== PART 1: APACHE KAFKA ========== */}
            <SectionDivider title="Part 1: Apache Kafka" />

            {/* What Is Kafka — pure markdown */}
            <MarkdownSection content={whatIsKafka} />

            {/* Core Architecture — INTERACTIVE cluster diagram */}
            <MarkdownSection content={coreArchitecture.split('```mermaid')[0]} />
            <ComponentLoader>
              <KafkaClusterDiagram />
            </ComponentLoader>
            <MarkdownSection content={coreArchitecture.split('---')[1] || ''} />

            {/* How Messages Flow — INTERACTIVE step animation */}
            <ComponentLoader>
              <MessageFlowAnimation />
            </ComponentLoader>

            {/* Kafka Internals — markdown (commit log, offsets, replication) */}
            <MarkdownSection content={kafkaInternals} />

            {/* Partitioning Strategies — markdown + annotated code */}
            <MarkdownSection content={partitioningStrategies} />

            {/* Delivery Semantics — INTERACTIVE simulator */}
            <MarkdownSection content={deliverySemantics.split('### At-Most-Once')[0]} />
            <ComponentLoader>
              <DeliverySemanticsSimulator />
            </ComponentLoader>

            {/* Consumer Groups — INTERACTIVE simulator */}
            <MarkdownSection content={consumerGroups.split('```mermaid')[0]} />
            <ComponentLoader>
              <ConsumerGroupSimulator />
            </ComponentLoader>
            <MarkdownSection content={consumerGroups.split('---').pop() || ''} />

            {/* Kafka Connect & Streams — mostly markdown */}
            <MarkdownSection content={kafkaConnectStreams} />

            {/* Multi-Region — markdown + callout */}
            <MarkdownSection content={multiRegion} />

            {/* CompletableFuture — markdown */}
            <MarkdownSection content={completableFuture} />

            {/* Performance Tuning — markdown (TODO: add slider dashboard later) */}
            <MarkdownSection content={performanceTuning} />

            {/* ========== PART 2: APACHE AVRO ========== */}
            <SectionDivider title="Part 2: Apache Avro" />

            <MarkdownSection content={whatIsAvro} />
            <MarkdownSection content={avroComparison} />
            <MarkdownSection content={avroSchema} />
            <MarkdownSection content={schemaEvolution} />
            <MarkdownSection content={compatibilityRules} />
            <MarkdownSection content={avroCodeExample} />

            {/* ========== PART 3: SCHEMA REGISTRY ========== */}
            <SectionDivider title="Part 3: Confluent Schema Registry" />

            <MarkdownSection content={whatIsSchemaRegistry} />
            <MarkdownSection content={schemaIdMagicByte} />
            <MarkdownSection content={compatibilityModes} />
            <MarkdownSection content={producerConsumerIntegration} />

            {/* ========== PART 4: INTERVIEW Q&A ========== */}
            <SectionDivider title="Part 4: Interview Q&A" />

            <ComponentLoader>
              <FlashcardDeck items={qaItems} />
            </ComponentLoader>

            {/* ========== PART 5: HOW ANSHUL USED IT ========== */}
            <SectionDivider title="Part 5: How Anshul Used It at Walmart" />

            <Callout title="Walmart Audit Logging System">
              <p>Anshul designed a <strong>three-tier Kafka-based audit logging system</strong> processing 2M+ events/day across US, Canada, and Mexico.</p>
            </Callout>

            <MarkdownSection content={howAnshulUsedIt} />

            {/* Quick Reference */}
            <MarkdownSection content={quickReference} />
          </div>

          <GuideNavigation prev={prevGuide} next={nextGuide} />
        </article>

        {/* Right TOC */}
        <div className="hidden xl:block w-[200px] shrink-0">
          <TableOfContents headings={headings} activeHeading={activeHeading} />
        </div>
      </div>
    </div>
  )
}

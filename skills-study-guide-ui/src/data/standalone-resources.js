/**
 * Standalone resource files from public/resources/
 * Categorized for sidebar navigation.
 * 70 files total.
 */
// Pinned files — always visible at top, no collapsing
export const pinnedResources = [
  { file: '06-rippling-behavioral-gaps.md', label: 'Rippling Behavioral Gaps' },
  { file: 'BEHAVIORAL-STORIES-STAR.md', label: 'Behavioral Stories (STAR)' },
  { file: 'GCC-BULLET-1-DISTRIBUTED-PIPELINE.md', label: 'GCC Bullet 1: Distributed Pipeline' },
  { file: 'GCC-BULLET-2-CLICKHOUSE-MIGRATION.md', label: 'GCC Bullet 2: ClickHouse Migration' },
  { file: 'GCC-BULLET-3-DATA-PLATFORM.md', label: 'GCC Bullet 3: Data Platform' },
  { file: 'GCC-BULLET-4-DUAL-DATABASE-API.md', label: 'GCC Bullet 4: Dual Database API' },
  { file: 'GCC-BULLET-5-ASSETS-DISCOVERY.md', label: 'GCC Bullet 5: Assets Discovery' },
  { file: 'GCC-BULLET-6-FAKE-FOLLOWER-ML.md', label: 'GCC Bullet 6: Fake Follower ML' },
  { file: 'GCC-INTERVIEW-MASTERCLASS.md', label: 'GCC Interview Masterclass' },
  { file: 'SCENARIO-BASED-QUESTIONS.md', label: 'Scenario Based Questions' },
  { file: 'SCALING-AND-FAILURE-MODES-DEEP-DIVE.md', label: 'Scaling & Failure Modes' },
  { file: 'SYSTEM_INTERCONNECTIVITY.md', label: 'System Interconnectivity' },
  { file: 'TECHNICAL-CONCEPTS-REVIEW.md', label: 'Technical Concepts Review' },
]

export const resourceCategories = [
  {
    name: 'Walmart Interview',
    icon: 'Building2',
    resources: [
      { file: 'WALMART_MASTER_PORTFOLIO.md', label: 'Master Portfolio' },
      { file: 'WALMART_SYSTEM_ARCHITECTURE.md', label: 'System Architecture' },
      { file: 'WALMART_SYSTEM_DESIGN_EXAMPLES.md', label: 'System Design Examples' },
      { file: 'WALMART_PROJECTS_DEEP_DIVE.md', label: 'Projects Deep Dive' },
      { file: 'WALMART_INTERVIEW_ALL_QUESTIONS.md', label: 'All Interview Questions' },
      { file: 'WALMART_GOOGLEYNESS_INDEX.md', label: 'Googleyness Index' },
      { file: 'WALMART_GOOGLEYNESS_QUESTIONS.md', label: 'Googleyness Questions' },
      { file: 'WALMART_HIRING_MANAGER_GUIDE.md', label: 'Hiring Manager Guide' },
      { file: 'WALMART_LEADERSHIP_STORIES.md', label: 'Leadership Stories' },
      { file: 'WALMART_METRICS_CHEATSHEET.md', label: 'Metrics Cheatsheet' },
      { file: 'WALMART_RESUME_TO_CODE_MAPPING.md', label: 'Resume to Code Mapping' },
    ],
  },
  {
    name: 'GCC Interview',
    icon: 'Briefcase',
    resources: [
      { file: 'GCC-INTERVIEW-MASTERCLASS.md', label: 'Interview Masterclass' },
      { file: 'GCC-SYSTEM-ARCHITECTURE.md', label: 'System Architecture' },
      { file: 'GCC-RESUME-BULLETS-DEEP-DIVE.md', label: 'Resume Bullets Deep Dive' },
      { file: 'GCC-PAYU-INTERVIEW-PREP.md', label: 'PayU Interview Prep' },
      { file: 'GCC-BULLET-1-DISTRIBUTED-PIPELINE.md', label: 'Bullet 1: Distributed Pipeline' },
      { file: 'GCC-BULLET-2-CLICKHOUSE-MIGRATION.md', label: 'Bullet 2: ClickHouse Migration' },
      { file: 'GCC-BULLET-3-DATA-PLATFORM.md', label: 'Bullet 3: Data Platform' },
      { file: 'GCC-BULLET-4-DUAL-DATABASE-API.md', label: 'Bullet 4: Dual Database API' },
      { file: 'GCC-BULLET-5-ASSETS-DISCOVERY.md', label: 'Bullet 5: Assets Discovery' },
      { file: 'GCC-BULLET-6-FAKE-FOLLOWER-ML.md', label: 'Bullet 6: Fake Follower ML' },
      { file: 'ANALYSIS_beat.md', label: 'Analysis: Beat' },
      { file: 'ANALYSIS_coffee.md', label: 'Analysis: Coffee' },
      { file: 'ANALYSIS_event_grpc.md', label: 'Analysis: Event gRPC' },
      { file: 'ANALYSIS_fake_follower_analysis.md', label: 'Analysis: Fake Follower' },
      { file: 'ANALYSIS_saas_gateway.md', label: 'Analysis: SaaS Gateway' },
      { file: 'ANALYSIS_stir.md', label: 'Analysis: Stir' },
      { file: 'BEAT_ADVANCED_FEATURES.md', label: 'Beat Advanced Features' },
    ],
  },
  {
    name: 'General Interview',
    icon: 'GraduationCap',
    resources: [
      { file: 'INTERVIEW-MASTER-INDEX.md', label: 'Master Index' },
      { file: 'INTERVIEW-PREP-COMPLETE-GUIDE.md', label: 'Complete Guide' },
      { file: 'INTERVIEW-PREP-PART2-BULLETS-4-12.md', label: 'Prep Part 2: Bullets 4-12' },
      { file: 'INTERVIEW-PREP-PART3-NEW-BULLETS-6-12.md', label: 'Prep Part 3: Bullets 6-12' },
      { file: 'INTERVIEW-MASTERCLASS-KAFKA-AUDIT.md', label: 'Masterclass: Kafka Audit' },
      { file: 'INTERVIEW-PLAYBOOK-HOW-TO-SPEAK.md', label: 'Playbook: How to Speak' },
      { file: 'INTERVIEW-ROUND-GUIDE.md', label: 'Round Guide' },
      { file: 'INTERVIEW_STORY_VERIFICATION.md', label: 'Story Verification' },
      { file: 'BEHAVIORAL-STORIES-STAR.md', label: 'Behavioral Stories (STAR)' },
      { file: 'MOCK-INTERVIEW-QUESTIONS.md', label: 'Mock Interview Questions' },
      { file: 'SCENARIO-BASED-QUESTIONS.md', label: 'Scenario Based Questions' },
      { file: 'QUICK-REFERENCE-CARDS.md', label: 'Quick Reference Cards' },
      { file: 'GAP-ANALYSIS.md', label: 'Gap Analysis' },
      { file: 'MASTER_PORTFOLIO_SUMMARY.md', label: 'Master Portfolio Summary' },
    ],
  },
  {
    name: 'Google Prep',
    icon: 'Search',
    resources: [
      { file: 'GOOGLE_INTERVIEW_MASTER_GUIDE.md', label: 'Master Guide' },
      { file: 'GOOGLE_INTERVIEW_DETAILED.md', label: 'Detailed Guide' },
      { file: 'GOOGLE_INTERVIEW_PREP.md', label: 'Interview Prep' },
      { file: 'GOOGLE_INTERVIEW_SCRIPTS.md', label: 'Interview Scripts' },
      { file: 'GOOGLE_L4_FINAL_PREP.md', label: 'L4 Final Prep' },
      { file: 'GOOGLEYNESS_ALL_QUESTIONS.md', label: 'Googleyness All Questions' },
      { file: 'README_GOOGLE_PREP.md', label: 'Google Prep Readme' },
      { file: 'PREVIOUS_WORK_GOOGLEYNESS_INDEX.md', label: 'Previous Work Googleyness Index' },
      { file: 'PREVIOUS_WORK_PROJECTS_DEEP_DIVE.md', label: 'Previous Work Projects Deep Dive' },
    ],
  },
  {
    name: 'Code Analysis',
    icon: 'FileCode',
    resources: [
      { file: '01-AUDIT-API-LOGS-GCS-SINK-ANALYSIS.md', label: 'Audit API Logs GCS Sink' },
      { file: '02-AUDIT-API-LOGS-SRV-ANALYSIS.md', label: 'Audit API Logs Service' },
      { file: '03-CP-NRTI-APIS-ANALYSIS.md', label: 'CP NRTI APIs' },
      { file: '04-DV-API-COMMON-LIBRARIES-ANALYSIS.md', label: 'DV API Common Libraries' },
      { file: '05-INVENTORY-EVENTS-SRV-ANALYSIS.md', label: 'Inventory Events Service' },
      { file: '06-INVENTORY-STATUS-SRV-ANALYSIS.md', label: 'Inventory Status Service' },
    ],
  },
  {
    name: 'Deep Dives',
    icon: 'BookOpen',
    resources: [
      { file: 'DEEP-DIVE-BULLET-1-2-10-KAFKA-AUDIT.md', label: 'Kafka Audit Deep Dive' },
      { file: 'DEEP-DIVE-BULLET-4-SPRINGBOOT3-JAVA17.md', label: 'Spring Boot 3 & Java 17 Deep Dive' },
      { file: 'SCALING-AND-FAILURE-MODES-DEEP-DIVE.md', label: 'Scaling & Failure Modes' },
      { file: 'TECHNICAL-CONCEPTS-REVIEW.md', label: 'Technical Concepts Review' },
      { file: 'SYSTEM_INTERCONNECTIVITY.md', label: 'System Interconnectivity' },
    ],
  },
  {
    name: 'Resume & PRs',
    icon: 'FileText',
    resources: [
      { file: 'RESUME-ANALYSIS-AND-IMPROVED-VERSION.md', label: 'Resume Analysis & Improved Version' },
      { file: 'RESUME-ANALYSIS-AND-RECOMMENDATIONS.md', label: 'Resume Analysis & Recommendations' },
      { file: 'RESUME_TO_CODE_MAPPING.md', label: 'Resume to Code Mapping' },
      { file: 'PR-PROJECT-MAPPING.md', label: 'PR Project Mapping' },
      { file: 'RELEVANT-PRS-BY-BULLET.md', label: 'Relevant PRs by Bullet' },
      { file: 'ALL-PRS-COMPLETE-LIST.md', label: 'All PRs Complete List' },
      { file: 'README.md', label: 'Resources Readme' },
    ],
  },
  {
    name: 'Rippling Prep',
    icon: 'Zap',
    resources: [
      { file: '06-rippling-behavioral-gaps.md', label: 'Behavioral Gaps' },
    ],
  },
]

/**
 * Flat lookup: filename -> { category, label }
 * Useful for breadcrumbs and page titles.
 */
export const resourceLookup = new Map()
for (const cat of resourceCategories) {
  for (const r of cat.resources) {
    resourceLookup.set(r.file, { category: cat.name, label: r.label, icon: cat.icon })
  }
}

/**
 * Get resource metadata by filename (with or without .md extension).
 */
export function getResource(filename) {
  const key = filename.endsWith('.md') ? filename : `${filename}.md`
  return resourceLookup.get(key) || null
}

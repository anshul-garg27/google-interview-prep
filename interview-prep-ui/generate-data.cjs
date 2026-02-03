// Script to generate data.js from markdown files
const fs = require('fs');
const path = require('path');

const parentDir = path.join(__dirname, '..');

const documents = [
  // Master Guides
  { id: 'GOOGLE_INTERVIEW_MASTER_GUIDE', title: 'Google Interview Master Guide', category: 'master', file: 'GOOGLE_INTERVIEW_MASTER_GUIDE.md', badge: 'Start Here' },
  { id: 'README_GOOGLE_PREP', title: 'README - Complete Prep Guide', category: 'master', file: 'README_GOOGLE_PREP.md' },

  // Google Interview Prep
  { id: 'GOOGLEYNESS_ALL_QUESTIONS', title: '60+ Googleyness Questions', category: 'google-interview', file: 'GOOGLEYNESS_ALL_QUESTIONS.md', badge: 'Must Read' },
  { id: 'GOOGLE_L4_FINAL_PREP', title: 'L4 Final Prep Guide', category: 'google-interview', file: 'GOOGLE_L4_FINAL_PREP.md' },
  { id: 'GOOGLE_INTERVIEW_SCRIPTS', title: 'Interview Scripts', category: 'google-interview', file: 'GOOGLE_INTERVIEW_SCRIPTS.md' },
  { id: 'GOOGLE_INTERVIEW_PREP', title: 'STAR Stories', category: 'google-interview', file: 'GOOGLE_INTERVIEW_PREP.md' },
  { id: 'GOOGLE_INTERVIEW_DETAILED', title: 'Detailed Breakdown', category: 'google-interview', file: 'GOOGLE_INTERVIEW_DETAILED.md' },

  // Walmart Interview Prep - Summary
  { id: 'WALMART_INTERVIEW_ALL_QUESTIONS', title: 'Walmart - 60+ Questions with STAR Answers', category: 'walmart-interview', file: 'WALMART_INTERVIEW_ALL_QUESTIONS.md', badge: 'Must Read' },
  { id: 'WALMART_GOOGLEYNESS_QUESTIONS', title: 'Walmart - Googleyness Questions', category: 'walmart-interview', file: 'WALMART_GOOGLEYNESS_QUESTIONS.md' },
  { id: 'WALMART_HIRING_MANAGER_GUIDE', title: 'Walmart - Hiring Manager Guide', category: 'walmart-interview', file: 'WALMART_HIRING_MANAGER_GUIDE.md' },
  { id: 'WALMART_LEADERSHIP_STORIES', title: 'Walmart - Leadership Stories', category: 'walmart-interview', file: 'WALMART_LEADERSHIP_STORIES.md' },

  // Walmart Interview Prep - Detailed (Bullets 1-12)
  { id: 'INTERVIEW_PREP_COMPLETE_GUIDE', title: 'Walmart - Detailed Prep (Bullets 1-3)', category: 'walmart-detailed', file: 'INTERVIEW-PREP-COMPLETE-GUIDE.md', badge: 'Deep Dive' },
  { id: 'INTERVIEW_PREP_PART2', title: 'Walmart - Detailed Prep (Bullets 4-5)', category: 'walmart-detailed', file: 'INTERVIEW-PREP-PART2-BULLETS-4-12.md' },
  { id: 'INTERVIEW_PREP_PART3', title: 'Walmart - Detailed Prep (Bullets 6-12)', category: 'walmart-detailed', file: 'INTERVIEW-PREP-PART3-NEW-BULLETS-6-12.md' },

  // Walmart Technical
  { id: 'WALMART_RESUME_TO_CODE_MAPPING', title: 'Walmart - Resume ↔ Code Mapping', category: 'walmart-technical', file: 'WALMART_RESUME_TO_CODE_MAPPING.md' },
  { id: 'WALMART_SYSTEM_ARCHITECTURE', title: 'Walmart - System Architecture', category: 'walmart-technical', file: 'WALMART_SYSTEM_ARCHITECTURE.md' },
  { id: 'WALMART_SYSTEM_DESIGN_EXAMPLES', title: 'Walmart - System Design Examples', category: 'walmart-technical', file: 'WALMART_SYSTEM_DESIGN_EXAMPLES.md' },
  { id: 'WALMART_METRICS_CHEATSHEET', title: 'Walmart - Metrics Cheatsheet', category: 'walmart-technical', file: 'WALMART_METRICS_CHEATSHEET.md' },

  // Walmart Portfolio & Analysis
  { id: 'WALMART_MASTER_PORTFOLIO', title: 'Walmart - Master Portfolio', category: 'walmart-portfolio', file: 'WALMART_MASTER_PORTFOLIO.md' },
  { id: 'RESUME_ANALYSIS', title: 'Walmart - Resume Analysis & Recommendations', category: 'walmart-portfolio', file: 'RESUME-ANALYSIS-AND-RECOMMENDATIONS.md' },

  // Walmart Microservices Analysis
  { id: 'ANALYSIS_01_AUDIT_GCS_SINK', title: 'Walmart - Audit API Logs GCS Sink', category: 'walmart-microservices', file: '01-AUDIT-API-LOGS-GCS-SINK-ANALYSIS.md' },
  { id: 'ANALYSIS_02_AUDIT_SRV', title: 'Walmart - Audit API Logs Service', category: 'walmart-microservices', file: '02-AUDIT-API-LOGS-SRV-ANALYSIS.md' },
  { id: 'ANALYSIS_03_CP_NRTI', title: 'Walmart - CP NRTI APIs', category: 'walmart-microservices', file: '03-CP-NRTI-APIS-ANALYSIS.md' },
  { id: 'ANALYSIS_04_DV_COMMON', title: 'Walmart - DV API Common Libraries', category: 'walmart-microservices', file: '04-DV-API-COMMON-LIBRARIES-ANALYSIS.md' },
  { id: 'ANALYSIS_05_INVENTORY_EVENTS', title: 'Walmart - Inventory Events Service', category: 'walmart-microservices', file: '05-INVENTORY-EVENTS-SRV-ANALYSIS.md' },
  { id: 'ANALYSIS_06_INVENTORY_STATUS', title: 'Walmart - Inventory Status Service', category: 'walmart-microservices', file: '06-INVENTORY-STATUS-SRV-ANALYSIS.md' },

  // Google Technical (Previous Work)
  { id: 'RESUME_TO_CODE_MAPPING', title: 'Previous - Resume ↔ Code Mapping', category: 'google-technical', file: 'RESUME_TO_CODE_MAPPING.md' },
  { id: 'SYSTEM_INTERCONNECTIVITY', title: 'Previous - System Architecture', category: 'google-technical', file: 'SYSTEM_INTERCONNECTIVITY.md' },
  { id: 'BEAT_ADVANCED_FEATURES', title: 'Previous - ML/Stats Features', category: 'google-technical', file: 'BEAT_ADVANCED_FEATURES.md' },

  // Google Analysis (Previous Work)
  { id: 'MASTER_PORTFOLIO_SUMMARY', title: 'Previous - Portfolio Summary', category: 'google-analysis', file: 'MASTER_PORTFOLIO_SUMMARY.md' },
  { id: 'ANALYSIS_beat', title: 'Previous - Beat Analysis', category: 'google-analysis', file: 'ANALYSIS_beat.md' },
  { id: 'ANALYSIS_stir', title: 'Previous - Stir Analysis', category: 'google-analysis', file: 'ANALYSIS_stir.md' },
  { id: 'ANALYSIS_event_grpc', title: 'Previous - Event-grpc Analysis', category: 'google-analysis', file: 'ANALYSIS_event_grpc.md' },
  { id: 'ANALYSIS_fake_follower', title: 'Previous - Fake Follower Analysis', category: 'google-analysis', file: 'ANALYSIS_fake_follower_analysis.md' },
  { id: 'ANALYSIS_coffee', title: 'Previous - Coffee Analysis', category: 'google-analysis', file: 'ANALYSIS_coffee.md' },
  { id: 'ANALYSIS_saas_gateway', title: 'Previous - SaaS Gateway Analysis', category: 'google-analysis', file: 'ANALYSIS_saas_gateway.md' },
];

const output = documents.map(doc => {
  const filePath = path.join(parentDir, doc.file);
  let content = '';

  try {
    content = fs.readFileSync(filePath, 'utf8');
    // JSON.stringify handles escaping, no manual escaping needed
  } catch (e) {
    console.error(`Could not read ${doc.file}:`, e.message);
    content = `# ${doc.title}\n\nContent not available.`;
  }

  return {
    id: doc.id,
    title: doc.title,
    category: doc.category,
    badge: doc.badge || null,
    content: content
  };
});

const jsContent = `// Auto-generated from markdown files
export const documents = ${JSON.stringify(output, null, 2)};
`;

fs.writeFileSync(path.join(__dirname, 'src', 'data.js'), jsContent);
console.log('Generated src/data.js successfully!');

// Script to generate data.js from markdown files
const fs = require('fs');
const path = require('path');

const parentDir = path.join(__dirname, '..');

const documents = [
  // Interview Prep
  { id: 'GOOGLEYNESS_ALL_QUESTIONS', title: '60+ Googleyness Questions', category: 'interview', file: 'GOOGLEYNESS_ALL_QUESTIONS.md', badge: 'Must Read' },
  { id: 'GOOGLE_L4_FINAL_PREP', title: 'L4 Final Prep Guide', category: 'interview', file: 'GOOGLE_L4_FINAL_PREP.md' },
  { id: 'GOOGLE_INTERVIEW_SCRIPTS', title: 'Interview Scripts', category: 'interview', file: 'GOOGLE_INTERVIEW_SCRIPTS.md' },
  { id: 'GOOGLE_INTERVIEW_PREP', title: 'STAR Stories', category: 'interview', file: 'GOOGLE_INTERVIEW_PREP.md' },
  { id: 'GOOGLE_INTERVIEW_DETAILED', title: 'Detailed Breakdown', category: 'interview', file: 'GOOGLE_INTERVIEW_DETAILED.md' },

  // Technical
  { id: 'RESUME_TO_CODE_MAPPING', title: 'Resume ↔ Code Mapping', category: 'technical', file: 'RESUME_TO_CODE_MAPPING.md' },
  { id: 'SYSTEM_INTERCONNECTIVITY', title: 'System Architecture', category: 'technical', file: 'SYSTEM_INTERCONNECTIVITY.md' },
  { id: 'BEAT_ADVANCED_FEATURES', title: 'ML/Stats Features', category: 'technical', file: 'BEAT_ADVANCED_FEATURES.md' },

  // Analysis
  { id: 'MASTER_PORTFOLIO_SUMMARY', title: 'Portfolio Summary', category: 'analysis', file: 'MASTER_PORTFOLIO_SUMMARY.md' },
  { id: 'ANALYSIS_beat', title: 'Beat Analysis', category: 'analysis', file: 'ANALYSIS_beat.md' },
  { id: 'ANALYSIS_stir', title: 'Stir Analysis', category: 'analysis', file: 'ANALYSIS_stir.md' },
  { id: 'ANALYSIS_event_grpc', title: 'Event-grpc Analysis', category: 'analysis', file: 'ANALYSIS_event_grpc.md' },
  { id: 'ANALYSIS_fake_follower', title: 'Fake Follower Analysis', category: 'analysis', file: 'ANALYSIS_fake_follower_analysis.md' },
  { id: 'ANALYSIS_coffee', title: 'Coffee Analysis', category: 'analysis', file: 'ANALYSIS_coffee.md' },
  { id: 'ANALYSIS_saas_gateway', title: 'SaaS Gateway Analysis', category: 'analysis', file: 'ANALYSIS_saas_gateway.md' },
];

const output = documents.map(doc => {
  const filePath = path.join(parentDir, doc.file);
  let content = '';

  try {
    content = fs.readFileSync(filePath, 'utf8');
    // Escape backticks and ${} for template literals
    content = content.replace(/\\/g, '\\\\').replace(/`/g, '\\`').replace(/\$\{/g, '\\${');
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

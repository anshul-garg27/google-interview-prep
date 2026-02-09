const fs = require('fs');
const path = require('path');

const guidesDir = path.join(__dirname, '..', '..', 'skills-study-guide');
const publicDir = path.join(__dirname, '..', 'public', 'guides');
const srcDir = path.join(__dirname, '..', 'src');

fs.mkdirSync(publicDir, { recursive: true });

const CATEGORIES = {
  '01': 'Infrastructure', '02': 'Infrastructure', '03': 'System Design',
  '04': 'Infrastructure', '05': 'Infrastructure',
  '06': 'Databases & Storage', '07': 'Languages', '08': 'Data & Observability',
  '09': 'Languages', '10': 'Languages', '11': 'System Design',
  '12': 'Databases & Storage', '13': 'Interview Prep', '14': 'Interview Prep',
  '15': 'Languages',
};

const CATEGORY_ORDER = ['Infrastructure', 'System Design', 'Databases & Storage', 'Languages', 'Data & Observability', 'Interview Prep'];

const TAGS = {
  '01': ['kafka', 'avro', 'schema-registry', 'streaming'],
  '02': ['grpc', 'protobuf', 'http2', 'rpc'],
  '03': ['system-design', 'distributed', 'microservices', 'eda'],
  '04': ['kubernetes', 'docker', 'istio', 'service-mesh'],
  '05': ['aws', 's3', 'lambda', 'sqs', 'kinesis'],
  '06': ['clickhouse', 'redis', 'rabbitmq', 'olap'],
  '07': ['go', 'goroutines', 'channels', 'concurrency'],
  '08': ['prometheus', 'grafana', 'airflow', 'dbt'],
  '09': ['java', 'java17', 'concurrency', 'gc'],
  '10': ['spring-boot', 'spring', 'hibernate', 'jpa'],
  '11': ['api-design', 'openapi', 'kafka-smt', 'disaster-recovery'],
  '12': ['sql', 'postgresql', 'indexing', 'mvcc'],
  '13': ['behavioral', 'star', 'googleyness', 'leadership'],
  '14': ['dsa', 'algorithms', 'leetcode', 'patterns'],
  '15': ['python', 'fastapi', 'asyncio', 'async'],
};

const files = fs.readdirSync(guidesDir)
  .filter(f => /^\d{2}-/.test(f) && f.endsWith('.md'))
  .sort();

const index = files.map(file => {
  const content = fs.readFileSync(path.join(guidesDir, file), 'utf8');
  const lines = content.split('\n');

  const titleLine = lines.find(l => /^# /.test(l));
  const title = titleLine
    ? titleLine.replace(/^# /, '').replace(/ --.*$/, '').replace(/ — .*$/, '').trim()
    : file;

  const contextLine = lines.find(l => /^\*\*Context:\*\*/.test(l));
  const description = contextLine ? contextLine.replace(/^\*\*Context:\*\*\s*/, '') : '';

  const h2s = (content.match(/^## /gm) || []).length;
  const h3s = (content.match(/^### /gm) || []).length;
  const codeBlockMarkers = (content.match(/^```/gm) || []).length;
  const mermaidCount = (content.match(/^```mermaid/gm) || []).length;
  const tableRows = (content.match(/^\|/gm) || []).length;
  const wordCount = content.split(/\s+/).length;

  const qaMatches = content.match(/\*\*Q\d+[.:]/g) || content.match(/^##+ Q\d+/gm) || [];

  const sections = [];
  lines.forEach(line => {
    const match = line.match(/^(#{2,3})\s+(.+)$/);
    if (match) {
      const text = match[2].replace(/[*`[\]]/g, '').trim();
      sections.push({
        level: match[1].length,
        title: text,
        id: text.toLowerCase().replace(/\s+/g, '-').replace(/[^\w-]/g, ''),
      });
    }
  });

  // Copy to public
  fs.copyFileSync(path.join(guidesDir, file), path.join(publicDir, file));

  const num = file.match(/^(\d{2})/)[1];
  const slug = file.replace(/\.md$/, '');

  return {
    num,
    slug,
    file,
    title,
    description,
    category: CATEGORIES[num] || 'General',
    tags: TAGS[num] || [],
    stats: {
      wordCount,
      sections: h2s + h3s,
      codeBlocks: Math.floor(codeBlockMarkers / 2) - mermaidCount,
      mermaidDiagrams: mermaidCount,
      tableRows,
      qaItems: qaMatches.length,
    },
    readingTime: Math.ceil(wordCount / 200),
    sections,
  };
});

// Write the index
fs.writeFileSync(
  path.join(srcDir, 'guide-index.json'),
  JSON.stringify(index, null, 2)
);

// Write category order
fs.writeFileSync(
  path.join(srcDir, 'categories.json'),
  JSON.stringify(CATEGORY_ORDER, null, 2)
);

console.log(`Generated index for ${index.length} guides (${files.length} files copied to public/guides/)`);

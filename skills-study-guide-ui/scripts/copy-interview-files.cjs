const fs = require('fs');
const path = require('path');

const workExDir = path.join(__dirname, '..', '..');
const publicBase = path.join(__dirname, '..', 'public');

// ============ 1. COPY PROJECT FOLDERS ============
const projects = [
  { slug: 'kafka-audit-logging', folder: '01-kafka-audit-logging' },
  { slug: 'spring-boot-3-migration', folder: '02-spring-boot-3-migration' },
  { slug: 'dsd-notification-system', folder: '03-dsd-notification-system' },
  { slug: 'common-library-jar', folder: '04-common-library-jar' },
  { slug: 'openapi-dc-inventory', folder: '05-openapi-dc-inventory' },
  { slug: 'transaction-event-history', folder: '06-transaction-event-history' },
  { slug: 'observability', folder: '07-observability' },
  { slug: 'gcc-beat-scraping', folder: '08-gcc-beat-scraping' },
  { slug: 'gcc-event-grpc', folder: '09-gcc-event-grpc' },
  { slug: 'gcc-stir-data-platform', folder: '10-gcc-stir-data-platform' },
  { slug: 'gcc-coffee-saas-api', folder: '11-gcc-coffee-saas-api' },
  { slug: 'gcc-fake-follower-ml', folder: '12-gcc-fake-follower-ml' },
];

let totalProjectFiles = 0;
for (const project of projects) {
  const srcDir = path.join(workExDir, project.folder);
  const destDir = path.join(publicBase, 'interview', project.slug);
  fs.mkdirSync(destDir, { recursive: true });

  if (!fs.existsSync(srcDir)) {
    console.warn(`  SKIP: ${project.folder} (not found)`);
    continue;
  }

  const files = fs.readdirSync(srcDir).filter(f => f.endsWith('.md'));
  for (const file of files) {
    fs.copyFileSync(path.join(srcDir, file), path.join(destDir, file));
    totalProjectFiles++;
  }
  console.log(`  ${project.slug}: ${files.length} files`);
}
console.log(`Copied ${totalProjectFiles} project files across ${projects.length} projects`);

// ============ 1b. GENERATE PROJECT FILES INDEX ============
const projectFilesIndex = {};
for (const project of projects) {
  const srcDir = path.join(workExDir, project.folder);
  if (!fs.existsSync(srcDir)) continue;
  const files = fs.readdirSync(srcDir)
    .filter(f => f.endsWith('.md') && f !== '00-INTERVIEW-MASTER.md')
    .sort();
  projectFilesIndex[project.slug] = files.map(f => ({
    file: f,
    label: f.replace(/\.md$/, '').replace(/^\d+-/, '').replace(/[-_]/g, ' ').replace(/\b\w/g, c => c.toUpperCase()),
  }));
}
fs.writeFileSync(
  path.join(__dirname, '..', 'src', 'project-files.json'),
  JSON.stringify(projectFilesIndex, null, 2)
);
console.log(`Generated project-files.json for ${Object.keys(projectFilesIndex).length} projects`);

// ============ 2. COPY STANDALONE MD FILES ============
const resourceDir = path.join(publicBase, 'resources');
fs.mkdirSync(resourceDir, { recursive: true });

const standaloneFiles = fs.readdirSync(workExDir).filter(f =>
  f.endsWith('.md') && !f.startsWith('.')
);

let copiedResources = 0;
for (const file of standaloneFiles) {
  const srcPath = path.join(workExDir, file);
  // Only copy files, not directories
  if (fs.statSync(srcPath).isFile()) {
    fs.copyFileSync(srcPath, path.join(resourceDir, file));
    copiedResources++;
  }
}
console.log(`Copied ${copiedResources} standalone resource files`);

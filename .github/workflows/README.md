# GitHub Actions Workflows

## Deploy Workflow

Automatically builds and deploys the Interview Prep UI to GitHub Pages.

### Triggers

- **Push to main**: Automatically builds and deploys
- **Pull Request**: Builds and lints (no deployment)
- **Manual**: Can be triggered manually from Actions tab

### Jobs

1. **Build**: Compiles React app with Vite
   - Installs dependencies
   - Generates data.js from markdown files
   - Builds production bundle
   - Uploads artifacts

2. **Deploy**: Deploys to GitHub Pages (main branch only)
   - Downloads build artifacts
   - Configures GitHub Pages
   - Deploys to Pages environment

3. **Lint**: Runs ESLint for code quality
   - Checks code formatting
   - Ensures best practices

### Setup Required

1. Enable GitHub Pages in repository settings:
   - Go to Settings → Pages
   - Source: GitHub Actions

2. The workflow will automatically deploy on every push to main

### URLs

- **Production**: https://anshul-garg27.github.io/google-interview-prep/
- **Local Dev**: http://localhost:5173/

### Environment Variables

- `NODE_ENV=production`: Used to set base path for GitHub Pages

### Technologies

- Node.js 20
- npm ci (clean install)
- Vite 7
- React 19
- GitHub Pages

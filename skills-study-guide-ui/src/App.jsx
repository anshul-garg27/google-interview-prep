import { Routes, Route, Navigate } from 'react-router-dom'
import Layout from './components/layout/Layout'
import HomePage from './components/home/HomePage'
import GuidePage from './components/guide/GuidePage'
import DemoPage from './components/demo/DemoPage'
import InterviewPrepPage from './components/interview/InterviewPrepPage'
import ResourcePage from './components/resource/ResourcePage'
import ProjectFilePage from './components/interview/ProjectFilePage'
import { ErrorBoundary } from './components/ErrorBoundary'

export default function App() {
  return (
    <Layout>
      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/demo" element={<DemoPage />} />
        <Route path="/interview/:projectSlug/:filename" element={
          <ErrorBoundary>
            <ProjectFilePage />
          </ErrorBoundary>
        } />
        <Route path="/interview/:projectSlug" element={
          <ErrorBoundary>
            <InterviewPrepPage />
          </ErrorBoundary>
        } />
        <Route path="/resource/:filename" element={
          <ErrorBoundary>
            <ResourcePage />
          </ErrorBoundary>
        } />
        <Route path="/guide/:slug" element={
          <ErrorBoundary>
            <GuidePage />
          </ErrorBoundary>
        } />
        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </Layout>
  )
}

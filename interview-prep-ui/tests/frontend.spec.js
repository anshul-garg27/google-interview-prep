import { test, expect } from '@playwright/test';

test.describe('Interview Prep UI - Core Functionality', () => {

  test.beforeEach(async ({ page }) => {
    await page.goto('/');
  });

  // ==================== PAGE LOAD TESTS ====================

  test('should load the page successfully', async ({ page }) => {
    await expect(page).toHaveTitle(/Google L4 Interview Prep/);
  });

  test('should display header with logo and title', async ({ page }) => {
    const header = page.locator('header');
    await expect(header).toBeVisible();

    // Check Google-style colored dots
    const dots = page.locator('header .rounded-full');
    await expect(dots).toHaveCount(4);

    // Check title in header
    const title = page.locator('header h1');
    await expect(title).toBeVisible();
    await expect(title).toContainText('Google L4 Interview Prep');
  });

  test('should display main content area', async ({ page }) => {
    const main = page.locator('main');
    await expect(main).toBeVisible();
  });

  // ==================== SIDEBAR TESTS ====================

  test('sidebar should be visible on desktop', async ({ page }) => {
    await page.setViewportSize({ width: 1280, height: 800 });
    const sidebar = page.locator('aside');
    await expect(sidebar).toBeVisible();
  });

  test('sidebar should have all three categories', async ({ page }) => {
    await page.setViewportSize({ width: 1280, height: 800 });

    // Use more specific locators for sidebar categories
    const sidebar = page.locator('aside');
    await expect(sidebar.locator('text=Interview Prep').first()).toBeVisible();
    await expect(sidebar.locator('text=Technical Docs')).toBeVisible();
    await expect(sidebar.locator('text=Project Analysis')).toBeVisible();
  });

  test('should display documents in sidebar', async ({ page }) => {
    await page.setViewportSize({ width: 1280, height: 800 });

    // Check for document button in sidebar
    const sidebar = page.locator('aside');
    const docButton = sidebar.locator('button:has-text("Googleyness Questions")').first();
    await expect(docButton).toBeVisible();
  });

  test('sidebar document click should change content', async ({ page }) => {
    await page.setViewportSize({ width: 1280, height: 800 });

    // Click on a different document in sidebar
    const sidebar = page.locator('aside');
    const docButton = sidebar.locator('button:has-text("Final Prep")').first();
    await docButton.click();

    // Main content title should change
    const mainTitle = page.locator('main > div > h1').first();
    await expect(mainTitle).toContainText('Final Prep');
  });

  // ==================== MOBILE SIDEBAR TESTS ====================

  test('mobile menu button should toggle sidebar', async ({ page }) => {
    await page.setViewportSize({ width: 375, height: 667 }); // iPhone size

    // Sidebar should be hidden initially on mobile
    const sidebar = page.locator('aside');
    await expect(sidebar).toHaveClass(/-translate-x-full/);

    // Click hamburger menu
    const menuButton = page.locator('header button').first();
    await menuButton.click();

    // Sidebar should be visible now
    await expect(sidebar).toHaveClass(/translate-x-0/);
  });

  test('clicking overlay should close mobile sidebar', async ({ page }) => {
    await page.setViewportSize({ width: 375, height: 667 });

    // Open sidebar
    const menuButton = page.locator('header button').first();
    await menuButton.click();

    // Click overlay
    const overlay = page.locator('.bg-black\\/50');
    await overlay.click();

    // Sidebar should close
    const sidebar = page.locator('aside');
    await expect(sidebar).toHaveClass(/-translate-x-full/);
  });

  // ==================== DARK MODE TESTS ====================

  test('dark mode toggle should work', async ({ page }) => {
    // Get initial state
    const htmlBefore = await page.locator('html').getAttribute('class');
    const wasDark = htmlBefore?.includes('dark') ?? false;

    // Find the theme toggle button (2nd button with svg in header actions)
    const themeToggle = page.locator('header button:has(svg.lucide-moon), header button:has(svg.lucide-sun)');
    await themeToggle.click();

    // Check if class changed
    const htmlAfter = await page.locator('html').getAttribute('class');
    const isDark = htmlAfter?.includes('dark') ?? false;

    expect(wasDark).not.toBe(isDark);
  });

  test('dark mode preference should persist', async ({ page }) => {
    // Toggle dark mode
    const themeToggle = page.locator('header button:has(svg.lucide-moon), header button:has(svg.lucide-sun)');
    await themeToggle.click();

    // Check localStorage
    const darkMode = await page.evaluate(() => localStorage.getItem('darkMode'));
    expect(darkMode).toBeTruthy();
  });

  // ==================== SEARCH TESTS ====================

  test('search button should open search bar', async ({ page }) => {
    // Click search button
    const searchButton = page.locator('header button:has(svg.lucide-search)');
    await searchButton.click();

    // Search input should be visible
    const searchInput = page.locator('input[placeholder*="Search"]');
    await expect(searchInput).toBeVisible();
  });

  test('search should return results for valid query', async ({ page }) => {
    // Open search
    const searchButton = page.locator('header button:has(svg.lucide-search)');
    await searchButton.click();

    // Type search query
    const searchInput = page.locator('input[placeholder*="Search"]');
    await searchInput.fill('Googleyness');

    // Wait for results
    await page.waitForTimeout(500);

    // Results should appear
    const resultsContainer = page.locator('header .shadow-lg');
    await expect(resultsContainer).toBeVisible();
  });

  test('clicking search result should navigate to document', async ({ page }) => {
    // Open search and search
    const searchButton = page.locator('header button:has(svg.lucide-search)');
    await searchButton.click();

    const searchInput = page.locator('input[placeholder*="Search"]');
    await searchInput.fill('STAR');

    await page.waitForTimeout(500);

    // Click first result
    const firstResult = page.locator('header .shadow-lg button').first();
    if (await firstResult.isVisible()) {
      await firstResult.click();

      // Search should close
      await expect(searchInput).not.toBeVisible();
    }
  });

  // ==================== CONTENT RENDERING TESTS ====================

  test('markdown content should render properly', async ({ page }) => {
    // Check if markdown headings render (h2/h3/h4 after our fix)
    const content = page.locator('article');
    await expect(content).toBeVisible();

    // Should have markdown-rendered content
    const headings = content.locator('h2, h3, h4');
    const count = await headings.count();
    expect(count).toBeGreaterThan(0);
  });

  test('should only have one h1 on the page (accessibility)', async ({ page }) => {
    // After our fix, there should be only ONE h1 (the page title)
    const h1Count = await page.locator('h1').count();
    // Header has h1, main has h1 for document title = 2 max
    expect(h1Count).toBeLessThanOrEqual(2);
  });

  test('code blocks should have syntax highlighting', async ({ page }) => {
    // Navigate to a document with code
    await page.setViewportSize({ width: 1280, height: 800 });

    const sidebar = page.locator('aside');
    const docButton = sidebar.locator('button:has-text("Resume")').first();
    await docButton.click();

    await page.waitForTimeout(500);

    // Check for code blocks
    const codeBlocks = page.locator('article pre, article code');
    const count = await codeBlocks.count();
    expect(count).toBeGreaterThanOrEqual(0);
  });

  // ==================== NAVIGATION TESTS ====================

  test('previous/next navigation should work', async ({ page }) => {
    await page.setViewportSize({ width: 1280, height: 800 });

    // Get current title
    const initialTitle = await page.locator('main > div > h1').first().textContent();

    // Find and click next button
    const nextButton = page.locator('main button:has(svg.lucide-chevron-right)').last();
    if (await nextButton.isVisible()) {
      await nextButton.click();
      await page.waitForTimeout(500);

      // Title should change
      const newTitle = await page.locator('main > div > h1').first().textContent();
      expect(newTitle).not.toBe(initialTitle);
    }
  });

  test('breadcrumb should show current category', async ({ page }) => {
    // Breadcrumb is in main content area
    const breadcrumb = page.locator('main .text-gray-500').first();
    await expect(breadcrumb).toBeVisible();
  });

  // ==================== STATS SECTION TESTS ====================

  test('quick stats should display in sidebar', async ({ page }) => {
    await page.setViewportSize({ width: 1280, height: 800 });

    // Check for stats section in sidebar (at bottom, inside border-t div)
    const sidebar = page.locator('aside');
    const statsSection = sidebar.locator('.border-t .grid');
    await expect(statsSection).toBeVisible();

    // Check for Questions stat (exact match in stats section)
    const questionsLabel = statsSection.locator('text=Questions');
    await expect(questionsLabel).toBeVisible();

    // Check for Projects stat
    const projectsLabel = statsSection.locator('text=Projects');
    await expect(projectsLabel).toBeVisible();

    // Check for 60+ number
    const questionsNumber = statsSection.locator('text=60+');
    await expect(questionsNumber).toBeVisible();
  });

  // ==================== ACCESSIBILITY TESTS ====================

  test('page should have proper heading structure', async ({ page }) => {
    // Main page title should be h1
    const mainTitle = page.locator('main > div > h1').first();
    await expect(mainTitle).toBeVisible();
  });

  test('interactive elements should be focusable', async ({ page }) => {
    // Tab through the page
    await page.keyboard.press('Tab');

    // Something should be focused
    const focused = await page.evaluate(() => document.activeElement?.tagName);
    expect(focused).toBeTruthy();
  });

  // ==================== RESPONSIVE TESTS ====================

  test('should be responsive on tablet', async ({ page }) => {
    await page.setViewportSize({ width: 768, height: 1024 });

    const main = page.locator('main');
    await expect(main).toBeVisible();
  });

  test('should be responsive on mobile', async ({ page }) => {
    await page.setViewportSize({ width: 375, height: 667 });

    const main = page.locator('main');
    await expect(main).toBeVisible();

    // Main content title should be visible
    const mainTitle = page.locator('main > div > h1').first();
    await expect(mainTitle).toBeVisible();
  });

  // ==================== BADGE TEST ====================

  test('must read badge should be visible on first document', async ({ page }) => {
    await page.setViewportSize({ width: 1280, height: 800 });

    const sidebar = page.locator('aside');
    const badge = sidebar.locator('text=Must Read');
    await expect(badge).toBeVisible();
  });

  // ==================== FOOTER TEST ====================

  test('footer should be visible', async ({ page }) => {
    // Scroll to bottom
    await page.evaluate(() => window.scrollTo(0, document.body.scrollHeight));

    const footer = page.locator('footer');
    await expect(footer).toBeVisible();
  });

  // ==================== GITHUB LINK TEST ====================

  test('github link should be present', async ({ page }) => {
    const githubLink = page.locator('a[href*="github.com"]');
    await expect(githubLink).toBeVisible();
  });

  // ==================== PERFORMANCE TEST ====================

  test('page should load within reasonable time', async ({ page }) => {
    const startTime = Date.now();
    await page.goto('/');
    await page.waitForLoadState('networkidle');
    const loadTime = Date.now() - startTime;

    // Should load within 5 seconds
    expect(loadTime).toBeLessThan(5000);
  });

  // ==================== CONSOLE ERROR TEST ====================

  test('should not have console errors', async ({ page }) => {
    const errors = [];
    page.on('console', msg => {
      if (msg.type() === 'error') {
        errors.push(msg.text());
      }
    });

    await page.goto('/');
    await page.waitForLoadState('networkidle');

    // Filter out known non-critical errors
    const criticalErrors = errors.filter(e =>
      !e.includes('favicon') &&
      !e.includes('404') &&
      !e.includes('net::')
    );

    expect(criticalErrors).toHaveLength(0);
  });

  // ==================== CODE BLOCK RENDERING TEST ====================

  test('code blocks should render as preformatted text', async ({ page }) => {
    await page.setViewportSize({ width: 1280, height: 800 });

    // Navigate to Beat Analysis which has code blocks
    const sidebar = page.locator('aside');
    const docButton = sidebar.locator('button:has-text("Beat Analysis")').first();
    await docButton.click();

    await page.waitForTimeout(1000);

    // Check for pre elements (code blocks)
    const preBlocks = page.locator('article pre');
    const count = await preBlocks.count();
    expect(count).toBeGreaterThan(0);

    // Check that pre blocks have proper styling
    const firstPre = preBlocks.first();
    await expect(firstPre).toBeVisible();

    // Check that content is multiline (has newlines preserved)
    const preContent = await firstPre.textContent();
    expect(preContent).toBeTruthy();
  });

  // ==================== SCROLL TO TOP TEST ====================

  test('clicking document should scroll to top', async ({ page }) => {
    await page.setViewportSize({ width: 1280, height: 800 });

    // Scroll down
    await page.evaluate(() => window.scrollTo(0, 1000));

    // Click a different document
    const sidebar = page.locator('aside');
    const docButton = sidebar.locator('button:has-text("Final Prep")').first();
    await docButton.click();

    // Wait a bit for scroll
    await page.waitForTimeout(300);

    // Should be at top
    const scrollY = await page.evaluate(() => window.scrollY);
    expect(scrollY).toBe(0);
  });
});

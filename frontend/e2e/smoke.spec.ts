import { test, expect } from '@playwright/test';
import { login, trackPageHealth } from './helpers';

// Phase 1 smoke suite: automates the manual browser-walks we kept doing by hand.
// Runs against the production-served path (Go backend serving dist under /web),
// NOT `ng serve` — so the backend CSP actually applies.

test('login → dashboard (cookie/JWT auth flow)', async ({ page }) => {
  await login(page);
  await expect(page).toHaveURL(/\/web\/dashboard/);
  await expect(page.locator('app-root')).toBeVisible();
});

test('enforcing CSP header present + no CSP violations / uncaught errors across key pages', async ({ page }) => {
  const health = trackPageHealth(page);

  await login(page);

  // The enforcing Content-Security-Policy must be on the served document (#124).
  const resp = await page.goto('dashboard');
  const csp = resp?.headers()['content-security-policy'] ?? '';
  expect(csp, 'enforcing CSP header missing').toContain("default-src 'self'");
  expect(csp).toContain("object-src 'none'");
  expect(csp).toContain("frame-ancestors 'none'");

  // Walk the main authenticated pages; let async errors/violations surface.
  for (const path of ['dashboard', 'sources', 'medical-history', 'labs']) {
    await page.goto(path);
    await expect(page.locator('app-root')).toBeVisible();
    await page.waitForTimeout(1500);
  }

  expect(health.cspViolations, `enforcing CSP violations:\n${health.cspViolations.join('\n')}`).toEqual([]);
  expect(health.pageErrors, `uncaught page errors:\n${health.pageErrors.join('\n')}`).toEqual([]);
});

test('SPA serves under /web and lforms wc-lhc-form web component registers', async ({ page }) => {
  // lforms web-component bundle is served (regression guard for the asset wiring).
  const js = await page.request.get('assets/js/lforms/lhc-forms.js');
  expect(js.status()).toBe(200);

  // The lforms scripts load on every page (index.html); the login page is enough.
  await page.goto('auth/signin');
  await expect(page.locator('app-root')).toBeVisible();

  // lhc-forms.js must register the <wc-lhc-form> custom element the wizards depend on.
  // (This is the lforms-42 smoke that doesn't require the deep modal navigation.)
  await expect
    .poll(() => page.evaluate(() => !!(window as unknown as { customElements?: CustomElementRegistry }).customElements?.get('wc-lhc-form')), {
      timeout: 20_000,
      message: 'wc-lhc-form custom element was never registered by lhc-forms.js',
    })
    .toBe(true);
});

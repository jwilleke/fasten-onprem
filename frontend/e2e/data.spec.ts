import { test, expect } from '@playwright/test';
import { login, trackPageHealth } from './helpers';
import { API_BASE } from './constants';

// Phase 3 (#131): with a synthetic Synthea bundle seeded in global-setup, verify the source is
// present and that the data-bearing pages render cleanly (no CSP violations / uncaught JS) — i.e.
// real FHIR content displays without errors, not just the empty-account state the smoke suite sees.
test('seeded FHIR data: source present + records pages render clean', async ({ page }) => {
  const health = trackPageHealth(page);
  await login(page);

  // The seeded manual source should be listed. page.request shares the browser's auth cookie.
  const resp = await page.request.get(`${API_BASE}/secure/source`);
  expect(resp.ok(), `GET /api/secure/source -> ${resp.status()}`).toBeTruthy();
  const body = await resp.json();
  const sources = Array.isArray(body) ? body : (body?.data ?? []);
  expect(sources.length, 'expected ≥1 seeded source (Synthea bundle)').toBeGreaterThan(0);

  // Data-bearing pages render without CSP violations or uncaught errors.
  for (const path of ['medical-history', 'labs']) {
    await page.goto(path);
    await expect(page.locator('app-root')).toBeVisible();
    await page.waitForTimeout(1500);
  }
  expect(health.cspViolations, `CSP violations:\n${health.cspViolations.join('\n')}`).toEqual([]);
  expect(health.pageErrors, `uncaught page errors:\n${health.pageErrors.join('\n')}`).toEqual([]);
});

import { request, FullConfig } from '@playwright/test';
import { API_BASE, E2E_USER, E2E_PASS } from './constants';

// Runs once after the webServer (Go backend) is up, before any tests.
// Seeds the test account via the public signup API. The first user becomes admin.
// Idempotent-ish: if the account already exists (reused dev server / non-fresh DB),
// signup returns non-2xx and we ignore it — the login specs will still authenticate.
export default async function globalSetup(_config: FullConfig) {
  const ctx = await request.newContext();
  try {
    const res = await ctx.post(`${API_BASE}/auth/signup`, {
      data: { username: E2E_USER, password: E2E_PASS },
    });
    if (res.ok()) {
      console.log(`[e2e] seeded account "${E2E_USER}"`);
    } else {
      console.log(`[e2e] signup returned ${res.status()} (account likely already exists) — continuing`);
    }
  } catch (e) {
    console.log(`[e2e] signup request failed (${e}) — continuing; login may still work`);
  } finally {
    await ctx.dispose();
  }
}

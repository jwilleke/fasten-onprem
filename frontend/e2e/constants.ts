// Shared constants for the E2E suite (#114-era browser-interaction automation).
// The backend runs on :9191 serving the SPA under /web (see config.e2e.yaml).
export const BASE_URL = 'http://localhost:9191/web/';
export const API_BASE = 'http://localhost:9191/api';

// Seeded by global-setup via POST /api/auth/signup (first user => admin).
export const E2E_USER = 'e2e';
export const E2E_PASS = 'e2e-test-pass-1234'; // not a secret; local throwaway DB only

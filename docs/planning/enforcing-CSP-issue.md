# Enforcing Content-Security-Policy — issue & working plan

> **Status:** in progress · CSP is currently **report-only** (safe) on `main`.
> **Tracking issue:** [#124](https://github.com/jwilleke/yourphr/issues/124) · Related: #105 (H4), #113, #120, #103 (H2), #114.
> **Owner:** living doc — updated with findings as we work through it.
> Companion to [architecture-security-review.md](./architecture-security-review.md) and
> [Standards-Conformance.md](./Standards-Conformance.md).

## Goal

Ship an **enforcing** `Content-Security-Policy` for the YourPHR SPA (served by the Go
backend under `/web/`) that blocks injected/XSS scripts and other classic web attacks,
**without breaking the app**. The security headers themselves (H4, #105) already ship;
this is about flipping the CSP from report-only to enforcing.

## The core problem (clearly stated)

Two things make a *fully strict* `script-src` hard for this app:

1. **The base-href bootstrap must stay inline.** `index.html` has two inline `<script>`
   blocks that run at parse time — the one that `document.write`s `<base href="/web/">`,
   and the lforms web-components guard. They **cannot be externalized**: a relative
   `<script src>` resolves against the page URL *before* `<base href>` exists, 404s, and
   is served as `text/html` → "Refused to execute script". That was the #113 outage.
   - ✅ **Solved** by allowlisting the two inline scripts in `script-src` by their
     **sha256 hash** (deterministic — the browser computes the same hash and runs them).

2. **Inline event handlers are incompatible with hash/nonce `script-src`.** CSP hashes
   and nonces **only cover `<script>` blocks** — they do **not** cover inline event-handler
   attributes (`onclick=`, `onload=`, …). Allowing those requires `'unsafe-hashes'` **plus
   a hash of every handler**. This app's runtime DOM contains inline handlers that are
   **not in our templates** (`grep src/app/**.html` finds none) → they're injected by a
   third-party widget (lforms / dwv) or via `innerHTML`. So enforcing a strict `script-src`
   blocks them and breaks the page.
   - ❌ **Unsolved** — this is the real blocker, and it may involve code we don't own.

**Key reframe:** the #1 XSS payoff — **session-token theft — is already closed by #103**
(HttpOnly cookie, no token in JS). So a fully-strict `script-src` is *incremental*
hardening; it must not block the safe, high-value parts of the CSP.

## What we've tried

### Attempt 1 — externalize the inline bootstrap scripts (#113)

Moved the base-href + lforms-guard scripts to `assets/js/*.js` to satisfy `script-src 'self'`.
**Broke prod** (the relative-load bootstrap paradox above). Reverted in **#120**.

### Attempt 2 — sha256-hash the inline scripts, flip to enforcing (#124, commit `936a8dfd`)

Kept the scripts inline, allowlisted by hash. **The hashes worked** (app loaded, no MIME
breakage — #113 fully fixed). But enforcing `script-src` then **blocked inline event
handlers**, breaking interactivity/render. Reverted (`4f3d5440`).

## Findings (live console, Attempt 2)

| Symptom | Cause | Action |
|---|---|---|
| App loaded, assets OK | inline-script **hashes accepted** ✅ | keep the hash approach |
| `Executing inline event handler violates script-src` (dashboard:18/45) | inline `onX=` handlers (3rd-party/runtime, **not our templates**) | the real blocker → staged approach |
| `oauth4webapi … Unexpected token 'export'` | ESM loaded as classic script — **bundling bug, not CSP** | investigate separately (IdpConnect path) |
| `site.webmanifest` blocked (default-src) | manifest 302'd to Authentik; no `manifest-src` | add `manifest-src 'self'` (non-breaking) |
| Perplexity-CDN font blocked | browser **extension**, not the app | ignore |

## Decision: staged enforcement (two policies at once)

CSP supports a `Content-Security-Policy` (enforce) **and** a
`Content-Security-Policy-Report-Only` (observe) header simultaneously.

**Enforce now** — safe, high-value, unaffected by inline handlers:
`default-src 'self'`, `frame-ancestors 'none'` (clickjacking), `base-uri 'self'`
(base-tag injection), `form-action 'self'` (form hijack), `object-src 'none'` (plugins),
`connect-src 'self' https://wallet.hello.coop https://issuer.hello.coop`,
`img-src 'self' data: https:`, `font-src 'self' data:`, `style-src 'self' 'unsafe-inline'`,
`manifest-src 'self'`, and **`script-src 'self' 'unsafe-inline'`** (still blocks
*cross-origin* script injection; permissive on inline so the app + handlers work).

**Observe** (report-only) the strict target:
`script-src 'self' 'sha256-66XQ…' 'sha256-EnWZB…'` — keeps visibility on what fully-strict
would block; tighten later.

**Rationale:** delivers ~80% of CSP's value **enforced and safe today**, vs. zero
(report-only-everything) or an outage (strict-everything). Token-theft XSS is already
mitigated by #103.

## Validation discipline (the lesson)

We caused two outages by **deploying to validate**. Going forward:
**validate the production-served path *locally* first** — `make build-frontend`, point the
Go backend's `web.src.frontend.path` at the built `dist/`, open `http://localhost:9090/web`
with the console open, walk every major page (dashboard, sources, records, **lforms
questionnaire**, reports, settings). Deploy only when the local console is clean.
`ng serve` does NOT apply the backend CSP and serves at `/`, so it cannot validate this.

## Path to fully-strict `script-src` (separate, lower priority)

1. Enumerate every inline event handler (local browser walk) and its source.
2. Eliminate ours (refactor `onX=` → Angular `(event)` bindings); for third-party widgets,
   replace/fork or accept `'unsafe-hashes'` + enumerated handler hashes.
3. Move script loading to **nonce + `'strict-dynamic'`** (needs per-request `index.html`
   templating in the backend; currently static `StaticFS` / `c.File`).
4. Likely paired with the **Angular upgrade (#114)**, which may change the widget stack.

## Worklog

- **2026-06-07** — Attempt 2 (hashes) reverted (`4f3d5440`); root cause = inline event
  handlers vs strict `script-src`. Decided on **staged enforcement**. This doc created.
  *Next: implement the staged CSP + validate the production-served path locally before deploy.*

## Related issues

- **#124** — this work (re-enable enforcing CSP). · **#105** — H4 security headers (shipped).
- **#113 / #120** — the externalize-scripts outage + revert. · **#103** — HttpOnly cookie
  (closes the primary XSS vector). · **#114** — Angular upgrade (enables the strict path).

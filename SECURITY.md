# Security Policy

YourPHR is a self-hosted **Personal Health Record** viewer — a community continuation of Fasten OnPrem (GPL v3). It stores and displays **Protected Health Information (PHI)** and personally identifying information (PII), so security and privacy are first-order concerns.

## ⚠️ Never include real PHI/PII in a report

When reporting a bug or a vulnerability — in a GitHub issue, a security advisory, logs, screenshots, or attached files — **never include real patient data, personal identifiers, access tokens, secrets, or database files.** Reproduce with **synthetic data only** (e.g. Synthea-generated FHIR bundles). A leak of real PHI is irreversible. If you believe real PHI has actually been exposed somewhere, say so privately (see below) — but do not attach the data itself.

## Reporting a vulnerability

Please report security vulnerabilities **privately**, not in a public issue:

- **Preferred:** GitHub **private vulnerability reporting** — go to the repository's **Security** tab → **Report a vulnerability** (GitHub Security Advisories). This keeps the report private until a fix is ready.
- Include: the affected version/commit, a **synthetic-data** reproduction, the impact, and any suggested fix.

We aim to acknowledge reports within a few days and to coordinate a fix and a disclosure timeline with you. There is no paid bug-bounty program.

## Supported versions

YourPHR ships from `main` (rolling release). Security fixes land on `main` and the published `ghcr.io/jwilleke/yourphr:main-<N>` images — run the latest. Older image tags are not separately patched.

## For operators — handling secrets & data

- **Never commit** secrets, keys, `.env`, real FHIR bundles, or the SQLite DB. See the "NEVER commit personal health data or unencrypted secrets" section of `CLAUDE.md`. `*.db`, `/db/`, `config.dev.yaml`, `certs/`, and key files are gitignored — keep them that way.
- DB encryption is **on by default**; **override** the placeholder `jwt.issuer.key` in any real deployment (via `config.dev.yaml` or an environment variable).
- The app is meant to run behind your own authentication and network controls (e.g. a reverse proxy / forward-auth) on a trusted network — it is not hardened for direct exposure to the public internet.

## Attribution

Built on Fasten OnPrem by Jason Kulatunga ([@AnalogJ](https://github.com/AnalogJ)) and contributors (GPL v3); attribution retained.

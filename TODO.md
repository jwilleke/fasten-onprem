# TODO

<!-- RESUME:START -->
## ▶ Resume here — 2026-06-12

- Last worked on: **FMH EHI classification & patient-legible display — Phase 1 SHIPPED to main** ([#266](https://github.com/jwilleke/yourphr/issues/266) epic, [#267](https://github.com/jwilleke/yourphr/issues/267) Phase 1; commits pushed `26ad8524..67581513`). Backend `pkg/condition.Classify` synthesizes `Condition.category`, derives display state (Active/Remission/Resolved/Unknown/RuledOut), and splits real health problems from FMH "Personal Health Conditions"; new `GET /api/secure/conditions/classified`. Dashboard: Current Medical Concerns = active problem-list-items (badges), collapsed Past Health Problems (date ranges), new **Patient Profile** section. No dedup. Verified: backend build + vet + condition tests; full frontend suite 460 pass / 0 fail.
- Branch / state: main, pushed. CI/Docker building → Flux will deploy a new `main-<N>`.
- **Next step:** once deployed, **verify on the live instance** (real FMH data) that Current Medical Concerns no longer shows social/admin items and that Patient Profile / Past Health Problems populate correctly — then **close [#267](https://github.com/jwilleke/yourphr/issues/267)** (left open pending this live check).
- Then: **Phase 2** of [#266](https://github.com/jwilleke/yourphr/issues/266) — generic provenance + reference resolvers ("who said this"; the `Encounter/<patient>_<id>` quirk), which also unblocks [#264](https://github.com/jwilleke/yourphr/issues/264). Phase 3 = detail-card legibility; Phase 4 = generalize + conformance remodeling (US Core smoking status).
- Reference: real FMH export + Veradigm EHI guide are in `private/phi/` (gitignored). Design: `docs/your-phr-dashboard/classification-and-display-architecture.md`; vendor mapping: `docs/vendors/followmyhealth-ehi-mapping.md`.
<!-- RESUME:END -->

> Generated from GitHub issue labels — do not hand-edit. The `▶ Resume here` pointer is written at session end.

## 🔴 P0 — Security & Critical

_none — security board clean (0 open Dependabot / code-scanning alerts)._

## 🟠 P1

- [#266](https://github.com/jwilleke/yourphr/issues/266) [EPIC] FollowMyHealth/Veradigm EHI classification & patient-legible display (two-layer)
- [#267](https://github.com/jwilleke/yourphr/issues/267) [FEATURE] Phase 1 — Condition classifier & Patient Profile (shipped; pending live verification)
- [#262](https://github.com/jwilleke/yourphr/issues/262) [EPIC] Patient-legible display — health info a normal person can actually use
- [#264](https://github.com/jwilleke/yourphr/issues/264) [FEATURE] Medication card display gaps + reference-resolution blocker
- [#254](https://github.com/jwilleke/yourphr/issues/254) [FEATURE] Support C-CDA / CCD document import and parsing
- [#255](https://github.com/jwilleke/yourphr/issues/255) [FEATURE] Support for PDFs
- [#250](https://github.com/jwilleke/yourphr/issues/250) [FEATURE] Add CMS Blue Button 2.0 as a SMART-on-FHIR sync source
- [#249](https://github.com/jwilleke/yourphr/issues/249) [FEATURE] Surface the 6 remaining US Core 9.0.0 Must-Support display gaps
- [#136](https://github.com/jwilleke/yourphr/issues/136) [EPIC] Support US Core 9.0.0 (profiles + Must-Support)

## 🟡 P2

- [#244](https://github.com/jwilleke/yourphr/issues/244) [EPIC] Per-profile dashboard widgets (US Core display end-state)
- [#256](https://github.com/jwilleke/yourphr/issues/256) [FEATURE] Sharing PHR data
- [#253](https://github.com/jwilleke/yourphr/issues/253) [EPIC] Support manual data entry and user-created records
- [#252](https://github.com/jwilleke/yourphr/issues/252) [FEATURE] Harden re-import dedup against stale overwrites
- [#251](https://github.com/jwilleke/yourphr/issues/251) [FEATURE] Explore Apple Health's supported-institution list as a provider-catalog source
- [#241](https://github.com/jwilleke/yourphr/issues/241) [chore] release-please: authenticate with a PAT / GitHub App token
- [#209](https://github.com/jwilleke/yourphr/issues/209) [EPIC] Migrate to Bootstrap 5
- [#261](https://github.com/jwilleke/yourphr/issues/261) [FEATURE] Settings under User Profile
- [#20](https://github.com/jwilleke/yourphr/issues/20) [EPIC] SMART on FHIR — live provider sync
- [#14](https://github.com/jwilleke/yourphr/issues/14) [FEATURE] User Profile Update (PII)
- [#53](https://github.com/jwilleke/yourphr/issues/53) [SMART] Veradigm/FollowMyHealth integration — blocked on vendor app approval

## ⏸ Deferred

- [#239](https://github.com/jwilleke/yourphr/issues/239) [chore] Revisit gofhir-models 0.1.x once encoding/json/v2 is default in Go
- [#131](https://github.com/jwilleke/yourphr/issues/131) [FEATURE] E2E testing — lforms questionnaire render + interact
- [#263](https://github.com/jwilleke/yourphr/issues/263) [FEATURE] Message Provider

## ❓ Needs triage

_none — all open issues carry a priority band._

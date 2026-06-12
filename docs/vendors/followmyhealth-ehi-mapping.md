# FollowMyHealth / Veradigm EHI export — mapping notes

Practical notes for handling FollowMyHealth (Veradigm, formerly Allscripts) **EHI export** bundles in YourPHR. These are the source-specific facts the [classification & display architecture](../your-phr-dashboard/classification-and-display-architecture.md) Layer 1 adapter relies on.

> **Attribution.** Mappings below are derived from Veradigm's *FollowMyHealth EHI Export Data Guide* (v2, Feb 2024) plus analysis of a real export. The guide is published at `veradigm.com` under `/img/legal/onc/VeradigmFMH_EHI_Export_Data_Guide_v2.pdf` (an ONC-mandated EHI disclosure) and is also linked from each export bundle's `link` with relation `service-doc`. The PDF is **not committed** (its footer marks it Veradigm-proprietary); it is kept under `private/phi/` for reference. This file paraphrases only the interoperability facts we implement.

## What the export is

- A **FHIR R4 `Bundle`** (`type: collection`) in JSON, produced by FollowMyHealth's **EHI export** tool (patient-initiated "export all my data"). Accompanied by document/image files in the same `.zip`.
- **Not** a SMART/Bulk Data API pull. YourPHR ingests it via the manual-upload path.
- Veradigm states it *"uses extensions where internal structures don't map perfectly to standard FHIR resources or value sets"* — and does so pervasively. FollowMyHealth is a **PHR, not an EHR**; the export is a best-effort capture of data from upstream EHR/PM systems.

## Detection signals

Treat a bundle/resource as a FollowMyHealth EHI export when any of:

- the bundle `link[]` has relation `service-doc` with a `fhir.followmyhealth.com` URL;
- resource `identifier[].system` / `code.coding[].system` URLs under `fhir.followmyhealth.com` (e.g. `…/id/global`, `…/id/interface`, `…/id/translation`);
- `identifier[].assigner.display == "FollowMyHealth"`.

## Condition — the big one

FollowMyHealth has **two distinct PHR sections** that it both exports as FHIR `Condition`:

- **My Health > Conditions > Health Conditions** → real medical problems.
- **My Health > Conditions > Personal Health Conditions** → social / lifestyle / administrative items (employment, education, marital status, tobacco/alcohol/substance use, household, etc.).

The guide provides **no** `Condition.category` and no standard mapping for the distinction — the only discriminator is a proprietary identifier. So Layer 1 **synthesizes** the standard `Condition.category` that conformant EHRs already populate.

### Discriminating signals (observed)

| Signal | "Health Condition" (real problem) | "Personal Health Condition" (profile) |
|---|---|---|
| `identifier[].value` interface tell | `…:HealthCondition,<n>` | `…:PersonalHealthConsideration,<n>` |
| `code.coding[]` | standard terminology (ICD‑9/ICD‑10/SNOMED) present | **none** — `code.text` only |
| `recorder` / `asserter` | often a `Practitioner` (sometimes `Patient` = self-reported) | absent |

Presence of a **coded diagnosis** is the most reliable separator (it also catches patient self-reported items the interface tell misses); the interface tell is a strong secondary confirmation. Synthesis rule (see the [Phase 1 spec](../your-phr-dashboard/phase-1-condition-classifier-spec.md) for the exact decision table):

- standard code or `HealthCondition` tell → `category = problem-list-item` (tier clinician)
- coded but patient-asserted, no standard code → `problem-list-item` (tier self-reported)
- text-only + `PersonalHealthConsideration` + no clinician recorder → `category = sdoh` (Patient Profile)
- ambiguous → default to a health problem (never bury a possible diagnosis)

### clinicalStatus / abatement / verificationStatus

- All exported conditions tend to be `clinicalStatus: active` regardless of section — so status alone cannot separate real problems from profile items.
- `abatement[x]` **is** populated for ended conditions (e.g. `abatementDateTime`) and is the native "until"/end-date signal (FHIR con-4: an abated condition's status should be inactive/resolved/remission).
- Honor `verificationStatus`: `entered-in-error` → omit; `refuted` → ruled out.

### Codes

- Real problems carry **ICD‑9 (`…/icd9cm`)** + **ICD‑10 (`hl7.org/fhir/sid/icd-10-cm`)** codings, plus a bare `display`-only coding.
- FollowMyHealth also emits a **proprietary code system** `https://fhir.followmyhealth.com/id/translation` whose `code` is an internal **UUID** (not a real terminology code). Show its `display`; never present the UUID as a "code". Filter `standardCodings` to recognized systems only.

## Reference formats (resolution quirk)

| Reference type | Format | Resolution |
|---|---|---|
| Patient / Practitioner / Organization | `Type/<id>` | direct id match |
| **Encounter** | `Encounter/<patientId>_<encounterId>` (underscore-joined) | **strip the `<patientId>_` prefix, then match `Encounter.id`** |

Naive resolution searches for the whole `patientId_encounterId` blob as an id and silently finds nothing. The `fullUrl` of every resource is also patient-scoped: `…/api/<ResourceType>/<patientId>/<resourceId>` — the middle segment is always the patient, **not** a provider.

## Provenance ("who said this")

- FollowMyHealth EHI exports contain **no `Provenance` resources**, and `Condition`s do **not** link to an `Encounter`.
- But many real conditions carry `recorder`/`asserter` → a `Practitioner` reference that **resolves** to a named `Practitioner` in the bundle, or → `Patient` (self-reported).
- So the provenance chain for FollowMyHealth is: `asserter`/`recorder` → named Practitioner or "Self-reported"; floor = "Source: FollowMyHealth" (the aggregator — never invent an originating clinic).

## Other documented quirks (long tail — handle on demand)

- **Encounter:** not visible in the PHR; `status` always `Unknown` (FollowMyHealth has no encounter-status concept).
- **DiagnosticReport / Observation:** custom extensions `…/ObservationExtension/CollectedOnDate` and `…/OrderedOnDate` carry dates FHIR R4 lacks fields for.
- **FamilyMemberHistory:** relationship "all family" has no FHIR code-system equivalent → emitted as `Display` only; custom `FamilyMemberHistoryExtension`.
- **Immunization:** statuses without a clean FHIR mapping are coerced to the nearest value with the original in `StatusReason`.
- **Document Reference:** `Note`-type docs are `text/plain`; `TransitionOfCare` docs are `application/xml` (C-CDA) — candidates for the C-CDA converter (#254), though re-importing them risks double-reporting (no-dedup decision).

## See also

- Architecture: [`docs/your-phr-dashboard/classification-and-display-architecture.md`](../your-phr-dashboard/classification-and-display-architecture.md)
- Phase 1 spec: [`docs/your-phr-dashboard/phase-1-condition-classifier-spec.md`](../your-phr-dashboard/phase-1-condition-classifier-spec.md)
- Related issues: #262 (legible display), #264 (reference resolution), #53 (Veradigm/FollowMyHealth), #254 (C-CDA import)

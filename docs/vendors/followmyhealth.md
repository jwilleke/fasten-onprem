# FollowMyHealth

Vendor reference for **FollowMyHealth**, the patient portal whose FHIR R4 exports are the primary real-world data YourPHR is being hardened against.

## Overview

**FollowMyHealth** (FMH) is a patient-engagement platform and personal health record (PHR) portal. It gives patients a single account to view records pulled from their providers' EHRs, message care teams, request prescription refills, schedule appointments, and download/export their data. It is **not** an EHR itself — it is a patient-facing aggregation layer that connects to provider systems and presents the patient's longitudinal record.

It is **proprietary, closed-source, vendor-hosted software.** There is no open-source implementation, and patients/developers cannot self-host it. YourPHR's relationship to FollowMyHealth is one-directional: we consume the **patient-initiated data export** (a FHIR R4 bundle) that FMH lets patients download.

## Ownership & History

- **Jardogs LLC** (Bloomington, Illinois; founded ~2009) built FollowMyHealth.
- **Allscripts acquired Jardogs in 2013** and operated FollowMyHealth as its patient-engagement product.
- Allscripts **rebranded to Veradigm in 2023**, so FollowMyHealth is today a **Veradigm** product. See [`veradigm-allscripts.md`](./veradigm-allscripts.md) for the corporate lineage and developer/partner access details.

## Products

FollowMyHealth is itself a suite (patient web portal, mobile apps, and the provider-/payer-facing engagement tooling). For YourPHR only two things matter:

- The **patient data export** — a downloadable archive containing a FHIR R4 JSON bundle plus the document files it references.
- The **SMART on FHIR API** — for live sync (see Known API Issues; not currently usable in this fork).

## Contact

| Purpose | Where |
|---|---|
| Patient support | Through the FollowMyHealth portal / the patient's provider |
| Developer / API support | `VeradigmConnect@veradigm.com` (Veradigm Connect program) |
| Developer portal | `https://developer.veradigm.com` |
| Partner / production-access request | `https://developer.veradigm.com/Account/PartnerRequest` |
| Open support ticket (this project) | **#17849** — Veradigm `unauthorized_client` on connect; steps to reproduce are in [`../FHIR/fhir-testing.md`](../FHIR/fhir-testing.md) and mirrored on [#53](https://github.com/jwilleke/yourphr/issues/53) |

## API & Integration

- **FHIR base / discovery:** FollowMyHealth advertises SMART endpoints under `fhir.followmyhealth.com` with authorization at `muauthentication.followmyhealth.com` (`/api/access` token endpoint, `/api/jwks`). A captured discovery document is in [`../FHIR/fhir-test-discovery-example.md`](../FHIR/fhir-test-discovery-example.md).
- **Access model:** a registered app is **"Test Only"** by default and reaches only Veradigm **test** organizations (synthetic patients, no PHI). Reaching real patient data requires Veradigm to grant **production** access — a partner request with a multi-day review. See the test-vs-real matrix in [`../FHIR/fhir-testing.md`](../FHIR/fhir-testing.md).
- **Auth:** authorization-code + PKCE. The discovery is internally **contradictory** — it advertises a `client-public` capability _and_ confidential-only token-endpoint auth methods (`client-confidential-symmetric`/`-asymmetric`), which is part of why the connect flow stalls (see below).

## Known API Issues

These are the concrete problems YourPHR has hit with real FollowMyHealth data and APIs. Most have been mitigated on the **import/display** side; live sync remains externally blocked.

1. **Live sync is blocked at authorization.** Connecting a FollowMyHealth source fails with `unauthorized_client`; the flow never reaches token exchange. This needs Veradigm to authorize the app (ticket **#17849**). Tracked in [#53](https://github.com/jwilleke/yourphr/issues/53). Because `fasten-sources` is stubbed in this fork, **manual bundle upload is the supported import path.**
2. **Compound reference ids break resource links.** FMH references resources as `Type/{patientId}_{resourceId}` but stores each resource under the bare `{resourceId}`, so links (e.g. Procedure → Encounter) dangle. Fixed in [#196](https://github.com/jwilleke/yourphr/issues/196) by associating the stripped suffix as well. **Existing imports must be re-imported** for the corrected links.
3. **Non-US-Core Encounters.** Often no `type[]`, a `class` with a `system` but no `code`/`display`, and only a `location`. The card now falls back through `type` → `serviceType` → `class` → `location` rather than rendering blank ([#54](https://github.com/jwilleke/yourphr/issues/54) follow-up).
4. **DocumentReference titles.** Documents carry a generic `type.text` (e.g. `"HIPAA"`) and put the real name in `description` / `type.coding[0].display` / `content[0].attachment.title`. Titling off `type.text` labelled thousands of distinct documents identically — fixed in the card ([#198](https://github.com/jwilleke/yourphr/issues/198)) and the backend sort_title ([#201](https://github.com/jwilleke/yourphr/issues/201)).
5. **Document bodies are not in the bundle.** Each DocumentReference's `content[0].attachment.url` is a **relative path** (`./…`) to a sibling file (`.txt`/`.jpg`/`.xml`/`.html`) in the export archive — _not_ inline data. A single JSON-bundle upload therefore carries **metadata only**; rendering the bodies would require a separate file-upload import path. (In one real export, 5,562 such files backed 5,564 DocumentReferences.)
6. **Text-only CodeableConcepts and local code systems.** Many coded fields have `text` but no `coding[]`, or use local systems (e.g. `https://fhir.followmyhealth.com/id/translation`). Search-parameter extraction now indexes text-only concepts ([#171](https://github.com/jwilleke/yourphr/issues/171)).
7. **Encounter dates are usually — but not always — real.** `period.start` is the genuine visit date for ~95% of encounters; a small number of "catch-all" encounters carry the record-creation date instead. YourPHR displays the record-stated `period`, and deliberately does **not** infer a date from linked resources (that would fabricate data).

## Relevance to YourPHR

FollowMyHealth is the **reference non-US-Core dataset** for this fork. The project's near-term mission — immediate, complete patient access to records — depends on rendering FMH exports faithfully even though they diverge from strict US Core. When fixing display issues, prefer graceful fallbacks for missing US-Core fields over assuming conformance. The fixes above (and the broader effort under EPIC [#2](https://github.com/jwilleke/yourphr/issues/2)) come directly from analyzing real FMH exports.

## References / Related issues

- Internal: [`../FHIR/fhir-testing.md`](../FHIR/fhir-testing.md), [`../FHIR/fhir-test-discovery-example.md`](../FHIR/fhir-test-discovery-example.md), [`veradigm-allscripts.md`](./veradigm-allscripts.md), [`../Roadmap.md`](../Roadmap.md)
- Issues: [#2](https://github.com/jwilleke/yourphr/issues/2) (standalone EPIC), [#53](https://github.com/jwilleke/yourphr/issues/53) (SMART sync, blocked), [#54](https://github.com/jwilleke/yourphr/issues/54) (display), [#171](https://github.com/jwilleke/yourphr/issues/171) (search extraction), [#196](https://github.com/jwilleke/yourphr/issues/196) (compound refs), [#198](https://github.com/jwilleke/yourphr/issues/198) / [#201](https://github.com/jwilleke/yourphr/issues/201) (document titles)
- External: `https://www.followmyhealth.com`, `https://developer.veradigm.com`

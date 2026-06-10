# Veradigm (formerly Allscripts)

Vendor reference for **Veradigm Inc.** — the company formerly known as **Allscripts Healthcare Solutions** — and the owner of [FollowMyHealth](./followmyhealth.md), the patient portal whose exports YourPHR consumes.

## Overview

**Veradigm** is a US health-IT company providing electronic health records, health data and analytics, payer/life-sciences solutions, and patient engagement. For YourPHR it matters as the **owner and gatekeeper of FollowMyHealth** and of the SMART on FHIR APIs and developer/partner program that govern live sync.

Like FollowMyHealth, Veradigm's products are **proprietary and closed-source.**

## Ownership & History

- **Allscripts** was founded in 1986 and grew into one of the largest US EHR / health-IT vendors (HQ Chicago). It expanded heavily by acquisition — including **Eclipsys (2010)**, **dbMotion**, **Jardogs / FollowMyHealth (2013)**, and **Practice Fusion (2018)**.
- **2022:** Allscripts **sold its hospital and large-physician-practice EHR business** (Sunrise, TouchWorks, and related products) to **Constellation Software / N. Harris Computer Corp**; that business was renamed **Altera Digital Health**.
- **2023:** the remaining company — data/analytics, payer & life-sciences, ambulatory EHR (Professional EHR, Practice Fusion), and patient engagement (FollowMyHealth) — **rebranded as Veradigm Inc.**
- **2024 onward:** Veradigm faced **Nasdaq listing-compliance problems** stemming from delayed financial filings (a previously disclosed software issue affecting revenue reporting). Treat the company's current listing/financial status as something to **verify against a current source** rather than rely on this doc — it is the kind of fact that ages quickly.

> **Note on identifiers (this repo):** YourPHR keeps upstream technical identifiers tied to the original project (e.g. `fasten-sources`). The "Veradigm" / "Allscripts" / "FollowMyHealth" names here refer to the external vendor, not to anything we rename in code.

## Products

Relevant to YourPHR:

- **FollowMyHealth** — patient portal / PHR; see [`followmyhealth.md`](./followmyhealth.md).
- **Professional EHR** and **Practice Fusion** — ambulatory EHRs that can be upstream sources of the data surfaced in FollowMyHealth.
- **Veradigm Connect** — the developer program / portal that issues API credentials and grants production access.

(The former Sunrise/TouchWorks hospital EHRs now belong to **Altera Digital Health** and are out of scope here.)

## Contact

| Purpose | Where |
|---|---|
| Developer program / API support | `VeradigmConnect@veradigm.com` |
| Developer portal | `https://developer.veradigm.com` |
| Partner / production-access request | `https://developer.veradigm.com/Account/PartnerRequest` |
| This project's open ticket | **#17849** (`unauthorized_client` on connect) — see [`../FHIR/fhir-testing.md`](../FHIR/fhir-testing.md) and [#53](https://github.com/jwilleke/yourphr/issues/53) |

## API & Integration

- **SMART on FHIR:** Veradigm exposes SMART-on-FHIR endpoints (FollowMyHealth's are under `fhir.followmyhealth.com` with auth at `muauthentication.followmyhealth.com`). Authorization is authorization-code + PKCE.
- **Test vs. production:** registered apps start **Test Only** and can reach only Veradigm **test** organizations (synthetic patients, no PHI). Production access is granted per-app via a **partner request** with a multi-day review. The full test-vs-real breakdown is in [`../FHIR/fhir-testing.md`](../FHIR/fhir-testing.md).
- **App types:** Veradigm's discovery advertises both public and confidential client capabilities; in practice the published metadata is **inconsistent** (see Known API Issues).

## Known API Issues

1. **`unauthorized_client` blocks the connect flow.** A registered app cannot complete authorization against FollowMyHealth; the flow never reaches token exchange. This requires Veradigm to authorize the app (ticket **#17849**). Tracked in [#53](https://github.com/jwilleke/yourphr/issues/53). **YourPHR cannot resolve this unilaterally** — it depends on Veradigm granting access.
2. **Contradictory discovery metadata.** The discovery document advertises a `client-public` capability _and_ confidential-only token-endpoint auth methods (`client-confidential-symmetric`/`-asymmetric`). Because the handshake never reaches token exchange, it is **not** yet confirmed whether a confidential client is actually required — so YourPHR has deliberately **not** built confidential-client support speculatively.
3. **`fasten-sources` is stubbed in this fork.** The upstream provider-client package went private; this repo replaces it with a local stub that has no real OAuth clients. Combined with (1), **live provider sync is non-functional** here, and **manual FHIR bundle upload is the supported import path.** Implementing a real SMART client for Veradigm/FollowMyHealth is a roadmap item.
4. **Data-shape issues** (non-US-Core Encounters, compound reference ids, document titles/bodies, text-only concepts) live with the data itself — see the Known API Issues in [`followmyhealth.md`](./followmyhealth.md).

## Relevance to YourPHR

Veradigm is the **external dependency** standing between YourPHR and live FollowMyHealth sync. Everything we can do without Veradigm's cooperation — faithfully importing and displaying the patient-initiated export — is in scope and actively worked (EPIC [#2](https://github.com/jwilleke/yourphr/issues/2)). Everything that requires Veradigm to authorize an app ([#53](https://github.com/jwilleke/yourphr/issues/53)) is **blocked on the vendor**, not on us.

## References / Related issues

- Internal: [`followmyhealth.md`](./followmyhealth.md), [`../FHIR/fhir-testing.md`](../FHIR/fhir-testing.md), [`../FHIR/fhir-test-discovery-example.md`](../FHIR/fhir-test-discovery-example.md), [`../Roadmap.md`](../Roadmap.md)
- Issues: [#2](https://github.com/jwilleke/yourphr/issues/2) (standalone EPIC), [#53](https://github.com/jwilleke/yourphr/issues/53) (SMART sync, blocked), [#54](https://github.com/jwilleke/yourphr/issues/54) (display)
- External: `https://veradigm.com`, `https://developer.veradigm.com`, `https://developer.veradigm.com/Account/PartnerRequest`

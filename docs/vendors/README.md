# Vendors

Reference notes on the external health-IT vendors whose data and APIs YourPHR interoperates with. Each doc follows the same shape: **Overview · Ownership & History · Products · Contact · API & Integration · Known API Issues · Relevance to YourPHR · References**.

| Vendor | Doc | Why it matters to YourPHR |
|---|---|---|
| **FollowMyHealth** | [`followmyhealth.md`](./followmyhealth.md) | Patient portal; its FHIR R4 export is the primary real-world (non-US-Core) dataset YourPHR is hardened against. |
| **Veradigm** (formerly **Allscripts**) | [`veradigm-allscripts.md`](./veradigm-allscripts.md) | Owns FollowMyHealth and the SMART/FHIR developer program; the external gatekeeper for live sync ([#53](https://github.com/jwilleke/yourphr/issues/53)). |

See also: [`../FHIR/fhir-testing.md`](../FHIR/fhir-testing.md) (test-vs-real environments) and [`../FHIR/fhir-test-discovery-example.md`](../FHIR/fhir-test-discovery-example.md) (a captured FollowMyHealth discovery document).

# ClientID Friction

By far the biggest issue in the project has been the otaining of ClientID to be able to make requests.

Most systems require you to register your app separately with each EHR vendor, and some make you register per health system. Here's how to get around a lot of that:

- Start with sandboxes to build and test without any approvals — use the free SMART Health IT sandbox at launch.smarthealthit.org or Epic's sandbox. These let you experiment immediately.
- For Epic specifically, patient-facing apps that stick to US Core data can use Automatic Client ID Distribution. Register once in their sandbox, meet a few criteria, and they push your client ID to hundreds of Epic organizations automatically — often live in 48 hours.

## Workarounds

Target systems that support standalone patient launch (where the patient logs into their portal and authorizes your app).

Focus on patient-facing SMART on FHIR apps — that's the cleanest way for your open source PHR to pull data from major EHRs like Epic, Cerner/Oracle Health, Athenahealth, and others without needing direct vendor partnerships everywhere.

Most systems require you to register your app separately with each EHR vendor, and some make you register per health system. Here's how to get around a lot of that:

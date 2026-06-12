# Local fhir-converter container (C-CDA → FHIR)

How to build and run the **Metriport fhir-converter** locally — the sidecar that converts C-CDA / CCD documents to FHIR R4 for manual import ([#254](https://github.com/jwilleke/yourphr/issues/254)). Validated end-to-end in #254 Phase 0.

## What it is

- The converter is **Metriport's `fhir-converter`** — a Node/Express HTTP service (Handlebars templates descended from Microsoft's open-source FHIR Converter) that maps C-CDA R2.1 documents to a FHIR R4 `Bundle`.
- It is a **separate process / container**, not part of the Go binary. YourPHR's backend POSTs the raw CCD to it and feeds the returned FHIR through the existing import pipeline. See `backend/pkg/web/handler/cda_converter.go`.
- **License:** the `fhir-converter` package is **AGPL-3.0-only**. Running it as an isolated process keeps it an independent work (no GPLv3 combined-work entanglement).
- **It is stateless** — no database, no persistence, holds no PHI beyond the in-flight request.

## Run it locally

There is **no published image** — build from source (Node 18; a few minutes the first time):

```bash
git clone --depth 1 https://github.com/metriport/metriport.git /tmp/metriport
cd /tmp/metriport/packages/fhir-converter
docker compose up --build -d        # serves on http://localhost:8777
```

If the build fails with `DeadlineExceeded` pulling the base image, pre-pull then rebuild:

```bash
docker pull node:18 && docker pull node:18-slim
docker compose up --build -d
```

Health check (the service answers `200` on `/`):

```bash
curl -s -o /dev/null -w '%{http_code}\n' http://localhost:8777/
```

Stop and remove it:

```bash
docker compose -f /tmp/metriport/packages/fhir-converter/docker-compose.yml down
```

## Convert a document

- **Endpoint:** `POST /api/convert/cda/ccd.hbs?patientId=<id>`
- **Content-Type:** `text/plain`; send the raw CDA XML as the body (use `curl --data-binary`, not `-d`, so XML whitespace is preserved).
- **`patientId`** becomes the FHIR `Patient.id`. YourPHR passes a **stable** id derived from the CDA `recordTarget` so re-imports stay idempotent — never a random value.
- **Response** is the FHIR bundle wrapped in an envelope:

```json
{ "fhirResource": { "resourceType": "Bundle", "type": "batch", "entry": [ ... ] } }
```

The backend unwraps `.fhirResource` before importing.

```bash
# fetch a synthetic CCD and convert it
curl -sL "https://raw.githubusercontent.com/HL7/CDA-ccda-2.1/master/examples/C-CDA_R2-1_CCD.xml" -o sample.xml
curl -s -X POST "http://localhost:8777/api/convert/cda/ccd.hbs?patientId=test-1" \
  -H "Content-Type: text/plain" --data-binary @sample.xml -o bundle.json
```

Use **synthetic** CCDs only (HL7 examples, Synthea). Never run real patient documents through a casual local setup, and never commit them.

## Gotchas

- **No prebuilt image** — `ghcr.io/metriport/fhir-converter` does not exist; you must build from source.
- **Port mapping is indirect** — the app binds `8080` inside the container; compose publishes it as `8777`. From the host, use `8777`. Bare `docker run` needs `-p 8777:8080`.
- **Node 18** — older than YourPHR's Node 24, but isolated in its own container, so no conflict.
- **Deterministic ids** — the converter mints resource ids with UUID v3 (content-derived), so the same CCD yields the same ids on every run (verified in Phase 0).
- Templates are baked into the image; no separate mount needed.

## How YourPHR uses it

Conversion is **opt-in** and **disabled by default** — a stock install and FHIR/NDJSON upload are unaffected. Enable it by pointing the backend at a running converter (`config.yaml` / `config.dev.yaml`):

```yaml
cda_converter:
  enabled: true
  url: 'http://localhost:8777'        # dev; 'http://yourphr-cda-converter:8777' in k8s
  timeout_seconds: 60
```

With it enabled, a `.xml` / `.ccd` upload (manual upload or the UI) is detected, converted by this sidecar, and imported — the raw CCD never leaves the instance. With it disabled, a CCD upload returns a clear "C-CDA import is not enabled" error.

## Production (Phase 2)

The sidecar is intended to be packaged and deployed like the relay (`yourphr-relay`): build/push `ghcr.io/jwilleke/yourphr-cda-converter` via a workflow mirroring `docker-relay.yaml`, deploy a Deployment + Service via `jwilleke/mj-infra-flux`, and reach it at the in-cluster DNS name. It must be **internal-only (no Ingress)** — the raw CCD is full PHI and must never leave the cluster. Tracked in [#254](https://github.com/jwilleke/yourphr/issues/254) Phase 2.

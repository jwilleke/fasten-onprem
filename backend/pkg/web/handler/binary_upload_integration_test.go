package handler

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/fastenhealth/fasten-onprem/backend/pkg/models"
	"github.com/fastenhealth/fasten-sources/clients/factory"
	sourcePkg "github.com/fastenhealth/fasten-sources/pkg"
	"github.com/sirupsen/logrus"
)

// Integration: a wrapped binary bundle must be ingestable by the REAL manual-import client — i.e. the
// production parser can extract the Patient id from it (the same call the upload handler makes after
// maybeWrapBinary). This is the integration risk the unit tests don't cover.
func TestWrapBinaryAsFHIR_IngestableByImportClient(t *testing.T) {
	pdf := []byte("%PDF-1.7 synthetic content for ingest test")
	wrapped, err := wrapBinaryAsFHIR(pdf, "report.pdf", "application/pdf", time.Date(2026, time.June, 14, 12, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("wrap: %v", err)
	}

	cred := models.SourceCredential{PlatformType: sourcePkg.PlatformTypeManual}
	client, err := factory.GetSourceClient("", context.Background(), logrus.NewEntry(logrus.New()), &cred)
	if err != nil {
		t.Fatalf("build manual client: %v", err)
	}

	patientID, _, err := client.ExtractPatientId(bytes.NewReader(wrapped))
	if err != nil {
		t.Fatalf("ExtractPatientId on wrapped bundle: %v", err)
	}
	if patientID != "upload-"+contentHash(pdf) {
		t.Errorf("import client extracted patient id %q, want %q — wrapped bundle is not ingestable", patientID, "upload-"+contentHash(pdf))
	}
}

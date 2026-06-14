package handler

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/fastenhealth/fasten-onprem/backend/pkg/document"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Uninterpreted binary upload (#255). Manual upload otherwise expects FHIR JSON/NDJSON (or a C-CDA,
// handled earlier by maybeConvertCDA). When a RAW binary is uploaded — a PDF, DICOM image, or a
// plain image — we do NOT parse or interpret it (the content is arbitrary and almost never maps to
// concrete clinical categories). Instead we wrap it, as-is, into a minimal FHIR Bundle (a Binary
// holding the bytes + a DocumentReference describing it) and feed it through the existing import
// pipeline unchanged. The result displays via the existing fhir-pdf/fhir-dicom/fhir-img viewers and
// is linkable via the existing reference graph. Provenance is the patient (DocumentReference.author
// = Patient), so the provenance resolver reports it as patient-sourced for free.

// detectBinaryType sniffs a raw upload by magic bytes and returns its MIME type if it is a supported
// uninterpreted binary, or "" for anything else (FHIR JSON/NDJSON/CDA pass through untouched).
func detectBinaryType(data []byte) string {
	switch {
	case bytes.HasPrefix(data, []byte("%PDF-")):
		return "application/pdf"
	case len(data) >= 132 && bytes.Equal(data[128:132], []byte("DICM")): // DICOM preamble: 128-byte lead-in then "DICM"
		return "application/dicom"
	case bytes.HasPrefix(data, []byte{0xFF, 0xD8, 0xFF}):
		return "image/jpeg"
	case bytes.HasPrefix(data, []byte{0x89, 'P', 'N', 'G', 0x0D, 0x0A, 0x1A, 0x0A}):
		return "image/png"
	}
	return ""
}

// contentHash returns a stable 16-hex-char id derived from the file content, so re-uploading the
// same file yields the same Patient/Binary/DocumentReference ids and upserts idempotently (no dup).
func contentHash(data []byte) string {
	sum := sha256.Sum256(data)
	return fmt.Sprintf("%x", sum[:8])
}

// wrapBinaryAsFHIR builds a minimal FHIR R4 collection Bundle for a raw uploaded binary: a Patient
// (so the import pipeline can extract a patient id), a Binary holding the base64 bytes, and a
// DocumentReference describing + pointing at it. Only explicit, known metadata is set — filename as
// the title, the detected content type, the upload time — everything else is left absent (no
// guessing). The DocumentReference is marked as a manual upload and authored by the patient.
func wrapBinaryAsFHIR(data []byte, filename, contentType string, now time.Time) ([]byte, error) {
	hash := contentHash(data)
	patientRef := "Patient/upload-" + hash
	binaryID := "binary-" + hash
	title := filename
	if title == "" {
		title = "Uploaded file"
	}

	resource := func(r map[string]interface{}) map[string]interface{} {
		return map[string]interface{}{"resource": r}
	}
	bundle := map[string]interface{}{
		"resourceType": "Bundle",
		"type":         "collection",
		"entry": []interface{}{
			resource(map[string]interface{}{
				"resourceType": "Patient",
				"id":           "upload-" + hash,
			}),
			resource(map[string]interface{}{
				"resourceType": "Binary",
				"id":           binaryID,
				"contentType":  contentType,
				"data":         base64.StdEncoding.EncodeToString(data),
			}),
			resource(map[string]interface{}{
				"resourceType": "DocumentReference",
				"id":           "docref-" + hash,
				"status":       "current",
				"identifier":   []interface{}{map[string]interface{}{"system": document.ManualUploadIDSystem, "value": filename}},
				"subject":      map[string]interface{}{"reference": patientRef},
				"author":       []interface{}{map[string]interface{}{"reference": patientRef, "display": "Uploaded by patient"}},
				"date":         now.UTC().Format(time.RFC3339),
				"type":         map[string]interface{}{"text": title},
				"content": []interface{}{map[string]interface{}{"attachment": map[string]interface{}{
					"contentType": contentType,
					"url":         "Binary/" + binaryID,
					"title":       title,
					"size":        len(data),
				}}},
			}),
		},
	}
	return json.Marshal(bundle)
}

// maybeWrapBinary inspects the uploaded file; if it is a raw binary (PDF/DICOM/image) it wraps it as
// a FHIR Bundle and returns a new temp file holding that JSON. FHIR JSON/NDJSON uploads are returned
// unchanged (rewound). Mirrors maybeConvertCDA, and runs after it in the pipeline.
func maybeWrapBinary(c *gin.Context, logger *logrus.Entry, bundleFile *os.File) (*os.File, error) {
	data, err := io.ReadAll(bundleFile)
	if err != nil {
		return nil, fmt.Errorf("reading uploaded file: %w", err)
	}
	contentType := detectBinaryType(data)
	if contentType == "" {
		if _, err := bundleFile.Seek(0, io.SeekStart); err != nil {
			return nil, err
		}
		return bundleFile, nil
	}

	filename := ""
	if fh, ferr := c.FormFile("file"); ferr == nil {
		filename = fh.Filename
	}
	logger.Infof("detected raw %s upload (%d bytes, filename=%q) — wrapping as DocumentReference + Binary", contentType, len(data), filename)

	wrapped, err := wrapBinaryAsFHIR(data, filename, contentType, time.Now())
	if err != nil {
		return nil, err
	}

	out, err := os.CreateTemp("", "fasten-binary-wrapped-*.json")
	if err != nil {
		return nil, fmt.Errorf("creating wrapped temp file: %w", err)
	}
	if _, err := out.Write(wrapped); err != nil {
		out.Close()
		return nil, fmt.Errorf("writing wrapped bundle: %w", err)
	}
	if _, err := out.Seek(0, io.SeekStart); err != nil {
		out.Close()
		return nil, err
	}

	// best-effort cleanup of the original raw-binary temp file
	origName := bundleFile.Name()
	_ = bundleFile.Close()
	_ = os.Remove(origName)
	return out, nil
}

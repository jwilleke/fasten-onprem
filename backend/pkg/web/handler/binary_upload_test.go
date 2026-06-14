package handler

import (
	"encoding/base64"
	"encoding/json"
	"testing"
	"time"

	"github.com/fastenhealth/fasten-onprem/backend/pkg/document"
)

func TestDetectBinaryType(t *testing.T) {
	dicom := make([]byte, 140)
	copy(dicom[128:], []byte("DICM"))

	cases := []struct {
		name string
		data []byte
		want string
	}{
		{"pdf", []byte("%PDF-1.7\n..."), "application/pdf"},
		{"dicom", dicom, "application/dicom"},
		{"jpeg", []byte{0xFF, 0xD8, 0xFF, 0xE0}, "image/jpeg"},
		{"png", []byte{0x89, 'P', 'N', 'G', 0x0D, 0x0A, 0x1A, 0x0A}, "image/png"},
		{"fhir json (passes through)", []byte(`{"resourceType":"Bundle"}`), ""},
		{"ndjson (passes through)", []byte(`{"resourceType":"Patient","id":"1"}`), ""},
		{"empty", []byte{}, ""},
	}
	for _, c := range cases {
		if got := detectBinaryType(c.data); got != c.want {
			t.Errorf("%s: detectBinaryType = %q, want %q", c.name, got, c.want)
		}
	}
}

func TestWrapBinaryAsFHIR(t *testing.T) {
	data := []byte("%PDF-1.7 synthetic pdf bytes")
	now := time.Date(2026, time.June, 14, 12, 0, 0, 0, time.UTC)

	out, err := wrapBinaryAsFHIR(data, "lab-report.pdf", "application/pdf", now)
	if err != nil {
		t.Fatalf("wrap: %v", err)
	}

	var bundle struct {
		ResourceType string `json:"resourceType"`
		Entry        []struct {
			Resource json.RawMessage `json:"resource"`
		} `json:"entry"`
	}
	if err := json.Unmarshal(out, &bundle); err != nil {
		t.Fatalf("unmarshal bundle: %v", err)
	}
	if bundle.ResourceType != "Bundle" || len(bundle.Entry) != 3 {
		t.Fatalf("want Bundle with 3 entries, got %s/%d", bundle.ResourceType, len(bundle.Entry))
	}

	types := map[string]map[string]interface{}{}
	for _, e := range bundle.Entry {
		var r map[string]interface{}
		json.Unmarshal(e.Resource, &r)
		types[r["resourceType"].(string)] = r
	}

	// Patient present (so ExtractPatientId works) + Binary holds the exact bytes as base64.
	if _, ok := types["Patient"]; !ok {
		t.Error("bundle missing a Patient resource")
	}
	bin := types["Binary"]
	if bin["contentType"] != "application/pdf" {
		t.Errorf("binary contentType = %v", bin["contentType"])
	}
	if decoded, _ := base64.StdEncoding.DecodeString(bin["data"].(string)); string(decoded) != string(data) {
		t.Error("binary data does not round-trip the uploaded bytes")
	}

	// DocumentReference: filename title, manual-upload marker, patient author, points at the Binary.
	dr := types["DocumentReference"]
	if dr["type"].(map[string]interface{})["text"] != "lab-report.pdf" {
		t.Errorf("docref title = %v, want filename", dr["type"])
	}
	ids := dr["identifier"].([]interface{})
	if ids[0].(map[string]interface{})["system"] != document.ManualUploadIDSystem {
		t.Errorf("docref missing manual-upload marker, got %v", ids)
	}
	author := dr["author"].([]interface{})[0].(map[string]interface{})
	if author["reference"].(string)[:8] != "Patient/" {
		t.Errorf("docref author should be the Patient, got %v", author["reference"])
	}
}

// Re-wrapping identical bytes yields identical resource ids (idempotent upsert, no duplicate).
func TestWrapBinaryAsFHIR_StableIDs(t *testing.T) {
	data := []byte("%PDF-identical")
	now := time.Date(2026, time.June, 14, 12, 0, 0, 0, time.UTC)
	later := now.Add(48 * time.Hour)

	id := func(out []byte) string {
		var b struct {
			Entry []struct {
				Resource struct {
					ResourceType string `json:"resourceType"`
					ID           string `json:"id"`
				} `json:"resource"`
			} `json:"entry"`
		}
		json.Unmarshal(out, &b)
		for _, e := range b.Entry {
			if e.Resource.ResourceType == "DocumentReference" {
				return e.Resource.ID
			}
		}
		return ""
	}

	a, _ := wrapBinaryAsFHIR(data, "x.pdf", "application/pdf", now)
	b, _ := wrapBinaryAsFHIR(data, "x.pdf", "application/pdf", later)
	if id(a) == "" || id(a) != id(b) {
		t.Errorf("DocumentReference id should be stable for identical content: %q vs %q", id(a), id(b))
	}
}

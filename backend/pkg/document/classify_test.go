package document

import (
	"encoding/json"
	"os"
	"testing"
)

// loadFixture reads the synthetic FollowMyHealth-shaped DocumentReference fixtures and wraps each as
// an InputResource keyed by its FHIR id. All values are synthetic.
func loadFixture(t *testing.T) []InputResource {
	t.Helper()
	data, err := os.ReadFile("testdata/fmh_documents.json")
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	var raws []json.RawMessage
	if err := json.Unmarshal(data, &raws); err != nil {
		t.Fatalf("unmarshal fixture: %v", err)
	}
	inputs := make([]InputResource, 0, len(raws))
	for _, r := range raws {
		var meta struct {
			ID string `json:"id"`
		}
		if err := json.Unmarshal(r, &meta); err != nil {
			t.Fatalf("unmarshal id: %v", err)
		}
		inputs = append(inputs, InputResource{
			SourceResourceType: "DocumentReference",
			SourceResourceID:   meta.ID,
			SourceID:           "synthetic-source",
			Raw:                r,
		})
	}
	return inputs
}

func byID(t *testing.T) map[string]ClassifiedDocument {
	t.Helper()
	out := map[string]ClassifiedDocument{}
	for _, d := range Classify(loadFixture(t)) {
		out[d.SourceResourceID] = d
	}
	return out
}

func TestClassify_Categories(t *testing.T) {
	got := byID(t)

	// entered-in-error is omitted; the remaining 6 are classified.
	if _, present := got["entered-in-error"]; present {
		t.Errorf("entered-in-error document should be omitted from output")
	}
	if len(got) != 6 {
		t.Errorf("expected 6 classified documents (7 fixtures minus entered-in-error), got %d", len(got))
	}

	cases := map[string]string{
		"ccd-clinical":        CategoryClinical, // application/xml C-CDA
		"html-clinical":       CategoryClinical, // text/html
		"exercise-note":       CategoryActivity, // Note tell + text/plain
		"sleep-note":          CategoryActivity, // Note tell + text/plain
		"ambiguous-plain":     CategoryClinical, // text/plain but NO Note tell -> safety bias, not buried
		"conformant-declared": CategoryClinical, // declared category honored, never re-synthesized to activity
	}
	for id, want := range cases {
		d, ok := got[id]
		if !ok {
			t.Errorf("%s: missing from output", id)
			continue
		}
		if d.Category != want {
			t.Errorf("%s: category = %q, want %q", id, d.Category, want)
		}
	}
}

func TestClassify_DisplayFields(t *testing.T) {
	got := byID(t)

	if d := got["ccd-clinical"]; d.Title != "SYNTHETIC HOSPITAL - Continuity of Care Document" || d.ContentType != "application/xml" {
		t.Errorf("ccd-clinical display fields wrong: %+v", d)
	}
	if d := got["ambiguous-plain"]; d.Title != "Some plain-text document" {
		t.Errorf("ambiguous-plain title (type.text) = %q", d.Title)
	}
	if d := got["exercise-note"]; d.Date != "2018-09-24" || d.URL == "" {
		t.Errorf("exercise-note date/url wrong: %+v", d)
	}
}

// The wearable signal is the FMH Note tell, not bare MIME: a plain-text document with no tell and no
// declared category must NOT be routed to activity (protects conformant sources).
func TestClassify_PlainTextWithoutTellStaysClinical(t *testing.T) {
	got := byID(t)
	if d := got["ambiguous-plain"]; d.Category != CategoryClinical {
		t.Errorf("plain-text with no Note tell should stay %q, got %q", CategoryClinical, d.Category)
	}
}

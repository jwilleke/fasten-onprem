// Package document classifies a patient's DocumentReference resources so real clinical documents are
// not buried under a flood of wearable/lifestyle entries.
//
// The problem (FollowMyHealth/Veradigm EHI): DocumentReference is ~85% of the export, with NO
// category and no coded type — a handful of genuine clinical documents (C-CDA / HTML) mixed with
// thousands of wearable "Notes" (text/plain "Exercise"/"Sleep" entries). Layer 1 synthesizes the
// standard DocumentReference.category the source omitted, from explicit signals only:
//
//   - a clinical MIME type (application/xml C-CDA, text/html) => a real clinical document;
//   - the FollowMyHealth interface tell identifier "…:Note,<n>" => a wearable/lifestyle note (activity).
//
// Like the Condition classifier, it FIRST honors a category a conformant source already declared
// (never re-categorizing it) and synthesizes only when none is present. The wearable signal is the
// FMH-specific Note tell, so bare MIME alone never routes a conformant source's plain-text note to
// activity. Safety bias on ambiguity is toward clinical-document — burying a possible real document
// under "activity" is worse than the reverse. One row per input (no dedup); the record marked
// entered-in-error is the only one dropped. Pure, stateless, source-agnostic — no DB, no HTTP. See
// docs/your-phr-dashboard/classification-and-display-architecture.md and docs/vendors/followmyhealth-ehi-mapping.md.
package document

import (
	"encoding/json"
	"strings"
)

// Synthesized DocumentReference.category values.
const (
	CategoryClinical     = "clinical-document" // a genuine clinical document -> Documents section
	CategoryActivity     = "activity"          // wearable / lifestyle note -> Activity (out of the clinical view)
	CategoryManualUpload = "manual-upload"     // a raw file the patient uploaded -> not interpreted
)

// ManualUploadIDSystem marks a DocumentReference as a patient-uploaded file (#255). Such uploads are
// arbitrary and not interpretable into clinical categories, so the classifier leaves them as
// CategoryManualUpload rather than guessing at clinical-vs-activity.
const ManualUploadIDSystem = "https://yourphr.org/id/manual-upload"

// InputResource is one stored DocumentReference row: authoritative type/id/source from the DB row
// plus the full FHIR JSON body.
type InputResource struct {
	SourceResourceType string
	SourceResourceID   string
	SourceID           string
	Raw                json.RawMessage
}

// ClassifiedDocument is one DocumentReference with its synthesized category and the display fields a
// document list needs. The raw record is never mutated; this is a read-time view-model.
type ClassifiedDocument struct {
	SourceResourceType string `json:"sourceResourceType"`
	SourceResourceID   string `json:"sourceResourceId"`
	SourceID           string `json:"sourceId"`
	Title              string `json:"title"`
	Category           string `json:"category"`
	ContentType        string `json:"contentType,omitempty"`
	URL                string `json:"url,omitempty"`
	Date               string `json:"date,omitempty"`
	Status             string `json:"status,omitempty"`
}

// Classify returns one ClassifiedDocument per input (in input order), except resources the record
// marks entered-in-error, which are omitted.
func Classify(resources []InputResource) []ClassifiedDocument {
	out := make([]ClassifiedDocument, 0, len(resources))
	for _, res := range resources {
		var raw rawDocumentReference
		if err := json.Unmarshal(res.Raw, &raw); err != nil {
			continue // unparseable record — skip rather than emit garbage
		}
		if strings.ToLower(raw.Status) == "entered-in-error" {
			continue // the record says this was a mistake — honor it, omit entirely
		}

		out = append(out, ClassifiedDocument{
			SourceResourceType: res.SourceResourceType,
			SourceResourceID:   res.SourceResourceID,
			SourceID:           res.SourceID,
			Title:              raw.title(),
			Category:           classify(&raw),
			ContentType:        raw.primaryContentType(),
			URL:                raw.primaryURL(),
			Date:               raw.Date,
			Status:             strings.ToLower(raw.Status),
		})
	}
	return out
}

// classify honors a category a conformant source already declared (Layer 1 never re-categorizes a
// conformant source) and synthesizes one only when the source omitted it (the FollowMyHealth case).
// A patient-uploaded file is never interpreted — it is left as CategoryManualUpload.
func classify(raw *rawDocumentReference) string {
	if isManualUpload(raw.Identifier) {
		return CategoryManualUpload
	}
	if existing := existingDocCategory(raw.Category); existing != "" {
		return existing
	}
	return synthesize(raw)
}

// synthesize derives the category from explicit signals, with a default-to-clinical safety bias.
func synthesize(raw *rawDocumentReference) string {
	if isClinicalMIME(raw.primaryContentType()) {
		return CategoryClinical // C-CDA / HTML — a real clinical document
	}
	if interfaceTell(raw.Identifier) == "Note" {
		return CategoryActivity // FollowMyHealth wearable/lifestyle "Note"
	}
	return CategoryClinical // safety bias: never bury a possible real document under activity
}

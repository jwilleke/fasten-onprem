package document

import "strings"

// Minimal FHIR R4 DocumentReference shapes — only the fields the classifier needs. JSON unmarshalling
// ignores absent fields. See docs/vendors/followmyhealth-ehi-mapping.md (DocumentReference section).

type fhirCoding struct {
	System  string `json:"system"`
	Code    string `json:"code"`
	Display string `json:"display"`
}

type fhirCodeableConcept struct {
	Text   string       `json:"text"`
	Coding []fhirCoding `json:"coding"`
}

type fhirIdentifier struct {
	System string `json:"system"`
	Value  string `json:"value"`
}

type fhirAttachment struct {
	ContentType string `json:"contentType"`
	URL         string `json:"url"`
	Title       string `json:"title"`
}

type fhirContent struct {
	Attachment *fhirAttachment `json:"attachment"`
}

type rawDocumentReference struct {
	ResourceType string                `json:"resourceType"`
	ID           string                `json:"id"`
	Status       string                `json:"status"`
	Identifier   []fhirIdentifier      `json:"identifier"`
	Category     []fhirCodeableConcept `json:"category"`
	Type         *fhirCodeableConcept  `json:"type"`
	Date         string                `json:"date"`
	Content      []fhirContent         `json:"content"`
}

// existingDocCategory returns the category a conformant source already declared, or "" when none is
// present (the FollowMyHealth case). A recognized US Core clinical-note maps explicitly; any other
// declared category is treated as a managed clinical document — a source that bothered to categorize
// is not synthesized into "activity". Honoring this is how Layer 1 never re-categorizes a conformant
// source.
func existingDocCategory(cats []fhirCodeableConcept) string {
	if len(cats) == 0 {
		return ""
	}
	for _, cc := range cats {
		for _, c := range cc.Coding {
			if strings.ToLower(c.Code) == "clinical-note" {
				return CategoryClinical
			}
		}
	}
	return CategoryClinical
}

// isClinicalMIME reports whether a content type marks a genuine clinical document (C-CDA XML or HTML),
// as opposed to the text/plain wearable "Notes".
func isClinicalMIME(contentType string) bool {
	ct := strings.ToLower(contentType)
	return strings.Contains(ct, "xml") || strings.Contains(ct, "html")
}

// interfaceTell extracts the FollowMyHealth interface object-type from an identifier value of the form
// "<n>:<guid>:<Type>,<id>" (e.g. "Note"). Returns "" when no such tell is present (conformant sources).
func interfaceTell(ids []fhirIdentifier) string {
	for _, id := range ids {
		colon := strings.LastIndex(id.Value, ":")
		if colon == -1 {
			continue
		}
		rest := id.Value[colon+1:]
		comma := strings.Index(rest, ",")
		if comma == -1 {
			continue
		}
		if t := rest[:comma]; t != "" {
			return t
		}
	}
	return ""
}

// primaryContentType returns the first content attachment's MIME type.
func (r *rawDocumentReference) primaryContentType() string {
	for _, c := range r.Content {
		if c.Attachment != nil && c.Attachment.ContentType != "" {
			return c.Attachment.ContentType
		}
	}
	return ""
}

// primaryURL returns the first content attachment's url.
func (r *rawDocumentReference) primaryURL() string {
	for _, c := range r.Content {
		if c.Attachment != nil && c.Attachment.URL != "" {
			return c.Attachment.URL
		}
	}
	return ""
}

// title prefers type.text, then a type coding display, then the attachment title.
func (r *rawDocumentReference) title() string {
	if r.Type != nil {
		if r.Type.Text != "" {
			return r.Type.Text
		}
		for _, c := range r.Type.Coding {
			if c.Display != "" {
				return c.Display
			}
		}
	}
	for _, c := range r.Content {
		if c.Attachment != nil && c.Attachment.Title != "" {
			return c.Attachment.Title
		}
	}
	return "Untitled document"
}

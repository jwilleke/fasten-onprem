package provenance

import "encoding/json"

// ExtractRequest builds a provenance Request from any FHIR resource's raw JSON, reading the superset
// of common author/informant/encounter fields and author-time elements across resource types. It is
// the generic entry point for the read path (one call works for all ~70 types); absent fields are
// simply skipped, so it never fabricates an author or a time.
//
// Author priority mirrors the bespoke wirings: asserter → recorder → requester → informationSource →
// author[]. A Patient reference among these resolves to "Self-reported" in the ladder. Author time is
// the first present of recordedDate / authoredOn / dateAsserted / issued / recorded / date.
func ExtractRequest(raw json.RawMessage, targetType, targetID, sourceLabel string) Request {
	var r struct {
		Asserter          *Reference  `json:"asserter"`
		Recorder          *Reference  `json:"recorder"`
		Requester         *Reference  `json:"requester"`
		InformationSource *Reference  `json:"informationSource"`
		Author            []Reference `json:"author"`
		Encounter         *Reference  `json:"encounter"`

		RecordedDate string `json:"recordedDate"`
		AuthoredOn   string `json:"authoredOn"`
		DateAsserted string `json:"dateAsserted"`
		Issued       string `json:"issued"`
		Recorded     string `json:"recorded"`
		Date         string `json:"date"`
	}
	_ = json.Unmarshal(raw, &r) // absent/!= fields stay zero — never fabricated

	var authors []Reference
	for _, ref := range []*Reference{r.Asserter, r.Recorder, r.Requester, r.InformationSource} {
		if ref != nil && ref.Reference != "" {
			authors = append(authors, *ref)
		}
	}
	for _, a := range r.Author {
		if a.Reference != "" {
			authors = append(authors, a)
		}
	}

	enc := Reference{}
	if r.Encounter != nil {
		enc = *r.Encounter
	}

	return Request{
		Authors:      authors,
		Encounter:    enc,
		TargetType:   targetType,
		TargetID:     targetID,
		SourceLabel:  sourceLabel,
		AuthoredTime: firstNonEmpty(r.RecordedDate, r.AuthoredOn, r.DateAsserted, r.Issued, r.Recorded, r.Date),
	}
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}

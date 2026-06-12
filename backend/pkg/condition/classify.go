// Package condition classifies a patient's Condition resources for legible display: it synthesizes
// the standard Condition.category that non-conformant sources (notably FollowMyHealth/Veradigm EHI
// exports) omit, derives a display state from clinicalStatus/abatement/verificationStatus, and
// separates real health problems from social/administrative "Personal Health Conditions".
//
// Classify is a pure, stateless derivation over the raw FHIR JSON — no database, no HTTP — so the
// rules are unit-testable in isolation. The "no guessing" principle is load-bearing: category,
// tier, and state come only from explicit signals in the record; nothing is fabricated, and the
// only resource dropped is one the record itself marks entered-in-error.
//
// Unlike medications, conditions are NOT deduped or merged — one output row per input Condition
// (report facts as the source provided them). See
// docs/your-phr-dashboard/phase-1-condition-classifier-spec.md.
package condition

import (
	"encoding/json"
	"time"
)

// Synthesized FHIR Condition.category values.
const (
	CategoryProblem       = "problem-list-item" // a health problem
	CategorySDOH          = "sdoh"              // social / personal profile
	CategoryHealthConcern = "health-concern"
)

// Display tiers.
const (
	TierClinician    = "clinician"     // coded diagnosis, clinician-attributed
	TierSelfReported = "self-reported" // real condition, patient-asserted
	TierProfile      = "profile"       // personal / social / administrative
)

// State drives where a health problem displays. Derived from clinicalStatus (primary) with
// abatement as a date source + non-conformance safety net, gated first by verificationStatus.
const (
	StateActive    = "Active"    // active / recurrence / relapse
	StateRemission = "Remission" // shown under Current, badged "in remission since <abatement>"
	StateResolved  = "Resolved"  // resolved / inactive (or abated with no status) -> Past Health Problems
	StateUnknown   = "Unknown"   // status absent/unrecognized -> shown, never assumed
	StateRuledOut  = "RuledOut"  // verificationStatus=refuted -> not a current problem
)

// InputResource is one stored Condition row: authoritative type/id/source from the DB row plus the
// full FHIR JSON body.
type InputResource struct {
	SourceResourceType string
	SourceResourceID   string
	SourceID           string
	Raw                json.RawMessage
}

// Coding is a fidelity passthrough of an original FHIR coding.
type Coding struct {
	System  string `json:"system,omitempty"`
	Code    string `json:"code,omitempty"`
	Display string `json:"display,omitempty"`
}

// ClassifiedCondition is one Condition with its synthesized category + tier + state and the display
// fields Phase 1 needs. The raw record is never mutated; this is a read-time view-model.
type ClassifiedCondition struct {
	SourceResourceType string   `json:"sourceResourceType"`
	SourceResourceID   string   `json:"sourceResourceId"`
	SourceID           string   `json:"sourceId"`
	Title              string   `json:"title"`
	Category           string   `json:"category"`
	Tier               string   `json:"tier"`
	State              string   `json:"state"`
	SelfReported       bool     `json:"selfReported"`
	ClinicalStatus     string   `json:"clinicalStatus,omitempty"`
	VerificationStatus string   `json:"verificationStatus,omitempty"`
	Onset              string   `json:"onset,omitempty"`
	Recorded           string   `json:"recorded,omitempty"`
	Abated             string   `json:"abated,omitempty"`
	Note               string   `json:"note,omitempty"`
	StandardCodings    []Coding `json:"standardCodings,omitempty"`
}

// Classify returns one ClassifiedCondition per input (in input order), except resources the record
// marks entered-in-error, which are omitted. `now` is reserved for future date-based rules and kept
// for signature symmetry with medication.Reconcile.
func Classify(resources []InputResource, now time.Time) []ClassifiedCondition {
	out := make([]ClassifiedCondition, 0, len(resources))
	for _, res := range resources {
		var raw rawCondition
		if err := json.Unmarshal(res.Raw, &raw); err != nil {
			continue // unparseable record — skip rather than emit garbage
		}

		verif := conceptCode(raw.VerificationStatus)
		if verif == "entered-in-error" {
			continue // the record says this was a mistake — honor it, omit entirely (FHIR con-5)
		}

		tier, category, selfReported := classify(&raw)
		state := resolveState(&raw, verif)

		out = append(out, ClassifiedCondition{
			SourceResourceType: res.SourceResourceType,
			SourceResourceID:   res.SourceResourceID,
			SourceID:           res.SourceID,
			Title:              raw.title(),
			Category:           category,
			Tier:               tier,
			State:              state,
			SelfReported:       selfReported,
			ClinicalStatus:     conceptCode(raw.ClinicalStatus),
			VerificationStatus: verif,
			Onset:              raw.onset(),
			Recorded:           raw.RecordedDate,
			Abated:             raw.abated(),
			Note:               raw.noteText(),
			StandardCodings:    standardCodings(raw.Code),
		})
	}
	return out
}

// classify assigns the tier + synthesized category from explicit signals, first-match-wins, with a
// default-to-health-item safety bias: only agreeing signals demote an item to the Patient Profile;
// anything ambiguous stays a health problem (never bury a possible diagnosis).
func classify(raw *rawCondition) (tier, category string, selfReported bool) {
	tell := vendorTell(raw.Identifier)
	stdCode := hasStandardCode(raw.Code)
	anyCoding := hasAnyCoding(raw.Code)
	patientAsserted := refIsType(raw.Asserter, "Patient") || (raw.Asserter == nil && refIsType(raw.Recorder, "Patient"))
	clinicianRecorder := refIsType(raw.Asserter, "Practitioner") || refIsType(raw.Recorder, "Practitioner")

	switch {
	case stdCode || tell == "HealthCondition":
		return TierClinician, CategoryProblem, false
	case anyCoding && patientAsserted:
		return TierSelfReported, CategoryProblem, true
	case !anyCoding && tell == "PersonalHealthConsideration" && !clinicianRecorder:
		return TierProfile, CategorySDOH, false
	default:
		return TierClinician, CategoryProblem, false // safety bias
	}
}

// resolveState derives the display state. verificationStatus gates first (refuted -> RuledOut);
// otherwise clinicalStatus is authoritative, with abatement implying Resolved only when no status
// is present (FHIR con-4: an abated condition's status should already be inactive/resolved/remission).
func resolveState(raw *rawCondition, verif string) string {
	if verif == "refuted" {
		return StateRuledOut
	}
	switch conceptCode(raw.ClinicalStatus) {
	case "active", "recurrence", "relapse":
		return StateActive
	case "remission":
		return StateRemission
	case "resolved", "inactive":
		return StateResolved
	default:
		if raw.abated() != "" {
			return StateResolved
		}
		return StateUnknown
	}
}

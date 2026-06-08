import * as _ from "lodash";
import {fhirVersions, ResourceType} from '../constants';
import {ReferenceModel} from '../datatypes/reference-model';
import {CodableConceptModel} from '../datatypes/codable-concept-model';
import {FastenDisplayModel} from '../fasten/fasten-display-model';
import {FastenOptions} from '../fasten/fasten-options';

// A single Provenance.agent — who was involved, how, and on whose behalf.
export interface ProvenanceAgent {
  type: string | undefined           // agent.type display (e.g. Author / Transmitter)
  who: ReferenceModel | undefined    // agent.who — Practitioner/Organization/Patient/PractitionerRole/RelatedPerson/Device (1..1)
  on_behalf_of: ReferenceModel | undefined   // agent.onBehalfOf — Organization
}

// US Core 9.0.0 Provenance (#162). Provenance is US-Core-required: the audit trail of who/what
// created a record and when. Must-Support: target[] (1..*), recorded (1..1), agent[] (1..*) with
// agent.type / agent.who / agent.onBehalfOf. occurred[x] and activity are not MS (activity shown
// when present). https://hl7.org/fhir/us/core/StructureDefinition-us-core-provenance.html
export class ProvenanceModel extends FastenDisplayModel {
  targets: ReferenceModel[] = []
  recorded: string | undefined
  agents: ProvenanceAgent[] = []
  activity: CodableConceptModel | undefined

  constructor(fhirResource: any, fhirVersion?: fhirVersions, fastenOptions?: FastenOptions) {
    super(fastenOptions)
    this.source_resource_type = ResourceType.Provenance

    this.targets = _.get(fhirResource, 'target', []);
    this.recorded = _.get(fhirResource, 'recorded');
    this.agents = _.get(fhirResource, 'agent', []).map((agent: any): ProvenanceAgent => ({
      type: _.get(agent, 'type.coding.0.display') || _.get(agent, 'type.coding.0.code') || _.get(agent, 'type.text'),
      who: _.get(agent, 'who'),
      on_behalf_of: _.get(agent, 'onBehalfOf'),
    }));
    this.activity = _.get(fhirResource, 'activity') ? new CodableConceptModel(_.get(fhirResource, 'activity')) : undefined;
  }
}

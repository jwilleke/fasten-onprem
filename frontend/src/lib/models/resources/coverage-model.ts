import {fhirVersions, ResourceType} from '../constants';
import * as _ from "lodash";
import {CodableConceptModel} from '../datatypes/codable-concept-model';
import {ReferenceModel} from '../datatypes/reference-model';
import {FastenDisplayModel} from '../fasten/fasten-display-model';
import {FastenOptions} from '../fasten/fasten-options';

export class CoverageModel extends FastenDisplayModel {
  title: string | undefined
  status: string | undefined
  coverage_type: CodableConceptModel | undefined        // US Core MS: type
  subscriber_id: string | undefined                     // US Core MS: subscriberId
  beneficiary: ReferenceModel | undefined               // US Core MS: beneficiary (Patient, 1..1)
  relationship: CodableConceptModel | undefined         // US Core MS: relationship
  period_start: string | undefined
  period_end: string | undefined
  payors: ReferenceModel[] | undefined                  // US Core MS: payor (Organization)
  coverage_class: { type?: CodableConceptModel, value?: string, name?: string }[] | undefined  // class (group/plan)
  dependent: string | undefined
  order: number | undefined

  constructor(fhirResource: any, fhirVersion?: fhirVersions, fastenOptions?: FastenOptions) {
    super(fastenOptions)
    this.source_resource_type = ResourceType.Coverage
    this.resourceDTO(fhirResource, fhirVersion || fhirVersions.R4);
  }

  commonDTO(fhirResource: any){
    this.coverage_type = _.get(fhirResource, 'type');
    this.title =
      _.get(fhirResource, 'type.text') ||
      _.get(fhirResource, 'type.coding.0.display') ||
      'Coverage';
    this.status = _.get(fhirResource, 'status', '');
    this.subscriber_id = _.get(fhirResource, 'subscriberId');
    this.beneficiary = _.get(fhirResource, 'beneficiary');
    this.relationship = _.get(fhirResource, 'relationship');
    this.period_start = _.get(fhirResource, 'period.start');
    this.period_end = _.get(fhirResource, 'period.end');
    this.payors = _.get(fhirResource, 'payor');
    this.dependent = _.get(fhirResource, 'dependent');
    this.order = _.get(fhirResource, 'order');
    this.coverage_class = _.get(fhirResource, 'class', []).map((item: any) => {
      return {
        type: _.get(item, 'type'),
        value: _.get(item, 'value'),
        name: _.get(item, 'name'),
      };
    });
  };

  resourceDTO(fhirResource: any, fhirVersion: fhirVersions){
    switch (fhirVersion) {
      case fhirVersions.DSTU2: {
        this.commonDTO(fhirResource)
        return
      }
      case fhirVersions.STU3: {
        this.commonDTO(fhirResource)
        return
      }
      case fhirVersions.R4: {
        this.commonDTO(fhirResource)
        return
      }

      default:
        throw Error('Unrecognized the fhir version property type.');
    }
  };
}

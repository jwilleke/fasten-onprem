import * as _ from "lodash";
import {fhirVersions, ResourceType} from '../constants';
import {ReferenceModel} from '../datatypes/reference-model';
import {FastenDisplayModel} from '../fasten/fasten-display-model';
import {FastenOptions} from '../fasten/fasten-options';

// A node in the QuestionnaireResponse answer tree: a question (text) with its answer(s) and/or
// nested sub-items (groups). Answers are flattened to display strings for rendering.
export interface QuestionnaireResponseItem {
  linkId: string | undefined
  text: string | undefined
  answers: string[]
  items: QuestionnaireResponseItem[]
}

// US Core 9.0.0 QuestionnaireResponse (#160). Must-Support: questionnaire (canonical, 1..1),
// status (1..1), subject (1..1), authored (1..1), author (0..1), and the item[] answer tree
// (item.linkId / item.text / item.answer.value[x] / nested item.item).
// https://hl7.org/fhir/us/core/StructureDefinition-us-core-questionnaireresponse.html
export class QuestionnaireResponseModel extends FastenDisplayModel {
  status: string | undefined
  questionnaire: string | undefined        // canonical reference to the Questionnaire being answered
  subject: ReferenceModel | undefined
  authored: string | undefined
  author: ReferenceModel | undefined
  items: QuestionnaireResponseItem[] = []

  constructor(fhirResource: any, fhirVersion?: fhirVersions, fastenOptions?: FastenOptions) {
    super(fastenOptions)
    this.source_resource_type = ResourceType.QuestionnaireResponse

    this.status = _.get(fhirResource, 'status');
    // R4 uses `questionnaire` (canonical string); STU3/DSTU2 used a Reference — accept either.
    this.questionnaire = _.get(fhirResource, 'questionnaire') || _.get(fhirResource, 'questionnaire.reference');
    this.subject = _.get(fhirResource, 'subject');
    this.authored = _.get(fhirResource, 'authored');
    this.author = _.get(fhirResource, 'author');
    this.items = QuestionnaireResponseModel.parseItems(_.get(fhirResource, 'item', []));
  }

  static parseItems(rawItems: any[]): QuestionnaireResponseItem[] {
    return (rawItems || []).map((item: any): QuestionnaireResponseItem => ({
      linkId: _.get(item, 'linkId'),
      text: _.get(item, 'text'),
      answers: _.get(item, 'answer', [])
        .map((answer: any) => QuestionnaireResponseModel.answerDisplay(answer))
        .filter((s: string | undefined): s is string => !!s),
      items: QuestionnaireResponseModel.parseItems(_.get(item, 'item', [])),
    }));
  }

  // answer.value[x] -> a human-readable string.
  static answerDisplay(answer: any): string | undefined {
    if (_.has(answer, 'valueString')) { return _.get(answer, 'valueString') }
    if (_.has(answer, 'valueBoolean')) { return _.get(answer, 'valueBoolean') ? 'Yes' : 'No' }
    if (_.has(answer, 'valueInteger')) { return String(_.get(answer, 'valueInteger')) }
    if (_.has(answer, 'valueDecimal')) { return String(_.get(answer, 'valueDecimal')) }
    if (_.has(answer, 'valueDate')) { return _.get(answer, 'valueDate') }
    if (_.has(answer, 'valueDateTime')) { return _.get(answer, 'valueDateTime') }
    if (_.has(answer, 'valueTime')) { return _.get(answer, 'valueTime') }
    if (_.has(answer, 'valueUri')) { return _.get(answer, 'valueUri') }
    if (_.has(answer, 'valueCoding')) {
      return _.get(answer, 'valueCoding.display') || _.get(answer, 'valueCoding.code');
    }
    if (_.has(answer, 'valueQuantity')) {
      return [_.get(answer, 'valueQuantity.value'), _.get(answer, 'valueQuantity.unit')].filter((v) => v != null).join(' ');
    }
    if (_.has(answer, 'valueReference')) {
      return _.get(answer, 'valueReference.display') || _.get(answer, 'valueReference.reference');
    }
    return undefined;
  }
}

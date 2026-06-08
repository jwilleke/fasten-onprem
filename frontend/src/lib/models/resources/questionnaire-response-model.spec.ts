import { QuestionnaireResponseModel } from './questionnaire-response-model';
import * as example1Fixture from '../../fixtures/r4/resources/questionnaireResponse/example1.json';

describe('QuestionnaireResponseModel', () => {
  it('should create an instance', () => {
    expect(new QuestionnaireResponseModel({})).toBeTruthy();
  });

  it('parses metadata + the nested item/answer tree from r4 example1', () => {
    const m = new QuestionnaireResponseModel(example1Fixture);
    expect(m.status).toBe('completed');
    expect(m.subject).toEqual({ reference: 'Patient/f201', display: 'Roel' });
    expect(m.authored).toBe('2013-06-18T00:00:00+01:00');
    expect(m.author).toEqual({ reference: 'Practitioner/f201' });

    // top-level item "1" is a group with a nested question "1.1"
    expect(m.items.length).toBeGreaterThan(0);
    const group = m.items[0];
    expect(group.linkId).toBe('1');
    expect(group.items.length).toBeGreaterThan(0);
    const question = group.items[0];
    expect(question.text).toBe('Do you have allergies?');
    expect(question.answers).toEqual(['I am allergic to house dust']);
  });

  it('renders the supported answer value[x] types as display strings', () => {
    expect(QuestionnaireResponseModel.answerDisplay({ valueString: 'hello' })).toBe('hello');
    expect(QuestionnaireResponseModel.answerDisplay({ valueBoolean: true })).toBe('Yes');
    expect(QuestionnaireResponseModel.answerDisplay({ valueBoolean: false })).toBe('No');
    expect(QuestionnaireResponseModel.answerDisplay({ valueInteger: 3 })).toBe('3');
    expect(QuestionnaireResponseModel.answerDisplay({ valueDate: '2020-01-02' })).toBe('2020-01-02');
    expect(QuestionnaireResponseModel.answerDisplay({ valueCoding: { code: 'x', display: 'Severe' } })).toBe('Severe');
    expect(QuestionnaireResponseModel.answerDisplay({ valueQuantity: { value: 5, unit: 'mg' } })).toBe('5 mg');
    expect(QuestionnaireResponseModel.answerDisplay({})).toBeUndefined();
  });

  it('defaults to empty items / undefined metadata for a bare resource', () => {
    const m = new QuestionnaireResponseModel({});
    expect(m.items).toEqual([]);
    expect(m.status).toBeUndefined();
    expect(m.subject).toBeUndefined();
  });
});

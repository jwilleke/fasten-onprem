import { ImmunizationModel } from './immunization-model';
import * as example1Fixture from "../../fixtures/r4/resources/immunization/example1.json"

describe('ImmunizationModel', () => {
  it('should create an instance', () => {
    expect(new ImmunizationModel({})).toBeTruthy();
  });
  describe('with r4', () => {

    it('should parse example1.json', () => {
      const expected = new ImmunizationModel({})

      expected.title = 'Fluvax (Influenza)'
      expected.status = 'completed'
      expected.primary_source = true // US Core MS: example1 has primarySource: true
      expected.provided_date  = '2013-01-10'
      // manufacturerText: string | undefined
      expected.has_lot_number = true
      expected.lot_number = 'AAJN11K'
      expected.lot_number_expiration_date = '2015-02-15'
      expected.has_dose_quantity = true
      expected.dose_quantity = { value: 5, system: 'http://unitsofmeasure.org', code: 'mg' }
      // requester: string | undefined
      // reported: string | undefined
      // performer: string | undefined
      expected.route = [ { system: 'http://terminology.hl7.org/CodeSystem/v3-RouteOfAdministration', code: 'IM', display: 'Injection, intramuscular' }]
      expected.has_route = true
      expected.site =  [{ system: 'http://terminology.hl7.org/CodeSystem/v3-ActSite', code: 'LA', display: 'left arm' } ]
      expected.has_site = true
      expected.patient = { reference: 'Patient/example' }
      expected.note = [ { text: 'Notes on adminstration of vaccine' } ]

      expect(new ImmunizationModel(example1Fixture)).toEqual(expected);
    });

    // US Core 9.0.0 Must-Support: statusReason (why a dose was not given) and primarySource.
    it('should capture statusReason and primarySource (a not-done immunization)', () => {
      const model = new ImmunizationModel({
        resourceType: 'Immunization',
        status: 'not-done',
        statusReason: { text: 'Patient allergy to vaccine component' },
        primarySource: false,
        vaccineCode: { text: 'Influenza vaccine' },
        patient: { reference: 'Patient/example' },
      })
      expect(model.status).toEqual('not-done')
      expect(model.status_reason).toEqual({ text: 'Patient allergy to vaccine component' })
      expect(model.primary_source).toEqual(false)
      expect(model.title).toEqual('Influenza vaccine')
    });
  })
});

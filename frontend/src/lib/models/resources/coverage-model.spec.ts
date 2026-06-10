import { CoverageModel } from './coverage-model';
import * as example1Fixture from "../../fixtures/r4/resources/coverage/example1.json"

describe('CoverageModel', () => {
  it('should create an instance', () => {
    expect(new CoverageModel({})).toBeTruthy();
  });

  describe('with r4', () => {
    it('should parse example1.json US Core Must-Support elements', () => {
      const model = new CoverageModel(example1Fixture)
      expect(model.status).toEqual('active')
      expect(model.title).toEqual('extended healthcare') // type.coding[0].display
      expect(model.beneficiary).toEqual({ reference: 'Patient/4' })
      expect(model.payors).toEqual([{ reference: 'Organization/2' }])
      expect(model.period_start).toEqual('2011-05-23')
      expect(model.period_end).toEqual('2012-05-23')
      // class (group): surfaced for the card
      expect(model.coverage_class?.length).toBeGreaterThan(0)
      expect(model.coverage_class?.[0]?.value).toEqual('CB135')
    });
  })
});

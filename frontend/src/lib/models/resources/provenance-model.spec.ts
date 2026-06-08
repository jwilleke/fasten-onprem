import { ProvenanceModel } from './provenance-model';
import * as example1Fixture from '../../fixtures/r4/resources/provenance/example1.json';

describe('ProvenanceModel', () => {
  it('should create an instance', () => {
    expect(new ProvenanceModel({})).toBeTruthy();
  });

  it('parses target, recorded, agents and activity from r4 example1', () => {
    const m = new ProvenanceModel(example1Fixture);
    expect(m.recorded).toBe('2016-05-26T12:46:25Z');
    expect(m.targets).toEqual([{ reference: 'DocumentReference/episode-summary' }]);
    expect(m.agents.length).toBe(2);
    expect(m.agents[0].type).toBe('Author');
    expect(m.agents[0].who).toEqual({ reference: 'Practitioner/practitioner-1', display: 'Dr. Adam Everyman' });
    expect(m.agents[0].on_behalf_of).toEqual({ reference: 'Organization/org-1', display: 'Good Health Clinic' });
    expect(m.agents[1].type).toBe('Transmitter');
    expect(m.agents[1].who).toEqual({ reference: 'Organization/org-1', display: 'Good Health Clinic' });
    expect(m.agents[1].on_behalf_of).toBeUndefined();
    expect(m.activity?.coding?.[0]?.display).toBe('create');
  });

  it('defaults to empty arrays / undefined for a bare resource', () => {
    const m = new ProvenanceModel({});
    expect(m.targets).toEqual([]);
    expect(m.agents).toEqual([]);
    expect(m.recorded).toBeUndefined();
    expect(m.activity).toBeUndefined();
  });
});

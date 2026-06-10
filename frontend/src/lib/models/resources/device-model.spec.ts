import { DeviceModel } from './device-model';
import {AdverseEventModel} from './adverse-event-model';
import {CodableConceptModel} from '../datatypes/codable-concept-model';
import * as fixture from "../../fixtures/r4/resources/device/example1.json"
import * as fixture2 from "../../fixtures/r4/resources/device/example2.json"


describe('DeviceModel', () => {
  it('should create an instance', () => {
    expect(new DeviceModel({})).toBeTruthy();
  });

  describe('with r4', () => {

    it('should parse example1.json', () => {
      const expected = new DeviceModel({})


      expect(new DeviceModel(fixture)).toEqual(expected);
    });

    it('should parse example2.json', () => {
      const expected = new DeviceModel({})
      expected.status = 'active'


      expect(new DeviceModel(fixture2)).toEqual(expected);
    });

    // US Core 9.0.0 (Implantable Device) Must-Support: type, udiCarrier.deviceIdentifier, patient.
    it('should capture patient + implantable-device identifiers', () => {
      const model = new DeviceModel({
        resourceType: 'Device',
        status: 'active',
        type: { coding: [{ system: 'http://snomed.info/sct', code: '468063009', display: 'Coronary artery stent' }], text: 'Coronary stent' },
        patient: { reference: 'Patient/example' },
        distinctIdentifier: 'A9999',
        lotNumber: 'LOT123',
        serialNumber: 'SN456',
        manufactureDate: '2022-01-15',
        udiCarrier: { deviceIdentifier: '00844588003288', carrierHRF: '(01)00844588003288(17)220101' },
      })
      expect(model.model).toEqual('Coronary stent')
      expect(model.patient).toEqual({ reference: 'Patient/example' })
      expect(model.distinct_identifier).toEqual('A9999')
      expect(model.lot_number).toEqual('LOT123')
      expect(model.serial_number).toEqual('SN456')
      expect(model.manufacture_date).toEqual('2022-01-15')
      expect(model.get_udi).toEqual('00844588003288')
    });
  })

});

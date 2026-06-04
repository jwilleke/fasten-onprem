import { quantityR4Factory } from "src/lib/fixtures/factories/r4/datatypes/quantity-r4-factory";
import { QuantityModel } from "./quantity-model";

describe('QuantityModel', () => {
  it('should create an instance', () => {
    expect(new QuantityModel({})).toBeTruthy();
  });

  it('returns the correct visualization types', () => {
    expect(new QuantityModel({}).visualizationTypes()).toEqual(['bar', 'table']);
  });

  it('returns the correct display', () => {
    const quantity = new QuantityModel(quantityR4Factory.value(8).comparator('<').build());
    const quantity2 = new QuantityModel(quantityR4Factory.value(8.2).comparator('<').unit('g').build());
    const quantity3 = new QuantityModel(quantityR4Factory.value(9.5).unit('g').build());

    expect(quantity.display()).toEqual('< 8')
    expect(quantity2.display()).toEqual('< 8.2 g')
    expect(quantity3.display()).toEqual('9.5 g')
  });

  describe('valueObject', () => {
    it('sets value if there is no comparator', () => {
      const quantity = new QuantityModel(quantityR4Factory.value(6.3).build());

      expect(quantity.valueObject()).toEqual({ value: 6.3 });
    });

    it('sets range correctly if there is a comparator', () => {
      const quantity = new QuantityModel(quantityR4Factory.value(8).comparator('<').build());
      const quantity2 = new QuantityModel(quantityR4Factory.value(8).comparator('>').build());

      expect(quantity.valueObject()).toEqual({ range: { low: null, high: 8 } });
      expect(quantity2.valueObject()).toEqual({ range: { low: 8, high: null } });
    });
  });

});

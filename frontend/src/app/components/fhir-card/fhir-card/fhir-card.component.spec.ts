import { ComponentFixture, TestBed } from '@angular/core/testing';

import { FhirCardComponent } from './fhir-card.component';
import {FhirCardOutletDirective} from './fhir-card-outlet.directive';

describe('FhirResourceComponent', () => {
  let component: FhirCardComponent;
  let fixture: ComponentFixture<FhirCardComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ FhirCardComponent, FhirCardOutletDirective ],
    })
    .compileComponents();

    fixture = TestBed.createComponent(FhirCardComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('renders the "who said this" provenance line when the model carries one', () => {
    // No source_resource_type → typeLookup returns null → no dynamic card loaded, just the provenance line.
    component.displayModel = {provenance: {kind: 'source', display: 'Source: FollowMyHealth', level: 4}} as any;
    component.ngOnChanges({} as any);
    fixture.detectChanges();
    const el: HTMLElement = fixture.nativeElement;
    expect(el.querySelector('.fhir-card-provenance')?.textContent).toContain('Source: FollowMyHealth');
  });

  it('shows no provenance line when the model has none', () => {
    component.displayModel = {source_resource_type: undefined} as any;
    component.ngOnChanges({} as any);
    fixture.detectChanges();
    const el: HTMLElement = fixture.nativeElement;
    expect(el.querySelector('.fhir-card-provenance')).toBeNull();
  });
});

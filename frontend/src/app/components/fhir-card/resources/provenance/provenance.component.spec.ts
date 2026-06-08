import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ProvenanceComponent } from './provenance.component';
import {RouterTestingModule} from '@angular/router/testing';

describe('ProvenanceComponent', () => {
  let component: ProvenanceComponent;
  let fixture: ComponentFixture<ProvenanceComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ RouterTestingModule, ProvenanceComponent ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(ProvenanceComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

import { ComponentFixture, TestBed } from '@angular/core/testing';

import { QuestionnaireResponseComponent } from './questionnaire-response.component';
import {RouterTestingModule} from '@angular/router/testing';

describe('QuestionnaireResponseComponent', () => {
  let component: QuestionnaireResponseComponent;
  let fixture: ComponentFixture<QuestionnaireResponseComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ RouterTestingModule, QuestionnaireResponseComponent ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(QuestionnaireResponseComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

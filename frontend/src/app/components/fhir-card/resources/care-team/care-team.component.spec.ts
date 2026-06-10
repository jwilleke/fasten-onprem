import { ComponentFixture, TestBed } from '@angular/core/testing';
import { NgbCollapseModule } from '@ng-bootstrap/ng-bootstrap';
import { RouterTestingModule } from '@angular/router/testing';

import { CareTeamComponent } from './care-team.component';

describe('CareTeamComponent', () => {
  let component: CareTeamComponent;
  let fixture: ComponentFixture<CareTeamComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [CareTeamComponent, NgbCollapseModule, RouterTestingModule]
    })
    .compileComponents();

    fixture = TestBed.createComponent(CareTeamComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

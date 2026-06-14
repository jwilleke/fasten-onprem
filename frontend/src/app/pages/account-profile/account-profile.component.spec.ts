import {ComponentFixture, TestBed} from '@angular/core/testing';
import {CommonModule} from '@angular/common';
import {FormsModule} from '@angular/forms';
import {of, throwError} from 'rxjs';
import {RouterTestingModule} from '@angular/router/testing';

import {AccountProfileComponent} from './account-profile.component';
import {FastenApiService} from '../../services/fasten-api.service';
import {ReportHeaderComponent} from 'src/app/components/report-header/report-header.component';

describe('AccountProfileComponent', () => {
  let component: AccountProfileComponent;
  let fixture: ComponentFixture<AccountProfileComponent>;
  let api: jasmine.SpyObj<FastenApiService>;

  beforeEach(async () => {
    api = jasmine.createSpyObj('FastenApiService', ['getCurrentUser', 'deleteAccount', 'getSummary', 'getResources', 'changePassword']);
    api.getCurrentUser.and.returnValue(of({username: 'jim', full_name: 'Jim Willeke', email: 'jim@example.com', role: 'admin'}));
    api.deleteAccount.and.returnValue(of(true));
    api.changePassword.and.returnValue(of(true));
    // ReportHeaderComponent (rendered via <report-header>) calls these on init.
    api.getSummary.and.returnValue(of({sources: []} as any));
    api.getResources.and.returnValue(of([]));

    await TestBed.configureTestingModule({
      declarations: [AccountProfileComponent, ReportHeaderComponent],
      imports: [CommonModule, FormsModule, RouterTestingModule],
      providers: [{provide: FastenApiService, useValue: api}],
    }).compileComponents();

    fixture = TestBed.createComponent(AccountProfileComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('creates and loads the current user', () => {
    expect(component).toBeTruthy();
    expect(api.getCurrentUser).toHaveBeenCalled();
    expect(component.user.username).toBe('jim');
    expect(component.loading.page).toBeFalse();
  });

  it('computes initials from the full name', () => {
    expect(component.initials).toBe('JW');
  });

  it('falls back to the first two letters when there is only one name part', () => {
    component.user = {username: 'jim'};
    expect(component.initials).toBe('JI');
  });

  it('delegates account deletion to the API', () => {
    component.deleteAccount();
    expect(api.deleteAccount).toHaveBeenCalled();
  });

  it('rejects a password change when the new passwords do not match (no API call)', () => {
    component.pw = {current: 'old', next: 'new12345', confirm: 'different'};
    component.changePassword();
    expect(api.changePassword).not.toHaveBeenCalled();
    expect(component.pwError).toContain('do not match');
  });

  it('changes the password and clears the form on success', () => {
    component.pw = {current: 'old', next: 'new12345', confirm: 'new12345'};
    component.changePassword();
    expect(api.changePassword).toHaveBeenCalledWith('old', 'new12345');
    expect(component.pwSuccess).toBeTrue();
    expect(component.pw.current).toBe('');
  });

  it('surfaces the server error message on failure', () => {
    api.changePassword.and.returnValue(throwError(() => ({error: {error: 'current password is incorrect'}})));
    component.pw = {current: 'wrong', next: 'new12345', confirm: 'new12345'};
    component.changePassword();
    expect(component.pwError).toBe('current password is incorrect');
    expect(component.pwSuccess).toBeFalse();
  });
});

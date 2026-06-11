import {ComponentFixture, TestBed} from '@angular/core/testing';
import {of, throwError} from 'rxjs';
import {RouterTestingModule} from '@angular/router/testing';
import {HttpClient} from '@angular/common/http';

import {ProfileSummaryWidgetComponent} from './profile-summary-widget.component';
import {FastenApiService} from '../../services/fasten-api.service';
import {HTTP_CLIENT_TOKEN} from '../../dependency-injection';
import {ResourceFhir} from '../../models/fasten/resource_fhir';
import {DashboardWidgetConfig} from '../../models/widget/dashboard-widget-config';

function resource(partial: Partial<ResourceFhir>): ResourceFhir {
  return new ResourceFhir({
    source_id: 'src-1',
    source_resource_type: 'Condition',
    source_resource_id: partial.source_resource_id || 'r1',
    sort_title: 'Resource',
    sort_date: null,
    resource_raw: {resourceType: 'Condition', id: partial.source_resource_id || 'r1'},
    ...partial,
  });
}

function configFor(from: string, extra: Partial<DashboardWidgetConfig> = {}): DashboardWidgetConfig {
  return {
    item_type: 'profile-summary-widget',
    title_text: 'Problems',
    description_text: '',
    width: 4,
    height: 4,
    queries: [{q: {select: [], from, where: {}}}],
    ...extra,
  } as DashboardWidgetConfig;
}

describe('ProfileSummaryWidgetComponent', () => {
  let component: ProfileSummaryWidgetComponent;
  let fixture: ComponentFixture<ProfileSummaryWidgetComponent>;
  let api: jasmine.SpyObj<FastenApiService>;

  beforeEach(async () => {
    api = jasmine.createSpyObj('FastenApiService', ['getResources']);
    api.getResources.and.returnValue(of([]));

    await TestBed.configureTestingModule({
      imports: [ProfileSummaryWidgetComponent, RouterTestingModule],
      providers: [
        {provide: FastenApiService, useValue: api},
        {provide: HTTP_CLIENT_TOKEN, useClass: HttpClient},
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(ProfileSummaryWidgetComponent);
    component = fixture.componentInstance;
  });

  it('creates and fetches the configured resource type', () => {
    component.widgetConfig = configFor('Condition');
    fixture.detectChanges();
    expect(component).toBeTruthy();
    expect(api.getResources).toHaveBeenCalledWith('Condition');
    expect(component.loading).toBeFalse();
  });

  it('is empty (no fetch) when the config carries no resource type', () => {
    component.widgetConfig = configFor('');
    fixture.detectChanges();
    expect(api.getResources).not.toHaveBeenCalled();
    expect(component.isEmpty).toBeTrue();
    expect(component.loading).toBeFalse();
  });

  it('counts all results, shows the recent N newest-first, and deep-links each item', () => {
    api.getResources.and.returnValue(of([
      resource({source_resource_id: 'old', sort_date: new Date('2020-01-01') as any}),
      resource({source_resource_id: 'new', sort_date: new Date('2024-01-01') as any}),
    ]));
    component.widgetConfig = configFor('Condition');
    fixture.detectChanges();

    expect(component.totalCount).toBe(2);
    expect(component.rows.length).toBe(2);
    // newest first
    expect(component.rows[0].link).toEqual(['/explore', 'src-1', 'resource', 'new']);
    expect(component.rows[1].link).toEqual(['/explore', 'src-1', 'resource', 'old']);
    expect(component.isEmpty).toBeFalse();
  });

  it('caps the list at MAX_ROWS while keeping the full count', () => {
    const many = Array.from({length: 10}, (_, i) =>
      resource({source_resource_id: `r${i}`, sort_date: new Date(2000 + i, 0, 1) as any}));
    api.getResources.and.returnValue(of(many));
    component.widgetConfig = configFor('Condition');
    fixture.detectChanges();

    expect(component.totalCount).toBe(10);
    expect(component.rows.length).toBe(6);
  });

  it('surfaces a status badge from the parsed model where present', () => {
    api.getResources.and.returnValue(of([
      resource({
        source_resource_id: 'c1',
        resource_raw: {
          resourceType: 'Condition',
          id: 'c1',
          clinicalStatus: {coding: [{code: 'active'}]},
        } as any,
      }),
    ]));
    component.widgetConfig = configFor('Condition');
    fixture.detectChanges();
    expect(component.rows[0].status).toBe('active');
  });

  it('renders the missing-data marker (empty) for a category with no records', () => {
    api.getResources.and.returnValue(of([]));
    component.widgetConfig = configFor('Condition');
    fixture.detectChanges();
    expect(component.isEmpty).toBeTrue();
    const el: HTMLElement = fixture.nativeElement;
    expect(el.querySelector('app-missing-data')).toBeTruthy();
  });

  it('exposes view-all only when a route is configured', () => {
    component.widgetConfig = configFor('Condition');
    fixture.detectChanges();
    expect(component.viewAllRoute).toBeUndefined();

    const withRoute = TestBed.createComponent(ProfileSummaryWidgetComponent);
    withRoute.componentInstance.widgetConfig = configFor('Condition', {
      parsing: {view_all_route: '/medical-history'} as any,
    });
    withRoute.detectChanges();
    expect(withRoute.componentInstance.viewAllRoute).toBe('/medical-history');
  });

  it('is empty when the API errors', () => {
    api.getResources.and.returnValue(throwError(() => new Error('boom')));
    component.widgetConfig = configFor('Condition');
    fixture.detectChanges();
    expect(component.isEmpty).toBeTrue();
    expect(component.loading).toBeFalse();
  });
});

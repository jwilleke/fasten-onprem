import { TestBed } from '@angular/core/testing';

import { DashboardPreferencesService } from './dashboard-preferences.service';

describe('DashboardPreferencesService', () => {
  let service: DashboardPreferencesService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(DashboardPreferencesService);
    localStorage.removeItem('dashboard_tile_order');
    localStorage.removeItem('dashboard_tile_order.testuser');
  });

  afterEach(() => {
    localStorage.removeItem('dashboard_tile_order');
    localStorage.removeItem('dashboard_tile_order.testuser');
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  it('returns null when no order is saved', () => {
    expect(service.getTileOrder()).toBeNull();
    expect(service.hasCustomTileOrder()).toBeFalse();
  });

  it('round-trips a tile order', () => {
    service.setTileOrder(['b', 'a', 'c']);
    expect(service.getTileOrder()).toEqual(['b', 'a', 'c']);
    expect(service.hasCustomTileOrder()).toBeTrue();
  });

  it('scopes orders per user', () => {
    service.setTileOrder(['a'], 'testuser');
    expect(service.getTileOrder('testuser')).toEqual(['a']);
    expect(service.getTileOrder()).toBeNull();
  });

  it('reset removes the saved order', () => {
    service.setTileOrder(['a', 'b']);
    service.resetTileOrder();
    expect(service.getTileOrder()).toBeNull();
  });

  it('ignores corrupt stored values', () => {
    localStorage.setItem('dashboard_tile_order', 'not-json{');
    expect(service.getTileOrder()).toBeNull();
  });
});

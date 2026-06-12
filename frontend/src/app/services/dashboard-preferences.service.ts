import {Injectable} from '@angular/core';

// Persists the user's dashboard tile order locally (per browser, scoped per
// user when known). Backend user_settings persistence is a follow-up (#209).
@Injectable({
  providedIn: 'root'
})
export class DashboardPreferencesService {
  private static STORAGE_KEY_PREFIX = 'dashboard_tile_order';

  private storageKey(userId?: string): string {
    return userId ? `${DashboardPreferencesService.STORAGE_KEY_PREFIX}.${userId}` : DashboardPreferencesService.STORAGE_KEY_PREFIX;
  }

  getTileOrder(userId?: string): string[] | null {
    try {
      const raw = localStorage.getItem(this.storageKey(userId));
      if (!raw) return null;
      const parsed = JSON.parse(raw);
      return Array.isArray(parsed) ? parsed.filter((id) => typeof id === 'string') : null;
    } catch (e) {
      return null;
    }
  }

  setTileOrder(tileIds: string[], userId?: string): void {
    localStorage.setItem(this.storageKey(userId), JSON.stringify(tileIds));
  }

  resetTileOrder(userId?: string): void {
    localStorage.removeItem(this.storageKey(userId));
  }

  hasCustomTileOrder(userId?: string): boolean {
    return this.getTileOrder(userId) != null;
  }
}

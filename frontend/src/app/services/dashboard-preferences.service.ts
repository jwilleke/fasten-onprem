import {Injectable} from '@angular/core';

// Persists the user's dashboard tile order + colors locally (per browser, scoped
// per user when known). Backend user_settings persistence is a follow-up (#209).
@Injectable({
  providedIn: 'root'
})
export class DashboardPreferencesService {
  private static STORAGE_KEY_PREFIX = 'dashboard_tile_order';
  private static COLORS_KEY_PREFIX = 'dashboard_tile_colors';

  private storageKey(userId?: string): string {
    return userId ? `${DashboardPreferencesService.STORAGE_KEY_PREFIX}.${userId}` : DashboardPreferencesService.STORAGE_KEY_PREFIX;
  }

  private colorsKey(userId?: string): string {
    return userId ? `${DashboardPreferencesService.COLORS_KEY_PREFIX}.${userId}` : DashboardPreferencesService.COLORS_KEY_PREFIX;
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

  // --- tile colors (map of tileId -> palette color name) ---

  getTileColors(userId?: string): Record<string, string> {
    try {
      const raw = localStorage.getItem(this.colorsKey(userId));
      if (!raw) return {};
      const parsed = JSON.parse(raw);
      return parsed && typeof parsed === 'object' ? parsed : {};
    } catch (e) {
      return {};
    }
  }

  setTileColor(tileId: string, color: string, userId?: string): void {
    const colors = this.getTileColors(userId);
    colors[tileId] = color;
    localStorage.setItem(this.colorsKey(userId), JSON.stringify(colors));
  }

  resetTileColors(userId?: string): void {
    localStorage.removeItem(this.colorsKey(userId));
  }

  hasCustomTileColors(userId?: string): boolean {
    return Object.keys(this.getTileColors(userId)).length > 0;
  }
}

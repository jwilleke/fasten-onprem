import {Component, OnInit} from '@angular/core';
import {Router} from '@angular/router';
import {CdkDragDrop, moveItemInArray} from '@angular/cdk/drag-drop';
import {FastenApiService} from '../../services/fasten-api.service';
import {AuthService} from '../../services/auth.service';
import {DashboardPreferencesService} from '../../services/dashboard-preferences.service';
import {Summary} from '../../models/fasten/summary';
import {ClassifiedCondition} from '../../models/fasten/classified-condition';

// A large-icon category tile. Labels pair plain language (label) with the
// standardized clinical term (clinicalLabel) so the record stays legible to
// the patient and useful to providers (#262).
// The palette a patient can pick a tile color from (matches the SCSS .tile-color-* classes).
export const TILE_PALETTE = ['amber', 'blue', 'red', 'green', 'teal', 'purple', 'pink', 'gray']

// A large-icon category tile. Labels pair plain language (label) with the
// standardized clinical term (clinicalLabel) so the record stays legible to
// the patient and useful to providers (#262). color is the default palette
// hue (patient-overridable); unit is the noun in the count sub-line.
export interface DashboardTile {
  id: string
  label: string
  clinicalLabel: string
  icon: string
  route: string
  resourceTypes: string[]
  count: number
  color: string
  unit: string
  // countKey tiles get their count from the condition classifier (not the summary resource counts).
  countKey?: 'concerns' | 'profile'
}

export const DEFAULT_TILES: DashboardTile[] = [
  {id: 'concerns', label: 'Medical Concerns', clinicalLabel: 'Active health problems', icon: 'fa-solid fa-heart-pulse', route: '/medical-history', resourceTypes: [], count: 0, color: 'red', unit: 'active', countKey: 'concerns'},
  {id: 'patient-profile', label: 'Patient Profile', clinicalLabel: 'Personal & social info', icon: 'fa-solid fa-id-card', route: '/medical-history', resourceTypes: [], count: 0, color: 'gray', unit: 'items', countKey: 'profile'},
  {id: 'medications', label: 'Medications', clinicalLabel: 'Prescriptions & medication statements', icon: 'fa-solid fa-pills', route: '/medications', resourceTypes: ['MedicationRequest', 'MedicationStatement', 'Medication', 'MedicationAdministration', 'MedicationDispense'], count: 0, color: 'amber', unit: 'records'},
  {id: 'allergies', label: 'Allergies', clinicalLabel: 'Allergies & intolerances', icon: 'fa-solid fa-triangle-exclamation', route: '/medical-history', resourceTypes: ['AllergyIntolerance'], count: 0, color: 'pink', unit: 'recorded'},
  {id: 'lab-results', label: 'Lab Results', clinicalLabel: 'Observations & diagnostic reports', icon: 'fa-solid fa-flask', route: '/labs', resourceTypes: ['Observation', 'DiagnosticReport'], count: 0, color: 'blue', unit: 'results'},
  {id: 'immunizations', label: 'Immunizations', clinicalLabel: 'Vaccinations', icon: 'fa-solid fa-syringe', route: '/medical-history', resourceTypes: ['Immunization'], count: 0, color: 'green', unit: 'on file'},
  {id: 'visits', label: 'Visits & Notes', clinicalLabel: 'Encounters', icon: 'fa-solid fa-notes-medical', route: '/medical-history', resourceTypes: ['Encounter'], count: 0, color: 'teal', unit: 'encounters'},
  {id: 'procedures', label: 'Procedures', clinicalLabel: 'Procedures & surgeries', icon: 'fa-solid fa-user-nurse', route: '/medical-history', resourceTypes: ['Procedure'], count: 0, color: 'purple', unit: 'procedures'},
  {id: 'documents', label: 'Documents', clinicalLabel: 'Clinical documents & files', icon: 'fa-solid fa-file-medical', route: '/medical-history', resourceTypes: ['DocumentReference', 'Media', 'Binary'], count: 0, color: 'gray', unit: 'documents'},
  {id: 'care-team', label: 'Care Team', clinicalLabel: 'Practitioners & organizations', icon: 'fa-solid fa-user-doctor', route: '/practitioners', resourceTypes: ['Practitioner', 'Organization', 'CareTeam'], count: 0, color: 'teal', unit: 'practitioners'},
]

@Component({
    selector: 'app-dashboard',
    templateUrl: './dashboard.component.html',
    styleUrls: ['./dashboard.component.scss'],
    standalone: false
})
export class DashboardComponent implements OnInit {
  loading = false

  lastUpdated: Date = null
  sourceCount = 0

  tiles: DashboardTile[] = []
  customizing = false
  hasCustomOrder = false
  hasCustomColors = false
  palette = TILE_PALETTE

  private userId: string = undefined

  constructor(
    private fastenApi: FastenApiService,
    private authService: AuthService,
    private dashboardPreferences: DashboardPreferencesService,
    private router: Router,
  ) { }

  ngOnInit() {
    this.loading = true

    // paint immediately with the default order; re-apply the saved order once
    // the user id arrives (preferences are scoped per user)
    this.tiles = DEFAULT_TILES.map((tile) => ({...tile}))

    this.authService.GetCurrentUser().then((claims) => {
      this.userId = claims?.sub
    }).catch(() => {
      this.userId = undefined
    }).finally(() => {
      this.tiles = this.applySavedOrder(this.tiles)
      this.applySavedColors(this.tiles)
      this.hasCustomOrder = this.dashboardPreferences.hasCustomTileOrder(this.userId)
      this.hasCustomColors = this.dashboardPreferences.hasCustomTileColors(this.userId)
    })

    this.fastenApi.getSummary().subscribe({
      next: (summary: Summary) => {
        this.loading = false
        this.populateTileCounts(summary)

        this.sourceCount = summary.sources?.length || 0
        if (summary.sources && summary.sources.length > 0) {
          this.lastUpdated = summary.sources.reduce((latest, source) => {
            const sourceDate = new Date(source.updated_at);
            return sourceDate > latest ? sourceDate : latest;
          }, new Date(0));
        }
      },
      error: () => { this.loading = false }
    })

    // The backend classifier synthesizes Condition.category + display state and separates real
    // health problems from social/administrative "Personal Health Conditions". The Medical Concerns
    // and Patient Profile tiles are counted from it (not the raw Condition count).
    this.fastenApi.getClassifiedConditions().subscribe({
      next: (rows: ClassifiedCondition[]) => {
        const problems = (rows || []).filter((r) => r.category === 'problem-list-item')
        // Concerns = Active + Remission (still tracked) + Unknown (shown, never assumed); RuledOut is
        // excluded, Resolved is past, and entered-in-error never reaches the client.
        const concerns = problems.filter((r) => r.state === 'Active' || r.state === 'Remission' || r.state === 'Unknown')
        const profile = (rows || []).filter((r) => r.category === 'sdoh' || r.category === 'health-concern')
        this.setTileCount('concerns', concerns.length)
        this.setTileCount('profile', profile.length)
      },
    })
  }

  private setTileCount(countKey: 'concerns' | 'profile', count: number) {
    const tile = this.tiles.find((t) => t.countKey === countKey)
    if (tile) tile.count = count
  }

  private populateTileCounts(summary: Summary) {
    const countsByType: Record<string, number> = {}
    for (const typeCount of (summary.resource_type_counts || [])) {
      countsByType[typeCount.resource_type] = typeCount.count
    }
    for (const tile of this.tiles) {
      if (tile.countKey) continue // counted from the classifier, not summary resource counts
      tile.count = tile.resourceTypes.reduce((sum, resourceType) => sum + (countsByType[resourceType] || 0), 0)
    }
  }

  private applySavedOrder(tiles: DashboardTile[]): DashboardTile[] {
    const savedOrder = this.dashboardPreferences.getTileOrder(this.userId)
    if (!savedOrder?.length) return tiles
    // saved ids first (in saved order), then any tiles added since the save
    const byId = new Map(tiles.map((tile) => [tile.id, tile]))
    const ordered = savedOrder.map((id) => byId.get(id)).filter((tile) => !!tile)
    const remaining = tiles.filter((tile) => !savedOrder.includes(tile.id))
    return [...ordered, ...remaining]
  }

  private applySavedColors(tiles: DashboardTile[]): void {
    const savedColors = this.dashboardPreferences.getTileColors(this.userId)
    for (const tile of tiles) {
      const color = savedColors[tile.id]
      if (color && this.palette.includes(color)) tile.color = color
    }
  }

  setTileColor(tile: DashboardTile, color: string) {
    tile.color = color
    this.dashboardPreferences.setTileColor(tile.id, color, this.userId)
    this.hasCustomColors = true
  }

  openTile(tile: DashboardTile) {
    if (this.customizing) return
    this.router.navigate([tile.route])
  }

  toggleCustomizing() {
    this.customizing = !this.customizing
  }

  onTileDrop(event: CdkDragDrop<DashboardTile[]>) {
    moveItemInArray(this.tiles, event.previousIndex, event.currentIndex)
    this.dashboardPreferences.setTileOrder(this.tiles.map((tile) => tile.id), this.userId)
    this.hasCustomOrder = true
  }

  resetTileOrder() {
    this.dashboardPreferences.resetTileOrder(this.userId)
    this.dashboardPreferences.resetTileColors(this.userId)
    this.hasCustomOrder = false
    this.hasCustomColors = false
    const countsById = new Map(this.tiles.map((tile) => [tile.id, tile.count]))
    this.tiles = DEFAULT_TILES.map((tile) => ({...tile, count: countsById.get(tile.id) ?? 0}))
  }
}

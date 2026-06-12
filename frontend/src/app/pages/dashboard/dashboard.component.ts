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
export interface DashboardTile {
  id: string
  label: string
  clinicalLabel: string
  icon: string
  route: string
  resourceTypes: string[]
  count: number
}

export const DEFAULT_TILES: DashboardTile[] = [
  {id: 'health-issues', label: 'Health Issues', clinicalLabel: 'Conditions & diagnoses', icon: 'fa-solid fa-heart-pulse', route: '/medical-history', resourceTypes: ['Condition'], count: 0},
  {id: 'medications', label: 'Medications', clinicalLabel: 'Prescriptions & medication statements', icon: 'fa-solid fa-pills', route: '/medications', resourceTypes: ['MedicationRequest', 'MedicationStatement', 'Medication', 'MedicationAdministration', 'MedicationDispense'], count: 0},
  {id: 'allergies', label: 'Allergies', clinicalLabel: 'Allergies & intolerances', icon: 'fa-solid fa-triangle-exclamation', route: '/medical-history', resourceTypes: ['AllergyIntolerance'], count: 0},
  {id: 'lab-results', label: 'Lab Results', clinicalLabel: 'Observations & diagnostic reports', icon: 'fa-solid fa-flask', route: '/labs', resourceTypes: ['Observation', 'DiagnosticReport'], count: 0},
  {id: 'immunizations', label: 'Immunizations', clinicalLabel: 'Vaccinations', icon: 'fa-solid fa-syringe', route: '/medical-history', resourceTypes: ['Immunization'], count: 0},
  {id: 'visits', label: 'Visits & Notes', clinicalLabel: 'Encounters', icon: 'fa-solid fa-notes-medical', route: '/medical-history', resourceTypes: ['Encounter'], count: 0},
  {id: 'procedures', label: 'Procedures', clinicalLabel: 'Procedures & surgeries', icon: 'fa-solid fa-user-nurse', route: '/medical-history', resourceTypes: ['Procedure'], count: 0},
  {id: 'documents', label: 'Documents', clinicalLabel: 'Clinical documents & files', icon: 'fa-solid fa-file-medical', route: '/medical-history', resourceTypes: ['DocumentReference', 'Media', 'Binary'], count: 0},
  {id: 'care-team', label: 'Care Team', clinicalLabel: 'Practitioners & organizations', icon: 'fa-solid fa-user-doctor', route: '/practitioners', resourceTypes: ['Practitioner', 'Organization', 'CareTeam'], count: 0},
]

@Component({
    selector: 'app-dashboard',
    templateUrl: './dashboard.component.html',
    styleUrls: ['./dashboard.component.scss'],
    standalone: false
})
export class DashboardComponent implements OnInit {
  loading = false
  concernsLoading = false

  lastUpdated: Date = null
  sourceCount = 0

  tiles: DashboardTile[] = []
  customizing = false
  hasCustomOrder = false

  // Conditions, classified by the backend (Layer 1) and grouped here for display (Layer 2).
  currentConcerns: ClassifiedCondition[] = []  // active / remission / unknown problem-list-items
  pastProblems: ClassifiedCondition[] = []     // resolved problem-list-items
  profileItems: ClassifiedCondition[] = []     // social / personal "Patient Profile" items
  showPastProblems = false

  private userId: string = undefined

  constructor(
    private fastenApi: FastenApiService,
    private authService: AuthService,
    private dashboardPreferences: DashboardPreferencesService,
    private router: Router,
  ) { }

  ngOnInit() {
    this.loading = true
    this.concernsLoading = true

    // paint immediately with the default order; re-apply the saved order once
    // the user id arrives (preferences are scoped per user)
    this.tiles = DEFAULT_TILES.map((tile) => ({...tile}))

    this.authService.GetCurrentUser().then((claims) => {
      this.userId = claims?.sub
    }).catch(() => {
      this.userId = undefined
    }).finally(() => {
      this.tiles = this.applySavedOrder(this.tiles)
      this.hasCustomOrder = this.dashboardPreferences.hasCustomTileOrder(this.userId)
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
    // health problems from social/administrative "Personal Health Conditions". Grouping here is the
    // source-agnostic display mapping (Layer 2).
    this.fastenApi.getClassifiedConditions().subscribe({
      next: (rows: ClassifiedCondition[]) => {
        this.concernsLoading = false
        const problems = (rows || []).filter((r) => r.category === 'problem-list-item')
        // Current = Active + Remission (still tracked) + Unknown (shown, never assumed); RuledOut is
        // intentionally excluded from both lists, and entered-in-error never reaches the client.
        this.currentConcerns = problems.filter((r) => r.state === 'Active' || r.state === 'Remission' || r.state === 'Unknown')
        this.pastProblems = problems.filter((r) => r.state === 'Resolved')
        this.profileItems = (rows || []).filter((r) => r.category === 'sdoh' || r.category === 'health-concern')
      },
      error: () => { this.concernsLoading = false }
    })
  }

  private populateTileCounts(summary: Summary) {
    const countsByType: Record<string, number> = {}
    for (const typeCount of (summary.resource_type_counts || [])) {
      countsByType[typeCount.resource_type] = typeCount.count
    }
    for (const tile of this.tiles) {
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

  openTile(tile: DashboardTile) {
    if (this.customizing) return
    this.router.navigate([tile.route])
  }

  openConcern(concern: ClassifiedCondition) {
    this.router.navigate(['/explore', concern.sourceId, 'resource', concern.sourceResourceId])
  }

  togglePastProblems() {
    this.showPastProblems = !this.showPastProblems
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
    this.hasCustomOrder = false
    const countsById = new Map(this.tiles.map((tile) => [tile.id, tile.count]))
    this.tiles = DEFAULT_TILES.map((tile) => ({...tile, count: countsById.get(tile.id) ?? 0}))
  }
}

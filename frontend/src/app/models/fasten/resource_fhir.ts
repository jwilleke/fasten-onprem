import {Provenance} from './provenance';

export class ResourceFhir {
  user_id?: string
  source_id = ""
  source_resource_type = ""
  source_resource_id = ""

  fhir_version = ""
  resource_raw: IResourceRaw
  related_resources?: ResourceFhir[] = []

  sort_title = ""
  sort_date: Date = null

  provenance?: Provenance  // "who said this" — resolved on the generic read path (#271)

  constructor(object?: any) {
    return Object.assign(this, object)
  }
}


//This is the "raw" Fhir resource
export interface IResourceRaw {
  resourceType: string
  id?: string
  meta?: IResourceMetaRaw
}
// This is the "raw" Fhir Bundle resource
export interface IResourceBundleRaw {
  resourceType: string
  id?: string
  entry: IResourceBundleEntryRaw[]
  total?: number
  link?: IResourceBundleLinkRaw[]
  meta?: IResourceMetaRaw
}

export interface IResourceBundleLinkRaw {
  id?: string
  relation: string
  url: string
}

export interface IResourceBundleEntryRaw {
  id?: string
  fullUrl?: string
  resource: IResourceRaw
}

export interface IResourceMetaRaw {
  id?: string
  versionId?: string
  lastUpdated: string
}

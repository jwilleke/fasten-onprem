import {ResourceFhir} from './resource_fhir';

export class ResourceGraphResponse {
  results: Record<string, ResourceFhir[]>
}

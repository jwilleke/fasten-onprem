import {FastenOptions} from './fasten-options';
import {Provenance} from './provenance';
import {ResourceType} from '../constants';

export class FastenDisplayModel {
  source_resource_type: ResourceType | undefined
  source_resource_id: string | undefined
  source_id: string | undefined
  sort_title: string | undefined
  sort_date: Date | undefined

  // "Who said this" — resolved at read time on the generic resource path (#271). Undefined for models
  // not built from that path (e.g. storybook fixtures). Rendered once by the fhir-card host.
  provenance: Provenance | undefined

  related_resources: Record<string, FastenDisplayModel[]> = {}

  constructor(options?: FastenOptions) {}
}

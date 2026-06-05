import {PatientAccessBrand, PatientAccessPortal} from '../patient-access-brands';

export interface ConnectGatewayEndpointListDisplayItem {
  id: string;
  platform_type: string;
}

export interface ConnectGatewayPortalListDisplayItem extends PatientAccessPortal {
  endpoints?: ConnectGatewayEndpointListDisplayItem[];
}

export interface ConnectGatewayBrandListDisplayItem extends PatientAccessBrand {
  portals: ConnectGatewayPortalListDisplayItem[];
  hidden: boolean;
}



export class ConnectGatewaySourceSearchResult {
  _index: string;
  _type: string;
  _id: string;
  _score: number;
  _source: ConnectGatewayBrandListDisplayItem;
  sort: string[];
  highlight?: {
    aliases?: string[];
  }
}

export class ConnectGatewaySourceSearchAggregation {
  sum_other_doc_count: number;
  buckets: {
    key: string;
    doc_count: number;
  }[]
}

export class ConnectGatewaySourceSearch {
  _scroll_id: string;
  took: number;
  timed_out: boolean;
  hits: {
    total: {
      value: number;
      relation: string;
    };
    max_score: number;
    hits: ConnectGatewaySourceSearchResult[];
  };
  aggregations: {
    by_platform_type: ConnectGatewaySourceSearchAggregation
    by_category: ConnectGatewaySourceSearchAggregation
  };
}

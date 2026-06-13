// Mirrors the backend handler.resourceListItem returned by GET /api/secure/resources/recent and
// /api/secure/resources/search — a lightweight resource row (title + date + identity), no raw body.

export interface ResourceListItem {
  source_id: string;
  source_resource_type: string;
  source_resource_id: string;
  title: string;
  date?: string;
}

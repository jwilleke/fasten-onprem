import {DashboardWidgetQuery} from './dashboard-widget-query';
import * as _ from 'lodash';

export class DashboardWidgetConfig {
  id?: string
  item_type: "image-list-group-widget" | "complex-line-widget" | "donut-chart-widget" | "dual-gauges-widget" | "grouped-bar-chart-widget" | "patient-vitals-widget" | "simple-line-chart-widget" | "table-widget" | "records-summary-widget" | "medications-widget" | "profile-summary-widget"

  title_text: string
  description_text: string

  queries:  {
    q: DashboardWidgetQuery,
    conditional_formats?: [],
    dataset_options?: {
      label?: string,
      borderWidth?: number,
      borderColor?: string,
      fill?: boolean,
      backgroundColor?: string,
    }
    // type?: "line",
    // style?: {
    //   "palette": "grey" | "pastel" | "light" | "default"
    // }
  }[]


  //used for display purposes within the Dashboard, not for the actual chart
  // minWidth?: number
  // minHeight?: number
  width: number
  height: number
  x?: number
  y?: number

  // `view_all_route` is consumed by the profile-summary-widget (#245): when set, the widget renders a
  // "View all →" header link to that route. Optional — category-level list pages are a later phase.
  parsing?: {label?: string, xAxisKey?: string, yAxisKey?: string, view_all_route?: string} | Record<string, string>
}

export enum ToastType {
  Error = "error",
  Success = "success",
  Info = "info"
}

export class ToastNotification {
  event_date: Date = new Date()
  title?: string
  message: string
  type: ToastType = ToastType.Info
  displayClass = 'demo-static-toast'
  autohide = true
  link?: {
    text: string,
    url: string
  }
}

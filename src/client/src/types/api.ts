export interface CustomApiError {
  title?: string
  status: number
  detail?: string
  clientCode?: ErrorCodes
  extensions?: Record<string, string>
  errorMessage?: string
}
export type ErrorCodes = '400' | '500' | '501' | '401' | '403'

export type ApiResponse<T> = {
  data?: T
  error: boolean
  status: number
} & CustomApiError

export type CreateEntityResponse = {
  id: number
}

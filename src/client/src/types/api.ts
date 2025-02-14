export type HackerApiResponse<T> = {
  statusCode: number
  status: string
  message: string
  data: T
}

export type HackerApiResponse<T> = {
  statusCode: number
  status: string
  message: string
  data: DataBody<T>
}

export type DataBody<T> = {
  body: T
}

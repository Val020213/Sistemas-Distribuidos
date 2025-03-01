export type HackerApiResponse<T> = {
  statusCode: number
  status: string
  message: string
  data: T
}

export type GoogleApiSearchResponse = {
  spelling: {
    correctedQuery: string
    htmlCorrectedQuery: string
  }
  items: GoogleSearchResults[]
}

export type GoogleSearchResults = {
  title: string
  htmlTitle: string
  link: string
  snippet: string
}

export type UrlDataType = {
  key: number
  url: string
  status: 'complete' | 'in_progress' | 'error'
}

export type UrlStatus = UrlDataType['status']

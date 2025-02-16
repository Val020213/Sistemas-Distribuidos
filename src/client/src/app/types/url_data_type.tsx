export type UrlDataType = {
  url: string
  status: 'complete' | 'in_progress' | 'error'
  content?: string
}

export type UrlStatus = UrlDataType['status']

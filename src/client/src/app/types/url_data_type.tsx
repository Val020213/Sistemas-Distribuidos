export type UrlDataType = {
  id: number | string
  url: string
  status: 'complete' | 'in_progress' | 'error'
  content?: string
}

export type UrlStatus = UrlDataType['status']

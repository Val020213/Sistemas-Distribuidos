export type UrlDataType = {
  id: number | string
  url: string
  status: 'pending' | 'scrapped' | 'error'
}

export type UrlStatus = UrlDataType['status']

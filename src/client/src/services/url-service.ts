'use server'
import { UrlDataType } from '@/app/types/url_data_type'
import { backendRoutes, tagsRoutes } from '@/routes/routes'
import { HackerApiResponse } from '@/types/api'

export async function fetchUrlService(
  url: string
): Promise<HackerApiResponse<string>> {
  const response = await fetch(`${backendRoutes.fetch}`, {
    method: 'POST',
    body: JSON.stringify({ url: url }),
    next: {
      tags: [tagsRoutes.fetch],
    },
  })
  return response.json()
}

export async function listUrlService(): Promise<
  HackerApiResponse<UrlDataType[]>
> {
  console.log('listUrlService - Request to address:', backendRoutes.list)
  const response = await fetch(`${backendRoutes.list}`, {
    method: 'GET',
    next: {
      tags: [tagsRoutes.list],
    },
  })
  console.log('listUrlService - Response:', response)
  return response.json()
}

export async function downloadUrlService(
  url: string
): Promise<HackerApiResponse<UrlDataType>> {
  const response = await fetch(backendRoutes.download, {
    method: 'POST',
    body: JSON.stringify({ url: url }),
    next: {
      tags: [tagsRoutes.download],
    },
  })
  return response.json()
}

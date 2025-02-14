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
  const response = await fetch(`${backendRoutes.list}`, {
    method: 'GET',
    next: {
      tags: [tagsRoutes.list],
    },
  })

  return response.json()
}

export async function downloadUrlService(
  id: string
): Promise<HackerApiResponse<UrlDataType>> {
  const response = await fetch(`${backendRoutes.download.replace(':id', id)}`, {
    method: 'GET',
    next: {
      tags: [tagsRoutes.download],
    },
  })
  return response.json()
}

'use server'
import { UrlDataType } from '@/app/types/url_data_type'
import { backendRoutes, tagsRoutes } from '@/routes/routes'
import { HackerApiResponse } from '@/types/api'

export async function fetchUrlService(
  url: string
): Promise<HackerApiResponse<string>> {
  const response = await fetch(`${backendRoutes.fetch}?url=${url}`, {
    method: 'GET',
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
      revalidate: 1,
    },
  })
  return response.json()
}
export async function downloadUrlService(id: string): Promise<string> {
  const response = await fetch(`${backendRoutes.download}?id=${id}`, {
    method: 'GET',
    next: {
      tags: [tagsRoutes.download],
    },
  })
  return response.text()
}

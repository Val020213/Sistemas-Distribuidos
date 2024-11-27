'use server'
import { backendRoutes } from '@/routes/routes'
import { HackerApiResponse } from '@/types/api'

export async function fetchUrlService(
  url: string
): Promise<HackerApiResponse<string>> {
  const response = await fetch(`${backendRoutes.fetch}?url=${url}`, {
    method: 'GET',
    next: {
      tags: ['fetchUrlService'],
    },
  })
  return response.json()
}

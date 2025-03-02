'use server'
import { UrlDataType } from '@/app/types/url_data_type'
import { backendRoutes, tagsRoutes } from '@/routes/routes'
import { GoogleApiSearchResponse, HackerApiResponse } from '@/types/api'

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
  console.log('AGUA', response)
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
  console.log('AGUA', response)
  return response.json()
}

export async function webSearchApi(
  search: string
): Promise<GoogleApiSearchResponse> {
  const apiKey = process.env.NEXT_GOOGLE_API_KEY

  if (!apiKey) {
    throw new Error('Missing Google API Key')
  }
  
  const response = await fetch(
    `${backendRoutes.search}?key=${apiKey}&cx=a5e01a1aa147a40fc&q=${encodeURIComponent(search)}`
  )
  console.log(response)
  if (!response.ok) {
    throw new Error(`Error: ${response.statusText}`)
  }
  return response.json()
}

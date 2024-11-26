import { backendRoutes } from '@/routes/routes'

export const fetchUrlService = async (url: string) => {
  const response = await fetch(`${backendRoutes.fetch}?url=${url}`, {
    method: 'GET',
    mode: 'no-cors',
    next: {
      tags: ['fetchUrlService'],
    },
  })
  return response
}

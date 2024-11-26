import { authenticationSchemaType } from '@/modules/system/security/authentication/authenticationSchema'
import { buildApiResponseAsync, handleApiError } from './api'
import { getToken } from './cookies'

export async function signing(data: authenticationSchemaType) {
  const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/v1/auth/token`, {
    cache: 'no-store',
    headers: {
      'Content-Type': 'application/json',
    },
    method: 'POST',
    body: JSON.stringify(data),
  })

  if (!res.ok) await handleApiError(res)
  return res
}

export async function signingOut() {
  const token = await getToken()

  const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/v1/user/logout`, {
    method: 'POST',
    headers: {
      Authorization: `Bearer ${token}`,
    },
  })

  if (!res.ok) await handleApiError(res)
  return res
}

const fetchWithAuth = async (url: string, options: RequestInit) => {
  const res = await fetch(url, {
    ...options,
    headers: {
      ...options.headers,
      Authorization: `Bearer ${await getToken()}`,
    },
  })
  if (!res.ok) await handleApiError(res)
  return buildApiResponseAsync(res.json())
}

export default fetchWithAuth

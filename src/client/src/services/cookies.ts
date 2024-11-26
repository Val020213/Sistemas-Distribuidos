'use server'

import { cookies, headers } from 'next/headers'
import { getAuthCookieName } from './api'

export const getToken = async () => {
  try {
    const cookieStore = cookies()
    const token = (await cookieStore).get(getAuthCookieName()!)?.value ?? ''
    return token
  } catch (e) {
    console.log('error getting token', e)
    throw new Error('Error getting token')
  }
}

export const getTokenHeader = async () => {
  const headersList = headers()
  try {
    return { Cookie: (await headersList).get('Cookie') ?? '' }
  } catch (e) {
    console.log('error getting token', e)
    return { Cookie: '' }
  }
}

export const deleteAccessToken = async () => {
  try {
    const cookieStore = cookies()
    ;(await cookieStore).delete(getAuthCookieName()!)
  } catch (e) {
    console.log('error deleting token', e)
    return ''
  }
}

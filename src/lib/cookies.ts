'use server'
import { cookies } from 'next/headers'

export async function saveUserData(key: string, value: string) {
  const cookieStore = await cookies()
  cookieStore.set(key, value, { path: '/', maxAge: 60 * 60 * 24 * 7 }) // 1 week
}

export async function getUserData(key: string): Promise<string | undefined> {
  const cookieStore = await cookies()
  return cookieStore.get(key)?.value
}

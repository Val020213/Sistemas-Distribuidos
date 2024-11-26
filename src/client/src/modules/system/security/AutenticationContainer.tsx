'use client'

import { GlobalContext } from '@/core/global-context/GlobalContext'
import { ReactNode, useContext, useEffect } from 'react'
import AuthenticationModal from './modal/AutenticationModal'
import { getUserData } from '@/lib/cookies'
import useShowHackerMessage from '@/hooks/useShowHackerMessage'

const HackerProtect = ({ children }: { children: ReactNode }) => {
  const globalContext = useContext(GlobalContext)
  const hackerMessage = useShowHackerMessage()

  async function getUserFromCookies() {
    const name = await getUserData('name')
    const token = await getUserData('token')
    if (name && token) {
      globalContext.user[1]({
        name,
        token,
      })
    }
  }

  useEffect(() => {
    getUserFromCookies()
  }, [])

  useEffect(() => {
    if (!globalContext.user[0]?.token.length) return

    let timeout: NodeJS.Timeout
    let warningTimeout: NodeJS.Timeout

    const resetTimeout = () => {
      clearTimeout(timeout)
      clearTimeout(warningTimeout)

      warningTimeout = setTimeout(() => {
        hackerMessage('Your session will expire in 3 minutes.')
      }, 7 * 60 * 1000) // 7 minutes to show the warning message

      timeout = setTimeout(() => {
        globalContext.user[1](undefined)
      }, 10 * 60 * 1000) // 10 minutes to log out
    }

    const addEventListeners = () => {
      window.addEventListener('mousemove', resetTimeout)
      window.addEventListener('keydown', resetTimeout)
    }

    const removeEventListeners = () => {
      window.removeEventListener('mousemove', resetTimeout)
      window.removeEventListener('keydown', resetTimeout)
    }

    resetTimeout()
    addEventListeners()

    return () => {
      clearTimeout(timeout)
      clearTimeout(warningTimeout)
      removeEventListeners()
    }
  }, [globalContext.user[0]])

  return (
    <>
      {!globalContext.user[0]?.token.length && <AuthenticationModal />}
      {children}
    </>
  )
}

export default HackerProtect

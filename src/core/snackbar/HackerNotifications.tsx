'use client'
import { ReactNode, createContext, useContext, useMemo, useState } from 'react'
import Alert from './HackerAlerts'

export type HackerNotificationsType = 'info' | 'success' | 'error' | 'warn'
type ReducedHackerNotificationsBuilder = (
  message: string,
  _timeout?: number
) => void
type HackerNotificationsBuilder = (
  message: string,
  type?: HackerNotificationsType,
  _timeout?: number
) => void

interface HackerNotificationsContextType {
  showHackerNotifications: HackerNotificationsBuilder
  success: ReducedHackerNotificationsBuilder
  error: ReducedHackerNotificationsBuilder
  warn: ReducedHackerNotificationsBuilder
  info: ReducedHackerNotificationsBuilder
}

const HackerNotificationsContext =
  createContext<HackerNotificationsContextType>(
    {} as HackerNotificationsContextType
  )

export const HackerNotificationsTimeout = 999999

const HackerNotificationsProvider: React.FC<{ children: ReactNode }> = ({
  children,
}) => {
  const [open, setOpen] = useState(false),
    [type, setType] = useState<HackerNotificationsType>('info'),
    [message, setMessage] = useState(''),
    [timeout, setTimeout] = useState(HackerNotificationsTimeout),
    showHackerNotifications = (
      message: string,
      type = 'info' as HackerNotificationsType,
      _timeout = HackerNotificationsTimeout
    ) => {
      setType(type)
      setMessage(message)
      setOpen(true)
      setTimeout(_timeout)
    },
    success = (message: string, _timeout = timeout) =>
      showHackerNotifications(message, 'success', _timeout),
    error = (message: string, _timeout = timeout) =>
      showHackerNotifications(message, 'error', _timeout),
    warn = (message: string, _timeout = timeout) =>
      showHackerNotifications(message, 'warn', _timeout),
    info = (message: string, _timeout = timeout) =>
      showHackerNotifications(message, 'info', _timeout)

  const providerProps = useMemo(
    () => ({
      showHackerNotifications,
      success,
      error,
      warn,
      info,
    }),
    []
  )

  function handleClose() {
    setOpen(false)
  }

  return (
    <HackerNotificationsContext.Provider value={providerProps}>
      {children}
      {open && <Alert message={message} type={type} onClose={handleClose} />}
    </HackerNotificationsContext.Provider>
  )
}

export const useHackerNotifications = () =>
  useContext(HackerNotificationsContext)
export default HackerNotificationsProvider

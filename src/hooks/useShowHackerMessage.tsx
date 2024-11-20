import {
  HackerNotificationsTimeout,
  useHackerNotifications,
} from '@/core/snackbar/HackerNotifications'

const useShowHackerMessage = () => {
  const hackerNotifications = useHackerNotifications()
  const showHackerMessage = (
    message: string,
    type: 'info' | 'success' | 'error' | 'warn' = 'error',
    timeout = HackerNotificationsTimeout
  ) => {
    hackerNotifications.showHackerNotifications(message, type, timeout)
  }
  return showHackerMessage
}

export default useShowHackerMessage

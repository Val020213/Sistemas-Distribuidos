import { IconButton, Snackbar, Stack, Typography } from '@mui/material'
import React, { useEffect, useState } from 'react'
import {
  HackerNotificationsTimeout,
  HackerNotificationsType,
} from './HackerNotifications'
import { tailwindColors } from '@/theme/tailwindColors'

type Props = {
  message: string
  type: HackerNotificationsType
  autoHideDuration?: number
  onClose?: () => void
}

const getTypeColor = (type: HackerNotificationsType) => {
  switch (type) {
    case 'info':
      return 'blue'
    case 'warn':
      return 'yellow'
    case 'error':
      return 'red'
    case 'success':
      return 'green'
    default:
      return 'green'
  }
}

const HackerAlert = ({
  message,
  type,
  onClose,
  autoHideDuration = HackerNotificationsTimeout,
}: Props) => {
  const [displayedMessage, setDisplayedMessage] = useState<string>('')

  useEffect(() => {
    let i = 0
    const typingInterval = setInterval(() => {
      if (i < message.length) {
        setDisplayedMessage((prev) => prev + message[prev.length])
        i++
      } else {
        clearInterval(typingInterval)
      }
    }, 50)

    return () => clearInterval(typingInterval)
  }, [message])

  return (
    <Snackbar
      open={true}
      autoHideDuration={autoHideDuration}
      onClose={onClose}
      anchorOrigin={{ vertical: 'bottom', horizontal: 'center' }}
    >
      <Stack
        display={'flex'}
        alignItems={'flex-start'}
        justifyContent={'space-between'}
        sx={{
          animation: `slideIn 0.3s ease-out`,
          fontFamily: '"Courier New", Courier, monospace',
          backgroundColor: '#000',
          color: tailwindColors[getTypeColor(type)][500],
          padding: '16px',
          borderRadius: '4px',
          width: { xs: '100vw', sm: '80vw', md: '50vw', lg: '40vw' },
          boxShadow: `0 0 10px ${tailwindColors[getTypeColor(type)][500]}`,
          display: 'flex',
          justifyContent: 'space-between',
        }}
      >
        <Typography
          whiteSpace={'break-spaces'}
          width={'100%'}
          fontSize={'0.825rem'}
        >
          {type === 'success' && '> SUCCESS: '}
          {type === 'error' && '> ERROR: '}
          {type === 'warn' && '> WARNING: '}
          {type === 'info' && '> INFO: '}
          {displayedMessage}
          <Typography
            component={'span'}
            style={{
              display: 'inline-block',
              width: '8px',
              height: '14px',
              backgroundColor: 'currentColor',
              marginLeft: '4px',
              animation: `blink 0.7s infinite`,
            }}
          />
        </Typography>
        <IconButton
          onClick={onClose}
          disableRipple
          sx={{
            color: 'currentcolor',
            fontSize: '0.875rem',
            p: 0,
            m: 0,
          }}
        >
          x
        </IconButton>
      </Stack>
    </Snackbar>
  )
}

export default HackerAlert

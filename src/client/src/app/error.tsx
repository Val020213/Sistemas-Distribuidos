'use client'
import { HackerButton } from '@/core/button/HackerButton'
import { tailwindColors } from '@/theme/tailwindColors'
import { Box, Stack, Typography } from '@mui/material'
import { useEffect, useMemo, useState } from 'react'

export default function Error({
  error,
  reset,
}: {
  error: Error & { digest?: string }
  reset: () => void
}) {
  const message = useMemo(() => {
    return (
      '> ERROR_CRÍTICO: M4lfuncionamiento_del_Sistema.\n> 1ntentando restablecer la con3xión...\n> Todos los sistemas están fuera de línea.\n> Iniciar anulación manual.\n> ' +
      'Excepción Encontrada: ' +
      error.message
    )
  }, [error.message])

  const [displayText, setDisplayText] = useState<string>('')

  useEffect(() => {
    let i = 0
    const typingEffect = setInterval(() => {
      if (i < message.length) {
        setDisplayText((prevText) => prevText + message[prevText.length])
        i++
      } else {
        clearInterval(typingEffect)
      }
    }, 40)

    return () => clearInterval(typingEffect)
  }, [message])

  return (
    <Stack justifyContent={'center'} alignItems={'center'}>
      <Stack maxWidth={'md'} spacing={4}>
        <Typography variant={'h3'} fontWeight={'bold'}>
          SYSTEM FAILURE
        </Typography>
        <Stack>
          <Box>
            <Typography
              variant="body1"
              component="pre"
              className="animate-flicker"
              style={{
                whiteSpace: 'pre-wrap',
                fontFamily: 'monospace',
                color: tailwindColors['green'][500],
                marginRight: '64px',
              }}
            >
              {`
           ________       ___    ___ ________  _________  _______  ______ _____      
          |\\   ____\\     |\\  \\  /  /|\\   ____\\/___    ___\\\  _____ \\ \\   _    _  \\    
          \\ \\  \\___|_    \\ \\  \\/  / | \\  \\___\|/__  \\  \\_/ \\  \\____ \\ \\  \\\\\\__\\ \\  \\   
           \\ \\_____  \\    \\ \\    / / \\ \\_____  \\  \\ \\  \\ \\ \\  _____ \\\ \\  \\\\|__| \\  \\  
            \\|____|\\  \\    \\/  /  /   \\|____|\\  \\  \\ \\  \\ \\ \\  \\____ \\ \\  \\    \\ \\  \\ 
              ____\\_\\  \\ __/  / /       ____\\_\\  \\  \\ \\__\\ \\ \\_______ \\ \\__\\    \\ \\__\\
            |\\_________\\\\___/ /       |\\_________\\   \\|__|  \\|_______| \\|__|     \\|__|
            \\|_________\\|___|/        \\|_________|                                   
        
            D O W N
          `}
            </Typography>
          </Box>

          <Box
            className="animate-glitch"
            sx={{
              background: tailwindColors['black'],
              height: 160,
              overflowX: 'hidden',
              border: '1px solid',
              borderColor: tailwindColors['green'][500],
              position: 'relative',
              padding: '8px',
              borderRadius: '2px',
            }}
          >
            <Stack sx={{ overflowY: 'auto' }}>
              <Typography
                variant="body1"
                component="pre"
                style={{ fontFamily: 'monospace', position: 'relative' }}
              >
                {displayText}{' '}
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
            </Stack>
          </Box>
        </Stack>

        <HackerButton variant="Button" onClick={() => reset()} fullWidth>
          REINTENTAR
        </HackerButton>
      </Stack>
    </Stack>
  )
}

'use client'

import { useState, useEffect } from 'react'
import { Button, Typography } from '@mui/material'
import { styled, keyframes, Stack } from '@mui/system'
import { tailwindColors } from '@/theme/tailwindColors'

const gradientMove = keyframes`
  0% { background-position: 0% 50%; }
  50% { background-position: 100% 50%; }
  100% { background-position: 0% 50%; }
`

const textColorChange = keyframes`
  0%, 100% { color: #33ff33;  }
  20% { color: #ff9933; }
  40% { color: #3399ff; }
  60% { color: #cc33ff; }
  80% { color: #ff3333; }
`

const iconColorChange = keyframes`
  0%, 100% { color: #33ff33; filter: drop-shadow(0 0 3px #33ff33); }
  20% { color: #ff9933; filter: drop-shadow(0 0 3px #ff9933); }
  40% { color: #3399ff; filter: drop-shadow(0 0 3px #3399ff); }
  60% { color: #cc33ff; filter: drop-shadow(0 0 3px #cc33ff); }
  80% { color: #ff3333; filter: drop-shadow(0 0 3px #ff3333); }
`

const glitch = keyframes`
  0% { transform: translate(0); }
  20% { transform: translate(-2px, 2px); }
  40% { transform: translate(-2px, -2px); }
  60% { transform: translate(2px, 2px); }
  80% { transform: translate(2px, -2px); }
  100% { transform: translate(0); }
`

const RetroButton = styled(Button)(({ theme }) => ({
  position: 'relative',
  overflow: 'hidden',
  fontFamily: 'monospace',
  px: 2,
  gap: 2,
  borderRadius: 0.5,
  border: `1px solid ${tailwindColors.green[400]}`,
  background: 'rgba(10, 10, 20, 0.8)',
  transition: 'all 0.3s',
  boxShadow: '0 0 5px rgba(255, 255, 255, 0.1)',
  '&::before': {
    content: '""',
    position: 'absolute',
    top: 0,
    left: 0,
    right: 0,
    bottom: 0,
    background:
      'linear-gradient(-45deg, #33ff33, #ff9933, #3399ff, #cc33ff, #ff3333)',
    backgroundSize: '400% 400%',
    animation: `${gradientMove} 8s ease infinite`,
    opacity: 0.3,
    zIndex: 0,
  },
  '&::after': {
    content: '""',
    position: 'absolute',
    top: 0,
    left: 0,
    right: 0,
    bottom: 0,
    background:
      'repeating-linear-gradient(transparent, transparent 2px, rgba(255, 255, 255, 0.1) 2px, rgba(255, 255, 255, 0.1) 4px)',
    opacity: 0.1,
    pointerEvents: 'none',
    zIndex: 1,
  },
}))

const AnimatedIcon = styled('div')({
  animation: `${iconColorChange} 8s linear infinite`,
})

const AnimatedText = styled(Typography)({
  animation: `${textColorChange} 8s linear infinite`,
  fontWeight: 'bold',
  letterSpacing: '0.05em',
})

const ButtonLine = styled('div')({
  position: 'absolute',
  bottom: 0,
  left: 0,
  width: '100%',
  height: '4px',
  background:
    'linear-gradient(90deg, #33ff33, #ff9933, #3399ff, #cc33ff, #ff3333)',
  backgroundSize: '400% 400%',
  transform: 'scaleX(0)',
  transformOrigin: 'left',
  transition: 'transform 0.3s ease-out',
})

type Props = {
  icons: React.ReactNode[]
  text: string
  onClick: () => void
}

export default function MuiRetroHackerButton({
  icons,
  text,
  onClick,
}: Readonly<Props>) {
  const [isHovered, setIsHovered] = useState(false)
  const [currentIcon, setCurrentIcon] = useState(0)

  useEffect(() => {
    if (!isHovered) {
      const interval = setInterval(() => {
        setCurrentIcon((prev) => (prev + 1) % icons.length)
      }, 3000)
      return () => clearInterval(interval)
    }
  }, [isHovered])

  return (
    <RetroButton
      fullWidth
      onMouseEnter={() => setIsHovered(true)}
      onMouseLeave={() => setIsHovered(false)}
      onClick={() => onClick()}
    >
      <AnimatedIcon>{icons[currentIcon]}</AnimatedIcon>
      <AnimatedText variant="button" ml={1}>
        {text}
      </AnimatedText>

      <ButtonLine
        sx={{
          transform: isHovered ? 'scaleX(1)' : 'scaleX(0)',
          animation: isHovered ? `${gradientMove} 3s ease infinite` : 'none',
        }}
      />
    </RetroButton>
  )
}

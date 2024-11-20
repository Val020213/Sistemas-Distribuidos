import { tailwindAnimations } from '@/theme/animations'
import { tailwindColors } from '@/theme/tailwindColors'
import { Theme } from '@emotion/react'
import { Box, Button, ButtonProps, SxProps, Typography } from '@mui/material'
import Image from 'next/image'
import React from 'react'
import disableIcon from '@/assets/jinx icon.svg'

type Props = Omit<Omit<ButtonProps, 'variant'>, 'color'> & {
  variant?: 'Card' | 'Button'
  color?: 'blue' | 'green' | 'red' | 'yellow' | 'purple' | 'orange'
  icon?: React.ReactNode
  children: React.ReactNode
}

export const HackerButton: React.FC<Props> = ({
  variant = 'Card',
  color = 'green',
  icon,
  children,
  ...others
}) => {
  return (
    <Button
      variant="outlined"
      className="group"
      {...others}
      sx={{
        height: variant === 'Card' ? '6rem' : '48px',
        color: `${tailwindColors[color][400]}`,
        borderColor: `${tailwindColors[color][400]}`,
        backgroundColor: 'black',
        '&:hover': {
          backgroundColor: `${tailwindColors[color][900]}`,
          color: `${tailwindColors[color][300]}`,
        },
        transition: 'all 0.3s',
        overflow: 'hidden',
        position: 'relative',
        borderRadius: variant === 'Button' ? 0.5 : 'inherit',
        px: variant === 'Card' ? 4 : 2,
        display: 'flex',
        flexDirection: variant === 'Card' ? 'column' : 'row',
        justifyContent: 'center',
        alignItems: 'center',
        gap: 2,
        '&.Mui-disabled': {
          backgroundColor: 'gray',
          color: 'lightgray',
          borderColor: 'lightgray',
        },
        ...(others.sx as SxProps<Theme>),
      }}
    >
      <Box
        sx={{
          position: 'absolute',
          inset: 0,
          backgroundColor: `${tailwindColors[color][400]}1A`,
          opacity: 0.1,
          animation: tailwindAnimations.pulse,
        }}
      />

      {others.disabled ? (
        <>
          <Box
            sx={{
              position: 'absolute',
              inset: 0,
              justifyContent: 'center',
              alignItems: 'center',
              display: 'flex',
              color: 'gray',
            }}
          >
            <Image
              src={disableIcon}
              alt="disable icon"
              width={48}
              height={48}
              style={{ filter: 'grayscale(1)' }}
            />
          </Box>
        </>
      ) : (
        <>
          {icon}
          <Typography fontFamily={'monospace'}>{children}</Typography>
          <Box
            sx={{
              position: 'absolute',
              bottom: 0,
              left: 0,
              width: '100%',
              height: '4px',
              backgroundColor: `${tailwindColors[color][400]}`,
              transform: 'scaleX(0)',
              transition: 'transform 0.3s',
              '.group:hover &': {
                transform: 'scaleX(1)',
              },
            }}
          />
        </>
      )}
    </Button>
  )
}

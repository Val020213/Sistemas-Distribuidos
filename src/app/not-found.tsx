'use client'
import { HackerButton } from '@/core/button/HackerButton'
import { tailwindColors } from '@/theme/tailwindColors'
import { Box, Stack, Typography } from '@mui/material'
import { useRouter } from 'next/navigation'
import image from '@/assets/Frame 1.png'
import Image from 'next/image'

export default function NotFound() {
  const navigate = useRouter()
  return (
    <Stack justifyContent={'center'} alignItems={'center'} height={'100%'}>
      <Stack maxWidth={'md'} justifyContent={'space-between'} height={'100%'}>
        <Stack spacing={4} position={'relative'}>
          <Typography variant={'h3'} fontWeight={'bold'}>
            Sistema de Redirección
          </Typography>
          <Box>
            <Typography
              variant="body1"
              className="animate-flicker"
              style={{
                whiteSpace: 'pre-wrap',
                fontFamily: 'monospace',
                color: tailwindColors['green'][500],
                marginRight: '64px',
              }}
            >
              {String.raw`
 ___   ___  ________  ___   ___                                                             
|\  \ |\  \|\   __  \|\  \ |\  \                                                            
\ \  \\_\  \ \  \|\  \ \  \\_\  \                                                           
 \ \______  \ \  \\\  \ \______  \                                                          
  \|_____|\  \ \  \\\  \|_____|\  \                                                         
         \ \__\ \_______\     \ \__\                                                        
          \|__|\|_______|      \|__|                                                        
                                                                                   
 ________   ________  _________        ________ ________  ___  ___  ________   ________     
|\   ___  \|\   __  \|\___   ___\     |\  _____\\   __  \|\  \|\  \|\   ___  \|\   ___ \    
\ \  \\ \  \ \  \|\  \|___ \  \_|     \ \  \__/\ \  \|\  \ \  \\\  \ \  \\ \  \ \  \_|\ \   
 \ \  \\ \  \ \  \\\  \   \ \  \       \ \   __\\ \  \\\  \ \  \\\  \ \  \\ \  \ \  \ \\ \  
  \ \  \\ \  \ \  \\\  \   \ \  \       \ \  \_| \ \  \\\  \ \  \\\  \ \  \\ \  \ \  \_\\ \ 
   \ \__\\ \__\ \_______\   \ \__\       \ \__\   \ \_______\ \_______\ \__\\ \__\ \_______\
    \|__| \|__|\|_______|    \|__|        \|__|    \|_______|\|_______|\|__| \|__|\|_______|
    
  `}
            </Typography>
            <Stack>
              <Typography textAlign={'end'}>
                404 | Página no encontrada
              </Typography>
            </Stack>
          </Box>
          <Box
            sx={{
              position: 'absolute',
              top: 0,
              right: 0,
            }}
          >
            <Image src={image} width={300} height={300} alt="" />
          </Box>
        </Stack>
        <HackerButton
          variant="Button"
          onClick={() => navigate.back()}
          fullWidth
          sx={{
            mb: 4,
          }}
        >
          Regresar
        </HackerButton>
      </Stack>
    </Stack>
  )
}

import { HackerButton } from '@/core/button/HackerButton'
import { clientRoutes } from '@/routes/routes'
import { Grid2 as Grid, List, ListItem, Stack, Typography } from '@mui/material'
import {
  ChartNoAxesCombined,
  Play,
  Power,
  Rocket,
  Settings,
} from 'lucide-react'
import Link from 'next/link'

const SystemContainer = () => {
  return (
    <Stack maxWidth={'md'} spacing={4} height={'100%'}>
      <Typography variant="h3" className="animate-glitch">
        SISTEMA DE SCRAPPING V1.0
        <Typography
          component={'span'}
          style={{
            fontSize: 'inherit',
            display: 'inline-block',
            marginLeft: '4px',
            animation: `blink 0.7s infinite`,
          }}
        >
          {'_'}
        </Typography>
      </Typography>
      <Stack height={'100%'} justifyContent={'center'}>
        <Grid container spacing={{ xs: 2, md: 3 }} columns={{ xs: 1, sm: 3 }}>
          <Grid size={{ xs: 1, sm: 1.5 }}>
            <Link href={clientRoutes.service.scrapper}>
              <HackerButton icon={<Play size={24} />} fullWidth>
                Iniciar Scrapper
              </HackerButton>
            </Link>
          </Grid>
          <Grid size={{ xs: 1, sm: 1.5 }}>
            <Link href={clientRoutes.service.dashboard}>
              <HackerButton
                color="orange"
                icon={<ChartNoAxesCombined size={24} />}
                fullWidth
              >
                Panel de Control
              </HackerButton>
            </Link>
          </Grid>
          <Grid size={{ xs: 1 }}>
            <Link href={clientRoutes.service.configuration}>
              <HackerButton
                color="blue"
                icon={<Settings size={24} />}
                fullWidth
              >
                Configuración
              </HackerButton>
            </Link>
          </Grid>
          <Grid size={{ xs: 1 }}>
            <Link href={clientRoutes.service.support}>
              <HackerButton
                color="purple"
                icon={<Rocket size={24} />}
                fullWidth
              >
                Contactar Soporte
              </HackerButton>
            </Link>
          </Grid>
          <Grid size={{ xs: 1 }}>
            <Link href={clientRoutes.service.exit}>
              <HackerButton color="red" icon={<Power size={24} />} fullWidth>
                Exit
              </HackerButton>
            </Link>
          </Grid>
        </Grid>
      </Stack>
      <Stack
        direction={{ xs: 'column', sm: 'row' }}
        spacing={{ xs: 4, sm: 0 }}
        sx={{
          animation: `glitch 10s infinite`,
          pb: { xs: 16, sm: 0 },
        }}
      >
        <Stack width={'100%'}>
          <Typography>&gt; Últimas Operaciones:</Typography>
          <List dense>
            <ListItem> • Operación: MVP con Beltran</ListItem>
            <ListItem> • Proyecto: Scraper distribuido</ListItem>
            <ListItem> • Misión: Aprobar distribuidos</ListItem>
          </List>
        </Stack>
        <Stack width={'100%'} spacing={1}>
          <Typography>&gt; Únete A La Legión:</Typography>
          <Typography>
            Estamos unidos en nuestra misión. Estamos en todas partes. No
            estamos en ningún lugar. Somos tú.
          </Typography>
        </Stack>
      </Stack>
    </Stack>
  )
}

export default SystemContainer

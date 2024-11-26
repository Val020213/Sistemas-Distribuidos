'use client'
import { tailwindColors } from '@/theme/tailwindColors'

import {
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  IconButton,
  ModalProps,
  Stack,
  Typography,
} from '@mui/material'
import { X } from 'lucide-react'
import { ReactNode } from 'react'

type Props = ModalProps & {
  title: string
  actions?: ReactNode
  width?:
    | 'sx'
    | 'sm'
    | 'md'
    | 'lg'
    | 'xl'
    | number
    | string
    | Record<string, string | number>
}

const HackerDialog = ({ title, actions, children, ...others }: Props) => {
  return (
    <Dialog
      {...others}
      sx={{
        ...others.sx,
        '& .MuiPaper-root': {
          backgroundColor: tailwindColors.gray[900],
          width: others.width ?? {
            xs: '100%',
            sm: '80%',
            md: '60%',
            lg: '50%',
          },
          boxShadow: `0px 0px 16px 2px ${tailwindColors.green[500]}`,
        },
      }}
    >
      <DialogTitle color={tailwindColors.green[500]}>
        <Stack
          direction={'row'}
          alignItems={'center'}
          justifyContent={'space-between'}
        >
          <Typography component={'h6'} fontSize={'1rem'}>
            {title}
          </Typography>
          <IconButton disableRipple>
            <X color={tailwindColors.green[500]} size={18} />
          </IconButton>
        </Stack>
      </DialogTitle>

      <DialogContent
        sx={{
          width: '100%',
        }}
      >
        {children}
      </DialogContent>
      <DialogActions>{actions}</DialogActions>
    </Dialog>
  )
}

export default HackerDialog

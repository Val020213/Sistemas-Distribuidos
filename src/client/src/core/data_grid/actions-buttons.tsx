import { Box } from '@mui/material'
import React from 'react'
import { HackerButton } from '../button/HackerButton'
import { FileDown } from 'lucide-react'

interface Props {
  disabled?: boolean
  onDownload?: () => void
  buttonsVariant?: 'text' | 'outlined'
}

const ActionsButtons = ({
  disabled,
  onDownload,
  buttonsVariant = 'outlined',
}: Props) => {
  return (
    <Box
      sx={{
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center',
        width: '100%',
        height: '100%',
      }}
    >
      <Box
        sx={{
          display: 'flex',
          justifyContent: 'center',
          alignItems: 'center',
          gap: 1,
        }}
      >
        {onDownload && (
          <HackerButton
            disabled={disabled}
            variant="Button"
            buttonVariant={buttonsVariant}
            sx={{
              width: '40px !important',
              height: '40px !important',
              px: '1px !important',
            }}
            onClick={onDownload}
          >
            <FileDown />
          </HackerButton>
        )}
      </Box>
    </Box>
  )
}

export default ActionsButtons

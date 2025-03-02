import { HackerButton } from '@/core/button/HackerButton'
import RHFInputWithLabel from '@/core/inputs/rhf/RHFInputWithLabel'
import { Stack } from '@mui/material'
import { useFormContext } from 'react-hook-form'
import { searchUrlSchemaType } from './searchUrlSchema'
import { ReactNode } from 'react'

const SearchUrlForm = ({
  onClose,
  children,
}: {
  onClose: () => void
  children?: ReactNode
}) => {
  const {
    formState: { isValid, isLoading, isDirty },
  } = useFormContext<searchUrlSchemaType>()

  return (
    <Stack spacing={4} alignItems={'end'}>
      <RHFInputWithLabel
        label={"> Motor_d3_B'usquedA"}
        name="searchObjetive"
        autoComplete="new-password"
        required
      />
      {children}
      <Stack direction={'row'} spacing={2}>
        <HackerButton
          variant="Button"
          disabled={!isValid || isLoading || !isDirty}
          type="submit"
          sx={{
            width: 180,
          }}
        >
          Atacar
        </HackerButton>
        <HackerButton
          variant="Button"
          onClick={onClose}
          color="red"
          sx={{
            width: 180,
          }}
        >
          Cancelar
        </HackerButton>
      </Stack>
    </Stack>
  )
}

export default SearchUrlForm

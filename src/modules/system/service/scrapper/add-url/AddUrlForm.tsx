import { HackerButton } from '@/core/button/HackerButton'
import RHFInputWithLabel from '@/core/inputs/rhf/RHFInputWithLabel'
import { Stack } from '@mui/material'
import { useFormContext } from 'react-hook-form'
import { addUrlSchemaType } from './addUrlSchema'

const AddUrlForm = ({ onClose }: { onClose: () => void }) => {
  const {
    formState: { isValid, isLoading },
  } = useFormContext<addUrlSchemaType>()

  return (
    <Stack spacing={4} alignItems={'end'}>
      <RHFInputWithLabel
        label={'Url'}
        name="url"
        autoComplete="new-password"
        required
      />
      <Stack direction={'row'} spacing={2}>
        <HackerButton
          variant="Button"
          disabled={!isValid || isLoading}
          type="submit"
          sx={{
            width: 180,
          }}
        >
          Añadir
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

export default AddUrlForm

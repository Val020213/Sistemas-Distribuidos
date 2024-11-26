import { HackerButton } from '@/core/button/HackerButton'
import RHFInputWithLabel from '@/core/inputs/rhf/RHFInputWithLabel'
import RHFPasswordInput from '@/core/inputs/rhf/RHFPasswordInput'
import { Stack } from '@mui/material'
import { useFormContext } from 'react-hook-form'

const RegisterForm = () => {
  const {
    formState: { isValid, isLoading },
  } = useFormContext()
  return (
    <Stack spacing={4}>
      <RHFInputWithLabel
        label={'Nombre del miembro'}
        name="username"
        autoComplete="new-password"
        required
        size="medium"
      />
      <RHFPasswordInput
        label={'Contraseña'}
        name="password"
        autoComplete="new-password"
        required
        size="medium"
      />
      <RHFPasswordInput
        label={'Contraseña'}
        name="confirmPassword"
        autoComplete="new-password"
        required
        size="medium"
      />
      <HackerButton
        variant="Button"
        disabled={!isValid || isLoading}
        type="submit"
      >
        Continuar
      </HackerButton>
    </Stack>
  )
}

export default RegisterForm

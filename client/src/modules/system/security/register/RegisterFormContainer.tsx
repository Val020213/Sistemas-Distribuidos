import { useForm } from 'react-hook-form'
import RegisterForm from './RegisterForm'
import { zodResolver } from '@hookform/resolvers/zod'
import useShowHackerMessage from '@/hooks/useShowHackerMessage'
import registerSchema, { registerSchemaType } from './registerSchema'
import { FormProvider } from '@/core/form/form-provider'

const RegisterFormContainer = () => {
  const hackerMessage = useShowHackerMessage()
  const methods = useForm<registerSchemaType>({
    resolver: zodResolver(registerSchema),
    defaultValues: {
      username: '',
      password: '',
    },
    mode: 'onChange',
  })

  function onSubmit(data: registerSchemaType) {
    hackerMessage(data.username + ' ' + data.password, 'info')
  }

  return (
    <FormProvider props={methods} onSubmit={onSubmit}>
      <RegisterForm />
    </FormProvider>
  )
}

export default RegisterFormContainer

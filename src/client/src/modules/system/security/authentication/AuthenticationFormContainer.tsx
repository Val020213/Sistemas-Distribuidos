import { useForm } from 'react-hook-form'
import authenticationSchema, {
  authenticationSchemaType,
} from './authenticationSchema'
import { FormProvider } from '@/core/form/form-provider'
import AuthenticationForm from './AuthenticationForm'
import { zodResolver } from '@hookform/resolvers/zod'
import useShowHackerMessage from '@/hooks/useShowHackerMessage'
import { useContext } from 'react'
import { GlobalContext } from '@/core/global-context/GlobalContext'
import { saveUserData } from '@/lib/cookies'

const AuthenticationFormContainer = () => {
  const hackerMessage = useShowHackerMessage()
  const globalContext = useContext(GlobalContext)
  const methods = useForm<authenticationSchemaType>({
    resolver: zodResolver(authenticationSchema),
    defaultValues: {
      username: '',
      password: '',
    },
    mode: 'onChange',
  })

  async function onSubmit(data: authenticationSchemaType) {
    hackerMessage(data.username + ' ' + data.password, 'info')
    // get token
    globalContext.user[1]({
      name: data.username,
      token: '132',
    })
    // save token
    await saveUserData('name', data.username)
    await saveUserData('token', '132')
  }

  return (
    <FormProvider props={methods} onSubmit={onSubmit}>
      <AuthenticationForm />
    </FormProvider>
  )
}

export default AuthenticationFormContainer

import { zodResolver } from '@hookform/resolvers/zod'
import { useForm } from 'react-hook-form'
import addUrlSchema, { addUrlSchemaType } from './addUrlSchema'
import { FormProvider } from '@/core/form/form-provider'
import AddUrlForm from './AddUrlForm'
import useShowHackerMessage from '@/hooks/useShowHackerMessage'
import { useRouter } from 'next/navigation'

const AddUrlFormContainer = () => {
  const hackerMessage = useShowHackerMessage()
  const navigate = useRouter()
  const methods = useForm<addUrlSchemaType>({
    resolver: zodResolver(addUrlSchema),
    defaultValues: {
      url: '',
    },
    mode: 'onChange',
  })

  function onSubmit(data: addUrlSchemaType) {
    hackerMessage(data.url, 'warn')
  }
  function onClose() {
    navigate.back()
  }

  return (
    <FormProvider props={methods} onSubmit={onSubmit}>
      <AddUrlForm onClose={onClose} />
    </FormProvider>
  )
}

export default AddUrlFormContainer

import { zodResolver } from '@hookform/resolvers/zod'
import { useForm } from 'react-hook-form'
import addUrlSchema, { addUrlSchemaType } from './addUrlSchema'
import { FormProvider } from '@/core/form/form-provider'
import AddUrlForm from './AddUrlForm'
import useShowHackerMessage from '@/hooks/useShowHackerMessage'
import { useRouter } from 'next/navigation'
import { fetchUrlService } from '@/services/url-service'

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

  async function onSubmit(data: addUrlSchemaType) {
    try {
      const response = await fetchUrlService(data.url)
      console.log(response)
      if (response.status === 200) {
        hackerMessage('Url agregada correctamente', 'success')
        navigate.back()
      } else {
        hackerMessage('Error al agregar la url', 'error')
      }
    } catch (error: unknown) {
      console.log(error)
      hackerMessage(`Error al agregar la url ${error}`, 'error')
    }
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

import { zodResolver } from '@hookform/resolvers/zod'
import { useForm } from 'react-hook-form'
import addUrlSchema, { addUrlSchemaType } from './addUrlSchema'
import { FormProvider } from '@/core/form/form-provider'
import AddUrlForm from './AddUrlForm'
import useShowHackerMessage from '@/hooks/useShowHackerMessage'
import { fetchUrlService } from '@/services/url-service'
import { tagsRoutes } from '@/routes/routes'
import { revalidateServerTags } from '@/routes/cache'

const AddUrlFormContainer = ({ onClose }: { onClose: () => void }) => {
  const hackerMessage = useShowHackerMessage()
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
      if (response.statusCode === 200) {
        hackerMessage('Url agregada correctamente', 'success')
        await revalidateServerTags(tagsRoutes.list)
        onClose()
      } else {
        hackerMessage(response.message, 'error')
      }
    } catch (error: unknown) {
      console.log(error)
      hackerMessage(`Error al agregar la url ${error}`, 'error')
    }
  }

  return (
    <FormProvider props={methods} onSubmit={onSubmit}>
      <AddUrlForm onClose={onClose} />
    </FormProvider>
  )
}

export default AddUrlFormContainer

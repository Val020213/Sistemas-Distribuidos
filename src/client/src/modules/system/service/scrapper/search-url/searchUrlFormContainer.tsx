'use client'
import { zodResolver } from '@hookform/resolvers/zod'
import { useForm } from 'react-hook-form'
import { FormProvider } from '@/core/form/form-provider'
import useShowHackerMessage from '@/hooks/useShowHackerMessage'

import searchUrlSchema, { searchUrlSchemaType } from './searchUrlSchema'
import SearchUrlForm from './searchUrlForm'
import { useState } from 'react'
import SearchResult from './searchResult'
import { GoogleApiSearchResponse } from '@/types/api'
import { fetchUrlService, webSearchApi } from '@/services/url-service'
import { tagsRoutes } from '@/routes/routes'
import { revalidateServerTags } from '@/routes/cache'

const SearchUrlFormContainer = ({ onClose }: { onClose: () => void }) => {
  const [apiResponse, setApiResponse] = useState<GoogleApiSearchResponse>()

  const hackerMessage = useShowHackerMessage()
  const methods = useForm<searchUrlSchemaType>({
    resolver: zodResolver(searchUrlSchema),
    defaultValues: {
      searchObjetive: '',
    },
    mode: 'onChange',
  })

  async function onSubmit(data: searchUrlSchemaType) {
    try {
      const response = await webSearchApi(data.searchObjetive)
      setApiResponse(response)
    } catch (error: unknown) {
      console.log(error)
      hackerMessage(
        `Han interceptado el ataque, intente de nuevo 8er8&*#@(187) ${error}`,
        'error'
      )
    }
  }

  async function onAttack(url: string) {
    try {
      const response = await fetchUrlService(url)
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
      <SearchUrlForm onClose={onClose}>
        <SearchResult response={apiResponse} onAttack={onAttack} />
      </SearchUrlForm>
    </FormProvider>
  )
}

export default SearchUrlFormContainer

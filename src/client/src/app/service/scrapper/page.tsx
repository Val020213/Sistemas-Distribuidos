import { ScrapperContainer } from '@/modules/system/service/scrapper/ScrapperContainer'
import { listUrlService } from '@/services/url-service'
import { Metadata } from 'next'

export const metadata: Metadata = {
  title: 'Scrapper Service',
  description: 'Scrapper url service',
}

export default async function Page() {
  const response = await listUrlService()

  if (response.statusCode !== 200) {
    console.error('Request Error', response)
    throw new Error('No data received from server')
  }

  return <ScrapperContainer data={response.data ?? []} />
}

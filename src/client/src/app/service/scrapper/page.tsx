import { ScrapperContainer } from '@/modules/system/service/scrapper/ScrapperContainer'
import { listUrlService } from '@/services/url-service'
import { Metadata } from 'next'
import { SearchParams } from 'next/dist/server/request/search-params'

export const metadata: Metadata = {
  title: 'Scrapper Service',
  description: 'Scrapper url service',
}

type Props = {
  readonly searchParams: SearchParams
}

export default async function Page({ searchParams }: Props) {
  const response = await listUrlService()

  if (response.statusCode !== 200) {
    return new Error('Error fetching data')
  }
  return (
    <ScrapperContainer
      data={response.data.body}
      searchParams={await searchParams}
    />
  )
}

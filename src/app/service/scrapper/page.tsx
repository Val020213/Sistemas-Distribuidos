import { UrlDataType } from '@/app/types/url_data_type'
import { ScrapperContainer } from '@/modules/system/service/scrapper/ScrapperContainer'
import { Metadata } from 'next'
import { generateMockData } from './generateMockData'
import { SearchParams } from 'next/dist/server/request/search-params'

export const metadata: Metadata = {
  title: 'Scrapper Service',
  description: 'Scrapper url service',
}

type Props = {
  readonly searchParams: SearchParams
}

export default async function Page({ searchParams }: Props) {
  const mockUrlData: UrlDataType[] = generateMockData(100)

  return (
    <ScrapperContainer data={mockUrlData} searchParams={await searchParams} />
  )
}

'use client'
import { UrlDataType } from '@/app/types/url_data_type'
import { ScrapperContainer } from '@/modules/system/service/scrapper/ScrapperContainer'
import { backendRoutes } from '@/routes/routes'
import { listUrlService } from '@/services/url-service'
import { HackerApiResponse } from '@/types/api'
// import { Metadata } from 'next'
import React, { useEffect, useState } from 'react'

// export const metadata: Metadata = {
//   title: 'Scrapper Service',
//   description: 'Scrapper url service',
// }

export default function Page() {
  const [response, setResponse] = useState<
    HackerApiResponse<UrlDataType[]> | undefined
  >()

  useEffect(() => {
    try {
      const fetchData = async () => {
        const result = await listUrlService()
        setResponse(result)
      }
      fetchData()
    } catch (error) {
      console.error(error)
      setResponse(undefined)
    }
  }, [])

  if (!response) {
    return <div>{backendRoutes.list}</div>
  }
  return <ScrapperContainer data={response.data.body} />
}

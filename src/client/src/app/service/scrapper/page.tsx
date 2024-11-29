'use client'
import { ScrapperContainer } from '@/modules/system/service/scrapper/ScrapperContainer'
import { listUrlService } from '@/services/url-service'
// import { Metadata } from 'next'
import React, { useEffect, useState } from 'react'

// export const metadata: Metadata = {
//   title: 'Scrapper Service',
//   description: 'Scrapper url service',
// }

export default function Page() {
  const [response, setResponse] = useState({})

  useEffect(() => {
    try {
      const fetchData = async () => {
        const result = await listUrlService()
        setResponse(result)
      }
      fetchData()
    } catch (error) {
      console.error(error)
      setResponse({ statusCode: 500, data: error })
    }
  }, [])

  if (response === {}) {
    return <div>Loading...</div>
  }

  if (response.statusCode !== 200) {
    return <div>{response.data}</div>
  }
  return <ScrapperContainer data={response.data.body} />
}

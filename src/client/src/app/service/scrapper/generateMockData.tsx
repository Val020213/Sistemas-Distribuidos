import { UrlDataType } from '@/app/types/url_data_type'

export function generateMockData(count: number): UrlDataType[] {
  const mockData: UrlDataType[] = []
  for (let i = 0; i < count; i++) {
    mockData.push({
      url: `https://example.com/page${i + 1}`,
      status: i % 2 ? 'complete' : 'in_progress',
    })
  }
  return mockData
}

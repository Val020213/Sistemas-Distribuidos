export const clientRoutes = {
  home: '/',
  service: {
    scrapper: 'service/scrapper',
    configuration: 'service/configuration',
    support: 'service/support',
    dashboard: 'service/dashboard',
    exit: 'service/exit',
  },
}

export const backendRoutes = {
  fetch: `${process.env.NEXT_PUBLIC_API_URL}/tasks`,
  list: `${process.env.NEXT_PUBLIC_API_URL}/tasks`,
  download: `${process.env.NEXT_PUBLIC_API_URL}/tasks/:id/content`,
}

export const tagsRoutes = {
  fetch: 'fetchUrlService',
  list: 'listUrlService',
  download: 'downloadUrlService',
}

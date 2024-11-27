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
  fetch: `${process.env.NEXT_PUBLIC_API_URL}/fetch`,
  list: `${process.env.NEXT_PUBLIC_API_URL}/list`,
}

export const tagsRoutes = {
  fetch: 'fetchUrlService',
  list: 'listUrlService',
}

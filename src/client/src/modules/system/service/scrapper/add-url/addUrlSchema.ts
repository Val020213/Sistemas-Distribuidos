import { z } from 'zod'

const addUrlSchema = z.object({
  url: z.string().min(1, { message: 'url is required' }),
})

type addUrlSchemaType = z.infer<typeof addUrlSchema>

export default addUrlSchema
export type { addUrlSchemaType }

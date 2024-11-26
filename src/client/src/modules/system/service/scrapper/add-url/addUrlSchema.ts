import { z } from 'zod'

const addUrlSchema = z.object({
  url: z.string().url({ message: 'Url invalida, revise el dato introducido' }),
})

type addUrlSchemaType = z.infer<typeof addUrlSchema>

export default addUrlSchema
export type { addUrlSchemaType }

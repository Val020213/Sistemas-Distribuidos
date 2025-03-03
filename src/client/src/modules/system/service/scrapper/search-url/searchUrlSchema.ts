import { z } from 'zod'

const searchUrlSchema = z.object({
  searchObjetive: z
    .string()
    .min(1, { message: 'El objetivo de búsqueda no puede estar vacío' }),
})

type searchUrlSchemaType = z.infer<typeof searchUrlSchema>

export default searchUrlSchema
export type { searchUrlSchemaType }

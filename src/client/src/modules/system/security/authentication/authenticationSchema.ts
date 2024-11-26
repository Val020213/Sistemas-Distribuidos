import { z } from 'zod'

const authenticationSchema = z.object({
  username: z.string().min(1, {
    message: 'El nombre de usuario es requerido',
  }),

  password: z.string().min(1, {
    message: 'La contraseña requerida',
  }),
})

type authenticationSchemaType = z.infer<typeof authenticationSchema>

export default authenticationSchema
export type { authenticationSchemaType }

import { z } from 'zod'

const registerSchema = z
  .object({
    username: z.string().min(1, {
      message: 'El nombre de usuario es requerido',
    }),
    password: z
      .string()
      .min(8, { message: 'La contraseña debe tener al menos 8 caracteres' })
      .regex(/[A-Z]/, {
        message: 'La contraseña debe contener al menos una letra mayúscula',
      })
      .regex(/[a-z]/, {
        message: 'La contraseña debe contener al menos una letra minúscula',
      })
      .regex(/[0-9]/, {
        message: 'La contraseña debe contener al menos un número',
      })
      .regex(/[^a-zA-Z0-9]/, {
        message: 'La contraseña debe contener al menos un carácter especial',
      }),

    confirmPassword: z.string().min(1, {
      message: 'La confirmación de la contraseña es requerida',
    }),
  })
  .refine((data) => data.password === data.confirmPassword, {
    message: 'Las contraseñas no coinciden',
    path: ['confirmPassword'],
  })

type registerSchemaType = z.infer<typeof registerSchema>

export default registerSchema
export type { registerSchemaType }

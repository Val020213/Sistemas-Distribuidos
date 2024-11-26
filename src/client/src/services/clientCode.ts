import { ErrorCodes } from '@/types/api'

export const API_ERROR_CODES: Record<ErrorCodes, string> = {
  '400': 'Ha ocurrido un error. Inténtelo de nuevo más tarde.',
  '401': 'No autorizado para realizar esta acción.',
  '403': 'No tiene permiso para realizar esta acción.',
  '500': 'Problemas con los servidores',
  '501': 'Acción no autorizada',
}

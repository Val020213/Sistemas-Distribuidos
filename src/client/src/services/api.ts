import { ApiResponse, CustomApiError, ErrorCodes } from '@/types/api'
import { API_ERROR_CODES } from './clientCode'

export const handleApiServerError = async <T>(
  response: Response
): Promise<ApiResponse<T>> => {
  if (response.status === 401) {
    return {
      error: true,
      title: 'Unauthorized',
      status: 401,
      detail: 'You are not authorized to access this resource',
    }
  }

  try {
    const error = await response.json()
    error.error = true
    console.log(response.url, error)
    return error
  } catch {
    console.log(response.url, response)
    return Promise.resolve({ status: 500, error: true })
  }
}

/**
 * Builds an API response from an awaitable promise.
 * @template T The type of the data in the response.
 * @param {Promise<T>} awaitable - The awaitable promise that resolves to the data.
 * @returns {Promise<ApiResponse<T>>} A promise that resolves to the API response.
 * @throws {ApiError} If the awaitable promise throws an ApiError.
 */
export const buildApiResponseAsync = async <T>(
  awaitable: Promise<T>
): Promise<ApiResponse<T>> => {
  try {
    const data = await awaitable
    return Promise.resolve({ data, error: false, status: 200 })
  } catch (e) {
    // If the error is an ApiError, return it with the appropriate error fields.
    if (isApiError(e)) {
      return { ...e, title: e.title, error: true }
    }
    // If the error is not an ApiError, return a generic 500 error.
    return { status: 500, error: true }
  }
}

export const getAuthCookieName = () => {
  return process.env.NEXT_PUBLIC_SESSION_COOKIE_NAME
}

export class ApiError extends Error {
  title: string
  status: number
  detail: string
  clientCode: ErrorCodes
  extensions?: Record<string, string>

  constructor(error: CustomApiError) {
    super()
    this.title = error.title ?? 'Error inesperado'
    this.status = error.status ?? 500
    this.detail =
      error.detail ?? 'Ha ocurrido un error inesperado, inténtelo mas tarde.'
    this.clientCode = error.clientCode ?? '500'
    if (error.extensions) this.extensions = error.extensions
  }

  toString() {
    return `${this.title} - ${this.status} - ${this.detail.toString()}`
  }
}

/**
 * Handles API errors by throwing an instance of ApiError with the appropriate
 * error details.
 *
 * @param {Response} response - The response object containing the error details.
 * @returns {Promise<void>} - Throws an instance of ApiError with the error details.
 * @throws {ApiError} - Throws an instance of ApiError with the error details.
 */

export const handleApiError = async (response: Response) => {
  // Handle 401 Unauthorized errors
  if (response.status === 401) {
    throw new ApiError({
      title: 'Unauthorized',
      status: 401,
      detail: resolveErrorCodes('401'), // Resolve error code
      clientCode: '401', // Set client code
    })
  }
  // Handle 500 Internal Server Errors
  if (response.status === 500) {
    throw new ApiError({
      title: 'Internal Server Error',
      status: 500,
      detail: resolveErrorCodes('500'), // Resolve error code
      clientCode: '500', // Set client code
    })
  }
  // Handle 403 Forbidden errors
  if (response.status === 403) {
    throw new ApiError({
      title: 'Forbidden',
      status: 403,
      detail: resolveErrorCodes('403'), // Resolve error code
      clientCode: '403', // Set client code
    })
  }

  // Handle all other errors
  const error = (await response.json()) as ApiError // Parse error response
  error.detail = resolveErrorCodes(error.clientCode) // Resolve error code
  throw new ApiError(error) // Throw ApiError with error details
}

/**
 * Resolves the error code from the client code.
 *
 * @param {ErrorCodes} clientCode - The client code for the error.
 * @return {string} The corresponding error message for the client code.
 */
const resolveErrorCodes = (clientCode: ErrorCodes): string => {
  // Check if the client code exists in the API_ERROR_CODES object
  if (clientCode in API_ERROR_CODES) {
    // If it exists, return the corresponding error message
    return API_ERROR_CODES[clientCode]
  }
  // If it doesn't exist, return a generic error message
  return 'Ha ocurrido un error inesperado, inténtelo mas tarde.'
}

/**
 * Determines if the given error is an instance of ApiError.
 *
 * @param {unknown} error - The error to check.
 * @returns {boolean} True if the error is an instance of ApiError, false otherwise.
 */
export const isApiError = (error: unknown): error is ApiError => {
  /**
   * Checks if the given error is an instance of ApiError.
   *
   * @param {unknown} error - The error to check.
   * @returns {boolean} True if the error is an instance of ApiError, false otherwise.
   */
  /**
   * Checks if the given error is an instance of ApiError.
   *
   * @param {unknown} error - The error to check.
   * @returns {boolean} True if the error is an instance of ApiError, false otherwise.
   */
  const isApiErrorInstance = (error: unknown): error is ApiError => {
    /**
     * Checks if the given error is an instance of ApiError.
     *
     * @param {unknown} error - The error to check.
     * @returns {boolean} True if the error is an instance of ApiError, false otherwise.
     */
    return error instanceof ApiError
  }

  // Check if the error is an instance of ApiError
  return isApiErrorInstance(error)
}

/**
 * This function is used to retrieve the error message from the given error.
 * The error can be of different types, such as a string, an instance of the Error class,
 * or an instance of the ApiError class.
 *
 * @param {unknown} error - The error to retrieve the message from.
 * @returns {string | undefined} The error message if it exists, otherwise undefined.
 */
export const getErrorMessage = (error: unknown): string | undefined => {
  // First, check if the error is a string. If it is, return it as the error message.
  if (typeof error === 'string') {
    // The error is a string, so return it as the error message.
    return error
  }

  // Next, check if the error is an instance of the Error class.
  // If it is, return the message property of the error object.
  if (error instanceof Error) {
    // The error is an instance of the Error class, so return its message property.
    return error.message
  }

  // Finally, check if the error is an instance of the ApiError class.
  // If it is, return the detail property of the error object.
  if (error instanceof ApiError) {
    // The error is an instance of the ApiError class, so return its detail property.
    return error.detail
  }

  // If none of the above conditions are met, return undefined.
  // This means that the error is not a string, an instance of the Error class,
  // or an instance of the ApiError class, so there is no error message to return.
  return undefined
}

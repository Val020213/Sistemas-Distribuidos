import { Dispatch, SetStateAction } from 'react'

export type GlobalContextProps = {
  user: [User | undefined, Dispatch<SetStateAction<User | undefined>>]
}

export type User = {
  name: string
  rol?: string
  token: string
}

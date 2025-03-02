import { tailwindColors } from '@/theme/tailwindColors'
import { InputBase, Stack, SxProps, Theme } from '@mui/material'
import { SearchCodeIcon } from 'lucide-react'

type Props = {
  search: string
  setSearch: (value: string) => void
  placeholder: string
  sx?: SxProps
}

export const SearchBar = ({ search, setSearch, placeholder, sx }: Props) => {
  return (
    <Stack
      direction={'row'}
      alignItems={'center'}
      borderColor={`${tailwindColors.green[500]}`}
      border={1}
      borderRadius={0.5}
      paddingY={0.5}
      paddingX={2}
      sx={{
        bgcolor: tailwindColors.gray[900],
        '&:hover': {
          borderColor: `${tailwindColors.green[500]}`,
          boxShadow: `0 0 0 2px ${tailwindColors.gray[600]}`,
        },
        '&:focus-within': {
          borderColor: `${tailwindColors.green[500]}`,
          boxShadow: `0 0 0 2px ${tailwindColors.green[500]}`,
        },
        ...sx,
      }}
    >
      <InputBase
        value={search}
        placeholder={placeholder}
        onChange={(e) => setSearch(e.target.value)}
        sx={{
          color: `${tailwindColors.green[500]} !important`,
          flex: 1,
          fontSize: '0.875rem',
        }}
      />
      <SearchCodeIcon size="20px" />
    </Stack>
  )
}

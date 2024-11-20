import { Box, Stack, SxProps } from '@mui/material'
import { HackerButton } from '../button/HackerButton'
import { ArrowLeft, ArrowRight } from 'lucide-react'

interface Props {
  onClickPrev: () => void
  onClickNext: () => void
  onChangeCurrentPage: (
    event: React.MouseEvent<HTMLButtonElement> | null,
    page: number
  ) => void
  numbers: number[]
  currentPage: number
  sx?: SxProps
  disabledNext?: boolean
  disabledPrev?: boolean
}

export default function CustomPagination({
  currentPage,
  numbers,
  onChangeCurrentPage,
  onClickNext,
  onClickPrev,
  disabledNext,
  disabledPrev,
  sx,
}: Readonly<Props>) {
  return (
    <Box
      sx={{
        display: 'flex',
        gap: 1.3,
        width: '100%',
        ...sx,
      }}
    >
      <HackerButton
        size="small"
        variant="Button"
        onClick={onClickPrev}
        disabled={disabledPrev}
      >
        <ArrowLeft />
      </HackerButton>
      <Stack
        direction={'row'}
        spacing={1.3}
        overflow={'scroll'}
        sx={{
          '::-webkit-scrollbar': {
            display: 'none',
          },
          msOverflowStyle: 'none' /* IE and Edge */,
          scrollbarWidth: 'none' /* Firefox */,
        }}
      >
        {numbers.map((number) => (
          <HackerButton
            size="small"
            variant={'Button'}
            color={currentPage === number - 1 ? 'orange' : 'green'}
            key={number}
            onClick={(event) => onChangeCurrentPage(event, number - 1)}
            sx={{
              minWidth: 48,
              borderRadius: 1,
            }}
          >
            {number}
          </HackerButton>
        ))}
      </Stack>
      <HackerButton
        size="small"
        variant="Button"
        onClick={onClickNext}
        disabled={disabledNext}
      >
        <ArrowRight />
      </HackerButton>
    </Box>
  )
}

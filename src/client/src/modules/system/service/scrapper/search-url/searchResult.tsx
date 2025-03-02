import ActionsButtons from '@/core/data_grid/actions-buttons'
import { tailwindColors } from '@/theme/tailwindColors'
import { GoogleApiSearchResponse, GoogleSearchResults } from '@/types/api'
import { Box, Stack, Typography } from '@mui/material'

type Props = {
  response?: GoogleApiSearchResponse
  onAttack: (url: string) => void
}

const ResultCard = ({
  result,
  onAttack,
}: {
  result: GoogleSearchResults
  onAttack: (url: string) => void
}) => {
  return (
    <Stack
      sx={{
        position: 'relative',
        p: 2,
        '&:hover': {
          backgroundColor: tailwindColors.green[900],
        },
      }}
    >
      <Typography fontWeight={'bold'}>{result.title}</Typography>
      <Typography
        component={'a'}
        href={result.link}
        color={`${tailwindColors.green[300]}`}
      >
        {result.link}
      </Typography>
      <Typography variant="body2">{result.snippet}</Typography>
      <Box
        sx={{
          position: 'absolute',
          top: 4,
          right: 4,
        }}
      >
        <ActionsButtons
          buttonsVariant="text"
          onDownload={() => onAttack(result.link)}
        />
      </Box>
    </Stack>
  )
}

const SearchResult = ({ response, onAttack }: Props) => {
  return (
    <Stack
      sx={{
        borderRadius: 0,
        padding: 1,
        border: `1px dashed ${tailwindColors.green[400]}`,
        width: '100%',
        height: '350px',
        overflowY: 'scroll',
      }}
    >
      <Typography fontSize={'0.625rem'}>
        {response &&
          `${response.items.length} objetivos para atacar encontrados`}
      </Typography>
      <Stack mt={3}>
        {response?.items.map((result, index) => (
          <ResultCard
            key={result.snippet}
            result={result}
            onAttack={onAttack}
          />
        ))}
      </Stack>
    </Stack>
  )
}
export default SearchResult

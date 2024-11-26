import { UrlStatus } from '@/app/types/url_data_type'
import { tailwindColors } from '@/theme/tailwindColors'
import { Chip } from '@mui/material'

const UrlStatusChip = ({ urlStatus }: { urlStatus: UrlStatus }) => {
  const color =
    urlStatus === 'pending'
      ? 'yellow'
      : urlStatus === 'scrapped'
      ? 'green'
      : 'red'
  return (
    <Chip
      sx={{
        bgcolor: `${tailwindColors[color][900]}`,
        color: `${tailwindColors[color][300]}`,
        borderRadius: 0,
        width: '100px',
        height: '30px',
        fontSize: '0.625rem',
      }}
      label={urlStatus.toUpperCase()}
    />
  )
}

export default UrlStatusChip

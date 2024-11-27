'use client'
import { UrlDataType } from '@/app/types/url_data_type'
import ActionsButtons from '@/core/data_grid/actions-buttons'
import UrlStatusChip from '@/core/data_grid/url-chip'
import HackerDataGrid from '@/core/data_grid/hacker-data-grid'
import useShowHackerMessage from '@/hooks/useShowHackerMessage'
import { Box, Stack } from '@mui/material'
import { GridColDef } from '@mui/x-data-grid'
import { HackerButton } from '@/core/button/HackerButton'
import Link from 'next/link'
import AddUrl from './add-url/AddUrl'
import { SearchParams } from 'next/dist/server/request/search-params'

type Props = {
  data: UrlDataType[]
  searchParams?: SearchParams
}

export const ScrapperContainer = ({ data, searchParams }: Props) => {
  const hackerMessages = useShowHackerMessage()
  const columns: GridColDef<UrlDataType>[] = [
    {
      field: 'id',
      headerName: 'ID',
      headerClassName: 'header-class',
    },
    {
      field: 'url',
      headerClassName: 'header-class',
      headerName: 'URL',
      flex: 2,
    },
    {
      field: 'status',
      headerClassName: 'header-class',
      headerName: 'STATUS',
      flex: 1,
      renderCell: (params) => {
        return (
          <Stack
            justifyContent={'center'}
            alignItems={'center'}
            width={'100%'}
            height={'100%'}
          >
            <UrlStatusChip urlStatus={params.value} />
          </Stack>
        )
      },
    },
    {
      field: 'actions',
      headerClassName: 'header-class',
      headerName: 'Action',
      align: 'center',
      headerAlign: 'center',
      sortable: false,
      disableColumnMenu: true,
      disableExport: true,
      renderCell: (params) => {
        return (
          <ActionsButtons
            onDownload={() =>
              hackerMessages(
                `Trying to download file with id ${params.id}...\n> Error: this functionality is in progress`
              )
            }
          />
        )
      },
    },
  ]

  return (
    <Stack spacing={4} sx={{ maxWidth: 'md', width: '100%' }}>
      <Stack spacing={4} direction={'row'}>
        <Link href={'/'}>
          <HackerButton
            variant="Button"
            color="green"
            sx={{
              width: '168px',
            }}
          >
            &lt; Ir al Sistema
          </HackerButton>
        </Link>
        <Link href={'?currentModal=addUrl'}>
          <HackerButton
            variant="Button"
            color="green"
            sx={{
              width: '168px',
            }}
          >
            + Agregar URL
          </HackerButton>
        </Link>
      </Stack>
      <HackerDataGrid columns={columns} data={data} />
      <Box height={8} />
      <AddUrl currentModal={searchParams?.currentModal as string} />
    </Stack>
  )
}

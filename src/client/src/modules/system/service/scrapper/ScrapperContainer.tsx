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
import { downloadUrlService } from '@/services/url-service'
import { useState } from 'react'

type Props = {
  data: UrlDataType[]
}

export const ScrapperContainer = ({ data }: Props) => {
  const hackerMessages = useShowHackerMessage()

  async function handleDownload(id: string) {
    const response = await downloadUrlService(id)

    if (response.data) {
      const url = window.URL.createObjectURL(
        new Blob([response.data.content ?? ''])
      )
      const link = document.createElement('a')
      link.href = url
      link.setAttribute('download', `${id}.html`)
      document.body.appendChild(link)
      link.click()
      hackerMessages('Descargando archivo', 'success')
      return
    }

    hackerMessages('Error al descargar archivo', 'error')
  }

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
            disabled={params.row.status !== 'complete'}
            onDownload={() => handleDownload(params.row.id.toString())}
          />
        )
      },
    },
  ]

  const [openModal, setOpenModal] = useState<string>('')

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

        <HackerButton
          variant="Button"
          color="green"
          sx={{
            width: '168px',
          }}
          onClick={() => setOpenModal('addUrl')}
        >
          + Agregar URL
        </HackerButton>
      </Stack>
      <HackerDataGrid columns={columns} data={data} />
      <Box height={8} />
      <AddUrl currentModal={openModal} onClose={() => setOpenModal('')} />
    </Stack>
  )
}

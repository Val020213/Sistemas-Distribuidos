'use client'
import {
  DataGrid,
  DataGridProps,
  GridColDef,
  gridPageCountSelector,
  GridPagination,
  GridRowClassNameParams,
  GridValidRowModel,
  useGridApiContext,
  useGridSelector,
} from '@mui/x-data-grid'
import {
  Box,
  InputBase,
  LinearProgress,
  Stack,
  styled,
  TablePaginationProps,
  Typography,
} from '@mui/material'
import { GridOverlay } from '@mui/x-data-grid'
import { useMemo, useState } from 'react'
import HackerPagination from './CustomPagination'
import { tailwindColors } from '@/theme/tailwindColors'
import { SearchCodeIcon } from 'lucide-react'

type Props = Omit<DataGridProps, 'pagination' | 'rows'> & {
  title?: string
  data: GridValidRowModel[]
  columns: GridColDef[]
  size?: 'small' | 'medium' | 'large'
}

export const dataGridSlots = {
  loadingOverlay: CustomLoadingOverlay,
  noRowsOverlay: CustomNoRowsOverlay,
  pagination: CustomPagination,
}

export const MAX_PAGINATION_NUMBERS = 5

export const addOddOrEvenClassname = (params: GridRowClassNameParams) => {
  return params.indexRelativeToCurrentPage % 2 === 0 ? 'odd-row' : 'even-row'
}

export const StyledDataGrid = styled(DataGrid)(() => ({
  border: 'none',
  width: '100%',
  borderRadius: '0px !important',
  '& .header-class': {
    fontWeight: 600,
    backgroundColor: tailwindColors.green[900],
    borderRadius: '0px !important',
    color: tailwindColors.green[300],
    borderTop: `2px solid ${tailwindColors.green[500]}`,
    borderBottom: `2px solid ${tailwindColors.green[500]}`,
  },
  '& .MuiSvgIcon-root': {
    color: tailwindColors.green[500],
  },
  '& .MuiDataGrid-sortIcon': {
    color: tailwindColors.green[500],
  },
  '& .even-row': {
    backgroundColor: tailwindColors.gray[900],
  },
  '& .odd-row': {
    backgroundColor: tailwindColors.gray[800],
  },
  '& .MuiDataGrid-cell': {
    color: tailwindColors.green[400],
    borderBottom: `1px solid ${tailwindColors.gray[900]}`,
    borderTop: `1px solid ${tailwindColors.gray[900]}`,
  },
  '& .MuiDataGrid-filler': {
    backgroundColor: tailwindColors.gray[900],
  },
  '& .MuiDataGrid-columnHeader:focus': {
    outline: 'none',
  },
  '& .MuiTablePagination-displayedRows': {
    color: tailwindColors.green[400],
  },
}))

const HackerDataGrid: React.FC<Props> = ({
  title,
  data,
  columns,
  size = 'small',
  ...other
}) => {
  const overlayHeight = {
    small: '300px',
    medium: '468px',
    large: '600px',
  }

  const [search, setSearch] = useState<string>('')

  const dataContent = useMemo(() => {
    if (search.length > 0) {
      return data.filter((row) => {
        return Object.values(row).some((value) =>
          value.toString().toLowerCase().includes(search.toLowerCase())
        )
      })
    }
    return data
  }, [data, search])

  return (
    <Stack
      position="relative"
      width={'100%'}
      spacing={3}
      minHeight={overlayHeight[size]}
      justifyContent={'space-between'}
      bgcolor={tailwindColors.black}
      paddingTop={4}
      sx={{
        animation: 'glow 1.5s ease-in-out infinite alternate',
      }}
    >
      <Stack px={4} spacing={3}>
        <Typography variant="h3" lineHeight="1.75rem">
          {title?.toUpperCase() ?? 'SISTEMA DE DATOS'}
        </Typography>
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
          }}
        >
          <InputBase
            value={search}
            placeholder="Buscar url"
            onChange={(e) => setSearch(e.target.value)}
            sx={{
              color: `${tailwindColors.green[500]} !important`,
              flex: 1,
              fontSize: '0.875rem',
            }}
          />
          <SearchCodeIcon size="20px" />
        </Stack>
      </Stack>

      <StyledDataGrid
        {...other}
        rows={dataContent}
        columns={columns}
        slots={dataGridSlots}
        sx={{
          '--DataGrid-overlayHeight': overlayHeight[size],
        }}
        initialState={{
          pagination: { paginationModel: { pageSize: 10 } },
        }}
        pagination
        columnHeaderHeight={40}
        getRowClassName={addOddOrEvenClassname}
        rowHeight={size === 'small' ? 52 : 60}
        loading={other.loading}
      />
    </Stack>
  )
}
function CustomLoadingOverlay() {
  return (
    <GridOverlay>
      <Box sx={{ position: 'absolute', top: 0, width: '100%' }}>
        <LinearProgress />
      </Box>
    </GridOverlay>
  )
}

export default HackerDataGrid

const StyledGridOverlay = styled('div')(({ theme }) => ({
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'center',
  justifyContent: 'center',
  height: '100%',
  '& .no-rows-primary': {
    fill: '#3D4751',
    ...theme.applyStyles('light', {
      fill: '#AEB8C2',
    }),
  },
  '& .no-rows-secondary': {
    fill: '#1D2126',
    ...theme.applyStyles('light', {
      fill: '#E8EAED',
    }),
  },
}))

function CustomNoRowsOverlay() {
  return (
    <StyledGridOverlay
      sx={{
        '& .no-rows-primary': {
          fill: tailwindColors.green[500],
        },
        '& .no-rows-secondary': {
          fill: tailwindColors.green[900],
        },
      }}
    >
      <svg
        xmlns="http://www.w3.org/2000/svg"
        fill="none"
        width={180}
        viewBox="0 0 452 257"
        aria-hidden
        focusable="false"
      >
        <path
          className="no-rows-primary"
          d="M348 69c-46.392 0-84 37.608-84 84s37.608 84 84 84 84-37.608 84-84-37.608-84-84-84Zm-104 84c0-57.438 46.562-104 104-104s104 46.562 104 104-46.562 104-104 104-104-46.562-104-104Z"
        />
        <path
          className="no-rows-primary"
          d="M308.929 113.929c3.905-3.905 10.237-3.905 14.142 0l63.64 63.64c3.905 3.905 3.905 10.236 0 14.142-3.906 3.905-10.237 3.905-14.142 0l-63.64-63.64c-3.905-3.905-3.905-10.237 0-14.142Z"
        />
        <path
          className="no-rows-primary"
          d="M308.929 191.711c-3.905-3.906-3.905-10.237 0-14.142l63.64-63.64c3.905-3.905 10.236-3.905 14.142 0 3.905 3.905 3.905 10.237 0 14.142l-63.64 63.64c-3.905 3.905-10.237 3.905-14.142 0Z"
        />
        <path
          className="no-rows-secondary"
          d="M0 10C0 4.477 4.477 0 10 0h380c5.523 0 10 4.477 10 10s-4.477 10-10 10H10C4.477 20 0 15.523 0 10ZM0 59c0-5.523 4.477-10 10-10h231c5.523 0 10 4.477 10 10s-4.477 10-10 10H10C4.477 69 0 64.523 0 59ZM0 106c0-5.523 4.477-10 10-10h203c5.523 0 10 4.477 10 10s-4.477 10-10 10H10c-5.523 0-10-4.477-10-10ZM0 153c0-5.523 4.477-10 10-10h195.5c5.523 0 10 4.477 10 10s-4.477 10-10 10H10c-5.523 0-10-4.477-10-10ZM0 200c0-5.523 4.477-10 10-10h203c5.523 0 10 4.477 10 10s-4.477 10-10 10H10c-5.523 0-10-4.477-10-10ZM0 247c0-5.523 4.477-10 10-10h231c5.523 0 10 4.477 10 10s-4.477 10-10 10H10c-5.523 0-10-4.477-10-10Z"
        />
      </svg>
      <Box sx={{ mt: 2 }}>
        <Typography component={'h2'} fontWeight={500} sx={{}}>
          SIN URLS QUE MOSTRAR
        </Typography>
      </Box>
    </StyledGridOverlay>
  )
}

function Pagination({
  page,
  onPageChange,
  count,
}: Pick<
  TablePaginationProps & { count: number },
  'page' | 'onPageChange' | 'className' | 'count'
>) {
  const apiRef = useGridApiContext()
  const pageCount = useGridSelector(apiRef, gridPageCountSelector)
  const pageInterval = useMemo(() => {
    return Math.floor(page / MAX_PAGINATION_NUMBERS)
  }, [page])
  return (
    <Stack p={4} pr={2}>
      <HackerPagination
        currentPage={page}
        numbers={Array.from(
          {
            length: Math.floor(
              (count - pageInterval * MAX_PAGINATION_NUMBERS) /
                MAX_PAGINATION_NUMBERS
            ),
          },
          (_, i) => pageInterval * MAX_PAGINATION_NUMBERS + i + 1
        )}
        onChangeCurrentPage={onPageChange}
        onClickNext={() =>
          onPageChange(null, (pageInterval + 1) * MAX_PAGINATION_NUMBERS)
        }
        onClickPrev={() =>
          onPageChange(
            null,
            Math.max(pageInterval * MAX_PAGINATION_NUMBERS - 1, 0)
          )
        }
        disabledNext={(pageInterval + 1) * MAX_PAGINATION_NUMBERS >= pageCount}
        disabledPrev={pageInterval === 0}
        sx={{
          display: 'flex',
          gap: 1.3,
          width: '100%',
          justifyContent: 'center',
        }}
      />
    </Stack>
  )
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
function CustomPagination(props: any) {
  return <GridPagination ActionsComponent={Pagination} {...props} />
}

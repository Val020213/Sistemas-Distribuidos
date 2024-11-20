import { useState } from 'react'

const useCustomPagination = <T>(rows: T[], rowsPerPage: number) => {
  const [currentPage, setCurrentPage] = useState(1)

  const lastIndex = currentPage * rowsPerPage
  const firstIndex = lastIndex - rowsPerPage

  const rowData = rows.slice(firstIndex, lastIndex)

  const nPage = Math.ceil(rows.length / rowsPerPage)

  const numbers = [...Array(nPage + 1).keys()].slice(1)

  const changeCurrentPage = (id: number) => {
    setCurrentPage(id)
  }

  const prevPage = () => {
    if (currentPage !== 1) {
      setCurrentPage(currentPage - 1)
    }
  }

  const nextPage = () => {
    if (currentPage !== nPage) {
      const newCurrentPage = currentPage + 1
      setCurrentPage(newCurrentPage)
    }
  }

  return {
    changeCurrentPage,
    prevPage,
    nextPage,
    currentPage,
    numbers,
    rowData,
  }
}

export default useCustomPagination

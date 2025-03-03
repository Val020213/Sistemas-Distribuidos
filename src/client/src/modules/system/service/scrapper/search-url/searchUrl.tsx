import HackerDialog from '@/core/dialog/HackerDIalog'
import SearchUrlFormContainer from './searchUrlFormContainer'

const SearchUrl = ({
  currentModal,
  onClose,
}: {
  currentModal?: string
  onClose: () => void
}) => {
  return (
    <HackerDialog
      title=">_://Iniciar búsqueda"
      open={currentModal === 'searchUrl'}
      onClose={onClose}
    >
      <SearchUrlFormContainer onClose={onClose} />
    </HackerDialog>
  )
}

export default SearchUrl

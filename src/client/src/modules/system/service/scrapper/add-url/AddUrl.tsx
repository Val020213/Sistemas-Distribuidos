import HackerDialog from '@/core/dialog/HackerDIalog'
import AddUrlFormContainer from './AddUrlFormContainer'

const AddUrl = ({
  currentModal,
  onClose,
}: {
  currentModal?: string
  onClose: () => void
}) => {
  return (
    <HackerDialog
      title="Procesar Url"
      open={currentModal === 'addUrl'}
      onClose={onClose}
    >
      <AddUrlFormContainer onClose={onClose} />
    </HackerDialog>
  )
}

export default AddUrl

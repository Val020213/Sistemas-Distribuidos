import HackerDialog from '@/core/dialog/HackerDIalog'
import AddUrlFormContainer from './AddUrlFormContainer'

const AddUrl = ({ currentModal }: { currentModal?: string }) => {
  return (
    <HackerDialog title="Procesar Url" open={currentModal === 'addUrl'}>
      <AddUrlFormContainer />
    </HackerDialog>
  )
}

export default AddUrl

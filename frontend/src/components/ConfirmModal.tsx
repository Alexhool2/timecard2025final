import { Button, Modal } from "react-bootstrap"

interface ConfirmModalProps {
  show: boolean
  onHide: () => void
  onConfirm: () => void
  title: string
  body: string
}

const ConfirmModal = ({
  show,
  onHide,
  onConfirm,
  title,
  body,
}: ConfirmModalProps) => {
  return (
    <Modal show={show} onHide={onHide}>
      <Modal.Header closeButton>
        <Modal.Title>{title}</Modal.Title>
        <Modal.Body>{body}</Modal.Body>
        <Modal.Footer>
          <Button variant="secondary" onClick={onHide}>
            No
          </Button>
          <Button variant="primary" onClick={onConfirm}>
            Yes
          </Button>
        </Modal.Footer>
      </Modal.Header>
    </Modal>
  )
}

export default ConfirmModal

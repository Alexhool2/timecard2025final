import { Container, Row, Col, Form, Button, Modal } from "react-bootstrap"
import UserSearch from "../components/UserSearch"
import { useState } from "react"
import { changePassword } from "../network/timecard.api"
import { useNavigate } from "react-router-dom"

const ResetPassword = () => {
  const [userId, setUserId] = useState<number | null>(null)

  const [validated, setValidated] = useState(false)
  const [showModal, setShowModal] = useState(false)
  const [newPassword, setNewPassword] = useState<string>("")
  const navigate = useNavigate()

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    event.stopPropagation()

    if (!userId) {
      alert("Please select a user first")
    }
    if (userId === null) {
      console.log("user id is null")
      return
    }

    const form = event?.currentTarget
    if (form.checkValidity())
      try {
        await changePassword(userId, newPassword)

        setShowModal(true)
      } catch (error) {
        console.log("Error creating User", error)
      }
    setValidated(true)
  }

  return (
    <Container className="mt-5">
      <Row className="justify-content-center">
        <Col md={6}>
          <h3 className="text-center">Select user</h3>
          <UserSearch onSelectUser={(id) => setUserId(id)} />
        </Col>
      </Row>
      <Row className="justify-content-center mt-3">
        <Form
          noValidate
          validated={validated}
          onSubmit={handleSubmit}
          className="w-50"
        >
          <Row className="mb-3 justify-content-center">
            <Form.Group as={Col} md="12" controlId="validationCustom01">
              <Form.Label>New Password</Form.Label>
              <Form.Control
                required
                type="text"
                placeholder="Enter new password"
                value={newPassword}
                onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
                  setNewPassword(e.target.value)
                }
              />
              <Form.Control.Feedback>Looks good!</Form.Control.Feedback>
            </Form.Group>
          </Row>
          <Row className="justify-content-center">
            <Col md="12" className="text-center">
              <Button variant="dark" type="submit">
                Change Password
              </Button>
            </Col>
          </Row>
        </Form>
      </Row>
      <Modal
        show={showModal}
        onHide={() => {
          setShowModal(false)
          navigate("/users")
        }}
      >
        <Modal.Header closeButton>
          <Modal.Title>Password reset successfully</Modal.Title>
        </Modal.Header>
        <Modal.Body>
          <p>
            <strong>UserId:</strong>
            {userId}
          </p>

          <p>
            <strong>New Password:</strong>
            {newPassword}
          </p>
          <p>Please save the new information</p>
          <Modal.Footer>
            <Button
              variant="dark"
              onClick={() => {
                setShowModal(false)
                navigate("/users")
              }}
            >
              OK
            </Button>
          </Modal.Footer>
        </Modal.Body>
      </Modal>
    </Container>
  )
}

export default ResetPassword

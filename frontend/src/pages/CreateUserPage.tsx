import { useState } from "react"
import Button from "react-bootstrap/Button"
import Col from "react-bootstrap/Col"
import Form from "react-bootstrap/Form"
import InputGroup from "react-bootstrap/InputGroup"
import Row from "react-bootstrap/Row"
import { createUser, CreateUserInterface } from "../network/timecard.api"
import { useNavigate } from "react-router-dom"
import { Modal } from "react-bootstrap"

function CreateUserPage() {
  const [validated, setValidated] = useState(false)
  const [showModal, setShowModal] = useState(false)

  const [formData, setFormData] = useState<CreateUserInterface>({
    firstName: "",
    lastName: "",
    userName: "",
    password: "",
    email: "",
    role: "",
  })
  const navigate = useNavigate()

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    event.stopPropagation()

    const form = event?.currentTarget
    if (form.checkValidity())
      try {
        await createUser(formData)

        setShowModal(true)
      } catch (error) {
        console.log("Error creating User", error)
      }
    setValidated(true)
  }

  const handleInputChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = event.target
    setFormData((prevData) => ({ ...prevData, [name]: value }))
  }
  return (
    <div className="d-flex justify-content-center align-items-center min-vh-100">
      <Form
        noValidate
        validated={validated}
        onSubmit={handleSubmit}
        className="w-50"
      >
        <Row className="mb-3 justify-content-center">
          <Form.Group as={Col} md="12" controlId="validationCustom01">
            <Form.Label>First name</Form.Label>
            <Form.Control
              required
              type="text"
              placeholder="First name"
              name="firstName"
              value={formData.firstName}
              onChange={handleInputChange}
            />
            <Form.Control.Feedback>Looks good!</Form.Control.Feedback>
          </Form.Group>
        </Row>
        <Row className="mb-3 justify-content-center">
          <Form.Group as={Col} md="12" controlId="validationCustom02">
            <Form.Label>Last name</Form.Label>
            <Form.Control
              required
              type="text"
              placeholder="Last name"
              name="lastName"
              value={formData.lastName}
              onChange={handleInputChange}
            />
            <Form.Control.Feedback>Looks good!</Form.Control.Feedback>
          </Form.Group>
        </Row>
        <Row className="mb-3 justify-content-center">
          <Form.Group as={Col} md="12" controlId="validationCustomUsername">
            <Form.Label>Username</Form.Label>
            <InputGroup hasValidation>
              <InputGroup.Text id="inputGroupPrepend">@</InputGroup.Text>
              <Form.Control
                type="text"
                placeholder="Username"
                aria-describedby="inputGroupPrepend"
                required
                name="userName"
                value={formData.userName}
                onChange={handleInputChange}
              />
              <Form.Control.Feedback type="invalid">
                Please choose a username.
              </Form.Control.Feedback>
            </InputGroup>
          </Form.Group>
        </Row>
        <Row className="mb-3 justify-content-center">
          <Form.Group as={Col} md="12" controlId="validationCustom03">
            <Form.Label>Password</Form.Label>
            <Form.Control
              type="password"
              placeholder="Password"
              required
              name="password"
              value={formData.password}
              onChange={handleInputChange}
            />
            <Form.Control.Feedback type="invalid">
              Please provide a valid password
            </Form.Control.Feedback>
          </Form.Group>
        </Row>
        <Row className="mb-3 justify-content-center">
          <Form.Group as={Col} md="12" controlId="validationCustom04">
            <Form.Label>Email</Form.Label>
            <Form.Control
              type="email"
              placeholder="email"
              required
              name="email"
              value={formData.email}
              onChange={handleInputChange}
            />
            <Form.Control.Feedback type="invalid">
              Please provide a valid email.
            </Form.Control.Feedback>
          </Form.Group>
        </Row>
        <Row className="mb-3 justify-content-center">
          <Form.Group as={Col} md="8" controlId="validationCustomOnRole">
            <Form.Label>Role</Form.Label>
            <Form.Select
              aria-label="Default select example"
              required
              name="role"
              value={formData.role}
              onChange={(e: React.ChangeEvent<HTMLSelectElement>) => {
                setFormData((prev) => ({ ...prev, role: e.target.value }))
              }}
            >
              <option value="">Select a role...</option>
              <option value="user">user</option>
              <option value="admin">admin</option>
            </Form.Select>
            <Form.Control.Feedback type="invalid">
              Please select a role.
            </Form.Control.Feedback>
          </Form.Group>
        </Row>
        <Form.Group className="mb-3">
          <Form.Check
            required
            label="Agree to terms and conditions"
            feedback="You must agree before submitting."
            feedbackType="invalid"
          />
        </Form.Group>
        <Button type="submit">Submit form</Button>
      </Form>
      <Modal
        show={showModal}
        onHide={() => {
          setShowModal(false)
          navigate("/users")
        }}
      >
        <Modal.Header closeButton>
          <Modal.Title>User Created successfully</Modal.Title>
        </Modal.Header>
        <Modal.Body>
          <p>
            <strong>Username:</strong>
            {formData.userName}
          </p>

          <p>
            <strong>Password:</strong>
            {formData.password}
          </p>
          <p>Please save the login information</p>
          <Modal.Footer>
            <Button
              variant="primary"
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
    </div>
  )
}

export default CreateUserPage

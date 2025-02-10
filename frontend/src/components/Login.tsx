import { useState } from "react"
import { useForm } from "react-hook-form"
import { LoginCredentials } from "../models/user"
import * as TimeCardApi from "../network/timecard.api"
import { UnauthorizedError } from "../errors/http_errors"
import { Alert, Button, Form, Modal } from "react-bootstrap"
import TextInputField from "./form/TextInputField"
import { LoggedUser } from "../network/timecard.api"

interface LoginProps {
  onLoginSuccessful: (user: LoggedUser) => void
  onDismiss: () => void
}

const Login = ({ onLoginSuccessful, onDismiss }: LoginProps) => {
  const [errorText, setErrorText] = useState<string | null>(null)

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<LoginCredentials>()

  async function onSubmit(credentials: LoginCredentials) {
    try {
      const user = await TimeCardApi.login(credentials)

      // change user state on app
      onLoginSuccessful(user)
    } catch (error) {
      console.error("Login error details:", error)
      if (error instanceof UnauthorizedError) {
        setErrorText(error.message)
      } else {
        setErrorText("An unexpected error occurred. Please try again.")
      }
    }
  }

  return (
    <Modal show onHide={onDismiss}>
      <Modal.Header closeButton>
        <Modal.Title>Login</Modal.Title>
      </Modal.Header>

      <Modal.Body>
        {errorText && <Alert variant="danger">{errorText}</Alert>}
        <Form onSubmit={handleSubmit(onSubmit)}>
          <TextInputField
            name="userName"
            label="Username"
            type="text"
            placeholder="Username"
            register={register}
            registerOptions={{ required: "Required" }}
            error={errors.userName}
          />
          <TextInputField
            name="password"
            label="Password"
            type="password"
            placeholder="Password"
            register={register}
            registerOptions={{ required: "Required" }}
            error={errors.password}
          />
          <Button
            type="submit"
            disabled={isSubmitting}
            className="w-100 bg-dark"
          >
            {isSubmitting ? "Loading..." : "Login"}
          </Button>
        </Form>
      </Modal.Body>
    </Modal>
  )
}

export default Login

import Navbar from "react-bootstrap/Navbar"
import Container from "react-bootstrap/Container"
import { Link } from "react-router-dom"
import { Nav } from "react-bootstrap"
import NavBarLoggedInView from "./NavBarLoggedInView"
import NavBarLoggedOutView from "./NavBarLoggedOutView"
import { LoggedUser } from "../network/timecard.api"

interface NavbarProps {
  loggedInUser: LoggedUser | null
  onLoginClicked: () => void
  onLogoutSuccessful: () => void
}

const NavBar = ({
  loggedInUser,
  onLoginClicked,
  onLogoutSuccessful,
}: NavbarProps) => {
  return (
    <Navbar bg="light" data-bs-theme="light" expand="sm" sticky="top">
      <Container>
        <Navbar.Brand as={Link} to="/">
          TimeCard App
        </Navbar.Brand>
        <Navbar.Toggle aria-controls="main-navbar" />
        <Navbar.Collapse id="main-navbar">
          <Nav>
            {loggedInUser?.role === "admin" && (
              <Nav.Link as={Link} to={"/userinfo"}>
                User info
              </Nav.Link>
            )}
            {loggedInUser?.role === "admin" && (
              <Nav.Link as={Link} to={"/createuser"}>
                Create New user
              </Nav.Link>
            )}
            {loggedInUser?.role === "admin" && (
              <Nav.Link as={Link} to={"/reset-password"}>
                Change password
              </Nav.Link>
            )}
          </Nav>
          <Nav className="ms-auto">
            {loggedInUser ? (
              <NavBarLoggedInView
                user={loggedInUser}
                onLogoutSuccessful={onLogoutSuccessful}
              />
            ) : (
              <NavBarLoggedOutView onLoginClicked={onLoginClicked} />
            )}
          </Nav>
        </Navbar.Collapse>
      </Container>
    </Navbar>
  )
}

export default NavBar

import { Button, Navbar } from "react-bootstrap"
import * as TimeCardApi from "../network/timecard.api"
import { LoggedUser } from "../network/timecard.api"

interface NavBarLoggedInViewProps {
  user: LoggedUser
  onLogoutSuccessful: () => void
}

const NavBarLoggedInView = ({
  user,
  onLogoutSuccessful,
}: NavBarLoggedInViewProps) => {
  async function logout() {
    try {
      await TimeCardApi.logout()
      sessionStorage.clear()
      localStorage.clear()
      localStorage.removeItem("eventId")

      onLogoutSuccessful()
    } catch (error) {
      console.log(error)
      alert(error)
    }
  }
  return (
    <>
      <Navbar.Text className="me-2">Signed in as: {user.userName}</Navbar.Text>
      <Button className="bg-dark" onClick={logout}>
        Logout
      </Button>
    </>
  )
}

export default NavBarLoggedInView

import { Button } from "react-bootstrap"

interface NavBarLoggedOutViewProps {
  onLoginClicked: () => void
}

const NavBarLoggedOutView = ({ onLoginClicked }: NavBarLoggedOutViewProps) => {
  return (
    <>
      <Button className="bg-dark" onClick={onLoginClicked}>
        Log In
      </Button>
    </>
  )
}

export default NavBarLoggedOutView

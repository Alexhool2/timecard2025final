import { useState, useEffect } from "react"
import { Container, Row, Col } from "react-bootstrap"
import Clock from "react-clock"
import "react-clock/dist/Clock.css"
import "../styles/global.css"
import StartTimeButton from "../components/StartTimeButton"
import StartLunchButton from "../components/StartLunchButton"
import EndLunchButton from "../components/EndLunchButton"
import EndTimeButton from "../components/EndTimeButton"
import { useEvent } from "../context/useEvent"
import { getLoggedInUser, LoggedUser } from "../network/timecard.api"

interface MainPageProps {
  user: LoggedUser
}

const MainPage = ({ user }: MainPageProps) => {
  const [time, setTime] = useState(new Date())
  const { eventId, setEventId } = useEvent()

  useEffect(() => {
    const timer = setInterval(() => {
      setTime(new Date())
    }, 1000)

    return () => clearInterval(timer)
  }, [])

  useEffect(() => {
    async function fetchEventId() {
      if (!eventId) {
        const storedEventId = localStorage.getItem("eventId")
        if (storedEventId) {
          setEventId(parseInt(storedEventId, 10))
        } else {
          // if cannot find on local storage check in backend
          try {
            const { eventID } = await getLoggedInUser()
            if (eventID !== null && eventID !== -1) {
              setEventId(eventID)
            }
          } catch (error) {
            console.error("Error fetching event ID:", error)
          }
        }
      }
    }
    fetchEventId()
  }, [eventId, setEventId])

  return (
    <Container className="text-center">
      <Row className="my-4">
        <Col>
          <h1>Welcome, {user.userName}</h1>
        </Col>
      </Row>
      <Row className="my-4 justify-content-center">
        <Col className="text-center">
          <div className="clock-container">
            <Clock value={time} />
          </div>
        </Col>
      </Row>
      <Row className="my-4">
        <Col className="d-grid gap-2 d-sm-flex justify-content-center">
          <StartTimeButton userId={user.id} />

          <StartLunchButton />

          <EndLunchButton />

          <EndTimeButton />
        </Col>
      </Row>
    </Container>
  )
}

export default MainPage

import { useState } from "react"
import { Button, Col, Row, Table } from "react-bootstrap"
import Calendar from "react-calendar"

const API_URL = import.meta.env.VITE_API_URL

interface EventData {
  eventId: number | null
  userId: number | null
  date: string | null
  startTime: string | null
  startLunch: string | null
  endLunch: string | null
  endTime: string | null
  userName: string | null
  firstName: string | null
  lastName: string | null
  role: string | null
}

const SearchSingleDate = ({ userId }: { userId: number | null }) => {
  const [selectedDate, setSelectedDate] = useState<Date>(new Date())
  const [events, setEvents] = useState<EventData | null>(null)
  const [hasSearched, setHasSearched] = useState(false)

  const fetchEvents = async () => {
    if (!userId || typeof userId !== "number") {
      setHasSearched(false)
      setEvents(null)
      return
    }
    try {
      const adjustedDate = new Date(selectedDate)
      adjustedDate.setMinutes(
        adjustedDate.getMinutes() - adjustedDate.getTimezoneOffset()
      )
      const dateStr = adjustedDate.toISOString().split("T")[0]

      const response = await fetch(
        `${API_URL}/api/v1/event/user/${userId}/date?date=${dateStr}`,
        {
          method: "GET",
          credentials: "include",
        }
      )
      setHasSearched(true)
      if (!response.ok) {
        setEvents(null)
        throw new Error(`Err ${response.status}: ${await response.text()}`)
      }
      const data = await response.json()
      if (data.user_id != userId) {
        setEvents(null)
        return
      }

      const formattedEvent: EventData = {
        eventId: data.event_id || null,
        userId: data.user_id || null,
        date: data.date || null,
        startTime: data.start_time || null,
        startLunch: data.start_lunch || null,
        endLunch: data.end_lunch || null,
        endTime: data.end_time || null,
        userName: data.user_name || null,
        firstName: data.first_name || null,
        lastName: data.last_name || null,
        role: data.role || null,
      }

      setEvents(formattedEvent)
    } catch {
      setEvents(null)
    }
  }

  const formatTime = (time: string | null) => {
    if (!time) return "Not registered"
    const date = new Date(time)
    return date.toISOString().split("T")[1].split(".")[0]
  }
  return (
    <>
      <Row className="justify-content-center mt-3">
        <Col
          md={6}
          className="d-flex flex-column align-items-center text-center"
        >
          <h5 className="text-center">Select single date</h5>
          <Calendar
            onChange={(date) => setSelectedDate(date as Date)}
            value={selectedDate}
          />
        </Col>
      </Row>

      <Row className="justify-content-center mt-3">
        <Col md={6} className="text-center">
          <Button variant="dark" onClick={fetchEvents} disabled={!userId}>
            Search
          </Button>
          {hasSearched && userId && events == null && (
            <div className="mt-3">
              <p className="text-danger">Event not found for this date</p>
            </div>
          )}
        </Col>
      </Row>

      {events && (
        <Row className="mt-4">
          <Col>
            <p className="fw-bold fs-5 text-dark">
              Date:{" "}
              {events.date
                ? new Date(events.date).toISOString().split("T")[0]
                : "No date"}
            </p>
          </Col>
        </Row>
      )}

      {events && (
        <Row className="justify-content-center mt-4">
          <Col md={8}>
            <Table striped bordered hover>
              <thead>
                <tr>
                  <th>Hour</th>
                  <th>Register</th>
                </tr>
              </thead>
              <tbody>
                <tr>
                  <td>Started</td>
                  <td>{formatTime(events.startTime)}</td>
                </tr>
                <tr>
                  <td>Started lunch</td>
                  <td>{formatTime(events.startLunch)}</td>
                </tr>
                <tr>
                  <td>Finished lunch</td>
                  <td>{formatTime(events.endLunch)}</td>
                </tr>
                <tr>
                  <td>Finished</td>
                  <td>{formatTime(events.endTime)}</td>
                </tr>
              </tbody>
            </Table>
          </Col>
        </Row>
      )}
    </>
  )
}

export default SearchSingleDate

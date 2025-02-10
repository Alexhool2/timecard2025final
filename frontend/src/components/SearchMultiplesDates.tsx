import { useState } from "react"
import { Button, Col, Row, Table } from "react-bootstrap"
import Calendar from "react-calendar"

const API_URL = import.meta.env.VITE_API_URL

interface EventData {
  event_id: number
  userId: number
  date: string
  start_time: string | null
  start_lunch: string | null
  end_lunch: string | null
  end_time: string | null
  userName: string
  firstName: string
  lastName: string
  role: string
}

const SearchMultiplesDates = ({ userId }: { userId: number | null }) => {
  const [selectedDates, setSelectedDates] = useState<[Date, Date] | []>([])
  const [events, setEvents] = useState<EventData[] | null>(null)
  const [hasSearched, setHasSearched] = useState(false)

  const adjustTimezone = (date: Date) => {
    const localDate = new Date(date)
    localDate.setMinutes(localDate.getMinutes() - localDate.getTimezoneOffset())
    return localDate.toISOString().split("T")[0]
  }

  const fetchEvents = async () => {
    if (!userId || selectedDates.length !== 2) {
      setHasSearched(false)
      return
    }

    try {
      const [startDate, endDate] = selectedDates
      const startDateStr = adjustTimezone(startDate)
      const endDateStr = adjustTimezone(endDate)

      const response = await fetch(
        `${API_URL}/api/v1/event/user/${userId}/period?start_date=${startDateStr}&end_date=${endDateStr}`,
        {
          method: "GET",
          credentials: "include",
        }
      )
      setHasSearched(true)
      if (!response.ok) {
        throw new Error(`Err ${response.status}: ${await response.text()}`)
      }
      const data = await response.json()

      setEvents(data.length ? data : null)
      if (data.length > 30) {
        alert("only possible to search in the range of 30 days")
        return
      }

      const tomorrow = new Date()
      tomorrow.setDate(tomorrow.getDate() + 1)
      tomorrow.setHours(0, 0, 0, 0)

      if (startDate > tomorrow || endDate > tomorrow) {
        alert("You cannot search for future dates, please try again")
        return
      }
    } catch (error) {
      console.error("Error finding events:", error)
      setEvents(null)
    }
  }

  const formatTime = (time: string | null) => {
    const date = new Date(time || "")
    if (isNaN(date.getTime())) return "Not registered"
    return date.toISOString().split("T")[1].split(".")[0]
  }
  return (
    <>
      <Row className="justify-content-center mt-3">
        <Col
          md={6}
          className="d-flex flex-column align-items-center text-center"
        >
          <h5 className="text-center">Select date range</h5>
          <Calendar
            selectRange={true}
            onChange={(date) => {
              if (Array.isArray(date) && date.length === 2) {
                setSelectedDates(date as [Date, Date])
              } else {
                setSelectedDates([])
              }
            }}
            value={selectedDates.length === 2 ? selectedDates : undefined}
          />
        </Col>
      </Row>

      <Row className="justify-content-center mt-3">
        <Col md={6} className="text-center">
          <Button
            variant="dark"
            onClick={fetchEvents}
            disabled={!userId || selectedDates.length !== 2}
          >
            Search
          </Button>
          {hasSearched && userId && events == null && (
            <div className="mt-3">
              <p className="text-danger">No events found for this period</p>
            </div>
          )}
        </Col>
      </Row>

      {events && selectedDates.length === 2 && (
        <Row className="mt-4">
          <Col>
            <p className="fw-bold fs-5 text-dark">
              Period: {adjustTimezone(selectedDates[0])} to{" "}
              {adjustTimezone(selectedDates[1])}
            </p>
          </Col>
        </Row>
      )}

      {events && events.length > 0 && (
        <Row className="justify-content-center mt-4">
          <Col md={10}>
            <Table striped bordered hover>
              <thead>
                <tr>
                  <th>Date</th>
                  <th>Started</th>
                  <th>Started lunch</th>
                  <th>Finished lunch</th>
                  <th>Finished</th>
                </tr>
              </thead>
              <tbody>
                {events.map((event, index) => (
                  <tr key={`event-${event.event_id}-${event.date}-${index}`}>
                    <td>{adjustTimezone(new Date(event.date))}</td>
                    <td>
                      {event.event_id === 0
                        ? "Absent"
                        : formatTime(event.start_time)}
                    </td>
                    <td>
                      {event.event_id === 0
                        ? "Absent"
                        : formatTime(event.start_lunch)}
                    </td>
                    <td>
                      {event.event_id === 0
                        ? "Absent"
                        : formatTime(event.end_lunch)}
                    </td>
                    <td>
                      {event.event_id === 0
                        ? "Absent"
                        : formatTime(event.end_time)}
                    </td>
                  </tr>
                ))}
              </tbody>
            </Table>
          </Col>
        </Row>
      )}
    </>
  )
}

export default SearchMultiplesDates

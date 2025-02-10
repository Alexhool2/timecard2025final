import { useState } from "react"
import { Container, Row, Col } from "react-bootstrap"
import "react-calendar/dist/Calendar.css"
import UserSearch from "../components/UserSearch"
import SearchSingleDate from "../components/SearchSingleDate"
import SearchMultiplesDates from "../components/SearchMultiplesDates"

const TimeReportPage = () => {
  const [userId, setUserId] = useState<number | null>(null)

  return (
    <Container className="mt-4">
      <Row className="justify-content-center">
        <Col md={6}>
          <h2 className="text-center">Select user</h2>
          <UserSearch onSelectUser={(id) => setUserId(id)} />
        </Col>
      </Row>
      <SearchSingleDate userId={userId} />
      <SearchMultiplesDates userId={userId} />
    </Container>
  )
}

export default TimeReportPage

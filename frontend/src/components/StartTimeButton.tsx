import { useState } from "react"
import { useEvent } from "../context/useEvent"
import { Button } from "react-bootstrap"
import ConfirmModal from "./ConfirmModal"

const API_URL = import.meta.env.VITE_API_URL

interface StartTimeProps {
  userId: number
}

const StartTimeButton = ({ userId }: StartTimeProps) => {
  const { setEventId } = useEvent()
  const [showModal, setShowModal] = useState(false)

  const handleStart = async () => {
    try {
      const response = await fetch(`${API_URL}/api/v1/event/create`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ userID: userId }),
        credentials: "include",
      })

      if (!response.ok) {
        const errorMessage = await response.json()
        console.error("Backend error:", errorMessage)

        if (
          errorMessage.error &&
          errorMessage.error.includes("Duplicate entry")
        ) {
          alert("Start time already registered.")
          return
        }
      }

      const data = await response.json()

      setEventId(data.eventId)

      setTimeout(() => {
        alert(`Event created successfully! Start time: ${data.startTime}`)
      }, 500)
    } catch (error) {
      console.error("Error submitting event:", error)
      alert("An unexpected error occurred. Please try again.")
    }
  }

  return (
    <>
      <Button className="bg-dark mx-1" onClick={() => setShowModal(true)}>
        Start time
      </Button>
      <ConfirmModal
        show={showModal}
        onHide={() => setShowModal(false)}
        onConfirm={() => {
          setShowModal(false)
          handleStart()
        }}
        title="Confirm Start Time"
        body="Do you confirm starting time?"
      ></ConfirmModal>
    </>
  )
}

export default StartTimeButton

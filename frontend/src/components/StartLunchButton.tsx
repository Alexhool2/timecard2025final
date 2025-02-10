import { useState } from "react"
import { useEvent } from "../context/useEvent"
import { Button } from "react-bootstrap"
import ConfirmModal from "./ConfirmModal"

const API_URL = import.meta.env.VITE_API_URL

const StartLunchButton = () => {
  const { eventId } = useEvent()
  const [showModal, setShowModal] = useState(false)

  const handleStartLunch = async () => {
    if (!eventId) {
      console.log(eventId)

      alert("Did you already punch start time?.")
      return
    }
    try {
      const response = await fetch(
        `${API_URL}/api/v1/event/start-lunch/${eventId}`,
        {
          method: "PUT",
          headers: {
            "Content-Type": "application/json",
          },

          credentials: "include",
        }
      )

      if (!response.ok) {
        const errorMessage = await response.json()
        console.error("Backend error:", errorMessage)

        if (
          errorMessage.error &&
          errorMessage.error.includes("already exists")
        ) {
          alert("Start lunch already registered.")
          return
        }
      }
      const data = await response.json()

      setTimeout(() => {
        alert(`Lunch time started successfully! Lunch time: ${data.startLunch}`)
      }, 500)
    } catch (error) {
      console.error("Error submitting event:", error)
      alert("An unexpected error occurred. Please try again.")
    }
  }
  return (
    <>
      <Button className="bg-dark mx-1" onClick={() => setShowModal(true)}>
        Start lunch
      </Button>
      <ConfirmModal
        show={showModal}
        onHide={() => setShowModal(false)}
        onConfirm={() => {
          setShowModal(false)
          handleStartLunch()
        }}
        title="Confirm Start Lunch"
        body="Do you confirm starting lunch?"
      ></ConfirmModal>
    </>
  )
}

export default StartLunchButton

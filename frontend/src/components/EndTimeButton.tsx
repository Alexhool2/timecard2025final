import { useState } from "react"
import { useEvent } from "../context/useEvent"
import { Button } from "react-bootstrap"
import ConfirmModal from "./ConfirmModal"

const API_URL = import.meta.env.VITE_API_URL

const EndTimeButton = () => {
  const { eventId } = useEvent()
  const [showModal, setShowModal] = useState(false)

  const handleEndTime = async () => {
    if (!eventId) {
      console.log("evento nao encontrado")
    }
    try {
      const response = await fetch(
        `${API_URL}/api/v1/event/end-time/${eventId}`,
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
          alert("End time already registered.")
          return
        }
      }
      const data = await response.json()
      setTimeout(() => {
        alert(`End Time ended successfully! time: ${data.endTime}`)
      }, 500)
    } catch (error) {
      console.error("Error submitting event:", error)
      alert("An unexpected error occurred. Please try again.")
    }
  }

  return (
    <>
      <Button className="bg-dark mx-1" onClick={() => setShowModal(true)}>
        End time
      </Button>
      <ConfirmModal
        show={showModal}
        onHide={() => setShowModal(false)}
        onConfirm={() => {
          setShowModal(false)
          handleEndTime()
        }}
        title="Confirm End time"
        body="Do you confirm finishing work?"
      ></ConfirmModal>
    </>
  )
}

export default EndTimeButton

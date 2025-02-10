import { useState } from "react"
import { useEvent } from "../context/useEvent"
import { Button } from "react-bootstrap"
import ConfirmModal from "./ConfirmModal"

const API_URL = import.meta.env.VITE_API_URL

const EndLunchButton = () => {
  const { eventId } = useEvent()
  const [showModal, setShowModal] = useState(false)

  const handleEndLunch = async () => {
    if (!eventId) {
      alert("Did you already punch start time?.")
      return
    }

    try {
      const response = await fetch(
        `${API_URL}/api/v1/event/end-lunch/${eventId}`,
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
          alert("End lunch already registered.")
          return
        } else if (
          errorMessage.error.includes("start_lunch needs to be first")
        ) {
          alert("Did you punch start lunch?. Could not find your data")
        }
      }
      const data = await response.json()
      setTimeout(() => {
        alert(`Lunch ended successfully! End lunch time: ${data.endLunch}`)
      }, 500)
    } catch (error) {
      console.error("Error submitting event:", error)
      alert("An unexpected error occurred. Please try again.")
    }
  }
  return (
    <>
      <Button className="bg-dark mx-1" onClick={() => setShowModal(true)}>
        End lunch
      </Button>
      <ConfirmModal
        show={showModal}
        onHide={() => setShowModal(false)}
        onConfirm={() => {
          setShowModal(false)
          handleEndLunch()
        }}
        title="Confirm End Lunch"
        body="Do you confirm ending your lunch time?"
      ></ConfirmModal>
    </>
  )
}

export default EndLunchButton

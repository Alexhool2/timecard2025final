import { useContext } from "react"
import { EventContext } from "./EventProvider"

export const useEvent = () => {
  const context = useContext(EventContext)
  if (!context) {
    throw new Error("useEvent must be used within an EventProvider")
  }
  return context
}

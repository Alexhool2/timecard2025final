import React, { createContext, useState, useEffect } from "react"

interface EventContextType {
  eventId: number | null
  setEventId: (id: number | null) => void
}

export const EventContext = createContext<EventContextType | undefined>(
  undefined
)

export const EventProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const [eventId, setEventId] = useState<number | null>(() => {
    // Inicializa o eventId a partir do Local Storage
    const storedEventId = localStorage.getItem("eventId")
    return storedEventId ? parseInt(storedEventId, 10) : null
  })

  useEffect(() => {
    // Salva o eventId no Local Storage sempre que ele mudar
    if (eventId !== null && eventId !== -1) {
      localStorage.setItem("eventId", eventId.toString())
    } else {
      localStorage.removeItem("eventId")
    }
  }, [eventId])

  return (
    <EventContext.Provider value={{ eventId, setEventId }}>
      {children}
    </EventContext.Provider>
  )
}

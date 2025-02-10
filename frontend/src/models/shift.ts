export interface Shift {
  id: number
  userId: number
  date: string // time.time mudar para iso
  startTime?: string
  startLunch?: string
  endLunch?: string
  endTime?: string
  createdAt: string
}

// originalmente esta desta forma em go

// type Event struct {
// 	ID          int          `json:"id"`
// 	UserID      int          `json:"userId"`
// 	Date        time.Time    `json:"date"`
// 	Start_time  sql.NullTime `json:"start_time"`
// 	Start_lunch sql.NullTime `json:"start_lunch"`
// 	End_lunch   sql.NullTime `json:"end_lunch"`
// 	End_time    sql.NullTime `json:"end_time"`
// 	CreatedAt   string       `json:"createdAt"`
// }

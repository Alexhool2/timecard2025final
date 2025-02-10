package types

import (
	"database/sql"
	"encoding/json"
	"time"
)

type UserStore interface {
	GetUserByUserName(userName string) (*User, error)
	GetAllUsers() ([]User, error)
	GetUserById(id int) (*User, error)
	CreateUser(User) error
	UpdatePassword(userID int, hashedPassword string) error
	GetDynamicUsers(search string) ([]User, error)
}
type EventStore interface {
	CreateEvent(CreateEventPayload) (int, error)
	GetEventById(id int) (*Event, error)
	GetEventByDate(date time.Time, userId int) (*EventWithUser, error)
	GetEventBySelectedDates(userId int, startDate, endDate time.Time) ([]EventWithUser, error)
	UpdateStartLunchEvent(id int, startLunch time.Time) error
	UpdateEndLunchEvent(id int, endLunch time.Time) error
	UpdateEndTimeEvent(id int, endLunch time.Time) error
	GetStartLunch(id int) (*Event, error)
	GetEndLunch(id int) (*Event, error)
	GetEndTime(id int) (*Event, error)
	GetTodayEvent(userID int) (eventID int, err error)
}

type EventID struct {
	ID int `json:"id"`
}

type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	UserName  string `json:"userName"`
	Password  string `json:"-"`
	Email     string `json:"email"`
	IsAdmin   bool   `json:"isAdmin"`
	CreatedAt string `json:"createdAt"`
	Role      Role   `json:"role"`
}

type Event struct {
	ID          int          `json:"id"`
	UserID      int          `json:"userId"`
	Date        time.Time    `json:"date"`
	Start_time  sql.NullTime `json:"start_time"`
	Start_lunch sql.NullTime `json:"start_lunch"`
	End_lunch   sql.NullTime `json:"end_lunch"`
	End_time    sql.NullTime `json:"end_time"`
	CreatedAt   string       `json:"createdAt"`
}

type EventWithUser struct {
	EventID    int        `json:"event_id"`
	UserID     int        `json:"user_id"`
	Date       time.Time  `json:"date"`
	StartTime  *time.Time `json:"start_time"`
	StartLunch *time.Time `json:"start_lunch"`
	EndLunch   *time.Time `json:"end_lunch"`
	EndTime    *time.Time `json:"end_time"`
	Role       Role       `json:"role"`
	UserName   string     `json:"user_name"`
	FirstName  string     `json:"first_name"`
	LastName   string     `json:"last_name"`
}

type LoginUserPayload struct {
	UserName string `json:"userName" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type CreateEventPayload struct {
	UserID     int       `json:"userId" validate:"required"`
	Date       time.Time `json:"date"`
	Start_time time.Time `json:"start_time"`
}

type RegisterUserPayload struct {
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	UserName  string `json:"userName" validate:"required"`
	Password  string `json:"password" validate:"required,min=6,max=130"`
	Email     string `json:"email" validate:"required,email"`
	Role      Role   `json:"role" validate:"required,role"`
}

type ResetPassword struct {
	NewPassword string `json:"newPassword" validate:"required,min=6,max=130"`
}

type Role string

type NullTime struct {
	sql.NullTime
}

func (nt NullTime) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(nt.Time)
}

func (nt *NullTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		nt.Valid = false
		return nil
	}
	if err := json.Unmarshal(data, &nt.Time); err != nil {
		return err
	}
	nt.Valid = true
	return nil
}

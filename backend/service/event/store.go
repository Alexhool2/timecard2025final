package event

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/alexhool2/TimeCard/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateEvent(event types.CreateEventPayload) (int, error) {
	query := `
	INSERT INTO event (userId, date, start_time)
	VALUES (?, ?, ?)
`

	result, err := s.db.Exec(query,
		event.UserID,
		event.Date.Format("2006-01-02"),
		event.Start_time.Format("2006-01-02T15:04:05"),
	)
	if err != nil {
		log.Printf("Database error %v", err)
		return 0, fmt.Errorf("failed to insert event: %w", err)
	}
	eventID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to retrive event Id: %v", err)
	}
	return int(eventID), nil
}

// Query to get start lunch event
func (s *Store) GetStartLunch(id int) (*types.Event, error) {
	query := `
		SELECT id, userId, date, start_lunch FROM event WHERE id = ?
	`

	row := s.db.QueryRow(query, id)
	var event types.Event
	err := row.Scan(
		&event.ID,
		&event.UserID,
		&event.Date,
		&event.Start_lunch,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("event with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to get event: %v", err)
	}

	if event.Start_lunch.Valid {
		return nil, fmt.Errorf("start_lunch already exists for event with id %d", id)
	}

	return &event, nil
}

// update startLunch query
func (s *Store) UpdateStartLunchEvent(id int, startLunch time.Time) error {
	formattedTime := startLunch.Format("2006-01-02T15:04:05")
	query := `UPDATE event SET start_lunch = ? WHERE id = ?`
	_, err := s.db.Exec(query, formattedTime, id)
	if err != nil {
		return fmt.Errorf("failed to update start_lunch for event with id %d: %v", id, err)
	}
	return nil
}

func (s *Store) GetEndLunch(id int) (*types.Event, error) {
	query := `
		SELECT id, userId, date, start_lunch,end_lunch FROM event WHERE id = ?
	`

	row := s.db.QueryRow(query, id)
	var event types.Event
	err := row.Scan(
		&event.ID,
		&event.UserID,
		&event.Date,
		&event.Start_lunch,
		&event.End_lunch,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("event with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to get event: %v", err)
	}

	if !event.Start_lunch.Valid {
		return nil, fmt.Errorf("start_lunch needs to be first")
	}

	if event.End_lunch.Valid {
		return nil, fmt.Errorf("end_lunch already exists for event with id %d", id)
	}

	return &event, nil
}

func (s *Store) UpdateEndLunchEvent(id int, endLunch time.Time) error {
	formattedTime := endLunch.Format("2006-01-02T15:04:05")
	query := `UPDATE event SET end_lunch = ? WHERE id = ?`
	_, err := s.db.Exec(query, formattedTime, id)
	if err != nil {
		return fmt.Errorf("failed to update end_lunch for event with id %d: %v", id, err)
	}
	return nil
}

func (s *Store) GetEndTime(id int) (*types.Event, error) {
	query := `
		SELECT id, userId, date, start_time,end_time FROM event WHERE id = ?
	`

	row := s.db.QueryRow(query, id)
	var event types.Event
	err := row.Scan(
		&event.ID,
		&event.UserID,
		&event.Date,
		&event.Start_time,
		&event.End_time,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("event with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to get event: %v", err)
	}

	if !event.Start_time.Valid {
		return nil, fmt.Errorf("start_time needs to be first")
	}

	if event.End_time.Valid {
		return nil, fmt.Errorf("end_time already exists for event with id %d", id)
	}

	return &event, nil
}

func (s *Store) UpdateEndTimeEvent(id int, endTime time.Time) error {
	formattedTime := endTime.Format("2006-01-02T15:04:05")

	query := `UPDATE event SET end_time = ? WHERE id = ?`
	_, err := s.db.Exec(query, formattedTime, id)
	if err != nil {
		return fmt.Errorf("failed to update end_time for event with id %d: %v", id, err)
	}

	return nil
}

func (s *Store) GetEventById(id int) (*types.Event, error) {
	query := `
		SELECT id, userId, date FROM event WHERE id = ?
	`

	row := s.db.QueryRow(query, id)
	var event types.Event
	err := row.Scan(
		&event.ID,
		&event.UserID,
		&event.Date,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("event with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to get event: %v", err)
	}
	return &event, nil
}

func (s *Store) GetEventByDate(date time.Time, userId int) (*types.EventWithUser, error) {
	formattedDate := date.Format("2006-01-02")
	query :=
		`
	SELECT
	 e.id as event_id,
	e.userId,
	e.date,
	e.start_time,
	e.start_lunch,
	e.end_lunch,
	e.end_time,
	u.role as user_role,
	u.firstName as user_first_name,
	u.lastName as user_last_name,
	u.userName as user_name
	
	FROM
	event e
	JOIN users u ON e.userId = u.id
     WHERE
	e.date = ? and e.userId = ?
	`

	row := s.db.QueryRow(query, formattedDate, userId)
	var event types.EventWithUser

	err := row.Scan(
		&event.EventID,
		&event.UserID,
		&event.Date,
		&event.StartTime,
		&event.StartLunch,
		&event.EndLunch,
		&event.EndTime,
		&event.Role,
		&event.FirstName,
		&event.LastName,
		&event.UserName,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no event found for this date %v", date)
		}
		return nil, fmt.Errorf("failed to get event %v", err)
	}
	return &event, nil

}

func (s *Store) GetEventBySelectedDates(userId int, startDate, endDate time.Time) ([]types.EventWithUser, error) {
	formattedStartDate := startDate.Format("2006-01-02")
	formattedEndDate := endDate.Format("2006-01-02")

	query := `
	SELECT
		e.id as event_id,
		e.userId,
		e.date,
		e.start_time,
		e.start_lunch,
		e.end_lunch,
		e.end_time,
		u.role as user_role,
		u.firstName as user_first_name,
		u.lastName as user_last_name,
		u.userName as user_name
	FROM event e
	JOIN users u ON e.userId = u.id
	WHERE e.date BETWEEN ? AND ? AND e.userId = ?
	ORDER BY e.date
	`

	rows, err := s.db.Query(query, formattedStartDate, formattedEndDate, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to query events: %v", err)
	}
	defer rows.Close()

	// Create a map to facilitate event lookup by date
	eventsMap := make(map[string]types.EventWithUser)

	// List of events returned from database
	var events []types.EventWithUser

	for rows.Next() {
		var event types.EventWithUser
		err := rows.Scan(
			&event.EventID,
			&event.UserID,
			&event.Date,
			&event.StartTime,
			&event.StartLunch,
			&event.EndLunch,
			&event.EndTime,
			&event.Role,
			&event.FirstName,
			&event.LastName,
			&event.UserName,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan event: %v", err)
		}

		eventsMap[event.Date.Format("2006-01-02")] = event
	}

	// Create final list of events, including "Absent" days
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format("2006-01-02")

		// If event exists for the date, add it to the list
		if event, found := eventsMap[dateStr]; found {
			events = append(events, event)
		} else {
			// Otherwise, add an "Absent" event
			events = append(events, types.EventWithUser{
				EventID:    0, // ID 0 indicates a placeholder event
				UserID:     userId,
				Date:       d,
				StartTime:  nil,
				StartLunch: nil,
				EndLunch:   nil,
				EndTime:    nil,
				Role:       "",
				FirstName:  "",
				LastName:   "",
				UserName:   "",
			})
		}
	}

	return events, nil
}

func (s *Store) GetTodayEvent(userID int) (eventID int, err error) {

	query := `SELECT id FROM event WHERE userId = ? ORDER BY date DESC LIMIT 1`
	err = s.db.QueryRow(query, userID).Scan(&eventID)
	if err != nil {
		return 0, fmt.Errorf("failed to get event for user %d: %v", userID, err)
	}
	return eventID, nil
}

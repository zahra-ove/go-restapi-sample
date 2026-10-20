package models

import (
	"time"

	"github.com/zahra-ove/go-restapi-sample/db"
)

type Event struct {
	ID          int64
	Name        string    `binding:"required"`
	Description string    `binding:"required"`
	Location    string    `binding:"required"`
	DateTime    time.Time `binding:"required"`
	UserID      int64
}

func (e *Event) Save() error {
	query := `INSERT INTO events(name, description, location, datetime, user_id) VALUES (?, ?, ?, ?, ?)`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(e.Name, e.Description, e.Location, e.DateTime, e.UserID)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return err
	}

	e.ID = id
	return err
}

func GetAllEvents() ([]Event, error) {
	rows, err := db.DB.Query(`SELECT id, name, description, location, datetime, user_id FROM events ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := []Event{}
	for rows.Next() {
		var event Event
		if err := rows.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserID); err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, rows.Err()
}

func GetEvent(eventId int64) (*Event, error) {
	row := db.DB.QueryRow("SELECT id, name, description, location, datetime, user_id FROM events WHERE id = ?", eventId)

	var event Event
	err := row.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserID)

	if err != nil {
		return nil, err
	}

	return &event, nil
}

func (event Event) Update() error {
	query := `UPDATE events 
				SET name = ?, description = ?, location = ?, dateTime = ? 
				WHERE id = ?`

	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(event.Name, event.Description, event.Location, event.DateTime, event.ID)
	return err
}

func (event Event) Delete() error {
	query := `DELETE FROM events 
				WHERE id = ?`

	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(event.ID)
	return err
}

func (event Event) Register(userId int64) error {
	query := `INSERT INTO registrations(event_id, user_id) VALUES (?, ?)`
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	_, err = stmt.Exec(event.ID, userId)
	if err != nil {
		return err
	}

	return nil
}

func (event Event) Cancel(userId int64) error {
	query := `DELETE FROM registrations WHERE event_id = ? AND user_id = ?`
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	_, err = stmt.Exec(event.ID, userId)
	if err != nil {
		return err
	}

	return nil
}

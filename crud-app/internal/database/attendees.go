package database

import (
	"context"
	"database/sql"
	"time"
)

type AttendeeModel struct {
	DB *sql.DB
}

type Attendee struct {
	Id       int    `json:"id"`
	UserId   int    `json:"userId"`
	EventId  int    `json:"eventId"`
}

func (m *AttendeeModel) Insert(attendee *Attendee) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	INSERT INTO attendees (user_id, event_id)
	VALUES ($1, $2)
	RETURNING id
	`

	err := m.DB.QueryRowContext(ctx, query, attendee.UserId, attendee.EventId).Scan(&attendee.Id)
	if err != nil {
		return err
	}

	return nil
}

func (m *AttendeeModel) GetByEventAndAttendee(eventId int, userId int) (*Attendee, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	SELECT id, user_id, event_id
	FROM attendees
	WHERE event_id = $1 AND user_id = $2
	`

	var attendee Attendee
	err := m.DB.QueryRowContext(ctx, query, eventId, userId).Scan(&attendee.Id, &attendee.UserId, &attendee.EventId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &attendee, nil
}

func (m *AttendeeModel) GetAttendeesByEvent(eventId int) ([]*Attendee, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	SELECT id, user_id, event_id
	FROM attendees
	WHERE event_id = $1
	`

	rows, err := m.DB.QueryContext(ctx, query, eventId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	attendees := []*Attendee{}
	
	for rows.Next() {
		var attendee Attendee
		err := rows.Scan(&attendee.Id, &attendee.UserId, &attendee.EventId)
		if err != nil {
			return nil, err
		}
		attendees = append(attendees, &attendee)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return attendees, nil
}

func (m *AttendeeModel) Delete(eventId int, userId int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	DELETE FROM attendees
	WHERE event_id = $1 AND user_id = $2
	`

	_, err := m.DB.ExecContext(ctx, query, eventId, userId)
	if err != nil {
		return err
	}

	return nil
}

func (m *AttendeeModel) GetEventsByAttendee(userId int) ([]*Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	SELECT e.id, e.title, e.description, e.date
	FROM events e
	INNER JOIN attendees a ON e.id = a.event_id
	WHERE a.user_id = $1
	`

	rows, err := m.DB.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := []*Event{}
	
	for rows.Next() {
		var event Event
		err := rows.Scan(&event.Id, &event.Name, &event.Description, &event.Date)
		if err != nil {
			return nil, err
		}
		events = append(events, &event)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}	

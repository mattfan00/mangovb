package vb

import "time"

type Event struct {
	Id       string `db:"id"`
	Name     string `db:"name"`
	Location string `db:"location"`
	// TODO: use a string here and make start_date a DATE type in DB
	StartDate   time.Time `db:"start_date"`
	StartTime   string    `db:"start_time"`
	EndTime     string    `db:"end_time"`
	SkillLevel  string
	Price       float64   `db:"price"`
	IsAvailable bool      `db:"is_available"`
	SpotsLeft   int       `db:"spots_left"`
	Url         string    `db:"url"`
	Source      string    `db:"source"`
	UpdatedOn   time.Time `db:"updated_on"`
}

type EventNotifType int

const (
	NewEvent EventNotifType = iota
	LimitedSpots
)

type EventNotif struct {
	Type      EventNotifType `db:"type"`
	EventId   string         `db:"event_id"`
	Event     Event
	CreatedOn time.Time `db:"created_on"`
}

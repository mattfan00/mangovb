package store

import (
	"github.com/jmoiron/sqlx"
	vb "github.com/mattfan00/nycvbtracker"
)

type EventStore struct {
	db *sqlx.DB
}

func NewEventStore(db *sqlx.DB) *EventStore {
	return &EventStore{
		db: db,
	}
}

func (es *EventStore) Insert(event vb.Event) error {
	stmt := "INSERT INTO event VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"

	_, err := es.db.Exec(stmt,
		event.Id,
		event.Source,
		event.Name,
		event.Location,
		event.StartDate,
		event.StartTime,
		event.EndTime,
		event.Price,
		event.IsAvailable,
		event.SpotsLeft,
		event.Url,
	)

	return err
}

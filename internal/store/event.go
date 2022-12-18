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

func (es *EventStore) InsertMultiple(events []vb.Event) error {
	tx, err := es.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt := "INSERT INTO event VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"

	for _, event := range events {
		_, err = tx.Exec(stmt,
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
		if err != nil {
			return err
		}
	}

	err = tx.Commit()

	return err
}

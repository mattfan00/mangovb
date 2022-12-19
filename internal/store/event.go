package store

import (
	vb "github.com/mattfan00/nycvbtracker"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
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
	baseInsert := sq.
		Insert("event").
		Columns("id", "source", "name", "location", "start_date", "start_time", "end_time", "price", "is_available", "spots_left", "url", "scraped_on")

	tx, err := es.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, event := range events {
		stmt, args, err := baseInsert.Values(
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
			event.ScrapedOn,
		).ToSql()
		if err != nil {
			return err
		}

		_, err = tx.Exec(stmt, args...)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()

	return err
}

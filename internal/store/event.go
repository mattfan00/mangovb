package store

import (
	"fmt"

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

func (es *EventStore) UpsertMultiple(events []vb.Event) error {
	baseInsert := sq.
		Insert("event").
		Columns("id", "source", "name", "location", "start_date", "start_time", "end_time", "price", "is_available", "spots_left", "url", "updated_on").
		Suffix(fmt.Sprintf("ON CONFLICT(id) DO UPDATE SET %s", setMap(map[string]interface{}{
			"name":         "excluded.name",
			"location":     "excluded.location",
			"start_date":   "excluded.start_date",
			"start_time":   "excluded.start_time",
			"end_time":     "excluded.end_time",
			"price":        "excluded.price",
			"is_available": "excluded.is_available",
			"spots_left":   "excluded.spots_left",
			"updated_on":   "excluded.updated_on",
		})))

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
			event.UpdatedOn,
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

func (es *EventStore) GetLatest() ([]vb.Event, error) {
	subquery := sq.Select().
		Column("*").
		Column("ROW_NUMBER() OVER (PARTITION BY id ORDER BY updated_on DESC) AS rn").
		From("event")

	stmt, args, err := sq.Select().
		Columns("id", "source", "name", "location", "start_date", "start_time", "end_time", "price", "is_available", "spots_left", "url", "updated_on").
		FromSelect(subquery, "t").
		Where(sq.And{
			sq.Eq{"rn": 1},
		}).
		ToSql()
	if err != nil {
		return []vb.Event{}, err
	}

	events := []vb.Event{}
	err = es.db.Select(&events, stmt, args...)

	return events, err
}

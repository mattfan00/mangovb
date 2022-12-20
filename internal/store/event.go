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

func (es *EventStore) GetLatestByIds(ids []string) (map[string]vb.Event, error) {
	subquery := sq.Select().
		Column("*").
		Column("ROW_NUMBER() OVER (PARTITION BY ID ORDER BY scraped_on DESC) AS rn").
		From("event")

	stmt, args, err := sq.Select().
		Columns("id", "source", "name", "location", "start_date", "start_time", "end_time", "price", "is_available", "spots_left", "url").
		FromSelect(subquery, "t").
		Where(sq.And{
			sq.Eq{"rn": 1},
			sq.Eq{"id": ids},
		}).
		ToSql()
	if err != nil {
		return map[string]vb.Event{}, err
	}

	events := []vb.Event{}
	err = es.db.Select(&events, stmt, args...)
	if err != nil {
		return map[string]vb.Event{}, err
	}

	eventMap := map[string]vb.Event{}
	for _, event := range events {
		eventMap[event.Id] = event
	}

	return eventMap, nil
}

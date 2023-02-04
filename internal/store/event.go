package store

import (
	"fmt"
	"time"

	vb "github.com/mattfan00/mangovb"

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
	updatedOn := time.Now()

	baseInsert := sq.
		Insert("event").
		Columns("id", "source_id", "name", "location", "start_time", "end_time", "skill_level", "price", "is_available", "spots_left", "url", "updated_on").
		Suffix(fmt.Sprintf("ON CONFLICT(id) DO UPDATE SET %s", setMap(map[string]interface{}{
			"name":         "excluded.name",
			"location":     "excluded.location",
			"start_time":   "excluded.start_time",
			"end_time":     "excluded.end_time",
			"skill_level":  "excluded.skill_level",
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
			event.SourceId,
			event.Name,
			event.Location,
			event.StartTime,
			event.EndTime,
			event.SkillLevel,
			event.Price,
			event.IsAvailable,
			event.SpotsLeft,
			event.Url,
			updatedOn,
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

func (es *EventStore) GetLatest(sort bool) ([]vb.Event, error) {
    baseSelect := sq.Select().
		Columns("id", "source_id", "name", "location", "start_time", "end_time", "skill_level", "price", "is_available", "spots_left", "url", "updated_on").
        From("event").
        Where("updated_on = (SELECT MAX(updated_on) FROM event)")

    if sort {
        baseSelect = baseSelect.OrderBy("start_time, name")
    }

    stmt, args, err := baseSelect.ToSql()

	if err != nil {
		return []vb.Event{}, err
	}

	events := []vb.Event{}
	err = es.db.Select(&events, stmt, args...)

	return events, err
}

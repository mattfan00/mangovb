package store

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	vb "github.com/mattfan00/mangovb"
)

type EventNotifStore struct {
	db *sqlx.DB
}

func NewEventNotifStore(db *sqlx.DB) *EventNotifStore {
	return &EventNotifStore{
		db: db,
	}
}

func (ens *EventNotifStore) InsertMultiple(notifs []vb.EventNotif) error {
	baseInsert := sq.
		Insert("event_notif").
		Columns("event_id", "type_id", "created_on")

	tx, err := ens.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, notif := range notifs {
		stmt, args, err := baseInsert.Values(
			notif.EventId,
			notif.TypeId,
			notif.CreatedOn,
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

func (ens *EventNotifStore) GetByEventIds(ids []string) (map[string][]vb.EventNotif, error) {
	stmt, args, err := sq.
		Select("event_id", "type_id", "created_on").
		From("event_notif").
		Where(sq.Eq{"event_id": ids}).
		ToSql()
	if err != nil {
		return map[string][]vb.EventNotif{}, err
	}

	foundNotifs := []vb.EventNotif{}
	err = ens.db.Select(&foundNotifs, stmt, args...)
	if err != nil {
		return map[string][]vb.EventNotif{}, err
	}

	notifMap := map[string][]vb.EventNotif{}
	for _, notif := range foundNotifs {
		notifMap[notif.EventId] = append(notifMap[notif.EventId], notif)
	}

	return notifMap, nil
}

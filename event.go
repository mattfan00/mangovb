package vb

import "time"

type Event struct {
	Id          string          `db:"id"`
	Name        string          `db:"name"`
	Location    string          `db:"location"`
	StartDate   time.Time       `db:"start_date"`
	StartTime   string          `db:"start_time"`
	EndTime     string          `db:"end_time"`
	SkillLevel  EventSkillLevel `db:"skill_level"`
	Price       float64         `db:"price"`
	IsAvailable bool            `db:"is_available"`
	SpotsLeft   int             `db:"spots_left"`
	Url         string          `db:"url"`
	SourceId    EventSource     `db:"source_id"`
	UpdatedOn   time.Time       `db:"updated_on"`
}

type EventSource int

const (
	NyUrban EventSource = iota
)

var EventSourceMap = map[EventSource]string{
	NyUrban: "NY Urban",
}

type EventSkillLevel int

const (
	None EventSkillLevel = iota
	Beginner
	Intermediate
	AdvancedIntermediate
	Advanced
)

type EventNotifType int

const (
	NewEvent EventNotifType = iota
	LimitedSpots
)

type EventNotif struct {
	TypeId    EventNotifType `db:"type_id"`
	EventId   string         `db:"event_id"`
	Event     Event
	CreatedOn time.Time `db:"created_on"`
}

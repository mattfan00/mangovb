package vb

import "time"

type Event struct {
	Id          string          `db:"id" json:"id"`
	Name        string          `db:"name" json:"name"`
	Location    string          `db:"location" json:"location"`
	StartTime   time.Time       `db:"start_time" json:"startTime"`
	EndTime     time.Time       `db:"end_time" json:"endTime"`
	SkillLevel  EventSkillLevel `db:"skill_level" json:"skillLevel"`
	Price       float64         `db:"price" json:"price"`
	IsAvailable bool            `db:"is_available" json:"isAvailable"`
	SpotsLeft   int             `db:"spots_left" json:"spotsLeft"`
	Url         string          `db:"url" json:"url"`
	SourceId    EventSource     `db:"source_id" json:"sourceId"`
	UpdatedOn   time.Time       `db:"updated_on" json:"updatedOn"`
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
